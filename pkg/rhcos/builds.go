// package rhcos contains APIs for interacting with the RHEL (or Fedora) CoreOS
// bootimages embedded as stream metadata JSON with the installer
// For more information, see docs/dev/pinned-coreos.md

package rhcos

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/url"

	"github.com/coreos/stream-metadata-go/stream"
	"github.com/openshift/installer/data"
	"github.com/pkg/errors"
)

// FetchRawCoreOSStream returns the raw stream metadata for the
// bootimages embedded in the installer.
func FetchRawCoreOSStream(ctx context.Context) ([]byte, error) {
	file, err := data.Assets.Open(getStreamFileName())
	if err != nil {
		return nil, errors.Wrapf(err, "failed to read embedded CoreOS stream metadata")
	}
	defer file.Close()

	body, err := ioutil.ReadAll(file)
	if err != nil {
		return nil, errors.Wrap(err, "failed to read CoreOS stream metadata")
	}
	return body, nil
}

// FetchCoreOSBuild returns the pinned version of RHEL/Fedora CoreOS used
// by the installer to provision the bootstrap node and control plane currently.
// For more information, see e.g. https://github.com/openshift/enhancements/pull/201
func FetchCoreOSBuild(ctx context.Context) (*stream.Stream, error) {
	body, err := FetchRawCoreOSStream(ctx)
	if err != nil {
		return nil, err
	}
	var st stream.Stream
	if err := json.Unmarshal(body, &st); err != nil {
		return nil, errors.Wrap(err, "failed to parse CoreOS stream metadata")
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
// mainly abstracts over e.g. `qcow2.xz` and `qcow2.gz`.  (FCOS uses
// xz, RHCOS uses gzip right now)
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
