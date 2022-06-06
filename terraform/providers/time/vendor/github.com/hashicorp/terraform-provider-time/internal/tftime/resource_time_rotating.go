package tftime

import (
	"context"
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/customdiff"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func resourceTimeRotating() *schema.Resource {
	return &schema.Resource{
		Create: resourceTimeRotatingCreate,
		Read:   resourceTimeRotatingRead,
		Update: resourceTimeRotatingUpdate,
		Delete: schema.Noop,

		CustomizeDiff: customdiff.Sequence(
			customdiff.If(resourceTimeRotatingConditionExpirationChange,
				func(ctx context.Context, diff *schema.ResourceDiff, meta interface{}) error {
					if diff.Id() == "" {
						return nil
					}

					timestamp, err := time.Parse(time.RFC3339, diff.Id())

					if err != nil {
						return fmt.Errorf("error parsing timestamp (%s): %s", diff.Id(), err)
					}

					var rotationTimestamp time.Time

					if v, ok := diff.GetOk("rotation_days"); ok {
						rotationTimestamp = timestamp.AddDate(0, 0, v.(int))
					}

					if v, ok := diff.GetOk("rotation_hours"); ok {
						rotationTimestamp = timestamp.Add(time.Duration(v.(int)) * time.Hour)
					}

					if v, ok := diff.GetOk("rotation_minutes"); ok {
						rotationTimestamp = timestamp.Add(time.Duration(v.(int)) * time.Minute)
					}

					if v, ok := diff.GetOk("rotation_months"); ok {
						rotationTimestamp = timestamp.AddDate(0, v.(int), 0)
					}

					if v, ok := diff.GetOk("rotation_rfc3339"); ok {
						var err error
						rotationTimestamp, err = time.Parse(time.RFC3339, v.(string))

						if err != nil {
							return fmt.Errorf("error parsing rotation_rfc3339 (%s): %s", v.(string), err)
						}
					}

					if v, ok := diff.GetOk("rotation_years"); ok {
						rotationTimestamp = timestamp.AddDate(v.(int), 0, 0)
					}

					if err := diff.SetNew("rotation_rfc3339", rotationTimestamp.Format(time.RFC3339)); err != nil {
						return fmt.Errorf("error setting new rotation_rfc3339: %s", err)
					}

					return nil
				},
			),
			customdiff.ForceNewIf("rotation_rfc3339", func(ctx context.Context, diff *schema.ResourceDiff, meta interface{}) bool {
				now := time.Now().UTC()
				rotationTimestamp, err := time.Parse(time.RFC3339, diff.Get("rotation_rfc3339").(string))

				if err != nil {
					return false
				}

				return now.After(rotationTimestamp)
			}),
		),

		Importer: &schema.ResourceImporter{
			State: func(d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
				idParts := strings.Split(d.Id(), ",")

				if len(idParts) != 2 && len(idParts) != 6 {
					return nil, fmt.Errorf("Unexpected format of ID (%q), expected BASETIMESTAMP,YEARS,MONTHS,DAYS,HOURS,MINUTES or BASETIMESTAMP,ROTATIONTIMESTAMP", d.Id())
				}

				if len(idParts) == 2 {
					if idParts[0] == "" || idParts[1] == "" {
						return nil, fmt.Errorf("Unexpected format of ID (%q), expected BASETIMESTAMP,ROTATIONTIMESTAMP", d.Id())
					}

					baseRfc3339 := idParts[0]
					rotationRfc3339 := idParts[1]

					d.SetId(baseRfc3339)

					if err := d.Set("rotation_rfc3339", rotationRfc3339); err != nil {
						return nil, fmt.Errorf("error setting rotation_rfc3339: %s", err)
					}

					return []*schema.ResourceData{d}, nil
				}

				if idParts[0] == "" || (idParts[1] == "" && idParts[2] == "" && idParts[3] == "" && idParts[4] == "" && idParts[5] == "") {
					return nil, fmt.Errorf("Unexpected format of ID (%q), expected BASETIMESTAMP,YEARS,MONTHS,DAYS,HOURS,MINUTES, where at least one rotation value is non-empty", d.Id())
				}

				baseRfc3339 := idParts[0]
				rotationYears, _ := strconv.Atoi(idParts[1])
				rotationMonths, _ := strconv.Atoi(idParts[2])
				rotationDays, _ := strconv.Atoi(idParts[3])
				rotationHours, _ := strconv.Atoi(idParts[4])
				rotationMinutes, _ := strconv.Atoi(idParts[5])

				d.SetId(baseRfc3339)

				if rotationYears > 0 {
					if err := d.Set("rotation_years", rotationYears); err != nil {
						return nil, fmt.Errorf("error setting rotation_years: %s", err)
					}
				}

				if rotationMonths > 0 {
					if err := d.Set("rotation_months", rotationMonths); err != nil {
						return nil, fmt.Errorf("error setting rotation_months: %s", err)
					}
				}

				if rotationDays > 0 {
					if err := d.Set("rotation_days", rotationDays); err != nil {
						return nil, fmt.Errorf("error setting rotation_days: %s", err)
					}
				}

				if rotationHours > 0 {
					if err := d.Set("rotation_hours", rotationHours); err != nil {
						return nil, fmt.Errorf("error setting rotation_hours: %s", err)
					}
				}

				if rotationMinutes > 0 {
					if err := d.Set("rotation_minutes", rotationMinutes); err != nil {
						return nil, fmt.Errorf("error setting rotation_minutes: %s", err)
					}
				}

				timestamp, err := time.Parse(time.RFC3339, d.Id())

				if err != nil {
					return nil, fmt.Errorf("error parsing timestamp (%s): %s", d.Id(), err)
				}

				var rotationTimestamp time.Time

				if v, ok := d.GetOk("rotation_days"); ok {
					rotationTimestamp = timestamp.AddDate(0, 0, v.(int))
				}

				if v, ok := d.GetOk("rotation_hours"); ok {
					rotationTimestamp = timestamp.Add(time.Duration(v.(int)) * time.Hour)
				}

				if v, ok := d.GetOk("rotation_minutes"); ok {
					rotationTimestamp = timestamp.Add(time.Duration(v.(int)) * time.Minute)
				}

				if v, ok := d.GetOk("rotation_months"); ok {
					rotationTimestamp = timestamp.AddDate(0, v.(int), 0)
				}

				if v, ok := d.GetOk("rotation_years"); ok {
					rotationTimestamp = timestamp.AddDate(v.(int), 0, 0)
				}

				if err := d.Set("rotation_rfc3339", rotationTimestamp.Format(time.RFC3339)); err != nil {
					return nil, fmt.Errorf("error setting rotation_rfc3339: %s", err)
				}

				return []*schema.ResourceData{d}, nil
			},
		},

		Schema: map[string]*schema.Schema{
			"day": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"rotation_days": {
				Type:     schema.TypeInt,
				Optional: true,
				AtLeastOneOf: []string{
					"rotation_days",
					"rotation_hours",
					"rotation_minutes",
					"rotation_months",
					"rotation_rfc3339",
					"rotation_years",
				},
				ValidateFunc: validation.IntAtLeast(1),
			},
			"rotation_hours": {
				Type:     schema.TypeInt,
				Optional: true,
				AtLeastOneOf: []string{
					"rotation_days",
					"rotation_hours",
					"rotation_minutes",
					"rotation_months",
					"rotation_rfc3339",
					"rotation_years",
				},
				ValidateFunc: validation.IntAtLeast(1),
			},
			"rotation_minutes": {
				Type:     schema.TypeInt,
				Optional: true,
				AtLeastOneOf: []string{
					"rotation_days",
					"rotation_hours",
					"rotation_minutes",
					"rotation_months",
					"rotation_rfc3339",
					"rotation_years",
				},
				ValidateFunc: validation.IntAtLeast(1),
			},
			"rotation_months": {
				Type:     schema.TypeInt,
				Optional: true,
				AtLeastOneOf: []string{
					"rotation_days",
					"rotation_hours",
					"rotation_minutes",
					"rotation_months",
					"rotation_rfc3339",
					"rotation_years",
				},
				ValidateFunc: validation.IntAtLeast(1),
			},
			"rotation_rfc3339": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				AtLeastOneOf: []string{
					"rotation_days",
					"rotation_hours",
					"rotation_minutes",
					"rotation_months",
					"rotation_rfc3339",
					"rotation_years",
				},
				ValidateFunc: validation.IsRFC3339Time,
			},
			"rotation_years": {
				Type:     schema.TypeInt,
				Optional: true,
				AtLeastOneOf: []string{
					"rotation_days",
					"rotation_hours",
					"rotation_minutes",
					"rotation_months",
					"rotation_rfc3339",
					"rotation_years",
				},
				ValidateFunc: validation.IntAtLeast(1),
			},
			"hour": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"triggers": {
				Type:     schema.TypeMap,
				Optional: true,
				ForceNew: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"minute": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"month": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"rfc3339": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ForceNew:     true,
				ValidateFunc: validation.IsRFC3339Time,
			},
			"second": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"unix": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"year": {
				Type:     schema.TypeInt,
				Computed: true,
			},
		},
	}
}

func resourceTimeRotatingCreate(d *schema.ResourceData, m interface{}) error {
	timestamp := time.Now().UTC()

	if v, ok := d.GetOk("rfc3339"); ok {
		var err error
		timestamp, err = time.Parse(time.RFC3339, v.(string))

		if err != nil {
			return fmt.Errorf("error parsing rfc3339 (%s): %s", v.(string), err)
		}
	}

	d.SetId(timestamp.Format(time.RFC3339))

	var rotationTimestamp time.Time

	if v, ok := d.GetOk("rotation_days"); ok {
		rotationTimestamp = timestamp.AddDate(0, 0, v.(int))
	}

	if v, ok := d.GetOk("rotation_hours"); ok {
		rotationTimestamp = timestamp.Add(time.Duration(v.(int)) * time.Hour)
	}

	if v, ok := d.GetOk("rotation_minutes"); ok {
		rotationTimestamp = timestamp.Add(time.Duration(v.(int)) * time.Minute)
	}

	if v, ok := d.GetOk("rotation_months"); ok {
		rotationTimestamp = timestamp.AddDate(0, v.(int), 0)
	}

	if v, ok := d.GetOk("rotation_rfc3339"); ok {
		var err error
		rotationTimestamp, err = time.Parse(time.RFC3339, v.(string))

		if err != nil {
			return fmt.Errorf("error parsing rotation_rfc3339 (%s): %s", v.(string), err)
		}
	}

	if v, ok := d.GetOk("rotation_years"); ok {
		rotationTimestamp = timestamp.AddDate(v.(int), 0, 0)
	}

	if err := d.Set("rotation_rfc3339", rotationTimestamp.Format(time.RFC3339)); err != nil {
		return fmt.Errorf("error setting rotation_rfc3339: %s", err)
	}

	return resourceTimeRotatingRead(d, m)
}

func resourceTimeRotatingRead(d *schema.ResourceData, m interface{}) error {
	timestamp, err := time.Parse(time.RFC3339, d.Id())

	if err != nil {
		return fmt.Errorf("error parsing timestamp (%s): %s", d.Id(), err)
	}

	if v, ok := d.GetOk("rotation_rfc3339"); ok && !d.IsNewResource() {
		now := time.Now().UTC()
		rotationTimestamp, err := time.Parse(time.RFC3339, v.(string))

		if err != nil {
			return fmt.Errorf("error parsing rotation_rfc3339 (%s): %s", v.(string), err)
		}

		if now.After(rotationTimestamp) {
			log.Printf("[INFO] Expiration timestamp (%s) is after current timestamp (%s), removing from state", v.(string), now.Format(time.RFC3339))
			d.SetId("")
			return nil
		}
	}

	if err := d.Set("day", timestamp.Day()); err != nil {
		return fmt.Errorf("error setting day: %s", err)
	}

	if err := d.Set("hour", timestamp.Hour()); err != nil {
		return fmt.Errorf("error setting hour: %s", err)
	}

	if err := d.Set("minute", timestamp.Minute()); err != nil {
		return fmt.Errorf("error setting minute: %s", err)
	}

	if err := d.Set("month", int(timestamp.Month())); err != nil {
		return fmt.Errorf("error setting month: %s", err)
	}

	if err := d.Set("rfc3339", timestamp.Format(time.RFC3339)); err != nil {
		return fmt.Errorf("error setting rfc3339: %s", err)
	}

	if err := d.Set("second", timestamp.Second()); err != nil {
		return fmt.Errorf("error setting second: %s", err)
	}

	if err := d.Set("unix", timestamp.Unix()); err != nil {
		return fmt.Errorf("error setting unix: %s", err)
	}

	if err := d.Set("year", timestamp.Year()); err != nil {
		return fmt.Errorf("error setting year: %s", err)
	}

	return nil
}

func resourceTimeRotatingUpdate(d *schema.ResourceData, m interface{}) error {
	if d.HasChanges("rotation_days", "rotation_hours", "rotation_minutes", "rotation_months", "rotation_years") {
		timestamp, err := time.Parse(time.RFC3339, d.Id())

		if err != nil {
			return fmt.Errorf("error parsing timestamp (%s): %s", d.Id(), err)
		}

		var rotationTimestamp time.Time

		if v, ok := d.GetOk("rotation_days"); ok {
			rotationTimestamp = timestamp.AddDate(0, 0, v.(int))
		}

		if v, ok := d.GetOk("rotation_hours"); ok {
			rotationTimestamp = timestamp.Add(time.Duration(v.(int)) * time.Hour)
		}

		if v, ok := d.GetOk("rotation_minutes"); ok {
			rotationTimestamp = timestamp.Add(time.Duration(v.(int)) * time.Minute)
		}

		if v, ok := d.GetOk("rotation_months"); ok {
			rotationTimestamp = timestamp.AddDate(0, v.(int), 0)
		}

		if v, ok := d.GetOk("rotation_years"); ok {
			rotationTimestamp = timestamp.AddDate(v.(int), 0, 0)
		}

		if err := d.Set("rotation_rfc3339", rotationTimestamp.Format(time.RFC3339)); err != nil {
			return fmt.Errorf("error setting rotation_rfc3339: %s", err)
		}
	}

	return resourceTimeRotatingRead(d, m)
}

func resourceTimeRotatingConditionExpirationChange(ctx context.Context, diff *schema.ResourceDiff, meta interface{}) bool {
	return diff.HasChange("rotation_days") ||
		diff.HasChange("rotation_hours") ||
		diff.HasChange("rotation_minutes") ||
		diff.HasChange("rotation_months") ||
		diff.HasChange("rotation_rfc3339") ||
		diff.HasChange("rotation_years")
}
