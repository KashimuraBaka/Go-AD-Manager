package mysql

import (
	"fmt"
	"sync"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type MysqlConnInfo struct {
	Address  string
	Port     int
	UserName string
	PassWord string
	DataBase string
}

type DataBase struct {
	sync.RWMutex
	*gorm.DB
}

func CrateDatabase(c MysqlConnInfo) (*DataBase, error) {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local", c.UserName, c.PassWord, c.Address, c.Port, c.DataBase)
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		QueryFields: true, // 打印SQL语句
		// SkipDefaultTransaction: true, // 禁用事务
	})
	if err != nil {
		return nil, err
	}
	return &DataBase{DB: db}, nil
}

func (db *DataBase) Clear(tableName string) error {
	return db.Exec("DELETE FROM " + tableName).Error
}
