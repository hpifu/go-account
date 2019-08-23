package account

type VerifyAuthCodeReq struct {
	Type  string `json:"type,omitempty"`
	Phone string `json:"phone,omitempty"`
	Email string `json:"email,omitempty"`
	Code  string `json:"code,omitempty"`
}

type VerifyAuthCodeRes struct {
	OK  bool   `json:"ok"`
	Tip string `json:"tip"`
}
