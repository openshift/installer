package endpoints

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"

	"github.com/IBM-Cloud/bluemix-go/bmxerror"
	"github.com/IBM-Cloud/bluemix-go/helpers"
)

//EndpointLocator ...
type EndpointLocator interface {
	AccountManagementEndpoint() (string, error)
	CertificateManagerEndpoint() (string, error)
	CFAPIEndpoint() (string, error)
	ContainerEndpoint() (string, error)
	ContainerRegistryEndpoint() (string, error)
	CisEndpoint() (string, error)
	GlobalSearchEndpoint() (string, error)
	GlobalTaggingEndpoint() (string, error)
	IAMEndpoint() (string, error)
	IAMPAPEndpoint() (string, error)
	ICDEndpoint() (string, error)
	MCCPAPIEndpoint() (string, error)
	ResourceManagementEndpoint() (string, error)
	ResourceControllerEndpoint() (string, error)
	ResourceCatalogEndpoint() (string, error)
	UAAEndpoint() (string, error)
	CseEndpoint() (string, error)
	SchematicsEndpoint() (string, error)
	UserManagementEndpoint() (string, error)
	HpcsEndpoint() (string, error)
	FunctionsEndpoint() (string, error)
	SatelliteEndpoint() (string, error)
}

const (
	//ErrCodeServiceEndpoint ...
	ErrCodeServiceEndpoint = "ServiceEndpointDoesnotExist"
)

var regionToEndpoint = map[string]map[string]string{
	"cf": {
		"us-south": "https://api.ng.bluemix.net",
		"us-east":  "https://api.us-east.bluemix.net",
		"eu-gb":    "https://api.eu-gb.bluemix.net",
		"au-syd":   "https://api.au-syd.bluemix.net",
		"eu-de":    "https://api.eu-de.bluemix.net",
		"jp-tok":   "https://api.jp-tok.bluemix.net",
	},
	"cr": {
		"us-south": "us.icr.io",
		"us-east":  "us.icr.io",
		"eu-de":    "de.icr.io",
		"au-syd":   "au.icr.io",
		"eu-gb":    "uk.icr.io",
		"jp-tok":   "jp.icr.io",
		"jp-osa":   "jp2.icr.io",
	},
	"uaa": {
		"us-south": "https://iam.cloud.ibm.com/cloudfoundry/login/us-south",
		"us-east":  "https://iam.cloud.ibm.com/cloudfoundry/login/us-east",
		"eu-gb":    "https://iam.cloud.ibm.com/cloudfoundry/login/uk-south",
		"au-syd":   "https://iam.cloud.ibm.com/cloudfoundry/login/ap-south",
		"eu-de":    "https://iam.cloud.ibm.com/cloudfoundry/login/eu-central",
	},
}
var privateRegions = map[string][]string{
	"accounts":              {"us-south", "us-east"},
	"certificate-manager":   {"us-south", "us-east", "eu-gb", "eu-de", "jp-tok", "au-syd", "jp-osa"},
	"icd":                   {"us-south", "us-east", "eu-gb", "eu-de", "jp-tok", "au-syd", "osl01", "seo01", "che01", "ca-tor"},
	"schematics":            {"us-south", "us-east", "eu-de", "eu-gb"},
	"global-search-tagging": {"us-south", "us-east"},
	"container":             {"us-south", "us-east", "eu-gb", "eu-de", "jp-tok", "au-syd", "jp-osa", "ca-tor"},
	"iam":                   {"us-south", "us-east"},
	"resource":              {"us-south", "us-east"},
	"satellite":             {"us-south", "us-east", "eu-gb", "eu-de"},
}
var cloudEndpoint = "cloud.ibm.com"

func contructEndpoint(subdomain, domain string) string {
	endpoint := fmt.Sprintf("https://%s.%s", subdomain, domain)
	return endpoint
}

func validateRegion(region string, regionList []string) (string, error) {
	for _, a := range regionList {
		if a == region {
			return a, nil
		}
	}
	return "", fmt.Errorf("The given region %s doesnot support private endpoints", region)
}

func init() {
	//TODO populate the endpoints which can be retrieved from given endpoints dynamically
	//Example - UAA can be found from the CF endpoint
}

type endpointLocator struct {
	region        string
	visibility    string
	endpointsFile map[string]interface{}
}

//NewEndpointLocator ...
func NewEndpointLocator(region, visibility, file string) EndpointLocator {
	var fileMap map[string]interface{}
	if f := helpers.EnvFallBack([]string{"IBMCLOUD_ENDPOINTS_FILE_PATH", "IC_ENDPOINTS_FILE_PATH"}, file); f != "" {
		jsonFile, err := os.Open(f)
		if err != nil {
			log.Fatalf("Unable to open endpoints file %s", err)
		}
		defer jsonFile.Close()
		bytes, err := ioutil.ReadAll(jsonFile)
		if err != nil {
			log.Fatalf("Unable to read endpoints file %s", err)
		}
		err = json.Unmarshal([]byte(bytes), &fileMap)
		if err != nil {
			log.Fatalf("Unable to unmarshal endpoints file %s", err)
		}
	}
	return &endpointLocator{region: region, visibility: visibility, endpointsFile: fileMap}
}

func (e *endpointLocator) AccountManagementEndpoint() (string, error) {
	//As the current list of regionToEndpoint above is not exhaustive we allow to read endpoints from the env
	endpoint := helpers.EnvFallBack([]string{"IBMCLOUD_ACCOUNT_MANAGEMENT_API_ENDPOINT"}, "")
	if endpoint != "" {
		return endpoint, nil
	}
	if e.endpointsFile != nil && e.visibility != "public-and-private" {
		url := fileFallBack(e.endpointsFile, e.visibility, "IBMCLOUD_ACCOUNT_MANAGEMENT_API_ENDPOINT", e.region, "")
		if url != "" {
			return url, nil
		}
	}
	if e.visibility == "private" || e.visibility == "public-and-private" {
		r, err := validateRegion(e.region, privateRegions["accounts"])
		if err != nil {
			r = "us-south" // As there is no global private endpoint making default region to us-south
			log.Printf("[ WARN ] There is no private endpoint support for this region %s, Defaulting to us-south", e.region)
		}
		return contructEndpoint(fmt.Sprintf("private.%s.accounts", r), cloudEndpoint), nil
	}
	return contructEndpoint("accounts", cloudEndpoint), nil
}

func (e *endpointLocator) CertificateManagerEndpoint() (string, error) {
	//As the current list of regionToEndpoint above is not exhaustive we allow to read endpoints from the env
	endpoint := helpers.EnvFallBack([]string{"IBMCLOUD_CERTIFICATE_MANAGER_API_ENDPOINT"}, "")
	if endpoint != "" {
		return endpoint, nil
	}
	if e.endpointsFile != nil && e.visibility != "public-and-private" {
		url := fileFallBack(e.endpointsFile, e.visibility, "IBMCLOUD_CERTIFICATE_MANAGER_API_ENDPOINT", e.region, "")
		if url != "" {
			return url, nil
		}
	}
	if e.visibility == "private" {
		return contructEndpoint(fmt.Sprintf("private.%s.certificate-manager", e.region), cloudEndpoint), nil
	}
	if e.visibility == "public-and-private" {
		r, err := validateRegion(e.region, privateRegions["certificate-manager"])
		if err != nil {
			return contructEndpoint(fmt.Sprintf("%s.certificate-manager", e.region), cloudEndpoint), nil
		}
		return contructEndpoint(fmt.Sprintf("private.%s.certificate-manager", r), cloudEndpoint), nil
	}
	return contructEndpoint(fmt.Sprintf("%s.certificate-manager", e.region), cloudEndpoint), nil
}

func (e *endpointLocator) CFAPIEndpoint() (string, error) {
	//As the current list of regionToEndpoint above is not exhaustive we allow to read endpoints from the env
	endpoint := helpers.EnvFallBack([]string{"IBMCLOUD_CF_API_ENDPOINT"}, "")
	if endpoint != "" {
		return endpoint, nil
	}
	if e.visibility == "private" {
		return "", bmxerror.New(ErrCodeServiceEndpoint, fmt.Sprintf("Private Endpoints is not supported by this service"))
	}
	if ep, ok := regionToEndpoint["cf"][e.region]; ok {
		return ep, nil
	}
	return "", bmxerror.New(ErrCodeServiceEndpoint, fmt.Sprintf("Cloud Foundry endpoint doesn't exist for region: %q", e.region))
}

func (e *endpointLocator) ContainerEndpoint() (string, error) {
	//As the current list of regionToEndpoint above is not exhaustive we allow to read endpoints from the env
	endpoint := helpers.EnvFallBack([]string{"IBMCLOUD_CS_API_ENDPOINT"}, "")
	if endpoint != "" {
		return endpoint, nil
	}
	if e.endpointsFile != nil && e.visibility != "public-and-private" {
		url := fileFallBack(e.endpointsFile, e.visibility, "IBMCLOUD_CS_API_ENDPOINT", e.region, "")
		if url != "" {
			return url, nil
		}
	}
	if e.visibility == "private" {
		return contructEndpoint(fmt.Sprintf("private.%s.containers", e.region), fmt.Sprintf("%s/global", cloudEndpoint)), nil
	}
	if e.visibility == "public-and-private" {
		r, err := validateRegion(e.region, privateRegions["container"])
		if err != nil {
			return contructEndpoint("containers", fmt.Sprintf("%s/global", cloudEndpoint)), nil
		}
		return contructEndpoint(fmt.Sprintf("private.%s.containers", r), fmt.Sprintf("%s/global", cloudEndpoint)), nil
	}
	return contructEndpoint("containers", fmt.Sprintf("%s/global", cloudEndpoint)), nil
}

func (e *endpointLocator) SchematicsEndpoint() (string, error) {
	//As the current list of regionToEndpoint above is not exhaustive we allow to read endpoints from the env
	endpoint := helpers.EnvFallBack([]string{"IBMCLOUD_SCHEMATICS_API_ENDPOINT"}, "")
	if endpoint != "" {
		return endpoint, nil
	}
	if e.endpointsFile != nil && e.visibility != "public-and-private" {
		url := fileFallBack(e.endpointsFile, e.visibility, "IBMCLOUD_SCHEMATICS_API_ENDPOINT", e.region, "")
		if url != "" {
			return url, nil
		}
	}
	if e.visibility == "private" || e.visibility == "public-and-private" {
		r, err := validateRegion(e.region, privateRegions["schematics"])
		if err != nil {
			r = "us-south"
			log.Printf("[ WARN ] There is no private endpoint support for this region %s, Defaulting to us-south", e.region)
		}
		if r == "us-south" || r == "us-east" {
			return contructEndpoint("private-us.schematics", cloudEndpoint), nil
		}
		if r == "eu-gb" || r == "eu-de" {
			return contructEndpoint("private-eu.schematics", cloudEndpoint), nil
		}
	}
	return contructEndpoint(fmt.Sprintf("%s.schematics", e.region), cloudEndpoint), nil
}

func (e *endpointLocator) ContainerRegistryEndpoint() (string, error) {
	//As the current list of regionToEndpoint above is not exhaustive we allow to read endpoints from the env
	endpoint := helpers.EnvFallBack([]string{"IBMCLOUD_CR_API_ENDPOINT"}, "")
	if endpoint != "" {
		return endpoint, nil
	}
	if e.endpointsFile != nil && e.visibility != "public-and-private" {
		url := fileFallBack(e.endpointsFile, e.visibility, "IBMCLOUD_CR_API_ENDPOINT", e.region, "")
		if url != "" {
			return url, nil
		}
	}
	if ep, ok := regionToEndpoint["cr"][e.region]; ok {
		return fmt.Sprintf("https://%s", ep), nil
	}
	if e.visibility == "private" {
		if ep, ok := regionToEndpoint["cr"][e.region]; ok {
			return contructEndpoint("private", ep), nil
		}
		return "", bmxerror.New(ErrCodeServiceEndpoint, fmt.Sprintf("Container Registry private endpoint doesn't exist for region: %q", e.region))
	}
	if e.visibility == "public-and-private" {
		if ep, ok := regionToEndpoint["cr"][e.region]; ok {
			return contructEndpoint("private", ep), nil
		}
		return "", bmxerror.New(ErrCodeServiceEndpoint, fmt.Sprintf("Container Registry endpoint doesn't exist for region: %q", e.region))
	}
	return "", bmxerror.New(ErrCodeServiceEndpoint, fmt.Sprintf("Container Registry endpoint doesn't exist for region: %q", e.region))
}

// Not used in Provider as we have migrated to go-sdk
func (e *endpointLocator) CisEndpoint() (string, error) {
	//As the current list of regionToEndpoint above is not exhaustive we allow to read endpoints from the env
	endpoint := helpers.EnvFallBack([]string{"IBMCLOUD_CIS_API_ENDPOINT"}, "")
	if endpoint != "" {
		return endpoint, nil
	}
	if e.endpointsFile != nil && e.visibility != "public-and-private" {
		url := fileFallBack(e.endpointsFile, e.visibility, "IBMCLOUD_CIS_API_ENDPOINT", e.region, "")
		if url != "" {
			return url, nil
		}
	}
	if e.visibility == "private" || e.visibility == "public-and-private" {
		return contructEndpoint("api.private.cis", cloudEndpoint), nil
	}
	return contructEndpoint("api.cis", cloudEndpoint), nil
}

func (e *endpointLocator) GlobalSearchEndpoint() (string, error) {
	//As the current list of regionToEndpoint above is not exhaustive we allow to read endpoints from the env
	endpoint := helpers.EnvFallBack([]string{"IBMCLOUD_GS_API_ENDPOINT"}, "")
	if endpoint != "" {
		return endpoint, nil
	}
	if e.endpointsFile != nil && e.visibility != "public-and-private" {
		url := fileFallBack(e.endpointsFile, e.visibility, "IBMCLOUD_GS_API_ENDPOINT", e.region, "")
		if url != "" {
			return url, nil
		}
	}
	if e.visibility == "private" || e.visibility == "public-and-private" {
		r, err := validateRegion(e.region, privateRegions["global-search-tagging"])
		if err != nil {
			r = "us-south" // As there is no global private endpoint making default region to us-south
			log.Printf("[ WARN ] There is no private endpoint support for this region %s, Defaulting to us-south", e.region)
		}
		return contructEndpoint(fmt.Sprintf("api.private.%s.global-search-tagging", r), cloudEndpoint), nil
	}
	return contructEndpoint("api.global-search-tagging", cloudEndpoint), nil
}

func (e *endpointLocator) GlobalTaggingEndpoint() (string, error) {
	//As the current list of regionToEndpoint above is not exhaustive we allow to read endpoints from the env
	endpoint := helpers.EnvFallBack([]string{"IBMCLOUD_GT_API_ENDPOINT"}, "")
	if endpoint != "" {
		return endpoint, nil
	}
	if e.endpointsFile != nil && e.visibility != "public-and-private" {
		url := fileFallBack(e.endpointsFile, e.visibility, "IBMCLOUD_GT_API_ENDPOINT", e.region, "")
		if url != "" {
			return url, nil
		}
	}
	if e.visibility == "private" || e.visibility == "public-and-private" {
		r, err := validateRegion(e.region, privateRegions["global-search-tagging"])
		if err != nil {
			r = "us-south" // As there is no global private endpoint making default region to us-south
			log.Printf("[ WARN ] There is no private endpoint support for this region %s, Defaulting to us-south", e.region)
		}
		return contructEndpoint(fmt.Sprintf("tags.private.%s.global-search-tagging", r), cloudEndpoint), nil
	}
	return contructEndpoint("tags.global-search-tagging", cloudEndpoint), nil
}

func (e *endpointLocator) IAMEndpoint() (string, error) {
	//As the current list of regionToEndpoint above is not exhaustive we allow to read endpoints from the env
	endpoint := helpers.EnvFallBack([]string{"IBMCLOUD_IAM_API_ENDPOINT"}, "")
	if endpoint != "" {
		return endpoint, nil
	}
	if e.endpointsFile != nil && e.visibility != "public-and-private" {
		url := fileFallBack(e.endpointsFile, e.visibility, "IBMCLOUD_IAM_API_ENDPOINT", e.region, "")
		if url != "" {
			return url, nil
		}
	}
	if e.visibility == "private" || e.visibility == "public-and-private" {
		r, err := validateRegion(e.region, privateRegions["iam"])
		if err != nil {
			return contructEndpoint("private.iam", cloudEndpoint), nil
		}
		return contructEndpoint(fmt.Sprintf("private.%s.iam", r), cloudEndpoint), nil
	}
	return contructEndpoint("iam", cloudEndpoint), nil
}

func (e *endpointLocator) IAMPAPEndpoint() (string, error) {
	//As the current list of regionToEndpoint above is not exhaustive we allow to read endpoints from the env
	endpoint := helpers.EnvFallBack([]string{"IBMCLOUD_IAMPAP_API_ENDPOINT"}, "")
	if endpoint != "" {
		return endpoint, nil
	}
	if e.endpointsFile != nil && e.visibility != "public-and-private" {
		url := fileFallBack(e.endpointsFile, e.visibility, "IBMCLOUD_IAMPAP_API_ENDPOINT", e.region, "")
		if url != "" {
			return url, nil
		}
	}
	if e.visibility == "private" || e.visibility == "public-and-private" {
		r, err := validateRegion(e.region, privateRegions["iam"])
		if err != nil {
			return contructEndpoint("private.iam", cloudEndpoint), nil
		}
		return contructEndpoint(fmt.Sprintf("private.%s.iam", r), cloudEndpoint), nil
	}
	return contructEndpoint("iam", cloudEndpoint), nil
}

func (e *endpointLocator) ICDEndpoint() (string, error) {
	//As the current list of regionToEndpoint above is not exhaustive we allow to read endpoints from the env
	endpoint := helpers.EnvFallBack([]string{"IBMCLOUD_ICD_API_ENDPOINT"}, "")
	if endpoint != "" {
		return endpoint, nil
	}
	if e.endpointsFile != nil && e.visibility != "public-and-private" {
		url := fileFallBack(e.endpointsFile, e.visibility, "IBMCLOUD_ICD_API_ENDPOINT", e.region, "")
		if url != "" {
			return url, nil
		}
	}
	if e.visibility == "private" {
		return contructEndpoint(fmt.Sprintf("api.%s.private.databases", e.region), cloudEndpoint), nil
	}
	if e.visibility == "public-and-private" {
		r, err := validateRegion(e.region, privateRegions["icd"])
		if err != nil {
			return contructEndpoint(fmt.Sprintf("api.%s.databases", e.region), cloudEndpoint), nil
		}
		return contructEndpoint(fmt.Sprintf("api.%s.private.databases", r), cloudEndpoint), nil
	}
	return contructEndpoint(fmt.Sprintf("api.%s.databases", e.region), cloudEndpoint), nil
}

func (e *endpointLocator) MCCPAPIEndpoint() (string, error) {
	endpoint := helpers.EnvFallBack([]string{"IBMCLOUD_MCCP_API_ENDPOINT"}, "")
	if endpoint != "" {
		return endpoint, nil
	}
	if e.endpointsFile != nil && e.visibility != "public-and-private" {
		url := fileFallBack(e.endpointsFile, e.visibility, "IBMCLOUD_MCCP_API_ENDPOINT", e.region, "")
		if url != "" {
			return url, nil
		}
	}
	if e.visibility == "private" {
		return "", bmxerror.New(ErrCodeServiceEndpoint, fmt.Sprintf("Private Endpoints is not supported by this service for the region %s", e.region))
	}
	return contructEndpoint(fmt.Sprintf("mccp.%s.cf", e.region), cloudEndpoint), nil
}

func (e *endpointLocator) ResourceManagementEndpoint() (string, error) {
	endpoint := helpers.EnvFallBack([]string{"IBMCLOUD_RESOURCE_MANAGEMENT_API_ENDPOINT"}, "")
	if endpoint != "" {
		return endpoint, nil
	}
	if e.endpointsFile != nil && e.visibility != "public-and-private" {
		url := fileFallBack(e.endpointsFile, e.visibility, "IBMCLOUD_RESOURCE_MANAGEMENT_API_ENDPOINT", e.region, "")
		if url != "" {
			return url, nil
		}
	}
	if e.visibility == "private" {
		r, err := validateRegion(e.region, privateRegions["resource"])
		if err != nil {
			fmt.Println("Private Endpint supports only us-south and us-east region specific endpoint")
			return contructEndpoint("private.resource-controller", cloudEndpoint), nil
		}
		return contructEndpoint(fmt.Sprintf("private.%s.resource-controller", r), cloudEndpoint), nil
	}
	if e.visibility == "public-and-private" {
		r, err := validateRegion(e.region, privateRegions["resource"])
		if err != nil {
			return contructEndpoint("resource-controller", cloudEndpoint), nil
		}
		return contructEndpoint(fmt.Sprintf("private.%s.resource-controller", r), cloudEndpoint), nil
	}
	return contructEndpoint("resource-controller", cloudEndpoint), nil
}

func (e *endpointLocator) ResourceControllerEndpoint() (string, error) {
	endpoint := helpers.EnvFallBack([]string{"IBMCLOUD_RESOURCE_CONTROLLER_API_ENDPOINT"}, "")
	if endpoint != "" {
		return endpoint, nil
	}
	if e.endpointsFile != nil && e.visibility != "public-and-private" {
		url := fileFallBack(e.endpointsFile, e.visibility, "IBMCLOUD_RESOURCE_CONTROLLER_API_ENDPOINT", e.region, "")
		if url != "" {
			return url, nil
		}
	}
	if e.visibility == "private" {
		r, err := validateRegion(e.region, privateRegions["resource"])
		if err != nil {
			fmt.Println("Private Endpint supports only us-south and us-east region specific endpoint")
			return contructEndpoint("private.resource-controller", cloudEndpoint), nil
		}
		return contructEndpoint(fmt.Sprintf("private.%s.resource-controller", r), cloudEndpoint), nil
	}
	if e.visibility == "public-and-private" {
		r, err := validateRegion(e.region, privateRegions["resource"])
		if err != nil {
			return contructEndpoint("resource-controller", cloudEndpoint), nil
		}
		return contructEndpoint(fmt.Sprintf("private.%s.resource-controller", r), cloudEndpoint), nil
	}
	return contructEndpoint("resource-controller", cloudEndpoint), nil
}

func (e *endpointLocator) ResourceCatalogEndpoint() (string, error) {
	endpoint := helpers.EnvFallBack([]string{"IBMCLOUD_RESOURCE_CATALOG_API_ENDPOINT"}, "")
	if endpoint != "" {
		return endpoint, nil
	}
	if e.endpointsFile != nil && e.visibility != "public-and-private" {
		url := fileFallBack(e.endpointsFile, e.visibility, "IBMCLOUD_RESOURCE_CATALOG_API_ENDPOINT", e.region, "")
		if url != "" {
			return url, nil
		}
	}
	if e.visibility == "private" || e.visibility == "public-and-private" {
		r, err := validateRegion(e.region, privateRegions["resource"])
		if err != nil {
			r = "us-south"
		}
		return contructEndpoint(fmt.Sprintf("private.%s.globalcatalog", r), cloudEndpoint), nil
	}
	return contructEndpoint("globalcatalog", cloudEndpoint), nil
}

func (e *endpointLocator) UAAEndpoint() (string, error) {
	endpoint := helpers.EnvFallBack([]string{"IBMCLOUD_UAA_ENDPOINT"}, "")
	if endpoint != "" {
		return endpoint, nil
	}
	if e.endpointsFile != nil && e.visibility != "public-and-private" {
		url := fileFallBack(e.endpointsFile, e.visibility, "IBMCLOUD_UAA_ENDPOINT", e.region, "")
		if url != "" {
			return url, nil
		}
	}
	if e.visibility == "private" {
		return "", bmxerror.New(ErrCodeServiceEndpoint, fmt.Sprintf("Private Endpoints is not supported by this service for the region %s", e.region))
	}
	if ep, ok := regionToEndpoint["uaa"][e.region]; ok {
		//As the current list of regionToEndpoint above is not exhaustive we allow to read endpoints from the env
		return ep, nil
	}
	return "", bmxerror.New(ErrCodeServiceEndpoint, fmt.Sprintf("UAA endpoint doesn't exist for region: %q", e.region))
}

func (e *endpointLocator) CseEndpoint() (string, error) {
	endpoint := helpers.EnvFallBack([]string{"IBMCLOUD_CSE_ENDPOINT"}, "")
	if endpoint != "" {
		return endpoint, nil
	}
	if e.endpointsFile != nil && e.visibility != "public-and-private" {
		url := fileFallBack(e.endpointsFile, e.visibility, "IBMCLOUD_CSE_ENDPOINT", e.region, "")
		if url != "" {
			return url, nil
		}
	}
	if e.visibility == "private" {
		return "", bmxerror.New(ErrCodeServiceEndpoint, fmt.Sprintf("Private Endpoints is not supported by this service"))
	}
	return contructEndpoint("api.serviceendpoint", cloudEndpoint), nil
}

func (e *endpointLocator) UserManagementEndpoint() (string, error) {
	endpoint := helpers.EnvFallBack([]string{"IBMCLOUD_USER_MANAGEMENT_ENDPOINT"}, "")
	if endpoint != "" {
		return endpoint, nil
	}
	if e.endpointsFile != nil && e.visibility != "public-and-private" {
		url := fileFallBack(e.endpointsFile, e.visibility, "IBMCLOUD_USER_MANAGEMENT_ENDPOINT", e.region, "")
		if url != "" {
			return url, nil
		}
	}
	if e.visibility == "private" || e.visibility == "public-and-private" {
		r, err := validateRegion(e.region, privateRegions["resource"])
		if err != nil {
			r = "us-south"
		}
		return contructEndpoint(fmt.Sprintf("private.%s.user-management", r), cloudEndpoint), nil
	}
	return contructEndpoint("user-management", cloudEndpoint), nil
}

func (e *endpointLocator) HpcsEndpoint() (string, error) {
	endpoint := helpers.EnvFallBack([]string{"IBMCLOUD_HPCS_API_ENDPOINT"}, "")
	if endpoint != "" {
		return endpoint, nil
	}
	if e.endpointsFile != nil && e.visibility != "public-and-private" {
		url := fileFallBack(e.endpointsFile, e.visibility, "IBMCLOUD_HPCS_API_ENDPOINT", e.region, "")
		if url != "" {
			return url, nil
		}
	}
	if e.visibility == "private" {
		return "", bmxerror.New(ErrCodeServiceEndpoint, fmt.Sprintf("Private Endpoints is not supported by this service for the region %s", e.region))
	}
	return fmt.Sprintf("https://%s.broker.hs-crypto.cloud.ibm.com/crypto_v2/", e.region), nil
}

func (e *endpointLocator) FunctionsEndpoint() (string, error) {
	endpoint := helpers.EnvFallBack([]string{"IBMCLOUD_FUNCTIONS_API_ENDPOINT"}, "")
	if endpoint != "" {
		return endpoint, nil
	}
	if e.endpointsFile != nil && e.visibility != "public-and-private" {
		url := fileFallBack(e.endpointsFile, e.visibility, "IBMCLOUD_FUNCTIONS_API_ENDPOINT", e.region, "")
		if url != "" {
			return url, nil
		}
	}
	if e.visibility == "private" {
		return "", bmxerror.New(ErrCodeServiceEndpoint, fmt.Sprintf("Private Endpoints is not supported by this service for the region %s", e.region))
	}
	return contructEndpoint(fmt.Sprintf("%s.functions", e.region), cloudEndpoint), nil
}

func (e *endpointLocator) SatelliteEndpoint() (string, error) {
	//As the current list of regionToEndpoint above is not exhaustive we allow to read endpoints from the env
	endpoint := helpers.EnvFallBack([]string{"IBMCLOUD_SAT_API_ENDPOINT"}, "")
	if endpoint != "" {
		return endpoint, nil
	}
	if e.endpointsFile != nil && e.visibility != "public-and-private" {
		url := fileFallBack(e.endpointsFile, e.visibility, "IBMCLOUD_SAT_API_ENDPOINT", e.region, "")
		if url != "" {
			return url, nil
		}
	}
	if e.visibility == "private" {
		return contructEndpoint(fmt.Sprintf("private.%s.api.link.satellite", e.region), fmt.Sprintf("%s", cloudEndpoint)), nil
	}
	if e.visibility == "public-and-private" {
		r, err := validateRegion(e.region, privateRegions["satellite"])
		if err != nil {
			return contructEndpoint("api.link.satellite", fmt.Sprintf("%s", cloudEndpoint)), nil
		}
		return contructEndpoint(fmt.Sprintf("private.%s.api.link.satellite", r), fmt.Sprintf("%s", cloudEndpoint)), nil
	}
	return contructEndpoint("api.link.satellite", fmt.Sprintf("%s", cloudEndpoint)), nil
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
