// Copyright IBM Corp. 2017, 2021 All Rights Reserved.
// Licensed under the Mozilla Public License v2.0

package ibm

import (
	"fmt"
	"io/ioutil"
	"path"
	"strings"

	apigatewaysdk "github.com/IBM/apigateway-go-sdk"
	"github.com/ghodss/yaml"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceIBMApiGatewayEndPoint() *schema.Resource {

	return &schema.Resource{
		Create:   resourceIBMApiGatewayEndPointCreate,
		Read:     resourceIBMApiGatewayEndPointGet,
		Update:   resourceIBMApiGatewayEndPointUpdate,
		Delete:   resourceIBMApiGatewayEndPointDelete,
		Importer: &schema.ResourceImporter{},
		Exists:   resourceIBMApiGatewayEndPointExists,
		Schema: map[string]*schema.Schema{
			"service_instance_crn": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Api Gateway Service Instance Crn",
			},
			"open_api_doc_name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Json File path",
			},
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Endpoint name",
			},
			"routes": {
				Type:        schema.TypeList,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Optional:    true,
				Description: "Invokable routes for an endpoint",
			},
			"managed": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
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
			"provider_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "user-defined",
				Description: " Provider ID of an endpoint allowable values user-defined and whisk",
			},
			"endpoint_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Endpoint ID",
			},
			"type": {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "unshare",
				Description: "Action type of Endpoint ALoowable values are share, unshare, manage, unmanage",
			},
		},
	}
}

func resourceIBMApiGatewayEndPointCreate(d *schema.ResourceData, meta interface{}) error {
	sess, err := meta.(ClientSession).BluemixSession()
	if err != nil {
		return err
	}
	endpointservice, err := meta.(ClientSession).APIGateway()
	if err != nil {
		return err
	}
	payload := &apigatewaysdk.CreateEndpointOptions{}

	oauthtoken := sess.Config.IAMAccessToken
	oauthtoken = strings.Replace(oauthtoken, "Bearer ", "", -1)
	payload.Authorization = &oauthtoken

	serviceInstanceCrn := d.Get("service_instance_crn").(string)
	payload.ServiceInstanceCrn = &serviceInstanceCrn
	payload.ParentCrn = &serviceInstanceCrn

	var name string
	if v, ok := d.GetOk("name"); ok && v != nil {
		name = v.(string)
		payload.Name = &name
	}

	openAPIDocName := d.Get("open_api_doc_name").(string)
	var document []byte
	// set to true as placeholder for logic control swtich
	if true {
		ext := path.Ext(openAPIDocName)
		if strings.ToLower(ext) == ".json" {
			data, err := ioutil.ReadFile(openAPIDocName)
			if err != nil {
				fmt.Println("Error uploading file", err)
				return err
			}
			document = data
		} else if strings.ToLower(ext) == ".yaml" || strings.ToLower(ext) == ".yml" {
			data, err := ioutil.ReadFile(openAPIDocName)
			if err != nil {
				fmt.Println("Error uploading file", err)
				return err
			}
			y2j, yErr := yaml.YAMLToJSON(data)
			if yErr != nil {
				fmt.Println("Error parsing yaml file", err)
				return err
			}
			document = y2j
		} else {
			return fmt.Errorf("File extension type must be json or yaml")

		}
	}
	payload.OpenApiDoc = string(document)

	var managed bool
	if m, ok := d.GetOk("managed"); ok && m != nil {
		managed = m.(bool)
		payload.Managed = &managed
	}
	var routes []string
	if r, ok := d.GetOk("routes"); ok && r != nil {
		routes = r.([]string)
		payload.Routes = routes
	}

	result, response, err := endpointservice.CreateEndpoint(payload)
	if err != nil {
		return fmt.Errorf("Error creating Endpoint: %s,%s", err, response)
	}

	d.SetId(fmt.Sprintf("%s//%s", *result.ServiceInstanceCrn, *result.ArtifactID))

	return resourceIBMApiGatewayEndPointGet(d, meta)
}

func resourceIBMApiGatewayEndPointGet(d *schema.ResourceData, meta interface{}) error {
	sess, err := meta.(ClientSession).BluemixSession()
	if err != nil {
		return err
	}
	endpointservice, err := meta.(ClientSession).APIGateway()
	if err != nil {
		return err
	}

	parts := d.Id()
	partslist := strings.Split(parts, "//")

	serviceInstanceCrn := partslist[0]
	apiID := partslist[1]

	oauthtoken := sess.Config.IAMAccessToken
	oauthtoken = strings.Replace(oauthtoken, "Bearer ", "", -1)

	payload := apigatewaysdk.GetEndpointOptions{
		ServiceInstanceCrn: &serviceInstanceCrn,
		ID:                 &apiID,
		Authorization:      &oauthtoken,
	}
	result, response, err := endpointservice.GetEndpoint(&payload)
	if err != nil {
		if response != nil && response.StatusCode == 404 {
			d.SetId("")
			return nil
		}
		return fmt.Errorf("Error Getting Endpoint: %s\n%s", err, response)
	}
	d.Set("service_instance_crn", serviceInstanceCrn)
	d.Set("endpoint_id", apiID)
	if result.Routes != nil {
		d.Set("routes", result.Routes)
	}
	if result.Name != nil {
		d.Set("name", result.Name)
	}
	if result.Managed != nil {
		d.Set("managed", result.Managed)
	}
	d.Set("provider_id", result.ProviderID)
	d.Set("shared", result.Shared)
	d.Set("base_path", result.BasePath)
	return nil
}

func resourceIBMApiGatewayEndPointUpdate(d *schema.ResourceData, meta interface{}) error {
	sess, err := meta.(ClientSession).BluemixSession()
	if err != nil {
		return err
	}
	endpointservice, err := meta.(ClientSession).APIGateway()
	if err != nil {
		return err
	}
	//payload for updating endpoint
	payload := &apigatewaysdk.UpdateEndpointOptions{}

	parts := d.Id()
	partslist := strings.Split(parts, "//")
	serviceInstanceCrn := partslist[0]
	apiID := partslist[1]

	oauthtoken := sess.Config.IAMAccessToken
	oauthtoken = strings.Replace(oauthtoken, "Bearer ", "", -1)
	payload.Authorization = &oauthtoken

	payload.ID = &apiID
	payload.NewArtifactID = &apiID

	payload.ServiceInstanceCrn = &serviceInstanceCrn
	payload.NewParentCrn = &serviceInstanceCrn
	payload.NewServiceInstanceCrn = &serviceInstanceCrn

	name := d.Get("name").(string)
	payload.NewName = &name

	managed := d.Get("managed").(bool)

	openAPIDocName := d.Get("open_api_doc_name").(string)
	var document []byte
	// set to true as placeholder for logic control swtich
	if true {
		ext := path.Ext(openAPIDocName)
		if strings.ToLower(ext) == ".json" {
			data, err := ioutil.ReadFile(openAPIDocName)
			if err != nil {
				fmt.Println("Error uploading file", err)
				return err
			}
			document = data
		} else if strings.ToLower(ext) == ".yaml" || strings.ToLower(ext) == ".yml" {
			data, err := ioutil.ReadFile(openAPIDocName)
			if err != nil {
				fmt.Println("Error uploading file", err)
				return err
			}
			y2j, yErr := yaml.YAMLToJSON(data)
			if yErr != nil {
				fmt.Println("Error parsing yaml file", err)
				return err
			}
			document = y2j
		} else {
			return fmt.Errorf("File extension type must be json or yaml")

		}
	}
	payload.NewOpenApiDoc = string(document)

	//payload for updating action of endpoint
	actionPayload := &apigatewaysdk.EndpointActionsOptions{}

	actionPayload.ServiceInstanceCrn = &serviceInstanceCrn
	actionPayload.Authorization = &oauthtoken

	actionPayload.ID = &apiID
	providerID := d.Get("provider_id").(string)
	actionPayload.ProviderID = &providerID

	actionType := d.Get("type").(string)
	actionPayload.Type = &actionType

	update := false

	if d.HasChange("name") {
		name := d.Get("name").(string)
		payload.NewName = &name
		update = true
	}
	if d.HasChange("provider_id") {
		providerID := d.Get("provider_id").(string)
		actionPayload.ProviderID = &providerID
	}
	if d.HasChange("type") {
		actionType := d.Get("type").(string)

		if managed == false && actionType == "share" {
			return fmt.Errorf("Endpoint %s not managed", apiID)
		}
		actionPayload.Type = &actionType

		_, response, err := endpointservice.EndpointActions(actionPayload)
		if err != nil {
			return fmt.Errorf("Error updating Endpoint Action: %s,%s", err, response)
		}
	}

	if d.HasChange("open_api_doc_name") {
		openAPIDocName := d.Get("open_api_doc_name").(string)
		var document []byte
		// set to true as placeholder for logic control swtich
		if true {
			ext := path.Ext(openAPIDocName)
			if strings.ToLower(ext) == ".json" {
				data, err := ioutil.ReadFile(openAPIDocName)
				if err != nil {
					fmt.Println("Error uploading file", err)
					return err
				}
				document = data
			} else if strings.ToLower(ext) == ".yaml" || strings.ToLower(ext) == ".yml" {
				data, err := ioutil.ReadFile(openAPIDocName)
				if err != nil {
					fmt.Println("Error uploading file", err)
					return err
				}
				y2j, yErr := yaml.YAMLToJSON(data)
				if yErr != nil {
					fmt.Println("Error parsing yaml file", err)
					return err
				}
				document = y2j
			} else {
				return fmt.Errorf("File extension type must be json or yaml")

			}
		}
		payload.NewOpenApiDoc = string(document)
		update = true
	}
	if d.HasChange("routes") {
		routes := d.Get("routes").([]string)
		payload.NewRoutes = routes
		update = true
	}
	if update {
		_, response, err := endpointservice.UpdateEndpoint(payload)
		if err != nil {
			return fmt.Errorf("Error updating Endpoint: %s,%s", err, response)
		}
	}
	return resourceIBMApiGatewayEndPointGet(d, meta)
}
func resourceIBMApiGatewayEndPointDelete(d *schema.ResourceData, meta interface{}) error {
	sess, err := meta.(ClientSession).BluemixSession()
	if err != nil {
		return err
	}
	endpointservice, err := meta.(ClientSession).APIGateway()
	if err != nil {
		return err
	}

	parts := d.Id()
	partslist := strings.Split(parts, "//")
	serviceInstanceCrn := partslist[0]
	apiID := partslist[1]

	oauthtoken := sess.Config.IAMAccessToken
	oauthtoken = strings.Replace(oauthtoken, "Bearer ", "", -1)

	payload := apigatewaysdk.DeleteEndpointOptions{
		ServiceInstanceCrn: &serviceInstanceCrn,
		ID:                 &apiID,
		Authorization:      &oauthtoken,
	}

	response, err := endpointservice.DeleteEndpoint(&payload)

	if err != nil {
		if response != nil && response.StatusCode == 404 {
			return nil
		}
		return fmt.Errorf("Error deleting Endpoint: %s\n%s", err, response)
	}
	d.SetId("")

	return nil
}

func resourceIBMApiGatewayEndPointExists(d *schema.ResourceData, meta interface{}) (bool, error) {
	sess, err := meta.(ClientSession).BluemixSession()
	if err != nil {
		return false, err
	}
	endpointservice, err := meta.(ClientSession).APIGateway()
	if err != nil {
		return false, err
	}

	parts := d.Id()
	partslist := strings.Split(parts, "//")
	serviceInstanceCrn := partslist[0]
	apiID := partslist[1]

	oauthtoken := sess.Config.IAMAccessToken
	oauthtoken = strings.Replace(oauthtoken, "Bearer ", "", -1)

	payload := apigatewaysdk.GetEndpointOptions{
		ServiceInstanceCrn: &serviceInstanceCrn,
		ID:                 &apiID,
		Authorization:      &oauthtoken,
	}
	_, response, err := endpointservice.GetEndpoint(&payload)
	if err != nil {
		if response != nil && response.StatusCode == 404 {
			return false, nil
		}
		return false, err
	}
	return true, nil
}
