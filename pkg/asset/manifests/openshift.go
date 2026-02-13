package manifests

import (
	"context"
	"encoding/base64"
	"fmt"
	"os"
	"path"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/gophercloud/utils/v2/openstack/clientconfig"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"k8s.io/apimachinery/pkg/util/sets"
	"sigs.k8s.io/yaml"

	"github.com/openshift/installer/pkg/asset"
	"github.com/openshift/installer/pkg/asset/installconfig"
	installconfigaws "github.com/openshift/installer/pkg/asset/installconfig/aws"
	"github.com/openshift/installer/pkg/asset/installconfig/gcp"
	"github.com/openshift/installer/pkg/asset/installconfig/ibmcloud"
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
	nutanixtypes "github.com/openshift/installer/pkg/types/nutanix"
	openstacktypes "github.com/openshift/installer/pkg/types/openstack"
	ovirttypes "github.com/openshift/installer/pkg/types/ovirt"
	powervctypes "github.com/openshift/installer/pkg/types/powervc"
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
		&FeatureGate{},

		&openshift.CloudCredsSecret{},
		&openshift.KubeadminPasswordSecret{},
		&openshift.RoleCloudCredsSecretReader{},
		&openshift.BaremetalConfig{},
		new(rhcos.Image),
		&openshift.AzureCloudProviderSecret{},
		&openshift.VSphereMachineAPICredsSecret{},
		&openshift.VSphereCSIDriverCredsSecret{},
		&openshift.VSphereCloudControllerCredsSecret{},
		&openshift.VSphereDiagnosticsCredsSecret{},
	}
}

// Generate generates the respective operator config.yml files
//
//nolint:gocyclo
func (o *Openshift) Generate(ctx context.Context, dependencies asset.Parents) error {
	installConfig := &installconfig.InstallConfig{}
	clusterID := &installconfig.ClusterID{}
	kubeadminPassword := &password.KubeadminPassword{}
	openshiftInstall := &openshiftinstall.Config{}
	featureGate := &FeatureGate{}
	dependencies.Get(installConfig, kubeadminPassword, clusterID, openshiftInstall, featureGate)
	var cloudCreds cloudCredsSecretData
	platform := installConfig.Config.Platform.Name()
	switch platform {
	case awstypes.Name:
		ssn, err := installConfig.AWS.Session(ctx)
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
		session, err := gcp.GetSession(ctx)
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
		client, err := ibmcloud.NewClient(installConfig.Config.Platform.IBMCloud.ServiceEndpoints)
		if err != nil {
			return err
		}
		cloudCreds = cloudCredsSecretData{
			IBMCloud: &IBMCloudCredsSecretData{
				Base64encodeAPIKey: base64.StdEncoding.EncodeToString([]byte(client.GetAPIKey())),
			},
		}
	case openstacktypes.Name, powervctypes.Name:
		opts := new(clientconfig.ClientOpts)
		opts.Cloud = installConfig.Config.Platform.OpenStack.Cloud
		cloud, err := clientconfig.GetCloudFromYAML(opts)
		if err != nil {
			return err
		}

		var caCert []byte
		if cloud.CACertFile != "" {
			var err error
			caCert, err = os.ReadFile(cloud.CACertFile)
			if err != nil {
				return err
			}
			// We need to replace the local cacert path with one that is used in OpenShift
			cloud.CACertFile = "/etc/kubernetes/static-pod-resources/configmaps/cloud-config/ca-bundle.pem"
		}

		// Application credentials are easily rotated in the event of a leak and should be preferred. Encourage their use.
		authTypes := sets.New(clientconfig.AuthPassword, clientconfig.AuthV2Password, clientconfig.AuthV3Password)
		if cloud.AuthInfo != nil && authTypes.Has(cloud.AuthType) {
			logrus.Warnf(
				"clouds.yaml file is using %q type auth. Consider using the %q auth type instead to rotate credentials more easily.",
				cloud.AuthType,
				clientconfig.AuthV3ApplicationCredential,
			)
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
		cloudProviderConfEncoded := base64.StdEncoding.EncodeToString(cloudProviderConf)
		caCertEncoded := base64.StdEncoding.EncodeToString(caCert)
		cloudCreds = cloudCredsSecretData{
			OpenStack: &OpenStackCredsSecretData{
				Base64encodeCloudsYAML: credsEncoded,
				Base64encodeCloudsConf: cloudProviderConfEncoded,
				Base64encodeCACert:     caCertEncoded,
			},
		}
	case vspheretypes.Name:
		// Check if any vCenter has component credentials defined
		hasComponentCredentials := false
		for _, vCenter := range installConfig.Config.VSphere.VCenters {
			if vCenter.ComponentCredentials != nil {
				hasComponentCredentials = true
				break
			}
		}

		if hasComponentCredentials {
			// Generate per-component secrets
			cloudCreds = generateVSphereComponentSecrets(installConfig.Config.VSphere.VCenters)
		} else {
			// Legacy mode: generate single shared secret
			vsphereCredList := make([]*VSphereCredsSecretData, 0)

			for _, vCenter := range installConfig.Config.VSphere.VCenters {
				vsphereCred := VSphereCredsSecretData{
					VCenter:              vCenter.Server,
					Base64encodeUsername: base64.StdEncoding.EncodeToString([]byte(vCenter.Username)),
					Base64encodePassword: base64.StdEncoding.EncodeToString([]byte(vCenter.Password)),
				}
				vsphereCredList = append(vsphereCredList, &vsphereCred)
			}

			cloudCreds = cloudCredsSecretData{
				VSphere: &vsphereCredList,
			}
		}
	case ovirttypes.Name:
		conf, err := ovirt.NewConfig()
		if err != nil {
			return err
		}

		if len(conf.CABundle) == 0 && len(conf.CAFile) > 0 {
			content, err := os.ReadFile(conf.CAFile)
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
	case nutanixtypes.Name:
		// Format credentials as JSON array according to Nutanix format
		credentialsData := fmt.Sprintf(`[{
		"type": "basic_auth",
		"data": {
			"prismCentral": {
				"username": "%s",
				"password": "%s"
			}
		}
	}]`,
			installConfig.Config.Platform.Nutanix.PrismCentral.Username,
			installConfig.Config.Platform.Nutanix.PrismCentral.Password,
		)

		cloudCreds = cloudCredsSecretData{
			Nutanix: &NutanixCredsSecretData{
				Base64encodeCredentials: base64.StdEncoding.EncodeToString([]byte(credentialsData)),
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
	vsphereMachineAPICredsSecret := &openshift.VSphereMachineAPICredsSecret{}
	vsphereCSIDriverCredsSecret := &openshift.VSphereCSIDriverCredsSecret{}
	vsphereCloudControllerCredsSecret := &openshift.VSphereCloudControllerCredsSecret{}
	vsphereDiagnosticsCredsSecret := &openshift.VSphereDiagnosticsCredsSecret{}

	dependencies.Get(
		cloudCredsSecret,
		kubeadminPasswordSecret,
		roleCloudCredsSecretReader,
		baremetalConfig,
		rhcosImage,
		vsphereMachineAPICredsSecret,
		vsphereCSIDriverCredsSecret,
		vsphereCloudControllerCredsSecret,
		vsphereDiagnosticsCredsSecret)

	assetData := map[string][]byte{
		"99_kubeadmin-password-secret.yaml": applyTemplateData(kubeadminPasswordSecret.Files()[0].Data, templateData),
	}

	switch platform {
	case vspheretypes.Name:
		// Check if using component credentials
		hasComponentCredentials := false
		for _, vCenter := range installConfig.Config.VSphere.VCenters {
			if vCenter.ComponentCredentials != nil {
				hasComponentCredentials = true
				break
			}
		}

		if hasComponentCredentials {
			// Render component-specific secrets when using component credentials
			if installConfig.Config.CredentialsMode != types.ManualCredentialsMode {
				assetData["99_vsphere-creds-machine-api.yaml"] = applyTemplateData(vsphereMachineAPICredsSecret.Files()[0].Data, templateData)
				assetData["99_vsphere-creds-csi-driver.yaml"] = applyTemplateData(vsphereCSIDriverCredsSecret.Files()[0].Data, templateData)
				assetData["99_vsphere-creds-cloud-controller.yaml"] = applyTemplateData(vsphereCloudControllerCredsSecret.Files()[0].Data, templateData)
				assetData["99_vsphere-creds-diagnostics.yaml"] = applyTemplateData(vsphereDiagnosticsCredsSecret.Files()[0].Data, templateData)
			}
			// Note: Role is not needed for component-specific secrets since CCO doesn't use them
		} else {
			// Legacy mode: render single shared secret
			if installConfig.Config.CredentialsMode != types.ManualCredentialsMode {
				assetData["99_cloud-creds-secret.yaml"] = applyTemplateData(cloudCredsSecret.Files()[0].Data, templateData)
			}
			assetData["99_role-cloud-creds-secret-reader.yaml"] = applyTemplateData(roleCloudCredsSecretReader.Files()[0].Data, templateData)
		}
	case awstypes.Name, openstacktypes.Name, powervctypes.Name, azuretypes.Name, gcptypes.Name, ibmcloudtypes.Name, ovirttypes.Name:
		if installConfig.Config.CredentialsMode != types.ManualCredentialsMode {
			assetData["99_cloud-creds-secret.yaml"] = applyTemplateData(cloudCredsSecret.Files()[0].Data, templateData)
		}
		assetData["99_role-cloud-creds-secret-reader.yaml"] = applyTemplateData(roleCloudCredsSecretReader.Files()[0].Data, templateData)
	case baremetaltypes.Name:
		bmTemplateData := baremetalTemplateData{
			Baremetal:                 installConfig.Config.Platform.BareMetal,
			ProvisioningOSDownloadURL: rhcosImage.ControlPlane,
		}
		assetData["99_baremetal-provisioning-config.yaml"] = applyTemplateData(baremetalConfig.Files()[0].Data, bmTemplateData)
	}

	o.FileList = []*asset.File{}
	for name, data := range assetData {
		if len(data) == 0 {
			continue
		}
		o.FileList = append(o.FileList, &asset.File{
			Filename: path.Join(openshiftManifestDir, name),
			Data:     data,
		})
	}

	o.FileList = append(o.FileList, openshiftInstall.Files()...)
	o.FileList = append(o.FileList, featureGate.Files()...)

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

// generateVSphereComponentSecrets creates per-component credential data for vSphere.
// For each component, it uses component-specific credentials if available,
// otherwise falls back to the main vCenter credentials.
func generateVSphereComponentSecrets(vcenters []vspheretypes.VCenter) cloudCredsSecretData {
	// Create separate credential lists for each component
	machineAPICreds := make([]*VSphereCredsSecretData, 0)
	csiDriverCreds := make([]*VSphereCredsSecretData, 0)
	cloudControllerCreds := make([]*VSphereCredsSecretData, 0)
	diagnosticsCreds := make([]*VSphereCredsSecretData, 0)

	for _, vcenter := range vcenters {
		// Fallback to main credentials if component credentials not specified
		machineAPIUser := vcenter.Username
		machineAPIPassword := vcenter.Password
		csiDriverUser := vcenter.Username
		csiDriverPassword := vcenter.Password
		cloudControllerUser := vcenter.Username
		cloudControllerPassword := vcenter.Password
		diagnosticsUser := vcenter.Username
		diagnosticsPassword := vcenter.Password

		// Override with component-specific credentials if provided
		if vcenter.ComponentCredentials != nil {
			if vcenter.ComponentCredentials.MachineAPI != nil {
				machineAPIUser = vcenter.ComponentCredentials.MachineAPI.User
				machineAPIPassword = vcenter.ComponentCredentials.MachineAPI.Password
			}
			if vcenter.ComponentCredentials.CSIDriver != nil {
				csiDriverUser = vcenter.ComponentCredentials.CSIDriver.User
				csiDriverPassword = vcenter.ComponentCredentials.CSIDriver.Password
			}
			if vcenter.ComponentCredentials.CloudController != nil {
				cloudControllerUser = vcenter.ComponentCredentials.CloudController.User
				cloudControllerPassword = vcenter.ComponentCredentials.CloudController.Password
			}
			if vcenter.ComponentCredentials.Diagnostics != nil {
				diagnosticsUser = vcenter.ComponentCredentials.Diagnostics.User
				diagnosticsPassword = vcenter.ComponentCredentials.Diagnostics.Password
			}
		}

		// Create credential entries for each component
		machineAPICreds = append(machineAPICreds, &VSphereCredsSecretData{
			VCenter:              vcenter.Server,
			Base64encodeUsername: base64.StdEncoding.EncodeToString([]byte(machineAPIUser)),
			Base64encodePassword: base64.StdEncoding.EncodeToString([]byte(machineAPIPassword)),
		})

		csiDriverCreds = append(csiDriverCreds, &VSphereCredsSecretData{
			VCenter:              vcenter.Server,
			Base64encodeUsername: base64.StdEncoding.EncodeToString([]byte(csiDriverUser)),
			Base64encodePassword: base64.StdEncoding.EncodeToString([]byte(csiDriverPassword)),
		})

		cloudControllerCreds = append(cloudControllerCreds, &VSphereCredsSecretData{
			VCenter:              vcenter.Server,
			Base64encodeUsername: base64.StdEncoding.EncodeToString([]byte(cloudControllerUser)),
			Base64encodePassword: base64.StdEncoding.EncodeToString([]byte(cloudControllerPassword)),
		})

		diagnosticsCreds = append(diagnosticsCreds, &VSphereCredsSecretData{
			VCenter:              vcenter.Server,
			Base64encodeUsername: base64.StdEncoding.EncodeToString([]byte(diagnosticsUser)),
			Base64encodePassword: base64.StdEncoding.EncodeToString([]byte(diagnosticsPassword)),
		})
	}

	// Return component-specific credentials
	return cloudCredsSecretData{
		VSphereComponents: &VSphereComponentCredsSecretData{
			MachineAPI:      &machineAPICreds,
			CSIDriver:       &csiDriverCreds,
			CloudController: &cloudControllerCreds,
			Diagnostics:     &diagnosticsCreds,
		},
	}
}
