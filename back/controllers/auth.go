package controllers

import (
	"errors"
	"fmt"
	"log"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/sti26/ovpn_freeipa_mgmt/libs"
	"github.com/sti26/ovpn_freeipa_mgmt/models"
)

func ValidateUser(c *gin.Context, ipaClient *libs.FreeIPA, user, pass, allowGroup *string) (string, error) {
	// Validate form input
	if strings.Trim(*user, " ") == "" || strings.Trim(*pass, " ") == "" {
		log.Println("[JSON_RPC] [Warn] Login/Password can't be empty")
		return "", errors.New("empty login/password is not allowed")
	}

	token, err := ipaClient.PostAuth(user, pass)
	if err != nil {
		return "", err
	}

	resp, _, err := ipaClient.Jrpc(c, token, "user_find", []interface{}{}, map[string]interface{}{
		"uid":      *user,
		"in_group": *allowGroup,
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
		return "", fmt.Errorf("[JSON_RPC] [Warn] %s not found in %s group", *user, *allowGroup)
	}

	return token, nil
}
