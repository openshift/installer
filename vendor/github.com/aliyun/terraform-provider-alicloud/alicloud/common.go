package alicloud

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/user"
	"path/filepath"
	"reflect"
	"runtime"
	"strconv"
	"strings"
	"time"

	"github.com/denverdino/aliyungo/cs"

	"github.com/aliyun/aliyun-datahub-sdk-go/datahub"
	sls "github.com/aliyun/aliyun-log-go-sdk"
	"github.com/aliyun/aliyun-oss-go-sdk/oss"
	"github.com/aliyun/aliyun-tablestore-go-sdk/tablestore"
	"github.com/aliyun/fc-go-sdk"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"

	"gopkg.in/yaml.v2"

	"math"

	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/requests"
	"github.com/denverdino/aliyungo/common"
	"github.com/google/uuid"
	"github.com/mitchellh/go-homedir"
)

type InstanceNetWork string

const (
	ClassicNet = InstanceNetWork("classic")
	VpcNet     = InstanceNetWork("vpc")
)

type PayType string

const (
	PrePaid  = PayType("PrePaid")
	PostPaid = PayType("PostPaid")
	Prepaid  = PayType("Prepaid")
	Postpaid = PayType("Postpaid")
)

const (
	NormalMode = "normal"
	SafetyMode = "safety"
)

type DdosbgpInsatnceType string

const (
	Enterprise   = DdosbgpInsatnceType("Enterprise")
	Professional = DdosbgpInsatnceType("Professional")
)

type DdosbgpInstanceIpType string

const (
	IPv4 = DdosbgpInstanceIpType("IPv4")
	IPv6 = DdosbgpInstanceIpType("IPv6")
)

type NetType string

const (
	Internet = NetType("Internet")
	Intranet = NetType("Intranet")
)

type NetworkType string

const (
	Classic         = NetworkType("Classic")
	Vpc             = NetworkType("Vpc")
	ClassicInternet = NetworkType("classic_internet")
	ClassicIntranet = NetworkType("classic_intranet")
	PUBLIC          = NetworkType("PUBLIC")
	PRIVATE         = NetworkType("PRIVATE")
)

type NodeType string

const (
	WORKER = NodeType("WORKER")
	KIBANA = NodeType("KIBANA")
)

type ActionType string

const (
	OPEN  = ActionType("OPEN")
	CLOSE = ActionType("CLOSE")
)

type TimeType string

const (
	Hour  = TimeType("Hour")
	Day   = TimeType("Day")
	Week  = TimeType("Week")
	Month = TimeType("Month")
	Year  = TimeType("Year")
)

type IpVersion string

const (
	IPV4 = IpVersion("ipv4")
	IPV6 = IpVersion("ipv6")
)

type Status string

const (
	Pending     = Status("Pending")
	Creating    = Status("Creating")
	Running     = Status("Running")
	Available   = Status("Available")
	Unavailable = Status("Unavailable")
	Modifying   = Status("Modifying")
	Deleting    = Status("Deleting")
	Starting    = Status("Starting")
	Stopping    = Status("Stopping")
	Stopped     = Status("Stopped")
	Normal      = Status("Normal")
	Changing    = Status("Changing")
	Online      = Status("online")
	Configuring = Status("configuring")

	Associating   = Status("Associating")
	Unassociating = Status("Unassociating")
	InUse         = Status("InUse")
	DiskInUse     = Status("In_use")

	Active   = Status("Active")
	Inactive = Status("Inactive")
	Idle     = Status("Idle")

	SoldOut = Status("SoldOut")

	InService      = Status("InService")
	Removing       = Status("Removing")
	DisabledStatus = Status("Disabled")

	Init            = Status("Init")
	Provisioning    = Status("Provisioning")
	Updating        = Status("Updating")
	FinancialLocked = Status("FinancialLocked")

	PUBLISHED   = Status("Published")
	NOPUBLISHED = Status("NonPublished")

	Deleted = Status("Deleted")
	Null    = Status("Null")

	Enable = Status("Enable")
	BINDED = Status("BINDED")
)

type IPType string

const (
	Inner   = IPType("Inner")
	Private = IPType("Private")
	Public  = IPType("Public")
)

type ResourceType string

const (
	ResourceTypeInstance      = ResourceType("Instance")
	ResourceTypeDisk          = ResourceType("Disk")
	ResourceTypeVSwitch       = ResourceType("VSwitch")
	ResourceTypeRds           = ResourceType("Rds")
	ResourceTypePolarDB       = ResourceType("PolarDB")
	IoOptimized               = ResourceType("IoOptimized")
	ResourceTypeRkv           = ResourceType("KVStore")
	ResourceTypeFC            = ResourceType("FunctionCompute")
	ResourceTypeElasticsearch = ResourceType("Elasticsearch")
	ResourceTypeSlb           = ResourceType("Slb")
	ResourceTypeMongoDB       = ResourceType("MongoDB")
	ResourceTypeGpdb          = ResourceType("Gpdb")
	ResourceTypeHBase         = ResourceType("HBase")
	ResourceTypeAdb           = ResourceType("ADB")
	ResourceTypeCassandra     = ResourceType("Cassandra")
)

type InternetChargeType string

const (
	PayByBandwidth = InternetChargeType("PayByBandwidth")
	PayByTraffic   = InternetChargeType("PayByTraffic")
	PayBy95        = InternetChargeType("PayBy95")
)

type AccountSite string

const (
	DomesticSite = AccountSite("Domestic")
	IntlSite     = AccountSite("International")
)

const (
	SnapshotCreatingInProcessing = Status("progressing")
	SnapshotCreatingAccomplished = Status("accomplished")
	SnapshotCreatingFailed       = Status("failed")

	SnapshotPolicyCreating  = Status("Creating")
	SnapshotPolicyAvailable = Status("available")
	SnapshotPolicyNormal    = Status("Normal")
)

// timeout for common product, ecs e.g.
const DefaultTimeout = 120
const Timeout5Minute = 300
const DefaultTimeoutMedium = 500

// timeout for long time progerss product, rds e.g.
const DefaultLongTimeout = 1000

const DefaultIntervalMini = 2

const DefaultIntervalShort = 5

const DefaultIntervalMedium = 10

const DefaultIntervalLong = 20

const (
	PageSizeSmall  = 10
	PageSizeMedium = 20
	PageSizeLarge  = 50
	PageSizeXLarge = 100
)

// Protocol represents network protocol
type Protocol string

// Constants of protocol definition
const (
	Http  = Protocol("http")
	Https = Protocol("https")
	Tcp   = Protocol("tcp")
	Udp   = Protocol("udp")
	All   = Protocol("all")
	Icmp  = Protocol("icmp")
	Gre   = Protocol("gre")
)

// ValidProtocols network protocol list
var ValidProtocols = []Protocol{Http, Https, Tcp, Udp}

// simple array value check method, support string type only
func isProtocolValid(value string) bool {
	res := false
	for _, v := range ValidProtocols {
		if string(v) == value {
			res = true
		}
	}
	return res
}

// default region for all resource
const DEFAULT_REGION = "cn-beijing"

const INT_MAX = 2147483647

// symbol of multiIZ
const MULTI_IZ_SYMBOL = "MAZ"

const COMMA_SEPARATED = ","

const COLON_SEPARATED = ":"

const SLASH_SEPARATED = "/"

const LOCAL_HOST_IP = "127.0.0.1"

// Takes the result of flatmap.Expand for an array of strings
// and returns a []string
func expandStringList(configured []interface{}) []string {
	vs := make([]string, 0, len(configured))
	for _, v := range configured {
		if v == nil {
			continue
		}
		vs = append(vs, v.(string))
	}
	return vs
}

// Takes list of string to strings. Expand to an array
// of raw strings and returns a []interface{}
func convertListStringToListInterface(list []string) []interface{} {
	vs := make([]interface{}, 0, len(list))
	for _, v := range list {
		vs = append(vs, v)
	}
	return vs
}

func expandIntList(configured []interface{}) []int {
	vs := make([]int, 0, len(configured))
	for _, v := range configured {
		vs = append(vs, v.(int))
	}
	return vs
}

// Convert the result for an array and returns a Json string
func convertListToJsonString(configured []interface{}) string {
	if len(configured) < 1 {
		return ""
	}
	result := "["
	for i, v := range configured {
		if v == nil {
			continue
		}
		result += "\"" + v.(string) + "\""
		if i < len(configured)-1 {
			result += ","
		}
	}
	result += "]"
	return result
}

func convertJsonStringToStringList(src interface{}) (result []interface{}) {
	if err, ok := src.([]interface{}); !ok {
		panic(err)
	}
	for _, v := range src.([]interface{}) {
		result = append(result, fmt.Sprint(formatInt(v)))
	}
	return
}

// Convert the result for an array and returns a comma separate
func convertListToCommaSeparate(configured []interface{}) string {
	if len(configured) < 1 {
		return ""
	}
	result := ""
	for i, v := range configured {
		rail := ","
		if i == len(configured)-1 {
			rail = ""
		}
		result += v.(string) + rail
	}
	return result
}

func convertBoolToString(configured bool) string {
	return strconv.FormatBool(configured)
}

func convertIntergerToString(configured int) string {
	return strconv.Itoa(configured)
}

func convertFloat64ToString(configured float64) string {
	return strconv.FormatFloat(configured, 'E', -1, 64)
}

func convertJsonStringToList(configured string) ([]interface{}, error) {
	result := make([]interface{}, 0)
	if err := json.Unmarshal([]byte(configured), &result); err != nil {
		return nil, err
	}

	return result, nil
}

func convertMaptoJsonString(m map[string]interface{}) (string, error) {
	sm := make(map[string]string, len(m))
	for k, v := range m {
		sm[k] = v.(string)
	}

	if result, err := json.Marshal(sm); err != nil {
		return "", err
	} else {
		return string(result), nil
	}
}

func convertMapFloat64ToJsonString(m map[string]interface{}) (string, error) {
	sm := make(map[string]json.Number, len(m))

	for k, v := range m {
		sm[k] = v.(json.Number)
	}

	if result, err := json.Marshal(sm); err != nil {
		return "", err
	} else {
		return string(result), nil
	}
}

func StringPointer(s string) *string {
	return &s
}

func BoolPointer(b bool) *bool {
	return &b
}

func Int32Pointer(i int32) *int32 {
	return &i
}

func Int64Pointer(i int64) *int64 {
	return &i
}

func IntMin(x, y int) int {
	if x < y {
		return x
	}
	return y
}

const ServerSideEncryptionAes256 = "AES256"
const ServerSideEncryptionKMS = "KMS"

type OptimizedType string

const (
	IOOptimized   = OptimizedType("optimized")
	NoneOptimized = OptimizedType("none")
)

type TagResourceType string

const (
	TagResourceImage         = TagResourceType("image")
	TagResourceInstance      = TagResourceType("instance")
	TagResourceAcl           = TagResourceType("acl")
	TagResourceCertificate   = TagResourceType("certificate")
	TagResourceSnapshot      = TagResourceType("snapshot")
	TagResourceKeypair       = TagResourceType("keypair")
	TagResourceDisk          = TagResourceType("disk")
	TagResourceSecurityGroup = TagResourceType("securitygroup")
	TagResourceEni           = TagResourceType("eni")
	TagResourceCdn           = TagResourceType("DOMAIN")
	TagResourceVpc           = TagResourceType("VPC")
	TagResourceVSwitch       = TagResourceType("VSWITCH")
	TagResourceRouteTable    = TagResourceType("ROUTETABLE")
	TagResourceEip           = TagResourceType("EIP")
	TagResourcePlugin        = TagResourceType("plugin")
	TagResourceApiGroup      = TagResourceType("apiGroup")
	TagResourceApp           = TagResourceType("app")
	TagResourceTopic         = TagResourceType("topic")
	TagResourceConsumerGroup = TagResourceType("consumergroup")
	TagResourceCluster       = TagResourceType("cluster")
)

type KubernetesNodeType string

const (
	KubernetesNodeMaster = ResourceType("Master")
	KubernetesNodeWorker = ResourceType("Worker")
)

func getPagination(pageNumber, pageSize int) (pagination common.Pagination) {
	pagination.PageSize = pageSize
	pagination.PageNumber = pageNumber
	return
}

const CharityPageUrl = "http://promotion.alicdn.com/help/oss/error.html"

func userDataHashSum(user_data string) string {
	// Check whether the user_data is not Base64 encoded.
	// Always calculate hash of base64 decoded value since we
	// check against double-encoding when setting it
	v, base64DecodeError := base64.StdEncoding.DecodeString(user_data)
	if base64DecodeError != nil {
		v = []byte(user_data)
	}
	return string(v)
}

// Remove useless blank in the string.
func Trim(v string) string {
	if len(v) < 1 {
		return v
	}
	return strings.Trim(v, " ")
}

func ConvertIntegerToInt(value requests.Integer) (v int, err error) {
	if strings.TrimSpace(string(value)) == "" {
		return
	}
	v, err = strconv.Atoi(string(value))
	if err != nil {
		return v, fmt.Errorf("Converting integer %s to int got an error: %#v.", value, err)
	}
	return
}

func GetUserHomeDir() (string, error) {
	usr, err := user.Current()
	if err != nil {
		return "", fmt.Errorf("Get current user got an error: %#v.", err)
	}
	return usr.HomeDir, nil
}

func writeToFile(filePath string, data interface{}) error {
	var out string
	switch data.(type) {
	case string:
		out = data.(string)
		break
	case nil:
		return nil
	default:
		bs, err := json.MarshalIndent(data, "", "\t")
		if err != nil {
			return fmt.Errorf("MarshalIndent data %#v got an error: %#v", data, err)
		}
		out = string(bs)
	}

	if strings.HasPrefix(filePath, "~") {
		home, err := GetUserHomeDir()
		if err != nil {
			return err
		}
		if home != "" {
			filePath = strings.Replace(filePath, "~", home, 1)
		}
	}

	if _, err := os.Stat(filePath); err == nil {
		if err := os.Remove(filePath); err != nil {
			return err
		}
	}

	return ioutil.WriteFile(filePath, []byte(out), 422)
}

type Invoker struct {
	catchers []*Catcher
}

type Catcher struct {
	Reason           string
	RetryCount       int
	RetryWaitSeconds int
}

var ClientErrorCatcher = Catcher{AliyunGoClientFailure, 10, 5}
var ServiceBusyCatcher = Catcher{"ServiceUnavailable", 10, 5}
var ThrottlingCatcher = Catcher{Throttling, 50, 2}

func NewInvoker() Invoker {
	i := Invoker{}
	i.AddCatcher(ClientErrorCatcher)
	i.AddCatcher(ServiceBusyCatcher)
	i.AddCatcher(ThrottlingCatcher)
	return i
}

func (a *Invoker) AddCatcher(catcher Catcher) {
	a.catchers = append(a.catchers, &catcher)
}

func (a *Invoker) Run(f func() error) error {
	err := f()

	if err == nil {
		return nil
	}

	for _, catcher := range a.catchers {
		if IsExpectedErrors(err, []string{catcher.Reason}) {
			catcher.RetryCount--

			if catcher.RetryCount <= 0 {
				return fmt.Errorf("Retry timeout and got an error: %#v.", err)
			} else {
				time.Sleep(time.Duration(catcher.RetryWaitSeconds) * time.Second)
				return a.Run(f)
			}
		}
	}
	return err
}

func buildClientToken(action string) string {
	token := strings.TrimSpace(fmt.Sprintf("TF-%s-%d-%s", action, time.Now().Unix(), strings.Trim(uuid.New().String(), "-")))
	if len(token) > 64 {
		token = token[0:64]
	}
	return token
}

func getNextpageNumber(number requests.Integer) (requests.Integer, error) {
	page, err := strconv.Atoi(string(number))
	if err != nil {
		return "", err
	}
	return requests.NewInteger(page + 1), nil
}

func terraformToAPI(field string) string {
	var result string
	for _, v := range strings.Split(field, "_") {
		if len(v) > 0 {
			result = fmt.Sprintf("%s%s%s", result, strings.ToUpper(string(v[0])), v[1:])
		}
	}
	return result
}

func compareJsonTemplateAreEquivalent(tem1, tem2 string) (bool, error) {
	obj1 := make(map[string]interface{})
	err := json.Unmarshal([]byte(tem1), &obj1)
	if err != nil {
		return false, err
	}

	canonicalJson1, _ := json.Marshal(obj1)

	obj2 := make(map[string]interface{})
	err = json.Unmarshal([]byte(tem2), &obj2)
	if err != nil {
		return false, err
	}

	canonicalJson2, _ := json.Marshal(obj2)

	equal := bytes.Compare(canonicalJson1, canonicalJson2) == 0
	if !equal {
		log.Printf("[DEBUG] Canonical template are not equal.\nFirst: %s\nSecond: %s\n",
			canonicalJson1, canonicalJson2)
	}
	return equal, nil
}

func compareYamlTemplateAreEquivalent(tem1, tem2 string) (bool, error) {
	var obj1 interface{}
	err := yaml.Unmarshal([]byte(tem1), &obj1)
	if err != nil {
		return false, err
	}

	canonicalYaml1, _ := yaml.Marshal(obj1)

	var obj2 interface{}
	err = yaml.Unmarshal([]byte(tem2), &obj2)
	if err != nil {
		return false, err
	}

	canonicalYaml2, _ := yaml.Marshal(obj2)

	equal := bytes.Compare(canonicalYaml1, canonicalYaml2) == 0
	if !equal {
		log.Printf("[DEBUG] Canonical template are not equal.\nFirst: %s\nSecond: %s\n",
			canonicalYaml1, canonicalYaml2)
	}
	return equal, nil
}

// loadFileContent returns contents of a file in a given path
func loadFileContent(v string) ([]byte, error) {
	filename, err := homedir.Expand(v)
	if err != nil {
		return nil, err
	}
	fileContent, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	return fileContent, nil
}

func debugOn() bool {
	for _, part := range strings.Split(os.Getenv("DEBUG"), ",") {
		if strings.TrimSpace(part) == "terraform" {
			return true
		}
	}
	return false
}

func addDebug(action, content interface{}, requestInfo ...interface{}) {
	if debugOn() {
		trace := "[DEBUG TRACE]:\n"
		for skip := 1; skip < 5; skip++ {
			_, filepath, line, _ := runtime.Caller(skip)
			trace += fmt.Sprintf("%s:%d\n", filepath, line)
		}

		if len(requestInfo) > 0 {
			var request = struct {
				Domain     string
				Version    string
				UserAgent  string
				ActionName string
				Method     string
				Product    string
				Region     string
				AK         string
			}{}
			switch requestInfo[0].(type) {
			case *requests.RpcRequest:
				tmp := requestInfo[0].(*requests.RpcRequest)
				request.Domain = tmp.GetDomain()
				request.Version = tmp.GetVersion()
				request.ActionName = tmp.GetActionName()
				request.Method = tmp.GetMethod()
				request.Product = tmp.GetProduct()
				request.Region = tmp.GetRegionId()
			case *requests.RoaRequest:
				tmp := requestInfo[0].(*requests.RoaRequest)
				request.Domain = tmp.GetDomain()
				request.Version = tmp.GetVersion()
				request.ActionName = tmp.GetActionName()
				request.Method = tmp.GetMethod()
				request.Product = tmp.GetProduct()
				request.Region = tmp.GetRegionId()
			case *requests.CommonRequest:
				tmp := requestInfo[0].(*requests.CommonRequest)
				request.Domain = tmp.GetDomain()
				request.Version = tmp.GetVersion()
				request.ActionName = tmp.GetActionName()
				request.Method = tmp.GetMethod()
				request.Product = tmp.GetProduct()
				request.Region = tmp.GetRegionId()
			case *fc.Client:
				client := requestInfo[0].(*fc.Client)
				request.Version = client.Config.APIVersion
				request.Product = "FC"
				request.ActionName = fmt.Sprintf("%s", action)
			case *sls.Client:
				request.Product = "LOG"
				request.ActionName = fmt.Sprintf("%s", action)
			case *tablestore.TableStoreClient:
				request.Product = "OTS"
				request.ActionName = fmt.Sprintf("%s", action)
			case *oss.Client:
				request.Product = "OSS"
				request.ActionName = fmt.Sprintf("%s", action)
			case *datahub.DataHub:
				request.Product = "DataHub"
				request.ActionName = fmt.Sprintf("%s", action)
			case *cs.Client:
				request.Product = "CS"
				request.ActionName = fmt.Sprintf("%s", action)
			}

			requestContent := ""
			if len(requestInfo) > 1 {
				requestContent = fmt.Sprintf("%#v", requestInfo[1])
			}

			if len(requestInfo) == 1 {
				if v, ok := requestInfo[0].(map[string]interface{}); ok {
					if res, err := json.Marshal(&v); err == nil {
						requestContent = string(res)
					}
					if res, err := json.Marshal(&content); err == nil {
						content = string(res)
					}
				}
			}

			content = fmt.Sprintf("%vDomain:%v, Version:%v, ActionName:%v, Method:%v, Product:%v, Region:%v\n\n"+
				"*************** %s Request ***************\n%#v\n",
				content, request.Domain, request.Version, request.ActionName,
				request.Method, request.Product, request.Region, request.ActionName, requestContent)
		}

		//fmt.Printf(DefaultDebugMsg, action, content, trace)
		log.Printf(DefaultDebugMsg, action, content, trace)
	}
}

// Return a ComplexError which including extra error message, error occurred file and path
func GetFunc(level int) string {
	pc, _, _, ok := runtime.Caller(level)
	if !ok {
		log.Printf("[ERROR] runtime.Caller error in GetFuncName.")
		return ""
	}
	return strings.TrimPrefix(filepath.Ext(runtime.FuncForPC(pc).Name()), ".")
}

func ParseResourceId(id string, length int) (parts []string, err error) {
	parts = strings.Split(id, ":")

	if len(parts) != length {
		err = WrapError(fmt.Errorf("Invalid Resource Id %s. Expected parts' length %d, got %d", id, length, len(parts)))
	}
	return parts, err
}

func ParseSlbListenerId(id string) (parts []string, err error) {
	parts = strings.Split(id, ":")
	if len(parts) != 2 && len(parts) != 3 {
		err = WrapError(fmt.Errorf("Invalid alicloud_slb_listener Id %s. Expected Id format is <slb id>:<protocol>:< frontend>.", id))
	}
	return parts, err
}

func GetCenChildInstanceType(id string) (c string, e error) {
	if strings.HasPrefix(id, "vpc") {
		return ChildInstanceTypeVpc, nil
	} else if strings.HasPrefix(id, "vbr") {
		return ChildInstanceTypeVbr, nil
	} else if strings.HasPrefix(id, "ccn") {
		return ChildInstanceTypeCcn, nil
	} else {
		return c, fmt.Errorf("CEN child instance ID invalid. Now, it only supports VPC or VBR or CCN instance.")
	}
}

func BuildStateConf(pending, target []string, timeout, delay time.Duration, f resource.StateRefreshFunc) *resource.StateChangeConf {
	return &resource.StateChangeConf{
		Pending:    pending,
		Target:     target,
		Refresh:    f,
		Timeout:    timeout,
		Delay:      delay,
		MinTimeout: 3 * time.Second,
	}
}

func incrementalWait(firstDuration time.Duration, increaseDuration time.Duration) func() {
	retryCount := 1
	return func() {
		var waitTime time.Duration
		if retryCount == 1 {
			waitTime = firstDuration
		} else if retryCount > 1 {
			waitTime += increaseDuration
		}
		time.Sleep(waitTime)
		retryCount++
	}
}

// If auto renew, the period computed from computePeriodByUnit will be changed
// This method used to compute a period accourding to current period and unit
func computePeriodByUnit(createTime, endTime interface{}, currentPeriod int, periodUnit string) (int, error) {
	var createTimeStr, endTimeStr string
	switch value := createTime.(type) {
	case int64:
		createTimeStr = time.Unix(createTime.(int64), 0).Format(time.RFC3339)
		endTimeStr = time.Unix(endTime.(int64), 0).Format(time.RFC3339)
	case string:
		createTimeStr = createTime.(string)
		endTimeStr = endTime.(string)
	default:
		return 0, WrapError(fmt.Errorf("Unsupported time type: %#v", value))
	}
	// currently, there is time value does not format as standard RFC3339
	UnStandardRFC3339 := "2006-01-02T15:04Z07:00"
	create, err := time.Parse(time.RFC3339, createTimeStr)
	if err != nil {
		log.Printf("Parase the CreateTime %#v failed and error is: %#v.", createTime, err)
		create, err = time.Parse(UnStandardRFC3339, createTimeStr)
		if err != nil {
			return 0, WrapError(err)
		}
	}
	end, err := time.Parse(time.RFC3339, endTimeStr)
	if err != nil {
		log.Printf("Parase the EndTime %#v failed and error is: %#v.", endTime, err)
		end, err = time.Parse(UnStandardRFC3339, endTimeStr)
		if err != nil {
			return 0, WrapError(err)
		}
	}
	var period int
	switch periodUnit {
	case "Month":
		period = int(math.Floor(end.Sub(create).Hours() / 24 / 30))
	case "Week":
		period = int(math.Floor(end.Sub(create).Hours() / 24 / 7))
	case "Year":
		period = int(math.Floor(end.Sub(create).Hours() / 24 / 365))
	default:
		err = fmt.Errorf("Unexpected period unit %s", periodUnit)
	}
	// The period at least is 1
	if period < 1 {
		period = 1
	}
	if period > 12 {
		period = 12
	}
	// period can not be modified and if the new period is changed, using the previous one.
	if currentPeriod > 0 && currentPeriod != period {
		period = currentPeriod
	}
	return period, WrapError(err)
}

func checkWaitForReady(object interface{}, conditions map[string]interface{}) (bool, map[string]interface{}, error) {
	if conditions == nil {
		return false, nil, nil
	}
	objectType := reflect.TypeOf(object)
	objectValue := reflect.ValueOf(object)
	values := make(map[string]interface{})
	for key, value := range conditions {
		if _, ok := objectType.FieldByName(key); ok {
			current := objectValue.FieldByName(key)
			values[key] = current
			if fmt.Sprintf("%v", current) != fmt.Sprintf("%v", value) {
				return false, values, nil
			}
		} else {
			return false, values, WrapError(fmt.Errorf("There is missing attribute %s in the object.", key))
		}
	}
	return true, values, nil
}

// When  using teadsl, we need to convert float, int64 and int32 to int for comparison.
func formatInt(src interface{}) int {
	if src == nil {
		return 0
	}
	attrType := reflect.TypeOf(src)
	switch attrType.String() {
	case "float64":
		return int(src.(float64))
	case "float32":
		return int(src.(float32))
	case "int64":
		return int(src.(int64))
	case "int32":
		return int(src.(int32))
	case "int":
		return src.(int)
	case "string":
		v, err := strconv.Atoi(src.(string))
		if err != nil {
			panic(err)
		}
		return v
	case "json.Number":
		v, err := strconv.Atoi(src.(json.Number).String())
		if err != nil {
			panic(err)
		}
		return v
	default:
		panic(fmt.Sprintf("Not support type %s", attrType.String()))
	}
	return 0
}

func convertArrayObjectToJsonString(src interface{}) (string, error) {
	res, err := json.Marshal(&src)
	if err != nil {
		return "", err
	}
	return string(res), nil
}

func convertArrayToString(src interface{}, sep string) string {
	if src == nil {
		return ""
	}
	items := make([]string, 0)
	for _, v := range src.([]interface{}) {
		items = append(items, fmt.Sprint(v))
	}
	return strings.Join(items, sep)
}
