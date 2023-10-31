package mdb

type UserInfo struct {
	UserID int    `json:"userid"`
	Name   string `json:"name"`
	Title  string `json:"title"`
}

type DriverInfo struct {
	ID           int    `json:"id"`
	MachineAlias string `json:"machine_alias"`
	SN           string `json:"sn"`
}

type AttendanceData struct {
	Users   []UserInfo
	Drivers []DriverInfo
}
