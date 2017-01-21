package m3u8

import (
	"m3u8"
	"testing"
)

func ExampleReader() {
	in := `#EXTM3U

#EXTINF:123, Sample artist - Sample title
C:\Path\Sample.mp3

#EXTINF:321,Example Artist - Example title
/home/user/Music/example.ogg
`
	r := m3u8.NewReader(strings.NewReader(in))

	for {
		src, err := r.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println(src)
	}

	// Output:
	// [C:\Path\Sample.mp3 (123 * time.second) "Sample artist - Sample title"]
}
