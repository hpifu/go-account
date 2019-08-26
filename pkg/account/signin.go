package account

type SignInReq struct {
	Username string `json:"username,omitempty" http:"body"`
	Password string `json:"password,omitempty" http:"body"`
}

type SignInRes struct {
	Valid bool   `json:"valid"`
	Token string `json:"token"`
}
