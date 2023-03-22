package azure

import (
	"context"
	"fmt"
	"net/url"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore/arm"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/policy"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/compute/armcompute/v4"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/storage/armstorage"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azblob"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	utilerrors "k8s.io/apimachinery/pkg/util/errors"

	azuresession "github.com/openshift/installer/pkg/asset/installconfig/azure"
	"github.com/openshift/installer/pkg/gather"
	"github.com/openshift/installer/pkg/gather/providers"
	"github.com/openshift/installer/pkg/types"
	"github.com/openshift/installer/pkg/types/azure"
)

// Gather holds options for resources we want to gather.
type Gather struct {
	resourceGroupName     string
	logger                logrus.FieldLogger
	serialLogBundle       string
	directory             string
	virtualMachinesClient *armcompute.VirtualMachinesClient
	accountsClient        *armstorage.AccountsClient
}

// New returns a Azure Gather from ClusterMetadata.
func New(logger logrus.FieldLogger, serialLogBundle string, bootstrap string, masters []string, metadata *types.ClusterMetadata) (providers.Gather, error) {
	cloudName := metadata.Azure.CloudName
	if cloudName == "" {
		cloudName = azure.PublicCloud
	}

	resourceGroupName := metadata.Azure.ResourceGroupName
	if resourceGroupName == "" {
		resourceGroupName = metadata.InfraID + "-rg"
	}

	session, err := azuresession.GetSession(cloudName, metadata.Azure.ARMEndpoint)
	if err != nil {
		return nil, err
	}

	accountClientOptions := arm.ClientOptions{
		ClientOptions: policy.ClientOptions{
			// NOTE: the api version must support AzureStack
			APIVersion: "2019-04-01",
			Cloud:      session.CloudConfig,
		},
	}
	accountsClient, err := armstorage.NewAccountsClient(session.Credentials.SubscriptionID, session.TokenCreds, &accountClientOptions)
	if err != nil {
		return nil, err
	}

	vmClientOptions := arm.ClientOptions{
		ClientOptions: policy.ClientOptions{
			// NOTE: the api version must both support AzureStack and BootDignosticsData
			APIVersion: "2020-06-01",
			Cloud:      session.CloudConfig,
		},
	}
	virtualMachinesClient, err := armcompute.NewVirtualMachinesClient(session.Credentials.SubscriptionID, session.TokenCreds, &vmClientOptions)
	if err != nil {
		return nil, err
	}

	gather := &Gather{
		resourceGroupName:     resourceGroupName,
		logger:                logger,
		serialLogBundle:       serialLogBundle,
		directory:             filepath.Dir(serialLogBundle),
		accountsClient:        accountsClient,
		virtualMachinesClient: virtualMachinesClient,
	}

	return gather, nil
}

// Run is the entrypoint to start the gather process.
func (g *Gather) Run() error {
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Minute)
	defer cancel()

	sharedKeyCredentials, err := getSharedKeyCredentials(ctx, g)
	if err != nil {
		return err
	}

	virtualMachines, err := getVirtualMachines(ctx, g)
	if err != nil {
		return err
	}

	var files []string

	// We can only get the serial log from VM's with boot diagnostics enabled
	bootDiagnostics := getBootDiagnostics(ctx, virtualMachines, g)
	for _, bootDiagnostic := range bootDiagnostics {
		if bootDiagnostic.ConsoleScreenshotBlobURI != nil {
			files = append(files, *bootDiagnostic.ConsoleScreenshotBlobURI)
		}

		if bootDiagnostic.SerialConsoleLogBlobURI != nil {
			files = append(files, *bootDiagnostic.SerialConsoleLogBlobURI)
		}
	}

	err = downloadFiles(ctx, files, sharedKeyCredentials, g)
	if err != nil {
		return err
	}

	return nil
}

func getSharedKeyCredentials(ctx context.Context, g *Gather) ([]*azblob.SharedKeyCredential, error) {
	pager := g.accountsClient.NewListByResourceGroupPager(g.resourceGroupName, nil)

	var sharedKeyCredentials []*azblob.SharedKeyCredential
	for pager.More() {
		accountListResult, err := pager.NextPage(ctx)
		if err != nil {
			return nil, errors.Wrap(err, "could not find any storage account keys")
		}
		for _, account := range accountListResult.Value {
			keyResults, err := g.accountsClient.ListKeys(ctx, g.resourceGroupName, *account.Name, nil)
			if err != nil {
				g.logger.Debugf("failed to list keys: %s", err.Error())
				continue
			}
			if keyResults.Keys != nil {
				for _, key := range keyResults.Keys {
					if key.Value != nil {
						sharedKeyCredential, err := azblob.NewSharedKeyCredential(*account.Name, *key.Value)
						if err != nil {
							g.logger.Debugf("failed to get shared key: %s", err.Error())
							continue
						}
						sharedKeyCredentials = append(sharedKeyCredentials, sharedKeyCredential)
					}
				}
			}
		}
	}

	return sharedKeyCredentials, nil
}

func getVirtualMachines(ctx context.Context, g *Gather) ([]*armcompute.VirtualMachine, error) {
	vmsPager := g.virtualMachinesClient.NewListPager(g.resourceGroupName, nil)

	var virtualMachines []*armcompute.VirtualMachine
	for vmsPager.More() {
		vmsPage, err := vmsPager.NextPage(ctx)
		if err != nil {
			g.logger.Debugf("failed to get vm: %s", err.Error())
			return nil, err
		}
		virtualMachines = append(virtualMachines, vmsPage.Value...)
	}

	return virtualMachines, nil
}

func getBootDiagnostics(ctx context.Context, virtualMachines []*armcompute.VirtualMachine, g *Gather) []*armcompute.BootDiagnosticsInstanceView {
	var bootDiagnostics []*armcompute.BootDiagnosticsInstanceView
	for _, vm := range virtualMachines {
		if vm.Properties.DiagnosticsProfile != nil &&
			vm.Properties.DiagnosticsProfile.BootDiagnostics != nil &&
			vm.Properties.DiagnosticsProfile.BootDiagnostics.Enabled != nil &&
			*vm.Properties.DiagnosticsProfile.BootDiagnostics.Enabled {
			instanceView, err := g.virtualMachinesClient.InstanceView(ctx, g.resourceGroupName, *vm.Name, nil)
			if err != nil {
				g.logger.Debugf("failed to get instance view: %v", err)
				continue
			}
			if instanceView.BootDiagnostics != nil {
				bootDiagnostics = append(bootDiagnostics, instanceView.BootDiagnostics)
			}
		}
	}

	return bootDiagnostics
}

func downloadFiles(ctx context.Context, fileURIs []string, sharedKeyCredentials []*azblob.SharedKeyCredential, g *Gather) error {
	var errs []error

	serialLogBundleDir := filepath.Join(g.directory, strings.TrimSuffix(filepath.Base(g.serialLogBundle), ".tar.gz"))

	err := os.MkdirAll(serialLogBundleDir, 0o755)
	if err != nil && !errors.Is(err, os.ErrExist) {
		return err
	}

	files := make([]string, 0, len(fileURIs))
	for _, fileURI := range fileURIs {
		filePath, err := downloadFile(ctx, fileURI, serialLogBundleDir, sharedKeyCredentials, g)
		if err != nil {
			errs = append(errs, err)
			continue
		}
		files = append(files, filePath)
	}

	if len(files) > 0 {
		err := gather.CreateArchive(files, g.serialLogBundle)
		if err != nil {
			g.logger.Debugf("failed to create archive: %s", err.Error())
			errs = append(errs, err)
		}
	}

	err = gather.DeleteArchiveDirectory(serialLogBundleDir)
	if err != nil {
		g.logger.Debugf("failed to remove archive directory: %v", err)
	}

	return utilerrors.NewAggregate(errs)
}

func downloadFile(ctx context.Context, fileURI string, filePathDir string, sharedKeyCredentials []*azblob.SharedKeyCredential, g *Gather) (string, error) {
	g.logger.Debugf("attemping to download %s", fileURI)

	uri, err := url.ParseRequestURI(fileURI)
	if err != nil {
		return "", err
	}
	uriParts := strings.Split(uri.Path, "/")
	containerName := uriParts[len(uriParts)-2]
	blobName := uriParts[len(uriParts)-1]
	filePath := filepath.Join(filePathDir, blobName)

	accountURL := fmt.Sprintf("%s://%s/", uri.Scheme, uri.Host)
	for _, credential := range sharedKeyCredentials {
		if !strings.HasPrefix(uri.Host, credential.AccountName()) {
			continue
		}
		blobClient, err := azblob.NewClientWithSharedKeyCredential(accountURL, credential, nil)
		if err != nil {
			g.logger.Debugf("failed to create blob client: %s", err.Error())
			continue
		}

		file, err := os.Create(filePath)
		if err != nil {
			g.logger.Debugf("failed to create file: %s", err.Error())
			return "", err
		}
		defer file.Close()

		_, err = blobClient.DownloadFile(ctx, containerName, blobName, file, nil)
		if err != nil {
			return "", err
		}

		return filePath, nil
	}

	return "", errors.Errorf("unable to download file: %s", filePath)
}
