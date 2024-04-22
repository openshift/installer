package release

import (
	"archive/tar"
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"path"
	"sort"
	"strconv"
	"strings"
	"sync"
	"text/tabwriter"
	"time"

	"github.com/MakeNowJust/heredoc"
	"github.com/blang/semver"
	"github.com/distribution/distribution/v3"
	"github.com/distribution/distribution/v3/manifest/manifestlist"
	units "github.com/docker/go-units"
	digest "github.com/opencontainers/go-digest"
	"github.com/spf13/cobra"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/transport"
	"k8s.io/klog/v2"

	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/util/duration"
	"k8s.io/apimachinery/pkg/util/sets"
	"k8s.io/cli-runtime/pkg/genericclioptions"
	"k8s.io/cli-runtime/pkg/genericiooptions"
	kcmdutil "k8s.io/kubectl/pkg/cmd/util"
	"k8s.io/kubectl/pkg/util/templates"
	"sigs.k8s.io/yaml"

	configv1 "github.com/openshift/api/config/v1"
	imageapi "github.com/openshift/api/image/v1"
	configv1client "github.com/openshift/client-go/config/clientset/versioned"
	"github.com/openshift/library-go/pkg/features"
	"github.com/openshift/library-go/pkg/image/dockerv1client"
	imagereference "github.com/openshift/library-go/pkg/image/reference"
	"github.com/openshift/oc/pkg/cli/image/extract"
	"github.com/openshift/oc/pkg/cli/image/imagesource"
	imageinfo "github.com/openshift/oc/pkg/cli/image/info"
	imagemanifest "github.com/openshift/oc/pkg/cli/image/manifest"
)

func NewInfoOptions(streams genericiooptions.IOStreams) *InfoOptions {
	return &InfoOptions{
		IOStreams:              streams,
		KubeTemplatePrintFlags: *genericclioptions.NewKubeTemplatePrintFlags(),
		ParallelOptions:        imagemanifest.ParallelOptions{MaxPerRegistry: 4},
	}
}

func NewInfo(f kcmdutil.Factory, streams genericiooptions.IOStreams) *cobra.Command {
	o := NewInfoOptions(streams)
	cmd := &cobra.Command{
		Use:   "info IMAGE [--changes-from=IMAGE] [--verify|--commits|--pullspecs]",
		Short: "Display information about a release",
		Long: templates.LongDesc(`
			Show information about an OpenShift release.

			This command retrieves, verifies, and formats the information describing an OpenShift update.
			Updates are delivered as container images with metadata describing the component images and
			the configuration necessary to install the system operators. A release image is usually
			referenced via its content digest, which allows this command and the update infrastructure to
			validate that updates have not been tampered with.

			If no arguments are specified the release of the currently connected cluster is displayed.
			Specify one or more images via pull spec to see details of each release image. You may also
			pass a semantic version (4.11.2) as an argument, and if cluster version object has seen such a
			version in the upgrades channel it will find the release info for that version.

			The --commits flag will display the Git commit IDs and repository URLs for the source of each
			component image. The --pullspecs flag will display the full component image pull spec. --size
			will show a breakdown of each image, their layers, and the total size of the payload.
			--contents shows the configuration that will be applied to the cluster when the update is run.
			If you have specified two images the difference between the first and second image will be
			shown. You may use -o name, -o digest, or -o pullspec to output the tag name, digest for
			image, or pullspec of the images referenced in the release image.

			The --verify flag will display one summary line per input release image and verify the
			integrity of each. The command will return an error if the release has been tampered with.
			Passing a pull spec with a digest (e.g. quay.io/openshift/release@sha256:a9bc...) instead of
			a tag when verifying an image is recommended since it ensures an attacker cannot trick you
			into installing an older, potentially vulnerable version.

			The --bugs and --changelog flags will use git to clone the source of the release and display
			the code changes that occurred between the two release arguments. This operation is slow
			and requires sufficient disk space on the selected drive to clone all repositories.

			If the specified image supports multiple operating systems, the image that matches the
			current operating system will be chosen. Otherwise you must pass --filter-by-os to
			select the desired image.
		`),
		Example: templates.Examples(`
			# Show information about the cluster's current release
			oc adm release info

			# Show the source code that comprises a release
			oc adm release info 4.11.2 --commit-urls

			# Show the source code difference between two releases
			oc adm release info 4.11.0 4.11.2 --commits

			# Show where the images referenced by the release are located
			oc adm release info quay.io/openshift-release-dev/ocp-release:4.11.2 --pullspecs

			# Show information about linux/s390x image
			# Note: Wildcard filter is not supported; pass a single os/arch to extract
			oc adm release info quay.io/openshift-release-dev/ocp-release:4.11.2 --filter-by-os=linux/s390x

		`),
		Run: func(cmd *cobra.Command, args []string) {
			kcmdutil.CheckErr(o.Complete(f, cmd, args))
			kcmdutil.CheckErr(o.Validate())
			kcmdutil.CheckErr(o.Run())
		},
	}
	flags := cmd.Flags()
	o.SecurityOptions.Bind(flags)
	o.FilterOptions.Bind(flags)
	o.ParallelOptions.Bind(flags)
	o.KubeTemplatePrintFlags.AddFlags(cmd)

	flags.StringVar(&o.From, "changes-from", o.From, "Show changes from this image to the requested image.")

	flags.BoolVar(&o.Verify, "verify", o.Verify, "Generate bug listings from the changelogs in the git repositories extracted to this path.")

	flags.StringVar(&o.ICSPFile, "icsp-file", o.ICSPFile, "Path to an ImageContentSourcePolicy file. If set, data from this file will be used to find alternative locations for images.")
	flags.MarkDeprecated("icsp-file", "support for it will be removed in a future release. Use --idms-file instead.")
	flags.StringVar(&o.IDMSFile, "idms-file", o.IDMSFile, "Path to an ImageDigestMirrorSet file. If set, data from this file will be used to find alternative locations for images.")

	flags.BoolVar(&o.ShowContents, "contents", o.ShowContents, "Display the contents of a release.")
	flags.BoolVar(&o.ShowCommit, "commits", o.ShowCommit, "Display information about the source an image was created with.")
	flags.BoolVar(&o.ShowCommitURL, "commit-urls", o.ShowCommitURL, "Display a link (if possible) to the source code.")
	flags.BoolVar(&o.ShowPullSpec, "pullspecs", o.ShowPullSpec, "Display the pull spec of each image instead of the digest.")
	flags.BoolVar(&o.ShowSize, "size", o.ShowSize, "Display the size of each image including overlap.")
	flags.StringVar(&o.ImageFor, "image-for", o.ImageFor, "Print the pull spec of the specified image or an error if it does not exist.")
	flags.StringVarP(&o.Output, "output", "o", o.Output, "Display the release info in an alternative format: digest|json|name|pullspec|template|jsonpath.")
	flags.StringVar(&o.ChangelogDir, "changelog", o.ChangelogDir, "Generate changelog output from the git directories extracted to this path.")
	flags.StringVar(&o.BugsDir, "bugs", o.BugsDir, "Generate bug listings from the changelogs in the git repositories extracted to this path.")
	flags.BoolVar(&o.IncludeImages, "include-images", o.IncludeImages, "When displaying JSON output of a release output the images the release references.")
	flags.StringVar(&o.FileDir, "dir", o.FileDir, "The directory on disk that file:// images will be copied under.")
	flags.BoolVar(&o.SkipBugCheck, "skip-bug-check", o.SkipBugCheck, "Do not check bug statuses when running generating bug listing with --output=name")
	return cmd
}

type InfoOptions struct {
	genericiooptions.IOStreams
	genericclioptions.KubeTemplatePrintFlags

	Images  []string
	From    string
	FileDir string

	Output        string
	ImageFor      string
	IncludeImages bool
	ShowContents  bool
	ShowCommit    bool
	ShowCommitURL bool
	ShowPullSpec  bool
	ShowSize      bool
	Verify        bool
	ICSPFile      string
	IDMSFile      string

	ChangelogDir string
	BugsDir      string
	SkipBugCheck bool

	ParallelOptions imagemanifest.ParallelOptions
	SecurityOptions imagemanifest.SecurityOptions
	FilterOptions   imagemanifest.FilterOptions
}

func findSemanticVersionArgs(args []string) map[string]semver.Version {
	semvers := make(map[string]semver.Version)
	for _, arg := range args {
		v, err := semver.Parse(arg)
		if err != nil {
			continue
		}
		semvers[arg] = v
	}
	return semvers
}

func findImageInClusterVersion(cv *configv1.ClusterVersion, arg string) (string, bool) {
	if cv.Status.Desired.Version == arg && len(cv.Status.Desired.Image) > 0 {
		return cv.Status.Desired.Image, true
	}
	for _, available := range cv.Status.AvailableUpdates {
		if available.Version == arg && len(available.Image) > 0 {
			return available.Image, true
		}
	}
	for _, history := range cv.Status.History {
		if history.Version == arg && len(history.Image) > 0 {
			return history.Image, true
		}
	}
	return "", false
}

type versionNode struct {
	Version string `json:"version"`
	Payload string `json:"payload"`
}

type versionGraph struct {
	Nodes []versionNode `json:"nodes"`
}

const defaultGraphURL = "https://api.openshift.com/api/upgrades_info/v1/graph"

// replaceStableSemanticArgs attempts to look up known major versions in existing public stable
// channels.
// TODO: perfom graph lookups from the cluster's graph endpoint and channel in preference
func replaceStableSemanticArgs(args []string, semanticArgs map[string]semver.Version, graphURL string) error {
	if len(graphURL) == 0 {
		graphURL = defaultGraphURL
	}
	u, err := url.Parse(graphURL)
	if err != nil {
		return err
	}

	transport, err := transport.HTTPWrappersForConfig(
		&transport.Config{
			UserAgent: rest.DefaultKubernetesUserAgent() + "(release-info)",
		},
		http.DefaultTransport,
	)
	if err != nil {
		return err
	}
	client := &http.Client{Transport: transport}

	for i, arg := range args {
		v, ok := semanticArgs[arg]
		if !ok {
			continue
		}
		if v.Major != 4 {
			continue
		}

		var found bool
		for _, stream := range []string{"fast", "stable", "candidate"} {
			u.RawQuery = url.Values{"channel": []string{fmt.Sprintf("%s-%d.%d", stream, v.Major, v.Minor)}}.Encode()
			if err := func() error {
				req, err := http.NewRequest("GET", u.String(), nil)
				if err != nil {
					return err
				}
				req.Header.Set("Accept", "application/json")
				resp, err := client.Do(req)
				if err != nil {
					return err
				}
				defer resp.Body.Close()
				switch resp.StatusCode {
				case http.StatusOK:
				default:
					io.Copy(io.Discard, resp.Body)
					return fmt.Errorf("unable to retrieve status for %q: %d", arg, resp.StatusCode)
				}
				data, err := io.ReadAll(resp.Body)
				if err != nil {
					return err
				}
				var versions versionGraph
				if err := json.Unmarshal(data, &versions); err != nil {
					return err
				}
				for _, version := range versions.Nodes {
					if version.Version == arg && len(version.Payload) > 0 {
						delete(semanticArgs, arg)
						args[i] = version.Payload
						found = true
						break
					}
				}
				return nil
			}(); err != nil {
				return err
			}
			if found {
				break
			}
		}
	}
	return nil
}

func replaceClusterSemanticArgs(f kcmdutil.Factory, args []string, semanticArgs map[string]semver.Version) ([]string, error) {
	cfg, err := f.ToRESTConfig()
	if err != nil {
		return args, fmt.Errorf("info expects one argument, or a connection to an OpenShift 4.x server: %v", err)
	}
	client, err := configv1client.NewForConfig(cfg)
	if err != nil {
		return args, fmt.Errorf("info expects one argument, or a connection to an OpenShift 4.x server: %v", err)
	}
	cv, err := client.ConfigV1().ClusterVersions().Get(context.TODO(), "version", metav1.GetOptions{})
	if err != nil {
		if errors.IsNotFound(err) || errors.ReasonForError(err) == metav1.StatusReasonUnknown {
			klog.V(2).Infof("Unable to find cluster version object from cluster: %v", err)
			return args, fmt.Errorf("info expects one argument, or a connection to an OpenShift 4.x server")
		}
		return args, fmt.Errorf("info expects one argument, or a connection to an OpenShift 4.x server: %v", err)
	}

	if len(args) == 0 {
		image := cv.Status.Desired.Image
		if len(image) == 0 && cv.Spec.DesiredUpdate != nil {
			image = cv.Spec.DesiredUpdate.Image
		}
		if len(image) == 0 {
			return nil, fmt.Errorf("the server is not reporting a release image at this time, please specify an image to view")
		}
		return []string{image}, nil
	}

	for i, arg := range args {
		if _, ok := semanticArgs[arg]; !ok {
			continue
		}
		image, ok := findImageInClusterVersion(cv, arg)
		if !ok {
			continue
		}
		delete(semanticArgs, arg)
		klog.V(2).Infof("Replaced argument %q with %q", arg, image)
		args[i] = image
	}
	return args, nil
}

func findArgumentsFromCluster(f kcmdutil.Factory, args []string) ([]string, error) {
	semanticArgs := findSemanticVersionArgs(args)
	if len(semanticArgs) == 0 && len(args) > 0 {
		return args, nil
	}
	klog.V(4).Infof("Found semantic versions: %v", semanticArgs)
	// attempt to find semantic args from the cluster
	args, clusterErr := replaceClusterSemanticArgs(f, args, semanticArgs)
	if len(semanticArgs) == 0 {
		return args, clusterErr
	}
	// if any semantic args remain, try to fetch them from the api endpoint out of a stable channel
	err := replaceStableSemanticArgs(args, semanticArgs, defaultGraphURL)
	if len(semanticArgs) == 0 || err != nil {
		if clusterErr != nil {
			klog.V(2).Infof("Ignored error retrieving semantic versions from cluster version: %v", err)
		}
		return args, err
	}
	// if there are any semantic args left, error
	for arg := range semanticArgs {
		return nil, fmt.Errorf("the semantic version %q is not present in the cluster version status or in the official versions list, cannot be resolved", arg)
	}
	return args, nil
}

func (o *InfoOptions) Complete(f kcmdutil.Factory, cmd *cobra.Command, args []string) error {
	args, err := findArgumentsFromCluster(f, args)
	if err != nil {
		return err
	}
	if len(args) < 1 {
		return fmt.Errorf("info expects at least one argument, a release image pull spec")
	}
	o.Images = args
	if len(o.From) == 0 && len(o.Images) == 2 && !o.Verify {
		o.From = o.Images[0]
		o.Images = o.Images[1:]
	}
	return o.FilterOptions.Complete(cmd.Flags())
}

func (o *InfoOptions) Validate() error {
	count := 0
	if len(o.ImageFor) > 0 {
		count++
	}
	if o.ShowCommit {
		count++
	}
	if o.ShowCommitURL {
		count++
	}
	if o.ShowPullSpec {
		count++
	}
	if o.ShowContents {
		count++
	}
	if o.ShowSize {
		count++
	}
	if o.Verify {
		count++
	}
	if count > 1 {
		return fmt.Errorf("only one of --commits, --commit-urls, --pullspecs, --contents, --size, --verify may be specified")
	}
	if len(o.ImageFor) > 0 && len(o.Output) > 0 {
		return fmt.Errorf("--output and --image-for may not both be specified")
	}
	if o.SkipBugCheck && len(o.BugsDir) == 0 {
		return fmt.Errorf("--skip-bug-check requires --bugs")
	}
	if o.SkipBugCheck && o.Output != "name" && o.Output != "json" {
		return fmt.Errorf("--skip-bug-check requires --output to be set to 'name' or 'json'")
	}
	if len(o.ChangelogDir) > 0 || len(o.BugsDir) > 0 {
		if len(o.From) == 0 {
			return fmt.Errorf("--changelog/--bugs require --changes-from")
		}
	}
	if len(o.ChangelogDir) > 0 && len(o.BugsDir) > 0 {
		return fmt.Errorf("--changelog and --bugs may not both be specified")
	}
	switch {
	case len(o.BugsDir) > 0:
		switch o.Output {
		case "", "name", "json":
		default:
			return fmt.Errorf("--output only supports 'name' or 'json' for --bugs")
		}
	case len(o.ChangelogDir) > 0:
		switch o.Output {
		case "", "json":
		default:
			return fmt.Errorf("--output only supports 'json' for --changelog")
		}
	default:
		output := strings.SplitN(o.Output, "=", 2)[0]
		if len(output) > 0 && !stringArrContains(o.allowedFormats(), output) {
			return fmt.Errorf("--output only supports %s", strings.Join(o.allowedFormats(), ", "))
		}
	}

	if len(o.Images) == 0 {
		return fmt.Errorf("must specify a release image as an argument")
	}
	if len(o.From) > 0 && len(o.Images) != 1 {
		return fmt.Errorf("must specify a single release image as argument when comparing to another release image")
	}

	if len(o.ICSPFile) > 0 && len(o.IDMSFile) > 0 {
		return fmt.Errorf("icsp-file and idms-file are mutually exclusive")
	}

	return o.FilterOptions.Validate()
}

func (o *InfoOptions) Run() error {
	fetchImages := o.ShowSize || o.Verify || o.IncludeImages

	if len(o.From) > 0 && !o.Verify {
		if o.ShowContents {
			return diffContents(o.From, o.Images[0], o.Out)
		}

		var baseRelease *ReleaseInfo
		var baseErr error
		done := make(chan struct{})
		go func() {
			defer close(done)
			baseRelease, baseErr = o.LoadReleaseInfo(o.From, fetchImages)
		}()

		release, err := o.LoadReleaseInfo(o.Images[0], fetchImages)
		if err != nil {
			return err
		}

		<-done
		if baseErr != nil {
			return baseErr
		}

		diff, err := calculateDiff(baseRelease, release)
		if err != nil {
			return err
		}
		if len(o.BugsDir) > 0 {
			return describeBugs(o.Out, o.ErrOut, diff, o.BugsDir, o.Output, o.SkipBugCheck)
		}
		if len(o.ChangelogDir) > 0 {
			return describeChangelog(o.Out, o.ErrOut, release, diff, o.ChangelogDir, o.Output)
		}
		return describeReleaseDiff(o.Out, diff, o.ShowCommit, o.Output)
	}

	var exitErr error
	for _, image := range o.Images {
		release, err := o.LoadReleaseInfo(image, fetchImages)
		if err != nil {
			exitErr = kcmdutil.ErrExit
			fmt.Fprintf(o.ErrOut, "error: %v\n", err)
			continue
		}
		if o.Verify {
			fmt.Fprintf(o.Out, "%s %s %s\n", release.Digest, release.References.CreationTimestamp.UTC().Format(time.RFC3339), release.PreferredName())
			continue
		}
		if err := o.describeImage(release); err != nil {
			exitErr = kcmdutil.ErrExit
			fmt.Fprintf(o.ErrOut, "error: %v\n", err)
			continue
		}
	}
	return exitErr
}

func (opt *InfoOptions) allowedFormats() []string {
	formats := []string{"json", "pullspec", "digest", "name"}
	formats = append(formats, opt.KubeTemplatePrintFlags.AllowedFormats()...)
	return formats
}

func diffContents(a, b string, out io.Writer) error {
	fmt.Fprintf(out, `To see the differences between these releases, run:

  %[1]s adm release extract %[2]s --to=/tmp/old
  %[1]s adm release extract %[3]s --to=/tmp/new
  diff /tmp/old /tmp/new

`, os.Args[0], a, b)
	return nil
}

func (o *InfoOptions) describeImage(release *ReleaseInfo) error {
	if o.ShowContents {
		_, err := io.Copy(o.Out, newContentStreamForRelease(release))
		return err
	}
	output := strings.SplitN(o.Output, "=", 2)
	switch output[0] {
	case "json":
		data, err := json.MarshalIndent(release, "", "  ")
		if err != nil {
			return err
		}
		fmt.Fprintln(o.Out, string(data))
		return nil
	case "name":
		for _, tag := range release.References.Spec.Tags {
			fmt.Fprintf(o.Out, "%s\n", tag.Name)
		}
		return nil
	case "pullspec":
		for _, tag := range release.References.Spec.Tags {
			if tag.From != nil && tag.From.Kind == "DockerImage" {
				fmt.Fprintf(o.Out, "%s\n", tag.From.Name)
			}
		}
		return nil
	case "digest":
		for _, tag := range release.References.Spec.Tags {
			if tag.From != nil && tag.From.Kind == "DockerImage" {
				if ref, err := imagereference.Parse(tag.From.Name); err != nil {
					fmt.Fprintf(o.ErrOut, "error: %s is not a valid reference: %v\n", tag.Name, err)
				} else if len(ref.ID) == 0 {
					fmt.Fprintf(o.ErrOut, "error: %s does not point to a digest\n", tag.Name)
				} else {
					fmt.Fprintf(o.Out, "%s\n", ref.ID)
				}
			}
		}
		return nil
	case "":
	default:
		p, err := o.KubeTemplatePrintFlags.ToPrinter(o.Output)
		if genericclioptions.IsNoCompatiblePrinterError(err) {
			return fmt.Errorf("output mode only supports %s", strings.Join(o.allowedFormats(), ", "))
		}
		if err != nil {
			return err
		}
		data, err := json.MarshalIndent(release, "", "  ")
		if err != nil {
			return err
		}
		return p.PrintObj(&runtime.Unknown{Raw: data}, o.Out)
	}
	if len(o.ImageFor) > 0 {
		spec, err := findImageSpec(release.References, o.ImageFor, release.Image)
		if err != nil {
			return err
		}
		fmt.Fprintln(o.Out, spec)
		return nil
	}
	return describeReleaseInfo(o.Out, release, o.ShowCommit, o.ShowCommitURL, o.ShowPullSpec, o.ShowSize)
}

func findImageSpec(image *imageapi.ImageStream, tagName, imageName string) (string, error) {
	for _, tag := range image.Spec.Tags {
		if tag.Name == tagName {
			if tag.From != nil && tag.From.Kind == "DockerImage" && len(tag.From.Name) > 0 {
				return tag.From.Name, nil
			}
		}
	}
	return "", fmt.Errorf("no image tag %q exists in the release image %s", tagName, imageName)
}

func calculateDiff(from, to *ReleaseInfo) (*ReleaseDiff, error) {
	diff := &ReleaseDiff{
		From:             from,
		To:               to,
		ChangedImages:    make(map[string]*ImageReferenceDiff),
		ChangedManifests: make(map[string]*ReleaseManifestDiff),
	}
	for i := range from.References.Spec.Tags {
		tag := &from.References.Spec.Tags[i]
		if tag.From == nil || tag.From.Kind != "DockerImage" {
			continue
		}
		diff.ChangedImages[tag.Name] = &ImageReferenceDiff{
			Name: tag.Name,
			From: tag,
		}
	}
	for i := range to.References.Spec.Tags {
		tag := &to.References.Spec.Tags[i]
		if tag.From == nil || tag.From.Kind != "DockerImage" {
			continue
		}
		if exists, ok := diff.ChangedImages[tag.Name]; ok {
			exists.To = tag
			continue
		}
		diff.ChangedImages[tag.Name] = &ImageReferenceDiff{
			Name: tag.Name,
			To:   tag,
		}
	}
	for k, v := range diff.ChangedImages {
		if v.From != nil && v.To != nil && v.From.From.Name == v.To.From.Name {
			delete(diff.ChangedImages, k)
		}
	}
	for name, manifest := range from.ManifestFiles {
		diff.ChangedManifests[name] = &ReleaseManifestDiff{
			Filename: name,
			From:     manifest,
		}
	}
	for name, manifest := range to.ManifestFiles {
		if exists, ok := diff.ChangedManifests[name]; ok {
			exists.To = manifest
			continue
		}
		diff.ChangedManifests[name] = &ReleaseManifestDiff{
			Filename: name,
			To:       manifest,
		}
	}
	for k, v := range diff.ChangedManifests {
		if bytes.Equal(v.From, v.To) {
			delete(diff.ChangedManifests, k)
		}
	}

	return diff, nil
}

type ReleaseDiff struct {
	From *ReleaseInfo `json:"from"`
	To   *ReleaseInfo `json:"to"`

	ChangedImages    map[string]*ImageReferenceDiff  `json:"changedImages"`
	ChangedManifests map[string]*ReleaseManifestDiff `json:"changedManifests"`
}

type ImageReferenceDiff struct {
	Name string `json:"name"`

	From *imageapi.TagReference `json:"from"`
	To   *imageapi.TagReference `json:"to"`
}

type ReleaseManifestDiff struct {
	Filename string `json:"filename"`

	From []byte `json:"from"`
	To   []byte `json:"to"`
}

type ReleaseInfo struct {
	Image              string                            `json:"image"`
	ImageRef           imagesource.TypedImageReference   `json:"-"`
	Digest             digest.Digest                     `json:"digest"`
	ContentDigest      digest.Digest                     `json:"contentDigest"`
	ManifestListDigest digest.Digest                     `json:"listDigest"`
	Config             *dockerv1client.DockerImageConfig `json:"config"`
	Metadata           *CincinnatiMetadata               `json:"metadata"`
	References         *imageapi.ImageStream             `json:"references"`

	// This field is deprecated, does not contain display names. Is replaced by
	// ComponentVersions.
	DeprecatedComponentVersions map[string]string `json:"versions"`
	ComponentVersions           ComponentVersions `json:"displayVersions"`

	Images map[string]*Image `json:"images"`

	RawMetadata   map[string][]byte `json:"-"`
	ManifestFiles map[string][]byte `json:"-"`
	UnknownFiles  []string          `json:"-"`

	Warnings []string `json:"warnings"`
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

func (i *ReleaseInfo) PreferredName() string {
	if i.Metadata != nil {
		return i.Metadata.Version
	}
	return i.References.Name
}

func (i *ReleaseInfo) Platform() string {
	config := i.Config
	if config == nil {
		config = &dockerv1client.DockerImageConfig{}
	}

	os := config.OS
	if len(os) == 0 {
		os = "unknown"
	}
	arch := config.Architecture
	if len(arch) == 0 {
		arch = "unknown"
	}
	if len(i.ManifestListDigest) == 0 {
		return fmt.Sprintf("%s/%s", os, arch)
	}
	return fmt.Sprintf("multi (%s/%s)", os, arch)
}

func (o *InfoOptions) LoadReleaseInfo(image string, retrieveImages bool) (*ReleaseInfo, error) {
	ref, err := imagesource.ParseReference(image)
	if err != nil {
		return nil, err
	}

	verifier := imagemanifest.NewVerifier()
	opts := extract.NewExtractOptions(genericiooptions.IOStreams{Out: o.Out, ErrOut: o.ErrOut})
	opts.SecurityOptions = o.SecurityOptions
	opts.FilterOptions = o.FilterOptions
	opts.FileDir = o.FileDir
	opts.ICSPFile = o.ICSPFile
	opts.IDMSFile = o.IDMSFile

	release := &ReleaseInfo{
		Image:    image,
		ImageRef: ref,

		RawMetadata: make(map[string][]byte),
	}

	opts.ImageMetadataCallback = func(m *extract.Mapping, dgst, contentDigest digest.Digest, config *dockerv1client.DockerImageConfig, manifestListDigest digest.Digest) {
		verifier.Verify(dgst, contentDigest)
		release.ManifestListDigest = manifestListDigest
		release.Digest = dgst
		release.ContentDigest = contentDigest
		release.Config = config
	}
	opts.OnlyFiles = true
	opts.Mappings = []extract.Mapping{
		{
			ImageRef: ref,

			From:        "release-manifests/",
			To:          ".",
			LayerFilter: extract.NewPositionLayerFilter(-1),
		},
	}
	var errs []error
	opts.TarEntryCallback = func(hdr *tar.Header, _ extract.LayerInfo, r io.Reader) (bool, error) {
		switch hdr.Name {
		case "image-references":
			data, err := io.ReadAll(r)
			if err != nil {
				errs = append(errs, fmt.Errorf("unable to read release image-references: %v", err))
				return true, nil
			}
			release.RawMetadata[hdr.Name] = data
			is, err := readReleaseImageReferences(data)
			if err != nil {
				errs = append(errs, err)
				return true, nil
			}
			release.References = is
		case "release-metadata":
			data, err := io.ReadAll(r)
			if err != nil {
				errs = append(errs, fmt.Errorf("unable to read release metadata: %v", err))
				return true, nil
			}
			release.RawMetadata[hdr.Name] = data
			m := &CincinnatiMetadata{}
			if err := json.Unmarshal(data, m); err != nil {
				errs = append(errs, fmt.Errorf("invalid release metadata: %v", err))
				return true, nil
			}
			release.Metadata = m
		default:
			if ext := path.Ext(hdr.Name); len(ext) > 0 && (ext == ".yaml" || ext == ".yml" || ext == ".json") {
				klog.V(4).Infof("Found manifest %s", hdr.Name)
				data, err := io.ReadAll(r)
				if err != nil {
					errs = append(errs, fmt.Errorf("unable to read release manifest %q: %v", hdr.Name, err))
					return true, nil
				}
				if release.ManifestFiles == nil {
					release.ManifestFiles = make(map[string][]byte)
				}
				release.ManifestFiles[hdr.Name] = data
			} else {
				release.UnknownFiles = append(release.UnknownFiles, hdr.Name)
			}
		}
		return true, nil
	}
	if err := opts.Run(); err != nil {
		return nil, err
	}
	if len(errs) > 0 {
		return nil, fmt.Errorf("release image could not be read: %s", errorList(errs))
	}

	if release.References == nil {
		return nil, fmt.Errorf("release image did not contain an image-references file")
	}

	release.ComponentVersions, errs = readComponentVersions(release.References, o.ErrOut)
	for _, err := range errs {
		release.Warnings = append(release.Warnings, err.Error())
	}

	if retrieveImages {
		var lock sync.Mutex
		release.Images = make(map[string]*Image)
		images := make(map[string]imagesource.TypedImageReference)
		r := &imageinfo.ImageRetriever{
			FileDir:         opts.FileDir,
			SecurityOptions: o.SecurityOptions,
			ParallelOptions: o.ParallelOptions,
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

				return nil, fmt.Errorf("no matching image found in manifest list")
			},
			ImageMetadataCallback: func(name string, image *imageinfo.Image, err error) error {
				if image != nil {
					verifier.Verify(image.Digest, image.ContentDigest)
				}
				lock.Lock()
				defer lock.Unlock()
				if err != nil {
					release.Warnings = append(release.Warnings, fmt.Sprintf("tag %q: %v", name, err))
					return nil
				}
				copied := Image(*image)
				release.Images[name] = &copied
				return nil
			},
		}
		for _, tag := range release.References.Spec.Tags {
			if tag.From == nil || tag.From.Kind != "DockerImage" {
				continue
			}
			ref, err := imagereference.Parse(tag.From.Name)
			if err != nil {
				release.Warnings = append(release.Warnings, fmt.Sprintf("tag %q has an invalid reference: %v", tag.Name, err))
				continue
			}
			images[tag.Name] = imagesource.TypedImageReference{Type: imagesource.DestinationRegistry, Ref: ref}
		}
		if _, err := r.Images(context.TODO(), images); err != nil {
			return nil, err
		}
	}

	if !verifier.Verified() {
		err := fmt.Errorf("the release image failed content verification and may have been tampered with")
		if !o.SecurityOptions.SkipVerification {
			return nil, err
		}
		fmt.Fprintf(o.ErrOut, "warning: %v\n", err)
	}

	sort.Strings(release.Warnings)

	return release, nil
}

func readComponentVersions(is *imageapi.ImageStream, errOut io.Writer) (ComponentVersions, []error) {
	var errs []error
	combined := make(map[string]sets.String)
	combinedDisplayNames := make(map[string]sets.String)
	kubectlVersions := make(map[string]sets.Set[string])
	for _, tag := range is.Spec.Tags {
		versions, ok := tag.Annotations[annotationBuildVersions]
		if !ok {
			continue
		}
		all, err := parseComponentVersionsLabel(versions, tag.Annotations[annotationBuildVersionsDisplayNames])
		if err != nil {
			errs = append(errs, fmt.Errorf("the referenced image %s had an invalid version annotation: %v", tag.Name, err))
		}
		for k, v := range all {
			if k == "kubectl" {
				if val, exist := kubectlVersions[v.Version]; exist {
					kubectlVersions[v.Version] = val.Insert(tag.Name)
				} else {
					kubectlVersions[v.Version] = sets.New[string](tag.Name)
				}
				if tag.Name != "cli" && tag.Name != "cli-artifacts" {
					continue
				}
			}
			existing, ok := combined[k]
			if !ok {
				existing = sets.NewString()
				combined[k] = existing
			}
			existing.Insert(v.Version)

			existingDisplayName, ok := combinedDisplayNames[k]
			if !ok {
				existingDisplayName = sets.NewString()
				combinedDisplayNames[k] = existingDisplayName
			}
			existingDisplayName.Insert(v.DisplayName)
		}
	}

	multiples := sets.NewString()
	if len(kubectlVersions) > 2 {
		var kubectlFormattedWarning []string
		for k, v := range kubectlVersions {
			kubectlFormattedWarning = append(kubectlFormattedWarning, fmt.Sprintf("%s (%s)", k, strings.Join(v.UnsortedList(), ", ")))
		}

		fmt.Fprintf(errOut, "warning: multiple versions reported for the kubectl: %s\n", strings.Join(kubectlFormattedWarning, ", "))
	}
	var out ComponentVersions
	var keys []string
	for k := range combined {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	for _, k := range keys {
		v := combined[k]
		if v.Len() > 1 {
			multiples = multiples.Insert(k)
		}
		if _, ok := out[k]; ok {
			continue
		}
		version := v.List()[0]
		if out == nil {
			out = make(ComponentVersions)
		}
		out[k] = ComponentVersion{Version: version}
	}
	for _, k := range keys {
		v, ok := combinedDisplayNames[k]
		if !ok {
			continue
		}
		if v.Len() > 1 {
			multiples = multiples.Insert(k)
		}
		version, ok := out[k]
		if !ok {
			continue
		}
		if len(version.DisplayName) == 0 {
			version.DisplayName = v.List()[0]
		}
		out[k] = version
	}

	if len(multiples) > 0 {
		errs = append(errs, fmt.Errorf("multiple versions or display names reported for the following component(s): %v", strings.Join(multiples.List(), ",  ")))
	}
	return out, errs
}

func errorList(errs []error) string {
	if len(errs) == 1 {
		return errs[0].Error()
	}
	buf := &bytes.Buffer{}
	fmt.Fprintf(buf, "\n\n")
	for _, err := range errs {
		fmt.Fprintf(buf, "* %v\n", err)
	}
	return buf.String()
}

func stringArrContains(arr []string, s string) bool {
	for _, item := range arr {
		if item == s {
			return true
		}
	}
	return false
}

func describeReleaseDiff(out io.Writer, diff *ReleaseDiff, showCommit bool, outputMode string) error {
	switch outputMode {
	case "json":
		data, err := json.MarshalIndent(diff, "", "  ")
		if err != nil {
			return err
		}
		fmt.Fprintln(out, string(data))
		return nil
	case "":
		// print human readable output
	default:
		return fmt.Errorf("unrecognized output mode: %s", outputMode)
	}
	if diff.To.Digest == diff.From.Digest {
		fmt.Fprintf(out, "Releases are identical\n")
		return nil
	}
	w := tabwriter.NewWriter(out, 0, 4, 2, ' ', 0)
	defer w.Flush()
	now := time.Now()
	fmt.Fprintf(w, "\tFROM\tTO\n")
	fmt.Fprintf(w, "Name:\t%s\t%s\n", diff.From.PreferredName(), diff.To.PreferredName())
	fmt.Fprintf(w, "Created:\t%s\t%s\n", duration.ShortHumanDuration(now.Sub(diff.From.Config.Created)), duration.ShortHumanDuration(now.Sub(diff.To.Config.Created)))
	if from, to := diff.From.Platform(), diff.To.Platform(); from != to {
		fmt.Fprintf(w, "OS/Arch:\t%s\t%s\n", from, to)
	}

	switch {
	case diff.From.Metadata != nil && diff.To.Metadata != nil:
		fmt.Fprintln(w)
		fmt.Fprintf(w, "Version:\t%s\t%s\n", diff.From.Metadata.Version, diff.To.Metadata.Version)
		canUpgrade := "No"
		if stringArrContains(diff.To.Metadata.Previous, diff.From.Metadata.Version) {
			canUpgrade = "Yes"
		}
		fmt.Fprintf(w, "Upgrade From:\t\t%s\n", canUpgrade)
	case diff.From.Metadata != nil && diff.To.Metadata == nil:
		fmt.Fprintf(w, "Has Release Metadata:\tYes\t\n")
	case diff.From.Metadata == nil && diff.To.Metadata != nil:
		fmt.Fprintf(w, "Has Release Metadata:\t\tYes\n")
	}

	if len(diff.ChangedImages) > 0 {
		var keys []string
		maxLen := 0
		for k := range diff.ChangedImages {
			if len(k) > maxLen {
				maxLen = len(k)
			}
			keys = append(keys, k)
		}
		justify := func(s string) string {
			return s + strings.Repeat(" ", maxLen-len(s))
		}
		sort.Strings(keys)
		var rebuilt []string
		writeTabSection(w, func(w io.Writer) {
			count := 0
			for _, k := range keys {
				if image := diff.ChangedImages[k]; image.To != nil && image.From != nil {
					if !codeChanged(image.From, image.To) {
						rebuilt = append(rebuilt, k)
						continue
					}
					if count == 0 {
						fmt.Fprintln(w)
						fmt.Fprintf(w, "Images Changed:\n")
					}
					count++
					old, new := digestOrRef(image.From.From.Name), digestOrRef(image.To.From.Name)
					if old != new {
						if showCommit {
							fmt.Fprintf(w, "  %s\t%s\n", justify(image.Name), gitDiffOrCommit(image.From, image.To))
						} else {
							fmt.Fprintf(w, "  %s\t%s\t%s\n", justify(image.Name), old, new)
						}
					}
				}
			}
		})

		if len(rebuilt) > 0 {
			writeTabSection(w, func(w io.Writer) {
				count := 0
				for _, k := range rebuilt {
					if image := diff.ChangedImages[k]; image.To != nil && image.From != nil {
						if count == 0 {
							fmt.Fprintln(w)
							fmt.Fprintf(w, "Images Rebuilt:\n")
						}
						count++
						old, new := digestOrRef(image.From.From.Name), digestOrRef(image.To.From.Name)
						if old != new {
							if showCommit {
								fmt.Fprintf(w, "  %s\t%s\n", justify(image.Name), gitDiffOrCommit(image.From, image.To))
							} else {
								fmt.Fprintf(w, "  %s\t%s\t%s\n", justify(image.Name), old, new)
							}
						}
					}
				}
			})
		}

		writeTabSection(w, func(w io.Writer) {
			count := 0
			for _, k := range keys {
				if image := diff.ChangedImages[k]; image.From == nil {
					if count == 0 {
						fmt.Fprintln(w)
						fmt.Fprintf(w, "Images Added:\n")
					}
					count++
					if showCommit {
						fmt.Fprintf(w, "  %s\t%s\n", justify(image.Name), repoAndCommit(image.To))
					} else {
						fmt.Fprintf(w, "  %s\t%s\n", justify(image.Name), digestOrRef(image.To.From.Name))
					}
				}
			}
		})

		writeTabSection(w, func(w io.Writer) {
			count := 0
			for _, k := range keys {
				if image := diff.ChangedImages[k]; image.To == nil {
					if count == 0 {
						fmt.Fprintln(w)
						fmt.Fprintf(w, "Images Removed:\n")
					}
					count++
					fmt.Fprintf(w, "  %s\n", justify(image.Name))
				}
			}
		})
	}

	printFeatureSetSection(diff, w)

	fmt.Fprintln(w)
	return nil
}

func printFeatureSetSection(diff *ReleaseDiff, w io.Writer) {
	featureSetDiff, err := calculateFeatureSetDiff(diff)
	if err != nil {
		fmt.Fprintf(w, "error reading .From release feature info: %v", err)
		return
	}

	if false {
		// TODO if we have a use for this on the CLI, then we can try to produce an ascii table
		// for now, this does nothing and we'll only render it for mardown.
		content, err := produceDiffMarkdown(nil, featureSetDiff)
		if err != nil {
			fmt.Fprintf(w, "error producing diff: %v", err)
			return
		}
		w.Write(content)
		fmt.Fprintln(w)
		fmt.Fprintln(w)
	}
}

func calculateFeatureSetDiff(diff *ReleaseDiff) (*features.ReleaseFeatureDiffInfo, error) {
	toFiles := features.FilenameToContent{}
	for filename, featureGateBytes := range diff.To.ManifestFiles {
		if !strings.Contains(filename, "featureGate-") {
			continue
		}
		toFiles[filename] = featureGateBytes
	}
	fromFiles := features.FilenameToContent{}
	for filename, featureGateBytes := range diff.From.ManifestFiles {
		if !strings.Contains(filename, "featureGate-") {
			continue
		}
		fromFiles[filename] = featureGateBytes
	}

	toReleaseFeatureInfo, err := features.ReadReleaseFeatureInfo(context.TODO(), toFiles)
	if err != nil {
		return nil, fmt.Errorf("error reading .To release feature info: %v", err)
	}
	fromReleaseFeatureInfo, err := features.ReadReleaseFeatureInfo(context.TODO(), fromFiles)
	if err != nil {
		return nil, fmt.Errorf("error reading .From release feature info: %v", err)
	}

	return toReleaseFeatureInfo.CalculateDiff(context.TODO(), fromReleaseFeatureInfo), nil
}

func markdownFeatureSetSection(info *ReleaseInfo, diff *ReleaseDiff, w io.Writer) {
	testReportBytes := info.ManifestFiles["0000_50_tests_test-reporting.yaml"]
	testReport := &configv1.TestReporting{}
	if len(testReportBytes) > 0 {
		if err := yaml.Unmarshal(testReportBytes, testReport); err != nil {
			fmt.Fprintf(w, "error reading test-reporting.yaml release feature info: %v", err)
		}
	}

	featureSetDiff, err := calculateFeatureSetDiff(diff)
	if err != nil {
		fmt.Fprintf(w, "error reading .From release feature info: %v", err)
		return
	}

	if len(featureSetDiff.GetOrderedFeatureGates()) == 0 {
		return
	}

	fmt.Fprintln(w)
	fmt.Fprintln(w, "### FeatureGate Changes")

	content, err := produceDiffMarkdown(testReport, featureSetDiff)
	if err != nil {
		fmt.Fprintf(w, "error producing diff: %v", err)
		return
	}
	w.Write(content)
	fmt.Fprintln(w)
	fmt.Fprintln(w)
}

// FeatureGatesForFeatureSetChanged returns newly enabled featureGates, newly disabled featureGates, and newly removed featureGates
func FeatureGatesForFeatureSetChanged(diff *ReleaseDiff, featureSetName string) (*configv1.FeatureGate, sets.Set[configv1.FeatureGateName], sets.Set[configv1.FeatureGateName], sets.Set[string], bool) {
	manifestName := fmt.Sprintf("0000_50_cluster-config-api_featureGate-%v.yaml", featureSetName)
	changed := diff.ChangedManifests[manifestName]
	if changed == nil || len(changed.To) == 0 {
		klog.V(2).Infof("featuregate manifest is missing in To")
		return nil, nil, nil, nil, false
	}
	toFeatureGates, err := ReadFeatureGateV1(changed.To)
	if err != nil {
		klog.V(2).Infof("failed to decode To %v: %v", manifestName, err)
		return nil, nil, nil, nil, false
	}

	toEntireSet := sets.Set[string]{}
	toEnabledSet := sets.Set[configv1.FeatureGateName]{}
	toDisabledSet := sets.Set[configv1.FeatureGateName]{}
	for _, curr := range toFeatureGates.Status.FeatureGates[0].Enabled {
		toEntireSet.Insert(string(curr.Name))
		toEnabledSet.Insert(curr.Name)
	}
	for _, curr := range toFeatureGates.Status.FeatureGates[0].Disabled {
		toEntireSet.Insert(string(curr.Name))
		toDisabledSet.Insert(curr.Name)
	}

	fromEntireSet := sets.Set[string]{}
	fromEnabledSet := sets.Set[configv1.FeatureGateName]{}
	fromDisabledSet := sets.Set[configv1.FeatureGateName]{}
	if len(changed.From) > 0 {
		fromFeatureGates, err := ReadFeatureGateV1(changed.From)
		if err != nil {
			klog.V(2).Infof("failed to decode From %v: %v", manifestName, err)
			return nil, nil, nil, nil, false
		}
		for _, curr := range fromFeatureGates.Status.FeatureGates[0].Enabled {
			fromEntireSet.Insert(string(curr.Name))
			fromEnabledSet.Insert(curr.Name)
		}
		for _, curr := range fromFeatureGates.Status.FeatureGates[0].Disabled {
			fromEntireSet.Insert(string(curr.Name))
			fromDisabledSet.Insert(curr.Name)
		}
	}

	newlyEnabledFeatureGates := toEnabledSet.Difference(fromEnabledSet)
	newlyDisabledFeatureGates := toDisabledSet.Difference(fromDisabledSet)
	removedFeatureGates := fromEntireSet.Difference(toEntireSet)
	haveFeaturesChange := len(newlyEnabledFeatureGates) > 0 || len(newlyDisabledFeatureGates) > 0 || len(removedFeatureGates) > 0
	return toFeatureGates, newlyEnabledFeatureGates, newlyDisabledFeatureGates, removedFeatureGates, haveFeaturesChange
}

func repoAndCommit(ref *imageapi.TagReference) string {
	repo := ref.Annotations[annotationBuildSourceLocation]
	commit := ref.Annotations[annotationBuildSourceCommit]
	if len(repo) == 0 || len(commit) == 0 {
		return "<unknown>"
	}
	return urlForRepoAndCommit(repo, commit)
}

func gitDiffOrCommit(from, to *imageapi.TagReference) string {
	oldRepo, newRepo := from.Annotations[annotationBuildSourceLocation], to.Annotations[annotationBuildSourceLocation]
	oldCommit, newCommit := from.Annotations[annotationBuildSourceCommit], to.Annotations[annotationBuildSourceCommit]
	if len(newRepo) == 0 || len(newCommit) == 0 {
		return "<unknown>"
	}
	if oldRepo == newRepo {
		if oldCommit == newCommit {
			return urlForRepoAndCommit(newRepo, newCommit)
		}
		return urlForRepoAndCommitRange(newRepo, oldCommit, newCommit)
	}
	if len(oldCommit) == 0 {
		return fmt.Sprintf("%s <unknown> -> %s", oldRepo, urlForRepoAndCommit(newRepo, newCommit))
	}
	if oldCommit == newCommit {
		return fmt.Sprintf("%s -> %s", oldRepo, urlForRepoAndCommit(newRepo, newCommit))
	}
	return fmt.Sprintf("%s -> %s", urlForRepoAndCommit(oldRepo, oldCommit), urlForRepoAndCommit(newRepo, newCommit))
}

func urlForRepoAndCommit(repo, commit string) string {
	if strings.HasPrefix(repo, urlGithubPrefix) {
		if u, err := url.Parse(repo); err == nil {
			u.Path = path.Join(u.Path, "commit", fmt.Sprintf("%s", commit))
			return u.String()
		}
	}
	return fmt.Sprintf("%s %s", repo, commit)
}

func urlForRepoAndCommitRange(repo, from, to string) string {
	if strings.HasPrefix(repo, urlGithubPrefix) {
		if u, err := url.Parse(repo); err == nil {
			u.Path = path.Join(u.Path, "compare", fmt.Sprintf("%s...%s", from, to))
			return u.String()
		}
	}
	return fmt.Sprintf("%s %s %s", repo, from, to)
}

func codeChanged(from, to *imageapi.TagReference) bool {
	oldCommit, newCommit := from.Annotations[annotationBuildSourceCommit], to.Annotations[annotationBuildSourceCommit]
	return len(oldCommit) > 0 && len(newCommit) > 0 && oldCommit != newCommit
}

func describeReleaseInfo(out io.Writer, release *ReleaseInfo, showCommit, showCommitURL, pullSpec, showSize bool) error {
	w := tabwriter.NewWriter(out, 0, 4, 1, ' ', 0)
	defer w.Flush()
	now := time.Now()
	fmt.Fprintf(w, "Name:\t%s\n", release.PreferredName())
	if len(release.ManifestListDigest) == 0 {
		fmt.Fprintf(w, "Digest:\t%s\n", release.Digest)
	} else {
		// if the image is manifestlisted, use the manifestlist digest.
		fmt.Fprintf(w, "Digest:\t%s\n", release.ManifestListDigest.String())
	}
	fmt.Fprintf(w, "Created:\t%s\n", release.Config.Created.UTC().Truncate(time.Second).Format(time.RFC3339))
	fmt.Fprintf(w, "OS/Arch:\t%s\n", release.Platform())
	fmt.Fprintf(w, "Manifests:\t%d\n", len(release.ManifestFiles))
	if len(release.UnknownFiles) > 0 {
		fmt.Fprintf(w, "Metadata files:\t%d\n", len(release.UnknownFiles))
	}

	fmt.Fprintln(w)
	refExact := release.ImageRef
	refExact.Ref.Tag = ""
	if len(release.ManifestListDigest) == 0 {
		refExact.Ref.ID = release.Digest.String()
	} else {
		// if the image is manifestlisted, use the manifestlist digest.
		refExact.Ref.ID = release.ManifestListDigest.String()
	}
	fmt.Fprintf(w, "Pull From:\t%s\n", refExact)

	if m := release.Metadata; m != nil {
		fmt.Fprintln(w)
		fmt.Fprintf(w, "Release Metadata:\n")
		fmt.Fprintf(w, "  Version:\t%s\n", m.Version)
		if len(m.Previous) > 0 {
			fmt.Fprintf(w, "  Upgrades:\t%s\n", strings.Join(sortSemanticVersions(m.Previous), ", "))
		} else {
			fmt.Fprintf(w, "  Upgrades:\t<none>\n")
		}
		var keys []string
		for k, v := range m.Metadata {
			switch t := v.(type) {
			case string:
				if len(t) == 0 {
					continue
				}
			case []interface{}:
				if len(t) == 0 {
					continue
				}
			case map[string]interface{}:
				if len(t) == 0 {
					continue
				}
			}
			keys = append(keys, k)
		}
		sort.Strings(keys)
		writeTabSection(w, func(w io.Writer) {
			fmt.Fprintf(w, "  Metadata:\n")
			for _, k := range keys {
				fmt.Fprintf(w, "    %s:\t%v\n", k, m.Metadata[k])
			}
		})
	}
	if len(release.ComponentVersions) > 0 {
		fmt.Fprintln(w)
		fmt.Fprintf(w, "Component Versions:\n")
		keys := release.ComponentVersions.OrderedKeys()
		for _, key := range keys {
			fmt.Fprintf(w, "  %s\t%s\t%s\n", key, release.ComponentVersions[key].Version, release.ComponentVersions[key].DisplayName)
		}
	}
	writeTabSection(w, func(w io.Writer) {
		fmt.Fprintln(w)
		fmt.Fprintf(w, "Images:\n")
		switch {
		case showSize:
			layerCount := make(map[string]int)
			baseLayer := make(map[string]int)
			totalSize := int64(0)
			for _, image := range release.Images {
				for i, layer := range image.Layers {
					digest := layer.Digest.String()
					if i == 0 {
						baseLayer[digest] = 0
					}
					count := layerCount[digest]
					if count == 0 {
						totalSize += layer.Size
					}
					layerCount[digest] = count + 1
				}
			}

			var baseHeader string
			if len(baseLayer) > 1 {
				baseHeader = "BASE"
			}
			fmt.Fprintf(w, "  NAME\t AGE\t LAYERS\t SIZE MB\t UNIQUE MB\t %s\n", baseHeader)
			coveredLayer := make(map[string]struct{})
			currentBase := 1
			for _, tag := range release.References.Spec.Tags {
				if tag.From == nil || tag.From.Kind != "DockerImage" {
					continue
				}

				image, ok := release.Images[tag.Name]
				if !ok {
					fmt.Fprintf(w, "  %s\t\t\t\t\t\n", tag.Name)
					continue
				}

				// create a column for a small number of unique base layers that visually indicates
				// which base layer belongs to which image
				var base string
				if len(baseLayer) > 1 {
					if baseIndex, ok := baseLayer[image.Layers[0].Digest.String()]; ok {
						if baseIndex == 0 {
							baseLayer[image.Layers[0].Digest.String()] = currentBase
							baseIndex = currentBase
							currentBase++
						}
						if len(baseLayer) <= 5 {
							base = strings.Repeat(" ", baseIndex-1) + string(rune('A'+baseIndex-1))
						} else {
							base = strconv.Itoa(baseIndex)
						}
					}
				}

				// count the size of the image and the unique size of the image, to give a better
				// idea of which images impact the payload the most
				unshared := int64(0)
				size := int64(0)
				for _, layer := range image.Layers {
					size += layer.Size
					if layerCount[layer.Digest.String()] > 1 {
						continue
					}
					unshared += layer.Size
				}
				// if this image has no unique layers, find the top-most layer and if this is the
				// first time it has been shown print the top layer size (as a reasonable proxy
				// for how much this image in particular contributes)
				if unshared == 0 {
					top := image.Layers[len(image.Layers)-1]
					if _, ok := coveredLayer[top.Digest.String()]; !ok {
						unshared = top.Size
						coveredLayer[top.Digest.String()] = struct{}{}
					}
				}
				age := ""
				if image.Config != nil && !image.Config.Created.IsZero() {
					age = duration.ShortHumanDuration(now.Sub(image.Config.Created))
				}
				fmt.Fprintf(w, "  %s\t%4s\t%7d\t%8.1f\t%10.1f\t %s\n", tag.Name, age, len(image.Layers), float64(size)/1024/1024, float64(unshared)/1024/1024, base)
			}
			fmt.Fprintln(w)
			if len(baseLayer) > 1 {
				fmt.Fprintf(w, "  %s across %d layers, %d different base images\n", units.HumanSize(float64(totalSize)), len(layerCount), len(baseLayer))
			} else {
				fmt.Fprintf(w, "  %s across %d layers\n", units.HumanSize(float64(totalSize)), len(layerCount))
			}

		case showCommitURL:
			fmt.Fprintf(w, "  NAME\tURL\t\n")
			for _, tag := range release.References.Spec.Tags {
				if tag.From == nil || tag.From.Kind != "DockerImage" {
					continue
				}
				var url string
				u, err := sourceLocationAsURL(tag.Annotations[annotationBuildSourceLocation])
				if err == nil && u.Host == "github.com" {
					url = fmt.Sprintf("https://github.com%s/commit/%s", u.Path, tag.Annotations[annotationBuildSourceCommit])
				} else {
					url = fmt.Sprintf("%s %s", tag.Annotations[annotationBuildSourceLocation], tag.Annotations[annotationBuildSourceCommit])
				}
				fmt.Fprintf(w, "  %s\t%s\n", tag.Name, url)
			}

		case showCommit:
			fmt.Fprintf(w, "  NAME\tREPO\tCOMMIT\t\n")
			for _, tag := range release.References.Spec.Tags {
				if tag.From == nil || tag.From.Kind != "DockerImage" {
					continue
				}
				fmt.Fprintf(w, "  %s\t%s\t%s\n", tag.Name, tag.Annotations[annotationBuildSourceLocation], tag.Annotations[annotationBuildSourceCommit])
			}

		case pullSpec:
			fmt.Fprintf(w, "  NAME\tPULL SPEC\n")
			for _, tag := range release.References.Spec.Tags {
				if tag.From == nil || tag.From.Kind != "DockerImage" {
					continue
				}
				fmt.Fprintf(w, "  %s\t%s\n", tag.Name, tag.From.Name)
			}

		default:
			fmt.Fprintf(w, "  NAME\tDIGEST\n")
			for _, tag := range release.References.Spec.Tags {
				if tag.From == nil || tag.From.Kind != "DockerImage" {
					continue
				}
				var id string
				if ref, err := imagereference.Parse(tag.From.Name); err == nil {
					id = ref.ID
				}
				if len(id) == 0 {
					id = tag.From.Name
				}
				fmt.Fprintf(w, "  %s\t%s\n", tag.Name, id)
			}
		}
	})
	if len(release.Warnings) > 0 {
		writeTabSection(w, func(w io.Writer) {
			fmt.Fprintln(w)
			fmt.Fprintf(w, "Warnings:\n")
			for _, warning := range release.Warnings {
				fmt.Fprintf(w, "* %s\n", warning)
			}
		})
	}
	fmt.Fprintln(w)
	return nil
}

func writeTabSection(out io.Writer, fn func(w io.Writer)) {
	w := tabwriter.NewWriter(out, 0, 4, 1, ' ', 0)
	fn(w)
	w.Flush()
}

func sortSemanticVersions(versionStrings []string) []string {
	var versions []semver.Version
	for _, version := range versionStrings {
		v, err := semver.Parse(version)
		if err != nil {
			return versionStrings
		}
		versions = append(versions, v)
	}
	semver.Sort(versions)
	versionStrings = make([]string, 0, len(versions))
	for _, v := range versions {
		versionStrings = append(versionStrings, v.String())
	}
	return versionStrings
}

func digestOrRef(ref string) string {
	if ref, err := imagereference.Parse(ref); err == nil && len(ref.ID) > 0 {
		return ref.ID
	}
	return ref
}

// replaceUnsafeInput prevents HTML blocks from being started in markdown for external
// inputs, but allows entities and quotes
var replaceUnsafeInput = strings.NewReplacer(
	`<`, "&lt;",
	`>`, "&gt;",
)

func describeChangelog(out, errOut io.Writer, releaseInfo *ReleaseInfo, diff *ReleaseDiff, dir, format string) error {
	if diff.To.Digest == diff.From.Digest {
		return fmt.Errorf("releases are identical")
	}

	var hasError bool

	var added, removed []string
	for k, imageDiff := range diff.ChangedImages {
		switch {
		case imageDiff.From == nil:
			added = append(added, k)
		case imageDiff.To == nil:
			removed = append(removed, k)
		}
	}
	codeChanges, imageChanges, incorrectImageChanges := releaseDiffContentChanges(diff)

	sort.Strings(added)
	sort.Strings(removed)

	switch format {
	case "json":
		changeLog := &ChangeLog{
			From: ChangeLogReleaseInfo{
				Name:    diff.From.PreferredName(),
				Created: diff.From.References.CreationTimestamp.UTC(),
				Digest:  diff.From.Digest,
			},
			To: ChangeLogReleaseInfo{
				Name:    diff.To.PreferredName(),
				Created: diff.To.References.CreationTimestamp.UTC(),
				Digest:  diff.To.Digest,
			},
			Components:    []ChangeLogComponentInfo{},
			NewImages:     []ChangeLogImageInfo{},
			RemovedImages: []ChangeLogImageInfo{},
			UpdatedImages: []ChangeLogImageInfo{},
		}
		if release, ok := diff.From.References.Annotations[annotationReleaseFromRelease]; ok {
			changeLog.From.PromotedFrom = release
		}
		if release, ok := diff.To.References.Annotations[annotationReleaseFromRelease]; ok {
			changeLog.To.PromotedFrom = release
		}
		if keys := diff.To.ComponentVersions.OrderedKeys(); len(keys) > 0 {
			for _, key := range keys {
				version := diff.To.ComponentVersions[key]
				old, ok := diff.From.ComponentVersions[key]
				if !ok || old == version {
					changeLog.Components = append(changeLog.Components, ChangeLogComponentInfo{
						Name:    componentDisplayName(key, version.DisplayName),
						Version: version.Version,
					})
					continue
				}
				changeLog.Components = append(changeLog.Components, ChangeLogComponentInfo{
					Name:    componentDisplayName(key, version.DisplayName),
					Version: version.Version,
					From:    old.Version,
				})
			}
		}
		if len(added) > 0 {
			for _, k := range added {
				changeLog.NewImages = append(changeLog.NewImages, refToChangeInfo(diff.ChangedImages[k].To))
			}
		}
		if len(removed) > 0 {
			for _, k := range removed {
				changeLog.RemovedImages = append(changeLog.RemovedImages, ChangeLogImageInfo{Name: k})
			}
		}
		if len(imageChanges) > 0 || len(incorrectImageChanges) > 0 {
			for _, change := range imageChanges {
				changeLog.RebuiltImages = append(changeLog.RebuiltImages, refToChangeInfo(diff.ChangedImages[change.Name].To))
			}
			for _, k := range incorrectImageChanges {
				changeLog.RebuiltImages = append(changeLog.RebuiltImages, ChangeLogImageInfo{Name: k})
			}
		}
		for _, change := range codeChanges {
			u, commits, err := commitsForRepo(dir, change, out, errOut)
			if err != nil {
				fmt.Fprintf(errOut, "error: %v\n", err)
				hasError = true
				continue
			}
			if len(commits) > 0 {
				info := ChangeLogImageInfo{
					Name:    strings.Join(change.ImagesAffected, ", "),
					Commits: []CommitInfo{},
				}
				if u.Host == "github.com" {
					info.Path = fmt.Sprintf("https://github.com%s/tree/%s", u.Path, change.To)
				}
				for _, commit := range commits {
					commitInfo := CommitInfo{
						Bugs:       commit.Refs.GetRefsForSource(Bugzilla),
						Issues:     commit.Refs.GetRefsForSource(Jira),
						Subject:    replaceUnsafeInput.Replace(commit.Subject),
						CommitDate: commit.CommitDate,
					}
					switch {
					case commit.PullRequest > 0:
						commitInfo.PullID = commit.PullRequest
						commitInfo.PullURL = fmt.Sprintf("https://%s%s/pull/%d", u.Host, u.Path, commit.PullRequest)
					case u.Host == "github.com":
						commitInfo.CommitID = commit.Commit[:8]
						commitInfo.CommitURL = fmt.Sprintf("https://%s%s/commit/%s", u.Host, u.Path, commit.Commit[:8])
					default:
						commitInfo.CommitID = commit.Commit[:8]
					}
					info.Commits = append(info.Commits, commitInfo)
				}
				if u.Host == "github.com" {
					info.FullChangeLog = fmt.Sprintf("https://%s%s/compare/%s...%s", u.Host, u.Path, change.From, change.To)
				}
				changeLog.UpdatedImages = append(changeLog.UpdatedImages, info)
			}
		}
		data, err := json.MarshalIndent(changeLog, "", "  ")
		if err != nil {
			return err
		}
		fmt.Fprintln(out, string(data))

	default:
		fmt.Fprintf(out, heredoc.Docf(`
		# %s

		Created: %s

		Image Digest: %s

	`, diff.To.PreferredName(), diff.To.References.CreationTimestamp.UTC(), "`"+diff.To.Digest+"`"))

		if release, ok := diff.To.References.Annotations[annotationReleaseFromRelease]; ok {
			fmt.Fprintf(out, "Promoted from %s\n\n", release)
		}
		fmt.Fprintln(out)
		fmt.Fprintf(out, "## Changes from %s\n\n", diff.From.PreferredName())

		if keys := diff.To.ComponentVersions.OrderedKeys(); len(keys) > 0 {
			fmt.Fprintf(out, "### Components\n\n")
			for _, key := range keys {
				version := diff.To.ComponentVersions[key]
				old, ok := diff.From.ComponentVersions[key]
				if !ok || old == version {
					fmt.Fprintf(out, "* %s %s\n", componentDisplayName(key, version.DisplayName), version)
					continue
				}
				fmt.Fprintf(out, "* %s upgraded from %s to %s\n", componentDisplayName(key, version.DisplayName), old, version)
			}
			fmt.Fprintln(out)
			fmt.Fprintln(out)
		}

		markdownFeatureSetSection(releaseInfo, diff, out)

		if len(added) > 0 {
			fmt.Fprintf(out, "### New images\n\n")
			for _, k := range added {
				fmt.Fprintf(out, "* %s\n", refToShortDescription(diff.ChangedImages[k].To))
			}
			fmt.Fprintln(out)
			fmt.Fprintln(out)
		}

		if len(removed) > 0 {
			fmt.Fprintf(out, "### Removed images\n\n")
			for _, k := range removed {
				fmt.Fprintf(out, "* %s\n", k)
			}
			fmt.Fprintln(out)
			fmt.Fprintln(out)
		}

		if len(imageChanges) > 0 || len(incorrectImageChanges) > 0 {
			fmt.Fprintf(out, "### Rebuilt images without code change\n\n")
			for _, change := range imageChanges {
				fmt.Fprintf(out, "* %s\n", refToShortDescription(diff.ChangedImages[change.Name].To))
			}
			for _, k := range incorrectImageChanges {
				fmt.Fprintf(out, "* %s\n", k)
			}
			fmt.Fprintln(out)
			fmt.Fprintln(out)
		}

		for _, change := range codeChanges {
			u, commits, err := commitsForRepo(dir, change, out, errOut)
			if err != nil {
				fmt.Fprintf(errOut, "error: %v\n", err)
				hasError = true
				continue
			}
			if len(commits) > 0 {
				if u.Host == "github.com" {
					fmt.Fprintf(out, "### [%s](https://github.com%s/tree/%s)\n\n", strings.Join(change.ImagesAffected, ", "), u.Path, change.To)
				} else {
					fmt.Fprintf(out, "### %s\n\n", strings.Join(change.ImagesAffected, ", "))
				}
				for _, commit := range commits {
					fmt.Fprintf(out, "*")
					commit.Refs.PrintRefs(out)
					fmt.Fprintf(out, " %s", replaceUnsafeInput.Replace(commit.Subject))
					switch {
					case commit.PullRequest > 0:
						fmt.Fprintf(out, " [#%d](%s)", commit.PullRequest, fmt.Sprintf("https://%s%s/pull/%d", u.Host, u.Path, commit.PullRequest))
					case u.Host == "github.com":
						commit := commit.Commit[:8]
						fmt.Fprintf(out, " [%s](%s)", commit, fmt.Sprintf("https://%s%s/commit/%s", u.Host, u.Path, commit))
					default:
						fmt.Fprintf(out, " %s", commit.Commit[:8])
					}
					fmt.Fprintf(out, "\n")
				}
				if u.Host == "github.com" {
					fmt.Fprintf(out, "* [Full changelog](%s)\n\n", fmt.Sprintf("https://%s%s/compare/%s...%s", u.Host, u.Path, change.From, change.To))
				} else {
					fmt.Fprintf(out, "* %s from %s to %s\n\n", change.Repo, change.FromShort(), change.ToShort())
				}
				fmt.Fprintln(out)
			}
		}
	}
	if hasError {
		return kcmdutil.ErrExit
	}
	return nil
}

func describeBugs(out, errOut io.Writer, diff *ReleaseDiff, dir string, format string, skipBugCheck bool) error {
	if diff.To.Digest == diff.From.Digest {
		return fmt.Errorf("releases are identical")
	}

	var hasError bool
	codeChanges, _, _ := releaseDiffContentChanges(diff)

	bugIDs := make(map[string]Ref)
	for _, change := range codeChanges {
		_, commits, err := commitsForRepo(dir, change, out, errOut)
		if err != nil {
			fmt.Fprintf(errOut, "error: %v\n", err)
			hasError = true
			continue
		}
		for _, commit := range commits {
			if len(commit.Refs.Refs) == 0 {
				continue
			}
			for _, v := range commit.Refs.Refs {
				// the describeBugs function only returns bug references because it's used from the --bugs argument
				// which prints bugs (not all changes) in the payload, so ignore jira refs to non bugs
				if v.Source == Jira && !(strings.HasPrefix(v.ID, "OCPBUGS-") || strings.HasPrefix(v.ID, "OCPBUGSM-")) {
					continue
				}
				if _, ok := bugIDs[v.ID]; !ok {
					bugIDs[v.ID] = v
				}
			}
		}
	}

	bugs := make(map[string]RefRemoteInfo)
	var valid []Ref
	if skipBugCheck {
		valid = GetRefList(bugIDs)
	} else {
		allBugIDs := GetRefList(bugIDs)
		for len(allBugIDs) > 0 {
			var next []Ref
			if len(allBugIDs) > 10 {
				next = allBugIDs[:10]
				allBugIDs = allBugIDs[10:]
			} else {
				next = allBugIDs
				allBugIDs = nil
			}

			bugList, err := RetrieveRefs(next)
			if err != nil {

			}
			for _, bug := range bugList.Refs {
				bugs[bug.ID] = bug
			}
		}

		for _, id := range GetRefList(bugIDs) {
			if _, ok := bugs[id.ID]; !ok {
				fmt.Fprintf(errOut, "error: Bug %s was not retrieved\n", id.ID)
				hasError = true
				continue
			}
			valid = append(valid, id)
		}
	}

	if len(valid) > 0 {
		switch format {
		case "name":
			for _, b := range valid {
				fmt.Fprintln(out, b.ID)
			}
		case "json":
			var printedBugs []RefRemoteInfo
			for _, v := range valid {
				if skipBugCheck {
					printedBugs = append(printedBugs, RefRemoteInfo{ID: v.ID, Source: v.Source})
				} else {
					if bug, ok := bugs[v.ID]; ok {
						printedBugs = append(printedBugs, bug)
					}
				}
			}
			data, err := json.MarshalIndent(printedBugs, "", "  ")
			if err != nil {
				return err
			}
			fmt.Fprintln(out, string(data))
		default:
			tw := tabwriter.NewWriter(out, 0, 0, 1, ' ', 0)
			fmt.Fprintln(tw, "ID\tSTATUS\tPRIORITY\tSUMMARY")
			for _, v := range valid {
				if bug, ok := bugs[v.ID]; ok {
					fmt.Fprintf(tw, "%s\t%s\t%s\t%s\n", v.ID, bug.Status, bug.Priority, bug.Summary)
				}
			}
			tw.Flush()
		}
	}

	if hasError {
		return kcmdutil.ErrExit
	}
	return nil
}

type ImageChange struct {
	Name     string
	From, To imagereference.DockerImageReference
}

type CodeChange struct {
	Repo     string
	From, To string

	AlternateRepos []string

	ImagesAffected []string
}

func (c CodeChange) FromShort() string {
	if len(c.From) > 8 {
		return c.From[:8]
	}
	return c.From
}

func (c CodeChange) ToShort() string {
	if len(c.To) > 8 {
		return c.To[:8]
	}
	return c.To
}

func commitsForRepo(dir string, change CodeChange, out, errOut io.Writer) (*url.URL, []MergeCommit, error) {
	u, err := sourceLocationAsURL(change.Repo)
	if err != nil {
		return nil, nil, fmt.Errorf("The source repository cannot be parsed %s: %v", change.Repo, err)
	}
	g, err := ensureCloneForRepo(dir, change.Repo, change.AlternateRepos, errOut, errOut)
	if err != nil {
		return nil, nil, err
	}
	commits, err := mergeLogForRepo(g, change.Repo, change.From, change.To)
	if err != nil {
		return nil, nil, fmt.Errorf("Could not load commits for %s: %v", change.Repo, err)
	}
	return u, commits, nil
}

func releaseDiffContentChanges(diff *ReleaseDiff) ([]CodeChange, []ImageChange, []string) {
	var imageChanges []ImageChange
	var unexpectedChanges []string
	var keys []string
	repoToCommitsToImages := make(map[string]map[string][]string)
	for k, imageDiff := range diff.ChangedImages {
		from, to := imageDiff.From, imageDiff.To
		switch {
		case from == nil, to == nil:
		default:
			newRepo := to.Annotations[annotationBuildSourceLocation]
			oldCommit, newCommit := from.Annotations[annotationBuildSourceCommit], to.Annotations[annotationBuildSourceCommit]
			if len(oldCommit) == 0 || oldCommit == newCommit {
				if from.From != nil && to.From != nil {
					if fromRef, err := imagereference.Parse(from.From.Name); err == nil {
						if toRef, err := imagereference.Parse(to.From.Name); err == nil {
							if len(fromRef.ID) > 0 && fromRef.ID == toRef.ID {
								// no change or only location changed
								break
							}
							imageChanges = append(imageChanges, ImageChange{
								Name: imageDiff.Name,
								From: fromRef,
								To:   toRef,
							})
							break
						}
					}
				}
				// before or after tag did not have a valid image reference
				unexpectedChanges = append(unexpectedChanges, k)
				break
			}
			commitRange, ok := repoToCommitsToImages[newRepo]
			if !ok {
				commitRange = make(map[string][]string)
				repoToCommitsToImages[newRepo] = commitRange
			}
			rangeID := fmt.Sprintf("%s..%s", oldCommit, newCommit)
			commitRange[rangeID] = append(commitRange[rangeID], k)
			keys = append(keys, k)
		}
	}
	sort.Slice(imageChanges, func(i, j int) bool {
		return imageChanges[i].Name < imageChanges[j].Name
	})
	sort.Strings(unexpectedChanges)
	sort.Strings(keys)
	var codeChanges []CodeChange
	for _, key := range keys {
		imageDiff := diff.ChangedImages[key]
		from, to := imageDiff.From, imageDiff.To
		oldRepo, newRepo := from.Annotations[annotationBuildSourceLocation], to.Annotations[annotationBuildSourceLocation]
		oldCommit, newCommit := from.Annotations[annotationBuildSourceCommit], to.Annotations[annotationBuildSourceCommit]

		if len(newRepo) == 0 || len(oldRepo) == 0 {
			continue
		}

		var alternateRepos []string
		if len(oldRepo) > 0 && oldRepo != newRepo {
			alternateRepos = append(alternateRepos, oldRepo)
		}

		// only display a given chunk of changes once
		commitRange := fmt.Sprintf("%s..%s", oldCommit, newCommit)
		allKeys := repoToCommitsToImages[newRepo][commitRange]
		if len(allKeys) == 0 {
			continue
		}
		repoToCommitsToImages[newRepo][commitRange] = nil
		sort.Strings(allKeys)

		codeChanges = append(codeChanges, CodeChange{
			Repo:           newRepo,
			From:           oldCommit,
			To:             newCommit,
			AlternateRepos: alternateRepos,
			ImagesAffected: allKeys,
		})
	}
	return codeChanges, imageChanges, unexpectedChanges
}

func refToShortDescription(ref *imageapi.TagReference) string {
	if from := ref.From; from != nil {
		name := ref.Name
		if u, err := sourceLocationAsURL(ref.Annotations[annotationBuildSourceLocation]); err == nil {
			if u.Host == "github.com" {
				if commit, ok := ref.Annotations[annotationBuildSourceCommit]; ok {
					shortCommit := commit
					if len(shortCommit) > 8 {
						shortCommit = shortCommit[:8]
					}
					name = fmt.Sprintf("[%s](https://github.com%s) git [%s](https://github.com%s/commit/%s)", name, u.Path, shortCommit, u.Path, commit)
				} else {
					name = fmt.Sprintf("[%s](https://github.com%s)", name, u.Path)
				}
			}
		}
		imageRef, err := imagereference.Parse(from.Name)
		if err == nil {
			switch {
			case len(imageRef.ID) > 0:
				return fmt.Sprintf("%s `%s`", name, imageRef.ID)
			case len(imageRef.Tag) > 0:
				return fmt.Sprintf("%s `:%s`", name, imageRef.Tag)
			default:
				return fmt.Sprintf("%s `%s`", name, imageRef.Exact())
			}
		}
		return fmt.Sprintf("%s `%s`", name, from.Name)
	}
	return ref.Name
}

func componentDisplayName(key, displayName string) string {
	if len(displayName) > 0 {
		return displayName
	}
	parts := strings.Split(key, "-")
	for i, part := range parts {
		if len(part) > 0 {
			parts[i] = strings.ToUpper(part[:1]) + part[1:]
		}
	}
	return strings.Join(parts, " ")
}

func orderedKeys(m map[string]string) []string {
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	return keys
}

type contentStream struct {
	current []byte
	parts   [][]byte
}

func (s *contentStream) Read(p []byte) (int, error) {
	remaining := len(p)
	count := 0
	for remaining > 0 {
		// find the next buffer, if we have nothing
		if len(s.current) == 0 {
			if len(s.parts) == 0 {
				return count, io.EOF
			}
			s.current = s.parts[0]
			s.parts = s.parts[1:]
		}

		have := len(s.current)

		// fill the buffer completely
		if have >= remaining {
			copy(p, s.current[:remaining])
			s.current = s.current[remaining:]
			return count + remaining, nil
		}

		// fill the buffer with whatever we have left
		copy(p, s.current[:have])
		s.current = nil
		p = p[have:]
		count += have
		remaining -= have
	}
	return count, nil
}

func newContentStreamForRelease(image *ReleaseInfo) io.Reader {
	names := make([]string, 0, len(image.ManifestFiles))
	for name := range image.ManifestFiles {
		names = append(names, name)
	}
	sort.Strings(names)

	rawNames := make([]string, 0, len(image.RawMetadata))
	for name := range image.RawMetadata {
		rawNames = append(rawNames, name)
	}
	sort.Strings(rawNames)

	data := make([][]byte, 0, (len(names)+len(rawNames))*3)

	for _, name := range rawNames {
		content := image.RawMetadata[name]
		data = append(data, []byte(fmt.Sprintf("# %s\n", name)), content)
		if len(content) > 0 && !bytes.HasSuffix(content, []byte("\n")) {
			data = append(data, []byte("\n"))
		}
	}
	for _, name := range names {
		content := image.ManifestFiles[name]
		data = append(data, []byte(fmt.Sprintf("# %s\n", name)), content)
		if len(content) > 0 && !bytes.HasSuffix(content, []byte("\n")) {
			data = append(data, []byte("\n"))
		}
	}
	return &contentStream{parts: data}
}

type ChangeLog struct {
	From ChangeLogReleaseInfo `json:"from"`
	To   ChangeLogReleaseInfo `json:"to"`

	Components    []ChangeLogComponentInfo `json:"components,omitempty"`
	NewImages     []ChangeLogImageInfo     `json:"newImages,omitempty"`
	RemovedImages []ChangeLogImageInfo     `json:"removedImages,omitempty"`
	RebuiltImages []ChangeLogImageInfo     `json:"rebuiltImages,omitempty"`
	UpdatedImages []ChangeLogImageInfo     `json:"updatedImages,omitempty"`
}

type ChangeLogReleaseInfo struct {
	Name         string        `json:"name"`
	Created      time.Time     `json:"created"`
	Digest       digest.Digest `json:"digest"`
	PromotedFrom string        `json:"promotedFrom,omitempty"`
}

type ChangeLogComponentInfo struct {
	Name    string `json:"name"`
	Version string `json:"version"`
	From    string `json:"from,omitempty"`
}

type ChangeLogImageInfo struct {
	Name          string       `json:"name"`
	Path          string       `json:"path"`
	ShortCommit   string       `json:"shortCommit,omitempty"`
	Commit        string       `json:"commit,omitempty"`
	ImageRef      string       `json:"imageRef,omitempty"`
	Commits       []CommitInfo `json:"commits,omitempty"`
	FullChangeLog string       `json:"fullChangeLog,omitempty"`
}

type CommitInfo struct {
	Bugs       map[string]string `json:"bugs,omitempty"`
	Issues     map[string]string `json:"issues,omitempty"`
	Subject    string            `json:"subject,omitempty"`
	PullID     int               `json:"pullID,omitempty"`
	PullURL    string            `json:"pullURL,omitempty"`
	CommitID   string            `json:"commitID,omitempty"`
	CommitURL  string            `json:"commitURL,omitempty"`
	CommitDate time.Time         `json:"commitDate,omitempty"`
}

func refToChangeInfo(ref *imageapi.TagReference) ChangeLogImageInfo {
	info := ChangeLogImageInfo{Name: ref.Name}

	if from := ref.From; from != nil {
		if u, err := sourceLocationAsURL(ref.Annotations[annotationBuildSourceLocation]); err == nil {
			if u.Host == "github.com" {
				if commit, ok := ref.Annotations[annotationBuildSourceCommit]; ok {
					shortCommit := commit
					if len(shortCommit) > 8 {
						shortCommit = shortCommit[:8]
					}
					info.Commit = commit
					info.ShortCommit = shortCommit
				}
				info.Path = fmt.Sprintf("https://github.com%s", u.Path)
			}
		}
		imageRef, err := imagereference.Parse(from.Name)
		if err == nil {
			switch {
			case len(imageRef.ID) > 0:
				info.ImageRef = imageRef.ID
			case len(imageRef.Tag) > 0:
				info.ImageRef = imageRef.Tag
			default:
				info.ImageRef = imageRef.Exact()
			}
			return info
		}
		info.ImageRef = from.Name
	}
	return info
}
