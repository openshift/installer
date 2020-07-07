package google

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func dataSourceGoogleFirebaseWebApp() *schema.Resource {
	// Generate datasource schema from resource
	dsSchema := datasourceSchemaFromResourceSchema(resourceFirebaseWebApp().Schema)

	// Set 'Required' schema elements
	addRequiredFieldsToSchema(dsSchema, "app_id")

	return &schema.Resource{
		Read:   dataSourceGoogleFirebaseWebAppRead,
		Schema: dsSchema,
	}
}

func dataSourceGoogleFirebaseWebAppRead(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	appId := d.Get("app_id")
	project, err := getProject(d, config)
	if err != nil {
		return err
	}
	name := fmt.Sprintf("projects/%s/webApps/%s", project, appId.(string))
	d.SetId(name)
	d.Set("name", name)
	return resourceFirebaseWebAppRead(d, meta)
}
