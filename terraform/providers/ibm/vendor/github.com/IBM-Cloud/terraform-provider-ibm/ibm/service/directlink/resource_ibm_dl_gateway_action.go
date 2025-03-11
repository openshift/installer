// Copyright IBM Corp. 2023 All Rights Reserved.
// Licensed under the Mozilla Public License v2.0

package directlink

import (
	"context"
	"fmt"
	"log"
	"os"
	"reflect"
	"time"

	"github.com/IBM/networking-go-sdk/directlinkv1"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/customdiff"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/flex"
	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/validate"
)

func ResourceIBMDLGatewayAction() *schema.Resource {
	return &schema.Resource{
		Create:   resourceIBMdlGatewayCreateAction,
		Read:     resourceIBMdlGatewayActionRead,
		Update:   resourceIBMdlGatewayActionUpdate,
		Delete:   resourceIBMdlGatewayActionDelete,
		Importer: &schema.ResourceImporter{},

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(60 * time.Minute),
			Delete: schema.DefaultTimeout(60 * time.Minute),
			Update: schema.DefaultTimeout(60 * time.Minute),
		},

		CustomizeDiff: customdiff.Sequence(
			func(_ context.Context, diff *schema.ResourceDiff, v interface{}) error {
				return flex.ResourceTagsCustomizeDiff(diff)
			},
		),
		Schema: map[string]*schema.Schema{
			dlGatewayId: {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The Direct Link gateway identifier",
			},
			dlCustomerAction: {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validate.InvokeValidator("ibm_dl_gateway_action", dlCustomerAction),
				Description:  "customer action on provider call",
			},
			dlAuthenticationKey: {
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    false,
				Sensitive:   true,
				Description: "BGP MD5 authentication key",
			},
			dlExportRouteFilters: {
				Type:        schema.TypeList,
				Optional:    true,
				ForceNew:    false,
				Description: "List Export Route Filters for a Direct Link gateway",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						dlExportRouteFilterId: {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Export route Filter identifier",
						},
						dlAction: {
							Type:         schema.TypeString,
							Required:     true,
							ValidateFunc: validate.InvokeValidator("ibm_dl_gateway_action", dlAction),
							Description:  "Determines whether the  routes that match the prefix-set will be permit or deny",
						},
						dlBefore: {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Identifier of the next route filter to be considered",
						},
						dlCreatedAt: {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The date and time of the export route filter was created",
						},
						dlGe: {
							Type:        schema.TypeInt,
							Optional:    true,
							Description: "The minimum matching length of the prefix-set",
						},
						dlLe: {
							Type:        schema.TypeInt,
							Optional:    true,
							Description: "The maximum matching length of the prefix-set",
						},
						dlPrefix: {
							Type:        schema.TypeString,
							Required:    true,
							Description: "IP prefix representing an address and mask length of the prefix-set",
						},
						dlUpdatedAt: {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The date and time of the export route filter was last updated",
						},
					},
				},
			},
			dlImportRouteFilters: {
				Type:        schema.TypeList,
				Optional:    true,
				ForceNew:    false,
				Description: "List Import Route Filters for a Direct Link gateway",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						dlImportRouteFilterId: {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Import route Filter identifier",
						},
						dlAction: {
							Type:         schema.TypeString,
							Required:     true,
							ValidateFunc: validate.InvokeValidator("ibm_dl_gateway_action", dlAction),
							Description:  "Determines whether the  routes that match the prefix-set will be permit or deny",
						},
						dlBefore: {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Identifier of the next route filter to be considered",
						},
						dlCreatedAt: {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The date and time of the export route filter was created",
						},
						dlGe: {
							Type:        schema.TypeInt,
							Optional:    true,
							Description: "The minimum matching length of the prefix-set",
						},
						dlLe: {
							Type:        schema.TypeInt,
							Optional:    true,
							Description: "The maximum matching length of the prefix-set",
						},
						dlPrefix: {
							Type:        schema.TypeString,
							Required:    true,
							Description: "IP prefix representing an address and mask length of the prefix-set",
						},
						dlUpdatedAt: {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The date and time of the export route filter was last updated",
						},
					},
				},
			},
			dlDefault_export_route_filter: {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validate.InvokeValidator("ibm_dl_gateway_action", dlDefault_export_route_filter),
				Description:  "The default directional route filter action that applies to routes that do not match any directional route filters",
			},
			dlDefault_import_route_filter: {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validate.InvokeValidator("ibm_dl_gateway_action", dlDefault_import_route_filter),
				Description:  "The default directional route filter action that applies to routes that do not match any directional route filters",
			},
			dlAsPrepends: {
				Type:        schema.TypeList,
				Optional:    true,
				ForceNew:    false,
				Description: "List of AS Prepend configuration information",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						dlCreatedAt: {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The date and time AS Prepend was created",
						},
						dlResourceId: {
							Type:        schema.TypeString,
							Optional:    true,
							ForceNew:    false,
							Computed:    true,
							Description: "The unique identifier for this AS Prepend",
						},
						dlLength: {
							Type:         schema.TypeInt,
							Required:     true,
							ForceNew:     false,
							ValidateFunc: validate.InvokeValidator("ibm_dl_gateway_action", dlLength),
							Description:  "Number of times the ASN to appended to the AS Path",
						},
						dlPolicy: {
							Type:         schema.TypeString,
							Required:     true,
							ForceNew:     false,
							ValidateFunc: validate.InvokeValidator("ibm_dl_gateway_action", dlPolicy),
							Description:  "Route type this AS Prepend applies to",
						},
						dlPrefix: {
							Type:        schema.TypeString,
							Optional:    true,
							ForceNew:    false,
							Description: "Comma separated list of prefixes this AS Prepend applies to. Maximum of 10 prefixes. If not specified, this AS Prepend applies to all prefixes",
							Deprecated:  "prefix will be deprecated and support will be removed. Use specific_prefixes instead",
						},
						dlSpecificPrefixes: {
							Type:        schema.TypeList,
							Description: "Array of prefixes this AS Prepend applies to",
							Optional:    true,
							ForceNew:    false,
							MinItems:    1,
							MaxItems:    10,
							Elem:        &schema.Schema{Type: schema.TypeString},
						},
						dlUpdatedAt: {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The date and time AS Prepend was updated",
						},
					},
				},
			},
			dlBfdInterval: {
				Type:         schema.TypeInt,
				Optional:     true,
				ForceNew:     false,
				Description:  "BFD Interval",
				ValidateFunc: validate.InvokeValidator("ibm_dl_gateway_action", dlBfdInterval),
			},
			dlBfdMultiplier: {
				Type:         schema.TypeInt,
				Optional:     true,
				ForceNew:     false,
				Description:  "BFD Multiplier",
				ValidateFunc: validate.InvokeValidator("ibm_dl_gateway_action", dlBfdMultiplier),
			},
			dlBfdStatus: {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Gateway BFD status",
			},
			dlBfdStatusUpdatedAt: {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Date and time BFD status was updated",
			},
			dlBgpAsn: {
				Type:        schema.TypeInt,
				Optional:    true,
				Description: "BGP ASN",
			},
			dlBgpBaseCidr: {
				Type:             schema.TypeString,
				Optional:         true,
				ForceNew:         false,
				DiffSuppressFunc: flex.ApplyOnce,
				Description:      "BGP base CIDR",
			},
			dlPort: {
				Type:          schema.TypeString,
				Optional:      true,
				ForceNew:      true,
				Computed:      true,
				Description:   "Gateway port",
				ConflictsWith: []string{"location_name", "cross_connect_router", "carrier_name", "customer_name"},
			},
			dlConnectionMode: {
				Type:         schema.TypeString,
				Computed:     true,
				Optional:     true,
				ForceNew:     false,
				Description:  "Type of services this Gateway is attached to. Mode transit means this Gateway will be attached to Transit Gateway Service and direct means this Gateway will be attached to vpc or classic connection",
				ValidateFunc: validate.InvokeValidator("ibm_dl_gateway_action", dlConnectionMode),
			},
			dlCrossConnectRouter: {
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				Description: "Cross connect router",
			},
			dlGlobal: {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: "Gateways with global routing (true) can connect to networks outside their associated region",
			},
			dlLocationName: {
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				Computed:    true,
				Description: "Gateway location",
			},
			dlMetered: {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: "Metered billing option",
			},
			dlName: {
				Type:         schema.TypeString,
				Computed:     true,
				Optional:     true,
				Description:  "The unique user-defined name for this gateway",
				ValidateFunc: validate.InvokeValidator("ibm_dl_gateway_action", dlName),
			},
			dlType: {
				Type:         schema.TypeString,
				Computed:     true,
				Optional:     true,
				Description:  "Gateway type",
				ValidateFunc: validate.InvokeValidator("ibm_dl_gateway_action", dlType),
			},
			dlCarrierName: {
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				Description: "Carrier name",
			},
			dlCustomerName: {
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				Description: "Customer name",
			},
			dlSpeedMbps: {
				Type:        schema.TypeInt,
				Optional:    true,
				Description: "Gateway speed in megabits per second",
			},
			dlBgpCerCidr: {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "BGP customer edge router CIDR",
			},
			dlLoaRejectReason: {
				Type:        schema.TypeString,
				Computed:    true,
				Optional:    true,
				ForceNew:    false,
				Description: "Loa reject reason",
			},
			dlBgpIbmCidr: {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "BGP IBM CIDR",
			},
			dlResourceGroup: {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				ForceNew:    true,
				Description: "Gateway resource group",
			},
			dlOperationalStatus: {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Gateway operational status",
			},
			dlProviderAPIManaged: {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "Indicates whether gateway was created through a provider portal",
			},
			dlVlan: {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "VLAN allocated for this gateway",
			},
			dlBgpIbmAsn: {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "IBM BGP ASN",
			},
			dlBgpStatus: {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Gateway BGP status",
			},
			dlBgpStatusUpdatedAt: {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Date and time BGP status was updated",
			},
			dlChangeRequest: {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Changes pending approval for provider managed Direct Link Connect gateways",
			},
			dlCompletionNoticeRejectReason: {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Reason for completion notice rejection",
			},
			dlCreatedAt: {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The date and time resource was created",
			},
			dlCrn: {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The CRN (Cloud Resource Name) of this gateway",
			},
			dlLinkStatus: {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Gateway link status",
			},
			dlLinkStatusUpdatedAt: {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Date and time Link status was updated",
			},
			dlLocationDisplayName: {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Gateway location long name",
			},
			dlTags: {
				Type:        schema.TypeSet,
				Optional:    true,
				Computed:    true,
				Elem:        &schema.Schema{Type: schema.TypeString, ValidateFunc: validate.InvokeValidator("ibm_dl_gateway_action", dlTags)},
				Set:         flex.ResourceIBMVPCHash,
				Description: "Tags for the direct link gateway",
			},
			flex.ResourceControllerURL: {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The URL of the IBM Cloud dashboard that can be used to explore and view details about this instance",
			},
			flex.ResourceName: {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The name of the resource",
			},
			flex.ResourceCRN: {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The crn of the resource",
			},
			flex.ResourceStatus: {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The status of the resource",
			},
			flex.ResourceGroupName: {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The resource group name in which resource is provisioned",
			},
		},
	}
}

func ResourceIBMDLGatewayActionValidator() *validate.ResourceValidator {

	validateSchema := make([]validate.ValidateSchema, 0)
	dlTypeAllowedValues := "dedicated, connect"
	dlConnectionModeAllowedValues := "direct, transit"
	dlPolicyAllowedValues := "export, import"
	dlActionValues := "permit, deny"
	dlCustomerActionValues := "create_gateway_approve, create_gateway_reject, delete_gateway_approve, delete_gateway_reject, update_attributes_approve, update_attributes_reject"
	validateSchema = append(validateSchema,
		validate.ValidateSchema{
			Identifier:                 dlCustomerAction,
			ValidateFunctionIdentifier: validate.ValidateAllowedStringValue,
			Type:                       validate.TypeString,
			Required:                   true,
			AllowedValues:              dlCustomerActionValues})
	validateSchema = append(validateSchema,
		validate.ValidateSchema{
			Identifier:                 dlAction,
			ValidateFunctionIdentifier: validate.ValidateAllowedStringValue,
			Type:                       validate.TypeString,
			Required:                   true,
			AllowedValues:              dlActionValues})
	validateSchema = append(validateSchema,
		validate.ValidateSchema{
			Identifier:                 dlDefault_export_route_filter,
			ValidateFunctionIdentifier: validate.ValidateAllowedStringValue,
			Type:                       validate.TypeString,
			Required:                   true,
			AllowedValues:              dlActionValues})
	validateSchema = append(validateSchema,
		validate.ValidateSchema{
			Identifier:                 dlDefault_import_route_filter,
			ValidateFunctionIdentifier: validate.ValidateAllowedStringValue,
			Type:                       validate.TypeString,
			Required:                   true,
			AllowedValues:              dlActionValues})
	validateSchema = append(validateSchema,
		validate.ValidateSchema{
			Identifier:                 dlType,
			ValidateFunctionIdentifier: validate.ValidateAllowedStringValue,
			Type:                       validate.TypeString,
			Required:                   true,
			AllowedValues:              dlTypeAllowedValues})
	validateSchema = append(validateSchema,
		validate.ValidateSchema{
			Identifier:                 dlName,
			ValidateFunctionIdentifier: validate.ValidateRegexpLen,
			Type:                       validate.TypeString,
			Required:                   true,
			Regexp:                     `^([a-zA-Z]|[a-zA-Z][-_a-zA-Z0-9]*[a-zA-Z0-9])$`,
			MinValueLength:             1,
			MaxValueLength:             63})
	validateSchema = append(validateSchema,
		validate.ValidateSchema{
			Identifier:                 "tags",
			ValidateFunctionIdentifier: validate.ValidateRegexpLen,
			Type:                       validate.TypeString,
			Optional:                   true,
			Regexp:                     `^[A-Za-z0-9:_ .-]+$`,
			MinValueLength:             1,
			MaxValueLength:             128})
	validateSchema = append(validateSchema,
		validate.ValidateSchema{
			Identifier:                 dlConnectionMode,
			ValidateFunctionIdentifier: validate.ValidateAllowedStringValue,
			Type:                       validate.TypeString,
			Required:                   true,
			AllowedValues:              dlConnectionModeAllowedValues})
	validateSchema = append(validateSchema,
		validate.ValidateSchema{
			Identifier:                 dlBfdInterval,
			ValidateFunctionIdentifier: validate.IntBetween,
			Type:                       validate.TypeInt,
			Required:                   true,
			MinValue:                   "300",
			MaxValue:                   "255000"})
	validateSchema = append(validateSchema,
		validate.ValidateSchema{
			Identifier:                 dlBfdMultiplier,
			ValidateFunctionIdentifier: validate.IntBetween,
			Type:                       validate.TypeInt,
			Required:                   true,
			MinValue:                   "1",
			MaxValue:                   "255"})
	validateSchema = append(validateSchema,
		validate.ValidateSchema{
			Identifier:                 dlPolicy,
			ValidateFunctionIdentifier: validate.ValidateAllowedStringValue,
			Type:                       validate.TypeString,
			Required:                   true,
			AllowedValues:              dlPolicyAllowedValues})

	validateSchema = append(validateSchema,
		validate.ValidateSchema{
			Identifier:                 dlLength,
			ValidateFunctionIdentifier: validate.IntBetween,
			Type:                       validate.TypeInt,
			Required:                   true,
			MinValue:                   "3",
			MaxValue:                   "10"})

	ibmISDLGatewayResourceValidator := validate.ResourceValidator{ResourceName: "ibm_dl_gateway", Schema: validateSchema}
	return &ibmISDLGatewayResourceValidator
}

func resourceIBMdlGatewayCreateAction(d *schema.ResourceData, meta interface{}) error {
	directLink, err := directlinkClient(meta)
	if err != nil {
		return err
	}

	createGatewayActionOptionsModel := &directlinkv1.CreateGatewayActionOptions{}
	createGatewayActionOptionsModel.Action = NewStrPointer(d.Get(dlCustomerAction).(string))
	createGatewayActionOptionsModel.ID = NewStrPointer(d.Get(dlGatewayId).(string))
	global := d.Get(dlGlobal).(bool)
	createGatewayActionOptionsModel.Global = &global
	metered := d.Get(dlMetered).(bool)
	createGatewayActionOptionsModel.Metered = &metered

	var bfdConfig directlinkv1.GatewayBfdConfigActionTemplate
	isBfdInterval := false
	if bfdInterval, ok := d.GetOk(dlBfdInterval); ok {
		isBfdInterval = true
		bfdConfig.Interval = NewInt64Pointer(int64(bfdInterval.(int)))
	}
	if bfdMultiplier, ok := d.GetOk(dlBfdMultiplier); ok {
		bfdConfig.Multiplier = NewInt64Pointer(int64(bfdMultiplier.(int)))
	} else if isBfdInterval {
		// Set the default value for multiplier if interval is set
		multiplier := int64(3)
		bfdConfig.Multiplier = &multiplier
	}

	asPrependsCreateItems := make([]directlinkv1.AsPrependTemplate, 0)
	if asPrependsInput, ok := d.GetOk(dlAsPrepends); ok {
		asPrependsItems := asPrependsInput.([]interface{})

		for _, asPrependItem := range asPrependsItems {
			i := asPrependItem.(map[string]interface{})

			// Construct an instance of the AsPrependTemplate model
			asPrependTemplateModel := new(directlinkv1.AsPrependTemplate)
			asPrependTemplateModel.Length = NewInt64Pointer(int64(i[dlLength].(int)))
			asPrependTemplateModel.Policy = NewStrPointer(i[dlPolicy].(string))
			asPrependTemplateModel.Prefix = nil
			asPrependTemplateModel.SpecificPrefixes = nil
			_, prefix_ok := i[dlPrefix]
			if prefix_ok && (len(i[dlPrefix].(string)) > 0) {
				asPrependTemplateModel.Prefix = NewStrPointer(i[dlPrefix].(string))
				asPrependTemplateModel.SpecificPrefixes = nil
			}

			sp_prefixOk, ok := i[dlSpecificPrefixes]
			if ok && len(sp_prefixOk.([]interface{})) > 0 {
				asPrependTemplateModel.Prefix = nil
				asPrependTemplateModel.SpecificPrefixes = flex.ExpandStringList(sp_prefixOk.([]interface{}))
			}
			asPrependsCreateItems = append(asPrependsCreateItems, *asPrependTemplateModel)
		}
	}
	exportRouteFiltersCreateList := make([]directlinkv1.GatewayTemplateRouteFilter, 0)
	if exportRouteFiltersInputList, ok := d.GetOk(dlExportRouteFilters); ok {
		exportRouteFilters := exportRouteFiltersInputList.([]interface{})

		for _, exportRouteFilter := range exportRouteFilters {
			filtersData := exportRouteFilter.(map[string]interface{})

			// Construct an Export Route Filters List
			exportRouteFilterTemplateModel := new(directlinkv1.GatewayTemplateRouteFilter)
			exportRouteFilterTemplateModel.Action = NewStrPointer(filtersData[dlAction].(string))
			exportRouteFilterTemplateModel.Prefix = NewStrPointer(filtersData[dlPrefix].(string))
			exportRouteFilterTemplateModel.Ge = nil
			exportRouteFilterTemplateModel.Le = nil
			if _, ok := filtersData[dlGe]; ok {
				exportRouteFilterTemplateModel.Ge = NewInt64Pointer(int64(filtersData[dlGe].(int)))
			}
			if _, ok := filtersData[dlLe]; ok {
				exportRouteFilterTemplateModel.Le = NewInt64Pointer(int64(filtersData[dlLe].(int)))
			}
			exportRouteFiltersCreateList = append(exportRouteFiltersCreateList, *exportRouteFilterTemplateModel)
		}
	}
	importRouteFiltersCreateList := make([]directlinkv1.GatewayTemplateRouteFilter, 0)
	if importRouteFiltersInputList, ok := d.GetOk(dlImportRouteFilters); ok {
		importRouteFilters := importRouteFiltersInputList.([]interface{})

		for _, importRouteFilter := range importRouteFilters {
			filtersData := importRouteFilter.(map[string]interface{})

			// Construct an Import Route Filters List
			importRouteFilterTemplateModel := new(directlinkv1.GatewayTemplateRouteFilter)
			importRouteFilterTemplateModel.Action = NewStrPointer(filtersData[dlAction].(string))
			importRouteFilterTemplateModel.Prefix = NewStrPointer(filtersData[dlPrefix].(string))
			importRouteFilterTemplateModel.Ge = nil
			importRouteFilterTemplateModel.Le = nil
			if _, ok := filtersData[dlGe]; ok {
				importRouteFilterTemplateModel.Ge = NewInt64Pointer(int64(filtersData[dlGe].(int)))
			}
			if _, ok := filtersData[dlLe]; ok {
				importRouteFilterTemplateModel.Le = NewInt64Pointer(int64(filtersData[dlLe].(int)))
			}
			importRouteFiltersCreateList = append(importRouteFiltersCreateList, *importRouteFilterTemplateModel)
		}
	}
	if _, ok := d.GetOk(dlResourceGroup); ok {
		resourceGroup := d.Get(dlResourceGroup).(string)
		createGatewayActionOptionsModel.ResourceGroup = &directlinkv1.ResourceGroupIdentity{ID: &resourceGroup}
	}
	if authKeyCrn, ok := d.GetOk(dlAuthenticationKey); ok {
		authKeyCrnStr := authKeyCrn.(string)
		createGatewayActionOptionsModel.AuthenticationKey = &directlinkv1.GatewayActionTemplateAuthenticationKey{Crn: &authKeyCrnStr}
	}
	if connectionMode, ok := d.GetOk(dlConnectionMode); ok {
		createGatewayActionOptionsModel.ConnectionMode = NewStrPointer(connectionMode.(string))
	}
	if !reflect.DeepEqual(bfdConfig, directlinkv1.GatewayBfdConfigActionTemplate{}) {
		createGatewayActionOptionsModel.BfdConfig = &bfdConfig
	}
	if len(asPrependsCreateItems) > 0 {
		createGatewayActionOptionsModel.AsPrepends = asPrependsCreateItems
	}
	if len(exportRouteFiltersCreateList) > 0 {
		createGatewayActionOptionsModel.ExportRouteFilters = exportRouteFiltersCreateList
	}
	if len(importRouteFiltersCreateList) > 0 {
		createGatewayActionOptionsModel.ImportRouteFilters = importRouteFiltersCreateList
	}
	if default_export_route_filter, ok := d.GetOk(dlDefault_export_route_filter); ok {
		createGatewayActionOptionsModel.DefaultExportRouteFilter = NewStrPointer(default_export_route_filter.(string))
	}
	if default_import_route_filter, ok := d.GetOk(dlDefault_import_route_filter); ok {
		createGatewayActionOptionsModel.DefaultImportRouteFilter = NewStrPointer(default_import_route_filter.(string))
	}
	gateway, response, err := directLink.CreateGatewayAction(createGatewayActionOptionsModel)
	if err != nil {
		return fmt.Errorf("[DEBUG] Direct Link Gateway Action err %s\n%s", err, response)
	}

	d.SetId(*gateway.ID)
	_, err = isWaitForDirectLinkAvailableforAction(directLink, d.Id(), d.Timeout(schema.TimeoutCreate))
	if err != nil {
		return err
	}
	v := os.Getenv("IC_ENV_TAGS")
	if _, ok := d.GetOk(dlTags); ok || v != "" {
		oldList, newList := d.GetChange(dlTags)
		err = flex.UpdateTagsUsingCRN(oldList, newList, meta, *gateway.Crn)
		if err != nil {
			log.Printf(
				"Error on create of resource direct link gateway %s tags: %s", d.Id(), err)
		}
	}
	return resourceIBMdlGatewayActionRead(d, meta)
}
func resourceIBMdlGatewayActionExportRouteFiltersRead(d *schema.ResourceData, meta interface{}) error {
	directLink, err := directlinkClient(meta)
	if err != nil {
		return err
	}

	gatewayId := d.Id()
	listGatewayExportRouteFiltersOptionsModel := &directlinkv1.ListGatewayExportRouteFiltersOptions{GatewayID: &gatewayId}
	exportRouteFilterList, response, err := directLink.ListGatewayExportRouteFilters(listGatewayExportRouteFiltersOptionsModel)
	if err != nil {
		log.Println("[WARN] Error listing Direct Link Export Route Filters", response, err)
		return err
	}
	exportRouteFilters := make([]map[string]interface{}, 0)
	for _, instance := range exportRouteFilterList.ExportRouteFilters {
		routeFilter := map[string]interface{}{}
		if instance.ID != nil {
			routeFilter[dlExportRouteFilterId] = *instance.ID
		}
		if instance.Action != nil {
			routeFilter[dlAction] = *instance.Action
		}
		if instance.Before != nil {
			routeFilter[dlBefore] = *instance.Before
		}
		if instance.CreatedAt != nil {
			routeFilter[dlCreatedAt] = instance.CreatedAt.String()
		}
		if instance.Prefix != nil {
			routeFilter[dlPrefix] = *instance.Prefix
		}
		if instance.UpdatedAt != nil {
			routeFilter[dlUpdatedAt] = instance.UpdatedAt.String()
		}
		if instance.Ge != nil {
			routeFilter[dlGe] = *instance.Ge
		}
		if instance.Le != nil {
			routeFilter[dlLe] = *instance.Le
		}
		exportRouteFilters = append(exportRouteFilters, routeFilter)
	}
	d.Set(dlExportRouteFilters, exportRouteFilters)
	return nil
}

func resourceIBMdlGatewayActionImportRouteFiltersRead(d *schema.ResourceData, meta interface{}) error {
	directLink, err := directlinkClient(meta)
	if err != nil {
		return err
	}

	gatewayId := d.Id()
	listGatewayImportRouteFiltersOptionsModel := &directlinkv1.ListGatewayImportRouteFiltersOptions{GatewayID: &gatewayId}
	importRouteFilterList, response, err := directLink.ListGatewayImportRouteFilters(listGatewayImportRouteFiltersOptionsModel)
	if err != nil {
		log.Println("[WARN] Error  while listing Direct Link Import Route Filters", response, err)
		return err
	}
	importRouteFilters := make([]map[string]interface{}, 0)
	for _, instance := range importRouteFilterList.ImportRouteFilters {
		routeFilter := map[string]interface{}{}
		if instance.ID != nil {
			routeFilter[dlImportRouteFilterId] = *instance.ID
		}
		if instance.Action != nil {
			routeFilter[dlAction] = *instance.Action
		}
		if instance.Before != nil {
			routeFilter[dlBefore] = *instance.Before
		}
		if instance.CreatedAt != nil {
			routeFilter[dlCreatedAt] = instance.CreatedAt.String()
		}
		if instance.Prefix != nil {
			routeFilter[dlPrefix] = *instance.Prefix
		}
		if instance.UpdatedAt != nil {
			routeFilter[dlUpdatedAt] = instance.UpdatedAt.String()
		}
		if instance.Ge != nil {
			routeFilter[dlGe] = *instance.Ge
		}
		if instance.Le != nil {
			routeFilter[dlLe] = *instance.Le
		}
		importRouteFilters = append(importRouteFilters, routeFilter)
	}
	d.Set(dlImportRouteFilters, importRouteFilters)
	return nil
}

func resourceIBMdlGatewayActionRead(d *schema.ResourceData, meta interface{}) error {
	directLink, err := directlinkClient(meta)
	if err != nil {
		return err
	}
	ID := d.Id()
	getOptions := &directlinkv1.GetGatewayOptions{
		ID: &ID,
	}
	instanceIntf, response, err := directLink.GetGateway(getOptions)

	if (err != nil) || (instanceIntf == nil) {
		if response != nil && response.StatusCode == 404 {
			d.SetId("")
			return nil
		}
		return fmt.Errorf("[ERROR] Error Getting Direct Link Gateway: %s\n%s", err, response)
	}
	instance := instanceIntf.(*directlinkv1.GetGatewayResponse)
	if instance.Name != nil {
		d.Set(dlName, *instance.Name)
	}
	if instance.Crn != nil {
		d.Set(dlCrn, *instance.Crn)
	}
	if instance.BgpAsn != nil {
		d.Set(dlBgpAsn, *instance.BgpAsn)
	}
	if instance.BgpIbmCidr != nil {
		d.Set(dlBgpIbmCidr, *instance.BgpIbmCidr)
	}
	if instance.BgpIbmAsn != nil {
		d.Set(dlBgpIbmAsn, *instance.BgpIbmAsn)
	}
	if instance.Metered != nil {
		d.Set(dlMetered, *instance.Metered)
	}
	if instance.CrossConnectRouter != nil {
		d.Set(dlCrossConnectRouter, *instance.CrossConnectRouter)
	}
	if instance.BgpBaseCidr != nil {
		d.Set(dlBgpBaseCidr, *instance.BgpBaseCidr)
	}
	if instance.BgpCerCidr != nil {
		d.Set(dlBgpCerCidr, *instance.BgpCerCidr)
	}
	if instance.ProviderApiManaged != nil {
		d.Set(dlProviderAPIManaged, *instance.ProviderApiManaged)
	}
	if instance.Type != nil {
		d.Set(dlType, *instance.Type)
	}
	if instance.SpeedMbps != nil {
		d.Set(dlSpeedMbps, *instance.SpeedMbps)
	}
	if instance.OperationalStatus != nil {
		d.Set(dlOperationalStatus, *instance.OperationalStatus)
	}
	if instance.BgpStatus != nil {
		d.Set(dlBgpStatus, *instance.BgpStatus)
	}
	if instance.BgpStatusUpdatedAt != nil {
		d.Set(dlBgpStatusUpdatedAt, instance.BgpStatusUpdatedAt.String())
	}
	if instance.CompletionNoticeRejectReason != nil {
		d.Set(dlCompletionNoticeRejectReason, *instance.CompletionNoticeRejectReason)
	}
	if instance.LocationName != nil {
		d.Set(dlLocationName, *instance.LocationName)
	}
	if instance.LocationDisplayName != nil {
		d.Set(dlLocationDisplayName, *instance.LocationDisplayName)
	}
	if instance.Vlan != nil {
		d.Set(dlVlan, *instance.Vlan)
	}
	if instance.Global != nil {
		d.Set(dlGlobal, *instance.Global)
	}
	if instance.Port != nil {
		d.Set(dlPort, *instance.Port.ID)
	}
	if instance.LinkStatus != nil {
		d.Set(dlLinkStatus, *instance.LinkStatus)
	}
	if instance.LinkStatusUpdatedAt != nil {
		d.Set(dlLinkStatusUpdatedAt, instance.LinkStatusUpdatedAt.String())
	}
	if instance.CreatedAt != nil {
		d.Set(dlCreatedAt, instance.CreatedAt.String())
	}
	if instance.AuthenticationKey != nil {
		d.Set(dlAuthenticationKey, *instance.AuthenticationKey.Crn)
	}
	if instance.ConnectionMode != nil {
		d.Set(dlConnectionMode, *instance.ConnectionMode)
	}
	if instance.DefaultExportRouteFilter != nil {
		d.Set(dlDefault_export_route_filter, *instance.DefaultExportRouteFilter)
	}
	if instance.DefaultImportRouteFilter != nil {
		d.Set(dlDefault_import_route_filter, *instance.DefaultImportRouteFilter)
	}
	asPrependList := make([]map[string]interface{}, 0)
	if len(instance.AsPrepends) > 0 {
		for _, asPrepend := range instance.AsPrepends {
			asPrependItem := map[string]interface{}{}
			asPrependItem[dlResourceId] = asPrepend.ID
			asPrependItem[dlLength] = asPrepend.Length
			asPrependItem[dlPrefix] = asPrepend.Prefix
			asPrependItem[dlSpecificPrefixes] = asPrepend.SpecificPrefixes
			asPrependItem[dlPolicy] = asPrepend.Policy
			asPrependItem[dlCreatedAt] = asPrepend.CreatedAt.String()
			asPrependItem[dlUpdatedAt] = asPrepend.UpdatedAt.String()

			asPrependList = append(asPrependList, asPrependItem)
		}

	}
	d.Set(dlAsPrepends, asPrependList)

	if instance.ChangeRequest != nil {
		gatewayChangeRequestIntf := instance.ChangeRequest
		gatewayChangeRequest := gatewayChangeRequestIntf.(*directlinkv1.GatewayChangeRequest)
		d.Set(dlChangeRequest, *gatewayChangeRequest.Type)
	}
	tags, err := flex.GetTagsUsingCRN(meta, *instance.Crn)
	if err != nil {
		log.Printf(
			"Error on get of resource direct link gateway (%s) tags: %s", d.Id(), err)
	}
	d.Set(dlTags, tags)
	controller, err := flex.GetBaseController(meta)
	if err != nil {
		return err
	}
	d.Set(flex.ResourceControllerURL, controller+"/interconnectivity/direct-link")
	d.Set(flex.ResourceName, *instance.Name)
	d.Set(flex.ResourceCRN, *instance.Crn)
	d.Set(flex.ResourceStatus, *instance.OperationalStatus)
	if instance.ResourceGroup != nil {
		rg := instance.ResourceGroup
		d.Set(dlResourceGroup, *rg.ID)
		d.Set(flex.ResourceGroupName, *rg.ID)
	}

	//Show the BFD Config parameters if set
	if instance.BfdConfig != nil {
		if instance.BfdConfig.Interval != nil {
			d.Set(dlBfdInterval, *instance.BfdConfig.Interval)
		}

		if instance.BfdConfig.Multiplier != nil {
			d.Set(dlBfdMultiplier, *instance.BfdConfig.Multiplier)
		}

		if instance.BfdConfig.BfdStatus != nil {
			d.Set(dlBfdStatus, *instance.BfdConfig.BfdStatus)
		}

		if instance.BfdConfig.BfdStatusUpdatedAt != nil {
			d.Set(dlBfdStatusUpdatedAt, instance.BfdConfig.BfdStatusUpdatedAt.String())
		}
	}
	resourceIBMdlGatewayActionExportRouteFiltersRead(d, meta)
	resourceIBMdlGatewayActionImportRouteFiltersRead(d, meta)
	return nil
}

func isWaitForDirectLinkAvailableforAction(client *directlinkv1.DirectLinkV1, id string, timeout time.Duration) (interface{}, error) {
	log.Printf("Waiting for direct link (%s) to be provisioned.", id)
	stateConf := &resource.StateChangeConf{
		Pending:    []string{"retry", dlGatewayProvisioning},
		Target:     []string{dlGatewayProvisioningDone, ""},
		Refresh:    isDirectLinkRefreshFuncforAction(client, id),
		Timeout:    timeout,
		Delay:      10 * time.Second,
		MinTimeout: 10 * time.Second,
	}
	return stateConf.WaitForState()
}
func isDirectLinkRefreshFuncforAction(client *directlinkv1.DirectLinkV1, id string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		getOptions := &directlinkv1.GetGatewayOptions{
			ID: &id,
		}

		instanceIntf, response, err := client.GetGateway(getOptions)
		if (err != nil) || (instanceIntf == nil) {
			return nil, "", fmt.Errorf("[ERROR] Error Getting Direct Link: %s\n%s", err, response)
		}
		instance := instanceIntf.(*directlinkv1.GetGatewayResponse)
		if *instance.OperationalStatus == "provisioned" || *instance.OperationalStatus == "failed" || *instance.OperationalStatus == "create_rejected" {
			return instance, dlGatewayProvisioningDone, nil
		}
		return instance, dlGatewayProvisioning, nil
	}
}

func isWaitForDirectLinkActionAvailable(client *directlinkv1.DirectLinkV1, id string, timeout time.Duration) (interface{}, error) {
	log.Printf("Waiting for direct link (%s) to be provisioned.", id)
	stateConf := &resource.StateChangeConf{
		Pending:    []string{"retry", dlGatewayActionUpdate},
		Target:     []string{dlGatewayActionUpdateDone, ""},
		Refresh:    isDirectLinkRefreshActionFunc(client, id),
		Timeout:    timeout,
		Delay:      10 * time.Second,
		MinTimeout: 10 * time.Second,
	}
	return stateConf.WaitForState()
}
func isDirectLinkRefreshActionFunc(client *directlinkv1.DirectLinkV1, id string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		getOptions := &directlinkv1.GetGatewayOptions{
			ID: &id,
		}
		instanceIntf, response, err := client.GetGateway(getOptions)

		if (err != nil) || (instanceIntf == nil) {
			return nil, "", fmt.Errorf("[ERROR] Error Getting Direct Link: %s\n%s", err, response)
		}
		instance := instanceIntf.(*directlinkv1.GetGatewayResponse)
		if instance.ChangeRequest != nil {
			gatewayChangeRequestIntf := instance.ChangeRequest
			gatewayChangeRequest := gatewayChangeRequestIntf.(*directlinkv1.GatewayChangeRequest)
			if *gatewayChangeRequest.Type == "update_attributes" {
				return instance, dlGatewayActionUpdateDone, nil
			}
		}
		return instance, dlGatewayActionUpdate, nil
	}
}

func isWaitForDirectLinkDeleteActionAvailable(client *directlinkv1.DirectLinkV1, id string, timeout time.Duration) (interface{}, error) {
	log.Printf("Waiting for direct link (%s) to be provisioned.", id)
	stateConf := &resource.StateChangeConf{
		Pending:    []string{"retry", dlGatewayDeleteActionUpdate},
		Target:     []string{dlGatewayDeleteActionUpdateDone, ""},
		Refresh:    isDirectLinkRefreshDeleteActionFunc(client, id),
		Timeout:    timeout,
		Delay:      10 * time.Second,
		MinTimeout: 10 * time.Second,
	}
	return stateConf.WaitForState()
}
func isDirectLinkRefreshDeleteActionFunc(client *directlinkv1.DirectLinkV1, id string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		getOptions := &directlinkv1.GetGatewayOptions{
			ID: &id,
		}
		instanceIntf, response, err := client.GetGateway(getOptions)

		if (err != nil) || (instanceIntf == nil) {
			return nil, "", fmt.Errorf("[ERROR] Error Getting Direct Link: %s\n%s", err, response)
		}
		instance := instanceIntf.(*directlinkv1.GetGatewayResponse)
		if instance.ChangeRequest != nil {
			gatewayChangeRequestIntf := instance.ChangeRequest
			gatewayChangeRequest := gatewayChangeRequestIntf.(*directlinkv1.GatewayChangeRequest)
			if *gatewayChangeRequest.Type == "delete_gateway" {
				return instance, dlGatewayDeleteActionUpdateDone, nil
			}
		}
		return instance, dlGatewayDeleteActionUpdate, nil
	}
}

func resourceIBMdlGatewayActionUpdate(d *schema.ResourceData, meta interface{}) error {

	directLink, err := directlinkClient(meta)
	if err != nil {
		return err
	}
	createGatewayActionOptionsModel := &directlinkv1.CreateGatewayActionOptions{}
	action := d.Get(dlCustomerAction).(string)
	createGatewayActionOptionsModel.Action = &action
	gatewayId := d.Get(dlGatewayId).(string)
	createGatewayActionOptionsModel.ID = &gatewayId

	updateList := make([]directlinkv1.GatewayActionTemplateUpdatesItemIntf, 0)
	getOptions := &directlinkv1.GetGatewayOptions{
		ID: &gatewayId,
	}

	if action == "update_attributes_approve" || action == "update_attributes_reject" {
		_, err = isWaitForDirectLinkActionAvailable(directLink, gatewayId, d.Timeout(schema.TimeoutCreate))
		if err != nil {
			return err
		}
		instanceIntf, response, err := directLink.GetGateway(getOptions)

		if (err != nil) || (instanceIntf == nil) {
			if response != nil && response.StatusCode == 404 {
				d.SetId("")
				return nil
			}
			return fmt.Errorf("[ERROR] Error Getting Direct Link Gateway : %s\n%s", err, response)
		}
		instance := instanceIntf.(*directlinkv1.GetGatewayResponse)
		if instance.ChangeRequest != nil {
			gatewayChangeRequestIntf := instance.ChangeRequest
			gatewayChangeRequest := gatewayChangeRequestIntf.(*directlinkv1.GatewayChangeRequest)
			d.Set(dlChangeRequest, *gatewayChangeRequest.Type)

			bgpIPUpdate := new(directlinkv1.GatewayActionTemplateUpdatesItemGatewayClientBGPIPUpdate)

			for _, updatechangeReq := range gatewayChangeRequest.Updates {
				updatechangeReqList := updatechangeReq.(*directlinkv1.GatewayChangeRequestUpdatesItem)
				if updatechangeReqList.SpeedMbps != nil {
					speedUpdate := new(directlinkv1.GatewayActionTemplateUpdatesItemGatewayClientSpeedUpdate)
					speedUpdate.SpeedMbps = updatechangeReqList.SpeedMbps
					updateList = append(updateList, speedUpdate)
				}
				if updatechangeReqList.BgpIbmCidr != nil {

					bgpIPUpdate.BgpIbmCidr = updatechangeReqList.BgpIbmCidr
					updateList = append(updateList, bgpIPUpdate)
				}
				if updatechangeReqList.BgpCerCidr != nil {

					bgpIPUpdate.BgpCerCidr = updatechangeReqList.BgpCerCidr
					updateList = append(updateList, bgpIPUpdate)
				}
				if updatechangeReqList.BgpAsn != nil {
					bgpAsnUpdate := new(directlinkv1.GatewayActionTemplateUpdatesItemGatewayClientBGPASNUpdate)
					bgpAsnUpdate.BgpAsn = updatechangeReqList.BgpAsn
					updateList = append(updateList, bgpAsnUpdate)
				}
				if updatechangeReqList.Vlan != nil {
					vlanUpdate := new(directlinkv1.GatewayActionTemplateUpdatesItemGatewayClientVLANUpdate)
					vlanUpdate.Vlan = updatechangeReqList.Vlan
					updateList = append(updateList, vlanUpdate)
				}

			}
		}
		createGatewayActionOptionsModel.Updates = updateList
		gateway, response, err := directLink.CreateGatewayAction(createGatewayActionOptionsModel)
		if err != nil {
			return fmt.Errorf("[DEBUG] Direct Link Gateway update_attributes_approve err %s\n%s", err, response)
		}
		d.SetId(*gateway.ID)
		_, err = isWaitForDirectLinkAvailableforAction(directLink, d.Id(), d.Timeout(schema.TimeoutCreate))
		if err != nil {
			return err
		}
		return resourceIBMdlGatewayActionRead(d, meta)
	}
	if action == "delete_gateway_approve" || action == "delete_gateway_reject" {
		_, err = isWaitForDirectLinkDeleteActionAvailable(directLink, gatewayId, d.Timeout(schema.TimeoutCreate))
		if err != nil {
			return err
		}
		_, response, err := directLink.CreateGatewayAction(createGatewayActionOptionsModel)
		if err != nil {
			return fmt.Errorf("[DEBUG] delete_gateway_approve failed with error  %s\n%s", err, response)
		}
	}
	return nil
}
func resourceIBMdlGatewayActionDelete(d *schema.ResourceData, meta interface{}) error {

	directLink, err := directlinkClient(meta)
	if err != nil {
		return err
	}
	createGatewayActionOptionsModel := &directlinkv1.CreateGatewayActionOptions{}
	action := d.Get(dlCustomerAction).(string)
	createGatewayActionOptionsModel.Action = &action
	gatewayId := d.Get(dlGatewayId).(string)
	createGatewayActionOptionsModel.ID = &gatewayId

	if action != "delete_gateway_approve" {
		delete_action := "delete_gateway_approve"
		createGatewayActionOptionsModel.Action = &delete_action
	}
	_, response, err := directLink.CreateGatewayAction(createGatewayActionOptionsModel)
	if err != nil {
		fmt.Printf("[DEBUG] delete_gateway_approve failed, may be gateway deleted already %s\n%s", err, response)
	}
	d.SetId("")
	return nil
}
