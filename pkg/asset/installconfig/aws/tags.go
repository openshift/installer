package aws

import (
	"strings"

	"github.com/aws/aws-sdk-go-v2/service/ec2/types"
	"k8s.io/utils/ptr"
)

const (
	// TagNameKubernetesClusterPrefix is the tag name prefix used by CCM
	// to differentiate multiple logically independent clusters running in the same AZ.
	TagNameKubernetesClusterPrefix = "kubernetes.io/cluster/"

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

// HasTagKeyPrefix returns true if there is a tag with a given key prefix.
func (t Tags) HasTagKeyPrefix(prefix string) bool {
	for key := range t {
		if strings.HasPrefix(key, prefix) {
			return true
		}
	}
	return false
}
