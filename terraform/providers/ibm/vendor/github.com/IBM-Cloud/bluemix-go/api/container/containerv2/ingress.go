package containerv2

import (
	"fmt"
	"strconv"

	"github.com/IBM-Cloud/bluemix-go/client"
)

// Secret struct holding details for a single secret
type Secret struct {
	Cluster              string `json:"cluster" description:"name of secret"`
	Name                 string `json:"name" description:"name of secret"`
	Namespace            string `json:"namespace" description:"namespace of secret"`
	Domain               string `json:"domain" description:"domain the cert belongs to"`
	CRN                  string `json:"crn" description:"crn of the certificate in certificate manager"`
	ExpiresOn            string `json:"expiresOn" description:"expiration date of the certificate"`
	Status               string `json:"status" description:"status of secret. Will be used for displaying callback operations to user"`
	UserManaged          bool   `json:"userManaged" description:"true or false. Used to show which certs and secrets are system generated and which are not"`
	Persistence          bool   `json:"persistence" description:"true or false. Persist the secret even if a user attempts to delete it"`
	Type                 string `json:"type" description:"supported types include TLS and Opaque"`
	SecretType           string `json:"secretType" description:"secrets manager type for secret"`
	LastUpdatedTimestamp string `json:"lastUpdatedTimestamp" description:"last updated timestamp for type tls secrets"`
	Fields               Fields `json:"fields" description:"fields in secret"`
}

// Secret struct holding details for a single secret field
type Field struct {
	Name                 string `json:"name" description:"name of secret field"`
	CRN                  string `json:"crn" description:"crn of secret field"`
	ExpiresOn            string `json:"expiresOn" description:"expiration date of the secret"`
	SecretType           string `json:"secretType" description:"secrets manager type for secret"`
	LastUpdatedTimestamp string `json:"lastUpdatedTimestamp" description:"last updated timestamp for type tls secrets"`
}

// Secrets struct for a secret array
type Fields []Field

// Secrets struct for a secret array
type Secrets []Secret

// SecretCreateConfig the secret create request
type SecretCreateConfig struct {
	Cluster     string     `json:"cluster" description:"name of secret" binding:"required"`
	Name        string     `json:"name" description:"name of secret" binding:"required"`
	Namespace   string     `json:"namespace" description:"namespace of secret. Optional, if none specified it will be placed in the ibm-cert-store namespace"`
	CRN         string     `json:"crn" description:"crn of the certificate in secret manager"`
	Persistence bool       `json:"persistence" description:"true or false. Persist the secret even if a user attempts to delete it"`
	Type        string     `json:"type" description:"TLS or Opaque. Defaults to TLS if none specified."`
	FieldsToAdd []FieldAdd `json:"add" description:"fields to add to secret of type opaque."`
}

// FieldAdd the secret field add request
type FieldAdd struct {
	Name         string `json:"name" description:"name of secret field. Cannot append prefix when setting this."`
	CRN          string `json:"crn" description:"crn of secret field"`
	AppendPrefix bool   `json:"append_prefix" description:"true or false. Append the secret name in secret manager as a prefix to secret type. Cannot set name when appending prefix."`
}

// FieldRemove the secret field remove request
type FieldRemove struct {
	Name string `json:"name" description:"name of secret field"`
}

// SecretDeleteConfig the secret delete request
type SecretDeleteConfig struct {
	Cluster   string `json:"cluster" description:"name of secret" binding:"required"`
	Name      string `json:"name" description:"name of secret" binding:"required"`
	Namespace string `json:"namespace" description:"namespace of secret" binding:"required"`
}

// SecretUpdateConfig secret update request
type SecretUpdateConfig struct {
	Cluster        string        `json:"cluster" description:"name of secret" binding:"required"`
	Name           string        `json:"name" description:"name of secret" binding:"required"`
	Namespace      string        `json:"namespace" description:"namespace of secret" binding:"required"`
	CRN            string        `json:"crn" description:"crn of the certificate in secret manager"`
	FieldsToAdd    []FieldAdd    `json:"add" description:"fields to add to secret"`
	FieldsToRemove []FieldRemove `json:"remove" description:"fields to remove from secret"`
}

// Instance struct holding details for a single instance
type Instance struct {
	Cluster         string `json:"cluster" description:"id of cluster"`
	Name            string `json:"name" description:"name of instance"`
	CRN             string `json:"crn" description:"crn of the instance"`
	SecretGroupID   string `json:"secretGroupID" description:"ID of the secret group where secrets will be stored"`
	SecretGroupName string `json:"secretGroupName" description:"name of the secret group where secrets will be stored"`
	CallbackChannel string `json:"callbackChannel" description:"callback channel of the instance"`
	UserManaged     bool   `json:"userManaged" description:"true or false. Used to show which certs and secrets are system generated and which are not"`
	IsDefault       bool   `json:"isDefault" description:"true or false. Used to show which instance subdomains certificates are uploaded into"`
	Type            string `json:"type" description:"designates instance type as either certificate manager instance or secrets manager instance"`
	Status          string `json:"status" description:"Used to show the status indicating if the instance is registered to the cluster or not"`
}

// Instances struct for a secret array
type Instances []Instance

// InstanceRegisterConfig the instance register request
type InstanceRegisterConfig struct {
	Cluster       string `json:"cluster" description:"id of cluster" binding:"required"`
	CRN           string `json:"crn" description:"crn of the instance"`
	IsDefault     bool   `json:"isDefault" description:"true or false. Used to show which instance subdomains certificates are uploaded into"`
	SecretGroupID string `json:"secretGroupID" description:"ID of the secret group where secrets will be stored"`
}

// InstanceDeleteConfig the instance delete request
type InstanceDeleteConfig struct {
	Cluster string `json:"cluster" description:"id of cluster" binding:"required"`
	Name    string `json:"name" description:"name of instance" binding:"required"`
}

// InstanceUpdateConfig instance update request
type InstanceUpdateConfig struct {
	Cluster       string `json:"cluster" description:"id of cluster" binding:"required"`
	Name          string `json:"name" description:"name of instance" binding:"required"`
	IsDefault     bool   `json:"isDefault" description:"true or false. Used to show which instance subdomains certificates are uploaded into"`
	SecretGroupID string `json:"secretGroupID" description:"ID of the secret group where secrets will be stored"`
}

type ingress struct {
	client *client.Client
}

// Ingress interface
type Ingress interface {
	CreateIngressSecret(req SecretCreateConfig) (response Secret, err error)
	UpdateIngressSecret(req SecretUpdateConfig) (response Secret, err error)
	AddIngressSecretField(req SecretUpdateConfig) (response Secret, err error)
	RemoveIngressSecretField(req SecretUpdateConfig) (response Secret, err error)
	DeleteIngressSecret(req SecretDeleteConfig) (err error)
	GetIngressSecretList(clusterNameOrID string, showDeleted bool) (response Secrets, err error)
	GetIngressSecret(clusterNameOrID, secretName, secretNamespace string) (response Secret, err error)
	RegisterIngressInstance(req InstanceRegisterConfig) (response Instance, err error)
	UpdateIngressInstance(req InstanceUpdateConfig) (err error)
	DeleteIngressInstance(req InstanceDeleteConfig) (err error)
	GetIngressInstance(clusterNameOrID, instanceName string) (response Instance, err error)
	GetIngressInstanceList(clusterNameOrID string, showDeleted bool) (response Instances, err error)
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

// AddIngressSecretField adds secret fields to an existing secret
func (r *ingress) AddIngressSecretField(req SecretUpdateConfig) (response Secret, err error) {
	_, err = r.client.Post("/ingress/v2/secret/addField", req, &response)
	return
}

// RemoveIngressSecretField removes secret fields from an existing secret
func (r *ingress) RemoveIngressSecretField(req SecretUpdateConfig) (response Secret, err error) {
	_, err = r.client.Post("/ingress/v2/secret/removeField", req, &response)
	return
}

// DeleteIngressSecret deletes the ingress secret from the cluster
func (r *ingress) DeleteIngressSecret(req SecretDeleteConfig) (err error) {
	_, err = r.client.Post("/ingress/v2/secret/deleteSecret", req, nil)
	return
}

func (r *ingress) RegisterIngressInstance(req InstanceRegisterConfig) (response Instance, err error) {
	_, err = r.client.Post("/ingress/v2/secret/registerInstance", req, &response)
	return
}

func (r *ingress) UpdateIngressInstance(req InstanceUpdateConfig) (err error) {
	_, err = r.client.Post("/ingress/v2/secret/updateInstance", req, nil)
	return
}

func (r *ingress) DeleteIngressInstance(req InstanceDeleteConfig) (err error) {
	_, err = r.client.Post("/ingress/v2/secret/unregisterInstance", req, nil)
	return
}

func (r *ingress) GetIngressInstance(clusterNameOrID, instanceName string) (response Instance, err error) {
	_, err = r.client.Get(fmt.Sprintf("/ingress/v2/secret/getInstance?cluster=%s&name=%s", clusterNameOrID, instanceName), &response)
	return
}

func (r *ingress) GetIngressInstanceList(clusterNameOrID string, showDeleted bool) (response Instances, err error) {
	deleted := strconv.FormatBool(showDeleted)
	_, err = r.client.Get(fmt.Sprintf("/ingress/v2/secret/getInstances?cluster=%s&showDeleted=%s", clusterNameOrID, deleted), &response)
	return
}
