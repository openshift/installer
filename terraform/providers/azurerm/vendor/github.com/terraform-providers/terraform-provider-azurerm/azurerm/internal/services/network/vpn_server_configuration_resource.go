package network

import (
	"fmt"
	"log"
	"time"

	"github.com/Azure/azure-sdk-for-go/services/network/mgmt/2020-05-01/network"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/tf"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/network/parse"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tags"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/timeouts"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func resourceVPNServerConfiguration() *schema.Resource {
	return &schema.Resource{
		Create: resourceVPNServerConfigurationCreateUpdate,
		Read:   resourceVPNServerConfigurationRead,
		Update: resourceVPNServerConfigurationCreateUpdate,
		Delete: resourceVPNServerConfigurationDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(90 * time.Minute),
			Read:   schema.DefaultTimeout(5 * time.Minute),
			Update: schema.DefaultTimeout(90 * time.Minute),
			Delete: schema.DefaultTimeout(90 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringIsNotEmpty,
			},

			"resource_group_name": azure.SchemaResourceGroupName(),

			"location": azure.SchemaLocation(),

			"vpn_authentication_types": {
				Type:     schema.TypeList,
				Required: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
					ValidateFunc: validation.StringInSlice([]string{
						string(network.AAD),
						string(network.Certificate),
						string(network.Radius),
					}, false),
				},

				// StatusCode=400 -- Original Error: Code="MultipleVpnAuthenticationTypesNotSupprtedOnVpnServerConfiguration"
				// Message="VpnServerConfiguration XXX/acctestrg-191125124621329676 supports single VpnAuthenticationType at a time.
				// Customer has specified 3 number of VpnAuthenticationTypes."
				MaxItems: 1,
			},

			// Optional
			"azure_active_directory_authentication": {
				Type:     schema.TypeList,
				Optional: true,
				MinItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"audience": {
							Type:     schema.TypeString,
							Required: true,
						},

						"issuer": {
							Type:     schema.TypeString,
							Required: true,
						},

						"tenant": {
							Type:     schema.TypeString,
							Required: true,
						},
					},
				},
			},

			"client_revoked_certificate": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:     schema.TypeString,
							Required: true,
						},

						"thumbprint": {
							Type:     schema.TypeString,
							Required: true,
						},
					},
				},
			},

			"client_root_certificate": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:     schema.TypeString,
							Required: true,
						},

						"public_cert_data": {
							Type:     schema.TypeString,
							Required: true,
						},
					},
				},
			},

			"ipsec_policy": {
				Type:     schema.TypeList,
				Optional: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"dh_group": {
							Type:     schema.TypeString,
							Required: true,
							ValidateFunc: validation.StringInSlice([]string{
								string(network.DHGroup1),
								string(network.DHGroup2),
								string(network.DHGroup14),
								string(network.DHGroup24),
								string(network.DHGroup2048),
								string(network.ECP256),
								string(network.ECP384),
								string(network.None),
							}, false),
						},

						"ike_encryption": {
							Type:     schema.TypeString,
							Required: true,
							ValidateFunc: validation.StringInSlice([]string{
								string(network.AES128),
								string(network.AES192),
								string(network.AES256),
								string(network.DES),
								string(network.DES3),
								string(network.GCMAES128),
								string(network.GCMAES256),
							}, false),
						},

						"ike_integrity": {
							Type:     schema.TypeString,
							Required: true,
							ValidateFunc: validation.StringInSlice([]string{
								string(network.IkeIntegrityGCMAES128),
								string(network.IkeIntegrityGCMAES256),
								string(network.IkeIntegrityMD5),
								string(network.IkeIntegritySHA1),
								string(network.IkeIntegritySHA256),
								string(network.IkeIntegritySHA384),
							}, false),
						},

						"ipsec_encryption": {
							Type:     schema.TypeString,
							Required: true,
							ValidateFunc: validation.StringInSlice([]string{
								string(network.IpsecEncryptionAES128),
								string(network.IpsecEncryptionAES192),
								string(network.IpsecEncryptionAES256),
								string(network.IpsecEncryptionDES),
								string(network.IpsecEncryptionDES3),
								string(network.IpsecEncryptionGCMAES128),
								string(network.IpsecEncryptionGCMAES192),
								string(network.IpsecEncryptionGCMAES256),
								string(network.IpsecEncryptionNone),
							}, false),
						},

						"ipsec_integrity": {
							Type:     schema.TypeString,
							Required: true,
							ValidateFunc: validation.StringInSlice([]string{
								string(network.IpsecIntegrityGCMAES128),
								string(network.IpsecIntegrityGCMAES192),
								string(network.IpsecIntegrityGCMAES256),
								string(network.IpsecIntegrityMD5),
								string(network.IpsecIntegritySHA1),
								string(network.IpsecIntegritySHA256),
							}, false),
						},

						"pfs_group": {
							Type:     schema.TypeString,
							Required: true,
							ValidateFunc: validation.StringInSlice([]string{
								string(network.PfsGroupECP256),
								string(network.PfsGroupECP384),
								string(network.PfsGroupNone),
								string(network.PfsGroupPFS1),
								string(network.PfsGroupPFS2),
								string(network.PfsGroupPFS14),
								string(network.PfsGroupPFS24),
								string(network.PfsGroupPFS2048),
								string(network.PfsGroupPFSMM),
							}, false),
						},

						"sa_lifetime_seconds": {
							Type:     schema.TypeInt,
							Required: true,
						},

						"sa_data_size_kilobytes": {
							Type:     schema.TypeInt,
							Required: true,
						},
					},
				},
			},

			"radius": {
				Type:     schema.TypeList,
				Optional: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"server": {
							Type:     schema.TypeList,
							Optional: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"address": {
										Type:         schema.TypeString,
										Required:     true,
										ValidateFunc: validation.StringIsNotEmpty,
									},

									"secret": {
										Type:         schema.TypeString,
										Required:     true,
										ValidateFunc: validation.StringIsNotEmpty,
										Sensitive:    true,
									},

									"score": {
										Type:         schema.TypeInt,
										Required:     true,
										ValidateFunc: validation.IntBetween(1, 30),
									},
								},
							},
						},

						"client_root_certificate": {
							Type:     schema.TypeSet,
							Optional: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"name": {
										Type:     schema.TypeString,
										Required: true,
									},

									"thumbprint": {
										Type:     schema.TypeString,
										Required: true,
									},
								},
							},
						},

						"server_root_certificate": {
							Type:     schema.TypeSet,
							Required: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"name": {
										Type:     schema.TypeString,
										Required: true,
									},

									"public_cert_data": {
										Type:     schema.TypeString,
										Required: true,
									},
								},
							},
						},
					},
				},
				ConflictsWith: []string{
					"radius_server",
				},
			},

			"radius_server": {
				Type:       schema.TypeList,
				Optional:   true,
				MaxItems:   1,
				Deprecated: "Deprecated in favour of `radius`",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"address": {
							Type:         schema.TypeString,
							Required:     true,
							ValidateFunc: validation.StringIsNotEmpty,
						},

						"secret": {
							Type:         schema.TypeString,
							Required:     true,
							ValidateFunc: validation.StringIsNotEmpty,
							Sensitive:    true,
						},

						"client_root_certificate": {
							Type:     schema.TypeSet,
							Optional: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"name": {
										Type:     schema.TypeString,
										Required: true,
									},

									"thumbprint": {
										Type:     schema.TypeString,
										Required: true,
									},
								},
							},
						},

						"server_root_certificate": {
							Type:     schema.TypeSet,
							Required: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"name": {
										Type:     schema.TypeString,
										Required: true,
									},

									"public_cert_data": {
										Type:     schema.TypeString,
										Required: true,
									},
								},
							},
						},
					},
				},
				ConflictsWith: []string{
					"radius",
				},
			},

			"vpn_protocols": {
				Type:     schema.TypeSet,
				Optional: true,
				Computed: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
					ValidateFunc: validation.StringInSlice([]string{
						string(network.VpnGatewayTunnelingProtocolIkeV2),
						string(network.VpnGatewayTunnelingProtocolOpenVPN),
					}, false),
				},
			},

			"tags": tags.Schema(),
		},
	}
}

func resourceVPNServerConfigurationCreateUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Network.VpnServerConfigurationsClient
	ctx, cancel := timeouts.ForCreateUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	name := d.Get("name").(string)
	resourceGroup := d.Get("resource_group_name").(string)

	if d.IsNewResource() {
		existing, err := client.Get(ctx, resourceGroup, name)
		if err != nil {
			if !utils.ResponseWasNotFound(existing.Response) {
				return fmt.Errorf("Error checking for presence of existing VPN Server Configuration %q (Resource Group %q): %+v", name, resourceGroup, err)
			}
		}

		if existing.ID != nil && *existing.ID != "" {
			return tf.ImportAsExistsError("azurerm_vpn_server_configuration", *existing.ID)
		}
	}

	aadAuthenticationRaw := d.Get("azure_active_directory_authentication").([]interface{})
	aadAuthentication := expandVpnServerConfigurationAADAuthentication(aadAuthenticationRaw)

	clientRevokedCertsRaw := d.Get("client_revoked_certificate").(*schema.Set).List()
	clientRevokedCerts := expandVpnServerConfigurationClientRevokedCertificates(clientRevokedCertsRaw)

	clientRootCertsRaw := d.Get("client_root_certificate").(*schema.Set).List()
	clientRootCerts := expandVpnServerConfigurationClientRootCertificates(clientRootCertsRaw)

	ipSecPoliciesRaw := d.Get("ipsec_policy").([]interface{})
	ipSecPolicies := expandVpnServerConfigurationIPSecPolicies(ipSecPoliciesRaw)

	radiusRaw := d.Get("radius").([]interface{})
	if len(radiusRaw) == 0 {
		radiusRaw = d.Get("radius_server").([]interface{})
	}
	radius := expandVpnServerConfigurationRadius(radiusRaw)

	vpnProtocolsRaw := d.Get("vpn_protocols").(*schema.Set).List()
	vpnProtocols := expandVpnServerConfigurationVPNProtocols(vpnProtocolsRaw)

	supportsAAD := false
	supportsCertificates := false
	supportsRadius := false

	vpnAuthenticationTypesRaw := d.Get("vpn_authentication_types").([]interface{})
	vpnAuthenticationTypes := make([]network.VpnAuthenticationType, 0)
	for _, v := range vpnAuthenticationTypesRaw {
		authType := network.VpnAuthenticationType(v.(string))

		switch authType {
		case network.AAD:
			supportsAAD = true

		case network.Certificate:
			supportsCertificates = true

		case network.Radius:
			supportsRadius = true

		default:
			return fmt.Errorf("Unsupported `vpn_authentication_type`: %q", authType)
		}

		vpnAuthenticationTypes = append(vpnAuthenticationTypes, authType)
	}

	props := network.VpnServerConfigurationProperties{
		AadAuthenticationParameters:  aadAuthentication,
		VpnAuthenticationTypes:       &vpnAuthenticationTypes,
		VpnClientRootCertificates:    clientRootCerts,
		VpnClientRevokedCertificates: clientRevokedCerts,
		VpnClientIpsecPolicies:       ipSecPolicies,
		VpnProtocols:                 vpnProtocols,
	}

	if supportsAAD && aadAuthentication == nil {
		return fmt.Errorf("`azure_active_directory_authentication` must be specified when `vpn_authentication_type` is set to `AAD`")
	}

	// parameter:VpnServerConfigVpnClientRootCertificates is not specified when VpnAuthenticationType as Certificate is selected.
	if supportsCertificates && len(clientRootCertsRaw) == 0 {
		return fmt.Errorf("`client_root_certificate` must be specified when `vpn_authentication_type` is set to `Certificate`")
	}

	if supportsRadius {
		if radius == nil {
			return fmt.Errorf("`radius` must be specified when `vpn_authentication_type` is set to `Radius`")
		}

		if radius.servers != nil && len(*radius.servers) != 0 {
			props.RadiusServers = radius.servers
		}

		props.RadiusServerAddress = utils.String(radius.address)
		props.RadiusServerSecret = utils.String(radius.secret)

		props.RadiusClientRootCertificates = radius.clientRootCertificates
		props.RadiusServerRootCertificates = radius.serverRootCertificates
	}

	location := azure.NormalizeLocation(d.Get("location").(string))
	t := d.Get("tags").(map[string]interface{})
	parameters := network.VpnServerConfiguration{
		Location:                         utils.String(location),
		VpnServerConfigurationProperties: &props,
		Tags:                             tags.Expand(t),
	}

	future, err := client.CreateOrUpdate(ctx, resourceGroup, name, parameters)
	if err != nil {
		return fmt.Errorf("Error creating VPN Server Configuration %q (Resource Group %q): %+v", name, resourceGroup, err)
	}

	if err := future.WaitForCompletionRef(ctx, client.Client); err != nil {
		return fmt.Errorf("Error waiting for creation of VPN Server Configuration %q (Resource Group %q): %+v", name, resourceGroup, err)
	}

	resp, err := client.Get(ctx, resourceGroup, name)
	if err != nil {
		return fmt.Errorf("Error retrieving VPN Server Configuration %q (Resource Group %q): %+v", name, resourceGroup, err)
	}

	d.SetId(*resp.ID)

	return resourceVPNServerConfigurationRead(d, meta)
}

func resourceVPNServerConfigurationRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Network.VpnServerConfigurationsClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.VpnServerConfigurationID(d.Id())
	if err != nil {
		return err
	}

	resp, err := client.Get(ctx, id.ResourceGroup, id.Name)
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			log.Printf("[DEBUG] VPN Server Configuration %q was not found in Resource Group %q - removing from state!", id.Name, id.ResourceGroup)
			d.SetId("")
			return nil
		}

		return fmt.Errorf("Error retrieving VPN Server Configuration %q (Resource Group %q): %+v", id.Name, id.ResourceGroup, err)
	}

	d.Set("name", id.Name)
	d.Set("resource_group_name", id.ResourceGroup)

	if location := resp.Location; location != nil {
		d.Set("location", azure.NormalizeLocation(*location))
	}

	if props := resp.VpnServerConfigurationProperties; props != nil {
		flattenedAADAuthentication := flattenVpnServerConfigurationAADAuthentication(props.AadAuthenticationParameters)
		if err := d.Set("azure_active_directory_authentication", flattenedAADAuthentication); err != nil {
			return fmt.Errorf("Error setting `azure_active_directory_authentication`: %+v", err)
		}

		flattenedClientRootCerts := flattenVpnServerConfigurationClientRootCertificates(props.VpnClientRootCertificates)
		if err := d.Set("client_root_certificate", flattenedClientRootCerts); err != nil {
			return fmt.Errorf("Error setting `client_root_certificate`: %+v", err)
		}

		flattenedClientRevokedCerts := flattenVpnServerConfigurationClientRevokedCertificates(props.VpnClientRevokedCertificates)
		if err := d.Set("client_revoked_certificate", flattenedClientRevokedCerts); err != nil {
			return fmt.Errorf("Error setting `client_revoked_certificate`: %+v", err)
		}

		flattenedIPSecPolicies := flattenVpnServerConfigurationIPSecPolicies(props.VpnClientIpsecPolicies)
		if err := d.Set("ipsec_policy", flattenedIPSecPolicies); err != nil {
			return fmt.Errorf("Error setting `ipsec_policy`: %+v", err)
		}

		flattenedRadius := flattenVpnServerConfigurationRadius(props)
		if len(flattenedRadius) > 0 {
			if flattenedRadius[0].(map[string]interface{})["server"] != nil {
				if err := d.Set("radius", flattenedRadius); err != nil {
					return fmt.Errorf("Error setting `radius`: %+v", err)
				}
			} else {
				if err := d.Set("radius_server", flattenedRadius); err != nil {
					return fmt.Errorf("Error setting `radius_server`: %+v", err)
				}
			}
		}

		vpnAuthenticationTypes := make([]interface{}, 0)
		if props.VpnAuthenticationTypes != nil {
			for _, v := range *props.VpnAuthenticationTypes {
				vpnAuthenticationTypes = append(vpnAuthenticationTypes, string(v))
			}
		}
		if err := d.Set("vpn_authentication_types", vpnAuthenticationTypes); err != nil {
			return fmt.Errorf("Error setting `vpn_authentication_types`: %+v", err)
		}

		flattenedVpnProtocols := flattenVpnServerConfigurationVPNProtocols(props.VpnProtocols)
		if err := d.Set("vpn_protocols", schema.NewSet(schema.HashString, flattenedVpnProtocols)); err != nil {
			return fmt.Errorf("Error setting `vpn_protocols`: %+v", err)
		}
	}

	return tags.FlattenAndSet(d, resp.Tags)
}

func resourceVPNServerConfigurationDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Network.VpnServerConfigurationsClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.VpnServerConfigurationID(d.Id())
	if err != nil {
		return err
	}

	future, err := client.Delete(ctx, id.ResourceGroup, id.Name)
	if err != nil {
		return fmt.Errorf("Error deleting VPN Server Configuration %q (Resource Group %q): %+v", id.Name, id.ResourceGroup, err)
	}

	if err = future.WaitForCompletionRef(ctx, client.Client); err != nil {
		return fmt.Errorf("Error waiting for deletion of VPN Server Configuration %q (Resource Group %q): %+v", id.Name, id.ResourceGroup, err)
	}

	return nil
}

func expandVpnServerConfigurationAADAuthentication(input []interface{}) *network.AadAuthenticationParameters {
	if len(input) == 0 {
		return nil
	}

	v := input[0].(map[string]interface{})
	return &network.AadAuthenticationParameters{
		AadAudience: utils.String(v["audience"].(string)),
		AadIssuer:   utils.String(v["issuer"].(string)),
		AadTenant:   utils.String(v["tenant"].(string)),
	}
}

func flattenVpnServerConfigurationAADAuthentication(input *network.AadAuthenticationParameters) []interface{} {
	if input == nil {
		return []interface{}{}
	}

	audience := ""
	if input.AadAudience != nil {
		audience = *input.AadAudience
	}

	issuer := ""
	if input.AadIssuer != nil {
		issuer = *input.AadIssuer
	}

	tenant := ""
	if input.AadTenant != nil {
		tenant = *input.AadTenant
	}

	return []interface{}{
		map[string]interface{}{
			"audience": audience,
			"issuer":   issuer,
			"tenant":   tenant,
		},
	}
}

func expandVpnServerConfigurationClientRootCertificates(input []interface{}) *[]network.VpnServerConfigVpnClientRootCertificate {
	clientRootCertificates := make([]network.VpnServerConfigVpnClientRootCertificate, 0)

	for _, v := range input {
		raw := v.(map[string]interface{})
		clientRootCertificates = append(clientRootCertificates, network.VpnServerConfigVpnClientRootCertificate{
			Name:           utils.String(raw["name"].(string)),
			PublicCertData: utils.String(raw["public_cert_data"].(string)),
		})
	}

	return &clientRootCertificates
}

func flattenVpnServerConfigurationClientRootCertificates(input *[]network.VpnServerConfigVpnClientRootCertificate) []interface{} {
	if input == nil {
		return []interface{}{}
	}

	output := make([]interface{}, 0)

	for _, v := range *input {
		name := ""
		if v.Name != nil {
			name = *v.Name
		}

		publicCertData := ""
		if v.PublicCertData != nil {
			publicCertData = *v.PublicCertData
		}

		output = append(output, map[string]interface{}{
			"name":             name,
			"public_cert_data": publicCertData,
		})
	}

	return output
}

func expandVpnServerConfigurationClientRevokedCertificates(input []interface{}) *[]network.VpnServerConfigVpnClientRevokedCertificate {
	clientRevokedCertificates := make([]network.VpnServerConfigVpnClientRevokedCertificate, 0)

	for _, v := range input {
		raw := v.(map[string]interface{})
		clientRevokedCertificates = append(clientRevokedCertificates, network.VpnServerConfigVpnClientRevokedCertificate{
			Name:       utils.String(raw["name"].(string)),
			Thumbprint: utils.String(raw["thumbprint"].(string)),
		})
	}

	return &clientRevokedCertificates
}

func flattenVpnServerConfigurationClientRevokedCertificates(input *[]network.VpnServerConfigVpnClientRevokedCertificate) []interface{} {
	if input == nil {
		return []interface{}{}
	}

	output := make([]interface{}, 0)
	for _, v := range *input {
		name := ""
		if v.Name != nil {
			name = *v.Name
		}

		thumbprint := ""
		if v.Thumbprint != nil {
			thumbprint = *v.Thumbprint
		}

		output = append(output, map[string]interface{}{
			"name":       name,
			"thumbprint": thumbprint,
		})
	}
	return output
}

func expandVpnServerConfigurationIPSecPolicies(input []interface{}) *[]network.IpsecPolicy {
	ipSecPolicies := make([]network.IpsecPolicy, 0)

	for _, raw := range input {
		v := raw.(map[string]interface{})
		ipSecPolicies = append(ipSecPolicies, network.IpsecPolicy{
			DhGroup:             network.DhGroup(v["dh_group"].(string)),
			IkeEncryption:       network.IkeEncryption(v["ike_encryption"].(string)),
			IkeIntegrity:        network.IkeIntegrity(v["ike_integrity"].(string)),
			IpsecEncryption:     network.IpsecEncryption(v["ipsec_encryption"].(string)),
			IpsecIntegrity:      network.IpsecIntegrity(v["ipsec_integrity"].(string)),
			PfsGroup:            network.PfsGroup(v["pfs_group"].(string)),
			SaLifeTimeSeconds:   utils.Int32(int32(v["sa_lifetime_seconds"].(int))),
			SaDataSizeKilobytes: utils.Int32(int32(v["sa_data_size_kilobytes"].(int))),
		})
	}

	return &ipSecPolicies
}

func flattenVpnServerConfigurationIPSecPolicies(input *[]network.IpsecPolicy) []interface{} {
	if input == nil {
		return []interface{}{}
	}

	output := make([]interface{}, 0)
	for _, v := range *input {
		saDataSizeKilobytes := 0
		if v.SaDataSizeKilobytes != nil {
			saDataSizeKilobytes = int(*v.SaDataSizeKilobytes)
		}

		saLifeTimeSeconds := 0
		if v.SaLifeTimeSeconds != nil {
			saLifeTimeSeconds = int(*v.SaLifeTimeSeconds)
		}

		output = append(output, map[string]interface{}{
			"dh_group":               string(v.DhGroup),
			"ipsec_encryption":       string(v.IpsecEncryption),
			"ipsec_integrity":        string(v.IpsecIntegrity),
			"ike_encryption":         string(v.IkeEncryption),
			"ike_integrity":          string(v.IkeIntegrity),
			"pfs_group":              string(v.PfsGroup),
			"sa_data_size_kilobytes": saDataSizeKilobytes,
			"sa_lifetime_seconds":    saLifeTimeSeconds,
		})
	}
	return output
}

type vpnServerConfigurationRadius struct {
	address                string
	secret                 string
	servers                *[]network.RadiusServer
	clientRootCertificates *[]network.VpnServerConfigRadiusClientRootCertificate
	serverRootCertificates *[]network.VpnServerConfigRadiusServerRootCertificate
}

func expandVpnServerConfigurationRadius(input []interface{}) *vpnServerConfigurationRadius {
	if len(input) == 0 {
		return nil
	}

	val := input[0].(map[string]interface{})

	clientRootCertificates := make([]network.VpnServerConfigRadiusClientRootCertificate, 0)
	clientRootCertsRaw := val["client_root_certificate"].(*schema.Set).List()
	for _, raw := range clientRootCertsRaw {
		v := raw.(map[string]interface{})
		clientRootCertificates = append(clientRootCertificates, network.VpnServerConfigRadiusClientRootCertificate{
			Name:       utils.String(v["name"].(string)),
			Thumbprint: utils.String(v["thumbprint"].(string)),
		})
	}

	serverRootCertificates := make([]network.VpnServerConfigRadiusServerRootCertificate, 0)
	serverRootCertsRaw := val["server_root_certificate"].(*schema.Set).List()
	for _, raw := range serverRootCertsRaw {
		v := raw.(map[string]interface{})
		serverRootCertificates = append(serverRootCertificates, network.VpnServerConfigRadiusServerRootCertificate{
			Name:           utils.String(v["name"].(string)),
			PublicCertData: utils.String(v["public_cert_data"].(string)),
		})
	}

	radiusServers := make([]network.RadiusServer, 0)
	address := ""
	secret := ""

	if val["server"] != nil {
		radiusServersRaw := val["server"].([]interface{})
		for _, raw := range radiusServersRaw {
			v := raw.(map[string]interface{})
			radiusServers = append(radiusServers, network.RadiusServer{
				RadiusServerAddress: utils.String(v["address"].(string)),
				RadiusServerSecret:  utils.String(v["secret"].(string)),
				RadiusServerScore:   utils.Int64(int64(v["score"].(int))),
			})
		}
	} else {
		address = val["address"].(string)
		secret = val["secret"].(string)
	}

	return &vpnServerConfigurationRadius{
		address:                address,
		secret:                 secret,
		servers:                &radiusServers,
		clientRootCertificates: &clientRootCertificates,
		serverRootCertificates: &serverRootCertificates,
	}
}

func flattenVpnServerConfigurationRadius(input *network.VpnServerConfigurationProperties) []interface{} {
	if input == nil || (input.RadiusServerAddress == nil && input.RadiusServers == nil) || input.RadiusServerRootCertificates == nil || len(*input.RadiusServerRootCertificates) == 0 {
		return []interface{}{}
	}

	clientRootCertificates := make([]interface{}, 0)
	if input.RadiusClientRootCertificates != nil {
		for _, v := range *input.RadiusClientRootCertificates {
			name := ""
			if v.Name != nil {
				name = *v.Name
			}

			thumbprint := ""
			if v.Thumbprint != nil {
				thumbprint = *v.Thumbprint
			}

			clientRootCertificates = append(clientRootCertificates, map[string]interface{}{
				"name":       name,
				"thumbprint": thumbprint,
			})
		}
	}

	serverRootCertificates := make([]interface{}, 0)
	if input.RadiusServerRootCertificates != nil {
		for _, v := range *input.RadiusServerRootCertificates {
			name := ""
			if v.Name != nil {
				name = *v.Name
			}

			publicCertData := ""
			if v.PublicCertData != nil {
				publicCertData = *v.PublicCertData
			}

			serverRootCertificates = append(serverRootCertificates, map[string]interface{}{
				"name":             name,
				"public_cert_data": publicCertData,
			})
		}
	}

	schema := map[string]interface{}{
		"client_root_certificate": clientRootCertificates,
		"server_root_certificate": serverRootCertificates,
	}

	if input.RadiusServerAddress != nil && *input.RadiusServerAddress != "" {
		schema["address"] = *input.RadiusServerAddress
	}

	if input.RadiusServerSecret != nil && *input.RadiusServerSecret != "" {
		schema["secret"] = *input.RadiusServerSecret
	}

	if input.RadiusServers != nil && len(*input.RadiusServers) > 0 {
		servers := make([]interface{}, 0)

		for _, v := range *input.RadiusServers {
			address := ""
			if v.RadiusServerAddress != nil {
				address = *v.RadiusServerAddress
			}

			secret := ""
			if v.RadiusServerSecret != nil {
				secret = *v.RadiusServerSecret
			}

			score := 0
			if v.RadiusServerScore != nil {
				score = int(*v.RadiusServerScore)
			}

			servers = append(servers, map[string]interface{}{
				"address": address,
				"secret":  secret,
				"score":   score,
			})
		}

		schema["server"] = servers
	}

	return []interface{}{
		schema,
	}
}

func expandVpnServerConfigurationVPNProtocols(input []interface{}) *[]network.VpnGatewayTunnelingProtocol {
	vpnProtocols := make([]network.VpnGatewayTunnelingProtocol, 0)

	for _, v := range input {
		vpnProtocols = append(vpnProtocols, network.VpnGatewayTunnelingProtocol(v.(string)))
	}

	return &vpnProtocols
}

func flattenVpnServerConfigurationVPNProtocols(input *[]network.VpnGatewayTunnelingProtocol) []interface{} {
	if input == nil {
		return []interface{}{}
	}

	output := make([]interface{}, 0)

	for _, v := range *input {
		output = append(output, string(v))
	}

	return output
}
