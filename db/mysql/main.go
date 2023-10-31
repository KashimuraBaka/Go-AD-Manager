package mysql

var DB *DataBase

func init() {
	db, err := CrateDatabase(MysqlConnInfo{
		Address:  "sev.kashimura.cc",
		Port:     1551,
		UserName: "Cirno",
		PassWord: "Baka@9999",
		DataBase: "web",
	})
	if err != nil {
		panic(err)
	}
	DB = db
}
