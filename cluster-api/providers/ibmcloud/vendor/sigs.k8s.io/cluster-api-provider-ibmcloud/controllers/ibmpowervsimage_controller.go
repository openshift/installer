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

package controllers

import (
	"context"
	"fmt"
	"time"

	"github.com/IBM-Cloud/power-go-client/power/models"

	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/tools/record"
	"k8s.io/klog/v2"

	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"

	capiv1beta1 "sigs.k8s.io/cluster-api/api/v1beta1"
	clusterv1util "sigs.k8s.io/cluster-api/util"
	"sigs.k8s.io/cluster-api/util/conditions"

	infrav1beta2 "sigs.k8s.io/cluster-api-provider-ibmcloud/api/v1beta2"
	"sigs.k8s.io/cluster-api-provider-ibmcloud/cloud/scope"
	"sigs.k8s.io/cluster-api-provider-ibmcloud/pkg/endpoints"
)

// IBMPowerVSImageReconciler reconciles a IBMPowerVSImage object.
type IBMPowerVSImageReconciler struct {
	client.Client
	Recorder        record.EventRecorder
	ServiceEndpoint []endpoints.ServiceEndpoint
	Scheme          *runtime.Scheme
}

//+kubebuilder:rbac:groups=infrastructure.cluster.x-k8s.io,resources=ibmpowervsimages,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=infrastructure.cluster.x-k8s.io,resources=ibmpowervsimages/status,verbs=get;update;patch

// Reconcile implements controller runtime Reconciler interface and handles reconciliation logic for IBMPowerVSImage.
func (r *IBMPowerVSImageReconciler) Reconcile(ctx context.Context, req ctrl.Request) (_ ctrl.Result, reterr error) {
	log := ctrl.LoggerFrom(ctx)

	ibmImage := &infrav1beta2.IBMPowerVSImage{}
	err := r.Get(ctx, req.NamespacedName, ibmImage)
	if err != nil {
		if apierrors.IsNotFound(err) {
			return ctrl.Result{}, nil
		}
		return ctrl.Result{}, err
	}

	var cluster *infrav1beta2.IBMPowerVSCluster
	scopeParams := scope.PowerVSImageScopeParams{
		Client:          r.Client,
		Logger:          log,
		IBMPowerVSImage: ibmImage,
		ServiceEndpoint: r.ServiceEndpoint,
	}

	// Externally managed clusters might not be available during image deletion. Get the cluster only when image is still not deleted.
	if ibmImage.DeletionTimestamp.IsZero() {
		cluster, err = scope.GetClusterByName(ctx, r.Client, ibmImage.Namespace, ibmImage.Spec.ClusterName)
		if err != nil {
			return ctrl.Result{}, err
		}
		scopeParams.Zone = cluster.Spec.Zone
	}

	// Create the scope
	imageScope, err := scope.NewPowerVSImageScope(scopeParams)
	if err != nil {
		return ctrl.Result{}, fmt.Errorf("failed to create scope: %w", err)
	}

	// Always close the scope when exiting this function so we can persist any IBMPowerVSImage changes.
	defer func() {
		if imageScope != nil {
			if err := imageScope.Close(); err != nil && reterr == nil {
				reterr = err
			}
		}
	}()

	// Handle deleted clusters.
	if !ibmImage.DeletionTimestamp.IsZero() {
		return r.reconcileDelete(imageScope)
	}

	return r.reconcile(cluster, imageScope)
}

func (r *IBMPowerVSImageReconciler) reconcile(cluster *infrav1beta2.IBMPowerVSCluster, imageScope *scope.PowerVSImageScope) (ctrl.Result, error) {
	if controllerutil.AddFinalizer(imageScope.IBMPowerVSImage, infrav1beta2.IBMPowerVSImageFinalizer) {
		return ctrl.Result{}, nil
	}

	// Create new labels section for IBMPowerVSImage metadata if nil.
	if imageScope.IBMPowerVSImage.Labels == nil {
		imageScope.IBMPowerVSImage.Labels = make(map[string]string)
	}

	if _, ok := imageScope.IBMPowerVSImage.Labels[capiv1beta1.ClusterNameLabel]; !ok {
		imageScope.IBMPowerVSImage.Labels[capiv1beta1.ClusterNameLabel] = imageScope.IBMPowerVSImage.Spec.ClusterName
	}

	if r.shouldAdopt(*imageScope.IBMPowerVSImage) {
		imageScope.Info("Image Controller has not yet set OwnerRef")
		imageScope.IBMPowerVSImage.OwnerReferences = clusterv1util.EnsureOwnerRef(imageScope.IBMPowerVSImage.OwnerReferences, metav1.OwnerReference{
			APIVersion: infrav1beta2.GroupVersion.String(),
			Kind:       "IBMPowerVSCluster",
			Name:       cluster.Name,
			UID:        cluster.UID,
		})
		return ctrl.Result{}, nil
	}

	if jobID := imageScope.GetJobID(); jobID != "" {
		job, err := imageScope.IBMPowerVSClient.GetJob(jobID)
		if err != nil {
			imageScope.Info("Unable to get job details")
			return ctrl.Result{RequeueAfter: 2 * time.Minute}, err
		}
		switch *job.Status.State {
		case "completed":
			conditions.MarkTrue(imageScope.IBMPowerVSImage, infrav1beta2.ImageImportedCondition)
		case "failed":
			imageScope.SetNotReady()
			imageScope.SetImageState(string(infrav1beta2.PowerVSImageStateFailed))
			conditions.MarkFalse(imageScope.IBMPowerVSImage, infrav1beta2.ImageImportedCondition, infrav1beta2.ImageImportFailedReason, capiv1beta1.ConditionSeverityError, "%s", job.Status.Message)
			return ctrl.Result{RequeueAfter: 2 * time.Minute}, fmt.Errorf("failed to import image, message: %s", job.Status.Message)
		case "queued":
			imageScope.SetNotReady()
			imageScope.SetImageState(string(infrav1beta2.PowerVSImageStateQue))
			conditions.MarkFalse(imageScope.IBMPowerVSImage, infrav1beta2.ImageImportedCondition, string(infrav1beta2.PowerVSImageStateQue), capiv1beta1.ConditionSeverityInfo, "%s", job.Status.Message)
			return ctrl.Result{RequeueAfter: 2 * time.Minute}, nil
		default:
			imageScope.SetNotReady()
			imageScope.SetImageState(string(infrav1beta2.PowerVSImageStateImporting))
			conditions.MarkFalse(imageScope.IBMPowerVSImage, infrav1beta2.ImageImportedCondition, *job.Status.State, capiv1beta1.ConditionSeverityInfo, "%s", job.Status.Message)
			return ctrl.Result{RequeueAfter: 2 * time.Minute}, nil
		}
	}

	img, jobRef, err := r.getOrCreate(imageScope)
	if err != nil {
		imageScope.Error(err, "Unable to import image")
		return ctrl.Result{}, fmt.Errorf("failed to reconcile Image for IBMPowerVSImage %s/%s: %w", imageScope.IBMPowerVSImage.Namespace, imageScope.IBMPowerVSImage.Name, err)
	}

	if jobRef != nil {
		imageScope.SetJobID(*jobRef.ID)
	}
	return reconcileImage(img, imageScope)
}

func reconcileImage(img *models.ImageReference, imageScope *scope.PowerVSImageScope) (_ ctrl.Result, reterr error) {
	if img != nil {
		image, err := imageScope.IBMPowerVSClient.GetImage(*img.ImageID)
		if err != nil {
			imageScope.Info("Unable to get image details")
			return ctrl.Result{}, err
		}

		imageScope.SetImageID(image.ImageID)
		imageScope.Info("ImageID", imageScope.GetImageID())
		imageScope.SetImageState(image.State)
		imageScope.Info("ImageState", image.State)

		switch imageScope.GetImageState() {
		case infrav1beta2.PowerVSImageStateQue:
			imageScope.Info("Image is in queued state")
			imageScope.SetNotReady()
			conditions.MarkFalse(imageScope.IBMPowerVSImage, infrav1beta2.ImageReadyCondition, infrav1beta2.ImageNotReadyReason, capiv1beta1.ConditionSeverityWarning, "")
			return ctrl.Result{RequeueAfter: 1 * time.Minute}, nil
		case infrav1beta2.PowerVSImageStateACTIVE:
			imageScope.Info("Image is in active state")
			imageScope.SetReady()
			conditions.MarkTrue(imageScope.IBMPowerVSImage, infrav1beta2.ImageReadyCondition)
		default:
			imageScope.SetNotReady()
			imageScope.Info("PowerVS image state is undefined", "state", image.State, "image-id", imageScope.GetImageID())
			conditions.MarkUnknown(imageScope.IBMPowerVSImage, infrav1beta2.ImageReadyCondition, "", "")
		}
	}

	// Requeue after 1 minute if image is not ready to update status of the image properly.
	if !imageScope.IsReady() {
		imageScope.Info("Image is not yet ready")
		return ctrl.Result{RequeueAfter: 1 * time.Minute}, nil
	}

	return ctrl.Result{}, nil
}

func (r *IBMPowerVSImageReconciler) reconcileDelete(scope *scope.PowerVSImageScope) (_ ctrl.Result, reterr error) {
	scope.Info("Handling deleted IBMPowerVSImage")

	defer func() {
		if reterr == nil {
			// Image is deleted so remove the finalizer.
			controllerutil.RemoveFinalizer(scope.IBMPowerVSImage, infrav1beta2.IBMPowerVSImageFinalizer)
		}
	}()

	if scope.GetImageID() == "" {
		scope.Info("ImageID is not yet set, hence not invoking the PowerVS API to delete the image")
		if scope.GetJobID() == "" {
			scope.Info("JobID is not yet set, hence not invoking the PowerVS API to delete the image import job")
			return ctrl.Result{}, nil
		}
		if err := scope.DeleteImportJob(); err != nil {
			scope.Error(err, "Error deleting IBMPowerVSImage Import Job")
			return ctrl.Result{}, fmt.Errorf("error deleting IBMPowerVSImage Import Job: %w", err)
		}
		return ctrl.Result{}, nil
	}

	if scope.IBMPowerVSImage.Spec.DeletePolicy != string(infrav1beta2.DeletePolicyRetain) {
		if err := scope.DeleteImage(); err != nil {
			scope.Error(err, "Error deleting IBMPowerVSImage")
			return ctrl.Result{}, fmt.Errorf("error deleting IBMPowerVSImage %v: %w", klog.KObj(scope.IBMPowerVSImage), err)
		}
	}
	return ctrl.Result{}, nil
}

func (r *IBMPowerVSImageReconciler) getOrCreate(scope *scope.PowerVSImageScope) (*models.ImageReference, *models.JobReference, error) {
	image, job, err := scope.CreateImageCOSBucket()
	return image, job, err
}

func (r *IBMPowerVSImageReconciler) shouldAdopt(i infrav1beta2.IBMPowerVSImage) bool {
	return !clusterv1util.HasOwner(i.OwnerReferences, infrav1beta2.GroupVersion.String(), []string{"IBMPowerVSCluster"})
}

// SetupWithManager sets up the controller with the Manager.
func (r *IBMPowerVSImageReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&infrav1beta2.IBMPowerVSImage{}).
		Complete(r)
}
