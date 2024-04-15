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

package roleassignments

import (
	"context"

	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/authorization/armauthorization/v2"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/compute/armcompute/v5"
	"github.com/pkg/errors"
	"sigs.k8s.io/cluster-api-provider-azure/azure"
	"sigs.k8s.io/cluster-api-provider-azure/azure/services/async"
	"sigs.k8s.io/cluster-api-provider-azure/azure/services/scalesets"
	"sigs.k8s.io/cluster-api-provider-azure/azure/services/virtualmachines"
	"sigs.k8s.io/cluster-api-provider-azure/util/tele"
)

const serviceName = "roleassignments"

// RoleAssignmentScope defines the scope interface for a role assignment service.
type RoleAssignmentScope interface {
	azure.AsyncStatusUpdater
	azure.Authorizer
	RoleAssignmentSpecs(principalID *string) []azure.ResourceSpecGetter
	HasSystemAssignedIdentity() bool
	RoleAssignmentResourceType() string
	Name() string
	ResourceGroup() string
}

// Service provides operations on Azure resources.
type Service struct {
	Scope                 RoleAssignmentScope
	virtualMachinesGetter async.Getter
	async.Reconciler
	virtualMachineScaleSetGetter async.Getter
}

// New creates a new service.
func New(scope RoleAssignmentScope) (*Service, error) {
	client, err := newClient(scope)
	if err != nil {
		return nil, errors.Wrap(err, "failed to create roleassignments service")
	}
	scaleSetsClient, err := scalesets.NewClient(scope, scope.DefaultedAzureCallTimeout())
	if err != nil {
		return nil, errors.Wrap(err, "failed to create scalesets service")
	}
	virtualMachinesClient, err := virtualmachines.NewClient(scope, scope.DefaultedAzureCallTimeout())
	if err != nil {
		return nil, errors.Wrap(err, "failed to create virtualmachines service")
	}
	return &Service{
		Scope:                        scope,
		virtualMachinesGetter:        virtualMachinesClient,
		virtualMachineScaleSetGetter: scaleSetsClient,
		Reconciler: async.New[armauthorization.RoleAssignmentsClientCreateResponse,
			armauthorization.RoleAssignmentsClientDeleteResponse](scope, client, nil),
	}, nil
}

// Name returns the service name.
func (s *Service) Name() string {
	return serviceName
}

// Reconcile idempotently creates or updates a role assignment.
func (s *Service) Reconcile(ctx context.Context) error {
	ctx, log, done := tele.StartSpanWithLogger(ctx, "roleassignments.Service.Reconcile")
	defer done()

	ctx, cancel := context.WithTimeout(ctx, s.Scope.DefaultedAzureServiceReconcileTimeout())
	defer cancel()

	log.V(2).Info("reconciling role assignment")

	// Return early if the identity is not system assigned as there will be no
	// role assignment spec in this case.
	if !s.Scope.HasSystemAssignedIdentity() {
		log.V(2).Info("no role assignment spec to reconcile")
		return nil
	}

	var principalID *string
	resourceType := s.Scope.RoleAssignmentResourceType()
	switch resourceType {
	case azure.VirtualMachine:
		ID, err := s.getVMPrincipalID(ctx)
		if err != nil {
			return errors.Wrap(err, "failed to assign role to system assigned identity")
		}
		principalID = ID
	case azure.VirtualMachineScaleSet:
		ID, err := s.getVMSSPrincipalID(ctx)
		if err != nil {
			return errors.Wrap(err, "failed to assign role to system assigned identity")
		}
		principalID = ID
	default:
		return errors.Errorf("unexpected resource type %q. Expected one of [%s, %s]", resourceType,
			azure.VirtualMachine, azure.VirtualMachineScaleSet)
	}

	for _, roleAssignmentSpec := range s.Scope.RoleAssignmentSpecs(principalID) {
		log.V(2).Info("Creating role assignment")
		if roleAssignmentSpec.ResourceName() == "" {
			log.V(2).Info("RoleAssignmentName is empty. This is not expected and will cause this System Assigned Identity to have no permissions.")
		}
		_, err := s.CreateOrUpdateResource(ctx, roleAssignmentSpec, serviceName)
		if err != nil {
			return errors.Wrapf(err, "cannot assign role to %s system assigned identity", resourceType)
		}
	}

	return nil
}

// getVMPrincipalID returns the VM principal ID.
func (s *Service) getVMPrincipalID(ctx context.Context) (*string, error) {
	ctx, log, done := tele.StartSpanWithLogger(ctx, "roleassignments.Service.getVMPrincipalID")
	defer done()
	log.V(2).Info("fetching principal ID for VM")
	spec := &virtualmachines.VMSpec{
		Name:          s.Scope.Name(),
		ResourceGroup: s.Scope.ResourceGroup(),
	}

	resultVMIface, err := s.virtualMachinesGetter.Get(ctx, spec)
	if err != nil {
		return nil, errors.Wrap(err, "failed to get principal ID for VM")
	}
	resultVM, ok := resultVMIface.(armcompute.VirtualMachine)
	if !ok {
		return nil, errors.Errorf("%T is not an armcompute.VirtualMachine", resultVMIface)
	}
	return resultVM.Identity.PrincipalID, nil
}

// getVMSSPrincipalID returns the VMSS principal ID.
func (s *Service) getVMSSPrincipalID(ctx context.Context) (*string, error) {
	ctx, log, done := tele.StartSpanWithLogger(ctx, "roleassignments.Service.getVMPrincipalID")
	defer done()
	log.V(2).Info("fetching principal ID for VMSS")
	spec := &scalesets.ScaleSetSpec{
		Name:          s.Scope.Name(),
		ResourceGroup: s.Scope.ResourceGroup(),
	}

	resultVMSSIface, err := s.virtualMachineScaleSetGetter.Get(ctx, spec)
	if err != nil {
		return nil, errors.Wrap(err, "failed to get principal ID for VMSS")
	}
	resultVMSS, ok := resultVMSSIface.(armcompute.VirtualMachineScaleSet)
	if !ok {
		return nil, errors.Errorf("%T is not an armcompute.VirtualMachineScaleSet", resultVMSSIface)
	}

	return resultVMSS.Identity.PrincipalID, nil
}

// Delete is a no-op as the role assignments get deleted as part of VM deletion.
func (s *Service) Delete(ctx context.Context) error {
	_, _, done := tele.StartSpanWithLogger(ctx, "roleassignments.Service.Delete")
	defer done()
	return nil
}

// IsManaged returns always returns true as CAPZ does not support BYO role assignments.
func (s *Service) IsManaged(ctx context.Context) (bool, error) {
	return true, nil
}
