package manifests

import (
	"context"
	"encoding/base64"
	"io/ioutil"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/ghodss/yaml"
	"github.com/gophercloud/utils/openstack/clientconfig"
	"github.com/pkg/errors"

	"github.com/openshift/installer/pkg/asset"
	"github.com/openshift/installer/pkg/asset/installconfig"
	installconfigaws "github.com/openshift/installer/pkg/asset/installconfig/aws"
	"github.com/openshift/installer/pkg/asset/installconfig/gcp"
	"github.com/openshift/installer/pkg/asset/installconfig/ibmcloud"
	kubeconfig "github.com/openshift/installer/pkg/asset/installconfig/kubevirt"
	"github.com/openshift/installer/pkg/asset/installconfig/ovirt"
	"github.com/openshift/installer/pkg/asset/machines"
	osmachine "github.com/openshift/installer/pkg/asset/machines/openstack"
	openstackmanifests "github.com/openshift/installer/pkg/asset/manifests/openstack"
	"github.com/openshift/installer/pkg/asset/openshiftinstall"
	"github.com/openshift/installer/pkg/asset/password"
	"github.com/openshift/installer/pkg/asset/rhcos"
	"github.com/openshift/installer/pkg/asset/templates/content/openshift"
	"github.com/openshift/installer/pkg/types"
	awstypes "github.com/openshift/installer/pkg/types/aws"
	azuretypes "github.com/openshift/installer/pkg/types/azure"
	baremetaltypes "github.com/openshift/installer/pkg/types/baremetal"
	gcptypes "github.com/openshift/installer/pkg/types/gcp"
	ibmcloudtypes "github.com/openshift/installer/pkg/types/ibmcloud"
	kubevirttypes "github.com/openshift/installer/pkg/types/kubevirt"
	openstacktypes "github.com/openshift/installer/pkg/types/openstack"
	ovirttypes "github.com/openshift/installer/pkg/types/ovirt"
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
		&openshiftinstall.Config{},

		&openshift.CloudCredsSecret{},
		&openshift.KubeadminPasswordSecret{},
		&openshift.RoleCloudCredsSecretReader{},
		&openshift.BaremetalConfig{},
		new(rhcos.Image),
		&openshift.AzureCloudProviderSecret{},
	}
}

// Generate generates the respective operator config.yml files
func (o *Openshift) Generate(dependencies asset.Parents) error {
	installConfig := &installconfig.InstallConfig{}
	clusterID := &installconfig.ClusterID{}
	kubeadminPassword := &password.KubeadminPassword{}
	openshiftInstall := &openshiftinstall.Config{}
	dependencies.Get(installConfig, kubeadminPassword, clusterID, openshiftInstall)
	var cloudCreds cloudCredsSecretData
	platform := installConfig.Config.Platform.Name()
	switch platform {
	case awstypes.Name:
		ssn, err := installConfig.AWS.Session(context.TODO())
		if err != nil {
			return err
		}
		creds, err := ssn.Config.Credentials.Get()
		if err != nil {
			return err
		}
		if !installconfigaws.IsStaticCredentials(creds) {
			switch {
			case installConfig.Config.CredentialsMode == "":
				return errors.Errorf("AWS credentials provided by %s are not valid for default credentials mode", creds.ProviderName)
			case installConfig.Config.CredentialsMode != types.ManualCredentialsMode:
				return errors.Errorf("AWS credentials provided by %s are not valid for %s credentials mode", creds.ProviderName, installConfig.Config.CredentialsMode)
			}
		}
		cloudCreds = cloudCredsSecretData{
			AWS: &AwsCredsSecretData{
				Base64encodeAccessKeyID:     base64.StdEncoding.EncodeToString([]byte(creds.AccessKeyID)),
				Base64encodeSecretAccessKey: base64.StdEncoding.EncodeToString([]byte(creds.SecretAccessKey)),
			},
		}

	case azuretypes.Name:
		resourceGroupName := installConfig.Config.Azure.ClusterResourceGroupName(clusterID.InfraID)
		session, err := installConfig.Azure.Session()
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
	case gcptypes.Name:
		session, err := gcp.GetSession(context.TODO())
		if err != nil {
			return err
		}
		creds := session.Credentials.JSON
		cloudCreds = cloudCredsSecretData{
			GCP: &GCPCredsSecretData{
				Base64encodeServiceAccount: base64.StdEncoding.EncodeToString(creds),
			},
		}
	case ibmcloudtypes.Name:
		client, err := ibmcloud.NewClient()
		if err != nil {
			return err
		}
		cloudCreds = cloudCredsSecretData{
			IBMCloud: &IBMCloudCredsSecretData{
				Base64encodeAPIKey: base64.StdEncoding.EncodeToString([]byte(client.Authenticator.ApiKey)),
			},
		}
	case openstacktypes.Name:
		opts := new(clientconfig.ClientOpts)
		opts.Cloud = installConfig.Config.Platform.OpenStack.Cloud
		cloud, err := clientconfig.GetCloudFromYAML(opts)
		if err != nil {
			return err
		}

		// We need to replace the local cacert path with one that is used in OpenShift
		if cloud.CACertFile != "" {
			cloud.CACertFile = "/etc/kubernetes/static-pod-resources/configmaps/cloud-config/ca-bundle.pem"
		}

		clouds := make(map[string]map[string]*clientconfig.Cloud)
		clouds["clouds"] = map[string]*clientconfig.Cloud{
			osmachine.CloudName: cloud,
		}

		marshalled, err := yaml.Marshal(clouds)
		if err != nil {
			return err
		}

		cloudProviderConf, err := openstackmanifests.CloudProviderConfigSecret(cloud)
		if err != nil {
			return err
		}

		credsEncoded := base64.StdEncoding.EncodeToString(marshalled)
		credsINIEncoded := base64.StdEncoding.EncodeToString(cloudProviderConf)
		cloudCreds = cloudCredsSecretData{
			OpenStack: &OpenStackCredsSecretData{
				Base64encodeCloudCreds:    credsEncoded,
				Base64encodeCloudCredsINI: credsINIEncoded,
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
	case ovirttypes.Name:
		conf, err := ovirt.NewConfig()
		if err != nil {
			return err
		}

		if len(conf.CABundle) == 0 && len(conf.CAFile) > 0 {
			content, err := ioutil.ReadFile(conf.CAFile)
			if err != nil {
				return errors.Wrapf(err, "failed to read the cert file: %s", conf.CAFile)
			}
			conf.CABundle = strings.TrimSpace(string(content))
		}

		cloudCreds = cloudCredsSecretData{
			Ovirt: &OvirtCredsSecretData{
				Base64encodeURL:      base64.StdEncoding.EncodeToString([]byte(conf.URL)),
				Base64encodeUsername: base64.StdEncoding.EncodeToString([]byte(conf.Username)),
				Base64encodePassword: base64.StdEncoding.EncodeToString([]byte(conf.Password)),
				Base64encodeInsecure: base64.StdEncoding.EncodeToString([]byte(strconv.FormatBool(conf.Insecure))),
				Base64encodeCABundle: base64.StdEncoding.EncodeToString([]byte(conf.CABundle)),
			},
		}
	case kubevirttypes.Name:
		kubeconfigContent, err := kubeconfig.LoadKubeConfigContent()
		if err != nil {
			return err
		}
		cloudCreds = cloudCredsSecretData{
			Kubevirt: &KubevirtCredsSecretData{
				Base64encodedKubeconfig: base64.StdEncoding.EncodeToString(kubeconfigContent),
			},
		}
	}

	templateData := &openshiftTemplateData{
		CloudCreds:                   cloudCreds,
		Base64EncodedKubeadminPwHash: base64.StdEncoding.EncodeToString(kubeadminPassword.PasswordHash),
	}

	cloudCredsSecret := &openshift.CloudCredsSecret{}
	kubeadminPasswordSecret := &openshift.KubeadminPasswordSecret{}
	roleCloudCredsSecretReader := &openshift.RoleCloudCredsSecretReader{}
	baremetalConfig := &openshift.BaremetalConfig{}
	rhcosImage := new(rhcos.Image)

	dependencies.Get(
		cloudCredsSecret,
		kubeadminPasswordSecret,
		roleCloudCredsSecretReader,
		baremetalConfig,
		rhcosImage)

	assetData := map[string][]byte{
		"99_kubeadmin-password-secret.yaml": applyTemplateData(kubeadminPasswordSecret.Files()[0].Data, templateData),
	}

	switch platform {
	case awstypes.Name, openstacktypes.Name, vspheretypes.Name, azuretypes.Name, gcptypes.Name, ibmcloudtypes.Name, ovirttypes.Name, kubevirttypes.Name:
		if installConfig.Config.CredentialsMode != types.ManualCredentialsMode {
			assetData["99_cloud-creds-secret.yaml"] = applyTemplateData(cloudCredsSecret.Files()[0].Data, templateData)
		}
		assetData["99_role-cloud-creds-secret-reader.yaml"] = applyTemplateData(roleCloudCredsSecretReader.Files()[0].Data, templateData)
	case baremetaltypes.Name:
		bmTemplateData := baremetalTemplateData{
			Baremetal:                 installConfig.Config.Platform.BareMetal,
			ProvisioningOSDownloadURL: string(*rhcosImage),
		}
		assetData["99_baremetal-provisioning-config.yaml"] = applyTemplateData(baremetalConfig.Files()[0].Data, bmTemplateData)
	}

	if platform == azuretypes.Name && installConfig.Config.Azure.IsARO() {
		// config is used to created compatible secret to trigger azure cloud
		// controller config merge behaviour
		// https://github.com/openshift/origin/blob/90c050f5afb4c52ace82b15e126efe98fa798d88/vendor/k8s.io/legacy-cloud-providers/azure/azure_config.go#L83
		session, err := installConfig.Azure.Session()
		if err != nil {
			return err
		}
		config := struct {
			AADClientID     string `json:"aadClientId" yaml:"aadClientId"`
			AADClientSecret string `json:"aadClientSecret" yaml:"aadClientSecret"`
		}{
			AADClientID:     session.Credentials.ClientID,
			AADClientSecret: session.Credentials.ClientSecret,
		}

		b, err := yaml.Marshal(config)
		if err != nil {
			return err
		}

		azureCloudProviderSecret := &openshift.AzureCloudProviderSecret{}
		dependencies.Get(azureCloudProviderSecret)
		for _, f := range azureCloudProviderSecret.Files() {
			name := strings.TrimSuffix(filepath.Base(f.Filename), ".template")
			assetData[name] = applyTemplateData(f.Data, map[string]string{
				"CloudConfig": string(b),
			})
		}
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

	o.FileList = append(o.FileList, openshiftInstall.Files()...)

	asset.SortFiles(o.FileList)

	return nil
}

// Files returns the files generated by the asset.
func (o *Openshift) Files() []*asset.File {
	return o.FileList
}

// Load returns the openshift asset from disk.
func (o *Openshift) Load(f asset.FileFetcher) (bool, error) {
	yamlFileList, err := f.FetchByPattern(filepath.Join(openshiftManifestDir, "*.yaml"))
	if err != nil {
		return false, errors.Wrap(err, "failed to load *.yaml files")
	}
	ymlFileList, err := f.FetchByPattern(filepath.Join(openshiftManifestDir, "*.yml"))
	if err != nil {
		return false, errors.Wrap(err, "failed to load *.yml files")
	}
	jsonFileList, err := f.FetchByPattern(filepath.Join(openshiftManifestDir, "*.json"))
	if err != nil {
		return false, errors.Wrap(err, "failed to load *.json files")
	}
	fileList := append(yamlFileList, ymlFileList...)
	fileList = append(fileList, jsonFileList...)

	for _, file := range fileList {
		if machines.IsMachineManifest(file) {
			continue
		}

		o.FileList = append(o.FileList, file)
	}

	asset.SortFiles(o.FileList)
	return len(o.FileList) > 0, nil
}
