package routers

import "github.com/sti26/ovpn_freeipa_mgmt/libs"

type Routers struct {
	Ipa  *libs.FreeIPA
	Cfg  *libs.GlobalConfig
	Ovpn *libs.OpenVPN
}
