package ibmcloud

import (
	"bytes"
	"context"
	"fmt"
	"net/http"
	"os"
	"strings"
	"time"

	ibmcrn "github.com/IBM-Cloud/bluemix-go/crn"
	"github.com/IBM/go-sdk-core/v5/core"
	cosaws "github.com/IBM/ibm-cos-sdk-go/aws"
	"github.com/IBM/ibm-cos-sdk-go/aws/credentials/ibmiam"
	cossession "github.com/IBM/ibm-cos-sdk-go/aws/session"
	ibms3 "github.com/IBM/ibm-cos-sdk-go/service/s3"
	kpclient "github.com/IBM/keyprotect-go-client"
	"github.com/IBM/networking-go-sdk/dnsrecordsv1"
	"github.com/IBM/networking-go-sdk/dnssvcsv1"
	"github.com/IBM/networking-go-sdk/dnszonesv1"
	"github.com/IBM/networking-go-sdk/zonesv1"
	"github.com/IBM/platform-services-go-sdk/globalcatalogv1"
	"github.com/IBM/platform-services-go-sdk/iamidentityv1"
	"github.com/IBM/platform-services-go-sdk/iampolicymanagementv1"
	"github.com/IBM/platform-services-go-sdk/resourcecontrollerv2"
	"github.com/IBM/platform-services-go-sdk/resourcemanagerv2"
	"github.com/IBM/vpc-go-sdk/vpcv1"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"k8s.io/utils/ptr"

	configv1 "github.com/openshift/api/config/v1"
	"github.com/openshift/installer/pkg/asset/installconfig/ibmcloud/responses"
	"github.com/openshift/installer/pkg/types"
	ibmcloudtypes "github.com/openshift/installer/pkg/types/ibmcloud"
)

//go:generate mockgen -source=./client.go -destination=./mock/ibmcloudclient_generated.go -package=mock

// API represents the calls made to the API.
type API interface {
	CreateCOSBucket(ctx context.Context, cosInstanceID string, bucketName string, region string) error
	CreateCOSInstance(ctx context.Context, cosName string, resourceGroupID string) (*resourcecontrollerv2.ResourceInstance, error)
	CreateCOSObject(ctx context.Context, sourceData []byte, fileName string, cosInstanceID string, bucketName string, region string) error
	CreateCISDNSRecord(ctx context.Context, cisInstanceCRN string, zoneID string, recordName string, cname string) error
	CreateDNSServicesDNSRecord(ctx context.Context, dnsInstanceID string, zoneID string, recordName string, cname string) error
	CreateIAMAuthorizationPolicy(tx context.Context, sourceServiceName string, sourceServiceResourceType string, targetServiceName string, targetServiceInstanceID string, roles []string) error
	CreateResourceGroup(ctx context.Context, rgName string) error
	GetAPIKey() string
	GetAuthenticatorAPIKeyDetails(ctx context.Context) (*iamidentityv1.APIKey, error)
	GetCISInstance(ctx context.Context, crnstr string) (*resourcecontrollerv2.ResourceInstance, error)
	GetCOSBucketByName(ctx context.Context, cosInstanceID string, bucketName string, region string) (*ibms3.Bucket, error)
	GetCOSInstanceByName(ctx context.Context, cosName string) (*resourcecontrollerv2.ResourceInstance, error)
	GetDNSInstance(ctx context.Context, crnstr string) (*resourcecontrollerv2.ResourceInstance, error)
	GetDNSInstancePermittedNetworks(ctx context.Context, dnsID string, dnsZone string) ([]string, error)
	GetDedicatedHostByName(ctx context.Context, name string, region string) (*vpcv1.DedicatedHost, error)
	GetDedicatedHostProfiles(ctx context.Context, region string) ([]vpcv1.DedicatedHostProfile, error)
	GetDNSRecordsByName(ctx context.Context, crnstr string, zoneID string, recordName string) ([]dnsrecordsv1.DnsrecordDetails, error)
	GetDNSZoneIDByName(ctx context.Context, name string, publish types.PublishingStrategy) (string, error)
	GetDNSZones(ctx context.Context, publish types.PublishingStrategy) ([]responses.DNSZoneResponse, error)
	GetEncryptionKey(ctx context.Context, keyCRN string) (*responses.EncryptionKeyResponse, error)
	GetIBMCloudRegions(ctx context.Context) (map[string]string, error)
	GetLoadBalancer(ctx context.Context, loadBalancerID string) (*vpcv1.LoadBalancer, error)
	GetResourceGroups(ctx context.Context) ([]resourcemanagerv2.ResourceGroup, error)
	GetResourceGroup(ctx context.Context, nameOrID string) (*resourcemanagerv2.ResourceGroup, error)
	GetSSHKeyByPublicKey(ctx context.Context, publicKey string) (*vpcv1.Key, error)
	GetSubnet(ctx context.Context, subnetID string) (*vpcv1.Subnet, error)
	GetSubnetByName(ctx context.Context, subnetName string, region string) (*vpcv1.Subnet, error)
	GetVSIProfiles(ctx context.Context) ([]vpcv1.InstanceProfile, error)
	GetVPC(ctx context.Context, vpcID string) (*vpcv1.VPC, error)
	GetVPCs(ctx context.Context, region string) ([]vpcv1.VPC, error)
	GetVPCByName(ctx context.Context, vpcName string) (*vpcv1.VPC, error)
	GetVPCZonesForRegion(ctx context.Context, region string) ([]string, error)
	SetVPCServiceURLForRegion(ctx context.Context, region string) error
}

// Client makes calls to the IBM Cloud API.
type Client struct {
	apiKey        string
	catalogAPI    *globalcatalogv1.GlobalCatalogV1
	managementAPI *resourcemanagerv2.ResourceManagerV2
	controllerAPI *resourcecontrollerv2.ResourceControllerV2
	vpcAPI        *vpcv1.VpcV1

	// A set of overriding endpoints for IBM Cloud Services
	serviceEndpoints []configv1.IBMCloudServiceEndpoint
	// Cache endpoint override for IBM Cloud IAM, if one was provided in serviceEndpoints
	iamServiceEndpointOverride string
}

// InstanceType is the IBM Cloud network services type being used
type InstanceType string

const (
	// CISInstanceType is a Cloud Internet Services InstanceType
	CISInstanceType InstanceType = "CIS"
	// DNSInstanceType is a DNS Services InstanceType
	DNSInstanceType InstanceType = "DNS"

	// cisServiceID is the Cloud Internet Services' catalog service ID.
	cisServiceID = "75874a60-cb12-11e7-948e-37ac098eb1b9"
	// cosServiceID is the Cloud Object Storage catalog service ID.
	cosServiceID = "dff97f5c-bc5e-4455-b470-411c3edbe49c"
	// cosSerivcePlanID is the Cloud Object Storage catalog Standard service plan ID.
	cosServicePlanID = "744bfc56-d12c-4866-88d5-dac9139e0e5d"
	// dnsServiceID is the DNS Services' catalog service ID.
	dnsServiceID = "b4ed8a30-936f-11e9-b289-1d079699cbe5"

	// hyperProtectCRNServiceName is the service name within the IBM Cloud CRN for the Hyper Protect service.
	hyperProtectCRNServiceName = "hs-crypto"
	// keyProtectCRNServiceName is the service name within the IBM Cloud CRN for the Key Protect service.
	keyProtectCRNServiceName = "kms"

	// cosDefaultURLTEmplate is the default URL endpoint template, with region subsitution, for IBM Cloud Object Storage service.
	cosDefaultURLTemplate = "s3.%s.cloud-object-storage.appdomain.cloud"
	// hyperProtectDefaultURLTemplate is the default URL endpoint template, with region substitution, for IBM Cloud Hyper Protect service.
	hyperProtectDefaultURLTemplate = "https://api.%s.hs-crypto.cloud.ibm.com"
	// iamTokenDefaultURL is the default URL endpoint for IBM Cloud IAM token service.
	iamTokenDefaultURL = "https://iam.cloud.ibm.com/identity/token" // #nosec G101 // this is the URL for IBM Cloud IAM tokens, not a credential
	// iamTokenPath is the URL path, to add to override IAM endpoints, for the IBM Cloud IAM token service.
	iamTokenPath = "identity/token" // #nosec G101 // this is the URI path for IBM Cloud IAM tokens, not a credential
	// keyProtectDefaultURLTemplate is the default URL endpoint template, with region substitution, for IBM Cloud Key Protect service.
	keyProtectDefaultURLTemplate = "https://%s.kms.cloud.ibm.com"
)

// COSResourceNotFoundError represents an error for a COS resource that is not found.
type COSResourceNotFoundError struct{}

// Error returns the error message for the COSResourceNotFoundError error type.
func (e *COSResourceNotFoundError) Error() string {
	return "COS Resource Not Found"
}

// VPCResourceNotFoundError represents an error for a VPC resoruce that is not found.
type VPCResourceNotFoundError struct{}

// Error returns the error message for the VPCResourceNotFoundError error type.
func (e *VPCResourceNotFoundError) Error() string {
	return "Not Found"
}

// NewClient initializes a client with any provided endpoint overrides.
func NewClient(endpoints []configv1.IBMCloudServiceEndpoint) (*Client, error) {
	apiKey := os.Getenv("IC_API_KEY")

	client := &Client{
		apiKey:           apiKey,
		serviceEndpoints: endpoints,
	}

	// Look for an override to IBM Cloud IAM service, preventing searching each time its necessary
	for _, endpoint := range endpoints {
		if endpoint.Name == configv1.IBMCloudServiceIAM {
			client.iamServiceEndpointOverride = endpoint.URL
			break
		}
	}

	if err := client.loadSDKServices(); err != nil {
		return nil, errors.Wrap(err, "failed to load IBM SDK services")
	}

	return client, nil
}

func (c *Client) loadSDKServices() error {
	servicesToLoad := []func() error{
		c.loadGlobalCatalogAPI,
		c.loadResourceManagementAPI,
		c.loadResourceControllerAPI,
		c.loadVPCV1API,
	}

	// Call all the load functions.
	for _, fn := range servicesToLoad {
		if err := fn(); err != nil {
			return err
		}
	}

	return nil
}

// CreateCOSBucket will create a new COS Bucket in the COS Instance, based on the Instance ID.
func (c *Client) CreateCOSBucket(ctx context.Context, cosInstanceID string, bucketName string, region string) error {
	localContext, cancel := context.WithTimeout(ctx, 1*time.Minute)
	defer cancel()
	cosClient := c.getCOSClient(cosInstanceID, region)
	// Setup Location Constraint, for now we always default to the Smart tier
	locationConstraint := fmt.Sprintf("%s-%s", region, "smart")

	options := &ibms3.CreateBucketInput{
		Bucket: ptr.To(bucketName),
		CreateBucketConfiguration: &ibms3.CreateBucketConfiguration{
			LocationConstraint: ptr.To(locationConstraint),
		},
	}
	if _, err := cosClient.CreateBucketWithContext(localContext, options); err != nil {
		return fmt.Errorf("failed to create cos bucket: %w", err)
	}
	return nil
}

// CreateCOSInstance will create a new COS Instance and return the ResourceInstance details.
func (c *Client) CreateCOSInstance(ctx context.Context, cosName string, resourceGroupID string) (*resourcecontrollerv2.ResourceInstance, error) {
	localContext, cancel := context.WithTimeout(ctx, 1*time.Minute)
	defer cancel()

	createOptions := c.controllerAPI.NewCreateResourceInstanceOptions(cosName, "Global", resourceGroupID, cosServicePlanID)

	instance, _, err := c.controllerAPI.CreateResourceInstanceWithContext(localContext, createOptions)
	if err != nil {
		return nil, fmt.Errorf("failed creating new COS instance: %w", err)
	}

	return instance, nil
}

// CreateCOSObject will (upload) create a new COS object in the designated COS instance and bucket.
func (c *Client) CreateCOSObject(ctx context.Context, sourceData []byte, fileName string, cosInstanceID string, bucketName string, region string) error {
	cosClient := c.getCOSClient(cosInstanceID, region)

	options := &ibms3.PutObjectInput{
		Body:   cosaws.ReadSeekCloser(bytes.NewReader(sourceData)),
		Bucket: cosaws.String(bucketName),
		Key:    cosaws.String(fileName),
	}

	if _, err := cosClient.PutObject(options); err != nil {
		return fmt.Errorf("failed creating cos object: %w", err)
	}
	return nil
}

// CreateCISDNSRecord creates a DNS Record in the IBM Cloud Internet Services (CIS) zone based on the provided Load Balancer.
func (c *Client) CreateCISDNSRecord(ctx context.Context, cisInstanceCRN string, zoneID string, recordName string, cname string) error {
	localContext, cancel := context.WithTimeout(ctx, 1*time.Minute)
	defer cancel()

	dnsRecordsService, err := c.getDNSRecordsAPI(cisInstanceCRN, zoneID)
	if err != nil {
		return fmt.Errorf("failed to create dns records service to create dns record: %w", err)
	}

	// Build the new DNS Record options.
	recordOptions := dnsRecordsService.NewCreateDnsRecordOptions()
	recordOptions.SetName(recordName)
	recordOptions.SetType(dnsrecordsv1.CreateDnsRecordOptions_Type_Cname)
	recordOptions.SetContent(cname)

	// Create new DNS Record.
	logrus.Debugf("creating cis dns record: recordName=%s, cname=%s", recordName, cname)
	recordDetails, _, err := dnsRecordsService.CreateDnsRecordWithContext(localContext, recordOptions)
	if err != nil {
		return fmt.Errorf("failed to create dns record %s: %w", recordName, err)
	}
	logrus.Debugf("created new dns record: recordName=%s, recordID=%s", recordName, *recordDetails.Result.ID)
	return nil
}

// CreateDNSServicesDNSRecord create a DNS Record in the DNS Serivces zone, based on provided the Load Balancer.
func (c *Client) CreateDNSServicesDNSRecord(ctx context.Context, dnsInstanceID string, zoneID string, recordName string, cname string) error {
	localContext, cancel := context.WithTimeout(ctx, 1*time.Minute)
	defer cancel()

	dnsService, err := c.getDNSServicesAPI()
	if err != nil {
		return fmt.Errorf("failed to create dns services service to create dns record: %w", err)
	}

	// Build the new DNS Record options.
	cnameRecord, err := dnsService.NewResourceRecordInputRdataRdataCnameRecord(cname)
	if err != nil {
		return fmt.Errorf("failed to create rdata cname record for dns services dns record: %w", err)
	}
	recordOptions := dnsService.NewCreateResourceRecordOptions(dnsInstanceID, zoneID)
	recordOptions.SetName(recordName)
	recordOptions.SetRdata(cnameRecord)
	recordOptions.SetTTL(60)
	recordOptions.SetType("CNAME")

	// Create new DNS Record.
	logrus.Debugf("creating dns services dns record: recordName=%s, cname=%s", recordName, cname)
	recordDetails, _, err := dnsService.CreateResourceRecordWithContext(localContext, recordOptions)
	if err != nil {
		return fmt.Errorf("failed to create dns record %s: %w", recordName, err)
	}
	logrus.Debugf("created new dns record: recordName=%s, recordID=%s", recordName, *recordDetails.ID)
	return nil
}

// CreateIAMAuthorizationPolicy creates a new IAM Authorization policy for read access to VPC to a COS Instance.
func (c *Client) CreateIAMAuthorizationPolicy(ctx context.Context, sourceServiceName string, sourceServiceResourceType string, targetServiceName string, targetServiceInstanceID string, roles []string) error {
	accountIDKeyPtr := ptr.To("accountId")
	resourceTypeKeyPtr := ptr.To("resourceType")
	serviceInstanceKeyPtr := ptr.To("serviceInstance")
	serviceNameKeyPtr := ptr.To("serviceName")
	stringEqualsOperatorPtr := ptr.To("stringEquals")

	apiKeyDetails, err := c.GetAuthenticatorAPIKeyDetails(ctx)
	if err != nil {
		return fmt.Errorf("failed collecting account ID: %w", err)
	}

	policyRoles := make([]iampolicymanagementv1.Roles, 0, len(roles))

	for _, role := range roles {
		policyRoles = append(policyRoles, iampolicymanagementv1.Roles{
			RoleID: ptr.To(role),
		})
	}

	policyControl := &iampolicymanagementv1.Control{
		Grant: &iampolicymanagementv1.Grant{
			Roles: policyRoles,
		},
	}

	// setup the source service policy details (VPC Custom Image)
	policySubject := &iampolicymanagementv1.V2PolicySubject{
		Attributes: []iampolicymanagementv1.V2PolicySubjectAttribute{
			{
				Key:      serviceNameKeyPtr,
				Operator: stringEqualsOperatorPtr,
				Value:    ptr.To(sourceServiceName),
			},
			{
				Key:      accountIDKeyPtr,
				Operator: stringEqualsOperatorPtr,
				Value:    apiKeyDetails.AccountID,
			},
			{
				Key:      resourceTypeKeyPtr,
				Operator: stringEqualsOperatorPtr,
				Value:    ptr.To(sourceServiceResourceType),
			},
		},
	}

	// setup the target resource policy details (COS)
	policyResource := &iampolicymanagementv1.V2PolicyResource{
		Attributes: []iampolicymanagementv1.V2PolicyResourceAttribute{
			{
				Key:      serviceNameKeyPtr,
				Operator: stringEqualsOperatorPtr,
				Value:    ptr.To(targetServiceName),
			},
			{
				Key:      accountIDKeyPtr,
				Operator: stringEqualsOperatorPtr,
				Value:    apiKeyDetails.AccountID,
			},
			{
				Key:      serviceInstanceKeyPtr,
				Operator: stringEqualsOperatorPtr,
				Value:    ptr.To(targetServiceInstanceID),
			},
		},
	}

	authenticator, err := NewIamAuthenticator(c.GetAPIKey(), c.iamServiceEndpointOverride)
	if err != nil {
		return fmt.Errorf("failed setting up IAM policy management service: %w", err)
	}
	iamPolicyManagementServiceOptions := &iampolicymanagementv1.IamPolicyManagementV1Options{
		Authenticator: authenticator,
	}
	if c.iamServiceEndpointOverride != "" {
		iamPolicyManagementServiceOptions.URL = c.iamServiceEndpointOverride
	}

	iamPolicyManagementService, err := iampolicymanagementv1.NewIamPolicyManagementV1(iamPolicyManagementServiceOptions)
	if err != nil {
		return fmt.Errorf("failed creation IAM policy management service: %w", err)
	}

	options := iamPolicyManagementService.NewCreateV2PolicyOptions(policyControl, "authorization")
	options.SetSubject(policySubject)
	options.SetResource(policyResource)

	if _, _, err = iamPolicyManagementService.CreateV2Policy(options); err != nil {
		return fmt.Errorf("failed creating IAM authorization policy: %w", err)
	}

	return nil
}

// CreateResourceGroup creates a new IBM Cloud Resource Group.
func (c *Client) CreateResourceGroup(ctx context.Context, rgName string) error {
	localContext, cancel := context.WithTimeout(ctx, 1*time.Minute)
	defer cancel()

	// Get the Account ID
	apiKeyDetails, err := c.GetAuthenticatorAPIKeyDetails(localContext)
	if err != nil {
		return fmt.Errorf("failed retrieving account ID: %w", err)
	}

	createRGOptions := c.managementAPI.NewCreateResourceGroupOptions()
	createRGOptions.SetName(rgName)
	createRGOptions.SetAccountID(*apiKeyDetails.AccountID)

	// Create the Resource Group
	if _, _, err = c.managementAPI.CreateResourceGroupWithContext(localContext, createRGOptions); err != nil {
		return fmt.Errorf("failed creating new resource group: %w", err)
	}

	return nil
}

// GetAPIKey gets the API Key.
func (c *Client) GetAPIKey() string {
	return c.apiKey
}

// GetAuthenticatorAPIKeyDetails gets detailed information on the API key used
// for authentication to the IBM Cloud APIs
func (c *Client) GetAuthenticatorAPIKeyDetails(ctx context.Context) (*iamidentityv1.APIKey, error) {
	authenticator, err := NewIamAuthenticator(c.GetAPIKey(), c.iamServiceEndpointOverride)
	if err != nil {
		return nil, err
	}
	newIamIdentityV1Options := iamidentityv1.IamIdentityV1Options{
		Authenticator: authenticator,
	}
	// If an IAM service endpoint override was provided, pass it along to override the default IAM service endpoint
	if c.iamServiceEndpointOverride != "" {
		newIamIdentityV1Options.URL = c.iamServiceEndpointOverride
	}
	iamIdentityService, err := iamidentityv1.NewIamIdentityV1(&newIamIdentityV1Options)
	if err != nil {
		return nil, err
	}

	options := iamIdentityService.NewGetAPIKeysDetailsOptions()
	options.SetIamAPIKey(c.GetAPIKey())
	details, _, err := iamIdentityService.GetAPIKeysDetailsWithContext(ctx, options)
	if err != nil {
		return nil, err
	}
	return details, nil
}

// getInstance gets a specific DNS or CIS instance by its CRN.
func (c *Client) getInstance(ctx context.Context, crnstr string, iType InstanceType) (*resourcecontrollerv2.ResourceInstance, error) {
	_, cancel := context.WithTimeout(ctx, 1*time.Minute)
	defer cancel()

	options := c.controllerAPI.NewGetResourceInstanceOptions(crnstr)
	resourceInstance, _, err := c.controllerAPI.GetResourceInstance(options)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to get %s instances", iType)
	}

	return resourceInstance, nil
}

// GetCISInstance gets a specific Cloud Internet Services by its CRN.
func (c *Client) GetCISInstance(ctx context.Context, crnstr string) (*resourcecontrollerv2.ResourceInstance, error) {
	return c.getInstance(ctx, crnstr, CISInstanceType)
}

// GetCOSBucketByName will get the COS Bucket that matches the name provided.
func (c *Client) GetCOSBucketByName(ctx context.Context, cosInstanceID string, bucketName string, region string) (*ibms3.Bucket, error) {
	localContext, cancel := context.WithTimeout(ctx, 1*time.Minute)
	defer cancel()

	cosClient := c.getCOSClient(cosInstanceID, region)

	options := &ibms3.ListBucketsInput{
		IBMServiceInstanceId: ptr.To(cosInstanceID),
	}
	bucketOutput, err := cosClient.ListBucketsWithContext(localContext, options)
	if err != nil {
		return nil, fmt.Errorf("failed listing buckets: %w", err)
	}

	for _, bucket := range bucketOutput.Buckets {
		if *bucket.Name == bucketName {
			return bucket, nil
		}
	}

	return nil, fmt.Errorf("failed to find bucket '%s' in instance %s", bucketName, cosInstanceID)
}

// getCOSClient returns a new IBM Cloud COS client session.
func (c *Client) getCOSClient(cosInstanceID string, region string) *ibms3.S3 {
	config := cosaws.NewConfig()

	// If an IAM service endpoint override was provided, use it to build the auth endpoint for the COS client (default is used for empty string)
	var authEndpoint string
	if c.iamServiceEndpointOverride != "" {
		authEndpoint = fmt.Sprintf("%s/%s", c.iamServiceEndpointOverride, iamTokenPath)
	}

	// Setup IAM credentials for COS client, passing in the IAM auth endpoint
	config.WithCredentials(ibmiam.NewStaticCredentials(cosaws.NewConfig(), authEndpoint, c.apiKey, cosInstanceID))
	config.WithEndpoint(fmt.Sprintf(cosDefaultURLTemplate, region))

	// If a COS service endpoint override was specified, set it in the COS config
	if overrideURL := ibmcloudtypes.CheckServiceEndpointOverride(configv1.IBMCloudServiceCOS, c.serviceEndpoints); overrideURL != "" {
		config.WithEndpoint(overrideURL)
	}

	sess := cossession.Must(cossession.NewSession())
	return ibms3.New(sess, config)
}

// GetCOSInstanceByName will get the COS Instance (ResourceInstance) that matches the name provided.
func (c *Client) GetCOSInstanceByName(ctx context.Context, cosName string) (*resourcecontrollerv2.ResourceInstance, error) {
	localContext, cancel := context.WithTimeout(ctx, 1*time.Minute)
	defer cancel()

	options := c.controllerAPI.NewListResourceInstancesOptions()
	options.SetResourceID(cosServiceID)

	listResourceInstanceResponse, _, err := c.controllerAPI.ListResourceInstancesWithContext(localContext, options)
	if err != nil {
		return nil, fmt.Errorf("failed to list cos instances: %w", err)
	}
	for _, instance := range listResourceInstanceResponse.Resources {
		if *instance.Name == cosName {
			return &instance, nil
		}
	}

	return nil, &COSResourceNotFoundError{}
}

// GetDNSInstance gets a specific DNS Services instance by its CRN.
func (c *Client) GetDNSInstance(ctx context.Context, crnstr string) (*resourcecontrollerv2.ResourceInstance, error) {
	return c.getInstance(ctx, crnstr, DNSInstanceType)
}

// GetDNSInstancePermittedNetworks gets the permitted VPC networks for a DNS Services instance
func (c *Client) GetDNSInstancePermittedNetworks(ctx context.Context, dnsID string, dnsZone string) ([]string, error) {
	_, cancel := context.WithTimeout(ctx, 1*time.Minute)
	defer cancel()

	dnsService, err := c.getDNSServicesAPI()
	if err != nil {
		return nil, fmt.Errorf("failed to create dns services service to retrieve permitted networks: %w", err)
	}

	listPermittedNetworksOptions := dnsService.NewListPermittedNetworksOptions(dnsID, dnsZone)
	permittedNetworks, _, err := dnsService.ListPermittedNetworksWithContext(ctx, listPermittedNetworksOptions)
	if err != nil {
		return nil, err
	}

	networks := []string{}
	for _, network := range permittedNetworks.PermittedNetworks {
		networks = append(networks, *network.PermittedNetwork.VpcCrn)
	}
	return networks, nil
}

// GetDedicatedHostByName gets dedicated host by name.
func (c *Client) GetDedicatedHostByName(ctx context.Context, name string, region string) (*vpcv1.DedicatedHost, error) {
	err := c.SetVPCServiceURLForRegion(ctx, region)
	if err != nil {
		return nil, err
	}

	options := c.vpcAPI.NewListDedicatedHostsOptions()
	dhosts, _, err := c.vpcAPI.ListDedicatedHostsWithContext(ctx, options)
	if err != nil {
		return nil, errors.Wrap(err, "failed to list dedicated hosts")
	}

	for _, dhost := range dhosts.DedicatedHosts {
		if *dhost.Name == name {
			return &dhost, nil
		}
	}

	return nil, fmt.Errorf("dedicated host %q not found", name)
}

// GetDedicatedHostProfiles gets a list of profiles supported in a region.
func (c *Client) GetDedicatedHostProfiles(ctx context.Context, region string) ([]vpcv1.DedicatedHostProfile, error) {
	err := c.SetVPCServiceURLForRegion(ctx, region)
	if err != nil {
		return nil, err
	}

	profilesOptions := c.vpcAPI.NewListDedicatedHostProfilesOptions()
	profiles, _, err := c.vpcAPI.ListDedicatedHostProfilesWithContext(ctx, profilesOptions)
	if err != nil {
		return nil, err
	}

	return profiles.Profiles, nil
}

// GetDNSRecordsByName gets DNS records in specific Cloud Internet Services instance
// by its CRN, zone ID, and DNS record name.
func (c *Client) GetDNSRecordsByName(ctx context.Context, crnstr string, zoneID string, recordName string) ([]dnsrecordsv1.DnsrecordDetails, error) {
	dnsRecordsService, err := c.getDNSRecordsAPI(crnstr, zoneID)
	if err != nil {
		return nil, fmt.Errorf("failed to create dns records service to retrieve dns record: %w", err)
	}

	// Get CIS DNS records by name
	records, _, err := dnsRecordsService.ListAllDnsRecordsWithContext(ctx, &dnsrecordsv1.ListAllDnsRecordsOptions{
		Name: core.StringPtr(recordName),
	})
	if err != nil {
		return nil, errors.Wrap(err, "could not retrieve DNS records")
	}

	return records.Result, nil
}

// GetDNSZoneIDByName gets the DNS (Internal) or CIS zone ID from its domain name.
func (c *Client) GetDNSZoneIDByName(ctx context.Context, name string, publish types.PublishingStrategy) (string, error) {
	zones, err := c.GetDNSZones(ctx, publish)
	if err != nil {
		return "", err
	}

	for _, z := range zones {
		if z.Name == name {
			return z.ID, nil
		}
	}

	return "", fmt.Errorf("DNS zone %q not found", name)
}

// GetDNSZones returns all of the active DNS zones managed by DNS or CIS.
func (c *Client) GetDNSZones(ctx context.Context, publish types.PublishingStrategy) ([]responses.DNSZoneResponse, error) {
	_, cancel := context.WithTimeout(ctx, 1*time.Minute)
	defer cancel()

	if publish == types.InternalPublishingStrategy {
		return c.getDNSDNSZones(ctx)
	}
	return c.getCISDNSZones(ctx)
}

func (c *Client) getDNSDNSZones(ctx context.Context) ([]responses.DNSZoneResponse, error) {
	options := c.controllerAPI.NewListResourceInstancesOptions()
	options.SetResourceID(dnsServiceID)

	listResourceInstancesResponse, _, err := c.controllerAPI.ListResourceInstances(options)
	if err != nil {
		return nil, errors.Wrap(err, "failed to get dns instance")
	}

	var allZones []responses.DNSZoneResponse
	for _, instance := range listResourceInstancesResponse.Resources {
		dnsZoneService, err := c.getDNSZonesAPI()
		if err != nil {
			return nil, fmt.Errorf("failed to create dns zones service to retrieve dns zone: %w", err)
		}

		listZonesOptions := dnsZoneService.NewListDnszonesOptions(*instance.GUID)
		result, _, err := dnsZoneService.ListDnszones(listZonesOptions)
		if result == nil {
			return nil, err
		}

		for _, zone := range result.Dnszones {
			stateLower := strings.ToLower(*zone.State)
			// DNS Zones can be 'pending_network_add' (without a permitted network, added during TF)
			if stateLower == dnszonesv1.Dnszone_State_Active || stateLower == dnszonesv1.Dnszone_State_PendingNetworkAdd {
				zoneStruct := responses.DNSZoneResponse{
					Name:            *zone.Name,
					ID:              *zone.ID,
					InstanceID:      *instance.GUID,
					InstanceCRN:     *instance.CRN,
					InstanceName:    *instance.Name,
					ResourceGroupID: *instance.ResourceGroupID,
				}
				allZones = append(allZones, zoneStruct)
			}
		}
	}

	return allZones, nil
}

func (c *Client) getCISDNSZones(ctx context.Context) ([]responses.DNSZoneResponse, error) {
	options := c.controllerAPI.NewListResourceInstancesOptions()
	options.SetResourceID(cisServiceID)

	listResourceInstancesResponse, _, err := c.controllerAPI.ListResourceInstances(options)
	if err != nil {
		return nil, errors.Wrap(err, "failed to get cis instance")
	}

	var allZones []responses.DNSZoneResponse
	for _, instance := range listResourceInstancesResponse.Resources {
		zonesService, err := c.getCISZonesAPI(instance.CRN)
		if err != nil {
			return nil, fmt.Errorf("failed to create zones service to retrieve dns zone: %w", err)
		}

		listZonesOptions := zonesService.NewListZonesOptions()
		listZonesResponse, _, err := zonesService.ListZones(listZonesOptions)

		if listZonesResponse == nil {
			return nil, err
		}

		for _, zone := range listZonesResponse.Result {
			if *zone.Status == "active" {
				zoneStruct := responses.DNSZoneResponse{
					Name:            *zone.Name,
					ID:              *zone.ID,
					InstanceID:      *instance.GUID,
					InstanceCRN:     *instance.CRN,
					InstanceName:    *instance.Name,
					ResourceGroupID: *instance.ResourceGroupID,
				}
				allZones = append(allZones, zoneStruct)
			}
		}
	}

	return allZones, nil
}

// GetEncryptionKey gets data for an encryption key
func (c *Client) GetEncryptionKey(ctx context.Context, keyCRN string) (*responses.EncryptionKeyResponse, error) {
	// Parse out the IBM Cloud details from CRN
	crn, err := ibmcrn.Parse(keyCRN)
	if err != nil {
		return nil, err
	}

	var keyResponse *responses.EncryptionKeyResponse

	keyAPI, err := c.getKeyServiceAPI(crn)
	if err != nil {
		return nil, err
	}

	key, err := keyAPI.GetKey(ctx, crn.Resource)
	if err != nil {
		return nil, err
	}

	if key != nil {
		keyResponse = &responses.EncryptionKeyResponse{
			ID:      key.ID,
			Type:    key.Type,
			CRN:     key.CRN,
			State:   key.State,
			Deleted: key.Deleted,
		}
	}

	return keyResponse, nil
}

// GetIBMCloudRegions gets the Regions for IBM Cloud, mapped by shortname to descriptive name.
func (c *Client) GetIBMCloudRegions(ctx context.Context) (map[string]string, error) {
	regions := make(map[string]string, 0)
	// Global Catalog is used to collect the IBM Cloud Regions, via multiple commands.
	geographyOptions := c.catalogAPI.NewListCatalogEntriesOptions()
	geographyOptions.SetQ("kind:geography")

	geographyResult, _, err := c.catalogAPI.ListCatalogEntriesWithContext(ctx, geographyOptions)
	if err != nil {
		return nil, fmt.Errorf("failed to list ibmcloud georaphic regions: %w", err)
	}

	countryList, err := c.GetChildrenFromParents(ctx, geographyResult.Resources, "country")
	if err != nil {
		return nil, fmt.Errorf("failed to list ibmcloud countrys: %w", err)
	}
	metroList, err := c.GetChildrenFromParents(ctx, countryList, "metro")
	if err != nil {
		return nil, fmt.Errorf("failed to list ibmcloud metros: %w", err)
	}

	regionList, err := c.GetChildrenFromParents(ctx, metroList, "region")
	if err != nil {
		return nil, fmt.Errorf("failed to list ibmcloud regions: %w", err)
	}

	for _, child := range regionList {
		if child.ID == nil {
			continue
		}
		var description string
		// Attempt to collect the descriptive Region name, although this is based on language, so leaves room for improvement.
		if _, ok := child.OverviewUI["en"]; ok {
			description = *child.OverviewUI["en"].Description
		}
		regions[*child.ID] = description
	}

	return regions, nil
}

// GetChildrenFromParents fetches the children from the IBM Catalog using the given list of parents and the specified kind.
func (c *Client) GetChildrenFromParents(ctx context.Context, parentList []globalcatalogv1.CatalogEntry, kind string) ([]globalcatalogv1.CatalogEntry, error) {
	var childrenList []globalcatalogv1.CatalogEntry

	for num, parent := range parentList {
		if parent.ID == nil {
			fmt.Printf("skipping parent num %d of type %s due to unexpected nil ID\n", num, kind)
			continue
		}

		childrenOptions := c.catalogAPI.NewGetChildObjectsOptions(*parent.ID, kind)
		childrenOptions.SetInclude("*")
		childrenResult, _, err := c.catalogAPI.GetChildObjectsWithContext(ctx, childrenOptions)
		if err != nil {
			return childrenList, fmt.Errorf("failed to retrieve child entries of type %s for %s: %w", kind, *parent.ID, err)
		}
		childrenList = append(childrenList, childrenResult.Resources...)
	}
	return childrenList, nil
}

// GetLoadBalancer gets a VPC Load Balancer by ID.
func (c *Client) GetLoadBalancer(ctx context.Context, loadBalancerID string) (*vpcv1.LoadBalancer, error) {
	localContext, cancel := context.WithTimeout(ctx, 1*time.Minute)
	defer cancel()

	options := c.vpcAPI.NewGetLoadBalancerOptions(loadBalancerID)
	loadBalancer, _, err := c.vpcAPI.GetLoadBalancerWithContext(localContext, options)
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve load balancer: %w", err)
	}
	return loadBalancer, nil
}

// GetResourceGroup gets a resource group by its name or ID.
func (c *Client) GetResourceGroup(ctx context.Context, nameOrID string) (*resourcemanagerv2.ResourceGroup, error) {
	groups, err := c.GetResourceGroups(ctx)
	if err != nil {
		return nil, err
	}

	for idx, rg := range groups {
		if *rg.ID == nameOrID || *rg.Name == nameOrID {
			return &groups[idx], nil
		}
	}
	return nil, fmt.Errorf("resource group %q not found", nameOrID)
}

// GetResourceGroups gets the list of resource groups.
func (c *Client) GetResourceGroups(ctx context.Context) ([]resourcemanagerv2.ResourceGroup, error) {
	localContext, cancel := context.WithTimeout(ctx, 1*time.Minute)
	defer cancel()

	apikey, err := c.GetAuthenticatorAPIKeyDetails(ctx)
	if err != nil {
		return nil, err
	}

	options := c.managementAPI.NewListResourceGroupsOptions()
	options.SetAccountID(*apikey.AccountID)

	listResourceGroupsResponse, _, err := c.managementAPI.ListResourceGroupsWithContext(localContext, options)
	if err != nil {
		return nil, err
	}
	return listResourceGroupsResponse.Resources, nil
}

// GetSSHKeyByPublicKey gets an SSH Key by its public key.
func (c *Client) GetSSHKeyByPublicKey(ctx context.Context, publicKey string) (*vpcv1.Key, error) {
	localContext, cancel := context.WithTimeout(ctx, 1*time.Minute)
	defer cancel()

	keysPager, err := c.vpcAPI.NewKeysPager(&vpcv1.ListKeysOptions{})
	if err != nil {
		return nil, fmt.Errorf("failed to create keys pager to list keys: %w", err)
	}

	keys, err := keysPager.GetAllWithContext(localContext)
	if err != nil {
		return nil, fmt.Errorf("failed to list all keys from key pager: %w", err)
	}

	for _, k := range keys {
		if k.PublicKey != nil && strings.TrimSpace(*k.PublicKey) == strings.TrimSpace(publicKey) {
			return ptr.To(k), nil
		}
	}

	return nil, nil
}

// GetSubnet gets a subnet by its ID.
func (c *Client) GetSubnet(ctx context.Context, subnetID string) (*vpcv1.Subnet, error) {
	_, cancel := context.WithTimeout(ctx, 1*time.Minute)
	defer cancel()

	subnet, detailedResponse, err := c.vpcAPI.GetSubnet(&vpcv1.GetSubnetOptions{ID: &subnetID})
	if detailedResponse.GetStatusCode() == http.StatusNotFound {
		return nil, &VPCResourceNotFoundError{}
	}
	return subnet, err
}

// GetSubnetByName gets a subnet by its Name.
func (c *Client) GetSubnetByName(ctx context.Context, subnetName string, region string) (*vpcv1.Subnet, error) {
	localContext, cancel := context.WithTimeout(ctx, 1*time.Minute)
	defer cancel()

	err := c.SetVPCServiceURLForRegion(ctx, region)
	if err != nil {
		return nil, err
	}

	listSubnetsOptions := c.vpcAPI.NewListSubnetsOptions()
	subnetsPager, err := c.vpcAPI.NewSubnetsPager(listSubnetsOptions)
	if err != nil {
		return nil, err
	}

	subnets, err := subnetsPager.GetAllWithContext(localContext)
	if err != nil {
		return nil, err
	}
	for _, subnet := range subnets {
		if subnetName == *subnet.Name {
			return &subnet, nil
		}
	}
	return nil, &VPCResourceNotFoundError{}
}

// GetVSIProfiles gets a list of all VSI profiles.
func (c *Client) GetVSIProfiles(ctx context.Context) ([]vpcv1.InstanceProfile, error) {
	listInstanceProfilesOptions := c.vpcAPI.NewListInstanceProfilesOptions()
	profiles, _, err := c.vpcAPI.ListInstanceProfilesWithContext(ctx, listInstanceProfilesOptions)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to list vpc vsi profiles using: %s", c.vpcAPI.Service.Options.URL)
	}
	return profiles.Profiles, nil
}

// GetVPC gets a VPC by its ID.
func (c *Client) GetVPC(ctx context.Context, vpcID string) (*vpcv1.VPC, error) {
	_, cancel := context.WithTimeout(ctx, 1*time.Minute)
	defer cancel()

	regions, err := c.getVPCRegions(ctx)
	if err != nil {
		return nil, err
	}

	for _, region := range regions {
		err := c.vpcAPI.SetServiceURL(fmt.Sprintf("%s/v1", *region.Endpoint))
		if err != nil {
			return nil, errors.Wrap(err, "failed to set vpc api service url")
		}

		if vpc, detailedResponse, err := c.vpcAPI.GetVPC(c.vpcAPI.NewGetVPCOptions(vpcID)); err != nil {
			if detailedResponse != nil {
				// If the response code signifies the VPC was not found, simply move on to the next region; otherwise we log the response
				if detailedResponse.GetStatusCode() != http.StatusNotFound {
					logrus.Warnf("Unexpected response while checking VPC %s in %s region: %s", vpcID, *region.Name, detailedResponse)
				}
			} else {
				logrus.Warnf("Failure collecting VPC %s in %s: %q", vpcID, *region.Name, err)
			}
		} else if vpc != nil {
			return vpc, nil
		}
	}

	return nil, &VPCResourceNotFoundError{}
}

// GetVPCs gets all VPCs in a region
func (c *Client) GetVPCs(ctx context.Context, region string) ([]vpcv1.VPC, error) {
	_, cancel := context.WithTimeout(ctx, 1*time.Minute)
	defer cancel()

	err := c.SetVPCServiceURLForRegion(ctx, region)
	if err != nil {
		return nil, errors.Wrap(err, "failed to set vpc api service url")
	}

	allVPCs := []vpcv1.VPC{}
	if vpcs, detailedResponse, err := c.vpcAPI.ListVpcs(c.vpcAPI.NewListVpcsOptions()); err != nil {
		if detailedResponse.GetStatusCode() != http.StatusNotFound {
			return nil, err
		}
	} else if vpcs != nil {
		allVPCs = append(allVPCs, vpcs.Vpcs...)
	}
	return allVPCs, nil
}

// GetVPCByName gets a VPC by its name.
func (c *Client) GetVPCByName(ctx context.Context, vpcName string) (*vpcv1.VPC, error) {
	_, cancel := context.WithTimeout(ctx, 1*time.Minute)
	defer cancel()

	regions, err := c.getVPCRegions(ctx)
	if err != nil {
		return nil, err
	}

	for _, region := range regions {
		err := c.vpcAPI.SetServiceURL(fmt.Sprintf("%s/v1", *region.Endpoint))
		if err != nil {
			return nil, fmt.Errorf("failed to set vpc api service url: %w", err)
		}

		if vpcs, detailedResponse, err := c.vpcAPI.ListVpcsWithContext(ctx, c.vpcAPI.NewListVpcsOptions()); err != nil {
			if detailedResponse != nil {
				// If the response code signifies no VPCs were not found, we simply move on to the next region; otherwise log the response
				if detailedResponse.GetStatusCode() != http.StatusNotFound {
					logrus.Warnf("Unexpected response while checking %s region: %s", *region.Name, detailedResponse)
				}
			} else {
				logrus.Warnf("Failure collecting VPCs in %s: %q", *region.Name, err)
			}
		} else {
			for _, vpc := range vpcs.Vpcs {
				if *vpc.Name == vpcName {
					return &vpc, nil
				}
			}
		}
	}

	return nil, &VPCResourceNotFoundError{}
}

// GetVPCZonesForRegion gets the supported zones for a VPC region.
func (c *Client) GetVPCZonesForRegion(ctx context.Context, region string) ([]string, error) {
	_, cancel := context.WithTimeout(ctx, 1*time.Minute)
	defer cancel()

	regionZonesOptions := c.vpcAPI.NewListRegionZonesOptions(region)
	zones, _, err := c.vpcAPI.ListRegionZonesWithContext(ctx, regionZonesOptions)
	if err != nil {
		return nil, err
	}

	response := make([]string, len(zones.Zones))
	for idx, zone := range zones.Zones {
		response[idx] = *zone.Name
	}
	return response, err
}

func (c *Client) getVPCRegions(ctx context.Context) ([]vpcv1.Region, error) {
	listRegionsOptions := c.vpcAPI.NewListRegionsOptions()
	listRegionsResponse, _, err := c.vpcAPI.ListRegionsWithContext(ctx, listRegionsOptions)
	if err != nil {
		return nil, errors.Wrap(err, "failed to list vpc regions")
	}

	return listRegionsResponse.Regions, nil
}

func (c *Client) loadGlobalCatalogAPI() error {
	authenticator, err := NewIamAuthenticator(c.GetAPIKey(), c.iamServiceEndpointOverride)
	if err != nil {
		return fmt.Errorf("failed generating authenticator for global catalog api: %w", err)
	}
	options := &globalcatalogv1.GlobalCatalogV1Options{
		Authenticator: authenticator,
	}
	if overrideURL := ibmcloudtypes.CheckServiceEndpointOverride(configv1.IBMCloudServiceGlobalCatalog, c.serviceEndpoints); overrideURL != "" {
		options.URL = overrideURL
	}
	globalCatalogV1Service, err := globalcatalogv1.NewGlobalCatalogV1(options)
	if err != nil {
		return fmt.Errorf("failed creating global catalog api service: %w", err)
	}
	c.catalogAPI = globalCatalogV1Service
	return nil
}

func (c *Client) loadResourceManagementAPI() error {
	authenticator, err := NewIamAuthenticator(c.GetAPIKey(), c.iamServiceEndpointOverride)
	if err != nil {
		return err
	}
	options := &resourcemanagerv2.ResourceManagerV2Options{
		Authenticator: authenticator,
	}
	if overrideURL := ibmcloudtypes.CheckServiceEndpointOverride(configv1.IBMCloudServiceResourceManager, c.serviceEndpoints); overrideURL != "" {
		options.URL = overrideURL
	}
	resourceManagerV2Service, err := resourcemanagerv2.NewResourceManagerV2(options)
	if err != nil {
		return err
	}
	c.managementAPI = resourceManagerV2Service
	return nil
}

func (c *Client) loadResourceControllerAPI() error {
	authenticator, err := NewIamAuthenticator(c.GetAPIKey(), c.iamServiceEndpointOverride)
	if err != nil {
		return err
	}
	options := &resourcecontrollerv2.ResourceControllerV2Options{
		Authenticator: authenticator,
	}
	if overrideURL := ibmcloudtypes.CheckServiceEndpointOverride(configv1.IBMCloudServiceResourceController, c.serviceEndpoints); overrideURL != "" {
		options.URL = overrideURL
	}
	resourceControllerV2Service, err := resourcecontrollerv2.NewResourceControllerV2(options)
	if err != nil {
		return err
	}
	c.controllerAPI = resourceControllerV2Service
	return nil
}

func (c *Client) loadVPCV1API() error {
	authenticator, err := NewIamAuthenticator(c.GetAPIKey(), c.iamServiceEndpointOverride)
	if err != nil {
		return err
	}
	options := &vpcv1.VpcV1Options{
		Authenticator: authenticator,
	}
	if overrideURL := ibmcloudtypes.CheckServiceEndpointOverride(configv1.IBMCloudServiceVPC, c.serviceEndpoints); overrideURL != "" {
		options.URL = overrideURL
	}
	vpcService, err := vpcv1.NewVpcV1(options)
	if err != nil {
		return err
	}
	c.vpcAPI = vpcService
	return nil
}

// getCISZonesAPI returns a new Zones service. This function isolates the management of Zones service, as it only required for Public (External) clusters.
// Zones is used to manage zones within IBM Cloud Internet Services (CIS) and is reliant on the CIS service endpoint.
func (c *Client) getCISZonesAPI(instanceCRN *string) (*zonesv1.ZonesV1, error) {
	authenticator, err := NewIamAuthenticator(c.GetAPIKey(), c.iamServiceEndpointOverride)
	if err != nil {
		return nil, fmt.Errorf("failed to create iam authenticator for zones service: %w", err)
	}
	options := &zonesv1.ZonesV1Options{
		Authenticator: authenticator,
		Crn:           instanceCRN,
	}

	// If a CIS service endpoint override was provided, pass it along to override the default Zones service, as zonesv1 is provided via IBM Cloud CIS endpoint.
	if overrideURL := ibmcloudtypes.CheckServiceEndpointOverride(configv1.IBMCloudServiceCIS, c.serviceEndpoints); overrideURL != "" {
		options.URL = overrideURL
	}

	zonesService, err := zonesv1.NewZonesV1(options)
	if err != nil {
		return nil, fmt.Errorf("failed to create zones service: %w", err)
	}
	return zonesService, nil
}

// getDNSZonesAPI returns a new DNS Zones service. This function isolates the management of DNS Zones service, as it is only required for Private (Internal) clusters currently.
// DNS Zones is used to manage zones within IBM Cloud DNS Services and is reliant on the DNS Services service endpoint.
func (c *Client) getDNSZonesAPI() (*dnszonesv1.DnsZonesV1, error) {
	authenticator, err := NewIamAuthenticator(c.GetAPIKey(), c.iamServiceEndpointOverride)
	if err != nil {
		return nil, fmt.Errorf("failed to create iam authenticator for dns zones service: %w", err)
	}
	options := &dnszonesv1.DnsZonesV1Options{
		Authenticator: authenticator,
	}

	// If a DNS Services service endpoint override was provided, pass it along to override the default DNS Zones service, as dnszonesv1 is provided via IBM Cloud DNS Services endpoint.
	if overrideURL := ibmcloudtypes.CheckServiceEndpointOverride(configv1.IBMCloudServiceDNSServices, c.serviceEndpoints); overrideURL != "" {
		options.URL = overrideURL
	}

	dnsZoneService, err := dnszonesv1.NewDnsZonesV1(options)
	if err != nil {
		return nil, fmt.Errorf("failed to create dns zones service: %w", err)
	}
	return dnsZoneService, nil
}

// getDNSServicesAPI returns a new DNS Services service. This function isolates the management of DNS Services service, as it is only required for Private (Internal) clusters currently.
func (c *Client) getDNSServicesAPI() (*dnssvcsv1.DnsSvcsV1, error) {
	authenticator, err := NewIamAuthenticator(c.GetAPIKey(), c.iamServiceEndpointOverride)
	if err != nil {
		return nil, err
	}
	options := &dnssvcsv1.DnsSvcsV1Options{
		Authenticator: authenticator,
	}

	// If a DNS Services service endpoint override was provided, pass it along to override the default DNS Services service endpoint.
	if overrideURL := ibmcloudtypes.CheckServiceEndpointOverride(configv1.IBMCloudServiceDNSServices, c.serviceEndpoints); overrideURL != "" {
		options.URL = overrideURL
	}

	dnsService, err := dnssvcsv1.NewDnsSvcsV1(options)
	if err != nil {
		return nil, err
	}
	return dnsService, nil
}

// getDNSRecordsAPI returns a new DNS Records service. This function isolates the management of DNS Records service, as it is only required for Public (External) clusters.
// DNS Records is used to manage DNS records within IBM Cloud Internet Services (CIS) and is reliant on the CIS service endpoint.
func (c *Client) getDNSRecordsAPI(instanceCRN string, zoneID string) (*dnsrecordsv1.DnsRecordsV1, error) {
	authenticator, err := NewIamAuthenticator(c.GetAPIKey(), c.iamServiceEndpointOverride)
	if err != nil {
		return nil, fmt.Errorf("failed to create iam authenticator for dns records service: %w", err)
	}
	// Set CIS DNS record service options
	options := &dnsrecordsv1.DnsRecordsV1Options{
		Authenticator:  authenticator,
		Crn:            core.StringPtr(instanceCRN),
		ZoneIdentifier: core.StringPtr(zoneID),
	}

	// If a CIS service endpoint override was provided, pass it along to override the default DNS Records service, as dnsrecordsv1 is provided via IBM CIS endpoint.
	if overrideURL := ibmcloudtypes.CheckServiceEndpointOverride(configv1.IBMCloudServiceCIS, c.serviceEndpoints); overrideURL != "" {
		options.URL = overrideURL
	}

	dnsRecordsService, err := dnsrecordsv1.NewDnsRecordsV1(options)
	if err != nil {
		return nil, fmt.Errorf("failed to create dns records service: %w", err)
	}
	return dnsRecordsService, nil
}

func (c *Client) getKeyServiceAPI(crn ibmcrn.CRN) (*kpclient.Client, error) {
	var clientConfig kpclient.ClientConfig

	switch crn.ServiceName {
	case hyperProtectCRNServiceName:
		clientConfig = kpclient.ClientConfig{
			BaseURL:    fmt.Sprintf(hyperProtectDefaultURLTemplate, crn.Region),
			APIKey:     c.apiKey,
			TokenURL:   iamTokenDefaultURL,
			InstanceID: crn.ServiceInstance,
		}

		// Override HyperProtect service URL, if one was provided
		if overrideURL := ibmcloudtypes.CheckServiceEndpointOverride(configv1.IBMCloudServiceHyperProtect, c.serviceEndpoints); overrideURL != "" {
			clientConfig.BaseURL = overrideURL
		}
	case keyProtectCRNServiceName:
		clientConfig = kpclient.ClientConfig{
			BaseURL:    fmt.Sprintf(keyProtectDefaultURLTemplate, crn.Region),
			APIKey:     c.apiKey,
			TokenURL:   iamTokenDefaultURL,
			InstanceID: crn.ServiceInstance,
		}

		// Override KeyProtect service URL, if one was provided
		if overrideURL := ibmcloudtypes.CheckServiceEndpointOverride(configv1.IBMCloudServiceKeyProtect, c.serviceEndpoints); overrideURL != "" {
			clientConfig.BaseURL = overrideURL
		}
	default:
		return nil, fmt.Errorf("unknown key service for provided encryption key: %s", crn)
	}

	// Override IAM token URL, if an IAM service override URL was provided
	if overrideURL := ibmcloudtypes.CheckServiceEndpointOverride(configv1.IBMCloudServiceIAM, c.serviceEndpoints); overrideURL != "" {
		// Construct the token URL using the overridden IAM URL and the token path
		clientConfig.TokenURL = fmt.Sprintf("%s/%s", overrideURL, iamTokenPath)
	}

	return kpclient.New(clientConfig, kpclient.DefaultTransport())
}

// SetVPCServiceURLForRegion will set the VPC Service URL to a specific IBM Cloud Region, in order to access Region scoped resources
func (c *Client) SetVPCServiceURLForRegion(ctx context.Context, region string) error {
	regionOptions := c.vpcAPI.NewGetRegionOptions(region)
	vpcRegion, _, err := c.vpcAPI.GetRegionWithContext(ctx, regionOptions)
	if err != nil {
		return err
	}
	err = c.vpcAPI.SetServiceURL(fmt.Sprintf("%s/v1", *vpcRegion.Endpoint))
	if err != nil {
		return err
	}
	return nil
}
