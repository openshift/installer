package imagesource

import (
	"path/filepath"
	"strings"
)

// generateDigestPath generates a platform-specific file path for the given digest.
// This uses `filepath.Join` for the `elem` parameter.
// Digests are typically in the format of `algo:hash`. Some platforms, such as
// Windows, do not allow for the use of `:` characters in file paths. In this case,
// the `:` character is replaced with `-`.
func generateDigestPath(digest string, elem ...string) string {
	dgst := strings.ReplaceAll(digest, ":", "-")
	return filepath.Join(append(elem, dgst)...)
}
