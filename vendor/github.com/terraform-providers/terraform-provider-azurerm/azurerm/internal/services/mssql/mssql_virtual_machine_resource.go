package mssql

import (
	"context"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/Azure/azure-sdk-for-go/services/preview/sqlvirtualmachine/mgmt/2017-03-01-preview/sqlvirtualmachine"
	"github.com/hashicorp/go-azure-helpers/response"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"

	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/tf"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	parseCompute "github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/compute/parse"
	computeValidate "github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/compute/validate"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/mssql/helper"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/mssql/parse"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/mssql/validate"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tags"
	azSchema "github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tf/schema"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tf/suppress"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/timeouts"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func resourceMsSqlVirtualMachine() *schema.Resource {
	return &schema.Resource{
		Create: resourceMsSqlVirtualMachineCreateUpdate,
		Read:   resourceMsSqlVirtualMachineRead,
		Update: resourceMsSqlVirtualMachineCreateUpdate,
		Delete: resourceMsSqlVirtualMachineDelete,

		CustomizeDiff: resourceMsSqlVirtualMachineCustomDiff,

		Importer: azSchema.ValidateResourceIDPriorToImport(func(id string) error {
			_, err := parse.SqlVirtualMachineID(id)
			return err
		}),

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(60 * time.Minute),
			Read:   schema.DefaultTimeout(5 * time.Minute),
			Update: schema.DefaultTimeout(60 * time.Minute),
			Delete: schema.DefaultTimeout(60 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			"virtual_machine_id": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: computeValidate.VirtualMachineID,
			},

			"sql_license_type": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
				ValidateFunc: validation.StringInSlice([]string{
					string(sqlvirtualmachine.PAYG),
					string(sqlvirtualmachine.AHUB),
				}, false),
			},

			"auto_backup": {
				Type:     schema.TypeList,
				Optional: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"encryption_enabled": {
							Type:     schema.TypeBool,
							Optional: true,
							Default:  false,
						},

						"encryption_password": {
							Type:         schema.TypeString,
							Optional:     true,
							Sensitive:    true,
							ValidateFunc: validation.StringIsNotEmpty,
						},

						"manual_schedule": {
							Type:     schema.TypeList,
							Optional: true,
							MaxItems: 1,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"full_backup_frequency": {
										Type:             schema.TypeString,
										Required:         true,
										DiffSuppressFunc: suppress.CaseDifference,
										ValidateFunc: validation.StringInSlice([]string{
											string(sqlvirtualmachine.Daily),
											string(sqlvirtualmachine.Weekly),
										}, false),
									},

									"full_backup_start_hour": {
										Type:         schema.TypeInt,
										Required:     true,
										ValidateFunc: validation.IntBetween(0, 23),
									},

									"full_backup_window_in_hours": {
										Type:         schema.TypeInt,
										Required:     true,
										ValidateFunc: validation.IntBetween(1, 23),
									},

									"log_backup_frequency_in_minutes": {
										Type:         schema.TypeInt,
										Required:     true,
										ValidateFunc: validation.IntBetween(5, 60),
									},
								},
							},
						},

						"retention_period_in_days": {
							Type:         schema.TypeInt,
							Required:     true,
							ValidateFunc: validation.IntBetween(1, 30),
						},

						"storage_blob_endpoint": {
							Type:         schema.TypeString,
							Required:     true,
							ValidateFunc: validation.IsURLWithHTTPS,
						},

						"storage_account_access_key": {
							Type:         schema.TypeString,
							Required:     true,
							ValidateFunc: validation.StringIsNotEmpty,
						},

						"system_databases_backup_enabled": {
							Type:     schema.TypeBool,
							Optional: true,
						},
					},
				},
			},

			"auto_patching": {
				Type:     schema.TypeList,
				Optional: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"day_of_week": {
							Type:     schema.TypeString,
							Required: true,
							ValidateFunc: validation.StringInSlice([]string{
								string(sqlvirtualmachine.Monday),
								string(sqlvirtualmachine.Tuesday),
								string(sqlvirtualmachine.Wednesday),
								string(sqlvirtualmachine.Thursday),
								string(sqlvirtualmachine.Friday),
								string(sqlvirtualmachine.Saturday),
								string(sqlvirtualmachine.Sunday),
							}, false),
						},

						"maintenance_window_duration_in_minutes": {
							Type:         schema.TypeInt,
							Required:     true,
							ValidateFunc: validation.IntBetween(30, 180),
						},

						"maintenance_window_starting_hour": {
							Type:         schema.TypeInt,
							Required:     true,
							ValidateFunc: validation.IntBetween(0, 23),
						},
					},
				},
			},

			"key_vault_credential": {
				Type:     schema.TypeList,
				Optional: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:         schema.TypeString,
							Required:     true,
							ValidateFunc: validation.StringIsNotEmpty,
							// api will add updated credential name, and return "sqlvmName:name1,sqlvmName:name2"
							DiffSuppressFunc: mssqlVMCredentialNameDiffSuppressFunc,
						},

						"key_vault_url": {
							Type:         schema.TypeString,
							Required:     true,
							ForceNew:     true,
							Sensitive:    true,
							ValidateFunc: validation.IsURLWithHTTPS,
						},

						"service_principal_name": {
							Type:         schema.TypeString,
							Required:     true,
							ForceNew:     true,
							Sensitive:    true,
							ValidateFunc: validation.StringIsNotEmpty,
						},

						"service_principal_secret": {
							Type:         schema.TypeString,
							Required:     true,
							ForceNew:     true,
							Sensitive:    true,
							ValidateFunc: validation.StringIsNotEmpty,
						},
					},
				},
			},

			"r_services_enabled": {
				Type:     schema.TypeBool,
				Optional: true,
			},

			"sql_connectivity_port": {
				Type:         schema.TypeInt,
				Optional:     true,
				Default:      1433,
				ValidateFunc: validation.IntBetween(1024, 65535),
			},

			"sql_connectivity_type": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  string(sqlvirtualmachine.PRIVATE),
				ValidateFunc: validation.StringInSlice([]string{
					string(sqlvirtualmachine.LOCAL),
					string(sqlvirtualmachine.PRIVATE),
					string(sqlvirtualmachine.PUBLIC),
				}, false),
			},

			"sql_connectivity_update_password": {
				Type:         schema.TypeString,
				Optional:     true,
				Sensitive:    true,
				ValidateFunc: validation.StringIsNotEmpty,
			},

			"sql_connectivity_update_username": {
				Type:         schema.TypeString,
				Optional:     true,
				Sensitive:    true,
				ValidateFunc: validate.SqlVirtualMachineLoginUserName,
			},

			"storage_configuration": {
				Type:     schema.TypeList,
				Optional: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"disk_type": {
							Type:     schema.TypeString,
							Required: true,
							ValidateFunc: validation.StringInSlice([]string{
								string(sqlvirtualmachine.NEW),
								string(sqlvirtualmachine.EXTEND),
								string(sqlvirtualmachine.ADD),
							}, false),
						},
						"storage_workload_type": {
							Type:     schema.TypeString,
							Required: true,
							ValidateFunc: validation.StringInSlice([]string{
								string(sqlvirtualmachine.GENERAL),
								string(sqlvirtualmachine.OLTP),
								string(sqlvirtualmachine.DW),
							}, false),
						},
						"data_settings":    helper.StorageSettingSchema(),
						"log_settings":     helper.StorageSettingSchema(),
						"temp_db_settings": helper.StorageSettingSchema(),
					},
				},
			},

			"tags": tags.Schema(),
		},
	}
}

func resourceMsSqlVirtualMachineCustomDiff(d *schema.ResourceDiff, _ interface{}) error {
	// ForceNew when removing the auto_backup block.
	// See https://github.com/Azure/azure-rest-api-specs/issues/12818#issuecomment-773727756
	old, new := d.GetChange("auto_backup")
	if len(old.([]interface{})) == 1 && len(new.([]interface{})) == 0 {
		return d.ForceNew("auto_backup")
	}

	encryptionEnabled := d.Get("auto_backup.0.encryption_enabled")
	v, ok := d.GetOk("auto_backup.0.encryption_password")

	if encryptionEnabled.(bool) && (!ok || v.(string) == "") {
		return fmt.Errorf("auto_backup: `encryption_password` is required when `encryption_enabled` is true")
	}

	if !encryptionEnabled.(bool) && ok && v.(string) != "" {
		return fmt.Errorf("auto_backup: `encryption_enabled` must be true when `encryption_password` is set")
	}

	return nil
}

func resourceMsSqlVirtualMachineCreateUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).MSSQL.VirtualMachinesClient
	vmclient := meta.(*clients.Client).Compute.VMClient
	ctx, cancel := timeouts.ForCreate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	vmId := d.Get("virtual_machine_id").(string)
	id, err := parseCompute.VirtualMachineID(vmId)
	if err != nil {
		return err
	}

	if d.IsNewResource() {
		existing, err := client.Get(ctx, id.ResourceGroup, id.Name, "*")
		if err != nil {
			if !utils.ResponseWasNotFound(existing.Response) {
				return fmt.Errorf("checking for present of existing Sql Virtual Machine (Sql Virtual Machine Name %q / Resource Group %q): %+v", id.Name, id.ResourceGroup, err)
			}
		}
		if existing.ID != nil && *existing.ID != "" {
			return tf.ImportAsExistsError("azurerm_mssql_virtual_machine", *existing.ID)
		}
	}

	// get location from vm
	respvm, err := vmclient.Get(ctx, id.ResourceGroup, id.Name, "")
	if err != nil {
		return fmt.Errorf("making Read request on Azure Virtual Machine %s: %+v", id.Name, err)
	}

	if *respvm.Location == "" {
		return fmt.Errorf("Location is empty from making Read request on Azure Virtual Machine %s: %+v", id.Name, err)
	}

	parameters := sqlvirtualmachine.SQLVirtualMachine{
		Location: utils.String(*respvm.Location),
		Properties: &sqlvirtualmachine.Properties{
			VirtualMachineResourceID:   utils.String(d.Get("virtual_machine_id").(string)),
			SQLServerLicenseType:       sqlvirtualmachine.SQLServerLicenseType(d.Get("sql_license_type").(string)),
			SQLManagement:              sqlvirtualmachine.Full,
			AutoBackupSettings:         expandSqlVirtualMachineAutoBackupSettings(d.Get("auto_backup").([]interface{})),
			AutoPatchingSettings:       expandSqlVirtualMachineAutoPatchingSettings(d.Get("auto_patching").([]interface{})),
			KeyVaultCredentialSettings: expandSqlVirtualMachineKeyVaultCredential(d.Get("key_vault_credential").([]interface{})),
			ServerConfigurationsManagementSettings: &sqlvirtualmachine.ServerConfigurationsManagementSettings{
				AdditionalFeaturesServerConfigurations: &sqlvirtualmachine.AdditionalFeaturesServerConfigurations{
					IsRServicesEnabled: utils.Bool(d.Get("r_services_enabled").(bool)),
				},
				SQLConnectivityUpdateSettings: &sqlvirtualmachine.SQLConnectivityUpdateSettings{
					Port:                  utils.Int32(int32(d.Get("sql_connectivity_port").(int))),
					ConnectivityType:      sqlvirtualmachine.ConnectivityType(d.Get("sql_connectivity_type").(string)),
					SQLAuthUpdatePassword: utils.String(d.Get("sql_connectivity_update_password").(string)),
					SQLAuthUpdateUserName: utils.String(d.Get("sql_connectivity_update_username").(string)),
				},
			},
			StorageConfigurationSettings: expandSqlVirtualMachineStorageConfigurationSettings(d.Get("storage_configuration").([]interface{})),
		},
		Tags: tags.Expand(d.Get("tags").(map[string]interface{})),
	}

	future, err := client.CreateOrUpdate(ctx, id.ResourceGroup, id.Name, parameters)
	if err != nil {
		return fmt.Errorf("creating Sql Virtual Machine (Sql Virtual Machine Name %q / Resource Group %q): %+v", id.Name, id.ResourceGroup, err)
	}
	if err = future.WaitForCompletionRef(ctx, client.Client); err != nil {
		return fmt.Errorf("waiting for creation of Sql Virtual Machine (Sql Virtual Machine Name %q / Resource Group %q): %+v", id.Name, id.ResourceGroup, err)
	}

	resp, err := client.Get(ctx, id.ResourceGroup, id.Name, "*")
	if err != nil {
		return fmt.Errorf("retrieving Sql Virtual Machine (Sql Virtual Machine Name %q / Resource Group %q): %+v", id.Name, id.ResourceGroup, err)
	}
	if resp.ID == nil {
		return fmt.Errorf("Cannot read Sql Virtual Machine (Sql Virtual Machine Name %q / Resource Group %q) ID", id.Name, id.ResourceGroup)
	}

	d.SetId(*resp.ID)

	// Wait for the auto backup settings to take effect
	// See: https://github.com/Azure/azure-rest-api-specs/issues/12818
	if autoBackup := d.Get("auto_backup"); (d.IsNewResource() && len(autoBackup.([]interface{})) > 0) || (!d.IsNewResource() && d.HasChange("auto_backup")) {
		log.Printf("[DEBUG] Waiting for SQL Virtual Machine %q AutoBackupSettings to take effect", d.Id())
		stateConf := &resource.StateChangeConf{
			Pending:                   []string{"Retry", "Pending"},
			Target:                    []string{"Updated"},
			Refresh:                   resourceMsSqlVirtualMachineAutoBackupSettingsRefreshFunc(ctx, client, d),
			MinTimeout:                1 * time.Minute,
			ContinuousTargetOccurence: 2,
		}

		if d.IsNewResource() {
			stateConf.Timeout = d.Timeout(schema.TimeoutCreate)
		} else {
			stateConf.Timeout = d.Timeout(schema.TimeoutUpdate)
		}

		if _, err := stateConf.WaitForState(); err != nil {
			return fmt.Errorf("waiting for SQL Virtual Machine %q AutoBackupSettings to take effect: %+v", d.Id(), err)
		}
	}

	return resourceMsSqlVirtualMachineRead(d, meta)
}

func resourceMsSqlVirtualMachineRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).MSSQL.VirtualMachinesClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.SqlVirtualMachineID(d.Id())
	if err != nil {
		return err
	}

	resp, err := client.Get(ctx, id.ResourceGroup, id.Name, "*")
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			log.Printf("[INFO] Sql Virtual Machine %q does not exist - removing from state", d.Id())
			d.SetId("")
			return nil
		}
		return fmt.Errorf("reading Sql Virtual Machine (Sql Virtual Machine Name %q / Resource Group %q): %+v", id.Name, id.ResourceGroup, err)
	}

	if props := resp.Properties; props != nil {
		d.Set("virtual_machine_id", props.VirtualMachineResourceID)
		d.Set("sql_license_type", string(props.SQLServerLicenseType))
		if err := d.Set("auto_backup", flattenSqlVirtualMachineAutoBackup(props.AutoBackupSettings, d)); err != nil {
			return fmt.Errorf("setting `auto_backup`: %+v", err)
		}

		if err := d.Set("auto_patching", flattenSqlVirtualMachineAutoPatching(props.AutoPatchingSettings)); err != nil {
			return fmt.Errorf("setting `auto_patching`: %+v", err)
		}

		if err := d.Set("key_vault_credential", flattenSqlVirtualMachineKeyVaultCredential(props.KeyVaultCredentialSettings, d)); err != nil {
			return fmt.Errorf("setting `key_vault_credential`: %+v", err)
		}

		if mgmtSettings := props.ServerConfigurationsManagementSettings; mgmtSettings != nil {
			if cfgs := mgmtSettings.AdditionalFeaturesServerConfigurations; cfgs != nil {
				d.Set("r_services_enabled", mgmtSettings.AdditionalFeaturesServerConfigurations.IsRServicesEnabled)
			}
			if scus := mgmtSettings.SQLConnectivityUpdateSettings; scus != nil {
				d.Set("sql_connectivity_port", mgmtSettings.SQLConnectivityUpdateSettings.Port)
				d.Set("sql_connectivity_type", mgmtSettings.SQLConnectivityUpdateSettings.ConnectivityType)
			}
		}

		// `storage_configuration.0.storage_workload_type` is in a different spot than the rest of the `storage_configuration`
		// so we'll grab that here and pass it along
		storageWorkloadType := ""
		if props.ServerConfigurationsManagementSettings != nil && props.ServerConfigurationsManagementSettings.SQLWorkloadTypeUpdateSettings != nil {
			storageWorkloadType = string(props.ServerConfigurationsManagementSettings.SQLWorkloadTypeUpdateSettings.SQLWorkloadType)
		}

		if err := d.Set("storage_configuration", flattenSqlVirtualMachineStorageConfigurationSettings(props.StorageConfigurationSettings, storageWorkloadType)); err != nil {
			return fmt.Errorf("error setting `storage_configuration`: %+v", err)
		}
	}
	return tags.FlattenAndSet(d, resp.Tags)
}

func resourceMsSqlVirtualMachineDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).MSSQL.VirtualMachinesClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.SqlVirtualMachineID(d.Id())
	if err != nil {
		return err
	}

	future, err := client.Delete(ctx, id.ResourceGroup, id.Name)
	if err != nil {
		return fmt.Errorf("deleting Sql Virtual Machine (Sql Virtual Machine Name %q / Resource Group %q): %+v", id.Name, id.ResourceGroup, err)
	}

	if err = future.WaitForCompletionRef(ctx, client.Client); err != nil {
		if !response.WasNotFound(future.Response()) {
			return fmt.Errorf("waiting for deleting Sql Virtual Machine (Sql Virtual Machine Name %q / Resource Group %q): %+v", id.Name, id.ResourceGroup, err)
		}
	}

	return nil
}

func resourceMsSqlVirtualMachineAutoBackupSettingsRefreshFunc(ctx context.Context, client *sqlvirtualmachine.SQLVirtualMachinesClient, d *schema.ResourceData) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		id, err := parse.SqlVirtualMachineID(d.Id())
		if err != nil {
			return nil, "Error", err
		}

		resp, err := client.Get(ctx, id.ResourceGroup, id.Name, "*")
		if err != nil {
			return nil, "Retry", err
		}

		if props := resp.Properties; props != nil {
			autoBackupSettings := flattenSqlVirtualMachineAutoBackup(props.AutoBackupSettings, d)

			if len(autoBackupSettings) == 0 {
				// auto backup was nil or disabled in the response
				if v, ok := d.GetOk("auto_backup"); !ok || len(v.([]interface{})) == 0 {
					// also disabled in the config
					return resp, "Updated", nil
				}
				return resp, "Pending", nil
			}

			if v, ok := d.GetOk("auto_backup"); !ok || len(v.([]interface{})) == 0 {
				// still waiting for it to be disabled
				return resp, "Pending", nil
			}

			// check each property in the auto_backup block for drift
			for prop, val := range autoBackupSettings[0].(map[string]interface{}) {
				v := d.Get(fmt.Sprintf("auto_backup.0.%s", prop))
				switch prop {
				case "manual_schedule":
					if m := val.([]interface{}); len(m) > 0 {
						if b, ok := d.GetOk("auto_backup.0.manual_schedule"); !ok || len(b.([]interface{})) == 0 {
							// manual schedule disabled in config but still showing in response
							return resp, "Pending", nil
						}
						// check each property in the manual_schedule block for drift
						for prop2, val2 := range m[0].(map[string]interface{}) {
							v2 := d.Get(fmt.Sprintf("auto_backup.0.manual_schedule.0.%s", prop2))
							switch prop2 {
							case "full_backup_frequency":
								if !strings.EqualFold(v2.(string), val2.(string)) {
									return resp, "Pending", nil
								}
							default:
								if v2 != val2 {
									return resp, "Pending", nil
								}
							}
						}
					} else if b, ok := d.GetOk("auto_backup.0.manual_schedule"); ok || len(b.([]interface{})) > 0 {
						// manual schedule set in config but not reflecting in response
						return resp, "Pending", nil
					}
				default:
					if v != val {
						return resp, "Pending", nil
					}
				}
			}

			return resp, "Updated", nil
		}

		return resp, "Retry", nil
	}
}

func expandSqlVirtualMachineAutoBackupSettings(input []interface{}) *sqlvirtualmachine.AutoBackupSettings {
	ret := sqlvirtualmachine.AutoBackupSettings{
		Enable: utils.Bool(false),
	}

	if len(input) > 0 {
		config := input[0].(map[string]interface{})
		ret.Enable = utils.Bool(true)

		if v, ok := config["retention_period_in_days"]; ok {
			ret.RetentionPeriod = utils.Int32(int32(v.(int)))
		}
		if v, ok := config["storage_blob_endpoint"]; ok {
			ret.StorageAccountURL = utils.String(v.(string))
		}
		if v, ok := config["storage_account_access_key"]; ok {
			ret.StorageAccessKey = utils.String(v.(string))
		}

		v, ok := config["encryption_enabled"]
		enableEncryption := ok && v.(bool)
		ret.EnableEncryption = utils.Bool(enableEncryption)
		if v, ok := config["encryption_password"]; enableEncryption && ok {
			ret.Password = utils.String(v.(string))
		}

		if v, ok := config["system_databases_backup_enabled"]; ok {
			ret.BackupSystemDbs = utils.Bool(v.(bool))
		}

		ret.BackupScheduleType = sqlvirtualmachine.Automated
		if v, ok := config["manual_schedule"]; ok && len(v.([]interface{})) > 0 {
			manualSchedule := v.([]interface{})[0].(map[string]interface{})
			ret.BackupScheduleType = sqlvirtualmachine.Manual

			if v, ok := manualSchedule["full_backup_frequency"]; ok {
				ret.FullBackupFrequency = sqlvirtualmachine.FullBackupFrequencyType(v.(string))
			}
			if v, ok := manualSchedule["full_backup_start_hour"]; ok {
				ret.FullBackupStartTime = utils.Int32(int32(v.(int)))
			}
			if v, ok := manualSchedule["full_backup_window_in_hours"]; ok {
				ret.FullBackupWindowHours = utils.Int32(int32(v.(int)))
			}
			if v, ok := manualSchedule["log_backup_frequency_in_minutes"]; ok {
				ret.LogBackupFrequency = utils.Int32(int32(v.(int)))
			}
		}
	}

	return &ret
}

func flattenSqlVirtualMachineAutoBackup(autoBackup *sqlvirtualmachine.AutoBackupSettings, d *schema.ResourceData) []interface{} {
	if autoBackup == nil || autoBackup.Enable == nil || !*autoBackup.Enable {
		return []interface{}{}
	}

	manualSchedule := make([]interface{}, 0)
	if strings.EqualFold(string(autoBackup.BackupScheduleType), string(sqlvirtualmachine.Manual)) {
		var fullBackupStartHour int
		if autoBackup.FullBackupStartTime != nil {
			fullBackupStartHour = int(*autoBackup.FullBackupStartTime)
		}

		var fullBackupWindowHours int
		if autoBackup.FullBackupWindowHours != nil {
			fullBackupWindowHours = int(*autoBackup.FullBackupWindowHours)
		}

		var logBackupFrequency int
		if autoBackup.LogBackupFrequency != nil {
			logBackupFrequency = int(*autoBackup.LogBackupFrequency)
		}

		manualSchedule = []interface{}{
			map[string]interface{}{
				"full_backup_frequency":           string(autoBackup.FullBackupFrequency),
				"full_backup_start_hour":          fullBackupStartHour,
				"full_backup_window_in_hours":     fullBackupWindowHours,
				"log_backup_frequency_in_minutes": logBackupFrequency,
			},
		}
	}

	var retentionPeriod int
	if autoBackup.RetentionPeriod != nil {
		retentionPeriod = int(*autoBackup.RetentionPeriod)
	}

	// Password, StorageAccessKey, StorageAccountURL are not returned, so we try to copy them
	// from existing config as a best effort.
	encryptionPassword := d.Get("auto_backup.0.encryption_password").(string)
	storageKey := d.Get("auto_backup.0.storage_account_access_key").(string)
	blobEndpoint := d.Get("auto_backup.0.storage_blob_endpoint").(string)

	return []interface{}{
		map[string]interface{}{
			"encryption_enabled":              autoBackup.EnableEncryption != nil && *autoBackup.EnableEncryption,
			"encryption_password":             encryptionPassword,
			"manual_schedule":                 manualSchedule,
			"retention_period_in_days":        retentionPeriod,
			"storage_account_access_key":      storageKey,
			"storage_blob_endpoint":           blobEndpoint,
			"system_databases_backup_enabled": autoBackup.BackupSystemDbs != nil && *autoBackup.BackupSystemDbs,
		},
	}
}

func expandSqlVirtualMachineAutoPatchingSettings(input []interface{}) *sqlvirtualmachine.AutoPatchingSettings {
	if len(input) == 0 {
		return nil
	}
	autoPatchingSetting := input[0].(map[string]interface{})

	return &sqlvirtualmachine.AutoPatchingSettings{
		Enable:                        utils.Bool(true),
		MaintenanceWindowDuration:     utils.Int32(int32(autoPatchingSetting["maintenance_window_duration_in_minutes"].(int))),
		MaintenanceWindowStartingHour: utils.Int32(int32(autoPatchingSetting["maintenance_window_starting_hour"].(int))),
		DayOfWeek:                     sqlvirtualmachine.DayOfWeek(autoPatchingSetting["day_of_week"].(string)),
	}
}

func flattenSqlVirtualMachineAutoPatching(autoPatching *sqlvirtualmachine.AutoPatchingSettings) []interface{} {
	if autoPatching == nil || autoPatching.Enable == nil || !*autoPatching.Enable {
		return []interface{}{}
	}

	var startHour int32
	if autoPatching.MaintenanceWindowStartingHour != nil {
		startHour = *autoPatching.MaintenanceWindowStartingHour
	}

	var duration int32
	if autoPatching.MaintenanceWindowDuration != nil {
		duration = *autoPatching.MaintenanceWindowDuration
	}

	return []interface{}{
		map[string]interface{}{
			"day_of_week":                            string(autoPatching.DayOfWeek),
			"maintenance_window_starting_hour":       startHour,
			"maintenance_window_duration_in_minutes": duration,
		},
	}
}

func expandSqlVirtualMachineKeyVaultCredential(input []interface{}) *sqlvirtualmachine.KeyVaultCredentialSettings {
	if len(input) == 0 {
		return nil
	}
	keyVaultCredentialSetting := input[0].(map[string]interface{})

	return &sqlvirtualmachine.KeyVaultCredentialSettings{
		Enable:                 utils.Bool(true),
		CredentialName:         utils.String(keyVaultCredentialSetting["name"].(string)),
		AzureKeyVaultURL:       utils.String(keyVaultCredentialSetting["key_vault_url"].(string)),
		ServicePrincipalName:   utils.String(keyVaultCredentialSetting["service_principal_name"].(string)),
		ServicePrincipalSecret: utils.String(keyVaultCredentialSetting["service_principal_secret"].(string)),
	}
}

func flattenSqlVirtualMachineKeyVaultCredential(keyVault *sqlvirtualmachine.KeyVaultCredentialSettings, d *schema.ResourceData) []interface{} {
	if keyVault == nil || keyVault.Enable == nil || !*keyVault.Enable {
		return []interface{}{}
	}

	name := ""
	if keyVault.CredentialName != nil {
		name = *keyVault.CredentialName
	}

	keyVaultUrl := ""
	if v, ok := d.GetOk("key_vault_credential.0.key_vault_url"); ok {
		keyVaultUrl = v.(string)
	}

	servicePrincipalName := ""
	if v, ok := d.GetOk("key_vault_credential.0.service_principal_name"); ok {
		servicePrincipalName = v.(string)
	}

	servicePrincipalSecret := ""
	if v, ok := d.GetOk("key_vault_credential.0.service_principal_secret"); ok {
		servicePrincipalSecret = v.(string)
	}

	return []interface{}{
		map[string]interface{}{
			"name":                     name,
			"key_vault_url":            keyVaultUrl,
			"service_principal_name":   servicePrincipalName,
			"service_principal_secret": servicePrincipalSecret,
		},
	}
}

func mssqlVMCredentialNameDiffSuppressFunc(_, old, new string, _ *schema.ResourceData) bool {
	oldNamelist := strings.Split(old, ",")
	for _, n := range oldNamelist {
		cur := strings.Split(n, ":")
		if len(cur) > 1 && cur[1] == new {
			return true
		}
	}
	return false
}

func expandSqlVirtualMachineStorageConfigurationSettings(input []interface{}) *sqlvirtualmachine.StorageConfigurationSettings {
	if len(input) == 0 || input[0] == nil {
		return nil
	}
	storageSettings := input[0].(map[string]interface{})

	return &sqlvirtualmachine.StorageConfigurationSettings{
		DiskConfigurationType: sqlvirtualmachine.DiskConfigurationType(storageSettings["disk_type"].(string)),
		StorageWorkloadType:   sqlvirtualmachine.StorageWorkloadType(storageSettings["storage_workload_type"].(string)),
		SQLDataSettings:       expandSqlVirtualMachineDataStorageSettings(storageSettings["data_settings"].([]interface{})),
		SQLLogSettings:        expandSqlVirtualMachineDataStorageSettings(storageSettings["log_settings"].([]interface{})),
		SQLTempDbSettings:     expandSqlVirtualMachineDataStorageSettings(storageSettings["temp_db_settings"].([]interface{})),
	}
}

func flattenSqlVirtualMachineStorageConfigurationSettings(input *sqlvirtualmachine.StorageConfigurationSettings, storageWorkloadType string) []interface{} {
	if input == nil {
		return []interface{}{}
	}

	output := map[string]interface{}{
		"storage_workload_type": storageWorkloadType,
		"disk_type":             string(input.DiskConfigurationType),
		"data_settings":         flattenSqlVirtualMachineStorageSettings(input.SQLDataSettings),
		"log_settings":          flattenSqlVirtualMachineStorageSettings(input.SQLLogSettings),
		"temp_db_settings":      flattenSqlVirtualMachineStorageSettings(input.SQLTempDbSettings),
	}

	if output["storage_workload_type"].(string) == "" && output["disk_type"] == "" &&
		len(output["data_settings"].([]interface{})) == 0 &&
		len(output["log_settings"].([]interface{})) == 0 &&
		len(output["temp_db_settings"].([]interface{})) == 0 {
		return []interface{}{}
	}

	return []interface{}{output}
}

func expandSqlVirtualMachineDataStorageSettings(input []interface{}) *sqlvirtualmachine.SQLStorageSettings {
	if len(input) == 0 || input[0] == nil {
		return nil
	}
	dataStorageSettings := input[0].(map[string]interface{})

	return &sqlvirtualmachine.SQLStorageSettings{
		Luns:            expandSqlVirtualMachineStorageSettingsLuns(dataStorageSettings["luns"].([]interface{})),
		DefaultFilePath: utils.String(dataStorageSettings["default_file_path"].(string)),
	}
}

func expandSqlVirtualMachineStorageSettingsLuns(input []interface{}) *[]int32 {
	expandedLuns := make([]int32, 0)
	for i := range input {
		if input[i] != nil {
			expandedLuns = append(expandedLuns, int32(input[i].(int)))
		}
	}

	return &expandedLuns
}

func flattenSqlVirtualMachineStorageSettings(input *sqlvirtualmachine.SQLStorageSettings) []interface{} {
	if input == nil || input.Luns == nil {
		return []interface{}{}
	}
	attrs := make(map[string]interface{})

	if input.Luns != nil {
		attrs["luns"] = *input.Luns
	}

	if input.DefaultFilePath != nil {
		attrs["default_file_path"] = *input.DefaultFilePath
	}

	return []interface{}{attrs}
}
