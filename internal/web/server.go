package web

import (
	"context"
	"crypto/sha256"
	"crypto/subtle"
	"fmt"
	"html/template"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/prophittcorey/muse"
	"github.com/prophittcorey/muse/internal/audio"
)

type route struct {
	Path    string
	Handler http.HandlerFunc
}

type routecollection []route

func (rs *routecollection) register(r route) {
	*rs = append(*rs, r)
}

var (
	authRequired         bool
	expectedUsernameHash [32]byte
	expectedPasswordHash [32]byte
)

func SetAuth(auth string) {
	if before, after, ok := strings.Cut(auth, ":"); ok {
		expectedUsernameHash = sha256.Sum256([]byte(before))
		expectedPasswordHash = sha256.Sum256([]byte(after))

		authRequired = true
	}
}

func checkAuth(w http.ResponseWriter, r *http.Request) bool {
	if !authRequired {
		return true
	}

	if username, password, ok := r.BasicAuth(); ok {
		usernameHash := sha256.Sum256([]byte(username))
		passwordHash := sha256.Sum256([]byte(password))

		usernameMatch := (subtle.ConstantTimeCompare(usernameHash[:], expectedUsernameHash[:]) == 1)
		passwordMatch := (subtle.ConstantTimeCompare(passwordHash[:], expectedPasswordHash[:]) == 1)

		if usernameMatch && passwordMatch {
			return true
		}
	}

	w.Header().Set("WWW-Authenticate", `Basic realm="restricted", charset="UTF-8"`)
	http.Error(w, "Unauthorized", http.StatusUnauthorized)

	return false
}

// Serve will load audio from a directory using any number of file globs. This
// function blocks.
func Serve(directory string, globs ...string) error {
	patterns := []string{}

	for _, glob := range globs {
		patterns = append(patterns, fmt.Sprintf("%s/%s", directory, glob))
	}

	if !audio.Scan(patterns...) {
		return fmt.Errorf("web: no music files were found")
	}

	return listenAndServe()
}

func listenAndServe() error {
	log.SetFlags(log.Flags() &^ (log.Ldate | log.Ltime))

	mux := http.NewServeMux()

	for _, r := range routes {
		mux.HandleFunc(r.Path, r.Handler)
	}

	srv := &http.Server{
		Addr:              fmt.Sprintf("%s:%s", os.Getenv("HOST"), os.Getenv("PORT")),
		ReadTimeout:       3 * time.Second,
		ReadHeaderTimeout: 3 * time.Second,
		WriteTimeout:      5 * time.Second,
		IdleTimeout:       30 * time.Second,
		Handler:           mux,
	}

	log.Printf("Listening on %s:%s\n", os.Getenv("DOMAIN"), os.Getenv("PORT"))

	// Run our server in a goroutine so that it doesn't block.
	go func() {
		if err := srv.ListenAndServe(); err != nil {
			if err != http.ErrServerClosed {
				log.Printf("web: srv.ListenAndServe returned an error; %s\n", err)
			}
		}
	}()

	c := make(chan os.Signal, 1)

	// We'll accept graceful shutdowns when quit via SIGINT (Ctrl+C)
	// SIGKILL, SIGQUIT or SIGTERM (Ctrl+/) will not be caught.
	signal.Notify(c, os.Interrupt, syscall.SIGTERM, syscall.SIGQUIT)

	// Block until we receive our signal.
	<-c

	log.Println("Server is shutting down.")

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)

	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		return fmt.Errorf(`web: failed to shut down with grace; %s`, err)
	}

	log.Println("Server has exited with grace.")

	return nil /* if we made it this far all is well */
}

var (
	templates    *template.Template
	routes       = routecollection{}
	inlinedCache = map[string]template.HTML{}
)

func setdefault(key, value string) {
	if v := os.Getenv(key); len(v) == 0 {
		if err := os.Setenv(key, value); err != nil {
			log.Printf("web: failed to set a default environment variable; %s", err)
		}
	}
}

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

func init() {
	setdefault("HOST", "127.0.0.1")
	setdefault("PORT", "3000")
	setdefault("DOMAIN", "localhost")

	tmpls := []string{
		"templates/pages/*.tmpl",
		"templates/partials/*.tmpl",
	}

	templates = template.New("").Funcs(template.FuncMap{
		"app_name": func() string {
			return "Muse"
		},
		"ran_at": func() int64 {
			return muse.RanAt
		},
		"embed": func(t string, name string) template.HTML {
			if templateHTML, ok := inlinedCache[t+name]; ok {
				return templateHTML
			}

			bs, err := muse.FS.ReadFile("assets/" + t + "/" + name)

			if err != nil {
				log.Fatal(err)
			}

			// #nosec
			inlined := template.HTML(string(bs))

			inlinedCache[t+name] = inlined

			return inlined
		},
	})

	if _, err := templates.ParseFS(muse.FS, tmpls...); err != nil {
		log.Fatal(err)
	}
}
