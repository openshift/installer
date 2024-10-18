package nutanix

import (
	"context"
	"fmt"
	"strconv"
	"strings"
	"time"

	nutanixclientv3 "github.com/nutanix-cloud-native/prism-go-client/v3"
	"github.com/sirupsen/logrus"
	"k8s.io/apimachinery/pkg/util/validation/field"

	"github.com/openshift/installer/pkg/rhcos"
	"github.com/openshift/installer/pkg/types"
	nutanixtypes "github.com/openshift/installer/pkg/types/nutanix"
)

// Validate executes platform-specific validation.
func Validate(ic *types.InstallConfig) error {
	if ic.Platform.Nutanix == nil {
		return field.Required(field.NewPath("platform", "nutanix"), "nutanix validation requires a nutanix platform configuration")
	}

	return nil
}

// ValidateForProvisioning performs platform validation specifically for installer-
// provisioned infrastructure. In this case, self-hosted networking is a requirement
// when the installer creates infrastructure for nutanix clusters.
func ValidateForProvisioning(ic *types.InstallConfig) error {
	errList := field.ErrorList{}
	parentPath := field.NewPath("platform", "nutanix")

	if ic.Platform.Nutanix == nil {
		return field.Required(parentPath, "nutanix validation requires a nutanix platform configuration")
	}

	ctx, cancel := context.WithTimeout(context.TODO(), 60*time.Second)
	defer cancel()
	p := ic.Platform.Nutanix
	nc, err := nutanixtypes.CreateNutanixClient(ctx,
		p.PrismCentral.Endpoint.Address,
		strconv.Itoa(int(p.PrismCentral.Endpoint.Port)),
		p.PrismCentral.Username,
		p.PrismCentral.Password)
	if err != nil {
		return field.Invalid(parentPath.Child("prismCentral"), p.PrismCentral.Endpoint.Address, fmt.Sprintf("failed to connect to the prism-central with the configured credentials: %v", err))
	}

	// validate whether a prism element with the UUID actually exists
	for _, pe := range p.PrismElements {
		_, err = nc.V3.GetCluster(ctx, pe.UUID)
		if err != nil {
			errList = append(errList, field.Invalid(parentPath.Child("prismElements"), pe.UUID, fmt.Sprintf("the prism element %s's UUID does not correspond to a valid prism element in Prism: %v", pe.Name, err)))
		}
	}

	// validate whether a subnet with the UUID actually exists
	for _, subnetUUID := range p.SubnetUUIDs {
		_, err = nc.V3.GetSubnet(ctx, subnetUUID)
		if err != nil {
			errList = append(errList, field.Invalid(parentPath.Child("subnetUUIDs"), subnetUUID, fmt.Sprintf("the subnet UUID does not correspond to a valid subnet in Prism: %v", err)))
		}
	}

	// validate PreloadedOSImageName if configured
	if p.PreloadedOSImageName != "" {
		err = validatePreloadedImage(ctx, nc, p)
		if err != nil {
			errList = append(errList, field.Invalid(parentPath.Child("preloadedOSImageName"), p.PreloadedOSImageName, fmt.Sprintf("fail to validate the preloaded rhcos image: %v", err)))
		}
	}

	// validate each FailureDomain configuration
	for _, fd := range p.FailureDomains {
		// validate whether the prism element with the UUID exists
		_, err = nc.V3.GetCluster(ctx, fd.PrismElement.UUID)
		if err != nil {
			errList = append(errList, field.Invalid(parentPath.Child("failureDomains", "prismElements"), fd.PrismElement.UUID,
				fmt.Sprintf("the failure domain %s configured prism element UUID does not correspond to a valid prism element in Prism: %v", fd.Name, err)))
		}

		// validate whether a subnet with the UUID actually exists
		for _, subnetUUID := range fd.SubnetUUIDs {
			_, err = nc.V3.GetSubnet(ctx, subnetUUID)
			if err != nil {
				errList = append(errList, field.Invalid(parentPath.Child("failureDomains", "subnetUUIDs"), subnetUUID,
					fmt.Sprintf("the failure domain %s configured subnet UUID does not correspond to a valid subnet in Prism: %v", fd.Name, err)))
			}
		}

		// validate each configured dataSource image exists
		for _, dsImgRef := range fd.DataSourceImages {
			switch {
			case dsImgRef.UUID != "":
				if _, err = nc.V3.GetImage(ctx, dsImgRef.UUID); err != nil {
					errMsg := fmt.Sprintf("failureDomain %q: failed to find the dataSource image with uuid %s: %v", fd.Name, dsImgRef.UUID, err)
					errList = append(errList, field.Invalid(parentPath.Child("failureDomains", "dataSourceImage", "uuid"), dsImgRef.UUID, errMsg))
				}
			case dsImgRef.Name != "":
				if dsImgUUID, err := nutanixtypes.FindImageUUIDByName(ctx, nc, dsImgRef.Name); err != nil {
					errMsg := fmt.Sprintf("failureDomain %q: failed to find the dataSource image with name %q: %v", fd.Name, dsImgRef.UUID, err)
					errList = append(errList, field.Invalid(parentPath.Child("failureDomains", "dataSourceImage", "name"), dsImgRef.Name, errMsg))
				} else {
					dsImgRef.UUID = *dsImgUUID
				}
			default:
				errList = append(errList, field.Required(parentPath.Child("failureDomains", "dataSourceImage"), "both the dataSourceImage's uuid and name are empty, you need to configure one."))
			}
		}
	}

	return errList.ToAggregate()
}

func validatePreloadedImage(ctx context.Context, nc *nutanixclientv3.Client, p *nutanixtypes.Platform) error {
	// retrieve the rhcos release version
	rhcosStream, err := rhcos.FetchCoreOSBuild(ctx)
	if err != nil {
		return err
	}

	arch, ok := rhcosStream.Architectures["x86_64"]
	if !ok {
		return fmt.Errorf("unable to find the x86_64 rhcos architecture")
	}
	artifacts, ok := arch.Artifacts["nutanix"]
	if !ok {
		return fmt.Errorf("unable to find the x86_64 nutanix rhcos artifacts")
	}
	rhcosReleaseVersion := artifacts.Release

	// retrieve the rhcos version number from rhcosReleaseVersion
	rhcosVerNum, err := strconv.Atoi(strings.Split(rhcosReleaseVersion, ".")[0])
	if err != nil {
		return fmt.Errorf("failed to get the rhcos image version number from the version string %s: %w", rhcosReleaseVersion, err)
	}

	// retrieve the rhcos version number from the preloaded image object
	imgUUID, err := nutanixtypes.FindImageUUIDByName(ctx, nc, p.PreloadedOSImageName)
	if err != nil {
		return err
	}
	imgResp, err := nc.V3.GetImage(ctx, *imgUUID)
	if err != nil {
		return fmt.Errorf("failed to retrieve the rhcos image with uuid %s: %w", *imgUUID, err)
	}
	imgSource := *imgResp.Status.Resources.SourceURI

	si := strings.LastIndex(imgSource, "/rhcos-")
	if si < 0 {
		return fmt.Errorf("failed to get the rhcos image version from the preloaded image %s object's source_uri %s", p.PreloadedOSImageName, imgSource)
	}
	verStr := strings.Split(imgSource[si+7:], ".")[0]
	imgVerNum, err := strconv.Atoi(verStr)
	if err != nil {
		return fmt.Errorf("failed to get the rhcos image version number from the version string %s: %w", verStr, err)
	}

	// verify that the image version numbers are compactible
	versionDiff := rhcosVerNum - imgVerNum
	switch {
	case versionDiff < 0:
		return fmt.Errorf("the preloaded image's rhcos version: %v is too many revisions ahead the installer bundled rhcos version: %v", imgVerNum, rhcosVerNum)
	case versionDiff >= 2:
		return fmt.Errorf("the preloaded image's rhcos version: %v is too many revisions behind the installer bundled rhcos version: %v", imgVerNum, rhcosVerNum)
	case versionDiff == 1:
		logrus.Warnf("the preloaded image's rhcos version: %v is behind the installer bundled rhcos version: %v, installation may fail", imgVerNum, rhcosVerNum)
	}

	return nil
}
