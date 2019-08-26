package account

type VerifyReq struct {
	Field string `json:"field,omitempty" http:"param"`
	Value string `json:"value,omitempty" http:"param"`
}

type VerifyRes struct {
	OK  bool   `json:"ok"`
	Tip string `json:"tip"`
}
