package operatorroles

import (
	cmv1 "github.com/openshift-online/ocm-sdk-go/clustersmgmt/v1"
)

func GetOperatorRolesArnsMap(cluster *cmv1.Cluster) map[string]string {
	operatorRolesMap := map[string]string{}
	for _, role := range cluster.AWS().STS().OperatorIAMRoles() {
		operatorRolesMap[role.Name()] = role.RoleARN()
	}
	return operatorRolesMap
}
