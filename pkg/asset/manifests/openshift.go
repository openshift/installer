package manifests

import (
	"encoding/base64"
	"path/filepath"

	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/ghodss/yaml"

	"github.com/gophercloud/utils/openstack/clientconfig"
	"github.com/openshift/installer/pkg/asset"
	"github.com/openshift/installer/pkg/asset/installconfig"
	"github.com/openshift/installer/pkg/asset/installconfig/azure"
	"github.com/openshift/installer/pkg/asset/machines"

	osmachine "github.com/openshift/installer/pkg/asset/machines/openstack"
	"github.com/openshift/installer/pkg/asset/password"
	"github.com/openshift/installer/pkg/asset/templates/content/openshift"
	awstypes "github.com/openshift/installer/pkg/types/aws"
	azuretypes "github.com/openshift/installer/pkg/types/azure"
	openstacktypes "github.com/openshift/installer/pkg/types/openstack"
	vspheretypes "github.com/openshift/installer/pkg/types/vsphere"
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
		&installconfig.ClusterID{},
		&password.KubeadminPassword{},

		&openshift.BindingDiscovery{},
		&openshift.CloudCredsSecret{},
		&openshift.KubeadminPasswordSecret{},
		&openshift.RoleCloudCredsSecretReader{},
		&openshift.RolebindingCloudCredsSecretReader{},
	}
}

// Generate generates the respective operator config.yml files
func (o *Openshift) Generate(dependencies asset.Parents) error {
	installConfig := &installconfig.InstallConfig{}
	clusterID := &installconfig.ClusterID{}
	kubeadminPassword := &password.KubeadminPassword{}
	dependencies.Get(installConfig, kubeadminPassword, clusterID)
	var cloudCreds cloudCredsSecretData
	platform := installConfig.Config.Platform.Name()
	switch platform {
	case awstypes.Name:
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

	case azuretypes.Name:
		resourceGroupName := clusterID.InfraID + "-rg"
		session, err := azure.GetSession()
		if err != nil {
			return err
		}
		creds := session.Credentials
		cloudCreds = cloudCredsSecretData{
			Azure: &AzureCredsSecretData{
				Base64encodeSubscriptionID: base64.StdEncoding.EncodeToString([]byte(creds.SubscriptionID)),
				Base64encodeClientID:       base64.StdEncoding.EncodeToString([]byte(creds.ClientID)),
				Base64encodeClientSecret:   base64.StdEncoding.EncodeToString([]byte(creds.ClientSecret)),
				Base64encodeTenantID:       base64.StdEncoding.EncodeToString([]byte(creds.TenantID)),
				Base64encodeResourcePrefix: base64.StdEncoding.EncodeToString([]byte(clusterID.InfraID)),
				Base64encodeResourceGroup:  base64.StdEncoding.EncodeToString([]byte(resourceGroupName)),
				Base64encodeRegion:         base64.StdEncoding.EncodeToString([]byte(installConfig.Config.Azure.Region)),
			},
		}
	case openstacktypes.Name:
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

		marshalled, err := yaml.Marshal(clouds)
		if err != nil {
			return err
		}

		credsEncoded := base64.StdEncoding.EncodeToString(marshalled)
		cloudCreds = cloudCredsSecretData{
			OpenStack: &OpenStackCredsSecretData{
				Base64encodeCloudCreds: credsEncoded,
			},
		}
	case vspheretypes.Name:
		cloudCreds = cloudCredsSecretData{
			VSphere: &VSphereCredsSecretData{
				VCenter:              installConfig.Config.VSphere.VCenter,
				Base64encodeUsername: base64.StdEncoding.EncodeToString([]byte(installConfig.Config.VSphere.Username)),
				Base64encodePassword: base64.StdEncoding.EncodeToString([]byte(installConfig.Config.VSphere.Password)),
			},
		}
	}

	templateData := &openshiftTemplateData{
		CloudCreds:                   cloudCreds,
		Base64EncodedKubeadminPwHash: base64.StdEncoding.EncodeToString(kubeadminPassword.PasswordHash),
	}

	bindingDiscovery := &openshift.BindingDiscovery{}
	cloudCredsSecret := &openshift.CloudCredsSecret{}
	kubeadminPasswordSecret := &openshift.KubeadminPasswordSecret{}
	roleCloudCredsSecretReader := &openshift.RoleCloudCredsSecretReader{}
	rolebindingCloudCredsSecretReader := &openshift.RolebindingCloudCredsSecretReader{}
	dependencies.Get(
		bindingDiscovery,
		cloudCredsSecret,
		kubeadminPasswordSecret,
		roleCloudCredsSecretReader,
		rolebindingCloudCredsSecretReader)

	assetData := map[string][]byte{
		"99_binding-discovery.yaml":         []byte(bindingDiscovery.Files()[0].Data),
		"99_kubeadmin-password-secret.yaml": applyTemplateData(kubeadminPasswordSecret.Files()[0].Data, templateData),
	}

	switch platform {
	case awstypes.Name, openstacktypes.Name, vspheretypes.Name, azuretypes.Name:
		assetData["99_cloud-creds-secret.yaml"] = applyTemplateData(cloudCredsSecret.Files()[0].Data, templateData)
		assetData["99_role-cloud-creds-secret-reader.yaml"] = applyTemplateData(roleCloudCredsSecretReader.Files()[0].Data, templateData)
	}

	switch platform {
	case openstacktypes.Name:
		assetData["99_rolebinding-cloud-creds-secret-reader.yaml"] = applyTemplateData(rolebindingCloudCredsSecretReader.Files()[0].Data, templateData)
	}

	o.FileList = []*asset.File{}
	for name, data := range assetData {
		if len(data) == 0 {
			continue
		}
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

	for _, file := range fileList {
		if machines.IsMachineManifest(file) {
			continue
		}

		o.FileList = append(o.FileList, file)
	}

	asset.SortFiles(o.FileList)
	return len(o.FileList) > 0, nil
}
