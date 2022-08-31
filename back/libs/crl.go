package libs

import (
	"bytes"
	"crypto/x509"
	"crypto/x509/pkix"
	b64 "encoding/base64"
	"io"
	"net/http"
	"net/url"
	"os"
	"strings"
)

type Crl struct {
	path string
}

func (crl *Crl) Download(u *url.URL) error {

	resp, err := http.Get(u.String())
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	buf := new(bytes.Buffer)

	_, err = io.Copy(buf, resp.Body)
	if err != nil {
		return err
	}

	pem := crl.derToPem(buf)

	err = os.WriteFile(crl.path, []byte(pem), 0644)
	if err != nil {
		return err
	}

	return nil
}

func (crl *Crl) GetInfo() (*pkix.TBSCertificateList, error) {
	b, err := os.ReadFile(crl.path)
	if err != nil {
		return nil, err
	}

	b, err = crl.pemToDer(b)
	if err != nil {
		return nil, err
	}

	data, err := x509.ParseDERCRL(b)
	if err != nil {
		return nil, err
	}

	return &data.TBSCertList, nil
}

func (crl *Crl) derToPem(der *bytes.Buffer) string {

	return "-----BEGIN X509 CRL-----\n" +
		b64.StdEncoding.EncodeToString(der.Bytes()) +
		"\n-----END X509 CRL-----\n"
}

func (crl *Crl) pemToDer(pem []byte) ([]byte, error) {

	data := strings.Split(string(pem), "-----")

	der, err := b64.StdEncoding.DecodeString(data[2])
	if err != nil {
		return nil, err
	}

	return der, nil
}
