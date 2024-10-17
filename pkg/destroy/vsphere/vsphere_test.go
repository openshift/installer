package vsphere

import (
	"context"
	"errors"
	"fmt"
	"io"
	"testing"

	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"github.com/vmware/govmomi/vim25/mo"
	vspheretypes "github.com/vmware/govmomi/vim25/types"
	gomock "go.uber.org/mock/gomock"

	"github.com/openshift/installer/pkg/destroy/vsphere/mock"
	"github.com/openshift/installer/pkg/types"
	"github.com/openshift/installer/pkg/types/vsphere"
)

type editMetadataFuncs []func(m *types.ClusterMetadata)

type testCase struct {
	name      string
	editFuncs editMetadataFuncs
	errorMsg  string
}

const (
	infraID       = "infra-id"
	listFailsID   = "list-fails-infra-id"
	deleteFailsID = "delete-fails-infra-id"
	stopFailsID   = "stop-VM-fails-infra-id"
)

var (
	nullLogger = func() logrus.FieldLogger {
		logger := logrus.StandardLogger()
		logger.SetOutput(io.Discard)
		return logger
	}()
	runningVM = func() mo.VirtualMachine {
		vm := mo.VirtualMachine{}
		vm.Name = "runningVM"
		vm.Summary.Runtime.PowerState = "poweredOn"
		return vm
	}()
	stoppedVM = func() mo.VirtualMachine {
		vm := mo.VirtualMachine{}
		vm.Name = "stoppedVM"
		vm.Summary.Runtime.PowerState = "poweredOff"
		return vm
	}()
	failVM = func() mo.VirtualMachine {
		vm := mo.VirtualMachine{}
		vm.Name = "failVM"
		vm.Summary.Runtime.PowerState = "unknown"
		return vm
	}()
)

func newDefaultMetadata() types.ClusterMetadata {
	metadata := types.ClusterMetadata{
		ClusterName: "cluster-name",
		ClusterID:   "cluster-id",
		InfraID:     infraID,
	}
	metadata.VSphere = &vsphere.Metadata{
		VCenter:           "vCenter",
		Username:          "username",
		Password:          "password",
		TerraformPlatform: "vsphere",
	}
	return metadata
}

func TestVsphereDeleteFolder(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	vsphereClient := mock.NewMockAPI(mockCtrl)

	const (
		oneFolderID    = "folder-infra-id"
		manyFoldersID  = "folders-infra-id"
		manyChildrenID = "children-infra-id"
	)

	validFolder := mo.Folder{}
	validFolder.Name = "valid-folder"
	childrenFolder := mo.Folder{}
	childrenFolder.Name = "invalid-folder"
	childrenFolder.ChildEntity = []vspheretypes.ManagedObjectReference{
		{Type: "child-type", Value: "child-value"},
		{Type: "child-type", Value: "child-value"},
	}
	failFolder := mo.Folder{}
	failFolder.Name = "fail-folder"

	manyFoldersPerTag := func(m *types.ClusterMetadata) {
		m.InfraID = manyFoldersID
	}
	oneFolderPerTag := func(m *types.ClusterMetadata) {
		m.InfraID = oneFolderID
	}
	manyChildrenFolder := func(m *types.ClusterMetadata) {
		m.InfraID = manyChildrenID
	}
	listFails := func(m *types.ClusterMetadata) {
		m.InfraID = listFailsID
	}
	deleteFails := func(m *types.ClusterMetadata) {
		m.InfraID = deleteFailsID
	}

	cases := []testCase{
		{
			name:      "Delete folder when no folder present",
			editFuncs: editMetadataFuncs{},
			errorMsg:  "",
		},
		{
			name:      "Delete empty folder",
			editFuncs: editMetadataFuncs{oneFolderPerTag},
			errorMsg:  "",
		},
		{
			name:      "Delete non-empty folder fails",
			editFuncs: editMetadataFuncs{manyChildrenFolder},
			errorMsg:  "Expected Folder .* to be empty",
		},
		{
			name:      "Delete folder fails when listing",
			editFuncs: editMetadataFuncs{listFails},
			errorMsg:  "list attached objects",
		},
		{
			name:      "Delete folders Zoning Terraform",
			editFuncs: editMetadataFuncs{manyFoldersPerTag},
			errorMsg:  "",
		},
		{
			name:      "Delete folder fails",
			editFuncs: editMetadataFuncs{deleteFails},
			errorMsg:  "some vsphere error",
		},
	}

	vsphereClient.
		EXPECT().
		ListFolders(gomock.Any(), gomock.Eq(infraID)).
		Return([]mo.Folder{}, nil).
		AnyTimes()
	vsphereClient.
		EXPECT().
		ListFolders(gomock.Any(), gomock.Eq(oneFolderID)).
		Return([]mo.Folder{validFolder}, nil).
		AnyTimes()
	vsphereClient.
		EXPECT().
		ListFolders(gomock.Any(), gomock.Eq(manyFoldersID)).
		Return([]mo.Folder{validFolder, validFolder}, nil).
		AnyTimes()
	vsphereClient.
		EXPECT().
		ListFolders(gomock.Any(), gomock.Eq(manyChildrenID)).
		Return([]mo.Folder{validFolder, childrenFolder, failFolder}, nil).
		AnyTimes()
	vsphereClient.
		EXPECT().
		ListFolders(gomock.Any(), gomock.Eq(deleteFailsID)).
		Return([]mo.Folder{validFolder, failFolder, validFolder, childrenFolder}, nil).
		AnyTimes()
	vsphereClient.
		EXPECT().
		ListFolders(gomock.Any(), gomock.Any()).
		Return(nil, errors.New("list attached objects infra-id: vsphere error")).
		AnyTimes()

	vsphereClient.
		EXPECT().
		DeleteFolder(gomock.Any(), gomock.Eq(validFolder)).
		Return(nil).
		AnyTimes()
	vsphereClient.
		EXPECT().
		DeleteFolder(gomock.Any(), gomock.Eq(childrenFolder)).
		Times(0) // Should not delete a folder with children
	vsphereClient.
		EXPECT().
		DeleteFolder(gomock.Any(), gomock.Any()).
		Return(errors.New("some vsphere error deleting Folder")).
		AnyTimes()

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			editedMetadata := newDefaultMetadata()
			for _, edit := range tc.editFuncs {
				edit(&editedMetadata)
			}
			uninstaller := newWithClient(nullLogger, &editedMetadata, []API{vsphereClient})
			assert.NotNil(t, uninstaller)
			err := uninstaller.deleteFolder(context.TODO())
			if tc.errorMsg != "" {
				assert.Regexp(t, tc.errorMsg, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestVsphereStopVirtualMachines(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	vsphereClient := mock.NewMockAPI(mockCtrl)

	const (
		stoppedVMsID = "stopped-VM-infra-id"
		runningVMsID = "running-VM-infra-id"
		mixedVMsID   = "mixed-VM-infra-id"
	)

	stoppedVMs := func(m *types.ClusterMetadata) {
		m.InfraID = stoppedVMsID
	}
	runningVMs := func(m *types.ClusterMetadata) {
		m.InfraID = runningVMsID
	}
	mixedVMs := func(m *types.ClusterMetadata) {
		m.InfraID = mixedVMsID
	}
	listFails := func(m *types.ClusterMetadata) {
		m.InfraID = listFailsID
	}
	stopFails := func(m *types.ClusterMetadata) {
		m.InfraID = stopFailsID
	}

	cases := []testCase{
		{
			name:      "Stop VMs when no VM present",
			editFuncs: editMetadataFuncs{},
			errorMsg:  "",
		},
		{
			name:      "Stop VMs when none running",
			editFuncs: editMetadataFuncs{stoppedVMs},
			errorMsg:  "",
		},
		{
			name:      "Stop VMs when all running",
			editFuncs: editMetadataFuncs{runningVMs},
			errorMsg:  "",
		},
		{
			name:      "Stop VMs when some running",
			editFuncs: editMetadataFuncs{mixedVMs},
			errorMsg:  "",
		},
		{
			name:      "Stop VMs fails when listing fails",
			editFuncs: editMetadataFuncs{listFails},
			errorMsg:  "some vsphere error",
		},
		{
			name:      "Stop VMs fails",
			editFuncs: editMetadataFuncs{stopFails},
			errorMsg:  "some vsphere error",
		},
	}

	vsphereClient.
		EXPECT().
		ListVirtualMachines(gomock.Any(), gomock.Eq(infraID)).
		Return([]mo.VirtualMachine{}, nil).
		AnyTimes()
	vsphereClient.
		EXPECT().
		ListVirtualMachines(gomock.Any(), gomock.Eq(stoppedVMsID)).
		Return([]mo.VirtualMachine{stoppedVM, stoppedVM}, nil).
		AnyTimes()
	vsphereClient.
		EXPECT().
		ListVirtualMachines(gomock.Any(), gomock.Eq(runningVMsID)).
		Return([]mo.VirtualMachine{runningVM, runningVM}, nil).
		AnyTimes()
	vsphereClient.
		EXPECT().
		ListVirtualMachines(gomock.Any(), gomock.Eq(mixedVMsID)).
		Return([]mo.VirtualMachine{runningVM, runningVM, stoppedVM}, nil).
		AnyTimes()
	vsphereClient.
		EXPECT().
		ListVirtualMachines(gomock.Any(), gomock.Eq(stopFailsID)).
		Return([]mo.VirtualMachine{runningVM, failVM, stoppedVM, failVM}, nil).
		AnyTimes()
	vsphereClient.
		EXPECT().
		ListVirtualMachines(gomock.Any(), gomock.Eq(listFailsID)).
		Return(nil, errors.New("some vsphere error listing VMs")).
		AnyTimes()

	vsphereClient.
		EXPECT().
		StopVirtualMachine(gomock.Any(), gomock.Eq(runningVM)).
		Return(nil).
		AnyTimes()
	vsphereClient.
		EXPECT().
		StopVirtualMachine(gomock.Any(), gomock.Eq(stoppedVM)).
		Times(0) // Should not try to stop a VM that is not running
	vsphereClient.
		EXPECT().
		StopVirtualMachine(gomock.Any(), gomock.Eq(failVM)).
		Return(errors.New("some vsphere error stopping VM")).
		AnyTimes()
	vsphereClient.
		EXPECT().
		StopVirtualMachine(gomock.Any(), gomock.Any()).
		Times(0)

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			editedMetadata := newDefaultMetadata()
			for _, edit := range tc.editFuncs {
				edit(&editedMetadata)
			}
			uninstaller := newWithClient(nullLogger, &editedMetadata, []API{vsphereClient})
			assert.NotNil(t, uninstaller)
			err := uninstaller.stopVirtualMachines(context.TODO())
			if tc.errorMsg != "" {
				assert.Regexp(t, tc.errorMsg, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestVsphereDeleteVirtualMachines(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	vsphereClient := mock.NewMockAPI(mockCtrl)

	const someVMsID = "some-VM-infra-id"

	someVMs := func(m *types.ClusterMetadata) {
		m.InfraID = someVMsID
	}
	listFails := func(m *types.ClusterMetadata) {
		m.InfraID = listFailsID
	}
	deleteFails := func(m *types.ClusterMetadata) {
		m.InfraID = deleteFailsID
	}

	cases := []testCase{
		{
			name:      "Delete VMs when none present",
			editFuncs: editMetadataFuncs{},
			errorMsg:  "",
		},
		{
			name:      "Delete VMs when some present",
			editFuncs: editMetadataFuncs{someVMs},
			errorMsg:  "",
		},
		{
			name:      "Delete VMs fails when listing fails",
			editFuncs: editMetadataFuncs{listFails},
			errorMsg:  "some vsphere error",
		},
		{
			name:      "Delete VMs fails when some fail",
			editFuncs: editMetadataFuncs{deleteFails},
			errorMsg:  "some vsphere error",
		},
	}

	vsphereClient.
		EXPECT().
		ListVirtualMachines(gomock.Any(), gomock.Eq(infraID)).
		Return([]mo.VirtualMachine{}, nil).
		AnyTimes()
	vsphereClient.
		EXPECT().
		ListVirtualMachines(gomock.Any(), gomock.Eq(someVMsID)).
		Return([]mo.VirtualMachine{stoppedVM, stoppedVM}, nil).
		AnyTimes()
	vsphereClient.
		EXPECT().
		ListVirtualMachines(gomock.Any(), gomock.Eq(deleteFailsID)).
		Return([]mo.VirtualMachine{stoppedVM, failVM, stoppedVM, failVM}, nil).
		AnyTimes()
	vsphereClient.
		EXPECT().
		ListVirtualMachines(gomock.Any(), gomock.Eq(listFailsID)).
		Return(nil, errors.New("some vsphere error listing VMs")).
		AnyTimes()

	vsphereClient.
		EXPECT().
		DeleteVirtualMachine(gomock.Any(), gomock.Eq(runningVM)).
		Times(0) // Not import but should not happen
	vsphereClient.
		EXPECT().
		DeleteVirtualMachine(gomock.Any(), gomock.Eq(stoppedVM)).
		AnyTimes()
	vsphereClient.
		EXPECT().
		DeleteVirtualMachine(gomock.Any(), gomock.Eq(failVM)).
		Return(errors.New("some vsphere error deleting VM")).
		AnyTimes()
	vsphereClient.
		EXPECT().
		StopVirtualMachine(gomock.Any(), gomock.Any()).
		Times(0)

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			editedMetadata := newDefaultMetadata()
			for _, edit := range tc.editFuncs {
				edit(&editedMetadata)
			}
			uninstaller := newWithClient(nullLogger, &editedMetadata, []API{vsphereClient})
			assert.NotNil(t, uninstaller)
			err := uninstaller.deleteVirtualMachines(context.TODO())
			if tc.errorMsg != "" {
				assert.Regexp(t, tc.errorMsg, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestDeleteStoragePolicy(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	vsphereClient := mock.NewMockAPI(mockCtrl)

	policyFails := fmt.Sprintf("openshift-storage-policy-%s", deleteFailsID)

	deleteFails := func(m *types.ClusterMetadata) {
		m.InfraID = deleteFailsID
	}

	cases := []testCase{
		{
			name:      "Delete Storage Policy succeeds",
			editFuncs: editMetadataFuncs{},
			errorMsg:  "",
		},
		{
			name:      "Delete Storage Policy fails",
			editFuncs: editMetadataFuncs{deleteFails},
			errorMsg:  "some vsphere error",
		},
	}

	vsphereClient.
		EXPECT().
		DeleteStoragePolicy(gomock.Any(), gomock.Eq(policyFails)).
		Return(errors.New("some vsphere error deleting Storage Policy")).
		AnyTimes()
	vsphereClient.
		EXPECT().
		DeleteStoragePolicy(gomock.Any(), gomock.Any()).
		Return(nil).
		AnyTimes()

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			editedMetadata := newDefaultMetadata()
			for _, edit := range tc.editFuncs {
				edit(&editedMetadata)
			}
			uninstaller := newWithClient(nullLogger, &editedMetadata, []API{vsphereClient})
			assert.NotNil(t, uninstaller)
			err := uninstaller.deleteStoragePolicy(context.TODO())
			if tc.errorMsg != "" {
				assert.Regexp(t, tc.errorMsg, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestDeleteTag(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	vsphereClient := mock.NewMockAPI(mockCtrl)

	deleteFails := func(m *types.ClusterMetadata) {
		m.InfraID = deleteFailsID
	}

	cases := []testCase{
		{
			name:      "Delete Tag succeeds",
			editFuncs: editMetadataFuncs{},
			errorMsg:  "",
		},
		{
			name:      "Delete Tag fails",
			editFuncs: editMetadataFuncs{deleteFails},
			errorMsg:  "some vsphere error",
		},
	}

	vsphereClient.
		EXPECT().
		DeleteTag(gomock.Any(), gomock.Eq(deleteFailsID)).
		Return(errors.New("some vsphere error deleting Tag")).
		AnyTimes()
	vsphereClient.
		EXPECT().
		DeleteTag(gomock.Any(), gomock.Any()).
		Return(nil).
		AnyTimes()

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			editedMetadata := newDefaultMetadata()
			for _, edit := range tc.editFuncs {
				edit(&editedMetadata)
			}
			uninstaller := newWithClient(nullLogger, &editedMetadata, []API{vsphereClient})
			assert.NotNil(t, uninstaller)
			err := uninstaller.deleteTag(context.TODO())
			if tc.errorMsg != "" {
				assert.Regexp(t, tc.errorMsg, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestDeleteTagCategory(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	vsphereClient := mock.NewMockAPI(mockCtrl)

	tagFails := fmt.Sprintf("openshift-%s", deleteFailsID)

	deleteFails := func(m *types.ClusterMetadata) {
		m.InfraID = deleteFailsID
	}

	cases := []testCase{
		{
			name:      "Delete Tag Category succeeds",
			editFuncs: editMetadataFuncs{},
			errorMsg:  "",
		},
		{
			name:      "Delete Tag Category fails",
			editFuncs: editMetadataFuncs{deleteFails},
			errorMsg:  "some vsphere error",
		},
	}

	vsphereClient.
		EXPECT().
		DeleteTagCategory(gomock.Any(), gomock.Eq(tagFails)).
		Return(errors.New("some vsphere error deleting Tag Category")).
		AnyTimes()
	vsphereClient.
		EXPECT().
		DeleteTagCategory(gomock.Any(), gomock.Any()).
		Return(nil).
		AnyTimes()

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			editedMetadata := newDefaultMetadata()
			for _, edit := range tc.editFuncs {
				edit(&editedMetadata)
			}
			uninstaller := newWithClient(nullLogger, &editedMetadata, []API{vsphereClient})
			assert.NotNil(t, uninstaller)
			err := uninstaller.deleteTagCategory(context.TODO())
			if tc.errorMsg != "" {
				assert.Regexp(t, tc.errorMsg, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
