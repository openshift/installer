package connectivity

// Region represents ECS region
type Region string

// Constants of region definition
const (
	Hangzhou    = Region("cn-hangzhou")
	Qingdao     = Region("cn-qingdao")
	Beijing     = Region("cn-beijing")
	Hongkong    = Region("cn-hongkong")
	Shenzhen    = Region("cn-shenzhen")
	Shanghai    = Region("cn-shanghai")
	Zhangjiakou = Region("cn-zhangjiakou")
	Huhehaote   = Region("cn-huhehaote")
	ChengDu     = Region("cn-chengdu")
	HeYuan      = Region("cn-heyuan")
	WuLanChaBu  = Region("cn-wulanchabu")
	GuangZhou   = Region("cn-guangzhou")

	APSouthEast1 = Region("ap-southeast-1")
	APNorthEast1 = Region("ap-northeast-1")
	APSouthEast2 = Region("ap-southeast-2")
	APSouthEast3 = Region("ap-southeast-3")
	APSouthEast5 = Region("ap-southeast-5")

	APSouth1 = Region("ap-south-1")

	USWest1 = Region("us-west-1")
	USEast1 = Region("us-east-1")

	MEEast1 = Region("me-east-1")

	EUCentral1 = Region("eu-central-1")
	EUWest1    = Region("eu-west-1")

	ShenZhenFinance     = Region("cn-shenzhen-finance-1")
	ShanghaiFinance     = Region("cn-shanghai-finance-1")
	ShanghaiFinance1Pub = Region("cn-shanghai-finance-1-pub")
	CnNorth2Gov1        = Region("cn-north-2-gov-1")
)

var ValidRegions = []Region{
	Hangzhou, Qingdao, Beijing, Shenzhen, Hongkong, Shanghai, Zhangjiakou, Huhehaote, ChengDu, HeYuan, WuLanChaBu, GuangZhou,
	USWest1, USEast1,
	APNorthEast1, APSouthEast1, APSouthEast2, APSouthEast3, APSouthEast5,
	APSouth1,
	MEEast1,
	EUCentral1, EUWest1,
	ShenZhenFinance, ShanghaiFinance, CnNorth2Gov1, ShanghaiFinance1Pub,
}

var EcsClassicSupportedRegions = []Region{Shenzhen, Shanghai, Beijing, Qingdao, Hangzhou, Hongkong, USWest1, APSouthEast1}
var EcsSpotNoSupportedRegions = []Region{APSouth1}
var SlbGuaranteedSupportedRegions = []Region{Qingdao, Beijing, Hangzhou, Shanghai, Shenzhen, Zhangjiakou, Huhehaote, APSouthEast1, USEast1}
var DrdsSupportedRegions = []Region{Beijing, Shenzhen, Hangzhou, Qingdao, Hongkong, Shanghai, Huhehaote, Zhangjiakou, APSouthEast1}
var DrdsClassicNoSupportedRegions = []Region{Hongkong}
var GpdbSupportedRegions = []Region{Beijing, Shenzhen, Hangzhou, Shanghai, Hongkong}

// Some Ram resources only one can be owned by one account at the same time,
// skipped here to avoid multi regions concurrency conflict.
var RamNoSkipRegions = []Region{Hangzhou, EUCentral1, APSouth1}
var CenNoSkipRegions = []Region{Shanghai, EUCentral1, APSouth1}
var KmsSkippedRegions = []Region{Beijing, Shanghai}

// Actiontrail only one can be owned by one account at the same time,
// skipped here to avoid multi regions concurrency conflict.
var ActiontrailNoSkipRegions = []Region{Hangzhou, EUCentral1, APSouth1}
var FcNoSupportedRegions = []Region{MEEast1}
var DatahubSupportedRegions = []Region{Beijing, Hangzhou, Shanghai, Shenzhen, APSouthEast1}
var RdsClassicNoSupportedRegions = []Region{APSouth1, APSouthEast2, APSouthEast3, APNorthEast1, EUCentral1, EUWest1, MEEast1}
var RdsMultiAzNoSupportedRegions = []Region{Qingdao, APNorthEast1, APSouthEast5, MEEast1}
var RdsPPASNoSupportedRegions = []Region{Qingdao, USEast1, APNorthEast1, EUCentral1, MEEast1, APSouthEast2, APSouthEast3, APSouth1, APSouthEast5, ChengDu, EUWest1}
var RouteTableNoSupportedRegions = []Region{Beijing, Hangzhou, Shenzhen}
var ApiGatewayNoSupportedRegions = []Region{Zhangjiakou, Huhehaote, USEast1, USWest1, EUWest1, MEEast1}
var OtsHighPerformanceNoSupportedRegions = []Region{Qingdao, Zhangjiakou, Huhehaote, Hongkong, APSouthEast1, APSouthEast2, APSouthEast5, APNorthEast1, EUCentral1, MEEast1, APSouth1}
var OtsCapacityNoSupportedRegions = []Region{APSouthEast1, USWest1, USEast1}
var PrivateIpNoSupportedRegions = []Region{Beijing, Hangzhou, Shenzhen}
var SwarmSupportedRegions = []Region{Qingdao, Beijing, Zhangjiakou, Huhehaote, Hangzhou, Shanghai, Shenzhen, Hongkong, APNorthEast1, APSouthEast1, APSouthEast2,
	APSouthEast3, USWest1, USEast1, EUCentral1}
var ManagedKubernetesSupportedRegions = []Region{Beijing, Hangzhou, Shanghai, Shenzhen, ChengDu, Hongkong, APSouthEast1, APSouthEast2, USWest1}
var ServerlessKubernetesSupportedRegions = []Region{Beijing, Hangzhou, Shanghai, APSouthEast1, APSouthEast3, APSouthEast5, APSouth1, Huhehaote}
var KubernetesSupportedRegions = []Region{Beijing, Zhangjiakou, Huhehaote, Hangzhou, Shanghai, Shenzhen, Hongkong, APNorthEast1, APSouthEast1,
	APSouthEast2, APSouthEast3, APSouthEast5, APSouth1, USEast1, USWest1, EUWest1, MEEast1, EUCentral1}
var NasClassicSupportedRegions = []Region{Hangzhou, Qingdao, Beijing, Hongkong, Shenzhen, Shanghai, Zhangjiakou, Huhehaote, ShenZhenFinance, ShanghaiFinance}
var CasClassicSupportedRegions = []Region{Hangzhou, APSouth1, MEEast1, EUCentral1, APNorthEast1, APSouthEast2}
var CRNoSupportedRegions = []Region{Beijing, Hangzhou, Qingdao, Huhehaote, Zhangjiakou}
var MongoDBClassicNoSupportedRegions = []Region{Huhehaote, Zhangjiakou, APSouthEast2, APSouthEast3, APSouthEast5, APSouth1, USEast1, USWest1, APNorthEast1}
var MongoDBMultiAzSupportedRegions = []Region{Hangzhou, Beijing, Shenzhen, EUCentral1}
var DdoscooSupportedRegions = []Region{Hangzhou, APSouthEast1}
var DdosbgpSupportedRegions = []Region{Hangzhou, Beijing, Shenzhen, Qingdao, Shanghai, Zhangjiakou, Huhehaote}

//var NetworkAclSupportedRegions = []Region{Hangzhou, Beijing, Shanghai, Hongkong, APSouthEast5, APSouth1}
var EssScalingConfigurationMultiSgSupportedRegions = []Region{APSouthEast1, APSouth1}
var SlbClassicNoSupportedRegions = []Region{APNorthEast1, APSouthEast2, APSouthEast3, APSouthEast5, APSouth1, USEast1, MEEast1, EUCentral1, EUWest1, Huhehaote, Zhangjiakou}
var NasNoSupportedRegions = []Region{Qingdao, APSouth1, APSouthEast3, APSouthEast5}
var OssVersioningSupportedRegions = []Region{APSouth1}
var OssSseSupportedRegions = []Region{Qingdao, Hangzhou, Beijing, Shanghai, Shenzhen, Hongkong, APNorthEast1, APSouth1, USEast1}
var GpdbClassicNoSupportedRegions = []Region{APSouthEast2, APSouthEast3, APSouthEast5, APSouth1, USEast1, USWest1, APNorthEast1, EUCentral1}
var OnsNoSupportRegions = []Region{APSouthEast5}
var AlikafkaSupportedRegions = []Region{Hangzhou, Qingdao, Beijing, Hongkong, Shenzhen, Shanghai, Zhangjiakou, Huhehaote, ChengDu, HeYuan, APNorthEast1, APSouthEast1, APSouthEast3, EUCentral1, EUWest1, USEast1, USWest1}
var SmartagSupportedRegions = []Region{Shanghai, ShanghaiFinance, Hongkong, APSouthEast1, APSouthEast2, APSouthEast3, APSouthEast5, EUCentral1, APNorthEast1}
var YundunDbauditSupportedRegions = []Region{Hangzhou, Beijing, Shanghai}
var HttpHttpsHealthCheckMehtodSupportedRegions = []Region{Hangzhou, Beijing, Shanghai, EUWest1, ChengDu, Qingdao, Hongkong, Shenzhen, APSouthEast5, Zhangjiakou, Huhehaote, MEEast1, APSouth1, EUCentral1, USWest1, APSouthEast3, APSouthEast2, APSouthEast1, APNorthEast1}
var HBaseClassicSupportedRegions = []Region{Hangzhou, Beijing, Shanghai, Shenzhen}
var EdasSupportedRegions = []Region{Hangzhou, Beijing, Shanghai, Shenzhen, Zhangjiakou, Qingdao, Hongkong}
var CloudConfigSupportedRegions = []Region{Shanghai, APSouthEast1}
var DBReadwriteSplittingConnectionSupportedRegions = []Region{APSouthEast1}
var KVstoreClassicNetworkInstanceSupportRegions = []Region{Hangzhou, Beijing, Shanghai, APSouthEast1, USEast1, USWest1}
var MaxComputeSupportRegions = []Region{}
var FnfSupportRegions = []Region{Hangzhou, Beijing, Shanghai, Shenzhen, USWest1}
var PrivateLinkRegions = []Region{EUCentral1}
var BrainIndustrialRegions = []Region{Hangzhou}
var EciContainerGroupRegions = []Region{Beijing}
var TsdbInstanceSupportRegions = []Region{Beijing, Hangzhou, Shenzhen, Shanghai, ShenZhenFinance, Qingdao, Zhangjiakou, ShanghaiFinance, Hongkong, USWest1, APNorthEast1, EUWest1, APSouthEast1, APSouthEast2, APSouthEast3, EUCentral1, APSouthEast5, Zhangjiakou, CnNorth2Gov1}
var EipanycastSupportRegions = []Region{Hangzhou}
var VpcIpv6SupportRegions = []Region{Hangzhou, Shanghai, Shenzhen, Beijing, Huhehaote, Hongkong, APSouthEast1}
var EssdSupportRegions = []Region{Zhangjiakou, Huhehaote}
var AdbReserverUnSupportRegions = []Region{EUCentral1}
var KmsKeyUnSupportRegions = []Region{EUCentral1}
var DmSupportRegions = []Region{Hangzhou, APSouthEast1, APSouthEast2}
var BssOpenApiSupportRegions = []Region{Hangzhou, Shanghai, APSouthEast1}
var EipAddressBGPProSupportRegions = []Region{Hongkong}
var CenTransitRouterVpcAttachmentSupportRegions = []Region{Hangzhou}
var ARMSSupportRegions = []Region{Hangzhou, Shanghai, Beijing, APSouthEast1}
var SaeSupportRegions = []Region{Hangzhou, Shanghai, Beijing, Zhangjiakou, Shenzhen, USWest1}
var EcdSupportRegions = []Region{Hangzhou, Shanghai, Beijing, Shenzhen, Hongkong, APSouthEast1, APSouthEast2}
var EcpSupportRegions = []Region{Hangzhou, Shanghai, Beijing, Shenzhen}
var CsgSupportRegions = []Region{Shanghai}
var SddpSupportRegions = []Region{Hangzhou, Zhangjiakou, APSouthEast1}
var CddcSupportRegions = []Region{Shenzhen, Beijing, APSouth1, EUWest1, APNorthEast1, MEEast1, ChengDu, Qingdao, Shanghai, Hongkong, HeYuan, APSouthEast1, APSouthEast2, APSouthEast3, EUCentral1, Huhehaote, APSouthEast5, USEast1, Zhangjiakou, USWest1, Hangzhou}
