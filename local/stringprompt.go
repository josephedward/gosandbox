package local

import (
    "bufio"
    "fmt"
    "os"
    "strings"
)

// StringPrompt asks for a string value using the label
func StringPrompt(label string) string {
    var s string
    r := bufio.NewReader(os.Stdin)
    for {
        fmt.Fprint(os.Stderr, label+" ")
        s, _ = r.ReadString('\n')
        if s != "" {
            break
        }
    }
    return strings.TrimSpace(s)
}

// func main() {
//     name := StringPrompt("What is your name?")
//     fmt.Printf("Hello, %s!\n", name)
// }


func PanicIfErr(err error) {
	if err != nil {
		fmt.Println(err)
		panic(err)
	}
}