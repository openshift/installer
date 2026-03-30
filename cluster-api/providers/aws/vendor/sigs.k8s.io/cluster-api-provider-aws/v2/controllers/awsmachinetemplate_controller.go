/*
Copyright 2025 The Kubernetes Authors.

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

package controllers

import (
	"context"
	"fmt"
	"strings"

	"github.com/aws/aws-sdk-go-v2/service/ec2"
	ec2types "github.com/aws/aws-sdk-go-v2/service/ec2/types"
	"github.com/pkg/errors"
	corev1 "k8s.io/api/core/v1"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/utils/ptr"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/builder"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller"
	"sigs.k8s.io/controller-runtime/pkg/event"
	"sigs.k8s.io/controller-runtime/pkg/handler"
	"sigs.k8s.io/controller-runtime/pkg/predicate"

	infrav1 "sigs.k8s.io/cluster-api-provider-aws/v2/api/v1beta2"
	ekscontrolplanev1 "sigs.k8s.io/cluster-api-provider-aws/v2/controlplane/eks/api/v1beta2"
	"sigs.k8s.io/cluster-api-provider-aws/v2/pkg/cloud/scope"
	ec2service "sigs.k8s.io/cluster-api-provider-aws/v2/pkg/cloud/services/ec2"
	"sigs.k8s.io/cluster-api-provider-aws/v2/pkg/logger"
	"sigs.k8s.io/cluster-api-provider-aws/v2/pkg/record"
	"sigs.k8s.io/cluster-api-provider-aws/v2/util/paused"
	controlplanev1 "sigs.k8s.io/cluster-api/api/controlplane/kubeadm/v1beta2"
	clusterv1 "sigs.k8s.io/cluster-api/api/core/v1beta2"
	"sigs.k8s.io/cluster-api/util"
	"sigs.k8s.io/cluster-api/util/predicates"
)

const (
	// awsMachineTemplateKind is the Kind name for AWSMachineTemplate resources.
	awsMachineTemplateKind = "AWSMachineTemplate"
)

// AWSMachineTemplateReconciler reconciles AWSMachineTemplate objects.
//
// This controller automatically populates capacity information for AWSMachineTemplate resources
// to enable autoscaling from zero.
//
// See: https://github.com/kubernetes-sigs/cluster-api/blob/main/docs/proposals/20210310-opt-in-autoscaling-from-zero.md
type AWSMachineTemplateReconciler struct {
	client.Client
	WatchFilterValue string
}

// SetupWithManager sets up the controller with the Manager.
func (r *AWSMachineTemplateReconciler) SetupWithManager(ctx context.Context, mgr ctrl.Manager, options controller.Options) error {
	log := logger.FromContext(ctx)

	b := ctrl.NewControllerManagedBy(mgr).
		For(&infrav1.AWSMachineTemplate{}).
		WithOptions(options).
		WithEventFilter(predicates.ResourceNotPausedAndHasFilterLabel(mgr.GetScheme(), log.GetLogger(), r.WatchFilterValue)).
		Watches(
			&clusterv1.MachineDeployment{},
			handler.EnqueueRequestsFromMapFunc(r.machineDeploymentToAWSMachineTemplate),
			// Only emit events for creation to reconcile in case the MachineDeployment got created after the AWSMachineTemplate was reconciled.
			builder.WithPredicates(resourceCreatedPredicate),
		).
		Watches(
			&clusterv1.MachineSet{},
			handler.EnqueueRequestsFromMapFunc(r.machineSetToAWSMachineTemplate),
			// Only emit events for creation to reconcile in case the MachineSet got created after the AWSMachineTemplate was reconciled.
			builder.WithPredicates(resourceCreatedPredicate),
		)

	// Watch KubeadmControlPlane if they exist.
	if _, err := mgr.GetRESTMapper().RESTMapping(schema.GroupKind{Group: controlplanev1.GroupVersion.Group, Kind: "KubeadmControlPlane"}, controlplanev1.GroupVersion.Version); err == nil {
		b = b.Watches(&controlplanev1.KubeadmControlPlane{},
			handler.EnqueueRequestsFromMapFunc(r.kubeadmControlPlaneToAWSMachineTemplate),
			// Only emit events for creation to reconcile in case the KubeadmControlPlane got created after the AWSMachineTemplate was reconciled.
			builder.WithPredicates(resourceCreatedPredicate),
		)
	}

	_, err := b.Build(r)
	if err != nil {
		return errors.Wrap(err, "failed setting up with a controller manager")
	}

	return nil
}

// +kubebuilder:rbac:groups=infrastructure.cluster.x-k8s.io,resources=awsmachinetemplates,verbs=get;list;watch
// +kubebuilder:rbac:groups=infrastructure.cluster.x-k8s.io,resources=awsmachinetemplates/status,verbs=get;update;patch
// +kubebuilder:rbac:groups=infrastructure.cluster.x-k8s.io,resources=awsclusters,verbs=get;list;watch
// +kubebuilder:rbac:groups=cluster.x-k8s.io,resources=clusters,verbs=get;list;watch
// +kubebuilder:rbac:groups=cluster.x-k8s.io,resources=machinedeployments;machinesets,verbs=get;list;watch
// +kubebuilder:rbac:groups="",resources=events,verbs=get;list;watch;create;update;patch

// Reconcile populates capacity information for AWSMachineTemplate.
func (r *AWSMachineTemplateReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	log := logger.FromContext(ctx)

	// Fetch the AWSMachineTemplate
	awsMachineTemplate := &infrav1.AWSMachineTemplate{}
	if err := r.Get(ctx, req.NamespacedName, awsMachineTemplate); err != nil {
		if apierrors.IsNotFound(err) {
			return ctrl.Result{}, nil
		}
		return ctrl.Result{}, err
	}

	// Get instance type from spec
	instanceType := awsMachineTemplate.Spec.Template.Spec.InstanceType
	if instanceType == "" {
		return ctrl.Result{}, nil
	}

	// Check if capacity and nodeInfo are already populated
	// This avoids unnecessary AWS API calls when the status is already populated
	if len(awsMachineTemplate.Status.Capacity) > 0 &&
		awsMachineTemplate.Status.NodeInfo != nil && awsMachineTemplate.Status.NodeInfo.OperatingSystem != "" && awsMachineTemplate.Status.NodeInfo.Architecture != "" {
		return ctrl.Result{}, nil
	}

	// Get the owner cluster
	cluster, err := util.GetOwnerCluster(ctx, r.Client, awsMachineTemplate.ObjectMeta)
	if err != nil {
		return ctrl.Result{}, err
	}
	if cluster == nil {
		return ctrl.Result{}, nil
	}

	// Check if the resource is paused
	if isPaused, conditionChanged, err := paused.EnsurePausedCondition(ctx, r.Client, cluster, awsMachineTemplate); err != nil || isPaused || conditionChanged {
		return ctrl.Result{}, err
	}

	// Find the region by checking ownerReferences
	region, err := r.getRegion(ctx, cluster)
	if err != nil {
		return ctrl.Result{}, err
	}
	if region == "" {
		return ctrl.Result{}, nil
	}

	// Create global scope for this region
	// Reference: exp/instancestate/awsinstancestate_controller.go:68-76
	globalScope, err := scope.NewGlobalScope(scope.GlobalScopeParams{
		ControllerName: "awsmachinetemplate",
		Region:         region,
	})
	if err != nil {
		record.Warnf(awsMachineTemplate, "AWSSessionFailed", "Failed to create AWS session for region %q: %v", region, err)
		return ctrl.Result{}, nil
	}

	// Create EC2 client from global scope
	ec2Client := ec2.NewFromConfig(globalScope.Session())

	// Query instance type capacity
	capacity, err := r.getInstanceTypeCapacity(ctx, ec2Client, instanceType)
	if err != nil {
		record.Warnf(awsMachineTemplate, "CapacityQueryFailed", "Failed to query capacity for instance type %q: %v", instanceType, err)
		return ctrl.Result{}, nil
	}

	// Query node info (architecture and OS)
	nodeInfo, err := r.getNodeInfo(ctx, ec2Client, awsMachineTemplate, instanceType)
	if err != nil {
		record.Warnf(awsMachineTemplate, "NodeInfoQueryFailed", "Failed to query node info for instance type %q: %v", instanceType, err)
		return ctrl.Result{}, nil
	}

	// Save original before modifying, then update all status fields at once
	original := awsMachineTemplate.DeepCopy()
	if len(capacity) > 0 {
		awsMachineTemplate.Status.Capacity = capacity
	}
	if nodeInfo != nil && (nodeInfo.Architecture != "" || nodeInfo.OperatingSystem != "") {
		awsMachineTemplate.Status.NodeInfo = nodeInfo
	}
	if err := r.Status().Patch(ctx, awsMachineTemplate, client.MergeFrom(original)); err != nil {
		return ctrl.Result{}, errors.Wrap(err, "failed to update AWSMachineTemplate status")
	}

	log.Info("Successfully populated capacity and nodeInfo", "instanceType", instanceType, "region", region, "capacity", capacity, "nodeInfo", nodeInfo)
	return ctrl.Result{}, nil
}

// getRegion finds the region by checking the template's owner cluster reference.
func (r *AWSMachineTemplateReconciler) getRegion(ctx context.Context, cluster *clusterv1.Cluster) (string, error) {
	if cluster == nil {
		return "", errors.New("no owner cluster found")
	}

	// Get region from AWSCluster (standard EC2-based cluster)
	if cluster.Spec.InfrastructureRef.IsDefined() && cluster.Spec.InfrastructureRef.Kind == "AWSCluster" {
		awsCluster := &infrav1.AWSCluster{}
		if err := r.Get(ctx, client.ObjectKey{
			Namespace: cluster.Namespace,
			Name:      cluster.Spec.InfrastructureRef.Name,
		}, awsCluster); err != nil {
			if !apierrors.IsNotFound(err) {
				return "", errors.Wrapf(err, "failed to get AWSCluster %s/%s", cluster.Namespace, cluster.Spec.InfrastructureRef.Name)
			}
		} else if awsCluster.Spec.Region != "" {
			return awsCluster.Spec.Region, nil
		}
	}

	// Get region from AWSManagedControlPlane (EKS cluster)
	if cluster.Spec.ControlPlaneRef.IsDefined() && cluster.Spec.ControlPlaneRef.Kind == "AWSManagedControlPlane" {
		awsManagedCP := &ekscontrolplanev1.AWSManagedControlPlane{}
		if err := r.Get(ctx, client.ObjectKey{
			Namespace: cluster.Namespace,
			Name:      cluster.Spec.ControlPlaneRef.Name,
		}, awsManagedCP); err != nil {
			if !apierrors.IsNotFound(err) {
				return "", errors.Wrapf(err, "failed to get AWSManagedControlPlane %s/%s", cluster.Namespace, cluster.Spec.ControlPlaneRef.Name)
			}
		} else if awsManagedCP.Spec.Region != "" {
			return awsManagedCP.Spec.Region, nil
		}
	}

	return "", nil
}

// getInstanceTypeCapacity queries AWS EC2 API for instance type capacity information.
// Returns the resource list (CPU, Memory).
func (r *AWSMachineTemplateReconciler) getInstanceTypeCapacity(ctx context.Context, ec2Client *ec2.Client, instanceType string) (corev1.ResourceList, error) {
	// Query instance type information
	input := &ec2.DescribeInstanceTypesInput{
		InstanceTypes: []ec2types.InstanceType{ec2types.InstanceType(instanceType)},
	}

	result, err := ec2Client.DescribeInstanceTypes(ctx, input)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to describe instance type %q", instanceType)
	}

	if len(result.InstanceTypes) == 0 {
		return nil, errors.Errorf("no information found for instance type %q", instanceType)
	}

	// Extract capacity information
	info := result.InstanceTypes[0]
	resourceList := corev1.ResourceList{}

	// CPU
	if info.VCpuInfo != nil && info.VCpuInfo.DefaultVCpus != nil {
		resourceList[corev1.ResourceCPU] = *resource.NewQuantity(int64(*info.VCpuInfo.DefaultVCpus), resource.DecimalSI)
	}

	// Memory
	if info.MemoryInfo != nil && info.MemoryInfo.SizeInMiB != nil {
		resourceList[corev1.ResourceMemory] = resource.MustParse(fmt.Sprintf("%dMi", *info.MemoryInfo.SizeInMiB))
	}

	return resourceList, nil
}

// getNodeInfo queries node information (architecture and OS) for the AWSMachineTemplate.
// It attempts to resolve nodeInfo using three strategies in order of priority:
//  1. Directly from explicitly specified AMI ID
//  2. From default AMI lookup (requires Kubernetes version from owner MachineDeployment/KubeadmControlPlane)
//  3. From instance type architecture (OS cannot be determined, only architecture)
func (r *AWSMachineTemplateReconciler) getNodeInfo(ctx context.Context, ec2Client *ec2.Client, template *infrav1.AWSMachineTemplate, instanceType string) (*infrav1.NodeInfo, error) {
	// Strategy 1: Extract nodeInfo from the AMI if an ID is set.
	if amiID := ptr.Deref(template.Spec.Template.Spec.AMI.ID, ""); amiID != "" {
		result, err := ec2Client.DescribeImages(ctx, &ec2.DescribeImagesInput{
			ImageIds: []string{amiID},
		})
		if err != nil {
			return nil, errors.Wrapf(err, "failed to describe AMI %q", amiID)
		}
		if len(result.Images) == 0 {
			return nil, errors.Errorf("no information found for AMI %q", amiID)
		}
		// Extract nodeInfo directly from the image object (no additional API call needed)
		return r.extractNodeInfoFromImage(result.Images[0]), nil
	}

	// No explicit AMI ID specified, query instance type to determine architecture
	// This architecture will be used to lookup default AMI (Strategy 2) or as fallback (Strategy 3)
	result, err := ec2Client.DescribeInstanceTypes(ctx, &ec2.DescribeInstanceTypesInput{
		InstanceTypes: []ec2types.InstanceType{ec2types.InstanceType(instanceType)},
	})
	if err != nil {
		return nil, errors.Wrapf(err, "failed to describe instance type %q", instanceType)
	}

	if len(result.InstanceTypes) == 0 {
		return nil, errors.Errorf("no information found for instance type %q", instanceType)
	}

	instanceTypeInfo := result.InstanceTypes[0]

	// Instance type must support exactly one architecture
	if instanceTypeInfo.ProcessorInfo == nil || len(instanceTypeInfo.ProcessorInfo.SupportedArchitectures) != 1 {
		return nil, errors.Errorf("instance type must support exactly one architecture, got %d", len(instanceTypeInfo.ProcessorInfo.SupportedArchitectures))
	}

	// Map EC2 architecture type to architecture tag for AMI lookup
	var architecture string
	var nodeInfoArch infrav1.Architecture
	switch instanceTypeInfo.ProcessorInfo.SupportedArchitectures[0] {
	case ec2types.ArchitectureTypeX8664:
		architecture = ec2service.Amd64ArchitectureTag
		nodeInfoArch = infrav1.ArchitectureAmd64
	case ec2types.ArchitectureTypeArm64:
		architecture = ec2service.Arm64ArchitectureTag
		nodeInfoArch = infrav1.ArchitectureArm64
	default:
		return nil, errors.Errorf("unsupported architecture: %v", instanceTypeInfo.ProcessorInfo.SupportedArchitectures[0])
	}

	// Strategy 2: Try to get Kubernetes version and lookup default AMI
	kubernetesVersion, err := r.getKubernetesVersion(ctx, template)
	if err == nil && kubernetesVersion != "" {
		// Attempt AMI lookup with the version
		image, err := ec2service.DefaultAMILookup(
			ec2Client,
			template.Spec.Template.Spec.ImageLookupOrg,
			template.Spec.Template.Spec.ImageLookupBaseOS,
			kubernetesVersion,
			architecture,
			template.Spec.Template.Spec.ImageLookupFormat,
		)
		if err == nil && image != nil {
			// Successfully found AMI, extract accurate nodeInfo from it
			return r.extractNodeInfoFromImage(*image), nil
		}
		// AMI lookup failed, fall through to Strategy 3
	}

	// Strategy 3: Fallback to instance type architecture only
	// Note: OS cannot be determined from instance type alone, only architecture
	return &infrav1.NodeInfo{
		Architecture: nodeInfoArch,
	}, nil
}

// extractNodeInfoFromImage extracts nodeInfo (architecture and OS) from an EC2 image.
// This is a pure function with no AWS API calls.
func (r *AWSMachineTemplateReconciler) extractNodeInfoFromImage(image ec2types.Image) *infrav1.NodeInfo {
	nodeInfo := &infrav1.NodeInfo{}

	// Extract architecture from AMI
	switch image.Architecture {
	case ec2types.ArchitectureValuesX8664:
		nodeInfo.Architecture = infrav1.ArchitectureAmd64
	case ec2types.ArchitectureValuesArm64:
		nodeInfo.Architecture = infrav1.ArchitectureArm64
	}

	// Determine OS - default to Linux, change to Windows if detected
	// Most AMIs are Linux-based, so we initialize with Linux as the default
	nodeInfo.OperatingSystem = infrav1.OperatingSystemLinux

	// Check Platform field (most reliable for Windows detection)
	if image.Platform == ec2types.PlatformValuesWindows {
		nodeInfo.OperatingSystem = infrav1.OperatingSystemWindows
		return nodeInfo
	}

	// Check PlatformDetails field for Windows indication
	if image.PlatformDetails != nil {
		platformDetails := strings.ToLower(*image.PlatformDetails)
		if strings.Contains(platformDetails, string(infrav1.OperatingSystemWindows)) {
			nodeInfo.OperatingSystem = infrav1.OperatingSystemWindows
		}
	}

	return nodeInfo
}

// getKubernetesVersion attempts to find the Kubernetes version by querying MachineDeployments
// or KubeadmControlPlanes that reference this AWSMachineTemplate.
func (r *AWSMachineTemplateReconciler) getKubernetesVersion(ctx context.Context, template *infrav1.AWSMachineTemplate) (string, error) {
	listOpts, err := getParentListOptions(template.ObjectMeta)
	if err != nil {
		return "", errors.Wrap(err, "failed to get parent list options")
	}

	// Try to find version from MachineSet first
	machineSetList := &clusterv1.MachineSetList{}
	if err := r.List(ctx, machineSetList, listOpts...); err != nil {
		return "", errors.Wrap(err, "failed to list MachineSets")
	}

	// Find MachineSets that reference this AWSMachineTemplate
	for _, ms := range machineSetList.Items {
		if ms.Spec.Template.Spec.InfrastructureRef.Kind == awsMachineTemplateKind &&
			ms.Spec.Template.Spec.InfrastructureRef.Name == template.Name &&
			ms.Spec.Template.Spec.Version != "" {
			return ms.Spec.Template.Spec.Version, nil
		}
	}

	// If not found, try MachineDeployment.
	machineDeploymentList := &clusterv1.MachineDeploymentList{}
	if err := r.List(ctx, machineDeploymentList, listOpts...); err != nil {
		return "", errors.Wrap(err, "failed to list MachineDeployments")
	}

	// Find MachineDeployments that reference this AWSMachineTemplate
	for _, md := range machineDeploymentList.Items {
		if md.Spec.Template.Spec.InfrastructureRef.Kind == awsMachineTemplateKind &&
			md.Spec.Template.Spec.InfrastructureRef.Name == template.Name &&
			md.Spec.Template.Spec.Version != "" {
			return md.Spec.Template.Spec.Version, nil
		}
	}

	// If not found, try KubeadmControlPlane
	kcpList := &controlplanev1.KubeadmControlPlaneList{}
	if err := r.List(ctx, kcpList, listOpts...); err != nil {
		return "", errors.Wrap(err, "failed to list KubeadmControlPlanes")
	}

	// Find KubeadmControlPlanes that reference this AWSMachineTemplate
	for _, kcp := range kcpList.Items {
		if kcp.Spec.MachineTemplate.Spec.InfrastructureRef.Kind == awsMachineTemplateKind &&
			kcp.Spec.MachineTemplate.Spec.InfrastructureRef.Name == template.Name &&
			kcp.Spec.Version != "" {
			return kcp.Spec.Version, nil
		}
	}

	return "", errors.New("no MachineDeployment or KubeadmControlPlane found referencing this AWSMachineTemplate with a version")
}

func getParentListOptions(obj metav1.ObjectMeta) ([]client.ListOption, error) {
	listOpts := []client.ListOption{
		client.InNamespace(obj.Namespace),
	}

	for _, ref := range obj.GetOwnerReferences() {
		if ref.Kind != "Cluster" {
			continue
		}
		gv, err := schema.ParseGroupVersion(ref.APIVersion)
		if err != nil {
			return nil, errors.WithStack(err)
		}
		if gv.Group == clusterv1.GroupVersion.Group {
			listOpts = append(listOpts, client.MatchingLabels{
				clusterv1.ClusterNameLabel: ref.Name,
			})
			break
		}
	}
	return listOpts, nil
}

// kubeadmControlPlaneToAWSMachineTemplate maps KubeadmControlPlane to AWSMachineTemplate reconcile requests.
// This enables the controller to reconcile AWSMachineTemplate when its owner KubeadmControlPlane is created or updated,
// ensuring that nodeInfo can be populated even if the cache hasn't synced yet.
func (r *AWSMachineTemplateReconciler) kubeadmControlPlaneToAWSMachineTemplate(ctx context.Context, o client.Object) []ctrl.Request {
	kcp, ok := o.(*controlplanev1.KubeadmControlPlane)
	if !ok {
		return nil
	}

	// Check if it references an AWSMachineTemplate
	if kcp.Spec.MachineTemplate.Spec.InfrastructureRef.Kind != awsMachineTemplateKind {
		return nil
	}

	// Return reconcile request for the referenced AWSMachineTemplate
	return []ctrl.Request{
		{
			NamespacedName: client.ObjectKey{
				Namespace: kcp.Namespace,
				Name:      kcp.Spec.MachineTemplate.Spec.InfrastructureRef.Name,
			},
		},
	}
}

// machineDeploymentToAWSMachineTemplate maps MachineDeployment to AWSMachineTemplate reconcile requests.
// This enables the controller to reconcile AWSMachineTemplate when its owner MachineDeployment is created or updated,
// ensuring that nodeInfo can be populated even if the cache hasn't synced yet.
func (r *AWSMachineTemplateReconciler) machineDeploymentToAWSMachineTemplate(ctx context.Context, o client.Object) []ctrl.Request {
	md, ok := o.(*clusterv1.MachineDeployment)
	if !ok {
		return nil
	}

	// Check if it references an AWSMachineTemplate
	if md.Spec.Template.Spec.InfrastructureRef.Kind != awsMachineTemplateKind {
		return nil
	}

	// Return reconcile request for the referenced AWSMachineTemplate
	return []ctrl.Request{
		{
			NamespacedName: client.ObjectKey{
				Namespace: md.Namespace,
				Name:      md.Spec.Template.Spec.InfrastructureRef.Name,
			},
		},
	}
}

// machineSetToAWSMachineTemplate maps MachineSet to AWSMachineTemplate reconcile requests.
// This enables the controller to reconcile AWSMachineTemplate when its owner MachineSet is created or updated,
// ensuring that nodeInfo can be populated even if the cache hasn't synced yet.
func (r *AWSMachineTemplateReconciler) machineSetToAWSMachineTemplate(ctx context.Context, o client.Object) []ctrl.Request {
	md, ok := o.(*clusterv1.MachineSet)
	if !ok {
		return nil
	}

	// Check if it references an AWSMachineTemplate
	if md.Spec.Template.Spec.InfrastructureRef.Kind != awsMachineTemplateKind {
		return nil
	}

	// Return reconcile request for the referenced AWSMachineTemplate
	return []ctrl.Request{
		{
			NamespacedName: client.ObjectKey{
				Namespace: md.Namespace,
				Name:      md.Spec.Template.Spec.InfrastructureRef.Name,
			},
		},
	}
}

var resourceCreatedPredicate = predicate.Funcs{
	CreateFunc:  func(e event.CreateEvent) bool { return true },
	UpdateFunc:  func(e event.UpdateEvent) bool { return false },
	DeleteFunc:  func(e event.DeleteEvent) bool { return false },
	GenericFunc: func(e event.GenericEvent) bool { return true },
}
