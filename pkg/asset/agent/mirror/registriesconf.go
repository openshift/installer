package mirror

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"regexp"

	"github.com/containers/image/v5/pkg/sysregistriesv2"
	operatorv1alpha1 "github.com/openshift/api/operator/v1alpha1"
	"github.com/pelletier/go-toml"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/yaml"

	configv1 "github.com/openshift/api/config/v1"
	"github.com/openshift/installer/pkg/asset"
	"github.com/openshift/installer/pkg/asset/agent"
	"github.com/openshift/installer/pkg/asset/agent/joiner"
	"github.com/openshift/installer/pkg/asset/agent/workflow"
	"github.com/openshift/installer/pkg/asset/ignition/bootstrap"
	"github.com/openshift/installer/pkg/asset/releaseimage"
	"github.com/openshift/installer/pkg/types"
)

var (
	// RegistriesConfFilename defines the name of the file on disk
	RegistriesConfFilename = filepath.Join(mirrorConfigDir, "registries.conf")
)

// The default registries.conf file is the podman default as it appears in
// CoreOS, with no unqualified-search-registries.
const defaultRegistriesConf = `
# NOTE: RISK OF USING UNQUALIFIED IMAGE NAMES
# We recommend always using fully qualified image names including the registry
# server (full dns name), namespace, image name, and tag
# (e.g., registry.redhat.io/ubi8/ubi:latest). Pulling by digest (i.e.,
# quay.io/repository/name@digest) further eliminates the ambiguity of tags.
# When using short names, there is always an inherent risk that the image being
# pulled could be spoofed. For example, a user wants to pull an image named
# 'foobar' from a registry and expects it to come from myregistry.com. If
# myregistry.com is not first in the search list, an attacker could place a
# different 'foobar' image at a registry earlier in the search list. The user
# would accidentally pull and run the attacker's image and code rather than the
# intended content. We recommend only adding registries which are completely
# trusted (i.e., registries which don't allow unknown or anonymous users to
# create accounts with arbitrary names). This will prevent an image from being
# spoofed, squatted or otherwise made insecure.  If it is necessary to use one
# of these registries, it should be added at the end of the list.
#
# # An array of host[:port] registries to try when pulling an unqualified image, in order.

unqualified-search-registries = []

# [[registry]]
# # The "prefix" field is used to choose the relevant [[registry]] TOML table;
# # (only) the TOML table with the longest match for the input image name
# # (taking into account namespace/repo/tag/digest separators) is used.
# # 
# # The prefix can also be of the form: *.example.com for wildcard subdomain
# # matching.
# #
# # If the prefix field is missing, it defaults to be the same as the "location" field.
# prefix = "example.com/foo"
#
# # If true, unencrypted HTTP as well as TLS connections with untrusted
# # certificates are allowed.
# insecure = false
#
# # If true, pulling images with matching names is forbidden.
# blocked = false
#
# # The physical location of the "prefix"-rooted namespace.
# #
# # By default, this is equal to "prefix" (in which case "prefix" can be omitted
# # and the [[registry]] TOML table can only specify "location").
# #
# # Example: Given
# #   prefix = "example.com/foo"
# #   location = "internal-registry-for-example.net/bar"
# # requests for the image example.com/foo/myimage:latest will actually work with the
# # internal-registry-for-example.net/bar/myimage:latest image.
#
# # The location can be empty iff prefix is in a
# # wildcarded format: "*.example.com". In this case, the input reference will
# # be used as-is without any rewrite.
# location = internal-registry-for-example.com/bar"
#
# # (Possibly-partial) mirrors for the "prefix"-rooted namespace.
# #
# # The mirrors are attempted in the specified order; the first one that can be
# # contacted and contains the image will be used (and if none of the mirrors contains the image,
# # the primary location specified by the "registry.location" field, or using the unmodified
# # user-specified reference, is tried last).
# #
# # Each TOML table in the "mirror" array can contain the following fields, with the same semantics
# # as if specified in the [[registry]] TOML table directly:
# # - location
# # - insecure
# [[registry.mirror]]
# location = "example-mirror-0.local/mirror-for-foo"
# [[registry.mirror]]
# location = "example-mirror-1.local/mirrors/foo"
# insecure = true
# # Given the above, a pull of example.com/foo/image:latest will try:
# # 1. example-mirror-0.local/mirror-for-foo/image:latest
# # 2. example-mirror-1.local/mirrors/foo/image:latest
# # 3. internal-registry-for-example.net/bar/image:latest
# # in order, and use the first one that exists.
`

// RegistriesConf generates the registries.conf file.
type RegistriesConf struct {
	File         *asset.File
	Config       *sysregistriesv2.V2RegistriesConf
	MirrorConfig []RegistriesConfig
}

// RegistriesConfig holds the data extracted from registries.conf
type RegistriesConfig struct {
	Location string
	Mirrors  []string
}

var _ asset.WritableAsset = (*RegistriesConf)(nil)

// Name returns a human friendly name for the asset.
func (*RegistriesConf) Name() string {
	return "Mirror Registries Config"
}

// Dependencies returns all of the dependencies directly needed to generate
// the asset.
func (*RegistriesConf) Dependencies() []asset.Asset {
	return []asset.Asset{
		&workflow.AgentWorkflow{},
		&joiner.ClusterInfo{},
		&agent.OptionalInstallConfig{},
		&releaseimage.Image{},
	}
}

// Generate generates the registries.conf file from install-config.
func (i *RegistriesConf) Generate(_ context.Context, dependencies asset.Parents) error {
	agentWorkflow := &workflow.AgentWorkflow{}
	clusterInfo := &joiner.ClusterInfo{}
	installConfig := &agent.OptionalInstallConfig{}
	releaseImage := &releaseimage.Image{}
	dependencies.Get(installConfig, releaseImage, agentWorkflow, clusterInfo)

	var imageDigestSources []types.ImageDigestSource
	var deprecatedImageContentSources []types.ImageContentSource
	var image string

	switch agentWorkflow.Workflow {
	case workflow.AgentWorkflowTypeInstall:
		if installConfig.Supplied {
			imageDigestSources = installConfig.Config.ImageDigestSources
			deprecatedImageContentSources = installConfig.Config.DeprecatedImageContentSources
		}
		image = releaseImage.PullSpec

	case workflow.AgentWorkflowTypeAddNodes:
		imageDigestSources = clusterInfo.ImageDigestSources
		deprecatedImageContentSources = clusterInfo.DeprecatedImageContentSources
		image = clusterInfo.ReleaseImage

	default:
		return fmt.Errorf("AgentWorkflowType value not supported: %s", agentWorkflow.Workflow)
	}

	if len(deprecatedImageContentSources) == 0 && len(imageDigestSources) == 0 {
		return i.generateDefaultRegistriesConf()
	}

	err := i.generateRegistriesConf(imageDigestSources, deprecatedImageContentSources)
	if err != nil {
		return err
	}

	if !i.releaseImageIsSameInRegistriesConf(image) {
		logrus.Warnf("The imageDigestSources configuration in install-config.yaml should have at least one source field matching the releaseImage value %s", releaseImage.PullSpec)
	}

	registriesData, err := toml.Marshal(i.Config)
	if err != nil {
		return err
	}

	i.File = &asset.File{
		Filename: RegistriesConfFilename,
		Data:     registriesData,
	}

	return nil
}

func (i *RegistriesConf) generateRegistriesConf(imageDigestSources []types.ImageDigestSource, deprecatedImageContentSources []types.ImageContentSource) error {
	if len(deprecatedImageContentSources) != 0 && len(imageDigestSources) != 0 {
		return fmt.Errorf("invalid install-config.yaml, cannot set imageContentSources and imageDigestSources at the same time")
	}

	digestMirrorSources := []types.ImageDigestSource{}
	if len(deprecatedImageContentSources) > 0 {
		digestMirrorSources = bootstrap.ContentSourceToDigestMirror(deprecatedImageContentSources)
	} else if len(imageDigestSources) > 0 {
		digestMirrorSources = append(digestMirrorSources, imageDigestSources...)
	}

	registries := &sysregistriesv2.V2RegistriesConf{
		Registries: []sysregistriesv2.Registry{},
	}
	for _, group := range bootstrap.MergedMirrorSets(digestMirrorSources) {
		if len(group.Mirrors) == 0 {
			continue
		}

		registry := sysregistriesv2.Registry{}
		registry.Endpoint.Location = group.Source
		registry.MirrorByDigestOnly = true
		registry.Blocked = group.SourcePolicy == configv1.NeverContactSource
		for _, mirror := range group.Mirrors {
			registry.Mirrors = append(registry.Mirrors, sysregistriesv2.Endpoint{Location: mirror})
		}
		registries.Registries = append(registries.Registries, registry)
	}
	i.Config = registries
	i.setMirrorConfig(i.Config)

	return nil
}

// HasMirrors returns whether there are any mirrors configured.
func (i *RegistriesConf) HasMirrors() bool {
	return len(i.MirrorConfig) > 0
}

// GetICSPContents converts the data in registries.conf into ICSP format.
func (i *RegistriesConf) GetICSPContents() ([]byte, error) {
	icsp := operatorv1alpha1.ImageContentSourcePolicy{
		TypeMeta: metav1.TypeMeta{
			APIVersion: operatorv1alpha1.SchemeGroupVersion.String(),
			Kind:       "ImageContentSourcePolicy",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name: "image-policy",
			// not namespaced
		},
	}

	icsp.Spec.RepositoryDigestMirrors = make([]operatorv1alpha1.RepositoryDigestMirrors, len(i.MirrorConfig))
	for i, mirrorRegistries := range i.MirrorConfig {
		icsp.Spec.RepositoryDigestMirrors[i] = operatorv1alpha1.RepositoryDigestMirrors{Source: mirrorRegistries.Location, Mirrors: mirrorRegistries.Mirrors}
	}

	// Convert to json first so json tags are handled
	jsonData, err := json.Marshal(&icsp)
	if err != nil {
		return nil, err
	}
	contents, err := yaml.JSONToYAML(jsonData)
	if err != nil {
		return nil, err
	}

	return contents, nil
}

// Files returns the files generated by the asset.
func (i *RegistriesConf) Files() []*asset.File {
	if i.File != nil {
		return []*asset.File{i.File}
	}
	return []*asset.File{}
}

// Load returns RegistriesConf asset from the disk.
func (i *RegistriesConf) Load(f asset.FileFetcher) (bool, error) {
	ctx := context.TODO()

	releaseImage := &releaseimage.Image{}
	if err := releaseImage.Generate(ctx, asset.Parents{}); err != nil {
		return false, fmt.Errorf("failed to generate the release image asset: %w", err)
	}

	file, err := f.FetchByName(RegistriesConfFilename)
	if err != nil {
		if os.IsNotExist(err) {
			return false, nil
		}
		return false, errors.Wrap(err, fmt.Sprintf("failed to load %s file", RegistriesConfFilename))
	}

	registriesConf := &sysregistriesv2.V2RegistriesConf{}
	if err := toml.Unmarshal(file.Data, registriesConf); err != nil {
		return false, errors.Wrapf(err, "failed to unmarshal %s", RegistriesConfFilename)
	}

	i.File, i.Config = file, registriesConf
	i.setMirrorConfig(i.Config)

	if string(i.File.Data) != defaultRegistriesConf {
		if i.validateRegistriesConf() {
			if !i.releaseImageIsSameInRegistriesConf(releaseImage.PullSpec) {
				logrus.Warnf("%s should have an entry matching the releaseImage %s", RegistriesConfFilename, releaseImage.PullSpec)
			}
		}
	}

	return true, nil
}

func (i *RegistriesConf) validateRegistriesConf() bool {
	for _, registry := range i.Config.Registries {
		if registry.Endpoint.Location == "" {
			logrus.Warnf("Location key not found in %s", RegistriesConfFilename)
			return false
		}
	}
	return true
}

func (i *RegistriesConf) releaseImageIsSameInRegistriesConf(releaseImage string) bool {
	return GetMirrorFromRelease(releaseImage, i) != ""
}

func (i *RegistriesConf) generateDefaultRegistriesConf() error {
	i.File = &asset.File{
		Filename: RegistriesConfFilename,
		Data:     []byte(defaultRegistriesConf),
	}
	registriesConf := &sysregistriesv2.V2RegistriesConf{}
	if err := toml.Unmarshal([]byte(defaultRegistriesConf), registriesConf); err != nil {
		return errors.Wrapf(err, "failed to unmarshal %s", RegistriesConfFilename)
	}
	i.Config = registriesConf
	return nil
}

func (i *RegistriesConf) setMirrorConfig(registriesConf *sysregistriesv2.V2RegistriesConf) {
	mirrorConfig := make([]RegistriesConfig, len(registriesConf.Registries))
	for i, reg := range registriesConf.Registries {
		mirrorConfig[i] = RegistriesConfig{
			Location: reg.Location,
		}
		for _, mirror := range reg.Mirrors {
			mirrorConfig[i].Mirrors = append(mirrorConfig[i].Mirrors, mirror.Location)
		}
	}
	i.MirrorConfig = mirrorConfig
}

// GetMirrorFromRelease gets the matching mirror configured for the releaseImage.
// If multiple mirrors are configured the first one is returned.
func GetMirrorFromRelease(releaseImage string, registriesConfig *RegistriesConf) string {
	source := regexp.MustCompile(`^(.+?)(@sha256)?:(.+)`).FindStringSubmatch(releaseImage)
	for _, config := range registriesConfig.MirrorConfig {
		if config.Location == source[1] {
			// include the tag with the build release image
			switch len(source) {
			case 4:
				// Has Sha256
				return fmt.Sprintf("%s%s:%s", config.Mirrors[0], source[2], source[3])
			case 3:
				return fmt.Sprintf("%s:%s", config.Mirrors[0], source[2])
			}
		}
	}

	return ""
}
