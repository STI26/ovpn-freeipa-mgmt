package routers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/sti26/ovpn_freeipa_mgmt/libs"
	"github.com/sti26/ovpn_freeipa_mgmt/models"
	"github.com/sti26/ovpn_freeipa_mgmt/serializers"
)

func (r *Routers) AppGetUsers(c *gin.Context) {

	var h models.Headers
	c.BindHeader(&h)

	fullInfo, _ := c.GetQuery("all")
	fullInfoBool, _ := strconv.ParseBool(fullInfo)

	resp, code, err := r.Ipa.Jrpc(c, h.Authorization, "user_find", []interface{}{}, map[string]interface{}{
		"in_group": *r.Cfg.IPAFilterUserGroup,
		"all":      fullInfoBool,
	})
	if err != nil {
		c.JSON(code, map[string]string{"error": err.Error()})
		return
	}

	var obj models.RespObjUserFind
	err = libs.ParseResponse(resp, &obj)
	if err != nil {
		c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
		return
	}

	data := serializers.UserListSerialazer(&obj, fullInfoBool)

	c.JSON(http.StatusOK, data)
}
