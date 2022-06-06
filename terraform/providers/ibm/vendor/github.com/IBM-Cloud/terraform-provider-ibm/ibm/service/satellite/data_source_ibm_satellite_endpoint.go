// Copyright IBM Corp. 2021 All Rights Reserved.
// Licensed under the Mozilla Public License v2.0

package satellite

import (
	"context"
	"fmt"
	"log"

	"github.com/IBM-Cloud/container-services-go-sdk/satellitelinkv1"
	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/conns"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func DataSourceIBMSatelliteEndpoint() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceIbmSatelliteEndpointRead,

		Schema: map[string]*schema.Schema{
			"location": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The Location ID.",
			},
			"endpoint_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The Endpoint ID.",
			},
			"connection_type": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The type of the endpoint.",
			},
			"display_name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The display name of the endpoint. Endpoint names must start with a letter and end with an alphanumeric character, can contain letters, numbers, and hyphen (-), and must be 63 characters or fewer.",
			},
			"server_host": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The host name or IP address of the server endpoint. For 'http-tunnel' protocol, server_host can start with '*.' , which means a wildcard to it's sub domains. Such as '*.example.com' can accept request to 'api.example.com' and 'www.example.com'.",
			},
			"server_port": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "The port number of the server endpoint. For 'http-tunnel' protocol, server_port can be 0, which means any port. Such as 0 is good for 80 (http) and 443 (https).",
			},
			"sni": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The server name indicator (SNI) which used to connect to the server endpoint. Only useful if server side requires SNI.",
			},
			"client_protocol": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The protocol in the client application side.",
			},
			"client_mutual_auth": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "Whether enable mutual auth in the client application side, when client_protocol is 'tls' or 'https', this field is required.",
			},
			"server_protocol": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The protocol in the server application side. This parameter will change to default value if it is omitted even when using PATCH API. If client_protocol is 'udp', server_protocol must be 'udp'. If client_protocol is 'tcp'/'http', server_protocol could be 'tcp'/'tls' and default to 'tcp'. If client_protocol is 'tls'/'https', server_protocol could be 'tcp'/'tls' and default to 'tls'. If client_protocol is 'http-tunnel', server_protocol must be 'tcp'.",
			},
			"server_mutual_auth": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "Whether enable mutual auth in the server application side, when client_protocol is 'tls', this field is required.",
			},
			"reject_unauth": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "Whether reject any connection to the server application which is not authorized with the list of supplied CAs in the fields certs.server_cert.",
			},
			"timeout": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "The inactivity timeout in the Endpoint side.",
			},
			"created_by": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The service or person who created the endpoint. Must be 1000 characters or fewer.",
			},
			"sources": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"source_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The Source ID.",
						},
						"enabled": {
							Type:        schema.TypeBool,
							Computed:    true,
							Description: "Whether the source is enabled for the endpoint.",
						},
						"last_change": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The last time modify the Endpoint configurations.",
						},
						"pending": {
							Type:        schema.TypeBool,
							Computed:    true,
							Description: "Whether the source has been enabled on this endpoint.",
						},
					},
				},
			},
			"connector_port": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "The connector port.",
			},
			"crn": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Service instance associated with this location.",
			},
			"service_name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The service name of the endpoint.",
			},
			"client_host": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The hostname which Satellite Link server listen on for the on-location endpoint, or the hostname which the connector server listen on for the on-cloud endpoint destiantion.",
			},
			"client_port": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "The port which Satellite Link server listen on for the on-location, or the port which the connector server listen on for the on-cloud endpoint destiantion.",
			},
			"certs": {
				Type: schema.TypeList,
				//MaxItems:    1,
				Computed:    true,
				Description: "The certs. Once it is generated, this field will always be defined even it is unused until the cert/key is deleted.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"client": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "The CA which Satellite Link trust when receiving the connection from the client application.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"cert": {
										Type:        schema.TypeList,
										Computed:    true,
										Description: "The root cert or the self-signed cert of the client application.",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"filename": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "The filename of the cert.",
												},
											},
										},
									},
								},
							},
						},
						"server": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "The CA which Satellite Link trust when sending the connection to server application.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"cert": {
										Type:        schema.TypeList,
										Computed:    true,
										Description: "The root cert or the self-signed cert of the server application.",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"filename": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "The filename of the cert.",
												},
											},
										},
									},
								},
							},
						},
						"connector": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "The cert which Satellite Link connector provide to identify itself for connecting to the client/server application.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"cert": {
										Type:        schema.TypeList,
										Computed:    true,
										Description: "The end-entity cert of the connector.",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"filename": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "The filename of the cert.",
												},
											},
										},
									},
									"key": {
										Type:        schema.TypeList,
										Computed:    true,
										Description: "The private key of the connector.",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"filename": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "The name of the key.",
												},
											},
										},
									},
								},
							},
						},
					},
				},
			},
			"status": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Whether the Endpoint is active or not.",
			},
			"created_at": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The time when the Endpoint is created.",
			},
			"last_change": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The last time modify the Endpoint configurations.",
			},
			"performance": {
				Type: schema.TypeList,
				//MaxItems:    1,
				Computed:    true,
				Description: "The last performance data of the endpoint.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"connection": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Concurrent connections number of moment when probe read the data.",
						},
						"rx_bandwidth": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Average Receive (to Cloud) Bandwidth of last two minutes, unit is Byte/s.",
						},
						"tx_bandwidth": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Average Transmitted (to Location) Bandwidth of last two minutes, unit is Byte/s.",
						},
						"bandwidth": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Average Tatal Bandwidth of last two minutes, unit is Byte/s.",
						},
						"connectors": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "The last performance data of the endpoint from each Connector.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"connector": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The name of the connector reported the performance data.",
									},
									"connections": {
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "Concurrent connections number of moment when probe read the data from the Connector.",
									},
									"rx_bw": {
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "Average Transmitted (to Location) Bandwidth of last two minutes read from the Connector, unit is Byte/s.",
									},
									"tx_bw": {
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "Average Transmitted (to Location) Bandwidth of last two minutes read from the Connector, unit is Byte/s.",
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

func dataSourceIbmSatelliteEndpointRead(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	satelliteLinkClient, err := meta.(conns.ClientSession).SatellitLinkClientSession()
	if err != nil {
		return diag.FromErr(err)
	}

	getEndpointsOptions := &satellitelinkv1.GetEndpointsOptions{}

	getEndpointsOptions.SetLocationID(d.Get("location").(string))
	getEndpointsOptions.SetEndpointID(d.Get("endpoint_id").(string))

	endpoint, response, err := satelliteLinkClient.GetEndpointsWithContext(context, getEndpointsOptions)
	if err != nil {
		log.Printf("[DEBUG] GetEndpointsWithContext failed %s\n%s", err, response)
		return diag.FromErr(fmt.Errorf("GetEndpointsWithContext failed %s\n%s", err, response))
	}

	d.SetId(fmt.Sprintf("%s/%s", *endpoint.LocationID, *endpoint.EndpointID))
	if err = d.Set("connection_type", endpoint.ConnType); err != nil {
		return diag.FromErr(fmt.Errorf("[ERROR] Error setting connection_type: %s", err))
	}
	if err = d.Set("display_name", endpoint.DisplayName); err != nil {
		return diag.FromErr(fmt.Errorf("[ERROR] Error setting display_name: %s", err))
	}
	if err = d.Set("server_host", endpoint.ServerHost); err != nil {
		return diag.FromErr(fmt.Errorf("[ERROR] Error setting server_host: %s", err))
	}
	if err = d.Set("server_port", endpoint.ServerPort); err != nil {
		return diag.FromErr(fmt.Errorf("[ERROR] Error setting server_port: %s", err))
	}
	if err = d.Set("sni", endpoint.Sni); err != nil {
		return diag.FromErr(fmt.Errorf("[ERROR] Error setting sni: %s", err))
	}
	if err = d.Set("client_protocol", endpoint.ClientProtocol); err != nil {
		return diag.FromErr(fmt.Errorf("[ERROR] Error setting client_protocol: %s", err))
	}
	if err = d.Set("client_mutual_auth", endpoint.ClientMutualAuth); err != nil {
		return diag.FromErr(fmt.Errorf("[ERROR] Error setting client_mutual_auth: %s", err))
	}
	if err = d.Set("server_protocol", endpoint.ServerProtocol); err != nil {
		return diag.FromErr(fmt.Errorf("[ERROR] Error setting server_protocol: %s", err))
	}
	if err = d.Set("server_mutual_auth", endpoint.ServerMutualAuth); err != nil {
		return diag.FromErr(fmt.Errorf("[ERROR] Error setting server_mutual_auth: %s", err))
	}
	if err = d.Set("reject_unauth", endpoint.RejectUnauth); err != nil {
		return diag.FromErr(fmt.Errorf("[ERROR] Error setting reject_unauth: %s", err))
	}
	if err = d.Set("timeout", endpoint.Timeout); err != nil {
		return diag.FromErr(fmt.Errorf("[ERROR] Error setting timeout: %s", err))
	}
	if err = d.Set("created_by", endpoint.CreatedBy); err != nil {
		return diag.FromErr(fmt.Errorf("[ERROR] Error setting created_by: %s", err))
	}

	if endpoint.Sources != nil {
		err = d.Set("sources", dataSourceEndpointFlattenSources(endpoint.Sources))
		if err != nil {
			return diag.FromErr(fmt.Errorf("[ERROR] Error setting sources %s", err))
		}
	}
	if err = d.Set("connector_port", endpoint.ConnectorPort); err != nil {
		return diag.FromErr(fmt.Errorf("[ERROR] Error setting connector_port: %s", err))
	}
	if err = d.Set("crn", endpoint.Crn); err != nil {
		return diag.FromErr(fmt.Errorf("[ERROR] Error setting crn: %s", err))
	}
	if err = d.Set("service_name", endpoint.ServiceName); err != nil {
		return diag.FromErr(fmt.Errorf("[ERROR] Error setting service_name: %s", err))
	}
	if err = d.Set("client_host", endpoint.ClientHost); err != nil {
		return diag.FromErr(fmt.Errorf("[ERROR] Error setting client_host: %s", err))
	}
	if err = d.Set("client_port", endpoint.ClientPort); err != nil {
		return diag.FromErr(fmt.Errorf("[ERROR] Error setting client_port: %s", err))
	}

	if endpoint.Certs != nil {
		err = d.Set("certs", dataSourceEndpointFlattenCerts(*endpoint.Certs))
		if err != nil {
			return diag.FromErr(fmt.Errorf("[ERROR] Error setting certs %s", err))
		}
	}
	if err = d.Set("status", endpoint.Status); err != nil {
		return diag.FromErr(fmt.Errorf("[ERROR] Error setting status: %s", err))
	}
	if err = d.Set("created_at", endpoint.CreatedAt); err != nil {
		return diag.FromErr(fmt.Errorf("[ERROR] Error setting created_at: %s", err))
	}
	if err = d.Set("last_change", endpoint.LastChange); err != nil {
		return diag.FromErr(fmt.Errorf("[ERROR] Error setting last_change: %s", err))
	}

	if endpoint.Performance != nil {
		err = d.Set("performance", dataSourceEndpointFlattenPerformance(*endpoint.Performance))
		if err != nil {
			return diag.FromErr(fmt.Errorf("[ERROR] Error setting performance %s", err))
		}
	}

	return nil
}

func dataSourceEndpointFlattenSources(result []satellitelinkv1.SourceStatusObject) (sources []map[string]interface{}) {
	for _, sourcesItem := range result {
		sources = append(sources, dataSourceEndpointSourcesToMap(sourcesItem))
	}

	return sources
}

func dataSourceEndpointSourcesToMap(sourcesItem satellitelinkv1.SourceStatusObject) (sourcesMap map[string]interface{}) {
	sourcesMap = map[string]interface{}{}

	if sourcesItem.SourceID != nil {
		sourcesMap["source_id"] = sourcesItem.SourceID
	}
	if sourcesItem.Enabled != nil {
		sourcesMap["enabled"] = sourcesItem.Enabled
	}
	if sourcesItem.LastChange != nil {
		sourcesMap["last_change"] = sourcesItem.LastChange
	}
	if sourcesItem.Pending != nil {
		sourcesMap["pending"] = sourcesItem.Pending
	}

	return sourcesMap
}

func dataSourceEndpointFlattenCerts(result satellitelinkv1.EndpointCerts) (finalList []map[string]interface{}) {
	finalList = []map[string]interface{}{}
	finalMap := dataSourceEndpointCertsToMap(result)
	finalList = append(finalList, finalMap)

	return finalList
}

func dataSourceEndpointCertsToMap(certsItem satellitelinkv1.EndpointCerts) (certsMap map[string]interface{}) {
	certsMap = map[string]interface{}{}

	if certsItem.Client != nil {
		clientList := []map[string]interface{}{}
		clientMap := dataSourceEndpointCertsClientToMap(*certsItem.Client)
		clientList = append(clientList, clientMap)
		certsMap["client"] = clientList
	}
	if certsItem.Server != nil {
		serverList := []map[string]interface{}{}
		serverMap := dataSourceEndpointCertsServerToMap(*certsItem.Server)
		serverList = append(serverList, serverMap)
		certsMap["server"] = serverList
	}
	if certsItem.Connector != nil {
		connectorList := []map[string]interface{}{}
		connectorMap := dataSourceEndpointCertsConnectorToMap(*certsItem.Connector)
		connectorList = append(connectorList, connectorMap)
		certsMap["connector"] = connectorList
	}

	return certsMap
}

func dataSourceEndpointCertsClientToMap(clientItem satellitelinkv1.EndpointCertsClient) (clientMap map[string]interface{}) {
	clientMap = map[string]interface{}{}

	if clientItem.Cert != nil {
		certList := []map[string]interface{}{}
		certMap := dataSourceEndpointClientCertToMap(*clientItem.Cert)
		certList = append(certList, certMap)
		clientMap["cert"] = certList
	}

	return clientMap
}

func dataSourceEndpointClientCertToMap(certItem satellitelinkv1.EndpointCertsClientCert) (certMap map[string]interface{}) {
	certMap = map[string]interface{}{}

	if certItem.Filename != nil {
		certMap["filename"] = certItem.Filename
	}

	return certMap
}

func dataSourceEndpointCertsServerToMap(serverItem satellitelinkv1.EndpointCertsServer) (serverMap map[string]interface{}) {
	serverMap = map[string]interface{}{}

	if serverItem.Cert != nil {
		certList := []map[string]interface{}{}
		certMap := dataSourceEndpointServerCertToMap(*serverItem.Cert)
		certList = append(certList, certMap)
		serverMap["cert"] = certList
	}

	return serverMap
}

func dataSourceEndpointServerCertToMap(certItem satellitelinkv1.EndpointCertsServerCert) (certMap map[string]interface{}) {
	certMap = map[string]interface{}{}

	if certItem.Filename != nil {
		certMap["filename"] = certItem.Filename
	}

	return certMap
}

func dataSourceEndpointCertsConnectorToMap(connectorItem satellitelinkv1.EndpointCertsConnector) (connectorMap map[string]interface{}) {
	connectorMap = map[string]interface{}{}

	if connectorItem.Cert != nil {
		certList := []map[string]interface{}{}
		certMap := dataSourceEndpointConnectorCertToMap(*connectorItem.Cert)
		certList = append(certList, certMap)
		connectorMap["cert"] = certList
	}
	if connectorItem.Key != nil {
		keyList := []map[string]interface{}{}
		keyMap := dataSourceEndpointConnectorKeyToMap(*connectorItem.Key)
		keyList = append(keyList, keyMap)
		connectorMap["key"] = keyList
	}

	return connectorMap
}

func dataSourceEndpointConnectorCertToMap(certItem satellitelinkv1.EndpointCertsConnectorCert) (certMap map[string]interface{}) {
	certMap = map[string]interface{}{}

	if certItem.Filename != nil {
		certMap["filename"] = certItem.Filename
	}

	return certMap
}

func dataSourceEndpointConnectorKeyToMap(keyItem satellitelinkv1.EndpointCertsConnectorKey) (keyMap map[string]interface{}) {
	keyMap = map[string]interface{}{}

	if keyItem.Filename != nil {
		keyMap["filename"] = keyItem.Filename
	}

	return keyMap
}

func dataSourceEndpointFlattenPerformance(result satellitelinkv1.EndpointPerformance) (finalList []map[string]interface{}) {
	finalList = []map[string]interface{}{}
	finalMap := dataSourceEndpointPerformanceToMap(result)
	finalList = append(finalList, finalMap)

	return finalList
}

func dataSourceEndpointPerformanceToMap(performanceItem satellitelinkv1.EndpointPerformance) (performanceMap map[string]interface{}) {
	performanceMap = map[string]interface{}{}

	if performanceItem.Connection != nil {
		performanceMap["connection"] = performanceItem.Connection
	}
	if performanceItem.RxBandwidth != nil {
		performanceMap["rx_bandwidth"] = performanceItem.RxBandwidth
	}
	if performanceItem.TxBandwidth != nil {
		performanceMap["tx_bandwidth"] = performanceItem.TxBandwidth
	}
	if performanceItem.Bandwidth != nil {
		performanceMap["bandwidth"] = performanceItem.Bandwidth
	}
	if performanceItem.Connectors != nil {
		connectorsList := []map[string]interface{}{}
		for _, connectorsItem := range performanceItem.Connectors {
			connectorsList = append(connectorsList, dataSourceEndpointPerformanceConnectorsToMap(connectorsItem))
		}
		performanceMap["connectors"] = connectorsList
	}

	return performanceMap
}

func dataSourceEndpointPerformanceConnectorsToMap(connectorsItem satellitelinkv1.EndpointPerformanceConnectorsItem) (connectorsMap map[string]interface{}) {
	connectorsMap = map[string]interface{}{}

	if connectorsItem.Connector != nil {
		connectorsMap["connector"] = connectorsItem.Connector
	}
	if connectorsItem.Connections != nil {
		connectorsMap["connections"] = connectorsItem.Connections
	}
	if connectorsItem.RxBW != nil {
		connectorsMap["rx_bw"] = connectorsItem.RxBW
	}
	if connectorsItem.TxBW != nil {
		connectorsMap["tx_bw"] = connectorsItem.TxBW
	}

	return connectorsMap
}
