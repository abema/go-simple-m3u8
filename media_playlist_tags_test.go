package m3u8

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestMediaPlaylistTags(t *testing.T) {
	t.Run("getters", func(t *testing.T) {
		tags := MediaPlaylistTags{
			"EXT-X-VERSION":                []string{"4"},
			"EXT-X-TARGETDURATION":         []string{"12"},
			"EXT-X-MEDIA-SEQUENCE":         []string{"12345"},
			"EXT-X-DISCONTINUITY-SEQUENCE": []string{"123"},
		}
		assert.Equal(t, int(4), tags.Version())
		assert.Equal(t, 12, tags.TargetDuration())
		assert.Equal(t, int64(12345), tags.MediaSequence())
		assert.Equal(t, int64(123), tags.DiscontinuitySequence())
	})

	t.Run("setters", func(t *testing.T) {
		tags := make(MediaPlaylistTags)
		tags.SetMediaSequence(12345)
		tags.SetDiscontinuitySequence(123)
		assert.Equal(t, MediaPlaylistTags{
			"EXT-X-MEDIA-SEQUENCE":         []string{"12345"},
			"EXT-X-DISCONTINUITY-SEQUENCE": []string{"123"},
		}, tags)
	})
}

func TestDateRangeAttrs(t *testing.T) {
	t.Run("Decode", func(t *testing.T) {
		t.Run("cue_out", func(t *testing.T) {
			attrs := DateRangeAttrs{
				"ID":               `"4"`,
				"START-DATE":       `"2023-05-12T05:09:20.988Z"`,
				"PLANNED-DURATION": "60.026",
				"SCTE35-OUT":       "0xFC306A",
			}
			values, err := attrs.Decode()
			require.NoError(t, err)
			assert.Equal(t, "4", values.EventID)
			assert.Equal(t, time.Date(2023, time.May, 12, 5, 9, 20, 988e6, time.UTC), values.StartDate)
			assert.True(t, values.EndDate.IsZero())
			assert.Zero(t, values.Duration)
			assert.Equal(t, 60.026, values.PlannedDuration)
			assert.Equal(t, []byte{0xFC, 0x30, 0x6A}, values.SCTE35Out)
		})

		t.Run("cue_in", func(t *testing.T) {
			attrs := DateRangeAttrs{
				"ID":         `"4"`,
				"START-DATE": `"2023-05-12T05:09:20.988Z"`,
				"END-DATE":   `"2023-05-12T05:10:21.015Z"`,
				"DURATION":   "60.026",
			}
			values, err := attrs.Decode()
			require.NoError(t, err)
			assert.Equal(t, "4", values.EventID)
			assert.Equal(t, time.Date(2023, time.May, 12, 5, 9, 20, 988e6, time.UTC), values.StartDate)
			assert.Equal(t, time.Date(2023, time.May, 12, 5, 10, 21, 015e6, time.UTC), values.EndDate)
			assert.Equal(t, 60.026, values.Duration)
			assert.Zero(t, values.PlannedDuration)
			assert.Nil(t, values.SCTE35Out)
		})

		t.Run("invalid_start_date", func(t *testing.T) {
			attrs := DateRangeAttrs{
				"ID":         `"4"`,
				"START-DATE": `"foo"`,
				"END-DATE":   `"2023-05-12T05:10:21.015Z"`,
				"DURATION":   "60.026",
			}
			_, err := attrs.Decode()
			require.Error(t, err)
		})

		t.Run("invalid_end_date", func(t *testing.T) {
			attrs := DateRangeAttrs{
				"ID":         `"4"`,
				"START-DATE": `"2023-05-12T05:09:20.988Z"`,
				"END-DATE":   `"bar"`,
				"DURATION":   "60.026",
			}
			_, err := attrs.Decode()
			require.Error(t, err)
		})

		t.Run("invalid_scte35_out", func(t *testing.T) {
			attrs := DateRangeAttrs{
				"ID":               `"4"`,
				"START-DATE":       `"2023-05-12T05:09:20.988Z"`,
				"PLANNED-DURATION": "60.026",
				"SCTE35-OUT":       "FC306A",
			}
			_, err := attrs.Decode()
			require.Error(t, err)
		})
	})
}
