package registryclient

import (
	"context"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"sync"

	"github.com/distribution/distribution/v3"
	"github.com/distribution/distribution/v3/registry/api/errcode"
	"github.com/distribution/distribution/v3/registry/client"
	"github.com/distribution/distribution/v3/registry/client/auth"
	"github.com/opencontainers/go-digest"
	"github.com/openshift/library-go/pkg/image/reference"
	"k8s.io/klog/v2"

	distributionreference "github.com/distribution/distribution/v3/reference"
)

// AlternateBlobSourceStrategy is consulted when a repository cannot be reached to find alternate
// repositories that may be able to serve a given content-addressed blob. The strategy is consulted
// at most twice - once before any request is made to a given repository. If FirstRequest() returns a
// list of alternates, OnFailure is not invoked.
type AlternateBlobSourceStrategy interface {
	// FirstRequest returns the set of locations that should be searched in a preferred order. If locator
	// is not included in the response it will not be searched. If alternateRepositories is an empty list
	// no lookup will be performed and requests will exit with an error. If alternateRepositories is nil
	// and err is nil, OnFailure will be invoked if the first request fails.
	FirstRequest(ctx context.Context, locator reference.DockerImageReference) (alternateRepositories []reference.DockerImageReference, err error)
	// OnFailure is invoked if FirstRequest returned no error and a nil list of locations if and only if
	// an API call fails on the specified request. The result of alternateRepositories is cached for
	// subsequent calls to that repository.
	OnFailure(ctx context.Context, locator reference.DockerImageReference) (alternateRepositories []reference.DockerImageReference, err error)
}

// ManifestWithLocationService extends the ManifestService to allow clients to retrieve a manifest and
// get the location of the mirrored manifest. Not all ManifestServices returned from a Repository will
// support this interface and it must be conditional.
type ManifestWithLocationService interface {
	distribution.ManifestService

	// GetWithLocation returns the registry URL the provided manifest digest was retrieved from which may be Repository.Named(),
	// or one of the blob mirrors if alternate location for blob sources was provided. It returns an error if the digest could not be
	// located - if an error is returned the source reference (Repository.Named()) will be set.
	GetWithLocation(ctx context.Context, dgst digest.Digest, options ...distribution.ManifestServiceOption) (distribution.Manifest, reference.DockerImageReference, error)
}

// RepositoryWithLocation extends the Repository and allows clients to know which repository registry this talks to
// as primary (as a complement to Named() which does not include the URL).
type RepositoryWithLocation interface {
	distribution.Repository

	// Ref returns the DockerImageReference representing this repository.
	Ref() reference.DockerImageReference
}

// blobMirroredRepoRetriever allows a caller to retrieve a distribution.Repository. It may perform
// requests to authorize the client and will return an error if it fails.
type blobMirroredRepoRetriever interface {
	connectToRegistry(context.Context, repositoryLocator, bool) (RepositoryWithLocation, error)
}

// repositoryLocator caches the components necessary to connect to a single image repository.
type repositoryLocator struct {
	// ref is the full image reference as it is provided by the client.
	ref reference.DockerImageReference
	// url may specify a default protocol (http) instead of (https), but is otherwise calculated
	// by taking ref.Registry and applying it to url.Host
	url *url.URL
	// named is the image repository path on the server (namespace and name in ref terms) and is
	// required for the distribution registryclient.
	named distributionreference.Named
}

// blobMirroredRepository provides failover lookup behavior for blobs in a given repository on
// errors by delegating to the provided strategy for the first request or when a failure occurs.
// The strategy is expected to return a set of alternate locations to consume content from,
// which may not include the original source. Only requests made for content addressable blobs
// may be consulted in this fashion (anything via digest) - everything else must use source().
type blobMirroredRepository struct {
	locator  repositoryLocator
	insecure bool

	strategy  AlternateBlobSourceStrategy
	retriever blobMirroredRepoRetriever

	lock  sync.Mutex
	order []reference.DockerImageReference
	repos map[reference.DockerImageReference]RepositoryWithLocation
}

// Named returns the name of the repository.
func (r *blobMirroredRepository) Named() distributionreference.Named {
	return r.locator.named
}

// Named returns the name of the repository.
func (r *blobMirroredRepository) Ref() reference.DockerImageReference {
	return r.locator.ref
}

// Manifests wraps the manifest service in a blobMirroredManifest for shared retries.
func (r *blobMirroredRepository) Manifests(ctx context.Context, options ...distribution.ManifestServiceOption) (distribution.ManifestService, error) {
	return &blobMirroredManifest{repo: r, options: options}, nil
}

// Blobs wraps the blob service in a blobMirroredBlobstore for shared retries.
func (r *blobMirroredRepository) Blobs(ctx context.Context) distribution.BlobStore {
	return blobMirroredBlobstore{repo: r}
}

// Tags lists the tags under the named repository.
func (r *blobMirroredRepository) Tags(ctx context.Context) distribution.TagService {
	return blobMirroredTags{repo: r}
}

var (
	errNoValidAlternates = fmt.Errorf("no valid alterative sources for this content located")
	errNoValidSource     = fmt.Errorf("no source repository defined for accessing the repository")
)

// initialRepos returns a list of locations to attempt to access, a boolean indicating that alternates
// were suggested, or an error.
func (r *blobMirroredRepository) initialRepos(ctx context.Context) ([]reference.DockerImageReference, bool, error) {
	if r.strategy == nil {
		return []reference.DockerImageReference{r.locator.ref}, false, nil
	}

	// protect FirstRequest being called only one at a time and r.order writes
	r.lock.Lock()
	defer r.lock.Unlock()
	if r.order != nil {
		return r.order, true, nil
	}
	alternates, err := r.strategy.FirstRequest(ctx, r.locator.ref)
	if err != nil {
		return nil, false, err
	}
	if len(alternates) == 0 {
		return []reference.DockerImageReference{r.locator.ref}, false, nil
	}
	r.order = alternates
	return r.order, len(alternates) > 0, nil
}

// errorRepos returns a list of alternate registries to search for the provided content.
func (r *blobMirroredRepository) errorRepos(ctx context.Context) ([]reference.DockerImageReference, error) {
	if r.strategy == nil {
		return nil, nil
	}

	// TODO: potentially filter certain types of errors, maybe even per method type, if we ever
	// retry non-idempotent operations
	// protect OnFailure being called one at a time and r.order writes
	r.lock.Lock()
	defer r.lock.Unlock()
	if r.order != nil {
		return nil, nil
	}
	alternates, err := r.strategy.OnFailure(ctx, r.locator.ref)
	if err != nil {
		return nil, err
	}
	r.order = alternates
	return r.order, nil
}

// attemptRepos will invoke fn on all repos until fn returns no error. fn is expected to be idempotent.
func (r *blobMirroredRepository) attemptRepos(ctx context.Context, repos []reference.DockerImageReference, fn func(r RepositoryWithLocation) error) error {
	var firstErr error
	for _, ref := range repos {
		klog.V(5).Infof("Attempting to connect to %s", ref)
		repo, err := r.connect(ctx, ref)
		if err != nil {
			if firstErr == nil {
				firstErr = err
			}
			continue
		}
		if err := fn(repo); err != nil {
			if firstErr == nil {
				firstErr = err
			}
			continue
		}
		return nil
	}
	return firstErr
}

// isRequestError reports whether the registry rejected the request or was not
// able to serve any data. True for an error from io.Copy indicates that
// nothing was copied.
func isRequestError(err error) bool {
	var errorCode errcode.ErrorCode
	responseError := &client.UnexpectedHTTPResponseError{}
	statusError := &client.UnexpectedHTTPStatusError{}
	return errors.As(err, &errcode.Errors{}) ||
		errors.As(err, &errcode.Error{}) ||
		errors.As(err, &errorCode) ||
		errors.As(err, &responseError) ||
		errors.As(err, &statusError) ||
		errors.Is(err, auth.ErrNoBasicAuthCredentials)
}

// attemptFirstConnectedRepo will invoke fn on the first repo that successfully connects.
func (r *blobMirroredRepository) attemptFirstConnectedRepo(ctx context.Context, repos []reference.DockerImageReference, fn func(r RepositoryWithLocation) error) error {
	var firstErr error
	for _, ref := range repos {
		klog.V(5).Infof("Attempting to connect to %s", ref)
		repo, err := r.connect(ctx, ref)
		if err != nil {
			if firstErr == nil {
				firstErr = err
			}
			continue
		}
		if err := fn(repo); err != nil {
			if !isRequestError(err) {
				// The request may have been halfway through
				// and we cannot make another attempt.
				return err
			}
			// The registry replied with an error like 4xx or 5xx,
			// i.e. it hasn't served the blob and we can make
			// another attempt with a different registry.
			if firstErr == nil {
				firstErr = err
			}
			continue
		}
		return nil
	}
	return firstErr
}

// alternates accesses the set of repositories that may be valid alternatives for accessing content
func (r *blobMirroredRepository) alternates(ctx context.Context, fn func(r RepositoryWithLocation) error) error {
	repos, loaded, err := r.initialRepos(ctx)
	if err != nil {
		return err
	}
	if attemptErr := r.attemptRepos(ctx, repos, fn); attemptErr != nil {
		if loaded {
			return attemptErr
		}
		alternates, err := r.errorRepos(ctx)
		if err != nil {
			return err
		}
		if len(alternates) == 0 {
			return attemptErr
		}
		if alternateErr := r.attemptRepos(ctx, alternates, fn); alternateErr != nil {
			return attemptErr
		}
	}
	return nil
}

// firstConnectedAlternate invokes fn on the first alternate that can be connected to. Use when the
// function can only be invoked once (such as a method with side effects, like ServeBlob which writes
// to the response).
func (r *blobMirroredRepository) firstConnectedAlternate(ctx context.Context, fn func(r RepositoryWithLocation) error) error {
	repos, loaded, err := r.initialRepos(ctx)
	if err != nil {
		return err
	}
	if len(repos) == 0 {
		return errNoValidAlternates
	}
	if attemptErr := r.attemptFirstConnectedRepo(ctx, repos, fn); attemptErr != nil {
		if loaded {
			return attemptErr
		}
		alternates, err := r.errorRepos(ctx)
		if err != nil {
			return err
		}
		if alternateErr := r.attemptFirstConnectedRepo(ctx, alternates, fn); alternateErr != nil {
			return attemptErr
		}
	}
	return nil
}

// source connects to the original repository or returns an error. It will always use the same value
// of insecure as the original repository. Use when the request should only go to the initial repo.
func (r *blobMirroredRepository) source(ctx context.Context, fn func(r distribution.Repository) error) error {
	repo, err := r.connect(ctx, r.locator.ref)
	if err != nil {
		return err
	}
	return fn(repo)
}

// connect reuses or creates a connection to the provided reference, returning a repository instance
// or an error. This method expects that the connection only talks to the provided registry.
func (r *blobMirroredRepository) connect(ctx context.Context, ref reference.DockerImageReference) (RepositoryWithLocation, error) {
	r.lock.Lock()
	defer r.lock.Unlock()

	repo, ok := r.repos[ref]
	if ok {
		return repo, nil
	}
	locator := repositoryLocator{
		ref: ref,
	}
	repo, err := r.retriever.connectToRegistry(ctx, locator, ref != r.locator.ref || r.insecure)
	if err != nil {
		return nil, err
	}
	if r.repos == nil {
		r.repos = make(map[reference.DockerImageReference]RepositoryWithLocation)
	}
	r.repos[ref] = repo
	return repo, nil
}

// blobMirroredManifest will sequentially retry manifest operations on a set of repositories determined
// by the repository list, caching manifest services locally as needed (manifest service is assumed
// to have local state and does so in the registry client). The individual manifest service is not
// thread safe, but methods on this interface are thread safe.
type blobMirroredManifest struct {
	repo    *blobMirroredRepository
	options []distribution.ManifestServiceOption

	lock  sync.Mutex
	cache map[distribution.Repository]distribution.ManifestService
}

var _ distribution.ManifestService = &blobMirroredManifest{}
var _ ManifestWithLocationService = &blobMirroredManifest{}

// init retrieves or caches a manifets service for the provided repository, since each manifest
// service has local state.
func (f *blobMirroredManifest) init(ctx context.Context, r distribution.Repository) (distribution.ManifestService, error) {
	f.lock.Lock()
	defer f.lock.Unlock()

	ms := f.cache[r]
	if ms != nil {
		return ms, nil
	}
	ms, err := r.Manifests(ctx, f.options...)
	if err != nil {
		return nil, err
	}
	if f.cache == nil {
		f.cache = make(map[distribution.Repository]distribution.ManifestService)
	}
	f.cache[r] = ms
	return ms, nil
}

// alternates invokes fn once per alternate repo until fn returns without error.
func (f *blobMirroredManifest) alternates(ctx context.Context, fn func(m distribution.ManifestService, repo RepositoryWithLocation) error) error {
	return f.repo.alternates(ctx, func(repo RepositoryWithLocation) error {
		ms, err := f.init(ctx, repo)
		if err != nil {
			return err
		}
		return fn(ms, repo)
	})
}

// source invokes fn against the primary location.
func (f *blobMirroredManifest) source(ctx context.Context, fn func(r distribution.ManifestService) error) error {
	return f.repo.source(ctx, func(r distribution.Repository) error {
		ms, err := f.init(ctx, r)
		if err != nil {
			return err
		}
		return fn(ms)
	})
}

func (f *blobMirroredManifest) Put(ctx context.Context, manifest distribution.Manifest, options ...distribution.ManifestServiceOption) (digest.Digest, error) {
	var dgst digest.Digest
	err := f.source(ctx, func(r distribution.ManifestService) error {
		var err error
		dgst, err = r.Put(ctx, manifest, options...)
		return err
	})
	return dgst, err
}

func (f *blobMirroredManifest) Delete(ctx context.Context, dgst digest.Digest) error {
	return f.source(ctx, func(r distribution.ManifestService) error {
		return r.Delete(ctx, dgst)
	})
}

func (f *blobMirroredManifest) Exists(ctx context.Context, dgst digest.Digest) (bool, error) {
	var ok bool
	err := f.alternates(ctx, func(m distribution.ManifestService, repo RepositoryWithLocation) error {
		var err error
		ok, err = m.Exists(ctx, dgst)
		return err
	})
	return ok, err
}

func (f *blobMirroredManifest) Get(ctx context.Context, dgst digest.Digest, options ...distribution.ManifestServiceOption) (distribution.Manifest, error) {
	var manifest distribution.Manifest
	err := f.alternates(ctx, func(m distribution.ManifestService, repo RepositoryWithLocation) error {
		var err error
		manifest, err = m.Get(ctx, dgst, options...)
		klog.V(5).Infof("get manifest for %s served from %#v: %v", dgst, m, err)
		return err
	})
	return manifest, err
}

func (f *blobMirroredManifest) GetWithLocation(ctx context.Context, dgst digest.Digest, options ...distribution.ManifestServiceOption) (distribution.Manifest, reference.DockerImageReference, error) {
	var manifest distribution.Manifest
	var ref = f.repo.locator.ref
	err := f.alternates(ctx, func(m distribution.ManifestService, repo RepositoryWithLocation) error {
		var err error
		manifest, err = m.Get(ctx, dgst, options...)
		klog.V(5).Infof("get manifest for %s served from %#v: %v", dgst, m, err)
		if err == nil {
			ref = repo.Ref()
		}
		return err
	})
	return manifest, ref, err
}

// blobMirroredBlobstore wraps the blob store and invokes retries on the repo.
type blobMirroredBlobstore struct {
	repo *blobMirroredRepository
}

var _ distribution.BlobService = blobMirroredBlobstore{}

func (f blobMirroredBlobstore) Get(ctx context.Context, dgst digest.Digest) ([]byte, error) {
	var data []byte
	err := f.repo.alternates(ctx, func(r RepositoryWithLocation) error {
		var err error
		data, err = r.Blobs(ctx).Get(ctx, dgst)
		klog.V(5).Infof("get for %s served from %s: %v", dgst, r.Named(), err)
		return err
	})
	return data, err
}

func (f blobMirroredBlobstore) Stat(ctx context.Context, dgst digest.Digest) (distribution.Descriptor, error) {
	var desc distribution.Descriptor
	err := f.repo.alternates(ctx, func(r RepositoryWithLocation) error {
		var err error
		desc, err = r.Blobs(ctx).Stat(ctx, dgst)
		klog.V(5).Infof("stat for %s served from %s: %v", dgst, r.Named(), err)
		return err
	})
	return desc, err
}

func (f blobMirroredBlobstore) ServeBlob(ctx context.Context, w http.ResponseWriter, req *http.Request, dgst digest.Digest) error {
	err := f.repo.firstConnectedAlternate(ctx, func(r RepositoryWithLocation) error {
		err := r.Blobs(ctx).ServeBlob(ctx, w, req, dgst)
		klog.V(5).Infof("blob %s served from %s: %v", dgst, r.Named(), err)
		return err
	})
	return err
}

func (f blobMirroredBlobstore) Open(ctx context.Context, dgst digest.Digest) (io.ReadSeekCloser, error) {
	var rsc io.ReadSeekCloser
	err := f.repo.alternates(ctx, func(r RepositoryWithLocation) error {
		var err error
		rsc, err = r.Blobs(ctx).Open(ctx, dgst)
		if err != nil {
			klog.V(5).Infof("open %s from %s: %v", dgst, r.Named(), err)
			return err
		}

		// Distribution's implementation of Open doesn't send any requests to
		// the registry. We need the reader to send a request to see if the
		// registry can serve the blob.
		_, err = rsc.Read([]byte{})
		klog.V(5).Infof("open (read) %s from %s: %v", dgst, r.Named(), err)
		return err
	})
	return rsc, err
}

func (f blobMirroredBlobstore) Create(ctx context.Context, options ...distribution.BlobCreateOption) (distribution.BlobWriter, error) {
	var bw distribution.BlobWriter
	err := f.repo.source(ctx, func(r distribution.Repository) error {
		var err error
		bw, err = r.Blobs(ctx).Create(ctx, options...)
		return err
	})
	return bw, err
}

func (f blobMirroredBlobstore) Put(ctx context.Context, mediaType string, p []byte) (distribution.Descriptor, error) {
	var desc distribution.Descriptor
	err := f.repo.source(ctx, func(r distribution.Repository) error {
		var err error
		desc, err = r.Blobs(ctx).Put(ctx, mediaType, p)
		return err
	})
	return desc, err
}

func (f blobMirroredBlobstore) Resume(ctx context.Context, id string) (distribution.BlobWriter, error) {
	var bw distribution.BlobWriter
	err := f.repo.source(ctx, func(r distribution.Repository) error {
		var err error
		bw, err = r.Blobs(ctx).Resume(ctx, id)
		return err
	})
	return bw, err
}

func (f blobMirroredBlobstore) Delete(ctx context.Context, dgst digest.Digest) error {
	return f.repo.source(ctx, func(r distribution.Repository) error {
		return r.Blobs(ctx).Delete(ctx, dgst)
	})
}

// blobMirroredTags lazily accesses the source repository
type blobMirroredTags struct {
	repo *blobMirroredRepository
}

var _ distribution.TagService = blobMirroredTags{}

func (f blobMirroredTags) Get(ctx context.Context, tag string) (distribution.Descriptor, error) {
	var desc distribution.Descriptor
	err := f.repo.source(ctx, func(r distribution.Repository) error {
		var err error
		desc, err = r.Tags(ctx).Get(ctx, tag)
		return err
	})
	return desc, err
}

func (f blobMirroredTags) All(ctx context.Context) ([]string, error) {
	var tags []string
	err := f.repo.source(ctx, func(r distribution.Repository) error {
		var err error
		tags, err = r.Tags(ctx).All(ctx)
		return err
	})
	return tags, err
}

func (f blobMirroredTags) Lookup(ctx context.Context, digest distribution.Descriptor) ([]string, error) {
	var tags []string
	err := f.repo.source(ctx, func(r distribution.Repository) error {
		var err error
		tags, err = r.Tags(ctx).Lookup(ctx, digest)
		return err
	})
	return tags, err
}

func (f blobMirroredTags) Tag(ctx context.Context, tag string, desc distribution.Descriptor) error {
	return f.repo.source(ctx, func(r distribution.Repository) error {
		return r.Tags(ctx).Tag(ctx, tag, desc)
	})
}

func (f blobMirroredTags) Untag(ctx context.Context, tag string) error {
	return f.repo.source(ctx, func(r distribution.Repository) error {
		return r.Tags(ctx).Untag(ctx, tag)
	})
}
