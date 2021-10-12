package connectivity

import (
	roaCS "github.com/alibabacloud-go/cs-20151215/v2/client"
	openapi "github.com/alibabacloud-go/darabonba-openapi/client"
	roa "github.com/alibabacloud-go/tea-roa/client"
	rpc "github.com/alibabacloud-go/tea-rpc/client"
	"github.com/alibabacloud-go/tea/tea"
	"github.com/aliyun/alibaba-cloud-sdk-go/sdk"
	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/endpoints"
	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/requests"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/adb"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/alidns"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/alikafka"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/bssopenapi"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/cbn"
	cdn_new "github.com/aliyun/alibaba-cloud-sdk-go/services/cdn"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/cloudapi"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/cms"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/cr"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/cr_ee"
	officalCS "github.com/aliyun/alibaba-cloud-sdk-go/services/cs"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/ddosbgp"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/ddoscoo"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/dds"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/drds"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/ecs"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/edas"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/elasticsearch"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/emr"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/ess"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/gpdb"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/hbase"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/market"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/maxcompute"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/ots"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/polardb"
	r_kvstore "github.com/aliyun/alibaba-cloud-sdk-go/services/r-kvstore"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/ram"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/rds"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/slb"
	slsPop "github.com/aliyun/alibaba-cloud-sdk-go/services/sls"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/smartag"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/sts"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/vpc"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/yundun_bastionhost"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/yundun_dbaudit"
	"github.com/aliyun/aliyun-datahub-sdk-go/datahub"
	sls "github.com/aliyun/aliyun-log-go-sdk"
	ali_mns "github.com/aliyun/aliyun-mns-go-sdk"
	"github.com/aliyun/aliyun-oss-go-sdk/oss"
	"github.com/aliyun/aliyun-tablestore-go-sdk/tablestore"
	"github.com/aliyun/fc-go-sdk"
	"github.com/denverdino/aliyungo/cdn"
	"github.com/denverdino/aliyungo/cs"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"

	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"
	"regexp"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/aliyun/alibaba-cloud-sdk-go/services/cassandra"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/dcdn"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/eci"
)

type AliyunClient struct {
	Region   Region
	RegionId string
	SourceIp string
	//In order to build ots table client, add accesskey and secretkey in aliyunclient temporarily.
	AccessKey                    string
	SecretKey                    string
	SecurityToken                string
	OtsInstanceName              string
	accountIdMutex               sync.RWMutex
	config                       *Config
	teaSdkConfig                 rpc.Config
	teaRoaSdkConfig              roa.Config
	accountId                    string
	ecsconn                      *ecs.Client
	essconn                      *ess.Client
	rdsconn                      *rds.Client
	vpcconn                      *vpc.Client
	slbconn                      *slb.Client
	alikafkaconn                 *alikafka.Client
	ossconn                      *oss.Client
	dnsconn                      *alidns.Client
	ramconn                      *ram.Client
	csconn                       *cs.Client
	officalCSConn                *officalCS.Client
	roaCSConn                    *roaCS.Client
	cdnconn_new                  *cdn_new.Client
	crconn                       *cr.Client
	creeconn                     *cr_ee.Client
	cdnconn                      *cdn.CdnClient
	otsconn                      *ots.Client
	cmsconn                      *cms.Client
	logconn                      *sls.Client
	fcconn                       *fc.Client
	cenconn                      *cbn.Client
	logpopconn                   *slsPop.Client
	ddsconn                      *dds.Client
	gpdbconn                     *gpdb.Client
	stsconn                      *sts.Client
	rkvconn                      *r_kvstore.Client
	polarDBconn                  *polardb.Client
	dhconn                       datahub.DataHubApi
	mnsconn                      *ali_mns.MNSClient
	cloudapiconn                 *cloudapi.Client
	teaConn                      *rpc.Client
	tablestoreconnByInstanceName map[string]*tablestore.TableStoreClient
	csprojectconnByKey           map[string]*cs.ProjectClient
	drdsconn                     *drds.Client
	elasticsearchconn            *elasticsearch.Client
	ddoscooconn                  *ddoscoo.Client
	ddosbgpconn                  *ddosbgp.Client
	bssopenapiconn               *bssopenapi.Client
	emrconn                      *emr.Client
	sagconn                      *smartag.Client
	dbauditconn                  *yundun_dbaudit.Client
	bastionhostconn              *yundun_bastionhost.Client
	marketconn                   *market.Client
	hbaseconn                    *hbase.Client
	adbconn                      *adb.Client
	cbnConn                      *cbn.Client
	maxcomputeconn               *maxcompute.Client
	dnsConn                      *alidns.Client
	edasconn                     *edas.Client
	bssopenapiConn               *bssopenapi.Client
	alidnsConn                   *alidns.Client
	ddoscooConn                  *ddoscoo.Client
	cassandraConn                *cassandra.Client
	eciConn                      *eci.Client
	ecsConn                      *ecs.Client
	dcdnConn                     *dcdn.Client
	cmsConn                      *cms.Client
	r_kvstoreConn                *r_kvstore.Client
	maxcomputeConn               *maxcompute.Client
}

type ApiVersion string

const (
	ApiVersion20140526 = ApiVersion("2014-05-26")
	ApiVersion20160815 = ApiVersion("2016-08-15")
	ApiVersion20140515 = ApiVersion("2014-05-15")
)

const businessInfoKey = "Terraform"

const DefaultClientRetryCountSmall = 5

const DefaultClientRetryCountMedium = 10

const DefaultClientRetryCountLarge = 15

const Terraform = "HashiCorp-Terraform"

const Provider = "Terraform-Provider"

const Module = "Terraform-Module"

var goSdkMutex = sync.RWMutex{} // The Go SDK is not thread-safe
var loadSdkfromRemoteMutex = sync.Mutex{}
var loadSdkEndpointMutex = sync.Mutex{}

// The main version number that is being run at the moment.
var providerVersion = "1.132.0"
var terraformVersion = strings.TrimSuffix(schema.Provider{}.TerraformVersion, "-dev")

// Temporarily maintain map for old ecs client methods and store special endpoint information
var EndpointMap = map[string]string{
	"cn-shenzhen-su18-b01":        "ecs.aliyuncs.com",
	"cn-beijing":                  "ecs.aliyuncs.com",
	"cn-shenzhen-st4-d01":         "ecs.aliyuncs.com",
	"cn-haidian-cm12-c01":         "ecs.aliyuncs.com",
	"cn-hangzhou-internal-prod-1": "ecs.aliyuncs.com",
	"cn-qingdao":                  "ecs.aliyuncs.com",
	"cn-shanghai":                 "ecs.aliyuncs.com",
	"cn-shanghai-finance-1":       "ecs.aliyuncs.com",
	"cn-hongkong":                 "ecs.aliyuncs.com",
	"us-west-1":                   "ecs.aliyuncs.com",
	"cn-shenzhen":                 "ecs.aliyuncs.com",
	"cn-shanghai-et15-b01":        "ecs.aliyuncs.com",
	"cn-hangzhou-bj-b01":          "ecs.aliyuncs.com",
	"cn-zhangbei-na61-b01":        "ecs.aliyuncs.com",
	"cn-shenzhen-finance-1":       "ecs.aliyuncs.com",
	"cn-shanghai-et2-b01":         "ecs.aliyuncs.com",
	"ap-southeast-1":              "ecs.aliyuncs.com",
	"cn-beijing-nu16-b01":         "ecs.aliyuncs.com",
	"us-east-1":                   "ecs.aliyuncs.com",
	"cn-fujian":                   "ecs.aliyuncs.com",
	"cn-hangzhou":                 "ecs.aliyuncs.com",
}

// Client for AliyunClient
func (c *Config) Client() (*AliyunClient, error) {
	// Get the auth and region. This can fail if keys/regions were not
	// specified and we're attempting to use the environment.
	if !c.SkipRegionValidation {
		err := c.loadAndValidate()
		if err != nil {
			return nil, err
		}
	}
	loadLocalEndpoint = hasLocalEndpoint()
	if hasLocalEndpoint() {
		if err := c.loadEndpointFromLocal(); err != nil {
			return nil, err
		}
	}
	teaSdkConfig, err := c.getTeaDslSdkConfig(true)
	if err != nil {
		return nil, err
	}
	teaRoaSdkConfig, err := c.getTeaRoaDslSdkConfig(true)
	if err != nil {
		return nil, err
	}
	return &AliyunClient{
		config:                       c,
		teaSdkConfig:                 teaSdkConfig,
		teaRoaSdkConfig:              teaRoaSdkConfig,
		SourceIp:                     c.SourceIp,
		Region:                       c.Region,
		RegionId:                     c.RegionId,
		AccessKey:                    c.AccessKey,
		SecretKey:                    c.SecretKey,
		SecurityToken:                c.SecurityToken,
		OtsInstanceName:              c.OtsInstanceName,
		accountId:                    c.AccountId,
		tablestoreconnByInstanceName: make(map[string]*tablestore.TableStoreClient),
		csprojectconnByKey:           make(map[string]*cs.ProjectClient),
	}, nil
}

func (client *AliyunClient) WithEcsClient(do func(*ecs.Client) (interface{}, error)) (interface{}, error) {
	// Initialize the ECS client if necessary
	if client.ecsconn == nil {
		endpoint := client.config.EcsEndpoint
		if endpoint == "" {
			endpoint = loadEndpoint(client.config.RegionId, ECSCode)
			if endpoint == "" {
				endpoint = EndpointMap[client.config.RegionId]
			}
		}
		if endpoint != "" {
			endpoints.AddEndpointMapping(client.config.RegionId, string(ECSCode), endpoint)
		}
		ecsconn, err := ecs.NewClientWithOptions(client.config.RegionId, client.getSdkConfig().WithTimeout(time.Duration(60)*time.Second), client.config.getAuthCredential(true))
		if err != nil {
			return nil, fmt.Errorf("unable to initialize the ECS client: %#v", err)
		}

		//if _, err := ecsconn.DescribeRegions(ecs.CreateDescribeRegionsRequest()); err != nil {
		//	return nil, err
		//}
		ecsconn.SetReadTimeout(time.Duration(client.config.ClientReadTimeout) * time.Millisecond)
		ecsconn.SetConnectTimeout(time.Duration(client.config.ClientConnectTimeout) * time.Millisecond)
		ecsconn.AppendUserAgent(Terraform, terraformVersion)
		ecsconn.AppendUserAgent(Provider, providerVersion)
		ecsconn.AppendUserAgent(Module, client.config.ConfigurationSource)
		client.ecsconn = ecsconn
	}

	return do(client.ecsconn)
}

func (client *AliyunClient) WithOfficalCSClient(do func(*officalCS.Client) (interface{}, error)) (interface{}, error) {
	goSdkMutex.Lock()
	defer goSdkMutex.Unlock()

	// Initialize the CS client if necessary
	if client.officalCSConn == nil {
		endpoint := client.config.CsEndpoint
		if endpoint == "" {
			endpoint = loadEndpoint(client.config.RegionId, CONTAINCode)
		}
		if endpoint != "" {
			endpoints.AddEndpointMapping(client.config.RegionId, string(CONTAINCode), endpoint)
		}
		csconn, err := officalCS.NewClientWithOptions(client.config.RegionId, client.getSdkConfig(), client.config.getAuthCredential(true))
		if err != nil {
			return nil, fmt.Errorf("unable to initialize the CS client: %#v", err)
		}

		csconn.SetReadTimeout(time.Duration(client.config.ClientReadTimeout) * time.Millisecond)
		csconn.SetConnectTimeout(time.Duration(client.config.ClientConnectTimeout) * time.Millisecond)
		csconn.AppendUserAgent(Terraform, terraformVersion)
		csconn.AppendUserAgent(Provider, providerVersion)
		csconn.AppendUserAgent(Module, client.config.ConfigurationSource)
		client.officalCSConn = csconn
	}

	return do(client.officalCSConn)
}

func (client *AliyunClient) WithRdsClient(do func(*rds.Client) (interface{}, error)) (interface{}, error) {
	// Initialize the RDS client if necessary
	if client.rdsconn == nil {
		endpoint := client.config.RdsEndpoint
		if endpoint == "" {
			endpoint = loadEndpoint(client.config.RegionId, RDSCode)
		}
		if endpoint != "" {
			endpoints.AddEndpointMapping(client.config.RegionId, string(RDSCode), endpoint)
		}
		rdsconn, err := rds.NewClientWithOptions(client.config.RegionId, client.getSdkConfig(), client.config.getAuthCredential(true))
		if err != nil {
			return nil, fmt.Errorf("unable to initialize the RDS client: %#v", err)
		}

		rdsconn.SetReadTimeout(time.Duration(client.config.ClientReadTimeout) * time.Millisecond)
		rdsconn.SetConnectTimeout(time.Duration(client.config.ClientConnectTimeout) * time.Millisecond)
		rdsconn.AppendUserAgent(Terraform, terraformVersion)
		rdsconn.AppendUserAgent(Provider, providerVersion)
		rdsconn.AppendUserAgent(Module, client.config.ConfigurationSource)
		client.rdsconn = rdsconn
	}

	return do(client.rdsconn)
}

func (client *AliyunClient) WithPolarDBClient(do func(*polardb.Client) (interface{}, error)) (interface{}, error) {
	// Initialize the PolarDB client if necessary
	if client.polarDBconn == nil {
		endpoint := client.config.PolarDBEndpoint
		if endpoint == "" {
			endpoint = loadEndpoint(client.config.RegionId, POLARDBCode)
			if endpoint == "" {
				endpoint = fmt.Sprintf("%s.polardb.aliyuncs.com", client.config.RegionId)
			}
		}

		polarDBconn, err := polardb.NewClientWithOptions(client.config.RegionId, client.getSdkConfig(), client.config.getAuthCredential(true))
		if err != nil {
			return nil, fmt.Errorf("unable to initialize the PolarDB client: %#v", err)

		}

		polarDBconn.SetReadTimeout(time.Duration(client.config.ClientReadTimeout) * time.Millisecond)
		polarDBconn.SetConnectTimeout(time.Duration(client.config.ClientConnectTimeout) * time.Millisecond)
		polarDBconn.AppendUserAgent(Terraform, terraformVersion)
		polarDBconn.AppendUserAgent(Provider, providerVersion)
		polarDBconn.AppendUserAgent(Module, client.config.ConfigurationSource)
		client.polarDBconn = polarDBconn
	}

	return do(client.polarDBconn)
}

func (client *AliyunClient) WithSlbClient(do func(*slb.Client) (interface{}, error)) (interface{}, error) {
	// Initialize the SLB client if necessary
	if client.slbconn == nil {
		endpoint := client.config.SlbEndpoint
		if endpoint == "" {
			endpoint = loadEndpoint(client.config.RegionId, SLBCode)
		}
		if endpoint != "" {
			endpoints.AddEndpointMapping(client.config.RegionId, string(SLBCode), endpoint)
		}
		slbconn, err := slb.NewClientWithOptions(client.config.RegionId, client.getSdkConfig(), client.config.getAuthCredential(true))
		if err != nil {
			return nil, fmt.Errorf("unable to initialize the SLB client: %#v", err)
		}

		slbconn.SetReadTimeout(time.Duration(client.config.ClientReadTimeout) * time.Millisecond)
		slbconn.SetConnectTimeout(time.Duration(client.config.ClientConnectTimeout) * time.Millisecond)
		slbconn.AppendUserAgent(Terraform, terraformVersion)
		slbconn.AppendUserAgent(Provider, providerVersion)
		slbconn.AppendUserAgent(Module, client.config.ConfigurationSource)
		client.slbconn = slbconn
	}

	return do(client.slbconn)
}

func (client *AliyunClient) WithVpcClient(do func(*vpc.Client) (interface{}, error)) (interface{}, error) {
	// Initialize the VPC client if necessary
	if client.vpcconn == nil {
		endpoint := client.config.VpcEndpoint
		if endpoint == "" {
			endpoint = loadEndpoint(client.config.RegionId, VPCCode)
		}
		if endpoint != "" {
			endpoints.AddEndpointMapping(client.config.RegionId, string(VPCCode), endpoint)
		}
		vpcconn, err := vpc.NewClientWithOptions(client.config.RegionId, client.getSdkConfig().WithTimeout(time.Duration(60)*time.Second), client.config.getAuthCredential(true))
		if err != nil {
			return nil, fmt.Errorf("unable to initialize the VPC client: %#v", err)
		}

		vpcconn.SetReadTimeout(time.Duration(client.config.ClientReadTimeout) * time.Millisecond)
		vpcconn.SetConnectTimeout(time.Duration(client.config.ClientConnectTimeout) * time.Millisecond)
		vpcconn.AppendUserAgent(Terraform, terraformVersion)
		vpcconn.AppendUserAgent(Provider, providerVersion)
		vpcconn.AppendUserAgent(Module, client.config.ConfigurationSource)
		client.vpcconn = vpcconn
	}

	return do(client.vpcconn)
}

func (client *AliyunClient) NewEcsClient() (*rpc.Client, error) {
	productCode := "ecs"
	endpoint := ""
	if v, ok := client.config.Endpoints[productCode]; !ok || v.(string) == "" {
		if err := client.loadEndpoint(productCode); err != nil {
			return nil, err
		}
	}
	if v, ok := client.config.Endpoints[productCode]; ok && v.(string) != "" {
		endpoint = v.(string)
		if endpoint == "ecs-cn-hangzhou.aliyuncs.com" {
			endpoint = "ecs.aliyuncs.com"
		}
	}
	if endpoint == "" {
		return nil, fmt.Errorf("[ERROR] missing the product %s endpoint.", productCode)
	}

	sdkConfig := client.teaSdkConfig
	sdkConfig.SetEndpoint(endpoint).SetReadTimeout(60000)

	conn, err := rpc.NewClient(&sdkConfig)
	if err != nil {
		return nil, fmt.Errorf("unable to initialize the %s client: %#v", productCode, err)
	}

	return conn, nil
}

func (client *AliyunClient) NewVpcClient() (*rpc.Client, error) {
	productCode := "vpc"
	endpoint := ""
	if v, ok := client.config.Endpoints[productCode]; !ok || v.(string) == "" {
		if err := client.loadEndpoint(productCode); err != nil {
			return nil, err
		}
	}
	if v, ok := client.config.Endpoints[productCode]; ok && v.(string) != "" {
		endpoint = v.(string)
	}
	if endpoint == "" {
		return nil, fmt.Errorf("[ERROR] missing the product %s endpoint.", productCode)
	}

	sdkConfig := client.teaSdkConfig
	sdkConfig.SetEndpoint(endpoint)

	conn, err := rpc.NewClient(&sdkConfig)
	if err != nil {
		return nil, fmt.Errorf("unable to initialize the %s client: %#v", productCode, err)
	}

	return conn, nil
}

func (client *AliyunClient) NewRdsClient() (*rpc.Client, error) {
	productCode := "rds"
	endpoint := ""
	if v, ok := client.config.Endpoints[productCode]; !ok || v.(string) == "" {
		if err := client.loadEndpoint(productCode); err != nil {
			return nil, err
		}
	}
	if v, ok := client.config.Endpoints[productCode]; ok && v.(string) != "" {
		endpoint = v.(string)
	}
	if endpoint == "" {
		return nil, fmt.Errorf("[ERROR] missing the product %s endpoint.", productCode)
	}

	sdkConfig := client.teaSdkConfig
	sdkConfig.SetEndpoint(endpoint)

	conn, err := rpc.NewClient(&sdkConfig)
	if err != nil {
		return nil, fmt.Errorf("unable to initialize the %s client: %#v", productCode, err)
	}

	return conn, nil
}

func (client *AliyunClient) NewPolarDBClient() (*rpc.Client, error) {
	productCode := "polardb"
	endpoint := ""
	if v, ok := client.config.Endpoints[productCode]; !ok || v.(string) == "" {
		if err := client.loadEndpoint(productCode); err != nil {
			return nil, err
		}
	}
	if v, ok := client.config.Endpoints[productCode]; ok && v.(string) != "" {
		endpoint = v.(string)
	}
	if endpoint == "" {
		return nil, fmt.Errorf("[ERROR] missing the product %s endpoint.", productCode)
	}

	sdkConfig := client.teaSdkConfig
	sdkConfig.SetEndpoint(endpoint)

	conn, err := rpc.NewClient(&sdkConfig)
	if err != nil {
		return nil, fmt.Errorf("unable to initialize the %s client: %#v", productCode, err)
	}

	return conn, nil
}

func (client *AliyunClient) WithCenClient(do func(*cbn.Client) (interface{}, error)) (interface{}, error) {
	// Initialize the CEN client if necessary
	if client.cenconn == nil {
		endpoint := client.config.CenEndpoint
		if endpoint == "" {
			endpoint = loadEndpoint(client.config.RegionId, CbnCode)
		}
		if endpoint != "" {
			endpoints.AddEndpointMapping(client.config.RegionId, string(CbnCode), endpoint)
		}
		cenconn, err := cbn.NewClientWithOptions(client.config.RegionId, client.getSdkConfig(), client.config.getAuthCredential(true))
		if err != nil {
			return nil, fmt.Errorf("unable to initialize the CEN client: %#v", err)
		}

		cenconn.SetReadTimeout(time.Duration(client.config.ClientReadTimeout) * time.Millisecond)
		cenconn.SetConnectTimeout(time.Duration(client.config.ClientConnectTimeout) * time.Millisecond)
		cenconn.AppendUserAgent(Terraform, terraformVersion)
		cenconn.AppendUserAgent(Provider, providerVersion)
		cenconn.AppendUserAgent(Module, client.config.ConfigurationSource)
		client.cenconn = cenconn
	}

	return do(client.cenconn)
}

func (client *AliyunClient) WithEssClient(do func(*ess.Client) (interface{}, error)) (interface{}, error) {
	// Initialize the ESS client if necessary
	if client.essconn == nil {
		endpoint := client.config.EssEndpoint
		if endpoint == "" {
			endpoint = loadEndpoint(client.config.RegionId, ESSCode)
		}
		if endpoint != "" {
			endpoints.AddEndpointMapping(client.config.RegionId, string(ESSCode), endpoint)
		}
		essconn, err := ess.NewClientWithOptions(client.config.RegionId, client.getSdkConfig(), client.config.getAuthCredential(true))
		if err != nil {
			return nil, fmt.Errorf("unable to initialize the ESS client: %#v", err)
		}

		essconn.SetReadTimeout(time.Duration(client.config.ClientReadTimeout) * time.Millisecond)
		essconn.SetConnectTimeout(time.Duration(client.config.ClientConnectTimeout) * time.Millisecond)
		essconn.AppendUserAgent(Terraform, terraformVersion)
		essconn.AppendUserAgent(Provider, providerVersion)
		essconn.AppendUserAgent(Module, client.config.ConfigurationSource)
		client.essconn = essconn
	}

	return do(client.essconn)
}

func (client *AliyunClient) WithOssClient(do func(*oss.Client) (interface{}, error)) (interface{}, error) {
	goSdkMutex.Lock()
	defer goSdkMutex.Unlock()

	// Initialize the OSS client if necessary
	if client.ossconn == nil {
		schma := strings.ToLower(client.config.Protocol)
		endpoint := client.config.OssEndpoint
		if endpoint == "" {
			endpoint = loadEndpoint(client.config.RegionId, OSSCode)
		}
		if endpoint == "" {
			endpointItem, err := client.describeEndpointForService(strings.ToLower(string(OSSCode)))
			if err != nil {
				log.Printf("describeEndpointForService got an error: %#v.", err)
			}
			endpoint = endpointItem
			if endpoint == "" {
				endpoint = fmt.Sprintf("oss-%s.aliyuncs.com", client.RegionId)
			}
		}
		if !strings.HasPrefix(endpoint, "http") {
			endpoint = fmt.Sprintf("%s://%s", schma, endpoint)
		}

		clientOptions := []oss.ClientOption{oss.UserAgent(client.getUserAgent()),
			oss.SecurityToken(client.config.SecurityToken)}
		proxy, err := client.getHttpProxy()
		if proxy != nil {
			skip, err := client.skipProxy(endpoint)
			if err != nil {
				return nil, err
			}
			if !skip {
				clientOptions = append(clientOptions, oss.Proxy(proxy.String()))
			}
		}

		ossconn, err := oss.New(endpoint, client.config.AccessKey, client.config.SecretKey, clientOptions...)
		if err != nil {
			return nil, fmt.Errorf("unable to initialize the OSS client: %#v", err)
		}

		client.ossconn = ossconn
	}

	return do(client.ossconn)
}

func (client *AliyunClient) WithOssBucketByName(bucketName string, do func(*oss.Bucket) (interface{}, error)) (interface{}, error) {
	return client.WithOssClient(func(ossClient *oss.Client) (interface{}, error) {
		bucket, err := client.ossconn.Bucket(bucketName)
		if err != nil {
			return nil, fmt.Errorf("unable to get the bucket %s: %#v", bucketName, err)
		}
		return do(bucket)
	})
}

func (client *AliyunClient) WithDnsClient(do func(*alidns.Client) (interface{}, error)) (interface{}, error) {
	// Initialize the DNS client if necessary
	if client.dnsconn == nil {
		endpoint := client.config.DnsEndpoint
		if endpoint == "" {
			endpoint = loadEndpoint(client.config.RegionId, DNSCode)
		}
		if endpoint != "" {
			endpoints.AddEndpointMapping(client.config.RegionId, string(DNSCode), endpoint)
		}

		dnsconn, err := alidns.NewClientWithOptions(client.config.RegionId, client.getSdkConfig(), client.config.getAuthCredential(true))
		if err != nil {
			return nil, fmt.Errorf("unable to initialize the DNS client: %#v", err)
		}
		dnsconn.SetReadTimeout(time.Duration(client.config.ClientReadTimeout) * time.Millisecond)
		dnsconn.SetConnectTimeout(time.Duration(client.config.ClientConnectTimeout) * time.Millisecond)
		dnsconn.AppendUserAgent(Terraform, terraformVersion)
		dnsconn.AppendUserAgent(Provider, providerVersion)
		dnsconn.AppendUserAgent(Module, client.config.ConfigurationSource)
		client.dnsconn = dnsconn
	}

	return do(client.dnsconn)
}

func (client *AliyunClient) WithRamClient(do func(*ram.Client) (interface{}, error)) (interface{}, error) {
	// Initialize the RAM client if necessary
	if client.ramconn == nil {
		endpoint := client.config.RamEndpoint
		if endpoint == "" {
			endpoint = loadEndpoint(client.config.RegionId, RAMCode)
		}
		if strings.HasPrefix(endpoint, "http") {
			endpoint = fmt.Sprintf("https://%s", strings.TrimPrefix(endpoint, "http://"))
		}
		if endpoint != "" {
			endpoints.AddEndpointMapping(client.config.RegionId, string(RAMCode), endpoint)
		}

		ramconn, err := ram.NewClientWithOptions(client.config.RegionId, client.getSdkConfig(), client.config.getAuthCredential(true))
		if err != nil {
			return nil, fmt.Errorf("unable to initialize the RAM client: %#v", err)
		}
		ramconn.SetReadTimeout(time.Duration(client.config.ClientReadTimeout) * time.Millisecond)
		ramconn.SetConnectTimeout(time.Duration(client.config.ClientConnectTimeout) * time.Millisecond)
		ramconn.AppendUserAgent(Terraform, terraformVersion)
		ramconn.AppendUserAgent(Provider, providerVersion)
		ramconn.AppendUserAgent(Module, client.config.ConfigurationSource)
		client.ramconn = ramconn
	}

	return do(client.ramconn)
}

func (client *AliyunClient) WithCsClient(do func(*cs.Client) (interface{}, error)) (interface{}, error) {
	goSdkMutex.Lock()
	defer goSdkMutex.Unlock()

	// Initialize the CS client if necessary
	if client.csconn == nil {
		csconn := cs.NewClientForAussumeRole(client.config.AccessKey, client.config.SecretKey, client.config.SecurityToken)
		csconn.SetUserAgent(client.getUserAgent())
		endpoint := client.config.CsEndpoint
		if endpoint == "" {
			endpoint = loadEndpoint(client.config.RegionId, CONTAINCode)
		}
		if endpoint != "" {
			if !strings.HasPrefix(endpoint, "http") {
				endpoint = fmt.Sprintf("https://%s", strings.TrimPrefix(endpoint, "://"))
			}
			csconn.SetEndpoint(endpoint)
		}
		client.csconn = csconn
	}

	return do(client.csconn)
}

func (client *AliyunClient) NewRoaCsClient() (*roaCS.Client, error) {
	endpoint := client.config.CsEndpoint
	if endpoint == "" {
		endpoint = OpenAckService
	}
	// Initialize the CS client if necessary
	roaCSConn, err := roaCS.NewClient(&openapi.Config{
		AccessKeyId:     tea.String(client.config.AccessKey),
		AccessKeySecret: tea.String(client.config.SecretKey),
		SecurityToken:   tea.String(client.config.SecurityToken),
		RegionId:        tea.String(client.config.RegionId),
		UserAgent:       tea.String(client.getUserAgent()),
		Endpoint:        tea.String(endpoint),
		ReadTimeout:     tea.Int(client.config.ClientReadTimeout),
		ConnectTimeout:  tea.Int(client.config.ClientConnectTimeout),
	})
	if err != nil {
		return nil, err
	}

	return roaCSConn, nil
}

func (client *AliyunClient) WithCrClient(do func(*cr.Client) (interface{}, error)) (interface{}, error) {
	// Initialize the CR client if necessary
	if client.crconn == nil {
		endpoint := client.config.CrEndpoint
		if endpoint == "" {
			endpoint = loadEndpoint(client.config.RegionId, CRCode)
			if endpoint == "" {
				endpoint = fmt.Sprintf("cr.%s.aliyuncs.com", client.config.RegionId)
			}
		}
		if endpoint != "" {
			endpoints.AddEndpointMapping(client.config.RegionId, string(CRCode), endpoint)
		}
		crconn, err := cr.NewClientWithOptions(client.config.RegionId, client.getSdkConfig(), client.config.getAuthCredential(true))
		if err != nil {
			return nil, fmt.Errorf("unable to initialize the CR client: %#v", err)
		}
		crconn.SetReadTimeout(time.Duration(client.config.ClientReadTimeout) * time.Millisecond)
		crconn.SetConnectTimeout(time.Duration(client.config.ClientConnectTimeout) * time.Millisecond)
		crconn.AppendUserAgent(Terraform, terraformVersion)
		crconn.AppendUserAgent(Provider, providerVersion)
		crconn.AppendUserAgent(Module, client.config.ConfigurationSource)
		client.crconn = crconn
	}

	return do(client.crconn)
}

func (client *AliyunClient) WithCrEEClient(do func(*cr_ee.Client) (interface{}, error)) (interface{}, error) {
	// Initialize the CR EE client if necessary
	if client.creeconn == nil {
		endpoint := client.config.CrEndpoint
		if endpoint == "" {
			endpoint = loadEndpoint(client.config.RegionId, CRCode)
			if endpoint == "" {
				endpoint = fmt.Sprintf("cr.%s.aliyuncs.com", client.config.RegionId)
			}
		}
		if endpoint != "" {
			endpoints.AddEndpointMapping(client.config.RegionId, string(CRCode), endpoint)
		}
		creeconn, err := cr_ee.NewClientWithOptions(client.config.RegionId, client.getSdkConfig(), client.config.getAuthCredential(true))
		if err != nil {
			return nil, fmt.Errorf("unable to initialize the CR EE client: %#v", err)
		}
		creeconn.SetReadTimeout(time.Duration(client.config.ClientReadTimeout) * time.Millisecond)
		creeconn.SetConnectTimeout(time.Duration(client.config.ClientConnectTimeout) * time.Millisecond)
		creeconn.AppendUserAgent(Terraform, terraformVersion)
		creeconn.AppendUserAgent(Provider, providerVersion)
		creeconn.AppendUserAgent(Module, client.config.ConfigurationSource)
		client.creeconn = creeconn
	}

	return do(client.creeconn)
}

func (client *AliyunClient) WithCdnClient(do func(*cdn.CdnClient) (interface{}, error)) (interface{}, error) {
	goSdkMutex.Lock()
	defer goSdkMutex.Unlock()

	// Initialize the CDN client if necessary
	if client.cdnconn == nil {
		cdnconn := cdn.NewClient(client.config.AccessKey, client.config.SecretKey)
		cdnconn.SetBusinessInfo(businessInfoKey)
		cdnconn.SetUserAgent(client.getUserAgent())
		cdnconn.SetSecurityToken(client.config.SecurityToken)
		endpoint := client.config.CdnEndpoint
		if endpoint == "" {
			endpoint = loadEndpoint(client.config.RegionId, CDNCode)
		}
		if endpoint != "" && !strings.HasPrefix(endpoint, "http") {
			cdnconn.SetEndpoint(fmt.Sprintf("https://%s", strings.TrimPrefix(endpoint, "://")))
		}
		client.cdnconn = cdnconn
	}
	return do(client.cdnconn)
}

func (client *AliyunClient) WithCdnClient_new(do func(*cdn_new.Client) (interface{}, error)) (interface{}, error) {
	// Initialize the CDN client if necessary
	if client.cdnconn_new == nil {
		endpoint := client.config.CdnEndpoint
		if endpoint == "" {
			endpoint = loadEndpoint(client.config.RegionId, CDNCode)
		}
		if endpoint != "" {
			endpoints.AddEndpointMapping(client.config.RegionId, string(CDNCode), endpoint)
		}
		cdnconn, err := cdn_new.NewClientWithOptions(client.config.RegionId, client.getSdkConfig(), client.config.getAuthCredential(true))
		if err != nil {
			return nil, fmt.Errorf("unable to initialize the CDN client: %#v", err)
		}
		cdnconn.SetReadTimeout(time.Duration(client.config.ClientReadTimeout) * time.Millisecond)
		cdnconn.SetConnectTimeout(time.Duration(client.config.ClientConnectTimeout) * time.Millisecond)
		cdnconn.AppendUserAgent(Terraform, terraformVersion)
		cdnconn.AppendUserAgent(Provider, providerVersion)
		cdnconn.AppendUserAgent(Module, client.config.ConfigurationSource)
		client.cdnconn_new = cdnconn
	}

	return do(client.cdnconn_new)
}

func (client *AliyunClient) WithOtsClient(do func(*ots.Client) (interface{}, error)) (interface{}, error) {
	// Initialize the OTS client if necessary
	if client.otsconn == nil {
		endpoint := client.config.OtsEndpoint
		if endpoint == "" {
			endpoint = loadEndpoint(client.config.RegionId, OTSCode)
		}
		if endpoint != "" {
			endpoints.AddEndpointMapping(client.config.RegionId, string(OTSCode), endpoint)
		}
		otsconn, err := ots.NewClientWithOptions(client.config.RegionId, client.getSdkConfig(), client.config.getAuthCredential(true))
		if err != nil {
			return nil, fmt.Errorf("unable to initialize the OTS client: %#v", err)
		}

		otsconn.SetReadTimeout(time.Duration(client.config.ClientReadTimeout) * time.Millisecond)
		otsconn.SetConnectTimeout(time.Duration(client.config.ClientConnectTimeout) * time.Millisecond)
		otsconn.AppendUserAgent(Terraform, terraformVersion)
		otsconn.AppendUserAgent(Provider, providerVersion)
		otsconn.AppendUserAgent(Module, client.config.ConfigurationSource)
		client.otsconn = otsconn
	}

	return do(client.otsconn)
}

func (client *AliyunClient) WithCmsClient(do func(*cms.Client) (interface{}, error)) (interface{}, error) {
	// Initialize the CMS client if necessary
	if client.cmsconn == nil {
		endpoint := client.config.CmsEndpoint
		if endpoint == "" {
			endpoint = loadEndpoint(client.config.RegionId, CMSCode)
		}
		if endpoint != "" {
			endpoints.AddEndpointMapping(client.config.RegionId, string(CMSCode), endpoint)
		}
		cmsconn, err := cms.NewClientWithOptions(client.config.RegionId, client.getSdkConfig(), client.config.getAuthCredential(true))
		if err != nil {
			return nil, fmt.Errorf("unable to initialize the CMS client: %#v", err)
		}

		cmsconn.SetReadTimeout(time.Duration(client.config.ClientReadTimeout) * time.Millisecond)
		cmsconn.SetConnectTimeout(time.Duration(client.config.ClientConnectTimeout) * time.Millisecond)
		cmsconn.AppendUserAgent(Terraform, terraformVersion)
		cmsconn.AppendUserAgent(Provider, providerVersion)
		cmsconn.AppendUserAgent(Module, client.config.ConfigurationSource)
		client.cmsconn = cmsconn
	}

	return do(client.cmsconn)
}

func (client *AliyunClient) WithStsClient(do func(*sts.Client) (interface{}, error)) (interface{}, error) {
	// Initialize the STS client if necessary
	if client.stsconn == nil {
		endpoint := client.config.StsEndpoint
		if endpoint == "" {
			endpoint = loadEndpoint(client.config.RegionId, STSCode)
		}
		if endpoint != "" {
			endpoints.AddEndpointMapping(client.config.RegionId, string(STSCode), endpoint)
		}
		stsconn, err := sts.NewClientWithOptions(client.config.RegionId, client.getSdkConfig(), client.config.getAuthCredential(true))
		if err != nil {
			return nil, fmt.Errorf("unable to initialize the STS client: %#v", err)
		}

		stsconn.SetReadTimeout(time.Duration(client.config.ClientReadTimeout) * time.Millisecond)
		stsconn.SetConnectTimeout(time.Duration(client.config.ClientConnectTimeout) * time.Millisecond)
		stsconn.AppendUserAgent(Terraform, terraformVersion)
		stsconn.AppendUserAgent(Provider, providerVersion)
		stsconn.AppendUserAgent(Module, client.config.ConfigurationSource)
		client.stsconn = stsconn
	}

	return do(client.stsconn)
}

func (client *AliyunClient) WithLogPopClient(do func(*slsPop.Client) (interface{}, error)) (interface{}, error) {
	// Initialize the HBase client if necessary
	if client.logpopconn == nil {
		logpopconn, err := slsPop.NewClientWithOptions(client.config.RegionId, client.getSdkConfig(), client.config.getAuthCredential(true))
		if err != nil {
			return nil, fmt.Errorf("unable to initialize the sls client: %#v", err)
		}
		logpopconn.SetReadTimeout(time.Duration(client.config.ClientReadTimeout) * time.Millisecond)
		logpopconn.SetConnectTimeout(time.Duration(client.config.ClientConnectTimeout) * time.Millisecond)
		logpopconn.AppendUserAgent(Terraform, terraformVersion)
		logpopconn.AppendUserAgent(Provider, providerVersion)
		logpopconn.AppendUserAgent(Module, client.config.ConfigurationSource)
		logpopconn.Domain = "sls.aliyuncs.com"
		client.logpopconn = logpopconn
	}

	return do(client.logpopconn)
}

func (client *AliyunClient) WithLogClient(do func(*sls.Client) (interface{}, error)) (interface{}, error) {
	goSdkMutex.Lock()
	defer goSdkMutex.Unlock()

	// Initialize the LOG client if necessary
	if client.logconn == nil {
		endpoint := client.config.LogEndpoint
		if endpoint == "" {
			endpoint = loadEndpoint(client.config.RegionId, LOGCode)
			if endpoint == "" {
				endpoint = fmt.Sprintf("%s.log.aliyuncs.com", client.config.RegionId)
			}
		}
		if !strings.HasPrefix(endpoint, "http") {
			endpoint = fmt.Sprintf("https://%s", strings.TrimPrefix(endpoint, "://"))
		}
		client.logconn = &sls.Client{
			AccessKeyID:     client.config.AccessKey,
			AccessKeySecret: client.config.SecretKey,
			Endpoint:        endpoint,
			SecurityToken:   client.config.SecurityToken,
			UserAgent:       client.getUserAgent(),
		}
	}

	return do(client.logconn)
}

func (client *AliyunClient) WithDrdsClient(do func(*drds.Client) (interface{}, error)) (interface{}, error) {
	// Initialize the DRDS client if necessary
	if client.drdsconn == nil {
		endpoint := client.config.DrdsEndpoint
		if endpoint == "" {
			endpoint = loadEndpoint(client.config.RegionId, DRDSCode)
			if endpoint == "" {
				endpoint = fmt.Sprintf("%s.drds.aliyuncs.com", client.config.RegionId)
			}
		}

		drdsconn, err := drds.NewClientWithOptions(client.config.RegionId, client.getSdkConfig(), client.config.getAuthCredential(true))
		if err != nil {
			return nil, fmt.Errorf("unable to initialize the DRDS client: %#v", err)

		}
		drdsconn.SetReadTimeout(time.Duration(client.config.ClientReadTimeout) * time.Millisecond)
		drdsconn.SetConnectTimeout(time.Duration(client.config.ClientConnectTimeout) * time.Millisecond)
		drdsconn.AppendUserAgent(Terraform, terraformVersion)
		drdsconn.AppendUserAgent(Provider, providerVersion)
		drdsconn.AppendUserAgent(Module, client.config.ConfigurationSource)
		client.drdsconn = drdsconn
	}

	return do(client.drdsconn)
}

func (client *AliyunClient) WithDdsClient(do func(*dds.Client) (interface{}, error)) (interface{}, error) {
	// Initialize the DDS client if necessary
	if client.ddsconn == nil {
		endpoint := client.config.DdsEndpoint
		if endpoint == "" {
			endpoint = loadEndpoint(client.config.RegionId, DDSCode)
		}
		if endpoint != "" {
			endpoints.AddEndpointMapping(client.config.RegionId, string(DDSCode), endpoint)
		}
		ddsconn, err := dds.NewClientWithOptions(client.config.RegionId, client.getSdkConfig(), client.config.getAuthCredential(true))
		if err != nil {
			return nil, fmt.Errorf("unable to initialize the DDS client: %#v", err)
		}

		ddsconn.SetReadTimeout(time.Duration(client.config.ClientReadTimeout) * time.Millisecond)
		ddsconn.SetConnectTimeout(time.Duration(client.config.ClientConnectTimeout) * time.Millisecond)
		ddsconn.AppendUserAgent(Terraform, terraformVersion)
		ddsconn.AppendUserAgent(Provider, providerVersion)
		ddsconn.AppendUserAgent(Module, client.config.ConfigurationSource)
		client.ddsconn = ddsconn
	}

	return do(client.ddsconn)
}

func (client *AliyunClient) WithGpdbClient(do func(*gpdb.Client) (interface{}, error)) (interface{}, error) {
	// Initialize the GPDB client if necessary
	if client.gpdbconn == nil {
		endpoint := client.config.GpdbEnpoint
		if endpoint == "" {
			endpoint = loadEndpoint(client.config.RegionId, GPDBCode)
		}
		if endpoint != "" {
			endpoints.AddEndpointMapping(client.config.RegionId, string(GPDBCode), endpoint)
		}
		gpdbconn, err := gpdb.NewClientWithOptions(client.config.RegionId, client.getSdkConfig(), client.config.getAuthCredential(true))
		if err != nil {
			return nil, fmt.Errorf("unable to initialize the GPDB client: %#v", err)
		}
		gpdbconn.SetReadTimeout(time.Duration(client.config.ClientReadTimeout) * time.Millisecond)
		gpdbconn.SetConnectTimeout(time.Duration(client.config.ClientConnectTimeout) * time.Millisecond)
		gpdbconn.AppendUserAgent(Terraform, terraformVersion)
		gpdbconn.AppendUserAgent(Provider, providerVersion)
		gpdbconn.AppendUserAgent(Module, client.config.ConfigurationSource)
		client.gpdbconn = gpdbconn
	}

	return do(client.gpdbconn)
}

func (client *AliyunClient) NewGpdbClient() (*rpc.Client, error) {
	productCode := "gpdb"
	endpoint := ""
	if v, ok := client.config.Endpoints[productCode]; !ok || v.(string) == "" {
		if err := client.loadEndpoint(productCode); err != nil {
			return nil, err
		}
	}
	if v, ok := client.config.Endpoints[productCode]; ok && v.(string) != "" {
		endpoint = v.(string)
	}
	if endpoint == "" {
		return nil, fmt.Errorf("[ERROR] missing the product %s endpoint.", productCode)
	}

	sdkConfig := client.teaSdkConfig
	sdkConfig.SetEndpoint(endpoint)

	conn, err := rpc.NewClient(&sdkConfig)
	if err != nil {
		return nil, fmt.Errorf("unable to initialize the %s client: %#v", productCode, err)
	}

	return conn, nil
}

func (client *AliyunClient) WithRkvClient(do func(*r_kvstore.Client) (interface{}, error)) (interface{}, error) {
	// Initialize the RKV client if necessary
	if client.rkvconn == nil {
		endpoint := client.config.KVStoreEndpoint
		if endpoint == "" {
			endpoint = loadEndpoint(client.config.RegionId, KVSTORECode)
		}
		if endpoint != "" {
			endpoints.AddEndpointMapping(client.config.RegionId, fmt.Sprintf("R-%s", string(KVSTORECode)), endpoint)
		}
		rkvconn, err := r_kvstore.NewClientWithOptions(client.config.RegionId, client.getSdkConfig(), client.config.getAuthCredential(true))
		if err != nil {
			return nil, fmt.Errorf("unable to initialize the RKV client: %#v", err)
		}
		rkvconn.SetReadTimeout(time.Duration(client.config.ClientReadTimeout) * time.Millisecond)
		rkvconn.SetConnectTimeout(time.Duration(client.config.ClientConnectTimeout) * time.Millisecond)
		rkvconn.AppendUserAgent(Terraform, terraformVersion)
		rkvconn.AppendUserAgent(Provider, providerVersion)
		rkvconn.AppendUserAgent(Module, client.config.ConfigurationSource)
		client.rkvconn = rkvconn
	}

	return do(client.rkvconn)
}

func (client *AliyunClient) WithFcClient(do func(*fc.Client) (interface{}, error)) (interface{}, error) {
	goSdkMutex.Lock()
	defer goSdkMutex.Unlock()

	// Initialize the FC client if necessary
	if client.fcconn == nil {
		endpoint := client.config.FcEndpoint
		if endpoint == "" {
			endpoint = loadEndpoint(client.config.RegionId, FCCode)
			if endpoint == "" {
				endpoint = fmt.Sprintf("%s.fc.aliyuncs.com", client.config.RegionId)
			}
		}
		if strings.HasPrefix(endpoint, "http") {
			endpoint = strings.TrimPrefix(strings.TrimPrefix(endpoint, "http://"), "https://")
		}
		accountId, err := client.AccountId()
		if err != nil {
			return nil, err
		}

		config := client.getSdkConfig()
		clientOptions := []fc.ClientOption{fc.WithSecurityToken(client.config.SecurityToken), fc.WithTransport(config.HttpTransport),
			fc.WithTimeout(30), fc.WithRetryCount(DefaultClientRetryCountSmall)}
		fcconn, err := fc.NewClient(fmt.Sprintf("https://%s.%s", accountId, endpoint), string(ApiVersion20160815), client.config.AccessKey, client.config.SecretKey, clientOptions...)
		if err != nil {
			return nil, fmt.Errorf("unable to initialize the FC client: %#v", err)
		}

		fcconn.Config.UserAgent = client.getUserAgent()
		fcconn.Config.SecurityToken = client.config.SecurityToken
		client.fcconn = fcconn
	}

	return do(client.fcconn)
}

func (client *AliyunClient) WithCloudApiClient(do func(*cloudapi.Client) (interface{}, error)) (interface{}, error) {
	// Initialize the Cloud API client if necessary
	if client.cloudapiconn == nil {
		endpoint := client.config.ApigatewayEndpoint
		if endpoint == "" {
			endpoint = loadEndpoint(client.RegionId, CLOUDAPICode)
		}
		if endpoint != "" {
			endpoints.AddEndpointMapping(client.RegionId, "CLOUDAPI", endpoint)
		}
		cloudapiconn, err := cloudapi.NewClientWithOptions(client.RegionId, client.getSdkConfig(), client.config.getAuthCredential(true))
		if err != nil {
			return nil, fmt.Errorf("unable to initialize the CloudAPI client: %#v", err)
		}
		cloudapiconn.SetReadTimeout(time.Duration(client.config.ClientReadTimeout) * time.Millisecond)
		cloudapiconn.SetConnectTimeout(time.Duration(client.config.ClientConnectTimeout) * time.Millisecond)
		cloudapiconn.AppendUserAgent(Terraform, terraformVersion)
		cloudapiconn.AppendUserAgent(Provider, providerVersion)
		cloudapiconn.AppendUserAgent(Module, client.config.ConfigurationSource)
		client.cloudapiconn = cloudapiconn
	}

	return do(client.cloudapiconn)
}

func (client *AliyunClient) NewTeaCommonClient(endpoint string) (*rpc.Client, error) {
	sdkConfig := client.teaSdkConfig
	sdkConfig.SetEndpoint(endpoint)

	conn, err := rpc.NewClient(&sdkConfig)
	if err != nil {
		return nil, fmt.Errorf("unable to initialize the tea client: %#v", err)
	}

	return conn, nil
}
func (client *AliyunClient) NewTeaRoaCommonClient(endpoint string) (*roa.Client, error) {
	sdkConfig := client.teaRoaSdkConfig
	sdkConfig.SetEndpoint(endpoint)

	conn, err := roa.NewClient(&sdkConfig)
	if err != nil {
		return nil, fmt.Errorf("unable to initialize the tea roa client: %#v", err)
	}

	return conn, nil
}

func (client *AliyunClient) WithDataHubClient(do func(api datahub.DataHubApi) (interface{}, error)) (interface{}, error) {
	goSdkMutex.Lock()
	defer goSdkMutex.Unlock()

	// Initialize the DataHub client if necessary
	if client.dhconn == nil {
		endpoint := client.config.DatahubEndpoint
		if endpoint == "" {
			endpoint = loadEndpoint(client.RegionId, DATAHUBCode)
		}
		if endpoint == "" {
			if client.RegionId == string(APSouthEast1) {
				endpoint = "dh-singapore.aliyuncs.com"
			} else {
				endpoint = fmt.Sprintf("dh-%s.aliyuncs.com", client.RegionId)
			}
		}
		if !strings.HasPrefix(endpoint, "http") {
			endpoint = fmt.Sprintf("https://%s", endpoint)
		}

		account := datahub.NewStsCredential(client.config.AccessKey, client.config.SecretKey, client.config.SecurityToken)
		config := &datahub.Config{
			UserAgent: client.getUserAgent(),
		}

		client.dhconn = datahub.NewClientWithConfig(endpoint, config, account)
	}

	return do(client.dhconn)
}

func (client *AliyunClient) WithMnsClient(do func(*ali_mns.MNSClient) (interface{}, error)) (interface{}, error) {
	goSdkMutex.Lock()
	defer goSdkMutex.Unlock()

	// Initialize the MNS client if necessary
	if client.mnsconn == nil {
		endpoint := client.config.MnsEndpoint
		if endpoint == "" {
			endpoint = loadEndpoint(client.config.RegionId, MNSCode)
			if endpoint == "" {
				endpoint = fmt.Sprintf("%s.aliyuncs.com", client.config.RegionId)
			}
		}

		accountId, err := client.AccountId()
		if err != nil {
			return nil, err
		}
		if strings.HasPrefix(endpoint, "http") {
			endpoint = strings.TrimPrefix(strings.TrimPrefix(endpoint, "http://"), "https://")
		}
		mnsUrl := fmt.Sprintf("https://%s.mns.%s", accountId, endpoint)

		mnsClient := ali_mns.NewAliMNSClientWithToken(mnsUrl, client.config.AccessKey, client.config.SecretKey, client.config.SecurityToken)
		proxy, err := client.getHttpProxy()
		if proxy != nil {
			skip, err := client.skipProxy(endpoint)
			if err != nil {
				return nil, err
			}
			if !skip {
				mnsClient.SetProxy(proxy.String())
			}
		}
		client.mnsconn = &mnsClient
	}

	return do(client.mnsconn)
}

func (client *AliyunClient) WithElasticsearchClient(do func(*elasticsearch.Client) (interface{}, error)) (interface{}, error) {
	// Initialize the Elasticsearch client if necessary
	if client.elasticsearchconn == nil {
		endpoint := client.config.ElasticsearchEndpoint
		if endpoint == "" {
			endpoint = loadEndpoint(client.config.RegionId, ELASTICSEARCHCode)
		}
		if endpoint != "" {
			endpoints.AddEndpointMapping(client.config.RegionId, string(ELASTICSEARCHCode), endpoint)
		}
		elasticsearchconn, err := elasticsearch.NewClientWithOptions(client.config.RegionId, client.getSdkConfig(), client.config.getAuthCredential(true))
		if err != nil {
			return nil, fmt.Errorf("unable to initialize the Elasticsearch client: %#v", err)
		}
		elasticsearchconn.SetReadTimeout(time.Duration(client.config.ClientReadTimeout) * time.Millisecond)
		elasticsearchconn.SetConnectTimeout(time.Duration(client.config.ClientConnectTimeout) * time.Millisecond)
		elasticsearchconn.AppendUserAgent(Terraform, terraformVersion)
		elasticsearchconn.AppendUserAgent(Provider, providerVersion)
		elasticsearchconn.AppendUserAgent(Module, client.config.ConfigurationSource)
		client.elasticsearchconn = elasticsearchconn
	}

	return do(client.elasticsearchconn)
}

func (client *AliyunClient) NewElasticsearchClient() (*roa.Client, error) {
	productCode := "elasticsearch"
	endpoint := ""
	if v, ok := client.config.Endpoints[productCode]; !ok || v.(string) == "" {
		if err := client.loadEndpoint(productCode); err != nil {
			return nil, err
		}
	}

	if v, ok := client.config.Endpoints[productCode]; ok && v.(string) != "" {
		endpoint = v.(string)
	}
	if endpoint == "" {
		return nil, fmt.Errorf("[ERROR] misssing the product %s endpoint.", productCode)
	}
	roaSdkConfig := client.teaRoaSdkConfig
	roaSdkConfig.SetEndpoint(endpoint)

	conn, err := roa.NewClient(&roaSdkConfig)
	if err != nil {
		return nil, fmt.Errorf("unable to initialize the %s client: %#v", productCode, err)
	}
	return conn, err
}

func (client *AliyunClient) WithMnsQueueManager(do func(ali_mns.AliQueueManager) (interface{}, error)) (interface{}, error) {
	return client.WithMnsClient(func(mnsClient *ali_mns.MNSClient) (interface{}, error) {
		queueManager := ali_mns.NewMNSQueueManager(*mnsClient)
		return do(queueManager)
	})
}

func (client *AliyunClient) WithMnsTopicManager(do func(ali_mns.AliTopicManager) (interface{}, error)) (interface{}, error) {
	return client.WithMnsClient(func(mnsClient *ali_mns.MNSClient) (interface{}, error) {
		topicManager := ali_mns.NewMNSTopicManager(*mnsClient)
		return do(topicManager)
	})
}

func (client *AliyunClient) WithMnsSubscriptionManagerByTopicName(topicName string, do func(ali_mns.AliMNSTopic) (interface{}, error)) (interface{}, error) {
	return client.WithMnsClient(func(mnsClient *ali_mns.MNSClient) (interface{}, error) {
		subscriptionManager := ali_mns.NewMNSTopic(topicName, *mnsClient)
		return do(subscriptionManager)
	})
}

func (client *AliyunClient) WithTableStoreClient(instanceName string, do func(*tablestore.TableStoreClient) (interface{}, error)) (interface{}, error) {
	goSdkMutex.Lock()
	defer goSdkMutex.Unlock()

	// Initialize the TABLESTORE client if necessary
	tableStoreClient, ok := client.tablestoreconnByInstanceName[instanceName]
	if !ok {
		endpoint := client.config.OtsEndpoint
		if endpoint == "" {
			endpoint = loadEndpoint(client.RegionId, OTSCode)
		}
		if endpoint == "" {
			endpoint = fmt.Sprintf("%s.%s.ots.aliyuncs.com", instanceName, client.RegionId)
		}
		if !strings.HasPrefix(endpoint, "https") && !strings.HasPrefix(endpoint, "http") {
			endpoint = fmt.Sprintf("https://%s", endpoint)
		}

		tableStoreClient = tablestore.NewClientWithConfig(endpoint, instanceName, client.config.AccessKey, client.config.SecretKey, client.config.SecurityToken, tablestore.NewDefaultTableStoreConfig())
		client.tablestoreconnByInstanceName[instanceName] = tableStoreClient
	}

	return do(tableStoreClient)
}

func (client *AliyunClient) WithCsProjectClient(clusterId, endpoint string, clusterCerts cs.ClusterCerts, do func(*cs.ProjectClient) (interface{}, error)) (interface{}, error) {
	goSdkMutex.Lock()
	defer goSdkMutex.Unlock()

	// Initialize the PROJECT client if necessary
	key := fmt.Sprintf("%s|%s|%s|%s|%s", clusterId, endpoint, clusterCerts.CA, clusterCerts.Cert, clusterCerts.Key)
	csProjectClient, ok := client.csprojectconnByKey[key]
	if !ok {
		var err error
		csProjectClient, err = cs.NewProjectClient(clusterId, endpoint, clusterCerts)
		if err != nil {
			return nil, fmt.Errorf("Getting Application Client failed by cluster id %s: %#v.", clusterCerts, err)
		}
		csProjectClient.SetDebug(false)
		csProjectClient.SetUserAgent(client.getUserAgent())
		client.csprojectconnByKey[key] = csProjectClient
	}

	return do(csProjectClient)
}

func (client *AliyunClient) NewCommonRequest(product, serviceCode, schema string, apiVersion ApiVersion) (*requests.CommonRequest, error) {
	endpoint := ""
	product = strings.ToLower(product)
	if client.config.Endpoints[product] == nil {
		if err := client.loadEndpoint(product); err != nil {
			return nil, err
		}
	}
	if client.config.Endpoints[product] != nil && client.config.Endpoints[product].(string) != "" {
		endpoint = client.config.Endpoints[product].(string)
	}
	request := requests.NewCommonRequest()
	// Use product code to find product domain
	if endpoint != "" {
		request.Domain = endpoint
	} else {
		// When getting endpoint failed by location, using custom endpoint instead
		request.Domain = fmt.Sprintf("%s.%s.aliyuncs.com", strings.ToLower(serviceCode), client.RegionId)
	}
	request.Version = string(apiVersion)
	request.RegionId = client.RegionId
	request.Product = product
	request.Scheme = schema
	request.SetReadTimeout(time.Duration(client.config.ClientReadTimeout) * time.Millisecond)
	request.SetConnectTimeout(time.Duration(client.config.ClientConnectTimeout) * time.Millisecond)
	request.AppendUserAgent(Terraform, terraformVersion)
	request.AppendUserAgent(Provider, providerVersion)
	request.AppendUserAgent(Module, client.config.ConfigurationSource)
	return request, nil
}

func (client *AliyunClient) AccountId() (string, error) {
	client.accountIdMutex.Lock()
	defer client.accountIdMutex.Unlock()

	if client.accountId == "" {
		log.Printf("[DEBUG] account_id not provided, attempting to retrieve it automatically...")
		identity, err := client.GetCallerIdentity()
		if err != nil {
			return "", err
		}
		if identity.AccountId == "" {
			return "", fmt.Errorf("caller identity doesn't contain any AccountId")
		}
		client.accountId = identity.AccountId
	}
	return client.accountId, nil
}

func (client *AliyunClient) getSdkConfig() *sdk.Config {
	return sdk.NewConfig().
		WithMaxRetryTime(DefaultClientRetryCountSmall).
		WithTimeout(time.Duration(30) * time.Second).
		WithEnableAsync(false).
		WithGoRoutinePoolSize(100).
		WithMaxTaskQueueSize(10000).
		WithDebug(false).
		WithHttpTransport(client.getTransport()).
		WithScheme(client.config.Protocol)
}

func (client *AliyunClient) getUserAgent() string {
	return fmt.Sprintf("%s/%s %s/%s %s/%s", Terraform, terraformVersion, Provider, providerVersion, Module, client.config.ConfigurationSource)
}

func (client *AliyunClient) getTransport() *http.Transport {
	handshakeTimeout, err := strconv.Atoi(os.Getenv("TLSHandshakeTimeout"))
	if err != nil {
		handshakeTimeout = 120
	}
	transport := &http.Transport{}
	transport.TLSHandshakeTimeout = time.Duration(handshakeTimeout) * time.Second

	return transport
}

func (client *AliyunClient) getHttpProxy() (proxy *url.URL, err error) {
	if client.config.Protocol == "HTTPS" {
		if rawurl := os.Getenv("HTTPS_PROXY"); rawurl != "" {
			proxy, err = url.Parse(rawurl)
		} else if rawurl := os.Getenv("https_proxy"); rawurl != "" {
			proxy, err = url.Parse(rawurl)
		}
	} else {
		if rawurl := os.Getenv("HTTP_PROXY"); rawurl != "" {
			proxy, err = url.Parse(rawurl)
		} else if rawurl := os.Getenv("http_proxy"); rawurl != "" {
			proxy, err = url.Parse(rawurl)
		}
	}
	return proxy, err
}

func (client *AliyunClient) skipProxy(endpoint string) (bool, error) {
	var urls []string
	if rawurl := os.Getenv("NO_PROXY"); rawurl != "" {
		urls = strings.Split(rawurl, ",")
	} else if rawurl := os.Getenv("no_proxy"); rawurl != "" {
		urls = strings.Split(rawurl, ",")
	}
	for _, value := range urls {
		if strings.HasPrefix(value, "*") {
			value = fmt.Sprintf(".%s", value)
		}
		noProxyReg, err := regexp.Compile(value)
		if err != nil {
			return false, err
		}
		if noProxyReg.MatchString(endpoint) {
			return true, nil
		}
	}
	return false, nil
}

func (client *AliyunClient) GetCallerIdentity() (*sts.GetCallerIdentityResponse, error) {
	args := sts.CreateGetCallerIdentityRequest()

	endpoint := client.config.StsEndpoint
	if endpoint == "" {
		endpoint = loadEndpoint(client.config.RegionId, STSCode)
	}
	if endpoint != "" {
		endpoints.AddEndpointMapping(client.config.RegionId, string(STSCode), endpoint)
	}
	stsClient, err := sts.NewClientWithOptions(client.config.RegionId, client.getSdkConfig(), client.config.getAuthCredential(true))
	if err != nil {
		return nil, fmt.Errorf("unable to initialize the STS client: %#v", err)
	}

	stsClient.SetReadTimeout(time.Duration(client.config.ClientReadTimeout) * time.Millisecond)
	stsClient.SetConnectTimeout(time.Duration(client.config.ClientConnectTimeout) * time.Millisecond)
	stsClient.AppendUserAgent(Terraform, terraformVersion)
	stsClient.AppendUserAgent(Provider, providerVersion)
	stsClient.AppendUserAgent(Module, client.config.ConfigurationSource)

	identity, err := stsClient.GetCallerIdentity(args)
	if err != nil {
		return nil, err
	}
	if identity == nil {
		return nil, fmt.Errorf("caller identity not found")
	}
	return identity, err
}

func (client *AliyunClient) WithDdoscooClient(do func(*ddoscoo.Client) (interface{}, error)) (interface{}, error) {
	// Initialize the ddoscoo client if necessary
	if client.ddoscooconn == nil {
		endpoint := client.config.DdoscooEndpoint
		if endpoint == "" {
			endpoint = loadEndpoint(client.config.RegionId, DDOSCOOCode)
		}
		if endpoint != "" {
			endpoints.AddEndpointMapping(client.config.RegionId, string(DDOSCOOCode), endpoint)
		}

		ddoscooconn, err := ddoscoo.NewClientWithOptions(client.config.RegionId, client.getSdkConfig(), client.config.getAuthCredential(true))
		if err != nil {
			return nil, fmt.Errorf("unable to initialize the DDOSCOO client: %#v", err)
		}
		ddoscooconn.SetReadTimeout(time.Duration(client.config.ClientReadTimeout) * time.Millisecond)
		ddoscooconn.SetConnectTimeout(time.Duration(client.config.ClientConnectTimeout) * time.Millisecond)
		ddoscooconn.AppendUserAgent(Terraform, terraformVersion)
		ddoscooconn.AppendUserAgent(Provider, providerVersion)
		ddoscooconn.AppendUserAgent(Module, client.config.ConfigurationSource)
		client.ddoscooconn = ddoscooconn

	}

	return do(client.ddoscooconn)
}

func (client *AliyunClient) WithDdosbgpClient(do func(*ddosbgp.Client) (interface{}, error)) (interface{}, error) {
	// Initialize the ddosbgp client if necessary
	if client.ddosbgpconn == nil {
		endpoint := client.config.DdosbgpEndpoint
		if endpoint == "" {
			endpoint = loadEndpoint(client.config.RegionId, DDOSBGPCode)
		}
		if endpoint != "" {
			endpoints.AddEndpointMapping(client.config.RegionId, string(DDOSBGPCode), endpoint)
		}

		ddosbgpconn, err := ddosbgp.NewClientWithOptions(client.config.RegionId, client.getSdkConfig(), client.config.getAuthCredential(true))
		if err != nil {
			return nil, fmt.Errorf("unable to initialize the DDOSBGP client: %#v", err)
		}
		ddosbgpconn.SetReadTimeout(time.Duration(client.config.ClientReadTimeout) * time.Millisecond)
		ddosbgpconn.SetConnectTimeout(time.Duration(client.config.ClientConnectTimeout) * time.Millisecond)
		ddosbgpconn.AppendUserAgent(Terraform, terraformVersion)
		ddosbgpconn.AppendUserAgent(Provider, providerVersion)
		ddosbgpconn.AppendUserAgent(Module, client.config.ConfigurationSource)
		client.ddosbgpconn = ddosbgpconn
	}

	return do(client.ddosbgpconn)
}

func (client *AliyunClient) WithBssopenapiClient(do func(*bssopenapi.Client) (interface{}, error)) (interface{}, error) {
	// Initialize the bssopenapi client if necessary
	if client.bssopenapiconn == nil {
		endpoint := client.config.BssOpenApiEndpoint
		if endpoint == "" {
			endpoint = loadEndpoint(client.config.RegionId, BSSOPENAPICode)
		}
		if endpoint != "" {
			endpoints.AddEndpointMapping(client.config.RegionId, string(BSSOPENAPICode), endpoint)
		}

		// bss endpoint depends on the account type.
		// Domestic account is business.aliyuncs.com (region is cn-hangzhou) and International account is business.ap-southeast-1.aliyuncs.com (region is ap-southeast-1)
		bssopenapiconn, err := bssopenapi.NewClientWithOptions(string(Hangzhou), client.getSdkConfig(), client.config.getAuthCredential(true))
		if err != nil {
			return nil, fmt.Errorf("unable to initialize the BSSOPENAPI client: %#v", err)
		}
		bssopenapiconn.SetReadTimeout(time.Duration(client.config.ClientReadTimeout) * time.Millisecond)
		bssopenapiconn.SetConnectTimeout(time.Duration(client.config.ClientConnectTimeout) * time.Millisecond)
		bssopenapiconn.AppendUserAgent(Terraform, terraformVersion)
		bssopenapiconn.AppendUserAgent(Provider, providerVersion)
		bssopenapiconn.AppendUserAgent(Module, client.config.ConfigurationSource)
		client.bssopenapiconn = bssopenapiconn
	}

	return do(client.bssopenapiconn)
}

func (client *AliyunClient) WithAlikafkaClient(do func(*alikafka.Client) (interface{}, error)) (interface{}, error) {
	// Initialize the alikafka client if necessary
	if client.alikafkaconn == nil {
		endpoint := client.config.AlikafkaEndpoint
		if endpoint == "" {
			endpoint = loadEndpoint(client.config.RegionId, ALIKAFKACode)
		}
		if endpoint != "" {
			endpoints.AddEndpointMapping(client.config.RegionId, string(ALIKAFKACode), endpoint)
		}
		alikafkaconn, err := alikafka.NewClientWithOptions(client.config.RegionId, client.getSdkConfig(), client.config.getAuthCredential(true))
		if err != nil {
			return nil, fmt.Errorf("unable to initialize the ALIKAFKA client: %#v", err)
		}
		alikafkaconn.SetReadTimeout(time.Duration(client.config.ClientReadTimeout) * time.Millisecond)
		alikafkaconn.SetConnectTimeout(time.Duration(client.config.ClientConnectTimeout) * time.Millisecond)
		alikafkaconn.AppendUserAgent(Terraform, terraformVersion)
		alikafkaconn.AppendUserAgent(Provider, providerVersion)
		alikafkaconn.AppendUserAgent(Module, client.config.ConfigurationSource)
		client.alikafkaconn = alikafkaconn
	}

	return do(client.alikafkaconn)
}

func (client *AliyunClient) WithEmrClient(do func(*emr.Client) (interface{}, error)) (interface{}, error) {
	if client.emrconn == nil {
		endpoint := client.config.EmrEndpoint
		if endpoint == "" {
			endpoint = loadEndpoint(client.config.RegionId, EMRCode)
		}
		if endpoint != "" {
			endpoints.AddEndpointMapping(client.config.RegionId, string(EMRCode), endpoint)
		}
		emrConn, err := emr.NewClientWithOptions(client.RegionId, client.getSdkConfig(), client.config.getAuthCredential(true))
		if err != nil {
			return nil, fmt.Errorf("unable to initialize the E-MapReduce client: %#v", err)
		}
		emrConn.SetReadTimeout(time.Duration(client.config.ClientReadTimeout) * time.Millisecond)
		emrConn.SetConnectTimeout(time.Duration(client.config.ClientConnectTimeout) * time.Millisecond)
		emrConn.AppendUserAgent(Terraform, terraformVersion)
		emrConn.AppendUserAgent(Provider, providerVersion)
		emrConn.AppendUserAgent(Module, client.config.ConfigurationSource)
		client.emrconn = emrConn
	}

	return do(client.emrconn)
}

func (client *AliyunClient) WithSagClient(do func(*smartag.Client) (interface{}, error)) (interface{}, error) {
	// Initialize the SAG client if necessary
	if client.sagconn == nil {
		endpoint := client.config.SagEndpoint
		if endpoint == "" {
			endpoint = loadEndpoint(client.config.RegionId, SAGCode)
		}
		if endpoint != "" {
			endpoints.AddEndpointMapping(client.config.RegionId, string(SAGCode), endpoint)
		}
		sagconn, err := smartag.NewClientWithOptions(client.config.RegionId, client.getSdkConfig(), client.config.getAuthCredential(true))
		if err != nil {
			return nil, fmt.Errorf("unable to initialize the SAG client: %#v", err)
		}
		sagconn.SetReadTimeout(time.Duration(client.config.ClientReadTimeout) * time.Millisecond)
		sagconn.SetConnectTimeout(time.Duration(client.config.ClientConnectTimeout) * time.Millisecond)
		sagconn.AppendUserAgent(Terraform, terraformVersion)
		sagconn.AppendUserAgent(Provider, providerVersion)
		sagconn.AppendUserAgent(Module, client.config.ConfigurationSource)
		client.sagconn = sagconn
	}

	return do(client.sagconn)
}

func (client *AliyunClient) WithDbauditClient(do func(*yundun_dbaudit.Client) (interface{}, error)) (interface{}, error) {
	// Initialize the ddoscoo client if necessary
	if client.dbauditconn == nil {
		dbauditconn, err := yundun_dbaudit.NewClientWithOptions(client.config.RegionId, client.getSdkConfig(), client.config.getAuthCredential(true))
		if err != nil {
			return nil, fmt.Errorf("unable to initialize the DBAUDIT client: %#v", err)
		}
		dbauditconn.SetReadTimeout(time.Duration(client.config.ClientReadTimeout) * time.Millisecond)
		dbauditconn.SetConnectTimeout(time.Duration(client.config.ClientConnectTimeout) * time.Millisecond)
		dbauditconn.AppendUserAgent(Terraform, terraformVersion)
		dbauditconn.AppendUserAgent(Provider, providerVersion)
		dbauditconn.AppendUserAgent(Module, client.config.ConfigurationSource)
		client.dbauditconn = dbauditconn
	}

	return do(client.dbauditconn)
}

func (client *AliyunClient) WithBastionhostClient(do func(*yundun_bastionhost.Client) (interface{}, error)) (interface{}, error) {
	if client.bastionhostconn == nil {
		bastionhostconn, err := yundun_bastionhost.NewClientWithOptions(client.config.RegionId, client.getSdkConfig(), client.config.getAuthCredential(true))
		if err != nil {
			return nil, fmt.Errorf("unable to initialize the BASTIONHOST client: %#v", err)
		}
		bastionhostconn.SetReadTimeout(time.Duration(client.config.ClientReadTimeout) * time.Millisecond)
		bastionhostconn.SetConnectTimeout(time.Duration(client.config.ClientConnectTimeout) * time.Millisecond)
		bastionhostconn.AppendUserAgent(Terraform, terraformVersion)
		bastionhostconn.AppendUserAgent(Provider, providerVersion)
		bastionhostconn.AppendUserAgent(Module, client.config.ConfigurationSource)
		client.bastionhostconn = bastionhostconn
	}

	return do(client.bastionhostconn)
}

func (client *AliyunClient) WithMarketClient(do func(*market.Client) (interface{}, error)) (interface{}, error) {
	// Initialize the Market API client if necessary
	if client.marketconn == nil {
		endpoint := client.config.MarketEndpoint
		if endpoint == "" {
			endpoint = loadEndpoint(client.RegionId, MARKETCode)
		}
		if endpoint != "" {
			endpoints.AddEndpointMapping(client.RegionId, "MARKET", endpoint)
		}
		marketconn, err := market.NewClientWithOptions(client.RegionId, client.getSdkConfig(), client.config.getAuthCredential(true))
		if err != nil {
			return nil, fmt.Errorf("unable to initialize the Market client: %#v", err)
		}
		marketconn.SetReadTimeout(time.Duration(client.config.ClientReadTimeout) * time.Millisecond)
		marketconn.SetConnectTimeout(time.Duration(client.config.ClientConnectTimeout) * time.Millisecond)
		marketconn.AppendUserAgent(Terraform, terraformVersion)
		marketconn.AppendUserAgent(Provider, providerVersion)
		marketconn.AppendUserAgent(Module, client.config.ConfigurationSource)
		client.marketconn = marketconn
	}

	return do(client.marketconn)
}

func (client *AliyunClient) WithHbaseClient(do func(*hbase.Client) (interface{}, error)) (interface{}, error) {
	// Initialize the HBase client if necessary
	if client.hbaseconn == nil {
		endpoint := client.config.HBaseEndpoint
		if endpoint == "" {
			endpoint = loadEndpoint(client.config.RegionId, HBASECode)
		}
		if endpoint != "" {
			endpoints.AddEndpointMapping(client.config.RegionId, string(HBASECode), endpoint)
		}
		hbaseconn, err := hbase.NewClientWithOptions(client.config.RegionId, client.getSdkConfig(), client.config.getAuthCredential(true))
		if err != nil {
			return nil, fmt.Errorf("unable to initialize the hbase client: %#v", err)
		}
		hbaseconn.SetReadTimeout(time.Duration(client.config.ClientReadTimeout) * time.Millisecond)
		hbaseconn.SetConnectTimeout(time.Duration(client.config.ClientConnectTimeout) * time.Millisecond)
		hbaseconn.AppendUserAgent(Terraform, terraformVersion)
		hbaseconn.AppendUserAgent(Provider, providerVersion)
		hbaseconn.AppendUserAgent(Module, client.config.ConfigurationSource)
		client.hbaseconn = hbaseconn
	}

	return do(client.hbaseconn)
}

func (client *AliyunClient) NewAdbClient() (*rpc.Client, error) {
	productCode := "ads"
	endpoint := ""
	if client.config.Endpoints[productCode] == nil {
		if err := client.loadEndpoint(productCode); err != nil {
			return nil, err
		}
	}
	if client.config.Endpoints[productCode] != nil && client.config.Endpoints[productCode].(string) != "" {
		endpoint = client.config.Endpoints[productCode].(string)
	}
	if endpoint == "" {
		return nil, fmt.Errorf("[ERROR] missing the product %s endpoint.", productCode)
	}

	sdkConfig := client.teaSdkConfig
	sdkConfig.SetEndpoint(endpoint)

	conn, err := rpc.NewClient(&sdkConfig)
	if err != nil {
		return nil, fmt.Errorf("unable to initialize the %s client: %#v", productCode, err)
	}

	return conn, nil
}

func (client *AliyunClient) WithAdbClient(do func(*adb.Client) (interface{}, error)) (interface{}, error) {
	// Initialize the adb client if necessary
	if client.adbconn == nil {
		endpoint := client.config.AdbEndpoint
		if endpoint == "" {
			endpoint = loadEndpoint(client.config.RegionId, ADBCode)
			if endpoint == "" {
				endpoint = fmt.Sprintf("%s.adb.aliyuncs.com", client.config.RegionId)
			}
		}

		adbconn, err := adb.NewClientWithOptions(client.config.RegionId, client.getSdkConfig(), client.config.getAuthCredential(true))
		if err != nil {
			return nil, fmt.Errorf("unable to initialize the adb client: %#v", err)

		}
		adbconn.SetReadTimeout(time.Duration(client.config.ClientReadTimeout) * time.Millisecond)
		adbconn.SetConnectTimeout(time.Duration(client.config.ClientConnectTimeout) * time.Millisecond)
		adbconn.AppendUserAgent(Terraform, terraformVersion)
		adbconn.AppendUserAgent(Provider, providerVersion)
		adbconn.AppendUserAgent(Module, client.config.ConfigurationSource)
		client.adbconn = adbconn
	}

	return do(client.adbconn)
}

func (client *AliyunClient) NewCbnClient() (*rpc.Client, error) {
	productCode := "cbn"
	endpoint := ""
	if v, ok := client.config.Endpoints[productCode]; !ok || v.(string) == "" {
		if err := client.loadEndpoint(productCode); err != nil {
			return nil, err
		}
	}
	if v, ok := client.config.Endpoints[productCode]; ok && v.(string) != "" {
		endpoint = v.(string)
	}
	if endpoint == "" {
		return nil, fmt.Errorf("[ERROR] missing the product %s endpoint.", productCode)
	}
	sdkConfig := client.teaSdkConfig
	sdkConfig.SetEndpoint(endpoint)
	conn, err := rpc.NewClient(&sdkConfig)
	if err != nil {
		return nil, fmt.Errorf("unable to initialize the %s client: %#v", productCode, err)
	}
	return conn, nil
}

func (client *AliyunClient) WithCbnClient(do func(*cbn.Client) (interface{}, error)) (interface{}, error) {
	if client.cbnConn == nil {
		endpoint := client.config.CbnEndpoint
		if endpoint == "" {
			endpoint = loadEndpoint(client.config.RegionId, CbnCode)
			// compatible with cen
			if endpoint == "" {
				endpoint = "cbn.aliyuncs.com"
			}
		}
		if strings.HasPrefix(endpoint, "http") {
			endpoint = fmt.Sprintf("https://%s", strings.TrimPrefix(endpoint, "http://"))
		}
		if endpoint != "" {
			endpoints.AddEndpointMapping(client.config.RegionId, string(CbnCode), endpoint)
		}

		cbnConn, err := cbn.NewClientWithOptions(client.config.RegionId, client.getSdkConfig(), client.config.getAuthCredential(true))
		if err != nil {
			return nil, fmt.Errorf("unable to initialize the Cbnclient: %#v", err)
		}
		cbnConn.SetReadTimeout(time.Duration(client.config.ClientReadTimeout) * time.Millisecond)
		cbnConn.SetConnectTimeout(time.Duration(client.config.ClientConnectTimeout) * time.Millisecond)
		cbnConn.AppendUserAgent(Terraform, terraformVersion)
		cbnConn.AppendUserAgent(Provider, providerVersion)
		cbnConn.AppendUserAgent(Module, client.config.ConfigurationSource)
		client.cbnConn = cbnConn
	}
	return do(client.cbnConn)
}

func (client *AliyunClient) WithEdasClient(do func(*edas.Client) (interface{}, error)) (interface{}, error) {
	// Initialize the edas client if necessary
	if client.edasconn == nil {
		endpoint := client.config.edasEndpoint
		if endpoint == "" {
			endpoint = loadEndpoint(client.config.RegionId, EDASCode)
		}
		if endpoint != "" {
			endpoints.AddEndpointMapping(client.config.RegionId, string(EDASCode), endpoint)
		}
		edasconn, err := edas.NewClientWithOptions(client.config.RegionId, client.getSdkConfig().WithTimeout(time.Duration(60)*time.Second), client.config.getAuthCredential(true))
		if err != nil {
			return nil, fmt.Errorf("unable to initialize the ALIKAFKA client: %#v", err)
		}
		edasconn.SetReadTimeout(time.Duration(client.config.ClientReadTimeout) * time.Millisecond)
		edasconn.SetConnectTimeout(time.Duration(client.config.ClientConnectTimeout) * time.Millisecond)
		edasconn.AppendUserAgent(Terraform, terraformVersion)
		edasconn.AppendUserAgent(Provider, providerVersion)
		edasconn.AppendUserAgent(Module, client.config.ConfigurationSource)
		client.edasconn = edasconn
	}

	return do(client.edasconn)
}

func (client *AliyunClient) WithAlidnsClient(do func(*alidns.Client) (interface{}, error)) (interface{}, error) {
	if client.alidnsConn == nil {
		endpoint := client.config.AlidnsEndpoint
		if endpoint == "" {
			endpoint = loadEndpoint(client.config.RegionId, AlidnsCode)
		}
		if strings.HasPrefix(endpoint, "http") {
			endpoint = fmt.Sprintf("https://%s", strings.TrimPrefix(endpoint, "http://"))
		}
		if endpoint != "" {
			endpoints.AddEndpointMapping(client.config.RegionId, string(AlidnsCode), endpoint)
		}

		alidnsConn, err := alidns.NewClientWithOptions(client.config.RegionId, client.getSdkConfig(), client.config.getAuthCredential(true))
		if err != nil {
			return nil, fmt.Errorf("unable to initialize the Alidnsclient: %#v", err)
		}
		alidnsConn.SetReadTimeout(time.Duration(client.config.ClientReadTimeout) * time.Millisecond)
		alidnsConn.SetConnectTimeout(time.Duration(client.config.ClientConnectTimeout) * time.Millisecond)
		alidnsConn.AppendUserAgent(Terraform, terraformVersion)
		alidnsConn.AppendUserAgent(Provider, providerVersion)
		alidnsConn.AppendUserAgent(Module, client.config.ConfigurationSource)
		client.alidnsConn = alidnsConn
	}
	return do(client.alidnsConn)
}

func (client *AliyunClient) WithCassandraClient(do func(*cassandra.Client) (interface{}, error)) (interface{}, error) {
	if client.cassandraConn == nil {
		endpoint := client.config.CassandraEndpoint
		if endpoint == "" {
			endpoint = loadEndpoint(client.config.RegionId, CassandraCode)
			endpoints.AddEndpointMapping(client.config.RegionId, string(CassandraCode), endpoint)
		}
		cassandraConn, err := cassandra.NewClientWithOptions(client.config.RegionId, client.getSdkConfig(), client.config.getAuthCredential(true))
		if err != nil {
			return nil, fmt.Errorf("unable to initialize the Cassandraclient: %#v", err)
		}
		cassandraConn.SetReadTimeout(time.Duration(client.config.ClientReadTimeout) * time.Millisecond)
		cassandraConn.SetConnectTimeout(time.Duration(client.config.ClientConnectTimeout) * time.Millisecond)
		cassandraConn.AppendUserAgent(Terraform, terraformVersion)
		cassandraConn.AppendUserAgent(Provider, providerVersion)
		cassandraConn.AppendUserAgent(Module, client.config.ConfigurationSource)
		client.cassandraConn = cassandraConn
	}
	return do(client.cassandraConn)
}

func (client *AliyunClient) WithEciClient(do func(*eci.Client) (interface{}, error)) (interface{}, error) {
	if client.eciConn == nil {
		endpoint := client.config.EciEndpoint
		if endpoint == "" {
			endpoint = loadEndpoint(client.config.RegionId, EciCode)
		}
		if strings.HasPrefix(endpoint, "http") {
			endpoint = fmt.Sprintf("https://%s", strings.TrimPrefix(endpoint, "http://"))
		}
		if endpoint != "" {
			endpoints.AddEndpointMapping(client.config.RegionId, string(EciCode), endpoint)
		}

		eciConn, err := eci.NewClientWithOptions(client.config.RegionId, client.getSdkConfig(), client.config.getAuthCredential(true))
		if err != nil {
			return nil, fmt.Errorf("unable to initialize the Eciclient: %#v", err)
		}
		eciConn.SetReadTimeout(time.Duration(client.config.ClientReadTimeout) * time.Millisecond)
		eciConn.SetConnectTimeout(time.Duration(client.config.ClientConnectTimeout) * time.Millisecond)
		eciConn.AppendUserAgent(Terraform, terraformVersion)
		eciConn.AppendUserAgent(Provider, providerVersion)
		eciConn.AppendUserAgent(Module, client.config.ConfigurationSource)
		client.eciConn = eciConn
	}
	return do(client.eciConn)
}

func (client *AliyunClient) WithDcdnClient(do func(*dcdn.Client) (interface{}, error)) (interface{}, error) {
	if client.dcdnConn == nil {
		endpoint := client.config.DcdnEndpoint
		if endpoint == "" {
			endpoint = loadEndpoint(client.config.RegionId, DcdnCode)
		}
		if strings.HasPrefix(endpoint, "http") {
			endpoint = fmt.Sprintf("https://%s", strings.TrimPrefix(endpoint, "http://"))
		}
		if endpoint != "" {
			endpoints.AddEndpointMapping(client.config.RegionId, string(DcdnCode), endpoint)
		}

		dcdnConn, err := dcdn.NewClientWithOptions(client.config.RegionId, client.getSdkConfig(), client.config.getAuthCredential(true))
		if err != nil {
			return nil, fmt.Errorf("unable to initialize the Dcdnclient: %#v", err)
		}
		dcdnConn.SetReadTimeout(time.Duration(client.config.ClientReadTimeout) * time.Millisecond)
		dcdnConn.SetConnectTimeout(time.Duration(client.config.ClientConnectTimeout) * time.Millisecond)
		dcdnConn.AppendUserAgent(Terraform, terraformVersion)
		dcdnConn.AppendUserAgent(Provider, providerVersion)
		dcdnConn.AppendUserAgent(Module, client.config.ConfigurationSource)
		client.dcdnConn = dcdnConn
	}
	return do(client.dcdnConn)
}

func (client *AliyunClient) WithRKvstoreClient(do func(*r_kvstore.Client) (interface{}, error)) (interface{}, error) {
	if client.r_kvstoreConn == nil {
		endpoint := client.config.RKvstoreEndpoint
		if endpoint == "" {
			endpoint = loadEndpoint(client.config.RegionId, RKvstoreCode)
		}
		if strings.HasPrefix(endpoint, "http") {
			endpoint = fmt.Sprintf("https://%s", strings.TrimPrefix(endpoint, "http://"))
		}
		if endpoint != "" {
			endpoints.AddEndpointMapping(client.config.RegionId, string(RKvstoreCode), endpoint)
		}

		r_kvstoreConn, err := r_kvstore.NewClientWithOptions(client.config.RegionId, client.getSdkConfig(), client.config.getAuthCredential(true))
		if err != nil {
			return nil, fmt.Errorf("unable to initialize the RKvstoreclient: %#v", err)
		}
		r_kvstoreConn.SetReadTimeout(time.Duration(client.config.ClientReadTimeout) * time.Millisecond)
		r_kvstoreConn.SetConnectTimeout(time.Duration(client.config.ClientConnectTimeout) * time.Millisecond)
		r_kvstoreConn.AppendUserAgent(Terraform, terraformVersion)
		r_kvstoreConn.AppendUserAgent(Provider, providerVersion)
		r_kvstoreConn.AppendUserAgent(Module, client.config.ConfigurationSource)
		client.r_kvstoreConn = r_kvstoreConn
	}
	return do(client.r_kvstoreConn)
}

func (client *AliyunClient) NewOnsClient() (*rpc.Client, error) {
	productCode := "ons"
	endpoint := ""
	if v, ok := client.config.Endpoints[productCode]; !ok || v.(string) == "" {
		if err := client.loadEndpoint(productCode); err != nil {
			return nil, err
		}
	}
	if v, ok := client.config.Endpoints[productCode]; ok && v.(string) != "" {
		endpoint = v.(string)
	}
	if endpoint == "" {
		return nil, fmt.Errorf("[ERROR] missing the product %s endpoint.", productCode)
	}
	sdkConfig := client.teaSdkConfig
	sdkConfig.SetEndpoint(endpoint)
	conn, err := rpc.NewClient(&sdkConfig)
	if err != nil {
		return nil, fmt.Errorf("unable to initialize the %s client: %#v", productCode, err)
	}
	return conn, nil
}

func (client *AliyunClient) NewCmsClient() (*rpc.Client, error) {
	productCode := "cms"
	endpoint := ""
	if v, ok := client.config.Endpoints[productCode]; !ok || v.(string) == "" {
		if err := client.loadEndpoint(productCode); err != nil {
			return nil, err
		}
	}
	if v, ok := client.config.Endpoints[productCode]; ok && v.(string) != "" {
		endpoint = v.(string)
	}
	if endpoint == "" {
		return nil, fmt.Errorf("[ERROR] missing the product %s endpoint.", productCode)
	}
	sdkConfig := client.teaSdkConfig
	sdkConfig.SetEndpoint(endpoint)
	conn, err := rpc.NewClient(&sdkConfig)
	if err != nil {
		return nil, fmt.Errorf("unable to initialize the %s client: %#v", productCode, err)
	}
	return conn, nil
}

func (client *AliyunClient) NewConfigClient() (*rpc.Client, error) {
	productCode := "config"
	endpoint := ""
	if v, ok := client.config.Endpoints[productCode]; !ok || v.(string) == "" {
		if err := client.loadEndpoint(productCode); err != nil {
			endpoint = "config.cn-shanghai.aliyuncs.com"
			client.config.Endpoints[productCode] = endpoint
			log.Printf("[ERROR] loading %s endpoint got an error: %#v. Using the central endpoint %s instead.", productCode, err, endpoint)
		}
	}
	if v, ok := client.config.Endpoints[productCode]; ok && v.(string) != "" {
		endpoint = v.(string)
	}
	if endpoint == "" {
		return nil, fmt.Errorf("[ERROR] missing the product %s endpoint.", productCode)
	}
	sdkConfig := client.teaSdkConfig
	sdkConfig.SetEndpoint(endpoint)
	conn, err := rpc.NewClient(&sdkConfig)
	if err != nil {
		return nil, fmt.Errorf("unable to initialize the %s client: %#v", productCode, err)
	}
	return conn, nil
}

func (client *AliyunClient) NewWafClient() (*rpc.Client, error) {
	productCode := "waf"
	endpoint := ""
	if v, ok := client.config.Endpoints[productCode]; !ok || v.(string) == "" {
		if err := client.loadEndpoint(productCode); err != nil {
			return nil, err
		}
	}
	if v, ok := client.config.Endpoints[productCode]; ok && v.(string) != "" {
		endpoint = v.(string)
	}
	if endpoint == "" {
		return nil, fmt.Errorf("[ERROR] missing the product %s endpoint.", productCode)
	}
	sdkConfig := client.teaSdkConfig
	sdkConfig.SetEndpoint(endpoint)
	conn, err := rpc.NewClient(&sdkConfig)
	if err != nil {
		return nil, fmt.Errorf("unable to initialize the %s client: %#v", productCode, err)
	}
	return conn, nil
}
func (client *AliyunClient) NewBssopenapiClient() (*rpc.Client, error) {
	productCode := "bssopenapi"
	// bss endpoint depends on the account type.
	// Domestic account is business.aliyuncs.com and International account is business.ap-southeast-1.aliyuncs.com
	endpoint := BssOpenAPIEndpointDomestic
	sdkConfig := client.teaSdkConfig
	sdkConfig.SetEndpoint(endpoint)
	conn, err := rpc.NewClient(&sdkConfig)
	if err != nil {
		return nil, fmt.Errorf("unable to initialize the %s client: %#v", productCode, err)
	}
	return conn, nil
}

func (client *AliyunClient) NewFnfClient() (*rpc.Client, error) {
	productCode := "fnf"
	endpoint := ""
	if v, ok := client.config.Endpoints[productCode]; !ok || v.(string) == "" {
		if err := client.loadEndpoint(productCode); err != nil {
			return nil, err
		}
	}
	if v, ok := client.config.Endpoints[productCode]; ok && v.(string) != "" {
		endpoint = v.(string)
	}
	if endpoint == "" {
		return nil, fmt.Errorf("[ERROR] missing the product %s endpoint.", productCode)
	}
	sdkConfig := client.teaSdkConfig
	sdkConfig.SetEndpoint(endpoint)
	conn, err := rpc.NewClient(&sdkConfig)
	if err != nil {
		return nil, fmt.Errorf("unable to initialize the %s client: %#v", productCode, err)
	}
	return conn, nil
}

func (client *AliyunClient) NewRosClient() (*rpc.Client, error) {
	productCode := "ros"
	endpoint := ""
	if v, ok := client.config.Endpoints[productCode]; !ok || v.(string) == "" {
		if err := client.loadEndpoint(productCode); err != nil {
			return nil, err
		}
	}
	if v, ok := client.config.Endpoints[productCode]; ok && v.(string) != "" {
		endpoint = v.(string)
	}
	if endpoint == "" {
		return nil, fmt.Errorf("[ERROR] missing the product %s endpoint.", productCode)
	}
	sdkConfig := client.teaSdkConfig
	sdkConfig.SetEndpoint(endpoint)
	conn, err := rpc.NewClient(&sdkConfig)
	if err != nil {
		return nil, fmt.Errorf("unable to initialize the %s client: %#v", productCode, err)
	}
	return conn, nil
}

func (client *AliyunClient) NewPvtzClient() (*rpc.Client, error) {
	productCode := "pvtz"
	endpoint := ""
	if v, ok := client.config.Endpoints[productCode]; !ok || v.(string) == "" {
		if err := client.loadEndpoint(productCode); err != nil {
			endpoint := "pvtz.aliyuncs.com"
			client.config.Endpoints[productCode] = endpoint
			log.Printf("[ERROR] loading %s endpoint got an error: %#v. Using the central endpoint %s instead.", productCode, err, endpoint)
		}
	}
	if v, ok := client.config.Endpoints[productCode]; ok && v.(string) != "" {
		endpoint = v.(string)
	}
	if endpoint == "" {
		return nil, fmt.Errorf("[ERROR] missing the product %s endpoint.", productCode)
	}
	sdkConfig := client.teaSdkConfig
	sdkConfig.SetEndpoint(endpoint)
	conn, err := rpc.NewClient(&sdkConfig)
	if err != nil {
		return nil, fmt.Errorf("unable to initialize the %s client: %#v", productCode, err)
	}
	return conn, nil
}

func (client *AliyunClient) NewPrivatelinkClient() (*rpc.Client, error) {
	productCode := "privatelink"
	endpoint := ""
	if v, ok := client.config.Endpoints[productCode]; !ok || v.(string) == "" {
		if err := client.loadEndpoint(productCode); err != nil {
			return nil, err
		}
	}
	if v, ok := client.config.Endpoints[productCode]; ok && v.(string) != "" {
		endpoint = v.(string)
	}
	if endpoint == "" {
		return nil, fmt.Errorf("[ERROR] missing the product %s endpoint.", productCode)
	}
	sdkConfig := client.teaSdkConfig
	sdkConfig.SetEndpoint(endpoint)
	conn, err := rpc.NewClient(&sdkConfig)
	if err != nil {
		return nil, fmt.Errorf("unable to initialize the %s client: %#v", productCode, err)
	}
	return conn, nil
}

func (client *AliyunClient) NewDcdnClient() (*rpc.Client, error) {
	productCode := "dcdn"
	endpoint := ""
	if v, ok := client.config.Endpoints[productCode]; !ok || v.(string) == "" {
		if err := client.loadEndpoint(productCode); err != nil {
			endpoint = "dcdn.aliyuncs.com"
			client.config.Endpoints[productCode] = endpoint
			log.Printf("[ERROR] loading %s endpoint got an error: %#v. Using the central endpoint %s instead.", productCode, err, endpoint)
		}
	}
	if v, ok := client.config.Endpoints[productCode]; ok && v.(string) != "" {
		endpoint = v.(string)
	}
	if endpoint == "" {
		return nil, fmt.Errorf("[ERROR] missing the product %s endpoint.", productCode)
	}
	sdkConfig := client.teaSdkConfig
	sdkConfig.SetEndpoint(endpoint)
	conn, err := rpc.NewClient(&sdkConfig)
	if err != nil {
		return nil, fmt.Errorf("unable to initialize the %s client: %#v", productCode, err)
	}
	return conn, nil
}

func (client *AliyunClient) NewOdpsClient() (*rpc.Client, error) {
	productCode := "odps"
	endpoint := ""
	if v, ok := client.config.Endpoints[productCode]; !ok || v.(string) == "" {
		if err := client.loadEndpoint(productCode); err != nil {
			return nil, err
		}
	}
	if v, ok := client.config.Endpoints[productCode]; ok && v.(string) != "" {
		endpoint = v.(string)
	}
	if endpoint == "" {
		return nil, fmt.Errorf("[ERROR] missing the product %s endpoint.", productCode)
	}
	sdkConfig := client.teaSdkConfig
	sdkConfig.SetEndpoint(endpoint)
	conn, err := rpc.NewClient(&sdkConfig)
	if err != nil {
		return nil, fmt.Errorf("unable to initialize the %s client: %#v", productCode, err)
	}
	return conn, nil
}

func (client *AliyunClient) NewRessharingClient() (*rpc.Client, error) {
	productCode := "ressharing"
	endpoint := ""
	if v, ok := client.config.Endpoints[productCode]; !ok || v.(string) == "" {
		if err := client.loadEndpoint(productCode); err != nil {
			return nil, err
		}
	}
	if v, ok := client.config.Endpoints[productCode]; ok && v.(string) != "" {
		endpoint = v.(string)
	}
	if endpoint == "" {
		return nil, fmt.Errorf("[ERROR] missing the product %s endpoint.", productCode)
	}
	sdkConfig := client.teaSdkConfig
	sdkConfig.SetEndpoint(endpoint)
	conn, err := rpc.NewClient(&sdkConfig)
	if err != nil {
		return nil, fmt.Errorf("unable to initialize the %s client: %#v", productCode, err)
	}
	return conn, nil
}

func (client *AliyunClient) NewGaplusClient() (*rpc.Client, error) {
	productCode := "gaplus"
	endpoint := ""
	if v, ok := client.config.Endpoints[productCode]; !ok || v.(string) == "" {
		if err := client.loadEndpoint(productCode); err != nil {
			return nil, err
		}
	}
	if v, ok := client.config.Endpoints[productCode]; ok && v.(string) != "" {
		endpoint = v.(string)
	}
	if endpoint == "" {
		return nil, fmt.Errorf("[ERROR] missing the product %s endpoint.", productCode)
	}
	sdkConfig := client.teaSdkConfig
	sdkConfig.SetEndpoint(endpoint)
	conn, err := rpc.NewClient(&sdkConfig)
	if err != nil {
		return nil, fmt.Errorf("unable to initialize the %s client: %#v", productCode, err)
	}
	return conn, nil
}

func (client *AliyunClient) NewEciClient() (*rpc.Client, error) {
	productCode := "eci"
	endpoint := ""
	if v, ok := client.config.Endpoints[productCode]; !ok || v.(string) == "" {
		if err := client.loadEndpoint(productCode); err != nil {
			return nil, err
		}
	}
	if v, ok := client.config.Endpoints[productCode]; ok && v.(string) != "" {
		endpoint = v.(string)
	}
	if endpoint == "" {
		return nil, fmt.Errorf("[ERROR] missing the product %s endpoint.", productCode)
	}
	sdkConfig := client.teaSdkConfig
	sdkConfig.SetEndpoint(endpoint)
	conn, err := rpc.NewClient(&sdkConfig)
	if err != nil {
		return nil, fmt.Errorf("unable to initialize the %s client: %#v", productCode, err)
	}
	return conn, nil
}

func (client *AliyunClient) NewActiontrailClient() (*rpc.Client, error) {
	productCode := "actiontrail"
	endpoint := ""
	if v, ok := client.config.Endpoints[productCode]; !ok || v.(string) == "" {
		if err := client.loadEndpoint(productCode); err != nil {
			return nil, err
		}
	}
	if v, ok := client.config.Endpoints[productCode]; ok && v.(string) != "" {
		endpoint = v.(string)
	}
	if endpoint == "" {
		return nil, fmt.Errorf("[ERROR] missing the product %s endpoint.", productCode)
	}
	sdkConfig := client.teaSdkConfig
	sdkConfig.SetEndpoint(endpoint)
	conn, err := rpc.NewClient(&sdkConfig)
	if err != nil {
		return nil, fmt.Errorf("unable to initialize the %s client: %#v", productCode, err)
	}
	return conn, nil
}

func (client *AliyunClient) NewMseClient() (*rpc.Client, error) {
	productCode := "mse"
	endpoint := ""
	if client.config.Endpoints[productCode] == nil {
		if err := client.loadEndpoint(productCode); err != nil {
			return nil, err
		}
	}
	if client.config.Endpoints[productCode] != nil && client.config.Endpoints[productCode].(string) != "" {
		endpoint = client.config.Endpoints[productCode].(string)
	}
	if endpoint == "" {
		return nil, fmt.Errorf("[ERROR] missing the product %s endpoint.", productCode)
	}
	sdkConfig := client.teaSdkConfig
	sdkConfig.SetEndpoint(endpoint)
	conn, err := rpc.NewClient(&sdkConfig)
	if err != nil {
		return nil, fmt.Errorf("unable to initialize the %s client: %#v", productCode, err)
	}
	return conn, nil
}

func (client *AliyunClient) NewHitsdbClient() (*rpc.Client, error) {
	productCode := "hitsdb"
	endpoint := ""
	if v, ok := client.config.Endpoints[productCode]; !ok || v.(string) == "" {
		if err := client.loadEndpoint(productCode); err != nil {
			return nil, err
		}
	}
	if v, ok := client.config.Endpoints[productCode]; ok && v.(string) != "" {
		endpoint = v.(string)
	}
	if endpoint == "" {
		return nil, fmt.Errorf("[ERROR] missing the product %s endpoint.", productCode)
	}
	sdkConfig := client.teaSdkConfig
	sdkConfig.SetEndpoint(endpoint)
	conn, err := rpc.NewClient(&sdkConfig)
	if err != nil {
		return nil, fmt.Errorf("unable to initialize the %s client: %#v", productCode, err)
	}
	return conn, nil
}

func (client *AliyunClient) NewAistudioClient() (*rpc.Client, error) {
	productCode := "aistudio"
	endpoint := ""
	if v, ok := client.config.Endpoints[productCode]; !ok || v.(string) == "" {
		if err := client.loadEndpoint(productCode); err != nil {
			endpoint = "brain-industrial.cn-hangzhou.aliyuncs.com"
			client.config.Endpoints[productCode] = endpoint
			log.Printf("[ERROR] loading %s endpoint got an error: %#v. Using the central endpoint %s instead.", productCode, err, endpoint)
		}
	}
	if v, ok := client.config.Endpoints[productCode]; ok && v.(string) != "" {
		endpoint = v.(string)
	}
	if endpoint == "" {
		return nil, fmt.Errorf("[ERROR] missing the product %s endpoint.", productCode)
	}
	sdkConfig := client.teaSdkConfig
	sdkConfig.SetEndpoint(endpoint)
	conn, err := rpc.NewClient(&sdkConfig)
	if err != nil {
		return nil, fmt.Errorf("unable to initialize the %s client: %#v", productCode, err)
	}
	return conn, nil
}

func (client *AliyunClient) NewEipanycastClient() (*rpc.Client, error) {
	productCode := "eipanycast"
	endpoint := ""
	if v, ok := client.config.Endpoints[productCode]; !ok || v.(string) == "" {
		if err := client.loadEndpoint(productCode); err != nil {
			return nil, err
		}
	}
	if v, ok := client.config.Endpoints[productCode]; ok && v.(string) != "" {
		endpoint = v.(string)
	}
	if endpoint == "" {
		return nil, fmt.Errorf("[ERROR] missing the product %s endpoint.", productCode)
	}
	sdkConfig := client.teaSdkConfig
	sdkConfig.SetEndpoint(endpoint)
	conn, err := rpc.NewClient(&sdkConfig)
	if err != nil {
		return nil, fmt.Errorf("unable to initialize the %s client: %#v", productCode, err)
	}
	return conn, nil
}

func (client *AliyunClient) NewOosClient() (*rpc.Client, error) {
	productCode := "oos"
	endpoint := ""
	if v, ok := client.config.Endpoints[productCode]; !ok || v.(string) == "" {
		if err := client.loadEndpoint(productCode); err != nil {
			return nil, err
		}
	}
	if v, ok := client.config.Endpoints[productCode]; ok && v.(string) != "" {
		endpoint = v.(string)
	}
	if endpoint == "" {
		return nil, fmt.Errorf("[ERROR] missing the product %s endpoint.", productCode)
	}
	sdkConfig := client.teaSdkConfig
	sdkConfig.SetEndpoint(endpoint)
	conn, err := rpc.NewClient(&sdkConfig)
	if err != nil {
		return nil, fmt.Errorf("unable to initialize the %s client: %#v", productCode, err)
	}
	return conn, nil
}

func (client *AliyunClient) NewImsClient() (*rpc.Client, error) {
	productCode := "ims"
	endpoint := ""
	if v, ok := client.config.Endpoints[productCode]; !ok || v.(string) == "" {
		if err := client.loadEndpoint(productCode); err != nil {
			endpoint = "ims.aliyuncs.com"
			client.config.Endpoints[productCode] = endpoint
			log.Printf("[ERROR] loading %s endpoint got an error: %#v. Using the central endpoint %s instead.", productCode, err, endpoint)
		}
	}
	if v, ok := client.config.Endpoints[productCode]; ok && v.(string) != "" {
		endpoint = v.(string)
	}
	if endpoint == "" {
		return nil, fmt.Errorf("[ERROR] missing the product %s endpoint.", productCode)
	}
	sdkConfig := client.teaSdkConfig
	sdkConfig.SetEndpoint(endpoint)
	conn, err := rpc.NewClient(&sdkConfig)
	if err != nil {
		return nil, fmt.Errorf("unable to initialize the %s client: %#v", productCode, err)
	}
	return conn, nil
}

func (client *AliyunClient) NewRamClient() (*rpc.Client, error) {
	productCode := "ram"
	endpoint := ""
	if v, ok := client.config.Endpoints[productCode]; !ok || v.(string) == "" {
		if err := client.loadEndpoint(productCode); err != nil {
			endpoint = "ram.aliyuncs.com"
			client.config.Endpoints[productCode] = endpoint
			log.Printf("[ERROR] loading %s endpoint got an error: %#v. Using the central endpoint %s instead.", productCode, err, endpoint)
		}
	}
	if v, ok := client.config.Endpoints[productCode]; ok && v.(string) != "" {
		endpoint = v.(string)
	}
	if endpoint == "" {
		return nil, fmt.Errorf("[ERROR] missing the product %s endpoint.", productCode)
	}
	sdkConfig := client.teaSdkConfig
	sdkConfig.SetEndpoint(endpoint)
	conn, err := rpc.NewClient(&sdkConfig)
	if err != nil {
		return nil, fmt.Errorf("unable to initialize the %s client: %#v", productCode, err)
	}
	return conn, nil
}

func (client *AliyunClient) NewResourcemanagerClient() (*rpc.Client, error) {
	productCode := "resourcemanager"
	endpoint := ""
	if v, ok := client.config.Endpoints[productCode]; !ok || v.(string) == "" {
		if err := client.loadEndpoint(productCode); err != nil {
			endpoint = "resourcemanager.aliyuncs.com"
			client.config.Endpoints[productCode] = endpoint
			log.Printf("[ERROR] loading %s endpoint got an error: %#v. Using the central endpoint %s instead.", productCode, err, endpoint)
		}
	}
	if v, ok := client.config.Endpoints[productCode]; ok && v.(string) != "" {
		endpoint = v.(string)
	}
	if endpoint == "" {
		return nil, fmt.Errorf("[ERROR] missing the product %s endpoint.", productCode)
	}
	sdkConfig := client.teaSdkConfig
	sdkConfig.SetEndpoint(endpoint)
	conn, err := rpc.NewClient(&sdkConfig)
	if err != nil {
		return nil, fmt.Errorf("unable to initialize the %s client: %#v", productCode, err)
	}
	return conn, nil
}

func (client *AliyunClient) NewQuotasClient() (*rpc.Client, error) {
	productCode := "quotas"
	endpoint := ""
	if v, ok := client.config.Endpoints[productCode]; !ok || v.(string) == "" {
		if err := client.loadEndpoint(productCode); err != nil {
			endpoint = "quotas.aliyuncs.com"
			client.config.Endpoints[productCode] = endpoint
			log.Printf("[ERROR] loading %s endpoint got an error: %#v. Using the central endpoint %s instead.", productCode, err, endpoint)
		}
	}
	if v, ok := client.config.Endpoints[productCode]; ok && v.(string) != "" {
		endpoint = v.(string)
	}
	if endpoint == "" {
		return nil, fmt.Errorf("[ERROR] missing the product %s endpoint.", productCode)
	}
	sdkConfig := client.teaSdkConfig
	sdkConfig.SetEndpoint(endpoint)
	conn, err := rpc.NewClient(&sdkConfig)
	if err != nil {
		return nil, fmt.Errorf("unable to initialize the %s client: %#v", productCode, err)
	}
	return conn, nil
}

func (client *AliyunClient) NewNasClient() (*rpc.Client, error) {
	productCode := "nas"
	endpoint := ""
	if v, ok := client.config.Endpoints[productCode]; !ok || v.(string) == "" {
		if err := client.loadEndpoint(productCode); err != nil {
			return nil, err
		}
	}
	if v, ok := client.config.Endpoints[productCode]; ok && v.(string) != "" {
		endpoint = v.(string)
	}
	if endpoint == "" {
		return nil, fmt.Errorf("[ERROR] missing the product %s endpoint.", productCode)
	}
	sdkConfig := client.teaSdkConfig
	sdkConfig.SetEndpoint(endpoint)
	conn, err := rpc.NewClient(&sdkConfig)
	if err != nil {
		return nil, fmt.Errorf("unable to initialize the %s client: %#v", productCode, err)
	}
	return conn, nil
}

func (client *AliyunClient) NewDmsenterpriseClient() (*rpc.Client, error) {
	productCode := "dmsenterprise"
	endpoint := ""
	if v, ok := client.config.Endpoints[productCode]; !ok || v.(string) == "" {
		if err := client.loadEndpoint(productCode); err != nil {
			endpoint = "dms-enterprise.aliyuncs.com"
			client.config.Endpoints[productCode] = endpoint
			log.Printf("[ERROR] loading %s endpoint got an error: %#v. Using the central endpoint %s instead.", productCode, err, endpoint)
		}
	}
	if v, ok := client.config.Endpoints[productCode]; ok && v.(string) != "" {
		endpoint = v.(string)
	}
	if endpoint == "" {
		return nil, fmt.Errorf("[ERROR] missing the product %s endpoint.", productCode)
	}
	sdkConfig := client.teaSdkConfig
	sdkConfig.SetEndpoint(endpoint)
	conn, err := rpc.NewClient(&sdkConfig)
	if err != nil {
		return nil, fmt.Errorf("unable to initialize the %s client: %#v", productCode, err)
	}
	return conn, nil
}

func (client *AliyunClient) NewHcsSgwClient() (*rpc.Client, error) {
	productCode := "hcs_sgw"
	endpoint := ""
	if v, ok := client.config.Endpoints[productCode]; !ok || v.(string) == "" {
		if err := client.loadEndpoint(productCode); err != nil {
			return nil, err
		}
	}
	if v, ok := client.config.Endpoints[productCode]; ok && v.(string) != "" {
		endpoint = v.(string)
	}
	if endpoint == "" {
		return nil, fmt.Errorf("[ERROR] missing the product %s endpoint.", productCode)
	}
	sdkConfig := client.teaSdkConfig
	sdkConfig.SetEndpoint(endpoint)
	conn, err := rpc.NewClient(&sdkConfig)
	if err != nil {
		return nil, fmt.Errorf("unable to initialize the %s client: %#v", productCode, err)
	}
	return conn, nil
}

func (client *AliyunClient) NewAdsClient() (*rpc.Client, error) {
	productCode := "ads"
	endpoint := ""
	if v, ok := client.config.Endpoints[productCode]; !ok || v.(string) == "" {
		if err := client.loadEndpoint(productCode); err != nil {
			return nil, err
		}
	}
	if v, ok := client.config.Endpoints[productCode]; ok && v.(string) != "" {
		endpoint = v.(string)
	}
	if endpoint == "" {
		return nil, fmt.Errorf("[ERROR] missing the product %s endpoint.", productCode)
	}
	sdkConfig := client.teaSdkConfig
	sdkConfig.SetEndpoint(endpoint)
	conn, err := rpc.NewClient(&sdkConfig)
	if err != nil {
		return nil, fmt.Errorf("unable to initialize the %s client: %#v", productCode, err)
	}
	return conn, nil
}

func (client *AliyunClient) NewDdoscooClient() (*rpc.Client, error) {
	productCode := "ddoscoo"
	endpoint := ""
	if v, ok := client.config.Endpoints[productCode]; !ok || v.(string) == "" {
		if err := client.loadEndpoint(productCode); err != nil {
			return nil, err
		}
	}
	if v, ok := client.config.Endpoints[productCode]; ok && v.(string) != "" {
		endpoint = v.(string)
	}
	if endpoint == "" {
		return nil, fmt.Errorf("[ERROR] missing the product %s endpoint.", productCode)
	}
	sdkConfig := client.teaSdkConfig
	sdkConfig.SetEndpoint(endpoint)
	conn, err := rpc.NewClient(&sdkConfig)
	if err != nil {
		return nil, fmt.Errorf("unable to initialize the %s client: %#v", productCode, err)
	}
	return conn, nil
}

func (client *AliyunClient) NewSlbClient() (*rpc.Client, error) {
	productCode := "slb"
	endpoint := ""
	if v, ok := client.config.Endpoints[productCode]; !ok || v.(string) == "" {
		if err := client.loadEndpoint(productCode); err != nil {
			return nil, err
		}
	}
	if v, ok := client.config.Endpoints[productCode]; ok && v.(string) != "" {
		endpoint = v.(string)
	}
	if endpoint == "" {
		return nil, fmt.Errorf("[ERROR] missing the product %s endpoint.", productCode)
	}
	sdkConfig := client.teaSdkConfig
	sdkConfig.SetEndpoint(endpoint)
	conn, err := rpc.NewClient(&sdkConfig)
	if err != nil {
		return nil, fmt.Errorf("unable to initialize the %s client: %#v", productCode, err)
	}
	return conn, nil
}

func (client *AliyunClient) NewKmsClient() (*rpc.Client, error) {
	productCode := "kms"
	endpoint := ""
	if v, ok := client.config.Endpoints[productCode]; !ok || v.(string) == "" {
		if err := client.loadEndpoint(productCode); err != nil {
			endpoint = "kms.cn-beijing.aliyuncs.com"
			client.config.Endpoints[productCode] = endpoint
			log.Printf("[ERROR] loading %s endpoint got an error: %#v. Using the central endpoint %s instead.", productCode, err, endpoint)
		}
	}
	if v, ok := client.config.Endpoints[productCode]; ok && v.(string) != "" {
		endpoint = v.(string)
	}
	if endpoint == "" {
		return nil, fmt.Errorf("[ERROR] missing the product %s endpoint.", productCode)
	}
	sdkConfig := client.teaSdkConfig
	sdkConfig.SetEndpoint(endpoint)
	conn, err := rpc.NewClient(&sdkConfig)
	if err != nil {
		return nil, fmt.Errorf("unable to initialize the %s client: %#v", productCode, err)
	}
	return conn, nil
}

func (client *AliyunClient) NewAlidnsClient() (*rpc.Client, error) {
	productCode := "alidns"
	endpoint := ""
	if v, ok := client.config.Endpoints[productCode]; !ok || v.(string) == "" {
		if err := client.loadEndpoint(productCode); err != nil {
			endpoint = "alidns.aliyuncs.com"
			client.config.Endpoints[productCode] = endpoint
			log.Printf("[ERROR] loading %s endpoint got an error: %#v. Using the central endpoint %s instead.", productCode, err, endpoint)
		}
	}
	if v, ok := client.config.Endpoints[productCode]; ok && v.(string) != "" {
		endpoint = v.(string)
	}
	if endpoint == "" {
		return nil, fmt.Errorf("[ERROR] missing the product %s endpoint.", productCode)
	}
	sdkConfig := client.teaSdkConfig
	sdkConfig.SetEndpoint(endpoint)
	conn, err := rpc.NewClient(&sdkConfig)
	if err != nil {
		return nil, fmt.Errorf("unable to initialize the %s client: %#v", productCode, err)
	}
	return conn, nil
}

func (client *AliyunClient) NewHbaseClient() (*rpc.Client, error) {
	productCode := "hbase"
	endpoint := ""
	if v, ok := client.config.Endpoints[productCode]; !ok || v.(string) == "" {
		if err := client.loadEndpoint(productCode); err != nil {
			return nil, err
		}
	}
	if v, ok := client.config.Endpoints[productCode]; ok && v.(string) != "" {
		endpoint = v.(string)
	}
	if endpoint == "" {
		return nil, fmt.Errorf("[ERROR] missing the product %s endpoint.", productCode)
	}
	sdkConfig := client.teaSdkConfig
	sdkConfig.SetEndpoint(endpoint)
	conn, err := rpc.NewClient(&sdkConfig)
	if err != nil {
		return nil, fmt.Errorf("unable to initialize the %s client: %#v", productCode, err)
	}
	return conn, nil
}

func (client *AliyunClient) NewDmClient() (*rpc.Client, error) {
	productCode := "dm"
	endpoint := ""
	if v, ok := client.config.Endpoints[productCode]; !ok || v.(string) == "" {
		if err := client.loadEndpoint(productCode); err != nil {
			return nil, err
		}
	}
	if v, ok := client.config.Endpoints[productCode]; ok && v.(string) != "" {
		endpoint = v.(string)
	}
	if endpoint == "" {
		return nil, fmt.Errorf("[ERROR] missing the product %s endpoint.", productCode)
	}
	sdkConfig := client.teaSdkConfig
	sdkConfig.SetEndpoint(endpoint)
	conn, err := rpc.NewClient(&sdkConfig)
	if err != nil {
		return nil, fmt.Errorf("unable to initialize the %s client: %#v", productCode, err)
	}
	return conn, nil
}

func (client *AliyunClient) NewEventbridgeClient() (*rpc.Client, error) {
	productCode := "eventbridge"
	endpoint := ""
	if v, ok := client.config.Endpoints[productCode]; !ok || v.(string) == "" {
		if err := client.loadEndpoint(productCode); err != nil {
			return nil, err
		}
	}
	if v, ok := client.config.Endpoints[productCode]; ok && v.(string) != "" {
		endpoint = v.(string)
	}
	if endpoint == "" {
		return nil, fmt.Errorf("[ERROR] missing the product %s endpoint.", productCode)
	}
	sdkConfig := client.teaSdkConfig
	sdkConfig.SetEndpoint(endpoint)
	conn, err := rpc.NewClient(&sdkConfig)
	if err != nil {
		return nil, fmt.Errorf("unable to initialize the %s client: %#v", productCode, err)
	}
	return conn, nil
}

func (client *AliyunClient) NewOnsproxyClient() (*rpc.Client, error) {
	productCode := "onsproxy"
	endpoint := ""
	if v, ok := client.config.Endpoints[productCode]; !ok || v.(string) == "" {
		if err := client.loadEndpoint(productCode); err != nil {
			return nil, err
		}
	}
	if v, ok := client.config.Endpoints[productCode]; ok && v.(string) != "" {
		endpoint = v.(string)
	}
	if endpoint == "" {
		return nil, fmt.Errorf("[ERROR] missing the product %s endpoint.", productCode)
	}
	sdkConfig := client.teaSdkConfig
	sdkConfig.SetEndpoint(endpoint)
	conn, err := rpc.NewClient(&sdkConfig)
	if err != nil {
		return nil, fmt.Errorf("unable to initialize the %s client: %#v", productCode, err)
	}
	return conn, nil
}

func (client *AliyunClient) NewCdsClient() (*rpc.Client, error) {
	productCode := "cds"
	endpoint := ""
	if v, ok := client.config.Endpoints[productCode]; !ok || v.(string) == "" {
		if err := client.loadEndpoint(productCode); err != nil {
			endpoint = "cassandra.aliyuncs.com"
			client.config.Endpoints[productCode] = endpoint
			log.Printf("[ERROR] loading %s endpoint got an error: %#v. Using the central endpoint %s instead.", productCode, err, endpoint)
		}
	}
	if v, ok := client.config.Endpoints[productCode]; ok && v.(string) != "" {
		endpoint = v.(string)
	}
	if endpoint == "" {
		return nil, fmt.Errorf("[ERROR] missing the product %s endpoint.", productCode)
	}
	sdkConfig := client.teaSdkConfig
	sdkConfig.SetEndpoint(endpoint)
	conn, err := rpc.NewClient(&sdkConfig)
	if err != nil {
		return nil, fmt.Errorf("unable to initialize the %s client: %#v", productCode, err)
	}
	return conn, nil
}

func (client *AliyunClient) NewHbrClient() (*rpc.Client, error) {
	productCode := "hbr"
	endpoint := ""
	if v, ok := client.config.Endpoints[productCode]; !ok || v.(string) == "" {
		if err := client.loadEndpoint(productCode); err != nil {
			endpoint = fmt.Sprintf("hbr.%s.aliyuncs.com", client.config.RegionId)
			client.config.Endpoints[productCode] = endpoint
			log.Printf("[ERROR] loading %s endpoint got an error: %#v. Using the endpoint %s instead.", productCode, err, endpoint)
		}
	}
	if v, ok := client.config.Endpoints[productCode]; ok && v.(string) != "" {
		endpoint = v.(string)
	}
	if endpoint == "" {
		return nil, fmt.Errorf("[ERROR] missing the product %s endpoint.", productCode)
	}
	sdkConfig := client.teaSdkConfig
	sdkConfig.SetEndpoint(endpoint)
	conn, err := rpc.NewClient(&sdkConfig)
	if err != nil {
		return nil, fmt.Errorf("unable to initialize the %s client: %#v", productCode, err)
	}
	return conn, nil
}

func (client *AliyunClient) NewCasClient() (*rpc.Client, error) {
	productCode := "cas"
	endpoint := ""
	if v, ok := client.config.Endpoints[productCode]; !ok || v.(string) == "" {
		if err := client.loadEndpoint(productCode); err != nil {
			endpoint = "cas.aliyuncs.com"
			client.config.Endpoints[productCode] = endpoint
			log.Printf("[ERROR] loading %s endpoint got an error: %#v. Using the endpoint %s instead.", productCode, err, endpoint)
		}
	}
	if v, ok := client.config.Endpoints[productCode]; ok && v.(string) != "" {
		endpoint = v.(string)
	}
	if endpoint == "" {
		return nil, fmt.Errorf("[ERROR] missing the product %s endpoint.", productCode)
	}
	sdkConfig := client.teaSdkConfig
	sdkConfig.SetEndpoint(endpoint)
	conn, err := rpc.NewClient(&sdkConfig)
	if err != nil {
		return nil, fmt.Errorf("unable to initialize the %s client: %#v", productCode, err)
	}
	return conn, nil
}

func (client *AliyunClient) NewArmsClient() (*rpc.Client, error) {
	productCode := "arms"
	endpoint := ""
	if v, ok := client.config.Endpoints[productCode]; !ok || v.(string) == "" {
		if err := client.loadEndpoint(productCode); err != nil {
			endpoint = fmt.Sprintf("arms.%s.aliyuncs.com", client.config.RegionId)
			client.config.Endpoints[productCode] = endpoint
			log.Printf("[ERROR] loading %s endpoint got an error: %#v. Using the endpoint %s instead.", productCode, err, endpoint)
		}
	}
	if v, ok := client.config.Endpoints[productCode]; ok && v.(string) != "" {
		endpoint = v.(string)
	}
	if endpoint == "" {
		return nil, fmt.Errorf("[ERROR] missing the product %s endpoint.", productCode)
	}
	sdkConfig := client.teaSdkConfig
	sdkConfig.SetEndpoint(endpoint)
	conn, err := rpc.NewClient(&sdkConfig)
	if err != nil {
		return nil, fmt.Errorf("unable to initialize the %s client: %#v", productCode, err)
	}
	return conn, nil
}

func (client *AliyunClient) NewCloudfwClient() (*rpc.Client, error) {
	productCode := "cloudfw"
	endpoint := ""
	if v, ok := client.config.Endpoints[productCode]; !ok || v.(string) == "" {
		if err := client.loadEndpoint(productCode); err != nil {
			endpoint = "cloudfw.aliyuncs.com"
			client.config.Endpoints[productCode] = endpoint
			log.Printf("[ERROR] loading %s endpoint got an error: %#v. Using the endpoint %s instead.", productCode, err, endpoint)
		}
	}
	if v, ok := client.config.Endpoints[productCode]; ok && v.(string) != "" {
		endpoint = v.(string)
	}
	if endpoint == "" {
		return nil, fmt.Errorf("[ERROR] missing the product %s endpoint.", productCode)
	}
	sdkConfig := client.teaSdkConfig
	sdkConfig.SetEndpoint(endpoint)
	conn, err := rpc.NewClient(&sdkConfig)
	if err != nil {
		return nil, fmt.Errorf("unable to initialize the %s client: %#v", productCode, err)
	}
	return conn, nil
}

func (client *AliyunClient) NewServerlessClient() (*roa.Client, error) {
	productCode := "serverless"
	endpoint := ""
	if v, ok := client.config.Endpoints[productCode]; !ok || v.(string) == "" {
		if err := client.loadEndpoint(productCode); err != nil {
			endpoint = fmt.Sprintf("serverless.%s.aliyuncs.com", client.config.RegionId)
			client.config.Endpoints[productCode] = endpoint
			log.Printf("[ERROR] loading %s endpoint got an error: %#v. Using the endpoint %s instead.", productCode, err, endpoint)
		}
	}
	if v, ok := client.config.Endpoints[productCode]; ok && v.(string) != "" {
		endpoint = v.(string)
	}
	if endpoint == "" {
		return nil, fmt.Errorf("[ERROR] missing the product %s endpoint.", productCode)
	}
	sdkConfig := client.teaRoaSdkConfig
	sdkConfig.SetEndpoint(endpoint)
	conn, err := roa.NewClient(&sdkConfig)
	if err != nil {
		return nil, fmt.Errorf("unable to initialize the %s client: %#v", productCode, err)
	}
	return conn, nil
}

func (client *AliyunClient) NewAlbClient() (*rpc.Client, error) {
	productCode := "alb"
	endpoint := ""
	if v, ok := client.config.Endpoints[productCode]; !ok || v.(string) == "" {
		if err := client.loadEndpoint(productCode); err != nil {
			return nil, err
		}
	}
	if v, ok := client.config.Endpoints[productCode]; ok && v.(string) != "" {
		endpoint = v.(string)
	}
	if endpoint == "" {
		return nil, fmt.Errorf("[ERROR] missing the product %s endpoint.", productCode)
	}
	sdkConfig := client.teaSdkConfig
	sdkConfig.SetEndpoint(endpoint)
	conn, err := rpc.NewClient(&sdkConfig)
	if err != nil {
		return nil, fmt.Errorf("unable to initialize the %s client: %#v", productCode, err)
	}
	return conn, nil
}

func (client *AliyunClient) NewRedisaClient() (*rpc.Client, error) {
	productCode := "redisa"
	endpoint := ""
	if v, ok := client.config.Endpoints[productCode]; !ok || v.(string) == "" {
		if err := client.loadEndpoint(productCode); err != nil {
			endpoint = "r-kvstore.aliyuncs.com"
			client.config.Endpoints[productCode] = endpoint
			log.Printf("[ERROR] loading %s endpoint got an error: %#v. Using the endpoint %s instead.", productCode, err, endpoint)
		}
	}
	if v, ok := client.config.Endpoints[productCode]; ok && v.(string) != "" {
		endpoint = v.(string)
	}
	if endpoint == "" {
		return nil, fmt.Errorf("[ERROR] missing the product %s endpoint.", productCode)
	}
	sdkConfig := client.teaSdkConfig
	sdkConfig.SetEndpoint(endpoint)
	conn, err := rpc.NewClient(&sdkConfig)
	if err != nil {
		return nil, fmt.Errorf("unable to initialize the %s client: %#v", productCode, err)
	}
	return conn, nil
}

func (client *AliyunClient) NewGwsecdClient() (*rpc.Client, error) {
	productCode := "gwsecd"
	endpoint := ""
	if v, ok := client.config.Endpoints[productCode]; !ok || v.(string) == "" {
		if err := client.loadEndpoint(productCode); err != nil {
			endpoint = fmt.Sprintf("gwsecd.%s.aliyuncs.com", client.config.RegionId)
			client.config.Endpoints[productCode] = endpoint
			log.Printf("[ERROR] loading %s endpoint got an error: %#v. Using the endpoint %s instead.", productCode, err, endpoint)
		}
	}
	if v, ok := client.config.Endpoints[productCode]; ok && v.(string) != "" {
		endpoint = v.(string)
	}
	if endpoint == "" {
		return nil, fmt.Errorf("[ERROR] missing the product %s endpoint.", productCode)
	}
	sdkConfig := client.teaSdkConfig
	sdkConfig.SetEndpoint(endpoint)
	conn, err := rpc.NewClient(&sdkConfig)
	if err != nil {
		return nil, fmt.Errorf("unable to initialize the %s client: %#v", productCode, err)
	}
	return conn, nil
}

func (client *AliyunClient) NewCloudphoneClient() (*rpc.Client, error) {
	productCode := "cloudphone"
	endpoint := ""
	if v, ok := client.config.Endpoints[productCode]; !ok || v.(string) == "" {
		if err := client.loadEndpoint(productCode); err != nil {
			endpoint = fmt.Sprintf("cloudphone.%s.aliyuncs.com", client.RegionId)
			client.config.Endpoints[productCode] = endpoint
			log.Printf("[ERROR] loading %s endpoint got an error: %#v. Using the endpoint %s instead.", productCode, err, endpoint)
		}
	}
	if v, ok := client.config.Endpoints[productCode]; ok && v.(string) != "" {
		endpoint = v.(string)
	}
	if endpoint == "" {
		return nil, fmt.Errorf("[ERROR] missing the product %s endpoint.", productCode)
	}
	sdkConfig := client.teaSdkConfig
	sdkConfig.SetEndpoint(endpoint)
	conn, err := rpc.NewClient(&sdkConfig)
	if err != nil {
		return nil, fmt.Errorf("unable to initialize the %s client: %#v", productCode, err)
	}
	return conn, nil
}

func (client *AliyunClient) NewScdnClient() (*rpc.Client, error) {
	productCode := "scdn"
	endpoint := ""
	if v, ok := client.config.Endpoints[productCode]; !ok || v.(string) == "" {
		if err := client.loadEndpoint(productCode); err != nil {
			endpoint = "scdn.aliyuncs.com"
			client.config.Endpoints[productCode] = endpoint
			log.Printf("[ERROR] loading %s endpoint got an error: %#v. Using the endpoint %s instead.", productCode, err, endpoint)
		}
	}
	if v, ok := client.config.Endpoints[productCode]; ok && v.(string) != "" {
		endpoint = v.(string)
	}
	if endpoint == "" {
		return nil, fmt.Errorf("[ERROR] missing the product %s endpoint.", productCode)
	}
	sdkConfig := client.teaSdkConfig
	sdkConfig.SetEndpoint(endpoint)
	conn, err := rpc.NewClient(&sdkConfig)
	if err != nil {
		return nil, fmt.Errorf("unable to initialize the %s client: %#v", productCode, err)
	}
	return conn, nil
}

func (client *AliyunClient) NewDataworkspublicClient() (*rpc.Client, error) {
	productCode := "dataworkspublic"
	endpoint := ""
	if v, ok := client.config.Endpoints[productCode]; !ok || v.(string) == "" {
		if err := client.loadEndpoint(productCode); err != nil {
			endpoint = fmt.Sprintf("dataworks.%s.aliyuncs.com", client.config.RegionId)
			client.config.Endpoints[productCode] = endpoint
			log.Printf("[ERROR] loading %s endpoint got an error: %#v. Using the endpoint %s instead.", productCode, err, endpoint)
		}
	}
	if v, ok := client.config.Endpoints[productCode]; ok && v.(string) != "" {
		endpoint = v.(string)
	}
	if endpoint == "" {
		return nil, fmt.Errorf("[ERROR] missing the product %s endpoint.", productCode)
	}
	sdkConfig := client.teaSdkConfig
	sdkConfig.SetEndpoint(endpoint)
	conn, err := rpc.NewClient(&sdkConfig)
	if err != nil {
		return nil, fmt.Errorf("unable to initialize the %s client: %#v", productCode, err)
	}
	return conn, nil
}

func (client *AliyunClient) NewCdnClient() (*rpc.Client, error) {
	productCode := "cdn"
	endpoint := ""
	if v, ok := client.config.Endpoints[productCode]; !ok || v.(string) == "" {
		if err := client.loadEndpoint(productCode); err != nil {
			endpoint = "cdn.aliyuncs.com"
			client.config.Endpoints[productCode] = endpoint
			log.Printf("[ERROR] loading %s endpoint got an error: %#v. Using the central endpoint %s instead.", productCode, err, endpoint)
		}
	}
	if v, ok := client.config.Endpoints[productCode]; ok && v.(string) != "" {
		endpoint = v.(string)
	}
	if endpoint == "" {
		return nil, fmt.Errorf("[ERROR] missing the product %s endpoint.", productCode)
	}

	sdkConfig := client.teaSdkConfig
	sdkConfig.SetEndpoint(endpoint)
	conn, err := rpc.NewClient(&sdkConfig)
	if err != nil {
		return nil, fmt.Errorf("unable to initialize the %s client: %#v", productCode, err)
	}
	return conn, nil
}

func (client *AliyunClient) NewCddcClient() (*rpc.Client, error) {
	productCode := "cddc"
	endpoint := ""
	if v, ok := client.config.Endpoints[productCode]; !ok || v.(string) == "" {
		if err := client.loadEndpoint(productCode); err != nil {
			endpoint = fmt.Sprintf("cddc.%s.aliyuncs.com", client.config.RegionId)
			client.config.Endpoints[productCode] = endpoint
			log.Printf("[ERROR] loading %s endpoint got an error: %#v. Using the endpoint %s instead.", productCode, err, endpoint)
		}
	}
	if v, ok := client.config.Endpoints[productCode]; ok && v.(string) != "" {
		endpoint = v.(string)
	}
	if endpoint == "" {
		return nil, fmt.Errorf("[ERROR] missing the product %s endpoint.", productCode)
	}
	sdkConfig := client.teaSdkConfig
	sdkConfig.SetEndpoint(endpoint)
	conn, err := rpc.NewClient(&sdkConfig)
	if err != nil {
		return nil, fmt.Errorf("unable to initialize the %s client: %#v", productCode, err)
	}
	return conn, nil
}

func (client *AliyunClient) NewMscopensubscriptionClient() (*rpc.Client, error) {
	productCode := "mscsub"
	endpoint := ""
	if v, ok := client.config.Endpoints[productCode]; !ok || v.(string) == "" {
		if err := client.loadEndpoint(productCode); err != nil {
			endpoint = "mscopensubscription.aliyuncs.com"
			client.config.Endpoints[productCode] = endpoint
			log.Printf("[ERROR] loading %s endpoint got an error: %#v. Using the endpoint %s instead.", productCode, err, endpoint)
		}
	}
	if v, ok := client.config.Endpoints[productCode]; ok && v.(string) != "" {
		endpoint = v.(string)
	}
	if endpoint == "" {
		return nil, fmt.Errorf("[ERROR] missing the product %s endpoint.", productCode)
	}
	sdkConfig := client.teaSdkConfig
	sdkConfig.SetEndpoint(endpoint)
	conn, err := rpc.NewClient(&sdkConfig)
	if err != nil {
		return nil, fmt.Errorf("unable to initialize the %s client: %#v", productCode, err)
	}
	return conn, nil
}

func (client *AliyunClient) NewSddpClient() (*rpc.Client, error) {
	productCode := "sddp"
	endpoint := ""
	if v, ok := client.config.Endpoints[productCode]; !ok || v.(string) == "" {
		if err := client.loadEndpoint(productCode); err != nil {
			endpoint = fmt.Sprintf("sddp.%s.aliyuncs.com", client.config.RegionId)
			client.config.Endpoints[productCode] = endpoint
			log.Printf("[ERROR] loading %s endpoint got an error: %#v. Using the endpoint %s instead.", productCode, err, endpoint)
		}
	}
	if v, ok := client.config.Endpoints[productCode]; ok && v.(string) != "" {
		endpoint = v.(string)
	}
	sdkConfig := client.teaSdkConfig
	sdkConfig.SetEndpoint(endpoint)
	conn, err := rpc.NewClient(&sdkConfig)
	if err != nil {
		return nil, fmt.Errorf("unable to initialize the %s client: %#v", productCode, err)
	}
	return conn, nil
}

func (client *AliyunClient) NewBastionhostClient() (*rpc.Client, error) {
	productCode := "bastionhost"
	endpoint := ""
	if v, ok := client.config.Endpoints[productCode]; !ok || v.(string) == "" {
		if err := client.loadEndpoint(productCode); err != nil {
			endpoint = fmt.Sprintf("bastionhost.%s.aliyuncs.com", client.config.RegionId)
			client.config.Endpoints[productCode] = endpoint
			log.Printf("[ERROR] loading %s endpoint got an error: %#v. Using the endpoint %s instead.", productCode, err, endpoint)
		}
	}
	if v, ok := client.config.Endpoints[productCode]; ok && v.(string) != "" {
		endpoint = v.(string)
	}
	if endpoint == "" {
		return nil, fmt.Errorf("[ERROR] missing the product %s endpoint.", productCode)
	}
	sdkConfig := client.teaSdkConfig
	sdkConfig.SetEndpoint(endpoint)
	conn, err := rpc.NewClient(&sdkConfig)
	if err != nil {
		return nil, fmt.Errorf("unable to initialize the %s client: %#v", productCode, err)
	}
	return conn, nil
}
