package audio

type Song struct {
	Title     string
	Artist    string
	Album     string
	Thumbnail []byte
	Path      string
}

func (s Song) Load() error {
	return nil
}
