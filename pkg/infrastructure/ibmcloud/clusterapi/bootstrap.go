package clusterapi

import (
	"context"
	"fmt"
	"reflect"

	"github.com/IBM/vpc-go-sdk/vpcv1"
	"github.com/sirupsen/logrus"

	ibmcloudic "github.com/openshift/installer/pkg/asset/installconfig/ibmcloud"
	ibmcloudtypes "github.com/openshift/installer/pkg/types/ibmcloud"
)

func cleanupIgnitionCOSBucket(ctx context.Context, client ibmcloudic.API, instanceID string, bucketName string, region string) error {
	// First, check whether the bucket exists.
	bucketDetails, err := client.GetCOSBucketByName(ctx, instanceID, bucketName, region)
	switch {
	case err != nil:
		return fmt.Errorf("failed checking for ignition bucket: %w", err)
	case bucketDetails == nil:
		logrus.Debugf("ignition bucket not found, skipping cleanup")
		return nil
	default:
		logrus.Debugf("found ignition bucket %s", *bucketDetails.Name)
	}

	// Buckets need to be empty to complete deletion, remove all objects in the bucket.
	cosObjectsDetails, err := client.ListCOSObjects(ctx, instanceID, bucketName, region)
	switch {
	case err != nil:
		return fmt.Errorf("failed collecting objects in ignition bucket: %w", err)
	case cosObjectsDetails == nil:
		logrus.Debugf("ignition bucket appears to be empty %s", *bucketDetails.Name)
	default:
		for _, object := range cosObjectsDetails.Contents {
			if object.Key == nil {
				return fmt.Errorf("cos object has no key in ignition bucket")
			}
			logrus.Debugf("deleting cos object %s in ignition bucket", *object.Key)
			if err = client.DeleteCOSObject(ctx, instanceID, bucketName, *object.Key, region); err != nil {
				return fmt.Errorf("failed to delete object %s in ignition bucket: %w", *object.Key, err)
			}
		}
	}

	// Finally, delete the bucket itself.
	logrus.Debugf("deleting ignition bucket %s", bucketName)
	if err = client.DeleteCOSBucket(ctx, instanceID, bucketName, region); err != nil {
		return fmt.Errorf("failed to delete ignition cos bucket: %w", err)
	}
	return nil
}

func cleanupBootstrapSecurityGroup(ctx context.Context, client ibmcloudic.API, vpcID string, securityGroupName string, region string) error {
	sgDetails, err := client.GetSecurityGroupByName(ctx, securityGroupName, vpcID, region)
	switch {
	case err != nil:
		return fmt.Errorf("failed retrieving security group as %s for destroy bootstrap: %w", securityGroupName, err)
	case sgDetails != nil:
		// After finding the bootstrap SG, check if we have to detach the bootstrap network interface (which should be the only attached target) from the SG first.
		logrus.Debugf("checking bootstrap network interface of bootstrap instance for bootstrap security group %s", *sgDetails.ID)
		for _, target := range sgDetails.Targets {
			// Check if the target is a SG Network Interface reference type first, then default to the generic SG target.
			networkInterfaceTarget, ok := target.(*vpcv1.SecurityGroupTargetReferenceNetworkInterfaceReferenceTargetContext)
			switch {
			case ok:
				logrus.Debugf("removing bootstrap network interface: %s", *networkInterfaceTarget.ID)
				if err = client.DeleteSecurityGroupTargetBinding(ctx, *sgDetails.ID, *networkInterfaceTarget.ID, region); err != nil {
					return fmt.Errorf("failed to detach bootstrap network interface %s from bootstrap security group %s: %w", *networkInterfaceTarget.ID, *sgDetails.ID, err)
				}
				logrus.Debugf("removed bootstrap network interface %s from bootstrap security group %s", *networkInterfaceTarget.ID, *sgDetails.ID)
			case reflect.TypeOf(target).String() == ibmcloudtypes.IBMCloudInfrastructureSecurityGroupTargetReference:
				targetReference := target.(*vpcv1.SecurityGroupTargetReference)
				switch {
				case targetReference.ResourceType != nil && *targetReference.ResourceType == vpcv1.SecurityGroupTargetReferenceBareMetalServerNetworkInterfaceReferenceTargetContextResourceTypeNetworkInterfaceConst:
					logrus.Debugf("removing bootstrap generic network interface: %s", *targetReference.ID)
					if err = client.DeleteSecurityGroupTargetBinding(ctx, *sgDetails.ID, *targetReference.ID, region); err != nil {
						return fmt.Errorf("failed to detach bootstrap generic network interface %s from bootstrap security group %s: %w", *targetReference.ID, *sgDetails.ID, err)
					}
					logrus.Debugf("removed bootstrap generic network interface %s from bootstrap security group %s", *targetReference.ID, *sgDetails.ID)
				default:
					targetReferenceType := ""
					if targetReference.ResourceType != nil {
						targetReferenceType = *targetReference.ResourceType
					}
					return fmt.Errorf("unexpected targetReference %s=%s", *targetReference.ID, targetReferenceType)
				}
			default:
				return fmt.Errorf("unexpected target attached to bootstrap security group %s=%s", *sgDetails.ID, reflect.TypeOf(target).String())
			}
		}

		logrus.Debugf("deleting bootstrap security group: %s=%s", *sgDetails.ID, securityGroupName)
		err = client.DeleteSecurityGroup(ctx, *sgDetails.ID, region)
		if err != nil {
			return fmt.Errorf("failed to delete bootstrap security group %s: %w", *sgDetails.ID, err)
		}
		logrus.Debugf("deleted bootstrap security group: %s", *sgDetails.ID)
	default:
		logrus.Debugf("no bootstrap security group found for the cluster as %s, skipping security group cleanup", securityGroupName)
	}
	return nil
}
