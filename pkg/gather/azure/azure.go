package azure

import (
	"bytes"
	"context"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/Azure/azure-sdk-for-go/profiles/2018-03-01/compute/mgmt/compute"
	"github.com/Azure/azure-sdk-for-go/profiles/2018-03-01/storage/mgmt/storage"
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
		screenBmp := to.String(bootDiagnostic.ConsoleScreenshotBlobURI)
		files = append(files, screenBmp)

		serialLog := to.String(bootDiagnostic.SerialConsoleLogBlobURI)
		files = append(files, serialLog)
	}

	err = downloadFiles(ctx, files, sharedKeyCredentials, g)
	if err != nil {
		return err
	}

	return nil
}

func getSharedKeyCredentials(ctx context.Context, g *Gather) ([]*azblob.SharedKeyCredential, error) {
	accountListResult, err := g.accountsClient.ListByResourceGroup(ctx, g.resourceGroupName)
	if err != nil {
		return nil, errors.Wrap(err, "could not find any storage accounts")
	}

	var sharedKeyCredentials []*azblob.SharedKeyCredential
	if accountListResult.Value != nil {
		for _, account := range *accountListResult.Value {
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
		for _, virtualMachine := range vmsPage.Values() {
			virtualMachines = append(virtualMachines, virtualMachine)
		}
	}

	return virtualMachines, nil
}

func getBootDiagnostics(ctx context.Context, virtualMachines []compute.VirtualMachine, g *Gather) []*compute.BootDiagnosticsInstanceView {
	var bootDiagnostics []*compute.BootDiagnosticsInstanceView
	for _, vm := range virtualMachines {
		if vm.VirtualMachineProperties.DiagnosticsProfile != nil &&
			vm.VirtualMachineProperties.DiagnosticsProfile.BootDiagnostics != nil &&
			to.Bool(vm.VirtualMachineProperties.DiagnosticsProfile.BootDiagnostics.Enabled) {
			instanceView, err := g.virtualMachinesClient.InstanceView(ctx, g.resourceGroupName, to.String(vm.Name))
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
	var files []string

	for _, fileURI := range fileURIs {
		filePath, err := downloadFile(ctx, fileURI, sharedKeyCredentials, g)
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

	serialLogBundleDir := filepath.Join(g.directory, strings.TrimSuffix(filepath.Base(g.serialLogBundle), ".tar.gz"))
	err := gather.DeleteArchiveDirectory(serialLogBundleDir)
	if err != nil {
		g.logger.Debugf("failed to remove archive directory: %v", err)
	}

	return utilerrors.NewAggregate(errs)
}

func downloadFile(ctx context.Context, fileURI string, sharedKeyCredentials []*azblob.SharedKeyCredential, g *Gather) (string, error) {
	directory := g.directory
	g.logger.Debugf("attemping to download %s", fileURI)

	serialLogBundleDir := strings.TrimSuffix(filepath.Base(g.serialLogBundle), ".tar.gz")
	filePathDir := filepath.Join(directory, serialLogBundleDir)
	filePath := filepath.Join(filePathDir, filepath.Base(fileURI))

	err := os.MkdirAll(filePathDir, 0755)
	if err != nil && !errors.Is(err, os.ErrExist) {
		return "", err
	}

	for _, credential := range sharedKeyCredentials {
		blobClient, err := azblob.NewBlobClientWithSharedKey(fileURI, credential, nil)
		if err != nil {
			g.logger.Debugf("failed to create blob client: %s", err.Error())
			continue
		}
		dr, err := blobClient.Download(ctx, nil)
		if err != nil {
			continue
		}

		data := &bytes.Buffer{}
		reader := dr.Body(&azblob.RetryReaderOptions{MaxRetryRequests: 3})
		_, err = data.ReadFrom(reader)
		if err != nil {
			g.logger.Debugf("failed to read: %s", err.Error())
			return "", err
		}
		err = reader.Close()
		if err != nil {
			return "", err
		}

		file, err := os.Create(filePath)
		if err != nil {
			g.logger.Debugf("failed to create file: %s", err.Error())
			return "", err
		}

		_, err = file.Write(data.Bytes())
		if err != nil {
			g.logger.Debugf("failed to write to file: %s", err.Error())
			file.Close()
			return "", err
		}

		file.Close()
		return filePath, nil
	}

	return "", errors.Errorf("unable to download file: %s", filePath)
}
