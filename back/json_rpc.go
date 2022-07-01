package main

import (
	"errors"
	"log"
	"net/http"
	"strings"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

func RPCAuthMiddleware(redirect bool) gin.HandlerFunc {
	return func(c *gin.Context) {
		session := sessions.Default(c)
		user := session.Get(gin.AuthUserKey)

		if user == nil {
			if redirect {
				c.Redirect(http.StatusPermanentRedirect, "/login")
			} else {
				c.AbortWithStatus(http.StatusUnauthorized)
			}
			return
		}

		// Continue down the chain to handler etc
		c.Next()
	}
}

func ValidateUser(c *gin.Context, user, pass *string) error {
	session := sessions.Default(c)

	// Validate form input
	if strings.Trim(*user, " ") == "" || strings.Trim(*pass, " ") == "" {
		log.Println("[RPC] [Warn] Login/Password can't be empty")
		return errors.New("empty login/password is not allowed")
	}

	err := ipaSession.post_auth(user, pass)
	if err != nil {
		log.Println("[RPC] [Warn] ", err.Error())
		return err
	}

	ipaSession.rpc(c, "user_find", []interface{}{}, map[string]interface{}{
		"uid": "adm",
		"all": true,
	})

	session.Set(gin.AuthUserKey, *user)
	session.Save()

	return nil
}

// func auth(c *gin.Context, conn *net.Conn, user, pass *string) error {

// 	// Check credentials
// 	if err := conn.Bind(config.userDN(c, *user), *pass); err != nil {
// 		log.Printf("[RPC] [Error] %s", err.Error())

// 		return errors.New("invalid credentials")
// 	}

// 	// Check filter
// 	request := ldap.NewSearchRequest(
// 		config.baseDN(c),
// 		ldap.ScopeWholeSubtree,
// 		ldap.DerefAlways,
// 		0,
// 		0,
// 		false,
// 		config.filter(c, *user),
// 		[]string{"dn"},
// 		nil,
// 	)
// 	response, err := conn.Search(request)
// 	if err != nil {

// 		log.Printf("[RPC] [Error] %s", err.Error())

// 		return errors.New("invalid filter")

// 	} else if len(response.Entries) != 1 {

// 		return errors.New("access denied")
// 	}

// 	log.Printf("[RPC] [Log] Authenticated successfuly!")

// 	return nil
// }
