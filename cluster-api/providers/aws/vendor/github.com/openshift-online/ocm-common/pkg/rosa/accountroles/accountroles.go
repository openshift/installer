package accountroles

import (
	awsUtils "github.com/openshift-online/ocm-common/pkg/aws/utils"
	cmv1 "github.com/openshift-online/ocm-sdk-go/clustersmgmt/v1"
)

const (
	InstallerAccountRole    = "installer"
	ControlPlaneAccountRole = "instance_controlplane"
	WorkerAccountRole       = "instance_worker"
	SupportAccountRole      = "support"
)

type AccountRole struct {
	Name string
	Flag string
}

var AccountRoles = map[string]AccountRole{
	InstallerAccountRole:    {Name: "Installer", Flag: "role-arn"},
	ControlPlaneAccountRole: {Name: "ControlPlane", Flag: "controlplane-iam-role"},
	WorkerAccountRole:       {Name: "Worker", Flag: "worker-iam-role"},
	SupportAccountRole:      {Name: "Support", Flag: "support-role-arn"},
}

func GetPathFromAccountRole(cluster *cmv1.Cluster, roleNameSuffix string) (string, error) {
	accRoles := GetAccountRolesArnsMap(cluster)
	if accRoles[roleNameSuffix] == "" {
		return "", nil
	}
	return awsUtils.GetPathFromArn(accRoles[roleNameSuffix])
}

func GetAccountRolesArnsMap(cluster *cmv1.Cluster) map[string]string {
	return map[string]string{
		AccountRoles[InstallerAccountRole].Name:    cluster.AWS().STS().RoleARN(),
		AccountRoles[SupportAccountRole].Name:      cluster.AWS().STS().SupportRoleARN(),
		AccountRoles[ControlPlaneAccountRole].Name: cluster.AWS().STS().InstanceIAMRoles().MasterRoleARN(),
		AccountRoles[WorkerAccountRole].Name:       cluster.AWS().STS().InstanceIAMRoles().WorkerRoleARN(),
	}
}
