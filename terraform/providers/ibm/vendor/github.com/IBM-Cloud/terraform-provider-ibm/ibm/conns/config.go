// Copyright IBM Corp. 2017, 2021 All Rights Reserved.
// Licensed under the Mozilla Public License v2.0

package conns

import (
	"crypto/tls"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net"
	gohttp "net/http"
	"os"
	"strings"
	"time"

	// Added code for the Power Colo Offering

	"github.com/IBM-Cloud/container-services-go-sdk/kubernetesserviceapiv1"
	"github.com/IBM-Cloud/container-services-go-sdk/satellitelinkv1"
	apigateway "github.com/IBM/apigateway-go-sdk/apigatewaycontrollerapiv1"
	"github.com/IBM/appconfiguration-go-admin-sdk/appconfigurationv1"
	appid "github.com/IBM/appid-management-go-sdk/appidmanagementv4"
	"github.com/IBM/cloud-databases-go-sdk/clouddatabasesv5"
	"github.com/IBM/container-registry-go-sdk/containerregistryv1"
	"github.com/IBM/go-sdk-core/v5/core"
	cosconfig "github.com/IBM/ibm-cos-sdk-go-config/resourceconfigurationv1"
	kp "github.com/IBM/keyprotect-go-client"
	cisalertsv1 "github.com/IBM/networking-go-sdk/alertsv1"
	cisoriginpull "github.com/IBM/networking-go-sdk/authenticatedoriginpullapiv1"
	ciscachev1 "github.com/IBM/networking-go-sdk/cachingapiv1"
	cisipv1 "github.com/IBM/networking-go-sdk/cisipapiv1"
	ciscustompagev1 "github.com/IBM/networking-go-sdk/custompagesv1"
	dlProviderV2 "github.com/IBM/networking-go-sdk/directlinkproviderv2"
	dl "github.com/IBM/networking-go-sdk/directlinkv1"
	cisdnsbulkv1 "github.com/IBM/networking-go-sdk/dnsrecordbulkv1"
	cisdnsrecordsv1 "github.com/IBM/networking-go-sdk/dnsrecordsv1"
	dns "github.com/IBM/networking-go-sdk/dnssvcsv1"
	cisedgefunctionv1 "github.com/IBM/networking-go-sdk/edgefunctionsapiv1"
	cisfiltersv1 "github.com/IBM/networking-go-sdk/filtersv1"
	cisfirewallrulesv1 "github.com/IBM/networking-go-sdk/firewallrulesv1"
	cisglbhealthcheckv1 "github.com/IBM/networking-go-sdk/globalloadbalancermonitorv1"
	cisglbpoolv0 "github.com/IBM/networking-go-sdk/globalloadbalancerpoolsv0"
	cisglbv1 "github.com/IBM/networking-go-sdk/globalloadbalancerv1"
	cislogpushjobsapiv1 "github.com/IBM/networking-go-sdk/logpushjobsapiv1"
	cismtlsv1 "github.com/IBM/networking-go-sdk/mtlsv1"
	cispagerulev1 "github.com/IBM/networking-go-sdk/pageruleapiv1"
	cisrangeappv1 "github.com/IBM/networking-go-sdk/rangeapplicationsv1"
	cisroutingv1 "github.com/IBM/networking-go-sdk/routingv1"
	cissslv1 "github.com/IBM/networking-go-sdk/sslcertificateapiv1"
	tg "github.com/IBM/networking-go-sdk/transitgatewayapisv1"
	cisuarulev1 "github.com/IBM/networking-go-sdk/useragentblockingrulesv1"
	ciswafgroupv1 "github.com/IBM/networking-go-sdk/wafrulegroupsapiv1"
	ciswafpackagev1 "github.com/IBM/networking-go-sdk/wafrulepackagesapiv1"
	ciswafrulev1 "github.com/IBM/networking-go-sdk/wafrulesapiv1"
	ciswebhooksv1 "github.com/IBM/networking-go-sdk/webhooksv1"
	cisaccessrulev1 "github.com/IBM/networking-go-sdk/zonefirewallaccessrulesv1"
	cislockdownv1 "github.com/IBM/networking-go-sdk/zonelockdownv1"
	cisratelimitv1 "github.com/IBM/networking-go-sdk/zoneratelimitsv1"
	cisdomainsettingsv1 "github.com/IBM/networking-go-sdk/zonessettingsv1"
	ciszonesv1 "github.com/IBM/networking-go-sdk/zonesv1"
	"github.com/IBM/platform-services-go-sdk/atrackerv1"
	"github.com/IBM/platform-services-go-sdk/atrackerv2"
	"github.com/IBM/platform-services-go-sdk/catalogmanagementv1"
	"github.com/IBM/platform-services-go-sdk/contextbasedrestrictionsv1"
	"github.com/IBM/platform-services-go-sdk/enterprisemanagementv1"
	"github.com/IBM/platform-services-go-sdk/globaltaggingv1"
	iamaccessgroups "github.com/IBM/platform-services-go-sdk/iamaccessgroupsv2"
	iamidentity "github.com/IBM/platform-services-go-sdk/iamidentityv1"
	iampolicymanagement "github.com/IBM/platform-services-go-sdk/iampolicymanagementv1"
	ibmcloudshellv1 "github.com/IBM/platform-services-go-sdk/ibmcloudshellv1"
	resourcecontroller "github.com/IBM/platform-services-go-sdk/resourcecontrollerv2"
	resourcemanager "github.com/IBM/platform-services-go-sdk/resourcemanagerv2"
	"github.com/IBM/push-notifications-go-sdk/pushservicev1"
	"github.com/IBM/scc-go-sdk/v3/adminserviceapiv1"
	"github.com/IBM/scc-go-sdk/v3/configurationgovernancev1"
	"github.com/IBM/scc-go-sdk/v3/posturemanagementv2"
	schematicsv1 "github.com/IBM/schematics-go-sdk/schematicsv1"
	"github.com/IBM/secrets-manager-go-sdk/secretsmanagerv1"
	vpc "github.com/IBM/vpc-go-sdk/vpcv1"
	"github.com/apache/openwhisk-client-go/whisk"
	jwt "github.com/golang-jwt/jwt"
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
	"github.com/IBM/continuous-delivery-go-sdk/cdtektonpipelinev2"
	"github.com/IBM/continuous-delivery-go-sdk/cdtoolchainv2"
	"github.com/IBM/event-notifications-go-admin-sdk/eventnotificationsv1"
	"github.com/IBM/eventstreams-go-sdk/pkg/schemaregistryv1"
	"github.com/IBM/ibm-hpcs-uko-sdk/ukov4"
	"github.com/IBM/scc-go-sdk/v3/posturemanagementv1"
)

// RetryAPIDelay - retry api delay
const RetryAPIDelay = 5 * time.Second

//BluemixRegion ...
var BluemixRegion string

var (
	errEmptyBluemixCredentials = errors.New("ibmcloud_api_key or bluemix_api_key or iam_token and iam_refresh_token must be provided. Please see the documentation on how to configure it")
)

//UserConfig ...
type UserConfig struct {
	UserID      string
	UserEmail   string
	UserAccount string
	CloudName   string `default:"bluemix"`
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

	//TrustedProfileToken Token
	IAMTrustedProfileID string

	//IAM Refresh Token
	IAMRefreshToken string

	// Zone
	Zone          string
	Visibility    string
	EndpointsFile string
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
	AppIDAPI() (*appid.AppIDManagementV4, error)
	BluemixSession() (*bxsession.Session, error)
	BluemixAcccountAPI() (accountv2.AccountServiceAPI, error)
	BluemixAcccountv1API() (accountv1.AccountServiceAPI, error)
	BluemixUserDetails() (*UserConfig, error)
	ContainerAPI() (containerv1.ContainerServiceAPI, error)
	VpcContainerAPI() (containerv2.ContainerServiceAPI, error)
	ContainerRegistryV1() (*containerregistryv1.ContainerRegistryV1, error)
	FunctionClient() (*whisk.Client, error)
	GlobalSearchAPI() (globalsearchv2.GlobalSearchServiceAPI, error)
	GlobalTaggingAPI() (globaltaggingv3.GlobalTaggingServiceAPI, error)
	GlobalTaggingAPIv1() (globaltaggingv1.GlobalTaggingV1, error)
	ICDAPI() (icdv4.ICDServiceAPI, error)
	CloudDatabasesV5() (*clouddatabasesv5.CloudDatabasesV5, error)
	IAMPolicyManagementV1API() (*iampolicymanagement.IamPolicyManagementV1, error)
	IAMAccessGroupsV2() (*iamaccessgroups.IamAccessGroupsV2, error)
	MccpAPI() (mccpv2.MccpServiceAPI, error)
	ResourceCatalogAPI() (catalog.ResourceCatalogAPI, error)
	ResourceManagementAPIv2() (managementv2.ResourceManagementAPIv2, error)
	ResourceControllerAPI() (controller.ResourceControllerAPI, error)
	ResourceControllerAPIV2() (controllerv2.ResourceControllerAPIV2, error)
	SoftLayerSession() *slsession.Session
	IBMPISession() (*ibmpisession.IBMPISession, error)
	UserManagementAPI() (usermanagementv2.UserManagementAPI, error)
	PushServiceV1() (*pushservicev1.PushServiceV1, error)
	EventNotificationsApiV1() (*eventnotificationsv1.EventNotificationsV1, error)
	AppConfigurationV1() (*appconfigurationv1.AppConfigurationV1, error)
	CertificateManagerAPI() (certificatemanager.CertificateManagerServiceAPI, error)
	KeyProtectAPI() (*kp.Client, error)
	KeyManagementAPI() (*kp.Client, error)
	VpcV1API() (*vpc.VpcV1, error)
	APIGateway() (*apigateway.ApiGatewayControllerApiV1, error)
	PrivateDNSClientSession() (*dns.DnsSvcsV1, error)
	CosConfigV1API() (*cosconfig.ResourceConfigurationV1, error)
	DirectlinkV1API() (*dl.DirectLinkV1, error)
	DirectlinkProviderV2API() (*dlProviderV2.DirectLinkProviderV2, error)
	TransitGatewayV1API() (*tg.TransitGatewayApisV1, error)
	HpcsEndpointAPI() (hpcs.HPCSV2, error)
	UkoV4() (*ukov4.UkoV4, error)
	FunctionIAMNamespaceAPI() (functions.FunctionServiceAPI, error)
	CisZonesV1ClientSession() (*ciszonesv1.ZonesV1, error)
	CisAlertsSession() (*cisalertsv1.AlertsV1, error)
	CisOrigAuthSession() (*cisoriginpull.AuthenticatedOriginPullApiV1, error)
	CisDNSRecordClientSession() (*cisdnsrecordsv1.DnsRecordsV1, error)
	CisDNSRecordBulkClientSession() (*cisdnsbulkv1.DnsRecordBulkV1, error)
	CisGLBClientSession() (*cisglbv1.GlobalLoadBalancerV1, error)
	CisGLBPoolClientSession() (*cisglbpoolv0.GlobalLoadBalancerPoolsV0, error)
	CisGLBHealthCheckClientSession() (*cisglbhealthcheckv1.GlobalLoadBalancerMonitorV1, error)
	CisIPClientSession() (*cisipv1.CisIpApiV1, error)
	CisPageRuleClientSession() (*cispagerulev1.PageRuleApiV1, error)
	CisLogpushJobsSession() (*cislogpushjobsapiv1.LogpushJobsApiV1, error)
	CisRLClientSession() (*cisratelimitv1.ZoneRateLimitsV1, error)
	CisEdgeFunctionClientSession() (*cisedgefunctionv1.EdgeFunctionsApiV1, error)
	CisSSLClientSession() (*cissslv1.SslCertificateApiV1, error)
	CisWAFPackageClientSession() (*ciswafpackagev1.WafRulePackagesApiV1, error)
	CisDomainSettingsClientSession() (*cisdomainsettingsv1.ZonesSettingsV1, error)
	CisRoutingClientSession() (*cisroutingv1.RoutingV1, error)
	CisWAFGroupClientSession() (*ciswafgroupv1.WafRuleGroupsApiV1, error)
	CisCacheClientSession() (*ciscachev1.CachingApiV1, error)
	CisMtlsSession() (*cismtlsv1.MtlsV1, error)
	CisWebhookSession() (*ciswebhooksv1.WebhooksV1, error)
	CisCustomPageClientSession() (*ciscustompagev1.CustomPagesV1, error)
	CisAccessRuleClientSession() (*cisaccessrulev1.ZoneFirewallAccessRulesV1, error)
	CisUARuleClientSession() (*cisuarulev1.UserAgentBlockingRulesV1, error)
	CisLockdownClientSession() (*cislockdownv1.ZoneLockdownV1, error)
	CisRangeAppClientSession() (*cisrangeappv1.RangeApplicationsV1, error)
	CisWAFRuleClientSession() (*ciswafrulev1.WafRulesApiV1, error)
	IAMIdentityV1API() (*iamidentity.IamIdentityV1, error)
	IBMCloudShellV1() (*ibmcloudshellv1.IBMCloudShellV1, error)
	ResourceManagerV2API() (*resourcemanager.ResourceManagerV2, error)
	CatalogManagementV1() (*catalogmanagementv1.CatalogManagementV1, error)
	EnterpriseManagementV1() (*enterprisemanagementv1.EnterpriseManagementV1, error)
	ResourceControllerV2API() (*resourcecontroller.ResourceControllerV2, error)
	SecretsManagerV1() (*secretsmanagerv1.SecretsManagerV1, error)
	SchematicsV1() (*schematicsv1.SchematicsV1, error)
	SatelliteClientSession() (*kubernetesserviceapiv1.KubernetesServiceApiV1, error)
	SatellitLinkClientSession() (*satellitelinkv1.SatelliteLinkV1, error)
	CisFiltersSession() (*cisfiltersv1.FiltersV1, error)
	CisFirewallRulesSession() (*cisfirewallrulesv1.FirewallRulesV1, error)
	AtrackerV1() (*atrackerv1.AtrackerV1, error)
	AtrackerV2() (*atrackerv2.AtrackerV2, error)
	ESschemaRegistrySession() (*schemaregistryv1.SchemaregistryV1, error)
	AdminServiceApiV1() (*adminserviceapiv1.AdminServiceApiV1, error)
	ConfigurationGovernanceV1() (*configurationgovernancev1.ConfigurationGovernanceV1, error)
	PostureManagementV1() (*posturemanagementv1.PostureManagementV1, error)
	ContextBasedRestrictionsV1() (*contextbasedrestrictionsv1.ContextBasedRestrictionsV1, error)
	PostureManagementV2() (*posturemanagementv2.PostureManagementV2, error)
	CdToolchainV2() (*cdtoolchainv2.CdToolchainV2, error)
	CdTektonPipelineV2() (*cdtektonpipelinev2.CdTektonPipelineV2, error)
}

type clientSession struct {
	session *Session

	appidErr error
	appidAPI *appid.AppIDManagementV4

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

	ibmCloudShellClient    *ibmcloudshellv1.IBMCloudShellV1
	ibmCloudShellClientErr error

	userManagementErr error
	userManagementAPI usermanagementv2.UserManagementAPI

	icdConfigErr  error
	icdServiceAPI icdv4.ICDServiceAPI

	cloudDatabasesClientErr error
	cloudDatabasesClient    *clouddatabasesv5.CloudDatabasesV5

	resourceControllerConfigErr  error
	resourceControllerServiceAPI controller.ResourceControllerAPI

	resourceControllerConfigErrv2  error
	resourceControllerServiceAPIv2 controllerv2.ResourceControllerAPIV2

	resourceManagementConfigErrv2  error
	resourceManagementServiceAPIv2 managementv2.ResourceManagementAPIv2

	resourceCatalogConfigErr  error
	resourceCatalogServiceAPI catalog.ResourceCatalogAPI

	ibmpiConfigErr error
	ibmpiSession   *ibmpisession.IBMPISession

	kpErr error
	kpAPI *kp.API

	kmsErr error
	kmsAPI *kp.API

	hpcsEndpointErr error
	hpcsEndpointAPI hpcs.HPCSV2

	ukoClient    *ukov4.UkoV4
	ukoClientErr error

	pDNSClient *dns.DnsSvcsV1
	pDNSErr    error

	bluemixSessionErr error

	pushServiceClient    *pushservicev1.PushServiceV1
	pushServiceClientErr error

	eventNotificationsApiClient    *eventnotificationsv1.EventNotificationsV1
	eventNotificationsApiClientErr error

	appConfigurationClient    *appconfigurationv1.AppConfigurationV1
	appConfigurationClientErr error

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

	// CIS Alerts
	cisAlertsClient *cisalertsv1.AlertsV1
	cisAlertsErr    error

	// CIS Authenticated Origin Pull
	cisOriginAuthClient  *cisoriginpull.AuthenticatedOriginPullApiV1
	cisOriginAuthPullErr error

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

	// CIS LogpushJobs service option
	cisLogpushJobsClient *cislogpushjobsapiv1.LogpushJobsApiV1
	cisLogpushJobsErr    error

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

	//Satellite service
	satelliteClient    *kubernetesserviceapiv1.KubernetesServiceApiV1
	satelliteClientErr error

	//IAM Policy Management
	iamPolicyManagementErr error
	iamPolicyManagementAPI *iampolicymanagement.IamPolicyManagementV1

	//IAM Access Groups
	iamAccessGroupsErr error
	iamAccessGroupsAPI *iamaccessgroups.IamAccessGroupsV2

	// MTLS Session options
	cisMtlsClient *cismtlsv1.MtlsV1
	cisMtlsErr    error

	// CIS Webhooks options
	cisWebhooksClient *ciswebhooksv1.WebhooksV1
	cisWebhooksErr    error

	// CIS Filters options
	cisFiltersClient *cisfiltersv1.FiltersV1
	cisFiltersErr    error

	// CIS FirewallRules options
	cisFirewallRulesClient *cisfirewallrulesv1.FirewallRulesV1
	cisFirewallRulesErr    error

	//Atracker
	atrackerClient    *atrackerv1.AtrackerV1
	atrackerClientErr error

	atrackerClientV2    *atrackerv2.AtrackerV2
	atrackerClientV2Err error

	//Satellite link service
	satelliteLinkClient    *satellitelinkv1.SatelliteLinkV1
	satelliteLinkClientErr error

	esSchemaRegistryClient *schemaregistryv1.SchemaregistryV1
	esSchemaRegistryErr    error

	// Security and Compliance Center (SCC) Admin
	adminServiceApiClient    *adminserviceapiv1.AdminServiceApiV1
	adminServiceApiClientErr error

	// Security and Compliance Center (SCC) Governance
	configServiceApiClient    *configurationgovernancev1.ConfigurationGovernanceV1
	configServiceApiClientErr error

	//Security and Compliance Center (SCC) Compliance posture
	postureManagementClientErr error
	postureManagementClient    *posturemanagementv1.PostureManagementV1

	//Security and Compliance Center (SCC) Compliance posture v2
	postureManagementClientv2    *posturemanagementv2.PostureManagementV2
	postureManagementClientErrv2 error

	// context Based Restrictions (CBR)
	contextBasedRestrictionsClient    *contextbasedrestrictionsv1.ContextBasedRestrictionsV1
	contextBasedRestrictionsClientErr error

	// CD Toolchain
	cdToolchainClient    *cdtoolchainv2.CdToolchainV2
	cdToolchainClientErr error

	// CD Tekton Pipeline
	cdTektonPipelineClient    *cdtektonpipelinev2.CdTektonPipelineV2
	cdTektonPipelineClientErr error
}

// AppIDAPI provides AppID Service APIs ...
func (session clientSession) AppIDAPI() (*appid.AppIDManagementV4, error) {
	return session.appidAPI, session.appidErr
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

// UKO
func (session clientSession) UkoV4() (*ukov4.UkoV4, error) {
	return session.ukoClient, session.ukoClientErr
}

// UserManagementAPI provides User management APIs ...
func (sess clientSession) UserManagementAPI() (usermanagementv2.UserManagementAPI, error) {
	return sess.userManagementAPI, sess.userManagementErr
}

// IAM Policy Management
func (sess clientSession) IAMPolicyManagementV1API() (*iampolicymanagement.IamPolicyManagementV1, error) {
	return sess.iamPolicyManagementAPI, sess.iamPolicyManagementErr
}

// IAMAccessGroupsV2 provides IAM AG APIs ...
func (sess clientSession) IAMAccessGroupsV2() (*iamaccessgroups.IamAccessGroupsV2, error) {
	return sess.iamAccessGroupsAPI, sess.iamAccessGroupsErr
}

// IBM Cloud Shell
func (session clientSession) IBMCloudShellV1() (*ibmcloudshellv1.IBMCloudShellV1, error) {
	return session.ibmCloudShellClient, session.ibmCloudShellClientErr
}

// IcdAPI provides IBM Cloud Databases APIs ...
func (sess clientSession) ICDAPI() (icdv4.ICDServiceAPI, error) {
	return sess.icdServiceAPI, sess.icdConfigErr
}

// The IBM Cloud Databases API
func (session clientSession) CloudDatabasesV5() (*clouddatabasesv5.CloudDatabasesV5, error) {
	return session.cloudDatabasesClient, session.cloudDatabasesClientErr
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

func (session clientSession) EventNotificationsApiV1() (*eventnotificationsv1.EventNotificationsV1, error) {
	return session.eventNotificationsApiClient, session.eventNotificationsApiClientErr
}

func (session clientSession) AppConfigurationV1() (*appconfigurationv1.AppConfigurationV1, error) {
	return session.appConfigurationClient, session.appConfigurationClientErr
}

func (sess clientSession) KeyProtectAPI() (*kp.Client, error) {
	return sess.kpAPI, sess.kpErr
}

func (sess clientSession) KeyManagementAPI() (*kp.Client, error) {
	if sess.kmsErr == nil {
		var clientConfig *kp.ClientConfig
		if sess.kmsAPI.Config.APIKey != "" {
			clientConfig = &kp.ClientConfig{
				BaseURL:  EnvFallBack([]string{"IBMCLOUD_KP_API_ENDPOINT"}, sess.kmsAPI.Config.BaseURL),
				APIKey:   sess.kmsAPI.Config.APIKey, //pragma: allowlist secret
				Verbose:  kp.VerboseFailOnly,
				TokenURL: sess.kmsAPI.Config.TokenURL,
			}
		} else {
			clientConfig = &kp.ClientConfig{
				BaseURL:       EnvFallBack([]string{"IBMCLOUD_KP_API_ENDPOINT"}, sess.kmsAPI.Config.BaseURL),
				Authorization: sess.session.BluemixSession.Config.IAMAccessToken, //pragma: allowlist secret
				Verbose:       kp.VerboseFailOnly,
				TokenURL:      sess.kmsAPI.Config.TokenURL,
			}
		}

		kpClient, err := kp.New(*clientConfig, DefaultTransport())
		if err != nil {
			sess.kpErr = fmt.Errorf("[ERROR] Error occured while configuring Key Protect Service: %q", err)
		}
		return kpClient, nil
	}
	return sess.kmsAPI, sess.kmsErr
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
	return sess.ibmpiSession, sess.ibmpiConfigErr
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

// CIS Alerts
func (sess clientSession) CisAlertsSession() (*cisalertsv1.AlertsV1, error) {
	if sess.cisAlertsErr != nil {
		return sess.cisAlertsClient, sess.cisAlertsErr
	}
	return sess.cisAlertsClient.Clone(), nil
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

// CIS Authenticated Origin Pull
func (sess clientSession) CisOrigAuthSession() (*cisoriginpull.AuthenticatedOriginPullApiV1, error) {
	if sess.cisOriginAuthPullErr != nil {
		return sess.cisOriginAuthClient, sess.cisOriginAuthPullErr
	}
	return sess.cisOriginAuthClient.Clone(), nil
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

// Satellite Link
func (session clientSession) SatellitLinkClientSession() (*satellitelinkv1.SatelliteLinkV1, error) {
	return session.satelliteLinkClient, session.satelliteLinkClientErr
}

var cloudEndpoint = "cloud.ibm.com"

// Session to the Satellite client
func (sess clientSession) SatelliteClientSession() (*kubernetesserviceapiv1.KubernetesServiceApiV1, error) {
	return sess.satelliteClient, sess.satelliteClientErr
}

// CIS LogPushJob
func (sess clientSession) CisLogpushJobsSession() (*cislogpushjobsapiv1.LogpushJobsApiV1, error) {
	if sess.cisLogpushJobsErr != nil {
		return sess.cisLogpushJobsClient, sess.cisLogpushJobsErr
	}
	return sess.cisLogpushJobsClient.Clone(), nil
}

// CIS MTLS session
func (sess clientSession) CisMtlsSession() (*cismtlsv1.MtlsV1, error) {
	if sess.cisMtlsErr != nil {
		return sess.cisMtlsClient, sess.cisMtlsErr
	}
	return sess.cisMtlsClient.Clone(), nil
}

// CIS Webhooks
func (sess clientSession) CisWebhookSession() (*ciswebhooksv1.WebhooksV1, error) {
	if sess.cisWebhooksErr != nil {
		return sess.cisWebhooksClient, sess.cisWebhooksErr
	}
	return sess.cisWebhooksClient.Clone(), nil
}

// CIS Filters
func (sess clientSession) CisFiltersSession() (*cisfiltersv1.FiltersV1, error) {
	if sess.cisFiltersErr != nil {
		return sess.cisFiltersClient, sess.cisFiltersErr
	}
	return sess.cisFiltersClient.Clone(), nil
}

// CIS FirewallRules
func (sess clientSession) CisFirewallRulesSession() (*cisfirewallrulesv1.FirewallRulesV1, error) {
	if sess.cisFirewallRulesErr != nil {
		return sess.cisFirewallRulesClient, sess.cisFirewallRulesErr
	}
	return sess.cisFirewallRulesClient.Clone(), nil
}

// Activity Tracker API
func (session clientSession) AtrackerV1() (*atrackerv1.AtrackerV1, error) {
	return session.atrackerClient, session.atrackerClientErr
}

func (session clientSession) AtrackerV2() (*atrackerv2.AtrackerV2, error) {
	return session.atrackerClientV2, session.atrackerClientV2Err
}

func (session clientSession) ESschemaRegistrySession() (*schemaregistryv1.SchemaregistryV1, error) {
	return session.esSchemaRegistryClient, session.esSchemaRegistryErr
}

//Security and Compliance center Admin API
func (session clientSession) AdminServiceApiV1() (*adminserviceapiv1.AdminServiceApiV1, error) {
	return session.adminServiceApiClient, session.adminServiceApiClientErr
}

func (session clientSession) ConfigurationGovernanceV1() (*configurationgovernancev1.ConfigurationGovernanceV1, error) {
	return session.configServiceApiClient, session.configServiceApiClientErr
}

// Security and Compliance center Posture Management
func (session clientSession) PostureManagementV1() (*posturemanagementv1.PostureManagementV1, error) {
	if session.postureManagementClientErr != nil {
		return session.postureManagementClient, session.postureManagementClientErr
	}
	return session.postureManagementClient.Clone(), nil
}

//Security and Compliance center Posture Management v2
func (session clientSession) PostureManagementV2() (*posturemanagementv2.PostureManagementV2, error) {
	if session.postureManagementClientErrv2 != nil {
		return session.postureManagementClientv2, session.postureManagementClientErrv2
	}
	return session.postureManagementClientv2.Clone(), nil
}

// Context Based Restrictions
func (session clientSession) ContextBasedRestrictionsV1() (*contextbasedrestrictionsv1.ContextBasedRestrictionsV1, error) {
	return session.contextBasedRestrictionsClient, session.contextBasedRestrictionsClientErr
}

// CD Toolchain
func (session clientSession) CdToolchainV2() (*cdtoolchainv2.CdToolchainV2, error) {
	return session.cdToolchainClient, session.cdToolchainClientErr
}

// CD Tekton Pipeline
func (session clientSession) CdTektonPipelineV2() (*cdtektonpipelinev2.CdTektonPipelineV2, error) {
	return session.cdTektonPipelineClient, session.cdTektonPipelineClientErr
}

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
		session.iamAccessGroupsErr = errEmptyBluemixCredentials
		session.icdConfigErr = errEmptyBluemixCredentials
		session.resourceCatalogConfigErr = errEmptyBluemixCredentials
		session.resourceManagerErr = errEmptyBluemixCredentials
		session.resourceManagementConfigErrv2 = errEmptyBluemixCredentials
		session.resourceControllerConfigErr = errEmptyBluemixCredentials
		session.resourceControllerConfigErrv2 = errEmptyBluemixCredentials
		session.enterpriseManagementClientErr = errEmptyBluemixCredentials
		session.resourceControllerErr = errEmptyBluemixCredentials
		session.catalogManagementClientErr = errEmptyBluemixCredentials
		session.ibmpiConfigErr = errEmptyBluemixCredentials
		session.userManagementErr = errEmptyBluemixCredentials
		session.certManagementErr = errEmptyBluemixCredentials
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
		session.cisAlertsErr = errEmptyBluemixCredentials
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
		session.cisMtlsErr = errEmptyBluemixCredentials
		session.cisAccessRuleErr = errEmptyBluemixCredentials
		session.cisUARuleErr = errEmptyBluemixCredentials
		session.cisLockdownErr = errEmptyBluemixCredentials
		session.cisRangeAppErr = errEmptyBluemixCredentials
		session.cisWAFRuleErr = errEmptyBluemixCredentials
		session.iamIdentityErr = errEmptyBluemixCredentials
		session.secretsManagerClientErr = errEmptyBluemixCredentials
		session.cisFiltersErr = errEmptyBluemixCredentials
		session.cisWebhooksErr = errEmptyBluemixCredentials
		session.cisLogpushJobsErr = errEmptyBluemixCredentials
		session.schematicsClientErr = errEmptyBluemixCredentials
		session.satelliteClientErr = errEmptyBluemixCredentials
		session.iamPolicyManagementErr = errEmptyBluemixCredentials
		session.satelliteLinkClientErr = errEmptyBluemixCredentials
		session.esSchemaRegistryErr = errEmptyBluemixCredentials
		session.contextBasedRestrictionsClientErr = errEmptyBluemixCredentials
		session.postureManagementClientErr = errEmptyBluemixCredentials
		session.postureManagementClientErrv2 = errEmptyBluemixCredentials
		session.configServiceApiClientErr = errEmptyBluemixCredentials
		session.cdTektonPipelineClientErr = errEmptyBluemixCredentials
		session.cdToolchainClientErr = errEmptyBluemixCredentials

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
				session.bmxUserFetchErr = fmt.Errorf("[ERROR] Error occured while fetching auth key for account user details: %q", err)
				session.functionConfigErr = fmt.Errorf("[ERROR] Error occured while fetching auth key for function: %q", err)
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
				session.functionConfigErr = fmt.Errorf("[ERROR] Error occured while fetching auth key for function: %q", err)
			}
		}
	}

	if c.IAMTrustedProfileID == "" && sess.BluemixSession.Config.IAMAccessToken != "" && sess.BluemixSession.Config.BluemixAPIKey == "" {
		err := RefreshToken(sess.BluemixSession)
		if err != nil {
			for count := c.RetryCount; count >= 0; count-- {
				if err == nil || !isRetryable(err) {
					break
				}
				time.Sleep(c.RetryDelay)
				log.Printf("Retrying refresh token %d", count)
				err = RefreshToken(sess.BluemixSession)
			}
			if err != nil {
				return nil, fmt.Errorf("[ERROR] Error occured while refreshing the token: %q", err)
			}
		}

	}
	userConfig, err := fetchUserDetails(sess.BluemixSession, c.RetryCount, c.RetryDelay)
	if err != nil {
		session.bmxUserFetchErr = fmt.Errorf("[ERROR] Error occured while fetching account user details: %q", err)
	}
	session.bmxUserDetails = userConfig

	if sess.SoftLayerSession != nil && sess.SoftLayerSession.IAMToken != "" {
		sess.SoftLayerSession.IAMToken = sess.BluemixSession.Config.IAMAccessToken
		sess.SoftLayerSession.IAMRefreshToken = sess.BluemixSession.Config.IAMRefreshToken
	}

	session.functionClient, session.functionConfigErr = FunctionClient(sess.BluemixSession.Config)

	BluemixRegion = sess.BluemixSession.Config.Region
	var fileMap map[string]interface{}
	if f := EnvFallBack([]string{"IBMCLOUD_ENDPOINTS_FILE_PATH", "IC_ENDPOINTS_FILE_PATH"}, c.EndpointsFile); f != "" {
		jsonFile, err := os.Open(f)
		if err != nil {
			log.Fatalf("Unable to open Endpoints File %s", err)
		}
		defer jsonFile.Close()
		bytes, err := ioutil.ReadAll(jsonFile)
		if err != nil {
			log.Fatalf("Unable to read Endpoints File %s", err)
		}
		err = json.Unmarshal([]byte(bytes), &fileMap)
		if err != nil {
			log.Fatalf("Unable to unmarshal Endpoints File %s", err)
		}
	}
	accv1API, err := accountv1.New(sess.BluemixSession)
	if err != nil {
		session.accountV1ConfigErr = fmt.Errorf("[ERROR] Error occured while configuring Bluemix Accountv1 Service: %q", err)
	}
	session.bmxAccountv1ServiceAPI = accv1API

	accAPI, err := accountv2.New(sess.BluemixSession)
	if err != nil {
		session.accountConfigErr = fmt.Errorf("[ERROR] Error occured while configuring  Account Service: %q", err)
	}
	session.bmxAccountServiceAPI = accAPI

	cfAPI, err := mccpv2.New(sess.BluemixSession)
	if err != nil {
		session.cfConfigErr = fmt.Errorf("[ERROR] Error occured while configuring MCCP service: %q", err)
	}
	session.cfServiceAPI = cfAPI

	clusterAPI, err := containerv1.New(sess.BluemixSession)
	if err != nil {
		session.csConfigErr = fmt.Errorf("[ERROR] Error occured while configuring Container Service for K8s cluster: %q", err)
	}
	session.csServiceAPI = clusterAPI

	v2clusterAPI, err := containerv2.New(sess.BluemixSession)
	if err != nil {
		session.csv2ConfigErr = fmt.Errorf("[ERROR] Error occured while configuring vpc Container Service for K8s cluster: %q", err)
	}
	session.csv2ServiceAPI = v2clusterAPI

	hpcsAPI, err := hpcs.New(sess.BluemixSession)
	if err != nil {
		session.hpcsEndpointErr = fmt.Errorf("[ERROR] Error occured while configuring hpcs Endpoint: %q", err)
	}
	session.hpcsEndpointAPI = hpcsAPI

	kpurl := ContructEndpoint(fmt.Sprintf("%s.kms", c.Region), cloudEndpoint)
	if c.Visibility == "private" || c.Visibility == "public-and-private" {
		kpurl = ContructEndpoint(fmt.Sprintf("private.%s.kms", c.Region), cloudEndpoint)
	}
	if fileMap != nil && c.Visibility != "public-and-private" {
		kpurl = fileFallBack(fileMap, c.Visibility, "IBMCLOUD_KP_API_ENDPOINT", c.Region, kpurl)
	}
	var options kp.ClientConfig
	if c.BluemixAPIKey != "" {
		options = kp.ClientConfig{
			BaseURL: EnvFallBack([]string{"IBMCLOUD_KP_API_ENDPOINT"}, kpurl),
			APIKey:  sess.BluemixSession.Config.BluemixAPIKey, //pragma: allowlist secret
			// InstanceID:    "42fET57nnadurKXzXAedFLOhGqETfIGYxOmQXkFgkJV9",
			Verbose: kp.VerboseFailOnly,
		}

	} else {
		options = kp.ClientConfig{
			BaseURL:       EnvFallBack([]string{"IBMCLOUD_KP_API_ENDPOINT"}, kpurl),
			Authorization: sess.BluemixSession.Config.IAMAccessToken,
			// InstanceID:    "42fET57nnadurKXzXAedFLOhGqETfIGYxOmQXkFgkJV9",
			Verbose: kp.VerboseFailOnly,
		}
	}
	kpAPIclient, err := kp.New(options, DefaultTransport())
	if err != nil {
		session.kpErr = fmt.Errorf("[ERROR] Error occured while configuring Key Protect Service: %q", err)
	}
	session.kpAPI = kpAPIclient

	iamURL := iamidentity.DefaultServiceURL
	if c.Visibility == "private" || c.Visibility == "public-and-private" {
		if c.Region == "us-south" || c.Region == "us-east" {
			iamURL = ContructEndpoint(fmt.Sprintf("private.%s.iam", c.Region), cloudEndpoint)
		} else {
			iamURL = ContructEndpoint("private.iam", cloudEndpoint)
		}
	}
	if fileMap != nil && c.Visibility != "public-and-private" {
		iamURL = fileFallBack(fileMap, c.Visibility, "IBMCLOUD_IAM_API_ENDPOINT", c.Region, iamURL)
	}

	// KEY MANAGEMENT Service
	kmsurl := ContructEndpoint(fmt.Sprintf("%s.kms", c.Region), cloudEndpoint)
	if c.Visibility == "private" || c.Visibility == "public-and-private" {
		kmsurl = ContructEndpoint(fmt.Sprintf("private.%s.kms", c.Region), cloudEndpoint)
	}
	if fileMap != nil && c.Visibility != "public-and-private" {
		kmsurl = fileFallBack(fileMap, c.Visibility, "IBMCLOUD_KP_API_ENDPOINT", c.Region, kmsurl)
	}
	var kmsOptions kp.ClientConfig
	if c.BluemixAPIKey != "" {
		kmsOptions = kp.ClientConfig{
			BaseURL: EnvFallBack([]string{"IBMCLOUD_KP_API_ENDPOINT"}, kmsurl),
			APIKey:  sess.BluemixSession.Config.BluemixAPIKey, //pragma: allowlist secret
			// InstanceID:    "5af62d5d-5d90-4b84-bbcd-90d2123ae6c8",
			Verbose:  kp.VerboseFailOnly,
			TokenURL: EnvFallBack([]string{"IBMCLOUD_IAM_API_ENDPOINT"}, iamURL) + "/identity/token",
		}

	} else {
		kmsOptions = kp.ClientConfig{
			BaseURL:       EnvFallBack([]string{"IBMCLOUD_KP_API_ENDPOINT"}, kmsurl),
			Authorization: sess.BluemixSession.Config.IAMAccessToken,
			// InstanceID:    "5af62d5d-5d90-4b84-bbcd-90d2123ae6c8",
			Verbose:  kp.VerboseFailOnly,
			TokenURL: EnvFallBack([]string{"IBMCLOUD_IAM_API_ENDPOINT"}, iamURL) + "/identity/token",
		}
	}
	kmsAPIclient, err := kp.New(kmsOptions, DefaultTransport())
	if err != nil {
		session.kmsErr = fmt.Errorf("[ERROR] Error occured while configuring key Service: %q", err)
	}
	session.kmsAPI = kmsAPIclient

	var authenticator core.Authenticator

	if c.BluemixAPIKey != "" || sess.BluemixSession.Config.IAMRefreshToken != "" {
		if c.BluemixAPIKey != "" {
			authenticator = &core.IamAuthenticator{
				ApiKey: c.BluemixAPIKey,
				URL:    EnvFallBack([]string{"IBMCLOUD_IAM_API_ENDPOINT"}, iamURL),
			}
		} else {
			// Construct the IamAuthenticator with the IAM refresh token.
			authenticator = &core.IamAuthenticator{
				RefreshToken: sess.BluemixSession.Config.IAMRefreshToken,
				ClientId:     "bx",
				ClientSecret: "bx",
				URL:          EnvFallBack([]string{"IBMCLOUD_IAM_API_ENDPOINT"}, iamURL),
			}
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
	ukoClientOptions := &ukov4.UkoV4Options{
		Authenticator: authenticator,
	}

	// Construct the service client.
	session.ukoClient, err = ukov4.NewUkoV4(ukoClientOptions)
	if err == nil {
		// Enable retries for API calls
		session.ukoClient.Service.EnableRetries(c.RetryCount, c.RetryDelay)
		// Add custom header for analytics
		session.ukoClient.SetDefaultHeaders(gohttp.Header{
			"X-Original-User-Agent": {fmt.Sprintf("terraform-provider-ibm/%s", version.Version)},
		})
	} else {
		session.ukoClientErr = fmt.Errorf("Error occurred while configuring HPCS UKO service: %q", err)
	}

	// APPID Service
	appIDEndpoint := fmt.Sprintf("https://%s.appid.cloud.ibm.com", c.Region)
	if c.Visibility == "private" {
		session.appidErr = fmt.Errorf("App Id resources doesnot support private endpoints")
	}
	if fileMap != nil && c.Visibility != "public-and-private" {
		appIDEndpoint = fileFallBack(fileMap, c.Visibility, "IBMCLOUD_APPID_MANAGEMENT_API_ENDPOINT", c.Region, appIDEndpoint)
	}
	appIDClientOptions := &appid.AppIDManagementV4Options{
		Authenticator: authenticator,
		URL:           EnvFallBack([]string{"IBMCLOUD_APPID_MANAGEMENT_API_ENDPOINT"}, appIDEndpoint),
	}
	appIDClient, err := appid.NewAppIDManagementV4(appIDClientOptions)
	if err != nil {
		session.appidErr = fmt.Errorf("error occured while configuring AppID service: #{err}")
	}
	if appIDClient != nil && appIDClient.Service != nil {
		appIDClient.Service.EnableRetries(c.RetryCount, c.RetryDelay)
		appIDClient.SetDefaultHeaders(gohttp.Header{
			"X-Original-User-Agent": {fmt.Sprintf("terraform-provider-ibm/%s", version.Version)},
		})
	}
	session.appidAPI = appIDClient

	// Construct an "options" struct for creating Context Based Restrictions service client.
	cbrURL := contextbasedrestrictionsv1.DefaultServiceURL
	if c.Visibility == "private" || c.Visibility == "public-and-private" {
		session.contextBasedRestrictionsClientErr = fmt.Errorf("Context Based Restrictions Service API does not support private endpoints") //return this error if private endpoints are not supported
	}
	if fileMap != nil && c.Visibility != "public-and-private" {
		cbrURL = fileFallBack(fileMap, c.Visibility, "IBMCLOUD_CONTEXT_BASED_RESTRICTIONS_ENDPOINT", c.Region, cbrURL)
	}
	contextBasedRestrictionsClientOptions := &contextbasedrestrictionsv1.ContextBasedRestrictionsV1Options{
		Authenticator: authenticator,
		URL:           EnvFallBack([]string{"IBMCLOUD_CONTEXT_BASED_RESTRICTIONS_ENDPOINT"}, cbrURL),
	}

	// Construct the service client.
	session.contextBasedRestrictionsClient, err = contextbasedrestrictionsv1.NewContextBasedRestrictionsV1(contextBasedRestrictionsClientOptions)
	if err == nil && session.contextBasedRestrictionsClient != nil {
		// Enable retries for API calls
		session.contextBasedRestrictionsClient.Service.EnableRetries(c.RetryCount, c.RetryDelay)
		// Add custom header for analytics
		session.contextBasedRestrictionsClient.SetDefaultHeaders(gohttp.Header{
			"X-Original-User-Agent": {fmt.Sprintf("terraform-provider-ibm/%s", version.Version)},
		})
	} else {
		session.contextBasedRestrictionsClientErr = fmt.Errorf("[ERROR] Error occurred while configuring Context Based Restrictions service: %q", err)
	}

	// CATALOG MANAGEMENT Service
	catalogManagementURL := "https://cm.globalcatalog.cloud.ibm.com/api/v1-beta"
	if c.Visibility == "private" {
		session.catalogManagementClientErr = fmt.Errorf("Catalog Management resource doesnot support private endpoints")
	}
	if fileMap != nil && c.Visibility != "public-and-private" {
		catalogManagementURL = fileFallBack(fileMap, c.Visibility, "IBMCLOUD_CATALOG_MANAGEMENT_API_ENDPOINT", c.Region, catalogManagementURL)
	}
	catalogManagementClientOptions := &catalogmanagementv1.CatalogManagementV1Options{
		URL:           EnvFallBack([]string{"IBMCLOUD_CATALOG_MANAGEMENT_API_ENDPOINT"}, catalogManagementURL),
		Authenticator: authenticator,
	}
	// Construct the service client.
	session.catalogManagementClient, err = catalogmanagementv1.NewCatalogManagementV1(catalogManagementClientOptions)
	if err != nil {
		session.catalogManagementClientErr = fmt.Errorf("[ERROR] Error occurred while configuring Catalog Management API service: %q", err)
	}
	if session.catalogManagementClient != nil && session.catalogManagementClient.Service != nil {
		// Enable retries for API calls
		session.catalogManagementClient.Service.EnableRetries(c.RetryCount, c.RetryDelay)
		// Add custom header for analytics
		session.catalogManagementClient.SetDefaultHeaders(gohttp.Header{
			"X-Original-User-Agent": {fmt.Sprintf("terraform-provider-ibm/%s", version.Version)},
		})
	}

	// ATRACKER Service
	var atrackerClientURL string
	atrackerClientURL, err = atrackerv1.GetServiceURLForRegion(c.Region)
	if err != nil {
		session.atrackerClientErr = err
	}
	if c.Visibility == "private" || c.Visibility == "public-and-private" {
		atrackerClientURL, err = atrackerv1.GetServiceURLForRegion("private." + c.Region)
		if err != nil && c.Visibility == "public-and-private" {
			atrackerClientURL, err = atrackerv1.GetServiceURLForRegion(c.Region)
			if err != nil {
				session.atrackerClientErr = err
			}
		}
	}
	if fileMap != nil && c.Visibility != "public-and-private" {
		atrackerClientURL = fileFallBack(fileMap, c.Visibility, "IBMCLOUD_ATRACKER_API_ENDPOINT", c.Region, atrackerClientURL)
	}
	atrackerClientOptions := &atrackerv1.AtrackerV1Options{
		Authenticator: authenticator,
		URL:           EnvFallBack([]string{"IBMCLOUD_ATRACKER_API_ENDPOINT"}, atrackerClientURL),
	}
	// Construct the service client.
	session.atrackerClient, err = atrackerv1.NewAtrackerV1(atrackerClientOptions)
	if err != nil {
		session.atrackerClientErr = fmt.Errorf("[ERROR] Error occurred while configuring Activity Tracker API service: %q", err)
	}
	if session.atrackerClient != nil && session.atrackerClient.Service != nil {
		// Enable retries for API calls
		session.atrackerClient.Service.EnableRetries(c.RetryCount, c.RetryDelay)
		// Add custom header for analytics
		session.atrackerClient.SetDefaultHeaders(gohttp.Header{
			"X-Original-User-Agent": {fmt.Sprintf("terraform-provider-ibm/%s", version.Version)},
		})
	}
	// Version 2 Atracker
	var atrackerClientV2URL string
	if c.Visibility == "private" || c.Visibility == "public-and-private" {
		atrackerClientV2URL, err = atrackerv2.GetServiceURLForRegion("private." + c.Region)
		if err != nil && c.Visibility == "public-and-private" {
			atrackerClientV2URL, err = atrackerv2.GetServiceURLForRegion(c.Region)
		}
	} else {
		atrackerClientV2URL, err = atrackerv2.GetServiceURLForRegion(c.Region)
	}
	if err != nil {
		atrackerClientV2URL = atrackerv2.DefaultServiceURL
	}
	if fileMap != nil && c.Visibility != "public-and-private" {
		atrackerClientV2URL = fileFallBack(fileMap, c.Visibility, "IBMCLOUD_ATRACKER_API_ENDPOINT", c.Region, atrackerClientV2URL)
	}
	atrackerClientV2Options := &atrackerv2.AtrackerV2Options{
		Authenticator: authenticator,
		URL:           EnvFallBack([]string{"IBMCLOUD_ATRACKER_API_ENDPOINT"}, atrackerClientV2URL),
	}
	session.atrackerClientV2, err = atrackerv2.NewAtrackerV2(atrackerClientV2Options)
	if err == nil {
		// Enable retries for API calls
		session.atrackerClientV2.Service.EnableRetries(c.RetryCount, c.RetryDelay)
		// Add custom header for analytics
		session.atrackerClientV2.SetDefaultHeaders(gohttp.Header{
			"X-Original-User-Agent": {fmt.Sprintf("terraform-provider-ibm/%s", version.Version)},
		})
	} else {
		session.atrackerClientV2Err = fmt.Errorf("Error occurred while configuring Activity Tracker API Version 2 service: %q", err)
	}

	// SCC ADMIN Service
	var adminServiceApiClientURL string
	if c.Visibility == "private" || c.Visibility == "public-and-private" {
		adminServiceApiClientURL, err = adminserviceapiv1.GetServiceURLForRegion("private." + c.Region)
		if err != nil && c.Visibility == "public-and-private" {
			adminServiceApiClientURL, err = adminserviceapiv1.GetServiceURLForRegion(c.Region)
		}
	} else {
		adminServiceApiClientURL, err = adminserviceapiv1.GetServiceURLForRegion(c.Region)
	}
	if err != nil {
		adminServiceApiClientURL = adminserviceapiv1.DefaultServiceURL
	}
	adminServiceApiClientOptions := &adminserviceapiv1.AdminServiceApiV1Options{
		Authenticator: authenticator,
		URL:           EnvFallBack([]string{"IBMCLOUD_SCC_ADMIN_API_ENDPOINT"}, adminServiceApiClientURL),
	}

	// Construct the service client.
	session.adminServiceApiClient, err = adminserviceapiv1.NewAdminServiceApiV1(adminServiceApiClientOptions)
	if err == nil {
		// Enable retries for API calls
		session.adminServiceApiClient.Service.EnableRetries(c.RetryCount, c.RetryDelay)
		// Add custom header for analytics
		session.adminServiceApiClient.SetDefaultHeaders(gohttp.Header{
			"X-Original-User-Agent": {fmt.Sprintf("terraform-provider-ibm/%s", version.Version)},
		})
	} else {
		session.adminServiceApiClientErr = fmt.Errorf("[ERROR] Error occurred while configuring Admin Service API service: %q", err)
	}

	// SCHEMATICS Service
	// schematicsEndpoint := "https://schematics.cloud.ibm.com"
	schematicsEndpoint := ContructEndpoint(fmt.Sprintf("%s.schematics", c.Region), cloudEndpoint)
	if c.Visibility == "private" || c.Visibility == "public-and-private" {
		schematicsEndpoint = ContructEndpoint(fmt.Sprintf("private-%s.schematics", c.Region), cloudEndpoint)
	}
	if fileMap != nil && c.Visibility != "public-and-private" {
		schematicsEndpoint = fileFallBack(fileMap, c.Visibility, "IBMCLOUD_SCHEMATICS_API_ENDPOINT", c.Region, schematicsEndpoint)
	}
	schematicsClientOptions := &schematicsv1.SchematicsV1Options{
		Authenticator: authenticator,
		URL:           EnvFallBack([]string{"IBMCLOUD_SCHEMATICS_API_ENDPOINT"}, schematicsEndpoint),
	}
	// Construct the service client.
	schematicsClient, err := schematicsv1.NewSchematicsV1(schematicsClientOptions)
	if err != nil {
		session.schematicsClientErr = fmt.Errorf("[ERROR] Error occurred while configuring Schematics Service API service: %q", err)
	}
	// Enable retries for API calls
	if schematicsClient != nil && schematicsClient.Service != nil {
		schematicsClient.Service.EnableRetries(c.RetryCount, c.RetryDelay)
		schematicsClient.SetDefaultHeaders(gohttp.Header{
			"X-Original-User-Agent": {fmt.Sprintf("terraform-provider-ibm/%s", version.Version)},
		})
	}
	session.schematicsClient = schematicsClient

	// VPC Service
	vpcurl := ContructEndpoint(fmt.Sprintf("%s.iaas", c.Region), fmt.Sprintf("%s/v1", cloudEndpoint))
	if c.Visibility == "private" || c.Visibility == "public-and-private" {
		vpcurl = ContructEndpoint(fmt.Sprintf("%s.private.iaas", c.Region), fmt.Sprintf("%s/v1", cloudEndpoint))
	}
	if fileMap != nil && c.Visibility != "public-and-private" {
		vpcurl = fileFallBack(fileMap, c.Visibility, "IBMCLOUD_IS_NG_API_ENDPOINT", c.Region, vpcurl)
	}
	vpcoptions := &vpc.VpcV1Options{
		URL:           EnvFallBack([]string{"IBMCLOUD_IS_NG_API_ENDPOINT"}, vpcurl),
		Authenticator: authenticator,
	}
	vpcclient, err := vpc.NewVpcV1(vpcoptions)
	if err != nil {
		session.vpcErr = fmt.Errorf("[ERROR] Error occured while configuring vpc service: %q", err)
	}
	if vpcclient != nil && vpcclient.Service != nil {
		vpcclient.Service.EnableRetries(c.RetryCount, c.RetryDelay)
		vpcclient.SetDefaultHeaders(gohttp.Header{
			"X-Original-User-Agent": {fmt.Sprintf("terraform-provider-ibm/%s", version.Version)},
		})
	}
	session.vpcAPI = vpcclient

	// PUSH NOTIFICATIONS Service
	pnurl := fmt.Sprintf("https://%s.imfpush.cloud.ibm.com/imfpush/v1", c.Region)
	if c.Visibility == "private" {
		session.pushServiceClientErr = fmt.Errorf("Push Notifications Service API doesnot support private endpoints")
	}
	if fileMap != nil && c.Visibility != "public-and-private" {
		pnurl = fileFallBack(fileMap, c.Visibility, "IBMCLOUD_PUSH_API_ENDPOINT", c.Region, pnurl)
	}
	pushNotificationOptions := &pushservicev1.PushServiceV1Options{
		URL:           EnvFallBack([]string{"IBMCLOUD_PUSH_API_ENDPOINT"}, pnurl),
		Authenticator: authenticator,
	}
	pnclient, err := pushservicev1.NewPushServiceV1(pushNotificationOptions)
	if err != nil {
		session.pushServiceClientErr = fmt.Errorf("[ERROR] Error occured while configuring Push Notifications service: %q", err)
	}
	if pnclient != nil && pnclient.Service != nil {
		// Enable retries for API calls
		pnclient.Service.EnableRetries(c.RetryCount, c.RetryDelay)
		pnclient.SetDefaultHeaders(gohttp.Header{
			"X-Original-User-Agent": {fmt.Sprintf("terraform-provider-ibm/%s", version.Version)},
		})
	}
	session.pushServiceClient = pnclient
	// event notifications
	enurl := fmt.Sprintf("https://%s.event-notifications.cloud.ibm.com/event-notifications", c.Region)
	if c.Visibility == "private" {
		session.eventNotificationsApiClientErr = fmt.Errorf("Event Notifications Service does not support private endpoints")
	}
	if fileMap != nil && c.Visibility != "public-and-private" {
		enurl = fileFallBack(fileMap, c.Visibility, "IBMCLOUD_EVENT_NOTIFICATIONS_API_ENDPOINT", c.Region, enurl)
	}
	enClientOptions := &eventnotificationsv1.EventNotificationsV1Options{
		Authenticator: authenticator,
		URL:           EnvFallBack([]string{"IBMCLOUD_EVENT_NOTIFICATIONS_API_ENDPOINT"}, enurl),
	}
	// Construct the service client.
	session.eventNotificationsApiClient, err = eventnotificationsv1.NewEventNotificationsV1(enClientOptions)
	if err != nil {
		// Enable {
		session.eventNotificationsApiClientErr = fmt.Errorf("[ERROR] Error occurred while configuring Event Notifications service: %q", err)
	}
	if session.eventNotificationsApiClient != nil && session.eventNotificationsApiClient.Service != nil {
		// Enable retries for API calls
		session.eventNotificationsApiClient.Service.EnableRetries(c.RetryCount, c.RetryDelay)
		session.eventNotificationsApiClient.SetDefaultHeaders(gohttp.Header{
			"X-Original-User-Agent": {fmt.Sprintf("terraform-provider-ibm/%s", version.Version)},
		})
	}

	// APP CONFIGURATION Service
	if c.Visibility == "private" {
		session.appConfigurationClientErr = fmt.Errorf("[ERROR] App Configuration Service API doesnot support private endpoints")
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
		session.appConfigurationClientErr = fmt.Errorf("[ERROR] Error occurred while configuring App Configuration service: %q", err)
	}

	// CONTAINER REGISTRY Service
	// Construct an "options" struct for creating the service client.
	containerRegistryClientURL, err := containerregistryv1.GetServiceURLForRegion(c.Region)
	if err != nil {
		containerRegistryClientURL = containerregistryv1.DefaultServiceURL
	}
	if c.Visibility == "private" || c.Visibility == "public-and-private" {
		containerRegistryClientURL, err = GetPrivateServiceURLForRegion(c.Region)
		if err != nil {
			containerRegistryClientURL, _ = GetPrivateServiceURLForRegion("global")
		}
	}
	if fileMap != nil && c.Visibility != "public-and-private" {
		containerRegistryClientURL = fileFallBack(fileMap, c.Visibility, "IBMCLOUD_CR_API_ENDPOINT", c.Region, containerRegistryClientURL)
	}
	containerRegistryClientOptions := &containerregistryv1.ContainerRegistryV1Options{
		Authenticator: authenticator,
		URL:           EnvFallBack([]string{"IBMCLOUD_CR_API_ENDPOINT"}, containerRegistryClientURL),
		Account:       core.StringPtr(userConfig.UserAccount),
	}
	// Construct the service client.
	session.containerRegistryClient, err = containerregistryv1.NewContainerRegistryV1(containerRegistryClientOptions)
	if err != nil {
		session.containerRegistryClientErr = fmt.Errorf("[ERROR] Error occurred while configuring IBM Cloud Container Registry API service: %q", err)
	}
	if session.containerRegistryClient != nil && session.containerRegistryClient.Service != nil {
		// Enable retries for API calls
		session.containerRegistryClient.Service.EnableRetries(c.RetryCount, c.RetryDelay)
		// Add custom header for analytics
		session.containerRegistryClient.SetDefaultHeaders(gohttp.Header{
			"X-Original-User-Agent": {fmt.Sprintf("terraform-provider-ibm/%s", version.Version)},
		})
	}

	// OBJECT STORAGE Service
	cosconfigurl := "https://config.cloud-object-storage.cloud.ibm.com/v1"
	if fileMap != nil && c.Visibility != "public-and-private" {
		cosconfigurl = fileFallBack(fileMap, c.Visibility, "IBMCLOUD_COS_CONFIG_ENDPOINT", c.Region, cosconfigurl)
	}
	cosconfigoptions := &cosconfig.ResourceConfigurationV1Options{
		Authenticator: authenticator,
		URL:           EnvFallBack([]string{"IBMCLOUD_COS_CONFIG_ENDPOINT"}, cosconfigurl),
	}
	cosconfigclient, err := cosconfig.NewResourceConfigurationV1(cosconfigoptions)
	if err != nil {
		session.cosConfigErr = fmt.Errorf("[ERROR] Error occured while configuring COS config service: %q", err)
	}
	session.cosConfigAPI = cosconfigclient

	globalSearchAPI, err := globalsearchv2.New(sess.BluemixSession)
	if err != nil {
		session.globalSearchConfigErr = fmt.Errorf("[ERROR] Error occured while configuring Global Search: %q", err)
	}
	session.globalSearchServiceAPI = globalSearchAPI
	// Global Tagging Bluemix-go
	globalTaggingAPI, err := globaltaggingv3.New(sess.BluemixSession)
	if err != nil {
		session.globalTaggingConfigErr = fmt.Errorf("[ERROR] Error occured while configuring Global Tagging: %q", err)
	}
	session.globalTaggingServiceAPI = globalTaggingAPI

	// GLOBAL TAGGING Service
	globalTaggingEndpoint := "https://tags.global-search-tagging.cloud.ibm.com"
	if c.Visibility == "private" || c.Visibility == "public-and-private" {
		var globalTaggingRegion string
		if c.Region != "us-south" && c.Region != "us-east" {
			globalTaggingRegion = "us-south"
		} else {
			globalTaggingRegion = c.Region
		}
		globalTaggingEndpoint = ContructEndpoint(fmt.Sprintf("tags.private.%s", globalTaggingRegion), fmt.Sprintf("global-search-tagging.%s", cloudEndpoint))
	}
	if fileMap != nil && c.Visibility != "public-and-private" {
		globalTaggingEndpoint = fileFallBack(fileMap, c.Visibility, "IBMCLOUD_GT_API_ENDPOINT", c.Region, globalTaggingEndpoint)
	}
	globalTaggingV1Options := &globaltaggingv1.GlobalTaggingV1Options{
		URL:           EnvFallBack([]string{"IBMCLOUD_GT_API_ENDPOINT"}, globalTaggingEndpoint),
		Authenticator: authenticator,
	}
	globalTaggingAPIV1, err := globaltaggingv1.NewGlobalTaggingV1(globalTaggingV1Options)
	if err != nil {
		session.globalTaggingConfigErrV1 = fmt.Errorf("[ERROR] Error occured while configuring Global Tagging: %q", err)
	}
	if globalTaggingAPIV1 != nil && globalTaggingAPIV1.Service != nil {
		session.globalTaggingServiceAPIV1 = *globalTaggingAPIV1
		session.globalTaggingServiceAPIV1.Service.EnableRetries(c.RetryCount, c.RetryDelay)
		session.globalTaggingServiceAPIV1.SetDefaultHeaders(gohttp.Header{
			"X-Original-User-Agent": {fmt.Sprintf("terraform-provider-ibm/%s", version.Version)},
		})
	}

	icdAPI, err := icdv4.New(sess.BluemixSession)
	if err != nil {
		session.icdConfigErr = fmt.Errorf("[ERROR] Error occured while configuring IBM Cloud Database Services: %q", err)
	}
	session.icdServiceAPI = icdAPI

	var cloudDatabasesEndpoint string

	if c.Visibility == "private" || c.Visibility == "public-and-private" {
		cloudDatabasesEndpoint = fmt.Sprintf("https://api.%s.private.databases.cloud.ibm.com/v5/ibm", c.Region)
	} else {
		cloudDatabasesEndpoint = fmt.Sprintf("https://api.%s.databases.cloud.ibm.com/v5/ibm", c.Region)
	}

	// Construct an "options" struct for creating the service client.
	cloudDatabasesClientOptions := &clouddatabasesv5.CloudDatabasesV5Options{
		URL:           EnvFallBack([]string{"IBMCLOUD_DATABASES_API_ENDPOINT"}, cloudDatabasesEndpoint),
		Authenticator: authenticator,
	}

	// Construct the service client.
	session.cloudDatabasesClient, err = clouddatabasesv5.NewCloudDatabasesV5(cloudDatabasesClientOptions)
	if err == nil {
		// Enable retries for API calls
		session.cloudDatabasesClient.Service.EnableRetries(c.RetryCount, c.RetryDelay)
		// Add custom header for analytics
		session.cloudDatabasesClient.SetDefaultHeaders(gohttp.Header{
			"X-Original-User-Agent": {fmt.Sprintf("terraform-provider-ibm/%s", version.Version)},
		})
	} else {
		session.cloudDatabasesClientErr = fmt.Errorf("Error occurred while configuring The IBM Cloud Databases API service: %q", err)
	}

	resourceCatalogAPI, err := catalog.New(sess.BluemixSession)
	if err != nil {
		session.resourceCatalogConfigErr = fmt.Errorf("[ERROR] Error occured while configuring Resource Catalog service: %q", err)
	}
	session.resourceCatalogServiceAPI = resourceCatalogAPI

	resourceManagementAPIv2, err := managementv2.New(sess.BluemixSession)
	if err != nil {
		session.resourceManagementConfigErrv2 = fmt.Errorf("[ERROR] Error occured while configuring Resource Management service: %q", err)
	}
	session.resourceManagementServiceAPIv2 = resourceManagementAPIv2

	resourceControllerAPI, err := controller.New(sess.BluemixSession)
	if err != nil {
		session.resourceControllerConfigErr = fmt.Errorf("[ERROR] Error occured while configuring Resource Controller service: %q", err)
	}
	session.resourceControllerServiceAPI = resourceControllerAPI

	ResourceControllerAPIv2, err := controllerv2.New(sess.BluemixSession)
	if err != nil {
		session.resourceControllerConfigErrv2 = fmt.Errorf("[ERROR] Error occured while configuring Resource Controller v2 service: %q", err)
	}
	session.resourceControllerServiceAPIv2 = ResourceControllerAPIv2

	userManagementAPI, err := usermanagementv2.New(sess.BluemixSession)
	if err != nil {
		session.userManagementErr = fmt.Errorf("[ERROR] Error occured while configuring user management service: %q", err)
	}
	session.userManagementAPI = userManagementAPI

	certManagementAPI, err := certificatemanager.New(sess.BluemixSession)
	if err != nil {
		session.certManagementErr = fmt.Errorf("[ERROR] Error occured while configuring Certificate manager service: %q", err)
	}
	session.certManagementAPI = certManagementAPI

	namespaceFunction, err := functions.New(sess.BluemixSession)
	if err != nil {
		session.functionIAMNamespaceErr = fmt.Errorf("[ERROR] Error occured while configuring Cloud Funciton Service : %q", err)
	}
	session.functionIAMNamespaceAPI = namespaceFunction

	//  API GATEWAY service
	apicurl := ContructEndpoint(fmt.Sprintf("api.%s.apigw", c.Region), fmt.Sprintf("%s/controller", cloudEndpoint))
	if c.Visibility == "private" || c.Visibility == "public-and-private" {
		apicurl = ContructEndpoint(fmt.Sprintf("api.private.%s.apigw", c.Region), fmt.Sprintf("%s/controller", cloudEndpoint))
	}
	if fileMap != nil && c.Visibility != "public-and-private" {
		apicurl = fileFallBack(fileMap, c.Visibility, "IBMCLOUD_API_GATEWAY_ENDPOINT", c.Region, apicurl)
	}
	APIGatewayControllerAPIV1Options := &apigateway.ApiGatewayControllerApiV1Options{
		URL:           EnvFallBack([]string{"IBMCLOUD_API_GATEWAY_ENDPOINT"}, apicurl),
		Authenticator: &core.NoAuthAuthenticator{},
	}
	apigatewayAPI, err := apigateway.NewApiGatewayControllerApiV1(APIGatewayControllerAPIV1Options)
	if err != nil {
		session.apigatewayErr = fmt.Errorf("[ERROR] Error occured while configuring  APIGateway service: %q", err)
	}
	session.apigatewayAPI = apigatewayAPI

	// POWER SYSTEMS Service
	piURL := ContructEndpoint(c.Region, "power-iaas.cloud.ibm.com")
	ibmPIOptions := &ibmpisession.IBMPIOptions{
		Authenticator: authenticator,
		Debug:         os.Getenv("TF_LOG") != "",
		Region:        c.Region,
		URL:           EnvFallBack([]string{"IBMCLOUD_PI_API_ENDPOINT"}, piURL),
		UserAccount:   userConfig.UserAccount,
		Zone:          c.Zone,
	}
	ibmpisession, err := ibmpisession.NewIBMPISession(ibmPIOptions)
	if err != nil {
		session.ibmpiConfigErr = fmt.Errorf("Error occured while configuring ibmpisession: %q", err)
	}
	session.ibmpiSession = ibmpisession

	// PRIVATE DNS Service
	pdnsURL := dns.DefaultServiceURL
	if c.Visibility == "private" || c.Visibility == "public-and-private" {
		pdnsURL = ContructEndpoint("api.private.dns-svcs", fmt.Sprintf("%s/v1", cloudEndpoint))
	}
	if fileMap != nil && c.Visibility != "public-and-private" {
		pdnsURL = fileFallBack(fileMap, c.Visibility, "IBMCLOUD_PRIVATE_DNS_API_ENDPOINT", c.Region, pdnsURL)
	}
	dnsOptions := &dns.DnsSvcsV1Options{
		URL:           EnvFallBack([]string{"IBMCLOUD_PRIVATE_DNS_API_ENDPOINT"}, pdnsURL),
		Authenticator: authenticator,
	}
	session.pDNSClient, session.pDNSErr = dns.NewDnsSvcsV1(dnsOptions)
	if session.pDNSErr != nil {
		session.pDNSErr = fmt.Errorf("[ERROR] Error occured while configuring PrivateDNS Service: %s", session.pDNSErr)
	}
	if session.pDNSClient != nil && session.pDNSClient.Service != nil {
		session.pDNSClient.Service.EnableRetries(c.RetryCount, c.RetryDelay)
		session.pDNSClient.SetDefaultHeaders(gohttp.Header{
			"X-Original-User-Agent": {fmt.Sprintf("terraform-provider-ibm/%s", version.Version)},
		})
	}

	// DIRECT LINK Service
	ver := time.Now().Format("2006-01-02")
	dlURL := dl.DefaultServiceURL
	if c.Visibility == "private" || c.Visibility == "public-and-private" {
		dlURL = ContructEndpoint("private.directlink", fmt.Sprintf("%s/v1", cloudEndpoint))
	}
	if fileMap != nil && c.Visibility != "public-and-private" {
		dlURL = fileFallBack(fileMap, c.Visibility, "IBMCLOUD_DL_API_ENDPOINT", c.Region, dlURL)
	}
	directlinkOptions := &dl.DirectLinkV1Options{
		URL:           EnvFallBack([]string{"IBMCLOUD_DL_API_ENDPOINT"}, dlURL),
		Authenticator: authenticator,
		Version:       &ver,
	}
	session.directlinkAPI, session.directlinkErr = dl.NewDirectLinkV1(directlinkOptions)
	if session.directlinkErr != nil {
		session.directlinkErr = fmt.Errorf("[ERROR] Error occured while configuring Direct Link Service: %s", session.directlinkErr)
	}
	if session.directlinkAPI != nil && session.directlinkAPI.Service != nil {
		session.directlinkAPI.Service.EnableRetries(c.RetryCount, c.RetryDelay)
		session.directlinkAPI.SetDefaultHeaders(gohttp.Header{
			"X-Original-User-Agent": {fmt.Sprintf("terraform-provider-ibm/%s", version.Version)},
		})
	}

	// DIRECT LINK PROVIDER Service
	dlproviderURL := dlProviderV2.DefaultServiceURL
	if c.Visibility == "private" || c.Visibility == "public-and-private" {
		dlproviderURL = ContructEndpoint("private.directlink", fmt.Sprintf("%s/provider/v2", cloudEndpoint))
	}
	if fileMap != nil && c.Visibility != "public-and-private" {
		dlproviderURL = fileFallBack(fileMap, c.Visibility, "IBMCLOUD_DL_PROVIDER_API_ENDPOINT", c.Region, dlproviderURL)
	}
	directLinkProviderV2Options := &dlProviderV2.DirectLinkProviderV2Options{
		URL:           EnvFallBack([]string{"IBMCLOUD_DL_PROVIDER_API_ENDPOINT"}, dlproviderURL),
		Authenticator: authenticator,
		Version:       &ver,
	}
	session.dlProviderAPI, session.dlProviderErr = dlProviderV2.NewDirectLinkProviderV2(directLinkProviderV2Options)
	if session.dlProviderErr != nil {
		session.dlProviderErr = fmt.Errorf("[ERROR] Error occured while configuring Direct Link Provider Service: %s", session.dlProviderErr)
	}
	if session.dlProviderAPI != nil && session.dlProviderAPI.Service != nil {
		session.dlProviderAPI.Service.EnableRetries(c.RetryCount, c.RetryDelay)
		session.dlProviderAPI.SetDefaultHeaders(gohttp.Header{
			"X-Original-User-Agent": {fmt.Sprintf("terraform-provider-ibm/%s", version.Version)},
		})
	}

	// TRANSIT GATEWAY Service
	tgURL := tg.DefaultServiceURL
	if c.Visibility == "private" || c.Visibility == "public-and-private" {
		tgURL = ContructEndpoint("private.transit", fmt.Sprintf("%s/v1", cloudEndpoint))
	}
	if fileMap != nil && c.Visibility != "public-and-private" {
		tgURL = fileFallBack(fileMap, c.Visibility, "IBMCLOUD_TG_API_ENDPOINT", c.Region, tgURL)
	}
	transitgatewayOptions := &tg.TransitGatewayApisV1Options{
		URL:           EnvFallBack([]string{"IBMCLOUD_TG_API_ENDPOINT"}, tgURL),
		Authenticator: authenticator,
		Version:       CreateVersionDate(),
	}
	session.transitgatewayAPI, session.transitgatewayErr = tg.NewTransitGatewayApisV1(transitgatewayOptions)
	if session.transitgatewayErr != nil {
		session.transitgatewayErr = fmt.Errorf("[ERROR] Error occured while configuring Transit Gateway Service: %s", session.transitgatewayErr)
	}
	if session.transitgatewayAPI != nil && session.transitgatewayAPI.Service != nil {
		session.transitgatewayAPI.Service.EnableRetries(c.RetryCount, c.RetryDelay)
		// session.transitgatewayAPI.SetDefaultHeaders(gohttp.Header{
		// 	"X-Original-User-Agent": {fmt.Sprintf("terraform-provider-ibm/%s", version.Version)},
		// })
	}

	// CIS Service instances starts here.
	cisURL := ContructEndpoint("api.cis", cloudEndpoint)
	if c.Visibility == "private" {
		// cisURL = ContructEndpoint("api.private.cis", cloudEndpoint)
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
		session.cisFiltersErr = fmt.Errorf("CIS Service doesnt support private endpoints.")
		session.cisWebhooksErr = fmt.Errorf("CIS Service doesnt support private endpoints.")
		session.cisMtlsErr = fmt.Errorf("CIS Service doesnt support private endpoints.")

	}
	if fileMap != nil && c.Visibility != "public-and-private" {
		cisURL = fileFallBack(fileMap, c.Visibility, "IBMCLOUD_CIS_API_ENDPOINT", c.Region, cisURL)
	}
	cisEndPoint := EnvFallBack([]string{"IBMCLOUD_CIS_API_ENDPOINT"}, cisURL)

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
		session.cisZonesV1Client.SetDefaultHeaders(gohttp.Header{
			"X-Original-User-Agent": {fmt.Sprintf("terraform-provider-ibm/%s", version.Version)},
		})
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
		session.cisDNSErr = fmt.Errorf("[ERROR] Error occured while configuring CIS DNS Service: %s", session.cisDNSErr)
	}
	if session.cisDNSRecordsClient != nil && session.cisDNSRecordsClient.Service != nil {
		session.cisDNSRecordsClient.Service.EnableRetries(c.RetryCount, c.RetryDelay)
		session.cisDNSRecordsClient.SetDefaultHeaders(gohttp.Header{
			"X-Original-User-Agent": {fmt.Sprintf("terraform-provider-ibm/%s", version.Version)},
		})
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
		session.cisDNSRecordBulkClient.SetDefaultHeaders(gohttp.Header{
			"X-Original-User-Agent": {fmt.Sprintf("terraform-provider-ibm/%s", version.Version)},
		})
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
			fmt.Errorf("[ERROR] Error occured while configuring CIS GLB Pool service: %s",
				session.cisGLBPoolErr)
	}
	if session.cisGLBPoolClient != nil && session.cisGLBPoolClient.Service != nil {
		session.cisGLBPoolClient.Service.EnableRetries(c.RetryCount, c.RetryDelay)
		session.cisGLBPoolClient.SetDefaultHeaders(gohttp.Header{
			"X-Original-User-Agent": {fmt.Sprintf("terraform-provider-ibm/%s", version.Version)},
		})
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
			fmt.Errorf("[ERROR] Error occured while configuring CIS GLB service: %s",
				session.cisGLBErr)
	}
	if session.cisGLBClient != nil && session.cisGLBClient.Service != nil {
		session.cisGLBClient.Service.EnableRetries(c.RetryCount, c.RetryDelay)
		session.cisGLBClient.SetDefaultHeaders(gohttp.Header{
			"X-Original-User-Agent": {fmt.Sprintf("terraform-provider-ibm/%s", version.Version)},
		})
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
			fmt.Errorf("[ERROR] Error occured while configuring CIS GLB Health Check service: %s",
				session.cisGLBHealthCheckErr)
	}
	if session.cisGLBHealthCheckClient != nil && session.cisGLBHealthCheckClient.Service != nil {
		session.cisGLBHealthCheckClient.Service.EnableRetries(c.RetryCount, c.RetryDelay)
		session.cisGLBHealthCheckClient.SetDefaultHeaders(gohttp.Header{
			"X-Original-User-Agent": {fmt.Sprintf("terraform-provider-ibm/%s", version.Version)},
		})
	}

	// IBM Network CIS IP
	cisIPOpt := &cisipv1.CisIpApiV1Options{
		URL:           cisEndPoint,
		Authenticator: authenticator,
	}
	session.cisIPClient, session.cisIPErr = cisipv1.NewCisIpApiV1(cisIPOpt)
	if session.cisIPErr != nil {
		session.cisIPErr = fmt.Errorf("[ERROR] Error occured while configuring CIS IP service: %s",
			session.cisIPErr)
	}
	if session.cisIPClient != nil && session.cisIPClient.Service != nil {
		session.cisIPClient.Service.EnableRetries(c.RetryCount, c.RetryDelay)
		session.cisIPClient.SetDefaultHeaders(gohttp.Header{
			"X-Original-User-Agent": {fmt.Sprintf("terraform-provider-ibm/%s", version.Version)},
		})
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
		session.cisRLClient.SetDefaultHeaders(gohttp.Header{
			"X-Original-User-Agent": {fmt.Sprintf("terraform-provider-ibm/%s", version.Version)},
		})
	}
	// IBM Network CIS Alerts
	cisAlertsOpt := &cisalertsv1.AlertsV1Options{
		URL:           cisEndPoint,
		Crn:           core.StringPtr(""),
		Authenticator: authenticator,
	}
	session.cisAlertsClient, session.cisAlertsErr = cisalertsv1.NewAlertsV1(cisAlertsOpt)
	if session.cisAlertsErr != nil {
		session.cisAlertsErr =
			fmt.Errorf("[ERROR] Error occured while configuring CIS Alerts : %s",
				session.cisAlertsErr)
	}
	if session.cisAlertsClient != nil && session.cisAlertsClient.Service != nil {
		session.cisAlertsClient.Service.EnableRetries(c.RetryCount, c.RetryDelay)
		session.cisAlertsClient.SetDefaultHeaders(gohttp.Header{
			"X-Original-User-Agent": {fmt.Sprintf("terraform-provider-ibm/%s", version.Version)},
		})
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
		session.cisPageRuleClient.SetDefaultHeaders(gohttp.Header{
			"X-Original-User-Agent": {fmt.Sprintf("terraform-provider-ibm/%s", version.Version)},
		})
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
			fmt.Errorf("[ERROR] Error occured while configuring CIS Edge Function service: %s",
				session.cisEdgeFunctionErr)
	}
	if session.cisEdgeFunctionClient != nil && session.cisEdgeFunctionClient.Service != nil {
		session.cisEdgeFunctionClient.Service.EnableRetries(c.RetryCount, c.RetryDelay)
		session.cisEdgeFunctionClient.SetDefaultHeaders(gohttp.Header{
			"X-Original-User-Agent": {fmt.Sprintf("terraform-provider-ibm/%s", version.Version)},
		})
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
			fmt.Errorf("[ERROR] Error occured while configuring CIS SSL certificate service: %s",
				session.cisSSLErr)
	}
	if session.cisSSLClient != nil && session.cisSSLClient.Service != nil {
		session.cisSSLClient.Service.EnableRetries(c.RetryCount, c.RetryDelay)
		session.cisSSLClient.SetDefaultHeaders(gohttp.Header{
			"X-Original-User-Agent": {fmt.Sprintf("terraform-provider-ibm/%s", version.Version)},
		})
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
			fmt.Errorf("[ERROR] Error occured while configuration CIS WAF Package service: %s",
				session.cisWAFPackageErr)
	}
	if session.cisWAFPackageClient != nil && session.cisWAFPackageClient.Service != nil {
		session.cisWAFPackageClient.Service.EnableRetries(c.RetryCount, c.RetryDelay)
		session.cisWAFPackageClient.SetDefaultHeaders(gohttp.Header{
			"X-Original-User-Agent": {fmt.Sprintf("terraform-provider-ibm/%s", version.Version)},
		})
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
			fmt.Errorf("[ERROR] Error occured while configuring CIS Domain Settings service: %s",
				session.cisDomainSettingsErr)
	}
	if session.cisDomainSettingsClient != nil && session.cisDomainSettingsClient.Service != nil {
		session.cisDomainSettingsClient.Service.EnableRetries(c.RetryCount, c.RetryDelay)
		session.cisDomainSettingsClient.SetDefaultHeaders(gohttp.Header{
			"X-Original-User-Agent": {fmt.Sprintf("terraform-provider-ibm/%s", version.Version)},
		})
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
			fmt.Errorf("[ERROR] Error occured while configuring CIS Routing service: %s",
				session.cisRoutingErr)
	}
	if session.cisRoutingClient != nil && session.cisRoutingClient.Service != nil {
		session.cisRoutingClient.Service.EnableRetries(c.RetryCount, c.RetryDelay)
		session.cisRoutingClient.SetDefaultHeaders(gohttp.Header{
			"X-Original-User-Agent": {fmt.Sprintf("terraform-provider-ibm/%s", version.Version)},
		})
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
			fmt.Errorf("[ERROR] Error occured while configuring CIS WAF Group service: %s",
				session.cisWAFGroupErr)
	}
	if session.cisWAFGroupClient != nil && session.cisWAFGroupClient.Service != nil {
		session.cisWAFGroupClient.Service.EnableRetries(c.RetryCount, c.RetryDelay)
		session.cisWAFGroupClient.SetDefaultHeaders(gohttp.Header{
			"X-Original-User-Agent": {fmt.Sprintf("terraform-provider-ibm/%s", version.Version)},
		})
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
			fmt.Errorf("[ERROR] Error occured while configuring CIS Caching service: %s",
				session.cisCacheErr)
	}
	if session.cisCacheClient != nil && session.cisCacheClient.Service != nil {
		session.cisCacheClient.Service.EnableRetries(c.RetryCount, c.RetryDelay)
		session.cisCacheClient.SetDefaultHeaders(gohttp.Header{
			"X-Original-User-Agent": {fmt.Sprintf("terraform-provider-ibm/%s", version.Version)},
		})
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
			fmt.Errorf("[ERROR] Error occured while configuring CIS Custom Pages service: %s",
				session.cisCustomPageErr)
	}
	if session.cisCustomPageClient != nil && session.cisCustomPageClient.Service != nil {
		session.cisCustomPageClient.Service.EnableRetries(c.RetryCount, c.RetryDelay)
		session.cisCustomPageClient.SetDefaultHeaders(gohttp.Header{
			"X-Original-User-Agent": {fmt.Sprintf("terraform-provider-ibm/%s", version.Version)},
		})
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
			fmt.Errorf("[ERROR] Error occured while configuring CIS Firewall Access Rule service: %s",
				session.cisAccessRuleErr)
	}
	if session.cisAccessRuleClient != nil && session.cisAccessRuleClient.Service != nil {
		session.cisAccessRuleClient.Service.EnableRetries(c.RetryCount, c.RetryDelay)
		session.cisAccessRuleClient.SetDefaultHeaders(gohttp.Header{
			"X-Original-User-Agent": {fmt.Sprintf("terraform-provider-ibm/%s", version.Version)},
		})
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
			fmt.Errorf("[ERROR] Error occured while configuring CIS Firewall User Agent Blocking Rule service: %s",
				session.cisUARuleErr)
	}
	if session.cisUARuleClient != nil && session.cisUARuleClient.Service != nil {
		session.cisUARuleClient.Service.EnableRetries(c.RetryCount, c.RetryDelay)
		session.cisUARuleClient.SetDefaultHeaders(gohttp.Header{
			"X-Original-User-Agent": {fmt.Sprintf("terraform-provider-ibm/%s", version.Version)},
		})
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
			fmt.Errorf("[ERROR] Error occured while configuring CIS Firewall Lockdown Rule service: %s",
				session.cisLockdownErr)
	}
	if session.cisLockdownClient != nil && session.cisLockdownClient.Service != nil {
		session.cisLockdownClient.Service.EnableRetries(c.RetryCount, c.RetryDelay)
		session.cisLockdownClient.SetDefaultHeaders(gohttp.Header{
			"X-Original-User-Agent": {fmt.Sprintf("terraform-provider-ibm/%s", version.Version)},
		})
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
			fmt.Errorf("[ERROR] Error occured while configuring CIS Range Application rule service: %s",
				session.cisRangeAppErr)
	}
	if session.cisRangeAppClient != nil && session.cisRangeAppClient.Service != nil {
		session.cisRangeAppClient.Service.EnableRetries(c.RetryCount, c.RetryDelay)
		session.cisRangeAppClient.SetDefaultHeaders(gohttp.Header{
			"X-Original-User-Agent": {fmt.Sprintf("terraform-provider-ibm/%s", version.Version)},
		})
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
		session.cisWAFRuleClient.SetDefaultHeaders(gohttp.Header{
			"X-Original-User-Agent": {fmt.Sprintf("terraform-provider-ibm/%s", version.Version)},
		})
	}

	// IBM Network CIS LogpushJobs
	cisLogpushJobOpt := &cislogpushjobsapiv1.LogpushJobsApiV1Options{
		URL:           cisEndPoint,
		Crn:           core.StringPtr(""),
		ZoneID:        core.StringPtr(""),
		Dataset:       core.StringPtr(""),
		Authenticator: authenticator,
	}
	session.cisLogpushJobsClient, session.cisLogpushJobsErr = cislogpushjobsapiv1.NewLogpushJobsApiV1(cisLogpushJobOpt)
	if session.cisLogpushJobsErr != nil {
		session.cisLogpushJobsErr =
			fmt.Errorf("[ERROR] Error occured while configuring CIS LogpushJobs : %s",
				session.cisLogpushJobsErr)
	}
	if session.cisLogpushJobsClient != nil && session.cisLogpushJobsClient.Service != nil {
		session.cisLogpushJobsClient.Service.EnableRetries(c.RetryCount, c.RetryDelay)
		session.cisLogpushJobsClient.SetDefaultHeaders(gohttp.Header{
			"X-Original-User-Agent": {fmt.Sprintf("terraform-provider-ibm/%s", version.Version)},
		})
	}

	// IBM MTLS Session
	cisMtlsOpt := &cismtlsv1.MtlsV1Options{
		URL:           cisEndPoint,
		Crn:           core.StringPtr(""),
		Authenticator: authenticator,
	}
	session.cisMtlsClient, session.cisMtlsErr = cismtlsv1.NewMtlsV1(cisMtlsOpt)
	if session.cisMtlsErr != nil {
		session.cisMtlsErr =
			fmt.Errorf("[ERROR] Error occured while configuring CIS MTLS : %s",
				session.cisMtlsErr)
	}
	if session.cisMtlsClient != nil && session.cisMtlsClient.Service != nil {
		session.cisMtlsClient.Service.EnableRetries(c.RetryCount, c.RetryDelay)
		session.cisMtlsClient.SetDefaultHeaders(gohttp.Header{
			"X-Original-User-Agent": {fmt.Sprintf("terraform-provider-ibm/%s", version.Version)},
		})
	}

	// IBM Network CIS Webhooks
	cisWebhooksOpt := &ciswebhooksv1.WebhooksV1Options{
		URL:           cisEndPoint,
		Crn:           core.StringPtr(""),
		Authenticator: authenticator,
	}
	session.cisWebhooksClient, session.cisWebhooksErr = ciswebhooksv1.NewWebhooksV1(cisWebhooksOpt)
	if session.cisWebhooksErr != nil {
		session.cisWebhooksErr =
			fmt.Errorf("[ERROR] Error occured while configuring CIS Webhooks : %s",
				session.cisWebhooksErr)
	}
	if session.cisWebhooksClient != nil && session.cisWebhooksClient.Service != nil {
		session.cisWebhooksClient.Service.EnableRetries(c.RetryCount, c.RetryDelay)
		session.cisWebhooksClient.SetDefaultHeaders(gohttp.Header{
			"X-Original-User-Agent": {fmt.Sprintf("terraform-provider-ibm/%s", version.Version)},
		})
	}
	// IBM Network CIS Filters
	cisFiltersOpt := &cisfiltersv1.FiltersV1Options{
		URL:           cisEndPoint,
		Authenticator: authenticator,
	}
	session.cisFiltersClient, session.cisFiltersErr = cisfiltersv1.NewFiltersV1(cisFiltersOpt)
	if session.cisFiltersErr != nil {
		session.cisFiltersErr =
			fmt.Errorf("[ERROR] Error occured while configuring CIS Filters : %s",
				session.cisFiltersErr)
	}
	if session.cisFiltersClient != nil && session.cisFiltersClient.Service != nil {
		session.cisFiltersClient.Service.EnableRetries(c.RetryCount, c.RetryDelay)
		session.cisFiltersClient.SetDefaultHeaders(gohttp.Header{
			"X-Original-User-Agent": {fmt.Sprintf("terraform-provider-ibm/%s", version.Version)},
		})
	}

	// IBM Network CIS Firewall rules
	cisFirewallrulesOpt := &cisfirewallrulesv1.FirewallRulesV1Options{
		URL:           cisEndPoint,
		Authenticator: authenticator,
	}
	session.cisFirewallRulesClient, session.cisFirewallRulesErr = cisfirewallrulesv1.NewFirewallRulesV1(cisFirewallrulesOpt)
	if session.cisFirewallRulesErr != nil {
		session.cisFirewallRulesErr =
			fmt.Errorf("[ERROR] Error occured while configuring CIS Firewall rules : %s",
				session.cisFirewallRulesErr)
	}
	if session.cisFirewallRulesClient != nil && session.cisFirewallRulesClient.Service != nil {
		session.cisFirewallRulesClient.Service.EnableRetries(c.RetryCount, c.RetryDelay)
		session.cisFirewallRulesClient.SetDefaultHeaders(gohttp.Header{
			"X-Original-User-Agent": {fmt.Sprintf("terraform-provider-ibm/%s", version.Version)},
		})
	}

	// IBM Network CIS Authenticated Origin Pull
	cisOriginAuthOptions := &cisoriginpull.AuthenticatedOriginPullApiV1Options{
		URL:            cisEndPoint,
		Authenticator:  authenticator,
		Crn:            core.StringPtr(""),
		ZoneIdentifier: core.StringPtr(""),
	}

	session.cisOriginAuthClient, session.cisOriginAuthPullErr =
		cisoriginpull.NewAuthenticatedOriginPullApiV1(cisOriginAuthOptions)
	if session.cisOriginAuthPullErr != nil {
		session.cisOriginAuthPullErr = fmt.Errorf(
			"Error occured while configuring CIS Authenticated Origin Pullservice: %s",
			session.cisOriginAuthPullErr)
	}
	if session.cisOriginAuthClient != nil && session.cisOriginAuthClient.Service != nil {
		session.cisOriginAuthClient.Service.EnableRetries(c.RetryCount, c.RetryDelay)
		session.cisOriginAuthClient.SetDefaultHeaders(gohttp.Header{
			"X-Original-User-Agent": {fmt.Sprintf("terraform-provider-ibm/%s", version.Version)},
		})
	}

	// IAM IDENTITY Service
	// iamIdenityURL := fmt.Sprintf("https://%s.iam.cloud.ibm.com/v1", c.Region)
	iamIdenityURL := iamidentity.DefaultServiceURL
	if c.Visibility == "private" || c.Visibility == "public-and-private" {
		if c.Region == "us-south" || c.Region == "us-east" {
			iamIdenityURL = ContructEndpoint(fmt.Sprintf("private.%s.iam", c.Region), cloudEndpoint)
		} else {
			iamIdenityURL = ContructEndpoint("private.iam", cloudEndpoint)
		}
	}
	if fileMap != nil && c.Visibility != "public-and-private" {
		iamIdenityURL = fileFallBack(fileMap, c.Visibility, "IBMCLOUD_IAM_API_ENDPOINT", c.Region, iamIdenityURL)
	}
	iamIdentityOptions := &iamidentity.IamIdentityV1Options{
		Authenticator: authenticator,
		URL:           EnvFallBack([]string{"IBMCLOUD_IAM_API_ENDPOINT"}, iamIdenityURL),
	}
	iamIdentityClient, err := iamidentity.NewIamIdentityV1(iamIdentityOptions)
	if err != nil {
		session.iamIdentityErr = fmt.Errorf("[ERROR] Error occured while configuring IAM Identity service: %q", err)
	}
	if iamIdentityClient != nil && iamIdentityClient.Service != nil {
		iamIdentityClient.Service.EnableRetries(c.RetryCount, c.RetryDelay)
		iamIdentityClient.SetDefaultHeaders(gohttp.Header{
			"X-Original-User-Agent": {fmt.Sprintf("terraform-provider-ibm/%s", version.Version)},
		})
	}
	session.iamIdentityAPI = iamIdentityClient

	// IAM POLICY MANAGEMENT Service
	iamPolicyManagementURL := iampolicymanagement.DefaultServiceURL
	if c.Visibility == "private" || c.Visibility == "public-and-private" {
		if c.Region == "us-south" || c.Region == "us-east" {
			iamPolicyManagementURL = ContructEndpoint(fmt.Sprintf("private.%s.iam", c.Region), cloudEndpoint)
		} else {
			iamPolicyManagementURL = ContructEndpoint("private.iam", cloudEndpoint)
		}
	}
	if fileMap != nil && c.Visibility != "public-and-private" {
		iamPolicyManagementURL = fileFallBack(fileMap, c.Visibility, "IBMCLOUD_IAM_API_ENDPOINT", c.Region, iamPolicyManagementURL)
	}
	iamPolicyManagementOptions := &iampolicymanagement.IamPolicyManagementV1Options{
		Authenticator: authenticator,
		URL:           EnvFallBack([]string{"IBMCLOUD_IAM_API_ENDPOINT"}, iamPolicyManagementURL),
	}
	iamPolicyManagementClient, err := iampolicymanagement.NewIamPolicyManagementV1(iamPolicyManagementOptions)
	if err != nil {
		session.iamPolicyManagementErr = fmt.Errorf("[ERROR] Error occured while configuring IAM Policy Management service: %q", err)
	}
	if iamPolicyManagementClient != nil && iamPolicyManagementClient.Service != nil {
		iamPolicyManagementClient.Service.EnableRetries(c.RetryCount, c.RetryDelay)
		iamPolicyManagementClient.SetDefaultHeaders(gohttp.Header{
			"X-Original-User-Agent": {fmt.Sprintf("terraform-provider-ibm/%s", version.Version)},
		})
	}
	session.iamPolicyManagementAPI = iamPolicyManagementClient

	// IAM ACCESS GROUP
	iamAccessGroupsURL := iamaccessgroups.DefaultServiceURL
	if c.Visibility == "private" || c.Visibility == "public-and-private" {
		if c.Region == "us-south" || c.Region == "us-east" {
			iamAccessGroupsURL = ContructEndpoint(fmt.Sprintf("private.%s.iam", c.Region), cloudEndpoint)
		} else {
			iamAccessGroupsURL = ContructEndpoint("private.iam", cloudEndpoint)
		}
	}
	if fileMap != nil && c.Visibility != "public-and-private" {
		iamAccessGroupsURL = fileFallBack(fileMap, c.Visibility, "IBMCLOUD_IAM_API_ENDPOINT", c.Region, iamAccessGroupsURL)
	}
	iamAccessGroupsOptions := &iamaccessgroups.IamAccessGroupsV2Options{
		Authenticator: authenticator,
		URL:           EnvFallBack([]string{"IBMCLOUD_IAM_API_ENDPOINT"}, iamAccessGroupsURL),
	}
	iamAccessGroupsClient, err := iamaccessgroups.NewIamAccessGroupsV2(iamAccessGroupsOptions)
	if err != nil {
		session.iamAccessGroupsErr = fmt.Errorf("[ERROR] Error occured while configuring IAM Access Group service: %q", err)
	}
	if iamAccessGroupsClient != nil && iamAccessGroupsClient.Service != nil {
		iamAccessGroupsClient.Service.EnableRetries(c.RetryCount, c.RetryDelay)
		iamAccessGroupsClient.SetDefaultHeaders(gohttp.Header{
			"X-Original-User-Agent": {fmt.Sprintf("terraform-provider-ibm/%s", version.Version)},
		})
	}
	session.iamAccessGroupsAPI = iamAccessGroupsClient

	// RESOURCE MANAGEMENT Service
	rmURL := resourcemanager.DefaultServiceURL
	if c.Visibility == "private" {
		if c.Region == "us-south" || c.Region == "us-east" {
			rmURL = ContructEndpoint(fmt.Sprintf("private.%s.resource-controller", c.Region), fmt.Sprintf("%s", cloudEndpoint))
		} else {
			fmt.Println("Private Endpint supports only us-south and us-east region specific endpoint")
			rmURL = ContructEndpoint("private.us-south.resource-controller", fmt.Sprintf("%s", cloudEndpoint))
		}
	}
	if c.Visibility == "public-and-private" {
		if c.Region == "us-south" || c.Region == "us-east" {
			rmURL = ContructEndpoint(fmt.Sprintf("private.%s.resource-controller", c.Region), fmt.Sprintf("%s", cloudEndpoint))
		} else {
			rmURL = resourcemanager.DefaultServiceURL
		}
	}
	if fileMap != nil && c.Visibility != "public-and-private" {
		rmURL = fileFallBack(fileMap, c.Visibility, "IBMCLOUD_RESOURCE_MANAGEMENT_API_ENDPOINT", c.Region, rmURL)
	}
	resourceManagerOptions := &resourcemanager.ResourceManagerV2Options{
		Authenticator: authenticator,
		URL:           EnvFallBack([]string{"IBMCLOUD_RESOURCE_MANAGEMENT_API_ENDPOINT"}, rmURL),
	}
	resourceManagerClient, err := resourcemanager.NewResourceManagerV2(resourceManagerOptions)
	if err != nil {
		session.resourceManagerErr = fmt.Errorf("[ERROR] Error occured while configuring Resource Manager service: %q", err)
	}
	if resourceManagerClient != nil && resourceManagerClient.Service != nil {
		resourceManagerClient.Service.EnableRetries(c.RetryCount, c.RetryDelay)
		resourceManagerClient.SetDefaultHeaders(gohttp.Header{
			"X-Original-User-Agent": {fmt.Sprintf("terraform-provider-ibm/%s", version.Version)},
		})
	}
	session.resourceManagerAPI = resourceManagerClient

	//CLOUD SHELL Service
	cloudShellUrl := ibmcloudshellv1.DefaultServiceURL
	if fileMap != nil && c.Visibility != "public-and-private" {
		cloudShellUrl = fileFallBack(fileMap, c.Visibility, "IBMCLOUD_CLOUD_SHELL_API_ENDPOINT", c.Region, cloudShellUrl)
	}
	ibmCloudShellClientOptions := &ibmcloudshellv1.IBMCloudShellV1Options{
		Authenticator: authenticator,
		URL:           EnvFallBack([]string{"IBMCLOUD_CLOUD_SHELL_API_ENDPOINT"}, cloudShellUrl),
	}
	session.ibmCloudShellClient, err = ibmcloudshellv1.NewIBMCloudShellV1(ibmCloudShellClientOptions)
	if err != nil {
		session.ibmCloudShellClientErr = fmt.Errorf("[ERROR] Error occurred while configuring IBM Cloud Shell service: %q", err)
	}
	if session.ibmCloudShellClient != nil && session.ibmCloudShellClient.Service != nil {
		session.ibmCloudShellClient.Service.EnableRetries(c.RetryCount, c.RetryDelay)
		session.ibmCloudShellClient.SetDefaultHeaders(gohttp.Header{
			"X-Original-User-Agent": {fmt.Sprintf("terraform-provider-ibm/%s", version.Version)},
		})
	}

	// ENTERPRISE Service
	enterpriseURL := enterprisemanagementv1.DefaultServiceURL
	if c.Visibility == "private" {
		if c.Region == "us-south" || c.Region == "us-east" || c.Region == "eu-fr" {
			enterpriseURL = ContructEndpoint(fmt.Sprintf("private.%s.enterprise", c.Region), fmt.Sprintf("%s/v1", cloudEndpoint))
		} else {
			fmt.Println("Private Endpint supports only us-south and us-east region specific endpoint")
			enterpriseURL = ContructEndpoint("private.us-south.enterprise", fmt.Sprintf("%s/v1", cloudEndpoint))
		}
	}
	if c.Visibility == "public-and-private" {
		if c.Region == "us-south" || c.Region == "us-east" || c.Region == "eu-fr" {
			enterpriseURL = ContructEndpoint(fmt.Sprintf("private.%s.enterprise", c.Region),
				fmt.Sprintf("%s/v1", cloudEndpoint))
		} else {
			enterpriseURL = enterprisemanagementv1.DefaultServiceURL
		}
	}
	if fileMap != nil && c.Visibility != "public-and-private" {
		enterpriseURL = fileFallBack(fileMap, c.Visibility, "IBMCLOUD_ENTERPRISE_API_ENDPOINT", c.Region, enterpriseURL)
	}
	enterpriseManagementClientOptions := &enterprisemanagementv1.EnterpriseManagementV1Options{
		Authenticator: authenticator,
		URL:           EnvFallBack([]string{"IBMCLOUD_ENTERPRISE_API_ENDPOINT"}, enterpriseURL),
	}
	enterpriseManagementClient, err := enterprisemanagementv1.NewEnterpriseManagementV1(enterpriseManagementClientOptions)
	if err != nil {
		session.enterpriseManagementClientErr = fmt.Errorf("[ERROR] Error occurred while configuring IBM Cloud Enterprise Management API service: %q", err)
	}
	if enterpriseManagementClient != nil && enterpriseManagementClient.Service != nil {
		enterpriseManagementClient.Service.EnableRetries(c.RetryCount, c.RetryDelay)
		enterpriseManagementClient.SetDefaultHeaders(gohttp.Header{
			"X-Original-User-Agent": {fmt.Sprintf("terraform-provider-ibm/%s", version.Version)},
		})
	}
	session.enterpriseManagementClient = enterpriseManagementClient

	// RESOURCE CONTROLLER Service
	rcURL := resourcecontroller.DefaultServiceURL
	if c.Visibility == "private" {
		if c.Region == "us-south" || c.Region == "us-east" {
			rcURL = ContructEndpoint(fmt.Sprintf("private.%s.resource-controller", c.Region), cloudEndpoint)
		} else {
			fmt.Println("Private Endpint supports only us-south and us-east region specific endpoint")
			rcURL = ContructEndpoint("private.us-south.resource-controller", cloudEndpoint)
		}
	}
	if c.Visibility == "public-and-private" {
		if c.Region == "us-south" || c.Region == "us-east" {
			rcURL = ContructEndpoint(fmt.Sprintf("private.%s.resource-controller", c.Region), cloudEndpoint)
		} else {
			rcURL = resourcecontroller.DefaultServiceURL
		}
	}
	if fileMap != nil && c.Visibility != "public-and-private" {
		rcURL = fileFallBack(fileMap, c.Visibility, "IBMCLOUD_RESOURCE_CONTROLLER_API_ENDPOINT", c.Region, rcURL)
	}
	resourceControllerOptions := &resourcecontroller.ResourceControllerV2Options{
		Authenticator: authenticator,
		URL:           EnvFallBack([]string{"IBMCLOUD_RESOURCE_CONTROLLER_API_ENDPOINT"}, rcURL),
	}
	resourceControllerClient, err := resourcecontroller.NewResourceControllerV2(resourceControllerOptions)
	if err != nil {
		session.resourceControllerErr = fmt.Errorf("[ERROR] Error occured while configuring Resource Controller service: %q", err)
	}
	if resourceControllerClient != nil && resourceControllerClient.Service != nil {
		resourceControllerClient.Service.EnableRetries(c.RetryCount, c.RetryDelay)
		resourceControllerClient.SetDefaultHeaders(gohttp.Header{
			"X-Original-User-Agent": {fmt.Sprintf("terraform-provider-ibm/%s", version.Version)},
		})
	}
	session.resourceControllerAPI = resourceControllerClient

	// SECRETS MANAGER Service
	secretsManagerClientOptions := &secretsmanagerv1.SecretsManagerV1Options{
		Authenticator: authenticator,
	}
	/// Construct the service client.
	session.secretsManagerClient, err = secretsmanagerv1.NewSecretsManagerV1(secretsManagerClientOptions)
	if err != nil {
		session.secretsManagerClientErr = fmt.Errorf("[ERROR] Error occurred while configuring IBM Cloud Secrets Manager API service: %q", err)
	}
	if session.secretsManagerClient != nil && session.secretsManagerClient.Service != nil {
		// Enable retries for API calls
		session.secretsManagerClient.Service.EnableRetries(c.RetryCount, c.RetryDelay)
		// Add custom header for analytics
		session.secretsManagerClient.SetDefaultHeaders(gohttp.Header{
			"X-Original-User-Agent": {fmt.Sprintf("terraform-provider-ibm/%s", version.Version)},
		})
	}

	// SATELLITE Service
	containerEndpoint := kubernetesserviceapiv1.DefaultServiceURL
	if c.Visibility == "private" || c.Visibility == "public-and-private" {
		containerEndpoint = ContructEndpoint(fmt.Sprintf("private.%s.containers", c.Region), fmt.Sprintf("%s/global", cloudEndpoint))
	}
	if fileMap != nil && c.Visibility != "public-and-private" {
		containerEndpoint = fileFallBack(fileMap, c.Visibility, "IBMCLOUD_SATELLITE_API_ENDPOINT", c.Region, containerEndpoint)
	}
	kubernetesServiceV1Options := &kubernetesserviceapiv1.KubernetesServiceApiV1Options{
		URL:           EnvFallBack([]string{"IBMCLOUD_SATELLITE_API_ENDPOINT"}, containerEndpoint),
		Authenticator: authenticator,
	}
	session.satelliteClient, err = kubernetesserviceapiv1.NewKubernetesServiceApiV1(kubernetesServiceV1Options)
	if err != nil {
		session.satelliteClientErr = fmt.Errorf("[ERROR] Error occured while configuring satellite client: %q", err)
	}

	// Enable retries for API calls
	if session.satelliteClient != nil && session.satelliteClient.Service != nil {
		session.satelliteClient.Service.EnableRetries(c.RetryCount, c.RetryDelay)
		session.satelliteClient.SetDefaultHeaders(gohttp.Header{
			"X-Original-User-Agent": {fmt.Sprintf("terraform-provider-ibm/%s", version.Version)},
		})
	}

	// SATELLITE LINK Service
	// Construct an "options" struct for creating the service client.
	satelliteLinkEndpoint := satellitelinkv1.DefaultServiceURL
	if c.Visibility == "private" || c.Visibility == "public-and-private" {
		satelliteLinkEndpoint = ContructEndpoint("private.api.link.satellite", cloudEndpoint)
	}
	if fileMap != nil && c.Visibility != "public-and-private" {
		satelliteLinkEndpoint = fileFallBack(fileMap, c.Visibility, "IBMCLOUD_SATELLITE_LINK_API_ENDPOINT", c.Region, satelliteLinkEndpoint)
	}
	satelliteLinkClientOptions := &satellitelinkv1.SatelliteLinkV1Options{
		URL:           EnvFallBack([]string{"IBMCLOUD_SATELLITE_LINK_API_ENDPOINT"}, satelliteLinkEndpoint),
		Authenticator: authenticator,
	}
	session.satelliteLinkClient, err = satellitelinkv1.NewSatelliteLinkV1(satelliteLinkClientOptions)
	if err != nil {
		session.satelliteLinkClientErr = fmt.Errorf("[ERROR] Error occurred while configuring Satellite Link service: %q", err)
	}
	if session.satelliteLinkClient != nil && session.satelliteLinkClient.Service != nil {
		// Enable retries for API calls
		session.satelliteLinkClient.Service.EnableRetries(c.RetryCount, c.RetryDelay)
		// Add custom header for analytics
		session.satelliteLinkClient.SetDefaultHeaders(gohttp.Header{
			"X-Original-User-Agent": {fmt.Sprintf("terraform-provider-ibm/%s", version.Version)},
		})
	}

	esSchemaRegistryV1Options := &schemaregistryv1.SchemaregistryV1Options{
		Authenticator: authenticator,
	}
	session.esSchemaRegistryClient, err = schemaregistryv1.NewSchemaregistryV1(esSchemaRegistryV1Options)
	if err != nil {
		session.esSchemaRegistryErr = fmt.Errorf("[ERROR] Error occured while configuring Event Streams schema registry: %q", err)
	}
	if session.esSchemaRegistryClient != nil && session.esSchemaRegistryClient.Service != nil {
		session.esSchemaRegistryClient.Service.EnableRetries(c.RetryCount, c.RetryDelay)
		session.esSchemaRegistryClient.SetDefaultHeaders(gohttp.Header{
			"X-Original-User-Agent": {fmt.Sprintf("terraform-provider-ibm/%s", version.Version)},
		})
	}

	// Governance Service
	var configServiceApiClientURL string
	if c.Visibility == "private" || c.Visibility == "public-and-private" {
		configServiceApiClientURL, err = configurationgovernancev1.GetServiceURLForRegion("private." + c.Region)
		if err != nil && c.Visibility == "public-and-private" {
			configServiceApiClientURL, err = configurationgovernancev1.GetServiceURLForRegion(c.Region)
		}
	} else {
		configServiceApiClientURL, err = configurationgovernancev1.GetServiceURLForRegion(c.Region)
	}
	if err != nil {
		configServiceApiClientURL = configurationgovernancev1.DefaultServiceURL
	}
	configServiceApiClientOptions := &configurationgovernancev1.ConfigurationGovernanceV1Options{
		Authenticator: authenticator,
		URL:           EnvFallBack([]string{"IBMCLOUD_CONFIGURATION_GOVERNANCE_API_ENDPOINT"}, configServiceApiClientURL),
	}
	session.configServiceApiClient, err = configurationgovernancev1.NewConfigurationGovernanceV1(configServiceApiClientOptions)
	if err == nil {
		// Enable retries for API calls
		session.configServiceApiClient.Service.EnableRetries(c.RetryCount, c.RetryDelay)
		// Add custom header for analytics
		session.configServiceApiClient.SetDefaultHeaders(gohttp.Header{
			"X-Original-User-Agent": {fmt.Sprintf("terraform-provider-ibm/%s", version.Version)},
		})
	} else {
		session.configServiceApiClientErr = fmt.Errorf("Error occurred while configuring Config Service API service: %q", err)
	}

	//COMPLIANCE Service
	// Construct an "options" struct for creating the service client.
	var postureManagementClientURL string
	if c.Visibility == "public" || c.Visibility == "public-and-private" {
		postureManagementClientURL, err = posturemanagementv1.GetServiceURLForRegion(c.Region)
	} else {
		session.postureManagementClientErr = fmt.Errorf("[ERROR] Error occurred while configuring Security Insights Findings API service: `%v` visibility not supported", c.Visibility)
	}
	if err != nil {
		postureManagementClientURL = posturemanagementv1.DefaultServiceURL
	}
	if fileMap != nil && c.Visibility != "public-and-private" {
		postureManagementClientURL = fileFallBack(fileMap, c.Visibility, "IBMCLOUD_COMPLIANCE_API_ENDPOINT", c.Region, postureManagementClientURL)
	}
	postureManagementClientOptions := &posturemanagementv1.PostureManagementV1Options{
		Authenticator: authenticator,
		URL:           EnvFallBack([]string{"IBMCLOUD_COMPLIANCE_API_ENDPOINT"}, postureManagementClientURL),
		AccountID:     core.StringPtr(userConfig.UserAccount),
	}

	// Construct the service client.
	session.postureManagementClient, err = posturemanagementv1.NewPostureManagementV1(postureManagementClientOptions)
	if err != nil {
		session.postureManagementClientErr = fmt.Errorf("[ERROR] Error occurred while configuring Posture Management service: %q", err)
	}
	if session.postureManagementClient != nil && session.postureManagementClient.Service != nil {
		// Enable retries for API calls
		session.postureManagementClient.Service.EnableRetries(c.RetryCount, c.RetryDelay)
		// Add custom header for analytics
		session.postureManagementClient.SetDefaultHeaders(gohttp.Header{
			"X-Original-User-Agent": {fmt.Sprintf("terraform-provider-ibm/%s", version.Version)},
		})
	}

	//COMPLIANCE Service v2 version
	// Construct an "options" struct for creating the service client.
	var postureManagementClientURLv2 string
	if c.Visibility == "public" || c.Visibility == "public-and-private" {
		postureManagementClientURLv2, err = posturemanagementv2.GetServiceURLForRegion(c.Region)
	} else {
		session.postureManagementClientErrv2 = fmt.Errorf("[ERROR] Error occurred while configuring Security Compliance Centre API service: `%v` visibility not supported", c.Visibility)
	}
	if err != nil {
		session.postureManagementClientErrv2 = fmt.Errorf("[ERROR] Error occurred while configuring Security Posture Management API service:  `%s` region not supported", c.Region)
	}
	if fileMap != nil && c.Visibility != "public-and-private" {
		postureManagementClientURLv2 = fileFallBack(fileMap, c.Visibility, "IBMCLOUD_COMPLIANCE_API_ENDPOINT", c.Region, postureManagementClientURLv2)
	}
	postureManagementClientOptionsv2 := &posturemanagementv2.PostureManagementV2Options{
		Authenticator: authenticator,
		URL:           EnvFallBack([]string{"IBMCLOUD_COMPLIANCE_API_ENDPOINT"}, postureManagementClientURLv2),
	}

	// Construct the service client.
	session.postureManagementClientv2, err = posturemanagementv2.NewPostureManagementV2(postureManagementClientOptionsv2)
	if err != nil {
		session.postureManagementClientErrv2 = fmt.Errorf("[ERROR] Error occurred while configuring Posture Management v2 service: %q", err)
	}
	if session.postureManagementClientv2 != nil && session.postureManagementClientv2.Service != nil {
		// Enable retries for API calls
		session.postureManagementClientv2.Service.EnableRetries(c.RetryCount, c.RetryDelay)
		// Add custom header for analytics
		session.postureManagementClientv2.SetDefaultHeaders(gohttp.Header{
			"X-Original-User-Agent": {fmt.Sprintf("terraform-provider-ibm/%s", version.Version)},
		})
	}

	// Construct an "options" struct for creating the service client.
	var cdToolchainClientURL string
	if c.Visibility == "private" || c.Visibility == "public-and-private" {
		cdToolchainClientURL, err = cdtoolchainv2.GetServiceURLForRegion("private." + c.Region)
		if err != nil && c.Visibility == "public-and-private" {
			cdToolchainClientURL, err = cdtoolchainv2.GetServiceURLForRegion(c.Region)
		}
	} else {
		cdToolchainClientURL, err = cdtoolchainv2.GetServiceURLForRegion(c.Region)
	}
	if err != nil {
		cdToolchainClientURL = cdtoolchainv2.DefaultServiceURL
	}
	if fileMap != nil && c.Visibility != "public-and-private" {
		cdToolchainClientURL = fileFallBack(fileMap, c.Visibility, "IBMCLOUD_TOOLCHAIN_ENDPOINT", c.Region, cdToolchainClientURL)
	}
	cdToolchainClientOptions := &cdtoolchainv2.CdToolchainV2Options{
		Authenticator: authenticator,
		URL:           EnvFallBack([]string{"IBMCLOUD_TOOLCHAIN_ENDPOINT"}, cdToolchainClientURL),
	}

	// Construct the service client.
	session.cdToolchainClient, err = cdtoolchainv2.NewCdToolchainV2(cdToolchainClientOptions)
	if err == nil {
		// Enable retries for API calls
		session.cdToolchainClient.Service.EnableRetries(c.RetryCount, c.RetryDelay)
		// Add custom header for analytics
		session.cdToolchainClient.SetDefaultHeaders(gohttp.Header{
			"X-Original-User-Agent": {fmt.Sprintf("terraform-provider-ibm/%s", version.Version)},
		})
	} else {
		session.cdToolchainClientErr = fmt.Errorf("Error occurred while configuring Toolchain service: %q", err)
	}

	// Construct an "options" struct for creating the tekton pipeline service client.
	var cdTektonPipelineClientURL string
	if c.Visibility == "private" || c.Visibility == "public-and-private" {
		cdTektonPipelineClientURL, err = cdtektonpipelinev2.GetServiceURLForRegion("private." + c.Region)
		if err != nil && c.Visibility == "public-and-private" {
			cdTektonPipelineClientURL, err = cdtektonpipelinev2.GetServiceURLForRegion(c.Region)
		}
	} else {
		cdTektonPipelineClientURL, err = cdtektonpipelinev2.GetServiceURLForRegion(c.Region)
	}
	if err != nil {
		cdTektonPipelineClientURL = cdtektonpipelinev2.DefaultServiceURL
	}
	if fileMap != nil && c.Visibility != "public-and-private" {
		cdTektonPipelineClientURL = fileFallBack(fileMap, c.Visibility, "IBMCLOUD_TEKTON_PIPELINE_ENDPOINT", c.Region, cdTektonPipelineClientURL)
	}
	cdTektonPipelineClientOptions := &cdtektonpipelinev2.CdTektonPipelineV2Options{
		Authenticator: authenticator,
		URL:           EnvFallBack([]string{"IBMCLOUD_TEKTON_PIPELINE_ENDPOINT"}, cdTektonPipelineClientURL),
	}
	// Construct the service client.
	session.cdTektonPipelineClient, err = cdtektonpipelinev2.NewCdTektonPipelineV2(cdTektonPipelineClientOptions)
	if err == nil {
		// Enable retries for API calls
		session.cdTektonPipelineClient.Service.EnableRetries(c.RetryCount, c.RetryDelay)
		// Add custom header for analytics
		session.cdTektonPipelineClient.SetDefaultHeaders(gohttp.Header{
			"X-Original-User-Agent": {fmt.Sprintf("terraform-provider-ibm/%s", version.Version)},
		})
	} else {
		session.cdTektonPipelineClientErr = fmt.Errorf("Error occurred while configuring CD Tekton Pipeline service: %q", err)
	}

	if os.Getenv("TF_LOG") != "" {
		logDestination := log.Writer()
		goLogger := log.New(logDestination, "", log.LstdFlags)
		core.SetLogger(core.NewLogger(core.LevelDebug, goLogger, goLogger))
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
		softlayerSession.IAMRefreshToken = c.IAMRefreshToken
	}
	if c.SoftLayerAPIKey != "" && c.SoftLayerUserName != "" {
		log.Println("Configuring SoftLayer Session with API key")
		softlayerSession.APIKey = c.SoftLayerAPIKey
		softlayerSession.UserName = c.SoftLayerUserName
	}
	softlayerSession.AppendUserAgent(fmt.Sprintf("terraform-provider-ibm/%s", version.Version))
	ibmSession.SoftLayerSession = softlayerSession

	if c.IAMTrustedProfileID == "" && (c.IAMToken != "" && c.IAMRefreshToken == "") || (c.IAMToken == "" && c.IAMRefreshToken != "") {
		return nil, fmt.Errorf("iam_token and iam_refresh_token must be provided")
	}
	if c.IAMTrustedProfileID != "" && c.IAMToken == "" {
		return nil, fmt.Errorf("iam_token and iam_profile_id must be provided")
	}

	if c.IAMToken != "" {
		log.Println("Configuring IBM Cloud Session with token")
		var sess *bxsession.Session
		bmxConfig := &bluemix.Config{
			IAMAccessToken:  c.IAMToken,
			IAMRefreshToken: c.IAMRefreshToken,
			//Comment out debug mode for v0.12
			Debug:         os.Getenv("TF_LOG") != "",
			HTTPTimeout:   c.BluemixTimeout,
			Region:        c.Region,
			ResourceGroup: c.ResourceGroup,
			RetryDelay:    &c.RetryDelay,
			MaxRetries:    &c.RetryCount,
			Visibility:    c.Visibility,
			EndpointsFile: c.EndpointsFile,
			UserAgent:     fmt.Sprintf("terraform-provider-ibm/%s", version.Version),
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
			Debug:         os.Getenv("TF_LOG") != "",
			HTTPTimeout:   c.BluemixTimeout,
			Region:        c.Region,
			ResourceGroup: c.ResourceGroup,
			RetryDelay:    &c.RetryDelay,
			MaxRetries:    &c.RetryCount,
			Visibility:    c.Visibility,
			EndpointsFile: c.EndpointsFile,
			UserAgent:     fmt.Sprintf("terraform-provider-ibm/%s", version.Version),
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
			"User-Agent":            []string{http.UserAgent()},
			"X-Original-User-Agent": []string{config.UserAgent},
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
			"User-Agent":            []string{http.UserAgent()},
			"X-Original-User-Agent": []string{http.UserAgent()},
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
		user.UserEmail = email.(string)
	}
	user.UserID = claims["id"].(string)
	user.UserAccount = claims["account"].(map[string]interface{})["bss"].(string)
	iss := claims["iss"].(string)
	if strings.Contains(iss, "https://iam.cloud.ibm.com") {
		user.CloudName = "bluemix"
	} else {
		user.CloudName = "staging"
	}
	user.cloudType = "public"

	user.generation = 2
	return &user, nil
}

func RefreshToken(sess *bxsession.Session) error {
	config := sess.Config
	tokenRefresher, err := authentication.NewIAMAuthRepository(config, &rest.Client{
		DefaultHeader: gohttp.Header{
			"User-Agent":            []string{http.UserAgent()},
			"X-Original-User-Agent": []string{config.UserAgent},
		},
	})
	if err != nil {
		return err
	}
	_, err = tokenRefresher.RefreshToken()
	return err
}

func EnvFallBack(envs []string, defaultValue string) string {
	for _, k := range envs {
		if v := os.Getenv(k); v != "" {
			return v
		}
	}
	return defaultValue
}
func fileFallBack(fileMap map[string]interface{}, visibility, key, region, defaultValue string) string {
	if val, ok := fileMap[key]; ok {
		if v, ok := val.(map[string]interface{})[visibility]; ok {
			if r, ok := v.(map[string]interface{})[region]; ok && r.(string) != "" {
				return r.(string)
			}
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

func ContructEndpoint(subdomain, domain string) string {
	endpoint := fmt.Sprintf("https://%s.%s", subdomain, domain)
	return endpoint
}
