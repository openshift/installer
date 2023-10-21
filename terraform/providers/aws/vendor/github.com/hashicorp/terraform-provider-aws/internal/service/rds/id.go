package rds

import (
	"fmt"
	"strings"
)

func ProxyEndpointParseID(id string) (string, string, error) {
	idParts := strings.SplitN(id, "/", 2)
	if len(idParts) != 2 || idParts[0] == "" || idParts[1] == "" {
		return "", "", fmt.Errorf("unexpected format of ID (%s), expected db_proxy_name/db_proxy_endpoint_name", id)
	}
	return idParts[0], idParts[1], nil
}

const clusterRoleAssociationResourceIDSeparator = ","

func ClusterRoleAssociationCreateResourceID(dbClusterID, roleARN string) string {
	parts := []string{dbClusterID, roleARN}
	id := strings.Join(parts, clusterRoleAssociationResourceIDSeparator)

	return id
}

func ClusterRoleAssociationParseResourceID(id string) (string, string, error) {
	parts := strings.Split(id, clusterRoleAssociationResourceIDSeparator)

	if len(parts) == 2 && parts[0] != "" && parts[1] != "" {
		return parts[0], parts[1], nil
	}

	return "", "", fmt.Errorf("unexpected format for ID (%[1]s), expected DBCLUSTERID%[2]sROLEARN", id, clusterRoleAssociationResourceIDSeparator)
}
