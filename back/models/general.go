package models

type Headers struct {
	Authorization string `json:"Authorization"`
}

type UserCertificate struct {
	Base64 string `json:"__base64__"`
}
