package azure

import (
	"context"
	"fmt"
	"time"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/arm"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/authorization/armauthorization/v3"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/msi/armmsi"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"k8s.io/utils/ptr"

	"github.com/openshift/installer/pkg/types"
)

type identityInput struct {
	installConfig *types.InstallConfig

	region            string
	resourceGroupName string
	subscriptionID    string

	tokenCredential azcore.TokenCredential
	infraID         string
	clientOpts      *arm.ClientOptions

	tags map[string]*string
}

// handleIdentity checks if a user-assigned identity is needed
// and creates one if appropriate.
func handleIdentity(ctx context.Context, in identityInput) error {
	if in.installConfig.CreateAzureIdentity() {
		return createUserAssignedIdentity(ctx, in)
	}

	return nil
}

func createUserAssignedIdentity(ctx context.Context, in identityInput) error {
	userAssignedIdentityName := fmt.Sprintf("%s-identity", in.infraID)
	armmsiClientFactory, err := armmsi.NewClientFactory(
		in.subscriptionID,
		in.tokenCredential,
		in.clientOpts,
	)
	if err != nil {
		return fmt.Errorf("failed to create armmsi client: %w", err)
	}
	_, err = armmsiClientFactory.NewUserAssignedIdentitiesClient().CreateOrUpdate(
		ctx,
		in.resourceGroupName,
		userAssignedIdentityName,
		armmsi.Identity{
			Location: ptr.To(in.region),
			Tags:     in.tags,
		},
		nil,
	)
	if err != nil {
		return fmt.Errorf("failed to create user assigned identity %s: %w", userAssignedIdentityName, err)
	}
	userAssignedIdentity, err := armmsiClientFactory.NewUserAssignedIdentitiesClient().Get(
		ctx,
		in.resourceGroupName,
		userAssignedIdentityName,
		nil,
	)
	if err != nil {
		return fmt.Errorf("failed to get user assigned identity %s: %w", userAssignedIdentityName, err)
	}
	principalID := *userAssignedIdentity.Properties.PrincipalID

	logrus.Debugf("UserAssignedIdentity.ID=%s", *userAssignedIdentity.ID)
	logrus.Debugf("PrinciapalID=%s", principalID)

	clientFactory, err := armauthorization.NewClientFactory(
		in.subscriptionID,
		in.tokenCredential,
		in.clientOpts,
	)
	if err != nil {
		return fmt.Errorf("failed to create armauthorization client: %w", err)
	}

	roleDefinitionsClient := clientFactory.NewRoleDefinitionsClient()

	var contributor *armauthorization.RoleDefinition
	rgScope := fmt.Sprintf("/subscriptions/%s/resourceGroups/%s", in.subscriptionID, in.resourceGroupName)
	roleDefinitionsPager := roleDefinitionsClient.NewListPager(rgScope, nil)
	for roleDefinitionsPager.More() {
		roleDefinitionsList, err := roleDefinitionsPager.NextPage(ctx)
		if err != nil {
			return fmt.Errorf("failed to find any role definitions: %w", err)
		}
		for _, roleDefinition := range roleDefinitionsList.Value {
			if *roleDefinition.Properties.RoleName == "Contributor" {
				contributor = roleDefinition
				break
			}
		}
	}
	if contributor == nil {
		return fmt.Errorf("failed to find contributor definition")
	}

	roleAssignmentsClient := clientFactory.NewRoleAssignmentsClient()
	roleAssignmentUUID := uuid.New().String()

	// XXX: Azure doesn't like creating an identity and immediately
	// creating a role assignment for the identity. There can be
	// replication delays. So, retry every 10 seconds for a minute until
	// the role assignment gets created.
	//
	// See https://aka.ms/docs-principaltype
	for i := 0; i < retryCount; i++ {
		_, err = roleAssignmentsClient.Create(ctx, rgScope, roleAssignmentUUID,
			armauthorization.RoleAssignmentCreateParameters{
				Properties: &armauthorization.RoleAssignmentProperties{
					PrincipalID:      ptr.To(principalID),
					RoleDefinitionID: contributor.ID,
				},
			},
			nil,
		)
		if err == nil {
			break
		}
		time.Sleep(retryTime)
	}
	if err != nil {
		return fmt.Errorf("failed to create role assignment: %w", err)
	}
	return nil
}
