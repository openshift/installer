package openstack

import (
	"fmt"
	"strings"

	"github.com/gophercloud/gophercloud/openstack/identity/v3/users"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceIdentityUserMembershipV3() *schema.Resource {
	return &schema.Resource{
		Create: resourceIdentityUserMembershipV3Create,
		Read:   resourceIdentityUserMembershipV3Read,
		Delete: resourceIdentityUserMembershipV3Delete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"user_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"group_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
		},
	}
}

func resourceIdentityUserMembershipV3Create(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	identityClient, err := config.IdentityV3Client(GetRegion(d, config))
	if err != nil {
		return fmt.Errorf("Error creating OpenStack identity client: %s", err)
	}

	userID := d.Get("user_id").(string)
	groupID := d.Get("group_id").(string)

	if err := users.AddToGroup(identityClient, groupID, userID).ExtractErr(); err != nil {
		return fmt.Errorf("Error creating openstack_identity_user_membership_v3: %s", err)
	}

	id := fmt.Sprintf("%s/%s", userID, groupID)
	d.SetId(id)

	return resourceIdentityUserMembershipV3Read(d, meta)
}

func resourceIdentityUserMembershipV3Read(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	identityClient, err := config.IdentityV3Client(GetRegion(d, config))
	if err != nil {
		return fmt.Errorf("Error creating OpenStack identity client: %s", err)
	}

	userID, groupID, err := parseUserMembershipID(d.Id())
	if err != nil {
		return CheckDeleted(d, err, "Error parsing ID of openstack_identity_user_membership_v3")
	}

	userMembership, err := users.IsMemberOfGroup(identityClient, groupID, userID).Extract()
	if err != nil || !userMembership {
		return CheckDeleted(d, err, "Error retrieving openstack_identity_user_membership_v3")
	}

	d.Set("region", GetRegion(d, config))
	d.Set("user_id", userID)
	d.Set("group_id", groupID)

	return nil
}

func resourceIdentityUserMembershipV3Delete(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	identityClient, err := config.IdentityV3Client(GetRegion(d, config))
	if err != nil {
		return fmt.Errorf("Error creating OpenStack identity client: %s", err)
	}

	userID, groupID, err := parseUserMembershipID(d.Id())
	if err != nil {
		return CheckDeleted(d, err, "Error parsing ID of openstack_identity_user_membership_v3")
	}

	if err := users.RemoveFromGroup(identityClient, groupID, userID).ExtractErr(); err != nil {
		return CheckDeleted(d, err, "Error removing openstack_identity_user_membership_v3")
	}

	return nil
}

func parseUserMembershipID(id string) (string, string, error) {
	idParts := strings.Split(id, "/")
	if len(idParts) < 2 {
		return "", "", fmt.Errorf("Unable to determine user membership ID %s", id)
	}

	userID := idParts[0]
	groupID := idParts[1]

	return userID, groupID, nil
}
