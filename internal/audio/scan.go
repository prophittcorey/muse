package audio

import (
	"log"
	"path/filepath"
)

func Scan(globs ...string) []*Song {
	added := map[string]struct{}{}

	files := []*Song{}

	for _, glob := range globs {
		if matches, err := filepath.Glob(glob); err == nil {
			for _, match := range matches {
				match := match

				if _, ok := added[match]; ok {
					continue
				}

				added[match] = struct{}{}

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
