package provider

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceShuffle() *schema.Resource {
	return &schema.Resource{
		Description: "The resource `random_shuffle` generates a random permutation of a list of strings " +
			"given as an argument.",
		Create: CreateShuffle,
		Read:   schema.Noop,
		Delete: schema.RemoveFromState,

		Schema: map[string]*schema.Schema{
			"keepers": {
				Description: "Arbitrary map of values that, when changed, will trigger recreation of " +
					"resource. See [the main provider documentation](../index.html) for more information.",
				Type:     schema.TypeMap,
				Optional: true,
				ForceNew: true,
			},

			"seed": {
				Description: "Arbitrary string with which to seed the random number generator, in order to " +
					"produce less-volatile permutations of the list.\n" +
					"\n" +
					"**Important:** Even with an identical seed, it is not guaranteed that the same permutation " +
					"will be produced across different versions of Terraform. This argument causes the " +
					"result to be *less volatile*, but not fixed for all time.",
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},

			"input": {
				Description: "The list of strings to shuffle.",
				Type:        schema.TypeList,
				Required:    true,
				ForceNew:    true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},

			"result_count": {
				Description: "The number of results to return. Defaults to the number of items in the " +
					"`input` list. If fewer items are requested, some elements will be excluded from the " +
					"result. If more items are requested, items will be repeated in the result but not more " +
					"frequently than the number of items in the input list.",
				Type:     schema.TypeInt,
				Optional: true,
				ForceNew: true,
			},

			"result": {
				Description: "Random permutation of the list of strings given in `input`.",
				Type:        schema.TypeList,
				Computed:    true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},

			"id": {
				Description: "A static value used internally by Terraform, this should not be referenced in configurations.",
				Type:        schema.TypeString,
				Computed:    true,
			},
		},
	}
}

func CreateShuffle(d *schema.ResourceData, _ interface{}) error {
	input := d.Get("input").([]interface{})
	seed := d.Get("seed").(string)

	resultCount := d.Get("result_count").(int)
	if resultCount == 0 {
		resultCount = len(input)
	}
	result := make([]interface{}, 0, resultCount)

	if len(input) > 0 {
		rand := NewRand(seed)

		// Keep producing permutations until we fill our result
	Batches:
		for {
			perm := rand.Perm(len(input))

			for _, i := range perm {
				result = append(result, input[i])

				if len(result) >= resultCount {
					break Batches
				}
			}
		}

	}

	d.SetId("-")
	d.Set("result", result)

	return nil
}
