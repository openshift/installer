/*
Copyright 2020 The Kubernetes Authors.

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

	"github.com/go-logr/logr"
	"github.com/google/go-cmp/cmp"
	"github.com/pkg/errors"
	corev1 "k8s.io/api/core/v1"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/client-go/tools/record"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
	"sigs.k8s.io/controller-runtime/pkg/handler"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
	"sigs.k8s.io/controller-runtime/pkg/source"

	infrav1 "sigs.k8s.io/cluster-api-provider-aws/api/v1beta1"
	"sigs.k8s.io/cluster-api-provider-aws/controllers"
	ekscontrolplanev1 "sigs.k8s.io/cluster-api-provider-aws/controlplane/eks/api/v1beta1"
	expinfrav1 "sigs.k8s.io/cluster-api-provider-aws/exp/api/v1beta1"
	"sigs.k8s.io/cluster-api-provider-aws/pkg/cloud"
	"sigs.k8s.io/cluster-api-provider-aws/pkg/cloud/scope"
	"sigs.k8s.io/cluster-api-provider-aws/pkg/cloud/services"
	asg "sigs.k8s.io/cluster-api-provider-aws/pkg/cloud/services/autoscaling"
	"sigs.k8s.io/cluster-api-provider-aws/pkg/cloud/services/ec2"
	"sigs.k8s.io/cluster-api-provider-aws/pkg/cloud/services/userdata"
	clusterv1 "sigs.k8s.io/cluster-api/api/v1beta1"
	expclusterv1 "sigs.k8s.io/cluster-api/exp/api/v1beta1"
	"sigs.k8s.io/cluster-api/util"
	"sigs.k8s.io/cluster-api/util/conditions"
	"sigs.k8s.io/cluster-api/util/predicates"
)

// AWSMachinePoolReconciler reconciles a AWSMachinePool object.
type AWSMachinePoolReconciler struct {
	client.Client
	Recorder          record.EventRecorder
	WatchFilterValue  string
	asgServiceFactory func(cloud.ClusterScoper) services.ASGInterface
	ec2ServiceFactory func(scope.EC2Scope) services.EC2Interface
}

func (r *AWSMachinePoolReconciler) getASGService(scope cloud.ClusterScoper) services.ASGInterface {
	if r.asgServiceFactory != nil {
		return r.asgServiceFactory(scope)
	}
	return asg.NewService(scope)
}

func (r *AWSMachinePoolReconciler) getEC2Service(scope scope.EC2Scope) services.EC2Interface {
	if r.ec2ServiceFactory != nil {
		return r.ec2ServiceFactory(scope)
	}

	return ec2.NewService(scope)
}

// +kubebuilder:rbac:groups=infrastructure.cluster.x-k8s.io,resources=awsmachinepools,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=infrastructure.cluster.x-k8s.io,resources=awsmachinepools/status,verbs=get;update;patch
// +kubebuilder:rbac:groups=cluster.x-k8s.io,resources=machinepools;machinepools/status,verbs=get;list;watch;patch
// +kubebuilder:rbac:groups="",resources=events,verbs=get;list;watch;create;update;patch
// +kubebuilder:rbac:groups="",resources=secrets;,verbs=get;list;watch
// +kubebuilder:rbac:groups="",resources=namespaces,verbs=get;list;watch

// Reconcile is the reconciliation loop for AWSMachinePool.
func (r *AWSMachinePoolReconciler) Reconcile(ctx context.Context, req ctrl.Request) (_ ctrl.Result, reterr error) {
	log := ctrl.LoggerFrom(ctx)

	// Fetch the AWSMachinePool .
	awsMachinePool := &expinfrav1.AWSMachinePool{}
	err := r.Get(ctx, req.NamespacedName, awsMachinePool)
	if err != nil {
		if apierrors.IsNotFound(err) {
			return ctrl.Result{}, nil
		}
		return ctrl.Result{}, err
	}

	// Fetch the CAPI MachinePool
	machinePool, err := getOwnerMachinePool(ctx, r.Client, awsMachinePool.ObjectMeta)
	if err != nil {
		return reconcile.Result{}, err
	}
	if machinePool == nil {
		log.Info("MachinePool Controller has not yet set OwnerRef")
		return reconcile.Result{}, nil
	}
	log = log.WithValues("machinePool", machinePool.Name)

	// Fetch the Cluster.
	cluster, err := util.GetClusterFromMetadata(ctx, r.Client, machinePool.ObjectMeta)
	if err != nil {
		log.Info("MachinePool is missing cluster label or cluster does not exist")
		return reconcile.Result{}, nil
	}

	log = log.WithValues("cluster", cluster.Name)

	infraCluster, err := r.getInfraCluster(ctx, log, cluster, awsMachinePool)
	if err != nil {
		return ctrl.Result{}, errors.New("error getting infra provider cluster or control plane object")
	}
	if infraCluster == nil {
		log.Info("AWSCluster or AWSManagedControlPlane is not ready yet")
		return ctrl.Result{}, nil
	}

	// Create the machine pool scope
	machinePoolScope, err := scope.NewMachinePoolScope(scope.MachinePoolScopeParams{
		Client:         r.Client,
		Cluster:        cluster,
		MachinePool:    machinePool,
		InfraCluster:   infraCluster,
		AWSMachinePool: awsMachinePool,
	})
	if err != nil {
		log.Error(err, "failed to create scope")
		return ctrl.Result{}, err
	}

	// Always close the scope when exiting this function so we can persist any AWSMachine changes.
	defer func() {
		// set Ready condition before AWSMachinePool is patched
		conditions.SetSummary(machinePoolScope.AWSMachinePool,
			conditions.WithConditions(
				expinfrav1.ASGReadyCondition,
				expinfrav1.LaunchTemplateReadyCondition,
			),
			conditions.WithStepCounterIfOnly(
				expinfrav1.ASGReadyCondition,
				expinfrav1.LaunchTemplateReadyCondition,
			),
		)

		if err := machinePoolScope.Close(); err != nil && reterr == nil {
			reterr = err
		}
	}()

	switch infraScope := infraCluster.(type) {
	case *scope.ManagedControlPlaneScope:
		if !awsMachinePool.ObjectMeta.DeletionTimestamp.IsZero() {
			return r.reconcileDelete(machinePoolScope, infraScope, infraScope)
		}

		return r.reconcileNormal(ctx, machinePoolScope, infraScope, infraScope)
	case *scope.ClusterScope:
		if !awsMachinePool.ObjectMeta.DeletionTimestamp.IsZero() {
			return r.reconcileDelete(machinePoolScope, infraScope, infraScope)
		}

		return r.reconcileNormal(ctx, machinePoolScope, infraScope, infraScope)
	default:
		return ctrl.Result{}, errors.New("infraCluster has unknown type")
	}
}

func (r *AWSMachinePoolReconciler) SetupWithManager(ctx context.Context, mgr ctrl.Manager, options controller.Options) error {
	return ctrl.NewControllerManagedBy(mgr).
		WithOptions(options).
		For(&expinfrav1.AWSMachinePool{}).
		Watches(
			&source.Kind{Type: &expclusterv1.MachinePool{}},
			handler.EnqueueRequestsFromMapFunc(machinePoolToInfrastructureMapFunc(expinfrav1.GroupVersion.WithKind("AWSMachinePool"))),
		).
		WithEventFilter(predicates.ResourceNotPausedAndHasFilterLabel(ctrl.LoggerFrom(ctx), r.WatchFilterValue)).
		Complete(r)
}

func (r *AWSMachinePoolReconciler) reconcileNormal(ctx context.Context, machinePoolScope *scope.MachinePoolScope, clusterScope cloud.ClusterScoper, ec2Scope scope.EC2Scope) (ctrl.Result, error) {
	clusterScope.Info("Reconciling AWSMachinePool")

	// If the AWSMachine is in an error state, return early.
	if machinePoolScope.HasFailed() {
		machinePoolScope.Info("Error state detected, skipping reconciliation")

		// TODO: If we are in a failed state, delete the secret regardless of instance state

		return ctrl.Result{}, nil
	}

	// If the AWSMachinepool doesn't have our finalizer, add it
	controllerutil.AddFinalizer(machinePoolScope.AWSMachinePool, expinfrav1.MachinePoolFinalizer)

	// Register finalizer immediately to avoid orphaning AWS resources
	if err := machinePoolScope.PatchObject(); err != nil {
		return ctrl.Result{}, err
	}

	if !machinePoolScope.Cluster.Status.InfrastructureReady {
		machinePoolScope.Info("Cluster infrastructure is not ready yet")
		conditions.MarkFalse(machinePoolScope.AWSMachinePool, expinfrav1.ASGReadyCondition, infrav1.WaitingForClusterInfrastructureReason, clusterv1.ConditionSeverityInfo, "")
		return ctrl.Result{}, nil
	}

	// Make sure bootstrap data is available and populated
	if machinePoolScope.MachinePool.Spec.Template.Spec.Bootstrap.DataSecretName == nil {
		machinePoolScope.Info("Bootstrap data secret reference is not yet available")
		conditions.MarkFalse(machinePoolScope.AWSMachinePool, expinfrav1.ASGReadyCondition, infrav1.WaitingForBootstrapDataReason, clusterv1.ConditionSeverityInfo, "")
		return ctrl.Result{}, nil
	}

	if err := r.reconcileLaunchTemplate(machinePoolScope, ec2Scope); err != nil {
		r.Recorder.Eventf(machinePoolScope.AWSMachinePool, corev1.EventTypeWarning, "FailedLaunchTemplateReconcile", "Failed to reconcile launch template: %v", err)
		machinePoolScope.Error(err, "failed to reconcile launch template")
		return ctrl.Result{}, err
	}

	// set the LaunchTemplateReady condition
	conditions.MarkTrue(machinePoolScope.AWSMachinePool, expinfrav1.LaunchTemplateReadyCondition)

	// Initialize ASG client
	asgsvc := r.getASGService(clusterScope)

	// Find existing ASG
	asg, err := r.findASG(machinePoolScope, asgsvc)
	if err != nil {
		conditions.MarkUnknown(machinePoolScope.AWSMachinePool, expinfrav1.ASGReadyCondition, expinfrav1.ASGNotFoundReason, err.Error())
		return ctrl.Result{}, err
	}

	if asg == nil {
		// Create new ASG
		if _, err := r.createPool(machinePoolScope, clusterScope); err != nil {
			conditions.MarkFalse(machinePoolScope.AWSMachinePool, expinfrav1.ASGReadyCondition, expinfrav1.ASGProvisionFailedReason, clusterv1.ConditionSeverityError, err.Error())
			return ctrl.Result{}, err
		}
		return ctrl.Result{}, nil
	}

	if scope.ReplicasExternallyManaged(machinePoolScope.MachinePool) {
		// Set MachinePool replicas to the ASG DesiredCapacity
		if *machinePoolScope.MachinePool.Spec.Replicas != *asg.DesiredCapacity {
			machinePoolScope.Info("Setting MachinePool replicas to ASG DesiredCapacity",
				"local", machinePoolScope.MachinePool.Spec.Replicas,
				"external", asg.DesiredCapacity)
			machinePoolScope.MachinePool.Spec.Replicas = asg.DesiredCapacity
			if err := machinePoolScope.PatchCAPIMachinePoolObject(ctx); err != nil {
				return ctrl.Result{}, err
			}
		}
	}

	if err := r.updatePool(machinePoolScope, clusterScope, asg); err != nil {
		machinePoolScope.Error(err, "error updating AWSMachinePool")
		return ctrl.Result{}, err
	}

	err = r.reconcileTags(machinePoolScope, clusterScope, ec2Scope)
	if err != nil {
		return ctrl.Result{}, errors.Wrap(err, "error updating tags")
	}

	// Make sure Spec.ProviderID is always set.
	machinePoolScope.AWSMachinePool.Spec.ProviderID = asg.ID
	providerIDList := make([]string, len(asg.Instances))

	for i, ec2 := range asg.Instances {
		providerIDList[i] = fmt.Sprintf("aws:///%s/%s", ec2.AvailabilityZone, ec2.ID)
	}

	machinePoolScope.SetAnnotation("cluster-api-provider-aws", "true")

	machinePoolScope.AWSMachinePool.Spec.ProviderIDList = providerIDList
	machinePoolScope.AWSMachinePool.Status.Replicas = int32(len(providerIDList))
	machinePoolScope.AWSMachinePool.Status.Ready = true
	conditions.MarkTrue(machinePoolScope.AWSMachinePool, expinfrav1.ASGReadyCondition)

	err = machinePoolScope.UpdateInstanceStatuses(ctx, asg.Instances)
	if err != nil {
		machinePoolScope.Info("Failed updating instances", "instances", asg.Instances)
	}

	return ctrl.Result{}, nil
}

func (r *AWSMachinePoolReconciler) reconcileDelete(machinePoolScope *scope.MachinePoolScope, clusterScope cloud.ClusterScoper, ec2Scope scope.EC2Scope) (ctrl.Result, error) {
	clusterScope.Info("Handling deleted AWSMachinePool")

	ec2Svc := r.getEC2Service(ec2Scope)
	asgSvc := r.getASGService(clusterScope)

	asg, err := r.findASG(machinePoolScope, asgSvc)
	if err != nil {
		return ctrl.Result{}, err
	}

	if asg == nil {
		machinePoolScope.V(2).Info("Unable to locate ASG")
		r.Recorder.Eventf(machinePoolScope.AWSMachinePool, corev1.EventTypeNormal, "NoASGFound", "Unable to find matching ASG")
	} else {
		machinePoolScope.SetASGStatus(asg.Status)
		switch asg.Status {
		case expinfrav1.ASGStatusDeleteInProgress:
			// ASG is already deleting
			machinePoolScope.SetNotReady()
			conditions.MarkFalse(machinePoolScope.AWSMachinePool, expinfrav1.ASGReadyCondition, expinfrav1.ASGDeletionInProgress, clusterv1.ConditionSeverityWarning, "")
			r.Recorder.Eventf(machinePoolScope.AWSMachinePool, corev1.EventTypeWarning, "DeletionInProgress", "ASG deletion in progress: %q", asg.Name)
			machinePoolScope.Info("ASG is already deleting", "name", asg.Name)
		default:
			machinePoolScope.Info("Deleting ASG", "id", asg.Name, "status", asg.Status)
			if err := asgSvc.DeleteASGAndWait(asg.Name); err != nil {
				r.Recorder.Eventf(machinePoolScope.AWSMachinePool, corev1.EventTypeWarning, "FailedDelete", "Failed to delete ASG %q: %v", asg.Name, err)
				return ctrl.Result{}, errors.Wrap(err, "failed to delete ASG")
			}
		}
	}

	launchTemplateID := machinePoolScope.AWSMachinePool.Status.LaunchTemplateID
	launchTemplate, _, err := ec2Svc.GetLaunchTemplate(machinePoolScope.Name())
	if err != nil {
		return ctrl.Result{}, err
	}

	if launchTemplate == nil {
		machinePoolScope.V(2).Info("Unable to locate launch template")
		r.Recorder.Eventf(machinePoolScope.AWSMachinePool, corev1.EventTypeNormal, "NoASGFound", "Unable to find matching ASG")
		controllerutil.RemoveFinalizer(machinePoolScope.AWSMachinePool, expinfrav1.MachinePoolFinalizer)
		return ctrl.Result{}, nil
	}

	machinePoolScope.Info("deleting launch template", "name", launchTemplate.Name)
	if err := ec2Svc.DeleteLaunchTemplate(launchTemplateID); err != nil {
		r.Recorder.Eventf(machinePoolScope.AWSMachinePool, corev1.EventTypeWarning, "FailedDelete", "Failed to delete launch template %q: %v", launchTemplate.Name, err)
		return ctrl.Result{}, errors.Wrap(err, "failed to delete ASG")
	}

	machinePoolScope.Info("successfully deleted AutoScalingGroup and Launch Template")

	// remove finalizer
	controllerutil.RemoveFinalizer(machinePoolScope.AWSMachinePool, expinfrav1.MachinePoolFinalizer)

	return ctrl.Result{}, nil
}

func (r *AWSMachinePoolReconciler) updatePool(machinePoolScope *scope.MachinePoolScope, clusterScope cloud.ClusterScoper, existingASG *expinfrav1.AutoScalingGroup) error {
	if asgNeedsUpdates(machinePoolScope, existingASG) {
		machinePoolScope.Info("updating AutoScalingGroup")
		asgSvc := r.getASGService(clusterScope)

		if err := asgSvc.UpdateASG(machinePoolScope); err != nil {
			r.Recorder.Eventf(machinePoolScope.AWSMachinePool, corev1.EventTypeWarning, "FailedUpdate", "Failed to update ASG: %v", err)
			return errors.Wrap(err, "unable to update ASG")
		}
	}

	return nil
}

func (r *AWSMachinePoolReconciler) createPool(machinePoolScope *scope.MachinePoolScope, clusterScope cloud.ClusterScoper) (*expinfrav1.AutoScalingGroup, error) {
	clusterScope.Info("Initializing ASG client")

	asgsvc := r.getASGService(clusterScope)

	machinePoolScope.Info("Creating Autoscaling Group")
	asg, err := asgsvc.CreateASG(machinePoolScope)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to create AWSMachinePool")
	}

	return asg, nil
}

func (r *AWSMachinePoolReconciler) findASG(machinePoolScope *scope.MachinePoolScope, asgsvc services.ASGInterface) (*expinfrav1.AutoScalingGroup, error) {
	// Query the instance using tags.
	asg, err := asgsvc.GetASGByName(machinePoolScope)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to query AWSMachinePool by name")
	}

	return asg, nil
}

func (r *AWSMachinePoolReconciler) reconcileLaunchTemplate(machinePoolScope *scope.MachinePoolScope, ec2Scope scope.EC2Scope) error {
	bootstrapData, err := machinePoolScope.GetRawBootstrapData()
	if err != nil {
		r.Recorder.Eventf(machinePoolScope.AWSMachinePool, corev1.EventTypeWarning, "FailedGetBootstrapData", err.Error())
	}
	bootstrapDataHash := userdata.ComputeHash(bootstrapData)

	ec2svc := r.getEC2Service(ec2Scope)

	machinePoolScope.Info("checking for existing launch template")
	launchTemplate, launchTemplateUserDataHash, err := ec2svc.GetLaunchTemplate(machinePoolScope.Name())
	if err != nil {
		conditions.MarkUnknown(machinePoolScope.AWSMachinePool, expinfrav1.LaunchTemplateReadyCondition, expinfrav1.LaunchTemplateNotFoundReason, err.Error())
		return err
	}

	imageID, err := ec2svc.DiscoverLaunchTemplateAMI(machinePoolScope)
	if err != nil {
		conditions.MarkFalse(machinePoolScope.AWSMachinePool, expinfrav1.LaunchTemplateReadyCondition, expinfrav1.LaunchTemplateCreateFailedReason, clusterv1.ConditionSeverityError, err.Error())
		return err
	}

	if launchTemplate == nil {
		machinePoolScope.Info("no existing launch template found, creating")
		launchTemplateID, err := ec2svc.CreateLaunchTemplate(machinePoolScope, imageID, bootstrapData)
		if err != nil {
			conditions.MarkFalse(machinePoolScope.AWSMachinePool, expinfrav1.LaunchTemplateReadyCondition, expinfrav1.LaunchTemplateCreateFailedReason, clusterv1.ConditionSeverityError, err.Error())
			return err
		}

		machinePoolScope.SetLaunchTemplateIDStatus(launchTemplateID)
		return machinePoolScope.PatchObject()
	}

	// LaunchTemplateID is set during LaunchTemplate creation, but for a scenario such as `clusterctl move`, status fields become blank.
	// If launchTemplate already exists but LaunchTemplateID field in the status is empty, get the ID and update the status.
	if machinePoolScope.AWSMachinePool.Status.LaunchTemplateID == "" {
		launchTemplateID, err := ec2svc.GetLaunchTemplateID(machinePoolScope.Name())
		if err != nil {
			conditions.MarkUnknown(machinePoolScope.AWSMachinePool, expinfrav1.LaunchTemplateReadyCondition, expinfrav1.LaunchTemplateNotFoundReason, err.Error())
			return err
		}
		machinePoolScope.SetLaunchTemplateIDStatus(launchTemplateID)
		return machinePoolScope.PatchObject()
	}

	annotation, err := r.machinePoolAnnotationJSON(machinePoolScope.AWSMachinePool, TagsLastAppliedAnnotation)
	if err != nil {
		return err
	}

	// Check if the instance tags were changed. If they were, create a new LaunchTemplate.
	tagsChanged, _, _, _ := tagsChanged(annotation, machinePoolScope.AdditionalTags()) // nolint:dogsled

	needsUpdate, err := ec2svc.LaunchTemplateNeedsUpdate(machinePoolScope, &machinePoolScope.AWSMachinePool.Spec.AWSLaunchTemplate, launchTemplate)
	if err != nil {
		return err
	}

	// If there is a change: before changing the template, check if there exist an ongoing instance refresh,
	// because only 1 instance refresh can be "InProgress". If template is updated when refresh cannot be started,
	// that change will not trigger a refresh. Do not start an instance refresh if only userdata changed.
	if needsUpdate || tagsChanged || *imageID != *launchTemplate.AMI.ID {
		asgSvc := r.getASGService(ec2Scope)
		canStart, err := asgSvc.CanStartASGInstanceRefresh(machinePoolScope)
		if err != nil {
			return err
		}
		if !canStart {
			conditions.MarkFalse(machinePoolScope.AWSMachinePool, expinfrav1.InstanceRefreshStartedCondition, expinfrav1.InstanceRefreshNotReadyReason, clusterv1.ConditionSeverityWarning, "")
			return errors.New("Cannot start a new instance refresh. Unfinished instance refresh exist")
		}
	}

	// Create a new launch template version if there's a difference in configuration, tags,
	// userdata, OR we've discovered a new AMI ID.
	if needsUpdate || tagsChanged || *imageID != *launchTemplate.AMI.ID || launchTemplateUserDataHash != bootstrapDataHash {
		machinePoolScope.Info("creating new version for launch template", "existing", launchTemplate, "incoming", machinePoolScope.AWSMachinePool.Spec.AWSLaunchTemplate)
		// There is a limit to the number of Launch Template Versions.
		// We ensure that the number of versions does not grow without bound by following a simple rule: Before we create a new version, we delete one old version, if there is at least one old version that is not in use.
		if err := ec2svc.PruneLaunchTemplateVersions(machinePoolScope.AWSMachinePool.Status.LaunchTemplateID); err != nil {
			return err
		}
		if err := ec2svc.CreateLaunchTemplateVersion(machinePoolScope, imageID, bootstrapData); err != nil {
			return err
		}
	}

	// After creating a new version of launch template, instance refresh is required
	// to trigger a rolling replacement of all previously launched instances.
	// If ONLY the userdata changed, previously launched instances continue to use the old launch
	// template.
	//
	// FIXME(dlipovetsky,sedefsavas): If the controller terminates, or the StartASGInstanceRefresh returns an error,
	// this conditional will not evaluate to true the next reconcile. If any machines use an older
	// Launch Template version, and the difference between the older and current versions is _more_
	// than userdata, we should start an Instance Refresh.
	if needsUpdate || tagsChanged || *imageID != *launchTemplate.AMI.ID {
		machinePoolScope.Info("starting instance refresh", "number of instances", machinePoolScope.MachinePool.Spec.Replicas)
		asgSvc := r.getASGService(ec2Scope)
		if err := asgSvc.StartASGInstanceRefresh(machinePoolScope); err != nil {
			conditions.MarkFalse(machinePoolScope.AWSMachinePool, expinfrav1.InstanceRefreshStartedCondition, expinfrav1.InstanceRefreshFailedReason, clusterv1.ConditionSeverityError, err.Error())
			return err
		}
		conditions.MarkTrue(machinePoolScope.AWSMachinePool, expinfrav1.InstanceRefreshStartedCondition)
	}

	return nil
}

func (r *AWSMachinePoolReconciler) reconcileTags(machinePoolScope *scope.MachinePoolScope, clusterScope cloud.ClusterScoper, ec2Scope scope.EC2Scope) error {
	ec2Svc := r.getEC2Service(ec2Scope)
	asgSvc := r.getASGService(clusterScope)

	launchTemplateID := machinePoolScope.AWSMachinePool.Status.LaunchTemplateID
	asgName := machinePoolScope.Name()
	additionalTags := machinePoolScope.AdditionalTags()

	tagsChanged, err := r.ensureTags(ec2Svc, asgSvc, machinePoolScope.AWSMachinePool, &launchTemplateID, &asgName, additionalTags)
	if err != nil {
		return err
	}
	if tagsChanged {
		r.Recorder.Eventf(machinePoolScope.AWSMachinePool, corev1.EventTypeNormal, "UpdatedTags", "updated tags on resources")
	}
	return nil
}

// asgNeedsUpdates compares incoming AWSMachinePool and compares against existing ASG.
func asgNeedsUpdates(machinePoolScope *scope.MachinePoolScope, existingASG *expinfrav1.AutoScalingGroup) bool {
	if !scope.ReplicasExternallyManaged(machinePoolScope.MachinePool) {
		if machinePoolScope.MachinePool.Spec.Replicas != nil {
			if existingASG.DesiredCapacity == nil || *machinePoolScope.MachinePool.Spec.Replicas != *existingASG.DesiredCapacity {
				return true
			}
		} else if existingASG.DesiredCapacity != nil {
			return true
		}
	}

	if machinePoolScope.AWSMachinePool.Spec.MaxSize != existingASG.MaxSize {
		return true
	}

	if machinePoolScope.AWSMachinePool.Spec.MinSize != existingASG.MinSize {
		return true
	}

	if machinePoolScope.AWSMachinePool.Spec.CapacityRebalance != existingASG.CapacityRebalance {
		return true
	}

	if !cmp.Equal(machinePoolScope.AWSMachinePool.Spec.MixedInstancesPolicy, existingASG.MixedInstancesPolicy) {
		machinePoolScope.Info("got a mixed diff here", "incoming", machinePoolScope.AWSMachinePool.Spec.MixedInstancesPolicy, "existing", existingASG.MixedInstancesPolicy)
		return true
	}

	// todo subnet diff

	return false
}

// getOwnerMachinePool returns the MachinePool object owning the current resource.
func getOwnerMachinePool(ctx context.Context, c client.Client, obj metav1.ObjectMeta) (*expclusterv1.MachinePool, error) {
	for _, ref := range obj.OwnerReferences {
		if ref.Kind != "MachinePool" {
			continue
		}
		gv, err := schema.ParseGroupVersion(ref.APIVersion)
		if err != nil {
			return nil, errors.WithStack(err)
		}
		if gv.Group == expclusterv1.GroupVersion.Group {
			return getMachinePoolByName(ctx, c, obj.Namespace, ref.Name)
		}
	}
	return nil, nil
}

// getMachinePoolByName finds and return a Machine object using the specified params.
func getMachinePoolByName(ctx context.Context, c client.Client, namespace, name string) (*expclusterv1.MachinePool, error) {
	m := &expclusterv1.MachinePool{}
	key := client.ObjectKey{Name: name, Namespace: namespace}
	if err := c.Get(ctx, key, m); err != nil {
		return nil, err
	}
	return m, nil
}

func machinePoolToInfrastructureMapFunc(gvk schema.GroupVersionKind) handler.MapFunc {
	return func(o client.Object) []reconcile.Request {
		m, ok := o.(*expclusterv1.MachinePool)
		if !ok {
			panic(fmt.Sprintf("Expected a MachinePool but got a %T", o))
		}

		gk := gvk.GroupKind()
		// Return early if the GroupKind doesn't match what we expect
		infraGK := m.Spec.Template.Spec.InfrastructureRef.GroupVersionKind().GroupKind()
		if gk != infraGK {
			return nil
		}

		return []reconcile.Request{
			{
				NamespacedName: client.ObjectKey{
					Namespace: m.Namespace,
					Name:      m.Spec.Template.Spec.InfrastructureRef.Name,
				},
			},
		}
	}
}

func (r *AWSMachinePoolReconciler) getInfraCluster(ctx context.Context, log logr.Logger, cluster *clusterv1.Cluster, awsMachinePool *expinfrav1.AWSMachinePool) (scope.EC2Scope, error) {
	var clusterScope *scope.ClusterScope
	var managedControlPlaneScope *scope.ManagedControlPlaneScope
	var err error

	if cluster.Spec.ControlPlaneRef != nil && cluster.Spec.ControlPlaneRef.Kind == controllers.AWSManagedControlPlaneRefKind {
		controlPlane := &ekscontrolplanev1.AWSManagedControlPlane{}
		controlPlaneName := client.ObjectKey{
			Namespace: awsMachinePool.Namespace,
			Name:      cluster.Spec.ControlPlaneRef.Name,
		}

		if err := r.Get(ctx, controlPlaneName, controlPlane); err != nil {
			// AWSManagedControlPlane is not ready
			return nil, nil // nolint:nilerr
		}

		managedControlPlaneScope, err = scope.NewManagedControlPlaneScope(scope.ManagedControlPlaneScopeParams{
			Client:         r.Client,
			Logger:         &log,
			Cluster:        cluster,
			ControlPlane:   controlPlane,
			ControllerName: "awsManagedControlPlane",
		})
		if err != nil {
			return nil, err
		}

		return managedControlPlaneScope, nil
	}

	awsCluster := &infrav1.AWSCluster{}

	infraClusterName := client.ObjectKey{
		Namespace: awsMachinePool.Namespace,
		Name:      cluster.Spec.InfrastructureRef.Name,
	}

	if err := r.Client.Get(ctx, infraClusterName, awsCluster); err != nil {
		// AWSCluster is not ready
		return nil, nil // nolint:nilerr
	}

	// Create the cluster scope
	clusterScope, err = scope.NewClusterScope(scope.ClusterScopeParams{
		Client:         r.Client,
		Logger:         &log,
		Cluster:        cluster,
		AWSCluster:     awsCluster,
		ControllerName: "awsmachine",
	})
	if err != nil {
		return nil, err
	}

	return clusterScope, nil
}
