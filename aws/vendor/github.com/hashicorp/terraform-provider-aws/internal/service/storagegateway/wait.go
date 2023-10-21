package storagegateway

import (
	"context"
	"time"

	"github.com/aws/aws-sdk-go/service/storagegateway"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
)

const (
	gatewayConnectedMinTimeout                = 10 * time.Second
	gatewayConnectedContinuousTargetOccurence = 6
	gatewayJoinDomainJoinedTimeout            = 5 * time.Minute
	storediSCSIVolumeAvailableTimeout         = 5 * time.Minute
	nfsFileShareAvailableDelay                = 5 * time.Second
	nfsFileShareDeletedDelay                  = 5 * time.Second
	smbFileShareAvailableDelay                = 5 * time.Second
	smbFileShareDeletedDelay                  = 5 * time.Second
	fileSystemAssociationAvailableDelay       = 5 * time.Second
	fileSystemAssociationDeletedDelay         = 5 * time.Second
)

func waitGatewayConnected(ctx context.Context, conn *storagegateway.StorageGateway, gatewayARN string, timeout time.Duration) (*storagegateway.DescribeGatewayInformationOutput, error) {
	stateConf := &retry.StateChangeConf{
		Pending:                   []string{storagegateway.ErrorCodeGatewayNotConnected},
		Target:                    []string{gatewayStatusConnected},
		Refresh:                   statusGateway(ctx, conn, gatewayARN),
		Timeout:                   timeout,
		MinTimeout:                gatewayConnectedMinTimeout,
		ContinuousTargetOccurence: gatewayConnectedContinuousTargetOccurence, // Gateway activations can take a few seconds and can trigger a reboot of the Gateway
	}

	outputRaw, err := stateConf.WaitForStateContext(ctx)

	switch output := outputRaw.(type) {
	case *storagegateway.DescribeGatewayInformationOutput:
		return output, err
	default:
		return nil, err
	}
}

func waitGatewayJoinDomainJoined(ctx context.Context, conn *storagegateway.StorageGateway, volumeARN string) (*storagegateway.DescribeSMBSettingsOutput, error) { //nolint:unparam
	stateConf := &retry.StateChangeConf{
		Pending: []string{storagegateway.ActiveDirectoryStatusJoining},
		Target:  []string{storagegateway.ActiveDirectoryStatusJoined},
		Refresh: statusGatewayJoinDomain(ctx, conn, volumeARN),
		Timeout: gatewayJoinDomainJoinedTimeout,
	}

	outputRaw, err := stateConf.WaitForStateContext(ctx)

	if output, ok := outputRaw.(*storagegateway.DescribeSMBSettingsOutput); ok {
		return output, err
	}

	return nil, err
}

// waitStorediSCSIVolumeAvailable waits for a StoredIscsiVolume to return Available
func waitStorediSCSIVolumeAvailable(ctx context.Context, conn *storagegateway.StorageGateway, volumeARN string) (*storagegateway.DescribeStorediSCSIVolumesOutput, error) {
	stateConf := &retry.StateChangeConf{
		Pending: []string{"BOOTSTRAPPING", "CREATING", "RESTORING"},
		Target:  []string{"AVAILABLE"},
		Refresh: statusStorediSCSIVolume(ctx, conn, volumeARN),
		Timeout: storediSCSIVolumeAvailableTimeout,
	}

	outputRaw, err := stateConf.WaitForStateContext(ctx)

	if output, ok := outputRaw.(*storagegateway.DescribeStorediSCSIVolumesOutput); ok {
		return output, err
	}

	return nil, err
}

func waitNFSFileShareCreated(ctx context.Context, conn *storagegateway.StorageGateway, arn string, timeout time.Duration) (*storagegateway.NFSFileShareInfo, error) {
	stateConf := &retry.StateChangeConf{
		Pending: []string{fileShareStatusCreating},
		Target:  []string{fileShareStatusAvailable},
		Refresh: statusNFSFileShare(ctx, conn, arn),
		Timeout: timeout,
		Delay:   nfsFileShareAvailableDelay,
	}

	outputRaw, err := stateConf.WaitForStateContext(ctx)

	if output, ok := outputRaw.(*storagegateway.NFSFileShareInfo); ok {
		return output, err
	}

	return nil, err
}

func waitNFSFileShareDeleted(ctx context.Context, conn *storagegateway.StorageGateway, arn string, timeout time.Duration) (*storagegateway.NFSFileShareInfo, error) {
	stateConf := &retry.StateChangeConf{
		Pending:        []string{fileShareStatusAvailable, fileShareStatusDeleting, fileShareStatusForceDeleting},
		Target:         []string{},
		Refresh:        statusNFSFileShare(ctx, conn, arn),
		Timeout:        timeout,
		Delay:          nfsFileShareDeletedDelay,
		NotFoundChecks: 1,
	}

	outputRaw, err := stateConf.WaitForStateContext(ctx)

	if output, ok := outputRaw.(*storagegateway.NFSFileShareInfo); ok {
		return output, err
	}

	return nil, err
}

func waitNFSFileShareUpdated(ctx context.Context, conn *storagegateway.StorageGateway, arn string, timeout time.Duration) (*storagegateway.NFSFileShareInfo, error) {
	stateConf := &retry.StateChangeConf{
		Pending: []string{fileShareStatusUpdating},
		Target:  []string{fileShareStatusAvailable},
		Refresh: statusNFSFileShare(ctx, conn, arn),
		Timeout: timeout,
		Delay:   nfsFileShareAvailableDelay,
	}

	outputRaw, err := stateConf.WaitForStateContext(ctx)

	if output, ok := outputRaw.(*storagegateway.NFSFileShareInfo); ok {
		return output, err
	}

	return nil, err
}

func waitSMBFileShareCreated(ctx context.Context, conn *storagegateway.StorageGateway, arn string, timeout time.Duration) (*storagegateway.SMBFileShareInfo, error) {
	stateConf := &retry.StateChangeConf{
		Pending: []string{fileShareStatusCreating},
		Target:  []string{fileShareStatusAvailable},
		Refresh: statusSMBFileShare(ctx, conn, arn),
		Timeout: timeout,
		Delay:   smbFileShareAvailableDelay,
	}

	outputRaw, err := stateConf.WaitForStateContext(ctx)

	if output, ok := outputRaw.(*storagegateway.SMBFileShareInfo); ok {
		return output, err
	}

	return nil, err
}

func waitSMBFileShareDeleted(ctx context.Context, conn *storagegateway.StorageGateway, arn string, timeout time.Duration) (*storagegateway.SMBFileShareInfo, error) {
	stateConf := &retry.StateChangeConf{
		Pending:        []string{fileShareStatusAvailable, fileShareStatusDeleting, fileShareStatusForceDeleting},
		Target:         []string{},
		Refresh:        statusSMBFileShare(ctx, conn, arn),
		Timeout:        timeout,
		Delay:          smbFileShareDeletedDelay,
		NotFoundChecks: 1,
	}

	outputRaw, err := stateConf.WaitForStateContext(ctx)

	if output, ok := outputRaw.(*storagegateway.SMBFileShareInfo); ok {
		return output, err
	}

	return nil, err
}

func waitSMBFileShareUpdated(ctx context.Context, conn *storagegateway.StorageGateway, arn string, timeout time.Duration) (*storagegateway.SMBFileShareInfo, error) {
	stateConf := &retry.StateChangeConf{
		Pending: []string{fileShareStatusUpdating},
		Target:  []string{fileShareStatusAvailable},
		Refresh: statusSMBFileShare(ctx, conn, arn),
		Timeout: timeout,
		Delay:   smbFileShareAvailableDelay,
	}

	outputRaw, err := stateConf.WaitForStateContext(ctx)

	if output, ok := outputRaw.(*storagegateway.SMBFileShareInfo); ok {
		return output, err
	}

	return nil, err
}

func waitFileSystemAssociationAvailable(ctx context.Context, conn *storagegateway.StorageGateway, fileSystemArn string, timeout time.Duration) (*storagegateway.FileSystemAssociationInfo, error) { //nolint:unparam
	stateConf := &retry.StateChangeConf{
		Pending: []string{fileSystemAssociationStatusCreating, fileSystemAssociationStatusUpdating},
		Target:  []string{fileSystemAssociationStatusAvailable},
		Refresh: statusFileSystemAssociation(ctx, conn, fileSystemArn),
		Timeout: timeout,
		Delay:   fileSystemAssociationAvailableDelay,
	}

	outputRaw, err := stateConf.WaitForStateContext(ctx)

	if output, ok := outputRaw.(*storagegateway.FileSystemAssociationInfo); ok {
		return output, err
	}

	return nil, err
}

func waitFileSystemAssociationDeleted(ctx context.Context, conn *storagegateway.StorageGateway, fileSystemArn string, timeout time.Duration) (*storagegateway.FileSystemAssociationInfo, error) {
	stateConf := &retry.StateChangeConf{
		Pending:        []string{fileSystemAssociationStatusAvailable, fileSystemAssociationStatusDeleting, fileSystemAssociationStatusForceDeleting},
		Target:         []string{},
		Refresh:        statusFileSystemAssociation(ctx, conn, fileSystemArn),
		Timeout:        timeout,
		Delay:          fileSystemAssociationDeletedDelay,
		NotFoundChecks: 1,
	}

	outputRaw, err := stateConf.WaitForStateContext(ctx)

	if output, ok := outputRaw.(*storagegateway.FileSystemAssociationInfo); ok {
		return output, err
	}

	return nil, err
}
