package downloader

import (
	"bufio"
	"context"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/lrstanley/go-ytdlp"
)

// Ensure ytdlp is installed (could be called on startup)
func InstallYtDlp(ctx context.Context) (string, error) {
	// Try strict install first
	path, err := ytdlp.Install(ctx, nil)
	if err == nil {
		// Update check (optional, don't fail if this fails)
		go func() {
			cmd := exec.CommandContext(context.Background(), path.Executable, "-U")
			cmd.Run()
		}()
		return path.Executable, nil
	}

	// Fallback: Check local cache for existing binary
	// Only do this if Install failed (likely network timeout)
	cacheDir, _ := os.UserCacheDir()
	ytdlpCache := filepath.Join(cacheDir, "go-ytdlp")

	// Find any file starting with "yt-dlp-"
	entries, _ := os.ReadDir(ytdlpCache)
	for _, entry := range entries {
		name := entry.Name()
		if strings.HasPrefix(name, "yt-dlp-") && !strings.HasSuffix(name, ".tmp") {
			fullPath := filepath.Join(ytdlpCache, name)
			fmt.Printf("Recovered using cached yt-dlp: %s\n", fullPath)
			return fullPath, nil
		}
	}

	return "", fmt.Errorf("failed to install yt-dlp and no local binary found: %v", err)
}

// downloadWithSubtitles executes the download using yt-dlp
func (d *Downloader) downloadWithSubtitles(ctx context.Context, id string, opts DownloadOptions) error {
	dl := d.GetDownload(id)
	if dl == nil {
		return fmt.Errorf("download not found: %s", id)
	}

	// Determine format
	format := "bestvideo+bestaudio/best"
	if opts.Format != "" {
		format = opts.Format
	}

	// Locate binary
	binPath := d.BinPath
	if binPath == "" {
		// Fallback (should not happen if initialized correctly)
		var err error
		binPath, err = InstallYtDlp(ctx)
		if err != nil {
			return err
		}
		d.BinPath = binPath // cache it
	}

	// Prepare args
	var args []string

	// Format
	args = append(args, "--format", format)

	// Paths
	output := filepath.Join(opts.OutputDir, opts.OutputTemplate)
	args = append(args, "--output", output)
	args = append(args, "--no-overwrites")

	// Networking / Anti-Bot
	// 1. Cookies (Best method)
	if opts.UseCookies {
		browser := opts.BrowserName
		if browser == "" {
			browser = "chrome" // default
		}
		args = append(args, "--cookies-from-browser", browser)
	}

	// 2. User Agent
	if opts.UserAgent != "" {
		args = append(args, "--user-agent", opts.UserAgent)
	} else {
		// Default robust UA
		args = append(args, "--user-agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.0.0 Safari/537.36")
	}

	// 3. Proxy
	if opts.ProxyURL != "" {
		args = append(args, "--proxy", opts.ProxyURL)
	}

	// 4. Rate Limit
	if opts.RateLimit != "" {
		args = append(args, "--limit-rate", opts.RateLimit)
	}

	// Common headers
	// Mimic browser aggressively
	args = append(args, "--add-header", "Accept-Language:en-US,en;q=0.9")
	args = append(args, "--add-header", "Accept:text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.7")
	args = append(args, "--add-header", "DNT:1")
	args = append(args, "--add-header", "Sec-Fetch-Mode:navigate")
	args = append(args, "--referer", dl.URL)
	args = append(args, "--no-check-certificates")

	// Subtitles
	if len(opts.SubtitleLangs) > 0 {
		args = append(args, "--sub-langs", strings.Join(opts.SubtitleLangs, ","))
	} else {
		args = append(args, "--sub-langs", "all")
	}

	if opts.SubtitleFormat != "" {
		args = append(args, "--convert-subs", opts.SubtitleFormat)
	}

	if opts.DownloadSubs {
		args = append(args, "--write-subs")
	}
	if opts.DownloadAutoSubs {
		args = append(args, "--write-auto-subs")
	}
	if opts.EmbedSubtitles {
		args = append(args, "--embed-subs")
	}

	args = append(args, "--no-abort-on-error")
	args = append(args, "--ignore-errors")

	// Experimental Impersonation (for TLS Fingerprinting)
	if opts.Impersonate != "" {
		args = append(args, "--impersonate", opts.Impersonate)
	}

	args = append(args, "--newline") // Critical for parsing
	args = append(args, "--progress")

	// URL must be last
	args = append(args, dl.URL)

	// Execute
	d.mu.Lock()
	dl.Status = "downloading" // ensure status
	d.mu.Unlock()

	cmd := exec.CommandContext(ctx, binPath, args...)

	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return err
	}
	// Also capture stderr using same pipe or separate?
	// yt-dlp prints progress to stdout, errors to stderr.
	// combinedOutput is easiest but we need streaming.
	// We will capture stderr separately for error reporting.
	cmd.Stderr = cmd.Stdout // Merge them for simplicity in this context

	if err := cmd.Start(); err != nil {
		d.mu.Lock()
		dl.Status = "failed"
		dl.Error = err.Error()
		d.mu.Unlock()
		return err
	}

	// Parser regex items
	// [download]  23.5% of 10.00MiB at 2.00MiB/s ETA 00:05
	reProgress := regexp.MustCompile(`\[download\]\s+(\d+\.?\d*)%`)
	reETA := regexp.MustCompile(`ETA\s+(\d+:\d+)`)
	reSpeed := regexp.MustCompile(`at\s+(\d+\.?\d*\w+/s)`)

	// Scan output
	// Scan output
	var outputLog strings.Builder
	scanner := bufio.NewScanner(stdout)
	for scanner.Scan() {
		line := scanner.Text()
		outputLog.WriteString(line + "\n")

		// Parse progress
		if strings.Contains(line, "[download]") {
			// Extract percent
			if matches := reProgress.FindStringSubmatch(line); len(matches) > 1 {
				if p, err := strconv.ParseFloat(matches[1], 64); err == nil {
					d.mu.Lock()
					dl.Progress = p / 100.0

					// Extract ETA
					if etaMatch := reETA.FindStringSubmatch(line); len(etaMatch) > 1 {
						dl.ETA = etaMatch[1]
					}
					// Extract Speed
					if speedMatch := reSpeed.FindStringSubmatch(line); len(speedMatch) > 1 {
						dl.Speed = speedMatch[1]
					}

					d.mu.Unlock()
				}
			}
		}
	}

	if err := cmd.Wait(); err != nil {
		d.mu.Lock()
		dl.Status = "failed"
		dl.Error = outputLog.String() // Provide full output as error
		d.mu.Unlock()
		return fmt.Errorf("yt-dlp error: %v, out: %s", err, outputLog.String())
	}

	d.mu.Lock()
	dl.Status = "completed"
	dl.Progress = 1.0
	dl.CompletedAt = time.Now()
	d.mu.Unlock()

	return nil
}
