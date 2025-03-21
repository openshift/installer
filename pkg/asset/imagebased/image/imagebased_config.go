package image

import (
	"context"
	"errors"
	"fmt"
	"net"
	"net/url"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	dockerref "github.com/containers/image/v5/docker/reference"
	"github.com/sirupsen/logrus"
	"github.com/thoas/go-funk"
	"k8s.io/apimachinery/pkg/util/validation/field"
	"sigs.k8s.io/yaml"

	"github.com/openshift/installer/pkg/asset"
	"github.com/openshift/installer/pkg/asset/agent/manifests/staticnetworkconfig"
	"github.com/openshift/installer/pkg/types/imagebased"
	"github.com/openshift/installer/pkg/validate"
)

const (
	coreosInstallerArgsValuesRegex = `^[A-Za-z0-9@!#$%*()_+-=//.,";':{}\[\]]+$`
)

var (
	configFilename              = "image-based-installation-config.yaml"
	allowedFlags                = []string{"--append-karg", "--delete-karg", "--save-partlabel", "--save-partindex"}
	defaultExtraPartitionLabel  = "var-lib-containers"
	defaultExtraPartitionStart  = "-40G"
	defaultExtraPartitionNumber = uint(5)
)

// ImageBasedInstallationConfig reads the image-based-installation-config.yaml file.
type ImageBasedInstallationConfig struct { // nolint:revive // although this name stutters it is useful to convey that it's an image-based installer related struct
	File     *asset.File
	Config   *imagebased.InstallationConfig
	Template string
}

var _ asset.WritableAsset = (*ImageBasedInstallationConfig)(nil)

// Name returns a human friendly name for the asset.
func (*ImageBasedInstallationConfig) Name() string {
	return "Image-based Installation ISO Config"
}

// Dependencies returns all of the dependencies directly needed to generate
// the asset.
func (*ImageBasedInstallationConfig) Dependencies() []asset.Asset {
	return []asset.Asset{}
}

// Generate generates the Image-based Installation Config YAML manifest.
func (i *ImageBasedInstallationConfig) Generate(_ context.Context, dependencies asset.Parents) error {
	configTemplate := `#
# Note: This is a sample ImageBasedInstallationConfig file showing
# which fields are available to aid you in creating your
# own image-based-installation-config.yaml file.
#
apiVersion: v1beta1
kind: ImageBasedInstallationConfig
metadata:
  name: example-image-based-installation-config
# The following fields are required
seedImage: quay.io/openshift-kni/seed-image:4.16.0
seedVersion: 4.16.0
installationDisk: /dev/vda
pullSecret: '<your-pull-secret>'
# networkConfig is optional and contains the network configuration for the host in NMState format.
# See https://nmstate.io/examples.html for examples.
# networkConfig:
#   interfaces:
#     - name: eth0
#       type: ethernet
#       state: up
#       mac-address: 00:00:00:00:00:00
#       ipv4:
#         enabled: true
#         address:
#           - ip: 192.168.122.2
#             prefix-length: 23
#         dhcp: false
`

	i.Template = configTemplate

	// Set the File field correctly with the generated image-based installation config YAML content.
	i.File = &asset.File{
		Filename: configFilename,
		Data:     []byte(i.Template),
	}

	return nil
}

// PersistToFile writes the image-based-installation-config.yaml file to the assets folder.
func (i *ImageBasedInstallationConfig) PersistToFile(directory string) error {
	if i.File == nil {
		return nil
	}

	configPath := filepath.Join(directory, configFilename)
	err := os.WriteFile(configPath, i.File.Data, 0o600)
	if err != nil {
		return err
	}

	return nil
}

// Files returns the files generated by the asset.
func (i *ImageBasedInstallationConfig) Files() []*asset.File {
	if i.File != nil {
		return []*asset.File{i.File}
	}
	return []*asset.File{}
}

// Load returns image-based installation ISO config asset from the disk.
func (i *ImageBasedInstallationConfig) Load(f asset.FileFetcher) (bool, error) {
	file, err := f.FetchByName(configFilename)
	if err != nil {
		if os.IsNotExist(err) {
			return false, nil
		}
		return false, fmt.Errorf("failed to load %s file: %w", configFilename, err)
	}
	config := &imagebased.InstallationConfig{
		ExtraPartitionLabel:  defaultExtraPartitionLabel,
		ExtraPartitionStart:  defaultExtraPartitionStart,
		ExtraPartitionNumber: defaultExtraPartitionNumber,
	}
	if err := yaml.UnmarshalStrict(file.Data, config); err != nil {
		return false, fmt.Errorf("failed to unmarshal %s: %w", configFilename, err)
	}

	i.File, i.Config = file, config
	if err = i.finish(); err != nil {
		return false, err
	}
	return true, nil
}

func (i *ImageBasedInstallationConfig) finish() error {
	if err := i.validate().ToAggregate(); err != nil {
		return fmt.Errorf("invalid Image-based Installation ISO Config: %w", err)
	}
	return nil
}

func (i *ImageBasedInstallationConfig) validate() field.ErrorList {
	var allErrs field.ErrorList

	if err := i.validatePullSecret(); err != nil {
		allErrs = append(allErrs, err...)
	}
	if err := i.validateSSHKey(); err != nil {
		allErrs = append(allErrs, err...)
	}
	if err := i.validateSeedImage(); err != nil {
		allErrs = append(allErrs, err...)
	}
	if err := i.validateSeedVersion(); err != nil {
		allErrs = append(allErrs, err...)
	}
	if err := i.validateInstallationDisk(); err != nil {
		allErrs = append(allErrs, err...)
	}
	if err := i.validateExtraPartitionStart(); err != nil {
		allErrs = append(allErrs, err...)
	}
	if err := i.validateAdditionalTrustBundle(); err != nil {
		allErrs = append(allErrs, err...)
	}
	if err := i.validateNetworkConfig(); err != nil {
		allErrs = append(allErrs, err...)
	}
	if err := i.validateImageDigestSources(); err != nil {
		allErrs = append(allErrs, err...)
	}
	if err := i.validateProxy(); err != nil {
		allErrs = append(allErrs, err...)
	}
	if err := i.validateCoreosInstallerArgs(); err != nil {
		allErrs = append(allErrs, err...)
	}

	return allErrs
}

func (i *ImageBasedInstallationConfig) validatePullSecret() field.ErrorList {
	var allErrs field.ErrorList

	pullSecretPath := field.NewPath("pullSecret")

	if i.Config.PullSecret == "" {
		allErrs = append(allErrs, field.Required(pullSecretPath, "you must specify a pullSecret"))
		return allErrs
	}

	if err := validate.ImagePullSecret(i.Config.PullSecret); err != nil {
		allErrs = append(allErrs, field.Invalid(pullSecretPath, i.Config.PullSecret, err.Error()))
	}

	return allErrs
}

func (i *ImageBasedInstallationConfig) validateSSHKey() field.ErrorList {
	var allErrs field.ErrorList

	// empty SSHKey is fine
	if i.Config.SSHKey == "" {
		return nil
	}

	sshKeyPath := field.NewPath("sshKey")

	if err := validate.SSHPublicKey(i.Config.SSHKey); err != nil {
		allErrs = append(allErrs, field.Invalid(sshKeyPath, i.Config.SSHKey, err.Error()))
	}

	return allErrs
}

func (i *ImageBasedInstallationConfig) validateAdditionalTrustBundle() field.ErrorList {
	var allErrs field.ErrorList

	// empty AdditionalTrustBundle is fine
	if i.Config.AdditionalTrustBundle == "" {
		return nil
	}

	additionalTrustBundlePath := field.NewPath("additionalTrustBundle")

	if err := validate.CABundle(i.Config.AdditionalTrustBundle); err != nil {
		allErrs = append(allErrs, field.Invalid(additionalTrustBundlePath, i.Config.AdditionalTrustBundle, err.Error()))
	}

	return allErrs
}

func (i *ImageBasedInstallationConfig) validateSeedImage() field.ErrorList {
	var allErrs field.ErrorList

	seedImagePath := field.NewPath("seedImage")

	if i.Config.SeedImage == "" {
		allErrs = append(allErrs, field.Required(seedImagePath, "you must specify a seedImage"))
	}

	return allErrs
}

func (i *ImageBasedInstallationConfig) validateSeedVersion() field.ErrorList {
	var allErrs field.ErrorList

	seedVersionPath := field.NewPath("seedVersion")

	if i.Config.SeedVersion == "" {
		allErrs = append(allErrs, field.Required(seedVersionPath, "you must specify a seedVersion"))
	}

	return allErrs
}

func (i *ImageBasedInstallationConfig) validateInstallationDisk() field.ErrorList {
	var allErrs field.ErrorList

	installationDiskPath := field.NewPath("installationDisk")

	if i.Config.InstallationDisk == "" {
		allErrs = append(allErrs, field.Required(installationDiskPath, "you must specify an installationDisk"))
	}

	return allErrs
}

func (i *ImageBasedInstallationConfig) validateNetworkConfig() field.ErrorList {
	var allErrs field.ErrorList

	// empty NetworkConfig is fine
	if i.Config.NetworkConfig == nil || i.Config.NetworkConfig.String() == "" {
		return nil
	}

	networkConfig := field.NewPath("networkConfig")

	staticNetworkConfigGenerator := staticnetworkconfig.New(logrus.StandardLogger(), staticnetworkconfig.Config{MaxConcurrentGenerations: 2})

	// Validate the network config using nmstatectl.
	if err := staticNetworkConfigGenerator.ValidateNMStateYaml(context.Background(), i.Config.NetworkConfig.String()); err != nil {
		allErrs = append(allErrs, field.Invalid(networkConfig, i.Config.NetworkConfig, err.Error()))
	}

	return allErrs
}

func (i *ImageBasedInstallationConfig) validateImageDigestSources() field.ErrorList {
	allErrs := field.ErrorList{}

	fldPath := field.NewPath("imageDigestSources")

	for gidx, group := range i.Config.ImageDigestSources {
		groupf := fldPath.Index(gidx)
		if err := validateNamedRepository(group.Source); err != nil {
			allErrs = append(allErrs, field.Invalid(groupf.Child("source"), group.Source, err.Error()))
		}

		for midx, mirror := range group.Mirrors {
			if err := validateNamedRepository(mirror); err != nil {
				allErrs = append(allErrs, field.Invalid(groupf.Child("mirrors").Index(midx), mirror, err.Error()))
				continue
			}
		}
	}
	return allErrs
}

func validateNamedRepository(r string) error {
	ref, err := dockerref.ParseNamed(r)
	if err != nil {
		// If a mirror name is provided without the named reference,
		// then the name is not considered canonical and will cause
		// an error. e.g. registry.lab.redhat.com:5000 will result
		// in an error. Instead we will check whether the input is
		// a valid hostname as a workaround.
		if errors.Is(err, dockerref.ErrNameNotCanonical) {
			// If the hostname string contains a port, lets attempt
			// to split them
			host, _, err := net.SplitHostPort(r)
			if err != nil {
				host = r
			}
			if err = validate.Host(host); err != nil {
				return fmt.Errorf("the repository provided is invalid: %w", err)
			}
			return nil
		}
		return fmt.Errorf("failed to parse: %w", err)
	}
	if !dockerref.IsNameOnly(ref) {
		return errors.New("must be repository--not reference")
	}
	return nil
}

func (i *ImageBasedInstallationConfig) validateProxy() field.ErrorList {
	allErrs := field.ErrorList{}

	// empty Proxy is fine
	if i.Config.Proxy == nil {
		return nil
	}

	fldPath := field.NewPath("proxy")

	if i.Config.Proxy.HTTPProxy == "" && i.Config.Proxy.HTTPSProxy == "" {
		allErrs = append(allErrs, field.Required(fldPath, "must include httpProxy or httpsProxy"))
	}

	if i.Config.Proxy.HTTPProxy != "" {
		allErrs = append(allErrs, validateURI(i.Config.Proxy.HTTPProxy, fldPath.Child("httpProxy"), []string{"http"})...)
	}
	if i.Config.Proxy.HTTPSProxy != "" {
		allErrs = append(allErrs, validateURI(i.Config.Proxy.HTTPSProxy, fldPath.Child("httpsProxy"), []string{"http", "https"})...)
	}
	if i.Config.Proxy.NoProxy != "" && i.Config.Proxy.NoProxy != "*" {
		for idx, v := range strings.Split(i.Config.Proxy.NoProxy, ",") {
			v = strings.TrimSpace(v)
			errDomain := validate.NoProxyDomainName(v)
			_, _, errCIDR := net.ParseCIDR(v)
			ip := net.ParseIP(v)
			if errDomain != nil && errCIDR != nil && ip == nil {
				allErrs = append(allErrs, field.Invalid(fldPath.Child("noProxy"), i.Config.Proxy.NoProxy, fmt.Sprintf(
					"each element of noProxy must be a IP, CIDR or domain without wildcard characters, which is violated by element %d %q", idx, v)))
			}
		}
	}

	return allErrs
}

func (i *ImageBasedInstallationConfig) validateCoreosInstallerArgs() field.ErrorList {
	var allErrs field.ErrorList
	argsRe := regexp.MustCompile("^-+.*")
	valuesRe := regexp.MustCompile(coreosInstallerArgsValuesRegex)

	coreosInstallerParamsPath := field.NewPath("CoreosInstallerArgs")

	for _, arg := range i.Config.CoreosInstallerArgs {
		if argsRe.MatchString(arg) {
			if !funk.ContainsString(allowedFlags, arg) {
				allErrs = append(allErrs, field.Required(coreosInstallerParamsPath,
					fmt.Sprintf("found unexpected flag %s for coreosInstallerArgs - allowed flags are %v", arg, allowedFlags)))
			}
			continue
		}

		if !valuesRe.MatchString(arg) {
			allErrs = append(allErrs, field.Required(coreosInstallerParamsPath, fmt.Sprintf("found unexpected chars in value %s for installer", arg)))
		}
	}

	return allErrs
}

func validateURI(uri string, fldPath *field.Path, schemes []string) field.ErrorList {
	parsed, err := url.ParseRequestURI(uri)
	if err != nil {
		return field.ErrorList{field.Invalid(fldPath, uri, err.Error())}
	}
	for _, scheme := range schemes {
		if scheme == parsed.Scheme {
			return nil
		}
	}
	return field.ErrorList{field.NotSupported(fldPath, parsed.Scheme, schemes)}
}

func (i *ImageBasedInstallationConfig) validateExtraPartitionStart() field.ErrorList {
	var allErrs field.ErrorList
	extraPartitionStartPath := field.NewPath("ExtraPartitionStart")

	start := i.Config.ExtraPartitionStart
	if start == "" {
		allErrs = append(allErrs, field.Required(extraPartitionStartPath, "partition start sector cannot be empty"))
		return allErrs
	}

	if start == "0" {
		return allErrs
	}

	// Matches patterns like: 10K, +10K, -20M, 1G, +1G, -2T, 3P, +3P.
	//
	// First group is optional: + or -.
	// Second group is numeric value.
	// Third group is K,M,G,T,P suffix.
	validFormat := regexp.MustCompile(`^([+-])?(\d+)([KMGTP])$`)

	if !validFormat.MatchString(start) {
		allErrs = append(allErrs, field.Invalid(extraPartitionStartPath, start, "partition start must be '0' or match pattern [+-]?<number>[KMGTP]"))
	}

	return allErrs
}
