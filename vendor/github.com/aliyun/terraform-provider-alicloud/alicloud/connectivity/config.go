package connectivity

import (
	"fmt"
	"log"

	roa "github.com/alibabacloud-go/tea-roa/client"

	"encoding/json"
	"net/http"
	"strings"

	rpc "github.com/alibabacloud-go/tea-rpc/client"
	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/auth"
	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/auth/credentials"
	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/requests"
	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/responses"
	credential "github.com/aliyun/credentials-go/credentials"
	"github.com/jmespath/go-jmespath"
)

var securityCredURL = "http://100.100.100.200/latest/meta-data/ram/security-credentials/"

// Config of aliyun
type Config struct {
	SourceIp             string
	AccessKey            string
	SecretKey            string
	EcsRoleName          string
	Region               Region
	RegionId             string
	SecurityToken        string
	OtsInstanceName      string
	AccountId            string
	Protocol             string
	ClientReadTimeout    int
	ClientConnectTimeout int

	RamRoleArn               string
	RamRoleSessionName       string
	RamRolePolicy            string
	RamRoleSessionExpiration int
	Endpoints                map[string]interface{}
	RKvstoreEndpoint         string
	EcsEndpoint              string
	RdsEndpoint              string
	SlbEndpoint              string
	VpcEndpoint              string
	CenEndpoint              string
	EssEndpoint              string
	OssEndpoint              string
	OnsEndpoint              string
	AlikafkaEndpoint         string
	DnsEndpoint              string
	RamEndpoint              string
	CsEndpoint               string
	CrEndpoint               string
	CdnEndpoint              string
	KmsEndpoint              string
	OtsEndpoint              string
	CmsEndpoint              string
	PvtzEndpoint             string
	StsEndpoint              string
	LogEndpoint              string
	DrdsEndpoint             string
	DdsEndpoint              string
	GpdbEnpoint              string
	KVStoreEndpoint          string
	PolarDBEndpoint          string
	FcEndpoint               string
	ApigatewayEndpoint       string
	DatahubEndpoint          string
	MnsEndpoint              string
	LocationEndpoint         string
	ElasticsearchEndpoint    string
	NasEndpoint              string
	BssOpenApiEndpoint       string
	DdoscooEndpoint          string
	DdosbgpEndpoint          string
	SagEndpoint              string
	EmrEndpoint              string
	CasEndpoint              string
	MarketEndpoint           string
	HBaseEndpoint            string
	AdbEndpoint              string
	MaxComputeEndpoint       string

	edasEndpoint                string
	SkipRegionValidation        bool
	ConfigurationSource         string
	CbnEndpoint                 string
	DmsEnterpriseEndpoint       string
	WafOpenapiEndpoint          string
	ResourcemanagerEndpoint     string
	BssopenapiEndpoint          string
	AlidnsEndpoint              string
	CassandraEndpoint           string
	EciEndpoint                 string
	OosEndpoint                 string
	DcdnEndpoint                string
	MseEndpoint                 string
	ActiontrailEndpoint         string
	ConfigEndpoint              string
	FnfEndpoint                 string
	RosEndpoint                 string
	PrivatelinkEndpoint         string
	MaxcomputeEndpoint          string
	ResourcesharingEndpoint     string
	GaEndpoint                  string
	HitsdbEndpoint              string
	BrainIndustrialEndpoint     string
	EipanycastEndpoint          string
	ImsEndpoint                 string
	QuotasEndpoint              string
	SgwEndpoint                 string
	ScdnEndpoint                string
	DmEndpoint                  string
	EventbridgeEndpoint         string
	OnsproxyEndpoint            string
	CdsEndpoint                 string
	HbrEndpoint                 string
	ArmsEndpoint                string
	CloudfwEndpoint             string
	ServerlessEndpoint          string
	AlbEndpoint                 string
	RedisaEndpoint              string
	GwsecdEndpoint              string
	CloudphoneEndpoint          string
	DataworkspublicEndpoint     string
	HcsSgwEndpoint              string
	CddcEndpoint                string
	MscopensubscriptionEndpoint string
	SddpEndpoint                string
	BastionhostEndpoint         string
}

func (c *Config) loadAndValidate() error {
	err := c.validateRegion()
	if err != nil {
		return err
	}

	return nil
}

func (c *Config) validateRegion() error {

	for _, valid := range ValidRegions {
		if c.Region == valid {
			return nil
		}
	}

	return fmt.Errorf("Invalid Alibaba Cloud region: %s", c.RegionId)
}

func (c *Config) getAuthCredential(stsSupported bool) auth.Credential {
	if c.AccessKey != "" && c.SecretKey != "" {
		if stsSupported && c.SecurityToken != "" {
			return credentials.NewStsTokenCredential(c.AccessKey, c.SecretKey, c.SecurityToken)
		}
		if c.RamRoleArn != "" {
			log.Printf("[INFO] Assume RAM Role specified in provider block assume_role { ... }")
			return credentials.NewRamRoleArnWithPolicyCredential(
				c.AccessKey, c.SecretKey, c.RamRoleArn,
				c.RamRoleSessionName, c.RamRolePolicy, c.RamRoleSessionExpiration)
		}
		return credentials.NewAccessKeyCredential(c.AccessKey, c.SecretKey)
	}
	if c.EcsRoleName != "" {
		return credentials.NewEcsRamRoleCredential(c.EcsRoleName)
	}

	return credentials.NewAccessKeyCredential(c.AccessKey, c.SecretKey)
}

// getAuthCredentialByEcsRoleName aims to access meta to get sts credential
// Actually, the job should be done by sdk, but currently not all resources and products support alibaba-cloud-sdk-go,
// and their go sdk does support ecs role name.
// This method is a temporary solution and it should be removed after all go sdk support ecs role name
// The related PR: https://github.com/aliyun/terraform-provider-alicloud/pull/731
func (c *Config) getAuthCredentialByEcsRoleName() (accessKey, secretKey, token string, err error) {
	if c.AccessKey != "" {
		return c.AccessKey, c.SecretKey, c.SecurityToken, nil
	}
	if c.EcsRoleName == "" {
		return
	}
	requestUrl := securityCredURL + c.EcsRoleName
	httpRequest, err := http.NewRequest(requests.GET, requestUrl, strings.NewReader(""))
	if err != nil {
		err = fmt.Errorf("build sts requests err: %s", err.Error())
		return
	}
	httpClient := &http.Client{}
	httpResponse, err := httpClient.Do(httpRequest)
	if err != nil {
		err = fmt.Errorf("get Ecs sts token err : %s", err.Error())
		return
	}

	response := responses.NewCommonResponse()
	err = responses.Unmarshal(response, httpResponse, "")
	if err != nil {
		err = fmt.Errorf("Unmarshal Ecs sts token response err : %s", err.Error())
		return
	}

	if response.GetHttpStatus() != http.StatusOK {
		err = fmt.Errorf("get Ecs sts token err, httpStatus: %d, message = %s", response.GetHttpStatus(), response.GetHttpContentString())
		return
	}
	var data interface{}
	err = json.Unmarshal(response.GetHttpContentBytes(), &data)
	if err != nil {
		err = fmt.Errorf("refresh Ecs sts token err, json.Unmarshal fail: %s", err.Error())
		return
	}
	code, err := jmespath.Search("Code", data)
	if err != nil {
		err = fmt.Errorf("refresh Ecs sts token err, fail to get Code: %s", err.Error())
		return
	}
	if code.(string) != "Success" {
		err = fmt.Errorf("refresh Ecs sts token err, Code is not Success")
		return
	}
	accessKeyId, err := jmespath.Search("AccessKeyId", data)
	if err != nil {
		err = fmt.Errorf("refresh Ecs sts token err, fail to get AccessKeyId: %s", err.Error())
		return
	}
	accessKeySecret, err := jmespath.Search("AccessKeySecret", data)
	if err != nil {
		err = fmt.Errorf("refresh Ecs sts token err, fail to get AccessKeySecret: %s", err.Error())
		return
	}
	securityToken, err := jmespath.Search("SecurityToken", data)
	if err != nil {
		err = fmt.Errorf("refresh Ecs sts token err, fail to get SecurityToken: %s", err.Error())
		return
	}

	if accessKeyId == nil || accessKeySecret == nil || securityToken == nil {
		err = fmt.Errorf("there is no any available accesskey, secret and security token for Ecs role %s", c.EcsRoleName)
		return
	}

	return accessKeyId.(string), accessKeySecret.(string), securityToken.(string), nil
}

func (c *Config) MakeConfigByEcsRoleName() error {
	accessKey, secretKey, token, err := c.getAuthCredentialByEcsRoleName()
	if err != nil {
		return err
	}
	c.AccessKey, c.SecretKey, c.SecurityToken = accessKey, secretKey, token
	return nil
}

func (c *Config) getTeaDslSdkConfig(stsSupported bool) (config rpc.Config, err error) {
	config.SetRegionId(c.RegionId)
	config.SetUserAgent(fmt.Sprintf("%s/%s %s/%s %s/%s", Terraform, terraformVersion, Provider, providerVersion, Module, c.ConfigurationSource))
	credential, err := credential.NewCredential(c.getCredentialConfig(stsSupported))
	config.SetCredential(credential).
		SetRegionId(c.RegionId).
		SetProtocol(c.Protocol).
		SetReadTimeout(c.ClientReadTimeout).
		SetConnectTimeout(c.ClientConnectTimeout).
		SetMaxIdleConns(500)
	return
}
func (c *Config) getTeaRoaDslSdkConfig(stsSupported bool) (config roa.Config, err error) {
	config.SetRegionId(c.RegionId)
	config.SetUserAgent(fmt.Sprintf("%s/%s %s/%s %s/%s", Terraform, terraformVersion, Provider, providerVersion, Module, c.ConfigurationSource))
	credential, err := credential.NewCredential(c.getCredentialConfig(stsSupported))
	config.SetCredential(credential).
		SetRegionId(c.RegionId).
		SetProtocol(c.Protocol).
		SetReadTimeout(c.ClientReadTimeout).
		SetConnectTimeout(c.ClientConnectTimeout).
		SetMaxIdleConns(500)
	return
}
func (c *Config) getCredentialConfig(stsSupported bool) *credential.Config {
	credentialType := ""
	credentialConfig := &credential.Config{}
	if c.AccessKey != "" && c.SecretKey != "" {
		credentialType = "access_key"
		credentialConfig.AccessKeyId = &c.AccessKey     // AccessKeyId
		credentialConfig.AccessKeySecret = &c.SecretKey // AccessKeySecret

		if stsSupported && c.SecurityToken != "" {
			credentialType = "sts"
			credentialConfig.SecurityToken = &c.SecurityToken // STS Token
		} else if c.RamRoleArn != "" {
			log.Printf("[INFO] Assume RAM Role specified in provider block assume_role { ... }")
			credentialType = "ram_role_arn"
			credentialConfig.RoleArn = &c.RamRoleArn
			credentialConfig.RoleSessionName = &c.RamRoleSessionName
			credentialConfig.RoleSessionExpiration = &c.RamRoleSessionExpiration
			credentialConfig.Policy = &c.RamRolePolicy
		}
	} else if c.EcsRoleName != "" {
		credentialType = "ecs_ram_role"
		credentialConfig.RoleName = &c.EcsRoleName
	}

	credentialConfig.Type = &credentialType
	return credentialConfig
}
