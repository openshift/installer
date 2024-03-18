package machines

import (
	"context"
	"fmt"
	"path/filepath"
	"strings"

	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/utils/ptr"
	capa "sigs.k8s.io/cluster-api-provider-aws/v2/api/v1beta2"
	capi "sigs.k8s.io/cluster-api/api/v1beta1"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/yaml"

	configv1 "github.com/openshift/api/config/v1"
	"github.com/openshift/installer/pkg/asset"
	"github.com/openshift/installer/pkg/asset/installconfig"
	icazure "github.com/openshift/installer/pkg/asset/installconfig/azure"
	"github.com/openshift/installer/pkg/asset/machines/aws"
	"github.com/openshift/installer/pkg/asset/machines/azure"
	"github.com/openshift/installer/pkg/asset/machines/gcp"
	"github.com/openshift/installer/pkg/asset/machines/openstack"
	"github.com/openshift/installer/pkg/asset/machines/powervs"
	vspherecapi "github.com/openshift/installer/pkg/asset/machines/vsphere"
	"github.com/openshift/installer/pkg/asset/manifests/capiutils"
	"github.com/openshift/installer/pkg/asset/rhcos"
	"github.com/openshift/installer/pkg/clusterapi"
	rhcosutils "github.com/openshift/installer/pkg/rhcos"
	awstypes "github.com/openshift/installer/pkg/types/aws"
	awsdefaults "github.com/openshift/installer/pkg/types/aws/defaults"
	azuretypes "github.com/openshift/installer/pkg/types/azure"
	azuredefaults "github.com/openshift/installer/pkg/types/azure/defaults"
	gcptypes "github.com/openshift/installer/pkg/types/gcp"
	openstacktypes "github.com/openshift/installer/pkg/types/openstack"
	powervstypes "github.com/openshift/installer/pkg/types/powervs"
	vspheretypes "github.com/openshift/installer/pkg/types/vsphere"
)

var _ asset.WritableRuntimeAsset = (*ClusterAPI)(nil)

var machineManifestDir = filepath.Join(capiutils.ManifestDir, "machines")

// ClusterAPI is the asset for CAPI control-plane manifests.
type ClusterAPI struct {
	FileList []*asset.RuntimeFile
}

// Name returns a human friendly name for the operator.
func (c *ClusterAPI) Name() string {
	return "Cluster API Machine Manifests"
}

// Dependencies returns all of the dependencies directly needed by the
// ClusterAPI machines asset.
func (c *ClusterAPI) Dependencies() []asset.Asset {
	return []asset.Asset{
		&installconfig.InstallConfig{},
		&installconfig.ClusterID{},
		new(rhcos.Image),
	}
}

// Generate generates Cluster API machine manifests.
//
//nolint:gocyclo
func (c *ClusterAPI) Generate(dependencies asset.Parents) error {
	installConfig := &installconfig.InstallConfig{}
	clusterID := &installconfig.ClusterID{}
	rhcosImage := new(rhcos.Image)
	dependencies.Get(installConfig, clusterID, rhcosImage)

	// If the feature gate is not enabled, do not generate any manifests.
	if !capiutils.IsEnabled(installConfig) {
		return nil
	}

	c.FileList = []*asset.RuntimeFile{}

	var err error
	ic := installConfig.Config
	pool := *ic.ControlPlane
	ctx := context.TODO()

	switch ic.Platform.Name() {
	case awstypes.Name:
		subnets := map[string]string{}
		if len(ic.Platform.AWS.Subnets) > 0 {
			subnetMeta, err := installConfig.AWS.PrivateSubnets(ctx)
			if err != nil {
				return err
			}
			for id, subnet := range subnetMeta {
				subnets[subnet.Zone.Name] = id
			}
		}

		mpool := defaultAWSMachinePoolPlatform("master")

		osImage := strings.SplitN(string(*rhcosImage), ",", 2)
		osImageID := osImage[0]
		if len(osImage) == 2 {
			osImageID = "" // the AMI will be generated later on
		}
		mpool.AMIID = osImageID

		mpool.Set(ic.Platform.AWS.DefaultMachinePlatform)
		mpool.Set(pool.Platform.AWS)
		zoneDefaults := false
		if len(mpool.Zones) == 0 {
			if len(subnets) > 0 {
				for zone := range subnets {
					mpool.Zones = append(mpool.Zones, zone)
				}
			} else {
				mpool.Zones, err = installConfig.AWS.AvailabilityZones(ctx)
				if err != nil {
					return err
				}
				zoneDefaults = true
			}
		}

		if mpool.InstanceType == "" {
			topology := configv1.HighlyAvailableTopologyMode
			if pool.Replicas != nil && *pool.Replicas == 1 {
				topology = configv1.SingleReplicaTopologyMode
			}
			mpool.InstanceType, err = aws.PreferredInstanceType(ctx, installConfig.AWS, awsdefaults.InstanceTypes(installConfig.Config.Platform.AWS.Region, installConfig.Config.ControlPlane.Architecture, topology), mpool.Zones)
			if err != nil {
				logrus.Warn(errors.Wrap(err, "failed to find default instance type"))
				mpool.InstanceType = awsdefaults.InstanceTypes(installConfig.Config.Platform.AWS.Region, installConfig.Config.ControlPlane.Architecture, topology)[0]
			}
		}

		// if the list of zones is the default we need to try to filter the list in case there are some zones where the instance might not be available
		if zoneDefaults {
			mpool.Zones, err = aws.FilterZonesBasedOnInstanceType(ctx, installConfig.AWS, mpool.InstanceType, mpool.Zones)
			if err != nil {
				logrus.Warn(errors.Wrap(err, "failed to filter zone list"))
			}
		}

		pool.Platform.AWS = &mpool
		awsMachines, err := aws.GenerateMachines(
			clusterID.InfraID,
			installConfig.Config.Platform.AWS.Region,
			subnets,
			&pool,
			"master",
			installConfig.Config.Platform.AWS.UserTags,
		)
		if err != nil {
			return errors.Wrap(err, "failed to create master machine objects")
		}
		c.FileList = append(c.FileList, awsMachines...)

		// TODO(vincepri): The following code is almost duplicated from aws.AWSMachines.
		// Refactor and generalize around a bootstrap pool, with a single machine and
		// a custom openshift label to determine the bootstrap machine role, so we can
		// delete the machine when the stage is complete.
		bootstrapAWSMachine := &capa.AWSMachine{
			ObjectMeta: metav1.ObjectMeta{
				Name: capiutils.GenerateBoostrapMachineName(clusterID.InfraID),
				Labels: map[string]string{
					"cluster.x-k8s.io/control-plane": "",
					"install.openshift.io/bootstrap": "",
				},
			},
			Spec: capa.AWSMachineSpec{
				Ignition:             &capa.Ignition{Version: "3.2"},
				UncompressedUserData: ptr.To(true),
				InstanceType:         mpool.InstanceType,
				AMI:                  capa.AMIReference{ID: ptr.To(mpool.AMIID)},
				SSHKeyName:           ptr.To(""),
				IAMInstanceProfile:   fmt.Sprintf("%s-master-profile", clusterID.InfraID),
				PublicIP:             ptr.To(true),
				RootVolume: &capa.Volume{
					Size:          int64(mpool.EC2RootVolume.Size),
					Type:          capa.VolumeType(mpool.EC2RootVolume.Type),
					IOPS:          int64(mpool.EC2RootVolume.IOPS),
					Encrypted:     ptr.To(true),
					EncryptionKey: mpool.KMSKeyARN,
				},
			},
		}
		bootstrapAWSMachine.SetGroupVersionKind(capa.GroupVersion.WithKind("AWSMachine"))
		// TODO(mtulio): add public ip pool
		if installConfig.Config.Platform.AWS.PublicIpv4Pool != "" {
			bootstrapAWSMachine.Spec.PublicIpv4Pool = ptr.To(installConfig.Config.Platform.AWS.PublicIpv4Pool)
		}

		// Handle additional security groups.
		for _, sg := range mpool.AdditionalSecurityGroupIDs {
			bootstrapAWSMachine.Spec.AdditionalSecurityGroups = append(
				bootstrapAWSMachine.Spec.AdditionalSecurityGroups,
				capa.AWSResourceReference{ID: ptr.To(sg)},
			)
		}

		c.FileList = append(c.FileList, &asset.RuntimeFile{
			File:   asset.File{Filename: fmt.Sprintf("10_inframachine_%s.yaml", bootstrapAWSMachine.Name)},
			Object: bootstrapAWSMachine,
		})

		bootstrapMachine := &capi.Machine{
			ObjectMeta: metav1.ObjectMeta{
				Name: bootstrapAWSMachine.Name,
				Labels: map[string]string{
					"cluster.x-k8s.io/control-plane": "",
				},
			},
			Spec: capi.MachineSpec{
				ClusterName: clusterID.InfraID,
				Bootstrap: capi.Bootstrap{
					DataSecretName: ptr.To(fmt.Sprintf("%s-%s", clusterID.InfraID, "bootstrap")),
				},
				InfrastructureRef: v1.ObjectReference{
					APIVersion: capa.GroupVersion.String(),
					Kind:       "AWSMachine",
					Name:       bootstrapAWSMachine.Name,
				},
			},
		}
		bootstrapMachine.SetGroupVersionKind(capi.GroupVersion.WithKind("Machine"))

		c.FileList = append(c.FileList, &asset.RuntimeFile{
			File:   asset.File{Filename: fmt.Sprintf("10_machine_%s.yaml", bootstrapMachine.Name)},
			Object: bootstrapMachine,
		})
	case azuretypes.Name:
		mpool := defaultAzureMachinePoolPlatform()
		mpool.InstanceType = azuredefaults.ControlPlaneInstanceType(
			installConfig.Config.Platform.Azure.CloudName,
			installConfig.Config.Platform.Azure.Region,
			installConfig.Config.ControlPlane.Architecture,
		)
		mpool.OSDisk.DiskSizeGB = 1024
		if installConfig.Config.Platform.Azure.CloudName == azuretypes.StackCloud {
			mpool.OSDisk.DiskSizeGB = azuredefaults.AzurestackMinimumDiskSize
		}
		mpool.Set(ic.Platform.Azure.DefaultMachinePlatform)
		mpool.Set(pool.Platform.Azure)

		session, err := installConfig.Azure.Session()
		if err != nil {
			return fmt.Errorf("failed to fetch session: %w", err)
		}
		client := icazure.NewClient(session)

		if len(mpool.Zones) == 0 {
			// if no azs are given we set to []string{""} for convenience over later operations.
			// It means no-zoned for the machine API
			mpool.Zones = []string{""}
		}
		if len(mpool.Zones) == 0 {
			azs, err := client.GetAvailabilityZones(context.TODO(), ic.Platform.Azure.Region, mpool.InstanceType)
			if err != nil {
				return fmt.Errorf("failed to fetch availability zones: %w", err)
			}
			mpool.Zones = azs
			if len(azs) == 0 {
				// if no azs are given we set to []string{""} for convenience over later operations.
				// It means no-zoned for the machine API
				mpool.Zones = []string{""}
			}
		}
		// client.GetControlPlaneSubnet(context.TODO(), ic.Platform.Azure.ResourceGroupName, ic.Platform.Azure.VirtualNetwork, )

		if mpool.OSImage.Publisher != "" {
			img, ierr := client.GetMarketplaceImage(context.TODO(), ic.Platform.Azure.Region, mpool.OSImage.Publisher, mpool.OSImage.Offer, mpool.OSImage.SKU, mpool.OSImage.Version)
			if ierr != nil {
				return fmt.Errorf("failed to fetch marketplace image: %w", ierr)
			}
			// Publisher is case-sensitive and matched against exactly. Also the
			// Plan's publisher might not be exactly the same as the Image's
			// publisher
			if img.Plan != nil && img.Plan.Publisher != nil {
				mpool.OSImage.Publisher = *img.Plan.Publisher
			}
		}
		pool.Platform.Azure = &mpool
		subnet := ic.Azure.ControlPlaneSubnet

		capabilities, err := client.GetVMCapabilities(context.TODO(), mpool.InstanceType, installConfig.Config.Platform.Azure.Region)
		if err != nil {
			return err
		}
		hyperVGen, err := icazure.GetHyperVGenerationVersion(capabilities, "")
		if err != nil {
			return err
		}
		useImageGallery := installConfig.Azure.CloudName != azuretypes.StackCloud
		masterUserDataSecretName := "master-user-data"

		azureMachines, err := azure.GenerateMachines(installConfig.Config.Platform.Azure, &pool, masterUserDataSecretName, clusterID.InfraID, "master", capabilities, useImageGallery, installConfig.Config.Platform.Azure.UserTags, hyperVGen, subnet, ic.Azure.ResourceGroupName)
		if err != nil {
			return fmt.Errorf("failed to create master machine objects: %w", err)
		}

		c.FileList = append(c.FileList, azureMachines...)
	case gcptypes.Name:
		// Generate GCP master machines using ControPlane machinepool
		mpool := defaultGCPMachinePoolPlatform(pool.Architecture)
		mpool.Set(ic.Platform.GCP.DefaultMachinePlatform)
		mpool.Set(pool.Platform.GCP)
		if len(mpool.Zones) == 0 {
			azs, err := gcp.ZonesForInstanceType(ic.Platform.GCP.ProjectID, ic.Platform.GCP.Region, mpool.InstanceType)
			if err != nil {
				return errors.Wrap(err, "failed to fetch availability zones")
			}
			mpool.Zones = azs
		}
		pool.Platform.GCP = &mpool

		gcpMachines, err := gcp.GenerateMachines(
			installConfig,
			clusterID.InfraID,
			&pool,
			string(*rhcosImage),
		)
		if err != nil {
			return fmt.Errorf("failed to create master machine objects %w", err)
		}
		c.FileList = append(c.FileList, gcpMachines...)

		// Generate GCP bootstrap machines
		bootstrapMachines, err := gcp.GenerateBootstrapMachines(
			capiutils.GenerateBoostrapMachineName(clusterID.InfraID),
			installConfig,
			clusterID.InfraID,
			&pool,
			string(*rhcosImage),
		)
		if err != nil {
			return fmt.Errorf("failed to create bootstrap machine objects %w", err)
		}
		c.FileList = append(c.FileList, bootstrapMachines...)
	case vspheretypes.Name:
		mpool := defaultVSphereMachinePoolPlatform()
		mpool.NumCPUs = 4
		mpool.NumCoresPerSocket = 4
		mpool.MemoryMiB = 16384
		mpool.Set(ic.Platform.VSphere.DefaultMachinePlatform)
		mpool.Set(pool.Platform.VSphere)

		platform := ic.VSphere

		for _, v := range platform.VCenters {
			err := installConfig.VSphere.Networks(ctx, v, platform.FailureDomains)
			if err != nil {
				return err
			}
		}

		// The machinepool has no zones defined, there are FailureDomains
		// This is a vSphere zonal installation. Generate machinepool zone
		// list.

		fdCount := int64(len(ic.Platform.VSphere.FailureDomains))
		var idx int64
		if len(mpool.Zones) == 0 && len(ic.VSphere.FailureDomains) != 0 {
			for i := int64(0); i < *(ic.ControlPlane.Replicas); i++ {
				idx = i
				if idx >= fdCount {
					idx = i % fdCount
				}
				mpool.Zones = append(mpool.Zones, ic.VSphere.FailureDomains[idx].Name)
			}
		}

		pool.Platform.VSphere = &mpool
		templateName := clusterID.InfraID + "-rhcos"

		c.FileList, err = vspherecapi.GenerateMachines(ctx, clusterID.InfraID, ic, &pool, templateName, "master", installConfig.VSphere)
		if err != nil {
			return fmt.Errorf("unable to generate CAPI machines for vSphere %w", err)
		}
	case openstacktypes.Name:
		mpool := defaultOpenStackMachinePoolPlatform()
		mpool.Set(ic.Platform.OpenStack.DefaultMachinePlatform)
		mpool.Set(pool.Platform.OpenStack)
		pool.Platform.OpenStack = &mpool

		imageName, _ := rhcosutils.GenerateOpenStackImageName(string(*rhcosImage), clusterID.InfraID)

		for _, role := range []string{"master", "bootstrap"} {
			openStackMachines, err := openstack.GenerateMachines(
				clusterID.InfraID,
				ic,
				&pool,
				imageName,
				role,
			)
			if err != nil {
				return fmt.Errorf("failed to create machine objects: %w", err)
			}
			c.FileList = append(c.FileList, openStackMachines...)
		}
	case powervstypes.Name:
		// Generate PowerVS master machines using ControPlane machinepool
		mpool := defaultPowerVSMachinePoolPlatform(ic)
		mpool.Set(ic.Platform.PowerVS.DefaultMachinePlatform)
		mpool.Set(pool.Platform.PowerVS)
		pool.Platform.PowerVS = &mpool

		powervsMachines, err := powervs.GenerateMachines(
			clusterID.InfraID,
			&pool,
			"master",
		)
		if err != nil {
			return fmt.Errorf("failed to create master machine objects %w", err)
		}

		c.FileList = append(c.FileList, powervsMachines...)
	default:
		// TODO: support other platforms
	}

	// Create the machine manifests.
	for _, m := range c.FileList {
		objData, err := yaml.Marshal(m.Object)
		if err != nil {
			return errors.Wrapf(err, "failed to marshal Cluster API machine manifest %s", m.Filename)
		}
		m.Data = objData

		// If the filename is already a path, do not append the manifest dir.
		if filepath.Dir(m.Filename) == machineManifestDir {
			continue
		}
		m.Filename = filepath.Join(machineManifestDir, m.Filename)
	}
	asset.SortManifestFiles(c.FileList)
	return nil
}

// Files returns the files generated by the asset.
func (c *ClusterAPI) Files() []*asset.File {
	files := []*asset.File{}
	for _, f := range c.FileList {
		files = append(files, &f.File)
	}
	return files
}

// RuntimeFiles returns the files generated by the asset.
func (c *ClusterAPI) RuntimeFiles() []*asset.RuntimeFile {
	return c.FileList
}

// Load returns the openshift asset from disk.
func (c *ClusterAPI) Load(f asset.FileFetcher) (bool, error) {
	yamlFileList, err := f.FetchByPattern(filepath.Join(machineManifestDir, "*.yaml"))
	if err != nil {
		return false, errors.Wrap(err, "failed to load *.yaml files")
	}
	ymlFileList, err := f.FetchByPattern(filepath.Join(machineManifestDir, "*.yml"))
	if err != nil {
		return false, errors.Wrap(err, "failed to load *.yml files")
	}
	jsonFileList, err := f.FetchByPattern(filepath.Join(machineManifestDir, "*.json"))
	if err != nil {
		return false, errors.Wrap(err, "failed to load *.json files")
	}
	fileList := append(yamlFileList, ymlFileList...) //nolint:gocritic
	fileList = append(fileList, jsonFileList...)

	for _, file := range fileList {
		u := &unstructured.Unstructured{}
		if err := yaml.Unmarshal(file.Data, u); err != nil {
			return false, errors.Wrap(err, "failed to unmarshal file")
		}
		obj, err := clusterapi.Scheme.New(u.GroupVersionKind())
		if err != nil {
			return false, errors.Wrap(err, "failed to create object")
		}
		if err := clusterapi.Scheme.Convert(u, obj, nil); err != nil {
			return false, errors.Wrap(err, "failed to convert object")
		}
		c.FileList = append(c.FileList, &asset.RuntimeFile{
			File: asset.File{
				Filename: file.Filename,
				Data:     file.Data},
			Object: obj.(client.Object),
		})
	}

	asset.SortManifestFiles(c.FileList)
	return len(c.FileList) > 0, nil
}
