/*
Copyright 2019 The Kubernetes Authors.

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
	"encoding/json"
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	ignTypes "github.com/flatcar-linux/ignition/config/v2_3/types"
	"github.com/go-logr/logr"
	"github.com/google/go-cmp/cmp"
	"github.com/pkg/errors"
	corev1 "k8s.io/api/core/v1"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/client-go/tools/record"
	"k8s.io/utils/pointer"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
	"sigs.k8s.io/controller-runtime/pkg/event"
	"sigs.k8s.io/controller-runtime/pkg/handler"
	"sigs.k8s.io/controller-runtime/pkg/predicate"
	"sigs.k8s.io/controller-runtime/pkg/source"

	infrav1 "sigs.k8s.io/cluster-api-provider-aws/api/v1beta1"
	ekscontrolplanev1 "sigs.k8s.io/cluster-api-provider-aws/controlplane/eks/api/v1beta1"
	"sigs.k8s.io/cluster-api-provider-aws/feature"
	"sigs.k8s.io/cluster-api-provider-aws/pkg/cloud"
	"sigs.k8s.io/cluster-api-provider-aws/pkg/cloud/scope"
	"sigs.k8s.io/cluster-api-provider-aws/pkg/cloud/services"
	"sigs.k8s.io/cluster-api-provider-aws/pkg/cloud/services/ec2"
	"sigs.k8s.io/cluster-api-provider-aws/pkg/cloud/services/elb"
	"sigs.k8s.io/cluster-api-provider-aws/pkg/cloud/services/instancestate"
	"sigs.k8s.io/cluster-api-provider-aws/pkg/cloud/services/s3"
	"sigs.k8s.io/cluster-api-provider-aws/pkg/cloud/services/secretsmanager"
	"sigs.k8s.io/cluster-api-provider-aws/pkg/cloud/services/ssm"
	"sigs.k8s.io/cluster-api-provider-aws/pkg/cloud/services/userdata"
	clusterv1 "sigs.k8s.io/cluster-api/api/v1beta1"
	"sigs.k8s.io/cluster-api/controllers/noderefutil"
	capierrors "sigs.k8s.io/cluster-api/errors"
	"sigs.k8s.io/cluster-api/util"
	"sigs.k8s.io/cluster-api/util/annotations"
	"sigs.k8s.io/cluster-api/util/conditions"
	"sigs.k8s.io/cluster-api/util/predicates"
)

// InstanceIDIndex defines the aws machine controller's instance ID index.
const InstanceIDIndex = ".spec.instanceID"

// AWSMachineReconciler reconciles a AwsMachine object.
type AWSMachineReconciler struct {
	client.Client
	Log                          logr.Logger
	Recorder                     record.EventRecorder
	ec2ServiceFactory            func(scope.EC2Scope) services.EC2Interface
	elbServiceFactory            func(scope.ELBScope) services.ELBInterface
	secretsManagerServiceFactory func(cloud.ClusterScoper) services.SecretInterface
	SSMServiceFactory            func(cloud.ClusterScoper) services.SecretInterface
	objectStoreServiceFactory    func(cloud.ClusterScoper) services.ObjectStoreInterface
	Endpoints                    []scope.ServiceEndpoint
	WatchFilterValue             string
}

const (
	// AWSManagedControlPlaneRefKind is the string value indicating that a cluster is AWS managed.
	AWSManagedControlPlaneRefKind = "AWSManagedControlPlane"
)

func (r *AWSMachineReconciler) getEC2Service(scope scope.EC2Scope) services.EC2Interface {
	if r.ec2ServiceFactory != nil {
		return r.ec2ServiceFactory(scope)
	}

	return ec2.NewService(scope)
}

func (r *AWSMachineReconciler) getSecretsManagerService(scope cloud.ClusterScoper) services.SecretInterface {
	if r.secretsManagerServiceFactory != nil {
		return r.secretsManagerServiceFactory(scope)
	}

	return secretsmanager.NewService(scope)
}

func (r *AWSMachineReconciler) getSSMService(scope cloud.ClusterScoper) services.SecretInterface {
	if r.SSMServiceFactory != nil {
		return r.SSMServiceFactory(scope)
	}
	return ssm.NewService(scope)
}

func (r *AWSMachineReconciler) getSecretService(machineScope *scope.MachineScope, scope cloud.ClusterScoper) (services.SecretInterface, error) {
	switch machineScope.SecureSecretsBackend() {
	case infrav1.SecretBackendSSMParameterStore:
		return r.getSSMService(scope), nil
	case infrav1.SecretBackendSecretsManager:
		return r.getSecretsManagerService(scope), nil
	}
	return nil, errors.New("invalid secret backend")
}

func (r *AWSMachineReconciler) getELBService(elbScope scope.ELBScope) services.ELBInterface {
	if r.elbServiceFactory != nil {
		return r.elbServiceFactory(elbScope)
	}
	return elb.NewService(elbScope)
}

func (r *AWSMachineReconciler) getObjectStoreService(scope scope.S3Scope) services.ObjectStoreInterface {
	if r.objectStoreServiceFactory != nil {
		return r.objectStoreServiceFactory(scope)
	}

	return s3.NewService(scope)
}

// +kubebuilder:rbac:groups=infrastructure.cluster.x-k8s.io,resources=awsmachines,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=infrastructure.cluster.x-k8s.io,resources=awsmachines/status,verbs=get;update;patch
// +kubebuilder:rbac:groups=cluster.x-k8s.io,resources=machines;machines/status,verbs=get;list;watch
// +kubebuilder:rbac:groups="",resources=secrets;,verbs=get;list;watch
// +kubebuilder:rbac:groups="",resources=namespaces,verbs=get;list;watch
// +kubebuilder:rbac:groups="",resources=events,verbs=get;list;watch;create;update;patch

func (r *AWSMachineReconciler) Reconcile(ctx context.Context, req ctrl.Request) (_ ctrl.Result, reterr error) {
	log := ctrl.LoggerFrom(ctx)

	// Fetch the AWSMachine instance.
	awsMachine := &infrav1.AWSMachine{}
	err := r.Get(ctx, req.NamespacedName, awsMachine)
	if err != nil {
		if apierrors.IsNotFound(err) {
			return ctrl.Result{}, nil
		}
		return ctrl.Result{}, err
	}

	// Fetch the Machine.
	machine, err := util.GetOwnerMachine(ctx, r.Client, awsMachine.ObjectMeta)
	if err != nil {
		return ctrl.Result{}, err
	}
	if machine == nil {
		log.Info("Machine Controller has not yet set OwnerRef")
		return ctrl.Result{}, nil
	}

	log = log.WithValues("machine", machine.Name)

	// Fetch the Cluster.
	cluster, err := util.GetClusterFromMetadata(ctx, r.Client, machine.ObjectMeta)
	if err != nil {
		log.Info("Machine is missing cluster label or cluster does not exist")
		return ctrl.Result{}, nil
	}

	if annotations.IsPaused(cluster, awsMachine) {
		log.Info("AWSMachine or linked Cluster is marked as paused. Won't reconcile")
		return ctrl.Result{}, nil
	}

	log = log.WithValues("cluster", cluster.Name)

	infraCluster, err := r.getInfraCluster(ctx, log, cluster, awsMachine)
	if err != nil {
		return ctrl.Result{}, errors.New("error getting infra provider cluster or control plane object")
	}
	if infraCluster == nil {
		log.Info("AWSCluster or AWSManagedControlPlane is not ready yet")
		return ctrl.Result{}, nil
	}

	// Create the machine scope
	machineScope, err := scope.NewMachineScope(scope.MachineScopeParams{
		Client:       r.Client,
		Cluster:      cluster,
		Machine:      machine,
		InfraCluster: infraCluster,
		AWSMachine:   awsMachine,
	})
	if err != nil {
		log.Error(err, "failed to create scope")
		return ctrl.Result{}, err
	}

	// Always close the scope when exiting this function so we can persist any AWSMachine changes.
	defer func() {
		if err := machineScope.Close(); err != nil && reterr == nil {
			reterr = err
		}
	}()

	switch infraScope := infraCluster.(type) {
	case *scope.ManagedControlPlaneScope:
		if !awsMachine.ObjectMeta.DeletionTimestamp.IsZero() {
			return r.reconcileDelete(machineScope, infraScope, infraScope, nil, nil)
		}

		return r.reconcileNormal(ctx, machineScope, infraScope, infraScope, nil, nil)
	case *scope.ClusterScope:
		if !awsMachine.ObjectMeta.DeletionTimestamp.IsZero() {
			return r.reconcileDelete(machineScope, infraScope, infraScope, infraScope, infraScope)
		}

		return r.reconcileNormal(ctx, machineScope, infraScope, infraScope, infraScope, infraScope)
	default:
		return ctrl.Result{}, errors.New("infraCluster has unknown type")
	}
}

func (r *AWSMachineReconciler) SetupWithManager(ctx context.Context, mgr ctrl.Manager, options controller.Options) error {
	log := ctrl.LoggerFrom(ctx)
	AWSClusterToAWSMachines := r.AWSClusterToAWSMachines(log)

	controller, err := ctrl.NewControllerManagedBy(mgr).
		WithOptions(options).
		For(&infrav1.AWSMachine{}).
		Watches(
			&source.Kind{Type: &clusterv1.Machine{}},
			handler.EnqueueRequestsFromMapFunc(util.MachineToInfrastructureMapFunc(infrav1.GroupVersion.WithKind("AWSMachine"))),
		).
		Watches(
			&source.Kind{Type: &infrav1.AWSCluster{}},
			handler.EnqueueRequestsFromMapFunc(AWSClusterToAWSMachines),
		).
		WithEventFilter(predicates.ResourceNotPausedAndHasFilterLabel(log, r.WatchFilterValue)).
		WithEventFilter(
			predicate.Funcs{
				// Avoid reconciling if the event triggering the reconciliation is related to incremental status updates
				// for AWSMachine resources only
				UpdateFunc: func(e event.UpdateEvent) bool {
					if e.ObjectOld.GetObjectKind().GroupVersionKind().Kind != "AWSMachine" {
						return true
					}

					oldMachine := e.ObjectOld.(*infrav1.AWSMachine).DeepCopy()
					newMachine := e.ObjectNew.(*infrav1.AWSMachine).DeepCopy()

					oldMachine.Status = infrav1.AWSMachineStatus{}
					newMachine.Status = infrav1.AWSMachineStatus{}

					oldMachine.ObjectMeta.ResourceVersion = ""
					newMachine.ObjectMeta.ResourceVersion = ""

					return !cmp.Equal(oldMachine, newMachine)
				},
			},
		).
		Build(r)
	if err != nil {
		return err
	}

	// Add index to AWSMachine to find by providerID
	if err := mgr.GetFieldIndexer().IndexField(ctx, &infrav1.AWSMachine{},
		InstanceIDIndex,
		r.indexAWSMachineByInstanceID,
	); err != nil {
		return errors.Wrap(err, "error setting index fields")
	}

	requeueAWSMachinesForUnpausedCluster := r.requeueAWSMachinesForUnpausedCluster(log)
	return controller.Watch(
		&source.Kind{Type: &clusterv1.Cluster{}},
		handler.EnqueueRequestsFromMapFunc(requeueAWSMachinesForUnpausedCluster),
		predicates.ClusterUnpausedAndInfrastructureReady(log),
	)
}

func (r *AWSMachineReconciler) reconcileDelete(machineScope *scope.MachineScope, clusterScope cloud.ClusterScoper, ec2Scope scope.EC2Scope, elbScope scope.ELBScope, objectStoreScope scope.S3Scope) (ctrl.Result, error) {
	machineScope.Info("Handling deleted AWSMachine")

	ec2Service := r.getEC2Service(ec2Scope)

	if err := r.deleteBootstrapData(machineScope, clusterScope, objectStoreScope); err != nil {
		machineScope.Error(err, "unable to delete machine")
		return ctrl.Result{}, err
	}

	instance, err := r.findInstance(machineScope, ec2Service)
	if err != nil && err != ec2.ErrInstanceNotFoundByID {
		machineScope.Error(err, "query to find instance failed")
		return ctrl.Result{}, err
	}

	if instance == nil {
		// The machine was never created or was deleted by some other entity
		// One way to reach this state:
		// 1. Scale deployment to 0
		// 2. Rename EC2 machine, and delete ProviderID from spec of both Machine
		// and AWSMachine
		// 3. Issue a delete
		// 4. Scale controller deployment to 1
		machineScope.V(2).Info("Unable to locate EC2 instance by ID or tags")
		r.Recorder.Eventf(machineScope.AWSMachine, corev1.EventTypeWarning, "NoInstanceFound", "Unable to find matching EC2 instance")
		controllerutil.RemoveFinalizer(machineScope.AWSMachine, infrav1.MachineFinalizer)
		return ctrl.Result{}, nil
	}

	machineScope.V(3).Info("EC2 instance found matching deleted AWSMachine", "instance-id", instance.ID)

	if err := r.reconcileLBAttachment(machineScope, elbScope, instance); err != nil {
		// We are tolerating AccessDenied error, so this won't block for users with older version of IAM;
		// all the other errors are blocking.
		if !elb.IsAccessDenied(err) && !elb.IsNotFound(err) {
			conditions.MarkFalse(machineScope.AWSMachine, infrav1.ELBAttachedCondition, "DeletingFailed", clusterv1.ConditionSeverityWarning, err.Error())
			return ctrl.Result{}, errors.Errorf("failed to reconcile LB attachment: %+v", err)
		}
	}

	if machineScope.IsControlPlane() {
		conditions.MarkFalse(machineScope.AWSMachine, infrav1.ELBAttachedCondition, clusterv1.DeletedReason, clusterv1.ConditionSeverityInfo, "")
	}

	if feature.Gates.Enabled(feature.EventBridgeInstanceState) {
		instancestateSvc := instancestate.NewService(ec2Scope)
		instancestateSvc.RemoveInstanceFromEventPattern(instance.ID)
	}

	// Check the instance state. If it's already shutting down or terminated,
	// do nothing. Otherwise attempt to delete it.
	// This decision is based on the ec2-instance-lifecycle graph at
	// https://docs.aws.amazon.com/AWSEC2/latest/UserGuide/ec2-instance-lifecycle.html
	switch instance.State {
	case infrav1.InstanceStateShuttingDown, infrav1.InstanceStateTerminated:
		machineScope.Info("EC2 instance is shutting down or already terminated", "instance-id", instance.ID)
	default:
		machineScope.Info("Terminating EC2 instance", "instance-id", instance.ID)

		// Set the InstanceReadyCondition and patch the object before the blocking operation
		conditions.MarkFalse(machineScope.AWSMachine, infrav1.InstanceReadyCondition, clusterv1.DeletingReason, clusterv1.ConditionSeverityInfo, "")
		if err := machineScope.PatchObject(); err != nil {
			machineScope.Error(err, "failed to patch object")
			return ctrl.Result{}, err
		}

		if err := ec2Service.TerminateInstanceAndWait(instance.ID); err != nil {
			machineScope.Error(err, "failed to terminate instance")
			conditions.MarkFalse(machineScope.AWSMachine, infrav1.InstanceReadyCondition, "DeletingFailed", clusterv1.ConditionSeverityWarning, err.Error())
			r.Recorder.Eventf(machineScope.AWSMachine, corev1.EventTypeWarning, "FailedTerminate", "Failed to terminate instance %q: %v", instance.ID, err)
			return ctrl.Result{}, err
		}
		conditions.MarkFalse(machineScope.AWSMachine, infrav1.InstanceReadyCondition, clusterv1.DeletedReason, clusterv1.ConditionSeverityInfo, "")

		// If the AWSMachine specifies NetworkStatus Interfaces, detach the cluster's core Security Groups from them as part of deletion.
		if len(machineScope.AWSMachine.Spec.NetworkInterfaces) > 0 {
			core, err := ec2Service.GetCoreSecurityGroups(machineScope)
			if err != nil {
				machineScope.Error(err, "failed to get core security groups to detach from instance's network interfaces")
				return ctrl.Result{}, err
			}

			machineScope.V(3).Info(
				"Detaching security groups from provided network interface",
				"groups", core,
				"instanceID", instance.ID,
			)

			conditions.MarkFalse(machineScope.AWSMachine, infrav1.SecurityGroupsReadyCondition, clusterv1.DeletingReason, clusterv1.ConditionSeverityInfo, "")
			if err := machineScope.PatchObject(); err != nil {
				return ctrl.Result{}, err
			}

			for _, id := range machineScope.AWSMachine.Spec.NetworkInterfaces {
				if err := ec2Service.DetachSecurityGroupsFromNetworkInterface(core, id); err != nil {
					machineScope.Error(err, "failed to detach security groups from instance's network interfaces")
					conditions.MarkFalse(machineScope.AWSMachine, infrav1.SecurityGroupsReadyCondition, "DeletingFailed", clusterv1.ConditionSeverityWarning, err.Error())
					return ctrl.Result{}, err
				}
			}
			conditions.MarkFalse(machineScope.AWSMachine, infrav1.SecurityGroupsReadyCondition, clusterv1.DeletedReason, clusterv1.ConditionSeverityInfo, "")
		}

		machineScope.Info("EC2 instance successfully terminated", "instance-id", instance.ID)
		r.Recorder.Eventf(machineScope.AWSMachine, corev1.EventTypeNormal, "SuccessfulTerminate", "Terminated instance %q", instance.ID)
	}

	// Instance is deleted so remove the finalizer.
	controllerutil.RemoveFinalizer(machineScope.AWSMachine, infrav1.MachineFinalizer)

	return ctrl.Result{}, nil
}

// findInstance queries the EC2 apis and retrieves the instance if it exists.
// If providerID is empty, finds instance by tags and if it cannot be found, returns empty instance with nil error.
// If providerID is set, either finds the instance by ID or returns error.
func (r *AWSMachineReconciler) findInstance(scope *scope.MachineScope, ec2svc services.EC2Interface) (*infrav1.Instance, error) {
	var instance *infrav1.Instance

	// Parse the ProviderID.
	pid, err := noderefutil.NewProviderID(scope.GetProviderID())
	if err != nil {
		if !errors.Is(err, noderefutil.ErrEmptyProviderID) {
			return nil, errors.Wrapf(err, "failed to parse Spec.ProviderID")
		}
		// If the ProviderID is empty, try to query the instance using tags.
		// If an instance cannot be found, GetRunningInstanceByTags returns empty instance with nil error.
		instance, err = ec2svc.GetRunningInstanceByTags(scope)
		if err != nil {
			return nil, errors.Wrapf(err, "failed to query AWSMachine instance by tags")
		}
	} else {
		// If the ProviderID is populated, describe the instance using the ID.
		// InstanceIfExists() returns error (ErrInstanceNotFoundByID or ErrDescribeInstance) if the instance could not be found.
		instance, err = ec2svc.InstanceIfExists(pointer.StringPtr(pid.ID()))
		if err != nil {
			return nil, err
		}
	}

	// The only case where the instance is nil here is when the providerId is empty and instance could not be found by tags.
	return instance, nil
}

func (r *AWSMachineReconciler) reconcileNormal(_ context.Context, machineScope *scope.MachineScope, clusterScope cloud.ClusterScoper, ec2Scope scope.EC2Scope, elbScope scope.ELBScope, objectStoreScope scope.S3Scope) (ctrl.Result, error) {
	machineScope.Info("Reconciling AWSMachine")

	// If the AWSMachine is in an error state, return early.
	if machineScope.HasFailed() {
		machineScope.Info("Error state detected, skipping reconciliation")

		// If we are in a failed state, delete the secret regardless of instance state.
		if err := r.deleteBootstrapData(machineScope, clusterScope, objectStoreScope); err != nil {
			machineScope.Error(err, "unable to reconcile machine")
			return ctrl.Result{}, err
		}

		return ctrl.Result{}, nil
	}

	if !machineScope.Cluster.Status.InfrastructureReady {
		machineScope.Info("Cluster infrastructure is not ready yet")
		conditions.MarkFalse(machineScope.AWSMachine, infrav1.InstanceReadyCondition, infrav1.WaitingForClusterInfrastructureReason, clusterv1.ConditionSeverityInfo, "")
		return ctrl.Result{}, nil
	}

	// Make sure bootstrap data is available and populated.
	if machineScope.Machine.Spec.Bootstrap.DataSecretName == nil {
		machineScope.Info("Bootstrap data secret reference is not yet available")
		conditions.MarkFalse(machineScope.AWSMachine, infrav1.InstanceReadyCondition, infrav1.WaitingForBootstrapDataReason, clusterv1.ConditionSeverityInfo, "")
		return ctrl.Result{}, nil
	}

	ec2svc := r.getEC2Service(ec2Scope)

	// Find existing instance
	instance, err := r.findInstance(machineScope, ec2svc)
	if err != nil {
		machineScope.Error(err, "unable to find instance")
		conditions.MarkUnknown(machineScope.AWSMachine, infrav1.InstanceReadyCondition, infrav1.InstanceNotFoundReason, err.Error())
		return ctrl.Result{}, err
	}

	// If the AWSMachine doesn't have our finalizer, add it.
	controllerutil.AddFinalizer(machineScope.AWSMachine, infrav1.MachineFinalizer)
	// Register the finalizer after first read operation from AWS to avoid orphaning AWS resources on delete
	if err := machineScope.PatchObject(); err != nil {
		machineScope.Error(err, "unable to patch object")
		return ctrl.Result{}, err
	}

	// Create new instance since providerId is nil and instance could not be found by tags.
	if instance == nil {
		// Avoid a flickering condition between InstanceProvisionStarted and InstanceProvisionFailed if there's a persistent failure with createInstance
		if conditions.GetReason(machineScope.AWSMachine, infrav1.InstanceReadyCondition) != infrav1.InstanceProvisionFailedReason {
			conditions.MarkFalse(machineScope.AWSMachine, infrav1.InstanceReadyCondition, infrav1.InstanceProvisionStartedReason, clusterv1.ConditionSeverityInfo, "")
			if patchErr := machineScope.PatchObject(); err != nil {
				machineScope.Error(patchErr, "failed to patch conditions")
				return ctrl.Result{}, patchErr
			}
		}

		var objectStoreSvc services.ObjectStoreInterface

		if objectStoreScope != nil {
			objectStoreSvc = r.getObjectStoreService(objectStoreScope)
		}

		instance, err = r.createInstance(ec2svc, machineScope, clusterScope, objectStoreSvc)
		if err != nil {
			machineScope.Error(err, "unable to create instance")
			conditions.MarkFalse(machineScope.AWSMachine, infrav1.InstanceReadyCondition, infrav1.InstanceProvisionFailedReason, clusterv1.ConditionSeverityError, err.Error())
			return ctrl.Result{}, err
		}
	}
	if feature.Gates.Enabled(feature.EventBridgeInstanceState) {
		instancestateSvc := instancestate.NewService(ec2Scope)
		if err := instancestateSvc.AddInstanceToEventPattern(instance.ID); err != nil {
			return ctrl.Result{}, errors.Wrap(err, "failed to add instance to Event Bridge instance state rule")
		}
	}

	// Make sure Spec.ProviderID and Spec.InstanceID are always set.
	machineScope.SetProviderID(instance.ID, instance.AvailabilityZone)
	machineScope.SetInstanceID(instance.ID)

	// See https://docs.aws.amazon.com/AWSEC2/latest/UserGuide/ec2-instance-lifecycle.html

	// Sets the AWSMachine status Interruptible, when the SpotMarketOptions is enabled for AWSMachine, Interruptible is set as true.
	machineScope.SetInterruptible()

	existingInstanceState := machineScope.GetInstanceState()
	machineScope.SetInstanceState(instance.State)

	// Proceed to reconcile the AWSMachine state.
	if existingInstanceState == nil || *existingInstanceState != instance.State {
		machineScope.Info("EC2 instance state changed", "state", instance.State, "instance-id", *machineScope.GetInstanceID())
	}

	switch instance.State {
	case infrav1.InstanceStatePending:
		machineScope.SetNotReady()
		conditions.MarkFalse(machineScope.AWSMachine, infrav1.InstanceReadyCondition, infrav1.InstanceNotReadyReason, clusterv1.ConditionSeverityWarning, "")
	case infrav1.InstanceStateStopping, infrav1.InstanceStateStopped:
		machineScope.SetNotReady()
		conditions.MarkFalse(machineScope.AWSMachine, infrav1.InstanceReadyCondition, infrav1.InstanceStoppedReason, clusterv1.ConditionSeverityError, "")
	case infrav1.InstanceStateRunning:
		machineScope.SetReady()
		conditions.MarkTrue(machineScope.AWSMachine, infrav1.InstanceReadyCondition)
	case infrav1.InstanceStateShuttingDown, infrav1.InstanceStateTerminated:
		machineScope.SetNotReady()
		machineScope.Info("Unexpected EC2 instance termination", "state", instance.State, "instance-id", *machineScope.GetInstanceID())
		r.Recorder.Eventf(machineScope.AWSMachine, corev1.EventTypeWarning, "InstanceUnexpectedTermination", "Unexpected EC2 instance termination")
		conditions.MarkFalse(machineScope.AWSMachine, infrav1.InstanceReadyCondition, infrav1.InstanceTerminatedReason, clusterv1.ConditionSeverityError, "")
	default:
		machineScope.SetNotReady()
		machineScope.Info("EC2 instance state is undefined", "state", instance.State, "instance-id", *machineScope.GetInstanceID())
		r.Recorder.Eventf(machineScope.AWSMachine, corev1.EventTypeWarning, "InstanceUnhandledState", "EC2 instance state is undefined")
		machineScope.SetFailureReason(capierrors.UpdateMachineError)
		machineScope.SetFailureMessage(errors.Errorf("EC2 instance state %q is undefined", instance.State))
		conditions.MarkUnknown(machineScope.AWSMachine, infrav1.InstanceReadyCondition, "", "")
	}

	// reconcile the deletion of the bootstrap data secret now that we have updated instance state
	if deleteSecretErr := r.deleteBootstrapData(machineScope, clusterScope, objectStoreScope); deleteSecretErr != nil {
		r.Log.Error(deleteSecretErr, "unable to delete secrets")
		return ctrl.Result{}, deleteSecretErr
	}

	if instance.State == infrav1.InstanceStateTerminated {
		machineScope.SetFailureReason(capierrors.UpdateMachineError)
		machineScope.SetFailureMessage(errors.Errorf("EC2 instance state %q is unexpected", instance.State))
	}

	// tasks that can take place during all known instance states
	if machineScope.InstanceIsInKnownState() {
		_, err = r.ensureTags(ec2svc, machineScope.AWSMachine, machineScope.GetInstanceID(), machineScope.AdditionalTags())
		if err != nil {
			machineScope.Error(err, "failed to ensure tags")
			return ctrl.Result{}, err
		}

		if instance != nil {
			r.ensureStorageTags(ec2svc, instance, machineScope.AWSMachine)
		}

		if err := r.reconcileLBAttachment(machineScope, elbScope, instance); err != nil {
			machineScope.Error(err, "failed to reconcile LB attachment")
			return ctrl.Result{}, err
		}
	}

	// tasks that can only take place during operational instance states
	if machineScope.InstanceIsOperational() {
		machineScope.SetAddresses(instance.Addresses)

		existingSecurityGroups, err := ec2svc.GetInstanceSecurityGroups(*machineScope.GetInstanceID())
		if err != nil {
			machineScope.Error(err, "unable to get instance security groups")
			return ctrl.Result{}, err
		}

		// Ensure that the security groups are correct.
		_, err = r.ensureSecurityGroups(ec2svc, machineScope, machineScope.AWSMachine.Spec.AdditionalSecurityGroups, existingSecurityGroups)
		if err != nil {
			conditions.MarkFalse(machineScope.AWSMachine, infrav1.SecurityGroupsReadyCondition, infrav1.SecurityGroupsFailedReason, clusterv1.ConditionSeverityError, err.Error())
			machineScope.Error(err, "unable to ensure security groups")
			return ctrl.Result{}, err
		}
		conditions.MarkTrue(machineScope.AWSMachine, infrav1.SecurityGroupsReadyCondition)
	}

	return ctrl.Result{}, nil
}

func (r *AWSMachineReconciler) deleteEncryptedBootstrapDataSecret(machineScope *scope.MachineScope, clusterScope cloud.ClusterScoper) error {
	secretSvc, secretBackendErr := r.getSecretService(machineScope, clusterScope)
	if secretBackendErr != nil {
		machineScope.Error(secretBackendErr, "unable to get secret service backend")
		return secretBackendErr
	}

	// do nothing if there isn't a secret
	if machineScope.GetSecretPrefix() == "" {
		return nil
	}
	if machineScope.GetSecretCount() == 0 {
		return errors.New("secretPrefix present, but secretCount is not set")
	}

	// Do nothing if the AWSMachine is not in a failed state, and is operational from an EC2 perspective, but does not have a node reference
	if !machineScope.HasFailed() && machineScope.InstanceIsOperational() && machineScope.Machine.Status.NodeRef == nil && !machineScope.AWSMachineIsDeleted() {
		return nil
	}
	machineScope.Info("Deleting unneeded entry from AWS Secret", "secretPrefix", machineScope.GetSecretPrefix())
	if err := secretSvc.Delete(machineScope); err != nil {
		machineScope.Info("Unable to delete entries from AWS Secret containing encrypted userdata", "secretPrefix", machineScope.GetSecretPrefix())
		r.Recorder.Eventf(machineScope.AWSMachine, corev1.EventTypeWarning, "FailedDeleteEncryptedBootstrapDataSecrets", "AWS Secret entries containing userdata not deleted")
		return err
	}
	r.Recorder.Eventf(machineScope.AWSMachine, corev1.EventTypeNormal, "SuccessfulDeleteEncryptedBootstrapDataSecrets", "AWS Secret entries containing userdata deleted")

	machineScope.DeleteSecretPrefix()
	machineScope.SetSecretCount(0)

	return nil
}

func (r *AWSMachineReconciler) createInstance(ec2svc services.EC2Interface, machineScope *scope.MachineScope, clusterScope cloud.ClusterScoper, objectStoreSvc services.ObjectStoreInterface) (*infrav1.Instance, error) {
	machineScope.Info("Creating EC2 instance")

	userData, userDataFormat, userDataErr := r.resolveUserData(machineScope, clusterScope, objectStoreSvc)
	if userDataErr != nil {
		return nil, errors.Wrapf(userDataErr, "failed to resolve userdata")
	}

	instance, err := ec2svc.CreateInstance(machineScope, userData, userDataFormat)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to create AWSMachine instance")
	}

	return instance, nil
}

func (r *AWSMachineReconciler) resolveUserData(machineScope *scope.MachineScope, clusterScope cloud.ClusterScoper, objectStoreSvc services.ObjectStoreInterface) ([]byte, string, error) {
	userData, userDataFormat, err := machineScope.GetRawBootstrapDataWithFormat()
	if err != nil {
		r.Recorder.Eventf(machineScope.AWSMachine, corev1.EventTypeWarning, "FailedGetBootstrapData", err.Error())
		return nil, "", err
	}

	if machineScope.UseSecretsManager(userDataFormat) {
		userData, err = r.cloudInitUserData(machineScope, clusterScope, userData)
	}

	if machineScope.UseIgnition(userDataFormat) {
		userData, err = r.ignitionUserData(machineScope, objectStoreSvc, userData)
	}

	return userData, userDataFormat, err
}

func (r *AWSMachineReconciler) cloudInitUserData(machineScope *scope.MachineScope, clusterScope cloud.ClusterScoper, userData []byte) ([]byte, error) {
	secretSvc, secretBackendErr := r.getSecretService(machineScope, clusterScope)
	if secretBackendErr != nil {
		machineScope.Error(secretBackendErr, "unable to reconcile machine")
		return nil, secretBackendErr
	}

	compressedUserData, compressErr := userdata.GzipBytes(userData)
	if compressErr != nil {
		return nil, compressErr
	}
	prefix, chunks, serviceErr := secretSvc.Create(machineScope, compressedUserData)
	// Only persist the AWS Secret Backend entries if there is at least one
	if chunks > 0 {
		machineScope.SetSecretPrefix(prefix)
		machineScope.SetSecretCount(chunks)
	}
	// Register the Secret ARN immediately to avoid orphaning whatever AWS resources have been created
	if err := machineScope.PatchObject(); err != nil {
		return nil, err
	}
	if serviceErr != nil {
		r.Recorder.Eventf(machineScope.AWSMachine, corev1.EventTypeWarning, "FailedCreateAWSSecrets", serviceErr.Error())
		machineScope.Error(serviceErr, "Failed to create AWS Secret entry", "secretPrefix", prefix)
		return nil, serviceErr
	}
	encryptedCloudInit, err := secretSvc.UserData(machineScope.GetSecretPrefix(), machineScope.GetSecretCount(), machineScope.InfraCluster.Region(), r.Endpoints)
	if err != nil {
		r.Recorder.Eventf(machineScope.AWSMachine, corev1.EventTypeWarning, "FailedGenerateAWSSecretsCloudInit", err.Error())
		return nil, err
	}
	return encryptedCloudInit, nil
}

func (r *AWSMachineReconciler) ignitionUserData(scope *scope.MachineScope, objectStoreSvc services.ObjectStoreInterface, userData []byte) ([]byte, error) {
	if objectStoreSvc == nil {
		return nil, errors.New("object store service not available")
	}

	objectURL, err := objectStoreSvc.Create(scope, userData)
	if err != nil {
		return nil, errors.Wrap(err, "creating userdata object")
	}

	ignData := &ignTypes.Config{
		Ignition: ignTypes.Ignition{
			Version: "2.3.0",
			Config: ignTypes.IgnitionConfig{
				Append: []ignTypes.ConfigReference{
					{
						Source: objectURL,
					},
				},
			},
		},
	}

	ignitionUserData, err := json.Marshal(ignData)
	if err != nil {
		r.Recorder.Eventf(scope.AWSMachine, corev1.EventTypeWarning, "FailedGenerateIgnition", err.Error())
		return nil, errors.Wrap(err, "serializing generated data")
	}

	return ignitionUserData, nil
}

func (r *AWSMachineReconciler) deleteBootstrapData(machineScope *scope.MachineScope, clusterScope cloud.ClusterScoper, objectStoreScope scope.S3Scope) error {
	if !machineScope.AWSMachine.Spec.CloudInit.InsecureSkipSecretsManager {
		if err := r.deleteEncryptedBootstrapDataSecret(machineScope, clusterScope); err != nil {
			return err
		}
	}

	if objectStoreScope != nil {
		// Bootstrap data will be removed from S3 if it is already populated.
		if err := r.deleteIgnitionBootstrapDataFromS3(machineScope, r.getObjectStoreService(objectStoreScope)); err != nil {
			return err
		}
	}

	return nil
}

func (r *AWSMachineReconciler) deleteIgnitionBootstrapDataFromS3(machineScope *scope.MachineScope, objectStoreSvc services.ObjectStoreInterface) error {
	// Do nothing if the AWSMachine is not in a failed state, and is operational from an EC2 perspective, but does not have a node reference
	if !machineScope.HasFailed() && machineScope.InstanceIsOperational() && machineScope.Machine.Status.NodeRef == nil && !machineScope.AWSMachineIsDeleted() {
		return nil
	}

	// If bootstrap data has not been populated yet, we cannot determine it's format, so there is probably nothing to do.
	if machineScope.Machine.Spec.Bootstrap.DataSecretName == nil {
		return nil
	}

	machineScope.Info("Deleting unneeded entry from AWS S3", "secretPrefix", machineScope.GetSecretPrefix())

	_, userDataFormat, err := machineScope.GetRawBootstrapDataWithFormat()
	if err != nil {
		r.Recorder.Eventf(machineScope.AWSMachine, corev1.EventTypeWarning, "FailedGetBootstrapData", err.Error())
		return err
	}

	if !machineScope.UseIgnition(userDataFormat) {
		return nil
	}

	if err := objectStoreSvc.Delete(machineScope); err != nil {
		return errors.Wrap(err, "deleting bootstrap data object")
	}

	return nil
}

func (r *AWSMachineReconciler) reconcileLBAttachment(machineScope *scope.MachineScope, elbScope scope.ELBScope, i *infrav1.Instance) error {
	if !machineScope.IsControlPlane() {
		return nil
	}

	elbsvc := r.getELBService(elbScope)

	// In order to prevent sending request to a "not-ready" control plane machines, it is required to remove the machine
	// from the ELB as soon as the machine gets deleted or when the machine is in a not running state.
	if !machineScope.AWSMachine.DeletionTimestamp.IsZero() || !machineScope.InstanceIsRunning() {
		registered, err := elbsvc.IsInstanceRegisteredWithAPIServerELB(i)
		if err != nil {
			r.Recorder.Eventf(machineScope.AWSMachine, corev1.EventTypeWarning, "FailedDetachControlPlaneELB",
				"Failed to deregister control plane instance %q from load balancer: failed to determine registration status: %v", i.ID, err)
			return errors.Wrapf(err, "could not deregister control plane instance %q from load balancer - error determining registration status", i.ID)
		}
		if !registered {
			// Already deregistered - nothing more to do
			return nil
		}

		if err := elbsvc.DeregisterInstanceFromAPIServerELB(i); err != nil {
			r.Recorder.Eventf(machineScope.AWSMachine, corev1.EventTypeWarning, "FailedDetachControlPlaneELB",
				"Failed to deregister control plane instance %q from load balancer: %v", i.ID, err)
			conditions.MarkFalse(machineScope.AWSMachine, infrav1.ELBAttachedCondition, infrav1.ELBDetachFailedReason, clusterv1.ConditionSeverityError, err.Error())
			return errors.Wrapf(err, "could not deregister control plane instance %q from load balancer", i.ID)
		}
		r.Recorder.Eventf(machineScope.AWSMachine, corev1.EventTypeNormal, "SuccessfulDetachControlPlaneELB",
			"Control plane instance %q is de-registered from load balancer", i.ID)
		return nil
	}

	registered, err := elbsvc.IsInstanceRegisteredWithAPIServerELB(i)
	if err != nil {
		r.Recorder.Eventf(machineScope.AWSMachine, corev1.EventTypeWarning, "FailedAttachControlPlaneELB",
			"Failed to register control plane instance %q with load balancer: failed to determine registration status: %v", i.ID, err)
		return errors.Wrapf(err, "could not register control plane instance %q with load balancer - error determining registration status", i.ID)
	}
	if registered {
		// Already registered - nothing more to do
		return nil
	}

	if err := elbsvc.RegisterInstanceWithAPIServerELB(i); err != nil {
		r.Recorder.Eventf(machineScope.AWSMachine, corev1.EventTypeWarning, "FailedAttachControlPlaneELB",
			"Failed to register control plane instance %q with load balancer: %v", i.ID, err)
		conditions.MarkFalse(machineScope.AWSMachine, infrav1.ELBAttachedCondition, infrav1.ELBAttachFailedReason, clusterv1.ConditionSeverityError, err.Error())
		return errors.Wrapf(err, "could not register control plane instance %q with load balancer", i.ID)
	}
	r.Recorder.Eventf(machineScope.AWSMachine, corev1.EventTypeNormal, "SuccessfulAttachControlPlaneELB",
		"Control plane instance %q is registered with load balancer", i.ID)
	conditions.MarkTrue(machineScope.AWSMachine, infrav1.ELBAttachedCondition)
	return nil
}

// AWSClusterToAWSMachines is a handler.ToRequestsFunc to be used to enqeue requests for reconciliation
// of AWSMachines.
func (r *AWSMachineReconciler) AWSClusterToAWSMachines(log logr.Logger) handler.MapFunc {
	return func(o client.Object) []ctrl.Request {
		c, ok := o.(*infrav1.AWSCluster)
		if !ok {
			panic(fmt.Sprintf("Expected a AWSCluster but got a %T", o))
		}

		log := log.WithValues("objectMapper", "awsClusterToAWSMachine", "namespace", c.Namespace, "awsCluster", c.Name)

		// Don't handle deleted AWSClusters
		if !c.ObjectMeta.DeletionTimestamp.IsZero() {
			log.V(4).Info("AWSCluster has a deletion timestamp, skipping mapping.")
			return nil
		}

		cluster, err := util.GetOwnerCluster(context.TODO(), r.Client, c.ObjectMeta)
		switch {
		case apierrors.IsNotFound(err) || cluster == nil:
			log.V(4).Info("Cluster for AWSCluster not found, skipping mapping.")
			return nil
		case err != nil:
			log.Error(err, "Failed to get owning cluster, skipping mapping.")
			return nil
		}

		return r.requestsForCluster(log, cluster.Namespace, cluster.Name)
	}
}

func (r *AWSMachineReconciler) requeueAWSMachinesForUnpausedCluster(log logr.Logger) handler.MapFunc {
	return func(o client.Object) []ctrl.Request {
		c, ok := o.(*clusterv1.Cluster)
		if !ok {
			panic(fmt.Sprintf("Expected a Cluster but got a %T", o))
		}

		log := log.WithValues("objectMapper", "clusterToAWSMachine", "namespace", c.Namespace, "cluster", c.Name)

		// Don't handle deleted clusters
		if !c.ObjectMeta.DeletionTimestamp.IsZero() {
			log.V(4).Info("Cluster has a deletion timestamp, skipping mapping.")
			return nil
		}

		return r.requestsForCluster(log, c.Namespace, c.Name)
	}
}

func (r *AWSMachineReconciler) requestsForCluster(log logr.Logger, namespace, name string) []ctrl.Request {
	labels := map[string]string{clusterv1.ClusterLabelName: name}
	machineList := &clusterv1.MachineList{}
	if err := r.Client.List(context.TODO(), machineList, client.InNamespace(namespace), client.MatchingLabels(labels)); err != nil {
		log.Error(err, "Failed to get owned Machines, skipping mapping.")
		return nil
	}

	result := make([]ctrl.Request, 0, len(machineList.Items))
	for _, m := range machineList.Items {
		log.WithValues("machine", m.Name)
		if m.Spec.InfrastructureRef.GroupVersionKind().Kind != "AWSMachine" {
			log.V(4).Info("Machine has an InfrastructureRef for a different type, will not add to reconciliation request.")
			continue
		}
		if m.Spec.InfrastructureRef.Name == "" {
			log.V(4).Info("Machine has an InfrastructureRef with an empty name, will not add to reconciliation request.")
			continue
		}
		log.WithValues("awsMachine", m.Spec.InfrastructureRef.Name)
		log.V(4).Info("Adding AWSMachine to reconciliation request.")
		result = append(result, ctrl.Request{NamespacedName: client.ObjectKey{Namespace: m.Namespace, Name: m.Spec.InfrastructureRef.Name}})
	}
	return result
}

func (r *AWSMachineReconciler) getInfraCluster(ctx context.Context, log logr.Logger, cluster *clusterv1.Cluster, awsMachine *infrav1.AWSMachine) (scope.EC2Scope, error) {
	var clusterScope *scope.ClusterScope
	var managedControlPlaneScope *scope.ManagedControlPlaneScope
	var err error

	if cluster.Spec.ControlPlaneRef != nil && cluster.Spec.ControlPlaneRef.Kind == "AWSManagedControlPlane" {
		controlPlane := &ekscontrolplanev1.AWSManagedControlPlane{}
		controlPlaneName := client.ObjectKey{
			Namespace: awsMachine.Namespace,
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
			Endpoints:      r.Endpoints,
		})
		if err != nil {
			return nil, err
		}

		return managedControlPlaneScope, nil
	}

	awsCluster := &infrav1.AWSCluster{}

	infraClusterName := client.ObjectKey{
		Namespace: awsMachine.Namespace,
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

func (r *AWSMachineReconciler) indexAWSMachineByInstanceID(o client.Object) []string {
	awsMachine, ok := o.(*infrav1.AWSMachine)
	if !ok {
		r.Log.Error(errors.New("incorrect type"), "expected an AWSMachine", "type", fmt.Sprintf("%T", o))
		return nil
	}

	if awsMachine.Spec.InstanceID != nil {
		return []string{*awsMachine.Spec.InstanceID}
	}

	return nil
}

func (r *AWSMachineReconciler) ensureStorageTags(ec2svc services.EC2Interface, instance *infrav1.Instance, machine *infrav1.AWSMachine) {
	annotations, err := r.machineAnnotationJSON(machine, VolumeTagsLastAppliedAnnotation)
	if err != nil {
		r.Log.Error(err, "Failed to fetch the annotations for volume tags")
	}
	for _, volumeID := range instance.VolumeIDs {
		if subAnnotation, ok := annotations[volumeID].(map[string]interface{}); ok {
			newAnnotation, err := r.ensureVolumeTags(ec2svc, aws.String(volumeID), subAnnotation, machine.Spec.AdditionalTags)
			if err != nil {
				r.Log.Error(err, "Failed to fetch the changed volume tags in EC2 instance")
			}
			annotations[volumeID] = newAnnotation
		} else {
			newAnnotation, err := r.ensureVolumeTags(ec2svc, aws.String(volumeID), make(map[string]interface{}), machine.Spec.AdditionalTags)
			if err != nil {
				r.Log.Error(err, "Failed to fetch the changed volume tags in EC2 instance")
			}
			annotations[volumeID] = newAnnotation
		}

		// We also need to update the annotation if anything changed.
		err = r.updateMachineAnnotationJSON(machine, VolumeTagsLastAppliedAnnotation, annotations)
		if err != nil {
			r.Log.Error(err, "Failed to fetch the changed volume tags in EC2 instance")
		}
	}
}
