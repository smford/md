package image

import (
	"bytes"
	"context"
	"encoding/base64"
	"errors"
	"fmt"
	"image"
	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"
	"io"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strings"
	"time"

	_ "golang.org/x/image/webp"
)

const (
	// MaxImageSize limits downloaded image payload to 25MB to prevent memory exhaustion.
	MaxImageSize = 25 * 1024 * 1024
	// DefaultTimeout defines HTTP fetch timeout for remote images.
	DefaultTimeout = 10 * time.Second
)

// Info contains metadata about an image.
type Info struct {
	Width    int
	Height   int
	Format   string
	ByteSize int
}

// Fetch loads an image from a local file, HTTP/HTTPS URL, or base64 data URI.
// If src is a relative path, it will be resolved relative to basePath.
func Fetch(ctx context.Context, src string, basePath string, client *http.Client) ([]byte, string, error) {
	src = strings.TrimSpace(src)
	if src == "" {
		return nil, "", errors.New("empty image source")
	}

	// 1. Data URI: data:image/png;base64,...
	if strings.HasPrefix(src, "data:") {
		return parseDataURI(src)
	}

	// 2. HTTP/HTTPS URL
	if strings.HasPrefix(src, "http://") || strings.HasPrefix(src, "https://") {
		return fetchRemote(ctx, src, client)
	}

	// 3. Local filesystem
	return fetchLocal(src, basePath)
}

func parseDataURI(dataURI string) ([]byte, string, error) {
	commaIdx := strings.Index(dataURI, ",")
	if commaIdx == -1 {
		return nil, "", errors.New("invalid data URI format")
	}

	meta := dataURI[:commaIdx]
	dataStr := dataURI[commaIdx+1:]

	var filename string
	if strings.Contains(meta, "image/png") {
		filename = "image.png"
	} else if strings.Contains(meta, "image/jpeg") {
		filename = "image.jpg"
	} else if strings.Contains(meta, "image/gif") {
		filename = "image.gif"
	} else if strings.Contains(meta, "image/webp") {
		filename = "image.webp"
	} else if strings.Contains(meta, "image/svg") {
		filename = "image.svg"
	} else {
		filename = "image.bin"
	}

	if strings.Contains(meta, ";base64") {
		data, err := base64.StdEncoding.DecodeString(dataStr)
		if err != nil {
			return nil, "", fmt.Errorf("failed to decode base64 data URI: %w", err)
		}
		return data, filename, nil
	}

	// URL-encoded or raw data
	unescaped, err := url.QueryUnescape(dataStr)
	if err != nil {
		return nil, "", fmt.Errorf("failed to unescape data URI: %w", err)
	}
	return []byte(unescaped), filename, nil
}

func fetchRemote(ctx context.Context, rawURL string, client *http.Client) ([]byte, string, error) {
	if client == nil {
		client = &http.Client{Timeout: DefaultTimeout}
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, rawURL, nil)
	if err != nil {
		return nil, "", fmt.Errorf("invalid request: %w", err)
	}

	// SRE Best Practice: Set descriptive User-Agent so remote CDNs don't block request
	req.Header.Set("User-Agent", "md-terminal-viewer/1.0 (Darwin/macOS; iTerm2)")

	resp, err := client.Do(req)
	if err != nil {
		return nil, "", fmt.Errorf("http request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return nil, "", fmt.Errorf("http error %d: %s", resp.StatusCode, resp.Status)
	}

	// Enforce max size limit to prevent memory exhaustion
	limitReader := io.LimitReader(resp.Body, MaxImageSize+1)
	data, err := io.ReadAll(limitReader)
	if err != nil {
		return nil, "", fmt.Errorf("reading response body: %w", err)
	}
	if len(data) > MaxImageSize {
		return nil, "", fmt.Errorf("image exceeds maximum allowed size of %d bytes", MaxImageSize)
	}

	// Extract filename from URL path
	filename := "image.png"
	if parsed, err := url.Parse(rawURL); err == nil {
		base := filepath.Base(parsed.Path)
		if base != "" && base != "." && base != "/" {
			filename = base
		}
	}

	return data, filename, nil
}

func fetchLocal(filePath, basePath string) ([]byte, string, error) {
	// If path is relative and basePath is given, join them
	targetPath := filePath
	if !filepath.IsAbs(targetPath) && basePath != "" {
		targetPath = filepath.Join(basePath, targetPath)
	}

	cleanPath := filepath.Clean(targetPath)
	stat, err := os.Stat(cleanPath)
	if err != nil {
		return nil, "", fmt.Errorf("file not found: %w", err)
	}
	if stat.IsDir() {
		return nil, "", fmt.Errorf("path is a directory: %s", cleanPath)
	}
	if stat.Size() > MaxImageSize {
		return nil, "", fmt.Errorf("image exceeds maximum allowed size (%d bytes)", MaxImageSize)
	}

	data, err := os.ReadFile(cleanPath)
	if err != nil {
		return nil, "", fmt.Errorf("reading image file: %w", err)
	}

	return data, filepath.Base(cleanPath), nil
}

// GetDimensions reads the intrinsic width, height, and format of the image bytes.
func GetDimensions(data []byte) (Info, error) {
	cfg, format, err := image.DecodeConfig(bytes.NewReader(data))
	if err != nil {
		return Info{ByteSize: len(data)}, err
	}
	return Info{
		Width:    cfg.Width,
		Height:   cfg.Height,
		Format:   format,
		ByteSize: len(data),
	}, nil
}

// ITerm2Options configures the inline image rendering for iTerm2.
type ITerm2Options struct {
	Width               string // e.g. "auto", "100%", "80", "400px"
	Height              string // e.g. "auto", "20", "300px"
	PreserveAspectRatio bool   // default true
	InTmux              bool   // wraps in tmux DCS passthrough
}

// FormatITerm2 generates the OSC 1337 escape sequence for displaying an inline image.
func FormatITerm2(data []byte, filename string, opts ITerm2Options) string {
	b64Data := base64.StdEncoding.EncodeToString(data)
	b64Name := base64.StdEncoding.EncodeToString([]byte(filename))

	width := opts.Width
	if width == "" {
		width = "auto"
	}
	height := opts.Height
	if height == "" {
		height = "auto"
	}

	aspect := "1"
	if !opts.PreserveAspectRatio {
		aspect = "0"
	}

	// Standard iTerm2 OSC 1337 escape sequence:
	// \033]1337;File=name=<b64>;size=<len>;inline=1;width=<w>;height=<h>;preserveAspectRatio=1:<payload>\a
	seq := fmt.Sprintf(
		"\x1b]1337;File=name=%s;size=%d;inline=1;preserveAspectRatio=%s;width=%s;height=%s:%s\a",
		b64Name,
		len(data),
		aspect,
		width,
		height,
		b64Data,
	)

	// If inside tmux, wrap in DCS passthrough
	if opts.InTmux {
		escaped := strings.ReplaceAll(seq, "\x1b", "\x1b\x1b")
		return "\x1bPtmux;" + escaped + "\x1b\\"
	}

	return seq
}

// FormatFallback renders a clean ASCII / Unicode placeholder for terminals that don't support inline images.
func FormatFallback(altText, src string, errMsg string) string {
	var sb strings.Builder
	if altText == "" {
		altText = "Untitled Image"
	}

	if errMsg != "" {
		sb.WriteString(fmt.Sprintf("  ┌─ ⚠️  [Image Error] %s ─┐\n", altText))
		sb.WriteString(fmt.Sprintf("  │  Source: %s\n", src))
		sb.WriteString(fmt.Sprintf("  │  Reason: %s\n", errMsg))
		sb.WriteString("  └──────────────────────────────────────────────────────────┘")
	} else {
		sb.WriteString(fmt.Sprintf("  ┌─ 🖼️  [Image: %s] ─┐\n", altText))
		sb.WriteString(fmt.Sprintf("  │  Source: %s\n", src))
		sb.WriteString("  └──────────────────────────────────────────────────────────┘")
	}

	return sb.String()
}
