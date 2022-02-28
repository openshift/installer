package openstack

import (
	"context"
	"log"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"

	"github.com/gophercloud/gophercloud/openstack/keymanager/v1/orders"
)

func resourceKeyManagerOrderV1() *schema.Resource {
	ret := &schema.Resource{
		CreateContext: resourceKeyManagerOrderV1Create,
		ReadContext:   resourceKeyManagerOrderV1Read,
		DeleteContext: resourceKeyManagerOrderV1Delete,

		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(30 * time.Minute),
			Delete: schema.DefaultTimeout(30 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			"meta": {
				Type:     schema.TypeList,
				Required: true,
				ForceNew: true,
				MaxItems: 1,
				MinItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"algorithm": {
							Type:     schema.TypeString,
							Required: true,
							ForceNew: true,
						},
						"bit_length": {
							Type:     schema.TypeInt,
							Required: true,
							ForceNew: true,
						},
						"expiration": {
							Type:         schema.TypeString,
							Optional:     true,
							ForceNew:     true,
							ValidateFunc: validation.IsRFC3339Time,
						},
						"mode": {
							Type:     schema.TypeString,
							Optional: true,
							ForceNew: true,
						},
						"name": {
							Type:     schema.TypeString,
							Optional: true,
							ForceNew: true,
						},
						"payload_content_type": {
							Type:     schema.TypeString,
							Optional: true,
							ForceNew: true,
							ValidateFunc: validation.StringInSlice([]string{
								"text/plain", "text/plain;charset=utf-8", "text/plain; charset=utf-8", "application/octet-stream", "application/pkcs8",
							}, true),
						},
					},
				},
			},
			"type": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
				ValidateFunc: validation.StringInSlice([]string{
					"asymmetric", "key",
				}, false),
			},
			"container_ref": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"created": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"creator_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"order_ref": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"secret_ref": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"status": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"sub_status": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"sub_status_message": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"updated": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}

	return ret
}

func resourceKeyManagerOrderV1Create(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*Config)
	kmClient, err := config.KeyManagerV1Client(GetRegion(d, config))
	if err != nil {
		return diag.Errorf("Error creating OpenStack KeyManager client: %s", err)
	}

	orderType := keyManagerOrderV1OrderType(d.Get("type").(string))
	metaOpts := expandKeyManagerOrderV1Meta(d.Get("meta").([]interface{}))
	createOpts := orders.CreateOpts{
		Type: orderType,
		Meta: metaOpts,
	}

	log.Printf("[DEBUG] Create Options for resource_keymanager_order_v1: %#v", createOpts)

	var order *orders.Order
	order, err = orders.Create(kmClient, createOpts).Extract()
	if err != nil {
		return diag.Errorf("Error creating openstack_keymanager_order_v1: %s", err)
	}

	uuid := keyManagerOrderV1GetUUIDfromOrderRef(order.OrderRef)

	stateConf := &resource.StateChangeConf{
		Pending:    []string{"PENDING"},
		Target:     []string{"ACTIVE"},
		Refresh:    keyManagerOrderV1WaitForOrderCreation(kmClient, uuid),
		Timeout:    d.Timeout(schema.TimeoutCreate),
		Delay:      0,
		MinTimeout: 2 * time.Second,
	}

	_, err = stateConf.WaitForStateContext(ctx)
	if err != nil {
		return diag.Errorf("Error waiting for openstack_keymanager_order_v1: %s", err)
	}

	d.SetId(uuid)

	return resourceKeyManagerOrderV1Read(ctx, d, meta)
}

func resourceKeyManagerOrderV1Read(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*Config)
	kmClient, err := config.KeyManagerV1Client(GetRegion(d, config))
	if err != nil {
		return diag.Errorf("Error creating OpenStack barbican client: %s", err)
	}

	order, err := orders.Get(kmClient, d.Id()).Extract()
	if err != nil {
		return diag.FromErr(CheckDeleted(d, err, "Error retrieving openstack_keymanager_order_v1"))
	}

	log.Printf("[DEBUG] Retrieved openstack_keymanager_order_v1 %s: %#v", d.Id(), order)

	d.Set("container_ref", order.ContainerRef)
	d.Set("created", order.Created.Format(time.RFC3339))
	d.Set("creator_id", order.CreatorID)
	d.Set("order_ref", order.OrderRef)
	d.Set("secret_ref", order.SecretRef)
	d.Set("status", order.Status)
	d.Set("sub_status", order.SubStatus)
	d.Set("sub_status_message", order.SubStatusMessage)
	d.Set("type", order.Type)
	d.Set("updated", order.Updated.Format(time.RFC3339))
	if err := d.Set("meta", flattenKeyManagerOrderV1Meta(order.Meta)); err != nil {
		return diag.Errorf("error setting meta for resource %s: %s", d.Id(), err)
	}

	return nil
}

func resourceKeyManagerOrderV1Delete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*Config)
	kmClient, err := config.KeyManagerV1Client(GetRegion(d, config))
	if err != nil {
		return diag.Errorf("Error creating OpenStack barbican client: %s", err)
	}

	stateConf := &resource.StateChangeConf{
		Pending:    []string{"PENDING"},
		Target:     []string{"DELETED"},
		Refresh:    keyManagerOrderV1WaitForOrderDeletion(kmClient, d.Id()),
		Timeout:    d.Timeout(schema.TimeoutDelete),
		Delay:      0,
		MinTimeout: 2 * time.Second,
	}

	if _, err = stateConf.WaitForStateContext(ctx); err != nil {
		return diag.FromErr(err)
	}

	return nil
}
