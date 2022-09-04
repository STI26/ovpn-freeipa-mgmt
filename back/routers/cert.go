package routers

import (
	"io/fs"
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/sti26/ovpn_freeipa_mgmt/libs"
	"github.com/sti26/ovpn_freeipa_mgmt/models"
	"github.com/sti26/ovpn_freeipa_mgmt/serializers"
)

func (r *Routers) AppGetCerts(c *gin.Context) {

	var h models.Headers
	c.BindHeader(&h)

	subject, _ := c.GetQuery("subject")

	resp, code, err := r.Ipa.Jrpc(c, h.Authorization, "cert_find", []interface{}{}, map[string]interface{}{
		"cacn":    *r.Cfg.IPAcacn,
		"subject": subject,
		"exactly": true,
	})
	if err != nil {
		log.Println("[JSON_RPC] [Warn] ", err)
		c.JSON(code, map[string]string{"error": err.Error()})
		return
	}

	var obj models.RespObjCertsRequest
	err = libs.ParseResponse(resp, &obj)
	if err != nil {
		log.Println("[parseResponse] [Warn] ", err)
		c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
		return
	}

	data := serializers.CertListSerialazer(&obj, r.getKeys())

	c.JSON(http.StatusOK, data)
}

func (r *Routers) AppRevokeCert(c *gin.Context) {

	certID, err := strconv.Atoi(c.Param("cid"))
	if err != nil {
		c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
		return
	}

	var h models.Headers
	c.BindHeader(&h)

	_, code, err := r.Ipa.Jrpc(c, h.Authorization, "cert_revoke",
		[]interface{}{
			certID,
		},
		map[string]interface{}{
			"cacn":              *r.Cfg.IPAcacn,
			"revocation_reason": 0,
		},
	)
	if err != nil {
		log.Println("[JSON_RPC] [Warn] ", err)
		c.JSON(code, map[string]string{"error": err.Error()})
		return
	}

	c.JSON(code, map[string]string{"error": ""})
}

func (r *Routers) AppGetCrl(c *gin.Context) {

	var h models.Headers
	c.BindHeader(&h)

	// Check Authentication
	_, code, err := r.Ipa.Jrpc(c, h.Authorization, "ping", []any{}, &gin.H{})
	if err != nil {
		log.Println("[JSON_RPC] [Warn] ", err, "|", code)
		c.JSON(code, map[string]string{"error": err.Error()})
		return
	}

	crlInfo, err := r.Ovpn.GetCrlInfo()
	if err != nil {
		log.Println("[GetCrlInfo] [Warn] ", err)
		c.JSON(code, map[string]string{"error": err.Error()})
		return
	}

	c.JSON(code, *serializers.CertRevokedListSerialazer(crlInfo))
}

func (r *Routers) getKeys() *[]fs.DirEntry {

	keys, err := os.ReadDir(*r.Cfg.OvpnKeys)
	if err != nil {
		log.Println("[readKeys] [Warn] ", err)
		return &[]fs.DirEntry{}
	}

	return &keys
}
