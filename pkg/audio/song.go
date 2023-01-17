package audio

import (
	"bytes"
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"io"
	"log"
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

type Picture struct {
	Mime        string
	Type        string
	Description string
	Data        []byte
}

type Tag struct {
	Header  Header
	Title   string
	Artist  string
	Album   string
	Date    string
	Picture *Picture
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

	log.Println("HEADER ", t.Header.Size)

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

		log.Println("GOT: ", id, size)
		// log.Println(decode(data[1:], data[0]))

		switch id {
		case "TALB":
			t.Album = decode(data[1:], data[0])
		case "TDRC":
			t.Date = decode(data[1:], data[0])
		case "TPE1", "TOPE":
			t.Artist = decode(data[1:], data[0])
		case "TIT2":
			t.Title = decode(data[1:], data[0])
		case "APIC":
			/*
				41 50 49 43  //APIC
				00 08 5A 04  //Frame Size
				00 03        //Flags: Unsynchronisation | Data Length Indicator.
				00 02 19 F5  //4 bytes data length

				00           //1 byte text encoding (ISO-8859-1)
				69 6D 61 67  //image/jpeg
				65 2F 6A 70  // ”
				65 67        // ”
			*/

			encoding := data[0]

			mime, data, _ := bytes.Cut(data[1:], []byte{0})

			pictype := data[0] /* special byte */

			description, data, _ := bytes.Cut(data[1:], []byte{0})

			log.Println(decode(mime, encoding), pictype, decode(description, encoding))

			t.Picture = &Picture{
				Description: decode(description, encoding),
				Mime:        decode(mime, encoding),
				Type:        string(pictype),
				Data:        data,
			}

			// For debugging purposes..
			// f, _ := os.Create(fmt.Sprintf("/tmp/%s.png", "album"))
			// f.Write(data)
			// f.Close()
		default:
		}
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
	ID   string
	Path string
	Tag  *Tag
}

func (s *Song) Load() error {
	f, err := os.Open(s.Path)

	if err != nil {
		return err
	}

	defer f.Close()

	hash := md5.Sum([]byte(s.Path))

	s.ID = hex.EncodeToString(hash[:])

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

	log.Printf("==\n%s\n", s.Path)

	if err := s.Tag.ParseFrames(bytes.NewReader(frames)); err != nil {
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
