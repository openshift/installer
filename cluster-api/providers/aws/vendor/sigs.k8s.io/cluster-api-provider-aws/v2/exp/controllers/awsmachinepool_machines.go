package controllers

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/go-logr/logr"
	"github.com/pkg/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/klog/v2"
	"k8s.io/utils/ptr"
	"sigs.k8s.io/controller-runtime/pkg/client"

	infrav1 "sigs.k8s.io/cluster-api-provider-aws/v2/api/v1beta2"
	expinfrav1 "sigs.k8s.io/cluster-api-provider-aws/v2/exp/api/v1beta2"
	"sigs.k8s.io/cluster-api-provider-aws/v2/feature"
	"sigs.k8s.io/cluster-api-provider-aws/v2/pkg/cloud/services"
	"sigs.k8s.io/cluster-api-provider-aws/v2/pkg/cloud/services/ec2"
	clusterv1 "sigs.k8s.io/cluster-api/api/v1beta1"
	expclusterv1 "sigs.k8s.io/cluster-api/exp/api/v1beta1"
	"sigs.k8s.io/cluster-api/util"
	"sigs.k8s.io/cluster-api/util/labels/format"
)

func createAWSMachinesIfNotExists(ctx context.Context, awsMachineList *infrav1.AWSMachineList, mp *expclusterv1.MachinePool, infraMachinePoolMeta *metav1.ObjectMeta, infraMachinePoolType *metav1.TypeMeta, existingASG *expinfrav1.AutoScalingGroup, l logr.Logger, client client.Client, ec2Svc services.EC2Interface) error {
	if !feature.Gates.Enabled(feature.MachinePoolMachines) {
		return errors.New("createAWSMachinesIfNotExists must not be called unless the MachinePoolMachines feature gate is enabled")
	}

	l.V(4).Info("Creating missing AWSMachines")

	providerIDToAWSMachine := make(map[string]infrav1.AWSMachine, len(awsMachineList.Items))
	for i := range awsMachineList.Items {
		awsMachine := awsMachineList.Items[i]
		if awsMachine.Spec.ProviderID == nil || *awsMachine.Spec.ProviderID == "" {
			continue
		}
		providerID := *awsMachine.Spec.ProviderID
		providerIDToAWSMachine[providerID] = awsMachine
	}

	for i := range existingASG.Instances {
		instanceID := existingASG.Instances[i].ID
		providerID := fmt.Sprintf("aws:///%s/%s", existingASG.Instances[i].AvailabilityZone, instanceID)

		instanceLogger := l.WithValues("providerID", providerID, "instanceID", instanceID, "asg", existingASG.Name)
		instanceLogger.V(4).Info("Checking if machine pool AWSMachine is up to date")
		if _, exists := providerIDToAWSMachine[providerID]; exists {
			continue
		}

		instance, err := ec2Svc.InstanceIfExists(&instanceID)
		if errors.Is(err, ec2.ErrInstanceNotFoundByID) {
			instanceLogger.V(4).Info("Instance not found, it may have already been deleted")
			continue
		}
		if err != nil {
			return fmt.Errorf("failed to look up EC2 instance %q: %w", instanceID, err)
		}

		securityGroups := make([]infrav1.AWSResourceReference, 0, len(instance.SecurityGroupIDs))
		for j := range instance.SecurityGroupIDs {
			securityGroups = append(securityGroups, infrav1.AWSResourceReference{
				ID: aws.String(instance.SecurityGroupIDs[j]),
			})
		}

		awsMachine := &infrav1.AWSMachine{
			ObjectMeta: metav1.ObjectMeta{
				Namespace:    mp.Namespace,
				GenerateName: fmt.Sprintf("%s-", existingASG.Name),
				Labels: map[string]string{
					clusterv1.MachinePoolNameLabel: format.MustFormatValue(mp.Name),
					clusterv1.ClusterNameLabel:     mp.Spec.ClusterName,
				},
				OwnerReferences: []metav1.OwnerReference{
					{
						APIVersion:         infraMachinePoolType.APIVersion,
						Kind:               infraMachinePoolType.Kind,
						Name:               infraMachinePoolMeta.Name,
						BlockOwnerDeletion: ptr.To(true),
						UID:                infraMachinePoolMeta.UID,
					},
				},
			},
			Spec: infrav1.AWSMachineSpec{
				ProviderID: aws.String(providerID),
				InstanceID: aws.String(instanceID),

				// Store some extra fields for informational purposes (not needed by CAPA)
				AMI: infrav1.AMIReference{
					ID: aws.String(instance.ImageID),
				},
				InstanceType:             instance.Type,
				PublicIP:                 aws.Bool(instance.PublicIP != nil),
				SSHKeyName:               instance.SSHKeyName,
				InstanceMetadataOptions:  instance.InstanceMetadataOptions,
				IAMInstanceProfile:       instance.IAMProfile,
				AdditionalSecurityGroups: securityGroups,
				Subnet:                   &infrav1.AWSResourceReference{ID: aws.String(instance.SubnetID)},
				RootVolume:               instance.RootVolume,
				NonRootVolumes:           instance.NonRootVolumes,
				NetworkInterfaces:        instance.NetworkInterfaces,
				CloudInit:                infrav1.CloudInit{},
				SpotMarketOptions:        instance.SpotMarketOptions,
				Tenancy:                  instance.Tenancy,
			},
		}
		instanceLogger.V(4).Info("Creating AWSMachine")
		if err := client.Create(ctx, awsMachine); err != nil {
			return fmt.Errorf("failed to create AWSMachine: %w", err)
		}
	}
	return nil
}

func deleteOrphanedAWSMachines(ctx context.Context, awsMachineList *infrav1.AWSMachineList, existingASG *expinfrav1.AutoScalingGroup, l logr.Logger, client client.Client) error {
	if !feature.Gates.Enabled(feature.MachinePoolMachines) {
		return errors.New("deleteOrphanedAWSMachines must not be called unless the MachinePoolMachines feature gate is enabled")
	}

	l.V(4).Info("Deleting orphaned AWSMachines")
	providerIDToInstance := make(map[string]infrav1.Instance, len(existingASG.Instances))
	for i := range existingASG.Instances {
		providerID := fmt.Sprintf("aws:///%s/%s", existingASG.Instances[i].AvailabilityZone, existingASG.Instances[i].ID)
		providerIDToInstance[providerID] = existingASG.Instances[i]
	}

	for i := range awsMachineList.Items {
		awsMachine := awsMachineList.Items[i]
		if awsMachine.Spec.ProviderID == nil || *awsMachine.Spec.ProviderID == "" {
			continue
		}

		providerID := *awsMachine.Spec.ProviderID
		if _, exists := providerIDToInstance[providerID]; exists {
			continue
		}

		machine, err := util.GetOwnerMachine(ctx, client, awsMachine.ObjectMeta)
		if err != nil {
			return fmt.Errorf("failed to get owner Machine for %s/%s: %w", awsMachine.Namespace, awsMachine.Name, err)
		}
		machineLogger := l.WithValues("machine", klog.KObj(machine), "awsmachine", klog.KObj(&awsMachine), "ProviderID", providerID)
		machineLogger.V(4).Info("Deleting orphaned Machine")
		if machine == nil {
			machineLogger.Info("No machine owner found for AWSMachine, deleting AWSMachine anyway.")
			if err := client.Delete(ctx, &awsMachine); err != nil {
				return fmt.Errorf("failed to delete orphan AWSMachine %s/%s: %w", awsMachine.Namespace, awsMachine.Name, err)
			}
			machineLogger.V(4).Info("Deleted AWSMachine")
			continue
		}

		if err := client.Delete(ctx, machine); err != nil {
			return fmt.Errorf("failed to delete orphan Machine %s/%s: %w", machine.Namespace, machine.Name, err)
		}
		machineLogger.V(4).Info("Deleted Machine")
	}
	return nil
}

func getAWSMachines(ctx context.Context, mp *expclusterv1.MachinePool, kubeClient client.Client) (*infrav1.AWSMachineList, error) {
	if !feature.Gates.Enabled(feature.MachinePoolMachines) {
		return nil, errors.New("getAWSMachines must not be called unless the MachinePoolMachines feature gate is enabled")
	}

	awsMachineList := &infrav1.AWSMachineList{}
	labels := map[string]string{
		clusterv1.MachinePoolNameLabel: format.MustFormatValue(mp.Name),
		clusterv1.ClusterNameLabel:     mp.Spec.ClusterName,
	}
	if err := kubeClient.List(ctx, awsMachineList, client.InNamespace(mp.Namespace), client.MatchingLabels(labels)); err != nil {
		return nil, err
	}
	return awsMachineList, nil
}

func reconcileDeleteAWSMachines(ctx context.Context, mp *expclusterv1.MachinePool, client client.Client, l logr.Logger) error {
	if !feature.Gates.Enabled(feature.MachinePoolMachines) {
		return errors.New("reconcileDeleteAWSMachines must not be called unless the MachinePoolMachines feature gate is enabled")
	}

	awsMachineList, err := getAWSMachines(ctx, mp, client)
	if err != nil {
		return err
	}
	for i := range awsMachineList.Items {
		awsMachine := awsMachineList.Items[i]
		if awsMachine.DeletionTimestamp.IsZero() {
			continue
		}
		logger := l.WithValues("awsmachine", klog.KObj(&awsMachine))
		// delete the owner Machine resource for the AWSMachine so that CAPI can clean up gracefully
		machine, err := util.GetOwnerMachine(ctx, client, awsMachine.ObjectMeta)
		if err != nil {
			logger.V(2).Info("Failed to get owner Machine", "err", err.Error())
			continue
		}

		if err := client.Delete(ctx, machine); err != nil {
			logger.V(2).Info("Failed to delete owner Machine", "err", err.Error())
		}
	}
	return nil
}
