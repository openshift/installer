//go:build !windows
// +build !windows

package imagesource

import "path/filepath"

// generateDigestPath generates a platform-specific file path for the given digest.
// This uses `filepath.Join` for the `elem` parameter.
// Digests are typically in the format of `algo:hash`. Some platforms, such as
// Windows, do not allow for the use of `:` characters in file paths. In this case,
// the `:` character is replaced with `-`.
func generateDigestPath(digest string, elem ...string) string {
	return filepath.Join(append(elem, digest)...)
}

// reconstructDigest is a no-op on Unix. See file_windows.go for more details.
func reconstructDigest(digest string) string {
	return digest
}
