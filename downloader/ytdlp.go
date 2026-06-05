package downloader

import (
	"bufio"
	"fmt"
	"os/exec"
	"regexp"
	"strconv"
	"strings"
)

type ProgressMsg struct {
	Percent float64
	Log     string
	Done    bool
	Err     error
}

func Download(url, mediaType, outputPath string, progress chan ProgressMsg) {
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
