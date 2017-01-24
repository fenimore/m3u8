// The source:
//
//#EXTM3U
// #EXTINF:123, Sample artist - Sample title
// /home/user/Music/sample.mp3
package m3u8

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"io"
	_ "log"
)

// A ParseError is returned for parsing errors.
// The first line is 1.  The first column is 0.
type ParseError struct {
	Line int   // Line where the error occurred
	Err  error // The actual error
}

func (e *ParseError) Error() string {
	return fmt.Sprintf("line %d: %s", e.Line, e.Err)
}

// These are the errors that can be returned in ParseError.Error
var (
	ErrNoHeader = errors.New("no #EXTM3U header")
	ErrNoSource = errors.New("no source found")
)

// A Reader reads sources from a m3u8-encoded file.
//
// As returned by NewReader, a Reader expects input conforming to RFC SIKE.
// The exported fields can be changed to customize the details before the
// first call to Read or ReadAll.
//
//
type Reader struct {
	// Newline is the field delimiter.
	// It is set to \n ('\n'), by NewReader.
	Newline rune
	// Comment, if not 0, is the comment character. Lines beginning with the
	// Comment character without preceding whitespace are ignored.
	// With leading whitespace the Comment character becomes part of the
	// field, even if TrimLeadingSpace is true.
	Comment rune
	// Setting bool
	// TrimLeadingSpace bool

	line   int
	abs    string // absolute path
	reader *bufio.Reader
	src    bytes.Buffer // the source of Read
}

// NewReader returns a new Reader that reads from r.
func NewReader(r io.Reader) *Reader {
	return &Reader{
		Newline: '\n',
		Comment: '#',
		reader:  bufio.NewReader(r),
	}
}

// error creates a new ParseError based on err.

func (r *Reader) error(err error) error {
	return &ParseError{
		Line: r.line,
		Err:  err,
	}
}

// Read reads one line from r. The line is a string.
// string representing one field.
func (r *Reader) Read() (src []string, err error) {
	for {
		src, err = r.parsePlaylist()
		if src != nil {
			break
		}
		if err != nil {
			return nil, err
		}
	}
	return src, nil
}

// parseSrc reads and parses a single m3u8 src
func (r *Reader) parsePlaylist() (src []string, err error) {
	r.line++

	rune, err := r.readRune()
	if err != nil {
		return nil, err
	}

	if rune == r.Comment {
		r.line++
		return nil, r.skip('\n')
	}
	if rune == r.Newline {
		return nil, r.skip('\n')
	}

	r.reader.UnreadRune()
	// At this point we have at least another Playlist.
	for {
		have, delim, err := r.parseSrc()
		if have {
			src = append(src, r.src.String())
		}
		// FIXME, except where the comment was?
		if delim == '\n' || err == io.EOF {
			r.line++
			return src, err
		} else if err != nil {
			return nil, err
		}
	}
}

func (r *Reader) parseMeta() (have bool, delim rune, err error) {
	// switch on possible Meta information
	return have, delim, err
}

// TODO: get meta data
func (r *Reader) parseSrc() (have bool, delim rune, err error) {
	r.src.Reset()

	rune, err := r.readRune()

	if err == io.EOF {
		return true, 0, err
	}
	if err != nil {
		return false, 0, err
	}

	switch rune {
	case r.Newline:
		// Check below
	case '\n':
		return false, rune, nil
	//case r.Comment:

	default:
		for {
			r.src.WriteRune(rune)
			rune, err = r.readRune()
			if err != nil || rune == r.Newline {
				break
			}
		}
	}

	if err == io.EOF {
		return true, 0, err
	}
	if err != nil {
		return false, 0, err
	}
	return true, rune, nil
}

func (r *Reader) readRune() (rune, error) {
	rune, _, err := r.reader.ReadRune()
	// should I handle strange things?
	return rune, err
}

// skip reads runes up to and including the rune delim or until error.
func (r *Reader) skip(delim rune) error {
	for {
		rune, _, err := r.reader.ReadRune()
		if err != nil {
			return err
		}
		if rune == delim {
			return nil
		}
	}
}
