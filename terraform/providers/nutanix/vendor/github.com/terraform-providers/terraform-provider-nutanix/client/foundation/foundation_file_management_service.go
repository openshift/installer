package foundation

import (
	"context"
	"fmt"
	"net/http"
	"os"

	"github.com/terraform-providers/terraform-provider-nutanix/client"
)

// Interface for file management apis of foundation
type FileManagementService interface {
	ListNOSPackages(context.Context) (*ListNOSPackagesResponse, error)
	ListHypervisorISOs(context.Context) (*ListHypervisorISOsResponse, error)
	UploadImage(context.Context, string, string, string) (*UploadImageResponse, error)
	DeleteImage(context.Context, string, string) error
}

// FileManagementOperations implements FileManagementService interface
type FileManagementOperations struct {
	client *client.Client
}

//ListNOSPackages lists the available AOS packages file names in Foundation
func (fmo FileManagementOperations) ListNOSPackages(ctx context.Context) (*ListNOSPackagesResponse, error) {
	path := "/enumerate_nos_packages"
	req, err := fmo.client.NewUnAuthRequest(ctx, http.MethodGet, path, nil)
	if err != nil {
		return nil, err
	}
	listNOSPackagesResponse := new(ListNOSPackagesResponse)
	return listNOSPackagesResponse, fmo.client.Do(ctx, req, listNOSPackagesResponse)
}

//ListHypervisorISOs lists the hypervisor ISOs available in Foundation
func (fmo FileManagementOperations) ListHypervisorISOs(ctx context.Context) (*ListHypervisorISOsResponse, error) {
	path := "/enumerate_hypervisor_isos"
	req, err := fmo.client.NewUnAuthRequest(ctx, http.MethodGet, path, nil)
	if err != nil {
		return nil, err
	}
	listHypervisorISOsResponse := new(ListHypervisorISOsResponse)
	return listHypervisorISOsResponse, fmo.client.Do(ctx, req, listHypervisorISOsResponse)
}

//UploadImage uploads the image to foundation vm as per installer type
func (fmo FileManagementOperations) UploadImage(ctx context.Context, installerType, fileName, source string) (*UploadImageResponse, error) {
	path := fmt.Sprintf("/upload?installer_type=%s&filename=%s", installerType, fileName)

	// open file. The source should be complete file path eg./Users/...
	file, err := os.Open(source)
	if err != nil {
		return nil, fmt.Errorf("error while opening file: %s", err)
	}
	defer file.Close()

	req, err := fmo.client.NewUnAuthUploadRequest(ctx, http.MethodPost, path, file)
	if err != nil {
		return nil, err
	}

	uploadImageResponse := new(UploadImageResponse)
	return uploadImageResponse, fmo.client.Do(ctx, req, uploadImageResponse)
}

func (fmo FileManagementOperations) DeleteImage(ctx context.Context, installerType, fileName string) error {
	path := "/delete/"

	// create body
	body := make(map[string]string)
	body["installer_type"] = installerType
	body["filename"] = fileName

	req, err := fmo.client.NewUnAuthFormEncodedRequest(ctx, http.MethodPost, path, body)
	if err != nil {
		return err
	}
	return fmo.client.Do(ctx, req, nil)
}
