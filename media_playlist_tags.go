package m3u8

import (
	"encoding/hex"
	"errors"
	"strconv"
	"strings"
	"time"
)

// MediaPlaylistTags represents the tags of a media playlist.
type MediaPlaylistTags Tags

// Raw returns the raw tags.
func (tags MediaPlaylistTags) Raw() Tags {
	return Tags(tags)
}

// Set sets the tag.
// If the tag already exists, it will be overwritten.
func (tags MediaPlaylistTags) Set(tag *Tag) {
	tags.Raw().Set(tag)
}

// Remove removes the tag.
// If the tag does not exist, it will do nothing.
// If the tag exists multiple times, all of them will be removed.
func (tags MediaPlaylistTags) Remove(name string) {
	tags.Raw().Remove(name)
}

// Version returns the value of the EXT-X-VERSION tag.
func (tags MediaPlaylistTags) Version() int {
	values, ok := tags[TagExtXVersion]
	if !ok || len(values) == 0 {
		return 1
	}
	value := values[0]
	version, err := strconv.Atoi(value)
	if err != nil {
		return 1
	}
	return version
}

// TargetDuration returns the value of the EXT-X-TARGETDURATION tag.
func (tags MediaPlaylistTags) TargetDuration() int {
	values, ok := tags[TagExtXTargetDuration]
	if !ok || len(values) == 0 {
		return 0
	}
	value := values[0]
	duration, err := strconv.Atoi(value)
	if err != nil {
		return 0
	}
	return duration
}

// MediaSequence returns the value of the EXT-X-MEDIA-SEQUENCE tag.
func (tags MediaPlaylistTags) MediaSequence() int64 {
	values, ok := tags[TagExtXMediaSequence]
	if !ok || len(values) == 0 {
		return 0
	}
	value := values[0]
	sequence, err := strconv.ParseInt(value, 10, 64)
	if err != nil {
		return 0
	}
	return sequence
}

// SetMediaSequence sets the value of the EXT-X-MEDIA-SEQUENCE tag.
func (tags MediaPlaylistTags) SetMediaSequence(sequence int64) {
	tags[TagExtXMediaSequence] = []string{strconv.FormatInt(sequence, 10)}
}

// TargetDuration returns the value of the EXT-X-TARGETDURATION tag.
func (tags MediaPlaylistTags) DiscontinuitySequence() int64 {
	values, ok := tags[TagExtXDiscontinuitySequence]
	if !ok || len(values) == 0 {
		return 0
	}
	value := values[0]
	sequence, err := strconv.ParseInt(value, 10, 64)
	if err != nil {
		return 0
	}
	return sequence
}

// SetDiscontinuitySequence sets the value of the EXT-X-DISCONTINUITY-SEQUENCE tag.
func (tags MediaPlaylistTags) SetDiscontinuitySequence(sequence int64) {
	tags[TagExtXDiscontinuitySequence] = []string{strconv.FormatInt(sequence, 10)}
}

// DateRangeAttrs represents the attributes of the EXT-X-DATERANGE tag.
type DateRangeAttrs Attributes

// EventID returns the value of the ID attribute.
func (attrs DateRangeAttrs) EventID() string {
	return strings.Trim(attrs["ID"], `"`)
}

// StartDate returns the value of the START-DATE attribute.
func (attrs DateRangeAttrs) StartDate() (time.Time, error) {
	value := strings.Trim(attrs["START-DATE"], `"`)
	if value == "" {
		return time.Time{}, nil
	}
	return time.Parse(time.RFC3339Nano, value)
}

// EndDate returns the value of the END-DATE attribute.
func (attrs DateRangeAttrs) EndDate() (time.Time, error) {
	value := strings.Trim(attrs["END-DATE"], `"`)
	if value == "" {
		return time.Time{}, nil
	}
	return time.Parse(time.RFC3339Nano, value)
}

// Duration returns the value of the DURATION attribute.
func (attrs DateRangeAttrs) Duration() (float64, error) {
	value := attrs["DURATION"]
	if value == "" {
		return 0, nil
	}
	return strconv.ParseFloat(value, 64)
}

// PlannedDuration returns the value of the PLANNED-DURATION attribute.
func (attrs DateRangeAttrs) PlannedDuration() (float64, error) {
	value := attrs["PLANNED-DURATION"]
	if value == "" {
		return 0, nil
	}
	return strconv.ParseFloat(value, 64)
}

// SCTE35Out returns the value of the SCTE35-OUT attribute.
func (attrs DateRangeAttrs) SCTE35Out() ([]byte, error) {
	value := attrs["SCTE35-OUT"]
	if value == "" {
		return nil, nil
	}
	if !strings.HasPrefix(value, "0x") {
		return nil, errors.New("unknown prefix")
	}
	return hex.DecodeString(value[2:])
}

// DateRangeAttrValues represents the attribute values of the EXT-X-DATERANGE tag.
type DateRangeAttrValues struct {
	EventID         string
	StartDate       time.Time
	EndDate         time.Time
	Duration        float64
	PlannedDuration float64
	SCTE35Out       []byte
}

// Decode decodes all the attributes.
func (attrs DateRangeAttrs) Decode() (*DateRangeAttrValues, error) {
	var err error
	values := new(DateRangeAttrValues)
	values.EventID = attrs.EventID()
	values.StartDate, err = attrs.StartDate()
	if err != nil {
		return nil, err
	}
	values.EndDate, err = attrs.EndDate()
	if err != nil {
		return nil, err
	}
	values.Duration, err = attrs.Duration()
	if err != nil {
		return nil, err
	}
	values.PlannedDuration, err = attrs.PlannedDuration()
	if err != nil {
		return nil, err
	}
	values.SCTE35Out, err = attrs.SCTE35Out()
	if err != nil {
		return nil, err
	}
	return values, nil
}
