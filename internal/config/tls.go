package config

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"os"
)

type TLSConfigError struct {
	inner error
}

func (e *TLSConfigError) Error() string {
	return fmt.Sprintf("TLS config: %s", e.inner.Error())
}

func NewTLSConfig(params TlsCfgParameters) (*tls.Config, error) {
	rootCAs := x509.NewCertPool()
	clientCAs := x509.NewCertPool()

	brokerCert, err := os.ReadFile(params.BrokerCertFile)
	if err != nil {
		return nil, &TLSConfigError{err}
	}
	clientCert, err := os.ReadFile(params.ClientCA)
	if err != nil {
		return nil, &TLSConfigError{err}
	}

	rootCAs.AppendCertsFromPEM(brokerCert)
	clientCAs.AppendCertsFromPEM(clientCert)

	cert, err := tls.LoadX509KeyPair(params.ClientCertFile, params.ClientPrivateKey)
	if err != nil {
		return nil, &TLSConfigError{err}
	}

	cert.Leaf, err = x509.ParseCertificate(cert.Certificate[0])
	if err != nil {
		return nil, &TLSConfigError{err}
	}

	return &tls.Config{
		RootCAs:      rootCAs,
		ClientAuth:   tls.RequireAndVerifyClientCert,
		ClientCAs:    clientCAs,
		Certificates: []tls.Certificate{cert},
	}, nil

}
