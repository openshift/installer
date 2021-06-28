// Copyright IBM Corp. 2017, 2021 All Rights Reserved.
// Licensed under the Mozilla Public License v2.0

package ibm

import (
	"crypto/tls"
	"errors"
	"fmt"
	"log"
	"net"
	gohttp "net/http"
	"os"
	"strings"
	"time"

	// Added code for the Power Colo Offering

	apigateway "github.com/IBM/apigateway-go-sdk"
	"github.com/IBM/appconfiguration-go-admin-sdk/appconfigurationv1"
	"github.com/IBM/container-registry-go-sdk/containerregistryv1"
	"github.com/IBM/go-sdk-core/v4/core"
	cosconfig "github.com/IBM/ibm-cos-sdk-go-config/resourceconfigurationv1"
	kp "github.com/IBM/keyprotect-go-client"
	ciscachev1 "github.com/IBM/networking-go-sdk/cachingapiv1"
	cisipv1 "github.com/IBM/networking-go-sdk/cisipapiv1"
	ciscustompagev1 "github.com/IBM/networking-go-sdk/custompagesv1"
	dlProviderV2 "github.com/IBM/networking-go-sdk/directlinkproviderv2"
	dl "github.com/IBM/networking-go-sdk/directlinkv1"
	cisdnsbulkv1 "github.com/IBM/networking-go-sdk/dnsrecordbulkv1"
	cisdnsrecordsv1 "github.com/IBM/networking-go-sdk/dnsrecordsv1"
	dns "github.com/IBM/networking-go-sdk/dnssvcsv1"
	cisedgefunctionv1 "github.com/IBM/networking-go-sdk/edgefunctionsapiv1"
	cisglbhealthcheckv1 "github.com/IBM/networking-go-sdk/globalloadbalancermonitorv1"
	cisglbpoolv0 "github.com/IBM/networking-go-sdk/globalloadbalancerpoolsv0"
	cisglbv1 "github.com/IBM/networking-go-sdk/globalloadbalancerv1"
	cispagerulev1 "github.com/IBM/networking-go-sdk/pageruleapiv1"
	cisrangeappv1 "github.com/IBM/networking-go-sdk/rangeapplicationsv1"
	cisroutingv1 "github.com/IBM/networking-go-sdk/routingv1"
	cissslv1 "github.com/IBM/networking-go-sdk/sslcertificateapiv1"
	tg "github.com/IBM/networking-go-sdk/transitgatewayapisv1"
	cisuarulev1 "github.com/IBM/networking-go-sdk/useragentblockingrulesv1"
	ciswafgroupv1 "github.com/IBM/networking-go-sdk/wafrulegroupsapiv1"
	ciswafpackagev1 "github.com/IBM/networking-go-sdk/wafrulepackagesapiv1"
	ciswafrulev1 "github.com/IBM/networking-go-sdk/wafrulesapiv1"
	cisaccessrulev1 "github.com/IBM/networking-go-sdk/zonefirewallaccessrulesv1"
	cislockdownv1 "github.com/IBM/networking-go-sdk/zonelockdownv1"
	cisratelimitv1 "github.com/IBM/networking-go-sdk/zoneratelimitsv1"
	cisdomainsettingsv1 "github.com/IBM/networking-go-sdk/zonessettingsv1"
	ciszonesv1 "github.com/IBM/networking-go-sdk/zonesv1"
	"github.com/IBM/platform-services-go-sdk/catalogmanagementv1"
	"github.com/IBM/platform-services-go-sdk/enterprisemanagementv1"
	"github.com/IBM/platform-services-go-sdk/globaltaggingv1"
	iamidentity "github.com/IBM/platform-services-go-sdk/iamidentityv1"
	iampolicymanagement "github.com/IBM/platform-services-go-sdk/iampolicymanagementv1"
	resourcecontroller "github.com/IBM/platform-services-go-sdk/resourcecontrollerv2"
	resourcemanager "github.com/IBM/platform-services-go-sdk/resourcemanagerv2"
	"github.com/IBM/push-notifications-go-sdk/pushservicev1"
	schematicsv1 "github.com/IBM/schematics-go-sdk/schematicsv1"
	"github.com/IBM/secrets-manager-go-sdk/secretsmanagerv1"
	vpcclassic "github.com/IBM/vpc-go-sdk/vpcclassicv1"
	vpc "github.com/IBM/vpc-go-sdk/vpcv1"
	"github.com/apache/openwhisk-client-go/whisk"
	jwt "github.com/dgrijalva/jwt-go"
	slsession "github.com/softlayer/softlayer-go/session"

	bluemix "github.com/IBM-Cloud/bluemix-go"
	"github.com/IBM-Cloud/bluemix-go/api/account/accountv1"
	"github.com/IBM-Cloud/bluemix-go/api/account/accountv2"
	"github.com/IBM-Cloud/bluemix-go/api/certificatemanager"
	"github.com/IBM-Cloud/bluemix-go/api/cis/cisv1"
	"github.com/IBM-Cloud/bluemix-go/api/container/containerv1"
	"github.com/IBM-Cloud/bluemix-go/api/container/containerv2"
	"github.com/IBM-Cloud/bluemix-go/api/functions"
	"github.com/IBM-Cloud/bluemix-go/api/globalsearch/globalsearchv2"
	"github.com/IBM-Cloud/bluemix-go/api/globaltagging/globaltaggingv3"
	"github.com/IBM-Cloud/bluemix-go/api/hpcs"
	"github.com/IBM-Cloud/bluemix-go/api/iam/iamv1"
	"github.com/IBM-Cloud/bluemix-go/api/iamuum/iamuumv2"
	"github.com/IBM-Cloud/bluemix-go/api/icd/icdv4"
	"github.com/IBM-Cloud/bluemix-go/api/mccp/mccpv2"
	"github.com/IBM-Cloud/bluemix-go/api/resource/resourcev1/catalog"
	"github.com/IBM-Cloud/bluemix-go/api/resource/resourcev1/controller"
	"github.com/IBM-Cloud/bluemix-go/api/resource/resourcev2/controllerv2"
	"github.com/IBM-Cloud/bluemix-go/api/resource/resourcev2/managementv2"
	"github.com/IBM-Cloud/bluemix-go/api/usermanagement/usermanagementv2"
	"github.com/IBM-Cloud/bluemix-go/authentication"
	"github.com/IBM-Cloud/bluemix-go/bmxerror"
	"github.com/IBM-Cloud/bluemix-go/http"
	"github.com/IBM-Cloud/bluemix-go/rest"
	bxsession "github.com/IBM-Cloud/bluemix-go/session"
	ibmpisession "github.com/IBM-Cloud/power-go-client/ibmpisession"
	"github.com/IBM-Cloud/terraform-provider-ibm/version"
)

// RetryAPIDelay - retry api delay
const RetryAPIDelay = 5 * time.Second

//BluemixRegion ...
var BluemixRegion string

var (
	errEmptySoftLayerCredentials = errors.New("iaas_classic_username and iaas_classic_api_key must be provided. Please see the documentation on how to configure them")
	errEmptyBluemixCredentials   = errors.New("ibmcloud_api_key or bluemix_api_key or iam_token and iam_refresh_token must be provided. Please see the documentation on how to configure it")
)

//UserConfig ...
type UserConfig struct {
	userID      string
	userEmail   string
	userAccount string
	cloudName   string `default:"bluemix"`
	cloudType   string `default:"public"`
	generation  int    `default:"2"`
}

//Config stores user provider input
type Config struct {
	//BluemixAPIKey is the Bluemix api key
	BluemixAPIKey string
	//Bluemix region
	Region string
	//Resource group id
	ResourceGroup string
	//Bluemix API timeout
	BluemixTimeout time.Duration

	//Softlayer end point url
	SoftLayerEndpointURL string

	//Softlayer API timeout
	SoftLayerTimeout time.Duration

	// Softlayer User Name
	SoftLayerUserName string

	// Softlayer API Key
	SoftLayerAPIKey string

	//Retry Count for API calls
	//Unexposed in the schema at this point as they are used only during session creation for a few calls
	//When sdk implements it we an expose them for expected behaviour
	//https://github.com/softlayer/softlayer-go/issues/41
	RetryCount int
	//Constant Retry Delay for API calls
	RetryDelay time.Duration

	// FunctionNameSpace ...
	FunctionNameSpace string

	//Riaas End point
	RiaasEndPoint string

	//Generation
	Generation int

	//IAM Token
	IAMToken string

	//IAM Refresh Token
	IAMRefreshToken string

	// PowerService Instance
	PowerServiceInstance string

	// Zone
	Zone       string
	Visibility string
}

//Session stores the information required for communication with the SoftLayer and Bluemix API
type Session struct {
	// SoftLayerSesssion is the the SoftLayer session used to connect to the SoftLayer API
	SoftLayerSession *slsession.Session

	// BluemixSession is the the Bluemix session used to connect to the Bluemix API
	BluemixSession *bxsession.Session
}

// ClientSession ...
type ClientSession interface {
	BluemixSession() (*bxsession.Session, error)
	BluemixAcccountAPI() (accountv2.AccountServiceAPI, error)
	BluemixAcccountv1API() (accountv1.AccountServiceAPI, error)
	BluemixUserDetails() (*UserConfig, error)
	ContainerAPI() (containerv1.ContainerServiceAPI, error)
	VpcContainerAPI() (containerv2.ContainerServiceAPI, error)
	ContainerRegistryV1() (*containerregistryv1.ContainerRegistryV1, error)
	CisAPI() (cisv1.CisServiceAPI, error)
	FunctionClient() (*whisk.Client, error)
	GlobalSearchAPI() (globalsearchv2.GlobalSearchServiceAPI, error)
	GlobalTaggingAPI() (globaltaggingv3.GlobalTaggingServiceAPI, error)
	GlobalTaggingAPIv1() (globaltaggingv1.GlobalTaggingV1, error)
	ICDAPI() (icdv4.ICDServiceAPI, error)
	IAMAPI() (iamv1.IAMServiceAPI, error)
	IAMPolicyManagementV1API() (*iampolicymanagement.IamPolicyManagementV1, error)
	IAMUUMAPIV2() (iamuumv2.IAMUUMServiceAPIv2, error)
	MccpAPI() (mccpv2.MccpServiceAPI, error)
	ResourceCatalogAPI() (catalog.ResourceCatalogAPI, error)
	ResourceManagementAPIv2() (managementv2.ResourceManagementAPIv2, error)
	ResourceControllerAPI() (controller.ResourceControllerAPI, error)
	ResourceControllerAPIV2() (controllerv2.ResourceControllerAPIV2, error)
	SoftLayerSession() *slsession.Session
	IBMPISession() (*ibmpisession.IBMPISession, error)
	UserManagementAPI() (usermanagementv2.UserManagementAPI, error)
	PushServiceV1() (*pushservicev1.PushServiceV1, error)
	AppConfigurationV1() (*appconfigurationv1.AppConfigurationV1, error)
	CertificateManagerAPI() (certificatemanager.CertificateManagerServiceAPI, error)
	keyProtectAPI() (*kp.Client, error)
	keyManagementAPI() (*kp.Client, error)
	VpcClassicV1API() (*vpcclassic.VpcClassicV1, error)
	VpcV1API() (*vpc.VpcV1, error)
	APIGateway() (*apigateway.ApiGatewayControllerApiV1, error)
	PrivateDNSClientSession() (*dns.DnsSvcsV1, error)
	CosConfigV1API() (*cosconfig.ResourceConfigurationV1, error)
	DirectlinkV1API() (*dl.DirectLinkV1, error)
	DirectlinkProviderV2API() (*dlProviderV2.DirectLinkProviderV2, error)
	TransitGatewayV1API() (*tg.TransitGatewayApisV1, error)
	HpcsEndpointAPI() (hpcs.HPCSV2, error)
	FunctionIAMNamespaceAPI() (functions.FunctionServiceAPI, error)
	CisZonesV1ClientSession() (*ciszonesv1.ZonesV1, error)
	CisDNSRecordClientSession() (*cisdnsrecordsv1.DnsRecordsV1, error)
	CisDNSRecordBulkClientSession() (*cisdnsbulkv1.DnsRecordBulkV1, error)
	CisGLBClientSession() (*cisglbv1.GlobalLoadBalancerV1, error)
	CisGLBPoolClientSession() (*cisglbpoolv0.GlobalLoadBalancerPoolsV0, error)
	CisGLBHealthCheckClientSession() (*cisglbhealthcheckv1.GlobalLoadBalancerMonitorV1, error)
	CisIPClientSession() (*cisipv1.CisIpApiV1, error)
	CisPageRuleClientSession() (*cispagerulev1.PageRuleApiV1, error)
	CisRLClientSession() (*cisratelimitv1.ZoneRateLimitsV1, error)
	CisEdgeFunctionClientSession() (*cisedgefunctionv1.EdgeFunctionsApiV1, error)
	CisSSLClientSession() (*cissslv1.SslCertificateApiV1, error)
	CisWAFPackageClientSession() (*ciswafpackagev1.WafRulePackagesApiV1, error)
	CisDomainSettingsClientSession() (*cisdomainsettingsv1.ZonesSettingsV1, error)
	CisRoutingClientSession() (*cisroutingv1.RoutingV1, error)
	CisWAFGroupClientSession() (*ciswafgroupv1.WafRuleGroupsApiV1, error)
	CisCacheClientSession() (*ciscachev1.CachingApiV1, error)
	CisCustomPageClientSession() (*ciscustompagev1.CustomPagesV1, error)
	CisAccessRuleClientSession() (*cisaccessrulev1.ZoneFirewallAccessRulesV1, error)
	CisUARuleClientSession() (*cisuarulev1.UserAgentBlockingRulesV1, error)
	CisLockdownClientSession() (*cislockdownv1.ZoneLockdownV1, error)
	CisRangeAppClientSession() (*cisrangeappv1.RangeApplicationsV1, error)
	CisWAFRuleClientSession() (*ciswafrulev1.WafRulesApiV1, error)
	IAMIdentityV1API() (*iamidentity.IamIdentityV1, error)
	ResourceManagerV2API() (*resourcemanager.ResourceManagerV2, error)
	CatalogManagementV1() (*catalogmanagementv1.CatalogManagementV1, error)
	EnterpriseManagementV1() (*enterprisemanagementv1.EnterpriseManagementV1, error)
	ResourceControllerV2API() (*resourcecontroller.ResourceControllerV2, error)
	SecretsManagerV1() (*secretsmanagerv1.SecretsManagerV1, error)
	SchematicsV1() (*schematicsv1.SchematicsV1, error)
}

type clientSession struct {
	session *Session

	apigatewayErr error
	apigatewayAPI *apigateway.ApiGatewayControllerApiV1

	accountConfigErr     error
	bmxAccountServiceAPI accountv2.AccountServiceAPI

	accountV1ConfigErr     error
	bmxAccountv1ServiceAPI accountv1.AccountServiceAPI

	bmxUserDetails  *UserConfig
	bmxUserFetchErr error

	csConfigErr  error
	csServiceAPI containerv1.ContainerServiceAPI

	csv2ConfigErr  error
	csv2ServiceAPI containerv2.ContainerServiceAPI

	containerRegistryClientErr error
	containerRegistryClient    *containerregistryv1.ContainerRegistryV1

	certManagementErr error
	certManagementAPI certificatemanager.CertificateManagerServiceAPI

	cfConfigErr  error
	cfServiceAPI mccpv2.MccpServiceAPI

	cisConfigErr  error
	cisServiceAPI cisv1.CisServiceAPI

	functionConfigErr error
	functionClient    *whisk.Client

	globalSearchConfigErr  error
	globalSearchServiceAPI globalsearchv2.GlobalSearchServiceAPI

	globalTaggingConfigErr  error
	globalTaggingServiceAPI globaltaggingv3.GlobalTaggingServiceAPI

	globalTaggingConfigErrV1  error
	globalTaggingServiceAPIV1 globaltaggingv1.GlobalTaggingV1

	iamUUMConfigErrV2  error
	iamUUMServiceAPIV2 iamuumv2.IAMUUMServiceAPIv2

	iamConfigErr  error
	iamServiceAPI iamv1.IAMServiceAPI

	userManagementErr error
	userManagementAPI usermanagementv2.UserManagementAPI

	icdConfigErr  error
	icdServiceAPI icdv4.ICDServiceAPI

	resourceControllerConfigErr  error
	resourceControllerServiceAPI controller.ResourceControllerAPI

	resourceControllerConfigErrv2  error
	resourceControllerServiceAPIv2 controllerv2.ResourceControllerAPIV2

	resourceManagementConfigErrv2  error
	resourceManagementServiceAPIv2 managementv2.ResourceManagementAPIv2

	resourceCatalogConfigErr  error
	resourceCatalogServiceAPI catalog.ResourceCatalogAPI

	powerConfigErr error
	ibmpiConfigErr error
	ibmpiSession   *ibmpisession.IBMPISession

	kpErr error
	kpAPI *kp.API

	kmsErr error
	kmsAPI *kp.API

	hpcsEndpointErr error
	hpcsEndpointAPI hpcs.HPCSV2

	pDNSClient *dns.DnsSvcsV1
	pDNSErr    error

	bluemixSessionErr error

	pushServiceClient    *pushservicev1.PushServiceV1
	pushServiceClientErr error

	appConfigurationClient    *appconfigurationv1.AppConfigurationV1
	appConfigurationClientErr error

	vpcClassicErr error
	vpcClassicAPI *vpcclassic.VpcClassicV1

	vpcErr error
	vpcAPI *vpc.VpcV1

	directlinkAPI *dl.DirectLinkV1
	directlinkErr error
	dlProviderAPI *dlProviderV2.DirectLinkProviderV2
	dlProviderErr error

	cosConfigErr error
	cosConfigAPI *cosconfig.ResourceConfigurationV1

	transitgatewayAPI *tg.TransitGatewayApisV1
	transitgatewayErr error

	functionIAMNamespaceAPI functions.FunctionServiceAPI
	functionIAMNamespaceErr error

	// CIS Zones
	cisZonesErr      error
	cisZonesV1Client *ciszonesv1.ZonesV1

	// CIS dns service options
	cisDNSErr           error
	cisDNSRecordsClient *cisdnsrecordsv1.DnsRecordsV1

	// CIS dns bulk service options
	cisDNSBulkErr          error
	cisDNSRecordBulkClient *cisdnsbulkv1.DnsRecordBulkV1

	// CIS Global Load Balancer Pool service options
	cisGLBPoolErr    error
	cisGLBPoolClient *cisglbpoolv0.GlobalLoadBalancerPoolsV0

	// CIS GLB service options
	cisGLBErr    error
	cisGLBClient *cisglbv1.GlobalLoadBalancerV1

	// CIS GLB health check service options
	cisGLBHealthCheckErr    error
	cisGLBHealthCheckClient *cisglbhealthcheckv1.GlobalLoadBalancerMonitorV1

	// CIS IP service options
	cisIPErr    error
	cisIPClient *cisipv1.CisIpApiV1

	// CIS Zone Rate Limits service options
	cisRLErr    error
	cisRLClient *cisratelimitv1.ZoneRateLimitsV1

	// CIS Page Rules service options
	cisPageRuleErr    error
	cisPageRuleClient *cispagerulev1.PageRuleApiV1

	// CIS Edge Functions service options
	cisEdgeFunctionErr    error
	cisEdgeFunctionClient *cisedgefunctionv1.EdgeFunctionsApiV1

	// CIS SSL certificate service options
	cisSSLErr    error
	cisSSLClient *cissslv1.SslCertificateApiV1

	// CIS WAF Package service options
	cisWAFPackageErr    error
	cisWAFPackageClient *ciswafpackagev1.WafRulePackagesApiV1

	// CIS Zone Setting service options
	cisDomainSettingsErr    error
	cisDomainSettingsClient *cisdomainsettingsv1.ZonesSettingsV1

	// CIS Routing service options
	cisRoutingErr    error
	cisRoutingClient *cisroutingv1.RoutingV1

	// CIS WAF Group service options
	cisWAFGroupErr    error
	cisWAFGroupClient *ciswafgroupv1.WafRuleGroupsApiV1

	// CIS Caching service options
	cisCacheErr    error
	cisCacheClient *ciscachev1.CachingApiV1

	// CIS Custom Pages service options
	cisCustomPageErr    error
	cisCustomPageClient *ciscustompagev1.CustomPagesV1

	// CIS Firewall Access rule service option
	cisAccessRuleErr    error
	cisAccessRuleClient *cisaccessrulev1.ZoneFirewallAccessRulesV1

	// CIS User Agent Blocking Rule service option
	cisUARuleErr    error
	cisUARuleClient *cisuarulev1.UserAgentBlockingRulesV1

	// CIS Firewall Lockdwon Rule service option
	cisLockdownErr    error
	cisLockdownClient *cislockdownv1.ZoneLockdownV1

	// CIS Range app service option
	cisRangeAppErr    error
	cisRangeAppClient *cisrangeappv1.RangeApplicationsV1

	// CIS WAF rule service options
	cisWAFRuleErr    error
	cisWAFRuleClient *ciswafrulev1.WafRulesApiV1
	//IAM Identity Option
	iamIdentityErr error
	iamIdentityAPI *iamidentity.IamIdentityV1

	//Resource Manager Option
	resourceManagerErr error
	resourceManagerAPI *resourcemanager.ResourceManagerV2

	//Catalog Management Option
	catalogManagementClient    *catalogmanagementv1.CatalogManagementV1
	catalogManagementClientErr error

	enterpriseManagementClient    *enterprisemanagementv1.EnterpriseManagementV1
	enterpriseManagementClientErr error

	//Resource Controller Option
	resourceControllerErr   error
	resourceControllerAPI   *resourcecontroller.ResourceControllerV2
	secretsManagerClient    *secretsmanagerv1.SecretsManagerV1
	secretsManagerClientErr error

	// Schematics service options
	schematicsClient    *schematicsv1.SchematicsV1
	schematicsClientErr error

	//IAM Policy Management
	iamPolicyManagementErr error
	iamPolicyManagementAPI *iampolicymanagement.IamPolicyManagementV1
}

func (session clientSession) CatalogManagementV1() (*catalogmanagementv1.CatalogManagementV1, error) {
	return session.catalogManagementClient, session.catalogManagementClientErr
}

// BluemixAcccountAPI ...
func (sess clientSession) BluemixAcccountAPI() (accountv2.AccountServiceAPI, error) {
	return sess.bmxAccountServiceAPI, sess.accountConfigErr
}

// BluemixAcccountAPI ...
func (sess clientSession) BluemixAcccountv1API() (accountv1.AccountServiceAPI, error) {
	return sess.bmxAccountv1ServiceAPI, sess.accountV1ConfigErr
}

// BluemixSession to provide the Bluemix Session
func (sess clientSession) BluemixSession() (*bxsession.Session, error) {
	return sess.session.BluemixSession, sess.bluemixSessionErr
}

// BluemixUserDetails ...
func (sess clientSession) BluemixUserDetails() (*UserConfig, error) {
	return sess.bmxUserDetails, sess.bmxUserFetchErr
}

// ContainerAPI provides Container Service APIs ...
func (sess clientSession) ContainerAPI() (containerv1.ContainerServiceAPI, error) {
	return sess.csServiceAPI, sess.csConfigErr
}

// VpcContainerAPI provides v2Container Service APIs ...
func (sess clientSession) VpcContainerAPI() (containerv2.ContainerServiceAPI, error) {
	return sess.csv2ServiceAPI, sess.csv2ConfigErr
}

// ContainerRegistryV1 provides Container Registry Service APIs ...
func (session clientSession) ContainerRegistryV1() (*containerregistryv1.ContainerRegistryV1, error) {
	return session.containerRegistryClient, session.containerRegistryClientErr
}

// SchematicsAPI provides schematics Service APIs ...
func (sess clientSession) SchematicsV1() (*schematicsv1.SchematicsV1, error) {
	return sess.schematicsClient, sess.schematicsClientErr
}

// CisAPI provides Cloud Internet Services APIs ...
func (sess clientSession) CisAPI() (cisv1.CisServiceAPI, error) {
	return sess.cisServiceAPI, sess.cisConfigErr
}

// FunctionClient ...
func (sess clientSession) FunctionClient() (*whisk.Client, error) {
	return sess.functionClient, sess.functionConfigErr
}

// GlobalSearchAPI provides Global Search  APIs ...
func (sess clientSession) GlobalSearchAPI() (globalsearchv2.GlobalSearchServiceAPI, error) {
	return sess.globalSearchServiceAPI, sess.globalSearchConfigErr
}

// GlobalTaggingAPI provides Global Search  APIs ...
func (sess clientSession) GlobalTaggingAPI() (globaltaggingv3.GlobalTaggingServiceAPI, error) {
	return sess.globalTaggingServiceAPI, sess.globalTaggingConfigErr
}

// GlobalTaggingAPIV1 provides Platform-go Global Tagging  APIs ...
func (sess clientSession) GlobalTaggingAPIv1() (globaltaggingv1.GlobalTaggingV1, error) {
	return sess.globalTaggingServiceAPIV1, sess.globalTaggingConfigErrV1
}

// HpcsEndpointAPI provides Hpcs Endpoint generator APIs ...
func (sess clientSession) HpcsEndpointAPI() (hpcs.HPCSV2, error) {
	return sess.hpcsEndpointAPI, sess.hpcsEndpointErr
}

// IAMAPI provides IAM PAP APIs ...
func (sess clientSession) IAMAPI() (iamv1.IAMServiceAPI, error) {
	return sess.iamServiceAPI, sess.iamConfigErr
}

// UserManagementAPI provides User management APIs ...
func (sess clientSession) UserManagementAPI() (usermanagementv2.UserManagementAPI, error) {
	return sess.userManagementAPI, sess.userManagementErr
}

// IAM Policy Management
func (sess clientSession) IAMPolicyManagementV1API() (*iampolicymanagement.IamPolicyManagementV1, error) {
	return sess.iamPolicyManagementAPI, sess.iamPolicyManagementErr
}

// IAMUUMAPIV2 provides IAM UUM APIs ...
func (sess clientSession) IAMUUMAPIV2() (iamuumv2.IAMUUMServiceAPIv2, error) {
	return sess.iamUUMServiceAPIV2, sess.iamUUMConfigErrV2
}

// IcdAPI provides IBM Cloud Databases APIs ...
func (sess clientSession) ICDAPI() (icdv4.ICDServiceAPI, error) {
	return sess.icdServiceAPI, sess.icdConfigErr
}

// MccpAPI provides Multi Cloud Controller Proxy APIs ...
func (sess clientSession) MccpAPI() (mccpv2.MccpServiceAPI, error) {
	return sess.cfServiceAPI, sess.cfConfigErr
}

// ResourceCatalogAPI ...
func (sess clientSession) ResourceCatalogAPI() (catalog.ResourceCatalogAPI, error) {
	return sess.resourceCatalogServiceAPI, sess.resourceCatalogConfigErr
}

// ResourceManagementAPIv2 ...
func (sess clientSession) ResourceManagementAPIv2() (managementv2.ResourceManagementAPIv2, error) {
	return sess.resourceManagementServiceAPIv2, sess.resourceManagementConfigErrv2
}

// ResourceControllerAPI ...
func (sess clientSession) ResourceControllerAPI() (controller.ResourceControllerAPI, error) {
	return sess.resourceControllerServiceAPI, sess.resourceControllerConfigErr
}

// ResourceControllerAPIv2 ...
func (sess clientSession) ResourceControllerAPIV2() (controllerv2.ResourceControllerAPIV2, error) {
	return sess.resourceControllerServiceAPIv2, sess.resourceControllerConfigErrv2
}

// SoftLayerSession providers SoftLayer Session
func (sess clientSession) SoftLayerSession() *slsession.Session {
	return sess.session.SoftLayerSession
}

// CertManagementAPI provides Certificate  management APIs ...
func (sess clientSession) CertificateManagerAPI() (certificatemanager.CertificateManagerServiceAPI, error) {
	return sess.certManagementAPI, sess.certManagementErr
}

//apigatewayAPI provides API Gateway APIs
func (sess clientSession) APIGateway() (*apigateway.ApiGatewayControllerApiV1, error) {
	return sess.apigatewayAPI, sess.apigatewayErr
}

func (session clientSession) PushServiceV1() (*pushservicev1.PushServiceV1, error) {
	return session.pushServiceClient, session.pushServiceClientErr
}

func (session clientSession) AppConfigurationV1() (*appconfigurationv1.AppConfigurationV1, error) {
	return session.appConfigurationClient, session.appConfigurationClientErr
}

func (sess clientSession) keyProtectAPI() (*kp.Client, error) {
	return sess.kpAPI, sess.kpErr
}

func (sess clientSession) keyManagementAPI() (*kp.Client, error) {
	return sess.kmsAPI, sess.kmsErr
}

func (sess clientSession) VpcClassicV1API() (*vpcclassic.VpcClassicV1, error) {
	return sess.vpcClassicAPI, sess.vpcClassicErr
}

func (sess clientSession) VpcV1API() (*vpc.VpcV1, error) {
	return sess.vpcAPI, sess.vpcErr
}

func (sess clientSession) DirectlinkV1API() (*dl.DirectLinkV1, error) {
	return sess.directlinkAPI, sess.directlinkErr
}
func (sess clientSession) DirectlinkProviderV2API() (*dlProviderV2.DirectLinkProviderV2, error) {
	return sess.dlProviderAPI, sess.dlProviderErr
}
func (sess clientSession) CosConfigV1API() (*cosconfig.ResourceConfigurationV1, error) {
	return sess.cosConfigAPI, sess.cosConfigErr
}

func (sess clientSession) TransitGatewayV1API() (*tg.TransitGatewayApisV1, error) {
	return sess.transitgatewayAPI, sess.transitgatewayErr
}

// Session to the Power Colo Service

func (sess clientSession) IBMPISession() (*ibmpisession.IBMPISession, error) {
	return sess.ibmpiSession, sess.powerConfigErr
}

// Private DNS Service

func (sess clientSession) PrivateDNSClientSession() (*dns.DnsSvcsV1, error) {
	return sess.pDNSClient, sess.pDNSErr
}

// Session to the Namespace cloud function

func (sess clientSession) FunctionIAMNamespaceAPI() (functions.FunctionServiceAPI, error) {
	return sess.functionIAMNamespaceAPI, sess.functionIAMNamespaceErr
}

// CIS Zones Service
func (sess clientSession) CisZonesV1ClientSession() (*ciszonesv1.ZonesV1, error) {
	if sess.cisZonesErr != nil {
		return sess.cisZonesV1Client, sess.cisZonesErr
	}
	return sess.cisZonesV1Client.Clone(), nil
}

// CIS DNS Service
func (sess clientSession) CisDNSRecordClientSession() (*cisdnsrecordsv1.DnsRecordsV1, error) {
	if sess.cisDNSErr != nil {
		return sess.cisDNSRecordsClient, sess.cisDNSErr
	}
	return sess.cisDNSRecordsClient.Clone(), nil
}

// CIS DNS Bulk Service
func (sess clientSession) CisDNSRecordBulkClientSession() (*cisdnsbulkv1.DnsRecordBulkV1, error) {
	if sess.cisDNSBulkErr != nil {
		return sess.cisDNSRecordBulkClient, sess.cisDNSBulkErr
	}
	return sess.cisDNSRecordBulkClient.Clone(), nil
}

// CIS GLB Pool
func (sess clientSession) CisGLBPoolClientSession() (*cisglbpoolv0.GlobalLoadBalancerPoolsV0, error) {
	if sess.cisGLBPoolErr != nil {
		return sess.cisGLBPoolClient, sess.cisGLBPoolErr
	}
	return sess.cisGLBPoolClient.Clone(), nil
}

// CIS GLB
func (sess clientSession) CisGLBClientSession() (*cisglbv1.GlobalLoadBalancerV1, error) {
	if sess.cisGLBErr != nil {
		return sess.cisGLBClient, sess.cisGLBErr
	}
	return sess.cisGLBClient.Clone(), nil
}

// CIS GLB Health Check/Monitor
func (sess clientSession) CisGLBHealthCheckClientSession() (*cisglbhealthcheckv1.GlobalLoadBalancerMonitorV1, error) {
	if sess.cisGLBHealthCheckErr != nil {
		return sess.cisGLBHealthCheckClient, sess.cisGLBHealthCheckErr
	}
	return sess.cisGLBHealthCheckClient.Clone(), nil
}

// CIS Zone Rate Limits
func (sess clientSession) CisRLClientSession() (*cisratelimitv1.ZoneRateLimitsV1, error) {
	if sess.cisRLErr != nil {
		return sess.cisRLClient, sess.cisRLErr
	}
	return sess.cisRLClient.Clone(), nil
}

// CIS IP
func (sess clientSession) CisIPClientSession() (*cisipv1.CisIpApiV1, error) {
	if sess.cisIPErr != nil {
		return sess.cisIPClient, sess.cisIPErr
	}
	return sess.cisIPClient.Clone(), nil
}

// CIS Page Rules
func (sess clientSession) CisPageRuleClientSession() (*cispagerulev1.PageRuleApiV1, error) {
	if sess.cisPageRuleErr != nil {
		return sess.cisPageRuleClient, sess.cisPageRuleErr
	}
	return sess.cisPageRuleClient.Clone(), nil
}

// CIS Edge Function
func (sess clientSession) CisEdgeFunctionClientSession() (*cisedgefunctionv1.EdgeFunctionsApiV1, error) {
	if sess.cisEdgeFunctionErr != nil {
		return sess.cisEdgeFunctionClient, sess.cisEdgeFunctionErr
	}
	return sess.cisEdgeFunctionClient.Clone(), nil
}

// CIS SSL certificate
func (sess clientSession) CisSSLClientSession() (*cissslv1.SslCertificateApiV1, error) {
	if sess.cisSSLErr != nil {
		return sess.cisSSLClient, sess.cisSSLErr
	}
	return sess.cisSSLClient.Clone(), nil
}

// CIS WAF Packages
func (sess clientSession) CisWAFPackageClientSession() (*ciswafpackagev1.WafRulePackagesApiV1, error) {
	if sess.cisWAFPackageErr != nil {
		return sess.cisWAFPackageClient, sess.cisWAFPackageErr
	}
	return sess.cisWAFPackageClient.Clone(), nil
}

// CIS Zone Settings
func (sess clientSession) CisDomainSettingsClientSession() (*cisdomainsettingsv1.ZonesSettingsV1, error) {
	if sess.cisDomainSettingsErr != nil {
		return sess.cisDomainSettingsClient, sess.cisDomainSettingsErr
	}
	return sess.cisDomainSettingsClient.Clone(), nil
}

// CIS Routing
func (sess clientSession) CisRoutingClientSession() (*cisroutingv1.RoutingV1, error) {
	if sess.cisRoutingErr != nil {
		return sess.cisRoutingClient, sess.cisRoutingErr
	}
	return sess.cisRoutingClient.Clone(), nil
}

// CIS WAF Group
func (sess clientSession) CisWAFGroupClientSession() (*ciswafgroupv1.WafRuleGroupsApiV1, error) {
	if sess.cisWAFGroupErr != nil {
		return sess.cisWAFGroupClient, sess.cisWAFGroupErr
	}
	return sess.cisWAFGroupClient.Clone(), nil
}

// CIS Cache service
func (sess clientSession) CisCacheClientSession() (*ciscachev1.CachingApiV1, error) {
	if sess.cisCacheErr != nil {
		return sess.cisCacheClient, sess.cisCacheErr
	}
	return sess.cisCacheClient.Clone(), nil
}

// CIS Zone Settings
func (sess clientSession) CisCustomPageClientSession() (*ciscustompagev1.CustomPagesV1, error) {
	if sess.cisCustomPageErr != nil {
		return sess.cisCustomPageClient, sess.cisCustomPageErr
	}
	return sess.cisCustomPageClient.Clone(), nil
}

// CIS Firewall access rule
func (sess clientSession) CisAccessRuleClientSession() (*cisaccessrulev1.ZoneFirewallAccessRulesV1, error) {
	if sess.cisAccessRuleErr != nil {
		return sess.cisAccessRuleClient, sess.cisAccessRuleErr
	}
	return sess.cisAccessRuleClient.Clone(), nil
}

// CIS User Agent Blocking rule
func (sess clientSession) CisUARuleClientSession() (*cisuarulev1.UserAgentBlockingRulesV1, error) {
	if sess.cisUARuleErr != nil {
		return sess.cisUARuleClient, sess.cisUARuleErr
	}
	return sess.cisUARuleClient.Clone(), nil
}

// CIS Firewall Lockdown rule
func (sess clientSession) CisLockdownClientSession() (*cislockdownv1.ZoneLockdownV1, error) {
	if sess.cisLockdownErr != nil {
		return sess.cisLockdownClient, sess.cisLockdownErr
	}
	return sess.cisLockdownClient.Clone(), nil
}

// CIS Range app rule
func (sess clientSession) CisRangeAppClientSession() (*cisrangeappv1.RangeApplicationsV1, error) {
	if sess.cisRangeAppErr != nil {
		return sess.cisRangeAppClient, sess.cisRangeAppErr
	}
	return sess.cisRangeAppClient.Clone(), nil
}

// CIS WAF Rule
func (sess clientSession) CisWAFRuleClientSession() (*ciswafrulev1.WafRulesApiV1, error) {
	if sess.cisWAFRuleErr != nil {
		return sess.cisWAFRuleClient, sess.cisWAFRuleErr
	}
	return sess.cisWAFRuleClient.Clone(), nil
}

// IAM Identity Session
func (sess clientSession) IAMIdentityV1API() (*iamidentity.IamIdentityV1, error) {
	return sess.iamIdentityAPI, sess.iamIdentityErr
}

// ResourceMAanger Session
func (sess clientSession) ResourceManagerV2API() (*resourcemanager.ResourceManagerV2, error) {
	return sess.resourceManagerAPI, sess.resourceManagerErr
}

func (session clientSession) EnterpriseManagementV1() (*enterprisemanagementv1.EnterpriseManagementV1, error) {
	return session.enterpriseManagementClient, session.enterpriseManagementClientErr
}

// ResourceController Session
func (sess clientSession) ResourceControllerV2API() (*resourcecontroller.ResourceControllerV2, error) {
	return sess.resourceControllerAPI, sess.resourceControllerErr
}

// SecretsManager Session
func (session clientSession) SecretsManagerV1() (*secretsmanagerv1.SecretsManagerV1, error) {
	return session.secretsManagerClient, session.secretsManagerClientErr
}

var cloudEndpoint = "cloud.ibm.com"

// ClientSession configures and returns a fully initialized ClientSession
func (c *Config) ClientSession() (interface{}, error) {
	sess, err := newSession(c)
	if err != nil {
		return nil, err
	}
	log.Printf("[INFO] Configured Region: %s\n", c.Region)
	session := clientSession{
		session: sess,
	}

	if sess.BluemixSession == nil {
		//Can be nil only  if bluemix_api_key is not provided
		log.Println("Skipping Bluemix Clients configuration")
		session.bluemixSessionErr = errEmptyBluemixCredentials
		session.accountConfigErr = errEmptyBluemixCredentials
		session.accountV1ConfigErr = errEmptyBluemixCredentials
		session.csConfigErr = errEmptyBluemixCredentials
		session.csv2ConfigErr = errEmptyBluemixCredentials
		session.containerRegistryClientErr = errEmptyBluemixCredentials
		session.kpErr = errEmptyBluemixCredentials
		session.pushServiceClientErr = errEmptyBluemixCredentials
		session.appConfigurationClientErr = errEmptyBluemixCredentials
		session.kmsErr = errEmptyBluemixCredentials
		session.cfConfigErr = errEmptyBluemixCredentials
		session.cisConfigErr = errEmptyBluemixCredentials
		session.functionConfigErr = errEmptyBluemixCredentials
		session.globalSearchConfigErr = errEmptyBluemixCredentials
		session.globalTaggingConfigErr = errEmptyBluemixCredentials
		session.globalTaggingConfigErrV1 = errEmptyBluemixCredentials
		session.hpcsEndpointErr = errEmptyBluemixCredentials
		session.iamConfigErr = errEmptyBluemixCredentials
		session.iamUUMConfigErrV2 = errEmptyBluemixCredentials
		session.icdConfigErr = errEmptyBluemixCredentials
		session.resourceCatalogConfigErr = errEmptyBluemixCredentials
		session.resourceManagerErr = errEmptyBluemixCredentials
		session.resourceManagementConfigErrv2 = errEmptyBluemixCredentials
		session.resourceControllerConfigErr = errEmptyBluemixCredentials
		session.resourceControllerConfigErrv2 = errEmptyBluemixCredentials
		session.enterpriseManagementClientErr = errEmptyBluemixCredentials
		session.resourceControllerErr = errEmptyBluemixCredentials
		session.catalogManagementClientErr = errEmptyBluemixCredentials
		session.powerConfigErr = errEmptyBluemixCredentials
		session.ibmpiConfigErr = errEmptyBluemixCredentials
		session.userManagementErr = errEmptyBluemixCredentials
		session.certManagementErr = errEmptyBluemixCredentials
		session.vpcClassicErr = errEmptyBluemixCredentials
		session.vpcErr = errEmptyBluemixCredentials
		session.apigatewayErr = errEmptyBluemixCredentials
		session.pDNSErr = errEmptyBluemixCredentials
		session.bmxUserFetchErr = errEmptyBluemixCredentials
		session.directlinkErr = errEmptyBluemixCredentials
		session.dlProviderErr = errEmptyBluemixCredentials
		session.cosConfigErr = errEmptyBluemixCredentials
		session.transitgatewayErr = errEmptyBluemixCredentials
		session.functionIAMNamespaceErr = errEmptyBluemixCredentials
		session.cisDNSErr = errEmptyBluemixCredentials
		session.cisDNSBulkErr = errEmptyBluemixCredentials
		session.cisGLBPoolErr = errEmptyBluemixCredentials
		session.cisGLBErr = errEmptyBluemixCredentials
		session.cisGLBHealthCheckErr = errEmptyBluemixCredentials
		session.cisIPErr = errEmptyBluemixCredentials
		session.cisZonesErr = errEmptyBluemixCredentials
		session.cisRLErr = errEmptyBluemixCredentials
		session.cisPageRuleErr = errEmptyBluemixCredentials
		session.cisEdgeFunctionErr = errEmptyBluemixCredentials
		session.cisSSLErr = errEmptyBluemixCredentials
		session.cisWAFPackageErr = errEmptyBluemixCredentials
		session.cisDomainSettingsErr = errEmptyBluemixCredentials
		session.cisRoutingErr = errEmptyBluemixCredentials
		session.cisWAFGroupErr = errEmptyBluemixCredentials
		session.cisCacheErr = errEmptyBluemixCredentials
		session.cisCustomPageErr = errEmptyBluemixCredentials
		session.cisAccessRuleErr = errEmptyBluemixCredentials
		session.cisUARuleErr = errEmptyBluemixCredentials
		session.cisLockdownErr = errEmptyBluemixCredentials
		session.cisRangeAppErr = errEmptyBluemixCredentials
		session.cisWAFRuleErr = errEmptyBluemixCredentials
		session.iamIdentityErr = errEmptyBluemixCredentials
		session.secretsManagerClientErr = errEmptyBluemixCredentials
		session.schematicsClientErr = errEmptyBluemixCredentials
		session.iamPolicyManagementErr = errEmptyBluemixCredentials

		return session, nil
	}

	if sess.BluemixSession.Config.BluemixAPIKey != "" {
		err = authenticateAPIKey(sess.BluemixSession)
		if err != nil {
			for count := c.RetryCount; count >= 0; count-- {
				if err == nil || !isRetryable(err) {
					break
				}
				time.Sleep(c.RetryDelay)
				log.Printf("Retrying IAM Authentication %d", count)
				err = authenticateAPIKey(sess.BluemixSession)
			}
			if err != nil {
				session.bmxUserFetchErr = fmt.Errorf("Error occured while fetching auth key for account user details: %q", err)
				session.functionConfigErr = fmt.Errorf("Error occured while fetching auth key for function: %q", err)
				session.powerConfigErr = fmt.Errorf("Error occured while fetching the auth key for power iaas: %q", err)
				session.ibmpiConfigErr = fmt.Errorf("Error occured while fetching the auth key for power iaas: %q", err)
			}
		}
		err = authenticateCF(sess.BluemixSession)
		if err != nil {
			for count := c.RetryCount; count >= 0; count-- {
				if err == nil || !isRetryable(err) {
					break
				}
				time.Sleep(c.RetryDelay)
				log.Printf("Retrying CF Authentication %d", count)
				err = authenticateCF(sess.BluemixSession)
			}
			if err != nil {
				session.functionConfigErr = fmt.Errorf("Error occured while fetching auth key for function: %q", err)
			}
		}
	}

	if sess.BluemixSession.Config.IAMAccessToken != "" && sess.BluemixSession.Config.BluemixAPIKey == "" {
		err := refreshToken(sess.BluemixSession)
		if err != nil {
			for count := c.RetryCount; count >= 0; count-- {
				if err == nil || !isRetryable(err) {
					break
				}
				time.Sleep(c.RetryDelay)
				log.Printf("Retrying refresh token %d", count)
				err = refreshToken(sess.BluemixSession)
			}
			if err != nil {
				return nil, fmt.Errorf("Error occured while refreshing the token: %q", err)
			}
		}

	}
	userConfig, err := fetchUserDetails(sess.BluemixSession, c.RetryCount, c.RetryDelay)
	if err != nil {
		session.bmxUserFetchErr = fmt.Errorf("Error occured while fetching account user details: %q", err)
	}
	session.bmxUserDetails = userConfig

	if sess.SoftLayerSession != nil && sess.SoftLayerSession.IAMToken != "" {
		sess.SoftLayerSession.IAMToken = sess.BluemixSession.Config.IAMAccessToken
	}

	session.functionClient, session.functionConfigErr = FunctionClient(sess.BluemixSession.Config)

	BluemixRegion = sess.BluemixSession.Config.Region

	accv1API, err := accountv1.New(sess.BluemixSession)
	if err != nil {
		session.accountV1ConfigErr = fmt.Errorf("Error occured while configuring Bluemix Accountv1 Service: %q", err)
	}
	session.bmxAccountv1ServiceAPI = accv1API

	accAPI, err := accountv2.New(sess.BluemixSession)
	if err != nil {
		session.accountConfigErr = fmt.Errorf("Error occured while configuring  Account Service: %q", err)
	}
	session.bmxAccountServiceAPI = accAPI

	cfAPI, err := mccpv2.New(sess.BluemixSession)
	if err != nil {
		session.cfConfigErr = fmt.Errorf("Error occured while configuring MCCP service: %q", err)
	}
	session.cfServiceAPI = cfAPI

	clusterAPI, err := containerv1.New(sess.BluemixSession)
	if err != nil {
		session.csConfigErr = fmt.Errorf("Error occured while configuring Container Service for K8s cluster: %q", err)
	}
	session.csServiceAPI = clusterAPI

	v2clusterAPI, err := containerv2.New(sess.BluemixSession)
	if err != nil {
		session.csv2ConfigErr = fmt.Errorf("Error occured while configuring vpc Container Service for K8s cluster: %q", err)
	}
	session.csv2ServiceAPI = v2clusterAPI

	hpcsAPI, err := hpcs.New(sess.BluemixSession)
	if err != nil {
		session.hpcsEndpointErr = fmt.Errorf("Error occured while configuring hpcs Endpoint: %q", err)
	}
	session.hpcsEndpointAPI = hpcsAPI

	kpurl := contructEndpoint(fmt.Sprintf("%s.kms", c.Region), cloudEndpoint)
	if c.Visibility == "private" || c.Visibility == "public-and-private" {
		kpurl = contructEndpoint(fmt.Sprintf("private.%s.kms", c.Region), cloudEndpoint)
	}
	var options kp.ClientConfig
	if c.BluemixAPIKey != "" {
		options = kp.ClientConfig{
			BaseURL: envFallBack([]string{"IBMCLOUD_KP_API_ENDPOINT"}, kpurl),
			APIKey:  sess.BluemixSession.Config.BluemixAPIKey, //pragma: allowlist secret
			// InstanceID:    "42fET57nnadurKXzXAedFLOhGqETfIGYxOmQXkFgkJV9",
			Verbose: kp.VerboseFailOnly,
		}

	} else {
		options = kp.ClientConfig{
			BaseURL:       envFallBack([]string{"IBMCLOUD_KP_API_ENDPOINT"}, kpurl),
			Authorization: sess.BluemixSession.Config.IAMAccessToken,
			// InstanceID:    "42fET57nnadurKXzXAedFLOhGqETfIGYxOmQXkFgkJV9",
			Verbose: kp.VerboseFailOnly,
		}
	}
	kpAPIclient, err := kp.New(options, kp.DefaultTransport())
	if err != nil {
		session.kpErr = fmt.Errorf("Error occured while configuring Key Protect Service: %q", err)
	}
	session.kpAPI = kpAPIclient

	kmsurl := contructEndpoint(fmt.Sprintf("%s.kms", c.Region), cloudEndpoint)
	if c.Visibility == "private" || c.Visibility == "public-and-private" {
		kmsurl = contructEndpoint(fmt.Sprintf("private.%s.kms", c.Region), cloudEndpoint)
	}
	var kmsOptions kp.ClientConfig
	if c.BluemixAPIKey != "" {
		kmsOptions = kp.ClientConfig{
			BaseURL: envFallBack([]string{"IBMCLOUD_KP_API_ENDPOINT"}, kmsurl),
			APIKey:  sess.BluemixSession.Config.BluemixAPIKey, //pragma: allowlist secret
			// InstanceID:    "5af62d5d-5d90-4b84-bbcd-90d2123ae6c8",
			Verbose: kp.VerboseFailOnly,
		}

	} else {
		kmsOptions = kp.ClientConfig{
			BaseURL:       envFallBack([]string{"IBMCLOUD_KP_API_ENDPOINT"}, kmsurl),
			Authorization: sess.BluemixSession.Config.IAMAccessToken,
			// InstanceID:    "5af62d5d-5d90-4b84-bbcd-90d2123ae6c8",
			Verbose: kp.VerboseFailOnly,
		}
	}
	kmsAPIclient, err := kp.New(kmsOptions, DefaultTransport())
	if err != nil {
		session.kmsErr = fmt.Errorf("Error occured while configuring key Service: %q", err)
	}
	session.kmsAPI = kmsAPIclient

	var authenticator core.Authenticator

	if c.BluemixAPIKey != "" {
		authenticator = &core.IamAuthenticator{
			ApiKey: c.BluemixAPIKey,
			URL:    envFallBack([]string{"IBMCLOUD_IAM_API_ENDPOINT"}, "https://iam.cloud.ibm.com") + "/identity/token",
		}
	} else if strings.HasPrefix(sess.BluemixSession.Config.IAMAccessToken, "Bearer") {
		authenticator = &core.BearerTokenAuthenticator{
			BearerToken: sess.BluemixSession.Config.IAMAccessToken[7:],
		}
	} else {
		authenticator = &core.BearerTokenAuthenticator{
			BearerToken: sess.BluemixSession.Config.IAMAccessToken,
		}
	}

	// Construct an "options" struct for creating the service client.
	catalogManagementURL := "https://cm.globalcatalog.cloud.ibm.com/api/v1-beta"
	if c.Visibility == "private" {
		session.catalogManagementClientErr = fmt.Errorf("Catalog Management resource doesnot support private endpoints")
	}
	catalogManagementClientOptions := &catalogmanagementv1.CatalogManagementV1Options{
		URL:           envFallBack([]string{"IBMCLOUD_CATALOG_MANAGEMENT_API_ENDPOINT"}, catalogManagementURL),
		Authenticator: authenticator,
	}

	// Construct the service client.
	session.catalogManagementClient, err = catalogmanagementv1.NewCatalogManagementV1(catalogManagementClientOptions)
	if err == nil {
		// Enable retries for API calls
		session.catalogManagementClient.Service.EnableRetries(c.RetryCount, c.RetryDelay)
		// Add custom header for analytics
		session.catalogManagementClient.SetDefaultHeaders(gohttp.Header{
			"X-Original-User-Agent": {fmt.Sprintf("terraform-provider-ibm/%s", version.Version)},
		})
	} else {
		session.catalogManagementClientErr = fmt.Errorf("Error occurred while configuring Catalog Management API service: %q", err)
	}
	schematicsEndpoint := "https://schematics.cloud.ibm.com"
	if c.Visibility == "private" || c.Visibility == "public-and-private" {
		if c.Region == "us-south" || c.Region == "us-east" {
			schematicsEndpoint = contructEndpoint("private-us.schematics", cloudEndpoint)
		} else if c.Region == "eu-gb" || c.Region == "eu-de" {
			schematicsEndpoint = contructEndpoint("private-eu.schematics", cloudEndpoint)
		} else {
			schematicsEndpoint = "https://schematics.cloud.ibm.com"
		}
	}
	schematicsClientOptions := &schematicsv1.SchematicsV1Options{
		Authenticator: authenticator,
		URL:           envFallBack([]string{"IBMCLOUD_SCHEMATICS_API_ENDPOINT"}, schematicsEndpoint),
	}

	// Construct the service client.
	schematicsClient, err := schematicsv1.NewSchematicsV1(schematicsClientOptions)
	// Enable retries for API calls
	if schematicsClient != nil && schematicsClient.Service != nil {
		schematicsClient.Service.EnableRetries(c.RetryCount, c.RetryDelay)
		if err != nil {
			session.schematicsClientErr = fmt.Errorf("Error occurred while configuring Schematics Service API service: %q", err)
		}
	}
	session.schematicsClient = schematicsClient

	vpcclassicurl := contructEndpoint(fmt.Sprintf("%s.iaas", c.Region), fmt.Sprintf("%s/v1", cloudEndpoint))
	if c.Visibility == "private" {
		if c.Region == "us-south" || c.Region == "us-east" {
			vpcclassicurl = contructEndpoint(fmt.Sprintf("%s.private.iaas", c.Region), fmt.Sprintf("%s/v1", cloudEndpoint))
		} else {
			session.vpcClassicErr = fmt.Errorf("VPC Classic supports private endpoints only in us-south and us-east")
		}
	}
	if c.Visibility == "public-and-private" {
		if c.Region == "us-south" || c.Region == "us-east" {
			vpcclassicurl = contructEndpoint(fmt.Sprintf("%s.private.iaas", c.Region), fmt.Sprintf("%s/v1", cloudEndpoint))
		} else {
			vpcclassicurl = contructEndpoint(fmt.Sprintf("%s.iaas", c.Region), fmt.Sprintf("%s/v1", cloudEndpoint))
		}
	}
	vpcclassicoptions := &vpcclassic.VpcClassicV1Options{
		URL:           envFallBack([]string{"IBMCLOUD_IS_API_ENDPOINT"}, vpcclassicurl),
		Authenticator: authenticator,
	}
	vpcclassicclient, err := vpcclassic.NewVpcClassicV1(vpcclassicoptions)
	if err != nil {
		session.vpcErr = fmt.Errorf("Error occured while configuring vpc classic service: %q", err)
	}
	if vpcclassicclient != nil && vpcclassicclient.Service != nil {
		vpcclassicclient.Service.EnableRetries(c.RetryCount, c.RetryDelay)
	}

	session.vpcClassicAPI = vpcclassicclient

	vpcurl := contructEndpoint(fmt.Sprintf("%s.iaas", c.Region), fmt.Sprintf("%s/v1", cloudEndpoint))
	if c.Visibility == "private" {
		if c.Region == "us-south" || c.Region == "us-east" {
			vpcurl = contructEndpoint(fmt.Sprintf("%s.private.iaas", c.Region), fmt.Sprintf("%s/v1", cloudEndpoint))
		} else {
			session.vpcErr = fmt.Errorf("VPC supports private endpoints only in us-south and us-east")
		}
	}
	if c.Visibility == "public-and-private" {
		if c.Region == "us-south" || c.Region == "us-east" {
			vpcurl = contructEndpoint(fmt.Sprintf("%s.private.iaas", c.Region), fmt.Sprintf("%s/v1", cloudEndpoint))
		}
		vpcurl = contructEndpoint(fmt.Sprintf("%s.iaas", c.Region), fmt.Sprintf("%s/v1", cloudEndpoint))
	}
	vpcoptions := &vpc.VpcV1Options{
		URL:           envFallBack([]string{"IBMCLOUD_IS_NG_API_ENDPOINT"}, vpcurl),
		Authenticator: authenticator,
	}
	vpcclient, err := vpc.NewVpcV1(vpcoptions)
	if err != nil {
		session.vpcErr = fmt.Errorf("Error occured while configuring vpc service: %q", err)
	}
	if vpcclient != nil && vpcclient.Service != nil {
		vpcclient.Service.EnableRetries(c.RetryCount, c.RetryDelay)
	}
	session.vpcAPI = vpcclient

	pnurl := fmt.Sprintf("https://%s.imfpush.cloud.ibm.com/imfpush/v1", c.Region)
	if c.Visibility == "private" {
		session.pushServiceClientErr = fmt.Errorf("Push Notifications Service API doesnot support private endpoints")
	}
	pushNotificationOptions := &pushservicev1.PushServiceV1Options{
		URL:           envFallBack([]string{"IBMCLOUD_PUSH_API_ENDPOINT"}, pnurl),
		Authenticator: authenticator,
	}
	pnclient, err := pushservicev1.NewPushServiceV1(pushNotificationOptions)
	if pnclient != nil {
		// Enable retries for API calls
		pnclient.EnableRetries(c.RetryCount, c.RetryDelay)
		session.pushServiceClient = pnclient
	} else {
		session.pushServiceClientErr = fmt.Errorf("Error occured while configuring push notification service: %q", err)
	}
	if c.Visibility == "private" {
		session.appConfigurationClientErr = fmt.Errorf("App Configuration Service API doesnot support private endpoints")
	}
	appConfigurationClientOptions := &appconfigurationv1.AppConfigurationV1Options{
		Authenticator: authenticator,
	}
	appConfigClient, err := appconfigurationv1.NewAppConfigurationV1(appConfigurationClientOptions)
	if appConfigClient != nil {
		// Enable retries for API calls
		appConfigClient.Service.EnableRetries(c.RetryCount, c.RetryDelay)
		session.appConfigurationClient = appConfigClient
	} else {
		session.appConfigurationClientErr = fmt.Errorf("Error occurred while configuring App Configuration service: %q", err)
	}
	// Construct an "options" struct for creating the service client.
	containerRegistryClientURL, err := containerregistryv1.GetServiceURLForRegion(c.Region)
	if err != nil {
		containerRegistryClientURL = containerregistryv1.DefaultServiceURL
	}
	if c.Visibility == "private" || c.Visibility == "public-and-private" {
		containerRegistryClientURL, err = GetPrivateServiceURLForRegion(c.Region)
		if err != nil {
			containerRegistryClientURL, _ = GetPrivateServiceURLForRegion("us-south")
		}
	}
	containerRegistryClientOptions := &containerregistryv1.ContainerRegistryV1Options{
		Authenticator: authenticator,
		URL:           envFallBack([]string{"IBMCLOUD_CR_API_ENDPOINT"}, containerRegistryClientURL),
		Account:       core.StringPtr(userConfig.userAccount),
	}

	// Construct the service client.
	session.containerRegistryClient, err = containerregistryv1.NewContainerRegistryV1(containerRegistryClientOptions)
	if err == nil {
		// Enable retries for API calls
		session.containerRegistryClient.Service.EnableRetries(c.RetryCount, c.RetryDelay)
		// Add custom header for analytics
		session.containerRegistryClient.SetDefaultHeaders(gohttp.Header{
			"X-Original-User-Agent": {fmt.Sprintf("terraform-provider-ibm/%s", version.Version)},
		})
	} else {
		session.containerRegistryClientErr = fmt.Errorf("Error occurred while configuring IBM Cloud Container Registry API service: %q", err)
	}

	//cosconfigurl := fmt.Sprintf("https://%s.iaas.cloud.ibm.com/v1", c.Region)
	cosconfigoptions := &cosconfig.ResourceConfigurationV1Options{
		Authenticator: authenticator,
		URL:           envFallBack([]string{"IBMCLOUD_COS_CONFIG_ENDPOINT"}, "https://config.cloud-object-storage.cloud.ibm.com/v1"),
	}
	cosconfigclient, err := cosconfig.NewResourceConfigurationV1(cosconfigoptions)
	if err != nil {
		session.cosConfigErr = fmt.Errorf("Error occured while configuring COS config service: %q", err)
	}
	session.cosConfigAPI = cosconfigclient

	cisAPI, err := cisv1.New(sess.BluemixSession)
	if err != nil {
		session.cisConfigErr = fmt.Errorf("Error occured while configuring Cloud Internet Services: %q", err)
	}
	session.cisServiceAPI = cisAPI

	globalSearchAPI, err := globalsearchv2.New(sess.BluemixSession)
	if err != nil {
		session.globalSearchConfigErr = fmt.Errorf("Error occured while configuring Global Search: %q", err)
	}
	session.globalSearchServiceAPI = globalSearchAPI

	globalTaggingAPI, err := globaltaggingv3.New(sess.BluemixSession)
	if err != nil {
		session.globalTaggingConfigErr = fmt.Errorf("Error occured while configuring Global Tagging: %q", err)
	}
	session.globalTaggingServiceAPI = globalTaggingAPI

	globalTaggingEndpoint := "https://tags.global-search-tagging.cloud.ibm.com"
	if c.Visibility == "private" || c.Visibility == "public-and-private" {
		var globalTaggingRegion string
		if c.Region != "us-south" && c.Region != "us-east" {
			globalTaggingRegion = "us-south"
		} else {
			globalTaggingRegion = c.Region
		}
		globalTaggingEndpoint = contructEndpoint(fmt.Sprintf("tags.private.%s", globalTaggingRegion), fmt.Sprintf("global-search-tagging.%s", cloudEndpoint))
	}

	globalTaggingV1Options := &globaltaggingv1.GlobalTaggingV1Options{
		URL:           envFallBack([]string{"IBMCLOUD_GT_API_ENDPOINT"}, globalTaggingEndpoint),
		Authenticator: authenticator,
	}

	globalTaggingAPIV1, err := globaltaggingv1.NewGlobalTaggingV1(globalTaggingV1Options)
	if err != nil {
		session.globalTaggingConfigErrV1 = fmt.Errorf("Error occured while configuring Global Tagging: %q", err)
	}
	if globalTaggingAPIV1 != nil {
		session.globalTaggingServiceAPIV1 = *globalTaggingAPIV1
		session.globalTaggingServiceAPIV1.Service.EnableRetries(c.RetryCount, c.RetryDelay)
	}

	iam, err := iamv1.New(sess.BluemixSession)
	if err != nil {
		session.iamConfigErr = fmt.Errorf("Error occured while configuring Bluemix IAM Service: %q", err)
	}
	session.iamServiceAPI = iam

	iamuumv2, err := iamuumv2.New(sess.BluemixSession)
	if err != nil {
		session.iamUUMConfigErrV2 = fmt.Errorf("Error occured while configuring Bluemix IAMUUM Service: %q", err)
	}
	session.iamUUMServiceAPIV2 = iamuumv2

	icdAPI, err := icdv4.New(sess.BluemixSession)
	if err != nil {
		session.icdConfigErr = fmt.Errorf("Error occured while configuring IBM Cloud Database Services: %q", err)
	}
	session.icdServiceAPI = icdAPI

	resourceCatalogAPI, err := catalog.New(sess.BluemixSession)
	if err != nil {
		session.resourceCatalogConfigErr = fmt.Errorf("Error occured while configuring Resource Catalog service: %q", err)
	}
	session.resourceCatalogServiceAPI = resourceCatalogAPI

	resourceManagementAPIv2, err := managementv2.New(sess.BluemixSession)
	if err != nil {
		session.resourceManagementConfigErrv2 = fmt.Errorf("Error occured while configuring Resource Management service: %q", err)
	}
	session.resourceManagementServiceAPIv2 = resourceManagementAPIv2

	resourceControllerAPI, err := controller.New(sess.BluemixSession)
	if err != nil {
		session.resourceControllerConfigErr = fmt.Errorf("Error occured while configuring Resource Controller service: %q", err)
	}
	session.resourceControllerServiceAPI = resourceControllerAPI

	ResourceControllerAPIv2, err := controllerv2.New(sess.BluemixSession)
	if err != nil {
		session.resourceControllerConfigErrv2 = fmt.Errorf("Error occured while configuring Resource Controller v2 service: %q", err)
	}
	session.resourceControllerServiceAPIv2 = ResourceControllerAPIv2

	userManagementAPI, err := usermanagementv2.New(sess.BluemixSession)
	if err != nil {
		session.userManagementErr = fmt.Errorf("Error occured while configuring user management service: %q", err)
	}
	session.userManagementAPI = userManagementAPI
	certManagementAPI, err := certificatemanager.New(sess.BluemixSession)
	if err != nil {
		session.certManagementErr = fmt.Errorf("Error occured while configuring Certificate manager service: %q", err)
	}
	session.certManagementAPI = certManagementAPI

	namespaceFunction, err := functions.New(sess.BluemixSession)
	if err != nil {
		session.functionIAMNamespaceErr = fmt.Errorf("Error occured while configuring Cloud Funciton Service : %q", err)
	}
	session.functionIAMNamespaceAPI = namespaceFunction

	apicurl := contructEndpoint(fmt.Sprintf("api.%s.apigw", c.Region), fmt.Sprintf("%s/controller", cloudEndpoint))
	if c.Visibility == "private" || c.Visibility == "public-and-private" {
		apicurl = contructEndpoint(fmt.Sprintf("api.private.%s.apigw", c.Region), fmt.Sprintf("%s/controller", cloudEndpoint))
	}
	APIGatewayControllerAPIV1Options := &apigateway.ApiGatewayControllerApiV1Options{
		URL:           envFallBack([]string{"IBMCLOUD_API_GATEWAY_ENDPOINT"}, apicurl),
		Authenticator: &core.NoAuthAuthenticator{},
	}
	apigatewayAPI, err := apigateway.NewApiGatewayControllerApiV1(APIGatewayControllerAPIV1Options)
	if err != nil {
		session.apigatewayErr = fmt.Errorf("Error occured while configuring  APIGateway service: %q", err)
	}
	session.apigatewayAPI = apigatewayAPI

	ibmpisession, err := ibmpisession.New(sess.BluemixSession.Config.IAMAccessToken, c.Region, false, 90000000000, session.bmxUserDetails.userAccount, c.Zone)
	if err != nil {
		session.ibmpiConfigErr = err
		return nil, err
	}

	session.ibmpiSession = ibmpisession

	pdnsURL := dns.DefaultServiceURL
	if c.Visibility == "private" || c.Visibility == "public-and-private" {
		pdnsURL = contructEndpoint("api.private.dns-svcs", fmt.Sprintf("%s/v1", cloudEndpoint))
	}
	dnsOptions := &dns.DnsSvcsV1Options{
		URL:           envFallBack([]string{"IBMCLOUD_PRIVATE_DNS_API_ENDPOINT"}, pdnsURL),
		Authenticator: authenticator,
	}

	session.pDNSClient, session.pDNSErr = dns.NewDnsSvcsV1(dnsOptions)
	if session.pDNSErr != nil {
		session.pDNSErr = fmt.Errorf("Error occured while configuring PrivateDNS Service: %s", session.pDNSErr)
	}
	if session.pDNSClient != nil && session.pDNSClient.Service != nil {
		session.pDNSClient.Service.EnableRetries(c.RetryCount, c.RetryDelay)
	}

	ver := time.Now().Format("2006-01-02")

	dlURL := dl.DefaultServiceURL
	if c.Visibility == "private" || c.Visibility == "public-and-private" {
		dlURL = contructEndpoint("private.directlink", fmt.Sprintf("%s/v1", cloudEndpoint))
	}
	directlinkOptions := &dl.DirectLinkV1Options{
		URL:           envFallBack([]string{"IBMCLOUD_DL_API_ENDPOINT"}, dlURL),
		Authenticator: authenticator,
		Version:       &ver,
	}

	session.directlinkAPI, session.directlinkErr = dl.NewDirectLinkV1(directlinkOptions)
	if session.directlinkErr != nil {
		session.directlinkErr = fmt.Errorf("Error occured while configuring Direct Link Service: %s", session.directlinkErr)
	}
	if session.directlinkAPI != nil && session.directlinkAPI.Service != nil {
		session.directlinkAPI.Service.EnableRetries(c.RetryCount, c.RetryDelay)
	}

	//Direct link provider
	dlproviderURL := dlProviderV2.DefaultServiceURL
	if c.Visibility == "private" || c.Visibility == "public-and-private" {
		dlproviderURL = contructEndpoint("private.directlink", fmt.Sprintf("%s/provider/v2", cloudEndpoint))
	}
	directLinkProviderV2Options := &dlProviderV2.DirectLinkProviderV2Options{
		URL:           envFallBack([]string{"IBMCLOUD_DL_PROVIDER_API_ENDPOINT"}, dlproviderURL),
		Authenticator: authenticator,
		Version:       &ver,
	}

	session.dlProviderAPI, session.dlProviderErr = dlProviderV2.NewDirectLinkProviderV2(directLinkProviderV2Options)
	if session.dlProviderErr != nil {
		session.dlProviderErr = fmt.Errorf("Error occured while configuring Direct Link Provider Service: %s", session.dlProviderErr)
	}
	if session.dlProviderAPI != nil && session.dlProviderAPI.Service != nil {
		session.dlProviderAPI.Service.EnableRetries(c.RetryCount, c.RetryDelay)
	}

	tgURL := tg.DefaultServiceURL
	if c.Visibility == "private" || c.Visibility == "public-and-private" {
		tgURL = contructEndpoint("private.transit", fmt.Sprintf("%s/v1", cloudEndpoint))
	}
	transitgatewayOptions := &tg.TransitGatewayApisV1Options{
		URL:           envFallBack([]string{"IBMCLOUD_TG_API_ENDPOINT"}, tgURL),
		Authenticator: authenticator,
		Version:       CreateVersionDate(),
	}

	session.transitgatewayAPI, session.transitgatewayErr = tg.NewTransitGatewayApisV1(transitgatewayOptions)
	if session.transitgatewayErr != nil {
		session.transitgatewayErr = fmt.Errorf("Error occured while configuring Transit Gateway Service: %s", session.transitgatewayErr)
	}
	if session.transitgatewayAPI != nil && session.transitgatewayAPI.Service != nil {
		session.transitgatewayAPI.Service.EnableRetries(c.RetryCount, c.RetryDelay)
	}

	// CIS Service instances starts here.
	cisURL := contructEndpoint("api.cis", cloudEndpoint)
	if c.Visibility == "private" {
		// cisURL = contructEndpoint("api.private.cis", cloudEndpoint)
		session.cisZonesErr = fmt.Errorf("CIS Service doesnt support private endpoints.")
		session.cisDNSBulkErr = fmt.Errorf("CIS Service doesnt support private endpoints.")
		session.cisGLBPoolErr = fmt.Errorf("CIS Service doesnt support private endpoints.")
		session.cisGLBErr = fmt.Errorf("CIS Service doesnt support private endpoints.")
		session.cisGLBHealthCheckErr = fmt.Errorf("CIS Service doesnt support private endpoints.")
		session.cisIPErr = fmt.Errorf("CIS Service doesnt support private endpoints.")
		session.cisRLErr = fmt.Errorf("CIS Service doesnt support private endpoints.")
		session.cisPageRuleErr = fmt.Errorf("CIS Service doesnt support private endpoints.")
		session.cisEdgeFunctionErr = fmt.Errorf("CIS Service doesnt support private endpoints.")
		session.cisSSLErr = fmt.Errorf("CIS Service doesnt support private endpoints.")
		session.cisWAFPackageErr = fmt.Errorf("CIS Service doesnt support private endpoints.")
		session.cisDomainSettingsErr = fmt.Errorf("CIS Service doesnt support private endpoints.")
		session.cisRoutingErr = fmt.Errorf("CIS Service doesnt support private endpoints.")
		session.cisWAFGroupErr = fmt.Errorf("CIS Service doesnt support private endpoints.")
		session.cisCacheErr = fmt.Errorf("CIS Service doesnt support private endpoints.")
		session.cisCustomPageErr = fmt.Errorf("CIS Service doesnt support private endpoints.")
		session.cisAccessRuleErr = fmt.Errorf("CIS Service doesnt support private endpoints.")
		session.cisUARuleErr = fmt.Errorf("CIS Service doesnt support private endpoints.")
		session.cisLockdownErr = fmt.Errorf("CIS Service doesnt support private endpoints.")
		session.cisRangeAppErr = fmt.Errorf("CIS Service doesnt support private endpoints.")
		session.cisWAFRuleErr = fmt.Errorf("CIS Service doesnt support private endpoints.")
	}
	cisEndPoint := envFallBack([]string{"IBMCLOUD_CIS_API_ENDPOINT"}, cisURL)

	// IBM Network CIS Zones service
	cisZonesV1Opt := &ciszonesv1.ZonesV1Options{
		URL:           cisEndPoint,
		Crn:           core.StringPtr(""),
		Authenticator: authenticator,
	}
	session.cisZonesV1Client, session.cisZonesErr = ciszonesv1.NewZonesV1(cisZonesV1Opt)
	if session.cisZonesErr != nil {
		session.cisZonesErr = fmt.Errorf(
			"Error occured while configuring CIS Zones service: %s",
			session.cisZonesErr)
	}
	if session.cisZonesV1Client != nil && session.cisZonesV1Client.Service != nil {
		session.cisZonesV1Client.Service.EnableRetries(c.RetryCount, c.RetryDelay)
	}

	// IBM Network CIS DNS Record service
	cisDNSRecordsOpt := &cisdnsrecordsv1.DnsRecordsV1Options{
		URL:            cisEndPoint,
		Crn:            core.StringPtr(""),
		ZoneIdentifier: core.StringPtr(""),
		Authenticator:  authenticator,
	}
	session.cisDNSRecordsClient, session.cisDNSErr = cisdnsrecordsv1.NewDnsRecordsV1(cisDNSRecordsOpt)
	if session.cisDNSErr != nil {
		session.cisDNSErr = fmt.Errorf("Error occured while configuring CIS DNS Service: %s", session.cisDNSErr)
	}
	if session.cisDNSRecordsClient != nil && session.cisDNSRecordsClient.Service != nil {
		session.cisDNSRecordsClient.Service.EnableRetries(c.RetryCount, c.RetryDelay)
	}

	// IBM Network CIS DNS Record bulk service
	cisDNSRecordBulkOpt := &cisdnsbulkv1.DnsRecordBulkV1Options{
		URL:            cisEndPoint,
		Crn:            core.StringPtr(""),
		ZoneIdentifier: core.StringPtr(""),
		Authenticator:  authenticator,
	}
	session.cisDNSRecordBulkClient, session.cisDNSBulkErr = cisdnsbulkv1.NewDnsRecordBulkV1(cisDNSRecordBulkOpt)
	if session.cisDNSBulkErr != nil {
		session.cisDNSBulkErr = fmt.Errorf(
			"Error occured while configuration CIS DNS bulk service : %s",
			session.cisDNSBulkErr)
	}
	if session.cisDNSRecordBulkClient != nil && session.cisDNSRecordBulkClient.Service != nil {
		session.cisDNSRecordBulkClient.Service.EnableRetries(c.RetryCount, c.RetryDelay)
	}

	// IBM Network CIS Global load balancer pool
	cisGLBPoolOpt := &cisglbpoolv0.GlobalLoadBalancerPoolsV0Options{
		URL:           cisEndPoint,
		Crn:           core.StringPtr(""),
		Authenticator: authenticator,
	}
	session.cisGLBPoolClient, session.cisGLBPoolErr =
		cisglbpoolv0.NewGlobalLoadBalancerPoolsV0(cisGLBPoolOpt)
	if session.cisGLBPoolErr != nil {
		session.cisGLBPoolErr =
			fmt.Errorf("Error occured while configuring CIS GLB Pool service: %s",
				session.cisGLBPoolErr)
	}
	if session.cisGLBPoolClient != nil && session.cisGLBPoolClient.Service != nil {
		session.cisGLBPoolClient.Service.EnableRetries(c.RetryCount, c.RetryDelay)
	}

	// IBM Network CIS Global load balancer
	cisGLBOpt := &cisglbv1.GlobalLoadBalancerV1Options{
		URL:            cisEndPoint,
		Authenticator:  authenticator,
		Crn:            core.StringPtr(""),
		ZoneIdentifier: core.StringPtr(""),
	}
	session.cisGLBClient, session.cisGLBErr = cisglbv1.NewGlobalLoadBalancerV1(cisGLBOpt)
	if session.cisGLBErr != nil {
		session.cisGLBErr =
			fmt.Errorf("Error occured while configuring CIS GLB service: %s",
				session.cisGLBErr)
	}
	if session.cisGLBClient != nil && session.cisGLBClient.Service != nil {
		session.cisGLBClient.Service.EnableRetries(c.RetryCount, c.RetryDelay)
	}

	// IBM Network CIS Global load balancer health check/monitor
	cisGLBHealthCheckOpt := &cisglbhealthcheckv1.GlobalLoadBalancerMonitorV1Options{
		URL:           cisEndPoint,
		Crn:           core.StringPtr(""),
		Authenticator: authenticator,
	}
	session.cisGLBHealthCheckClient, session.cisGLBHealthCheckErr =
		cisglbhealthcheckv1.NewGlobalLoadBalancerMonitorV1(cisGLBHealthCheckOpt)
	if session.cisGLBHealthCheckErr != nil {
		session.cisGLBHealthCheckErr =
			fmt.Errorf("Error occured while configuring CIS GLB Health Check service: %s",
				session.cisGLBHealthCheckErr)
	}
	if session.cisGLBHealthCheckClient != nil && session.cisGLBHealthCheckClient.Service != nil {
		session.cisGLBHealthCheckClient.Service.EnableRetries(c.RetryCount, c.RetryDelay)
	}

	// IBM Network CIS IP
	cisIPOpt := &cisipv1.CisIpApiV1Options{
		URL:           cisEndPoint,
		Authenticator: authenticator,
	}
	session.cisIPClient, session.cisIPErr = cisipv1.NewCisIpApiV1(cisIPOpt)
	if session.cisIPErr != nil {
		session.cisIPErr = fmt.Errorf("Error occured while configuring CIS IP service: %s",
			session.cisIPErr)
	}
	if session.cisIPClient != nil && session.cisIPClient.Service != nil {
		session.cisIPClient.Service.EnableRetries(c.RetryCount, c.RetryDelay)
	}

	// IBM Network CIS Zone Rate Limit
	cisRLOpt := &cisratelimitv1.ZoneRateLimitsV1Options{
		URL:            cisEndPoint,
		Crn:            core.StringPtr(""),
		ZoneIdentifier: core.StringPtr(""),
		Authenticator:  authenticator,
	}
	session.cisRLClient, session.cisRLErr = cisratelimitv1.NewZoneRateLimitsV1(cisRLOpt)
	if session.cisRLErr != nil {
		session.cisRLErr = fmt.Errorf(
			"Error occured while cofiguring CIS Zone Rate Limit service: %s",
			session.cisRLErr)
	}
	if session.cisRLClient != nil && session.cisRLClient.Service != nil {
		session.cisRLClient.Service.EnableRetries(c.RetryCount, c.RetryDelay)
	}

	// IBM Network CIS Page Rules
	cisPageRuleOpt := &cispagerulev1.PageRuleApiV1Options{
		URL:           cisEndPoint,
		Crn:           core.StringPtr(""),
		ZoneID:        core.StringPtr(""),
		Authenticator: authenticator,
	}
	session.cisPageRuleClient, session.cisPageRuleErr = cispagerulev1.NewPageRuleApiV1(cisPageRuleOpt)
	if session.cisPageRuleErr != nil {
		session.cisPageRuleErr = fmt.Errorf(
			"Error occured while cofiguring CIS Page Rule service: %s",
			session.cisPageRuleErr)
	}
	if session.cisPageRuleClient != nil && session.cisPageRuleClient.Service != nil {
		session.cisPageRuleClient.Service.EnableRetries(c.RetryCount, c.RetryDelay)
	}

	// IBM Network CIS Edge Function
	cisEdgeFunctionOpt := &cisedgefunctionv1.EdgeFunctionsApiV1Options{
		URL:            cisEndPoint,
		Crn:            core.StringPtr(""),
		ZoneIdentifier: core.StringPtr(""),
		Authenticator:  authenticator,
	}
	session.cisEdgeFunctionClient, session.cisEdgeFunctionErr =
		cisedgefunctionv1.NewEdgeFunctionsApiV1(cisEdgeFunctionOpt)
	if session.cisEdgeFunctionErr != nil {
		session.cisEdgeFunctionErr =
			fmt.Errorf("Error occured while configuring CIS Edge Function service: %s",
				session.cisEdgeFunctionErr)
	}
	if session.cisEdgeFunctionClient != nil && session.cisEdgeFunctionClient.Service != nil {
		session.cisEdgeFunctionClient.Service.EnableRetries(c.RetryCount, c.RetryDelay)
	}

	// IBM Network CIS SSL certificate
	cisSSLOpt := &cissslv1.SslCertificateApiV1Options{
		URL:            cisEndPoint,
		Crn:            core.StringPtr(""),
		ZoneIdentifier: core.StringPtr(""),
		Authenticator:  authenticator,
	}

	session.cisSSLClient, session.cisSSLErr = cissslv1.NewSslCertificateApiV1(cisSSLOpt)
	if session.cisSSLErr != nil {
		session.cisSSLErr =
			fmt.Errorf("Error occured while configuring CIS SSL certificate service: %s",
				session.cisSSLErr)
	}
	if session.cisSSLClient != nil && session.cisSSLClient.Service != nil {
		session.cisSSLClient.Service.EnableRetries(c.RetryCount, c.RetryDelay)
	}

	// IBM Network CIS WAF Package
	cisWAFPackageOpt := &ciswafpackagev1.WafRulePackagesApiV1Options{
		URL:           cisEndPoint,
		Crn:           core.StringPtr(""),
		ZoneID:        core.StringPtr(""),
		Authenticator: authenticator,
	}
	session.cisWAFPackageClient, session.cisWAFPackageErr =
		ciswafpackagev1.NewWafRulePackagesApiV1(cisWAFPackageOpt)
	if session.cisWAFPackageErr != nil {
		session.cisWAFPackageErr =
			fmt.Errorf("Error occured while configuration CIS WAF Package service: %s",
				session.cisWAFPackageErr)
	}
	if session.cisWAFPackageClient != nil && session.cisWAFPackageClient.Service != nil {
		session.cisWAFPackageClient.Service.EnableRetries(c.RetryCount, c.RetryDelay)
	}

	// IBM Network CIS Domain settings
	cisDomainSettingsOpt := &cisdomainsettingsv1.ZonesSettingsV1Options{
		URL:            cisEndPoint,
		Crn:            core.StringPtr(""),
		ZoneIdentifier: core.StringPtr(""),
		Authenticator:  authenticator,
	}
	session.cisDomainSettingsClient, session.cisDomainSettingsErr =
		cisdomainsettingsv1.NewZonesSettingsV1(cisDomainSettingsOpt)
	if session.cisDomainSettingsErr != nil {
		session.cisDomainSettingsErr =
			fmt.Errorf("Error occured while configuring CIS Domain Settings service: %s",
				session.cisDomainSettingsErr)
	}
	if session.cisDomainSettingsClient != nil && session.cisDomainSettingsClient.Service != nil {
		session.cisDomainSettingsClient.Service.EnableRetries(c.RetryCount, c.RetryDelay)
	}

	// IBM Network CIS Routing
	cisRoutingOpt := &cisroutingv1.RoutingV1Options{
		URL:            cisEndPoint,
		Crn:            core.StringPtr(""),
		ZoneIdentifier: core.StringPtr(""),
		Authenticator:  authenticator,
	}
	session.cisRoutingClient, session.cisRoutingErr =
		cisroutingv1.NewRoutingV1(cisRoutingOpt)
	if session.cisRoutingErr != nil {
		session.cisRoutingErr =
			fmt.Errorf("Error occured while configuring CIS Routing service: %s",
				session.cisRoutingErr)
	}
	if session.cisRoutingClient != nil && session.cisRoutingClient.Service != nil {
		session.cisRoutingClient.Service.EnableRetries(c.RetryCount, c.RetryDelay)
	}

	// IBM Network CIS WAF Group
	cisWAFGroupOpt := &ciswafgroupv1.WafRuleGroupsApiV1Options{
		URL:           cisEndPoint,
		Crn:           core.StringPtr(""),
		ZoneID:        core.StringPtr(""),
		Authenticator: authenticator,
	}
	session.cisWAFGroupClient, session.cisWAFGroupErr =
		ciswafgroupv1.NewWafRuleGroupsApiV1(cisWAFGroupOpt)
	if session.cisWAFGroupErr != nil {
		session.cisWAFGroupErr =
			fmt.Errorf("Error occured while configuring CIS WAF Group service: %s",
				session.cisWAFGroupErr)
	}
	if session.cisWAFGroupClient != nil && session.cisWAFGroupClient.Service != nil {
		session.cisWAFGroupClient.Service.EnableRetries(c.RetryCount, c.RetryDelay)
	}

	// IBM Network CIS Cache service
	cisCacheOpt := &ciscachev1.CachingApiV1Options{
		URL:           cisEndPoint,
		Crn:           core.StringPtr(""),
		ZoneID:        core.StringPtr(""),
		Authenticator: authenticator,
	}
	session.cisCacheClient, session.cisCacheErr =
		ciscachev1.NewCachingApiV1(cisCacheOpt)
	if session.cisCacheErr != nil {
		session.cisCacheErr =
			fmt.Errorf("Error occured while configuring CIS Caching service: %s",
				session.cisCacheErr)
	}
	if session.cisCacheClient != nil && session.cisCacheClient.Service != nil {
		session.cisCacheClient.Service.EnableRetries(c.RetryCount, c.RetryDelay)
	}

	// IBM Network CIS Custom pages service
	cisCustomPageOpt := &ciscustompagev1.CustomPagesV1Options{
		URL:            cisEndPoint,
		Crn:            core.StringPtr(""),
		ZoneIdentifier: core.StringPtr(""),
		Authenticator:  authenticator,
	}

	session.cisCustomPageClient, session.cisCustomPageErr =
		ciscustompagev1.NewCustomPagesV1(cisCustomPageOpt)
	if session.cisCustomPageErr != nil {
		session.cisCustomPageErr =
			fmt.Errorf("Error occured while configuring CIS Custom Pages service: %s",
				session.cisCustomPageErr)
	}
	if session.cisCustomPageClient != nil && session.cisCustomPageClient.Service != nil {
		session.cisCustomPageClient.Service.EnableRetries(c.RetryCount, c.RetryDelay)
	}

	// IBM Network CIS Firewall Access rule
	cisAccessRuleOpt := &cisaccessrulev1.ZoneFirewallAccessRulesV1Options{
		URL:            cisEndPoint,
		Crn:            core.StringPtr(""),
		ZoneIdentifier: core.StringPtr(""),
		Authenticator:  authenticator,
	}
	session.cisAccessRuleClient, session.cisAccessRuleErr =
		cisaccessrulev1.NewZoneFirewallAccessRulesV1(cisAccessRuleOpt)
	if session.cisAccessRuleErr != nil {
		session.cisAccessRuleErr =
			fmt.Errorf("Error occured while configuring CIS Firewall Access Rule service: %s",
				session.cisAccessRuleErr)
	}
	if session.cisAccessRuleClient != nil && session.cisAccessRuleClient.Service != nil {
		session.cisAccessRuleClient.Service.EnableRetries(c.RetryCount, c.RetryDelay)
	}

	// IBM Network CIS Firewall User Agent Blocking rule
	cisUARuleOpt := &cisuarulev1.UserAgentBlockingRulesV1Options{
		URL:            cisEndPoint,
		Crn:            core.StringPtr(""),
		ZoneIdentifier: core.StringPtr(""),
		Authenticator:  authenticator,
	}
	session.cisUARuleClient, session.cisUARuleErr =
		cisuarulev1.NewUserAgentBlockingRulesV1(cisUARuleOpt)
	if session.cisUARuleErr != nil {
		session.cisUARuleErr =
			fmt.Errorf("Error occured while configuring CIS Firewall User Agent Blocking Rule service: %s",
				session.cisUARuleErr)
	}
	if session.cisUARuleClient != nil && session.cisUARuleClient.Service != nil {
		session.cisUARuleClient.Service.EnableRetries(c.RetryCount, c.RetryDelay)
	}

	// IBM Network CIS Firewall Lockdown rule
	cisLockdownOpt := &cislockdownv1.ZoneLockdownV1Options{
		URL:            cisEndPoint,
		Crn:            core.StringPtr(""),
		ZoneIdentifier: core.StringPtr(""),
		Authenticator:  authenticator,
	}
	session.cisLockdownClient, session.cisLockdownErr =
		cislockdownv1.NewZoneLockdownV1(cisLockdownOpt)
	if session.cisLockdownErr != nil {
		session.cisLockdownErr =
			fmt.Errorf("Error occured while configuring CIS Firewall Lockdown Rule service: %s",
				session.cisLockdownErr)
	}
	if session.cisLockdownClient != nil && session.cisLockdownClient.Service != nil {
		session.cisLockdownClient.Service.EnableRetries(c.RetryCount, c.RetryDelay)
	}

	// IBM Network CIS Range Application rule
	cisRangeAppOpt := &cisrangeappv1.RangeApplicationsV1Options{
		URL:            cisEndPoint,
		Crn:            core.StringPtr(""),
		ZoneIdentifier: core.StringPtr(""),
		Authenticator:  authenticator,
	}
	session.cisRangeAppClient, session.cisRangeAppErr =
		cisrangeappv1.NewRangeApplicationsV1(cisRangeAppOpt)
	if session.cisRangeAppErr != nil {
		session.cisRangeAppErr =
			fmt.Errorf("Error occured while configuring CIS Range Application rule service: %s",
				session.cisRangeAppErr)
	}
	if session.cisRangeAppClient != nil && session.cisRangeAppClient.Service != nil {
		session.cisRangeAppClient.Service.EnableRetries(c.RetryCount, c.RetryDelay)
	}

	// IBM Network CIS WAF Rule Service
	cisWAFRuleOpt := &ciswafrulev1.WafRulesApiV1Options{
		URL:           cisEndPoint,
		Crn:           core.StringPtr(""),
		ZoneID:        core.StringPtr(""),
		Authenticator: authenticator,
	}
	session.cisWAFRuleClient, session.cisWAFRuleErr =
		ciswafrulev1.NewWafRulesApiV1(cisWAFRuleOpt)
	if session.cisWAFRuleErr != nil {
		session.cisWAFRuleErr = fmt.Errorf(
			"Error occured while configuring CIS WAF Rules service: %s",
			session.cisWAFRuleErr)
	}
	if session.cisWAFRuleClient != nil && session.cisWAFRuleClient.Service != nil {
		session.cisWAFRuleClient.Service.EnableRetries(c.RetryCount, c.RetryDelay)
	}

	// iamIdenityURL := fmt.Sprintf("https://%s.iam.cloud.ibm.com/v1", c.Region)
	iamURL := iamidentity.DefaultServiceURL
	if c.Visibility == "private" || c.Visibility == "public-and-private" {
		if c.Region == "us-south" || c.Region == "us-east" {
			iamURL = contructEndpoint(fmt.Sprintf("private.%s.iam", c.Region), cloudEndpoint)
		} else {
			iamURL = contructEndpoint("private.iam", cloudEndpoint)
		}
	}
	iamIdentityOptions := &iamidentity.IamIdentityV1Options{
		Authenticator: authenticator,
		URL:           envFallBack([]string{"IBMCLOUD_IAM_API_ENDPOINT"}, iamURL),
	}
	iamIdentityClient, err := iamidentity.NewIamIdentityV1(iamIdentityOptions)
	if err != nil {
		session.iamIdentityErr = fmt.Errorf("Error occured while configuring IAM Identity service: %q", err)
	}
	if iamIdentityClient != nil && iamIdentityClient.Service != nil {
		iamIdentityClient.Service.EnableRetries(c.RetryCount, c.RetryDelay)
	}
	session.iamIdentityAPI = iamIdentityClient

	iamPolicyManagementURL := iampolicymanagement.DefaultServiceURL
	if c.Visibility == "private" || c.Visibility == "public-and-private" {
		if c.Region == "us-south" || c.Region == "us-east" {
			iamPolicyManagementURL = contructEndpoint(fmt.Sprintf("private.%s.iam", c.Region), cloudEndpoint)
		} else {
			iamPolicyManagementURL = contructEndpoint("private.iam", cloudEndpoint)
		}
	}
	iamPolicyManagementOptions := &iampolicymanagement.IamPolicyManagementV1Options{
		Authenticator: authenticator,
		URL:           envFallBack([]string{"IBMCLOUD_IAM_API_ENDPOINT"}, iamPolicyManagementURL),
	}
	iamPolicyManagementClient, err := iampolicymanagement.NewIamPolicyManagementV1(iamPolicyManagementOptions)
	if err != nil {
		session.iamPolicyManagementErr = fmt.Errorf("Error occured while configuring IAM Policy Management service: %q", err)
	}
	if iamPolicyManagementClient != nil && iamPolicyManagementClient.Service != nil {
		iamPolicyManagementClient.Service.EnableRetries(c.RetryCount, c.RetryDelay)
	}
	session.iamPolicyManagementAPI = iamPolicyManagementClient

	rmURL := resourcemanager.DefaultServiceURL
	if c.Visibility == "private" {
		if c.Region == "us-south" || c.Region == "us-east" {
			rmURL = contructEndpoint(fmt.Sprintf("private.%s.resource-controller", c.Region), fmt.Sprintf("%s/v2", cloudEndpoint))
		} else {
			fmt.Println("Private Endpint supports only us-south and us-east region specific endpoint")
			rmURL = contructEndpoint("private.us-south.resource-controller", fmt.Sprintf("%s/v2", cloudEndpoint))
		}
	}
	if c.Visibility == "public-and-private" {
		if c.Region == "us-south" || c.Region == "us-east" {
			rmURL = contructEndpoint(fmt.Sprintf("private.%s.resource-controller", c.Region), fmt.Sprintf("%s/v2", cloudEndpoint))
		} else {
			rmURL = resourcemanager.DefaultServiceURL
		}
	}
	resourceManagerOptions := &resourcemanager.ResourceManagerV2Options{
		Authenticator: authenticator,
		URL:           envFallBack([]string{"IBMCLOUD_RESOURCE_MANAGEMENT_API_ENDPOINT"}, rmURL),
	}
	resourceManagerClient, err := resourcemanager.NewResourceManagerV2(resourceManagerOptions)
	if err != nil {
		session.resourceManagerErr = fmt.Errorf("Error occured while configuring Resource Manager service: %q", err)
	}
	if resourceManagerClient != nil {
		resourceManagerClient.EnableRetries(c.RetryCount, c.RetryDelay)
	}
	session.resourceManagerAPI = resourceManagerClient

	enterpriseURL := enterprisemanagementv1.DefaultServiceURL
	if c.Visibility == "private" {
		if c.Region == "us-south" || c.Region == "us-east" || c.Region == "eu-fr" {
			enterpriseURL = contructEndpoint(fmt.Sprintf("private.%s.enterprise", c.Region), fmt.Sprintf("%s/v1", cloudEndpoint))
		} else {
			fmt.Println("Private Endpint supports only us-south and us-east region specific endpoint")
			enterpriseURL = contructEndpoint("private.us-south.enterprise", fmt.Sprintf("%s/v1", cloudEndpoint))
		}
	}
	if c.Visibility == "public-and-private" {
		if c.Region == "us-south" || c.Region == "us-east" || c.Region == "eu-fr" {
			enterpriseURL = contructEndpoint(fmt.Sprintf("private.%s.enterprise", c.Region),
				fmt.Sprintf("%s/v1", cloudEndpoint))
		} else {
			enterpriseURL = enterprisemanagementv1.DefaultServiceURL
		}
	}
	enterpriseManagementClientOptions := &enterprisemanagementv1.EnterpriseManagementV1Options{
		Authenticator: authenticator,
		URL:           envFallBack([]string{"IBMCLOUD_ENTERPRISE_API_ENDPOINT"}, enterpriseURL),
	}
	enterpriseManagementClient, err := enterprisemanagementv1.NewEnterpriseManagementV1(enterpriseManagementClientOptions)
	if err == nil {
		enterpriseManagementClient.EnableRetries(c.RetryCount, c.RetryDelay)
	} else {
		session.enterpriseManagementClientErr = fmt.Errorf("Error occurred while configuring IBM Cloud Enterprise Management API service: %q", err)
	}
	session.enterpriseManagementClient = enterpriseManagementClient

	// resource controller API
	rcURL := resourcecontroller.DefaultServiceURL
	if c.Visibility == "private" {
		if c.Region == "us-south" || c.Region == "us-east" {
			rcURL = contructEndpoint(fmt.Sprintf("private.%s.resource-controller", c.Region), cloudEndpoint)
		} else {
			fmt.Println("Private Endpint supports only us-south and us-east region specific endpoint")
			rcURL = contructEndpoint("private.us-south.resource-controller", cloudEndpoint)
		}
	}
	if c.Visibility == "public-and-private" {
		if c.Region == "us-south" || c.Region == "us-east" {
			rcURL = contructEndpoint(fmt.Sprintf("private.%s.resource-controller", c.Region), cloudEndpoint)
		} else {
			rcURL = resourcecontroller.DefaultServiceURL
		}
	}
	resourceControllerOptions := &resourcecontroller.ResourceControllerV2Options{
		Authenticator: authenticator,
		URL:           envFallBack([]string{"IBMCLOUD_RESOURCE_CONTROLLER_API_ENDPOINT"}, rcURL),
	}
	resourceControllerClient, err := resourcecontroller.NewResourceControllerV2(resourceControllerOptions)
	if err != nil {
		session.resourceControllerErr = fmt.Errorf("Error occured while configuring Resource Controller service: %q", err)
	}
	if resourceControllerClient != nil {
		resourceControllerClient.EnableRetries(c.RetryCount, c.RetryDelay)
	}
	session.resourceControllerAPI = resourceControllerClient
	// var authenticator2 *core.BearerTokenAuthenticator
	// Construct an "options" struct for creating the service client.
	secretsManagerClientOptions := &secretsmanagerv1.SecretsManagerV1Options{
		Authenticator: authenticator,
	}

	/// Construct the service client.
	session.secretsManagerClient, err = secretsmanagerv1.NewSecretsManagerV1(secretsManagerClientOptions)
	if err == nil {
		// Enable retries for API calls
		session.secretsManagerClient.Service.EnableRetries(c.RetryCount, c.RetryDelay)
		// Add custom header for analytics
		session.secretsManagerClient.SetDefaultHeaders(gohttp.Header{
			"X-Original-User-Agent": {fmt.Sprintf("terraform-provider-ibm/%s", version.Version)},
		})
	} else {
		session.secretsManagerClientErr = fmt.Errorf("Error occurred while configuring IBM Cloud Secrets Manager API service: %q", err)
	}

	return session, nil
}

// CreateVersionDate requires mandatory version attribute. Any date from 2019-12-13 up to the currentdate may be provided. Specify the current date to request the latest version.
func CreateVersionDate() *string {
	version := time.Now().Format("2006-01-02")
	return &version
}

func newSession(c *Config) (*Session, error) {
	ibmSession := &Session{}

	softlayerSession := &slsession.Session{
		Endpoint:  c.SoftLayerEndpointURL,
		Timeout:   c.SoftLayerTimeout,
		UserName:  c.SoftLayerUserName,
		APIKey:    c.SoftLayerAPIKey,
		Debug:     os.Getenv("TF_LOG") != "",
		Retries:   c.RetryCount,
		RetryWait: c.RetryDelay,
	}

	if c.IAMToken != "" {
		log.Println("Configuring SoftLayer Session with token")
		softlayerSession.IAMToken = c.IAMToken
	}
	if c.SoftLayerAPIKey != "" && c.SoftLayerUserName != "" {
		log.Println("Configuring SoftLayer Session with API key")
		softlayerSession.APIKey = c.SoftLayerAPIKey
		softlayerSession.UserName = c.SoftLayerUserName
	}
	softlayerSession.AppendUserAgent(fmt.Sprintf("terraform-provider-ibm/%s", version.Version))
	ibmSession.SoftLayerSession = softlayerSession

	if (c.IAMToken != "" && c.IAMRefreshToken == "") || (c.IAMToken == "" && c.IAMRefreshToken != "") {
		return nil, fmt.Errorf("iam_token and iam_refresh_token must be provided")
	}

	if c.IAMToken != "" && c.IAMRefreshToken != "" {
		log.Println("Configuring IBM Cloud Session with token")
		var sess *bxsession.Session
		bmxConfig := &bluemix.Config{
			IAMAccessToken:  c.IAMToken,
			IAMRefreshToken: c.IAMRefreshToken,
			//Comment out debug mode for v0.12
			//Debug:         os.Getenv("TF_LOG") != "",
			HTTPTimeout:   c.BluemixTimeout,
			Region:        c.Region,
			ResourceGroup: c.ResourceGroup,
			RetryDelay:    &c.RetryDelay,
			MaxRetries:    &c.RetryCount,
			Visibility:    c.Visibility,
		}
		sess, err := bxsession.New(bmxConfig)
		if err != nil {
			return nil, err
		}
		ibmSession.BluemixSession = sess
	}

	if c.BluemixAPIKey != "" {
		log.Println("Configuring IBM Cloud Session with API key")
		var sess *bxsession.Session
		bmxConfig := &bluemix.Config{
			BluemixAPIKey: c.BluemixAPIKey,
			//Comment out debug mode for v0.12
			//Debug:         os.Getenv("TF_LOG") != "",
			HTTPTimeout:   c.BluemixTimeout,
			Region:        c.Region,
			ResourceGroup: c.ResourceGroup,
			RetryDelay:    &c.RetryDelay,
			MaxRetries:    &c.RetryCount,
			Visibility:    c.Visibility,
			//PowerServiceInstance: c.PowerServiceInstance,
		}
		sess, err := bxsession.New(bmxConfig)
		if err != nil {
			return nil, err
		}
		ibmSession.BluemixSession = sess
	}

	return ibmSession, nil
}

func authenticateAPIKey(sess *bxsession.Session) error {
	config := sess.Config
	tokenRefresher, err := authentication.NewIAMAuthRepository(config, &rest.Client{
		DefaultHeader: gohttp.Header{
			"User-Agent": []string{http.UserAgent()},
		},
	})
	if err != nil {
		return err
	}
	return tokenRefresher.AuthenticateAPIKey(config.BluemixAPIKey)
}

func authenticateCF(sess *bxsession.Session) error {
	config := sess.Config
	tokenRefresher, err := authentication.NewUAARepository(config, &rest.Client{
		DefaultHeader: gohttp.Header{
			"User-Agent": []string{http.UserAgent()},
		},
	})
	if err != nil {
		return err
	}
	return tokenRefresher.AuthenticateAPIKey(config.BluemixAPIKey)
}

func fetchUserDetails(sess *bxsession.Session, retries int, retryDelay time.Duration) (*UserConfig, error) {
	config := sess.Config
	user := UserConfig{}
	var bluemixToken string

	if strings.HasPrefix(config.IAMAccessToken, "Bearer") {
		bluemixToken = config.IAMAccessToken[7:len(config.IAMAccessToken)]
	} else {
		bluemixToken = config.IAMAccessToken
	}

	token, err := jwt.Parse(bluemixToken, func(token *jwt.Token) (interface{}, error) {
		return "", nil
	})
	//TODO validate with key
	if err != nil && !strings.Contains(err.Error(), "key is of invalid type") {
		if retries > 0 {
			if config.BluemixAPIKey != "" {
				time.Sleep(retryDelay)
				log.Printf("Retrying authentication for user details %d", retries)
				_ = authenticateAPIKey(sess)
				return fetchUserDetails(sess, retries-1, retryDelay)
			}
		}
		return &user, err
	}
	claims := token.Claims.(jwt.MapClaims)
	if email, ok := claims["email"]; ok {
		user.userEmail = email.(string)
	}
	user.userID = claims["id"].(string)
	user.userAccount = claims["account"].(map[string]interface{})["bss"].(string)
	iss := claims["iss"].(string)
	if strings.Contains(iss, "https://iam.cloud.ibm.com") {
		user.cloudName = "bluemix"
	} else {
		user.cloudName = "staging"
	}
	user.cloudType = "public"

	user.generation = 2
	return &user, nil
}

func refreshToken(sess *bxsession.Session) error {
	config := sess.Config
	tokenRefresher, err := authentication.NewIAMAuthRepository(config, &rest.Client{
		DefaultHeader: gohttp.Header{
			"User-Agent": []string{http.UserAgent()},
		},
	})
	if err != nil {
		return err
	}
	_, err = tokenRefresher.RefreshToken()
	return err
}

func envFallBack(envs []string, defaultValue string) string {
	for _, k := range envs {
		if v := os.Getenv(k); v != "" {
			return v
		}
	}
	return defaultValue
}

// DefaultTransport ...
func DefaultTransport() gohttp.RoundTripper {
	transport := &gohttp.Transport{
		Proxy:               gohttp.ProxyFromEnvironment,
		DisableKeepAlives:   true,
		MaxIdleConnsPerHost: -1,
		TLSClientConfig: &tls.Config{
			InsecureSkipVerify: false,
		},
	}
	return transport
}

func isRetryable(err error) bool {
	if bmErr, ok := err.(bmxerror.RequestFailure); ok {
		switch bmErr.StatusCode() {
		case 408, 504, 599, 429, 500, 502, 520, 503:
			return true
		}
	}

	if netErr, ok := err.(net.Error); ok && netErr.Timeout() {
		return true
	}

	if netErr, ok := err.(*net.OpError); ok && netErr.Timeout() {
		return true
	}

	if netErr, ok := err.(net.UnknownNetworkError); ok && netErr.Timeout() {
		return true
	}

	return false
}

func contructEndpoint(subdomain, domain string) string {
	endpoint := fmt.Sprintf("https://%s.%s", subdomain, domain)
	return endpoint
}
