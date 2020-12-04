package v1beta1

import (
	"context"
	"encoding/json"
	"net/http"

	osconfigv1 "github.com/openshift/api/config/v1"
	utilerrors "k8s.io/apimachinery/pkg/util/errors"
	"k8s.io/klog/v2"
	"sigs.k8s.io/controller-runtime/pkg/webhook/admission"
)

// machineSetValidatorHandler validates MachineSet API resources.
// implements type Handler interface.
// https://godoc.org/github.com/kubernetes-sigs/controller-runtime/pkg/webhook/admission#Handler
type machineSetValidatorHandler struct {
	*admissionHandler
}

// machineSetDefaulterHandler defaults MachineSet API resources.
// implements type Handler interface.
// https://godoc.org/github.com/kubernetes-sigs/controller-runtime/pkg/webhook/admission#Handler
type machineSetDefaulterHandler struct {
	*admissionHandler
}

// NewMachineSetValidator returns a new machineSetValidatorHandler.
func NewMachineSetValidator() (*machineSetValidatorHandler, error) {
	infra, err := getInfra()
	if err != nil {
		return nil, err
	}

	return createMachineSetValidator(infra.Status.PlatformStatus.Type, infra.Status.InfrastructureName), nil
}

func createMachineSetValidator(platform osconfigv1.PlatformType, clusterID string) *machineSetValidatorHandler {
	return &machineSetValidatorHandler{
		admissionHandler: &admissionHandler{
			clusterID:         clusterID,
			webhookOperations: getMachineValidatorOperation(platform),
		},
	}
}

// NewMachineSetDefaulter returns a new machineSetDefaulterHandler.
func NewMachineSetDefaulter() (*machineSetDefaulterHandler, error) {
	infra, err := getInfra()
	if err != nil {
		return nil, err
	}

	return createMachineSetDefaulter(infra.Status.PlatformStatus, infra.Status.InfrastructureName), nil
}

func createMachineSetDefaulter(platformStatus *osconfigv1.PlatformStatus, clusterID string) *machineSetDefaulterHandler {
	return &machineSetDefaulterHandler{
		admissionHandler: &admissionHandler{
			clusterID:         clusterID,
			webhookOperations: getMachineDefaulterOperation(platformStatus),
		},
	}
}

// Handle handles HTTP requests for admission webhook servers.
func (h *machineSetValidatorHandler) Handle(ctx context.Context, req admission.Request) admission.Response {
	ms := &MachineSet{}

	if err := h.decoder.Decode(req, ms); err != nil {
		return admission.Errored(http.StatusBadRequest, err)
	}

	klog.V(3).Infof("Validate webhook called for MachineSet: %s", ms.GetName())

	ok, warnings, errs := h.validateMachineSet(ms)
	if !ok {
		return responseWithWarnings(admission.Denied(errs.Error()), warnings)
	}

	return responseWithWarnings(admission.Allowed("MachineSet valid"), warnings)
}

// Handle handles HTTP requests for admission webhook servers.
func (h *machineSetDefaulterHandler) Handle(ctx context.Context, req admission.Request) admission.Response {
	ms := &MachineSet{}

	if err := h.decoder.Decode(req, ms); err != nil {
		return admission.Errored(http.StatusBadRequest, err)
	}

	klog.V(3).Infof("Mutate webhook called for MachineSet: %s", ms.GetName())

	ok, warnings, errs := h.defaultMachineSet(ms)
	if !ok {
		return responseWithWarnings(admission.Denied(errs.Error()), warnings)
	}

	marshaledMachineSet, err := json.Marshal(ms)
	if err != nil {
		return responseWithWarnings(admission.Errored(http.StatusInternalServerError, err), warnings)
	}
	return responseWithWarnings(admission.PatchResponseFromRaw(req.Object.Raw, marshaledMachineSet), warnings)
}

func (h *machineSetValidatorHandler) validateMachineSet(ms *MachineSet) (bool, []string, utilerrors.Aggregate) {
	var errs []error

	// Create a Machine from the MachineSet and validate the Machine template
	m := &Machine{Spec: ms.Spec.Template.Spec}
	ok, warnings, err := h.webhookOperations(m, h.clusterID)
	if !ok {
		errs = append(errs, err.Errors()...)
	}

	if len(errs) > 0 {
		return false, warnings, utilerrors.NewAggregate(errs)
	}
	return true, warnings, nil
}

func (h *machineSetDefaulterHandler) defaultMachineSet(ms *MachineSet) (bool, []string, utilerrors.Aggregate) {
	// Create a Machine from the MachineSet and default the Machine template
	m := &Machine{Spec: ms.Spec.Template.Spec}
	ok, warnings, err := h.webhookOperations(m, h.clusterID)
	if !ok {
		return false, warnings, utilerrors.NewAggregate(err.Errors())
	}

	// Restore the defaulted template
	ms.Spec.Template.Spec = m.Spec
	return true, warnings, nil
}
