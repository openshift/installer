package powervs

import (
	"context"
	"fmt"

	"github.com/IBM-Cloud/power-go-client/clients/instance"
	"github.com/IBM-Cloud/power-go-client/ibmpisession"
	"github.com/IBM/go-sdk-core/v5/core"
	"github.com/sirupsen/logrus"

	powervsconfig "github.com/openshift/installer/pkg/asset/installconfig/powervs"
)

// GetBootImageFromWorkspace retrieves a boot image from the PowerVS workspace.
// If an active image is found in the workspace, it returns that image name.
// If no images are found, it returns a default name in the format "rhcos-{clusterID}".
func GetBootImageFromWorkspace(ctx context.Context, serviceInstanceGUID string, zone string, clusterID string) (string, error) {
	// Create a new PowerVS client
	client, err := powervsconfig.NewClient()
	if err != nil {
		return "", fmt.Errorf("failed to create PowerVS client: %w", err)
	}

	// Create authenticator
	authenticator := &core.IamAuthenticator{
		ApiKey: client.GetAPIKey(),
	}

	// Create PI session
	piSession, err := ibmpisession.NewIBMPISession(&ibmpisession.IBMPIOptions{
		Authenticator: authenticator,
		UserAccount:   client.BXCli.User.Account,
		Zone:          zone,
		Debug:         false,
	})
	if err != nil {
		return "", fmt.Errorf("failed to create PowerVS session: %w", err)
	}

	// Create image client
	imageClient := instance.NewIBMPIImageClient(ctx, piSession, serviceInstanceGUID)
	if imageClient == nil {
		return "", fmt.Errorf("failed to create PowerVS image client")
	}

	// Get all images from the workspace
	images, err := imageClient.GetAll()
	if err != nil {
		return "", fmt.Errorf("failed to list images from PowerVS workspace: %w", err)
	}

	// If no images found in workspace, use default naming pattern
	if images == nil || len(images.Images) == 0 {
		defaultImage := fmt.Sprintf("rhcos-%s", clusterID)
		logrus.Infof("No images found in PowerVS workspace, using default image name: %s", defaultImage)
		return defaultImage, nil
	}

	// Find the first active image
	for _, image := range images.Images {
		if image.State != nil && *image.State == "active" && image.Name != nil {
			logrus.Infof("Selected PowerVS boot image from workspace: %s", *image.Name)
			return *image.Name, nil
		}
	}

	// If no active images found, use default naming pattern
	defaultImage := fmt.Sprintf("rhcos-%s", clusterID)
	logrus.Infof("No active images found in PowerVS workspace, using default image name: %s", defaultImage)
	return defaultImage, nil
}

// Made with Bob
