package yt

import (
	"context"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"time"

	"github.com/kkdai/youtube/v2"
)

type DownloadOptions struct {
	URL string

	OutputDir  string
	OutputName string
	Timeout    time.Duration
}

type Result struct {
	VideoID    string
	Title      string
	Author     string
	FilePath   string
	FileName   string
	MimeType   string
	Quality    string
	Itag       int
	ContentLen int64
}

func Download(opts DownloadOptions) (*Result, error) {
	if strings.TrimSpace(opts.URL) == "" {
		return nil, errors.New("ytdl: missing URL")
	}

	ctx := context.Background()
	if opts.Timeout > 0 {
		var cancel context.CancelFunc
		ctx, cancel = context.WithTimeout(ctx, opts.Timeout)
		defer cancel()
	}

	client := &youtube.Client{}

	video, err := client.GetVideoContext(ctx, opts.URL)
	if err != nil {
		return nil, fmt.Errorf("ytdl: get video: %w", err)
	}

	format := bestProgressive(video.Formats)
	if format == nil {
		return nil, errors.New("ytdl: no progressive stream found")
	}

	stream, size, err := client.GetStreamContext(ctx, video, format)
	if err != nil {
		return nil, fmt.Errorf("ytdl: get stream: %w", err)
	}
	defer stream.Close()

	dir := opts.OutputDir
	if dir == "" {
		dir = "."
	}

	if err := os.MkdirAll(dir, 0o755); err != nil {
		return nil, fmt.Errorf("ytdl: create output dir: %w", err)
	}

	name := opts.OutputName
	if name == "" {
		ext := extensionFromMime(format.MimeType)
		if ext == "" {
			ext = ".mp4"
		}
		name = slug(video.Title) + ext
	}

	path := filepath.Join(dir, name)

	out, err := os.Create(path)
	if err != nil {
		return nil, fmt.Errorf("ytdl: create file: %w", err)
	}
	defer out.Close()

	if _, err := io.Copy(out, stream); err != nil {
		return nil, fmt.Errorf("ytdl: write file: %w", err)
	}

	return &Result{
		VideoID:    video.ID,
		Title:      video.Title,
		Author:     video.Author,
		FilePath:   path,
		FileName:   name,
		MimeType:   format.MimeType,
		Quality:    format.Quality,
		Itag:       format.ItagNo,
		ContentLen: size,
	}, nil
}

func bestProgressive(formats youtube.FormatList) *youtube.Format {
	var best *youtube.Format
	bestScore := -1

	for i := range formats {
		f := &formats[i]

		// Progressive means audio + video together.
		if f.AudioChannels <= 0 {
			continue
		}

		score := scoreFormat(f)
		if score > bestScore {
			bestScore = score
			best = f
		}
	}

	return best
}

func scoreFormat(f *youtube.Format) int {
	score := 0

	score += parseQuality(f.Quality) * 1000

	if strings.Contains(strings.ToLower(f.MimeType), "video/mp4") {
		score += 100
	}

	score += f.Bitrate / 1000
	score += f.AudioChannels * 10

	return score
}

var qualityRE = regexp.MustCompile(`(\d+)`)

func parseQuality(q string) int {
	m := qualityRE.FindStringSubmatch(strings.ToLower(q))
	if len(m) < 2 {
		switch strings.ToLower(q) {
		case "small":
			return 240
		case "medium":
			return 360
		case "large":
			return 480
		case "hd720":
			return 720
		case "hd1080":
			return 1080
		default:
			return 0
		}
	}

	var n int
	fmt.Sscanf(m[1], "%d", &n)
	return n
}

func extensionFromMime(mime string) string {
	m := strings.ToLower(mime)
	switch {
	case strings.Contains(m, "mp4"):
		return ".mp4"
	case strings.Contains(m, "webm"):
		return ".webm"
	case strings.Contains(m, "3gpp"):
		return ".3gp"
	default:
		return ""
	}
}

func slug(s string) string {
	s = strings.ToLower(s)

	var b strings.Builder
	lastDash := false

	for _, r := range s {

		switch {
		case r >= 'a' && r <= 'z':
			b.WriteRune(r)
			lastDash = false

		case r >= '0' && r <= '9':
			b.WriteRune(r)
			lastDash = false

		default:
			if !lastDash {
				b.WriteRune('-')
				lastDash = true
			}
		}
	}

	out := strings.Trim(b.String(), "-")

	if out == "" {
		return "video"
	}

	return out
}
