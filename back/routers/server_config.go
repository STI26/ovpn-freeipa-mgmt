package routers

import (
	"bytes"
	"crypto/x509"
	"encoding/pem"
	"log"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/sti26/ovpn_freeipa_mgmt/libs"
	"github.com/sti26/ovpn_freeipa_mgmt/models"
	"github.com/sti26/ovpn_freeipa_mgmt/serializers"
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

func (r *Routers) AppCreateServerCert(c *gin.Context) {

	userID := c.Param("uid")

	var h models.Headers
	c.BindHeader(&h)

	// Create csr and key
	csr, key, err := libs.NewCSR(userID)
	if err != nil {
		log.Println("[createCSR] [Warn] ", err)
		c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
		return
	}

	// Create new certificate
	servers, _ := r.Ipa.GetIPAServers()
	resp, code, err := r.Ipa.Jrpc(c, h.Authorization, "cert_request",
		[]interface{}{
			csr.Raw,
		},
		map[string]interface{}{
			"cacn":         *r.Cfg.IPAcacn,
			"request_type": "pkcs10",
			"principal":    libs.GetPrincipal(userID, servers[0]),
			"add":          true,
		},
	)
	if err != nil {
		log.Println("[JSON_RPC] [Warn] ", err)
		c.JSON(code, map[string]string{"error": err.Error()})
		return
	}

	var obj models.RespObjCertRequest
	err = libs.ParseResponse(resp, &obj)
	if err != nil {
		log.Println("[parseResponse] [Warn] ", err)
		c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
		return
	}

	// Save Certificate
	cert := serializers.CertSerialazer(&obj)
	err = os.WriteFile(r.Ovpn.Config.Cert, []byte(cert), 0644)
	if err != nil {
		log.Println("[WriteFile] [Warn] ", err)
		c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
		return
	}

	// Convert key from der to pem
	var PrivateKeyRow bytes.Buffer

	b, _ := x509.MarshalPKCS8PrivateKey(key)
	var privateKey = &pem.Block{
		Type:  "PRIVATE KEY",
		Bytes: b,
	}

	err = pem.Encode(&PrivateKeyRow, privateKey)
	if err != nil {
		log.Println("[pem.Encode] [Warn] ", err)
		c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
		return
	}

	// Save key
	err = os.WriteFile(r.Ovpn.Config.Key, PrivateKeyRow.Bytes(), 0644)
	if err != nil {
		log.Println("[WriteFile] [Warn] ", err)
		c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, map[string]string{"error": ""})
}

func (r *Routers) AppCreateCA(c *gin.Context) {
	var h models.Headers
	c.BindHeader(&h)

	// Check Authentication
	if _, code, err := r.Ipa.Jrpc(c, h.Authorization, "ping", []any{}, &gin.H{}); err != nil {
		log.Println("[JSON_RPC] [Warn] ", err, "|", code)
		c.JSON(code, map[string]string{"error": err.Error()})
		return
	}

	if err := r.Ovpn.UpdateCA(); err != nil {
		log.Println("[UpdateCA] [Warn] ", err)
		c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, map[string]string{"error": ""})
}

func (r *Routers) AppCreateCrl(c *gin.Context) {
	var h models.Headers
	c.BindHeader(&h)

	// Check Authentication
	if _, code, err := r.Ipa.Jrpc(c, h.Authorization, "ping", []any{}, &gin.H{}); err != nil {
		log.Println("[JSON_RPC] [Warn] ", err, "|", code)
		c.JSON(code, map[string]string{"error": err.Error()})
		return
	}

	if err := r.Ovpn.UpdateCrl(); err != nil {
		log.Println("[UpdateCrl] [Warn] ", err)
		c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, map[string]string{"error": ""})
}

func (r *Routers) AppCreateDH(c *gin.Context) {
	var h models.Headers
	c.BindHeader(&h)

	// Check Authentication
	if _, code, err := r.Ipa.Jrpc(c, h.Authorization, "ping", []any{}, &gin.H{}); err != nil {
		log.Println("[JSON_RPC] [Warn] ", err, "|", code)
		c.JSON(code, map[string]string{"error": err.Error()})
		return
	}

	// TODO: ...

	c.JSON(http.StatusOK, map[string]string{"error": ""})
}

func (r *Routers) AppCreateTlsAuth(c *gin.Context) {
	var h models.Headers
	c.BindHeader(&h)

	// Check Authentication
	if _, code, err := r.Ipa.Jrpc(c, h.Authorization, "ping", []any{}, &gin.H{}); err != nil {
		log.Println("[JSON_RPC] [Warn] ", err, "|", code)
		c.JSON(code, map[string]string{"error": err.Error()})
		return
	}

	// TODO: ...

	c.JSON(http.StatusOK, map[string]string{"error": ""})
}
