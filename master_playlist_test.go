package m3u8

import (
	"bytes"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const sampleMaster01Input = `#EXTM3U
`

const sampleMaster01Output = `#EXTM3U
`

const sampleMaster02Input = `#EXTM3U
#EXT-X-STREAM-INF:BANDWIDTH=1280000,AVERAGE-BANDWIDTH=1000000
http://example.com/low.m3u8

#EXT-X-STREAM-INF:BANDWIDTH=2560000,AVERAGE-BANDWIDTH=2000000
http://example.com/mid.m3u8

#EXT-X-STREAM-INF:BANDWIDTH=7680000,AVERAGE-BANDWIDTH=6000000
http://example.com/hi.m3u8

#EXT-X-STREAM-INF:BANDWIDTH=65000,CODECS="mp4a.40.5"
http://example.com/audio-only.m3u8
`

const sampleMaster02Output = `#EXTM3U
#EXT-X-STREAM-INF:AVERAGE-BANDWIDTH=1000000,BANDWIDTH=1280000
http://example.com/low.m3u8
#EXT-X-STREAM-INF:AVERAGE-BANDWIDTH=2000000,BANDWIDTH=2560000
http://example.com/mid.m3u8
#EXT-X-STREAM-INF:AVERAGE-BANDWIDTH=6000000,BANDWIDTH=7680000
http://example.com/hi.m3u8
#EXT-X-STREAM-INF:BANDWIDTH=65000,CODECS="mp4a.40.5"
http://example.com/audio-only.m3u8
`

const sampleAlternativeStreamInput = `#EXTM3U
#EXT-X-MEDIA:TYPE=AUDIO,GROUP-ID="aac",NAME="English",DEFAULT=YES,AUTOSELECT=YES,LANGUAGE="en",URI="main/english-audio.m3u8"
#EXT-X-MEDIA:TYPE=AUDIO,GROUP-ID="aac",NAME="Deutsch",DEFAULT=NO,AUTOSELECT=YES,LANGUAGE="de",URI="main/german-audio.m3u8"
#EXT-X-STREAM-INF:BANDWIDTH=1280000,CODECS="avc1.4d401e",AUDIO="aac"
low-video.m3u8
#EXT-X-STREAM-INF:BANDWIDTH=2560000,CODECS="avc1.4d401e",AUDIO="aac"
middle-video.m3u8
#EXT-X-STREAM-INF:BANDWIDTH=7680000,CODECS="avc1.4d401e",AUDIO="aac"
high-video.m3u8
`

const sampleAlternativeStreamOutput = `#EXTM3U
#EXT-X-MEDIA:GROUP-ID="aac",AUTOSELECT=YES,DEFAULT=YES,LANGUAGE="en",NAME="English",TYPE=AUDIO,URI="main/english-audio.m3u8"
#EXT-X-MEDIA:GROUP-ID="aac",AUTOSELECT=YES,DEFAULT=NO,LANGUAGE="de",NAME="Deutsch",TYPE=AUDIO,URI="main/german-audio.m3u8"
#EXT-X-STREAM-INF:AUDIO="aac",BANDWIDTH=1280000,CODECS="avc1.4d401e"
low-video.m3u8
#EXT-X-STREAM-INF:AUDIO="aac",BANDWIDTH=2560000,CODECS="avc1.4d401e"
middle-video.m3u8
#EXT-X-STREAM-INF:AUDIO="aac",BANDWIDTH=7680000,CODECS="avc1.4d401e"
high-video.m3u8
`

const sampleIFrameOnlyInput = `#EXTM3U
#EXT-X-STREAM-INF:BANDWIDTH=1280000
low-audio-video.m3u8
#EXT-X-I-FRAME-STREAM-INF:BANDWIDTH=86000,URI="low-iframe.m3u8"
#EXT-X-STREAM-INF:BANDWIDTH=2560000
middle-audio-video.m3u8
#EXT-X-I-FRAME-STREAM-INF:BANDWIDTH=150000,URI="middle-iframe.m3u8"
#EXT-X-STREAM-INF:BANDWIDTH=7680000
high-audio-video.m3u8
#EXT-X-I-FRAME-STREAM-INF:BANDWIDTH=550000,URI="high-iframe.m3u8"
#EXT-X-STREAM-INF:BANDWIDTH=65000,CODECS="mp4a.40.5"
audio-only.m3u8
`

const sampleIFrameOnlyOutput = `#EXTM3U
#EXT-X-STREAM-INF:BANDWIDTH=1280000
low-audio-video.m3u8
#EXT-X-STREAM-INF:BANDWIDTH=2560000
middle-audio-video.m3u8
#EXT-X-STREAM-INF:BANDWIDTH=7680000
high-audio-video.m3u8
#EXT-X-STREAM-INF:BANDWIDTH=65000,CODECS="mp4a.40.5"
audio-only.m3u8
#EXT-X-I-FRAME-STREAM-INF:BANDWIDTH=86000,URI="low-iframe.m3u8"
#EXT-X-I-FRAME-STREAM-INF:BANDWIDTH=150000,URI="middle-iframe.m3u8"
#EXT-X-I-FRAME-STREAM-INF:BANDWIDTH=550000,URI="high-iframe.m3u8"
`

func TestDecodeMasterPlaylist(t *testing.T) {
	for idx, testData := range []struct {
		input   string
		output  string
		streams int
	}{
		{input: sampleMaster01Input, output: sampleMaster01Output, streams: 0},
		{input: sampleMaster02Input, output: sampleMaster02Output, streams: 4},
		{input: sampleAlternativeStreamInput, output: sampleAlternativeStreamOutput, streams: 3},
		{input: sampleIFrameOnlyInput, output: sampleIFrameOnlyOutput, streams: 4},
	} {
		t.Run(fmt.Sprintf("%d", idx), func(t *testing.T) {
			r := bytes.NewReader([]byte(testData.input))
			playlist, err := DecodeMasterPlaylist(r)
			require.NoError(t, err)
			require.Len(t, playlist.Streams, testData.streams)
			w := bytes.NewBuffer(nil)
			require.NoError(t, playlist.Encode(w))
			assert.Equal(t, testData.output, w.String())
		})
	}
}
