package web

import (
	"log"
	"net/http"
	"strings"

	"github.com/prophittcorey/muse"
	"github.com/prophittcorey/muse/internal/audio"
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
						"Tracks":       audio.Tracks.All,
						"DefaultTrack": audio.Tracks.All[0],
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

	defaultAlbumArt, err := muse.FS.ReadFile("assets/images/album.png")

	if err != nil {
		log.Printf("routes.go: failed to load default album art; %s\n", err)
	}

	// The thumbnails route is used to load embedded image assets if available. If one is not
	// available a default image will be used.
	routes.register(route{
		Path: "/thumbnail/",
		Handler: logger(func(w http.ResponseWriter, r *http.Request) {
			switch r.Method {
			case "GET":
				id := strings.TrimPrefix(r.URL.Path, "/thumbnail/")

				if t := audio.Tracks.Find(id); t != nil && t.Tag.Picture != nil {
					w.Write(t.Tag.Picture.Data)
					return
				}

				w.Write(defaultAlbumArt)
			default:
				http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			}
		}),
	})

	// The track route will stream the actual music file given the ID.
	routes.register(route{
		Path: "/track/",
		Handler: logger(func(w http.ResponseWriter, r *http.Request) {
			switch r.Method {
			case "GET":
				id := strings.TrimPrefix(r.URL.Path, "/track/")

				if t := audio.Tracks.Find(id); t != nil {
					http.ServeFile(w, r, t.Path)
					return
				}

				http.NotFound(w, r)
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
