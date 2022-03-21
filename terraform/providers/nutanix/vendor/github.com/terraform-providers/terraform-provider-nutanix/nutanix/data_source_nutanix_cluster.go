package nutanix

import (
	"fmt"
	"log"
	"strconv"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	v3 "github.com/terraform-providers/terraform-provider-nutanix/client/v3"
	"github.com/terraform-providers/terraform-provider-nutanix/utils"
)

func dataSourceNutanixCluster() *schema.Resource {
	return &schema.Resource{
		Read:          dataSourceNutanixClusterRead,
		SchemaVersion: 1,
		StateUpgraders: []schema.StateUpgrader{
			{
				Type:    resourceNutanixDatasourceClusterResourceV0().CoreConfigSchema().ImpliedType(),
				Upgrade: resourceDatasourceClusterStateUpgradeV0,
				Version: 0,
			},
		},
		Schema: map[string]*schema.Schema{
			"cluster_id": {
				Type:          schema.TypeString,
				Computed:      true,
				Optional:      true,
				ConflictsWith: []string{"name"},
			},
			"metadata": {
				Type:     schema.TypeMap,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"last_update_time": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"kind": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"uuid": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"creation_time": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"spec_version": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"spec_hash": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"name": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
			"categories": categoriesSchema(),
			"project_reference": {
				Type:     schema.TypeMap,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"kind": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"uuid": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"name": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
			"owner_reference": {
				Type:     schema.TypeMap,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"kind": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"uuid": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"name": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
			"api_version": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"name": {
				Type:          schema.TypeString,
				Optional:      true,
				Computed:      true,
				ConflictsWith: []string{"cluster_id"},
			},

			// COMPUTED
			"state": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"nodes": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"ip": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"version": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"type": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
			"gpu_driver_version": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"client_auth": {
				Type:     schema.TypeMap,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"status": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"ca_chain": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"name": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
			"authorized_public_key_list": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"key": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"name": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
			"software_map_ncc": {
				Type:     schema.TypeMap,
				Computed: true,
			},
			"software_map_nos": {
				Type:     schema.TypeMap,
				Computed: true,
			},
			"encryption_status": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"ssl_key_type": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"ssl_key_name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"ssl_key_signing_info": {
				Type:     schema.TypeMap,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"city": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"common_name_suffix": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"state": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"country_code": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"common_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"organization": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"email_address": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
			"ssl_key_expire_datetime": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"service_list": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"supported_information_verbosity": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"certification_signing_info": {
				Type:     schema.TypeMap,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"city": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"common_name_suffix": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"state": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"country_code": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"common_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"organization": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"email_address": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
			"operation_mode": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"ca_certificate_list": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"ca_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"certificate": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
			"enabled_feature_list": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"is_available": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"build": {
				Type:     schema.TypeMap,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"commit_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"full_version": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"commit_date": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"version": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"short_commit_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"build_type": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
			"timezone": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"cluster_arch": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"management_server_list": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"ip": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"drs_enabled": {
							Type:     schema.TypeBool,
							Computed: true,
						},
						"status_list": {
							Type:     schema.TypeList,
							Computed: true,
							Elem:     &schema.Schema{Type: schema.TypeString},
						},
						"type": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
			"masquerading_port": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"masquerading_ip": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"external_ip": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"http_proxy_list": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"credentials": {
							Type:     schema.TypeMap,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"username": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"password": {
										Type:     schema.TypeString,
										Computed: true,
									},
								},
							},
						},
						"proxy_type_list": {
							Type:     schema.TypeList,
							Computed: true,
							Elem:     &schema.Schema{Type: schema.TypeString},
						},
						"address": {
							Type:     schema.TypeMap,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"ip": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"fqdn": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"port": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"ipv6": {
										Type:     schema.TypeString,
										Computed: true,
									},
								},
							},
						},
					},
				},
			},
			"smtp_server_type": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"smtp_server_email_address": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"smtp_server_credentials": {
				Type:     schema.TypeMap,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"username": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"password": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
			"smtp_server_proxy_type_list": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"smtp_server_address": {
				Type:     schema.TypeMap,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"ip": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"fqdn": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"port": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"ipv6": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
			"ntp_server_ip_list": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"external_subnet": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"external_data_services_ip": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"internal_subnet": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"domain_server_nameserver": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"domain_server_name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"domain_server_credentials": {
				Type:     schema.TypeMap,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"username": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"password": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
			"nfs_subnet_whitelist": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"name_server_ip_list": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"http_proxy_whitelist": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"target": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"target_type": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
			"analysis_vm_efficiency_map": {
				Type:     schema.TypeMap,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"bully_vm_num": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"constrained_vm_num": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"dead_vm_num": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"inefficient_vm_num": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"overprovisioned_vm_num": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func dataSourceNutanixClusterRead(d *schema.ResourceData, meta interface{}) error {
	// Get client connection
	conn := meta.(*Client).API

	c, ok := d.GetOk("cluster_id")
	var v *v3.ClusterIntentResponse
	var err error
	if ok {
		// Make request to the API
		v, err = conn.V3.GetCluster(c.(string))
		if err != nil {
			return err
		}
	} else {
		n, ok := d.GetOk("name")
		if !ok {
			return fmt.Errorf("please provide the cluster_id or name attribute")
		}
		v, err = findClusterByName(conn, n.(string))
		if err != nil {
			return err
		}
	}

	m, c := setRSEntityMetadata(v.Metadata)

	if err := d.Set("metadata", m); err != nil {
		return err
	}

	if err := d.Set("categories", c); err != nil {
		return err
	}

	if err := d.Set("api_version", utils.StringValue(v.APIVersion)); err != nil {
		return err
	}

	if err := d.Set("project_reference", flattenReferenceValues(v.Metadata.ProjectReference)); err != nil {
		return err
	}

	if err := d.Set("owner_reference", flattenReferenceValues(v.Metadata.OwnerReference)); err != nil {
		return err
	}

	if err := d.Set("name", utils.StringValue(v.Status.Name)); err != nil {
		return err
	}

	if err := d.Set("state", utils.StringValue(v.Status.State)); err != nil {
		return err
	}

	nodes := make([]map[string]interface{}, 0)

	if v.Status.Resources.Nodes != nil {
		if v.Status.Resources.Nodes.HypervisorServerList != nil {
			nodes = make([]map[string]interface{}, len(v.Status.Resources.Nodes.HypervisorServerList))

			for k, v := range v.Status.Resources.Nodes.HypervisorServerList {
				node := make(map[string]interface{})
				node["ip"] = utils.StringValue(v.IP)
				node["version"] = utils.StringValue(v.Version)
				node["type"] = utils.StringValue(v.Type)
				nodes[k] = node
			}
		}
	}

	if err := d.Set("nodes", nodes); err != nil {
		return err
	}

	config := v.Status.Resources.Config
	if err := d.Set("gpu_driver_version", utils.StringValue(config.GpuDriverVersion)); err != nil {
		return err
	}

	clientAuth := make(map[string]interface{})
	if config.ClientAuth != nil {
		clientAuth["status"] = utils.StringValue(config.ClientAuth.Status)
		clientAuth["ca_chain"] = utils.StringValue(config.ClientAuth.CaChain)
		clientAuth["name"] = utils.StringValue(config.ClientAuth.Name)
	}

	if err := d.Set("client_auth", clientAuth); err != nil {
		return err
	}

	authPublicKey := make([]map[string]interface{}, 0)
	if config.AuthorizedPublicKeyList != nil {
		authPublicKey = make([]map[string]interface{}, len(config.AuthorizedPublicKeyList))

		for k, v := range config.AuthorizedPublicKeyList {
			auth := make(map[string]interface{})
			auth["key"] = utils.StringValue(v.Key)
			auth["name"] = utils.StringValue(v.Name)
			authPublicKey[k] = auth
		}
	}

	if err := d.Set("authorized_public_key_list", authPublicKey); err != nil {
		return err
	}

	ncc := make(map[string]interface{})
	nos := make(map[string]interface{})

	if config.SoftwareMap != nil {
		if config.SoftwareMap.NCC != nil {
			ncc["software_type"] = utils.StringValue(config.SoftwareMap.NCC.SoftwareType)
			ncc["status"] = utils.StringValue(config.SoftwareMap.NCC.Status)
			ncc["version"] = utils.StringValue(config.SoftwareMap.NCC.Version)
		}

		if config.SoftwareMap.NOS != nil {
			nos["software_type"] = utils.StringValue(config.SoftwareMap.NOS.SoftwareType)
			nos["status"] = utils.StringValue(config.SoftwareMap.NOS.Status)
			nos["version"] = utils.StringValue(config.SoftwareMap.NOS.Version)
		}
	}

	if err := d.Set("software_map_ncc", ncc); err != nil {
		return err
	}

	if err := d.Set("software_map_nos", nos); err != nil {
		return err
	}

	if err := d.Set("encryption_status", utils.StringValue(config.EncryptionStatus)); err != nil {
		return err
	}

	signingInfo := make(map[string]interface{})

	if config.SslKey != nil {
		if err := d.Set("ssl_key_type", utils.StringValue(config.SslKey.KeyType)); err != nil {
			return err
		}

		if err := d.Set("ssl_key_name", utils.StringValue(config.SslKey.KeyName)); err != nil {
			return err
		}

		if config.SslKey.SigningInfo != nil {
			signingInfo["city"] = utils.StringValue(config.SslKey.SigningInfo.City)
			signingInfo["common_name_suffix"] = utils.StringValue(config.SslKey.SigningInfo.CommonNameSuffix)
			signingInfo["state"] = utils.StringValue(config.SslKey.SigningInfo.State)
			signingInfo["country_code"] = utils.StringValue(config.SslKey.SigningInfo.CountryCode)
			signingInfo["common_name"] = utils.StringValue(config.SslKey.SigningInfo.CommonName)
			signingInfo["organization"] = utils.StringValue(config.SslKey.SigningInfo.Organization)
			signingInfo["email_address"] = utils.StringValue(config.SslKey.SigningInfo.EmailAddress)
		}

		if err := d.Set("ssl_key_signing_info", signingInfo); err != nil {
			return err
		}

		if err := d.Set("ssl_key_expire_datetime", utils.StringValue(config.SslKey.ExpireDatetime)); err != nil {
			return err
		}
	} else {
		if err := d.Set("ssl_key_type", ""); err != nil {
			return err
		}

		if err := d.Set("ssl_key_name", ""); err != nil {
			return err
		}

		if err := d.Set("ssl_key_signing_info", signingInfo); err != nil {
			return err
		}

		if err := d.Set("ssl_key_expire_datetime", ""); err != nil {
			return err
		}
	}

	if err := d.Set("service_list", utils.StringValueSlice(config.ServiceList)); err != nil {
		return err
	}

	if err := d.Set("supported_information_verbosity", utils.StringValue(config.SupportedInformationVerbosity)); err != nil {
		return err
	}

	certSigning := make(map[string]interface{})
	if config.CertificationSigningInfo != nil {
		certSigning["city"] = utils.StringValue(config.CertificationSigningInfo.City)
		certSigning["common_name_suffix"] = utils.StringValue(config.CertificationSigningInfo.CommonNameSuffix)
		certSigning["state"] = utils.StringValue(config.CertificationSigningInfo.State)
		certSigning["country_code"] = utils.StringValue(config.CertificationSigningInfo.CountryCode)
		certSigning["common_name"] = utils.StringValue(config.CertificationSigningInfo.CommonName)
		certSigning["organization"] = utils.StringValue(config.CertificationSigningInfo.Organization)
		certSigning["email_address"] = utils.StringValue(config.CertificationSigningInfo.EmailAddress)
	}

	if err := d.Set("certification_signing_info", certSigning); err != nil {
		return err
	}

	if err := d.Set("operation_mode", utils.StringValue(config.OperationMode)); err != nil {
		return err
	}

	caCert := make([]map[string]interface{}, 0)

	if config.CaCertificateList != nil {
		caCert = make([]map[string]interface{}, len(config.CaCertificateList))

		for k, v := range config.CaCertificateList {
			ca := make(map[string]interface{})
			ca["ca_name"] = utils.StringValue(v.CaName)
			ca["certificate"] = utils.StringValue(v.Certificate)
			caCert[k] = ca
		}
	}

	if err := d.Set("ca_certificate_list", caCert); err != nil {
		return err
	}

	if err := d.Set("enabled_feature_list", utils.StringValueSlice(config.EnabledFeatureList)); err != nil {
		return err
	}

	if err := d.Set("is_available", utils.BoolValue(config.IsAvailable)); err != nil {
		return err
	}

	build := make(map[string]interface{})
	if config.Build != nil {
		build["commit_id"] = utils.StringValue(config.Build.CommitID)
		build["full_version"] = utils.StringValue(config.Build.FullVersion)
		build["commit_date"] = utils.StringValue(config.Build.CommitDate)
		build["version"] = utils.StringValue(config.Build.Version)
		build["short_commit_id"] = utils.StringValue(config.Build.ShortCommitID)
		build["build_type"] = utils.StringValue(config.Build.BuildType)
	}

	if err := d.Set("build", build); err != nil {
		return err
	}

	if err := d.Set("timezone", utils.StringValue(config.Timezone)); err != nil {
		return err
	}

	if err := d.Set("cluster_arch", utils.StringValue(config.ClusterArch)); err != nil {
		return err
	}

	managementServer := make([]map[string]interface{}, 0)

	if config.ManagementServerList != nil {
		managementServer = make([]map[string]interface{}, len(config.ManagementServerList))

		for k, v := range config.ManagementServerList {
			manage := make(map[string]interface{})
			manage["ip"] = utils.StringValue(v.IP)
			manage["drs_enabled"] = utils.BoolValue(v.DrsEnabled)
			manage["status_list"] = utils.StringValueSlice(v.StatusList)
			manage["type"] = utils.StringValue(v.Type)
			managementServer[k] = manage
		}
	}

	if err := d.Set("management_server_list", managementServer); err != nil {
		return err
	}

	network := v.Status.Resources.Network
	if err := d.Set("masquerading_port", utils.Int64Value(network.MasqueradingPort)); err != nil {
		return err
	}

	if err := d.Set("masquerading_ip", utils.StringValue(network.MasqueradingIP)); err != nil {
		return err
	}

	if err := d.Set("external_ip", utils.StringValue(network.ExternalIP)); err != nil {
		return err
	}

	httpProxy := make([]map[string]interface{}, 0)

	if network.HTTPProxyList != nil {
		httpProxy = make([]map[string]interface{}, len(network.HTTPProxyList))

		for k, v := range network.HTTPProxyList {
			http := make(map[string]interface{})
			creds := make(map[string]interface{})
			addr := make(map[string]interface{})

			if v.Credentials != nil {
				creds["username"] = utils.StringValue(v.Credentials.Username)
				creds["password"] = utils.StringValue(v.Credentials.Password)
				http["credentials"] = creds
			}

			http["proxy_type_list"] = utils.StringValueSlice(v.ProxyTypeList)
			addr["ip"] = utils.StringValue(v.Address.IP)
			addr["fqdn"] = utils.StringValue(v.Address.FQDN)
			addr["port"] = strconv.Itoa(int(utils.Int64Value(v.Address.Port)))
			addr["ipv6"] = utils.StringValue(v.Address.IPV6)
			http["address"] = addr

			httpProxy[k] = http
		}
	}

	if err := d.Set("http_proxy_list", httpProxy); err != nil {
		return err
	}

	smtpServCreds := make(map[string]interface{})
	smtpServAddr := make(map[string]interface{})

	if network.SMTPServer != nil {
		if err := d.Set("smtp_server_type", utils.StringValue(network.SMTPServer.Type)); err != nil {
			return err
		}

		if err := d.Set("smtp_server_email_address", utils.StringValue(network.SMTPServer.EmailAddress)); err != nil {
			return err
		}

		if network.SMTPServer.Server != nil {
			if err := d.Set("smtp_server_proxy_type_list", utils.StringValueSlice(network.SMTPServer.Server.ProxyTypeList)); err != nil {
				return err
			}

			if network.SMTPServer.Server.Credentials != nil {
				smtpServCreds["username"] = utils.StringValue(network.SMTPServer.Server.Credentials.Username)
				smtpServCreds["password"] = utils.StringValue(network.SMTPServer.Server.Credentials.Password)
			}

			smtpServAddr["ip"] = utils.StringValue(network.SMTPServer.Server.Address.IP)
			smtpServAddr["fqdn"] = utils.StringValue(network.SMTPServer.Server.Address.FQDN)
			smtpServAddr["port"] = strconv.Itoa(int(utils.Int64Value(network.SMTPServer.Server.Address.Port)))
			smtpServAddr["ipv6"] = utils.StringValue(network.SMTPServer.Server.Address.IPV6)
		}

		if err := d.Set("smtp_server_credentials", smtpServCreds); err != nil {
			return err
		}

		if err := d.Set("smtp_server_address", smtpServAddr); err != nil {
			return err
		}
	} else {
		if err := d.Set("smtp_server_type", ""); err != nil {
			return err
		}
		if err := d.Set("smtp_server_email_address", ""); err != nil {
			return err
		}
		if err := d.Set("smtp_server_credentials", smtpServCreds); err != nil {
			return err
		}
		if err := d.Set("smtp_server_proxy_type_list", make([]string, 0)); err != nil {
			return err
		}
		if err := d.Set("smtp_server_address", smtpServAddr); err != nil {
			return err
		}
	}

	if err := d.Set("ntp_server_ip_list", utils.StringValueSlice(network.NameServerIPList)); err != nil {
		return err
	}

	if err := d.Set("external_subnet", utils.StringValue(network.ExternalSubnet)); err != nil {
		return err
	}

	if err := d.Set("external_data_services_ip", utils.StringValue(network.ExternalDataServicesIP)); err != nil {
		return err
	}

	if err := d.Set("internal_subnet", utils.StringValue(network.InternalSubnet)); err != nil {
		return err
	}

	domain := network.DomainServer
	domServCreds := make(map[string]interface{})

	if domain != nil {
		if err := d.Set("domain_server_nameserver", utils.StringValue(domain.Nameserver)); err != nil {
			return err
		}

		if err := d.Set("domain_server_name", utils.StringValue(domain.Name)); err != nil {
			return err
		}

		domServCreds["username"] = utils.StringValue(domain.DomainCredentials.Username)
		domServCreds["password"] = utils.StringValue(domain.DomainCredentials.Password)

		if err := d.Set("domain_server_credentials", domServCreds); err != nil {
			return err
		}
	} else {
		if err := d.Set("domain_server_nameserver", ""); err != nil {
			return err
		}
		if err := d.Set("domain_server_name", ""); err != nil {
			return err
		}
		if err := d.Set("domain_server_credentials", domServCreds); err != nil {
			return err
		}
	}

	if err := d.Set("nfs_subnet_whitelist", utils.StringValueSlice(network.NFSSubnetWhitelist)); err != nil {
		return err
	}

	if err := d.Set("name_server_ip_list", utils.StringValueSlice(network.NameServerIPList)); err != nil {
		return err
	}

	httpWhiteList := make([]map[string]interface{}, 0)
	if network.HTTPProxyWhitelist != nil {
		httpWhiteList = make([]map[string]interface{}, len(network.HTTPProxyWhitelist))

		for k, v := range network.HTTPProxyWhitelist {
			http := make(map[string]interface{})
			http["target"] = utils.StringValue(v.Target)
			http["target_type"] = utils.StringValue(v.TargetType)
			httpWhiteList[k] = http
		}
	}

	if err := d.Set("http_proxy_whitelist", httpWhiteList); err != nil {
		return err
	}

	analysis := make(map[string]interface{})
	if v.Status.Resources.Analysis != nil {
		analysis["bully_vm_num"] = utils.StringValue(v.Status.Resources.Analysis.VMEfficiencyMap.BullyVMNum)
		analysis["constrained_vm_num"] = utils.StringValue(v.Status.Resources.Analysis.VMEfficiencyMap.ConstrainedVMNum)
		analysis["dead_vm_num"] = utils.StringValue(v.Status.Resources.Analysis.VMEfficiencyMap.DeadVMNum)
		analysis["inefficient_vm_num"] = utils.StringValue(v.Status.Resources.Analysis.VMEfficiencyMap.InefficientVMNum)
		analysis["overprovisioned_vm_num"] = utils.StringValue(v.Status.Resources.Analysis.VMEfficiencyMap.OverprovisionedVMNum)
	}

	if err := d.Set("analysis_vm_efficiency_map", analysis); err != nil {
		return err
	}

	cUUID := utils.StringValue(v.Metadata.UUID)
	if err := d.Set("cluster_id", cUUID); err != nil {
		return err
	}

	d.SetId(cUUID)
	return nil
}

func findClusterByName(conn *v3.Client, name string) (*v3.ClusterIntentResponse, error) {
	filter := fmt.Sprintf("name==%s", name)
	resp, err := conn.V3.ListAllCluster(filter)
	if err != nil {
		return nil, err
	}
	entities := resp.Entities

	found := make([]*v3.ClusterIntentResponse, 0)
	for _, v := range entities {
		if *v.Status.Name == name {
			found = append(found, &v3.ClusterIntentResponse{
				Status:     v.Status,
				Spec:       v.Spec,
				Metadata:   v.Metadata,
				APIVersion: v.APIVersion,
			})
		}
	}

	if len(found) > 1 {
		return nil, fmt.Errorf("your query returned more than one result. Please use cluster_id argument instead")
	}

	if len(found) == 0 {
		return nil, fmt.Errorf("did not find cluster with name %s", name)
	}

	return found[0], nil
}

func resourceDatasourceClusterStateUpgradeV0(is map[string]interface{}, meta interface{}) (map[string]interface{}, error) {
	log.Printf("[DEBUG] Entering resourceDatasourceClusterStateUpgradeV0")
	return resourceNutanixCategoriesMigrateState(is, meta)
}

func resourceNutanixDatasourceClusterResourceV0() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"cluster_id": {
				Type:          schema.TypeString,
				Computed:      true,
				Optional:      true,
				ConflictsWith: []string{"name"},
			},
			"metadata": {
				Type:     schema.TypeMap,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"last_update_time": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"kind": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"uuid": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"creation_time": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"spec_version": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"spec_hash": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"name": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
			"categories": {
				Type:     schema.TypeMap,
				Optional: true,
				Computed: true,
			},
			"project_reference": {
				Type:     schema.TypeMap,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"kind": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"uuid": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"name": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
			"owner_reference": {
				Type:     schema.TypeMap,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"kind": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"uuid": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"name": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
			"api_version": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"name": {
				Type:          schema.TypeString,
				Optional:      true,
				Computed:      true,
				ConflictsWith: []string{"cluster_id"},
			},

			// COMPUTED
			"state": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"nodes": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"ip": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"version": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"type": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
			"gpu_driver_version": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"client_auth": {
				Type:     schema.TypeMap,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"status": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"ca_chain": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"name": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
			"authorized_public_key_list": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"key": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"name": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
			"software_map_ncc": {
				Type:     schema.TypeMap,
				Computed: true,
			},
			"software_map_nos": {
				Type:     schema.TypeMap,
				Computed: true,
			},
			"encryption_status": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"ssl_key_type": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"ssl_key_name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"ssl_key_signing_info": {
				Type:     schema.TypeMap,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"city": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"common_name_suffix": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"state": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"country_code": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"common_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"organization": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"email_address": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
			"ssl_key_expire_datetime": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"service_list": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"supported_information_verbosity": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"certification_signing_info": {
				Type:     schema.TypeMap,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"city": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"common_name_suffix": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"state": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"country_code": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"common_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"organization": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"email_address": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
			"operation_mode": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"ca_certificate_list": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"ca_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"certificate": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
			"enabled_feature_list": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"is_available": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"build": {
				Type:     schema.TypeMap,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"commit_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"full_version": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"commit_date": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"version": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"short_commit_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"build_type": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
			"timezone": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"cluster_arch": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"management_server_list": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"ip": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"drs_enabled": {
							Type:     schema.TypeBool,
							Computed: true,
						},
						"status_list": {
							Type:     schema.TypeList,
							Computed: true,
							Elem:     &schema.Schema{Type: schema.TypeString},
						},
						"type": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
			"masquerading_port": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"masquerading_ip": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"external_ip": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"http_proxy_list": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"credentials": {
							Type:     schema.TypeMap,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"username": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"password": {
										Type:     schema.TypeString,
										Computed: true,
									},
								},
							},
						},
						"proxy_type_list": {
							Type:     schema.TypeList,
							Computed: true,
							Elem:     &schema.Schema{Type: schema.TypeString},
						},
						"address": {
							Type:     schema.TypeMap,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"ip": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"fqdn": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"port": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"ipv6": {
										Type:     schema.TypeString,
										Computed: true,
									},
								},
							},
						},
					},
				},
			},
			"smtp_server_type": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"smtp_server_email_address": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"smtp_server_credentials": {
				Type:     schema.TypeMap,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"username": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"password": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
			"smtp_server_proxy_type_list": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"smtp_server_address": {
				Type:     schema.TypeMap,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"ip": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"fqdn": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"port": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"ipv6": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
			"ntp_server_ip_list": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"external_subnet": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"external_data_services_ip": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"internal_subnet": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"domain_server_nameserver": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"domain_server_name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"domain_server_credentials": {
				Type:     schema.TypeMap,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"username": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"password": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
			"nfs_subnet_whitelist": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"name_server_ip_list": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"http_proxy_whitelist": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"target": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"target_type": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
			"analysis_vm_efficiency_map": {
				Type:     schema.TypeMap,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"bully_vm_num": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"constrained_vm_num": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"dead_vm_num": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"inefficient_vm_num": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"overprovisioned_vm_num": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
		},
	}
}
