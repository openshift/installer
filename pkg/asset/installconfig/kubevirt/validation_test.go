package kubevirt

import (
	"fmt"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"

	"github.com/openshift/installer/pkg/asset/installconfig/kubevirt/mock"
	"github.com/openshift/installer/pkg/ipnet"
	"github.com/openshift/installer/pkg/types"
	"github.com/openshift/installer/pkg/types/kubevirt"
	corev1 "k8s.io/api/core/v1"
)

var (
	validNamespace        = "valid-namespace"
	validStorageClass     = "valid-storage-class"
	validNetworkName      = "valid-network-name"
	validAPIVIP           = "192.168.123.15"
	validIngressVIP       = "192.168.123.20"
	validAccessMode       = "valid-access-mode"
	validMachineCIDR      = "192.168.123.0/24"
	invalidKubeconfigPath = "invalid-kubeconfig-path"
	invalidNamespace      = "invalid-namespace"
	invalidStorageClass   = "invalid-storage-class"
	invalidNetworkName    = "invalid-network-name"
	invalidAPIVIP         = "invalid-api-vip"
	invalidIngressVIP     = "invalid-ingress-vip"
	invalidAccessMode     = "invalid-access-mode"
	invalidMachineCIDR    = "10.0.0.0/16"
	namespaceStruct       = &corev1.Namespace{}
	kubeMacPoolLabels     = map[string]string{"mutatevirtualmachines.kubemacpool.io": "allocate"}
)

func validInstallConfig() *types.InstallConfig {
	return &types.InstallConfig{
		Networking: &types.Networking{
			MachineNetwork: []types.MachineNetworkEntry{
				{CIDR: *ipnet.MustParseCIDR(invalidMachineCIDR)},
				{CIDR: *ipnet.MustParseCIDR(validMachineCIDR)},
				{CIDR: *ipnet.MustParseCIDR(invalidMachineCIDR)},
			},
		},
		Platform: types.Platform{
			Kubevirt: &kubevirt.Platform{
				Namespace:                  validNamespace,
				StorageClass:               validStorageClass,
				NetworkName:                validNetworkName,
				APIVIP:                     validAPIVIP,
				IngressVIP:                 validIngressVIP,
				PersistentVolumeAccessMode: validAccessMode,
			},
		},
	}
}

func TestKubevirtInstallConfigValidation(t *testing.T) {
	cases := []struct {
		name           string
		edit           func(ic *types.InstallConfig)
		expectedError  bool
		expectedErrMsg string
		expectClient   func(kubevirtClient *mock.MockClient)
	}{
		{
			name:           "valid",
			edit:           nil,
			expectedError:  false,
			expectedErrMsg: "",
			expectClient: func(kubevirtClient *mock.MockClient) {
				kubevirtClient.EXPECT().GetNetworkAttachmentDefinition(gomock.Any(), validNetworkName, validNamespace).Return(nil, nil).AnyTimes()
				kubevirtClient.EXPECT().GetStorageClass(gomock.Any(), validStorageClass).Return(nil, nil).AnyTimes()
				namespaceStruct.Labels = kubeMacPoolLabels
				kubevirtClient.EXPECT().GetNamespace(gomock.Any(), validNamespace).Return(namespaceStruct, nil).AnyTimes()
			},
		},
		{
			name: "valid one machine network",
			edit: func(ic *types.InstallConfig) {
				ic.Networking.MachineNetwork = []types.MachineNetworkEntry{
					{CIDR: *ipnet.MustParseCIDR(validMachineCIDR)},
				}
			},
			expectedError:  false,
			expectedErrMsg: "",
			expectClient: func(kubevirtClient *mock.MockClient) {
				kubevirtClient.EXPECT().GetNetworkAttachmentDefinition(gomock.Any(), validNetworkName, validNamespace).Return(nil, nil).AnyTimes()
				kubevirtClient.EXPECT().GetStorageClass(gomock.Any(), validStorageClass).Return(nil, nil).AnyTimes()
				namespaceStruct.Labels = kubeMacPoolLabels
				kubevirtClient.EXPECT().GetNamespace(gomock.Any(), validNamespace).Return(namespaceStruct, nil).AnyTimes()
			},
		},
		{
			name:           "invalid storage class",
			edit:           func(ic *types.InstallConfig) { ic.Platform.Kubevirt.StorageClass = invalidStorageClass },
			expectedError:  true,
			expectedErrMsg: "platform.kubevirt.storageClass: Invalid value: \"invalid-storage-class\": failed to get StorageClass from InfraCluster, with error: test",
			expectClient: func(kubevirtClient *mock.MockClient) {
				kubevirtClient.EXPECT().GetNetworkAttachmentDefinition(gomock.Any(), validNetworkName, validNamespace).Return(nil, nil).AnyTimes()
				kubevirtClient.EXPECT().GetStorageClass(gomock.Any(), invalidStorageClass).Return(nil, fmt.Errorf("test")).AnyTimes()
				namespaceStruct.Labels = kubeMacPoolLabels
				kubevirtClient.EXPECT().GetNamespace(gomock.Any(), validNamespace).Return(namespaceStruct, nil).AnyTimes()
			},
		},
		{
			name:           "invalid network name",
			edit:           func(ic *types.InstallConfig) { ic.Platform.Kubevirt.NetworkName = invalidNetworkName },
			expectedError:  true,
			expectedErrMsg: "platform.kubevirt.networkName: Invalid value: \"invalid-network-name\": failed to get network-attachment-definition from InfraCluster, with error: test",
			expectClient: func(kubevirtClient *mock.MockClient) {
				kubevirtClient.EXPECT().GetNetworkAttachmentDefinition(gomock.Any(), invalidNetworkName, validNamespace).Return(nil, fmt.Errorf("test")).AnyTimes()
				kubevirtClient.EXPECT().GetStorageClass(gomock.Any(), validStorageClass).Return(nil, nil).AnyTimes()
				namespaceStruct.Labels = kubeMacPoolLabels
				kubevirtClient.EXPECT().GetNamespace(gomock.Any(), validNamespace).Return(namespaceStruct, nil).AnyTimes()
			},
		},
		{
			name:           "invalid kubemacpool namespace not found",
			edit:           nil,
			expectedError:  true,
			expectedErrMsg: "platform.kubevirt.namespace: Invalid value: \"valid-namespace\": failed to get namepsace, with error: test",
			expectClient: func(kubevirtClient *mock.MockClient) {
				kubevirtClient.EXPECT().GetNetworkAttachmentDefinition(gomock.Any(), validNetworkName, validNamespace).Return(nil, nil).AnyTimes()
				kubevirtClient.EXPECT().GetStorageClass(gomock.Any(), validStorageClass).Return(nil, nil).AnyTimes()
				namespaceStruct.Labels = kubeMacPoolLabels
				kubevirtClient.EXPECT().GetNamespace(gomock.Any(), validNamespace).Return(nil, fmt.Errorf("test")).AnyTimes()
			},
		},
		{
			name:           "invalid kubemacpool Labels nil",
			edit:           nil,
			expectedError:  true,
			expectedErrMsg: "platform.kubevirt.namespace: Invalid value: \"valid-namespace\": KubeMacPool component is not enabled for the namespace, the namespace must have label \"mutatevirtualmachines.kubemacpool.io: allocate\"",
			expectClient: func(kubevirtClient *mock.MockClient) {
				kubevirtClient.EXPECT().GetNetworkAttachmentDefinition(gomock.Any(), validNetworkName, validNamespace).Return(nil, nil).AnyTimes()
				kubevirtClient.EXPECT().GetStorageClass(gomock.Any(), validStorageClass).Return(nil, nil).AnyTimes()
				namespaceStruct.Labels = nil
				kubevirtClient.EXPECT().GetNamespace(gomock.Any(), validNamespace).Return(namespaceStruct, nil).AnyTimes()
			},
		},
		{
			name:           "invalid kubemacpool Labels empty",
			edit:           nil,
			expectedError:  true,
			expectedErrMsg: "platform.kubevirt.namespace: Invalid value: \"valid-namespace\": KubeMacPool component is not enabled for the namespace, the namespace must have label \"mutatevirtualmachines.kubemacpool.io: allocate\"",
			expectClient: func(kubevirtClient *mock.MockClient) {
				kubevirtClient.EXPECT().GetNetworkAttachmentDefinition(gomock.Any(), validNetworkName, validNamespace).Return(nil, nil).AnyTimes()
				kubevirtClient.EXPECT().GetStorageClass(gomock.Any(), validStorageClass).Return(nil, nil).AnyTimes()
				namespaceStruct.Labels = map[string]string{}
				kubevirtClient.EXPECT().GetNamespace(gomock.Any(), validNamespace).Return(namespaceStruct, nil).AnyTimes()
			},
		},
		{
			name:           "invalid kubemacpool wrong label val",
			edit:           nil,
			expectedError:  true,
			expectedErrMsg: "platform.kubevirt.namespace: Invalid value: \"valid-namespace\": KubeMacPool component is not enabled for the namespace, the namespace must have label \"mutatevirtualmachines.kubemacpool.io: allocate\"",
			expectClient: func(kubevirtClient *mock.MockClient) {
				kubevirtClient.EXPECT().GetNetworkAttachmentDefinition(gomock.Any(), validNetworkName, validNamespace).Return(nil, nil).AnyTimes()
				kubevirtClient.EXPECT().GetStorageClass(gomock.Any(), validStorageClass).Return(nil, nil).AnyTimes()
				namespaceStruct.Labels = map[string]string{"mutatevirtualmachines.kubemacpool.io": "wrong value"}
				kubevirtClient.EXPECT().GetNamespace(gomock.Any(), validNamespace).Return(namespaceStruct, nil).AnyTimes()
			},
		},
	}
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			installConfig := validInstallConfig()
			if tc.edit != nil {
				tc.edit(installConfig)
			}

			kubevirtClient := mock.NewMockClient(mockCtrl)
			if tc.expectClient != nil {
				tc.expectClient(kubevirtClient)
			}

			errs := Validate(installConfig, kubevirtClient)
			if tc.expectedError {
				assert.Regexp(t, tc.expectedErrMsg, errs)
			} else {
				assert.Empty(t, errs)
			}
		})
	}
}
