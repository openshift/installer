package controllers

import (
	"context"
	"fmt"
	"time"

	"github.com/blang/semver"
	cmv1 "github.com/openshift-online/ocm-sdk-go/clustersmgmt/v1"
	"github.com/openshift/rosa/pkg/ocm"
	"github.com/pkg/errors"
	corev1 "k8s.io/api/core/v1"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/client-go/tools/record"
	"k8s.io/klog/v2"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/client/apiutil"
	"sigs.k8s.io/controller-runtime/pkg/controller"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
	"sigs.k8s.io/controller-runtime/pkg/handler"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"

	rosacontrolplanev1 "sigs.k8s.io/cluster-api-provider-aws/v2/controlplane/rosa/api/v1beta2"
	expinfrav1 "sigs.k8s.io/cluster-api-provider-aws/v2/exp/api/v1beta2"
	"sigs.k8s.io/cluster-api-provider-aws/v2/pkg/cloud/scope"
	"sigs.k8s.io/cluster-api-provider-aws/v2/pkg/logger"
	"sigs.k8s.io/cluster-api-provider-aws/v2/pkg/rosa"
	clusterv1 "sigs.k8s.io/cluster-api/api/v1beta1"
	expclusterv1 "sigs.k8s.io/cluster-api/exp/api/v1beta1"
	"sigs.k8s.io/cluster-api/util"
	"sigs.k8s.io/cluster-api/util/annotations"
	"sigs.k8s.io/cluster-api/util/conditions"
	"sigs.k8s.io/cluster-api/util/predicates"
)

// ROSAMachinePoolReconciler reconciles a RosaMachinePool object.
type ROSAMachinePoolReconciler struct {
	client.Client
	Recorder         record.EventRecorder
	WatchFilterValue string
	Endpoints        []scope.ServiceEndpoint
}

// SetupWithManager is used to setup the controller.
func (r *ROSAMachinePoolReconciler) SetupWithManager(ctx context.Context, mgr ctrl.Manager, options controller.Options) error {
	log := logger.FromContext(ctx)

	gvk, err := apiutil.GVKForObject(new(expinfrav1.ROSAMachinePool), mgr.GetScheme())
	if err != nil {
		return errors.Wrapf(err, "failed to find GVK for RosaMachinePool")
	}
	rosaControlPlaneToRosaMachinePoolMap := rosaControlPlaneToRosaMachinePoolMapFunc(r.Client, gvk, log)
	return ctrl.NewControllerManagedBy(mgr).
		For(&expinfrav1.ROSAMachinePool{}).
		WithOptions(options).
		WithEventFilter(predicates.ResourceNotPausedAndHasFilterLabel(log.GetLogger(), r.WatchFilterValue)).
		Watches(
			&expclusterv1.MachinePool{},
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

// Reconcile reconciles RosaMachinePool.
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
		return reconcile.Result{}, nil
	}

	if annotations.IsPaused(cluster, rosaMachinePool) {
		log.Info("Reconciliation is paused for this object")
		return ctrl.Result{}, nil
	}

	log = log.WithValues("cluster", klog.KObj(cluster))

	controlPlaneKey := client.ObjectKey{
		Namespace: rosaMachinePool.Namespace,
		Name:      cluster.Spec.ControlPlaneRef.Name,
	}
	controlPlane := &rosacontrolplanev1.ROSAControlPlane{}
	if err := r.Client.Get(ctx, controlPlaneKey, controlPlane); err != nil {
		log.Info("Failed to retrieve ControlPlane from MachinePool")
		return reconcile.Result{}, nil
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
		return ctrl.Result{}, errors.Wrap(err, "failed to create scope")
	}

	rosaControlPlaneScope, err := scope.NewROSAControlPlaneScope(scope.ROSAControlPlaneScopeParams{
		Client:         r.Client,
		Cluster:        cluster,
		ControlPlane:   controlPlane,
		ControllerName: "rosaControlPlane",
	})
	if err != nil {
		return ctrl.Result{}, errors.Wrap(err, "failed to create control plane scope")
	}

	if !controlPlane.Status.Ready {
		log.Info("Control plane is not ready yet")
		err := machinePoolScope.RosaMchinePoolReadyFalse(expinfrav1.WaitingForRosaControlPlaneReason, "")
		return ctrl.Result{}, err
	}

	defer func() {
		conditions.SetSummary(machinePoolScope.RosaMachinePool, conditions.WithConditions(expinfrav1.RosaMachinePoolReadyCondition), conditions.WithStepCounter())

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
	machinePoolScope.Info("Reconciling RosaMachinePool")

	if controllerutil.AddFinalizer(machinePoolScope.RosaMachinePool, expinfrav1.RosaMachinePoolFinalizer) {
		if err := machinePoolScope.PatchObject(); err != nil {
			return ctrl.Result{}, err
		}
	}

	ocmClient, err := rosa.NewOCMClient(ctx, rosaControlPlaneScope)
	if err != nil {
		// TODO: need to expose in status, as likely the credentials are invalid
		return ctrl.Result{}, fmt.Errorf("failed to create OCM client: %w", err)
	}

	failureMessage, err := validateMachinePoolSpec(machinePoolScope)
	if err != nil {
		return ctrl.Result{}, fmt.Errorf("failed to validate ROSAMachinePool.spec: %w", err)
	}
	if failureMessage != nil {
		machinePoolScope.RosaMachinePool.Status.FailureMessage = failureMessage
		// dont' requeue because input is invalid and manual intervention is needed.
		return ctrl.Result{}, nil
	} else {
		machinePoolScope.RosaMachinePool.Status.FailureMessage = nil
	}

	rosaMachinePool := machinePoolScope.RosaMachinePool
	machinePool := machinePoolScope.MachinePool
	controlPlane := machinePoolScope.ControlPlane

	createdNodePool, found, err := ocmClient.GetNodePool(*controlPlane.Status.ID, rosaMachinePool.Spec.NodePoolName)
	if err != nil {
		return ctrl.Result{}, err
	}
	if found {
		// TODO (alberto): discover and store providerIDs from aws so the CAPI controller can match then to Nodes and report readiness.
		rosaMachinePool.Status.Replicas = int32(createdNodePool.Status().CurrentReplicas())
		if createdNodePool.Replicas() == createdNodePool.Status().CurrentReplicas() && createdNodePool.Status().Message() == "" {
			conditions.MarkTrue(rosaMachinePool, expinfrav1.RosaMachinePoolReadyCondition)
			rosaMachinePool.Status.Ready = true

			if err := r.reconcileMachinePoolVersion(machinePoolScope, ocmClient, createdNodePool); err != nil {
				return ctrl.Result{}, err
			}

			return ctrl.Result{}, nil
		}

		conditions.MarkFalse(rosaMachinePool,
			expinfrav1.RosaMachinePoolReadyCondition,
			createdNodePool.Status().Message(),
			clusterv1.ConditionSeverityInfo,
			"")

		machinePoolScope.Info("waiting for NodePool to become ready", "state", createdNodePool.Status().Message())
		// Requeue so that status.ready is set to true when the nodepool is fully created.
		return ctrl.Result{RequeueAfter: time.Second * 60}, nil
	}

	npBuilder := cmv1.NewNodePool()
	npBuilder.ID(rosaMachinePool.Spec.NodePoolName).
		Labels(rosaMachinePool.Spec.Labels).
		AutoRepair(rosaMachinePool.Spec.AutoRepair)

	if rosaMachinePool.Spec.Autoscaling != nil {
		npBuilder = npBuilder.Autoscaling(
			cmv1.NewNodePoolAutoscaling().
				MinReplica(rosaMachinePool.Spec.Autoscaling.MinReplicas).
				MaxReplica(rosaMachinePool.Spec.Autoscaling.MaxReplicas))
	} else {
		replicas := 1
		if machinePool.Spec.Replicas != nil {
			replicas = int(*machinePool.Spec.Replicas)
		}
		npBuilder = npBuilder.Replicas(replicas)
	}

	if rosaMachinePool.Spec.Subnet != "" {
		npBuilder.Subnet(rosaMachinePool.Spec.Subnet)
	}

	npBuilder.AWSNodePool(cmv1.NewAWSNodePool().InstanceType(rosaMachinePool.Spec.InstanceType))
	if rosaMachinePool.Spec.Version != "" {
		npBuilder.Version(cmv1.NewVersion().ID(ocm.CreateVersionID(rosaMachinePool.Spec.Version, ocm.DefaultChannelGroup)))
	}

	nodePoolSpec, err := npBuilder.Build()
	if err != nil {
		return ctrl.Result{}, fmt.Errorf("failed to build rosa nodepool: %w", err)
	}

	createdNodePool, err = ocmClient.CreateNodePool(*controlPlane.Status.ID, nodePoolSpec)
	if err != nil {
		return ctrl.Result{}, fmt.Errorf("failed to create nodepool: %w", err)
	}

	machinePoolScope.RosaMachinePool.Status.ID = createdNodePool.ID()

	return ctrl.Result{}, nil
}

func (r *ROSAMachinePoolReconciler) reconcileDelete(
	ctx context.Context, machinePoolScope *scope.RosaMachinePoolScope,
	rosaControlPlaneScope *scope.ROSAControlPlaneScope,
) error {
	machinePoolScope.Info("Reconciling deletion of RosaMachinePool")

	ocmClient, err := rosa.NewOCMClient(ctx, rosaControlPlaneScope)
	if err != nil {
		// TODO: need to expose in status, as likely the credentials are invalid
		return fmt.Errorf("failed to create OCM client: %w", err)
	}

	nodePool, found, err := ocmClient.GetNodePool(*machinePoolScope.ControlPlane.Status.ID, machinePoolScope.NodePoolName())
	if err != nil {
		return err
	}
	if found {
		if err := ocmClient.DeleteNodePool(*machinePoolScope.ControlPlane.Status.ID, nodePool.ID()); err != nil {
			return err
		}
	}

	controllerutil.RemoveFinalizer(machinePoolScope.RosaMachinePool, expinfrav1.RosaMachinePoolFinalizer)

	return nil
}

func (r *ROSAMachinePoolReconciler) reconcileMachinePoolVersion(machinePoolScope *scope.RosaMachinePoolScope, ocmClient *ocm.Client, nodePool *cmv1.NodePool) error {
	version := machinePoolScope.RosaMachinePool.Spec.Version
	if version == "" {
		version = machinePoolScope.ControlPlane.Spec.Version
	}

	if version == rosa.RawVersionID(nodePool.Version()) {
		conditions.MarkFalse(machinePoolScope.RosaMachinePool, expinfrav1.RosaMachinePoolUpgradingCondition, "upgraded", clusterv1.ConditionSeverityInfo, "")
		return nil
	}

	clusterID := *machinePoolScope.ControlPlane.Status.ID
	_, scheduledUpgrade, err := ocmClient.GetHypershiftNodePoolUpgrade(clusterID, machinePoolScope.ControlPlane.Spec.RosaClusterName, nodePool.ID())
	if err != nil {
		return fmt.Errorf("failed to get existing scheduled upgrades: %w", err)
	}

	if scheduledUpgrade == nil {
		policy, err := ocmClient.BuildNodeUpgradePolicy(version, nodePool.ID(), ocm.UpgradeScheduling{
			AutomaticUpgrades: false,
			// The OCM API places guardrails around the minimum and maximum delay that a user can request,
			// for the next run of the upgrade, which is [5min,6mo]. Set our next run request to something
			// slightly longer than 5min to make sure we account for the latency between when we send this
			// request and when the server processes it.
			// https://gitlab.cee.redhat.com/service/uhc-clusters-service/-/blob/master/cmd/clusters-service/servecmd/apiserver/upgrade_policy_handlers.go
			NextRun: time.Now().Add(6 * time.Minute),
		})
		if err != nil {
			return fmt.Errorf("failed to create nodePool upgrade schedule to version %s: %w", version, err)
		}

		scheduledUpgrade, err = ocmClient.ScheduleNodePoolUpgrade(clusterID, nodePool.ID(), policy)
		if err != nil {
			return fmt.Errorf("failed to schedule nodePool upgrade to version %s: %w", version, err)
		}
	}

	condition := &clusterv1.Condition{
		Type:    expinfrav1.RosaMachinePoolUpgradingCondition,
		Status:  corev1.ConditionTrue,
		Reason:  string(scheduledUpgrade.State().Value()),
		Message: fmt.Sprintf("Upgrading to version %s", scheduledUpgrade.Version()),
	}
	conditions.Set(machinePoolScope.RosaMachinePool, condition)

	// if nodePool is already upgrading to another version we need to wait until the current upgrade is finished, return an error to requeue and try later.
	if scheduledUpgrade.Version() != version {
		return fmt.Errorf("there is already a %s upgrade to version %s", scheduledUpgrade.State().Value(), scheduledUpgrade.Version())
	}

	return nil
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

		managedPoolForClusterList := expclusterv1.MachinePoolList{}
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
