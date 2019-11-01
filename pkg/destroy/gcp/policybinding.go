package gcp

import (
	"fmt"
	"strings"

	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"

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

func clearIAMPolicyBindings(policy *resourcemanager.Policy, clusterID string, logger logrus.FieldLogger) bool {
	removedBindings := false
	for _, binding := range policy.Bindings {
		members := []string{}
		for _, member := range binding.Members {
			if strings.HasPrefix(member, fmt.Sprintf("serviceAccount:%s", clusterID)) {
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
	if !clearIAMPolicyBindings(policy, o.ClusterID, o.Logger) {
		deletedPolicy := o.setPendingItems("iampolicy", []string{})
		if len(deletedPolicy) > 0 {
			o.Logger.Infof("Deleted IAM project role bindings")
		}
		return nil
	}
	o.setPendingItems("iampolicy", []string{"policy"})
	err = o.setProjectIAMPolicy(policy)
	return aggregateError([]error{err}, 1)
}
