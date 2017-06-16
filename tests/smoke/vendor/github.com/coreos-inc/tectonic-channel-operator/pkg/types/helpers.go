package types

import (
	"encoding/json"
	"fmt"
	"path"
	"strings"

	"k8s.io/client-go/pkg/runtime"
)

// UpdateDesiredVersion updates the desired version of the TPR's spec.
func UpdateDesiredVersion(object *runtime.Unstructured, desiredVersion interface{}) error {
	spec, ok := object.Object[JSONNameVersionSpec]
	if !ok {
		spec = make(map[string]interface{})
		object.Object[JSONNameVersionSpec] = spec
	}

	specBody, ok := spec.(map[string]interface{})
	if !ok {
		return fmt.Errorf("expect %q to be a map, got: %T", JSONNameVersionSpec, spec)
	}

	specBody[JSONNameVersionSpecDesiredVersion] = desiredVersion
	return nil
}

// UpdateTargetVersion updates the target version of the TPR's status.
func UpdateTargetVersion(object *runtime.Unstructured, targetVersion interface{}) error {
	status, ok := object.Object[JSONNameVersionStatus]
	if !ok {
		status = make(map[string]interface{})
		object.Object[JSONNameVersionStatus] = status
	}

	statusBody, ok := status.(map[string]interface{})
	if !ok {
		return fmt.Errorf("expect %q to be a map, got: %T", JSONNameVersionStatus, status)
	}

	statusBody[JSONNameVersionStatusTargetVersion] = targetVersion
	return nil
}

// UpdateSpecPaused updates the "paused" boolean in the TPR's spec.
func UpdateSpecPaused(object *runtime.Unstructured, paused interface{}) error {
	spec, ok := object.Object[JSONNameVersionSpec]
	if !ok {
		spec = make(map[string]interface{})
		object.Object[JSONNameVersionSpec] = spec
	}

	specBody, ok := spec.(map[string]interface{})
	if !ok {
		return fmt.Errorf("expect %q to be a map, got: %T", JSONNameVersionSpec, spec)
	}

	specBody[JSONNameVersionSpecPaused] = paused
	return nil
}

// UpdateStatusPaused updates the "paused" boolean in the TPR's status.
func UpdateStatusPaused(object *runtime.Unstructured, paused interface{}) error {
	status, ok := object.Object[JSONNameVersionStatus]
	if !ok {
		status = make(map[string]interface{})
		object.Object[JSONNameVersionStatus] = status
	}

	statusBody, ok := status.(map[string]interface{})
	if !ok {
		return fmt.Errorf("expect %q to be a map, got: %T", JSONNameVersionStatus, status)
	}

	statusBody[JSONNameVersionStatusPaused] = paused
	return nil
}

// UpdateCurrentRemoveTargetVersion updates the current version of the TPR's status to be targetVersion,
// and cleanup the targetVersion.
func UpdateCurrentRemoveTargetVersion(object *runtime.Unstructured, _ interface{}) error {
	status, ok := object.Object[JSONNameVersionStatus]
	if !ok {
		status = make(map[string]interface{})
		object.Object[JSONNameVersionStatus] = status
	}

	statusBody, ok := status.(map[string]interface{})
	if !ok {
		return fmt.Errorf("expect %q to be a map, got: %T", JSONNameVersionStatus, status)
	}

	targetVersion, ok := statusBody[JSONNameVersionStatusTargetVersion]
	if !ok || targetVersion == "" {
		return fmt.Errorf("targetVersion is empty")
	}

	statusBody[JSONNameVersionStatusCurrentVersion] = targetVersion
	delete(statusBody, JSONNameVersionStatusTargetVersion)

	return nil
}

// SetUpdateTrigger sets the "triggerUpdate" in the TCO config.
func SetUpdateTrigger(object *runtime.Unstructured, value interface{}) error {
	var tcoConfig ChannelOperatorConfig

	b, err := json.Marshal(object.Object)
	if err != nil {
		return fmt.Errorf("failed to marshal the TPR: %v", err)
	}

	if err := json.Unmarshal(b, &tcoConfig); err != nil {
		return fmt.Errorf("failed to unmarshal TCO config: %v", err)
	}

	v, ok := value.(bool)
	if !ok {
		return fmt.Errorf("expect argument to be a boolean, got: %T", value)
	}

	tcoConfig.TriggerUpdate = v

	b, err = json.Marshal(&tcoConfig)
	if err != nil {
		return fmt.Errorf("failed to marshal TCO config: %v", err)
	}

	if err := json.Unmarshal(b, object); err != nil {
		return fmt.Errorf("failed to unmarshal the TPR: %v", err)
	}

	return nil
}

// SetUpdateCheckTrigger sets the "triggerUpdateCheck" in the TCO config.
func SetUpdateCheckTrigger(object *runtime.Unstructured, value interface{}) error {
	var tcoConfig ChannelOperatorConfig

	b, err := json.Marshal(object.Object)
	if err != nil {
		return fmt.Errorf("failed to marshal the TPR: %v", err)
	}

	if err := json.Unmarshal(b, &tcoConfig); err != nil {
		return fmt.Errorf("failed to unmarshal TCO config: %v", err)
	}

	v, ok := value.(bool)
	if !ok {
		return fmt.Errorf("expect argument to be a boolean, got: %T", value)
	}

	tcoConfig.TriggerUpdateCheck = v

	b, err = json.Marshal(&tcoConfig)
	if err != nil {
		return fmt.Errorf("failed to marshal TCO config: %v", err)
	}

	if err := json.Unmarshal(b, object); err != nil {
		return fmt.Errorf("failed to unmarshal the TPR: %v", err)
	}

	return nil
}

// SetAutoUpdate sets the "autoUpdate" in the TCO config.
func SetAutoUpdate(object *runtime.Unstructured, value interface{}) error {
	var tcoConfig ChannelOperatorConfig

	b, err := json.Marshal(object.Object)
	if err != nil {
		return fmt.Errorf("failed to marshal the TPR: %v", err)
	}

	if err := json.Unmarshal(b, &tcoConfig); err != nil {
		return fmt.Errorf("failed to unmarshal TCO config: %v", err)
	}

	v, ok := value.(bool)
	if !ok {
		return fmt.Errorf("expect argument to be a boolean, got: %T", value)
	}

	tcoConfig.AutoUpdate = v

	b, err = json.Marshal(&tcoConfig)
	if err != nil {
		return fmt.Errorf("failed to marshal TCO config: %v", err)
	}

	if err := json.Unmarshal(b, object); err != nil {
		return fmt.Errorf("failed to unmarshal the TPR: %v", err)
	}

	return nil
}

// CreateTaskStatuses creates the "taskStatuses" in the AppVersionStatus
// with the given task list.
func CreateTaskStatuses(object *runtime.Unstructured, value interface{}) error {
	var appVersion AppVersion

	b, err := json.Marshal(object.Object)
	if err != nil {
		return fmt.Errorf("failed to marshal the TPR: %v", err)
	}

	if err := json.Unmarshal(b, &appVersion); err != nil {
		return fmt.Errorf("failed to unmarshal AppVersion: %v", err)
	}

	statuses, ok := value.(TaskStatusList)
	if !ok {
		return fmt.Errorf("expect argument to be a boolean, got: %T", value)
	}

	appVersion.Status.TaskStatuses = statuses

	b, err = json.Marshal(&appVersion)
	if err != nil {
		return fmt.Errorf("failed to marshal AppVersion: %v", err)
	}

	if err := json.Unmarshal(b, object); err != nil {
		return fmt.Errorf("failed to unmarshal the TPR: %v", err)
	}

	return nil
}

// SetTaskStatus overwrites the update task, returns an error if the it does not exist.
func SetTaskStatus(object *runtime.Unstructured, value interface{}) error {
	var appVersion AppVersion

	b, err := json.Marshal(object.Object)
	if err != nil {
		return fmt.Errorf("failed to marshal the TPR: %v", err)
	}

	if err := json.Unmarshal(b, &appVersion); err != nil {
		return fmt.Errorf("failed to unmarshal AppVersion: %v", err)
	}

	status, ok := value.(TaskStatus)
	if !ok {
		return fmt.Errorf("expect argument to be a TaskStatus, got: %T", value)
	}

	var found bool
	for i, v := range appVersion.Status.TaskStatuses {
		if v.Name == status.Name {
			appVersion.Status.TaskStatuses[i] = status
			found = true
		}
	}
	if !found {
		return fmt.Errorf("%q is not found in TaskStatus", status.Name)
	}

	b, err = json.Marshal(&appVersion)
	if err != nil {
		return fmt.Errorf("failed to marshal AppVersion: %v", err)
	}

	if err := json.Unmarshal(b, object); err != nil {
		return fmt.Errorf("failed to unmarshal the TPR: %v", err)
	}

	return nil
}

// ResetFailureStatus fills the AppVersion.failureStatus field.
func ResetFailureStatus(object *runtime.Unstructured, value interface{}) error {
	var appVersion AppVersion

	b, err := json.Marshal(object.Object)
	if err != nil {
		return fmt.Errorf("failed to marshal the TPR: %v", err)
	}

	if err := json.Unmarshal(b, &appVersion); err != nil {
		return fmt.Errorf("failed to unmarshal AppVersion: %v", err)
	}

	appVersion.Status.FailureStatus = nil

	b, err = json.Marshal(&appVersion)
	if err != nil {
		return fmt.Errorf("failed to marshal AppVersion: %v", err)
	}

	if err := json.Unmarshal(b, object); err != nil {
		return fmt.Errorf("failed to unmarshal the TPR: %v", err)
	}

	return nil
}

// SetFailureStatus fills the AppVersion.failureStatus field.
func SetFailureStatus(object *runtime.Unstructured, value interface{}) error {
	var appVersion AppVersion

	b, err := json.Marshal(object.Object)
	if err != nil {
		return fmt.Errorf("failed to marshal the TPR: %v", err)
	}

	if err := json.Unmarshal(b, &appVersion); err != nil {
		return fmt.Errorf("failed to unmarshal AppVersion: %v", err)
	}

	status, ok := value.(FailureStatus)
	if !ok {
		return fmt.Errorf("expect argument to be a FailureStatus, got: %T", value)
	}

	appVersion.Status.FailureStatus = &status

	b, err = json.Marshal(&appVersion)
	if err != nil {
		return fmt.Errorf("failed to marshal AppVersion: %v", err)
	}

	if err := json.Unmarshal(b, object); err != nil {
		return fmt.Errorf("failed to unmarshal the TPR: %v", err)
	}

	return nil
}

// VersionToName creates a name for the TectonicVersion for the given version,
// e.g. "1.4.5" becomes "1-4-5".
func VersionToName(version string) string {
	return strings.Replace(version, ".", "-", -1)
}

// ApplyTectonicVersionMetadata makes sure the TectonicVersion's metadata is right,
// such as Kind, APIVersion, namespace, name, etc.
func ApplyTectonicVersionMetadata(spec *TectonicVersion) *TectonicVersion {
	spec.Kind = TectonicVersionTPRKind
	spec.APIVersion = path.Join(TectonicAPIGroup, TectonicVersionTPRVersion)
	spec.Namespace = TectonicNamespace
	spec.Name = VersionToName(spec.Version)

	if spec.Labels == nil {
		spec.Labels = make(map[string]string)
	}
	spec.Labels[LabelKeyManagedByChannelOperator] = LabelValueManagedByChannelOperator
	return spec
}
