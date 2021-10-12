package alicloud

import (
	"fmt"
	"log"
	"strings"
	"time"

	util "github.com/alibabacloud-go/tea-utils/service"

	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceAlicloudDBAccountPrivilege() *schema.Resource {
	return &schema.Resource{
		Create: resourceAlicloudDBAccountPrivilegeCreate,
		Read:   resourceAlicloudDBAccountPrivilegeRead,
		Update: resourceAlicloudDBAccountPrivilegeUpdate,
		Delete: resourceAlicloudDBAccountPrivilegeDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"instance_id": {
				Type:     schema.TypeString,
				ForceNew: true,
				Required: true,
			},

			"account_name": {
				Type:     schema.TypeString,
				ForceNew: true,
				Required: true,
			},

			"privilege": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringInSlice([]string{"ReadOnly", "ReadWrite", "DDLOnly", "DMLOnly", "DBOwner"}, false),
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

func resourceAlicloudDBAccountPrivilegeCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	rdsService := RdsService{client}
	instanceId := d.Get("instance_id").(string)
	account := d.Get("account_name").(string)
	privilege := d.Get("privilege").(string)
	dbList := d.Get("db_names").(*schema.Set).List()
	// wait instance running before granting
	if err := rdsService.WaitForDBInstance(instanceId, Running, DefaultLongTimeout); err != nil {
		return WrapError(err)
	}
	d.SetId(fmt.Sprintf("%s%s%s%s%s", instanceId, COLON_SEPARATED, account, COLON_SEPARATED, privilege))

	if len(dbList) > 0 {
		for _, db := range dbList {
			if err := resource.Retry(10*time.Minute, func() *resource.RetryError {
				if err := rdsService.GrantAccountPrivilege(d.Id(), db.(string)); err != nil {
					if IsExpectedErrors(err, OperationDeniedDBStatus) || NeedRetry(err) {
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

	return resourceAlicloudDBAccountPrivilegeRead(d, meta)
}

func resourceAlicloudDBAccountPrivilegeRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	rsdService := RdsService{client}
	parts, err := ParseResourceId(d.Id(), 3)
	if err != nil {
		return WrapError(err)
	}
	if !d.IsNewResource() && len(d.Get("db_names").(*schema.Set).List()) < 0 {
		log.Println("[WARN] Resource alicloud_db_account_privilege has no any databaces in this state.")
		return nil
	}
	object, err := rsdService.DescribeDBAccountPrivilege(d.Id())
	if err != nil {
		if NotFoundError(err) {
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}

	d.Set("instance_id", object["DBInstanceId"])
	d.Set("account_name", object["AccountName"])
	d.Set("privilege", parts[2])
	var names []string
	databasePrivilege := object["DatabasePrivileges"].(map[string]interface{})["DatabasePrivilege"].([]interface{})
	for _, pri := range databasePrivilege {
		pri := pri.(map[string]interface{})
		if pri["AccountPrivilege"] == parts[2] {
			if len(d.Get("db_names").(*schema.Set).List()) > 0 {
				for _, name := range d.Get("db_names").(*schema.Set).List() {
					if pri["DBName"].(string) == name.(string) {
						names = append(names, pri["DBName"].(string))
						break
					}
				}
			} else {
				names = append(names, pri["DBName"].(string))
			}
		}
	}

	if len(names) < 1 && strings.HasPrefix(object["DBInstanceId"].(string), "pgm-") {
		conn, err := client.NewRdsClient()
		if err != nil {
			return WrapError(err)
		}
		action := "DescribeDatabases"
		request := map[string]interface{}{
			"RegionId":     client.RegionId,
			"DBInstanceId": object["DBInstanceId"],
			"SourceIp":     client.SourceIp,
		}
		runtime := util.RuntimeOptions{}
		runtime.SetAutoretry(true)
		err = resource.Retry(5*time.Minute, func() *resource.RetryError {
			response, err := conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2014-08-15"), StringPointer("AK"), nil, request, &runtime)
			if err != nil {
				if IsExpectedErrors(err, []string{"InternalError", "OperationDenied.DBInstanceStatus"}) {
					return resource.RetryableError(WrapErrorf(err, DefaultErrorMsg, object["DBInstanceId"], action, AlibabaCloudSdkGoERROR))
				}
				return resource.NonRetryableError(WrapErrorf(err, DefaultErrorMsg, object["DBInstanceId"], action, AlibabaCloudSdkGoERROR))
			}
			addDebug(action, response, request)
			dataBases := response["Databases"].(map[string]interface{})["Database"].([]interface{})
			for _, db := range dataBases {
				db := db.(map[string]interface{})
				accountPrivilegeInfos := db["Accounts"].(map[string]interface{})["AccountPrivilegeInfo"].([]interface{})
				for _, account := range accountPrivilegeInfos {
					account := account.(map[string]interface{})
					if account["Account"] == object["AccountName"] && (account["AccountPrivilege"] == parts[2] || account["AccountPrivilege"] == "ALL") {
						names = append(names, db["DBName"].(string))
					}
				}
			}
			return nil
		})
		if err != nil {
			return WrapError(err)
		}
	}

	d.Set("db_names", names)

	return nil
}

func resourceAlicloudDBAccountPrivilegeUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	rdsService := RdsService{client}
	d.Partial(true)

	if d.HasChange("db_names") {
		parts := strings.Split(d.Id(), COLON_SEPARATED)

		o, n := d.GetChange("db_names")
		os := o.(*schema.Set)
		ns := n.(*schema.Set)
		remove := os.Difference(ns).List()
		add := ns.Difference(os).List()

		if len(remove) > 0 {
			if strings.HasPrefix(d.Id(), "pgm-") {
				return WrapError(fmt.Errorf("At present, the PostgreSql database does not support revoking the current privilege."))
			}
			// wait instance running before revoking
			if err := rdsService.WaitForDBInstance(parts[0], Running, DefaultTimeoutMedium); err != nil {
				return WrapError(err)
			}
			for _, db := range remove {
				if err := rdsService.RevokeAccountPrivilege(d.Id(), db.(string)); err != nil {
					return WrapError(err)
				}
			}
		}

		if len(add) > 0 {
			// wait instance running before granting
			if err := rdsService.WaitForDBInstance(parts[0], Running, DefaultTimeoutMedium); err != nil {
				return WrapError(err)
			}
			for _, db := range add {
				if err := rdsService.GrantAccountPrivilege(d.Id(), db.(string)); err != nil {
					return WrapError(err)
				}
			}
		}
		d.SetPartial("db_names")
	}

	d.Partial(false)
	return resourceAlicloudDBAccountPrivilegeRead(d, meta)
}

func resourceAlicloudDBAccountPrivilegeDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	rdsService := RdsService{client}
	parts, err := ParseResourceId(d.Id(), 3)
	if err != nil {
		return WrapError(err)
	}
	object, err := rdsService.DescribeDBAccountPrivilege(d.Id())
	if err != nil {
		if NotFoundError(err) {
			return nil
		}
		return WrapError(err)
	}
	if strings.HasPrefix(d.Id(), "pgm-") {
		return nil
	}
	dbNames := make([]string, 0)
	databasePrivileges := object["DatabasePrivileges"].(map[string]interface{})["DatabasePrivilege"].([]interface{})
	if len(databasePrivileges) > 0 && len(d.Get("db_names").(*schema.Set).List()) > 0 {
		for _, pri := range databasePrivileges {
			pri := pri.(map[string]interface{})
			if pri["AccountPrivilege"] == parts[2] {
				dbName := pri["DBName"].(string)
				for _, name := range d.Get("db_names").(*schema.Set).List() {
					if dbName == name.(string) {
						if err := rdsService.RevokeAccountPrivilege(d.Id(), dbName); err != nil {
							return WrapError(err)
						}
						dbNames = append(dbNames, dbName)
						break
					}
				}
			}
		}
	}

	for _, dbName := range dbNames {
		if err := rdsService.WaitForAccountPrivilege(d.Id(), dbName, Deleted, DefaultTimeoutMedium); err != nil {
			return WrapError(err)
		}
	}
	return nil
}
