package m3u8

import (
	"bytes"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var sampleMedia01Input = `#EXTM3U
#EXT-X-VERSION:3
#EXT-X-MEDIA-SEQUENCE:7794
#EXT-X-TARGETDURATION:15

#EXT-X-KEY:METHOD=AES-128,URI="https://drm.example.com/key.php?r=52"

#EXTINF:2.833,
http://media.example.com/sequence1-A.ts
#EXTINF:15.0,
http://media.example.com/sequence1-B.ts
#EXTINF:13.333,
http://media.example.com/sequence1-C.ts

#EXT-X-KEY:METHOD=AES-128,URI="https://drm.example.com/key.php?r=53"

#EXTINF:15.0,
http://media.example.com/sequence2-A.ts
`

var sampleMedia01Output = `#EXTM3U
#EXT-X-VERSION:3
#EXT-X-TARGETDURATION:15
#EXT-X-MEDIA-SEQUENCE:7794
#EXT-X-KEY:METHOD=AES-128,URI="https://drm.example.com/key.php?r=52"
#EXTINF:2.833,
http://media.example.com/sequence1-A.ts
#EXTINF:15.0,
http://media.example.com/sequence1-B.ts
#EXTINF:13.333,
http://media.example.com/sequence1-C.ts
#EXT-X-KEY:METHOD=AES-128,URI="https://drm.example.com/key.php?r=53"
#EXTINF:15.0,
http://media.example.com/sequence2-A.ts
`

var sampleMedia02Input = `#EXTM3U
#EXT-X-VERSION:3
#EXT-X-TARGETDURATION:15
#EXT-X-KEY:METHOD=NONE
#EXTINF:2.833,
http://media.example.com/sequence1-A.ts
#EXTINF:15.0,
http://media.example.com/sequence1-B.ts
#EXTINF:13.333,
http://media.example.com/sequence1-C.ts
#EXTINF:15.0,
http://media.example.com/sequence2-A.ts
#EXT-X-ENDLIST
`

var sampleMedia02Output = `#EXTM3U
#EXT-X-VERSION:3
#EXT-X-TARGETDURATION:15
#EXT-X-KEY:METHOD=NONE
#EXTINF:2.833,
http://media.example.com/sequence1-A.ts
#EXTINF:15.0,
http://media.example.com/sequence1-B.ts
#EXTINF:13.333,
http://media.example.com/sequence1-C.ts
#EXTINF:15.0,
http://media.example.com/sequence2-A.ts
#EXT-X-ENDLIST
`

var sampleMedia03Input = `#EXTM3U
#EXT-X-VERSION:3
#EXT-X-TARGETDURATION:10
#EXT-X-KEY:METHOD=NONE
#EXTINF:2.833,
http://media.example.com/sequence1-A.ts
#EXTINF:10.0,
http://media.example.com/sequence1-B.ts
#EXTINF:8.333,
http://media.example.com/sequence1-C.ts
#EXT-OATCLS-SCTE35:E4eq3p3EuL2CDpJgMLapii+Uu/phHQwUL6W1JUGnttg=
#EXT-X-BLACKOUT:TYPE=NETWORK_END
#EXT-X-DISCONTINUITY
#EXTINF:10.0,
http://media.example.com/sequence2-A.ts
#EXT-X-ENDLIST
`

var sampleMedia03Output = `#EXTM3U
#EXT-X-VERSION:3
#EXT-X-TARGETDURATION:10
#EXT-X-KEY:METHOD=NONE
#EXTINF:2.833,
http://media.example.com/sequence1-A.ts
#EXTINF:10.0,
http://media.example.com/sequence1-B.ts
#EXTINF:8.333,
http://media.example.com/sequence1-C.ts
#EXT-OATCLS-SCTE35:E4eq3p3EuL2CDpJgMLapii+Uu/phHQwUL6W1JUGnttg=
#EXT-X-BLACKOUT:TYPE=NETWORK_END
#EXT-X-DISCONTINUITY
#EXTINF:10.0,
http://media.example.com/sequence2-A.ts
#EXT-X-ENDLIST
`

var sampleCue01 = `#EXTM3U
#EXT-X-VERSION:3
#EXT-X-TARGETDURATION:5
#EXT-X-MEDIA-SEQUENCE:2680
#EXT-X-CUE-OUT-CONT:Duration=30,ElapsedTime=20
#EXTINF:5,
http://media.example.com/segment2680.ts
#EXT-X-CUE-OUT-CONT:Duration=30,ElapsedTime=25
#EXTINF:5,
http://media.example.com/segment2681.ts
#EXT-X-CUE-IN
#EXT-X-CUE-OUT:60
#EXTINF:5,
http://media.example.com/segment2682.ts
#EXT-X-CUE-OUT-CONT:Duration=30,ElapsedTime=5
#EXTINF:5,
http://media.example.com/segment2683.ts
`

var sampleDateRange01 = `#EXTM3U
#EXT-X-VERSION:3
#EXT-X-TARGETDURATION:5
#EXT-X-MEDIA-SEQUENCE:2680
#EXT-X-PROGRAM-DATE-TIME:2024-01-01T01:00:50.000Z
#EXT-X-DATERANGE:ID="100",START-DATE="2024-01-01T01:00:00.000Z",END-DATE="2024-01-01T01:01:00.000Z",DURATION=60.000
#EXTINF:5,
http://media.example.com/segment2680.ts
#EXTINF:5,
http://media.example.com/segment2681.ts
#EXT-X-DATERANGE:ID="100",START-DATE="2024-01-01T01:00:00.000Z",END-DATE="2024-01-01T01:01:00.000Z",DURATION=60.000
#EXT-X-DATERANGE:ID="200",START-DATE="2024-01-01T01:01:00.000Z",END-DATE="2024-01-01T01:02:30.000Z",DURATION=90.000
#EXTINF:5,
http://media.example.com/segment2682.ts
#EXTINF:5,
http://media.example.com/segment2683.ts
`

func TestDecodeMediaPlaylist(t *testing.T) {
	t.Run("simple", func(t *testing.T) {
		for idx, testData := range []struct {
			input          string
			output         string
			targetDuration int
		}{
			{input: sampleMedia01Input, output: sampleMedia01Output, targetDuration: 15},
			{input: sampleMedia02Input, output: sampleMedia02Output, targetDuration: 15},
			{input: sampleMedia03Input, output: sampleMedia03Output, targetDuration: 10},
		} {
			t.Run(fmt.Sprintf("%d", idx), func(t *testing.T) {
				r := bytes.NewReader([]byte(testData.input))
				playlist, err := DecodeMediaPlaylist(r)
				require.NoError(t, err)
				assert.Equal(t, testData.targetDuration, playlist.Tags.TargetDuration())
				require.Len(t, playlist.Segments, 4)
				assert.Equal(t, 2.833, playlist.Segments[0].Tags.ExtInfValue())
				w := bytes.NewBuffer(nil)
				require.NoError(t, playlist.Encode(w))
				assert.Equal(t, testData.output, w.String())
			})
		}
	})

	t.Run("cue", func(t *testing.T) {
		for idx, testData := range []struct {
			input  string
			output string
		}{
			{input: sampleCue01, output: sampleCue01},
			{input: sampleDateRange01, output: sampleDateRange01},
		} {
			t.Run(fmt.Sprintf("%d", idx), func(t *testing.T) {
				r := bytes.NewReader([]byte(testData.input))
				playlist, err := DecodeMediaPlaylist(r)
				require.NoError(t, err)
				require.NotEmpty(t, playlist.Segments)
				assert.Equal(t, 5.0, playlist.Segments[0].Tags.ExtInfValue())
				w := bytes.NewBuffer(nil)
				require.NoError(t, playlist.Encode(w))
				assert.Equal(t, testData.output, w.String())
			})
		}
	})
}
