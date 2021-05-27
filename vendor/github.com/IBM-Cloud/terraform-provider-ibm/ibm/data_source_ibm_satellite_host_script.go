// Copyright IBM Corp. 2017, 2021 All Rights Reserved.
// Licensed under the Mozilla Public License v2.0

package ibm

import (
	"fmt"
	"io/ioutil"
	"log"
	"path/filepath"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	homedir "github.com/mitchellh/go-homedir"
	"github.ibm.com/ibmcloud/kubernetesservice-go-sdk/kubernetesserviceapiv1"
)

func dataSourceIBMSatelliteAttachHostScript() *schema.Resource {
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
				Type:     schema.TypeString,
				Required: true,
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

	satClient, err := meta.(ClientSession).SatelliteClientSession()
	if err != nil {
		return err
	}

	getSatLocOptions := &kubernetesserviceapiv1.GetSatelliteLocationOptions{
		Controller: &location,
	}

	locData, _, err := satClient.GetSatelliteLocation(getSatLocOptions)
	if err != nil {
		return err
	}

	// script labels
	labels := make(map[string]string)
	if v, ok := d.GetOk("labels"); ok {
		l := v.(*schema.Set)
		labels = flattenHostLabels(l.List())
		d.Set("labels", l)
	}

	if len(scriptDir) == 0 {
		scriptDir, err = homedir.Dir()
		if err != nil {
			return fmt.Errorf("Error fetching homedir: %s", err)
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
		return fmt.Errorf("Error Generating Satellite Registration Script: %s\n%s", err, resp)
	}

	lines := strings.Split(string(resp), "\n")
	for i, line := range lines {
		if strings.Contains(line, "API_URL=") {
			i = i + 1
			if strings.ToLower(hostProvider) == "aws" {
				lines[i] = "yum update -y\nyum-config-manager --enable '*'\nyum repolist all\nyum install container-selinux -y"
			} else if strings.ToLower(hostProvider) == "ibm" {
				lines[i] = "subscription-manager refresh\nsubscription-manager repos --enable=*\n"
			} else if strings.ToLower(hostProvider) == "azure" {
				lines[i] = fmt.Sprintf(`#Grow the base volume group first
echo -e "r\ne\ny\nw\ny\ny\n" | gdisk /dev/sda
#mark result as true as this returns a non-0 RC when syncing disks
echo -e "n\n\n\n\n\nw\n" | fdisk /dev/sda || true
partx -l /dev/sda || true
partx -v -a /dev/sda || true
pvcreate /dev/sda5
vgextend rootvg /dev/sda5
# Grow the TMP LV
lvextend -L+10G /dev/rootvg/tmplv
xfs_growfs /dev/rootvg/tmplv
# Grow the var LV
lvextend -L+20G /dev/rootvg/varlv
xfs_growfs /dev/rootvg/varlv
yum update --disablerepo=* --enablerepo="*microsoft*" -y
yum-config-manager --enable '*'
yum repolist all
yum install container-selinux -y
				`)
			}
		}
	}

	scriptContent := strings.Join(lines, "\n")
	err = ioutil.WriteFile(scriptPath, []byte(scriptContent), 0644)
	if err != nil {
		return fmt.Errorf("Error Creating Satellite Attach Host Script: %s", err)
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
