package ocm

import (
	cmv1 "github.com/openshift-online/ocm-sdk-go/clustersmgmt/v1"
)

func (c *Client) GetLimitedSupportReasons(clusterID string) (
	limitedSupportReasons []*cmv1.LimitedSupportReason, err error) {
	collection := c.ocm.ClustersMgmt().V1().
		Clusters().
		Cluster(clusterID).
		LimitedSupportReasons()
	page := 1
	size := 100
	search := "override_enabled='f'"
	for {
		response, err := collection.List().
			Page(page).
			Parameter("search", search).
			Size(size).
			Send()
		if err != nil {
			return nil, handleErr(response.Error(), err)
		}
		limitedSupportReasons = append(limitedSupportReasons, response.Items().Slice()...)
		if response.Size() < size {
			break
		}
		page++
	}
	return limitedSupportReasons, nil
}
