package release

import (
	"context"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strings"
	"sync"
	"time"

	digest "github.com/opencontainers/go-digest"
	operatorv1alpha1 "github.com/openshift/api/operator/v1alpha1"
	"github.com/spf13/cobra"
	corev1 "k8s.io/api/core/v1"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/cli-runtime/pkg/genericiooptions"
	corev1client "k8s.io/client-go/kubernetes/typed/core/v1"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/transport"
	"k8s.io/client-go/util/retry"
	"k8s.io/klog/v2"
	kcmdutil "k8s.io/kubectl/pkg/cmd/util"
	"k8s.io/kubectl/pkg/util/templates"
	"sigs.k8s.io/yaml"

	apicfgv1 "github.com/openshift/api/config/v1"
	imagev1 "github.com/openshift/api/image/v1"
	imageclient "github.com/openshift/client-go/image/clientset/versioned"
	"github.com/openshift/library-go/pkg/image/dockerv1client"
	imagereference "github.com/openshift/library-go/pkg/image/reference"
	"github.com/openshift/library-go/pkg/manifest"
	"github.com/openshift/library-go/pkg/verify"
	"github.com/openshift/library-go/pkg/verify/store/configmap"
	"github.com/openshift/library-go/pkg/verify/store/sigstore"
	"github.com/openshift/library-go/pkg/verify/util"
	"github.com/openshift/oc/pkg/cli/image/extract"
	"github.com/openshift/oc/pkg/cli/image/imagesource"
	imagemanifest "github.com/openshift/oc/pkg/cli/image/manifest"
	"github.com/openshift/oc/pkg/cli/image/mirror"
)

// configFilesBaseDir is created under '--to-dir', when specified, to contain release image
// signature files. It is not used when '--release-image-signature-to-dir` is specified
// which takes precedence over '--to-dir'.
const configFilesBaseDir = "config"

// maxDigestHashLen is used to truncate digest hash portion before using as part of
// signature file name.
const maxDigestHashLen = 16

// signatureFileNameFmt defines format of the release image signature file name.
const signatureFileNameFmt = "signature-%s-%s.json"

// instructionTypeICSP defines the printed out instruction type for ImageContentSourcePolicy.
const instructionTypeICSP = "icsp"

// instructionTypeIDMS defines the printed out instruction type for ImageDigestMirrorSet.
const instructionTypeIDMS = "idms"

// instructionTypeNone will not print out mirror instruction for ImageDigestMirrorSet or ImageContentSourcePolicy.
const instructionTypeNone = "none"

// archMap maps Go architecture strings to OpenShift supported values for any that differ.
var archMap = map[string]string{
	"amd64": "x86_64",
	"arm64": "aarch64",
}

// NewMirrorOptions creates the options for mirroring a release.
func NewMirrorOptions(streams genericiooptions.IOStreams) *MirrorOptions {
	return &MirrorOptions{
		IOStreams:       streams,
		ParallelOptions: imagemanifest.ParallelOptions{MaxPerRegistry: 6},
	}
}

// NewMirror creates a command to mirror an existing release.
//
// # Example command to mirror a release to a local repository to work offline
//
//	$ oc adm release mirror \
//	    --from=registry.ci.openshift.org/openshift/v4.11 \
//	    --to=mycompany.com/myrepository/repo
func NewMirror(f kcmdutil.Factory, streams genericiooptions.IOStreams) *cobra.Command {
	o := NewMirrorOptions(streams)
	cmd := &cobra.Command{
		Use:   "mirror",
		Short: "Mirror a release to a different image registry location",
		Long: templates.LongDesc(`
			Mirror an OpenShift release image to another registry and produce a configuration
			manifest containing the release image signature.

			Copies the images and update payload for a given release from one registry to another.
			By default this command will not alter the payload and will print out the configuration
			that must be applied to a cluster to use the mirror, but you may opt to rewrite the
			update to point to the new location and lose the cryptographic integrity of the update.

			Creates a release image signature config map that can be saved to a directory, applied
			directly to a connected cluster, or both.

			The common use for this command is to mirror a specific OpenShift release version to
			a private registry and create a signature config map for use in a disconnected or
			offline context. The command copies all images that are part of a release into the
			target repository and then prints the correct information to give to OpenShift to use
			that content offline. An alternate mode is to specify --to-image-stream, which imports
			the images directly into an OpenShift image stream.

			You may use --to-dir to specify a directory to download release content into, and add
			the file:// prefix to the --to flag. The command will print the 'oc image mirror' command
			that can be used to upload the release to another registry.

			You may use --apply-release-image-signature, --release-image-signature-to-dir, or both
			to control the handling of the signature config map. Option
			--apply-release-image-signature will apply the config map directly to a connected
			cluster while --release-image-signature-to-dir specifies an export target directory. If
			--release-image-signature-to-dir is not specified but --to-dir is,
			--release-image-signature-to-dir defaults to a 'config' subdirectory of --to-dir.
			The --overwrite option only applies when --apply-release-image-signature is specified
			and indicates to update an exisiting config map if one is found. A config map written to a
			directory will always replace onethat already exists.
		`),
		Example: templates.Examples(`
			# Perform a dry run showing what would be mirrored, including the mirror objects
			oc adm release mirror 4.11.0 --to myregistry.local/openshift/release \
				--release-image-signature-to-dir /tmp/releases --dry-run

			# Mirror a release into the current directory
			oc adm release mirror 4.11.0 --to file://openshift/release \
				--release-image-signature-to-dir /tmp/releases

			# Mirror a release to another directory in the default location
			oc adm release mirror 4.11.0 --to-dir /tmp/releases

			# Upload a release from the current directory to another server
			oc adm release mirror --from file://openshift/release --to myregistry.com/openshift/release \
				--release-image-signature-to-dir /tmp/releases

			# Mirror the 4.11.0 release to repository registry.example.com and apply signatures to connected cluster
			oc adm release mirror --from=quay.io/openshift-release-dev/ocp-release:4.11.0-x86_64 \
				--to=registry.example.com/your/repository --apply-release-image-signature
		`),
		Run: func(cmd *cobra.Command, args []string) {
			kcmdutil.CheckErr(o.Complete(cmd, f, args))
			kcmdutil.CheckErr(o.Validate())
			kcmdutil.CheckErr(o.Run(cmd.Context()))
		},
	}
	flags := cmd.Flags()
	o.SecurityOptions.Bind(flags)
	o.ParallelOptions.Bind(flags)

	flags.StringVar(&o.From, "from", o.From, "Image containing the release payload.")
	flags.StringVar(&o.To, "to", o.To, "An image repository to push to.")
	flags.StringVar(&o.ToImageStream, "to-image-stream", o.ToImageStream, "An image stream to tag images into.")
	flags.StringVar(&o.FromDir, "from-dir", o.FromDir, "A directory to import images from.")
	flags.StringVar(&o.ToDir, "to-dir", o.ToDir, "A directory to export images to.")
	flags.BoolVar(&o.ToMirror, "to-mirror", o.ToMirror, "Output the mirror mappings instead of mirroring.")
	flags.BoolVar(&o.DryRun, "dry-run", o.DryRun, "Display information about the mirror without actually executing it.")
	flags.BoolVar(&o.KeepManifestList, "keep-manifest-list", o.KeepManifestList, "If an image is part of a manifest list, always mirror the list even if only one image is found.")
	flags.BoolVar(&o.ApplyReleaseImageSignature, "apply-release-image-signature", o.ApplyReleaseImageSignature, "Apply release image signature to connected cluster.")
	flags.StringVar(&o.ReleaseImageSignatureToDir, "release-image-signature-to-dir", o.ReleaseImageSignatureToDir, "A directory to export release image signature to.")

	flags.StringVar(&o.PrintImageSourceInstructions, "print-mirror-instructions", o.PrintImageSourceInstructions, "Print instructions of ImageContentSourcePolicy or ImageDigestMirrorSet for using images from mirror registries. The valid values are 'icsp', 'idms' and 'none'. Default value is icsp.")
	flags.BoolVar(&o.SkipRelease, "skip-release-image", o.SkipRelease, "Do not push the release image.")
	flags.StringVar(&o.ToRelease, "to-release-image", o.ToRelease, "Specify an alternate locations for the release image instead as tag 'release' in --to.")
	flags.BoolVar(&o.Overwrite, "overwrite", o.Overwrite, "Used with --apply-release-image-signature to update an existing signature configmap.")
	return cmd
}

type MirrorOptions struct {
	genericiooptions.IOStreams

	SecurityOptions imagemanifest.SecurityOptions
	ParallelOptions imagemanifest.ParallelOptions

	From    string
	FromDir string

	To            string
	ToImageStream string

	// modifies the targets
	ToRelease   string
	SkipRelease bool

	ToMirror bool
	ToDir    string

	KeepManifestList bool

	ApplyReleaseImageSignature bool
	ReleaseImageSignatureToDir string
	Overwrite                  bool

	DryRun                       bool
	PrintImageSourceInstructions string

	ImageClientFn  func() (imageclient.Interface, string, error)
	CoreV1ClientFn func() (corev1client.ConfigMapInterface, error)

	ImageStream *imagev1.ImageStream
	TargetFn    func(component string) imagereference.DockerImageReference
}

func (o *MirrorOptions) Complete(cmd *cobra.Command, f kcmdutil.Factory, args []string) error {
	switch {
	case len(args) == 0 && len(o.From) == 0:
		return fmt.Errorf("must specify a release image with --from")
	case len(args) == 1 && len(o.From) == 0:
		o.From = args[0]
	case len(args) == 1 && len(o.From) > 0:
		return fmt.Errorf("you may not specify an argument and --from")
	case len(args) > 1:
		return fmt.Errorf("only one argument is accepted")
	}

	args, err := findArgumentsFromCluster(f, []string{o.From})
	if err != nil {
		return err
	}
	if len(args) != 1 {
		return fmt.Errorf("only one release image may be mirrored")
	}
	o.From = args[0]

	o.ImageClientFn = func() (imageclient.Interface, string, error) {
		cfg, err := f.ToRESTConfig()
		if err != nil {
			return nil, "", err
		}
		client, err := imageclient.NewForConfig(cfg)
		if err != nil {
			return nil, "", err
		}
		ns, _, err := f.ToRawKubeConfigLoader().Namespace()
		if err != nil {
			return nil, "", err
		}
		return client, ns, nil
	}
	o.CoreV1ClientFn = func() (corev1client.ConfigMapInterface, error) {
		cfg, err := f.ToRESTConfig()
		if err != nil {
			return nil, err
		}
		coreClient, err := corev1client.NewForConfig(cfg)
		if err != nil {
			return nil, err
		}
		client := coreClient.ConfigMaps(configmap.NamespaceLabelConfigMap)
		return client, nil
	}
	if o.PrintImageSourceInstructions == "" {
		o.PrintImageSourceInstructions = instructionTypeICSP
	}
	instructionType := strings.ToLower(o.PrintImageSourceInstructions)
	switch instructionType {
	case instructionTypeICSP:
		fmt.Fprintf(o.ErrOut, "Flag --print-mirror-instructions's value 'icsp' has been deprecated. Use 'idms' instead to allow the printing of instructions for ImageDigestSources and ImageDigestMirrorSet.\n")
	case instructionTypeIDMS, instructionTypeNone:
	default:
		return fmt.Errorf("--print-mirror-instructions must be one of icsp, idms, none")
	}
	o.PrintImageSourceInstructions = instructionType

	return nil
}

func (o *MirrorOptions) Validate() error {
	if len(o.From) == 0 && o.ImageStream == nil {
		return fmt.Errorf("must specify a release image with --from")
	}

	outputs := 0
	if len(o.To) > 0 {
		outputs++
	}
	if len(o.ToImageStream) > 0 {
		outputs++
	}
	if len(o.ToDir) > 0 {
		if outputs == 0 {
			outputs++
		}
	}
	if o.ToMirror {
		if outputs == 0 {
			outputs++
		}
	}
	if outputs != 1 {
		return fmt.Errorf("must specify an image repository or image stream to mirror the release to")
	}

	if o.SkipRelease && len(o.ToRelease) > 0 {
		return fmt.Errorf("--skip-release-image and --to-release-image may not both be specified")
	}

	if len(o.ReleaseImageSignatureToDir) == 0 && len(o.ToDir) > 0 {
		o.ReleaseImageSignatureToDir = filepath.Join(o.ToDir, configFilesBaseDir)
	}

	if o.Overwrite && !o.ApplyReleaseImageSignature {
		return fmt.Errorf("--overwite is only valid when --apply-release-image-signature is specified")
	}
	return nil
}

const replaceComponentMarker = "X-X-X-X-X-X-X"
const replaceVersionMarker = "V-V-V-V-V-V-V"

// verifyClientBuilder is a wrapper around the operator's HTTPClient method.
// It is used by the image verifier to get an up-to-date http client.
type verifyClientBuilder struct {
	builder func() (*http.Client, error)
}

func (vcb *verifyClientBuilder) HTTPClient() (*http.Client, error) {
	return vcb.builder()
}

func createSignatureFileName(digest string) (string, error) {
	parts := strings.SplitN(digest, ":", 3)
	if len(parts) != 2 || len(parts[0]) == 0 || len(parts[1]) == 0 {
		return "", fmt.Errorf("the provided digest, %s, must be of the form ALGO:HASH", digest)
	}
	algo, hash := parts[0], parts[1]

	if len(hash) > maxDigestHashLen {
		hash = hash[:maxDigestHashLen]
	}
	return fmt.Sprintf(signatureFileNameFmt, algo, hash), nil
}

// handleSignatures implements the image release signature configmap specific logic.
// Signature configmaps may be written to a directory or applied to a cluster.
func (o *MirrorOptions) handleSignatures(context context.Context, signaturesByDigest map[string][][]byte) error {
	var client corev1client.ConfigMapInterface
	if !o.DryRun && o.ApplyReleaseImageSignature {
		var err error
		client, err = o.CoreV1ClientFn()
		if err != nil {
			return fmt.Errorf("creating a Kubernetes API client: %v", err)
		}
	}
	for digest, signatures := range signaturesByDigest {
		cmData, err := verify.GetSignaturesAsConfigmap(digest, signatures)
		if err != nil {
			return fmt.Errorf("converting signatures to a configmap: %v", err)
		}
		if o.ApplyReleaseImageSignature {
			if o.DryRun {
				fmt.Fprintf(o.Out, "info: Create or configure configmap %s\n", cmData.Name)
			} else {
				var create bool = true
				if o.Overwrite {
					// An error is returned if the configmap does not exist in which case we will
					// attempt to create the manifest.
					if _, err := client.Get(context, cmData.Name, metav1.GetOptions{}); err == nil {
						create = false
						if _, err := client.Update(context, cmData, metav1.UpdateOptions{}); err != nil {
							return fmt.Errorf("updating configmap %s: %v", cmData.Name, err)
						} else {
							fmt.Fprintf(o.Out, "configmap/%s configured\n", cmData.Name)
						}
					}
				}
				if create {
					if _, err := client.Create(context, cmData, metav1.CreateOptions{}); err != nil {
						return fmt.Errorf("creating configmap %s: %v", cmData.Name, err)
					} else {
						fmt.Fprintf(o.Out, "configmap/%s created\n", cmData.Name)
					}
				}
			}
		}
		if len(o.ReleaseImageSignatureToDir) > 0 {
			fileName, err := createSignatureFileName(digest)
			if err != nil {
				return fmt.Errorf("creating filename: %v", err)
			}
			fullName := filepath.Join(o.ReleaseImageSignatureToDir, fileName)
			if o.DryRun {
				fmt.Fprintf(o.Out, "info: Write configmap signature file %s\n", fullName)
			} else {
				cmDataBytes, err := util.ConfigMapAsBytes(cmData)
				if err != nil {
					return fmt.Errorf("marshaling configmap YAML: %v", err)
				}
				if err := os.MkdirAll(filepath.Dir(fullName), 0750); err != nil {
					return err
				}
				if err := os.WriteFile(fullName, cmDataBytes, 0640); err != nil {
					return err
				}
				fmt.Fprintf(o.Out, "Configmap signature file %s created\n", fullName)
			}
		}
	}
	return nil
}

func (o *MirrorOptions) Run(ctx context.Context) error {
	var recreateRequired bool
	var hasPrefix bool
	var targetFn func(name string) imagesource.TypedImageReference
	var dst string
	if len(o.ToImageStream) > 0 {
		dst = imagereference.DockerImageReference{
			Registry:  "example.com",
			Namespace: "somenamespace",
			Name:      "mirror",
		}.Exact()
	} else {
		dst = o.To
	}
	if len(dst) == 0 {
		if len(o.ToDir) > 0 {
			dst = "file://openshift/release"
		} else {
			dst = "openshift/release"
		}
	}

	var toDisk bool
	var version string
	var archExt string
	if strings.Contains(dst, "${component}") {
		format := strings.Replace(dst, "${component}", replaceComponentMarker, -1)
		format = strings.Replace(format, "${version}", replaceVersionMarker, -1)
		dstRef, err := imagesource.ParseReference(format)
		if err != nil {
			return fmt.Errorf("--to must be a valid image reference: %v", err)
		}
		toDisk = dstRef.Type == imagesource.DestinationFile
		targetFn = func(name string) imagesource.TypedImageReference {
			if len(name) == 0 {
				name = "release"
			}
			value := strings.Replace(dst, "${component}", name, -1)
			value = strings.Replace(value, "${version}", version, -1)
			value = value + archExt
			ref, err := imagesource.ParseReference(value)
			if err != nil {
				klog.Fatalf("requested component %q could not be injected into %s: %v", name, dst, err)
			}
			return ref
		}
		replaceCount := strings.Count(dst, "${component}")
		recreateRequired = replaceCount > 1 || (replaceCount == 1 && !strings.Contains(dstRef.Ref.Tag, replaceComponentMarker))

	} else {
		ref, err := imagesource.ParseReference(dst)
		if err != nil {
			return fmt.Errorf("--to must be a valid image repository: %v", err)
		}
		toDisk = ref.Type == imagesource.DestinationFile
		if len(ref.Ref.ID) > 0 || len(ref.Ref.Tag) > 0 {
			return fmt.Errorf("--to must be to an image repository and may not contain a tag or digest")
		}
		targetFn = func(name string) imagesource.TypedImageReference {
			copied := ref
			if len(name) > 0 {
				copied.Ref.Tag = fmt.Sprintf("%s%s-%s", version, archExt, name)
			} else {
				copied.Ref.Tag = fmt.Sprintf("%s%s", version, archExt)
			}
			return copied
		}
		hasPrefix = true
	}

	o.TargetFn = func(name string) imagereference.DockerImageReference {
		ref := targetFn(name)
		return ref.Ref
	}

	if recreateRequired {
		return fmt.Errorf("when mirroring to multiple repositories, use the new release command with --from-release and --mirror")
	}

	var releaseDigest string
	var manifests []manifest.Manifest
	is := o.ImageStream
	if is == nil {
		// load image references
		extractOpts := NewExtractOptions(genericiooptions.IOStreams{ErrOut: o.ErrOut}, true)
		extractOpts.Directory = ""
		extractOpts.ParallelOptions = o.ParallelOptions
		extractOpts.SecurityOptions = o.SecurityOptions
		if o.KeepManifestList {
			// we'll always use manifests from the linux/amd64 image, since the manifests
			// won't differ between architectures, at least for now
			re, err := regexp.Compile("linux/amd64")
			if err != nil {
				return err
			}
			extractOpts.FilterOptions.OSFilter = re
		}
		extractOpts.ImageMetadataCallback = func(m *extract.Mapping, dgst, contentDigest digest.Digest, config *dockerv1client.DockerImageConfig, manifestListDigest digest.Digest) {
			releaseDigest = contentDigest.String()
			if config != nil {
				// Use 'multi' instead of config.Architecture if keeping the ManifestList.
				if o.KeepManifestList && manifestListDigest != "" {
					archExt = "-multi"
				} else if val, ok := archMap[config.Architecture]; ok {
					archExt = "-" + val
				} else {
					archExt = "-" + config.Architecture
				}
			} else {
				fmt.Fprintf(o.ErrOut, "warning: Unable to retrieve image release architecture\n")
			}
		}
		extractOpts.FileDir = o.FromDir
		extractOpts.From = o.From
		if err := extractOpts.Run(ctx); err != nil {
			return fmt.Errorf("unable to retrieve release image info: %v", err)
		}

		is = extractOpts.ImageReferences
		if is == nil {
			return errors.New("unable to load image-references from release payload")
		}
		manifests = extractOpts.Manifests
	}
	version = is.Name

	// sourceFn is given a chance to rewrite source mappings
	sourceFn := func(ref imagesource.TypedImageReference) imagesource.TypedImageReference {
		return ref
	}

	httpClientConstructor := sigstore.NewCachedHTTPClientConstructor(o.HTTPClient, nil)

	// Attempt to load a verifier as defined by the release being mirrored
	imageVerifier, err := verify.NewFromManifests(manifests, httpClientConstructor.HTTPClient)
	if err != nil {
		return fmt.Errorf("Unable to load configmap verifier: %v", err)
	}
	if imageVerifier != nil {
		klog.V(4).Infof("Verifying release authenticity: %v", imageVerifier)
	} else {
		fmt.Fprintf(o.ErrOut, "warning: No release authenticity verification is configured, all releases are considered unverified\n")
		imageVerifier = verify.Reject
	}
	// verify the provided payload
	ctx, cancelFn := context.WithCancel(context.Background())
	defer cancelFn()
	if err := imageVerifier.Verify(ctx, releaseDigest); err != nil {
		fmt.Fprintf(o.ErrOut, "warning: An image was retrieved that failed verification: %v\n", err)
	}
	var mappings []mirror.Mapping
	if len(o.From) > 0 {
		src := o.From
		srcRef, err := imagesource.ParseReference(src)
		if err != nil {
			return fmt.Errorf("invalid --from: %v", err)
		}

		// if the source ref is a file type, provide a function that checks the local file store for a given manifest
		// before continuing, to allow mirroring an entire release to disk in a single file://REPO.
		if srcRef.Type == imagesource.DestinationFile {
			if repo, err := (&imagesource.Options{FileDir: o.FromDir}).Repository(context.TODO(), srcRef); err == nil {
				sourceFn = func(ref imagesource.TypedImageReference) imagesource.TypedImageReference {
					if ref.Type == imagesource.DestinationFile || len(ref.Ref.ID) == 0 {
						return ref
					}
					manifests, err := repo.Manifests(context.TODO())
					if err != nil {
						klog.V(2).Infof("Unable to get local manifest service: %v", err)
						return ref
					}
					ok, err := manifests.Exists(context.TODO(), digest.Digest(ref.Ref.ID))
					if err != nil {
						klog.V(2).Infof("Unable to get check for local manifest: %v", err)
						return ref
					}
					if !ok {
						return ref
					}
					updated := srcRef
					updated.Ref.Tag = ""
					updated.Ref.ID = ref.Ref.ID
					klog.V(2).Infof("Rewrote %s to %s", ref, updated)
					return updated
				}
			} else {
				klog.V(2).Infof("Unable to build local file lookup: %v", err)
			}
		}

		if len(o.ToRelease) > 0 {
			dstRef, err := imagesource.ParseReference(o.ToRelease)
			if err != nil {
				return fmt.Errorf("invalid --to-release-image: %v", err)
			}
			mappings = append(mappings, mirror.Mapping{
				Source:      srcRef,
				Destination: dstRef,
				Name:        o.ToRelease,
			})
		} else if !o.SkipRelease {
			dstRef := targetFn("")
			mappings = append(mappings, mirror.Mapping{
				Source:      srcRef,
				Destination: dstRef,
				Name:        "release",
			})
		}
	}

	repositories := make(map[string]struct{})

	// build the mapping list for mirroring and rewrite if necessary
	for i := range is.Spec.Tags {
		tag := &is.Spec.Tags[i]
		if tag.From == nil || tag.From.Kind != "DockerImage" {
			continue
		}
		from, err := imagereference.Parse(tag.From.Name)
		if err != nil {
			return fmt.Errorf("release tag %q is not valid: %v", tag.Name, err)
		}
		if len(from.Tag) > 0 || len(from.ID) == 0 {
			return fmt.Errorf("image-references should only contain pointers to images by digest: %s", tag.From.Name)
		}

		// Allow mirror refs to be sourced locally
		srcMirrorRef := imagesource.TypedImageReference{Ref: from, Type: imagesource.DestinationRegistry}
		srcMirrorRef = sourceFn(srcMirrorRef)

		// Create a unique map of repos as keys
		currentRepo := from.AsRepository().String()
		repositories[currentRepo] = struct{}{}

		dstMirrorRef := targetFn(tag.Name)
		mappings = append(mappings, mirror.Mapping{
			Source:      srcMirrorRef,
			Destination: dstMirrorRef,
			Name:        tag.Name,
		})
		klog.V(2).Infof("Mapping %#v", mappings[len(mappings)-1])

		dstRef := targetFn(tag.Name)
		dstRef.Ref.Tag = ""
		dstRef.Ref.ID = from.ID
		tag.From.Name = dstRef.Ref.Exact()
	}

	if len(mappings) == 0 {
		fmt.Fprintf(o.ErrOut, "warning: Release image contains no image references - is this a valid release?\n")
	}

	if o.ToMirror {
		for _, mapping := range mappings {
			fmt.Fprintf(o.Out, "%s %s\n", mapping.Source.String(), mapping.Destination.String())
		}
		return nil
	}

	if len(o.ToImageStream) > 0 {
		remaining := make(map[string]mirror.Mapping)
		for _, mapping := range mappings {
			remaining[mapping.Name] = mapping
		}
		client, ns, err := o.ImageClientFn()
		if err != nil {
			return err
		}
		hasErrors := make(map[string]error)
		maxPerIteration := 12

		importMode := imagev1.ImportModeLegacy
		if o.KeepManifestList {
			importMode = imagev1.ImportModePreserveOriginal
		}
		for retries := 4; (len(remaining) > 0 || len(hasErrors) > 0) && retries > 0; {
			if len(remaining) == 0 {
				for _, mapping := range mappings {
					if _, ok := hasErrors[mapping.Name]; ok {
						remaining[mapping.Name] = mapping
						delete(hasErrors, mapping.Name)
					}
				}
				retries--
			}
			err := retry.RetryOnConflict(retry.DefaultRetry, func() error {
				isi := &imagev1.ImageStreamImport{
					ObjectMeta: metav1.ObjectMeta{
						Name: o.ToImageStream,
					},
					Spec: imagev1.ImageStreamImportSpec{
						Import: !o.DryRun,
					},
				}
				for _, mapping := range remaining {
					if mapping.Source.Type != imagesource.DestinationRegistry {
						return fmt.Errorf("source mapping %s must point to a registry", mapping.Source)
					}
					isi.Spec.Images = append(isi.Spec.Images, imagev1.ImageImportSpec{
						From: corev1.ObjectReference{
							Kind: "DockerImage",
							Name: mapping.Source.Ref.Exact(),
						},
						To: &corev1.LocalObjectReference{
							Name: mapping.Name,
						},
						ImportPolicy: imagev1.TagImportPolicy{
							ImportMode: importMode,
						},
					})
					if len(isi.Spec.Images) > maxPerIteration {
						break
					}
				}

				// use RESTClient directly here to be able to extend request timeout
				result := &imagev1.ImageStreamImport{}
				if err := client.ImageV1().RESTClient().Post().
					Namespace(ns).
					Resource(imagev1.Resource("imagestreamimports").Resource).
					Body(isi).
					// this instructs the api server to allow our request to take up to an hour - chosen as a high boundary
					Timeout(3 * time.Minute).
					Do(context.TODO()).
					Into(result); err != nil {
					return err
				}

				for i, image := range result.Status.Images {
					name := result.Spec.Images[i].To.Name
					klog.V(4).Infof("Import result for %s: %#v", name, image.Status)
					if image.Status.Status == metav1.StatusSuccess {
						delete(remaining, name)
						delete(hasErrors, name)
					} else {
						delete(remaining, name)
						err := apierrors.FromObject(&image.Status)
						hasErrors[name] = err
						klog.V(2).Infof("Failed to import %s as tag %s: %v", remaining[name].Source, name, err)
					}
				}
				return nil
			})
			if err != nil {
				return err
			}
		}

		if len(hasErrors) > 0 {
			var messages []string
			for k, v := range hasErrors {
				messages = append(messages, fmt.Sprintf("%s: %v", k, v))
			}
			sort.Strings(messages)
			if len(messages) == 1 {
				return fmt.Errorf("unable to import a release image: %s", messages[0])
			}
			return fmt.Errorf("unable to import some release images:\n* %s", strings.Join(messages, "\n* "))
		}

		fmt.Fprintf(os.Stderr, "Mirrored %d images to %s/%s\n", len(mappings), ns, o.ToImageStream)
		return nil
	}

	fmt.Fprintf(os.Stderr, "info: Mirroring %d images to %s ...\n", len(mappings), dst)
	var lock sync.Mutex
	opts := mirror.NewMirrorImageOptions(genericiooptions.IOStreams{Out: o.Out, ErrOut: o.ErrOut})
	opts.SecurityOptions = o.SecurityOptions
	opts.ParallelOptions = o.ParallelOptions
	opts.Mappings = mappings
	opts.FromFileDir = o.FromDir
	opts.FileDir = o.ToDir
	opts.DryRun = o.DryRun
	opts.KeepManifestList = o.KeepManifestList
	opts.ManifestUpdateCallback = func(registry string, manifests map[digest.Digest]digest.Digest) error {
		lock.Lock()
		defer lock.Unlock()

		// when uploading to a schema1 registry, manifest ids change and we must remap them
		for i := range is.Spec.Tags {
			tag := &is.Spec.Tags[i]
			if tag.From == nil || tag.From.Kind != "DockerImage" {
				continue
			}
			ref, err := imagereference.Parse(tag.From.Name)
			if err != nil {
				return fmt.Errorf("unable to parse image reference %s (%s): %v", tag.Name, tag.From.Name, err)
			}
			if ref.Registry != registry {
				continue
			}
			if changed, ok := manifests[digest.Digest(ref.ID)]; ok {
				ref.ID = changed.String()
				klog.V(4).Infof("During mirroring, image %s was updated to digest %s", tag.From.Name, changed)
				tag.From.Name = ref.Exact()
			}
		}
		return nil
	}
	if err := opts.Run(); err != nil {
		return err
	}

	to := o.ToRelease
	if len(to) == 0 {
		to = targetFn("").Ref.Exact()
	}

	fmt.Fprintf(o.Out, "\nSuccess\nUpdate image:  %s\n", to)
	var toList []string
	if len(o.To) > 0 {
		toList = append(toList, o.To)
	}
	if len(o.ToRelease) > 0 {
		toList = append(toList, o.ToRelease)
	}
	if len(toList) > 0 {
		for _, t := range toList {
			if hasPrefix {
				fmt.Fprintf(o.Out, "Mirror prefix: %s\n", t)
			} else {
				fmt.Fprintf(o.Out, "Mirrored to: %s\n", t)
			}
		}
	}
	if toDisk {
		if len(o.ToDir) > 0 {
			fmt.Fprintf(o.Out, "\nTo upload local images to a registry, run:\n\n    oc image mirror --from-dir=%s 'file://%s*' REGISTRY/REPOSITORY\n\n", o.ToDir, to)
		} else {
			fmt.Fprintf(o.Out, "\nTo upload local images to a registry, run:\n\n    oc image mirror 'file://%s*' REGISTRY/REPOSITORY\n\n", to)
		}
	} else if len(toList) > 0 {
		if err := printImageMirrorInstructions(o.Out, o.From, toList, o.ReleaseImageSignatureToDir, repositories, o.PrintImageSourceInstructions); err != nil {
			return fmt.Errorf("error creating mirror usage instructions: %v", err)
		}
	}
	if o.ApplyReleaseImageSignature || len(o.ReleaseImageSignatureToDir) > 0 {
		signatures := imageVerifier.Signatures()
		if signatures == nil || len(signatures) == 0 {
			return errors.New("failed to retrieve cached signatures")
		}
		if _, ok := signatures[releaseDigest]; ok {
			err := o.handleSignatures(ctx, signatures)
			if err != nil {
				return fmt.Errorf("handling release image signatures: %v", err)
			}
		} else {
			digests := make([]string, 0, len(signatures))
			for digest := range signatures {
				digests = append(digests, digest)
			}
			sort.Strings(digests)
			return fmt.Errorf("no cached signatures for digest %s, just:\n%s", releaseDigest, strings.Join(digests, "\n"))
		}
	}
	return nil
}

type mirrorSet struct {
	source  string
	mirrors []string
}

// printImageMirrorInstructions provides examples to the user for using the new repository mirror.
// https://github.com/openshift/installer/blob/master/docs/dev/alternative_release_image_sources.md
func printImageMirrorInstructions(out io.Writer, from string, toList []string, signatureToDir string, repositories map[string]struct{}, printImageSourceInstructions string) error {

	if printImageSourceInstructions == instructionTypeNone {
		return nil
	}

	var (
		sources []mirrorSet
		err     error
	)

	for _, to := range toList {
		mirrorRef, err := imagesource.ParseReference(to)
		if err != nil {
			return fmt.Errorf("unable to parse image reference '%s': %v", to, err)
		}
		if mirrorRef.Type != imagesource.DestinationRegistry {
			return nil
		}
		mirrorRepo := mirrorRef.Ref.AsRepository().String()
		if len(from) != 0 {
			sourceRef, err := imagesource.ParseReference(from)
			if err != nil {
				return fmt.Errorf("unable to parse image reference '%s': %v", from, err)
			}
			if sourceRef.Type != imagesource.DestinationRegistry {
				return nil
			}
			sourceRepo := sourceRef.Ref.AsRepository().String()
			repositories[sourceRepo] = struct{}{}
		}

		if len(repositories) == 0 {
			return nil
		}

		for repository := range repositories {
			sources = append(sources, mirrorSet{
				source:  repository,
				mirrors: []string{mirrorRepo},
			})
		}
	}
	// de-duplicate sources
	uniqueSources := dedupeSortSources(sources)
	sources = uniqueSources

	// Create and display install-config.yaml example
	// print instructions for either ICSP or IDMS, should drop printICSPInstructions() when ICSP is no longer supported.
	if printImageSourceInstructions == instructionTypeIDMS {
		err = printIDMSInstructions(out, sources)
	} else if printImageSourceInstructions == instructionTypeICSP {
		err = printICSPInstructions(out, sources)
	}
	if len(signatureToDir) != 0 {
		fmt.Fprintf(out, "\n\nTo apply signature configmaps use 'oc apply' on files found in %s\n\n", signatureToDir)
	}
	return err
}

// Create and display install-config.yaml example using ImageContentSourcePolicy(ICSP)
func printICSPInstructions(out io.Writer, sources []mirrorSet) error {
	type installConfigSubsection struct {
		ImageContentSources []operatorv1alpha1.RepositoryDigestMirrors `json:"imageContentSources"`
	}
	icspSources := []operatorv1alpha1.RepositoryDigestMirrors{}
	for _, s := range sources {
		icspSources = append(icspSources, operatorv1alpha1.RepositoryDigestMirrors{
			Source:  s.source,
			Mirrors: s.mirrors,
		})
	}
	imageContentSources := installConfigSubsection{
		ImageContentSources: icspSources}
	installConfigExample, err := yaml.Marshal(imageContentSources)
	if err != nil {
		return fmt.Errorf("unable to marshal install-config.yaml example yaml: %v", err)
	}
	fmt.Fprintf(out, "\nTo use the new mirrored repository to install, add the following section to the install-config.yaml:\n\n")
	fmt.Fprint(out, string(installConfigExample))

	// Create and display ImageContentSourcePolicy example
	icsp := operatorv1alpha1.ImageContentSourcePolicy{
		TypeMeta: metav1.TypeMeta{
			APIVersion: operatorv1alpha1.GroupVersion.String(),
			Kind:       "ImageContentSourcePolicy"},
		ObjectMeta: metav1.ObjectMeta{
			Name: "example",
		},
		Spec: operatorv1alpha1.ImageContentSourcePolicySpec{
			RepositoryDigestMirrors: icspSources,
		},
	}

	// Create an unstructured object for removing creationTimestamp
	unstructuredObj := unstructured.Unstructured{}
	unstructuredObj.Object, err = runtime.DefaultUnstructuredConverter.ToUnstructured(&icsp)
	if err != nil {
		return fmt.Errorf("ToUnstructured error: %v", err)
	}
	delete(unstructuredObj.Object["metadata"].(map[string]interface{}), "creationTimestamp")

	icspExample, err := yaml.Marshal(unstructuredObj.Object)
	if err != nil {
		return fmt.Errorf("unable to marshal ImageContentSourcePolicy example yaml: %v", err)
	}
	fmt.Fprintf(out, "\n\nTo use the new mirrored repository for upgrades, use the following to create an ImageContentSourcePolicy:\n\n")
	fmt.Fprint(out, string(icspExample))

	return nil
}

// Create and display install-config.yaml example using ImageDigestMirrorSet(IDMS)
func printIDMSInstructions(out io.Writer, sources []mirrorSet) error {
	type installConfigSubsection struct {
		ImageDigestSources []apicfgv1.ImageDigestMirrors `json:"imageDigestSources"`
	}
	idmsSources := convertMirrorSetToImageDigestMirrors(sources)
	imageDigestSources := installConfigSubsection{
		ImageDigestSources: idmsSources}
	installConfigExample, err := yaml.Marshal(imageDigestSources)
	if err != nil {
		return fmt.Errorf("unable to marshal install-config.yaml example yaml: %v", err)
	}
	fmt.Fprintf(out, "\nTo use the new mirrored repository to install, add the following section to the install-config.yaml:\n\n")
	fmt.Fprint(out, string(installConfigExample))

	// Create and display ImageDigestMirrorSet example
	idms := apicfgv1.ImageDigestMirrorSet{
		TypeMeta: metav1.TypeMeta{
			APIVersion: apicfgv1.GroupVersion.String(),
			Kind:       "ImageDigestMirrorSet"},
		ObjectMeta: metav1.ObjectMeta{
			Name: "example",
		},
		Spec: apicfgv1.ImageDigestMirrorSetSpec{
			ImageDigestMirrors: idmsSources,
		},
	}

	// Create an unstructured object for removing creationTimestamp, status
	unstructuredObj := unstructured.Unstructured{}
	unstructuredObj.Object, err = runtime.DefaultUnstructuredConverter.ToUnstructured(&idms)
	if err != nil {
		return fmt.Errorf("ToUnstructured error: %v", err)
	}
	delete(unstructuredObj.Object["metadata"].(map[string]interface{}), "creationTimestamp")
	delete(unstructuredObj.Object, "status")

	idmsExample, err := yaml.Marshal(unstructuredObj.Object)
	if err != nil {
		return fmt.Errorf("unable to marshal ImageDigestMirrorSet example yaml: %v", err)
	}
	fmt.Fprintf(out, "\n\nTo use the new mirrored repository for upgrades, use the following to create an ImageDigestMirrorSet:\n\n")
	fmt.Fprint(out, string(idmsExample))

	return nil
}

// HTTPClient provides a method for generating an HTTP client
// with the proxy and trust settings, if set in the cluster.
func (o *MirrorOptions) HTTPClient() (*http.Client, error) {
	transport, err := transport.HTTPWrappersForConfig(
		&transport.Config{
			UserAgent: rest.DefaultKubernetesUserAgent() + "(release-mirror)",
		},
		http.DefaultTransport,
	)
	if err != nil {
		return nil, err
	}
	return &http.Client{
		Transport: transport,
	}, nil
}

func convertMirrorSetToImageDigestMirrors(sources []mirrorSet) []apicfgv1.ImageDigestMirrors {
	idmsSources := []apicfgv1.ImageDigestMirrors{}
	for _, s := range sources {
		mirrors := []apicfgv1.ImageMirror{}
		for _, m := range s.mirrors {
			mirrors = append(mirrors, apicfgv1.ImageMirror(m))
		}
		idmsSources = append(idmsSources, apicfgv1.ImageDigestMirrors{
			Source:  s.source,
			Mirrors: mirrors,
		})
	}
	return idmsSources
}

func dedupeSortSources(sources []mirrorSet) []mirrorSet {
	unique := make(map[string][]string)
	var uniqueSources []mirrorSet
	for _, s := range sources {
		if mirrors, ok := unique[s.source]; ok {
			for _, m := range s.mirrors {
				if !contains(mirrors, m) {
					mirrors = append(mirrors, m)
				}
			}
			unique[s.source] = mirrors
		} else {
			unique[s.source] = s.mirrors
		}
	}
	for s, m := range unique {
		uniqueSources = append(uniqueSources, mirrorSet{
			source:  s,
			mirrors: m,
		})
	}
	sort.Slice(uniqueSources, func(i, j int) bool {
		return uniqueSources[i].source < sources[j].source
	})
	return uniqueSources
}

func contains(s []string, str string) bool {
	for _, v := range s {
		if v == str {
			return true
		}
	}
	return false
}
