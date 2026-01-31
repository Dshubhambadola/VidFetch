package downloader

import (
	"context"
	"fmt"
	"os/exec"
	"strings"
	"time"
)

type Updater struct {
	BinPath         string
	lastUpdateCheck time.Time
	updateInterval  time.Duration
}

// NewUpdater creates a new updater instance
func NewUpdater(binPath string) *Updater {
	return &Updater{
		BinPath:        binPath,
		updateInterval: 24 * time.Hour,
	}
}

// CheckUpdate checks if an update is available and updates if needed
// mode: "stable" or "nightly"
func (u *Updater) CheckAndUpdate(ctx context.Context, mode string) (string, error) {
	// Skip if checked recently (debounce)
	if time.Since(u.lastUpdateCheck) < u.updateInterval {
		return "Skipped (recently checked)", nil
	}

	fmt.Println("Checking for yt-dlp updates...")
	u.lastUpdateCheck = time.Now()

	var cmd *exec.Cmd
	if mode == "nightly" {
		// Update to nightly build
		cmd = exec.CommandContext(ctx, u.BinPath, "--update-to", "nightly")
	} else {
		// Update to stable
		cmd = exec.CommandContext(ctx, u.BinPath, "-U")
	}

	output, err := cmd.CombinedOutput()
	if err != nil {
		return "", fmt.Errorf("update failed: %v, output: %s", err, string(output))
	}

	outStr := string(output)
	if strings.Contains(outStr, "up-to-date") {
		return "Up to date", nil
	}

	return fmt.Sprintf("Updated: %s", outStr), nil
}

// GetVersion returns the current version of yt-dlp
func (u *Updater) GetVersion(ctx context.Context) (string, error) {
	if u.BinPath == "" {
		return "", fmt.Errorf("binary path not set")
	}

	cmd := exec.CommandContext(ctx, u.BinPath, "--version")
	output, err := cmd.Output()
	if err != nil {
		return "", err
	}

	return strings.TrimSpace(string(output)), nil
}
