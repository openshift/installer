package credentialsrequest

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParseCredentialRequestBytes(t *testing.T) {
	cases := []struct {
		name    string
		data    string
		wantNil bool
		wantErr string
	}{
		{
			name: "valid credentials request",
			data: `apiVersion: cloudcredential.openshift.io/v1
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
		{
			name: "errors on credentials request without provider spec",
			data: `apiVersion: cloudcredential.openshift.io/v1
kind: CredentialsRequest
metadata:
  name: cloud-credentials
  namespace: openshift-ingress-operator
spec:
  secretRef:
    name: cloud-credentials
    namespace: openshift-ingress-operator
`,
			wantErr: "has no provider spec",
		},
		{
			name:    "errors on invalid yaml",
			data:    "not valid yaml: [",
			wantErr: "failed to parse",
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			req, err := parseCredentialRequestBytes([]byte(tc.data), "test.yaml")
			if tc.wantErr != "" {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tc.wantErr)
				return
			}
			assert.NoError(t, err)
			if tc.wantNil {
				assert.Nil(t, req)
			} else {
				assert.NotNil(t, req)
			}
		})
	}
}

func TestParseCredentialRequestBytesMultiple(t *testing.T) {
	files := map[string]string{
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
	}

	var requests []CredentialRequest
	for name, content := range files {
		req, err := parseCredentialRequestBytes([]byte(content), name)
		if !assert.NoError(t, err, "file %s", name) {
			return
		}
		if req != nil {
			requests = append(requests, *req)
		}
	}
	assert.Len(t, requests, 2)
}

func TestParseCredentialRequestBytesFields(t *testing.T) {
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

	r, err := parseCredentialRequestBytes([]byte(content), "ingress.yaml")
	if !assert.NoError(t, err) {
		return
	}
	if !assert.NotNil(t, r) {
		return
	}

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

func TestParseCredentialRequestBytesLoadRoundTrip(t *testing.T) {
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
  providerSpec:
    apiVersion: cloudcredential.openshift.io/v1
    kind: AWSProviderSpec
    statementEntries:
      - effect: Allow
        action:
          - elasticloadbalancing:DescribeLoadBalancers
        resource: "*"
`
	tmpDir := t.TempDir()
	filename := "99_credentials-request_openshift-ingress-operator-cloud-credentials.yaml"
	if err := os.WriteFile(filepath.Join(tmpDir, filename), []byte(content), 0o600); err != nil {
		t.Fatalf("failed to write test file: %v", err)
	}

	data, err := os.ReadFile(filepath.Join(tmpDir, filename))
	if err != nil {
		t.Fatalf("failed to read test file: %v", err)
	}

	req, err := parseCredentialRequestBytes(data, filename)
	assert.NoError(t, err)
	assert.NotNil(t, req)
	assert.IsType(t, &AWSProviderSpec{}, req.ProviderSpec)
}
