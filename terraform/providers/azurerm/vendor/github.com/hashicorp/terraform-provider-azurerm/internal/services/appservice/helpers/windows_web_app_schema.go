package helpers

import (
	"fmt"
	"strings"

	"github.com/Azure/azure-sdk-for-go/services/web/mgmt/2021-03-01/web" // nolint: staticcheck
	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/apimanagement/parse"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/apimanagement/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

type SiteConfigWindows struct {
	AlwaysOn                 bool                      `tfschema:"always_on"`
	ApiManagementConfigId    string                    `tfschema:"api_management_api_id"`
	ApiDefinition            string                    `tfschema:"api_definition_url"`
	AppCommandLine           string                    `tfschema:"app_command_line"`
	AutoHeal                 bool                      `tfschema:"auto_heal_enabled"`
	AutoHealSettings         []AutoHealSettingWindows  `tfschema:"auto_heal_setting"`
	UseManagedIdentityACR    bool                      `tfschema:"container_registry_use_managed_identity"`
	ContainerRegistryUserMSI string                    `tfschema:"container_registry_managed_identity_client_id"`
	DefaultDocuments         []string                  `tfschema:"default_documents"`
	Http2Enabled             bool                      `tfschema:"http2_enabled"`
	IpRestriction            []IpRestriction           `tfschema:"ip_restriction"`
	ScmUseMainIpRestriction  bool                      `tfschema:"scm_use_main_ip_restriction"`
	ScmIpRestriction         []IpRestriction           `tfschema:"scm_ip_restriction"`
	LoadBalancing            string                    `tfschema:"load_balancing_mode"`
	LocalMysql               bool                      `tfschema:"local_mysql_enabled"`
	ManagedPipelineMode      string                    `tfschema:"managed_pipeline_mode"`
	RemoteDebugging          bool                      `tfschema:"remote_debugging_enabled"`
	RemoteDebuggingVersion   string                    `tfschema:"remote_debugging_version"`
	ScmType                  string                    `tfschema:"scm_type"`
	Use32BitWorker           bool                      `tfschema:"use_32_bit_worker"`
	WebSockets               bool                      `tfschema:"websockets_enabled"`
	FtpsState                string                    `tfschema:"ftps_state"`
	HealthCheckPath          string                    `tfschema:"health_check_path"`
	HealthCheckEvictionTime  int                       `tfschema:"health_check_eviction_time_in_min"`
	WorkerCount              int                       `tfschema:"worker_count"`
	ApplicationStack         []ApplicationStackWindows `tfschema:"application_stack"`
	VirtualApplications      []VirtualApplication      `tfschema:"virtual_application"`
	MinTlsVersion            string                    `tfschema:"minimum_tls_version"`
	ScmMinTlsVersion         string                    `tfschema:"scm_minimum_tls_version"`
	Cors                     []CorsSetting             `tfschema:"cors"`
	DetailedErrorLogging     bool                      `tfschema:"detailed_error_logging_enabled"`
	WindowsFxVersion         string                    `tfschema:"windows_fx_version"`
	VnetRouteAllEnabled      bool                      `tfschema:"vnet_route_all_enabled"`
	// TODO new properties / blocks
	// SiteLimits []SiteLimitsSettings `tfschema:"site_limits"` // TODO - ASE related for limiting App resource consumption
	// PushSettings - Supported in SDK, but blocked by manual step needed for connecting app to notification hub.
}

func SiteConfigSchemaWindows() *pluginsdk.Schema {
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

				"application_stack": windowsApplicationStackSchema(),

				"app_command_line": {
					Type:     pluginsdk.TypeString,
					Optional: true,
				},

				"auto_heal_enabled": {
					Type:     pluginsdk.TypeBool,
					Optional: true,
					Default:  false,
					RequiredWith: []string{
						"site_config.0.auto_heal_setting",
					},
				},

				"auto_heal_setting": autoHealSettingSchemaWindows(),

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
					Default:  string(web.SiteLoadBalancingLeastRequests),
					ValidateFunc: validation.StringInSlice([]string{
						string(web.SiteLoadBalancingLeastRequests),
						string(web.SiteLoadBalancingWeightedRoundRobin),
						string(web.SiteLoadBalancingLeastResponseTime),
						string(web.SiteLoadBalancingWeightedTotalTraffic),
						string(web.SiteLoadBalancingRequestHash),
						string(web.SiteLoadBalancingPerSiteRoundRobin),
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
					Default:  true, // Variable default value depending on several factors, such as plan type?
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

				"virtual_application": virtualApplicationsSchema(),

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

				"windows_fx_version": {
					Type:     pluginsdk.TypeString,
					Computed: true,
				},
			},
		},
	}
}

func SiteConfigSchemaWindowsComputed() *pluginsdk.Schema {
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

				"application_stack": windowsApplicationStackSchemaComputed(),

				"app_command_line": {
					Type:     pluginsdk.TypeString,
					Computed: true,
				},

				"auto_heal_enabled": {
					Type:     pluginsdk.TypeBool,
					Computed: true,
				},

				"auto_heal_setting": autoHealSettingSchemaWindowsComputed(),

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

				"virtual_application": virtualApplicationsSchemaComputed(),

				"detailed_error_logging_enabled": {
					Type:     pluginsdk.TypeBool,
					Computed: true,
				},

				"windows_fx_version": {
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

func ExpandSiteConfigWindows(siteConfig []SiteConfigWindows, existing *web.SiteConfig, metadata sdk.ResourceMetaData, servicePlan web.AppServicePlan) (*web.SiteConfig, *string, error) {
	if len(siteConfig) == 0 {
		return nil, nil, nil
	}

	expanded := &web.SiteConfig{}
	if existing != nil {
		expanded = existing
	}

	winSiteConfig := siteConfig[0]

	currentStack := ""
	if len(winSiteConfig.ApplicationStack) == 1 {
		winAppStack := winSiteConfig.ApplicationStack[0]
		currentStack = winAppStack.CurrentStack
	}

	if servicePlan.Sku != nil && servicePlan.Sku.Name != nil {
		if isFreeOrSharedServicePlan(*servicePlan.Sku.Name) {
			if winSiteConfig.AlwaysOn {
				return nil, nil, fmt.Errorf("always_on cannot be set to true when using Free, F1, D1 Sku")
			}
			if expanded.AlwaysOn != nil && *expanded.AlwaysOn {
				// TODO - This code will not work as expected due to the service plan always being updated before the App, so the apply will always error out on the Service Plan change if `always_on` is incorrectly set.
				// Need to investigate if there's a way to avoid users needing to run 2 applies for this.
				return nil, nil, fmt.Errorf("always_on feature has to be turned off before switching to a free/shared Sku")
			}
		}
	}
	expanded.AlwaysOn = pointer.To(winSiteConfig.AlwaysOn)

	if metadata.ResourceData.HasChange("site_config.0.api_management_api_id") {
		expanded.APIManagementConfig = &web.APIManagementConfig{
			ID: pointer.To(winSiteConfig.ApiManagementConfigId),
		}
	}

	if metadata.ResourceData.HasChange("site_config.0.api_definition_url") {
		expanded.APIDefinition = &web.APIDefinitionInfo{
			URL: pointer.To(winSiteConfig.ApiDefinition),
		}
	}

	if metadata.ResourceData.HasChange("site_config.0.app_command_line") {
		expanded.AppCommandLine = pointer.To(winSiteConfig.AppCommandLine)
	}

	if metadata.ResourceData.HasChange("site_config.0.application_stack") {
		if len(winSiteConfig.ApplicationStack) == 1 {
			winAppStack := winSiteConfig.ApplicationStack[0]
			if winAppStack.NetFrameworkVersion != "" {
				expanded.NetFrameworkVersion = pointer.To(winAppStack.NetFrameworkVersion)
				if currentStack == "" {
					currentStack = CurrentStackDotNet
				}
			}
			if winAppStack.NetCoreVersion != "" {
				expanded.NetFrameworkVersion = pointer.To(winAppStack.NetCoreVersion)
				if currentStack == "" {
					currentStack = CurrentStackDotNetCore
				}
			}
			if winAppStack.NodeVersion != "" {
				// Note: node version is now exclusively controlled via app_setting.WEBSITE_NODE_DEFAULT_VERSION
				if currentStack == "" {
					currentStack = CurrentStackNode
				}
			}
			if winAppStack.PhpVersion != "" {
				if winAppStack.PhpVersion != PhpVersionOff {
					expanded.PhpVersion = pointer.To(winAppStack.PhpVersion)
				} else {
					expanded.PhpVersion = pointer.To("")
				}
				if currentStack == "" {
					currentStack = CurrentStackPhp
				}
			}
			if winAppStack.PythonVersion != "" || winAppStack.Python {
				expanded.PythonVersion = pointer.To(winAppStack.PythonVersion)
				if currentStack == "" {
					currentStack = CurrentStackPython
				}
			}
			if winAppStack.JavaVersion != "" {
				expanded.JavaVersion = pointer.To(winAppStack.JavaVersion)
				switch {
				case winAppStack.JavaEmbeddedServer:
					expanded.JavaContainer = pointer.To(JavaContainerEmbeddedServer)
					expanded.JavaContainerVersion = pointer.To(JavaContainerEmbeddedServerVersion)
				case winAppStack.TomcatVersion != "":
					expanded.JavaContainer = pointer.To(JavaContainerTomcat)
					expanded.JavaContainerVersion = pointer.To(winAppStack.TomcatVersion)
				case winAppStack.JavaContainer != "":
					expanded.JavaContainer = pointer.To(winAppStack.JavaContainer)
					expanded.JavaContainerVersion = pointer.To(winAppStack.JavaContainerVersion)
				}

				if currentStack == "" {
					currentStack = CurrentStackJava
				}
			}
			if winAppStack.DockerContainerName != "" || winAppStack.DockerContainerRegistry != "" || winAppStack.DockerContainerTag != "" {
				if winAppStack.DockerContainerRegistry != "" {
					expanded.WindowsFxVersion = pointer.To(fmt.Sprintf("DOCKER|%s/%s:%s", winAppStack.DockerContainerRegistry, winAppStack.DockerContainerName, winAppStack.DockerContainerTag))
				} else {
					expanded.WindowsFxVersion = pointer.To(fmt.Sprintf("DOCKER|%s:%s", winAppStack.DockerContainerName, winAppStack.DockerContainerTag))
				}
			}
		} else {
			expanded.WindowsFxVersion = pointer.To("")
		}
	}

	if metadata.ResourceData.HasChange("site_config.0.virtual_application") {
		expanded.VirtualApplications = expandVirtualApplicationsForUpdate(winSiteConfig.VirtualApplications)
	} else {
		expanded.VirtualApplications = expandVirtualApplications(winSiteConfig.VirtualApplications)
	}

	expanded.AcrUseManagedIdentityCreds = pointer.To(winSiteConfig.UseManagedIdentityACR)

	if metadata.ResourceData.HasChange("site_config.0.container_registry_managed_identity_client_id") {
		expanded.AcrUserManagedIdentityID = pointer.To(winSiteConfig.ContainerRegistryUserMSI)
	}

	if metadata.ResourceData.HasChange("site_config.0.default_documents") {
		expanded.DefaultDocuments = &winSiteConfig.DefaultDocuments
	}

	expanded.HTTP20Enabled = pointer.To(winSiteConfig.Http2Enabled)

	if metadata.ResourceData.HasChange("site_config.0.ip_restriction") {
		ipRestrictions, err := ExpandIpRestrictions(winSiteConfig.IpRestriction)
		if err != nil {
			return nil, nil, fmt.Errorf("expanding IP Restrictions: %+v", err)
		}
		expanded.IPSecurityRestrictions = ipRestrictions
	}

	expanded.ScmIPSecurityRestrictionsUseMain = pointer.To(winSiteConfig.ScmUseMainIpRestriction)

	if metadata.ResourceData.HasChange("site_config.0.scm_ip_restriction") {
		scmIpRestrictions, err := ExpandIpRestrictions(winSiteConfig.ScmIpRestriction)
		if err != nil {
			return nil, nil, fmt.Errorf("expanding SCM IP Restrictions: %+v", err)
		}
		expanded.ScmIPSecurityRestrictions = scmIpRestrictions
	}

	expanded.LocalMySQLEnabled = pointer.To(winSiteConfig.LocalMysql)

	if metadata.ResourceData.HasChange("site_config.0.load_balancing_mode") {
		expanded.LoadBalancing = web.SiteLoadBalancing(winSiteConfig.LoadBalancing)
	}

	if metadata.ResourceData.HasChange("site_config.0.managed_pipeline_mode") {
		expanded.ManagedPipelineMode = web.ManagedPipelineMode(winSiteConfig.ManagedPipelineMode)
	}

	expanded.RemoteDebuggingEnabled = pointer.To(winSiteConfig.RemoteDebugging)

	if metadata.ResourceData.HasChange("site_config.0.remote_debugging_version") {
		expanded.RemoteDebuggingVersion = pointer.To(winSiteConfig.RemoteDebuggingVersion)
	}

	expanded.Use32BitWorkerProcess = pointer.To(winSiteConfig.Use32BitWorker)

	expanded.WebSocketsEnabled = pointer.To(winSiteConfig.WebSockets)

	if metadata.ResourceData.HasChange("site_config.0.ftps_state") {
		expanded.FtpsState = web.FtpsState(winSiteConfig.FtpsState)
	}

	if metadata.ResourceData.HasChange("site_config.0.health_check_path") {
		expanded.HealthCheckPath = pointer.To(winSiteConfig.HealthCheckPath)
	}

	if metadata.ResourceData.HasChange("site_config.0.worker_count") {
		expanded.NumberOfWorkers = pointer.To(int32(winSiteConfig.WorkerCount))
	}

	if metadata.ResourceData.HasChange("site_config.0.minimum_tls_version") {
		expanded.MinTLSVersion = web.SupportedTLSVersions(winSiteConfig.MinTlsVersion)
	}

	if metadata.ResourceData.HasChange("site_config.0.scm_minimum_tls_version") {
		expanded.ScmMinTLSVersion = web.SupportedTLSVersions(winSiteConfig.ScmMinTlsVersion)
	}

	if metadata.ResourceData.HasChange("site_config.0.cors") {
		cors := ExpandCorsSettings(winSiteConfig.Cors)
		if cors == nil {
			cors = &web.CorsSettings{
				AllowedOrigins: &[]string{},
			}
		}
		expanded.Cors = cors
	}

	if metadata.ResourceData.HasChange("site_config.0.auto_heal_enabled") {
		expanded.AutoHealEnabled = pointer.To(winSiteConfig.AutoHeal)
	}

	if metadata.ResourceData.HasChange("site_config.0.auto_heal_setting") {
		expanded.AutoHealRules = expandAutoHealSettingsWindows(winSiteConfig.AutoHealSettings)
	}

	if metadata.ResourceData.HasChange("site_config.0.vnet_route_all_enabled") {
		expanded.VnetRouteAllEnabled = pointer.To(winSiteConfig.VnetRouteAllEnabled)
	}

	return expanded, &currentStack, nil
}

func FlattenSiteConfigWindows(appSiteConfig *web.SiteConfig, currentStack string, healthCheckCount *int) ([]SiteConfigWindows, error) {
	if appSiteConfig == nil {
		return nil, nil
	}

	siteConfig := SiteConfigWindows{
		AlwaysOn:                 pointer.From(appSiteConfig.AlwaysOn),
		AppCommandLine:           pointer.From(appSiteConfig.AppCommandLine),
		AutoHeal:                 pointer.From(appSiteConfig.AutoHealEnabled),
		AutoHealSettings:         flattenAutoHealSettingsWindows(appSiteConfig.AutoHealRules),
		ContainerRegistryUserMSI: pointer.From(appSiteConfig.AcrUserManagedIdentityID),
		DetailedErrorLogging:     pointer.From(appSiteConfig.DetailedErrorLoggingEnabled),
		FtpsState:                string(appSiteConfig.FtpsState),
		HealthCheckPath:          pointer.From(appSiteConfig.HealthCheckPath),
		HealthCheckEvictionTime:  pointer.From(healthCheckCount),
		Http2Enabled:             pointer.From(appSiteConfig.HTTP20Enabled),
		IpRestriction:            FlattenIpRestrictions(appSiteConfig.IPSecurityRestrictions),
		LoadBalancing:            string(appSiteConfig.LoadBalancing),
		LocalMysql:               pointer.From(appSiteConfig.LocalMySQLEnabled),
		ManagedPipelineMode:      string(appSiteConfig.ManagedPipelineMode),
		MinTlsVersion:            string(appSiteConfig.MinTLSVersion),
		WorkerCount:              int(pointer.From(appSiteConfig.NumberOfWorkers)),
		RemoteDebugging:          pointer.From(appSiteConfig.RemoteDebuggingEnabled),
		RemoteDebuggingVersion:   strings.ToUpper(pointer.From(appSiteConfig.RemoteDebuggingVersion)),
		ScmIpRestriction:         FlattenIpRestrictions(appSiteConfig.ScmIPSecurityRestrictions),
		ScmMinTlsVersion:         string(appSiteConfig.ScmMinTLSVersion),
		ScmType:                  string(appSiteConfig.ScmType),
		ScmUseMainIpRestriction:  pointer.From(appSiteConfig.ScmIPSecurityRestrictionsUseMain),
		Use32BitWorker:           pointer.From(appSiteConfig.Use32BitWorkerProcess),
		UseManagedIdentityACR:    pointer.From(appSiteConfig.AcrUseManagedIdentityCreds),
		VirtualApplications:      flattenVirtualApplications(appSiteConfig.VirtualApplications),
		WebSockets:               pointer.From(appSiteConfig.WebSocketsEnabled),
		VnetRouteAllEnabled:      pointer.From(appSiteConfig.VnetRouteAllEnabled),
	}

	if appSiteConfig.APIManagementConfig != nil && appSiteConfig.APIManagementConfig.ID != nil {
		apiId, err := parse.ApiIDInsensitively(*appSiteConfig.APIManagementConfig.ID)
		if err != nil {
			return nil, fmt.Errorf("could not parse API Management ID: %+v", err)
		}
		siteConfig.ApiManagementConfigId = apiId.ID()
	}

	if appSiteConfig.APIDefinition != nil && appSiteConfig.APIDefinition.URL != nil {
		siteConfig.ApiDefinition = *appSiteConfig.APIDefinition.URL
	}

	if appSiteConfig.DefaultDocuments != nil {
		siteConfig.DefaultDocuments = *appSiteConfig.DefaultDocuments
	}

	if appSiteConfig.NumberOfWorkers != nil {
		siteConfig.WorkerCount = int(*appSiteConfig.NumberOfWorkers)
	}

	var winAppStack ApplicationStackWindows
	if currentStack == CurrentStackDotNetCore {
		winAppStack.NetCoreVersion = pointer.From(appSiteConfig.NetFrameworkVersion)
	}
	if currentStack == CurrentStackDotNet {
		winAppStack.NetFrameworkVersion = pointer.From(appSiteConfig.NetFrameworkVersion)
	}
	winAppStack.PhpVersion = pointer.From(appSiteConfig.PhpVersion)
	if winAppStack.PhpVersion == "" {
		winAppStack.PhpVersion = PhpVersionOff
	}
	winAppStack.PythonVersion = pointer.From(appSiteConfig.PythonVersion) // This _should_ always be `""`
	winAppStack.Python = currentStack == CurrentStackPython
	winAppStack.JavaVersion = pointer.From(appSiteConfig.JavaVersion)
	switch pointer.From(appSiteConfig.JavaContainer) {
	case JavaContainerTomcat:
		winAppStack.TomcatVersion = *appSiteConfig.JavaContainerVersion
	case JavaContainerEmbeddedServer:
		winAppStack.JavaEmbeddedServer = true
	}

	siteConfig.WindowsFxVersion = pointer.From(appSiteConfig.WindowsFxVersion)
	if siteConfig.WindowsFxVersion != "" {
		// Decode the string to docker values
		parts := strings.Split(strings.TrimPrefix(siteConfig.WindowsFxVersion, "DOCKER|"), ":")
		if len(parts) == 2 {
			winAppStack.DockerContainerTag = parts[1]
			path := strings.Split(parts[0], "/")
			if len(path) > 1 {
				winAppStack.DockerContainerRegistry = path[0]
				winAppStack.DockerContainerName = strings.TrimPrefix(parts[0], fmt.Sprintf("%s/", path[0]))
			} else {
				winAppStack.DockerContainerName = path[0]
			}
		}
	}
	winAppStack.CurrentStack = currentStack

	siteConfig.ApplicationStack = []ApplicationStackWindows{winAppStack}

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

	return []SiteConfigWindows{siteConfig}, nil
}
