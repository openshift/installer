package openstack

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/gophercloud/gophercloud/openstack/db/v1/users"
)

func resourceDatabaseUserV1() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceDatabaseUserV1Create,
		ReadContext:   resourceDatabaseUserV1Read,
		DeleteContext: resourceDatabaseUserV1Delete,

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(10 * time.Minute),
			Delete: schema.DefaultTimeout(10 * time.Minute),
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
				Required: true,
				ForceNew: true,
			},

			"instance_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			"password": {
				Type:      schema.TypeString,
				Required:  true,
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
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Set:      schema.HashString,
			},
		},
	}
}

func resourceDatabaseUserV1Create(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*Config)
	DatabaseV1Client, err := config.DatabaseV1Client(GetRegion(d, config))
	if err != nil {
		return diag.Errorf("Error creating OpenStack database client: %s", err)
	}

	userName := d.Get("name").(string)
	rawDatabases := d.Get("databases").(*schema.Set).List()
	instanceID := d.Get("instance_id").(string)

	var usersList users.BatchCreateOpts
	usersList = append(usersList, users.CreateOpts{
		Name:      userName,
		Password:  d.Get("password").(string),
		Host:      d.Get("host").(string),
		Databases: expandDatabaseUserV1Databases(rawDatabases),
	})

	err = users.Create(DatabaseV1Client, instanceID, usersList).ExtractErr()
	if err != nil {
		return diag.Errorf("Error creating openstack_db_user_v1: %s", err)
	}

	stateConf := &resource.StateChangeConf{
		Pending:    []string{"BUILD"},
		Target:     []string{"ACTIVE"},
		Refresh:    databaseUserV1StateRefreshFunc(DatabaseV1Client, instanceID, userName),
		Timeout:    d.Timeout(schema.TimeoutCreate),
		Delay:      10 * time.Second,
		MinTimeout: 3 * time.Second,
	}

	_, err = stateConf.WaitForStateContext(ctx)
	if err != nil {
		return diag.Errorf("Error waiting for openstack_db_user_v1 %s to be created: %s", userName, err)
	}

	// Store the ID now
	d.SetId(fmt.Sprintf("%s/%s", instanceID, userName))

	return resourceDatabaseUserV1Read(ctx, d, meta)
}

func resourceDatabaseUserV1Read(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*Config)
	DatabaseV1Client, err := config.DatabaseV1Client(GetRegion(d, config))
	if err != nil {
		return diag.Errorf("Error creating OpenStack database client: %s", err)
	}

	userID := strings.SplitN(d.Id(), "/", 2)
	if len(userID) != 2 {
		return diag.Errorf("Invalid openstack_db_user_v1 ID: %s", d.Id())
	}

	instanceID := userID[0]
	userName := userID[1]

	exists, userObj, err := databaseUserV1Exists(DatabaseV1Client, instanceID, userName)
	if err != nil {
		return diag.Errorf("Error checking if openstack_db_user_v1 %s exists: %s", d.Id(), err)
	}

	if !exists {
		d.SetId("")
		return nil
	}

	d.Set("name", userName)

	databases := flattenDatabaseUserV1Databases(userObj.Databases)
	if err := d.Set("databases", databases); err != nil {
		return diag.Errorf("Unable to set databases: %s", err)
	}

	return nil
}

func resourceDatabaseUserV1Delete(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*Config)
	DatabaseV1Client, err := config.DatabaseV1Client(GetRegion(d, config))
	if err != nil {
		return diag.Errorf("Error creating OpenStack database client: %s", err)
	}

	userID := strings.SplitN(d.Id(), "/", 2)
	if len(userID) != 2 {
		return diag.Errorf("Invalid openstack_db_user_v1 ID: %s", d.Id())
	}

	instanceID := userID[0]
	userName := userID[1]

	exists, _, err := databaseUserV1Exists(DatabaseV1Client, instanceID, userName)
	if err != nil {
		return diag.Errorf("Error checking if openstack_db_user_v1 %s exists: %s", d.Id(), err)
	}

	if !exists {
		return nil
	}

	err = users.Delete(DatabaseV1Client, instanceID, userName).ExtractErr()
	if err != nil {
		return diag.Errorf("Error deleting openstack_db_user_v1 %s: %s", d.Id(), err)
	}

	return nil
}
