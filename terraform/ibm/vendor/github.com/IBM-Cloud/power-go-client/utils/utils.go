package utils

import (
	"net/url"
	"os"
	"reflect"

	"github.com/IBM-Cloud/power-go-client/helpers"
)

// GetNext ...
func GetNext(next interface{}) string {
	if reflect.ValueOf(next).IsNil() {
		return ""
	}

	u, err := url.Parse(reflect.ValueOf(next).Elem().FieldByName("Href").Elem().String())
	if err != nil {
		return ""
	}

	q := u.Query()
	return q.Get("start")
}

// GetEndpoint ...
func GetEndpoint(generation int, regionName string) string {

	switch generation {
	case 1:
		ep := getGCEndpoint(regionName)
		return helpers.EnvFallBack([]string{"IBMCLOUD_IS_API_ENDPOINT"}, ep)
	case 2:
		ep := getNGEndpoint(regionName)
		return helpers.EnvFallBack([]string{"IBMCLOUD_IS_NG_API_ENDPOINT"}, ep)
	}
	ep := getNGEndpoint(regionName)
	return helpers.EnvFallBack([]string{"IBMCLOUD_IS_NG_API_ENDPOINT"}, ep)
}

func getGCEndpoint(regionName string) string {
	if url := os.Getenv("IBMCLOUD_IS_API_ENDPOINT"); url != "" {
		return url
	}
	return regionName + ".iaas.cloud.ibm.com"
}

// For Power-IAAS
func getNGEndpoint(regionName string) string {
	if url := os.Getenv("IBMCLOUD_IS_NG_API_ENDPOINT"); url != "" {
		return url
	}
	return regionName + ".power-iaas.cloud.ibm.com"
}

func GetPowerEndPoint(regionName string) string {
	if url := os.Getenv("IBMCLOUD_POWER_API_ENDPOINT"); url != "" {
		return url
	}
	return regionName + ".power-iaas.cloud.ibm.com"

}
