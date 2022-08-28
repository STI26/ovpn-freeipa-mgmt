package models

type RespObjCaRequest struct {
	Count   int    `json:"count"`
	Message any    `json:"messages"`
	Result  []Ca   `json:"result"`
	Summary string `json:"summary"`
}

type RespObjCertRequest struct {
	Count   int    `json:"count"`
	Message any    `json:"messages"`
	Result  Cert   `json:"result"`
	Summary string `json:"summary"`
}

type RespObjCertsRequest struct {
	Count   int    `json:"count"`
	Message any    `json:"messages"`
	Result  []Cert `json:"result"`
	Summary string `json:"summary"`
}

type Cert struct {
	Cacn             string `json:"cacn"`
	Issuer           string `json:"issuer"`
	SerialNumber     int    `json:"serial_number"`
	SerialNumberHex  string `json:"serial_number_hex"`
	Status           string `json:"status"`
	Subject          string `json:"subject"`
	ValidNotAfter    string `json:"valid_not_after"`
	ValidNotBefore   string `json:"valid_not_before"`
	Certificate      string `json:"certificate,omitempty"`
	RevocationReason int    `json:"revocation_reason,omitempty"`
	CertificateChain []struct {
		Base64 string `json:"__base64__"`
	} `json:"certificate_chain,omitempty"`
}

type Ca struct {
	Cn               []string `json:"cn"`
	IpaCaIssuerDn    []string `json:"ipacaissuerdn"`
	IpaCaSubjectDn   []string `json:"ipacasubjectdn"`
	Dn               string   `json:"dn"`
	Certificate      string   `json:"certificate,omitempty"`
	CertificateChain []struct {
		Base64 string `json:"__base64__"`
	} `json:"certificate_chain,omitempty"`
}
