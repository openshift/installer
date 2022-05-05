// Copyright IBM Corp. 2017, 2021 All Rights Reserved.
// Licensed under the Mozilla Public License v2.0

package functions

import (
	"fmt"
	"log"

	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/conns"
	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/flex"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func DataSourceIBMFunctionPackage() *schema.Resource {
	return &schema.Resource{

		Read: dataSourceIBMFunctionPackageRead,

		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Name of the package.",
			},
			"namespace": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Name of the namespace.",
			},
			"publish": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "Package Visibility.",
			},
			"version": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Semantic version of the package.",
			},

			"annotations": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "All annotations set on package by user and those set by the IBM Cloud Function backend/API.",
			},

			"parameters": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "All parameters set on package by user and those set by the IBM Cloud Function backend/API.",
			},

			"bind_package_name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Name of binded package.",
			},

			"package_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func dataSourceIBMFunctionPackageRead(d *schema.ResourceData, meta interface{}) error {
	functionNamespaceAPI, err := meta.(conns.ClientSession).FunctionIAMNamespaceAPI()
	if err != nil {
		return err
	}

	bxSession, err := meta.(conns.ClientSession).BluemixSession()
	if err != nil {
		return err
	}
	namespace := d.Get("namespace").(string)
	wskClient, err := conns.SetupOpenWhiskClientConfig(namespace, bxSession, functionNamespaceAPI)
	if err != nil {
		return err

	}

	packageService := wskClient.Packages
	name := d.Get("name").(string)
	pkg, _, err := packageService.Get(name)
	if err != nil {
		return fmt.Errorf("[ERROR] Error retrieving IBM Cloud Function package %s : %s", name, err)
	}

	d.SetId(pkg.Name)
	d.Set("name", pkg.Name)
	d.Set("namespace", namespace)
	d.Set("publish", pkg.Publish)
	d.Set("version", pkg.Version)
	d.Set("package_id", pkg.Name)
	annotations, err := flex.FlattenAnnotations(pkg.Annotations)
	if err != nil {
		log.Printf(
			"An error occured during reading of package (%s) annotations : %s", d.Id(), err)
	}
	d.Set("annotations", annotations)
	parameters, err := flex.FlattenParameters(pkg.Parameters)
	if err != nil {
		log.Printf(
			"An error occured during reading of package (%s) parameters : %s", d.Id(), err)
	}
	d.Set("parameters", parameters)

	if !flex.IsEmpty(*pkg.Binding) {
		d.Set("bind_package_name", fmt.Sprintf("/%s/%s", pkg.Binding.Namespace, pkg.Binding.Name))
	}
	return nil
}
