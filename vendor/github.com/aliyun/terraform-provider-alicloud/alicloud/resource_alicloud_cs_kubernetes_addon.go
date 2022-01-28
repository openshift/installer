package alicloud

import (
	"fmt"
	"time"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

const ResourceAlicloudCSKubernetesAddon = "resourceAlicloudCSKubernetesAddon"

func resourceAlicloudCSKubernetesAddon() *schema.Resource {
	return &schema.Resource{
		Create: resourceAlicloudCSKubernetesAddonCreate,
		Read:   resourceAlicloudCSKubernetesAddonRead,
		Update: resourceAlicloudCSKubernetesAddonUpdate,
		Delete: resourceAlicloudCSKubernetesAddonDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(10 * time.Minute),
			Update: schema.DefaultTimeout(10 * time.Minute),
			Delete: schema.DefaultTimeout(10 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			"cluster_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"version": {
				Type:     schema.TypeString,
				Required: true,
			},
			"next_version": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"can_upgrade": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"required": {
				Type:     schema.TypeBool,
				Computed: true,
			},
		},
	}
}

func resourceAlicloudCSKubernetesAddonRead(d *schema.ResourceData, meta interface{}) error {
	client, err := meta.(*connectivity.AliyunClient).NewRoaCsClient()
	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, ResourceAlicloudCSKubernetesAddon, "InitializeClient", err)
	}
	csClient := CsClient{client}

	parts, err := ParseResourceId(d.Id(), 2)
	if err != nil {
		return WrapError(err)
	}

	addon, err := csClient.DescribeCsKubernetesAddon(d.Id())
	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, ResourceAlicloudCSKubernetesAddon, "DescribeCsKubernetesAddon", err)
	}

	d.Set("cluster_id", parts[0])
	d.Set("name", addon.ComponentName)
	d.Set("version", addon.Version)
	d.Set("next_version", addon.NextVersion)
	d.Set("can_upgrade", addon.CanUpgrade)
	d.Set("required", addon.Required)

	return nil
}

func resourceAlicloudCSKubernetesAddonCreate(d *schema.ResourceData, meta interface{}) error {
	clusterId := d.Get("cluster_id").(string)
	name := d.Get("name").(string)
	version := d.Get("version").(string)

	client, err := meta.(*connectivity.AliyunClient).NewRoaCsClient()
	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, ResourceAlicloudCSKubernetesAddon, "InitializeClient", err)
	}
	csClient := CsClient{client}

	status, err := csClient.DescribeCsKubernetesAddonStatus(clusterId, name)
	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, ResourceAlicloudCSKubernetesAddon, "DescribeCsKubernetesAddonStatus", err)
	}

	changed := false

	if status.Version == "" {
		// Addon not installed
		changed = true
		err := csClient.installAddon(d)
		if err != nil {
			return WrapErrorf(err, DefaultErrorMsg, ResourceAlicloudCSKubernetesAddon, "installAddon", err)
		}
	} else if status.Version != "" && status.Version != version {
		// Addon has been installed
		changed = true
		err := csClient.upgradeAddon(d)
		if err != nil {
			return WrapErrorf(err, DefaultErrorMsg, ResourceAlicloudCSKubernetesAddon, "upgradeAddon", err)
		}
	}

	d.SetId(fmt.Sprintf("%s%s%s", clusterId, COLON_SEPARATED, name))

	stateConf := BuildStateConf([]string{"running", "Upgrading", "Pause"}, []string{"Success"}, d.Timeout(schema.TimeoutCreate), 10*time.Second, csClient.CsKubernetesAddonStateRefreshFunc(clusterId, name, []string{"Failed", "Canceled"}))
	if changed {
		if _, err := stateConf.WaitForState(); err != nil {
			return WrapErrorf(err, DefaultErrorMsg, ResourceAlicloudCSKubernetesAddon, "WaitForSuccessAfterCreate", d.Id())
		}
	}
	return resourceAlicloudCSKubernetesAddonRead(d, meta)
}

func resourceAlicloudCSKubernetesAddonUpdate(d *schema.ResourceData, meta interface{}) error {
	client, err := meta.(*connectivity.AliyunClient).NewRoaCsClient()
	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, ResourceAlicloudCSKubernetesAddon, "InitializeClient", err)
	}
	csClient := CsClient{client}

	if d.HasChange("version") {
		clusterId := d.Get("cluster_id").(string)
		name := d.Get("name").(string)

		err := csClient.upgradeAddon(d)
		if err != nil {
			return WrapErrorf(err, DefaultErrorMsg, ResourceAlicloudCSKubernetesAddon, "upgradeAddon", err)
		}

		stateConf := BuildStateConf([]string{"running", "Upgrading", "Pause"}, []string{"Success"}, d.Timeout(schema.TimeoutUpdate), 10*time.Second, csClient.CsKubernetesAddonStateRefreshFunc(clusterId, name, []string{"Failed", "Canceled"}))
		if _, err := stateConf.WaitForState(); err != nil {
			return WrapErrorf(err, DefaultErrorMsg, ResourceAlicloudCSKubernetesAddon, "WaitForSuccessAfterUpdate", d.Id())
		}
	}

	return resourceAlicloudCSKubernetesAddonRead(d, meta)
}

func resourceAlicloudCSKubernetesAddonDelete(d *schema.ResourceData, meta interface{}) error {
	clusterId := d.Get("cluster_id").(string)
	name := d.Get("name").(string)

	client, err := meta.(*connectivity.AliyunClient).NewRoaCsClient()
	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, ResourceAlicloudCSKubernetesAddon, "InitializeClient", err)
	}
	csClient := CsClient{client}

	err = csClient.uninstallAddon(d)
	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, ResourceAlicloudCSKubernetesAddon, "uninstallAddon", err)
	}

	stateConf := BuildStateConf([]string{"running", "Upgrading", "Pause"}, []string{"Success"}, d.Timeout(schema.TimeoutDelete), 10*time.Second, csClient.CsKubernetesAddonStateRefreshFunc(clusterId, name, []string{"Failed", "Canceled"}))
	if _, err := stateConf.WaitForState(); err != nil {
		return WrapErrorf(err, DefaultErrorMsg, ResourceAlicloudCSKubernetesAddon, "WaitForSuccessAfterDelete", d.Id())
	}

	return nil
}
