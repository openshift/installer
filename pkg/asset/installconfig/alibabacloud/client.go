package alibabacloud

import (
	"os"
	"strings"

	survey "github.com/AlecAivazis/survey/v2"
	"github.com/aliyun/alibaba-cloud-sdk-go/sdk"
	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/auth"
	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/auth/credentials"
	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/auth/credentials/provider"
	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/endpoints"
	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/requests"
	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/responses"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/ecs"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/pvtz"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/resourcemanager"
)

// Environmental virables
const (
	ENVAccessKeyID     = "ALIBABA_CLOUD_ACCESS_KEY_ID"
	ENVAccessKeySecret = "ALIBABA_CLOUD_ACCESS_KEY_SECRET"
)

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

	if _credential, ok := credential.(credentials.AccessKeyCredential); ok {
		client.AccessKeyID = _credential.AccessKeyId
		client.AccessKeySecret = _credential.AccessKeySecret
	}
	return
}

func getCredentials() (credential auth.Credential, err error) {
	// Get AccessKey and AccessKeySecret information from the enviroment
	// usage: https://github.com/aliyun/alibaba-cloud-sdk-go/blob/7259de46d58ef905c66e04babf791190512a85da/docs/2-Client-EN.md#1-environment-credentials
	credential, err = provider.NewEnvProvider().Resolve()
	if credential != nil {
		return credential, nil
	}

	// Get AccessKey and AccessKeySecret information from configuration file,default path:"~/.alibabacloud/credentials"
	// usage: https://github.com/aliyun/alibaba-cloud-sdk-go/blob/7259de46d58ef905c66e04babf791190512a85da/docs/2-Client-EN.md#2-credentials-file
	credential, err = provider.NewProfileProvider().Resolve()
	if credential != nil {
		return credential, nil
	}

	// Get AccessKey and AccessKeySecret information via interactive command line
	credential, err = askCredentials()
	if credential != nil {
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

	os.Setenv(ENVAccessKeyID, accessKeyID)
	os.Setenv(ENVAccessKeySecret, accessKeySecret)

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

func defaultEndpoint() map[string]string {
	return map[string]string{
		"pvtz":            "pvtz.aliyuncs.com",
		"resourcemanager": "resourcemanager.aliyuncs.com",
		"ecs":             "ecs.aliyuncs.com",
	}
}
