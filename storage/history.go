package storage

import (
	"encoding/json"
	"os"
	"path/filepath"
	"sync"
	"time"

	"github.com/shubhambadola/VidFetch/downloader"
)

type History struct {
	Downloads []downloader.Download `json:"downloads"`
	LastSync  time.Time             `json:"last_sync"`
	path      string
	mu        sync.RWMutex
}

func NewHistory(path string) (*History, error) {
	h := &History{
		Downloads: []downloader.Download{},
		path:      path,
	}

	// Ensure dir exists
	dir := filepath.Dir(path)
	os.MkdirAll(dir, 0755)

	if err := h.Load(); err != nil {
		if os.IsNotExist(err) {
			return h, nil
		}
		return nil, err
	}
	return h, nil
}

func (h *History) Load() error {
	h.mu.Lock()
	defer h.mu.Unlock()

	data, err := os.ReadFile(h.path)
	if err != nil {
		return err
	}
	return json.Unmarshal(data, h)
}

func (h *History) Save() error {
	h.mu.Lock()
	defer h.mu.Unlock()

	h.LastSync = time.Now()
	data, err := json.MarshalIndent(h, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(h.path, data, 0644)
}

func (h *History) Add(dl downloader.Download) error {
	h.mu.Lock()
	// Prepend
	h.Downloads = append([]downloader.Download{dl}, h.Downloads...)
	h.mu.Unlock()
	return h.Save()
}

func (h *History) Get() []downloader.Download {
	h.mu.RLock()
	defer h.mu.RUnlock()
	// Return copy
	dst := make([]downloader.Download, len(h.Downloads))
	copy(dst, h.Downloads)
	return dst
}
