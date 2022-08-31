package routers

import (
	"bytes"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/exec"
	"strings"

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

	type config struct {
		Server libs.Network `json:"server"`
		Local  string       `json:"local"`
	}
	var data config

	err := c.BindJSON(&data)
	if err != nil {
		c.JSON(http.StatusBadRequest, &gin.H{
			"error": err.Error(),
		})
		return
	}

	local := models.OvpnAvailableOptions["local"]
	local.Value = data.Local
	r.Ovpn.Optons = append(r.Ovpn.Optons, local)

	r.Ovpn.Server = data.Server

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
	params := map[string]interface{}{
		"cacn":         *r.Cfg.IPAcacn,
		"request_type": "pkcs10",
		"principal":    libs.GetPrincipal(userID, servers[0]),
		"add":          true,
	}
	if *r.Cfg.CaProfile != "" {
		params["profile_id"] = *r.Cfg.CaProfile
	}

	resp, code, err := r.Ipa.Jrpc(c, h.Authorization, "cert_request",
		[]interface{}{
			csr.Raw,
		},
		params,
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
	_cert, _ := r.Ovpn.GetOptionByKey("cert")
	libs.BackupFile(_cert.Value)

	err = os.WriteFile(_cert.Value, []byte(cert), 0644)
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
	_key, _ := r.Ovpn.GetOptionByKey("key")
	libs.BackupFile(_key.Value)

	err = os.WriteFile(_key.Value, PrivateKeyRow.Bytes(), 0644)
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

	resp, code, err := r.Ipa.Jrpc(c, h.Authorization, "ca_find", []any{}, &gin.H{
		"cn":  *r.Cfg.IPAcacn,
		"all": true,
	})
	if err != nil {
		log.Println("[JSON_RPC] [Warn] ", err)
		c.JSON(code, map[string]string{"error": err.Error()})
		return
	}

	var obj models.RespObjCaRequest
	err = libs.ParseResponse(resp, &obj)
	if err != nil {
		log.Println("[parseResponse] [Warn] ", err)
		c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
		return
	}

	ca := serializers.CaSerialazer(&obj)
	_ca, _ := r.Ovpn.GetOptionByKey("ca")
	libs.BackupFile(_ca.Value)

	err = os.WriteFile(_ca.Value, []byte(ca), 0644)
	if err != nil {
		log.Println("[WriteFile] [Warn] ", err)
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

	_dh, _ := r.Ovpn.GetOptionByKey("dh")
	libs.BackupFile(_dh.Value)

	cmd := exec.Command("openssl", "dhparam", "-out", _dh.Value, "2048")
	output, _ := cmd.Output()
	fmt.Println(string(output))

	log.Printf("[Log] dh updated: %s (%s)", c.Request.UserAgent(), c.ClientIP())

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

	_tlsAuth, _ := r.Ovpn.GetOptionByKey("tls-auth")
	pathToTlsAuth := strings.Split(_tlsAuth.Value, " ")[0]
	libs.BackupFile(pathToTlsAuth)

	cmd := exec.Command("openvpn", "--genkey", "secret", pathToTlsAuth)
	if _, err := cmd.Output(); err != nil {
		log.Println("[TlsAuth] [Warn] ", err)
	}

	log.Printf("[Log] tls-auth updated: %s (%s)", c.Request.UserAgent(), c.ClientIP())

	c.JSON(http.StatusOK, map[string]string{"error": ""})
}
