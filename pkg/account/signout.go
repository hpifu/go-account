package account

type SignOutReq struct {
	Token string `json:"token,omitempty" http:"param"`
}

type SignOutRes struct {
	OK bool `json:"ok"`
}
