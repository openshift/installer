package terraform

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDiagnose(t *testing.T) {
	cases := []struct {
		input string
		err   string
	}{{
		input: `Error: Error creating Blob "rhcoskltwa.vhd" (Container "vhd" / Account "clusterkltwa"): Error copy/waiting: 
   on ../tmp/openshift-install-348626978/main.tf line 169, in resource "azurerm_storage_blob" "rhcos_image":"
  169: resource "azurerm_storage_blob" "rhcos_image" {
`,
		err: `error\(Timeout\) from Infrastructure Provider: Copying the VHD to user environment was too slow, and timeout was reached for the success\.`,
	}, {
		input: `Error: Error Creating/Updating Subnet "xxxx-master-subnet" (Virtual Network "xxxx-vnet" / Resource Group "xxxx-rg"): network.SubnetsClient#CreateOrUpdate: Failure sending request: StatusCode=0 -- Original Error: autorest/azure: Service returned an error. Status=<nil> Code="AnotherOperationInProgress" Message="Another operation on this or dependent resource is in progress. To retrieve status of the operation use uri: https://management.azure.com/subscriptions/d38f1e38-4bed-438e-b227-833f997adf6a/providers/Microsoft.Network/locations/eastus2/operations/62c8a417-7168-464f-83e6-96912bd6b30a?api-version=2019-09-01." Details=[]

  on ../tmp/openshift-install-513947104/vnet/vnet.tf line 10, in resource "azurerm_subnet" "master_subnet":"
  10: resource "azurerm_subnet" "master_subnet" {
`,
		err: `error\(AzureMultiOperationFailure\) from Infrastructure Provider: Creating Subnets failed because Azure could not process multiple operations\.`,
	}, {
		input: `Error: Error Creating/Updating Public IP "xxxx-bootstrap-pip-v4" (Resource Group "xxxx-rg"): network.PublicIPAddressesClient#CreateOrUpdate: Failure sending request: StatusCode=400 -- Original Error: Code="PublicIPCountLimitReached" Message="Cannot create more than 50 public IP addresses for this subscription in this region." Details=[]

  on ../tmp/openshift-install-172932975/bootstrap/main.tf line 65, in resource "azurerm_public_ip" "bootstrap_public_ip_v4":
  65: resource "azurerm_public_ip" "bootstrap_public_ip_v4" {
`,

		err: `error\(AzureQuotaLimitExceeded\) from Infrastructure Provider: Service limits exceeded for Public IPs in the the subscriptions for the region. Requesting increase in quota should fix the error\.`,
	}, {
		input: `Error: Code="OSProvisioningTimedOut" Message="OS Provisioning for VM 'xxxx-master-2' did not finish in the allotted time. The VM may still finish provisioning successfully. Please check provisioning state later. Also, make sure the image has been properly prepared (generalized).\\r\\n * Instructions for Windows: https://azure.microsoft.com/documentation/articles/virtual-machines-windows-upload-image/ \\r\\n * Instructions for Linux: https://azure.microsoft.com/documentation/articles/virtual-machines-linux-capture-image/ "

  on ../tmp/openshift-install-172932975/master/master.tf line 81, in resource "azurerm_virtual_machine" "master":
  81: resource "azurerm_virtual_machine" "master" {
`,

		err: `error\(AzureVirtualMachineFailure\) from Infrastructure Provider: Some virtual machines failed to provision in alloted time`,
	}, {
		input: `
Error: Error waiting for instance to create: Internal error. Please try again or contact Google Support. (Code: '8712799794455203922')


  on ../tmp/openshift-install-910996711/master/main.tf line 31, in resource "google_compute_instance" "master":
  31: resource "google_compute_instance" "master" {
`,

		err: `error\(GCPComputeBackendTimeout\) from Infrastructure Provider: GCP is experiencing backend service interuptions, the compute instance failed to create in reasonable time\.`,
	}, {
		input: `Error: Error reading Service Account "projects/project-id/serviceAccounts/xxxx-m@project-id.iam.gserviceaccount.com": googleapi: Error 503: The service is currently unavailable., backendError`,

		err: `error\(GCPBackendInternalError\) from Infrastructure Provider: GCP is experiencing backend service interuptions. Please try again or contact Google Support`,
	}, {
		input: `
Error: Error adding instances to InstanceGroup: googleapi: Error 503: Internal error. Please try again or contact Google Support. (Code: 'xxxx'), backendError

  on ../tmp/openshift-install-267295217/bootstrap/main.tf line 87, in resource "google_compute_instance_group" "bootstrap":
  87: resource "google_compute_instance_group" "bootstrap" {
`,

		err: `error\(GCPBackendInternalError\) from Infrastructure Provider: GCP is experiencing backend service interuptions. Please try again or contact Google Support`,
	}, {
		input: `
Error: Error applying IAM policy to project "project-id": Too many conflicts.  Latest error: Error setting IAM policy for project "project-id": googleapi: Error 409: There were concurrent policy changes. Please retry the whole read-modify-write with exponential backoff., aborted

  on ../tmp/openshift-install-392130810/master/main.tf line 26, in resource "google_project_iam_member" "master-service-account-user":
  26: resource "google_project_iam_member" "master-service-account-user" {
`,

		err: `error\(GCPTooManyIAMUpdatesInFlight\) from Infrastructure Provider: There are a lot of IAM updates to the project in flight. Failed after reaching a limit of read-modify-write on conflict backoffs\.`,
	}, {
		input: `
Error: Error retrieving resource group: resources.GroupsClient#Get: Failure responding to request: StatusCode=404 -- Original Error: autorest/azure: Service returned an error. Status=404 Code="ResourceGroupNotFound" Message="Resource group 'xxxxx-rg' could not be found."

  on ../tmp/openshift-install-424775273/main.tf line 124, in resource "azurerm_resource_group" "main":
 124: resource "azurerm_resource_group" "main" {
`,

		err: `error\(AzureEventualConsistencyFailure\) from Infrastructure Provider: Failed to find a resource that was recently created usualy caused by Azure's eventual consistency delays\.`,
	}, {
		input: `
Error: compute.VirtualMachinesClient#CreateOrUpdate: Failure sending request: StatusCode=0 -- Original Error: autorest/azure: Service returned an error. Status=<nil> Code="OperationNotAllowed" Message="Operation could not be completed as it results in exceeding approved Total Regional Cores quota. Additional details - Deployment Model: Resource Manager, Location: centralus, Current Limit: 200, Current Usage: 198, Additional Required: 8, (Minimum) New Limit Required: 206. Submit a request for Quota increase at https://aka.ms/ProdportalCRP/?#create/Microsoft.Support/Parameters/%7B%22subId%22:%225f675811-04fa-483f-9709-ffd8a9da03f0%22,%22pesId%22:%2206bfd9d3-516b-d5c6-5802-169c800dec89%22,%22supportTopicId%22:%22e12e3d1d-7fa0-af33-c6d0-3c50df9658a3%22%7D by specifying parameters listed in the ‘Details’ section for deployment to succeed. Please read more about quota limits at https://docs.microsoft.com/en-us/azure/azure-supportability/regional-quota-requests."

  on ../../../../tmp/openshift-install-941329162/master/master.tf line 81, in resource "azurerm_virtual_machine" "master":
  81: resource "azurerm_virtual_machine" "master" {
`,

		err: `error\(AzureQuotaLimitExceeded\) from Infrastructure Provider: Service limits exceeded for Virtual Machine cores in the the subscriptions for the region\. Requesting increase in quota should fix the error\.`,
	}}

	for _, test := range cases {
		t.Run("", func(t *testing.T) {
			err := Diagnose(test.input)
			if test.err == "" {
				assert.NoError(t, err)
			} else {
				assert.Regexp(t, test.err, err)
			}
		})
	}
}
