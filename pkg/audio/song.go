package audio

import (
	"fmt"
	"io"
	"os"
)

type Header struct {
	Tag      string
	Revision int
	Minor    int
	Flags    int
	Size     int
}

func (m Header) Version() string {
	return fmt.Sprintf("%sv2.%d.%d\n", m.Tag, m.Revision, m.Minor)
}

type Song struct {
	Path      string
	Title     string
	Artist    string
	Album     string
	Thumbnail []byte
	Header    Header
}

func (s Song) Load() error {
	f, err := os.Open(s.Path)

	if err != nil {
		return err
	}

	defer f.Close()

	bs := make([]byte, 10)

	if _, err = io.ReadFull(f, bs); err != nil {
		return err
	}

	s.Header = Header{
		Tag:      string(bs[0:3]),
		Revision: int(bs[3]),
		Minor:    int(bs[4]),
	}

	// TODO: Parse the flags and interpret them.
	// TODO: Parse the size.

	return nil
}
