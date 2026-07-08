package clusterapi

import (
	"context"
	"crypto"
	"crypto/rsa"
	"crypto/sha1" //nolint:gosec
	"crypto/tls"
	"crypto/x509"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"net/url"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/iam"
	iamtypes "github.com/aws/aws-sdk-go-v2/service/iam/types"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	s3types "github.com/aws/aws-sdk-go-v2/service/s3/types"
	jose "github.com/go-jose/go-jose/v4"
	"github.com/sirupsen/logrus"

	ccov1 "github.com/openshift/cloud-credential-operator/pkg/apis/cloudcredential/v1"
	"github.com/openshift/installer/pkg/asset/credentialsrequest"
	awsconfig "github.com/openshift/installer/pkg/asset/installconfig/aws"
	tlsconfig "github.com/openshift/installer/pkg/asset/tls"
	"github.com/openshift/installer/pkg/infrastructure/clusterapi"
	awstypes "github.com/openshift/installer/pkg/types/aws"
)

const (
	oidcClientID = "openshift"
	stsAudience  = "sts.amazonaws.com"
)

// provisionOIDCAndIAMRoles creates the OIDC infrastructure and IAM roles needed for
// STS-based credential provisioning. It creates the S3 bucket, OIDC provider, and IAM roles.
func provisionOIDCAndIAMRoles(ctx context.Context, in clusterapi.PreProvisionInput) error {
	infraID := in.InfraID
	p := in.InstallConfig.Config.Platform.AWS
	region := p.Region
	endpointOpts := awsconfig.EndpointOptions{
		Region:    region,
		Endpoints: p.ServiceEndpoints,
	}

	pubKeyData := in.BoundSASigningKey.PublicKeyData()
	if pubKeyData == nil {
		return fmt.Errorf("SA signing public key not available for STS provisioning")
	}

	saSigningPubKey, err := tlsconfig.PemToPublicKey(pubKeyData)
	if err != nil {
		return fmt.Errorf("failed to parse SA signing public key: %w", err)
	}

	partitionID, err := awsconfig.GetPartitionIDForRegion(ctx, region)
	if err != nil {
		return fmt.Errorf("failed to get partition for region %s: %w", region, err)
	}

	iamTags := buildIAMTags(infraID, p.UserTags)
	s3Tags := buildS3Tags(infraID, p.UserTags)

	iamClient, err := awsconfig.NewIAMClient(ctx, endpointOpts)
	if err != nil {
		return fmt.Errorf("failed to create IAM client: %w", err)
	}

	s3Client, err := awsconfig.NewS3Client(ctx, endpointOpts)
	if err != nil {
		return fmt.Errorf("failed to create S3 client: %w", err)
	}

	issuerURL, err := createOIDCBucket(ctx, s3Client, infraID, region, partitionID, saSigningPubKey, s3Tags)
	if err != nil {
		return fmt.Errorf("failed to create OIDC S3 bucket: %w", err)
	}
	logrus.Infof("Created OIDC S3 bucket: %s", issuerURL)

	oidcProviderARN, err := createOIDCProvider(ctx, iamClient, issuerURL, iamTags)
	if err != nil {
		return fmt.Errorf("failed to create IAM OIDC provider: %w", err)
	}
	logrus.Infof("Created IAM OIDC provider: %s", oidcProviderARN)

	if err := createSTSIAMRoles(ctx, iamClient, infraID, oidcProviderARN, issuerURL, in.CredentialsRequests.Requests, iamTags); err != nil {
		return fmt.Errorf("failed to create STS IAM roles: %w", err)
	}

	return nil
}

func buildIAMTags(infraID string, userTags map[string]string) []iamtypes.Tag {
	tags := []iamtypes.Tag{
		{
			Key:   aws.String(fmt.Sprintf("kubernetes.io/cluster/%s", infraID)),
			Value: aws.String("owned"),
		},
	}
	for k, v := range userTags {
		tags = append(tags, iamtypes.Tag{
			Key:   aws.String(k),
			Value: aws.String(v),
		})
	}
	return tags
}

func buildS3Tags(infraID string, userTags map[string]string) []s3types.Tag {
	tags := []s3types.Tag{
		{
			Key:   aws.String(fmt.Sprintf("kubernetes.io/cluster/%s", infraID)),
			Value: aws.String("owned"),
		},
	}
	for k, v := range userTags {
		tags = append(tags, s3types.Tag{
			Key:   aws.String(k),
			Value: aws.String(v),
		})
	}
	return tags
}

// createOIDCBucket creates the S3 bucket that serves the OIDC discovery
// document and JWKS for the cluster's service account token signing key.
func createOIDCBucket(ctx context.Context, client *s3.Client, infraID, region, partitionID string, pubKey *rsa.PublicKey, tags []s3types.Tag) (string, error) {
	bucketName := awstypes.OIDCBucketName(infraID)

	createInput := &s3.CreateBucketInput{
		Bucket: aws.String(bucketName),
	}
	if region != awstypes.UsEast1RegionID {
		createInput.CreateBucketConfiguration = &s3types.CreateBucketConfiguration{
			LocationConstraint: s3types.BucketLocationConstraint(region),
		}
	}

	if _, err := client.CreateBucket(ctx, createInput); err != nil {
		var boe *s3types.BucketAlreadyOwnedByYou
		if !errors.As(err, &boe) {
			return "", fmt.Errorf("failed to create bucket %s: %w", bucketName, err)
		}
		logrus.Debugf("Reusing existing OIDC bucket %s", bucketName)
	}

	if _, err := client.PutBucketTagging(ctx, &s3.PutBucketTaggingInput{
		Bucket:  aws.String(bucketName),
		Tagging: &s3types.Tagging{TagSet: tags},
	}); err != nil {
		return "", fmt.Errorf("failed to tag bucket: %w", err)
	}

	if _, err := client.PutPublicAccessBlock(ctx, &s3.PutPublicAccessBlockInput{
		Bucket: aws.String(bucketName),
		PublicAccessBlockConfiguration: &s3types.PublicAccessBlockConfiguration{
			BlockPublicAcls:       aws.Bool(true),
			BlockPublicPolicy:     aws.Bool(false),
			IgnorePublicAcls:      aws.Bool(true),
			RestrictPublicBuckets: aws.Bool(false),
		},
	}); err != nil {
		return "", fmt.Errorf("failed to set public access block: %w", err)
	}

	bucketARN := fmt.Sprintf("arn:%s:s3:::%s", partitionID, bucketName)
	policyDoc := fmt.Sprintf(`{
  "Version": "2012-10-17",
  "Statement": [
    {
      "Effect": "Allow",
      "Principal": "*",
      "Action": "s3:GetObject",
      "Resource": "%s/*"
    }
  ]
}`, bucketARN)
	if _, err := client.PutBucketPolicy(ctx, &s3.PutBucketPolicyInput{
		Bucket: aws.String(bucketName),
		Policy: aws.String(policyDoc),
	}); err != nil {
		return "", fmt.Errorf("failed to set bucket policy: %w", err)
	}

	issuerURL, err := awsconfig.OIDCIssuerURL(infraID, region)
	if err != nil {
		return "", err
	}

	jwksBytes, keyID, err := BuildJSONWebKeySet(pubKey)
	if err != nil {
		return "", fmt.Errorf("failed to build JWKS: %w", err)
	}

	discoveryDoc := buildDiscoveryDocument(issuerURL)

	if _, err := client.PutObject(ctx, &s3.PutObjectInput{
		Bucket:      aws.String(bucketName),
		Key:         aws.String(".well-known/openid-configuration"),
		Body:        strings.NewReader(discoveryDoc),
		ContentType: aws.String("application/json"),
	}); err != nil {
		return "", fmt.Errorf("failed to upload discovery document: %w", err)
	}

	if _, err := client.PutObject(ctx, &s3.PutObjectInput{
		Bucket:      aws.String(bucketName),
		Key:         aws.String("keys.json"),
		Body:        strings.NewReader(string(jwksBytes)),
		ContentType: aws.String("application/json"),
	}); err != nil {
		return "", fmt.Errorf("failed to upload JWKS: %w", err)
	}

	logrus.Debugf("Uploaded OIDC discovery document and JWKS (key-id=%s) to bucket %s", keyID, bucketName)
	return issuerURL, nil
}

func buildDiscoveryDocument(issuerURL string) string {
	return fmt.Sprintf(`{
  "issuer": "%s",
  "jwks_uri": "%s/keys.json",
  "response_types_supported": ["id_token"],
  "subject_types_supported": ["public"],
  "id_token_signing_alg_values_supported": ["RS256"],
  "claims_supported": ["aud", "exp", "sub", "iat", "iss"]
}`, issuerURL, issuerURL)
}

// BuildJSONWebKeySet builds a JWKS document for the given public key.
// The key ID is derived by hashing the PKIX DER encoding of the public key with SHA-256, matching the method used by kube-apiserver.
// Reference: https://github.com/kubernetes/kubernetes/blob/0f140bf1eeaf63c155f5eba1db8db9b5d52d5467/pkg/serviceaccount/jwt.go#L89-L111.
func BuildJSONWebKeySet(pubKey *rsa.PublicKey) ([]byte, string, error) {
	derBytes, err := x509.MarshalPKIXPublicKey(pubKey)
	if err != nil {
		return nil, "", fmt.Errorf("failed to serialize public key to DER format: %w", err)
	}

	hasher := crypto.SHA256.New()
	hasher.Write(derBytes) //nolint:errcheck
	keyID := base64.RawURLEncoding.EncodeToString(hasher.Sum(nil))

	jwk := jose.JSONWebKey{
		Key:       pubKey,
		KeyID:     keyID,
		Algorithm: string(jose.RS256),
		Use:       "sig",
	}

	jwks := jose.JSONWebKeySet{Keys: []jose.JSONWebKey{jwk}}
	data, err := json.Marshal(jwks)
	if err != nil {
		return nil, "", fmt.Errorf("failed to marshal JWKS: %w", err)
	}
	return data, keyID, nil
}

// createOIDCProvider creates an IAM OpenID Connect provider that trusts
// the given issuer URL for ServiceAccount token authentication.
func createOIDCProvider(ctx context.Context, client *iam.Client, issuerURL string, tags []iamtypes.Tag) (string, error) {
	thumbprint, err := getIssuerTLSThumbprint(ctx, issuerURL)
	if err != nil {
		return "", fmt.Errorf("failed to get TLS thumbprint for %s: %w", issuerURL, err)
	}

	result, err := client.CreateOpenIDConnectProvider(ctx, &iam.CreateOpenIDConnectProviderInput{
		Url:            aws.String(issuerURL),
		ClientIDList:   []string{oidcClientID, stsAudience},
		ThumbprintList: []string{thumbprint},
		Tags:           tags,
	})
	if err != nil {
		return "", fmt.Errorf("failed to create OIDC provider: %w", err)
	}

	return *result.OpenIDConnectProviderArn, nil
}

// getIssuerTLSThumbprint retrieves the SHA-1 fingerprint of the root CA
// certificate for the given issuer URL's TLS chain.
func getIssuerTLSThumbprint(ctx context.Context, issuerURL string) (string, error) {
	u, err := url.Parse(issuerURL)
	if err != nil {
		return "", fmt.Errorf("failed to parse issuer URL: %w", err)
	}

	host := u.Host
	if !strings.Contains(host, ":") {
		host += ":443"
	}

	dialer := &tls.Dialer{Config: &tls.Config{
		MinVersion: tls.VersionTLS12,
	}}
	conn, err := dialer.DialContext(ctx, "tcp", host)
	if err != nil {
		return "", fmt.Errorf("failed to connect to %s: %w", host, err)
	}
	defer conn.Close()

	tlsConn, ok := conn.(*tls.Conn)
	if !ok {
		return "", fmt.Errorf("connection to %s is not TLS", host)
	}
	certs := tlsConn.ConnectionState().PeerCertificates
	if len(certs) == 0 {
		return "", fmt.Errorf("no certificates found for %s", host)
	}

	rootCert := certs[len(certs)-1]
	fingerprint := sha1.Sum(rootCert.Raw) //nolint:gosec
	return fmt.Sprintf("%x", fingerprint), nil
}

func extractIssuerHost(issuerURL string) (string, error) {
	u, err := url.Parse(issuerURL)
	if err != nil {
		return "", fmt.Errorf("failed to parse issuer URL %s: %w", issuerURL, err)
	}
	return u.Host + u.Path, nil
}

// createSTSIAMRoles creates per-component IAM roles with trust policies
// that allow assuming the role via OIDC web identity tokens.
func createSTSIAMRoles(ctx context.Context, client *iam.Client, infraID, oidcProviderARN, issuerURL string, credReqs []credentialsrequest.CredentialRequest, tags []iamtypes.Tag) error {
	issuerHost, err := extractIssuerHost(issuerURL)
	if err != nil {
		return err
	}

	for _, cr := range credReqs {
		roleName := awstypes.STSRoleName(infraID, cr.SecretRefNamespace, cr.SecretRefName)
		logrus.Infof("Creating STS IAM role %s for %s/%s", roleName, cr.SecretRefNamespace, cr.SecretRefName)

		trustPolicy, err := buildTrustPolicy(oidcProviderARN, issuerHost, cr.SecretRefNamespace, cr.ServiceAccountNames)
		if err != nil {
			return fmt.Errorf("failed to build trust policy for role %s: %w", roleName, err)
		}

		if err := createOrUpdateRole(ctx, client, roleName, trustPolicy, tags); err != nil {
			return fmt.Errorf("failed to create role %s: %w", roleName, err)
		}

		awsSpec, ok := cr.ProviderSpec.(*credentialsrequest.AWSProviderSpec)
		if !ok {
			return fmt.Errorf("credential request %s/%s does not have an AWS provider spec", cr.SecretRefNamespace, cr.SecretRefName)
		}

		permPolicy, err := buildPermissionPolicy(awsSpec.StatementEntries)
		if err != nil {
			return fmt.Errorf("failed to build permission policy for role %s: %w", roleName, err)
		}

		policyName := roleName + "-policy"
		if _, err := client.PutRolePolicy(ctx, &iam.PutRolePolicyInput{
			RoleName:       aws.String(roleName),
			PolicyName:     aws.String(policyName),
			PolicyDocument: aws.String(permPolicy),
		}); err != nil {
			return fmt.Errorf("failed to put inline policy on role %s: %w", roleName, err)
		}
	}

	logrus.Infof("Created %d STS IAM roles", len(credReqs))
	return nil
}

func createOrUpdateRole(ctx context.Context, client *iam.Client, roleName, trustPolicy string, tags []iamtypes.Tag) error {
	_, err := client.GetRole(ctx, &iam.GetRoleInput{RoleName: aws.String(roleName)})
	if err != nil {
		var noSuchEntity *iamtypes.NoSuchEntityException
		if !errors.As(err, &noSuchEntity) {
			return fmt.Errorf("failed to check for existing role: %w", err)
		}

		if _, err := client.CreateRole(ctx, &iam.CreateRoleInput{
			RoleName:                 aws.String(roleName),
			AssumeRolePolicyDocument: aws.String(trustPolicy),
			Tags:                     tags,
		}); err != nil {
			return fmt.Errorf("failed to create role: %w", err)
		}
		return nil
	}

	if _, err := client.UpdateAssumeRolePolicy(ctx, &iam.UpdateAssumeRolePolicyInput{
		RoleName:       aws.String(roleName),
		PolicyDocument: aws.String(trustPolicy),
	}); err != nil {
		return fmt.Errorf("failed to update trust policy on existing role: %w", err)
	}
	return nil
}

const roleTrustPolicyTemplate = `{
  "Version": "2012-10-17",
  "Statement": [
    {
      "Effect": "Allow",
      "Principal": {
        "Federated": "%s"
      },
      "Action": "sts:AssumeRoleWithWebIdentity",
      "Condition": %s
    }
  ]
}`

const rolePermissionPolicyTemplate = `{
  "Version": "2012-10-17",
  "Statement": %s
}`

func buildTrustPolicy(oidcProviderARN, issuerHost, namespace string, serviceAccountNames []string) (string, error) {
	subjects := make([]string, 0, len(serviceAccountNames))
	for _, sa := range serviceAccountNames {
		subjects = append(subjects, fmt.Sprintf("system:serviceaccount:%s:%s", namespace, sa))
	}

	condition, err := json.Marshal(map[string]interface{}{
		"StringEquals": map[string]interface{}{
			issuerHost + ":sub": subjects,
		},
	})
	if err != nil {
		return "", fmt.Errorf("failed to marshal trust policy condition: %w", err)
	}

	return fmt.Sprintf(roleTrustPolicyTemplate, oidcProviderARN, string(condition)), nil
}

func buildPermissionPolicy(entries []ccov1.StatementEntry) (string, error) {
	type statement struct {
		Effect    string      `json:"Effect"`
		Action    []string    `json:"Action"`
		Resource  string      `json:"Resource"`
		Condition interface{} `json:"Condition,omitempty"`
	}

	stmts := make([]statement, 0, len(entries))
	for _, e := range entries {
		s := statement{
			Effect:   e.Effect,
			Action:   e.Action,
			Resource: e.Resource,
		}
		if len(e.PolicyCondition) > 0 {
			s.Condition = e.PolicyCondition
		}
		stmts = append(stmts, s)
	}

	stmtsJSON, err := json.Marshal(stmts)
	if err != nil {
		return "", fmt.Errorf("failed to marshal permission policy statements: %w", err)
	}

	return fmt.Sprintf(rolePermissionPolicyTemplate, string(stmtsJSON)), nil
}
