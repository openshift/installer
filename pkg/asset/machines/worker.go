package machines

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/serializer"
	"k8s.io/apimachinery/pkg/util/intstr"
	"k8s.io/utils/pointer"
	ipamv1 "sigs.k8s.io/cluster-api/api/ipam/v1beta1" //nolint:staticcheck //CORS-3563
	"sigs.k8s.io/yaml"

	configv1 "github.com/openshift/api/config/v1"
	"github.com/openshift/api/features"
	machinev1 "github.com/openshift/api/machine/v1"
	machinev1alpha1 "github.com/openshift/api/machine/v1alpha1"
	machinev1beta1 "github.com/openshift/api/machine/v1beta1"
	mcfgv1 "github.com/openshift/api/machineconfiguration/v1"
	baremetalapi "github.com/openshift/cluster-api-provider-baremetal/pkg/apis"
	baremetalprovider "github.com/openshift/cluster-api-provider-baremetal/pkg/apis/baremetal/v1alpha1"
	libvirtapi "github.com/openshift/cluster-api-provider-libvirt/pkg/apis"
	libvirtprovider "github.com/openshift/cluster-api-provider-libvirt/pkg/apis/libvirtproviderconfig/v1beta1"
	ovirtproviderapi "github.com/openshift/cluster-api-provider-ovirt/pkg/apis"
	ovirtprovider "github.com/openshift/cluster-api-provider-ovirt/pkg/apis/ovirtprovider/v1beta1"
	"github.com/openshift/installer/pkg/asset"
	"github.com/openshift/installer/pkg/asset/ignition/machine"
	"github.com/openshift/installer/pkg/asset/installconfig"
	icaws "github.com/openshift/installer/pkg/asset/installconfig/aws"
	icgcp "github.com/openshift/installer/pkg/asset/installconfig/gcp"
	powervsconfig "github.com/openshift/installer/pkg/asset/installconfig/powervs"
	"github.com/openshift/installer/pkg/asset/machines/aws"
	"github.com/openshift/installer/pkg/asset/machines/azure"
	"github.com/openshift/installer/pkg/asset/machines/baremetal"
	"github.com/openshift/installer/pkg/asset/machines/gcp"
	"github.com/openshift/installer/pkg/asset/machines/ibmcloud"
	"github.com/openshift/installer/pkg/asset/machines/machineconfig"
	"github.com/openshift/installer/pkg/asset/machines/nutanix"
	"github.com/openshift/installer/pkg/asset/machines/openstack"
	"github.com/openshift/installer/pkg/asset/machines/ovirt"
	"github.com/openshift/installer/pkg/asset/machines/powervs"
	"github.com/openshift/installer/pkg/asset/machines/vsphere"
	"github.com/openshift/installer/pkg/asset/rhcos"
	rhcosutils "github.com/openshift/installer/pkg/rhcos"
	"github.com/openshift/installer/pkg/types"
	awstypes "github.com/openshift/installer/pkg/types/aws"
	awsdefaults "github.com/openshift/installer/pkg/types/aws/defaults"
	azuretypes "github.com/openshift/installer/pkg/types/azure"
	azuredefaults "github.com/openshift/installer/pkg/types/azure/defaults"
	baremetaltypes "github.com/openshift/installer/pkg/types/baremetal"
	externaltypes "github.com/openshift/installer/pkg/types/external"
	gcptypes "github.com/openshift/installer/pkg/types/gcp"
	ibmcloudtypes "github.com/openshift/installer/pkg/types/ibmcloud"
	nonetypes "github.com/openshift/installer/pkg/types/none"
	nutanixtypes "github.com/openshift/installer/pkg/types/nutanix"
	openstacktypes "github.com/openshift/installer/pkg/types/openstack"
	ovirttypes "github.com/openshift/installer/pkg/types/ovirt"
	powervctypes "github.com/openshift/installer/pkg/types/powervc"
	powervstypes "github.com/openshift/installer/pkg/types/powervs"
	powervsdefaults "github.com/openshift/installer/pkg/types/powervs/defaults"
	vspheretypes "github.com/openshift/installer/pkg/types/vsphere"
	ibmcloudapi "github.com/openshift/machine-api-provider-ibmcloud/pkg/apis"
	ibmcloudprovider "github.com/openshift/machine-api-provider-ibmcloud/pkg/apis/ibmcloudprovider/v1"
)

const (
	// workerMachineSetFileName is the format string for constructing the worker MachineSet filenames.
	workerMachineSetFileName = "99_openshift-cluster-api_worker-machineset-%s.yaml"

	// workerMachineFileName is the format string for constructing the worker Machine filenames.
	workerMachineFileName = "99_openshift-cluster-api_worker-machines-%s.yaml"

	// workerUserDataFileName is the filename used for the worker user-data secret.
	workerUserDataFileName = "99_openshift-cluster-api_worker-user-data-secret.yaml"

	// decimalRootVolumeSize is the size in GB we use for some platforms.
	// See below.
	decimalRootVolumeSize = 120

	// powerOfTwoRootVolumeSize is the size in GB we use for other platforms.
	// The reasons for the specific choices between these two may boil down
	// to which section of code the person adding a platform was copy-pasting from.
	// https://github.com/openshift/openshift-docs/blob/main/modules/installation-requirements-user-infra.adoc#minimum-resource-requirements
	powerOfTwoRootVolumeSize = 128
)

var (
	workerMachineSetFileNamePattern = fmt.Sprintf(workerMachineSetFileName, "*")
	workerMachineFileNamePattern    = fmt.Sprintf(workerMachineFileName, "*")
	workerIPClaimFileNamePattern    = fmt.Sprintf(ipClaimFileName, "*worker*")
	workerIPAddressFileNamePattern  = fmt.Sprintf(ipAddressFileName, "*worker*")

	_ asset.WritableAsset = (*Worker)(nil)
)

func defaultAWSMachinePoolPlatform(poolName string) awstypes.MachinePool {
	defaultEBSType := awstypes.VolumeTypeGp3

	// gp3 is not offered in all local-zones locations used by Edge Pools.
	// Once it is available, it can be used as default for all machine pools.
	// https://aws.amazon.com/about-aws/global-infrastructure/localzones/features
	if poolName == types.MachinePoolEdgeRoleName {
		defaultEBSType = awstypes.VolumeTypeGp2
	}
	return awstypes.MachinePool{
		EC2RootVolume: awstypes.EC2RootVolume{
			Type: defaultEBSType,
			Size: decimalRootVolumeSize,
		},
	}
}

func defaultAzureMachinePoolPlatform(env azuretypes.CloudEnvironment) azuretypes.MachinePool {
	idType := azuretypes.VMIdentityUserAssigned
	if env == azuretypes.StackCloud {
		idType = azuretypes.VMIdentityNone
	}

	return azuretypes.MachinePool{
		OSDisk: azuretypes.OSDisk{
			DiskSizeGB: powerOfTwoRootVolumeSize,
			DiskType:   azuretypes.DefaultDiskType,
		},
		Identity: &azuretypes.VMIdentity{Type: idType},
	}
}

func defaultGCPMachinePoolPlatform(arch types.Architecture) gcptypes.MachinePool {
	return gcptypes.MachinePool{
		InstanceType: icgcp.DefaultInstanceTypeForArch(arch),
		OSDisk: gcptypes.OSDisk{
			DiskSizeGB: powerOfTwoRootVolumeSize,
			DiskType:   "pd-ssd",
		},
	}
}

func defaultIBMCloudMachinePoolPlatform() ibmcloudtypes.MachinePool {
	return ibmcloudtypes.MachinePool{
		InstanceType: "bx2-4x16",
	}
}

func defaultOpenStackMachinePoolPlatform() openstacktypes.MachinePool {
	return openstacktypes.MachinePool{
		Zones: []string{""},
	}
}

func defaultBareMetalMachinePoolPlatform() baremetaltypes.MachinePool {
	return baremetaltypes.MachinePool{}
}

func defaultOvirtMachinePoolPlatform() ovirttypes.MachinePool {
	return ovirttypes.MachinePool{
		CPU: &ovirttypes.CPU{
			Cores:   4,
			Sockets: 1,
			Threads: 1,
		},
		MemoryMB: 16348,
		OSDisk: &ovirttypes.Disk{
			SizeGB: decimalRootVolumeSize,
		},
		VMType:            ovirttypes.VMTypeServer,
		AutoPinningPolicy: ovirttypes.AutoPinningNone,
	}
}

func defaultVSphereMachinePoolPlatform() vspheretypes.MachinePool {
	return vspheretypes.MachinePool{
		NumCPUs:           4,
		NumCoresPerSocket: 4,
		MemoryMiB:         16384,
		OSDisk: vspheretypes.OSDisk{
			DiskSizeGB: decimalRootVolumeSize,
		},
	}
}

func defaultPowerVSMachinePoolPlatform(ic *types.InstallConfig) powervstypes.MachinePool {
	var (
		client   *powervsconfig.Client
		fallback = false
		sysTypes []string
		sysType  = "s922"
		err      error
	)

	// Update the saved session storage with the install config since the session
	// storage is used as the defaults.
	err = powervsconfig.UpdateSessionStoreToAuthFile(&powervsconfig.SessionStore{
		ID:                   ic.PowerVS.UserID,
		DefaultRegion:        ic.PowerVS.Region,
		DefaultZone:          ic.PowerVS.Zone,
		PowerVSResourceGroup: ic.PowerVS.PowerVSResourceGroup,
	})
	if err != nil {
		fallback = true
		logrus.Warnf("could not UpdateSessionStoreToAuthFile in defaultPowerVSMachinePoolPlatform")
	}

	client, err = powervsconfig.NewClient()
	if err != nil {
		fallback = true
		logrus.Warnf("could not get client in defaultPowerVSMachinePoolPlatform")
	} else {
		sysTypes, err = client.GetDatacenterSupportedSystems(context.Background(), ic.PowerVS.Zone)
		if err != nil {
			fallback = true
			logrus.Warnf("For given zone %v, GetDatacenterSupportedSystems returns %v", ic.PowerVS.Zone, err)
		} else {
			// Is the hardcoded default of s922 in the list?
			found := false
			for _, st := range sysTypes {
				if st == sysType {
					found = true
					break
				}
			}
			if !found {
				sysType = sysTypes[0]
			}
		}
	}

	if fallback {
		// Fallback to hardcoded list
		sysTypes, err = powervstypes.AvailableSysTypes(ic.PowerVS.Region, ic.PowerVS.Zone)
		if err == nil {
			sysType = sysTypes[0]
		} else {
			logrus.Warnf("For given zone %v, AvailableSysTypes returns %v", ic.PowerVS.Zone, err)
		}
	}

	logrus.Debugf("defaultPowerVSMachinePoolPlatform: using a default SysType of %s with values to choose from %v", sysType, sysTypes)

	return powervstypes.MachinePool{
		MemoryGiB:  32,
		Processors: intstr.FromString("0.5"),
		ProcType:   machinev1.PowerVSProcessorTypeShared,
		SysType:    sysType,
	}
}

func defaultNutanixMachinePoolPlatform() nutanixtypes.MachinePool {
	return nutanixtypes.MachinePool{
		NumCPUs:           4,
		NumCoresPerSocket: 1,
		MemoryMiB:         16384,
		OSDisk: nutanixtypes.OSDisk{
			DiskSizeGiB: decimalRootVolumeSize,
		},
	}
}

// awsSetPreferredInstanceByEdgeZone discovers supported instanceType for each edge pool
// using the existing preferred instance list used by worker compute pool.
// Each machine set in the edge pool, created for each zone, can use different instance
// types depending on the instance offerings in the location (Local Zones).
func awsSetPreferredInstanceByEdgeZone(ctx context.Context, defaultTypes []string, meta *icaws.Metadata, zones icaws.Zones) (ok bool) {
	allZonesFound := true
	for zone := range zones {
		preferredType, err := aws.PreferredInstanceType(ctx, meta, defaultTypes, []string{zone})
		if err != nil {
			logrus.Warnf("unable to select instanceType on the zone[%v] from the preferred list: %v. You must update the MachineSet manifest: %v", zone, defaultTypes, err)
			allZonesFound = false
			continue
		}
		if _, ok := zones[zone]; !ok {
			zones[zone] = &icaws.Zone{Name: zone}
		}
		zones[zone].PreferredInstanceType = preferredType
	}
	return allZonesFound
}

// Worker generates the machinesets for `worker` machine pool.
type Worker struct {
	UserDataFile       *asset.File
	MachineConfigFiles []*asset.File
	MachineSetFiles    []*asset.File
	MachineFiles       []*asset.File
	IPClaimFiles       []*asset.File
	IPAddrFiles        []*asset.File
}

// Name returns a human friendly name for the Worker Asset.
func (w *Worker) Name() string {
	return "Worker Machines"
}

// Dependencies returns all of the dependencies directly needed by the
// Worker asset
func (w *Worker) Dependencies() []asset.Asset {
	return []asset.Asset{
		&installconfig.ClusterID{},
		// PlatformCredsCheck just checks the creds (and asks, if needed)
		// We do not actually use it in this asset directly, hence
		// it is put in the dependencies but not fetched in Generate
		&installconfig.PlatformCredsCheck{},
		&installconfig.InstallConfig{},
		new(rhcos.Image),
		new(rhcos.Release),
		&machine.Worker{},
	}
}

// Generate generates the Worker asset.
//
//nolint:gocyclo
func (w *Worker) Generate(ctx context.Context, dependencies asset.Parents) error {
	clusterID := &installconfig.ClusterID{}
	installConfig := &installconfig.InstallConfig{}
	rhcosImage := new(rhcos.Image)
	rhcosRelease := new(rhcos.Release)
	wign := &machine.Worker{}
	dependencies.Get(clusterID, installConfig, rhcosImage, rhcosRelease, wign)

	workerUserDataSecretName := "worker-user-data"

	machines := []machinev1beta1.Machine{}
	machineConfigs := []*mcfgv1.MachineConfig{}
	machineSets := []runtime.Object{}
	var ipClaims []ipamv1.IPAddressClaim
	var ipAddrs []ipamv1.IPAddress
	var err error
	ic := installConfig.Config
	for _, pool := range ic.Compute {
		pool := pool // this makes golint happy... G601: Implicit memory aliasing in for loop. (gosec)
		if pool.Hyperthreading == types.HyperthreadingDisabled {
			ignHT, err := machineconfig.ForHyperthreadingDisabled("worker")
			if err != nil {
				return errors.Wrap(err, "failed to create ignition for hyperthreading disabled for worker machines")
			}
			machineConfigs = append(machineConfigs, ignHT)
		}
		if ic.SSHKey != "" {
			ignSSH, err := machineconfig.ForAuthorizedKeys(ic.SSHKey, "worker")
			if err != nil {
				return errors.Wrap(err, "failed to create ignition for authorized SSH keys for worker machines")
			}
			machineConfigs = append(machineConfigs, ignSSH)
		}
		if ic.FIPS {
			ignFIPS, err := machineconfig.ForFIPSEnabled("worker")
			if err != nil {
				return errors.Wrap(err, "failed to create ignition for FIPS enabled for worker machines")
			}
			machineConfigs = append(machineConfigs, ignFIPS)
		}
		if ic.Platform.Name() == powervstypes.Name {
			// always enable multipath for powervs.
			ignMultipath, err := machineconfig.ForMultipathEnabled("worker")
			if err != nil {
				return errors.Wrap(err, "failed to create ignition for multipath enabled for worker machines")
			}
			machineConfigs = append(machineConfigs, ignMultipath)

			// set SMT level if specified for powervs.
			if pool.Platform.PowerVS != nil && pool.Platform.PowerVS.SMTLevel != "" {
				ignPowerSMT, err := machineconfig.ForPowerSMT("worker", pool.Platform.PowerVS.SMTLevel)
				if err != nil {
					return errors.Wrap(err, "failed to create ignition for Power SMT for worker machines")
				}
				machineConfigs = append(machineConfigs, ignPowerSMT)
			}

			if installConfig.Config.Publish == types.InternalPublishingStrategy &&
				(len(installConfig.Config.ImageDigestSources) > 0 || len(installConfig.Config.DeprecatedImageContentSources) > 0) {
				ignChrony, err := machineconfig.ForCustomNTP("worker", powervsdefaults.DefaultNTPServer)
				if err != nil {
					return errors.Wrap(err, "failed to create ignition for custom NTP for worker machines")
				}
				machineConfigs = append(machineConfigs, ignChrony)

				ignRoutes, err := machineconfig.ForExtraRoutes("worker", powervsdefaults.DefaultExtraRoutes(), ic.MachineNetwork[0].CIDR.String())
				if err != nil {
					return errors.Wrap(err, "failed to create ignition for extra routes for worker machines")
				}
				machineConfigs = append(machineConfigs, ignRoutes)
			}
		}
		if installConfig.Config.EnabledFeatureGates().Enabled(features.FeatureGateMultiDiskSetup) {
			for i, diskSetup := range pool.DiskSetup {
				var dataDisk any
				var diskName string

				switch diskSetup.Type {
				case types.Etcd:
					diskName = diskSetup.Etcd.PlatformDiskID
				case types.Swap:
					diskName = diskSetup.Etcd.PlatformDiskID
				case types.UserDefined:
					diskName = diskSetup.UserDefined.PlatformDiskID
				default:
					// We shouldn't get here, but just in case
					return errors.Errorf("disk setup type %s is not supported", diskSetup.Type)
				}

				switch ic.Platform.Name() {
				// Each platform has their unique dataDisk type
				case azuretypes.Name:
					if i < len(pool.Platform.Azure.DataDisks) {
						dataDisk = pool.Platform.Azure.DataDisks[i]
					}
				case vspheretypes.Name:
					vsphereMachinePool := pool.Platform.VSphere
					for index, disk := range vsphereMachinePool.DataDisks {
						if disk.Name == diskName {
							dataDisk = vsphere.DiskInfo{
								Index: index,
								Disk:  disk,
							}
							break
						}
					}
				default:
					return errors.Errorf("disk setup for %s is not supported", ic.Platform.Name())
				}

				if dataDisk != nil {
					diskSetupIgn, err := NodeDiskSetup(installConfig, "worker", diskSetup, dataDisk)
					if err != nil {
						return errors.Wrap(err, "failed to create ignition to setup disks for compute")
					}
					machineConfigs = append(machineConfigs, diskSetupIgn)
				}
			}
		}
		// The maximum number of networks supported on ServiceNetwork is two, one IPv4 and one IPv6 network.
		// The cluster-network-operator handles the validation of this field.
		// Reference: https://github.com/openshift/cluster-network-operator/blob/fc3e0e25b4cfa43e14122bdcdd6d7f2585017d75/pkg/network/cluster_config.go#L45-L52
		if ic.Networking != nil && len(ic.Networking.ServiceNetwork) == 2 {
			// Only configure kernel args for dual-stack clusters.
			ignIPv6, err := machineconfig.ForDualStackAddresses("worker")
			if err != nil {
				return errors.Wrap(err, "failed to create ignition to configure IPv6 for worker machines")
			}
			machineConfigs = append(machineConfigs, ignIPv6)
		}

		switch ic.Platform.Name() {
		case awstypes.Name:
			var subnets icaws.SubnetsByZone
			switch pool.Name {
			case types.MachinePoolEdgeRoleName:
				subnets, err = aws.MachineSubnetsByZones(ctx, installConfig, awstypes.EdgeNodeSubnetRole)
				if err != nil {
					return err
				}
			default:
				subnets, err = aws.MachineSubnetsByZones(ctx, installConfig, awstypes.ClusterNodeSubnetRole)
				if err != nil {
					return err
				}
			}

			mpool := defaultAWSMachinePoolPlatform(pool.Name)

			osImage := strings.SplitN(rhcosImage.Compute, ",", 2)
			osImageID := osImage[0]
			if len(osImage) == 2 {
				osImageID = "" // the AMI will be generated later on
			}
			mpool.AMIID = osImageID

			mpool.Set(ic.Platform.AWS.DefaultMachinePlatform)
			mpool.Set(pool.Platform.AWS)
			zones := icaws.Zones{}
			zoneDefaults := false
			if len(mpool.Zones) == 0 {
				if len(subnets) > 0 {
					for _, subnet := range subnets {
						if subnet.Zone == nil {
							return errors.Wrapf(err, "failed to find zone attributes for subnet %s", subnet.ID)
						}
						mpool.Zones = append(mpool.Zones, subnet.Zone.Name)
						zones[subnet.Zone.Name] = subnets[subnet.Zone.Name].Zone
					}
				} else {
					mpool.Zones, err = installConfig.AWS.AvailabilityZones(ctx)
					if err != nil {
						return err
					}
					zoneDefaults = true
				}
			}

			// Requirements when using edge compute pools to populate machine sets.
			if pool.Name == types.MachinePoolEdgeRoleName {
				err = installConfig.AWS.SetZoneAttributes(ctx, mpool.Zones, zones)
				if err != nil {
					return errors.Wrap(err, "failed to retrieve zone attributes for edge compute pool")
				}

				if pool.Replicas == nil || *pool.Replicas == 0 {
					pool.Replicas = pointer.Int64(int64(len(mpool.Zones)))
				}
			}

			if mpool.InstanceType == "" {
				arch := installConfig.Config.ControlPlane.Architecture
				if len(installConfig.Config.Compute) > 0 {
					arch = installConfig.Config.Compute[0].Architecture
				}
				instanceTypes := awsdefaults.InstanceTypes(installConfig.Config.Platform.AWS.Region, arch, configv1.HighlyAvailableTopologyMode)
				switch pool.Name {
				case types.MachinePoolEdgeRoleName:
					if !awsSetPreferredInstanceByEdgeZone(ctx, instanceTypes, installConfig.AWS, zones) {
						// Using the default instance type from the non-edge pool often fails.
						return fmt.Errorf("failed to find instance type for one or more zones in the %s pool. Please specify an instance type in the install-config.yaml", pool.Name)
					}
				default:
					mpool.InstanceType, err = aws.PreferredInstanceType(ctx, installConfig.AWS, instanceTypes, mpool.Zones)
					if err != nil {
						logrus.Warn(errors.Wrapf(err, "failed to find default instance type for %s pool", pool.Name))
						mpool.InstanceType = instanceTypes[0]
					}
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
			sets, err := aws.MachineSets(&aws.MachineSetInput{
				ClusterID:                clusterID.InfraID,
				InstallConfigPlatformAWS: installConfig.Config.Platform.AWS,
				Subnets:                  subnets,
				Zones:                    zones,
				PublicSubnet:             awstypes.IsPublicOnlySubnetsEnabled(),
				Pool:                     &pool,
				Role:                     pool.Name,
				UserDataSecret:           workerUserDataSecretName,
			})
			if err != nil {
				return errors.Wrap(err, "failed to create worker machine objects")
			}
			for _, set := range sets {
				machineSets = append(machineSets, set)
			}
		case azuretypes.Name:
			mpool := defaultAzureMachinePoolPlatform(installConfig.Config.Platform.Azure.CloudName)
			mpool.InstanceType = azuredefaults.ComputeInstanceType(
				installConfig.Config.Platform.Azure.CloudName,
				installConfig.Config.Platform.Azure.Region,
				pool.Architecture,
			)
			mpool.Set(ic.Platform.Azure.DefaultMachinePlatform)
			mpool.Set(pool.Platform.Azure)

			client, err := installConfig.Azure.Client()
			if err != nil {
				return err
			}

			session, err := installConfig.Azure.Session()
			if err != nil {
				return err
			}

			if len(mpool.Zones) == 0 {
				azs, err := installConfig.Azure.VMAvailabilityZones(ctx, mpool.InstanceType)
				if err != nil {
					return errors.Wrap(err, "failed to fetch availability zones")
				}
				mpool.Zones = azs
				if len(azs) == 0 {
					// if no azs are given we set to []string{""} for convenience over later operations.
					// It means no-zoned for the machine API
					mpool.Zones = []string{""}
				}
			}
			subnetZones := []string{}
			if ic.Azure.OutboundType == azuretypes.NATGatewayMultiZoneOutboundType {
				subnetZones, err = installConfig.Azure.AvailabilityZones(ctx)
				if err != nil {
					return errors.Wrap(err, "failed to fetch availability zones")
				}
				computeSubnet := installConfig.Config.Azure.ComputeSubnetName(clusterID.InfraID)
				_, err := installConfig.Azure.GenerateZonesSubnetMap(installConfig.Config.Azure.Subnets, computeSubnet)
				if err != nil {
					return err
				}
			}

			if mpool.OSImage.Publisher != "" {
				img, ierr := client.GetMarketplaceImage(ctx, ic.Platform.Azure.Region, mpool.OSImage.Publisher, mpool.OSImage.Offer, mpool.OSImage.SKU, mpool.OSImage.Version)
				if ierr != nil {
					return fmt.Errorf("failed to fetch marketplace image: %w", ierr)
				}
				// Publisher is case-sensitive and matched against exactly. Also
				// the Plan's publisher might not be exactly the same as the
				// Image's publisher
				if img.Plan != nil && img.Plan.Publisher != nil {
					mpool.OSImage.Publisher = *img.Plan.Publisher
				}
			}
			pool.Platform.Azure = &mpool

			capabilities, err := installConfig.Azure.ComputeCapabilities()
			if err != nil {
				return err
			}

			sets, err := azure.MachineSets(clusterID.InfraID, installConfig, &pool, rhcosImage.Compute, "worker", workerUserDataSecretName, capabilities, subnetZones, session)
			if err != nil {
				return errors.Wrap(err, "failed to create worker machine objects")
			}
			for _, set := range sets {
				machineSets = append(machineSets, set)
			}
		case baremetaltypes.Name:
			mpool := defaultBareMetalMachinePoolPlatform()
			mpool.Set(ic.Platform.BareMetal.DefaultMachinePlatform)
			mpool.Set(pool.Platform.BareMetal)
			pool.Platform.BareMetal = &mpool

			enabledCaps := installConfig.Config.GetEnabledCapabilities()
			if enabledCaps.Has(configv1.ClusterVersionCapabilityMachineAPI) {
				// Use managed user data secret, since images used by MachineSet
				// are always up to date
				workerUserDataSecretName = "worker-user-data-managed"
				sets, err := baremetal.MachineSets(clusterID.InfraID, ic, &pool, "", "worker", workerUserDataSecretName)
				if err != nil {
					return errors.Wrap(err, "failed to create worker machine objects")
				}
				for _, set := range sets {
					machineSets = append(machineSets, set)
				}
			}
		case gcptypes.Name:
			mpool := defaultGCPMachinePoolPlatform(pool.Architecture)
			mpool.Set(ic.Platform.GCP.DefaultMachinePlatform)
			mpool.Set(pool.Platform.GCP)
			if len(mpool.Zones) == 0 {
				azs, err := gcp.ZonesForInstanceType(ic.Platform.GCP.ProjectID, ic.Platform.GCP.Region, mpool.InstanceType, ic.Platform.GCP.Endpoint)
				if err != nil {
					return errors.Wrap(err, "failed to fetch availability zones")
				}
				mpool.Zones = azs
			}
			pool.Platform.GCP = &mpool
			sets, err := gcp.MachineSets(clusterID.InfraID, ic, &pool, rhcosImage.Compute, "worker", workerUserDataSecretName)
			if err != nil {
				return errors.Wrap(err, "failed to create worker machine objects")
			}
			for _, set := range sets {
				machineSets = append(machineSets, set)
			}
		case ibmcloudtypes.Name:
			subnets := map[string]string{}
			if len(ic.Platform.IBMCloud.ComputeSubnets) > 0 {
				subnetMetas, err := installConfig.IBMCloud.ComputeSubnets(ctx)
				if err != nil {
					return err
				}
				for _, subnet := range subnetMetas {
					subnets[subnet.Zone] = subnet.Name
				}
			}
			mpool := defaultIBMCloudMachinePoolPlatform()
			mpool.Set(ic.Platform.IBMCloud.DefaultMachinePlatform)
			mpool.Set(pool.Platform.IBMCloud)
			if len(mpool.Zones) == 0 {
				azs, err := ibmcloud.AvailabilityZones(ic.Platform.IBMCloud.Region, ic.Platform.IBMCloud.ServiceEndpoints)
				if err != nil {
					return errors.Wrap(err, "failed to fetch availability zones")
				}
				mpool.Zones = azs
			}
			pool.Platform.IBMCloud = &mpool
			sets, err := ibmcloud.MachineSets(clusterID.InfraID, ic, subnets, &pool, "worker", workerUserDataSecretName)
			if err != nil {
				return errors.Wrap(err, "failed to create worker machine objects")
			}
			for _, set := range sets {
				machineSets = append(machineSets, set)
			}
		case openstacktypes.Name, powervctypes.Name:
			mpool := defaultOpenStackMachinePoolPlatform()
			mpool.Set(ic.Platform.OpenStack.DefaultMachinePlatform)
			mpool.Set(pool.Platform.OpenStack)
			pool.Platform.OpenStack = &mpool

			imageName, _ := rhcosutils.GenerateOpenStackImageName(rhcosImage.Compute, clusterID.InfraID)

			sets, err := openstack.MachineSets(ctx, clusterID.InfraID, ic, &pool, imageName, "worker", workerUserDataSecretName)
			if err != nil {
				return fmt.Errorf("failed to create worker machine objects: %w", err)
			}
			for _, set := range sets {
				machineSets = append(machineSets, set)
			}
		case vspheretypes.Name:
			mpool := defaultVSphereMachinePoolPlatform()
			mpool.Set(ic.Platform.VSphere.DefaultMachinePlatform)
			mpool.Set(pool.Platform.VSphere)
			pool.Platform.VSphere = &mpool

			sets, err := vsphere.MachineSets(clusterID.InfraID, ic, &pool, "worker", workerUserDataSecretName)
			if err != nil {
				return errors.Wrap(err, "failed to create worker machine objects")
			}
			for _, set := range sets {
				machineSets = append(machineSets, set)
			}

			// If static IPs are configured, we must generate worker machines and scale the machinesets to 0.
			if ic.Platform.VSphere.Hosts != nil {
				logrus.Debug("Generating worker machines with static IPs.")

				data, err := vsphere.Machines(clusterID.InfraID, ic, &pool, "worker", workerUserDataSecretName)
				if err != nil {
					return errors.Wrap(err, "failed to create worker machine objects")
				}

				machines = data.Machines
				ipClaims = data.IPClaims
				ipAddrs = data.IPAddresses

				logrus.Debugf("Generated %v worker machines.", len(machines))

				for _, ms := range sets {
					ms.Spec.Replicas = pointer.Int32(0)
				}
			}
		case ovirttypes.Name:
			mpool := defaultOvirtMachinePoolPlatform()
			mpool.Set(ic.Platform.Ovirt.DefaultMachinePlatform)
			mpool.Set(pool.Platform.Ovirt)
			pool.Platform.Ovirt = &mpool

			imageName, _ := rhcosutils.GenerateOpenStackImageName(rhcosImage.Compute, clusterID.InfraID)

			sets, err := ovirt.MachineSets(clusterID.InfraID, ic, &pool, imageName, "worker", workerUserDataSecretName)
			if err != nil {
				return errors.Wrap(err, "failed to create worker machine objects for ovirt provider")
			}
			for _, set := range sets {
				machineSets = append(machineSets, set)
			}
		case powervstypes.Name:
			mpool := defaultPowerVSMachinePoolPlatform(ic)
			mpool.Set(ic.Platform.PowerVS.DefaultMachinePlatform)
			mpool.Set(pool.Platform.PowerVS)
			pool.Platform.PowerVS = &mpool
			sets, err := powervs.MachineSets(clusterID.InfraID, ic, &pool, "worker", "worker-user-data")
			if err != nil {
				return errors.Wrap(err, "failed to create worker machine objects for powervs provider")
			}
			for _, set := range sets {
				machineSets = append(machineSets, set)
			}
		case externaltypes.Name, nonetypes.Name:
		case nutanixtypes.Name:
			mpool := defaultNutanixMachinePoolPlatform()
			mpool.Set(ic.Platform.Nutanix.DefaultMachinePlatform)
			mpool.Set(pool.Platform.Nutanix)
			if err = mpool.ValidateConfig(ic.Platform.Nutanix, "worker"); err != nil {
				return errors.Wrap(err, "failed to create worker machine objects")
			}
			pool.Platform.Nutanix = &mpool
			imageName := nutanixtypes.RHCOSImageName(ic.Platform.Nutanix, clusterID.InfraID)

			sets, err := nutanix.MachineSets(clusterID.InfraID, ic, &pool, imageName, "worker", workerUserDataSecretName)
			if err != nil {
				return errors.Wrap(err, "failed to create worker machine objects")
			}
			for _, set := range sets {
				machineSets = append(machineSets, set)
			}
		default:
			return fmt.Errorf("invalid Platform")
		}
	}

	data, err := UserDataSecret(workerUserDataSecretName, wign.File.Data)
	if err != nil {
		return errors.Wrap(err, "failed to create user-data secret for worker machines")
	}
	w.UserDataFile = &asset.File{
		Filename: filepath.Join(directory, workerUserDataFileName),
		Data:     data,
	}

	w.MachineConfigFiles, err = machineconfig.Manifests(machineConfigs, "worker", directory)
	if err != nil {
		return errors.Wrap(err, "failed to create MachineConfig manifests for worker machines")
	}

	w.MachineSetFiles = make([]*asset.File, len(machineSets))
	padFormat := fmt.Sprintf("%%0%dd", len(fmt.Sprintf("%d", len(machineSets))))
	for i, machineSet := range machineSets {
		data, err := yaml.Marshal(machineSet)
		if err != nil {
			return errors.Wrapf(err, "marshal worker %d", i)
		}

		padded := fmt.Sprintf(padFormat, i)
		w.MachineSetFiles[i] = &asset.File{
			Filename: filepath.Join(directory, fmt.Sprintf(workerMachineSetFileName, padded)),
			Data:     data,
		}
	}

	w.IPClaimFiles = make([]*asset.File, len(ipClaims))
	for i, claim := range ipClaims {
		data, err := yaml.Marshal(claim)
		if err != nil {
			return errors.Wrapf(err, "marshal ip claim %v", claim.Name)
		}

		w.IPClaimFiles[i] = &asset.File{
			Filename: filepath.Join(directory, fmt.Sprintf(ipClaimFileName, claim.Name)),
			Data:     data,
		}
	}

	w.IPAddrFiles = make([]*asset.File, len(ipAddrs))
	for i, address := range ipAddrs {
		data, err := yaml.Marshal(address)
		if err != nil {
			return errors.Wrapf(err, "marshal ip claim %v", address.Name)
		}

		w.IPAddrFiles[i] = &asset.File{
			Filename: filepath.Join(directory, fmt.Sprintf(ipAddressFileName, address.Name)),
			Data:     data,
		}
	}
	w.MachineFiles = make([]*asset.File, len(machines))
	for i, machineDef := range machines {
		data, err := yaml.Marshal(machineDef)
		if err != nil {
			return errors.Wrapf(err, "marshal master %d", i)
		}

		padded := fmt.Sprintf(padFormat, i)
		w.MachineFiles[i] = &asset.File{
			Filename: filepath.Join(directory, fmt.Sprintf(workerMachineFileName, padded)),
			Data:     data,
		}
	}
	return nil
}

// Files returns the files generated by the asset.
func (w *Worker) Files() []*asset.File {
	files := make([]*asset.File, 0, 1+len(w.MachineConfigFiles)+len(w.MachineSetFiles))
	if w.UserDataFile != nil {
		files = append(files, w.UserDataFile)
	}
	files = append(files, w.MachineConfigFiles...)
	files = append(files, w.MachineSetFiles...)
	files = append(files, w.MachineFiles...)
	files = append(files, w.IPClaimFiles...)
	files = append(files, w.IPAddrFiles...)
	return files
}

// Load reads the asset files from disk.
func (w *Worker) Load(f asset.FileFetcher) (found bool, err error) {
	file, err := f.FetchByName(filepath.Join(directory, workerUserDataFileName))
	if err != nil {
		if os.IsNotExist(err) {
			return false, nil
		}
		return false, err
	}
	w.UserDataFile = file

	w.MachineConfigFiles, err = machineconfig.Load(f, "worker", directory)
	if err != nil {
		return true, err
	}

	fileList, err := f.FetchByPattern(filepath.Join(directory, workerMachineSetFileNamePattern))
	if err != nil {
		return true, err
	}

	w.MachineSetFiles = fileList

	fileList, err = f.FetchByPattern(filepath.Join(directory, workerMachineFileNamePattern))
	if err != nil {
		return true, err
	}
	w.MachineFiles = fileList

	fileList, err = f.FetchByPattern(filepath.Join(directory, workerIPClaimFileNamePattern))
	if err != nil {
		return true, err
	}
	w.IPClaimFiles = fileList

	fileList, err = f.FetchByPattern(filepath.Join(directory, workerIPAddressFileNamePattern))
	if err != nil {
		return true, err
	}
	w.IPAddrFiles = fileList

	return true, nil
}

// MachineSets returns MachineSet manifest structures.
func (w *Worker) MachineSets() ([]machinev1beta1.MachineSet, error) {
	scheme := runtime.NewScheme()
	baremetalapi.AddToScheme(scheme)
	ibmcloudapi.AddToScheme(scheme)
	libvirtapi.AddToScheme(scheme)
	ovirtproviderapi.AddToScheme(scheme)
	scheme.AddKnownTypes(machinev1alpha1.GroupVersion,
		&machinev1alpha1.OpenstackProviderSpec{},
	)
	scheme.AddKnownTypes(machinev1beta1.SchemeGroupVersion,
		&machinev1beta1.AWSMachineProviderConfig{},
		&machinev1beta1.VSphereMachineProviderSpec{},
		&machinev1beta1.AzureMachineProviderSpec{},
		&machinev1beta1.GCPMachineProviderSpec{},
	)
	machinev1.Install(scheme)
	scheme.AddKnownTypes(machinev1.GroupVersion,
		&machinev1.NutanixMachineProviderConfig{},
		&machinev1.PowerVSMachineProviderConfig{},
	)
	machinev1beta1.AddToScheme(scheme)
	decoder := serializer.NewCodecFactory(scheme).UniversalDecoder(
		baremetalprovider.SchemeGroupVersion,
		ibmcloudprovider.SchemeGroupVersion,
		libvirtprovider.SchemeGroupVersion,
		machinev1.GroupVersion,
		machinev1alpha1.GroupVersion,
		ovirtprovider.SchemeGroupVersion,
		machinev1beta1.SchemeGroupVersion,
	)

	machineSets := []machinev1beta1.MachineSet{}
	for i, file := range w.MachineSetFiles {
		machineSet := &machinev1beta1.MachineSet{}
		err := yaml.Unmarshal(file.Data, &machineSet)
		if err != nil {
			return machineSets, errors.Wrapf(err, "unmarshal worker %d", i)
		}

		obj, _, err := decoder.Decode(machineSet.Spec.Template.Spec.ProviderSpec.Value.Raw, nil, nil)
		if err != nil {
			return machineSets, errors.Wrapf(err, "unmarshal worker %d", i)
		}

		machineSet.Spec.Template.Spec.ProviderSpec.Value = &runtime.RawExtension{Object: obj}
		machineSets = append(machineSets, *machineSet)
	}

	return machineSets, nil
}
