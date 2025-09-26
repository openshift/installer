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
)

type marketplaceStream map[string]*rhcos.Marketplace

// FetchRawCoreOSStream returns the raw stream metadata for the
// bootimages embedded in the installer.
func FetchRawCoreOSStream(ctx context.Context) ([]byte, error) {
	st, err := FetchCoreOSBuild(ctx)
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
func FetchCoreOSBuild(ctx context.Context) (*stream.Stream, error) {
	body, err := fetchRawCoreOSStream(ctx)
	if err != nil {
		return nil, err
	}
	var st stream.Stream
	if err := json.Unmarshal(body, &st); err != nil {
		return nil, fmt.Errorf("failed to parse CoreOS stream metadata: %w", err)
	}

	// Merge marketplace json file into stream json file
	mktBody, err := fetchRawMarketplaceStream()
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
		}
	}
	return &st, nil
}

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

// FindArtifactURL returns a single "disk" artifact type; this
// mainly abstracts over different compression formats like `qcow2.xz` and `qcow2.gz`.
//
// Use this function only for cases where there's a single artifact type, such
// as `qemu` and `openstack`.
//
// Some platforms have multiple artifact types; for example, `metal` has an ISO
// as well as PXE files.  This function will error in such a case.
func FindArtifactURL(artifacts stream.PlatformArtifacts) (string, error) {
	var artifact *stream.Artifact
	for _, v := range artifacts.Formats {
		if v.Disk != nil {
			if artifact != nil {
				return "", fmt.Errorf("multiple \"disk\" artifacts found")
			}
			artifact = v.Disk
		}
	}
	if artifact != nil {
		return FormatURLWithIntegrity(artifact)
	}
	return "", fmt.Errorf("no \"disk\" artifact found")
}

func fetchRawCoreOSStream(ctx context.Context) ([]byte, error) {
	file, err := data.Assets.Open(getStreamFileName())
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

func fetchRawMarketplaceStream() ([]byte, error) {
	file, err := data.Assets.Open(getMarketplaceStreamFileName())
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
