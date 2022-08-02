package serializers

import (
	"github.com/gin-gonic/gin"
	"github.com/sti26/ovpn_freeipa_mgmt/models"
)

func UserListSerialazer(respObj *models.RespObjUserFind, fullInfo bool) *gin.H {

	users := []map[string]any{}
	for _, i := range respObj.Result {
		users = append(users, map[string]any{
			"id":                   i.UidNumber[0],
			"name":                 i.UID[0],
			"numberOfCertificates": len(i.UserCertificate),
		})
	}

	return &gin.H{
		"error": "",
		"users": users,
	}
}
