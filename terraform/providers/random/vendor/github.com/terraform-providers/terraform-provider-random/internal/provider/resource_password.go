package provider

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourcePassword() *schema.Resource {
	return &schema.Resource{
		Description: "Identical to [random_string](string.html) with the exception that the result is " +
			"treated as sensitive and, thus, _not_ displayed in console output. Read more about sensitive " +
			"data handling in the [Terraform documentation](https://www.terraform.io/docs/language/state/sensitive-data.html).\n" +
			"\n" +
			"This resource *does* use a cryptographic random number generator.",
		Create: createStringFunc(true),
		Read:   readNil,
		Delete: schema.RemoveFromState,
		Schema: stringSchemaV1(true),
		Importer: &schema.ResourceImporter{
			StateContext: importStringFunc(true),
		},
	}
}
