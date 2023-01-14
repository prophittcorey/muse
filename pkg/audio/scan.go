package audio

import (
	"path/filepath"
)

func Scan(globs ...string) []File {
	files := []File{}

	for _, glob := range globs {
		if matches, err := filepath.Glob(glob); err == nil {
			for _, match := range matches {
				match := match
				files = append(files, File{Path: match})
			}
		}
	}

	return files
}
