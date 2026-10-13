package ocm

import (
	cmv1 "github.com/openshift-online/ocm-sdk-go/clustersmgmt/v1"
)

func (c *Client) GetInflightChecks(clusterID string) (
	inflightChecks []*cmv1.InflightCheck, err error) {
	collection := c.ocm.ClustersMgmt().V1().
		Clusters().
		Cluster(clusterID).
		InflightChecks()
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
		inflightChecks = append(inflightChecks, response.Items().Slice()...)
		if response.Size() < size {
			break
		}
		page++
	}
	return inflightChecks, nil
}
