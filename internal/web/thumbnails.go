package web

import (
	"net/http"
	"strings"
)

func init() {
	routes.register(route{
		Path: "/thumbnail/",
		Handler: logger(func(w http.ResponseWriter, r *http.Request) {
			switch r.Method {
			case "GET":
				id := strings.TrimPrefix(r.URL.Path, "/thumbnail/")

				for _, song := range MusicCollection {
					if song.ID == id {
						// TODO: Write the MIME type? Faster lookup with a cache?
						w.Write(song.Tag.Picture.Data)
					}
				}
			default:
				http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			}
		}),
	})
}
