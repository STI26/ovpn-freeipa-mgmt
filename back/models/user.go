package models

type RespObjUserFind struct {
	Count   int     `json:"count"`
	Message any     `json:"messages"`
	Result  []Users `json:"result"`
	Summary string  `json:"summary"`
}

type Users struct {
	DN               string            `json:"dn"`
	GidNumber        []string          `json:"gidnumber"`
	GivenName        []string          `json:"givenname"`
	HomeDirectory    []string          `json:"homedirectory"`
	KrbCanonicalName []string          `json:"krbcanonicalname"`
	KrbPrincipalName []string          `json:"krbprincipalname"`
	LoginShell       []string          `json:"loginshell"`
	Mail             []string          `json:"mail"`
	NsAccountLock    bool              `json:"nsaccountlock"`
	SN               []string          `json:"sn"`
	UID              []string          `json:"uid"`
	UidNumber        []string          `json:"uidnumber"`
	MemberOfGroup    []string          `json:"memberof_group"`
	UserCertificate  []UserCertificate `json:"usercertificate"`
}
