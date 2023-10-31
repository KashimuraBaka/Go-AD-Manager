package mdb

import (
	"path"
	"path/filepath"
)

var DB *DataBase

func init() {

	absPath, err := filepath.Abs(path.Join("static", "att2000.mdb"))
	if err != nil {
		return
	}
	DB = &DataBase{Path: absPath}
}
