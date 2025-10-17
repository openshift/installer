/*
Copyright 2025 Nutanix

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

	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
	kerrors "k8s.io/apimachinery/pkg/util/errors"
	coreinformers "k8s.io/client-go/informers/core/v1"
	"k8s.io/utils/ptr"
	capiv1 "sigs.k8s.io/cluster-api/api/v1beta1"
	"sigs.k8s.io/cluster-api/util/conditions"
	"sigs.k8s.io/cluster-api/util/patch"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller"
	ctrlutil "sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
	"sigs.k8s.io/controller-runtime/pkg/handler"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"

	infrav1 "github.com/nutanix-cloud-native/cluster-api-provider-nutanix/api/v1beta1"
)

// NutanixFailureDomainReconciler reconciles a NutanixFailureDomain object
type NutanixFailureDomainReconciler struct {
	client.Client
	SecretInformer    coreinformers.SecretInformer
	ConfigMapInformer coreinformers.ConfigMapInformer
	Scheme            *runtime.Scheme
	controllerConfig  *ControllerConfig
}

func NewNutanixFailureDomainReconciler(client client.Client, secretInformer coreinformers.SecretInformer, configMapInformer coreinformers.ConfigMapInformer, scheme *runtime.Scheme, copts ...ControllerConfigOpts) (*NutanixFailureDomainReconciler, error) {
	controllerConf := &ControllerConfig{}
	for _, opt := range copts {
		if err := opt(controllerConf); err != nil {
			return nil, err
		}
	}

	return &NutanixFailureDomainReconciler{
		Client:            client,
		SecretInformer:    secretInformer,
		ConfigMapInformer: configMapInformer,
		Scheme:            scheme,
		controllerConfig:  controllerConf,
	}, nil
}

// SetupWithManager sets up the NutanixFailureDomain controller with the Manager.
func (r *NutanixFailureDomainReconciler) SetupWithManager(ctx context.Context, mgr ctrl.Manager) error {
	copts := controller.Options{
		MaxConcurrentReconciles: r.controllerConfig.MaxConcurrentReconciles,
		RateLimiter:             r.controllerConfig.RateLimiter,
		SkipNameValidation:      ptr.To(r.controllerConfig.SkipNameValidation),
	}

	return ctrl.NewControllerManagedBy(mgr).
		Named("nutanixfailuredomain-controller").
		For(&infrav1.NutanixFailureDomain{}).
		Watches(
			&infrav1.NutanixMachine{},
			handler.EnqueueRequestsFromMapFunc(
				r.mapNutanixMachineToNutanixFailureDomain(),
			),
		).
		WithOptions(copts).
		Complete(r)
}

func (r *NutanixFailureDomainReconciler) mapNutanixMachineToNutanixFailureDomain() handler.MapFunc {
	return func(ctx context.Context, o client.Object) []ctrl.Request {
		log := ctrl.LoggerFrom(ctx)
		nm, ok := o.(*infrav1.NutanixMachine)
		if !ok {
			log.Error(fmt.Errorf("expected a NutanixMachine object but was %T", o), "unexpected type")
			return nil
		}

		reqs := make([]ctrl.Request, 0)
		if nm.Status.FailureDomain == nil {
			return reqs
		}

		// Fetch the NutanixFailureDomain object in the local namespace
		fdName := *nm.Status.FailureDomain
		nfd := &infrav1.NutanixFailureDomain{}
		nfdKey := client.ObjectKey{Name: fdName, Namespace: nm.Namespace}
		if err := r.Get(ctx, nfdKey, nfd); err != nil {
			log.Error(err, "Failed to fetch the nutanix failure domain object for nutanix machine")
			return nil
		}

		objKey := client.ObjectKey{Name: nfd.Name, Namespace: nfd.Namespace}
		reqs = append(reqs, ctrl.Request{NamespacedName: objKey})
		return reqs
	}
}

//+kubebuilder:rbac:groups="",resources=secrets,verbs=get;list;watch
//+kubebuilder:rbac:groups="",resources=configmaps,verbs=get;list;watch
// +kubebuilder:rbac:groups=infrastructure.cluster.x-k8s.io,resources=nutanixclusters,verbs=get;list
// +kubebuilder:rbac:groups=infrastructure.cluster.x-k8s.io,resources=nutanixmachines,verbs=get;list
// +kubebuilder:rbac:groups=infrastructure.cluster.x-k8s.io,resources=nutanixfailuredomains,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=infrastructure.cluster.x-k8s.io,resources=nutanixfailuredomains/status,verbs=get;update;patch
// +kubebuilder:rbac:groups=infrastructure.cluster.x-k8s.io,resources=nutanixfailuredomains/finalizers,verbs=get;update;patch

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
// TODO(user): Modify the Reconcile function to compare the state specified by
// the NutanixFailureDomain object against the actual cluster state, and then
// perform operations to make the cluster state reflect the state specified by
// the user.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.8.3/pkg/reconcile
func (r *NutanixFailureDomainReconciler) Reconcile(ctx context.Context, req ctrl.Request) (res ctrl.Result, reterr error) {
	log := ctrl.LoggerFrom(ctx)
	log.Info("Reconciling the NutanixFailureDomain")

	// Fetch the NutanixFailureDomain instance
	fd := &infrav1.NutanixFailureDomain{}
	if err := r.Get(ctx, req.NamespacedName, fd); err != nil {
		if errors.IsNotFound(err) {
			// Request object not found, could have been deleted after reconcile request.
			log.Info("The NutanixFailureDomain object not found. Ignoring since object must be deleted.")
			return reconcile.Result{}, nil
		}

		// Error reading the object - requeue the request.
		log.Error(err, "failed to fetch the NutanixFailureDomain object")
		return reconcile.Result{}, err
	}

	// Initialize the patch helper.
	patchHelper, err := patch.NewHelper(fd, r.Client)
	if err != nil {
		log.Error(err, "Failed to configure the patch helper")
		return ctrl.Result{Requeue: true}, nil
	}

	defer func() {
		// Always attempt to Patch the NutanixFailureDomain object and its status after each reconciliation.
		if err := patchHelper.Patch(ctx, fd); err != nil {
			reterr = kerrors.NewAggregate([]error{reterr, err})
			log.Error(reterr, "Failed to patch NutanixFailureDomain.")
		} else {
			log.Info("Patched NutanixFailureDomain.", "status", fd.Status, "finalizers", fd.Finalizers)
		}
	}()

	// Add finalizer first if not set yet
	if !ctrlutil.ContainsFinalizer(fd, infrav1.NutanixFailureDomainFinalizer) {
		if ctrlutil.AddFinalizer(fd, infrav1.NutanixFailureDomainFinalizer) {
			// Add finalizer first avoid the race condition between init and delete.
			log.Info("Added the finalizer to the object", "finalizers", fd.Finalizers)
			return reconcile.Result{}, nil
		}
	}

	// Handle deletion
	if !fd.DeletionTimestamp.IsZero() {
		return r.reconcileDelete(ctx, fd)
	}

	return r.reconcileNormal(ctx, fd)
}

func (r *NutanixFailureDomainReconciler) reconcileDelete(ctx context.Context, fd *infrav1.NutanixFailureDomain) (ctrl.Result, error) {
	log := ctrl.LoggerFrom(ctx)
	log.Info("Handling NutanixFailureDomain deletion")

	// Check if there are NutanixMachines using this failure domain
	// A map with nutanixmachine name as key and the cluster name as value
	ntxMachines := map[string]string{}
	ntxMachineList := &infrav1.NutanixMachineList{}

	if err := r.List(ctx, ntxMachineList, client.InNamespace(fd.Namespace)); err != nil {
		return ctrl.Result{}, err
	}

	for _, nm := range ntxMachineList.Items {
		if !nm.DeletionTimestamp.IsZero() {
			continue
		}

		if nm.Status.FailureDomain != nil && *nm.Status.FailureDomain == fd.Name {
			ntxMachines[nm.Name] = nm.GetLabels()[capiv1.ClusterNameLabel]
		}
	}

	if len(ntxMachines) == 0 {
		conditions.MarkTrue(fd, infrav1.FailureDomainSafeForDeletionCondition)

		// Remove the finalizer from the failure domain object
		ctrlutil.RemoveFinalizer(fd, infrav1.NutanixFailureDomainFinalizer)
		return ctrl.Result{}, nil
	}

	errMsg := fmt.Sprintf("The failure domain is used by machines: %v", ntxMachines)
	conditions.MarkFalse(fd, infrav1.FailureDomainSafeForDeletionCondition,
		infrav1.FailureDomainInUseReason, capiv1.ConditionSeverityError, "%s", errMsg)

	reterr := fmt.Errorf("the failure domain %q is not safe for deletion since it is in use", fd.Name)
	log.Error(reterr, errMsg)
	return ctrl.Result{}, reterr
}

func (r *NutanixFailureDomainReconciler) reconcileNormal(ctx context.Context, fd *infrav1.NutanixFailureDomain) (ctrl.Result, error) {
	log := ctrl.LoggerFrom(ctx)
	log.Info("Handling NutanixFailureDomain reconciling")

	// Remove the FailureDomainSafeForDeletionCondition if there are any
	conditions.Delete(fd, infrav1.FailureDomainSafeForDeletionCondition)

	return ctrl.Result{}, nil
}
