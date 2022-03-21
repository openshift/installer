package nutanix

import (
	"log"
	"strconv"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/terraform-providers/terraform-provider-nutanix/utils"
)

func dataSourceNutanixClusters() *schema.Resource {
	return &schema.Resource{
		Read:          dataSourceNutanixClustersRead,
		SchemaVersion: 1,
		StateUpgraders: []schema.StateUpgrader{
			{
				Type:    resourceNutanixDatasourceClustersResourceV0().CoreConfigSchema().ImpliedType(),
				Upgrade: resourceDatasourceClustersStateUpgradeV0,
				Version: 0,
			},
		},
		Schema: map[string]*schema.Schema{
			"api_version": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"entities": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
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
							Type:     schema.TypeString,
							Computed: true,
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
				},
			},
		},
	}
}

func dataSourceNutanixClustersRead(d *schema.ResourceData, meta interface{}) error {
	// Get client connection
	conn := meta.(*Client).API

	var filter string
	resp, err := conn.V3.ListAllCluster(filter)

	if err != nil {
		return err
	}

	if err := d.Set("api_version", resp.APIVersion); err != nil {
		return err
	}

	entities := make([]map[string]interface{}, len(resp.Entities))

	for k, v := range resp.Entities {
		entity := make(map[string]interface{})
		m, c := setRSEntityMetadata(v.Metadata)

		entity["metadata"] = m
		entity["project_reference"] = flattenReferenceValues(v.Metadata.ProjectReference)
		entity["owner_reference"] = flattenReferenceValues(v.Metadata.OwnerReference)
		entity["categories"] = c
		entity["api_version"] = utils.StringValue(v.APIVersion)
		entity["name"] = utils.StringValue(v.Status.Name)
		entity["state"] = utils.StringValue(v.Status.State)

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

		entity["nodes"] = nodes

		config := v.Status.Resources.Config
		entity["gpu_driver_version"] = utils.StringValue(config.GpuDriverVersion)

		clientAuth := make(map[string]interface{})
		if config.ClientAuth != nil {
			clientAuth["status"] = utils.StringValue(config.ClientAuth.Status)
			clientAuth["ca_chain"] = utils.StringValue(config.ClientAuth.CaChain)
			clientAuth["name"] = utils.StringValue(config.ClientAuth.Name)
		}

		entity["client_auth"] = clientAuth

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

		entity["authorized_public_key_list"] = authPublicKey

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

		entity["software_map_ncc"] = ncc
		entity["software_map_nos"] = nos

		entity["encryption_status"] = utils.StringValue(config.EncryptionStatus)

		signingInfo := make(map[string]interface{})

		if config.SslKey != nil {
			entity["ssl_key_type"] = utils.StringValue(config.SslKey.KeyType)
			entity["ssl_key_name"] = utils.StringValue(config.SslKey.KeyName)

			if config.SslKey.SigningInfo != nil {
				signingInfo["city"] = utils.StringValue(config.SslKey.SigningInfo.City)
				signingInfo["common_name_suffix"] = utils.StringValue(config.SslKey.SigningInfo.CommonNameSuffix)
				signingInfo["state"] = utils.StringValue(config.SslKey.SigningInfo.State)
				signingInfo["country_code"] = utils.StringValue(config.SslKey.SigningInfo.CountryCode)
				signingInfo["common_name"] = utils.StringValue(config.SslKey.SigningInfo.CommonName)
				signingInfo["organization"] = utils.StringValue(config.SslKey.SigningInfo.Organization)
				signingInfo["email_address"] = utils.StringValue(config.SslKey.SigningInfo.EmailAddress)
			}

			entity["ssl_key_signing_info"] = signingInfo
			entity["ssl_key_expire_datetime"] = utils.StringValue(config.SslKey.ExpireDatetime)
		} else {
			entity["ssl_key_type"] = ""
			entity["ssl_key_name"] = ""
			entity["ssl_key_signing_info"] = signingInfo
			entity["ssl_key_expire_datetime"] = ""
		}

		entity["service_list"] = utils.StringValueSlice(config.ServiceList)
		entity["supported_information_verbosity"] = utils.StringValue(config.SupportedInformationVerbosity)

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

		entity["certification_signing_info"] = certSigning
		entity["operation_mode"] = utils.StringValue(config.OperationMode)

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

		entity["ca_certificate_list"] = caCert

		entity["enabled_feature_list"] = utils.StringValueSlice(config.EnabledFeatureList)
		entity["is_available"] = utils.BoolValue(config.IsAvailable)

		build := make(map[string]interface{})
		if config.Build != nil {
			build["commit_id"] = utils.StringValue(config.Build.CommitID)
			build["full_version"] = utils.StringValue(config.Build.FullVersion)
			build["commit_date"] = utils.StringValue(config.Build.CommitDate)
			build["version"] = utils.StringValue(config.Build.Version)
			build["short_commit_id"] = utils.StringValue(config.Build.ShortCommitID)
			build["build_type"] = utils.StringValue(config.Build.BuildType)
		}

		entity["build"] = build

		entity["timezone"] = utils.StringValue(config.Timezone)
		entity["cluster_arch"] = utils.StringValue(config.ClusterArch)

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
		entity["management_server_list"] = managementServer

		network := v.Status.Resources.Network
		entity["masquerading_port"] = utils.Int64Value(network.MasqueradingPort)
		entity["masquerading_ip"] = utils.StringValue(network.MasqueradingIP)
		entity["external_ip"] = utils.StringValue(network.ExternalIP)

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
		entity["http_proxy_list"] = httpProxy

		smtpServCreds := make(map[string]interface{})
		smtpServAddr := make(map[string]interface{})
		if network.SMTPServer != nil {
			entity["smtp_server_type"] = utils.StringValue(network.SMTPServer.Type)
			entity["smtp_server_email_address"] = utils.StringValue(network.SMTPServer.EmailAddress)

			if network.SMTPServer.Server != nil {
				entity["smtp_server_proxy_type_list"] = utils.StringValueSlice(network.SMTPServer.Server.ProxyTypeList)

				if network.SMTPServer.Server.Credentials != nil {
					smtpServCreds["username"] = utils.StringValue(network.SMTPServer.Server.Credentials.Username)
					smtpServCreds["password"] = utils.StringValue(network.SMTPServer.Server.Credentials.Password)
				}
				smtpServAddr["ip"] = utils.StringValue(network.SMTPServer.Server.Address.IP)
				smtpServAddr["fqdn"] = utils.StringValue(network.SMTPServer.Server.Address.FQDN)
				smtpServAddr["port"] = strconv.Itoa(int(utils.Int64Value(network.SMTPServer.Server.Address.Port)))
				smtpServAddr["ipv6"] = utils.StringValue(network.SMTPServer.Server.Address.IPV6)
			}
			entity["smtp_server_credentials"] = smtpServCreds
			entity["smtp_server_address"] = smtpServAddr
		} else {
			entity["smtp_server_type"] = ""
			entity["smtp_server_email_address"] = ""
			entity["smtp_server_credentials"] = smtpServCreds
			entity["smtp_server_proxy_type_list"] = make([]string, 0)
			entity["smtp_server_address"] = smtpServAddr
		}

		entity["ntp_server_ip_list"] = utils.StringValueSlice(network.NameServerIPList)
		entity["external_subnet"] = utils.StringValue(network.ExternalSubnet)
		entity["external_data_services_ip"] = utils.StringValue(network.ExternalDataServicesIP)
		entity["internal_subnet"] = utils.StringValue(network.InternalSubnet)

		domain := network.DomainServer
		domServCreds := make(map[string]interface{})

		if domain != nil {
			entity["domain_server_nameserver"] = utils.StringValue(domain.Nameserver)
			entity["domain_server_name"] = utils.StringValue(domain.Name)

			domServCreds["username"] = utils.StringValue(domain.DomainCredentials.Username)
			domServCreds["password"] = utils.StringValue(domain.DomainCredentials.Password)
			entity["domain_server_credentials"] = domServCreds
		} else {
			entity["domain_server_nameserver"] = ""
			entity["domain_server_name"] = ""
			entity["domain_server_credentials"] = domServCreds
		}

		entity["nfs_subnet_whitelist"] = utils.StringValueSlice(network.NFSSubnetWhitelist)
		entity["name_server_ip_list"] = utils.StringValueSlice(network.NameServerIPList)

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
		entity["http_proxy_whitelist"] = httpWhiteList

		analysis := make(map[string]interface{})
		if v.Status.Resources.Analysis != nil {
			analysis["bully_vm_num"] = utils.StringValue(v.Status.Resources.Analysis.VMEfficiencyMap.BullyVMNum)
			analysis["constrained_vm_num"] = utils.StringValue(v.Status.Resources.Analysis.VMEfficiencyMap.ConstrainedVMNum)
			analysis["dead_vm_num"] = utils.StringValue(v.Status.Resources.Analysis.VMEfficiencyMap.DeadVMNum)
			analysis["inefficient_vm_num"] = utils.StringValue(v.Status.Resources.Analysis.VMEfficiencyMap.InefficientVMNum)
			analysis["overprovisioned_vm_num"] = utils.StringValue(v.Status.Resources.Analysis.VMEfficiencyMap.OverprovisionedVMNum)
		}
		entity["analysis_vm_efficiency_map"] = analysis

		entities[k] = entity
	}

	if err := d.Set("entities", entities); err != nil {
		return err
	}

	d.SetId(resource.UniqueId())

	return nil
}

func resourceDatasourceClustersStateUpgradeV0(is map[string]interface{}, meta interface{}) (map[string]interface{}, error) {
	log.Printf("[DEBUG] Entering resourceDatasourceClustersStateUpgradeV0")
	return resourceNutanixCategoriesMigrateState(is, meta)
}

func resourceNutanixDatasourceClustersResourceV0() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"api_version": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"entities": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
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
							Type:     schema.TypeString,
							Computed: true,
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
				},
			},
		},
	}
}
