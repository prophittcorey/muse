package audio

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"unicode/utf16"
	"unicode/utf8"
)

const (
	Unsynchronization = 1 << 7
	ExtendedHeader    = 1 << 6
	Experimental      = 1 << 5
	FooterPresent     = 1 << 4
)

type Tag struct {
	Header    Header
	Title     string
	Artist    string
	Album     string
	Thumbnail []byte
}

func (t Tag) String() string {
	return fmt.Sprintf("%s: %d bytes; extended %v", t.Header.Version(), t.Header.Size, t.Header.Flag(ExtendedHeader))
}

func (t *Tag) ParseFrames(r io.Reader) error {
	/*
		Frame ID      $xx xx xx xx  (four characters)
		Size      4 * %0xxxxxxx
		Flags         $xx xx
	*/

	for {
		header := make([]byte, 10)

		if _, err := io.ReadFull(r, header); err != nil {
			break
		}

		id := string(header[0:4])
		size := (int(header[4]) << 24) | (int(header[5]) << 16) | (int(header[6]) << 8) | int(header[7])
		data := make([]byte, size)

		if _, err := io.ReadFull(r, data); err != nil {
			break
		}

		switch id {
		case "TALB":
			t.Album = decode(data[1:], data[0])
		case "TPE1":
			t.Artist = decode(data[1:], data[0])
		case "TIT2":
			t.Title = decode(data[1:], data[0])
		case "APIC":
			t.Thumbnail = data
		default:
		}

		break
	}

	return nil
}

type Header struct {
	Tag      string
	Revision int
	Minor    int
	Flags    int
	Size     int
}

func (h Header) Version() string {
	return fmt.Sprintf("%sv2.%d.%d", h.Tag, h.Revision, h.Minor)
}

func (h Header) Flag(flag int) bool {
	return (h.Flags & flag) != 0
}

type Song struct {
	Path string
	Tag  *Tag
}

func (s *Song) Load() error {
	f, err := os.Open(s.Path)

	if err != nil {
		return err
	}

	defer f.Close()

	/* read and parse the header */

	bs := make([]byte, 10)

	if _, err = io.ReadFull(f, bs); err != nil {
		return err
	}

	s.Tag = &Tag{
		Header: Header{
			Tag:      string(bs[0:3]),
			Revision: int(bs[3]),
			Minor:    int(bs[4]),
			Flags:    int(bs[5]),
			Size:     (int(bs[6]) << 24) | (int(bs[7]) << 16) | (int(bs[8]) << 8) | int(bs[9]),
		},
	}

	/* if an extended header is present, skip it */

	if s.Tag.Header.Flag(ExtendedHeader) {
		bs := make([]byte, 6)

		if _, err = io.ReadFull(f, bs); err != nil {
			return err
		}

		totalsize := ((int(bs[0]) << 24) | (int(bs[1]) << 16) | (int(bs[2]) << 8) | int(bs[3])) - 6

		if totalsize > 0 {
			bs := make([]byte, totalsize)

			if _, err = io.ReadFull(f, bs); err != nil {
				return err
			}
		}
	}

	/* pull the tag frames out of the remaining portion of the file */

	framelen := s.Tag.Header.Size

	if s.Tag.Header.Flag(FooterPresent) {
		framelen -= 10
	}

	frames := make([]byte, framelen)

	if _, err = io.ReadFull(f, frames); err != nil {
		return err
	}

	return s.Tag.ParseFrames(bytes.NewReader(frames))
}

func decode(bs []byte, encoding byte) string {
	// 00 – ISO-8859-1 (ASCII).
	// 01 – UCS-2 (UTF-16 encoded Unicode with BOM), in ID3v2.2 and ID3v2.3.
	// 02 – UTF-16BE encoded Unicode without BOM, in ID3v2.4.
	// 03 – UTF-8 encoded Unicode, in ID3v2.4.

	switch encoding {
	case 0, 3:
		return string(bs)
	case 1, 2:
		u16s := make([]uint16, 1)
		ret := &bytes.Buffer{}
		b8buf := make([]byte, 4)

		for i := 0; i < len(bs); i += 2 {
			u16s[0] = uint16(bs[i]) + (uint16(bs[i+1]) << 8)
			r := utf16.Decode(u16s)
			n := utf8.EncodeRune(b8buf, r[0])
			ret.Write(b8buf[:n])
		}

		return ret.String()
	default:
		return ""
	}
}
