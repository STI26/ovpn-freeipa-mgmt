package models

type Headers struct {
	Authorization string `json:"Authorization"`
}

type UserCertificate struct {
	Base64 string `json:"__base64__"`
}

type RespObjPing struct {
	Message any    `json:"messages"`
	Summary string `json:"summary"`
}
