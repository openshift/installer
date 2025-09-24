package powervs

import (
	"fmt"

	survey "github.com/AlecAivazis/survey/v2"
	"github.com/IBM/go-sdk-core/v5/core"
	"github.com/IBM/vpc-go-sdk/vpcv1"
)

// GetUserSelectedPermittedNetwork gets the user-preferred permitted network for the DNS zone for the deployment.
func GetUserSelectedPermittedNetwork(resourceGroup string) (string, error) {
	var (
		client          *Client
		err             error
		options         *vpcv1.ListVpcsOptions
		resourceGroupID string
		response        *core.DetailedResponse
		vpc             vpcv1.VPC
		vpcChoice       string
		vpcCollection   *vpcv1.VPCCollection
	)
	client, err = NewClient()
	if err != nil {
		return "", err
	}
	vpcNames := []string{}
	listGroupOptions := client.managementAPI.NewListResourceGroupsOptions()
	groups, _, err := client.managementAPI.ListResourceGroups(listGroupOptions)
	if err != nil {
		return "", fmt.Errorf("failed to list resource groups: %w", err)
	}
	for _, group := range groups.Resources {
		if *group.Name == resourceGroup {
			resourceGroupID = *group.ID
		}
	}
	options = client.vpcAPI.NewListVpcsOptions()
	options.SetResourceGroupID(resourceGroupID)
	vpcCollection, response, err = client.vpcAPI.ListVpcs(options)
	if err != nil {
		return "", fmt.Errorf("failed to list vps: err = %w, response = %v", err, response)
	}
	for _, vpc = range vpcCollection.Vpcs {
		vpcNames = append(vpcNames, *vpc.Name)
	}
	err = survey.Ask([]*survey.Question{
		{
			Prompt: &survey.Select{
				Message: "Cluster VPC",
				Help:    "The VPC of the cluster. If you don't see your intended VPC listed, check whether the VPC belongs to the resource group selected and rerun the installer.",
				Default: "",
				Options: vpcNames,
			},
		},
	}, &vpcChoice)
	if err != nil {
		return "", fmt.Errorf("survey.ask failed with: %w", err)
	}
	return vpcChoice, nil
}
