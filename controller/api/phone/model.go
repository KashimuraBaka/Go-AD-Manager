package phone

type PhoneInfo struct {
	ID       int    `json:"id" gorm:"column:id"`
	Name     string `json:"name" gorm:"column:name"`
	Phone    string `json:"phone" gorm:"column:phone"`
	RecordIP string `json:"record_ip" gorm:"column:record_ip"`
}

func (PhoneInfo) TableName() string {
	return "client_phone_list"
}
