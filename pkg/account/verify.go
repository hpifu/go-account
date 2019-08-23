package account

type VerifyReqBody struct {
	Field string `json:"field,omitempty"`
	Value string `json:"value,omitempty"`
}

type VerifyResBody struct {
	OK  bool   `json:"ok"`
	Tip string `json:"tip"`
}
