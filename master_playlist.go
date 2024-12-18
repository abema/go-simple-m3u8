package m3u8

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"strings"
)

// MasterPlaylist represents a master playlist.
type MasterPlaylist struct {
	// Tags is a list of tags in the master playlist.
	// This list does not include stream tags.
	Tags Tags

	// Streams is a list of variant streams.
	Streams []*Stream

	// Alternatives contains alternative renditions.
	Alternatives Alternatives

	// IFrameStreams is a list of I-frame streams.
	IFrameStreams []*Stream
}

// Stream represents a variant stream.
type Stream struct {
	// Attributes is a list of attributes in the stream.
	Attributes StreamInfAttrs

	// URI is the URI of the media playlist.
	URI string
}

type Alternatives struct {
	Video          map[string][]*Alternative
	Audio          map[string][]*Alternative
	Subtitles      map[string][]*Alternative
	ClosedCaptions map[string][]*Alternative
}

type Alternative struct {
	// Attributes is a list of attributes in the alternative.
	Attributes MediaAttrs
}

// DecodeMasterPlaylist decodes a master playlist from io.Reader.
func DecodeMasterPlaylist(r io.Reader) (*MasterPlaylist, error) {
	scanner := bufio.NewScanner(r)
	var playlist MasterPlaylist
	playlist.Tags = make(Tags)
	playlist.Streams = make([]*Stream, 0)
	playlist.Alternatives = Alternatives{
		Video:          make(map[string][]*Alternative),
		Audio:          make(map[string][]*Alternative),
		Subtitles:      make(map[string][]*Alternative),
		ClosedCaptions: make(map[string][]*Alternative),
	}
	var streamInfAttrs StreamInfAttrs
	for scanner.Scan() {
		line := scanner.Text()
		if line == "" {
			continue
		}
		tagName := TagName(line)
		if tagName == "" {
			playlist.Streams = append(playlist.Streams, &Stream{
				Attributes: streamInfAttrs,
				URI:        line,
			})
			streamInfAttrs = nil
		} else if streamInfAttrs != nil {
			return nil, errors.New("invalid EXT-X-STREAM-INF tag")
		} else if tagName == TagExtXStreamInf {
			attrs, err := ParseTagAttributes(AttributeString(line))
			if err != nil {
				return nil, err
			}
			streamInfAttrs = StreamInfAttrs(attrs)
		} else if tagName == TagExtXIFrameStreamInf {
			attrs, err := ParseTagAttributes(AttributeString(line))
			if err != nil {
				return nil, err
			}
			uri := strings.Trim(attrs["URI"], "\"")
			delete(attrs, "URI")
			playlist.IFrameStreams = append(playlist.IFrameStreams, &Stream{
				Attributes: StreamInfAttrs(attrs),
				URI:        uri,
			})
		} else if tagName == TagExtXMedia {
			attrs, err := ParseTagAttributes(AttributeString(line))
			if err != nil {
				return nil, err
			}
			groupID := strings.Trim(attrs["GROUP-ID"], `"`)
			if groupID == "" {
				return nil, errors.New("missing GROUP-ID")
			}
			typ := attrs["TYPE"]
			switch MediaType(typ) {
			case MediaTypeVideo:
				if _, ok := playlist.Alternatives.Video[groupID]; !ok {
					playlist.Alternatives.Video[groupID] = make([]*Alternative, 0)
				}
				playlist.Alternatives.Video[groupID] = append(playlist.Alternatives.Video[groupID], &Alternative{
					Attributes: MediaAttrs(attrs),
				})
			case MediaTypeAudio:
				if _, ok := playlist.Alternatives.Audio[groupID]; !ok {
					playlist.Alternatives.Audio[groupID] = make([]*Alternative, 0)
				}
				playlist.Alternatives.Audio[groupID] = append(playlist.Alternatives.Audio[groupID], &Alternative{
					Attributes: MediaAttrs(attrs),
				})
			case MediaTypeSubtitles:
				if _, ok := playlist.Alternatives.Subtitles[groupID]; !ok {
					playlist.Alternatives.Subtitles[groupID] = make([]*Alternative, 0)
				}
				playlist.Alternatives.Subtitles[groupID] = append(playlist.Alternatives.Subtitles[groupID], &Alternative{
					Attributes: MediaAttrs(attrs),
				})
			case MediaTypeClosedCaptions:
				if _, ok := playlist.Alternatives.ClosedCaptions[groupID]; !ok {
					playlist.Alternatives.ClosedCaptions[groupID] = make([]*Alternative, 0)
				}
				playlist.Alternatives.ClosedCaptions[groupID] = append(playlist.Alternatives.ClosedCaptions[groupID], &Alternative{
					Attributes: MediaAttrs(attrs),
				})
			default:
				return nil, errors.New("invalid TYPE")
			}
		} else if tagName != "" {
			playlist.Tags.Add(&Tag{
				Name:       tagName,
				Attributes: AttributeString(line),
			})
		}
	}
	return &playlist, nil
}

// Encode encodes a master playlist to io.Writer.
func (playlist *MasterPlaylist) Encode(w io.Writer) error {
	for _, tag := range playlist.Tags.List() {
		err := tag.Encode(w)
		if err != nil {
			return err
		}
	}
	for groupID, alternatives := range playlist.Alternatives.Video {
		for _, alt := range alternatives {
			if err := encodeExtXMedia(w, MediaTypeVideo, groupID, alt.Attributes); err != nil {
				return err
			}
		}
	}
	for groupID, alternatives := range playlist.Alternatives.Audio {
		for _, alt := range alternatives {
			if err := encodeExtXMedia(w, MediaTypeAudio, groupID, alt.Attributes); err != nil {
				return err
			}
		}
	}
	for groupID, alternatives := range playlist.Alternatives.Subtitles {
		for _, alt := range alternatives {
			if err := encodeExtXMedia(w, MediaTypeSubtitles, groupID, alt.Attributes); err != nil {
				return err
			}
		}
	}
	for groupID, alternatives := range playlist.Alternatives.ClosedCaptions {
		for _, alt := range alternatives {
			if err := encodeExtXMedia(w, MediaTypeClosedCaptions, groupID, alt.Attributes); err != nil {
				return err
			}
		}
	}
	for _, stream := range playlist.Streams {
		if _, err := fmt.Fprintf(w, "#%s:%s\n", TagExtXStreamInf, Attributes(stream.Attributes).String()); err != nil {
			return err
		}
		if _, err := fmt.Fprintln(w, stream.URI); err != nil {
			return err
		}
	}
	for _, stream := range playlist.IFrameStreams {
		if _, err := fmt.Fprintf(w, "#%s:%s,URI=\"%s\"\n", TagExtXIFrameStreamInf, Attributes(stream.Attributes).String(), stream.URI); err != nil {
			return err
		}
	}
	return nil
}

func encodeExtXMedia(w io.Writer, typ MediaType, groupID string, attrs MediaAttrs) error {
	a := make(Attributes, len(attrs)-2)
	for k, v := range attrs {
		if k != "TYPE" && k != "GROUP-ID" {
			a[k] = v
		}
	}
	_, err := fmt.Fprintf(w, "#%s:TYPE=%s,GROUP-ID=\"%s\",%s\n", TagExtXMedia, typ, groupID, a.String())
	return err
}

// Type returns the type of the playlist.
func (playlist *MasterPlaylist) Type() PlaylistType {
	return PlaylistTypeMaster
}

// Master returns the master playlist.
// If the playlist is not a master playlist, it returns nil.
func (playlist *MasterPlaylist) Master() *MasterPlaylist {
	return playlist
}

// Media returns the media playlist.
// If the playlist is not a media playlist, it returns nil.
func (playlist *MasterPlaylist) Media() *MediaPlaylist {
	return nil
}
