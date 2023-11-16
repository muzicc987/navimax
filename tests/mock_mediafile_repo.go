package tests

import (
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/navidrome/navidrome/utils/slice"
	"golang.org/x/exp/maps"

	"github.com/navidrome/navidrome/model"
)

func CreateMockMediaFileRepo(mfl model.MediaFolderRepository) *MockMediaFileRepo {
	return &MockMediaFileRepo{
		data: make(map[string]*model.MediaFile),
	}
}

type MockMediaFileRepo struct {
	model.MediaFileRepository
	data map[string]*model.MediaFile
	err  bool
	mfl  MockMediaFolderRepo
}

func (m *MockMediaFileRepo) SetError(err bool) {
	m.err = err
}

func (m *MockMediaFileRepo) SetData(mfs model.MediaFiles) {
	m.data = make(map[string]*model.MediaFile)
	for i, mf := range mfs {
		m.data[mf.ID] = &mfs[i]
	}
}

func (m *MockMediaFileRepo) Exists(id string) (bool, error) {
	if m.err {
		return false, errors.New("error")
	}
	_, found := m.data[id]
	return found, nil
}

func (m *MockMediaFileRepo) Get(id string) (*model.MediaFile, error) {
	if m.err {
		return nil, errors.New("error")
	}
	if d, ok := m.data[id]; ok {
		return d, nil
	}
	return nil, model.ErrNotFound
}

func (m *MockMediaFileRepo) GetAll(...model.QueryOptions) (model.MediaFiles, error) {
	if m.err {
		return nil, errors.New("error")
	}
	values := maps.Values(m.data)
	return slice.Map(values, func(p *model.MediaFile) model.MediaFile {
		return *p
	}), nil
}

func (m *MockMediaFileRepo) Put(mf *model.MediaFile) error {
	if m.err {
		return errors.New("error")
	}
	if mf.ID == "" {
		mf.ID = uuid.NewString()
	}
	m.data[mf.ID] = mf
	return nil
}

func (m *MockMediaFileRepo) IncPlayCount(id string, timestamp time.Time) error {
	if m.err {
		return errors.New("error")
	}
	if d, ok := m.data[id]; ok {
		d.PlayCount++
		d.PlayDate = timestamp
		return nil
	}
	return model.ErrNotFound
}

func (m *MockMediaFileRepo) FindByAlbum(artistId string) (model.MediaFiles, error) {
	if m.err {
		return nil, errors.New("error")
	}
	var res = make(model.MediaFiles, len(m.data))
	i := 0
	for _, a := range m.data {
		if a.AlbumID == artistId {
			res[i] = *a
			i++
		}
	}

	return res, nil
}

func (m *MockMediaFileRepo) FindAllByParent(parent_id string) (model.MediaFiles, error) {
	if m.err {
		return nil, errors.New("error")
	}
	mfs := model.MediaFiles{}

	for _, folder := range m.mfl.Data {
		if folder.ParentId == parent_id && folder.Path != "" {
			file, ok := m.data[folder.ID]
			if !ok {
				return nil, model.ErrNotFound
			}

			mfs = append(mfs, *file)
		}
	}
	return mfs, nil
}

func (r *MockMediaFileRepo) DeleteByPath(basePath string) (int64, error) {
	return 0, errors.New(basePath)
}

var _ model.MediaFileRepository = (*MockMediaFileRepo)(nil)
