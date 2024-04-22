package manifest

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"regexp"
	"runtime"
	"sync"

	"github.com/spf13/pflag"

	"github.com/distribution/distribution/v3"
	"github.com/distribution/distribution/v3/manifest/manifestlist"
	"github.com/distribution/distribution/v3/manifest/ocischema"
	"github.com/distribution/distribution/v3/manifest/schema1"
	"github.com/distribution/distribution/v3/manifest/schema2"
	"github.com/distribution/distribution/v3/reference"
	"github.com/distribution/distribution/v3/registry/api/errcode"
	v2 "github.com/distribution/distribution/v3/registry/api/v2"

	"github.com/docker/libtrust"
	"github.com/opencontainers/go-digest"
	"k8s.io/client-go/rest"
	"k8s.io/klog/v2"

	imagespecv1 "github.com/opencontainers/image-spec/specs-go/v1"
	"github.com/openshift/library-go/pkg/image/dockerv1client"
	imagereference "github.com/openshift/library-go/pkg/image/reference"
	"github.com/openshift/library-go/pkg/image/registryclient"
	"github.com/openshift/oc/pkg/cli/image/manifest/dockercredentials"
	"github.com/openshift/oc/pkg/helpers/image/dockerlayer/add"
)

type ParallelOptions struct {
	MaxPerRegistry int
}

func (o *ParallelOptions) Bind(flags *pflag.FlagSet) {
	flags.IntVar(&o.MaxPerRegistry, "max-per-registry", o.MaxPerRegistry, "Number of concurrent requests allowed per registry.")
}

type SecurityOptions struct {
	RegistryConfig   string
	Insecure         bool
	SkipVerification bool
	CAData           string

	CachedContext *registryclient.Context
}

func (o *SecurityOptions) Bind(flags *pflag.FlagSet) {
	// TODO: remove REGISTRY_AUTH_PREFERENCE env variable support and support only podman in 4.15
	flags.StringVarP(&o.RegistryConfig, "registry-config", "a", o.RegistryConfig, "Path to your registry credentials. Alternatively REGISTRY_AUTH_FILE env variable can be also specified. Defaults to ${XDG_RUNTIME_DIR}/containers/auth.json, /run/containers/${UID}/auth.json, ${XDG_CONFIG_HOME}/containers/auth.json, ${DOCKER_CONFIG}, ~/.docker/config.json, ~/.dockercfg. The order can be changed via the REGISTRY_AUTH_PREFERENCE env variable (deprecated) to a \"docker\" value to prioritizes Docker credentials over Podman's.")
	flags.BoolVar(&o.Insecure, "insecure", o.Insecure, "Allow push and pull operations to registries to be made over HTTP")
	flags.BoolVar(&o.SkipVerification, "skip-verification", o.SkipVerification, "Skip verifying the integrity of the retrieved content. This is not recommended, but may be necessary when importing images from older image registries. Only bypass verification if the registry is known to be trustworthy.")
	flags.StringVar(&o.CAData, "certificate-authority", o.CAData, "The path to a certificate authority bundle to use when communicating with the managed container image registries. If --insecure is used, this flag will be ignored. ")
}

// ReferentialHTTPClient returns an http.Client that is appropriate for accessing
// blobs referenced outside of the registry (due to the present of the URLs attribute
// in the manifest reference for a layer).
func (o *SecurityOptions) ReferentialHTTPClient() (*http.Client, error) {
	ctx, err := o.Context()
	if err != nil {
		return nil, err
	}
	client := &http.Client{}
	if o.Insecure {
		client.Transport = ctx.InsecureTransport
	} else {
		client.Transport = ctx.Transport
	}
	return client, nil
}

type Verifier interface {
	Verify(dgst, contentDgst digest.Digest)
	Verified() bool
}

func NewVerifier() Verifier {
	return &verifier{}
}

type verifier struct {
	lock     sync.Mutex
	hadError bool
}

func (v *verifier) Verify(dgst, contentDgst digest.Digest) {
	if contentDgst == dgst {
		return
	}
	v.lock.Lock()
	defer v.lock.Unlock()
	v.hadError = true
}

func (v *verifier) Verified() bool {
	v.lock.Lock()
	defer v.lock.Unlock()
	return !v.hadError
}

func (o *SecurityOptions) Context() (*registryclient.Context, error) {
	if o.CachedContext != nil {
		return o.CachedContext, nil
	}
	ctx, err := o.NewContext()
	if err == nil {
		o.CachedContext = ctx
		o.CachedContext.Retries = 3
	}
	return ctx, err
}

func (o *SecurityOptions) NewContext() (*registryclient.Context, error) {
	userAgent := rest.DefaultKubernetesUserAgent()
	var rt http.RoundTripper
	var err error
	if len(o.CAData) > 0 {
		cadata, err := os.ReadFile(o.CAData)
		if err != nil {
			return nil, fmt.Errorf("failed to read registry ca bundle: %v", err)
		}

		rt, err = rest.TransportFor(&rest.Config{UserAgent: userAgent, TLSClientConfig: rest.TLSClientConfig{CAData: cadata}})
		if err != nil {
			return nil, err
		}

	} else {
		rt, err = rest.TransportFor(&rest.Config{UserAgent: userAgent})
		if err != nil {
			return nil, err
		}
	}
	insecureRT, err := rest.TransportFor(&rest.Config{TLSClientConfig: rest.TLSClientConfig{Insecure: true}, UserAgent: userAgent})
	if err != nil {
		return nil, err
	}
	credStoreFactory, err := dockercredentials.NewCredentialStoreFactory(o.RegistryConfig)
	if err != nil {
		if len(o.RegistryConfig) > 0 {
			return nil, fmt.Errorf("unable to load --registry-config: %v", err)
		}
		return nil, err
	}
	ctx := registryclient.NewContext(rt, insecureRT).WithCredentialsFactory(credStoreFactory)
	ctx.DisableDigestVerification = o.SkipVerification
	return ctx, nil
}

// FilterOptions assist in filtering out unneeded manifests from ManifestList objects.
type FilterOptions struct {
	FilterByOS      string
	DefaultOSFilter bool
	OSFilter        *regexp.Regexp
}

// Bind adds the options to the flag set.
func (o *FilterOptions) Bind(flags *pflag.FlagSet) {
	flags.StringVar(&o.FilterByOS, "filter-by-os", o.FilterByOS, "A regular expression to control which images are considered when multiple variants are available. Images will be passed as '<platform>/<architecture>[/<variant>]'.")
}

// Validate checks whether the flags are ready for use.
func (o *FilterOptions) Validate() error {
	pattern := o.FilterByOS
	if len(pattern) > 0 {
		re, err := regexp.Compile(pattern)
		if err != nil {
			return fmt.Errorf("--filter-by-os was not a valid regular expression: %v", err)
		}
		o.OSFilter = re
	}
	return nil
}

// Complete performs defaulting by OS.
func (o *FilterOptions) Complete(flags *pflag.FlagSet) error {
	pattern := o.FilterByOS
	if len(pattern) == 0 && !flags.Changed("filter-by-os") {
		o.DefaultOSFilter = true
		o.FilterByOS = regexp.QuoteMeta(fmt.Sprintf("%s/%s", runtime.GOOS, runtime.GOARCH))
	}
	return nil
}

// IsWildcardFilter returns true if the filter regex is set to a wildcard
func (o *FilterOptions) IsWildcardFilter() bool {
	wildcardFilter := ".*"
	if o.FilterByOS == wildcardFilter {
		return true
	}
	return false
}

// Include returns true if the provided manifest should be included, or the first image if the user didn't alter the
// default selection and there is only one image.
func (o *FilterOptions) Include(d *manifestlist.ManifestDescriptor, hasMultiple bool) bool {
	if o.OSFilter == nil {
		return true
	}
	if o.DefaultOSFilter && !hasMultiple {
		return true
	}
	s := PlatformSpecString(d.Platform)
	return o.OSFilter.MatchString(s)
}

func PlatformSpecString(platform manifestlist.PlatformSpec) string {
	if len(platform.Variant) > 0 {
		return fmt.Sprintf("%s/%s/%s", platform.OS, platform.Architecture, platform.Variant)
	}
	return fmt.Sprintf("%s/%s", platform.OS, platform.Architecture)
}

// IncludeAll returns true if the provided manifest matches the filter, or all if there was no filter.
func (o *FilterOptions) IncludeAll(d *manifestlist.ManifestDescriptor, hasMultiple bool) bool {
	if o.OSFilter == nil {
		return true
	}
	s := PlatformSpecString(d.Platform)
	return o.OSFilter.MatchString(s)
}

type FilterFunc func(*manifestlist.ManifestDescriptor, bool) bool

// PreferManifestList specifically requests a manifest list first
var PreferManifestList = distribution.WithManifestMediaTypes([]string{
	manifestlist.MediaTypeManifestList,
	imagespecv1.MediaTypeImageIndex,
	schema2.MediaTypeManifest,
	imagespecv1.MediaTypeImageManifest,
})

// IsManifestList returns if a given image is a manifestlist or not
func IsManifestList(ctx context.Context, from imagereference.DockerImageReference, repo distribution.Repository) (bool, error) {
	var srcDigest digest.Digest
	if len(from.ID) > 0 {
		srcDigest = digest.Digest(from.ID)
	} else if len(from.Tag) > 0 {
		desc, err := repo.Tags(ctx).Get(ctx, from.Tag)
		if err != nil {
			return false, err
		}
		srcDigest = desc.Digest
	} else {
		return false, fmt.Errorf("no tag or digest specified")
	}
	manifests, err := repo.Manifests(ctx)
	if err != nil {
		return false, err
	}
	srcManifest, err := manifests.Get(ctx, srcDigest, PreferManifestList)
	if err != nil {
		return false, err
	}

	_, ok := srcManifest.(*manifestlist.DeserializedManifestList)
	return ok, nil
}

// AllManifests returns all non-list manifests, the list manifest (if any), the digest the from refers to, or an error.
func AllManifests(ctx context.Context, from imagereference.DockerImageReference, repo distribution.Repository) (map[digest.Digest]distribution.Manifest, *manifestlist.DeserializedManifestList, digest.Digest, error) {
	var srcDigest digest.Digest
	if len(from.ID) > 0 {
		srcDigest = digest.Digest(from.ID)
	} else if len(from.Tag) > 0 {
		desc, err := repo.Tags(ctx).Get(ctx, from.Tag)
		if err != nil {
			return nil, nil, "", err
		}
		srcDigest = desc.Digest
	} else {
		return nil, nil, "", fmt.Errorf("no tag or digest specified")
	}
	manifests, err := repo.Manifests(ctx)
	if err != nil {
		return nil, nil, "", err
	}
	srcManifest, err := manifests.Get(ctx, srcDigest, PreferManifestList)
	if err != nil {
		return nil, nil, "", err
	}

	return ManifestsFromList(ctx, srcDigest, srcManifest, manifests, from)
}

type ManifestLocation struct {
	Manifest     digest.Digest
	ManifestList digest.Digest
}

func (m ManifestLocation) IsList() bool {
	return len(m.ManifestList) > 0
}

func (m ManifestLocation) ManifestListDigest() digest.Digest {
	if m.IsList() {
		return m.ManifestList
	}
	return ""
}

func (m ManifestLocation) String() string {
	if m.IsList() {
		return fmt.Sprintf("manifest %s in manifest list %s", m.Manifest, m.ManifestList)
	}
	return fmt.Sprintf("manifest %s", m.Manifest)
}

// FirstManifest returns the first manifest at the request location that matches the filter function.
func FirstManifest(ctx context.Context, from imagereference.DockerImageReference, repo distribution.Repository, filterFn FilterFunc) (distribution.Manifest, ManifestLocation, error) {
	var srcDigest digest.Digest
	if len(from.ID) > 0 {
		srcDigest = digest.Digest(from.ID)
	} else if len(from.Tag) > 0 {
		desc, err := repo.Tags(ctx).Get(ctx, from.Tag)
		if err != nil {
			return nil, ManifestLocation{}, err
		}
		srcDigest = desc.Digest
	} else {
		return nil, ManifestLocation{}, fmt.Errorf("no tag or digest specified")
	}
	manifests, err := repo.Manifests(ctx)
	if err != nil {
		return nil, ManifestLocation{}, err
	}
	srcManifest, err := manifests.Get(ctx, srcDigest, PreferManifestList)
	if err != nil {
		return nil, ManifestLocation{}, err
	}

	originalSrcDigest := srcDigest
	srcChildren, srcManifest, srcDigest, err := ProcessManifestList(ctx, srcDigest, srcManifest, manifests, from, filterFn, false)
	if err != nil {
		return nil, ManifestLocation{}, err
	}
	if srcManifest == nil {
		return nil, ManifestLocation{}, AllImageFilteredErr
	}
	if len(srcChildren) > 0 {
		// More than one match within the list, return the first
		childManifest := srcChildren[0]
		childDigest, err := registryclient.ContentDigestForManifest(childManifest, srcDigest.Algorithm())
		if err != nil {
			return nil, ManifestLocation{}, fmt.Errorf("could not generate digest for first manifest")
		}
		return childManifest, ManifestLocation{Manifest: childDigest, ManifestList: originalSrcDigest}, nil
	}
	if srcDigest != originalSrcDigest {
		// One match in list selected by ProcessManifestList
		return srcManifest, ManifestLocation{Manifest: srcDigest, ManifestList: originalSrcDigest}, nil
	}
	// Was not a list
	return srcManifest, ManifestLocation{Manifest: srcDigest}, nil
}

// ManifestToImageConfig takes an image manifest and converts it into a structured object.
func ManifestToImageConfig(ctx context.Context, srcManifest distribution.Manifest, blobs distribution.BlobService, location ManifestLocation) (*dockerv1client.DockerImageConfig, []distribution.Descriptor, error) {
	switch t := srcManifest.(type) {
	case *schema2.DeserializedManifest:
		if t.Config.MediaType != schema2.MediaTypeImageConfig {
			return nil, nil, fmt.Errorf("%s does not have the expected image configuration media type: %s", location, t.Config.MediaType)
		}
		configJSON, err := blobs.Get(ctx, t.Config.Digest)
		if err != nil {
			return nil, nil, fmt.Errorf("cannot retrieve image configuration for %s: %v", location, err)
		}
		klog.V(4).Infof("Raw image config json:\n%s", string(configJSON))
		config := &dockerv1client.DockerImageConfig{}
		if err := json.Unmarshal(configJSON, &config); err != nil {
			return nil, nil, fmt.Errorf("unable to parse image configuration: %v", err)
		}

		base := config
		layers := t.Layers
		base.Size = 0
		for _, layer := range t.Layers {
			base.Size += layer.Size
		}

		return base, layers, nil

	case *ocischema.DeserializedManifest:
		if t.Config.MediaType != imagespecv1.MediaTypeImageConfig {
			return nil, nil, fmt.Errorf("%s does not have the expected image configuration media type: %s", location, t.Config.MediaType)
		}
		configJSON, err := blobs.Get(ctx, t.Config.Digest)
		if err != nil {
			return nil, nil, fmt.Errorf("cannot retrieve image configuration for %s: %v", location, err)
		}
		klog.V(4).Infof("Raw image config json:\n%s", string(configJSON))
		config := &dockerv1client.DockerImageConfig{}
		if err := json.Unmarshal(configJSON, &config); err != nil {
			return nil, nil, fmt.Errorf("unable to parse image configuration: %v", err)
		}

		base := config
		layers := t.Layers
		base.Size = 0
		for _, layer := range t.Layers {
			base.Size += layer.Size
		}

		return base, layers, nil

	case *schema1.SignedManifest:
		if klog.V(4).Enabled() {
			_, configJSON, _ := srcManifest.Payload()
			klog.Infof("Raw image config json:\n%s", string(configJSON))
		}
		if len(t.History) == 0 {
			return nil, nil, fmt.Errorf("input image is in an unknown format: no v1Compatibility history")
		}
		config := &dockerv1client.DockerV1CompatibilityImage{}
		if err := json.Unmarshal([]byte(t.History[0].V1Compatibility), &config); err != nil {
			return nil, nil, err
		}

		base := &dockerv1client.DockerImageConfig{}
		if err := dockerv1client.Convert_DockerV1CompatibilityImage_to_DockerImageConfig(config, base); err != nil {
			return nil, nil, err
		}

		// schema1 layers are in reverse order
		layers := make([]distribution.Descriptor, 0, len(t.FSLayers))
		for i := len(t.FSLayers) - 1; i >= 0; i-- {
			layer := distribution.Descriptor{
				MediaType: schema2.MediaTypeLayer,
				Digest:    t.FSLayers[i].BlobSum,
				// size must be reconstructed from the blobs
			}
			// we must reconstruct the tar sum from the blobs
			add.AddLayerToConfig(base, layer, "")
			layers = append(layers, layer)
		}

		return base, layers, nil

	case *manifestlist.DeserializedManifestList:
		return nil, nil, fmt.Errorf("use --keep-manifest-list option for image manifest type %T from %s", srcManifest, location)
	default:
		return nil, nil, fmt.Errorf("unknown image manifest of type %T from %s", srcManifest, location)
	}
}

func ProcessManifestList(ctx context.Context, srcDigest digest.Digest, srcManifest distribution.Manifest, manifests distribution.ManifestService, ref imagereference.DockerImageReference, filterFn FilterFunc, keepManifestList bool) (children []distribution.Manifest, manifest distribution.Manifest, digest digest.Digest, err error) {
	var childManifests []distribution.Manifest
	switch t := srcManifest.(type) {
	case *manifestlist.DeserializedManifestList:
		manifestDigest := srcDigest
		manifestList := t

		filtered := make([]manifestlist.ManifestDescriptor, 0, len(t.Manifests))
		for _, manifest := range t.Manifests {
			if !filterFn(&manifest, len(t.Manifests) > 1) {
				klog.V(5).Infof("Skipping image %s for %#v from %s", manifest.Digest, manifest.Platform, ref)
				continue
			}
			klog.V(5).Infof("Including image %s for %#v from %s", manifest.Digest, manifest.Platform, ref)
			filtered = append(filtered, manifest)
		}

		if len(filtered) == 0 && !keepManifestList {
			return nil, nil, "", nil
		}

		// if we're not keeping manifest lists and this one has been filtered, make a new one with
		// just the filtered platforms.
		if len(filtered) != len(t.Manifests) && !keepManifestList {
			var err error
			t, err = manifestlist.FromDescriptors(filtered)
			if err != nil {
				return nil, nil, "", fmt.Errorf("unable to filter source image %s manifest list: %v", ref, err)
			}
			_, body, err := t.Payload()
			if err != nil {
				return nil, nil, "", fmt.Errorf("unable to filter source image %s manifest list (bad payload): %v", ref, err)
			}
			manifestList = t
			manifestDigest, err = registryclient.ContentDigestForManifest(t, srcDigest.Algorithm())
			if err != nil {
				return nil, nil, "", err
			}
			klog.V(5).Infof("Filtered manifest list to new digest %s:\n%s", manifestDigest, body)
		}

		for i, manifest := range filtered {
			childManifest, err := manifests.Get(ctx, manifest.Digest, PreferManifestList)
			if err != nil {
				return nil, nil, "", fmt.Errorf("unable to retrieve source image %s manifest #%d from manifest list: %v", ref, i+1, err)
			}
			childManifests = append(childManifests, childManifest)
		}

		switch {
		case len(childManifests) == 1 && !keepManifestList:
			// Just return the single platform specific image
			manifestDigest, err := registryclient.ContentDigestForManifest(childManifests[0], srcDigest.Algorithm())
			if err != nil {
				return nil, nil, "", err
			}
			klog.V(2).Infof("Chose %s/%s manifest from the manifest list.", t.Manifests[0].Platform.OS, t.Manifests[0].Platform.Architecture)
			return nil, childManifests[0], manifestDigest, nil
		default:
			return childManifests, manifestList, manifestDigest, nil
		}

	default:
		return nil, srcManifest, srcDigest, nil
	}
}

// ManifestsFromList returns a map of all image manifests for a given manifest. It returns the ManifestList and its digest if
// srcManifest is a list, or an error.
func ManifestsFromList(ctx context.Context, srcDigest digest.Digest, srcManifest distribution.Manifest, manifests distribution.ManifestService, ref imagereference.DockerImageReference) (map[digest.Digest]distribution.Manifest, *manifestlist.DeserializedManifestList, digest.Digest, error) {
	switch t := srcManifest.(type) {
	case *manifestlist.DeserializedManifestList:
		allManifests := make(map[digest.Digest]distribution.Manifest)
		manifestDigest := srcDigest
		manifestList := t

		for i, manifest := range t.Manifests {
			childManifest, err := manifests.Get(ctx, manifest.Digest, PreferManifestList)
			if err != nil {
				return nil, nil, "", fmt.Errorf("unable to retrieve source image %s manifest #%d from manifest list: %v", ref, i+1, err)
			}
			allManifests[manifest.Digest] = childManifest
		}

		return allManifests, manifestList, manifestDigest, nil

	default:
		return map[digest.Digest]distribution.Manifest{srcDigest: srcManifest}, nil, "", nil
	}
}

// TODO: Remove support for v2 schema in 4.9
func PutManifestInCompatibleSchema(
	ctx context.Context,
	srcManifest distribution.Manifest,
	tag string,
	toManifests distribution.ManifestService,
	ref reference.Named,
	blobs distribution.BlobService, // support schema2 -> schema1 downconversion
	configJSON []byte, // optional, if not passed blobs will be used
) (digest.Digest, error) {
	var options []distribution.ManifestServiceOption
	if len(tag) > 0 {
		klog.V(5).Infof("Put manifest %s:%s", ref, tag)
		options = []distribution.ManifestServiceOption{distribution.WithTag(tag)}
	} else {
		klog.V(5).Infof("Put manifest %s", ref)
	}
	switch t := srcManifest.(type) {
	case *schema1.SignedManifest:
		manifest, err := convertToSchema2(ctx, blobs, t)
		if err != nil {
			klog.V(2).Infof("Unable to convert manifest to schema2: %v", err)
			return toManifests.Put(ctx, t, distribution.WithTag(tag))
		}
		klog.Infof("warning: Digests are not preserved with schema version 1 images. Support for schema version 1 images will be removed in a future release")
		return toManifests.Put(ctx, manifest, options...)
	}

	toDigest, err := toManifests.Put(ctx, srcManifest, options...)
	if err == nil {
		return toDigest, nil
	}
	errs, ok := err.(errcode.Errors)
	if !ok || len(errs) == 0 {
		return toDigest, err
	}
	errCode, ok := errs[0].(errcode.Error)
	if !ok || errCode.ErrorCode() != v2.ErrorCodeManifestInvalid {
		return toDigest, err
	}
	// try downconverting to v2-schema1
	schema2Manifest, ok := srcManifest.(*schema2.DeserializedManifest)
	if !ok {
		return toDigest, err
	}
	tagRef, tagErr := reference.WithTag(ref, tag)
	if tagErr != nil {
		return toDigest, err
	}
	klog.V(5).Infof("Registry reported invalid manifest error, attempting to convert to v2schema1 as ref %s", tagRef)
	schema1Manifest, convertErr := convertToSchema1(ctx, blobs, configJSON, schema2Manifest, tagRef)
	if convertErr != nil {
		if klog.V(6).Enabled() {
			_, data, _ := schema2Manifest.Payload()
			klog.Infof("Input schema\n%s", string(data))
		}
		klog.V(2).Infof("Unable to convert manifest to schema1: %v", convertErr)
		return toDigest, err
	}
	if klog.V(6).Enabled() {
		_, data, _ := schema1Manifest.Payload()
		klog.Infof("Converted to v2schema1\n%s", string(data))
	}
	return toManifests.Put(ctx, schema1Manifest, distribution.WithTag(tag))
}

// convertToSchema2 attempts to build a v2 manifest from a v1 manifest, which requires reading blobs to get layer sizes.
// Requires the destination layers already exist in the target repository.
func convertToSchema2(ctx context.Context, blobs distribution.BlobService, srcManifest *schema1.SignedManifest) (distribution.Manifest, error) {
	if klog.V(6).Enabled() {
		klog.Infof("Up converting v1 schema image:\n%#v", srcManifest.Manifest)
	}

	config, layers, err := ManifestToImageConfig(ctx, srcManifest, blobs, ManifestLocation{})
	if err != nil {
		return nil, err
	}
	if klog.V(6).Enabled() {
		klog.Infof("Resulting schema: %#v", config)
	}
	// create synthetic history
	// TODO: create restored history?
	if len(config.History) == 0 {
		for i := len(config.History); i < len(layers); i++ {
			config.History = append(config.History, dockerv1client.DockerConfigHistory{
				Created: config.Created,
			})
		}
	}

	configJSON, err := json.Marshal(config)
	if err != nil {
		return nil, err
	}
	if klog.V(6).Enabled() {
		klog.Infof("Resulting config.json:\n%s", string(configJSON))
	}
	configDescriptor := distribution.Descriptor{
		Digest:    digest.FromBytes(configJSON),
		Size:      int64(len(configJSON)),
		MediaType: schema2.MediaTypeImageConfig,
	}
	b := schema2.NewManifestBuilder(configDescriptor, configJSON)
	_, err = blobs.Put(ctx, schema2.MediaTypeImageConfig, configJSON)
	if err != nil {
		return nil, err
	}
	for _, layer := range layers {
		desc, err := blobs.Stat(ctx, layer.Digest)
		if err != nil {
			return nil, err
		}
		desc.MediaType = schema2.MediaTypeLayer
		if err := b.AppendReference(desc); err != nil {
			return nil, err
		}
	}
	return b.Build(ctx)
}

// TODO: Remove support for v2 schema in 4.9
func convertToSchema1(ctx context.Context, blobs distribution.BlobService, configJSON []byte, schema2Manifest *schema2.DeserializedManifest, ref reference.Named) (distribution.Manifest, error) {
	if configJSON == nil {
		targetDescriptor := schema2Manifest.Target()
		config, err := blobs.Get(ctx, targetDescriptor.Digest)
		if err != nil {
			return nil, err
		}
		configJSON = config
	}
	trustKey, err := loadPrivateKey()
	if err != nil {
		return nil, err
	}
	if klog.V(6).Enabled() {
		klog.Infof("Down converting v2 schema image:\n%#v\n%s", schema2Manifest.Layers, configJSON)
	}
	builder := schema1.NewConfigManifestBuilder(blobs, trustKey, ref, configJSON)
	for _, d := range schema2Manifest.Layers {
		if err := builder.AppendReference(d); err != nil {
			return nil, err
		}
	}
	manifest, err := builder.Build(ctx)
	if err != nil {
		return nil, err
	}
	return manifest, nil
}

var (
	privateKeyLock sync.Mutex
	privateKey     libtrust.PrivateKey
)

// TODO: Remove support for v2 schema in 4.9
func loadPrivateKey() (libtrust.PrivateKey, error) {
	privateKeyLock.Lock()
	defer privateKeyLock.Unlock()
	if privateKey != nil {
		return privateKey, nil
	}
	trustKey, err := libtrust.GenerateECP256PrivateKey()
	if err != nil {
		return nil, err
	}
	privateKey = trustKey
	return privateKey, nil
}
