package gcp

import (
	"context"
	"strings"

	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	resourcemanager "google.golang.org/api/cloudresourcemanager/v1"
	"k8s.io/apimachinery/pkg/util/sets"
)

func (o *ClusterUninstaller) getProjectIAMPolicy(ctx context.Context) (*resourcemanager.Policy, error) {
	o.Logger.Debug("Fetching project IAM policy")
	ctx, cancel := context.WithTimeout(ctx, defaultTimeout)
	defer cancel()
	req := &resourcemanager.GetIamPolicyRequest{}
	policy, err := o.rmSvc.Projects.GetIamPolicy(o.ProjectID, req).Context(ctx).Do()
	if err != nil {
		return nil, errors.Wrapf(err, "failed to fetch project IAM policy")
	}
	return policy, nil
}

func (o *ClusterUninstaller) setProjectIAMPolicy(ctx context.Context, policy *resourcemanager.Policy) error {
	o.Logger.Debug("Setting project IAM policy")
	ctx, cancel := context.WithTimeout(ctx, defaultTimeout)
	defer cancel()
	req := &resourcemanager.SetIamPolicyRequest{Policy: policy}
	_, err := o.rmSvc.Projects.SetIamPolicy(o.ProjectID, req).Context(ctx).Do()
	if err != nil {
		return errors.Wrapf(err, "failed to set project IAM policy")
	}
	return nil
}

func (o *ClusterUninstaller) clearIAMPolicyBindings(policy *resourcemanager.Policy, emails sets.String, logger logrus.FieldLogger) bool {
	removedBindings := false
	for _, binding := range policy.Bindings {
		members := []string{}
		for _, member := range binding.Members {
			email := policyMemberToEmail(member)
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

// policyMemberToEmail takes member of IAM policy binding and converts it to service account email.
// https://cloud.google.com/iam/docs/reference/rest/v1/Policy#Binding
// see members[]
func policyMemberToEmail(member string) string {
	email := strings.TrimPrefix(strings.TrimPrefix(member, "deleted:"), "serviceAccount:")
	if idx := strings.Index(email, "?uid"); idx != -1 {
		email = email[:idx]
	}
	return email
}
