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

func NewTLSConfig(brokerCertFile string) (*tls.Config, error) {
	rootCAs := x509.NewCertPool()
	clientCAs := x509.NewCertPool()

	brokerCert, err := os.ReadFile(brokerCertFile)
	if err != nil {
		return nil, &TLSConfigError{err}
	}
	clientCert, err := os.ReadFile("certs/self-signed.crt")
	if err != nil {
		return nil, &TLSConfigError{err}
	}

	rootCAs.AppendCertsFromPEM(brokerCert)
	clientCAs.AppendCertsFromPEM(clientCert)

	cert, err := tls.LoadX509KeyPair("certs/client.crt", "certs/client.key")
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
