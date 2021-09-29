package ibmcloud

import (
	"net/http"
	"strings"

	"github.com/IBM/platform-services-go-sdk/iampolicymanagementv1"
	"github.com/pkg/errors"
)

const iamAuthorizationTypeName = "iam authorization"

// listIAMAuthorizations lists IAM authorizations
func (o *ClusterUninstaller) listIAMAuthorizations() (cloudResources, error) {
	o.Logger.Debugf("Listing IAM authorizations")
	ctx, cancel := o.contextWithTimeout()
	defer cancel()

	options := o.iamPolicyManagementSvc.NewListPoliciesOptions(o.AccountID)
	options.SetType("authorization")
	resources, _, err := o.iamPolicyManagementSvc.ListPoliciesWithContext(ctx, options)
	if err != nil {
		return nil, errors.Wrapf(err, "Failed to list IAM authorizations")
	}

	result := []cloudResource{}
	for _, policy := range resources.Policies {
		if o.policyMatches(policy) {
			result = append(result, cloudResource{
				key:      *policy.ID,
				name:     *policy.ID,
				status:   *policy.State,
				typeName: iamAuthorizationTypeName,
				id:       *policy.ID,
			})
		}
	}

	return cloudResources{}.insert(result...), nil
}

// policyMatches returns true if the IAM Policy matches the one set up to allow
// the VPC service to read from the COS bucket containing the uploaded RHCOS image.
func (o *ClusterUninstaller) policyMatches(policy iampolicymanagementv1.Policy) bool {
	// Ideally we would match using the description field of the Policy type. However,
	// this is not currently supported in the IBM Terraform Provider. An issue has
	// been opened for this: https://github.com/IBM-Cloud/terraform-provider-ibm/issues/2894
	// Once implemented, this code can be updaed to check for a well-known
	// description that includes the cluster's InfraID or something equivalent.
	onlyOneSubject := len(policy.Subjects) == 1
	subjectMatches := len(policy.Subjects[0].Attributes) == 3 &&
		*policy.Subjects[0].Attributes[0].Name == "accountId" &&
		*policy.Subjects[0].Attributes[0].Value == o.AccountID &&
		*policy.Subjects[0].Attributes[1].Name == "serviceName" &&
		*policy.Subjects[0].Attributes[1].Value == "is" &&
		*policy.Subjects[0].Attributes[2].Name == "resourceType" &&
		*policy.Subjects[0].Attributes[2].Value == "image"

	onlyOneResource := len(policy.Resources) == 1
	cosInstanceID, err := o.COSInstanceID()
	if err != nil {
		o.Logger.Warn("Unable to determine IAM policy match. Failed to obtain COS instance ID. ", err)
		return false
	}
	resourceMatches := len(policy.Resources[0].Attributes) == 3 &&
		*policy.Resources[0].Attributes[0].Name == "accountId" &&
		*policy.Resources[0].Attributes[0].Value == o.AccountID &&
		*policy.Resources[0].Attributes[1].Name == "serviceName" &&
		*policy.Resources[0].Attributes[1].Value == "cloud-object-storage" &&
		*policy.Resources[0].Attributes[2].Name == "serviceInstance" &&
		//strings.Contains(*policy.Resources[0].Attributes[2].Value, cosInstanceID) &&
		strings.Contains(cosInstanceID, *policy.Resources[0].Attributes[2].Value) &&
		*policy.Resources[0].Attributes[2].Operator == "stringEquals"

	return onlyOneSubject && onlyOneResource && subjectMatches && resourceMatches
}

func (o *ClusterUninstaller) deleteIAMAuthorization(item cloudResource) error {
	o.Logger.Debugf("Deleting IAM authorization %q", item.name)
	ctx, cancel := o.contextWithTimeout()
	defer cancel()

	options := o.iamPolicyManagementSvc.NewDeletePolicyOptions(item.id)
	details, err := o.iamPolicyManagementSvc.DeletePolicyWithContext(ctx, options)

	if err != nil && details != nil && details.StatusCode == http.StatusNotFound {
		// The resource is gone
		o.deletePendingItems(item.typeName, []cloudResource{item})
		o.Logger.Infof("Deleted IAM authorization %q", item.name)
		return nil
	}

	if err != nil && details != nil && details.StatusCode != http.StatusNotFound {
		return errors.Wrapf(err, "Failed to delete IAM authorization %s", item.name)
	}

	return nil
}

// destroyIAMAuthorizations removes all IAM authorization pertaining to the cluster ID.
func (o *ClusterUninstaller) destroyIAMAuthorizations() error {
	found, err := o.listIAMAuthorizations()
	if err != nil {
		return err
	}

	items := o.insertPendingItems(iamAuthorizationTypeName, found.list())

	for _, item := range items {
		if _, ok := found[item.key]; !ok {
			// This item has finished deletion.
			o.deletePendingItems(item.typeName, []cloudResource{item})
			o.Logger.Infof("Deleted IAM authorization %q", item.name)
			continue
		}
		err = o.deleteIAMAuthorization(item)
		if err != nil {
			o.errorTracker.suppressWarning(item.key, err, o.Logger)
		}
	}

	if items = o.getPendingItems(iamAuthorizationTypeName); len(items) > 0 {
		return errors.Errorf("%d items pending", len(items))
	}
	return nil
}
