// package rhcos contains APIs for interacting with the RHEL (or Fedora) CoreOS
// bootimages embedded as stream metadata JSON with the installer
// For more information, see docs/dev/pinned-coreos.md

package rhcos

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/fs"
	"net/url"

	"github.com/coreos/stream-metadata-go/stream"
	"github.com/coreos/stream-metadata-go/stream/rhcos"
	"github.com/sirupsen/logrus"

	"github.com/openshift/installer/data"
	"github.com/openshift/installer/pkg/types"
)

type marketplaceStream map[string]*rhcos.Marketplace

// FetchRawCoreOSStream returns the raw stream metadata for the
// bootimages embedded in the installer.
func FetchRawCoreOSStream(ctx context.Context, osImageStream types.OSImageStream) ([]byte, error) {
	st, err := FetchCoreOSBuild(ctx, osImageStream)
	if err != nil {
		return nil, fmt.Errorf("failed to get combined CoreOS build: %w", err)
	}
	rawStream, err := json.Marshal(st)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal combined CoreOS stream: %w", err)
	}
	return rawStream, nil
}

// FetchCoreOSBuild returns the pinned version of RHEL/Fedora CoreOS used
// by the installer to provision the bootstrap node and control plane currently.
// For more information, see e.g. https://github.com/openshift/enhancements/pull/201
func FetchCoreOSBuild(ctx context.Context, osImageStream types.OSImageStream) (*stream.Stream, error) {
	body, err := fetchRawCoreOSStream(osImageStream)
	if err != nil {
		return nil, err
	}
	var st stream.Stream
	if err := json.Unmarshal(body, &st); err != nil {
		return nil, fmt.Errorf("failed to parse CoreOS stream metadata: %w", err)
	}

	// Merge marketplace json file into stream json file
	mktBody, err := fetchRawMarketplaceStream(osImageStream)
	if err != nil {
		if errors.Is(err, fs.ErrNotExist) {
			logrus.Debug("No marketplace json file found: skipping merge.")
			return &st, nil
		}
		return nil, err
	}
	var mktSt marketplaceStream
	if err := json.Unmarshal(mktBody, &mktSt); err != nil {
		return nil, fmt.Errorf("failed to parse marketplace stream: %w", err)
	}

	for name, arch := range st.Architectures {
		if mkt, ok := mktSt[name]; ok {
			if arch.RHELCoreOSExtensions == nil {
				arch.RHELCoreOSExtensions = &rhcos.Extensions{}
			}
			arch.RHELCoreOSExtensions.Marketplace = mkt
			st.Architectures[name] = arch
		}
	}
	return &st, nil
}

// CompressedSHA256Param is the query parameter carrying the sha256 of an
// artifact as it is served, i.e. before decompression.  It is deliberately
// distinct from the "sha256" parameter set by FormatURLWithIntegrity, which by
// convention holds the digest of the *uncompressed* artifact.
const CompressedSHA256Param = "compressed-sha256"

// FormatURLWithIntegrity squashes an artifact into a URL string
// with the uncompressed sha256 as a query parameter.  This is necessary
// currently because various parts of the installer pass around this
// reference as a string, and it's also exposed to users via install-config overrides.
func FormatURLWithIntegrity(artifact *stream.Artifact) (string, error) {
	u, err := url.Parse(artifact.Location)
	if err != nil {
		return "", fmt.Errorf("failed to parse artifact URL: %v", err)
	}
	q := u.Query()
	q.Set("sha256", artifact.UncompressedSha256)
	u.RawQuery = q.Encode()
	return u.String(), nil
}

// FormatURLWithCompressedIntegrity squashes an artifact into a URL string with
// the sha256 of the artifact as served -- before any decompression -- as a
// query parameter.  Use this rather than FormatURLWithIntegrity for consumers
// that keep the artifact compressed, such as the GCP image upload, which needs
// the tar.gz byte for byte.
func FormatURLWithCompressedIntegrity(artifact *stream.Artifact) (string, error) {
	// Streams are not required to carry a digest, but an artifact that cannot be
	// verified must not be passed off as one that can: callers only skip
	// verification when the parameter is absent, which would be silent.
	if artifact.Sha256 == "" {
		return "", fmt.Errorf("artifact %s has no sha256 checksum", artifact.Location)
	}
	u, err := url.Parse(artifact.Location)
	if err != nil {
		return "", fmt.Errorf("failed to parse artifact URL: %v", err)
	}
	q := u.Query()
	q.Set(CompressedSHA256Param, artifact.Sha256)
	u.RawQuery = q.Encode()
	return u.String(), nil
}

// FindArtifactURL returns a single "disk" artifact type; this
// mainly abstracts over different compression formats like `qcow2.xz` and `qcow2.gz`.
//
// Use this function only for cases where there's a single artifact type, such
// as `qemu` and `openstack`.
//
// Some platforms have multiple artifact types; for example, `metal` has an ISO
// as well as PXE files.  This function will error in such a case.
func FindArtifactURL(artifacts stream.PlatformArtifacts) (string, error) {
	artifact, err := findDiskArtifact(artifacts)
	if err != nil {
		return "", err
	}
	return FormatURLWithIntegrity(artifact)
}

// FindCompressedArtifactURL is FindArtifactURL for consumers that use the
// artifact in the compressed form it is served in; the returned URL carries the
// compressed sha256 rather than the uncompressed one.
func FindCompressedArtifactURL(artifacts stream.PlatformArtifacts) (string, error) {
	artifact, err := findDiskArtifact(artifacts)
	if err != nil {
		return "", err
	}
	return FormatURLWithCompressedIntegrity(artifact)
}

// findDiskArtifact returns the single "disk" artifact among the given platform
// artifacts, erroring when there is not exactly one.
func findDiskArtifact(artifacts stream.PlatformArtifacts) (*stream.Artifact, error) {
	var artifact *stream.Artifact
	for _, v := range artifacts.Formats {
		if v.Disk != nil {
			if artifact != nil {
				return nil, fmt.Errorf("multiple \"disk\" artifacts found")
			}
			artifact = v.Disk
		}
	}
	if artifact == nil {
		return nil, fmt.Errorf("no \"disk\" artifact found")
	}
	return artifact, nil
}

func fetchRawCoreOSStream(osImageStream types.OSImageStream) ([]byte, error) {
	file, err := data.Assets.Open(getStreamFileName(osImageStream))
	if err != nil {
		return nil, fmt.Errorf("failed to read embedded CoreOS stream metadata: %w", err)
	}
	defer file.Close()

	body, err := io.ReadAll(file)
	if err != nil {
		return nil, fmt.Errorf("failed to read CoreOS stream metadata: %w", err)
	}
	return body, nil
}

func fetchRawMarketplaceStream(osImageStream types.OSImageStream) ([]byte, error) {
	file, err := data.Assets.Open(getMarketplaceStreamFileName(osImageStream))
	if err != nil {
		return nil, err
	}
	defer file.Close()

	body, err := io.ReadAll(file)
	if err != nil {
		return nil, err
	}
	return body, nil
}
