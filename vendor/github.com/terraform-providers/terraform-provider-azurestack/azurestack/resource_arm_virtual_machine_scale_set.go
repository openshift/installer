package azurestack

import (
	"bytes"
	"fmt"
	"log"
	"strings"

	"github.com/Azure/azure-sdk-for-go/profiles/2019-03-01/compute/mgmt/compute"
	"github.com/hashicorp/terraform-plugin-sdk/helper/hashcode"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/structure"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/suppress"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
	"github.com/terraform-providers/terraform-provider-azurestack/azurestack/helpers/azure"
)

func resourceArmVirtualMachineScaleSet() *schema.Resource {
	return &schema.Resource{
		Create: resourceArmVirtualMachineScaleSetCreate,
		Read:   resourceArmVirtualMachineScaleSetRead,
		Update: resourceArmVirtualMachineScaleSetCreate,
		Delete: resourceArmVirtualMachineScaleSetDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.NoZeroValues,
			},

			"location": locationSchema(),

			"resource_group_name": resourceGroupNameSchema(),

			// Not supported on profile: 2017-03-09
			// "zones": zonesSchema(),

			"identity": {
				Type:     schema.TypeList,
				Optional: true,
				Computed: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"type": {
							Type:             schema.TypeString,
							Required:         true,
							DiffSuppressFunc: suppress.CaseDifference,
							ValidateFunc: validation.StringInSlice([]string{
								"SystemAssigned",
							}, true),
						},
						"principal_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},

			"sku": {
				Type:     schema.TypeSet,
				Required: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:         schema.TypeString,
							Required:     true,
							ValidateFunc: validation.NoZeroValues,
						},

						"tier": {
							Type:             schema.TypeString,
							Optional:         true,
							Computed:         true,
							DiffSuppressFunc: suppress.CaseDifference,
						},

						"capacity": {
							Type:         schema.TypeInt,
							Required:     true,
							ValidateFunc: validation.IntAtLeast(0),
						},
					},
				},
				Set: resourceArmVirtualMachineScaleSetSkuHash,
			},

			"license_type": {
				Type:             schema.TypeString,
				Optional:         true,
				Computed:         true,
				DiffSuppressFunc: suppress.CaseDifference,
				ValidateFunc: validation.StringInSlice([]string{
					"Windows_Client",
					"Windows_Server",
				}, true),
			},

			"upgrade_policy_mode": {
				Type:     schema.TypeString,
				Required: true,
				ValidateFunc: validation.StringInSlice([]string{
					string(compute.Automatic),
					string(compute.Manual),
					// This Enum is not defined on the 2016-03-30 service
					// Using a regular string instead
					"Rolling",
				}, true),
			},

			"overprovision": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  true,
			},

			"single_placement_group": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  true,
				ForceNew: true,
			},

			"priority": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				Default:  "Regular",
				ValidateFunc: validation.StringInSlice([]string{
					// The Enums are not defined on the 2016-03-30 service
					// Using a regular string instead
					"Low",
					"Regular",
				}, true),
			},

			"os_profile": {
				Type:     schema.TypeList,
				Required: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"computer_name_prefix": {
							Type:     schema.TypeString,
							Required: true,
							ForceNew: true,
						},

						"admin_username": {
							Type:         schema.TypeString,
							Required:     true,
							ValidateFunc: validation.NoZeroValues,
						},

						"admin_password": {
							Type:         schema.TypeString,
							Optional:     true,
							Sensitive:    true,
							ValidateFunc: validation.NoZeroValues,
						},

						"custom_data": {
							Type:             schema.TypeString,
							Optional:         true,
							StateFunc:        userDataStateFunc,
							DiffSuppressFunc: userDataDiffSuppressFunc,
						},
					},
				},
			},

			"os_profile_secrets": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"source_vault_id": {
							Type:     schema.TypeString,
							Required: true,
						},

						"vault_certificates": {
							Type:     schema.TypeList,
							Optional: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"certificate_url": {
										Type:     schema.TypeString,
										Required: true,
									},
									"certificate_store": {
										Type:     schema.TypeString,
										Optional: true,
									},
								},
							},
						},
					},
				},
			},

			"os_profile_windows_config": {
				Type:     schema.TypeSet,
				Optional: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"provision_vm_agent": {
							Type:     schema.TypeBool,
							Optional: true,
						},
						"enable_automatic_upgrades": {
							Type:     schema.TypeBool,
							Optional: true,
						},
						"winrm": {
							Type:     schema.TypeList,
							Optional: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"protocol": {
										Type:     schema.TypeString,
										Required: true,
									},
									"certificate_url": {
										Type:     schema.TypeString,
										Optional: true,
									},
								},
							},
						},
						"additional_unattend_config": {
							Type:     schema.TypeList,
							Optional: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"pass": {
										Type:     schema.TypeString,
										Required: true,
									},
									"component": {
										Type:     schema.TypeString,
										Required: true,
									},
									"setting_name": {
										Type:     schema.TypeString,
										Required: true,
									},
									"content": {
										Type:     schema.TypeString,
										Required: true,
									},
								},
							},
						},
					},
				},
				Set: resourceArmVirtualMachineScaleSetOsProfileWindowsConfigHash,
			},

			"os_profile_linux_config": {
				Type:     schema.TypeSet,
				Optional: true,
				Computed: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"disable_password_authentication": {
							Type:     schema.TypeBool,
							Optional: true,
							Default:  false,
							ForceNew: true,
						},
						"ssh_keys": {
							Type:     schema.TypeList,
							Optional: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"path": {
										Type:     schema.TypeString,
										Required: true,
									},
									"key_data": {
										Type:     schema.TypeString,
										Optional: true,
									},
								},
							},
						},
					},
				},
				Set: resourceArmVirtualMachineScaleSetOsProfileLinuxConfigHash,
			},

			"network_profile": {
				Type:     schema.TypeSet,
				Required: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:         schema.TypeString,
							Required:     true,
							ValidateFunc: validation.NoZeroValues,
						},

						"primary": {
							Type:     schema.TypeBool,
							Required: true,
						},

						"network_security_group_id": {
							Type:         schema.TypeString,
							Optional:     true,
							ValidateFunc: azure.ValidateResourceID,
						},

						"ip_configuration": {
							Type:     schema.TypeList,
							Required: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"name": {
										Type:         schema.TypeString,
										Required:     true,
										ValidateFunc: validation.NoZeroValues,
									},

									"subnet_id": {
										Type:         schema.TypeString,
										Required:     true,
										ValidateFunc: azure.ValidateResourceID,
									},

									"application_gateway_backend_address_pool_ids": {
										Type:     schema.TypeSet,
										Optional: true,
										Elem:     &schema.Schema{Type: schema.TypeString},
										Set:      schema.HashString,
									},

									"load_balancer_backend_address_pool_ids": {
										Type:     schema.TypeSet,
										Optional: true,
										Elem:     &schema.Schema{Type: schema.TypeString},
										Set:      schema.HashString,
									},

									"load_balancer_inbound_nat_rules_ids": {
										Type:     schema.TypeSet,
										Optional: true,
										Computed: true,
										Elem:     &schema.Schema{Type: schema.TypeString},
										Set:      schema.HashString,
									},

									"primary": {
										Type:     schema.TypeBool,
										Optional: true,
									},

									"public_ip_address_configuration": {
										Type:     schema.TypeList,
										Optional: true,
										MaxItems: 1,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"name": {
													Type:     schema.TypeString,
													Required: true,
												},

												"idle_timeout": {
													Type:         schema.TypeInt,
													Required:     true,
													ValidateFunc: validation.IntBetween(4, 32),
												},

												"domain_name_label": {
													Type:     schema.TypeString,
													Required: true,
												},
											},
										},
									},
								},
							},
						},
					},
				},
				Set: resourceArmVirtualMachineScaleSetNetworkConfigurationHash,
			},

			//Not supported on profile: 2017-03-09
			// "boot_diagnostics": {
			// 	Type:     schema.TypeList,
			// 	Optional: true,
			// 	MaxItems: 1,
			// 	Elem: &schema.Resource{
			// 		Schema: map[string]*schema.Schema{
			// 			"enabled": {
			// 				Type:     schema.TypeBool,
			// 				Optional: true,
			// 				Default:  true,
			// 			},
			// 			"storage_uri": {
			// 				Type:     schema.TypeString,
			// 				Required: true,
			// 			},
			// 		},
			// 	},
			// },

			"storage_profile_os_disk": {
				Type:     schema.TypeSet,
				Required: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:     schema.TypeString,
							Optional: true,
						},

						"image": {
							Type:     schema.TypeString,
							Optional: true,
						},

						"vhd_containers": {
							Type:     schema.TypeSet,
							Optional: true,
							Elem:     &schema.Schema{Type: schema.TypeString},
							Set:      schema.HashString,
						},

						"managed_disk_type": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
							ValidateFunc: validation.StringInSlice([]string{
								string(compute.PremiumLRS),
								string(compute.StandardLRS),
							}, false),
						},

						"caching": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
						},

						"os_type": {
							Type:     schema.TypeString,
							Optional: true,
						},

						"create_option": {
							Type:     schema.TypeString,
							Required: true,
						},
					},
				},
				Set: resourceArmVirtualMachineScaleSetStorageProfileOsDiskHash,
			},

			"storage_profile_data_disk": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"lun": {
							Type:     schema.TypeInt,
							Required: true,
						},

						"create_option": {
							Type:     schema.TypeString,
							Required: true,
						},

						"caching": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
						},

						"disk_size_gb": {
							Type:         schema.TypeInt,
							Optional:     true,
							Computed:     true,
							ValidateFunc: validateDiskSizeGB,
						},

						"managed_disk_type": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
							ValidateFunc: validation.StringInSlice([]string{
								string(compute.PremiumLRS),
								string(compute.StandardLRS),
							}, false),
						},
					},
				},
			},

			"storage_profile_image_reference": {
				Type:     schema.TypeSet,
				Optional: true,
				Computed: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:     schema.TypeString,
							Optional: true,
						},

						"publisher": {
							Type:     schema.TypeString,
							Optional: true,
						},

						"offer": {
							Type:     schema.TypeString,
							Optional: true,
						},

						"sku": {
							Type:     schema.TypeString,
							Optional: true,
						},

						"version": {
							Type:     schema.TypeString,
							Optional: true,
						},
					},
				},
				Set: resourceArmVirtualMachineScaleSetStorageProfileImageReferenceHash,
			},

			"plan": {
				Type:     schema.TypeSet,
				Optional: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:     schema.TypeString,
							Required: true,
						},

						"publisher": {
							Type:     schema.TypeString,
							Required: true,
						},

						"product": {
							Type:     schema.TypeString,
							Required: true,
						},
					},
				},
			},

			"extension": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:     schema.TypeString,
							Required: true,
						},

						"publisher": {
							Type:     schema.TypeString,
							Required: true,
						},

						"type": {
							Type:     schema.TypeString,
							Required: true,
						},

						"type_handler_version": {
							Type:     schema.TypeString,
							Required: true,
						},

						"auto_upgrade_minor_version": {
							Type:     schema.TypeBool,
							Optional: true,
						},

						"settings": {
							Type:             schema.TypeString,
							Optional:         true,
							ValidateFunc:     validation.ValidateJsonString,
							DiffSuppressFunc: structure.SuppressJsonDiff,
						},

						"protected_settings": {
							Type:             schema.TypeString,
							Optional:         true,
							Sensitive:        true,
							ValidateFunc:     validation.ValidateJsonString,
							DiffSuppressFunc: structure.SuppressJsonDiff,
						},
					},
				},
				Set: resourceArmVirtualMachineScaleSetExtensionHash,
			},

			"tags": tagsSchema(),
		},
	}
}

func resourceArmVirtualMachineScaleSetCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ArmClient).vmScaleSetClient
	ctx := meta.(*ArmClient).StopContext

	log.Printf("[INFO] preparing arguments for Azure ARM Virtual Machine Scale Set creation.")

	name := d.Get("name").(string)
	location := azureStackNormalizeLocation(d.Get("location").(string))
	resGroup := d.Get("resource_group_name").(string)
	tags := d.Get("tags").(map[string]interface{})
	// zones := expandZones(d.Get("zones").([]interface{}))

	sku, err := expandVirtualMachineScaleSetSku(d)
	if err != nil {
		return err
	}

	storageProfile := compute.VirtualMachineScaleSetStorageProfile{}
	osDisk, err := expandAzureStackVirtualMachineScaleSetsStorageProfileOsDisk(d)
	if err != nil {
		return err
	}
	storageProfile.OsDisk = osDisk

	// Not supported
	if _, ok := d.GetOk("storage_profile_data_disk"); ok {
		dataDisks, err := expandAzureStackVirtualMachineScaleSetsStorageProfileDataDisk(d)
		if err != nil {
			return err
		}
		storageProfile.DataDisks = &dataDisks
	}

	if _, ok := d.GetOk("storage_profile_image_reference"); ok {
		imageRef, err := expandAzureRmVirtualMachineScaleSetStorageProfileImageReference(d)
		if err != nil {
			return err
		}
		storageProfile.ImageReference = imageRef
	}

	osProfile, err := expandAzureStackVirtualMachineScaleSetsOsProfile(d)
	if err != nil {
		return err
	}

	extensions, err := expandAzureStackVirtualMachineScaleSetExtensions(d)
	if err != nil {
		return err
	}

	updatePolicy := d.Get("upgrade_policy_mode").(string)
	overprovision := d.Get("overprovision").(bool)

	scaleSetProps := compute.VirtualMachineScaleSetProperties{
		UpgradePolicy: &compute.UpgradePolicy{
			Mode: compute.UpgradeMode(updatePolicy),
		},
		VirtualMachineProfile: &compute.VirtualMachineScaleSetVMProfile{
			NetworkProfile:   expandAzureRmVirtualMachineScaleSetNetworkProfile(d),
			StorageProfile:   &storageProfile,
			OsProfile:        osProfile,
			ExtensionProfile: extensions,
		},
		Overprovision: &overprovision,
	}

	properties := compute.VirtualMachineScaleSet{
		Name:                             &name,
		Location:                         &location,
		Tags:                             *expandTags(tags),
		Sku:                              sku,
		VirtualMachineScaleSetProperties: &scaleSetProps,
	}

	if _, ok := d.GetOk("identity"); ok {
		properties.Identity = expandAzureRmVirtualMachineScaleSetIdentity(d)
	}

	future, err := client.CreateOrUpdate(ctx, resGroup, name, properties)
	if err != nil {
		return err
	}

	if err = future.WaitForCompletionRef(ctx, client.Client); err != nil {
		return fmt.Errorf("Error Creating/Updating Virtual Machine Scale Set %s (resource group %s) ID: %+v", name, resGroup, err)
	}

	read, err := client.Get(ctx, resGroup, name)
	if err != nil {
		return err
	}
	if read.ID == nil {
		return fmt.Errorf("Cannot read Virtual Machine Scale Set %s (resource group %s) ID", name, resGroup)
	}

	d.SetId(*read.ID)

	return resourceArmVirtualMachineScaleSetRead(d, meta)
}

func resourceArmVirtualMachineScaleSetRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ArmClient).vmScaleSetClient
	ctx := meta.(*ArmClient).StopContext

	id, err := parseAzureResourceID(d.Id())
	if err != nil {
		return err
	}
	resGroup := id.ResourceGroup
	name := id.Path["virtualMachineScaleSets"]

	resp, err := client.Get(ctx, resGroup, name)
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			log.Printf("[INFO] AzureStack Virtual Machine Scale Set (%s) Not Found. Removing from State", name)
			d.SetId("")
			return nil
		}
		return fmt.Errorf("Error making Read request on Azure Virtual Machine Scale Set %s: %+v", name, err)
	}

	d.Set("name", resp.Name)
	d.Set("resource_group_name", resGroup)
	if location := resp.Location; location != nil {
		d.Set("location", azureStackNormalizeLocation(*location))
	}
	// d.Set("zones", resp.Zones)

	if err := d.Set("sku", flattenAzureRmVirtualMachineScaleSetSku(resp.Sku)); err != nil {
		return fmt.Errorf("[DEBUG] Error setting `sku`: %#v", err)
	}

	flattenedIdentity := flattenAzureRmVirtualMachineScaleSetIdentity(resp.Identity)
	if err := d.Set("identity", flattenedIdentity); err != nil {
		return fmt.Errorf("[DEBUG] Error setting `identity`: %+v", err)
	}

	if properties := resp.VirtualMachineScaleSetProperties; properties != nil {

		if upgradePolicy := properties.UpgradePolicy; upgradePolicy != nil {
			d.Set("upgrade_policy_mode", upgradePolicy.Mode)
		}
		d.Set("overprovision", properties.Overprovision)

		if profile := properties.VirtualMachineProfile; profile != nil {

			osProfile := flattenAzureStackVirtualMachineScaleSetOsProfile(d, profile.OsProfile)
			if err := d.Set("os_profile", osProfile); err != nil {
				return fmt.Errorf("[DEBUG] Error setting `os_profile`: %#v", err)
			}

			if osProfile := profile.OsProfile; osProfile != nil {
				if linuxConfiguration := osProfile.LinuxConfiguration; linuxConfiguration != nil {
					flattenedLinuxConfiguration := flattenAzureRmVirtualMachineScaleSetOsProfileLinuxConfig(linuxConfiguration)
					if err := d.Set("os_profile_linux_config", flattenedLinuxConfiguration); err != nil {
						return fmt.Errorf("[DEBUG] Error setting `os_profile_linux_config`: %#v", err)
					}
				}

				if secrets := osProfile.Secrets; secrets != nil {
					flattenedSecrets := flattenAzureRmVirtualMachineScaleSetOsProfileSecrets(secrets)
					if err := d.Set("os_profile_secrets", flattenedSecrets); err != nil {
						return fmt.Errorf("[DEBUG] Error setting `os_profile_secrets`: %#v", err)
					}

				}

				if windowsConfiguration := osProfile.WindowsConfiguration; windowsConfiguration != nil {
					flattenedWindowsConfiguration := flattenAzureRmVirtualMachineScaleSetOsProfileWindowsConfig(windowsConfiguration)
					if err := d.Set("os_profile_windows_config", flattenedWindowsConfiguration); err != nil {
						return fmt.Errorf("[DEBUG] Error setting `os_profile_windows_config`: %#v", err)
					}
				}
			}

			if networkProfile := profile.NetworkProfile; networkProfile != nil {
				flattenedNetworkProfile := flattenAzureRmVirtualMachineScaleSetNetworkProfile(networkProfile)
				if err := d.Set("network_profile", flattenedNetworkProfile); err != nil {
					return fmt.Errorf("[DEBUG] Error setting `network_profile`: %#v", err)
				}
			}

			if storageProfile := profile.StorageProfile; storageProfile != nil {

				if dataDisks := resp.VirtualMachineProfile.StorageProfile.DataDisks; dataDisks != nil {
					flattenedDataDisks := flattenAzureStackVirtualMachineScaleSetStorageProfileDataDisk(dataDisks)
					if err := d.Set("storage_profile_data_disk", flattenedDataDisks); err != nil {
						return fmt.Errorf("[DEBUG] Error setting `storage_profile_data_disk`: %#v", err)
					}
				}

				if imageRef := storageProfile.ImageReference; imageRef != nil {
					flattenedImageRef := flattenAzureRmVirtualMachineScaleSetStorageProfileImageReference(imageRef)
					if err := d.Set("storage_profile_image_reference", flattenedImageRef); err != nil {
						return fmt.Errorf("[DEBUG] Error setting `storage_profile_image_reference`: %#v", err)
					}
				}

				if osDisk := storageProfile.OsDisk; osDisk != nil {
					flattenedOSDisk := flattenAzureRmVirtualMachineScaleSetStorageProfileOSDisk(osDisk)
					if err := d.Set("storage_profile_os_disk", flattenedOSDisk); err != nil {
						return fmt.Errorf("[DEBUG] Error setting `storage_profile_os_disk`: %#v", err)
					}
				}
			}

			if extensionProfile := properties.VirtualMachineProfile.ExtensionProfile; extensionProfile != nil {
				extension, err := flattenAzureRmVirtualMachineScaleSetExtensionProfile(extensionProfile)
				if err != nil {
					return fmt.Errorf("[DEBUG] Error setting Virtual Machine Scale Set Extension Profile error: %#v", err)
				}
				if err := d.Set("extension", extension); err != nil {
					return fmt.Errorf("[DEBUG] Error setting `extension`: %#v", err)
				}
			}
		}
	}

	flattenAndSetTags(d, &resp.Tags)

	return nil
}

func resourceArmVirtualMachineScaleSetDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ArmClient).vmScaleSetClient
	ctx := meta.(*ArmClient).StopContext

	id, err := parseAzureResourceID(d.Id())
	if err != nil {
		return err
	}
	resGroup := id.ResourceGroup
	name := id.Path["virtualMachineScaleSets"]

	future, err := client.Delete(ctx, resGroup, name)
	if err != nil {
		return err
	}

	if err = future.WaitForCompletionRef(ctx, client.Client); err != nil {
		return fmt.Errorf("Error Deleting Virtual Machine Scale Set %s (resource group %s) ID: %+v", name, resGroup, err)

	}

	return nil
}

func flattenAzureRmVirtualMachineScaleSetIdentity(identity *compute.VirtualMachineScaleSetIdentity) []interface{} {
	if identity == nil {
		return make([]interface{}, 0)
	}

	result := make(map[string]interface{})
	result["type"] = string(identity.Type)
	if identity.PrincipalID != nil {
		result["principal_id"] = *identity.PrincipalID
	}

	return []interface{}{result}
}

func flattenAzureRmVirtualMachineScaleSetOsProfileLinuxConfig(config *compute.LinuxConfiguration) []interface{} {
	result := make(map[string]interface{})

	if v := config.DisablePasswordAuthentication; v != nil {
		result["disable_password_authentication"] = *v
	}

	if ssh := config.SSH; ssh != nil {
		if keys := ssh.PublicKeys; keys != nil {
			ssh_keys := make([]map[string]interface{}, 0, len(*keys))

			for _, i := range *keys {
				key := make(map[string]interface{})
				if i.Path != nil {
					key["path"] = *i.Path
				}
				if i.KeyData != nil {
					key["key_data"] = *i.KeyData
				}
				ssh_keys = append(ssh_keys, key)
			}

			result["ssh_keys"] = ssh_keys
		}
	}

	return []interface{}{result}
}

func flattenAzureRmVirtualMachineScaleSetOsProfileWindowsConfig(config *compute.WindowsConfiguration) []interface{} {
	result := make(map[string]interface{})

	if config.ProvisionVMAgent != nil {
		result["provision_vm_agent"] = *config.ProvisionVMAgent
	}

	if config.EnableAutomaticUpdates != nil {
		result["enable_automatic_upgrades"] = *config.EnableAutomaticUpdates
	}

	if config.WinRM != nil {
		listeners := make([]map[string]interface{}, 0, len(*config.WinRM.Listeners))
		for _, i := range *config.WinRM.Listeners {
			listener := make(map[string]interface{})
			listener["protocol"] = i.Protocol

			if i.CertificateURL != nil {
				listener["certificate_url"] = *i.CertificateURL
			}

			listeners = append(listeners, listener)
		}

		result["winrm"] = listeners
	}

	if config.AdditionalUnattendContent != nil {
		content := make([]map[string]interface{}, 0, len(*config.AdditionalUnattendContent))
		for _, i := range *config.AdditionalUnattendContent {
			c := make(map[string]interface{})
			c["pass"] = i.PassName
			c["component"] = i.ComponentName
			c["setting_name"] = i.SettingName

			if i.Content != nil {
				c["content"] = *i.Content
			}

			content = append(content, c)
		}

		result["additional_unattend_config"] = content
	}

	return []interface{}{result}
}

func flattenAzureRmVirtualMachineScaleSetOsProfileSecrets(secrets *[]compute.VaultSecretGroup) []map[string]interface{} {
	result := make([]map[string]interface{}, 0, len(*secrets))
	for _, secret := range *secrets {
		s := map[string]interface{}{
			"source_vault_id": *secret.SourceVault.ID,
		}

		if secret.VaultCertificates != nil {
			certs := make([]map[string]interface{}, 0, len(*secret.VaultCertificates))
			for _, cert := range *secret.VaultCertificates {
				vaultCert := make(map[string]interface{})
				vaultCert["certificate_url"] = *cert.CertificateURL

				if cert.CertificateStore != nil {
					vaultCert["certificate_store"] = *cert.CertificateStore
				}

				certs = append(certs, vaultCert)
			}

			s["vault_certificates"] = certs
		}

		result = append(result, s)
	}
	return result
}

func flattenAzureRmVirtualMachineScaleSetNetworkProfile(profile *compute.VirtualMachineScaleSetNetworkProfile) []map[string]interface{} {
	networkConfigurations := profile.NetworkInterfaceConfigurations
	result := make([]map[string]interface{}, 0, len(*networkConfigurations))
	for _, netConfig := range *networkConfigurations {
		s := map[string]interface{}{
			"name":    *netConfig.Name,
			"primary": *netConfig.VirtualMachineScaleSetNetworkConfigurationProperties.Primary,
		}

		if netConfig.VirtualMachineScaleSetNetworkConfigurationProperties.IPConfigurations != nil {
			ipConfigs := make([]map[string]interface{}, 0, len(*netConfig.VirtualMachineScaleSetNetworkConfigurationProperties.IPConfigurations))
			for _, ipConfig := range *netConfig.VirtualMachineScaleSetNetworkConfigurationProperties.IPConfigurations {
				config := make(map[string]interface{})
				config["name"] = *ipConfig.Name

				properties := ipConfig.VirtualMachineScaleSetIPConfigurationProperties

				if ipConfig.VirtualMachineScaleSetIPConfigurationProperties.Subnet != nil {
					config["subnet_id"] = *properties.Subnet.ID
				}

				addressPools := make([]interface{}, 0)
				if properties.ApplicationGatewayBackendAddressPools != nil {
					for _, pool := range *properties.ApplicationGatewayBackendAddressPools {
						addressPools = append(addressPools, *pool.ID)
					}
				}
				config["application_gateway_backend_address_pool_ids"] = schema.NewSet(schema.HashString, addressPools)

				if properties.LoadBalancerBackendAddressPools != nil {
					addressPools := make([]interface{}, 0, len(*properties.LoadBalancerBackendAddressPools))
					for _, pool := range *properties.LoadBalancerBackendAddressPools {
						addressPools = append(addressPools, *pool.ID)
					}
					config["load_balancer_backend_address_pool_ids"] = schema.NewSet(schema.HashString, addressPools)
				}

				if properties.LoadBalancerInboundNatPools != nil {
					inboundNatPools := make([]interface{}, 0, len(*properties.LoadBalancerInboundNatPools))
					for _, rule := range *properties.LoadBalancerInboundNatPools {
						inboundNatPools = append(inboundNatPools, *rule.ID)
					}
					config["load_balancer_inbound_nat_rules_ids"] = schema.NewSet(schema.HashString, inboundNatPools)
				}

				ipConfigs = append(ipConfigs, config)
			}

			s["ip_configuration"] = ipConfigs
		}

		result = append(result, s)
	}

	return result
}

func flattenAzureStackVirtualMachineScaleSetOsProfile(d *schema.ResourceData, profile *compute.VirtualMachineScaleSetOSProfile) []interface{} {
	result := make(map[string]interface{})

	result["computer_name_prefix"] = *profile.ComputerNamePrefix
	result["admin_username"] = *profile.AdminUsername

	// admin password isn't returned, so let's look it up
	if v, ok := d.GetOk("os_profile.0.admin_password"); ok {
		password := v.(string)
		result["admin_password"] = password
	}

	if profile.CustomData != nil {
		result["custom_data"] = *profile.CustomData
	} else {
		// look up the current custom data
		value := d.Get("os_profile.0.custom_data").(string)
		if !isBase64Encoded(value) {
			value = base64Encode(value)
		}
		result["custom_data"] = value
	}

	return []interface{}{result}
}

func flattenAzureRmVirtualMachineScaleSetStorageProfileOSDisk(profile *compute.VirtualMachineScaleSetOSDisk) []interface{} {
	result := make(map[string]interface{})

	if profile.Name != nil {
		result["name"] = *profile.Name
	}

	if profile.Image != nil {
		result["image"] = *profile.Image.URI
	}

	if profile.VhdContainers != nil {
		containers := make([]interface{}, 0, len(*profile.VhdContainers))
		for _, container := range *profile.VhdContainers {
			containers = append(containers, container)
		}
		result["vhd_containers"] = schema.NewSet(schema.HashString, containers)
	}

	if profile.ManagedDisk != nil {
		result["managed_disk_type"] = string(profile.ManagedDisk.StorageAccountType)
	}

	result["caching"] = profile.Caching
	result["create_option"] = profile.CreateOption
	result["os_type"] = profile.OsType

	return []interface{}{result}
}

func flattenAzureStackVirtualMachineScaleSetStorageProfileDataDisk(disks *[]compute.VirtualMachineScaleSetDataDisk) interface{} {
	result := make([]interface{}, len(*disks))
	for i, disk := range *disks {
		l := make(map[string]interface{})
		if disk.ManagedDisk != nil {
			l["managed_disk_type"] = string(disk.ManagedDisk.StorageAccountType)
		}

		l["create_option"] = disk.CreateOption
		l["caching"] = string(disk.Caching)
		if disk.DiskSizeGB != nil {
			l["disk_size_gb"] = *disk.DiskSizeGB
		}
		if v := disk.Lun; v != nil {
			l["lun"] = *v
		}

		result[i] = l
	}
	return result
}

func flattenAzureRmVirtualMachineScaleSetStorageProfileImageReference(image *compute.ImageReference) []interface{} {
	result := make(map[string]interface{})
	if image.Publisher != nil {
		result["publisher"] = *image.Publisher
	}
	if image.Offer != nil {
		result["offer"] = *image.Offer
	}
	if image.Sku != nil {
		result["sku"] = *image.Sku
	}
	if image.Version != nil {
		result["version"] = *image.Version
	}

	return []interface{}{result}
}

func flattenAzureRmVirtualMachineScaleSetSku(sku *compute.Sku) []interface{} {
	result := make(map[string]interface{})
	result["name"] = *sku.Name
	result["capacity"] = *sku.Capacity

	if *sku.Tier != "" {
		result["tier"] = *sku.Tier
	}

	return []interface{}{result}
}

func flattenAzureRmVirtualMachineScaleSetExtensionProfile(profile *compute.VirtualMachineScaleSetExtensionProfile) ([]map[string]interface{}, error) {
	if profile.Extensions == nil {
		return nil, nil
	}

	result := make([]map[string]interface{}, 0, len(*profile.Extensions))
	for _, extension := range *profile.Extensions {
		e := make(map[string]interface{})
		e["name"] = *extension.Name
		properties := extension.VirtualMachineScaleSetExtensionProperties
		if properties != nil {
			e["publisher"] = *properties.Publisher
			e["type"] = *properties.Type
			e["type_handler_version"] = *properties.TypeHandlerVersion
			if properties.AutoUpgradeMinorVersion != nil {
				e["auto_upgrade_minor_version"] = *properties.AutoUpgradeMinorVersion
			}

			if settings := properties.Settings; settings != nil {
				settingsVal := settings.(map[string]interface{})
				settingsJson, err := structure.FlattenJsonToString(settingsVal)
				if err != nil {
					return nil, err
				}
				e["settings"] = settingsJson
			}
		}

		result = append(result, e)
	}

	return result, nil
}

func resourceArmVirtualMachineScaleSetStorageProfileImageReferenceHash(v interface{}) int {
	var buf bytes.Buffer

	if m, ok := v.(map[string]interface{}); ok {
		if v, ok := m["publisher"]; ok {
			buf.WriteString(fmt.Sprintf("%s-", v.(string)))
		}
		if v, ok := m["offer"]; ok {
			buf.WriteString(fmt.Sprintf("%s-", v.(string)))
		}
		if v, ok := m["sku"]; ok {
			buf.WriteString(fmt.Sprintf("%s-", v.(string)))
		}
		if v, ok := m["version"]; ok {
			buf.WriteString(fmt.Sprintf("%s-", v.(string)))
		}
		if v, ok := m["id"]; ok {
			buf.WriteString(fmt.Sprintf("%s-", v.(string)))
		}
	}

	return hashcode.String(buf.String())
}

func resourceArmVirtualMachineScaleSetSkuHash(v interface{}) int {
	var buf bytes.Buffer

	if m, ok := v.(map[string]interface{}); ok {
		buf.WriteString(fmt.Sprintf("%s-", m["name"].(string)))
		buf.WriteString(fmt.Sprintf("%d-", m["capacity"].(int)))

		if v, ok := m["tier"]; ok {
			buf.WriteString(fmt.Sprintf("%s-", strings.ToLower(v.(string))))
		}
	}

	return hashcode.String(buf.String())
}

func resourceArmVirtualMachineScaleSetStorageProfileOsDiskHash(v interface{}) int {
	var buf bytes.Buffer

	if m, ok := v.(map[string]interface{}); ok {
		buf.WriteString(fmt.Sprintf("%s-", m["name"].(string)))

		if v, ok := m["vhd_containers"]; ok {
			buf.WriteString(fmt.Sprintf("%s-", v.(*schema.Set).List()))
		}
	}

	return hashcode.String(buf.String())
}

func resourceArmVirtualMachineScaleSetNetworkConfigurationHash(v interface{}) int {
	var buf bytes.Buffer

	if m, ok := v.(map[string]interface{}); ok {
		buf.WriteString(fmt.Sprintf("%s-", m["name"].(string)))
		// buf.WriteString(fmt.Sprintf("%t-", m["primary"].(bool)))
	}

	return hashcode.String(buf.String())
}

func resourceArmVirtualMachineScaleSetOsProfileLinuxConfigHash(v interface{}) int {
	var buf bytes.Buffer

	if m, ok := v.(map[string]interface{}); ok {
		buf.WriteString(fmt.Sprintf("%t-", m["disable_password_authentication"].(bool)))
	}

	return hashcode.String(buf.String())
}

func resourceArmVirtualMachineScaleSetOsProfileWindowsConfigHash(v interface{}) int {
	var buf bytes.Buffer

	if m, ok := v.(map[string]interface{}); ok {
		if v, ok := m["provision_vm_agent"]; ok {
			buf.WriteString(fmt.Sprintf("%t-", v.(bool)))
		}
		if v, ok := m["enable_automatic_upgrades"]; ok {
			buf.WriteString(fmt.Sprintf("%t-", v.(bool)))
		}
	}

	return hashcode.String(buf.String())
}

func resourceArmVirtualMachineScaleSetExtensionHash(v interface{}) int {
	var buf bytes.Buffer

	if m, ok := v.(map[string]interface{}); ok {
		buf.WriteString(fmt.Sprintf("%s-", m["name"].(string)))
		buf.WriteString(fmt.Sprintf("%s-", m["publisher"].(string)))
		buf.WriteString(fmt.Sprintf("%s-", m["type"].(string)))
		buf.WriteString(fmt.Sprintf("%s-", m["type_handler_version"].(string)))

		if v, ok := m["auto_upgrade_minor_version"]; ok {
			buf.WriteString(fmt.Sprintf("%t-", v.(bool)))
		}

		// we need to ensure the whitespace is consistent
		settings := m["settings"].(string)
		if settings != "" {
			expandedSettings, err := structure.ExpandJsonFromString(settings)
			if err == nil {
				serialisedSettings, err := structure.FlattenJsonToString(expandedSettings)
				if err == nil {
					buf.WriteString(fmt.Sprintf("%s-", serialisedSettings))
				}
			}
		}
	}

	return hashcode.String(buf.String())
}

func expandVirtualMachineScaleSetSku(d *schema.ResourceData) (*compute.Sku, error) {
	skuConfig := d.Get("sku").(*schema.Set).List()

	config := skuConfig[0].(map[string]interface{})

	name := config["name"].(string)
	tier := config["tier"].(string)
	capacity := int64(config["capacity"].(int))

	sku := &compute.Sku{
		Name:     &name,
		Capacity: &capacity,
	}

	if tier != "" {
		sku.Tier = &tier
	}

	return sku, nil
}

func expandAzureRmVirtualMachineScaleSetNetworkProfile(d *schema.ResourceData) *compute.VirtualMachineScaleSetNetworkProfile {
	scaleSetNetworkProfileConfigs := d.Get("network_profile").(*schema.Set).List()
	networkProfileConfig := make([]compute.VirtualMachineScaleSetNetworkConfiguration, 0, len(scaleSetNetworkProfileConfigs))

	for _, npProfileConfig := range scaleSetNetworkProfileConfigs {
		config := npProfileConfig.(map[string]interface{})

		name := config["name"].(string)
		primary := config["primary"].(bool)

		ipConfigurationConfigs := config["ip_configuration"].([]interface{})
		ipConfigurations := make([]compute.VirtualMachineScaleSetIPConfiguration, 0, len(ipConfigurationConfigs))
		for _, ipConfigConfig := range ipConfigurationConfigs {
			ipconfig := ipConfigConfig.(map[string]interface{})
			name := ipconfig["name"].(string)
			subnetId := ipconfig["subnet_id"].(string)

			ipConfiguration := compute.VirtualMachineScaleSetIPConfiguration{
				Name: &name,
				VirtualMachineScaleSetIPConfigurationProperties: &compute.VirtualMachineScaleSetIPConfigurationProperties{
					Subnet: &compute.APIEntityReference{
						ID: &subnetId,
					},
				},
			}

			if v := ipconfig["application_gateway_backend_address_pool_ids"]; v != nil {
				pools := v.(*schema.Set).List()
				resources := make([]compute.SubResource, 0, len(pools))
				for _, p := range pools {
					id := p.(string)
					resources = append(resources, compute.SubResource{
						ID: &id,
					})
				}
				ipConfiguration.ApplicationGatewayBackendAddressPools = &resources
			}

			if v := ipconfig["load_balancer_backend_address_pool_ids"]; v != nil {
				pools := v.(*schema.Set).List()
				resources := make([]compute.SubResource, 0, len(pools))
				for _, p := range pools {
					id := p.(string)
					resources = append(resources, compute.SubResource{
						ID: &id,
					})
				}
				ipConfiguration.LoadBalancerBackendAddressPools = &resources
			}

			if v := ipconfig["load_balancer_inbound_nat_rules_ids"]; v != nil {
				rules := v.(*schema.Set).List()
				rulesResources := make([]compute.SubResource, 0, len(rules))
				for _, m := range rules {
					id := m.(string)
					rulesResources = append(rulesResources, compute.SubResource{
						ID: &id,
					})
				}
				ipConfiguration.LoadBalancerInboundNatPools = &rulesResources
			}

			ipConfigurations = append(ipConfigurations, ipConfiguration)
		}

		nProfile := compute.VirtualMachineScaleSetNetworkConfiguration{
			Name: &name,
			VirtualMachineScaleSetNetworkConfigurationProperties: &compute.VirtualMachineScaleSetNetworkConfigurationProperties{
				Primary:          &primary,
				IPConfigurations: &ipConfigurations,
			},
		}

		networkProfileConfig = append(networkProfileConfig, nProfile)
	}

	return &compute.VirtualMachineScaleSetNetworkProfile{
		NetworkInterfaceConfigurations: &networkProfileConfig,
	}
}

func expandAzureStackVirtualMachineScaleSetsOsProfile(d *schema.ResourceData) (*compute.VirtualMachineScaleSetOSProfile, error) {
	osProfileConfigs := d.Get("os_profile").([]interface{})

	osProfileConfig := osProfileConfigs[0].(map[string]interface{})
	namePrefix := osProfileConfig["computer_name_prefix"].(string)
	username := osProfileConfig["admin_username"].(string)
	password := osProfileConfig["admin_password"].(string)
	customData := osProfileConfig["custom_data"].(string)

	osProfile := &compute.VirtualMachineScaleSetOSProfile{
		ComputerNamePrefix: &namePrefix,
		AdminUsername:      &username,
	}

	if password != "" {
		osProfile.AdminPassword = &password
	}

	if customData != "" {
		customData = base64Encode(customData)
		osProfile.CustomData = &customData
	}

	if _, ok := d.GetOk("os_profile_secrets"); ok {
		secrets := expandAzureRmVirtualMachineScaleSetOsProfileSecrets(d)
		if secrets != nil {
			osProfile.Secrets = secrets
		}
	}

	if _, ok := d.GetOk("os_profile_linux_config"); ok {
		linuxConfig, err := expandAzureRmVirtualMachineScaleSetOsProfileLinuxConfig(d)
		if err != nil {
			return nil, err
		}
		osProfile.LinuxConfiguration = linuxConfig
	}

	if _, ok := d.GetOk("os_profile_windows_config"); ok {
		winConfig, err := expandAzureRmVirtualMachineScaleSetOsProfileWindowsConfig(d)
		if err != nil {
			return nil, err
		}
		if winConfig != nil {
			osProfile.WindowsConfiguration = winConfig
		}
	}

	return osProfile, nil
}

func expandAzureRmVirtualMachineScaleSetIdentity(d *schema.ResourceData) *compute.VirtualMachineScaleSetIdentity {
	v := d.Get("identity")
	identities := v.([]interface{})
	identity := identities[0].(map[string]interface{})
	identityType := identity["type"].(string)
	return &compute.VirtualMachineScaleSetIdentity{
		Type: compute.ResourceIdentityType(identityType),
	}
}

func expandAzureStackVirtualMachineScaleSetsStorageProfileOsDisk(d *schema.ResourceData) (*compute.VirtualMachineScaleSetOSDisk, error) {
	osDiskConfigs := d.Get("storage_profile_os_disk").(*schema.Set).List()

	osDiskConfig := osDiskConfigs[0].(map[string]interface{})
	name := osDiskConfig["name"].(string)
	image := osDiskConfig["image"].(string)
	vhd_containers := osDiskConfig["vhd_containers"].(*schema.Set).List()
	caching := osDiskConfig["caching"].(string)
	osType := osDiskConfig["os_type"].(string)
	createOption := osDiskConfig["create_option"].(string)
	managedDiskType := osDiskConfig["managed_disk_type"].(string)

	if managedDiskType == "" && name == "" {
		return nil, fmt.Errorf("[ERROR] `name` must be set in `storage_profile_os_disk` for unmanaged disk")
	}

	osDisk := &compute.VirtualMachineScaleSetOSDisk{
		Name:         &name,
		Caching:      compute.CachingTypes(caching),
		OsType:       compute.OperatingSystemTypes(osType),
		CreateOption: compute.DiskCreateOptionTypes(createOption),
	}

	if image != "" {
		osDisk.Image = &compute.VirtualHardDisk{
			URI: &image,
		}
	}

	if len(vhd_containers) > 0 {
		var vhdContainers []string
		for _, v := range vhd_containers {
			str := v.(string)
			vhdContainers = append(vhdContainers, str)
		}
		osDisk.VhdContainers = &vhdContainers
	}

	managedDisk := &compute.VirtualMachineScaleSetManagedDiskParameters{}

	if managedDiskType != "" {
		if name != "" {
			return nil, fmt.Errorf("[ERROR] Conflict between `name` and `managed_disk_type` on `storage_profile_os_disk` (please remove name or set it to blank)")
		}

		osDisk.Name = nil
		managedDisk.StorageAccountType = compute.StorageAccountTypes(managedDiskType)
		osDisk.ManagedDisk = managedDisk
	}

	//BEGIN: code to be removed after GH-13016 is merged
	if image != "" && managedDiskType != "" {
		return nil, fmt.Errorf("[ERROR] Conflict between `image` and `managed_disk_type` on `storage_profile_os_disk` (only one or the other can be used)")
	}

	if len(vhd_containers) > 0 && managedDiskType != "" {
		return nil, fmt.Errorf("[ERROR] Conflict between `vhd_containers` and `managed_disk_type` on `storage_profile_os_disk` (only one or the other can be used)")
	}
	//END: code to be removed after GH-13016 is merged

	return osDisk, nil
}

func expandAzureStackVirtualMachineScaleSetsStorageProfileDataDisk(d *schema.ResourceData) ([]compute.VirtualMachineScaleSetDataDisk, error) {
	disks := d.Get("storage_profile_data_disk").([]interface{})
	dataDisks := make([]compute.VirtualMachineScaleSetDataDisk, 0, len(disks))
	for _, diskConfig := range disks {
		config := diskConfig.(map[string]interface{})

		createOption := config["create_option"].(string)
		managedDiskType := config["managed_disk_type"].(string)
		lun := int32(config["lun"].(int))

		dataDisk := compute.VirtualMachineScaleSetDataDisk{
			Lun:          &lun,
			CreateOption: compute.DiskCreateOptionTypes(createOption),
		}

		managedDiskVMSS := &compute.VirtualMachineScaleSetManagedDiskParameters{}

		if managedDiskType != "" {
			managedDiskVMSS.StorageAccountType = compute.StorageAccountTypes(managedDiskType)
		} else {
			managedDiskVMSS.StorageAccountType = compute.StandardLRS
		}

		// assume that data disks in VMSS can only be Managed Disks
		dataDisk.ManagedDisk = managedDiskVMSS
		if v := config["caching"].(string); v != "" {
			dataDisk.Caching = compute.CachingTypes(v)
		}

		if v := config["disk_size_gb"]; v != nil {
			diskSize := int32(config["disk_size_gb"].(int))
			dataDisk.DiskSizeGB = &diskSize
		}

		dataDisks = append(dataDisks, dataDisk)
	}

	return dataDisks, nil
}

func expandAzureRmVirtualMachineScaleSetStorageProfileImageReference(d *schema.ResourceData) (*compute.ImageReference, error) {
	storageImageRefs := d.Get("storage_profile_image_reference").(*schema.Set).List()

	storageImageRef := storageImageRefs[0].(map[string]interface{})

	imageID := storageImageRef["id"].(string)
	publisher := storageImageRef["publisher"].(string)

	imageReference := compute.ImageReference{}

	if imageID != "" && publisher != "" {
		return nil, fmt.Errorf("[ERROR] Conflict between `id` and `publisher` (only one or the other can be used)")
	}

	offer := storageImageRef["offer"].(string)
	sku := storageImageRef["sku"].(string)
	version := storageImageRef["version"].(string)

	imageReference.Publisher = utils.String(publisher)
	imageReference.Offer = utils.String(offer)
	imageReference.Sku = utils.String(sku)
	imageReference.Version = utils.String(version)

	return &imageReference, nil
}

func expandAzureRmVirtualMachineScaleSetOsProfileLinuxConfig(d *schema.ResourceData) (*compute.LinuxConfiguration, error) {
	osProfilesLinuxConfig := d.Get("os_profile_linux_config").(*schema.Set).List()

	linuxConfig := osProfilesLinuxConfig[0].(map[string]interface{})
	disablePasswordAuth := linuxConfig["disable_password_authentication"].(bool)

	linuxKeys := linuxConfig["ssh_keys"].([]interface{})
	sshPublicKeys := make([]compute.SSHPublicKey, 0, len(linuxKeys))
	for _, key := range linuxKeys {
		if key == nil {
			continue
		}
		sshKey := key.(map[string]interface{})
		path := sshKey["path"].(string)
		keyData := sshKey["key_data"].(string)

		sshPublicKey := compute.SSHPublicKey{
			Path:    &path,
			KeyData: &keyData,
		}

		sshPublicKeys = append(sshPublicKeys, sshPublicKey)
	}

	config := &compute.LinuxConfiguration{
		DisablePasswordAuthentication: &disablePasswordAuth,
		SSH: &compute.SSHConfiguration{
			PublicKeys: &sshPublicKeys,
		},
	}

	return config, nil
}

func expandAzureRmVirtualMachineScaleSetOsProfileWindowsConfig(d *schema.ResourceData) (*compute.WindowsConfiguration, error) {
	osProfilesWindowsConfig := d.Get("os_profile_windows_config").(*schema.Set).List()

	osProfileConfig := osProfilesWindowsConfig[0].(map[string]interface{})
	config := &compute.WindowsConfiguration{}

	if v := osProfileConfig["provision_vm_agent"]; v != nil {
		provision := v.(bool)
		config.ProvisionVMAgent = &provision
	}

	if v := osProfileConfig["enable_automatic_upgrades"]; v != nil {
		update := v.(bool)
		config.EnableAutomaticUpdates = &update
	}

	if v := osProfileConfig["winrm"]; v != nil {
		winRm := v.([]interface{})
		if len(winRm) > 0 {
			winRmListeners := make([]compute.WinRMListener, 0, len(winRm))
			for _, winRmConfig := range winRm {
				config := winRmConfig.(map[string]interface{})

				protocol := config["protocol"].(string)
				winRmListener := compute.WinRMListener{
					Protocol: compute.ProtocolTypes(protocol),
				}
				if v := config["certificate_url"].(string); v != "" {
					winRmListener.CertificateURL = &v
				}

				winRmListeners = append(winRmListeners, winRmListener)
			}
			config.WinRM = &compute.WinRMConfiguration{
				Listeners: &winRmListeners,
			}
		}
	}
	if v := osProfileConfig["additional_unattend_config"]; v != nil {
		additionalConfig := v.([]interface{})
		if len(additionalConfig) > 0 {
			additionalConfigContent := make([]compute.AdditionalUnattendContent, 0, len(additionalConfig))
			for _, addConfig := range additionalConfig {
				config := addConfig.(map[string]interface{})
				pass := config["pass"].(string)
				component := config["component"].(string)
				settingName := config["setting_name"].(string)
				content := config["content"].(string)

				addContent := compute.AdditionalUnattendContent{
					PassName:      compute.PassNames(pass),
					ComponentName: compute.ComponentNames(component),
					SettingName:   compute.SettingNames(settingName),
				}

				if content != "" {
					addContent.Content = &content
				}

				additionalConfigContent = append(additionalConfigContent, addContent)
			}
			config.AdditionalUnattendContent = &additionalConfigContent
		}
	}
	return config, nil
}

func expandAzureRmVirtualMachineScaleSetOsProfileSecrets(d *schema.ResourceData) *[]compute.VaultSecretGroup {
	secretsConfig := d.Get("os_profile_secrets").(*schema.Set).List()
	secrets := make([]compute.VaultSecretGroup, 0, len(secretsConfig))

	for _, secretConfig := range secretsConfig {
		config := secretConfig.(map[string]interface{})
		sourceVaultId := config["source_vault_id"].(string)

		vaultSecretGroup := compute.VaultSecretGroup{
			SourceVault: &compute.SubResource{
				ID: &sourceVaultId,
			},
		}

		if v := config["vault_certificates"]; v != nil {
			certsConfig := v.([]interface{})
			certs := make([]compute.VaultCertificate, 0, len(certsConfig))
			for _, certConfig := range certsConfig {
				config := certConfig.(map[string]interface{})

				certUrl := config["certificate_url"].(string)
				cert := compute.VaultCertificate{
					CertificateURL: &certUrl,
				}
				if v := config["certificate_store"].(string); v != "" {
					cert.CertificateStore = &v
				}

				certs = append(certs, cert)
			}
			vaultSecretGroup.VaultCertificates = &certs
		}

		secrets = append(secrets, vaultSecretGroup)
	}

	return &secrets
}

func expandAzureStackVirtualMachineScaleSetExtensions(d *schema.ResourceData) (*compute.VirtualMachineScaleSetExtensionProfile, error) {
	extensions := d.Get("extension").(*schema.Set).List()
	resources := make([]compute.VirtualMachineScaleSetExtension, 0, len(extensions))
	for _, e := range extensions {
		config := e.(map[string]interface{})
		name := config["name"].(string)
		publisher := config["publisher"].(string)
		t := config["type"].(string)
		version := config["type_handler_version"].(string)

		extension := compute.VirtualMachineScaleSetExtension{
			Name: &name,
			VirtualMachineScaleSetExtensionProperties: &compute.VirtualMachineScaleSetExtensionProperties{
				Publisher:          &publisher,
				Type:               &t,
				TypeHandlerVersion: &version,
			},
		}

		if u := config["auto_upgrade_minor_version"]; u != nil {
			upgrade := u.(bool)
			extension.VirtualMachineScaleSetExtensionProperties.AutoUpgradeMinorVersion = &upgrade
		}

		if s := config["settings"].(string); s != "" {
			settings, err := structure.ExpandJsonFromString(s)
			if err != nil {
				return nil, fmt.Errorf("unable to parse settings: %+v", err)
			}
			extension.VirtualMachineScaleSetExtensionProperties.Settings = &settings
		}

		if s := config["protected_settings"].(string); s != "" {
			protectedSettings, err := structure.ExpandJsonFromString(s)
			if err != nil {
				return nil, fmt.Errorf("unable to parse protected_settings: %+v", err)
			}
			extension.VirtualMachineScaleSetExtensionProperties.ProtectedSettings = &protectedSettings
		}

		resources = append(resources, extension)
	}

	return &compute.VirtualMachineScaleSetExtensionProfile{
		Extensions: &resources,
	}, nil
}
