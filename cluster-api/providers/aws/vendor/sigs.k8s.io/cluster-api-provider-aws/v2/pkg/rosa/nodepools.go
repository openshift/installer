package rosa

import (
	"time"

	cmv1 "github.com/openshift-online/ocm-sdk-go/clustersmgmt/v1"
)

// CreateNodePool adds a new node pool to the cluster.
func (c *RosaClient) CreateNodePool(clusterID string, nodePool *cmv1.NodePool) (*cmv1.NodePool, error) {
	response, err := c.ocm.ClustersMgmt().V1().
		Clusters().Cluster(clusterID).
		NodePools().
		Add().Body(nodePool).
		Send()
	if err != nil {
		return nil, handleErr(response.Error(), err)
	}
	return response.Body(), nil
}

// GetNodePools retrieves the list of node pools in the cluster.
func (c *RosaClient) GetNodePools(clusterID string) ([]*cmv1.NodePool, error) {
	response, err := c.ocm.ClustersMgmt().V1().
		Clusters().Cluster(clusterID).
		NodePools().
		List().Page(1).Size(-1).
		Send()
	if err != nil {
		return nil, handleErr(response.Error(), err)
	}
	return response.Items().Slice(), nil
}

// GetNodePool retrieves the details of the specified node pool.
func (c *RosaClient) GetNodePool(clusterID string, nodePoolID string) (*cmv1.NodePool, bool, error) {
	response, err := c.ocm.ClustersMgmt().V1().
		Clusters().Cluster(clusterID).
		NodePools().
		NodePool(nodePoolID).
		Get().
		Send()
	if response.Status() == 404 {
		return nil, false, nil
	}
	if err != nil {
		return nil, false, handleErr(response.Error(), err)
	}
	return response.Body(), true, nil
}

// UpdateNodePool updates the specified node pool.
func (c *RosaClient) UpdateNodePool(clusterID string, nodePool *cmv1.NodePool) (*cmv1.NodePool, error) {
	response, err := c.ocm.ClustersMgmt().V1().
		Clusters().Cluster(clusterID).
		NodePools().NodePool(nodePool.ID()).
		Update().Body(nodePool).
		Send()
	if err != nil {
		return nil, handleErr(response.Error(), err)
	}
	return response.Body(), nil
}

// DeleteNodePool deletes the specified node pool.
func (c *RosaClient) DeleteNodePool(clusterID string, nodePoolID string) error {
	response, err := c.ocm.ClustersMgmt().V1().
		Clusters().Cluster(clusterID).
		NodePools().NodePool(nodePoolID).
		Delete().
		Send()
	if err != nil {
		return handleErr(response.Error(), err)
	}
	return nil
}

// CheckNodePoolExistingScheduledUpgrade checks and returns the current upgrade schedule for the nodePool if any.
func (c *RosaClient) CheckNodePoolExistingScheduledUpgrade(clusterID string, nodePool *cmv1.NodePool) (*cmv1.NodePoolUpgradePolicy, error) {
	upgradePolicies, err := c.getNodePoolUpgradePolicies(clusterID, nodePool.ID())
	if err != nil {
		return nil, err
	}
	for _, upgradePolicy := range upgradePolicies {
		if upgradePolicy.UpgradeType() == cmv1.UpgradeTypeNodePool {
			return upgradePolicy, nil
		}
	}
	return nil, nil
}

// ScheduleNodePoolUpgrade schedules a new nodePool upgrade to the specified version at the specified time.
func (c *RosaClient) ScheduleNodePoolUpgrade(clusterID string, nodePool *cmv1.NodePool, version string, nextRun time.Time) (*cmv1.NodePoolUpgradePolicy, error) {
	// earliestNextRun is set to at least 5 min from now by the OCM API.
	// we set it to 6 min here to account for latencty.
	earliestNextRun := time.Now().Add(time.Minute * 6)
	if nextRun.Before(earliestNextRun) {
		nextRun = earliestNextRun
	}

	upgradePolicy, err := cmv1.NewNodePoolUpgradePolicy().
		UpgradeType(cmv1.UpgradeTypeNodePool).
		NodePoolID(nodePool.ID()).
		ScheduleType(cmv1.ScheduleTypeManual).
		Version(version).
		NextRun(nextRun).
		Build()
	if err != nil {
		return nil, err
	}

	response, err := c.ocm.ClustersMgmt().V1().
		Clusters().Cluster(clusterID).
		NodePools().
		NodePool(nodePool.ID()).UpgradePolicies().
		Add().Body(upgradePolicy).
		Send()
	if err != nil {
		return nil, handleErr(response.Error(), err)
	}

	return response.Body(), nil
}

func (c *RosaClient) getNodePoolUpgradePolicies(clusterID string, nodePoolID string) (nodePoolUpgradePolicies []*cmv1.NodePoolUpgradePolicy, err error) {
	collection := c.ocm.ClustersMgmt().V1().
		Clusters().
		Cluster(clusterID).NodePools().NodePool(nodePoolID).UpgradePolicies()
	page := 1
	size := 100
	for {
		response, err := collection.List().
			Page(page).
			Size(size).
			Send()
		if err != nil {
			return nil, handleErr(response.Error(), err)
		}
		nodePoolUpgradePolicies = append(nodePoolUpgradePolicies, response.Items().Slice()...)
		if response.Size() < size {
			break
		}
		page++
	}
	return
}
