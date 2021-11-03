package alibabacloud

import (
	"fmt"
	"os"
	"os/user"
	"path/filepath"
	"strings"

	survey "github.com/AlecAivazis/survey/v2"
	"github.com/aliyun/alibaba-cloud-sdk-go/sdk"
	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/auth"
	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/auth/credentials"
	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/auth/credentials/provider"
	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/endpoints"
	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/requests"
	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/responses"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/alidns"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/ecs"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/pvtz"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/resourcemanager"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/vpc"
	"github.com/aliyun/aliyun-oss-go-sdk/oss"
	"github.com/pkg/errors"
)

// Environment virables
const (
	envCredentialFile = "ALIBABA_CLOUD_CREDENTIALS_FILE"
)

// Credential configuration file template.
const configurationTemplate = `
[default]              
enable = true                    
type = access_key                
access_key_id = %s              
access_key_secret = %s
`

// Client makes calls to the Alibaba Cloud API.
type Client struct {
	sdk.Client
	RegionID        string
	AccessKeyID     string
	AccessKeySecret string
}

func newClientWithOptions(regionID string, config *sdk.Config, credential auth.Credential) (client *Client, err error) {
	client = &Client{
		RegionID: regionID,
	}
	err = client.InitWithOptions(regionID, config, credential)
	return
}

// NewClient initializes a client with a session.
func NewClient(regionID string) (client *Client, err error) {
	credential, err := getCredentials()
	if err != nil {
		return nil, err
	}

	config := sdk.NewConfig()

	client, err = newClientWithOptions(regionID, config, credential)

	switch credentialType := credential.(type) {
	case *credentials.AccessKeyCredential:
		{
			client.AccessKeyID = credentialType.AccessKeyId
			client.AccessKeySecret = credentialType.AccessKeySecret
		}
	default:
		errors.Errorf("Please use certification type AccessKey.")
	}

	return
}

func getCredentials() (credential auth.Credential, err error) {
	// Get AccessKey and AccessKeySecret information from the enviroment
	// usage: https://github.com/aliyun/alibaba-cloud-sdk-go/blob/7259de46d58ef905c66e04babf791190512a85da/docs/2-Client-EN.md#1-environment-credentials
	credential, err = provider.NewEnvProvider().Resolve()
	if err == nil && credential != nil {
		return credential, nil
	}

	// Get AccessKey and AccessKeySecret information from configuration file,default path:"~/.alibabacloud/credentials"
	// usage: https://github.com/aliyun/alibaba-cloud-sdk-go/blob/7259de46d58ef905c66e04babf791190512a85da/docs/2-Client-EN.md#2-credentials-file
	credential, err = provider.NewProfileProvider().Resolve()
	if err == nil && credential != nil {
		return credential, nil
	}

	// Get AccessKey and AccessKeySecret information via interactive command line
	credential, err = askCredentials()
	if err == nil && credential != nil {
		return credential, nil
	}

	return nil, err
}

func askCredentials() (auth.Credential, error) {
	var accessKeyID string

	err := survey.Ask([]*survey.Question{
		{
			Prompt: &survey.Input{
				Message: "Alibaba Cloud Access Key ID",
				Help:    "The AccessKey ID is used to identify a user.\nhttps://www.alibabacloud.com/help/doc-detail/53045.html",
			},
		},
	}, &accessKeyID)

	if err != nil {
		return nil, err
	}

	var accessKeySecret string
	err = survey.Ask([]*survey.Question{
		{
			Prompt: &survey.Password{
				Message: "Alibaba Cloud Secret Access Key",
				Help:    "The AccessKey secret is used to verify a user. You must keep your AccessKey secret strictly confidential.",
			},
		},
	}, &accessKeySecret)
	if err != nil {
		return nil, err
	}

	storeCredentials(accessKeyID, accessKeySecret)
	return credentials.NewAccessKeyCredential(accessKeyID, accessKeySecret), nil
}

func (client *Client) doActionWithSetDomain(request requests.AcsRequest, response responses.AcsResponse) (err error) {
	endpoint, err := endpoints.Resolve(&endpoints.ResolveParam{
		Product:  strings.ToLower(request.GetProduct()),
		RegionId: strings.ToLower(client.RegionID),
	})

	if err != nil {
		endpoint = defaultEndpoint()[strings.ToLower(request.GetProduct())]
	}

	request.SetDomain(endpoint)
	err = client.DoAction(request, response)
	return
}

// DescribeRegions gets the list of regions.
func (client *Client) DescribeRegions() (response *ecs.DescribeRegionsResponse, err error) {
	request := ecs.CreateDescribeRegionsRequest()
	request.AcceptLanguage = defaultAcceptLanguage
	response = &ecs.DescribeRegionsResponse{
		BaseResponse: &responses.BaseResponse{},
	}
	err = client.doActionWithSetDomain(request, response)
	return
}

// DescribeAvailableResource query available resources.
func (client *Client) DescribeAvailableResource(destinationResource string) (response *ecs.DescribeAvailableResourceResponse, err error) {
	request := ecs.CreateDescribeAvailableResourceRequest()
	request.RegionId = client.RegionID
	request.DestinationResource = destinationResource
	response = &ecs.DescribeAvailableResourceResponse{
		BaseResponse: &responses.BaseResponse{},
	}
	err = client.doActionWithSetDomain(request, response)
	return
}

// ListEnhanhcedNatGatewayAvailableZones query available zone for enhanhced NAT gateway.
func (client *Client) ListEnhanhcedNatGatewayAvailableZones() (response *vpc.ListEnhanhcedNatGatewayAvailableZonesResponse, err error) {
	request := vpc.CreateListEnhanhcedNatGatewayAvailableZonesRequest()
	request.RegionId = client.RegionID
	response = &vpc.ListEnhanhcedNatGatewayAvailableZonesResponse{
		BaseResponse: &responses.BaseResponse{},
	}
	err = client.doActionWithSetDomain(request, response)
	return
}

// DescribeAvailableInstanceType query available instance type of ECS.
func (client *Client) DescribeAvailableInstanceType(zoneID string, instanceType string) (response *ecs.DescribeAvailableResourceResponse, err error) {
	request := ecs.CreateDescribeAvailableResourceRequest()
	request.RegionId = client.RegionID
	request.ZoneId = zoneID
	request.DestinationResource = "InstanceType"
	request.InstanceType = instanceType
	response = &ecs.DescribeAvailableResourceResponse{
		BaseResponse: &responses.BaseResponse{},
	}
	err = client.doActionWithSetDomain(request, response)
	return
}

// ListResourceGroups gets the list of resource groups.
func (client *Client) ListResourceGroups() (response *resourcemanager.ListResourceGroupsResponse, err error) {
	request := resourcemanager.CreateListResourceGroupsRequest()
	request.Status = "OK"
	request.Scheme = "https"
	response = &resourcemanager.ListResourceGroupsResponse{
		BaseResponse: &responses.BaseResponse{},
	}
	err = client.doActionWithSetDomain(request, response)
	return
}

// ListPrivateZoneRegions gets the list of regions for privatzone.
func (client *Client) ListPrivateZoneRegions() (response *pvtz.DescribeRegionsResponse, err error) {
	request := pvtz.CreateDescribeRegionsRequest()
	request.AcceptLanguage = defaultAcceptLanguage
	response = &pvtz.DescribeRegionsResponse{
		BaseResponse: &responses.BaseResponse{},
	}
	err = client.doActionWithSetDomain(request, response)
	return
}

// ListDNSDomain get the list of domains.
func (client *Client) ListDNSDomain(baseDomain string) (response *alidns.DescribeDomainsResponse, err error) {
	request := alidns.CreateDescribeDomainsRequest()
	request.KeyWord = baseDomain
	response = &alidns.DescribeDomainsResponse{
		BaseResponse: &responses.BaseResponse{},
	}
	err = client.doActionWithSetDomain(request, response)
	return
}

// ListPrivateZones gets the list of privatzones.
func (client *Client) ListPrivateZones(zoneName string) (response *pvtz.DescribeZonesResponse, err error) {
	request := pvtz.CreateDescribeZonesRequest()
	request.Lang = "en"
	request.Keyword = zoneName
	response = &pvtz.DescribeZonesResponse{
		BaseResponse: &responses.BaseResponse{},
	}
	err = client.doActionWithSetDomain(request, response)
	return
}

// GetOSSObjectSignURL returns a presigned URL for a OSS object
func (client *Client) GetOSSObjectSignURL(bucketName string, objectName string) (signedURL string, err error) {
	endpoint := fmt.Sprintf("oss-%s.aliyuncs.com", client.RegionID)

	ossClient, err := oss.New(endpoint, client.AccessKeyID, client.AccessKeySecret)
	if err != nil {
		return "", err
	}

	bucket, err := ossClient.Bucket(bucketName)
	if err != nil {
		return "", err
	}

	signedURL, err = bucket.SignURL(objectName, oss.HTTPGet, 7200)
	return
}

func defaultEndpoint() map[string]string {
	return map[string]string{
		"pvtz":            "pvtz.aliyuncs.com",
		"resourcemanager": "resourcemanager.aliyuncs.com",
		"ecs":             "ecs.aliyuncs.com",
	}
}

func storeCredentials(accessKeyID string, accessKeySecret string) (err error) {
	dirPath, ok := os.LookupEnv(envCredentialFile)
	if !ok || dirPath == "" {
		user, err := user.Current()
		if err != nil {
			return err
		}
		dirPath = user.HomeDir
	}

	dirPath = filepath.Join(dirPath, ".alibabacloud")
	err = os.MkdirAll(dirPath, os.ModePerm)
	if err != nil {
		return err
	}

	filePath := filepath.Join(dirPath, "credentials")

	file, err := os.OpenFile(filePath, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0666)
	if err != nil {
		return err
	}
	defer file.Close()
	file.WriteString(fmt.Sprintf(configurationTemplate, accessKeyID, accessKeySecret))

	return nil
}
