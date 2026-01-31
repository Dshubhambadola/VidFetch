package downloader

import (
	"context"
	"fmt"
	"path/filepath"
	"strings"
	"time"

	"github.com/lrstanley/go-ytdlp"
)

// Ensure ytdlp is installed (could be called on startup)
func InstallYtDlp(ctx context.Context) error {
	_, err := ytdlp.Install(ctx, nil)
	return err
}

// downloadWithSubtitles executes the download using yt-dlp
func (d *Downloader) downloadWithSubtitles(ctx context.Context, id string, opts DownloadOptions) error {
	dl := d.GetDownload(id)
	if dl == nil {
		return fmt.Errorf("download not found: %s", id)
	}

	// Construct yt-dlp arguments
	cmd := ytdlp.New().
		Format("bestvideo+bestaudio/best"). // Default best quality
		Output(filepath.Join(opts.OutputDir, opts.OutputTemplate)).
		NoOverwrites().
		// Progress handling
		Progress().
		ProgressFunc(500*time.Millisecond, func(p ytdlp.ProgressUpdate) {
			d.handleProgress(id, p)
		}).
		// Subtitle Options
		SubLangs(strings.Join(opts.SubtitleLangs, ",")).
		ConvertSubs(opts.SubtitleFormat).
		// Error handling
		NoAbortOnError().
		IgnoreErrors().
		// Local config
		NoCheckCertificates()

	if opts.VideoFormat != "" {
		cmd = cmd.MergeOutputFormat(opts.VideoFormat)
	}

	if opts.DownloadSubs {
		cmd = cmd.WriteSubs()
	}
	if opts.DownloadAutoSubs {
		cmd = cmd.WriteAutoSubs()
	}
	if opts.EmbedSubtitles {
		cmd = cmd.EmbedSubs()
	}
	// If NOT saving external files separately AND embedding is on, we might need flags.
	// But usually WriteSubs + EmbedSubs keeps the external file unless we delete it?
	// actually yt-dlp might keep it.
	// We won't add complex cleanup logic yet unless requested.

	// Run command
	_, err := cmd.Run(ctx, dl.URL)
	if err != nil {
		d.mu.Lock()
		dl.Status = "failed"
		dl.Error = err.Error()
		d.mu.Unlock()
		return err
	}

	d.mu.Lock()
	dl.Status = "completed"
	dl.Progress = 100
	d.mu.Unlock()

	return nil
}
