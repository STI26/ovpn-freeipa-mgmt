package routers

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/sti26/ovpn_freeipa_mgmt/libs"
	"github.com/sti26/ovpn_freeipa_mgmt/models"
)

func (r *Routers) ValidateUser(c *gin.Context, user, pass *string) (string, error) {
	// Validate form input
	if strings.Trim(*user, " ") == "" || strings.Trim(*pass, " ") == "" {
		log.Println("[JSON_RPC] [Warn] Login/Password can't be empty")
		return "", errors.New("empty login/password is not allowed")
	}

	token, err := r.Ipa.PostAuth(user, pass)
	if err != nil {
		return "", err
	}

	resp, _, err := r.Ipa.Jrpc(c, token, "user_find", []interface{}{}, map[string]interface{}{
		"uid":      *user,
		"in_group": *r.Cfg.IPAAllowGroup,
	})
	if err != nil {
		log.Println("[JSON_RPC] [Warn] ", err)
		return "", err
	}

	var obj models.RespObjUserFind
	err = libs.ParseResponse(resp, &obj)
	if err != nil {
		log.Println("[JSON_RPC] [Error] ", err)
		return "", err
	}

	if len(obj.Result) != 1 {
		return "", fmt.Errorf("[JSON_RPC] [Warn] %s not found in %s group", *user, *r.Cfg.IPAAllowGroup)
	}

	return token, nil
}

func (r *Routers) AppLogin(c *gin.Context) {

	type creds struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}

	var data creds
	c.BindJSON(&data)

	token, err := r.ValidateUser(c, &data.Username, &data.Password)
	if err != nil {
		log.Printf("[Log] Login failed: %s (%s)", data.Username, c.ClientIP())

		obj := gin.H{
			"error": err.Error(),
			"token": "",
		}

		c.JSON(http.StatusForbidden, obj)
		return
	}

	log.Printf("[Log] Login success: %s (%s)", data.Username, c.ClientIP())

	c.JSON(http.StatusOK, map[string]string{"error": "", "token": token})
}

func (r *Routers) AppLogout(c *gin.Context) {
	type body struct {
		User  string `json:"user"`
		Token string `json:"token"`
	}

	var data body
	c.BindJSON(&data)

	_, code, err := r.Ipa.Jrpc(c, data.Token, "session_logout")
	if err != nil {
		log.Println("[JSON_RPC] [Warn] ", err)
		c.JSON(code, map[string]string{"error": err.Error()})
		return
	}

	log.Printf("[Log] Logout success: %s", data.User)

	c.JSON(http.StatusOK, map[string]string{"error": ""})
}
