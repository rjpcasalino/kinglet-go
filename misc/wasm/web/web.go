// A basic HTTP server.
// By default, it serves the current working directory on port 8080.
package main

import (
	"flag"
	"log"
	"net/http"
	"os"
	"os/exec"
	"runtime"
)

var (
	listen = flag.String("listen", ":8080", "listen address")
	dir    = flag.String("dir", ".", "directory to serve")
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
		go openBrowser("http://localhost" + *listen)
	}

	err = http.ListenAndServe(*listen, http.FileServer(http.Dir(*dir)))
	log.Fatalln(err)
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
		log.Printf("open browser: %v", err)
	}
}
