package gencrypto

import (
	"bytes"
	"context"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/x509"
	"encoding/base64"
	"encoding/pem"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"github.com/sirupsen/logrus"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"

	"github.com/openshift/installer/pkg/asset"
	"github.com/openshift/installer/pkg/asset/agent/joiner"
	"github.com/openshift/installer/pkg/asset/agent/workflow"
)

const (
	authTokenSecretNamespace = "openshift-config" //nolint:gosec // no sensitive info
	authTokenSecretName      = "agent-auth-token" //nolint:gosec // no sensitive info
	agentAuthKey             = "agentAuthToken"
	userAuthKey              = "userAuthToken"
	watcherAuthKey           = "watcherAuthToken"
	authTokenPublicDataKey   = "authTokenPublicKey"
	// AuthType holds the authenticator type for agent based installer.
	AuthType       = "agent-installer-local"
	agentPersona   = "agentAuth"
	userPersona    = "userAuth"
	watcherPersona = "watcherAuth"
)

// AuthConfig is an asset that generates ECDSA public/private keys, JWT token.
type AuthConfig struct {
	PublicKey, AgentAuthToken, UserAuthToken, WatcherAuthToken, AuthTokenExpiry, AuthType string
}

var _ asset.Asset = (*AuthConfig)(nil)

// Dependencies returns the assets on which the AuthConfig asset depends.
func (a *AuthConfig) Dependencies() []asset.Asset {
	return []asset.Asset{
		&workflow.AgentWorkflow{},
		&joiner.AddNodesConfig{},
	}
}

// Generate generates the auth config for agent installer APIs.
func (a *AuthConfig) Generate(_ context.Context, dependencies asset.Parents) error {
	agentWorkflow := &workflow.AgentWorkflow{}
	dependencies.Get(agentWorkflow)
	a.AuthType = AuthType

	publicKey, privateKey, err := keyPairPEM()
	if err != nil {
		return err
	}
	// Encode to Base64 (Standard encoding)
	encodedPubKeyPEM := base64.StdEncoding.EncodeToString([]byte(publicKey))

	a.PublicKey = encodedPubKeyPEM

	tokens := map[string]*string{
		agentPersona:   &a.AgentAuthToken,
		userPersona:    &a.UserAuthToken,
		watcherPersona: &a.WatcherAuthToken,
	}

	switch agentWorkflow.Workflow {
	case workflow.AgentWorkflowTypeInstall:
		// Auth tokens do not expire
		for persona, tokenField := range tokens {
			token, err := generateToken(persona, privateKey, nil)
			if err != nil {
				return err
			}
			*tokenField = token
		}

	case workflow.AgentWorkflowTypeAddNodes:
		addNodesConfig := &joiner.AddNodesConfig{}
		dependencies.Get(addNodesConfig)

		// Auth tokens expires after 48 hours
		expiry := time.Now().UTC().Add(48 * time.Hour)
		a.AuthTokenExpiry = expiry.Format(time.RFC3339)

		for persona, tokenField := range tokens {
			token, err := generateToken(persona, privateKey, &expiry)
			if err != nil {
				return err
			}
			*tokenField = token
		}

		err = a.createOrUpdateAuthTokenSecret(addNodesConfig.Params.Kubeconfig)
		if err != nil {
			return err
		}
	default:
		return fmt.Errorf("AgentWorkflowType value not supported: %s", agentWorkflow.Workflow)
	}
	return nil
}

// Name returns the human-friendly name of the asset.
func (*AuthConfig) Name() string {
	return "Agent Installer API Auth Config"
}

// keyPairPEM returns the public, private keys in PEM format.
func keyPairPEM() (string, string, error) {
	priv, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		return "", "", err
	}

	// encode private key to PEM string
	privBytes, err := x509.MarshalECPrivateKey(priv)
	if err != nil {
		return "", "", err
	}

	block := &pem.Block{
		Type:  "EC PRIVATE KEY",
		Bytes: privBytes,
	}

	var privKeyPEM bytes.Buffer
	err = pem.Encode(&privKeyPEM, block)
	if err != nil {
		return "", "", err
	}

	// encode public key to PEM string
	pubBytes, err := x509.MarshalPKIXPublicKey(priv.Public())
	if err != nil {
		return "", "", err
	}

	block = &pem.Block{
		Type:  "EC PUBLIC KEY",
		Bytes: pubBytes,
	}

	var pubKeyPEM bytes.Buffer
	err = pem.Encode(&pubKeyPEM, block)
	if err != nil {
		return "", "", err
	}

	return pubKeyPEM.String(), privKeyPEM.String(), nil
}

// generateToken returns a JWT token based on the private key.
func generateToken(userPersona string, privateKeyPem string, expiry *time.Time) (string, error) {
	// Create the JWT claims
	claims := jwt.MapClaims{
		"auth_scheme": userPersona,
	}
	// Set the expiry time if provided
	if expiry != nil {
		claims["exp"] = expiry.Unix()
	}

	// Create the token using the ES256 signing method and the claims
	token := jwt.NewWithClaims(jwt.SigningMethodES256, claims)

	priv, err := jwt.ParseECPrivateKeyFromPEM([]byte(privateKeyPem))
	if err != nil {
		return "", err
	}
	// Sign the token with the provided private key
	tokenString, err := token.SignedString(priv)
	if err != nil {
		return "", err
	}
	return tokenString, nil
}

func initClient(kubeconfig string) (*kubernetes.Clientset, error) {
	var err error
	var config *rest.Config
	if kubeconfig != "" {
		config, err = clientcmd.BuildConfigFromFlags("", kubeconfig)
	} else {
		config, err = rest.InClusterConfig()
	}
	if err != nil {
		return nil, err
	}

	k8sclientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		return nil, err
	}

	return k8sclientset, err
}

func (a *AuthConfig) createOrUpdateAuthTokenSecret(kubeconfigPath string) error {
	k8sclientset, err := initClient(kubeconfigPath)
	if err != nil {
		return err
	}
	// check if secret exists
	retrievedSecret, err := k8sclientset.CoreV1().Secrets(authTokenSecretNamespace).Get(context.Background(), authTokenSecretName, metav1.GetOptions{})
	// if the secret does not exist
	if err != nil {
		if errors.IsNotFound(err) {
			return a.createSecret(k8sclientset)
		}
		// Other errors while trying to get the secret
		return fmt.Errorf("unable to retrieve secret %s/%s: %w", authTokenSecretNamespace, authTokenSecretName, err)
	}

	// if the secret exists in the cluster, get the token
	retrievedAgentAuthToken, retrievedUserAuthToken, retrievedWatcherAuthToken, err := extractAuthTokensFromSecret(retrievedSecret)
	if err != nil {
		return err
	}
	// All auth tokens expire at the same time so we could only check any 1 token to get the expiry time
	expiryTime, err := ParseExpirationFromToken(retrievedAgentAuthToken)
	if err != nil {
		return err
	}
	// Calculate 24 hours before the expiration time
	thresholdTime := expiryTime.UTC().Add(-24 * time.Hour)
	// Check if current time is past the thresholdTime time of 24 hours
	if time.Now().UTC().After(thresholdTime) {
		// update the secret in the cluster with a new token from asset store
		err = a.refreshAuthTokenSecret(k8sclientset, retrievedSecret)
		if err != nil {
			return err
		}
	} else {
		// Update the token in asset store with the retrieved token from the cluster
		a.AgentAuthToken = retrievedAgentAuthToken
		a.UserAuthToken = retrievedUserAuthToken
		a.WatcherAuthToken = retrievedWatcherAuthToken
		// get the token expiry time of the retrieved token from the cluster
		a.AuthTokenExpiry = expiryTime.UTC().Format(time.RFC3339)

		retrievedPublicKey, err := extractPublicKeyFromSecret(retrievedSecret)
		if err != nil {
			return err
		}
		// Update the asset store with the retrieved public key associated with the valid token from the cluster
		a.PublicKey = retrievedPublicKey
		logrus.Infof("Reusing existing auth token (valid up to %s)", a.AuthTokenExpiry)
	}
	return err
}

func (a *AuthConfig) createSecret(k8sclientset kubernetes.Interface) error {
	// Create a Secret
	secret := &corev1.Secret{
		ObjectMeta: metav1.ObjectMeta{
			Name: authTokenSecretName,
			// only for informational purposes
			Annotations: map[string]string{
				"updatedAt": "", // Initially set to empty
				"expiresAt": a.AuthTokenExpiry,
			},
		},
		Type: corev1.SecretTypeOpaque,
		Data: map[string][]byte{
			agentAuthKey:           []byte(a.AgentAuthToken),
			userAuthKey:            []byte(a.UserAuthToken),
			watcherAuthKey:         []byte(a.WatcherAuthToken),
			authTokenPublicDataKey: []byte(a.PublicKey),
		},
	}
	_, err := k8sclientset.CoreV1().Secrets(authTokenSecretNamespace).Create(context.Background(), secret, metav1.CreateOptions{})
	if err != nil {
		return fmt.Errorf("failed to create secret: %w", err)
	}
	logrus.Infof("Generated auth token (valid up to %s)", a.AuthTokenExpiry)
	logrus.Infof("Created secret %s/%s", authTokenSecretNamespace, authTokenSecretName)

	return nil
}

func (a *AuthConfig) refreshAuthTokenSecret(k8sclientset kubernetes.Interface, retrievedSecret *corev1.Secret) error {
	retrievedSecret.Data[agentAuthKey] = []byte(a.AgentAuthToken)
	// Update userAuthKey and watcherAuthKey only if they exist in the secret
	if _, exists := retrievedSecret.Data[userAuthKey]; exists {
		retrievedSecret.Data[userAuthKey] = []byte(a.UserAuthToken)
	}
	if _, exists := retrievedSecret.Data[watcherAuthKey]; exists {
		retrievedSecret.Data[watcherAuthKey] = []byte(a.WatcherAuthToken)
	}
	retrievedSecret.Data[authTokenPublicDataKey] = []byte(a.PublicKey)
	// only for informational purposes
	retrievedSecret.Annotations["updatedAt"] = time.Now().UTC().Format(time.RFC3339)
	retrievedSecret.Annotations["expiresAt"] = a.AuthTokenExpiry

	_, err := k8sclientset.CoreV1().Secrets(authTokenSecretNamespace).Update(context.TODO(), retrievedSecret, metav1.UpdateOptions{})
	if err != nil {
		return err
	}
	logrus.Infof("Auth token regenerated (valid up to %s)", a.AuthTokenExpiry)
	logrus.Infof("Updated secret %s/%s", authTokenSecretNamespace, authTokenSecretName)
	return nil
}

// GetWatcherAuthTokenFromCluster returns a watcherAuth token string stored as the secret from the cluster.
func GetWatcherAuthTokenFromCluster(ctx context.Context, kubeconfigPath string) (string, error) {
	client, err := initClient(kubeconfigPath)
	if err != nil {
		return "", err
	}

	retrievedSecret, err := client.CoreV1().Secrets(authTokenSecretNamespace).Get(ctx, authTokenSecretName, metav1.GetOptions{})
	if err != nil {
		return "", err
	}
	_, _, watcherAuthToken, err := extractAuthTokensFromSecret(retrievedSecret)
	if err != nil {
		return "", err
	}
	return watcherAuthToken, err
}

func extractAuthTokensFromSecret(secret *corev1.Secret) (string, string, string, error) {
	// Check for agentAuthKey, which must exist in both old (4.17) and new versions (4.18+)
	existingAgentAuthToken, agentAuthTokenExists := secret.Data[agentAuthKey]
	if !agentAuthTokenExists || len(existingAgentAuthToken) == 0 {
		return "", "", "", fmt.Errorf("auth token secret %s/%s does not contain the key %s or is empty", authTokenSecretNamespace, authTokenSecretName, agentAuthKey)
	}

	existingUserAuthToken, userAuthTokenExists := secret.Data[userAuthKey]
	existingWatcherAuthToken, watcherAuthTokenExists := secret.Data[watcherAuthKey]

	// Handle old version compatibility for OCP 4.17
	if !userAuthTokenExists && !watcherAuthTokenExists {
		// For old version OCP 4.17, where only agentAuthToken is present
		return string(existingAgentAuthToken), "", "", nil
	}

	// Handle cases where new keys are missing in OCP 4.18+
	if (!userAuthTokenExists || len(existingUserAuthToken) == 0) || (!watcherAuthTokenExists || len(existingWatcherAuthToken) == 0) {
		return "", "", "", fmt.Errorf("auth token secret %s/%s is missing one or more required keys (%s, %s, or %s) or they are empty",
			authTokenSecretNamespace, authTokenSecretName, agentAuthKey, userAuthKey, watcherAuthKey)
	}

	// Return all keys if present
	return string(existingAgentAuthToken), string(existingUserAuthToken), string(existingWatcherAuthToken), nil
}

func extractPublicKeyFromSecret(secret *corev1.Secret) (string, error) {
	existingPublicKey, exists := secret.Data[authTokenPublicDataKey]
	if !exists || len(existingPublicKey) == 0 {
		return "", fmt.Errorf("auth token secret %s/%s does not contain the key %s or is empty", authTokenSecretNamespace, authTokenSecretName, authTokenPublicDataKey)
	}
	return string(existingPublicKey), nil
}
