/*
Copyright (c) 2020 Red Hat, Inc.

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
	"errors"
	"fmt"
	"strings"

	v1 "github.com/openshift-online/ocm-api-model/clientapi/clustersmgmt/v1"
	amsv1 "github.com/openshift-online/ocm-sdk-go/accountsmgmt/v1"
	cmv1 "github.com/openshift-online/ocm-sdk-go/clustersmgmt/v1"

	"github.com/openshift/rosa/pkg/aws"
)

const AcceleratedComputing = "accelerated_computing"

func (c *Client) GetMachineTypesInRegion(cloudProviderData *cmv1.CloudProviderData) (MachineTypeList, error) {
	collection := c.ocm.ClustersMgmt().V1().AWSInquiries().MachineTypes()
	page := 1
	size := 100

	var machineTypes MachineTypeList
	for {
		response, err := collection.Search().
			Parameter("order", "category asc").
			Body(cloudProviderData).
			Page(page).
			Size(size).
			Send()
		if err != nil {
			return MachineTypeList{}, err
		}

		response.Items().Each(func(item *cmv1.MachineType) bool {
			machineTypes.Items = append(machineTypes.Items, &MachineType{
				MachineType: item,
			})
			return true
		})

		if response.Size() < size {
			break
		}
		page++
	}

	return machineTypes, nil
}

func (c *Client) GetMachineTypes() (machineTypes MachineTypeList, err error) {
	collection := c.ocm.ClustersMgmt().V1().MachineTypes()
	page := 1
	size := 100
	for {
		var response *cmv1.MachineTypesListResponse
		response, err := collection.List().
			Search("cloud_provider.id = 'aws'").
			Order("category asc").
			Page(page).
			Size(size).
			Send()
		if err != nil {
			errMsg := response.Error().Reason()
			if errMsg == "" {
				errMsg = err.Error()
			}
			return MachineTypeList{}, errors.New(errMsg)
		}

		response.Items().Each(func(item *cmv1.MachineType) bool {
			machineTypes.Items = append(machineTypes.Items, &MachineType{
				MachineType: item,
			})
			return true
		})

		if response.Size() < size {
			break
		}
		page++
	}

	return
}

func getDefaultNodes(multiAZ bool) int {
	minimumNodes := 2
	if multiAZ {
		minimumNodes = 3
	}
	return minimumNodes
}

type MachineType struct {
	MachineType    *cmv1.MachineType
	Available      bool
	availableQuota int
}

func (mt MachineType) HasQuota(multiAZ bool) bool {
	return mt.MachineType.Category() != AcceleratedComputing || mt.availableQuota > getDefaultNodes(multiAZ)
}

// GetAvailableMachineTypesInRegion get the supported machine type in the region.
// The function triggers the 'api/clusters_mgmt/v1/aws_inquiries/machine_types'
// and passes a role ARN for STS clusters or access keys for non-STS clusters.
func (c *Client) GetAvailableMachineTypesInRegion(region string, availabilityZones []string, roleARN string,
	awsClient aws.Client, externalId string) (MachineTypeList, error) {
	cloudProviderDataBuilder, err := c.createCloudProviderDataBuilder(roleARN, awsClient, externalId)
	if err != nil {
		return MachineTypeList{}, err
	}
	if len(availabilityZones) > 0 {
		cloudProviderDataBuilder = cloudProviderDataBuilder.AvailabilityZones(availabilityZones...)
	}
	cloudProviderData, err := cloudProviderDataBuilder.Region(cmv1.NewCloudRegion().ID(region)).Build()
	if err != nil {
		return MachineTypeList{}, err
	}

	machineTypes, err := c.GetMachineTypesInRegion(cloudProviderData)
	if err != nil {
		return MachineTypeList{}, err
	}

	machineTypes.Region = region
	machineTypes.AvailabilityZones = availabilityZones

	quotaCosts, err := c.getQuotaCosts()
	if err != nil {
		return MachineTypeList{}, err
	}

	machineTypes.UpdateAvailableQuota(quotaCosts)
	return machineTypes, nil
}

func (c *Client) GetAvailableMachineTypes() (MachineTypeList, error) {
	machineTypes, err := c.GetMachineTypes()
	if err != nil {
		return MachineTypeList{}, err
	}

	quotaCosts, err := c.getQuotaCosts()
	if err != nil {
		return MachineTypeList{}, err
	}

	machineTypes.UpdateAvailableQuota(quotaCosts)
	return machineTypes, nil
}

func (c *Client) getQuotaCosts() (*amsv1.QuotaCostList, error) {
	acctResponse, err := c.ocm.AccountsMgmt().V1().CurrentAccount().
		Get().
		Send()
	if err != nil {
		return nil, handleErr(acctResponse.Error(), err)
	}
	organization := acctResponse.Body().Organization().ID()
	quotaCostResponse, err := c.ocm.AccountsMgmt().V1().Organizations().
		Organization(organization).
		QuotaCost().
		List().
		Parameter("fetchRelatedResources", true).
		Parameter("search", "quota_id~='gpu'").
		Page(1).
		Size(-1).
		Send()
	if err != nil {
		return nil, handleErr(quotaCostResponse.Error(), err)
	}
	quotaCosts := quotaCostResponse.Items()
	return quotaCosts, nil
}

// A list of MachineTypes with additional information
type MachineTypeList struct {
	Items             []*MachineType
	Region            string
	AvailabilityZones []string
}

// IDs extracts list of IDs from a MachineTypeList
func (mtl *MachineTypeList) IDs() []string {
	res := make([]string, len(mtl.Items))
	for i, v := range mtl.Items {
		res[i] = v.MachineType.ID()
	}
	return res
}

// Find returns the first MachineType matching the ID
func (mtl *MachineTypeList) Find(id string) *MachineType {
	for _, v := range mtl.Items {
		if v.MachineType.ID() == id {
			return v
		}
	}
	return nil
}

// Filter returns a new MachineTypeList with only elements for which fn returned true
func (mtl *MachineTypeList) Filter(fn func(*MachineType) bool) MachineTypeList {
	var res MachineTypeList
	for _, v := range mtl.Items {
		if fn(v) {
			res.Items = append(res.Items, v)
		}
	}
	return res
}

func (mtl *MachineTypeList) UpdateAvailableQuota(quotaCosts *amsv1.QuotaCostList) {
	for _, machineType := range mtl.Items {
		if machineType.MachineType.Category() != AcceleratedComputing {
			machineType.Available = true
			continue
		}
		quotaCosts.Each(func(quotaCost *amsv1.QuotaCost) bool {
			for _, relatedResource := range quotaCost.RelatedResources() {
				if machineType.MachineType.GenericName() == relatedResource.ResourceName() &&
					isCompatible(relatedResource) {
					availableQuota := (quotaCost.Allowed() - quotaCost.Consumed()) / relatedResource.Cost()
					machineType.Available = availableQuota > 1
					machineType.availableQuota = availableQuota
					return false
				}
			}
			return true
		})
	}
}

func (mtl *MachineTypeList) GetAvailableIDs(multiAZ bool) *MachineTypeList {
	list := mtl.Filter(func(mt *MachineType) bool {
		return mt.Available && mt.HasQuota(multiAZ)
	})
	return &list
}

func (mtl *MachineTypeList) GetWinLi(winLi string) *MachineTypeList {
	if !strings.EqualFold(winLi, string(v1.ImageTypeWindows)) {
		return mtl
	}

	list := mtl.Filter(func(mt *MachineType) bool {
		if mt.MachineType == nil {
			return false
		}
		features, ok := mt.MachineType.GetFeatures()
		if !ok {
			return false
		}
		winLi, ok := features.GetWinLI()
		if ok {
			return winLi
		}
		return false
	})
	return &list
}

// Validate AWS machine type is available with enough quota in the list
func (mtl *MachineTypeList) ValidateMachineType(machineType string, multiAZ bool) error {
	if machineType == "" {
		return nil
	}
	v := mtl.Find(machineType)

	if v == nil {
		return nil // Replaced not-found error with a preflight in CS (validateZoneSupportInstanceType)
	}

	if !v.HasQuota(multiAZ) {
		err := fmt.Errorf("insufficient quota for instance type: %s", machineType)
		return err
	}

	return nil
}
