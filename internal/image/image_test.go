package image

import (
	"bytes"
	"context"
	"encoding/base64"
	"image"
	"image/color"
	"image/png"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// Helper to generate a 2x2 PNG in memory
func generateTestPNG(t *testing.T) []byte {
	t.Helper()
	img := image.NewRGBA(image.Rect(0, 0, 2, 2))
	img.Set(0, 0, color.RGBA{R: 255, A: 255})
	img.Set(1, 1, color.RGBA{B: 255, A: 255})

	var buf bytes.Buffer
	if err := png.Encode(&buf, img); err != nil {
		t.Fatalf("failed to encode test png: %v", err)
	}
	return buf.Bytes()
}

func TestFetch_Local(t *testing.T) {
	pngData := generateTestPNG(t)
	tmpDir := t.TempDir()
	imgPath := filepath.Join(tmpDir, "test.png")
	if err := os.WriteFile(imgPath, pngData, 0600); err != nil {
		t.Fatalf("write file failed: %v", err)
	}

	ctx := context.Background()

	// 1. Absolute path
	got, name, err := Fetch(ctx, imgPath, "", nil)
	if err != nil {
		t.Fatalf("Fetch absolute failed: %v", err)
	}
	if !bytes.Equal(got, pngData) {
		t.Errorf("expected fetched data to match source")
	}
	if name != "test.png" {
		t.Errorf("expected filename test.png, got %s", name)
	}

	// 2. Relative path with basePath
	gotRel, _, err := Fetch(ctx, "test.png", tmpDir, nil)
	if err != nil {
		t.Fatalf("Fetch relative failed: %v", err)
	}
	if !bytes.Equal(gotRel, pngData) {
		t.Errorf("expected fetched relative data to match")
	}

	// 3. Nonexistent file
	_, _, err = Fetch(ctx, "missing.png", tmpDir, nil)
	if err == nil {
		t.Errorf("expected error for missing file, got nil")
	}
}

func TestFetch_DataURI(t *testing.T) {
	pngData := generateTestPNG(t)
	b64 := base64.StdEncoding.EncodeToString(pngData)
	dataURI := "data:image/png;base64," + b64

	ctx := context.Background()
	got, name, err := Fetch(ctx, dataURI, "", nil)
	if err != nil {
		t.Fatalf("Fetch data URI failed: %v", err)
	}
	if !bytes.Equal(got, pngData) {
		t.Errorf("expected decoded data to match")
	}
	if name != "image.png" {
		t.Errorf("expected name image.png, got %s", name)
	}
}

func TestFetch_Remote(t *testing.T) {
	pngData := generateTestPNG(t)
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/logo.png" {
			w.Header().Set("Content-Type", "image/png")
			_, _ = w.Write(pngData)
			return
		}
		http.NotFound(w, r)
	}))
	defer server.Close()

	ctx := context.Background()

	// Success case
	got, name, err := Fetch(ctx, server.URL+"/logo.png", "", server.Client())
	if err != nil {
		t.Fatalf("Fetch remote failed: %v", err)
	}
	if !bytes.Equal(got, pngData) {
		t.Errorf("expected fetched remote bytes to match")
	}
	if name != "logo.png" {
		t.Errorf("expected name logo.png, got %s", name)
	}

	// 404 case
	_, _, err = Fetch(ctx, server.URL+"/notfound.png", "", server.Client())
	if err == nil {
		t.Errorf("expected 404 error, got nil")
	}
}

func TestGetDimensions(t *testing.T) {
	pngData := generateTestPNG(t)
	info, err := GetDimensions(pngData)
	if err != nil {
		t.Fatalf("GetDimensions failed: %v", err)
	}
	if info.Width != 2 || info.Height != 2 {
		t.Errorf("expected 2x2, got %dx%d", info.Width, info.Height)
	}
	if info.Format != "png" {
		t.Errorf("expected format png, got %s", info.Format)
	}
}

func TestFormatITerm2(t *testing.T) {
	pngData := []byte("fake-png-payload")
	opts := ITerm2Options{
		Width:               "auto",
		Height:              "auto",
		PreserveAspectRatio: true,
		InTmux:              false,
	}

	seq := FormatITerm2(pngData, "diagram.png", opts)
	if !strings.HasPrefix(seq, "\x1b]1337;File=") {
		t.Errorf("expected OSC 1337 prefix, got %q", seq)
	}
	if !strings.HasSuffix(seq, "\a") {
		t.Errorf("expected BEL termination, got %q", seq)
	}
	if !strings.Contains(seq, "inline=1") {
		t.Errorf("expected inline=1")
	}

	// Test tmux DCS wrapper
	opts.InTmux = true
	tmuxSeq := FormatITerm2(pngData, "diagram.png", opts)
	if !strings.HasPrefix(tmuxSeq, "\x1bPtmux;") {
		t.Errorf("expected tmux DCS prefix, got %q", tmuxSeq)
	}
	if !strings.HasSuffix(tmuxSeq, "\x1b\\") {
		t.Errorf("expected tmux DCS suffix, got %q", tmuxSeq)
	}
}

func TestFormatFallback(t *testing.T) {
	out := FormatFallback("Alt Text", "test.png", "")
	if !strings.Contains(out, "Alt Text") || !strings.Contains(out, "test.png") {
		t.Errorf("unexpected fallback output: %s", out)
	}

	errOut := FormatFallback("Alt Text", "test.png", "not found")
	if !strings.Contains(errOut, "Image Error") || !strings.Contains(errOut, "not found") {
		t.Errorf("unexpected error fallback output: %s", errOut)
	}
}
