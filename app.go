package main

import (
	"context"
	"fmt"
	"log"

	"github.com/shubhambadola/VidFetch/downloader"
)

// App struct
type App struct {
	ctx        context.Context
	downloader *downloader.Downloader
}

// NewApp creates a new App application struct
func NewApp() *App {
	return &App{
		downloader: downloader.NewDownloader(),
	}
}

// startup is called when the app starts. The context is saved
// so we can call the runtime methods
func (a *App) startup(ctx context.Context) {
	a.ctx = ctx
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

	download, err := a.downloader.DownloadSynchronously(a.ctx, url, opts)
	if err != nil {
		return "", err
	}

	return fmt.Sprintf("Download started: %s", download.ID), nil
}

// GetProgress exposed to frontend
func (a *App) GetProgress(id string) (float64, string, string) {
	return a.downloader.GetProgress(id)
}
