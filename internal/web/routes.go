package web

import (
	"log"
	"net/http"
	"strings"

	"github.com/prophittcorey/muse"
)

func init() {
	// The index page handler. This is a bit special because it handles the index
	// page ("/") and any pages that don't match a registered route (serves as the
	// catch all handler).
	routes.register(route{
		Path: "/",
		Handler: logger(func(w http.ResponseWriter, r *http.Request) {
			switch r.Method {
			case "GET":
				if r.URL.Path != "/" {
					// NOTE: This is the catch all code path. We could do things here like
					// redirect old broken links, render a 404 page, etc.

					http.NotFound(w, r)
				} else {
					data := map[string]interface{}{
						"Tracks":       tracks,
						"DefaultTrack": tracks[0],
					}

					if err := templates.ExecuteTemplate(w, "pages/index.tmpl", data); err != nil {
						log.Printf("web: error rendering index page; %s", err)
					}
				}
			default:
				http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			}
		}),
	})

	// The thumbnails route is used to load embedded image assets if available. If one is not
	// available a default image will be used.
	routes.register(route{
		Path: "/thumbnail/",
		Handler: logger(func(w http.ResponseWriter, r *http.Request) {
			switch r.Method {
			case "GET":
				id := strings.TrimPrefix(r.URL.Path, "/thumbnail/")

				for _, song := range tracks {
					if song.ID == id {
						// TODO: Write the MIME type? Faster lookup with a cache?
						w.Write(song.Tag.Picture.Data)
					}
				}

				// TODO: Load a default image if none is available.
			default:
				http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			}
		}),
	})

	// Handles all asset requests. Since we're embedding all assets into the
	// binary we need to serve the assets from Go. If we were not, we could
	// directly access the assets via Nginx/Apache and leave the application
	// server alone. NOTE: Nginx/Apache can still cache the assets after they
	// leave the application server.
	routes.register(route{
		Path: "/assets/",
		Handler: func() func(http.ResponseWriter, *http.Request) {
			return logger(func(w http.ResponseWriter, r *http.Request) {
				neuter(http.FileServer(http.FS(muse.FS))).ServeHTTP(w, r)
			})
		}(),
	})
}