package m3u8

import (
	//"github.com/polypmer/m3u8"
	"strings"
	"testing"
)

func TestRead(t *testing.T) {
	in := `#EXTM3U

#EXTINF:123, Sample artist - Sample title
/home/user/Music/sample.mp3

#EXTINF:321,Example Artist - Example title
/home/user/Music/example.ogg
`
	r := NewReader(strings.NewReader(in))
	src, err := r.Read()
	if err != nil {
		t.Error(err)
	}
	if src[0] != "/home/user/Music/sample.mp3" {
		t.Error("DIdn't read the right src")
	}
}
