// Copyright IBM Corp. 2022 All Rights Reserved.
// Licensed under the Mozilla Public License v2.0

package vpc

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/flex"
	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/validate"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/customdiff"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/IBM/go-sdk-core/v5/core"
	"github.com/IBM/vpc-go-sdk/vpcv1"
)

const (
	isVPNServerStatusPending  = "pending"
	isVPNServerStatusUpdating = "updating"
	isVPNServerStatusStable   = "stable"
	isVPNServerStatusFailed   = "failed"

	isVPNServerStatusDeleting = "deleting"
	isVPNServerStatusDeleted  = "deleted"

	isVPNServerAccessTags    = "access_tags"
	isVPNServerUserTagType   = "user"
	isVPNServerAccessTagType = "access"
)

func ResourceIBMIsVPNServer() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceIBMIsVPNServerCreate,
		ReadContext:   resourceIBMIsVPNServerRead,
		UpdateContext: resourceIBMIsVPNServerUpdate,
		DeleteContext: resourceIBMIsVPNServerDelete,
		Importer:      &schema.ResourceImporter{},

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(10 * time.Minute),
			Update: schema.DefaultTimeout(10 * time.Minute),
			Delete: schema.DefaultTimeout(10 * time.Minute),
		},

		CustomizeDiff: customdiff.All(
			customdiff.Sequence(
				func(_ context.Context, diff *schema.ResourceDiff, v interface{}) error {
					return flex.ResourceValidateAccessTags(diff, v)
				}),
		),

		Schema: map[string]*schema.Schema{
			"certificate_crn": &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    false,
				Description: "The crn of certificate instance for this VPN server.",
			},
			"client_authentication": &schema.Schema{
				Type:        schema.TypeList,
				Required:    true,
				ForceNew:    false,
				MaxItems:    2,
				Description: "The methods used to authenticate VPN clients to this VPN server. VPN clients must authenticate against all provided methods.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"method": &schema.Schema{
							Type:         schema.TypeString,
							Required:     true,
							ValidateFunc: validate.ValidateAllowedStringValues([]string{"certificate", "username"}),
							Description:  "The type of authentication.",
						},
						"identity_provider": &schema.Schema{
							Type:         schema.TypeString,
							Optional:     true,
							ValidateFunc: validate.ValidateAllowedStringValues([]string{"iam"}),
							Description:  "The type of identity provider to be used by the VPN client.- `iam`: IBM identity and access managementThe enumerated values for this property are expected to expand in the future. When processing this property, check for and log unknown values. Optionally halt processing and surface the error, or bypass the route on which the unexpected property value was encountered.",
						},
						"client_ca_crn": &schema.Schema{
							Type:        schema.TypeString,
							Optional:    true,
							Description: "The crn of certificate instance to use for the VPN client certificate authority (CA).",
						},
					},
				},
			},
			"client_auto_delete": &schema.Schema{
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "If set to `true`, disconnected VPN clients will be automatically deleted after the `client_auto_delete_timeout` time has passed.",
			},
			"client_auto_delete_timeout": &schema.Schema{
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "Hours after which disconnected VPN clients will be automatically deleted. If `0`, disconnected VPN clients will be deleted immediately.",
			},
			"client_dns_server_ips": &schema.Schema{
				Type:        schema.TypeSet,
				Optional:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Set:         schema.HashString,
				Description: "The DNS server addresses that will be provided to VPN clients connected to this VPN server. The IP address. This property may add support for IPv6 addresses in the future. When processing a value in this property, verify that the address is in an expected format. If it is not, log an error. Optionally halt processing and surface the error, or bypass the resource on which the unexpected IP address format was encountered.",
			},
			"client_idle_timeout": &schema.Schema{
				Type:         schema.TypeInt,
				Optional:     true,
				Default:      600,
				ValidateFunc: validate.InvokeValidator("ibm_is_vpn_server", "client_idle_timeout"),
				Description:  "The seconds a VPN client can be idle before this VPN server will disconnect it.   Specify `0` to prevent the server from disconnecting idle clients.",
			},
			"client_ip_pool": &schema.Schema{
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     false,
				ValidateFunc: validate.InvokeValidator("ibm_is_vpn_server", "client_ip_pool"),
				Description:  "The VPN client IPv4 address pool, expressed in CIDR format. The request must not overlap with any existing address prefixes in the VPC or any of the following reserved address ranges:  - `127.0.0.0/8` (IPv4 loopback addresses)  - `161.26.0.0/16` (IBM services)  - `166.8.0.0/14` (Cloud Service Endpoints)  - `169.254.0.0/16` (IPv4 link-local addresses)  - `224.0.0.0/4` (IPv4 multicast addresses)The prefix length of the client IP address pool's CIDR must be between`/9` (8,388,608 addresses) and `/22` (1024 addresses). A CIDR block that contains twice the number of IP addresses that are required to enable the maximum number of concurrent connections is recommended.",
			},
			"created_at": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The date and time that the VPN server was created.",
			},
			"crn": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The CRN for this VPN server.",
			},
			"enable_split_tunneling": &schema.Schema{
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
				Description: "Indicates whether the split tunneling is enabled on this VPN server.",
			},
			"health_state": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The health of this resource.- `ok`: Healthy- `degraded`: Suffering from compromised performance, capacity, or connectivity- `faulted`: Completely unreachable, inoperative, or otherwise entirely incapacitated- `inapplicable`: The health state does not apply because of the current lifecycle state. A resource with a lifecycle state of `failed` or `deleting` will have a health state of `inapplicable`. A `pending` resource may also have this state.",
			},
			"health_reasons": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"code": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "A snake case string succinctly identifying the reason for this health state.",
						},

						"message": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "An explanation of the reason for this health state.",
						},

						"more_info": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Link to documentation about the reason for this health state.",
						},
					},
				},
			},
			"hostname": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Fully qualified domain name assigned to this VPN server.",
			},
			"href": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The URL for this VPN server.",
			},
			"vpn_server": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The unique identifier for this VPN server.",
			},
			"lifecycle_state": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The lifecycle state of the VPN server.",
			},
			"lifecycle_reasons": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "The reasons for the current lifecycle_state (if any).",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"code": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "A snake case string succinctly identifying the reason for this lifecycle state.",
						},

						"message": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "An explanation of the reason for this lifecycle state.",
						},

						"more_info": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Link to documentation about the reason for this lifecycle state.",
						},
					},
				},
			},
			"name": &schema.Schema{
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validate.InvokeValidator("ibm_is_vpn_server", "name"),
				Description:  "The user-defined name for this VPN server. If unspecified, the name will be a hyphenated list of randomly-selected words. Names must be unique within the VPC this VPN server is serving.",
			},
			"port": &schema.Schema{
				Type:         schema.TypeInt,
				Optional:     true,
				Default:      443,
				ValidateFunc: validate.InvokeValidator("ibm_is_vpn_server", "port"),
				Description:  "The port number to use for this VPN server.",
			},
			"private_ips": &schema.Schema{
				Type:        schema.TypeList,
				Computed:    true,
				Description: "The reserved IPs bound to this VPN server.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"address": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The IP address. This property may add support for IPv6 addresses in the future. When processing a value in this property, verify that the address is in an expected format. If it is not, log an error. Optionally halt processing and surface the error, or bypass the resource on which the unexpected IP address format was encountered.",
						},
						"deleted": &schema.Schema{
							Type:        schema.TypeList,
							Computed:    true,
							Description: "If present, this property indicates the referenced resource has been deleted and providessome supplementary information.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"more_info": &schema.Schema{
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Link to documentation about deleted resources.",
									},
								},
							},
						},
						"href": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The URL for this reserved IP.",
						},
						"id": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The unique identifier for this reserved IP.",
						},
						"name": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The user-defined or system-provided name for this reserved IP.",
						},
						"resource_type": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The resource type.",
						},
					},
				},
			},

			"protocol": &schema.Schema{
				Type:         schema.TypeString,
				Optional:     true,
				Default:      "udp",
				ValidateFunc: validate.InvokeValidator("ibm_is_vpn_server", "protocol"),
				Description:  "The transport protocol to use for this VPN server.",
			},

			"resource_group": &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				Computed:    true,
				Description: "The unique identifier for this resource group. The resource group to use. If unspecified, the account's [default resourcegroup](https://cloud.ibm.com/apidocs/resource-manager#introduction) is used.",
			},

			"security_groups": {
				Type:        schema.TypeSet,
				Optional:    true,
				Computed:    true,
				ForceNew:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Set:         schema.HashString,
				Description: "The unique identifier for this security group. The security groups to use for this VPN server. If unspecified, the VPC's default security group is used.",
			},

			"subnets": {
				Type:     schema.TypeSet,
				Required: true,
				MinItems: 1,
				MaxItems: 2,
				// ForceNew:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Set:         schema.HashString,
				Description: "The unique identifier for this subnet. The subnets to provision this VPN server in.  Use subnets in different zones for high availability.",
			},

			"vpc": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "The VPC this VPN server resides in.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"crn": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The CRN for this VPC.",
						},
						"deleted": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "If present, this property indicates the referenced resource has been deleted and providessome supplementary information.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"more_info": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Link to documentation about deleted resources.",
									},
								},
							},
						},
						"href": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The URL for this VPC.",
						},
						"id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The unique identifier for this VPC.",
						},
						"name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The unique user-defined name for this VPC.",
						},
					},
				},
			},

			"resource_type": &schema.Schema{
				Type:         schema.TypeString,
				Default:      "vpn_server",
				Optional:     true,
				ValidateFunc: validate.ValidateAllowedStringValues([]string{"vpn_server"}),
				Description:  "The type of resource referenced.",
			},

			isVPNServerAccessTags: {
				Type:        schema.TypeSet,
				Optional:    true,
				Computed:    true,
				Elem:        &schema.Schema{Type: schema.TypeString, ValidateFunc: validate.InvokeValidator("ibm_is_vpn_server", "accesstag")},
				Set:         flex.ResourceIBMVPCHash,
				Description: "List of access management tags",
			},
		},
	}
}

func ResourceIBMIsVPNServerValidator() *validate.ResourceValidator {
	validateSchema := make([]validate.ValidateSchema, 0)
	validateSchema = append(validateSchema,
		validate.ValidateSchema{
			Identifier:                 "client_ip_pool",
			ValidateFunctionIdentifier: validate.ValidateRegexp,
			Type:                       validate.TypeString,
			Required:                   true,
			Regexp:                     `^(([0-9]|[1-9][0-9]|1[0-9]{2}|2[0-4][0-9]|25[0-5])\.){3}([0-9]|[1-9][0-9]|1[0-9]{2}|2[0-4][0-9]|25[0-5])(\/(3[0-2]|[1-2][0-9]|[0-9]))$`,
		},
		validate.ValidateSchema{
			Identifier:                 "client_idle_timeout",
			ValidateFunctionIdentifier: validate.IntBetween,
			Type:                       validate.TypeInt,
			Optional:                   true,
			MinValue:                   "0",
			MaxValue:                   "28800",
		},
		validate.ValidateSchema{
			Identifier:                 "name",
			ValidateFunctionIdentifier: validate.ValidateRegexpLen,
			Type:                       validate.TypeString,
			Optional:                   true,
			Regexp:                     `^([a-z]|[a-z][-a-z0-9]*[a-z0-9])$`,
			MinValueLength:             1,
			MaxValueLength:             63,
		},
		validate.ValidateSchema{
			Identifier:                 "port",
			ValidateFunctionIdentifier: validate.IntBetween,
			Type:                       validate.TypeInt,
			Optional:                   true,
			MinValue:                   "1",
			MaxValue:                   "65535",
		},
		validate.ValidateSchema{
			Identifier:                 "protocol",
			ValidateFunctionIdentifier: validate.ValidateAllowedStringValue,
			Type:                       validate.TypeString,
			Optional:                   true,
			AllowedValues:              "tcp, udp",
		},
		validate.ValidateSchema{
			Identifier:                 "method",
			ValidateFunctionIdentifier: validate.ValidateAllowedStringValue,
			Type:                       validate.TypeString,
			Required:                   true,
			AllowedValues:              "certificate , username",
		},
		validate.ValidateSchema{
			Identifier:                 "accesstag",
			ValidateFunctionIdentifier: validate.ValidateRegexpLen,
			Type:                       validate.TypeString,
			Optional:                   true,
			Regexp:                     `^([A-Za-z0-9_.-]|[A-Za-z0-9_.-][A-Za-z0-9_ .-]*[A-Za-z0-9_.-]):([A-Za-z0-9_.-]|[A-Za-z0-9_.-][A-Za-z0-9_ .-]*[A-Za-z0-9_.-])$`,
			MinValueLength:             1,
			MaxValueLength:             128,
		},
	)

	resourceValidator := validate.ResourceValidator{ResourceName: "ibm_is_vpn_server", Schema: validateSchema}
	return &resourceValidator
}

func resourceIBMIsVPNServerCreate(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	sess, err := vpcClient(meta)
	if err != nil {
		return diag.FromErr(err)
	}

	createVPNServerOptions := &vpcv1.CreateVPNServerOptions{}
	crn_val := d.Get("certificate_crn").(string)
	certificateInstanceIdentity := &vpcv1.CertificateInstanceIdentity{}
	certificateInstanceIdentity.CRN = &crn_val
	createVPNServerOptions.Certificate = certificateInstanceIdentity

	var clientAuthentication []vpcv1.VPNServerAuthenticationPrototypeIntf
	clientAuthArray := d.Get("client_authentication").([]interface{})
	for _, clientauth := range clientAuthArray {
		clientAuth := clientauth.(map[string]interface{})
		method := clientAuth["method"].(string)
		clientAuthPrototype := &vpcv1.VPNServerAuthenticationPrototype{}
		clientAuthPrototype.Method = &method

		if method == "certificate" {
			if clientAuth["client_ca_crn"] != nil {
				crn_val := clientAuth["client_ca_crn"].(string)
				certificateInstanceIdentity := &vpcv1.CertificateInstanceIdentity{}
				certificateInstanceIdentity.CRN = &crn_val
				clientAuthPrototype.ClientCa = certificateInstanceIdentity

			} else {
				return diag.FromErr(fmt.Errorf("[ERROR] Error method type `certificate` should be passed with `client_ca_crn`"))
			}
		} else {
			if clientAuth["identity_provider"] != nil {
				providerType := clientAuth["identity_provider"].(string)
				clientAuthPrototype.IdentityProvider = &vpcv1.VPNServerAuthenticationByUsernameIDProvider{
					ProviderType: &providerType,
				}
			} else {
				return diag.FromErr(fmt.Errorf("[ERROR] Error method type `username` should be passed with `identity_provider`"))
			}

		}
		clientAuthentication = append(clientAuthentication, clientAuthPrototype)
	}
	createVPNServerOptions.ClientAuthentication = clientAuthentication

	if _, ok := d.GetOk("client_dns_server_ips"); ok {
		var clientDnsServerIps []vpcv1.IP
		clientDnsServerIpsArray := d.Get("client_dns_server_ips").(*schema.Set)
		for _, val := range clientDnsServerIpsArray.List() {
			value := val.(string)
			address := &vpcv1.IP{
				Address: &value,
			}
			clientDnsServerIps = append(clientDnsServerIps, *address)
		}
		createVPNServerOptions.SetClientDnsServerIps(clientDnsServerIps)
	}

	if _, ok := d.GetOk("client_idle_timeout"); ok {
		createVPNServerOptions.SetClientIdleTimeout(int64(d.Get("client_idle_timeout").(int)))
	}

	createVPNServerOptions.SetClientIPPool(d.Get("client_ip_pool").(string))

	if _, ok := d.GetOk("enable_split_tunneling"); ok {
		createVPNServerOptions.SetEnableSplitTunneling(d.Get("enable_split_tunneling").(bool))
	}

	if _, ok := d.GetOk("name"); ok {
		createVPNServerOptions.SetName(d.Get("name").(string))
	}

	if _, ok := d.GetOk("port"); ok {
		createVPNServerOptions.SetPort(int64(d.Get("port").(int)))
	}

	if _, ok := d.GetOk("protocol"); ok {
		createVPNServerOptions.SetProtocol(d.Get("protocol").(string))
	}

	if rg, ok := d.GetOk("resource_group"); ok {
		resourceGroup := rg.(string)
		createVPNServerOptions.ResourceGroup = &vpcv1.ResourceGroupIdentity{
			ID: &resourceGroup,
		}
	}

	if _, ok := d.GetOk("security_groups"); ok {
		sg := d.Get("security_groups").(*schema.Set)
		var securityGroups []vpcv1.SecurityGroupIdentityIntf
		for _, val := range sg.List() {
			value := val.(string)
			securityGroupIdentity := &vpcv1.SecurityGroupIdentity{
				ID: &value,
			}
			securityGroups = append(securityGroups, securityGroupIdentity)
		}
		createVPNServerOptions.SetSecurityGroups(securityGroups)
	}

	var subnets []vpcv1.SubnetIdentityIntf
	subnetsArray := d.Get("subnets").(*schema.Set)
	for _, val := range subnetsArray.List() {
		value := val.(string)
		subnetIdentity := &vpcv1.SubnetIdentity{
			ID: &value,
		}
		subnets = append(subnets, subnetIdentity)
	}
	createVPNServerOptions.SetSubnets(subnets)

	vpnServer, response, err := sess.CreateVPNServerWithContext(context, createVPNServerOptions)
	if err != nil {
		log.Printf("[DEBUG] CreateVPNServerWithContext failed %s\n%s", err, response)
		return diag.FromErr(fmt.Errorf("[ERROR] CreateVPNServerWithContext failed %s\n%s", err, response))
	}

	d.SetId(*vpnServer.ID)

	_, err = isWaitForVPNServerStable(context, sess, d, d.Timeout(schema.TimeoutCreate))
	if err != nil {
		return diag.FromErr(fmt.Errorf("[ERROR] VPNServer failed %s\n", err))
	}

	if _, ok := d.GetOk(isVPNServerAccessTags); ok {
		oldList, newList := d.GetChange(isVPNServerAccessTags)
		err = flex.UpdateGlobalTagsUsingCRN(oldList, newList, meta, *vpnServer.CRN, "", isVPNServerAccessTagType)
		if err != nil {
			log.Printf(
				"Error on create of resource vpc (%s) access tags: %s", d.Id(), err)
		}
	}

	return resourceIBMIsVPNServerRead(context, d, meta)
}

func isWaitForVPNServerStable(context context.Context, sess *vpcv1.VpcV1, d *schema.ResourceData, timeout time.Duration) (interface{}, error) {

	log.Printf("Waiting for VPN Server(%s) to be stable.", d.Id())
	stateConf := &resource.StateChangeConf{
		Pending: []string{isVPNServerStatusPending, isVPNServerStatusUpdating},
		Target:  []string{isVPNServerStatusStable, isVPNServerStatusFailed},
		Refresh: func() (interface{}, string, error) {
			getVPNServerOptions := &vpcv1.GetVPNServerOptions{}

			getVPNServerOptions.SetID(d.Id())

			vpnServer, response, err := sess.GetVPNServerWithContext(context, getVPNServerOptions)
			if err != nil {
				log.Printf("[DEBUG] GetVPNServerWithContext failed %s\n%s", err, response)
				return vpnServer, "", fmt.Errorf("Error Getting VPC Server: %s\n%s", err, response)
			}

			if *vpnServer.LifecycleState == "stable" || *vpnServer.LifecycleState == "failed" {
				return vpnServer, *vpnServer.LifecycleState, nil
			}
			return vpnServer, *vpnServer.LifecycleState, nil
		},
		Timeout:    timeout,
		Delay:      10 * time.Second,
		MinTimeout: 10 * time.Second,
	}

	return stateConf.WaitForState()
}

func resourceIBMIsVPNServerRead(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	sess, err := vpcClient(meta)
	if err != nil {
		return diag.FromErr(err)
	}

	getVPNServerOptions := &vpcv1.GetVPNServerOptions{}

	getVPNServerOptions.SetID(d.Id())

	vpnServer, response, err := sess.GetVPNServerWithContext(context, getVPNServerOptions)
	if err != nil {
		if response != nil && response.StatusCode == 404 {
			d.SetId("")
			return nil
		}
		log.Printf("[DEBUG] GetVPNServerWithContext failed %s\n%s", err, response)
		return diag.FromErr(fmt.Errorf("[ERROR] GetVPNServerWithContext failed %s\n%s", err, response))
	}

	if err = d.Set("vpn_server", d.Id()); err != nil {
		return diag.FromErr(fmt.Errorf("[ERROR] Error setting vpn_server: %s", err))
	}

	if err = d.Set("certificate_crn", *vpnServer.Certificate.CRN); err != nil {
		return diag.FromErr(fmt.Errorf("[ERROR] Error setting certificate: %s", err))
	}

	vpnServerAuthenticationPrototypeArray := make([]interface{}, len(vpnServer.ClientAuthentication))
	for i, clientAuthenticationItem := range vpnServer.ClientAuthentication {
		vpnServerAuthenticationPrototype := make(map[string]interface{})
		vpnServerAuthentication := clientAuthenticationItem.(*vpcv1.VPNServerAuthentication)
		if vpnServerAuthentication != nil {
			if vpnServerAuthentication.Method != nil {
				vpnServerAuthenticationPrototype["method"] = *vpnServerAuthentication.Method
				if vpnServerAuthentication.ClientCa != nil && vpnServerAuthentication.ClientCa.CRN != nil {
					vpnServerAuthenticationPrototype["client_ca_crn"] = *vpnServerAuthentication.ClientCa.CRN
				}
				if vpnServerAuthentication.IdentityProvider != nil {
					vpnServerAuthenticationByUsernameIDProvider := vpnServerAuthentication.IdentityProvider.(*vpcv1.VPNServerAuthenticationByUsernameIDProvider)
					vpnServerAuthenticationPrototype["identity_provider"] = *vpnServerAuthenticationByUsernameIDProvider.ProviderType
				}
			}
		}
		vpnServerAuthenticationPrototypeArray[i] = vpnServerAuthenticationPrototype
	}
	if err = d.Set("client_authentication", vpnServerAuthenticationPrototypeArray); err != nil {
		return diag.FromErr(fmt.Errorf("[ERROR] Error setting client_authentication: %s", err))
	}

	if err = d.Set("client_ip_pool", *vpnServer.ClientIPPool); err != nil {
		return diag.FromErr(fmt.Errorf("[ERROR] Error setting client_ip_pool: %s", err))
	}

	subnets := make([]string, 0)
	for i := 0; i < len(vpnServer.Subnets); i++ {
		subnets = append(subnets, string(*(vpnServer.Subnets[i].ID)))
	}
	if err = d.Set("subnets", subnets); err != nil {
		return diag.FromErr(fmt.Errorf("[ERROR] Error setting subnets: %s", err))
	}

	if vpnServer.ClientDnsServerIps != nil {
		// clientDnsServerIps := []map[string]interface{}{}
		clientDnsServerIps := make([]string, 0)
		for _, clientDnsServerIpsItem := range vpnServer.ClientDnsServerIps {
			if clientDnsServerIpsItem.Address != nil {
				clientDnsServerIps = append(clientDnsServerIps, *clientDnsServerIpsItem.Address)
			}
		}
		if err = d.Set("client_dns_server_ips", clientDnsServerIps); err != nil {
			return diag.FromErr(fmt.Errorf("[ERROR] Error setting client_dns_server_ips: %s", err))
		}
	}
	if err = d.Set("client_idle_timeout", flex.IntValue(vpnServer.ClientIdleTimeout)); err != nil {
		return diag.FromErr(fmt.Errorf("[ERROR] Error setting client_idle_timeout: %s", err))
	}
	if err = d.Set("enable_split_tunneling", vpnServer.EnableSplitTunneling); err != nil {
		return diag.FromErr(fmt.Errorf("[ERROR] Error setting enable_split_tunneling: %s", err))
	}
	if err = d.Set("name", vpnServer.Name); err != nil {
		return diag.FromErr(fmt.Errorf("[ERROR] Error setting name: %s", err))
	}
	if err = d.Set("port", flex.IntValue(vpnServer.Port)); err != nil {
		return diag.FromErr(fmt.Errorf("[ERROR] Error setting port: %s", err))
	}
	if err = d.Set("protocol", vpnServer.Protocol); err != nil {
		return diag.FromErr(fmt.Errorf("[ERROR] Error setting protocol: %s", err))
	}
	if vpnServer.ResourceGroup != nil && vpnServer.ResourceGroup.ID != nil {
		if err = d.Set("resource_group", vpnServer.ResourceGroup.ID); err != nil {
			return diag.FromErr(fmt.Errorf("[ERROR] Error setting resource_group: %s", err))
		}
	}
	if vpnServer.SecurityGroups != nil {
		securityGroups := make([]string, 0)
		for _, securityGroupsItem := range vpnServer.SecurityGroups {
			if securityGroupsItem.ID != nil {
				securityGroups = append(securityGroups, *securityGroupsItem.ID)
			}
		}
		if err = d.Set("security_groups", securityGroups); err != nil {
			return diag.FromErr(fmt.Errorf("[ERROR] Error setting security_groups: %s", err))
		}
	}
	if err = d.Set("client_auto_delete", vpnServer.ClientAutoDelete); err != nil {
		return diag.FromErr(fmt.Errorf("[ERROR] Error setting client_auto_delete: %s", err))
	}
	if err = d.Set("client_auto_delete_timeout", flex.IntValue(vpnServer.ClientAutoDeleteTimeout)); err != nil {
		return diag.FromErr(fmt.Errorf("[ERROR] Error setting client_auto_delete_timeout: %s", err))
	}
	if err = d.Set("created_at", flex.DateTimeToString(vpnServer.CreatedAt)); err != nil {
		return diag.FromErr(fmt.Errorf("[ERROR] Error setting created_at: %s", err))
	}
	if err = d.Set("crn", vpnServer.CRN); err != nil {
		return diag.FromErr(fmt.Errorf("[ERROR] Error setting crn: %s", err))
	}
	if err = d.Set("health_state", vpnServer.HealthState); err != nil {
		return diag.FromErr(fmt.Errorf("[ERROR] Error setting health_state: %s", err))
	}
	if vpnServer.HealthReasons != nil {
		if err := d.Set("health_reasons", resourceVPNServerFlattenHealthReasons(vpnServer.HealthReasons)); err != nil {
			return diag.FromErr(fmt.Errorf("[ERROR] Error setting health_reasons: %s", err))
		}
	}
	if err = d.Set("hostname", vpnServer.Hostname); err != nil {
		return diag.FromErr(fmt.Errorf("[ERROR] Error setting hostname: %s", err))
	}
	if err = d.Set("href", vpnServer.Href); err != nil {
		return diag.FromErr(fmt.Errorf("[ERROR] Error setting href: %s", err))
	}
	if err = d.Set("lifecycle_state", vpnServer.LifecycleState); err != nil {
		return diag.FromErr(fmt.Errorf("[ERROR] Error setting lifecycle_state: %s", err))
	}
	if vpnServer.LifecycleReasons != nil {
		if err := d.Set("lifecycle_reasons", resourceVPNServerFlattenLifecycleReasons(vpnServer.LifecycleReasons)); err != nil {
			return diag.FromErr(fmt.Errorf("[ERROR] Error setting lifecycle_reasons: %s", err))
		}
	}
	privateIps := []map[string]interface{}{}
	for _, privateIpsItem := range vpnServer.PrivateIps {
		privateIpsItemMap := resourceIBMIsVPNServerReservedIPReferenceToMap(privateIpsItem)
		privateIps = append(privateIps, privateIpsItemMap)
	}
	if err = d.Set("private_ips", privateIps); err != nil {
		return diag.FromErr(fmt.Errorf("[ERROR] Error setting private_ips: %s", err))
	}

	if vpnServer.VPC != nil {
		err = d.Set("vpc", dataSourceVPNServerFlattenVpcReference(vpnServer.VPC))
		if err != nil {
			return diag.FromErr(fmt.Errorf("[ERROR] Error setting the vpc: %s", err))
		}
	}

	if err = d.Set("resource_type", vpnServer.ResourceType); err != nil {
		return diag.FromErr(fmt.Errorf("[ERROR] Error setting resource_type: %s", err))
	}

	accesstags, err := flex.GetGlobalTagsUsingCRN(meta, *vpnServer.CRN, "", isVPNServerAccessTagType)
	if err != nil {
		log.Printf(
			"Error on get of resource vpn server (%s) access tags: %s", d.Id(), err)
	}
	d.Set(isVPNServerAccessTags, accesstags)

	return nil
}

func resourceIBMIsVPNServerReservedIPReferenceToMap(reservedIPReference vpcv1.ReservedIPReference) map[string]interface{} {
	reservedIPReferenceMap := map[string]interface{}{}

	reservedIPReferenceMap["address"] = reservedIPReference.Address
	if reservedIPReference.Deleted != nil {
		DeletedMap := resourceIBMIsVPNServerReservedIPReferenceDeletedToMap(*reservedIPReference.Deleted)
		reservedIPReferenceMap["deleted"] = []map[string]interface{}{DeletedMap}
	}
	reservedIPReferenceMap["href"] = reservedIPReference.Href
	reservedIPReferenceMap["id"] = reservedIPReference.ID
	reservedIPReferenceMap["name"] = reservedIPReference.Name
	reservedIPReferenceMap["resource_type"] = reservedIPReference.ResourceType

	return reservedIPReferenceMap
}

func resourceIBMIsVPNServerReservedIPReferenceDeletedToMap(reservedIPReferenceDeleted vpcv1.Deleted) map[string]interface{} {
	reservedIPReferenceDeletedMap := map[string]interface{}{}

	reservedIPReferenceDeletedMap["more_info"] = reservedIPReferenceDeleted.MoreInfo

	return reservedIPReferenceDeletedMap
}

func resourceIBMIsVPNServerVPCReference(result *vpcv1.VPCReference) (vpcs []map[string]interface{}) {
	vpcs = append(vpcs, resourceIBMIsVPNServerVPCReferenceToMap(*result))
	return vpcs
}

func resourceIBMIsVPNServerVPCReferenceToMap(vpcRef vpcv1.VPCReference) map[string]interface{} {
	vpcMap := map[string]interface{}{}
	vpcMap["crn"] = vpcRef.CRN
	if vpcRef.Deleted != nil {
		deletedMap := resourceIBMIsVPNServerVPCReferenceDeletedToMap(*vpcRef.Deleted)
		// vpcMap["deleted"] = []map[string]interface{}{deletedMap}
		vpcMap["deleted"] = deletedMap
	}
	vpcMap["href"] = vpcRef.Href
	vpcMap["name"] = vpcRef.Name
	return vpcMap
}

func resourceIBMIsVPNServerVPCReferenceDeletedToMap(vpcRefDeleted vpcv1.Deleted) map[string]interface{} {
	vpcRefDeletedMap := map[string]interface{}{}
	vpcRefDeletedMap["more_info"] = vpcRefDeleted.MoreInfo
	return vpcRefDeletedMap
}

func resourceIBMIsVPNServerUpdate(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	sess, err := vpcClient(meta)
	if err != nil {
		return diag.FromErr(err)
	}

	updateVPNServerOptions := &vpcv1.UpdateVPNServerOptions{}
	updateVPNServerOptions.SetID(d.Id())
	hasChange := false

	patchVals := &vpcv1.VPNServerPatch{}
	if d.HasChange("certificate_crn") {
		crn_val := d.Get("certificate_crn").(string)

		certificateInstanceIdentity := &vpcv1.CertificateInstanceIdentity{}
		certificateInstanceIdentity.CRN = &crn_val
		patchVals.Certificate = certificateInstanceIdentity
		hasChange = true
	}

	if d.HasChange("client_authentication") {
		var clientAuthentication []vpcv1.VPNServerAuthenticationPrototypeIntf
		clientAuthArray := d.Get("client_authentication").([]interface{})
		for _, clientauth := range clientAuthArray {
			clientAuth := clientauth.(map[string]interface{})
			method := clientAuth["method"].(string)
			clientAuthPrototype := &vpcv1.VPNServerAuthenticationPrototype{}
			clientAuthPrototype.Method = &method

			if method == "certificate" {
				if clientAuth["client_ca_crn"] != nil && clientAuth["client_ca_crn"] != "" {
					crn_val := clientAuth["client_ca_crn"].(string)
					certificateInstanceIdentity := &vpcv1.CertificateInstanceIdentity{}
					certificateInstanceIdentity.CRN = &crn_val
					clientAuthPrototype.ClientCa = certificateInstanceIdentity

				} else {
					return diag.FromErr(fmt.Errorf("[ERROR] Error method type `certificate` should be passed with `client_ca_crn`"))
				}
			} else {
				if clientAuth["identity_provider"] != nil && clientAuth["identity_provider"] != "" {
					providerType := clientAuth["identity_provider"].(string)
					clientAuthPrototype.IdentityProvider = &vpcv1.VPNServerAuthenticationByUsernameIDProvider{
						ProviderType: &providerType,
					}
				} else {
					return diag.FromErr(fmt.Errorf("[ERROR] Error method type `username` should be passed with `identity_provider`"))
				}

			}
			clientAuthentication = append(clientAuthentication, clientAuthPrototype)
		}
		patchVals.ClientAuthentication = clientAuthentication
		hasChange = true
	}

	if d.HasChange("client_ip_pool") {
		patchVals.ClientIPPool = core.StringPtr(d.Get("client_ip_pool").(string))
		hasChange = true
	}

	if d.HasChange("client_dns_server_ips") {
		var clientDnsServerIps []vpcv1.IP
		clientDnsServerIpsArray := d.Get("client_dns_server_ips").(*schema.Set)
		for _, val := range clientDnsServerIpsArray.List() {
			value := val.(string)
			address := &vpcv1.IP{
				Address: &value,
			}
			clientDnsServerIps = append(clientDnsServerIps, *address)
		}
		patchVals.ClientDnsServerIps = clientDnsServerIps
		hasChange = true
	}

	if d.HasChange("client_idle_timeout") {
		patchVals.ClientIdleTimeout = core.Int64Ptr(int64(d.Get("client_idle_timeout").(int)))
		hasChange = true
	}

	if d.HasChange("enable_split_tunneling") {
		patchVals.EnableSplitTunneling = core.BoolPtr(d.Get("enable_split_tunneling").(bool))
		hasChange = true
	}

	if d.HasChange("name") {
		patchVals.Name = core.StringPtr(d.Get("name").(string))
		hasChange = true
	}

	if d.HasChange("port") {
		patchVals.Port = core.Int64Ptr(int64(d.Get("port").(int)))
		hasChange = true
	}

	if d.HasChange("protocol") {
		patchVals.Protocol = core.StringPtr(d.Get("protocol").(string))
		hasChange = true
	}

	getVPNServerOptions := &vpcv1.GetVPNServerOptions{}
	getVPNServerOptions.SetID(d.Id())
	vpnServer, response, err := sess.GetVPNServerWithContext(context, getVPNServerOptions)
	if err != nil {
		if response != nil && response.StatusCode == 404 {
			d.SetId("")
			return nil
		}
		log.Printf("[DEBUG] GetVPNServerWithContext failed %s\n%s", err, response)
		return diag.FromErr(fmt.Errorf("[ERROR] GetVPNServerWithContext failed %s\n%s", err, response))
	}
	eTag := response.Headers.Get("ETag") // Getting Etag from the response headers.

	if d.HasChange(isVPNServerAccessTags) {
		oldList, newList := d.GetChange(isVPNServerAccessTags)
		err = flex.UpdateGlobalTagsUsingCRN(oldList, newList, meta, *vpnServer.CRN, "", isVPNServerAccessTagType)
		if err != nil {
			log.Printf(
				"Error on update of resource vpn server (%s) access tags: %s", d.Id(), err)
		}
	}
	// Upgrade or Downgrade of Subnet
	if d.HasChange("subnets") {
		var subnets []vpcv1.SubnetIdentityIntf
		subnetsArray := d.Get("subnets").(*schema.Set)
		for _, val := range subnetsArray.List() {
			value := val.(string)
			subnetIdentity := &vpcv1.SubnetIdentity{
				ID: &value,
			}
			subnets = append(subnets, subnetIdentity)
		}
		patchVals.Subnets = subnets
		hasChange = true
	}
	if hasChange {
		updateVPNServerOptions.IfMatch = &eTag // if-Match or Etag Change for Patch
		updateVPNServerOptions.VPNServerPatch, _ = patchVals.AsPatch()
		_, response, err := sess.UpdateVPNServerWithContext(context, updateVPNServerOptions)
		if err != nil {
			log.Printf("[DEBUG] UpdateVPNServerWithContext failed %s\n%s", err, response)
			return diag.FromErr(fmt.Errorf("[ERROR] UpdateVPNServerWithContext failed %s\n%s", err, response))
		}
		_, err = isWaitForVPNServerStable(context, sess, d, d.Timeout(schema.TimeoutUpdate))
		if err != nil {
			return diag.FromErr(fmt.Errorf("[ERROR] VPNServer failed %s\n", err))
		}
	}

	return resourceIBMIsVPNServerRead(context, d, meta)
}

func resourceIBMIsVPNServerDelete(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	sess, err := vpcClient(meta)
	if err != nil {
		return diag.FromErr(err)
	}

	getVPNServerOptions := &vpcv1.GetVPNServerOptions{}
	getVPNServerOptions.SetID(d.Id())
	_, response, err := sess.GetVPNServerWithContext(context, getVPNServerOptions)
	if err != nil {
		if response != nil && response.StatusCode == 404 {
			d.SetId("")
			return nil
		}
		log.Printf("[DEBUG] GetVPNServerWithContext failed %s\n%s", err, response)
		return diag.FromErr(fmt.Errorf("[ERROR] GetVPNServerWithContext failed %s\n%s", err, response))
	}
	etag := response.Headers.Get("Etag")
	deleteVPNServerOptions := &vpcv1.DeleteVPNServerOptions{}
	deleteVPNServerOptions.SetID(d.Id())
	deleteVPNServerOptions.SetIfMatch(etag)

	response, err = sess.DeleteVPNServerWithContext(context, deleteVPNServerOptions)
	if err != nil {
		log.Printf("[DEBUG] DeleteVPNServerWithContext failed %s\n%s", err, response)
		return diag.FromErr(fmt.Errorf("[ERROR] DeleteVPNServerWithContext failed %s\n%s", err, response))
	}

	_, err = isWaitForVPNServerDeleted(context, sess, d)
	if err != nil {
		return diag.FromErr(fmt.Errorf("[ERROR] VPNServer failed %s\n", err))
	}
	d.SetId("")

	return nil
}

func isWaitForVPNServerDeleted(context context.Context, sess *vpcv1.VpcV1, d *schema.ResourceData) (interface{}, error) {

	log.Printf("Waiting for VPN Server (%s) to be deleted.", d.Id())
	stateConf := &resource.StateChangeConf{
		Pending: []string{"retry", isVPNServerStatusDeleting},
		Target:  []string{isVPNServerStatusDeleted, isVPNServerStatusFailed},
		Refresh: func() (interface{}, string, error) {
			getVPNServerOptions := &vpcv1.GetVPNServerOptions{}
			getVPNServerOptions.SetID(d.Id())

			vpnServer, response, err := sess.GetVPNServerWithContext(context, getVPNServerOptions)
			if err != nil {
				if response != nil && response.StatusCode == 404 {
					return vpnServer, isVPNServerStatusDeleted, nil
				}
				return vpnServer, *vpnServer.LifecycleState, fmt.Errorf("The VPC vpn server %s failed to delete: %s\n%s", d.Id(), err, response)
			}
			return vpnServer, *vpnServer.LifecycleState, nil

		},
		Timeout:    d.Timeout(schema.TimeoutDelete),
		Delay:      10 * time.Second,
		MinTimeout: 10 * time.Second,
	}

	return stateConf.WaitForState()
}

func resourceVPNServerFlattenLifecycleReasons(lifecycleReasons []vpcv1.VPNServerLifecycleReason) (lifecycleReasonsList []map[string]interface{}) {
	lifecycleReasonsList = make([]map[string]interface{}, 0)
	for _, lr := range lifecycleReasons {
		currentLR := map[string]interface{}{}
		if lr.Code != nil && lr.Message != nil {
			currentLR[isInstanceLifecycleReasonsCode] = *lr.Code
			currentLR[isInstanceLifecycleReasonsMessage] = *lr.Message
			if lr.MoreInfo != nil {
				currentLR[isInstanceLifecycleReasonsMoreInfo] = *lr.MoreInfo
			}
			lifecycleReasonsList = append(lifecycleReasonsList, currentLR)
		}
	}
	return lifecycleReasonsList
}

func resourceVPNServerFlattenHealthReasons(healthReasons []vpcv1.VPNServerHealthReason) (healthReasonsList []map[string]interface{}) {
	healthReasonsList = make([]map[string]interface{}, 0)
	for _, hr := range healthReasons {
		currentHR := map[string]interface{}{}
		if hr.Code != nil && hr.Message != nil {
			currentHR[isVolumeHealthReasonsCode] = *hr.Code
			currentHR[isVolumeHealthReasonsMessage] = *hr.Message
			if hr.MoreInfo != nil {
				currentHR[isVolumeHealthReasonsMoreInfo] = *hr.MoreInfo
			}
			healthReasonsList = append(healthReasonsList, currentHR)
		}
	}
	return healthReasonsList
}
