package core

import (
	"flag"
	"fmt"
	"github.com/go-rod/rod"
	"github.com/go-rod/rod/lib/launcher"
	"github.com/go-rod/rod/lib/utils"
	"github.com/ysmood/leakless"
	"log"
	"net"
	"net/http"
	"os"
	"os/exec" 
)

func CustomLaunch() *rod.Browser {
	// get the browser executable path
	path := launcher.NewBrowser().MustGet()

	// use the FormatArgs to construct args, this line is optional, you can construct the args manually
	args := launcher.New().FormatArgs()

	var cmd *exec.Cmd
	if true { // decide whether to use leakless or not
		cmd = leakless.New().Command(path, args...)
	} else {
		cmd = exec.Command(path, args...)
	}

	parser := launcher.NewURLParser()
	cmd.Stderr = parser
	utils.E(cmd.Start())
	u := launcher.MustResolveURL(<-parser.URL)

	return rod.New().ControlURL(u).MustConnect()
}

func UseSystemBrowser() *rod.Browser {
	if path, exists := launcher.LookPath(); exists {
		u := launcher.New().Bin(path).MustLaunch()
		return rod.New().ControlURL(u).MustConnect()
	}
	return nil
}

func BrowserCliOutput() *rod.Browser {
	// Pipe the browser stderr and stdout to os.Stdout .
	u := launcher.New().Logger(os.Stdout).MustLaunch()
	return rod.New().ControlURL(u).MustConnect()
}

func Manual(u string) *rod.Browser {
	return rod.New().ControlURL(u).MustConnect()
}

var addr = flag.String("address", "localhost:7317", "the address to listen to")
var quiet = flag.Bool("quiet", false, "silence the log")
var allowAllPath = flag.Bool("allow-all", false, "allow all path set by the client")

// A server to help launch browser remotely
func Manager() {
	frameShell()
	flag.Parse()

	m := launcher.NewManager()

	// if !*quiet {
		m.Logger = log.New(os.Stdout, "", 0)
	// }

	if *allowAllPath {
		m.BeforeLaunch = func(l *launcher.Launcher, rw http.ResponseWriter, r *http.Request) {}
	}

	listen, err := net.Listen("tcp", *addr)
	if err != nil {
		utils.E(err)
	}

	// if !*quiet {
		fmt.Println("rod-manager listening on:", listen.Addr().String())
	// }

	srv := &http.Server{Handler: m}
	utils.E(srv.Serve(listen))

}

// func Remote() (*rod.Browser) {
// 	frameShell()
// 	// This example is to launch a browser remotely, not connect to a running browser remotely,
// 	// to connect to a running browser check the "../connect-browser" example.
// 	// Rod provides a docker image for beginers, run the below to start a launcher.Manager:
// 	//
// 	//     docker run -p 7317:7317 ghcr.io/go-rod/rod
// 	//
// 	// For more information, check the doc of launcher.Manager
// 	l := launcher.MustNewManaged("")

// 	// You can also set any flag remotely before you launch the remote browser.
// 	// Available flags: https://peter.sh/experiments/chromium-command-line-switches
// 	l.Set("disable-gpu").Delete("disable-gpu")

// 	// Launch with headful mode
// 	l.Headless(false).XVFB("--server-num=5", "--server-args=-screen 0 1600x900x16")

// 	browser := rod.New().Client(l.MustClient()).MustConnect()

// 	// You may want to start a server to watch the screenshots of the remote browser.
// 	launcher.Open(browser.ServeMonitor(""))

// 	fmt.Println(
// 		browser.MustPage("https://mdn.dev/").MustEval("() => document.title"),
// 	)
// 	// utils.Pause()
// 	return browser
// }



func frameShell(){
	fmt.Println("Before shell script:")
    cmd := exec.Command("bash", "-c", "./scripts/frame.sh")
    cmd.Stdin = os.Stdin
    cmd.Stdout = os.Stdout
    cmd.Stderr = os.Stderr
    _ = cmd.Run() // add error checking
    fmt.Println("After shell script")
}
