package pshell

import "fmt"

var Shell *ADControl

func init() {
	ps, err := CreateADControl()
	if err != nil {
		fmt.Println("[Powershell] Error:", err)
	}

	if err = ps.CreateSession("192.168.102.254", "Administrator", "RGrr2019"); err != nil {
		panic(err)
	}

	ps.GetUsers("")

	Shell = ps
}
