// Package nutanix collects Nutanix-specific configuration.
package nutanix

import (
	"context"
	"fmt"
	"sort"
	"time"

	"github.com/AlecAivazis/survey/v2"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	nutanixclientv3 "github.com/terraform-providers/terraform-provider-nutanix/client/v3"

	"github.com/openshift/installer/pkg/types/nutanix"
	nutanixtypes "github.com/openshift/installer/pkg/types/nutanix"
	"github.com/openshift/installer/pkg/validate"
)

// PrismCentralClient wraps a Nutanix V3 client
type PrismCentralClient struct {
	PrismCentral string
	Username     string
	Password     string
	Port         string
	V3Client     *nutanixclientv3.Client
}

// Platform collects Nutanix-specific configuration.
func Platform() (*nutanix.Platform, error) {
	nutanixClient, err := getClients()
	if err != nil {
		return nil, err
	}

	ctx := context.TODO()
	v3Client := nutanixClient.V3Client
	peUUID, err := getPrismElement(ctx, v3Client)
	if err != nil {
		return nil, err
	}

	subnetUUID, err := getSubnet(ctx, v3Client, peUUID)
	if err != nil {
		return nil, err
	}

	apiVIP, ingressVIP, err := getVIPs()
	if err != nil {
		return nil, errors.Wrap(err, "failed to get VIPs")
	}

	platform := &nutanix.Platform{
		PrismCentral:     nutanixClient.PrismCentral,
		Port:             nutanixClient.Port,
		Username:         nutanixClient.Username,
		Password:         nutanixClient.Password,
		PrismElementUUID: peUUID,
		SubnetUUID:       subnetUUID,
		APIVIP:           apiVIP,
		IngressVIP:       ingressVIP,
	}
	return platform, nil

}

// getClients() surveys the user for username, password, port & prism central.
// Validation on the three fields is performed by creating a client.
// If creating the client fails, an error is returned.
func getClients() (*PrismCentralClient, error) {
	var prismCentral, username, password, port string
	if err := survey.Ask([]*survey.Question{
		{
			Prompt: &survey.Input{
				Message: "Prism Central",
				Help:    "The domain name or IP address of the Prism Central to be used for installation.",
			},
			Validate: survey.ComposeValidators(survey.Required, func(ans interface{}) error {
				return validate.Host(ans.(string))
			}),
		},
	}, &prismCentral); err != nil {
		return nil, errors.Wrap(err, "failed UserInput")
	}

	if err := survey.Ask([]*survey.Question{
		{
			Prompt: &survey.Input{
				Message: "Port",
				Help:    "The port used to login to Prism Central.",
				Default: "9440",
			},
			Validate: survey.Required,
		},
	}, &port); err != nil {
		return nil, errors.Wrap(err, "failed UserInput")
	}

	if err := survey.Ask([]*survey.Question{
		{
			Prompt: &survey.Input{
				Message: "Username",
				Help:    "The username to login to Prism Central.",
			},
			Validate: survey.Required,
		},
	}, &username); err != nil {
		return nil, errors.Wrap(err, "failed UserInput")
	}

	if err := survey.Ask([]*survey.Question{
		{
			Prompt: &survey.Password{
				Message: "Password",
				Help:    "The password to login to Prism Central.",
			},
			Validate: survey.Required,
		},
	}, &password); err != nil {
		return nil, errors.Wrap(err, "failed UserInput")
	}

	// There is a noticeable delay when creating the client, so let the user know what's going on.
	logrus.Infof("Connecting to Prism Central %s", prismCentral)
	clientV3, err := nutanixtypes.CreateNutanixClient(context.TODO(),
		prismCentral,
		port,
		username,
		password,
	)
	if err != nil {
		return nil, errors.Wrapf(err, "unable to connect to Prism Central %s. Ensure provided information is correct", prismCentral)
	}

	return &PrismCentralClient{
		PrismCentral: prismCentral,
		Username:     username,
		Password:     password,
		Port:         port,
		V3Client:     clientV3,
	}, nil
}

func getPrismElement(ctx context.Context, client *nutanixclientv3.Client) (string, error) {
	ctx, cancel := context.WithTimeout(ctx, 60*time.Second)
	defer cancel()

	emptyFilter := ""
	allPrismElements, err := client.V3.ListAllCluster(emptyFilter)
	if err != nil {
		return "", errors.Wrap(err, "unable to list prism element clusters")
	}

	pes := allPrismElements.Entities
	if len(pes) == 0 {
		return "", errors.New("did not find any prism element clusters")
	} else if len(pes) == 1 {
		peName := *pes[0].Spec.Name
		peUUID := *pes[0].Metadata.UUID
		logrus.Infof("Defaulting to only available prism element: %s", peName)
		return peUUID, nil
	}

	peUUIDs := make(map[string]string)
	var peChoices []string
	for _, p := range pes {
		n := *p.Spec.Name
		peUUIDs[n] = *p.Metadata.UUID
		peChoices = append(peChoices, n)
	}

	var selectedPe string
	if err := survey.Ask([]*survey.Question{
		{
			Prompt: &survey.Select{
				Message: "Prism Element",
				Options: peChoices,
				Help:    "The Prism Element to be used for installation.",
			},
			Validate: survey.Required,
		},
	}, &selectedPe); err != nil {
		return "", errors.Wrap(err, "failed UserInput")
	}

	return peUUIDs[selectedPe], nil

}

func getSubnet(ctx context.Context, client *nutanixclientv3.Client, peUUID string) (string, error) {
	ctx, cancel := context.WithTimeout(ctx, 60*time.Second)
	defer cancel()
	emptyFilter := ""
	subnetsAll, err := client.V3.ListAllSubnet(emptyFilter)
	if err != nil {
		return "", errors.Wrap(err, "unable to list subnets")
	}

	subnets := subnetsAll.Entities
	// API returns an error when no results, but let's leave this in to be defensive.
	if len(subnets) == 0 {
		return "", errors.New("did not find any subnets")
	}
	if len(subnets) == 1 {
		n := *subnets[0].Spec.Name
		u := *subnets[0].Metadata.UUID
		logrus.Infof("Defaulting to only available network: %s", n)
		return u, nil
	}

	subnetUUIDs := make(map[string]string)
	var subnetChoices []string
	for _, s := range subnets {
		if *s.Spec.ClusterReference.UUID == peUUID {
			n := *s.Spec.Name
			subnetUUIDs[n] = *s.Metadata.UUID
			subnetChoices = append(subnetChoices, n)
		}
	}
	if len(subnetChoices) == 0 {
		return "", errors.New(fmt.Sprintf("could not find any subnets linked to Prism Element with UUID %s", peUUID))
	}
	sort.Strings(subnetChoices)

	var selectedSubnet string
	if err := survey.Ask([]*survey.Question{
		{
			Prompt: &survey.Select{
				Message: "Subnet",
				Options: subnetChoices,
				Help:    "The subnet to be used for installation.",
			},
			Validate: survey.Required,
		},
	}, &selectedSubnet); err != nil {
		return "", errors.Wrap(err, "failed UserInput")
	}

	return subnetUUIDs[selectedSubnet], nil
}

func getVIPs() (string, string, error) {
	var apiVIP, ingressVIP string

	if err := survey.Ask([]*survey.Question{
		{
			Prompt: &survey.Input{
				Message: "Virtual IP Address for API",
				Help:    "The VIP to be used for the OpenShift API.",
			},
			Validate: survey.ComposeValidators(survey.Required, func(ans interface{}) error {
				return validate.IP((ans).(string))
			}),
		},
	}, &apiVIP); err != nil {
		return "", "", errors.Wrap(err, "failed UserInput")
	}

	if err := survey.Ask([]*survey.Question{
		{
			Prompt: &survey.Input{
				Message: "Virtual IP Address for Ingress",
				Help:    "The VIP to be used for ingress to the cluster.",
			},
			Validate: survey.ComposeValidators(survey.Required, func(ans interface{}) error {
				if apiVIP == (ans.(string)) {
					return fmt.Errorf("%q should not be equal to the Virtual IP address for the API", ans.(string))
				}
				return validate.IP((ans).(string))
			}),
		},
	}, &ingressVIP); err != nil {
		return "", "", errors.Wrap(err, "failed UserInput")
	}

	return apiVIP, ingressVIP, nil
}
