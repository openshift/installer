package aws

import (
	"context"
	"encoding/json"
	"fmt"
	"os"

	"github.com/ghodss/yaml"
	"github.com/openshift/installer/pkg/installerassets"
	"github.com/pkg/errors"
	awsprovider "sigs.k8s.io/cluster-api-provider-aws/pkg/apis/awsproviderconfig/v1alpha1"
)

func getUserTags(ctx context.Context) (data []byte, err error) {
	userTags := map[string]string{}
	if value, ok := os.LookupEnv("_CI_ONLY_STAY_AWAY_OPENSHIFT_INSTALL_AWS_USER_TAGS"); ok {
		if err := json.Unmarshal([]byte(value), &userTags); err != nil {
			return nil, errors.Wrapf(err, "_CI_ONLY_STAY_AWAY_OPENSHIFT_INSTALL_AWS_USER_TAGS contains invalid JSON: %s", value)
		}
	}

	return yaml.Marshal(userTags)
}

func tagsFromUserTags(clusterID, clusterName string, userTags map[string]string) ([]awsprovider.TagSpecification, error) {
	tags := []awsprovider.TagSpecification{
		{Name: "tectonicClusterID", Value: clusterID},
		{Name: fmt.Sprintf("kubernetes.io/cluster/%s", clusterName), Value: "owned"},
	}
	forbiddenTags := map[string]bool{}
	for _, tag := range tags {
		forbiddenTags[tag.Name] = true
	}
	for key, value := range userTags {
		if forbiddenTags[key] {
			return nil, errors.Errorf("user tags may not clobber %s", key)
		}
		tags = append(tags, awsprovider.TagSpecification{Name: key, Value: value})
	}
	return tags, nil
}

func init() {
	installerassets.Defaults["aws/user-tags"] = getUserTags
}
