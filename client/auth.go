package client

import (
	"crypto/tls"
	"crypto/x509"
	"encoding/base64"
	"net/http"
	"os"
)

// AuthProvider defines the interface for applying authentication to HTTP requests.
// Implementations include BasicAuth, TokenAuth, and others.
type AuthProvider interface {
	Apply(req *http.Request) error
}

// BasicAuth provides HTTP Basic Authentication using a username and password.
type BasicAuth struct {
	Username string // The HTTP basic auth username
	Password string // The HTTP basic auth password
}

// Apply adds an Authorization header to the HTTP request using base64-encoded credentials.
func (a *BasicAuth) Apply(req *http.Request) error {
	credentials := a.Username + ":" + a.Password
	encoded := base64.StdEncoding.EncodeToString([]byte(credentials))
	req.Header.Set("Authorization", "Basic "+encoded)
	return nil
}

// TLSAuth configures an HTTP client with mutual TLS authentication.
// CertFile and KeyFile are the client cert and key.
// CAFile is the root certificate authority used for server validation.
type TLSAuth struct {
	CertFile           string // Path to the client certificate PEM file
	KeyFile            string // Path to the client private key PEM file
	CAFile             string // Path to the CA certificate PEM file
	InsecureSkipVerify bool   // Whether to skip TLS certificate validation (not recommended)
}

// ConfigureClient sets the client's HTTP transport to use the specified TLS configuration.
func (a *TLSAuth) ConfigureClient(client *http.Client) error {
	cert, err := tls.LoadX509KeyPair(a.CertFile, a.KeyFile)
	if err != nil {
		return err
	}

	caCert, err := os.ReadFile(a.CAFile)
	if err != nil {
		return err
	}

	caPool := x509.NewCertPool()
	caPool.AppendCertsFromPEM(caCert)

	tlsConfig := &tls.Config{
		Certificates:       []tls.Certificate{cert},
		RootCAs:            caPool,
		InsecureSkipVerify: a.InsecureSkipVerify,
	}

	client.Transport = &http.Transport{
		TLSClientConfig: tlsConfig,
	}

	return nil
}
