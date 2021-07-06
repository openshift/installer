// Copyright IBM Corp. 2017, 2021 All Rights Reserved.
// Licensed under the Mozilla Public License v2.0

package ibm

import (
	"fmt"
	"log"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func dataSourceIBMFunctionAction() *schema.Resource {
	return &schema.Resource{

		Read: dataSourceIBMFunctionActionRead,

		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Name of action.",
			},
			"namespace": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Name of the namespace.",
			},
			"limits": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"timeout": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "The timeout LIMIT in milliseconds after which the action is terminated.",
						},
						"memory": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "The maximum memory LIMIT in MB for the action (default 256.",
						},
						"log_size": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "The maximum log size LIMIT in MB for the action.",
						},
					},
				},
			},
			"exec": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"image": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Container image name when kind is 'blackbox'.",
						},
						"init": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Optional zipfile reference.",
						},
						"code": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The code to execute when kind is not 'blackbox'.",
						},
						"kind": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The type of action. Possible values:php:7.3, nodejs:8, swift:3, nodejs, blackbox, java, sequence, nodejs:10, python:3, python, python:2, swift, swift:4.2.",
						},
						"main": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The name of the action entry point (function or fully-qualified method name when applicable)",
						},
						"components": {
							Type:        schema.TypeList,
							Computed:    true,
							Elem:        &schema.Schema{Type: schema.TypeString},
							Description: "The List of fully qualified action",
						},
					},
				},
			},
			"publish": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "Action visibilty.",
			},
			"version": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Semantic version of the item.",
			},
			"annotations": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "All annotations set on action by user and those set by the IBM Cloud Function backend/API.",
			},
			"parameters": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "All paramters set on action by user and those set by the IBM Cloud Function backend/API.",
			},
			"action_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"target_endpoint_url": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Action target endpoint URL.",
			},
		},
	}
}

func dataSourceIBMFunctionActionRead(d *schema.ResourceData, meta interface{}) error {
	functionNamespaceAPI, err := meta.(ClientSession).FunctionIAMNamespaceAPI()
	if err != nil {
		return err
	}

	bxSession, err := meta.(ClientSession).BluemixSession()
	if err != nil {
		return err
	}
	namespace := d.Get("namespace").(string)
	wskClient, err := setupOpenWhiskClientConfig(namespace, bxSession, functionNamespaceAPI)
	if err != nil {
		return err

	}

	actionService := wskClient.Actions
	name := d.Get("name").(string)

	action, _, err := actionService.Get(name, true)
	if err != nil {
		return fmt.Errorf("Error retrieving IBM Cloud Function Action %s : %s", name, err)
	}

	temp := strings.Split(action.Namespace, "/")
	pkgName := ""
	if len(temp) == 2 {
		pkgName = temp[1]
		d.SetId(fmt.Sprintf("%s/%s", pkgName, action.Name))
		d.Set("name", fmt.Sprintf("%s/%s", pkgName, action.Name))
	} else {
		d.SetId(action.Name)
		d.Set("name", action.Name)
	}

	d.Set("namespace", namespace)
	d.Set("limits", flattenLimits(action.Limits))
	d.Set("exec", flattenExec(action.Exec, d))
	d.Set("publish", action.Publish)
	d.Set("version", action.Version)
	d.Set("action_id", action.Name)
	annotations, err := flattenAnnotations(action.Annotations)
	if err != nil {
		log.Printf(
			"An error occured during reading of action (%s) annotations : %s", d.Id(), err)
	}
	d.Set("annotations", annotations)
	parameters, err := flattenParameters(action.Parameters)
	if err != nil {
		log.Printf(
			"An error occured during reading of action (%s) parameters : %s", d.Id(), err)
	}
	d.Set("parameters", parameters)

	targetURL, err := action.ActionURL(wskClient.Config.Host, "/api", wskClient.Config.Version, pkgName)
	if err != nil {
		log.Printf(
			"An error occured during reading of action (%s) targetURL : %s", d.Id(), err)
	}
	d.Set("target_endpoint_url", targetURL)

	return nil
}
