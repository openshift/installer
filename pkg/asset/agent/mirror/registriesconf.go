package mirror

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/containers/image/pkg/sysregistriesv2"
	"github.com/openshift/installer/pkg/asset"
	"github.com/openshift/installer/pkg/asset/agent"
	"github.com/openshift/installer/pkg/asset/ignition/bootstrap"
	"github.com/pelletier/go-toml"
	"github.com/pkg/errors"
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
	MirrorConfig []RegistriesConfig
}

// RegistriesConfig holds the data extracted from registries.conf
type RegistriesConfig struct {
	Location string
	Mirror   string
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
		&agent.OptionalInstallConfig{},
	}
}

// Generate generates the registries.conf file from install-config.
func (i *RegistriesConf) Generate(dependencies asset.Parents) error {
	installConfig := &agent.OptionalInstallConfig{}
	dependencies.Get(installConfig)
	if !installConfig.Supplied || len(installConfig.Config.ImageContentSources) == 0 {
		i.File = &asset.File{
			Filename: RegistriesConfFilename,
			Data:     []byte(defaultRegistriesConf),
		}
		return i.finish()
	}

	registries := sysregistriesv2.V2RegistriesConf{
		Registries: []sysregistriesv2.Registry{},
	}
	for _, group := range bootstrap.MergedMirrorSets(installConfig.Config.ImageContentSources) {
		if len(group.Mirrors) == 0 {
			continue
		}

		registry := sysregistriesv2.Registry{}
		registry.Endpoint.Location = group.Source
		registry.MirrorByDigestOnly = true
		for _, mirror := range group.Mirrors {
			registry.Mirrors = append(registry.Mirrors, sysregistriesv2.Endpoint{Location: mirror})
		}
		registries.Registries = append(registries.Registries, registry)
	}

	data, err := toml.Marshal(registries)
	if err != nil {
		return err
	}

	i.File = &asset.File{
		Filename: RegistriesConfFilename,
		Data:     data,
	}

	return i.finish()
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

	file, err := f.FetchByName(RegistriesConfFilename)
	if err != nil {
		if os.IsNotExist(err) {
			return false, nil
		}
		return false, errors.Wrap(err, fmt.Sprintf("failed to load %s file", RegistriesConfFilename))
	}

	i.File = file

	if err = i.finish(); err != nil {
		return false, err
	}

	return true, nil
}

func (i *RegistriesConf) finish() error {

	config, err := extractLocationMirrorDataFromRegistries(i.File.Data)
	if err != nil {
		return errors.Wrap(err, fmt.Sprintf("failed to parse mirrors in %s", RegistriesConfFilename))
	}

	i.MirrorConfig = config

	return nil
}

// From assisted-service pkg/mirrorregistries/generator.go
func extractLocationMirrorDataFromRegistries(registriesConfToml []byte) ([]RegistriesConfig, error) {
	registries := sysregistriesv2.V2RegistriesConf{}
	err := toml.Unmarshal(registriesConfToml, &registries)
	if err != nil {
		return nil, err
	}

	registriesConfList := make([]RegistriesConfig, len(registries.Registries))
	for i, reg := range registries.Registries {
		registriesConfList[i] = RegistriesConfig{
			Location: reg.Location,
			Mirror:   reg.Mirrors[0].Location,
		}
	}

	return registriesConfList, nil
}
