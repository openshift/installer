package manifests

import (
	"context"
	"fmt"
	"path/filepath"
	"strings"

	"github.com/pkg/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/yaml"

	configv1 "github.com/openshift/api/config/v1"
	"github.com/openshift/installer/pkg/asset"
	"github.com/openshift/installer/pkg/asset/installconfig"
	icaws "github.com/openshift/installer/pkg/asset/installconfig/aws"
	icgcp "github.com/openshift/installer/pkg/asset/installconfig/gcp"
	icibmcloud "github.com/openshift/installer/pkg/asset/installconfig/ibmcloud"
	icpowervs "github.com/openshift/installer/pkg/asset/installconfig/powervs"
	"github.com/openshift/installer/pkg/types"
	awstypes "github.com/openshift/installer/pkg/types/aws"
	azuretypes "github.com/openshift/installer/pkg/types/azure"
	baremetaltypes "github.com/openshift/installer/pkg/types/baremetal"
	externaltypes "github.com/openshift/installer/pkg/types/external"
	gcptypes "github.com/openshift/installer/pkg/types/gcp"
	ibmcloudtypes "github.com/openshift/installer/pkg/types/ibmcloud"
	nonetypes "github.com/openshift/installer/pkg/types/none"
	nutanixtypes "github.com/openshift/installer/pkg/types/nutanix"
	openstacktypes "github.com/openshift/installer/pkg/types/openstack"
	ovirttypes "github.com/openshift/installer/pkg/types/ovirt"
	powervstypes "github.com/openshift/installer/pkg/types/powervs"
	vspheretypes "github.com/openshift/installer/pkg/types/vsphere"
)

var (
	dnsCfgFilename = filepath.Join(manifestDir, "cluster-dns-02-config.yml")

	combineGCPZoneInfo = func(project, zoneName string) string {
		return fmt.Sprintf("project/%s/managedZones/%s", project, zoneName)
	}
)

// DNS generates the cluster-dns-*.yml files.
type DNS struct {
	FileList []*asset.File
}

var _ asset.WritableAsset = (*DNS)(nil)

// Name returns a human friendly name for the asset.
func (*DNS) Name() string {
	return "DNS Config"
}

// Dependencies returns all of the dependencies directly needed to generate
// the asset.
func (*DNS) Dependencies() []asset.Asset {
	return []asset.Asset{
		&installconfig.InstallConfig{},
		&installconfig.ClusterID{},
		// PlatformCredsCheck just checks the creds (and asks, if needed)
		// We do not actually use it in this asset directly, hence
		// it is put in the dependencies but not fetched in Generate
		&installconfig.PlatformCredsCheck{},
	}
}

// Generate generates the DNS config and its CRD.
func (d *DNS) Generate(ctx context.Context, dependencies asset.Parents) error {
	installConfig := &installconfig.InstallConfig{}
	clusterID := &installconfig.ClusterID{}
	dependencies.Get(installConfig, clusterID)

	config := &configv1.DNS{
		TypeMeta: metav1.TypeMeta{
			APIVersion: configv1.SchemeGroupVersion.String(),
			Kind:       "DNS",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name: "cluster",
			// not namespaced
		},
		Spec: configv1.DNSSpec{
			BaseDomain: installConfig.Config.ClusterDomain(),
		},
	}

	switch installConfig.Config.Platform.Name() {
	case awstypes.Name:
		if installConfig.Config.Publish == types.ExternalPublishingStrategy {
			sess, err := installConfig.AWS.Session(ctx)
			if err != nil {
				return errors.Wrap(err, "failed to initialize session")
			}
			zone, err := icaws.GetPublicZone(sess, installConfig.Config.BaseDomain)
			if err != nil {
				return errors.Wrapf(err, "getting public zone for %q", installConfig.Config.BaseDomain)
			}
			config.Spec.PublicZone = &configv1.DNSZone{ID: strings.TrimPrefix(*zone.Id, "/hostedzone/")}
		}
		if hostedZone := installConfig.Config.AWS.HostedZone; hostedZone == "" {
			config.Spec.PrivateZone = &configv1.DNSZone{Tags: map[string]string{
				fmt.Sprintf("kubernetes.io/cluster/%s", clusterID.InfraID): "owned",
				"Name": fmt.Sprintf("%s-int", clusterID.InfraID),
			}}
		} else {
			config.Spec.PrivateZone = &configv1.DNSZone{ID: hostedZone}

			if r := installConfig.Config.AWS.HostedZoneRole; r != "" {
				config.Spec.Platform = configv1.DNSPlatformSpec{
					Type: configv1.AWSPlatformType,
					AWS: &configv1.AWSDNSSpec{
						PrivateZoneIAMRole: r,
					},
				}
			}
		}
	case azuretypes.Name:
		dnsConfig, err := installConfig.Azure.DNSConfig()
		if err != nil {
			return err
		}

		if installConfig.Config.Publish == types.ExternalPublishingStrategy ||
			(installConfig.Config.Publish == types.MixedPublishingStrategy && installConfig.Config.OperatorPublishingStrategy.Ingress != "Internal") {
			//currently, this guesses the azure resource IDs from known parameter.
			config.Spec.PublicZone = &configv1.DNSZone{
				ID: dnsConfig.GetDNSZoneID(installConfig.Config.Azure.BaseDomainResourceGroupName, installConfig.Config.BaseDomain),
			}
		}
		if installConfig.Azure.CloudName != azuretypes.StackCloud {
			config.Spec.PrivateZone = &configv1.DNSZone{
				ID: dnsConfig.GetPrivateDNSZoneID(installConfig.Config.Azure.ClusterResourceGroupName(clusterID.InfraID), installConfig.Config.ClusterDomain()),
			}
		}
	case gcptypes.Name:
		// We donot want to configure cloud DNS when `UserProvisionedDNS` is enabled.
		// So, do not set PrivateZone and PublicZone fields in the DNS manifest.
		if installConfig.Config.GCP.UserProvisionedDNS == gcptypes.UserProvisionedDNSEnabled {
			config.Spec.PublicZone = &configv1.DNSZone{ID: ""}
			config.Spec.PrivateZone = &configv1.DNSZone{ID: ""}
			break
		}
		client, err := icgcp.NewClient(context.Background())
		if err != nil {
			return err
		}

		// Set the public zone
		switch {
		case installConfig.Config.Publish != types.ExternalPublishingStrategy:
			// Do not use a public zone when not publishing externally.
		default:
			// Search the project for a zone with the specified base domain.
			zone, err := client.GetDNSZone(ctx, installConfig.Config.GCP.ProjectID, installConfig.Config.BaseDomain, true)
			if err != nil {
				return errors.Wrapf(err, "failed to get public zone for %q", installConfig.Config.BaseDomain)
			}
			config.Spec.PublicZone = &configv1.DNSZone{ID: zone.Name}
		}

		// Set the private zone
		privateZoneID := fmt.Sprintf("%s-private-zone", clusterID.InfraID)
		zone, err := client.GetDNSZone(ctx, installConfig.Config.GCP.ProjectID, installConfig.Config.ClusterDomain(), false)
		if err != nil {
			return errors.Wrapf(err, "failed to get private zone for %q", installConfig.Config.BaseDomain)
		}
		if zone != nil {
			privateZoneID = zone.Name
		}
		config.Spec.PrivateZone = &configv1.DNSZone{ID: privateZoneID}

	case ibmcloudtypes.Name:
		client, err := icibmcloud.NewClient(installConfig.Config.Platform.IBMCloud.ServiceEndpoints)
		if err != nil {
			return errors.Wrap(err, "failed to get IBM Cloud client")
		}

		zoneID, err := client.GetDNSZoneIDByName(ctx, installConfig.Config.BaseDomain, installConfig.Config.Publish)
		if err != nil {
			return errors.Wrap(err, "failed to get DNS zone ID")
		}

		if installConfig.Config.Publish == types.ExternalPublishingStrategy {
			config.Spec.PublicZone = &configv1.DNSZone{
				ID: zoneID,
			}
		}
		config.Spec.PrivateZone = &configv1.DNSZone{
			ID: zoneID,
		}
	case powervstypes.Name:
		client, err := icpowervs.NewClient()
		if err != nil {
			return errors.Wrap(err, "failed to get IBM PowerVS client")
		}

		zoneID, err := client.GetDNSZoneIDByName(ctx, installConfig.Config.BaseDomain, installConfig.Config.Publish)
		if err != nil {
			return errors.Wrap(err, "failed to get DNS zone ID")
		}

		if installConfig.Config.Publish == types.ExternalPublishingStrategy {
			config.Spec.PublicZone = &configv1.DNSZone{
				ID: zoneID,
			}
		}
		config.Spec.PrivateZone = &configv1.DNSZone{
			ID: zoneID,
		}
	case openstacktypes.Name, baremetaltypes.Name, externaltypes.Name, nonetypes.Name, vspheretypes.Name, ovirttypes.Name, nutanixtypes.Name:
	default:
		return errors.New("invalid Platform")
	}

	configData, err := yaml.Marshal(config)
	if err != nil {
		return errors.Wrapf(err, "failed to create %s manifests from InstallConfig", d.Name())
	}

	d.FileList = []*asset.File{
		{
			Filename: dnsCfgFilename,
			Data:     configData,
		},
	}

	return nil
}

// Files returns the files generated by the asset.
func (d *DNS) Files() []*asset.File {
	return d.FileList
}

// Load loads the already-rendered files back from disk.
func (d *DNS) Load(f asset.FileFetcher) (bool, error) {
	return false, nil
}
