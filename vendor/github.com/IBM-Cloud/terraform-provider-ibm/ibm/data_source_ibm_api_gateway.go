// Copyright IBM Corp. 2017, 2021 All Rights Reserved.
// Licensed under the Mozilla Public License v2.0

package ibm

import (
	"encoding/json"
	"fmt"
	"strings"

	apigatewaysdk "github.com/IBM/apigateway-go-sdk"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func dataSourceIBMApiGateway() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceIBMApiGatewayRead,
		Schema: map[string]*schema.Schema{
			"service_instance_crn": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Api Gateway Service Instance Crn",
			},
			"endpoints": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "List of all endpoints of an instance",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"endpoint_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Endpoint ID",
						},
						"name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Endpoint name",
						},
						"routes": {
							Type:        schema.TypeList,
							Elem:        &schema.Schema{Type: schema.TypeString},
							Computed:    true,
							Description: "Invokable routes for an endpoint",
						},
						"managed": {
							Type:        schema.TypeBool,
							Computed:    true,
							Description: "Managed indicates if endpoint is online or offline.",
						},
						"shared": {
							Type:        schema.TypeBool,
							Computed:    true,
							Description: "The Shared status of an endpoint",
						},
						"base_path": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: " Base path of an endpoint",
						},
						"managed_url": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Managed url for an endpoint",
						},
						"alias_url": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Alias Url for an endpoint",
						},
						"open_api_doc": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "API document that represents endpoint",
						},
						"subscriptions": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "List of all subscription of an endpoint",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"client_id": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Subscription Id",
									},
									"name": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Subscription name",
									},
									"type": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Subscription type",
									},
									"secret_provided": {
										Type:        schema.TypeBool,
										Computed:    true,
										Description: "Indicates if client secret is provided to subscription or not",
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

func dataSourceIBMApiGatewayRead(d *schema.ResourceData, meta interface{}) error {
	sess, err := meta.(ClientSession).BluemixSession()
	if err != nil {
		return err
	}
	endpointservice, err := meta.(ClientSession).APIGateway()
	if err != nil {
		return err
	}
	payload := &apigatewaysdk.GetAllEndpointsOptions{}
	oauthtoken := sess.Config.IAMAccessToken
	oauthtoken = strings.Replace(oauthtoken, "Bearer ", "", -1)
	serviceInstanceCrn := d.Get("service_instance_crn").(string)
	payload.Authorization = &oauthtoken
	payload.ServiceInstanceCrn = &serviceInstanceCrn
	allendpoints, response, err := endpointservice.GetAllEndpoints(payload)
	if err != nil {
		return fmt.Errorf("Error Getting All Endpoint: %s,%s", err, response)
	}
	endpointsMap := make([]map[string]interface{}, 0, len(*allendpoints))

	for _, endpoint := range *allendpoints {
		ArtifactID := endpoint.ArtifactID

		swaggerPayload := &apigatewaysdk.GetEndpointSwaggerOptions{}
		swaggerPayload.Authorization = &oauthtoken
		swaggerPayload.ID = ArtifactID
		swaggerPayload.ServiceInstanceCrn = &serviceInstanceCrn

		swagger, err := endpointservice.GetEndpointSwagger(swaggerPayload)
		if err != nil {
			return fmt.Errorf("Error Getting All Endpoint: %s,%s", err, swagger)
		}
		doc := swagger.Result
		str, err := json.Marshal(doc)
		if err != nil {
			fmt.Printf("error while json Marshal: %v", err)
		}
		swagger_document := string(str)
		SubscriptionPayload := &apigatewaysdk.GetAllSubscriptionsOptions{}
		SubscriptionPayload.ArtifactID = ArtifactID
		SubscriptionPayload.Authorization = &oauthtoken
		if v, ok := d.GetOk("type"); ok && v != nil {
			Type := v.(string)
			if Type == "internal" {
				Type = "bluemix"
			}
			SubscriptionPayload.Type = &Type
		}
		allsubscriptions, response, err := endpointservice.GetAllSubscriptions(SubscriptionPayload)
		if err != nil {
			return fmt.Errorf("Error Getting All Endpoint: %s %s", err, response)
		}
		subscriptionMap := make([]map[string]interface{}, 0, len(*allsubscriptions))
		for _, subscription := range *allsubscriptions {
			allsubscription := make(map[string]interface{})
			allsubscription["name"] = *subscription.Name
			allsubscription["client_id"] = subscription.ClientID
			if *subscription.Type == "bluemix" {
				*subscription.Type = "internal"
			}
			allsubscription["type"] = subscription.Type
			allsubscription["secret_provided"] = subscription.SecretProvided
			subscriptionMap = append(subscriptionMap, allsubscription)
		}
		result := make(map[string]interface{})
		result["endpoint_id"] = *endpoint.ArtifactID
		result["name"] = *endpoint.Name
		result["managed"] = endpoint.Managed
		result["shared"] = endpoint.Shared
		result["routes"] = endpoint.Routes
		result["managed_url"] = *endpoint.ManagedURL
		result["base_path"] = endpoint.BasePath
		result["alias_url"] = endpoint.AliasURL
		result["open_api_doc"] = swagger_document
		result["subscriptions"] = subscriptionMap
		endpointsMap = append(endpointsMap, result)
	}
	d.SetId(serviceInstanceCrn)
	d.Set("service_instance_crn", serviceInstanceCrn)
	d.Set("endpoints", endpointsMap)

	return nil
}
