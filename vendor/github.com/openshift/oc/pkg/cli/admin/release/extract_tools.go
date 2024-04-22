package release

import (
	"archive/tar"
	"archive/zip"
	"bufio"
	"bytes"
	"compress/gzip"
	"context"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"hash"
	"io"
	"os"
	"path/filepath"
	"regexp"
	"runtime"
	"sort"
	"strings"
	"sync"
	"syscall"

	"golang.org/x/crypto/ssh/terminal"

	"golang.org/x/crypto/openpgp"

	"k8s.io/klog/v2"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/sets"
	"k8s.io/cli-runtime/pkg/genericiooptions"
	appsv1client "k8s.io/client-go/kubernetes/typed/apps/v1"
	"k8s.io/client-go/rest"
	"k8s.io/utils/pointer"
	"sigs.k8s.io/yaml"

	"github.com/MakeNowJust/heredoc"
	configv1 "github.com/openshift/api/config/v1"
	configv1client "github.com/openshift/client-go/config/clientset/versioned/typed/config/v1"
	imagereference "github.com/openshift/library-go/pkg/image/reference"
	"github.com/openshift/library-go/pkg/manifest"
	"github.com/openshift/oc/pkg/cli/admin/internal/codesign"
	"github.com/openshift/oc/pkg/cli/image/extract"
	"github.com/openshift/oc/pkg/cli/image/imagesource"
	"github.com/openshift/oc/pkg/version"
)

// extractTarget describes how a file in the release image can be extracted to disk.
type extractTarget struct {
	OS       string
	Arch     string
	Command  string
	Optional bool
	NewArch  bool

	InjectReleaseImage   bool
	InjectReleaseVersion bool
	SignMachOBinary      bool

	ArchiveFormat     string
	AsArchive         bool
	AsZip             bool
	Readme            string
	LinkTo            []string
	TargetCommandName string

	Mapping extract.Mapping
}

// installConfig is a stub for loading what we need from install-configs without having to vendor all of github.com/openshift/installer/pkg/types.
type installConfig struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata"`
	Platform          map[string]interface{}                   `json:"platform"`
	Capabilities      *configv1.ClusterVersionCapabilitiesSpec `json:"capabilities,omitempty"`
	FeatureSet        configv1.FeatureSet                      `json:"featureSet,omitempty"`
}

// manifestInclusionConfiguration configures manifest inclusion, so
// callers can opt in to new filtering options instead of having to
// update existing call-sites, even if they do not need a new
// filtering option.
type manifestInclusionConfiguration struct {
	// ExcludeIdentifier, if non-nil, excludes manifests that match the exclusion identifier.
	ExcludeIdentifier *string

	// RequiredFeatureSet, if non-nil, excludes manifests unless they match the desired feature set.
	RequiredFeatureSet *string

	// Profile, if non-nil, excludes manifests unless they match the cluster profile.
	Profile *string

	// Capabilities, if non-nil, excludes manifests unless they match the enabled cluster capabilities.
	Capabilities *configv1.ClusterVersionCapabilitiesStatus

	// Overrides excludes manifests for overridden resources.
	Overrides []configv1.ComponentOverride

	// Platform, if non-nil, excludes CredentialsRequests manifests unless they match the infrastructure platform.
	Platform *string
}

type includer func(m *manifest.Manifest) error

// extractTools extracts all referenced commands as archives in the target dir.
func (o *ExtractOptions) extractTools() error {
	return o.extractCommand("")
}

var (
	readmeInstallUnix = heredoc.Doc(`
	# OpenShift Install

	The OpenShift installer \u0060openshift-install\u0060 makes it easy to get a cluster
	running on the public cloud or your local infrastructure.

	To learn more about installing OpenShift, visit [docs.openshift.com](https://docs.openshift.com)
	and select the version of OpenShift you are using.

	## Installing the tools

	After extracting this archive, you can move the \u0060openshift-install\u0060 binary
	to a location on your PATH such as \u0060/usr/local/bin\u0060, or keep it in a temporary
	directory and reference it via \u0060./openshift-install\u0060.

	## License

	OpenShift is licensed under the Apache Public License 2.0. The source code for this
	program is [located on github](https://github.com/openshift/installer).
	`)

	readmeCLIUnix = heredoc.Doc(`
	# OpenShift Clients

	The OpenShift client \u0060oc\u0060 simplifies working with Kubernetes and OpenShift
	clusters, offering a number of advantages over \u0060kubectl\u0060 such as easy login,
	kube config file management, and access to developer tools. The \u0060kubectl\u0060
	binary is included alongside for when strict Kubernetes compliance is necessary.

	To learn more about OpenShift, visit [docs.openshift.com](https://docs.openshift.com)
	and select the version of OpenShift you are using.

	## Installing the tools

	After extracting this archive, move the \u0060oc\u0060 and \u0060kubectl\u0060 binaries
	to a location on your PATH such as \u0060/usr/local/bin\u0060. Then run:

	    oc login [API_URL]

	to start a session against an OpenShift cluster. After login, run \u0060oc\u0060 and
	\u0060oc help\u0060 to learn more about how to get started with OpenShift.

	## License

	OpenShift is licensed under the Apache Public License 2.0. The source code for this
	program is [located on github](https://github.com/openshift/oc).
	`)

	readmeCLIWindows = heredoc.Doc(`
	# OpenShift Clients

	The OpenShift client \u0060oc.exe\u0060 simplifies working with Kubernetes and OpenShift
	clusters, offering a number of advantages over \u0060kubectl.exe\u0060 such as easy login,
	kube config file management, and access to developer tools.

	To learn more about OpenShift, visit [docs.openshift.com](https://docs.openshift.com)
	and select the version of OpenShift you are using.

	## Installing the tools

	After extracting this archive, move the \u0060oc.exe\u0060 binary	to a location on your
	PATH. Then run:

	    oc login [API_URL]

	to start a session against an OpenShift cluster. After login, run \u0060oc.exe\u0060 and
	\u0060oc.exe help\u0060 to learn more about how to get started with OpenShift.

	If you would like to use \u0060kubectl.exe\u0060 instead, copy the \u0060oc.exe\u0060 file
	and rename it to \u0060kubectl.exe\u0060. The interface will follow the conventions of that
	CLI.

	## License

	OpenShift is licensed under the Apache Public License 2.0. The source code for this
	program is [located on github](https://github.com/openshift/oc).
	`)

	readmeCCOUnix = heredoc.Doc(`
	# Cloud Credential Operator utility

	The ccoctl tool provides various commands to assist with the creating and maintenance of
	cloud credentials from outside the cluster (necessary when CCO is put in "Manual" mode).

	To learn more about OpenShift, visit [docs.openshift.com](https://docs.openshift.com)
	and select the version of OpenShift you are using.

	## Installing the tools

	After extracting this archive, move the \u0060ccoctl\u0060 binary to a location on your
	PATH such as \u0060/usr/local/bin\u0060, or keep it in a temporary directory and
	reference it via \u0060./ccoctl\u0060.

	## License

	OpenShift is licensed under the Apache Public License 2.0. The source code for this
	program is [located on github](https://github.com/openshift/cloud-credential-operator).
	`)

	// indicates that the architecture of the binary matches the release payload
	targetReleaseArch = "release-arch"
)

// extractTools extracts specific commands out of images referenced by the release image.
// TODO: in the future the metadata this command contains might be loaded from the release
//
//	image, but we must maintain compatibility with older payloads if so
func (o *ExtractOptions) extractCommand(command string) error {
	// Available targets is treated as a GA API and may not be changed without backwards
	// compatibility of at least N-2 releases.
	availableTargets := []extractTarget{
		{
			OS:      "darwin",
			Arch:    "amd64",
			Command: "oc",
			Mapping: extract.Mapping{Image: "cli-artifacts", From: "usr/share/openshift/mac/oc"},

			LinkTo:               []string{"kubectl"},
			Readme:               readmeCLIUnix,
			InjectReleaseVersion: true,
			ArchiveFormat:        "openshift-client-mac-%s.tar.gz",
		},
		{
			OS:      "darwin",
			Arch:    "arm64",
			Command: "oc",
			NewArch: true,
			Mapping: extract.Mapping{Image: "cli-artifacts", From: "usr/share/openshift/mac_arm64/oc"},

			LinkTo:               []string{"kubectl"},
			Readme:               readmeCLIUnix,
			InjectReleaseVersion: true,
			SignMachOBinary:      true,
			ArchiveFormat:        "openshift-client-mac-arm64-%s.tar.gz",
		},
		{
			OS:      "linux",
			Arch:    targetReleaseArch,
			Command: "oc",
			Mapping: extract.Mapping{Image: "cli", From: "usr/bin/oc"},

			LinkTo:               []string{"kubectl"},
			Readme:               readmeCLIUnix,
			InjectReleaseVersion: true,
			ArchiveFormat:        "openshift-client-linux-%s.tar.gz",
		},
		{
			OS:      "linux",
			Arch:    "amd64",
			Command: "oc",
			Mapping: extract.Mapping{Image: "cli-artifacts", From: "usr/share/openshift/linux_amd64/oc"},

			LinkTo:               []string{"kubectl"},
			Readme:               readmeCLIUnix,
			InjectReleaseVersion: true,
			ArchiveFormat:        "openshift-client-linux-amd64-%s.tar.gz",
		},
		{
			OS:      "linux",
			Arch:    "amd64",
			Command: "oc.rhel9",
			Mapping: extract.Mapping{Image: "cli-artifacts", From: "usr/share/openshift/linux_amd64/oc.rhel9"},

			LinkTo:               []string{"kubectl"},
			Readme:               readmeCLIUnix,
			InjectReleaseVersion: true,
			ArchiveFormat:        "openshift-client-linux-amd64-rhel9-%s.tar.gz",
			TargetCommandName:    "oc",
		},
		{
			OS:      "linux",
			Arch:    "amd64",
			Command: "oc.rhel8",
			Mapping: extract.Mapping{Image: "cli-artifacts", From: "usr/share/openshift/linux_amd64/oc.rhel8"},

			LinkTo:               []string{"kubectl"},
			Readme:               readmeCLIUnix,
			InjectReleaseVersion: true,
			ArchiveFormat:        "openshift-client-linux-amd64-rhel8-%s.tar.gz",
			TargetCommandName:    "oc",
		},
		{
			OS:      "linux",
			Arch:    "arm64",
			Command: "oc",
			NewArch: true,
			Mapping: extract.Mapping{Image: "cli-artifacts", From: "usr/share/openshift/linux_arm64/oc"},

			LinkTo:               []string{"kubectl"},
			Readme:               readmeCLIUnix,
			InjectReleaseVersion: true,
			ArchiveFormat:        "openshift-client-linux-arm64-%s.tar.gz",
		},
		{
			OS:      "linux",
			Arch:    "arm64",
			Command: "oc.rhel9",
			NewArch: true,
			Mapping: extract.Mapping{Image: "cli-artifacts", From: "usr/share/openshift/linux_arm64/oc.rhel9"},

			LinkTo:               []string{"kubectl"},
			Readme:               readmeCLIUnix,
			InjectReleaseVersion: true,
			ArchiveFormat:        "openshift-client-linux-arm64-rhel9-%s.tar.gz",
			TargetCommandName:    "oc",
		},
		{
			OS:      "linux",
			Arch:    "arm64",
			Command: "oc.rhel8",
			NewArch: true,
			Mapping: extract.Mapping{Image: "cli-artifacts", From: "usr/share/openshift/linux_arm64/oc.rhel8"},

			LinkTo:               []string{"kubectl"},
			Readme:               readmeCLIUnix,
			InjectReleaseVersion: true,
			ArchiveFormat:        "openshift-client-linux-arm64-rhel8-%s.tar.gz",
			TargetCommandName:    "oc",
		},
		{
			OS:      "linux",
			Arch:    "ppc64le",
			Command: "oc",
			NewArch: true,
			Mapping: extract.Mapping{Image: "cli-artifacts", From: "usr/share/openshift/linux_ppc64le/oc"},

			LinkTo:               []string{"kubectl"},
			Readme:               readmeCLIUnix,
			InjectReleaseVersion: true,
			ArchiveFormat:        "openshift-client-linux-ppc64le-%s.tar.gz",
		},
		{
			OS:      "linux",
			Arch:    "ppc64le",
			Command: "oc.rhel9",
			NewArch: true,
			Mapping: extract.Mapping{Image: "cli-artifacts", From: "usr/share/openshift/linux_ppc64le/oc.rhel9"},

			LinkTo:               []string{"kubectl"},
			Readme:               readmeCLIUnix,
			InjectReleaseVersion: true,
			ArchiveFormat:        "openshift-client-linux-ppc64le-rhel9-%s.tar.gz",
			TargetCommandName:    "oc",
		},
		{
			OS:      "linux",
			Arch:    "ppc64le",
			Command: "oc.rhel8",
			NewArch: true,
			Mapping: extract.Mapping{Image: "cli-artifacts", From: "usr/share/openshift/linux_ppc64le/oc.rhel8"},

			LinkTo:               []string{"kubectl"},
			Readme:               readmeCLIUnix,
			InjectReleaseVersion: true,
			ArchiveFormat:        "openshift-client-linux-ppc64le-rhel8-%s.tar.gz",
			TargetCommandName:    "oc",
		},
		{
			OS:      "windows",
			Arch:    "amd64",
			Command: "oc",
			Mapping: extract.Mapping{Image: "cli-artifacts", From: "usr/share/openshift/windows/oc.exe"},

			Readme:               readmeCLIWindows,
			InjectReleaseVersion: true,
			ArchiveFormat:        "openshift-client-windows-%s.zip",
			AsZip:                true,
		},
		{
			OS:      "darwin",
			Arch:    "amd64",
			Command: "openshift-install",
			Mapping: extract.Mapping{Image: "installer-artifacts", From: "usr/share/openshift/mac/openshift-install"},

			Readme:               readmeInstallUnix,
			InjectReleaseImage:   true,
			InjectReleaseVersion: true,
			ArchiveFormat:        "openshift-install-mac-%s.tar.gz",
		},
		{
			OS:      "darwin",
			Arch:    "arm64",
			Command: "openshift-install",
			NewArch: true,
			Mapping: extract.Mapping{Image: "installer-artifacts", From: "usr/share/openshift/mac_arm64/openshift-install"},

			Readme:               readmeInstallUnix,
			InjectReleaseImage:   true,
			InjectReleaseVersion: true,
			SignMachOBinary:      true,
			ArchiveFormat:        "openshift-install-mac-arm64-%s.tar.gz",
		},
		{
			OS:      "linux",
			Arch:    targetReleaseArch,
			Command: "openshift-install",
			Mapping: extract.Mapping{Image: "installer", From: "usr/bin/openshift-install"},

			Readme:               readmeInstallUnix,
			InjectReleaseImage:   true,
			InjectReleaseVersion: true,
			ArchiveFormat:        "openshift-install-linux-%s.tar.gz",
		},
		{
			OS:      "linux",
			Arch:    "amd64",
			Command: "openshift-install",
			NewArch: true,
			Mapping: extract.Mapping{Image: "installer-artifacts", From: "usr/share/openshift/linux_amd64/openshift-install"},

			Readme:               readmeInstallUnix,
			InjectReleaseImage:   true,
			InjectReleaseVersion: true,
			ArchiveFormat:        "openshift-install-linux-amd64-%s.tar.gz",
		},
		{
			OS:      "linux",
			Arch:    "arm64",
			Command: "openshift-install",
			NewArch: true,
			Mapping: extract.Mapping{Image: "installer-artifacts", From: "usr/share/openshift/linux_arm64/openshift-install"},

			Readme:               readmeInstallUnix,
			InjectReleaseImage:   true,
			InjectReleaseVersion: true,
			ArchiveFormat:        "openshift-install-linux-arm64-%s.tar.gz",
		},
		{
			OS:       "linux",
			Arch:     targetReleaseArch,
			Command:  "openshift-baremetal-install",
			Optional: true,
			Mapping:  extract.Mapping{Image: "baremetal-installer", From: "usr/bin/openshift-install"},

			Readme:               readmeInstallUnix,
			InjectReleaseImage:   true,
			InjectReleaseVersion: true,
			ArchiveFormat:        "openshift-baremetal-install-linux-%s.tar.gz",
		},
		{
			OS:       "linux",
			Arch:     targetReleaseArch,
			Command:  "openshift-install-fips",
			Optional: true,
			Mapping:  extract.Mapping{Image: "baremetal-installer", From: "usr/bin/openshift-install"},

			Readme:               readmeInstallUnix,
			InjectReleaseImage:   true,
			InjectReleaseVersion: true,
			ArchiveFormat:        "openshift-install-rhel-%s.tar.gz",
		},
		{
			OS:      "linux",
			Arch:    targetReleaseArch,
			Command: "ccoctl",
			NewArch: true,
			Mapping: extract.Mapping{Image: "cloud-credential-operator", From: "usr/bin/ccoctl"},

			Readme:        readmeCCOUnix,
			ArchiveFormat: "ccoctl-linux-%s.tar.gz",
		},
	}

	currentArch := runtime.GOARCH
	currentOS := runtime.GOOS
	if len(o.CommandOperatingSystem) > 0 {
		currentOS = o.CommandOperatingSystem
		if currentOS == "*" {
			currentArch = "*"
		} else {
			parsed := strings.Split(currentOS, "/")
			switch len(parsed) {
			case 1:
				// --command-os=linux
				currentArch = "amd64"
			case 2:
				// --command-os=linux/amd64
				currentOS = parsed[0]
				currentArch = parsed[1]
			default:
				return fmt.Errorf("invalid --command-os value")
			}

		}
	}
	if currentOS == "mac" {
		currentOS = "darwin"
	}

	// Select the subset of targets based on command line input
	var willArchive bool
	var targets []extractTarget

	// Filter by command, or gather all non-optional targets
	if len(command) > 0 {
		for _, target := range availableTargets {
			if target.Command == command {
				targets = append(targets, target)
			}
		}
	} else {
		for _, target := range availableTargets {
			if !target.Optional {
				targets = append(targets, target)
			}
		}
	}

	// If the user didn't specify a command, or the operating system is set
	// to '*', we'll produce an archive
	if len(command) == 0 || o.CommandOperatingSystem == "*" {
		for i := range targets {
			targets[i].AsArchive = true
			targets[i].AsZip = targets[i].OS == "windows"
		}
	}

	if len(targets) == 0 {
		switch {
		case len(command) > 0 && currentOS != "*":
			return fmt.Errorf("command %q does not support the operating system %q", o.Command, currentOS)
		case len(command) > 0:
			return fmt.Errorf("the supported commands are 'oc' and 'openshift-install'")
		default:
			return fmt.Errorf("no available commands")
		}
	}

	var hashFn = sha256.New
	var signer *openpgp.Entity
	if willArchive && len(o.SigningKey) > 0 {
		key, err := os.ReadFile(o.SigningKey)
		if err != nil {
			return err
		}
		keyring, err := openpgp.ReadArmoredKeyRing(bytes.NewBuffer(key))
		if err != nil {
			return err
		}
		for _, key := range keyring {
			if !key.PrivateKey.CanSign() {
				continue
			}
			fmt.Fprintf(o.Out, "Enter password for private key: ")
			password, err := terminal.ReadPassword(int(syscall.Stdin))
			fmt.Fprintln(o.Out)
			if err != nil {
				return err
			}
			if err := key.PrivateKey.Decrypt(password); err != nil {
				return fmt.Errorf("unable to decrypt signing key: %v", err)
			}
			for i, subkey := range key.Subkeys {
				if err := subkey.PrivateKey.Decrypt(password); err != nil {
					return fmt.Errorf("unable to decrypt signing subkey %d: %v", i, err)
				}
			}
			signer = key
			break
		}
		if signer == nil {
			return fmt.Errorf("no private key exists in %s capable of signing the output", o.SigningKey)
		}
	}

	// load the release image
	dir := o.Directory
	infoOptions := NewInfoOptions(o.IOStreams)
	infoOptions.SecurityOptions = o.SecurityOptions
	infoOptions.FilterOptions = o.FilterOptions
	infoOptions.FileDir = o.FileDir
	infoOptions.ICSPFile = o.ICSPFile
	infoOptions.IDMSFile = o.IDMSFile
	release, err := infoOptions.LoadReleaseInfo(o.From, false)
	if err != nil {
		return err
	}
	releaseArch := release.Config.Architecture
	releaseName := release.PreferredName()
	refExact := release.ImageRef
	refExact.Ref.Tag = ""
	// if the release image is manifestlist image, we'll not change digest with
	// arch based sha. Because we want that the extracted tool can be used for all archs.
	if len(release.ManifestListDigest) == 0 {
		refExact.Ref.ID = release.Digest.String()
	} else {
		// if the image is manifestlisted, use the manifestlist digest.
		refExact.Ref.ID = release.ManifestListDigest.String()
	}
	exactReleaseImage := refExact.String()

	targetArchCommands := make(map[string]struct{})
	// resolve target image references to their pull specs
	missing := sets.NewString()
	var validTargets []extractTarget
	for _, target := range targets {
		if currentOS != "*" && target.OS != currentOS {
			klog.V(2).Infof("Skipping %s, does not match current OS %s", target.ArchiveFormat, target.OS)
			continue
		}
		if currentArch != "*" && target.Arch != currentArch {
			if currentArch != releaseArch || target.Arch != targetReleaseArch {
				klog.V(2).Infof("Skipping %s, does not match current architecture %s", target.ArchiveFormat, target.Arch)
				continue
			}
		}

		if target.Arch == targetReleaseArch {
			targetArchCommands[target.Command] = struct{}{}
		}

		if target.OS == "linux" && target.Arch == releaseArch {
			if _, ok := targetArchCommands[target.Command]; ok {
				// Some target commands have release-arch types that defines extracting
				// the command in whatever release architecture type is set(e.g. linux/amd64)
				// But there is also another target type for these commands specifically set to
				// each architecture type(linux/amd64, linux/arm64) and it is expected that
				// one of these arch types collide with release-arch type. Thus,
				// to prevent duplicate extraction, we have to skip the one colliding with release-arch type.
				// However, we need to skip per command name because some command may not have release-arch type.
				klog.V(2).Infof("Skipping duplicate %s", target.ArchiveFormat)
				continue
			}
		}
		spec, err := findImageSpec(release.References, target.Mapping.Image, o.From)
		if err != nil && !target.NewArch {
			missing.Insert(target.Mapping.Image)
			continue
		}
		klog.V(2).Infof("Will extract %s from %s", target.Mapping.From, spec)
		ref, err := imagereference.Parse(spec)
		if err != nil {
			return err
		}
		target.Mapping.Image = spec
		target.Mapping.ImageRef = imagesource.TypedImageReference{Ref: ref, Type: imagesource.DestinationRegistry}
		// if the name of the extracted binary is set to different from the
		// actual command name, we set it to new target command name.
		if target.TargetCommandName != "" {
			target.Command = target.TargetCommandName
		}
		if target.AsArchive {
			willArchive = true
			target.Mapping.Name = fmt.Sprintf(target.ArchiveFormat, releaseName)
			target.Mapping.To = filepath.Join(dir, target.Mapping.Name)
		} else {
			target.Mapping.To = filepath.Join(dir, target.Command)
			target.Mapping.Name = fmt.Sprintf("%s-%s-%s", target.OS, target.Arch, target.Command)
		}
		validTargets = append(validTargets, target)
	}

	if len(validTargets) == 0 {
		if len(missing) == 1 {
			return fmt.Errorf("the image %q containing the desired command is not available", missing.List()[0])
		}
		if len(missing) > 1 {
			return fmt.Errorf("some required images are missing: %s", strings.Join(missing.List(), ", "))
		}
		return fmt.Errorf("the command '%s' is not available for %s/%s on %s", command, currentOS, currentArch, o.From)
	}
	if len(missing) > 0 {
		fmt.Fprintf(o.ErrOut, "warning: Some commands can not be extracted due to missing images: %s\n", strings.Join(missing.List(), ", "))
	}

	// will extract in parallel
	opts := extract.NewExtractOptions(genericiooptions.IOStreams{Out: o.Out, ErrOut: o.ErrOut})
	opts.ParallelOptions = o.ParallelOptions
	opts.SecurityOptions = o.SecurityOptions
	opts.FilterOptions = o.FilterOptions
	opts.ICSPFile = o.ICSPFile
	opts.IDMSFile = o.IDMSFile
	opts.OnlyFiles = true

	// create the mapping lookup of the valid targets
	var extractLock sync.Mutex
	targetsByName := make(map[string]extractTarget)
	for _, target := range validTargets {
		targetsByName[target.Mapping.Name] = target
		opts.Mappings = append(opts.Mappings, target.Mapping)
	}
	hashByTargetName := make(map[string]string)

	// ensure to is a directory
	if err := os.MkdirAll(dir, 0777); err != nil {
		return err
	}

	// as each layer is extracted, take the output binary and write it to disk
	opts.TarEntryCallback = func(hdr *tar.Header, layer extract.LayerInfo, r io.Reader) (bool, error) {
		// ensure we don't process the same mapping twice due to programmer error
		target, ok := func() (extractTarget, bool) {
			extractLock.Lock()
			defer extractLock.Unlock()
			target, ok := targetsByName[layer.Mapping.Name]
			return target, ok
		}()
		if !ok {
			return false, fmt.Errorf("unable to find target with mapping name %s", layer.Mapping.Name)
		}

		// open the file
		f, err := os.OpenFile(layer.Mapping.To, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0755)
		if err != nil {
			return false, err
		}

		// if we need to write an archive, wrap the file appropriately to create a single
		// entry
		var w io.Writer = f

		bw := bufio.NewWriterSize(w, 16*1024)
		w = bw

		var hash hash.Hash
		closeFn := func() error { return nil }
		if target.AsArchive {
			text := strings.Replace(target.Readme, `\u0060`, "`", -1)
			hash = hashFn()
			w = io.MultiWriter(hash, w)
			if target.AsZip {
				klog.V(2).Infof("Writing %s as a ZIP archive %s", hdr.Name, layer.Mapping.To)
				zw := zip.NewWriter(w)

				if len(text) > 0 {
					text = strings.Replace(text, "\n", "\r\n", -1)
					zh := &zip.FileHeader{
						Method:             zip.Deflate,
						Name:               "README.md",
						UncompressedSize64: uint64(len(text)),
						Modified:           hdr.ModTime,
					}
					zh.SetMode(os.FileMode(0755))

					fw, err := zw.CreateHeader(zh)
					if err != nil {
						return false, err
					}
					if _, err := fmt.Fprintf(fw, text); err != nil {
						return false, err
					}
				}

				zh := &zip.FileHeader{
					Method:             zip.Deflate,
					Name:               target.Command + ".exe",
					UncompressedSize64: uint64(hdr.Size),
					Modified:           hdr.ModTime,
				}
				zh.SetMode(os.FileMode(0755))

				fw, err := zw.CreateHeader(zh)
				if err != nil {
					return false, err
				}

				w = fw
				closeFn = func() error { return zw.Close() }

			} else {
				klog.V(2).Infof("Writing %s as a tar.gz archive %s", hdr.Name, layer.Mapping.To)
				gw, err := gzip.NewWriterLevel(w, 3)
				if err != nil {
					return false, err
				}
				tw := tar.NewWriter(gw)

				if len(text) > 0 {
					if err := tw.WriteHeader(&tar.Header{
						Name:     "README.md",
						Mode:     int64(os.FileMode(0644).Perm()),
						Size:     int64(len(text)),
						Typeflag: tar.TypeReg,
						ModTime:  hdr.ModTime,
					}); err != nil {
						return false, err
					}
					if _, err := fmt.Fprintf(tw, text); err != nil {
						return false, err
					}
				}

				if err := tw.WriteHeader(&tar.Header{
					Name:     target.Command,
					Mode:     int64(os.FileMode(0755).Perm()),
					Size:     hdr.Size,
					Typeflag: tar.TypeReg,
					ModTime:  hdr.ModTime,
				}); err != nil {
					return false, err
				}

				w = tw
				closeFn = func() error {
					for _, link := range target.LinkTo {
						if err := tw.WriteHeader(&tar.Header{
							Name:     link,
							Mode:     int64(os.FileMode(0755).Perm()),
							Size:     0,
							Typeflag: tar.TypeLink,
							ModTime:  hdr.ModTime,
							Linkname: target.Command,
						}); err != nil {
							return err
						}
					}
					if err := tw.Close(); err != nil {
						return err
					}
					return gw.Close()
				}
			}
		}

		// copy the input to disk
		replacements := []replacement{}
		if target.InjectReleaseImage {
			replacements = append(replacements, replacement{
				name:   "release image",
				marker: append([]byte{0}, []byte(releaseImageMarker[1:])...),
				value:  exactReleaseImage,
			})
		}
		if target.InjectReleaseVersion {
			replacements = append(replacements, replacement{
				name:   "release version",
				marker: append([]byte{0}, []byte(releaseVersionMarker[1:])...),
				value:  releaseName,
			})
		}
		err = copyAndReplace(o.ErrOut, w, r, 4*1024, replacements, target.Command)
		if err != nil {
			closeFn()
			f.Close()
			os.Remove(f.Name())
			return false, err
		}

		// ensure the file is written to disk
		if err := closeFn(); err != nil {
			return false, err
		}
		if err := bw.Flush(); err != nil {
			return false, err
		}
		if err := f.Close(); err != nil {
			return false, err
		}
		if err := os.Chtimes(f.Name(), hdr.ModTime, hdr.ModTime); err != nil {
			klog.V(2).Infof("Unable to set extracted file modification time: %v", err)
		}

		if (target.InjectReleaseVersion || target.InjectReleaseImage) && target.SignMachOBinary {
			if err = codesign.ResignMacho(layer.Mapping.To, target.AsArchive, target.Command, target.LinkTo); err != nil {
				klog.Infof("Unable to resign macho binaries:  %v", err)
			} else {
				if target.AsArchive {
					// Since we rewrite tarball after signing mach-o files in darwin/arm64,
					// we should reflect the modified hash sum of this tarball to prevent mismatches.
					h := hashFn()
					archived, err := os.Open(layer.Mapping.To)
					if err != nil {
						return false, err
					}
					if _, err = io.Copy(h, archived); err != nil {
						archived.Close()
						return false, err
					}
					hash = h
					archived.Close()
				}
			}
		}

		func() {
			extractLock.Lock()
			defer extractLock.Unlock()
			delete(targetsByName, layer.Mapping.Name)
			if hash != nil {
				hashByTargetName[layer.Mapping.To] = hex.EncodeToString(hash.Sum(nil))
			}
		}()

		return false, nil
	}
	if err := opts.Run(); err != nil {
		return err
	}

	if willArchive {
		buf := &bytes.Buffer{}
		fmt.Fprintf(buf, heredoc.Doc(`
			Client tools for OpenShift
			--------------------------

			These archives contain the client tooling for [OpenShift](https://docs.openshift.com).

			To verify the contents of this directory, use the 'gpg' and 'shasum' tools to
			ensure the archives you have downloaded match those published from this location.

			The openshift-install binary has been preconfigured to install the following release:

			---

		`))
		if err := describeReleaseInfo(buf, release, false, false, true, false); err != nil {
			return err
		}
		filename := "release.txt"
		if err := os.WriteFile(filepath.Join(dir, filename), buf.Bytes(), 0644); err != nil {
			return err
		}
		hash := hashFn()
		hash.Write(buf.Bytes())
		hashByTargetName[filename] = hex.EncodeToString(hash.Sum(nil))
	}

	// write a checksum of the tar files to disk as sha256sum.txt.asc
	if len(hashByTargetName) > 0 {
		var keys []string
		for k := range hashByTargetName {
			keys = append(keys, k)
		}
		sort.Strings(keys)
		var lines []string
		for _, k := range keys {
			hash := hashByTargetName[k]
			lines = append(lines, fmt.Sprintf("%s  %s", hash, filepath.Base(k)))
		}
		// ensure a trailing newline
		if len(lines[len(lines)-1]) != 0 {
			lines = append(lines, "")
		}
		// write the content manifest
		data := []byte(strings.Join(lines, "\n"))
		filename := "sha256sum.txt"
		if err := os.WriteFile(filepath.Join(dir, filename), data, 0644); err != nil {
			return fmt.Errorf("unable to write checksum file: %v", err)
		}
		// sign the content manifest
		if signer != nil {
			buf := &bytes.Buffer{}
			if err := openpgp.ArmoredDetachSign(buf, signer, bytes.NewBuffer(data), nil); err != nil {
				return fmt.Errorf("unable to sign the sha256sum.txt file: %v", err)
			}
			if err := os.WriteFile(filepath.Join(dir, filename+".asc"), buf.Bytes(), 0644); err != nil {
				return fmt.Errorf("unable to write signed manifest: %v", err)
			}
		}
	}

	// if we did not process some targets, report that to the user and error if necessary
	if len(targetsByName) > 0 {
		var missing []string
		for _, target := range targetsByName {
			if target.NewArch {
				continue
			}
			if command == "" && (strings.Contains(target.Mapping.From, "rhel9") || strings.Contains(target.Mapping.From, "rhel8")) {
				// if user explicitly wants to extract oc.rhel9(or installer.rhel9) via --command=oc.rhel9 and
				// if release does not have this binary, we can safely return error.
				// On the other hand, if user wants to extract all the tooling in older versions via --tools flag,
				// we shouldn't print any error indicating that oc.rhel9 does not exist in this release payload.
				continue
			}
			missing = append(missing, target.Mapping.From)
		}
		sort.Strings(missing)
		if len(missing) == 1 {
			return fmt.Errorf("image did not contain %s", missing[0])
		}
		if len(missing) > 1 {
			return fmt.Errorf("unable to find multiple files: %s", strings.Join(missing, ", "))
		}
	}

	return nil
}

const (
	// releaseImageMarker is the placeholder within a binary for the release image pullspec.
	releaseImageMarker = "!_RELEASE_IMAGE_LOCATION_\x00XXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXX\x00"
	// releaseVersionMarker is the placeholder within a binary for the release image version name string.
	releaseVersionMarker = "!_RELEASE_VERSION_LOCATION_\x00XXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXX\x00"
)

type replacement struct {
	name   string
	marker []byte
	value  string
}

// copyAndReplace performs a targeted replacement for binaries that
// contain special marker strings, replacing the first occurrence of each
// marker with a new string and a NUL terminating byte.  It logs a warning
// if any replacements are not performed.
func copyAndReplace(errorOutput io.Writer, w io.Writer, r io.Reader, bufferSize int, replacements []replacement, name string) error {
	if len(replacements) == 0 {
		_, err := io.Copy(w, r)
		return err
	}

	longestMarker := 0
	for _, replacement := range replacements {
		if len(replacement.value) > len(replacement.marker)-1 {
			return fmt.Errorf("the %s value has %d bytes, but the maximum replacement length is %d", replacement.name, len(replacement.value), len(replacement.marker)-1)
		}

		if len(replacement.marker) > longestMarker {
			longestMarker = len(replacement.marker)
			if bufferSize < longestMarker {
				return fmt.Errorf("the buffer size must be greater than %d bytes to find %s", len(replacement.marker), replacement.name)
			}
		}
	}

	offset := 0
	buf := make([]byte, bufferSize)
	remaining := len(replacements)
	matches := make([]bool, remaining)

	for {
		n, err := io.ReadFull(r, buf[offset:])

		// search in the buffer for the expected match
		end := offset + n
		if n > 0 {
			for i, matched := range matches {
				if !matched {
					replacement := replacements[i]
					index := bytes.Index(buf[:end], replacement.marker)
					if index != -1 {
						klog.V(2).Infof("Found match for %s at %d (len=%d, offset=%d, n=%d)", replacements[i].name, index, len(buf), offset, n)
						// the replacement starts at the beginning of the match, contains the replacement value and a terminating NUL byte
						copy(buf[index:index+len(replacement.value)], []byte(replacement.value))
						buf[index+len(replacement.value)] = 0x00
						matches[i] = true
						remaining -= 1
					}
				}
			}
		}

		// write everything that we have already searched (excluding the end of the buffer that will
		// be checked next pass)
		writeTo := end - longestMarker
		if writeTo < 0 {
			writeTo = 0
		}
		if remaining == 0 || err != nil {
			writeTo = end
		}
		offset = 0
		for offset < writeTo {
			n, wErr := w.Write(buf[offset:writeTo])
			if wErr != nil {
				return wErr
			}
			offset += n
		}

		if err != nil {
			if err == io.EOF || err == io.ErrUnexpectedEOF {
				remainingNames := make([]string, 0, len(replacements))
				for i, matched := range matches {
					if !matched {
						remainingNames = append(remainingNames, replacements[i].name)
					}
				}
				sort.Strings(remainingNames)
				if len(remainingNames) > 0 {
					fmt.Fprintf(errorOutput, "warning: Unable to make all expected replacements in %s.  Remaining: %s", name, strings.Join(remainingNames, ", "))
				}
				return nil
			}
			return err
		}

		// once we match each replacement, we can copy the rest of the file without processing
		if remaining == 0 {
			_, err := io.Copy(w, r)
			return err
		}

		// ensure the beginning of the buffer matches the end of the current buffer so that we
		// can search for matches that span buffers
		copy(buf[:writeTo], buf[writeTo:end])
		offset = end - writeTo
	}

}

func findClusterIncludeConfigFromInstallConfig(ctx context.Context, installConfigPath string) (manifestInclusionConfiguration, error) {
	config := manifestInclusionConfiguration{}

	clientVersion, reportedVersion, err := version.ExtractVersion()
	if err != nil {
		return config, err
	}
	if reportedVersion == "" {
		reportedVersion = clientVersion.String()
	}

	installConfigBytes, err := os.ReadFile(installConfigPath)
	if err != nil {
		return config, err
	}

	data := installConfig{}

	if err := yaml.Unmarshal(installConfigBytes, &data); err != nil {
		return config, fmt.Errorf("failed to parse %s: %w", installConfigPath, err)
	}

	if data.APIVersion != "v1" {
		return config, fmt.Errorf("unrecognized %s API version: %q (expected %q)", installConfigPath, data.APIVersion, "v1")
	}

	config.RequiredFeatureSet = pointer.String(string(data.FeatureSet))
	config.Profile = pointer.String(manifest.DefaultClusterProfile) // assumption, but there's no install-config data about profile to give us more insight
	for key := range data.Platform {
		config.Platform = pointer.String(key)
	}

	if data.Capabilities != nil {
		config.Capabilities = &configv1.ClusterVersionCapabilitiesStatus{}
		if enabled, ok := configv1.ClusterVersionCapabilitySets[data.Capabilities.BaselineCapabilitySet]; !ok {
			return config, fmt.Errorf("unrecognized baselineCapabilitySet %q", data.Capabilities.BaselineCapabilitySet)
		} else {
			if data.Capabilities.BaselineCapabilitySet == configv1.ClusterVersionCapabilitySetCurrent {
				klog.Infof("If the eventual cluster will not be the same minor version as this %s 'oc', the actual %s capability set may differ.", reportedVersion, data.Capabilities.BaselineCapabilitySet)
			}
			config.Capabilities.EnabledCapabilities = append(config.Capabilities.EnabledCapabilities, enabled...)
		}
		config.Capabilities.EnabledCapabilities = append(config.Capabilities.EnabledCapabilities, data.Capabilities.AdditionalEnabledCapabilities...)

		klog.Infof("If the eventual cluster will not be the same minor version as this %s 'oc', the known capability sets may differ.", reportedVersion)
		config.Capabilities.KnownCapabilities = configv1.KnownClusterVersionCapabilities
	}

	return config, nil
}

func findClusterIncludeConfig(ctx context.Context, restConfig *rest.Config) (manifestInclusionConfiguration, error) {
	config := manifestInclusionConfiguration{}

	client, err := configv1client.NewForConfig(restConfig)
	if err != nil {
		return config, err
	}

	if featureGate, err := client.FeatureGates().Get(ctx, "cluster", metav1.GetOptions{}); err != nil {
		return config, err
	} else {
		config.RequiredFeatureSet = pointer.String(string(featureGate.Spec.FeatureSet))
	}

	if clusterVersion, err := client.ClusterVersions().Get(ctx, "version", metav1.GetOptions{}); err != nil {
		return config, err
	} else {
		config.Overrides = clusterVersion.Spec.Overrides
		config.Capabilities = &clusterVersion.Status.Capabilities

		// FIXME: eventually pull in GetImplicitlyEnabledCapabilities from https://github.com/openshift/cluster-version-operator/blob/86e24d66119a73f50282b66a8d6f2e3518aa0e15/pkg/payload/payload.go#L237-L240 for cases where a minor update would implicitly enable some additional capabilities.  For now, 4.13 to 4.14 will always enable MachineAPI.
		currentVersion := clusterVersion.Status.Desired.Version
		matches := regexp.MustCompile(`^(\d+[.]\d+)[.].*`).FindStringSubmatch(currentVersion)
		if len(matches) < 2 {
			return config, fmt.Errorf("failed to parse major.minor version from ClusterVersion status.desired.version %q", currentVersion)
		} else if matches[1] == "4.13" {
			machineAPI := configv1.ClusterVersionCapability("MachineAPI")
			config.Capabilities.EnabledCapabilities = append(config.Capabilities.EnabledCapabilities, machineAPI)
			config.Capabilities.KnownCapabilities = append(config.Capabilities.KnownCapabilities, machineAPI)
		}
	}

	if infrastructure, err := client.Infrastructures().Get(ctx, "cluster", metav1.GetOptions{}); err != nil {
		return config, err
	} else if infrastructure.Status.PlatformStatus == nil {
		return config, fmt.Errorf("cluster infrastructure does not declare status.platformStatus: %v", infrastructure.Status)
	} else {
		config.Platform = pointer.String(strings.ToLower(string(infrastructure.Status.PlatformStatus.Type)))
	}

	appsClient, err := appsv1client.NewForConfig(restConfig)
	if err != nil {
		return config, err
	}

	if deployment, err := appsClient.Deployments("openshift-cluster-version").Get(ctx, "cluster-version-operator", metav1.GetOptions{}); err != nil {
		return config, err
	} else {
		for _, container := range deployment.Spec.Template.Spec.Containers {
			for _, env := range container.Env {
				if env.Name == "CLUSTER_PROFILE" {
					config.Profile = pointer.String(env.Value)
					break
				}
			}
		}
	}

	return config, nil
}

func newIncluder(config manifestInclusionConfiguration) includer {
	return func(m *manifest.Manifest) error {
		return m.Include(config.ExcludeIdentifier, config.RequiredFeatureSet, config.Profile, config.Capabilities, config.Overrides)
	}
}
