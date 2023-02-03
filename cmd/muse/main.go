package main

import (
	"flag"
	"log"
	"os"
	"strings"

	"github.com/prophittcorey/muse/internal/web"
)

func setenv(key, value string) {
	if err := os.Setenv(key, value); err != nil {
		log.Printf("muse: failed to set a default environment variable; %s", err)
	}
}

func main() {
	var (
		help   bool
		port   string
		host   string
		domain string
		globs  string
		dir    string
		auth   string
	)

	flag.BoolVar(&help, "h", false, "Displays the program's usage")
	flag.StringVar(&port, "port", "3000", "The port to run the server on (default: 3000)")
	flag.StringVar(&host, "host", "127.0.0.1", "The host to run the server on (default: 127.0.0.1)")
	flag.StringVar(&domain, "domain", "localhost", "The domain name for the server (default: localhost)")
	flag.StringVar(&dir, "dir", ".", "A base directory to base the glob patterns from (default: .)")
	flag.StringVar(&globs, "globs", "*.mp3,**/*.mp3,**/**/*.mp3", "A comma separated list of glob patterns (default: **/*.mp3)")
	flag.StringVar(&auth, "auth", "", "A username and password to use for basic auth (example: admin:qwerty)")

	flag.Parse()

	if help {
		flag.Usage()
		return
	}

	setenv("HOST", host)
	setenv("PORT", port)
	setenv("DOMAIN", domain)
	setenv("BASIC_AUTH", auth)

	web.SetAuth(auth)

	if err := web.Serve(dir, strings.Split(globs, ",")...); err != nil {
		log.Fatal(err)
	}
}
