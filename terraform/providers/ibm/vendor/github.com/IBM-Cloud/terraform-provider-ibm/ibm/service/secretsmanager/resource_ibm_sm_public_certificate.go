// Copyright IBM Corp. 2022 All Rights Reserved.
// Licensed under the Mozilla Public License v2.0

package secretsmanager

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/IBM-Cloud/bluemix-go/bmxerror"
	"github.com/akamai/AkamaiOPEN-edgegrid-golang/client-v1"
	"github.com/akamai/AkamaiOPEN-edgegrid-golang/edgegrid"
	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v5/pkg/dns"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"io/ioutil"
	"log"
	"reflect"
	"sort"
	"strings"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/conns"
	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/flex"
	"github.com/IBM/go-sdk-core/v5/core"
	"github.com/IBM/secrets-manager-go-sdk/v2/secretsmanagerv2"
)

func ResourceIbmSmPublicCertificate() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceIbmSmPublicCertificateCreate,
		ReadContext:   resourceIbmSmPublicCertificateRead,
		UpdateContext: resourceIbmSmPublicCertificateUpdate,
		DeleteContext: resourceIbmSmPublicCertificateDelete,
		Importer:      &schema.ResourceImporter{},

		Schema: map[string]*schema.Schema{
			"secret_type": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The secret type. Supported types are arbitrary, certificates (imported, public, and private), IAM credentials, key-value, and user credentials.",
			},
			"name": &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
				Description: "A human-readable name to assign to your secret.To protect your privacy, do not use personal data, such as your name or location, as a name for your secret.",
			},
			"description": &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				Description: "An extended description of your secret.To protect your privacy, do not use personal data, such as your name or location, as a description for your secret group.",
			},
			"secret_group_id": &schema.Schema{
				Type:        schema.TypeString,
				ForceNew:    true,
				Optional:    true,
				Computed:    true,
				Description: "A v4 UUID identifier, or `default` secret group.",
			},
			"labels": &schema.Schema{
				Type:        schema.TypeList,
				Optional:    true,
				Computed:    true,
				Description: "Labels that you can use to search for secrets in your instance.Up to 30 labels can be created.",
				Elem:        &schema.Schema{Type: schema.TypeString},
			},
			"common_name": &schema.Schema{
				Type:        schema.TypeString,
				ForceNew:    true,
				Required:    true,
				Description: "The Common Name (AKA CN) represents the server name that is protected by the SSL certificate.",
			},
			"alt_names": &schema.Schema{
				Type:             schema.TypeList,
				ForceNew:         true,
				Optional:         true,
				Computed:         true,
				Description:      "With the Subject Alternative Name field, you can specify additional host names to be protected by a single SSL certificate.",
				Elem:             &schema.Schema{Type: schema.TypeString},
				DiffSuppressFunc: altNamesDiffSuppress,
			},
			"key_algorithm": &schema.Schema{
				Type:        schema.TypeString,
				ForceNew:    true,
				Optional:    true,
				Default:     "RSA2048",
				Description: "The identifier for the cryptographic algorithm to be used to generate the public key that is associated with the certificate.The algorithm that you select determines the encryption algorithm (`RSA` or `ECDSA`) and key size to be used to generate keys and sign certificates. For longer living certificates, it is recommended to use longer keys to provide more encryption protection. Allowed values:  RSA2048, RSA4096, EC256, EC384.",
			},
			"ca": &schema.Schema{
				Type:        schema.TypeString,
				ForceNew:    true,
				Required:    true,
				Description: "The name of the certificate authority configuration.",
			},
			"dns": &schema.Schema{
				Type:        schema.TypeString,
				ForceNew:    true,
				Required:    true,
				Description: "The name of the DNS provider configuration.",
			},
			"bundle_certs": &schema.Schema{
				Type:        schema.TypeBool,
				ForceNew:    true,
				Optional:    true,
				Default:     true,
				Description: "Determines whether your issued certificate is bundled with intermediate certificates. Set to `false` for the certificate file to contain only the issued certificate.",
			},
			"rotation": &schema.Schema{
				Type:        schema.TypeList,
				MaxItems:    1,
				Optional:    true,
				Computed:    true,
				Description: "Determines whether Secrets Manager rotates your secrets automatically.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"auto_rotate": &schema.Schema{
							Type:        schema.TypeBool,
							Optional:    true,
							Computed:    true,
							Description: "Determines whether Secrets Manager rotates your secret automatically.Default is `false`. If `auto_rotate` is set to `true` the service rotates your certificate 31 days before it expires.",
						},
						"rotate_keys": &schema.Schema{
							Type:        schema.TypeBool,
							Optional:    true,
							Computed:    true,
							Description: "Determines whether Secrets Manager rotates the private key for your public certificate automatically.Default is `false`. If it is set to `true`, the service generates and stores a new private key for your rotated certificate.",
						},
					},
				},
			},
			"custom_metadata": &schema.Schema{
				Type:        schema.TypeMap,
				Optional:    true,
				Computed:    true,
				Description: "The secret metadata that a user can customize.",
				Elem:        &schema.Schema{Type: schema.TypeString},
			},
			"version_custom_metadata": &schema.Schema{
				Type:        schema.TypeMap,
				Optional:    true,
				Description: "The secret version metadata that a user can customize.",
				Elem:        &schema.Schema{Type: schema.TypeString},
			},
			"akamai": &schema.Schema{
				Type:        schema.TypeList,
				MaxItems:    1,
				Optional:    true,
				Description: "",
				ForceNew:    true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"edgerc": &schema.Schema{
							Type:        schema.TypeList,
							Optional:    true,
							Description: "Akamai credentials",
							MaxItems:    1,
							ForceNew:    true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"path_to_edgerc": &schema.Schema{
										Type:        schema.TypeString,
										Optional:    true,
										Description: "Path to Akamai's configuration file.",
										ForceNew:    true,
									},
									"config_section": &schema.Schema{
										Description: "The section of the edgerc file to use for configuration.",
										Optional:    true,
										Type:        schema.TypeString,
										Default:     "default",
										ForceNew:    true,
									},
								},
							},
						},
						"config": &schema.Schema{
							Type:        schema.TypeList,
							Optional:    true,
							Description: "Akamai credentials",
							MaxItems:    1,
							ForceNew:    true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"host": &schema.Schema{
										Type:      schema.TypeString,
										Optional:  true,
										ForceNew:  true,
										Sensitive: true,
									},
									"client_secret": &schema.Schema{
										Type:      schema.TypeString,
										Optional:  true,
										ForceNew:  true,
										Sensitive: true,
									},
									"access_token": &schema.Schema{
										Type:      schema.TypeString,
										Optional:  true,
										ForceNew:  true,
										Sensitive: true,
									},
									"client_token": &schema.Schema{
										Type:      schema.TypeString,
										Optional:  true,
										ForceNew:  true,
										Sensitive: true,
									},
								},
							},
						},
					},
				},
			},
			"created_by": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The unique identifier that is associated with the entity that created the secret.",
			},
			"created_at": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The date when a resource was created. The date format follows RFC 3339.",
			},
			"crn": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "A CRN that uniquely identifies an IBM Cloud resource.",
			},
			"downloaded": &schema.Schema{
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "Indicates whether the secret data that is associated with a secret version was retrieved in a call to the service API.",
			},
			"secret_id": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "A v4 UUID identifier.",
			},
			"locks_total": &schema.Schema{
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "The number of locks of the secret.",
			},
			"state": &schema.Schema{
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "The secret state that is based on NIST SP 800-57. States are integers and correspond to the `Pre-activation = 0`, `Active = 1`,  `Suspended = 2`, `Deactivated = 3`, and `Destroyed = 5` values.",
			},
			"state_description": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "A text representation of the secret state.",
			},
			"updated_at": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The date when a resource was recently modified. The date format follows RFC 3339.",
			},
			"versions_total": &schema.Schema{
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "The number of versions of the secret.",
			},
			"signing_algorithm": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The identifier for the cryptographic algorithm that was used by the issuing certificate authority to sign a certificate.",
			},
			"expiration_date": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The date a secret is expired. The date format follows RFC 3339.",
			},
			"issuance_info": &schema.Schema{
				Type:        schema.TypeList,
				Computed:    true,
				Description: "Issuance information that is associated with your certificate.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"auto_rotated": &schema.Schema{
							Type:        schema.TypeBool,
							Computed:    true,
							Description: "Indicates whether the issued certificate is configured with an automatic rotation policy.",
						},
						"challenges": &schema.Schema{
							Type:        schema.TypeList,
							Computed:    true,
							Description: "The set of challenges. It is returned only when ordering public certificates by using manual DNS configuration.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"domain": &schema.Schema{
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The challenge domain.",
									},
									"expiration": &schema.Schema{
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The challenge expiration date. The date format follows RFC 3339.",
									},
									"status": &schema.Schema{
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The challenge status.",
									},
									"txt_record_name": &schema.Schema{
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The TXT record name.",
									},
									"txt_record_value": &schema.Schema{
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The TXT record value.",
									},
								},
							},
						},
						"dns_challenge_validation_time": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The date that a user requests to validate DNS challenges for certificates that are ordered with a manual DNS provider. The date format follows RFC 3339.",
						},
						"error_code": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "A code that identifies an issuance error.This field, along with `error_message`, is returned when Secrets Manager successfully processes your request, but the certificate authority is unable to issue a certificate.",
						},
						"error_message": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "A human-readable message that provides details about the issuance error.",
						},
						"ordered_on": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The date when the certificate is ordered. The date format follows RFC 3339.",
						},
						"state": &schema.Schema{
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "The secret state that is based on NIST SP 800-57. States are integers and correspond to the `Pre-activation = 0`, `Active = 1`,  `Suspended = 2`, `Deactivated = 3`, and `Destroyed = 5` values.",
						},
						"state_description": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "A text representation of the secret state.",
						},
					},
				},
			},
			"issuer": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The distinguished name that identifies the entity that signed and issued the certificate.",
			},
			"serial_number": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The unique serial number that was assigned to a certificate by the issuing certificate authority.",
			},
			"validity": &schema.Schema{
				Type:        schema.TypeList,
				Computed:    true,
				Description: "The date and time that the certificate validity period begins and ends.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"not_before": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The date-time format follows RFC 3339.",
						},
						"not_after": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The date-time format follows RFC 3339.",
						},
					},
				},
			},
			"certificate": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Sensitive:   true,
				Description: "The PEM-encoded contents of your certificate.",
			},
			"intermediate": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Sensitive:   true,
				Description: "(Optional) The PEM-encoded intermediate certificate to associate with the root certificate.",
			},
			"private_key": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Sensitive:   true,
				Description: "(Optional) The PEM-encoded private key to associate with the certificate.",
			},
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(37 * time.Minute),
		},
	}
}

func resourceIbmSmPublicCertificateCreate(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	secretsManagerClient, err := meta.(conns.ClientSession).SecretsManagerV2()
	if err != nil {
		tfErr := flex.TerraformErrorf(err, "", PublicCertSecretResourceName, "create")
		return tfErr.GetDiag()
	}

	region := getRegion(secretsManagerClient, d)
	instanceId := d.Get("instance_id").(string)
	secretsManagerClient = getClientWithInstanceEndpoint(secretsManagerClient, instanceId, region, getEndpointType(secretsManagerClient, d))

	createSecretOptions := &secretsmanagerv2.CreateSecretOptions{}

	secretPrototypeModel, err := resourceIbmSmPublicCertificateMapToSecretPrototype(d)
	if err != nil {
		tfErr := flex.TerraformErrorf(err, "", PublicCertSecretResourceName, "create")
		return tfErr.GetDiag()
	}
	createSecretOptions.SetSecretPrototype(secretPrototypeModel)

	secretIntf, response, err := secretsManagerClient.CreateSecretWithContext(context, createSecretOptions)
	if err != nil {
		log.Printf("[DEBUG] CreateSecretWithContext failed %s\n%s", err, response)
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("CreateSecretWithContext failed: %s\n%s", err.Error(), response), PublicCertSecretResourceName, "create")
		return tfErr.GetDiag()
	}

	secret := secretIntf.(*secretsmanagerv2.PublicCertificate)
	d.SetId(fmt.Sprintf("%s/%s/%s", region, instanceId, *secret.ID))
	d.Set("secret_id", *secret.ID)

	if *secret.Dns == "manual" || *secret.Dns == "akamai" {
		_, err = waitForIbmSmPublicCertificateCreate(secretsManagerClient, d, "", "pre_activation")
	} else {
		_, err = waitForIbmSmPublicCertificateCreate(secretsManagerClient, d, "pre_activation", "active")
	}

	if err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("Error waiting for resource IbmSmPublicCertificate (%s) to be created: %s", d.Id(), err.Error()), PublicCertSecretResourceName, "create")
		return tfErr.GetDiag()
	}

	return resourceIbmSmPublicCertificateRead(context, d, meta)
}

func waitForIbmSmPublicCertificateCreate(secretsManagerClient *secretsmanagerv2.SecretsManagerV2, d *schema.ResourceData, pendingStatus string, targetStatus string) (interface{}, error) {
	getSecretOptions := &secretsmanagerv2.GetSecretOptions{}

	id := strings.Split(d.Id(), "/")
	secretId := id[2]

	getSecretOptions.SetID(secretId)
	stateConf := &resource.StateChangeConf{
		Pending: []string{pendingStatus},
		Target:  []string{targetStatus},
		Refresh: func() (interface{}, string, error) {
			stateObjIntf, response, err := secretsManagerClient.GetSecret(getSecretOptions)
			stateObj := stateObjIntf.(*secretsmanagerv2.PublicCertificate)
			if err != nil {
				if apiErr, ok := err.(bmxerror.RequestFailure); ok && apiErr.StatusCode() == 404 {
					return nil, "", fmt.Errorf("The instance %s does not exist anymore: %s\n%s", "getSecretOptions", err, response)
				}
				return nil, "", err
			}
			failStates := map[string]bool{"destroyed": true}
			if failStates[*stateObj.StateDescription] {
				return stateObj, *stateObj.StateDescription, fmt.Errorf("The instance %s failed: %s\n%s", "getSecretOptions", err, response)
			}
			return stateObj, *stateObj.StateDescription, nil
		},
		Timeout:    d.Timeout(schema.TimeoutCreate),
		Delay:      0 * time.Second,
		MinTimeout: 5 * time.Second,
	}

	return stateConf.WaitForState()
}

func resourceIbmSmPublicCertificateRead(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	secretsManagerClient, err := meta.(conns.ClientSession).SecretsManagerV2()
	if err != nil {
		tfErr := flex.TerraformErrorf(err, "", PublicCertSecretResourceName, "read")
		return tfErr.GetDiag()
	}

	id := strings.Split(d.Id(), "/")
	if len(id) != 3 {
		tfErr := flex.TerraformErrorf(nil, "Wrong format of resource ID. To import a secret use the format `<region>/<instance_id>/<secret_id>`", PublicCertSecretResourceName, "read")
		return tfErr.GetDiag()
	}
	region := id[0]
	instanceId := id[1]
	secretId := id[2]
	secretsManagerClient = getClientWithInstanceEndpoint(secretsManagerClient, instanceId, region, getEndpointType(secretsManagerClient, d))

	getSecretOptions := &secretsmanagerv2.GetSecretOptions{}

	getSecretOptions.SetID(secretId)

	secretIntf, response, err := secretsManagerClient.GetSecretWithContext(context, getSecretOptions)
	if err != nil {
		if response != nil && response.StatusCode == 404 {
			d.SetId("")
			return nil
		}
		log.Printf("[DEBUG] GetSecretWithContext failed %s\n%s", err, response)
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("GetSecretWithContext failed %s\n%s", err, response), PublicCertSecretResourceName, "read")
		return tfErr.GetDiag()
	}

	secret := secretIntf.(*secretsmanagerv2.PublicCertificate)

	if err = d.Set("secret_id", secretId); err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("Error setting secret_id"), PublicCertSecretResourceName, "read")
		return tfErr.GetDiag()
	}
	if err = d.Set("instance_id", instanceId); err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("Error setting instance_id"), PublicCertSecretResourceName, "read")
		return tfErr.GetDiag()
	}
	if err = d.Set("region", region); err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("Error setting region"), PublicCertSecretResourceName, "read")
		return tfErr.GetDiag()
	}
	if err = d.Set("created_by", secret.CreatedBy); err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("Error setting created_by"), PublicCertSecretResourceName, "read")
		return tfErr.GetDiag()
	}
	if err = d.Set("created_at", DateTimeToRFC3339(secret.CreatedAt)); err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("Error setting created_at"), PublicCertSecretResourceName, "read")
		return tfErr.GetDiag()
	}
	if err = d.Set("crn", secret.Crn); err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("Error setting crn"), PublicCertSecretResourceName, "read")
		return tfErr.GetDiag()
	}
	if err = d.Set("downloaded", secret.Downloaded); err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("Error setting downloaded"), PublicCertSecretResourceName, "read")
		return tfErr.GetDiag()
	}
	if err = d.Set("locks_total", flex.IntValue(secret.LocksTotal)); err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("Error setting locks_total"), PublicCertSecretResourceName, "read")
		return tfErr.GetDiag()
	}
	if err = d.Set("name", secret.Name); err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("Error setting name"), PublicCertSecretResourceName, "read")
		return tfErr.GetDiag()
	}
	if err = d.Set("secret_group_id", secret.SecretGroupID); err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("Error setting secret_group_id"), PublicCertSecretResourceName, "read")
		return tfErr.GetDiag()
	}
	if err = d.Set("secret_type", secret.SecretType); err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("Error setting secret_type"), PublicCertSecretResourceName, "read")
		return tfErr.GetDiag()
	}
	if err = d.Set("state", flex.IntValue(secret.State)); err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("Error setting state"), PublicCertSecretResourceName, "read")
		return tfErr.GetDiag()
	}
	if err = d.Set("state_description", secret.StateDescription); err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("Error setting state_description"), PublicCertSecretResourceName, "read")
		return tfErr.GetDiag()
	}
	if err = d.Set("updated_at", DateTimeToRFC3339(secret.UpdatedAt)); err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("Error setting updated_at"), PublicCertSecretResourceName, "read")
		return tfErr.GetDiag()
	}
	if err = d.Set("versions_total", flex.IntValue(secret.VersionsTotal)); err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("Error setting versions_total"), PublicCertSecretResourceName, "read")
		return tfErr.GetDiag()
	}
	if err = d.Set("common_name", secret.CommonName); err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("Error setting common_name"), PublicCertSecretResourceName, "read")
		return tfErr.GetDiag()
	}
	if secret.IssuanceInfo != nil {
		issuanceInfoMap, err := resourceIbmSmPublicCertificateCertificateIssuanceInfoToMap(secret.IssuanceInfo, d)
		if err != nil {
			tfErr := flex.TerraformErrorf(err, "", PublicCertSecretResourceName, "read")
			return tfErr.GetDiag()
		}
		if err = d.Set("issuance_info", []map[string]interface{}{issuanceInfoMap}); err != nil {
			tfErr := flex.TerraformErrorf(err, fmt.Sprintf("Error setting issuance_info"), PublicCertSecretResourceName, "read")
			return tfErr.GetDiag()
		}
	}

	if err = d.Set("key_algorithm", secret.KeyAlgorithm); err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("Error setting key_algorithm"), PublicCertSecretResourceName, "read")
		return tfErr.GetDiag()
	}
	if err = d.Set("ca", secret.Ca); err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("Error setting ca"), PublicCertSecretResourceName, "read")
		return tfErr.GetDiag()
	}
	if d.Get("dns").(string) != "akamai" {
		if err = d.Set("dns", secret.Dns); err != nil {
			tfErr := flex.TerraformErrorf(err, fmt.Sprintf("Error setting dns"), PublicCertSecretResourceName, "read")
			return tfErr.GetDiag()
		}
	}
	if err = d.Set("bundle_certs", secret.BundleCerts); err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("Error setting bundle_certs"), PublicCertSecretResourceName, "read")
		return tfErr.GetDiag()
	}
	rotationMap, err := resourceIbmSmPublicCertificateRotationPolicyToMap(secret.Rotation)
	if err != nil {
		tfErr := flex.TerraformErrorf(err, "", PublicCertSecretResourceName, "read")
		return tfErr.GetDiag()
	}
	if err = d.Set("rotation", []map[string]interface{}{rotationMap}); err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("Error setting rotation"), PublicCertSecretResourceName, "read")
		return tfErr.GetDiag()
	}
	if secret.CustomMetadata != nil {
		d.Set("custom_metadata", secret.CustomMetadata)
	}
	if err = d.Set("description", secret.Description); err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("Error setting description"), PublicCertSecretResourceName, "read")
		return tfErr.GetDiag()
	}
	if secret.Labels != nil {
		if err = d.Set("labels", secret.Labels); err != nil {
			tfErr := flex.TerraformErrorf(err, fmt.Sprintf("Error setting labels"), PublicCertSecretResourceName, "read")
			return tfErr.GetDiag()
		}
	}
	if err = d.Set("signing_algorithm", secret.SigningAlgorithm); err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("Error setting signing_algorithm"), PublicCertSecretResourceName, "read")
		return tfErr.GetDiag()
	}
	if secret.AltNames != nil {
		if err = d.Set("alt_names", secret.AltNames); err != nil {
			tfErr := flex.TerraformErrorf(err, fmt.Sprintf("Error setting alt_names"), PublicCertSecretResourceName, "read")
			return tfErr.GetDiag()
		}
	}
	if err = d.Set("expiration_date", DateTimeToRFC3339(secret.ExpirationDate)); err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("Error setting expiration_date"), PublicCertSecretResourceName, "read")
		return tfErr.GetDiag()
	}
	if err = d.Set("issuer", secret.Issuer); err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("Error setting issuer"), PublicCertSecretResourceName, "read")
		return tfErr.GetDiag()
	}
	if err = d.Set("serial_number", secret.SerialNumber); err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("Error setting serial_number"), PublicCertSecretResourceName, "read")
		return tfErr.GetDiag()
	}
	validityMap, err := resourceIbmSmPublicCertificateCertificateValidityToMap(secret.Validity)
	if err != nil {
		tfErr := flex.TerraformErrorf(err, "", PublicCertSecretResourceName, "read")
		return tfErr.GetDiag()
	}
	if err = d.Set("validity", []map[string]interface{}{validityMap}); err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("Error setting validity"), PublicCertSecretResourceName, "read")
		return tfErr.GetDiag()
	}
	if err = d.Set("certificate", secret.Certificate); err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("Error setting certificate"), PublicCertSecretResourceName, "read")
		return tfErr.GetDiag()
	}
	if err = d.Set("intermediate", secret.Intermediate); err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("Error setting intermediate"), PublicCertSecretResourceName, "read")
		return tfErr.GetDiag()
	}
	if err = d.Set("private_key", secret.PrivateKey); err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("Error setting private_key"), PublicCertSecretResourceName, "read")
		return tfErr.GetDiag()
	}

	if *secret.StateDescription == "active" {
		// Call get version metadata API to get the current version_custom_metadata
		getVersionMetdataOptions := &secretsmanagerv2.GetSecretVersionMetadataOptions{}
		getVersionMetdataOptions.SetSecretID(secretId)
		getVersionMetdataOptions.SetID("current")

		versionMetadataIntf, response, err := secretsManagerClient.GetSecretVersionMetadataWithContext(context, getVersionMetdataOptions)
		if err != nil {
			log.Printf("[DEBUG] GetSecretVersionMetadataWithContext failed %s\n%s", err, response)
			tfErr := flex.TerraformErrorf(err, fmt.Sprintf("GetSecretVersionMetadataWithContext failed %s\n%s", err, response), PublicCertSecretResourceName, "read")
			return tfErr.GetDiag()
		}

		versionMetadata := versionMetadataIntf.(*secretsmanagerv2.PublicCertificateVersionMetadata)
		if versionMetadata.VersionCustomMetadata != nil {
			if err = d.Set("version_custom_metadata", versionMetadata.VersionCustomMetadata); err != nil {
				tfErr := flex.TerraformErrorf(err, fmt.Sprintf("Error setting version_custom_metadata"), PublicCertSecretResourceName, "read")
				return tfErr.GetDiag()
			}
		}
	}

	if d.Get("dns").(string) == "akamai" && d.Get("state_description").(string) == "pre_activation" {
		err := setChallengesWithAkamaiAndValidateManualDns(context, d, meta, secret, secretsManagerClient)
		if err != nil {
			return err
		}
		return resourceIbmSmPublicCertificateRead(context, d, meta)
	}
	return nil
}

func resourceIbmSmPublicCertificateUpdate(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	secretsManagerClient, err := meta.(conns.ClientSession).SecretsManagerV2()
	if err != nil {
		tfErr := flex.TerraformErrorf(err, "", PublicCertSecretResourceName, "update")
		return tfErr.GetDiag()
	}

	id := strings.Split(d.Id(), "/")
	region := id[0]
	instanceId := id[1]
	secretId := id[2]
	secretsManagerClient = getClientWithInstanceEndpoint(secretsManagerClient, instanceId, region, getEndpointType(secretsManagerClient, d))

	updateSecretMetadataOptions := &secretsmanagerv2.UpdateSecretMetadataOptions{}

	updateSecretMetadataOptions.SetID(secretId)

	hasChange := false

	patchVals := &secretsmanagerv2.SecretMetadataPatch{}

	if d.HasChange("name") {
		patchVals.Name = core.StringPtr(d.Get("name").(string))
		hasChange = true
	}
	if d.HasChange("description") {
		patchVals.Description = core.StringPtr(d.Get("description").(string))
		hasChange = true
	}
	if d.HasChange("labels") {
		labels := d.Get("labels").([]interface{})
		labelsParsed := make([]string, len(labels))
		for i, v := range labels {
			labelsParsed[i] = fmt.Sprint(v)
		}
		patchVals.Labels = labelsParsed
		hasChange = true
	}
	if d.HasChange("custom_metadata") {
		patchVals.CustomMetadata = d.Get("custom_metadata").(map[string]interface{})
		hasChange = true
	}
	if d.HasChange("rotation") {
		RotationModel, err := resourceIbmSmPublicCertificateMapToRotationPolicy(d.Get("rotation").([]interface{})[0].(map[string]interface{}))
		if err != nil {
			log.Printf("[DEBUG] UpdateSecretMetadataWithContext failed: Reading Rotation parameter failed: %s", err)
			tfErr := flex.TerraformErrorf(err, fmt.Sprintf("UpdateSecretMetadataWithContext failed: Reading Rotation parameter failed: %s", err), PublicCertSecretResourceName, "update")
			return tfErr.GetDiag()
		}
		patchVals.Rotation = RotationModel
		hasChange = true
	}

	if hasChange {
		updateSecretMetadataOptions.SecretMetadataPatch, _ = patchVals.AsPatch()
		_, response, err := secretsManagerClient.UpdateSecretMetadataWithContext(context, updateSecretMetadataOptions)
		if err != nil {
			log.Printf("[DEBUG] UpdateSecretMetadataWithContext failed %s\n%s", err, response)
			tfErr := flex.TerraformErrorf(err, fmt.Sprintf("UpdateSecretMetadataWithContext failed %s\n%s", err, response), PublicCertSecretResourceName, "update")
			return tfErr.GetDiag()
		}
	}

	if d.HasChange("version_custom_metadata") {
		// Apply change to version_custom_metadata in current version
		secretVersionMetadataPatchModel := new(secretsmanagerv2.SecretVersionMetadataPatch)
		secretVersionMetadataPatchModel.VersionCustomMetadata = d.Get("version_custom_metadata").(map[string]interface{})
		secretVersionMetadataPatchModelAsPatch, _ := secretVersionMetadataAsPatchFunction(secretVersionMetadataPatchModel)

		updateSecretVersionOptions := &secretsmanagerv2.UpdateSecretVersionMetadataOptions{}
		updateSecretVersionOptions.SetSecretID(secretId)
		updateSecretVersionOptions.SetID("current")
		updateSecretVersionOptions.SetSecretVersionMetadataPatch(secretVersionMetadataPatchModelAsPatch)
		_, response, err := secretsManagerClient.UpdateSecretVersionMetadataWithContext(context, updateSecretVersionOptions)
		if err != nil {
			if hasChange {
				// Call the read function to update the Terraform state with the change already applied to the metadata
				resourceIbmSmPublicCertificateRead(context, d, meta)
			}
			log.Printf("[DEBUG] UpdateSecretVersionMetadataWithContext failed %s\n%s", err, response)
			tfErr := flex.TerraformErrorf(err, fmt.Sprintf("UpdateSecretVersionMetadataWithContext failed %s\n%s", err, response), PublicCertSecretResourceName, "update")
			return tfErr.GetDiag()
		}
	}

	return resourceIbmSmPublicCertificateRead(context, d, meta)
}

func resourceIbmSmPublicCertificateDelete(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	secretsManagerClient, err := meta.(conns.ClientSession).SecretsManagerV2()
	if err != nil {
		tfErr := flex.TerraformErrorf(err, "", PublicCertSecretResourceName, "delete")
		return tfErr.GetDiag()
	}

	id := strings.Split(d.Id(), "/")
	region := id[0]
	instanceId := id[1]
	secretId := id[2]
	secretsManagerClient = getClientWithInstanceEndpoint(secretsManagerClient, instanceId, region, getEndpointType(secretsManagerClient, d))

	deleteSecretOptions := &secretsmanagerv2.DeleteSecretOptions{}

	deleteSecretOptions.SetID(secretId)

	response, err := secretsManagerClient.DeleteSecretWithContext(context, deleteSecretOptions)
	if err != nil {
		log.Printf("[DEBUG] DeleteSecretWithContext failed %s\n%s", err, response)
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("DeleteSecretWithContext failed %s\n%s", err, response), PublicCertSecretResourceName, "delete")
		return tfErr.GetDiag()
	}

	d.SetId("")

	return nil
}

func resourceIbmSmPublicCertificateMapToSecretPrototype(d *schema.ResourceData) (secretsmanagerv2.SecretPrototypeIntf, error) {
	model := &secretsmanagerv2.PublicCertificatePrototype{}
	model.SecretType = core.StringPtr("public_cert")

	if _, ok := d.GetOk("name"); ok {
		model.Name = core.StringPtr(d.Get("name").(string))
	}
	if _, ok := d.GetOk("description"); ok {
		model.Description = core.StringPtr(d.Get("description").(string))
	}
	if _, ok := d.GetOk("secret_group_id"); ok {
		model.SecretGroupID = core.StringPtr(d.Get("secret_group_id").(string))
	}
	if _, ok := d.GetOk("labels"); ok {
		labels := d.Get("labels").([]interface{})
		labelsParsed := make([]string, len(labels))
		for i, v := range labels {
			labelsParsed[i] = fmt.Sprint(v)
		}
		model.Labels = labelsParsed
	}
	if _, ok := d.GetOk("common_name"); ok {
		model.CommonName = core.StringPtr(d.Get("common_name").(string))
	}
	if _, ok := d.GetOk("alt_names"); ok {
		altNames := d.Get("alt_names").([]interface{})
		altNamesParsed := make([]string, len(altNames))
		for i, v := range altNames {
			altNamesParsed[i] = fmt.Sprint(v)
		}
		model.AltNames = altNamesParsed
	}
	if _, ok := d.GetOk("key_algorithm"); ok {
		model.KeyAlgorithm = core.StringPtr(d.Get("key_algorithm").(string))
	}
	if _, ok := d.GetOk("ca"); ok {
		model.Ca = core.StringPtr(d.Get("ca").(string))
	}
	if _, ok := d.GetOk("dns"); ok {
		if d.Get("dns").(string) == "akamai" {
			model.Dns = core.StringPtr("manual")
		} else {
			model.Dns = core.StringPtr(d.Get("dns").(string))
		}
	}
	bundleCerts, ok := d.GetOkExists("bundle_certs")
	if ok {
		model.BundleCerts = core.BoolPtr(bundleCerts.(bool))
	}
	if _, ok := d.GetOk("rotation"); ok {
		RotationModel, err := resourceIbmSmPublicCertificateMapToPublicCertificateRotationPolicy(d.Get("rotation").([]interface{})[0].(map[string]interface{}))
		if err != nil {
			return model, err
		}
		model.Rotation = RotationModel
	}
	if _, ok := d.GetOk("custom_metadata"); ok {
		model.CustomMetadata = d.Get("custom_metadata").(map[string]interface{})
	}
	if _, ok := d.GetOk("version_custom_metadata"); ok {
		model.VersionCustomMetadata = d.Get("version_custom_metadata").(map[string]interface{})
	}
	return model, nil
}

func resourceIbmSmPublicCertificateMapToRotationPolicy(modelMap map[string]interface{}) (secretsmanagerv2.RotationPolicyIntf, error) {
	model := &secretsmanagerv2.RotationPolicy{}
	if modelMap["auto_rotate"] != nil {
		model.AutoRotate = core.BoolPtr(modelMap["auto_rotate"].(bool))
	}
	if modelMap["interval"] != nil {
		model.Interval = core.Int64Ptr(int64(modelMap["interval"].(int)))
	}
	if modelMap["unit"] != nil && modelMap["unit"].(string) != "" {
		model.Unit = core.StringPtr(modelMap["unit"].(string))
	}
	if modelMap["rotate_keys"] != nil {
		model.RotateKeys = core.BoolPtr(modelMap["rotate_keys"].(bool))
	}
	return model, nil
}

func resourceIbmSmPublicCertificateMapToCommonRotationPolicy(modelMap map[string]interface{}) (*secretsmanagerv2.CommonRotationPolicy, error) {
	model := &secretsmanagerv2.CommonRotationPolicy{}
	if modelMap["auto_rotate"] != nil {
		model.AutoRotate = core.BoolPtr(modelMap["auto_rotate"].(bool))
	}
	if modelMap["interval"] != nil {
		model.Interval = core.Int64Ptr(int64(modelMap["interval"].(int)))
	}
	if modelMap["unit"] != nil && modelMap["unit"].(string) != "" {
		model.Unit = core.StringPtr(modelMap["unit"].(string))
	}
	return model, nil
}

func resourceIbmSmPublicCertificateMapToPublicCertificateRotationPolicy(modelMap map[string]interface{}) (*secretsmanagerv2.PublicCertificateRotationPolicy, error) {
	model := &secretsmanagerv2.PublicCertificateRotationPolicy{}
	if modelMap["auto_rotate"] != nil {
		model.AutoRotate = core.BoolPtr(modelMap["auto_rotate"].(bool))
	}
	if modelMap["rotate_keys"] != nil {
		model.RotateKeys = core.BoolPtr(modelMap["rotate_keys"].(bool))
	}
	return model, nil
}

func resourceIbmSmPublicCertificateRotationPolicyToMap(modelIntf secretsmanagerv2.RotationPolicyIntf) (map[string]interface{}, error) {
	model := modelIntf.(*secretsmanagerv2.RotationPolicy)
	modelMap := make(map[string]interface{})

	if model.AutoRotate != nil {
		modelMap["auto_rotate"] = model.AutoRotate
	}
	if model.Interval != nil {
		modelMap["interval"] = flex.IntValue(model.Interval)
	}
	if model.Unit != nil {
		modelMap["unit"] = model.Unit
	}
	if model.RotateKeys != nil {
		modelMap["rotate_keys"] = model.RotateKeys
	}
	return modelMap, nil
}

func resourceIbmSmPublicCertificateCertificateIssuanceInfoToMap(model *secretsmanagerv2.CertificateIssuanceInfo, d *schema.ResourceData) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.AutoRotated != nil {
		modelMap["auto_rotated"] = model.AutoRotated
	}
	if model.Challenges != nil {
		challenges := []map[string]interface{}{}
		for _, challengesItem := range model.Challenges {
			challengesItemMap, err := resourceIbmSmPublicCertificateChallengeResourceToMap(&challengesItem)
			if err != nil {
				return modelMap, err
			}
			challenges = append(challenges, challengesItemMap)
		}
		modelMap["challenges"] = challenges
	} else {
		if d.Get("dns").(string) == "manual" {
			modelMap["challenges"] = d.Get("issuance_info").([]interface{})[0].(map[string]interface{})["challenges"]
		}
	}
	if model.DnsChallengeValidationTime != nil {
		modelMap["dns_challenge_validation_time"] = model.DnsChallengeValidationTime.String()
	}
	if model.ErrorCode != nil {
		modelMap["error_code"] = model.ErrorCode
	}
	if model.ErrorMessage != nil {
		modelMap["error_message"] = model.ErrorMessage
	}
	if model.OrderedOn != nil {
		modelMap["ordered_on"] = model.OrderedOn.String()
	}
	if model.State != nil {
		modelMap["state"] = flex.IntValue(model.State)
	}
	if model.StateDescription != nil {
		modelMap["state_description"] = model.StateDescription
	}
	return modelMap, nil
}

func resourceIbmSmPublicCertificateChallengeResourceToMap(model *secretsmanagerv2.ChallengeResource) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.Domain != nil {
		modelMap["domain"] = model.Domain
	}
	if model.Expiration != nil {
		modelMap["expiration"] = model.Expiration.String()
	}
	if model.Status != nil {
		modelMap["status"] = model.Status
	}
	if model.TxtRecordName != nil {
		modelMap["txt_record_name"] = model.TxtRecordName
	}
	if model.TxtRecordValue != nil {
		modelMap["txt_record_value"] = model.TxtRecordValue
	}
	return modelMap, nil
}

func resourceIbmSmPublicCertificateCertificateValidityToMap(model *secretsmanagerv2.CertificateValidity) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model == nil {
		modelMap["not_before"] = ""
		modelMap["not_after"] = ""
	} else {
		modelMap["not_before"] = model.NotBefore.String()
		modelMap["not_after"] = model.NotAfter.String()
	}
	return modelMap, nil
}

func altNamesDiffSuppress(key, oldValue, newValue string, d *schema.ResourceData) bool {
	lastDotIndex := strings.LastIndex(key, ".")
	if lastDotIndex != -1 {
		key = key[:lastDotIndex]
	}

	oldData, newData := d.GetChange(key)
	if oldData == nil || newData == nil {
		return false
	}

	oldAltNames, _ := ConcreteListToStringSlice(oldData.([]any), d.Get("common_name").(string))
	newAltNames, _ := ConcreteListToStringSlice(newData.([]any), d.Get("common_name").(string))

	sort.Strings(oldAltNames)
	sort.Strings(newAltNames)

	return reflect.DeepEqual(oldAltNames, newAltNames)
}

func ConcreteListToStringSlice(elements []interface{}, commonName string) ([]string, error) {
	// Create a new slice of strings with the same length as the input slice
	output := make([]string, len(elements))
	commonNameIndex := -1

	// Iterate over the input slice and cast each element to a string
	for i, elem := range elements {
		// Make sure the element is castable to string
		str, ok := elem.(string)
		if !ok {
			return nil, fmt.Errorf("cannot cast element %d to type []string", i)
		}
		if str == commonName {
			commonNameIndex = i
		}
		// Store the string in the output slice
		output[i] = str
	}

	if commonNameIndex != -1 {
		tempLst := append([]string{}, output[:commonNameIndex]...)
		output = append(tempLst, output[commonNameIndex+1:]...)
	}

	return output, nil
}

func setChallengesWithAkamaiAndValidateManualDns(context context.Context, d *schema.ResourceData, meta interface{}, secret *secretsmanagerv2.PublicCertificate, secretsManagerClient *secretsmanagerv2.SecretsManagerV2) diag.Diagnostics {
	config, err := configureAkamai(d)
	if err != nil {
		resourceIbmSmPublicCertificateDelete(context, d, meta)
		return err
	}

	successfullySetChallengeDomains := make(map[string]string)
	ttl := 120

	domains := []string{d.Get("common_name").(string)}
	if secret.AltNames != nil {
		domains = append(domains, secret.AltNames...)
	}

	for _, domainItem := range domains {
		domainForTxtRecordName := domainItem
		if domainItem[0] == '*' {
			domainForTxtRecordName = domainItem[2:]
		}

		txtRecordName := "_acme-challenge." + domainForTxtRecordName + "."

		txtRecordValuesChallenges, err := findAllTxtRecordValuesForDomain(domainItem, txtRecordName, secret, successfullySetChallengeDomains) // get txtRecordValues from our challenges
		if err != nil {
			resourceIbmSmPublicCertificateDelete(context, d, meta)
			return err
		}

		if len(txtRecordValuesChallenges) > 0 { // if we had not created already a dns record set for this txtRecordName (domain)
			zone, err := getZone(domainForTxtRecordName, domainItem, config)
			if err != nil {
				resourceIbmSmPublicCertificateDelete(context, d, meta)
				return err
			}

			txtRecordValuesAkamai, err := checkIfRecordExistsInAkamai(config, zone, txtRecordName) // get txtRecordValues from akamai
			if err != nil {
				resourceIbmSmPublicCertificateDelete(context, d, meta)
				return err
			}

			if len(txtRecordValuesAkamai) > 0 {
				txtRecordValuesToAdd := findTxtRecordValuesDifferences(txtRecordValuesAkamai, txtRecordValuesChallenges)
				if len(txtRecordValuesToAdd) > 0 {
					// there are already some txtRecordValues stored in akamai (len(txtRecordValuesAkamai > 0)) AND
					// there are *new* txtRecordValues in the challenges (len(txtRecordValuesToAdd) > 0) --> need to update
					// UPDATE
					txtRecordValuesAkamaiUpdated := append(txtRecordValuesAkamai, txtRecordValuesToAdd...)
					err := createOrUpdateAkamaiChallengeRecordSet(config, zone, txtRecordName, ttl, txtRecordValuesAkamaiUpdated, "PUT")
					//err := createOrUpdateAkamaiChallengeRecordSet(config, zone, txtRecordName, ttl, txtRecordValuesToAdd, "PUT")
					if err != nil {
						resourceIbmSmPublicCertificateDelete(context, d, meta)
						return err
					}
				}
			} else {
				// there is no txtRecordValues in akamai --> need to create
				// CREATE
				err := createOrUpdateAkamaiChallengeRecordSet(config, zone, txtRecordName, ttl, txtRecordValuesChallenges, "POST")
				if err != nil {
					resourceIbmSmPublicCertificateDelete(context, d, meta)
					return err
				}
			}
		}

	}

	for _, challengeItem := range secret.IssuanceInfo.Challenges {
		if _, exists := successfullySetChallengeDomains[*challengeItem.TxtRecordValue]; !exists {
			resourceIbmSmPublicCertificateDelete(context, d, meta)
			tfErr := flex.TerraformErrorf(nil, fmt.Sprintf("error: a dns record set in Akamai was not created for domain: %s", *challengeItem.Domain), PublicCertSecretResourceName, "read")
			return tfErr.GetDiag()
		}
	}

	return validateManualDns(context, d, secretsManagerClient)
}

func configureAkamai(d *schema.ResourceData) (edgegrid.Config, diag.Diagnostics) {
	var config edgegrid.Config
	var err error
	defaultErrMsg := "error configuring Akamai: One or more arguments are missing. Please verify that you provided either a path to your 'edgerc' file or all the config parameters ('host', 'client_secret', 'access_token' and 'client_token')"
	defaultTfErr := flex.TerraformErrorf(nil, defaultErrMsg, PublicCertSecretResourceName, "read")

	if len(d.Get("akamai").([]interface{})) == 0 || d.Get("akamai").([]interface{})[0] == nil {
		return config, defaultTfErr.GetDiag()
	}
	akamaiData := d.Get("akamai").([]interface{})[0].(map[string]interface{})

	if len(akamaiData["edgerc"].([]interface{})) > 0 {
		edgercData := akamaiData["edgerc"].([]interface{})[0].(map[string]interface{})
		edgerc := edgercData["path_to_edgerc"].(string)
		if edgerc == "" {
			return config, defaultTfErr.GetDiag()
		}
		configSection := edgercData["config_section"].(string)
		config, err = edgegrid.InitEdgeRc(edgerc, configSection)
		if err != nil {
			tfErr := flex.TerraformErrorf(nil, fmt.Sprintf("error initiating edgerc: %s", err), PublicCertSecretResourceName, "read")
			return config, tfErr.GetDiag()
		}
	} else if len(akamaiData["config"].([]interface{})) > 0 && akamaiData["config"].([]interface{})[0] != nil {
		akamaiDataConfig := akamaiData["config"].([]interface{})[0].(map[string]interface{})
		if akamaiDataConfig["host"] != "" && akamaiDataConfig["client_secret"] != "" && akamaiDataConfig["access_token"] != "" && akamaiDataConfig["client_token"] != "" {
			config.ClientSecret = akamaiDataConfig["client_secret"].(string)
			config.Host = akamaiDataConfig["host"].(string)
			config.AccessToken = akamaiDataConfig["access_token"].(string)
			config.ClientToken = akamaiDataConfig["client_token"].(string)
			if config.MaxBody == 0 {
				config.MaxBody = 131072
			}
		} else {
			return config, defaultTfErr.GetDiag()
		}
	} else {
		return config, defaultTfErr.GetDiag()
	}

	return config, nil

}

func checkIfRecordExistsInAkamai(config edgegrid.Config, zone string, txtRecordName string) ([]string, diag.Diagnostics) {
	req, err := client.NewRequest(config, "GET", fmt.Sprintf("/config-dns/v2/zones/%s/names/%s/types/TXT", zone, txtRecordName), nil)
	if err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("error creating akamai 'GET' request: %s", err), PublicCertSecretResourceName, "read")
		return nil, tfErr.GetDiag()
	}
	res, err := client.Do(config, req)
	if err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("error in performing akamai 'GET' request: %s", err), PublicCertSecretResourceName, "read")
		return nil, tfErr.GetDiag()
	}
	if res.StatusCode == 404 { // there is no record set, we need to create one
		return nil, nil
	} else if res.StatusCode == 200 {
		var recordData dns.RecordBody

		err := json.NewDecoder(res.Body).Decode(&recordData)
		if err != nil {
			tfErr := flex.TerraformErrorf(err, fmt.Sprintf("error in performing akamai 'GET' request: error in decoding JSON: %s", err), PublicCertSecretResourceName, "read")
			return nil, tfErr.GetDiag()
		}

		return recordData.Target, nil
	} else {
		body, err := ioutil.ReadAll(res.Body)
		if err != nil {
			fmt.Printf("Error reading response: %s\n", err.Error())
			tfErr := flex.TerraformErrorf(err, fmt.Sprintf("error in performing akamai 'GET' request: error reading error: %s", err), PublicCertSecretResourceName, "read")
			return nil, tfErr.GetDiag()
		}
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("error in performing akamai 'GET' request: %s", string(body)), PublicCertSecretResourceName, "read")
		return nil, tfErr.GetDiag()
	}
}

func createOrUpdateAkamaiChallengeRecordSet(config edgegrid.Config, zone string, txtRecordName string, ttl int, rdata []string, method string) diag.Diagnostics {
	type TXTRecordSet struct {
		Name  string   `json:"name"`
		Type  string   `json:"type"`
		TTL   int      `json:"ttl"`
		Rdata []string `json:"rdata"`
	}

	recordSetBody := TXTRecordSet{
		Name:  txtRecordName,
		Type:  "TXT",
		TTL:   ttl,
		Rdata: rdata,
	}

	jsonBody, err := json.Marshal(recordSetBody)
	if err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("error setting body for akamai request: %s", err), PublicCertSecretResourceName, "read")
		return tfErr.GetDiag()
	}
	req, err := client.NewRequest(config, method, fmt.Sprintf("/config-dns/v2/zones/%s/names/%s/types/TXT", zone, txtRecordName), bytes.NewReader(jsonBody))
	if err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("error creating akamai request: %s", err), PublicCertSecretResourceName, "read")
		return tfErr.GetDiag()
	}
	res, err := client.Do(config, req)
	if err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("error in akamai request: %s", err), PublicCertSecretResourceName, "read")
		return tfErr.GetDiag()
	}
	if res.StatusCode != 201 && res.StatusCode != 200 {
		body, err := ioutil.ReadAll(res.Body)
		if err != nil {
			fmt.Printf("Error reading response: %s\n", err.Error())
			return nil
		}
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("error from akamai in '%s' request: %s", method, string(body)), PublicCertSecretResourceName, "read")
		return tfErr.GetDiag()
	}
	return nil
}

func findAllTxtRecordValuesForDomain(domainItem string, txtRecordName string, secret *secretsmanagerv2.PublicCertificate, successfullySetChallengeDomains map[string]string) ([]string, diag.Diagnostics) {
	challenges := secret.IssuanceInfo.Challenges
	var txtRecordValues []string

	for _, challengesItem := range challenges {
		if *challengesItem.TxtRecordName == txtRecordName {
			if _, exists := successfullySetChallengeDomains[*challengesItem.TxtRecordValue]; !exists {
				txtRecordValues = append(txtRecordValues, *challengesItem.TxtRecordValue)
				successfullySetChallengeDomains[*challengesItem.TxtRecordValue] = *challengesItem.Domain
			} else { // if the txtRecordValue exists --> this txtRecordName (domain) has already been checked --> no need to continue
				return txtRecordValues, nil
			}
		}
	}

	if len(txtRecordValues) > 0 {
		return txtRecordValues, nil
	}
	tfErr := flex.TerraformErrorf(nil, fmt.Sprintf("failed to find a challenge for the domain: %s", domainItem), PublicCertSecretResourceName, "read")
	return nil, tfErr.GetDiag()
}

func findTxtRecordValuesDifferences(akamaiValues, challengesValues []string) []string {
	var differences []string

	if akamaiValues != nil {
		set := make(map[string]bool)
		for _, value := range akamaiValues {
			set[value] = true
		}

		for _, challengeValue := range challengesValues {
			stringChallengeValue := "\"" + challengeValue + "\""
			if !set[stringChallengeValue] {
				differences = append(differences, challengeValue)
			}
		}
	}

	return differences
}

func getZone(currentZone string, originalDomain string, config edgegrid.Config) (string, diag.Diagnostics) {
	req, err := client.NewRequest(config, "GET", fmt.Sprintf("/config-dns/v2/zones/%s", currentZone), nil)
	if err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("error creating akamai 'GET' zone request: %s", err), PublicCertSecretResourceName, "read")
		return "", tfErr.GetDiag()
	}
	res, err := client.Do(config, req)
	if err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("error in performing akamai 'GET' zone request: %s", err), PublicCertSecretResourceName, "read")
		return "", tfErr.GetDiag()
	}
	if res.StatusCode == 404 {
		zoneSplit := strings.Split(currentZone, ".")
		if len(zoneSplit) == 2 {
			tfErr := flex.TerraformErrorf(err, fmt.Sprintf("could not find a zone in Akamai for the domain: %s", originalDomain), PublicCertSecretResourceName, "read")
			return "", tfErr.GetDiag()
		}

		newZone := strings.Join(zoneSplit[1:], ".")
		return getZone(newZone, originalDomain, config)

	} else if res.StatusCode == 200 {
		return currentZone, nil
	} else {
		body, err := ioutil.ReadAll(res.Body)
		if err != nil {
			fmt.Printf("Error reading response: %s\n", err.Error())
			tfErr := flex.TerraformErrorf(err, fmt.Sprintf("error in performing akamai 'GET' zone request for zone: %s: error reading error: %s", currentZone, err), PublicCertSecretResourceName, "read")
			return "", tfErr.GetDiag()
		}
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("error in performing akamai 'GET' zone request for zone: %s:: %s", currentZone, string(body)), PublicCertSecretResourceName, "read")
		return "", tfErr.GetDiag()
	}
}
