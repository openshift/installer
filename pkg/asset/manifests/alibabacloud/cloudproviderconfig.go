package alibabacloud

import (
	"bytes"
	"encoding/json"
)

// CloudConfig wraps the settings for the Alibaba Cloud provider.
// ref: https://github.com/kubernetes/cloud-provider-alibaba-cloud/blob/d6d0962b4be051c7dc536dc7e49ad0aff018ef3b/cloud-controller-manager/alicloud.go#L90
type CloudConfig struct {
	Global GlobalConfig `json:"Global"`
}

// GlobalConfig wraps the settings 'Global' for the Alibaba Cloud provider .
type GlobalConfig struct {
	KubernetesClusterTag string `json:"kubernetesClusterTag"`
	NodeMonitorPeriod    int64  `json:"nodeMonitorPeriod"`
	NodeAddrSyncPeriod   int64  `json:"nodeAddrSyncPeriod"`
	UID                  string `json:"uid"`
	VpcID                string `json:"vpcid"`
	Region               string `json:"region"`
	ZoneID               string `json:"zoneid"`
	VswitchID            string `json:"vswitchid"`
	ClusterID            string `json:"clusterID"`
	RouteTableIDs        string `json:"routeTableIDs"`
	ServiceBackendType   string `json:"serviceBackendType"`

	DisablePublicSLB bool `json:"disablePublicSLB"`

	AccessKeyID     string `json:"accessKeyID"`
	AccessKeySecret string `json:"accessKeySecret"`
}

// JSON generates the cloud provider json config for the Alibaba Cloud platform.
// managed resource names are matching the convention defined by capz
func (params CloudConfig) JSON() (string, error) {
	buff := &bytes.Buffer{}
	encoder := json.NewEncoder(buff)
	encoder.SetIndent("", "\t")
	if err := encoder.Encode(params); err != nil {
		return "", err
	}
	return buff.String(), nil
}
