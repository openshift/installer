package manifests

import (
	"context"
	"fmt"
	"path/filepath"
	"strings"

	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
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
	dnstypes "github.com/openshift/installer/pkg/types/dns"
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
func (d *DNS) Generate(ctx context.Context, dependencies asset.Parents) error { //nolint:gocyclo
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
		// We do not want to configure cloud DNS when `UserProvisionedDNS` is enabled.
		// So, do not set PrivateZone and PublicZone fields in the DNS manifest.
		if installConfig.Config.AWS.UserProvisionedDNS == dnstypes.UserProvisionedDNSEnabled {
			config.Spec.PublicZone = nil
			config.Spec.PrivateZone = nil
			break
		}
		if installConfig.Config.Publish == types.ExternalPublishingStrategy {
			client, err := icaws.NewRoute53Client(ctx, icaws.EndpointOptions{
				Region:    installConfig.Config.AWS.Region,
				Endpoints: installConfig.Config.AWS.ServiceEndpoints,
			}, "")
			if err != nil {
				return fmt.Errorf("failed to create route 53 client: %w", err)
			}

			zone, err := icaws.GetPublicZone(ctx, client, installConfig.Config.BaseDomain)
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
			// We do not want to configure cloud DNS when `UserProvisionedDNS` is enabled.
			// So, do not set PrivateZone and PublicZone fields in the DNS manifest.
			if installConfig.Config.Azure.UserProvisionedDNS == dnstypes.UserProvisionedDNSEnabled {
				config.Spec.PublicZone = &configv1.DNSZone{ID: ""}
				config.Spec.PrivateZone = &configv1.DNSZone{ID: ""}
			}
		}
	case gcptypes.Name:
		// We do not want to configure cloud DNS when `UserProvisionedDNS` is enabled.
		// So, do not set PrivateZone and PublicZone fields in the DNS manifest.
		if installConfig.Config.GCP.UserProvisionedDNS == dnstypes.UserProvisionedDNSEnabled {
			config.Spec.PublicZone = nil
			config.Spec.PrivateZone = nil
			break
		}
		client, err := icgcp.NewClient(context.Background(), installConfig.Config.GCP.ServiceEndpoints)
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

			publicZoneName := fmt.Sprintf("projects/%s/managedZones/%s", installConfig.Config.GCP.ProjectID, zone.Name)
			logrus.Infof("generating GCP Public DNS Zone %s", publicZoneName)
			config.Spec.PublicZone = &configv1.DNSZone{ID: publicZoneName}
		}

		// Ingress operator can handle a zone with the following format:
		// projects/{projectID}/managedZones/{zoneID}. This will allow
		// the installer to pass the project without a new field in the
		// DNSZone struct.
		params, err := GetGCPPrivateZoneInfo(ctx, client, installConfig, clusterID.InfraID)
		if err != nil {
			return fmt.Errorf("failed to get private zone info: %w", err)
		}

		privateZoneName := fmt.Sprintf("projects/%s/managedZones/%s", params.Project, params.Name)
		logrus.Infof("generating GCP Private DNS Zone %s", privateZoneName)
		config.Spec.PrivateZone = &configv1.DNSZone{ID: privateZoneName}

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

// GCPNetworkName create the full resource name for a network.
func GCPNetworkName(project, network string) string {
	return fmt.Sprintf("https://www.googleapis.com/compute/v1/projects/%s/global/networks/%s", project, network)
}

// GCPDefaultPrivateZoneID returns the default name for a gcp private dns zone. This zone name will be used during
// installations where the user has not provided a private zone name (xpn installs only), no
// preexisting private dns zone is found (xpn installs only), and default installation cases.
func GCPDefaultPrivateZoneID(clusterID string) string {
	return fmt.Sprintf("%s-private-zone", clusterID)
}

// GetGCPPrivateZoneInfo attempts to find the name of the private zone for GCP installs. When a shared vpc install
// occurs, a precreated zone may be used. If a zone is found (in this instance), then the zone should be paired with
// the network that is supplied through the install config (when applicable).
func GetGCPPrivateZoneInfo(ctx context.Context, client *icgcp.Client, installConfig *installconfig.InstallConfig, clusterID string) (gcptypes.DNSZoneParams, error) {
	params := gcptypes.DNSZoneParams{
		// Force set the private zone ID to an empty string to ensure
		// the search for DNS zones looks for ANY not a specific zone.
		// This is required, because the user may enter no zone information
		// but still wish to bring a private zone during xpn installs (in this
		// case it must exist in the `projectID`).
		Name:             "",
		InstallerCreated: true,
		IsPublic:         false,
		BaseDomain:       installConfig.Config.ClusterDomain(),
		Project:          installConfig.Config.GCP.ProjectID,
	}

	if installConfig.Config.GCP.NetworkProjectID != "" && installConfig.Config.GCP.Network != "" {
		icdns := installConfig.Config.GCP.DNS
		if icdns != nil && icdns.PrivateZone != nil {
			if icdns.PrivateZone.ProjectID != "" {
				params.Project = icdns.PrivateZone.ProjectID
			}
			// Override the default with the name provided. If this zone does not exist, then
			// this should still be returned as valid.
			params.Name = icdns.PrivateZone.Name
		}

		zone, err := client.GetDNSZoneFromParams(ctx, params)
		if err != nil {
			// Currently, the only time that a private zone lookup will produce an error is if we
			// failed to find the dns zones. That should result in an error returned here too.
			return params, fmt.Errorf("private dns zone %s does not exist or is invalid: %w", params.Name, err)
		}
		if zone == nil {
			// CORS-4012: The user may specify a zone to be created if it does not exist.
			// Do not fail if the specified zone does not exist.
			if params.Name == "" {
				params.Name = GCPDefaultPrivateZoneID(clusterID)
			}
			return params, nil
		}

		params.Name = zone.Name
		params.InstallerCreated = false
	} else {
		params.Name = GCPDefaultPrivateZoneID(clusterID)
	}
	return params, nil
}

// Files returns the files generated by the asset.
func (d *DNS) Files() []*asset.File {
	return d.FileList
}

// Load loads the already-rendered files back from disk.
func (d *DNS) Load(f asset.FileFetcher) (bool, error) {
	return false, nil
}
