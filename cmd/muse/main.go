package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/prophittcorey/muse/internal/web"
	"github.com/prophittcorey/muse/pkg/audio"
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
	)

	flag.BoolVar(&help, "h", false, "Displays the program's usage")
	flag.StringVar(&port, "port", "3000", "The port to run the server on (default: 3000)")
	flag.StringVar(&host, "host", "127.0.0.1", "The host to run the server on (default: 127.0.0.1)")
	flag.StringVar(&domain, "domain", "localhost", "The domain name for the server (default: localhost)")
	flag.StringVar(&dir, "dir", ".", "A base directory to base the glob patterns from (default: .)")
	flag.StringVar(&globs, "globs", "**/*.mp3", "A comma separated list of glob patterns (default: **/*.mp3)")

	flag.Parse()

	if help {
		flag.Usage()
		return
	}

	setenv("HOST", host)
	setenv("PORT", port)
	setenv("DOMAIN", domain)

	patterns := []string{}

	for _, glob := range strings.Split(globs, ",") {
		patterns = append(patterns, fmt.Sprintf("%s/%s", dir, glob))
	}

	web.MusicCollection = audio.Scan(patterns...)

	if len(web.MusicCollection) == 0 {
		log.Fatalf("err: no music files were found for %s", globs)
		return
	}

	if err := web.ListenAndServe(); err != nil {
		log.Fatal(err)
	}
}
