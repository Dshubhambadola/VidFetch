# VidFetch 🚀

VidFetch is a modern, high-performance video downloader built with **Go (Wails)** and **React**. It leverages the power of `yt-dlp` to download high-quality videos from YouTube, Instagram, TikTok, and thousands of other sites, wrapped in a beautiful desktop UI.

![VidFetch UI](./frontend/public/app-icon.png)

## ✨ Features

- **🎥 High-Quality Downloads**: Support for 4K/8K video, audio-only extraction, and various formats (MP4, MKV, WEBM).
- **🚀 Turbo Engine**: Powered by `yt-dlp` for maximum speed and compatibility.
- **🛡️ Anti-Blocking System**:
  - **Browser Cookies**: Import cookies from Chrome/Firefox to access age-gated/premium content.
  - **TLS Impersonation**: Mimic real browsers (Chrome, Safari) to bypass strict firewalls (Cloudflare).
  - **Smart Headers**: Automatic injection of browser-like headers.
- **🔄 Auto-Update System**:
  - Automatically keeps the download engine (`yt-dlp`) up to date.
  - Switch between **Stable** and **Nightly** builds (Nightly recommended for YouTube).
- **📋 Queue & History**: Manage multiple downloads with accurate progress tracking, speed stats, and a history log.
- **🌗 Beautiful UI**: Clean, responsive interface with Dark/Light mode support.
- **📦 Batch Download**: Queue multiple URLs at once.

## 🛠️ Tech Stack

- **Backend**: Go (Golang)
- **Frontend**: React, TypeScript, TailwindCSS
- **Framework**: Wails v2 (Go + Webview)
- **Engine**: `yt-dlp` (embedded/managed)

## 📦 Installation & Build

### Prerequisites
- **Go** (1.20+)
- **Node.js** (18+)
- **yt-dlp**: (Optional) The app will attempt to install it automatically, but having it installed is recommended.

### Development
1. Clone the repo:
   ```bash
   git clone https://github.com/yourusername/VidFetch.git
   cd VidFetch
   ```
2. Install Wails:
   ```bash
   go install github.com/wailsapp/wails/v2/cmd/wails@latest
   ```
3. Run in Dev Mode:
   ```bash
   wails dev
   ```

### Production Build
To create a standalone application (`.app` on macOS, `.exe` on Windows):

```bash
wails build
```

The output will be in `build/bin/`.

## 📖 User Guide

### Basic Download
1. Paste a video URL in the input box.
2. Select desired Quality (e.g., 1080p).
3. Click **Download**.

### Anti-Blocking (Fixing Errors)
If you encounter "Video Unavailable" or "Connection Reset" errors:
1. Open **Settings** (Gear icon).
2. Enable **Use Browser Cookies** (Select your browser).
3. If that fails, try **TLS Impersonation**: Select "Chrome" or "Safari".
4. Check **Engine Updates** and update to the **Nightly** build.

## 🤝 Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

## 📄 License

MIT License
