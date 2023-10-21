package helpers

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/Azure/azure-sdk-for-go/services/web/mgmt/2021-03-01/web" // nolint: staticcheck
	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/apimanagement/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

type SiteConfigLinux struct {
	AlwaysOn                bool                    `tfschema:"always_on"`
	ApiManagementConfigId   string                  `tfschema:"api_management_api_id"`
	ApiDefinition           string                  `tfschema:"api_definition_url"`
	AppCommandLine          string                  `tfschema:"app_command_line"`
	AutoHeal                bool                    `tfschema:"auto_heal_enabled"`
	AutoHealSettings        []AutoHealSettingLinux  `tfschema:"auto_heal_setting"`
	UseManagedIdentityACR   bool                    `tfschema:"container_registry_use_managed_identity"`
	ContainerRegistryMSI    string                  `tfschema:"container_registry_managed_identity_client_id"`
	DefaultDocuments        []string                `tfschema:"default_documents"`
	Http2Enabled            bool                    `tfschema:"http2_enabled"`
	IpRestriction           []IpRestriction         `tfschema:"ip_restriction"`
	ScmUseMainIpRestriction bool                    `tfschema:"scm_use_main_ip_restriction"`
	ScmIpRestriction        []IpRestriction         `tfschema:"scm_ip_restriction"`
	LoadBalancing           string                  `tfschema:"load_balancing_mode"`
	LocalMysql              bool                    `tfschema:"local_mysql_enabled"`
	ManagedPipelineMode     string                  `tfschema:"managed_pipeline_mode"`
	RemoteDebugging         bool                    `tfschema:"remote_debugging_enabled"`
	RemoteDebuggingVersion  string                  `tfschema:"remote_debugging_version"`
	ScmType                 string                  `tfschema:"scm_type"`
	Use32BitWorker          bool                    `tfschema:"use_32_bit_worker"`
	WebSockets              bool                    `tfschema:"websockets_enabled"`
	FtpsState               string                  `tfschema:"ftps_state"`
	HealthCheckPath         string                  `tfschema:"health_check_path"`
	HealthCheckEvictionTime int                     `tfschema:"health_check_eviction_time_in_min"`
	NumberOfWorkers         int                     `tfschema:"worker_count"`
	ApplicationStack        []ApplicationStackLinux `tfschema:"application_stack"`
	MinTlsVersion           string                  `tfschema:"minimum_tls_version"`
	ScmMinTlsVersion        string                  `tfschema:"scm_minimum_tls_version"`
	Cors                    []CorsSetting           `tfschema:"cors"`
	DetailedErrorLogging    bool                    `tfschema:"detailed_error_logging_enabled"`
	LinuxFxVersion          string                  `tfschema:"linux_fx_version"`
	VnetRouteAllEnabled     bool                    `tfschema:"vnet_route_all_enabled"`
	// SiteLimits []SiteLimitsSettings `tfschema:"site_limits"` // TODO - New block to (possibly) support? No way to configure this in the portal?
}

func SiteConfigSchemaLinux() *pluginsdk.Schema {
	return &pluginsdk.Schema{
		Type:     pluginsdk.TypeList,
		Required: true,
		MaxItems: 1,
		Elem: &pluginsdk.Resource{
			Schema: map[string]*pluginsdk.Schema{
				"always_on": {
					Type:     pluginsdk.TypeBool,
					Optional: true,
					Default:  true,
				},

				"api_management_api_id": {
					Type:         pluginsdk.TypeString,
					Optional:     true,
					ValidateFunc: validate.ApiID,
				},

				"api_definition_url": {
					Type:         pluginsdk.TypeString,
					Optional:     true,
					ValidateFunc: validation.IsURLWithHTTPorHTTPS,
				},

				"app_command_line": {
					Type:     pluginsdk.TypeString,
					Optional: true,
				},

				"application_stack": linuxApplicationStackSchema(),

				"auto_heal_enabled": {
					Type:     pluginsdk.TypeBool,
					Optional: true,
					RequiredWith: []string{
						"site_config.0.auto_heal_setting",
					},
				},

				"auto_heal_setting": autoHealSettingSchemaLinux(),

				"container_registry_use_managed_identity": {
					Type:     pluginsdk.TypeBool,
					Optional: true,
					Default:  false,
				},

				"container_registry_managed_identity_client_id": {
					Type:         pluginsdk.TypeString,
					Optional:     true,
					ValidateFunc: validation.IsUUID,
				},

				"default_documents": {
					Type:     pluginsdk.TypeList,
					Optional: true,
					Computed: true,
					Elem: &pluginsdk.Schema{
						Type: pluginsdk.TypeString,
					},
				},

				"http2_enabled": {
					Type:     pluginsdk.TypeBool,
					Optional: true,
					Default:  false,
				},

				"ip_restriction": IpRestrictionSchema(),

				"scm_use_main_ip_restriction": {
					Type:     pluginsdk.TypeBool,
					Optional: true,
					Default:  false,
				},

				"scm_ip_restriction": IpRestrictionSchema(),

				"local_mysql_enabled": {
					Type:     pluginsdk.TypeBool,
					Optional: true,
					Default:  false,
				},

				"load_balancing_mode": {
					Type:     pluginsdk.TypeString,
					Optional: true,
					Default:  "LeastRequests",
					ValidateFunc: validation.StringInSlice([]string{
						"LeastRequests", // Service default
						"WeightedRoundRobin",
						"LeastResponseTime",
						"WeightedTotalTraffic",
						"RequestHash",
						"PerSiteRoundRobin",
					}, false),
				},

				"managed_pipeline_mode": {
					Type:     pluginsdk.TypeString,
					Optional: true,
					Default:  string(web.ManagedPipelineModeIntegrated),
					ValidateFunc: validation.StringInSlice([]string{
						string(web.ManagedPipelineModeClassic),
						string(web.ManagedPipelineModeIntegrated),
					}, false),
				},

				"remote_debugging_enabled": {
					Type:     pluginsdk.TypeBool,
					Optional: true,
					Default:  false,
				},

				"remote_debugging_version": {
					Type:     pluginsdk.TypeString,
					Optional: true,
					Computed: true,
					ValidateFunc: validation.StringInSlice([]string{
						"VS2017",
						"VS2019",
					}, false),
				},

				"scm_type": {
					Type:     pluginsdk.TypeString,
					Computed: true,
				},

				"use_32_bit_worker": {
					Type:     pluginsdk.TypeBool,
					Optional: true,
					Default:  true,
				},

				"websockets_enabled": {
					Type:     pluginsdk.TypeBool,
					Optional: true,
					Default:  false,
				},

				"ftps_state": {
					Type:     pluginsdk.TypeString,
					Optional: true,
					Default:  string(web.FtpsStateDisabled),
					ValidateFunc: validation.StringInSlice([]string{
						string(web.FtpsStateAllAllowed),
						string(web.FtpsStateDisabled),
						string(web.FtpsStateFtpsOnly),
					}, false),
				},

				"health_check_path": {
					Type:     pluginsdk.TypeString,
					Optional: true,
				},

				"health_check_eviction_time_in_min": {
					Type:         pluginsdk.TypeInt,
					Optional:     true,
					Computed:     true,
					ValidateFunc: validation.IntBetween(2, 10),
					Description:  "The amount of time in minutes that a node is unhealthy before being removed from the load balancer. Possible values are between `2` and `10`. Defaults to `10`. Only valid in conjunction with `health_check_path`",
				},

				"worker_count": {
					Type:         pluginsdk.TypeInt,
					Optional:     true,
					Computed:     true,
					ValidateFunc: validation.IntBetween(1, 100),
				},

				"minimum_tls_version": {
					Type:     pluginsdk.TypeString,
					Optional: true,
					Default:  string(web.SupportedTLSVersionsOneFullStopTwo),
					ValidateFunc: validation.StringInSlice([]string{
						string(web.SupportedTLSVersionsOneFullStopZero),
						string(web.SupportedTLSVersionsOneFullStopOne),
						string(web.SupportedTLSVersionsOneFullStopTwo),
					}, false),
				},

				"scm_minimum_tls_version": {
					Type:     pluginsdk.TypeString,
					Optional: true,
					Default:  string(web.SupportedTLSVersionsOneFullStopTwo),
					ValidateFunc: validation.StringInSlice([]string{
						string(web.SupportedTLSVersionsOneFullStopZero),
						string(web.SupportedTLSVersionsOneFullStopOne),
						string(web.SupportedTLSVersionsOneFullStopTwo),
					}, false),
				},

				"cors": CorsSettingsSchema(),

				"vnet_route_all_enabled": {
					Type:        pluginsdk.TypeBool,
					Optional:    true,
					Default:     false,
					Description: "Should all outbound traffic to have Virtual Network Security Groups and User Defined Routes applied? Defaults to `false`.",
				},

				"detailed_error_logging_enabled": {
					Type:     pluginsdk.TypeBool,
					Computed: true,
				},

				"linux_fx_version": {
					Type:     pluginsdk.TypeString,
					Computed: true,
				},
			},
		},
	}
}

func SiteConfigSchemaLinuxComputed() *pluginsdk.Schema {
	return &pluginsdk.Schema{
		Type:     pluginsdk.TypeList,
		Computed: true,
		Elem: &pluginsdk.Resource{
			Schema: map[string]*pluginsdk.Schema{
				"always_on": {
					Type:     pluginsdk.TypeBool,
					Computed: true,
				},

				"api_management_api_id": {
					Type:     pluginsdk.TypeString,
					Computed: true,
				},

				"api_definition_url": {
					Type:     pluginsdk.TypeString,
					Computed: true,
				},

				"app_command_line": {
					Type:     pluginsdk.TypeString,
					Computed: true,
				},

				"application_stack": linuxApplicationStackSchemaComputed(),

				"auto_heal_enabled": {
					Type:     pluginsdk.TypeBool,
					Computed: true,
				},

				"auto_heal_setting": autoHealSettingSchemaLinuxComputed(),

				"container_registry_use_managed_identity": {
					Type:     pluginsdk.TypeBool,
					Computed: true,
				},

				"container_registry_managed_identity_client_id": {
					Type:     pluginsdk.TypeString,
					Computed: true,
				},

				"default_documents": {
					Type:     pluginsdk.TypeList,
					Computed: true,
					Elem: &pluginsdk.Schema{
						Type: pluginsdk.TypeString,
					},
				},

				"http2_enabled": {
					Type:     pluginsdk.TypeBool,
					Computed: true,
				},

				"ip_restriction": IpRestrictionSchemaComputed(),

				"scm_use_main_ip_restriction": {
					Type:     pluginsdk.TypeBool,
					Computed: true,
				},

				"scm_ip_restriction": IpRestrictionSchemaComputed(),

				"local_mysql_enabled": {
					Type:     pluginsdk.TypeBool,
					Computed: true,
				},

				"load_balancing_mode": {
					Type:     pluginsdk.TypeString,
					Computed: true,
				},

				"managed_pipeline_mode": {
					Type:     pluginsdk.TypeString,
					Computed: true,
				},

				"remote_debugging_enabled": {
					Type:     pluginsdk.TypeBool,
					Computed: true,
				},

				"remote_debugging_version": {
					Type:     pluginsdk.TypeString,
					Computed: true,
				},

				"scm_type": {
					Type:     pluginsdk.TypeString,
					Computed: true,
				},

				"use_32_bit_worker": {
					Type:     pluginsdk.TypeBool,
					Computed: true,
				},

				"websockets_enabled": {
					Type:     pluginsdk.TypeBool,
					Computed: true,
				},

				"ftps_state": {
					Type:     pluginsdk.TypeString,
					Computed: true,
				},

				"health_check_path": {
					Type:     pluginsdk.TypeString,
					Computed: true,
				},

				"health_check_eviction_time_in_min": {
					Type:     pluginsdk.TypeInt,
					Computed: true,
				},

				"worker_count": {
					Type:     pluginsdk.TypeInt,
					Computed: true,
				},

				"minimum_tls_version": {
					Type:     pluginsdk.TypeString,
					Computed: true,
				},

				"scm_minimum_tls_version": {
					Type:     pluginsdk.TypeString,
					Computed: true,
				},

				"cors": CorsSettingsSchemaComputed(),

				"detailed_error_logging_enabled": {
					Type:     pluginsdk.TypeBool,
					Computed: true,
				},

				"linux_fx_version": {
					Type:     pluginsdk.TypeString,
					Computed: true,
				},

				"vnet_route_all_enabled": {
					Type:     pluginsdk.TypeBool,
					Computed: true,
				},
			},
		},
	}
}

type AutoHealSettingLinux struct {
	Triggers []AutoHealTriggerLinux `tfschema:"trigger"`
	Actions  []AutoHealActionLinux  `tfschema:"action"`
}

type AutoHealTriggerLinux struct {
	Requests     []AutoHealRequestTrigger    `tfschema:"requests"`
	StatusCodes  []AutoHealStatusCodeTrigger `tfschema:"status_code"` // 0 or more, ranges split by `-`, ranges cannot use sub-status or win32 code
	SlowRequests []AutoHealSlowRequest       `tfschema:"slow_request"`
}

type AutoHealActionLinux struct {
	ActionType         string `tfschema:"action_type"`                    // Enum - Only `Recycle` allowed
	MinimumProcessTime string `tfschema:"minimum_process_execution_time"` // Minimum uptime for process before action will trigger
}

func autoHealSettingSchemaLinux() *pluginsdk.Schema {
	return &pluginsdk.Schema{
		Type:     pluginsdk.TypeList,
		Optional: true,
		MaxItems: 1,
		Elem: &pluginsdk.Resource{
			Schema: map[string]*pluginsdk.Schema{
				"trigger": autoHealTriggerSchemaLinux(),

				"action": autoHealActionSchemaLinux(),
			},
		},
		RequiredWith: []string{
			"site_config.0.auto_heal_enabled",
		},
	}
}

func autoHealSettingSchemaLinuxComputed() *pluginsdk.Schema {
	return &pluginsdk.Schema{
		Type:     pluginsdk.TypeList,
		Computed: true,
		Elem: &pluginsdk.Resource{
			Schema: map[string]*pluginsdk.Schema{
				"trigger": autoHealTriggerSchemaLinuxComputed(),

				"action": autoHealActionSchemaLinuxComputed(),
			},
		},
	}
}

func autoHealActionSchemaLinux() *pluginsdk.Schema {
	return &pluginsdk.Schema{
		Type:     pluginsdk.TypeList,
		Optional: true,
		MaxItems: 1,
		Elem: &pluginsdk.Resource{
			Schema: map[string]*pluginsdk.Schema{
				"action_type": {
					Type:     pluginsdk.TypeString,
					Required: true,
					ValidateFunc: validation.StringInSlice([]string{
						string(web.AutoHealActionTypeRecycle),
					}, false),
				},

				"minimum_process_execution_time": {
					Type:     pluginsdk.TypeString,
					Optional: true,
					Computed: true,
					// ValidateFunc: // TODO - Time in hh:mm:ss, because why not...
				},
			},
		},
	}
}

func autoHealActionSchemaLinuxComputed() *pluginsdk.Schema {
	return &pluginsdk.Schema{
		Type:     pluginsdk.TypeList,
		Computed: true,
		Elem: &pluginsdk.Resource{
			Schema: map[string]*pluginsdk.Schema{
				"action_type": {
					Type:     pluginsdk.TypeString,
					Computed: true,
				},

				"minimum_process_execution_time": {
					Type:     pluginsdk.TypeString,
					Computed: true,
				},
			},
		},
	}
}

// (@jackofallops) - trigger schemas intentionally left long-hand for now
func autoHealTriggerSchemaLinux() *pluginsdk.Schema {
	return &pluginsdk.Schema{
		Type:     pluginsdk.TypeList,
		Optional: true,
		MaxItems: 1,
		Elem: &pluginsdk.Resource{
			Schema: map[string]*pluginsdk.Schema{
				"requests": {
					Type:     pluginsdk.TypeList,
					Optional: true,
					MaxItems: 1,
					Elem: &pluginsdk.Resource{
						Schema: map[string]*pluginsdk.Schema{
							"count": {
								Type:         pluginsdk.TypeInt,
								Required:     true,
								ValidateFunc: validation.IntAtLeast(1),
							},

							"interval": {
								Type:     pluginsdk.TypeString,
								Required: true,
								// ValidateFunc: validation.IsRFC3339Time, // TODO should be hh:mm:ss - This is too loose, need to improve?
							},
						},
					},
				},

				"status_code": {
					Type:     pluginsdk.TypeList,
					Optional: true,
					Elem: &pluginsdk.Resource{
						Schema: map[string]*pluginsdk.Schema{
							"status_code_range": {
								Type:         pluginsdk.TypeString,
								Required:     true,
								ValidateFunc: nil, // TODO - status code range validation
							},

							"count": {
								Type:         pluginsdk.TypeInt,
								Required:     true,
								ValidateFunc: validation.IntAtLeast(1),
							},

							"interval": {
								Type:     pluginsdk.TypeString,
								Required: true,
								// ValidateFunc: validation.IsRFC3339Time,
							},

							"sub_status": {
								Type:         pluginsdk.TypeInt,
								Optional:     true,
								ValidateFunc: nil, // TODO - no docs on this, needs investigation
							},

							"win32_status": {
								Type:         pluginsdk.TypeString,
								Optional:     true,
								ValidateFunc: nil, // TODO - no docs on this, needs investigation
							},

							"path": {
								Type:         pluginsdk.TypeString,
								Optional:     true,
								ValidateFunc: validation.StringIsNotEmpty,
							},
						},
					},
				},

				"slow_request": {
					Type:     pluginsdk.TypeList,
					Optional: true,
					MaxItems: 1,
					Elem: &pluginsdk.Resource{
						Schema: map[string]*pluginsdk.Schema{
							"time_taken": {
								Type:     pluginsdk.TypeString,
								Required: true,
								// ValidateFunc: validation.IsRFC3339Time,
							},

							"interval": {
								Type:     pluginsdk.TypeString,
								Required: true,
								// ValidateFunc: validation.IsRFC3339Time,
							},

							"count": {
								Type:         pluginsdk.TypeInt,
								Required:     true,
								ValidateFunc: validation.IntAtLeast(1),
							},

							"path": {
								Type:         pluginsdk.TypeString,
								Optional:     true,
								ValidateFunc: validation.StringIsNotEmpty,
							},
						},
					},
				},
			},
		},
	}
}

func autoHealTriggerSchemaLinuxComputed() *pluginsdk.Schema {
	return &pluginsdk.Schema{
		Type:     pluginsdk.TypeList,
		Computed: true,
		Elem: &pluginsdk.Resource{
			Schema: map[string]*pluginsdk.Schema{
				"requests": {
					Type:     pluginsdk.TypeList,
					Computed: true,
					Elem: &pluginsdk.Resource{
						Schema: map[string]*pluginsdk.Schema{
							"count": {
								Type:     pluginsdk.TypeInt,
								Computed: true,
							},

							"interval": {
								Type:     pluginsdk.TypeString,
								Computed: true,
							},
						},
					},
				},

				"status_code": {
					Type:     pluginsdk.TypeList,
					Computed: true,
					Elem: &pluginsdk.Resource{
						Schema: map[string]*pluginsdk.Schema{
							"status_code_range": {
								Type:     pluginsdk.TypeString,
								Computed: true,
							},

							"count": {
								Type:     pluginsdk.TypeInt,
								Computed: true,
							},

							"interval": {
								Type:     pluginsdk.TypeString,
								Computed: true,
							},

							"sub_status": {
								Type:     pluginsdk.TypeInt,
								Computed: true,
							},

							"win32_status": {
								Type:     pluginsdk.TypeString,
								Computed: true,
							},

							"path": {
								Type:     pluginsdk.TypeString,
								Computed: true,
							},
						},
					},
				},

				"slow_request": {
					Type:     pluginsdk.TypeList,
					Optional: true,
					Elem: &pluginsdk.Resource{
						Schema: map[string]*pluginsdk.Schema{
							"time_taken": {
								Type:     pluginsdk.TypeString,
								Computed: true,
							},

							"interval": {
								Type:     pluginsdk.TypeString,
								Computed: true,
							},

							"count": {
								Type:     pluginsdk.TypeInt,
								Computed: true,
							},

							"path": {
								Type:     pluginsdk.TypeString,
								Computed: true,
							},
						},
					},
				},
			},
		},
	}
}

func ExpandSiteConfigLinux(siteConfig []SiteConfigLinux, existing *web.SiteConfig, metadata sdk.ResourceMetaData, servicePlan web.AppServicePlan) (*web.SiteConfig, error) {
	if len(siteConfig) == 0 {
		return nil, nil
	}
	expanded := &web.SiteConfig{}
	if existing != nil {
		expanded = existing
	}

	linuxSiteConfig := siteConfig[0]

	if servicePlan.Sku != nil && servicePlan.Sku.Name != nil {
		if isFreeOrSharedServicePlan(*servicePlan.Sku.Name) {
			if linuxSiteConfig.AlwaysOn {
				return nil, fmt.Errorf("always_on cannot be set to true when using Free, F1, D1 Sku")
			}
			if expanded.AlwaysOn != nil && *expanded.AlwaysOn {
				return nil, fmt.Errorf("always_on feature has to be turned off before switching to a free/shared Sku")
			}
		}
	}
	expanded.AlwaysOn = pointer.To(linuxSiteConfig.AlwaysOn)

	if metadata.ResourceData.HasChange("site_config.0.api_management_api_id") {
		expanded.APIManagementConfig = &web.APIManagementConfig{
			ID: pointer.To(linuxSiteConfig.ApiManagementConfigId),
		}
	}

	if metadata.ResourceData.HasChange("site_config.0.api_definition_url") {
		expanded.APIDefinition = &web.APIDefinitionInfo{
			URL: pointer.To(linuxSiteConfig.ApiDefinition),
		}
	}

	if metadata.ResourceData.HasChange("site_config.0.app_command_line") {
		expanded.AppCommandLine = pointer.To(linuxSiteConfig.AppCommandLine)
	}

	if metadata.ResourceData.HasChange("site_config.0.application_stack") {
		if len(linuxSiteConfig.ApplicationStack) == 1 {
			linuxAppStack := linuxSiteConfig.ApplicationStack[0]
			if linuxAppStack.NetFrameworkVersion != "" {
				expanded.LinuxFxVersion = pointer.To(fmt.Sprintf("DOTNETCORE|%s", linuxAppStack.NetFrameworkVersion))
			}

			if linuxAppStack.GoVersion != "" {
				expanded.LinuxFxVersion = pointer.To(fmt.Sprintf("GO|%s", linuxAppStack.GoVersion))
			}

			if linuxAppStack.PhpVersion != "" {
				expanded.LinuxFxVersion = pointer.To(fmt.Sprintf("PHP|%s", linuxAppStack.PhpVersion))
			}

			if linuxAppStack.NodeVersion != "" {
				expanded.LinuxFxVersion = pointer.To(fmt.Sprintf("NODE|%s", linuxAppStack.NodeVersion))
			}

			if linuxAppStack.RubyVersion != "" {
				expanded.LinuxFxVersion = pointer.To(fmt.Sprintf("RUBY|%s", linuxAppStack.RubyVersion))
			}

			if linuxAppStack.PythonVersion != "" {
				expanded.LinuxFxVersion = pointer.To(fmt.Sprintf("PYTHON|%s", linuxAppStack.PythonVersion))
			}

			if linuxAppStack.JavaServer != "" {
				javaString, err := JavaLinuxFxStringBuilder(linuxAppStack.JavaVersion, linuxAppStack.JavaServer, linuxAppStack.JavaServerVersion)
				if err != nil {
					return nil, fmt.Errorf("could not build linuxFxVersion string: %+v", err)
				}
				expanded.LinuxFxVersion = javaString
			}

			if linuxAppStack.DockerImage != "" {
				expanded.LinuxFxVersion = pointer.To(fmt.Sprintf("DOCKER|%s:%s", linuxAppStack.DockerImage, linuxAppStack.DockerImageTag))
			}
		} else {
			expanded.LinuxFxVersion = pointer.To("")
		}
	}

	expanded.AcrUseManagedIdentityCreds = pointer.To(linuxSiteConfig.UseManagedIdentityACR)

	if metadata.ResourceData.HasChange("site_config.0.container_registry_managed_identity_client_id") {
		expanded.AcrUserManagedIdentityID = pointer.To(linuxSiteConfig.ContainerRegistryMSI)
	}

	if metadata.ResourceData.HasChange("site_config.0.default_documents") {
		expanded.DefaultDocuments = &linuxSiteConfig.DefaultDocuments
	}

	expanded.HTTP20Enabled = pointer.To(linuxSiteConfig.Http2Enabled)

	if metadata.ResourceData.HasChange("site_config.0.ip_restriction") {
		ipRestrictions, err := ExpandIpRestrictions(linuxSiteConfig.IpRestriction)
		if err != nil {
			return nil, err
		}
		expanded.IPSecurityRestrictions = ipRestrictions
	}

	expanded.ScmIPSecurityRestrictionsUseMain = pointer.To(linuxSiteConfig.ScmUseMainIpRestriction)

	if metadata.ResourceData.HasChange("site_config.0.scm_ip_restriction") {
		scmIpRestrictions, err := ExpandIpRestrictions(linuxSiteConfig.ScmIpRestriction)
		if err != nil {
			return nil, err
		}
		expanded.ScmIPSecurityRestrictions = scmIpRestrictions
	}

	expanded.LocalMySQLEnabled = pointer.To(linuxSiteConfig.LocalMysql)

	if metadata.ResourceData.HasChange("site_config.0.load_balancing_mode") {
		expanded.LoadBalancing = web.SiteLoadBalancing(linuxSiteConfig.LoadBalancing)
	}

	if metadata.ResourceData.HasChange("site_config.0.managed_pipeline_mode") {
		expanded.ManagedPipelineMode = web.ManagedPipelineMode(linuxSiteConfig.ManagedPipelineMode)
	}

	expanded.RemoteDebuggingEnabled = pointer.To(linuxSiteConfig.RemoteDebugging)

	if metadata.ResourceData.HasChange("site_config.0.remote_debugging_version") {
		expanded.RemoteDebuggingVersion = pointer.To(linuxSiteConfig.RemoteDebuggingVersion)
	}

	expanded.Use32BitWorkerProcess = pointer.To(linuxSiteConfig.Use32BitWorker)

	expanded.WebSocketsEnabled = pointer.To(linuxSiteConfig.WebSockets)

	if metadata.ResourceData.HasChange("site_config.0.ftps_state") {
		expanded.FtpsState = web.FtpsState(linuxSiteConfig.FtpsState)
	}

	if metadata.ResourceData.HasChange("site_config.0.health_check_path") {
		expanded.HealthCheckPath = pointer.To(linuxSiteConfig.HealthCheckPath)
	}

	if metadata.ResourceData.HasChange("site_config.0.worker_count") {
		expanded.NumberOfWorkers = pointer.To(int32(linuxSiteConfig.NumberOfWorkers))
	}

	if metadata.ResourceData.HasChange("site_config.0.minimum_tls_version") {
		expanded.MinTLSVersion = web.SupportedTLSVersions(linuxSiteConfig.MinTlsVersion)
	}

	if metadata.ResourceData.HasChange("site_config.0.scm_minimum_tls_version") {
		expanded.ScmMinTLSVersion = web.SupportedTLSVersions(linuxSiteConfig.ScmMinTlsVersion)
	}

	if metadata.ResourceData.HasChange("site_config.0.cors") {
		cors := ExpandCorsSettings(linuxSiteConfig.Cors)
		if cors == nil {
			cors = &web.CorsSettings{
				AllowedOrigins: &[]string{},
			}
		}
		expanded.Cors = cors
	}

	expanded.AutoHealEnabled = pointer.To(linuxSiteConfig.AutoHeal)

	if metadata.ResourceData.HasChange("site_config.0.auto_heal_setting") {
		expanded.AutoHealRules = expandAutoHealSettingsLinux(linuxSiteConfig.AutoHealSettings)
	}

	if metadata.ResourceData.HasChange("site_config.0.vnet_route_all_enabled") {
		expanded.VnetRouteAllEnabled = pointer.To(linuxSiteConfig.VnetRouteAllEnabled)
	}

	return expanded, nil
}

func FlattenSiteConfigLinux(appSiteConfig *web.SiteConfig, healthCheckCount *int) []SiteConfigLinux {
	if appSiteConfig == nil {
		return nil
	}

	siteConfig := SiteConfigLinux{
		AlwaysOn:                pointer.From(appSiteConfig.AlwaysOn),
		AppCommandLine:          pointer.From(appSiteConfig.AppCommandLine),
		AutoHeal:                pointer.From(appSiteConfig.AutoHealEnabled),
		AutoHealSettings:        flattenAutoHealSettingsLinux(appSiteConfig.AutoHealRules),
		ContainerRegistryMSI:    pointer.From(appSiteConfig.AcrUserManagedIdentityID),
		DetailedErrorLogging:    pointer.From(appSiteConfig.DetailedErrorLoggingEnabled),
		Http2Enabled:            pointer.From(appSiteConfig.HTTP20Enabled),
		IpRestriction:           FlattenIpRestrictions(appSiteConfig.IPSecurityRestrictions),
		ManagedPipelineMode:     string(appSiteConfig.ManagedPipelineMode),
		ScmType:                 string(appSiteConfig.ScmType),
		FtpsState:               string(appSiteConfig.FtpsState),
		HealthCheckPath:         pointer.From(appSiteConfig.HealthCheckPath),
		HealthCheckEvictionTime: pointer.From(healthCheckCount),
		LoadBalancing:           string(appSiteConfig.LoadBalancing),
		LocalMysql:              pointer.From(appSiteConfig.LocalMySQLEnabled),
		MinTlsVersion:           string(appSiteConfig.MinTLSVersion),
		NumberOfWorkers:         int(pointer.From(appSiteConfig.NumberOfWorkers)),
		RemoteDebugging:         pointer.From(appSiteConfig.RemoteDebuggingEnabled),
		RemoteDebuggingVersion:  strings.ToUpper(pointer.From(appSiteConfig.RemoteDebuggingVersion)),
		ScmIpRestriction:        FlattenIpRestrictions(appSiteConfig.ScmIPSecurityRestrictions),
		ScmMinTlsVersion:        string(appSiteConfig.ScmMinTLSVersion),
		ScmUseMainIpRestriction: pointer.From(appSiteConfig.ScmIPSecurityRestrictionsUseMain),
		Use32BitWorker:          pointer.From(appSiteConfig.Use32BitWorkerProcess),
		UseManagedIdentityACR:   pointer.From(appSiteConfig.AcrUseManagedIdentityCreds),
		WebSockets:              pointer.From(appSiteConfig.WebSocketsEnabled),
		VnetRouteAllEnabled:     pointer.From(appSiteConfig.VnetRouteAllEnabled),
	}

	if appSiteConfig.APIManagementConfig != nil && appSiteConfig.APIManagementConfig.ID != nil {
		siteConfig.ApiManagementConfigId = *appSiteConfig.APIManagementConfig.ID
	}

	if appSiteConfig.APIDefinition != nil && appSiteConfig.APIDefinition.URL != nil {
		siteConfig.ApiDefinition = *appSiteConfig.APIDefinition.URL
	}

	if appSiteConfig.DefaultDocuments != nil {
		siteConfig.DefaultDocuments = *appSiteConfig.DefaultDocuments
	}

	if appSiteConfig.LinuxFxVersion != nil {
		var linuxAppStack ApplicationStackLinux
		siteConfig.LinuxFxVersion = *appSiteConfig.LinuxFxVersion
		// Decode the string to docker values
		linuxAppStack = decodeApplicationStackLinux(siteConfig.LinuxFxVersion)
		siteConfig.ApplicationStack = []ApplicationStackLinux{linuxAppStack}
	}

	if appSiteConfig.Cors != nil {
		corsEmpty := false
		corsSettings := appSiteConfig.Cors
		cors := CorsSetting{}
		if corsSettings.SupportCredentials != nil {
			cors.SupportCredentials = *corsSettings.SupportCredentials
		}

		if corsSettings.AllowedOrigins != nil {
			if len(*corsSettings.AllowedOrigins) > 0 {
				cors.AllowedOrigins = *corsSettings.AllowedOrigins
			} else if !cors.SupportCredentials {
				corsEmpty = true
			}
		}
		if !corsEmpty {
			siteConfig.Cors = []CorsSetting{cors}
		}
	}

	return []SiteConfigLinux{siteConfig}
}

func expandAutoHealSettingsLinux(autoHealSettings []AutoHealSettingLinux) *web.AutoHealRules {
	if len(autoHealSettings) == 0 {
		return nil
	}

	result := &web.AutoHealRules{
		Triggers: &web.AutoHealTriggers{},
		Actions:  &web.AutoHealActions{},
	}

	autoHeal := autoHealSettings[0]

	triggers := autoHeal.Triggers[0]
	if len(triggers.Requests) == 1 {
		result.Triggers.Requests = &web.RequestsBasedTrigger{
			Count:        pointer.To(int32(triggers.Requests[0].Count)),
			TimeInterval: pointer.To(triggers.Requests[0].Interval),
		}
	}

	if len(triggers.SlowRequests) == 1 {
		result.Triggers.SlowRequests = &web.SlowRequestsBasedTrigger{
			TimeTaken:    pointer.To(triggers.SlowRequests[0].TimeTaken),
			TimeInterval: pointer.To(triggers.SlowRequests[0].Interval),
			Count:        pointer.To(int32(triggers.SlowRequests[0].Count)),
		}
		if triggers.SlowRequests[0].Path != "" {
			result.Triggers.SlowRequests.Path = pointer.To(triggers.SlowRequests[0].Path)
		}
	}

	if len(triggers.StatusCodes) > 0 {
		statusCodeTriggers := make([]web.StatusCodesBasedTrigger, 0)
		statusCodeRangeTriggers := make([]web.StatusCodesRangeBasedTrigger, 0)
		for _, s := range triggers.StatusCodes {
			statusCodeTrigger := web.StatusCodesBasedTrigger{}
			statusCodeRangeTrigger := web.StatusCodesRangeBasedTrigger{}
			parts := strings.Split(s.StatusCodeRange, "-")
			if len(parts) == 2 {
				statusCodeRangeTrigger.StatusCodes = pointer.To(s.StatusCodeRange)
				statusCodeRangeTrigger.Count = pointer.To(int32(s.Count))
				statusCodeRangeTrigger.TimeInterval = pointer.To(s.Interval)
				if s.Path != "" {
					statusCodeRangeTrigger.Path = pointer.To(s.Path)
				}
				statusCodeRangeTriggers = append(statusCodeRangeTriggers, statusCodeRangeTrigger)
			} else {
				statusCode, err := strconv.Atoi(s.StatusCodeRange)
				if err == nil {
					statusCodeTrigger.Status = pointer.To(int32(statusCode))
				}
				statusCodeTrigger.Count = pointer.To(int32(s.Count))
				statusCodeTrigger.TimeInterval = pointer.To(s.Interval)
				if s.Path != "" {
					statusCodeTrigger.Path = pointer.To(s.Path)
				}
				statusCodeTriggers = append(statusCodeTriggers, statusCodeTrigger)
			}
		}
		result.Triggers.StatusCodes = &statusCodeTriggers
		result.Triggers.StatusCodesRange = &statusCodeRangeTriggers
	}

	action := autoHeal.Actions[0]
	result.Actions.ActionType = web.AutoHealActionType(action.ActionType)
	result.Actions.MinProcessExecutionTime = pointer.To(action.MinimumProcessTime)

	return result
}

func flattenAutoHealSettingsLinux(autoHealRules *web.AutoHealRules) []AutoHealSettingLinux {
	if autoHealRules == nil {
		return nil
	}

	result := AutoHealSettingLinux{}

	// Triggers
	if autoHealRules.Triggers != nil {
		resultTrigger := AutoHealTriggerLinux{}
		triggers := *autoHealRules.Triggers
		if triggers.Requests != nil {
			count := 0
			if triggers.Requests.Count != nil {
				count = int(*triggers.Requests.Count)
			}
			resultTrigger.Requests = []AutoHealRequestTrigger{{
				Count:    count,
				Interval: pointer.From(triggers.Requests.TimeInterval),
			}}
		}

		statusCodeTriggers := make([]AutoHealStatusCodeTrigger, 0)
		if triggers.StatusCodes != nil {
			for _, s := range *triggers.StatusCodes {
				t := AutoHealStatusCodeTrigger{
					Interval: pointer.From(s.TimeInterval),
					Path:     pointer.From(s.Path),
				}

				if s.Status != nil {
					t.StatusCodeRange = strconv.Itoa(int(*s.Status))
				}

				if s.Count != nil {
					t.Count = int(*s.Count)
				}

				if s.SubStatus != nil {
					t.SubStatus = int(*s.SubStatus)
				}
				statusCodeTriggers = append(statusCodeTriggers, t)
			}
		}
		if triggers.StatusCodesRange != nil {
			for _, s := range *triggers.StatusCodesRange {
				t := AutoHealStatusCodeTrigger{
					Interval: pointer.From(s.TimeInterval),
					Path:     pointer.From(s.Path),
				}
				if s.Count != nil {
					t.Count = int(*s.Count)
				}

				if s.StatusCodes != nil {
					t.StatusCodeRange = *s.StatusCodes
				}
				statusCodeTriggers = append(statusCodeTriggers, t)
			}
		}
		resultTrigger.StatusCodes = statusCodeTriggers

		slowRequestTriggers := make([]AutoHealSlowRequest, 0)
		if triggers.SlowRequests != nil {
			slowRequestTriggers = append(slowRequestTriggers, AutoHealSlowRequest{
				TimeTaken: pointer.From(triggers.SlowRequests.TimeTaken),
				Interval:  pointer.From(triggers.SlowRequests.TimeInterval),
				Count:     int(pointer.From(triggers.SlowRequests.Count)),
				Path:      pointer.From(triggers.SlowRequests.Path),
			})
		}
		resultTrigger.SlowRequests = slowRequestTriggers
		result.Triggers = []AutoHealTriggerLinux{resultTrigger}
	}

	// Actions
	if autoHealRules.Actions != nil {
		actions := *autoHealRules.Actions

		result.Actions = []AutoHealActionLinux{{
			ActionType:         string(actions.ActionType),
			MinimumProcessTime: pointer.From(actions.MinProcessExecutionTime),
		}}
	}

	if result.Triggers != nil || result.Actions != nil {
		return []AutoHealSettingLinux{result}
	}

	return nil
}
