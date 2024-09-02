package metadata

import (
	"cmp"
	"fmt"
	"strings"

	"github.com/navidrome/navidrome/consts"
	"github.com/navidrome/navidrome/model"
	"github.com/navidrome/navidrome/model/id"
	"github.com/navidrome/navidrome/utils"
	"github.com/navidrome/navidrome/utils/str"
)

// FIXME Must be configurable
func (md Metadata) trackPID(mf model.MediaFile) string {
	// FIXME It will never reach the 3rd option
	value := cmp.Or(
		mf.MbzReleaseTrackID,
		fmt.Sprintf("%s\\%02d\\%02d", md.albumID(mf), mf.DiscNumber, mf.TrackNumber),
		fmt.Sprintf("%s\\%s", md.albumID(mf), md.mapTrackTitle()),
	)

	return id.NewHash(str.Clear(strings.ToLower(value)))
}

// FIXME Must be configurable
func (md Metadata) albumID(mf model.MediaFile) string {
	value := cmp.Or(
		mf.MbzAlbumID,
		func() string {
			parts := []string{
				strings.ToLower(md.mapAlbumName()),
				strings.ToLower(md.String(model.TagAlbumVersion)),
				md.String(model.TagReleaseDate),
			}
			return strings.Join(parts, "\\")
		}(),
	)
	return id.NewHash(str.Clear(strings.ToLower(value)))
}

// FIXME Must be configurable
func (md Metadata) artistID(name string) string {
	return id.NewHash(str.Clear(strings.ToLower(name)))
}

func (md Metadata) mapTrackTitle() string {
	if title := md.String(model.TagTitle); title != "" {
		return title
	}
	return utils.BaseName(md.FilePath())
}

func (md Metadata) mapAlbumName() string {
	return cmp.Or(
		md.String(model.TagAlbum),
		consts.UnknownAlbum,
	)
}
