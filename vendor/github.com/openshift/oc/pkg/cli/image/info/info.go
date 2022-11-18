package info

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"sort"
	"strings"
	"sync"
	"text/tabwriter"
	"time"

	"github.com/docker/distribution"
	"github.com/docker/distribution/manifest/manifestlist"
	units "github.com/docker/go-units"
	digest "github.com/opencontainers/go-digest"
	"github.com/spf13/cobra"

	"k8s.io/apimachinery/pkg/util/duration"
	"k8s.io/apimachinery/pkg/util/sets"
	"k8s.io/cli-runtime/pkg/genericclioptions"
	"k8s.io/klog/v2"
	kcmdutil "k8s.io/kubectl/pkg/cmd/util"
	"k8s.io/kubectl/pkg/util/templates"

	"github.com/openshift/library-go/pkg/image/dockerv1client"
	"github.com/openshift/library-go/pkg/image/registryclient"
	"github.com/openshift/oc/pkg/cli/image/imagesource"
	imagemanifest "github.com/openshift/oc/pkg/cli/image/manifest"
	"github.com/openshift/oc/pkg/cli/image/strategy"
	"github.com/openshift/oc/pkg/cli/image/workqueue"
)

func NewInfoOptions(streams genericclioptions.IOStreams) *InfoOptions {
	return &InfoOptions{
		IOStreams: streams,
	}
}

func NewInfo(f kcmdutil.Factory, streams genericclioptions.IOStreams) *cobra.Command {
	o := NewInfoOptions(streams)
	cmd := &cobra.Command{
		Use:   "info IMAGE [...]",
		Short: "Display information about an image",
		Long: templates.LongDesc(`
			Show information about an image in a remote image registry.

			This command will retrieve metadata about container images in a remote image
			registry. You may specify images by tag or digest and specify multiple at a
			time.

			Images in manifest list format will be shown for your current operating system.
			To see the image for a particular OS use the --filter-by-os=OS/ARCH flag.
		`),
		Example: templates.Examples(`
			# Show information about an image
			oc image info quay.io/openshift/cli:latest

			# Show information about images matching a wildcard
			oc image info quay.io/openshift/cli:4.*

			# Show information about a file mirrored to disk under DIR
			oc image info --dir=DIR file://library/busybox:latest

			# Select which image from a multi-OS image to show
			oc image info library/busybox:latest --filter-by-os=linux/arm64

		`),
		Run: func(cmd *cobra.Command, args []string) {
			kcmdutil.CheckErr(o.Complete(f, cmd, args))
			kcmdutil.CheckErr(o.Validate(cmd))
			kcmdutil.CheckErr(o.Run())
		},
	}
	flags := cmd.Flags()
	o.FilterOptions.Bind(flags)
	o.SecurityOptions.Bind(flags)
	flags.StringVarP(&o.Output, "output", "o", o.Output, "Print the image in an alternative format: json")
	flags.StringVar(&o.FileDir, "dir", o.FileDir, "The directory on disk that file:// images will be read from.")
	flags.StringVar(&o.ICSPFile, "icsp-file", o.ICSPFile, "Path to an ImageContentSourcePolicy file.  If set, data from this file will be used to find alternative locations for images.")
	flags.BoolVar(&o.ShowMultiArch, "show-multiarch", o.ShowMultiArch, "Show information even if the image is multiarch image. If not set, error is thrown for multiarch images.")

	return cmd
}

type InfoOptions struct {
	genericclioptions.IOStreams

	SecurityOptions imagemanifest.SecurityOptions
	FilterOptions   imagemanifest.FilterOptions

	Images        []string
	FileDir       string
	Output        string
	ICSPFile      string
	ShowMultiArch bool
}

func (o *InfoOptions) Complete(f kcmdutil.Factory, cmd *cobra.Command, args []string) error {
	if len(args) < 1 {
		return fmt.Errorf("info expects at least one argument, an image pull spec")
	}
	o.Images = args

	return nil
}

func (o *InfoOptions) Validate(cmd *cobra.Command) error {
	if len(o.Images) == 0 {
		return fmt.Errorf("must specify one or more images as arguments")
	}
	return o.FilterOptions.Validate()
}

func (o *InfoOptions) Run() error {
	// cache the context
	registryContext, err := o.SecurityOptions.Context()
	if err != nil {
		return err
	}
	if len(o.ICSPFile) > 0 {
		registryContext = registryContext.WithAlternateBlobSourceStrategy(strategy.NewICSPOnErrorStrategy(o.ICSPFile))
	}
	opts := &imagesource.Options{
		FileDir:         o.FileDir,
		Insecure:        o.SecurityOptions.Insecure,
		RegistryContext: registryContext,
	}

	hadError := false
	icspWarned := false
	for _, location := range o.Images {
		sources, err := imagesource.ParseSourceReference(location, opts.ExpandWildcard)
		if err != nil {
			return err
		}
		for _, src := range sources {
			if len(src.Ref.Tag) == 0 && len(src.Ref.ID) == 0 {
				return fmt.Errorf("--from must point to an image ID or image tag")
			}
			if !icspWarned && len(o.ICSPFile) > 0 && len(src.Ref.Tag) > 0 {
				fmt.Fprintf(o.ErrOut, "warning: --icsp-file only applies to images referenced by digest and will be ignored for tags\n")
				icspWarned = true
			}

			var images []*Image
			retriever := &ImageRetriever{
				FileDir:         o.FileDir,
				SecurityOptions: o.SecurityOptions,
				ManifestListCallback: func(from string, list *manifestlist.DeserializedManifestList, all map[digest.Digest]distribution.Manifest) (map[digest.Digest]distribution.Manifest, error) {
					filtered := make(map[digest.Digest]distribution.Manifest)
					for _, manifest := range list.Manifests {
						if !o.FilterOptions.Include(&manifest, len(list.Manifests) > 1) {
							klog.V(5).Infof("Skipping image for %#v from %s", manifest.Platform, from)
							continue
						}
						filtered[manifest.Digest] = all[manifest.Digest]
					}
					if len(filtered) == 1 {
						return filtered, nil
					}

					if o.ShowMultiArch {
						return filtered, nil
					}

					buf := &bytes.Buffer{}
					w := tabwriter.NewWriter(buf, 0, 0, 1, ' ', 0)
					fmt.Fprintf(w, "  OS\tDIGEST\n")
					for _, manifest := range list.Manifests {
						fmt.Fprintf(w, "  %s\t%s\n", imagemanifest.PlatformSpecString(manifest.Platform), manifest.Digest)
					}
					w.Flush()
					return nil, fmt.Errorf("the image is a manifest list and contains multiple images - use --filter-by-os to select from:\n\n%s\n", buf.String())
				},

				ImageMetadataCallback: func(from string, i *Image, err error) error {
					if err != nil {
						return err
					}
					images = append(images, i)
					return nil
				},
			}
			if _, err := retriever.Image(context.TODO(), src); err != nil {
				return err
			}

			switch o.Output {
			case "":
			case "json":
				var data []byte
				if len(images) == 1 {
					data, err = json.MarshalIndent(images[0], "", "  ")
				} else {
					data, err = json.MarshalIndent(images, "", "  ")
				}
				if err != nil {
					return err
				}
				fmt.Fprintf(o.Out, "%s", string(data))
				continue
			default:
				return fmt.Errorf("unrecognized --output, only 'json' is supported")
			}

			for _, img := range images {
				if err := describeImage(o.Out, img); err != nil {
					hadError = true
					if err != kcmdutil.ErrExit {
						fmt.Fprintf(o.ErrOut, "error: %v", err)
					}
				}
			}
		}
	}
	if hadError {
		return kcmdutil.ErrExit
	}
	return nil
}

type Image struct {
	Name          string                            `json:"name"`
	Ref           imagesource.TypedImageReference   `json:"-"`
	Digest        digest.Digest                     `json:"digest"`
	ContentDigest digest.Digest                     `json:"contentDigest"`
	ListDigest    digest.Digest                     `json:"listDigest"`
	MediaType     string                            `json:"mediaType"`
	Layers        []distribution.Descriptor         `json:"layers"`
	Config        *dockerv1client.DockerImageConfig `json:"config"`

	Manifest distribution.Manifest `json:"-"`
}

func describeImage(out io.Writer, image *Image) error {
	var err error

	w := tabwriter.NewWriter(out, 0, 4, 1, ' ', 0)
	defer w.Flush()
	fmt.Fprintf(w, "Name:\t%s\n", image.Name)
	if len(image.Ref.Ref.ID) == 0 || image.Ref.Ref.ID != image.Digest.String() {
		fmt.Fprintf(w, "Digest:\t%s\n", image.Digest)
	}
	if len(image.ListDigest) > 0 {
		fmt.Fprintf(w, "Manifest List:\t%s\n", image.ListDigest)
	}
	if image.ContentDigest != image.Digest {
		fmt.Fprintf(w, "Content Digest:\t%s\n\tERROR: the image contents do not match the requested digest, this image has been tampered with\n", image.ContentDigest)
		err = kcmdutil.ErrExit
	}

	fmt.Fprintf(w, "Media Type:\t%s\n", image.MediaType)
	if image.Config.Created.IsZero() {
		fmt.Fprintf(w, "Created:\t%s\n", "<unknown>")
	} else {
		fmt.Fprintf(w, "Created:\t%s ago\n", duration.ShortHumanDuration(time.Now().Sub(image.Config.Created)))
	}
	switch l := len(image.Layers); l {
	case 0:
		// legacy case, server does not know individual layers
		fmt.Fprintf(w, "Layer Size:\t%s\n", units.HumanSize(float64(image.Config.Size)))
	default:
		imageSize := fmt.Sprintf("%s in %d layers", units.HumanSize(float64(image.Config.Size)), len(image.Layers))
		if image.Config.Size == 0 {
			imageSize = fmt.Sprintf("%d layers (size unavailable)", len(image.Layers))
		}

		fmt.Fprintf(w, "Image Size:\t%s\n", imageSize)
		for i, layer := range image.Layers {
			layerSize := units.HumanSize(float64(layer.Size))
			if layer.Size == 0 {
				layerSize = "--"
			}

			if i == 0 {
				fmt.Fprintf(w, "%s\t%s\t%s\n", "Layers:", layerSize, layer.Digest)
			} else {
				fmt.Fprintf(w, "%s\t%s\t%s\n", "", layerSize, layer.Digest)
			}
		}
	}
	fmt.Fprintf(w, "OS:\t%s\n", image.Config.OS)
	fmt.Fprintf(w, "Arch:\t%s\n", image.Config.Architecture)
	if len(image.Config.Author) > 0 {
		fmt.Fprintf(w, "Author:\t%s\n", image.Config.Author)
	}

	config := image.Config.Config
	if config != nil {
		hasCommand := false
		if len(config.Entrypoint) > 0 {
			hasCommand = true
			fmt.Fprintf(w, "Entrypoint:\t%s\n", strings.Join(config.Entrypoint, " "))
		}
		if len(config.Cmd) > 0 {
			hasCommand = true
			fmt.Fprintf(w, "Command:\t%s\n", strings.Join(config.Cmd, " "))
		}
		if !hasCommand {
			fmt.Fprintf(w, "Command:\t%s\n", "<none>")
		}
		if len(config.WorkingDir) > 0 {
			fmt.Fprintf(w, "Working Dir:\t%s\n", config.WorkingDir)
		}
		if len(config.User) > 0 {
			fmt.Fprintf(w, "User:\t%s\n", config.User)
		}
		ports := sets.NewString()
		for k := range config.ExposedPorts {
			ports.Insert(k)
		}
		if len(ports) > 0 {
			fmt.Fprintf(w, "Exposes Ports:\t%s\n", strings.Join(ports.List(), ", "))
		}
	}

	if config != nil && len(config.Env) > 0 {
		for i, env := range config.Env {
			if i == 0 {
				fmt.Fprintf(w, "%s\t%s\n", "Environment:", env)
			} else {
				fmt.Fprintf(w, "%s\t%s\n", "", env)
			}
		}
	}

	if config != nil && len(config.Labels) > 0 {
		var keys []string
		for k := range config.Labels {
			keys = append(keys, k)
		}
		sort.Strings(keys)
		for i, key := range keys {
			if i == 0 {
				fmt.Fprintf(w, "%s\t%s=%s\n", "Labels:", key, config.Labels[key])
			} else {
				fmt.Fprintf(w, "%s\t%s=%s\n", "", key, config.Labels[key])
			}
		}
	}

	if config != nil && len(config.Volumes) > 0 {
		var keys []string
		for k := range config.Volumes {
			keys = append(keys, k)
		}
		sort.Strings(keys)
		for i, volume := range keys {
			if i == 0 {
				fmt.Fprintf(w, "%s\t%s\n", "Volumes:", volume)
			} else {
				fmt.Fprintf(w, "%s\t%s\n", "", volume)
			}
		}
	}

	fmt.Fprintln(w)
	return err
}

func writeTabSection(out io.Writer, fn func(w io.Writer)) {
	w := tabwriter.NewWriter(out, 0, 4, 1, ' ', 0)
	fn(w)
	w.Flush()
}

type ImageRetriever struct {
	FileDir         string
	SecurityOptions imagemanifest.SecurityOptions
	ParallelOptions imagemanifest.ParallelOptions
	// ImageMetadataCallback is invoked once per image retrieved, and may be called in parallel if
	// MaxPerRegistry is set higher than 1. If err is passed image is nil. If an error is returned
	// execution will stop.
	ImageMetadataCallback func(from string, image *Image, err error) error
	// ManifestListCallback, if specified, is invoked if the root image is a manifest list. If an
	// error returned processing stops. If zero manifests are returned the next item is rendered
	// and no ImageMetadataCallback calls occur. If more than one manifest is returned
	// ImageMetadataCallback will be invoked once for each item.
	ManifestListCallback func(from string, list *manifestlist.DeserializedManifestList, all map[digest.Digest]distribution.Manifest) (map[digest.Digest]distribution.Manifest, error)
}

// Image returns a single image matching ref.
func (o *ImageRetriever) Image(ctx context.Context, ref imagesource.TypedImageReference) (*Image, error) {
	images, err := o.Images(ctx, map[string]imagesource.TypedImageReference{"": ref})
	if err != nil {
		return nil, err
	}
	return images[""], nil
}

// Images invokes the retriever as specified and returns both the result of callbacks and a map
// of images invoked. It takes a value receiver because it mutates the original object.
func (o *ImageRetriever) Images(ctx context.Context, refs map[string]imagesource.TypedImageReference) (map[string]*Image, error) {
	fromContext, err := o.SecurityOptions.Context()
	if err != nil {
		return nil, err
	}
	fromOptions := &imagesource.Options{
		FileDir:         o.FileDir,
		Insecure:        o.SecurityOptions.Insecure,
		RegistryContext: fromContext,
	}

	var lock sync.Mutex
	images := make(map[string]*Image)
	callbackFn := func(name string, image *Image, err error) error {
		if o.ImageMetadataCallback != nil {
			if err := o.ImageMetadataCallback(name, image, err); err != nil {
				return err
			}
		}
		lock.Lock()
		defer lock.Unlock()
		images[name] = image
		return err
	}

	stopCh := make(chan struct{})
	defer close(stopCh)
	q := workqueue.New(o.ParallelOptions.MaxPerRegistry, stopCh)
	return images, q.Try(func(q workqueue.Try) {
		for key := range refs {
			name := key
			from := refs[key]
			q.Try(func() error {
				repo, err := fromOptions.Repository(ctx, from)
				if err != nil {
					return callbackFn(name, nil, fmt.Errorf("unable to connect to image repository %s: %v", from, err))
				}

				allManifests, manifestList, listDigest, err := imagemanifest.AllManifests(ctx, from.Ref, repo)
				if err != nil {
					if imagemanifest.IsImageForbidden(err) {
						msg := fmt.Sprintf("image %q does not exist or you don't have permission to access the repository", from)
						return callbackFn(name, nil, imagemanifest.NewImageForbidden(msg, err))
					}
					if imagemanifest.IsImageNotFound(err) {
						msg := fmt.Sprintf("image %q not found: %s", from, err.Error())
						return callbackFn(name, nil, imagemanifest.NewImageNotFound(msg, err))
					}
					return callbackFn(name, nil, fmt.Errorf("unable to read image %s: %v", from, err))
				}

				if o.ManifestListCallback != nil && manifestList != nil {
					allManifests, err = o.ManifestListCallback(name, manifestList, allManifests)
					if err != nil {
						return err
					}
				}

				if len(allManifests) == 0 {
					return imagemanifest.NewImageNotFound(fmt.Sprintf("no manifests could be found for %q", from), nil)
				}

				for srcDigest, srcManifest := range allManifests {
					contentDigest, contentErr := registryclient.ContentDigestForManifest(srcManifest, srcDigest.Algorithm())
					if contentErr != nil {
						return callbackFn(name, nil, contentErr)
					}

					imageConfig, layers, manifestErr := imagemanifest.ManifestToImageConfig(ctx, srcManifest, repo.Blobs(ctx), imagemanifest.ManifestLocation{ManifestList: listDigest, Manifest: srcDigest})
					mediaType, _, _ := srcManifest.Payload()
					if err := callbackFn(name, &Image{
						Name:          from.Ref.Exact(),
						Ref:           from,
						MediaType:     mediaType,
						Digest:        srcDigest,
						ContentDigest: contentDigest,
						ListDigest:    listDigest,
						Config:        imageConfig,
						Layers:        layers,
						Manifest:      srcManifest,
					}, manifestErr); err != nil {
						return err
					}
				}
				return nil
			})
		}
	})
}
