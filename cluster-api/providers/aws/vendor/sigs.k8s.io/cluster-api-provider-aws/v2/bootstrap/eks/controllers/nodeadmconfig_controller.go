/*
Copyright 2026 The Kubernetes Authors.

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
	"bytes"
	"context"
	"encoding/base64"
	"fmt"
	"time"

	"github.com/pkg/errors"
	corev1 "k8s.io/api/core/v1"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/klog/v2"
	"k8s.io/utils/ptr"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller"
	"sigs.k8s.io/controller-runtime/pkg/handler"
	"sigs.k8s.io/controller-runtime/pkg/source"

	eksbootstrapv1 "sigs.k8s.io/cluster-api-provider-aws/v2/bootstrap/eks/api/v1beta2"
	"sigs.k8s.io/cluster-api-provider-aws/v2/bootstrap/eks/internal/userdata"
	ekscontrolplanev1 "sigs.k8s.io/cluster-api-provider-aws/v2/controlplane/eks/api/v1beta2"
	"sigs.k8s.io/cluster-api-provider-aws/v2/pkg/logger"
	"sigs.k8s.io/cluster-api-provider-aws/v2/util/paused"
	clusterv1beta1 "sigs.k8s.io/cluster-api/api/core/v1beta1"
	clusterv1 "sigs.k8s.io/cluster-api/api/core/v1beta2"
	bsutil "sigs.k8s.io/cluster-api/bootstrap/util"
	"sigs.k8s.io/cluster-api/feature"
	"sigs.k8s.io/cluster-api/util"
	v1beta1conditions "sigs.k8s.io/cluster-api/util/deprecated/v1beta1/conditions"
	kubeconfigutil "sigs.k8s.io/cluster-api/util/kubeconfig"
	"sigs.k8s.io/cluster-api/util/patch"
	"sigs.k8s.io/cluster-api/util/predicates"
)

// NodeadmConfigReconciler reconciles a NodeadmConfig object.
type NodeadmConfigReconciler struct {
	client.Client
	Scheme           *runtime.Scheme
	WatchFilterValue string
}

// +kubebuilder:rbac:groups=bootstrap.cluster.x-k8s.io,resources=nodeadmconfigs,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=bootstrap.cluster.x-k8s.io,resources=nodeadmconfigs/status,verbs=get;update;patch
// +kubebuilder:rbac:groups=controlplane.cluster.x-k8s.io,resources=awsmanagedcontrolplanes,verbs=get;list;watch
// +kubebuilder:rbac:groups=cluster.x-k8s.io,resources=machines;machinepools;clusters,verbs=get;list;watch
// +kubebuilder:rbac:groups=cluster.x-k8s.io,resources=machinepools,verbs=get;list;watch
// +kubebuilder:rbac:groups="",resources=secrets,verbs=get;list;watch;create;update;delete;

func (r *NodeadmConfigReconciler) Reconcile(ctx context.Context, req ctrl.Request) (_ ctrl.Result, rerr error) {
	log := logger.FromContext(ctx)

	// get NodeadmConfig
	config := &eksbootstrapv1.NodeadmConfig{}
	if err := r.Client.Get(ctx, req.NamespacedName, config); err != nil {
		if apierrors.IsNotFound(err) {
			return ctrl.Result{}, nil
		}
		log.Error(err, "Failed to get config")
		return ctrl.Result{}, err
	}
	log = log.WithValues("NodeadmConfig", config.GetName())

	// check owner references and look up owning Machine object
	configOwner, err := bsutil.GetTypedConfigOwner(ctx, r.Client, config)
	if apierrors.IsNotFound(err) {
		// no error here, requeue until we find an owner
		log.Debug("NodeadmConfig failed to look up owner reference, re-queueing")
		return ctrl.Result{RequeueAfter: time.Minute}, nil
	}
	if err != nil {
		log.Error(err, "NodeadmConfig failed to get owner")
		return ctrl.Result{}, err
	}
	if configOwner == nil {
		// no error, requeue until we find an owner
		log.Debug("NodeadmConfig has no owner reference set, re-queueing")
		return ctrl.Result{RequeueAfter: time.Minute}, nil
	}

	log = log.WithValues(configOwner.GetKind(), configOwner.GetName())

	cluster, err := util.GetClusterByName(ctx, r.Client, configOwner.GetNamespace(), configOwner.ClusterName())
	if err != nil {
		if errors.Is(err, util.ErrNoCluster) {
			log.Info("NodeadmConfig does not belong to a cluster yet, re-queuing until it's part of a cluster")
			return ctrl.Result{RequeueAfter: time.Minute}, nil
		}
		if apierrors.IsNotFound(err) {
			log.Info("Cluster does not exist yet, re-queueing until it is created")
			return ctrl.Result{RequeueAfter: time.Minute}, nil
		}
		log.Error(err, "Could not get cluster with metadata")
		return ctrl.Result{}, err
	}
	log = log.WithValues("cluster", klog.KObj(cluster))

	if isPaused, conditionChanged, err := paused.EnsurePausedCondition(ctx, r.Client, cluster, config); err != nil || isPaused || conditionChanged {
		return ctrl.Result{}, err
	}

	patchHelper, err := patch.NewHelper(config, r.Client)
	if err != nil {
		return ctrl.Result{}, err
	}

	// set up defer block for updating config
	defer func() {
		v1beta1conditions.SetSummary(config,
			v1beta1conditions.WithConditions(
				eksbootstrapv1.DataSecretAvailableCondition,
			),
			v1beta1conditions.WithStepCounter(),
		)

		patchOpts := []patch.Option{}
		if rerr == nil {
			patchOpts = append(patchOpts, patch.WithStatusObservedGeneration{})
		}
		if err := patchHelper.Patch(ctx, config, patchOpts...); err != nil {
			log.Error(rerr, "Failed to patch config")
			if rerr == nil {
				rerr = err
			}
		}
	}()

	return r.joinWorker(ctx, cluster, config, configOwner)
}

func (r *NodeadmConfigReconciler) joinWorker(ctx context.Context, cluster *clusterv1.Cluster, config *eksbootstrapv1.NodeadmConfig, configOwner *bsutil.ConfigOwner) (ctrl.Result, error) {
	log := logger.FromContext(ctx)

	// only need to reconcile the secret for Machine kinds once, but MachinePools need updates for new launch templates
	if config.Status.DataSecretName != nil && configOwner.GetKind() == "Machine" {
		secretKey := client.ObjectKey{Namespace: config.Namespace, Name: *config.Status.DataSecretName}
		log = log.WithValues("data-secret-name", secretKey.Name)
		existingSecret := &corev1.Secret{}

		err := r.Client.Get(ctx, secretKey, existingSecret)
		if err != nil && !apierrors.IsNotFound(err) {
			log.Error(err, "unable to check for existing bootstrap secret")
			return ctrl.Result{}, err
		}
		if err == nil {
			// We already have a secret that we don't need to regenerate
			return ctrl.Result{}, nil
		}
	}

	if cluster.Spec.ControlPlaneRef.Kind != "AWSManagedControlPlane" {
		return ctrl.Result{}, errors.New("Cluster's controlPlaneRef needs to be an AWSManagedControlPlane in order to use the EKS bootstrap provider")
	}

	if !ptr.Deref(cluster.Status.Initialization.InfrastructureProvisioned, false) {
		log.Info("Cluster infrastructure is not ready")
		v1beta1conditions.MarkFalse(config,
			eksbootstrapv1.DataSecretAvailableCondition,
			eksbootstrapv1.WaitingForClusterInfrastructureReason,
			clusterv1beta1.ConditionSeverityInfo, "")
		return ctrl.Result{}, nil
	}

	if !ptr.Deref(cluster.Status.Initialization.ControlPlaneInitialized, false) {
		log.Info("Control Plane has not yet been initialized")
		v1beta1conditions.MarkFalse(config, eksbootstrapv1.DataSecretAvailableCondition, eksbootstrapv1.WaitingForControlPlaneInitializationReason, clusterv1beta1.ConditionSeverityInfo, "")
		return ctrl.Result{RequeueAfter: 30 * time.Second}, nil
	}

	controlPlane := &ekscontrolplanev1.AWSManagedControlPlane{}
	if err := r.Get(ctx, client.ObjectKey{Name: cluster.Spec.ControlPlaneRef.Name, Namespace: cluster.Namespace}, controlPlane); err != nil {
		return ctrl.Result{}, errors.Wrap(err, "failed to get control plane")
	}
	// Check if control plane is ready
	if !v1beta1conditions.IsTrue(controlPlane, ekscontrolplanev1.EKSControlPlaneReadyCondition) {
		log.Info("Waiting for control plane to be ready")
		v1beta1conditions.MarkFalse(
			config,
			eksbootstrapv1.DataSecretAvailableCondition,
			eksbootstrapv1.DataSecretGenerationFailedReason,
			clusterv1beta1.ConditionSeverityInfo,
			"Control plane is not initialized yet",
		)
		return ctrl.Result{RequeueAfter: 30 * time.Second}, nil
	}
	log.Info("Control plane is ready, proceeding with userdata generation")

	log.Info("Generating userdata")
	fileResolver := FileResolver{Client: r.Client}
	files, err := fileResolver.ResolveFiles(ctx, config.Namespace, config.Spec.Files)
	if err != nil {
		log.Info("Failed to resolve files for user data")
		v1beta1conditions.MarkFalse(config, eksbootstrapv1.DataSecretAvailableCondition, eksbootstrapv1.DataSecretGenerationFailedReason, clusterv1beta1.ConditionSeverityWarning, "%s", err.Error())
		return ctrl.Result{}, err
	}

	serviceCIDR := ""
	if len(cluster.Spec.ClusterNetwork.Services.CIDRBlocks) > 0 {
		serviceCIDR = cluster.Spec.ClusterNetwork.Services.CIDRBlocks[0]
	}
	nodeInput := &userdata.NodeadmInput{
		// AWSManagedControlPlane webhooks default and validate EKSClusterName
		ClusterName:        controlPlane.Spec.EKSClusterName,
		PreNodeadmCommands: config.Spec.PreNodeadmCommands,
		Users:              config.Spec.Users,
		NTP:                config.Spec.NTP,
		DiskSetup:          config.Spec.DiskSetup,
		Mounts:             config.Spec.Mounts,
		Files:              files,
		ServiceCIDR:        serviceCIDR,
		APIServerEndpoint:  cluster.Spec.ControlPlaneEndpoint.Host,
	}
	if config.Spec.Kubelet != nil {
		nodeInput.KubeletFlags = config.Spec.Kubelet.Flags
		if config.Spec.Kubelet.Config != nil {
			nodeInput.KubeletConfig = config.Spec.Kubelet.Config
		}
	}
	if config.Spec.Containerd != nil {
		nodeInput.ContainerdConfig = config.Spec.Containerd.Config
		if config.Spec.Containerd.BaseRuntimeSpec != nil {
			nodeInput.ContainerdBaseRuntimeSpec = config.Spec.Containerd.BaseRuntimeSpec
		}
	}
	if config.Spec.FeatureGates != nil {
		nodeInput.FeatureGates = config.Spec.FeatureGates
	}

	// Fetch CA cert from KubeConfig secret
	obj := client.ObjectKey{
		Namespace: cluster.Namespace,
		Name:      cluster.Name,
	}
	ca, err := extractCAFromSecret(ctx, r.Client, obj)
	if err != nil {
		log.Error(err, "Failed to extract CA from kubeconfig secret")
		v1beta1conditions.MarkFalse(config, eksbootstrapv1.DataSecretAvailableCondition,
			eksbootstrapv1.DataSecretGenerationFailedReason,
			clusterv1beta1.ConditionSeverityWarning,
			"Failed to extract CA from kubeconfig secret: %v", err)
		return ctrl.Result{}, err
	}
	nodeInput.CACert = ca

	log.Info("Generating nodeadm userdata",
		"cluster", controlPlane.Spec.EKSClusterName,
		"endpoint", nodeInput.APIServerEndpoint)
	// generate userdata
	userDataScript, err := userdata.NewNodeadmUserdata(nodeInput)
	if err != nil {
		log.Error(err, "Failed to create a worker join configuration")
		v1beta1conditions.MarkFalse(config, eksbootstrapv1.DataSecretAvailableCondition, eksbootstrapv1.DataSecretGenerationFailedReason, clusterv1beta1.ConditionSeverityWarning, "")
		return ctrl.Result{}, err
	}

	// store userdata as secret
	if err := r.storeBootstrapData(ctx, cluster, config, userDataScript); err != nil {
		log.Error(err, "Failed to store bootstrap data")
		v1beta1conditions.MarkFalse(config, eksbootstrapv1.DataSecretAvailableCondition, eksbootstrapv1.DataSecretGenerationFailedReason, clusterv1beta1.ConditionSeverityWarning, "")
		return ctrl.Result{}, err
	}

	v1beta1conditions.MarkTrue(config, eksbootstrapv1.DataSecretAvailableCondition)
	return ctrl.Result{}, nil
}

// storeBootstrapData creates a new secret with the data passed in as input,
// sets the reference in the configuration status and ready to true.
func (r *NodeadmConfigReconciler) storeBootstrapData(ctx context.Context, cluster *clusterv1.Cluster, config *eksbootstrapv1.NodeadmConfig, data []byte) error {
	log := logger.FromContext(ctx)

	// as secret creation and scope.Config status patch are not atomic operations
	// it is possible that secret creation happens but the config.Status patches are not applied
	secret := &corev1.Secret{}
	if err := r.Client.Get(ctx, client.ObjectKey{
		Name:      config.Name,
		Namespace: config.Namespace,
	}, secret); err != nil {
		if apierrors.IsNotFound(err) {
			if secret, err = r.createBootstrapSecret(ctx, cluster, config, data); err != nil {
				return errors.Wrap(err, "failed to create bootstrap data secret for NodeadmConfig")
			}
			log.Info("created bootstrap data secret for NodeadmConfig", "secret", klog.KObj(secret))
		} else {
			return errors.Wrap(err, "failed to get data secret for NodeadmConfig")
		}
	} else {
		updated, err := r.updateBootstrapSecret(ctx, secret, data)
		if err != nil {
			return errors.Wrap(err, "failed to update data secret for NodeadmConfig")
		}
		if updated {
			log.Info("updated bootstrap data secret for NodeadmConfig", "secret", klog.KObj(secret))
		} else {
			log.Trace("no change in bootstrap data secret for NodeadmConfig", "secret", klog.KObj(secret))
		}
	}

	config.Status.DataSecretName = ptr.To(secret.Name)
	config.Status.Initialization.DataSecretCreated = ptr.To(true)
	//nolint:staticcheck // we will support this implementation until CAPA is v1beta2 compliant
	config.Status.Ready = true
	v1beta1conditions.MarkTrue(config, eksbootstrapv1.DataSecretAvailableCondition)
	return nil
}

func (r *NodeadmConfigReconciler) createBootstrapSecret(ctx context.Context, cluster *clusterv1.Cluster, config *eksbootstrapv1.NodeadmConfig, data []byte) (*corev1.Secret, error) {
	secret := &corev1.Secret{
		ObjectMeta: metav1.ObjectMeta{
			Name:      config.Name,
			Namespace: config.Namespace,
			Labels: map[string]string{
				clusterv1.ClusterNameLabel: cluster.Name,
			},
			OwnerReferences: []metav1.OwnerReference{
				{
					APIVersion: eksbootstrapv1.GroupVersion.String(),
					Kind:       "NodeadmConfig",
					Name:       config.Name,
					UID:        config.UID,
					Controller: ptr.To[bool](true),
				},
			},
		},
		Data: map[string][]byte{
			"value": data,
		},
		Type: clusterv1.ClusterSecretType,
	}
	return secret, r.Client.Create(ctx, secret)
}

// Update the userdata in the bootstrap Secret.
func (r *NodeadmConfigReconciler) updateBootstrapSecret(ctx context.Context, secret *corev1.Secret, data []byte) (bool, error) {
	if secret.Data == nil {
		secret.Data = make(map[string][]byte)
	}
	if !bytes.Equal(secret.Data["value"], data) {
		secret.Data["value"] = data
		return true, r.Client.Update(ctx, secret)
	}
	return false, nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *NodeadmConfigReconciler) SetupWithManager(ctx context.Context, mgr ctrl.Manager, option controller.Options) error {
	b := ctrl.NewControllerManagedBy(mgr).
		For(&eksbootstrapv1.NodeadmConfig{}).
		WithOptions(option).
		WithEventFilter(predicates.ResourceHasFilterLabel(mgr.GetScheme(), logger.FromContext(ctx).GetLogger(), r.WatchFilterValue)).
		Watches(
			&clusterv1.Machine{},
			handler.EnqueueRequestsFromMapFunc(r.MachineToBootstrapMapFunc),
		)

	if feature.Gates.Enabled(feature.MachinePool) {
		b = b.Watches(
			&clusterv1.MachinePool{},
			handler.EnqueueRequestsFromMapFunc(r.MachinePoolToBootstrapMapFunc),
		)
	}

	c, err := b.Build(r)
	if err != nil {
		return errors.Wrap(err, "failed setting up with a controller manager")
	}

	err = c.Watch(
		source.Kind[client.Object](mgr.GetCache(), &clusterv1.Cluster{},
			handler.EnqueueRequestsFromMapFunc((r.ClusterToNodeadmConfigs)),
			predicates.ClusterPausedTransitionsOrInfrastructureProvisioned(mgr.GetScheme(), logger.FromContext(ctx).GetLogger())),
	)
	if err != nil {
		return errors.Wrap(err, "failed adding watch for Clusters to controller manager")
	}
	return nil
}

// MachineToBootstrapMapFunc is a handler.ToRequestsFunc to be used to enque requests for
// NodeadmConfig reconciliation.
func (r *NodeadmConfigReconciler) MachineToBootstrapMapFunc(_ context.Context, o client.Object) []ctrl.Request {
	result := []ctrl.Request{}

	m, ok := o.(*clusterv1.Machine)
	if !ok {
		klog.Errorf("Expected a Machine but got a %T", o)
		return result
	}
	if m.Spec.Bootstrap.ConfigRef.IsDefined() && m.Spec.Bootstrap.ConfigRef.APIGroup == eksbootstrapv1.GroupVersion.Group && m.Spec.Bootstrap.ConfigRef.Kind == eksbootstrapv1.NodeadmConfigKind {
		name := client.ObjectKey{Namespace: m.Namespace, Name: m.Spec.Bootstrap.ConfigRef.Name}
		result = append(result, ctrl.Request{NamespacedName: name})
	}
	return result
}

// MachinePoolToBootstrapMapFunc is a handler.ToRequestsFunc to be uses to enqueue requests
// for NodeadmConfig reconciliation.
func (r *NodeadmConfigReconciler) MachinePoolToBootstrapMapFunc(_ context.Context, o client.Object) []ctrl.Request {
	result := []ctrl.Request{}

	m, ok := o.(*clusterv1.MachinePool)
	if !ok {
		klog.Errorf("Expected a MachinePool but got a %T", o)
		return result
	}
	configRef := m.Spec.Template.Spec.Bootstrap.ConfigRef
	if configRef.IsDefined() && configRef.APIGroup == eksbootstrapv1.GroupVersion.Group && configRef.Kind == eksbootstrapv1.NodeadmConfigKind {
		name := client.ObjectKey{Namespace: m.Namespace, Name: configRef.Name}
		result = append(result, ctrl.Request{NamespacedName: name})
	}

	return result
}

// ClusterToNodeadmConfigs is a handler.ToRequestsFunc to be used to enqueue requests for
// NodeadmConfig reconciliation.
func (r *NodeadmConfigReconciler) ClusterToNodeadmConfigs(_ context.Context, o client.Object) []ctrl.Request {
	result := []ctrl.Request{}

	c, ok := o.(*clusterv1.Cluster)
	if !ok {
		klog.Errorf("Expected a Cluster but got a %T", o)
		return result
	}

	selectors := []client.ListOption{
		client.InNamespace(c.Namespace),
		client.MatchingLabels{
			clusterv1.ClusterNameLabel: c.Name,
		},
	}

	machineList := &clusterv1.MachineList{}
	if err := r.Client.List(context.Background(), machineList, selectors...); err != nil {
		return nil
	}

	for _, m := range machineList.Items {
		if m.Spec.Bootstrap.ConfigRef.IsDefined() &&
			m.Spec.Bootstrap.ConfigRef.APIGroup == eksbootstrapv1.GroupVersion.Group &&
			m.Spec.Bootstrap.ConfigRef.Kind == eksbootstrapv1.NodeadmConfigKind {
			name := client.ObjectKey{Namespace: m.Namespace, Name: m.Spec.Bootstrap.ConfigRef.Name}
			result = append(result, ctrl.Request{NamespacedName: name})
		}
	}

	return result
}

func extractCAFromSecret(ctx context.Context, c client.Client, obj client.ObjectKey) (string, error) {
	data, err := kubeconfigutil.FromSecret(ctx, c, obj)
	if err != nil {
		return "", errors.Wrapf(err, "failed to get kubeconfig secret %s", obj.Name)
	}
	config, err := clientcmd.Load(data)
	if err != nil {
		return "", errors.Wrapf(err, "failed to parse kubeconfig data from secret %s", obj.Name)
	}

	// Iterate through all clusters in the kubeconfig and use the first one with CA data
	for _, cluster := range config.Clusters {
		if len(cluster.CertificateAuthorityData) > 0 {
			return base64.StdEncoding.EncodeToString(cluster.CertificateAuthorityData), nil
		}
	}

	return "", fmt.Errorf("no cluster with CA data found in kubeconfig")
}
