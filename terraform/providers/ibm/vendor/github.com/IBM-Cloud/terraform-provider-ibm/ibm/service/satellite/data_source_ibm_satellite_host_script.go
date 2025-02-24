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
			"coreos_host": {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: "If true, returns a CoreOS ignition file for the host. Otherwise, returns a RHEL attach script",
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
			"host_link_agent_endpoint": {
				Description: "The satellite link agent endpoint, required for reduced firewall attach script",
				Type:        schema.TypeString,
				Optional:    true,
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
		labels = flex.FlattenKeyValues(l.List())
		d.Set("labels", l)
	}

	if len(scriptDir) == 0 {
		scriptDir, err = homedir.Dir()
		if err != nil {
			return fmt.Errorf("[ERROR] Error fetching homedir: %s", err)
		}
	}
	scriptDir, _ = filepath.Abs(scriptDir)
	var scriptPath string

	//Generate script
	createRegOptions := &kubernetesserviceapiv1.AttachSatelliteHostOptions{}
	createRegOptions.Controller = locData.ID
	createRegOptions.Labels = labels

	//check to see if host attach is CoreOS or RHEL
	var host_os string
	var coreos_enabled bool
	if _, ok := d.GetOk("coreos_host"); ok {
		coreos_enabled = d.Get("coreos_host").(bool)
		if coreos_enabled {
			host_os = "RHCOS"
			createRegOptions.OperatingSystem = &host_os
			scriptPath = filepath.Join(scriptDir, "addHost.ign")
		}
	} else {
		coreos_enabled = false
		host_os = "RHEL"
		createRegOptions.OperatingSystem = &host_os
		scriptPath = filepath.Join(scriptDir, "addHost.sh")
	}

	// If the user supplied link agent endpoint, use reduced firewall attach script
	if hlae, ok := d.GetOk("host_link_agent_endpoint"); ok {
		host_link_agent_endpoint := hlae.(string)
		createRegOptions.HostLinkAgentEndpoint = &host_link_agent_endpoint
	}

	resp, err := satClient.AttachSatelliteHost(createRegOptions)
	if err != nil {
		return fmt.Errorf("[ERROR] Error Generating Satellite Registration Script: %s\n%s", err, resp)
	}

	scriptContent := string(resp)

	//if this is a RHEL host, find insert point for custom code
	if !coreos_enabled {
		lines := strings.Split(scriptContent, "\n")
		var index int
		for i, line := range lines {
			if strings.Contains(line, `export OPERATING_SYSTEM`) {
				index = i
				break
			}
		}

		var insertionText string

		switch {
		case strings.ToLower(hostProvider) == "aws":
			insertionText = `
yum-config-manager --enable '*'
yum install container-selinux -y
`
		case strings.ToLower(hostProvider) == "ibm":
			insertionText = `
subscription-manager refresh
if [[ "${OPERATING_SYSTEM}" == "RHEL7" ]]; then
	subscription-manager repos --enable rhel-server-rhscl-7-rpms
	subscription-manager repos --enable rhel-7-server-optional-rpms
	subscription-manager repos --enable rhel-7-server-rh-common-rpms
	subscription-manager repos --enable rhel-7-server-supplementary-rpms
	subscription-manager repos --enable rhel-7-server-extras-rpms
elif [[ "${OPERATING_SYSTEM}" == "RHEL8" ]]; then
	subscription-manager release --set=8
	subscription-manager repos --disable='*eus*'
	subscription-manager repos --enable rhel-8-for-x86_64-baseos-rpms 
	subscription-manager repos --enable rhel-8-for-x86_64-appstream-rpms;
fi
yum install container-selinux -y
`
		case strings.ToLower(hostProvider) == "azure":
			insertionText = `
#if [[ "${OPERATING_SYSTEM}" == "RHEL8" ]]; then
#	update-alternatives --install /usr/bin/python3 python3 /usr/bin/python3.8 1
#	update-alternatives --set python3 /usr/bin/python3.8
#fi
yum install container-selinux -y
`
		case strings.ToLower(hostProvider) == "google":
			insertionText = `
#if [[ "${OPERATING_SYSTEM}" == "RHEL8" ]]; then
#	update-alternatives --install /usr/bin/python3 python3 /usr/bin/python3.8 1
#	update-alternatives --set python3 /usr/bin/python3.8
#fi
yum install container-selinux -y
`
		default:
			if script, ok := d.GetOk("custom_script"); ok {
				insertionText = script.(string)
			}
		}

		lines[index] = lines[index] + "\n" + insertionText
		scriptContent = strings.Join(lines, "\n")
	}

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
