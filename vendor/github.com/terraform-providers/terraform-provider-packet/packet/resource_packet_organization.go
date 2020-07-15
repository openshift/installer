package packet

import (
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/packethost/packngo"
)

func resourcePacketOrganization() *schema.Resource {
	return &schema.Resource{
		Create: resourcePacketOrganizationCreate,
		Read:   resourcePacketOrganizationRead,
		Update: resourcePacketOrganizationUpdate,
		Delete: resourcePacketOrganizationDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},

			"description": {
				Type:     schema.TypeString,
				Optional: true,
				Required: false,
			},

			"website": {
				Type:     schema.TypeString,
				Optional: true,
				Required: false,
			},

			"twitter": {
				Type:     schema.TypeString,
				Optional: true,
				Required: false,
			},

			"logo": {
				Type:     schema.TypeString,
				Optional: true,
				Required: false,
			},

			"created": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"updated": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourcePacketOrganizationCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*packngo.Client)

	createRequest := &packngo.OrganizationCreateRequest{
		Name: d.Get("name").(string),
	}

	if attr, ok := d.GetOk("website"); ok {
		createRequest.Website = attr.(string)
	}

	if attr, ok := d.GetOk("description"); ok {
		createRequest.Description = attr.(string)
	}

	if attr, ok := d.GetOk("twitter"); ok {
		createRequest.Twitter = attr.(string)
	}

	if attr, ok := d.GetOk("logo"); ok {
		createRequest.Logo = attr.(string)
	}

	org, _, err := client.Organizations.Create(createRequest)
	if err != nil {
		return friendlyError(err)
	}

	d.SetId(org.ID)

	return resourcePacketOrganizationRead(d, meta)
}

func resourcePacketOrganizationRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*packngo.Client)

	key, _, err := client.Organizations.Get(d.Id(), nil)
	if err != nil {
		err = friendlyError(err)

		// If the project somehow already destroyed, mark as succesfully gone.
		if isNotFound(err) {
			d.SetId("")

			return nil
		}

		return err
	}

	d.Set("id", key.ID)
	d.Set("name", key.Name)
	d.Set("description", key.Description)
	d.Set("website", key.Website)
	d.Set("twitter", key.Twitter)
	d.Set("logo", key.Logo)
	d.Set("created", key.Created)
	d.Set("updated", key.Updated)

	return nil
}

func resourcePacketOrganizationUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*packngo.Client)

	updateRequest := &packngo.OrganizationUpdateRequest{}

	if d.HasChange("name") {
		oName := d.Get("name").(string)
		updateRequest.Name = &oName
	}

	if d.HasChange("description") {
		oDescription := d.Get("description").(string)
		updateRequest.Description = &oDescription
	}

	if d.HasChange("website") {
		oWebsite := d.Get("website").(string)
		updateRequest.Website = &oWebsite
	}

	if d.HasChange("twitter") {
		oTwitter := d.Get("twitter").(string)
		updateRequest.Twitter = &oTwitter
	}

	if d.HasChange("logo") {
		oLogo := d.Get("logo").(string)
		updateRequest.Logo = &oLogo
	}
	_, _, err := client.Organizations.Update(d.Id(), updateRequest)
	if err != nil {
		return friendlyError(err)
	}

	return resourcePacketOrganizationRead(d, meta)
}

func resourcePacketOrganizationDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*packngo.Client)

	_, err := client.Organizations.Delete(d.Id())
	if err != nil {
		return friendlyError(err)
	}

	d.SetId("")
	return nil
}
