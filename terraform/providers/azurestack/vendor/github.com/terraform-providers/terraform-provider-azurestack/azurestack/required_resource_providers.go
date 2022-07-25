package azurestack

import (
	"context"
	"log"

	"github.com/Azure/azure-sdk-for-go/profiles/2017-03-09/resources/mgmt/resources"
	"github.com/hashicorp/go-azure-helpers/resourceproviders"
)

// requiredResourceProviders returns all of the Resource Providers used by the AzureStack Provider
// whilst all may not be used by every user - the intention is that we determine which should be
// registered such that we can avoid obscure errors where Resource Providers aren't registered.
// new Resource Providers should be added to this list as they're used in the Provider
// (this is the approach used by Microsoft in their tooling)
func requiredResourceProviders() map[string]struct{} {
	// NOTE: Resource Providers in this list are case sensitive
	return map[string]struct{}{
		"Microsoft.Authorization": {},
		"Microsoft.Compute":       {},
		"Microsoft.KeyVault":      {},
		"Microsoft.Network":       {},
		"Microsoft.Storage":       {},
	}
}

func ensureResourceProvidersAreRegistered(ctx context.Context, client resources.ProvidersClient, availableRPs []resources.Provider, requiredRPs map[string]struct{}) error {
	log.Printf("[DEBUG] Determining which Resource Providers require Registration")
	providersToRegister := resourceproviders.DetermineResourceProvidersRequiringRegistration(availableRPs, requiredRPs)

	if len(providersToRegister) > 0 {
		log.Printf("[DEBUG] Registering %d Resource Providers", len(providersToRegister))
		err := resourceproviders.RegisterForSubscription(ctx, client, providersToRegister)
		if err != nil {
			return err
		}
	} else {
		log.Printf("[DEBUG] All required Resource Providers are registered")
	}

	return nil
}
