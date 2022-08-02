package serializers

import (
	"github.com/gin-gonic/gin"
	"github.com/sti26/ovpn_freeipa_mgmt/models"
)

func HostListSerialazer(respObj *models.RespObjHostFind, fullInfo bool) *gin.H {

	if fullInfo {
		return &gin.H{
			"error": "",
			"hosts": respObj.Result,
		}
	}

	hosts := []map[string]any{}
	for _, i := range respObj.Result {
		hosts = append(hosts, map[string]any{
			"id":                   i.DN,
			"name":                 i.FQDN[0],
			"numberOfCertificates": len(i.UserCertificate),
		})
	}

	return &gin.H{
		"error": "",
		"hosts": hosts,
	}
}
