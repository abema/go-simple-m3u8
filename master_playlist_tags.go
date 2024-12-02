package m3u8

import (
	"errors"
	"strconv"
	"strings"
)

// StreamInfAttrs represents the attributes of the EXT-X-STREAM-INF tag.
type StreamInfAttrs Attributes

// Resolution returns the resolution of the stream.
func (attrs StreamInfAttrs) Resolution() (width, height int, err error) {
	return ParseResolution(attrs["RESOLUTION"])
}

// ParseResolution parses the resolution string.
func ParseResolution(resolution string) (width, height int, err error) {
	idx := strings.Index(resolution, "x")
	if idx == -1 || idx == 0 || idx == len(resolution)-1 {
		return 0, 0, errors.New("invalid resolution")
	}
	width, err = strconv.Atoi(resolution[:idx])
	if err != nil {
		return 0, 0, err
	}
	height, err = strconv.Atoi(resolution[idx+1:])
	if err != nil {
		return 0, 0, err
	}
	return width, height, nil
}

// MediaAttrs represents the attributes of the EXT-X-MEDIA tag.
type MediaAttrs Attributes
