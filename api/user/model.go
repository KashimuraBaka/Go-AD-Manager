package user

import "time"

type UserRecord struct {
	IP       string    `json:"ip" gorm:"column:ip;primarykey"`
	ReadTime time.Time `json:"readtime" gorm:"column:readtime"`
	ReadNum  int       `json:"readnum" gorm:"column:readnum"`
	Role     int       `json:"role" gorm:"column:role"`
}

func (UserRecord) TableName() string {
	return "user_record"
}

type UserLogon struct {
	UserName  string    `json:"username" gorm:"column:user;primarykey"`
	PassWord  string    `json:"password" gorm:"column:pwd"`
	LogonTime time.Time `json:"logon_time" gorm:"column:logontime"`
	IP        string    `json:"ip" gorm:"column:ip"`
}

func (UserLogon) TableName() string {
	return "user_logon"
}
