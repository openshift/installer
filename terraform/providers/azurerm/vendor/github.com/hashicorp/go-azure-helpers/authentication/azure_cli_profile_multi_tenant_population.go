// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package authentication

import (
	"fmt"
	"strings"

	"github.com/Azure/go-autorest/autorest/azure/cli"
)

func (a *azureCLIProfileMultiTenant) populateSubscriptionID() error {
	subscriptionId, err := a.findDefaultSubscriptionId()
	if err != nil {
		return err
	}

	a.subscriptionId = subscriptionId
	return nil
}

func (a *azureCLIProfileMultiTenant) populateTenantID() error {
	subscription, err := a.findSubscription(a.subscriptionId)
	if err != nil {
		return err
	}

	a.tenantId = subscription.TenantID
	return nil
}

func (a *azureCLIProfileMultiTenant) populateEnvironment() error {
	subscription, err := a.findSubscription(a.subscriptionId)
	if err != nil {
		return err
	}

	a.environment = normalizeEnvironmentName(subscription.EnvironmentName)
	return nil
}

func (a azureCLIProfileMultiTenant) findDefaultSubscriptionId() (string, error) {
	for _, subscription := range a.profile.Subscriptions {
		if subscription.IsDefault {
			return subscription.ID, nil
		}
	}

	return "", fmt.Errorf("No Subscription was Marked as Default in the Azure Profile.")
}

func (a azureCLIProfileMultiTenant) findSubscription(subscriptionId string) (*cli.Subscription, error) {
	for _, subscription := range a.profile.Subscriptions {
		if strings.EqualFold(subscription.ID, subscriptionId) {
			return &subscription, nil
		}
	}

	return nil, fmt.Errorf("Subscription %q was not found in your Azure CLI credentials. Please verify it exists in `az account list`.", subscriptionId)
}
