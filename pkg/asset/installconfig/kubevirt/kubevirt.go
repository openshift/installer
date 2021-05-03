package kubevirt

import (
	"github.com/AlecAivazis/survey/v2"

	"github.com/openshift/installer/pkg/types/kubevirt"
)

// Platform collects kubevirt-specific configuration.
func Platform() (*kubevirt.Platform, error) {
	var (
		namespace, apiVIP, ingressVIP, networkName string
		err                                        error
	)

	if namespace, err = selectNamespace(); err != nil {
		return nil, err
	}

	if apiVIP, err = selectAPIVIP(); err != nil {
		return nil, err
	}

	if ingressVIP, err = selectIngressVIP(); err != nil {
		return nil, err
	}

	if networkName, err = selectNetworkName(); err != nil {
		return nil, err
	}

	return &kubevirt.Platform{
		Namespace:   namespace,
		APIVIP:      apiVIP,
		IngressVIP:  ingressVIP,
		NetworkName: networkName,
	}, nil
}

func selectNamespace() (string, error) {
	var selectedNamespace string

	err := survey.Ask([]*survey.Question{
		{
			Prompt: &survey.Input{
				Message: "Namespace",
				Help:    "The namespace, in the infracluster, where all the resources of the tenantcluster would be created.",
			},
		},
	}, &selectedNamespace)

	return selectedNamespace, err
}

func selectAPIVIP() (string, error) {
	var selectedAPIVIP string

	err := survey.Ask([]*survey.Question{
		{
			Prompt: &survey.Input{
				Message: "API VIP",
				Help:    "The Virtual IP address used for external access to the OpenShift API.",
			},
		},
	}, &selectedAPIVIP)

	return selectedAPIVIP, err
}

func selectIngressVIP() (string, error) {
	var selectedIngressVIP string

	err := survey.Ask([]*survey.Question{
		{
			Prompt: &survey.Input{
				Message: "Ingress VIP",
				Help:    "An external IP which routes to the default ingress controller.",
			},
		},
	}, &selectedIngressVIP)

	return selectedIngressVIP, err
}

func selectNetworkName() (string, error) {
	var selectedNetworkName string

	err := survey.Ask([]*survey.Question{
		{
			Prompt: &survey.Input{
				Message: "Network Name",
				Help:    "The target network of all the network interfaces of the nodes.",
			},
		},
	}, &selectedNetworkName)

	return selectedNetworkName, err
}
