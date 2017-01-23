package m3u8

import (
	//"github.com/polypmer/m3u8"
	"io"
	"log"
	"strings"
	"testing"
)

func TestReader(t *testing.T) {
	//do nothing
}

func ExampleReader() {
	in := `#EXTM3U

#EXTINF:123, Sample artist - Sample title
C:\Path\Sample.mp3

#EXTINF:321,Example Artist - Example title
/home/user/Music/example.ogg
`
	r := NewReader(strings.NewReader(in))
	b := make([]byte, 0)

	for {
		_, err := r.Read(b)
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatal(err)
		}
		log.Println(string(b))
	}

	// Output:
	// [C:\Path\Sample.mp3 (123 * time.second) "Sample artist - Sample title"]
}
