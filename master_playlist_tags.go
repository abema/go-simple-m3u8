package m3u8

import (
	"errors"
	"strconv"
	"strings"
)

// MediaType represents the media type of the stream.
type MediaType string

const (
	MediaTypeAudio          MediaType = "AUDIO"
	MediaTypeVideo          MediaType = "VIDEO"
	MediaTypeSubtitles      MediaType = "SUBTITLES"
	MediaTypeClosedCaptions MediaType = "CLOSED-CAPTIONS"
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

// Bandwidth returns the bandwidth of the stream.
func (attrs StreamInfAttrs) Bandwidth() (int64, error) {
	return strconv.ParseInt(attrs["BANDWIDTH"], 10, 64)
}

// SetBandwidth sets the bandwidth of the stream.
func (attrs StreamInfAttrs) SetBandwidth(bandwidth int64) {
	attrs["BANDWIDTH"] = strconv.FormatInt(bandwidth, 10)
}

// AverageBandwidth returns the average bandwidth of the stream.
func (attrs StreamInfAttrs) AverageBandwidth() (int64, error) {
	return strconv.ParseInt(attrs["AVERAGE-BANDWIDTH"], 10, 64)
}

// SetAverageBandwidth sets the average bandwidth of the stream.
func (attrs StreamInfAttrs) SetAverageBandwidth(bandwidth int64) {
	attrs["AVERAGE-BANDWIDTH"] = strconv.FormatInt(bandwidth, 10)
}

// Codecs returns the codecs of the stream.
func (attrs StreamInfAttrs) Codecs() []string {
	return strings.Split(strings.Trim(attrs["CODECS"], `"`), ",")
}

// SetCodecs sets the codecs of the stream.
func (attrs StreamInfAttrs) SetCodecs(codecs []string) {
	attrs["CODECS"] = `"` + strings.Join(codecs, ",") + `"`
}

// FrameRate returns the frame rate of the stream.
func (attrs StreamInfAttrs) FrameRate() (float64, error) {
	return strconv.ParseFloat(attrs["FRAME-RATE"], 64)
}

// SetFrameRate sets the frame rate of the stream.
func (attrs StreamInfAttrs) SetFrameRate(frameRate float64) {
	attrs["FRAME-RATE"] = strconv.FormatFloat(frameRate, 'f', -1, 64)
}

// Audio returns the audio group of the stream.
func (attrs StreamInfAttrs) Audio() string {
	return strings.Trim(attrs["AUDIO"], `"`)
}

// SetAudio sets the audio group of the stream.
func (attrs StreamInfAttrs) SetAudio(audio string) {
	attrs["AUDIO"] = `"` + audio + `"`
}

// Video returns the video group of the stream.
func (attrs StreamInfAttrs) Video() string {
	return strings.Trim(attrs["VIDEO"], `"`)
}

// SetVideo sets the video group of the stream.
func (attrs StreamInfAttrs) SetVideo(video string) {
	attrs["VIDEO"] = `"` + video + `"`
}

// Subtitles returns the subtitles group of the stream.
func (attrs StreamInfAttrs) Subtitles() string {
	return strings.Trim(attrs["SUBTITLES"], `"`)
}

// SetSubtitles sets the subtitles group of the stream.
func (attrs StreamInfAttrs) SetSubtitles(subtitles string) {
	attrs["SUBTITLES"] = `"` + subtitles + `"`
}

// ClosedCaptions returns the closed captions group of the stream.
func (attrs StreamInfAttrs) ClosedCaptions() string {
	return strings.Trim(attrs["CLOSED-CAPTIONS"], `"`)
}

// SetClosedCaptions sets the closed captions group of the stream.
func (attrs StreamInfAttrs) SetClosedCaptions(closedCaptions string) {
	attrs["CLOSED-CAPTIONS"] = `"` + closedCaptions + `"`
}

// MediaAttrs represents the attributes of the EXT-X-MEDIA tag.
type MediaAttrs Attributes

// Type returns the type of the media.
func (attrs MediaAttrs) Type() MediaType {
	return MediaType(attrs["TYPE"])
}

// SetType sets the type of the media.
func (attrs MediaAttrs) SetType(mediaType MediaType) {
	attrs["TYPE"] = string(mediaType)
}

// URI returns the URI of the media.
func (attrs MediaAttrs) URI() string {
	return strings.Trim(attrs["URI"], `"`)
}

// SetURI sets the URI of the media.
func (attrs MediaAttrs) SetURI(uri string) {
	attrs["URI"] = `"` + uri + `"`
}

// GroupID returns the group ID of the media.
func (attrs MediaAttrs) GroupID() string {
	return strings.Trim(attrs["GROUP-ID"], `"`)
}

// SetGroupID sets the group ID of the media.
func (attrs MediaAttrs) SetGroupID(groupID string) {
	attrs["GROUP-ID"] = `"` + groupID + `"`
}

// Language returns the language of the media.
func (attrs MediaAttrs) Language() string {
	return strings.Trim(attrs["LANGUAGE"], `"`)
}

// SetLanguage sets the language of the media.
func (attrs MediaAttrs) SetLanguage(language string) {
	attrs["LANGUAGE"] = `"` + language + `"`
}

// AssocLanguage returns the associated language of the media.
func (attrs MediaAttrs) AssocLanguage() string {
	return strings.Trim(attrs["ASSOC-LANGUAGE"], `"`)
}

// SetAssocLanguage sets the associated language of the media.
func (attrs MediaAttrs) SetAssocLanguage(language string) {
	attrs["ASSOC-LANGUAGE"] = `"` + language + `"`
}

// Name returns the name of the media.
func (attrs MediaAttrs) Name() string {
	return strings.Trim(attrs["NAME"], `"`)
}

// SetName sets the name of the media.
func (attrs MediaAttrs) SetName(name string) {
	attrs["NAME"] = `"` + name + `"`
}

// Default returns the default flag of the media.
func (attrs MediaAttrs) Default() bool {
	return attrs["DEFAULT"] == "YES"
}

// SetDefault sets the default flag of the media.
func (attrs MediaAttrs) SetDefault(defaultFlag bool) {
	if defaultFlag {
		attrs["DEFAULT"] = "YES"
	} else {
		attrs["DEFAULT"] = "NO"
	}
}

// Autoselect returns the autoselect flag of the media.
func (attrs MediaAttrs) Autoselect() bool {
	return attrs["AUTOSELECT"] == "YES"
}

// SetAutoselect sets the autoselect flag of the media.
func (attrs MediaAttrs) SetAutoselect(autoselect bool) {
	if autoselect {
		attrs["AUTOSELECT"] = "YES"
	} else {
		attrs["AUTOSELECT"] = "NO"
	}
}
