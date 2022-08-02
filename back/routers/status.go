package routers

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/sti26/ovpn_freeipa_mgmt/models"
	"github.com/sti26/ovpn_freeipa_mgmt/serializers"
)

func (r *Routers) AppGetStatus(c *gin.Context) {
	var h models.Headers
	c.BindHeader(&h)

	// Check Authentication
	if _, code, err := r.Ipa.Jrpc(c, h.Authorization, "ping", []any{}, map[string]any{}); err != nil {
		log.Println("[JSON_RPC] [Warn] ", err, "|", code)
		c.JSON(code, map[string]string{"error": err.Error()})
		return
	}

	status, err := r.Ovpn.GetStatusInfo()
	if err != nil {
		c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
		return
	}

	data := serializers.StatusSerialazer(status)

	c.JSON(http.StatusOK, data)
}
