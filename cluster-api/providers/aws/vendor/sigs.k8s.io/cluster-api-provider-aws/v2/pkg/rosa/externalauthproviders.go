package rosa

import (
	"context"
	"fmt"
	"time"

	sdk "github.com/openshift-online/ocm-sdk-go"
	cmv1 "github.com/openshift-online/ocm-sdk-go/clustersmgmt/v1"

	"sigs.k8s.io/cluster-api-provider-aws/v2/pkg/cloud/scope"
)

// ExternalAuthClient handles externalAuth operations.
type ExternalAuthClient struct {
	ocm *sdk.Connection
}

// NewExternalAuthClient creates and return a new client to handle externalAuth operations.
func NewExternalAuthClient(ctx context.Context, rosaScope *scope.ROSAControlPlaneScope) (*ExternalAuthClient, error) {
	ocmConnection, err := newOCMRawConnection(ctx, rosaScope)
	if err != nil {
		return nil, err
	}
	return &ExternalAuthClient{
		ocm: ocmConnection,
	}, nil
}

// Close closes the underlying ocm connection.
func (c *ExternalAuthClient) Close() error {
	return c.ocm.Close()
}

// CreateExternalAuth creates a new external auth porivder.
func (c *ExternalAuthClient) CreateExternalAuth(clusterID string, externalAuth *cmv1.ExternalAuth) (*cmv1.ExternalAuth, error) {
	response, err := c.ocm.ClustersMgmt().V1().
		Clusters().Cluster(clusterID).
		ExternalAuthConfig().ExternalAuths().Add().Body(externalAuth).Send()
	if err != nil {
		return nil, handleErr(response.Error(), err)
	}
	return response.Body(), nil
}

// UpdateExternalAuth updates an existing external auth porivder.
func (c *ExternalAuthClient) UpdateExternalAuth(clusterID string, externalAuth *cmv1.ExternalAuth) (*cmv1.ExternalAuth, error) {
	response, err := c.ocm.ClustersMgmt().V1().
		Clusters().Cluster(clusterID).
		ExternalAuthConfig().ExternalAuths().
		ExternalAuth(externalAuth.ID()).
		Update().Body(externalAuth).Send()
	if err != nil {
		return nil, handleErr(response.Error(), err)
	}
	return response.Body(), nil
}

// GetExternalAuth retrieves the specified external auth porivder.
func (c *ExternalAuthClient) GetExternalAuth(clusterID string, externalAuthID string) (*cmv1.ExternalAuth, bool, error) {
	response, err := c.ocm.ClustersMgmt().V1().
		Clusters().Cluster(clusterID).ExternalAuthConfig().
		ExternalAuths().ExternalAuth(externalAuthID).
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

// ListExternalAuths lists all external auth porivder for the cluster.
func (c *ExternalAuthClient) ListExternalAuths(clusterID string) ([]*cmv1.ExternalAuth, error) {
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

// DeleteExternalAuth deletes the specified external auth porivder.
func (c *ExternalAuthClient) DeleteExternalAuth(clusterID string, externalAuthID string) error {
	response, err := c.ocm.ClustersMgmt().V1().
		Clusters().Cluster(clusterID).
		ExternalAuthConfig().ExternalAuths().
		ExternalAuth(externalAuthID).
		Delete().
		Send()
	if err != nil {
		return handleErr(response.Error(), err)
	}
	return nil
}

// CreateBreakGlassCredential creates a break glass credential.
func (c *ExternalAuthClient) CreateBreakGlassCredential(clusterID string, breakGlassCredential *cmv1.BreakGlassCredential) (*cmv1.BreakGlassCredential, error) {
	response, err := c.ocm.ClustersMgmt().V1().
		Clusters().Cluster(clusterID).BreakGlassCredentials().
		Add().Body(breakGlassCredential).Send()
	if err != nil {
		return nil, handleErr(response.Error(), err)
	}
	return response.Body(), nil
}

const pollInterval = 15 * time.Second

// PollKubeconfig continuously polls for the kubeconfig of the provided break glass credential.
func (c *ExternalAuthClient) PollKubeconfig(ctx context.Context, clusterID string, credentialID string) (kubeconfig string, err error) {
	ctx, cancel := context.WithTimeout(ctx, time.Minute*5)
	defer cancel()

	credentialClient := c.ocm.ClustersMgmt().V1().Clusters().
		Cluster(clusterID).BreakGlassCredentials().BreakGlassCredential(credentialID)
	response, err := credentialClient.Poll().
		Interval(pollInterval).
		Predicate(func(bgcgr *cmv1.BreakGlassCredentialGetResponse) bool {
			return bgcgr.Body().Status() == cmv1.BreakGlassCredentialStatusIssued && bgcgr.Body().Kubeconfig() != ""
		}).
		StartContext(ctx)
	if err != nil {
		err = fmt.Errorf("failed to poll kubeconfig for cluster '%s' with break glass credential '%s': %v",
			clusterID, credentialID, err)
		return
	}

	return response.Body().Kubeconfig(), nil
}
