package libs

import (
	"crypto/tls"
	"flag"
	"net/http"

	"github.com/thanhpk/randstr"
)

type GlobalConfig struct {
	ListenAddress *string

	IPADomain          *string
	IPAServer          *string
	IPAAllowGroup      *string
	IPAFilterUserGroup *string
	IPAFilterHostGroup *string
	IPAcacn            *string

	OvpnConf *string
	OvpnKeys *string

	Secret string

	ShowVersion *bool
	Version     *string
}

func (cfg *GlobalConfig) Init() {
	http.DefaultTransport.(*http.Transport).TLSClientConfig = &tls.Config{InsecureSkipVerify: true}

	cfg.ListenAddress = flag.String("addr", "127.0.0.1:8000", "Listening and serving address.")

	cfg.IPADomain = flag.String("ipa-domain", "", "Domain with IPA servers. (search by SRV record)")
	cfg.IPAServer = flag.String("ipa-server", "https://ipa.local.net", "FreeIPA server with a scheme.")
	cfg.IPAAllowGroup = flag.String("ipa-allowgroup", "admins", "IPA group with allowed access.")
	cfg.IPAFilterUserGroup = flag.String("ipa-usergroup", "", "IPA user group.")
	cfg.IPAFilterHostGroup = flag.String("ipa-hostgroup", "", "IPA host group.")
	cfg.IPAcacn = flag.String("ipa-cacn", "openVPN", "Name of issuing CA.")

	cfg.OvpnConf = flag.String("ovpn-serverconf", "/etc/openvpn/server/server.conf", "Path to openvpn server.conf file.")
	cfg.OvpnKeys = flag.String("ovpn-keys", "/etc/openvpn/keys", "Path to folder with user keys.")

	cfg.ShowVersion = flag.Bool("version", false, "Show version.")

	flag.Parse()

	cfg.Secret = randstr.String(16)
}
