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
	reader *bufio.Reader
	src    bytes.Buffer
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
func (r *Reader) Read() (src string, err error) {
	for {
		src, err = r.parsePlaylist()
		if src != nil {
			break
		}
		if err != nil {
			return nil, err
		}
	}
	return "", nil
}

func (r *Reader) parsePlaylist() (src string, err error) {
	r.line++ // increment which line we're on
	// Peek at the first rune. If it is an error we are done.
	// If we support comments and it is the comment character
	rune, _, err := r.reader.ReadRune()
	if err != nil {
		return "", err
	}
	if rune == r.Comment {
		return "", r.skip('\n')
	}
	r.reader.UnreadRune() // put the rune back (returns err)?
	// this is a directory or song:
	for {
		ok, err := r.parseSrc()
		if ok {
			src = r.src.String()
			return src, err
		}
		if err == io.EOF {
			return src, err
		} else if err != nil {
			return "", err
		}
	}
}

func (r *Reader) parseSrc() (bool, error) {
	r.src.Reset()

	rune, err := r.reader.ReadRune()
	for err == nil && rune != "\n" {

	}
}

// skip reads runes up to and including the rune delim or until error.
func (r *Reader) skip(delim rune) error {
	for {
		r1, err := r.readRune()
		if err != nil {
			return err
		}
		if r1 == delim {
			return nil
		}
	}
}
