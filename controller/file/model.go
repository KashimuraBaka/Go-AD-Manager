package file

import "time"

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
