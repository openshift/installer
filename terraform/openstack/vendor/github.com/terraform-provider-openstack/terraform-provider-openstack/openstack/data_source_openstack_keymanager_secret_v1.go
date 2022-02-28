package openstack

import (
	"context"
	"fmt"
	"log"
	"regexp"
	"strings"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"

	"github.com/gophercloud/gophercloud/openstack/keymanager/v1/acls"
	"github.com/gophercloud/gophercloud/openstack/keymanager/v1/secrets"
)

func getDateFilters() [4]string {
	return [4]string{
		string(secrets.DateFilterGT),
		string(secrets.DateFilterGTE),
		string(secrets.DateFilterLT),
		string(secrets.DateFilterLTE),
	}
}

func getDateFiltersRegexPreformatted() string {
	df := getDateFilters()
	return strings.Join(df[:], "|")
}

func dataSourceKeyManagerSecretV1() *schema.Resource {
	ret := &schema.Resource{
		ReadContext: dataSourceKeyManagerSecretV1Read,

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
			},

			"name": {
				Type:     schema.TypeString,
				Optional: true,
			},

			"bit_length": {
				Type:     schema.TypeInt,
				Optional: true,
			},

			"algorithm": {
				Type:     schema.TypeString,
				Optional: true,
			},

			"mode": {
				Type:     schema.TypeString,
				Optional: true,
			},

			"secret_type": {
				Type:     schema.TypeString,
				Optional: true,
				ValidateFunc: validation.StringInSlice([]string{
					"symmetric", "public", "private", "passphrase", "certificate", "opaque",
				}, false),
			},

			"acl_only": {
				Type:     schema.TypeBool,
				Optional: true,
			},

			"expiration_filter": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: dataSourceValidateDateFilter,
			},

			"created_at_filter": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: dataSourceValidateDateFilter,
			},

			"updated_at_filter": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: dataSourceValidateDateFilter,
			},

			// computed
			"acl": {
				Type:     schema.TypeList,
				Computed: true,
			},

			"secret_ref": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"creator_id": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"expiration": {
				Type:     schema.TypeString,
				Computed: true,
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

			"status": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"payload": {
				Type:      schema.TypeString,
				Computed:  true,
				Sensitive: true,
			},

			"payload_content_type": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"payload_content_encoding": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"metadata": {
				Type:     schema.TypeMap,
				Computed: true,
			},
		},
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

func dataSourceKeyManagerSecretV1Read(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*Config)
	kmClient, err := config.KeyManagerV1Client(GetRegion(d, config))
	if err != nil {
		return diag.Errorf("Error creating OpenStack barbican client: %s", err)
	}

	aclOnly := d.Get("acl_only").(bool)

	listOpts := secrets.ListOpts{
		Name:            d.Get("name").(string),
		Bits:            d.Get("bit_length").(int),
		Alg:             d.Get("algorithm").(string),
		Mode:            d.Get("mode").(string),
		SecretType:      secrets.SecretType(d.Get("secret_type").(string)),
		ACLOnly:         &aclOnly,
		CreatedQuery:    dataSourceParseDateFilter(d.Get("created_at_filter").(string)),
		UpdatedQuery:    dataSourceParseDateFilter(d.Get("updated_at_filter").(string)),
		ExpirationQuery: dataSourceParseDateFilter(d.Get("expiration_filter").(string)),
	}

	log.Printf("[DEBUG] %#+v List Options: %#v", dataSourceParseDateFilter(d.Get("updated_at_filter").(string)), listOpts)

	allPages, err := secrets.List(kmClient, listOpts).AllPages()
	if err != nil {
		return diag.Errorf("Unable to query openstack_keymanager_secret_v1 secrets: %s", err)
	}

	allSecrets, err := secrets.ExtractSecrets(allPages)
	if err != nil {
		return diag.Errorf("Unable to retrieve openstack_keymanager_secret_v1 secrets: %s", err)
	}

	if len(allSecrets) < 1 {
		return diag.Errorf("Your query returned no openstack_keymanager_secret_v1 results. " +
			"Please change your search criteria and try again")
	}

	if len(allSecrets) > 1 {
		log.Printf("[DEBUG] Multiple openstack_keymanager_secret_v1 results found: %#v", allSecrets)
		return diag.Errorf("Your query returned more than one result. Please try a more " +
			"specific search criteria")
	}

	secret := allSecrets[0]

	log.Printf("[DEBUG] Retrieved openstack_keymanager_secret_v1 %s: %#v", d.Id(), secret)

	uuid := keyManagerSecretV1GetUUIDfromSecretRef(secret.SecretRef)

	d.SetId(uuid)
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
		log.Printf("[DEBUG] Unable to get %s secret metadata: %s", uuid, err)
	}
	d.Set("metadata", metadataMap)

	if secret.Expiration == (time.Time{}) {
		d.Set("expiration", "")
	} else {
		d.Set("expiration", secret.Expiration.Format(time.RFC3339))
	}

	acl, err := acls.GetSecretACL(kmClient, d.Id()).Extract()
	if err != nil {
		log.Printf("[DEBUG] Unable to get %s secret acls: %s", uuid, err)
	}
	d.Set("acl", flattenKeyManagerV1ACLs(acl))

	// Set the region
	d.Set("region", GetRegion(d, config))

	return nil
}

func dataSourceParseDateFilter(date string) *secrets.DateQuery {
	// error checks are not necessary, since they were validated by terraform validate functions
	var parts []string
	if regexp.MustCompile("^" + getDateFiltersRegexPreformatted() + ":").Match([]byte(date)) {
		parts = strings.SplitN(date, ":", 2)
	} else {
		parts = []string{date}
	}

	var parsedTime time.Time
	var filter *secrets.DateQuery

	if len(parts) == 2 {
		parsedTime, _ = time.Parse(time.RFC3339, parts[1])

		filter = &secrets.DateQuery{Date: parsedTime, Filter: secrets.DateFilter(parts[0])}
	} else {
		parsedTime, _ = time.Parse(time.RFC3339, parts[0])

		filter = &secrets.DateQuery{Date: parsedTime}
	}

	if parsedTime == (time.Time{}) {
		return nil
	}

	return filter
}

func dataSourceValidateDateFilter(v interface{}, k string) (ws []string, errors []error) {
	var parts []string
	if regexp.MustCompile("^" + getDateFiltersRegexPreformatted() + ":").Match([]byte(v.(string))) {
		parts = strings.SplitN(v.(string), ":", 2)
	} else {
		parts = []string{v.(string)}
	}

	if len(parts) == 2 {
		supportedDateFilters := getDateFilters()
		if !strSliceContains(supportedDateFilters[:], parts[0]) {
			errors = append(errors, fmt.Errorf("Invalid %q date filter, supported: %+q", parts[0], supportedDateFilters))
		}

		_, err := time.Parse(time.RFC3339, parts[1])
		if err != nil {
			errors = append(errors, err)
		}

		return
	}

	_, err := time.Parse(time.RFC3339, parts[0])
	if err != nil {
		errors = append(errors, err)
	}

	return
}
