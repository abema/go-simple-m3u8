package m3u8

import (
	"strconv"
	"strings"
	"time"
)

// SegmentTags represents the tags of a segment.
type SegmentTags Tags

// Raw returns the raw tags.
func (tags SegmentTags) Raw() Tags {
	return Tags(tags)
}

// Set sets the tag.
// If the tag already exists, it will be overwritten.
func (tags SegmentTags) Set(tag *Tag) {
	tags.Raw().Set(tag)
}

// ProgramDateTime returns the value of the EXT-X-PROGRAM-DATE-TIME tag.
func (tags SegmentTags) ProgramDateTime() (time.Time, bool) {
	values, ok := tags[TagExtXProgramDateTime]
	if !ok || len(values) == 0 {
		return time.Time{}, false
	}
	value := values[0]
	t, err := time.Parse(time.RFC3339Nano, value)
	if err != nil {
		return time.Time{}, false
	}
	return t, true
}

// SetProgramDateTime sets the value of the EXT-X-PROGRAM-DATE-TIME tag.
func (tags SegmentTags) SetProgramDateTime(t time.Time) {
	tags[TagExtXProgramDateTime] = []string{t.Format("2006-01-02T15:04:05.999Z07:00")}
}

// ExtInfValue returns the value of the EXTINF tag.
func (tags SegmentTags) ExtInfValue() float64 {
	values, ok := tags[TagExtInf]
	if !ok || len(values) == 0 {
		return 0
	}
	value := values[0]
	idx := strings.Index(value, ",")
	if idx <= 0 {
		return 0
	}
	duration, err := strconv.ParseFloat(value[:idx], 64)
	if err != nil {
		return 0
	}
	return duration
}

// SetExtInfValue sets the value of the EXTINF tag.
func (tags SegmentTags) SetExtInfValue(duration float64, bitSize int) {
	tags[TagExtInf] = []string{strconv.FormatFloat(duration, 'f', -1, bitSize) + ","}
}
