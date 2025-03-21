/*
Copyright (c) 2021 Red Hat, Inc.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

  http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package ocm

import (
	"bytes"
	"fmt"

	amsv1 "github.com/openshift-online/ocm-sdk-go/accountsmgmt/v1"
	asv1 "github.com/openshift-online/ocm-sdk-go/addonsmgmt/v1"
	cmv1 "github.com/openshift-online/ocm-sdk-go/clustersmgmt/v1"
)

type AddOnBilling struct {
	BillingModel     string
	BillingAccountID string
}

var BillingOptions = []string{
	string(amsv1.BillingModelMarketplace),
	string(amsv1.BillingModelStandard),
	string(amsv1.BillingModelMarketplaceAWS),
	string(amsv1.BillingModelMarketplaceAzure),
	string(amsv1.BillingModelMarketplaceRHM),
}

var BillingModels = map[string]asv1.BillingModel{
	string(amsv1.BillingModelMarketplace):      asv1.BillingModelMarketplace,
	string(amsv1.BillingModelStandard):         asv1.BillingModelStandard,
	string(amsv1.BillingModelMarketplaceAWS):   asv1.BillingModelMarketplaceAws,
	string(amsv1.BillingModelMarketplaceAzure): asv1.BillingModelMarketplaceAzure,
	string(amsv1.BillingModelMarketplaceRHM):   asv1.BillingModelMarketplaceRhm,
}

type AddOnParam struct {
	Key string
	Val string
}

type AddOnResource struct {
	AddOn     *asv1.Addon
	AZType    string
	Available bool
}

// We customize here the marshalling of AddOnResource to JSON as the embedded AddOn struct does not
// export its member fields. We instead delegate to asv1.MarshalAddon for Marshalling the JSON of the
// embedded AddOn struct and then build a string representation of the struct to be returned as JSON.
func (ar *AddOnResource) MarshalJSON() ([]byte, error) {
	var b bytes.Buffer
	err := asv1.MarshalAddon(ar.AddOn, &b)
	if err != nil {
		return nil, err
	}

	json := fmt.Sprintf(`{"addon": %s, "az_type": "%s", "available":"%t"}`, b.String(), ar.AZType, ar.Available)
	return []byte(json), nil
}

type ClusterAddOn struct {
	ID    string `json:"id"`
	Name  string `json:"name"`
	State string `json:"state"`
}

func (c *Client) InstallAddOn(clusterID, addOnID string, params []AddOnParam, billing AddOnBilling) error {
	addOnInstallationBuilder := asv1.NewAddonInstallation().
		Addon(asv1.NewAddon().ID(addOnID))

	if len(params) > 0 {
		addOnParamList := make([]*asv1.AddonInstallationParameterBuilder, len(params))
		for i, param := range params {
			addOnParamList[i] = asv1.NewAddonInstallationParameter().Id(param.Key).Value(param.Val)
		}
		addOnInstallationBuilder = addOnInstallationBuilder.
			Parameters(asv1.NewAddonInstallationParameterList().Items(addOnParamList...))
	}

	billingModel, exists := BillingModels[billing.BillingModel]
	if !exists {
		return fmt.Errorf("'%s' is not an valid billing model", billing.BillingModel)
	}
	billingBuilder := asv1.NewAddonInstallationBilling().
		BillingModel(billingModel).
		BillingMarketplaceAccount(billing.BillingAccountID)
	addOnInstallationBuilder.Billing(billingBuilder)

	addOnInstallation, err := addOnInstallationBuilder.Build()
	if err != nil {
		return err
	}

	response, err := c.ocm.AddonsMgmt().V1().
		Clusters().
		Cluster(clusterID).
		Addons().
		Add().
		Body(addOnInstallation).
		Send()
	if err != nil {
		return handleErr(response.Error(), err)
	}

	return nil
}

func (c *Client) UninstallAddOn(clusterID, addOnID string) error {
	response, err := c.ocm.AddonsMgmt().V1().
		Clusters().
		Cluster(clusterID).
		Addons().
		Addon(addOnID).
		Delete().
		Send()
	if err != nil {
		return handleErr(response.Error(), err)
	}

	return nil
}

func (c *Client) GetAddOnInstallation(clusterID, addOnID string) (*asv1.AddonInstallation, error) {
	response, err := c.ocm.AddonsMgmt().V1().
		Clusters().
		Cluster(clusterID).
		Addons().
		Addon(addOnID).
		Get().
		Send()
	if err != nil {
		return nil, handleErr(response.Error(), err)
	}

	return response.Body(), nil
}

func (c *Client) UpdateAddOnInstallation(clusterID, addOnID string, params []AddOnParam) error {
	addOnInstallationBuilder := asv1.NewAddonInstallation().
		Addon(asv1.NewAddon().ID(addOnID))

	if len(params) > 0 {
		addOnParamList := make([]*asv1.AddonInstallationParameterBuilder, len(params))
		for i, param := range params {
			addOnParamList[i] = asv1.NewAddonInstallationParameter().Id(param.Key).Value(param.Val)
		}
		addOnInstallationBuilder = addOnInstallationBuilder.
			Parameters(asv1.NewAddonInstallationParameterList().Items(addOnParamList...))
	}

	addOnInstallation, err := addOnInstallationBuilder.Build()
	if err != nil {
		return err
	}

	response, err := c.ocm.AddonsMgmt().V1().Clusters().Cluster(clusterID).
		Addons().Addon(addOnID).
		Update().Body(addOnInstallation).Send()
	if err != nil {
		return handleErr(response.Error(), err)
	}

	return nil
}

func (c *Client) GetAddOnParameters(clusterID, addOnID string) (*asv1.AddonParameterList, error) {
	response, err := c.ocm.AddonsMgmt().V1().Clusters().
		Cluster(clusterID).AddonInquiries().AddonInquiry(addOnID).Get().Send()
	if err != nil {
		return nil, handleErr(response.Error(), err)
	}
	return response.Body().Parameters(), nil
}

// Get complete list of available add-ons for the current organization
func (c *Client) GetAvailableAddOns() ([]*AddOnResource, error) {
	// Get organization ID (used to get add-on quotas)
	acctResponse, err := c.ocm.AccountsMgmt().V1().CurrentAccount().
		Get().
		Send()
	if err != nil {
		return nil, handleErr(acctResponse.Error(), err)
	}
	organization := acctResponse.Body().Organization().ID()

	// Get a list of add-on quotas for the current organization
	quotaCostResponse, err := c.ocm.AccountsMgmt().V1().Organizations().
		Organization(organization).
		QuotaCost().
		List().
		Search("quota_id LIKE 'add-on%'").
		Parameter("fetchRelatedResources", true).
		Page(1).
		Size(-1).
		Send()
	if err != nil {
		return nil, handleErr(quotaCostResponse.Error(), err)
	}
	quotaCosts := quotaCostResponse.Items()

	// Get complete list of enabled add-ons
	addOnsResponse, err := c.ocm.AddonsMgmt().V1().Addons().
		List().
		Search("enabled='t'").
		Page(1).
		Size(-1).
		Send()
	if err != nil {
		return nil, handleErr(addOnsResponse.Error(), err)
	}

	var addOns []*AddOnResource

	// Populate enabled add-ons with if they are available for the current org
	addOnsResponse.Items().Each(func(addOn *asv1.Addon) bool {
		addOnResource := &AddOnResource{
			AddOn: addOn,
		}
		// Free add-ons are always available
		available := addOn.ResourceCost() == 0

		// Only return add-ons for which the org has quota
		quotaCosts.Each(func(quotaCost *amsv1.QuotaCost) bool {
			// Check all related resources to ensure we're checking the product of the correct addon
			for _, relatedResource := range quotaCost.RelatedResources() {
				// Only return compatible addons
				if addOn.ResourceName() == relatedResource.ResourceName() && isCompatible(relatedResource) {
					available = true

					// Addon is only available if quota allows it
					addOnResource.Available = quotaCost.Allowed()-quotaCost.Consumed() >= relatedResource.Cost()

					// Track AZ type so that we can compare against cluster
					addOnResource.AZType = relatedResource.AvailabilityZoneType()
					// Since add-on is considered available now, there's no need to check the other resources
					return false
				}
			}
			return true
		})

		// Only display add-ons that meet the above criteria
		if available {
			addOns = append(addOns, addOnResource)
		}

		return true
	})

	return addOns, nil
}

func (c *Client) GetAddOn(id string) (*asv1.Addon, error) {
	response, err := c.ocm.AddonsMgmt().V1().Addons().Addon(id).Get().Send()
	if err != nil {
		return nil, handleErr(response.Error(), err)
	}
	return response.Body(), nil
}

// Get all add-ons available for a cluster
func (c *Client) GetClusterAddOns(cluster *cmv1.Cluster) ([]*ClusterAddOn, error) {
	addOnResources, err := c.GetAvailableAddOns()
	if err != nil {
		return nil, err
	}

	// Get add-ons already installed on cluster
	addOnInstallationsResponse, err := c.ocm.AddonsMgmt().V1().Clusters().
		Cluster(cluster.ID()).
		Addons().
		List().
		Page(1).
		Size(-1).
		Send()
	if err != nil {
		return nil, handleErr(addOnInstallationsResponse.Error(), err)
	}
	addOnInstallations := addOnInstallationsResponse.Items()

	var clusterAddOns []*ClusterAddOn

	// Populate add-on installations with all add-on metadata
	for _, addOnResource := range addOnResources {
		// Ensure add-on is compatible with the cluster's availability zones
		if !(addOnResource.AZType == ANY ||
			(cluster.MultiAZ() && addOnResource.AZType == "multi") ||
			(!cluster.MultiAZ() && addOnResource.AZType == "single")) {
			continue
		}
		clusterAddOn := ClusterAddOn{
			ID:    addOnResource.AddOn.ID(),
			Name:  addOnResource.AddOn.Name(),
			State: "not installed",
		}
		if !addOnResource.Available {
			clusterAddOn.State = "unavailable"
		}

		// Get the state of add-on installations on the cluster
		addOnInstallations.Each(func(addOnInstallation *asv1.AddonInstallation) bool {
			if addOnResource.AddOn.ID() == addOnInstallation.Addon().ID() {
				clusterAddOn.State = string(addOnInstallation.State())
				if clusterAddOn.State == "" {
					clusterAddOn.State = string(asv1.AddonInstallationStateInstalling)
				}
			}
			return true
		})

		clusterAddOns = append(clusterAddOns, &clusterAddOn)
	}

	return clusterAddOns, nil
}

func (c *Client) AddClusterOperatorRole(cluster *cmv1.Cluster, role *cmv1.OperatorIAMRole) error {
	// Make sure the role doesn't exist already, to avoid conflicts
	operatorRoles := cluster.AWS().STS().OperatorIAMRoles()
	for _, item := range operatorRoles {
		if role.Name() == item.Name() &&
			role.Namespace() == item.Namespace() &&
			role.RoleARN() == item.RoleARN() {
			return nil
		}
	}

	response, err := c.ocm.ClustersMgmt().V1().
		Clusters().
		Cluster(cluster.ID()).
		STSOperatorRoles().
		Add().
		Body(role).
		Send()
	if err != nil {
		return handleErr(response.Error(), err)
	}

	return nil
}
