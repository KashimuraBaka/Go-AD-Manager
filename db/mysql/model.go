package mysql

type DomainUser struct {
	Name string `gorm:"column:name"`
	IP   string `gorm:"column:ip"`
	MAC  string `gorm:"column:mac"`
}
