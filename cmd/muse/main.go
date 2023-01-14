package main

import (
	"flag"
	"log"
	"os"

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
		dir    string
	)

	flag.BoolVar(&help, "h", false, "Displays the program's usage")
	flag.StringVar(&port, "port", "3000", "The port to run the server on (default: 3000)")
	flag.StringVar(&host, "host", "127.0.0.1", "The host to run the server on (default: 127.0.0.1)")
	flag.StringVar(&domain, "domain", "localhost", "The domain name for the server (default: localhost)")
	flag.StringVar(&dir, "dir", ".", "The directory to scan for mp3 files (default: .)")

	flag.Parse()

	if help {
		flag.Usage()
		return
	}

	setenv("HOST", host)
	setenv("PORT", port)
	setenv("DOMAIN", domain)

	if err := web.ListenAndServe(); err != nil {
		log.Fatal(err)
	}
}
