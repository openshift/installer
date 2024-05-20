package nutanix

import (
	"context"
	"fmt"
	"log"

	karbon "github.com/terraform-providers/terraform-provider-nutanix/client/karbon"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func resourceNutanixKarbonPrivateRegistry() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceNutanixKarbonPrivateRegistryCreate,
		ReadContext:   resourceNutanixKarbonPrivateRegistryRead,
		UpdateContext: resourceNutanixKarbonPrivateRegistryUpdate,
		DeleteContext: resourceNutanixKarbonPrivateRegistryDelete,
		Exists:        resourceNutanixKarbonPrivateRegistryExists,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		SchemaVersion: 1,
		Schema:        KarbonPrivateRegistryResourceMap(),
	}
}

func KarbonPrivateRegistryResourceMap() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"name": {
			Type:     schema.TypeString,
			Required: true,
			ForceNew: true,
		},
		"cert": {
			Type:     schema.TypeString,
			Optional: true,
			ForceNew: true,
		},
		"url": {
			Type:     schema.TypeString,
			Required: true,
			ForceNew: true,
		},
		"port": {
			Type:         schema.TypeInt,
			Required:     true,
			ValidateFunc: validation.IntAtLeast(1),
		},
		"endpoint": {
			Type:     schema.TypeString,
			Computed: true,
		},
		"username": {
			Type:     schema.TypeString,
			Optional: true,
			ForceNew: true,
		},
		"password": {
			Type:      schema.TypeString,
			Optional:  true,
			Sensitive: true,
			ForceNew:  true,
		},
	}
}

func resourceNutanixKarbonPrivateRegistryCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	log.Print("[Debug] Entering resourceNutanixKarbonPrivateRegistryCreate")
	// Get client connection
	client := meta.(*Client)
	conn := client.KarbonAPI
	setTimeout(meta)
	// Prepare request
	karbonPrivateRegistry := &karbon.PrivateRegistryIntentInput{}
	if name, ok := d.GetOk("name"); ok {
		n := name.(string)
		karbonPrivateRegistry.Name = &n
	} else {
		return diag.Errorf("error occurred during private registry creation: name must be set")
	}
	if url, ok := d.GetOk("url"); ok {
		u := url.(string)
		karbonPrivateRegistry.URL = &u
	} else {
		return diag.Errorf("error occurred during private registry creation: url must be set")
	}
	if port, ok := d.GetOk("port"); ok {
		p := int64(port.(int))
		karbonPrivateRegistry.Port = &p
	}

	if cert, ok := d.GetOk("cert"); ok {
		c := cert.(string)
		karbonPrivateRegistry.Cert = &c
	}
	if username, ok := d.GetOk("username"); ok {
		u := username.(string)
		karbonPrivateRegistry.Username = &u
	}
	if password, ok := d.GetOk("password"); ok {
		pw := password.(string)
		karbonPrivateRegistry.Password = &pw
	}
	createPrivateRegistryResponse, err := conn.PrivateRegistry.CreateKarbonPrivateRegistry(karbonPrivateRegistry)
	if err != nil {
		return diag.Errorf("error occurred during private registry creation: %s", err)
	}

	// Set terraform state id
	d.SetId(*createPrivateRegistryResponse.UUID)
	return resourceNutanixKarbonPrivateRegistryRead(ctx, d, meta)
}

func resourceNutanixKarbonPrivateRegistryRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	log.Print("[Debug] Entering resourceNutanixKarbonPrivateRegistryRead")
	// Get client connection
	conn := meta.(*Client).KarbonAPI
	setTimeout(meta)
	// Make request to the API
	var name interface{}
	var ok bool
	if name, ok = d.GetOk("name"); !ok {
		return diag.Errorf("cannot read private registry without name")
	}
	resp, err := conn.PrivateRegistry.GetKarbonPrivateRegistry(name.(string))
	if err != nil {
		d.SetId("")
		return nil
	}
	if err := d.Set("name", *resp.Name); err != nil {
		return diag.Errorf("error setting name for Karbon private registry %s: %s", d.Id(), err)
	}
	if err := d.Set("endpoint", *resp.Endpoint); err != nil {
		return diag.Errorf("error setting endpoint for Karbon private registry %s: %s", d.Id(), err)
	}
	d.SetId(*resp.UUID)
	return nil
}

func resourceNutanixKarbonPrivateRegistryUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	log.Print("[Debug] Entering resourceNutanixKarbonPrivateRegistryUpdate")
	return resourceNutanixKarbonPrivateRegistryRead(ctx, d, meta)
}

func resourceNutanixKarbonPrivateRegistryDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	log.Print("[Debug] Entering resourceNutanixKarbonPrivateRegistryDelete")
	client := meta.(*Client)
	conn := client.KarbonAPI
	setTimeout(meta)
	karbonPrivateRegistryName := d.Get("name").(string)

	_, err := conn.PrivateRegistry.DeleteKarbonPrivateRegistry(karbonPrivateRegistryName)
	if err != nil {
		return diag.Errorf("error while deleting Karbon Private Registry UUID(%s): %s", d.Id(), err)
	}
	d.SetId("")
	return nil
}

func resourceNutanixKarbonPrivateRegistryExists(d *schema.ResourceData, meta interface{}) (bool, error) {
	log.Print("[DEBUG] Entering resourceNutanixKarbonPrivateRegistryExists")
	conn := meta.(*Client).KarbonAPI
	setTimeout(meta)
	// Make request to the API
	var name interface{}
	var ok bool
	if name, ok = d.GetOk("name"); !ok {
		return false, fmt.Errorf("cannot read private registry without name")
	}
	_, err := conn.PrivateRegistry.GetKarbonPrivateRegistry(name.(string))
	if err != nil {
		d.SetId("")
		return false, nil
	}
	return true, nil
}
