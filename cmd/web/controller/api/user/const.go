package user

import "time"

type LoginParmas struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type UserLogonInfo struct {
	Name     string    `json:"name,omitempty"`
	IP       string    `json:"ip"`
	ReadTime time.Time `json:"read_time"`
	ReadNum  int       `json:"read_num"`
	Token    string    `json:"token,omitempty"`
}
