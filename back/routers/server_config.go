package routers

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/sti26/ovpn_freeipa_mgmt/models"
)

func (r *Routers) AppGetServerConfig(c *gin.Context) {
	var h models.Headers
	c.BindHeader(&h)

	// Check Authentication
	if _, code, err := r.Ipa.Jrpc(c, h.Authorization, "ping", []any{}, &gin.H{}); err != nil {
		log.Println("[JSON_RPC] [Warn] ", err, "|", code)
		c.JSON(code, map[string]string{"error": err.Error()})
		return
	}

	// Reload config from file
	err := r.Ovpn.Init(*r.Cfg.OvpnConf)
	if err != nil {
		c.JSON(http.StatusOK, &gin.H{
			"error":  err.Error(),
			"config": r.Ovpn.GetServerConfig(),
		})
		return
	}

	c.JSON(http.StatusOK, &gin.H{
		"error":  "",
		"config": r.Ovpn.GetServerConfig(),
	})
}

func (r *Routers) AppCreateServerConfig(c *gin.Context) {

	var h models.Headers
	c.BindHeader(&h)

	// Check Authentication
	if _, code, err := r.Ipa.Jrpc(c, h.Authorization, "ping", []any{}, &gin.H{}); err != nil {
		log.Println("[JSON_RPC] [Warn] ", err, "|", code)
		c.JSON(code, map[string]string{"error": err.Error()})
		return
	}

	err := c.BindJSON(r.Ovpn)
	if err != nil {
		c.JSON(http.StatusBadRequest, &gin.H{
			"error": err.Error(),
		})
		return
	}

	err = r.Ovpn.CreateServerConfig(*r.Cfg.OvpnConf)
	if err != nil {
		c.JSON(422, &gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, &gin.H{
		"error": "",
	})
}
