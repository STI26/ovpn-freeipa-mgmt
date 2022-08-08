package serializers

import (
	"fmt"
	"io/fs"

	"github.com/gin-gonic/gin"
	"github.com/sti26/ovpn_freeipa_mgmt/libs"
	"github.com/sti26/ovpn_freeipa_mgmt/models"
)

func CertListSerialazer(respObj *models.RespObjCertsRequest, keys *[]fs.DirEntry) *gin.H {

	certs := []map[string]interface{}{}
	for _, i := range respObj.Result {
		keyName := fmt.Sprintf("%d.key", i.SerialNumber)

		certs = append(certs, map[string]interface{}{
			"id":               i.SerialNumber,
			"status":           i.Status,
			"subject":          i.Subject,
			"valid_not_after":  i.ValidNotAfter,
			"valid_not_before": i.ValidNotBefore,
			"key_exists":       libs.KeyContains(keys, keyName),
		})
	}

	return &gin.H{
		"error":        "",
		"certificates": certs,
	}
}

func certSerialazer(cert string) string {
	res := "-----BEGIN CERTIFICATE-----"

	for i := 0; i < len(cert); i++ {
		if i%64 == 0 {
			res += "\n"
		}

		res += string(cert[i])
	}

	res += "\n-----END CERTIFICATE-----"

	return res
}

func CertsSerialazer(respObj *models.RespObjCertsRequest) string {

	if len(respObj.Result) != 1 {
		return ""
	}

	res := certSerialazer(respObj.Result[0].Certificate)

	return res
}

func CertSerialazer(respObj *models.RespObjCertRequest) string {

	res := certSerialazer(respObj.Result.Certificate)

	return res
}
