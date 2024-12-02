package m3u8

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"math"
	"regexp"
	"sort"
	"strings"
)

const (
	// Basic Tags
	TagExtM3U      = "EXTM3U"
	TagExtXVersion = "EXT-X-VERSION"

	// Master Playlist Tags
	TagExtXMedia           = "EXT-X-MEDIA"
	TagExtXStreamInf       = "EXT-X-STREAM-INF"
	TagExtXIFrameStreamInf = "EXT-X-I-FRAME-STREAM-INF"
	TagExtXSessionData     = "EXT-X-SESSION-DATA"
	TagExtXSessionKey      = "EXT-X-SESSION-KEY"

	// Media Playlist Tags
	TagExtXTargetDuration        = "EXT-X-TARGETDURATION"
	TagExtXMediaSequence         = "EXT-X-MEDIA-SEQUENCE"
	TagExtXDiscontinuitySequence = "EXT-X-DISCONTINUITY-SEQUENCE"
	TagExtXEndlist               = "EXT-X-ENDLIST"
	TagExtXPlaylistType          = "EXT-X-PLAYLIST-TYPE"
	TagExtXIFramesOnly           = "EXT-X-I-FRAMES-ONLY"

	// Media or Master Playlist Tags
	TagExtXIndependentSegments = "EXT-X-INDEPENDENT-SEGMENTS"
	TagExtXStart               = "EXT-X-START"

	// Segment Tags
	TagExtInf              = "EXTINF"
	TagExtXByteRange       = "EXT-X-BYTERANGE"
	TagExtXDiscontinuity   = "EXT-X-DISCONTINUITY"
	TagExtXKey             = "EXT-X-KEY"
	TagExtXMap             = "EXT-X-MAP"
	TagExtXProgramDateTime = "EXT-X-PROGRAM-DATE-TIME"
	TagExtXDateRange       = "EXT-X-DATERANGE"

	// Cue
	TagExtXSCTE35     = "EXT-OATCLS-SCTE35"
	TagExtXAsset      = "EXT-X-ASSET"
	TagExtXCueOut     = "EXT-X-CUE-OUT"
	TagExtXCueOutCont = "EXT-X-CUE-OUT-CONT"
	TagExtXCueIn      = "EXT-X-CUE-IN"
	TagExtXBlackout   = "EXT-X-BLACKOUT"
)

var tagOrderMap = map[string]int{
	// Basic Tags
	TagExtM3U:      0,
	TagExtXVersion: 1,

	// Media Playlist Tags
	TagExtXTargetDuration:        100,
	TagExtXPlaylistType:          101,
	TagExtXIFramesOnly:           102,
	TagExtXMediaSequence:         103,
	TagExtXDiscontinuitySequence: 104,
	TagExtXEndlist:               math.MaxInt32,

	// Media or Master Playlist Tags
	TagExtXIndependentSegments: 200,
	TagExtXStart:               201,

	// Cue
	TagExtXCueIn:      300,
	TagExtXSCTE35:     301,
	TagExtXAsset:      302,
	TagExtXCueOut:     303,
	TagExtXCueOutCont: 304,
	TagExtXBlackout:   305,

	// Segment Tags
	TagExtXDiscontinuity:   400,
	TagExtXKey:             401,
	TagExtXMap:             402,
	TagExtXProgramDateTime: 403,
	TagExtXDateRange:       404,
	TagExtInf:              405,
	TagExtXByteRange:       406,

	// Master Playlist Tags
	TagExtXMedia:           500,
	TagExtXStreamInf:       501,
	TagExtXIFrameStreamInf: 502,
	TagExtXSessionData:     503,
	TagExtXSessionKey:      504,
}

func getTagOrder(name string) int {
	if order, ok := tagOrderMap[name]; ok {
		return order
	}
	return math.MaxInt
}

var segmentTagSet = map[string]struct{}{
	TagExtInf:              {},
	TagExtXByteRange:       {},
	TagExtXDiscontinuity:   {},
	TagExtXKey:             {},
	TagExtXMap:             {},
	TagExtXProgramDateTime: {},
	TagExtXDateRange:       {},
	TagExtXSCTE35:          {},
	TagExtXAsset:           {},
	TagExtXCueOut:          {},
	TagExtXCueOutCont:      {},
	TagExtXCueIn:           {},
	TagExtXBlackout:        {},
}

// TagName extracts the tag name from the line.
func TagName(line string) string {
	if len(line) == 0 || line[0] != '#' {
		return ""
	}
	name := line[1:]
	idx := strings.Index(name, ":")
	if idx != -1 {
		name = name[:idx]
	}
	return name
}

// AttributeString extracts the tag attributes from the line.
func AttributeString(line string) string {
	if len(line) == 0 || line[0] != '#' {
		return ""
	}
	idx := strings.Index(line, ":")
	if idx == -1 {
		return ""
	}
	return line[idx+1:]
}

// IsSegmentTagName returns true if the name is a segment tag name.
func IsSegmentTagName(name string) bool {
	_, t := segmentTagSet[name]
	return t
}

// Attributes represents a set of attributes of a tag.
type Attributes map[string]string

// String encodes the attributes to a string.
func (attr Attributes) String() string {
	var buf bytes.Buffer
	keys := make([]string, 0, len(attr))
	for key := range attr {
		keys = append(keys, key)
	}
	sort.Strings(keys)
	for _, key := range keys {
		value := attr[key]
		if buf.Len() != 0 {
			buf.WriteString(",")
		}
		buf.WriteString(key)
		if value != "" {
			buf.WriteString("=")
			buf.WriteString(value)
		}
	}
	return buf.String()
}

var regexpFirstAttribute = regexp.MustCompile(`^([0-9A-Za-z-]+)(?:=("[^"]*"|[^",]*)|)(?:,|$)`)

// ParseTagAttributes parses the attributes and returns it as Attributes.
func ParseTagAttributes(attributes string) (Attributes, error) {
	m := make(Attributes)
	for len(attributes) != 0 {
		s := regexpFirstAttribute.FindStringSubmatch(attributes)
		if len(s) != 3 {
			return nil, errors.New("invalid HLS tag attributes")
		}
		m[s[1]] = s[2]
		attributes = attributes[len(s[0]):]
	}
	return m, nil
}

// Tag represents a tag.
type Tag struct {
	Name       string
	Attributes string
}

// Encode encodes the tag to io.Writer.
func (tag *Tag) Encode(w io.Writer) error {
	if _, err := fmt.Fprintf(w, "#%s", tag.Name); err != nil {
		return err
	}
	if tag.Attributes != "" {
		if _, err := fmt.Fprintf(w, ":%s", tag.Attributes); err != nil {
			return err
		}
	}
	if _, err := w.Write([]byte("\n")); err != nil {
		return err
	}
	return nil
}

// Tags represents a set of tags.
type Tags map[string][]string

// First returns the first tag of the name.
func (tags Tags) First(name string) *Tag {
	if attrsList, ok := tags[name]; ok && len(attrsList) > 0 {
		return &Tag{
			Name:       name,
			Attributes: attrsList[0],
		}
	}
	return nil
}

// Last returns the last tag of the name.
func (tags Tags) Last(name string) *Tag {
	if attrsList, ok := tags[name]; ok && len(attrsList) > 0 {
		return &Tag{
			Name:       name,
			Attributes: attrsList[len(attrsList)-1],
		}
	}
	return nil
}

// Set sets the tag.
// If the tag already exists, it will be overwritten.
func (tags Tags) Set(tag *Tag) {
	tags[tag.Name] = []string{tag.Attributes}
}

// Add adds the tag.
// If the tag already exists, it will be appended.
func (tags Tags) Add(tag *Tag) {
	tags[tag.Name] = append(tags[tag.Name], tag.Attributes)
}

// Remove removes the tag.
// If the tag does not exist, it will do nothing.
// If the tag exists multiple times, all of them will be removed.
func (tags Tags) RemoveByName(name string) {
	delete(tags, name)
}

// List returns the sorted list of tags.
func (tags Tags) List() []*Tag {
	list := make([]*Tag, 0, len(tags))
	for name, attrsList := range tags {
		for _, attrs := range attrsList {
			list = append(list, &Tag{
				Name:       name,
				Attributes: attrs,
			})
		}
	}
	sort.SliceStable(list, func(i, j int) bool {
		return getTagOrder(list[i].Name) < getTagOrder(list[j].Name)
	})
	return list
}
