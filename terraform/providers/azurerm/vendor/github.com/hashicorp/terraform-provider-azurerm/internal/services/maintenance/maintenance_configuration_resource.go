package maintenance

import (
	"fmt"
	"log"
	"regexp"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"github.com/hashicorp/go-azure-sdk/resource-manager/maintenance/2022-07-01-preview/maintenanceconfigurations"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/maintenance/migration"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/maintenance/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tags"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

func resourceArmMaintenanceConfiguration() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Create: resourceArmMaintenanceConfigurationCreateUpdate,
		Read:   resourceArmMaintenanceConfigurationRead,
		Update: resourceArmMaintenanceConfigurationCreateUpdate,
		Delete: resourceArmMaintenanceConfigurationDelete,

		SchemaVersion: 1,
		StateUpgraders: pluginsdk.StateUpgrades(map[int]pluginsdk.StateUpgrade{
			0: migration.ConfigurationV0ToV1{},
		}),

		Timeouts: &pluginsdk.ResourceTimeout{
			Create: pluginsdk.DefaultTimeout(30 * time.Minute),
			Read:   pluginsdk.DefaultTimeout(5 * time.Minute),
			Update: pluginsdk.DefaultTimeout(30 * time.Minute),
			Delete: pluginsdk.DefaultTimeout(30 * time.Minute),
		},

		Importer: pluginsdk.ImporterValidatingResourceId(func(id string) error {
			_, err := maintenanceconfigurations.ParseMaintenanceConfigurationIDInsensitively(id)
			return err
		}),

		Schema: map[string]*pluginsdk.Schema{
			"name": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringIsNotEmpty,
			},

			"location": commonschema.Location(),

			"resource_group_name": commonschema.ResourceGroupName(),

			"scope": {
				Type:     pluginsdk.TypeString,
				Required: true,
				ValidateFunc: validation.StringInSlice([]string{
					string(maintenanceconfigurations.MaintenanceScopeExtension),
					string(maintenanceconfigurations.MaintenanceScopeHost),
					string(maintenanceconfigurations.MaintenanceScopeInGuestPatch),
					string(maintenanceconfigurations.MaintenanceScopeOSImage),
					string(maintenanceconfigurations.MaintenanceScopeSQLDB),
					string(maintenanceconfigurations.MaintenanceScopeSQLManagedInstance),
				}, false),
			},

			"visibility": {
				Type:     pluginsdk.TypeString,
				Optional: true,
				Default:  string(maintenanceconfigurations.VisibilityCustom),
				ValidateFunc: validation.StringInSlice([]string{
					string(maintenanceconfigurations.VisibilityCustom),
					// Creating public configurations doesn't appear to be supported, API returns `Public Maintenance Configuration must set correct properties`
					// string(maintenance.VisibilityPublic),
				}, false),
			},

			"window": {
				Type:     pluginsdk.TypeList,
				Optional: true,
				MaxItems: 1,
				Elem: &pluginsdk.Resource{
					Schema: map[string]*pluginsdk.Schema{
						"start_date_time": {
							Type:     pluginsdk.TypeString,
							Required: true,
						},
						"expiration_date_time": {
							Type:     pluginsdk.TypeString,
							Optional: true,
						},
						"duration": {
							Type:     pluginsdk.TypeString,
							Optional: true,
							ValidateFunc: validation.StringMatch(
								regexp.MustCompile("^(0[0-9]|1[0-9]|2[0-3]):[0-5][0-9]$"),
								"duration must match the format HH:mm",
							),
						},
						"time_zone": {
							Type:         pluginsdk.TypeString,
							Required:     true,
							ValidateFunc: validate.MaintenanceTimeZone(),
						},
						"recur_every": {
							Type:         pluginsdk.TypeString,
							Optional:     true,
							ValidateFunc: validation.StringIsNotEmpty,
						},
					},
				},
			},

			"install_patches": {
				Type:     pluginsdk.TypeList,
				MinItems: 1,
				MaxItems: 1,
				Optional: true,
				Elem: &pluginsdk.Resource{
					Schema: map[string]*schema.Schema{
						"linux": {
							Type:     pluginsdk.TypeList,
							MinItems: 1,
							Optional: true,
							Elem: &pluginsdk.Resource{
								Schema: map[string]*schema.Schema{
									"classifications_to_include": {
										Type:     pluginsdk.TypeList,
										Optional: true,
										Elem: &pluginsdk.Schema{
											Type: pluginsdk.TypeString,
											ValidateFunc: validation.StringInSlice([]string{
												"Critical",
												"Security",
												"Other",
											}, false),
										},
									},
									"package_names_mask_to_exclude": {
										Type:     pluginsdk.TypeList,
										Optional: true,
										Elem: &pluginsdk.Schema{
											Type: pluginsdk.TypeString,
										},
									},
									"package_names_mask_to_include": {
										Type:     pluginsdk.TypeList,
										Optional: true,
										Elem: &pluginsdk.Schema{
											Type: pluginsdk.TypeString,
										},
									},
								},
							},
						},
						"windows": {
							Type:     pluginsdk.TypeList,
							Optional: true,
							MinItems: 1,
							Elem: &pluginsdk.Resource{
								Schema: map[string]*schema.Schema{
									"classifications_to_include": {
										Type:     pluginsdk.TypeList,
										Optional: true,
										Elem: &schema.Schema{
											Type: schema.TypeString,
											ValidateFunc: validation.StringInSlice([]string{
												"Critical",
												"Security",
												"UpdateRollup",
												"FeaturePack",
												"ServicePack",
												"Definition",
												"Tools",
												"Updates",
											}, false),
										},
									},
									"kb_numbers_to_exclude": {
										Type:     pluginsdk.TypeList,
										Optional: true,
										Elem: &pluginsdk.Schema{
											Type: pluginsdk.TypeString,
										},
									},
									"kb_numbers_to_include": {
										Type:     pluginsdk.TypeList,
										Optional: true,
										Elem: &pluginsdk.Schema{
											Type: pluginsdk.TypeString,
										},
									},
								},
							},
						},
						"reboot": {
							Type:         pluginsdk.TypeString,
							Optional:     true,
							ValidateFunc: validation.StringInSlice(maintenanceconfigurations.PossibleValuesForRebootOptions(), false),
						},
					},
				},
			},

			"in_guest_user_patch_mode": {
				Type:     pluginsdk.TypeString,
				Optional: true,
				ValidateFunc: validation.StringInSlice([]string{
					"Platform",
					"User",
				}, false),
			},

			"properties": {
				Type:     pluginsdk.TypeMap,
				Optional: true,
				Elem: &pluginsdk.Schema{
					Type:         pluginsdk.TypeString,
					ValidateFunc: validation.StringIsNotEmpty,
				},
			},

			"tags": tags.Schema(),
		},
	}
}

func resourceArmMaintenanceConfigurationCreateUpdate(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Maintenance.ConfigurationsClient
	subscriptionId := meta.(*clients.Client).Account.SubscriptionId
	ctx, cancel := timeouts.ForCreateUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id := maintenanceconfigurations.NewMaintenanceConfigurationID(subscriptionId, d.Get("resource_group_name").(string), d.Get("name").(string))
	if d.IsNewResource() {
		existing, err := client.Get(ctx, id)
		if err != nil {
			if !response.WasNotFound(existing.HttpResponse) {
				return fmt.Errorf("checking for presence of existing %s: %+v", id, err)
			}
		}
		if !response.WasNotFound(existing.HttpResponse) {
			return tf.ImportAsExistsError("azurerm_maintenance_configuration", id.ID())
		}
	}

	scope := maintenanceconfigurations.MaintenanceScope(d.Get("scope").(string))
	visibility := maintenanceconfigurations.Visibility(d.Get("visibility").(string))
	windowRaw := d.Get("window").([]interface{})
	window := expandMaintenanceConfigurationWindow(windowRaw)
	installPatches := expandMaintenanceConfigurationInstallPatches(d.Get("install_patches").([]interface{}))
	extensionProperties := expandExtensionProperties(d.Get("properties").(map[string]interface{}))

	if scope == maintenanceconfigurations.MaintenanceScopeInGuestPatch {
		if window == nil {
			return fmt.Errorf("`window` must be specified when `scope` is `InGuestPatch`")
		}
		if installPatches == nil {
			return fmt.Errorf("`install_patches` must be specified when `scope` is `InGuestPatch`")
		}
		if _, ok := (*extensionProperties)["InGuestPatchMode"]; !ok {
			if _, ok := d.GetOk("in_guest_user_patch_mode"); !ok {
				return fmt.Errorf("`in_guest_user_patch_mode` must be specified when `scope` is `InGuestPatch`")
			}
			(*extensionProperties)["InGuestPatchMode"] = d.Get("in_guest_user_patch_mode").(string)
		}
	}

	configuration := maintenanceconfigurations.MaintenanceConfiguration{
		Name:     utils.String(id.MaintenanceConfigurationName),
		Location: utils.String(location.Normalize(d.Get("location").(string))),
		Properties: &maintenanceconfigurations.MaintenanceConfigurationProperties{
			MaintenanceScope:    &scope,
			Visibility:          &visibility,
			Namespace:           utils.String("Microsoft.Maintenance"),
			MaintenanceWindow:   window,
			ExtensionProperties: extensionProperties,
			InstallPatches:      installPatches,
		},
		Tags: expandTags(d.Get("tags").(map[string]interface{})),
	}

	if _, err := client.CreateOrUpdate(ctx, id, configuration); err != nil {
		return fmt.Errorf("creating/updating %s: %+v", id, err)
	}

	d.SetId(id.ID())
	return resourceArmMaintenanceConfigurationRead(d, meta)
}

func resourceArmMaintenanceConfigurationRead(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Maintenance.ConfigurationsClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := maintenanceconfigurations.ParseMaintenanceConfigurationID(d.Id())
	if err != nil {
		return err
	}

	resp, err := client.Get(ctx, *id)
	if err != nil {
		if response.WasNotFound(resp.HttpResponse) {
			log.Printf("[INFO] maintenance %q does not exist - removing from state", d.Id())
			d.SetId("")
			return nil
		}
		return fmt.Errorf("retrieving %s: %+v", id, err)
	}

	d.Set("name", id.MaintenanceConfigurationName)
	d.Set("resource_group_name", id.ResourceGroupName)

	if model := resp.Model; model != nil {
		if props := model.Properties; props != nil {
			d.Set("scope", props.MaintenanceScope)
			d.Set("visibility", props.Visibility)

			properties := flattenExtensionProperties(props.ExtensionProperties)
			if properties["InGuestPatchMode"] != nil {
				d.Set("in_guest_user_patch_mode", properties["InGuestPatchMode"])
				delete(properties, "InGuestPatchMode")
			}
			d.Set("properties", properties)

			window := flattenMaintenanceConfigurationWindow(props.MaintenanceWindow)
			if err := d.Set("window", window); err != nil {
				return fmt.Errorf("setting `window`: %+v", err)
			}

			installPatches := flattenMaintenanceConfigurationInstallPatches(props.InstallPatches)
			if err := d.Set("install_patches", installPatches); err != nil {
				return fmt.Errorf("setting `install_patches`: %+v", err)
			}
		}
		d.Set("location", location.NormalizeNilable(model.Location))
		if err = tags.FlattenAndSet(d, flattenTags(model.Tags)); err != nil {
			return err
		}
	}
	return nil
}

func resourceArmMaintenanceConfigurationDelete(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Maintenance.ConfigurationsClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := maintenanceconfigurations.ParseMaintenanceConfigurationIDInsensitively(d.Id())
	if err != nil {
		return err
	}

	if _, err := client.Delete(ctx, *id); err != nil {
		return fmt.Errorf("deleting %s: %+v", id, err)
	}
	return nil
}

func expandMaintenanceConfigurationWindow(input []interface{}) *maintenanceconfigurations.MaintenanceWindow {
	if len(input) == 0 {
		return nil
	}

	v := input[0].(map[string]interface{})
	startDateTime := v["start_date_time"].(string)
	expirationDateTime := v["expiration_date_time"].(string)
	duration := v["duration"].(string)
	timeZone := v["time_zone"].(string)
	recurEvery := v["recur_every"].(string)
	window := maintenanceconfigurations.MaintenanceWindow{
		StartDateTime:      utils.String(startDateTime),
		ExpirationDateTime: utils.String(expirationDateTime),
		Duration:           utils.String(duration),
		TimeZone:           utils.String(timeZone),
		RecurEvery:         utils.String(recurEvery),
	}
	return &window
}

func flattenMaintenanceConfigurationWindow(input *maintenanceconfigurations.MaintenanceWindow) []interface{} {
	results := make([]interface{}, 0)

	if v := input; v != nil {
		output := make(map[string]interface{})

		if startDateTime := v.StartDateTime; startDateTime != nil {
			output["start_date_time"] = *startDateTime
		}

		if expirationDateTime := v.ExpirationDateTime; expirationDateTime != nil {
			output["expiration_date_time"] = *expirationDateTime
		}

		if duration := v.Duration; duration != nil {
			output["duration"] = *duration
		}

		if timeZone := v.TimeZone; timeZone != nil {
			output["time_zone"] = *timeZone
		}

		if recurEvery := v.RecurEvery; recurEvery != nil {
			output["recur_every"] = *recurEvery
		}

		results = append(results, output)
	}

	return results
}

func expandMaintenanceConfigurationInstallPatches(input []interface{}) *maintenanceconfigurations.InputPatchConfiguration {
	if len(input) == 0 {
		return nil
	}

	v, ok := input[0].(map[string]interface{})
	if !ok {
		return nil
	}
	rebootSetting := maintenanceconfigurations.RebootOptions(v["reboot"].(string))
	installPatches := maintenanceconfigurations.InputPatchConfiguration{
		WindowsParameters: expandMaintenanceConfigurationInstallPatchesWindows(v["windows"].([]interface{})),
		LinuxParameters:   expandMaintenanceConfigurationInstallPatchesLinux(v["linux"].([]interface{})),
		RebootSetting:     &rebootSetting,
	}
	return &installPatches
}

func flattenMaintenanceConfigurationInstallPatches(input *maintenanceconfigurations.InputPatchConfiguration) []interface{} {
	results := make([]interface{}, 0)

	if v := input; v != nil {
		output := make(map[string]interface{})

		if rebootSetting := v.RebootSetting; rebootSetting != nil {
			output["reboot"] = *rebootSetting
		}

		if windowsParameters := v.WindowsParameters; windowsParameters != nil {
			output["windows"] = flattenMaintenanceConfigurationInstallPatchesWindows(windowsParameters)
		}

		if linuxParameters := v.LinuxParameters; linuxParameters != nil {
			output["linux"] = flattenMaintenanceConfigurationInstallPatchesLinux(linuxParameters)
		}

		results = append(results, output)
	}

	return results
}

func expandMaintenanceConfigurationInstallPatchesWindows(input []interface{}) *maintenanceconfigurations.InputWindowsParameters {
	if len(input) == 0 {
		return nil
	}

	v, ok := input[0].(map[string]interface{})
	if !ok {
		return nil
	}
	windowsInput := maintenanceconfigurations.InputWindowsParameters{}
	if v, ok := v["classifications_to_include"]; ok {
		windowsInput.ClassificationsToInclude = utils.ExpandStringSlice(v.([]interface{}))
	}
	if v, ok := v["kb_numbers_to_exclude"]; ok {
		windowsInput.KbNumbersToExclude = utils.ExpandStringSlice(v.([]interface{}))
	}
	if v, ok := v["kb_numbers_to_include"]; ok {
		windowsInput.KbNumbersToInclude = utils.ExpandStringSlice(v.([]interface{}))
	}
	return &windowsInput
}

func flattenMaintenanceConfigurationInstallPatchesWindows(input *maintenanceconfigurations.InputWindowsParameters) []interface{} {
	results := make([]interface{}, 0)

	if v := input; v != nil {
		output := make(map[string]interface{})

		if classificationsToInclude := v.ClassificationsToInclude; classificationsToInclude != nil {
			output["classifications_to_include"] = utils.FlattenStringSlice(classificationsToInclude)
		}

		if kbNumbersToExclude := v.KbNumbersToExclude; kbNumbersToExclude != nil {
			output["kb_numbers_to_exclude"] = utils.FlattenStringSlice(kbNumbersToExclude)
		}

		if kbNumbersToInclude := v.KbNumbersToInclude; kbNumbersToInclude != nil {
			output["kb_numbers_to_include"] = utils.FlattenStringSlice(kbNumbersToInclude)
		}

		results = append(results, output)
	}

	return results
}

func expandMaintenanceConfigurationInstallPatchesLinux(input []interface{}) *maintenanceconfigurations.InputLinuxParameters {
	if len(input) == 0 {
		return nil
	}

	v, ok := input[0].(map[string]interface{})
	if !ok {
		return nil
	}
	linuxParameters := maintenanceconfigurations.InputLinuxParameters{}
	if v, ok := v["classifications_to_include"]; ok {
		linuxParameters.ClassificationsToInclude = utils.ExpandStringSlice(v.([]interface{}))
	}
	if v, ok := v["packages_to_exclude"]; ok {
		linuxParameters.PackageNameMasksToExclude = utils.ExpandStringSlice(v.([]interface{}))
	}
	if v, ok := v["packages_to_include"]; ok {
		linuxParameters.PackageNameMasksToInclude = utils.ExpandStringSlice(v.([]interface{}))
	}
	return &linuxParameters
}

func flattenMaintenanceConfigurationInstallPatchesLinux(input *maintenanceconfigurations.InputLinuxParameters) []interface{} {
	results := make([]interface{}, 0)

	if v := input; v != nil {
		output := make(map[string]interface{})

		if classificationsToInclude := v.ClassificationsToInclude; classificationsToInclude != nil {
			output["classifications_to_include"] = utils.FlattenStringSlice(classificationsToInclude)
		}

		if packageNameMasksToInclude := v.PackageNameMasksToInclude; packageNameMasksToInclude != nil {
			output["packages_to_exclude"] = utils.FlattenStringSlice(packageNameMasksToInclude)
		}

		if packageNameMasksToExclude := v.PackageNameMasksToExclude; packageNameMasksToExclude != nil {
			output["packages_to_include"] = utils.FlattenStringSlice(packageNameMasksToExclude)
		}

		results = append(results, output)
	}

	return results
}

func expandExtensionProperties(input map[string]interface{}) *map[string]string {
	output := make(map[string]string)
	for k, v := range input {
		output[k] = v.(string)
	}
	return &output
}

func flattenExtensionProperties(input *map[string]string) map[string]interface{} {
	output := make(map[string]interface{})
	if input != nil {
		for k, v := range *input {
			if k == "inGuestPatchMode" {
				output["InGuestPatchMode"] = v
				continue
			}
			output[k] = v
		}
	}
	return output
}
