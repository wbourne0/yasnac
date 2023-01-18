/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>
*/
package main

import (
	// "bytes"
	// "fmt"
	// "strings"

	// "os"

	"yasnac/cmd"
	// "yasnac/lib"
)

const str = `/JOB
//NAME GOTEST
//POS
///NPOS 3,0,0,0
///TOOL 0
///PULSE
C000=-4720,45472,-31419,-53,8048,2945
C001=-4720,49552,-31419,-53,8048,2945
C002=-1264,43504,-30699,-53,8048,2945
//INST
///DATE 2014/11/04 13:48
///ATTR 0,0,0,0,0,0,0,0,0,0,0,0,0,1,0,0
///FRAME BASE
NOP
*1
MOVL C000 V=11.7 CONT
MOVJ C001 VJ=1.00 CONT
MOVJ C002 VJ=0.78 CONT
MOVJ C000 VJ=0.78 CONT
JUMP *1
END
`

func main() {
	cmd.Execute()
	// dat, err := lib.NewERCMachine("/dev/ttyUSB0")
	// // dat, err := lib.NewERCMachine("/dev/serial/by-id/usb-Prolific_Technology_Inc._USB-Serial_Controller_D-if00-port0")

	// fmt.Println("opened")

	// if err != nil {
	// 	panic(err)
	// }

	// // if err := dat.WriteFile("GOTEST", bytes.NewBufferString(strings.ReplaceAll(str, "\n", "\r\n"))); err != nil {
	// // 	panic(err)
	// // }

	// // if err := dat.ReadFileTo("GOTEST", os.Stdout); err != nil {
	// // 	panic(err)
	// // }

	// b, err := dat.RunCommand("RJDIR *")
	// fmt.Println(string(b), err)

	// // b, err = dat.RunCommand("DELETE GOTEST")
	// // fmt.Println(string(b), err)

	// // b, err = dat.RunCommand("SVON 1")
	// // fmt.Println(string(b), err)

	// b, err = dat.RunCommand("RSTATS")
	// fmt.Println(string(b), err)
	// b, err = dat.RunCommand("RPOS")
	// fmt.Println(string(b), err)
}
