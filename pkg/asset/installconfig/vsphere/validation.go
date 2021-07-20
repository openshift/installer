package vsphere

import (
	"context"
	"crypto/x509"
	"fmt"
	"strconv"
	"time"

	"github.com/hashicorp/go-version"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"github.com/vmware/govmomi/find"
	"github.com/vmware/govmomi/vim25"
	"github.com/vmware/govmomi/vim25/mo"
	"k8s.io/apimachinery/pkg/util/validation/field"

	"github.com/openshift/installer/pkg/types"
	vspheretypes "github.com/openshift/installer/pkg/types/vsphere"
	"github.com/openshift/installer/pkg/types/vsphere/validation"
)

const (
	esxi67U3BuildNumber int = 14320388
)

// Validate executes platform-specific validation.
func Validate(ic *types.InstallConfig, certificates ...x509.Certificate) error {
	allErrs := field.ErrorList{}
	if ic.Platform.VSphere == nil {
		return errors.New(field.Required(field.NewPath("platform", "vsphere"), "vSphere validation requires a vSphere platform configuration").Error())
	}

	client, _, err := vspheretypes.CreateVSphereClients(context.Background(), ic.VSphere.VCenter, ic.VSphere.Username, ic.VSphere.Password, certificates...)
	if err != nil {
		return errors.Wrap(err, "unable to connect to vCenter API")
	}

	allErrs = append(allErrs, validatevCenterVersion(client, field.NewPath("platform").Child("vsphere"))...)
	allErrs = append(allErrs, validateESXiVersion(client, ic, field.NewPath("platform").Child("vsphere"))...)
	allErrs = append(allErrs, validation.ValidatePlatform(ic.Platform.VSphere, field.NewPath("platform").Child("vsphere"))...)

	return allErrs.ToAggregate()
}

// ValidateForProvisioning performs platform validation specifically for installer-
// provisioned infrastructure. In this case, self-hosted networking is a requirement
// when the installer creates infrastructure for vSphere clusters.
func ValidateForProvisioning(ic *types.InstallConfig, certificates ...x509.Certificate) error {
	allErrs := field.ErrorList{}
	if ic.Platform.VSphere == nil {
		return errors.New(field.Required(field.NewPath("platform", "vsphere"), "vSphere validation requires a vSphere platform configuration").Error())
	}

	allErrs = append(allErrs, validation.ValidateForProvisioning(ic.Platform.VSphere, field.NewPath("platform").Child("vsphere"))...)
	allErrs = append(allErrs, folderExists(ic, field.NewPath("platform").Child("vsphere").Child("folder"))...)

	return allErrs.ToAggregate()
}

// folderExists returns an error if a folder is specified in the vSphere platform but a folder with that name is not found in the datacenter.
func folderExists(ic *types.InstallConfig, fldPath *field.Path) field.ErrorList {
	allErrs := field.ErrorList{}
	cfg := ic.VSphere

	// If no folder is specified, skip this check as the folder will be created.
	if cfg.Folder == "" {
		return allErrs
	}

	vim25Client, _, err := vspheretypes.CreateVSphereClients(context.TODO(), cfg.VCenter, cfg.Username, cfg.Password)
	if err != nil {
		err = errors.Wrap(err, "unable to connect to vCenter API")
		return append(allErrs, field.InternalError(fldPath, err))
	}

	finder := find.NewFinder(vim25Client)

	ctx, cancel := context.WithTimeout(context.TODO(), 60*time.Second)
	defer cancel()

	if _, err = finder.Folder(ctx, cfg.Folder); err != nil {
		return append(allErrs, field.Invalid(fldPath, cfg.Folder, err.Error()))
	}
	return nil
}

func validatevCenterVersion(vim25 *vim25.Client, fldPath *field.Path) field.ErrorList {
	allErrs := field.ErrorList{}
	constraints, err := version.NewConstraint("<= 6.7")

	if err != nil {
		return append(allErrs, field.InternalError(fldPath, err))
	}

	vCenterVersion, err := version.NewVersion(vim25.ServiceContent.About.Version)
	if err != nil {
		return append(allErrs, field.InternalError(fldPath, err))
	}
	build, err := strconv.Atoi(vim25.ServiceContent.About.Build)
	if err != nil {
		return append(allErrs, field.InternalError(fldPath, err))
	}

	if constraints.Check(vCenterVersion) {
		logrus.Warnf("The out-of-tree storage driver requires virtual hardware 15 and vSphere 6.7 U3. Current vCenter version: %s, build: %s",
			vim25.ServiceContent.About.Version, vim25.ServiceContent.About.Build)

	} else {
		// This is the vCenter 6.7 U3 build number
		// Anything less than this version is unsupported with the
		// out-of-tree CSI.
		// https://kb.vmware.com/s/article/2143838
		// https://vsphere-csi-driver.sigs.k8s.io/supported_features_matrix.html
		if build < esxi67U3BuildNumber {
			logrus.Warnf("The out-of-tree storage driver requires virtual hardware 15 and vSphere 6.7 U3. Current vCenter version: %s, build: %s",
				vim25.ServiceContent.About.Version, vim25.ServiceContent.About.Build)
		}
	}
	return nil
}

func validateESXiVersion(vim25 *vim25.Client, ic *types.InstallConfig, fldPath *field.Path) field.ErrorList {

	allErrs := field.ErrorList{}
	cfg := ic.VSphere
	clusterPath := fmt.Sprintf("/%s/host/%s", cfg.Datacenter, cfg.Cluster)
	finder := find.NewFinder(vim25)

	ctx, cancel := context.WithTimeout(context.TODO(), 60*time.Second)
	defer cancel()

	clusters, err := finder.ClusterComputeResourceList(ctx, clusterPath)

	if err != nil {
		err = errors.Wrapf(err, "unable to find cluster on path: %s", clusterPath)
		return append(allErrs, field.InternalError(fldPath, err))
	}

	v67, err := version.NewVersion("6.7")
	if err != nil {
		return append(allErrs, field.InternalError(fldPath, err))
	}

	hosts, err := clusters[0].Hosts(context.TODO())
	if err != nil {
		err = errors.Wrapf(err, "unable to find hosts from cluster on path: %s", clusterPath)
		return append(allErrs, field.InternalError(fldPath, err))
	}

	for _, h := range hosts {
		var mh mo.HostSystem
		h.Properties(context.TODO(), h.Reference(), []string{"config.product"}, &mh)

		esxiHostVersion, err := version.NewVersion(mh.Config.Product.Version)
		if err != nil {
			return append(allErrs, field.InternalError(fldPath, err))
		}

		if esxiHostVersion.Equal(v67) {
			build, err := strconv.Atoi(mh.Config.Product.Build)
			if err != nil {
				return append(allErrs, field.InternalError(fldPath, err))
			}
			// This is the ESXi 6.7 U3 build number
			// Anything less than this version is unsupported with the
			// out-of-tree CSI.
			// https://kb.vmware.com/s/article/2143838
			// https://vsphere-csi-driver.sigs.k8s.io/supported_features_matrix.html
			if build < esxi67U3BuildNumber {
				logrus.Warnf("The out-of-tree storage driver requires virtual hardware 15 and vSphere 6.7 U3. The ESXi host: %s is version: %s and build: %s",
					h.Name(), mh.Config.Product.Version, mh.Config.Product.Build)
			}
		} else if esxiHostVersion.LessThan(v67) { // If ESXi host is <6.7 print warning
			logrus.Warnf("The out-of-tree storage driver requires virtual hardware 15 and vSphere 6.7 U3. The ESXi host: %s is version: %s and build: %s",
				h.Name(), mh.Config.Product.Version, mh.Config.Product.Build)
		}
	}
	return nil
}
