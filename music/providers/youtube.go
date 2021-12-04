package music

import (
	"net/http"
	"regexp"
	"strconv"
	"time"

	"github.com/XzFrosT/crunch/utils"
	"github.com/XzFrosT/crunch/utils/emojis"
	"github.com/Pauloo27/searchtube"
	"github.com/kkdai/youtube/v2"
	"github.com/pkg/errors"
)

type YoutubeProvider struct{}

var (
	videoRegex    = regexp.MustCompile(`^((?:https?:)?\/\/)?((?:www|m)\.)?((?:youtube\.com|youtu.be))(\/(?:[\w\-]+\?v=|embed\/|v\/)?)([\w\-]+)(\S+)?$`)
	playlistRegex = regexp.MustCompile(`[&?]list=([A-Za-z0-9_-]{18,42})(&.*)?$`)
	hlsRegex      = regexp.MustCompile(`(https?:\/\/(www\.)?[-a-zA-Z0-9@:%._\+~#=]{2,256}\.[a-z]{2,6}\b([-a-zA-Z0-9@:%_\+.~#,?&*//=]*)(.m3u8)\b([-a-zA-Z0-9@:%_\+.~#,?&//=]*))`)

	client = &youtube.Client{}
	cache  = make(map[string]*Song)
)

func (YoutubeProvider) DisplayName() string {
	return utils.Fmt("%s YouTube", emojis.Youtube)
}

func (YoutubeProvider) IsSupported(term string) bool {
	return videoRegex.MatchString(term) || !utils.LinkRegex.MatchString(term) || playlistRegex.MatchString(term)
}

func (YoutubeProvider) IsLoaded(s *Song) bool {
	return !s.Expires.IsZero() && !time.Now().Add(s.Duration).After(s.Expires)
}

func (ytProvider YoutubeProvider) Load(s *Song) error {
	loadedSong, err := ytProvider.handleVideo(s.URL, 0)
	if err != nil {
		return err
	}

	s.StreamingURL, s.IsOpus, s.Expires = loadedSong.StreamingURL, loadedSong.IsOpus, loadedSong.Expires
	s.Thumbnail = loadedSong.Thumbnail
	return nil
}

func (ytProvider *YoutubeProvider) Find(term string) (*QueryResult, error) {
	if playlistRegex.MatchString(term) {
		return ytProvider.handlePlaylist(term)
	}

	if videoRegex.MatchString(term) {
		if video, err := ytProvider.handleVideo(term, 0); err != nil {
			return nil, err
		} else {
			return &QueryResult{Songs: []*Song{video}}, nil
		}
	}

	items, err := searchtube.Search(term, 5)
	if err != nil {
		return nil, err
	}

	if len(items) < 1 {
		return nil, nil
	}

	result := &QueryResult{}
	for _, video := range items {
		duration, _ := video.GetDuration()

		result.Songs = append(result.Songs, &Song{
			Title:     video.Title,
			URL:       video.URL,
			Author:    video.Uploader,
			Thumbnail: video.Thumbnail,
			Duration:  duration,
			IsLive:    video.Live,
			provider:  ytProvider,
		})
	}

	return result, nil
}

func (ytProvider *YoutubeProvider) handlePlaylist(URL string) (*QueryResult, error) {
	playlist, err := client.GetPlaylist(URL)
	if err != nil {
		return nil, err
	}

	result := &QueryResult{
		Songs: []*Song{},
		Playlist: &Playlist{
			Title:  playlist.Title,
			Author: playlist.Author,
			URL:    utils.Fmt("https://youtube.com/playlist?list=%s", playlist.ID),
		},
	}

	for _, item := range playlist.Videos {
		result.Playlist.Duration += item.Duration

		result.Songs = append(result.Songs, &Song{
			Title:     item.Title,
			Author:    item.Author,
			Duration:  item.Duration,
			Thumbnail: utils.Fmt("https://img.youtube.com/vi/%s/mqdefault.jpg", item.ID),
			URL:       utils.Fmt("https://youtu.be/%s", item.ID),
			provider:  ytProvider,
		})
	}

	return result, nil
}

func (ytProvier YoutubeProvider) handleVideo(term string, attempts int) (song *Song, err error) {
	if term, err = youtube.ExtractVideoID(term); err != nil {
		return nil, err
	}

	if cached := cache[term]; cached != nil && ytProvier.IsLoaded(cached) {
		return cached, nil
	}

	video, err := client.GetVideo(term)
	if err != nil {
		return nil, err
	}

	streamingURL, isOpus := "", false
	if video.HLSManifestURL == "" {
		var format *youtube.Format

		if format, isOpus = video.Formats.FindByItag(251), true; format == nil { // Opus
			format, isOpus = video.Formats.FindByItag(140), false // M4a
		}

		if streamingURL, err = client.GetStreamURL(video, format); err != nil {
			return nil, err
		}
	} else {
		if streamingURL, err = getLiveURL(video.HLSManifestURL); err != nil {
			return nil, err
		}
	}

	expires, err := getExpires(streamingURL)
	if err != nil {
		if attempts >= 5 {
			return nil, err
		}
		return ytProvier.handleVideo(term, attempts+1)
	}

	song = &Song{
		Title:        video.Title,
		Author:       video.Author,
		URL:          utils.Fmt("https://youtu.be/%s", video.ID),
		Duration:     video.Duration,
		Thumbnail:    video.Thumbnails[len(video.Thumbnails)-1].URL,
		StreamingURL: streamingURL,
		Expires:      expires,
		IsLive:       video.HLSManifestURL != "",
		IsOpus:       isOpus,
		provider:     &ytProvier,
	}

	if !song.IsLive {
		cache[term] = song
	}
	return song, nil
}

func getExpires(streamingURL string) (time.Time, error) {
	response, err := http.Get(streamingURL)
	if err != nil {
		return time.Time{}, err
	}
	response.Body.Close()

	if response.StatusCode >= 400 {
		return time.Time{}, errors.Errorf("the server responded with unexpected %s status", response.Status)
	}

	expires, _ := strconv.Atoi(response.Request.URL.Query().Get("expire"))
	return time.Unix(int64(expires), 0), nil
}

func getLiveURL(manifestURL string) (string, error) {
	body, err := utils.FromWebString(manifestURL)
	if err != nil {
		return "", err
	}

	if hlsURL := hlsRegex.FindString(body); hlsURL != "" {
		return hlsURL, nil
	} else {
		return "", errors.New("no valid URL found within HLS")
	}
}