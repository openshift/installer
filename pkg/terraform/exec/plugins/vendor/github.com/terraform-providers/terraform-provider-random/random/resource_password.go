package random

import (
	"github.com/hashicorp/terraform/helper/schema"
)

func resourcePassword() *schema.Resource {
	return &schema.Resource{
		Create: createStringFunc(true),
		Read:   readNil,
		Delete: schema.RemoveFromState,
		Schema: stringSchemaV1(true),
	}
}
