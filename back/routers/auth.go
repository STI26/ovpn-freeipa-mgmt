package routers

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/sti26/ovpn_freeipa_mgmt/controllers"
)

func (r *Routers) AppLogin(c *gin.Context) {

	type creds struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}

	var data creds
	c.BindJSON(&data)

	token, err := controllers.ValidateUser(c, r.Ipa, &data.Username, &data.Password, r.Cfg.IPAAllowGroup)
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
