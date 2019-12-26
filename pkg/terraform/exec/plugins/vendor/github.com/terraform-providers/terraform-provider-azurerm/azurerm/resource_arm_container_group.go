package azurerm

import (
	"bytes"
	"fmt"
	"log"
	"strings"

	"github.com/hashicorp/terraform/helper/hashcode"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/suppress"

	"github.com/Azure/azure-sdk-for-go/services/containerinstance/mgmt/2018-10-01/containerinstance"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/helper/validation"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/tf"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/validate"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func resourceArmContainerGroup() *schema.Resource {
	return &schema.Resource{
		Create: resourceArmContainerGroupCreate,
		Read:   resourceArmContainerGroupRead,
		Delete: resourceArmContainerGroupDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validate.NoEmptyStrings,
			},

			"location": locationSchema(),

			"resource_group_name": resourceGroupNameSchema(),

			"ip_address_type": {
				Type:             schema.TypeString,
				Optional:         true,
				Default:          "Public",
				ForceNew:         true,
				DiffSuppressFunc: suppress.CaseDifference,
				ValidateFunc: validation.StringInSlice([]string{
					string(containerinstance.Public),
				}, true),
			},

			"os_type": {
				Type:             schema.TypeString,
				Required:         true,
				ForceNew:         true,
				DiffSuppressFunc: suppress.CaseDifference,
				ValidateFunc: validation.StringInSlice([]string{
					string(containerinstance.Windows),
					string(containerinstance.Linux),
				}, true),
			},

			"image_registry_credential": {
				Type:     schema.TypeList,
				Optional: true,
				ForceNew: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"server": {
							Type:         schema.TypeString,
							Required:     true,
							ForceNew:     true,
							ValidateFunc: validate.NoEmptyStrings,
						},

						"username": {
							Type:         schema.TypeString,
							Required:     true,
							ForceNew:     true,
							ValidateFunc: validate.NoEmptyStrings,
						},

						"password": {
							Type:         schema.TypeString,
							Required:     true,
							Sensitive:    true,
							ForceNew:     true,
							ValidateFunc: validate.NoEmptyStrings,
						},
					},
				},
			},

			"identity": {
				Type:     schema.TypeList,
				Optional: true,
				Computed: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"type": {
							Type:     schema.TypeString,
							Required: true,
							ValidateFunc: validation.StringInSlice([]string{
								"SystemAssigned",
								"UserAssigned",
								"SystemAssigned, UserAssigned",
							}, false),
						},
						"principal_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"identity_ids": {
							Type:     schema.TypeList,
							Optional: true,
							MinItems: 1,
							Elem: &schema.Schema{
								Type:         schema.TypeString,
								ValidateFunc: validation.NoZeroValues,
							},
						},
					},
				},
			},

			"tags": tagsForceNewSchema(),

			"restart_policy": {
				Type:             schema.TypeString,
				Optional:         true,
				ForceNew:         true,
				Default:          string(containerinstance.Always),
				DiffSuppressFunc: suppress.CaseDifference,
				ValidateFunc: validation.StringInSlice([]string{
					string(containerinstance.Always),
					string(containerinstance.Never),
					string(containerinstance.OnFailure),
				}, true),
			},

			"dns_name_label": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},

			"container": {
				Type:     schema.TypeList,
				Required: true,
				ForceNew: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:         schema.TypeString,
							Required:     true,
							ForceNew:     true,
							ValidateFunc: validate.NoEmptyStrings,
						},

						"image": {
							Type:         schema.TypeString,
							Required:     true,
							ForceNew:     true,
							ValidateFunc: validate.NoEmptyStrings,
						},

						"cpu": {
							Type:     schema.TypeFloat,
							Required: true,
							ForceNew: true,
						},

						"memory": {
							Type:     schema.TypeFloat,
							Required: true,
							ForceNew: true,
						},

						"gpu": {
							Type:     schema.TypeList,
							Optional: true,
							MaxItems: 1,
							ForceNew: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"count": {
										Type:     schema.TypeInt,
										Optional: true,
										ForceNew: true,
										ValidateFunc: validate.IntInSlice([]int{
											1,
											2,
											4,
										}),
									},

									"sku": {
										Type:     schema.TypeString,
										Optional: true,
										ForceNew: true,
										ValidateFunc: validation.StringInSlice([]string{
											"K80",
											"P100",
											"V100",
										}, false),
									},
								},
							},
						},

						"port": {
							Type:         schema.TypeInt,
							Optional:     true,
							ForceNew:     true,
							Computed:     true,
							Deprecated:   "Deprecated in favor of `ports`",
							ValidateFunc: validate.PortNumber,
						},

						"protocol": {
							Type:             schema.TypeString,
							Optional:         true,
							ForceNew:         true,
							Computed:         true,
							Deprecated:       "Deprecated in favor of `ports`",
							DiffSuppressFunc: suppress.CaseDifference,
							ValidateFunc: validation.StringInSlice([]string{
								string(containerinstance.TCP),
								string(containerinstance.UDP),
							}, true),
						},

						"ports": {
							Type:     schema.TypeSet,
							Optional: true,
							ForceNew: true,
							Computed: true,
							Set:      resourceArmContainerGroupPortsHash,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"port": {
										Type:         schema.TypeInt,
										Optional:     true,
										ForceNew:     true,
										Computed:     true,
										ValidateFunc: validate.PortNumber,
									},

									"protocol": {
										Type:     schema.TypeString,
										Optional: true,
										ForceNew: true,
										Computed: true,
										//Default:  string(containerinstance.TCP), restore in 2.0
										ValidateFunc: validation.StringInSlice([]string{
											string(containerinstance.TCP),
											string(containerinstance.UDP),
										}, false),
									},
								},
							},
						},

						"environment_variables": {
							Type:     schema.TypeMap,
							ForceNew: true,
							Optional: true,
						},

						"secure_environment_variables": {
							Type:      schema.TypeMap,
							Optional:  true,
							ForceNew:  true,
							Sensitive: true,
						},

						"command": {
							Type:       schema.TypeString,
							Optional:   true,
							Computed:   true,
							Deprecated: "Use `commands` instead.",
						},

						"commands": {
							Type:     schema.TypeList,
							Optional: true,
							Computed: true,
							Elem:     &schema.Schema{Type: schema.TypeString},
						},

						"volume": {
							Type:     schema.TypeList,
							Optional: true,
							ForceNew: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"name": {
										Type:         schema.TypeString,
										Required:     true,
										ForceNew:     true,
										ValidateFunc: validate.NoEmptyStrings,
									},

									"mount_path": {
										Type:         schema.TypeString,
										Required:     true,
										ForceNew:     true,
										ValidateFunc: validate.NoEmptyStrings,
									},

									"read_only": {
										Type:     schema.TypeBool,
										Optional: true,
										ForceNew: true,
										Default:  false,
									},

									"share_name": {
										Type:         schema.TypeString,
										Required:     true,
										ForceNew:     true,
										ValidateFunc: validate.NoEmptyStrings,
									},

									"storage_account_name": {
										Type:         schema.TypeString,
										Required:     true,
										ForceNew:     true,
										ValidateFunc: validate.NoEmptyStrings,
									},

									"storage_account_key": {
										Type:         schema.TypeString,
										Required:     true,
										ForceNew:     true,
										ValidateFunc: validate.NoEmptyStrings,
									},
								},
							},
						},

						"liveness_probe": azure.SchemaContainerGroupProbe(),

						"readiness_probe": azure.SchemaContainerGroupProbe(),
					},
				},
			},

			"diagnostics": {
				Type:     schema.TypeList,
				Optional: true,
				ForceNew: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"log_analytics": {
							Type:     schema.TypeList,
							Required: true,
							ForceNew: true,
							MaxItems: 1,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"workspace_id": {
										Type:         schema.TypeString,
										Required:     true,
										ForceNew:     true,
										ValidateFunc: validate.UUID,
									},

									"workspace_key": {
										Type:         schema.TypeString,
										Required:     true,
										Sensitive:    true,
										ForceNew:     true,
										ValidateFunc: validate.NoEmptyStrings,
									},

									"log_type": {
										Type:     schema.TypeString,
										Required: true,
										ForceNew: true,
										ValidateFunc: validation.StringInSlice([]string{
											string(containerinstance.ContainerInsights),
											string(containerinstance.ContainerInstanceLogs),
										}, false),
									},

									"metadata": {
										Type:     schema.TypeMap,
										Optional: true,
										ForceNew: true,
										Elem: &schema.Schema{
											Type: schema.TypeString,
										},
									},
								},
							},
						},
					},
				},
			},

			"ip_address": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"fqdn": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceArmContainerGroupCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ArmClient).containerGroupsClient
	ctx := meta.(*ArmClient).StopContext

	resGroup := d.Get("resource_group_name").(string)
	name := d.Get("name").(string)

	if requireResourcesToBeImported && d.IsNewResource() {
		existing, err := client.Get(ctx, resGroup, name)
		if err != nil {
			if !utils.ResponseWasNotFound(existing.Response) {
				return fmt.Errorf("Error checking for presence of existing Container Group %q (Resource Group %q): %s", name, resGroup, err)
			}
		}

		if existing.ID != nil && *existing.ID != "" {
			return tf.ImportAsExistsError("azurerm_container_group", *existing.ID)
		}
	}

	location := azureRMNormalizeLocation(d.Get("location").(string))
	OSType := d.Get("os_type").(string)
	IPAddressType := d.Get("ip_address_type").(string)
	tags := d.Get("tags").(map[string]interface{})
	restartPolicy := d.Get("restart_policy").(string)

	diagnosticsRaw := d.Get("diagnostics").([]interface{})
	diagnostics := expandContainerGroupDiagnostics(diagnosticsRaw)

	containers, containerGroupPorts, containerGroupVolumes := expandContainerGroupContainers(d)
	containerGroup := containerinstance.ContainerGroup{
		Name:     &name,
		Location: &location,
		Tags:     expandTags(tags),
		Identity: expandContainerGroupIdentity(d),
		ContainerGroupProperties: &containerinstance.ContainerGroupProperties{
			Containers:    containers,
			Diagnostics:   diagnostics,
			RestartPolicy: containerinstance.ContainerGroupRestartPolicy(restartPolicy),
			IPAddress: &containerinstance.IPAddress{
				Type:  containerinstance.ContainerGroupIPAddressType(IPAddressType),
				Ports: containerGroupPorts,
			},
			OsType:                   containerinstance.OperatingSystemTypes(OSType),
			Volumes:                  containerGroupVolumes,
			ImageRegistryCredentials: expandContainerImageRegistryCredentials(d),
		},
	}

	if dnsNameLabel := d.Get("dns_name_label").(string); dnsNameLabel != "" {
		containerGroup.ContainerGroupProperties.IPAddress.DNSNameLabel = &dnsNameLabel
	}

	future, err := client.CreateOrUpdate(ctx, resGroup, name, containerGroup)
	if err != nil {
		return fmt.Errorf("Error creating/updating container group %q (Resource Group %q): %+v", name, resGroup, err)
	}

	if err = future.WaitForCompletionRef(ctx, client.Client); err != nil {
		return fmt.Errorf("Error waiting for completion of container group %q (Resource Group %q): %+v", name, resGroup, err)
	}

	read, err := client.Get(ctx, resGroup, name)
	if err != nil {
		return err
	}

	if read.ID == nil {
		return fmt.Errorf("Cannot read container group %s (resource group %s) ID", name, resGroup)
	}

	d.SetId(*read.ID)

	return resourceArmContainerGroupRead(d, meta)
}

func resourceArmContainerGroupRead(d *schema.ResourceData, meta interface{}) error {
	ctx := meta.(*ArmClient).StopContext
	client := meta.(*ArmClient).containerGroupsClient

	id, err := parseAzureResourceID(d.Id())
	if err != nil {
		return err
	}

	resourceGroup := id.ResourceGroup
	name := id.Path["containerGroups"]

	resp, err := client.Get(ctx, resourceGroup, name)

	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			log.Printf("[DEBUG] Container Group %q was not found in Resource Group %q - removing from state!", name, resourceGroup)
			d.SetId("")
			return nil
		}
		return err
	}

	d.Set("name", name)
	d.Set("resource_group_name", resourceGroup)
	if location := resp.Location; location != nil {
		d.Set("location", azureRMNormalizeLocation(*location))
	}

	if err := d.Set("identity", flattenContainerGroupIdentity(d, resp.Identity)); err != nil {
		return fmt.Errorf("Error setting `identity`: %+v", err)
	}

	if props := resp.ContainerGroupProperties; props != nil {
		containerConfigs := flattenContainerGroupContainers(d, resp.Containers, props.IPAddress.Ports, props.Volumes)
		if err := d.Set("container", containerConfigs); err != nil {
			return fmt.Errorf("Error setting `container`: %+v", err)
		}

		if err := d.Set("image_registry_credential", flattenContainerImageRegistryCredentials(d, props.ImageRegistryCredentials)); err != nil {
			return fmt.Errorf("Error setting `image_registry_credential`: %+v", err)
		}

		if address := props.IPAddress; address != nil {
			d.Set("ip_address_type", address.Type)
			d.Set("ip_address", address.IP)
			d.Set("dns_name_label", address.DNSNameLabel)
			d.Set("fqdn", address.Fqdn)
		}

		d.Set("restart_policy", string(props.RestartPolicy))
		d.Set("os_type", string(props.OsType))

		if err := d.Set("diagnostics", flattenContainerGroupDiagnostics(d, props.Diagnostics)); err != nil {
			return fmt.Errorf("Error setting `diagnostics`: %+v", err)
		}
	}

	flattenAndSetTags(d, resp.Tags)

	return nil
}

func resourceArmContainerGroupDelete(d *schema.ResourceData, meta interface{}) error {
	ctx := meta.(*ArmClient).StopContext
	client := meta.(*ArmClient).containerGroupsClient

	id, err := parseAzureResourceID(d.Id())
	if err != nil {
		return err
	}

	resourceGroup := id.ResourceGroup
	name := id.Path["containerGroups"]

	resp, err := client.Delete(ctx, resourceGroup, name)
	if err != nil {
		if !utils.ResponseWasNotFound(resp.Response) {
			return fmt.Errorf("Error deleting Container Group %q (Resource Group %q): %+v", name, resourceGroup, err)
		}
	}

	return nil
}

func expandContainerGroupContainers(d *schema.ResourceData) (*[]containerinstance.Container, *[]containerinstance.Port, *[]containerinstance.Volume) {
	containersConfig := d.Get("container").([]interface{})
	containers := make([]containerinstance.Container, 0)
	containerGroupPorts := make([]containerinstance.Port, 0)
	containerGroupVolumes := make([]containerinstance.Volume, 0)

	for _, containerConfig := range containersConfig {
		data := containerConfig.(map[string]interface{})

		name := data["name"].(string)
		image := data["image"].(string)
		cpu := data["cpu"].(float64)
		memory := data["memory"].(float64)

		container := containerinstance.Container{
			Name: utils.String(name),
			ContainerProperties: &containerinstance.ContainerProperties{
				Image: utils.String(image),
				Resources: &containerinstance.ResourceRequirements{
					Requests: &containerinstance.ResourceRequests{
						MemoryInGB: utils.Float(memory),
						CPU:        utils.Float(cpu),
					},
				},
			},
		}

		if v, ok := data["gpu"]; ok {
			gpus := v.([]interface{})
			for _, gpuRaw := range gpus {
				v := gpuRaw.(map[string]interface{})
				gpuCount := int32(v["count"].(int))
				gpuSku := containerinstance.GpuSku(v["sku"].(string))

				gpus := containerinstance.GpuResource{
					Count: &gpuCount,
					Sku:   gpuSku,
				}
				container.Resources.Requests.Gpu = &gpus
			}
		}

		if v, ok := data["ports"].(*schema.Set); ok && len(v.List()) > 0 {
			var ports []containerinstance.ContainerPort
			for _, v := range v.List() {
				portObj := v.(map[string]interface{})

				port := int32(portObj["port"].(int))
				proto := portObj["protocol"].(string)

				ports = append(ports, containerinstance.ContainerPort{
					Port:     &port,
					Protocol: containerinstance.ContainerNetworkProtocol(proto),
				})
				containerGroupPorts = append(containerGroupPorts, containerinstance.Port{
					Port:     &port,
					Protocol: containerinstance.ContainerGroupNetworkProtocol(proto),
				})
			}
			container.Ports = &ports
		} else {
			if v := int32(data["port"].(int)); v != 0 {
				ports := []containerinstance.ContainerPort{
					{
						Port: &v,
					},
				}

				port := containerinstance.Port{
					Port: &v,
				}

				if v, ok := data["protocol"].(string); ok {
					ports[0].Protocol = containerinstance.ContainerNetworkProtocol(v)
					port.Protocol = containerinstance.ContainerGroupNetworkProtocol(v)
				}

				container.Ports = &ports
				containerGroupPorts = append(containerGroupPorts, port)
			}
		}

		// Set both sensitive and non-secure environment variables
		var envVars *[]containerinstance.EnvironmentVariable
		var secEnvVars *[]containerinstance.EnvironmentVariable

		// Expand environment_variables into slice
		if v, ok := data["environment_variables"]; ok {
			envVars = expandContainerEnvironmentVariables(v, false)
		}

		// Expand secure_environment_variables into slice
		if v, ok := data["secure_environment_variables"]; ok {
			secEnvVars = expandContainerEnvironmentVariables(v, true)
		}

		// Combine environment variable slices
		*envVars = append(*envVars, *secEnvVars...)

		// Set both secure and non secure environment variables
		container.EnvironmentVariables = envVars

		if v, ok := data["commands"]; ok {
			c := v.([]interface{})
			command := make([]string, 0)
			for _, v := range c {
				command = append(command, v.(string))
			}

			container.Command = &command
		}

		if container.Command == nil {
			if v := data["command"]; v != "" {
				command := strings.Split(v.(string), " ")
				container.Command = &command
			}
		}

		if v, ok := data["volume"]; ok {
			volumeMounts, containerGroupVolumesPartial := expandContainerVolumes(v)
			container.VolumeMounts = volumeMounts
			if containerGroupVolumesPartial != nil {
				containerGroupVolumes = append(containerGroupVolumes, *containerGroupVolumesPartial...)
			}
		}

		if v, ok := data["liveness_probe"]; ok {
			container.ContainerProperties.LivenessProbe = expandContainerProbe(v)
		}

		if v, ok := data["readiness_probe"]; ok {
			container.ContainerProperties.ReadinessProbe = expandContainerProbe(v)
		}

		containers = append(containers, container)
	}

	return &containers, &containerGroupPorts, &containerGroupVolumes
}

func expandContainerEnvironmentVariables(input interface{}, secure bool) *[]containerinstance.EnvironmentVariable {

	envVars := input.(map[string]interface{})
	output := make([]containerinstance.EnvironmentVariable, 0, len(envVars))

	if secure {

		for k, v := range envVars {
			ev := containerinstance.EnvironmentVariable{
				Name:        utils.String(k),
				SecureValue: utils.String(v.(string)),
			}

			output = append(output, ev)
		}

	} else {

		for k, v := range envVars {
			ev := containerinstance.EnvironmentVariable{
				Name:  utils.String(k),
				Value: utils.String(v.(string)),
			}

			output = append(output, ev)
		}
	}
	return &output
}

func expandContainerGroupIdentity(d *schema.ResourceData) *containerinstance.ContainerGroupIdentity {
	v := d.Get("identity")
	identities := v.([]interface{})
	if len(identities) == 0 {
		return nil
	}
	identity := identities[0].(map[string]interface{})
	identityType := containerinstance.ResourceIdentityType(identity["type"].(string))

	identityIds := make(map[string]*containerinstance.ContainerGroupIdentityUserAssignedIdentitiesValue)
	for _, id := range identity["identity_ids"].([]interface{}) {
		identityIds[id.(string)] = &containerinstance.ContainerGroupIdentityUserAssignedIdentitiesValue{}
	}

	cgIdentity := containerinstance.ContainerGroupIdentity{
		Type: identityType,
	}

	if cgIdentity.Type == containerinstance.UserAssigned || cgIdentity.Type == containerinstance.SystemAssignedUserAssigned {
		cgIdentity.UserAssignedIdentities = identityIds
	}

	return &cgIdentity
}

func expandContainerImageRegistryCredentials(d *schema.ResourceData) *[]containerinstance.ImageRegistryCredential {
	credsRaw := d.Get("image_registry_credential").([]interface{})
	if len(credsRaw) == 0 {
		return nil
	}

	output := make([]containerinstance.ImageRegistryCredential, 0, len(credsRaw))

	for _, c := range credsRaw {
		credConfig := c.(map[string]interface{})

		output = append(output, containerinstance.ImageRegistryCredential{
			Server:   utils.String(credConfig["server"].(string)),
			Password: utils.String(credConfig["password"].(string)),
			Username: utils.String(credConfig["username"].(string)),
		})
	}

	return &output
}

func expandContainerVolumes(input interface{}) (*[]containerinstance.VolumeMount, *[]containerinstance.Volume) {
	volumesRaw := input.([]interface{})

	if len(volumesRaw) == 0 {
		return nil, nil
	}

	volumeMounts := make([]containerinstance.VolumeMount, 0)
	containerGroupVolumes := make([]containerinstance.Volume, 0)

	for _, volumeRaw := range volumesRaw {
		volumeConfig := volumeRaw.(map[string]interface{})

		name := volumeConfig["name"].(string)
		mountPath := volumeConfig["mount_path"].(string)
		readOnly := volumeConfig["read_only"].(bool)
		shareName := volumeConfig["share_name"].(string)
		storageAccountName := volumeConfig["storage_account_name"].(string)
		storageAccountKey := volumeConfig["storage_account_key"].(string)

		vm := containerinstance.VolumeMount{
			Name:      utils.String(name),
			MountPath: utils.String(mountPath),
			ReadOnly:  utils.Bool(readOnly),
		}

		volumeMounts = append(volumeMounts, vm)

		cv := containerinstance.Volume{
			Name: utils.String(name),
			AzureFile: &containerinstance.AzureFileVolume{
				ShareName:          utils.String(shareName),
				ReadOnly:           utils.Bool(readOnly),
				StorageAccountName: utils.String(storageAccountName),
				StorageAccountKey:  utils.String(storageAccountKey),
			},
		}

		containerGroupVolumes = append(containerGroupVolumes, cv)
	}

	return &volumeMounts, &containerGroupVolumes
}

func expandContainerProbe(input interface{}) *containerinstance.ContainerProbe {
	probe := containerinstance.ContainerProbe{}
	probeRaw := input.([]interface{})

	if len(probeRaw) == 0 {
		return nil
	}

	for _, p := range probeRaw {
		probeConfig := p.(map[string]interface{})

		if v := probeConfig["initial_delay_seconds"].(int); v > 0 {
			probe.InitialDelaySeconds = utils.Int32(int32(v))
		}

		if v := probeConfig["period_seconds"].(int); v > 0 {
			probe.PeriodSeconds = utils.Int32(int32(v))
		}

		if v := probeConfig["failure_threshold"].(int); v > 0 {
			probe.FailureThreshold = utils.Int32(int32(v))
		}

		if v := probeConfig["success_threshold"].(int); v > 0 {
			probe.SuccessThreshold = utils.Int32(int32(v))
		}

		if v := probeConfig["timeout_seconds"].(int); v > 0 {
			probe.TimeoutSeconds = utils.Int32(int32(v))
		}

		commands := probeConfig["exec"].([]interface{})
		if len(commands) > 0 {
			exec := containerinstance.ContainerExec{
				Command: utils.ExpandStringArray(commands),
			}
			probe.Exec = &exec
		}

		httpRaw := probeConfig["http_get"].([]interface{})
		if len(httpRaw) > 0 {

			for _, httpget := range httpRaw {
				x := httpget.(map[string]interface{})

				path := x["path"].(string)
				port := x["port"].(int)
				scheme := x["scheme"].(string)

				probe.HTTPGet = &containerinstance.ContainerHTTPGet{
					Path:   utils.String(path),
					Port:   utils.Int32(int32(port)),
					Scheme: containerinstance.Scheme(scheme),
				}
			}
		}
	}
	return &probe
}

func flattenContainerGroupIdentity(d *schema.ResourceData, identity *containerinstance.ContainerGroupIdentity) []interface{} {
	if identity == nil {
		return make([]interface{}, 0)
	}

	result := make(map[string]interface{})
	result["type"] = string(identity.Type)
	if identity.PrincipalID != nil {
		result["principal_id"] = *identity.PrincipalID
	}

	identityIds := make([]string, 0)
	if identity.UserAssignedIdentities != nil {
		/*
			"userAssignedIdentities": {
			  "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/tomdevidentity/providers/Microsoft.ManagedIdentity/userAssignedIdentities/tom123": {
				"principalId": "00000000-0000-0000-0000-000000000000",
				"clientId": "00000000-0000-0000-0000-000000000000"
			  }
			}
		*/
		for key := range identity.UserAssignedIdentities {
			identityIds = append(identityIds, key)
		}
	}
	result["identity_ids"] = identityIds

	return []interface{}{result}
}

func flattenContainerImageRegistryCredentials(d *schema.ResourceData, input *[]containerinstance.ImageRegistryCredential) []interface{} {
	if input == nil {
		return nil
	}
	configsOld := d.Get("image_registry_credential").([]interface{})

	output := make([]interface{}, 0)
	for i, cred := range *input {
		credConfig := make(map[string]interface{})
		if cred.Server != nil {
			credConfig["server"] = *cred.Server
		}
		if cred.Username != nil {
			credConfig["username"] = *cred.Username
		}

		if len(configsOld) > i {
			data := configsOld[i].(map[string]interface{})
			oldServer := data["server"].(string)
			if cred.Server != nil && *cred.Server == oldServer {
				if v, ok := d.GetOk(fmt.Sprintf("image_registry_credential.%d.password", i)); ok {
					credConfig["password"] = v.(string)
				}
			}
		}

		output = append(output, credConfig)
	}
	return output
}

func flattenContainerGroupContainers(d *schema.ResourceData, containers *[]containerinstance.Container, containerGroupPorts *[]containerinstance.Port, containerGroupVolumes *[]containerinstance.Volume) []interface{} {

	//map old container names to index so we can look up things up
	nameIndexMap := map[string]int{}
	for i, c := range d.Get("container").([]interface{}) {
		cfg := c.(map[string]interface{})
		nameIndexMap[cfg["name"].(string)] = i
	}

	containerCfg := make([]interface{}, 0, len(*containers))
	for _, container := range *containers {

		//TODO fix this crash point
		name := *container.Name

		//get index from name
		index := nameIndexMap[name]

		containerConfig := make(map[string]interface{})
		containerConfig["name"] = name

		if v := container.Image; v != nil {
			containerConfig["image"] = *v
		}

		if resources := container.Resources; resources != nil {
			if resourceRequests := resources.Requests; resourceRequests != nil {
				if v := resourceRequests.CPU; v != nil {
					containerConfig["cpu"] = *v
				}
				if v := resourceRequests.MemoryInGB; v != nil {
					containerConfig["memory"] = *v
				}

				gpus := make([]interface{}, 0)
				if v := resourceRequests.Gpu; v != nil {
					gpu := make(map[string]interface{})
					if v.Count != nil {
						gpu["count"] = *v.Count
					}
					gpu["sku"] = string(v.Sku)
					gpus = append(gpus, gpu)
				}
				containerConfig["gpu"] = gpus
			}
		}

		if cPorts := container.Ports; cPorts != nil && len(*cPorts) > 0 {
			ports := make([]interface{}, 0)
			for _, p := range *cPorts {
				port := make(map[string]interface{})
				if v := p.Port; v != nil {
					port["port"] = int(*v)
				}
				port["protocol"] = string(p.Protocol)
				ports = append(ports, port)
			}
			containerConfig["ports"] = schema.NewSet(resourceArmContainerGroupPortsHash, ports)

			//old deprecated code
			containerPort := *(*cPorts)[0].Port
			containerConfig["port"] = containerPort
			// protocol isn't returned in container config, have to search in container group ports
			protocol := ""
			if containerGroupPorts != nil {
				for _, cgPort := range *containerGroupPorts {
					if *cgPort.Port == containerPort {
						protocol = string(cgPort.Protocol)
					}
				}
			}
			if protocol != "" {
				containerConfig["protocol"] = protocol
			}
		}

		if container.EnvironmentVariables != nil {
			if len(*container.EnvironmentVariables) > 0 {
				containerConfig["environment_variables"] = flattenContainerEnvironmentVariables(container.EnvironmentVariables, false, d, index)
			}
		}

		if container.EnvironmentVariables != nil {
			if len(*container.EnvironmentVariables) > 0 {
				containerConfig["secure_environment_variables"] = flattenContainerEnvironmentVariables(container.EnvironmentVariables, true, d, index)
			}
		}

		commands := make([]string, 0)
		if command := container.Command; command != nil {
			containerConfig["command"] = strings.Join(*command, " ")
			commands = *command
		}
		containerConfig["commands"] = commands

		if containerGroupVolumes != nil && container.VolumeMounts != nil {
			// Also pass in the container volume config from schema
			var containerVolumesConfig *[]interface{}
			containersConfigRaw := d.Get("container").([]interface{})
			for _, containerConfigRaw := range containersConfigRaw {
				data := containerConfigRaw.(map[string]interface{})
				nameRaw := data["name"].(string)
				if nameRaw == *container.Name {
					// found container config for current container
					// extract volume mounts from config
					if v, ok := data["volume"]; ok {
						containerVolumesRaw := v.([]interface{})
						containerVolumesConfig = &containerVolumesRaw
					}
				}
			}
			containerConfig["volume"] = flattenContainerVolumes(container.VolumeMounts, containerGroupVolumes, containerVolumesConfig)
		}

		containerConfig["liveness_probe"] = flattenContainerProbes(container.LivenessProbe)
		containerConfig["readiness_probe"] = flattenContainerProbes(container.ReadinessProbe)

		containerCfg = append(containerCfg, containerConfig)
	}

	return containerCfg
}

func flattenContainerEnvironmentVariables(input *[]containerinstance.EnvironmentVariable, isSecure bool, d *schema.ResourceData, oldContainerIndex int) map[string]interface{} {
	output := make(map[string]interface{})

	if input == nil {
		return output
	}

	if isSecure {
		for _, envVar := range *input {

			if envVar.Name != nil && envVar.Value == nil {
				if v, ok := d.GetOk(fmt.Sprintf("container.%d.secure_environment_variables.%s", oldContainerIndex, *envVar.Name)); ok {
					log.Printf("[DEBUG] SECURE    : Name: %s - Value: %s", *envVar.Name, v.(string))
					output[*envVar.Name] = v.(string)
				}
			}
		}
	} else {
		for _, envVar := range *input {
			if envVar.Name != nil && envVar.Value != nil {
				log.Printf("[DEBUG] NOT SECURE: Name: %s - Value: %s", *envVar.Name, *envVar.Value)
				output[*envVar.Name] = *envVar.Value
			}
		}
	}

	return output
}

func flattenContainerVolumes(volumeMounts *[]containerinstance.VolumeMount, containerGroupVolumes *[]containerinstance.Volume, containerVolumesConfig *[]interface{}) []interface{} {
	volumeConfigs := make([]interface{}, 0)

	if volumeMounts == nil {
		return volumeConfigs
	}

	for _, vm := range *volumeMounts {
		volumeConfig := make(map[string]interface{})
		if vm.Name != nil {
			volumeConfig["name"] = *vm.Name
		}
		if vm.MountPath != nil {
			volumeConfig["mount_path"] = *vm.MountPath
		}
		if vm.ReadOnly != nil {
			volumeConfig["read_only"] = *vm.ReadOnly
		}

		// find corresponding volume in container group volumes
		// and use the data
		if containerGroupVolumes != nil {
			for _, cgv := range *containerGroupVolumes {
				if cgv.Name == nil || vm.Name == nil {
					continue
				}

				if *cgv.Name == *vm.Name {
					if file := cgv.AzureFile; file != nil {
						if file.ShareName != nil {
							volumeConfig["share_name"] = *file.ShareName
						}
						if file.StorageAccountName != nil {
							volumeConfig["storage_account_name"] = *file.StorageAccountName
						}
						// skip storage_account_key, is always nil
					}
				}
			}
		}

		// find corresponding volume in config
		// and use the data
		if containerVolumesConfig != nil {
			for _, cvr := range *containerVolumesConfig {
				cv := cvr.(map[string]interface{})
				rawName := cv["name"].(string)
				if vm.Name != nil && *vm.Name == rawName {
					storageAccountKey := cv["storage_account_key"].(string)
					volumeConfig["storage_account_key"] = storageAccountKey
				}
			}
		}

		volumeConfigs = append(volumeConfigs, volumeConfig)
	}

	return volumeConfigs
}

func flattenContainerProbes(input *containerinstance.ContainerProbe) []interface{} {
	outputs := make([]interface{}, 0)
	if input == nil {
		return outputs
	}

	output := make(map[string]interface{})

	if v := input.Exec; v != nil {
		output["exec"] = *v.Command
	}

	httpGets := make([]interface{}, 0)
	if get := input.HTTPGet; get != nil {
		httpGet := make(map[string]interface{})

		if v := get.Path; v != nil {
			httpGet["path"] = *v
		}

		if v := get.Port; v != nil {
			httpGet["port"] = *v
		}

		if get.Scheme != "" {
			httpGet["scheme"] = get.Scheme
		}

		httpGets = append(httpGets, httpGet)
	}
	output["http_get"] = httpGets

	if v := input.FailureThreshold; v != nil {
		output["failure_threshold"] = *v
	}

	if v := input.InitialDelaySeconds; v != nil {
		output["initial_delay_seconds"] = *v
	}

	if v := input.PeriodSeconds; v != nil {
		output["period_seconds"] = *v
	}

	if v := input.SuccessThreshold; v != nil {
		output["success_threshold"] = *v
	}

	if v := input.TimeoutSeconds; v != nil {
		output["timeout_seconds"] = *v
	}

	outputs = append(outputs, output)
	return outputs
}

func expandContainerGroupDiagnostics(input []interface{}) *containerinstance.ContainerGroupDiagnostics {
	if len(input) == 0 {
		return nil
	}

	vs := input[0].(map[string]interface{})

	analyticsVs := vs["log_analytics"].([]interface{})
	analyticsV := analyticsVs[0].(map[string]interface{})

	workspaceId := analyticsV["workspace_id"].(string)
	workspaceKey := analyticsV["workspace_key"].(string)
	logType := containerinstance.LogAnalyticsLogType(analyticsV["log_type"].(string))

	metadataMap := analyticsV["metadata"].(map[string]interface{})
	metadata := make(map[string]*string)
	for k, v := range metadataMap {
		strValue := v.(string)
		metadata[k] = &strValue
	}

	return &containerinstance.ContainerGroupDiagnostics{
		LogAnalytics: &containerinstance.LogAnalytics{
			WorkspaceID:  utils.String(workspaceId),
			WorkspaceKey: utils.String(workspaceKey),
			LogType:      logType,
			Metadata:     metadata,
		},
	}
}

func flattenContainerGroupDiagnostics(d *schema.ResourceData, input *containerinstance.ContainerGroupDiagnostics) []interface{} {
	if input == nil {
		return []interface{}{}
	}

	logAnalytics := make([]interface{}, 0)

	if la := input.LogAnalytics; la != nil {
		output := make(map[string]interface{})

		output["log_type"] = string(la.LogType)

		metadata := make(map[string]interface{})
		for k, v := range la.Metadata {
			metadata[k] = *v
		}
		output["metadata"] = metadata

		if la.WorkspaceID != nil {
			output["workspace_id"] = *la.WorkspaceID
		}

		// the existing config may not exist at Import time, protect against it.
		workspaceKey := ""
		if existingDiags := d.Get("diagnostics").([]interface{}); len(existingDiags) > 0 {
			existingDiag := existingDiags[0].(map[string]interface{})
			if existingLA := existingDiag["log_analytics"].([]interface{}); len(existingLA) > 0 {
				vs := existingLA[0].(map[string]interface{})
				if key := vs["workspace_key"]; key != nil && key.(string) != "" {
					workspaceKey = key.(string)
				}
			}
		}
		output["workspace_key"] = workspaceKey

		logAnalytics = append(logAnalytics, output)
	}

	return []interface{}{
		map[string]interface{}{
			"log_analytics": logAnalytics,
		},
	}
}

func resourceArmContainerGroupPortsHash(v interface{}) int {
	var buf bytes.Buffer

	if m, ok := v.(map[string]interface{}); ok {
		buf.WriteString(fmt.Sprintf("%d-", m["port"].(int)))
		buf.WriteString(fmt.Sprintf("%s-", m["protocol"].(string)))
	}

	return hashcode.String(buf.String())
}
