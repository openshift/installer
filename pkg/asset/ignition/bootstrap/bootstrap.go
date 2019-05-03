package bootstrap

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"strings"
	"text/template"

	"github.com/containers/image/pkg/sysregistriesv2"
	"github.com/coreos/ignition/config/util"
	igntypes "github.com/coreos/ignition/config/v2_2/types"
	configv1 "github.com/openshift/api/config/v1"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"

	"github.com/openshift/installer/data"
	"github.com/openshift/installer/pkg/asset"
	"github.com/openshift/installer/pkg/asset/ignition"
	"github.com/openshift/installer/pkg/asset/installconfig"
	"github.com/openshift/installer/pkg/asset/kubeconfig"
	"github.com/openshift/installer/pkg/asset/machines"
	"github.com/openshift/installer/pkg/asset/manifests"
	"github.com/openshift/installer/pkg/asset/releaseimage"
	"github.com/openshift/installer/pkg/asset/tls"
	"github.com/openshift/installer/pkg/types"
)

const (
	rootDir              = "/opt/openshift"
	bootstrapIgnFilename = "bootstrap.ign"
	ignitionUser         = "core"
)

// bootstrapTemplateData is the data to use to replace values in bootstrap
// template files.
type bootstrapTemplateData struct {
	AdditionalTrustBundle string
	EtcdCluster           string
	PullSecret            string
	ReleaseImage          string
	Proxy                 *configv1.ProxyStatus
	Registries            []sysregistriesv2.Registry
}

// Bootstrap is an asset that generates the ignition config for bootstrap nodes.
type Bootstrap struct {
	Config *igntypes.Config
	File   *asset.File
}

var _ asset.WritableAsset = (*Bootstrap)(nil)

// Dependencies returns the assets on which the Bootstrap asset depends.
func (a *Bootstrap) Dependencies() []asset.Asset {
	return []asset.Asset{
		&installconfig.InstallConfig{},
		&kubeconfig.AdminClient{},
		&kubeconfig.Kubelet{},
		&kubeconfig.LoopbackClient{},
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
	}
}

// Generate generates the ignition config for the Bootstrap asset.
func (a *Bootstrap) Generate(dependencies asset.Parents) error {
	installConfig := &installconfig.InstallConfig{}
	proxy := &manifests.Proxy{}
	releaseImage := &releaseimage.Image{}
	dependencies.Get(installConfig, proxy, releaseImage)

	templateData, err := a.getTemplateData(installConfig.Config, releaseImage.PullSpec, installConfig.Config.ImageContentSources, proxy.Config)

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
		igntypes.PasswdUser{Name: "core", SSHAuthorizedKeys: []igntypes.SSHAuthorizedKey{igntypes.SSHAuthorizedKey(installConfig.Config.SSHKey)}},
	)

	data, err := json.Marshal(a.Config)
	if err != nil {
		return errors.Wrap(err, "failed to Marshal Ignition config")
	}
	a.File = &asset.File{
		Filename: bootstrapIgnFilename,
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
func (a *Bootstrap) getTemplateData(installConfig *types.InstallConfig, releaseImage string, imageSources []types.ImageContentSource, proxy *configv1.Proxy) (*bootstrapTemplateData, error) {
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

	return &bootstrapTemplateData{
		AdditionalTrustBundle: installConfig.AdditionalTrustBundle,
		PullSecret:            installConfig.PullSecret,
		ReleaseImage:          releaseImage,
		EtcdCluster:           strings.Join(etcdEndpoints, ","),
		Proxy:                 &proxy.Status,
		Registries:            registries,
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

	var mode int
	appendToFile := false
	if path.Base(path.Dir(uri)) == "bin" {
		mode = 0555
	} else if filename == "motd" {
		mode = 0644
		appendToFile = true
	} else {
		mode = 0600
	}
	ign := ignition.FileFromBytes(strings.TrimSuffix(base, ".template"), "root", mode, data)
	ign.Append = appendToFile
	a.Config.Storage.Files = append(a.Config.Storage.Files, ign)

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
		"keepalived.service": {},
		"coredns.service":    {},
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

			dropins := []igntypes.SystemdDropin{}
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

				dropins = append(dropins, igntypes.SystemdDropin{
					Name:     childName,
					Contents: string(contents),
				})
			}

			name := strings.TrimSuffix(childInfo.Name(), ".d")
			unit := igntypes.Unit{
				Name:    name,
				Dropins: dropins,
			}
			if _, ok := enabled[name]; ok {
				unit.Enabled = util.BoolToPtr(true)
			}
			a.Config.Systemd.Units = append(a.Config.Systemd.Units, unit)
		} else {
			name, contents, err := readFile(childInfo.Name(), file, templateData)
			if err != nil {
				return err
			}

			unit := igntypes.Unit{
				Name:     name,
				Contents: string(contents),
			}
			if _, ok := enabled[name]; ok {
				unit.Enabled = util.BoolToPtr(true)
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
	} {
		dependencies.Get(asset)
		a.Config.Storage.Files = append(a.Config.Storage.Files, ignition.FilesFromAsset(rootDir, "root", 0644, asset)...)
	}

	// These files are all added with mode 0600; use for secret keys and the like.
	for _, asset := range []asset.WritableAsset{
		&kubeconfig.AdminClient{},
		&kubeconfig.Kubelet{},
		&kubeconfig.LoopbackClient{},
		&tls.AdminKubeConfigCABundle{},
		&tls.AggregatorCA{},
		&tls.AggregatorCABundle{},
		&tls.AggregatorClientCertKey{},
		&tls.AggregatorSignerCertKey{},
		&tls.APIServerProxyCertKey{},
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
		a.Config.Storage.Files = append(a.Config.Storage.Files, ignition.FilesFromAsset(rootDir, "root", 0600, asset)...)
	}

	rootCA := &tls.RootCA{}
	dependencies.Get(rootCA)
	a.Config.Storage.Files = append(a.Config.Storage.Files, ignition.FileFromBytes(filepath.Join(rootDir, rootCA.CertFile().Filename), "root", 0644, rootCA.Cert()))
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
	return true, nil
}
