package web

import (
	"log"
	"net"
	"net/http"
	"strings"
	"time"
)

func neuter(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if strings.HasSuffix(r.URL.Path, "/") {
			http.NotFound(w, r)
			return
		}

		next.ServeHTTP(w, r)
	})
}

func getIP(r *http.Request) string {
	remoteIP := r.RemoteAddr

	if forwarded := r.Header["X-Forwarded-For"]; len(forwarded) != 0 {
		remoteIP = forwarded[0]
	}

	host, _, err := net.SplitHostPort(remoteIP)

	if err == nil {
		remoteIP = host
	}

	if strings.Contains(remoteIP, ",") {
		parts := strings.Split(remoteIP, ",")

		if len(parts) > 0 {
			return strings.TrimSpace(parts[0])
		}
	}

	return remoteIP
}

func logger(h http.HandlerFunc) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		ip := getIP(r)
		now := time.Now()

		log.Printf("Started %s %s %s\n", ip, r.Method, r.URL)
		h(w, r)
		log.Printf("Completed %s %s %s in %s\n", ip, r.Method, r.URL, time.Since(now))
	}
}
