package bootstrap

import (
	"fmt"

	"github.com/sirupsen/logrus"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/validation/field"
	"k8s.io/utils/ptr"
	"sigs.k8s.io/yaml"

	configv1 "github.com/openshift/api/config/v1"
	"github.com/openshift/installer/pkg/asset"
	"github.com/openshift/installer/pkg/asset/manifests"
	"github.com/openshift/installer/pkg/types"
)

// validateTopology checks the manifests consumed by ignition generation, including
// user edits made after create manifests, against the configured replica counts
// and the effective scheduler settings.
func validateTopology(files []*asset.File, installConfig *types.InstallConfig) error {
	// A missing scheduler manifest uses the API default, mastersSchedulable=false.
	var scheduler configv1.Scheduler
	schedulerFilename := "the default Scheduler configuration"
	// Both topology fields default to HighlyAvailable when omitted in the API.
	infrastructure := configv1.Infrastructure{
		Status: configv1.InfrastructureStatus{
			ControlPlaneTopology:   configv1.HighlyAvailableTopologyMode,
			InfrastructureTopology: configv1.HighlyAvailableTopologyMode,
		},
	}
	var infrastructureFilename string
	for _, file := range files {
		// Identify resources by their metadata so renamed manifests are checked too.
		var metadata metav1.PartialObjectMetadata
		if err := yaml.Unmarshal(file.Data, &metadata); err != nil {
			return fmt.Errorf("failed to unmarshal %s: %w", file.Filename, err)
		}
		if metadata.APIVersion != configv1.SchemeGroupVersion.String() || metadata.Name != "cluster" || metadata.Namespace != "" {
			continue
		}

		switch metadata.Kind {
		case "Scheduler":
			if err := yaml.Unmarshal(file.Data, &scheduler); err != nil {
				return fmt.Errorf("failed to unmarshal %s: %w", file.Filename, err)
			}
			schedulerFilename = file.Filename
		case "Infrastructure":
			if err := yaml.Unmarshal(file.Data, &infrastructure); err != nil {
				return fmt.Errorf("failed to unmarshal %s: %w", file.Filename, err)
			}
			infrastructureFilename = file.Filename
		}
	}

	if infrastructureFilename == "" {
		return fmt.Errorf("missing Infrastructure cluster manifest; topology cannot be validated before creating ignition configs")
	}

	errs := field.ErrorList{}
	expectedControlPlane, _ := manifests.DetermineTopologies(installConfig)
	if infrastructure.Status.ControlPlaneTopology != expectedControlPlane {
		errs = append(errs, field.Invalid(field.NewPath("status", "controlPlaneTopology"), infrastructure.Status.ControlPlaneTopology,
			fmt.Sprintf("must be %s for the control-plane and arbiter replicas in install-config", expectedControlPlane)))
	}

	infraPath := field.NewPath("status", "infrastructureTopology")
	switch infrastructure.Status.InfrastructureTopology {
	case configv1.SingleReplicaTopologyMode, configv1.HighlyAvailableTopologyMode:
		workerReplicas := int64(0)
		for _, pool := range installConfig.Compute {
			workerReplicas += ptr.Deref(pool.Replicas, 0)
		}
		schedulableReplicas := workerReplicas
		if scheduler.Spec.MastersSchedulable {
			schedulableReplicas += ptr.Deref(installConfig.ControlPlane.Replicas, 3)
		}
		// Zero configured workers can mean UPI, with workers provisioned outside
		// the installer. Only reject SingleReplica when multiple schedulable
		// nodes are known; zero or one does not establish the actual UPI count.
		// Arbiters do not provide infrastructure workload capacity.
		if infrastructure.Status.InfrastructureTopology == configv1.SingleReplicaTopologyMode && schedulableReplicas > 1 {
			errs = append(errs, field.Invalid(infraPath, infrastructure.Status.InfrastructureTopology,
				fmt.Sprintf("must be %s for %d configured workers and spec.mastersSchedulable=%t in %s; correct the infrastructure topology or scheduler settings before creating ignition configs",
					configv1.HighlyAvailableTopologyMode, workerReplicas, scheduler.Spec.MastersSchedulable, schedulerFilename)))
		} else if infrastructure.Status.InfrastructureTopology == configv1.HighlyAvailableTopologyMode && workerReplicas > 0 && schedulableReplicas == 1 {
			// Allow HA with a single schedulable node so users can plan for scale-up.
			logrus.Warnf("status.infrastructureTopology is HighlyAvailable in %s, but only one schedulable node is configured with spec.mastersSchedulable=false in %s; installation may not complete until additional schedulable nodes are available", infrastructureFilename, schedulerFilename)
		}
	default:
		errs = append(errs, field.NotSupported(infraPath, infrastructure.Status.InfrastructureTopology,
			[]string{string(configv1.SingleReplicaTopologyMode), string(configv1.HighlyAvailableTopologyMode)}))
	}
	if err := errs.ToAggregate(); err != nil {
		return fmt.Errorf("incompatible topology in %s: %w", infrastructureFilename, err)
	}
	return nil
}
