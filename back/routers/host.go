package routers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/sti26/ovpn_freeipa_mgmt/libs"
	"github.com/sti26/ovpn_freeipa_mgmt/models"
	"github.com/sti26/ovpn_freeipa_mgmt/serializers"
)

func (r *Routers) AppGetHosts(c *gin.Context) {

	var h models.Headers
	c.BindHeader(&h)

	fullInfo, _ := c.GetQuery("all")
	fullInfoBool, _ := strconv.ParseBool(fullInfo)

	resp, code, err := r.Ipa.Jrpc(c, h.Authorization, "host_find", []interface{}{}, map[string]interface{}{
		"in_hostgroup": *r.Cfg.IPAFilterHostGroup,
		"all":          fullInfoBool,
	})
	if err != nil {
		c.JSON(code, map[string]string{"error": err.Error()})
		return
	}

	var obj models.RespObjHostFind
	err = libs.ParseResponse(resp, &obj)
	if err != nil {
		c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
		return
	}

	data := serializers.HostListSerialazer(&obj, fullInfoBool)

	c.JSON(http.StatusOK, data)
}
