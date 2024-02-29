package ocm

import (
	"net/http"

	cmv1 "github.com/openshift-online/ocm-sdk-go/clustersmgmt/v1"
)

type KubeletConfigArgs struct {
	PodPidsLimit int
}

func (c *Client) GetClusterKubeletConfig(clusterID string) (*cmv1.KubeletConfig, error) {
	response, err := c.ocm.ClustersMgmt().V1().Clusters().Cluster(clusterID).KubeletConfig().Get().Send()

	if response.Status() == http.StatusNotFound {
		return nil, nil
	}

	if err != nil {
		return nil, handleErr(response.Error(), err)
	}

	return response.Body(), nil
}

func (c *Client) DeleteKubeletConfig(clusterID string) error {
	response, err := c.ocm.ClustersMgmt().V1().Clusters().Cluster(clusterID).KubeletConfig().Delete().Send()
	if err != nil {
		return handleErr(response.Error(), err)
	}
	return nil
}

func toOCMKubeletConfig(args KubeletConfigArgs) (*cmv1.KubeletConfig, error) {
	builder := &cmv1.KubeletConfigBuilder{}
	kubeletConfig, err := builder.PodPidsLimit(args.PodPidsLimit).Build()
	if err != nil {
		return nil, err
	}

	return kubeletConfig, nil
}

func (c *Client) CreateKubeletConfig(clusterID string, args KubeletConfigArgs) (*cmv1.KubeletConfig, error) {

	kubeletConfig, err := toOCMKubeletConfig(args)
	if err != nil {
		return nil, err
	}

	response, err := c.ocm.ClustersMgmt().V1().Clusters().Cluster(clusterID).
		KubeletConfig().Post().Body(kubeletConfig).Send()

	if err != nil {
		return nil, err
	}

	return response.Body(), nil
}

func (c *Client) UpdateKubeletConfig(clusterID string, args KubeletConfigArgs) (*cmv1.KubeletConfig, error) {
	kubeletConfig, err := toOCMKubeletConfig(args)
	if err != nil {
		return nil, err
	}

	response, err := c.ocm.ClustersMgmt().V1().Clusters().Cluster(clusterID).
		KubeletConfig().Update().Body(kubeletConfig).Send()

	if err != nil {
		return nil, err
	}

	return response.Body(), nil
}
