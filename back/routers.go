package main

import (
	"log"
	"net/http"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/csrf"
)

func AppLoginPage(c *gin.Context) {

	c.HTML(
		http.StatusOK,
		"login.html",
		gin.H{
			"error":          "",
			"title":          *config.Title,
			csrf.TemplateTag: csrf.TemplateField(c.Request),
		},
	)
}

func AppLogin(c *gin.Context) {

	username := c.PostForm("username")
	password := c.PostForm("password")

	if err := ValidateUser(c, &username, &password); err != nil {
		log.Printf("[Log] Login failed: %s (%s)", username, c.ClientIP())

		obj := gin.H{
			"error":          err.Error(),
			"title":          *config.Title,
			csrf.TemplateTag: csrf.TemplateField(c.Request),
		}

		c.HTML(http.StatusUnauthorized, "login.html", obj)
	}

	log.Printf("[Log] Login success: %s (%s)", username, c.ClientIP())

	c.Redirect(http.StatusMovedPermanently, "/login")
}

func AppLogout(c *gin.Context) {

	session := sessions.Default(c)
	user := session.Get(gin.AuthUserKey)
	session.Clear()

	if err := session.Save(); err != nil {
		log.Printf("[Log] Logout failed: %s", user)

		c.JSON(http.StatusBadRequest, map[string]bool{"ok": false})
		return
	}

	log.Printf("[Log] Logout success: %s", user)

	c.Redirect(http.StatusMovedPermanently, "/login")
}

func AppIndexPage(c *gin.Context) {
	c.Status(http.StatusOK)
}
