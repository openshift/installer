package bootstrap

import (
	"bytes"
	"crypto/x509"
	"encoding/json"
	"encoding/pem"
	"fmt"
	"io"
	"net"
	"os"
	"path"
	"path/filepath"
	"strings"
	"text/template"
	"time"

	"github.com/containers/image/pkg/sysregistriesv2"
	ignutil "github.com/coreos/ignition/v2/config/util"
	igntypes "github.com/coreos/ignition/v2/config/v3_2/types"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"github.com/vincent-petithory/dataurl"
	utilsnet "k8s.io/utils/net"

	configv1 "github.com/openshift/api/config/v1"
	"github.com/openshift/installer/data"
	"github.com/openshift/installer/pkg/asset"
	"github.com/openshift/installer/pkg/asset/ignition"
	"github.com/openshift/installer/pkg/asset/ignition/bootstrap/baremetal"
	"github.com/openshift/installer/pkg/asset/ignition/bootstrap/vsphere"
	mcign "github.com/openshift/installer/pkg/asset/ignition/machine"
	"github.com/openshift/installer/pkg/asset/installconfig"
	"github.com/openshift/installer/pkg/asset/kubeconfig"
	"github.com/openshift/installer/pkg/asset/machines"
	"github.com/openshift/installer/pkg/asset/manifests"
	"github.com/openshift/installer/pkg/asset/releaseimage"
	"github.com/openshift/installer/pkg/asset/rhcos"
	"github.com/openshift/installer/pkg/asset/tls"
	"github.com/openshift/installer/pkg/types"
	baremetaltypes "github.com/openshift/installer/pkg/types/baremetal"
	vspheretypes "github.com/openshift/installer/pkg/types/vsphere"
)

const (
	rootDir = "/opt/openshift"
)

var (
	commonEnabledServices = []string{
		"progress.service",
		"kubelet.service",
		"chown-gatewayd-key.service",
		"systemd-journal-gatewayd.socket",
		"approve-csr.service",
		// baremetal & openstack platform services
		"keepalived.service",
		"coredns.service",
		"ironic.service",
		"master-bmh-update.service",
	}
)

// bootstrapTemplateData is the data to use to replace values in bootstrap
// template files.
type bootstrapTemplateData struct {
	AdditionalTrustBundle string
	FIPS                  bool
	EtcdCluster           string
	PullSecret            string
	SSHKey                string
	ReleaseImage          string
	ClusterProfile        string
	Proxy                 *configv1.ProxyStatus
	Registries            []sysregistriesv2.Registry
	BootImage             string
	PlatformData          platformTemplateData
	BootstrapInPlace      *types.BootstrapInPlace
	UseIPv6ForNodeIP      bool
	UseDualForNodeIP      bool
	IsFCOS                bool
	IsSCOS                bool
	IsOKD                 bool
	BootstrapNodeIP       string
	APIServerURL          string
	APIIntServerURL       string
	FeatureSet            configv1.FeatureSet
}

// platformTemplateData is the data to use to replace values in bootstrap
// template files that are specific to one platform.
type platformTemplateData struct {
	BareMetal *baremetal.TemplateData
	VSphere   *vsphere.TemplateData
}

// Common is an asset that generates the ignition config for bootstrap nodes.
type Common struct {
	Config *igntypes.Config
	File   *asset.File
}

// Dependencies returns the assets on which the Bootstrap asset depends.
func (a *Common) Dependencies() []asset.Asset {
	return []asset.Asset{
		&baremetal.IronicCreds{},
		&CVOIgnore{},
		&installconfig.InstallConfig{},
		&kubeconfig.AdminInternalClient{},
		&kubeconfig.Kubelet{},
		&kubeconfig.LoopbackClient{},
		&mcign.MasterIgnitionCustomizations{},
		&mcign.WorkerIgnitionCustomizations{},
		&machines.Master{},
		&machines.Worker{},
		&manifests.Manifests{},
		&manifests.Openshift{},
		&manifests.Proxy{},
		&tls.AdminKubeConfigCABundle{},
		&tls.AggregatorCA{},
		&tls.AggregatorCABundle{},
		&tls.AggregatorClientCertKey{},
		&tls.AggregatorSignerCertKey{},
		&tls.APIServerProxyCertKey{},
		&tls.BootstrapSSHKeyPair{},
		&tls.BoundSASigningKey{},
		&tls.CloudProviderCABundle{},
		&tls.JournalCertKey{},
		&tls.KubeAPIServerLBCABundle{},
		&tls.KubeAPIServerExternalLBServerCertKey{},
		&tls.KubeAPIServerInternalLBServerCertKey{},
		&tls.KubeAPIServerLBSignerCertKey{},
		&tls.KubeAPIServerLocalhostCABundle{},
		&tls.KubeAPIServerLocalhostServerCertKey{},
		&tls.KubeAPIServerLocalhostSignerCertKey{},
		&tls.KubeAPIServerServiceNetworkCABundle{},
		&tls.KubeAPIServerServiceNetworkServerCertKey{},
		&tls.KubeAPIServerServiceNetworkSignerCertKey{},
		&tls.KubeAPIServerCompleteCABundle{},
		&tls.KubeAPIServerCompleteClientCABundle{},
		&tls.KubeAPIServerToKubeletCABundle{},
		&tls.KubeAPIServerToKubeletClientCertKey{},
		&tls.KubeAPIServerToKubeletSignerCertKey{},
		&tls.KubeControlPlaneCABundle{},
		&tls.KubeControlPlaneKubeControllerManagerClientCertKey{},
		&tls.KubeControlPlaneKubeSchedulerClientCertKey{},
		&tls.KubeControlPlaneSignerCertKey{},
		&tls.KubeletBootstrapCABundle{},
		&tls.KubeletClientCABundle{},
		&tls.KubeletClientCertKey{},
		&tls.KubeletCSRSignerCertKey{},
		&tls.KubeletServingCABundle{},
		&tls.MCSCertKey{},
		&tls.RootCA{},
		&tls.ServiceAccountKeyPair{},
		&releaseimage.Image{},
		new(rhcos.Image),
	}
}

// Generate generates the ignition config for the Bootstrap asset.
func (a *Common) generateConfig(dependencies asset.Parents, templateData *bootstrapTemplateData) error {
	installConfig := &installconfig.InstallConfig{}
	bootstrapSSHKeyPair := &tls.BootstrapSSHKeyPair{}
	dependencies.Get(installConfig, bootstrapSSHKeyPair)

	a.Config = &igntypes.Config{
		Ignition: igntypes.Ignition{
			Version: igntypes.MaxVersion.String(),
		},
	}

	if err := AddStorageFiles(a.Config, "/", "bootstrap/files", templateData); err != nil {
		return err
	}
	if err := AddSystemdUnits(a.Config, "bootstrap/systemd/units", templateData, commonEnabledServices); err != nil {
		return err
	}

	// Check for optional platform specific files/units
	platform := installConfig.Config.Platform.Name()
	platformFilePath := fmt.Sprintf("bootstrap/%s/files", platform)
	directory, err := data.Assets.Open(platformFilePath)
	if err == nil {
		directory.Close()
		err = AddStorageFiles(a.Config, "/", platformFilePath, templateData)
		if err != nil {
			return err
		}
	}

	platformUnitPath := fmt.Sprintf("bootstrap/%s/systemd/units", platform)
	directory, err = data.Assets.Open(platformUnitPath)
	if err == nil {
		directory.Close()
		if err = AddSystemdUnits(a.Config, platformUnitPath, templateData, commonEnabledServices); err != nil {
			return err
		}
	}

	a.addParentFiles(dependencies)

	a.Config.Passwd.Users = append(
		a.Config.Passwd.Users,
		igntypes.PasswdUser{Name: "core", SSHAuthorizedKeys: []igntypes.SSHAuthorizedKey{
			igntypes.SSHAuthorizedKey(installConfig.Config.SSHKey),
			igntypes.SSHAuthorizedKey(string(bootstrapSSHKeyPair.Public())),
		}},
	)

	return nil
}

func (a *Common) generateFile(filename string) error {
	data, err := ignition.Marshal(a.Config)
	if err != nil {
		return errors.Wrap(err, "failed to Marshal Ignition config")
	}
	a.File = &asset.File{
		Filename: filename,
		Data:     data,
	}
	return nil
}

// Files returns the files generated by the asset.
func (a *Common) Files() []*asset.File {
	if a.File != nil {
		return []*asset.File{a.File}
	}
	return nil
}

// getTemplateData returns the data to use to execute bootstrap templates.
func (a *Common) getTemplateData(dependencies asset.Parents, bootstrapInPlace bool) *bootstrapTemplateData {
	installConfig := &installconfig.InstallConfig{}
	proxy := &manifests.Proxy{}
	releaseImage := &releaseimage.Image{}
	rhcosImage := new(rhcos.Image)
	bootstrapSSHKeyPair := &tls.BootstrapSSHKeyPair{}
	ironicCreds := &baremetal.IronicCreds{}
	dependencies.Get(installConfig, proxy, releaseImage, rhcosImage, bootstrapSSHKeyPair, ironicCreds)

	etcdEndpoints := make([]string, *installConfig.Config.ControlPlane.Replicas)

	for i := range etcdEndpoints {
		etcdEndpoints[i] = fmt.Sprintf("https://etcd-%d.%s:2379", i, installConfig.Config.ClusterDomain())
	}

	registries := []sysregistriesv2.Registry{}
	for _, group := range MergedMirrorSets(installConfig.Config.ImageContentSources) {
		if len(group.Mirrors) == 0 {
			continue
		}

		registry := sysregistriesv2.Registry{}
		registry.Endpoint.Location = group.Source
		registry.MirrorByDigestOnly = true
		for _, mirror := range group.Mirrors {
			registry.Mirrors = append(registry.Mirrors, sysregistriesv2.Endpoint{Location: mirror})
		}
		registries = append(registries, registry)
	}

	// Generate platform-specific bootstrap data
	var platformData platformTemplateData

	switch installConfig.Config.Platform.Name() {
	case baremetaltypes.Name:
		platformData.BareMetal = baremetal.GetTemplateData(installConfig.Config.Platform.BareMetal, installConfig.Config.MachineNetwork, ironicCreds.Username, ironicCreds.Password)
	case vspheretypes.Name:
		platformData.VSphere = vsphere.GetTemplateData(installConfig.Config.Platform.VSphere)
	}

	bootstrapNodeIP := os.Getenv("OPENSHIFT_INSTALL_BOOTSTRAP_NODE_IP")
	if bootstrapNodeIP != "" && net.ParseIP(bootstrapNodeIP) == nil {
		logrus.Warnf("OPENSHIFT_INSTALL_BOOTSTRAP_NODE_IP must have valid ip address, given %s. Skipping it", bootstrapNodeIP)
		bootstrapNodeIP = ""
	}

	platformFirstAPIVIP := firstAPIVIP(&installConfig.Config.Platform)
	APIIntVIPonIPv6 := utilsnet.IsIPv6String(platformFirstAPIVIP)

	networkStack := 0
	for _, snet := range installConfig.Config.ServiceNetwork {
		if snet.IP.To4() != nil {
			networkStack |= 1
		} else {
			networkStack |= 2
		}
	}

	// Set cluster profile
	clusterProfile := ""
	if cp := os.Getenv("OPENSHIFT_INSTALL_EXPERIMENTAL_CLUSTER_PROFILE"); cp != "" {
		logrus.Warnf("Found override for Cluster Profile: %q", cp)
		clusterProfile = cp
	}
	var bootstrapInPlaceConfig *types.BootstrapInPlace
	if bootstrapInPlace {
		bootstrapInPlaceConfig = installConfig.Config.BootstrapInPlace
	}

	apiURL := fmt.Sprintf("api.%s", installConfig.Config.ClusterDomain())
	apiIntURL := fmt.Sprintf("api-int.%s", installConfig.Config.ClusterDomain())
	return &bootstrapTemplateData{
		AdditionalTrustBundle: installConfig.Config.AdditionalTrustBundle,
		FIPS:                  installConfig.Config.FIPS,
		PullSecret:            installConfig.Config.PullSecret,
		SSHKey:                installConfig.Config.SSHKey,
		ReleaseImage:          releaseImage.PullSpec,
		EtcdCluster:           strings.Join(etcdEndpoints, ","),
		Proxy:                 &proxy.Config.Status,
		Registries:            registries,
		BootImage:             string(*rhcosImage),
		PlatformData:          platformData,
		ClusterProfile:        clusterProfile,
		BootstrapInPlace:      bootstrapInPlaceConfig,
		UseIPv6ForNodeIP:      APIIntVIPonIPv6,
		UseDualForNodeIP:      networkStack == 3,
		IsFCOS:                installConfig.Config.IsFCOS(),
		IsSCOS:                installConfig.Config.IsSCOS(),
		IsOKD:                 installConfig.Config.IsOKD(),
		BootstrapNodeIP:       bootstrapNodeIP,
		APIServerURL:          apiURL,
		APIIntServerURL:       apiIntURL,
		FeatureSet:            installConfig.Config.FeatureSet,
	}
}

// AddStorageFiles adds files to a Ignition config.
// Parameters:
// config - the ignition config to be modified
// base - path were the files are written to in to config
// uri - path under data/data specifying the files to be included
// templateData - struct to used to render templates
func AddStorageFiles(config *igntypes.Config, base string, uri string, templateData interface{}) (err error) {
	file, err := data.Assets.Open(uri)
	if err != nil {
		return err
	}
	defer file.Close()

	info, err := file.Stat()
	if err != nil {
		return err
	}

	if info.IsDir() {
		children, err := file.Readdir(0)
		if err != nil {
			return err
		}
		if err = file.Close(); err != nil {
			return err
		}

		for _, childInfo := range children {
			name := childInfo.Name()
			err = AddStorageFiles(config, path.Join(base, name), path.Join(uri, name), templateData)
			if err != nil {
				return err
			}
		}
		return nil
	}

	name := info.Name()
	_, data, err := readFile(name, file, templateData)
	if err != nil {
		return err
	}

	filename := path.Base(uri)
	parentDir := path.Base(path.Dir(uri))

	var mode int
	appendToFile := false
	if parentDir == "bin" || parentDir == "dispatcher.d" {
		mode = 0555
	} else if filename == "motd" || filename == "containers.conf" {
		mode = 0644
		appendToFile = true
	} else {
		mode = 0600
	}
	ign := ignition.FileFromBytes(strings.TrimSuffix(base, ".template"), "root", mode, data)
	if appendToFile {
		ignition.ConvertToAppendix(&ign)
	}

	// Replace files that already exist in the slice with ones added later, otherwise append them
	config.Storage.Files = replaceOrAppend(config.Storage.Files, ign)

	return nil
}

// AddSystemdUnits adds systemd units to a Ignition config.
// Parameters:
// config - the ignition config to be modified
// uri - path under data/data specifying the systemd units files to be included
// templateData - struct to used to render templates
// enabledServices - a list of systemd units to be enabled by default
func AddSystemdUnits(config *igntypes.Config, uri string, templateData interface{}, enabledServices []string) (err error) {
	enabled := make(map[string]struct{}, len(enabledServices))
	for _, s := range enabledServices {
		enabled[s] = struct{}{}
	}

	directory, err := data.Assets.Open(uri)
	if err != nil {
		return err
	}
	defer directory.Close()

	children, err := directory.Readdir(0)
	if err != nil {
		return err
	}

	for _, childInfo := range children {
		dir := path.Join(uri, childInfo.Name())
		file, err := data.Assets.Open(dir)
		if err != nil {
			return err
		}
		defer file.Close()

		info, err := file.Stat()
		if err != nil {
			return err
		}

		if info.IsDir() {
			if dir := info.Name(); !strings.HasSuffix(dir, ".d") {
				logrus.Tracef("Ignoring internal asset directory %q while looking for systemd drop-ins", dir)
				continue
			}

			children, err := file.Readdir(0)
			if err != nil {
				return err
			}
			if err = file.Close(); err != nil {
				return err
			}

			dropins := []igntypes.Dropin{}
			for _, childInfo := range children {
				file, err := data.Assets.Open(path.Join(dir, childInfo.Name()))
				if err != nil {
					return err
				}
				defer file.Close()

				childName, contents, err := readFile(childInfo.Name(), file, templateData)
				if err != nil {
					return err
				}

				dropins = append(dropins, igntypes.Dropin{
					Name:     childName,
					Contents: ignutil.StrToPtr(string(contents)),
				})
			}

			name := strings.TrimSuffix(childInfo.Name(), ".d")
			unit := igntypes.Unit{
				Name:    name,
				Dropins: dropins,
			}
			if _, ok := enabled[name]; ok {
				unit.Enabled = ignutil.BoolToPtr(true)
			}
			config.Systemd.Units = append(config.Systemd.Units, unit)
		} else {
			name, contents, err := readFile(childInfo.Name(), file, templateData)
			if err != nil {
				return err
			}

			unit := igntypes.Unit{
				Name:     name,
				Contents: ignutil.StrToPtr(string(contents)),
			}
			if _, ok := enabled[name]; ok {
				unit.Enabled = ignutil.BoolToPtr(true)
			}
			config.Systemd.Units = append(config.Systemd.Units, unit)
		}
	}

	return nil
}

// Read data from the string reader, and, if the name ends with
// '.template', strip that extension from the name and render the
// template.
func readFile(name string, reader io.Reader, templateData interface{}) (finalName string, data []byte, err error) {
	data, err = io.ReadAll(reader)
	if err != nil {
		return name, []byte{}, err
	}

	if filepath.Ext(name) == ".template" {
		name = strings.TrimSuffix(name, ".template")
		tmpl := template.New(name)
		tmpl, err := tmpl.Parse(string(data))
		if err != nil {
			return name, data, err
		}
		stringData := applyTemplateData(tmpl, templateData)
		data = []byte(stringData)
	}

	return name, data, nil
}

func (a *Common) addParentFiles(dependencies asset.Parents) {
	// These files are all added with mode 0644, i.e. readable
	// by all processes on the system.
	for _, asset := range []asset.WritableAsset{
		&manifests.Manifests{},
		&manifests.Openshift{},
		&machines.Master{},
		&machines.Worker{},
		&mcign.MasterIgnitionCustomizations{},
		&mcign.WorkerIgnitionCustomizations{},
		&CVOIgnore{}, // this must come after manifests.Manifests so that it replaces cvo-overrides.yaml
	} {
		dependencies.Get(asset)

		// Replace files that already exist in the slice with ones added later, otherwise append them
		for _, file := range ignition.FilesFromAsset(rootDir, "root", 0644, asset) {
			a.Config.Storage.Files = replaceOrAppend(a.Config.Storage.Files, file)
		}
	}

	// These files are all added with mode 0600; use for secret keys and the like.
	for _, asset := range []asset.WritableAsset{
		&kubeconfig.AdminInternalClient{},
		&kubeconfig.Kubelet{},
		&kubeconfig.LoopbackClient{},
		&tls.AdminKubeConfigCABundle{},
		&tls.AggregatorCA{},
		&tls.AggregatorCABundle{},
		&tls.AggregatorClientCertKey{},
		&tls.AggregatorSignerCertKey{},
		&tls.APIServerProxyCertKey{},
		&tls.BoundSASigningKey{},
		&tls.CloudProviderCABundle{},
		&tls.KubeAPIServerLBCABundle{},
		&tls.KubeAPIServerExternalLBServerCertKey{},
		&tls.KubeAPIServerInternalLBServerCertKey{},
		&tls.KubeAPIServerLBSignerCertKey{},
		&tls.KubeAPIServerLocalhostCABundle{},
		&tls.KubeAPIServerLocalhostServerCertKey{},
		&tls.KubeAPIServerLocalhostSignerCertKey{},
		&tls.KubeAPIServerServiceNetworkCABundle{},
		&tls.KubeAPIServerServiceNetworkServerCertKey{},
		&tls.KubeAPIServerServiceNetworkSignerCertKey{},
		&tls.KubeAPIServerCompleteCABundle{},
		&tls.KubeAPIServerCompleteClientCABundle{},
		&tls.KubeAPIServerToKubeletCABundle{},
		&tls.KubeAPIServerToKubeletClientCertKey{},
		&tls.KubeAPIServerToKubeletSignerCertKey{},
		&tls.KubeControlPlaneCABundle{},
		&tls.KubeControlPlaneKubeControllerManagerClientCertKey{},
		&tls.KubeControlPlaneKubeSchedulerClientCertKey{},
		&tls.KubeControlPlaneSignerCertKey{},
		&tls.KubeletBootstrapCABundle{},
		&tls.KubeletClientCABundle{},
		&tls.KubeletClientCertKey{},
		&tls.KubeletCSRSignerCertKey{},
		&tls.KubeletServingCABundle{},
		&tls.MCSCertKey{},
		&tls.ServiceAccountKeyPair{},
		&tls.JournalCertKey{},
	} {
		dependencies.Get(asset)

		// Replace files that already exist in the slice with ones added later, otherwise append them
		for _, file := range ignition.FilesFromAsset(rootDir, "root", 0600, asset) {
			a.Config.Storage.Files = replaceOrAppend(a.Config.Storage.Files, file)
		}
	}

	rootCA := &tls.RootCA{}
	dependencies.Get(rootCA)
	a.Config.Storage.Files = replaceOrAppend(a.Config.Storage.Files, ignition.FileFromBytes(filepath.Join(rootDir, rootCA.CertFile().Filename), "root", 0644, rootCA.Cert()))
}

func replaceOrAppend(files []igntypes.File, file igntypes.File) []igntypes.File {
	for i, f := range files {
		if f.Node.Path == file.Node.Path {
			files[i] = file
			return files
		}
	}
	files = append(files, file)
	return files
}

func applyTemplateData(template *template.Template, templateData interface{}) string {
	buf := &bytes.Buffer{}
	if err := template.Execute(buf, templateData); err != nil {
		panic(err)
	}
	return buf.String()
}

// Load returns the bootstrap ignition from disk.
func (a *Common) load(f asset.FileFetcher, filename string) (found bool, err error) {
	file, err := f.FetchByName(filename)
	if err != nil {
		if os.IsNotExist(err) {
			return false, nil
		}
		return false, err
	}

	config := &igntypes.Config{}
	if err := json.Unmarshal(file.Data, config); err != nil {
		return false, errors.Wrapf(err, "failed to unmarshal %s", filename)
	}

	a.File, a.Config = file, config
	warnIfCertificatesExpired(a.Config)
	return true, nil
}

// warnIfCertificatesExpired checks for expired certificates and warns if so
func warnIfCertificatesExpired(config *igntypes.Config) {
	expiredCerts := 0
	for _, file := range config.Storage.Files {
		if filepath.Ext(file.Path) == ".crt" && file.Contents.Source != nil {
			fileName := path.Base(file.Path)
			decoded, err := dataurl.DecodeString(*file.Contents.Source)
			if err != nil {
				logrus.Debugf("Unable to decode certificate %s: %s", fileName, err.Error())
				continue
			}
			data := decoded.Data
			for {
				block, rest := pem.Decode(data)
				if block == nil {
					break
				}

				cert, err := x509.ParseCertificate(block.Bytes)
				if err == nil {
					if time.Now().UTC().After(cert.NotAfter) {
						logrus.Warnf("Bootstrap Ignition-Config Certificate %s expired at %s.", path.Base(file.Path), cert.NotAfter.Format(time.RFC3339))
						expiredCerts++
					}
				} else {
					logrus.Debugf("Unable to parse certificate %s: %s", fileName, err.Error())
					break
				}

				data = rest
			}
		}
	}

	if expiredCerts > 0 {
		logrus.Warnf("Bootstrap Ignition-Config: %d certificates expired. Installation attempts with the created Ignition-Configs will possibly fail.", expiredCerts)
	}
}

// APIVIPs returns the string representations of the platform's API VIPs
// It returns nil if the platform does not configure VIPs
func apiVIPs(p *types.Platform) []string {
	switch {
	case p == nil:
		return nil
	case p.BareMetal != nil:
		return p.BareMetal.APIVIPs
	case p.OpenStack != nil:
		return p.OpenStack.APIVIPs
	case p.VSphere != nil:
		return p.VSphere.APIVIPs
	case p.Ovirt != nil:
		return p.Ovirt.APIVIPs
	case p.Nutanix != nil:
		return p.Nutanix.APIVIPs
	default:
		return nil
	}
}

// firstAPIVIP returns the first VIP of the API server (e.g. in case of
// dual-stack)
func firstAPIVIP(p *types.Platform) string {
	for _, vip := range apiVIPs(p) {
		return vip
	}

	return ""
}
