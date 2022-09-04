package serializers

import (
	"crypto/x509/pkix"
	"math/big"

	"github.com/gin-gonic/gin"
)

func CertRevokedListSerialazer(respObj *pkix.TBSCertificateList) *gin.H {
	type RevokedCertificates struct {
		SerialNumber   big.Int `json:"serial_number"`
		RevocationTime string  `json:"revocation_time"`
	}

	type obj struct {
		Issuer              string                `json:"issuer"`
		ThisUpdate          string                `json:"this_time"`
		NextUpdate          string                `json:"next_time"`
		RevokedCertificates []RevokedCertificates `json:"revoked_certificates"`
	}

	crl := obj{
		Issuer:     respObj.Issuer.String(),
		ThisUpdate: respObj.ThisUpdate.String(),
		NextUpdate: respObj.NextUpdate.String(),
	}

	for _, i := range respObj.RevokedCertificates {
		crl.RevokedCertificates = append(
			crl.RevokedCertificates,
			RevokedCertificates{
				SerialNumber:   *i.SerialNumber,
				RevocationTime: i.RevocationTime.String(),
			},
		)
	}

	return &gin.H{
		"error": "",
		"crl":   crl,
	}
}
