package aws

import (
	"context"
	"testing"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/arn"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/service/iam"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"

	"github.com/openshift/installer/pkg/destroy/aws/mock"
)

type (
	editARNFuncs           []func(arn *arn.ARN)
	attachedRolePoliciesCB func(*iam.ListAttachedRolePoliciesOutput, bool) bool
	listAccessKeysCB       func(*iam.ListAccessKeysOutput, bool) bool
	listInstanceProfilesCB func(*iam.ListInstanceProfilesForRoleOutput, bool) bool
	listRolePoliciesCB     func(*iam.ListRolePoliciesOutput, bool) bool
	listUserPoliciesCB     func(*iam.ListUserPoliciesOutput, bool) bool
)

var (
	notFoundError = awserr.New(
		iam.ErrCodeNoSuchEntityException,
		"some aws iam error: not found",
		nil,
	)
	someIAMError = awserr.New(
		iam.ErrCodeDeleteConflictException,
		"some aws iam error: conflict",
		nil,
	)
)

func TestDeleteInstanceProfileByName(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	iamClient := mock.NewMockIAMAPI(mockCtrl)

	cases := []struct {
		name     string
		role     string
		errorMsg string
	}{
		{
			name:     "Delete Instance Profile succeeds",
			role:     "role-name",
			errorMsg: "",
		},
		{
			name:     "Delete Instance Profile succeeds when not found",
			role:     "role-not-found",
			errorMsg: "",
		},
		{
			name:     "Delete Instance Profile fails",
			role:     "role-fails",
			errorMsg: "some aws iam error",
		},
	}

	iamClient.
		EXPECT().
		DeleteInstanceProfile(gomock.Any(), gomock.Eq("role-name")).
		Return(nil).
		AnyTimes()
	iamClient.
		EXPECT().
		DeleteInstanceProfile(gomock.Any(), gomock.Eq("role-fails")).
		Return(someIAMError).
		AnyTimes()
	iamClient.
		EXPECT().
		DeleteInstanceProfile(gomock.Any(), gomock.Any()).
		Return(notFoundError).
		AnyTimes()

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			err := deleteIAMInstanceProfileByName(context.TODO(), iamClient, &tc.role, nullLogger)
			if tc.errorMsg != "" {
				assert.Regexp(t, tc.errorMsg, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestDeleteInstanceProfile(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	iamClient := mock.NewMockIAMAPI(mockCtrl)

	validProfileARN := func(arn *arn.ARN) {
		arn.Resource = "instance-profile/profile-name"
	}
	notFoundProfileARN := func(arn *arn.ARN) {
		arn.Resource = "instance-profile/not-found"
	}
	failGetProfileARN := func(arn *arn.ARN) {
		arn.Resource = "instance-profile/fail-get"
	}
	failDeleteProfileARN := func(arn *arn.ARN) {
		arn.Resource = "instance-profile/fail-delete"
	}
	failRemoveRoleARN := func(arn *arn.ARN) {
		arn.Resource = "instance-profile/role-fail"
	}
	notFoundRoleARN := func(arn *arn.ARN) {
		arn.Resource = "instance-profile/not-found-role"
	}
	invalidTypeARN := func(arn *arn.ARN) {
		arn.Resource = "random-type/random-name"
	}

	cases := []struct {
		name string
		// arn      arn.ARN
		editFuncs editARNFuncs
		errorMsg  string
	}{
		{
			name:      "Delete Instance Profile with no roles succeeds",
			editFuncs: editARNFuncs{validProfileARN},
			errorMsg:  "",
		},
		{
			name:      "Delete Instance Profile succeeds when profile not found",
			editFuncs: editARNFuncs{notFoundProfileARN},
			errorMsg:  "",
		},
		{
			name:      "Delete Instance Profile wrong resourceType fails",
			editFuncs: editARNFuncs{invalidTypeARN},
			errorMsg:  "ARN passed to deleteIAMInstanceProfile: ",
		},
		{
			name:      "Delete Instance Profile invalid resource fails",
			editFuncs: editARNFuncs{},
			errorMsg:  "does not contain the expected slash",
		},
		{
			name:      "Delete Instance Profile fails to get profile",
			editFuncs: editARNFuncs{failGetProfileARN},
			errorMsg:  "some aws iam error",
		},
		{
			name:      "Delete Instance Profile remove role fails",
			editFuncs: editARNFuncs{failRemoveRoleARN},
			errorMsg:  "dissociating .*: some aws iam error",
		},
		{
			name:      "Delete Instance Profile remove role not found fails",
			editFuncs: editARNFuncs{notFoundRoleARN},
			errorMsg:  "some aws iam error",
		},
		{
			name:      "Delete Instance Profile fails to delete",
			editFuncs: editARNFuncs{failDeleteProfileARN},
			errorMsg:  "some aws iam error",
		},
	}

	iamClient.
		EXPECT().
		GetInstanceProfile(gomock.Any(), gomock.Eq("profile-name")).
		Return(
			&iam.GetInstanceProfileOutput{
				InstanceProfile: &iam.InstanceProfile{
					InstanceProfileName: aws.String("profile-name"),
				},
			}, nil).
		AnyTimes()
	iamClient.
		EXPECT().
		GetInstanceProfile(gomock.Any(), gomock.Eq("fail-delete")).
		Return(
			&iam.GetInstanceProfileOutput{
				InstanceProfile: &iam.InstanceProfile{
					InstanceProfileName: aws.String("fail-delete"),
				},
			}, nil).
		AnyTimes()
	iamClient.
		EXPECT().
		GetInstanceProfile(gomock.Any(), gomock.Eq("role-fail")).
		Return(
			&iam.GetInstanceProfileOutput{
				InstanceProfile: &iam.InstanceProfile{
					InstanceProfileName: aws.String("profile-name"),
					Roles: []*iam.Role{
						{RoleName: aws.String("role-name")},
						{RoleName: aws.String("role-fail")},
					},
				},
			}, nil).
		AnyTimes()
	iamClient.
		EXPECT().
		GetInstanceProfile(gomock.Any(), gomock.Eq("not-found-role")).
		Return(
			&iam.GetInstanceProfileOutput{
				InstanceProfile: &iam.InstanceProfile{
					InstanceProfileName: aws.String("profile-name"),
					Roles: []*iam.Role{
						{RoleName: aws.String("role-name")},
						{RoleName: aws.String("role-not-found")},
						{RoleName: aws.String("role-fail")},
					},
				},
			}, nil).
		AnyTimes()
	iamClient.
		EXPECT().
		GetInstanceProfile(gomock.Any(), gomock.Eq("not-found")).
		Return(nil, notFoundError).
		AnyTimes()
	iamClient.
		EXPECT().
		GetInstanceProfile(gomock.Any(), gomock.Eq("fail-get")).
		Return(nil, someIAMError).
		AnyTimes()
	iamClient.
		EXPECT().
		GetInstanceProfile(gomock.Any(), gomock.Eq("random-name")).
		Times(0) // Should not be called with wrong resource type

	iamClient.
		EXPECT().
		RemoveRoleFromInstanceProfile(gomock.Any(), gomock.Any(), gomock.Eq("role-name")).
		Return(nil).
		AnyTimes()
	iamClient.
		EXPECT().
		RemoveRoleFromInstanceProfile(gomock.Any(), gomock.Any(), gomock.Eq("role-not-found")).
		Return(notFoundError).
		AnyTimes()
	iamClient.
		EXPECT().
		RemoveRoleFromInstanceProfile(gomock.Any(), gomock.Any(), gomock.Eq("role-fail")).
		Return(someIAMError).
		AnyTimes()

	iamClient.
		EXPECT().
		DeleteInstanceProfile(gomock.Any(), gomock.Eq("profile-name")).
		Return(nil).
		AnyTimes()
	iamClient.
		EXPECT().
		DeleteInstanceProfile(gomock.Any(), gomock.Eq("not-found")).
		Return(notFoundError).
		AnyTimes()
	iamClient.
		EXPECT().
		DeleteInstanceProfile(gomock.Any(), gomock.Eq("fail-delete")).
		Return(someIAMError).
		AnyTimes()
	iamClient.
		EXPECT().
		DeleteInstanceProfile(gomock.Any(), gomock.Any()).
		Times(0) // Should not be called when anything else fails

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			arn := arn.ARN{}
			for _, edit := range tc.editFuncs {
				edit(&arn)
			}
			err := deleteIAMInstanceProfile(context.TODO(), iamClient, arn, nullLogger)
			if tc.errorMsg != "" {
				assert.Regexp(t, tc.errorMsg, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestDeleteRolePolicies(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	iamClient := mock.NewMockIAMAPI(mockCtrl)

	cases := []struct {
		name     string
		role     string
		errorMsg string
	}{
		{
			name:     "Delete Role Policies succeeds",
			role:     "role-name",
			errorMsg: "",
		},
		{
			name:     "Delete Role Policies list fails",
			role:     "list-policies-fails",
			errorMsg: "listing IAM role policies",
		},
		{
			name:     "Delete Role Policies delete fails",
			role:     "delete-policies-fails",
			errorMsg: "deleting IAM role policy .*",
		},
		{
			name:     "Delete Role Policies delete fails when not found",
			role:     "delete-policies-fails-not-found",
			errorMsg: "deleting IAM role policy .*",
		},
	}

	iamClient.
		EXPECT().
		ListRolePoliciesPages(gomock.Any(), gomock.Eq("role-name"), gomock.Any()).
		DoAndReturn(
			func(_ context.Context, _ string, fn listRolePoliciesCB) error {
				results := &iam.ListRolePoliciesOutput{
					PolicyNames: []*string{aws.String("policy-name"), aws.String("policy-name")},
				}
				// Make sure to run the function passed in as argument
				fn(results, true)
				return nil
			}).
		AnyTimes()
	iamClient.
		EXPECT().
		ListRolePoliciesPages(gomock.Any(), gomock.Eq("delete-policies-fails"), gomock.Any()).
		DoAndReturn(
			func(_ context.Context, _ string, fn listRolePoliciesCB) error {
				results := &iam.ListRolePoliciesOutput{
					PolicyNames: []*string{aws.String("policy-name"), aws.String("policy-not-found"), aws.String("policy-fail"), aws.String("policy-fail")},
				}
				// Make sure to run the function passed in as argument
				fn(results, true)
				return nil
			}).
		AnyTimes()
	iamClient.
		EXPECT().
		ListRolePoliciesPages(gomock.Any(), gomock.Eq("delete-policies-fails-not-found"), gomock.Any()).
		DoAndReturn(
			func(_ context.Context, _ string, fn listRolePoliciesCB) error {
				results := &iam.ListRolePoliciesOutput{
					PolicyNames: []*string{aws.String("policy-name"), aws.String("policy-not-found"), aws.String("policy-name"), aws.String("policy-not-found")},
				}
				// Make sure to run the function passed in as argument
				fn(results, true)
				return nil
			}).
		AnyTimes()
	iamClient.
		EXPECT().
		ListRolePoliciesPages(gomock.Any(), gomock.Eq("list-policies-fails"), gomock.Any()).
		Return(someIAMError).AnyTimes()

	iamClient.
		EXPECT().
		DeleteRolePolicy(gomock.Any(), gomock.Any(), gomock.Eq("policy-name")).
		Return(nil).
		AnyTimes()
	iamClient.
		EXPECT().
		DeleteRolePolicy(gomock.Any(), gomock.Any(), gomock.Eq("policy-fail")).
		Return(someIAMError).
		AnyTimes()
	iamClient.
		EXPECT().
		DeleteRolePolicy(gomock.Any(), gomock.Any(), gomock.Eq("policy-not-found")).
		Return(notFoundError).
		AnyTimes()

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			err := deleteRolePolicies(context.TODO(), iamClient, tc.role, nullLogger)
			if tc.errorMsg != "" {
				assert.Regexp(t, tc.errorMsg, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestDeleteAttachedRolePolicies(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	iamClient := mock.NewMockIAMAPI(mockCtrl)

	cases := []struct {
		name     string
		role     string
		errorMsg string
	}{
		{
			name:     "Delete Attached Role Policies succeeds",
			role:     "role-name",
			errorMsg: "",
		},
		{
			name:     "Delete Attached Role Policies list fails",
			role:     "list-attached-fails",
			errorMsg: "listing attached IAM role policies",
		},
		{
			name:     "Delete Attached Role Policies detach fails",
			role:     "delete-attached-fails",
			errorMsg: "detaching IAM role policy .*",
		},
		{
			name:     "Delete Attached Role Policies delete fails when not found",
			role:     "delete-attached-fails-not-found",
			errorMsg: "detaching IAM role policy .*",
		},
	}

	iamClient.
		EXPECT().
		ListAttachedRolePoliciesPages(gomock.Any(), gomock.Eq("role-name"), gomock.Any()).
		DoAndReturn(
			func(_ context.Context, _ string, fn attachedRolePoliciesCB) error {
				results := &iam.ListAttachedRolePoliciesOutput{
					AttachedPolicies: []*iam.AttachedPolicy{
						{PolicyName: aws.String("policy-name"), PolicyArn: aws.String("policy-arn")},
						{PolicyName: aws.String("policy-name"), PolicyArn: aws.String("policy-arn")},
					},
				}
				fn(results, true)
				return nil
			}).
		AnyTimes()
	iamClient.
		EXPECT().
		ListAttachedRolePoliciesPages(gomock.Any(), gomock.Eq("delete-attached-fails"), gomock.Any()).
		DoAndReturn(
			func(_ context.Context, _ string, fn attachedRolePoliciesCB) error {
				results := &iam.ListAttachedRolePoliciesOutput{
					AttachedPolicies: []*iam.AttachedPolicy{
						{PolicyName: aws.String("policy-name"), PolicyArn: aws.String("policy-arn")},
						{PolicyName: aws.String("policy-fail"), PolicyArn: aws.String("policy-fail-arn")},
						{PolicyName: aws.String("policy-name"), PolicyArn: aws.String("policy-arn")},
					},
				}
				fn(results, true)
				return nil
			}).
		AnyTimes()
	iamClient.
		EXPECT().
		ListAttachedRolePoliciesPages(gomock.Any(), gomock.Eq("delete-attached-fails-not-found"), gomock.Any()).
		DoAndReturn(
			func(_ context.Context, _ string, fn attachedRolePoliciesCB) error {
				results := &iam.ListAttachedRolePoliciesOutput{
					AttachedPolicies: []*iam.AttachedPolicy{
						{PolicyName: aws.String("policy-name"), PolicyArn: aws.String("policy-arn")},
						{PolicyName: aws.String("policy-fail"), PolicyArn: aws.String("policy-not-found-arn")},
						{PolicyName: aws.String("policy-name"), PolicyArn: aws.String("policy-arn")},
					},
				}
				fn(results, true)
				return nil
			}).
		AnyTimes()
	iamClient.
		EXPECT().
		ListAttachedRolePoliciesPages(gomock.Any(), gomock.Eq("list-attached-fails"), gomock.Any()).
		Return(someIAMError).
		AnyTimes()

	iamClient.
		EXPECT().
		DetachRolePolicy(gomock.Any(), gomock.Any(), gomock.Eq("policy-arn")).
		Return(nil).
		AnyTimes()
	iamClient.
		EXPECT().
		DetachRolePolicy(gomock.Any(), gomock.Any(), gomock.Eq("policy-fail-arn")).
		Return(someIAMError).
		AnyTimes()
	iamClient.
		EXPECT().
		DetachRolePolicy(gomock.Any(), gomock.Any(), gomock.Eq("policy-not-found-arn")).
		Return(notFoundError).
		AnyTimes()

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			err := deleteAttachedRolePolicies(context.TODO(), iamClient, tc.role, nullLogger)
			if tc.errorMsg != "" {
				assert.Regexp(t, tc.errorMsg, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestDeleteInstanceProfilesForRole(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	iamClient := mock.NewMockIAMAPI(mockCtrl)

	cases := []struct {
		name     string
		role     string
		errorMsg string
	}{
		{
			name:     "Delete Instance Profile succeeds",
			role:     "role-name",
			errorMsg: "",
		},
		{
			name:     "Delete Instance Profile list fails",
			role:     "list-profiles-fails",
			errorMsg: "listing IAM instance profiles: ",
		},
		{
			name:     "Delete Instance Profile parse fails",
			role:     "parse-profiles-fails",
			errorMsg: "parse ARN for IAM instance profile",
		},
		{
			name:     "Delete Instance profile delete fails",
			role:     "delete-profiles-fails",
			errorMsg: "deleting .*",
		},
		{
			name:     "Delete Instance Profile delete fails when not found",
			role:     "delete-profiles-fails-not-found",
			errorMsg: "listing IAM instance profiles: .*",
		},
	}

	iamClient.
		EXPECT().
		ListInstanceProfilesForRolePages(gomock.Any(), gomock.Eq("role-name"), gomock.Any()).
		DoAndReturn(
			func(_ context.Context, _ string, fn listInstanceProfilesCB) error {
				results := &iam.ListInstanceProfilesForRoleOutput{
					InstanceProfiles: []*iam.InstanceProfile{
						{Arn: aws.String("arn:aws:iam:::instance-profile/profile-name")},
						{Arn: aws.String("arn:aws:iam:::instance-profile/profile-name")},
					},
				}
				fn(results, true)
				return nil
			}).
		AnyTimes()
	iamClient.
		EXPECT().
		ListInstanceProfilesForRolePages(gomock.Any(), gomock.Eq("list-profiles-fails"), gomock.Any()).
		Return(someIAMError).
		AnyTimes()
	iamClient.
		EXPECT().
		ListInstanceProfilesForRolePages(gomock.Any(), gomock.Eq("parse-profiles-fails"), gomock.Any()).
		DoAndReturn(
			func(_ context.Context, _ string, fn listInstanceProfilesCB) error {
				results := &iam.ListInstanceProfilesForRoleOutput{
					InstanceProfiles: []*iam.InstanceProfile{
						{Arn: aws.String("arn:aws:iam:profile-name")},
						{Arn: aws.String("profile-name")},
					},
				}
				fn(results, true)
				return nil
			}).
		AnyTimes()
	iamClient.
		EXPECT().
		ListInstanceProfilesForRolePages(gomock.Any(), gomock.Eq("delete-profiles-fails"), gomock.Any()).
		DoAndReturn(
			func(_ context.Context, _ string, fn listInstanceProfilesCB) error {
				results := &iam.ListInstanceProfilesForRoleOutput{
					InstanceProfiles: []*iam.InstanceProfile{
						{Arn: aws.String("arn:aws:iam:::instance-profile/profile-name")},
						{Arn: aws.String("arn:aws:iam:::instance-profile/profile-name-fails")},
						{Arn: aws.String("arn:aws:iam:::instance-profile/profile-name-not-found")},
					},
				}
				fn(results, true)
				return nil
			}).
		AnyTimes()
	iamClient.
		EXPECT().
		ListInstanceProfilesForRolePages(gomock.Any(), gomock.Eq("delete-profiles-fails-not-found"), gomock.Any()).
		Return(notFoundError).
		AnyTimes()

	// deleteIAMInstanceProfile is already being independently test, so here
	// we just care about two possible outcomes: it either succeeds or fails
	iamClient.
		EXPECT().
		GetInstanceProfile(gomock.Any(), gomock.Eq("profile-name")).
		DoAndReturn(
			func(_ context.Context, name string) (*iam.GetInstanceProfileOutput, error) {
				profile := &iam.GetInstanceProfileOutput{
					InstanceProfile: &iam.InstanceProfile{
						InstanceProfileName: aws.String(name),
					},
				}
				return profile, nil
			},
		).
		AnyTimes()
	iamClient.
		EXPECT().
		GetInstanceProfile(gomock.Any(), gomock.Any()).Return(nil, someIAMError).
		AnyTimes()
	iamClient.
		EXPECT().
		DeleteInstanceProfile(gomock.Any(), gomock.Any()).Return(nil).
		AnyTimes()

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			err := deleteInstanceProfilesForRole(context.TODO(), iamClient, tc.role, nullLogger)
			if tc.errorMsg != "" {
				assert.Regexp(t, tc.errorMsg, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestDeleteIAMRole(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	iamClient := mock.NewMockIAMAPI(mockCtrl)

	validRoleARN := func(arn *arn.ARN) {
		arn.Resource = "role/role-name"
	}
	failDeletePoliciesARN := func(arn *arn.ARN) {
		arn.Resource = "role/role-delete-policies-fails"
	}
	failDeleteAttachedARN := func(arn *arn.ARN) {
		arn.Resource = "role/role-delete-attached-fails"
	}
	failDeleteInstanceARN := func(arn *arn.ARN) {
		arn.Resource = "role/role-delete-instance-fails"
	}
	failDeleteRoleARN := func(arn *arn.ARN) {
		arn.Resource = "role/role-delete-fails"
	}
	notFoundRoleARN := func(arn *arn.ARN) {
		arn.Resource = "role/role-not-found-fails"
	}
	invalidTypeARN := func(arn *arn.ARN) {
		arn.Resource = "random-type/random-name"
	}

	cases := []struct {
		name      string
		editFuncs editARNFuncs
		errorMsg  string
	}{
		{
			name:      "Delete Role wrong resourceType fails",
			editFuncs: editARNFuncs{invalidTypeARN},
			errorMsg:  "ARN passed to deleteIAMRole: ",
		},
		{
			name:      "Delete Role invalid resource fails",
			editFuncs: editARNFuncs{},
			errorMsg:  "does not contain the expected slash",
		},
		{
			name:      "Delete Role delete role policies fails",
			editFuncs: editARNFuncs{failDeletePoliciesARN},
			errorMsg:  "listing IAM role policies: .*",
		},
		{
			name:      "Delete Role delete attached policies fails",
			editFuncs: editARNFuncs{failDeleteAttachedARN},
			errorMsg:  "some aws iam error",
		},
		{
			name:      "Delete Role delete instance profiles fails",
			editFuncs: editARNFuncs{failDeleteInstanceARN},
			errorMsg:  "listing IAM instance profiles: .*",
		},
		{
			name:      "Delete Role suceeds",
			editFuncs: editARNFuncs{validRoleARN},
			errorMsg:  "",
		},
		{
			name:      "Delete Role fails",
			editFuncs: editARNFuncs{failDeleteRoleARN},
			errorMsg:  "some aws iam error",
		},
		{
			name:      "Delete Role fails when role not found",
			editFuncs: editARNFuncs{notFoundRoleARN},
			errorMsg:  "some aws iam error",
		},
	}

	// deleteRolePolicies is already being independently tested, so here we
	// just care about two possible outcomes: it either succeeds or fails
	iamClient.
		EXPECT().
		ListRolePoliciesPages(gomock.Any(), gomock.Eq("role-delete-policies-fails"), gomock.Any()).
		Return(someIAMError).
		AnyTimes()
	iamClient.
		EXPECT().
		ListRolePoliciesPages(gomock.Any(), gomock.Eq("random-name"), gomock.Any()).
		Times(0) // Should not be called on wrong resource type
	iamClient.
		EXPECT().
		ListRolePoliciesPages(gomock.Any(), gomock.Any(), gomock.Any()).
		Return(nil).
		AnyTimes()
	iamClient.
		EXPECT().
		DeleteRolePolicy(gomock.Any(), gomock.Eq("random-name"), gomock.Any()).
		Times(0) // Should not be called on wrong resource type

	// deleteAttachedRolePolicies is already being independently test, so here
	// we just care about two possible outcomes: it either succeeds or fails
	iamClient.
		EXPECT().
		ListAttachedRolePoliciesPages(gomock.Any(), gomock.Eq("role-delete-attached-fails"), gomock.Any()).
		Return(someIAMError).
		AnyTimes()
	iamClient.
		EXPECT().
		ListAttachedRolePoliciesPages(gomock.Any(), gomock.Eq("random-name"), gomock.Any()).
		Times(0) // Should not be called on wrong resource type
	iamClient.
		EXPECT().
		ListAttachedRolePoliciesPages(gomock.Any(), gomock.Any(), gomock.Any()).
		Return(nil).
		AnyTimes()
	iamClient.
		EXPECT().
		DetachRolePolicy(gomock.Any(), gomock.Eq("random-name"), gomock.Any()).
		Times(0) // Should not be called on wrong resource type

	// deleteInstanceProfilesForRole is already being independently tested, so
	// here we just care about two possible outcomes: it either succeeds or
	// fails
	iamClient.
		EXPECT().
		ListInstanceProfilesForRolePages(gomock.Any(), gomock.Eq("role-delete-instance-fails"), gomock.Any()).
		Return(someIAMError).
		AnyTimes()
	iamClient.
		EXPECT().
		ListInstanceProfilesForRolePages(gomock.Any(), gomock.Eq("random-name"), gomock.Any()).
		Times(0) // Should not be called on wrong resource type
	iamClient.
		EXPECT().
		ListInstanceProfilesForRolePages(gomock.Any(), gomock.Any(), gomock.Any()).
		Return(nil).
		AnyTimes()

	iamClient.
		EXPECT().
		DeleteRole(gomock.Any(), gomock.Eq("role-name")).
		Return(nil).
		AnyTimes()
	iamClient.
		EXPECT().
		DeleteRole(gomock.Any(), gomock.Eq("role-delete-fails")).
		Return(someIAMError).
		AnyTimes()
	iamClient.
		EXPECT().
		DeleteRole(gomock.Any(), gomock.Eq("role-not-found-fails")).
		Return(notFoundError).
		AnyTimes()
	iamClient.
		EXPECT().
		DeleteRole(gomock.Any(), gomock.Eq("random-name")).
		Times(0) // Should not be called on wrong resource type

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			arn := arn.ARN{}
			for _, edit := range tc.editFuncs {
				edit(&arn)
			}
			err := deleteIAMRole(context.TODO(), iamClient, arn, nullLogger)
			if tc.errorMsg != "" {
				assert.Regexp(t, tc.errorMsg, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestDeleteUserPolicies(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	iamClient := mock.NewMockIAMAPI(mockCtrl)

	cases := []struct {
		name     string
		user     string
		errorMsg string
	}{
		{
			name:     "Delete User Policy succeeds",
			user:     "user-name",
			errorMsg: "",
		},
		{
			name:     "Delete User Policy list fails",
			user:     "user-list-fails",
			errorMsg: "listing IAM user policies",
		},
		{
			name:     "Delete User Policy list fails not found",
			user:     "user-list-fails-not-found",
			errorMsg: "listing IAM user policies",
		},
		{
			name:     "Delete User Policy delete fails",
			user:     "user-delete-fails",
			errorMsg: "deleting IAM user policy .*",
		},
		{
			name:     "Delete User Policy delete fails not found",
			user:     "user-delete-fails-not-found",
			errorMsg: "deleting IAM user policy .*",
		},
	}

	iamClient.
		EXPECT().
		ListUserPoliciesPages(gomock.Any(), gomock.Eq("user-name"), gomock.Any()).
		DoAndReturn(
			func(_ context.Context, _ string, fn listUserPoliciesCB) error {
				results := &iam.ListUserPoliciesOutput{
					PolicyNames: []*string{aws.String("policy-name"), aws.String("policy-name")},
				}
				fn(results, true)
				return nil
			}).
		AnyTimes()
	iamClient.
		EXPECT().
		ListUserPoliciesPages(gomock.Any(), gomock.Eq("user-delete-fails"), gomock.Any()).
		DoAndReturn(
			func(_ context.Context, _ string, fn listUserPoliciesCB) error {
				results := &iam.ListUserPoliciesOutput{
					PolicyNames: []*string{aws.String("policy-name"), aws.String("policy-fails")},
				}
				fn(results, true)
				return nil
			}).
		AnyTimes()
	iamClient.
		EXPECT().
		ListUserPoliciesPages(gomock.Any(), gomock.Eq("user-delete-fails-not-found"), gomock.Any()).
		DoAndReturn(
			func(_ context.Context, _ string, fn listUserPoliciesCB) error {
				results := &iam.ListUserPoliciesOutput{
					PolicyNames: []*string{aws.String("policy-name"), aws.String("policy-not-found")},
				}
				fn(results, true)
				return nil
			}).
		AnyTimes()
	iamClient.
		EXPECT().
		ListUserPoliciesPages(gomock.Any(), gomock.Eq("user-list-fails"), gomock.Any()).
		Return(someIAMError).
		AnyTimes()
	iamClient.
		EXPECT().
		ListUserPoliciesPages(gomock.Any(), gomock.Eq("user-list-fails-not-found"), gomock.Any()).
		Return(notFoundError).
		AnyTimes()

	iamClient.
		EXPECT().
		DeleteUserPolicy(gomock.Any(), gomock.Any(), gomock.Eq("policy-name")).
		Return(nil).
		AnyTimes()
	iamClient.
		EXPECT().
		DeleteUserPolicy(gomock.Any(), gomock.Any(), gomock.Eq("polify-fails")).
		Return(someIAMError).
		AnyTimes()
	iamClient.
		EXPECT().
		DeleteUserPolicy(gomock.Any(), gomock.Any(), gomock.Any()).
		Return(notFoundError).
		AnyTimes()

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			err := deleteUserPolicies(context.TODO(), iamClient, tc.user, nullLogger)
			if tc.errorMsg != "" {
				assert.Regexp(t, tc.errorMsg, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestDeleteAccessKeys(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	iamClient := mock.NewMockIAMAPI(mockCtrl)

	cases := []struct {
		name     string
		user     string
		errorMsg string
	}{
		{
			name:     "Delete Access Keys succeeds",
			user:     "user-name",
			errorMsg: "",
		},
		{
			name:     "Delete Access Keys list fails",
			user:     "user-list-fails",
			errorMsg: "listing IAM access keys: .*",
		},
		{
			name:     "Delete Access Keys list fails not found",
			user:     "user-list-fails-not-found",
			errorMsg: "listing IAM access keys: .*",
		},
		{
			name:     "Delete Access Keys delete fails",
			user:     "user-delete-fails",
			errorMsg: "deleting IAM access key .*",
		},
		{
			name:     "Delete Access Keys delete fails not found",
			user:     "user-delete-fails-not-found",
			errorMsg: "deleting IAM access key .*",
		},
	}

	iamClient.
		EXPECT().
		ListAccessKeysPages(gomock.Any(), gomock.Eq("user-name"), gomock.Any()).
		DoAndReturn(
			func(_ context.Context, _ string, fn listAccessKeysCB) error {
				results := &iam.ListAccessKeysOutput{
					AccessKeyMetadata: []*iam.AccessKeyMetadata{
						{AccessKeyId: aws.String("key-id")},
						{AccessKeyId: aws.String("key-id")},
					},
				}
				fn(results, true)
				return nil
			}).
		AnyTimes()
	iamClient.
		EXPECT().
		ListAccessKeysPages(gomock.Any(), gomock.Eq("user-delete-fails"), gomock.Any()).
		DoAndReturn(
			func(_ context.Context, _ string, fn listAccessKeysCB) error {
				results := &iam.ListAccessKeysOutput{
					AccessKeyMetadata: []*iam.AccessKeyMetadata{
						{AccessKeyId: aws.String("key-id")},
						{AccessKeyId: aws.String("key-delete-fails")},
					},
				}
				fn(results, true)
				return nil
			}).
		AnyTimes()
	iamClient.
		EXPECT().
		ListAccessKeysPages(gomock.Any(), gomock.Eq("user-delete-fails-not-found"), gomock.Any()).
		DoAndReturn(
			func(_ context.Context, _ string, fn listAccessKeysCB) error {
				results := &iam.ListAccessKeysOutput{
					AccessKeyMetadata: []*iam.AccessKeyMetadata{
						{AccessKeyId: aws.String("key-id")},
						{AccessKeyId: aws.String("key-not-found")},
					},
				}
				fn(results, true)
				return nil
			}).
		AnyTimes()
	iamClient.
		EXPECT().
		ListAccessKeysPages(gomock.Any(), gomock.Eq("user-list-fails"), gomock.Any()).
		Return(someIAMError).
		AnyTimes()
	iamClient.
		EXPECT().
		ListAccessKeysPages(gomock.Any(), gomock.Eq("user-list-fails-not-found"), gomock.Any()).
		Return(notFoundError).
		AnyTimes()

	iamClient.
		EXPECT().
		DeleteAccessKey(gomock.Any(), gomock.Any(), gomock.Eq("key-id")).
		Return(nil).
		AnyTimes()
	iamClient.
		EXPECT().
		DeleteAccessKey(gomock.Any(), gomock.Any(), gomock.Eq("key-delete-fails")).
		Return(someIAMError).
		AnyTimes()
	iamClient.
		EXPECT().
		DeleteAccessKey(gomock.Any(), gomock.Any(), gomock.Any()).
		Return(notFoundError).
		AnyTimes()

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			err := deleteAccessKeys(context.TODO(), iamClient, tc.user, nullLogger)
			if tc.errorMsg != "" {
				assert.Regexp(t, tc.errorMsg, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestDeleteIAMUser(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	iamClient := mock.NewMockIAMAPI(mockCtrl)

	cases := []struct {
		name     string
		user     string
		errorMsg string
	}{
		{
			name:     "Delete User succeeds",
			user:     "user-name",
			errorMsg: "",
		},
		{
			name:     "Delete User delete policies fails",
			user:     "user-policies-fails",
			errorMsg: "listing IAM user policies: .*",
		},
		{
			name:     "Delete User delete access keys fails",
			user:     "user-access-fails",
			errorMsg: "listing IAM access keys: .*",
		},
		{
			name:     "Delete User delete fails",
			user:     "user-delete-fails",
			errorMsg: "some aws iam error",
		},
		{
			name:     "Delete User delete fails not found",
			user:     "user-delete-fails-not-found",
			errorMsg: "some aws iam error",
		},
	}

	// deleteUserPolicies is already being independently tested, so here
	// we just care about two possible outcomes: it either succeeds or fails
	iamClient.
		EXPECT().
		ListUserPoliciesPages(gomock.Any(), gomock.Eq("user-policies-fails"), gomock.Any()).
		Return(someIAMError).
		AnyTimes()
	iamClient.
		EXPECT().
		ListUserPoliciesPages(gomock.Any(), gomock.Any(), gomock.Any()).
		Return(nil).
		AnyTimes()
	// deleteAccessKeys is already being independently tested, so here we just
	// care about two possible outcomes: it either succeeds or fails
	iamClient.
		EXPECT().
		ListAccessKeysPages(gomock.Any(), gomock.Eq("user-access-fails"), gomock.Any()).
		Return(someIAMError).
		AnyTimes()
	iamClient.
		EXPECT().
		ListAccessKeysPages(gomock.Any(), gomock.Any(), gomock.Any()).
		Return(nil).
		AnyTimes()

	iamClient.
		EXPECT().
		DeleteUser(gomock.Any(), gomock.Eq("user-name")).
		Return(nil).
		AnyTimes()
	iamClient.
		EXPECT().
		DeleteUser(gomock.Any(), gomock.Eq("user-delete-fails")).
		Return(someIAMError).
		AnyTimes()
	iamClient.
		EXPECT().
		DeleteUser(gomock.Any(), gomock.Any()).
		Return(notFoundError).
		AnyTimes()

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			err := deleteIAMUser(context.TODO(), iamClient, tc.user, nullLogger)
			if tc.errorMsg != "" {
				assert.Regexp(t, tc.errorMsg, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestDeleteIAM(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	iamClient := mock.NewMockIAMAPI(mockCtrl)

	invalidTypeARN := func(arn *arn.ARN) {
		arn.Resource = "random-type/random-name"
	}

	cases := []struct {
		name      string
		editFuncs editARNFuncs
		errorMsg  string
	}{
		{
			name:      "Delete IAM invalid resource type fails",
			editFuncs: editARNFuncs{invalidTypeARN},
			errorMsg:  "unrecognized EC2 resource type random-type",
		},
		{
			name:      "Delete IAM parse resource fails",
			editFuncs: editARNFuncs{},
			errorMsg:  "resource .* does not contain the expected slash",
		},
	}

	// All the other functions were independently tested, so we just check for
	// cases when deleteIAM might fail by itself

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			arn := arn.ARN{}
			for _, edit := range tc.editFuncs {
				edit(&arn)
			}
			err := deleteIAMWithClient(context.TODO(), iamClient, arn, nullLogger)
			if tc.errorMsg != "" {
				assert.Regexp(t, tc.errorMsg, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
