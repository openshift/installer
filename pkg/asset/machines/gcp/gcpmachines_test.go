// Package gcp generates Machine objects for gcp.
package gcp

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	compute "google.golang.org/api/compute/v1"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/utils/ptr"
	capg "sigs.k8s.io/cluster-api-provider-gcp/api/v1beta1"
	capi "sigs.k8s.io/cluster-api/api/v1beta1"

	"github.com/openshift/installer/pkg/asset/installconfig"
	gcpconsts "github.com/openshift/installer/pkg/constants/gcp"
	"github.com/openshift/installer/pkg/types"
	gcptypes "github.com/openshift/installer/pkg/types/gcp"
)

const (
	numReplicas = 3
)

func Test_GenerateMachines(t *testing.T) {
	cases := []struct {
		name              string
		installConfig     *installconfig.InstallConfig
		expectedGCPConfig *capg.GCPMachine
		expectedError     string
	}{
		{
			name:              "base configuration",
			installConfig:     getBaseInstallConfig(),
			expectedGCPConfig: getBaseGCPMachine(),
		},
		{
			name:              "additional labels",
			installConfig:     getICWithLabels(),
			expectedGCPConfig: getGCPMachineWithLabels(),
		},
		{
			name:              "onhostmaintenance",
			installConfig:     getICWithOnHostMaintenance(),
			expectedGCPConfig: getGCPMachineWithOnHostMaintenance(),
		},
		{
			name:              "confidentialcompute",
			installConfig:     getICWithConfidentialCompute(),
			expectedGCPConfig: getGCPMachineWithConfidentialCompute(),
		},
		{
			name:              "secureboot",
			installConfig:     getICWithSecureBoot(),
			expectedGCPConfig: getGCPMachineWithSecureBoot(),
		},
		{
			name:              "serviceaccount",
			installConfig:     getICWithServiceAccount(),
			expectedGCPConfig: getGCPMachineWithServiceAccount(),
		},
		{
			name:              "serviceaccount-controlplane-machine",
			installConfig:     getICWithServiceAccountControlPlaneMachine(),
			expectedGCPConfig: getGCPMachineWithServiceAccountControlPlaneMachine(),
		},
		{
			name:              "usertags",
			installConfig:     getICWithUserTags(),
			expectedGCPConfig: getGCPMachineWithTags(),
		},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			installConfig := tc.installConfig
			ic := installConfig.Config
			pool := ic.ControlPlane
			infraID := "012345678"
			rhcosImage := "rhcos-415-92-202311241643-0-gcp-x86-64"

			mpool := gcptypes.MachinePool{
				InstanceType: "n2-standard-4",
				OSDisk: gcptypes.OSDisk{
					DiskSizeGB: 128,
					DiskType:   "pd-ssd",
				},
			}
			mpool.Set(ic.Platform.GCP.DefaultMachinePlatform)
			mpool.Set(pool.Platform.GCP)
			pool.Platform.GCP = &mpool

			gcpMachines, err := GenerateMachines(
				installConfig,
				infraID,
				pool,
				rhcosImage,
			)

			if tc.expectedError != "" {
				assert.Equal(t, tc.expectedError, err.Error())
			} else {
				assert.NoError(t, err)
				assert.NotEmpty(t, gcpMachines)

				assert.Equal(t, numReplicas*2, len(gcpMachines))
				// Check first set of GCP and CAPI machines
				actualGCPMachine := gcpMachines[0].Object
				actualCapiMachine := gcpMachines[1].Object
				assert.Equal(t, tc.expectedGCPConfig, actualGCPMachine)
				assert.Equal(t, getBaseCapiMachine(), actualCapiMachine)
			}
		})
	}
}

func getBaseInstallConfig() *installconfig.InstallConfig {
	return &installconfig.InstallConfig{
		AssetBase: installconfig.AssetBase{
			Config: &types.InstallConfig{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "ocp-edge-cluster-0",
					Namespace: "cluster-0",
				},
				BaseDomain: "testing.com",
				ControlPlane: &types.MachinePool{
					Name:     "master",
					Replicas: ptr.To(int64(numReplicas)),
					Platform: types.MachinePoolPlatform{},
				},
				Platform: types.Platform{
					GCP: &gcptypes.Platform{
						ProjectID: "my-project",
						Region:    "us-east1",
					},
				},
			},
		},
	}
}

func getICWithLabels() *installconfig.InstallConfig {
	ic := getBaseInstallConfig()
	ic.Config.Platform.GCP.UserLabels = []gcptypes.UserLabel{{Key: "foo", Value: "bar"},
		{Key: "id", Value: "1234"}}
	return ic
}

func getICWithOnHostMaintenance() *installconfig.InstallConfig {
	ic := getBaseInstallConfig()
	ic.Config.Platform.GCP.DefaultMachinePlatform = &gcptypes.MachinePool{OnHostMaintenance: "Terminate"}
	return ic
}

func getICWithConfidentialCompute() *installconfig.InstallConfig {
	ic := getBaseInstallConfig()
	ic.Config.Platform.GCP.DefaultMachinePlatform = &gcptypes.MachinePool{ConfidentialCompute: "Enabled"}
	return ic
}

func getICWithSecureBoot() *installconfig.InstallConfig {
	ic := getBaseInstallConfig()
	ic.Config.Platform.GCP.DefaultMachinePlatform = &gcptypes.MachinePool{SecureBoot: "Enabled"}
	return ic
}

func getICWithServiceAccount() *installconfig.InstallConfig {
	ic := getBaseInstallConfig()
	ic.Config.Platform.GCP.DefaultMachinePlatform = &gcptypes.MachinePool{ServiceAccount: "user-service-account@some-project.iam.gserviceaccount.com"}
	return ic
}

func getICWithServiceAccountControlPlaneMachine() *installconfig.InstallConfig {
	ic := getBaseInstallConfig()
	ic.Config.Platform.GCP.DefaultMachinePlatform = &gcptypes.MachinePool{ServiceAccount: "user-service-account@some-project.iam.gserviceaccount.com"}
	ic.Config.ControlPlane = &types.MachinePool{
		Name:     "master",
		Replicas: ptr.To(int64(numReplicas)),
		Platform: types.MachinePoolPlatform{
			GCP: &gcptypes.MachinePool{
				ServiceAccount: "other-service-account@some-project.iam.gserviceaccount.com"},
		},
	}
	return ic
}

func getICWithUserTags() *installconfig.InstallConfig {
	ic := getBaseInstallConfig()
	ic.Config.Platform.GCP.UserTags = []gcptypes.UserTag{{ParentID: "my-project", Key: "foo", Value: "bar"},
		{ParentID: "other-project", Key: "id", Value: "1234"}}
	return ic
}

func getBaseGCPMachine() *capg.GCPMachine {
	subnet := "012345678-master-subnet"
	image := "rhcos-415-92-202311241643-0-gcp-x86-64"
	diskType := "pd-ssd"
	gcpMachine := &capg.GCPMachine{
		ObjectMeta: metav1.ObjectMeta{
			Name: "012345678-master-0",
			Labels: map[string]string{
				"cluster.x-k8s.io/control-plane": "",
			},
		},
		Spec: capg.GCPMachineSpec{
			InstanceType: "n2-standard-4",
			Subnet:       &subnet,
			Image:        &image,
			AdditionalLabels: capg.Labels{
				fmt.Sprintf(gcpconsts.ClusterIDLabelFmt, "012345678"): "owned",
			},
			RootDeviceSize: 128,
			RootDeviceType: ptr.To(capg.DiskType(diskType)),
			ServiceAccount: &capg.ServiceAccount{
				Email:  "012345678-m@my-project.iam.gserviceaccount.com",
				Scopes: []string{compute.CloudPlatformScope},
			},
			ResourceManagerTags: []capg.ResourceManagerTag{},
			IPForwarding:        ptr.To(capg.IPForwardingDisabled),
		},
	}
	gcpMachine.SetGroupVersionKind(capg.GroupVersion.WithKind("GCPMachine"))
	return gcpMachine
}

func getGCPMachineWithLabels() *capg.GCPMachine {
	gcpMachine := getBaseGCPMachine()
	gcpMachine.Spec.AdditionalLabels = capg.Labels{
		fmt.Sprintf(gcpconsts.ClusterIDLabelFmt, "012345678"): "owned",
		"foo": "bar",
		"id":  "1234"}
	return gcpMachine
}

func getGCPMachineWithOnHostMaintenance() *capg.GCPMachine {
	gcpMachine := getBaseGCPMachine()
	var maint capg.HostMaintenancePolicy = "Terminate"
	gcpMachine.Spec.OnHostMaintenance = &maint
	return gcpMachine
}

func getGCPMachineWithConfidentialCompute() *capg.GCPMachine {
	gcpMachine := getBaseGCPMachine()
	var cc capg.ConfidentialComputePolicy = "Enabled"
	gcpMachine.Spec.ConfidentialCompute = &cc
	return gcpMachine
}

func getGCPMachineWithSecureBoot() *capg.GCPMachine {
	gcpMachine := getBaseGCPMachine()
	secureBoot := capg.GCPShieldedInstanceConfig{SecureBoot: capg.SecureBootPolicy("Enabled")}
	gcpMachine.Spec.ShieldedInstanceConfig = &secureBoot
	return gcpMachine
}

func getGCPMachineWithServiceAccount() *capg.GCPMachine {
	gcpMachine := getBaseGCPMachine()
	gcpMachine.Spec.ServiceAccount = &capg.ServiceAccount{
		Email:  "user-service-account@some-project.iam.gserviceaccount.com",
		Scopes: []string{compute.CloudPlatformScope},
	}
	return gcpMachine
}

func getGCPMachineWithServiceAccountControlPlaneMachine() *capg.GCPMachine {
	gcpMachine := getBaseGCPMachine()
	gcpMachine.Spec.ServiceAccount = &capg.ServiceAccount{
		Email:  "other-service-account@some-project.iam.gserviceaccount.com",
		Scopes: []string{compute.CloudPlatformScope},
	}
	return gcpMachine
}

func getBaseCapiMachine() *capi.Machine {
	dataSecret := fmt.Sprintf("%s-master", "012345678")

	capiMachine := &capi.Machine{
		ObjectMeta: metav1.ObjectMeta{
			Name: "012345678-master-0",
			Labels: map[string]string{
				"cluster.x-k8s.io/control-plane": "",
			},
		},
		Spec: capi.MachineSpec{
			ClusterName: "012345678",
			Bootstrap: capi.Bootstrap{
				DataSecretName: ptr.To(dataSecret),
			},
			InfrastructureRef: v1.ObjectReference{
				APIVersion: capg.GroupVersion.String(),
				Kind:       "GCPMachine",
				Name:       "012345678-master-0",
			},
		},
	}
	capiMachine.SetGroupVersionKind(capi.GroupVersion.WithKind("Machine"))
	return capiMachine
}

func getGCPMachineWithTags() *capg.GCPMachine {
	gcpMachine := getBaseGCPMachine()
	gcpMachine.Spec.ResourceManagerTags = []capg.ResourceManagerTag{
		{ParentID: "my-project",
			Key:   "foo",
			Value: "bar"},
		{ParentID: "other-project",
			Key:   "id",
			Value: "1234"}}
	return gcpMachine
}
