package downloader

import (
	"fmt"

	"github.com/lrstanley/go-ytdlp"
)

// handleProgress updates the download state based on the callback from yt-dlp
func (d *Downloader) handleProgress(id string, update ytdlp.ProgressUpdate) {
	d.mu.Lock()
	defer d.mu.Unlock()

	dl, ok := d.downloads[id]
	if !ok {
		return
	}

	dl.Progress = update.Percent()
	dl.Downloaded = int64(update.DownloadedBytes)
	dl.FileSize = int64(update.TotalBytes)

	if update.Status != "" {
		// Map ytdlp status to our status if useful,
		// but usually "downloading" is fine until done.
		// update.Status might be "downloading", "finished", "processing"
	}

	// Format ETA
	eta := update.ETA()
	if eta > 0 {
		dl.ETA = eta.String()
	}

	// Speed isn't directly compatible or calculated?
	// We can leave speed empty for now or calculate it if needed.
	// For simplicity in Phase 1, just bytes/percent/ETA.
}

func formatBytes(bytes int64) string {
	const unit = 1024
	if bytes < unit {
		return fmt.Sprintf("%d B", bytes)
	}
	div, exp := int64(unit), 0
	for n := bytes / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}
	return fmt.Sprintf("%.1f %ciB", float64(bytes)/float64(div), "KMGTPE"[exp])
}
