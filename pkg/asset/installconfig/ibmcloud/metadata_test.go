package ibmcloud

import (
	"context"
	"fmt"
	"testing"

	"github.com/IBM/platform-services-go-sdk/iamidentityv1"
	"github.com/IBM/vpc-go-sdk/vpcv1"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
	"k8s.io/utils/ptr"

	"github.com/openshift/installer/pkg/asset/installconfig/ibmcloud/mock"
	"github.com/openshift/installer/pkg/asset/installconfig/ibmcloud/responses"
	"github.com/openshift/installer/pkg/types"
	ibmcloudtypes "github.com/openshift/installer/pkg/types/ibmcloud"
)

type editMetadata []func(m *Metadata)

var (
	// Shared test values.
	goodDomain = "domain.good.com"
	badDomain  = "domain.bad"
	region     = "us-south"

	// Account ID test values.
	newAccountID      = "new-account-id"
	existingAccountID = "existing-account-id"

	// CIS Instance CRN test values.
	newCISCRN      = "new-cis-crn"
	existingCISCRN = "existing-cis-crn"

	// DNS Instance test values.
	newDNSInstanceID       = "new-dns-instance-id"
	newDNSInstanceCRN      = "new-dns-instance-crn"
	existingDNSInstanceID  = "existing-dns-instance-id"
	existingDNSInstanceCRN = "existing-dns-instance-crn"
	unknownDNSInstanceID   = "unknown-dns-instance-id"
	unknownDNSInstanceCRN  = "unknown-dns-instance-crn"

	// Newly retrieved Compute Subnets.
	newComputeSubnet1Name     = "new-compute-subnet-1"
	newComputeSubnet1ID       = "new-compute-1-id"
	newComputeSubnet1CIDR     = "new-compute-1-cidr"
	newComputeSubnet1CRN      = "new-compute-1-crn"
	newComputeSubnet1VPCName  = "new-compute-1-vpc"
	newComputeSubnet1ZoneName = "new-compute-1-zone"
	newComputeSubnet2Name     = "new-compute-subnet-2"
	newComputeSubnet2ID       = "new-compute-2-id"
	newComputeSubnet2CIDR     = "new-compute-2-cidr"
	newComputeSubnet2CRN      = "new-compute-2-crn"
	newComputeSubnet2VPCName  = "new-compute-2-vpc"
	newComputeSubnet2ZoneName = "new-compute-2-zone"
	newComputeSubnets         = map[string]Subnet{
		newComputeSubnet1ID: {
			Name: newComputeSubnet1Name,
			ID:   newComputeSubnet1ID,
			CIDR: newComputeSubnet1CIDR,
			CRN:  newComputeSubnet1CRN,
			VPC:  newComputeSubnet1VPCName,
			Zone: newComputeSubnet1ZoneName,
		},
		newComputeSubnet2ID: {
			Name: newComputeSubnet2Name,
			ID:   newComputeSubnet2ID,
			CIDR: newComputeSubnet2CIDR,
			CRN:  newComputeSubnet2CRN,
			VPC:  newComputeSubnet2VPCName,
			Zone: newComputeSubnet2ZoneName,
		},
	}

	// Previously retrieved Compute Subnets.
	existingComputeSubnets = map[string]Subnet{
		"existing-compute-1-id": {
			Name: "existing-compute-subnet-1",
			ID:   "existing-compute-1-id",
			CIDR: "existing-compute-1-cidr",
			CRN:  "existing-compute-1-crn",
			VPC:  "existing-compute-1-vpc",
			Zone: "existing-compute-1-zone",
		},
		"existing-compute-2-id": {
			Name: "existing-compute-subnet-2",
			ID:   "existing-compute-2-id",
			CIDR: "existing-compute-2-cidr",
			CRN:  "existing-compute-2-crn",
			VPC:  "existing-compute-2-vpc",
			Zone: "existing-compute-2-zone",
		},
	}

	// Newly retrieved Control Plane Subnets.
	newControlPlaneSubnet1Name     = "new-cp-subnet-1"
	newControlPlaneSubnet1ID       = "new-cp-1-id"
	newControlPlaneSubnet1CIDR     = "new-cp-1-cidr"
	newControlPlaneSubnet1CRN      = "new-cp-1-crn"
	newControlPlaneSubnet1VPCName  = "new-cp-1-vpc"
	newControlPlaneSubnet1ZoneName = "new-cp-1-zone"
	newControlPlaneSubnet2Name     = "new-cp-subnet-2"
	newControlPlaneSubnet2ID       = "new-cp-2-id"
	newControlPlaneSubnet2CIDR     = "new-cp-2-cidr"
	newControlPlaneSubnet2CRN      = "new-cp-2-crn"
	newControlPlaneSubnet2VPCName  = "new-cp-2-vpc"
	newControlPlaneSubnet2ZoneName = "new-cp-2-zone"
	newControlPlaneSubnets         = map[string]Subnet{
		newControlPlaneSubnet1ID: {
			Name: newControlPlaneSubnet1Name,
			ID:   newControlPlaneSubnet1ID,
			CIDR: newControlPlaneSubnet1CIDR,
			CRN:  newControlPlaneSubnet1CRN,
			VPC:  newControlPlaneSubnet1VPCName,
			Zone: newControlPlaneSubnet1ZoneName,
		},
		newControlPlaneSubnet2ID: {
			Name: newControlPlaneSubnet2Name,
			ID:   newControlPlaneSubnet2ID,
			CIDR: newControlPlaneSubnet2CIDR,
			CRN:  newControlPlaneSubnet2CRN,
			VPC:  newControlPlaneSubnet2VPCName,
			Zone: newControlPlaneSubnet2ZoneName,
		},
	}

	// Previously retrieved Control Plane Subnets.
	existingControlPlaneSubnets = map[string]Subnet{
		"existing-cp-1-id": {
			Name: "existing-cp-subnet-1",
			ID:   "existing-cp-1-id",
			CIDR: "existing-cp-1-cidr",
			CRN:  "existing-cp-1-crn",
			VPC:  "existing-cp-1-vpc",
			Zone: "existing-cp-1-zone",
		},
		"existing-cp-2-id": {
			Name: "existing-cp-subnet-2",
			ID:   "existing-cp-2-id",
			CIDR: "existing-cp-2-cidr",
			CRN:  "existing-cp-2-crn",
			VPC:  "existing-cp-2-vpc",
			Zone: "existing-cp-2-zone",
		},
	}

	// Subnet Names for failure finding subnets.
	failedGetComputeSubnetName      = "failed-get-compute-subnet"
	failedGetControlPlaneSubnetName = "failed-get-cp-subnet"

	// Subnet Names and ID's for missing values.
	noIDComputeSubnetName       = "no-id-compute-subnet"
	incompleteComputeSubnetName = "incomplete-compute-subnet-name"
	noCIDRComputeSubnetID       = "no-cidr-compute-subnet-id"
	noCRNComputeSubnetID        = "no-crn-compute-subnet-id"
	noNameComputeSubnetID       = "no-name-compute-subnet-id"
	noVPCComputeSubnetID        = "no-vpc-compute-subnet-id"
	noZoneComputeSubnetID       = "no-zone-compute-subnet-id"

	// VPCReferences for Client Subnet responses.
	vpcReferenceComputeSubnet1      = vpcv1.VPCReference{Name: &newComputeSubnet1VPCName}
	vpcReferenceComputeSubnet2      = vpcv1.VPCReference{Name: &newComputeSubnet2VPCName}
	vpcReferenceControlPlaneSubnet1 = vpcv1.VPCReference{Name: &newControlPlaneSubnet1VPCName}
	vpcReferenceControlPlaneSubnet2 = vpcv1.VPCReference{Name: &newControlPlaneSubnet2VPCName}

	// ZoneRefences for Client Subnet responses.
	zoneReferenceComputeSubnet1      = vpcv1.ZoneReference{Name: &newComputeSubnet1ZoneName}
	zoneReferenceComputeSubnet2      = vpcv1.ZoneReference{Name: &newComputeSubnet2ZoneName}
	zoneReferenceControlPlaneSubnet1 = vpcv1.ZoneReference{Name: &newControlPlaneSubnet1ZoneName}
	zoneReferenceControlPlaneSubnet2 = vpcv1.ZoneReference{Name: &newControlPlaneSubnet2ZoneName}
)

func baseMetadata() *Metadata {
	return NewMetadata(&types.InstallConfig{
		BaseDomain: goodDomain,
		Platform: types.Platform{
			IBMCloud: &ibmcloudtypes.Platform{
				Region: region,
			},
		},
		Publish: types.ExternalPublishingStrategy,
	})
}

func setInternalPublishingStrategy(m *Metadata) {
	m.publishStrategy = types.InternalPublishingStrategy
}

func TestAccountID(t *testing.T) {
	testCases := []struct {
		name          string
		edits         editMetadata
		errorMsg      string
		expectedValue string
	}{
		{
			name:          "new accountID",
			expectedValue: newAccountID,
		},
		{
			name: "existing accountID",
			edits: editMetadata{
				func(m *Metadata) {
					m.accountID = existingAccountID
				},
			},
			expectedValue: existingAccountID,
		},
		{
			name:     "auth apikey error",
			errorMsg: "bad api key",
		},
	}

	// IBM Cloud Client Mocks.
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	ibmcloudClient := mock.NewMockAPI(mockCtrl)

	// Shared Mocks.
	ibmcloudClient.EXPECT().SetVPCServiceURLForRegion(gomock.Any(), "us-south").AnyTimes()

	// Mocks: new accountID.
	ibmcloudClient.EXPECT().GetAuthenticatorAPIKeyDetails(gomock.Any()).Return(&iamidentityv1.APIKey{AccountID: &newAccountID}, nil)

	// Mocks: existing accountID.
	// N/A.

	// Mocks: auth apikey error.
	ibmcloudClient.EXPECT().GetAuthenticatorAPIKeyDetails(gomock.Any()).Return(nil, fmt.Errorf("bad api key"))

	for _, tCase := range testCases {
		t.Run(tCase.name, func(t *testing.T) {
			metadata := baseMetadata()
			metadata.client = ibmcloudClient
			for _, edit := range tCase.edits {
				edit(metadata)
			}

			actualValue, err := metadata.AccountID(context.TODO())
			if err != nil {
				assert.Regexp(t, tCase.errorMsg, err)
			} else {
				assert.Equal(t, tCase.expectedValue, actualValue)
			}
		})
	}
}

func TestCreateDNSRecord(t *testing.T) {
	testCases := []struct {
		name         string
		edits        editMetadata
		errorMsg     string
		loadBalancer *vpcv1.LoadBalancer
	}{
		{
			name:     "test zone id failure",
			errorMsg: fmt.Sprintf("failed to retrieve dns zone by base domain %s for %s cluster: id not found", goodDomain, types.ExternalPublishingStrategy),
			loadBalancer: &vpcv1.LoadBalancer{
				Name: ptr.To(fmt.Sprintf("cluster-%s", KubernetesAPIPublicSuffix)),
			},
		},
		{
			name:     "test cis instance crn failure",
			errorMsg: "failed to retrieve cis instance crn for dns record: cis crn not found",
			loadBalancer: &vpcv1.LoadBalancer{
				Name: ptr.To(fmt.Sprintf("cluster-%s", KubernetesAPIPublicSuffix)),
			},
		},
		{
			name:     "test dns instance id failure",
			errorMsg: "failed to retrieve dns instance for dns record: dns services instance not found",
			edits: editMetadata{
				func(m *Metadata) {
					m.publishStrategy = types.InternalPublishingStrategy
				},
			},
			loadBalancer: &vpcv1.LoadBalancer{
				Name: ptr.To(fmt.Sprintf("cluster-%s", KubernetesAPIPublicSuffix)),
			},
		},
		{
			name: "test kube public api",
			loadBalancer: &vpcv1.LoadBalancer{
				Name:     ptr.To(fmt.Sprintf("cluster-%s", KubernetesAPIPublicSuffix)),
				Hostname: ptr.To("lb-hostname"),
			},
		},
		{
			name: "test kube private api",
			loadBalancer: &vpcv1.LoadBalancer{
				Name:     ptr.To(fmt.Sprintf("cluster-%s", KubernetesAPIPrivateSuffix)),
				Hostname: ptr.To("lb-hostname"),
			},
		},
		{
			name: "test dns services dns record",
			edits: editMetadata{
				func(m *Metadata) {
					m.publishStrategy = types.InternalPublishingStrategy
				},
			},
			loadBalancer: &vpcv1.LoadBalancer{
				Name:     ptr.To(fmt.Sprintf("cluster-%s", KubernetesAPIPublicSuffix)),
				Hostname: ptr.To("lb-hostname"),
			},
		},
	}

	clusterDomain := "cluster.test-domain.test"

	// IBM Cloud Client Mocks.
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	ibmcloudClient := mock.NewMockAPI(mockCtrl)

	// Mocks: test zone id failure
	ibmcloudClient.EXPECT().GetDNSZoneIDByName(gomock.Any(), goodDomain, types.ExternalPublishingStrategy).Return("", fmt.Errorf("id not found"))

	// Mocks: test cis instance crn failure
	ibmcloudClient.EXPECT().GetDNSZoneIDByName(gomock.Any(), goodDomain, types.ExternalPublishingStrategy).Return("zoneID", nil)
	ibmcloudClient.EXPECT().GetDNSZones(gomock.Any(), types.ExternalPublishingStrategy).Return(nil, fmt.Errorf("cis crn not found"))

	// Mocks: test dns instance id failure
	ibmcloudClient.EXPECT().GetDNSZoneIDByName(gomock.Any(), goodDomain, types.InternalPublishingStrategy).Return("zoneID", nil)
	ibmcloudClient.EXPECT().GetDNSZones(gomock.Any(), types.InternalPublishingStrategy).Return(nil, fmt.Errorf("dns services instance not found"))

	// Mocks: test kube public api
	ibmcloudClient.EXPECT().GetDNSZoneIDByName(gomock.Any(), goodDomain, types.ExternalPublishingStrategy).Return("zoneID", nil)
	ibmcloudClient.EXPECT().GetDNSZones(gomock.Any(), types.ExternalPublishingStrategy).Return([]responses.DNSZoneResponse{{Name: goodDomain, InstanceCRN: newCISCRN}}, nil)
	ibmcloudClient.EXPECT().CreateCISDNSRecord(gomock.Any(), newCISCRN, "zoneID", fmt.Sprintf("api.%s", clusterDomain), "lb-hostname").Return(nil)

	// Mocks: test kube private api
	ibmcloudClient.EXPECT().GetDNSZoneIDByName(gomock.Any(), goodDomain, types.ExternalPublishingStrategy).Return("zoneID", nil)
	ibmcloudClient.EXPECT().GetDNSZones(gomock.Any(), types.ExternalPublishingStrategy).Return([]responses.DNSZoneResponse{{Name: goodDomain, InstanceCRN: newCISCRN}}, nil)
	ibmcloudClient.EXPECT().CreateCISDNSRecord(gomock.Any(), newCISCRN, "zoneID", fmt.Sprintf("api-int.%s", clusterDomain), "lb-hostname").Return(nil)

	// Mocks: test dns services dns record
	ibmcloudClient.EXPECT().GetDNSZoneIDByName(gomock.Any(), goodDomain, types.InternalPublishingStrategy).Return("zoneID", nil)
	ibmcloudClient.EXPECT().GetDNSZones(gomock.Any(), types.InternalPublishingStrategy).Return([]responses.DNSZoneResponse{{Name: goodDomain, InstanceID: newDNSInstanceID, InstanceCRN: newDNSInstanceCRN}}, nil)
	ibmcloudClient.EXPECT().CreateDNSServicesDNSRecord(gomock.Any(), newDNSInstanceID, "zoneID", fmt.Sprintf("api.%s", clusterDomain), "lb-hostname")

	for _, tCase := range testCases {
		t.Run(tCase.name, func(t *testing.T) {
			metadata := baseMetadata()
			metadata.client = ibmcloudClient
			for _, edit := range tCase.edits {
				edit(metadata)
			}

			err := metadata.CreateDNSRecord(context.TODO(), clusterDomain, tCase.loadBalancer)
			if err != nil {
				assert.Regexp(t, tCase.errorMsg, err)
			}
		})
	}
}

func TestCISInstanceCRN(t *testing.T) {
	testCases := []struct {
		name          string
		edits         editMetadata
		errorMsg      string
		expectedValue string
	}{
		{
			name:          "new cis crn",
			expectedValue: newCISCRN,
		},
		{
			name: "existing cis crn",
			edits: editMetadata{
				func(m *Metadata) {
					m.cisInstanceCRN = existingCISCRN
				},
			},
			expectedValue: existingCISCRN,
		},
		{
			name:     "get dns zone error",
			errorMsg: "dns zone error",
		},
		{
			name:     "dns zone not found error",
			errorMsg: fmt.Sprintf("cisInstanceCRN unknown due to DNS zone %q not found", goodDomain),
		},
	}

	// IBM Cloud Client Mocks.
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	ibmcloudClient := mock.NewMockAPI(mockCtrl)

	// Shared Mocks.
	// N/A.

	// Mocks: new cis crn.
	ibmcloudClient.EXPECT().GetDNSZones(gomock.Any(), types.ExternalPublishingStrategy).Return([]responses.DNSZoneResponse{{Name: goodDomain, InstanceCRN: newCISCRN}}, nil)

	// Mocks: existing cis crn.
	// N/A.

	// Mocks: get dns zone error.
	ibmcloudClient.EXPECT().GetDNSZones(gomock.Any(), types.ExternalPublishingStrategy).Return(nil, fmt.Errorf("dns zone error"))

	// Mocks: dns zone not found error.
	ibmcloudClient.EXPECT().GetDNSZones(gomock.Any(), types.ExternalPublishingStrategy).Return([]responses.DNSZoneResponse{{Name: badDomain}}, nil)

	for _, tCase := range testCases {
		t.Run(tCase.name, func(t *testing.T) {
			metadata := baseMetadata()
			metadata.client = ibmcloudClient
			for _, edit := range tCase.edits {
				edit(metadata)
			}

			actualValue, err := metadata.CISInstanceCRN(context.TODO())
			if err != nil {
				assert.Regexp(t, tCase.errorMsg, err)
			} else {
				assert.Equal(t, tCase.expectedValue, actualValue)
			}
		})
	}
}

func TestSetCISInstanceCRN(t *testing.T) {
	testCases := []struct {
		name   string
		cisCRN string
	}{
		{
			name:   "set cis crn",
			cisCRN: "cis-instance-crn",
		},
	}

	for _, tCase := range testCases {
		t.Run(tCase.name, func(t *testing.T) {
			metadata := baseMetadata()

			metadata.cisInstanceCRN = tCase.cisCRN
			actualValue, err := metadata.CISInstanceCRN(context.TODO())
			assert.NoError(t, err)
			assert.Equal(t, tCase.cisCRN, actualValue)
		})
	}
}

func TestDNSInstance(t *testing.T) {
	testCases := []struct {
		name        string
		edits       editMetadata
		errorMsg    string
		expectedID  string
		expectedCRN string
	}{
		{
			name:        "new dns instance",
			expectedID:  newDNSInstanceID,
			expectedCRN: newDNSInstanceCRN,
		},
		{
			name: "existing dns instance",
			edits: editMetadata{
				func(m *Metadata) {
					m.dnsInstance = &DNSInstance{
						ID:  existingDNSInstanceID,
						CRN: existingDNSInstanceCRN,
					}
				},
			},
			expectedID:  existingDNSInstanceID,
			expectedCRN: existingDNSInstanceCRN,
		},
		{
			name:     "get dns zone error",
			errorMsg: "dns zone error",
		},
		{
			name:     "dns instance unknown id error",
			errorMsg: fmt.Sprintf("dnsInstance has unknown ID/CRN: %q - %q", "", unknownDNSInstanceCRN),
		},
		{
			name:     "dns instance unknown crn error",
			errorMsg: fmt.Sprintf("dnsInstance has unknown ID/CRN: %q - %q", unknownDNSInstanceID, ""),
		},
		{
			name:     "dns zone not found error",
			errorMsg: fmt.Sprintf("dnsInstance unknown due to DNS zone %q not found", goodDomain),
		},
	}

	// IBM Cloud Client Mocks.
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	ibmcloudClient := mock.NewMockAPI(mockCtrl)

	// Shared Mocks.
	// N/A.

	// Mocks: new dns instance.
	ibmcloudClient.EXPECT().GetDNSZones(gomock.Any(), types.InternalPublishingStrategy).Return([]responses.DNSZoneResponse{{Name: goodDomain, InstanceID: newDNSInstanceID, InstanceCRN: newDNSInstanceCRN}}, nil)

	// Mocks: existing dns instance.
	// N/A.

	// Mocks: get dns zone error.
	ibmcloudClient.EXPECT().GetDNSZones(gomock.Any(), types.InternalPublishingStrategy).Return(nil, fmt.Errorf("dns zone error"))

	// Mocks: dns instance unknown id error.
	ibmcloudClient.EXPECT().GetDNSZones(gomock.Any(), types.InternalPublishingStrategy).Return([]responses.DNSZoneResponse{{Name: goodDomain, InstanceCRN: unknownDNSInstanceCRN}}, nil)

	// Mocks: dns instance unknown crn error.
	ibmcloudClient.EXPECT().GetDNSZones(gomock.Any(), types.InternalPublishingStrategy).Return([]responses.DNSZoneResponse{{Name: goodDomain, InstanceID: unknownDNSInstanceID}}, nil)

	// Mocks: dns zone not found error.
	ibmcloudClient.EXPECT().GetDNSZones(gomock.Any(), types.InternalPublishingStrategy).Return([]responses.DNSZoneResponse{{Name: badDomain}}, nil)

	for _, tCase := range testCases {
		t.Run(tCase.name, func(t *testing.T) {
			metadata := baseMetadata()
			setInternalPublishingStrategy(metadata)
			metadata.client = ibmcloudClient
			for _, edit := range tCase.edits {
				edit(metadata)
			}

			actualDNS, err := metadata.DNSInstance(context.TODO())
			if err != nil {
				assert.Regexp(t, tCase.errorMsg, err)
			} else {
				assert.Equal(t, tCase.expectedID, actualDNS.ID)
				assert.Equal(t, tCase.expectedCRN, actualDNS.CRN)
			}
		})
	}
}

func TestSetDNSInstance(t *testing.T) {
	testCases := []struct {
		name   string
		dnsID  string
		dnsCRN string
	}{
		{
			name:   "set dns id/crn",
			dnsID:  "dns-instance-id",
			dnsCRN: "dns-instance-crn",
		},
	}

	for _, tCase := range testCases {
		t.Run(tCase.name, func(t *testing.T) {
			metadata := baseMetadata()
			setInternalPublishingStrategy(metadata)

			metadata.dnsInstance = &DNSInstance{
				ID:  tCase.dnsID,
				CRN: tCase.dnsCRN,
			}
			actualDNS, err := metadata.DNSInstance(context.TODO())
			assert.NoError(t, err)
			assert.Equal(t, tCase.dnsID, actualDNS.ID)
			assert.Equal(t, tCase.dnsCRN, actualDNS.CRN)
		})
	}
}

func TestComputeSubnets(t *testing.T) {
	testCases := []struct {
		name          string
		edits         editMetadata
		errorMsg      string
		expectedValue map[string]Subnet
	}{
		{
			name:          "no compute subnets",
			expectedValue: nil,
		},
		{
			name: "new compute subnets",
			edits: editMetadata{
				func(m *Metadata) {
					m.ComputeSubnetNames = []string{newComputeSubnet1Name, newComputeSubnet2Name}
				},
			},
			expectedValue: newComputeSubnets,
		},
		{
			name: "new single compute subnet",
			edits: editMetadata{
				func(m *Metadata) {
					m.ComputeSubnetNames = []string{newComputeSubnet2Name}
				},
			},
			expectedValue: map[string]Subnet{
				newComputeSubnet2ID: newComputeSubnets[newComputeSubnet2ID],
			},
		},
		{
			name: "existing compute subnets",
			edits: editMetadata{
				func(m *Metadata) {
					m.computeSubnets = existingComputeSubnets
				},
			},
			expectedValue: existingComputeSubnets,
		},
		{
			name: "failed getting compute subnet",
			edits: editMetadata{
				func(m *Metadata) {
					m.ComputeSubnetNames = []string{failedGetComputeSubnetName}
				},
			},
			errorMsg: fmt.Sprintf("getting subnet %s", failedGetComputeSubnetName),
		},
		{
			name: "compute subnet has no id",
			edits: editMetadata{
				func(m *Metadata) {
					m.ComputeSubnetNames = []string{noIDComputeSubnetName}
				},
			},
			errorMsg: fmt.Sprintf("%s has no ID", noIDComputeSubnetName),
		},
		{
			name: "compute subnet has no CIDR block",
			edits: editMetadata{
				func(m *Metadata) {
					m.ComputeSubnetNames = []string{incompleteComputeSubnetName}
				},
			},
			errorMsg: fmt.Sprintf("%s has no Ipv4CIDRBlock", noCIDRComputeSubnetID),
		},
		{
			name: "compute subnet has no CRN",
			edits: editMetadata{
				func(m *Metadata) {
					m.ComputeSubnetNames = []string{incompleteComputeSubnetName}
				},
			},
			errorMsg: fmt.Sprintf("%s has no CRN", noCRNComputeSubnetID),
		},
		{
			name: "compute subnet has no Name",
			edits: editMetadata{
				func(m *Metadata) {
					m.ComputeSubnetNames = []string{incompleteComputeSubnetName}
				},
			},
			errorMsg: fmt.Sprintf("%s has no Name", noNameComputeSubnetID),
		},
		{
			name: "compute subnet has no VPC",
			edits: editMetadata{
				func(m *Metadata) {
					m.ComputeSubnetNames = []string{incompleteComputeSubnetName}
				},
			},
			errorMsg: fmt.Sprintf("%s has no VPC", noVPCComputeSubnetID),
		},
		{
			name: "compute subnet has no Zone",
			edits: editMetadata{
				func(m *Metadata) {
					m.ComputeSubnetNames = []string{incompleteComputeSubnetName}
				},
			},
			errorMsg: fmt.Sprintf("%s has no Zone", noZoneComputeSubnetID),
		},
	}

	// IBM Cloud Client Mocks.
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	ibmcloudClient := mock.NewMockAPI(mockCtrl)

	// Shared Mocks.
	// N/A.

	// Mocks: no compute subnets.
	// N/A.

	// Mocks: new compute subnets.
	ibmcloudClient.EXPECT().GetSubnetByName(gomock.Any(), newComputeSubnet1Name, region).Return(
		&vpcv1.Subnet{
			Name:          &newComputeSubnet1Name,
			ID:            &newComputeSubnet1ID,
			Ipv4CIDRBlock: &newComputeSubnet1CIDR,
			CRN:           &newComputeSubnet1CRN,
			VPC:           &vpcReferenceComputeSubnet1,
			Zone:          &zoneReferenceComputeSubnet1,
		},
		nil,
	)
	ibmcloudClient.EXPECT().GetSubnetByName(gomock.Any(), newComputeSubnet2Name, region).Return(
		&vpcv1.Subnet{
			Name:          &newComputeSubnet2Name,
			ID:            &newComputeSubnet2ID,
			Ipv4CIDRBlock: &newComputeSubnet2CIDR,
			CRN:           &newComputeSubnet2CRN,
			VPC:           &vpcReferenceComputeSubnet2,
			Zone:          &zoneReferenceComputeSubnet2,
		},
		nil,
	)

	// Mocks: new single compute subnet.
	ibmcloudClient.EXPECT().GetSubnetByName(gomock.Any(), newComputeSubnet2Name, region).Return(
		&vpcv1.Subnet{
			Name:          &newComputeSubnet2Name,
			ID:            &newComputeSubnet2ID,
			Ipv4CIDRBlock: &newComputeSubnet2CIDR,
			CRN:           &newComputeSubnet2CRN,
			VPC:           &vpcReferenceComputeSubnet2,
			Zone:          &zoneReferenceComputeSubnet2,
		},
		nil,
	)

	// Mocks: existing compute subnets.
	// N/A.

	// Mocks: failed getting compute subnet.
	ibmcloudClient.EXPECT().GetSubnetByName(gomock.Any(), failedGetComputeSubnetName, region).Return(nil, &VPCResourceNotFoundError{})

	// Mocks: compute subnet has no ID.
	ibmcloudClient.EXPECT().GetSubnetByName(gomock.Any(), noIDComputeSubnetName, region).Return(
		&vpcv1.Subnet{
			Name: &noIDComputeSubnetName,
			// ID Skipped.
			Ipv4CIDRBlock: &newComputeSubnet1CIDR,
			CRN:           &newComputeSubnet1CRN,
			VPC:           &vpcReferenceComputeSubnet1,
			Zone:          &zoneReferenceComputeSubnet1,
		},
		nil,
	)

	// Mocks: compute subnet has no CIDR.
	ibmcloudClient.EXPECT().GetSubnetByName(gomock.Any(), incompleteComputeSubnetName, region).Return(
		&vpcv1.Subnet{
			Name: &incompleteComputeSubnetName,
			ID:   &noCIDRComputeSubnetID,
			// Ipv4CIDRBlock Skipped.
			CRN:  &newComputeSubnet1CRN,
			VPC:  &vpcReferenceComputeSubnet1,
			Zone: &zoneReferenceComputeSubnet1,
		},
		nil,
	)

	// Mocks: compute subnet has no CRN.
	ibmcloudClient.EXPECT().GetSubnetByName(gomock.Any(), incompleteComputeSubnetName, region).Return(
		&vpcv1.Subnet{
			Name:          &incompleteComputeSubnetName,
			ID:            &noCRNComputeSubnetID,
			Ipv4CIDRBlock: &newComputeSubnet1CIDR,
			// CRN Skipped.
			VPC:  &vpcReferenceComputeSubnet1,
			Zone: &zoneReferenceComputeSubnet1,
		},
		nil,
	)

	// Mocks: compute subnet has no Name.
	ibmcloudClient.EXPECT().GetSubnetByName(gomock.Any(), incompleteComputeSubnetName, region).Return(
		&vpcv1.Subnet{
			// Name Skipped.
			ID:            &noNameComputeSubnetID,
			Ipv4CIDRBlock: &newComputeSubnet1CIDR,
			CRN:           &newComputeSubnet1CRN,
			VPC:           &vpcReferenceComputeSubnet1,
			Zone:          &zoneReferenceComputeSubnet1,
		},
		nil,
	)

	// Mocks: compute subnet has no VPC.
	ibmcloudClient.EXPECT().GetSubnetByName(gomock.Any(), incompleteComputeSubnetName, region).Return(
		&vpcv1.Subnet{
			Name:          &incompleteComputeSubnetName,
			ID:            &noVPCComputeSubnetID,
			Ipv4CIDRBlock: &newComputeSubnet1CIDR,
			CRN:           &newComputeSubnet1CRN,
			// VPC Skipped.
			Zone: &zoneReferenceComputeSubnet1,
		},
		nil,
	)

	// Mocks: compute subnet has no Zone.
	ibmcloudClient.EXPECT().GetSubnetByName(gomock.Any(), incompleteComputeSubnetName, region).Return(
		&vpcv1.Subnet{
			Name:          &incompleteComputeSubnetName,
			ID:            &noZoneComputeSubnetID,
			Ipv4CIDRBlock: &newComputeSubnet1CIDR,
			CRN:           &newComputeSubnet1CRN,
			VPC:           &vpcReferenceComputeSubnet1,
			// Zone Skipped.
		},
		nil,
	)

	for _, tCase := range testCases {
		t.Run(tCase.name, func(t *testing.T) {
			metadata := baseMetadata()
			metadata.client = ibmcloudClient

			for _, edit := range tCase.edits {
				edit(metadata)
			}

			actualValue, err := metadata.ComputeSubnets(context.TODO())
			if err != nil {
				assert.Regexp(t, tCase.errorMsg, err)
			} else {
				assert.Equal(t, tCase.expectedValue, actualValue)
			}
		})
	}
}

func TestControlPlaneSubnets(t *testing.T) {
	testCases := []struct {
		name          string
		edits         editMetadata
		errorMsg      string
		expectedValue map[string]Subnet
	}{
		{
			name:          "no control plane subnets",
			expectedValue: nil,
		},
		{
			name: "new control plane subnets",
			edits: editMetadata{
				func(m *Metadata) {
					m.ControlPlaneSubnetNames = []string{newControlPlaneSubnet1Name, newControlPlaneSubnet2Name}
				},
			},
			expectedValue: newControlPlaneSubnets,
		},
		{
			name: "new single control plane subnet",
			edits: editMetadata{
				func(m *Metadata) {
					m.ControlPlaneSubnetNames = []string{newControlPlaneSubnet2Name}
				},
			},
			expectedValue: map[string]Subnet{
				newControlPlaneSubnet2ID: newControlPlaneSubnets[newControlPlaneSubnet2ID],
			},
		},
		{
			name: "existing control plane subnets",
			edits: editMetadata{
				func(m *Metadata) {
					m.controlPlaneSubnets = existingControlPlaneSubnets
				},
			},
			expectedValue: existingControlPlaneSubnets,
		},
		{
			name: "failed getting control plane subnet",
			edits: editMetadata{
				func(m *Metadata) {
					m.ControlPlaneSubnetNames = []string{failedGetControlPlaneSubnetName}
				},
			},
			errorMsg: fmt.Sprintf("getting subnet %v", failedGetControlPlaneSubnetName),
		},
	}

	// IBM Cloud Client Mocks.
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	ibmcloudClient := mock.NewMockAPI(mockCtrl)

	// Shared Mocks.
	// N/A.

	// Mocks: no control plane subnets.
	// N/A.

	// Mocks: new control plane subnets.
	ibmcloudClient.EXPECT().GetSubnetByName(gomock.Any(), newControlPlaneSubnet1Name, region).Return(
		&vpcv1.Subnet{
			Name:          &newControlPlaneSubnet1Name,
			ID:            &newControlPlaneSubnet1ID,
			Ipv4CIDRBlock: &newControlPlaneSubnet1CIDR,
			CRN:           &newControlPlaneSubnet1CRN,
			VPC:           &vpcReferenceControlPlaneSubnet1,
			Zone:          &zoneReferenceControlPlaneSubnet1,
		},
		nil,
	)
	ibmcloudClient.EXPECT().GetSubnetByName(gomock.Any(), newControlPlaneSubnet2Name, region).Return(
		&vpcv1.Subnet{
			Name:          &newControlPlaneSubnet2Name,
			ID:            &newControlPlaneSubnet2ID,
			Ipv4CIDRBlock: &newControlPlaneSubnet2CIDR,
			CRN:           &newControlPlaneSubnet2CRN,
			VPC:           &vpcReferenceControlPlaneSubnet2,
			Zone:          &zoneReferenceControlPlaneSubnet2,
		},
		nil,
	)

	// Mocks: new single control plane subnet.
	ibmcloudClient.EXPECT().GetSubnetByName(gomock.Any(), newControlPlaneSubnet2Name, region).Return(
		&vpcv1.Subnet{
			Name:          &newControlPlaneSubnet2Name,
			ID:            &newControlPlaneSubnet2ID,
			Ipv4CIDRBlock: &newControlPlaneSubnet2CIDR,
			CRN:           &newControlPlaneSubnet2CRN,
			VPC:           &vpcReferenceControlPlaneSubnet2,
			Zone:          &zoneReferenceControlPlaneSubnet2,
		},
		nil,
	)

	// Mocks: existing control plane subnets.
	// N/A.

	// Mocks: failed getting control plane subnet.
	ibmcloudClient.EXPECT().GetSubnetByName(gomock.Any(), failedGetControlPlaneSubnetName, region).Return(nil, &VPCResourceNotFoundError{})

	for _, tCase := range testCases {
		t.Run(tCase.name, func(t *testing.T) {
			metadata := baseMetadata()
			metadata.client = ibmcloudClient

			for _, edit := range tCase.edits {
				edit(metadata)
			}

			actualValue, err := metadata.ControlPlaneSubnets(context.TODO())
			if err != nil {
				assert.Regexp(t, tCase.errorMsg, err)
			} else {
				assert.Equal(t, tCase.expectedValue, actualValue)
			}
		})
	}
}

func TestClient(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	ibmcloudClient := mock.NewMockAPI(mockCtrl)

	testCases := []struct {
		name          string
		edits         editMetadata
		errorMsg      string
		expectedValue API
	}{
		{
			name: "existing client",
			edits: editMetadata{
				func(m *Metadata) {
					m.client = ibmcloudClient
				},
			},
			expectedValue: ibmcloudClient,
		},
	}

	// IBM Cloud Client Mocks.
	// N/A.

	for _, tCase := range testCases {
		t.Run(tCase.name, func(t *testing.T) {
			metadata := baseMetadata()

			actualValue, err := metadata.Client()
			if err != nil {
				assert.Regexp(t, tCase.errorMsg, err)
			} else {
				assert.Equal(t, tCase.expectedValue, actualValue)
			}
		})
	}
}
