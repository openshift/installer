package cache

import (
	"bytes"
	"compress/gzip"
	"crypto/sha256"
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"path/filepath"
	"testing"

	"github.com/klauspost/compress/zstd"
	"github.com/stretchr/testify/assert"
	"github.com/ulikunitz/xz"
)

// testPayload must be large enough that its compressed form exceeds 512 bytes,
// since cacheFile reads a 512-byte buffer for magic byte detection.
var testPayload = bytes.Repeat([]byte("this is a test payload for RHCOS image decompression testing\n"), 100)

func compressGz(t *testing.T, data []byte) []byte {
	t.Helper()
	var buf bytes.Buffer
	w := gzip.NewWriter(&buf)
	_, err := w.Write(data)
	assert.NoError(t, err)
	assert.NoError(t, w.Close())
	return buf.Bytes()
}

func compressXz(t *testing.T, data []byte) []byte {
	t.Helper()
	var buf bytes.Buffer
	w, err := xz.NewWriter(&buf)
	assert.NoError(t, err)
	_, err = w.Write(data)
	assert.NoError(t, err)
	assert.NoError(t, w.Close())
	return buf.Bytes()
}

func compressZstd(t *testing.T, data []byte) []byte {
	t.Helper()
	var buf bytes.Buffer
	w, err := zstd.NewWriter(&buf)
	assert.NoError(t, err)
	_, err = w.Write(data)
	assert.NoError(t, err)
	assert.NoError(t, w.Close())
	return buf.Bytes()
}

func sha256sum(data []byte) string {
	h := sha256.Sum256(data)
	return fmt.Sprintf("%x", h[:])
}

func TestCacheFileDecompression(t *testing.T) {
	tests := []struct {
		name       string
		sourceName string
		compress   func(*testing.T, []byte) []byte
	}{
		{
			name:       "gzip",
			sourceName: "image.qcow2.gz",
			compress:   compressGz,
		},
		{
			name:       "xz",
			sourceName: "image.qcow2.xz",
			compress:   compressXz,
		},
		{
			name:       "zstd",
			sourceName: "image.qcow2.zst",
			compress:   compressZstd,
		},
		{
			name:       "uncompressed",
			sourceName: "image.qcow2",
			compress:   func(_ *testing.T, data []byte) []byte { return data },
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			compressed := tt.compress(t, testPayload)
			checksum := sha256sum(testPayload)

			dir := t.TempDir()
			outPath := filepath.Join(dir, "image.qcow2")

			err := cacheFile(bytes.NewReader(compressed), outPath, checksum, tt.sourceName)
			if err != nil {
				t.Fatalf("cacheFile failed: %v", err)
			}

			got, err := os.ReadFile(outPath)
			if err != nil {
				t.Fatalf("reading output file: %v", err)
			}
			assert.Equal(t, testPayload, got)
		})
	}
}

func TestCacheFileChecksumMismatch(t *testing.T) {
	compressed := compressGz(t, testPayload)

	dir := t.TempDir()
	outPath := filepath.Join(dir, "image.qcow2")

	err := cacheFile(bytes.NewReader(compressed), outPath, "badchecksum", "image.qcow2.gz")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "Checksum mismatch")
}

func TestUncompressedName(t *testing.T) {
	tests := []struct {
		path string
		want string
	}{
		{"https://example.com/image.qcow2.gz", "image.qcow2"},
		{"https://example.com/image.qcow2.xz", "image.qcow2"},
		{"https://example.com/image.qcow2.zst", "image.qcow2"},
		{"https://example.com/image.qcow2", "image.qcow2"},
		{"https://example.com/image.raw.gz", "image.raw"},
		{"https://example.com/image.vmdk.zst", "image.vmdk"},
	}

	for _, tt := range tests {
		t.Run(tt.path, func(t *testing.T) {
			parsed, err := url.Parse(tt.path)
			if err != nil {
				t.Fatalf("parsing URL: %v", err)
			}
			u := &urlWithIntegrity{location: *parsed}
			assert.Equal(t, tt.want, u.uncompressedName())
		})
	}
}

func TestDownloadImageFile(t *testing.T) {
	tests := []struct {
		name     string
		filename string
		compress func(*testing.T, []byte) []byte
	}{
		{
			name:     "gzip",
			filename: "image.qcow2.gz",
			compress: compressGz,
		},
		{
			name:     "zstd",
			filename: "image.qcow2.zst",
			compress: compressZstd,
		},
		{
			name:     "uncompressed",
			filename: "image.qcow2",
			compress: func(_ *testing.T, data []byte) []byte { return data },
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			compressed := tt.compress(t, testPayload)
			checksum := sha256sum(testPayload)

			srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				_, err := w.Write(compressed)
				if err != nil {
					t.Fatalf("writing response: %v", err)
				}
			}))
			defer srv.Close()

			imageURL := fmt.Sprintf("%s/%s?sha256=%s", srv.URL, tt.filename, checksum)

			// Override HOME so os.UserCacheDir() uses a temp directory
			tmpDir := t.TempDir()
			t.Setenv("HOME", tmpDir)

			filePath, err := DownloadImageFile(imageURL, "test")
			if err != nil {
				t.Fatalf("DownloadImageFile failed: %v", err)
			}

			got, err := os.ReadFile(filePath)
			if err != nil {
				t.Fatalf("reading downloaded file: %v", err)
			}
			assert.Equal(t, testPayload, got)
		})
	}
}
