package rosa

import (
	"fmt"

	cmv1 "github.com/openshift-online/ocm-sdk-go/clustersmgmt/v1"
)

const (
	rosaCreatorArnProperty = "rosa_creator_arn"
)

// CreateCluster creates a new ROSA cluster using the specified spec.
func (c *RosaClient) CreateCluster(spec *cmv1.Cluster) (*cmv1.Cluster, error) {
	cluster, err := c.ocm.ClustersMgmt().V1().Clusters().
		Add().
		Body(spec).
		Send()
	if err != nil {
		return nil, handleErr(cluster.Error(), err)
	}

	clusterObject := cluster.Body()
	return clusterObject, nil
}

// DeleteCluster deletes the ROSA cluster.
func (c *RosaClient) DeleteCluster(clusterID string) error {
	response, err := c.ocm.ClustersMgmt().V1().Clusters().
		Cluster(clusterID).
		Delete().
		BestEffort(true).
		Send()
	if err != nil {
		return handleErr(response.Error(), err)
	}

	return nil
}

// GetCluster retrieves the ROSA/OCM cluster object.
func (c *RosaClient) GetCluster() (*cmv1.Cluster, error) {
	clusterKey := c.rosaScope.RosaClusterName()
	query := fmt.Sprintf("%s AND (id = '%s' OR name = '%s' OR external_id = '%s')",
		getClusterFilter(c.rosaScope.Identity.Arn),
		clusterKey, clusterKey, clusterKey,
	)
	response, err := c.ocm.ClustersMgmt().V1().Clusters().List().
		Search(query).
		Page(1).
		Size(1).
		Send()
	if err != nil {
		return nil, handleErr(response.Error(), err)
	}

	switch response.Total() {
	case 0:
		return nil, nil
	case 1:
		return response.Items().Slice()[0], nil
	default:
		return nil, fmt.Errorf("there are %d clusters with identifier or name '%s'", response.Total(), clusterKey)
	}
}

// Generate a query that filters clusters running on the current AWS session account.
func getClusterFilter(creatorArn *string) string {
	filter := "product.id = 'rosa'"
	if creatorArn != nil {
		filter = fmt.Sprintf("%s AND (properties.%s = '%s')",
			filter,
			rosaCreatorArnProperty,
			*creatorArn)
	}
	return filter
}
