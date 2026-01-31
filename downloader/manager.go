package downloader

import (
	"context"
	"fmt"
	"sync"
	"time"
)

// Download represents the state of a single download
type Download struct {
	ID            string          `json:"id"`
	URL           string          `json:"url"`
	Title         string          `json:"title"`
	Platform      string          `json:"platform"`
	Status        string          `json:"status"` // pending, downloading, merging, completed, failed
	Progress      float64         `json:"progress"`
	Speed         string          `json:"speed"`
	ETA           string          `json:"eta"`
	FileSize      int64           `json:"file_size"`
	Downloaded    int64           `json:"downloaded"`
	FilePath      string          `json:"file_path"`
	Thumbnail     string          `json:"thumbnail"`
	Duration      int             `json:"duration"`
	Quality       string          `json:"quality"`
	Format        string          `json:"format"`
	SubtitleCount int             `json:"subtitle_count"`
	SubtitleLangs []string        `json:"subtitle_langs"`
	CreatedAt     time.Time       `json:"created_at"`
	CompletedAt   time.Time       `json:"completed_at"`
	Error         string          `json:"error"`
	Options       DownloadOptions `json:"options"` // Store options for retry/resume
}

// DownloadOptions configures the download parameters
type DownloadOptions struct {
	// Quality settings
	Format      string `json:"format"`       // "best", "1080p", "720p", "audio"
	VideoFormat string `json:"video_format"` // "mp4", "mkv", "webm"
	AudioOnly   bool   `json:"audio_only"`

	// Subtitle settings
	DownloadSubs     bool     `json:"download_subs"`
	DownloadAutoSubs bool     `json:"download_auto_subs"`
	SubtitleLangs    []string `json:"subtitle_langs"`
	EmbedSubtitles   bool     `json:"embed_subtitles"`
	SubtitleFormat   string   `json:"subtitle_format"`
	ExternalSubFiles bool     `json:"external_sub_files"`

	// Output settings
	OutputDir      string `json:"output_dir"`
	OutputTemplate string `json:"output_template"`

	// Download behavior
	NoPlaylist    bool `json:"no_playlist"`
	PlaylistStart int  `json:"playlist_start"`
	PlaylistEnd   int  `json:"playlist_end"`
}

// Downloader manages the download queue and yt-dlp execution
type Downloader struct {
	mu         sync.RWMutex
	downloads  map[string]*Download
	queue      chan string // Queue of download IDs to process
	max        int
	OnComplete func(*Download) // Callback for persistence
}

func NewDownloader(maxConcurrent int) *Downloader {
	return &Downloader{
		downloads: make(map[string]*Download),
		queue:     make(chan string, 100), // Buffer for queue
		max:       maxConcurrent,
	}
}

// Start initializes the worker pool
func (d *Downloader) Start(ctx context.Context) {
	for i := 0; i < d.max; i++ {
		go d.worker(ctx)
	}
}

func (d *Downloader) worker(ctx context.Context) {
	for id := range d.queue {
		// Identify download
		dl := d.GetDownload(id)
		if dl == nil {
			continue
		}

		// Update status
		d.mu.Lock()
		dl.Status = "downloading"
		d.mu.Unlock()

		// Execute
		err := d.downloadWithSubtitles(ctx, dl.ID, dl.Options)

		d.mu.Lock()
		dl.CompletedAt = time.Now()
		if err != nil {
			dl.Status = "failed"
			dl.Error = err.Error()
		} else {
			dl.Status = "completed"
			dl.Progress = 1.0
		}

		// Callback if set
		if d.OnComplete != nil {
			go d.OnComplete(dl)
		}
		d.mu.Unlock()
	}
}

// QueueDownload adds a download to the queue
func (d *Downloader) QueueDownload(url string, opts DownloadOptions) string {
	d.mu.Lock()
	defer d.mu.Unlock()

	id := fmt.Sprintf("dl_%d", time.Now().UnixNano())
	dl := &Download{
		ID:            id,
		URL:           url,
		Status:        "pending",
		CreatedAt:     time.Now(),
		SubtitleLangs: opts.SubtitleLangs,
		Quality:       opts.Format,
		Options:       opts,
	}
	d.downloads[id] = dl

	// Push to queue (non-blocking if buffer not full, but we should handle it)
	go func() {
		d.queue <- id
	}()

	return id
}

// GetDownload retrieves a download by ID safely
func (d *Downloader) GetDownload(id string) *Download {
	d.mu.RLock()
	defer d.mu.RUnlock()
	if dl, ok := d.downloads[id]; ok {
		return dl
	}
	return nil
}

// AddDownload adds a new download to the map and returns its ID
func (d *Downloader) AddDownload(url string, opts DownloadOptions) *Download {
	d.mu.Lock()
	defer d.mu.Unlock()

	id := fmt.Sprintf("dl_%d", time.Now().UnixNano())
	dl := &Download{
		ID:            id,
		URL:           url,
		Status:        "pending",
		CreatedAt:     time.Now(),
		SubtitleLangs: opts.SubtitleLangs,
		Quality:       opts.Format,
		Options:       opts,
	}
	d.downloads[id] = dl
	return dl
}

// DownloadSynchronously aids the CLI by running the download immediately and blocking
func (d *Downloader) DownloadSynchronously(ctx context.Context, url string, opts DownloadOptions) (*Download, error) {
	dl := d.AddDownload(url, opts)
	dl.Status = "downloading"

	// Delegate to the internal download implementation
	// Note: downloadWithSubtitles updates the dl object directly
	err := d.downloadWithSubtitles(ctx, dl.ID, opts)
	return dl, err
}

// GetProgress returns safe copy of progress fields
func (d *Downloader) GetProgress(id string) (float64, string, string) {
	d.mu.RLock()
	defer d.mu.RUnlock()
	if dl, ok := d.downloads[id]; ok {
		return dl.Progress, dl.ETA, dl.Status
	}
	return 0, "", ""
}

// GetAllDownloads returns all downloads in key-random order
func (d *Downloader) GetAllDownloads() []Download {
	d.mu.RLock()
	defer d.mu.RUnlock()

	list := make([]Download, 0, len(d.downloads))
	for _, dl := range d.downloads {
		list = append(list, *dl)
	}
	return list
}
