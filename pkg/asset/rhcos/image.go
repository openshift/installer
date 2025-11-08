// Package rhcos contains assets for RHCOS.
package rhcos

import (
	"context"
	"fmt"
	"net/url"
	"os"
	"time"

	"github.com/coreos/stream-metadata-go/arch"
	"github.com/sirupsen/logrus"

	"github.com/openshift/installer/pkg/asset"
	"github.com/openshift/installer/pkg/asset/installconfig"
	icazure "github.com/openshift/installer/pkg/asset/installconfig/azure"
	"github.com/openshift/installer/pkg/rhcos"
	"github.com/openshift/installer/pkg/types"
	"github.com/openshift/installer/pkg/types/aws"
	"github.com/openshift/installer/pkg/types/azure"
	"github.com/openshift/installer/pkg/types/baremetal"
	"github.com/openshift/installer/pkg/types/external"
	"github.com/openshift/installer/pkg/types/gcp"
	"github.com/openshift/installer/pkg/types/ibmcloud"
	"github.com/openshift/installer/pkg/types/none"
	"github.com/openshift/installer/pkg/types/nutanix"
	"github.com/openshift/installer/pkg/types/openstack"
	"github.com/openshift/installer/pkg/types/ovirt"
	"github.com/openshift/installer/pkg/types/powervc"
	"github.com/openshift/installer/pkg/types/powervs"
	"github.com/openshift/installer/pkg/types/vsphere"
)

// Image is location of RHCOS image.
// This stores the location of the image based on the platform.
// eg. on AWS this contains ami-id, on Livirt this can be the URI for QEMU image etc.
type Image struct {
	ControlPlane string
	Compute      string
}

var _ asset.Asset = (*Image)(nil)

// Name returns the human-friendly name of the asset.
func (i *Image) Name() string {
	return "Image"
}

// Dependencies returns dependencies used by the RHCOS asset.
func (i *Image) Dependencies() []asset.Asset {
	return []asset.Asset{
		&installconfig.InstallConfig{},
	}
}

// Generate the RHCOS image location.
func (i *Image) Generate(ctx context.Context, p asset.Parents) error {
	if oi, ok := os.LookupEnv("OPENSHIFT_INSTALL_OS_IMAGE_OVERRIDE"); ok && oi != "" {
		logrus.Warn("Found override for OS Image. Please be warned, this is not advised")
		*i = *MakeAsset(oi)
		return nil
	}

	ic := &installconfig.InstallConfig{}
	p.Get(ic)
	config := ic.Config
	osimageControlPlane, err := osImage(ctx, ic, config.ControlPlane)
	if err != nil {
		return err
	}
	var computePool *types.MachinePool
	if len(config.Compute) > 0 {
		computePool = &config.Compute[0]
	} else {
		computePool = config.ControlPlane
	}
	osimageCompute, err := osImage(ctx, ic, computePool)
	if err != nil {
		return err
	}
	*i = Image{osimageControlPlane, osimageCompute}
	return nil
}

//nolint:gocyclo
func osImage(ctx context.Context, ic *installconfig.InstallConfig, machinePool *types.MachinePool) (string, error) {
	ctx, cancel := context.WithTimeout(ctx, 30*time.Second)
	defer cancel()

	nodeArch := machinePool.Architecture
	archName := arch.RpmArch(string(nodeArch))

	st, err := rhcos.FetchCoreOSBuild(ctx)
	if err != nil {
		return "", err
	}
	streamArch, err := st.GetArchitecture(archName)
	if err != nil {
		return "", err
	}
	streamArchPrefix := st.FormatPrefix(archName)
	platform := ic.Config.Platform
	switch platform.Name() {
	case aws.Name:
		region := platform.AWS.Region
		if !rhcos.AMIRegions(nodeArch).Has(region) {
			const globalResourceRegion = "us-east-1"
			logrus.Debugf("No AMI found in %s. Using AMI from %s.", region, globalResourceRegion)
			region = globalResourceRegion
		}
		osimage, err := st.GetAMI(archName, region)
		if err != nil {
			return "", err
		}
		if region != platform.AWS.Region {
			osimage = fmt.Sprintf("%s,%s", osimage, region)
		}
		return osimage, nil
	case gcp.Name:
		if streamArch.Images.Gcp != nil {
			img := streamArch.Images.Gcp
			return fmt.Sprintf("projects/%s/global/images/%s", img.Project, img.Name), nil
		}
		return "", fmt.Errorf("%s: No GCP build found", streamArchPrefix)
	case ibmcloud.Name:
		if a, ok := streamArch.Artifacts["ibmcloud"]; ok {
			return rhcos.FindArtifactURL(a)
		}
		return "", fmt.Errorf("%s: No ibmcloud build found", streamArchPrefix)
	case ovirt.Name, openstack.Name, powervc.Name:
		op := platform.OpenStack
		if op != nil {
			if oi := op.ClusterOSImage; oi != "" {
				return oi, nil
			}
		}
		if a, ok := streamArch.Artifacts["openstack"]; ok {
			return rhcos.FindArtifactURL(a)
		}
		return "", fmt.Errorf("%s: No openstack build found", streamArchPrefix)
	case azure.Name:
		ext := streamArch.RHELCoreOSExtensions
		if platform.Azure.CloudName == azure.StackCloud {
			return platform.Azure.ClusterOSImage, nil
		}
		if ext == nil {
			return "", fmt.Errorf("%s: No extensions found in stream", streamArchPrefix)
		}
		if ext.AzureDisk == nil {
			return "", fmt.Errorf("%s: No azure build found", streamArchPrefix)
		}
		azi := ext.AzureDisk.URL
		if mkt := ext.Marketplace; mkt == nil || mkt.Azure == nil || mkt.Azure.NoPurchasePlan == nil || mkt.Azure.NoPurchasePlan.Gen2 == nil {
			logrus.Warnf("%s: No default Azure marketplace image was found in stream", streamArchPrefix)
		} else {
			gen, err := getHyperVGeneration(ic.Azure, machinePool.Name)
			if err != nil {
				return "", fmt.Errorf("failed to get hyperVGeneration: %w", err)
			}
			azi = ext.Marketplace.Azure.NoPurchasePlan.Gen2.URN()
			if gen == "V1" {
				if mkt.Azure.NoPurchasePlan.Gen1 == nil {
					return "", fmt.Errorf("a HyperVGeneration 1 instance was selected but no Gen1 marketplace image is available")
				}
				azi = ext.Marketplace.Azure.NoPurchasePlan.Gen1.URN()
			}
		}
		return azi, nil
	case baremetal.Name:
		// Check for image URL override
		if oi := platform.BareMetal.ClusterOSImage; oi != "" {
			return oi, nil
		}
		// Use image from release payload
		return "", nil
	case vsphere.Name:
		// Check for image URL override
		if platform.VSphere.ClusterOSImage != "" {
			return platform.VSphere.ClusterOSImage, nil
		}

		if a, ok := streamArch.Artifacts["vmware"]; ok {
			// for an unknown reason vSphere OVAs are not
			// integrity checked. Instead of going through
			// FindArtifactURL just create the URL here.
			artifact := a.Formats["ova"].Disk
			u, err := url.Parse(artifact.Location)
			if err != nil {
				return "", err
			}

			// Add the sha256 query to the url
			// This will later be used in pkg/rhcos/cache/cache.go
			q := u.Query()
			q.Set("sha256", artifact.Sha256)

			u.RawQuery = q.Encode()

			return u.String(), nil
		}
		return "", fmt.Errorf("%s: No vmware build found", streamArchPrefix)
	case powervs.Name:
		// Check for image URL override
		if platform.PowerVS.ClusterOSImage != "" {
			return platform.PowerVS.ClusterOSImage, nil
		}

		if streamArch.Images.PowerVS != nil {
			var (
				vpcRegion string
				err       error
			)
			if platform.PowerVS.VPCRegion != "" {
				vpcRegion = platform.PowerVS.VPCRegion
			} else {
				vpcRegion = powervs.Regions[platform.PowerVS.Region].VPCRegion
			}
			vpcRegion, err = powervs.COSRegionForVPCRegion(vpcRegion)
			if err != nil {
				return "", fmt.Errorf("%s: No Power COS region found", streamArchPrefix)
			}
			img := streamArch.Images.PowerVS.Regions[vpcRegion]
			logrus.Debug("Power VS using image ", img.Object)
			return fmt.Sprintf("%s/%s", img.Bucket, img.Object), nil
		}

		return "", fmt.Errorf("%s: No Power VS build found", streamArchPrefix)
	case external.Name:
		return "", nil
	case none.Name:
		return "", nil
	case nutanix.Name:
		if platform.Nutanix != nil && platform.Nutanix.ClusterOSImage != "" {
			return platform.Nutanix.ClusterOSImage, nil
		}
		if a, ok := streamArch.Artifacts["nutanix"]; ok {
			return rhcos.FindArtifactURL(a)
		}
		return "", fmt.Errorf("%s: No nutanix build found", streamArchPrefix)
	default:
		return "", fmt.Errorf("invalid platform %v", platform.Name())
	}
}

// MakeAsset returns an Image asset with the given os image.
func MakeAsset(osImage string) *Image {
	return &Image{
		ControlPlane: osImage,
		Compute:      osImage,
	}
}

func getHyperVGeneration(metadata *icazure.Metadata, role string) (string, error) {
	if role == types.MachinePoolControlPlaneRoleName {
		return metadata.ControlPlaneHyperVGeneration()
	}
	return metadata.ComputeHyperVGeneration()
}
