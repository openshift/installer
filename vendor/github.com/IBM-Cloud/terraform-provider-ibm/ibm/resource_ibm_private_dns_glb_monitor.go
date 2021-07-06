// Copyright IBM Corp. 2017, 2021 All Rights Reserved.
// Licensed under the Mozilla Public License v2.0

package ibm

import (
	"fmt"
	"strings"
	"time"

	"github.com/IBM/networking-go-sdk/dnssvcsv1"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

const (
	ibmDNSGlbMonitor            = "ibm_dns_glb_monitor"
	pdnsGlbMonitorName          = "name"
	pdnsGlbMonitorID            = "monitor_id"
	pdnsGlbMonitorDescription   = "description"
	pdnsGlbMonitorType          = "type"
	pdnsGlbMonitorPort          = "port"
	pdnsGlbMonitorInterval      = "interval"
	pdnsGlbMonitorRetries       = "retries"
	pdnsGlbMonitorTimeout       = "timeout"
	pdnsGlbMonitorMethod        = "method"
	pdnsGlbMonitorPath          = "path"
	pdnsGlbMonitorAllowInsecure = "allow_insecure"
	pdnsGlbMonitorExpectedCodes = "expected_codes"
	pdnsGlbMonitorExpectedBody  = "expected_body"
	pdnsGlbMonitorHeaders       = "headers"
	pdnsGlbMonitorHeadersName   = "name"
	pdnsGlbMonitorHeadersValue  = "value"
	pdnsGlbMonitorCreatedOn     = "created_on"
	pdnsGlbMonitorModifiedOn    = "modified_on"
)

func resourceIBMPrivateDNSGLBMonitor() *schema.Resource {
	return &schema.Resource{
		Create:   resourceIBMPrivateDNSGLBMonitorCreate,
		Read:     resourceIBMPrivateDNSGLBMonitorRead,
		Update:   resourceIBMPrivateDNSGLBMonitorUpdate,
		Delete:   resourceIBMPrivateDNSGLBMonitorDelete,
		Exists:   resourceIBMPrivateDNSGLBMonitorExists,
		Importer: &schema.ResourceImporter{},

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(10 * time.Minute),
			Update: schema.DefaultTimeout(10 * time.Minute),
			Delete: schema.DefaultTimeout(10 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			pdnsGlbMonitorID: {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Monitor Id",
			},

			pdnsInstanceID: {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Instance Id",
			},

			pdnsGlbMonitorName: {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The unique identifier of a service instance.",
			},

			pdnsGlbMonitorDescription: {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Descriptive text of the load balancer monitor",
			},

			pdnsGlbMonitorType: {
				Type:         schema.TypeString,
				Optional:     true,
				Default:      "HTTP",
				ValidateFunc: InvokeValidator(ibmDNSGlbMonitor, pdnsGlbMonitorType),
				Description:  "The protocol to use for the health check",
			},

			pdnsGlbMonitorPort: {
				Type:        schema.TypeInt,
				Computed:    true,
				Optional:    true,
				Description: "Port number to connect to for the health check",
			},

			pdnsGlbMonitorInterval: {
				Type:        schema.TypeInt,
				Optional:    true,
				Default:     60,
				Description: "The interval between each health check",
			},

			pdnsGlbMonitorRetries: {
				Type:        schema.TypeInt,
				Optional:    true,
				Default:     1,
				Description: "The number of retries to attempt in case of a timeout before marking the origin as unhealthy",
			},

			pdnsGlbMonitorTimeout: {
				Type:        schema.TypeInt,
				Optional:    true,
				Default:     5,
				Description: "The timeout (in seconds) before marking the health check as failed",
			},

			pdnsGlbMonitorMethod: {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ValidateFunc: InvokeValidator(ibmDNSGlbMonitor, pdnsGlbMonitorMethod),
				Description:  "The method to use for the health check",
			},

			pdnsGlbMonitorPath: {
				Type:        schema.TypeString,
				Computed:    true,
				Optional:    true,
				Description: "The endpoint path to health check against",
			},

			pdnsGlbMonitorHeaders: {
				Type:        schema.TypeSet,
				Optional:    true,
				Description: "The HTTP request headers to send in the health check",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						pdnsGlbMonitorHeadersName: {
							Type:        schema.TypeString,
							Description: "The name of HTTP request header",
							Required:    true,
						},

						pdnsGlbMonitorHeadersValue: {
							Type:        schema.TypeList,
							Description: "The value of HTTP request header",
							Required:    true,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
					},
				},
			},

			pdnsGlbMonitorAllowInsecure: {
				Type:        schema.TypeBool,
				Computed:    true,
				Optional:    true,
				Description: "Do not validate the certificate when monitor use HTTPS. This parameter is currently only valid for HTTPS monitors.",
			},

			pdnsGlbMonitorExpectedCodes: {
				Type:         schema.TypeString,
				Computed:     true,
				Optional:     true,
				ValidateFunc: InvokeValidator(ibmDNSGlbMonitor, pdnsGlbMonitorExpectedCodes),
				Description:  "The expected HTTP response code or code range of the health check. This parameter is only valid for HTTP and HTTPS",
			},

			pdnsGlbMonitorExpectedBody: {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "A case-insensitive sub-string to look for in the response body",
			},

			pdnsGlbMonitorCreatedOn: {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "GLB Monitor creation date",
			},

			pdnsGlbMonitorModifiedOn: {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "GLB Monitor Modification date",
			},
		},
	}
}

func resourceIBMPrivateDNSGLBMonitorValidator() *ResourceValidator {
	monitorCheckTypes := "HTTP, HTTPS, TCP"
	methods := "GET, HEAD"
	expectedcode := "200,201,202,203,204,205,206,207,208,226,2xx"

	validateSchema := make([]ValidateSchema, 1)
	validateSchema = append(validateSchema,
		ValidateSchema{
			Identifier:                 pdnsGlbMonitorType,
			ValidateFunctionIdentifier: ValidateAllowedStringValue,
			Type:                       TypeString,
			Required:                   true,
			AllowedValues:              monitorCheckTypes})
	validateSchema = append(validateSchema,
		ValidateSchema{
			Identifier:                 pdnsGlbMonitorMethod,
			ValidateFunctionIdentifier: ValidateAllowedStringValue,
			Type:                       TypeString,
			Required:                   true,
			AllowedValues:              methods})
	validateSchema = append(validateSchema,
		ValidateSchema{
			Identifier:                 pdnsGlbMonitorExpectedCodes,
			ValidateFunctionIdentifier: ValidateAllowedStringValue,
			Type:                       TypeString,
			Required:                   true,
			AllowedValues:              expectedcode})
	dnsMonitorValidator := ResourceValidator{ResourceName: ibmDNSGlbMonitor, Schema: validateSchema}
	return &dnsMonitorValidator
}

func resourceIBMPrivateDNSGLBMonitorCreate(d *schema.ResourceData, meta interface{}) error {
	sess, err := meta.(ClientSession).PrivateDNSClientSession()
	if err != nil {
		return err
	}
	instanceID := d.Get(pdnsInstanceID).(string)
	createMonitorOptions := sess.NewCreateMonitorOptions(instanceID)

	monitorname := d.Get(pdnsGlbMonitorName).(string)
	monitorinterval := int64(d.Get(pdnsGlbMonitorInterval).(int))
	monitorretries := int64(d.Get(pdnsGlbMonitorRetries).(int))
	monitortimeout := int64(d.Get(pdnsGlbMonitorTimeout).(int))
	createMonitorOptions.SetName(monitorname)
	createMonitorOptions.SetInterval(monitorinterval)
	createMonitorOptions.SetRetries(monitorretries)
	createMonitorOptions.SetTimeout(monitortimeout)
	if monitordescription, ok := d.GetOk(pdnsGlbMonitorDescription); ok {
		createMonitorOptions.SetDescription(monitordescription.(string))
	}
	if Mtype, ok := d.GetOk(pdnsGlbMonitorType); ok {
		createMonitorOptions.SetType(Mtype.(string))
	}
	if Mport, ok := d.GetOk(pdnsGlbMonitorPort); ok {
		createMonitorOptions.SetPort(int64(Mport.(int)))
	}
	if monitorpath, ok := d.GetOk(pdnsGlbMonitorPath); ok {
		createMonitorOptions.SetPath((monitorpath).(string))
	}
	if monitorexpectedcodes, ok := d.GetOk(pdnsGlbMonitorExpectedCodes); ok {
		createMonitorOptions.SetExpectedCodes((monitorexpectedcodes).(string))
	}
	if monitormethod, ok := d.GetOk(pdnsGlbMonitorMethod); ok {
		createMonitorOptions.SetMethod((monitormethod).(string))
	}
	if monitorexpectedbody, ok := d.GetOk(pdnsGlbMonitorExpectedBody); ok {
		createMonitorOptions.SetExpectedBody((monitorexpectedbody).(string))
	}
	if monitorheaders, ok := d.GetOk(pdnsGlbMonitorHeaders); ok {
		expandedmonitorheaders, err := expandPDNSGLBMonitorsHeader(monitorheaders)
		if err != nil {
			return err
		}
		createMonitorOptions.SetHeadersVar(expandedmonitorheaders)
	}
	if monitorallowinsecure, ok := d.GetOkExists(pdnsGlbMonitorAllowInsecure); ok {
		createMonitorOptions.SetAllowInsecure((monitorallowinsecure).(bool))
	}

	response, detail, err := sess.CreateMonitor(createMonitorOptions)
	if err != nil {
		return fmt.Errorf("Error creating pdns GLB monitor:%s\n%s", err, detail)
	}

	d.SetId(fmt.Sprintf("%s/%s", instanceID, *response.ID))
	return resourceIBMPrivateDNSGLBMonitorRead(d, meta)
}

func expandPDNSGLBMonitorsHeader(header interface{}) ([]dnssvcsv1.HealthcheckHeader, error) {
	headers := header.(*schema.Set).List()
	expandheaders := make([]dnssvcsv1.HealthcheckHeader, 0)
	for _, v := range headers {
		locationConfig := v.(map[string]interface{})
		hname := locationConfig[pdnsGlbMonitorHeadersName].(string)
		headers := expandStringList(locationConfig[pdnsGlbMonitorHeadersValue].([]interface{}))
		headerItem := dnssvcsv1.HealthcheckHeader{
			Name:  &hname,
			Value: headers,
		}
		expandheaders = append(expandheaders, headerItem)
	}
	return expandheaders, nil
}

func resourceIBMPrivateDNSGLBMonitorRead(d *schema.ResourceData, meta interface{}) error {
	sess, err := meta.(ClientSession).PrivateDNSClientSession()
	if err != nil {
		return err
	}
	idset := strings.Split(d.Id(), "/")

	getMonitorOptions := sess.NewGetMonitorOptions(idset[0], idset[1])
	response, detail, err := sess.GetMonitor(getMonitorOptions)
	if err != nil {
		return fmt.Errorf("Error fetching pdns GLB Monitor:%s\n%s", err, detail)
	}
	d.Set(pdnsInstanceID, idset[0])
	d.Set(pdnsGlbMonitorID, response.ID)
	d.Set(pdnsGlbMonitorName, response.Name)
	d.Set(pdnsGlbMonitorCreatedOn, response.CreatedOn)
	d.Set(pdnsGlbMonitorModifiedOn, response.ModifiedOn)
	d.Set(pdnsGlbMonitorType, response.Type)
	d.Set(pdnsGlbMonitorPort, response.Port)
	if response.Path != nil {
		d.Set(pdnsGlbMonitorPath, response.Path)
	}
	if response.Interval != nil {
		d.Set(pdnsGlbMonitorInterval, response.Interval)
	}
	if response.Retries != nil {
		d.Set(pdnsGlbMonitorRetries, response.Retries)
	}
	if response.Timeout != nil {
		d.Set(pdnsGlbMonitorTimeout, response.Timeout)
	}
	if response.Method != nil {
		d.Set(pdnsGlbMonitorMethod, response.Method)
	}
	if response.ExpectedCodes != nil {
		d.Set(pdnsGlbMonitorExpectedCodes, response.ExpectedCodes)
	}
	if response.AllowInsecure != nil {
		d.Set(pdnsGlbMonitorAllowInsecure, response.AllowInsecure)
	}
	if response.Description != nil {
		d.Set(pdnsGlbMonitorDescription, response.Description)
	}
	if response.ExpectedBody != nil {
		d.Set(pdnsGlbMonitorExpectedBody, response.ExpectedBody)
	}

	d.Set(pdnsGlbMonitorHeaders, flattenDataSourceLoadBalancerHeader(response.HeadersVar))

	return nil
}

func flattenDataSourceLoadBalancerHeader(header []dnssvcsv1.HealthcheckHeader) interface{} {
	flattened := make([]interface{}, 0)

	for _, v := range header {
		cfg := map[string]interface{}{
			pdnsGlbMonitorHeadersName:  v.Name,
			pdnsGlbMonitorHeadersValue: flattenStringList(v.Value),
		}
		flattened = append(flattened, cfg)
	}
	return flattened
}

func resourceIBMPrivateDNSGLBMonitorUpdate(d *schema.ResourceData, meta interface{}) error {
	sess, err := meta.(ClientSession).PrivateDNSClientSession()
	if err != nil {
		return err
	}

	idset := strings.Split(d.Id(), "/")

	// Update PDNS GLB Monitor if attributes has any change

	if d.HasChange(pdnsGlbMonitorName) ||
		d.HasChange(pdnsGlbMonitorDescription) ||
		d.HasChange(pdnsGlbMonitorInterval) ||
		d.HasChange(pdnsGlbMonitorRetries) ||
		d.HasChange(pdnsGlbMonitorTimeout) ||
		d.HasChange(pdnsGlbMonitorExpectedBody) ||
		d.HasChange(pdnsGlbMonitorType) ||
		d.HasChange(pdnsGlbMonitorPort) ||
		d.HasChange(pdnsGlbMonitorPath) ||
		d.HasChange(pdnsGlbMonitorAllowInsecure) ||
		d.HasChange(pdnsGlbMonitorExpectedCodes) ||
		d.HasChange(pdnsGlbMonitorHeaders) {

		updateMonitorOptions := sess.NewUpdateMonitorOptions(idset[0], idset[1])
		uname := d.Get(pdnsGlbMonitorName).(string)
		udescription := d.Get(pdnsGlbMonitorDescription).(string)
		uinterval := int64(d.Get(pdnsGlbMonitorInterval).(int))
		uretries := int64(d.Get(pdnsGlbMonitorRetries).(int))
		utimeout := int64(d.Get(pdnsGlbMonitorTimeout).(int))
		updateMonitorOptions.SetName(uname)
		updateMonitorOptions.SetDescription(udescription)
		updateMonitorOptions.SetInterval(uinterval)
		updateMonitorOptions.SetRetries(uretries)
		updateMonitorOptions.SetTimeout(utimeout)

		if Mtype, ok := d.GetOk(pdnsGlbMonitorType); ok {
			updateMonitorOptions.SetType(Mtype.(string))
		}
		if Mport, ok := d.GetOk(pdnsGlbMonitorPort); ok {
			updateMonitorOptions.SetPort(int64(Mport.(int)))
		}
		if monitorpath, ok := d.GetOk(pdnsGlbMonitorPath); ok {
			updateMonitorOptions.SetPath((monitorpath).(string))
		}
		if monitorexpectedcodes, ok := d.GetOk(pdnsGlbMonitorExpectedCodes); ok {
			updateMonitorOptions.SetExpectedCodes((monitorexpectedcodes).(string))
		}
		if monitormethod, ok := d.GetOk(pdnsGlbMonitorMethod); ok {
			updateMonitorOptions.SetMethod((monitormethod).(string))
		}
		if monitorexpectedbody, ok := d.GetOk(pdnsGlbMonitorExpectedBody); ok {
			updateMonitorOptions.SetExpectedBody((monitorexpectedbody).(string))
		}
		if monitorheaders, ok := d.GetOk(pdnsGlbMonitorHeaders); ok {
			expandedmonitorheaders, err := expandPDNSGLBMonitorsHeader(monitorheaders)
			if err != nil {
				return err
			}
			updateMonitorOptions.SetHeadersVar(expandedmonitorheaders)
		}
		if monitorallowinsecure, ok := d.GetOkExists(pdnsGlbMonitorAllowInsecure); ok {
			updateMonitorOptions.SetAllowInsecure((monitorallowinsecure).(bool))
		}

		_, detail, err := sess.UpdateMonitor(updateMonitorOptions)

		if err != nil {
			return fmt.Errorf("Error updating pdns GLB Monitor:%s\n%s", err, detail)
		}
	}

	return resourceIBMPrivateDNSGLBMonitorRead(d, meta)
}

func resourceIBMPrivateDNSGLBMonitorDelete(d *schema.ResourceData, meta interface{}) error {
	sess, err := meta.(ClientSession).PrivateDNSClientSession()
	if err != nil {
		return err
	}

	idset := strings.Split(d.Id(), "/")

	DeleteMonitorOptions := sess.NewDeleteMonitorOptions(idset[0], idset[1])
	response, err := sess.DeleteMonitor(DeleteMonitorOptions)

	if err != nil {
		return fmt.Errorf("Error deleting pdns GLB Monitor:%s\n%s", err, response)
	}

	d.SetId("")
	return nil
}

func resourceIBMPrivateDNSGLBMonitorExists(d *schema.ResourceData, meta interface{}) (bool, error) {
	sess, err := meta.(ClientSession).PrivateDNSClientSession()
	if err != nil {
		return false, err
	}

	idset := strings.Split(d.Id(), "/")

	getMonitorOptions := sess.NewGetMonitorOptions(idset[0], idset[1])
	response, detail, err := sess.GetMonitor(getMonitorOptions)
	if err != nil {
		if response != nil && detail != nil && detail.StatusCode == 404 {
			return false, nil
		}
		return false, err
	}
	return true, nil
}
