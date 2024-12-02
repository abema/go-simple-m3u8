package m3u8

import (
	"bufio"
	"errors"
	"fmt"
	"io"
)

// ErrUnexpectedSegmentTags is returned when segment tags are found without a segment URI.
var ErrUnexpectedSegmentTags = errors.New("unexpected segment tags")

// MediaPlaylist represents a media playlist.
type MediaPlaylist struct {
	// Tags is a list of tags in the media playlist.
	// This list does not include segment tags.
	Tags MediaPlaylistTags

	// Segments is a list of segments in the media playlist.
	Segments []*Segment

	// EndList indicates that no more media segments will be added to the
	// media playlist file in the future.
	EndList bool
}

// Segment represents a media segment with its tags.
type Segment struct {
	// Tags is a list of tags in the segment.
	Tags SegmentTags

	// URI is the URI of the segment.
	URI string

	// Sequence is the media sequence number of the segment.
	// This field is set by DecodeMediaPlaylist.
	// When encoding a media playlist, this field is ignored.
	Sequence int64

	// DiscontinuitySequence is the discontinuity sequence number of the segment.
	// This field is set by DecodeMediaPlaylist.
	// When encoding a media playlist, this field is ignored.
	DiscontinuitySequence int64
}

// DecodeMediaPlaylist decodes a media playlist from io.Reader.
func DecodeMediaPlaylist(r io.Reader) (*MediaPlaylist, error) {
	scanner := bufio.NewScanner(r)
	var playlist MediaPlaylist
	playlist.Tags = make(MediaPlaylistTags)
	playlist.Segments = make([]*Segment, 0, 8)
	segmentTags := make(SegmentTags)
	for scanner.Scan() {
		line := scanner.Text()
		if line == "" {
			continue
		}
		tagName := TagName(line)
		if tagName == "" {
			playlist.Segments = append(playlist.Segments, &Segment{
				Tags: segmentTags,
				URI:  line,
			})
			segmentTags = make(SegmentTags)
		} else if IsSegmentTagName(tagName) {
			segmentTags.Raw().Add(&Tag{
				Name:       tagName,
				Attributes: AttributeString(line),
			})
		} else if tagName == TagExtXEndlist {
			playlist.EndList = true
		} else {
			playlist.Tags.Raw().Add(&Tag{
				Name:       tagName,
				Attributes: AttributeString(line),
			})
		}
	}
	if err := scanner.Err(); err != nil {
		return nil, err
	}
	sequence := playlist.Tags.MediaSequence()
	discSequence := playlist.Tags.DiscontinuitySequence()
	for _, segment := range playlist.Segments {
		if _, exists := segment.Tags[TagExtXDiscontinuity]; exists {
			discSequence++
		}
		segment.Sequence = sequence
		segment.DiscontinuitySequence = discSequence
		sequence++
	}
	if len(segmentTags) != 0 {
		return &playlist, ErrUnexpectedSegmentTags
	}
	return &playlist, nil
}

// Encode encodes a media playlist to io.Writer.
func (playlist *MediaPlaylist) Encode(w io.Writer) error {
	for _, tag := range playlist.Tags.Raw().List() {
		err := tag.Encode(w)
		if err != nil {
			return err
		}
	}
	for _, segment := range playlist.Segments {
		for _, tag := range segment.Tags.Raw().List() {
			err := tag.Encode(w)
			if err != nil {
				return err
			}
		}
		if _, err := fmt.Fprintf(w, "%s\n", segment.URI); err != nil {
			return err
		}
	}
	if playlist.EndList {
		if _, err := w.Write([]byte("#" + TagExtXEndlist + "\n")); err != nil {
			return err
		}
	}
	return nil
}
