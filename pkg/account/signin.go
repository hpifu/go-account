package account

type SignInReq struct {
	Username string `json:"username,omitempty"`
	Password string `json:"password,omitempty"`
}

type SignInRes struct {
	Valid bool   `json:"valid"`
	Token string `json:"token"`
}
