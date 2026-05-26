package matrix

import (
	"context"
	"encoding/json"
	"os"
	"sync"

	"maunium.net/go/mautrix/id"
)

var DefaultFileSyncStorePath = "sync_state.json"

type syncData struct {
	NextBatch map[string]string `json:"next_batch"`
	FilterID  map[string]string `json:"filter_id"`
}

type FileSyncStore struct {
	path  string
	mutex sync.Mutex
	data  syncData
}

func NewDefaultFileSyncStore() *FileSyncStore {
	return NewFileSyncStore(DefaultFileSyncStorePath)
}

func NewFileSyncStore(path string) *FileSyncStore {
	store := &FileSyncStore{
		path: path,
		data: syncData{
			NextBatch: map[string]string{},
			FilterID:  map[string]string{},
		},
	}

	file, err := os.ReadFile(path)
	if err == nil {
		_ = json.Unmarshal(file, &store.data)
	}

	return store
}

func (s *FileSyncStore) save() error {
	tmp := s.path + ".tmp"

	b, err := json.MarshalIndent(s.data, "", "  ")
	if err != nil {
		return err
	}

	if err := os.WriteFile(tmp, b, 0644); err != nil {
		return err
	}

	return os.Rename(tmp, s.path)
}

func mutate(s *FileSyncStore, handler func()) error {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	handler()
	return s.save()
}

func (s *FileSyncStore) SaveFilterID(_ context.Context, userID id.UserID, filterID string) error {
	return mutate(s, func() {
		s.data.FilterID[userID.String()] = filterID
	})
}

func (s *FileSyncStore) LoadFilterID(_ context.Context, userID id.UserID) (string, error) {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	return s.data.FilterID[userID.String()], nil
}

func (s *FileSyncStore) SaveNextBatch(_ context.Context, userID id.UserID, token string) error {
	return mutate(s, func() {
		s.data.NextBatch[userID.String()] = token
	})
}

func (s *FileSyncStore) LoadNextBatch(_ context.Context, userID id.UserID) (string, error) {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	return s.data.NextBatch[userID.String()], nil
}
