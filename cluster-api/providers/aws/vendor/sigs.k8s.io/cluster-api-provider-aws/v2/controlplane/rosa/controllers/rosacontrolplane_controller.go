/*
Copyright 2023 The Kubernetes Authors.

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
	"errors"
	"fmt"
	"net"
	"net/url"
	"strconv"
	"strings"
	"time"

	idputils "github.com/openshift-online/ocm-common/pkg/idp/utils"
	cmv1 "github.com/openshift-online/ocm-sdk-go/clustersmgmt/v1"
	rosaaws "github.com/openshift/rosa/pkg/aws"
	"github.com/openshift/rosa/pkg/ocm"
	"github.com/zgalor/weberr"
	corev1 "k8s.io/api/core/v1"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	restclient "k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/tools/clientcmd/api"
	"k8s.io/klog/v2"
	"k8s.io/utils/ptr"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
	"sigs.k8s.io/controller-runtime/pkg/handler"
	"sigs.k8s.io/controller-runtime/pkg/source"

	rosacontrolplanev1 "sigs.k8s.io/cluster-api-provider-aws/v2/controlplane/rosa/api/v1beta2"
	expinfrav1 "sigs.k8s.io/cluster-api-provider-aws/v2/exp/api/v1beta2"
	"sigs.k8s.io/cluster-api-provider-aws/v2/pkg/cloud/scope"
	"sigs.k8s.io/cluster-api-provider-aws/v2/pkg/logger"
	"sigs.k8s.io/cluster-api-provider-aws/v2/pkg/rosa"
	clusterv1 "sigs.k8s.io/cluster-api/api/v1beta1"
	"sigs.k8s.io/cluster-api/util"
	capiannotations "sigs.k8s.io/cluster-api/util/annotations"
	"sigs.k8s.io/cluster-api/util/conditions"
	"sigs.k8s.io/cluster-api/util/kubeconfig"
	"sigs.k8s.io/cluster-api/util/predicates"
	"sigs.k8s.io/cluster-api/util/secret"
)

const (
	rosaControlPlaneKind = "ROSAControlPlane"
	// ROSAControlPlaneFinalizer allows the controller to clean up resources on delete.
	ROSAControlPlaneFinalizer = "rosacontrolplane.controlplane.cluster.x-k8s.io"
)

type ROSAControlPlaneReconciler struct {
	client.Client
	WatchFilterValue string
	WaitInfraPeriod  time.Duration
	Endpoints        []scope.ServiceEndpoint
}

// SetupWithManager is used to setup the controller.
func (r *ROSAControlPlaneReconciler) SetupWithManager(ctx context.Context, mgr ctrl.Manager, options controller.Options) error {
	log := logger.FromContext(ctx)

	rosaControlPlane := &rosacontrolplanev1.ROSAControlPlane{}
	c, err := ctrl.NewControllerManagedBy(mgr).
		For(rosaControlPlane).
		WithOptions(options).
		WithEventFilter(predicates.ResourceNotPausedAndHasFilterLabel(log.GetLogger(), r.WatchFilterValue)).
		Build(r)

	if err != nil {
		return fmt.Errorf("failed setting up the AWSManagedControlPlane controller manager: %w", err)
	}

	if err = c.Watch(
		source.Kind(mgr.GetCache(), &clusterv1.Cluster{}),
		handler.EnqueueRequestsFromMapFunc(util.ClusterToInfrastructureMapFunc(ctx, rosaControlPlane.GroupVersionKind(), mgr.GetClient(), &expinfrav1.ROSACluster{})),
		predicates.ClusterUnpausedAndInfrastructureReady(log.GetLogger()),
	); err != nil {
		return fmt.Errorf("failed adding a watch for ready clusters: %w", err)
	}

	if err = c.Watch(
		source.Kind(mgr.GetCache(), &expinfrav1.ROSACluster{}),
		handler.EnqueueRequestsFromMapFunc(r.rosaClusterToROSAControlPlane(log)),
	); err != nil {
		return fmt.Errorf("failed adding a watch for ROSACluster")
	}

	return nil
}

// +kubebuilder:rbac:groups=core,resources=events,verbs=get;list;watch;create;patch
// +kubebuilder:rbac:groups="",resources=secrets,verbs=get;list;watch;create;update;delete;patch
// +kubebuilder:rbac:groups="",resources=namespaces,verbs=get;list;watch
// +kubebuilder:rbac:groups=cluster.x-k8s.io,resources=clusters;clusters/status,verbs=get;list;watch
// +kubebuilder:rbac:groups=cluster.x-k8s.io,resources=machinedeployments,verbs=get;list;watch
// +kubebuilder:rbac:groups=cluster.x-k8s.io,resources=machinepools,verbs=get;list;watch
// +kubebuilder:rbac:groups=controlplane.cluster.x-k8s.io,resources=rosacontrolplanes,verbs=get;list;watch;update;patch;delete
// +kubebuilder:rbac:groups=controlplane.cluster.x-k8s.io,resources=rosacontrolplanes/status,verbs=get;update;patch

// Reconcile will reconcile RosaControlPlane Resources.
func (r *ROSAControlPlaneReconciler) Reconcile(ctx context.Context, req ctrl.Request) (res ctrl.Result, reterr error) {
	log := logger.FromContext(ctx)

	// Get the control plane instance
	rosaControlPlane := &rosacontrolplanev1.ROSAControlPlane{}
	if err := r.Client.Get(ctx, req.NamespacedName, rosaControlPlane); err != nil {
		if apierrors.IsNotFound(err) {
			return ctrl.Result{}, nil
		}
		return ctrl.Result{Requeue: true}, nil
	}

	// Get the cluster
	cluster, err := util.GetOwnerCluster(ctx, r.Client, rosaControlPlane.ObjectMeta)
	if err != nil {
		log.Error(err, "Failed to retrieve owner Cluster from the API Server")
		return ctrl.Result{}, err
	}
	if cluster == nil {
		log.Info("Cluster Controller has not yet set OwnerRef")
		return ctrl.Result{}, nil
	}

	log = log.WithValues("cluster", klog.KObj(cluster))

	if capiannotations.IsPaused(cluster, rosaControlPlane) {
		log.Info("Reconciliation is paused for this object")
		return ctrl.Result{}, nil
	}

	rosaScope, err := scope.NewROSAControlPlaneScope(scope.ROSAControlPlaneScopeParams{
		Client:         r.Client,
		Cluster:        cluster,
		ControlPlane:   rosaControlPlane,
		ControllerName: strings.ToLower(rosaControlPlaneKind),
		Endpoints:      r.Endpoints,
		Logger:         log,
	})
	if err != nil {
		return ctrl.Result{}, fmt.Errorf("failed to create scope: %w", err)
	}

	// Always close the scope
	defer func() {
		if err := rosaScope.Close(); err != nil {
			reterr = errors.Join(reterr, err)
		}
	}()

	if !rosaControlPlane.ObjectMeta.DeletionTimestamp.IsZero() {
		// Handle deletion reconciliation loop.
		return r.reconcileDelete(ctx, rosaScope)
	}

	// Handle normal reconciliation loop.
	return r.reconcileNormal(ctx, rosaScope)
}

func (r *ROSAControlPlaneReconciler) reconcileNormal(ctx context.Context, rosaScope *scope.ROSAControlPlaneScope) (res ctrl.Result, reterr error) {
	rosaScope.Info("Reconciling ROSAControlPlane")

	// if !rosaScope.Cluster.Status.InfrastructureReady {
	//	rosaScope.Info("Cluster infrastructure is not ready yet")
	//	return ctrl.Result{RequeueAfter: r.WaitInfraPeriod}, nil
	//}
	if controllerutil.AddFinalizer(rosaScope.ControlPlane, ROSAControlPlaneFinalizer) {
		if err := rosaScope.PatchObject(); err != nil {
			return ctrl.Result{}, err
		}
	}

	ocmClient, err := rosa.NewOCMClient(ctx, rosaScope)
	if err != nil {
		// TODO: need to expose in status, as likely the credentials are invalid
		return ctrl.Result{}, fmt.Errorf("failed to create OCM client: %w", err)
	}

	creator, err := rosaaws.CreatorForCallerIdentity(rosaScope.Identity)
	if err != nil {
		return ctrl.Result{}, fmt.Errorf("failed to transform caller identity to creator: %w", err)
	}

	if validationMessage, validationError := validateControlPlaneSpec(ocmClient, rosaScope); validationError != nil {
		return ctrl.Result{}, fmt.Errorf("validate ROSAControlPlane.spec: %w", validationError)
	} else if validationMessage != "" {
		rosaScope.ControlPlane.Status.FailureMessage = ptr.To(validationMessage)
		// dont' requeue because input is invalid and manual intervention is needed.
		return ctrl.Result{}, nil
	} else {
		rosaScope.ControlPlane.Status.FailureMessage = nil
	}

	cluster, err := ocmClient.GetCluster(rosaScope.ControlPlane.Spec.RosaClusterName, creator)
	if weberr.GetType(err) != weberr.NotFound && err != nil {
		return ctrl.Result{}, err
	}

	if cluster != nil {
		rosaScope.ControlPlane.Status.ID = ptr.To(cluster.ID())
		rosaScope.ControlPlane.Status.ConsoleURL = cluster.Console().URL()
		rosaScope.ControlPlane.Status.OIDCEndpointURL = cluster.AWS().STS().OIDCEndpointURL()
		rosaScope.ControlPlane.Status.Ready = false

		switch cluster.Status().State() {
		case cmv1.ClusterStateReady:
			conditions.MarkTrue(rosaScope.ControlPlane, rosacontrolplanev1.ROSAControlPlaneReadyCondition)
			rosaScope.ControlPlane.Status.Ready = true

			apiEndpoint, err := buildAPIEndpoint(cluster)
			if err != nil {
				return ctrl.Result{}, err
			}
			rosaScope.ControlPlane.Spec.ControlPlaneEndpoint = *apiEndpoint

			if err := r.reconcileKubeconfig(ctx, rosaScope, ocmClient, cluster); err != nil {
				return ctrl.Result{}, fmt.Errorf("failed to reconcile kubeconfig: %w", err)
			}
			if err := r.reconcileClusterVersion(rosaScope, ocmClient, cluster); err != nil {
				return ctrl.Result{}, err
			}
			return ctrl.Result{}, nil
		case cmv1.ClusterStateError:
			errorMessage := cluster.Status().ProvisionErrorMessage()
			rosaScope.ControlPlane.Status.FailureMessage = &errorMessage

			conditions.MarkFalse(rosaScope.ControlPlane,
				rosacontrolplanev1.ROSAControlPlaneReadyCondition,
				string(cluster.Status().State()),
				clusterv1.ConditionSeverityError,
				cluster.Status().ProvisionErrorCode())
			// Cluster is in an unrecoverable state, returning nil error so that the request doesn't get requeued.
			return ctrl.Result{}, nil
		}

		conditions.MarkFalse(rosaScope.ControlPlane,
			rosacontrolplanev1.ROSAControlPlaneReadyCondition,
			string(cluster.Status().State()),
			clusterv1.ConditionSeverityInfo,
			cluster.Status().Description())

		rosaScope.Info("waiting for cluster to become ready", "state", cluster.Status().State())
		// Requeue so that status.ready is set to true when the cluster is fully created.
		return ctrl.Result{RequeueAfter: time.Second * 60}, nil
	}

	billingAccount := *rosaScope.Identity.Account
	if rosaScope.ControlPlane.Spec.BillingAccount != "" {
		billingAccount = rosaScope.ControlPlane.Spec.BillingAccount
	}

	spec := ocm.Spec{
		DryRun:                    ptr.To(false),
		Name:                      rosaScope.RosaClusterName(),
		Region:                    *rosaScope.ControlPlane.Spec.Region,
		MultiAZ:                   true,
		Version:                   ocm.CreateVersionID(rosaScope.ControlPlane.Spec.Version, ocm.DefaultChannelGroup),
		ChannelGroup:              "stable",
		Expiration:                time.Now().Add(1 * time.Hour),
		DisableWorkloadMonitoring: ptr.To(true),
		DefaultIngress:            ocm.NewDefaultIngressSpec(), // n.b. this is a no-op when it's set to the default value
		ComputeMachineType:        rosaScope.ControlPlane.Spec.InstanceType,

		SubnetIds:         rosaScope.ControlPlane.Spec.Subnets,
		AvailabilityZones: rosaScope.ControlPlane.Spec.AvailabilityZones,
		NetworkType:       rosaScope.ControlPlane.Spec.Network.NetworkType,
		IsSTS:             true,
		RoleARN:           *rosaScope.ControlPlane.Spec.InstallerRoleARN,
		SupportRoleARN:    *rosaScope.ControlPlane.Spec.SupportRoleARN,
		OperatorIAMRoles:  getOperatorIAMRole(*rosaScope.ControlPlane),
		WorkerRoleARN:     *rosaScope.ControlPlane.Spec.WorkerRoleARN,
		OidcConfigId:      *rosaScope.ControlPlane.Spec.OIDCID,
		Mode:              "auto",
		Hypershift: ocm.Hypershift{
			Enabled: true,
		},
		BillingAccount: billingAccount,
		AWSCreator:     creator,
	}

	_, machineCIDR, err := net.ParseCIDR(rosaScope.ControlPlane.Spec.Network.MachineCIDR)
	if err == nil {
		spec.MachineCIDR = *machineCIDR
	} else {
		// TODO: expose in status
		rosaScope.Error(err, "rosacontrolplane.spec.network.machineCIDR invalid", rosaScope.ControlPlane.Spec.Network.MachineCIDR)
		return ctrl.Result{}, nil
	}

	if rosaScope.ControlPlane.Spec.Network.PodCIDR != "" {
		_, podCIDR, err := net.ParseCIDR(rosaScope.ControlPlane.Spec.Network.PodCIDR)
		if err == nil {
			spec.PodCIDR = *podCIDR
		} else {
			// TODO: expose in status.
			rosaScope.Error(err, "rosacontrolplane.spec.network.podCIDR invalid", rosaScope.ControlPlane.Spec.Network.PodCIDR)
			return ctrl.Result{}, nil
		}
	}

	if rosaScope.ControlPlane.Spec.Network.ServiceCIDR != "" {
		_, serviceCIDR, err := net.ParseCIDR(rosaScope.ControlPlane.Spec.Network.ServiceCIDR)
		if err == nil {
			spec.ServiceCIDR = *serviceCIDR
		} else {
			// TODO: expose in status.
			rosaScope.Error(err, "rosacontrolplane.spec.network.serviceCIDR invalid", rosaScope.ControlPlane.Spec.Network.ServiceCIDR)
			return ctrl.Result{}, nil
		}
	}

	// Set autoscale replica
	if rosaScope.ControlPlane.Spec.Autoscaling != nil {
		spec.MaxReplicas = rosaScope.ControlPlane.Spec.Autoscaling.MaxReplicas
		spec.MinReplicas = rosaScope.ControlPlane.Spec.Autoscaling.MinReplicas
	}

	cluster, err = ocmClient.CreateCluster(spec)
	if err != nil {
		// TODO: need to expose in status, as likely the spec is invalid
		return ctrl.Result{}, fmt.Errorf("failed to create OCM cluster: %w", err)
	}

	rosaScope.Info("cluster created", "state", cluster.Status().State())
	clusterID := cluster.ID()
	rosaScope.ControlPlane.Status.ID = &clusterID

	return ctrl.Result{}, nil
}

func getOperatorIAMRole(rosaControlPlane rosacontrolplanev1.ROSAControlPlane) []ocm.OperatorIAMRole {
	return []ocm.OperatorIAMRole{
		{
			Name:      "cloud-credentials",
			Namespace: "openshift-ingress-operator",
			RoleARN:   rosaControlPlane.Spec.RolesRef.IngressARN,
		},
		{
			Name:      "installer-cloud-credentials",
			Namespace: "openshift-image-registry",
			RoleARN:   rosaControlPlane.Spec.RolesRef.ImageRegistryARN,
		},
		{
			Name:      "ebs-cloud-credentials",
			Namespace: "openshift-cluster-csi-drivers",
			RoleARN:   rosaControlPlane.Spec.RolesRef.StorageARN,
		},
		{
			Name:      "cloud-credentials",
			Namespace: "openshift-cloud-network-config-controller",
			RoleARN:   rosaControlPlane.Spec.RolesRef.NetworkARN,
		},
		{
			Name:      "kube-controller-manager",
			Namespace: "kube-system",
			RoleARN:   rosaControlPlane.Spec.RolesRef.KubeCloudControllerARN,
		},
		{
			Name:      "kms-provider",
			Namespace: "kube-system",
			RoleARN:   rosaControlPlane.Spec.RolesRef.KMSProviderARN,
		},
		{
			Name:      "control-plane-operator",
			Namespace: "kube-system",
			RoleARN:   rosaControlPlane.Spec.RolesRef.ControlPlaneOperatorARN,
		},
		{
			Name:      "capa-controller-manager",
			Namespace: "kube-system",
			RoleARN:   rosaControlPlane.Spec.RolesRef.NodePoolManagementARN,
		},
	}
}

func (r *ROSAControlPlaneReconciler) reconcileDelete(ctx context.Context, rosaScope *scope.ROSAControlPlaneScope) (res ctrl.Result, reterr error) {
	rosaScope.Info("Reconciling ROSAControlPlane delete")

	ocmClient, err := rosa.NewOCMClient(ctx, rosaScope)
	if err != nil {
		// TODO: need to expose in status, as likely the credentials are invalid
		return ctrl.Result{}, fmt.Errorf("failed to create OCM client: %w", err)
	}

	creator, err := rosaaws.CreatorForCallerIdentity(rosaScope.Identity)
	if err != nil {
		return ctrl.Result{}, fmt.Errorf("failed to transform caller identity to creator: %w", err)
	}

	cluster, err := ocmClient.GetCluster(rosaScope.ControlPlane.Spec.RosaClusterName, creator)
	if weberr.GetType(err) != weberr.NotFound && err != nil {
		return ctrl.Result{}, err
	}
	if cluster == nil {
		// cluster is fully deleted, remove finalizer.
		controllerutil.RemoveFinalizer(rosaScope.ControlPlane, ROSAControlPlaneFinalizer)
		return ctrl.Result{}, nil
	}

	if cluster.Status().State() != cmv1.ClusterStateUninstalling {
		if _, err := ocmClient.DeleteCluster(cluster.ID(), true, creator); err != nil {
			return ctrl.Result{}, err
		}
	}

	rosaScope.ControlPlane.Status.Ready = false
	rosaScope.Info("waiting for cluster to be deleted")
	// Requeue to remove the finalizer when the cluster is fully deleted.
	return ctrl.Result{RequeueAfter: time.Second * 60}, nil
}

func (r *ROSAControlPlaneReconciler) reconcileClusterVersion(rosaScope *scope.ROSAControlPlaneScope, ocmClient *ocm.Client, cluster *cmv1.Cluster) error {
	version := rosaScope.ControlPlane.Spec.Version
	if version == rosa.RawVersionID(cluster.Version()) {
		conditions.MarkFalse(rosaScope.ControlPlane, rosacontrolplanev1.ROSAControlPlaneUpgradingCondition, "upgraded", clusterv1.ConditionSeverityInfo, "")
		return nil
	}

	scheduledUpgrade, err := rosa.CheckExistingScheduledUpgrade(ocmClient, cluster)
	if err != nil {
		return fmt.Errorf("failed to get existing scheduled upgrades: %w", err)
	}

	if scheduledUpgrade == nil {
		scheduledUpgrade, err = rosa.ScheduleControlPlaneUpgrade(ocmClient, cluster, version, time.Now())
		if err != nil {
			return fmt.Errorf("failed to schedule control plane upgrade to version %s: %w", version, err)
		}
	}

	condition := &clusterv1.Condition{
		Type:    rosacontrolplanev1.ROSAControlPlaneUpgradingCondition,
		Status:  corev1.ConditionTrue,
		Reason:  string(scheduledUpgrade.State().Value()),
		Message: fmt.Sprintf("Upgrading to version %s", scheduledUpgrade.Version()),
	}
	conditions.Set(rosaScope.ControlPlane, condition)

	// if cluster is already upgrading to another version we need to wait until the current upgrade is finished, return an error to requeue and try later.
	if scheduledUpgrade.Version() != version {
		return fmt.Errorf("there is already a %s upgrade to version %s", scheduledUpgrade.State().Value(), scheduledUpgrade.Version())
	}

	return nil
}

func (r *ROSAControlPlaneReconciler) reconcileKubeconfig(ctx context.Context, rosaScope *scope.ROSAControlPlaneScope, ocmClient *ocm.Client, cluster *cmv1.Cluster) error {
	rosaScope.Debug("Reconciling ROSA kubeconfig for cluster", "cluster-name", rosaScope.RosaClusterName())

	clusterRef := client.ObjectKeyFromObject(rosaScope.Cluster)
	kubeconfigSecret, err := secret.GetFromNamespacedName(ctx, r.Client, clusterRef, secret.Kubeconfig)
	if err != nil {
		if !apierrors.IsNotFound(err) {
			return fmt.Errorf("failed to get kubeconfig secret: %w", err)
		}
	}

	// generate a new password for the cluster admin user, or retrieve an existing one.
	password, err := r.reconcileClusterAdminPassword(ctx, rosaScope)
	if err != nil {
		return fmt.Errorf("failed to reconcile cluster admin password secret: %w", err)
	}

	clusterName := rosaScope.RosaClusterName()
	userName := fmt.Sprintf("%s-capi-admin", clusterName)
	apiServerURL := cluster.API().URL()

	// create new user with admin privileges in the ROSA cluster if 'userName' doesn't already exist.
	err = rosa.CreateAdminUserIfNotExist(ocmClient, cluster.ID(), userName, password)
	if err != nil {
		return err
	}

	clientConfig := &restclient.Config{
		Host:     apiServerURL,
		Username: userName,
	}
	// request an acccess token using the credentials of the cluster admin user created earlier.
	// this token is used in the kubeconfig to authenticate with the API server.
	token, err := rosa.RequestToken(ctx, apiServerURL, userName, password, clientConfig)
	if err != nil {
		return fmt.Errorf("failed to request token: %w", err)
	}

	// create the kubeconfig spec.
	contextName := fmt.Sprintf("%s@%s", userName, clusterName)
	cfg := &api.Config{
		APIVersion: api.SchemeGroupVersion.Version,
		Clusters: map[string]*api.Cluster{
			clusterName: {
				Server: apiServerURL,
			},
		},
		Contexts: map[string]*api.Context{
			contextName: {
				Cluster:  clusterName,
				AuthInfo: userName,
			},
		},
		CurrentContext: contextName,
		AuthInfos: map[string]*api.AuthInfo{
			userName: {
				Token: token.AccessToken,
			},
		},
	}
	out, err := clientcmd.Write(*cfg)
	if err != nil {
		return fmt.Errorf("failed to serialize config to yaml: %w", err)
	}

	if kubeconfigSecret != nil {
		// update existing kubeconfig secret.
		kubeconfigSecret.Data[secret.KubeconfigDataName] = out
		if err := r.Client.Update(ctx, kubeconfigSecret); err != nil {
			return fmt.Errorf("failed to update kubeconfig secret: %w", err)
		}
	} else {
		// create new kubeconfig secret.
		controllerOwnerRef := *metav1.NewControllerRef(rosaScope.ControlPlane, rosacontrolplanev1.GroupVersion.WithKind("ROSAControlPlane"))
		kubeconfigSecret = kubeconfig.GenerateSecretWithOwner(clusterRef, out, controllerOwnerRef)
		if err := r.Client.Create(ctx, kubeconfigSecret); err != nil {
			return fmt.Errorf("failed to create kubeconfig secret: %w", err)
		}
	}

	rosaScope.ControlPlane.Status.Initialized = true
	return nil
}

// reconcileClusterAdminPassword generates and store the password of the cluster admin user in a secret which is used to request a token for kubeconfig auth.
// Since it is not possible to retrieve a user's password through the ocm API once created,
// we have to store the password in a secret as it is needed later to refresh the token.
func (r *ROSAControlPlaneReconciler) reconcileClusterAdminPassword(ctx context.Context, rosaScope *scope.ROSAControlPlaneScope) (string, error) {
	passwordSecret := rosaScope.ClusterAdminPasswordSecret()
	err := r.Client.Get(ctx, client.ObjectKeyFromObject(passwordSecret), passwordSecret)
	if err == nil {
		password := string(passwordSecret.Data["value"])
		return password, nil
	} else if !apierrors.IsNotFound(err) {
		return "", fmt.Errorf("failed to get cluster admin password secret: %w", err)
	}
	password, err := idputils.GenerateRandomPassword()
	// Generate a new password and create the secret
	if err != nil {
		return "", err
	}

	controllerOwnerRef := *metav1.NewControllerRef(rosaScope.ControlPlane, rosacontrolplanev1.GroupVersion.WithKind("ROSAControlPlane"))
	passwordSecret.Data = map[string][]byte{
		"value": []byte(password),
	}
	passwordSecret.OwnerReferences = []metav1.OwnerReference{
		controllerOwnerRef,
	}
	if err := r.Client.Create(ctx, passwordSecret); err != nil {
		return "", err
	}

	return password, nil
}

func validateControlPlaneSpec(ocmClient *ocm.Client, rosaScope *scope.ROSAControlPlaneScope) (string, error) {
	version := rosaScope.ControlPlane.Spec.Version
	valid, err := ocmClient.ValidateHypershiftVersion(version, "stable")
	if err != nil {
		return "", fmt.Errorf("failed to check if version is valid: %w", err)
	}
	if !valid {
		return fmt.Sprintf("version %s is not supported", version), nil
	}

	// TODO: add more input validations
	return "", nil
}

func (r *ROSAControlPlaneReconciler) rosaClusterToROSAControlPlane(log *logger.Logger) handler.MapFunc {
	return func(ctx context.Context, o client.Object) []ctrl.Request {
		rosaCluster, ok := o.(*expinfrav1.ROSACluster)
		if !ok {
			log.Error(fmt.Errorf("expected a ROSACluster but got a %T", o), "Expected ROSACluster")
			return nil
		}

		if !rosaCluster.ObjectMeta.DeletionTimestamp.IsZero() {
			log.Debug("ROSACluster has a deletion timestamp, skipping mapping")
			return nil
		}

		cluster, err := util.GetOwnerCluster(ctx, r.Client, rosaCluster.ObjectMeta)
		if err != nil {
			log.Error(err, "failed to get owning cluster")
			return nil
		}
		if cluster == nil {
			log.Debug("Owning cluster not set on ROSACluster, skipping mapping")
			return nil
		}

		controlPlaneRef := cluster.Spec.ControlPlaneRef
		if controlPlaneRef == nil || controlPlaneRef.Kind != rosaControlPlaneKind {
			log.Debug("ControlPlaneRef is nil or not ROSAControlPlane, skipping mapping")
			return nil
		}

		return []ctrl.Request{
			{
				NamespacedName: types.NamespacedName{
					Name:      controlPlaneRef.Name,
					Namespace: controlPlaneRef.Namespace,
				},
			},
		}
	}
}

func buildAPIEndpoint(cluster *cmv1.Cluster) (*clusterv1.APIEndpoint, error) {
	parsedURL, err := url.ParseRequestURI(cluster.API().URL())
	if err != nil {
		return nil, err
	}
	host, portStr, err := net.SplitHostPort(parsedURL.Host)
	if err != nil {
		return nil, err
	}

	port, err := strconv.Atoi(portStr)
	if err != nil {
		return nil, err
	}

	return &clusterv1.APIEndpoint{
		Host: host,
		Port: int32(port), // #nosec G109
	}, nil
}
