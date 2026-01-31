package main

import (
	"context"
	"fmt"
	"log"

	"github.com/shubhambadola/VidFetch/downloader"
	"github.com/shubhambadola/VidFetch/storage"
	"github.com/wailsapp/wails/v2/pkg/runtime"
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

		// Emit event to frontend if context is available
		// Only emit success event if actually completed
		if app.ctx != nil && dl.Status == "completed" {
			runtime.EventsEmit(app.ctx, "download-complete", dl)
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
	// Ensure yt-dlp is installed and get path
	go func() {
		path, err := downloader.InstallYtDlp(ctx)
		if err != nil {
			log.Printf("Failed to install yt-dlp: %v", err)
		} else {
			// Update the downloader with the cached path
			a.downloader.BinPath = path
			a.downloader.Updater = downloader.NewUpdater(path)
			log.Printf("yt-dlp ready at: %s", path)

			// Auto-check for updates on startup (async)
			go func() {
				msg, err := a.downloader.Updater.CheckAndUpdate(ctx, "stable")
				if err != nil {
					log.Printf("Startup update check warning: %v", err)
				} else {
					log.Printf("Startup update check: %s", msg)
				}
			}()
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

// DownloadVideoWithOptions allows frontend to specify options (Quality, etc.)
func (a *App) DownloadVideoWithOptions(url string, options downloader.DownloadOptions) (string, error) {
	// Set defaults for missing fields if necessary
	if options.OutputDir == "" {
		options.OutputDir = "./downloads"
	}
	if options.OutputTemplate == "" {
		options.OutputTemplate = "%(title)s.%(ext)s"
	}
	// Ensure subtitle settings are preserved if not explicitly set?
	// The frontend should probably send full state, but for safety:
	if len(options.SubtitleLangs) == 0 {
		options.SubtitleLangs = []string{"all"}
	}

	id := a.downloader.QueueDownload(url, options)
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

// CheckForUpdates triggers a manual update check
func (a *App) CheckForUpdates(mode string) (string, error) {
	if a.downloader.Updater == nil {
		return "", fmt.Errorf("engine initializing, please wait...")
	}
	return a.downloader.Updater.CheckAndUpdate(a.ctx, mode)
}

// GetYtdlpVersion returns the current version
func (a *App) GetYtdlpVersion() (string, error) {
	if a.downloader.Updater == nil {
		return "Initializing...", nil
	}
	return a.downloader.Updater.GetVersion(a.ctx)
}
