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

	// Alternatives is a list of alternative renditions.
	Alternatives map[string][]*Alternative

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
	playlist.Alternatives = make(map[string][]*Alternative)
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
			groupID, ok := attrs["GROUP-ID"]
			if !ok {
				return nil, errors.New("missing GROUP-ID")
			}
			delete(attrs, "GROUP-ID")
			if _, ok := playlist.Alternatives[groupID]; !ok {
				playlist.Alternatives[groupID] = make([]*Alternative, 0)
			}
			playlist.Alternatives[groupID] = append(playlist.Alternatives[groupID], &Alternative{
				Attributes: MediaAttrs(attrs),
			})
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
	for groupID, alternatives := range playlist.Alternatives {
		for _, alt := range alternatives {
			if _, err := fmt.Fprintf(w, "#%s:GROUP-ID=%s,%s\n", TagExtXMedia, groupID, Attributes(alt.Attributes).String()); err != nil {
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
