package main

import (
	"context"
	"fmt"
	"log"

	"github.com/shubhambadola/VidFetch/downloader"
	"github.com/shubhambadola/VidFetch/storage"
)

// App struct
type App struct {
	ctx        context.Context
	downloader *downloader.Downloader
	history    *storage.History
}

// NewApp creates a new App application struct
func NewApp() *App {
	// Initialize history
	hist, err := storage.NewHistory("history.json")
	if err != nil {
		log.Printf("Failed to load history: %v", err)
		// valid to continue with empty history
		hist, _ = storage.NewHistory("history.json")
	}

	app := &App{
		downloader: downloader.NewDownloader(3), // Max 3 concurrent
		history:    hist,
	}

	// Setup callback
	app.downloader.OnComplete = func(dl *downloader.Download) {
		// Save to history
		if err := app.history.Add(*dl); err != nil {
			log.Printf("Failed to save history: %v", err)
		}
	}

	return app
}

// startup is called when the app starts. The context is saved
// so we can call the runtime methods
func (a *App) startup(ctx context.Context) {
	a.ctx = ctx

	// Start downloader workers
	a.downloader.Start(ctx)

	// Ensure yt-dlp is installed
	go func() {
		if err := downloader.InstallYtDlp(ctx); err != nil {
			log.Printf("Failed to install yt-dlp: %v", err)
		}
	}()
}

// DownloadVideo is the method exposed to the frontend
func (a *App) DownloadVideo(url string) (string, error) {
	// Create default options for now (will expand later)
	opts := downloader.DownloadOptions{
		// Hardcoded output for Phase 2 test
		OutputDir:      "./downloads", // Should probably be absolute path in real app
		OutputTemplate: "%(title)s.%(ext)s",
		DownloadSubs:   true,
		EmbedSubtitles: true,
		SubtitleLangs:  []string{"all"},
		SubtitleFormat: "srt",
	}

	// Use QueueDownload instead of synchronous
	id := a.downloader.QueueDownload(url, opts)
	return fmt.Sprintf("Download queued: %s", id), nil
}

// GetHistory returns completed downloads
func (a *App) GetHistory() []downloader.Download {
	return a.history.Get()
}

// GetQueue returns active/pending downloads
func (a *App) GetQueue() []downloader.Download {
	return a.downloader.GetAllDownloads()
}

// GetProgress exposed to frontend
func (a *App) GetProgress(id string) (float64, string, string) {
	return a.downloader.GetProgress(id)
}
