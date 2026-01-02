package pshell

type ADUser struct {
	Name          string `json:"name"`
	Enabled       bool   `json:"enabled"`
	LogonCount    int    `json:"logoncount"`
	CanonicalName string `json:"CanonicalName"`
	LastLogon     int64  `json:"lastlogon"`
	BadPwdTime    int64  `json:"badpasswordtime"`
	BadPwdCount   int    `json:"badPwdCount"`
	PwdLastSet    int64  `json:"pwdLastSet"`
	Info          struct {
		OperatingSystem        string `json:"OperatingSystem"`
		OperatingSystemVersion string `json:"OperatingSystemVersion"`
		IP                     string `json:"ip"`
		MAC                    string `json:"mac"`
	} `json:"info"`
}
