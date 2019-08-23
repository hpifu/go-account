package account

type SignOutReq struct {
	Token string `json:"token,omitempty"`
}

type SignOutRes struct {
	OK bool `json:"ok"`
}
