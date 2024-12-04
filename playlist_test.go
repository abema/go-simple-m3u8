package m3u8

import (
	"bytes"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestDecodePlaylist(t *testing.T) {
	for idx, testData := range []struct {
		input        string
		output       string
		playlistType PlaylistType
	}{
		{input: sampleMaster01Input, output: sampleMaster01Output, playlistType: PlaylistTypeMaster},
		{input: sampleMaster02Input, output: sampleMaster02Output, playlistType: PlaylistTypeMaster},
		{input: sampleMedia01Input, output: sampleMedia01Output, playlistType: PlaylistTypeMedia},
		{input: sampleMedia02Input, output: sampleMedia02Output, playlistType: PlaylistTypeMedia},
	} {
		t.Run(fmt.Sprintf("%d", idx), func(t *testing.T) {
			r := bytes.NewReader([]byte(testData.input))
			playlist, err := DecodePlaylist(r)
			require.NoError(t, err)
			require.Equal(t, playlist.Type(), testData.playlistType)
			require.Equal(t, playlist.Master() != nil, testData.playlistType == PlaylistTypeMaster)
			require.Equal(t, playlist.Media() != nil, testData.playlistType == PlaylistTypeMedia)
			w := bytes.NewBuffer(nil)
			require.NoError(t, playlist.Encode(w))
			assert.Equal(t, testData.output, w.String())
		})
	}
}
