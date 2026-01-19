package util

import (
	"fmt"
	"strings"

	machinev1 "github.com/openshift/api/machine/v1beta1"
	computeservice "github.com/openshift/machine-api-provider-gcp/pkg/cloud/gcp/actuators/services/compute"
)

const (
	UEFICompatible = "UEFI_COMPATIBLE"
)

// IsUEFICompatible detects if the machine's boot disk was created with a UEFI image.
// Shielded VMs require UEFI-compatible images. In some instances customers may still
// be using non UEFI-compatible images. e.g OpenShift images listed on the GCP marketplace
// are not updated with every release, and the 4.8 image is used until 4.12 and was not
// created with UEFI support.
func IsUEFICompatible(gceService computeservice.GCPComputeService, providerConfig *machinev1.GCPMachineProviderSpec) (bool, error) {
	for _, disk := range providerConfig.Disks {
		if !disk.Boot {
			continue
		}

		// Parse the image reference from the disk.
		imageRef, err := parseImageReference(disk.Image, providerConfig.ProjectID)
		if err != nil {
			return false, fmt.Errorf("failed to parse disk image reference: %w", err)
		}

		if imageRef.IsFamily {
			return checkImageFamilyImageUEFICompatible(gceService, imageRef.Project, providerConfig.Zone, imageRef.Image)
		}
		return checkImageUEFICompatible(gceService, imageRef.Project, imageRef.Image)
	}
	return false, fmt.Errorf("no boot disk found")
}

// imageReference holds parsed details from a disk.Image string.
type imageReference struct {
	Project  string
	Image    string
	IsFamily bool // true if the image is specified as a family reference.
}

// parseImageReference extracts project and image information from the given image string.
// It supports various formats:
//   - "projects/{project}/global/images/{image}"
//   - "projects/{project}/global/images/family/{imageFamily}"
//   - "https://www.googleapis.com/compute/v1/projects/{project}/global/images/{image}"
//   - A simple image name without slashes, in which case providerProject is used.
func parseImageReference(imageStr, providerProject string) (*imageReference, error) {
	// If the image string does not contain a slash, assume it's a simple image name.
	if !strings.Contains(imageStr, "/") {
		return &imageReference{
			Project:  providerProject,
			Image:    imageStr,
			IsFamily: false,
		}, nil
	}

	// Check if the image string contains "projects/".
	if !strings.Contains(imageStr, "projects/") {
		return nil, fmt.Errorf("image string %q does not contain expected 'projects/' segment", imageStr)
	}

	// Split based on "projects/".
	parts := strings.SplitN(imageStr, "projects/", 2)
	if len(parts) < 2 {
		return nil, fmt.Errorf("unexpected format for image string: %q", imageStr)
	}

	// Split the remainder by "/".
	subParts := strings.Split(parts[1], "/")
	// Expected formats:
	// For non-family images: {project}/global/images/{image} => at least 4 parts.
	// For family images: {project}/global/images/family/{imageFamily} => at least 5 parts.
	if len(subParts) < 4 {
		return nil, fmt.Errorf("unexpected image path format in %q", imageStr)
	}

	// Determine if the image is specified as a family.
	isFamily := false
	imageName := ""
	if len(subParts) >= 5 && subParts[2] == "images" && subParts[3] == "family" {
		isFamily = true
		imageName = subParts[4]
	} else if subParts[1] == "global" && subParts[2] == "images" {
		imageName = subParts[3]
	} else {
		return nil, fmt.Errorf("unrecognized image path format in %q", imageStr)
	}

	return &imageReference{
		Project:  subParts[0],
		Image:    imageName,
		IsFamily: isFamily,
	}, nil
}

// checkImageUEFICompatible retrieves the image and checks its GuestOSFeatures for UEFI support.
func checkImageUEFICompatible(gceService computeservice.GCPComputeService, project, image string) (bool, error) {
	img, err := gceService.ImageGet(project, image)
	if err != nil {
		return false, fmt.Errorf("unable to retrieve image %q in project %q: %w", image, project, err)
	}
	for _, feat := range img.GuestOsFeatures {
		if strings.EqualFold(feat.Type, UEFICompatible) {
			return true, nil
		}
	}
	return false, nil
}

// checkImageFamilyImageUEFICompatible retrieves the image family and checks for UEFI support.
func checkImageFamilyImageUEFICompatible(gceService computeservice.GCPComputeService, project, zone, imageFamily string) (bool, error) {
	family, err := gceService.ImageFamilyGet(project, zone, imageFamily)
	if err != nil {
		return false, fmt.Errorf("unable to retrieve image family %q in project %q: %w", imageFamily, project, err)
	}
	for _, feat := range family.Image.GuestOsFeatures {
		if strings.EqualFold(feat.Type, UEFICompatible) {
			return true, nil
		}
	}
	return false, nil
}
