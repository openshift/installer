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
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/blang/semver"
	ignTypes "github.com/coreos/ignition/config/v2_3/types"
	ignV3Types "github.com/coreos/ignition/v2/config/v3_4/types"
	"github.com/go-logr/logr"
	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"github.com/pkg/errors"
	corev1 "k8s.io/api/core/v1"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	kerrors "k8s.io/apimachinery/pkg/util/errors"
	"k8s.io/client-go/tools/record"
	"k8s.io/klog/v2"
	"k8s.io/utils/ptr"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
	"sigs.k8s.io/controller-runtime/pkg/event"
	"sigs.k8s.io/controller-runtime/pkg/handler"
	"sigs.k8s.io/controller-runtime/pkg/predicate"
	"sigs.k8s.io/controller-runtime/pkg/source"

	infrav1 "sigs.k8s.io/cluster-api-provider-aws/v2/api/v1beta2"
	ekscontrolplanev1 "sigs.k8s.io/cluster-api-provider-aws/v2/controlplane/eks/api/v1beta2"
	"sigs.k8s.io/cluster-api-provider-aws/v2/feature"
	"sigs.k8s.io/cluster-api-provider-aws/v2/pkg/cloud"
	"sigs.k8s.io/cluster-api-provider-aws/v2/pkg/cloud/scope"
	"sigs.k8s.io/cluster-api-provider-aws/v2/pkg/cloud/services"
	"sigs.k8s.io/cluster-api-provider-aws/v2/pkg/cloud/services/ec2"
	"sigs.k8s.io/cluster-api-provider-aws/v2/pkg/cloud/services/elb"
	"sigs.k8s.io/cluster-api-provider-aws/v2/pkg/cloud/services/instancestate"
	"sigs.k8s.io/cluster-api-provider-aws/v2/pkg/cloud/services/s3"
	"sigs.k8s.io/cluster-api-provider-aws/v2/pkg/cloud/services/secretsmanager"
	"sigs.k8s.io/cluster-api-provider-aws/v2/pkg/cloud/services/ssm"
	"sigs.k8s.io/cluster-api-provider-aws/v2/pkg/cloud/services/userdata"
	"sigs.k8s.io/cluster-api-provider-aws/v2/pkg/logger"
	clusterv1 "sigs.k8s.io/cluster-api/api/v1beta1"
	"sigs.k8s.io/cluster-api/util"
	"sigs.k8s.io/cluster-api/util/annotations"
	"sigs.k8s.io/cluster-api/util/conditions"
	"sigs.k8s.io/cluster-api/util/predicates"
)

const (
	// InstanceIDIndex defines the aws machine controller's instance ID index.
	InstanceIDIndex = ".spec.instanceID"

	// DefaultReconcilerRequeue is the default value for the reconcile retry.
	DefaultReconcilerRequeue = 30 * time.Second
)

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
	TagUnmanagedNetworkResources bool
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

// +kubebuilder:rbac:groups=controlplane.cluster.x-k8s.io,resources=*,verbs=get;list;watch
// +kubebuilder:rbac:groups=infrastructure.cluster.x-k8s.io,resources=awsmachines,verbs=create;get;list;watch;update;patch;delete
// +kubebuilder:rbac:groups=infrastructure.cluster.x-k8s.io,resources=awsmachines/status,verbs=get;update;patch
// +kubebuilder:rbac:groups=cluster.x-k8s.io,resources=machines,verbs=get;list;watch;delete
// +kubebuilder:rbac:groups=cluster.x-k8s.io,resources=machines/status,verbs=get;list;watch
// +kubebuilder:rbac:groups="",resources=secrets;,verbs=get;list;watch
// +kubebuilder:rbac:groups="",resources=namespaces,verbs=get;list;watch
// +kubebuilder:rbac:groups="",resources=events,verbs=get;list;watch;create;update;patch

func (r *AWSMachineReconciler) Reconcile(ctx context.Context, req ctrl.Request) (_ ctrl.Result, reterr error) {
	log := logger.FromContext(ctx)

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

	log = log.WithValues("machine", klog.KObj(machine))

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

	log = log.WithValues("cluster", klog.KObj(cluster))

	infraCluster, err := r.getInfraCluster(ctx, log, cluster, awsMachine)
	if err != nil {
		return ctrl.Result{}, errors.Errorf("error getting infra provider cluster or control plane object: %v", err)
	}
	if infraCluster == nil {
		log.Info("AWSCluster or AWSManagedControlPlane is not ready yet")
		return ctrl.Result{}, nil
	}

	infrav1.SetDefaults_AWSMachineSpec(&awsMachine.Spec)

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
	log := logger.FromContext(ctx)
	AWSClusterToAWSMachines := r.AWSClusterToAWSMachines(log)

	controller, err := ctrl.NewControllerManagedBy(mgr).
		WithOptions(options).
		For(&infrav1.AWSMachine{}).
		Watches(
			&clusterv1.Machine{},
			handler.EnqueueRequestsFromMapFunc(util.MachineToInfrastructureMapFunc(infrav1.GroupVersion.WithKind("AWSMachine"))),
		).
		Watches(
			&infrav1.AWSCluster{},
			handler.EnqueueRequestsFromMapFunc(AWSClusterToAWSMachines),
		).
		WithEventFilter(predicates.ResourceNotPausedAndHasFilterLabel(mgr.GetScheme(), log.GetLogger(), r.WatchFilterValue)).
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
		source.Kind[client.Object](mgr.GetCache(), &clusterv1.Cluster{},
			handler.EnqueueRequestsFromMapFunc(requeueAWSMachinesForUnpausedCluster),
			predicates.ClusterPausedTransitionsOrInfrastructureReady(mgr.GetScheme(), log.GetLogger())),
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
		machineScope.Warn("Unable to locate EC2 instance by ID or tags")
		r.Recorder.Eventf(machineScope.AWSMachine, corev1.EventTypeWarning, "NoInstanceFound", "Unable to find matching EC2 instance")
		controllerutil.RemoveFinalizer(machineScope.AWSMachine, infrav1.MachineFinalizer)
		return ctrl.Result{}, nil
	}

	machineScope.Debug("EC2 instance found matching deleted AWSMachine", "instance-id", instance.ID)

	if err := r.reconcileLBAttachment(machineScope, elbScope, instance); err != nil {
		// We are tolerating AccessDenied error, so this won't block for users with older version of IAM;
		// all the other errors are blocking.
		// Because we are reconciling all load balancers, attempt to treat the error as a list of errors.
		if err = kerrors.FilterOut(err, elb.IsAccessDenied, elb.IsNotFound); err != nil {
			conditions.MarkFalse(machineScope.AWSMachine, infrav1.ELBAttachedCondition, "DeletingFailed", clusterv1.ConditionSeverityWarning, "%s", err.Error())
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
	case infrav1.InstanceStateShuttingDown:
		machineScope.Info("EC2 instance is shutting down or already terminated", "instance-id", instance.ID)
		// requeue reconciliation until we observe termination (or the instance can no longer be looked up)
		return ctrl.Result{RequeueAfter: time.Minute}, nil
	case infrav1.InstanceStateTerminated:
		machineScope.Info("EC2 instance terminated successfully", "instance-id", instance.ID)
		controllerutil.RemoveFinalizer(machineScope.AWSMachine, infrav1.MachineFinalizer)
		return ctrl.Result{}, nil
	default:
		machineScope.Info("Terminating EC2 instance", "instance-id", instance.ID)

		// Set the InstanceReadyCondition and patch the object before the blocking operation
		conditions.MarkFalse(machineScope.AWSMachine, infrav1.InstanceReadyCondition, clusterv1.DeletingReason, clusterv1.ConditionSeverityInfo, "")
		if err := machineScope.PatchObject(); err != nil {
			machineScope.Error(err, "failed to patch object")
			return ctrl.Result{}, err
		}

		if err := ec2Service.TerminateInstance(instance.ID); err != nil {
			machineScope.Error(err, "failed to terminate instance")
			conditions.MarkFalse(machineScope.AWSMachine, infrav1.InstanceReadyCondition, "DeletingFailed", clusterv1.ConditionSeverityWarning, "%s", err.Error())
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

			machineScope.Debug(
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
					conditions.MarkFalse(machineScope.AWSMachine, infrav1.SecurityGroupsReadyCondition, "DeletingFailed", clusterv1.ConditionSeverityWarning, "%s", err.Error())
					return ctrl.Result{}, err
				}
			}
			conditions.MarkFalse(machineScope.AWSMachine, infrav1.SecurityGroupsReadyCondition, clusterv1.DeletedReason, clusterv1.ConditionSeverityInfo, "")
		}

		// Release an Elastic IP when the machine has public IP Address (EIP) with a cluster-wide config
		// to consume from BYO IPv4 Pool.
		if machineScope.GetElasticIPPool() != nil {
			if err := ec2Service.ReleaseElasticIP(instance.ID); err != nil {
				machineScope.Error(err, "failed to release elastic IP address")
				return ctrl.Result{}, err
			}
		}

		machineScope.Info("EC2 instance successfully terminated", "instance-id", instance.ID)
		r.Recorder.Eventf(machineScope.AWSMachine, corev1.EventTypeNormal, "SuccessfulTerminate", "Terminated instance %q", instance.ID)

		// requeue reconciliation until we observe termination (or the instance can no longer be looked up)
		return ctrl.Result{RequeueAfter: time.Minute}, nil
	}
}

// findInstance queries the EC2 apis and retrieves the instance if it exists.
// If providerID is empty, finds instance by tags and if it cannot be found, returns empty instance with nil error.
// If providerID is set, either finds the instance by ID or returns error.
func (r *AWSMachineReconciler) findInstance(machineScope *scope.MachineScope, ec2svc services.EC2Interface) (*infrav1.Instance, error) {
	var instance *infrav1.Instance

	// Parse the ProviderID.
	pid, err := scope.NewProviderID(machineScope.GetProviderID())
	if err != nil {
		//nolint:staticcheck
		if !errors.Is(err, scope.ErrEmptyProviderID) {
			return nil, errors.Wrapf(err, "failed to parse Spec.ProviderID")
		}
		// If the ProviderID is empty, try to query the instance using tags.
		// If an instance cannot be found, GetRunningInstanceByTags returns empty instance with nil error.
		instance, err = ec2svc.GetRunningInstanceByTags(machineScope)
		if err != nil {
			return nil, errors.Wrapf(err, "failed to query AWSMachine instance by tags")
		}
	} else {
		// If the ProviderID is populated, describe the instance using the ID.
		// InstanceIfExists() returns error (ErrInstanceNotFoundByID or ErrDescribeInstance) if the instance could not be found.
		//nolint:staticcheck
		instance, err = ec2svc.InstanceIfExists(ptr.To[string](pid.ID()))
		if err != nil {
			return nil, err
		}
	}

	// The only case where the instance is nil here is when the providerId is empty and instance could not be found by tags.
	return instance, nil
}

//nolint:gocyclo
func (r *AWSMachineReconciler) reconcileNormal(_ context.Context, machineScope *scope.MachineScope, clusterScope cloud.ClusterScoper, ec2Scope scope.EC2Scope, elbScope scope.ELBScope, objectStoreScope scope.S3Scope) (ctrl.Result, error) {
	machineScope.Trace("Reconciling AWSMachine")

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
	if !machineScope.IsMachinePoolMachine() && machineScope.Machine.Spec.Bootstrap.DataSecretName == nil {
		machineScope.Info("Bootstrap data secret reference is not yet available")
		conditions.MarkFalse(machineScope.AWSMachine, infrav1.InstanceReadyCondition, infrav1.WaitingForBootstrapDataReason, clusterv1.ConditionSeverityInfo, "")
		return ctrl.Result{}, nil
	}

	ec2svc := r.getEC2Service(ec2Scope)

	// Find existing instance
	instance, err := r.findInstance(machineScope, ec2svc)
	if err != nil {
		machineScope.Error(err, "unable to find instance")
		conditions.MarkUnknown(machineScope.AWSMachine, infrav1.InstanceReadyCondition, infrav1.InstanceNotFoundReason, "%s", err.Error())
		return ctrl.Result{}, err
	}
	if instance == nil && machineScope.IsMachinePoolMachine() {
		err = errors.New("no instance found for machine pool")
		machineScope.Error(err, "unable to find instance")
		conditions.MarkUnknown(machineScope.AWSMachine, infrav1.InstanceReadyCondition, infrav1.InstanceNotFoundReason, "%s", err.Error())
		return ctrl.Result{}, err
	}

	// If the AWSMachine doesn't have our finalizer, add it.
	if controllerutil.AddFinalizer(machineScope.AWSMachine, infrav1.MachineFinalizer) {
		// Register the finalizer after first read operation from AWS to avoid orphaning AWS resources on delete
		if err := machineScope.PatchObject(); err != nil {
			machineScope.Error(err, "unable to patch object")
			return ctrl.Result{}, err
		}
	}

	// Create new instance since providerId is nil and instance could not be found by tags.
	if instance == nil {
		// Avoid a flickering condition between InstanceProvisionStarted and InstanceProvisionFailed if there's a persistent failure with createInstance
		if conditions.GetReason(machineScope.AWSMachine, infrav1.InstanceReadyCondition) != infrav1.InstanceProvisionFailedReason {
			conditions.MarkFalse(machineScope.AWSMachine, infrav1.InstanceReadyCondition, infrav1.InstanceProvisionStartedReason, clusterv1.ConditionSeverityInfo, "")
			if patchErr := machineScope.PatchObject(); patchErr != nil {
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
			conditions.MarkFalse(machineScope.AWSMachine, infrav1.InstanceReadyCondition, infrav1.InstanceProvisionFailedReason, clusterv1.ConditionSeverityError, "%s", err.Error())
			return ctrl.Result{}, err
		}
	}

	// BYO Public IPv4 Pool feature: allocates and associates an EIP to machine when PublicIP and
	// cluster-wide config Public IPv4 Pool configuration are set. The custom EIP is associated
	// after the instance is created and transictioned to Running state.
	// The CreateInstance() is enforcing to not assign public IP address when PublicIP is set with
	// BYOIpv4 Pool, preventing a duplicated EIP creation.
	if pool := machineScope.GetElasticIPPool(); pool != nil {
		requeue, err := ec2svc.ReconcileElasticIPFromPublicPool(pool, instance)
		if err != nil {
			machineScope.Error(err, "Failed to reconcile BYO Public IPv4")
			return ctrl.Result{}, err
		}
		if requeue {
			machineScope.Debug("Found instance in pending state while reconciling publicIpv4Pool, requeue", "instance", instance.ID)
			return ctrl.Result{RequeueAfter: DefaultReconcilerRequeue}, nil
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

	shouldRequeue := false
	switch instance.State {
	case infrav1.InstanceStatePending:
		machineScope.SetNotReady()
		shouldRequeue = true
		conditions.MarkFalse(machineScope.AWSMachine, infrav1.InstanceReadyCondition, infrav1.InstanceNotReadyReason, clusterv1.ConditionSeverityWarning, "")
	case infrav1.InstanceStateStopping, infrav1.InstanceStateStopped:
		machineScope.SetNotReady()
		conditions.MarkFalse(machineScope.AWSMachine, infrav1.InstanceReadyCondition, infrav1.InstanceStoppedReason, clusterv1.ConditionSeverityError, "")
	case infrav1.InstanceStateRunning:
		machineScope.SetReady()
		conditions.MarkTrue(machineScope.AWSMachine, infrav1.InstanceReadyCondition)
	case infrav1.InstanceStateShuttingDown, infrav1.InstanceStateTerminated:
		machineScope.SetNotReady()

		if machineScope.IsMachinePoolMachine() {
			// In an auto-scaling group, instance termination is perfectly normal on scale-down
			// and therefore should not be reported as error.
			machineScope.Info("EC2 instance of machine pool was terminated", "state", instance.State, "instance-id", *machineScope.GetInstanceID())
			r.Recorder.Eventf(machineScope.AWSMachine, corev1.EventTypeNormal, infrav1.InstanceTerminatedReason, "EC2 instance termination")
			conditions.MarkFalse(machineScope.AWSMachine, infrav1.InstanceReadyCondition, infrav1.InstanceTerminatedReason, clusterv1.ConditionSeverityInfo, "")
		} else {
			machineScope.Info("Unexpected EC2 instance termination", "state", instance.State, "instance-id", *machineScope.GetInstanceID())
			r.Recorder.Eventf(machineScope.AWSMachine, corev1.EventTypeWarning, "InstanceUnexpectedTermination", "Unexpected EC2 instance termination")
			conditions.MarkFalse(machineScope.AWSMachine, infrav1.InstanceReadyCondition, infrav1.InstanceTerminatedReason, clusterv1.ConditionSeverityError, "")
		}
	default:
		machineScope.SetNotReady()
		machineScope.Info("EC2 instance state is undefined", "state", instance.State, "instance-id", *machineScope.GetInstanceID())
		r.Recorder.Eventf(machineScope.AWSMachine, corev1.EventTypeWarning, "InstanceUnhandledState", "EC2 instance state is undefined")
		machineScope.SetFailureReason("UpdateError")
		machineScope.SetFailureMessage(errors.Errorf("EC2 instance state %q is undefined", instance.State))
		conditions.MarkUnknown(machineScope.AWSMachine, infrav1.InstanceReadyCondition, "", "")
	}

	// reconcile the deletion of the bootstrap data secret now that we have updated instance state
	if !machineScope.IsMachinePoolMachine() {
		if deleteSecretErr := r.deleteBootstrapData(machineScope, clusterScope, objectStoreScope); deleteSecretErr != nil {
			r.Log.Error(deleteSecretErr, "unable to delete secrets")
			return ctrl.Result{}, deleteSecretErr
		}

		// For machine pool machines, it is expected that the ASG terminates instances at any time,
		// so no error is logged for those.
		if instance.State == infrav1.InstanceStateTerminated {
			machineScope.SetFailureReason("UpdateError")
			machineScope.SetFailureMessage(errors.Errorf("EC2 instance state %q is unexpected", instance.State))
		}
	}

	// tasks that can take place during all known instance states
	if machineScope.InstanceIsInKnownState() {
		_, err = r.ensureTags(ec2svc, machineScope.AWSMachine, machineScope.GetInstanceID(), machineScope.AdditionalTags())
		if err != nil {
			machineScope.Error(err, "failed to ensure tags")
			return ctrl.Result{}, err
		}

		if instance != nil {
			r.ensureStorageTags(ec2svc, instance, machineScope.AWSMachine, machineScope.AdditionalTags())
		}

		if err := r.reconcileLBAttachment(machineScope, elbScope, instance); err != nil {
			// We are tolerating InstanceNotRunning error, so we don't report it as an error condition.
			// Because we are reconciling all load balancers, attempt to treat the error as a list of errors.
			if err := kerrors.FilterOut(err, elb.IsInstanceNotRunning); err != nil {
				machineScope.Error(err, "failed to reconcile LB attachment")
				return ctrl.Result{}, err
			}
			// Cannot attach non-running instances to LB
			shouldRequeue = true
		}
	}

	// tasks that can only take place during operational instance states
	if machineScope.InstanceIsOperational() {
		err := r.reconcileOperationalState(ec2svc, machineScope, instance)
		if err != nil {
			return ctrl.Result{}, err
		}
	}

	machineScope.Debug("done reconciling instance", "instance", instance)
	if shouldRequeue {
		machineScope.Debug("but find the instance is pending, requeue", "instance", instance.ID)
		return ctrl.Result{RequeueAfter: DefaultReconcilerRequeue}, nil
	}
	return ctrl.Result{}, nil
}

func (r *AWSMachineReconciler) reconcileOperationalState(ec2svc services.EC2Interface, machineScope *scope.MachineScope, instance *infrav1.Instance) error {
	machineScope.SetAddresses(instance.Addresses)

	existingSecurityGroups, err := ec2svc.GetInstanceSecurityGroups(*machineScope.GetInstanceID())
	if err != nil {
		machineScope.Error(err, "unable to get instance security groups")
		return err
	}

	// Ensure that the security groups are correct.
	_, err = r.ensureSecurityGroups(ec2svc, machineScope, machineScope.AWSMachine.Spec.AdditionalSecurityGroups, existingSecurityGroups)
	if err != nil {
		conditions.MarkFalse(machineScope.AWSMachine, infrav1.SecurityGroupsReadyCondition, infrav1.SecurityGroupsFailedReason, clusterv1.ConditionSeverityError, "%s", err.Error())
		machineScope.Error(err, "unable to ensure security groups")
		return err
	}
	conditions.MarkTrue(machineScope.AWSMachine, infrav1.SecurityGroupsReadyCondition)

	err = r.ensureInstanceMetadataOptions(ec2svc, instance, machineScope.AWSMachine)
	if err != nil {
		machineScope.Error(err, "failed to ensure instance metadata options")
		return err
	}

	return nil
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
		var ignitionStorageType infrav1.IgnitionStorageTypeOption
		if machineScope.AWSMachine.Spec.Ignition == nil {
			ignitionStorageType = infrav1.IgnitionStorageTypeOptionClusterObjectStore
		} else {
			ignitionStorageType = machineScope.AWSMachine.Spec.Ignition.StorageType
		}

		switch ignitionStorageType {
		case infrav1.IgnitionStorageTypeOptionClusterObjectStore:
			userData, err = r.generateIgnitionWithRemoteStorage(machineScope, objectStoreSvc, userData)
		case infrav1.IgnitionStorageTypeOptionUnencryptedUserData:
			// No further modifications to userdata are needed for plain storage in UnencryptedUserData.
		default:
			return nil, "", errors.Errorf("unsupported ignition storageType %q", ignitionStorageType)
		}
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

// generateIgnitionWithRemoteStorage uses a remote object storage (S3 bucket) and stores user data in it,
// then returns the config to instruct ignition on how to pull the user data from the bucket.
func (r *AWSMachineReconciler) generateIgnitionWithRemoteStorage(scope *scope.MachineScope, objectStoreSvc services.ObjectStoreInterface, userData []byte) ([]byte, error) {
	if objectStoreSvc == nil {
		return nil, errors.New("using Ignition by default requires a cluster wide object storage configured at `AWSCluster.Spec.Ignition.S3Bucket`. " +
			"You must configure one or instruct Ignition to use EC2 user data instead, by setting `AWSMachine.Spec.Ignition.StorageType` to `UnencryptedUserData`")
	}

	objectURL, err := objectStoreSvc.Create(scope, userData)
	if err != nil {
		return nil, errors.Wrap(err, "creating userdata object")
	}

	ignVersion := getIgnitionVersion(scope)
	semver, err := semver.ParseTolerant(ignVersion)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to parse ignition version %q", ignVersion)
	}

	switch semver.Major {
	case 2:
		ignData := &ignTypes.Config{
			Ignition: ignTypes.Ignition{
				Version: semver.String(),
				Config: ignTypes.IgnitionConfig{
					Append: []ignTypes.ConfigReference{
						{
							Source: objectURL,
						},
					},
				},
			},
		}

		return json.Marshal(ignData)
	case 3:
		ignData := &ignV3Types.Config{
			Ignition: ignV3Types.Ignition{
				Version: semver.String(),
				Config: ignV3Types.IgnitionConfig{
					Merge: []ignV3Types.Resource{
						{
							Source: aws.String(objectURL),
						},
					},
				},
			},
		}

		if scope.AWSMachine.Spec.Ignition.Proxy != nil {
			ignData.Ignition.Proxy = ignV3Types.Proxy{
				HTTPProxy:  scope.AWSMachine.Spec.Ignition.Proxy.HTTPProxy,
				HTTPSProxy: scope.AWSMachine.Spec.Ignition.Proxy.HTTPSProxy,
			}
			for _, noProxy := range scope.AWSMachine.Spec.Ignition.Proxy.NoProxy {
				ignData.Ignition.Proxy.NoProxy = append(ignData.Ignition.Proxy.NoProxy, ignV3Types.NoProxyItem(noProxy))
			}
		}

		if scope.AWSMachine.Spec.Ignition.TLS != nil {
			for _, cert := range scope.AWSMachine.Spec.Ignition.TLS.CASources {
				ignData.Ignition.Security.TLS.CertificateAuthorities = append(
					ignData.Ignition.Security.TLS.CertificateAuthorities,
					ignV3Types.Resource{Source: aws.String(string(cert))},
				)
			}
		}

		return json.Marshal(ignData)
	default:
		return nil, errors.Errorf("unsupported ignition version %q", ignVersion)
	}
}

func getIgnitionVersion(scope *scope.MachineScope) string {
	if scope.AWSMachine.Spec.Ignition == nil {
		scope.AWSMachine.Spec.Ignition = &infrav1.Ignition{}
	}
	if scope.AWSMachine.Spec.Ignition.Version == "" {
		scope.AWSMachine.Spec.Ignition.Version = infrav1.DefaultIgnitionVersion
	}
	return scope.AWSMachine.Spec.Ignition.Version
}

func (r *AWSMachineReconciler) deleteBootstrapData(machineScope *scope.MachineScope, clusterScope cloud.ClusterScoper, objectStoreScope scope.S3Scope) error {
	var userDataFormat string
	var err error
	if machineScope.Machine.Spec.Bootstrap.DataSecretName != nil {
		_, userDataFormat, err = machineScope.GetRawBootstrapDataWithFormat()
		if client.IgnoreNotFound(err) != nil {
			return errors.Wrap(err, "failed to get raw userdata")
		}
	}

	if machineScope.UseSecretsManager(userDataFormat) {
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

	// If bootstrap data has not been populated yet, we cannot determine its format, so there is probably nothing to do.
	if machineScope.Machine.Spec.Bootstrap.DataSecretName == nil {
		return nil
	}

	_, userDataFormat, err := machineScope.GetRawBootstrapDataWithFormat()
	if err != nil && !apierrors.IsNotFound(err) {
		r.Recorder.Eventf(machineScope.AWSMachine, corev1.EventTypeWarning, "FailedGetBootstrapData", err.Error())
		return err
	}

	// We only use an S3 bucket to store userdata if we use Ignition with StorageType ClusterObjectStore.
	if !machineScope.UseIgnition(userDataFormat) ||
		(machineScope.AWSMachine.Spec.Ignition != nil &&
			machineScope.AWSMachine.Spec.Ignition.StorageType != infrav1.IgnitionStorageTypeOptionClusterObjectStore) {
		return nil
	}

	if err := objectStoreSvc.Delete(machineScope); err != nil {
		return errors.Wrap(err, "deleting bootstrap data object")
	}

	return nil
}

// reconcileLBAttachment reconciles attachment to _all_ defined load balancers.
// Callers are expected to filter out known-good errors out of the aggregate error list.
func (r *AWSMachineReconciler) reconcileLBAttachment(machineScope *scope.MachineScope, elbScope scope.ELBScope, i *infrav1.Instance) error {
	if !machineScope.IsControlPlane() {
		return nil
	}

	elbsvc := r.getELBService(elbScope)

	errs := []error{}
	for _, lbSpec := range elbScope.ControlPlaneLoadBalancers() {
		if lbSpec == nil {
			continue
		}
		// In order to prevent sending request to a "not-ready" control plane machines, it is required to remove the machine
		// from the ELB as soon as the machine or infra machine gets deleted or when the machine is in a not running state.
		if machineScope.AWSMachineIsDeleted() || machineScope.MachineIsDeleted() || !machineScope.InstanceIsRunning() {
			if lbSpec.LoadBalancerType == infrav1.LoadBalancerTypeClassic {
				machineScope.Debug("deregistering from classic load balancer")
				return r.deregisterInstanceFromClassicLB(machineScope, elbsvc, i)
			}
			machineScope.Debug("deregistering from v2 load balancer")
			errs = append(errs, r.deregisterInstanceFromV2LB(machineScope, elbsvc, i, lbSpec))
			continue
		}

		if err := r.registerInstanceToLBs(machineScope, elbsvc, i, lbSpec); err != nil {
			errs = append(errs, errors.Wrapf(err, "could not register machine to load balancer"))
		}
	}

	return kerrors.NewAggregate(errs)
}

func (r *AWSMachineReconciler) registerInstanceToLBs(machineScope *scope.MachineScope, elbsvc services.ELBInterface, i *infrav1.Instance, lb *infrav1.AWSLoadBalancerSpec) error {
	switch lb.LoadBalancerType {
	case infrav1.LoadBalancerTypeClassic, "":
		machineScope.Debug("registering to classic load balancer")
		return r.registerInstanceToClassicLB(machineScope, elbsvc, i)
	case infrav1.LoadBalancerTypeELB, infrav1.LoadBalancerTypeALB, infrav1.LoadBalancerTypeNLB:
		machineScope.Debug("registering to v2 load balancer")
		return r.registerInstanceToV2LB(machineScope, elbsvc, i, lb)
	}

	return errors.Errorf("unknown load balancer type %q", lb.LoadBalancerType)
}

func (r *AWSMachineReconciler) registerInstanceToClassicLB(machineScope *scope.MachineScope, elbsvc services.ELBInterface, i *infrav1.Instance) error {
	registered, err := elbsvc.IsInstanceRegisteredWithAPIServerELB(i)
	if err != nil {
		r.Recorder.Eventf(machineScope.AWSMachine, corev1.EventTypeWarning, "FailedAttachControlPlaneELB",
			"Failed to register control plane instance %q with classic load balancer: failed to determine registration status: %v", i.ID, err)
		return errors.Wrapf(err, "could not register control plane instance %q with classic load balancer - error determining registration status", i.ID)
	}
	if registered {
		// Already registered - nothing more to do
		return nil
	}

	if err := elbsvc.RegisterInstanceWithAPIServerELB(i); err != nil {
		r.Recorder.Eventf(machineScope.AWSMachine, corev1.EventTypeWarning, "FailedAttachControlPlaneELB",
			"Failed to register control plane instance %q with classic load balancer: %v", i.ID, err)
		conditions.MarkFalse(machineScope.AWSMachine, infrav1.ELBAttachedCondition, infrav1.ELBAttachFailedReason, clusterv1.ConditionSeverityError, "%s", err.Error())
		return errors.Wrapf(err, "could not register control plane instance %q with classic load balancer", i.ID)
	}
	r.Recorder.Eventf(machineScope.AWSMachine, corev1.EventTypeNormal, "SuccessfulAttachControlPlaneELB",
		"Control plane instance %q is registered with classic load balancer", i.ID)
	conditions.MarkTrue(machineScope.AWSMachine, infrav1.ELBAttachedCondition)
	return nil
}

func (r *AWSMachineReconciler) registerInstanceToV2LB(machineScope *scope.MachineScope, elbsvc services.ELBInterface, instance *infrav1.Instance, lb *infrav1.AWSLoadBalancerSpec) error {
	_, registered, err := elbsvc.IsInstanceRegisteredWithAPIServerLB(instance, lb)
	if err != nil {
		r.Recorder.Eventf(machineScope.AWSMachine, corev1.EventTypeWarning, "FailedAttachControlPlaneELB",
			"Failed to register control plane instance %q with load balancer: failed to determine registration status: %v", instance.ID, err)
		return errors.Wrapf(err, "could not register control plane instance %q with load balancer - error determining registration status", instance.ID)
	}
	if registered {
		machineScope.Logger.Debug("Instance is already registered.", "instance", instance.ID)
		return nil
	}

	// See https://docs.aws.amazon.com/elasticloadbalancing/latest/application/target-group-register-targets.html#register-instances
	if ptr.Deref(machineScope.GetInstanceState(), infrav1.InstanceStatePending) != infrav1.InstanceStateRunning {
		r.Recorder.Eventf(machineScope.AWSMachine, corev1.EventTypeWarning, "FailedAttachControlPlaneELB",
			"Cannot register control plane instance %q with load balancer: instance is not running", instance.ID)
		conditions.MarkFalse(machineScope.AWSMachine, infrav1.ELBAttachedCondition, infrav1.ELBAttachFailedReason, clusterv1.ConditionSeverityInfo, "instance not running")
		return elb.NewInstanceNotRunning("instance is not running")
	}

	if err := elbsvc.RegisterInstanceWithAPIServerLB(instance, lb); err != nil {
		r.Recorder.Eventf(machineScope.AWSMachine, corev1.EventTypeWarning, "FailedAttachControlPlaneELB",
			"Failed to register control plane instance %q with load balancer: %v", instance.ID, err)
		conditions.MarkFalse(machineScope.AWSMachine, infrav1.ELBAttachedCondition, infrav1.ELBAttachFailedReason, clusterv1.ConditionSeverityError, "%s", err.Error())
		return errors.Wrapf(err, "could not register control plane instance %q with load balancer", instance.ID)
	}
	r.Recorder.Eventf(machineScope.AWSMachine, corev1.EventTypeNormal, "SuccessfulAttachControlPlaneELB",
		"Control plane instance %q is registered with load balancer", instance.ID)
	conditions.MarkTrue(machineScope.AWSMachine, infrav1.ELBAttachedCondition)
	return nil
}

func (r *AWSMachineReconciler) deregisterInstanceFromClassicLB(machineScope *scope.MachineScope, elbsvc services.ELBInterface, instance *infrav1.Instance) error {
	registered, err := elbsvc.IsInstanceRegisteredWithAPIServerELB(instance)
	if err != nil {
		r.Recorder.Eventf(machineScope.AWSMachine, corev1.EventTypeWarning, "FailedDetachControlPlaneELB",
			"Failed to deregister control plane instance %q from load balancer: failed to determine registration status: %v", instance.ID, err)
		return errors.Wrapf(err, "could not deregister control plane instance %q from load balancer - error determining registration status", instance.ID)
	}
	if !registered {
		machineScope.Logger.Debug("Instance is already registered.", "instance", instance.ID)
		return nil
	}

	if err := elbsvc.DeregisterInstanceFromAPIServerELB(instance); err != nil {
		r.Recorder.Eventf(machineScope.AWSMachine, corev1.EventTypeWarning, "FailedDetachControlPlaneELB",
			"Failed to deregister control plane instance %q from load balancer: %v", instance.ID, err)
		conditions.MarkFalse(machineScope.AWSMachine, infrav1.ELBAttachedCondition, infrav1.ELBDetachFailedReason, clusterv1.ConditionSeverityError, "%s", err.Error())
		return errors.Wrapf(err, "could not deregister control plane instance %q from load balancer", instance.ID)
	}

	r.Recorder.Eventf(machineScope.AWSMachine, corev1.EventTypeNormal, "SuccessfulDetachControlPlaneELB",
		"Control plane instance %q is de-registered from load balancer", instance.ID)
	return nil
}

func (r *AWSMachineReconciler) deregisterInstanceFromV2LB(machineScope *scope.MachineScope, elbsvc services.ELBInterface, i *infrav1.Instance, lb *infrav1.AWSLoadBalancerSpec) error {
	targetGroupARNs, registered, err := elbsvc.IsInstanceRegisteredWithAPIServerLB(i, lb)
	if err != nil {
		r.Recorder.Eventf(machineScope.AWSMachine, corev1.EventTypeWarning, "FailedDetachControlPlaneELB",
			"Failed to deregister control plane instance %q from load balancer: failed to determine registration status: %v", i.ID, err)
		return errors.Wrapf(err, "could not deregister control plane instance %q from load balancer - error determining registration status", i.ID)
	}
	if !registered {
		// Already deregistered - nothing more to do
		return nil
	}

	for _, targetGroupArn := range targetGroupARNs {
		if err := elbsvc.DeregisterInstanceFromAPIServerLB(targetGroupArn, i); err != nil {
			r.Recorder.Eventf(machineScope.AWSMachine, corev1.EventTypeWarning, "FailedDetachControlPlaneELB",
				"Failed to deregister control plane instance %q from load balancer: %v", i.ID, err)
			conditions.MarkFalse(machineScope.AWSMachine, infrav1.ELBAttachedCondition, infrav1.ELBDetachFailedReason, clusterv1.ConditionSeverityError, "%s", err.Error())
			return errors.Wrapf(err, "could not deregister control plane instance %q from load balancer", i.ID)
		}
	}

	r.Recorder.Eventf(machineScope.AWSMachine, corev1.EventTypeNormal, "SuccessfulDetachControlPlaneELB",
		"Control plane instance %q is de-registered from load balancer", i.ID)
	return nil
}

// AWSClusterToAWSMachines is a handler.ToRequestsFunc to be used to enqeue requests for reconciliation
// of AWSMachines.
func (r *AWSMachineReconciler) AWSClusterToAWSMachines(log logger.Wrapper) handler.MapFunc {
	return func(ctx context.Context, o client.Object) []ctrl.Request {
		c, ok := o.(*infrav1.AWSCluster)
		if !ok {
			klog.Errorf("Expected a AWSCluster but got a %T", o)
		}

		log := log.WithValues("objectMapper", "awsClusterToAWSMachine", "cluster", klog.KRef(c.Namespace, c.Name))

		// Don't handle deleted AWSClusters
		if !c.ObjectMeta.DeletionTimestamp.IsZero() {
			log.Trace("AWSCluster has a deletion timestamp, skipping mapping.")
			return nil
		}

		cluster, err := util.GetOwnerCluster(ctx, r.Client, c.ObjectMeta)
		switch {
		case apierrors.IsNotFound(err) || cluster == nil:
			log.Trace("Cluster for AWSCluster not found, skipping mapping.")
			return nil
		case err != nil:
			log.Error(err, "Failed to get owning cluster, skipping mapping.")
			return nil
		}

		return r.requestsForCluster(log, cluster.Namespace, cluster.Name)
	}
}

func (r *AWSMachineReconciler) requeueAWSMachinesForUnpausedCluster(log logger.Wrapper) handler.MapFunc {
	return func(ctx context.Context, o client.Object) []ctrl.Request {
		c, ok := o.(*clusterv1.Cluster)
		if !ok {
			klog.Errorf("Expected a Cluster but got a %T", o)
		}

		log := log.WithValues("objectMapper", "clusterToAWSMachine", "cluster", klog.KRef(c.Namespace, c.Name))

		// Don't handle deleted clusters
		if !c.ObjectMeta.DeletionTimestamp.IsZero() {
			log.Trace("Cluster has a deletion timestamp, skipping mapping.")
			return nil
		}

		return r.requestsForCluster(log, c.Namespace, c.Name)
	}
}

func (r *AWSMachineReconciler) requestsForCluster(log logger.Wrapper, namespace, name string) []ctrl.Request {
	labels := map[string]string{clusterv1.ClusterNameLabel: name}
	machineList := &clusterv1.MachineList{}
	if err := r.Client.List(context.TODO(), machineList, client.InNamespace(namespace), client.MatchingLabels(labels)); err != nil {
		log.Error(err, "Failed to get owned Machines, skipping mapping.")
		return nil
	}

	result := make([]ctrl.Request, 0, len(machineList.Items))
	for _, m := range machineList.Items {
		log.WithValues("machine", klog.KObj(&m))
		if m.Spec.InfrastructureRef.GroupVersionKind().Kind != "AWSMachine" {
			log.Trace("Machine has an InfrastructureRef for a different type, will not add to reconciliation request.")
			continue
		}
		if m.Spec.InfrastructureRef.Name == "" {
			log.Trace("Machine has an InfrastructureRef with an empty name, will not add to reconciliation request.")
			continue
		}
		log.WithValues("awsMachine", klog.KRef(m.Spec.InfrastructureRef.Namespace, m.Spec.InfrastructureRef.Name))
		log.Trace("Adding AWSMachine to reconciliation request.")
		result = append(result, ctrl.Request{NamespacedName: client.ObjectKey{Namespace: m.Namespace, Name: m.Spec.InfrastructureRef.Name}})
	}
	return result
}

func (r *AWSMachineReconciler) getInfraCluster(ctx context.Context, log *logger.Logger, cluster *clusterv1.Cluster, awsMachine *infrav1.AWSMachine) (scope.EC2Scope, error) {
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
			return nil, nil //nolint:nilerr
		}

		managedControlPlaneScope, err = scope.NewManagedControlPlaneScope(scope.ManagedControlPlaneScopeParams{
			Client:                       r.Client,
			Logger:                       log,
			Cluster:                      cluster,
			ControlPlane:                 controlPlane,
			ControllerName:               "awsManagedControlPlane",
			Endpoints:                    r.Endpoints,
			TagUnmanagedNetworkResources: r.TagUnmanagedNetworkResources,
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
		return nil, nil //nolint:nilerr
	}

	// Create the cluster scope
	clusterScope, err = scope.NewClusterScope(scope.ClusterScopeParams{
		Client:                       r.Client,
		Logger:                       log,
		Cluster:                      cluster,
		AWSCluster:                   awsCluster,
		ControllerName:               "awsmachine",
		TagUnmanagedNetworkResources: r.TagUnmanagedNetworkResources,
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

func (r *AWSMachineReconciler) ensureStorageTags(ec2svc services.EC2Interface, instance *infrav1.Instance, machine *infrav1.AWSMachine, additionalTags map[string]string) {
	prevAnnotations, err := r.machineAnnotationJSON(machine, VolumeTagsLastAppliedAnnotation)
	if err != nil {
		r.Log.Error(err, "Failed to fetch the annotations for volume tags")
	}
	annotations := make(map[string]interface{}, len(instance.VolumeIDs))
	for _, volumeID := range instance.VolumeIDs {
		if subAnnotation, ok := prevAnnotations[volumeID].(map[string]interface{}); ok {
			newAnnotation, err := r.ensureVolumeTags(ec2svc, aws.String(volumeID), subAnnotation, additionalTags)
			if err != nil {
				r.Log.Error(err, "Failed to fetch the changed volume tags in EC2 instance")
			}
			annotations[volumeID] = newAnnotation
		} else {
			newAnnotation, err := r.ensureVolumeTags(ec2svc, aws.String(volumeID), make(map[string]interface{}), additionalTags)
			if err != nil {
				r.Log.Error(err, "Failed to fetch the changed volume tags in EC2 instance")
			}
			annotations[volumeID] = newAnnotation
		}
	}

	if !cmp.Equal(prevAnnotations, annotations, cmpopts.EquateEmpty()) {
		// We also need to update the annotation if anything changed.
		err = r.updateMachineAnnotationJSON(machine, VolumeTagsLastAppliedAnnotation, annotations)
		if err != nil {
			r.Log.Error(err, "Failed to fetch the changed volume tags in EC2 instance")
		}
	}
}

func (r *AWSMachineReconciler) ensureInstanceMetadataOptions(ec2svc services.EC2Interface, instance *infrav1.Instance, machine *infrav1.AWSMachine) error {
	if cmp.Equal(machine.Spec.InstanceMetadataOptions, instance.InstanceMetadataOptions) {
		return nil
	}

	return ec2svc.ModifyInstanceMetadataOptions(instance.ID, machine.Spec.InstanceMetadataOptions)
}
