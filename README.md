# Social Media Tools

A sleek, terminal-based media downloader built with Go and [Bubble Tea](https://github.com/charmbracelet/bubbletea). Download videos, audio, and thumbnails from YouTube (and other platforms supported by yt-dlp) — all from a beautiful TUI.

## Features

- **Clipboard Detection** — Automatically detects URLs copied to your clipboard.
- **Multiple Formats** — Download as video, audio (MP3), or thumbnail.
- **Live Progress** — Real-time progress bar and streaming yt-dlp logs.
- **Smart Download Paths** — Automatically organises downloads by media type and platform (`~/Videos`, `~/Music`, `~/Pictures`), or save to the current directory.
- **Apple-inspired Dark Theme** — Polished terminal UI with rounded borders, gradient progress bar, and styled buttons.

## Prerequisites

| Dependency | Version | Purpose |
|---|---|---|
| [Go](https://go.dev/dl/) | ≥ 1.22 | Build & run |
| [yt-dlp](https://github.com/yt-dlp/yt-dlp) | latest | Media downloading backend |
| [FFmpeg](https://ffmpeg.org/) | latest | Audio extraction & muxing (required by yt-dlp) |

### Installing prerequisites

**Arch Linux:**
```bash
sudo pacman -S go yt-dlp ffmpeg
```

**Ubuntu / Debian:**
```bash
sudo apt install golang ffmpeg
# yt-dlp — install the latest version via pip or download the binary
pip install -U yt-dlp
```

**macOS (Homebrew):**
```bash
brew install go yt-dlp ffmpeg
```

**Windows (Scoop):**
```powershell
scoop install go yt-dlp ffmpeg
```

## Installation

```bash
# Clone the repository
git clone https://github.com/your-username/social-media-tools.git
cd social-media-tools

# Install Go dependencies
go mod download

# Build the binary
go build -o smt .
```

The binary `smt` will be created in the project root. You can move it to a directory in your `$PATH`:

```bash
# Optional: install system-wide
sudo mv smt /usr/local/bin/
```

## Usage

```bash
# Run directly
go run .

# Or use the built binary
./smt
```

### Workflow

1. **Copy a URL** to your clipboard (e.g. a YouTube link) — the app will detect it automatically.
2. If no URL is found, press **`n`** to manually enter one.
3. **Choose a format**: press `1` for Video, `2` for Audio, `3` for Thumbnail.
4. Watch the **live progress bar** and **streaming logs** as yt-dlp downloads your media.
5. When complete, you'll see the download path. Press **`n`** to download another, or **`q`** to quit.

### Keyboard Shortcuts

| Key | Action |
|---|---|
| `n` | Enter / paste a new URL |
| `1` `2` `3` | Select format (Video / Audio / Thumbnail) |
| `c` | Toggle download path (default vs. current directory) |
| `Ctrl+V` | Paste URL in the modal |
| `Enter` | Confirm URL |
| `Esc` | Close the URL modal |
| `q` / `Ctrl+C` | Quit |

### Default Download Locations

| Format | Path |
|---|---|
| Video | `~/Videos/social-media-tools/<Platform>/` |
| Audio | `~/Music/social-media-tools/<Platform>/` |
| Thumbnail | `~/Pictures/social-media-tools/<Platform>/` |

Press **`c`** to toggle saving to the **current working directory** instead.

## Project Structure

```
social-media-tools/
├── main.go                          # App entry point, Bubble Tea model & views
├── go.mod                           # Go module definition & dependencies
├── go.sum                           # Dependency checksums
├── .gitignore
│
├── downloader/
│   └── ytdlp.go                     # yt-dlp process runner — spawns the process,
│                                    # parses stdout/stderr for progress & logs
│
├── internal/
│   ├── clipboard/
│   │   └── clipboard.go             # Clipboard URL detection
│   └── storage/
│       ├── storage.go               # Download path resolution logic
│       └── storage_test.go          # Tests for path generation
│
├── platform/
│   ├── platform.go                  # Platform interface definition
│   └── youtube/
│       ├── youtube.go               # YouTube URL matcher
│       └── youtube_test.go          # Tests for URL matching
│
└── ui/
    └── styles/
        └── styles.go                # Lipgloss styles (colors, borders, buttons)
```

### Architecture

```
┌─────────────┐     ┌──────────────┐     ┌──────────────┐
│   main.go   │────▶│  downloader  │────▶│   yt-dlp     │
│  (Bubble Tea│◀────│  (goroutine) │◀────│  (process)   │
│   TUI)      │ chan │              │ pipe │              │
└──────┬──────┘     └──────────────┘     └──────────────┘
       │
       ├── internal/clipboard   (URL auto-detection)
       ├── internal/storage     (download path logic)
       ├── platform/*           (URL → platform matching)
       └── ui/styles            (theme & components)
```

- **main.go** — The Bubble Tea application. Manages UI state (detecting → selecting → downloading → finished), handles keyboard input, and renders the TUI.
- **downloader/ytdlp.go** — Spawns `yt-dlp` as a subprocess, captures both stdout (progress percentages) and stderr (info/warnings), and streams `ProgressMsg` structs over a channel back to the TUI.
- **internal/clipboard** — Reads the system clipboard on startup and checks for HTTP(S) URLs.
- **internal/storage** — Resolves the output directory based on media type, platform name, and whether the user toggled "save to PWD".
- **platform/** — Extensible platform matching. Currently supports YouTube. To add a new platform, implement the `Platform` interface and register it in `main.go`.
- **ui/styles** — Lipgloss style definitions for the Apple-inspired dark theme.

## Dependencies

| Package | Purpose |
|---|---|
| [`charmbracelet/bubbletea`](https://github.com/charmbracelet/bubbletea) | Terminal UI framework (Elm-architecture) |
| [`charmbracelet/bubbles`](https://github.com/charmbracelet/bubbles) | Pre-built TUI components (progress bar, text input) |
| [`charmbracelet/lipgloss`](https://github.com/charmbracelet/lipgloss) | Terminal styling & layout |
| [`atotto/clipboard`](https://github.com/atotto/clipboard) | Cross-platform clipboard access |

## Development

```bash
# Run with live reload (using air, optional)
go install github.com/air-verse/air@latest
air

# Run tests
go test ./...

# Build for release
go build -ldflags="-s -w" -o smt .
```

### Adding a New Platform

1. Create a new package under `platform/` (e.g. `platform/instagram/`).
2. Implement the `Platform` interface:
   ```go
   type Platform interface {
       Match(url string) bool
       GetName() string
   }
   ```
3. Register it in `main.go`:
   ```go
   platforms: []platform.Platform{
       &youtube.YouTube{},
       &instagram.Instagram{},  // Add here
   },
   ```

## License

This project is provided as-is. See the repository for license details.
