package routers

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/sti26/ovpn_freeipa_mgmt/libs"
	"github.com/sti26/ovpn_freeipa_mgmt/models"
)

func (r *Routers) AppGetVersion(c *gin.Context) {
	var h models.Headers
	c.BindHeader(&h)

	// Check Authentication
	resp, code, err := r.Ipa.Jrpc(c, h.Authorization, "ping", []any{}, map[string]any{})
	if err != nil {
		log.Println("[JSON_RPC] [Warn] ", err, "|", code)
		c.JSON(code, &gin.H{"error": err.Error()})
		return
	}

	var obj models.RespObjPing
	err = libs.ParseResponse(resp, &obj)
	if err != nil {
		log.Println("[parseResponse] [Warn] ", err)
		c.JSON(http.StatusBadRequest, &gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, &gin.H{
		"error":       "",
		"version":     r.Cfg.Version,
		"ipa_version": obj.Summary,
	})
}
