package eventnotification

import "github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

// Wrapper function around  deprecated GetOkExists function with same functionality
func GetFieldExists(d *schema.ResourceData, field string) (interface{}, bool) {
	return d.GetOkExists(field)
}
