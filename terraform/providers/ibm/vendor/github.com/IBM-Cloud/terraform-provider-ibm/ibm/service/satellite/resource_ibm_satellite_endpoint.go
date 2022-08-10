// Copyright IBM Corp. 2021 All Rights Reserved.
// Licensed under the Mozilla Public License v2.0

package satellite

import (
	"context"
	"fmt"
	"log"

	"github.com/IBM-Cloud/container-services-go-sdk/satellitelinkv1"
	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/conns"
	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/flex"
	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/validate"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/IBM/go-sdk-core/v5/core"
)

func ResourceIBMSatelliteEndpoint() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceIbmSatelliteEndpointCreate,
		ReadContext:   resourceIbmSatelliteEndpointRead,
		UpdateContext: resourceIbmSatelliteEndpointUpdate,
		DeleteContext: resourceIbmSatelliteEndpointDelete,
		Importer:      &schema.ResourceImporter{},

		Schema: map[string]*schema.Schema{
			"location": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The Location ID.",
			},
			"connection_type": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validate.InvokeValidator("ibm_satellite_endpoint", "connection_type"),
				Description:  "The type of the endpoint.",
			},
			"display_name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The display name of the endpoint. Endpoint names must start with a letter and end with an alphanumeric character, can contain letters, numbers, and hyphen (-), and must be 63 characters or fewer.",
			},
			"server_host": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The host name or IP address of the server endpoint. For 'http-tunnel' protocol, server_host can start with '*.' , which means a wildcard to it's sub domains. Such as '*.example.com' can accept request to 'api.example.com' and 'www.example.com'.",
			},
			"server_port": {
				Type:        schema.TypeInt,
				Required:    true,
				Description: "The port number of the server endpoint. For 'http-tunnel' protocol, server_port can be 0, which means any port. Such as 0 is good for 80 (http) and 443 (https).",
			},
			"sni": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The server name indicator (SNI) which used to connect to the server endpoint. Only useful if server side requires SNI.",
			},
			"client_protocol": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validate.InvokeValidator("ibm_satellite_endpoint", "client_protocol"),
				Description:  "The protocol in the client application side.",
			},
			"client_mutual_auth": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
				Description: "Whether enable mutual auth in the client application side, when client_protocol is 'tls' or 'https', this field is required.",
			},
			"server_protocol": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validate.InvokeValidator("ibm_satellite_endpoint", "server_protocol"),
				Description:  "The protocol in the server application side. This parameter will change to default value if it is omitted even when using PATCH API. If client_protocol is 'udp', server_protocol must be 'udp'. If client_protocol is 'tcp'/'http', server_protocol could be 'tcp'/'tls' and default to 'tcp'. If client_protocol is 'tls'/'https', server_protocol could be 'tcp'/'tls' and default to 'tls'. If client_protocol is 'http-tunnel', server_protocol must be 'tcp'.",
			},
			"server_mutual_auth": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
				Description: "Whether enable mutual auth in the server application side, when client_protocol is 'tls', this field is required.",
			},
			"reject_unauth": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
				Description: "Whether reject any connection to the server application which is not authorized with the list of supplied CAs in the fields certs.server_cert.",
			},
			"timeout": {
				Type:         schema.TypeInt,
				Optional:     true,
				ValidateFunc: validate.InvokeValidator("ibm_satellite_endpoint", "timeout"),
				Description:  "The inactivity timeout in the Endpoint side.",
			},
			"created_by": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The service or person who created the endpoint. Must be 1000 characters or fewer.",
			},
			"certs": {
				Type:        schema.TypeList,
				MaxItems:    1,
				Optional:    true,
				Description: "The certs.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"client": {
							Type:        schema.TypeList,
							MaxItems:    1,
							Optional:    true,
							Description: "The CA which Satellite Link trust when receiving the connection from the client application.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"cert": {
										Type:        schema.TypeList,
										MaxItems:    1,
										Optional:    true,
										Description: "The root cert or the self-signed cert of the client application.",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"filename": {
													Type:        schema.TypeString,
													Optional:    true,
													Description: "The filename of the cert.",
												},
												"file_contents": {
													Type:        schema.TypeString,
													Optional:    true,
													Description: "The content of the cert. The certificate file must be in Privacy-enhanced Electronic Mail (PEM) format.",
												},
											},
										},
									},
								},
							},
						},
						"server": {
							Type:        schema.TypeList,
							MaxItems:    1,
							Optional:    true,
							Description: "The CA which Satellite Link trust when sending the connection to server application.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"cert": {
										Type:        schema.TypeList,
										MaxItems:    1,
										Optional:    true,
										Description: "The root cert or the self-signed cert of the server application.",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"filename": {
													Type:        schema.TypeString,
													Optional:    true,
													Description: "The filename of the cert.",
												},
												"file_contents": {
													Type:        schema.TypeString,
													Optional:    true,
													Description: "The content of the cert. The certificate file must be in Privacy-enhanced Electronic Mail (PEM) format.",
												},
											},
										},
									},
								},
							},
						},
						"connector": {
							Type:        schema.TypeList,
							MaxItems:    1,
							Optional:    true,
							Description: "The cert which Satellite Link connector provide to identify itself for connecting to the client/server application.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"cert": {
										Type:        schema.TypeList,
										MaxItems:    1,
										Optional:    true,
										Description: "The end-entity cert. This is required when the key is defined.",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"filename": {
													Type:        schema.TypeString,
													Optional:    true,
													Description: "The filename of the cert.",
												},
												"file_contents": {
													Type:        schema.TypeString,
													Optional:    true,
													Description: "The content of the cert. The certificate file must be in Privacy-enhanced Electronic Mail (PEM) format.",
												},
											},
										},
									},
									"key": {
										Type:        schema.TypeList,
										MaxItems:    1,
										Optional:    true,
										Description: "The private key of the end-entity certificate. This is required when the cert is defined.",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"filename": {
													Type:        schema.TypeString,
													Optional:    true,
													Description: "The name of the key.",
												},
												"file_contents": {
													Type:        schema.TypeString,
													Optional:    true,
													Description: "The content of the key. The private key file must be in Privacy-enhanced Electronic Mail (PEM) format.",
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
			"sources": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"source_id": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "The Source ID.",
						},
						"enabled": {
							Type:        schema.TypeBool,
							Optional:    true,
							Description: "Whether the source is enabled for the endpoint.",
						},
						"last_change": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "The last time modify the Endpoint configurations.",
						},
						"pending": {
							Type:        schema.TypeBool,
							Optional:    true,
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
			"endpoint_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The Endpoint ID.",
			},
			"performance": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "The last performance data of the endpoint.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"connection": {
							Type:        schema.TypeInt,
							Optional:    true,
							Description: "Concurrent connections number of moment when probe read the data.",
						},
						"rx_bandwidth": {
							Type:        schema.TypeInt,
							Optional:    true,
							Description: "Average Receive (to Cloud) Bandwidth of last two minutes, unit is Byte/s.",
						},
						"tx_bandwidth": {
							Type:        schema.TypeInt,
							Optional:    true,
							Description: "Average Transmitted (to Location) Bandwidth of last two minutes, unit is Byte/s.",
						},
						"bandwidth": {
							Type:        schema.TypeInt,
							Optional:    true,
							Description: "Average Tatal Bandwidth of last two minutes, unit is Byte/s.",
						},
						"connectors": {
							Type:        schema.TypeList,
							Optional:    true,
							Description: "The last performance data of the endpoint from each Connector.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"connector": {
										Type:        schema.TypeString,
										Optional:    true,
										Description: "The name of the connector reported the performance data.",
									},
									"connections": {
										Type:        schema.TypeInt,
										Optional:    true,
										Description: "Concurrent connections number of moment when probe read the data from the Connector.",
									},
									"rx_bw": {
										Type:        schema.TypeInt,
										Optional:    true,
										Description: "Average Transmitted (to Location) Bandwidth of last two minutes read from the Connector, unit is Byte/s.",
									},
									"tx_bw": {
										Type:        schema.TypeInt,
										Optional:    true,
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

func ResourceIBMSatelliteEndpointValidator() *validate.ResourceValidator {
	validateSchema := make([]validate.ValidateSchema, 0)
	validateSchema = append(validateSchema,
		validate.ValidateSchema{
			Identifier:                 "connection_type",
			ValidateFunctionIdentifier: validate.ValidateAllowedStringValue,
			Type:                       validate.TypeString,
			Optional:                   true,
			AllowedValues:              "cloud, location",
		},
		validate.ValidateSchema{
			Identifier:                 "client_protocol",
			ValidateFunctionIdentifier: validate.ValidateAllowedStringValue,
			Type:                       validate.TypeString,
			Optional:                   true,
			AllowedValues:              "http, http-tunnel, https, tcp, tls, udp",
		},
		validate.ValidateSchema{
			Identifier:                 "server_protocol",
			ValidateFunctionIdentifier: validate.ValidateAllowedStringValue,
			Type:                       validate.TypeString,
			Optional:                   true,
			AllowedValues:              "tcp, tls, udp",
		},
		validate.ValidateSchema{
			Identifier:                 "timeout",
			ValidateFunctionIdentifier: validate.IntBetween,
			Type:                       validate.TypeInt,
			Optional:                   true,
			MinValue:                   "1",
			MaxValue:                   "180",
		},
	)

	resourceValidator := validate.ResourceValidator{ResourceName: "ibm_satellite_endpoint", Schema: validateSchema}
	return &resourceValidator
}

func resourceIbmSatelliteEndpointCreate(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	satelliteLinkClient, err := meta.(conns.ClientSession).SatellitLinkClientSession()
	if err != nil {
		return diag.FromErr(err)
	}

	createEndpointsOptions := &satellitelinkv1.CreateEndpointsOptions{}

	createEndpointsOptions.SetLocationID(d.Get("location").(string))
	if _, ok := d.GetOk("connection_type"); ok {
		createEndpointsOptions.SetConnType(d.Get("connection_type").(string))
	}
	if _, ok := d.GetOk("display_name"); ok {
		createEndpointsOptions.SetDisplayName(d.Get("display_name").(string))
	}
	if _, ok := d.GetOk("server_host"); ok {
		createEndpointsOptions.SetServerHost(d.Get("server_host").(string))
	}
	if _, ok := d.GetOk("server_port"); ok {
		createEndpointsOptions.SetServerPort(int64(d.Get("server_port").(int)))
	}
	if _, ok := d.GetOk("sni"); ok {
		createEndpointsOptions.SetSni(d.Get("sni").(string))
	}
	if _, ok := d.GetOk("client_protocol"); ok {
		createEndpointsOptions.SetClientProtocol(d.Get("client_protocol").(string))
	}
	if _, ok := d.GetOk("client_mutual_auth"); ok {
		createEndpointsOptions.SetClientMutualAuth(d.Get("client_mutual_auth").(bool))
	}
	if _, ok := d.GetOk("server_protocol"); ok {
		createEndpointsOptions.SetServerProtocol(d.Get("server_protocol").(string))
	}
	if _, ok := d.GetOk("server_mutual_auth"); ok {
		createEndpointsOptions.SetServerMutualAuth(d.Get("server_mutual_auth").(bool))
	}
	if _, ok := d.GetOk("reject_unauth"); ok {
		createEndpointsOptions.SetRejectUnauth(d.Get("reject_unauth").(bool))
	}
	if _, ok := d.GetOk("timeout"); ok {
		createEndpointsOptions.SetTimeout(int64(d.Get("timeout").(int)))
	}
	if _, ok := d.GetOk("created_by"); ok {
		createEndpointsOptions.SetCreatedBy(d.Get("created_by").(string))
	}
	if _, ok := d.GetOk("certs"); ok {
		certs := resourceIbmSatelliteEndpointMapToAdditionalNewEndpointRequestCerts(d.Get("certs.0").(map[string]interface{}))
		createEndpointsOptions.SetCerts(&certs)
	}

	endpoint, response, err := satelliteLinkClient.CreateEndpointsWithContext(context, createEndpointsOptions)
	if err != nil {
		log.Printf("[DEBUG] CreateEndpointsWithContext failed %s\n%s", err, response)
		return diag.FromErr(fmt.Errorf("CreateEndpointsWithContext failed %s\n%s", err, response))
	}

	d.SetId(fmt.Sprintf("%s/%s", *createEndpointsOptions.LocationID, *endpoint.EndpointID))

	return resourceIbmSatelliteEndpointRead(context, d, meta)
}

func resourceIbmSatelliteEndpointMapToAdditionalNewEndpointRequestCerts(additionalNewEndpointRequestCertsMap map[string]interface{}) satellitelinkv1.AdditionalNewEndpointRequestCerts {
	additionalNewEndpointRequestCerts := satellitelinkv1.AdditionalNewEndpointRequestCerts{}

	if additionalNewEndpointRequestCertsMap["client"] != nil {
		for _, clientItem := range additionalNewEndpointRequestCertsMap["client"].([]interface{}) {
			client := resourceIbmSatelliteEndpointMapToAdditionalNewEndpointRequestCertsClient(clientItem.(map[string]interface{}))
			additionalNewEndpointRequestCerts.Client = &client
		}
	}
	if additionalNewEndpointRequestCertsMap["server"] != nil {
		for _, serverItem := range additionalNewEndpointRequestCertsMap["server"].([]interface{}) {
			server := resourceIbmSatelliteEndpointMapToAdditionalNewEndpointRequestCertsServer(serverItem.(map[string]interface{}))
			additionalNewEndpointRequestCerts.Server = &server
		}
	}
	if additionalNewEndpointRequestCertsMap["connector"] != nil {
		for _, connectorItem := range additionalNewEndpointRequestCertsMap["connector"].([]interface{}) {
			connector := resourceIbmSatelliteEndpointMapToAdditionalNewEndpointRequestCertsConnector(connectorItem.(map[string]interface{}))
			additionalNewEndpointRequestCerts.Connector = &connector
		}
	}

	return additionalNewEndpointRequestCerts
}

func resourceIbmSatelliteEndpointMapToAdditionalNewEndpointRequestCertsClient(additionalNewEndpointRequestCertsClientMap map[string]interface{}) satellitelinkv1.AdditionalNewEndpointRequestCertsClient {
	additionalNewEndpointRequestCertsClient := satellitelinkv1.AdditionalNewEndpointRequestCertsClient{}

	if additionalNewEndpointRequestCertsClientMap["cert"] != nil {
		for _, certItem := range additionalNewEndpointRequestCertsClientMap["cert"].([]interface{}) {
			cert := resourceIbmSatelliteEndpointMapToAdditionalNewEndpointRequestCertsClientCert(certItem.(map[string]interface{}))
			additionalNewEndpointRequestCertsClient.Cert = &cert
		}
	}

	return additionalNewEndpointRequestCertsClient
}

func resourceIbmSatelliteEndpointMapToAdditionalNewEndpointRequestCertsClientCert(additionalNewEndpointRequestCertsClientCertMap map[string]interface{}) satellitelinkv1.AdditionalNewEndpointRequestCertsClientCert {
	additionalNewEndpointRequestCertsClientCert := satellitelinkv1.AdditionalNewEndpointRequestCertsClientCert{}

	if additionalNewEndpointRequestCertsClientCertMap["filename"] != nil {
		additionalNewEndpointRequestCertsClientCert.Filename = core.StringPtr(additionalNewEndpointRequestCertsClientCertMap["filename"].(string))
	}
	if additionalNewEndpointRequestCertsClientCertMap["file_contents"] != nil {
		additionalNewEndpointRequestCertsClientCert.FileContents = core.StringPtr(additionalNewEndpointRequestCertsClientCertMap["file_contents"].(string))
	}

	return additionalNewEndpointRequestCertsClientCert
}

func resourceIbmSatelliteEndpointMapToAdditionalNewEndpointRequestCertsServer(additionalNewEndpointRequestCertsServerMap map[string]interface{}) satellitelinkv1.AdditionalNewEndpointRequestCertsServer {
	additionalNewEndpointRequestCertsServer := satellitelinkv1.AdditionalNewEndpointRequestCertsServer{}

	if additionalNewEndpointRequestCertsServerMap["cert"] != nil {
		for _, serverItem := range additionalNewEndpointRequestCertsServerMap["cert"].([]interface{}) {
			server := resourceIbmSatelliteEndpointMapToAdditionalNewEndpointRequestCertsServerCert(serverItem.(map[string]interface{}))
			additionalNewEndpointRequestCertsServer.Cert = &server
		}
	}

	return additionalNewEndpointRequestCertsServer
}

func resourceIbmSatelliteEndpointMapToAdditionalNewEndpointRequestCertsServerCert(additionalNewEndpointRequestCertsServerCertMap map[string]interface{}) satellitelinkv1.AdditionalNewEndpointRequestCertsServerCert {
	additionalNewEndpointRequestCertsServerCert := satellitelinkv1.AdditionalNewEndpointRequestCertsServerCert{}

	if additionalNewEndpointRequestCertsServerCertMap["filename"] != nil {
		additionalNewEndpointRequestCertsServerCert.Filename = core.StringPtr(additionalNewEndpointRequestCertsServerCertMap["filename"].(string))
	}
	if additionalNewEndpointRequestCertsServerCertMap["file_contents"] != nil {
		additionalNewEndpointRequestCertsServerCert.FileContents = core.StringPtr(additionalNewEndpointRequestCertsServerCertMap["file_contents"].(string))
	}

	return additionalNewEndpointRequestCertsServerCert
}

func resourceIbmSatelliteEndpointMapToAdditionalNewEndpointRequestCertsConnector(additionalNewEndpointRequestCertsConnectorMap map[string]interface{}) satellitelinkv1.AdditionalNewEndpointRequestCertsConnector {
	additionalNewEndpointRequestCertsConnector := satellitelinkv1.AdditionalNewEndpointRequestCertsConnector{}

	if additionalNewEndpointRequestCertsConnectorMap["cert"] != nil {
		for _, connItem := range additionalNewEndpointRequestCertsConnectorMap["cert"].([]interface{}) {
			cert := resourceIbmSatelliteEndpointMapToAdditionalNewEndpointRequestCertsConnectorCert(connItem.(map[string]interface{}))
			additionalNewEndpointRequestCertsConnector.Cert = &cert
		}
	}
	if additionalNewEndpointRequestCertsConnectorMap["key"] != nil {
		for _, keyItem := range additionalNewEndpointRequestCertsConnectorMap["key"].([]interface{}) {
			key := resourceIbmSatelliteEndpointMapToAdditionalNewEndpointRequestCertsConnectorKey(keyItem.(map[string]interface{}))
			additionalNewEndpointRequestCertsConnector.Key = &key
		}
	}

	return additionalNewEndpointRequestCertsConnector
}

func resourceIbmSatelliteEndpointMapToAdditionalNewEndpointRequestCertsConnectorCert(additionalNewEndpointRequestCertsConnectorCertMap map[string]interface{}) satellitelinkv1.AdditionalNewEndpointRequestCertsConnectorCert {
	additionalNewEndpointRequestCertsConnectorCert := satellitelinkv1.AdditionalNewEndpointRequestCertsConnectorCert{}

	if additionalNewEndpointRequestCertsConnectorCertMap["filename"] != nil {
		additionalNewEndpointRequestCertsConnectorCert.Filename = core.StringPtr(additionalNewEndpointRequestCertsConnectorCertMap["filename"].(string))
	}
	if additionalNewEndpointRequestCertsConnectorCertMap["file_contents"] != nil {
		additionalNewEndpointRequestCertsConnectorCert.FileContents = core.StringPtr(additionalNewEndpointRequestCertsConnectorCertMap["file_contents"].(string))
	}

	return additionalNewEndpointRequestCertsConnectorCert
}

func resourceIbmSatelliteEndpointMapToAdditionalNewEndpointRequestCertsConnectorKey(additionalNewEndpointRequestCertsConnectorKeyMap map[string]interface{}) satellitelinkv1.AdditionalNewEndpointRequestCertsConnectorKey {
	additionalNewEndpointRequestCertsConnectorKey := satellitelinkv1.AdditionalNewEndpointRequestCertsConnectorKey{}

	if additionalNewEndpointRequestCertsConnectorKeyMap["filename"] != nil {
		additionalNewEndpointRequestCertsConnectorKey.Filename = core.StringPtr(additionalNewEndpointRequestCertsConnectorKeyMap["filename"].(string))
	}
	if additionalNewEndpointRequestCertsConnectorKeyMap["file_contents"] != nil {
		additionalNewEndpointRequestCertsConnectorKey.FileContents = core.StringPtr(additionalNewEndpointRequestCertsConnectorKeyMap["file_contents"].(string))
	}

	return additionalNewEndpointRequestCertsConnectorKey
}

func resourceIbmSatelliteEndpointRead(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	satelliteLinkClient, err := meta.(conns.ClientSession).SatellitLinkClientSession()
	if err != nil {
		return diag.FromErr(err)
	}

	parts, err := flex.SepIdParts(d.Id(), "/")
	if err != nil {
		return diag.FromErr(err)
	}

	getEndpointsOptions := &satellitelinkv1.GetEndpointsOptions{}
	getEndpointsOptions.SetLocationID(parts[0])
	getEndpointsOptions.SetEndpointID(parts[1])

	endpoint, response, err := satelliteLinkClient.GetEndpointsWithContext(context, getEndpointsOptions)
	if err != nil {
		if response != nil && response.StatusCode == 404 {
			d.SetId("")
			return nil
		}
		log.Printf("[DEBUG] ListEndpointsWithContext failed %s\n%s", err, response)
		return diag.FromErr(fmt.Errorf("ListEndpointsWithContext failed %s\n%s", err, response))
	}

	if endpoint.EndpointID != nil {
		d.Set("endpoint_id", *endpoint.EndpointID)
	}

	if err = d.Set("location", endpoint.LocationID); err != nil {
		return diag.FromErr(fmt.Errorf("[ERROR] Error setting location: %s", err))
	}
	if err = d.Set("connection_type", endpoint.ConnType); err != nil {
		return diag.FromErr(fmt.Errorf("[ERROR] Error setting connection_type: %s", err))
	}
	if err = d.Set("display_name", endpoint.DisplayName); err != nil {
		return diag.FromErr(fmt.Errorf("[ERROR] Error setting display_name: %s", err))
	}
	if err = d.Set("server_host", endpoint.ServerHost); err != nil {
		return diag.FromErr(fmt.Errorf("[ERROR] Error setting server_host: %s", err))
	}
	if err = d.Set("server_port", flex.IntValue(endpoint.ServerPort)); err != nil {
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
	if err = d.Set("timeout", flex.IntValue(endpoint.Timeout)); err != nil {
		return diag.FromErr(fmt.Errorf("[ERROR] Error setting timeout: %s", err))
	}
	if err = d.Set("created_by", endpoint.CreatedBy); err != nil {
		return diag.FromErr(fmt.Errorf("[ERROR] Error setting created_by: %s", err))
	}

	if endpoint.Sources != nil {
		sources := []map[string]interface{}{}
		for _, sourcesItem := range endpoint.Sources {
			sourcesItemMap := resourceIbmSatelliteEndpointSourceStatusObjectToMap(sourcesItem)
			sources = append(sources, sourcesItemMap)
		}
		if err = d.Set("sources", sources); err != nil {
			return diag.FromErr(fmt.Errorf("[ERROR] Error setting sources: %s", err))
		}
	}
	if err = d.Set("connector_port", flex.IntValue(endpoint.ConnectorPort)); err != nil {
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
	if err = d.Set("client_port", flex.IntValue(endpoint.ClientPort)); err != nil {
		return diag.FromErr(fmt.Errorf("[ERROR] Error setting client_port: %s", err))
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
		performanceMap := resourceIbmSatelliteEndpointEndpointPerformanceToMap(*endpoint.Performance)
		if err = d.Set("performance", []map[string]interface{}{performanceMap}); err != nil {
			return diag.FromErr(fmt.Errorf("[ERROR] Error setting performance: %s", err))
		}
	}

	return nil
}

func resourceIbmSatelliteEndpointAdditionalNewEndpointRequestCertsToMap(additionalNewEndpointRequestCerts satellitelinkv1.AdditionalNewEndpointRequestCerts) map[string]interface{} {
	additionalNewEndpointRequestCertsMap := map[string]interface{}{}

	if additionalNewEndpointRequestCerts.Client != nil {
		clientMap := resourceIbmSatelliteEndpointAdditionalNewEndpointRequestCertsClientToMap(*additionalNewEndpointRequestCerts.Client)
		additionalNewEndpointRequestCertsMap["client"] = []map[string]interface{}{clientMap}
	}
	if additionalNewEndpointRequestCerts.Server != nil {
		serverMap := resourceIbmSatelliteEndpointAdditionalNewEndpointRequestCertsServerToMap(*additionalNewEndpointRequestCerts.Server)
		additionalNewEndpointRequestCertsMap["server"] = []map[string]interface{}{serverMap}
	}
	if additionalNewEndpointRequestCerts.Connector != nil {
		connectorMap := resourceIbmSatelliteEndpointAdditionalNewEndpointRequestCertsConnectorToMap(*additionalNewEndpointRequestCerts.Connector)
		additionalNewEndpointRequestCertsMap["connector"] = []map[string]interface{}{connectorMap}
	}

	return additionalNewEndpointRequestCertsMap
}

func resourceIbmSatelliteEndpointAdditionalNewEndpointRequestCertsClientToMap(additionalNewEndpointRequestCertsClient satellitelinkv1.AdditionalNewEndpointRequestCertsClient) map[string]interface{} {
	additionalNewEndpointRequestCertsClientMap := map[string]interface{}{}

	if additionalNewEndpointRequestCertsClient.Cert != nil {
		certMap := resourceIbmSatelliteEndpointAdditionalNewEndpointRequestCertsClientCertToMap(*additionalNewEndpointRequestCertsClient.Cert)
		additionalNewEndpointRequestCertsClientMap["cert"] = []map[string]interface{}{certMap}
	}

	return additionalNewEndpointRequestCertsClientMap
}

func resourceIbmSatelliteEndpointAdditionalNewEndpointRequestCertsClientCertToMap(additionalNewEndpointRequestCertsClientCert satellitelinkv1.AdditionalNewEndpointRequestCertsClientCert) map[string]interface{} {
	additionalNewEndpointRequestCertsClientCertMap := map[string]interface{}{}

	additionalNewEndpointRequestCertsClientCertMap["filename"] = additionalNewEndpointRequestCertsClientCert.Filename
	additionalNewEndpointRequestCertsClientCertMap["file_contents"] = additionalNewEndpointRequestCertsClientCert.FileContents

	return additionalNewEndpointRequestCertsClientCertMap
}

func resourceIbmSatelliteEndpointAdditionalNewEndpointRequestCertsServerToMap(additionalNewEndpointRequestCertsServer satellitelinkv1.AdditionalNewEndpointRequestCertsServer) map[string]interface{} {
	additionalNewEndpointRequestCertsServerMap := map[string]interface{}{}

	if additionalNewEndpointRequestCertsServer.Cert != nil {
		certMap := resourceIbmSatelliteEndpointAdditionalNewEndpointRequestCertsServerCertToMap(*additionalNewEndpointRequestCertsServer.Cert)
		additionalNewEndpointRequestCertsServerMap["cert"] = []map[string]interface{}{certMap}
	}

	return additionalNewEndpointRequestCertsServerMap
}

func resourceIbmSatelliteEndpointAdditionalNewEndpointRequestCertsServerCertToMap(additionalNewEndpointRequestCertsServerCert satellitelinkv1.AdditionalNewEndpointRequestCertsServerCert) map[string]interface{} {
	additionalNewEndpointRequestCertsServerCertMap := map[string]interface{}{}

	additionalNewEndpointRequestCertsServerCertMap["filename"] = additionalNewEndpointRequestCertsServerCert.Filename
	additionalNewEndpointRequestCertsServerCertMap["file_contents"] = additionalNewEndpointRequestCertsServerCert.FileContents

	return additionalNewEndpointRequestCertsServerCertMap
}

func resourceIbmSatelliteEndpointAdditionalNewEndpointRequestCertsConnectorToMap(additionalNewEndpointRequestCertsConnector satellitelinkv1.AdditionalNewEndpointRequestCertsConnector) map[string]interface{} {
	additionalNewEndpointRequestCertsConnectorMap := map[string]interface{}{}

	if additionalNewEndpointRequestCertsConnector.Cert != nil {
		certMap := resourceIbmSatelliteEndpointAdditionalNewEndpointRequestCertsConnectorCertToMap(*additionalNewEndpointRequestCertsConnector.Cert)
		additionalNewEndpointRequestCertsConnectorMap["cert"] = []map[string]interface{}{certMap}
	}
	if additionalNewEndpointRequestCertsConnector.Key != nil {
		keyMap := resourceIbmSatelliteEndpointAdditionalNewEndpointRequestCertsConnectorKeyToMap(*additionalNewEndpointRequestCertsConnector.Key)
		additionalNewEndpointRequestCertsConnectorMap["key"] = []map[string]interface{}{keyMap}
	}

	return additionalNewEndpointRequestCertsConnectorMap
}

func resourceIbmSatelliteEndpointAdditionalNewEndpointRequestCertsConnectorCertToMap(additionalNewEndpointRequestCertsConnectorCert satellitelinkv1.AdditionalNewEndpointRequestCertsConnectorCert) map[string]interface{} {
	additionalNewEndpointRequestCertsConnectorCertMap := map[string]interface{}{}

	additionalNewEndpointRequestCertsConnectorCertMap["filename"] = additionalNewEndpointRequestCertsConnectorCert.Filename
	additionalNewEndpointRequestCertsConnectorCertMap["file_contents"] = additionalNewEndpointRequestCertsConnectorCert.FileContents

	return additionalNewEndpointRequestCertsConnectorCertMap
}

func resourceIbmSatelliteEndpointAdditionalNewEndpointRequestCertsConnectorKeyToMap(additionalNewEndpointRequestCertsConnectorKey satellitelinkv1.AdditionalNewEndpointRequestCertsConnectorKey) map[string]interface{} {
	additionalNewEndpointRequestCertsConnectorKeyMap := map[string]interface{}{}

	additionalNewEndpointRequestCertsConnectorKeyMap["filename"] = additionalNewEndpointRequestCertsConnectorKey.Filename
	additionalNewEndpointRequestCertsConnectorKeyMap["file_contents"] = additionalNewEndpointRequestCertsConnectorKey.FileContents

	return additionalNewEndpointRequestCertsConnectorKeyMap
}

func resourceIbmSatelliteEndpointSourceStatusObjectToMap(sourceStatusObject satellitelinkv1.SourceStatusObject) map[string]interface{} {
	sourceStatusObjectMap := map[string]interface{}{}

	sourceStatusObjectMap["source_id"] = sourceStatusObject.SourceID
	sourceStatusObjectMap["enabled"] = sourceStatusObject.Enabled
	sourceStatusObjectMap["last_change"] = sourceStatusObject.LastChange
	sourceStatusObjectMap["pending"] = sourceStatusObject.Pending

	return sourceStatusObjectMap
}

func resourceIbmSatelliteEndpointEndpointPerformanceToMap(endpointPerformance satellitelinkv1.EndpointPerformance) map[string]interface{} {
	endpointPerformanceMap := map[string]interface{}{}

	endpointPerformanceMap["connection"] = flex.IntValue(endpointPerformance.Connection)
	endpointPerformanceMap["rx_bandwidth"] = flex.IntValue(endpointPerformance.RxBandwidth)
	endpointPerformanceMap["tx_bandwidth"] = flex.IntValue(endpointPerformance.TxBandwidth)
	endpointPerformanceMap["bandwidth"] = flex.IntValue(endpointPerformance.Bandwidth)
	if endpointPerformance.Connectors != nil {
		connectors := []map[string]interface{}{}
		for _, connectorsItem := range endpointPerformance.Connectors {
			connectorsItemMap := resourceIbmSatelliteEndpointEndpointPerformanceConnectorsItemToMap(connectorsItem)
			connectors = append(connectors, connectorsItemMap)
		}
		endpointPerformanceMap["connectors"] = connectors
	}

	return endpointPerformanceMap
}

func resourceIbmSatelliteEndpointEndpointPerformanceConnectorsItemToMap(endpointPerformanceConnectorsItem satellitelinkv1.EndpointPerformanceConnectorsItem) map[string]interface{} {
	endpointPerformanceConnectorsItemMap := map[string]interface{}{}

	endpointPerformanceConnectorsItemMap["connector"] = endpointPerformanceConnectorsItem.Connector
	endpointPerformanceConnectorsItemMap["connections"] = flex.IntValue(endpointPerformanceConnectorsItem.Connections)
	endpointPerformanceConnectorsItemMap["rxBW"] = flex.IntValue(endpointPerformanceConnectorsItem.RxBW)
	endpointPerformanceConnectorsItemMap["txBW"] = flex.IntValue(endpointPerformanceConnectorsItem.TxBW)

	return endpointPerformanceConnectorsItemMap
}

func resourceIbmSatelliteEndpointUpdate(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	satelliteLinkClient, err := meta.(conns.ClientSession).SatellitLinkClientSession()
	if err != nil {
		return diag.FromErr(err)
	}

	updateEndpointsOptions := &satellitelinkv1.UpdateEndpointsOptions{}

	parts, err := flex.SepIdParts(d.Id(), "/")
	if err != nil {
		return diag.FromErr(err)
	}

	updateEndpointsOptions.SetLocationID(parts[0])
	updateEndpointsOptions.SetEndpointID(parts[1])

	hasChange := false

	if d.HasChange("location") {
		updateEndpointsOptions.SetLocationID(d.Get("location").(string))
		hasChange = true
	}
	if d.HasChange("display_name") {
		updateEndpointsOptions.SetDisplayName(d.Get("display_name").(string))
		hasChange = true
	}
	if d.HasChange("server_host") {
		updateEndpointsOptions.SetServerHost(d.Get("server_host").(string))
		hasChange = true
	}
	if d.HasChange("server_port") {
		updateEndpointsOptions.SetServerPort(int64(d.Get("server_port").(int)))
		hasChange = true
	}
	if d.HasChange("sni") {
		updateEndpointsOptions.SetSni(d.Get("sni").(string))
		hasChange = true
	}
	if d.HasChange("client_protocol") {
		updateEndpointsOptions.SetClientProtocol(d.Get("client_protocol").(string))
		hasChange = true
	}
	if d.HasChange("client_mutual_auth") {
		updateEndpointsOptions.SetClientMutualAuth(d.Get("client_mutual_auth").(bool))
		hasChange = true
	}
	if d.HasChange("server_protocol") {
		updateEndpointsOptions.SetServerProtocol(d.Get("server_protocol").(string))
		hasChange = true
	}
	if d.HasChange("server_mutual_auth") {
		updateEndpointsOptions.SetServerMutualAuth(d.Get("server_mutual_auth").(bool))
		hasChange = true
	}
	if d.HasChange("reject_unauth") {
		updateEndpointsOptions.SetRejectUnauth(d.Get("reject_unauth").(bool))
		hasChange = true
	}
	if d.HasChange("timeout") {
		updateEndpointsOptions.SetTimeout(int64(d.Get("timeout").(int)))
		hasChange = true
	}
	if d.HasChange("created_by") {
		updateEndpointsOptions.SetCreatedBy(d.Get("created_by").(string))
		hasChange = true
	}
	if d.HasChange("certs") {
		certs := resourceIbmSatelliteEndpointUpdateEndpointRequestCerts(d.Get("certs.0").(map[string]interface{}))
		updateEndpointsOptions.SetCerts(&certs)
		hasChange = true
	}

	if hasChange {
		_, response, err := satelliteLinkClient.UpdateEndpointsWithContext(context, updateEndpointsOptions)
		if err != nil {
			log.Printf("[DEBUG] UpdateEndpointsWithContext failed %s\n%s", err, response)
			return diag.FromErr(fmt.Errorf("UpdateEndpointsWithContext failed %s\n%s", err, response))
		}
	}

	return resourceIbmSatelliteEndpointRead(context, d, meta)
}

func resourceIbmSatelliteEndpointUpdateEndpointRequestCerts(udateEndpointRequestCertsMap map[string]interface{}) satellitelinkv1.UpdatedEndpointRequestCerts {
	updateEndpointRequestCerts := satellitelinkv1.UpdatedEndpointRequestCerts{}

	if udateEndpointRequestCertsMap["client"] != nil {
		for _, clientItem := range udateEndpointRequestCertsMap["client"].([]interface{}) {
			client := resourceIbmSatelliteEndpointMapToUpdateEndpointRequestCertsClient(clientItem.(map[string]interface{}))
			updateEndpointRequestCerts.Client = &client
		}
	}
	if udateEndpointRequestCertsMap["server"] != nil {
		for _, serverItem := range udateEndpointRequestCertsMap["server"].([]interface{}) {
			server := resourceIbmSatelliteEndpointMapToUpdateEndpointRequestCertsServer(serverItem.(map[string]interface{}))
			updateEndpointRequestCerts.Server = &server
		}
	}
	if udateEndpointRequestCertsMap["connector"] != nil {
		for _, connectorItem := range udateEndpointRequestCertsMap["connector"].([]interface{}) {
			connector := resourceIbmSatelliteEndpointMapToUpdateEndpointRequestCertsConnector(connectorItem.(map[string]interface{}))
			updateEndpointRequestCerts.Connector = &connector
		}
	}

	return updateEndpointRequestCerts
}

func resourceIbmSatelliteEndpointMapToUpdateEndpointRequestCertsClient(updateEndpointRequestCertsClientMap map[string]interface{}) satellitelinkv1.UpdatedEndpointRequestCertsClient {
	updateEndpointRequestCertsClient := satellitelinkv1.UpdatedEndpointRequestCertsClient{}

	if updateEndpointRequestCertsClientMap["cert"] != nil {
		for _, certItem := range updateEndpointRequestCertsClientMap["cert"].([]interface{}) {
			cert := resourceIbmSatelliteEndpointMapToUpdateEndpointRequestCertsClientCert(certItem.(map[string]interface{}))
			updateEndpointRequestCertsClient.Cert = &cert
		}
	}

	return updateEndpointRequestCertsClient
}

func resourceIbmSatelliteEndpointMapToUpdateEndpointRequestCertsClientCert(updateEndpointRequestCertsClientCertMap map[string]interface{}) satellitelinkv1.UpdatedEndpointRequestCertsClientCert {
	updateEndpointRequestCertsClientCert := satellitelinkv1.UpdatedEndpointRequestCertsClientCert{}

	if updateEndpointRequestCertsClientCertMap["filename"] != nil {
		updateEndpointRequestCertsClientCert.Filename = core.StringPtr(updateEndpointRequestCertsClientCertMap["filename"].(string))
	}
	if updateEndpointRequestCertsClientCertMap["file_contents"] != nil {
		updateEndpointRequestCertsClientCert.FileContents = core.StringPtr(updateEndpointRequestCertsClientCertMap["file_contents"].(string))
	}

	return updateEndpointRequestCertsClientCert
}

func resourceIbmSatelliteEndpointMapToUpdateEndpointRequestCertsServer(updateEndpointRequestCertsClientMap map[string]interface{}) satellitelinkv1.UpdatedEndpointRequestCertsServer {
	updateEndpointRequestCertsServer := satellitelinkv1.UpdatedEndpointRequestCertsServer{}

	if updateEndpointRequestCertsClientMap["cert"] != nil {
		for _, certItem := range updateEndpointRequestCertsClientMap["cert"].([]interface{}) {
			cert := resourceIbmSatelliteEndpointMapToUpdateEndpointRequestCertsServerCert(certItem.(map[string]interface{}))
			updateEndpointRequestCertsServer.Cert = &cert
		}
	}

	return updateEndpointRequestCertsServer
}

func resourceIbmSatelliteEndpointMapToUpdateEndpointRequestCertsServerCert(updateEndpointRequestCertsServerCertMap map[string]interface{}) satellitelinkv1.UpdatedEndpointRequestCertsServerCert {
	updateEndpointRequestCertsServerCert := satellitelinkv1.UpdatedEndpointRequestCertsServerCert{}

	if updateEndpointRequestCertsServerCertMap["filename"] != nil {
		updateEndpointRequestCertsServerCert.Filename = core.StringPtr(updateEndpointRequestCertsServerCertMap["filename"].(string))
	}
	if updateEndpointRequestCertsServerCertMap["file_contents"] != nil {
		updateEndpointRequestCertsServerCert.FileContents = core.StringPtr(updateEndpointRequestCertsServerCertMap["file_contents"].(string))
	}

	return updateEndpointRequestCertsServerCert
}

func resourceIbmSatelliteEndpointMapToUpdateEndpointRequestCertsConnector(updateEndpointRequestCertsConnectorMap map[string]interface{}) satellitelinkv1.UpdatedEndpointRequestCertsConnector {
	updateEndpointRequestCertsConnector := satellitelinkv1.UpdatedEndpointRequestCertsConnector{}

	if updateEndpointRequestCertsConnectorMap["cert"] != nil {
		for _, certItem := range updateEndpointRequestCertsConnectorMap["cert"].([]interface{}) {
			cert := resourceIbmSatelliteEndpointMapToUpdateEndpointRequestCertsConnectorCert(certItem.(map[string]interface{}))
			updateEndpointRequestCertsConnector.Cert = &cert
		}
	}

	return updateEndpointRequestCertsConnector
}

func resourceIbmSatelliteEndpointMapToUpdateEndpointRequestCertsConnectorCert(updateEndpointRequestCertsServerCertMap map[string]interface{}) satellitelinkv1.UpdatedEndpointRequestCertsConnectorCert {
	updateEndpointRequestCertsConnectorCert := satellitelinkv1.UpdatedEndpointRequestCertsConnectorCert{}

	if updateEndpointRequestCertsServerCertMap["filename"] != nil {
		updateEndpointRequestCertsConnectorCert.Filename = core.StringPtr(updateEndpointRequestCertsServerCertMap["filename"].(string))
	}
	if updateEndpointRequestCertsServerCertMap["file_contents"] != nil {
		updateEndpointRequestCertsConnectorCert.FileContents = core.StringPtr(updateEndpointRequestCertsServerCertMap["file_contents"].(string))
	}

	return updateEndpointRequestCertsConnectorCert
}

func resourceIbmSatelliteEndpointDelete(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	satelliteLinkClient, err := meta.(conns.ClientSession).SatellitLinkClientSession()
	if err != nil {
		return diag.FromErr(err)
	}

	deleteEndpointsOptions := &satellitelinkv1.DeleteEndpointsOptions{}

	parts, err := flex.SepIdParts(d.Id(), "/")
	if err != nil {
		return diag.FromErr(err)
	}

	deleteEndpointsOptions.SetLocationID(parts[0])
	deleteEndpointsOptions.SetEndpointID(parts[1])

	_, response, err := satelliteLinkClient.DeleteEndpointsWithContext(context, deleteEndpointsOptions)
	if err != nil {
		log.Printf("[DEBUG] DeleteEndpointsWithContext failed %s\n%s", err, response)
		return diag.FromErr(fmt.Errorf("DeleteEndpointsWithContext failed %s\n%s", err, response))
	}

	d.SetId("")

	return nil
}
