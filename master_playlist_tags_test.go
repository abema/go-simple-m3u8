package m3u8

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestStreamInfAttrs(t *testing.T) {
	t.Run("Resolution", func(t *testing.T) {
		t.Run("ok", func(t *testing.T) {
			streamInf := StreamInfAttrs{"RESOLUTION": "1920x1080"}
			w, h, err := streamInf.Resolution()
			require.NoError(t, err)
			assert.Equal(t, 1920, w)
			assert.Equal(t, 1080, h)
		})

		t.Run("invalid_separator", func(t *testing.T) {
			streamInf := StreamInfAttrs{"RESOLUTION": "1920_1080"}
			_, _, err := streamInf.Resolution()
			require.Error(t, err)
		})

		t.Run("invalid_width", func(t *testing.T) {
			streamInf := StreamInfAttrs{"RESOLUTION": "abcx1080"}
			_, _, err := streamInf.Resolution()
			require.Error(t, err)
		})

		t.Run("invalid_height", func(t *testing.T) {
			streamInf := StreamInfAttrs{"RESOLUTION": "1920x###"}
			_, _, err := streamInf.Resolution()
			require.Error(t, err)
		})
	})
}
