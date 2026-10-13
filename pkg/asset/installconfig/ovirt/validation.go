package ovirt

import (
	"fmt"

	"github.com/AlecAivazis/survey/v2"
	ovirtsdk "github.com/ovirt/go-ovirt"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"k8s.io/apimachinery/pkg/util/validation/field"

	"github.com/openshift/installer/pkg/types"
	"github.com/openshift/installer/pkg/types/ovirt"
	"github.com/openshift/installer/pkg/types/ovirt/validation"
)

// Validate executes ovirt specific validation
func Validate(ic *types.InstallConfig) error {
	allErrs := field.ErrorList{}
	ovirtPlatformPath := field.NewPath("platform", "ovirt")

	if ic.Platform.Ovirt == nil {
		return errors.New(field.Required(
			ovirtPlatformPath,
			"validation requires a Engine platform configuration").Error())
	}

	allErrs = append(
		allErrs,
		validation.ValidatePlatform(ic.Platform.Ovirt, ovirtPlatformPath, ic)...)

	con, err := NewConnection()
	if err != nil {
		return err
	}
	defer con.Close()

	if err := validateVNICProfile(*ic.Ovirt, con); err != nil {
		allErrs = append(
			allErrs,
			field.Invalid(ovirtPlatformPath.Child("vnicProfileID"), ic.Ovirt.VNICProfileID, err.Error()))
	}
	if ic.ControlPlane != nil && ic.ControlPlane.Platform.Ovirt != nil {
		allErrs = append(
			allErrs,
			validateMachinePool(con, field.NewPath("controlPlane", "platform", "ovirt"), ic.ControlPlane.Platform.Ovirt, *ic.Ovirt)...)
	}
	for idx, compute := range ic.Compute {
		fldPath := field.NewPath("compute").Index(idx)
		if compute.Platform.Ovirt != nil {
			allErrs = append(
				allErrs,
				validateMachinePool(con, fldPath.Child("platform", "ovirt"), compute.Platform.Ovirt, *ic.Ovirt)...)
		}
	}

	return allErrs.ToAggregate()
}

func validateMachinePool(con *ovirtsdk.Connection, child *field.Path, pool *ovirt.MachinePool, platform ovirt.Platform) field.ErrorList {
	allErrs := field.ErrorList{}
	allErrs = append(allErrs, validateInstanceTypeID(con, child, pool)...)
	allErrs = append(allErrs, validateMachineAffinityGroups(con, child, pool, platform)...)

	return allErrs
}

// validateExistingAffinityGroup checks that there is no affinity group with the same name in the cluster
func validateExistingAffinityGroup(con *ovirtsdk.Connection, platform ovirt.Platform) error {
	res, err := con.SystemService().ClustersService().
		ClusterService(platform.ClusterID).AffinityGroupsService().List().Send()
	if err != nil {
		return errors.Errorf("failed listing affinity groups for cluster %v", platform.ClusterID)
	}
	for _, ag := range res.MustGroups().Slice() {
		for _, agNew := range platform.AffinityGroups {
			if ag.MustName() == agNew.Name {
				return errors.Errorf(
					"affinity group %v already exist in cluster %v", agNew.Name, platform.ClusterID)
			}
		}
	}
	return nil
}

func validateClusterResources(con *ovirtsdk.Connection, ic *types.InstallConfig) error {
	mAgReplicas := make(map[string]int)
	for _, agn := range ic.ControlPlane.Platform.Ovirt.AffinityGroupsNames {
		mAgReplicas[agn] = mAgReplicas[agn] + int(*ic.ControlPlane.Replicas)
	}
	for _, compute := range ic.Compute {
		for _, agn := range compute.Platform.Ovirt.AffinityGroupsNames {
			mAgReplicas[agn] = mAgReplicas[agn] + int(*compute.Replicas)
		}
	}

	clusterName, err := GetClusterName(con, ic.Ovirt.ClusterID)
	if err != nil {
		return err
	}
	hosts, err := FindHostsInCluster(con, clusterName)
	if err != nil {
		return err
	}
	for _, ag := range ic.Ovirt.AffinityGroups {
		if _, found := mAgReplicas[ag.Name]; found {
			if len(hosts) < mAgReplicas[ag.Name] {
				msg := fmt.Sprintf("Affinity Group %v cannot be fulfilled, oVirt cluster doesn't"+
					"have enough hosts: found %v hosts but %v replicas assigned to affinity group",
					ag.Name, len(hosts), mAgReplicas[ag.Name])
				if ag.Enforcing {
					return fmt.Errorf(msg, ag)
				}
				logrus.Warning(msg)
			}
		}
	}
	return nil
}

func validateInstanceTypeID(con *ovirtsdk.Connection, child *field.Path, machinePool *ovirt.MachinePool) field.ErrorList {
	allErrs := field.ErrorList{}
	if machinePool.InstanceTypeID != "" {
		_, err := con.SystemService().InstanceTypesService().InstanceTypeService(machinePool.InstanceTypeID).Get().Send()
		if err != nil {
			allErrs = append(allErrs, field.NotFound(child.Child("instanceTypeID"), machinePool.InstanceTypeID))
		}
	}
	return allErrs
}

// validateMachineAffinityGroups checks that the affinity groups on the machine object exist in the oVirt cluster
// or created by the installer.
func validateMachineAffinityGroups(con *ovirtsdk.Connection, child *field.Path, machinePool *ovirt.MachinePool, platform ovirt.Platform) field.ErrorList {
	allErrs := field.ErrorList{}
	existingAG := make(map[string]int)

	res, err := con.SystemService().ClustersService().
		ClusterService(platform.ClusterID).AffinityGroupsService().List().Send()
	if err != nil {
		return append(
			allErrs,
			field.InternalError(
				child,
				errors.Errorf("failed listing affinity groups for cluster %v", platform.ClusterID)))
	}
	for _, ag := range res.MustGroups().Slice() {
		existingAG[ag.MustName()] = 0
	}
	// add affinity groups the installer creates
	for _, ag := range platform.AffinityGroups {
		existingAG[ag.Name] = 0
	}
	for _, ag := range machinePool.AffinityGroupsNames {
		if _, ok := existingAG[ag]; !ok {
			allErrs = append(
				allErrs,
				field.Invalid(
					child.Child("affinityGroupsNames"), ag,
					fmt.Sprintf(
						"Affinity Group %v doesn't exist in oVirt cluster or created by the installer", ag)))
		}
	}
	return allErrs
}

// authenticated takes an ovirt platform and validates
// its connection to the API by establishing
// the connection and authenticating successfully.
// The API connection is closed in the end and must leak
// or be reused in any way.
func authenticated(c *Config) survey.Validator {
	return func(val interface{}) error {
		connection, err := ovirtsdk.NewConnectionBuilder().
			URL(c.URL).
			Username(c.Username).
			Password(fmt.Sprint(val)).
			CAFile(c.CAFile).
			CACert([]byte(c.CABundle)).
			Insecure(c.Insecure).
			Build()

		if err != nil {
			return errors.Errorf("failed to construct connection to Engine platform %s", err)
		}

		defer connection.Close()

		err = connection.Test()
		if err != nil {
			return errors.Errorf("failed to connect to the Engine platform %s", err)
		}
		return nil
	}
}

// validate the provided vnic profile exists and belongs the the cluster network
func validateVNICProfile(platform ovirt.Platform, con *ovirtsdk.Connection) error {
	if platform.VNICProfileID != "" {
		profiles, err := FetchVNICProfileByClusterNetwork(con, platform.ClusterID, platform.NetworkName)
		if err != nil {
			return err
		}

		for _, p := range profiles {
			if platform.VNICProfileID == p.MustId() {
				return nil
			}
		}

		return fmt.Errorf(
			"vNic profile ID %s does not belong to cluster network %s",
			platform.VNICProfileID,
			platform.NetworkName)
	}
	return nil
}

// ValidateForProvisioning validates that the install config is valid for provisioning the cluster.
func ValidateForProvisioning(ic *types.InstallConfig) error {
	con, err := NewConnection()
	if err != nil {
		return err
	}
	defer con.Close()
	if err := validateClusterResources(con, ic); err != nil {
		return err
	}
	if err := validateExistingAffinityGroup(con, *ic.Ovirt); err != nil {
		return err
	}
	return nil
}
