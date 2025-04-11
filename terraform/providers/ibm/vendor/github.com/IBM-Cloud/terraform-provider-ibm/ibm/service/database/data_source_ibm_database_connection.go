// Copyright IBM Corp. 2024 All Rights Reserved.
// Licensed under the Mozilla Public License v2.0

package database

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/conns"
	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/flex"
	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/validate"
	"github.com/IBM/cloud-databases-go-sdk/clouddatabasesv5"
)

func DataSourceIBMDatabaseConnection() *schema.Resource {
	return &schema.Resource{
		ReadContext: DataSourceIBMDatabaseConnectionRead,

		Schema: map[string]*schema.Schema{
			"deployment_id": &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
				Description: "Deployment ID.",
				ValidateFunc: validate.InvokeDataSourceValidator(
					"ibm_database_connection",
					"deployment_id"),
			},
			"user_type": &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
				Description: "User type.",
			},
			"user_id": &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
				Description: "User ID.",
			},
			"endpoint_type": &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
				Description: "Endpoint Type. The endpoint must be enabled on the deployment before its connection information can be fetched.",
				ValidateFunc: validate.InvokeDataSourceValidator(
					"ibm_database_connection",
					"endpoint_type"),
			},
			"certificate_root": &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Optional certificate root path to prepend certificate names. Certificates would be stored in this directory for use by other commands.",
			},
			"postgres": &schema.Schema{
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"type": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Type of connection being described.",
						},
						"composed": &schema.Schema{
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
						"scheme": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Scheme/protocol for URI connection.",
						},
						"hosts": &schema.Schema{
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"hostname": &schema.Schema{
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Hostname for connection.",
									},
									"port": &schema.Schema{
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "Port number for connection.",
									},
								},
							},
						},
						"path": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Path for URI connection.",
						},
						"query_options": &schema.Schema{
							Type:        schema.TypeMap,
							Computed:    true,
							Description: "Query options to add to the URI connection.",
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
						"authentication": &schema.Schema{
							Type:        schema.TypeList,
							Computed:    true,
							Description: "Authentication data for Connection String.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"method": &schema.Schema{
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Authentication method for this credential.",
									},
									"username": &schema.Schema{
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Username part of credential.",
									},
									"password": &schema.Schema{
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Password part of credential.",
									},
								},
							},
						},
						"certificate": &schema.Schema{
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"name": &schema.Schema{
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Name associated with the certificate.",
									},
									"certificate_base64": &schema.Schema{
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Base64 encoded version of the certificate.",
									},
								},
							},
						},
						"ssl": &schema.Schema{
							Type:        schema.TypeBool,
							Computed:    true,
							Description: "Indicates ssl is required for the connection.",
						},
						"browser_accessible": &schema.Schema{
							Type:        schema.TypeBool,
							Computed:    true,
							Description: "Indicates the address is accessible by browser.",
						},
						"database": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Name of the database to use in the URI connection.",
						},
					},
				},
			},
			"cli": &schema.Schema{
				Type:        schema.TypeList,
				Computed:    true,
				Description: "CLI Connection.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"type": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Type of connection being described.",
						},
						"composed": &schema.Schema{
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
						"environment": &schema.Schema{
							Type:        schema.TypeMap,
							Computed:    true,
							Description: "A map of environment variables for a CLI connection.",
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
						"bin": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The name of the executable the CLI should run.",
						},
						"arguments": &schema.Schema{
							Type:        schema.TypeList,
							Computed:    true,
							Description: "Sets of arguments to call the executable with. The outer array corresponds to a possible way to call the CLI; the inner array is the set of arguments to use with that call.",
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
						"certificate": &schema.Schema{
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"name": &schema.Schema{
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Name associated with the certificate.",
									},
									"certificate_base64": &schema.Schema{
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Base64 encoded version of the certificate.",
									},
								},
							},
						},
					},
				},
			},
			"rediss": &schema.Schema{
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"type": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Type of connection being described.",
						},
						"composed": &schema.Schema{
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
						"scheme": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Scheme/protocol for URI connection.",
						},
						"hosts": &schema.Schema{
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"hostname": &schema.Schema{
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Hostname for connection.",
									},
									"port": &schema.Schema{
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "Port number for connection.",
									},
								},
							},
						},
						"path": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Path for URI connection.",
						},
						"query_options": &schema.Schema{
							Type:        schema.TypeMap,
							Computed:    true,
							Description: "Query options to add to the URI connection.",
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
						"authentication": &schema.Schema{
							Type:        schema.TypeList,
							Computed:    true,
							Description: "Authentication data for Connection String.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"method": &schema.Schema{
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Authentication method for this credential.",
									},
									"username": &schema.Schema{
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Username part of credential.",
									},
									"password": &schema.Schema{
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Password part of credential.",
									},
								},
							},
						},
						"certificate": &schema.Schema{
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"name": &schema.Schema{
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Name associated with the certificate.",
									},
									"certificate_base64": &schema.Schema{
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Base64 encoded version of the certificate.",
									},
								},
							},
						},
						"ssl": &schema.Schema{
							Type:        schema.TypeBool,
							Computed:    true,
							Description: "Indicates ssl is required for the connection.",
						},
						"browser_accessible": &schema.Schema{
							Type:        schema.TypeBool,
							Computed:    true,
							Description: "Indicates the address is accessible by browser.",
						},
						"database": &schema.Schema{
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Number of the database to use in the URI connection.",
						},
					},
				},
			},
			"https": &schema.Schema{
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"type": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Type of connection being described.",
						},
						"composed": &schema.Schema{
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
						"scheme": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Scheme/protocol for URI connection.",
						},
						"hosts": &schema.Schema{
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"hostname": &schema.Schema{
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Hostname for connection.",
									},
									"port": &schema.Schema{
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "Port number for connection.",
									},
								},
							},
						},
						"path": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Path for URI connection.",
						},
						"query_options": &schema.Schema{
							Type:        schema.TypeMap,
							Computed:    true,
							Description: "Query options to add to the URI connection.",
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
						"authentication": &schema.Schema{
							Type:        schema.TypeList,
							Computed:    true,
							Description: "Authentication data for Connection String.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"method": &schema.Schema{
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Authentication method for this credential.",
									},
									"username": &schema.Schema{
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Username part of credential.",
									},
									"password": &schema.Schema{
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Password part of credential.",
									},
								},
							},
						},
						"certificate": &schema.Schema{
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"name": &schema.Schema{
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Name associated with the certificate.",
									},
									"certificate_base64": &schema.Schema{
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Base64 encoded version of the certificate.",
									},
								},
							},
						},
						"ssl": &schema.Schema{
							Type:        schema.TypeBool,
							Computed:    true,
							Description: "Indicates ssl is required for the connection.",
						},
						"browser_accessible": &schema.Schema{
							Type:        schema.TypeBool,
							Computed:    true,
							Description: "Indicates the address is accessible by browser.",
						},
					},
				},
			},
			"amqps": &schema.Schema{
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"type": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Type of connection being described.",
						},
						"composed": &schema.Schema{
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
						"scheme": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Scheme/protocol for URI connection.",
						},
						"hosts": &schema.Schema{
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"hostname": &schema.Schema{
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Hostname for connection.",
									},
									"port": &schema.Schema{
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "Port number for connection.",
									},
								},
							},
						},
						"path": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Path for URI connection.",
						},
						"query_options": &schema.Schema{
							Type:        schema.TypeMap,
							Computed:    true,
							Description: "Query options to add to the URI connection.",
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
						"authentication": &schema.Schema{
							Type:        schema.TypeList,
							Computed:    true,
							Description: "Authentication data for Connection String.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"method": &schema.Schema{
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Authentication method for this credential.",
									},
									"username": &schema.Schema{
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Username part of credential.",
									},
									"password": &schema.Schema{
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Password part of credential.",
									},
								},
							},
						},
						"certificate": &schema.Schema{
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"name": &schema.Schema{
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Name associated with the certificate.",
									},
									"certificate_base64": &schema.Schema{
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Base64 encoded version of the certificate.",
									},
								},
							},
						},
						"ssl": &schema.Schema{
							Type:        schema.TypeBool,
							Computed:    true,
							Description: "Indicates ssl is required for the connection.",
						},
						"browser_accessible": &schema.Schema{
							Type:        schema.TypeBool,
							Computed:    true,
							Description: "Indicates the address is accessible by browser.",
						},
					},
				},
			},
			"mqtts": &schema.Schema{
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"type": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Type of connection being described.",
						},
						"composed": &schema.Schema{
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
						"scheme": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Scheme/protocol for URI connection.",
						},
						"hosts": &schema.Schema{
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"hostname": &schema.Schema{
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Hostname for connection.",
									},
									"port": &schema.Schema{
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "Port number for connection.",
									},
								},
							},
						},
						"path": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Path for URI connection.",
						},
						"query_options": &schema.Schema{
							Type:        schema.TypeMap,
							Computed:    true,
							Description: "Query options to add to the URI connection.",
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
						"authentication": &schema.Schema{
							Type:        schema.TypeList,
							Computed:    true,
							Description: "Authentication data for Connection String.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"method": &schema.Schema{
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Authentication method for this credential.",
									},
									"username": &schema.Schema{
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Username part of credential.",
									},
									"password": &schema.Schema{
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Password part of credential.",
									},
								},
							},
						},
						"certificate": &schema.Schema{
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"name": &schema.Schema{
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Name associated with the certificate.",
									},
									"certificate_base64": &schema.Schema{
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Base64 encoded version of the certificate.",
									},
								},
							},
						},
						"ssl": &schema.Schema{
							Type:        schema.TypeBool,
							Computed:    true,
							Description: "Indicates ssl is required for the connection.",
						},
						"browser_accessible": &schema.Schema{
							Type:        schema.TypeBool,
							Computed:    true,
							Description: "Indicates the address is accessible by browser.",
						},
					},
				},
			},
			"stomp_ssl": &schema.Schema{
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"type": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Type of connection being described.",
						},
						"composed": &schema.Schema{
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
						"scheme": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Scheme/protocol for URI connection.",
						},
						"hosts": &schema.Schema{
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"hostname": &schema.Schema{
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Hostname for connection.",
									},
									"port": &schema.Schema{
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "Port number for connection.",
									},
								},
							},
						},
						"path": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Path for URI connection.",
						},
						"query_options": &schema.Schema{
							Type:        schema.TypeMap,
							Computed:    true,
							Description: "Query options to add to the URI connection.",
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
						"authentication": &schema.Schema{
							Type:        schema.TypeList,
							Computed:    true,
							Description: "Authentication data for Connection String.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"method": &schema.Schema{
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Authentication method for this credential.",
									},
									"username": &schema.Schema{
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Username part of credential.",
									},
									"password": &schema.Schema{
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Password part of credential.",
									},
								},
							},
						},
						"certificate": &schema.Schema{
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"name": &schema.Schema{
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Name associated with the certificate.",
									},
									"certificate_base64": &schema.Schema{
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Base64 encoded version of the certificate.",
									},
								},
							},
						},
						"ssl": &schema.Schema{
							Type:        schema.TypeBool,
							Computed:    true,
							Description: "Indicates ssl is required for the connection.",
						},
						"browser_accessible": &schema.Schema{
							Type:        schema.TypeBool,
							Computed:    true,
							Description: "Indicates the address is accessible by browser.",
						},
					},
				},
			},
			"grpc": &schema.Schema{
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"type": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Type of connection being described.",
						},
						"composed": &schema.Schema{
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
						"scheme": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Scheme/protocol for URI connection.",
						},
						"hosts": &schema.Schema{
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"hostname": &schema.Schema{
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Hostname for connection.",
									},
									"port": &schema.Schema{
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "Port number for connection.",
									},
								},
							},
						},
						"path": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Path for URI connection.",
						},
						"query_options": &schema.Schema{
							Type:        schema.TypeMap,
							Computed:    true,
							Description: "Query options to add to the URI connection.",
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
						"authentication": &schema.Schema{
							Type:        schema.TypeList,
							Computed:    true,
							Description: "Authentication data for Connection String.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"method": &schema.Schema{
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Authentication method for this credential.",
									},
									"username": &schema.Schema{
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Username part of credential.",
									},
									"password": &schema.Schema{
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Password part of credential.",
									},
								},
							},
						},
						"certificate": &schema.Schema{
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"name": &schema.Schema{
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Name associated with the certificate.",
									},
									"certificate_base64": &schema.Schema{
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Base64 encoded version of the certificate.",
									},
								},
							},
						},
						"ssl": &schema.Schema{
							Type:        schema.TypeBool,
							Computed:    true,
							Description: "Indicates ssl is required for the connection.",
						},
						"browser_accessible": &schema.Schema{
							Type:        schema.TypeBool,
							Computed:    true,
							Description: "Indicates the address is accessible by browser.",
						},
					},
				},
			},
			"mongodb": &schema.Schema{
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"type": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Type of connection being described.",
						},
						"composed": &schema.Schema{
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
						"scheme": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Scheme/protocol for URI connection.",
						},
						"hosts": &schema.Schema{
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"hostname": &schema.Schema{
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Hostname for connection.",
									},
									"port": &schema.Schema{
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "Port number for connection.",
									},
								},
							},
						},
						"path": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Path for URI connection.",
						},
						"query_options": &schema.Schema{
							Type:        schema.TypeMap,
							Computed:    true,
							Description: "Query options to add to the URI connection.",
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
						"authentication": &schema.Schema{
							Type:        schema.TypeList,
							Computed:    true,
							Description: "Authentication data for Connection String.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"method": &schema.Schema{
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Authentication method for this credential.",
									},
									"username": &schema.Schema{
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Username part of credential.",
									},
									"password": &schema.Schema{
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Password part of credential.",
									},
								},
							},
						},
						"certificate": &schema.Schema{
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"name": &schema.Schema{
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Name associated with the certificate.",
									},
									"certificate_base64": &schema.Schema{
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Base64 encoded version of the certificate.",
									},
								},
							},
						},
						"ssl": &schema.Schema{
							Type:        schema.TypeBool,
							Computed:    true,
							Description: "Indicates ssl is required for the connection.",
						},
						"browser_accessible": &schema.Schema{
							Type:        schema.TypeBool,
							Computed:    true,
							Description: "Indicates the address is accessible by browser.",
						},
						"database": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Name of the database to use in the URI connection.",
						},
						"replica_set": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Name of the replica set to use in the URI connection.",
						},
					},
				},
			},
			"bi_connector": &schema.Schema{
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"type": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Type of connection being described.",
						},
						"composed": &schema.Schema{
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
						"scheme": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Scheme/protocol for URI connection.",
						},
						"hosts": &schema.Schema{
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"hostname": &schema.Schema{
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Hostname for connection.",
									},
									"port": &schema.Schema{
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "Port number for connection.",
									},
								},
							},
						},
						"path": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Path for URI connection.",
						},
						"query_options": &schema.Schema{
							Type:        schema.TypeMap,
							Computed:    true,
							Description: "Query options to add to the URI connection.",
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
						"authentication": &schema.Schema{
							Type:        schema.TypeList,
							Computed:    true,
							Description: "Authentication data for Connection String.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"method": &schema.Schema{
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Authentication method for this credential.",
									},
									"username": &schema.Schema{
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Username part of credential.",
									},
									"password": &schema.Schema{
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Password part of credential.",
									},
								},
							},
						},
						"certificate": &schema.Schema{
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"name": &schema.Schema{
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Name associated with the certificate.",
									},
									"certificate_base64": &schema.Schema{
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Base64 encoded version of the certificate.",
									},
								},
							},
						},
						"ssl": &schema.Schema{
							Type:        schema.TypeBool,
							Computed:    true,
							Description: "Indicates ssl is required for the connection.",
						},
						"browser_accessible": &schema.Schema{
							Type:        schema.TypeBool,
							Computed:    true,
							Description: "Indicates the address is accessible by browser.",
						},
					},
				},
			},
			"analytics": &schema.Schema{
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"type": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Type of connection being described.",
						},
						"composed": &schema.Schema{
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
						"scheme": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Scheme/protocol for URI connection.",
						},
						"hosts": &schema.Schema{
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"hostname": &schema.Schema{
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Hostname for connection.",
									},
									"port": &schema.Schema{
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "Port number for connection.",
									},
								},
							},
						},
						"path": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Path for URI connection.",
						},
						"query_options": &schema.Schema{
							Type:        schema.TypeMap,
							Computed:    true,
							Description: "Query options to add to the URI connection.",
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
						"authentication": &schema.Schema{
							Type:        schema.TypeList,
							Computed:    true,
							Description: "Authentication data for Connection String.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"method": &schema.Schema{
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Authentication method for this credential.",
									},
									"username": &schema.Schema{
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Username part of credential.",
									},
									"password": &schema.Schema{
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Password part of credential.",
									},
								},
							},
						},
						"certificate": &schema.Schema{
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"name": &schema.Schema{
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Name associated with the certificate.",
									},
									"certificate_base64": &schema.Schema{
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Base64 encoded version of the certificate.",
									},
								},
							},
						},
						"ssl": &schema.Schema{
							Type:        schema.TypeBool,
							Computed:    true,
							Description: "Indicates ssl is required for the connection.",
						},
						"browser_accessible": &schema.Schema{
							Type:        schema.TypeBool,
							Computed:    true,
							Description: "Indicates the address is accessible by browser.",
						},
					},
				},
			},
			"ops_manager": &schema.Schema{
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"type": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Type of connection being described.",
						},
						"composed": &schema.Schema{
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
						"scheme": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Scheme/protocol for URI connection.",
						},
						"hosts": &schema.Schema{
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"hostname": &schema.Schema{
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Hostname for connection.",
									},
									"port": &schema.Schema{
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "Port number for connection.",
									},
								},
							},
						},
						"path": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Path for URI connection.",
						},
						"query_options": &schema.Schema{
							Type:        schema.TypeMap,
							Computed:    true,
							Description: "Query options to add to the URI connection.",
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
						"authentication": &schema.Schema{
							Type:        schema.TypeList,
							Computed:    true,
							Description: "Authentication data for Connection String.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"method": &schema.Schema{
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Authentication method for this credential.",
									},
									"username": &schema.Schema{
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Username part of credential.",
									},
									"password": &schema.Schema{
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Password part of credential.",
									},
								},
							},
						},
						"certificate": &schema.Schema{
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"name": &schema.Schema{
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Name associated with the certificate.",
									},
									"certificate_base64": &schema.Schema{
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Base64 encoded version of the certificate.",
									},
								},
							},
						},
						"ssl": &schema.Schema{
							Type:        schema.TypeBool,
							Computed:    true,
							Description: "Indicates ssl is required for the connection.",
						},
						"browser_accessible": &schema.Schema{
							Type:        schema.TypeBool,
							Computed:    true,
							Description: "Indicates the address is accessible by browser.",
						},
					},
				},
			},
			"mysql": &schema.Schema{
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"type": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Type of connection being described.",
						},
						"composed": &schema.Schema{
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
						"scheme": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Scheme/protocol for URI connection.",
						},
						"hosts": &schema.Schema{
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"hostname": &schema.Schema{
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Hostname for connection.",
									},
									"port": &schema.Schema{
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "Port number for connection.",
									},
								},
							},
						},
						"path": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Path for URI connection.",
						},
						"query_options": &schema.Schema{
							Type:        schema.TypeMap,
							Computed:    true,
							Description: "Query options to add to the URI connection.",
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
						"authentication": &schema.Schema{
							Type:        schema.TypeList,
							Computed:    true,
							Description: "Authentication data for Connection String.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"method": &schema.Schema{
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Authentication method for this credential.",
									},
									"username": &schema.Schema{
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Username part of credential.",
									},
									"password": &schema.Schema{
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Password part of credential.",
									},
								},
							},
						},
						"certificate": &schema.Schema{
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"name": &schema.Schema{
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Name associated with the certificate.",
									},
									"certificate_base64": &schema.Schema{
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Base64 encoded version of the certificate.",
									},
								},
							},
						},
						"ssl": &schema.Schema{
							Type:        schema.TypeBool,
							Computed:    true,
							Description: "Indicates ssl is required for the connection.",
						},
						"browser_accessible": &schema.Schema{
							Type:        schema.TypeBool,
							Computed:    true,
							Description: "Indicates the address is accessible by browser.",
						},
						"database": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Name of the database to use in the URI connection.",
						},
					},
				},
			},
			"secure": &schema.Schema{
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"hosts": &schema.Schema{
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"hostname": &schema.Schema{
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Hostname for connection.",
									},
									"port": &schema.Schema{
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "Port number for connection.",
									},
								},
							},
						},
						"authentication": &schema.Schema{
							Type:        schema.TypeList,
							Computed:    true,
							Description: "Authentication data for Connection String.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"method": &schema.Schema{
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Authentication method for this credential.",
									},
									"username": &schema.Schema{
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Username part of credential.",
									},
									"password": &schema.Schema{
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Password part of credential.",
									},
								},
							},
						},
						"bundle": &schema.Schema{
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"name": &schema.Schema{
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Name associated with the certificate.",
									},
									"bundle_base64": &schema.Schema{
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Base64 encoded version of the certificate bundle.",
									},
								},
							},
						},
					},
				},
			},
			"emp": &schema.Schema{
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"type": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Type of connection being described.",
						},
						"composed": &schema.Schema{
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
						"scheme": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Scheme/protocol for URI connection.",
						},
						"hosts": &schema.Schema{
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"hostname": &schema.Schema{
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Hostname for connection.",
									},
									"port": &schema.Schema{
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "Port number for connection.",
									},
								},
							},
						},
						"path": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Path for URI connection.",
						},
						"query_options": &schema.Schema{
							Type:        schema.TypeMap,
							Computed:    true,
							Description: "Query options to add to the URI connection.",
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
						"authentication": &schema.Schema{
							Type:        schema.TypeList,
							Computed:    true,
							Description: "Authentication data for Connection String.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"method": &schema.Schema{
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Authentication method for this credential.",
									},
									"username": &schema.Schema{
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Username part of credential.",
									},
									"password": &schema.Schema{
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Password part of credential.",
									},
								},
							},
						},
						"certificate": &schema.Schema{
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"name": &schema.Schema{
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Name associated with the certificate.",
									},
									"certificate_base64": &schema.Schema{
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Base64 encoded version of the certificate.",
									},
								},
							},
						},
						"ssl": &schema.Schema{
							Type:        schema.TypeBool,
							Computed:    true,
							Description: "Indicates ssl is required for the connection.",
						},
						"browser_accessible": &schema.Schema{
							Type:        schema.TypeBool,
							Computed:    true,
							Description: "Indicates the address is accessible by browser.",
						},
					},
				},
			},
		},
	}
}
func DataSourceIBMDatabaseConnectionValidator() *validate.ResourceValidator {

	validateSchema := make([]validate.ValidateSchema, 0)

	validateSchema = append(validateSchema,
		validate.ValidateSchema{
			Identifier:                 "deployment_id",
			ValidateFunctionIdentifier: validate.ValidateCloudData,
			Type:                       validate.TypeString,
			Required:                   true,
			CloudDataType:              "cloud-database",
			CloudDataRange:             []string{"resolved_to:id"}})
	validateSchema = append(validateSchema,
		validate.ValidateSchema{
			Identifier:                 "endpoint_type",
			ValidateFunctionIdentifier: validate.ValidateAllowedStringValue,
			Type:                       validate.TypeString,
			Required:                   true,
			AllowedValues:              "public, private, public-and-private"})

	iBMDatabaseConnectionsValidator := validate.ResourceValidator{ResourceName: "ibm_database_connection", Schema: validateSchema}
	return &iBMDatabaseConnectionsValidator
}

func DataSourceIBMDatabaseConnectionRead(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cloudDatabasesClient, err := meta.(conns.ClientSession).CloudDatabasesV5()
	if err != nil {
		tfErr := flex.TerraformErrorf(err, err.Error(), "(Data) ibm_database_connection", "read")
		log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
		return tfErr.GetDiag()
	}

	getConnectionOptions := &clouddatabasesv5.GetConnectionOptions{}

	getConnectionOptions.SetID(d.Get("deployment_id").(string))
	getConnectionOptions.SetUserType(d.Get("user_type").(string))
	getConnectionOptions.SetUserID(d.Get("user_id").(string))
	getConnectionOptions.SetEndpointType(d.Get("endpoint_type").(string))

	if _, ok := d.GetOk("certificate_root"); ok {
		getConnectionOptions.SetCertificateRoot(d.Get("certificate_root").(string))
	}

	connection, response, err := cloudDatabasesClient.GetConnectionWithContext(context, getConnectionOptions)
	if err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("GetConnectionWithContext failed:  %s\n%s", err.Error(), response), "(Data) ibm_database_connection", "read")
		log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
		return tfErr.GetDiag()
	}

	d.SetId(DataSourceIBMDatabaseConnectionID(d))
	conn := connection.Connection.(*clouddatabasesv5.Connection)

	postgres := []map[string]interface{}{}
	if conn.Postgres != nil {
		modelMap, err := DataSourceIBMDatabaseConnectionPostgreSQLConnectionURIToMap(conn.Postgres)
		if err != nil {
			tfErr := flex.TerraformErrorf(err, err.Error(), "(Data) ibm_database_connection", "read")
			return tfErr.GetDiag()
		}
		postgres = append(postgres, modelMap)
	}
	if err = d.Set("postgres", postgres); err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("Error setting postgres: %s", err), "(Data) ibm_database_connection", "read")
		return tfErr.GetDiag()
	}

	cli := []map[string]interface{}{}
	if conn.Cli != nil {
		modelMap, err := DataSourceIBMDatabaseConnectionConnectionCliToMap(conn.Cli)
		if err != nil {
			tfErr := flex.TerraformErrorf(err, err.Error(), "(Data) ibm_database_connection", "read")
			return tfErr.GetDiag()
		}
		cli = append(cli, modelMap)
	}
	if err = d.Set("cli", cli); err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("Error setting cli: %s", err), "(Data) ibm_database_connection", "read")
		return tfErr.GetDiag()
	}

	rediss := []map[string]interface{}{}
	if conn.Rediss != nil {
		modelMap, err := DataSourceIBMDatabaseConnectionRedisConnectionURIToMap(conn.Rediss)
		if err != nil {
			tfErr := flex.TerraformErrorf(err, err.Error(), "(Data) ibm_database_connection", "read")
			return tfErr.GetDiag()
		}
		rediss = append(rediss, modelMap)
	}
	if err = d.Set("rediss", rediss); err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("Error setting rediss: %s", err), "(Data) ibm_database_connection", "read")
		return tfErr.GetDiag()
	}

	https := []map[string]interface{}{}
	if conn.HTTPS != nil {
		modelMap, err := DataSourceIBMDatabaseConnectionConnectionURIToMap(conn.HTTPS)
		if err != nil {
			tfErr := flex.TerraformErrorf(err, err.Error(), "(Data) ibm_database_connection", "read")
			return tfErr.GetDiag()
		}
		https = append(https, modelMap)
	}
	if err = d.Set("https", https); err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("Error setting https: %s", err), "(Data) ibm_database_connection", "read")
		return tfErr.GetDiag()
	}

	amqps := []map[string]interface{}{}
	if conn.Amqps != nil {
		modelMap, err := DataSourceIBMDatabaseConnectionConnectionURIToMap(conn.Amqps)
		if err != nil {
			tfErr := flex.TerraformErrorf(err, err.Error(), "(Data) ibm_database_connection", "read")
			return tfErr.GetDiag()
		}
		amqps = append(amqps, modelMap)
	}
	if err = d.Set("amqps", amqps); err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("Error setting amqps: %s", err), "(Data) ibm_database_connection", "read")
		return tfErr.GetDiag()
	}

	mqtts := []map[string]interface{}{}
	if conn.Mqtts != nil {
		modelMap, err := DataSourceIBMDatabaseConnectionConnectionURIToMap(conn.Mqtts)
		if err != nil {
			tfErr := flex.TerraformErrorf(err, err.Error(), "(Data) ibm_database_connection", "read")
			return tfErr.GetDiag()
		}
		mqtts = append(mqtts, modelMap)
	}
	if err = d.Set("mqtts", mqtts); err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("Error setting mqtts: %s", err), "(Data) ibm_database_connection", "read")
		return tfErr.GetDiag()
	}

	stompSsl := []map[string]interface{}{}
	if conn.StompSsl != nil {
		modelMap, err := DataSourceIBMDatabaseConnectionConnectionURIToMap(conn.StompSsl)
		if err != nil {
			tfErr := flex.TerraformErrorf(err, err.Error(), "(Data) ibm_database_connection", "read")
			return tfErr.GetDiag()
		}
		stompSsl = append(stompSsl, modelMap)
	}
	if err = d.Set("stomp_ssl", stompSsl); err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("Error setting stomp_ssl: %s", err), "(Data) ibm_database_connection", "read")
		return tfErr.GetDiag()
	}

	grpc := []map[string]interface{}{}
	if conn.Grpc != nil {
		modelMap, err := DataSourceIBMDatabaseConnectionConnectionURIToMap(conn.Grpc)
		if err != nil {
			tfErr := flex.TerraformErrorf(err, err.Error(), "(Data) ibm_database_connection", "read")
			return tfErr.GetDiag()
		}
		grpc = append(grpc, modelMap)
	}
	if err = d.Set("grpc", grpc); err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("Error setting grpc: %s", err), "(Data) ibm_database_connection", "read")
		return tfErr.GetDiag()
	}

	mongodb := []map[string]interface{}{}
	if conn.Mongodb != nil {
		modelMap, err := DataSourceIBMDatabaseConnectionMongoDbConnectionURIToMap(conn.Mongodb)
		if err != nil {
			tfErr := flex.TerraformErrorf(err, err.Error(), "(Data) ibm_database_connection", "read")
			return tfErr.GetDiag()
		}
		mongodb = append(mongodb, modelMap)
	}
	if err = d.Set("mongodb", mongodb); err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("Error setting mongodb: %s", err), "(Data) ibm_database_connection", "read")
		return tfErr.GetDiag()
	}

	biConnector := []map[string]interface{}{}
	if conn.BiConnector != nil {
		modelMap, err := DataSourceIBMDatabaseConnectionConnectionURIToMap(conn.BiConnector)
		if err != nil {
			tfErr := flex.TerraformErrorf(err, err.Error(), "(Data) ibm_database_connection", "read")
			return tfErr.GetDiag()
		}
		biConnector = append(biConnector, modelMap)
	}
	if err = d.Set("bi_connector", biConnector); err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("Error setting bi_connector: %s", err), "(Data) ibm_database_connection", "read")
		return tfErr.GetDiag()
	}

	analytics := []map[string]interface{}{}
	if conn.Analytics != nil {
		modelMap, err := DataSourceIBMDatabaseConnectionConnectionURIToMap(conn.Analytics)
		if err != nil {
			tfErr := flex.TerraformErrorf(err, err.Error(), "(Data) ibm_database_connection", "read")
			return tfErr.GetDiag()
		}
		analytics = append(analytics, modelMap)
	}
	if err = d.Set("analytics", analytics); err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("Error setting analytics: %s", err), "(Data) ibm_database_connection", "read")
		return tfErr.GetDiag()
	}

	opsManager := []map[string]interface{}{}
	if conn.OpsManager != nil {
		modelMap, err := DataSourceIBMDatabaseConnectionConnectionURIToMap(conn.OpsManager)
		if err != nil {
			tfErr := flex.TerraformErrorf(err, err.Error(), "(Data) ibm_database_connection", "read")
			return tfErr.GetDiag()
		}
		opsManager = append(opsManager, modelMap)
	}
	if err = d.Set("ops_manager", opsManager); err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("Error setting ops_manager: %s", err), "(Data) ibm_database_connection", "read")
		return tfErr.GetDiag()
	}

	mysql := []map[string]interface{}{}
	if conn.Mysql != nil {
		modelMap, err := DataSourceIBMDatabaseConnectionMySQLConnectionURIToMap(conn.Mysql)
		if err != nil {
			tfErr := flex.TerraformErrorf(err, err.Error(), "(Data) ibm_database_connection", "read")
			return tfErr.GetDiag()
		}
		mysql = append(mysql, modelMap)
	}
	if err = d.Set("mysql", mysql); err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("Error setting mysql: %s", err), "(Data) ibm_database_connection", "read")
		return tfErr.GetDiag()
	}

	secure := []map[string]interface{}{}
	if err = d.Set("secure", secure); err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("Error setting secure: %s", err), "(Data) ibm_database_connection", "read")
		return tfErr.GetDiag()
	}

	emp := []map[string]interface{}{}
	if conn.Emp != nil {
		modelMap, err := DataSourceIBMDatabaseConnectionConnectionURIToMap(conn.Emp)
		if err != nil {
			tfErr := flex.TerraformErrorf(err, err.Error(), "(Data) ibm_database_connection", "read")
			return tfErr.GetDiag()
		}
		emp = append(emp, modelMap)
	}
	if err = d.Set("emp", emp); err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("Error setting emp: %s", err), "(Data) ibm_database_connection", "read")
		return tfErr.GetDiag()
	}

	return nil
}

// DataSourceIBMDatabaseConnectionID returns a reasonable ID for the list.
func DataSourceIBMDatabaseConnectionID(d *schema.ResourceData) string {
	return time.Now().UTC().String()
}

func DataSourceIBMDatabaseConnectionPostgreSQLConnectionURIToMap(model *clouddatabasesv5.PostgreSQLConnectionURI) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.Type != nil {
		modelMap["type"] = *model.Type
	}
	if model.Composed != nil {
		modelMap["composed"] = model.Composed
	}
	if model.Scheme != nil {
		modelMap["scheme"] = *model.Scheme
	}
	if model.Hosts != nil {
		hosts := []map[string]interface{}{}
		for _, hostsItem := range model.Hosts {
			hostsItemMap, err := DataSourceIBMDatabaseConnectionConnectionHostToMap(&hostsItem)
			if err != nil {
				return modelMap, err
			}
			hosts = append(hosts, hostsItemMap)
		}
		modelMap["hosts"] = hosts
	}
	if model.Path != nil {
		modelMap["path"] = *model.Path
	}
	if model.QueryOptions != nil {
		queryOptionsMap := make(map[string]interface{}, len(model.QueryOptions))
		for _, _ = range model.QueryOptions {
		}
		modelMap["query_options"] = flex.Flatten(queryOptionsMap)
	}
	if model.Authentication != nil {
		authenticationMap, err := DataSourceIBMDatabaseConnectionConnectionAuthenticationToMap(model.Authentication)
		if err != nil {
			return modelMap, err
		}
		modelMap["authentication"] = []map[string]interface{}{authenticationMap}
	}
	if model.Certificate != nil {
		certificateMap, err := DataSourceIBMDatabaseConnectionConnectionCertificateToMap(model.Certificate)
		if err != nil {
			return modelMap, err
		}
		modelMap["certificate"] = []map[string]interface{}{certificateMap}
	}
	if model.Ssl != nil {
		modelMap["ssl"] = *model.Ssl
	}
	if model.BrowserAccessible != nil {
		modelMap["browser_accessible"] = *model.BrowserAccessible
	}
	if model.Database != nil {
		modelMap["database"] = *model.Database
	}
	return modelMap, nil
}

func DataSourceIBMDatabaseConnectionConnectionHostToMap(model *clouddatabasesv5.ConnectionHost) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.Hostname != nil {
		modelMap["hostname"] = *model.Hostname
	}
	if model.Port != nil {
		modelMap["port"] = *model.Port
	}
	return modelMap, nil
}

func DataSourceIBMDatabaseConnectionConnectionAuthenticationToMap(model *clouddatabasesv5.ConnectionAuthentication) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.Method != nil {
		modelMap["method"] = *model.Method
	}
	if model.Username != nil {
		modelMap["username"] = *model.Username
	}
	if model.Password != nil {
		modelMap["password"] = *model.Password
	}
	return modelMap, nil
}

func DataSourceIBMDatabaseConnectionConnectionCertificateToMap(model *clouddatabasesv5.ConnectionCertificate) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.Name != nil {
		modelMap["name"] = *model.Name
	}
	if model.CertificateBase64 != nil {
		modelMap["certificate_base64"] = *model.CertificateBase64
	}
	return modelMap, nil
}

func DataSourceIBMDatabaseConnectionConnectionCliToMap(model *clouddatabasesv5.ConnectionCli) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.Type != nil {
		modelMap["type"] = *model.Type
	}
	if model.Composed != nil {
		modelMap["composed"] = model.Composed
	}
	if model.Environment != nil {
		environmentMap := make(map[string]interface{}, len(model.Environment))
		for _, _ = range model.Environment {
		}
		modelMap["environment"] = flex.Flatten(environmentMap)
	}
	if model.Bin != nil {
		modelMap["bin"] = *model.Bin
	}
	if model.Arguments != nil {
	}
	if model.Certificate != nil {
		certificateMap, err := DataSourceIBMDatabaseConnectionConnectionCertificateToMap(model.Certificate)
		if err != nil {
			return modelMap, err
		}
		modelMap["certificate"] = []map[string]interface{}{certificateMap}
	}
	return modelMap, nil
}

func DataSourceIBMDatabaseConnectionRedisConnectionURIToMap(model *clouddatabasesv5.RedisConnectionURI) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.Type != nil {
		modelMap["type"] = *model.Type
	}
	if model.Composed != nil {
		modelMap["composed"] = model.Composed
	}
	if model.Scheme != nil {
		modelMap["scheme"] = *model.Scheme
	}
	if model.Hosts != nil {
		hosts := []map[string]interface{}{}
		for _, hostsItem := range model.Hosts {
			hostsItemMap, err := DataSourceIBMDatabaseConnectionConnectionHostToMap(&hostsItem)
			if err != nil {
				return modelMap, err
			}
			hosts = append(hosts, hostsItemMap)
		}
		modelMap["hosts"] = hosts
	}
	if model.Path != nil {
		modelMap["path"] = *model.Path
	}
	if model.QueryOptions != nil {
		queryOptionsMap := make(map[string]interface{}, len(model.QueryOptions))
		for _, _ = range model.QueryOptions {
		}
		modelMap["query_options"] = flex.Flatten(queryOptionsMap)
	}
	if model.Authentication != nil {
		authenticationMap, err := DataSourceIBMDatabaseConnectionConnectionAuthenticationToMap(model.Authentication)
		if err != nil {
			return modelMap, err
		}
		modelMap["authentication"] = []map[string]interface{}{authenticationMap}
	}
	if model.Certificate != nil {
		certificateMap, err := DataSourceIBMDatabaseConnectionConnectionCertificateToMap(model.Certificate)
		if err != nil {
			return modelMap, err
		}
		modelMap["certificate"] = []map[string]interface{}{certificateMap}
	}
	if model.Ssl != nil {
		modelMap["ssl"] = *model.Ssl
	}
	if model.BrowserAccessible != nil {
		modelMap["browser_accessible"] = *model.BrowserAccessible
	}
	if model.Database != nil {
		modelMap["database"] = *model.Database
	}
	return modelMap, nil
}

func DataSourceIBMDatabaseConnectionConnectionURIToMap(model *clouddatabasesv5.ConnectionURI) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.Type != nil {
		modelMap["type"] = *model.Type
	}
	if model.Composed != nil {
		modelMap["composed"] = model.Composed
	}
	if model.Scheme != nil {
		modelMap["scheme"] = *model.Scheme
	}
	if model.Hosts != nil {
		hosts := []map[string]interface{}{}
		for _, hostsItem := range model.Hosts {
			hostsItemMap, err := DataSourceIBMDatabaseConnectionConnectionHostToMap(&hostsItem)
			if err != nil {
				return modelMap, err
			}
			hosts = append(hosts, hostsItemMap)
		}
		modelMap["hosts"] = hosts
	}
	if model.Path != nil {
		modelMap["path"] = *model.Path
	}
	if model.QueryOptions != nil {
		queryOptionsMap := make(map[string]interface{}, len(model.QueryOptions))
		for _, _ = range model.QueryOptions {
		}
		modelMap["query_options"] = flex.Flatten(queryOptionsMap)
	}
	if model.Authentication != nil {
		authenticationMap, err := DataSourceIBMDatabaseConnectionConnectionAuthenticationToMap(model.Authentication)
		if err != nil {
			return modelMap, err
		}
		modelMap["authentication"] = []map[string]interface{}{authenticationMap}
	}
	if model.Certificate != nil {
		certificateMap, err := DataSourceIBMDatabaseConnectionConnectionCertificateToMap(model.Certificate)
		if err != nil {
			return modelMap, err
		}
		modelMap["certificate"] = []map[string]interface{}{certificateMap}
	}
	if model.Ssl != nil {
		modelMap["ssl"] = *model.Ssl
	}
	if model.BrowserAccessible != nil {
		modelMap["browser_accessible"] = *model.BrowserAccessible
	}
	return modelMap, nil
}

func DataSourceIBMDatabaseConnectionMongoDbConnectionURIToMap(model *clouddatabasesv5.MongoDbConnectionURI) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.Type != nil {
		modelMap["type"] = *model.Type
	}
	if model.Composed != nil {
		modelMap["composed"] = model.Composed
	}
	if model.Scheme != nil {
		modelMap["scheme"] = *model.Scheme
	}
	if model.Hosts != nil {
		hosts := []map[string]interface{}{}
		for _, hostsItem := range model.Hosts {
			hostsItemMap, err := DataSourceIBMDatabaseConnectionConnectionHostToMap(&hostsItem)
			if err != nil {
				return modelMap, err
			}
			hosts = append(hosts, hostsItemMap)
		}
		modelMap["hosts"] = hosts
	}
	if model.Path != nil {
		modelMap["path"] = *model.Path
	}
	if model.QueryOptions != nil {
		queryOptionsMap := make(map[string]interface{}, len(model.QueryOptions))
		for _, _ = range model.QueryOptions {
		}
		modelMap["query_options"] = flex.Flatten(queryOptionsMap)
	}
	if model.Authentication != nil {
		authenticationMap, err := DataSourceIBMDatabaseConnectionConnectionAuthenticationToMap(model.Authentication)
		if err != nil {
			return modelMap, err
		}
		modelMap["authentication"] = []map[string]interface{}{authenticationMap}
	}
	if model.Certificate != nil {
		certificateMap, err := DataSourceIBMDatabaseConnectionConnectionCertificateToMap(model.Certificate)
		if err != nil {
			return modelMap, err
		}
		modelMap["certificate"] = []map[string]interface{}{certificateMap}
	}
	if model.Ssl != nil {
		modelMap["ssl"] = *model.Ssl
	}
	if model.BrowserAccessible != nil {
		modelMap["browser_accessible"] = *model.BrowserAccessible
	}
	if model.Database != nil {
		modelMap["database"] = *model.Database
	}
	if model.ReplicaSet != nil {
		modelMap["replica_set"] = *model.ReplicaSet
	}
	return modelMap, nil
}

func DataSourceIBMDatabaseConnectionMySQLConnectionURIToMap(model *clouddatabasesv5.MySQLConnectionURI) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.Type != nil {
		modelMap["type"] = *model.Type
	}
	if model.Composed != nil {
		modelMap["composed"] = model.Composed
	}
	if model.Scheme != nil {
		modelMap["scheme"] = *model.Scheme
	}
	if model.Hosts != nil {
		hosts := []map[string]interface{}{}
		for _, hostsItem := range model.Hosts {
			hostsItemMap, err := DataSourceIBMDatabaseConnectionConnectionHostToMap(&hostsItem)
			if err != nil {
				return modelMap, err
			}
			hosts = append(hosts, hostsItemMap)
		}
		modelMap["hosts"] = hosts
	}
	if model.Path != nil {
		modelMap["path"] = *model.Path
	}
	if model.QueryOptions != nil {
		queryOptionsMap := make(map[string]interface{}, len(model.QueryOptions))
		for _, _ = range model.QueryOptions {
		}
		modelMap["query_options"] = flex.Flatten(queryOptionsMap)
	}
	if model.Authentication != nil {
		authenticationMap, err := DataSourceIBMDatabaseConnectionConnectionAuthenticationToMap(model.Authentication)
		if err != nil {
			return modelMap, err
		}
		modelMap["authentication"] = []map[string]interface{}{authenticationMap}
	}
	if model.Certificate != nil {
		certificateMap, err := DataSourceIBMDatabaseConnectionConnectionCertificateToMap(model.Certificate)
		if err != nil {
			return modelMap, err
		}
		modelMap["certificate"] = []map[string]interface{}{certificateMap}
	}
	if model.Ssl != nil {
		modelMap["ssl"] = *model.Ssl
	}
	if model.BrowserAccessible != nil {
		modelMap["browser_accessible"] = *model.BrowserAccessible
	}
	if model.Database != nil {
		modelMap["database"] = *model.Database
	}
	return modelMap, nil
}
