package openstack

import (
	"context"
	"log"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/gophercloud/gophercloud/openstack/db/v1/databases"
	"github.com/gophercloud/gophercloud/openstack/db/v1/instances"
	"github.com/gophercloud/gophercloud/openstack/db/v1/users"
)

func resourceDatabaseInstanceV1() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceDatabaseInstanceV1Create,
		ReadContext:   resourceDatabaseInstanceV1Read,
		DeleteContext: resourceDatabaseInstanceV1Delete,
		UpdateContext: resourceDatabaseInstanceUpdate,

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(30 * time.Minute),
			Delete: schema.DefaultTimeout(30 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},

			"name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			"flavor_id": {
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				Computed:    true,
				DefaultFunc: schema.EnvDefaultFunc("OS_FLAVOR_ID", nil),
			},

			"size": {
				Type:     schema.TypeInt,
				Required: true,
				ForceNew: true,
			},

			"datastore": {
				Type:     schema.TypeList,
				Required: true,
				ForceNew: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"version": {
							Type:     schema.TypeString,
							Required: true,
							ForceNew: true,
						},
						"type": {
							Type:     schema.TypeString,
							Required: true,
							ForceNew: true,
						},
					},
				},
			},

			"network": {
				Type:     schema.TypeList,
				Optional: true,
				ForceNew: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"uuid": {
							Type:     schema.TypeString,
							Optional: true,
							ForceNew: true,
						},
						"port": {
							Type:     schema.TypeString,
							Optional: true,
							ForceNew: true,
						},
						"fixed_ip_v4": {
							Type:     schema.TypeString,
							Optional: true,
							ForceNew: true,
						},
						"fixed_ip_v6": {
							Type:     schema.TypeString,
							Optional: true,
							ForceNew: true,
						},
					},
				},
			},

			"database": {
				Type:     schema.TypeList,
				Optional: true,
				ForceNew: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:     schema.TypeString,
							Required: true,
							ForceNew: true,
						},
						"charset": {
							Type:     schema.TypeString,
							Optional: true,
							ForceNew: true,
						},
						"collate": {
							Type:     schema.TypeString,
							Optional: true,
							ForceNew: true,
						},
					},
				},
			},

			"user": {
				Type:     schema.TypeList,
				Optional: true,
				ForceNew: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:     schema.TypeString,
							Required: true,
							ForceNew: true,
						},
						"password": {
							Type:      schema.TypeString,
							Optional:  true,
							ForceNew:  true,
							Sensitive: true,
						},
						"host": {
							Type:     schema.TypeString,
							Optional: true,
							ForceNew: true,
						},
						"databases": {
							Type:     schema.TypeSet,
							Optional: true,
							ForceNew: true,
							Elem:     &schema.Schema{Type: schema.TypeString},
						},
					},
				},
			},

			"configuration_id": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: false,
				ForceNew: false,
			},

			"addresses": {
				Type:     schema.TypeList,
				Optional: false,
				Computed: true,
				ForceNew: false,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
		},
	}
}

func resourceDatabaseInstanceV1Create(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*Config)
	DatabaseV1Client, err := config.DatabaseV1Client(GetRegion(d, config))
	if err != nil {
		return diag.Errorf("Error creating OpenStack database client: %s", err)
	}

	createOpts := &instances.CreateOpts{
		FlavorRef: d.Get("flavor_id").(string),
		Name:      d.Get("name").(string),
		Size:      d.Get("size").(int),
	}

	// datastore
	var datastore instances.DatastoreOpts
	if v, ok := d.GetOk("datastore"); ok {
		datastore = expandDatabaseInstanceV1Datastore(v.([]interface{}))
	}
	createOpts.Datastore = &datastore

	// networks
	var networks []instances.NetworkOpts
	if v, ok := d.GetOk("network"); ok {
		networks = expandDatabaseInstanceV1Networks(v.([]interface{}))
	}
	createOpts.Networks = networks

	// databases
	var dbs databases.BatchCreateOpts
	if v, ok := d.GetOk("database"); ok {
		dbs = expandDatabaseInstanceV1Databases(v.([]interface{}))
	}
	createOpts.Databases = dbs

	// users
	var userList users.BatchCreateOpts
	if v, ok := d.GetOk("user"); ok {
		userList = expandDatabaseInstanceV1Users(v.([]interface{}))
	}
	createOpts.Users = userList

	log.Printf("[DEBUG] openstack_db_instance_v1 create options: %#v", createOpts)

	instance, err := instances.Create(DatabaseV1Client, createOpts).Extract()
	if err != nil {
		return diag.Errorf("Error creating openstack_db_instance_v1: %s", err)
	}

	// Wait for the instance to become available.
	log.Printf("[DEBUG] Waiting for openstack_db_instance_v1 %s to become available", instance.ID)

	stateConf := &resource.StateChangeConf{
		Pending:    []string{"BUILD"},
		Target:     []string{"ACTIVE", "HEALTHY"},
		Refresh:    databaseInstanceV1StateRefreshFunc(DatabaseV1Client, instance.ID),
		Timeout:    d.Timeout(schema.TimeoutCreate),
		Delay:      10 * time.Second,
		MinTimeout: 3 * time.Second,
	}

	_, err = stateConf.WaitForStateContext(ctx)
	if err != nil {
		return diag.Errorf("Error waiting for openstack_db_instance_v1 %s to become ready: %s", instance.ID, err)
	}

	if configuration, ok := d.GetOk("configuration_id"); ok {
		log.Printf("[DEBUG] Attaching configuration %s to openstack_db_instance_v1 %s", configuration, instance.ID)
		err := instances.AttachConfigurationGroup(DatabaseV1Client, instance.ID, configuration.(string)).ExtractErr()
		if err != nil {
			return diag.Errorf("error attaching configuration group %s to openstack_db_instance_v1 %s: %s",
				configuration, instance.ID, err)
		}
	}

	// Store the ID now
	d.SetId(instance.ID)

	return resourceDatabaseInstanceV1Read(ctx, d, meta)
}

func resourceDatabaseInstanceV1Read(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*Config)
	DatabaseV1Client, err := config.DatabaseV1Client(GetRegion(d, config))
	if err != nil {
		return diag.Errorf("Error creating OpenStack database client: %s", err)
	}

	instance, err := instances.Get(DatabaseV1Client, d.Id()).Extract()
	if err != nil {
		return diag.FromErr(CheckDeleted(d, err, "Error retrieving openstack_db_instance_v1"))
	}

	log.Printf("[DEBUG] Retrieved openstack_db_instance_v1 %s: %#v", d.Id(), instance)

	d.Set("name", instance.Name)
	d.Set("flavor_id", instance.Flavor)
	d.Set("datastore", instance.Datastore)
	d.Set("region", GetRegion(d, config))
	d.Set("addresses", instance.IP)

	return nil
}

func resourceDatabaseInstanceUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*Config)
	DatabaseV1Client, err := config.DatabaseV1Client(GetRegion(d, config))
	if err != nil {
		return diag.Errorf("Error creating OpenStack database client: %s", err)
	}

	if d.HasChange("configuration_id") {
		o, n := d.GetChange("configuration_id")

		err := instances.DetachConfigurationGroup(DatabaseV1Client, d.Id()).ExtractErr()
		if err != nil {
			return diag.FromErr(err)
		}
		log.Printf("Detaching configuration %s from openstack_db_instance_v1 %s", o, d.Id())

		if n != "" {
			err := instances.AttachConfigurationGroup(DatabaseV1Client, d.Id(), n.(string)).ExtractErr()
			if err != nil {
				return diag.FromErr(err)
			}
			log.Printf("Attaching configuration to openstack_db_instance_v1 %s", d.Id())
		}
	}

	return resourceDatabaseInstanceV1Read(ctx, d, meta)
}

func resourceDatabaseInstanceV1Delete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*Config)
	DatabaseV1Client, err := config.DatabaseV1Client(GetRegion(d, config))
	if err != nil {
		return diag.Errorf("Error creating OpenStack database client: %s", err)
	}

	err = instances.Delete(DatabaseV1Client, d.Id()).ExtractErr()
	if err != nil {
		return diag.FromErr(CheckDeleted(d, err, "Error deleting openstack_db_instance_v1"))
	}

	stateConf := &resource.StateChangeConf{
		Pending:    []string{"ACTIVE", "SHUTDOWN"},
		Target:     []string{"DELETED"},
		Refresh:    databaseInstanceV1StateRefreshFunc(DatabaseV1Client, d.Id()),
		Timeout:    d.Timeout(schema.TimeoutDelete),
		Delay:      10 * time.Second,
		MinTimeout: 3 * time.Second,
	}

	_, err = stateConf.WaitForStateContext(ctx)
	if err != nil {
		return diag.Errorf("Error waiting for openstack_db_instance_v1 %s to Delete:  %s", d.Id(), err)
	}

	return nil
}
