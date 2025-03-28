// Copyright IBM Corp. 2017, 2021 All Rights Reserved.
// Licensed under the Mozilla Public License v2.0

package dnsservices

import (
	"fmt"
	"math/rand"
	"regexp"
	"strings"
	"time"

	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/conns"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

var allowedPrivateDomainRecordTypes = []string{
	"A", "AAAA", "CNAME", "MX", "PTR", "SRV", "TXT",
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

const (
	pdnsResourceRecordID = "resource_record_id"
	pdnsRecordType       = "type"
	pdnsRecordTTL        = "ttl"
	pdnsRecordName       = "name"
	pdnsRdata            = "rdata"
	pdnsMxPreference     = "preference"
	pdnsSrvPort          = "port"
	pdnsSrvPriority      = "priority"
	pdnsSrvWeight        = "weight"
	pdnsSrvProtocol      = "protocol"
	pdnsSrvService       = "service"
	pdnsRecordCreatedOn  = "created_on"
	pdnsRecordModifiedOn = "modified_on"
)

func caseDiffSuppress(_, old, new string, _ *schema.ResourceData) bool {
	return strings.ToUpper(old) == strings.ToUpper(new)
}

func ResourceIBMPrivateDNSResourceRecord() *schema.Resource {
	return &schema.Resource{
		Create:   resourceIBMPrivateDNSResourceRecordCreate,
		Read:     resourceIBMPrivateDNSResourceRecordRead,
		Update:   resourceIBMPrivateDNSResourceRecordUpdate,
		Delete:   resourceIBMPrivateDNSResourceRecordDelete,
		Exists:   resourceIBMPrivateDNSResourceRecordExists,
		Importer: &schema.ResourceImporter{},

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(10 * time.Minute),
			Update: schema.DefaultTimeout(10 * time.Minute),
			Delete: schema.DefaultTimeout(10 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			pdnsResourceRecordID: {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Resource record ID",
			},

			pdnsInstanceID: {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Instance ID",
			},

			pdnsZoneID: {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Zone ID",
			},

			pdnsRecordName: {
				Type:             schema.TypeString,
				Required:         true,
				DiffSuppressFunc: suppressPDNSRecordNameDiff,
				Description:      "DNS record name",
			},

			pdnsRecordType: {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
				ValidateFunc: func(val interface{}, field string) (warnings []string, errors []error) {
					value := val.(string)
					for _, rtype := range allowedPrivateDomainRecordTypes {
						if value == rtype {
							return
						}
					}

					errors = append(
						errors,
						fmt.Errorf("%s is not one of the valid domain record types: %s",
							value, strings.Join(allowedPrivateDomainRecordTypes, ", "),
						),
					)
					return
				},
				Description: "DNS record Type",
			},

			pdnsRdata: {
				Type:             schema.TypeString,
				Required:         true,
				DiffSuppressFunc: caseDiffSuppress,
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
				Description: "DNS record Data",
			},

			pdnsRecordTTL: {
				Type:     schema.TypeInt,
				Optional: true,
				Default:  900,
				DefaultFunc: func() (interface{}, error) {
					return 900, nil
				},
				Description: "DNS record TTL",
			},

			pdnsMxPreference: {
				Type:        schema.TypeInt,
				Optional:    true,
				Default:     0,
				Description: "DNS maximum preference",
			},

			pdnsSrvPort: {
				Type:        schema.TypeInt,
				Optional:    true,
				Description: "DNS server Port",
			},

			pdnsSrvPriority: {
				Type:        schema.TypeInt,
				Optional:    true,
				Default:     0,
				Description: "DNS server Priority",
			},

			pdnsSrvWeight: {
				Type:        schema.TypeInt,
				Optional:    true,
				Default:     0,
				Description: "DNS server weight",
			},

			pdnsSrvService: {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Service info",
			},

			pdnsSrvProtocol: {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Protocol",
			},

			pdnsRecordCreatedOn: {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Creation Data",
			},

			pdnsRecordModifiedOn: {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Modification date",
			},
		},
	}
}

func resourceIBMPrivateDNSResourceRecordCreate(d *schema.ResourceData, meta interface{}) error {
	sess, err := meta.(conns.ClientSession).PrivateDNSClientSession()
	if err != nil {
		return err
	}

	var (
		instanceID string
		zoneID     string
		recordType string
		name       string
		rdata      string
		service    string
		protocol   string
		ttl        int
		preference int
		port       int
		priority   int
		weight     int
	)

	instanceID = d.Get(pdnsInstanceID).(string)
	zoneID = d.Get(pdnsZoneID).(string)
	recordType = d.Get(pdnsRecordType).(string)
	name = d.Get(pdnsRecordName).(string)
	rdata = d.Get(pdnsRdata).(string)

	if v, ok := d.GetOk(pdnsRecordTTL); ok {
		ttl = v.(int)
	}

	createResourceRecordOptions := sess.NewCreateResourceRecordOptions(instanceID, zoneID, recordType)
	createResourceRecordOptions.SetName(name)
	createResourceRecordOptions.SetTTL(int64(ttl))

	switch recordType {
	case "A":
		resourceRecordAData, err := sess.NewResourceRecordInputRdataRdataARecord(rdata)
		if err != nil {
			return fmt.Errorf("[ERROR] Error creating pdns resource record A data:%s", err)
		}
		createResourceRecordOptions.SetRdata(resourceRecordAData)
	case "AAAA":
		resourceRecordAaaaData, err := sess.NewResourceRecordInputRdataRdataAaaaRecord(rdata)
		if err != nil {
			return fmt.Errorf("[ERROR] Error creating pdns resource record Aaaa data:%s", err)
		}
		createResourceRecordOptions.SetRdata(resourceRecordAaaaData)
	case "CNAME":
		resourceRecordCnameData, err := sess.NewResourceRecordInputRdataRdataCnameRecord(rdata)
		if err != nil {
			return fmt.Errorf("[ERROR] Error creating pdns resource record Cname data:%s", err)
		}
		createResourceRecordOptions.SetRdata(resourceRecordCnameData)
	case "PTR":
		resourceRecordPtrData, err := sess.NewResourceRecordInputRdataRdataPtrRecord(rdata)
		if err != nil {
			return fmt.Errorf("[ERROR] Error creating pdns resource record Ptr data:%s", err)
		}
		createResourceRecordOptions.SetRdata(resourceRecordPtrData)
	case "TXT":
		resourceRecordTxtData, err := sess.NewResourceRecordInputRdataRdataTxtRecord(rdata)
		if err != nil {
			return fmt.Errorf("[ERROR] Error creating pdns resource record Txt data:%s", err)
		}
		createResourceRecordOptions.SetRdata(resourceRecordTxtData)
	case "MX":
		if v, ok := d.GetOk(pdnsMxPreference); ok {
			preference = v.(int)
		}
		resourceRecordMxData, err := sess.NewResourceRecordInputRdataRdataMxRecord(rdata, int64(preference))
		if err != nil {
			return fmt.Errorf("[ERROR] Error creating pdns resource record Mx data:%s", err)
		}
		createResourceRecordOptions.SetRdata(resourceRecordMxData)
	case "SRV":
		if v, ok := d.GetOk(pdnsSrvPort); ok {
			port = v.(int)
		}
		if v, ok := d.GetOk(pdnsSrvPriority); ok {
			priority = v.(int)
		}
		if v, ok := d.GetOk(pdnsSrvWeight); ok {
			weight = v.(int)
		}
		resourceRecordSrvData, err := sess.NewResourceRecordInputRdataRdataSrvRecord(int64(port), int64(priority), rdata, int64(weight))
		if err != nil {
			return fmt.Errorf("[ERROR] Error creating pdns resource record Srv data:%s", err)
		}
		if v, ok := d.GetOk(pdnsSrvService); ok {
			service = v.(string)
		}
		if v, ok := d.GetOk(pdnsSrvProtocol); ok {
			protocol = v.(string)
		}
		createResourceRecordOptions.SetRdata(resourceRecordSrvData)
		createResourceRecordOptions.SetService(service)
		createResourceRecordOptions.SetProtocol(protocol)
	}
	rand.Seed(time.Now().UnixNano())
	randI := fmt.Sprint(rand.Intn(50))
	mk := "private_dns_resource_record_" + instanceID + zoneID + randI
	conns.IbmMutexKV.Lock(mk)
	defer conns.IbmMutexKV.Unlock(mk)
	response, detail, err := sess.CreateResourceRecord(createResourceRecordOptions)
	if err != nil {
		return fmt.Errorf("[ERROR] Error creating pdns resource record:%s\n%s", err, detail)
	}

	d.SetId(fmt.Sprintf("%s/%s/%s", instanceID, zoneID, *response.ID))

	return resourceIBMPrivateDNSResourceRecordRead(d, meta)
}

func resourceIBMPrivateDNSResourceRecordRead(d *schema.ResourceData, meta interface{}) error {
	idSet := strings.Split(d.Id(), "/")
	sess, err := meta.(conns.ClientSession).PrivateDNSClientSession()
	if err != nil {
		return err
	}
	getResourceRecordOptions := sess.NewGetResourceRecordOptions(idSet[0], idSet[1], idSet[2])
	response, detail, err := sess.GetResourceRecord(getResourceRecordOptions)
	if err != nil {
		return fmt.Errorf("[ERROR] Error reading pdns resource record:%s\n%s", err, detail)
	}

	// extract the record name by removing zone details
	var recordName string
	zone := strings.Split(idSet[1], ":")
	name := strings.Split(*response.Name, zone[0])
	name[0] = strings.Trim(name[0], ".")
	recordName = name[0]

	if *response.Type == "SRV" {
		// "_sip._udp.testsrv"
		temp := strings.Split(name[0], ".")
		recordName = temp[2]
	}

	d.Set(pdnsResourceRecordID, response.ID)
	d.Set(pdnsInstanceID, idSet[0])
	d.Set(pdnsZoneID, idSet[1])
	d.Set(pdnsRecordName, recordName)
	d.Set(pdnsRecordType, response.Type)
	d.Set(pdnsRecordTTL, response.TTL)
	d.Set(pdnsRecordCreatedOn, response.CreatedOn.String())
	d.Set(pdnsRecordModifiedOn, response.ModifiedOn.String())

	if *response.Type == "SRV" {
		data := response.Rdata
		d.Set(pdnsSrvPort, data["port"])
		d.Set(pdnsSrvPriority, data["priority"])
		d.Set(pdnsSrvWeight, data["weight"])
		d.Set(pdnsRdata, data["target"].(string))
		d.Set(pdnsSrvService, response.Service)
		d.Set(pdnsSrvProtocol, response.Protocol)
	}

	if *response.Type == "MX" {
		data := response.Rdata
		d.Set(pdnsMxPreference, data["preference"])
		d.Set(pdnsRdata, data["exchange"].(string))
	}
	if *response.Type == "A" || *response.Type == "AAAA" {
		data := response.Rdata
		d.Set(pdnsRdata, data["ip"].(string))
	}
	if *response.Type == "CNAME" {
		data := response.Rdata
		d.Set(pdnsRdata, data["cname"].(string))
	}
	if *response.Type == "PTR" {
		data := response.Rdata
		d.Set(pdnsRdata, data["ptrdname"].(string))
	}
	if *response.Type == "TXT" {
		data := response.Rdata
		d.Set(pdnsRdata, data["text"].(string))
	}

	return nil
}

func resourceIBMPrivateDNSResourceRecordUpdate(d *schema.ResourceData, meta interface{}) error {
	idSet := strings.Split(d.Id(), "/")

	sess, err := meta.(conns.ClientSession).PrivateDNSClientSession()
	if err != nil {
		return err
	}
	rand.Seed(time.Now().UnixNano())
	randI := fmt.Sprint(rand.Intn(50))
	mk := "private_dns_resource_record_" + idSet[0] + idSet[1] + randI
	conns.IbmMutexKV.Lock(mk)
	defer conns.IbmMutexKV.Unlock(mk)

	updateResourceRecordOptions := sess.NewUpdateResourceRecordOptions(idSet[0], idSet[1], idSet[2], "", nil)

	var rdata string

	if d.HasChange(pdnsRecordName) || d.HasChange(pdnsRecordTTL) || d.HasChange(pdnsRdata) ||
		d.HasChange(pdnsSrvPort) || d.HasChange(pdnsSrvPriority) ||
		d.HasChange(pdnsSrvWeight) || d.HasChange(pdnsSrvService) ||
		d.HasChange(pdnsSrvProtocol) || d.HasChange(pdnsMxPreference) {

		recordName := d.Get(pdnsRecordName).(string)
		recordType := d.Get(pdnsRecordType).(string)
		ttl := int64(d.Get(pdnsRecordTTL).(int))

		if recordType != "PTR" {
			updateResourceRecordOptions.SetName(recordName)
		}
		switch recordType {
		case "A":
			updateResourceRecordOptions.SetTTL(ttl)
			rdata = d.Get(pdnsRdata).(string)
			resourceRecordAData, err := sess.NewResourceRecordUpdateInputRdataRdataARecord(rdata)
			if err != nil {
				return fmt.Errorf("[ERROR] Error creating pdns resource record A data:%s", err)
			}
			updateResourceRecordOptions.SetRdata(resourceRecordAData)

		case "AAAA":
			updateResourceRecordOptions.SetTTL(ttl)
			rdata = d.Get(pdnsRdata).(string)
			resourceRecordAaaaData, err := sess.NewResourceRecordUpdateInputRdataRdataAaaaRecord(rdata)
			if err != nil {
				return fmt.Errorf("[ERROR] Error creating pdns resource record Aaaa data:%s", err)
			}
			updateResourceRecordOptions.SetRdata(resourceRecordAaaaData)

		case "CNAME":
			updateResourceRecordOptions.SetTTL(ttl)
			rdata = d.Get(pdnsRdata).(string)
			resourceRecordCnameData, err := sess.NewResourceRecordUpdateInputRdataRdataCnameRecord(rdata)
			if err != nil {
				return fmt.Errorf("[ERROR] Error creating pdns resource record Cname data:%s", err)
			}
			updateResourceRecordOptions.SetRdata(resourceRecordCnameData)

		case "PTR":
			updateResourceRecordOptions.SetTTL(ttl)

		case "TXT":
			updateResourceRecordOptions.SetTTL(ttl)
			rdata = d.Get(pdnsRdata).(string)
			resourceRecordTxtData, err := sess.NewResourceRecordUpdateInputRdataRdataTxtRecord(rdata)
			if err != nil {
				return fmt.Errorf("[ERROR] Error creating pdns resource record Txt data:%s", err)
			}
			updateResourceRecordOptions.SetRdata(resourceRecordTxtData)

		case "MX":
			updateResourceRecordOptions.SetTTL(ttl)
			rdata = d.Get(pdnsRdata).(string)
			preference := d.Get(pdnsMxPreference).(int)

			resourceRecordMxData, err := sess.NewResourceRecordUpdateInputRdataRdataMxRecord(rdata, int64(preference))
			if err != nil {
				return fmt.Errorf("[ERROR] Error creating pdns resource record Mx data:%s", err)
			}
			updateResourceRecordOptions.SetRdata(resourceRecordMxData)

		case "SRV":
			updateResourceRecordOptions.SetTTL(ttl)
			rdata = d.Get(pdnsRdata).(string)
			port := d.Get(pdnsSrvPort).(int)
			priority := d.Get(pdnsSrvPriority).(int)
			weight := d.Get(pdnsSrvWeight).(int)

			resourceRecordSrvData, err := sess.NewResourceRecordUpdateInputRdataRdataSrvRecord(int64(port), int64(priority), rdata, int64(weight))
			if err != nil {
				return fmt.Errorf("[ERROR] Error creating pdns resource record Srv data:%s", err)
			}
			updateResourceRecordOptions.SetRdata(resourceRecordSrvData)

			service := d.Get(pdnsSrvService).(string)
			protocol := d.Get(pdnsSrvProtocol).(string)
			updateResourceRecordOptions.SetService(service)
			updateResourceRecordOptions.SetProtocol(protocol)
		}

		_, detail, err := sess.UpdateResourceRecord(updateResourceRecordOptions)
		if err != nil {
			return fmt.Errorf("[ERROR] Error updating pdns resource record:%s\n%s", err, detail)
		}
	}

	return resourceIBMPrivateDNSResourceRecordRead(d, meta)
}

func resourceIBMPrivateDNSResourceRecordDelete(d *schema.ResourceData, meta interface{}) error {
	idSet := strings.Split(d.Id(), "/")

	sess, err := meta.(conns.ClientSession).PrivateDNSClientSession()
	if err != nil {
		return err
	}
	rand.Seed(time.Now().UnixNano())
	randI := fmt.Sprint(rand.Intn(50))
	deleteResourceRecordOptions := sess.NewDeleteResourceRecordOptions(idSet[0], idSet[1], idSet[2])
	mk := "private_dns_resource_record_" + idSet[0] + idSet[1] + randI
	conns.IbmMutexKV.Lock(mk)
	defer conns.IbmMutexKV.Unlock(mk)
	response, err := sess.DeleteResourceRecord(deleteResourceRecordOptions)
	if err != nil {
		return fmt.Errorf("[ERROR] Error deleting pdns resource record:%s\n%s", err, response)
	}

	d.SetId("")
	return nil
}

func resourceIBMPrivateDNSResourceRecordExists(d *schema.ResourceData, meta interface{}) (bool, error) {
	sess, err := meta.(conns.ClientSession).PrivateDNSClientSession()
	if err != nil {
		return false, err
	}

	idSet := strings.Split(d.Id(), "/")
	if len(idSet) < 3 {
		return false, fmt.Errorf("[ERROR] Incorrect ID %s: Id should be a combination of InstanceID/zoneID/recordID", d.Id())
	}
	rand.Seed(time.Now().UnixNano())
	randI := fmt.Sprint(rand.Intn(50))
	getResourceRecordOptions := sess.NewGetResourceRecordOptions(idSet[0], idSet[1], idSet[2])
	mk := "private_dns_resource_record_" + idSet[0] + idSet[1] + randI
	conns.IbmMutexKV.Lock(mk)
	defer conns.IbmMutexKV.Unlock(mk)
	_, response, err := sess.GetResourceRecord(getResourceRecordOptions)

	if err != nil {
		if response != nil && response.StatusCode == 404 {
			return false, nil
		}
		return false, err
	}
	return true, nil
}

func suppressPDNSRecordNameDiff(k, old, new string, d *schema.ResourceData) bool {
	// PDNS concantenates name with domain. So just check name is the same
	if strings.ToUpper(strings.SplitN(old, ".", 2)[0]) == strings.ToUpper(strings.SplitN(new, ".", 2)[0]) {
		return true
	}
	// If name is @, its replaced by the domain name. So ignore check.
	if new == "@" {
		return true
	}

	return false
}
