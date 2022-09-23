package cli

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
		Error(err.Error())
	}
}

func Success(message ...interface{}) {
	for _, msg := range message {
		s, ok := msg.(string) // the "ok" boolean will flag success.
		if ok {
			fmt.Println(Green + string(s) + Reset)
		} else {
			fmt.Println(msg)
		}
	}
}

func Error(message ...interface{}) {
	for _, msg := range message {
		s, ok := msg.(string) // the "ok" boolean will flag success.
		if ok {
			fmt.Println(Red + string(s) + Reset)
		} else {
			fmt.Println(msg)
		}
	}
}


																									   

func Welcome() {	
	fmt.Println(Green+"---------------------------------------------------------------------------------------------------------------"+Reset)
	fmt.Println(Cyan+" .oooooo.               .oooooo..o                             .o8  oooooooooo.                          "+Reset)
	fmt.Println(Cyan+" d8P'  'Y8b             d8P'    'Y8                             888  '888'   'Y8b                        "+Reset)
	fmt.Println(Cyan+" 888            .ooooo.  Y88bo.       .oooo.   ooo. .oo.    .oooo888   888     888  .ooooo.  oooo    ooo "+Reset)
	fmt.Println(Cyan+" 888           d88' '88b   Y8888o.    P  )88b   888PY88b   d88   888   888oooo888  d88   88b   88b..8P   "+Reset)
	fmt.Println(Cyan+" 888     ooooo 888   888      ''Y88b  .oP'888   888   888  888   888   888    '88b 888   888    Y888'    "+Reset)
	fmt.Println(Cyan+" 88.    .88'  888   888 oo     .d8P d8(  888   888   888  888   888   888    .88P 888   888  .o8''88b    "+Reset)
	fmt.Println(Cyan+" Y8bood8P'   'Y8bod8P' 8''88888P'  'Y888''8o o888o o888o 'Y8bod88P' o888bood8P'  'Y8bod8P' o88'   888o   "+Reset)
	fmt.Println(Green+"---------------------------------------------------------------------------------------------------------------"+Reset)
	fmt.Println(Yellow +"    ACloudGuru Sandbox Procurement Tool" + Reset)
	fmt.Println(Green+"---------------------------------------------------------------------------------------------------------------"+Reset)

}
