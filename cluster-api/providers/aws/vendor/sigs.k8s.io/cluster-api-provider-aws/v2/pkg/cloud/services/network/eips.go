/*
Copyright 2018 The Kubernetes Authors.

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

package network

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/pkg/errors"

	infrav1 "sigs.k8s.io/cluster-api-provider-aws/v2/api/v1beta2"
	"sigs.k8s.io/cluster-api-provider-aws/v2/pkg/cloud/awserrors"
	"sigs.k8s.io/cluster-api-provider-aws/v2/pkg/cloud/filter"
	"sigs.k8s.io/cluster-api-provider-aws/v2/pkg/cloud/services/wait"
	"sigs.k8s.io/cluster-api-provider-aws/v2/pkg/cloud/tags"
	"sigs.k8s.io/cluster-api-provider-aws/v2/pkg/record"
)

func (s *Service) getOrAllocateAddresses(num int, role string, pool *infrav1.ElasticIPPool) (eips []string, err error) {
	out, err := s.describeAddresses(role)
	if err != nil {
		record.Eventf(s.scope.InfraCluster(), "FailedDescribeAddresses", "Failed to query addresses for role %q: %v", role, err)
		return nil, errors.Wrap(err, "failed to query addresses")
	}

	// Reuse existing unallocated addreses with the same role.
	for _, address := range out.Addresses {
		if address.AssociationId == nil {
			eips = append(eips, aws.StringValue(address.AllocationId))
		}
	}

	// allocate addresses when needed.
	tagSpecifications := tags.BuildParamsToTagSpecification(ec2.ResourceTypeElasticIp, s.getEIPTagParams(role))
	for len(eips) < num {
		allocInput := &ec2.AllocateAddressInput{
			Domain: aws.String("vpc"),
			TagSpecifications: []*ec2.TagSpecification{
				tagSpecifications,
			},
		}

		// Set EIP to consume from BYO Public IPv4 pools when defined in NetworkSpec with preflight checks.
		// The checks makes sure there is free IPs available in the pool before allocating it.
		// The check also validate the fallback strategy to consume from Amazon pool when the
		// pool is exchausted.
		if err := s.setByoPublicIpv4(pool, allocInput); err != nil {
			return nil, err
		}

		ip, err := s.allocateAddress(allocInput)
		if err != nil {
			record.Warnf(s.scope.InfraCluster(), "FailedAllocateAddress", "Failed to allocate Elastic IP for %q: %v", role, err)
			return nil, fmt.Errorf("failed to allocate Elastic IP for %q: %w", role, err)
		}
		eips = append(eips, ip)
	}

	return eips, nil
}

func (s *Service) allocateAddress(alloc *ec2.AllocateAddressInput) (string, error) {
	out, err := s.EC2Client.AllocateAddressWithContext(context.TODO(), alloc)
	if err != nil {
		return "", err
	}

	return aws.StringValue(out.AllocationId), nil
}

func (s *Service) describeAddresses(role string) (*ec2.DescribeAddressesOutput, error) {
	x := []*ec2.Filter{filter.EC2.Cluster(s.scope.Name())}
	if role != "" {
		x = append(x, filter.EC2.ProviderRole(role))
	}

	return s.EC2Client.DescribeAddressesWithContext(context.TODO(), &ec2.DescribeAddressesInput{
		Filters: x,
	})
}

func (s *Service) disassociateAddress(ip *ec2.Address) error {
	err := wait.WaitForWithRetryable(wait.NewBackoff(), func() (bool, error) {
		_, err := s.EC2Client.DisassociateAddressWithContext(context.TODO(), &ec2.DisassociateAddressInput{
			AssociationId: ip.AssociationId,
		})
		if err != nil {
			cause, _ := awserrors.Code(errors.Cause(err))
			if cause != awserrors.AssociationIDNotFound {
				return false, err
			}
		}
		return true, nil
	}, awserrors.AuthFailure)
	if err != nil {
		record.Warnf(s.scope.InfraCluster(), "FailedDisassociateEIP", "Failed to disassociate Elastic IP %q: %v", *ip.AllocationId, err)
		return errors.Wrapf(err, "failed to disassociate Elastic IP %q", *ip.AllocationId)
	}
	return nil
}

// releaseAddress releases an given EIP address back to the pool.
func (s *Service) releaseAddress(ip *ec2.Address) error {
	if ip.AssociationId != nil {
		if _, err := s.EC2Client.DisassociateAddressWithContext(context.TODO(), &ec2.DisassociateAddressInput{
			AssociationId: ip.AssociationId,
		}); err != nil {
			record.Warnf(s.scope.InfraCluster(), "FailedDisassociateEIP", "Failed to disassociate Elastic IP %q: %v", *ip.AllocationId, err)
			return errors.Errorf("failed to disassociate Elastic IP %q with allocation ID %q: Still associated with association ID %q", *ip.PublicIp, *ip.AllocationId, *ip.AssociationId)
		}
	}

	if err := wait.WaitForWithRetryable(wait.NewBackoff(), func() (bool, error) {
		_, err := s.EC2Client.ReleaseAddressWithContext(context.TODO(), &ec2.ReleaseAddressInput{AllocationId: ip.AllocationId})
		if err != nil {
			if ip.AssociationId != nil {
				if s.disassociateAddress(ip) != nil {
					return false, err
				}
			}
			return false, err
		}
		return true, nil
	}, awserrors.AuthFailure, awserrors.InUseIPAddress); err != nil {
		record.Warnf(s.scope.InfraCluster(), "FailedReleaseEIP", "Failed to disassociate Elastic IP %q: %v", *ip.AllocationId, err)
		return errors.Wrapf(err, "failed to release ElasticIP %q", *ip.AllocationId)
	}

	s.scope.Info("released ElasticIP", "eip", *ip.PublicIp, "allocation-id", *ip.AllocationId)
	return nil
}

// releaseAddressesWithFilter discovery address to be released based in filters, returning no error,
// when all addresses have been released.
func (s *Service) releaseAddressesWithFilter(filters []*ec2.Filter) error {
	out, err := s.EC2Client.DescribeAddressesWithContext(context.TODO(), &ec2.DescribeAddressesInput{
		Filters: filters,
	})
	if err != nil {
		return errors.Wrapf(err, "failed to describe elastic IPs %q", err)
	}
	if out == nil {
		return nil
	}
	for i := range out.Addresses {
		if err := s.releaseAddress(out.Addresses[i]); err != nil {
			return err
		}
	}
	return nil
}

// releaseAddresses is default cluster release flow, discoverying and releasing all
// addresses associated and owned by the cluster tag.
func (s *Service) releaseAddresses() error {
	filters := []*ec2.Filter{filter.EC2.Cluster(s.scope.Name())}
	filters = append(filters, filter.EC2.ClusterOwned(s.scope.Name()))
	return s.releaseAddressesWithFilter(filters)
}

func (s *Service) getEIPTagParams(role string) infrav1.BuildParams {
	name := fmt.Sprintf("%s-eip-%s", s.scope.Name(), role)

	return infrav1.BuildParams{
		ClusterName: s.scope.Name(),
		Lifecycle:   infrav1.ResourceLifecycleOwned,
		Name:        aws.String(name),
		Role:        aws.String(role),
		Additional:  s.scope.AdditionalTags(),
	}
}

// GetOrAllocateAddresses exports the interface to allocate an address from external services.
func (s *Service) GetOrAllocateAddresses(pool *infrav1.ElasticIPPool, num int, role string) (eips []string, err error) {
	return s.getOrAllocateAddresses(num, role, pool)
}

// GetAddresses returns the address associated to a given role.
func (s *Service) GetAddresses(role string) (*ec2.DescribeAddressesOutput, error) {
	return s.describeAddresses(role)
}

// ReleaseAddressByRole releases EIP addresses filtering by tag CAPA provider role.
func (s *Service) ReleaseAddressByRole(role string) error {
	return s.releaseAddressesWithFilter([]*ec2.Filter{
		filter.EC2.ClusterOwned(s.scope.Name()),
		filter.EC2.ProviderRole(role),
	})
}

// setByoPublicIpv4 check if the config has Public IPv4 Pool defined, then
// check if there are IPs available to consume from allocation, otherwise
// fallback to Amazon pool when explicty failure isn't defined.
func (s *Service) setByoPublicIpv4(pool *infrav1.ElasticIPPool, alloc *ec2.AllocateAddressInput) error {
	if pool == nil {
		return nil
	}
	// check if pool has free IP.
	ok, err := s.publicIpv4PoolHasAtLeastNFreeIPs(pool, 1)
	if err != nil {
		record.Warnf(s.scope.InfraCluster(), "FailedAllocateEIP", "Failed to allocate Elastic IP from Public IPv4 pool %q: %w", *pool.PublicIpv4Pool, err)
		return fmt.Errorf("failed to update Elastic IP: %w", err)
	}

	// use the custom public ipv4 pool to the Elastic IP allocation.
	if ok {
		alloc.PublicIpv4Pool = pool.PublicIpv4Pool
		return nil
	}

	// default, don't change allocation config, use Amazon pool.
	return nil
}

// publicIpv4PoolHasAtLeastNFreeIPs check if there are N IPs address available in a Public IPv4 Pool.
func (s *Service) publicIpv4PoolHasAtLeastNFreeIPs(pool *infrav1.ElasticIPPool, want int64) (bool, error) {
	if pool == nil {
		return true, nil
	}
	if pool.PublicIpv4Pool == nil {
		return true, nil
	}
	publicIpv4Pool := pool.PublicIpv4Pool
	pools, err := s.EC2Client.DescribePublicIpv4Pools(&ec2.DescribePublicIpv4PoolsInput{
		PoolIds: []*string{publicIpv4Pool},
	})
	if err != nil {
		return false, fmt.Errorf("failed to describe Public IPv4 Pool %q: %w", *publicIpv4Pool, err)
	}
	if len(pools.PublicIpv4Pools) != 1 {
		return false, fmt.Errorf("unexpected number of configured Public IPv4 Pools. want 1, got %d", len(pools.PublicIpv4Pools))
	}

	freeIPs := aws.Int64Value(pools.PublicIpv4Pools[0].TotalAvailableAddressCount)
	hasFreeIPs := freeIPs >= want

	// force to fallback to Amazon pool when the custom pool is full.
	fallbackToAmazonPool := pool.PublicIpv4PoolFallBackOrder != nil && pool.PublicIpv4PoolFallBackOrder.Equal(infrav1.PublicIpv4PoolFallbackOrderAmazonPool)
	if !hasFreeIPs && fallbackToAmazonPool {
		s.scope.Debug(fmt.Sprintf("public IPv4 pool %q has reached the limit with %d IPs available, using user-defined fallback config %q", *publicIpv4Pool, freeIPs, pool.PublicIpv4PoolFallBackOrder.String()), "eip")
		return false, nil
	}
	if !hasFreeIPs {
		return false, fmt.Errorf("public IPv4 pool %q does not have enough free IP addresses: want %d, got %d", *publicIpv4Pool, want, freeIPs)
	}
	s.scope.Debug(fmt.Sprintf("public IPv4 Pool %q has %d IPs available", *publicIpv4Pool, freeIPs), "eip")
	return true, nil
}
