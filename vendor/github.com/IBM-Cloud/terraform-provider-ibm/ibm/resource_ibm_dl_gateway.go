// Copyright IBM Corp. 2017, 2021 All Rights Reserved.
// Licensed under the Mozilla Public License v2.0

package ibm

import (
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/IBM/networking-go-sdk/directlinkv1"
	"github.com/hashicorp/terraform-plugin-sdk/helper/customdiff"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

const (
	dlLoaRejectReason              = "loa_reject_reason"
	dlCustomerName                 = "customer_name"
	dlCarrierName                  = "carrier_name"
	dlResourceGroup                = "resource_group"
	dlBgpAsn                       = "bgp_asn"
	dlBgpBaseCidr                  = "bgp_base_cidr"
	dlBgpCerCidr                   = "bgp_cer_cidr"
	dlBgpIbmCidr                   = "bgp_ibm_cidr"
	dlCrossConnectRouter           = "cross_connect_router"
	dlGlobal                       = "global"
	dlLocationName                 = "location_name"
	dlName                         = "name"
	dlSpeedMbps                    = "speed_mbps"
	dlOperationalStatus            = "operational_status"
	dlBgpStatus                    = "bgp_status"
	dlLinkStatus                   = "link_status"
	dlType                         = "type"
	dlCrn                          = "crn"
	dlCreatedAt                    = "created_at"
	dlMetered                      = "metered"
	dlLocationDisplayName          = "location_display_name"
	dlBgpIbmAsn                    = "bgp_ibm_asn"
	dlCompletionNoticeRejectReason = "completion_notice_reject_reason"
	dlPort                         = "port"
	dlProviderAPIManaged           = "provider_api_managed"
	dlVlan                         = "vlan"
	dlTags                         = "tags"
	dlActive                       = "active"
	dlFallbackCak                  = "fallback_cak"
	dlPrimaryCak                   = "primary_cak"
	dlSakExpiryTime                = "sak_expiry_time"
	dlWindowSize                   = "window_size"
	dlMacSecConfig                 = "macsec_config"
	dlCipherSuite                  = "cipher_suite"
	dlConfidentialityOffset        = "confidentiality_offset"
	dlCryptographicAlgorithm       = "cryptographic_algorithm"
	dlKeyServerPriority            = "key_server_priority"
	dlMacSecConfigStatus           = "status"
	dlChangeRequest                = "change_request"
	dlGatewayProvisioning          = "configuring"
	dlGatewayProvisioningDone      = "provisioned"
	dlGatewayProvisioningRejected  = "create_rejected"
)

func resourceIBMDLGateway() *schema.Resource {
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
			func(diff *schema.ResourceDiff, v interface{}) error {
				return resourceTagsCustomizeDiff(diff)
			},
		),

		Schema: map[string]*schema.Schema{
			dlBgpAsn: {
				Type:        schema.TypeInt,
				Required:    true,
				ForceNew:    true,
				Description: "BGP ASN",
			},
			dlBgpBaseCidr: {
				Type:             schema.TypeString,
				Optional:         true,
				ForceNew:         false,
				DiffSuppressFunc: applyOnce,
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
				ValidateFunc: InvokeValidator("ibm_dl_gateway", dlName),
				// ValidateFunc: validateRegexpLen(1, 63, "^([a-zA-Z]|[a-zA-Z][-_a-zA-Z0-9]*[a-zA-Z0-9])$"),
			},
			dlCarrierName: {
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				Description: "Carrier name",
				// ValidateFunc: validateRegexpLen(1, 128, "^[a-z][A-Z][0-9][ -_]$"),
			},
			dlCustomerName: {
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				Description: "Customer name",
				// ValidateFunc: validateRegexpLen(1, 128, "^[a-z][A-Z][0-9][ -_]$"),
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
				ValidateFunc: InvokeValidator("ibm_dl_gateway", dlType),
				// ValidateFunc: validateAllowedStringValue([]string{"dedicated", "connect"}),
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
				ForceNew:    true,
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
				ForceNew:    true,
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
			dlLocationDisplayName: {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Gateway location long name",
			},
			dlTags: {
				Type:        schema.TypeSet,
				Optional:    true,
				Computed:    true,
				Elem:        &schema.Schema{Type: schema.TypeString, ValidateFunc: InvokeValidator("ibm_dl_gateway", "tag")},
				Set:         resourceIBMVPCHash,
				Description: "Tags for the direct link gateway",
			},
			ResourceControllerURL: {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The URL of the IBM Cloud dashboard that can be used to explore and view details about this instance",
			},

			ResourceName: {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The name of the resource",
			},

			ResourceCRN: {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The crn of the resource",
			},

			ResourceStatus: {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The status of the resource",
			},

			ResourceGroupName: {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The resource group name in which resource is provisioned",
			},
		},
	}
}

func resourceIBMDLGatewayValidator() *ResourceValidator {

	validateSchema := make([]ValidateSchema, 2)
	dlTypeAllowedValues := "dedicated, connect"

	validateSchema = append(validateSchema,
		ValidateSchema{
			Identifier:                 dlType,
			ValidateFunctionIdentifier: ValidateAllowedStringValue,
			Type:                       TypeString,
			Required:                   true,
			AllowedValues:              dlTypeAllowedValues})
	validateSchema = append(validateSchema,
		ValidateSchema{
			Identifier:                 dlName,
			ValidateFunctionIdentifier: ValidateRegexpLen,
			Type:                       TypeString,
			Required:                   true,
			Regexp:                     `^([a-zA-Z]|[a-zA-Z][-_a-zA-Z0-9]*[a-zA-Z0-9])$`,
			MinValueLength:             1,
			MaxValueLength:             63})
	validateSchema = append(validateSchema,
		ValidateSchema{
			Identifier:                 "tag",
			ValidateFunctionIdentifier: ValidateRegexpLen,
			Type:                       TypeString,
			Optional:                   true,
			Regexp:                     `^[A-Za-z0-9:_ .-]+$`,
			MinValueLength:             1,
			MaxValueLength:             128})

	ibmISDLGatewayResourceValidator := ResourceValidator{ResourceName: "ibm_dl_gateway", Schema: validateSchema}
	return &ibmISDLGatewayResourceValidator
}

func directlinkClient(meta interface{}) (*directlinkv1.DirectLinkV1, error) {
	sess, err := meta.(ClientSession).DirectlinkV1API()
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

	if dtype == "dedicated" {
		var crossConnectRouter, carrierName, locationName, customerName string
		if _, ok := d.GetOk(dlCarrierName); ok {
			carrierName = d.Get(dlCarrierName).(string)
			//		gatewayTemplateModel.CarrierName = &carrierName
		} else {
			err = fmt.Errorf("Error creating gateway, %s is a required field", dlCarrierName)
			log.Printf("%s is a required field", dlCarrierName)
			return err
		}
		if _, ok := d.GetOk(dlCrossConnectRouter); ok {
			crossConnectRouter = d.Get(dlCrossConnectRouter).(string)
			//	gatewayTemplateModel.CrossConnectRouter = &crossConnectRouter
		} else {
			err = fmt.Errorf("Error creating gateway, %s is a required field", dlCrossConnectRouter)
			log.Printf("%s is a required field", dlCrossConnectRouter)
			return err
		}
		if _, ok := d.GetOk(dlLocationName); ok {
			locationName = d.Get(dlLocationName).(string)
			//gatewayTemplateModel.LocationName = &locationName
		} else {
			err = fmt.Errorf("Error creating gateway, %s is a required field", dlLocationName)
			log.Printf("%s is a required field", dlLocationName)
			return err
		}
		if _, ok := d.GetOk(dlCustomerName); ok {
			customerName = d.Get(dlCustomerName).(string)
			//gatewayTemplateModel.CustomerName = &customerName
		} else {
			err = fmt.Errorf("Error creating gateway, %s is a required field", dlCustomerName)
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
			createGatewayOptionsModel.GatewayTemplate = gatewayConnectTemplateModel

		} else {
			err = fmt.Errorf("Error creating direct link connect gateway, %s is a required field", dlPort)
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
		err = UpdateTagsUsingCRN(oldList, newList, meta, *gateway.Crn)
		if err != nil {
			log.Printf(
				"Error on create of resource direct link gateway %s (%s) tags: %s", dtype, d.Id(), err)
		}
	}

	return resourceIBMdlGatewayRead(d, meta)
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

	instance, response, err := directLink.GetGateway(getOptions)
	if err != nil {
		if response != nil && response.StatusCode == 404 {
			d.SetId("")
			return nil
		}
		return fmt.Errorf("Error Getting Direct Link Gateway (%s Template): %s\n%s", dtype, err, response)
	}
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
	if instance.CreatedAt != nil {
		d.Set(dlCreatedAt, instance.CreatedAt.String())
	}
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
	tags, err := GetTagsUsingCRN(meta, *instance.Crn)
	if err != nil {
		log.Printf(
			"Error on get of resource direct link gateway (%s) tags: %s", d.Id(), err)
	}
	d.Set(dlTags, tags)
	controller, err := getBaseController(meta)
	if err != nil {
		return err
	}
	d.Set(ResourceControllerURL, controller+"/interconnectivity/direct-link")
	d.Set(ResourceName, *instance.Name)
	d.Set(ResourceCRN, *instance.Crn)
	d.Set(ResourceStatus, *instance.OperationalStatus)
	if instance.ResourceGroup != nil {
		rg := instance.ResourceGroup
		d.Set(dlResourceGroup, *rg.ID)
		d.Set(ResourceGroupName, *rg.ID)
	}

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
		instance, response, err := client.GetGateway(getOptions)
		if err != nil {
			return nil, "", fmt.Errorf("Error Getting Direct Link: %s\n%s", err, response)
		}
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
	instance, detail, err := directLink.GetGateway(getOptions)

	if err != nil {
		log.Printf("Error fetching Direct Link Gateway :%s", detail)
		return err
	}

	updateGatewayOptionsModel := &directlinkv1.UpdateGatewayOptions{}
	updateGatewayOptionsModel.ID = &ID
	dtype := *instance.Type

	if d.HasChange(dlTags) {
		oldList, newList := d.GetChange(dlTags)
		err = UpdateTagsUsingCRN(oldList, newList, meta, *instance.Crn)
		if err != nil {
			log.Printf(
				"Error on update of resource direct link gateway (%s) tags: %s", *instance.ID, err)
		}
	}

	if d.HasChange(dlName) {
		name := d.Get(dlName).(string)
		updateGatewayOptionsModel.Name = &name
	}
	if d.HasChange(dlSpeedMbps) {
		speed := int64(d.Get(dlSpeedMbps).(int))
		updateGatewayOptionsModel.SpeedMbps = &speed
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
	if d.HasChange(dlGlobal) {
		global := d.Get(dlGlobal).(bool)
		updateGatewayOptionsModel.Global = &global
	}
	if d.HasChange(dlMetered) {
		metered := d.Get(dlMetered).(bool)
		updateGatewayOptionsModel.Metered = &metered
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
			updateGatewayOptionsModel.MacsecConfig = gatewayMacsecConfigTemplatePatchModel
		} else {
			updateGatewayOptionsModel.MacsecConfig = nil
		}
	}
	_, response, err := directLink.UpdateGateway(updateGatewayOptionsModel)
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
	_, response, err := directLink.GetGateway(getOptions)
	if err != nil {
		if response != nil && response.StatusCode == 404 {
			d.SetId("")
			return false, nil
		}
		return false, fmt.Errorf("Error Getting Direct Link Gateway : %s\n%s", err, response)
	}
	return true, nil
}
