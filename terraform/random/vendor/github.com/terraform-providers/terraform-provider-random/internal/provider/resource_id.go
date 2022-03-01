package provider

import (
	"crypto/rand"
	"encoding/base64"
	"encoding/hex"
	"errors"
	"math/big"
	"strings"

	"github.com/hashicorp/errwrap"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceId() *schema.Resource {
	return &schema.Resource{
		Description: `
The resource ` + "`random_id`" + ` generates random numbers that are intended to be
used as unique identifiers for other resources.

This resource *does* use a cryptographic random number generator in order
to minimize the chance of collisions, making the results of this resource
when a 16-byte identifier is requested of equivalent uniqueness to a
type-4 UUID.

This resource can be used in conjunction with resources that have
the ` + "`create_before_destroy`" + ` lifecycle flag set to avoid conflicts with
unique names during the brief period where both the old and new resources
exist concurrently.
`,
		Create: CreateID,
		Read:   RepopulateEncodings,
		Delete: schema.RemoveFromState,
		Importer: &schema.ResourceImporter{
			State: ImportID,
		},

		Schema: map[string]*schema.Schema{
			"keepers": {
				Description: "Arbitrary map of values that, when changed, will trigger recreation of " +
					"resource. See [the main provider documentation](../index.html) for more information.",
				Type:     schema.TypeMap,
				Optional: true,
				ForceNew: true,
			},

			"byte_length": {
				Description: "The number of random bytes to produce. The minimum value is 1, which produces " +
					"eight bits of randomness.",
				Type:     schema.TypeInt,
				Required: true,
				ForceNew: true,
			},

			"prefix": {
				Description: "Arbitrary string to prefix the output value with. This string is supplied as-is, " +
					"meaning it is not guaranteed to be URL-safe or base64 encoded.",
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},

			"b64_url": {
				Description: "The generated id presented in base64, using the URL-friendly character set: " +
					"case-sensitive letters, digits and the characters `_` and `-`.",
				Type:     schema.TypeString,
				Computed: true,
			},

			"b64_std": {
				Description: "The generated id presented in base64 without additional transformations.",
				Type:        schema.TypeString,
				Computed:    true,
			},

			"hex": {
				Description: "The generated id presented in padded hexadecimal digits. This result will " +
					"always be twice as long as the requested byte length.",
				Type:     schema.TypeString,
				Computed: true,
			},

			"dec": {
				Description: "The generated id presented in non-padded decimal digits.",
				Type:        schema.TypeString,
				Computed:    true,
			},

			"id": {
				Description: "The generated id presented in base64 without additional transformations or prefix.",
				Type:        schema.TypeString,
				Computed:    true,
			},
		},
	}
}

func CreateID(d *schema.ResourceData, meta interface{}) error {
	byteLength := d.Get("byte_length").(int)
	bytes := make([]byte, byteLength)

	n, err := rand.Reader.Read(bytes)
	if n != byteLength {
		return errors.New("generated insufficient random bytes")
	}
	if err != nil {
		return errwrap.Wrapf("error generating random bytes: {{err}}", err)
	}

	b64Str := base64.RawURLEncoding.EncodeToString(bytes)
	d.SetId(b64Str)

	return RepopulateEncodings(d, meta)
}

func RepopulateEncodings(d *schema.ResourceData, _ interface{}) error {
	prefix := d.Get("prefix").(string)
	base64Str := d.Id()

	bytes, err := base64.RawURLEncoding.DecodeString(base64Str)
	if err != nil {
		return errwrap.Wrapf("Error decoding ID: {{err}}", err)
	}

	b64StdStr := base64.StdEncoding.EncodeToString(bytes)
	hexStr := hex.EncodeToString(bytes)

	bigInt := big.Int{}
	bigInt.SetBytes(bytes)
	decStr := bigInt.String()

	d.Set("b64_url", prefix+base64Str)
	d.Set("b64_std", prefix+b64StdStr)

	d.Set("hex", prefix+hexStr)
	d.Set("dec", prefix+decStr)

	return nil
}

func ImportID(d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	id := d.Id()

	sep := strings.LastIndex(id, ",")
	if sep != -1 {
		d.Set("prefix", id[:sep])
		id = id[sep+1:]
	}

	bytes, err := base64.RawURLEncoding.DecodeString(id)
	if err != nil {
		return nil, errwrap.Wrapf("Error decoding ID: {{err}}", err)
	}

	d.Set("byte_length", len(bytes))
	d.SetId(id)

	return []*schema.ResourceData{d}, nil
}
