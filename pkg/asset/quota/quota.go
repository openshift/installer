package quota

import (
	"context"
	"fmt"
	"strings"

	"github.com/openshift/installer/pkg/asset"
	"github.com/openshift/installer/pkg/asset/installconfig"
	"github.com/openshift/installer/pkg/asset/machines"
	"github.com/openshift/installer/pkg/asset/quota/gcp"
	"github.com/openshift/installer/pkg/diagnostics"
	"github.com/openshift/installer/pkg/quota"
	quotagcp "github.com/openshift/installer/pkg/quota/gcp"
	"github.com/openshift/installer/pkg/types/aws"
	"github.com/openshift/installer/pkg/types/azure"
	"github.com/openshift/installer/pkg/types/baremetal"
	typesgcp "github.com/openshift/installer/pkg/types/gcp"
	"github.com/openshift/installer/pkg/types/libvirt"
	"github.com/openshift/installer/pkg/types/none"
	"github.com/openshift/installer/pkg/types/openstack"
	"github.com/openshift/installer/pkg/types/ovirt"
	"github.com/openshift/installer/pkg/types/vsphere"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

// PlatformQuotaCheck is an asset that validates the install-config platform for
// any resource requirements based on the quotas available.
type PlatformQuotaCheck struct {
}

var _ asset.Asset = (*PlatformQuotaCheck)(nil)

// Dependencies returns the dependencies for PlatformQuotaCheck
func (a *PlatformQuotaCheck) Dependencies() []asset.Asset {
	return []asset.Asset{
		&installconfig.InstallConfig{},
		&machines.Master{},
		&machines.Worker{},
	}
}

// Generate queries for input from the user.
func (a *PlatformQuotaCheck) Generate(dependencies asset.Parents) error {
	ic := &installconfig.InstallConfig{}
	mastersAsset := &machines.Master{}
	workersAsset := &machines.Worker{}
	dependencies.Get(ic, mastersAsset, workersAsset)

	masters, err := mastersAsset.Machines()
	if err != nil {
		return err
	}

	workers, err := workersAsset.MachineSets()
	if err != nil {
		return err
	}

	platform := ic.Config.Platform.Name()
	switch platform {
	case typesgcp.Name:
		services := []string{"compute.googleapis.com", "iam.googleapis.com"}
		q, err := quotagcp.Load(context.TODO(), ic.Config.Platform.GCP.ProjectID, services...)
		if quotagcp.IsUnauthorized(err) {
			logrus.Warnf("Missing permissions to fetch Quotas and therefore will skip checking them: %v", err)
			return nil
		}
		if err != nil {
			return errors.Wrapf(err, "failed to load Quota for services: %s", strings.Join(services, ", "))
		}
		reports, err := quota.Check(q, gcp.Constraints(ic.Config, masters, workers))
		if err != nil {
			return summarizeFailingReport(reports)
		}
		summarizeReport(reports)
	case aws.Name, azure.Name, baremetal.Name, libvirt.Name, none.Name, openstack.Name, ovirt.Name, vsphere.Name:
		// no special provisioning requirements to check
	default:
		err = fmt.Errorf("unknown platform type %q", platform)
	}
	return err
}

// Name returns the human-friendly name of the asset.
func (a *PlatformQuotaCheck) Name() string {
	return "Platform Quota Check"
}

// summarizeFailingReport summarizes a report when there are failing constraints.
func summarizeFailingReport(reports []quota.ConstraintReport) *diagnostics.Err {
	dErr := &diagnostics.Err{Reason: "MissingQuota"}
	var messages []string
	for _, report := range reports {
		switch report.Result {
		case quota.NotAvailable:
			messages = append(messages, fmt.Sprintf("%s is not available in %s because %s", report.For.Name, report.For.Region, report.Message))
		case quota.Unknown:
			messages = append(messages, fmt.Sprintf("could not find any quota information for %s", report.For.Name))
		default:
			continue
		}
	}
	dErr.Message = strings.Join(messages, ", ")
	return dErr
}

// summarizeReport summarizes a report when there are availble.
func summarizeReport(reports []quota.ConstraintReport) {
	var low []string
	for _, report := range reports {
		switch report.Result {
		case quota.AvailableButLow:
			low = append(low, fmt.Sprintf("%s (%s)", report.For.Name, report.For.Region))
		default:
			continue
		}
	}
	logrus.Warnf("Following quotas %s are available but will be completely used pretty soon.", strings.Join(low, ", "))
}
