package m3u8

import (
	"bufio"
	"bytes"
	"io"
)

// PlaylistType represents the type of playlist.
type PlaylistType string

const (
	// PlaylistTypeMaster represents a master playlist.
	PlaylistTypeMaster PlaylistType = "master"
	// PlaylistTypeMedia represents a media playlist.
	PlaylistTypeMedia PlaylistType = "media"
)

type Playlist interface {
	// Encode encodes the playlist to io.Writer.
	Encode(w io.Writer) error

	// Type returns the type of the playlist.
	Type() PlaylistType

	// Master returns the master playlist.
	// If the playlist is not a master playlist, it returns nil.
	Master() *MasterPlaylist

	// Media returns the media playlist.
	// If the playlist is not a media playlist, it returns nil.
	Media() *MediaPlaylist
}

// DecodePlaylist detects the type of playlist and decodes it from io.Reader.
func DecodePlaylist(r io.Reader) (Playlist, error) {
	data, err := io.ReadAll(r)
	if err != nil {
		return nil, err
	}
	br := bytes.NewReader(data)

	var masterPlaylistTagCount int
	var mediaPlaylistTagCount int
	scanner := bufio.NewScanner(br)
	for scanner.Scan() {
		line := scanner.Text()
		if line == "" {
			continue
		}
		tagName := TagName(line)
		if isMasterPlaylistTag(tagName) {
			masterPlaylistTagCount++
		} else if isMediaPlaylistTag(tagName) || IsSegmentTagName(tagName) {
			mediaPlaylistTagCount++
		}
	}
	if err := scanner.Err(); err != nil {
		return nil, err
	}

	br.Seek(0, io.SeekStart)
	if masterPlaylistTagCount >= mediaPlaylistTagCount {
		return DecodeMasterPlaylist(br)
	} else {
		return DecodeMediaPlaylist(br)
	}
}
