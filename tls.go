package config

import (
	"crypto/tls"
	"crypto/x509"
	"errors"
	"io/ioutil"
)

const (
	_DEFAULT_SERVERNAME = "archer"

	_CFG_CERT    = "cert"
	_CFG_KEY     = "key"
	_CFG_ROOTCAS = "rootCAs"
)

type TLSConfigOptions struct {
	CertFilename    string
	KeyFilename     string
	RootCAsFilename string
}

func createDefaultTLSConfig() *tls.Config {

	c := &tls.Config{
		// Force use of the most recent version
		MinVersion: tls.VersionTLS12,

		Certificates: make([]tls.Certificate, 1),

		// We are always talking to / being a server named simply "archer". Effectively
		// we aren't using host name verification functionality from TLS
		ServerName: _DEFAULT_SERVERNAME,

		RootCAs:    x509.NewCertPool(),
		ClientAuth: tls.RequireAndVerifyClientCert,

		// Removing the RC4 and 3DES cipher suites from the defaults
		// ^^^ good move
		CipherSuites: []uint16{
			tls.TLS_RSA_WITH_AES_128_CBC_SHA,
			tls.TLS_RSA_WITH_AES_256_CBC_SHA,
			tls.TLS_ECDHE_ECDSA_WITH_AES_128_CBC_SHA,
			tls.TLS_ECDHE_ECDSA_WITH_AES_256_CBC_SHA,
			tls.TLS_ECDHE_RSA_WITH_AES_128_CBC_SHA,
			tls.TLS_ECDHE_RSA_WITH_AES_256_CBC_SHA,
			tls.TLS_ECDHE_RSA_WITH_AES_128_GCM_SHA256,
			tls.TLS_ECDHE_ECDSA_WITH_AES_128_GCM_SHA256,
		},
	}
	// Use the same CAs to verify the clients as we do for verifying servers since we
	// are a closed system
	c.ClientCAs = c.RootCAs

	return c
}

func CreateTLSConfig(o *TLSConfigOptions) (c *tls.Config, err error) {

	if o == nil {
		return nil, errors.New("No TLSConfigOptions value was provided")
	}

	if len(o.CertFilename) == 0 {
		return nil, errors.New("No cert value was found with a filename for a certificate .pem file")
	}

	if len(o.KeyFilename) == 0 {
		return nil, errors.New("No key value was found with a filename for a key .pem file")
	}

	if len(o.RootCAsFilename) == 0 {
		return nil, errors.New("No rootCAs value was found with a filename for a root CA chain .pem file")
	}

	c = createDefaultTLSConfig()

	c.Certificates[0], err = tls.LoadX509KeyPair(o.CertFilename, o.KeyFilename)
	if err != nil {
		return nil, err
	}

	data, err := ioutil.ReadFile(o.RootCAsFilename)
	if err != nil {
		return nil, err
	}

	if !c.RootCAs.AppendCertsFromPEM(data) {
		return nil, errors.New("Unable to parse root CAs from " + o.RootCAsFilename)
	}

	return c, nil
}

func (self *AclNode) TLSConfigOptions() (*TLSConfigOptions, error) {

	o := &TLSConfigOptions{
		CertFilename:    self.ChildAsString(_CFG_CERT),
		KeyFilename:     self.ChildAsString(_CFG_KEY),
		RootCAsFilename: self.ChildAsString(_CFG_ROOTCAS),
	}

	return o, nil
}

func (self *AclNode) TLSConfig() (c *tls.Config, err error) {

	tlsConfigOptions, err := self.TLSConfigOptions()
	if err != nil {
		return nil, err
	}

	return CreateTLSConfig(tlsConfigOptions)
}
