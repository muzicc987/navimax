package playback

import (
	"context"
	"mime"
	"os"
	"time"

	"github.com/faiface/beep"
	"github.com/faiface/beep/mp3"
	"github.com/faiface/beep/speaker"
	"github.com/navidrome/navidrome/log"
	"github.com/navidrome/navidrome/model"
	"github.com/navidrome/navidrome/model/request"
	"github.com/navidrome/navidrome/server/subsonic/responses"
)

type PlaybackDevice struct {
	Ctx        context.Context
	DataStore  model.DataStore
	Default    bool
	User       string
	Name       string
	Method     string
	DeviceName string
	Playlist   responses.JukeboxPlaylist
	Ctrl       *beep.Ctrl
}

func (pd *PlaybackDevice) Get(user string) (responses.JukeboxPlaylist, error) {
	log.Debug("processing Get action")
	return responses.JukeboxPlaylist{}, nil
}

func (pd *PlaybackDevice) Status(user string) (responses.JukeboxStatus, error) {
	log.Debug("processing Status action")
	return responses.JukeboxStatus{}, nil
}
func (pd *PlaybackDevice) Set(user string, ids []string) (responses.JukeboxStatus, error) {
	log.Debug("processing Set action.")

	mf, err := pd.DataStore.MediaFile(pd.Ctx).Get(ids[0])
	if err != nil {
		return responses.JukeboxStatus{}, err
	}

	log.Debug("Found mediafile: " + mf.Path)

	child := childFromMediaFile(pd.Ctx, *mf)

	pd.Playlist.Entry[0] = child
	pd.prepareSong(child.Path)

	return responses.JukeboxStatus{}, nil
}
func (pd *PlaybackDevice) Start(user string) (responses.JukeboxStatus, error) {
	log.Debug("processing Start action")
	pd.playHead()
	return responses.JukeboxStatus{}, nil
}
func (pd *PlaybackDevice) Stop(user string) (responses.JukeboxStatus, error) {
	log.Debug("processing Stop action")
	pd.pauseHead()
	return responses.JukeboxStatus{}, nil
}
func (pd *PlaybackDevice) Skip(user string, index int, offset int) (responses.JukeboxStatus, error) {
	log.Debug("processing Skip action")
	return responses.JukeboxStatus{}, nil
}
func (pd *PlaybackDevice) Add(user string, ids []string) (responses.JukeboxStatus, error) {
	log.Debug("processing Add action")
	// pd.Playlist.Entry = append(pd.Playlist.Entry, child)
	return responses.JukeboxStatus{}, nil
}
func (pd *PlaybackDevice) Clear(user string) (responses.JukeboxStatus, error) {
	log.Debug("processing Clear action")
	return responses.JukeboxStatus{}, nil
}
func (pd *PlaybackDevice) Remove(user string, index int) (responses.JukeboxStatus, error) {
	log.Debug("processing Remove action")
	return responses.JukeboxStatus{}, nil
}
func (pd *PlaybackDevice) Shuffle(user string) (responses.JukeboxStatus, error) {
	log.Debug("processing Shuffle action")
	return responses.JukeboxStatus{}, nil
}
func (pd *PlaybackDevice) SetGain(user string, gain float64) (responses.JukeboxStatus, error) {
	log.Debug("processing SetGain action")
	return responses.JukeboxStatus{}, nil
}

func (pd *PlaybackDevice) playHead() {
	speaker.Lock()
	pd.Ctrl.Paused = false
	speaker.Unlock()
}

func (pd *PlaybackDevice) pauseHead() {
	speaker.Lock()
	pd.Ctrl.Paused = true
	speaker.Unlock()
}

func (pd *PlaybackDevice) prepareSong(songname string) {
	log.Debug("Playing song: " + songname)
	f, err := os.Open(songname)
	if err != nil {
		log.Fatal(err)
	}

	streamer, format, err := mp3.Decode(f)
	if err != nil {
		log.Fatal(err)
	}

	pd.Ctrl.Streamer = streamer
	pd.Ctrl.Paused = true

	err = speaker.Init(format.SampleRate, format.SampleRate.N(time.Second/10))
	if err != nil {
		log.Fatal(err)
	}

	go func() {
		speaker.Play(pd.Ctrl)
	}()

}

func getTranscoding(ctx context.Context) (format string, bitRate int) {
	if trc, ok := request.TranscodingFrom(ctx); ok {
		format = trc.TargetFormat
	}
	if plr, ok := request.PlayerFrom(ctx); ok {
		bitRate = plr.MaxBitRate
	}
	return
}

// FIXME: this is a copy from subsonic/helpers.go consolidate.
func childFromMediaFile(ctx context.Context, mf model.MediaFile) responses.Child {
	child := responses.Child{}
	child.Id = mf.ID
	child.Title = mf.Title
	child.IsDir = false
	child.Parent = mf.AlbumID
	child.Album = mf.Album
	child.Year = mf.Year
	child.Artist = mf.Artist
	child.Genre = mf.Genre
	child.Track = mf.TrackNumber
	child.Duration = int(mf.Duration)
	child.Size = mf.Size
	child.Suffix = mf.Suffix
	child.BitRate = mf.BitRate
	child.CoverArt = mf.CoverArtID().String()
	child.ContentType = mf.ContentType()
	child.Path = mf.Path
	child.DiscNumber = mf.DiscNumber
	child.Created = &mf.CreatedAt
	child.AlbumId = mf.AlbumID
	child.ArtistId = mf.ArtistID
	child.Type = "music"
	child.PlayCount = mf.PlayCount
	if mf.PlayCount > 0 {
		child.Played = &mf.PlayDate
	}
	if mf.Starred {
		child.Starred = &mf.StarredAt
	}
	child.UserRating = mf.Rating

	format, _ := getTranscoding(ctx)
	if mf.Suffix != "" && format != "" && mf.Suffix != format {
		child.TranscodedSuffix = format
		child.TranscodedContentType = mime.TypeByExtension("." + format)
	}
	child.BookmarkPosition = mf.BookmarkPosition
	return child
}
