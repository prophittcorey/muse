package audio

import (
	"log"
	"path/filepath"
)

var (
	Tracks = &TrackCollection{
		lookup: map[string]*Track{},
	}
)

type TrackCollection struct {
	lookup map[string]*Track
	All    []*Track
}

// Insert adds a track to the collection.
func (t *TrackCollection) Insert(audio *Track) {
	if track := t.Find(audio.ID); track == nil {
		Tracks.lookup[audio.ID] = audio
		Tracks.All = append(Tracks.All, audio)
	}
}

// Find returns a track via it's ID. This is an O(1) lookup. Nil is returned
// if no track is found.
func (t *TrackCollection) Find(id string) *Track {
	if track, ok := t.lookup[id]; ok {
		return track
	}

	return nil
}

// Scan looks for and loads audio tracks. Returns true if any tracks were
// found and loaded.
func Scan(globs ...string) bool {
	added := map[string]struct{}{}

	for _, glob := range globs {
		if matches, err := filepath.Glob(glob); err == nil {
			for _, match := range matches {
				match := match

				if _, ok := added[match]; ok {
					continue
				}

				added[match] = struct{}{}

				song := Track{Path: match}

				if err := song.Load(); err == nil {
					Tracks.Insert(&song)
				} else {
					log.Printf("failed to load %s; %s\n", match, err)
				}
			}
		}
	}

	return len(Tracks.All) > 0
}
