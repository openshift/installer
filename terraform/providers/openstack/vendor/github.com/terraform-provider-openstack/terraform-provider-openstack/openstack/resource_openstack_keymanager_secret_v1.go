package openstack

import (
	"context"
	"log"
	"strings"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/customdiff"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"

	"github.com/gophercloud/gophercloud/openstack/keymanager/v1/acls"
	"github.com/gophercloud/gophercloud/openstack/keymanager/v1/secrets"
)

func resourceKeyManagerSecretV1() *schema.Resource {
	ret := &schema.Resource{
		CreateContext: resourceKeyManagerSecretV1Create,
		ReadContext:   resourceKeyManagerSecretV1Read,
		UpdateContext: resourceKeyManagerSecretV1Update,
		DeleteContext: resourceKeyManagerSecretV1Delete,

		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(30 * time.Minute),
			Update: schema.DefaultTimeout(30 * time.Minute),
			Delete: schema.DefaultTimeout(30 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},

			"name": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},

			"bit_length": {
				Type:     schema.TypeInt,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},

			"algorithm": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},

			"creator_id": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"mode": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},

			"secret_ref": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"secret_type": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				Computed: true,
				ValidateFunc: validation.StringInSlice([]string{
					"symmetric", "public", "private", "passphrase", "certificate", "opaque",
				}, false),
			},

			"status": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"payload": {
				Type:      schema.TypeString,
				Optional:  true,
				Sensitive: true,
				ForceNew:  true,
				Computed:  true,
				DiffSuppressFunc: func(k, o, n string, d *schema.ResourceData) bool {
					if strings.TrimSpace(o) == strings.TrimSpace(n) {
						return true
					}
					return false
				},
			},

			"payload_content_type": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				ValidateFunc: validation.StringInSlice([]string{
					"text/plain", "text/plain;charset=utf-8", "text/plain; charset=utf-8", "application/octet-stream", "application/pkcs8",
				}, true),
			},

			"payload_content_encoding": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				ValidateFunc: validation.StringInSlice([]string{
					"base64", "binary",
				}, false),
			},

			"metadata": {
				Type:     schema.TypeMap,
				Optional: true,
				ForceNew: false,
			},

			"acl": {
				Type:     schema.TypeList,
				Optional: true,
				Computed: true,
				MaxItems: 1,
			},

			"expiration": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				ValidateFunc: validation.IsRFC3339Time,
			},

			"created_at": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"updated_at": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"content_types": {
				Type:     schema.TypeMap,
				Computed: true,
			},

			"all_metadata": {
				Type:     schema.TypeMap,
				Computed: true,
			},
		},

		CustomizeDiff: customdiff.Sequence(
			// Clear the diff if the source payload is base64 encoded.
			func(ctx context.Context, diff *schema.ResourceDiff, v interface{}) error {
				return resourceSecretV1PayloadBase64CustomizeDiff(diff)
			},
		),
	}

	elem := &schema.Resource{
		Schema: make(map[string]*schema.Schema),
	}
	for _, aclOp := range getSupportedACLOperations() {
		elem.Schema[aclOp] = getACLSchema()
	}
	ret.Schema["acl"].Elem = elem

	return ret
}

func resourceKeyManagerSecretV1Create(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*Config)
	kmClient, err := config.KeyManagerV1Client(GetRegion(d, config))
	if err != nil {
		return diag.Errorf("Error creating OpenStack KeyManager client: %s", err)
	}

	var expiration *time.Time
	if v, err := time.Parse(time.RFC3339, d.Get("expiration").(string)); err == nil {
		expiration = &v
	}

	secretType := keyManagerSecretV1SecretType(d.Get("secret_type").(string))

	createOpts := secrets.CreateOpts{
		Name:       d.Get("name").(string),
		Algorithm:  d.Get("algorithm").(string),
		BitLength:  d.Get("bit_length").(int),
		Mode:       d.Get("mode").(string),
		Expiration: expiration,
		SecretType: secretType,
	}

	log.Printf("[DEBUG] Create Options for resource_keymanager_secret_v1: %#v", createOpts)

	var secret *secrets.Secret
	secret, err = secrets.Create(kmClient, createOpts).Extract()
	if err != nil {
		return diag.Errorf("Error creating openstack_keymanager_secret_v1: %s", err)
	}

	uuid := keyManagerSecretV1GetUUIDfromSecretRef(secret.SecretRef)

	stateConf := &resource.StateChangeConf{
		Pending:    []string{"PENDING"},
		Target:     []string{"ACTIVE"},
		Refresh:    keyManagerSecretV1WaitForSecretCreation(kmClient, uuid),
		Timeout:    d.Timeout(schema.TimeoutCreate),
		Delay:      0,
		MinTimeout: 2 * time.Second,
	}

	_, err = stateConf.WaitForStateContext(ctx)
	if err != nil {
		return diag.Errorf("Error waiting for openstack_keymanager_secret_v1: %s", err)
	}

	d.SetId(uuid)

	d.Partial(true)

	// set the acl first before uploading the payload
	if acl, ok := d.GetOk("acl"); ok {
		setOpts := expandKeyManagerV1ACLs(acl)
		_, err = acls.SetSecretACL(kmClient, uuid, setOpts).Extract()
		if err != nil {
			return diag.Errorf("Error settings ACLs for the openstack_keymanager_secret_v1: %s", err)
		}
	}

	// set the payload
	updateOpts := secrets.UpdateOpts{
		Payload:         d.Get("payload").(string),
		ContentType:     d.Get("payload_content_type").(string),
		ContentEncoding: d.Get("payload_content_encoding").(string),
	}
	err = secrets.Update(kmClient, uuid, updateOpts).Err
	if err != nil {
		return diag.Errorf("Error setting openstack_keymanager_secret_v1 payload: %s", err)
	}

	_, err = stateConf.WaitForStateContext(ctx)
	if err != nil {
		return diag.Errorf("Error waiting for openstack_keymanager_secret_v1: %s", err)
	}

	// set the metadata
	var metadataCreateOpts secrets.MetadataOpts
	metadataCreateOpts = flattenKeyManagerSecretV1Metadata(d)

	log.Printf("[DEBUG] Metadata Create Options for resource_keymanager_secret_metadata_v1 %s: %#v", uuid, metadataCreateOpts)

	if len(metadataCreateOpts) > 0 {
		_, err = secrets.CreateMetadata(kmClient, uuid, metadataCreateOpts).Extract()
		if err != nil {
			return diag.Errorf("Error creating metadata for openstack_keymanager_secret_v1 with ID %s: %s", uuid, err)
		}

		stateConf = &resource.StateChangeConf{
			Pending:    []string{"PENDING"},
			Target:     []string{"ACTIVE"},
			Refresh:    keyManagerSecretMetadataV1WaitForSecretMetadataCreation(kmClient, uuid),
			Timeout:    d.Timeout(schema.TimeoutCreate),
			Delay:      0,
			MinTimeout: 2 * time.Second,
		}

		_, err = stateConf.WaitForStateContext(ctx)
		if err != nil {
			return diag.Errorf("Error creating metadata for openstack_keymanager_secret_v1 %s: %s", uuid, err)
		}
	}

	d.Partial(false)

	return resourceKeyManagerSecretV1Read(ctx, d, meta)
}

func resourceKeyManagerSecretV1Read(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*Config)
	kmClient, err := config.KeyManagerV1Client(GetRegion(d, config))
	if err != nil {
		return diag.Errorf("Error creating OpenStack barbican client: %s", err)
	}

	secret, err := secrets.Get(kmClient, d.Id()).Extract()
	if err != nil {
		return diag.FromErr(CheckDeleted(d, err, "Error retrieving openstack_keymanager_secret_v1"))
	}

	log.Printf("[DEBUG] Retrieved openstack_keymanager_secret_v1 %s: %#v", d.Id(), secret)

	d.Set("name", secret.Name)

	d.Set("bit_length", secret.BitLength)
	d.Set("algorithm", secret.Algorithm)
	d.Set("creator_id", secret.CreatorID)
	d.Set("mode", secret.Mode)
	d.Set("secret_ref", secret.SecretRef)
	d.Set("secret_type", secret.SecretType)
	d.Set("status", secret.Status)
	d.Set("created_at", secret.Created.Format(time.RFC3339))
	d.Set("updated_at", secret.Updated.Format(time.RFC3339))
	d.Set("content_types", secret.ContentTypes)

	// don't fail, if the default key doesn't exist
	payloadContentType, _ := secret.ContentTypes["default"]
	d.Set("payload_content_type", payloadContentType)

	d.Set("payload", keyManagerSecretV1GetPayload(kmClient, d.Id()))
	metadataMap, err := secrets.GetMetadata(kmClient, d.Id()).Extract()
	if err != nil {
		log.Printf("[DEBUG] Unable to get %s secret metadata: %s", d.Id(), err)
	}
	d.Set("all_metadata", metadataMap)

	if secret.Expiration == (time.Time{}) {
		d.Set("expiration", "")
	} else {
		d.Set("expiration", secret.Expiration.Format(time.RFC3339))
	}

	acl, err := acls.GetSecretACL(kmClient, d.Id()).Extract()
	if err != nil {
		log.Printf("[DEBUG] Unable to get %s secret acls: %s", d.Id(), err)
	}
	d.Set("acl", flattenKeyManagerV1ACLs(acl))

	// Set the region
	d.Set("region", GetRegion(d, config))

	return nil
}

func resourceKeyManagerSecretV1Update(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*Config)
	kmClient, err := config.KeyManagerV1Client(GetRegion(d, config))
	if err != nil {
		return diag.Errorf("Error creating OpenStack barbican client: %s", err)
	}

	if d.HasChange("acl") {
		updateOpts := expandKeyManagerV1ACLs(d.Get("acl"))
		_, err := acls.UpdateSecretACL(kmClient, d.Id(), updateOpts).Extract()
		if err != nil {
			return diag.Errorf("Error updating openstack_keymanager_secret_v1 %s acl: %s", d.Id(), err)
		}
	}

	if d.HasChange("metadata") {
		var metadataToDelete []string
		var metadataToAdd []string
		var metadataToUpdate []string

		o, n := d.GetChange("metadata")
		oldMetadata := o.(map[string]interface{})
		newMetadata := n.(map[string]interface{})
		existingMetadata := d.Get("all_metadata").(map[string]interface{})

		// Determine if any metadata keys were removed from the configuration.
		// Then request those keys to be deleted.
		for oldKey := range oldMetadata {
			if _, ok := newMetadata[oldKey]; !ok {
				metadataToDelete = append(metadataToDelete, oldKey)
			}
		}

		log.Printf("[DEBUG] Deleting the following items from metadata for openstack_keymanager_secret_v1 %s: %v", d.Id(), metadataToDelete)

		for _, key := range metadataToDelete {
			err := secrets.DeleteMetadatum(kmClient, d.Id(), key).ExtractErr()
			if err != nil {
				return diag.Errorf("Error deleting openstack_keymanager_secret_v1 %s metadata %s: %s", d.Id(), key, err)
			}
		}

		// Determine if any metadata keys were updated or added in the configuration.
		// Then request those keys to be updated or added.
		for newKey, newValue := range newMetadata {
			if oldValue, ok := oldMetadata[newKey]; ok {
				if newValue != oldValue {
					metadataToUpdate = append(metadataToUpdate, newKey)
				}
			} else if existingValue, ok := existingMetadata[newKey]; ok {
				if newValue != existingValue {
					metadataToUpdate = append(metadataToUpdate, newKey)
				}
			} else {
				metadataToAdd = append(metadataToAdd, newKey)
			}
		}

		log.Printf("[DEBUG] Updating the following items in metadata for openstack_keymanager_secret_v1 %s: %v", d.Id(), metadataToUpdate)

		for _, key := range metadataToUpdate {
			var metadatumOpts secrets.MetadatumOpts
			metadatumOpts.Key = key
			metadatumOpts.Value = newMetadata[key].(string)
			_, err := secrets.UpdateMetadatum(kmClient, d.Id(), metadatumOpts).Extract()
			if err != nil {
				return diag.Errorf("Error updating openstack_keymanager_secret_v1 %s metadata %s: %s", d.Id(), key, err)
			}
		}

		log.Printf("[DEBUG] Adding the following items to metadata for openstack_keymanager_secret_v1 %s: %v", d.Id(), metadataToAdd)

		for _, key := range metadataToAdd {
			var metadatumOpts secrets.MetadatumOpts
			metadatumOpts.Key = key
			metadatumOpts.Value = newMetadata[key].(string)
			err := secrets.CreateMetadatum(kmClient, d.Id(), metadatumOpts).Err
			if err != nil {
				return diag.Errorf("Error adding openstack_keymanager_secret_v1 %s metadata %s: %s", d.Id(), key, err)
			}
		}
	}

	return resourceKeyManagerSecretV1Read(ctx, d, meta)
}

func resourceKeyManagerSecretV1Delete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*Config)
	kmClient, err := config.KeyManagerV1Client(GetRegion(d, config))
	if err != nil {
		return diag.Errorf("Error creating OpenStack barbican client: %s", err)
	}

	stateConf := &resource.StateChangeConf{
		Pending:    []string{"PENDING"},
		Target:     []string{"DELETED"},
		Refresh:    keyManagerSecretV1WaitForSecretDeletion(kmClient, d.Id()),
		Timeout:    d.Timeout(schema.TimeoutDelete),
		Delay:      0,
		MinTimeout: 2 * time.Second,
	}

	if _, err = stateConf.WaitForStateContext(ctx); err != nil {
		return diag.FromErr(err)
	}

	return nil
}
