package bootstrap

import (
	"bytes"
	"crypto/x509"
	"encoding/json"
	"encoding/pem"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"strings"
	"text/template"
	"time"

	"github.com/containers/image/pkg/sysregistriesv2"
	ignutil "github.com/coreos/ignition/v2/config/util"
	igntypes "github.com/coreos/ignition/v2/config/v3_1/types"
	configv1 "github.com/openshift/api/config/v1"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"github.com/vincent-petithory/dataurl"

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
	rootDir                     = "/opt/openshift"
	bootstrapIgnFilename        = "bootstrap.ign"
	bootstrapInPlaceIgnFilename = "bootstrap-in-place-for-live-iso.ign"
)

// bootstrapTemplateData is the data to use to replace values in bootstrap
// template files.
type bootstrapTemplateData struct {
	AdditionalTrustBundle string
	FIPS                  bool
	EtcdCluster           string
	PullSecret            string
	ReleaseImage          string
	ClusterProfile        string
	Proxy                 *configv1.ProxyStatus
	Registries            []sysregistriesv2.Registry
	BootImage             string
	PlatformData          platformTemplateData
	SingleNode            *SingleNodeBootstrapInPlace
}

// platformTemplateData is the data to use to replace values in bootstrap
// template files that are specific to one platform.
type platformTemplateData struct {
	BareMetal *baremetal.TemplateData
	VSphere   *vsphere.TemplateData
}

// Bootstrap is an asset that generates the ignition config for bootstrap nodes.
type Bootstrap struct {
	SingleNodeBootstrapInPlace bool
	Config                     *igntypes.Config
	File                       *asset.File
}

var _ asset.WritableAsset = (*Bootstrap)(nil)

// Dependencies returns the assets on which the Bootstrap asset depends.
func (a *Bootstrap) Dependencies() []asset.Asset {
	return []asset.Asset{
		&baremetal.IronicCreds{},
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
		&tls.EtcdCABundle{},
		&tls.EtcdMetricCABundle{},
		&tls.EtcdMetricSignerCertKey{},
		&tls.EtcdMetricSignerClientCertKey{},
		&tls.EtcdSignerCertKey{},
		&tls.EtcdSignerClientCertKey{},
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
func (a *Bootstrap) Generate(dependencies asset.Parents) error {
	installConfig := &installconfig.InstallConfig{}
	proxy := &manifests.Proxy{}
	releaseImage := &releaseimage.Image{}
	rhcosImage := new(rhcos.Image)
	bootstrapSSHKeyPair := &tls.BootstrapSSHKeyPair{}
	ironicCreds := &baremetal.IronicCreds{}
	dependencies.Get(installConfig, proxy, releaseImage, rhcosImage, bootstrapSSHKeyPair, ironicCreds)

	templateData, err := a.getTemplateData(installConfig.Config, releaseImage.PullSpec, installConfig.Config.ImageContentSources, proxy.Config, rhcosImage, ironicCreds)

	if err != nil {
		return errors.Wrap(err, "failed to get bootstrap templates")
	}

	a.Config = &igntypes.Config{
		Ignition: igntypes.Ignition{
			Version: igntypes.MaxVersion.String(),
		},
	}

	err = a.addStorageFiles("/", "bootstrap/files", templateData)
	if err != nil {
		return err
	}
	if templateData.SingleNode.BootstrapInPlace {
		err = a.addStorageFiles("/", "bootstrap/bootstrap-in-place", templateData)
		if err != nil {
			return err
		}
	}
	err = a.addSystemdUnits("bootstrap/systemd/units", templateData)
	if err != nil {
		return err
	}

	// Check for optional platform specific files/units
	platform := installConfig.Config.Platform.Name()
	platformFilePath := fmt.Sprintf("bootstrap/%s/files", platform)
	directory, err := data.Assets.Open(platformFilePath)
	if err == nil {
		directory.Close()
		err = a.addStorageFiles("/", platformFilePath, templateData)
		if err != nil {
			return err
		}
	}

	platformUnitPath := fmt.Sprintf("bootstrap/%s/systemd/units", platform)
	directory, err = data.Assets.Open(platformUnitPath)
	if err == nil {
		directory.Close()
		err = a.addSystemdUnits(platformUnitPath, templateData)
		if err != nil {
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

	data, err := ignition.Marshal(a.Config)
	if err != nil {
		return errors.Wrap(err, "failed to Marshal Ignition config")
	}
	fileName := bootstrapIgnFilename
	if a.SingleNodeBootstrapInPlace {
		if err := verifyBootstrapInPlace(installConfig.Config); err != nil {
			return err
		}
		fileName = bootstrapInPlaceIgnFilename
	}
	a.File = &asset.File{
		Filename: fileName,
		Data:     data,
	}

	return nil
}

// Name returns the human-friendly name of the asset.
func (a *Bootstrap) Name() string {
	return "Bootstrap Ignition Config"
}

// Files returns the files generated by the asset.
func (a *Bootstrap) Files() []*asset.File {
	if a.File != nil {
		return []*asset.File{a.File}
	}
	return []*asset.File{}
}

// getTemplateData returns the data to use to execute bootstrap templates.
func (a *Bootstrap) getTemplateData(installConfig *types.InstallConfig, releaseImage string, imageSources []types.ImageContentSource, proxy *configv1.Proxy, rhcosImage *rhcos.Image, ironicCreds *baremetal.IronicCreds) (*bootstrapTemplateData, error) {
	etcdEndpoints := make([]string, *installConfig.ControlPlane.Replicas)

	for i := range etcdEndpoints {
		etcdEndpoints[i] = fmt.Sprintf("https://etcd-%d.%s:2379", i, installConfig.ClusterDomain())
	}

	registries := []sysregistriesv2.Registry{}
	for _, group := range mergedMirrorSets(imageSources) {
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

	switch installConfig.Platform.Name() {
	case baremetaltypes.Name:
		platformData.BareMetal = baremetal.GetTemplateData(installConfig.Platform.BareMetal, installConfig.MachineNetwork, ironicCreds.Username, ironicCreds.Password)
	case vspheretypes.Name:
		platformData.VSphere = vsphere.GetTemplateData(installConfig.Platform.VSphere)
	}

	// Set cluster profile
	clusterProfile := ""
	if cp := os.Getenv("OPENSHIFT_INSTALL_EXPERIMENTAL_CLUSTER_PROFILE"); cp != "" {
		logrus.Warnf("Found override for Cluster Profile: %q", cp)
		clusterProfile = cp
	}
	var singelNodeTemplateData *SingleNodeBootstrapInPlace
	if a.SingleNodeBootstrapInPlace {
		singelNodeTemplateData = GetSingleNodeBootstrapInPlaceConfig()
	} else {
		singelNodeTemplateData = &SingleNodeBootstrapInPlace{BootstrapInPlace: false}
	}

	return &bootstrapTemplateData{
		AdditionalTrustBundle: installConfig.AdditionalTrustBundle,
		FIPS:                  installConfig.FIPS,
		PullSecret:            installConfig.PullSecret,
		ReleaseImage:          releaseImage,
		EtcdCluster:           strings.Join(etcdEndpoints, ","),
		Proxy:                 &proxy.Status,
		Registries:            registries,
		BootImage:             string(*rhcosImage),
		PlatformData:          platformData,
		ClusterProfile:        clusterProfile,
		SingleNode:            singelNodeTemplateData,
	}, nil
}

func (a *Bootstrap) addStorageFiles(base string, uri string, templateData *bootstrapTemplateData) (err error) {
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
			err = a.addStorageFiles(path.Join(base, name), path.Join(uri, name), templateData)
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
	} else if filename == "motd" {
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
	a.Config.Storage.Files = replaceOrAppend(a.Config.Storage.Files, ign)

	return nil
}

func (a *Bootstrap) addSystemdUnits(uri string, templateData *bootstrapTemplateData) (err error) {
	enabled := map[string]struct{}{
		"progress.service":                {},
		"kubelet.service":                 {},
		"chown-gatewayd-key.service":      {},
		"systemd-journal-gatewayd.socket": {},
		"approve-csr.service":             {},
		// baremetal & openstack platform services
		"keepalived.service":        {},
		"coredns.service":           {},
		"ironic.service":            {},
		"master-bmh-update.service": {},
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
			a.Config.Systemd.Units = append(a.Config.Systemd.Units, unit)
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
			a.Config.Systemd.Units = append(a.Config.Systemd.Units, unit)
		}
	}

	return nil
}

// Read data from the string reader, and, if the name ends with
// '.template', strip that extension from the name and render the
// template.
func readFile(name string, reader io.Reader, templateData interface{}) (finalName string, data []byte, err error) {
	data, err = ioutil.ReadAll(reader)
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

func (a *Bootstrap) addParentFiles(dependencies asset.Parents) {
	// These files are all added with mode 0644, i.e. readable
	// by all processes on the system.
	for _, asset := range []asset.WritableAsset{
		&manifests.Manifests{},
		&manifests.Openshift{},
		&machines.Master{},
		&machines.Worker{},
		&mcign.MasterIgnitionCustomizations{},
		&mcign.WorkerIgnitionCustomizations{},
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
		&tls.EtcdCABundle{},
		&tls.EtcdMetricCABundle{},
		&tls.EtcdMetricSignerCertKey{},
		&tls.EtcdMetricSignerClientCertKey{},
		&tls.EtcdSignerCertKey{},
		&tls.EtcdSignerClientCertKey{},
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
func (a *Bootstrap) Load(f asset.FileFetcher) (found bool, err error) {
	file, err := f.FetchByName(bootstrapIgnFilename)
	if err != nil {
		if os.IsNotExist(err) {
			return false, nil
		}
		return false, err
	}

	config := &igntypes.Config{}
	if err := json.Unmarshal(file.Data, config); err != nil {
		return false, errors.Wrapf(err, "failed to unmarshal %s", bootstrapIgnFilename)
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
