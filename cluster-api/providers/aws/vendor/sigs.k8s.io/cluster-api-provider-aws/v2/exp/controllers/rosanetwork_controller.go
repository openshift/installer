/*
Copyright The Kubernetes Authors.

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
	"maps"
	"slices"
	"strconv"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	cloudformationtypes "github.com/aws/aws-sdk-go-v2/service/cloudformation/types"
	"github.com/aws/smithy-go"
	"github.com/go-logr/logr"
	rosaCFNetwork "github.com/openshift/rosa/cmd/create/network"
	rosaAWSClient "github.com/openshift/rosa/pkg/aws"
	corev1 "k8s.io/api/core/v1"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"

	expinfrav1 "sigs.k8s.io/cluster-api-provider-aws/v2/exp/api/v1beta2"
	"sigs.k8s.io/cluster-api-provider-aws/v2/pkg/cloud/scope"
	"sigs.k8s.io/cluster-api-provider-aws/v2/pkg/logger"
	clusterv1 "sigs.k8s.io/cluster-api/api/v1beta1"
	"sigs.k8s.io/cluster-api/util/conditions"
	"sigs.k8s.io/cluster-api/util/predicates"
)

// ROSANetworkReconciler reconciles a ROSANetwork object.
type ROSANetworkReconciler struct {
	client.Client
	Log              logr.Logger
	Scheme           *runtime.Scheme
	awsClient        rosaAWSClient.Client
	cfStack          *cloudformationtypes.Stack
	WatchFilterValue string
}

// +kubebuilder:rbac:groups=infrastructure.cluster.x-k8s.io,resources=rosanetworks,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=infrastructure.cluster.x-k8s.io,resources=rosanetworks/status,verbs=get;update;patch

func (r *ROSANetworkReconciler) Reconcile(ctx context.Context, req ctrl.Request) (res ctrl.Result, reterr error) {
	log := logger.FromContext(ctx)

	// Get the rosanetwork instance
	rosaNetwork := &expinfrav1.ROSANetwork{}
	if err := r.Client.Get(ctx, req.NamespacedName, rosaNetwork); err != nil {
		if apierrors.IsNotFound(err) {
			return ctrl.Result{}, nil
		}
		log.Info("error getting ROSANetwork: %w", err)
		return ctrl.Result{Requeue: true}, nil
	}

	rosaNetworkScope, err := scope.NewROSANetworkScope(scope.ROSANetworkScopeParams{
		Client:         r.Client,
		ROSANetwork:    rosaNetwork,
		ControllerName: "rosanetwork",
		Logger:         log,
	})
	if err != nil {
		return ctrl.Result{}, fmt.Errorf("failed to create rosanetwork scope: %w", err)
	}

	// Create a new AWS/CloudFormation Client using the session cache
	if r.awsClient == nil {
		session := rosaNetworkScope.Session()
		logger := rosaNetworkScope.Logger.GetLogger()
		awsClient, err := rosaAWSClient.NewClient().
			CapaLogger(&logger).
			ExternalConfig(&session).
			Build()
		if err != nil {
			return ctrl.Result{}, fmt.Errorf("failed to create AWS Client: %w", err)
		}
		r.awsClient = awsClient
	}

	// Try to fetch CF stack with a given name
	r.cfStack, err = r.awsClient.GetCFStack(ctx, rosaNetworkScope.ROSANetwork.Spec.StackName)
	if err != nil {
		var apiErr smithy.APIError // in case the stack does not exist, AWS returns ValidationError
		if errors.As(err, &apiErr) && apiErr.ErrorCode() == "ValidationError" {
			r.cfStack = nil
		} else {
			return ctrl.Result{}, fmt.Errorf("error fetching CF stack details: %w", err)
		}
	}

	// Always close the scope
	defer func() {
		if err := rosaNetworkScope.PatchObject(); err != nil {
			reterr = errors.Join(reterr, err)
		}
	}()

	if !rosaNetwork.ObjectMeta.DeletionTimestamp.IsZero() {
		// Handle deletion reconciliation loop.
		return r.reconcileDelete(ctx, rosaNetworkScope)
	}

	// Handle normal reconciliation loop.
	return r.reconcileNormal(ctx, rosaNetworkScope)
}

func (r *ROSANetworkReconciler) reconcileNormal(ctx context.Context, rosaNetScope *scope.ROSANetworkScope) (res ctrl.Result, reterr error) {
	rosaNetScope.Info("Reconciling ROSANetwork")

	if controllerutil.AddFinalizer(rosaNetScope.ROSANetwork, expinfrav1.ROSANetworkFinalizer) {
		if err := rosaNetScope.PatchObject(); err != nil {
			return ctrl.Result{}, fmt.Errorf("failed to patch ROSANetwork: %w", err)
		}
	}

	if r.cfStack == nil { // The CF stack does not exist yet
		templateBody := string(rosaCFNetwork.CloudFormationTemplateFile)
		cfParams := map[string]string{
			"AvailabilityZoneCount": strconv.Itoa(rosaNetScope.ROSANetwork.Spec.AvailabilityZoneCount),
			"Region":                rosaNetScope.ROSANetwork.Spec.Region,
			"Name":                  rosaNetScope.ROSANetwork.Spec.StackName,
			"VpcCidr":               rosaNetScope.ROSANetwork.Spec.CIDRBlock,
		}
		// Explicitly specified AZs
		for i, zone := range rosaNetScope.ROSANetwork.Spec.AvailabilityZones {
			cfParams[fmt.Sprintf("AZ%d", i)] = zone
		}

		// Call the AWS CF stack create API
		_, err := r.awsClient.CreateStackWithParamsTags(ctx, templateBody, rosaNetScope.ROSANetwork.Spec.StackName, cfParams, rosaNetScope.ROSANetwork.Spec.StackTags)
		if err != nil {
			conditions.MarkFalse(rosaNetScope.ROSANetwork,
				expinfrav1.ROSANetworkReadyCondition,
				expinfrav1.ROSANetworkFailedReason,
				clusterv1.ConditionSeverityError,
				"%s",
				err.Error())
			return ctrl.Result{}, fmt.Errorf("failed to start CF stack creation: %w", err)
		}
		conditions.MarkFalse(rosaNetScope.ROSANetwork,
			expinfrav1.ROSANetworkReadyCondition,
			expinfrav1.ROSANetworkCreatingReason,
			clusterv1.ConditionSeverityInfo,
			"")
		return ctrl.Result{}, nil
	}
	// The cloudformation stack already exists
	if err := r.updateROSANetworkResources(ctx, rosaNetScope.ROSANetwork); err != nil {
		rosaNetScope.Info("error fetching CF stack resources: %w", err)
		return ctrl.Result{RequeueAfter: time.Second * 60}, nil
	}

	switch r.cfStack.StackStatus {
	case cloudformationtypes.StackStatusCreateInProgress: // Create in progress
		// Set the reason of false ROSANetworkReadyCondition to Creating
		conditions.MarkFalse(rosaNetScope.ROSANetwork,
			expinfrav1.ROSANetworkReadyCondition,
			expinfrav1.ROSANetworkCreatingReason,
			clusterv1.ConditionSeverityInfo,
			"")
		return ctrl.Result{RequeueAfter: time.Second * 60}, nil
	case cloudformationtypes.StackStatusCreateComplete: // Create complete
		if err := r.parseSubnets(rosaNetScope.ROSANetwork); err != nil {
			return ctrl.Result{}, fmt.Errorf("parsing stack subnets failed: %w", err)
		}

		// Set the reason of true ROSANetworkReadyCondition to Created
		// We have to use conditions.Set(), since conditions.MarkTrue() does not support setting reason
		conditions.Set(rosaNetScope.ROSANetwork,
			&clusterv1.Condition{
				Type:     expinfrav1.ROSANetworkReadyCondition,
				Status:   corev1.ConditionTrue,
				Reason:   expinfrav1.ROSANetworkCreatedReason,
				Severity: clusterv1.ConditionSeverityInfo,
			})
		return ctrl.Result{}, nil
	case cloudformationtypes.StackStatusCreateFailed: // Create failed
		// Set the reason of false ROSANetworkReadyCondition to Failed
		conditions.MarkFalse(rosaNetScope.ROSANetwork,
			expinfrav1.ROSANetworkReadyCondition,
			expinfrav1.ROSANetworkFailedReason,
			clusterv1.ConditionSeverityError,
			"")
		return ctrl.Result{}, fmt.Errorf("cloudformation stack %s creation failed, see the stack resources for more information", *r.cfStack.StackName)
	}

	return ctrl.Result{}, nil
}

func (r *ROSANetworkReconciler) reconcileDelete(ctx context.Context, rosaNetScope *scope.ROSANetworkScope) (res ctrl.Result, reterr error) {
	rosaNetScope.Info("Reconciling ROSANetwork delete")

	if r.cfStack != nil { // The CF stack still exists
		if err := r.updateROSANetworkResources(ctx, rosaNetScope.ROSANetwork); err != nil {
			rosaNetScope.Info("error fetching CF stack resources: %w", err)
			return ctrl.Result{RequeueAfter: time.Second * 60}, nil
		}

		switch r.cfStack.StackStatus {
		case cloudformationtypes.StackStatusDeleteInProgress: // Deletion in progress
			return ctrl.Result{RequeueAfter: time.Second * 60}, nil
		case cloudformationtypes.StackStatusDeleteFailed: // Deletion failed
			conditions.MarkFalse(rosaNetScope.ROSANetwork,
				expinfrav1.ROSANetworkReadyCondition,
				expinfrav1.ROSANetworkDeletionFailedReason,
				clusterv1.ConditionSeverityError,
				"")
			return ctrl.Result{}, fmt.Errorf("CF stack deletion failed")
		default: // All the other states
			err := r.awsClient.DeleteCFStack(ctx, rosaNetScope.ROSANetwork.Spec.StackName)
			if err != nil {
				conditions.MarkFalse(rosaNetScope.ROSANetwork,
					expinfrav1.ROSANetworkReadyCondition,
					expinfrav1.ROSANetworkDeletionFailedReason,
					clusterv1.ConditionSeverityError,
					"%s",
					err.Error())
				return ctrl.Result{}, fmt.Errorf("failed to start CF stack deletion: %w", err)
			}
			conditions.MarkFalse(rosaNetScope.ROSANetwork,
				expinfrav1.ROSANetworkReadyCondition,
				expinfrav1.ROSANetworkDeletingReason,
				clusterv1.ConditionSeverityInfo,
				"")
			return ctrl.Result{RequeueAfter: time.Second * 60}, nil
		}
	} else {
		controllerutil.RemoveFinalizer(rosaNetScope.ROSANetwork, expinfrav1.ROSANetworkFinalizer)
	}

	return ctrl.Result{}, nil
}

func (r *ROSANetworkReconciler) updateROSANetworkResources(ctx context.Context, rosaNet *expinfrav1.ROSANetwork) error {
	resources, err := r.awsClient.DescribeCFStackResources(ctx, rosaNet.Spec.StackName)
	if err != nil {
		return fmt.Errorf("error calling AWS DescribeStackResources(): %w", err)
	}

	rosaNet.Status.Resources = make([]expinfrav1.CFResource, len(*resources))
	for i, resource := range *resources {
		rosaNet.Status.Resources[i] = expinfrav1.CFResource{
			LogicalID:    aws.ToString(resource.LogicalResourceId),
			PhysicalID:   aws.ToString(resource.PhysicalResourceId),
			ResourceType: aws.ToString(resource.ResourceType),
			Status:       string(resource.ResourceStatus),
			Reason:       aws.ToString(resource.ResourceStatusReason),
		}
	}

	return nil
}

func (r *ROSANetworkReconciler) parseSubnets(rosaNet *expinfrav1.ROSANetwork) error {
	subnets := make(map[string]expinfrav1.ROSANetworkSubnet)

	for _, resource := range rosaNet.Status.Resources {
		if resource.ResourceType != "AWS::EC2::Subnet" {
			continue
		}

		az, err := r.awsClient.GetSubnetAvailabilityZone(resource.PhysicalID)
		if err != nil {
			return fmt.Errorf("failed to get AZ for subnet %s: %w", resource.PhysicalID, err)
		}

		subnet := subnets[az]
		subnet.AvailabilityZone = az

		if strings.HasPrefix(resource.LogicalID, "SubnetPrivate") {
			subnet.PrivateSubnet = resource.PhysicalID
		} else {
			subnet.PublicSubnet = resource.PhysicalID
		}

		subnets[az] = subnet
	}

	rosaNet.Status.Subnets = slices.Collect(maps.Values(subnets))

	return nil
}

// SetupWithManager is used to setup the controller.
func (r *ROSANetworkReconciler) SetupWithManager(ctx context.Context, mgr ctrl.Manager, options controller.Options) error {
	log := logger.FromContext(ctx)
	return ctrl.NewControllerManagedBy(mgr).
		WithOptions(options).
		WithEventFilter(predicates.ResourceHasFilterLabel(mgr.GetScheme(), log.GetLogger(), r.WatchFilterValue)).
		For(&expinfrav1.ROSANetwork{}).
		Complete(r)
}
