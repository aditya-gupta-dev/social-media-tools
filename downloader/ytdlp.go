package downloader

import (
	"bufio"
	"fmt"
	"os/exec"
	"regexp"
	"strconv"
	"strings"
	"encoding/json"
	"os"
	"path/filepath"


	yttranscript "github.com/aditya-gupta-dev/go-yt-transcript"
)

type ProgressMsg struct {
	Percent float64
	Log     string
	Done    bool
	Err     error
}

func Download(url, mediaType, outputPath string, progress chan ProgressMsg) {
	if mediaType == "transcript" {
		progress <- ProgressMsg{Log: "Fetching transcript metadata..."}
		snippets, err := yttranscript.GetTranscript(url)
		if err != nil {
			progress <- ProgressMsg{Err: fmt.Errorf("failed to fetch transcript: %w", err)}
			return
		}

		progress <- ProgressMsg{Log: "Parsing transcript...", Percent: 0.5}
		
		id, _ := yttranscript.ExtractVideoID(url)
		
		// Create a directory for this specific video ID
		videoDir := filepath.Join(outputPath, id)
		if err := os.MkdirAll(videoDir, 0755); err != nil {
			progress <- ProgressMsg{Err: fmt.Errorf("failed to create video directory: %w", err)}
			return
		}
		
		filename := filepath.Join(videoDir, "transcript.json")
		
		file, err := os.Create(filename)
		if err != nil {
			progress <- ProgressMsg{Err: fmt.Errorf("failed to create file: %w", err)}
			return
		}
		defer file.Close()

		progress <- ProgressMsg{Log: "Saving to file...", Percent: 0.8}
		
		encoder := json.NewEncoder(file)
		encoder.SetIndent("", "  ")
		if err := encoder.Encode(snippets); err != nil {
			progress <- ProgressMsg{Err: fmt.Errorf("failed to write json: %w", err)}
			return
		}

		progress <- ProgressMsg{Log: "Done!", Done: true, Percent: 1.0}
		return
	}

	var args []string
	args = append(args, "--newline", "--progress")
	args = append(args, "-o", outputPath+"/%(title)s.%(ext)s")

	switch mediaType {
	case "audio":
		args = append(args, "-x", "--audio-format", "mp3")
	case "video":
		args = append(args, "-f", "bestvideo+bestaudio/best")
	case "thumbnail":
		args = append(args, "--write-thumbnail", "--skip-download")
	}

	args = append(args, url)

	cmd := exec.Command("yt-dlp", args...)
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		progress <- ProgressMsg{Err: fmt.Errorf("stdout pipe: %w", err)}
		return
	}

	stderr, err := cmd.StderrPipe()
	if err != nil {
		progress <- ProgressMsg{Err: fmt.Errorf("stderr pipe: %w", err)}
		return
	}

	if err := cmd.Start(); err != nil {
		progress <- ProgressMsg{Err: fmt.Errorf("failed to start yt-dlp: %w", err)}
		return
	}

	// Matches both "45.2%" and "100%"
	re := regexp.MustCompile(`\[download\]\s+(\d+(?:\.\d+)?)%`)

	// Read stderr in a separate goroutine so it doesn't block
	stderrDone := make(chan struct{})
	go func() {
		defer close(stderrDone)
		sc := bufio.NewScanner(stderr)
		for sc.Scan() {
			line := strings.TrimSpace(sc.Text())
			if line != "" {
				progress <- ProgressMsg{Log: line}
			}
		}
	}()

	// Read stdout (progress lines) in the main goroutine
	scanner := bufio.NewScanner(stdout)
	for scanner.Scan() {
		line := scanner.Text()
		msg := ProgressMsg{Log: line}

		matches := re.FindStringSubmatch(line)
		if len(matches) > 1 {
			p, _ := strconv.ParseFloat(matches[1], 64)
			msg.Percent = p / 100.0
		}

		progress <- msg
	}

	// Wait for stderr goroutine to finish before calling cmd.Wait
	<-stderrDone

	if err := cmd.Wait(); err != nil {
		progress <- ProgressMsg{Err: fmt.Errorf("yt-dlp exited with error: %w", err)}
	} else {
		progress <- ProgressMsg{Done: true, Percent: 1.0}
	}
}
