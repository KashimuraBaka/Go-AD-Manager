package pshell

import (
	"bytes"
	"fmt"
	"io"

	"golang.org/x/text/encoding/simplifiedchinese"
	"golang.org/x/text/transform"
)

var Shell *ADControl

func init() {
	ps, err := CreateADControl()
	if err != nil {
		fmt.Println("[Powershell] Error:", err)
	}

	if err = ps.CreateSession("192.168.102.254", "Administrator", "RGrr2019"); err != nil {
		b := []byte(err.Error())
		reader := transform.NewReader(bytes.NewReader(b), simplifiedchinese.GBK.NewDecoder())
		d, _ := io.ReadAll(reader)
		fmt.Println("[Powershell] Error", string(d))
	}

	Shell = ps
}
