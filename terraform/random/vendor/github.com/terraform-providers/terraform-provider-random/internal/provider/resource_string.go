package provider

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceString() *schema.Resource {
	return &schema.Resource{
		Description: "The resource `random_string` generates a random permutation of alphanumeric " +
			"characters and optionally special characters.\n" +
			"\n" +
			"This resource *does* use a cryptographic random number generator.\n" +
			"\n" +
			"Historically this resource's intended usage has been ambiguous as the original example used " +
			"it in a password. For backwards compatibility it will continue to exist. For unique ids please " +
			"use [random_id](id.html), for sensitive random values please use [random_password](password.html).",
		Create:        createStringFunc(false),
		Read:          readNil,
		Delete:        schema.RemoveFromState,
		MigrateState:  resourceRandomStringMigrateState,
		SchemaVersion: 1,
		Schema:        stringSchemaV1(false),
		Importer: &schema.ResourceImporter{
			StateContext: importStringFunc(false),
		},
	}
}
