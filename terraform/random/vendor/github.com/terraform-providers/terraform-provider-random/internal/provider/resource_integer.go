package provider

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/hashicorp/errwrap"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceInteger() *schema.Resource {
	return &schema.Resource{
		Description: "The resource `random_integer` generates random values from a given range, described " +
			"by the `min` and `max` attributes of a given resource.\n" +
			"\n" +
			"This resource can be used in conjunction with resources that have the `create_before_destroy` " +
			"lifecycle flag set, to avoid conflicts with unique names during the brief period where both the " +
			"old and new resources exist concurrently.",
		Create: CreateInteger,
		Read:   schema.Noop,
		Delete: schema.RemoveFromState,
		Importer: &schema.ResourceImporter{
			State: ImportInteger,
		},

		Schema: map[string]*schema.Schema{
			"keepers": {
				Description: "Arbitrary map of values that, when changed, will trigger recreation of " +
					"resource. See [the main provider documentation](../index.html) for more information.",
				Type:     schema.TypeMap,
				Optional: true,
				ForceNew: true,
			},

			"min": {
				Description: "The minimum inclusive value of the range.",
				Type:        schema.TypeInt,
				Required:    true,
				ForceNew:    true,
			},

			"max": {
				Description: "The maximum inclusive value of the range.",
				Type:        schema.TypeInt,
				Required:    true,
				ForceNew:    true,
			},

			"seed": {
				Description: "A custom seed to always produce the same value.",
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
			},

			"result": {
				Description: "The random integer result.",
				Type:        schema.TypeInt,
				Computed:    true,
			},

			"id": {
				Description: "The string representation of the integer result.",
				Type:        schema.TypeString,
				Computed:    true,
			},
		},
		UseJSONNumber: true,
	}
}

func CreateInteger(d *schema.ResourceData, meta interface{}) error {
	min := d.Get("min").(int)
	max := d.Get("max").(int)
	seed := d.Get("seed").(string)

	if max <= min {
		return fmt.Errorf("Minimum value needs to be smaller than maximum value")
	}
	rand := NewRand(seed)
	number := rand.Intn((max+1)-min) + min

	d.Set("result", number)
	d.SetId(strconv.Itoa(number))

	return nil
}

func ImportInteger(d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	parts := strings.Split(d.Id(), ",")
	if len(parts) != 3 && len(parts) != 4 {
		return nil, fmt.Errorf("Invalid import usage: expecting {result},{min},{max} or {result},{min},{max},{seed}")
	}

	result, err := strconv.Atoi(parts[0])
	if err != nil {
		return nil, errwrap.Wrapf("Error parsing \"result\": {{err}}", err)
	}
	d.Set("result", result)

	min, err := strconv.Atoi(parts[1])
	if err != nil {
		return nil, errwrap.Wrapf("Error parsing \"min\": {{err}}", err)
	}
	d.Set("min", min)

	max, err := strconv.Atoi(parts[2])
	if err != nil {
		return nil, errwrap.Wrapf("Error parsing \"max\": {{err}}", err)
	}
	d.Set("max", max)

	if len(parts) == 4 {
		d.Set("seed", parts[3])
	}
	d.SetId(parts[0])

	return []*schema.ResourceData{d}, nil
}
