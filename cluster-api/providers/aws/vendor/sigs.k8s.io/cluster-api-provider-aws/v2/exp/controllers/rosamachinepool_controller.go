package controllers

import (
	"context"
	"fmt"
	"time"

	"github.com/aws/aws-sdk-go-v2/service/ec2"
	ec2types "github.com/aws/aws-sdk-go-v2/service/ec2/types"
	"github.com/blang/semver"
	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	cmv1 "github.com/openshift-online/ocm-sdk-go/clustersmgmt/v1"
	"github.com/openshift/rosa/pkg/ocm"
	"github.com/pkg/errors"
	corev1 "k8s.io/api/core/v1"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/client-go/tools/record"
	"k8s.io/klog/v2"
	"k8s.io/utils/ptr"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/client/apiutil"
	"sigs.k8s.io/controller-runtime/pkg/controller"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
	"sigs.k8s.io/controller-runtime/pkg/handler"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"

	rosacontrolplanev1 "sigs.k8s.io/cluster-api-provider-aws/v2/controlplane/rosa/api/v1beta2"
	expinfrav1 "sigs.k8s.io/cluster-api-provider-aws/v2/exp/api/v1beta2"
	"sigs.k8s.io/cluster-api-provider-aws/v2/exp/utils"
	"sigs.k8s.io/cluster-api-provider-aws/v2/pkg/cloud"
	"sigs.k8s.io/cluster-api-provider-aws/v2/pkg/cloud/scope"
	stsservice "sigs.k8s.io/cluster-api-provider-aws/v2/pkg/cloud/services/sts"
	"sigs.k8s.io/cluster-api-provider-aws/v2/pkg/logger"
	"sigs.k8s.io/cluster-api-provider-aws/v2/pkg/rosa"
	"sigs.k8s.io/cluster-api-provider-aws/v2/util/paused"
	clusterv1beta1 "sigs.k8s.io/cluster-api/api/core/v1beta1"
	clusterv1 "sigs.k8s.io/cluster-api/api/core/v1beta2"
	"sigs.k8s.io/cluster-api/util"
	"sigs.k8s.io/cluster-api/util/annotations"
	v1beta1conditions "sigs.k8s.io/cluster-api/util/deprecated/v1beta1/conditions"
	"sigs.k8s.io/cluster-api/util/deprecated/v1beta1/patch"
	"sigs.k8s.io/cluster-api/util/predicates"
)

// ROSAMachinePoolReconciler reconciles a ROSAMachinePool object.
type ROSAMachinePoolReconciler struct {
	client.Client
	Recorder         record.EventRecorder
	WatchFilterValue string
	NewStsClient     func(cloud.ScopeUsage, cloud.Session, logger.Wrapper, runtime.Object) stsservice.STSClient
	NewOCMClient     func(ctx context.Context, rosaScope *scope.ROSAControlPlaneScope) (rosa.OCMClient, error)
}

// SetupWithManager is used to setup the controller.
func (r *ROSAMachinePoolReconciler) SetupWithManager(ctx context.Context, mgr ctrl.Manager, options controller.Options) error {
	log := logger.FromContext(ctx)
	r.NewOCMClient = rosa.NewWrappedOCMClient
	r.NewStsClient = scope.NewSTSClient

	gvk, err := apiutil.GVKForObject(new(expinfrav1.ROSAMachinePool), mgr.GetScheme())
	if err != nil {
		return errors.Wrapf(err, "failed to find GVK for ROSAMachinePool")
	}
	rosaControlPlaneToRosaMachinePoolMap := rosaControlPlaneToRosaMachinePoolMapFunc(r.Client, gvk, log)
	return ctrl.NewControllerManagedBy(mgr).
		For(&expinfrav1.ROSAMachinePool{}).
		WithOptions(options).
		WithEventFilter(predicates.ResourceHasFilterLabel(mgr.GetScheme(), log.GetLogger(), r.WatchFilterValue)).
		Watches(
			&clusterv1.MachinePool{},
			handler.EnqueueRequestsFromMapFunc(machinePoolToInfrastructureMapFunc(gvk)),
		).
		Watches(
			&rosacontrolplanev1.ROSAControlPlane{},
			handler.EnqueueRequestsFromMapFunc(rosaControlPlaneToRosaMachinePoolMap),
		).
		Complete(r)
}

// +kubebuilder:rbac:groups=cluster.x-k8s.io,resources=machinepools;machinepools/status,verbs=get;list;watch;patch
// +kubebuilder:rbac:groups=core,resources=events,verbs=get;list;watch;create;patch
// +kubebuilder:rbac:groups=controlplane.cluster.x-k8s.io,resources=rosacontrolplanes;rosacontrolplanes/status,verbs=get;list;watch
// +kubebuilder:rbac:groups=infrastructure.cluster.x-k8s.io,resources=rosamachinepools,verbs=get;list;watch;update;patch;delete
// +kubebuilder:rbac:groups=infrastructure.cluster.x-k8s.io,resources=rosamachinepools/status,verbs=get;update;patch
// +kubebuilder:rbac:groups=infrastructure.cluster.x-k8s.io,resources=rosamachinepools/finalizers,verbs=update

// Reconcile reconciles ROSAMachinePool.
func (r *ROSAMachinePoolReconciler) Reconcile(ctx context.Context, req ctrl.Request) (_ ctrl.Result, reterr error) {
	log := logger.FromContext(ctx)

	rosaMachinePool := &expinfrav1.ROSAMachinePool{}
	if err := r.Get(ctx, req.NamespacedName, rosaMachinePool); err != nil {
		if apierrors.IsNotFound(err) {
			return ctrl.Result{}, nil
		}
		return ctrl.Result{Requeue: true}, nil
	}

	machinePool, err := getOwnerMachinePool(ctx, r.Client, rosaMachinePool.ObjectMeta)
	if err != nil {
		log.Error(err, "Failed to retrieve owner MachinePool from the API Server")
		return ctrl.Result{}, err
	}
	if machinePool == nil {
		log.Info("MachinePool Controller has not yet set OwnerRef")
		return ctrl.Result{}, nil
	}

	log = log.WithValues("MachinePool", klog.KObj(machinePool))

	cluster, err := util.GetClusterFromMetadata(ctx, r.Client, machinePool.ObjectMeta)
	if err != nil {
		log.Info("Failed to retrieve Cluster from MachinePool")
		return ctrl.Result{}, nil
	}

	if isPaused, conditionChanged, err := paused.EnsurePausedCondition(ctx, r.Client, cluster, rosaMachinePool); err != nil || isPaused || conditionChanged {
		return ctrl.Result{}, err
	}

	log = log.WithValues("cluster", klog.KObj(cluster))

	controlPlaneKey := client.ObjectKey{
		Namespace: rosaMachinePool.Namespace,
		Name:      cluster.Spec.ControlPlaneRef.Name,
	}
	controlPlane := &rosacontrolplanev1.ROSAControlPlane{}
	if err := r.Client.Get(ctx, controlPlaneKey, controlPlane); err != nil {
		if apierrors.IsNotFound(err) && !rosaMachinePool.DeletionTimestamp.IsZero() {
			// When the ROSAControlPlane is not found and the ROSAMachinePool CR is marked for deletion,
			// it indicates that the ROSAControlPlane (and its associated NodePools) has already been deleted,
			// while the ROSAMachinePool remains pending — since a ROSA-HCP cluster cannot exist without a NodePool.
			// To handle this scenario, we trigger deletion of the ROSAControlPlane CR to initiate cleanup of the ROSA-HCP,
			// relying on OCM to cascade-delete the related NodePools.
			// Note: This state should rarely occur. However, during smoke tests, the ROSAMachinePool reconcile cycle
			// may occasionally lag behind the deletion of the NodePools and ROSAControlPlane.
			log.Info("RosaControlPlane not found, RosaMachinePool is deleted")
			patchHelper, err := patch.NewHelper(rosaMachinePool, r.Client)
			if err != nil {
				return ctrl.Result{}, errors.Wrap(err, "failed to init RosaMachinePool patch helper")
			}

			controllerutil.RemoveFinalizer(rosaMachinePool, expinfrav1.RosaMachinePoolFinalizer)
			return ctrl.Result{}, patchHelper.Patch(ctx, rosaMachinePool, patch.WithOwnedConditions{Conditions: []clusterv1beta1.ConditionType{
				expinfrav1.RosaMachinePoolReadyCondition,
			}})
		}

		log.Info("Failed to retrieve ControlPlane from MachinePool")
		return ctrl.Result{}, err
	}

	machinePoolScope, err := scope.NewRosaMachinePoolScope(scope.RosaMachinePoolScopeParams{
		Client:          r.Client,
		ControllerName:  "rosamachinepool",
		Cluster:         cluster,
		ControlPlane:    controlPlane,
		MachinePool:     machinePool,
		RosaMachinePool: rosaMachinePool,
		Logger:          log,
	})
	if err != nil {
		return ctrl.Result{}, errors.Wrap(err, "failed to create rosaMachinePool scope")
	}

	rosaControlPlaneScope, err := scope.NewROSAControlPlaneScope(scope.ROSAControlPlaneScopeParams{
		Client:         r.Client,
		Cluster:        cluster,
		ControlPlane:   controlPlane,
		ControllerName: "rosaControlPlane",
		NewStsClient:   r.NewStsClient,
	})
	if err != nil {
		return ctrl.Result{}, errors.Wrap(err, "failed to create rosaControlPlane scope")
	}

	if !controlPlane.Status.Ready && controlPlane.ObjectMeta.DeletionTimestamp.IsZero() {
		log.Info("Control plane is not ready yet")
		err := machinePoolScope.RosaMachinePoolReadyFalse(expinfrav1.WaitingForRosaControlPlaneReason, "")
		return ctrl.Result{}, err
	}

	defer func() {
		v1beta1conditions.SetSummary(machinePoolScope.RosaMachinePool, v1beta1conditions.WithConditions(expinfrav1.RosaMachinePoolReadyCondition), v1beta1conditions.WithStepCounter())

		if err := machinePoolScope.Close(); err != nil && reterr == nil {
			reterr = err
		}
	}()

	if !rosaMachinePool.ObjectMeta.DeletionTimestamp.IsZero() {
		return ctrl.Result{}, r.reconcileDelete(ctx, machinePoolScope, rosaControlPlaneScope)
	}

	return r.reconcileNormal(ctx, machinePoolScope, rosaControlPlaneScope)
}

func (r *ROSAMachinePoolReconciler) reconcileNormal(ctx context.Context,
	machinePoolScope *scope.RosaMachinePoolScope,
	rosaControlPlaneScope *scope.ROSAControlPlaneScope,
) (ctrl.Result, error) {
	machinePoolScope.Info("Reconciling ROSAMachinePool")

	if controllerutil.AddFinalizer(machinePoolScope.RosaMachinePool, expinfrav1.RosaMachinePoolFinalizer) {
		if err := machinePoolScope.PatchObject(); err != nil {
			return ctrl.Result{}, err
		}
	}
	if r.NewOCMClient == nil {
		return ctrl.Result{}, fmt.Errorf("failed to create OCM client: NewOCMClient is nil")
	}

	ocmClient, err := r.NewOCMClient(ctx, rosaControlPlaneScope)
	if err != nil || ocmClient == nil {
		// TODO: need to expose in status, as likely the credentials are invalid
		return ctrl.Result{}, fmt.Errorf("failed to create OCM client: %w", err)
	}

	failureMessage, err := validateMachinePoolSpec(machinePoolScope)
	if err != nil {
		return ctrl.Result{}, fmt.Errorf("failed to validate ROSAMachinePool.spec: %w", err)
	}
	if failureMessage != nil {
		// dont' requeue because input is invalid and manual intervention is needed.
		return ctrl.Result{}, nil
	}
	machinePoolScope.RosaMachinePool.Status.FailureMessage = nil

	rosaMachinePool := machinePoolScope.RosaMachinePool
	machinePool := machinePoolScope.MachinePool

	if rosaMachinePool.Spec.Autoscaling != nil && !annotations.ReplicasManagedByExternalAutoscaler(machinePool) {
		// make sure cluster.x-k8s.io/replicas-managed-by annotation is set on CAPI MachinePool when autoscaling is enabled.
		annotations.AddAnnotations(machinePool, map[string]string{
			clusterv1.ReplicasManagedByAnnotation: "rosa",
		})
		if err := machinePoolScope.PatchCAPIMachinePoolObject(ctx); err != nil {
			return ctrl.Result{}, err
		}
	}

	nodePool, found, err := ocmClient.GetNodePool(machinePoolScope.ControlPlane.Status.ID, rosaMachinePool.Spec.NodePoolName)
	if err != nil {
		return ctrl.Result{}, err
	}
	if found {
		if rosaMachinePool.Spec.AvailabilityZone == "" {
			// reflect the current AvailabilityZone in the spec if not set.
			rosaMachinePool.Spec.AvailabilityZone = nodePool.AvailabilityZone()
		}

		nodePool, err := r.updateNodePool(machinePoolScope, ocmClient, nodePool)
		if err != nil {
			return ctrl.Result{}, fmt.Errorf("failed to ensure rosaMachinePool: %w", err)
		}

		currentReplicas := int32(nodePool.Status().CurrentReplicas()) //#nosec G115
		if annotations.ReplicasManagedByExternalAutoscaler(machinePool) {
			// Set MachinePool replicas to rosa autoscaling replicas
			if *machinePool.Spec.Replicas != currentReplicas {
				machinePoolScope.Info("Setting MachinePool replicas to rosa autoscaling replicas",
					"local", *machinePool.Spec.Replicas,
					"external", currentReplicas)
				machinePool.Spec.Replicas = &currentReplicas
				if err := machinePoolScope.PatchCAPIMachinePoolObject(ctx); err != nil {
					return ctrl.Result{}, err
				}
			}
		}
		if err := r.reconcileProviderIDList(ctx, machinePoolScope, nodePool); err != nil {
			return ctrl.Result{}, fmt.Errorf("failed to reconcile ProviderIDList: %w", err)
		}

		rosaMachinePool.Status.Replicas = currentReplicas
		if rosa.IsNodePoolReady(nodePool) {
			v1beta1conditions.MarkTrue(rosaMachinePool, expinfrav1.RosaMachinePoolReadyCondition)
			rosaMachinePool.Status.Ready = true

			if err := r.reconcileMachinePoolVersion(machinePoolScope, ocmClient, nodePool); err != nil {
				return ctrl.Result{}, err
			}

			return ctrl.Result{}, nil
		}

		v1beta1conditions.MarkFalse(rosaMachinePool,
			expinfrav1.RosaMachinePoolReadyCondition,
			nodePool.Status().Message(),
			clusterv1beta1.ConditionSeverityInfo,
			"")

		machinePoolScope.Info("waiting for NodePool to become ready", "state", nodePool.Status().Message())
		// Requeue so that status.ready is set to true when the nodepool is fully created.
		return ctrl.Result{RequeueAfter: time.Second * 60}, nil
	}

	npBuilder := nodePoolBuilder(rosaMachinePool.Spec, machinePool.Spec, machinePoolScope.ControlPlane.Spec.ChannelGroup)
	nodePoolSpec, err := npBuilder.Build()
	if err != nil {
		return ctrl.Result{}, fmt.Errorf("failed to build rosa nodepool: %w", err)
	}

	nodePool, err = ocmClient.CreateNodePool(machinePoolScope.ControlPlane.Status.ID, nodePoolSpec)
	if err != nil {
		v1beta1conditions.MarkFalse(rosaMachinePool,
			expinfrav1.RosaMachinePoolReadyCondition,
			expinfrav1.RosaMachinePoolReconciliationFailedReason,
			clusterv1beta1.ConditionSeverityError,
			"failed to create ROSAMachinePool: %s", err.Error())
		return ctrl.Result{}, fmt.Errorf("failed to create nodepool: %w", err)
	}

	machinePoolScope.RosaMachinePool.Status.ID = nodePool.ID()
	return ctrl.Result{}, nil
}

func (r *ROSAMachinePoolReconciler) reconcileDelete(
	ctx context.Context, machinePoolScope *scope.RosaMachinePoolScope,
	rosaControlPlaneScope *scope.ROSAControlPlaneScope,
) error {
	machinePoolScope.Info("Reconciling deletion of RosaMachinePool")

	ocmClient, err := r.NewOCMClient(ctx, rosaControlPlaneScope)
	if err != nil || ocmClient == nil {
		// TODO: need to expose in status, as likely the credentials are invalid
		return fmt.Errorf("failed to create OCM client: %w", err)
	}

	nodePool, found, err := ocmClient.GetNodePool(machinePoolScope.ControlPlane.Status.ID, machinePoolScope.NodePoolName())
	if err != nil {
		return err
	}
	if found {
		if err := ocmClient.DeleteNodePool(machinePoolScope.ControlPlane.Status.ID, nodePool.ID()); err != nil {
			return err
		}
		machinePoolScope.Info("Successfully deleted NodePool")
	}

	controllerutil.RemoveFinalizer(machinePoolScope.RosaMachinePool, expinfrav1.RosaMachinePoolFinalizer)

	return nil
}

func (r *ROSAMachinePoolReconciler) reconcileMachinePoolVersion(machinePoolScope *scope.RosaMachinePoolScope, ocmClient rosa.OCMClient, nodePool *cmv1.NodePool) error {
	version := machinePoolScope.RosaMachinePool.Spec.Version
	if version == "" || version == rosa.RawVersionID(nodePool.Version()) {
		machinePoolScope.RosaMachinePool.Status.AvailableUpgrades = nodePool.Version().AvailableUpgrades()
		v1beta1conditions.MarkFalse(machinePoolScope.RosaMachinePool, expinfrav1.RosaMachinePoolUpgradingCondition, "upgraded", clusterv1beta1.ConditionSeverityInfo, "")
		return nil
	}

	clusterID := machinePoolScope.ControlPlane.Status.ID
	_, scheduledUpgrade, err := ocmClient.GetHypershiftNodePoolUpgrade(clusterID, machinePoolScope.ControlPlane.Spec.RosaClusterName, nodePool.ID())
	if err != nil {
		return fmt.Errorf("failed to get existing scheduled upgrades: %w", err)
	}

	if scheduledUpgrade == nil {
		scheduledUpgrade, err = rosa.ScheduleNodePoolUpgrade(ocmClient, clusterID, nodePool, version, time.Now())
		if err != nil {
			return fmt.Errorf("failed to schedule nodePool upgrade to version %s: %w", version, err)
		}
	}

	condition := &clusterv1beta1.Condition{
		Type:    expinfrav1.RosaMachinePoolUpgradingCondition,
		Status:  corev1.ConditionTrue,
		Reason:  string(scheduledUpgrade.State().Value()),
		Message: fmt.Sprintf("Upgrading to version %s", scheduledUpgrade.Version()),
	}
	v1beta1conditions.Set(machinePoolScope.RosaMachinePool, condition)

	// if nodePool is already upgrading to another version we need to wait until the current upgrade is finished, return an error to requeue and try later.
	if scheduledUpgrade.Version() != version {
		return fmt.Errorf("there is already a %s upgrade to version %s", scheduledUpgrade.State().Value(), scheduledUpgrade.Version())
	}

	return nil
}

func (r *ROSAMachinePoolReconciler) shouldUpdateRosaReplicas(machinePoolScope *scope.RosaMachinePoolScope, nodePool *cmv1.NodePool) bool {
	if machinePoolScope.MachinePool.Spec.Replicas == nil || machinePoolScope.RosaMachinePool.Spec.Autoscaling != nil || annotations.ReplicasManagedByExternalAutoscaler(machinePoolScope.MachinePool) {
		return false
	}

	return nodePool.Replicas() != int(*machinePoolScope.MachinePool.Spec.Replicas)
}

func (r *ROSAMachinePoolReconciler) updateNodePool(machinePoolScope *scope.RosaMachinePoolScope, ocmClient rosa.OCMClient, nodePool *cmv1.NodePool) (*cmv1.NodePool, error) {
	machinePool := machinePoolScope.RosaMachinePool.DeepCopy()
	// default all fields before comparing, so that nil/unset fields don't cause an unnecessary update call.
	machinePool.Default()
	desiredSpec := machinePool.Spec

	specDiff := computeSpecDiff(desiredSpec, nodePool)
	// Replicas are not part of RosaMachinePoolSpec
	if specDiff == "" && !r.shouldUpdateRosaReplicas(machinePoolScope, nodePool) {
		// no changes detected.
		return nodePool, nil
	}
	machinePoolScope.Info("MachinePool spec diff detected", "diff", specDiff)

	// zero-out fields that shouldn't be part of the update call.
	desiredSpec.Version = ""
	desiredSpec.AdditionalSecurityGroups = nil
	desiredSpec.AdditionalTags = nil

	npBuilder := nodePoolBuilder(desiredSpec, machinePoolScope.MachinePool.Spec, machinePoolScope.ControlPlane.Spec.ChannelGroup)
	nodePoolSpec, err := npBuilder.Build()
	if err != nil {
		return nil, fmt.Errorf("failed to build nodePool spec: %w", err)
	}

	updatedNodePool, err := ocmClient.UpdateNodePool(machinePoolScope.ControlPlane.Status.ID, nodePoolSpec)
	if err != nil {
		v1beta1conditions.MarkFalse(machinePoolScope.RosaMachinePool,
			expinfrav1.RosaMachinePoolReadyCondition,
			expinfrav1.RosaMachinePoolReconciliationFailedReason,
			clusterv1beta1.ConditionSeverityError,
			"failed to update ROSAMachinePool: %s", err.Error())
		return nil, fmt.Errorf("failed to update nodePool: %w", err)
	}

	return updatedNodePool, nil
}

func computeSpecDiff(desiredSpec expinfrav1.RosaMachinePoolSpec, nodePool *cmv1.NodePool) string {
	currentSpec := utils.NodePoolToRosaMachinePoolSpec(nodePool)

	ignoredFields := []string{
		"ProviderIDList",           // providerIDList is set by the controller.
		"Version",                  // Version changes are reconciled separately.
		"AdditionalTags",           // AdditionalTags day2 changes not supported.
		"AdditionalSecurityGroups", // AdditionalSecurityGroups day2 changes not supported.
	}

	return cmp.Diff(desiredSpec, currentSpec,
		cmpopts.EquateEmpty(), // ensures empty non-nil slices and nil slices are considered equal.
		cmpopts.IgnoreFields(currentSpec, ignoredFields...))
}

func validateMachinePoolSpec(machinePoolScope *scope.RosaMachinePoolScope) (*string, error) {
	if machinePoolScope.RosaMachinePool.Spec.Version == "" {
		return nil, nil
	}

	version, err := semver.Parse(machinePoolScope.RosaMachinePool.Spec.Version)
	if err != nil {
		return nil, fmt.Errorf("failed to parse MachinePool version: %w", err)
	}
	minSupportedVersion, maxSupportedVersion, err := rosa.MachinePoolSupportedVersionsRange(machinePoolScope.ControlPlane.Spec.Version)
	if err != nil {
		return nil, fmt.Errorf("failed to get supported machinePool versions range: %w", err)
	}

	if version.GT(*maxSupportedVersion) || version.LT(*minSupportedVersion) {
		message := fmt.Sprintf("version %s is not supported, should be in the range: >= %s and <= %s", version, minSupportedVersion, maxSupportedVersion)
		return &message, nil
	}

	// TODO: add more input validations
	return nil, nil
}

func nodePoolBuilder(rosaMachinePoolSpec expinfrav1.RosaMachinePoolSpec, machinePoolSpec clusterv1.MachinePoolSpec, controlPlaneChannelGroup rosacontrolplanev1.ChannelGroupType) *cmv1.NodePoolBuilder {
	npBuilder := cmv1.NewNodePool().ID(rosaMachinePoolSpec.NodePoolName).
		Labels(rosaMachinePoolSpec.Labels).
		AutoRepair(rosaMachinePoolSpec.AutoRepair)

	if rosaMachinePoolSpec.TuningConfigs != nil {
		npBuilder = npBuilder.TuningConfigs(rosaMachinePoolSpec.TuningConfigs...)
	}

	if len(rosaMachinePoolSpec.Taints) > 0 {
		taintBuilders := []*cmv1.TaintBuilder{}
		for _, taint := range rosaMachinePoolSpec.Taints {
			newTaintBuilder := cmv1.NewTaint().Key(taint.Key).Value(taint.Value).Effect(string(taint.Effect))
			taintBuilders = append(taintBuilders, newTaintBuilder)
		}
		npBuilder = npBuilder.Taints(taintBuilders...)
	}

	if rosaMachinePoolSpec.Autoscaling != nil {
		npBuilder = npBuilder.Autoscaling(
			cmv1.NewNodePoolAutoscaling().
				MinReplica(rosaMachinePoolSpec.Autoscaling.MinReplicas).
				MaxReplica(rosaMachinePoolSpec.Autoscaling.MaxReplicas))
	} else {
		replicas := 1
		if machinePoolSpec.Replicas != nil {
			replicas = int(*machinePoolSpec.Replicas)
		}
		npBuilder = npBuilder.Replicas(replicas)
	}

	if rosaMachinePoolSpec.Subnet != "" {
		npBuilder.Subnet(rosaMachinePoolSpec.Subnet)
	}

	awsNodePool := cmv1.NewAWSNodePool().InstanceType(rosaMachinePoolSpec.InstanceType)
	if rosaMachinePoolSpec.AdditionalSecurityGroups != nil {
		awsNodePool = awsNodePool.AdditionalSecurityGroupIds(rosaMachinePoolSpec.AdditionalSecurityGroups...)
	}
	if rosaMachinePoolSpec.AdditionalTags != nil {
		awsNodePool = awsNodePool.Tags(rosaMachinePoolSpec.AdditionalTags)
	}
	if rosaMachinePoolSpec.VolumeSize > 75 {
		awsNodePool = awsNodePool.RootVolume(cmv1.NewAWSVolume().Size(rosaMachinePoolSpec.VolumeSize))
	}
	if rosaMachinePoolSpec.CapacityReservationID != "" {
		capacityReservation := cmv1.NewAWSCapacityReservation().Id(rosaMachinePoolSpec.CapacityReservationID)
		awsNodePool = awsNodePool.CapacityReservation(capacityReservation)
	}
	npBuilder.AWSNodePool(awsNodePool)

	if rosaMachinePoolSpec.Version != "" {
		npBuilder.Version(cmv1.NewVersion().ID(ocm.CreateVersionID(rosaMachinePoolSpec.Version, string(controlPlaneChannelGroup))))
	}

	if rosaMachinePoolSpec.NodeDrainGracePeriod != nil {
		valueBuilder := cmv1.NewValue().Value(rosaMachinePoolSpec.NodeDrainGracePeriod.Minutes()).Unit("minutes")
		npBuilder.NodeDrainGracePeriod(valueBuilder)
	}

	if rosaMachinePoolSpec.UpdateConfig != nil {
		configMgmtBuilder := cmv1.NewNodePoolManagementUpgrade()

		if rollingUpdate := rosaMachinePoolSpec.UpdateConfig.RollingUpdate; rollingUpdate != nil {
			if rollingUpdate.MaxSurge != nil {
				configMgmtBuilder = configMgmtBuilder.MaxSurge(rollingUpdate.MaxSurge.String())
			}
			if rollingUpdate.MaxUnavailable != nil {
				configMgmtBuilder = configMgmtBuilder.MaxUnavailable(rollingUpdate.MaxUnavailable.String())
			}
		}

		npBuilder = npBuilder.ManagementUpgrade(configMgmtBuilder)
	}

	if rosaMachinePoolSpec.ImageType == string(cmv1.ImageTypeWindows) {
		npBuilder = npBuilder.ImageType(cmv1.ImageTypeWindows)
	} else if rosaMachinePoolSpec.ImageType == string(cmv1.ImageTypeDefault) {
		npBuilder = npBuilder.ImageType(cmv1.ImageTypeDefault)
	}

	return npBuilder
}

func (r *ROSAMachinePoolReconciler) reconcileProviderIDList(ctx context.Context, machinePoolScope *scope.RosaMachinePoolScope, nodePool *cmv1.NodePool) error {
	tags := nodePool.AWSNodePool().Tags()
	if len(tags) == 0 {
		// can't identify EC2 instances belonging to this NodePool without tags.
		return nil
	}

	ec2Svc := scope.NewEC2Client(machinePoolScope, machinePoolScope, &machinePoolScope.Logger, machinePoolScope.InfraCluster())
	response, err := ec2Svc.DescribeInstances(ctx, &ec2.DescribeInstancesInput{
		Filters: buildEC2FiltersFromTags(tags),
	})
	if err != nil {
		return err
	}

	var providerIDList []string
	for _, reservation := range response.Reservations {
		for _, instance := range reservation.Instances {
			providerID := scope.GenerateProviderID(*instance.Placement.AvailabilityZone, *instance.InstanceId)
			providerIDList = append(providerIDList, providerID)
		}
	}

	machinePoolScope.RosaMachinePool.Spec.ProviderIDList = providerIDList
	return nil
}

func buildEC2FiltersFromTags(tags map[string]string) []ec2types.Filter {
	filters := make([]ec2types.Filter, len(tags)+1)
	for key, value := range tags {
		filters = append(filters, ec2types.Filter{
			Name: ptr.To(fmt.Sprintf("tag:%s", key)),
			Values: []string{
				value,
			},
		})
	}

	// only list instances that are running or just started
	filters = append(filters, ec2types.Filter{
		Name: ptr.To("instance-state-name"),
		Values: []string{
			"running", "pending",
		},
	})

	return filters
}

func rosaControlPlaneToRosaMachinePoolMapFunc(c client.Client, gvk schema.GroupVersionKind, log logger.Wrapper) handler.MapFunc {
	return func(ctx context.Context, o client.Object) []reconcile.Request {
		rosaControlPlane, ok := o.(*rosacontrolplanev1.ROSAControlPlane)
		if !ok {
			klog.Errorf("Expected a RosaControlPlane but got a %T", o)
		}

		if !rosaControlPlane.ObjectMeta.DeletionTimestamp.IsZero() {
			return nil
		}

		clusterKey, err := GetOwnerClusterKey(rosaControlPlane.ObjectMeta)
		if err != nil {
			log.Error(err, "couldn't get ROSA control plane owner ObjectKey")
			return nil
		}
		if clusterKey == nil {
			return nil
		}

		managedPoolForClusterList := clusterv1.MachinePoolList{}
		if err := c.List(
			ctx, &managedPoolForClusterList, client.InNamespace(clusterKey.Namespace), client.MatchingLabels{clusterv1.ClusterNameLabel: clusterKey.Name},
		); err != nil {
			log.Error(err, "couldn't list pools for cluster")
			return nil
		}

		mapFunc := machinePoolToInfrastructureMapFunc(gvk)

		var results []ctrl.Request
		for i := range managedPoolForClusterList.Items {
			rosaMachinePool := mapFunc(ctx, &managedPoolForClusterList.Items[i])
			results = append(results, rosaMachinePool...)
		}

		return results
	}
}
