package main

import (
	"fmt"
	"strings"

	m3u8 "github.com/abema/go-simple-m3u8"
)

const sampleData = `#EXTM3U
#EXT-X-STREAM-INF:BANDWIDTH=1280000,AVERAGE-BANDWIDTH=1000000
http://example.com/low.m3u8
#EXT-X-STREAM-INF:BANDWIDTH=2560000,AVERAGE-BANDWIDTH=2000000
http://example.com/mid.m3u8
#EXT-X-STREAM-INF:BANDWIDTH=7680000,AVERAGE-BANDWIDTH=6000000
http://example.com/hi.m3u8
#EXT-X-STREAM-INF:BANDWIDTH=65000,CODECS="mp4a.40.5"
http://example.com/audio-only.m3u8
`

func main() {
	playlist, err := m3u8.DecodePlaylist(strings.NewReader(sampleData))
	if err != nil {
		panic(err)
	}
	fmt.Println("Type:", playlist.Type())
	fmt.Println("Tags:", len(playlist.Master().Tags))
	for name, values := range playlist.Master().Tags {
		fmt.Printf("  %s: %d\n", name, len(values))
	}
	fmt.Println("Streams:")
	for i, stream := range playlist.Master().Streams {
		fmt.Printf("  %d:\n", i)
		height, width, _ := stream.Attributes.Resolution()
		fmt.Println("    Height:", height)
		fmt.Println("    Width:", width)
		fmt.Println("    URI:", stream.URI)
	}
}

/* Output:
Type: master
Tags: 1
  EXTM3U: 1
Streams:
  0:
    Height: 0
    Width: 0
    URI: http://example.com/low.m3u8
  1:
    Height: 0
    Width: 0
    URI: http://example.com/mid.m3u8
  2:
    Height: 0
    Width: 0
    URI: http://example.com/hi.m3u8
  3:
    Height: 0
    Width: 0
    URI: http://example.com/audio-only.m3u8
*/
