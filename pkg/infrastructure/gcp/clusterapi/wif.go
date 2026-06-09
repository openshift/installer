package clusterapi

import (
	"context"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"encoding/base64"
	"encoding/json"
	"encoding/pem"
	"errors"
	"fmt"
	"math/big"
	"net/http"
	"slices"
	"time"

	"cloud.google.com/go/storage"
	"github.com/sirupsen/logrus"
	resourcemanager "google.golang.org/api/cloudresourcemanager/v3"
	"google.golang.org/api/googleapi"
	"google.golang.org/api/iam/v1"
	"google.golang.org/api/option"
	"k8s.io/apimachinery/pkg/util/wait"

	icgcp "github.com/openshift/installer/pkg/asset/installconfig/gcp"
	"github.com/openshift/installer/pkg/infrastructure/clusterapi"
	gcptypes "github.com/openshift/installer/pkg/types/gcp"
)

const (
	wifRetryTime  = 10 * time.Second
	wifRetryCount = 12
)

// WIFPoolName returns the deterministic WIF pool ID for a cluster.
func WIFPoolName(infraID string) string {
	return fmt.Sprintf("%s-wif-pool", infraID)
}

// WIFProviderName returns the deterministic WIF OIDC provider ID for a cluster.
func WIFProviderName(infraID string) string {
	return fmt.Sprintf("%s-oidc-provider", infraID)
}

// OIDCBucketName returns the deterministic GCS bucket name for OIDC discovery.
func OIDCBucketName(infraID string) string {
	return fmt.Sprintf("%s-oidc", infraID)
}

// OIDCIssuerURL returns the deterministic OIDC issuer URL for a cluster.
func OIDCIssuerURL(infraID string) string {
	return fmt.Sprintf("https://storage.googleapis.com/%s", OIDCBucketName(infraID))
}

// ProvisionWIF creates WIF infrastructure for installer-provisioned mode.
// Requires a user-provided bound SA signing key for JWKS generation.
func ProvisionWIF(ctx context.Context, in clusterapi.PreProvisionInput) error {
	if in.BoundSASigningKey == nil || len(in.BoundSASigningKey.Files()) == 0 {
		return fmt.Errorf("bound service account signing key is required for WIF provisioning; provide bound-service-account-signing-key.key in the asset directory")
	}

	var publicKeyPEM []byte
	for _, f := range in.BoundSASigningKey.Files() {
		if f.Filename == "tls/bound-service-account-signing-key.pub" {
			publicKeyPEM = f.Data
			break
		}
	}
	if len(publicKeyPEM) == 0 {
		return fmt.Errorf("failed to find bound SA signing public key")
	}

	platform := in.InstallConfig.Config.Platform.GCP
	projectID := platform.ProjectID
	infraID := in.InfraID
	region := platform.Region

	iamOpts := []option.ClientOption{}
	storageOpts := []option.ClientOption{}
	crmOpts := []option.ClientOption{}
	if gcptypes.ShouldUseEndpointForInstaller(platform.Endpoint) {
		iamOpts = append(iamOpts, icgcp.CreateEndpointOption(platform.Endpoint.Name, icgcp.ServiceNameGCPIAM))
		storageOpts = append(storageOpts, icgcp.CreateEndpointOption(platform.Endpoint.Name, icgcp.ServiceNameGCPStorage))
		crmOpts = append(crmOpts, icgcp.CreateEndpointOption(platform.Endpoint.Name, icgcp.ServiceNameGCPCloudResource))
	}

	iamSvc, err := icgcp.GetIAMService(ctx, iamOpts...)
	if err != nil {
		return fmt.Errorf("failed to create IAM service: %w", err)
	}

	storageClient, err := icgcp.GetStorageService(ctx, storageOpts...)
	if err != nil {
		return fmt.Errorf("failed to create storage client: %w", err)
	}

	crmSvc, err := icgcp.GetCloudResourceService(ctx, crmOpts...)
	if err != nil {
		return fmt.Errorf("failed to create resource manager service: %w", err)
	}

	projectNumber, err := GetProjectNumber(ctx, crmSvc, projectID)
	if err != nil {
		return fmt.Errorf("failed to get project number: %w", err)
	}

	poolName := WIFPoolName(infraID)
	providerName := WIFProviderName(infraID)
	issuerURL := OIDCIssuerURL(infraID)

	logrus.Infof("Creating workload identity pool %s", poolName)
	if err := createWorkloadIdentityPool(ctx, iamSvc, projectID, infraID, poolName); err != nil {
		return fmt.Errorf("failed to create workload identity pool: %w", err)
	}

	logrus.Infof("Creating OIDC discovery bucket %s", OIDCBucketName(infraID))
	if err := createOIDCBucket(ctx, storageClient, projectID, region, infraID, issuerURL, publicKeyPEM); err != nil {
		return fmt.Errorf("failed to create OIDC bucket: %w", err)
	}

	logrus.Infof("Creating OIDC provider %s", providerName)
	if err := createOIDCProvider(ctx, iamSvc, projectID, projectNumber, infraID, poolName, providerName, issuerURL); err != nil {
		return fmt.Errorf("failed to create OIDC provider: %w", err)
	}

	logrus.Infof("Binding service accounts to workload identity pool")
	masterSAEmail := gcptypes.GetDefaultServiceAccount(platform, infraID, "master")
	workerSAEmail := gcptypes.GetDefaultServiceAccount(platform, infraID, "worker")

	bindings := []struct {
		k8sNamespace string
		k8sSAName    string
		gcpSAEmail   string
	}{
		{"openshift-cloud-credential-operator", "cloud-credential-operator", masterSAEmail},
		{"openshift-image-registry", "cluster-image-registry-operator", masterSAEmail},
		{"openshift-ingress-operator", "ingress-operator", masterSAEmail},
		{"openshift-machine-api", "machine-api-controllers", masterSAEmail},
		{"openshift-cloud-controller-manager", "cloud-controller-manager", masterSAEmail},
		{"openshift-cluster-csi-drivers", "gcp-pd-csi-driver-operator", workerSAEmail},
	}

	for _, b := range bindings {
		if err := bindServiceAccountToWIF(ctx, iamSvc, projectID, projectNumber, poolName, b.k8sNamespace, b.k8sSAName, b.gcpSAEmail); err != nil {
			return fmt.Errorf("failed to bind SA %s/%s to %s: %w", b.k8sNamespace, b.k8sSAName, b.gcpSAEmail, err)
		}
	}

	logrus.Infof("WIF provisioning complete")
	return nil
}

// GetProjectNumber resolves a project ID to its numeric project number.
func GetProjectNumber(ctx context.Context, crmSvc *resourcemanager.Service, projectID string) (string, error) {
	ctx, cancel := context.WithTimeout(ctx, time.Minute)
	defer cancel()

	project, err := crmSvc.Projects.Get(fmt.Sprintf("projects/%s", projectID)).Context(ctx).Do()
	if err != nil {
		return "", fmt.Errorf("failed to get project %s: %w", projectID, err)
	}

	// project.Name is "projects/{number}"
	if len(project.Name) > 9 {
		return project.Name[9:], nil
	}
	return "", fmt.Errorf("unexpected project name format: %s", project.Name)
}

// GetBYOIssuerURL describes an existing WIF provider and returns its OIDC issuer URL.
func GetBYOIssuerURL(ctx context.Context, iamSvc *iam.Service, projectID, poolID, providerID string) (string, error) {
	ctx, cancel := context.WithTimeout(ctx, time.Minute)
	defer cancel()

	name := fmt.Sprintf("projects/%s/locations/global/workloadIdentityPools/%s/providers/%s",
		projectID, poolID, providerID)
	provider, err := iamSvc.Projects.Locations.WorkloadIdentityPools.Providers.Get(name).Context(ctx).Do()
	if err != nil {
		return "", fmt.Errorf("failed to get WIF provider %s: %w", name, err)
	}
	if provider.Oidc == nil || provider.Oidc.IssuerUri == "" {
		return "", fmt.Errorf("WIF provider %s does not have an OIDC issuer URI", name)
	}
	return provider.Oidc.IssuerUri, nil
}

func createWorkloadIdentityPool(ctx context.Context, iamSvc *iam.Service, projectID, infraID, poolName string) error {
	ctx, cancel := context.WithTimeout(ctx, 2*time.Minute)
	defer cancel()

	parent := fmt.Sprintf("projects/%s/locations/global", projectID)
	pool := &iam.WorkloadIdentityPool{
		DisplayName: fmt.Sprintf("%s WIF Pool", infraID),
		Description: "Created by OpenShift Installer",
	}

	op, err := iamSvc.Projects.Locations.WorkloadIdentityPools.Create(parent, pool).
		WorkloadIdentityPoolId(poolName).Context(ctx).Do()
	if err != nil {
		if isAlreadyExists(err) {
			logrus.Debugf("WIF pool %s already exists, continuing", poolName)
			return nil
		}
		return fmt.Errorf("failed to create WIF pool: %w", err)
	}

	return waitForIAMOperation(ctx, iamSvc, op.Name)
}

func createOIDCBucket(ctx context.Context, storageClient *storage.Client, projectID, region, infraID, issuerURL string, publicKeyPEM []byte) error {
	bucketName := OIDCBucketName(infraID)
	bucket := storageClient.Bucket(bucketName)

	if err := bucket.Create(ctx, projectID, &storage.BucketAttrs{
		Location:                 region,
		UniformBucketLevelAccess: storage.UniformBucketLevelAccess{Enabled: true},
		Labels:                   map[string]string{"kubernetes-io-cluster-" + infraID: "owned"},
	}); err != nil {
		return fmt.Errorf("failed to create OIDC bucket %s: %w", bucketName, err)
	}

	policy, err := bucket.IAM().Policy(ctx)
	if err != nil {
		return fmt.Errorf("failed to get OIDC bucket IAM policy: %w", err)
	}
	policy.Add("allUsers", "roles/storage.objectViewer")
	if err := bucket.IAM().SetPolicy(ctx, policy); err != nil {
		return fmt.Errorf("failed to set OIDC bucket IAM policy: %w", err)
	}

	discoveryDoc := generateOIDCDiscoveryDoc(issuerURL)
	wtr := bucket.Object(".well-known/openid-configuration").NewWriter(ctx)
	wtr.ContentType = "application/json"
	if _, err := wtr.Write(discoveryDoc); err != nil {
		return fmt.Errorf("failed to write OIDC discovery doc: %w", err)
	}
	if err := wtr.Close(); err != nil {
		return fmt.Errorf("failed to close OIDC discovery doc writer: %w", err)
	}

	jwksData, err := GenerateJWKS(publicKeyPEM)
	if err != nil {
		return fmt.Errorf("failed to generate JWKS: %w", err)
	}

	jwksWriter := bucket.Object("keys.json").NewWriter(ctx)
	jwksWriter.ContentType = "application/json"
	if _, err := jwksWriter.Write(jwksData); err != nil {
		return fmt.Errorf("failed to write JWKS: %w", err)
	}
	if err := jwksWriter.Close(); err != nil {
		return fmt.Errorf("failed to close JWKS writer: %w", err)
	}

	return nil
}

func createOIDCProvider(ctx context.Context, iamSvc *iam.Service, projectID, projectNumber, infraID, poolName, providerName, issuerURL string) error {
	ctx, cancel := context.WithTimeout(ctx, 2*time.Minute)
	defer cancel()

	parent := fmt.Sprintf("projects/%s/locations/global/workloadIdentityPools/%s", projectID, poolName)
	audience := BuildAudienceURI(projectNumber, poolName, providerName)

	provider := &iam.WorkloadIdentityPoolProvider{
		DisplayName: fmt.Sprintf("%s OIDC Provider", infraID),
		Description: "Created by OpenShift Installer",
		Oidc: &iam.Oidc{
			IssuerUri:        issuerURL,
			AllowedAudiences: []string{audience},
		},
		AttributeMapping: map[string]string{
			"google.subject": "assertion.sub",
		},
	}

	op, err := iamSvc.Projects.Locations.WorkloadIdentityPools.Providers.Create(parent, provider).
		WorkloadIdentityPoolProviderId(providerName).Context(ctx).Do()
	if err != nil {
		if isAlreadyExists(err) {
			logrus.Debugf("OIDC provider %s already exists, continuing", providerName)
			return nil
		}
		return fmt.Errorf("failed to create OIDC provider: %w", err)
	}

	return waitForIAMOperation(ctx, iamSvc, op.Name)
}

func bindServiceAccountToWIF(ctx context.Context, iamSvc *iam.Service, projectID, projectNumber, poolName, k8sNamespace, k8sSAName, gcpSAEmail string) error {
	ctx, cancel := context.WithTimeout(ctx, time.Minute)
	defer cancel()

	saResource := fmt.Sprintf("projects/%s/serviceAccounts/%s", projectID, gcpSAEmail)
	member := fmt.Sprintf(
		"principal://iam.googleapis.com/projects/%s/locations/global/workloadIdentityPools/%s/subject/system:serviceaccount:%s:%s",
		projectNumber, poolName, k8sNamespace, k8sSAName,
	)

	policy, err := iamSvc.Projects.ServiceAccounts.GetIamPolicy(saResource).Context(ctx).Do()
	if err != nil {
		return fmt.Errorf("failed to get SA IAM policy for %s: %w", gcpSAEmail, err)
	}

	role := "roles/iam.workloadIdentityUser"
	var binding *iam.Binding
	for _, b := range policy.Bindings {
		if b.Role == role {
			if slices.Contains(b.Members, member) {
				logrus.Debugf("WIF binding already exists for %s/%s on %s", k8sNamespace, k8sSAName, gcpSAEmail)
				return nil
			}
			binding = b
			break
		}
	}

	if binding == nil {
		binding = &iam.Binding{
			Role:    role,
			Members: []string{member},
		}
		policy.Bindings = append(policy.Bindings, binding)
	} else {
		binding.Members = append(binding.Members, member)
	}

	_, err = iamSvc.Projects.ServiceAccounts.SetIamPolicy(saResource, &iam.SetIamPolicyRequest{
		Policy: policy,
	}).Context(ctx).Do()
	if err != nil {
		return fmt.Errorf("failed to set SA IAM policy for %s: %w", gcpSAEmail, err)
	}

	logrus.Debugf("Bound %s/%s to %s via WIF", k8sNamespace, k8sSAName, gcpSAEmail)
	return nil
}

func waitForIAMOperation(ctx context.Context, iamSvc *iam.Service, opName string) error {
	backoff := wait.Backoff{
		Duration: wifRetryTime,
		Factor:   1.5,
		Jitter:   0.1,
		Steps:    wifRetryCount,
	}

	var lastErr error
	if waitErr := wait.ExponentialBackoffWithContext(ctx, backoff, func(ctx context.Context) (bool, error) {
		op, err := iamSvc.Projects.Locations.WorkloadIdentityPools.Operations.Get(opName).Context(ctx).Do()
		if err != nil {
			lastErr = err
			return false, nil
		}
		if op.Done {
			if op.Error != nil {
				return false, fmt.Errorf("operation %s failed: %s", opName, op.Error.Message)
			}
			return true, nil
		}
		return false, nil
	}); waitErr != nil {
		if wait.Interrupted(waitErr) {
			return fmt.Errorf("timed out waiting for operation %s: %w", opName, lastErr)
		}
		return waitErr
	}
	return nil
}

func generateOIDCDiscoveryDoc(issuerURL string) []byte {
	doc := map[string]any{
		"issuer":                                issuerURL,
		"jwks_uri":                              fmt.Sprintf("%s/keys.json", issuerURL),
		"response_types_supported":              []string{"id_token"},
		"subject_types_supported":               []string{"public"},
		"id_token_signing_alg_values_supported": []string{"RS256"},
	}
	data, _ := json.MarshalIndent(doc, "", "  ")
	return data
}

// GenerateJWKS produces a JSON Web Key Set from an RSA public key PEM.
func GenerateJWKS(publicKeyPEM []byte) ([]byte, error) {
	block, _ := pem.Decode(publicKeyPEM)
	if block == nil {
		return nil, fmt.Errorf("failed to decode PEM block")
	}

	pub, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		return nil, fmt.Errorf("failed to parse public key: %w", err)
	}

	rsaPub, ok := pub.(*rsa.PublicKey)
	if !ok {
		return nil, fmt.Errorf("public key is not RSA")
	}

	// Compute KID as base64url(SHA256(DER-encoded public key))
	derBytes, err := x509.MarshalPKIXPublicKey(rsaPub)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal public key: %w", err)
	}
	hash := sha256.Sum256(derBytes)
	kid := base64.RawURLEncoding.EncodeToString(hash[:])

	jwks := map[string]any{
		"keys": []map[string]string{
			{
				"kty": "RSA",
				"alg": "RS256",
				"use": "sig",
				"kid": kid,
				"n":   base64.RawURLEncoding.EncodeToString(rsaPub.N.Bytes()),
				"e":   base64.RawURLEncoding.EncodeToString(big.NewInt(int64(rsaPub.E)).Bytes()),
			},
		},
	}

	return json.Marshal(jwks)
}

func isAlreadyExists(err error) bool {
	var ae *googleapi.Error
	return errors.As(err, &ae) && ae.Code == http.StatusConflict
}

// BuildAudienceURI constructs the WIF audience URI from project number, pool ID, and provider ID.
func BuildAudienceURI(projectNumber, poolID, providerID string) string {
	return fmt.Sprintf(
		"//iam.googleapis.com/projects/%s/locations/global/workloadIdentityPools/%s/providers/%s",
		projectNumber, poolID, providerID,
	)
}
