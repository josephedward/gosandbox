package core

import (
	"fmt"
)

var Reset = "\033[0m"
var Red = "\033[31m"
var Green = "\033[32m"
var Yellow = "\033[33m"
var Blue = "\033[34m"
var Purple = "\033[35m"
var Cyan = "\033[36m"
var Gray = "\033[37m"
var White = "\033[97m"

func PrintIfErr(err error) {
	if err != nil {
		fmt.Println(Red + err.Error() + Reset)
	}
}

func Success(message ...interface{}) {
	for _, msg := range message {
		s, ok := msg.(string) // the "ok" boolean will flag success.
		if ok {
			fmt.Print(Green + string(s) + Reset)
		} else {
			fmt.Println(msg)
		}
	}
}
