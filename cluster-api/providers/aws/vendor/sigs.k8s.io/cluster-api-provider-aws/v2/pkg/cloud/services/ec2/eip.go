package ec2

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ec2"

	infrav1 "sigs.k8s.io/cluster-api-provider-aws/v2/api/v1beta2"
	"sigs.k8s.io/cluster-api-provider-aws/v2/pkg/record"
)

func getElasticIPRoleName(instanceID string) string {
	return fmt.Sprintf("ec2-%s", instanceID)
}

// ReconcileElasticIPFromPublicPool reconciles the elastic IP from a custom Public IPv4 Pool.
func (s *Service) ReconcileElasticIPFromPublicPool(pool *infrav1.ElasticIPPool, instance *infrav1.Instance) error {
	// TODO: check if the instance is in the state allowing EIP association.
	// Expected instance states: pending or running
	// https://docs.aws.amazon.com/AWSEC2/latest/UserGuide/ec2-instance-lifecycle.html
	if err := s.getAndAssociateAddressesToInstance(pool, getElasticIPRoleName(instance.ID), instance.ID); err != nil {
		return fmt.Errorf("failed to reconcile Elastic IP: %w", err)
	}
	return nil
}

// ReleaseElasticIP releases a specific elastic IP based on the instance role.
func (s *Service) ReleaseElasticIP(instanceID string) error {
	return s.netService.ReleaseAddressByRole(getElasticIPRoleName(instanceID))
}

// getAndAssociateAddressesToInstance find or create an EIP from an instance and role.
func (s *Service) getAndAssociateAddressesToInstance(pool *infrav1.ElasticIPPool, role string, instance string) (err error) {
	eips, err := s.netService.GetOrAllocateAddresses(pool, 1, role)
	if err != nil {
		record.Warnf(s.scope.InfraCluster(), "FailedAllocateEIP", "Failed to get Elastic IP for %q: %v", role, err)
		return err
	}
	if len(eips) != 1 {
		record.Warnf(s.scope.InfraCluster(), "FailedAllocateEIP", "Failed to allocate Elastic IP for %q: %v", role, err)
		return fmt.Errorf("unexpected number of Elastic IP to instance %q, got %d: %w", instance, len(eips), err)
	}
	_, err = s.EC2Client.AssociateAddressWithContext(context.TODO(), &ec2.AssociateAddressInput{
		InstanceId:   aws.String(instance),
		AllocationId: aws.String(eips[0]),
	})
	if err != nil {
		record.Warnf(s.scope.InfraCluster(), "FailedAssociateEIP", "Failed to associate Elastic IP for %q: %v", role, err)
		return fmt.Errorf("failed to associate Elastic IP %q to instance %q: %w", eips[0], instance, err)
	}
	return nil
}
