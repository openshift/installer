package provider

import (
	"context"
	"fmt"
	"strings"

	petname "github.com/dustinkirkland/golang-petname"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourcePet() *schema.Resource {
	// This is necessary to ensure each call to petname is properly randomised:
	// the library uses `rand.Intn()` and does NOT seed `rand.Seed()` by default,
	// so this call takes care of that.
	petname.NonDeterministicMode()

	return &schema.Resource{
		Description: "The resource `random_pet` generates random pet names that are intended to be used as " +
			"unique identifiers for other resources.\n" +
			"\n" +
			"This resource can be used in conjunction with resources that have the `create_before_destroy` " +
			"lifecycle flag set, to avoid conflicts with unique names during the brief period where both the old " +
			"and new resources exist concurrently.",
		CreateContext: CreatePet,
		ReadContext:   schema.NoopContext,
		DeleteContext: RemoveResourceFromState,

		Schema: map[string]*schema.Schema{
			"keepers": {
				Description: "Arbitrary map of values that, when changed, will trigger recreation of " +
					"resource. See [the main provider documentation](../index.html) for more information.",
				Type:     schema.TypeMap,
				Optional: true,
				ForceNew: true,
			},

			"length": {
				Description: "The length (in words) of the pet name.",
				Type:        schema.TypeInt,
				Optional:    true,
				Default:     2,
				ForceNew:    true,
			},

			"prefix": {
				Description: "A string to prefix the name with.",
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
			},

			"separator": {
				Description: "The character to separate words in the pet name.",
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "-",
				ForceNew:    true,
			},

			"id": {
				Description: "The random pet name",
				Type:        schema.TypeString,
				Computed:    true,
			},
		},
	}
}

func CreatePet(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	length := d.Get("length").(int)
	separator := d.Get("separator").(string)
	prefix := d.Get("prefix").(string)

	pet := strings.ToLower(petname.Generate(length, separator))

	if prefix != "" {
		pet = fmt.Sprintf("%s%s%s", prefix, separator, pet)
	}

	d.SetId(pet)

	return nil
}
