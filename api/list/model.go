package list

import "time"

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

type DownloadFile struct {
	ID          int       `json:"id" gorm:"column:id;primaryKey"`
	Name        string    `json:"name" gorm:"column:name"`
	FileName    string    `json:"filename" gorm:"column:fname"`
	Remark      string    `json:"remark" gorm:"column:rname"`
	Size        int64     `json:"size" gorm:"column:size"`
	User        string    `json:"user" gorm:"column:user"`
	Time        time.Time `json:"time" gorm:"column:time"`
	DownloadNum int       `json:"download_num" gorm:"column:dtimes"`
}

func (DownloadFile) TableName() string {
	return "list_file"
}
