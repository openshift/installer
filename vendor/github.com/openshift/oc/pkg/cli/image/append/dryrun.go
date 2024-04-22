package append

import (
	"context"
	"net/http"

	"k8s.io/klog/v2"

	"github.com/distribution/distribution/v3"
	digest "github.com/opencontainers/go-digest"

	"github.com/openshift/library-go/pkg/image/registryclient"
)

// dryRunManifestService emulates a remote registry for dry run behavior
type dryRunManifestService struct{}

func (s *dryRunManifestService) Exists(ctx context.Context, dgst digest.Digest) (bool, error) {
	panic("not implemented")
}

func (s *dryRunManifestService) Get(ctx context.Context, dgst digest.Digest, options ...distribution.ManifestServiceOption) (distribution.Manifest, error) {
	panic("not implemented")
}

func (s *dryRunManifestService) Put(ctx context.Context, manifest distribution.Manifest, options ...distribution.ManifestServiceOption) (digest.Digest, error) {
	klog.V(4).Infof("Manifest: %#v", manifest.References())
	return registryclient.ContentDigestForManifest(manifest, digest.SHA256)
}

func (s *dryRunManifestService) Delete(ctx context.Context, dgst digest.Digest) error {
	panic("not implemented")
}

// dryRunBlobStore emulates a remote registry for dry run behavior
type dryRunBlobStore struct {
	layers []distribution.Descriptor
}

func (s *dryRunBlobStore) Stat(ctx context.Context, dgst digest.Digest) (distribution.Descriptor, error) {
	for _, layer := range s.layers {
		if layer.Digest == dgst {
			return layer, nil
		}
	}
	return distribution.Descriptor{}, distribution.ErrBlobUnknown
}

func (s *dryRunBlobStore) Get(ctx context.Context, dgst digest.Digest) ([]byte, error) {
	panic("not implemented")
}

func (s *dryRunBlobStore) Open(ctx context.Context, dgst digest.Digest) (distribution.ReadSeekCloser, error) {
	panic("not implemented")
}

func (s *dryRunBlobStore) Put(ctx context.Context, mediaType string, p []byte) (distribution.Descriptor, error) {
	return distribution.Descriptor{
		MediaType: mediaType,
		Size:      int64(len(p)),
		Digest:    digest.SHA256.FromBytes(p),
	}, nil
}

func (s *dryRunBlobStore) Create(ctx context.Context, options ...distribution.BlobCreateOption) (distribution.BlobWriter, error) {
	panic("not implemented")
}

func (s *dryRunBlobStore) Resume(ctx context.Context, id string) (distribution.BlobWriter, error) {
	panic("not implemented")
}

func (s *dryRunBlobStore) ServeBlob(ctx context.Context, w http.ResponseWriter, r *http.Request, dgst digest.Digest) error {
	panic("not implemented")
}

func (s *dryRunBlobStore) Delete(ctx context.Context, dgst digest.Digest) error {
	panic("not implemented")
}
