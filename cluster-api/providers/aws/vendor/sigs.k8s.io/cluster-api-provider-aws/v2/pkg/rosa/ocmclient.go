// Package rosa provides a way to interact with the Red Hat OpenShift Service on AWS (ROSA) API.
package rosa

import (
	"context"

	v1 "github.com/openshift-online/ocm-sdk-go/clustersmgmt/v1"
	"github.com/openshift/rosa/pkg/aws"
	"github.com/openshift/rosa/pkg/ocm"

	"sigs.k8s.io/cluster-api-provider-aws/v2/pkg/cloud/scope"
)

type ocmclient struct {
	ocmClient *ocm.Client
}

// OCMClient wraps ocm.Client methods that we use in interface, so we are able to mock it.
// We should get rid of this once ocm.Client has its own interface.
type OCMClient interface {
	AckVersionGate(clusterID string, gateID string) error
	AddHTPasswdUser(username string, password string, clusterID string, idpID string) error
	CreateNodePool(clusterID string, nodePool *v1.NodePool) (*v1.NodePool, error)
	CreateIdentityProvider(clusterID string, idp *v1.IdentityProvider) (*v1.IdentityProvider, error)
	CreateCluster(config ocm.Spec) (*v1.Cluster, error)
	CreateUser(clusterID string, group string, user *v1.User) (*v1.User, error)
	DeleteCluster(clusterKey string, bestEffort bool, creator *aws.Creator) (*v1.Cluster, error)
	DeleteNodePool(clusterID string, nodePoolID string) error
	DeleteUser(clusterID string, group string, username string) error
	GetCluster(clusterKey string, creator *aws.Creator) (*v1.Cluster, error)
	GetControlPlaneUpgradePolicies(clusterID string) (controlPlaneUpgradePolicies []*v1.ControlPlaneUpgradePolicy, err error)
	GetHTPasswdUserList(clusterID string, htpasswdIDPId string) (*v1.HTPasswdUserList, error)
	GetIdentityProviders(clusterID string) ([]*v1.IdentityProvider, error)
	GetMissingGateAgreementsHypershift(clusterID string, upgradePolicy *v1.ControlPlaneUpgradePolicy) ([]*v1.VersionGate, error)
	GetNodePool(clusterID string, nodePoolID string) (*v1.NodePool, bool, error)
	GetHypershiftNodePoolUpgrade(clusterID string, clusterKey string, nodePoolID string) (*v1.NodePool, *v1.NodePoolUpgradePolicy, error)
	GetUser(clusterID string, group string, username string) (*v1.User, error)
	ScheduleHypershiftControlPlaneUpgrade(clusterID string, upgradePolicy *v1.ControlPlaneUpgradePolicy) (*v1.ControlPlaneUpgradePolicy, error)
	ScheduleNodePoolUpgrade(clusterID string, nodePoolID string, upgradePolicy *v1.NodePoolUpgradePolicy) (*v1.NodePoolUpgradePolicy, error)
	UpdateNodePool(clusterID string, nodePool *v1.NodePool) (*v1.NodePool, error)
	UpdateCluster(clusterKey string, creator *aws.Creator, config ocm.Spec) error
	ValidateHypershiftVersion(versionRawID string, channelGroup string) (bool, error)
}

func (c *ocmclient) AckVersionGate(clusterID string, gateID string) error {
	return c.ocmClient.AckVersionGate(clusterID, gateID)
}

func (c *ocmclient) AddHTPasswdUser(username string, password string, clusterID string, idpID string) error {
	return c.ocmClient.AddHTPasswdUser(username, password, clusterID, idpID)
}
func (c *ocmclient) CreateIdentityProvider(clusterID string, idp *v1.IdentityProvider) (*v1.IdentityProvider, error) {
	return c.ocmClient.CreateIdentityProvider(clusterID, idp)
}
func (c *ocmclient) CreateNodePool(clusterID string, nodePool *v1.NodePool) (*v1.NodePool, error) {
	return c.ocmClient.CreateNodePool(clusterID, nodePool)
}

func (c *ocmclient) CreateCluster(config ocm.Spec) (*v1.Cluster, error) {
	return c.ocmClient.CreateCluster(config)
}
func (c *ocmclient) CreateUser(clusterID string, group string, user *v1.User) (*v1.User, error) {
	return c.ocmClient.CreateUser(clusterID, group, user)
}

func (c *ocmclient) DeleteUser(clusterID string, group string, username string) error {
	return c.ocmClient.DeleteUser(clusterID, group, username)
}

func (c *ocmclient) DeleteNodePool(clusterID string, nodePoolID string) error {
	return c.ocmClient.DeleteNodePool(clusterID, nodePoolID)
}

func (c *ocmclient) DeleteCluster(clusterKey string, bestEffort bool, creator *aws.Creator) (*v1.Cluster, error) {
	return c.ocmClient.DeleteCluster(clusterKey, bestEffort, creator)
}

func (c *ocmclient) GetIdentityProviders(clusterID string) ([]*v1.IdentityProvider, error) {
	return c.ocmClient.GetIdentityProviders(clusterID)
}

func (c *ocmclient) GetControlPlaneUpgradePolicies(clusterID string) (controlPlaneUpgradePolicies []*v1.ControlPlaneUpgradePolicy, err error) {
	return c.ocmClient.GetControlPlaneUpgradePolicies(clusterID)
}

func (c *ocmclient) GetHTPasswdUserList(clusterID string, htpasswdIDPId string) (*v1.HTPasswdUserList, error) {
	return c.ocmClient.GetHTPasswdUserList(clusterID, htpasswdIDPId)
}

func (c *ocmclient) GetMissingGateAgreementsHypershift(clusterID string, upgradePolicy *v1.ControlPlaneUpgradePolicy) ([]*v1.VersionGate, error) {
	return c.ocmClient.GetMissingGateAgreementsHypershift(clusterID, upgradePolicy)
}

func (c *ocmclient) GetNodePool(clusterID string, nodePoolID string) (*v1.NodePool, bool, error) {
	return c.ocmClient.GetNodePool(clusterID, nodePoolID)
}

func (c *ocmclient) GetHypershiftNodePoolUpgrade(clusterID string, clusterKey string, nodePoolID string) (*v1.NodePool, *v1.NodePoolUpgradePolicy, error) {
	return c.ocmClient.GetHypershiftNodePoolUpgrade(clusterID, clusterKey, nodePoolID)
}

func (c *ocmclient) GetCluster(clusterKey string, creator *aws.Creator) (*v1.Cluster, error) {
	return c.ocmClient.GetCluster(clusterKey, creator)
}

func (c *ocmclient) GetUser(clusterID string, group string, username string) (*v1.User, error) {
	return c.ocmClient.GetUser(clusterID, group, username)
}

func (c *ocmclient) ScheduleNodePoolUpgrade(clusterID string, nodePoolID string, upgradePolicy *v1.NodePoolUpgradePolicy) (*v1.NodePoolUpgradePolicy, error) {
	return c.ocmClient.ScheduleNodePoolUpgrade(clusterID, nodePoolID, upgradePolicy)
}

func (c *ocmclient) ScheduleHypershiftControlPlaneUpgrade(clusterID string, upgradePolicy *v1.ControlPlaneUpgradePolicy) (*v1.ControlPlaneUpgradePolicy, error) {
	return c.ocmClient.ScheduleHypershiftControlPlaneUpgrade(clusterID, upgradePolicy)
}

func (c *ocmclient) UpdateCluster(clusterKey string, creator *aws.Creator, config ocm.Spec) error {
	return c.ocmClient.UpdateCluster(clusterKey, creator, config)
}

func (c *ocmclient) UpdateNodePool(clusterID string, nodePool *v1.NodePool) (*v1.NodePool, error) {
	return c.ocmClient.UpdateNodePool(clusterID, nodePool)
}

func (c *ocmclient) ValidateHypershiftVersion(versionRawID string, channelGroup string) (bool, error) {
	return c.ocmClient.ValidateHypershiftVersion(versionRawID, channelGroup)
}

// NewMockOCMClient creates a new empty ocm.Client without any real connection.
func NewMockOCMClient(ctx context.Context, rosaScope *scope.ROSAControlPlaneScope) (OCMClient, error) {
	return &ocmclient{ocmClient: &ocm.Client{}}, nil
}
