package core

import (
	"os"
	"os/exec"
	"github.com/go-rod/rod"
	"github.com/go-rod/rod/lib/launcher"
	"github.com/go-rod/rod/lib/utils"
	"github.com/ysmood/leakless"
	"flag"
	"fmt"
	"log"
	"net"
	"net/http"
)

func CustomLaunch() *rod.Browser{
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


func BrowserCliOutput() *rod.Browser{
	// Pipe the browser stderr and stdout to os.Stdout .
	u := launcher.New().Logger(os.Stdout).MustLaunch()
	return rod.New().ControlURL(u).MustConnect()
}

func Manual(u string) *rod.Browser{
	return rod.New().ControlURL(u).MustConnect()
}


var addr = flag.String("address", "localhost:7317", "the address to listen to")
var quiet = flag.Bool("quiet", false, "silence the log")
var allowAllPath = flag.Bool("allow-all", false, "allow all path set by the client")
// A server to help launch browser remotely
func manager() {
	flag.Parse()

	m := launcher.NewManager()

	if !*quiet {
		m.Logger = log.New(os.Stdout, "", 0)
	}

	if *allowAllPath {
		m.BeforeLaunch = func(l *launcher.Launcher, rw http.ResponseWriter, r *http.Request) {}
	}

	l, err := net.Listen("tcp", *addr)
	if err != nil {
		utils.E(err)
	}

	if !*quiet {
		fmt.Println("rod-manager listening on:", l.Addr().String())
	}

	srv := &http.Server{Handler: m}
	utils.E(srv.Serve(l))
}
