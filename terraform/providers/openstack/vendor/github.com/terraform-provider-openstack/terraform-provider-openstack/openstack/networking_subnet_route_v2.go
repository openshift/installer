package openstack

import (
	"fmt"
	"strings"
)

func resourceNetworkingSubnetRouteV2BuildID(subnetID, dstCIDR, nextHop string) string {
	return fmt.Sprintf("%s-route-%s-%s", subnetID, dstCIDR, nextHop)
}

func resourceNetworkingSubnetRouteV2ParseID(routeID string) (string, string, string, error) {
	routeIDAllParts := strings.Split(routeID, "-route-")
	if len(routeIDAllParts) != 2 {
		return "", "", "", fmt.Errorf("invalid ID format: %s", routeID)
	}

	routeIDLastPart := routeIDAllParts[1]
	routeIDLastParts := strings.Split(routeIDLastPart, "-")
	if len(routeIDLastParts) != 2 {
		return "", "", "", fmt.Errorf("invalid last part format for %s: %s", routeID, routeIDLastPart)
	}

	return routeIDAllParts[0], routeIDLastParts[0], routeIDLastParts[1], nil
}
