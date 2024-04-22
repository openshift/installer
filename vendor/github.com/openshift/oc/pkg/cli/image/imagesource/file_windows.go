package imagesource

import (
	"fmt"
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

// reconstructDigest replaces ":" with "-" which makes the digest parsing nondeterministic.
// A digest regex corresponds to `[a-z0-9](?:[.+_-][a-z0-9])*:[a-zA-Z0-9=_-]+`.
// Making "sha256-sha256:659d863050ebd58ebd2ea4b4e89226b494ba39457932ebf11669763eea3b2ed0" a valid digest.
// After the ":" replacement the filename gets translated into
// "sha256-sha256-659d863050ebd58ebd2ea4b4e89226b494ba39457932ebf11669763eea3b2ed0" which has
// two matching digests:
// - sha256-sha256:659d863050ebd58ebd2ea4b4e89226b494ba39457932ebf11669763eea3b2ed0
// - sha256:sha256-659d863050ebd58ebd2ea4b4e89226b494ba39457932ebf11669763eea3b2ed0
//
// From https://github.com/opencontainers/image-spec/blob/main/descriptor.md#digests:
// digest                ::= algorithm ":" encoded
// algorithm             ::= algorithm-component (algorithm-separator algorithm-component)*
// algorithm-component   ::= [a-z0-9]+
// algorithm-separator   ::= [+._-]
// encoded               ::= [a-zA-Z0-9=_-]+
//
// Given only sha256 and sha512 are currently supported (based on the linked descriptor.md#digests)
// we can parse these two cases and panic for the rest.
func reconstructDigest(digest string) string {
	if strings.HasPrefix(digest, "sha256-") {
		return "sha256:" + strings.TrimPrefix(digest, "sha256-")
	}
	if strings.HasPrefix(digest, "sha512-") {
		return "sha512:" + strings.TrimPrefix(digest, "sha512-")
	}
	panic(fmt.Sprintf("Only 'sha256-' and 'sha512-' digest prefixes are supported on Windows. Got %q", digest))
}
