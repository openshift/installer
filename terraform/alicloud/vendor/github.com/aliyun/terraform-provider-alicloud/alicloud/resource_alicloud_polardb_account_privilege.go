package alicloud

import (
	"fmt"
	"strings"
	"time"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
)

func resourceAlicloudPolarDBAccountPrivilege() *schema.Resource {
	return &schema.Resource{
		Create: resourceAlicloudPolarDBAccountPrivilegeCreate,
		Read:   resourceAlicloudPolarDBAccountPrivilegeRead,
		Update: resourceAlicloudPolarDBAccountPrivilegeUpdate,
		Delete: resourceAlicloudPolarDBAccountPrivilegeDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"db_cluster_id": {
				Type:     schema.TypeString,
				ForceNew: true,
				Required: true,
			},

			"account_name": {
				Type:     schema.TypeString,
				ForceNew: true,
				Required: true,
			},

			"account_privilege": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringInSlice([]string{"ReadOnly", "ReadWrite", "DMLOnly", "DDLOnly"}, false),
				Default:      "ReadOnly",
				ForceNew:     true,
			},

			"db_names": {
				Type:     schema.TypeSet,
				Required: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
				MinItems: 1,
			},
		},
	}
}

func resourceAlicloudPolarDBAccountPrivilegeCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	polarDBService := PolarDBService{client}
	clusterId := d.Get("db_cluster_id").(string)
	account := d.Get("account_name").(string)
	privilege := d.Get("account_privilege").(string)
	dbList := d.Get("db_names").(*schema.Set).List()
	// wait instance running before granting
	if err := polarDBService.WaitForPolarDBInstance(clusterId, Running, DefaultLongTimeout); err != nil {
		return WrapError(err)
	}
	d.SetId(fmt.Sprintf("%s%s%s%s%s", clusterId, COLON_SEPARATED, account, COLON_SEPARATED, privilege))

	if len(dbList) > 0 {
		for _, db := range dbList {
			if err := resource.Retry(10*time.Minute, func() *resource.RetryError {
				if err := polarDBService.GrantPolarDBAccountPrivilege(d.Id(), db.(string)); err != nil {
					if IsExpectedErrors(err, OperationDeniedDBStatus) {
						return resource.RetryableError(err)
					}
					return resource.NonRetryableError(err)
				}
				return nil
			}); err != nil {
				return WrapError(err)
			}
		}
	}

	return resourceAlicloudPolarDBAccountPrivilegeRead(d, meta)
}

func resourceAlicloudPolarDBAccountPrivilegeRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	polarDBService := PolarDBService{client}
	parts, err := ParseResourceId(d.Id(), 3)
	if err != nil {
		return WrapError(err)
	}
	object, err := polarDBService.DescribePolarDBAccountPrivilege(d.Id())
	if err != nil {
		if NotFoundError(err) {
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}

	d.Set("db_cluster_id", parts[0])
	d.Set("account_name", object.AccountName)
	d.Set("account_privilege", parts[2])
	var names []string
	for _, pri := range object.DatabasePrivileges {
		if pri.AccountPrivilege == parts[2] {
			names = append(names, pri.DBName)
		}
	}
	d.Set("db_names", names)

	return nil
}

func resourceAlicloudPolarDBAccountPrivilegeUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	PolarDBService := PolarDBService{client}
	d.Partial(true)

	if d.HasChange("db_names") {
		parts := strings.Split(d.Id(), COLON_SEPARATED)

		o, n := d.GetChange("db_names")
		os := o.(*schema.Set)
		ns := n.(*schema.Set)
		remove := os.Difference(ns).List()
		add := ns.Difference(os).List()

		if len(remove) > 0 {
			// wait instance running before revoking
			if err := PolarDBService.WaitForPolarDBInstance(parts[0], Running, DefaultTimeoutMedium); err != nil {
				return WrapError(err)
			}
			for _, db := range remove {
				if err := PolarDBService.RevokePolarDBAccountPrivilege(d.Id(), db.(string)); err != nil {
					return WrapError(err)
				}
			}
		}

		if len(add) > 0 {
			// wait instance running before granting
			if err := PolarDBService.WaitForPolarDBInstance(parts[0], Running, DefaultTimeoutMedium); err != nil {
				return WrapError(err)
			}
			for _, db := range add {
				if err := PolarDBService.GrantPolarDBAccountPrivilege(d.Id(), db.(string)); err != nil {
					return WrapError(err)
				}
			}
		}
		d.SetPartial("db_names")
	}

	d.Partial(false)
	return resourceAlicloudPolarDBAccountPrivilegeRead(d, meta)
}

func resourceAlicloudPolarDBAccountPrivilegeDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	PolarDBService := PolarDBService{client}
	parts, err := ParseResourceId(d.Id(), 3)
	if err != nil {
		return WrapError(err)
	}
	object, err := PolarDBService.DescribePolarDBAccountPrivilege(d.Id())
	if err != nil {
		if NotFoundError(err) {
			return nil
		}
		return WrapError(err)
	}
	var dbName string

	if len(object.DatabasePrivileges) > 0 {
		for _, pri := range object.DatabasePrivileges {
			if pri.AccountPrivilege == parts[2] {
				dbName = pri.DBName
				if err := PolarDBService.RevokePolarDBAccountPrivilege(d.Id(), pri.DBName); err != nil {
					return WrapError(err)
				}
			}
		}
	}

	return PolarDBService.WaitForPolarDBAccountPrivilege(d.Id(), dbName, Deleted, DefaultTimeoutMedium)
}
