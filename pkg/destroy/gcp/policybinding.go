package gcp

import (
	"strings"

	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"k8s.io/apimachinery/pkg/util/sets"

	resourcemanager "google.golang.org/api/cloudresourcemanager/v1"
)

func (o *ClusterUninstaller) getProjectIAMPolicy() (*resourcemanager.Policy, error) {
	o.Logger.Debug("Fetching project IAM policy")
	ctx, cancel := o.contextWithTimeout()
	defer cancel()
	req := &resourcemanager.GetIamPolicyRequest{}
	policy, err := o.rmSvc.Projects.GetIamPolicy(o.ProjectID, req).Context(ctx).Do()
	if err != nil {
		return nil, errors.Wrapf(err, "failed to fetch project IAM policy")
	}
	return policy, nil
}

func (o *ClusterUninstaller) setProjectIAMPolicy(policy *resourcemanager.Policy) error {
	o.Logger.Debug("Setting project IAM policy")
	ctx, cancel := o.contextWithTimeout()
	defer cancel()
	req := &resourcemanager.SetIamPolicyRequest{Policy: policy}
	_, err := o.rmSvc.Projects.SetIamPolicy(o.ProjectID, req).Context(ctx).Do()
	if err != nil {
		return errors.Wrapf(err, "failed to set project IAM policy")
	}
	return nil
}

func clearIAMPolicyBindings(policy *resourcemanager.Policy, emails sets.String, logger logrus.FieldLogger) bool {
	removedBindings := false
	for _, binding := range policy.Bindings {
		members := []string{}
		for _, member := range binding.Members {
			email := strings.TrimPrefix(strings.TrimPrefix(member, "deleted:"), "serviceAccount:")
			if emails.Has(email) {
				logger.Debugf("IAM: removing %s from role %s", member, binding.Role)
				removedBindings = true
				continue
			}
			members = append(members, member)
		}
		binding.Members = members
	}
	return removedBindings
}

// destroyIAMPolicyBindings removes any role bindings from the project policy to
// service accounts that start with the cluster's infra ID.
func (o *ClusterUninstaller) destroyIAMPolicyBindings() error {
	policy, err := o.getProjectIAMPolicy()
	if err != nil {
		return err
	}

	sas := o.getPendingItems("serviceaccount_binding")
	emails := sets.NewString()
	for _, item := range sas {
		emails.Insert(item.url)
	}

	if !clearIAMPolicyBindings(policy, emails, o.Logger) {
		pendingPolicy := o.getPendingItems("iampolicy")
		if len(pendingPolicy) > 0 {
			o.Logger.Infof("Deleted IAM project role bindings")
			o.deletePendingItems("iampolicy", pendingPolicy)
		}
		return nil
	}
	o.insertPendingItems("iampolicy", []cloudResource{{key: "policy", name: "policy", typeName: "iampolicy"}})
	err = o.setProjectIAMPolicy(policy)
	return aggregateError([]error{err}, 1)
}
