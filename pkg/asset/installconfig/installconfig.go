package installconfig

import (
	"context"
	"fmt"

	"github.com/pkg/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/validation/field"

	"github.com/openshift/installer/pkg/asset"
	"github.com/openshift/installer/pkg/asset/installconfig/aws"
	icazure "github.com/openshift/installer/pkg/asset/installconfig/azure"
	icgcp "github.com/openshift/installer/pkg/asset/installconfig/gcp"
	icibmcloud "github.com/openshift/installer/pkg/asset/installconfig/ibmcloud"
	icnutanix "github.com/openshift/installer/pkg/asset/installconfig/nutanix"
	icopenstack "github.com/openshift/installer/pkg/asset/installconfig/openstack"
	icovirt "github.com/openshift/installer/pkg/asset/installconfig/ovirt"
	icpowervs "github.com/openshift/installer/pkg/asset/installconfig/powervs"
	icvsphere "github.com/openshift/installer/pkg/asset/installconfig/vsphere"
	"github.com/openshift/installer/pkg/types"
	"github.com/openshift/installer/pkg/types/defaults"
	"github.com/openshift/installer/pkg/types/validation"
)

const (
	installConfigFilename = "install-config.yaml"
)

// InstallConfig generates the install-config.yaml file.
type InstallConfig struct {
	AssetBase
	AWS      *aws.Metadata        `json:"aws,omitempty"`
	Azure    *icazure.Metadata    `json:"azure,omitempty"`
	IBMCloud *icibmcloud.Metadata `json:"ibmcloud,omitempty"`
	PowerVS  *icpowervs.Metadata  `json:"powervs,omitempty"`
	VSphere  *icvsphere.Metadata  `json:"vsphere,omitempty"`
}

var _ asset.WritableAsset = (*InstallConfig)(nil)

// MakeAsset returns an InstallConfig asset containing a given InstallConfig CR.
func MakeAsset(config *types.InstallConfig) *InstallConfig {
	return &InstallConfig{
		AssetBase: AssetBase{
			Config: config,
		},
	}
}

// Dependencies returns all of the dependencies directly needed by an
// InstallConfig asset.
func (a *InstallConfig) Dependencies() []asset.Asset {
	return []asset.Asset{
		&sshPublicKey{},
		&baseDomain{},
		&clusterName{},
		&pullSecret{},
		&platform{},
	}
}

// Generate the install-config.yaml file.
func (a *InstallConfig) Generate(ctx context.Context, parents asset.Parents) error {
	sshPublicKey := &sshPublicKey{}
	baseDomain := &baseDomain{}
	clusterName := &clusterName{}
	pullSecret := &pullSecret{}
	platform := &platform{}
	parents.Get(
		sshPublicKey,
		baseDomain,
		clusterName,
		pullSecret,
		platform,
	)

	a.Config = &types.InstallConfig{
		TypeMeta: metav1.TypeMeta{
			APIVersion: types.InstallConfigVersion,
		},
		ObjectMeta: metav1.ObjectMeta{
			Name: clusterName.ClusterName,
		},
		SSHKey:     sshPublicKey.Key,
		BaseDomain: baseDomain.BaseDomain,
		Publish:    baseDomain.Publish,
		PullSecret: pullSecret.PullSecret,
	}

	a.Config.AWS = platform.AWS
	a.Config.None = platform.None
	a.Config.OpenStack = platform.OpenStack
	a.Config.VSphere = platform.VSphere
	a.Config.Azure = platform.Azure
	a.Config.GCP = platform.GCP
	a.Config.IBMCloud = platform.IBMCloud
	a.Config.BareMetal = platform.BareMetal
	a.Config.Ovirt = platform.Ovirt
	a.Config.PowerVS = platform.PowerVS
	a.Config.Nutanix = platform.Nutanix

	defaults.SetInstallConfigDefaults(a.Config)

	return a.finish(ctx, "")
}

// Load returns the installconfig from disk.
func (a *InstallConfig) Load(f asset.FileFetcher) (found bool, err error) {
	ctx := context.TODO()
	found, err = a.LoadFromFile(f)
	if found && err == nil {
		if err := a.finish(ctx, installConfigFilename); err != nil {
			return false, errors.Wrap(err, asset.InstallConfigError)
		}
	}

	return found, err
}

// finishGCP will set default values in the install config that require api calls rather than static checks.
func (a *InstallConfig) finishGCP() error {
	if endpoint := a.Config.Platform.GCP.Endpoint; endpoint != nil && endpoint.ClusterUseOnly == nil {
		client, err := icgcp.NewClient(context.TODO(), nil)
		if err != nil {
			return err
		}
		defaultClusterUseOnly := true
		if _, err := client.GetRegions(context.TODO(), a.Config.Platform.GCP.ProjectID); err != nil {
			defaultClusterUseOnly = false
		}
		a.Config.Platform.GCP.Endpoint.ClusterUseOnly = &defaultClusterUseOnly
	}
	return nil
}

// finishAWS set defaults for AWS Platform before the config validation.
func (a *InstallConfig) finishAWS() error {
	// Set the Default Edge Compute pool when the subnets in AWS Local Zones are defined,
	// when installing a cluster in existing VPC.
	if len(a.Config.Platform.AWS.VPC.Subnets) > 0 {
		edgeSubnets, err := a.AWS.EdgeSubnets(context.TODO())
		if err != nil {
			return errors.Wrap(err, fmt.Sprintf("unable to load edge subnets: %v", err))
		}
		totalEdgeSubnets := int64(len(edgeSubnets))
		if totalEdgeSubnets == 0 {
			return nil
		}
		if edgePool := defaults.CreateEdgeMachinePoolDefaults(a.Config.Compute, &a.Config.Platform, totalEdgeSubnets); edgePool != nil {
			a.Config.Compute = append(a.Config.Compute, *edgePool)
		}
	}
	return nil
}

func (a *InstallConfig) finish(ctx context.Context, filename string) error {
	if a.Config.AWS != nil {
		a.AWS = aws.NewMetadata(a.Config.Platform.AWS.Region, a.Config.Platform.AWS.VPC.Subnets, a.Config.AWS.ServiceEndpoints)
		if err := a.finishAWS(); err != nil {
			return err
		}
	}
	if a.Config.Azure != nil {
		a.Azure = icazure.NewMetadata(a.Config.Azure.CloudName, a.Config.Azure.ARMEndpoint)
	}
	if a.Config.GCP != nil {
		if err := a.finishGCP(); err != nil {
			return err
		}
	}
	if a.Config.IBMCloud != nil {
		a.IBMCloud = icibmcloud.NewMetadata(a.Config)
	}
	if a.Config.PowerVS != nil {
		a.PowerVS = icpowervs.NewMetadata(a.Config)
	}
	if a.Config.VSphere != nil {
		a.VSphere = icvsphere.NewMetadata()

		for _, v := range a.Config.VSphere.VCenters {
			_ = a.VSphere.AddCredentials(v.Server, v.Username, v.Password)
		}
	}

	if err := validation.ValidateInstallConfig(a.Config, false).ToAggregate(); err != nil {
		if filename == "" {
			return errors.Wrap(err, "invalid install config")
		}
		return errors.Wrapf(err, "invalid %q file", filename)
	}

	if err := a.platformValidation(ctx); err != nil {
		return err
	}
	return a.RecordFile()
}

// platformValidation runs validations that require connecting to the
// underlying platform. In some cases, platforms also duplicate validations
// that have already been checked by validation.ValidateInstallConfig().
func (a *InstallConfig) platformValidation(ctx context.Context) error {
	if a.Config.Platform.Azure != nil {
		client, err := a.Azure.Client()
		if err != nil {
			return err
		}
		return icazure.Validate(client, a.Config)
	}
	if a.Config.Platform.GCP != nil {
		client, err := icgcp.NewClient(ctx, a.Config.GCP.Endpoint)
		if err != nil {
			return err
		}
		return icgcp.Validate(client, a.Config)
	}
	if a.Config.Platform.IBMCloud != nil {
		// Validate the Service Endpoints now, before performing any additional validation of the InstallConfig
		err := icibmcloud.ValidateServiceEndpoints(a.Config)
		if err != nil {
			return err
		}
		client, err := icibmcloud.NewClient(a.Config.Platform.IBMCloud.ServiceEndpoints)
		if err != nil {
			return err
		}
		return icibmcloud.Validate(client, a.Config)
	}
	if a.Config.Platform.AWS != nil {
		return aws.Validate(ctx, a.AWS, a.Config)
	}
	if a.Config.Platform.VSphere != nil {
		return icvsphere.Validate(a.Config)
	}
	if a.Config.Platform.Ovirt != nil {
		return icovirt.Validate(a.Config)
	}
	if a.Config.Platform.OpenStack != nil {
		return icopenstack.Validate(ctx, a.Config)
	}
	if a.Config.Platform.PowerVS != nil {
		return icpowervs.Validate(a.Config)
	}
	if a.Config.Platform.Nutanix != nil {
		return icnutanix.Validate(a.Config)
	}
	return field.ErrorList{}.ToAggregate()
}
