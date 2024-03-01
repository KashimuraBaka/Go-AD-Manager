package mysql

var DB *DataBase

func init() {
	db, err := CrateDatabase(MysqlConnInfo{
		Address:  "sev.kashimura.cn",
		Port:     3306,
		UserName: "Cirno",
		PassWord: "Baka@99999",
		DataBase: "web",
	})
	if err != nil {
		// panic(err)
	}
	DB = db
}
