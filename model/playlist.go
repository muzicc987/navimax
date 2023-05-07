package model

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/navidrome/navidrome/model/criteria"
	"golang.org/x/exp/slices"
)

type Playlist struct {
	ID        string         `structs:"id" json:"id"          orm:"column(id)"`
	Name      string         `structs:"name" json:"name"`
	Comment   string         `structs:"comment" json:"comment"`
	Duration  float32        `structs:"duration" json:"duration"`
	Size      int64          `structs:"size" json:"size"`
	SongCount int            `structs:"song_count" json:"songCount"`
	OwnerName string         `structs:"-" json:"ownerName"`
	OwnerID   string         `structs:"owner_id" json:"ownerId"  orm:"column(owner_id)"`
	Public    bool           `structs:"public" json:"public"`
	Tracks    PlaylistTracks `structs:"-" json:"tracks,omitempty"`
	Path      string         `structs:"path" json:"path"`
	Sync      bool           `structs:"sync" json:"sync"`
	CreatedAt time.Time      `structs:"created_at" json:"createdAt"`
	UpdatedAt time.Time      `structs:"updated_at" json:"updatedAt"`

	// External Info
	ExternalAgent       string `structs:"external_agent" json:"external_agent"`
	ExternalId          string `structs:"external_id" json:"externalId"`
	ExternalSync        bool   `structs:"external_sync" json:"externalSync"`
	ExternalSyncable    bool   `structs:"external_syncable" json:"externalSyncable"`
	ExternalUrl         string `structs:"external_url" json:"externalUrl"`
	ExternalRecommended bool   `structs:"external_recommended" json:"externalRecommended"`

	// SmartPlaylist attributes
	Rules       *criteria.Criteria `structs:"-" json:"rules"`
	EvaluatedAt time.Time          `structs:"evaluated_at" json:"evaluatedAt"`
}

func (pls Playlist) IsSmartPlaylist() bool {
	return pls.Rules != nil && pls.Rules.Expression != nil
}

func (pls Playlist) MediaFiles() MediaFiles {
	if len(pls.Tracks) == 0 {
		return nil
	}
	return pls.Tracks.MediaFiles()
}

func (pls *Playlist) RemoveTracks(idxToRemove []int) {
	var newTracks PlaylistTracks
	for i, t := range pls.Tracks {
		if slices.Contains(idxToRemove, i) {
			continue
		}
		newTracks = append(newTracks, t)
	}
	pls.Tracks = newTracks
}

// ToM3U8 exports the playlist to the Extended M3U8 format, as specified in
// https://docs.fileformat.com/audio/m3u/#extended-m3u
func (pls *Playlist) ToM3U8() string {
	buf := strings.Builder{}
	buf.WriteString("#EXTM3U\n")
	buf.WriteString(fmt.Sprintf("#PLAYLIST:%s\n", pls.Name))
	for _, t := range pls.Tracks {
		buf.WriteString(fmt.Sprintf("#EXTINF:%.f,%s - %s\n", t.Duration, t.Artist, t.Title))
		buf.WriteString(t.Path + "\n")
	}
	return buf.String()
}

func (pls *Playlist) AddTracks(mediaFileIds []string) {
	pos := len(pls.Tracks)
	for _, mfId := range mediaFileIds {
		pos++
		t := PlaylistTrack{
			ID:          strconv.Itoa(pos),
			MediaFileID: mfId,
			MediaFile:   MediaFile{ID: mfId},
			PlaylistID:  pls.ID,
		}
		pls.Tracks = append(pls.Tracks, t)
	}
}

func (pls *Playlist) AddMediaFiles(mfs MediaFiles) {
	pos := len(pls.Tracks)
	for _, mf := range mfs {
		pos++
		t := PlaylistTrack{
			ID:          strconv.Itoa(pos),
			MediaFileID: mf.ID,
			MediaFile:   mf,
			PlaylistID:  pls.ID,
		}
		pls.Tracks = append(pls.Tracks, t)
	}
}

func (pls Playlist) CoverArtID() ArtworkID {
	return artworkIDFromPlaylist(pls)
}

type Playlists []Playlist

type PlaylistRepository interface {
	ResourceRepository
	CountAll(options ...QueryOptions) (int64, error)
	Exists(id string) (bool, error)
	Put(pls *Playlist) error
	Get(id string) (*Playlist, error)
	GetSyncedPlaylists() (Playlists, error)
	GetWithTracks(id string, refreshSmartPlaylist bool) (*Playlist, error)
	GetAll(options ...QueryOptions) (Playlists, error)
	CheckExternalIds(agent string, ids []string) ([]string, error)
	GetByExternalInfo(agent, id string) (*Playlist, error)
	GetRecommended(userId, agent string) (*Playlist, error)
	FindByPath(path string) (*Playlist, error)
	Delete(id string) error
	Tracks(playlistId string, refreshSmartPlaylist bool) PlaylistTrackRepository
}

type PlaylistTrack struct {
	ID          string `json:"id"          orm:"column(id)"`
	MediaFileID string `json:"mediaFileId" orm:"column(media_file_id)"`
	PlaylistID  string `json:"playlistId" orm:"column(playlist_id)"`
	MediaFile
}

type PlaylistTracks []PlaylistTrack

func (plt PlaylistTracks) MediaFiles() MediaFiles {
	mfs := make(MediaFiles, len(plt))
	for i, t := range plt {
		mfs[i] = t.MediaFile
	}
	return mfs
}

type PlaylistTrackRepository interface {
	ResourceRepository
	GetAll(options ...QueryOptions) (PlaylistTracks, error)
	GetAlbumIDs(options ...QueryOptions) ([]string, error)
	Add(mediaFileIds []string) (int, error)
	AddAlbums(albumIds []string) (int, error)
	AddArtists(artistIds []string) (int, error)
	AddDiscs(discs []DiscID) (int, error)
	Delete(id ...string) error
	DeleteAll() error
	Reorder(pos int, newPos int) error
}
