package routers

import (
	"bytes"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"text/template"

	"github.com/gin-gonic/gin"
	"github.com/sti26/ovpn_freeipa_mgmt/libs"
	"github.com/sti26/ovpn_freeipa_mgmt/models"
	"github.com/sti26/ovpn_freeipa_mgmt/serializers"
)

func (r *Routers) AppGetConfig(c *gin.Context) {

	userID := c.Param("uid")

	var h models.Headers
	c.BindHeader(&h)

	// Check Authentication
	if _, code, err := r.Ipa.Jrpc(c, h.Authorization, "ping", []any{}, &gin.H{}); err != nil {
		log.Println("[JSON_RPC] [Warn] ", err, "|", code)
		c.JSON(code, map[string]string{"error": err.Error()})
		return
	}

	// Get routes
	_ccd, _ := r.Ovpn.GetOptionByKey("client-config-dir")
	path := filepath.Join(_ccd.Value, userID)
	routes, err := os.ReadFile(path)
	if err != nil {
		log.Println("[ReadFile] [Warn] ", err)
	}

	// Get IP
	_ipp, _ := r.Ovpn.GetOptionByKey("ifconfig-pool-persist")
	ipp := libs.IfconfigPoolPersist{
		Path:    _ipp.Value,
		Network: r.Ovpn.Server,
	}
	ip, err := ipp.GetIP(userID)
	if err != nil {
		log.Println("[ReadFile] [Warn] ", err)
	}

	c.JSON(http.StatusOK, gin.H{
		"error": "",
		"config": gin.H{
			"routes": string(routes),
			"ip":     ip,
		},
	})
}

func (r *Routers) AppDownloadConfig(c *gin.Context) {

	userID := c.Param("uid")

	certID, err := strconv.Atoi(c.Param("cid"))
	if err != nil {
		c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
		return
	}

	var h models.Headers
	c.BindHeader(&h)

	// Get user certificate
	resp, code, err := r.Ipa.Jrpc(c, h.Authorization, "cert_find", []interface{}{}, map[string]interface{}{
		"cacn":              *r.Cfg.IPAcacn,
		"subject":           userID,
		"min_serial_number": certID,
		"max_serial_number": certID,
		"all":               true,
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

	cert := serializers.CertsSerialazer(&obj)

	// Get user key
	path := filepath.Join(*r.Cfg.OvpnKeys, fmt.Sprintf("%d.key", certID))
	key, err := os.ReadFile(path)
	if err != nil {
		log.Println("[readKey] [Warn] ", err)
		c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
		return
	}

	// Get user tls-auth
	_tlsAuth, _ := r.Ovpn.GetOptionByKey("tls-auth")
	var tlsAuth []byte
	if _tlsAuth != nil {
		pathToTlsAuth := strings.Split(_tlsAuth.Value, " ")[0]
		tlsAuth, err = os.ReadFile(pathToTlsAuth)
		if err != nil {
			log.Println("[readTlsAuth] [Warn] ", err)
			c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
			return
		}
	}

	// Get user ca
	_ca, _ := r.Ovpn.GetOptionByKey("ca")
	ca, err := os.ReadFile(_ca.Value)
	if err != nil {
		log.Println("[readCA] [Warn] ", err)
		c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
		return
	}

	// Fill template
	config := new(bytes.Buffer)
	tmpl, err := template.ParseFiles("assets/client.tmpl")
	if err != nil {
		log.Println("[parseTemplate] [Warn] ", err)
		c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
		return
	}

	err = tmpl.Execute(config, map[string]interface{}{
		"serverConfig": r.Ovpn.Optons,
		"ca":           string(ca),
		"tlsAuth":      string(tlsAuth),
		"cert":         cert,
		"key":          string(key),
	})
	if err != nil {
		log.Println("[tmpl.Execute] [Warn] ", err)
		c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
		return
	}

	c.Header("Content-Type", "text/plain")
	c.Header("Content-Disposition", "attachment;filename="+userID+".ovpn")

	c.String(http.StatusOK, config.String())
}

func (r *Routers) AppCreateConfig(c *gin.Context) {

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
	os.MkdirAll(*r.Cfg.OvpnKeys, 0644)
	path := filepath.Join(*r.Cfg.OvpnKeys, fmt.Sprintf("%d.key", obj.Result.SerialNumber))

	err = os.WriteFile(path, PrivateKeyRow.Bytes(), 0644)
	if err != nil {
		log.Println("[WriteFile] [Warn] ", err)
		c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
		return
	}

	// Add IP
	_ipp, _ := r.Ovpn.GetOptionByKey("ifconfig-pool-persist")
	ipp := libs.IfconfigPoolPersist{
		Path:    _ipp.Value,
		Network: r.Ovpn.Server,
	}
	_, err = ipp.AddIP(userID)
	if err != nil {
		log.Println("[WriteFile] [Warn] ", err)
		c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
		return
	}

	// Add Default route
	_ccd, _ := r.Ovpn.GetOptionByKey("client-config-dir")
	os.MkdirAll(_ccd.Value, 0644)
	path = filepath.Join(_ccd.Value, userID)

	err = os.WriteFile(path, []byte("push \"redirect-gateway def1\"\n"), 0644)
	if err != nil {
		log.Println("[WriteFile] [Warn] ", err)
		c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, map[string]string{"error": ""})
}

func (r *Routers) AppUpdateConfig(c *gin.Context) {

	userID := c.Param("uid")

	var data struct {
		Routes string `json:"routes"`
		IP     string `json:"ip"`
	}
	c.BindJSON(&data)

	var h models.Headers
	c.BindHeader(&h)

	// Check Authentication
	if _, code, err := r.Ipa.Jrpc(c, h.Authorization, "ping", []any{}, &gin.H{}); err != nil {
		log.Println("[JSON_RPC] [Warn] ", err, "|", code)
		c.JSON(code, map[string]string{"error": err.Error()})
		return
	}

	// Update routes
	_ccd, _ := r.Ovpn.GetOptionByKey("client-config-dir")
	path := filepath.Join(_ccd.Value, userID)
	err := os.WriteFile(path, []byte(data.Routes), 0644)
	if err != nil {
		log.Println("[WriteFile] [Warn] ", err)
		c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
		return
	}

	// Update IP
	_ipp, _ := r.Ovpn.GetOptionByKey("ifconfig-pool-persist")
	ipp := libs.IfconfigPoolPersist{
		Path:    _ipp.Value,
		Network: r.Ovpn.Server,
	}
	ip, _ := ipp.GetIP(userID)
	if ip != data.IP {
		err := ipp.UpdateIP(userID, data.IP)
		if err != nil {
			log.Println("[ReadFile] [Warn] ", err)
			c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
			return
		}
	}

	c.JSON(http.StatusOK, map[string]string{"error": ""})
}

func (r *Routers) AppDeleteConfig(c *gin.Context) {

	userID := c.Param("uid")

	var h models.Headers
	c.BindHeader(&h)

	// Check Authentication
	if _, code, err := r.Ipa.Jrpc(c, h.Authorization, "ping", []any{}, &gin.H{}); err != nil {
		log.Println("[JSON_RPC] [Warn] ", err, "|", code)
		c.JSON(code, map[string]string{"error": err.Error()})
		return
	}

	_ipp, _ := r.Ovpn.GetOptionByKey("ifconfig-pool-persist")
	ipp := libs.IfconfigPoolPersist{
		Path:    _ipp.Value,
		Network: r.Ovpn.Server,
	}
	err := ipp.DeleteIP(userID)
	if err != nil {
		log.Println("[WriteFile] [Warn] ", err)
		c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
		return
	}

	_ccd, _ := r.Ovpn.GetOptionByKey("client-config-dir")
	path := filepath.Join(_ccd.Value, userID)
	err = os.RemoveAll(path)
	if err != nil {
		log.Println("[RemoveFile] [Warn] ", err)
		c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, map[string]string{"error": ""})
}
