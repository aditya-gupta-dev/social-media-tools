package downloader

import (
	"bufio"
	"os/exec"
	"regexp"
	"strconv"
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
		progress <- ProgressMsg{Err: err}
		return
	}

	if err := cmd.Start(); err != nil {
		progress <- ProgressMsg{Err: err}
		return
	}

	scanner := bufio.NewScanner(stdout)
	re := regexp.MustCompile(`\[download\]\s+(\d+\.\d+)%`)

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

	if err := cmd.Wait(); err != nil {
		progress <- ProgressMsg{Err: err}
	} else {
		progress <- ProgressMsg{Done: true}
	}
}
