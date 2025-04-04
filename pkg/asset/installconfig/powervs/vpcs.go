package powervs

import (
	"context"
	"fmt"
	"net/http"

	survey "github.com/AlecAivazis/survey/v2"
	"github.com/IBM-Cloud/bluemix-go/crn"
	"github.com/IBM/go-sdk-core/v5/core"
	"github.com/IBM/vpc-go-sdk/vpcv1"
)

// GetUserSelectedPermittedNetwork gets the user-preferred permitted network for the DNS zone for the deployment.
func GetUserSelectedPermittedNetwork(zoneID string, dnsCRN crn.CRN, region string) (string, error) {
	var (
		detailedResponse       *core.DetailedResponse
		vpcChoice              string
		permittedNetworkCRNs   []string
		permittedNetworkCRN    string
		permittedNetworkCRNCRN crn.CRN
		err                    error
		client                 *Client
		vpcOptions             *vpcv1.GetVPCOptions
		vpc                    *vpcv1.VPC
	)
	vpcNames := []string{}
	client, err = NewClient()
	if err != nil {
		return "", err
	}
	permittedNetworkCRNs, err = client.GetDNSInstancePermittedNetworks(context.TODO(), dnsCRN.ServiceInstance, zoneID)
	if err != nil {
		return "", err
	}
	for _, permittedNetworkCRN = range permittedNetworkCRNs {
		permittedNetworkCRNCRN, err = crn.Parse(permittedNetworkCRN)
		if err != nil {
			return "", fmt.Errorf("failed to parse DNSInstanceCRN: %w", err)
		}
		vpcOptions = client.vpcAPI.NewGetVPCOptions(permittedNetworkCRNCRN.Resource)
		if vpc, detailedResponse, err = client.vpcAPI.GetVPC(vpcOptions); err != nil {
			if detailedResponse.GetStatusCode() != http.StatusNotFound {
				return "", err
			}
		} else if vpc != nil {
			vpcNames = append(vpcNames, *vpc.Name)
		}
	}
	err = survey.Ask([]*survey.Question{
		{
			Prompt: &survey.Select{
				Message: "Permitted Network",
				Help:    "The VPC of the cluster. If you don't see your intended VPC listed, add the VPC to Permitted Networks of the DNS Zone and rerun the installer.",
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
