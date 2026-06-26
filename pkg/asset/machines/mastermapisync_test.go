package machines

import (
	"testing"

	"github.com/stretchr/testify/assert"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/utils/ptr"
	capa "sigs.k8s.io/cluster-api-provider-aws/v2/api/v1beta2"

	machinev1beta1 "github.com/openshift/api/machine/v1beta1"
	"github.com/openshift/installer/pkg/asset"
)

func makeTestMAPIConfig(instanceType string, volumeSize int64, volumeType string) *machinev1beta1.AWSMachineProviderConfig {
	return &machinev1beta1.AWSMachineProviderConfig{
		TypeMeta: metav1.TypeMeta{
			APIVersion: "machine.openshift.io/v1beta1",
			Kind:       "AWSMachineProviderConfig",
		},
		InstanceType: instanceType,
		AMI: machinev1beta1.AWSResourceReference{
			ID: ptr.To("ami-12345"),
		},
		BlockDevices: []machinev1beta1.BlockDeviceMappingSpec{
			{
				EBS: &machinev1beta1.EBSBlockDeviceSpec{
					VolumeSize: ptr.To(volumeSize),
					VolumeType: ptr.To(volumeType),
					Iops:       ptr.To(int64(3000)),
					Encrypted:  ptr.To(true),
				},
			},
		},
		IAMInstanceProfile: &machinev1beta1.AWSResourceReference{
			ID: ptr.To("test-cluster-master-profile"),
		},
		MetadataServiceOptions: machinev1beta1.MetadataServiceOptions{
			Authentication: machinev1beta1.MetadataServiceAuthentication("Optional"),
		},
		Placement: machinev1beta1.Placement{
			Region:           "us-east-1",
			AvailabilityZone: "us-east-1a",
		},
		Subnet: machinev1beta1.AWSResourceReference{
			Filters: []machinev1beta1.Filter{
				{Name: "tag:Name", Values: []string{"test-cluster-subnet-private-us-east-1a"}},
			},
		},
	}
}

func makeTestCAPIAWSMachine(name, instanceType string, volumeSize int64, volumeType string) *capa.AWSMachine {
	return &capa.AWSMachine{
		TypeMeta: metav1.TypeMeta{
			APIVersion: "infrastructure.cluster.x-k8s.io/v1beta2",
			Kind:       "AWSMachine",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name: name,
			Labels: map[string]string{
				"cluster.x-k8s.io/control-plane": "",
			},
		},
		Spec: capa.AWSMachineSpec{
			InstanceType:       instanceType,
			AMI:                capa.AMIReference{ID: ptr.To("ami-12345")},
			IAMInstanceProfile: "test-cluster-master-profile",
			PublicIP:           ptr.To(false),
			SSHKeyName:         ptr.To(""),
			RootVolume: &capa.Volume{
				Size:      volumeSize,
				Type:      capa.VolumeType(volumeType),
				IOPS:      3000,
				Encrypted: ptr.To(true),
			},
			InstanceMetadataOptions: &capa.InstanceMetadataOptions{
				HTTPTokens:   capa.HTTPTokensStateOptional,
				HTTPEndpoint: capa.InstanceMetadataEndpointStateEnabled,
			},
			Subnet: &capa.AWSResourceReference{
				Filters: []capa.Filter{
					{Name: "tag:Name", Values: []string{"test-cluster-subnet-private-us-east-1a"}},
				},
			},
		},
	}
}


func TestSyncAWSFields_NoDrift(t *testing.T) {
	mapi := makeTestMAPIConfig("m6i.xlarge", 120, "gp3")
	capi := makeTestCAPIAWSMachine("test-cluster-master-0", "m6i.xlarge", 120, "gp3")

	drifts := syncAWSFields(mapi, capi)

	assert.Empty(t, drifts, "expected no drift when specs match")
	assert.Equal(t, "m6i.xlarge", capi.Spec.InstanceType)
	assert.Equal(t, int64(120), capi.Spec.RootVolume.Size)
}

func TestSyncAWSFields_InstanceTypeDrift(t *testing.T) {
	mapi := makeTestMAPIConfig("m6i.4xlarge", 120, "gp3")
	capi := makeTestCAPIAWSMachine("test-cluster-master-0", "m6i.xlarge", 120, "gp3")

	drifts := syncAWSFields(mapi, capi)

	assert.Len(t, drifts, 1)
	assert.Contains(t, drifts[0], "instanceType")
	assert.Contains(t, drifts[0], "m6i.4xlarge")
	assert.Equal(t, "m6i.4xlarge", capi.Spec.InstanceType)
}

func TestSyncAWSFields_RootVolumeDrift(t *testing.T) {
	mapi := makeTestMAPIConfig("m6i.xlarge", 200, "io1")
	mapi.BlockDevices[0].EBS.Iops = ptr.To(int64(5000))
	mapi.BlockDevices[0].EBS.ThroughputMib = ptr.To(int32(500))

	capi := makeTestCAPIAWSMachine("test-cluster-master-0", "m6i.xlarge", 120, "gp3")

	drifts := syncAWSFields(mapi, capi)

	assert.Len(t, drifts, 4) // size, type, iops, throughput
	assert.Equal(t, int64(200), capi.Spec.RootVolume.Size)
	assert.Equal(t, capa.VolumeType("io1"), capi.Spec.RootVolume.Type)
	assert.Equal(t, int64(5000), capi.Spec.RootVolume.IOPS)
	assert.Equal(t, ptr.To(int64(500)), capi.Spec.RootVolume.Throughput)
}

func TestSyncAWSFields_MultipleFieldDrift(t *testing.T) {
	mapi := makeTestMAPIConfig("m6i.4xlarge", 200, "io2")
	mapi.AMI.ID = ptr.To("ami-custom")
	mapi.IAMInstanceProfile = &machinev1beta1.AWSResourceReference{
		ID: ptr.To("custom-profile"),
	}
	mapi.MetadataServiceOptions.Authentication = "Required"
	mapi.PublicIP = ptr.To(true)

	capi := makeTestCAPIAWSMachine("test-cluster-master-0", "m6i.xlarge", 120, "gp3")

	drifts := syncAWSFields(mapi, capi)

	assert.Greater(t, len(drifts), 4)
	assert.Equal(t, "m6i.4xlarge", capi.Spec.InstanceType)
	assert.Equal(t, ptr.To("ami-custom"), capi.Spec.AMI.ID)
	assert.Equal(t, int64(200), capi.Spec.RootVolume.Size)
	assert.Equal(t, capa.VolumeType("io2"), capi.Spec.RootVolume.Type)
	assert.Equal(t, "custom-profile", capi.Spec.IAMInstanceProfile)
	assert.Equal(t, capa.HTTPTokensState("required"), capi.Spec.InstanceMetadataOptions.HTTPTokens)
	assert.Equal(t, ptr.To(true), capi.Spec.PublicIP)
}

func TestSyncAWSFields_AMIDrift(t *testing.T) {
	mapi := makeTestMAPIConfig("m6i.xlarge", 120, "gp3")
	mapi.AMI.ID = ptr.To("ami-custom-image")

	capi := makeTestCAPIAWSMachine("test-cluster-master-0", "m6i.xlarge", 120, "gp3")
	capi.Spec.AMI.ID = ptr.To("ami-12345")

	drifts := syncAWSFields(mapi, capi)

	assert.Len(t, drifts, 1)
	assert.Contains(t, drifts[0], "ami.id")
	assert.Equal(t, ptr.To("ami-custom-image"), capi.Spec.AMI.ID)
}

func TestSyncAWSFields_SubnetIDDrift(t *testing.T) {
	mapi := makeTestMAPIConfig("m6i.xlarge", 120, "gp3")
	mapi.Subnet = machinev1beta1.AWSResourceReference{
		ID: ptr.To("subnet-abc123"),
	}

	capi := makeTestCAPIAWSMachine("test-cluster-master-0", "m6i.xlarge", 120, "gp3")

	drifts := syncAWSFields(mapi, capi)

	assert.Len(t, drifts, 1)
	assert.Contains(t, drifts[0], "subnet.id")
	assert.Equal(t, ptr.To("subnet-abc123"), capi.Spec.Subnet.ID)
	assert.Nil(t, capi.Spec.Subnet.Filters)
}

func TestSyncAWSFields_KMSKeyDrift(t *testing.T) {
	mapi := makeTestMAPIConfig("m6i.xlarge", 120, "gp3")
	mapi.BlockDevices[0].EBS.KMSKey = machinev1beta1.AWSResourceReference{
		ARN: ptr.To("arn:aws:kms:us-east-1:123456789:key/my-key"),
	}

	capi := makeTestCAPIAWSMachine("test-cluster-master-0", "m6i.xlarge", 120, "gp3")

	drifts := syncAWSFields(mapi, capi)

	assert.Len(t, drifts, 1)
	assert.Contains(t, drifts[0], "rootVolume.encryptionKey")
	assert.Equal(t, "arn:aws:kms:us-east-1:123456789:key/my-key", capi.Spec.RootVolume.EncryptionKey)
}

func TestIndexAWSMachinesByName(t *testing.T) {
	master0 := makeTestCAPIAWSMachine("test-cluster-master-0", "m6i.xlarge", 120, "gp3")
	master1 := makeTestCAPIAWSMachine("test-cluster-master-1", "m6i.xlarge", 120, "gp3")
	bootstrap := makeTestCAPIAWSMachine("test-cluster-bootstrap", "m6i.xlarge", 120, "gp3")

	files := []*asset.RuntimeFile{
		{Object: master0},
		{Object: master1},
		{Object: bootstrap},
	}

	result := indexAWSMachinesByName(files)

	assert.Len(t, result, 2)
	assert.Contains(t, result, "test-cluster-master-0")
	assert.Contains(t, result, "test-cluster-master-1")
	assert.NotContains(t, result, "test-cluster-bootstrap")
}

func TestSyncAWSFields_CPUOptionsDrift(t *testing.T) {
	mapi := makeTestMAPIConfig("m6i.xlarge", 120, "gp3")
	cc := machinev1beta1.AWSConfidentialComputePolicy("AMDEncryptedVirtualizationNestedPaging")
	mapi.CPUOptions = &machinev1beta1.CPUOptions{
		ConfidentialCompute: &cc,
	}

	capi := makeTestCAPIAWSMachine("test-cluster-master-0", "m6i.xlarge", 120, "gp3")

	drifts := syncAWSFields(mapi, capi)

	assert.Len(t, drifts, 1)
	assert.Contains(t, drifts[0], "cpuOptions.confidentialCompute")
	assert.Equal(t, capa.AWSConfidentialComputePolicy("AMDEncryptedVirtualizationNestedPaging"), capi.Spec.CPUOptions.ConfidentialCompute)
}

func TestFindAWSMachineByIndex(t *testing.T) {
	machines := map[string]*capa.AWSMachine{
		"test-cluster-master-0": makeTestCAPIAWSMachine("test-cluster-master-0", "m6i.xlarge", 120, "gp3"),
		"test-cluster-master-1": makeTestCAPIAWSMachine("test-cluster-master-1", "m6i.xlarge", 120, "gp3"),
	}

	m, found := findAWSMachineByIndex(machines, 0)
	assert.True(t, found)
	assert.Equal(t, "test-cluster-master-0", m.Name)

	m, found = findAWSMachineByIndex(machines, 1)
	assert.True(t, found)
	assert.Equal(t, "test-cluster-master-1", m.Name)

	_, found = findAWSMachineByIndex(machines, 5)
	assert.False(t, found)
}
