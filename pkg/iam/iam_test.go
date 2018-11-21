package iam

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_getPermissionsFromPolicyDoc(t *testing.T) {
	type args struct {
		userPolicyDocs *[]string
	}
	tests := []struct {
		name string
		args args
		want []string
	}{
		// Test 1 is testing this encoded policy document:
		// {"Version":"2012-10-17","Statement":[{"Effect":"Allow","Action":"*","Resource":"*"}]}
		{
			name: "Allow all actions",
			args: args{userPolicyDocs: &[]string{"%7B%0A%20%20%22Version%22%3A%20%222012-10-17%22%2C%0A%20%20%22Statement%22%3A%20%5B%0A%20%20%20%20%7B%0A%20%20%20%20%20%20%22Effect%22%3A%20%22Allow%22%2C%0A%20%20%20%20%20%20%22Action%22%3A%20%22%2A%22%2C%0A%20%20%20%20%20%20%22Resource%22%3A%20%22%2A%22%0A%20%20%20%20%7D%0A%20%20%5D%0A%7D"}},
			want: []string{"*"},
		},

		// Test 2 is testing this encoded policy document:
		// {"Version":"2012-10-17","Statement":[{"Effect":"Allow","Action":"s3:*","Resource":"*"}]}
		{
			name: "Allow all S3 actions",
			args: args{userPolicyDocs: &[]string{"%7B%0A%20%20%22Version%22%3A%20%222012-10-17%22%2C%0A%20%20%22Statement%22%3A%20%5B%0A%20%20%20%20%7B%0A%20%20%20%20%20%20%22Effect%22%3A%20%22Allow%22%2C%0A%20%20%20%20%20%20%22Action%22%3A%20%22s3%3A%2A%22%2C%0A%20%20%20%20%20%20%22Resource%22%3A%20%22%2A%22%0A%20%20%20%20%7D%0A%20%20%5D%0A%7D"}},
			want: []string{"s3:*"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := getPermissionsFromPolicyDoc(tt.args.userPolicyDocs)
			assert.Equal(t, tt.want, got)
		})
	}
}

func Test_checkOffPermissionsByPrefix(t *testing.T) {
	type args struct {
		checklist map[string]bool
		prefix    string
	}
	tests := []struct {
		name string
		args args
		want map[string]bool
	}{
		{
			name: "Check off all actions",
			args: args{
				checklist: map[string]bool{
					"s3:DeleteBucket":       false,
					"s3:CreateBucket":       false,
					"s3:CreateObject":       false,
					"iam:GetRole":           false,
					"route53:GetHostedZone": false,
				},
				prefix: "*",
			},
			want: map[string]bool{
				"s3:DeleteBucket":       true,
				"s3:CreateBucket":       true,
				"s3:CreateObject":       true,
				"iam:GetRole":           true,
				"route53:GetHostedZone": true,
			},
		},
		{
			name: "Check off all actions related to S3",
			args: args{
				checklist: map[string]bool{
					"s3:DeleteBucket":       false,
					"s3:CreateBucket":       false,
					"s3:CreateObject":       false,
					"iam:GetRole":           false,
					"route53:GetHostedZone": false,
				},
				prefix: "s3:*",
			},
			want: map[string]bool{
				"s3:DeleteBucket":       true,
				"s3:CreateBucket":       true,
				"s3:CreateObject":       true,
				"iam:GetRole":           false,
				"route53:GetHostedZone": false,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := checkOffPermissionsByPrefix(tt.args.checklist, tt.args.prefix)
			assert.Equal(t, tt.want, got)
		})
	}
}
