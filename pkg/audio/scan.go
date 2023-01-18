package audio

import (
	"log"
	"path/filepath"
)

// TODO: Make sure the songs are unique...
func Scan(globs ...string) []*Song {
	files := []*Song{}

	for _, glob := range globs {
		if matches, err := filepath.Glob(glob); err == nil {
			for _, match := range matches {
				match := match

				song := Song{Path: match}

				if err := song.Load(); err == nil {
					files = append(files, &song)
				} else {
					log.Printf("failed to load %s; %s\n", match, err)
				}
			}
		}
	}

	return files
}
