// Copyright IBM Corp. 2017, 2021 All Rights Reserved.
// Licensed under the Mozilla Public License v2.0

package cis

import (
	"fmt"
	"log"

	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/conns"
	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/flex"
	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/validate"
	"github.com/IBM/go-sdk-core/v5/core"
	cisrangeappv1 "github.com/IBM/networking-go-sdk/rangeapplicationsv1"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

const (
	ibmCISRangeApp                     = "ibm_cis_range_app"
	cisRangeAppID                      = "app_id"
	cisRangeAppProtocol                = "protocol"
	cisRangeAppDNS                     = "dns"
	cisRangeAppDNSType                 = "dns_type"
	cisRangeAppOriginDirect            = "origin_direct"
	cisRangeAppOriginDNS               = "origin_dns"
	cisRangeAppOriginPort              = "origin_port"
	cisRangeAppIPFirewall              = "ip_firewall"
	cisRangeAppProxyProtocol           = "proxy_protocol"
	cisRangeAppProxyProtocolOff        = "off"
	cisRangeAppProxyProtocolV1         = "v1"
	cisRangeAppProxyProtocolV2         = "v2"
	cisRangeAppProxyProtocolSimple     = "simple"
	cisRangeAppEdgeIPsType             = "edge_ips_type"
	cisRangeAppEdgeIPsTypeDynamic      = "dynamic"
	cisRangeAppEdgeIPsConnectivity     = "edge_ips_connectivity"
	cisRangeAppEdgeIPsConnectivityIPv4 = "ipv4"
	cisRangeAppEdgeIPsConnectivityIPv6 = "ipv6"
	cisRangeAppEdgeIPsConnectivityAll  = "all"
	cisRangeAppTrafficType             = "traffic_type"
	cisRangeAppTrafficTypeDirect       = "direct"
	cisRangeAppTrafficTypeHTTP         = "http"
	cisRangeAppTrafficTypeHTTPS        = "https"
	cisRangeAppTLS                     = "tls"
	cisRangeAppTLSOff                  = "off"
	cisRangeAppTLSFlexible             = "flexible"
	cisRangeAppTLSFull                 = "full"
	cisRangeAppTLSStrict               = "strict"
	cisRangeAppCreatedOn               = "created_on"
	cisRangeAppModifiedOn              = "modified_on"
)

func ResourceIBMCISRangeApp() *schema.Resource {
	return &schema.Resource{
		Create:   ResourceIBMCISRangeAppCreate,
		Read:     ResourceIBMCISRangeAppRead,
		Update:   ResourceIBMCISRangeAppUpdate,
		Delete:   ResourceIBMCISRangeAppDelete,
		Exists:   ResourceIBMCISRangeAppExists,
		Importer: &schema.ResourceImporter{},
		Schema: map[string]*schema.Schema{
			cisID: {
				Type:        schema.TypeString,
				Required:    true,
				Description: "CIS Intance CRN",
				ValidateFunc: validate.InvokeValidator(ibmCISRangeApp,
					"cis_id"),
			},
			cisDomainID: {
				Type:             schema.TypeString,
				Required:         true,
				Description:      "CIS Domain ID",
				DiffSuppressFunc: suppressDomainIDDiff,
			},
			cisRangeAppID: {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Application identifier",
			},
			cisRangeAppProtocol: {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Defines the protocol and port for this application",
			},
			cisRangeAppDNS: {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Name of the DNS record for this application",
			},
			cisRangeAppDNSType: {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Type of the DNS record for this application",
			},
			cisRangeAppOriginDirect: {
				Type:         schema.TypeList,
				Optional:     true,
				AtLeastOneOf: []string{cisRangeAppOriginDirect, cisRangeAppOriginDNS},
				Description:  "IP address and port of the origin for this Range application.",
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			cisRangeAppOriginDNS: {
				Type:         schema.TypeString,
				Optional:     true,
				AtLeastOneOf: []string{cisRangeAppOriginDirect, cisRangeAppOriginDNS},
				Description:  "DNS record pointing to the origin for this Range application.",
			},
			cisRangeAppOriginPort: {
				Type:          schema.TypeInt,
				Optional:      true,
				ConflictsWith: []string{cisRangeAppOriginDirect},
				Description:   "Port at the origin that listens to traffic",
			},
			cisRangeAppIPFirewall: {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: "Enables the IP Firewall for this application. Only available for TCP applications.",
			},
			cisRangeAppProxyProtocol: {
				Type:         schema.TypeString,
				Optional:     true,
				Description:  "Allows for the true client IP to be passed to the service.",
				ValidateFunc: validate.InvokeValidator(ibmCISRangeApp, cisRangeAppProxyProtocol),
			},
			cisRangeAppEdgeIPsType: {
				Type:         schema.TypeString,
				Optional:     true,
				Default:      cisRangeAppEdgeIPsTypeDynamic,
				Description:  "The type of edge IP configuration.",
				ValidateFunc: validate.InvokeValidator(ibmCISRangeApp, cisRangeAppEdgeIPsType),
			},
			cisRangeAppEdgeIPsConnectivity: {
				Type:         schema.TypeString,
				Optional:     true,
				Default:      cisRangeAppEdgeIPsConnectivityAll,
				Description:  "Specifies the IP version.",
				ValidateFunc: validate.InvokeValidator(ibmCISRangeApp, cisRangeAppEdgeIPsConnectivity),
			},
			cisRangeAppTrafficType: {
				Type:         schema.TypeString,
				Optional:     true,
				Default:      cisRangeAppTrafficTypeDirect,
				Description:  "Configure how traffic is handled at the edge.",
				ValidateFunc: validate.InvokeValidator(ibmCISRangeApp, cisRangeAppTrafficType),
			},
			cisRangeAppTLS: {
				Type:         schema.TypeString,
				Optional:     true,
				Default:      cisRangeAppTLSOff,
				Description:  "Configure if and how TLS connections are terminated at the edge.",
				ValidateFunc: validate.InvokeValidator(ibmCISRangeApp, cisRangeAppTLS),
			},
			cisRangeAppCreatedOn: {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "created on date",
			},
			cisRangeAppModifiedOn: {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "modified on date",
			},
		},
	}
}
func ResourceIBMCISRangeAppValidator() *validate.ResourceValidator {

	validateSchema := make([]validate.ValidateSchema, 0)
	proxyProtocol := "off, v1, v2, simple"
	connectivity := "ipv4, ipv6, all"
	trafficType := "direct, http, https"
	tls := "off, flexible, full, strict"
	validateSchema = append(validateSchema,
		validate.ValidateSchema{
			Identifier:                 "cis_id",
			ValidateFunctionIdentifier: validate.ValidateCloudData,
			Type:                       validate.TypeString,
			CloudDataType:              "ResourceInstance",
			CloudDataRange:             []string{"service:internet-svcs"},
			Required:                   true})
	validateSchema = append(validateSchema,
		validate.ValidateSchema{
			Identifier:                 cisRangeAppProxyProtocol,
			ValidateFunctionIdentifier: validate.ValidateAllowedStringValue,
			Type:                       validate.TypeString,
			Required:                   true,
			AllowedValues:              proxyProtocol})
	validateSchema = append(validateSchema,
		validate.ValidateSchema{
			Identifier:                 cisRangeAppEdgeIPsConnectivity,
			ValidateFunctionIdentifier: validate.ValidateAllowedStringValue,
			Type:                       validate.TypeString,
			Required:                   true,
			AllowedValues:              connectivity})
	validateSchema = append(validateSchema,
		validate.ValidateSchema{
			Identifier:                 cisRangeAppEdgeIPsType,
			ValidateFunctionIdentifier: validate.ValidateAllowedStringValue,
			Type:                       validate.TypeString,
			Required:                   true,
			AllowedValues:              "dynamic"})
	validateSchema = append(validateSchema,
		validate.ValidateSchema{
			Identifier:                 cisRangeAppTrafficType,
			ValidateFunctionIdentifier: validate.ValidateAllowedStringValue,
			Type:                       validate.TypeString,
			Required:                   true,
			AllowedValues:              trafficType})
	validateSchema = append(validateSchema,
		validate.ValidateSchema{
			Identifier:                 cisRangeAppTLS,
			ValidateFunctionIdentifier: validate.ValidateAllowedStringValue,
			Type:                       validate.TypeString,
			Required:                   true,
			AllowedValues:              tls})

	ibmCISRangeAppResourceValidator := validate.ResourceValidator{ResourceName: ibmCISRangeApp, Schema: validateSchema}
	return &ibmCISRangeAppResourceValidator
}
func ResourceIBMCISRangeAppCreate(d *schema.ResourceData, meta interface{}) error {
	cisClient, err := meta.(conns.ClientSession).CisRangeAppClientSession()
	if err != nil {
		return err
	}

	crn := d.Get(cisID).(string)
	zoneID, _, _ := flex.ConvertTftoCisTwoVar(d.Get(cisDomainID).(string))
	cisClient.Crn = core.StringPtr(crn)
	cisClient.ZoneIdentifier = core.StringPtr(zoneID)

	protocol := d.Get(cisRangeAppProtocol).(string)
	dns := d.Get(cisRangeAppDNS).(string)
	dnsType := d.Get(cisRangeAppDNSType).(string)

	dnsOpt := &cisrangeappv1.RangeAppReqDns{
		Type: &dnsType,
		Name: &dns,
	}

	opt := cisClient.NewCreateRangeAppOptions(protocol, dnsOpt)

	if v, ok := d.GetOk(cisRangeAppOriginDirect); ok {
		opt.SetOriginDirect(flex.ExpandStringList(v.([]interface{})))
	}
	if v, ok := d.GetOk(cisRangeAppOriginDNS); ok {
		originDNSOpt := &cisrangeappv1.RangeAppReqOriginDns{
			Name: core.StringPtr(v.(string)),
		}
		opt.SetOriginDns(originDNSOpt)
	}
	if v, ok := d.GetOk(cisRangeAppOriginPort); ok {
		opt.SetOriginPort(int64(v.(int)))
	}
	if v, ok := d.GetOkExists(cisRangeAppIPFirewall); ok {
		opt.SetIpFirewall(v.(bool))
	}
	if v, ok := d.GetOk(cisRangeAppProxyProtocol); ok {
		opt.SetProxyProtocol(v.(string))
	}
	edgeIPsOpt := &cisrangeappv1.RangeAppReqEdgeIps{
		Type:         core.StringPtr(cisRangeAppEdgeIPsTypeDynamic),
		Connectivity: core.StringPtr(cisRangeAppEdgeIPsConnectivityAll),
	}
	if v, ok := d.GetOk(cisRangeAppEdgeIPsType); ok {
		edgeIPsOpt.Type = core.StringPtr(v.(string))
	}
	if v, ok := d.GetOk(cisRangeAppEdgeIPsType); ok {
		edgeIPsOpt.Connectivity = core.StringPtr(v.(string))
	}
	if v, ok := d.GetOk(cisRangeAppTrafficType); ok {
		opt.SetTrafficType(v.(string))
	}
	if v, ok := d.GetOk(cisRangeAppTLS); ok {
		opt.SetTls(v.(string))
	}

	result, resp, err := cisClient.CreateRangeApp(opt)
	if err != nil {
		return fmt.Errorf("[ERROR] Failed to create range application: %v", resp)
	}
	d.SetId(flex.ConvertCisToTfThreeVar(*result.Result.ID, zoneID, crn))
	return ResourceIBMCISRangeAppRead(d, meta)
}

func ResourceIBMCISRangeAppRead(d *schema.ResourceData, meta interface{}) error {
	cisClient, err := meta.(conns.ClientSession).CisRangeAppClientSession()
	if err != nil {
		return err
	}

	rangeAppID, zoneID, crn, _ := flex.ConvertTfToCisThreeVar(d.Id())
	cisClient.Crn = core.StringPtr(crn)
	cisClient.ZoneIdentifier = core.StringPtr(zoneID)

	opt := cisClient.NewGetRangeAppOptions(rangeAppID)
	result, resp, err := cisClient.GetRangeApp(opt)
	if err != nil {
		return fmt.Errorf("[ERROR] Failed to read range application: %v", resp)
	}
	d.Set(cisID, crn)
	d.Set(cisDomainID, zoneID)
	d.Set(cisRangeAppID, result.Result.ID)
	d.Set(cisRangeAppProtocol, result.Result.Protocol)
	d.Set(cisRangeAppDNSType, result.Result.Dns.Type)
	d.Set(cisRangeAppDNS, result.Result.Dns.Name)
	d.Set(cisRangeAppOriginDirect, flex.FlattenStringList(result.Result.OriginDirect))
	d.Set(cisRangeAppProxyProtocol, result.Result.ProxyProtocol)
	d.Set(cisRangeAppIPFirewall, result.Result.IpFirewall)
	d.Set(cisRangeAppTrafficType, result.Result.TrafficType)
	d.Set(cisRangeAppEdgeIPsType, result.Result.EdgeIps.Type)
	d.Set(cisRangeAppEdgeIPsConnectivity, result.Result.EdgeIps.Connectivity)
	d.Set(cisRangeAppTLS, result.Result.Tls)
	d.Set(cisRangeAppCreatedOn, result.Result.CreatedOn.String())
	d.Set(cisRangeAppModifiedOn, result.Result.ModifiedOn.String())
	return nil
}

func ResourceIBMCISRangeAppUpdate(d *schema.ResourceData, meta interface{}) error {
	cisClient, err := meta.(conns.ClientSession).CisRangeAppClientSession()
	if err != nil {
		return err
	}

	if d.HasChange(cisRangeAppOriginDirect) ||
		d.HasChange(cisRangeAppOriginDNS) ||
		d.HasChange(cisRangeAppOriginPort) ||
		d.HasChange(cisRangeAppIPFirewall) ||
		d.HasChange(cisRangeAppProxyProtocol) ||
		d.HasChange(cisRangeAppEdgeIPsType) ||
		d.HasChange(cisRangeAppEdgeIPsConnectivity) ||
		d.HasChange(cisRangeAppTLS) ||
		d.HasChange(cisRangeAppTrafficType) {

		rangeAppID, zoneID, crn, _ := flex.ConvertTfToCisThreeVar(d.Id())
		cisClient.Crn = core.StringPtr(crn)
		cisClient.ZoneIdentifier = core.StringPtr(zoneID)
		protocol := d.Get(cisRangeAppProtocol).(string)
		dns := d.Get(cisRangeAppDNS).(string)
		dnsType := d.Get(cisRangeAppDNSType).(string)

		dnsOpt := &cisrangeappv1.RangeAppReqDns{
			Type: &dnsType,
			Name: &dns,
		}

		opt := cisClient.NewUpdateRangeAppOptions(rangeAppID, protocol, dnsOpt)

		if v, ok := d.GetOk(cisRangeAppOriginDirect); ok {
			opt.SetOriginDirect(flex.ExpandStringList(v.([]interface{})))
		}
		if v, ok := d.GetOk(cisRangeAppOriginDNS); ok {
			originDNSOpt := &cisrangeappv1.RangeAppReqOriginDns{
				Name: core.StringPtr(v.(string)),
			}
			opt.SetOriginDns(originDNSOpt)
		}
		if v, ok := d.GetOk(cisRangeAppOriginPort); ok {
			opt.SetOriginPort(int64(v.(int)))
		}
		if v, ok := d.GetOkExists(cisRangeAppIPFirewall); ok {
			opt.SetIpFirewall(v.(bool))
		}
		if v, ok := d.GetOk(cisRangeAppProxyProtocol); ok {
			opt.SetProxyProtocol(v.(string))
		}
		edgeIPsOpt := &cisrangeappv1.RangeAppReqEdgeIps{
			Type:         core.StringPtr(cisRangeAppEdgeIPsTypeDynamic),
			Connectivity: core.StringPtr(cisRangeAppEdgeIPsConnectivityAll),
		}
		if v, ok := d.GetOk(cisRangeAppEdgeIPsType); ok {
			edgeIPsOpt.Type = core.StringPtr(v.(string))
		}
		if v, ok := d.GetOk(cisRangeAppEdgeIPsType); ok {
			edgeIPsOpt.Connectivity = core.StringPtr(v.(string))
		}
		if v, ok := d.GetOk(cisRangeAppTrafficType); ok {
			opt.SetTrafficType(v.(string))
		}
		if v, ok := d.GetOk(cisRangeAppTLS); ok {
			opt.SetTls(v.(string))
		}
		_, resp, err := cisClient.UpdateRangeApp(opt)
		if err != nil {
			return fmt.Errorf("[ERROR] Failed to update range application: %v", resp)
		}
	}
	return ResourceIBMCISRangeAppRead(d, meta)
}

func ResourceIBMCISRangeAppDelete(d *schema.ResourceData, meta interface{}) error {
	cisClient, err := meta.(conns.ClientSession).CisRangeAppClientSession()
	if err != nil {
		return err
	}

	rangeAppID, zoneID, cisID, _ := flex.ConvertTfToCisThreeVar(d.Id())
	cisClient.Crn = core.StringPtr(cisID)
	cisClient.ZoneIdentifier = core.StringPtr(zoneID)
	opt := cisClient.NewDeleteRangeAppOptions(rangeAppID)
	_, resp, err := cisClient.DeleteRangeApp(opt)
	if err != nil {
		return fmt.Errorf("[ERROR] Failed to delete range application: %v", resp)
	}
	return nil
}

func ResourceIBMCISRangeAppExists(d *schema.ResourceData, meta interface{}) (bool, error) {
	cisClient, err := meta.(conns.ClientSession).CisRangeAppClientSession()
	if err != nil {
		return false, err
	}
	rangeAppID, zoneID, cisID, _ := flex.ConvertTfToCisThreeVar(d.Id())
	cisClient.Crn = core.StringPtr(cisID)
	cisClient.ZoneIdentifier = core.StringPtr(zoneID)
	opt := cisClient.NewGetRangeAppOptions(rangeAppID)
	_, resp, err := cisClient.GetRangeApp(opt)
	if err != nil {
		if resp != nil && resp.StatusCode == 404 {
			log.Println("range application is not found")
			return false, nil
		}
		return false, fmt.Errorf("[ERROR] Failed to getting existing range application: %v", err)
	}
	return true, nil
}
