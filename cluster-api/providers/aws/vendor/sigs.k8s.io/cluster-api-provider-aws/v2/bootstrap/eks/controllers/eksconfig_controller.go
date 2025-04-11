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

// Package controllers provides a way to reconcile EKSConfig objects.
package controllers

import (
	"bytes"
	"context"
	"time"

	"github.com/pkg/errors"
	corev1 "k8s.io/api/core/v1"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
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
	clusterv1 "sigs.k8s.io/cluster-api/api/v1beta1"
	bsutil "sigs.k8s.io/cluster-api/bootstrap/util"
	expclusterv1 "sigs.k8s.io/cluster-api/exp/api/v1beta1"
	"sigs.k8s.io/cluster-api/feature"
	"sigs.k8s.io/cluster-api/util"
	"sigs.k8s.io/cluster-api/util/annotations"
	"sigs.k8s.io/cluster-api/util/conditions"
	"sigs.k8s.io/cluster-api/util/patch"
	"sigs.k8s.io/cluster-api/util/predicates"
)

// EKSConfigReconciler reconciles a EKSConfig object.
type EKSConfigReconciler struct {
	client.Client
	Scheme           *runtime.Scheme
	WatchFilterValue string
}

// +kubebuilder:rbac:groups=bootstrap.cluster.x-k8s.io,resources=eksconfigs,verbs=get;list;watch;update;patch
// +kubebuilder:rbac:groups=bootstrap.cluster.x-k8s.io,resources=eksconfigs/status,verbs=get;update;patch
// +kubebuilder:rbac:groups=controlplane.cluster.x-k8s.io,resources=awsmanagedcontrolplanes,verbs=get;list;watch
// +kubebuilder:rbac:groups=cluster.x-k8s.io,resources=machines;machinepools;clusters,verbs=get;list;watch
// +kubebuilder:rbac:groups=cluster.x-k8s.io,resources=machinepools,verbs=get;list;watch
// +kubebuilder:rbac:groups="",resources=secrets,verbs=get;list;watch;create;update;delete;

func (r *EKSConfigReconciler) Reconcile(ctx context.Context, req ctrl.Request) (_ ctrl.Result, rerr error) {
	log := logger.FromContext(ctx)

	// get EKSConfig
	config := &eksbootstrapv1.EKSConfig{}
	if err := r.Client.Get(ctx, req.NamespacedName, config); err != nil {
		if apierrors.IsNotFound(err) {
			return ctrl.Result{}, nil
		}
		log.Error(err, "Failed to get config")
		return ctrl.Result{}, err
	}
	log = log.WithValues("EKSConfig", config.GetName())

	// check owner references and look up owning Machine object
	configOwner, err := bsutil.GetTypedConfigOwner(ctx, r.Client, config)
	if apierrors.IsNotFound(err) {
		// no error here, requeue until we find an owner
		log.Debug("eksconfig failed to look up owner reference, re-queueing")
		return ctrl.Result{RequeueAfter: time.Minute}, nil
	}
	if err != nil {
		log.Error(err, "eksconfig failed to get owner")
		return ctrl.Result{}, err
	}
	if configOwner == nil {
		// no error, requeue until we find an owner
		log.Debug("eksconfig has no owner reference set, re-queueing")
		return ctrl.Result{RequeueAfter: time.Minute}, nil
	}

	log = log.WithValues(configOwner.GetKind(), configOwner.GetName())

	cluster, err := util.GetClusterByName(ctx, r.Client, configOwner.GetNamespace(), configOwner.ClusterName())
	if err != nil {
		if errors.Is(err, util.ErrNoCluster) {
			log.Info("EKSConfig does not belong to a cluster yet, re-queuing until it's part of a cluster")
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

	if annotations.IsPaused(cluster, config) {
		log.Info("Reconciliation is paused for this object")
		return ctrl.Result{}, nil
	}

	patchHelper, err := patch.NewHelper(config, r.Client)
	if err != nil {
		return ctrl.Result{}, err
	}

	// set up defer block for updating config
	defer func() {
		conditions.SetSummary(config,
			conditions.WithConditions(
				eksbootstrapv1.DataSecretAvailableCondition,
			),
			conditions.WithStepCounter(),
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

	return ctrl.Result{}, r.joinWorker(ctx, cluster, config, configOwner)
}

func (r *EKSConfigReconciler) resolveFiles(ctx context.Context, cfg *eksbootstrapv1.EKSConfig) ([]eksbootstrapv1.File, error) {
	collected := make([]eksbootstrapv1.File, 0, len(cfg.Spec.Files))

	for i := range cfg.Spec.Files {
		in := cfg.Spec.Files[i]
		if in.ContentFrom != nil {
			data, err := r.resolveSecretFileContent(ctx, cfg.Namespace, in)
			if err != nil {
				return nil, errors.Wrapf(err, "failed to resolve file source")
			}
			in.ContentFrom = nil
			in.Content = string(data)
		}
		collected = append(collected, in)
	}

	return collected, nil
}

func (r *EKSConfigReconciler) resolveSecretFileContent(ctx context.Context, ns string, source eksbootstrapv1.File) ([]byte, error) {
	secret := &corev1.Secret{}
	key := types.NamespacedName{Namespace: ns, Name: source.ContentFrom.Secret.Name}
	if err := r.Client.Get(ctx, key, secret); err != nil {
		if apierrors.IsNotFound(err) {
			return nil, errors.Wrapf(err, "secret not found: %s", key)
		}
		return nil, errors.Wrapf(err, "failed to retrieve Secret %q", key)
	}
	data, ok := secret.Data[source.ContentFrom.Secret.Key]
	if !ok {
		return nil, errors.Errorf("secret references non-existent secret key: %q", source.ContentFrom.Secret.Key)
	}
	return data, nil
}

func (r *EKSConfigReconciler) joinWorker(ctx context.Context, cluster *clusterv1.Cluster, config *eksbootstrapv1.EKSConfig, configOwner *bsutil.ConfigOwner) error {
	log := logger.FromContext(ctx)

	// only need to reconcile the secret for Machine kinds once, but MachinePools need updates for new launch templates
	if config.Status.DataSecretName != nil && configOwner.GetKind() == "Machine" {
		secretKey := client.ObjectKey{Namespace: config.Namespace, Name: *config.Status.DataSecretName}
		log = log.WithValues("data-secret-name", secretKey.Name)
		existingSecret := &corev1.Secret{}

		// No error here means the Secret exists and we have no
		// reason to proceed.
		err := r.Client.Get(ctx, secretKey, existingSecret)
		switch {
		case err == nil:
			return nil
		case !apierrors.IsNotFound(err):
			log.Error(err, "unable to check for existing bootstrap secret")
			return err
		}
	}

	if cluster.Spec.ControlPlaneRef == nil || cluster.Spec.ControlPlaneRef.Kind != "AWSManagedControlPlane" {
		return errors.New("Cluster's controlPlaneRef needs to be an AWSManagedControlPlane in order to use the EKS bootstrap provider")
	}

	if !cluster.Status.InfrastructureReady {
		log.Info("Cluster infrastructure is not ready")
		conditions.MarkFalse(config,
			eksbootstrapv1.DataSecretAvailableCondition,
			eksbootstrapv1.WaitingForClusterInfrastructureReason,
			clusterv1.ConditionSeverityInfo, "")
		return nil
	}

	if !conditions.IsTrue(cluster, clusterv1.ControlPlaneInitializedCondition) {
		log.Info("Control Plane has not yet been initialized")
		conditions.MarkFalse(config, eksbootstrapv1.DataSecretAvailableCondition, eksbootstrapv1.WaitingForControlPlaneInitializationReason, clusterv1.ConditionSeverityInfo, "")
		return nil
	}

	controlPlane := &ekscontrolplanev1.AWSManagedControlPlane{}
	if err := r.Get(ctx, client.ObjectKey{Name: cluster.Spec.ControlPlaneRef.Name, Namespace: cluster.Spec.ControlPlaneRef.Namespace}, controlPlane); err != nil {
		return err
	}

	log.Info("Generating userdata")
	files, err := r.resolveFiles(ctx, config)
	if err != nil {
		log.Info("Failed to resolve files for user data")
		conditions.MarkFalse(config, eksbootstrapv1.DataSecretAvailableCondition, eksbootstrapv1.DataSecretGenerationFailedReason, clusterv1.ConditionSeverityWarning, "%s", err.Error())
		return err
	}

	nodeInput := &userdata.NodeInput{
		// AWSManagedControlPlane webhooks default and validate EKSClusterName
		ClusterName:              controlPlane.Spec.EKSClusterName,
		KubeletExtraArgs:         config.Spec.KubeletExtraArgs,
		ContainerRuntime:         config.Spec.ContainerRuntime,
		DNSClusterIP:             config.Spec.DNSClusterIP,
		DockerConfigJSON:         config.Spec.DockerConfigJSON,
		APIRetryAttempts:         config.Spec.APIRetryAttempts,
		UseMaxPods:               config.Spec.UseMaxPods,
		PreBootstrapCommands:     config.Spec.PreBootstrapCommands,
		PostBootstrapCommands:    config.Spec.PostBootstrapCommands,
		BootstrapCommandOverride: config.Spec.BootstrapCommandOverride,
		NTP:                      config.Spec.NTP,
		Users:                    config.Spec.Users,
		DiskSetup:                config.Spec.DiskSetup,
		Mounts:                   config.Spec.Mounts,
		Files:                    files,
	}
	if config.Spec.PauseContainer != nil {
		nodeInput.PauseContainerAccount = &config.Spec.PauseContainer.AccountNumber
		nodeInput.PauseContainerVersion = &config.Spec.PauseContainer.Version
	}

	// Check if IPv6 was provided to the user configuration first
	// If not, we also check if the cluster is ipv6 based.
	if config.Spec.ServiceIPV6Cidr != nil && *config.Spec.ServiceIPV6Cidr != "" {
		nodeInput.ServiceIPV6Cidr = config.Spec.ServiceIPV6Cidr
		nodeInput.IPFamily = ptr.To[string]("ipv6")
	}

	// we don't want to override any manually set configuration options.
	if config.Spec.ServiceIPV6Cidr == nil && controlPlane.Spec.NetworkSpec.VPC.IsIPv6Enabled() {
		log.Info("Adding ipv6 data to userdata....")
		nodeInput.ServiceIPV6Cidr = ptr.To[string](controlPlane.Spec.NetworkSpec.VPC.IPv6.CidrBlock)
		nodeInput.IPFamily = ptr.To[string]("ipv6")
	}

	// generate userdata
	userDataScript, err := userdata.NewNode(nodeInput)
	if err != nil {
		log.Error(err, "Failed to create a worker join configuration")
		conditions.MarkFalse(config, eksbootstrapv1.DataSecretAvailableCondition, eksbootstrapv1.DataSecretGenerationFailedReason, clusterv1.ConditionSeverityWarning, "")
		return err
	}

	// store userdata as secret
	if err := r.storeBootstrapData(ctx, cluster, config, userDataScript); err != nil {
		log.Error(err, "Failed to store bootstrap data")
		conditions.MarkFalse(config, eksbootstrapv1.DataSecretAvailableCondition, eksbootstrapv1.DataSecretGenerationFailedReason, clusterv1.ConditionSeverityWarning, "")
		return err
	}

	return nil
}

func (r *EKSConfigReconciler) SetupWithManager(ctx context.Context, mgr ctrl.Manager, option controller.Options) error {
	b := ctrl.NewControllerManagedBy(mgr).
		For(&eksbootstrapv1.EKSConfig{}).
		WithOptions(option).
		WithEventFilter(predicates.ResourceNotPausedAndHasFilterLabel(mgr.GetScheme(), logger.FromContext(ctx).GetLogger(), r.WatchFilterValue)).
		Watches(
			&clusterv1.Machine{},
			handler.EnqueueRequestsFromMapFunc(r.MachineToBootstrapMapFunc),
		)

	if feature.Gates.Enabled(feature.MachinePool) {
		b = b.Watches(
			&expclusterv1.MachinePool{},
			handler.EnqueueRequestsFromMapFunc(r.MachinePoolToBootstrapMapFunc),
		)
	}

	c, err := b.Build(r)
	if err != nil {
		return errors.Wrap(err, "failed setting up with a controller manager")
	}

	err = c.Watch(
		source.Kind[client.Object](mgr.GetCache(), &clusterv1.Cluster{},
			handler.EnqueueRequestsFromMapFunc((r.ClusterToEKSConfigs)),
			predicates.ClusterPausedTransitionsOrInfrastructureReady(mgr.GetScheme(), logger.FromContext(ctx).GetLogger())),
	)
	if err != nil {
		return errors.Wrap(err, "failed adding watch for Clusters to controller manager")
	}

	return nil
}

// storeBootstrapData creates a new secret with the data passed in as input,
// sets the reference in the configuration status and ready to true.
func (r *EKSConfigReconciler) storeBootstrapData(ctx context.Context, cluster *clusterv1.Cluster, config *eksbootstrapv1.EKSConfig, data []byte) error {
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
				return errors.Wrap(err, "failed to create bootstrap data secret for EKSConfig")
			}
			log.Info("created bootstrap data secret for EKSConfig", "secret", klog.KObj(secret))
		} else {
			return errors.Wrap(err, "failed to get data secret for EKSConfig")
		}
	} else {
		updated, err := r.updateBootstrapSecret(ctx, secret, data)
		if err != nil {
			return errors.Wrap(err, "failed to update data secret for EKSConfig")
		}
		if updated {
			log.Info("updated bootstrap data secret for EKSConfig", "secret", klog.KObj(secret))
		} else {
			log.Trace("no change in bootstrap data secret for EKSConfig", "secret", klog.KObj(secret))
		}
	}

	config.Status.DataSecretName = ptr.To[string](secret.Name)
	config.Status.Ready = true
	conditions.MarkTrue(config, eksbootstrapv1.DataSecretAvailableCondition)
	return nil
}

// MachineToBootstrapMapFunc is a handler.ToRequestsFunc to be used to enqueue requests
// for EKSConfig reconciliation.
func (r *EKSConfigReconciler) MachineToBootstrapMapFunc(_ context.Context, o client.Object) []ctrl.Request {
	result := []ctrl.Request{}

	m, ok := o.(*clusterv1.Machine)
	if !ok {
		klog.Errorf("Expected a Machine but got a %T", o)
	}
	if m.Spec.Bootstrap.ConfigRef != nil && m.Spec.Bootstrap.ConfigRef.GroupVersionKind() == eksbootstrapv1.GroupVersion.WithKind("EKSConfig") {
		name := client.ObjectKey{Namespace: m.Namespace, Name: m.Spec.Bootstrap.ConfigRef.Name}
		result = append(result, ctrl.Request{NamespacedName: name})
	}
	return result
}

// MachinePoolToBootstrapMapFunc is a handler.ToRequestsFunc to be uses to enqueue requests
// for EKSConfig reconciliation.
func (r *EKSConfigReconciler) MachinePoolToBootstrapMapFunc(_ context.Context, o client.Object) []ctrl.Request {
	result := []ctrl.Request{}

	m, ok := o.(*expclusterv1.MachinePool)
	if !ok {
		klog.Errorf("Expected a MachinePool but got a %T", o)
	}
	configRef := m.Spec.Template.Spec.Bootstrap.ConfigRef
	if configRef != nil && configRef.GroupVersionKind().GroupKind() == eksbootstrapv1.GroupVersion.WithKind("EKSConfig").GroupKind() {
		name := client.ObjectKey{Namespace: m.Namespace, Name: configRef.Name}
		result = append(result, ctrl.Request{NamespacedName: name})
	}

	return result
}

// ClusterToEKSConfigs is a handler.ToRequestsFunc to be used to enqueue requests for
// EKSConfig reconciliation.
func (r *EKSConfigReconciler) ClusterToEKSConfigs(_ context.Context, o client.Object) []ctrl.Request {
	result := []ctrl.Request{}

	c, ok := o.(*clusterv1.Cluster)
	if !ok {
		klog.Errorf("Expected a Cluster but got a %T", o)
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
		if m.Spec.Bootstrap.ConfigRef != nil &&
			m.Spec.Bootstrap.ConfigRef.GroupVersionKind().GroupKind() == eksbootstrapv1.GroupVersion.WithKind("EKSConfig").GroupKind() {
			name := client.ObjectKey{Namespace: m.Namespace, Name: m.Spec.Bootstrap.ConfigRef.Name}
			result = append(result, ctrl.Request{NamespacedName: name})
		}
	}

	return result
}

// Create the Secret containing bootstrap userdata.
func (r *EKSConfigReconciler) createBootstrapSecret(ctx context.Context, cluster *clusterv1.Cluster, config *eksbootstrapv1.EKSConfig, data []byte) (*corev1.Secret, error) {
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
					Kind:       "EKSConfig",
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
func (r *EKSConfigReconciler) updateBootstrapSecret(ctx context.Context, secret *corev1.Secret, data []byte) (bool, error) {
	if secret.Data == nil {
		secret.Data = make(map[string][]byte)
	}
	if !bytes.Equal(secret.Data["value"], data) {
		secret.Data["value"] = data
		return true, r.Client.Update(ctx, secret)
	}
	return false, nil
}
