package elb

import (
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/elbv2"

	infrav1 "sigs.k8s.io/cluster-api-provider-aws/v2/api/v1beta2"
)

func getElasticIPRoleName() string {
	return fmt.Sprintf("lb-%s", infrav1.APIServerRoleTagValue)
}

// allocatePublicIpv4AddressFromByoIPPool claims for Elastic IPs from an user-defined public IPv4 pool,
// allocating it to the NetworkMapping structure from an Network Load Balancer.
func (s *Service) allocatePublicIpv4AddressFromByoIPPool(input *elbv2.CreateLoadBalancerInput) error {
	// Custom Public IPv4 Pool isn't set.
	if s.scope.VPC().GetPublicIpv4Pool() == nil {
		return nil
	}

	// Only NLB is supported
	if input.Type == nil {
		return fmt.Errorf("PublicIpv4Pool is supported only when the Load Balancer type is %q", elbv2.LoadBalancerTypeEnumNetwork)
	}
	if *input.Type != string(elbv2.LoadBalancerTypeEnumNetwork) {
		return fmt.Errorf("PublicIpv4Pool is not supported with Load Balancer type %s. Use Network Load Balancer instead", *input.Type)
	}

	// Custom SubnetMappings should not be defined or overridden by user-defined mapping.
	if len(input.SubnetMappings) > 0 {
		return fmt.Errorf("PublicIpv4Pool is mutually exclusive with SubnetMappings")
	}

	eips, err := s.netService.GetOrAllocateAddresses(s.scope.VPC().GetElasticIPPool(), len(input.Subnets), getElasticIPRoleName())
	if err != nil {
		return fmt.Errorf("failed to allocate address from Public IPv4 Pool %q to role %s: %w", *s.scope.VPC().GetPublicIpv4Pool(), getElasticIPRoleName(), err)
	}
	if len(eips) != len(input.Subnets) {
		return fmt.Errorf("number of allocated EIP addresses (%d) from pool %q must match with the subnet count (%d)", len(eips), *s.scope.VPC().GetPublicIpv4Pool(), len(input.Subnets))
	}
	for cnt, sb := range input.Subnets {
		input.SubnetMappings = append(input.SubnetMappings, &elbv2.SubnetMapping{
			SubnetId:     aws.String(*sb),
			AllocationId: aws.String(eips[cnt]),
		})
	}
	// Subnets and SubnetMappings are mutual exclusive. Cleaning Subnets when BYO IP is defined,
	// and SubnetMappings are mounted.
	input.Subnets = []*string{}

	return nil
}
