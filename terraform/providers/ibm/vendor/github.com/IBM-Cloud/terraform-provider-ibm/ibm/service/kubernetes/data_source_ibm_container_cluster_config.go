// Copyright IBM Corp. 2017, 2021 All Rights Reserved.
// Licensed under the Mozilla Public License v2.0

package kubernetes

import (
	"fmt"
	"log"
	"path/filepath"
	"regexp"
	"strings"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	homedir "github.com/mitchellh/go-homedir"

	v1 "github.com/IBM-Cloud/bluemix-go/api/container/containerv1"
	"github.com/IBM-Cloud/bluemix-go/helpers"
	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/conns"
)

func DataSourceIBMContainerClusterConfig() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceIBMContainerClusterConfigRead,

		Schema: map[string]*schema.Schema{

			"org_guid": {
				Description: "The bluemix organization guid this cluster belongs to",
				Type:        schema.TypeString,
				Optional:    true,
				Deprecated:  "This field is deprecated",
			},
			"space_guid": {
				Description: "The bluemix space guid this cluster belongs to",
				Type:        schema.TypeString,
				Optional:    true,
				Deprecated:  "This field is deprecated",
			},
			"account_guid": {
				Description: "The bluemix account guid this cluster belongs to",
				Type:        schema.TypeString,
				Optional:    true,
				Deprecated:  "This field is deprecated",
			},
			"region": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The cluster region",
				Deprecated:  "This field is deprecated",
			},
			"resource_group_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "ID of the resource group.",
			},
			"cluster_name_id": {
				Description: "The name/id of the cluster",
				Type:        schema.TypeString,
				Required:    true,
			},
			"config_dir": {
				Description: "The directory where the cluster config to be downloaded. Default is home directory ",
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
			},
			"download": {
				Description: "If set to false will not download the config, otherwise they are downloaded each time but onto the same path for a given cluster name/id",
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     true,
			},
			"admin": {
				Description: "If set to true will download the config for admin",
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
			},
			"network": {
				Description: "If set to true will download the Calico network config with the Admin config",
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
			},
			"config_file_path": {
				Description: "The absolute path to the kubernetes config yml file ",
				Type:        schema.TypeString,
				Computed:    true,
			},
			"calico_config_file_path": {
				Description: "The absolute path to the calico network config file ",
				Type:        schema.TypeString,
				Computed:    true,
			},
			"admin_key": {
				Type:      schema.TypeString,
				Computed:  true,
				Sensitive: true,
			},
			"admin_certificate": {
				Type:      schema.TypeString,
				Computed:  true,
				Sensitive: true,
			},
			"ca_certificate": {
				Type:      schema.TypeString,
				Computed:  true,
				Sensitive: true,
			},
			"host": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"token": {
				Type:      schema.TypeString,
				Computed:  true,
				Sensitive: true,
			},
		},
	}
}

func dataSourceIBMContainerClusterConfigRead(d *schema.ResourceData, meta interface{}) error {
	csClient, err := meta.(conns.ClientSession).VpcContainerAPI()
	if err != nil {
		return err
	}
	csAPI := csClient.Clusters()
	name := d.Get("cluster_name_id").(string)
	download := d.Get("download").(bool)
	admin := d.Get("admin").(bool)
	configDir := d.Get("config_dir").(string)
	network := d.Get("network").(bool)

	clusterId := "Cluster_Config_" + name
	conns.IbmMutexKV.Lock(clusterId)
	defer conns.IbmMutexKV.Unlock(clusterId)

	if len(configDir) == 0 {
		configDir, err = homedir.Dir()
		if err != nil {
			return fmt.Errorf("[ERROR] Error fetching homedir: %s", err)
		}
	}
	configDir, _ = filepath.Abs(configDir)

	var configPath string
	if !download {
		log.Println("Skipping download of the cluster config", "Going to check if it already exists")
		expectedDir := v1.ComputeClusterConfigDir(configDir, name, admin)
		configPath = filepath.Join(expectedDir, "config.yml")
		if !helpers.FileExists(configPath) {
			return fmt.Errorf(`[ERROR] Couldn't  find the cluster config at expected path %s. Please set "download" to true to download the new config`, configPath)
		}
		d.Set("config_file_path", configPath)

	} else {
		targetEnv, err := getVpcClusterTargetHeader(d, meta)
		if err != nil {
			return err
		}
		if network {
			// For the Network config we need to gather the certs so we must override the admin value
			var calicoConfigFilePath string
			var clusterKeyDetails v1.ClusterKeyInfo
			err = resource.Retry(5*time.Minute, func() *resource.RetryError {
				var err error
				calicoConfigFilePath, clusterKeyDetails, err = csAPI.StoreConfigDetail(name, configDir, admin || true, network, targetEnv)
				if err != nil {
					log.Printf("[DEBUG] Failed to fetch cluster config err %s", err)
					if strings.Contains(err.Error(), "Could not login to openshift account runtime error:") {
						return resource.RetryableError(err)
					}
					if intermittentUserLookupFailure, _ := regexp.MatchString("Error: lookup of user for \"(.+)\" failed", err.Error()); intermittentUserLookupFailure {
						// Intermittent error resulting from synchronisation delay
						return resource.RetryableError(err)
					}
					return resource.NonRetryableError(err)
				}
				return nil
			})
			if conns.IsResourceTimeoutError(err) {
				calicoConfigFilePath, clusterKeyDetails, err = csAPI.StoreConfigDetail(name, configDir, admin || true, network, targetEnv)
			}
			if err != nil {
				return fmt.Errorf("[ERROR] Error downloading the cluster config [%s]: %s", name, err)
			}
			d.Set("calico_config_file_path", calicoConfigFilePath)
			d.Set("admin_key", clusterKeyDetails.AdminKey)
			d.Set("admin_certificate", clusterKeyDetails.Admin)
			d.Set("ca_certificate", clusterKeyDetails.ClusterCACertificate)
			d.Set("host", clusterKeyDetails.Host)
			d.Set("token", clusterKeyDetails.Token)
			d.Set("config_file_path", clusterKeyDetails.FilePath)

		} else {
			var clusterKeyDetails v1.ClusterKeyInfo
			err = resource.Retry(5*time.Minute, func() *resource.RetryError {
				var err error
				clusterKeyDetails, err = csAPI.GetClusterConfigDetail(name, configDir, admin, targetEnv)
				if err != nil {
					log.Printf("[DEBUG] Failed to fetch cluster config err %s", err)
					if strings.Contains(err.Error(), "Could not login to openshift account runtime error:") {
						return resource.RetryableError(err)
					}
					if intermittentUserLookupFailure, _ := regexp.MatchString("Error: lookup of user for \"(.+)\" failed", err.Error()); intermittentUserLookupFailure {
						// Intermittent error resulting from synchronisation delay
						return resource.RetryableError(err)
					}
					return resource.NonRetryableError(err)
				}
				return nil
			})
			if conns.IsResourceTimeoutError(err) {
				clusterKeyDetails, err = csAPI.GetClusterConfigDetail(name, configDir, admin, targetEnv)
			}
			if err != nil {
				return fmt.Errorf("[ERROR] Error downloading the cluster config [%s]: %s", name, err)
			}
			d.Set("admin_key", clusterKeyDetails.AdminKey)
			d.Set("admin_certificate", clusterKeyDetails.Admin)
			d.Set("ca_certificate", clusterKeyDetails.ClusterCACertificate)
			d.Set("host", clusterKeyDetails.Host)
			d.Set("token", clusterKeyDetails.Token)
			d.Set("config_file_path", clusterKeyDetails.FilePath)
		}
	}

	d.SetId(name)
	d.Set("config_dir", configDir)
	return nil
}
