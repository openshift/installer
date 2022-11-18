package append

import (
	"bytes"
	"context"
	"net/http"

	"github.com/docker/distribution"
	"github.com/docker/distribution/reference"
	digest "github.com/opencontainers/go-digest"

	"github.com/openshift/oc/pkg/helpers/image/dockerlayer"
)

// scratchRepo can serve the scratch image blob.
type scratchRepo struct{}

var _ distribution.Repository = scratchRepo{}

func (_ scratchRepo) Named() reference.Named { panic("not implemented") }
func (_ scratchRepo) Tags(ctx context.Context) distribution.TagService {
	panic("not implemented")
}
func (_ scratchRepo) Manifests(ctx context.Context, options ...distribution.ManifestServiceOption) (distribution.ManifestService, error) {
	panic("not implemented")
}

func (r scratchRepo) Blobs(ctx context.Context) distribution.BlobStore { return r }

func (_ scratchRepo) Stat(ctx context.Context, dgst digest.Digest) (distribution.Descriptor, error) {
	if dgst != dockerlayer.GzippedEmptyLayerDigest {
		return distribution.Descriptor{}, distribution.ErrBlobUnknown
	}
	return distribution.Descriptor{
		MediaType: "application/vnd.docker.image.rootfs.diff.tar.gzip",
		Digest:    digest.Digest(dockerlayer.GzippedEmptyLayerDigest),
		Size:      int64(len(dockerlayer.GzippedEmptyLayer)),
	}, nil
}

func (_ scratchRepo) Get(ctx context.Context, dgst digest.Digest) ([]byte, error) {
	if dgst != dockerlayer.GzippedEmptyLayerDigest {
		return nil, distribution.ErrBlobUnknown
	}
	return dockerlayer.GzippedEmptyLayer, nil
}

type nopCloseBuffer struct {
	*bytes.Buffer
}

func (_ nopCloseBuffer) Seek(offset int64, whence int) (int64, error) {
	return 0, nil
}

func (_ nopCloseBuffer) Close() error {
	return nil
}

func (_ scratchRepo) Open(ctx context.Context, dgst digest.Digest) (distribution.ReadSeekCloser, error) {
	if dgst != dockerlayer.GzippedEmptyLayerDigest {
		return nil, distribution.ErrBlobUnknown
	}
	return nopCloseBuffer{bytes.NewBuffer(dockerlayer.GzippedEmptyLayer)}, nil
}

func (_ scratchRepo) Put(ctx context.Context, mediaType string, p []byte) (distribution.Descriptor, error) {
	panic("not implemented")
}

func (_ scratchRepo) Create(ctx context.Context, options ...distribution.BlobCreateOption) (distribution.BlobWriter, error) {
	panic("not implemented")
}

func (_ scratchRepo) Resume(ctx context.Context, id string) (distribution.BlobWriter, error) {
	panic("not implemented")
}

func (_ scratchRepo) ServeBlob(ctx context.Context, w http.ResponseWriter, r *http.Request, dgst digest.Digest) error {
	panic("not implemented")
}

func (_ scratchRepo) Delete(ctx context.Context, dgst digest.Digest) error {
	panic("not implemented")
}
