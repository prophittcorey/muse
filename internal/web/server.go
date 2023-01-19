package web

import (
	"context"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/prophittcorey/muse"
	"github.com/prophittcorey/muse/pkg/audio"
)

var (
	Tracks []*audio.Song
)

type route struct {
	Path    string
	Handler http.HandlerFunc
}

type routecollection []route

func (rs *routecollection) register(r route) {
	*rs = append(*rs, r)
}

// ListenAndServce begins listening and responding to web requests. This
// function blocks.
func ListenAndServe() error {
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
	templates *template.Template
	routes    = routecollection{}
)

func setdefault(key, value string) {
	if v := os.Getenv(key); len(v) == 0 {
		if err := os.Setenv(key, value); err != nil {
			log.Printf("web: failed to set a default environment variable; %s", err)
		}
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
	})

	if _, err := templates.ParseFS(muse.FS, tmpls...); err != nil {
		log.Fatal(err)
	}
}
