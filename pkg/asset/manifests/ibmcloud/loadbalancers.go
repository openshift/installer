package ibmcloud

import (
	"fmt"

	"k8s.io/utils/ptr"
	capibmcloud "sigs.k8s.io/cluster-api-provider-ibmcloud/api/v1beta2"

	ibmcloudic "github.com/openshift/installer/pkg/asset/installconfig/ibmcloud"
	"github.com/openshift/installer/pkg/types"
)

const (
	// healthMonitorURLReadyz is the health monitoring URL used to report when healthy/ready.
	healthMonitorURLReadyz = "/readyz"
)

func getLoadBalancers(infraID string, securityGroups []capibmcloud.VPCResource, subnets []capibmcloud.VPCResource, publish types.PublishingStrategy) []capibmcloud.VPCLoadBalancerSpec {
	loadBalancers := make([]capibmcloud.VPCLoadBalancerSpec, 0, 2)

	loadBalancers = append(loadBalancers, buildPrivateLoadBalancer(infraID, securityGroups, subnets))
	if publish == types.ExternalPublishingStrategy {
		loadBalancers = append(loadBalancers, buildPublicLoadBalancer(infraID, securityGroups, subnets))
	}

	return loadBalancers
}

func buildPrivateLoadBalancer(infraID string, securityGroups []capibmcloud.VPCResource, subnets []capibmcloud.VPCResource) capibmcloud.VPCLoadBalancerSpec {
	kubeAPIBackendPoolNamePtr := ptr.To(fmt.Sprintf("%s-%s", infraID, ibmcloudic.KubernetesAPIPrivateSuffix))
	machineConfigBackendPoolNamePtr := ptr.To(fmt.Sprintf("%s-%s", infraID, ibmcloudic.MachineConfigSuffix))

	return capibmcloud.VPCLoadBalancerSpec{
		Name:   fmt.Sprintf("%s-%s", infraID, ibmcloudic.KubernetesAPIPrivateSuffix),
		Public: ptr.To(false),
		AdditionalListeners: []capibmcloud.AdditionalListenerSpec{
			{
				DefaultPoolName: kubeAPIBackendPoolNamePtr,
				Port:            ibmcloudic.KubernetesAPIPort,
				Protocol:        &capibmcloud.VPCLoadBalancerListenerProtocolTCP,
			},
			{
				DefaultPoolName: machineConfigBackendPoolNamePtr,
				Port:            ibmcloudic.MachineConfigServerPort,
				Protocol:        &capibmcloud.VPCLoadBalancerListenerProtocolTCP,
			},
		},
		BackendPools: []capibmcloud.VPCLoadBalancerBackendPoolSpec{
			{
				// Kubernetes API pool
				Name:      kubeAPIBackendPoolNamePtr,
				Algorithm: capibmcloud.VPCLoadBalancerBackendPoolAlgorithmRoundRobin,
				Protocol:  capibmcloud.VPCLoadBalancerBackendPoolProtocolTCP,
				HealthMonitor: capibmcloud.VPCLoadBalancerHealthMonitorSpec{
					Delay:   60,
					Retries: 5,
					Timeout: 30,
					Type:    capibmcloud.VPCLoadBalancerBackendPoolHealthMonitorTypeHTTPS,
					URLPath: ptr.To(healthMonitorURLReadyz),
				},
			},
			{
				// Machine Config Server pool
				Name:      machineConfigBackendPoolNamePtr,
				Algorithm: capibmcloud.VPCLoadBalancerBackendPoolAlgorithmRoundRobin,
				Protocol:  capibmcloud.VPCLoadBalancerBackendPoolProtocolTCP,
				HealthMonitor: capibmcloud.VPCLoadBalancerHealthMonitorSpec{
					Delay:   60,
					Retries: 5,
					Timeout: 30,
					Type:    capibmcloud.VPCLoadBalancerBackendPoolHealthMonitorTypeTCP,
					URLPath: ptr.To(healthMonitorURLReadyz),
				},
			},
		},
		SecurityGroups: securityGroups,
		Subnets:        subnets,
	}
}

func buildPublicLoadBalancer(infraID string, securityGroups []capibmcloud.VPCResource, subnets []capibmcloud.VPCResource) capibmcloud.VPCLoadBalancerSpec {
	backendPoolNamePtr := ptr.To(fmt.Sprintf("%s-%s", infraID, ibmcloudic.KubernetesAPIPublicSuffix))

	return capibmcloud.VPCLoadBalancerSpec{
		Name:   fmt.Sprintf("%s-%s", infraID, ibmcloudic.KubernetesAPIPublicSuffix),
		Public: ptr.To(true),
		AdditionalListeners: []capibmcloud.AdditionalListenerSpec{
			{
				DefaultPoolName: backendPoolNamePtr,
				Port:            ibmcloudic.KubernetesAPIPort,
				Protocol:        &capibmcloud.VPCLoadBalancerListenerProtocolTCP,
			},
		},
		BackendPools: []capibmcloud.VPCLoadBalancerBackendPoolSpec{
			{
				// Kubernetes API pool
				Name:      backendPoolNamePtr,
				Algorithm: capibmcloud.VPCLoadBalancerBackendPoolAlgorithmRoundRobin,
				Protocol:  capibmcloud.VPCLoadBalancerBackendPoolProtocolTCP,
				HealthMonitor: capibmcloud.VPCLoadBalancerHealthMonitorSpec{
					Delay:   60,
					Retries: 5,
					Timeout: 30,
					Type:    capibmcloud.VPCLoadBalancerBackendPoolHealthMonitorTypeHTTPS,
					URLPath: ptr.To(healthMonitorURLReadyz),
				},
			},
		},
		SecurityGroups: securityGroups,
		Subnets:        subnets,
	}
}
