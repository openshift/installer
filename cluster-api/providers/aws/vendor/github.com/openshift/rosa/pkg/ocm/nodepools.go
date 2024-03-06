package ocm

import cmv1 "github.com/openshift-online/ocm-sdk-go/clustersmgmt/v1"

func (c *Client) CreateNodePool(clusterID string, nodePool *cmv1.NodePool) (*cmv1.NodePool, error) {
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

func (c *Client) GetNodePools(clusterID string) ([]*cmv1.NodePool, error) {
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

func (c *Client) GetNodePool(clusterID string, nodePoolID string) (*cmv1.NodePool, bool, error) {
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

func (c *Client) UpdateNodePool(clusterID string, nodePool *cmv1.NodePool) (*cmv1.NodePool, error) {
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

func (c *Client) DeleteNodePool(clusterID string, nodePoolID string) error {
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
