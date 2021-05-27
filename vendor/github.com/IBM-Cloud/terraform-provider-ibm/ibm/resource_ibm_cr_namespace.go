// Copyright IBM Corp. 2021 All Rights Reserved.
// Licensed under the Mozilla Public License v2.0

package ibm

import (
	"context"
	"fmt"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"

	"github.com/IBM/container-registry-go-sdk/containerregistryv1"
)

func resourceIBMCrNamespace() *schema.Resource {
	return &schema.Resource{
		Create:   resourceIBMCrNamespaceCreate,
		Read:     resourceIBMCrNamespaceRead,
		Update:   resourceIBMCrNamespaceUpdate,
		Delete:   resourceIBMCrNamespaceDelete,
		Importer: &schema.ResourceImporter{},

		Schema: map[string]*schema.Schema{
			"name": &schema.Schema{
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: InvokeValidator("ibm_cr_namespace", "name"),
				Description:  "The name of the namespace.",
			},
			"resource_group_id": &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				ForceNew:    true,
				Description: "The ID of the resource group that the namespace will be created within.",
			},
			"tags": &schema.Schema{
				Type:        schema.TypeSet,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Set:         schema.HashString,
				Optional:    true,
				Description: "List of tags",
			},
			"account": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The IBM Cloud account that owns the namespace.",
			},
			"created_date": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "When the namespace was created.",
			},
			"crn": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "If the namespace has been assigned to a resource group, this is the IBM Cloud CRN representing the namespace.",
			},
			"resource_created_date": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "When the namespace was assigned to a resource group.",
			},
			"updated_date": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "When the namespace was last updated.",
			},
			// HAND-ADDED DEPRECATED FIELDS, TO BE DELETED IN FUTURE
			"created_on": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "When the namespace was created.",
				Deprecated:  "This field is deprecated",
			},
			"updated_on": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "When the namespace was last updated.",
				Deprecated:  "This field is deprecated",
			},
		},
	}
}

func resourceIBMCrNamespaceValidator() *ResourceValidator {
	validateSchema := make([]ValidateSchema, 1)
	validateSchema = append(validateSchema,
		ValidateSchema{
			Identifier:                 "name",
			ValidateFunctionIdentifier: ValidateRegexpLen,
			Type:                       TypeString,
			Required:                   true,
			Regexp:                     `^[a-z0-9]+[a-z0-9_-]+[a-z0-9]+$`,
			MinValueLength:             4,
			MaxValueLength:             30,
		},
	)

	resourceValidator := ResourceValidator{ResourceName: "ibm_cr_namespace", Schema: validateSchema}
	return &resourceValidator
}

func resourceIBMCrNamespaceCreate(d *schema.ResourceData, meta interface{}) error {
	containerRegistryClient, err := meta.(ClientSession).ContainerRegistryV1()
	if err != nil {
		return err
	}

	createNamespaceOptions := &containerregistryv1.CreateNamespaceOptions{}

	createNamespaceOptions.SetName(d.Get("name").(string))
	if _, ok := d.GetOk("resource_group_id"); ok {
		createNamespaceOptions.SetXAuthResourceGroup(d.Get("resource_group_id").(string))
	} else {
		defaultRg, err := defaultResourceGroup(meta)
		if err != nil {
			return err
		}
		createNamespaceOptions.SetXAuthResourceGroup(defaultRg)
	}

	namespace, response, err := containerRegistryClient.CreateNamespaceWithContext(context.TODO(), createNamespaceOptions)
	if err != nil {
		log.Printf("[DEBUG] CreateNamespaceWithContext failed %s\n%s", err, response)
		return err
	}

	d.SetId(*namespace.Namespace)

	return resourceIBMCrNamespaceRead(d, meta)
}

func resourceIBMCrNamespaceRead(d *schema.ResourceData, meta interface{}) error {
	containerRegistryClient, err := meta.(ClientSession).ContainerRegistryV1()
	if err != nil {
		return err
	}

	listNamespaceDetailsOptions := &containerregistryv1.ListNamespaceDetailsOptions{}

	namespaceDetailsList, response, err := containerRegistryClient.ListNamespaceDetailsWithContext(context.TODO(), listNamespaceDetailsOptions)

	if err != nil {
		if response != nil && response.StatusCode == 404 {
			d.SetId("")
			return nil
		}
		log.Printf("[DEBUG] ListNamespaceDetailsWithContext failed %s\n%s", err, response)
		return err
	}

	var namespaceDetails containerregistryv1.NamespaceDetails
	for _, namespaceDetails = range namespaceDetailsList {
		if *namespaceDetails.Name == d.Id() {
			break
		}
	}
	if namespaceDetails.Name == nil || *namespaceDetails.Name != d.Id() {
		d.SetId("")
		return nil
	}

	if err = d.Set("name", namespaceDetails.Name); err != nil {
		return fmt.Errorf("Error setting name: %s", err)
	}
	if err = d.Set("resource_group_id", namespaceDetails.ResourceGroup); err != nil {
		return fmt.Errorf("Error setting resource_group_id: %s", err)
	}
	if err = d.Set("account", namespaceDetails.Account); err != nil {
		return fmt.Errorf("Error setting account: %s", err)
	}
	if err = d.Set("created_date", namespaceDetails.CreatedDate); err != nil {
		return fmt.Errorf("Error setting created_date: %s", err)
	}
	if err = d.Set("crn", namespaceDetails.CRN); err != nil {
		return fmt.Errorf("Error setting crn: %s", err)
	}
	if err = d.Set("resource_created_date", namespaceDetails.ResourceCreatedDate); err != nil {
		return fmt.Errorf("Error setting resource_created_date: %s", err)
	}
	if err = d.Set("updated_date", namespaceDetails.UpdatedDate); err != nil {
		return fmt.Errorf("Error setting updated_date: %s", err)
	}
	// HAND-ADDED DEPRECATED FIELDS, TO BE DELETED IN FUTURE
	if err = d.Set("updated_on", namespaceDetails.UpdatedDate); err != nil {
		return fmt.Errorf("Error setting updated_date: %s", err)
	}
	if err = d.Set("created_on", namespaceDetails.CreatedDate); err != nil {
		return fmt.Errorf("Error setting created_date: %s", err)
	}

	return nil
}

// Dummy update method just for local tags
func resourceIBMCrNamespaceUpdate(d *schema.ResourceData, meta interface{}) error {
	return resourceIBMCrNamespaceRead(d, meta)
}

func resourceIBMCrNamespaceDelete(d *schema.ResourceData, meta interface{}) error {
	containerRegistryClient, err := meta.(ClientSession).ContainerRegistryV1()
	if err != nil {
		return err
	}

	deleteNamespaceOptions := &containerregistryv1.DeleteNamespaceOptions{}

	deleteNamespaceOptions.SetName(d.Id())

	response, err := containerRegistryClient.DeleteNamespaceWithContext(context.TODO(), deleteNamespaceOptions)
	if err != nil {
		log.Printf("[DEBUG] DeleteNamespaceWithContext failed %s\n%s", err, response)
		return err
	}

	d.SetId("")

	return nil
}
