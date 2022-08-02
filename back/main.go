package main

import (
	"log"
	"net"
	"strings"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/sti26/ovpn_freeipa_mgmt/libs"
	"github.com/sti26/ovpn_freeipa_mgmt/routers"
)

var config = libs.GlobalConfig{}

var ipaClient libs.FreeIPA
var ovpn libs.OpenVPN

func setupRouter(rts *routers.Routers) *gin.Engine {

	r := gin.Default()
	cfg := cors.DefaultConfig()
	cfg.AllowAllOrigins = true
	cfg.AddAllowHeaders("*")
	// config.AllowHeaders = []string{"Content-Type", "Authorization"}
	r.Use(cors.New(cfg))

	gRoot := r.Group("/")
	{
		gRoot.POST("/login", rts.AppLogin)
		gRoot.POST("/logout", rts.AppLogout)

		gRoot.GET("/users/:uid/config/:cid", rts.AppDownloadConfig)
		gRoot.GET("/users/:uid/config", rts.AppGetConfig)
		gRoot.POST("/users/:uid/config", rts.AppCreateConfig)
		gRoot.PUT("/users/:uid/config/:cid", rts.AppUpdateConfig)
		gRoot.DELETE("/users/:uid/config/:cid", rts.AppDeleteConfig)

		gRoot.GET("/users", rts.AppGetUsers)
		gRoot.GET("/hosts", rts.AppGetHosts)
		gRoot.GET("/certs", rts.AppGetCerts)
		gRoot.DELETE("/certs/:cid", rts.AppRevokeCert)
		gRoot.GET("/status", rts.AppGetStatus)
		gRoot.GET("/config", rts.AppGetServerConfig)
	}

	return r
}

func main() {
	gin.SetMode(gin.ReleaseMode)

	config.Init()
	ipaClient = libs.FreeIPA{Domain: *config.IPADomain, Server: *config.IPAServer, Secret: config.Secret}
	ipaServers, _ := ipaClient.GetIPAServers()

	ovpn = libs.OpenVPN{IPAHosts: ipaServers}
	ovpn.Init(*config.OvpnConf)

	mask := net.ParseIP(ovpn.Server.Mask)
	sz, _ := net.IPMask(mask.To4()).Size()

	rts := routers.Routers{Ipa: &ipaClient, Cfg: &config, Ovpn: &ovpn}

	r := setupRouter(&rts)

	log.Printf(
		"~~~~~~ OpenVPN Managment ~~~~~~\n"+
			"\tListen address:        %s\n"+
			"\tIPA Servers:           %s\n"+
			"\tIPA Allow Group:       %s\n"+
			"\tOpenVPN Network:       %s/%d\n",
		*config.ListenAddress,
		strings.Join(ipaServers, ", "),
		*config.IPAAllowGroup,
		ovpn.Server.IP, sz,
	)

	r.Run(*config.ListenAddress)
}
