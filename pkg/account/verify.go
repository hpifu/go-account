package account

type VerifyReq struct {
	Field string `json:"field,omitempty"`
	Value string `json:"value,omitempty"`
}

type VerifyRes struct {
	OK  bool   `json:"ok"`
	Tip string `json:"tip"`
}
