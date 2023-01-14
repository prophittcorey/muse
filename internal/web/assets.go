package web

import (
	"net/http"

	"github.com/prophittcorey/muse"
)

func init() {
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
