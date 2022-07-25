package provider

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"math/big"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
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
		CreateContext: CreateID,
		ReadContext:   RepopulateEncodings,
		DeleteContext: RemoveResourceFromState,
		Importer: &schema.ResourceImporter{
			StateContext: ImportID,
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

func CreateID(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	byteLength := d.Get("byte_length").(int)
	bytes := make([]byte, byteLength)

	n, err := rand.Reader.Read(bytes)
	if n != byteLength {
		return append(diags, diag.Errorf("generated insufficient random bytes: %s", err)...)
	}
	if err != nil {
		return append(diags, diag.Errorf("error generating random bytes: %s", err)...)
	}

	b64Str := base64.RawURLEncoding.EncodeToString(bytes)
	d.SetId(b64Str)

	repopEncsDiags := RepopulateEncodings(ctx, d, meta)
	if repopEncsDiags != nil {
		return append(diags, repopEncsDiags...)
	}

	return diags
}

func RepopulateEncodings(_ context.Context, d *schema.ResourceData, _ interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	prefix := d.Get("prefix").(string)
	base64Str := d.Id()

	bytes, err := base64.RawURLEncoding.DecodeString(base64Str)
	if err != nil {
		return append(diags, diag.Errorf("error decoding ID: %s", err)...)
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

func ImportID(_ context.Context, d *schema.ResourceData, _ interface{}) ([]*schema.ResourceData, error) {
	id := d.Id()

	sep := strings.LastIndex(id, ",")
	if sep != -1 {
		d.Set("prefix", id[:sep])
		id = id[sep+1:]
	}

	bytes, err := base64.RawURLEncoding.DecodeString(id)
	if err != nil {
		return nil, fmt.Errorf("error decoding ID: %w", err)
	}

	d.Set("byte_length", len(bytes))
	d.SetId(id)

	return []*schema.ResourceData{d}, nil
}
