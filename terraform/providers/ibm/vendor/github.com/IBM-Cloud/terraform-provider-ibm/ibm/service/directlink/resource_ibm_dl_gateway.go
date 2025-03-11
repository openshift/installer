// Copyright IBM Corp. 2017, 2021, 2023 All Rights Reserved.
// Licensed under the Mozilla Public License v2.0

package directlink

import (
	"context"
	"fmt"
	"log"
	"os"
	"reflect"
	"strings"
	"time"

	"github.com/IBM/go-sdk-core/v3/core"
	"github.com/IBM/networking-go-sdk/directlinkv1"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/customdiff"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/conns"
	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/flex"
	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/validate"
)

func ResourceIBMDLGateway() *schema.Resource {
	return &schema.Resource{
		Create:   resourceIBMdlGatewayCreate,
		Read:     resourceIBMdlGatewayRead,
		Delete:   resourceIBMdlGatewayDelete,
		Exists:   resourceIBMdlGatewayExists,
		Update:   resourceIBMdlGatewayUpdate,
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
			dlAuthenticationKey: {
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    false,
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
							ValidateFunc: validate.InvokeValidator("ibm_dl_gateway", dlAction),
							Description:  "Determines whether the  routes that match the prefix-set will be permit or deny",
						},
						dlBefore: {
							Type:        schema.TypeString,
							Optional:    true,
							Computed:    true,
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
							ValidateFunc: validate.InvokeValidator("ibm_dl_gateway", dlAction),
							Description:  "Determines whether the  routes that match the prefix-set will be permit or deny",
						},
						dlBefore: {
							Type:        schema.TypeString,
							Optional:    true,
							Computed:    true,
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
				Computed:     true,
				ValidateFunc: validate.InvokeValidator("ibm_dl_gateway", dlDefault_export_route_filter),
				Description:  "The default directional route filter action that applies to routes that do not match any directional route filters",
			},
			dlDefault_import_route_filter: {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ValidateFunc: validate.InvokeValidator("ibm_dl_gateway", dlDefault_import_route_filter),
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
							ValidateFunc: validate.InvokeValidator("ibm_dl_gateway", dlLength),
							Description:  "Number of times the ASN to appended to the AS Path",
						},
						dlPolicy: {
							Type:         schema.TypeString,
							Required:     true,
							ForceNew:     false,
							ValidateFunc: validate.InvokeValidator("ibm_dl_gateway", dlPolicy),
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
				ValidateFunc: validate.InvokeValidator("ibm_dl_gateway", dlBfdInterval),
			},
			dlBfdMultiplier: {
				Type:         schema.TypeInt,
				Optional:     true,
				ForceNew:     false,
				Description:  "BFD Multiplier",
				ValidateFunc: validate.InvokeValidator("ibm_dl_gateway", dlBfdMultiplier),
			},
			dlBfdStatus: {
				Type:        schema.TypeString,
				Computed:    true,
				Optional:    true,
				Description: "Gateway BFD status",
			},
			dlBfdStatusUpdatedAt: {
				Type:        schema.TypeString,
				Computed:    true,
				Optional:    true,
				Description: "Date and time BFD status was updated",
			},
			dlBgpAsn: {
				Type:        schema.TypeInt,
				Required:    true,
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
				ValidateFunc: validate.InvokeValidator("ibm_dl_gateway", dlConnectionMode),
			},
			dlCrossConnectRouter: {
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				Description: "Cross connect router",
			},
			dlGlobal: {
				Type:        schema.TypeBool,
				Required:    true,
				ForceNew:    false,
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
				Required:    true,
				ForceNew:    false,
				Description: "Metered billing option",
			},
			dlName: {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     false,
				Description:  "The unique user-defined name for this gateway",
				ValidateFunc: validate.InvokeValidator("ibm_dl_gateway", dlName),
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
				Required:    true,
				ForceNew:    false,
				Description: "Gateway speed in megabits per second",
			},
			dlType: {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				Description:  "Gateway type",
				ValidateFunc: validate.InvokeValidator("ibm_dl_gateway", dlType),
			},
			dlMacSecConfig: {
				Type:        schema.TypeList,
				MinItems:    0,
				MaxItems:    1,
				Optional:    true,
				ForceNew:    false,
				Description: "MACsec configuration information",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						dlActive: {
							Type:        schema.TypeBool,
							Required:    true,
							ForceNew:    false,
							Description: "Indicate whether MACsec protection should be active (true) or inactive (false) for this MACsec enabled gateway",
						},
						dlPrimaryCak: {
							Type:        schema.TypeString,
							Required:    true,
							ForceNew:    false,
							Description: "Desired primary connectivity association key. Keys for a MACsec configuration must have names with an even number of characters from [0-9a-fA-F]",
						},
						dlFallbackCak: {
							Type:        schema.TypeString,
							Optional:    true,
							ForceNew:    false,
							Description: "Fallback connectivity association key. Keys used for MACsec configuration must have names with an even number of characters from [0-9a-fA-F]",
						},
						dlWindowSize: {
							Type:        schema.TypeInt,
							Optional:    true,
							ForceNew:    false,
							Default:     148809600,
							Description: "Replay protection window size",
						},
						dlActiveCak: {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Active connectivity association key.",
						},
						dlSakExpiryTime: {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Secure Association Key (SAK) expiry time in seconds",
						},
						dlCipherSuite: {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "SAK cipher suite",
						},
						dlConfidentialityOffset: {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Confidentiality Offset",
						},
						dlCryptographicAlgorithm: {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Cryptographic Algorithm",
						},
						dlKeyServerPriority: {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Key Server Priority",
						},
						dlMacSecConfigStatus: {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The current status of MACsec on the device for this gateway",
						},
						dlSecurityPolicy: {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Packets without MACsec headers are not dropped when security_policy is should_secure.",
						},
					},
				},
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
				Type:          schema.TypeInt,
				Optional:      true,
				Computed:      true,
				Description:   "VLAN allocated for this gateway",
				ConflictsWith: []string{"remove_vlan"},
				ValidateFunc:  validate.InvokeValidator("ibm_dl_gateway", dlVlan),
			},
			dlRemoveVlan: {
				Type:          schema.TypeBool,
				Optional:      true,
				Default:       false,
				ConflictsWith: []string{"vlan"},
				Description:   "Remove VLAN allocated for this dedicated gateway",
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
				Elem:        &schema.Schema{Type: schema.TypeString, ValidateFunc: validate.InvokeValidator("ibm_dl_gateway", dlTags)},
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

func ResourceIBMDLGatewayValidator() *validate.ResourceValidator {

	validateSchema := make([]validate.ValidateSchema, 0)
	dlTypeAllowedValues := "dedicated, connect"
	dlConnectionModeAllowedValues := "direct, transit"
	dlPolicyAllowedValues := "export, import"
	dlActionValues := "permit, deny"

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

	validateSchema = append(validateSchema,
		validate.ValidateSchema{
			Identifier:                 dlVlan,
			ValidateFunctionIdentifier: validate.IntBetween,
			Type:                       validate.TypeInt,
			MinValue:                   "2",
			MaxValue:                   "3967"})

	ibmISDLGatewayResourceValidator := validate.ResourceValidator{ResourceName: "ibm_dl_gateway", Schema: validateSchema}
	return &ibmISDLGatewayResourceValidator
}

func directlinkClient(meta interface{}) (*directlinkv1.DirectLinkV1, error) {
	sess, err := meta.(conns.ClientSession).DirectlinkV1API()
	return sess, err
}

func resourceIBMdlGatewayCreate(d *schema.ResourceData, meta interface{}) error {
	directLink, err := directlinkClient(meta)
	if err != nil {
		return err
	}
	dtype := d.Get(dlType).(string)
	createGatewayOptionsModel := &directlinkv1.CreateGatewayOptions{}
	name := d.Get(dlName).(string)
	speed := int64(d.Get(dlSpeedMbps).(int))
	global := d.Get(dlGlobal).(bool)
	bgpAsn := int64(d.Get(dlBgpAsn).(int))
	metered := d.Get(dlMetered).(bool)

	var bfdConfig directlinkv1.GatewayBfdConfigTemplate
	isBfdInterval := false

	if bfdInterval, ok := d.GetOk(dlBfdInterval); ok {
		isBfdInterval = true
		interval := int64(bfdInterval.(int))
		bfdConfig.Interval = &interval
	}

	if bfdMultiplier, ok := d.GetOk(dlBfdMultiplier); ok {
		multiplier := int64(bfdMultiplier.(int))
		bfdConfig.Multiplier = &multiplier
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
	if dtype == "dedicated" {
		var crossConnectRouter, carrierName, locationName, customerName string
		if _, ok := d.GetOk(dlCarrierName); ok {
			carrierName = d.Get(dlCarrierName).(string)
		} else {
			err = fmt.Errorf("[ERROR] Error creating gateway, %s is a required field", dlCarrierName)
			log.Printf("%s is a required field", dlCarrierName)
			return err
		}
		if _, ok := d.GetOk(dlCrossConnectRouter); ok {
			crossConnectRouter = d.Get(dlCrossConnectRouter).(string)
		} else {
			err = fmt.Errorf("[ERROR] Error creating gateway, %s is a required field", dlCrossConnectRouter)
			log.Printf("%s is a required field", dlCrossConnectRouter)
			return err
		}
		if _, ok := d.GetOk(dlLocationName); ok {
			locationName = d.Get(dlLocationName).(string)
		} else {
			err = fmt.Errorf("[ERROR] Error creating gateway, %s is a required field", dlLocationName)
			log.Printf("%s is a required field", dlLocationName)
			return err
		}
		if _, ok := d.GetOk(dlCustomerName); ok {
			customerName = d.Get(dlCustomerName).(string)
		} else {
			err = fmt.Errorf("[ERROR] Error creating gateway, %s is a required field", dlCustomerName)
			log.Printf("%s is a required field", dlCustomerName)
			return err
		}
		gatewayDedicatedTemplateModel, _ := directLink.NewGatewayTemplateGatewayTypeDedicatedTemplate(bgpAsn, global, metered, name, speed, dtype, carrierName, crossConnectRouter, customerName, locationName)

		if _, ok := d.GetOk(dlBgpIbmCidr); ok {
			bgpIbmCidr := d.Get(dlBgpIbmCidr).(string)
			gatewayDedicatedTemplateModel.BgpIbmCidr = &bgpIbmCidr

		}
		if _, ok := d.GetOk(dlBgpCerCidr); ok {
			bgpCerCidr := d.Get(dlBgpCerCidr).(string)
			gatewayDedicatedTemplateModel.BgpCerCidr = &bgpCerCidr

		}
		if _, ok := d.GetOk(dlResourceGroup); ok {
			resourceGroup := d.Get(dlResourceGroup).(string)
			gatewayDedicatedTemplateModel.ResourceGroup = &directlinkv1.ResourceGroupIdentity{ID: &resourceGroup}

		}
		if _, ok := d.GetOk(dlBgpBaseCidr); ok {
			bgpBaseCidr := d.Get(dlBgpBaseCidr).(string)
			gatewayDedicatedTemplateModel.BgpBaseCidr = &bgpBaseCidr
		}
		if _, ok := d.GetOk(dlMacSecConfig); ok {
			// Construct an instance of the GatewayMacsecConfigTemplate model
			gatewayMacsecConfigTemplateModel := new(directlinkv1.GatewayMacsecConfigTemplate)
			activebool := d.Get("macsec_config.0.active").(bool)
			gatewayMacsecConfigTemplateModel.Active = &activebool

			// Construct an instance of the GatewayMacsecCak model
			gatewayMacsecCakModel := new(directlinkv1.GatewayMacsecConfigTemplatePrimaryCak)
			primaryCakstr := d.Get("macsec_config.0.primary_cak").(string)
			gatewayMacsecCakModel.Crn = &primaryCakstr
			gatewayMacsecConfigTemplateModel.PrimaryCak = gatewayMacsecCakModel

			if fallbackCak, ok := d.GetOk("macsec_config.0.fallback_cak"); ok {
				// Construct an instance of the GatewayMacsecCak model
				gatewayMacsecCakModel := new(directlinkv1.GatewayMacsecConfigTemplateFallbackCak)
				fallbackCakstr := fallbackCak.(string)
				gatewayMacsecCakModel.Crn = &fallbackCakstr
				gatewayMacsecConfigTemplateModel.FallbackCak = gatewayMacsecCakModel
			}
			if windowSize, ok := d.GetOk("macsec_config.0.window_size"); ok {
				windowSizeint := int64(windowSize.(int))
				gatewayMacsecConfigTemplateModel.WindowSize = &windowSizeint
			}
			gatewayDedicatedTemplateModel.MacsecConfig = gatewayMacsecConfigTemplateModel
		}

		if authKeyCrn, ok := d.GetOk(dlAuthenticationKey); ok {
			authKeyCrnStr := authKeyCrn.(string)
			gatewayDedicatedTemplateModel.AuthenticationKey = &directlinkv1.GatewayTemplateAuthenticationKey{Crn: &authKeyCrnStr}
		}

		if connectionMode, ok := d.GetOk(dlConnectionMode); ok {
			connectionModeStr := connectionMode.(string)
			gatewayDedicatedTemplateModel.ConnectionMode = &connectionModeStr
		}

		if !reflect.DeepEqual(bfdConfig, directlinkv1.GatewayBfdConfigTemplate{}) {
			gatewayDedicatedTemplateModel.BfdConfig = &bfdConfig
		}

		if len(asPrependsCreateItems) > 0 {
			gatewayDedicatedTemplateModel.AsPrepends = asPrependsCreateItems
		}
		if len(exportRouteFiltersCreateList) > 0 {
			gatewayDedicatedTemplateModel.ExportRouteFilters = exportRouteFiltersCreateList
		}
		if len(importRouteFiltersCreateList) > 0 {
			gatewayDedicatedTemplateModel.ImportRouteFilters = importRouteFiltersCreateList
		}
		if default_export_route_filter, ok := d.GetOk(dlDefault_export_route_filter); ok {
			gatewayDedicatedTemplateModel.DefaultExportRouteFilter = NewStrPointer(default_export_route_filter.(string))
		}
		if default_import_route_filter, ok := d.GetOk(dlDefault_import_route_filter); ok {
			gatewayDedicatedTemplateModel.DefaultImportRouteFilter = NewStrPointer(default_import_route_filter.(string))
		}

		if vlan, ok := d.GetOk(dlVlan); ok {
			mapped_vlan := int64(vlan.(int))
			gatewayDedicatedTemplateModel.Vlan = &mapped_vlan
		}
		createGatewayOptionsModel.GatewayTemplate = gatewayDedicatedTemplateModel

	} else if dtype == "connect" {
		var portID string
		if _, ok := d.GetOk(dlPort); ok {
			portID = d.Get(dlPort).(string)
		}
		if portID != "" {
			portIdentity, _ := directLink.NewGatewayPortIdentity(portID)
			gatewayConnectTemplateModel, _ := directLink.NewGatewayTemplateGatewayTypeConnectTemplate(bgpAsn, global, metered, name, speed, dtype, portIdentity)

			if _, ok := d.GetOk(dlBgpIbmCidr); ok {
				bgpIbmCidr := d.Get(dlBgpIbmCidr).(string)
				gatewayConnectTemplateModel.BgpIbmCidr = &bgpIbmCidr

			}
			if _, ok := d.GetOk(dlBgpBaseCidr); ok {
				bgpBaseCidr := d.Get(dlBgpBaseCidr).(string)
				gatewayConnectTemplateModel.BgpBaseCidr = &bgpBaseCidr
			}
			if _, ok := d.GetOk(dlBgpCerCidr); ok {
				bgpCerCidr := d.Get(dlBgpCerCidr).(string)
				gatewayConnectTemplateModel.BgpCerCidr = &bgpCerCidr

			}
			if _, ok := d.GetOk(dlResourceGroup); ok {
				resourceGroup := d.Get(dlResourceGroup).(string)
				gatewayConnectTemplateModel.ResourceGroup = &directlinkv1.ResourceGroupIdentity{ID: &resourceGroup}

			}

			if authKeyCrn, ok := d.GetOk(dlAuthenticationKey); ok {
				authKeyCrnStr := authKeyCrn.(string)
				gatewayConnectTemplateModel.AuthenticationKey = &directlinkv1.GatewayTemplateAuthenticationKey{Crn: &authKeyCrnStr}
			}

			if connectionMode, ok := d.GetOk(dlConnectionMode); ok {
				connectionModeStr := connectionMode.(string)
				gatewayConnectTemplateModel.ConnectionMode = &connectionModeStr
			}

			if !reflect.DeepEqual(bfdConfig, directlinkv1.GatewayBfdConfigTemplate{}) {
				gatewayConnectTemplateModel.BfdConfig = &bfdConfig
			}

			if len(asPrependsCreateItems) > 0 {
				gatewayConnectTemplateModel.AsPrepends = asPrependsCreateItems
			}
			if len(exportRouteFiltersCreateList) > 0 {
				gatewayConnectTemplateModel.ExportRouteFilters = exportRouteFiltersCreateList
			}
			if len(importRouteFiltersCreateList) > 0 {
				gatewayConnectTemplateModel.ImportRouteFilters = importRouteFiltersCreateList
			}
			if default_export_route_filter, ok := d.GetOk(dlDefault_export_route_filter); ok {
				gatewayConnectTemplateModel.DefaultExportRouteFilter = NewStrPointer(default_export_route_filter.(string))
			}
			if default_import_route_filter, ok := d.GetOk(dlDefault_import_route_filter); ok {
				gatewayConnectTemplateModel.DefaultImportRouteFilter = NewStrPointer(default_import_route_filter.(string))
			}
			createGatewayOptionsModel.GatewayTemplate = gatewayConnectTemplateModel

		} else {
			err = fmt.Errorf("[ERROR] Error creating direct link connect gateway, %s is a required field", dlPort)
			return err
		}
	}

	gateway, response, err := directLink.CreateGateway(createGatewayOptionsModel)
	if err != nil {
		return fmt.Errorf("[DEBUG] Create Direct Link Gateway (%s) err %s\n%s", dtype, err, response)
	}
	d.SetId(*gateway.ID)

	log.Printf("[INFO] Created Direct Link Gateway (%s Template) : %s", dtype, *gateway.ID)
	if dtype == "connect" {
		getPortOptions := directLink.NewGetPortOptions(*gateway.Port.ID)
		port, response, err := directLink.GetPort(getPortOptions)
		if err != nil {
			return fmt.Errorf("[ERROR] Error getting port %s %s", response, err)
		}
		if port != nil && port.ProviderName != nil && !strings.Contains(strings.ToLower(*port.ProviderName), "netbond") && !strings.Contains(strings.ToLower(*port.ProviderName), "megaport") {
			_, err = isWaitForDirectLinkAvailable(directLink, d.Id(), d.Timeout(schema.TimeoutCreate))
			if err != nil {
				return err
			}
		}

	}

	v := os.Getenv("IC_ENV_TAGS")
	if _, ok := d.GetOk(dlTags); ok || v != "" {
		oldList, newList := d.GetChange(dlTags)
		err = flex.UpdateTagsUsingCRN(oldList, newList, meta, *gateway.Crn)
		if err != nil {
			log.Printf(
				"Error on create of resource direct link gateway %s (%s) tags: %s", dtype, d.Id(), err)
		}
	}

	return resourceIBMdlGatewayRead(d, meta)
}

func resourceIBMdlGatewayExportRouteFiltersRead(d *schema.ResourceData, meta interface{}) error {
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

func resourceIBMdlGatewayImportRouteFiltersRead(d *schema.ResourceData, meta interface{}) error {
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

func resourceIBMdlGatewayRead(d *schema.ResourceData, meta interface{}) error {
	dtype := d.Get(dlType).(string)
	log.Printf("[INFO] Inside resourceIBMdlGatewayRead: %s", dtype)

	directLink, err := directlinkClient(meta)
	if err != nil {
		return err
	}

	ID := d.Id()

	getOptions := &directlinkv1.GetGatewayOptions{
		ID: &ID,
	}
	log.Printf("[INFO] Calling getgateway api: %s", dtype)

	instanceIntf, response, err := directLink.GetGateway(getOptions)
	if (err != nil) || (instanceIntf == nil) {
		if response != nil && response.StatusCode == 404 {
			d.SetId("")
			return nil
		}
		return fmt.Errorf("[ERROR] Error Getting Direct Link Gateway (%s Template): %s\n%s", dtype, err, response)
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
	} else {
		d.Set(dlVlan, nil)
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

	if dtype == "dedicated" {
		if instance.MacsecConfig != nil {
			macsecList := make([]map[string]interface{}, 0)
			currentMacSec := map[string]interface{}{}
			// Construct an instance of the GatewayMacsecConfigTemplate model
			gatewayMacsecConfigTemplateModel := instance.MacsecConfig
			if gatewayMacsecConfigTemplateModel.Active != nil {
				currentMacSec[dlActive] = *gatewayMacsecConfigTemplateModel.Active
			}
			if gatewayMacsecConfigTemplateModel.ActiveCak != nil {
				if gatewayMacsecConfigTemplateModel.ActiveCak.Crn != nil {
					currentMacSec[dlActiveCak] = *gatewayMacsecConfigTemplateModel.ActiveCak.Crn
				}
			}
			if gatewayMacsecConfigTemplateModel.PrimaryCak != nil {
				currentMacSec[dlPrimaryCak] = *gatewayMacsecConfigTemplateModel.PrimaryCak.Crn
			}
			if gatewayMacsecConfigTemplateModel.FallbackCak != nil {
				if gatewayMacsecConfigTemplateModel.FallbackCak.Crn != nil {
					currentMacSec[dlFallbackCak] = *gatewayMacsecConfigTemplateModel.FallbackCak.Crn
				}
			}
			if gatewayMacsecConfigTemplateModel.SakExpiryTime != nil {
				currentMacSec[dlSakExpiryTime] = *gatewayMacsecConfigTemplateModel.SakExpiryTime
			}
			if gatewayMacsecConfigTemplateModel.SecurityPolicy != nil {
				currentMacSec[dlSecurityPolicy] = *gatewayMacsecConfigTemplateModel.SecurityPolicy
			}
			if gatewayMacsecConfigTemplateModel.WindowSize != nil {
				currentMacSec[dlWindowSize] = *gatewayMacsecConfigTemplateModel.WindowSize
			}
			if gatewayMacsecConfigTemplateModel.CipherSuite != nil {
				currentMacSec[dlCipherSuite] = *gatewayMacsecConfigTemplateModel.CipherSuite
			}
			if gatewayMacsecConfigTemplateModel.ConfidentialityOffset != nil {
				currentMacSec[dlConfidentialityOffset] = *gatewayMacsecConfigTemplateModel.ConfidentialityOffset
			}
			if gatewayMacsecConfigTemplateModel.CryptographicAlgorithm != nil {
				currentMacSec[dlCryptographicAlgorithm] = *gatewayMacsecConfigTemplateModel.CryptographicAlgorithm
			}
			if gatewayMacsecConfigTemplateModel.KeyServerPriority != nil {
				currentMacSec[dlKeyServerPriority] = *gatewayMacsecConfigTemplateModel.KeyServerPriority
			}
			if gatewayMacsecConfigTemplateModel.Status != nil {
				currentMacSec[dlMacSecConfigStatus] = *gatewayMacsecConfigTemplateModel.Status
			}
			macsecList = append(macsecList, currentMacSec)
			d.Set(dlMacSecConfig, macsecList)
		}
	}
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
	resourceIBMdlGatewayExportRouteFiltersRead(d, meta)
	resourceIBMdlGatewayImportRouteFiltersRead(d, meta)
	return nil
}
func isWaitForDirectLinkAvailable(client *directlinkv1.DirectLinkV1, id string, timeout time.Duration) (interface{}, error) {
	log.Printf("Waiting for direct link (%s) to be provisioned.", id)
	stateConf := &resource.StateChangeConf{
		Pending:    []string{"retry", dlGatewayProvisioning},
		Target:     []string{dlGatewayProvisioningDone, ""},
		Refresh:    isDirectLinkRefreshFunc(client, id),
		Timeout:    timeout,
		Delay:      10 * time.Second,
		MinTimeout: 10 * time.Second,
	}
	return stateConf.WaitForState()
}
func isDirectLinkRefreshFunc(client *directlinkv1.DirectLinkV1, id string) resource.StateRefreshFunc {
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

func resourceIBMdlGatewayUpdate(d *schema.ResourceData, meta interface{}) error {

	directLink, err := directlinkClient(meta)
	if err != nil {
		return err
	}

	ID := d.Id()
	getOptions := &directlinkv1.GetGatewayOptions{
		ID: &ID,
	}
	instanceIntf, detail, err := directLink.GetGateway(getOptions)

	if (err != nil) || (instanceIntf == nil) {
		log.Printf("Error fetching Direct Link Gateway :%s", detail)
		return err
	}
	instance := instanceIntf.(*directlinkv1.GetGatewayResponse)
	gatewayPatchTemplateModel := map[string]interface{}{}
	dtype := *instance.Type

	if d.HasChange(dlTags) {
		oldList, newList := d.GetChange(dlTags)
		err = flex.UpdateTagsUsingCRN(oldList, newList, meta, *instance.Crn)
		if err != nil {
			log.Printf(
				"Error on update of resource direct link gateway (%s) tags: %s", *instance.ID, err)
		}
	}

	if d.HasChange(dlName) {
		name := d.Get(dlName).(string)
		gatewayPatchTemplateModel["name"] = &name
	}
	if d.HasChange(dlSpeedMbps) {
		speed := int64(d.Get(dlSpeedMbps).(int))
		gatewayPatchTemplateModel["speed_mbps"] = &speed
	}
	if d.HasChange(dlBgpAsn) {
		bgpAsn := int64(d.Get(dlBgpAsn).(int))
		gatewayPatchTemplateModel["bgp_asn"] = &bgpAsn
	}
	if d.HasChange(dlBgpCerCidr) {
		bgpCerCidr := d.Get(dlBgpCerCidr).(string)
		gatewayPatchTemplateModel["bgp_cer_cidr"] = &bgpCerCidr
	}
	if d.HasChange(dlBgpIbmCidr) {
		bgpIbmCidr := d.Get(dlBgpIbmCidr).(string)
		gatewayPatchTemplateModel["bgp_ibm_cidr"] = &bgpIbmCidr
	}
	if d.HasChange(dlAsPrepends) {
		listGatewayAsPrependsOptions := directLink.NewListGatewayAsPrependsOptions(ID)
		_, response, operationErr := directLink.ListGatewayAsPrepends(listGatewayAsPrependsOptions)
		if operationErr != nil {
			log.Printf("[DEBUG] Error listing Direct Link Gateway AS Prepends err %s\n%s", err, response)
			return fmt.Errorf("[ERROR] Error listing Direct Link Gateway AS Prepends err %s\n%s", err, response)
		}
		etag := response.GetHeaders().Get("etag")
		asPrependsCreateItems := make([]directlinkv1.AsPrependPrefixArrayTemplate, 0)
		if asPrependsInput, ok := d.GetOk(dlAsPrepends); ok {
			asPrependsItems := asPrependsInput.([]interface{})

			for _, asPrependItem := range asPrependsItems {
				i := asPrependItem.(map[string]interface{})

				// Construct an instance of the AsPrependTemplate model
				asPrependTemplateModel := new(directlinkv1.AsPrependPrefixArrayTemplate)
				asPrependTemplateModel.Length = NewInt64Pointer(int64(i[dlLength].(int)))
				asPrependTemplateModel.Policy = NewStrPointer(i[dlPolicy].(string))
				asPrependTemplateModel.SpecificPrefixes = nil
				_, prefix_ok := i[dlPrefix]

				sp_prefixOk, ok := i[dlSpecificPrefixes]
				if ok && len(sp_prefixOk.([]interface{})) > 0 {
					asPrependTemplateModel.SpecificPrefixes = flex.ExpandStringList(sp_prefixOk.([]interface{}))
				} else if prefix_ok && (len(i[dlPrefix].(string)) > 0) {
					asPrependTemplateModel.SpecificPrefixes = strings.Split(i[dlPrefix].(string), ",")
				}

				asPrependsCreateItems = append(asPrependsCreateItems, *asPrependTemplateModel)
			}
		}

		// Construct an instance of the AsPrependPrefixArrayTemplate model

		replaceGatewayAsPrependsOptionsModel := new(directlinkv1.ReplaceGatewayAsPrependsOptions)
		replaceGatewayAsPrependsOptionsModel.GatewayID = &ID
		replaceGatewayAsPrependsOptionsModel.IfMatch = core.StringPtr(etag)
		replaceGatewayAsPrependsOptionsModel.AsPrepends = asPrependsCreateItems
		replaceGatewayAsPrependsOptionsModel.Headers = map[string]string{"If-Match": etag}

		_, responseRep, operationErr := directLink.ReplaceGatewayAsPrepends(replaceGatewayAsPrependsOptionsModel)
		if operationErr != nil {
			log.Printf("[DEBUG] Error while replacing AS Prepends to a gateway id %s %s\n%s", ID, operationErr, responseRep)
			return fmt.Errorf("[ERROR] Error while replacing AS Prepends to a gateway id %s %s\n%s", ID, operationErr, responseRep)
		}
	}
	if d.HasChange(dlExportRouteFilters) {

		listGatewayExportRouteFiltersOptionsModel := &directlinkv1.ListGatewayExportRouteFiltersOptions{GatewayID: &ID}
		_, response, operationErr := directLink.ListGatewayExportRouteFilters(listGatewayExportRouteFiltersOptionsModel)
		if operationErr != nil {
			log.Printf("[DEBUG] Error listing the Direct Link Export Route Filters  %s\n%s", operationErr, response)
			return fmt.Errorf("[ERROR] Error listing Direct Link Gateway Export Route Filters %s\n%s", operationErr, response)
		}
		etag := response.GetHeaders().Get("etag")
		exportRouteFiltersReplaceList := make([]directlinkv1.GatewayTemplateRouteFilter, 0)
		if exportRouteFiltersInputList, ok := d.GetOk(dlExportRouteFilters); ok {
			exportRouteFilters := exportRouteFiltersInputList.([]interface{})

			for _, exportRouteFilter := range exportRouteFilters {
				filtersData := exportRouteFilter.(map[string]interface{})

				// Construct an instance of the Export Route Fileter Template model
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
				exportRouteFiltersReplaceList = append(exportRouteFiltersReplaceList, *exportRouteFilterTemplateModel)
			}
		}
		replaceGatewayExportRouteFiltersOptionsModel := new(directlinkv1.ReplaceGatewayExportRouteFiltersOptions)
		replaceGatewayExportRouteFiltersOptionsModel.GatewayID = core.StringPtr(ID)
		replaceGatewayExportRouteFiltersOptionsModel.ExportRouteFilters = exportRouteFiltersReplaceList
		replaceGatewayExportRouteFiltersOptionsModel.IfMatch = core.StringPtr(etag)
		replaceGatewayExportRouteFiltersOptionsModel.Headers = map[string]string{"If-Match": etag}

		/* after updating the Asprepends , waiting  for
		   gateway to move to provisioned state
		*/
		_, err = isWaitForDirectLinkAvailable(directLink, ID, d.Timeout(schema.TimeoutCreate))
		if err != nil {
			return err
		}
		_, response, err := directLink.ReplaceGatewayExportRouteFilters(replaceGatewayExportRouteFiltersOptionsModel)
		if err != nil {
			log.Printf("[DEBUG] Error while replacing Export Route Fileter to a gateway id %s %s\n%s", ID, err, response)
			return fmt.Errorf("[ERROR] Error while replacing Export Route Fileter to a gateway id %s %s\n%s", ID, err, response)
		}
	}

	if d.HasChange(dlImportRouteFilters) {
		listGatewayImportRouteFiltersOptionsModel := &directlinkv1.ListGatewayImportRouteFiltersOptions{GatewayID: &ID}
		_, response, operationErr := directLink.ListGatewayImportRouteFilters(listGatewayImportRouteFiltersOptionsModel)
		if operationErr != nil {
			log.Printf("[DEBUG] Error listing the Direct Link Import Route Filters  %s\n%s", operationErr, response)
			return fmt.Errorf("[ERROR] Error listing Direct Link Gateway Import Route Filters %s\n%s", operationErr, response)
		}
		etag := response.GetHeaders().Get("etag")
		importRouteFiltersReplaceList := make([]directlinkv1.GatewayTemplateRouteFilter, 0)
		if importRouteFiltersInputList, ok := d.GetOk(dlImportRouteFilters); ok {
			importRouteFilters := importRouteFiltersInputList.([]interface{})

			for _, importRouteFilter := range importRouteFilters {
				filtersData := importRouteFilter.(map[string]interface{})

				// Construct an instance of the Export Route Fileter Template model
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
				importRouteFiltersReplaceList = append(importRouteFiltersReplaceList, *importRouteFilterTemplateModel)
			}
		}
		replaceGatewayImportRouteFiltersOptionsModel := new(directlinkv1.ReplaceGatewayImportRouteFiltersOptions)
		replaceGatewayImportRouteFiltersOptionsModel.GatewayID = core.StringPtr(ID)
		replaceGatewayImportRouteFiltersOptionsModel.ImportRouteFilters = importRouteFiltersReplaceList
		replaceGatewayImportRouteFiltersOptionsModel.IfMatch = core.StringPtr(etag)
		replaceGatewayImportRouteFiltersOptionsModel.Headers = map[string]string{"If-Match": etag}

		/* after updating the export route filter , waiting  for
		   gatewat to move to provisioned state
		*/
		_, err = isWaitForDirectLinkAvailable(directLink, ID, d.Timeout(schema.TimeoutCreate))
		if err != nil {
			return err
		}
		_, response, err := directLink.ReplaceGatewayImportRouteFilters(replaceGatewayImportRouteFiltersOptionsModel)
		if err != nil {
			log.Printf("[DEBUG] Error while replacing Import Route Fileter to a gateway id %s %s\n%s", ID, err, response)
			return fmt.Errorf("[ERROR] Error while replacing Import Route Fileter to a gateway id %s %s\n%s", ID, err, response)
		}

	}
	/*
		NOTE: Operational Status cannot be maintained in terraform. The status keeps changing automatically in server side.
		Hence, cannot be maintained in terraform.
		Operational Status and LoaRejectReason are linked.
		Hence, a user cannot update through terraform.

		if d.HasChange(dlOperationalStatus) {
			if _, ok := d.GetOk(dlOperationalStatus); ok {
				operStatus := d.Get(dlOperationalStatus).(string)
				updateGatewayOptionsModel.OperationalStatus = &operStatus
			}
			if _, ok := d.GetOk(dlLoaRejectReason); ok {
				loaRejectReason := d.Get(dlLoaRejectReason).(string)
				updateGatewayOptionsModel.LoaRejectReason = &loaRejectReason
			}
		}
	*/
	if d.HasChange(dlDefault_export_route_filter) {
		gatewayPatchTemplateModel["default_export_route_filter"] = NewStrPointer(d.Get(dlDefault_export_route_filter).(string))
	}
	if d.HasChange(dlDefault_import_route_filter) {
		gatewayPatchTemplateModel["default_import_route_filter"] = NewStrPointer(d.Get(dlDefault_import_route_filter).(string))
	}
	if d.HasChange(dlGlobal) {
		global := d.Get(dlGlobal).(bool)
		gatewayPatchTemplateModel["global"] = &global
	}
	if d.HasChange(dlMetered) {
		metered := d.Get(dlMetered).(bool)
		gatewayPatchTemplateModel["metered"] = &metered
	}
	if d.HasChange(dlAuthenticationKey) {
		authenticationKeyCrn := d.Get(dlAuthenticationKey).(string)
		authenticationKeyPatchTemplate := new(directlinkv1.GatewayPatchTemplateAuthenticationKey)
		authenticationKeyPatchTemplate.Crn = &authenticationKeyCrn
		gatewayPatchTemplateModel["authentication_key"] = authenticationKeyPatchTemplate
	}

	if mode, ok := d.GetOk(dlConnectionMode); ok && d.HasChange(dlConnectionMode) {
		updatedConnectionMode := mode.(string)
		gatewayPatchTemplateModel["connection_mode"] = &updatedConnectionMode
	}

	var updatedBfdConfig directlinkv1.GatewayBfdPatchTemplate
	if bfdInterval, ok := d.GetOk(dlBfdInterval); ok && d.HasChange(dlBfdInterval) {
		updatedBfdInterval := bfdInterval.(int64)
		updatedBfdConfig.Interval = &updatedBfdInterval
	}

	if bfdMultiplier, ok := d.GetOk(dlBfdMultiplier); ok && d.HasChange(dlBfdMultiplier) {
		updatedbfdMultiplier := bfdMultiplier.(int64)
		updatedBfdConfig.Multiplier = &updatedbfdMultiplier
	}

	if !reflect.DeepEqual(updatedBfdConfig, directlinkv1.GatewayBfdPatchTemplate{}) {
		gatewayPatchTemplateModel["bfd_config"] = &updatedBfdConfig
	}

	if dtype == "dedicated" {
		if d.HasChange(dlMacSecConfig) && !d.IsNewResource() {
			// Construct an instance of the GatewayMacsecConfigTemplate model
			gatewayMacsecConfigTemplatePatchModel := new(directlinkv1.GatewayMacsecConfigPatchTemplate)
			if d.HasChange("macsec_config.0.active") {
				activebool := d.Get("macsec_config.0.active").(bool)
				gatewayMacsecConfigTemplatePatchModel.Active = &activebool
			}
			if d.HasChange("macsec_config.0.primary_cak") {
				// Construct an instance of the GatewayMacsecCak model
				gatewayMacsecCakModel := new(directlinkv1.GatewayMacsecConfigPatchTemplatePrimaryCak)
				primaryCakstr := d.Get("macsec_config.0.primary_cak").(string)
				gatewayMacsecCakModel.Crn = &primaryCakstr
				gatewayMacsecConfigTemplatePatchModel.PrimaryCak = gatewayMacsecCakModel
			}
			if d.HasChange("macsec_config.0.fallback_cak") {
				// Construct an instance of the GatewayMacsecCak model
				gatewayMacsecCakModel := new(directlinkv1.GatewayMacsecConfigPatchTemplateFallbackCak)
				if _, ok := d.GetOk("macsec_config.0.fallback_cak"); ok {
					fallbackCakstr := d.Get("macsec_config.0.fallback_cak").(string)
					gatewayMacsecCakModel.Crn = &fallbackCakstr
					gatewayMacsecConfigTemplatePatchModel.FallbackCak = gatewayMacsecCakModel
				} else {
					fallbackCakstr := ""
					gatewayMacsecCakModel.Crn = &fallbackCakstr
				}
				gatewayMacsecConfigTemplatePatchModel.FallbackCak = gatewayMacsecCakModel
			}
			if d.HasChange("macsec_config.0.window_size") {
				if _, ok := d.GetOk("macsec_config.0.window_size"); ok {
					windowSizeint := int64(d.Get("macsec_config.0.window_size").(int))
					gatewayMacsecConfigTemplatePatchModel.WindowSize = &windowSizeint
				}
			}
			gatewayPatchTemplateModel["macsec_config"] = gatewayMacsecConfigTemplatePatchModel
		} else {
			gatewayPatchTemplateModel["macsec_config"] = nil
		}
		if d.HasChange(dlVlan) {
			if _, ok := d.GetOk(dlVlan); ok {
				vlan := int64(d.Get(dlVlan).(int))
				gatewayPatchTemplateModel["vlan"] = &vlan
			}
		}
		if removeVlanOk, ok := d.GetOk(dlRemoveVlan); ok && !d.IsNewResource() {
			removeVlan := removeVlanOk.(bool)
			if removeVlan {
				gatewayPatchTemplateModel["vlan"] = nil
			}
		}
	}
	name := d.Get(dlName).(string)
	gatewayPatchTemplateModel["name"] = &name

	patchGatewayOptions := directLink.NewUpdateGatewayOptions(ID, gatewayPatchTemplateModel)
	_, response, err := directLink.UpdateGateway(patchGatewayOptions)
	if err != nil {
		log.Printf("[DEBUG] Update Direct Link Gateway err %s\n%s", err, response)
		return err
	}

	return resourceIBMdlGatewayRead(d, meta)
}

func resourceIBMdlGatewayDelete(d *schema.ResourceData, meta interface{}) error {

	directLink, err := directlinkClient(meta)
	if err != nil {
		return err
	}

	ID := d.Id()
	delOptions := &directlinkv1.DeleteGatewayOptions{
		ID: &ID,
	}
	response, err := directLink.DeleteGateway(delOptions)

	if err != nil && response.StatusCode != 404 {
		log.Printf("Error deleting Direct Link Gateway : %s", response)
		return err
	}

	d.SetId("")
	return nil
}

func resourceIBMdlGatewayExists(d *schema.ResourceData, meta interface{}) (bool, error) {
	directLink, err := directlinkClient(meta)
	if err != nil {
		return false, err
	}

	ID := d.Id()

	getOptions := &directlinkv1.GetGatewayOptions{
		ID: &ID,
	}
	instanceIntf, response, err := directLink.GetGateway(getOptions)

	if (err != nil) || (instanceIntf == nil) {
		if response != nil && response.StatusCode == 404 {
			d.SetId("")
			return false, nil
		}
		return false, fmt.Errorf("[ERROR] Error Getting Direct Link Gateway : %s\n%s", err, response)
	}
	_ = instanceIntf.(*directlinkv1.GetGatewayResponse)
	return true, nil
}
