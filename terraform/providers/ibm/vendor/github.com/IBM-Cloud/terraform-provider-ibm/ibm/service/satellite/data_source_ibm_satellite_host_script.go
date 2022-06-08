// Copyright IBM Corp. 2017, 2021 All Rights Reserved.
// Licensed under the Mozilla Public License v2.0

package satellite

import (
	"fmt"
	"io/ioutil"
	"log"
	"path/filepath"
	"strings"
	"time"

	"github.com/IBM-Cloud/container-services-go-sdk/kubernetesserviceapiv1"
	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/conns"
	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/flex"
	"github.com/IBM/go-sdk-core/v5/core"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	homedir "github.com/mitchellh/go-homedir"
)

func DataSourceIBMSatelliteAttachHostScript() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceIBMSatelliteAttachHostScriptRead,

		Schema: map[string]*schema.Schema{
			"location": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "A unique name for the new Satellite location",
			},
			"description": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "A unique name for the new Satellite location",
			},
			"labels": {
				Type:        schema.TypeSet,
				Optional:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Set:         schema.HashString,
				Description: "List of labels for the attach host",
			},
			"host_provider": {
				Type:         schema.TypeString,
				Optional:     true,
				ExactlyOneOf: []string{"host_provider", "custom_script"},
			},
			"script_dir": {
				Description: "The directory where the satellite attach host script to be downloaded. Default is home directory",
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
			},
			"script_path": {
				Description: "The absolute path to the generated host script file",
				Type:        schema.TypeString,
				Computed:    true,
			},
			"host_script": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Attach host script content",
			},
			"custom_script": {
				Description:  "The custom script that has to be appended to generated host script file",
				Type:         schema.TypeString,
				Optional:     true,
				ExactlyOneOf: []string{"host_provider", "custom_script"},
			},
		},
	}
}

func dataSourceIBMSatelliteAttachHostScriptRead(d *schema.ResourceData, meta interface{}) error {
	var scriptDir string
	location := d.Get("location").(string)
	hostProvider := d.Get("host_provider").(string)

	if _, ok := d.GetOk("script_dir"); ok {
		scriptDir = d.Get("script_dir").(string)
	}

	satClient, err := meta.(conns.ClientSession).SatelliteClientSession()
	if err != nil {
		return err
	}

	var locData *kubernetesserviceapiv1.MultishiftGetController
	var response *core.DetailedResponse
	getSatLocOptions := &kubernetesserviceapiv1.GetSatelliteLocationOptions{
		Controller: &location,
	}

	err = resource.Retry(1*time.Minute, func() *resource.RetryError {
		locData, response, err = satClient.GetSatelliteLocation(getSatLocOptions)
		if err != nil || locData == nil {
			if response != nil && response.StatusCode == 404 {
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		return nil
	})

	if conns.IsResourceTimeoutError(err) {
		locData, response, err = satClient.GetSatelliteLocation(getSatLocOptions)
	}
	if err != nil || locData == nil {
		return fmt.Errorf("[ERROR] Error getting Satellite location (%s): %s\n%s", location, err, response)
	}

	// script labels
	labels := make(map[string]string)
	if v, ok := d.GetOk("labels"); ok {
		l := v.(*schema.Set)
		labels = flex.FlattenHostLabels(l.List())
		d.Set("labels", l)
	}

	if len(scriptDir) == 0 {
		scriptDir, err = homedir.Dir()
		if err != nil {
			return fmt.Errorf("[ERROR] Error fetching homedir: %s", err)
		}
	}
	scriptDir, _ = filepath.Abs(scriptDir)
	scriptPath := filepath.Join(scriptDir, "addHost.sh")

	//Generate script
	createRegOptions := &kubernetesserviceapiv1.AttachSatelliteHostOptions{}
	createRegOptions.Controller = locData.ID
	createRegOptions.Labels = labels

	resp, err := satClient.AttachSatelliteHost(createRegOptions)
	if err != nil {
		return fmt.Errorf("[ERROR] Error Generating Satellite Registration Script: %s\n%s", err, resp)
	}

	lines := strings.Split(string(resp), "\n")
	for i, line := range lines {
		if strings.Contains(line, "API_URL=") {
			i = i + 1
			if script, ok := d.GetOk("custom_script"); ok {
				lines[i] = script.(string)
			} else {
				if strings.ToLower(hostProvider) == "aws" {
					lines[i] = "yum update -y\nyum-config-manager --enable '*'\nyum repolist all\nyum install container-selinux -y"
				} else if strings.ToLower(hostProvider) == "ibm" {
					lines[i] = `subscription-manager refresh
subscription-manager repos --enable rhel-server-rhscl-7-rpms
subscription-manager repos --enable rhel-7-server-optional-rpms
subscription-manager repos --enable rhel-7-server-rh-common-rpms
subscription-manager repos --enable rhel-7-server-supplementary-rpms
subscription-manager repos --enable rhel-7-server-extras-rpms`
				} else if strings.ToLower(hostProvider) == "azure" {
					lines[i] = fmt.Sprintf(`yum update --disablerepo=* --enablerepo="*microsoft*" -y
yum-config-manager --enable '*'
yum repolist all
yum install container-selinux -y
					`)
				} else if strings.ToLower(hostProvider) == "google" {
					lines[i] = fmt.Sprintf(`yum update --disablerepo=* --enablerepo="*" -y
yum repolist all
yum install container-selinux -y
yum install subscription-manager -y
	`)
				} else {
					lines[i] = "subscription-manager refresh\nyum update -y\n"
				}
			}

		}
	}

	scriptContent := strings.Join(lines, "\n")
	err = ioutil.WriteFile(scriptPath, []byte(scriptContent), 0644)
	if err != nil {
		return fmt.Errorf("[ERROR] Error Creating Satellite Attach Host Script: %s", err)
	}

	d.Set("location", location)
	d.Set("host_script", scriptContent)
	d.Set("host_provider", hostProvider)
	d.Set("script_dir", scriptDir)
	d.Set("script_path", scriptPath)
	d.SetId(*locData.ID)

	log.Printf("[INFO] Generated satellite location script : %s", *locData.Name)

	return nil
}
