// A basic HTTP server.
// By default, it serves the current working directory on port 8080.
package main

import (
	"flag"
	"log"
	"net"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
)

var (
	defaultDir = func() string {
		exe, err := os.Executable()
		if err != nil {
			return "."
		}
		return filepath.Dir(exe)
	}()
	listen = flag.String("listen", ":8080", "listen address")
	dir    = flag.String("dir", defaultDir, "directory to serve")
	open   = flag.Bool("open", true, "open a browser once the server is running")
)

func main() {
	flag.Parse()
	absDir, err := os.Getwd()
	if err == nil {
		log.Printf("serving %s", absDir)
	}
	log.Printf("listening on %q...", *listen)

	if *open {
		go openBrowser(browserURL(*listen))
	}

	err = http.ListenAndServe(*listen, http.FileServer(http.Dir(*dir)))
	log.Fatalln(err)
}

func browserURL(addr string) string {
	addr = strings.TrimSpace(addr)

	if !strings.Contains(addr, "://") && strings.Contains(addr, ":") {
		host, port, err := net.SplitHostPort(addr)
		if err == nil {
			if host == "" || host == "0.0.0.0" || host == "::" || host == "[::]" {
				host = "localhost"
			}
			return "http://" + host + ":" + port
		}
	}

	if !strings.HasPrefix(addr, "http://") && !strings.HasPrefix(addr, "https://") {
		return "http://" + addr
	}
	return addr
}

func openBrowser(url string) {
	var cmd *exec.Cmd

	switch runtime.GOOS {
	case "darwin":
		cmd = exec.Command("open", url)
	case "windows":
		cmd = exec.Command("rundll32", "url.dll,FileProtocolHandler", url)
	default:
		cmd = exec.Command("xdg-open", url)
	}

	if err := cmd.Start(); err != nil {
		log.Printf("could not automatically open browser: %v", err)
	}
}
