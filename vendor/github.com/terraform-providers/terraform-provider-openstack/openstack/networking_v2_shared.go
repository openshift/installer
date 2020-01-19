package openstack

import "github.com/hashicorp/terraform-plugin-sdk/helper/schema"

func networkingV2ReadAttributesTags(d *schema.ResourceData, tags []string) {
	expandObjectReadTags(d, tags)
}

func networkingV2UpdateAttributesTags(d *schema.ResourceData) []string {
	return expandObjectUpdateTags(d)
}

func networkingV2AttributesTags(d *schema.ResourceData) []string {
	return expandObjectTags(d)
}
