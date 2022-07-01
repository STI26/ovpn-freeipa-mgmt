package main

import (
	"embed"
	"html/template"
	"log"
	"net/http"
	"strings"

	"github.com/gin-contrib/cors"
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/csrf"
	adapter "github.com/gwatts/gin-adapter"
)

// content holds our static web server content.
//go:embed public/templates/* public/assets/login.min.css
var content embed.FS

var config = GlobalConfig{}
var ipaSession = IPASession{}

func setupRouter() *gin.Engine {

	store := cookie.NewStore([]byte(*config.CookieSecret))
	store.Options(sessions.Options{MaxAge: *config.SessionsMaxAge * 60 * 60})

	r := gin.Default()
	r.Use(cors.Default())
	r.Use(sessions.Sessions("ldapsession", store))

	addr := strings.Split(*config.ListenAddress, ":")
	r.SetTrustedProxies([]string{addr[0]})

	templ := template.Must(template.New("").ParseFS(content, "public/templates/*.html"))
	r.SetHTMLTemplate(templ)
	r.StaticFS("/static", http.FS(content))

	if *config.CSRF {
		csrfMiddleware := csrf.Protect([]byte(*config.CSRFSecret), csrf.Path("/"))
		r.Use(adapter.Wrap(csrfMiddleware))
	}

	gRoot := r.Group("/")
	{
		// Routes without auth
		gRoot.GET("/login", AppLoginPage)
		gRoot.POST("/login", AppLogin)
		gRoot.GET("/logout", AppLogout)

		// Routes with auth
		gRoot.GET("/", RPCAuthMiddleware(true), AppIndexPage)
		gRoot.GET("/auth", RPCAuthMiddleware(false), AppIndexPage)
	}

	return r
}

func main() {
	gin.SetMode(gin.ReleaseMode)

	config.init()

	r := setupRouter()

	log.Printf(
		"~~~~~~ OpenVPN Managment ~~~~~~\n"+
			"\tListenAddress:       %s\n"+
			"\tLdapAllowGroup:      %s\n",
		*config.ListenAddress, *config.LdapAllowGroup,
	)

	r.Run(*config.ListenAddress)
}
