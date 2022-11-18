package imagesource

import (
	"context"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strings"
	"time"

	"k8s.io/klog/v2"

	man "github.com/containers/image/v5/manifest"
	"github.com/docker/distribution"
	"github.com/docker/distribution/reference"
	"github.com/opencontainers/go-digest"
	godigest "github.com/opencontainers/go-digest"
)

type fileDriver struct {
	BaseDir string
}

func (d *fileDriver) Repository(ctx context.Context, server *url.URL, repoName string, insecure bool) (distribution.Repository, error) {
	klog.V(3).Infof("Repository %s %s", server, repoName)

	ref, err := reference.Parse(repoName)
	if err != nil {
		return nil, err
	}
	named, ok := ref.(reference.Named)
	if !ok {
		return nil, fmt.Errorf("%s is not a valid repository name", repoName)
	}

	repo := &fileRepository{
		repoName: named,
		repoPath: repoPathForName(repoName),
		basePath: d.BaseDir,
	}
	return repo, nil
}

func repoPathForName(repoName string) string {
	return strings.ReplaceAll(repoName, "/", string(filepath.Separator))
}

type fileRepository struct {
	basePath string
	repoPath string
	repoName reference.Named
}

// Named returns the name of the repository.
func (r *fileRepository) Named() reference.Named {
	return r.repoName
}

// Manifests returns a reference to this repository's manifest service.
// with the supplied options applied.
func (r *fileRepository) Manifests(ctx context.Context, options ...distribution.ManifestServiceOption) (distribution.ManifestService, error) {
	return &fileManifestService{r: r}, nil
}

// Blobs returns a reference to this repository's blob service.
func (r *fileRepository) Blobs(ctx context.Context) distribution.BlobStore {
	return &fileBlobStore{r: r}
}

// Tags returns a reference to this repositories tag service
func (r *fileRepository) Tags(ctx context.Context) distribution.TagService {
	return &fileTagStore{r: r}
}

type fileTagStore struct {
	r *fileRepository
}

// Get retrieves the descriptor identified by the tag. Some
// implementations may differentiate between "trusted" tags and
// "untrusted" tags. If a tag is "untrusted", the mapping will be returned
// as an ErrTagUntrusted error, with the target descriptor.
func (s *fileTagStore) Get(ctx context.Context, tag string) (distribution.Descriptor, error) {
	path := filepath.Join(s.r.basePath, "v2", s.r.repoPath, "manifests", tag)
	fi, err := os.Lstat(path)
	if err != nil {
		if os.IsNotExist(err) {
			return distribution.Descriptor{}, distribution.ErrBlobUnknown
		}
		return distribution.Descriptor{}, err
	}
	if fi.Mode()&os.ModeSymlink == 0 {
		return distribution.Descriptor{}, fmt.Errorf("not a symlink")
	}
	target, err := os.Readlink(path)
	if err != nil {
		return distribution.Descriptor{}, err
	}
	return distribution.Descriptor{
		Digest: godigest.Digest(filepath.Base(target)),
		Size:   fi.Size(),
	}, nil
}

// Tag associates the tag with the provided descriptor, updating the
// current association, if needed.
func (s *fileTagStore) Tag(ctx context.Context, tag string, desc distribution.Descriptor) error {
	return fmt.Errorf("tagging images in local file storage is not supported")
}

// Untag removes the given tag association
func (s *fileTagStore) Untag(ctx context.Context, tag string) error {
	return fmt.Errorf("removing tags from images in local file storage is not supported")
}

// All returns the set of tags managed by this tag service
func (s *fileTagStore) All(ctx context.Context) ([]string, error) {
	path := filepath.Join(s.r.basePath, "v2", s.r.repoPath, "manifests")
	fis, err := ioutil.ReadDir(path)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, nil
		}
		return nil, err
	}
	var names []string
	for _, fi := range fis {
		if fi.Mode()&os.ModeSymlink != 0 {
			names = append(names, filepath.Base(fi.Name()))
		}
	}
	return names, nil
}

// Lookup returns the set of tags referencing the given digest.
func (s *fileTagStore) Lookup(ctx context.Context, digest distribution.Descriptor) ([]string, error) {
	return nil, fmt.Errorf("retrieving tags for a digest in local file storage is not supported")
}

type fileManifestService struct {
	r *fileRepository
}

// Exists returns true if the manifest exists.
func (s *fileManifestService) Exists(ctx context.Context, dgst godigest.Digest) (bool, error) {
	path := generateDigestPath(dgst.String(), s.r.basePath, "v2", s.r.repoPath, "manifests")
	fi, err := os.Lstat(path)
	if err != nil {
		if os.IsNotExist(err) {
			return false, nil
		}
		return false, err
	}
	if fi.IsDir() || fi.Mode()&os.ModeSymlink != 0 {
		return false, fmt.Errorf("not a file: %s", fi.Name())
	}
	return true, nil
}

// Get retrieves the manifest specified by the given digest
func (s *fileManifestService) Get(ctx context.Context, dgst godigest.Digest, options ...distribution.ManifestServiceOption) (distribution.Manifest, error) {
	path := generateDigestPath(dgst.String(), s.r.basePath, "v2", s.r.repoPath, "manifests")
	data, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}

	manifest, desc, err := distribution.UnmarshalManifest(man.GuessMIMEType(data), data)
	if err != nil {
		return nil, err
	}
	klog.V(5).Infof("Read manifest %T from %s: %v", manifest, path, desc)
	return manifest, nil
}

// atomicWrite performs an atomic write and move of a file. It expects the destination
// to be a file or to be missing.
func atomicWrite(path string, data []byte) error {
	if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
		return err
	}
	if err := ioutil.WriteFile(path+".download", data, 0600); err != nil {
		return err
	}
	return os.Rename(path+".download", path)
}

// atomicWrite performs an atomic hardlink and move of a file. It expects the destination
// to be a file or to be missing.
func atomicLink(path string, sourcePath string) error {
	if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
		return err
	}
	if err := os.Link(sourcePath, path+".download"); err != nil {
		return err
	}
	return os.Rename(path+".download", path)
}

// atomicWrite performs an atomic symlink and move of a file. It expects the destination
// to be a file or to be missing.
func atomicSymlink(path string, sourcePath string) error {
	if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
		return err
	}
	if err := os.Symlink(sourcePath, path+".download"); err != nil {
		return err
	}
	return os.Rename(path+".download", path)
}

// Put creates or updates the given manifest returning the manifest digest
func (s *fileManifestService) Put(ctx context.Context, manifest distribution.Manifest, options ...distribution.ManifestServiceOption) (godigest.Digest, error) {
	_, payload, err := manifest.Payload()
	if err != nil {
		return "", err
	}
	dgst := godigest.FromBytes(payload)
	path := generateDigestPath(dgst.String(), s.r.basePath, "v2", s.r.repoPath, "blobs")
	if err := atomicWrite(path, payload); err != nil {
		return "", err
	}

	// set tags and manifests
	manifestPath := generateDigestPath(dgst.String(), s.r.basePath, "v2", s.r.repoPath, "manifests")
	if err := atomicLink(manifestPath, path); err != nil {
		return "", err
	}
	for _, option := range options {
		if opt, ok := option.(distribution.WithTagOption); ok {
			tagPath := filepath.Join(s.r.basePath, "v2", s.r.repoPath, "manifests", opt.Tag)
			if err := atomicSymlink(tagPath, generateDigestPath(dgst.String())); err != nil {
				return "", err
			}
		}
	}
	return dgst, nil
}

// Delete removes the manifest specified by the given digest. Deleting
// a manifest that doesn't exist will return ErrManifestNotFound
func (s *fileManifestService) Delete(ctx context.Context, dgst godigest.Digest) error {
	return fmt.Errorf("unimplemented")
}

type fileBlobStore struct {
	r *fileRepository
}

func (s *fileBlobStore) Stat(ctx context.Context, dgst godigest.Digest) (distribution.Descriptor, error) {
	path := generateDigestPath(dgst.String(), s.r.basePath, "v2", s.r.repoPath, "blobs")
	fi, err := os.Stat(path)
	if err != nil {
		if os.IsNotExist(err) {
			return distribution.Descriptor{}, distribution.ErrBlobUnknown
		}
		return distribution.Descriptor{}, err
	}
	if fi.IsDir() {
		return distribution.Descriptor{}, fmt.Errorf("not a file")
	}
	return distribution.Descriptor{
		Digest: dgst,
		Size:   fi.Size(),
	}, nil
}

func (s *fileBlobStore) Delete(ctx context.Context, dgst godigest.Digest) error {
	return fmt.Errorf("unimplemented")
}

func (s *fileBlobStore) Get(ctx context.Context, dgst godigest.Digest) ([]byte, error) {
	path := generateDigestPath(dgst.String(), s.r.basePath, "v2", s.r.repoPath, "blobs")
	return ioutil.ReadFile(path)
}

func (s *fileBlobStore) Open(ctx context.Context, dgst godigest.Digest) (distribution.ReadSeekCloser, error) {
	path := generateDigestPath(dgst.String(), s.r.basePath, "v2", s.r.repoPath, "blobs")
	return os.Open(path)
}

func (s *fileBlobStore) ServeBlob(ctx context.Context, w http.ResponseWriter, r *http.Request, dgst godigest.Digest) error {
	return fmt.Errorf("unimplemented")
}

func (s *fileBlobStore) Put(ctx context.Context, mediaType string, payload []byte) (distribution.Descriptor, error) {
	dgst := godigest.FromBytes(payload)
	path := generateDigestPath(dgst.String(), s.r.basePath, "v2", s.r.repoPath, "blobs")
	if err := atomicWrite(path, payload); err != nil {
		return distribution.Descriptor{}, err
	}
	return distribution.Descriptor{MediaType: mediaType, Size: int64(len(payload)), Digest: dgst}, nil
}

func (s *fileBlobStore) Create(ctx context.Context, options ...distribution.BlobCreateOption) (distribution.BlobWriter, error) {
	var opts distribution.CreateOptions
	for _, option := range options {
		err := option.Apply(&opts)
		if err != nil {
			return nil, err
		}
	}

	if opts.Mount.Stat == nil || len(opts.Mount.Stat.Digest) == 0 {
		dir := filepath.Join(s.r.basePath, "v2", s.r.repoPath, "blobs")
		if err := os.MkdirAll(dir, 0755); err != nil {
			return nil, err
		}
		f, err := ioutil.TempFile(dir, ".tmp-")
		if err != nil {
			return nil, err
		}
		return s.r.newWriterTemp(dir, f), nil
	}

	d := opts.Mount.Stat.Digest
	path := generateDigestPath(d.String(), s.r.basePath, "v2", s.r.repoPath, "blobs")
	if opts.Mount.ShouldMount {
		repoPath := repoPathForName(opts.Mount.From.Name())
		sourcePath := generateDigestPath(d.String(), s.r.basePath, "v2", repoPath, "blobs")
		klog.V(5).Infof("Try mount from %s to %s", sourcePath, path)
		if fi, err := os.Stat(sourcePath); err == nil && !fi.IsDir() {
			if err := atomicLink(path, sourcePath); err != nil {
				return nil, err
			}
			return nil, distribution.ErrBlobMounted{
				From:       opts.Mount.From,
				Descriptor: distribution.Descriptor{Digest: d, Size: fi.Size()},
			}
		}
	}
	return s.r.newWriter(path, d.String(), opts.Mount.Stat.Size), nil
}

func (s *fileBlobStore) Resume(ctx context.Context, id string) (distribution.BlobWriter, error) {
	return nil, fmt.Errorf("unimplemented")
}

// fileWriter attempts to save blobs to disk in the appropriate location.
type fileWriter struct {
	driver   *fileRepository
	path     string
	uploadID string

	f *os.File

	closed    bool
	committed bool
	cancelled bool
	size      int64
	digest    godigest.Digest
	startedAt time.Time
}

func (d *fileRepository) newWriter(path, uploadID string, size int64) distribution.BlobWriter {
	return &fileWriter{
		driver:   d,
		path:     path,
		uploadID: uploadID,
		size:     size,
	}
}

func (d *fileRepository) newWriterTemp(path string, f *os.File) distribution.BlobWriter {
	return &fileWriter{
		driver:   d,
		path:     path,
		f:        f,
		uploadID: filepath.Base(f.Name()),
	}
}

func (w *fileWriter) ID() string {
	return w.uploadID
}

func (w *fileWriter) StartedAt() time.Time {
	return w.startedAt
}

func (w *fileWriter) ReadFrom(r io.Reader) (int64, error) {
	switch {
	case w.closed:
		return 0, fmt.Errorf("already closed")
	case w.committed:
		return 0, fmt.Errorf("already committed")
	case w.cancelled:
		return 0, fmt.Errorf("already cancelled")
	}
	if w.startedAt.IsZero() {
		w.startedAt = time.Now()
	}

	if w.f != nil {
		blobDigest, n, err := digestCopy(w.f, r)
		if err != nil {
			w.f.Close()
			return 0, err
		}
		if err := w.f.Close(); err != nil {
			return 0, err
		}
		name := w.f.Name()
		w.f = nil
		path := generateDigestPath(blobDigest.String(), w.path)
		if err := os.Rename(name, path); err != nil {
			return 0, err
		}
		w.size = n
		w.digest = blobDigest
		return n, err
	}

	if err := os.MkdirAll(filepath.Dir(w.path), 0755); err != nil {
		return 0, err
	}
	f, err := os.OpenFile(w.path+".download", os.O_TRUNC|os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0600)
	if err != nil {
		return 0, err
	}
	n, err := io.Copy(f, r)
	if err != nil {
		f.Close()
		return 0, err
	}
	if err := f.Close(); err != nil {
		return 0, err
	}
	if err := os.Rename(w.path+".download", w.path); err != nil {
		return 0, err
	}
	return n, nil
}

func (w *fileWriter) Write(p []byte) (int, error) {
	return 0, fmt.Errorf("already closed")
}

func (w *fileWriter) Size() int64 {
	return w.size
}

func (w *fileWriter) Close() error {
	switch {
	case w.closed:
		return fmt.Errorf("already closed")
	}
	w.closed = true
	if w.f != nil {
		w.f.Close()
		if err := os.Remove(w.f.Name()); err != nil && !os.IsNotExist(err) {
			return err
		}
		return nil
	}
	return nil
}

func (w *fileWriter) Cancel(ctx context.Context) error {
	switch {
	case w.closed:
		return fmt.Errorf("already closed")
	case w.committed:
		return fmt.Errorf("already committed")
	}
	w.cancelled = true
	return nil
}

// TODO: verify uploaded descriptor matches
func (w *fileWriter) Commit(ctx context.Context, descriptor distribution.Descriptor) (distribution.Descriptor, error) {
	desc := descriptor
	switch {
	case w.closed:
		return desc, fmt.Errorf("already closed")
	case w.committed:
		return desc, fmt.Errorf("already committed")
	case w.cancelled:
		return desc, fmt.Errorf("already cancelled")
	}
	w.committed = true
	if w.size > 0 {
		desc.Size = w.size
	}
	if len(w.digest) > 0 {
		desc.Digest = w.digest
	}
	return desc, nil
}

// digestCopy reads all of src into dst. It will return the sha256 sum of the
// stream (the blobDigest) or an error.
func digestCopy(dst io.Writer, src io.Reader) (blobDigest digest.Digest, n int64, err error) {
	algo := digest.Canonical
	blobhash := algo.Hash()
	n, err = io.Copy(io.MultiWriter(dst, blobhash), src)
	blobDigest = digest.NewDigestFromBytes(algo, blobhash.Sum(make([]byte, 0, blobhash.Size())))
	return blobDigest, n, err
}
