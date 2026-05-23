package machines

import (
	"testing"

	"github.com/stretchr/testify/assert"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/utils/ptr"

	machinev1beta1 "github.com/openshift/api/machine/v1beta1"
)

func makeMAPIProviderConfig(instanceType string, volumeSize int64, volumeType string) *machinev1beta1.AWSMachineProviderConfig {
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

func makeTestCPMSProviderConfig(instanceType string, volumeSize int64, volumeType string) *machinev1beta1.AWSMachineProviderConfig {
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
			Region: "us-east-1",
		},
	}
}

func TestSyncCPMSAWSFields_NoDrift(t *testing.T) {
	mapi := makeMAPIProviderConfig("m6i.xlarge", 120, "gp3")
	cpms := makeTestCPMSProviderConfig("m6i.xlarge", 120, "gp3")

	drifts := syncCPMSAWSFields(mapi, cpms)

	assert.Empty(t, drifts, "expected no drift when specs match")
	assert.Equal(t, "m6i.xlarge", cpms.InstanceType)
}

func TestSyncCPMSAWSFields_InstanceTypeDrift(t *testing.T) {
	mapi := makeMAPIProviderConfig("m6i.4xlarge", 120, "gp3")
	cpms := makeTestCPMSProviderConfig("m6i.xlarge", 120, "gp3")

	drifts := syncCPMSAWSFields(mapi, cpms)

	assert.Len(t, drifts, 1)
	assert.Contains(t, drifts[0], "instanceType")
	assert.Equal(t, "m6i.4xlarge", cpms.InstanceType)
}

func TestSyncCPMSAWSFields_RootVolumeDrift(t *testing.T) {
	mapi := makeMAPIProviderConfig("m6i.xlarge", 200, "io2")
	mapi.BlockDevices[0].EBS.Iops = ptr.To(int64(5000))

	cpms := makeTestCPMSProviderConfig("m6i.xlarge", 120, "gp3")

	drifts := syncCPMSAWSFields(mapi, cpms)

	assert.Len(t, drifts, 3) // size, type, iops
	assert.Equal(t, ptr.To(int64(200)), cpms.BlockDevices[0].EBS.VolumeSize)
	assert.Equal(t, ptr.To("io2"), cpms.BlockDevices[0].EBS.VolumeType)
	assert.Equal(t, ptr.To(int64(5000)), cpms.BlockDevices[0].EBS.Iops)
}

func TestSyncCPMSAWSFields_AMIDrift(t *testing.T) {
	mapi := makeMAPIProviderConfig("m6i.xlarge", 120, "gp3")
	mapi.AMI.ID = ptr.To("ami-custom-image")

	cpms := makeTestCPMSProviderConfig("m6i.xlarge", 120, "gp3")

	drifts := syncCPMSAWSFields(mapi, cpms)

	assert.Len(t, drifts, 1)
	assert.Contains(t, drifts[0], "ami.id")
	assert.Equal(t, ptr.To("ami-custom-image"), cpms.AMI.ID)
}

func TestSyncCPMSAWSFields_IAMProfileDrift(t *testing.T) {
	mapi := makeMAPIProviderConfig("m6i.xlarge", 120, "gp3")
	mapi.IAMInstanceProfile = &machinev1beta1.AWSResourceReference{
		ID: ptr.To("custom-iam-profile"),
	}

	cpms := makeTestCPMSProviderConfig("m6i.xlarge", 120, "gp3")

	drifts := syncCPMSAWSFields(mapi, cpms)

	assert.Len(t, drifts, 1)
	assert.Contains(t, drifts[0], "iamInstanceProfile")
	assert.Equal(t, ptr.To("custom-iam-profile"), cpms.IAMInstanceProfile.ID)
}

func TestSyncCPMSAWSFields_MetadataAuthDrift(t *testing.T) {
	mapi := makeMAPIProviderConfig("m6i.xlarge", 120, "gp3")
	mapi.MetadataServiceOptions.Authentication = "Required"

	cpms := makeTestCPMSProviderConfig("m6i.xlarge", 120, "gp3")

	drifts := syncCPMSAWSFields(mapi, cpms)

	assert.Len(t, drifts, 1)
	assert.Contains(t, drifts[0], "metadataServiceOptions.authentication")
	assert.Equal(t, machinev1beta1.MetadataServiceAuthentication("Required"), cpms.MetadataServiceOptions.Authentication)
}

func TestSyncCPMSAWSFields_MultiFieldDrift(t *testing.T) {
	mapi := makeMAPIProviderConfig("m6i.4xlarge", 200, "io2")
	mapi.AMI.ID = ptr.To("ami-custom")
	mapi.IAMInstanceProfile = &machinev1beta1.AWSResourceReference{
		ID: ptr.To("custom-profile"),
	}
	mapi.MetadataServiceOptions.Authentication = "Required"
	mapi.PublicIP = ptr.To(true)

	cpms := makeTestCPMSProviderConfig("m6i.xlarge", 120, "gp3")

	drifts := syncCPMSAWSFields(mapi, cpms)

	assert.Greater(t, len(drifts), 4)
	assert.Equal(t, "m6i.4xlarge", cpms.InstanceType)
	assert.Equal(t, ptr.To("ami-custom"), cpms.AMI.ID)
	assert.Equal(t, ptr.To(int64(200)), cpms.BlockDevices[0].EBS.VolumeSize)
	assert.Equal(t, ptr.To("io2"), cpms.BlockDevices[0].EBS.VolumeType)
	assert.Equal(t, ptr.To("custom-profile"), cpms.IAMInstanceProfile.ID)
	assert.Equal(t, machinev1beta1.MetadataServiceAuthentication("Required"), cpms.MetadataServiceOptions.Authentication)
	assert.Equal(t, ptr.To(true), cpms.PublicIP)
}

func TestSyncCPMSAWSFields_KMSKeyDrift(t *testing.T) {
	mapi := makeMAPIProviderConfig("m6i.xlarge", 120, "gp3")
	mapi.BlockDevices[0].EBS.KMSKey = machinev1beta1.AWSResourceReference{
		ARN: ptr.To("arn:aws:kms:us-east-1:123456789:key/my-key"),
	}

	cpms := makeTestCPMSProviderConfig("m6i.xlarge", 120, "gp3")

	drifts := syncCPMSAWSFields(mapi, cpms)

	assert.Len(t, drifts, 1)
	assert.Contains(t, drifts[0], "rootVolume.kmsKey")
	assert.Equal(t, ptr.To("arn:aws:kms:us-east-1:123456789:key/my-key"), cpms.BlockDevices[0].EBS.KMSKey.ARN)
}

func TestSyncCPMSAWSFields_PublicIPDrift(t *testing.T) {
	mapi := makeMAPIProviderConfig("m6i.xlarge", 120, "gp3")
	mapi.PublicIP = ptr.To(true)

	cpms := makeTestCPMSProviderConfig("m6i.xlarge", 120, "gp3")

	drifts := syncCPMSAWSFields(mapi, cpms)

	assert.Len(t, drifts, 1)
	assert.Contains(t, drifts[0], "publicIP")
	assert.Equal(t, ptr.To(true), cpms.PublicIP)
}

func TestSyncCPMSAWSFields_SubnetNotSynced(t *testing.T) {
	mapi := makeMAPIProviderConfig("m6i.xlarge", 120, "gp3")
	mapi.Subnet = machinev1beta1.AWSResourceReference{
		ID: ptr.To("subnet-abc123"),
	}

	cpms := makeTestCPMSProviderConfig("m6i.xlarge", 120, "gp3")

	drifts := syncCPMSAWSFields(mapi, cpms)

	assert.Empty(t, drifts, "subnet should not be synced to CPMS (handled by FailureDomains)")
}

func TestSyncCPMSAWSFields_PlacementNotSynced(t *testing.T) {
	mapi := makeMAPIProviderConfig("m6i.xlarge", 120, "gp3")
	mapi.Placement.AvailabilityZone = "us-east-1a"

	cpms := makeTestCPMSProviderConfig("m6i.xlarge", 120, "gp3")

	drifts := syncCPMSAWSFields(mapi, cpms)

	assert.Empty(t, drifts, "placement AZ should not be synced to CPMS (handled by FailureDomains)")
}

func TestDecodeCPMSProviderSpec_NilInput(t *testing.T) {
	result, err := decodeCPMSProviderSpec(nil)
	assert.Nil(t, result)
	assert.NoError(t, err)
}

func TestDecodeCPMSProviderSpec_FromRawBytes(t *testing.T) {
	raw := []byte(`{
		"apiVersion": "machine.openshift.io/v1beta1",
		"kind": "AWSMachineProviderConfig",
		"instanceType": "m6i.xlarge",
		"blockDevices": [{"ebs": {"volumeSize": 120, "volumeType": "gp3"}}]
	}`)
	ext := &runtime.RawExtension{Raw: raw}

	result, err := decodeCPMSProviderSpec(ext)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, "m6i.xlarge", result.InstanceType)
	assert.Len(t, result.BlockDevices, 1)
	assert.Equal(t, ptr.To(int64(120)), result.BlockDevices[0].EBS.VolumeSize)
}
