package models

type RespObjHostFind struct {
	Count   int     `json:"count"`
	Message any     `json:"messages"`
	Result  []Hosts `json:"result"`
	Summary string  `json:"summary"`
}

type Hosts struct {
	DN                 string            `json:"dn"`
	FQDN               []string          `json:"fqdn"`
	KrbCanonicalName   []string          `json:"krbcanonicalname"`
	KrbPrincipalName   []string          `json:"krbprincipalname"`
	NsHardwarePlatform []string          `json:"nshardwareplatform"`
	NsOsVersion        []string          `json:"nsosversion"`
	UserCertificate    []UserCertificate `json:"usercertificate"`
}
