package azure

import (
	"fmt"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore/to"
)

// GetSystemTags gets the cluster owned system tags.
func GetSystemTags(clusterID string) map[string]*string {
	return map[string]*string{
		fmt.Sprintf("kubernetes.io_cluster.%s", clusterID): to.Ptr("owned"),
	}
}

// ConvertSDKTagsToCAPZTags converts sdk tags to capz tags.
func ConvertSDKTagsToCAPZTags(tags map[string]*string) map[string]string {
	capzTags := make(map[string]string, len(tags)+1)
	for k, v := range tags {
		if v != nil {
			capzTags[k] = *v
		}
	}
	return capzTags
}

// ConvertCAPZTagsToSDKTags converts capz tags to sdk tags.
func ConvertCAPZTagsToSDKTags(tags map[string]string) map[string]*string {
	sdkTags := make(map[string]*string, len(tags)+1)
	for k, v := range tags {
		sdkTags[k] = to.Ptr(v)
	}
	return sdkTags
}
