package libs

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"crypto/x509/pkix"
	"log"
)

func NewCSR(subject string) (*x509.CertificateRequest, *rsa.PrivateKey, error) {
	bits := 4096

	log.Printf("Generating %d-bit RSA key", bits)
	certKey, err := rsa.GenerateKey(rand.Reader, bits)
	if err != nil {
		return nil, nil, err
	}

	template := &x509.CertificateRequest{
		SignatureAlgorithm: x509.SHA256WithRSA,
		PublicKeyAlgorithm: x509.RSA,
		PublicKey:          &certKey.PublicKey,
		Subject:            pkix.Name{CommonName: subject},
	}

	log.Println("Generating CSR")
	csrDER, err := x509.CreateCertificateRequest(rand.Reader, template, certKey)
	if err != nil {
		return nil, nil, err
	}

	csr, err := x509.ParseCertificateRequest(csrDER)
	if err != nil {
		return nil, nil, err
	}
	return csr, certKey, nil
}
