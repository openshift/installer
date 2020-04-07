package manifests

import (
	"archive/tar"
        "bytes"
	"compress/gzip"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"text/template"

	"github.com/ghodss/yaml"
	"github.com/pkg/errors"

	configv1 "github.com/openshift/api/config/v1"
	"github.com/openshift/installer/pkg/asset"
	"github.com/openshift/installer/pkg/asset/installconfig"
	"github.com/openshift/installer/pkg/asset/templates/content/openshift"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

var (
	noCrdFilename = filepath.Join(manifestDir, "cluster-network-01-crd.yml")
	noCfgFilename = filepath.Join(manifestDir, "cluster-network-02-config.yml")
	noSNATCRFilename = filepath.Join(manifestDir, "cluster-network-27-snat-policy-cr.yaml")
)

var snatCRTmpl = template.Must(template.New("snat-cr").Parse(`apiVersion: aci.snat/v1
kind: SnatPolicy
metadata:
  name: installerclusterdefault
spec:
  snatIp:
    -  {{.snatIP}}
`))

var rdConfigTmpl = template.Must(template.New("rdconfig").Parse(`apiVersion: aci.snat/v1
kind: RdConfig
metadata:
  name: routingdomain-config
  namespace: aci-containers-system
spec:
  usersubnets:
  - {{ .neutronCIDR }}
  - 224.0.0.0/4
`))

var clusterNetwork03Tmpl = template.Must(template.New("cluster03").Parse(`apiVersion: operator.openshift.io/v1
kind: Network
metadata:
  name: cluster
spec:
  disableMultiNetwork: true
  clusterNetwork:
  - cidr: {{.clusterNet}}
    hostPrefix: {{.hostPrefix}}
  defaultNetwork:
    type: {{.netType}}
  serviceNetwork:
  - {{.svcNet}}
`))

// We need to manually create our CRDs first, so we can create the
// configuration instance of it in the installer. Other operators have
// their CRD created by the CVO, but we need to create the corresponding
// CRs in the installer, so we need the CRD to be there.
// The first CRD is the high-level Network.config.openshift.io object,
// which is stable and minimal. Administrators can configure the
// network in a more detailed manner with the operator-specific CR, which
// also needs to be done before the installer is run, so we provide both.

// Networking generates the cluster-network-*.yml files.
type Networking struct {
	Config   *configv1.Network
	FileList []*asset.File
}

var _ asset.WritableAsset = (*Networking)(nil)

// Name returns a human friendly name for the operator.
func (no *Networking) Name() string {
	return "Network Config"
}

// Dependencies returns all of the dependencies directly needed to generate
// network configuration.
func (no *Networking) Dependencies() []asset.Asset {
	return []asset.Asset{
		&installconfig.InstallConfig{},
		&openshift.NetworkCRDs{},
	}
}

// Generate generates the network operator config and its CRD and the SNAT Cluster CR if needed.
func (no *Networking) Generate(dependencies asset.Parents) error {
	installConfig := &installconfig.InstallConfig{}
	crds := &openshift.NetworkCRDs{}
	dependencies.Get(installConfig, crds)

	netConfig := installConfig.Config.Networking

	clusterNet := []configv1.ClusterNetworkEntry{}
	if len(netConfig.ClusterNetwork) > 0 {
		for _, net := range netConfig.ClusterNetwork {
			clusterNet = append(clusterNet, configv1.ClusterNetworkEntry{
				CIDR:       net.CIDR.String(),
				HostPrefix: uint32(net.HostPrefix),
			})
		}
	} else {
		return errors.Errorf("ClusterNetworks must be specified")
	}

	serviceNet := []string{}
	for _, sn := range netConfig.ServiceNetwork {
		serviceNet = append(serviceNet, sn.String())
	}

	no.Config = &configv1.Network{
		TypeMeta: metav1.TypeMeta{
			APIVersion: configv1.SchemeGroupVersion.String(),
			Kind:       "Network",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name: "cluster",
			// not namespaced
		},
		Spec: configv1.NetworkSpec{
			ClusterNetwork: clusterNet,
			ServiceNetwork: serviceNet,
			NetworkType:    netConfig.NetworkType,
			// Block all Service.ExternalIPs by default
			ExternalIP: &configv1.ExternalIPConfig{
				Policy: &configv1.ExternalIPPolicy{},
			},
		},
		Status: configv1.NetworkStatus{
			ClusterNetwork: clusterNet,
			ServiceNetwork: serviceNet,
			NetworkType:    netConfig.NetworkType,
		},
	}

	configData, err := yaml.Marshal(no.Config)
	if err != nil {
		return errors.Wrapf(err, "failed to create %s manifests from InstallConfig", no.Name())
	}

	crdContents := ""
	for _, crdFile := range crds.Files() {
		crdContents = fmt.Sprintf("%s\n---\n%s", crdContents, crdFile.Data)
	}

	no.FileList = []*asset.File{
                {
                        Filename: noCrdFilename,
                        Data:     []byte(crdContents),
                },
                {
                        Filename: noCfgFilename,
                        Data:     configData,
                },
        }

	// Untar and add acc-provision files
	r, _ := os.Open(installConfig.Config.Platform.OpenStack.AciNetExt.ProvisionTar)
	uncompressedStream, _ := gzip.NewReader(r)
	tarReader := tar.NewReader(uncompressedStream)
	var noRDconfigFilename string
	for true {
		header, err := tarReader.Next()
		if err == io.EOF {
			break
		}
		b, _ := ioutil.ReadAll(tarReader);

		// Save filenumber of hostagent daemonset for rdconfig CR
		if strings.Contains(header.Name, "DaemonSet-aci-containers-host"){
			hostAgentFileName := header.Name
			hyphenParsed := strings.Split(hostAgentFileName, "-")
			hostAgentFileNo, err := strconv.Atoi(hyphenParsed[2])
			if err != nil {
				return errors.Wrapf(err, "failed to decipher host agent  manifest file from acc-provision tar")
			}
			rdConfigFileNo := strconv.Itoa(hostAgentFileNo - 4)
			noRDconfigFilename = filepath.Join(manifestDir, "cluster-network-" + rdConfigFileNo + "-2-CustomResource-rdconfig.yaml")
		}

		// Edit cluster-network-03 file with correct fields
		if strings.Contains(header.Name, "cluster-network-03"){
			cluster03Data := &bytes.Buffer{}
			clusterNetworkCIDR := &netConfig.ClusterNetwork[0].CIDR
			data := map[string]string{"clusterNet": clusterNetworkCIDR.String(),
                		"hostPrefix":  strconv.Itoa(int(netConfig.ClusterNetwork[0].HostPrefix)),
                		"netType": netConfig.NetworkType, "svcNet": netConfig.ServiceNetwork[0].String()}
			if err := clusterNetwork03Tmpl.Execute(cluster03Data, data); err != nil {
				return errors.Wrapf(err, "failed to create cluster-network-03 manifests from InstallConfig")
			}
			b = cluster03Data.Bytes()

		}
		tempFile := &asset.File{Filename: filepath.Join(manifestDir, header.Name), Data: b}
		no.FileList = append(no.FileList, tempFile)
	}

	// Create SNAT Cluster CR file 
	if installConfig.Config.Platform.OpenStack.AciNetExt.ClusterSNATSubnet != "" {
		snatData := &bytes.Buffer{}
		data := map[string]string{"snatIP": installConfig.Config.Platform.OpenStack.AciNetExt.ClusterSNATSubnet}
		if err := snatCRTmpl.Execute(snatData, data); err != nil {
			return errors.Wrapf(err, "failed to create SNAT CR manifests from InstallConfig")
		}
		// add destIP if field present
		if installConfig.Config.Platform.OpenStack.AciNetExt.ClusterSNATDest != "" {
			dest := "  destIp:\n    -  " + installConfig.Config.Platform.OpenStack.AciNetExt.ClusterSNATDest + "\n"
			snatData.WriteString(dest)
		}
		snatFile := &asset.File{Filename: noSNATCRFilename, Data: snatData.Bytes()}
		no.FileList = append(no.FileList, snatFile)

		// Create yaml for rdConfig
		if noRDconfigFilename == "" {
			return errors.New("no manifest with DaemonSet-aci-containers-host found in acc-provision tar")
		}
                rdConfigData := &bytes.Buffer{}
                data = map[string]string{"neutronCIDR": installConfig.Config.Platform.OpenStack.AciNetExt.NeutronCIDR.String()}
                if err = rdConfigTmpl.Execute(rdConfigData, data); err != nil {
                        return errors.Wrapf(err, "failed to create rdconfig manifest from InstallConfig")
                }
                rdconfigFile := &asset.File{Filename: noRDconfigFilename, Data: rdConfigData.Bytes()}
                no.FileList = append(no.FileList, rdconfigFile)
	}

	return nil
}

// Files returns the files generated by the asset.
func (no *Networking) Files() []*asset.File {
	return no.FileList
}

// Load returns false since this asset is not written to disk by the installer.
func (no *Networking) Load(f asset.FileFetcher) (bool, error) {
	return false, nil
}
