package aws

import (
	"strings"

	"github.com/aws/aws-sdk-go-v2/service/ec2/types"
	"k8s.io/apimachinery/pkg/util/sets"
	"k8s.io/utils/ptr"
)

const (
	// TagNameKubernetesClusterPrefix is the tag name prefix used by CCM
	// to differentiate multiple logically independent clusters running in the same AZ.
	TagNameKubernetesClusterPrefix = "kubernetes.io/cluster/"

	// TagNameAWSProviderClusterPrefix is the tag prefix added by CAPI/CAPA to differentiate
	// cluster-api-provider-aws owned components.
	TagNameAWSProviderClusterPrefix = "sigs.k8s.io/cluster-api-provider-aws/cluster/"

	// TagValueOwned is the tag value to indicate that a resource is considered owned
	// and managed by the cluster.
	TagValueOwned = "owned"

	// TagValueShared is the tag value to indicate that a resource is considered shared
	// with the cluster.
	TagValueShared = "shared"

	// TagNameKubernetesUnmanaged is the tag name to indicate that a resource is unmanaged
	// by the cluster and should be ignored by CCM. For example, kubernetes.io/cluster/unmanaged=true.
	TagNameKubernetesUnmanaged = TagNameKubernetesClusterPrefix + "unmanaged"
)

// Tags represents AWS resource tags as a map.
// This helps avoid iterating over the tag list for every lookup.
type Tags map[string]string

// FromAWSTags converts a list of AWS tags into a map.
func FromAWSTags(awsTags []types.Tag) Tags {
	tags := make(Tags, len(awsTags))
	for _, tag := range awsTags {
		key, value := ptr.Deref(tag.Key, ""), ptr.Deref(tag.Value, "")
		if len(key) > 0 {
			tags[key] = value
		}
	}
	return tags
}

// GetTagKeysWithPrefix return tag keys matching a given prefix.
func (t Tags) GetTagKeysWithPrefix(prefix string) []string {
	keys := make([]string, 0)
	for key := range t {
		if strings.HasPrefix(key, prefix) {
			keys = append(keys, key)
		}
	}
	return keys
}

// HasTagKeyPrefix returns true if there is a tag with a given key prefix.
func (t Tags) HasTagKeyPrefix(prefix string) bool {
	keys := t.GetTagKeysWithPrefix(prefix)
	return len(keys) > 0
}

// getClusterIDsWithPrefix is a helper to extract cluster IDs from tags with a given prefix and resource lifecycle.
// The values of resourceLifeCycle can be "owned" or "shared".
func (t Tags) getClusterIDsWithPrefix(tagPrefix string, resourceLifeCycle string) []string {
	clusterIDs := make([]string, 0)
	keys := t.GetTagKeysWithPrefix(tagPrefix)
	for _, key := range keys {
		if value := t[key]; value == resourceLifeCycle {
			if clusterID := strings.TrimPrefix(key, tagPrefix); clusterID != "" {
				clusterIDs = append(clusterIDs, clusterID)
			}
		}
	}
	return clusterIDs
}

// GetClusterIDs returns the unique cluster IDs from either tag if any:
// - "kubernetes.io/cluster/<cluster-id>: <resourceLifeCycle>"
// - "sigs.k8s.io/cluster-api-provider-aws/cluster/<cluster-id>: <resourceLifeCycle>"
// The values of resourceLifeCycle can be "owned" or "shared".
func (t Tags) GetClusterIDs(resourceLifeCycle string) []string {
	clusterIDs := t.getClusterIDsWithPrefix(TagNameKubernetesClusterPrefix, resourceLifeCycle)
	clusterIDs = append(clusterIDs, t.getClusterIDsWithPrefix(TagNameAWSProviderClusterPrefix, resourceLifeCycle)...)

	uniqueClusterIDs := sets.New[string]()
	for _, id := range clusterIDs {
		uniqueClusterIDs.Insert(id)
	}
	return uniqueClusterIDs.UnsortedList()
}
