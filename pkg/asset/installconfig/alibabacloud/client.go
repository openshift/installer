package alibabacloud

import (
	"os"

	survey "github.com/AlecAivazis/survey/v2"
	"github.com/aliyun/alibaba-cloud-sdk-go/sdk"
	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/auth"
	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/auth/credentials"
	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/auth/credentials/provider"
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
	AccessKeyId     string
	AccessKeySecret string
}

func NewClientWithOptions(regionId string, config *sdk.Config, credential auth.Credential) (client *Client, err error) {
	client = &Client{}
	err = client.InitWithOptions(regionId, config, credential)
	return
}

func NewClient(regionId string) (client *Client, err error) {
	credential, err := getCredentials()
	if err != nil {
		return nil, err
	}

	config := sdk.NewConfig()

	client, err = NewClientWithOptions(regionId, config, credential)

	if _credential, ok := credential.(credentials.AccessKeyCredential); ok {
		client.AccessKeyId = _credential.AccessKeyId
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
	var access_key_id string

	err := survey.Ask([]*survey.Question{
		{
			Prompt: &survey.Input{
				Message: "Alibaba Cloud Access Key ID",
				Help:    "The AccessKey ID is used to identify a user.\nhttps://www.alibabacloud.com/help/doc-detail/53045.html",
			},
		},
	}, &access_key_id)

	if err != nil {
		return nil, err
	}

	var access_key_secret string
	err = survey.Ask([]*survey.Question{
		{
			Prompt: &survey.Password{
				Message: "Alibaba Cloud Secret Access Key",
				Help:    "The AccessKey secret is used to verify a user. You must keep your AccessKey secret strictly confidential.",
			},
		},
	}, &access_key_secret)

	if err != nil {
		return nil, err
	}

	os.Setenv(ENVAccessKeyID, access_key_id)
	os.Setenv(ENVAccessKeySecret, access_key_secret)

	return credentials.NewAccessKeyCredential(access_key_id, access_key_secret), nil
}

func (client *Client) DescribeRegions(region_id string) (response *ecs.DescribeRegionsResponse, err error) {
	request := ecs.CreateDescribeRegionsRequest()
	request.AcceptLanguage = DefaultAcceptLanguage
	request.RegionId = region_id
	response = &ecs.DescribeRegionsResponse{
		BaseResponse: &responses.BaseResponse{},
	}
	err = client.DoAction(request, response)
	return
}

func (client *Client) ListResourceGroups() (response *resourcemanager.ListResourceGroupsResponse, err error) {
	request := resourcemanager.CreateListResourceGroupsRequest()
	request.Status = "OK"
	request.Scheme = "https"
	response = &resourcemanager.ListResourceGroupsResponse{
		BaseResponse: &responses.BaseResponse{},
	}
	err = client.DoAction(request, response)
	return
}

func (client *Client) ListPrivateZoneRegions() (response *pvtz.DescribeRegionsResponse, err error) {
	request := pvtz.CreateDescribeRegionsRequest()
	request.AcceptLanguage = DefaultAcceptLanguage
	response = &pvtz.DescribeRegionsResponse{
		BaseResponse: &responses.BaseResponse{},
	}
	err = client.DoAction(request, response)
	return
}

func (client *Client) ListPrivateZones(zone_name string) (response *pvtz.DescribeZonesResponse, err error) {
	request := pvtz.CreateDescribeZonesRequest()
	request.Lang = "en"
	request.Keyword = zone_name
	response = &pvtz.DescribeZonesResponse{
		BaseResponse: &responses.BaseResponse{},
	}
	err = client.DoAction(request, response)
	return
}
