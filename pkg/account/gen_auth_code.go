package account

type GenAuthCodeReq struct {
	Type      string `json:"type" http:"body"`
	FirstName string `json:"firstName,omitempty" http:"body"`
	LastName  string `json:"lastName,omitempty" http:"body"`
	Email     string `json:"email,omitempty" http:"body"`
	Phone     string `json:"phone,omitempty" http:"body"`
}

type GenAuthCodeRes struct {
	OK bool `json:"ok"`
}
