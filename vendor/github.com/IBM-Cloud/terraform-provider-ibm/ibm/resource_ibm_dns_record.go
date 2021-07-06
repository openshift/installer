// Copyright IBM Corp. 2017, 2021 All Rights Reserved.
// Licensed under the Mozilla Public License v2.0

package ibm

import (
	"fmt"
	"log"
	"regexp"
	"strconv"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/softlayer/softlayer-go/datatypes"
	"github.com/softlayer/softlayer-go/services"
	"github.com/softlayer/softlayer-go/sl"
)

var allowedDomainRecordTypes = []string{
	"a", "aaaa", "cname", "mx", "ns", "ptr", "soa", "spf", "srv", "txt",
}
var ipv6Regexp *regexp.Regexp
var upcaseRegexp *regexp.Regexp

func init() {
	ipv6Regexp, _ = regexp.Compile(
		"[a-zA-Z0-9]{4}:[a-zA-Z0-9]{4}:[a-zA-Z0-9]{4}:[a-zA-Z0-9]{4}:" +
			"[a-zA-Z0-9]{4}:[a-zA-Z0-9]{4}:[a-zA-Z0-9]{4}:[a-zA-Z0-9]{4}",
	)
	upcaseRegexp, _ = regexp.Compile("[A-Z]")
}

func resourceIBMDNSRecord() *schema.Resource {
	return &schema.Resource{
		Exists:   resourceIBMDNSRecordExists,
		Create:   resourceIBMDNSRecordCreate,
		Read:     resourceIBMDNSRecordRead,
		Update:   resourceIBMDNSRecordUpdate,
		Delete:   resourceIBMDNSRecordDelete,
		Importer: &schema.ResourceImporter{},
		Schema: map[string]*schema.Schema{
			"data": {
				Type:     schema.TypeString,
				Required: true,
				ValidateFunc: func(val interface{}, field string) (warnings []string, errors []error) {
					value := val.(string)
					if ipv6Regexp.MatchString(value) && upcaseRegexp.MatchString(value) {
						errors = append(
							errors,
							fmt.Errorf(
								"IPv6 addresses in the data property cannot have upper case letters: %s",
								value,
							),
						)
					}
					return
				},
				Description: "DNS record data",
			},

			"domain_id": {
				Type:        schema.TypeInt,
				Required:    true,
				ForceNew:    true,
				Description: "Domain ID of dns record instance",
			},

			"expire": {
				Type:        schema.TypeInt,
				Optional:    true,
				Computed:    true,
				Description: "DNS record expiry info",
			},

			"host": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Hostname",
			},

			"mx_priority": {
				Type:        schema.TypeInt,
				Optional:    true,
				Default:     0,
				Description: "Maximum priority",
			},

			"refresh": {
				Type:        schema.TypeInt,
				Optional:    true,
				Computed:    true,
				Description: "refresh rate",
			},

			"responsible_person": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "Responsible person for DNS record",
			},

			"retry": {
				Type:        schema.TypeInt,
				Optional:    true,
				Computed:    true,
				Description: "Retry count",
			},

			"minimum_ttl": {
				Type:        schema.TypeInt,
				Optional:    true,
				Computed:    true,
				Description: "Minimun TTL configuration",
			},

			"ttl": {
				Type:     schema.TypeInt,
				Required: true,
				DefaultFunc: func() (interface{}, error) {
					return 86400, nil
				},
				Description: "TTL configuration",
			},

			"type": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
				ValidateFunc: func(val interface{}, field string) (warnings []string, errors []error) {
					value := val.(string)
					for _, rtype := range allowedDomainRecordTypes {
						if value == rtype {
							return
						}
					}

					errors = append(
						errors,
						fmt.Errorf("%s is not one of the valid domain record types: %s",
							value, strings.Join(allowedDomainRecordTypes, ", "),
						),
					)
					return
				},
				Description: "DNS record type",
			},

			"service": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "service info",
			},

			"protocol": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "protocol info",
			},

			"port": {
				Type:        schema.TypeInt,
				Optional:    true,
				Description: "port number",
			},

			"priority": {
				Type:        schema.TypeInt,
				Optional:    true,
				Default:     0,
				Description: "priority info",
			},

			"weight": {
				Type:        schema.TypeInt,
				Optional:    true,
				Default:     0,
				Description: "weight info",
			},

			"tags": {
				Type:        schema.TypeSet,
				Optional:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Set:         schema.HashString,
				Description: "tags set for the resource",
			},
		},
	}
}

//  Creates DNS Domain Resource Record
//  https://sldn.softlayer.com/reference/services/SoftLayer_Dns_Domain_ResourceRecord/createObject
func resourceIBMDNSRecordCreate(d *schema.ResourceData, meta interface{}) error {
	sess := meta.(ClientSession).SoftLayerSession()
	service := services.GetDnsDomainResourceRecordService(sess.SetRetries(0))

	opts := datatypes.Dns_Domain_ResourceRecord{
		Data:     sl.String(d.Get("data").(string)),
		DomainId: sl.Int(d.Get("domain_id").(int)),
		Host:     sl.String(d.Get("host").(string)),
		Ttl:      sl.Int(d.Get("ttl").(int)),
		Type:     sl.String(d.Get("type").(string)),
	}

	if expire, ok := d.GetOk("expire"); ok {
		opts.Expire = sl.Int(expire.(int))
	}

	if minimum, ok := d.GetOk("minimum_ttl"); ok {
		opts.Minimum = sl.Int(minimum.(int))
	}

	mxPriority := d.Get("mx_priority")
	opts.MxPriority = sl.Int(mxPriority.(int))

	if refresh, ok := d.GetOk("refresh"); ok {
		opts.Refresh = sl.Int(refresh.(int))
	}

	if responsiblePerson, ok := d.GetOk("responsible_person"); ok {
		opts.ResponsiblePerson = sl.String(responsiblePerson.(string))
	}

	if retry, ok := d.GetOk("retry"); ok {
		opts.Retry = sl.Int(retry.(int))
	}

	optsSrv := datatypes.Dns_Domain_ResourceRecord_SrvType{
		Dns_Domain_ResourceRecord: opts,
	}
	if *opts.Type == "srv" {
		if serviceName, ok := d.GetOk("service"); ok {
			optsSrv.Service = sl.String(serviceName.(string))
		}

		if protocol, ok := d.GetOk("protocol"); ok {
			optsSrv.Protocol = sl.String(protocol.(string))
		}

		optsSrv.Priority = sl.Int(d.Get("priority").(int))

		optsSrv.Weight = sl.Int(d.Get("weight").(int))

		if port, ok := d.GetOk("port"); ok {
			optsSrv.Port = sl.Int(port.(int))
		}
	}

	log.Printf("[INFO] Creating DNS Resource %s Record for '%d' dns domain", *opts.Type, d.Get("domain_id").(int))

	var err error
	var id int
	if *opts.Type == "srv" {
		var record datatypes.Dns_Domain_ResourceRecord_SrvType
		serviceSrv := services.GetDnsDomainResourceRecordSrvTypeService(sess)
		record, err = serviceSrv.CreateObject(&optsSrv)
		if record.Id != nil {
			id = *record.Id
		}
	} else {
		var record datatypes.Dns_Domain_ResourceRecord
		record, err = service.CreateObject(&opts)
		if record.Id != nil {
			id = *record.Id
		}
	}

	if err != nil {
		return fmt.Errorf("Error creating DNS Resource %s Record: %s", *opts.Type, err)
	}

	d.SetId(fmt.Sprintf("%d", id))

	log.Printf("[INFO] Dns Resource %s Record ID: %s", *opts.Type, d.Id())

	return resourceIBMDNSRecordRead(d, meta)
}

//  Reads DNS Domain Resource Record from SL system
//  https://sldn.softlayer.com/reference/services/SoftLayer_Dns_Domain_ResourceRecord/getObject
func resourceIBMDNSRecordRead(d *schema.ResourceData, meta interface{}) error {
	sess := meta.(ClientSession).SoftLayerSession()
	service := services.GetDnsDomainResourceRecordService(sess)

	id, err := strconv.Atoi(d.Id())
	if err != nil {
		return fmt.Errorf("Not a valid ID, must be an integer: %s", err)
	}
	result, err := service.Id(id).GetObject()
	if err != nil {
		return fmt.Errorf("Error retrieving DNS Resource Record: %s", err)
	}

	// Required fields
	d.Set("data", *result.Data)
	d.Set("domain_id", *result.DomainId)
	d.Set("host", *result.Host)
	d.Set("type", *result.Type)
	d.Set("ttl", *result.Ttl)

	// Optional fields
	d.Set("expire", sl.Get(result.Expire, nil))
	d.Set("minimum_ttl", sl.Get(result.Minimum, nil))
	d.Set("mx_priority", sl.Get(result.MxPriority, nil))
	d.Set("responsible_person", sl.Get(result.ResponsiblePerson, nil))
	d.Set("refresh", sl.Get(result.Refresh, nil))
	d.Set("retry", sl.Get(result.Retry, nil))

	if *result.Type == "srv" {
		d.Set("service", sl.Get(result.Service, nil))
		d.Set("protocol", sl.Get(result.Protocol, nil))
		d.Set("port", sl.Get(result.Port, nil))
		d.Set("priority", sl.Get(result.Priority, nil))
		d.Set("weight", sl.Get(result.Weight, nil))
	}

	return nil
}

//  Updates DNS Domain Resource Record in SL system
//  https://sldn.softlayer.com/reference/services/SoftLayer_Dns_Domain_ResourceRecord/editObject
func resourceIBMDNSRecordUpdate(d *schema.ResourceData, meta interface{}) error {
	sess := meta.(ClientSession).SoftLayerSession()
	service := services.GetDnsDomainResourceRecordService(sess)
	serviceNoRetry := services.GetDnsDomainResourceRecordService(sess.SetRetries(0))

	recordId, _ := strconv.Atoi(d.Id())
	record, err := service.Id(recordId).GetObject()
	if err != nil {
		return fmt.Errorf("Error retrieving DNS Resource Record: %s", err)
	}

	recordType := d.Get("type").(string)

	if data, ok := d.GetOk("data"); ok && d.HasChange("data") {
		record.Data = sl.String(data.(string))
	}

	if domain_id, ok := d.GetOk("domain_id"); ok && d.HasChange("domain_id") {
		record.DomainId = sl.Int(domain_id.(int))
	}

	if host, ok := d.GetOk("host"); ok && d.HasChange("host") {
		record.Host = sl.String(host.(string))
	}

	if ttl, ok := d.GetOk("ttl"); ok && d.HasChange("ttl") {
		record.Ttl = sl.Int(ttl.(int))
	}

	if expire, ok := d.GetOk("expire"); ok && d.HasChange("expire") {
		record.Expire = sl.Int(expire.(int))
	}

	if minimum_ttl, ok := d.GetOk("minimum_ttl"); ok && d.HasChange("minimum_ttl") {
		record.Minimum = sl.Int(minimum_ttl.(int))
	}

	if d.HasChange("mx_priority") {
		record.MxPriority = sl.Int(d.Get("mx_priority").(int))
	}

	if refresh, ok := d.GetOk("refresh"); ok && d.HasChange("refresh") {
		record.Refresh = sl.Int(refresh.(int))
	}

	if contact_email, ok := d.GetOk("responsible_person"); ok && d.HasChange("responsible_person") {
		record.ResponsiblePerson = sl.String(contact_email.(string))
	}

	if retry, ok := d.GetOk("retry"); ok && d.HasChange("retry") {
		record.Retry = sl.Int(retry.(int))
	}

	recordSrv := datatypes.Dns_Domain_ResourceRecord_SrvType{
		Dns_Domain_ResourceRecord: record,
	}
	if recordType == "srv" {
		if service, ok := d.GetOk("service"); ok && d.HasChange("service") {
			recordSrv.Service = sl.String(service.(string))
		}

		if d.HasChange("priority") {
			recordSrv.Priority = sl.Int(d.Get("priority").(int))
		}

		if protocol, ok := d.GetOk("protocol"); ok && d.HasChange("protocol") {
			recordSrv.Protocol = sl.String(protocol.(string))
		}

		if port, ok := d.GetOk("port"); ok && d.HasChange("port") {
			recordSrv.Port = sl.Int(port.(int))
		}

		if d.HasChange("weight") {
			recordSrv.Weight = sl.Int(d.Get("weight").(int))
		}
	}

	if recordType == "srv" {
		_, err = services.GetDnsDomainResourceRecordSrvTypeService(sess.SetRetries(0)).
			Id(recordId).EditObject(&recordSrv)
	} else {
		_, err = serviceNoRetry.Id(recordId).EditObject(&record)
	}

	if err != nil {
		return fmt.Errorf("Error editing DNS Resource %s Record %d: %s", recordType, recordId, err)
	}

	return nil
}

//  Deletes DNS Domain Resource Record in SL system
//  https://sldn.softlayer.com/reference/services/SoftLayer_Dns_Domain_ResourceRecord/deleteObject
func resourceIBMDNSRecordDelete(d *schema.ResourceData, meta interface{}) error {
	sess := meta.(ClientSession).SoftLayerSession()
	service := services.GetDnsDomainResourceRecordService(sess)

	id, err := strconv.Atoi(d.Id())
	if err != nil {
		return fmt.Errorf("Not a valid ID, must be an integer: %s", err)
	}

	_, err = service.Id(id).DeleteObject()

	if err != nil {
		return fmt.Errorf("Error deleting DNS Resource Record: %s", err)
	}

	return nil
}

// Exists function is called by refresh
// if the entity is absent - it is deleted from the .tfstate file
func resourceIBMDNSRecordExists(d *schema.ResourceData, meta interface{}) (bool, error) {
	sess := meta.(ClientSession).SoftLayerSession()
	service := services.GetDnsDomainResourceRecordService(sess)

	id, err := strconv.Atoi(d.Id())
	if err != nil {
		return false, fmt.Errorf("Not a valid ID, must be an integer: %s", err)
	}

	record, err := service.Id(id).GetObject()
	if err != nil {
		if apiErr, ok := err.(sl.Error); ok {
			if apiErr.StatusCode == 404 {
				return false, nil
			}
		}
		return false, fmt.Errorf("Error retrieving domain record info: %s", err)
	}
	return record.Id != nil && *record.Id == id, nil
}
