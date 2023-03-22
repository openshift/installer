package azure

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/Azure/azure-sdk-for-go/profiles/2018-03-01/storage/mgmt/storage"
	"github.com/Azure/azure-sdk-for-go/profiles/2020-09-01/compute/mgmt/compute"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azblob"
	"github.com/Azure/go-autorest/autorest/to"
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
	virtualMachinesClient compute.VirtualMachinesClient
	accountsClient        storage.AccountsClient
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

	accountsClient := storage.NewAccountsClientWithBaseURI(session.Environment.ResourceManagerEndpoint, session.Credentials.SubscriptionID)
	accountsClient.Authorizer = session.Authorizer

	virtualMachinesClient := compute.NewVirtualMachinesClientWithBaseURI(session.Environment.ResourceManagerEndpoint, session.Credentials.SubscriptionID)
	virtualMachinesClient.Authorizer = session.Authorizer

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

	storageAccounts, err := getStorageAccounts(ctx, g)
	if err != nil {
		return err
	}

	sharedKeyCredentials, err := getSharedKeyCredentials(ctx, storageAccounts, g)
	if err != nil {
		return err
	}

	virtualMachines, err := getVirtualMachines(ctx, g)
	if err != nil {
		return err
	}

	// We can only get the serial log from VM's with boot diagnostics enabled
	bootDiagnostics := getBootDiagnostics(ctx, virtualMachines, g)

	err = downloadFiles(ctx, bootDiagnostics, storageAccounts, sharedKeyCredentials, g)
	if err != nil {
		return err
	}

	return nil
}

func getStorageAccounts(ctx context.Context, g *Gather) ([]storage.Account, error) {
	accountListResult, err := g.accountsClient.ListByResourceGroup(ctx, g.resourceGroupName)
	if err != nil {
		return nil, errors.Wrap(err, "could not list storage accounts")
	}
	if accountListResult.Value == nil {
		return nil, fmt.Errorf("could not find any storage accounts")
	}
	return *accountListResult.Value, nil
}

func getSharedKeyCredentials(ctx context.Context, storageAccounts []storage.Account, g *Gather) ([]*azblob.SharedKeyCredential, error) {
	var sharedKeyCredentials []*azblob.SharedKeyCredential
	for _, account := range storageAccounts {
		accountName := to.String(account.Name)
		keyResults, err := g.accountsClient.ListKeys(ctx, g.resourceGroupName, accountName)
		if err != nil {
			g.logger.Debugf("failed to list keys: %s", err.Error())
			continue
		}
		if keyResults.Keys != nil {
			for _, key := range *keyResults.Keys {
				if key.Value != nil {
					sharedKeyCredential, err := azblob.NewSharedKeyCredential(accountName, to.String(key.Value))
					if err != nil {
						g.logger.Debugf("failed to get shared key: %s", err.Error())
						continue
					}
					sharedKeyCredentials = append(sharedKeyCredentials, sharedKeyCredential)
				}
			}
		}
	}

	return sharedKeyCredentials, nil
}

func getVirtualMachines(ctx context.Context, g *Gather) ([]compute.VirtualMachine, error) {
	vmsPage, err := g.virtualMachinesClient.List(ctx, g.resourceGroupName)
	if err != nil {
		return nil, err
	}

	var virtualMachines []compute.VirtualMachine
	for ; vmsPage.NotDone(); err = vmsPage.NextWithContext(ctx) {
		if err != nil {
			g.logger.Debugf("failed to get vm: %s", err.Error())
			continue
		}
		virtualMachines = append(virtualMachines, vmsPage.Values()...)
	}

	return virtualMachines, nil
}

func getBootDiagnostics(ctx context.Context, virtualMachines []compute.VirtualMachine, g *Gather) []string {
	var bootDiagnostics []string
	for _, vm := range virtualMachines {
		if vm.DiagnosticsProfile == nil ||
			vm.DiagnosticsProfile.BootDiagnostics == nil ||
			!to.Bool(vm.DiagnosticsProfile.BootDiagnostics.Enabled) {
			g.logger.Debugf("boot diagnostics not found or not enabled for %s", to.String(vm.Name))
			continue
		}
		var screenshotURI *string
		var serialLogURI *string
		instanceView, err := g.virtualMachinesClient.InstanceView(ctx, g.resourceGroupName, to.String(vm.Name))
		if err != nil {
			g.logger.Debugf("failed to get instance view: %v", err)
			continue
		}
		if instanceView.BootDiagnostics != nil {
			screenshotURI = instanceView.BootDiagnostics.ConsoleScreenshotBlobURI
			serialLogURI = instanceView.BootDiagnostics.SerialConsoleLogBlobURI

		}
		if screenshotURI == nil && serialLogURI == nil {
			data, err := g.virtualMachinesClient.RetrieveBootDiagnosticsData(ctx, g.resourceGroupName, *vm.Name, to.Int32Ptr(10))
			if err != nil {
				g.logger.Debugf("failed to get boot diagnostics data: %v", err)
				continue
			}
			screenshotURI = data.ConsoleScreenshotBlobURI
			serialLogURI = data.SerialConsoleLogBlobURI
		}
		if screenshotURI != nil {
			bootDiagnostics = append(bootDiagnostics, *screenshotURI)
		}
		if serialLogURI != nil {
			bootDiagnostics = append(bootDiagnostics, *serialLogURI)
		}
	}

	return bootDiagnostics
}

func downloadFiles(ctx context.Context, fileURIs []string, storageAccounts []storage.Account, sharedKeyCredentials []*azblob.SharedKeyCredential, g *Gather) error {
	var errs []error
	var files []string

	serialLogBundleDir := filepath.Join(g.directory, strings.TrimSuffix(filepath.Base(g.serialLogBundle), ".tar.gz"))
	filePathDir := filepath.Join(g.directory, serialLogBundleDir)
	err := os.MkdirAll(filePathDir, 0o755)
	if err != nil && !errors.Is(err, os.ErrExist) {
		return err
	}

	for _, fileURI := range fileURIs {
		var filePath string
		var ferr error
		if isBlobFromManagedAccount(fileURI, storageAccounts) {
			filePath, ferr = downloadFile(ctx, fileURI, filePathDir, g.logger)
		} else {
			filePath, ferr = downloadBlob(ctx, fileURI, filePathDir, sharedKeyCredentials, g.logger)
		}
		if ferr != nil {
			errs = append(errs, ferr)
			continue
		}
		files = append(files, filePath)
	}

	if len(files) > 0 {
		cerr := gather.CreateArchive(files, g.serialLogBundle)
		if cerr != nil {
			g.logger.Debugf("failed to create archive: %s", cerr.Error())
			errs = append(errs, cerr)
		}
	}

	err = gather.DeleteArchiveDirectory(serialLogBundleDir)
	if err != nil {
		g.logger.Debugf("failed to remove archive directory: %v", err)
	}

	return utilerrors.NewAggregate(errs)
}

func isBlobFromManagedAccount(blobURI string, storageAccounts []storage.Account) bool {
	for _, account := range storageAccounts {
		if account.PrimaryEndpoints != nil &&
			account.PrimaryEndpoints.Blob != nil &&
			strings.HasPrefix(blobURI, *account.PrimaryEndpoints.Blob) {
			return false
		}
	}
	return true
}

func downloadBlob(ctx context.Context, fileURI string, filePathDir string, sharedKeyCredentials []*azblob.SharedKeyCredential, logger logrus.FieldLogger) (string, error) {
	logger.Debugf("attemping to download %s from storage account", fileURI)

	uri, err := url.ParseRequestURI(fileURI)
	if err != nil {
		return "", err
	}
	filePath := filepath.Join(filePathDir, filepath.Base(uri.Path))

	for _, credential := range sharedKeyCredentials {
		// The credential is not for the blob's account
		if !strings.HasPrefix(uri.Host, credential.AccountName()) {
			continue
		}

		blobClient, err := azblob.NewBlobClientWithSharedKey(fileURI, credential, nil)
		if err != nil {
			logger.Debugf("failed to create blob client: %s", err.Error())
			continue
		}

		dr, err := blobClient.Download(ctx, nil)
		if err != nil {
			logger.Debugf("failed to download blob: %s", err.Error())
			continue
		}
		data := &bytes.Buffer{}
		reader := dr.Body(&azblob.RetryReaderOptions{MaxRetryRequests: 3})
		_, err = data.ReadFrom(reader)
		if err != nil {
			logger.Debugf("failed to read: %s", err.Error())
			return "", err
		}
		err = reader.Close()
		if err != nil {
			return "", err
		}

		file, err := os.Create(filePath)
		if err != nil {
			return "", errors.Wrapf(err, "failed to create file %s", filePath)
		}
		defer file.Close()

		_, err = file.Write(data.Bytes())
		if err != nil {
			logger.Debugf("failed to write to file: %s", err.Error())
			return "", err
		}

		return filePath, nil
	}

	return "", errors.Errorf("unable to download file: %s", filePath)
}

func downloadFile(ctx context.Context, fileURI string, filePathDir string, logger logrus.FieldLogger) (string, error) {
	logger.Debugln("attemping to download file from managed account")

	uri, err := url.ParseRequestURI(fileURI)
	if err != nil {
		return "", errors.Wrapf(err, "unable to parse file URI %s", fileURI)
	}
	filePath := filepath.Join(filePathDir, filepath.Base(uri.Path))

	req, err := http.NewRequestWithContext(ctx, "GET", fileURI, nil)
	if err != nil {
		return "", err
	}
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("unable to download file: %s", filePath)
	}

	file, err := os.Create(filePath)
	if err != nil {
		return "", errors.Wrapf(err, "failed to create file %s", filePath)
	}
	defer file.Close()

	_, err = io.Copy(file, resp.Body)
	if err != nil {
		return "", errors.Wrapf(err, "failed to write to file %s", filePath)
	}

	return filePath, nil
}
