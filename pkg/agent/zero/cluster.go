package zero

import (
	"context"

	"github.com/go-openapi/strfmt"
	"github.com/sirupsen/logrus"

	"github.com/openshift/assisted-service/client/installer"
	"github.com/openshift/assisted-service/models"
)

type ClusterZero struct {
	Ctx                   context.Context
	Api                   *zeroClient
	clusterZeroID         *strfmt.UUID
	clusterZeroInfraEnvID *strfmt.UUID
}

type zeroClient struct {
	Kube *clusterZeroKubeAPIClient
	Rest *nodeZeroRestClient
}

func NewClusterZero(ctx context.Context, assetDir string) (*ClusterZero, error) {

	czero := &ClusterZero{}
	capi := &zeroClient{}

	restclient, err := NewNodeZeroRestClient(ctx, assetDir)
	if err != nil {
		logrus.Fatal(err)
	}
	kubeclient, err := NewClusterZeroKubeAPIClient(ctx, assetDir)
	if err != nil {
		logrus.Fatal(err)
	}

	capi.Rest = restclient
	capi.Kube = kubeclient

	clusterZeroID, err := restclient.getClusterZeroClusterID()
	if err != nil {
		return nil, err
	}
	clusterZeroInfraEnvID, err := restclient.getClusterZeroInfraEnvID()
	if err != nil {
		return nil, err
	}

	czero.Ctx = ctx
	czero.Api = capi
	czero.clusterZeroID = clusterZeroID
	czero.clusterZeroInfraEnvID = clusterZeroInfraEnvID
	return czero, nil
}

func (czero *ClusterZero) Get() (*models.Cluster, error) {
	// GET /v2/clusters/{cluster_zero_id}
	getClusterParams := &installer.V2GetClusterParams{ClusterID: *czero.clusterZeroID}
	result, err := czero.Api.Rest.Client.Installer.V2GetCluster(czero.Ctx, getClusterParams)
	if err != nil {
		return nil, err
	}
	clusterZero := result.Payload
	return clusterZero, nil
}

// TODO(lranjbar): Print install status from the Cluster object
func (czero *ClusterZero) PrintInstallStatus(*models.Cluster) error {

	return nil
}

// TODO(lranjbar)[AGENT-172]: Need to parse the validations_info object returned by the REST API
// Example Response from /v2/clusters/:
// *models.Cluster I expect have a validations_info JSON object to marshal
// [
//   {
//     "api_vip": "192.168.111.5",
//     "base_dns_domain": "test.metalkube.org",
//     "cluster_networks": [
//       {
//         "cidr": "10.128.0.0/14",
//         "cluster_id": "bfe541fa-9494-4bcc-8c45-3ebba77a7344",
//         "host_prefix": 23
//       }
//     ],
//     "connectivity_majority_groups": "{\"192.168.110.0/23\":[\"303f1bcb-8f34-4b58-94d1-f3acc1ec3a10\",\"91774a06-ef2d-4f76-9df6-8dd58ccd5ef0\",\"c28d55ad-4054-4a37-be60-d44b37d561fc\",\"f89c67fc-8b12-416e-a762-a2de58204fff\",\"f93bf19f-50d1-4256-ab5b-26ec43dd88d4\"],\"IPv4\":[\"303f1bcb-8f34-4b58-94d1-f3acc1ec3a10\",\"91774a06-ef2d-4f76-9df6-8dd58ccd5ef0\",\"c28d55ad-4054-4a37-be60-d44b37d561fc\",\"f89c67fc-8b12-416e-a762-a2de58204fff\",\"f93bf19f-50d1-4256-ab5b-26ec43dd88d4\"],\"IPv6\":[\"303f1bcb-8f34-4b58-94d1-f3acc1ec3a10\",\"91774a06-ef2d-4f76-9df6-8dd58ccd5ef0\",\"c28d55ad-4054-4a37-be60-d44b37d561fc\",\"f89c67fc-8b12-416e-a762-a2de58204fff\",\"f93bf19f-50d1-4256-ab5b-26ec43dd88d4\"]}",
//     "controller_logs_collected_at": "0001-01-01T00:00:00.000Z",
//     "controller_logs_started_at": "0001-01-01T00:00:00.000Z",
//     "cpu_architecture": "x86_64",
//     "created_at": "2022-06-08T21:06:52.523858Z",
//     "deleted_at": null,
//     "disk_encryption": {
//       "enable_on": "none",
//       "mode": "tpmv2"
//     },
//     "email_domain": "Unknown",
//     "enabled_host_count": 5,
//     "feature_usage": "{\"auto assign role\":{\"id\":\"AUTO_ASSIGN_ROLE\",\"name\":\"auto assign role\"}}",
//     "host_networks": null,
//     "hosts": [],
//     "href": "/api/assisted-install/v2/clusters/bfe541fa-9494-4bcc-8c45-3ebba77a7344",
//     "hyperthreading": "all",
//     "id": "bfe541fa-9494-4bcc-8c45-3ebba77a7344",
//     "ignition_endpoint": {},
//     "image_info": {
//       "created_at": "2022-06-08T21:06:52.523858Z",
//       "expires_at": "0001-01-01T00:00:00.000Z"
//     },
//     "ingress_vip": "192.168.111.4",
//     "install_completed_at": "0001-01-01T00:00:00.000Z",
//     "install_started_at": "2022-06-08T21:08:38.876Z",
//     "kind": "Cluster",
//     "machine_networks": [
//       {
//         "cidr": "192.168.110.0/23",
//         "cluster_id": "bfe541fa-9494-4bcc-8c45-3ebba77a7344"
//       }
//     ],
//     "monitored_operators": [
//       {
//         "cluster_id": "bfe541fa-9494-4bcc-8c45-3ebba77a7344",
//         "name": "console",
//         "operator_type": "builtin",
//         "status_updated_at": "0001-01-01T00:00:00.000Z",
//         "timeout_seconds": 3600
//       }
//     ],
//     "name": "ostest",
//     "ocp_release_image": "registry.ci.openshift.org/ocp/release:4.11.0-0.nightly-2022-06-06-201913",
//     "openshift_cluster_id": "ec5e3943-7d7c-46a9-a6c6-193689137fbb",
//     "openshift_version": "4.11.0-0.nightly-2022-06-06-201913",
//     "platform": {
//       "ovirt": {},
//       "type": "baremetal"
//     },
//     "progress": {
//       "installing_stage_percentage": 52,
//       "preparing_for_installation_stage_percentage": 100,
//       "total_percentage": 46
//     },
//     "pull_secret_set": true,
//     "schedulable_masters": false,
//     "service_networks": [
//       {
//         "cidr": "172.30.0.0/16",
//         "cluster_id": "bfe541fa-9494-4bcc-8c45-3ebba77a7344"
//       }
//     ],
//     "ssh_public_key": "REDACTED",
//     "status": "installing",
//     "status_info": "Installation in progress",
//     "status_updated_at": "2022-06-08T21:09:29.176Z",
//     "total_host_count": 5,
//     "updated_at": "2022-06-08T21:09:29.179951Z",
//     "user_managed_networking": false,
//     "user_name": "admin",
//     "validations_info": "{\"configuration\":[{\"id\":\"pull-secret-set\",\"status\":\"success\",\"message\":\"The pull secret is set.\"}],\"hosts-data\":[{\"id\":\"all-hosts-are-ready-to-install\",\"status\":\"success\",\"message\":\"All hosts in the cluster are ready to install.\"},{\"id\":\"sufficient-masters-count\",\"status\":\"success\",\"message\":\"The cluster has a sufficient number of master candidates.\"}],\"network\":[{\"id\":\"api-vip-defined\",\"status\":\"success\",\"message\":\"The API virtual IP is defined.\"},{\"id\":\"api-vip-valid\",\"status\":\"success\",\"message\":\"api vip 192.168.111.5 belongs to the Machine CIDR and is not in use.\"},{\"id\":\"cluster-cidr-defined\",\"status\":\"success\",\"message\":\"The Cluster Network CIDR is defined.\"},{\"id\":\"dns-domain-defined\",\"status\":\"success\",\"message\":\"The base domain is defined.\"},{\"id\":\"ingress-vip-defined\",\"status\":\"success\",\"message\":\"The Ingress virtual IP is defined.\"},{\"id\":\"ingress-vip-valid\",\"status\":\"success\",\"message\":\"ingress vip 192.168.111.4 belongs to the Machine CIDR and is not in use.\"},{\"id\":\"machine-cidr-defined\",\"status\":\"success\",\"message\":\"The Machine Network CIDR is defined.\"},{\"id\":\"machine-cidr-equals-to-calculated-cidr\",\"status\":\"success\",\"message\":\"The Cluster Machine CIDR is equivalent to the calculated CIDR.\"},{\"id\":\"network-prefix-valid\",\"status\":\"success\",\"message\":\"The Cluster Network prefix is valid.\"},{\"id\":\"network-type-valid\",\"status\":\"success\",\"message\":\"The cluster has a valid network type\"},{\"id\":\"networks-same-address-families\",\"status\":\"success\",\"message\":\"Same address families for all networks.\"},{\"id\":\"no-cidrs-overlapping\",\"status\":\"success\",\"message\":\"No CIDRS are overlapping.\"},{\"id\":\"ntp-server-configured\",\"status\":\"success\",\"message\":\"No ntp problems found\"},{\"id\":\"service-cidr-defined\",\"status\":\"success\",\"message\":\"The Service Network CIDR is defined.\"}],\"operators\":[{\"id\":\"cnv-requirements-satisfied\",\"status\":\"success\",\"message\":\"cnv is disabled\"},{\"id\":\"lso-requirements-satisfied\",\"status\":\"success\",\"message\":\"lso is disabled\"},{\"id\":\"odf-requirements-satisfied\",\"status\":\"success\",\"message\":\"odf is disabled\"}]}",
//     "vip_dhcp_allocation": false
//   }
// ]
// TODO(lranjbar)[AGENT-172]: Need to parse the validations_info object returned by the REST API
// *models.Cluster I expect have a validations_info JSON object to marshal
func (czero *ClusterZero) ParseValidationInfo(*models.Cluster) (bool, error) {

	return false, nil
}
