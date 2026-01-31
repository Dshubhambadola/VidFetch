package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"time"

	"github.com/shubhambadola/VidFetch/downloader"
)

func main() {
	urlFlag := flag.String("url", "", "URL to download")
	outputDirFlag := flag.String("out", "./downloads", "Output directory")
	subsFlag := flag.Bool("subs", true, "Download subtitles")
	embedFlag := flag.Bool("embed", true, "Embed subtitles")

	// Anti-Blocking Flags
	cookies := flag.Bool("cookies", false, "Use Chrome browser cookies")
	proxy := flag.String("proxy", "", "Proxy URL")
	rateLimit := flag.String("limit", "", "Rate limit (e.g. 2M)")
	userAgent := flag.String("ua", "", "Custom User Agent")

	flag.Parse()

	if *urlFlag == "" {
		fmt.Println("Please provide a URL using -url")
		flag.PrintDefaults()
		os.Exit(1)
	}

	// Create output dir
	absPath, _ := filepath.Abs(*outputDirFlag)
	if err := os.MkdirAll(absPath, 0755); err != nil {
		log.Fatalf("Failed to create output directory: %v", err)
	}

	fmt.Printf("Initializing VidFetch Core...\n")
	ctx := context.Background()

	// Ensure yt-dlp is installed
	binPath, err := downloader.InstallYtDlp(ctx)
	if err != nil {
		log.Printf("Installing yt-dlp failed (might already be installed or network issue): %v", err)
	}

	dlr := downloader.NewDownloader(1)
	dlr.BinPath = binPath

	fmt.Printf("Starting download for: %s\n", *urlFlag)
	fmt.Printf("Output directory: %s\n", absPath)

	opts := downloader.DownloadOptions{
		OutputDir:        absPath,
		OutputTemplate:   "%(title)s.%(ext)s",
		DownloadSubs:     *subsFlag,
		DownloadAutoSubs: *subsFlag,
		EmbedSubtitles:   *embedFlag,
		SubtitleLangs:    []string{"all"}, // Default to all
		SubtitleFormat:   "srt",

		// Anti-Blocking
		UseCookies:  *cookies,
		BrowserName: "chrome",
		ProxyURL:    *proxy,
		RateLimit:   *rateLimit,
		UserAgent:   *userAgent,
	}

	start := time.Now()

	// Start download in goroutine
	done := make(chan error, 1)
	dlReady := make(chan string, 1)

	go func() {
		dl, err := dlr.DownloadSynchronously(ctx, *urlFlag, opts)
		if dl != nil {
			dlReady <- dl.ID
		}
		done <- err
	}()

	var currentID string
	// Wait for ID or completion
	select {
	case id := <-dlReady:
		currentID = id
	case err := <-done:
		if err != nil {
			log.Fatalf("Download failed: %v", err)
		}
		fmt.Printf("\nDownload completed successfully in %v\n", time.Since(start))
		return
	}

	ticker := time.NewTicker(500 * time.Millisecond)
	defer ticker.Stop()

	for {
		select {
		case err := <-done:
			if err != nil {
				log.Fatalf("Download failed: %v", err)
			}
			fmt.Printf("\nDownload completed successfully in %v\n", time.Since(start))
			return
		case <-ticker.C:
			if currentID != "" {
				prog, eta, status := dlr.GetProgress(currentID)
				fmt.Printf("\rProgress: %.1f%% | ETA: %s | Status: %s   ", prog*100, eta, status)
			}
		}
	}
}
