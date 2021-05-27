package quota

import (
	"context"
	"fmt"
	"strings"

	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"

	"github.com/openshift/installer/pkg/asset"
	"github.com/openshift/installer/pkg/asset/installconfig"
	configgcp "github.com/openshift/installer/pkg/asset/installconfig/gcp"
	openstackvalidation "github.com/openshift/installer/pkg/asset/installconfig/openstack/validation"
	"github.com/openshift/installer/pkg/asset/machines"
	"github.com/openshift/installer/pkg/asset/quota/aws"
	"github.com/openshift/installer/pkg/asset/quota/gcp"
	"github.com/openshift/installer/pkg/asset/quota/openstack"
	"github.com/openshift/installer/pkg/diagnostics"
	"github.com/openshift/installer/pkg/quota"
	quotaaws "github.com/openshift/installer/pkg/quota/aws"
	quotagcp "github.com/openshift/installer/pkg/quota/gcp"
	typesaws "github.com/openshift/installer/pkg/types/aws"
	"github.com/openshift/installer/pkg/types/azure"
	"github.com/openshift/installer/pkg/types/baremetal"
	typesgcp "github.com/openshift/installer/pkg/types/gcp"
	typesibmcloud "github.com/openshift/installer/pkg/types/ibmcloud"
	"github.com/openshift/installer/pkg/types/kubevirt"
	"github.com/openshift/installer/pkg/types/libvirt"
	"github.com/openshift/installer/pkg/types/none"
	typesopenstack "github.com/openshift/installer/pkg/types/openstack"
	"github.com/openshift/installer/pkg/types/ovirt"
	"github.com/openshift/installer/pkg/types/vsphere"
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

	// TODO: IBM[#87]: Add quota checks
	if ic.Config.Platform.Name() == typesibmcloud.Name {
		return nil
	}

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
	case typesaws.Name:
		if !quotaaws.SupportedRegions.Has(ic.AWS.Region) {
			logrus.Debugf("%s does not support API for checking quotas, therefore skipping.", ic.AWS.Region)
			return nil
		}
		services := []string{"ec2", "vpc"}
		session, err := ic.AWS.Session(context.TODO())
		if err != nil {
			return errors.Wrap(err, "failed to load AWS session")
		}
		q, err := quotaaws.Load(context.TODO(), session, ic.AWS.Region, services...)
		if quotaaws.IsUnauthorized(err) {
			logrus.Warnf("Missing permissions to fetch Quotas and therefore will skip checking them: %v, make sure you have `servicequotas:ListAWSDefaultServiceQuotas` permission available to the user.", err)
			return nil
		}
		if err != nil {
			return errors.Wrapf(err, "failed to load Quota for services: %s", strings.Join(services, ", "))
		}
		instanceTypes, err := aws.InstanceTypes(context.TODO(), session, ic.AWS.Region)
		if quotaaws.IsUnauthorized(err) {
			logrus.Warnf("Missing permissions to fetch instance types and therefore will skip checking Quotas: %v, make sure you have `ec2:DescribeInstanceTypes` permission available to the user.", err)
			return nil
		}
		if err != nil {
			return errors.Wrapf(err, "failed to load instance types for %s", ic.AWS.Region)
		}
		reports, err := quota.Check(q, aws.Constraints(ic.Config, masters, workers, instanceTypes))
		if err != nil {
			return summarizeFailingReport(reports)
		}
		summarizeReport(reports)
	case typesgcp.Name:
		services := []string{"compute.googleapis.com", "iam.googleapis.com"}
		q, err := quotagcp.Load(context.TODO(), ic.Config.Platform.GCP.ProjectID, services...)
		if quotagcp.IsUnauthorized(err) {
			logrus.Warnf("Missing permissions to fetch Quotas and therefore will skip checking them: %v, make sure you have `roles/servicemanagement.quotaViewer` assigned to the user.", err)
			return nil
		}
		if err != nil {
			return errors.Wrapf(err, "failed to load Quota for services: %s", strings.Join(services, ", "))
		}
		session, err := configgcp.GetSession(context.TODO())
		if err != nil {
			return errors.Wrap(err, "failed to load GCP session")
		}
		client, err := gcp.NewClient(context.TODO(), session, ic.Config.Platform.GCP.ProjectID)
		if err != nil {
			return errors.Wrap(err, "failed to create client for quota constraints")
		}
		reports, err := quota.Check(q, gcp.Constraints(client, ic.Config, masters, workers))
		if err != nil {
			return summarizeFailingReport(reports)
		}
		summarizeReport(reports)
	case typesibmcloud.Name:
		// TODO: IBM[#87]: Add quota checks
		return nil
	case typesopenstack.Name:
		ci, err := openstackvalidation.GetCloudInfo(ic.Config)
		if err != nil {
			return errors.Wrap(err, "failed to get cloud info")
		}
		if ci == nil {
			logrus.Warnf("Empty OpenStack cloud info and therefore will skip checking quota validation.")
			return nil
		}
		reports, err := quota.Check(ci.Quotas, openstack.Constraints(ci, masters, workers, ic.Config.NetworkType))
		if err != nil {
			return summarizeFailingReport(reports)
		}
		summarizeReport(reports)
	case azure.Name, baremetal.Name, libvirt.Name, none.Name, ovirt.Name, vsphere.Name, kubevirt.Name:
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
func summarizeFailingReport(reports []quota.ConstraintReport) error {
	var notavailable []string
	var unknown []string
	var regionMessage string
	for _, report := range reports {
		switch report.Result {
		case quota.NotAvailable:
			if report.For.Region != "" {
				regionMessage = " in " + report.For.Region
			} else {
				regionMessage = ""
			}
			notavailable = append(notavailable, fmt.Sprintf("%s is not available%s because %s", report.For.Name, regionMessage, report.Message))
		case quota.Unknown:
			unknown = append(unknown, report.For.Name)
		default:
			continue
		}
	}

	if len(notavailable) == 0 && len(unknown) > 0 {
		// all quotas are missing information so warn and skip
		logrus.Warnf("Failed to find information on quotas %s", strings.Join(unknown, ", "))
		return nil
	}

	msg := strings.Join(notavailable, ", ")
	if len(unknown) > 0 {
		msg = fmt.Sprintf("%s, and could not find information on %s", msg, strings.Join(unknown, ", "))
	}
	return &diagnostics.Err{Reason: "MissingQuota", Message: msg}
}

// summarizeReport summarizes a report when there are availble.
func summarizeReport(reports []quota.ConstraintReport) {
	var low []string
	var regionMessage string
	for _, report := range reports {
		switch report.Result {
		case quota.AvailableButLow:
			if report.For.Region != "" {
				regionMessage = " (" + report.For.Region + ")"
			} else {
				regionMessage = ""
			}
			low = append(low, fmt.Sprintf("%s%s", report.For.Name, regionMessage))
		default:
			continue
		}
	}
	if len(low) > 0 {
		logrus.Warnf("Following quotas %s are available but will be completely used pretty soon.", strings.Join(low, ", "))
	}
}
