package configs

import (
	"crypto/tls"
	"crypto/x509"
	"os"
)

type clientParamsTLS struct {
	brokerCA         string
	clientCA         string
	clientCert       string
	clientPrivateKey string
}

func NewTLSConfig(params clientParamsTLS) (*tls.Config, error) {
	rootCAs := x509.NewCertPool()
	clientCAs := x509.NewCertPool()

	brokerCert, err := os.ReadFile(params.brokerCA)
	if err != nil {
		return nil, &tlsConfigError{err}
	}
	clientCert, err := os.ReadFile(params.clientCA)
	if err != nil {
		return nil, &tlsConfigError{err}
	}

	rootCAs.AppendCertsFromPEM(brokerCert)
	clientCAs.AppendCertsFromPEM(clientCert)

	cert, err := tls.LoadX509KeyPair(params.clientCert, params.clientPrivateKey)
	if err != nil {
		return nil, &tlsConfigError{err}
	}

	cert.Leaf, err = x509.ParseCertificate(cert.Certificate[0])
	if err != nil {
		return nil, &tlsConfigError{err}
	}

	return &tls.Config{
		RootCAs:      rootCAs,
		ClientAuth:   tls.RequireAndVerifyClientCert,
		ClientCAs:    clientCAs,
		Certificates: []tls.Certificate{cert},
	}, nil

}
