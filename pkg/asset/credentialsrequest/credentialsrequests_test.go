package credentialsrequest

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParseCredentialsRequests(t *testing.T) {
	cases := []struct {
		name      string
		files     map[string]string
		wantCount int
		wantErr   string
	}{
		{
			name: "single valid credentials request",
			files: map[string]string{
				"openshift-ingress-operator.yaml": `apiVersion: cloudcredential.openshift.io/v1
kind: CredentialsRequest
metadata:
  name: cloud-credentials
  namespace: openshift-ingress-operator
spec:
  secretRef:
    name: cloud-credentials
    namespace: openshift-ingress-operator
  serviceAccountNames:
    - ingress-operator
  providerSpec:
    apiVersion: cloudcredential.openshift.io/v1
    kind: AWSProviderSpec
    statementEntries:
      - effect: Allow
        action:
          - elasticloadbalancing:DescribeLoadBalancers
          - route53:ListHostedZones
          - route53:ChangeResourceRecordSets
          - tag:GetResources
        resource: "*"
`,
			},
			wantCount: 1,
		},
		{
			name: "multiple credentials requests",
			files: map[string]string{
				"ingress.yaml": `apiVersion: cloudcredential.openshift.io/v1
kind: CredentialsRequest
metadata:
  name: cloud-credentials
  namespace: openshift-ingress-operator
spec:
  secretRef:
    name: cloud-credentials
    namespace: openshift-ingress-operator
  serviceAccountNames:
    - ingress-operator
  providerSpec:
    apiVersion: cloudcredential.openshift.io/v1
    kind: AWSProviderSpec
    statementEntries:
      - effect: Allow
        action:
          - elasticloadbalancing:DescribeLoadBalancers
        resource: "*"
`,
				"machine-api.yaml": `apiVersion: cloudcredential.openshift.io/v1
kind: CredentialsRequest
metadata:
  name: aws-cloud-credentials
  namespace: openshift-machine-api
spec:
  secretRef:
    name: aws-cloud-credentials
    namespace: openshift-machine-api
  serviceAccountNames:
    - machine-api-controllers
  providerSpec:
    apiVersion: cloudcredential.openshift.io/v1
    kind: AWSProviderSpec
    statementEntries:
      - effect: Allow
        action:
          - ec2:CreateTags
          - ec2:DescribeInstances
          - ec2:RunInstances
          - ec2:TerminateInstances
        resource: "*"
`,
			},
			wantCount: 2,
		},
		{
			name: "skips non-yaml files",
			files: map[string]string{
				"README.md": "not a yaml file",
				"ingress.yaml": `apiVersion: cloudcredential.openshift.io/v1
kind: CredentialsRequest
metadata:
  name: cloud-credentials
  namespace: openshift-ingress-operator
spec:
  secretRef:
    name: cloud-credentials
    namespace: openshift-ingress-operator
  serviceAccountNames:
    - ingress-operator
  providerSpec:
    apiVersion: cloudcredential.openshift.io/v1
    kind: AWSProviderSpec
    statementEntries:
      - effect: Allow
        action:
          - elasticloadbalancing:DescribeLoadBalancers
        resource: "*"
`,
			},
			wantCount: 1,
		},
		{
			name: "skips credentials request without provider spec",
			files: map[string]string{
				"no-provider.yaml": `apiVersion: cloudcredential.openshift.io/v1
kind: CredentialsRequest
metadata:
  name: cloud-credentials
  namespace: openshift-ingress-operator
spec:
  secretRef:
    name: cloud-credentials
    namespace: openshift-ingress-operator
`,
				"valid.yaml": `apiVersion: cloudcredential.openshift.io/v1
kind: CredentialsRequest
metadata:
  name: cloud-credentials
  namespace: openshift-machine-api
spec:
  secretRef:
    name: aws-cloud-credentials
    namespace: openshift-machine-api
  serviceAccountNames:
    - machine-api-controllers
  providerSpec:
    apiVersion: cloudcredential.openshift.io/v1
    kind: AWSProviderSpec
    statementEntries:
      - effect: Allow
        action:
          - ec2:DescribeInstances
        resource: "*"
`,
			},
			wantCount: 1,
		},
		{
			name:    "empty directory returns error",
			files:   map[string]string{},
			wantErr: "no credentials requests found",
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			tmpDir := t.TempDir()
			for name, content := range tc.files {
				if err := os.WriteFile(filepath.Join(tmpDir, name), []byte(content), 0o600); err != nil {
					t.Fatalf("failed to write test file: %v", err)
				}
			}

			requests, err := parseCredentialsRequests("aws", tmpDir)
			if tc.wantErr != "" {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tc.wantErr)
				return
			}

			assert.NoError(t, err)
			assert.Len(t, requests, tc.wantCount)
		})
	}
}

func TestParseCredentialsRequestFields(t *testing.T) {
	tmpDir := t.TempDir()
	content := `apiVersion: cloudcredential.openshift.io/v1
kind: CredentialsRequest
metadata:
  name: cloud-credentials
  namespace: openshift-ingress-operator
spec:
  secretRef:
    name: cloud-credentials
    namespace: openshift-ingress-operator
  serviceAccountNames:
    - ingress-operator
    - ingress-operator-sa2
  providerSpec:
    apiVersion: cloudcredential.openshift.io/v1
    kind: AWSProviderSpec
    statementEntries:
      - effect: Allow
        action:
          - elasticloadbalancing:DescribeLoadBalancers
          - route53:ListHostedZones
        resource: "*"
      - effect: Allow
        action:
          - s3:GetObject
        resource: "arn:aws:s3:::my-bucket/*"
`
	if err := os.WriteFile(filepath.Join(tmpDir, "ingress.yaml"), []byte(content), 0o600); err != nil {
		t.Fatalf("failed to write test file: %v", err)
	}

	requests, err := parseCredentialsRequests("aws", tmpDir)
	if !assert.NoError(t, err) {
		return
	}
	if !assert.Len(t, requests, 1) {
		return
	}

	r := requests[0]
	assert.Equal(t, "cloud-credentials", r.Name)
	assert.Equal(t, "openshift-ingress-operator", r.Namespace)
	assert.Equal(t, "cloud-credentials", r.SecretRefName)
	assert.Equal(t, "openshift-ingress-operator", r.SecretRefNamespace)
	assert.Equal(t, []string{"ingress-operator", "ingress-operator-sa2"}, r.ServiceAccountNames)
	assert.IsType(t, &AWSProviderSpec{}, r.ProviderSpec)
	awsSpec := r.ProviderSpec.(*AWSProviderSpec)
	assert.Len(t, awsSpec.StatementEntries, 2)
	assert.Equal(t, "Allow", awsSpec.StatementEntries[0].Effect)
	assert.Equal(t, []string{"elasticloadbalancing:DescribeLoadBalancers", "route53:ListHostedZones"}, awsSpec.StatementEntries[0].Action)
	assert.Equal(t, "*", awsSpec.StatementEntries[0].Resource)
	assert.Equal(t, "arn:aws:s3:::my-bucket/*", awsSpec.StatementEntries[1].Resource)
}
