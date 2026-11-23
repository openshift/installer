/*
Copyright 2022 The Kubernetes Authors.

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

package powervs

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/IBM-Cloud/power-go-client/power/models"

	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	kerrors "k8s.io/apimachinery/pkg/util/errors"
	"k8s.io/client-go/tools/record"
	"k8s.io/klog/v2"

	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"

	clusterv1 "sigs.k8s.io/cluster-api/api/core/v1beta2"
	clusterv1util "sigs.k8s.io/cluster-api/util"
	"sigs.k8s.io/cluster-api/util/conditions"
	deprecatedv1beta1conditions "sigs.k8s.io/cluster-api/util/conditions/deprecated/v1beta1"
	"sigs.k8s.io/cluster-api/util/finalizers"
	"sigs.k8s.io/cluster-api/util/patch"

	infrav1 "sigs.k8s.io/cluster-api-provider-ibmcloud/api/powervs/v1beta3"
	"sigs.k8s.io/cluster-api-provider-ibmcloud/pkg/cloud/endpoints"
	powervsscope "sigs.k8s.io/cluster-api-provider-ibmcloud/pkg/cloud/scope/powervs"
)

const (
	ibmPowerVSClusterKind = "IBMPowerVSCluster"
)

// IBMPowerVSImageReconciler reconciles a IBMPowerVSImage object.
type IBMPowerVSImageReconciler struct {
	client.Client
	Recorder        record.EventRecorder
	ServiceEndpoint []endpoints.ServiceEndpoint
	Scheme          *runtime.Scheme
	ClientBuilder   powervsscope.ClientBuilder
}

//+kubebuilder:rbac:groups=infrastructure.cluster.x-k8s.io,resources=ibmpowervsimages,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=infrastructure.cluster.x-k8s.io,resources=ibmpowervsimages/status,verbs=get;update;patch

// Reconcile implements controller runtime Reconciler interface and handles reconciliation logic for IBMPowerVSImage.
func (r *IBMPowerVSImageReconciler) Reconcile(ctx context.Context, req ctrl.Request) (_ ctrl.Result, reterr error) {
	log := ctrl.LoggerFrom(ctx)
	log.Info("Reconciling IBMPowerVSImage")
	defer log.Info("Finished reconciling IBMPowerVSImage")

	// Fetch the IBMPowerVSImage.
	ibmPowerVSImage := &infrav1.IBMPowerVSImage{}
	if err := r.Client.Get(ctx, req.NamespacedName, ibmPowerVSImage); err != nil {
		if apierrors.IsNotFound(err) {
			log.Info("IBMPowerVSImage not found")
			return ctrl.Result{}, nil
		}
		return ctrl.Result{}, fmt.Errorf("failed to get IBMPowerVSImage: %w", err)
	}

	// Add finalizer first if not set to avoid the race condition between init and delete.
	if finalizerAdded, err := finalizers.EnsureFinalizer(ctx, r.Client, ibmPowerVSImage, infrav1.IBMPowerVSImageFinalizer); err != nil || finalizerAdded {
		if err == nil {
			log.Info("Added finalizer to IBMPowerVSImage, requeuing")
		}
		return ctrl.Result{}, err
	}

	var cluster *infrav1.IBMPowerVSCluster
	scopeParams := powervsscope.ImageScopeParams{
		Client:          r.Client,
		IBMPowerVSImage: ibmPowerVSImage,
		ServiceEndpoint: r.ServiceEndpoint,
		ClientBuilder:   r.ClientBuilder,
		Recorder:        r.Recorder,
	}

	// Externally managed clusters might not be available during image deletion. Get the cluster only when image is still not deleted.
	if ibmPowerVSImage.DeletionTimestamp.IsZero() {
		var err error
		cluster, err = powervsscope.GetClusterByName(ctx, r.Client, ibmPowerVSImage.Namespace, ibmPowerVSImage.Spec.ClusterName)
		if err != nil {
			return ctrl.Result{}, err
		}
		scopeParams.Zone = cluster.Spec.Zone
	}

	// Initialize the patch helper
	patchHelper, err := patch.NewHelper(ibmPowerVSImage, r.Client)
	if err != nil {
		return ctrl.Result{}, fmt.Errorf("failed to init patch helper: %w", err)
	}

	// Always attempt to Patch the IBMPowerVSImage object and status after each reconciliation.
	defer func() {
		if err := patchIBMPowerVSImage(ctx, patchHelper, ibmPowerVSImage); err != nil {
			reterr = kerrors.NewAggregate([]error{reterr, err})
		}
	}()

	// Create the scope
	imageScope, err := powervsscope.NewPowerVSImageScope(ctx, scopeParams)
	if err != nil {
		if errors.Is(err, powervsscope.ErrWorkspaceNotInActiveState) {
			r.markCondition(ibmPowerVSImage, infrav1.WorkspaceReadyCondition, infrav1.ServiceInstanceReadyV1Beta2Condition, metav1.ConditionFalse, infrav1.WorkspaceNotReadyReason, clusterv1.ConditionSeverityError, err.Error())
		}
		r.markCondition(ibmPowerVSImage, infrav1.IBMPowerVSImageReadyCondition, infrav1.ImageReadyV1Beta2Condition, metav1.ConditionUnknown, infrav1.IBMPowerVSImageReadyUnknownReason, clusterv1.ConditionSeverityInfo, "Failed to create image scope")
		return ctrl.Result{}, fmt.Errorf("failed to create scope: %w", err)
	}

	// Handle deleted clusters.
	if !ibmPowerVSImage.DeletionTimestamp.IsZero() {
		return r.reconcileDelete(ctx, imageScope)
	}

	return r.reconcile(ctx, cluster, imageScope)
}

func (r *IBMPowerVSImageReconciler) reconcile(ctx context.Context, cluster *infrav1.IBMPowerVSCluster, imageScope *powervsscope.ImageScope) (ctrl.Result, error) {
	log := ctrl.LoggerFrom(ctx)

	// 1. Ensure Labels and Ownership
	if imageScope.IBMPowerVSImage.Labels == nil {
		imageScope.IBMPowerVSImage.Labels = make(map[string]string)
	}

	if _, ok := imageScope.IBMPowerVSImage.Labels[clusterv1.ClusterNameLabel]; !ok {
		imageScope.IBMPowerVSImage.Labels[clusterv1.ClusterNameLabel] = imageScope.IBMPowerVSImage.Spec.ClusterName
	}

	if r.shouldAdopt(*imageScope.IBMPowerVSImage) {
		log.Info("Setting OwnerRef on IBMPowerVSImage, requeuing")
		imageScope.IBMPowerVSImage.OwnerReferences = clusterv1util.EnsureOwnerRef(imageScope.IBMPowerVSImage.OwnerReferences, metav1.OwnerReference{
			APIVersion: infrav1.GroupVersion.String(),
			Kind:       ibmPowerVSClusterKind,
			Name:       cluster.Name,
			UID:        cluster.UID,
		})
		return ctrl.Result{}, nil
	}

	// 2. Mark Workspace Ready
	r.markCondition(imageScope.IBMPowerVSImage, infrav1.WorkspaceReadyCondition, infrav1.ServiceInstanceReadyV1Beta2Condition, metav1.ConditionTrue, infrav1.WorkspaceReadyReason, clusterv1.ConditionSeverityInfo, "")

	// 3. Import Job Polling Flow
	if jobID := imageScope.GetJobID(); jobID != "" {
		job, err := imageScope.IBMPowerVSClient.GetJob(ctx, jobID)
		if err != nil {
			log.Info("Unable to get job details", "jobID", jobID)
			return ctrl.Result{RequeueAfter: 2 * time.Minute}, err
		}

		if job.Status == nil || job.Status.State == nil {
			log.Info("Job status or state is currently nil, requeuing", "jobID", jobID)
			return ctrl.Result{RequeueAfter: 1 * time.Minute}, nil
		}

		imageScope.SetImageState(*job.Status.State)

		log.Info("Polling import job", "jobID", jobID, "state", *job.Status.State)
		switch imageScope.GetImageState() {
		case infrav1.PowerVSImageStateCompleted:
			log.Info("Import job completed, proceeding to reconcile image", "jobID", jobID)
			r.markCondition(imageScope.IBMPowerVSImage, infrav1.IBMPowerVSImageReadyCondition, infrav1.ImageImportedV1Beta2Condition, metav1.ConditionTrue, infrav1.IBMPowerVSImageReadyReason, clusterv1.ConditionSeverityInfo, "")

		case infrav1.PowerVSImageStateFailed:
			r.markCondition(imageScope.IBMPowerVSImage, infrav1.IBMPowerVSImageReadyCondition, infrav1.ImageImportedV1Beta2Condition, metav1.ConditionFalse, infrav1.IBMPowerVSImageImportFailedReason, clusterv1.ConditionSeverityError, job.Status.Message)
			return ctrl.Result{RequeueAfter: 2 * time.Minute}, fmt.Errorf("failed to import image: %s", job.Status.Message)

		case infrav1.PowerVSImageStateQueued:
			r.markCondition(imageScope.IBMPowerVSImage, infrav1.IBMPowerVSImageReadyCondition, infrav1.ImageImportedV1Beta2Condition, metav1.ConditionFalse, infrav1.IBMPowerVSImageQueuedReason, clusterv1.ConditionSeverityInfo, job.Status.Message)
			return ctrl.Result{RequeueAfter: 2 * time.Minute}, nil

		default: // Importing
			imageScope.SetImageState(string(infrav1.PowerVSImageStateImporting))
			r.markCondition(imageScope.IBMPowerVSImage, infrav1.IBMPowerVSImageReadyCondition, infrav1.ImageImportedV1Beta2Condition, metav1.ConditionFalse, infrav1.IBMPowerVSImageNotReadyReason, clusterv1.ConditionSeverityInfo, job.Status.Message)
			return ctrl.Result{RequeueAfter: 2 * time.Minute}, nil
		}
	}

	// 4. Trigger Initial Import
	img, jobRef, err := imageScope.GetOrImportImage(ctx)
	if err != nil {
		log.Error(err, "Unable to import image")
		return ctrl.Result{}, fmt.Errorf("failed to reconcile Image for %s/%s: %w", imageScope.IBMPowerVSImage.Namespace, imageScope.IBMPowerVSImage.Name, err)
	}

	if jobRef != nil && jobRef.ID != nil {
		log.Info("Import job submitted", "jobID", *jobRef.ID)
		imageScope.SetJobID(*jobRef.ID)
	} else if img == nil {
		log.Info("Import job already in progress (no jobRef returned), requeuing")
	}
	return r.reconcileImage(ctx, img, imageScope)
}

func (r *IBMPowerVSImageReconciler) reconcileImage(ctx context.Context, img *models.ImageReference, imageScope *powervsscope.ImageScope) (_ ctrl.Result, reterr error) {
	log := ctrl.LoggerFrom(ctx)

	if img != nil && img.ImageID != nil {
		image, err := imageScope.IBMPowerVSClient.GetImage(ctx, *img.ImageID)
		if err != nil {
			log.Info("Unable to get image details", "imageID", *img.ImageID)
			return ctrl.Result{}, err
		}

		imageScope.SetImageID(image.ImageID)
		imageScope.SetImageState(image.State)
		log.Info("Image status updated", "imageID", imageScope.GetImageID(), "state", image.State)

		switch imageScope.GetImageState() {
		case infrav1.PowerVSImageStateQueued:
			r.markCondition(imageScope.IBMPowerVSImage, infrav1.IBMPowerVSImageReadyCondition, infrav1.ImageReadyV1Beta2Condition, metav1.ConditionFalse, infrav1.IBMPowerVSImageNotReadyReason, clusterv1.ConditionSeverityWarning, "Image is queued")
			return ctrl.Result{RequeueAfter: 1 * time.Minute}, nil

		case infrav1.PowerVSImageStateACTIVE:
			imageScope.SetImageActive()
			r.markCondition(imageScope.IBMPowerVSImage, infrav1.IBMPowerVSImageReadyCondition, infrav1.ImageReadyV1Beta2Condition, metav1.ConditionTrue, infrav1.IBMPowerVSImageReadyReason, clusterv1.ConditionSeverityInfo, "")

		default:
			log.Info("PowerVS image state is undefined", "state", image.State, "imageID", imageScope.GetImageID())
			r.markCondition(imageScope.IBMPowerVSImage, infrav1.IBMPowerVSImageReadyCondition, infrav1.ImageReadyV1Beta2Condition, metav1.ConditionUnknown, infrav1.IBMPowerVSImageReadyUnknownReason, clusterv1.ConditionSeverityInfo, fmt.Sprintf("Unknown state: %s", image.State))
		}
	}

	if !imageScope.IsImageActive() {
		log.Info("Image is not yet ready, requeue", "state", imageScope.GetImageState())
		return ctrl.Result{RequeueAfter: 1 * time.Minute}, nil
	}

	return ctrl.Result{}, nil
}

func (r *IBMPowerVSImageReconciler) reconcileDelete(ctx context.Context, scope *powervsscope.ImageScope) (_ ctrl.Result, reterr error) {
	log := ctrl.LoggerFrom(ctx)
	log.Info("Reconciling IBMPowerVSImage delete")

	// 1. Signal that deletion is in progress
	r.markCondition(scope.IBMPowerVSImage, infrav1.IBMPowerVSImageReadyCondition, infrav1.ImageReadyV1Beta2Condition, metav1.ConditionFalse, infrav1.IBMPowerVSImageDeletingReason, clusterv1.ConditionSeverityInfo, "")

	// 2. Ensure finalizer is removed only on complete success
	defer func() {
		if reterr == nil {
			log.Info("IBMPowerVSImage deleted, removing finalizer")
			controllerutil.RemoveFinalizer(scope.IBMPowerVSImage, infrav1.IBMPowerVSImageFinalizer)
		}
	}()

	// 3. Handle cases where the Image was never fully imported
	if scope.GetImageID() == "" {
		log.Info("IBMPowerVSImage ImageID is not yet set, skipping PowerVS API image deletion")

		if scope.GetJobID() != "" {
			if err := scope.DeleteImportJob(ctx); err != nil {
				log.Error(err, "Error deleting IBMPowerVSImage Import Job")
				r.markCondition(scope.IBMPowerVSImage, infrav1.IBMPowerVSImageReadyCondition, infrav1.ImageReadyV1Beta2Condition, metav1.ConditionFalse, infrav1.IBMPowerVSImageDeletingReason, clusterv1.ConditionSeverityWarning, fmt.Sprintf("Failed to delete import job: %v", err))
				return ctrl.Result{}, fmt.Errorf("error deleting IBMPowerVSImage Import Job: %w", err)
			}
		} else {
			log.Info("JobID is not yet set, skipping PowerVS API job deletion")
		}
		return ctrl.Result{}, nil
	}

	// 4. Handle actual Image deletion (respecting the DeletePolicy)
	if scope.IBMPowerVSImage.Spec.DeletePolicy != infrav1.PowerVSImageDeletePolicyRetain {
		if err := scope.DeleteImage(ctx); err != nil {
			log.Error(err, "Error deleting IBMPowerVSImage")

			// Note: Replaced the accidental "InstanceDeletingReason" with the correct Image reason
			r.markCondition(scope.IBMPowerVSImage, infrav1.IBMPowerVSImageReadyCondition, infrav1.ImageReadyV1Beta2Condition, metav1.ConditionFalse, infrav1.IBMPowerVSImageDeletingReason, clusterv1.ConditionSeverityWarning, fmt.Sprintf("failed to delete IBMPowerVSImage: %v", err))

			return ctrl.Result{}, fmt.Errorf("error deleting IBMPowerVSImage %v: %w", klog.KObj(scope.IBMPowerVSImage), err)
		}
	} else {
		log.Info("Skipping PowerVS API image deletion due to retain policy", "imageID", scope.GetImageID())
	}

	return ctrl.Result{}, nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *IBMPowerVSImageReconciler) SetupWithManager(_ context.Context, mgr ctrl.Manager) error {
	if r.ClientBuilder == nil {
		r.ClientBuilder = powervsscope.ProdClientBuilder{}
	}
	return ctrl.NewControllerManagedBy(mgr).
		For(&infrav1.IBMPowerVSImage{}).
		Complete(r)
}

func (r *IBMPowerVSImageReconciler) shouldAdopt(i infrav1.IBMPowerVSImage) bool {
	return !clusterv1util.HasOwner(i.OwnerReferences, infrav1.GroupVersion.String(), []string{ibmPowerVSClusterKind})
}

func patchIBMPowerVSImage(ctx context.Context, patchHelper *patch.Helper, ibmPowerVSImage *infrav1.IBMPowerVSImage) error {
	// Before computing ready condition, make sure that ImageReady is always set.
	// NOTE: This is required because v1beta2 conditions comply to guideline requiring conditions to be set at the
	// first reconcile.
	if c := conditions.Get(ibmPowerVSImage, infrav1.IBMPowerVSImageReadyCondition); c == nil {
		if ibmPowerVSImage.Status.ImageState == infrav1.PowerVSImageStateACTIVE {
			conditions.Set(ibmPowerVSImage, metav1.Condition{
				Type:   infrav1.IBMPowerVSImageReadyCondition,
				Status: metav1.ConditionTrue,
				Reason: infrav1.IBMPowerVSImageReadyReason,
			})
		} else {
			conditions.Set(ibmPowerVSImage, metav1.Condition{
				Type:   infrav1.IBMPowerVSImageReadyCondition,
				Status: metav1.ConditionFalse,
				Reason: infrav1.IBMPowerVSImageNotReadyReason,
			})
		}
	}

	// always update the readyCondition.
	deprecatedv1beta1conditions.SetSummary(ibmPowerVSImage,
		deprecatedv1beta1conditions.WithConditions(
			infrav1.ImageReadyV1Beta2Condition,
		),
	)

	if err := conditions.SetSummaryCondition(ibmPowerVSImage, ibmPowerVSImage, infrav1.IBMPowerVSImageReadyCondition,
		conditions.ForConditionTypes{
			infrav1.WorkspaceReadyCondition,
		},
		// Using a custom merge strategy to override reasons applied during merge.
		conditions.CustomMergeStrategy{
			MergeStrategy: conditions.DefaultMergeStrategy(
				// Use custom reasons.
				conditions.ComputeReasonFunc(conditions.GetDefaultComputeMergeReasonFunc(
					infrav1.IBMPowerVSImageNotReadyReason,
					infrav1.IBMPowerVSImageReadyUnknownReason,
					infrav1.IBMPowerVSImageReadyReason,
				)),
			),
		},
	); err != nil {
		return fmt.Errorf("failed to set %s condition: %w", infrav1.IBMPowerVSImageReadyCondition, err)
	}

	// Patch the IBMPowerVSImage resource.
	return patchHelper.Patch(ctx, ibmPowerVSImage, patch.WithOwnedConditions{Conditions: []string{
		infrav1.IBMPowerVSImageReadyCondition,
		clusterv1.PausedCondition,
		infrav1.WorkspaceReadyCondition,
	}}, patch.Clusterv1ConditionsFieldPath{statusField, deprecatedStatus, v1beta2Version, deprecatedConditionsField})
}

// markCondition safely sets both the modern and legacy conditions for the image.
func (r *IBMPowerVSImageReconciler) markCondition(image *infrav1.IBMPowerVSImage, condType string, legacyCondType clusterv1.ConditionType, status metav1.ConditionStatus, reason string, severity clusterv1.ConditionSeverity, msg string) {
	conditions.Set(image, metav1.Condition{
		Type:    condType,
		Status:  status,
		Reason:  reason,
		Message: msg,
	})

	legacyReason := reason + "V1Beta2"

	// Retain the specific reason strings used by older clients where the legacy
	// reason differs from the default "<reason>V1Beta2" pattern.
	switch reason {
	case infrav1.IBMPowerVSImageNotReadyReason:
		legacyReason = infrav1.ImageNotReadyV1Beta2Reason
	case infrav1.IBMPowerVSImageQueuedReason:
		legacyReason = string(infrav1.PowerVSImageStateQueued)
	case infrav1.IBMPowerVSImageImportFailedReason:
		legacyReason = infrav1.ImageImportFailedV1Beta2Reason
	case infrav1.IBMPowerVSImageReadyUnknownReason:
		legacyReason = infrav1.ImageStateUnknownV1Beta2Reason
	}

	switch status {
	case metav1.ConditionUnknown:
		deprecatedv1beta1conditions.MarkUnknown(image, legacyCondType, legacyReason, "%s", msg)
	case metav1.ConditionTrue:
		deprecatedv1beta1conditions.MarkTrue(image, legacyCondType)
	default:
		deprecatedv1beta1conditions.MarkFalse(image, legacyCondType, legacyReason, severity, "%s", msg)
	}
}
