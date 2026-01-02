package mysql

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

type PhoneInfo struct {
	ID       int    `json:"id" gorm:"column:id"`
	Name     string `json:"name" gorm:"column:name"`
	Phone    string `json:"phone" gorm:"column:phone"`
	RecordIP string `json:"record_ip" gorm:"column:record_ip"`
}

func (PhoneInfo) TableName() string {
	return "client_phone_list"
}

type DomainUser struct {
	Name string `json:"name" gorm:"column:name"`
	IP   string `json:"ip" gorm:"column:ip"`
	Mac  string `json:"mac" gorm:"column:mac"`
}

func (DomainUser) TableName() string {
	return "list_domain"
}

type SystemUrl struct {
	Name string `json:"name" gorm:"column:name"`
	Url  string `json:"url" gorm:"column:url"`
}

func (SystemUrl) TableName() string {
	return "list_url"
}

type FileInfo struct {
	ID          int64     `json:"id" gorm:"column:id;primarykey"`
	Name        string    `json:"name" gorm:"column:name"`
	FileName    string    `json:"fname" gorm:"column:fname"`
	ReName      string    `json:"rname" gorm:"column:rname"`
	Size        int64     `json:"size" gorm:"column:size"`
	User        string    `json:"user" gorm:"column:user"`
	Date        time.Time `json:"time" gorm:"column:time"`
	DownloadNum int       `json:"downloadnum" gorm:"column:downloadnum"`
}

func (FileInfo) TableName() string {
	return "list_file"
}

type ACUser struct {
	Userid      int       `json:"userid" gorm:"column:userid;primaryKey"`
	Badgenumber int       `json:"badgenid" gorm:"column:badgenumber"`
	Name        string    `json:"name" gorm:"column:name"`
	Group       int       `json:"group" gorm:"column:group"`
	Authority   int       `json:"authority" gorm:"column:authority"`
	Device      int       `json:"device" gorm:"column:device"`
	CreateTime  time.Time `json:"create_time" gorm:"column:create_time"`
	DeleteTime  time.Time `json:"delete_time" gorm:"column:delete_time"`
}

func (ACUser) TableName() string {
	return "ac_user"
}

type ACAuthority struct {
	Authority int    `json:"authority" gorm:"column:authority;primaryKey"`
	Name      string `json:"name" gorm:"column:name"`
}

func (ACAuthority) TableName() string {
	return "ac_authority"
}

type ACDevice struct {
	Device int    `json:"device" gorm:"column:device;primaryKey"`
	Name   string `json:"name" gorm:"column:name"`
	SN     string `json:"sn" gorm:"column:sn"`
}

func (ACDevice) TableName() string {
	return "ac_device"
}

type ACGroup struct {
	Group int    `json:"group" gorm:"column:group;primaryKey"`
	Name  string `json:"name" gorm:"column:name"`
}

func (ACGroup) TableName() string {
	return "ac_group"
}
