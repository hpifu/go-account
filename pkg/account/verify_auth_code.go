package account

type VerifyAuthCodeReq struct {
	Type  string `json:"type,omitempty" http:"param"`
	Phone string `json:"phone,omitempty" http:"param"`
	Email string `json:"email,omitempty" http:"param"`
	Code  string `json:"code,omitempty" http:"param"`
}

type VerifyAuthCodeRes struct {
	OK  bool   `json:"ok"`
	Tip string `json:"tip"`
}
