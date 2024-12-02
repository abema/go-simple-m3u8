package m3u8

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestSegmentTags(t *testing.T) {
	t.Run("getters", func(t *testing.T) {
		tags := SegmentTags{
			"EXT-X-PROGRAM-DATE-TIME": []string{"2023-01-02T03:04:05.678Z"},
			"EXTINF":                  []string{"12.34,"},
		}
		pdt, ok := tags.ProgramDateTime()
		require.True(t, ok)
		assert.Equal(t, int64(1672628645678e+6), pdt.UnixNano())
		assert.Equal(t, 12.34, tags.ExtInfValue())
	})

	t.Run("setters", func(t *testing.T) {
		tags := make(SegmentTags)
		tags.SetProgramDateTime(time.Unix(1672628645, 678e+6).In(time.UTC))
		tags.SetExtInfValue(12.34, 64)
		tags.Set(&Tag{Name: "EXT-X-MY-TAG", Attributes: "test-value"})
		assert.Equal(t, SegmentTags{
			"EXT-X-PROGRAM-DATE-TIME": []string{"2023-01-02T03:04:05.678Z"},
			"EXTINF":                  []string{"12.34,"},
			"EXT-X-MY-TAG":            []string{"test-value"},
		}, tags)
	})

	t.Run("pdt_not_found", func(t *testing.T) {
		tags := SegmentTags{
			"EXTINF": []string{"12.34,"},
		}
		_, ok := tags.ProgramDateTime()
		assert.False(t, ok)
	})

	t.Run("invalid_pdt", func(t *testing.T) {
		tags := SegmentTags{
			"EXT-X-PROGRAM-DATE-TIME": []string{"invalid"},
		}
		_, ok := tags.ProgramDateTime()
		assert.False(t, ok)
	})
}
