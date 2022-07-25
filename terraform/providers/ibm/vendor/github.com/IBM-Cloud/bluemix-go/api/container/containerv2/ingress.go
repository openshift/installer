package containerv2

import (
	"fmt"
	"strconv"

	"github.com/IBM-Cloud/bluemix-go/client"
)

// Secret struct holding details for a single secret
type Secret struct {
	Cluster     string `json:"cluster" description:"name of secret"`
	Name        string `json:"name" description:"name of secret"`
	Namespace   string `json:"namespace" description:"namespace of secret"`
	Domain      string `json:"domain" description:"domain the cert belongs to"`
	CRN         string `json:"crn" description:"crn of the certificate in certificate manager"`
	ExpiresOn   string `json:"expiresOn" description:"expiration date of the certificate"`
	Status      string `json:"status" description:"status of  Will be used for displaying callback operations to user"`
	UserManaged bool   `json:"userManaged" description:"true or false. Used to show which certs and secrets are system generated and which are not"`
	Persistence bool   `json:"persistence" description:"true or false. Persist the secret even if a user attempts to delete it"`
}

// Secrets struct for a secret array
type Secrets []Secret

// SecretCreateConfig the secret create request
type SecretCreateConfig struct {
	Cluster     string `json:"cluster" description:"name of secret" binding:"required"`
	Name        string `json:"name" description:"name of secret" binding:"required"`
	Namespace   string `json:"namespace" description:"namespace of  Optional, if none specified it will be placed in the ibm-cert-store namespace"`
	CRN         string `json:"crn" description:"crn of the certificate in certificate manager"`
	Persistence bool   `json:"persistence" description:"true or false. Persist the secret even if a user attempts to delete it"`
}

// SecretDeleteConfig the secret delete request
type SecretDeleteConfig struct {
	Cluster   string `json:"cluster" description:"name of secret" binding:"required"`
	Name      string `json:"name" description:"name of secret" binding:"required"`
	Namespace string `json:"namespace" description:"namespace of secret" binding:"required"`
}

// SecretUpdateConfig secret update request
type SecretUpdateConfig struct {
	Cluster   string `json:"cluster" description:"name of secret" binding:"required"`
	Name      string `json:"name" description:"name of secret" binding:"required"`
	Namespace string `json:"namespace" description:"namespace of secret" binding:"required"`
	CRN       string `json:"crn" description:"crn of the certificate in certificate manager"`
}

type ingress struct {
	client *client.Client
}

//Ingress interface
type Ingress interface {
	CreateIngressSecret(req SecretCreateConfig) (response Secret, err error)
	UpdateIngressSecret(req SecretUpdateConfig) (response Secret, err error)
	DeleteIngressSecret(req SecretDeleteConfig) (err error)
	GetIngressSecretList(clusterNameOrID string, showDeleted bool) (response Secrets, err error)
	GetIngressSecret(clusterNameOrID, secretName, secretNamespace string) (response Secret, err error)
}

func newIngressAPI(c *client.Client) Ingress {
	return &ingress{
		client: c,
	}
}

// GetIngressSecretList returns a list of ingress secrets for a given cluster
func (r *ingress) GetIngressSecretList(clusterNameOrID string, showDeleted bool) (response Secrets, err error) {
	deleted := strconv.FormatBool(showDeleted)
	_, err = r.client.Get(fmt.Sprintf("/ingress/v2/secret/getSecrets?cluster=%s&showDeleted=%s", clusterNameOrID, deleted), &response)
	return
}

// GetIngressSecret returns a single ingress secret in a given cluster
func (r *ingress) GetIngressSecret(clusterNameOrID, secretName, secretNamespace string) (response Secret, err error) {
	_, err = r.client.Get(fmt.Sprintf("/ingress/v2/secret/getSecret?cluster=%s&name=%s&namespace=%s", clusterNameOrID, secretName, secretNamespace), &response)
	return
}

// CreateIngressSecret creates an ingress secret with the given name in the given namespace
func (r *ingress) CreateIngressSecret(req SecretCreateConfig) (response Secret, err error) {
	_, err = r.client.Post("/ingress/v2/secret/createSecret", req, &response)
	return
}

// UpdateIngressSecret updates an existing secret with new cert values
func (r *ingress) UpdateIngressSecret(req SecretUpdateConfig) (response Secret, err error) {
	_, err = r.client.Post("/ingress/v2/secret/updateSecret", req, &response)
	return
}

// DeleteIngressSecret deletes the ingress secret from the cluster
func (r *ingress) DeleteIngressSecret(req SecretDeleteConfig) (err error) {
	_, err = r.client.Post("/ingress/v2/secret/deleteSecret", req, nil)
	return
}
