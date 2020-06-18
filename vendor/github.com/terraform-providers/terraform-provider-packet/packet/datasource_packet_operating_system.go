package packet

import (
	"fmt"
	"log"
	"strings"

	"github.com/hashicorp/terraform/helper/schema"
	"github.com/packethost/packngo"
)

func dataSourceOperatingSystem() *schema.Resource {
	return &schema.Resource{
		Read: dataSourcePacketOperatingSystemRead,
		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"distro": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"version": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"provisionable_on": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"slug": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func dataSourcePacketOperatingSystemRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*packngo.Client)

	name, nameOK := d.GetOk("name")
	distro, distroOK := d.GetOk("distro")
	version, versionOK := d.GetOk("version")
	provisionableOn, provisionableOnOK := d.GetOk("provisionable_on")

	if !nameOK && !distroOK && !versionOK && !provisionableOnOK {
		return fmt.Errorf("One of name, distro, version, or provisionable_on must be assigned")
	}

	log.Println("[DEBUG] ******")
	log.Println("[DEBUG] params", name, distro, version, provisionableOn)
	log.Println("[DEBUG] ******")

	oss, _, err := client.OperatingSystems.List()
	if err != nil {
		return err
	}

	final := []packngo.OS{}
	temp := []packngo.OS{}

	if nameOK {
		for _, os := range oss {
			if strings.Contains(strings.ToLower(os.Name), strings.ToLower(name.(string))) {
				temp = append(temp, os)
			}
			final = temp
		}
	}

	if distroOK {
		temp = []packngo.OS{}
		if len(temp) == 0 {
			final = oss
		}
		for _, v := range final {
			if v.Distro == distro.(string) {
				temp = append(temp, v)
			}
		}
		final = temp
	}

	if versionOK {
		temp = []packngo.OS{}
		if len(final) == 0 {
			final = oss
		}
		for _, v := range final {
			if v.Version == version.(string) {
				temp = append(temp, v)
			}
		}
		final = temp
	}

	if provisionableOnOK {
		temp = []packngo.OS{}
		if len(final) == 0 {
			final = oss
		}
		for _, v := range final {
			for _, po := range v.ProvisionableOn {
				if po == provisionableOn.(string) {
					temp = append(temp, v)
				}
			}
		}
		final = temp
	}
	log.Println("[DEBUG] RESULTS:", final)

	if len(final) == 0 {
		return fmt.Errorf("There are no operating systems that match the search criteria")
	}

	if len(final) > 1 {
		return fmt.Errorf("There is more than one operating system that matches the search criteria")
	}
	for _, v := range final {
		d.Set("name", v.Name)
		d.Set("distro", v.Distro)
		d.Set("version", v.Version)
		d.Set("slug", v.Slug)
		d.SetId(v.Slug)
	}
	return nil
}
