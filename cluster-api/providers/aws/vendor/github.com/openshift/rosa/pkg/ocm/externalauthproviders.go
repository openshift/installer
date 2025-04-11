package ocm

import cmv1 "github.com/openshift-online/ocm-sdk-go/clustersmgmt/v1"

func (c *Client) CreateExternalAuth(clusterID string, ExternalAuth *cmv1.ExternalAuth) (*cmv1.ExternalAuth, error) {
	response, err := c.ocm.ClustersMgmt().V1().
		Clusters().Cluster(clusterID).
		ExternalAuthConfig().ExternalAuths().Add().Body(ExternalAuth).Send()
	if err != nil {
		return nil, handleErr(response.Error(), err)
	}
	return response.Body(), nil
}

func (c *Client) GetExternalAuth(clusterID string, externalAuthId string) (*cmv1.ExternalAuth, bool, error) {
	response, err := c.ocm.ClustersMgmt().V1().
		Clusters().Cluster(clusterID).ExternalAuthConfig().
		ExternalAuths().ExternalAuth(externalAuthId).
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

func (c *Client) GetExternalAuths(clusterID string) ([]*cmv1.ExternalAuth, error) {
	response, err := c.ocm.ClustersMgmt().V1().
		Clusters().Cluster(clusterID).
		ExternalAuthConfig().
		ExternalAuths().
		List().Page(1).Size(-1).
		Send()
	if err != nil {
		return nil, handleErr(response.Error(), err)
	}
	return response.Items().Slice(), nil
}

func (c *Client) DeleteExternalAuth(clusterID string, externalAuthId string) error {
	response, err := c.ocm.ClustersMgmt().V1().
		Clusters().Cluster(clusterID).
		ExternalAuthConfig().ExternalAuths().
		ExternalAuth(externalAuthId).
		Delete().
		Send()
	if err != nil {
		return handleErr(response.Error(), err)
	}
	return nil
}
