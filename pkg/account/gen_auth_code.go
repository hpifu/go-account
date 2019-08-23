package account

type GenAuthCodeReq struct {
	Type      string `json:"type"`
	FirstName string `json:"firstName,omitempty"`
	LastName  string `json:"lastName,omitempty"`
	Email     string `json:"email,omitempty"`
	Phone     string `json:"phone,omitempty"`
}

type GenAuthCodeRes struct {
	OK bool `json:"ok"`
}
