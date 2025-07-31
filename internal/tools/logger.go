package tools

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"os/exec"
	"regexp"
	"strconv"
	"strings"
)

const (
	colorReset = "\033[0m"
	colorGreen = "\033[32m"
	colorCyan  = "\033[36m"
)

// ProgressLogger holds the necessary data for logging ffmpeg progress.
type ProgressLogger struct {
	duration float64
	regex    *regexp.Regexp
}

// NewProgressLogger creates a new ProgressLogger.
func NewProgressLogger(videoPath string) (*ProgressLogger, error) {
	duration, err := getVideoDuration(videoPath)
	if err != nil {
		return nil, fmt.Errorf("failed to get video duration: %w", err)
	}

	return &ProgressLogger{
		duration: duration,
		regex:    regexp.MustCompile(`time=(\d{2}):(\d{2}):(\d{2})\.(\d{2})`),
	}, nil
}

// LogProgress reads from the ffmpeg stderr and logs the progress.
func (p *ProgressLogger) LogProgress(stderr io.Reader) {
	scanner := bufio.NewScanner(stderr)
	for scanner.Scan() {
		line := scanner.Text()

		if strings.Contains(line, "frame=") {
			matches := p.regex.FindStringSubmatch(line)
			if len(matches) >= 5 {
				hours, _ := strconv.Atoi(matches[1])
				minutes, _ := strconv.Atoi(matches[2])
				seconds, _ := strconv.Atoi(matches[3])
				currentTime := float64(hours*3600 + minutes*60 + seconds)
				percentage := (currentTime / p.duration) * 100
				log.Printf(colorGreen+"Transcoding progress: %.2f%%"+colorReset, percentage)
			}
		}

		if strings.Contains(line, "Opening '") && strings.Contains(line, ".ts' for writing") {
			re := regexp.MustCompile(`Opening '([^']*)' for writing`)
			matches := re.FindStringSubmatch(line)
			if len(matches) > 1 {
				log.Printf(colorCyan+"Processing chunk: %s"+colorReset, matches[1])
			}
		}
	}
}

// getVideoDuration uses ffprobe to get the duration of a video file.
func getVideoDuration(videoPath string) (float64, error) {
	cmd := exec.Command("ffprobe", "-v", "error", "-show_entries", "format=duration", "-of", "default=noprint_wrappers=1:nokey=1", videoPath)
	output, err := cmd.Output()
	if err != nil {
		return 0, err
	}

	duration, err := strconv.ParseFloat(strings.TrimSpace(string(output)), 64)
	if err != nil {
		return 0, err
	}

	return duration, nil
}