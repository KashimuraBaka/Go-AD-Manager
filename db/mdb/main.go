package mdb

var DB *DataBase

func init() {
	/* absPath, err := filepath.Abs(path.Join("static", "att2000.mdb"))
	if err != nil {
		return
	} */
	DB = &DataBase{Path: "F:\\DataBase\\att2000.mdb"}
}
