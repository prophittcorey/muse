package audio

import (
	"bytes"
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"io"
	"log"
	"os"
	"strings"
	"unicode/utf16"
	"unicode/utf8"
)

const (
	Unsynchronization = 1 << 7
	ExtendedHeader    = 1 << 6
	Experimental      = 1 << 5
	FooterPresent     = 1 << 4

	syncsafebytelen = 7
	normbytelen     = 8
)

var (
	nullbyte   = []byte{0}
	doublenull = []byte{0, 0}
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

	log.Printf("%s: %d bytes, %b\n", t.Header.Version(), t.Header.Size, t.Header.Flags)

	for {
		header := make([]byte, 10)

		if _, err := io.ReadFull(r, header); err != nil {
			break
		}

		bytelen := syncsafebytelen

		if t.Header.Revision < 4 {
			bytelen = normbytelen
		}

		id := string(header[0:4])

		size := bytestoint(header[4:8], bytelen)

		data := make([]byte, size)

		if n, err := io.ReadFull(r, data); err != nil {
			log.Printf("unexpected EOF; read %d / %d; %s\n", n, size, err)
			return err
		}

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

				Text encoding      $xx
				MIME Type          <text string> $00
				Picture Type       $xx
				Description        <text string according to encoding> $00 (00)
				Picture Data       <binary data>
			*/

			encoding := data[0]

			mime, data, _ := bytes.Cut(data[1:], nullbyte)

			pictype := data[0] /* special byte */

			/* the description may need a double null according to the text encoding */
			sep := nullbyte

			if encoding == 1 || encoding == 2 {
				sep = doublenull
			}

			description, data, _ := bytes.Cut(data[1:], sep)

			t.Picture = &Picture{
				Description: decode(description, encoding),
				Mime:        decode(mime, encoding),
				Type:        string(pictype),
				Data:        data,
			}
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

type Track struct {
	ID   string
	Path string
	Tag  *Tag
}

func (t *Track) Load() error {
	f, err := os.Open(t.Path)

	if err != nil {
		return err
	}

	defer f.Close()

	hash := md5.Sum([]byte(t.Path))

	t.ID = hex.EncodeToString(hash[:])

	/* read and parse the header */

	bs := make([]byte, 10)

	if _, err = io.ReadFull(f, bs); err != nil {
		return err
	}

	t.Tag = &Tag{
		Header: Header{
			Tag:      string(bs[0:3]),
			Revision: int(bs[3]),
			Minor:    int(bs[4]),
			Flags:    int(bs[5]),
			Size:     bytestoint(bs[6:10], syncsafebytelen),
		},
	}

	if t.Tag.Header.Tag != "ID3" || t.Tag.Header.Revision < 2 {
		return fmt.Errorf("error: tag version not supported; %s", t.Tag.Header.Version())
	}

	/* if an extended header is present, skip it */

	if t.Tag.Header.Flag(ExtendedHeader) {
		bs := make([]byte, 6)

		if _, err = io.ReadFull(f, bs); err != nil {
			return err
		}

		totalsize := bytestoint(bs[0:4], syncsafebytelen)

		if totalsize > 0 {
			bs := make([]byte, totalsize)

			if _, err = io.ReadFull(f, bs); err != nil {
				return err
			}

			t.Tag.Header.Size -= totalsize
		}
	}

	/* pull the tag frames out of the remaining portion of the file */

	framelen := t.Tag.Header.Size

	frames := make([]byte, framelen)

	if _, err = io.ReadFull(f, frames); err != nil {
		return err
	}

	log.Printf("%s (%d bytes)\n", t.Path, len(frames))

	return t.Tag.ParseFrames(bytes.NewReader(frames))
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

		// HACK: Not sure why UTF16 strings in ID3v2.3 have a null terminator.
		return strings.Replace(ret.String(), "\x00", "", -1)
	default:
		return ""
	}
}

func bytestoint(bs []byte, bitlen int) int {
	var i int

	for _, b := range bs {
		i = (i << bitlen) | int(b)
	}

	return i
}
