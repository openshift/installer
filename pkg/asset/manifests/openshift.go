package manifests

import (
	"encoding/base64"
	"path/filepath"

	"github.com/aws/aws-sdk-go/aws/session"

	// TODO(flaper87): Migrate to ghodss asap
	// This yaml is currently used only by the OpenStack
	// clouds serialization. We're working on migrating
	// clientconfig out of go-yaml. We'll use it here
	// until that happens.
	// https://github.com/openshift/installer/pull/854
	goyaml "gopkg.in/yaml.v2"

	"github.com/ghodss/yaml"
	"github.com/gophercloud/utils/openstack/clientconfig"
	machineapi "github.com/openshift/cluster-api/pkg/apis/machine/v1beta1"
	"github.com/openshift/installer/pkg/asset"
	"github.com/openshift/installer/pkg/asset/installconfig"
	"github.com/openshift/installer/pkg/asset/machines"
	osmachine "github.com/openshift/installer/pkg/asset/machines/openstack"
	"github.com/openshift/installer/pkg/asset/password"
	"github.com/openshift/installer/pkg/asset/templates/content/openshift"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	clusterapi "sigs.k8s.io/cluster-api/pkg/apis/cluster/v1alpha1"
)

const (
	openshiftManifestDir = "openshift"
)

var (
	_ asset.WritableAsset = (*Openshift)(nil)
)

// Openshift generates the dependent resource manifests for openShift (as against bootkube)
type Openshift struct {
	FileList []*asset.File
}

// Name returns a human friendly name for the operator
func (o *Openshift) Name() string {
	return "Openshift Manifests"
}

// Dependencies returns all of the dependencies directly needed by the
// Openshift asset
func (o *Openshift) Dependencies() []asset.Asset {
	return []asset.Asset{
		&installconfig.InstallConfig{},
		&ClusterK8sIO{},
		&machines.Worker{},
		&machines.Master{},
		&password.KubeadminPassword{},

		&openshift.BindingDiscovery{},
		&openshift.CloudCredsSecret{},
		&openshift.DeprecatedCloudCredsSecret{},
		&openshift.KubeadminPasswordSecret{},
		&openshift.RoleCloudCredsSecretReader{},
		&openshift.DeprecatedRoleCloudCredsSecretReader{},
	}
}

// Generate generates the respective operator config.yml files
func (o *Openshift) Generate(dependencies asset.Parents) error {
	installConfig := &installconfig.InstallConfig{}
	kubeadminPassword := &password.KubeadminPassword{}
	clusterk8sio := &ClusterK8sIO{}
	worker := &machines.Worker{}
	master := &machines.Master{}
	dependencies.Get(installConfig, clusterk8sio, worker, master, kubeadminPassword)
	var cloudCreds cloudCredsSecretData
	platform := installConfig.Config.Platform.Name()
	switch platform {
	case "aws":
		ssn := session.Must(session.NewSessionWithOptions(session.Options{
			SharedConfigState: session.SharedConfigEnable,
		}))
		creds, err := ssn.Config.Credentials.Get()
		if err != nil {
			return err
		}
		cloudCreds = cloudCredsSecretData{
			AWS: &AwsCredsSecretData{
				Base64encodeAccessKeyID:     base64.StdEncoding.EncodeToString([]byte(creds.AccessKeyID)),
				Base64encodeSecretAccessKey: base64.StdEncoding.EncodeToString([]byte(creds.SecretAccessKey)),
			},
		}
	case "openstack":
		opts := new(clientconfig.ClientOpts)
		opts.Cloud = installConfig.Config.Platform.OpenStack.Cloud
		cloud, err := clientconfig.GetCloudFromYAML(opts)
		if err != nil {
			return err
		}
		clouds := make(map[string]map[string]*clientconfig.Cloud)
		clouds["clouds"] = map[string]*clientconfig.Cloud{
			osmachine.CloudName: cloud,
		}

		marshalled, err := goyaml.Marshal(clouds)
		if err != nil {
			return err
		}

		credsEncoded := base64.StdEncoding.EncodeToString(marshalled)
		cloudCreds = cloudCredsSecretData{
			OpenStack: &OpenStackCredsSecretData{
				Base64encodeCloudCreds: credsEncoded,
			},
		}
	}

	templateData := &openshiftTemplateData{
		CloudCreds:                   cloudCreds,
		Base64EncodedKubeadminPwHash: base64.StdEncoding.EncodeToString(kubeadminPassword.PasswordHash),
	}

	bindingDiscovery := &openshift.BindingDiscovery{}
	deprecatedCloudCredsSecret := &openshift.DeprecatedCloudCredsSecret{}
	cloudCredsSecret := &openshift.CloudCredsSecret{}
	kubeadminPasswordSecret := &openshift.KubeadminPasswordSecret{}
	deprecatedRoleCloudCredsSecretReader := &openshift.DeprecatedRoleCloudCredsSecretReader{}
	roleCloudCredsSecretReader := &openshift.RoleCloudCredsSecretReader{}
	dependencies.Get(
		bindingDiscovery,
		cloudCredsSecret,
		kubeadminPasswordSecret,
		roleCloudCredsSecretReader)

	var masterMachines []byte
	var err error
	if master.Machines != nil {
		masterMachines, err = listFromMachines(master.Machines)
	} else {
		masterMachines, err = listFromMachinesDeprecated(master.MachinesDeprecated)
	}
	if err != nil {
		return err
	}

	assetData := map[string][]byte{
		"99_binding-discovery.yaml":                             []byte(bindingDiscovery.Files()[0].Data),
		"99_kubeadmin-password-secret.yaml":                     applyTemplateData(kubeadminPasswordSecret.Files()[0].Data, templateData),
		"99_openshift-cluster-api_cluster.yaml":                 clusterk8sio.Raw,
		"99_openshift-cluster-api_master-machines.yaml":         masterMachines,
		"99_openshift-cluster-api_master-user-data-secret.yaml": master.UserDataSecretRaw,
		"99_openshift-cluster-api_worker-machineset.yaml":       worker.MachineSetRaw,
		"99_openshift-cluster-api_worker-user-data-secret.yaml": worker.UserDataSecretRaw,
	}

	switch platform {
	case "aws", "openstack":
		assetData["99_cloud-creds-secret.yaml"] = applyTemplateData(cloudCredsSecret.Files()[0].Data, templateData)
		assetData["99_deprecated-cloud-creds-secret.yaml"] = applyTemplateData(deprecatedCloudCredsSecret.Files()[0].Data, templateData)
		assetData["99_role-cloud-creds-secret-reader.yaml"] = applyTemplateData(roleCloudCredsSecretReader.Files()[0].Data, templateData)
		assetData["99_deprecated-role-cloud-creds-secret-reader.yaml"] = applyTemplateData(deprecatedRoleCloudCredsSecretReader.Files()[0].Data, templateData)
	}

	o.FileList = []*asset.File{}
	for name, data := range assetData {
		o.FileList = append(o.FileList, &asset.File{
			Filename: filepath.Join(openshiftManifestDir, name),
			Data:     data,
		})
	}

	asset.SortFiles(o.FileList)

	return nil
}

// Files returns the files generated by the asset.
func (o *Openshift) Files() []*asset.File {
	return o.FileList
}

// Load returns the openshift asset from disk.
func (o *Openshift) Load(f asset.FileFetcher) (bool, error) {
	fileList, err := f.FetchByPattern(filepath.Join(openshiftManifestDir, "*"))
	if err != nil {
		return false, err
	}
	o.FileList = fileList
	asset.SortFiles(o.FileList)
	return len(fileList) > 0, nil
}

func listFromMachines(objs []machineapi.Machine) ([]byte, error) {
	list := &metav1.List{
		TypeMeta: metav1.TypeMeta{
			APIVersion: "v1",
			Kind:       "List",
		},
	}
	for idx := range objs {
		list.Items = append(list.Items, runtime.RawExtension{Object: &objs[idx]})
	}

	return yaml.Marshal(list)
}

func listFromMachinesDeprecated(objs []clusterapi.Machine) ([]byte, error) {
	list := &metav1.List{
		TypeMeta: metav1.TypeMeta{
			APIVersion: "v1",
			Kind:       "List",
		},
	}
	for idx := range objs {
		list.Items = append(list.Items, runtime.RawExtension{Object: &objs[idx]})
	}

	return yaml.Marshal(list)
}
