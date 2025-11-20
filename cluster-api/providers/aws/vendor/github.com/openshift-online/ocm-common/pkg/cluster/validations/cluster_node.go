package validations

import (
	"fmt"
)

const (
	SingleAZCount = 1
	MultiAZCount  = 3
)

// MinReplicasValidator is responsible for verifying whether the minReplicas value adheres to specific rules:
//
// * The minReplicas value must be a positive number.
// * If multiple availability zones are used, the Cluster must have a minimum of 3 compute nodes and the total count of nodes must be a multiple of 3.
// * If a single availability zone is used, the Cluster must include at least 2 compute nodes.
// * For Hosted clusters, the number of compute nodes must be a multiple of the number of private subnets.
func MinReplicasValidator(minReplicas int, multiAZ bool, isHostedCP bool, privateSubnetsCount int) error {
	if minReplicas < 0 {
		return fmt.Errorf("The value for the number of cluster nodes must be non-negative")
	}
	if isHostedCP {
		// This value should be validated in a previous step when checking the subnets
		if privateSubnetsCount < 1 {
			return fmt.Errorf("Hosted clusters require at least a private subnet")
		}

		if minReplicas%privateSubnetsCount != 0 {
			return fmt.Errorf("Hosted clusters require that the number of compute nodes be a multiple of "+
				"the number of private subnets %d, instead received: %d", privateSubnetsCount, minReplicas)
		}
		return nil
	}

	if multiAZ {
		if minReplicas < 3 {
			return fmt.Errorf("Multi AZ cluster requires at least 3 compute nodes")
		}
		if minReplicas%3 != 0 {
			return fmt.Errorf("Multi AZ clusters require that the number of compute nodes be a multiple of 3")
		}
	} else if minReplicas < 2 {
		return fmt.Errorf("Cluster requires at least 2 compute nodes")
	}
	return nil
}

// MaxReplicasValidator is responsible for verifying whether the maxReplicas value adheres to specific rules:
//
// * The maxReplicas value must be at least equal to minReplicas.
// * If multiple availability zones are used, the number of compute nodes must be a multiple of 3.
// * For Hosted clusters, the number of compute nodes must be a multiple of the number of private subnets.
//
// The assumtion here is that minReplicas was already validated
func MaxReplicasValidator(minReplicas int, maxReplicas int, multiAZ bool, isHostedCP bool, privateSubnetsCount int) error {
	if minReplicas > maxReplicas {
		return fmt.Errorf("max-replicas must be greater or equal to min-replicas")
	}

	if isHostedCP {
		if maxReplicas%privateSubnetsCount != 0 {
			return fmt.Errorf("Hosted clusters require that the number of compute nodes be a multiple of "+
				"the number of private subnets %d, instead received: %d", privateSubnetsCount, maxReplicas)
		}
		return nil
	}

	if multiAZ && maxReplicas%3 != 0 {
		return fmt.Errorf("Multi AZ clusters require that the number of compute nodes be a multiple of 3")
	}
	return nil
}

// ValidateAvailabilityZonesCount is responsible for verifying whether the number of availability zones adheres to specific rules:
//
// * The number of availability zones for a multi AZ cluster should be 3.
// * The number of availability zones for a single AZ cluster should be 1.
func ValidateAvailabilityZonesCount(multiAZ bool, availabilityZonesCount int) error {
	if multiAZ && availabilityZonesCount != MultiAZCount {
		return fmt.Errorf("The number of availability zones for a multi AZ cluster should be %d, "+
			"instead received: %d", MultiAZCount, availabilityZonesCount)
	}
	if !multiAZ && availabilityZonesCount != SingleAZCount {
		return fmt.Errorf("The number of availability zones for a single AZ cluster should be %d, "+
			"instead received: %d", SingleAZCount, availabilityZonesCount)
	}

	return nil
}
