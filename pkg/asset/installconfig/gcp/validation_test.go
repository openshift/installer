package gcp

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"strings"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	googleoauth "golang.org/x/oauth2/google"
	"google.golang.org/api/cloudresourcemanager/v1"
	compute "google.golang.org/api/compute/v1"
	dns "google.golang.org/api/dns/v1"
	"google.golang.org/api/googleapi"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/sets"
	"k8s.io/apimachinery/pkg/util/validation/field"

	"github.com/openshift/installer/pkg/asset/installconfig/gcp/mock"
	"github.com/openshift/installer/pkg/ipnet"
	"github.com/openshift/installer/pkg/types"
	"github.com/openshift/installer/pkg/types/gcp"
)

type editFunctions []func(ic *types.InstallConfig)

var (
	validNetworkName   = "valid-vpc"
	validProjectName   = "valid-project"
	invalidProjectName = "invalid-project"
	validRegion        = "us-east1"
	invalidRegion      = "us-east4"
	validZone          = "us-east1-b"
	validComputeSubnet = "valid-compute-subnet"
	validCPSubnet      = "valid-controlplane-subnet"
	validCIDR          = "10.0.0.0/16"
	validClusterName   = "valid-cluster"
	validPrivateZone   = "valid-short-private-zone"
	validPublicZone    = "valid-short-public-zone"
	invalidPublicZone  = "invalid-short-public-zone"
	validBaseDomain    = "example.installer.domain."
	validXpnSA         = "valid-example-sa@gcloud.serviceaccount.com"
	invalidXpnSA       = "invalid-example-sa@gcloud.serviceaccount.com"

	validPrivateDNSZone = dns.ManagedZone{
		Name:    validPrivateZone,
		DnsName: fmt.Sprintf("%s.%s", validClusterName, strings.TrimSuffix(validBaseDomain, ".")),
	}
	validPublicDNSZone = dns.ManagedZone{
		Name:    validPublicZone,
		DnsName: validBaseDomain,
	}

	invalidateMachineCIDR = func(ic *types.InstallConfig) {
		_, newCidr, _ := net.ParseCIDR("192.168.111.0/24")
		ic.MachineNetwork = []types.MachineNetworkEntry{
			{CIDR: ipnet.IPNet{IPNet: *newCidr}},
		}
	}

	validMachineTypes = func(ic *types.InstallConfig) {
		ic.Platform.GCP.DefaultMachinePlatform.InstanceType = "n1-standard-2"
		ic.ControlPlane.Platform.GCP.InstanceType = "n1-standard-4"
		ic.Compute[0].Platform.GCP.InstanceType = "n1-standard-2"
	}

	invalidateDefaultMachineTypes = func(ic *types.InstallConfig) {
		ic.Platform.GCP.DefaultMachinePlatform.InstanceType = "n1-standard-1"
	}

	invalidateControlPlaneMachineTypes = func(ic *types.InstallConfig) {
		ic.ControlPlane.Platform.GCP.InstanceType = "n1-standard-1"
	}

	invalidateComputeMachineTypes = func(ic *types.InstallConfig) {
		ic.Compute[0].Platform.GCP.InstanceType = "n1-standard-1"
	}

	undefinedDefaultMachineTypes = func(ic *types.InstallConfig) {
		ic.Platform.GCP.DefaultMachinePlatform.InstanceType = "n1-dne-1"
	}

	invalidateNetwork        = func(ic *types.InstallConfig) { ic.GCP.Network = "invalid-vpc" }
	invalidateComputeSubnet  = func(ic *types.InstallConfig) { ic.GCP.ComputeSubnet = "invalid-compute-subnet" }
	invalidateCPSubnet       = func(ic *types.InstallConfig) { ic.GCP.ControlPlaneSubnet = "invalid-cp-subnet" }
	invalidateRegion         = func(ic *types.InstallConfig) { ic.GCP.Region = invalidRegion }
	invalidateProject        = func(ic *types.InstallConfig) { ic.GCP.ProjectID = invalidProjectName }
	invalidateNetworkProject = func(ic *types.InstallConfig) { ic.GCP.NetworkProjectID = invalidProjectName }
	removeVPC                = func(ic *types.InstallConfig) { ic.GCP.Network = "" }
	removeSubnets            = func(ic *types.InstallConfig) { ic.GCP.ComputeSubnet, ic.GCP.ControlPlaneSubnet = "", "" }
	invalidClusterName       = func(ic *types.InstallConfig) { ic.ObjectMeta.Name = "testgoogletest" }
	validNetworkProject      = func(ic *types.InstallConfig) { ic.GCP.NetworkProjectID = validProjectName }
	validateXpnSA            = func(ic *types.InstallConfig) { ic.ControlPlane.Platform.GCP.ServiceAccount = validXpnSA }
	invalidateXpnSA          = func(ic *types.InstallConfig) { ic.ControlPlane.Platform.GCP.ServiceAccount = invalidXpnSA }

	machineTypeAPIResult = map[string]*compute.MachineType{
		"n1-standard-1": {GuestCpus: 1, MemoryMb: 3840},
		"n1-standard-2": {GuestCpus: 2, MemoryMb: 7680},
		"n1-standard-4": {GuestCpus: 4, MemoryMb: 15360},
	}

	subnetAPIResult = []*compute.Subnetwork{
		{
			Name:        validCPSubnet,
			IpCidrRange: validCIDR,
		},
		{
			Name:        validComputeSubnet,
			IpCidrRange: validCIDR,
		},
	}
)

func validInstallConfig() *types.InstallConfig {
	return &types.InstallConfig{
		BaseDomain: validBaseDomain,
		ObjectMeta: metav1.ObjectMeta{
			Name: validClusterName,
		},
		Networking: &types.Networking{
			MachineNetwork: []types.MachineNetworkEntry{
				{CIDR: *ipnet.MustParseCIDR(validCIDR)},
			},
		},
		Platform: types.Platform{
			GCP: &gcp.Platform{
				DefaultMachinePlatform: &gcp.MachinePool{},
				ProjectID:              validProjectName,
				Region:                 validRegion,
				Network:                validNetworkName,
				ComputeSubnet:          validComputeSubnet,
				ControlPlaneSubnet:     validCPSubnet,
			},
		},
		ControlPlane: &types.MachinePool{
			Platform: types.MachinePoolPlatform{
				GCP: &gcp.MachinePool{},
			},
		},
		Compute: []types.MachinePool{{
			Platform: types.MachinePoolPlatform{
				GCP: &gcp.MachinePool{},
			},
		}},
	}
}

func TestGCPInstallConfigValidation(t *testing.T) {
	cases := []struct {
		name           string
		edits          editFunctions
		expectedError  bool
		expectedErrMsg string
	}{
		{
			name:           "Valid network & subnets",
			edits:          editFunctions{},
			expectedError:  false,
			expectedErrMsg: "",
		},
		{
			name:           "Invalid ClusterName",
			edits:          editFunctions{invalidClusterName},
			expectedError:  true,
			expectedErrMsg: `clusterName: Invalid value: "testgoogletest": cluster name must not start with "goog" or contain variations of "google"`,
		},
		{
			name:           "Valid install config without network & subnets",
			edits:          editFunctions{removeVPC, removeSubnets},
			expectedError:  false,
			expectedErrMsg: "",
		},
		{
			name:           "Invalid subnet range",
			edits:          editFunctions{invalidateMachineCIDR},
			expectedError:  true,
			expectedErrMsg: "computeSubnet: Invalid value.*subnet CIDR range start 10.0.0.0 is outside of the specified machine networks",
		},
		{
			name:           "Invalid network",
			edits:          editFunctions{invalidateNetwork},
			expectedError:  true,
			expectedErrMsg: "network: Invalid value",
		},
		{
			name:           "Invalid compute subnet",
			edits:          editFunctions{invalidateComputeSubnet},
			expectedError:  true,
			expectedErrMsg: "computeSubnet: Invalid value",
		},
		{
			name:           "Invalid control plane subnet",
			edits:          editFunctions{invalidateCPSubnet},
			expectedError:  true,
			expectedErrMsg: "controlPlaneSubnet: Invalid value",
		},
		{
			name:           "Invalid both subnets",
			edits:          editFunctions{invalidateCPSubnet, invalidateComputeSubnet},
			expectedError:  true,
			expectedErrMsg: "computeSubnet: Invalid value.*controlPlaneSubnet: Invalid value",
		},
		{
			name:           "Valid machine types",
			edits:          editFunctions{validMachineTypes},
			expectedError:  false,
			expectedErrMsg: "",
		},
		{
			name:           "Invalid default machine type",
			edits:          editFunctions{invalidateDefaultMachineTypes},
			expectedError:  true,
			expectedErrMsg: `\[platform.gcp.defaultMachinePlatform.type: Invalid value: "n1-standard-1": instance type does not meet minimum resource requirements of 4 vCPUs, platform.gcp.defaultMachinePlatform.type: Invalid value: "n1-standard-1": instance type does not meet minimum resource requirements of 15360 MB Memory\]`,
		},
		{
			name:           "Invalid control plane machine types",
			edits:          editFunctions{invalidateControlPlaneMachineTypes},
			expectedError:  true,
			expectedErrMsg: `[controlPlane.platform.gcp.type: Invalid value: "n1\-standard\-1": instance type does not meet minimum resource requirements of 4 vCPUs, controlPlane.platform.gcp.type: Invalid value: "n1\-standard\-1": instance type does not meet minimum resource requirements of 15361 MB Memory]`,
		},
		{
			name:           "Invalid compute machine types",
			edits:          editFunctions{invalidateComputeMachineTypes},
			expectedError:  true,
			expectedErrMsg: `\[compute\[0\].platform.gcp.type: Invalid value: "n1-standard-1": instance type does not meet minimum resource requirements of 2 vCPUs, compute\[0\].platform.gcp.type: Invalid value: "n1-standard-1": instance type does not meet minimum resource requirements of 7680 MB Memory\]`,
		},
		{
			name:           "Undefined default machine types",
			edits:          editFunctions{undefinedDefaultMachineTypes},
			expectedError:  true,
			expectedErrMsg: `Internal error: 404`,
		},
		{
			name:           "Invalid region",
			edits:          editFunctions{invalidateRegion},
			expectedError:  true,
			expectedErrMsg: "could not find subnet valid-compute-subnet in network valid-vpc and region us-east4",
		},
		{
			name:           "Invalid project",
			edits:          editFunctions{invalidateProject},
			expectedError:  true,
			expectedErrMsg: "network: Invalid value",
		},
		{
			name:           "Invalid project & region",
			edits:          editFunctions{invalidateRegion, invalidateProject},
			expectedError:  true,
			expectedErrMsg: "network: Invalid value",
		},
		{
			name:           "Invalid project ID",
			edits:          editFunctions{invalidateProject, removeSubnets, removeVPC},
			expectedError:  true,
			expectedErrMsg: "platform.gcp.project: Invalid value: \"invalid-project\": invalid project ID",
		},
		{
			name:           "Invalid network project ID",
			edits:          editFunctions{invalidateNetworkProject},
			expectedError:  true,
			expectedErrMsg: "platform.gcp.networkProjectID: Invalid value: \"invalid-project\": invalid project ID",
		},
		{
			name:           "Valid Region",
			edits:          editFunctions{},
			expectedError:  false,
			expectedErrMsg: "",
		},
		{
			name:           "Invalid region not found",
			edits:          editFunctions{invalidateRegion, invalidateProject},
			expectedError:  true,
			expectedErrMsg: "platform.gcp.project: Invalid value: \"invalid-project\": invalid project ID",
		},
		{
			name:           "Region not validated",
			edits:          editFunctions{invalidateRegion},
			expectedError:  true,
			expectedErrMsg: "platform.gcp.region: Invalid value: \"us-east4\": invalid region",
		},
		{
			name:          "Valid XPN Service Account",
			edits:         editFunctions{validNetworkProject, validateXpnSA},
			expectedError: false,
		},
		{
			name:           "Invalid XPN Service Account",
			edits:          editFunctions{validNetworkProject, invalidateXpnSA},
			expectedError:  true,
			expectedErrMsg: "controlPlane.platform.gcp.serviceAccount: Internal error\"",
		},
	}
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	gcpClient := mock.NewMockAPI(mockCtrl)

	errNotFound := &googleapi.Error{Code: http.StatusNotFound}

	// Should get the list of projects.
	gcpClient.EXPECT().GetProjects(gomock.Any()).Return(map[string]string{"valid-project": "valid-project"}, nil).AnyTimes()
	gcpClient.EXPECT().GetProjectByID(gomock.Any(), "valid-project").Return(&cloudresourcemanager.Project{}, nil).AnyTimes()
	gcpClient.EXPECT().GetProjectByID(gomock.Any(), "invalid-project").Return(nil, errNotFound).AnyTimes()
	gcpClient.EXPECT().GetProjectByID(gomock.Any(), gomock.Any()).Return(nil, fmt.Errorf("error")).AnyTimes()

	// Should get the list of zones.
	gcpClient.EXPECT().GetZones(gomock.Any(), gomock.Any(), gomock.Any()).Return([]*compute.Zone{{Name: validZone}}, nil).AnyTimes()

	// When passed an invalid project, no regions will be returned
	gcpClient.EXPECT().GetRegions(gomock.Any(), invalidProjectName).Return(nil, fmt.Errorf("failed to get regions for project")).AnyTimes()
	// When passed a project that is valid but the region is not contained, an error should still occur
	gcpClient.EXPECT().GetRegions(gomock.Any(), validProjectName).Return([]string{validRegion}, nil).AnyTimes()

	// Should return the machine type as specified.
	for key, value := range machineTypeAPIResult {
		gcpClient.EXPECT().GetMachineTypeWithZones(gomock.Any(), gomock.Any(), gomock.Any(), key).Return(value, sets.New(validZone), nil).AnyTimes()
	}
	// When passed incorrect machine type, the API returns nil.
	gcpClient.EXPECT().GetMachineTypeWithZones(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(nil, nil, fmt.Errorf("404")).AnyTimes()

	// When passed the correct network & project, return an empty network, which should be enough to validate ok.
	gcpClient.EXPECT().GetNetwork(gomock.Any(), validNetworkName, validProjectName).Return(&compute.Network{}, nil).AnyTimes()

	// When passed an incorrect network or incorrect project, the API returns nil
	gcpClient.EXPECT().GetNetwork(gomock.Any(), gomock.Not(validNetworkName), gomock.Any()).Return(nil, fmt.Errorf("404")).AnyTimes()
	gcpClient.EXPECT().GetNetwork(gomock.Any(), gomock.Any(), gomock.Not(validProjectName)).Return(nil, fmt.Errorf("404")).AnyTimes()

	// When passed a correct network, project, & region, returns valid subnets.
	// We will test incorrect subnets, by changing the install config.
	gcpClient.EXPECT().GetSubnetworks(gomock.Any(), validNetworkName, validProjectName, validRegion).Return(subnetAPIResult, nil).AnyTimes()

	// When passed an incorrect network, project or region, return empty list.
	gcpClient.EXPECT().GetSubnetworks(gomock.Any(), gomock.Not(validNetworkName), gomock.Any(), gomock.Any()).Return([]*compute.Subnetwork{}, nil).AnyTimes()
	gcpClient.EXPECT().GetSubnetworks(gomock.Any(), gomock.Any(), gomock.Not(validProjectName), gomock.Any()).Return([]*compute.Subnetwork{}, nil).AnyTimes()
	gcpClient.EXPECT().GetSubnetworks(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Not(validRegion)).Return([]*compute.Subnetwork{}, nil).AnyTimes()

	// Return fake credentials when asked
	gcpClient.EXPECT().GetCredentials().Return(&googleoauth.Credentials{JSON: []byte("fake creds")}).AnyTimes()

	// Expected results for the managed zone tests
	gcpClient.EXPECT().GetDNSZoneByName(gomock.Any(), gomock.Any(), validPublicZone).Return(&validPublicDNSZone, nil).AnyTimes()
	gcpClient.EXPECT().GetDNSZoneByName(gomock.Any(), gomock.Any(), validPrivateZone).Return(&validPrivateDNSZone, nil).AnyTimes()
	gcpClient.EXPECT().GetDNSZoneByName(gomock.Any(), gomock.Any(), invalidPublicZone).Return(nil, fmt.Errorf("no matching DNS Zone found")).AnyTimes()

	gcpClient.EXPECT().GetServiceAccount(gomock.Any(), validProjectName, validXpnSA).Return(validXpnSA, nil).AnyTimes()
	gcpClient.EXPECT().GetServiceAccount(gomock.Any(), validProjectName, invalidXpnSA).Return("", fmt.Errorf("controlPlane.platform.gcp.serviceAccount: Internal error\"")).AnyTimes()

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			editedInstallConfig := validInstallConfig()
			for _, edit := range tc.edits {
				edit(editedInstallConfig)
			}

			errs := Validate(gcpClient, editedInstallConfig)
			if tc.expectedError {
				assert.Regexp(t, tc.expectedErrMsg, errs)
			} else {
				assert.Empty(t, errs)
			}
		})
	}
}

func TestValidatePreExistingPublicDNS(t *testing.T) {
	cases := []struct {
		name    string
		records []*dns.ResourceRecordSet
		err     string
	}{{
		name:    "no pre-existing",
		records: nil,
	}, {
		name:    "no pre-existing",
		records: []*dns.ResourceRecordSet{{Name: "api.another-cluster-name.base-domain."}},
	}, {
		name:    "pre-existing",
		records: []*dns.ResourceRecordSet{{Name: "api.cluster-name.base-domain."}},
		err:     `^metadata\.name: Invalid value: "cluster-name": record\(s\) \["api\.cluster-name\.base-domain\."\] already exists in DNS Zone \(project-id/zone-name\) and might be in use by another cluster, please remove it to continue$`,
	}, {
		name:    "pre-existing",
		records: []*dns.ResourceRecordSet{{Name: "api.cluster-name.base-domain."}, {Name: "api.cluster-name.base-domain."}},
		err:     `^metadata\.name: Invalid value: "cluster-name": record\(s\) \["api\.cluster-name\.base-domain\."\] already exists in DNS Zone \(project-id/zone-name\) and might be in use by another cluster, please remove it to continue$`,
	}}

	for _, test := range cases {
		t.Run(test.name, func(t *testing.T) {
			mockCtrl := gomock.NewController(t)
			defer mockCtrl.Finish()
			gcpClient := mock.NewMockAPI(mockCtrl)

			gcpClient.EXPECT().GetDNSZone(gomock.Any(), "project-id", "base-domain", true).Return(&dns.ManagedZone{Name: "zone-name"}, nil).AnyTimes()
			gcpClient.EXPECT().GetRecordSets(gomock.Any(), gomock.Eq("project-id"), gomock.Eq("zone-name")).Return(test.records, nil).AnyTimes()

			err := ValidatePreExistingPublicDNS(gcpClient, &types.InstallConfig{
				ObjectMeta: metav1.ObjectMeta{Name: "cluster-name"},
				BaseDomain: "base-domain",
				Platform:   types.Platform{GCP: &gcp.Platform{ProjectID: "project-id"}},
			})
			if test.err == "" {
				assert.True(t, err == nil)
			} else {
				assert.Regexp(t, test.err, err)
			}
		})
	}
}

func TestValidatePrivateDNSZone(t *testing.T) {
	cases := []struct {
		name    string
		records []*dns.ResourceRecordSet
		err     string
	}{{
		name:    "no pre-existing",
		records: nil,
	}, {
		name:    "no pre-existing",
		records: []*dns.ResourceRecordSet{{Name: "api.another-cluster-name.base-domain."}},
	}, {
		name:    "pre-existing",
		records: []*dns.ResourceRecordSet{{Name: "api.cluster-name.base-domain."}},
		err:     `^metadata\.name: Invalid value: "cluster-name": record\(s\) \["api\.cluster-name\.base-domain\."\] already exists in DNS Zone \(project-id/zone-name\) and might be in use by another cluster, please remove it to continue$`,
	}, {
		name:    "pre-existing",
		records: []*dns.ResourceRecordSet{{Name: "api.cluster-name.base-domain."}, {Name: "api.cluster-name.base-domain."}},
		err:     `^metadata\.name: Invalid value: "cluster-name": record\(s\) \["api\.cluster-name\.base-domain\."\] already exists in DNS Zone \(project-id/zone-name\) and might be in use by another cluster, please remove it to continue$`,
	}}

	for _, test := range cases {
		t.Run(test.name, func(t *testing.T) {
			mockCtrl := gomock.NewController(t)
			defer mockCtrl.Finish()
			gcpClient := mock.NewMockAPI(mockCtrl)

			gcpClient.EXPECT().GetDNSZone(gomock.Any(), "project-id", "cluster-name.base-domain", false).Return(&dns.ManagedZone{Name: "zone-name"}, nil).AnyTimes()
			gcpClient.EXPECT().GetRecordSets(gomock.Any(), gomock.Eq("project-id"), gomock.Eq("zone-name")).Return(test.records, nil).AnyTimes()

			err := ValidatePrivateDNSZone(gcpClient, &types.InstallConfig{
				ObjectMeta: metav1.ObjectMeta{Name: "cluster-name"},
				BaseDomain: "base-domain",
				Platform:   types.Platform{GCP: &gcp.Platform{ProjectID: "project-id", Network: "shared-vpc", NetworkProjectID: "test-network-project"}},
			})
			if test.err == "" {
				assert.True(t, err == nil)
			} else {
				assert.Regexp(t, test.err, err)
			}
		})
	}
}

func TestGCPEnabledServicesList(t *testing.T) {
	cases := []struct {
		name     string
		services []string
		err      string
	}{{
		name:     "No services present",
		services: nil,
		err:      "unable to fetch enabled services for project. Make sure 'serviceusage.googleapis.com' is enabled",
	}, {
		name:     "Service Usage missing",
		services: []string{"compute.googleapis.com"},
		err:      "unable to fetch enabled services for project. Make sure 'serviceusage.googleapis.com' is enabled",
	}, {
		name: "All pre-existing",
		services: []string{
			"compute.googleapis.com",
			"cloudresourcemanager.googleapis.com", "dns.googleapis.com",
			"iam.googleapis.com", "iamcredentials.googleapis.com", "serviceusage.googleapis.com",
			"deploymentmanager.googleapis.com",
		},
	}, {
		name:     "Some services present",
		services: []string{"compute.googleapis.com", "serviceusage.googleapis.com"},
		err:      "the following required services are not enabled in this project: cloudresourcemanager.googleapis.com,dns.googleapis.com,iam.googleapis.com,iamcredentials.googleapis.com",
	}}

	errForbidden := &googleapi.Error{Code: http.StatusForbidden}

	for _, test := range cases {
		t.Run(test.name, func(t *testing.T) {
			mockCtrl := gomock.NewController(t)
			defer mockCtrl.Finish()
			gcpClient := mock.NewMockAPI(mockCtrl)

			if !sets.NewString(test.services...).Has("serviceusage.googleapis.com") {
				gcpClient.EXPECT().GetEnabledServices(gomock.Any(), gomock.Any()).Return(nil, errForbidden).AnyTimes()
			} else {
				gcpClient.EXPECT().GetEnabledServices(gomock.Any(), gomock.Any()).Return(test.services, nil).AnyTimes()
			}
			err := ValidateEnabledServices(context.TODO(), gcpClient, "")
			if test.err == "" {
				assert.NoError(t, err)
			} else {
				assert.Regexp(t, test.err, err)
			}
		})
	}
}

func TestValidateCredentialMode(t *testing.T) {
	cases := []struct {
		name       string
		creds      types.CredentialsMode
		emptyCreds bool
		err        string
	}{{
		name:       "missing json with manual creds",
		creds:      types.ManualCredentialsMode,
		emptyCreds: true,
	}, {
		name:       "missing json without manual creds",
		creds:      types.PassthroughCredentialsMode,
		emptyCreds: true,
		err:        "credentialsMode: Forbidden: environmental authentication is only supported with Manual credentials mode",
	}, {
		name:       "supplied json with manual creds",
		creds:      types.ManualCredentialsMode,
		emptyCreds: false,
	}, {
		name:       "supplied json without manual creds",
		creds:      types.PassthroughCredentialsMode,
		emptyCreds: false,
	}}

	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	// Client where the creds are empty
	gcpClientEmptyCreds := mock.NewMockAPI(mockCtrl)
	gcpClientEmptyCreds.EXPECT().GetCredentials().Return(&googleoauth.Credentials{}).AnyTimes()

	// Client that contains creds
	gcpClientWithCreds := mock.NewMockAPI(mockCtrl)
	gcpClientWithCreds.EXPECT().GetCredentials().Return(&googleoauth.Credentials{JSON: []byte("fake creds")}).AnyTimes()

	for _, test := range cases {
		t.Run(test.name, func(t *testing.T) {

			ic := types.InstallConfig{
				ObjectMeta:      metav1.ObjectMeta{Name: "cluster-name"},
				BaseDomain:      "base-domain",
				Platform:        types.Platform{GCP: &gcp.Platform{ProjectID: "project-id"}},
				CredentialsMode: test.creds,
			}

			var err error
			if test.emptyCreds {
				err = ValidateCredentialMode(gcpClientEmptyCreds, &ic)
			} else {
				err = ValidateCredentialMode(gcpClientWithCreds, &ic)
			}

			if test.err == "" {
				assert.NoError(t, err)
			} else {
				assert.Regexp(t, test.err, err)
			}
		})
	}
}

func TestValidateZones(t *testing.T) {
	validZonesDefaultMachine := func(ic *types.InstallConfig) {
		ic.Platform.GCP.DefaultMachinePlatform.Zones = []string{"us-central1-a", "us-central1-c"}
	}
	validZonesControlPlane := func(ic *types.InstallConfig) {
		ic.ControlPlane.Platform.GCP.Zones = []string{"us-central1-a", "us-central1-b"}
	}
	validZonesCompute := func(ic *types.InstallConfig) {
		ic.Compute[0].Platform.GCP.Zones = []string{"us-central1-b", "us-central1-c", "us-central1-d"}
	}
	invalidZonesDefaultMachine := func(ic *types.InstallConfig) {
		ic.Platform.GCP.DefaultMachinePlatform.Zones = []string{"us-central1-a", "us-central1-x", "us-central1-y"}
	}
	invalidZonesControlPlane := func(ic *types.InstallConfig) {
		ic.ControlPlane.Platform.GCP.Zones = []string{"us-central1-d", "us-central1-x", "us-central1-y"}
	}
	invalidZonesCompute := func(ic *types.InstallConfig) {
		ic.Compute[0].Platform.GCP.Zones = []string{"us-central1-y", "us-central1-z", "us-central1-w"}
	}

	cases := []struct {
		name           string
		edits          editFunctions
		expectedError  bool
		expectedErrMsg string
	}{
		{
			name:           "Valid zones for defaultMachine",
			edits:          editFunctions{validZonesDefaultMachine},
			expectedError:  false,
			expectedErrMsg: "",
		},
		{
			name:           "Invalid zones for defaultMachine",
			edits:          editFunctions{invalidZonesDefaultMachine},
			expectedError:  true,
			expectedErrMsg: `^\[platform.gcp.defaultMachinePlatform.zones: Invalid value: \[\]string\{"us\-central1\-x", "us\-central1\-y"\}: zone\(s\) not found in region\]$`,
		},
		{
			name:           "Valid zones for controlPlane",
			edits:          editFunctions{validZonesControlPlane},
			expectedError:  false,
			expectedErrMsg: "",
		},
		{
			name:           "Invalid zones for controlPlane",
			edits:          editFunctions{invalidZonesControlPlane},
			expectedError:  true,
			expectedErrMsg: `^\[controlPlane.platform.gcp.zones: Invalid value: \[\]string\{"us\-central1\-x", "us\-central1\-y"\}: zone\(s\) not found in region\]$`,
		},
		{
			name:           "Valid zones for compute",
			edits:          editFunctions{validZonesCompute},
			expectedError:  false,
			expectedErrMsg: "",
		},
		{
			name:           "Invalid zones for compute",
			edits:          editFunctions{invalidZonesCompute},
			expectedError:  true,
			expectedErrMsg: `^\[compute\[0\].platform.gcp.zones: Invalid value: \[\]string\{"us\-central1\-w", "us\-central1\-y", "us\-central1\-z"\}: zone\(s\) not found in region\]$`,
		},
	}

	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	gcpClient := mock.NewMockAPI(mockCtrl)

	validZones := []*compute.Zone{
		{Name: "us-central1-a"},
		{Name: "us-central1-b"},
		{Name: "us-central1-c"},
		{Name: "us-central1-d"},
	}

	// Should get the list of zones.
	gcpClient.EXPECT().GetZones(gomock.Any(), gomock.Any(), gomock.Any()).Return(validZones, nil).AnyTimes()

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			editedInstallConfig := validInstallConfig()
			for _, edit := range tc.edits {
				edit(editedInstallConfig)
			}

			errs := validateZones(gcpClient, editedInstallConfig)
			if tc.expectedError {
				assert.Regexp(t, tc.expectedErrMsg, errs)
			} else {
				assert.Empty(t, errs)
			}
		})
	}
}

func TestValidateInstanceType(t *testing.T) {
	cases := []struct {
		name           string
		zones          []string
		instanceType   string
		expectedError  bool
		expectedErrMsg string
	}{
		{
			name:           "Valid instance type with min requirements and no zones specified",
			zones:          []string{},
			instanceType:   "n1-standard-4",
			expectedError:  false,
			expectedErrMsg: "",
		},
		{
			name:           "Valid instance type with min requirements and valid zones specified",
			zones:          []string{"a", "b"},
			instanceType:   "n1-standard-4",
			expectedError:  false,
			expectedErrMsg: "",
		},
		{
			name:           "Valid instance type with min requirements and invalid zones specified",
			zones:          []string{"a", "b", "d", "x", "y"},
			instanceType:   "n1-standard-4",
			expectedError:  true,
			expectedErrMsg: `^\[instance.type: Invalid value: "n1\-standard\-4": instance type not available in zones: \[x y\]\]$`,
		},
		{
			name:           "Valid instance fails min requirements and no zones specified",
			zones:          []string{},
			instanceType:   "n1-standard-2",
			expectedError:  true,
			expectedErrMsg: `^\[instance.type: Invalid value: "n1\-standard\-2": instance type does not meet minimum resource requirements of 4 vCPUs instance.type: Invalid value: "n1\-standard\-2": instance type does not meet minimum resource requirements of 15360 MB Memory\]$`,
		},
		{
			name:           "Valid instance fails min requirements and valid zones specified",
			zones:          []string{"a", "b"},
			instanceType:   "n1-standard-1",
			expectedError:  true,
			expectedErrMsg: `^\[instance.type: Invalid value: "n1\-standard\-1": instance type does not meet minimum resource requirements of 4 vCPUs instance.type: Invalid value: "n1\-standard\-1": instance type does not meet minimum resource requirements of 15360 MB Memory\]$`,
		},
		{
			name:           "Valid instance fails min requirements and invalid zones specified",
			zones:          []string{"a", "x", "y"},
			instanceType:   "n1-standard-1",
			expectedError:  true,
			expectedErrMsg: `^\[instance.type: Invalid value: "n1\-standard\-1": instance type not available in zones: \[x y\] instance.type: Invalid value: "n1\-standard\-1": instance type does not meet minimum resource requirements of 4 vCPUs instance.type: Invalid value: "n1\-standard\-1": instance type does not meet minimum resource requirements of 15360 MB Memory\]$`,
		},
		{
			name:           "Invalid instance and no zones specified",
			zones:          []string{},
			instanceType:   "invalid-instance-1",
			expectedError:  true,
			expectedErrMsg: `^\[<nil>: Internal error: 404\]$`,
		},
		{
			name:           "Invalid instance and valid zones specified",
			zones:          []string{"a", "b"},
			instanceType:   "invalid-instance-1",
			expectedError:  true,
			expectedErrMsg: `^\[<nil>: Internal error: 404\]$`,
		},
		{
			name:           "Invalid instance and invalid zones specified",
			zones:          []string{"x", "y", "z"},
			instanceType:   "invalid-instance-1",
			expectedError:  true,
			expectedErrMsg: `^\[<nil>: Internal error: 404\]$`,
		},
	}

	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	gcpClient := mock.NewMockAPI(mockCtrl)

	// Should return the machine type as specified.
	for key, value := range machineTypeAPIResult {
		gcpClient.EXPECT().GetMachineTypeWithZones(gomock.Any(), gomock.Any(), gomock.Any(), key).Return(value, sets.New("a", "b", "c", "d"), nil).AnyTimes()
	}
	// When passed incorrect machine type, the API returns nil.
	gcpClient.EXPECT().GetMachineTypeWithZones(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(nil, nil, fmt.Errorf("404")).AnyTimes()

	zones := []*compute.Zone{
		{Name: "a"},
		{Name: "b"},
		{Name: "c"},
	}
	gcpClient.EXPECT().GetZones(gomock.Any(), gomock.Any(), gomock.Any()).Return(zones, nil).AnyTimes()

	for _, test := range cases {
		t.Run(test.name, func(t *testing.T) {
			errs := ValidateInstanceType(gcpClient, field.NewPath("instance"), "project-id", "region", test.zones, test.instanceType, controlPlaneReq)
			if test.expectedError {
				assert.Regexp(t, test.expectedErrMsg, errs)
			} else {
				assert.Empty(t, errs)
			}
		})
	}
}

func TestValidateServiceAccountPresent(t *testing.T) {
	cases := []struct {
		name             string
		creds            *googleoauth.Credentials
		serviceAccount   string
		networkProjectID string
		expectedError    string
	}{
		{
			name:  "Test no network project ID",
			creds: &googleoauth.Credentials{},
		},
		{
			name:             "Test network project ID with service account",
			creds:            &googleoauth.Credentials{},
			serviceAccount:   "test-service-account",
			networkProjectID: "test-network-project",
		},
		{
			name:             "Test network project ID service account and creds",
			creds:            &googleoauth.Credentials{JSON: []byte("{}")},
			serviceAccount:   "test-service-account",
			networkProjectID: "test-network-project",
		},
		{
			name:             "Test network project ID no creds",
			creds:            &googleoauth.Credentials{JSON: nil},
			networkProjectID: "test-network-project",
			expectedError:    "controlPlane.platform.gcp.serviceAccount: Required value: service account must be provided when authentication credentials do not provide a service account",
		},
	}

	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	for _, test := range cases {
		gcpClient := mock.NewMockAPI(mockCtrl)
		if test.networkProjectID != "" {
			gcpClient.EXPECT().GetCredentials().Return(test.creds)
		}

		t.Run(test.name, func(t *testing.T) {
			ic := types.InstallConfig{
				ObjectMeta:      metav1.ObjectMeta{Name: "cluster-name"},
				BaseDomain:      "base-domain",
				Platform:        types.Platform{GCP: &gcp.Platform{ProjectID: "project-id", NetworkProjectID: test.networkProjectID}},
				CredentialsMode: types.PassthroughCredentialsMode,
				ControlPlane: &types.MachinePool{
					Platform: types.MachinePoolPlatform{
						GCP: &gcp.MachinePool{
							ServiceAccount: test.serviceAccount,
						},
					},
				},
			}

			errorList := validateServiceAccountPresent(gcpClient, &ic)
			if errorList == nil && test.expectedError == "" {
				assert.NoError(t, errorList.ToAggregate())
			} else {
				assert.Regexp(t, test.expectedError, errorList.ToAggregate())
			}
		})
	}
}
