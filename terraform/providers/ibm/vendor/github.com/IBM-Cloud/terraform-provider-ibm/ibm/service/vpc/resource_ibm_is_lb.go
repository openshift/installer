// Copyright IBM Corp. 2017, 2021 All Rights Reserved.
// Licensed under the Mozilla Public License v2.0

package vpc

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/flex"
	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/validate"
	"github.com/IBM/vpc-go-sdk/vpcv1"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/customdiff"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

const (
	isLBAvailability                 = "availability"
	isLBAccessMode                   = "access_mode"
	isLBAccessModes                  = "access_modes"
	isLBInstanceGroupsSupported      = "instance_groups_supported"
	isLBSourceIPPersistenceSupported = "source_ip_session_persistence_supported"
	isLBName                         = "name"
	isLBStatus                       = "status"
	isLBCrn                          = "crn"
	isLBTags                         = "tags"
	isLBType                         = "type"
	isLBSubnets                      = "subnets"
	isLBHostName                     = "hostname"
	isLBPublicIPs                    = "public_ips"
	isLBPrivateIPs                   = "private_ips"
	isLBListeners                    = "listeners"
	isLBPools                        = "pools"
	isLBOperatingStatus              = "operating_status"
	isLBDeleting                     = "deleting"
	isLBDeleted                      = "done"
	isLBProvisioning                 = "provisioning"
	isLBProvisioningDone             = "done"
	isLBResourceGroup                = "resource_group"
	isLBProfile                      = "profile"
	isLBRouteMode                    = "route_mode"
	isLBUdpSupported                 = "udp_supported"
	isLBLogging                      = "logging"
	isLBSecurityGroups               = "security_groups"
	isLBSecurityGroupsSupported      = "security_group_supported"

	isLBAccessTags = "access_tags"
)

func ResourceIBMISLB() *schema.Resource {
	return &schema.Resource{
		Create:   resourceIBMISLBCreate,
		Read:     resourceIBMISLBRead,
		Update:   resourceIBMISLBUpdate,
		Delete:   resourceIBMISLBDelete,
		Exists:   resourceIBMISLBExists,
		Importer: &schema.ResourceImporter{},

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(30 * time.Minute),
			Update: schema.DefaultTimeout(30 * time.Minute),
			Delete: schema.DefaultTimeout(30 * time.Minute),
		},

		CustomizeDiff: customdiff.All(
			customdiff.Sequence(
				func(_ context.Context, diff *schema.ResourceDiff, v interface{}) error {
					return flex.ResourceTagsCustomizeDiff(diff)
				},
			),
			customdiff.Sequence(
				func(_ context.Context, diff *schema.ResourceDiff, v interface{}) error {
					return flex.ResourceRouteModeValidate(diff)
				}),
			customdiff.Sequence(
				func(_ context.Context, diff *schema.ResourceDiff, v interface{}) error {
					return flex.ResourceValidateAccessTags(diff, v)
				}),
		),

		Schema: map[string]*schema.Schema{

			isLBName: {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     false,
				ValidateFunc: validate.InvokeValidator("ibm_is_lb", isLBName),
				Description:  "Load Balancer name",
			},

			isLBType: {
				Type:         schema.TypeString,
				ForceNew:     true,
				Optional:     true,
				Default:      "public",
				ValidateFunc: validate.InvokeValidator("ibm_is_lb", isLBType),
				Description:  "Load Balancer type",
			},

			isLBAvailability: {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The availability of this load balancer",
			},
			isLBAccessMode: {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The access mode of this load balancer",
			},
			isLBInstanceGroupsSupported: {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "Indicates whether this load balancer supports instance groups.",
			},
			isLBSourceIPPersistenceSupported: {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "Indicates whether this load balancer supports source IP session persistence.",
			},
			"dns": {
				Type:        schema.TypeList,
				Optional:    true,
				MinItems:    1,
				MaxItems:    1,
				Description: "The DNS configuration for this load balancer.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"instance_crn": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "The CRN for this DNS instance",
						},
						"zone_id": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "The unique identifier of the DNS zone.",
						},
					},
				},
			},
			isLBStatus: {
				Type:     schema.TypeString,
				Computed: true,
			},
			"failsafe_policy_actions": &schema.Schema{
				Type:        schema.TypeList,
				Computed:    true,
				Description: "The supported `failsafe_policy.action` values for this load balancer's pools.",
				Elem:        &schema.Schema{Type: schema.TypeString},
			},
			isLBCrn: {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The CRN for this Load Balancer",
			},

			isLBOperatingStatus: {
				Type:     schema.TypeString,
				Computed: true,
			},

			isLBPublicIPs: {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},

			isLBPrivateIPs: {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			isLBPrivateIPDetail: {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "The private IP addresses assigned to this load balancer.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						isLBPrivateIpAddress: {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The IP address to reserve, which must not already be reserved on the subnet.",
						},
						isLBPrivateIpHref: {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The URL for this reserved IP",
						},
						isLBPrivateIpName: {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The user-defined name for this reserved IP. If unspecified, the name will be a hyphenated list of randomly-selected words. Names must be unique within the subnet the reserved IP resides in. ",
						},
						isLBPrivateIpId: {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Identifies a reserved IP by a unique property.",
						},
						isLBPrivateIpResourceType: {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The resource type",
						},
					},
				},
			},
			isLBSubnets: {
				Type:        schema.TypeSet,
				Required:    true,
				ForceNew:    false,
				MinItems:    1,
				MaxItems:    15,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Set:         schema.HashString,
				Description: "Load Balancer subnets list",
			},

			isLBSecurityGroups: {
				Type:        schema.TypeSet,
				Computed:    true,
				Optional:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Set:         schema.HashString,
				Description: "Load Balancer securitygroups list",
			},

			isLBSecurityGroupsSupported: {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "Security Group Supported for this Load Balancer",
			},

			isLBProfile: {
				Type:          schema.TypeString,
				Optional:      true,
				Computed:      true,
				ForceNew:      true,
				Description:   "The profile to use for this load balancer.",
				ValidateFunc:  validate.InvokeValidator("ibm_is_lb", isLBProfile),
				ConflictsWith: []string{isLBLogging},
			},

			isLBTags: {
				Type:     schema.TypeSet,
				Optional: true,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString, ValidateFunc: validate.InvokeValidator("ibm_is_lb", "tags")},
				Set:      flex.ResourceIBMVPCHash,
			},

			isLBAccessTags: {
				Type:        schema.TypeSet,
				Optional:    true,
				Computed:    true,
				Elem:        &schema.Schema{Type: schema.TypeString, ValidateFunc: validate.InvokeValidator("ibm_is_lb", "accesstag")},
				Set:         flex.ResourceIBMVPCHash,
				Description: "List of access management tags",
			},

			isLBResourceGroup: {
				Type:     schema.TypeString,
				ForceNew: true,
				Optional: true,
				Computed: true,
			},

			isLBRouteMode: {
				Type:        schema.TypeBool,
				ForceNew:    true,
				Optional:    true,
				Default:     false,
				Description: "Indicates whether route mode is enabled for this load balancer",
			},

			isLBUdpSupported: {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "Indicates whether this load balancer supports UDP.",
			},

			isLBHostName: {
				Type:     schema.TypeString,
				Computed: true,
			},

			isLBLogging: {
				Type:          schema.TypeBool,
				Optional:      true,
				Default:       false,
				Description:   "Logging of Load Balancer",
				ConflictsWith: []string{isLBProfile},
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

			flex.ResourceGroupName: {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The resource group name in which resource is provisioned",
			},

			"version": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func ResourceIBMISLBValidator() *validate.ResourceValidator {

	validateSchema := make([]validate.ValidateSchema, 0)
	lbtype := "public, private, private_path"
	isLBProfileAllowedValues := "network-fixed, network-private-path"

	validateSchema = append(validateSchema,
		validate.ValidateSchema{
			Identifier:                 isLBName,
			ValidateFunctionIdentifier: validate.ValidateRegexpLen,
			Type:                       validate.TypeString,
			Required:                   true,
			Regexp:                     `^([a-z]|[a-z][-a-z0-9]*[a-z0-9])$`,
			MinValueLength:             1,
			MaxValueLength:             63})
	validateSchema = append(validateSchema,
		validate.ValidateSchema{
			Identifier:                 isLBType,
			ValidateFunctionIdentifier: validate.ValidateAllowedStringValue,
			Type:                       validate.TypeString,
			Required:                   true,
			AllowedValues:              lbtype})
	validateSchema = append(validateSchema,
		validate.ValidateSchema{
			Identifier:                 isLBProfile,
			ValidateFunctionIdentifier: validate.ValidateAllowedStringValue,
			Type:                       validate.TypeString,
			Required:                   false,
			AllowedValues:              isLBProfileAllowedValues})
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
			Identifier:                 "accesstag",
			ValidateFunctionIdentifier: validate.ValidateRegexpLen,
			Type:                       validate.TypeString,
			Optional:                   true,
			Regexp:                     `^([A-Za-z0-9_.-]|[A-Za-z0-9_.-][A-Za-z0-9_ .-]*[A-Za-z0-9_.-]):([A-Za-z0-9_.-]|[A-Za-z0-9_.-][A-Za-z0-9_ .-]*[A-Za-z0-9_.-])$`,
			MinValueLength:             1,
			MaxValueLength:             128})

	ibmISLBResourceValidator := validate.ResourceValidator{ResourceName: "ibm_is_lb", Schema: validateSchema}
	return &ibmISLBResourceValidator
}

func resourceIBMISLBCreate(d *schema.ResourceData, meta interface{}) error {

	name := d.Get(isLBName).(string)
	subnets := d.Get(isLBSubnets).(*schema.Set)

	var isLogging bool
	if lbLogging, ok := d.GetOk(isLBLogging); ok {
		isLogging = lbLogging.(bool)
	}

	var securityGroups *schema.Set
	if sg, ok := d.GetOk(isLBSecurityGroups); ok {
		securityGroups = sg.(*schema.Set)
	}

	// subnets := flex.ExpandStringList((d.Get(isLBSubnets).(*schema.Set)).List())
	var lbType, rg string
	isPrivatePath := false
	isPublic := true
	if types, ok := d.GetOk(isLBType); ok {
		lbType = types.(string)
	}

	if lbType == "private" {
		isPublic = false
	}

	if lbType == "private_path" {
		isPrivatePath = true
		isPublic = false
	}

	if grp, ok := d.GetOk(isLBResourceGroup); ok {
		rg = grp.(string)
	}

	err := lbCreate(d, meta, name, lbType, rg, subnets, isPublic, isPrivatePath, isLogging, securityGroups)
	if err != nil {
		return err
	}

	return resourceIBMISLBRead(d, meta)
}

func lbCreate(d *schema.ResourceData, meta interface{}, name, lbType, rg string, subnets *schema.Set, isPublic, isPrivatePath, isLogging bool, securityGroups *schema.Set) error {
	sess, err := vpcClient(meta)
	if err != nil {
		return err
	}

	options := &vpcv1.CreateLoadBalancerOptions{
		IsPublic: &isPublic,
		Name:     &name,
	}
	if isPrivatePath {
		options.IsPrivatePath = &isPrivatePath
	}

	if dnsIntf, ok := d.GetOk("dns"); ok {
		dnsMap := dnsIntf.([]interface{})[0].(map[string]interface{})
		dnsInstance, _ := dnsMap["instance_crn"].(string)
		zone, _ := dnsMap["zone_id"].(string)
		dns := &vpcv1.LoadBalancerDnsPrototype{
			Instance: &vpcv1.DnsInstanceIdentity{
				CRN: &dnsInstance,
			},
			Zone: &vpcv1.DnsZoneIdentity{
				ID: &zone,
			},
		}
		options.Dns = dns
	}

	if routeModeBool, ok := d.GetOk(isLBRouteMode); ok {
		routeMode := routeModeBool.(bool)
		options.RouteMode = &routeMode
	}

	if subnets.Len() != 0 {
		subnetobjs := make([]vpcv1.SubnetIdentityIntf, subnets.Len())
		for i, subnet := range subnets.List() {
			subnetstr := subnet.(string)
			subnetobjs[i] = &vpcv1.SubnetIdentity{
				ID: &subnetstr,
			}
		}
		options.Subnets = subnetobjs
	}

	if securityGroups != nil && securityGroups.Len() != 0 {
		securityGroupobjs := make([]vpcv1.SecurityGroupIdentityIntf, securityGroups.Len())
		for i, securityGroup := range securityGroups.List() {
			securityGroupstr := securityGroup.(string)
			securityGroupobjs[i] = &vpcv1.SecurityGroupIdentity{
				ID: &securityGroupstr,
			}
		}
		options.SecurityGroups = securityGroupobjs
	}

	if rg != "" {
		options.ResourceGroup = &vpcv1.ResourceGroupIdentity{
			ID: &rg,
		}
	}

	if _, ok := d.GetOk(isLBProfile); ok {
		profile := d.Get(isLBProfile).(string)
		// Construct an instance of the LoadBalancerPoolIdentityByName model
		loadBalancerProfileIdentityModel := new(vpcv1.LoadBalancerProfileIdentityByName)
		loadBalancerProfileIdentityModel.Name = &profile
		options.Profile = loadBalancerProfileIdentityModel
	} else {

		dataPath := &vpcv1.LoadBalancerLoggingDatapathPrototype{
			Active: &isLogging,
		}
		loadBalancerLogging := &vpcv1.LoadBalancerLoggingPrototype{
			Datapath: dataPath,
		}
		options.Logging = loadBalancerLogging
	}

	lb, response, err := sess.CreateLoadBalancer(options)
	if err != nil {
		return fmt.Errorf("[ERROR] Error while creating Load Balancer err %s\n%s", err, response)
	}
	d.SetId(*lb.ID)
	log.Printf("[INFO] Load Balancer : %s", *lb.ID)
	_, err = isWaitForLBAvailable(sess, d.Id(), d.Timeout(schema.TimeoutCreate))
	if err != nil {
		return err
	}
	v := os.Getenv("IC_ENV_TAGS")
	if _, ok := d.GetOk(isLBTags); ok || v != "" {
		oldList, newList := d.GetChange(isLBTags)
		err = flex.UpdateGlobalTagsUsingCRN(oldList, newList, meta, *lb.CRN, "", isUserTagType)
		if err != nil {
			log.Printf(
				"Error on create of resource vpc Load Balancer (%s) tags: %s", d.Id(), err)
		}
	}

	if _, ok := d.GetOk(isLBAccessTags); ok {
		oldList, newList := d.GetChange(isLBAccessTags)
		err = flex.UpdateGlobalTagsUsingCRN(oldList, newList, meta, *lb.CRN, "", isAccessTagType)
		if err != nil {
			log.Printf(
				"Error on create of resource load balancer (%s) access tags: %s", d.Id(), err)
		}
	}
	return nil
}

func resourceIBMISLBRead(d *schema.ResourceData, meta interface{}) error {
	id := d.Id()

	err := lbGet(d, meta, id)
	if err != nil {
		return err
	}

	return nil
}

func lbGet(d *schema.ResourceData, meta interface{}, id string) error {
	sess, err := vpcClient(meta)
	if err != nil {
		return err
	}
	getLoadBalancerOptions := &vpcv1.GetLoadBalancerOptions{
		ID: &id,
	}
	lb, response, err := sess.GetLoadBalancer(getLoadBalancerOptions)
	if err != nil {
		if response != nil && response.StatusCode == 404 {
			d.SetId("")
			return nil
		}
		return fmt.Errorf("[ERROR] Error getting Load Balancer : %s\n%s", err, response)
	}
	if lb.Availability != nil {
		d.Set(isLBAvailability, *lb.Availability)
	}
	if lb.AccessMode != nil {
		d.Set(isLBAccessMode, *lb.AccessMode)
	}
	if lb.InstanceGroupsSupported != nil {
		d.Set(isLBInstanceGroupsSupported, *lb.InstanceGroupsSupported)
	}
	if lb.SourceIPSessionPersistenceSupported != nil {
		d.Set(isLBSourceIPPersistenceSupported, *lb.SourceIPSessionPersistenceSupported)
	}

	dnsList := make([]map[string]interface{}, 0)
	if lb.Dns != nil {
		dns := map[string]interface{}{}
		dns["instance_crn"] = lb.Dns.Instance.CRN
		dns["zone_id"] = lb.Dns.Zone.ID
		dnsList = append(dnsList, dns)
		d.Set("dns", dnsList)
	} else {
		d.Set("dns", nil)
	}
	d.Set(isLBName, *lb.Name)
	if lb.IsPublic != nil && *lb.IsPublic {
		d.Set(isLBType, "public")
	} else {
		if lb.IsPrivatePath != nil && *lb.IsPrivatePath {
			d.Set(isLBType, "private_path")
		} else {
			d.Set(isLBType, "private")
		}
	}
	if lb.RouteMode != nil {
		d.Set(isLBRouteMode, *lb.RouteMode)
	}

	if err = d.Set("failsafe_policy_actions", lb.FailsafePolicyActions); err != nil {
		err = fmt.Errorf("Error setting failsafe_policy_actions: %s", err)
		return flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_is_lb", "read", "set-failsafe_policy_actions")
	}
	d.Set(isLBStatus, *lb.ProvisioningStatus)
	d.Set(isLBCrn, *lb.CRN)
	d.Set(isLBOperatingStatus, *lb.OperatingStatus)
	publicIpList := make([]string, 0)
	if lb.PublicIps != nil {
		for _, ip := range lb.PublicIps {
			if ip.Address != nil {
				pubip := *ip.Address
				publicIpList = append(publicIpList, pubip)
			}
		}
	}
	d.Set(isLBPublicIPs, publicIpList)
	privateIpList := make([]string, 0)
	privateIpDetailList := make([]map[string]interface{}, 0)
	if lb.PrivateIps != nil {
		for _, ip := range lb.PrivateIps {
			if ip.Address != nil {
				prip := *ip.Address
				privateIpList = append(privateIpList, prip)
			}
			currentPriIp := map[string]interface{}{}
			if ip.Address != nil {
				currentPriIp[isLBPrivateIpAddress] = ip.Address
			}
			if ip.Href != nil {
				currentPriIp[isLBPrivateIpHref] = ip.Href
			}
			if ip.Name != nil {
				currentPriIp[isLBPrivateIpName] = ip.Name
			}
			if ip.ID != nil {
				currentPriIp[isLBPrivateIpId] = ip.ID
			}
			if ip.ResourceType != nil {
				currentPriIp[isLBPrivateIpResourceType] = ip.ResourceType
			}
			privateIpDetailList = append(privateIpDetailList, currentPriIp)
		}
	}
	// isLBPrivateIPs is same as isLBPrivateIPDetail.[].address
	d.Set(isLBPrivateIPs, privateIpList)
	d.Set(isLBPrivateIPDetail, privateIpDetailList)
	if lb.Subnets != nil {
		subnetList := make([]string, 0)
		for _, subnet := range lb.Subnets {
			if subnet.ID != nil {
				sub := *subnet.ID
				subnetList = append(subnetList, sub)
			}
		}
		d.Set(isLBSubnets, subnetList)
	}

	d.Set(isLBSecurityGroupsSupported, false)
	if lb.SecurityGroups != nil {
		securitygroupList := make([]string, 0)
		for _, SecurityGroup := range lb.SecurityGroups {
			if SecurityGroup.ID != nil {
				securityGroupID := *SecurityGroup.ID
				securitygroupList = append(securitygroupList, securityGroupID)
			}
		}
		d.Set(isLBSecurityGroups, securitygroupList)
		d.Set(isLBSecurityGroupsSupported, true)
	}

	if lb.Profile != nil {
		profile := lb.Profile
		if profile.Name != nil {
			d.Set(isLBProfile, *lb.Profile.Name)
		}
	} else {
		if lb.Logging != nil && lb.Logging.Datapath != nil && lb.Logging.Datapath.Active != nil {
			d.Set(isLBLogging, *lb.Logging.Datapath.Active)
		}
	}

	d.Set(isLBResourceGroup, *lb.ResourceGroup.ID)
	d.Set(isLBHostName, *lb.Hostname)
	if lb.UDPSupported != nil {
		d.Set(isLBUdpSupported, *lb.UDPSupported)
	}
	tags, err := flex.GetGlobalTagsUsingCRN(meta, *lb.CRN, "", isUserTagType)
	if err != nil {
		log.Printf(
			"Error on get of resource vpc Load Balancer (%s) tags: %s", d.Id(), err)
	}
	d.Set(isLBTags, tags)
	accesstags, err := flex.GetGlobalTagsUsingCRN(meta, *lb.CRN, "", isAccessTagType)
	if err != nil {
		log.Printf(
			"Error on get of resource load balancer (%s) access tags: %s", d.Id(), err)
	}
	d.Set(isLBAccessTags, accesstags)
	controller, err := flex.GetBaseController(meta)
	if err != nil {
		return err
	}
	d.Set(flex.ResourceControllerURL, controller+"/vpc-ext/network/loadBalancers")
	d.Set(flex.ResourceName, *lb.Name)
	if lb.ResourceGroup != nil {
		d.Set(flex.ResourceGroupName, lb.ResourceGroup.Name)
	}
	if err = d.Set("version", response.Headers.Get("Etag")); err != nil {
		return fmt.Errorf("[ERROR] Error setting version: %s", err)
	}
	return nil
}

func resourceIBMISLBUpdate(d *schema.ResourceData, meta interface{}) error {

	id := d.Id()
	name := ""
	isLogging := false
	hasChanged := false
	hasChangedLog := false
	var remove, add []string
	hasChangedSecurityGroups := false

	if d.HasChange(isLBName) {
		name = d.Get(isLBName).(string)
		hasChanged = true
	}
	if d.HasChange(isLBLogging) {
		isLogging = d.Get(isLBLogging).(bool)
		hasChangedLog = true
	}
	if d.HasChange(isLBSecurityGroups) {
		oldSecurityGroups, newSecurityGroups := d.GetChange(isLBSecurityGroups)
		oSecurityGroups := oldSecurityGroups.(*schema.Set)
		nSecurityGroups := newSecurityGroups.(*schema.Set)
		remove = flex.ExpandStringList(oSecurityGroups.Difference(nSecurityGroups).List())
		add = flex.ExpandStringList(nSecurityGroups.Difference(oSecurityGroups).List())
		hasChangedSecurityGroups = true
	}

	err := lbUpdate(d, meta, id, name, hasChanged, isLogging, hasChangedLog, hasChangedSecurityGroups, remove, add)
	if err != nil {
		return err
	}

	return resourceIBMISLBRead(d, meta)
}

func lbUpdate(d *schema.ResourceData, meta interface{}, id, name string, hasChanged bool, isLogging bool, hasChangedLog bool, hasChangedSecurityGroups bool, remove, add []string) error {
	sess, err := vpcClient(meta)
	if err != nil {
		return err
	}
	if d.HasChange(isLBTags) || d.HasChange(isLBAccessTags) {
		getLoadBalancerOptions := &vpcv1.GetLoadBalancerOptions{
			ID: &id,
		}
		lb, response, err := sess.GetLoadBalancer(getLoadBalancerOptions)
		if err != nil {
			return fmt.Errorf("[ERROR] Error getting Load Balancer : %s\n%s", err, response)
		}
		if d.HasChange(isLBTags) {
			oldList, newList := d.GetChange(isLBTags)
			err = flex.UpdateGlobalTagsUsingCRN(oldList, newList, meta, *lb.CRN, "", isUserTagType)
			if err != nil {
				log.Printf(
					"Error on update of resource vpc Load Balancer (%s) tags: %s", d.Id(), err)
			}
		}
		if d.HasChange(isLBAccessTags) {
			oldList, newList := d.GetChange(isLBAccessTags)
			err := flex.UpdateGlobalTagsUsingCRN(oldList, newList, meta, *lb.CRN, "", isAccessTagType)
			if err != nil {
				log.Printf(
					"Error on update of resource load balancer (%s) access tags: %s", d.Id(), err)
			}
		}
	}

	if d.HasChange("dns") {
		updateLoadBalancerOptions := &vpcv1.UpdateLoadBalancerOptions{
			ID: &id,
		}
		dnsRemoved := false
		if _, ok := d.GetOk("dns"); !ok {
			dnsRemoved = true
		}
		dnsPatchModel := &vpcv1.LoadBalancerDnsPatch{}
		if d.HasChange("dns.0.instance_crn") {
			dnsInstanceCrn := d.Get("dns.0.instance_crn").(string)
			dnsPatchModel.Instance = &vpcv1.DnsInstanceIdentity{
				CRN: &dnsInstanceCrn,
			}
		}
		if d.HasChange("dns.0.zone_id") {
			dnsZoneId := d.Get("dns.0.zone_id").(string)
			dnsPatchModel.Zone = &vpcv1.DnsZoneIdentity{
				ID: &dnsZoneId,
			}
		}

		loadBalancerPatchModel := &vpcv1.LoadBalancerPatch{
			Dns: dnsPatchModel,
		}
		loadBalancerPatch, err := loadBalancerPatchModel.AsPatch()
		if err != nil {
			return fmt.Errorf("[ERROR] Error calling asPatch for LoadBalancerPatch: %s", err)
		}
		if dnsRemoved {
			loadBalancerPatch["dns"] = nil
		}
		updateLoadBalancerOptions.LoadBalancerPatch = loadBalancerPatch

		_, response, err := sess.UpdateLoadBalancer(updateLoadBalancerOptions)
		if err != nil {
			return fmt.Errorf("[ERROR] Error Updating subents in vpc Load Balancer : %s\n%s", err, response)
		}
		_, err = isWaitForLBAvailable(sess, d.Id(), d.Timeout(schema.TimeoutUpdate))
		if err != nil {
			return err
		}
	}

	if d.HasChange(isLBSubnets) {
		updateLoadBalancerOptions := &vpcv1.UpdateLoadBalancerOptions{
			ID: &id,
		}
		updateLoadBalancerOptions.SetIfMatch(d.Get("version").(string))
		loadBalancerPatchModel := &vpcv1.LoadBalancerPatch{}
		subnets := d.Get(isLBSubnets).(*schema.Set)
		if subnets.Len() != 0 {
			subnetobjs := make([]vpcv1.SubnetIdentityIntf, subnets.Len())
			for i, subnet := range subnets.List() {
				subnetstr := subnet.(string)
				subnetobjs[i] = &vpcv1.SubnetIdentity{
					ID: &subnetstr,
				}
			}
			loadBalancerPatchModel.Subnets = subnetobjs
		}
		loadBalancerPatch, err := loadBalancerPatchModel.AsPatch()
		if err != nil {
			return fmt.Errorf("[ERROR] Error calling asPatch for LoadBalancerPatch: %s", err)
		}
		updateLoadBalancerOptions.LoadBalancerPatch = loadBalancerPatch

		_, response, err := sess.UpdateLoadBalancer(updateLoadBalancerOptions)
		if err != nil {
			return fmt.Errorf("[ERROR] Error Updating subents in vpc Load Balancer : %s\n%s", err, response)
		}
		_, err = isWaitForLBAvailable(sess, d.Id(), d.Timeout(schema.TimeoutUpdate))
		if err != nil {
			return err
		}
	}

	if hasChanged {
		updateLoadBalancerOptions := &vpcv1.UpdateLoadBalancerOptions{
			ID: &id,
		}
		loadBalancerPatchModel := &vpcv1.LoadBalancerPatch{
			Name: &name,
		}
		loadBalancerPatch, err := loadBalancerPatchModel.AsPatch()
		if err != nil {
			return fmt.Errorf("[ERROR] Error calling asPatch for LoadBalancerPatch: %s", err)
		}
		updateLoadBalancerOptions.LoadBalancerPatch = loadBalancerPatch

		_, response, err := sess.UpdateLoadBalancer(updateLoadBalancerOptions)
		if err != nil {
			return fmt.Errorf("[ERROR] Error Updating vpc Load Balancer : %s\n%s", err, response)
		}
	}
	if hasChangedLog {
		updateLoadBalancerOptions := &vpcv1.UpdateLoadBalancerOptions{
			ID: &id,
		}
		dataPath := &vpcv1.LoadBalancerLoggingDatapathPatch{
			Active: &isLogging,
		}
		loadBalancerLogging := &vpcv1.LoadBalancerLoggingPatch{
			Datapath: dataPath,
		}
		loadBalancerPatchModel := &vpcv1.LoadBalancerPatch{
			Logging: loadBalancerLogging,
		}
		loadBalancerPatch, err := loadBalancerPatchModel.AsPatch()
		if err != nil {
			return fmt.Errorf("[ERROR] Error calling asPatch for LoadBalancerPatch: %s", err)
		}
		updateLoadBalancerOptions.LoadBalancerPatch = loadBalancerPatch

		_, response, err := sess.UpdateLoadBalancer(updateLoadBalancerOptions)
		if err != nil {
			return fmt.Errorf("[ERROR] Error Updating vpc Load Balancer : %s\n%s", err, response)
		}
	}

	if hasChangedSecurityGroups {

		if len(add) > 0 {
			for _, securityGroupID := range add {
				createSecurityGroupTargetBindingOptions := &vpcv1.CreateSecurityGroupTargetBindingOptions{}
				createSecurityGroupTargetBindingOptions.SecurityGroupID = &securityGroupID
				createSecurityGroupTargetBindingOptions.ID = &id
				_, response, err := sess.CreateSecurityGroupTargetBinding(createSecurityGroupTargetBindingOptions)
				if err != nil {
					return fmt.Errorf("[ERROR] Error while creating Security Group Target Binding %s\n%s", err, response)
				}
				_, err = isWaitForLBAvailable(sess, d.Id(), d.Timeout(schema.TimeoutUpdate))
				if err != nil {
					return err
				}
			}
		}

		if len(remove) > 0 {
			for _, securityGroupID := range remove {
				getSecurityGroupTargetOptions := &vpcv1.GetSecurityGroupTargetOptions{
					SecurityGroupID: &securityGroupID,
					ID:              &id,
				}
				_, response, err := sess.GetSecurityGroupTarget(getSecurityGroupTargetOptions)
				if err != nil {
					if response != nil && response.StatusCode == 404 {
						continue
					}
					return fmt.Errorf("[ERROR] Error Getting Security Group Target for this load balancer (%s): %s\n%s", securityGroupID, err, response)
				}
				deleteSecurityGroupTargetBindingOptions := sess.NewDeleteSecurityGroupTargetBindingOptions(securityGroupID, id)
				response, err = sess.DeleteSecurityGroupTargetBinding(deleteSecurityGroupTargetBindingOptions)
				if err != nil {
					return fmt.Errorf("[ERROR] Error Deleting Security Group Target for this load balancer : %s\n%s", err, response)
				}
				_, err = isWaitForLBAvailable(sess, d.Id(), d.Timeout(schema.TimeoutUpdate))
				if err != nil {
					return err
				}
			}
		}
	}
	return nil
}

func resourceIBMISLBDelete(d *schema.ResourceData, meta interface{}) error {
	id := d.Id()

	err := lbDelete(d, meta, id)
	if err != nil {
		return err
	}

	d.SetId("")
	return nil
}

func lbDelete(d *schema.ResourceData, meta interface{}, id string) error {
	sess, err := vpcClient(meta)
	if err != nil {
		return err
	}

	getLoadBalancerOptions := &vpcv1.GetLoadBalancerOptions{
		ID: &id,
	}
	_, response, err := sess.GetLoadBalancer(getLoadBalancerOptions)
	if err != nil {
		if response != nil && response.StatusCode == 404 {
			d.SetId("")
			return nil
		}
		return fmt.Errorf("[ERROR] Error Getting vpc load balancer(%s): %s\n%s", id, err, response)
	}

	deleteLoadBalancerOptions := &vpcv1.DeleteLoadBalancerOptions{
		ID: &id,
	}
	response, err = sess.DeleteLoadBalancer(deleteLoadBalancerOptions)
	if err != nil {
		return fmt.Errorf("[ERROR] Error Deleting vpc load balancer : %s\n%s", err, response)
	}
	_, err = isWaitForLBDeleted(sess, id, d.Timeout(schema.TimeoutDelete))
	if err != nil {
		return err
	}
	d.SetId("")
	return nil
}

func isWaitForLBDeleted(lbc *vpcv1.VpcV1, id string, timeout time.Duration) (interface{}, error) {
	log.Printf("Waiting for  (%s) to be deleted.", id)

	stateConf := &resource.StateChangeConf{
		Pending:    []string{"retry", isLBDeleting},
		Target:     []string{isLBDeleted, "failed"},
		Refresh:    isLBDeleteRefreshFunc(lbc, id),
		Timeout:    timeout,
		Delay:      10 * time.Second,
		MinTimeout: 10 * time.Second,
	}

	return stateConf.WaitForState()
}

func isLBDeleteRefreshFunc(lbc *vpcv1.VpcV1, id string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		log.Printf("[DEBUG] is lb delete function here")
		getLoadBalancerOptions := &vpcv1.GetLoadBalancerOptions{
			ID: &id,
		}
		lb, response, err := lbc.GetLoadBalancer(getLoadBalancerOptions)
		if err != nil {
			if response != nil && response.StatusCode == 404 {
				return lb, isLBDeleted, nil
			}
			return nil, "failed", fmt.Errorf("[ERROR] The vpc load balancer %s failed to delete: %s\n%s", id, err, response)
		}
		return lb, isLBDeleting, nil
	}
}

func resourceIBMISLBExists(d *schema.ResourceData, meta interface{}) (bool, error) {
	id := d.Id()

	exists, err := lbExists(d, meta, id)
	return exists, err

}

func lbExists(d *schema.ResourceData, meta interface{}, id string) (bool, error) {
	sess, err := vpcClient(meta)
	if err != nil {
		return false, err
	}
	getLoadBalancerOptions := &vpcv1.GetLoadBalancerOptions{
		ID: &id,
	}
	_, response, err := sess.GetLoadBalancer(getLoadBalancerOptions)
	if err != nil {
		if response != nil && response.StatusCode == 404 {
			return false, nil
		}
		return false, fmt.Errorf("[ERROR] Error getting vpc load balancer: %s\n%s", err, response)
	}
	return true, nil
}

func isWaitForLBAvailable(sess *vpcv1.VpcV1, lbId string, timeout time.Duration) (interface{}, error) {
	log.Printf("Waiting for load balancer (%s) to be available.", lbId)

	stateConf := &resource.StateChangeConf{
		Pending:    []string{"retry", isLBProvisioning, "update_pending"},
		Target:     []string{isLBProvisioningDone, ""},
		Refresh:    isLBRefreshFunc(sess, lbId),
		Timeout:    timeout,
		Delay:      10 * time.Second,
		MinTimeout: 10 * time.Second,
	}

	return stateConf.WaitForState()
}

func isLBRefreshFunc(sess *vpcv1.VpcV1, lbId string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {

		getlboptions := &vpcv1.GetLoadBalancerOptions{
			ID: &lbId,
		}
		lb, response, err := sess.GetLoadBalancer(getlboptions)
		if err != nil {
			return nil, "", fmt.Errorf("[ERROR] Error Getting Load Balancer : %s\n%s", err, response)
		}

		if *lb.ProvisioningStatus == "active" || *lb.ProvisioningStatus == "failed" {
			return lb, isLBProvisioningDone, nil
		}

		return lb, isLBProvisioning, nil
	}
}
