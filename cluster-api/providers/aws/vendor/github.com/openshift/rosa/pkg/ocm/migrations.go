package ocm

import (
	"net"
	"strings"

	cmv1 "github.com/openshift-online/ocm-sdk-go/clustersmgmt/v1"
	errors "github.com/zgalor/weberr"
)

const (
	SubnetConfigTransit    = "transit"
	SubnetConfigMasquerade = "masquerade"
	SubnetConfigJoin       = "join"

	OvnInternalSubnetsFlagName = "ovn-internal-subnets"
	NetworkTypeFlagName        = "network-type"
	NetworkTypeOvn             = "OVNKubernetes"
	NetworkTypeOvnAlias        = "OVN-Kubernetes"
)

func (c *Client) FetchClusterMigrations(clusterId string) (*cmv1.ClusterMigrationsListResponse,
	error) {
	collection := c.ocm.ClustersMgmt().V1().Clusters().Cluster(clusterId)
	resources := collection.Migrations().List()

	clusterMigrations, err := resources.Send()
	if err != nil {
		return nil, errors.UserWrapf(err, "Can't retrieve migrations for cluster '%s'\n", clusterId)
	}

	return clusterMigrations, nil
}

// SDN -> OVN migration key=value parser
func ParseAndValidateOvnInternalSubnets(internalSubnets string) (map[string]string, error) {
	subnetMap := make(map[string]string)
	if internalSubnets == "" || internalSubnets == `""` {
		return nil, errors.UserErrorf("Expected value for '%s'",
			OvnInternalSubnetsFlagName)
	} else {
		possibleConfigs := strings.Split(internalSubnets, ",")
		for i, config := range possibleConfigs {
			// Skip last one if empty, keep the ones leading up to it
			if config == "" && i == len(possibleConfigs)-1 {
				continue
			}
			if !strings.Contains(config, "=") {
				return nil, errors.UserErrorf("Expected key=value format for labels")
			}
			tokens := strings.Split(config, "=")
			// Validate key is correct
			if tokens[0] != SubnetConfigTransit && tokens[0] != SubnetConfigMasquerade && tokens[0] != SubnetConfigJoin {
				return nil, errors.UserErrorf("Expected 'join', 'transit', or "+
					"'masquerade' as keys for '%s'", OvnInternalSubnetsFlagName)
			}
			// Validate value
			_, _, err := net.ParseCIDR(tokens[1])
			if err != nil {
				return nil, errors.UserErrorf("Expected valid CIDR as value "+
					"for '%s'", tokens[0])
			}
			key := strings.TrimSpace(tokens[0])
			value := strings.TrimSpace(tokens[1])
			if _, exists := subnetMap[key]; exists {
				return nil, errors.UserErrorf("Duplicated subnet configuration key "+
					"'%s' used", key)
			}
			subnetMap[key] = value
		}
	}
	return subnetMap, nil
}
