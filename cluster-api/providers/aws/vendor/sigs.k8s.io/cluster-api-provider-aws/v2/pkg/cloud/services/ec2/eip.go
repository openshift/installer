package ec2

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ec2"
	"k8s.io/utils/ptr"

	infrav1 "sigs.k8s.io/cluster-api-provider-aws/v2/api/v1beta2"
	"sigs.k8s.io/cluster-api-provider-aws/v2/pkg/record"
)

func getElasticIPRoleName(instanceID string) string {
	return fmt.Sprintf("ec2-%s", instanceID)
}

// ReconcileElasticIPFromPublicPool reconciles the elastic IP from a custom Public IPv4 Pool.
func (s *Service) ReconcileElasticIPFromPublicPool(pool *infrav1.ElasticIPPool, instance *infrav1.Instance) (bool, error) {
	shouldRequeue := true
	// Should not happen
	if pool == nil {
		return shouldRequeue, fmt.Errorf("unexpected behavior, pool must be set when reconcile ElasticIPPool")
	}
	iip := ptr.Deref(instance.PublicIP, "")
	s.scope.Debug("Reconciling machine with custom Public IPv4 Pool", "instance-id", instance.ID, "instance-state", instance.State, "instance-public-ip", iip, "publicIpv4PoolID", pool.PublicIpv4Pool)

	// Requeue when the instance is not ready to be associated.
	if instance.State != infrav1.InstanceStateRunning {
		s.scope.Debug("Unable to reconcile Elastic IP Pool for instance", "instance-id", instance.ID, "instance-state", instance.State)
		return shouldRequeue, nil
	}

	// All done, must reconcile only when the instance is in running state.
	shouldRequeue = false

	// Prevent running association every reconciliation when it is already done.
	addrs, err := s.netService.GetAddresses(getElasticIPRoleName(instance.ID))
	if err != nil {
		s.scope.Error(err, "error checking if addresses exists for Elastic IP Pool to machine", "eip-role", getElasticIPRoleName(instance.ID))
		return shouldRequeue, err
	}
	if len(addrs.Addresses) > 0 {
		if len(addrs.Addresses) != 1 {
			return shouldRequeue, fmt.Errorf("unexpected number of EIPs allocated to the role. expected 1, got %d", len(addrs.Addresses))
		}
		addr := addrs.Addresses[0]
		// address is already associated.
		if addr.AssociationId != nil && addr.InstanceId != nil && *addr.InstanceId == instance.ID {
			s.scope.Debug("Machine is already associated with an Elastic IP with custom Public IPv4 pool", "eip", addr.AllocationId, "eip-address", addr.PublicIp, "eip-associationID", addr.AssociationId, "eip-instance", addr.InstanceId)
			return shouldRequeue, nil
		}
	}

	// Get existing, or allocate an EIP, then Associate to the machine.
	// Should requeue if any error is returned in the process.
	if err := s.getAndAssociateAddressesToInstance(pool, getElasticIPRoleName(instance.ID), instance.ID); err != nil {
		return true, fmt.Errorf("failed to reconcile Elastic IP: %w", err)
	}
	return shouldRequeue, nil
}

// ReleaseElasticIP releases a specific Elastic IP based on the instance role.
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
