package mdb

import (
	"database/sql"
	"errors"
	"fmt"

	_ "github.com/alexbrainman/odbc"
)

type DataBase struct {
	Path string
}

func (db DataBase) Connect() (*sql.DB, error) {
	return sql.Open("odbc", fmt.Sprintf("driver={Microsoft Access Driver (*.mdb, *.accdb)};dbq=%s", db.Path))
}

func (db DataBase) SelectUserInfo() (AttendanceData, error) {
	data := AttendanceData{}

	conn, err := db.Connect()
	if err != nil {
		return data, err
	}
	defer conn.Close()

	rows, err := conn.Query("SELECT USERID, NAME, TITLE FROM USERINFO")
	if err != nil {
		return data, err
	}

	for rows.Next() {
		var userid int
		var name sql.NullString
		var title sql.NullString
		if err := rows.Scan(&userid, &name, &title); err != nil {
			return data, err
		}
		data.Users = append(data.Users, UserInfo{
			UserID: userid,
			Name:   name.String,
			Title:  title.String,
		})
	}
	rows.Close()

	rows, err = conn.Query("SELECT ID, MachineAlias, sn FROM Machines")
	if err != nil {
		return data, err
	}

	for rows.Next() {
		var id int
		var ma sql.NullString
		var sn sql.NullString
		if err := rows.Scan(&id, &ma, &sn); err != nil {
			return data, err
		}
		data.Drivers = append(data.Drivers, DriverInfo{
			ID:           id,
			MachineAlias: ma.String,
			SN:           sn.String,
		})
	}

	return data, nil
}

func (db DataBase) CheckUserID(userid int) error {
	conn, err := db.Connect()
	if err != nil {
		return err
	}
	defer conn.Close()

	rows, err := conn.Query("SELECT * FROM USERINFO WHERE USERID = ?", userid)
	if err != nil {
		return err
	}
	defer rows.Close()

	if !rows.Next() {
		return errors.New("not found")
	}
	return nil
}

func (db DataBase) InsertCheckInOut(userid int, time string, device string) error {
	conn, err := db.Connect()
	if err != nil {
		return err
	}
	defer conn.Close()

	rows, err := conn.Exec(`
		INSERT INTO CHECKINOUT
		(USERID, CHECKTIME, VERIFYCODE, SENSORID, sn, UserExtFmt)
		VALUES (?, ?, 1, 1, ?, 1);
	`, userid, time, device)
	if err != nil {
		return err
	}
	if _, err = rows.RowsAffected(); err != nil {
		return err
	}
	return nil
}
