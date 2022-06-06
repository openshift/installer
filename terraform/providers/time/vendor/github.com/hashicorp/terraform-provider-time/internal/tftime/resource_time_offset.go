package tftime

import (
	"context"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/customdiff"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func resourceTimeOffset() *schema.Resource {
	return &schema.Resource{
		Create: resourceTimeOffsetCreate,
		Read:   resourceTimeOffsetRead,
		Update: resourceTimeOffsetUpdate,
		Delete: schema.Noop,

		CustomizeDiff: customdiff.Sequence(
			customdiff.If(resourceTimeOffsetConditionExpirationChange,
				func(ctx context.Context, diff *schema.ResourceDiff, meta interface{}) error {
					if diff.Id() == "" {
						return nil
					}

					timestamp, err := time.Parse(time.RFC3339, diff.Id())

					if err != nil {
						return fmt.Errorf("error parsing timestamp (%s): %s", diff.Id(), err)
					}

					if v, ok := diff.GetOk("offset_days"); ok {
						timestamp = timestamp.AddDate(0, 0, v.(int))
					}

					if v, ok := diff.GetOk("offset_hours"); ok {
						timestamp = timestamp.Add(time.Duration(v.(int)) * time.Hour)
					}

					if v, ok := diff.GetOk("offset_minutes"); ok {
						timestamp = timestamp.Add(time.Duration(v.(int)) * time.Minute)
					}

					if v, ok := diff.GetOk("offset_months"); ok {
						timestamp = timestamp.AddDate(0, v.(int), 0)
					}

					if v, ok := diff.GetOk("offset_seconds"); ok {
						timestamp = timestamp.Add(time.Duration(v.(int)) * time.Second)
					}

					if v, ok := diff.GetOk("offset_years"); ok {
						timestamp = timestamp.AddDate(v.(int), 0, 0)
					}

					if err := diff.SetNew("day", timestamp.Day()); err != nil {
						return fmt.Errorf("error setting new day: %s", err)
					}

					if err := diff.SetNew("hour", timestamp.Hour()); err != nil {
						return fmt.Errorf("error setting new hour: %s", err)
					}

					if err := diff.SetNew("minute", timestamp.Minute()); err != nil {
						return fmt.Errorf("error setting new minute: %s", err)
					}

					if err := diff.SetNew("month", int(timestamp.Month())); err != nil {
						return fmt.Errorf("error setting new month: %s", err)
					}

					if err := diff.SetNew("rfc3339", timestamp.Format(time.RFC3339)); err != nil {
						return fmt.Errorf("error setting new rfc3339: %s", err)
					}

					if err := diff.SetNew("second", timestamp.Second()); err != nil {
						return fmt.Errorf("error setting new second: %s", err)
					}

					if err := diff.SetNew("unix", timestamp.Unix()); err != nil {
						return fmt.Errorf("error setting new unix: %s", err)
					}

					if err := diff.SetNew("year", timestamp.Year()); err != nil {
						return fmt.Errorf("error setting new year: %s", err)
					}

					return nil
				},
			),
		),

		Importer: &schema.ResourceImporter{
			State: func(d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
				idParts := strings.Split(d.Id(), ",")

				if len(idParts) != 7 {
					return nil, fmt.Errorf("Unexpected format of ID (%q), expected BASETIMESTAMP,YEARS,MONTHS,DAYS,HOURS,MINUTES,SECONDS", d.Id())
				}

				if idParts[0] == "" || (idParts[1] == "" && idParts[2] == "" && idParts[3] == "" && idParts[4] == "" && idParts[5] == "" && idParts[6] == "") {
					return nil, fmt.Errorf("Unexpected format of ID (%q), expected BASETIMESTAMP,YEARS,MONTHS,DAYS,HOURS,MINUTES,SECONDS where at least one offset value is non-empty", d.Id())
				}

				baseRfc3339 := idParts[0]
				offsetYears, _ := strconv.Atoi(idParts[1])
				offsetMonths, _ := strconv.Atoi(idParts[2])
				offsetDays, _ := strconv.Atoi(idParts[3])
				offsetHours, _ := strconv.Atoi(idParts[4])
				offsetMinutes, _ := strconv.Atoi(idParts[5])
				offsetSeconds, _ := strconv.Atoi(idParts[6])

				d.SetId(baseRfc3339)

				if err := d.Set("base_rfc3339", baseRfc3339); err != nil {
					return nil, fmt.Errorf("error setting base_rfc3339: %s", err)
				}

				if offsetYears > 0 {
					if err := d.Set("offset_years", offsetYears); err != nil {
						return nil, fmt.Errorf("error setting offset_years: %s", err)
					}
				}

				if offsetMonths > 0 {
					if err := d.Set("offset_months", offsetMonths); err != nil {
						return nil, fmt.Errorf("error setting offset_months: %s", err)
					}
				}

				if offsetDays > 0 {
					if err := d.Set("offset_days", offsetDays); err != nil {
						return nil, fmt.Errorf("error setting offset_days: %s", err)
					}
				}

				if offsetHours > 0 {
					if err := d.Set("offset_hours", offsetHours); err != nil {
						return nil, fmt.Errorf("error setting offset_hours: %s", err)
					}
				}

				if offsetMinutes > 0 {
					if err := d.Set("offset_minutes", offsetMinutes); err != nil {
						return nil, fmt.Errorf("error setting offset_minutes: %s", err)
					}
				}

				if offsetSeconds > 0 {
					if err := d.Set("offset_seconds", offsetSeconds); err != nil {
						return nil, fmt.Errorf("error setting offset_seconds: %s", err)
					}
				}

				timestamp, err := time.Parse(time.RFC3339, d.Id())

				if err != nil {
					return nil, fmt.Errorf("error parsing base timestamp (%s): %s", d.Id(), err)
				}

				if v, ok := d.GetOk("offset_days"); ok {
					timestamp = timestamp.AddDate(0, 0, v.(int))
				}

				if v, ok := d.GetOk("offset_hours"); ok {
					timestamp = timestamp.Add(time.Duration(v.(int)) * time.Hour)
				}

				if v, ok := d.GetOk("offset_minutes"); ok {
					timestamp = timestamp.Add(time.Duration(v.(int)) * time.Minute)
				}

				if v, ok := d.GetOk("offset_months"); ok {
					timestamp = timestamp.AddDate(0, v.(int), 0)
				}

				if v, ok := d.GetOk("offset_seconds"); ok {
					timestamp = timestamp.Add(time.Duration(v.(int)) * time.Second)
				}

				if v, ok := d.GetOk("offset_years"); ok {
					timestamp = timestamp.AddDate(v.(int), 0, 0)
				}

				if err := d.Set("rfc3339", timestamp.Format(time.RFC3339)); err != nil {
					return nil, fmt.Errorf("error setting rfc3339: %s", err)
				}

				return []*schema.ResourceData{d}, nil
			},
		},

		Schema: map[string]*schema.Schema{
			"base_rfc3339": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ForceNew:     true,
				ValidateFunc: validation.IsRFC3339Time,
			},
			"day": {
				Type:     schema.TypeInt,
				Computed: true,
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
			"offset_days": {
				Type:     schema.TypeInt,
				Optional: true,
				AtLeastOneOf: []string{
					"offset_days",
					"offset_hours",
					"offset_minutes",
					"offset_months",
					"offset_seconds",
					"offset_years",
				},
			},
			"offset_hours": {
				Type:     schema.TypeInt,
				Optional: true,
				AtLeastOneOf: []string{
					"offset_days",
					"offset_hours",
					"offset_minutes",
					"offset_months",
					"offset_seconds",
					"offset_years",
				},
			},
			"offset_minutes": {
				Type:     schema.TypeInt,
				Optional: true,
				AtLeastOneOf: []string{
					"offset_days",
					"offset_hours",
					"offset_minutes",
					"offset_months",
					"offset_seconds",
					"offset_years",
				},
			},
			"offset_months": {
				Type:     schema.TypeInt,
				Optional: true,
				AtLeastOneOf: []string{
					"offset_days",
					"offset_hours",
					"offset_minutes",
					"offset_months",
					"offset_seconds",
					"offset_years",
				},
			},
			"offset_seconds": {
				Type:     schema.TypeInt,
				Optional: true,
				AtLeastOneOf: []string{
					"offset_days",
					"offset_hours",
					"offset_minutes",
					"offset_months",
					"offset_seconds",
					"offset_years",
				},
			},
			"offset_years": {
				Type:     schema.TypeInt,
				Optional: true,
				AtLeastOneOf: []string{
					"offset_days",
					"offset_hours",
					"offset_minutes",
					"offset_months",
					"offset_seconds",
					"offset_years",
				},
			},
			"rfc3339": {
				Type:     schema.TypeString,
				Computed: true,
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

func resourceTimeOffsetCreate(d *schema.ResourceData, m interface{}) error {
	timestamp := time.Now().UTC()

	if v, ok := d.GetOk("base_rfc3339"); ok {
		var err error
		timestamp, err = time.Parse(time.RFC3339, v.(string))

		if err != nil {
			return fmt.Errorf("error parsing base_rfc3339 (%s): %s", v.(string), err)
		}
	}

	d.SetId(timestamp.Format(time.RFC3339))

	if err := d.Set("base_rfc3339", timestamp.Format(time.RFC3339)); err != nil {
		return fmt.Errorf("error setting base_rfc3339: %s", err)
	}

	var offsetTimestamp time.Time

	if v, ok := d.GetOk("offset_days"); ok {
		offsetTimestamp = timestamp.AddDate(0, 0, v.(int))
	}

	if v, ok := d.GetOk("offset_hours"); ok {
		offsetTimestamp = timestamp.Add(time.Duration(v.(int)) * time.Hour)
	}

	if v, ok := d.GetOk("offset_minutes"); ok {
		offsetTimestamp = timestamp.Add(time.Duration(v.(int)) * time.Minute)
	}

	if v, ok := d.GetOk("offset_months"); ok {
		offsetTimestamp = timestamp.AddDate(0, v.(int), 0)
	}

	if v, ok := d.GetOk("offset_seconds"); ok {
		offsetTimestamp = timestamp.Add(time.Duration(v.(int)) * time.Second)
	}

	if v, ok := d.GetOk("offset_years"); ok {
		offsetTimestamp = timestamp.AddDate(v.(int), 0, 0)
	}

	if err := d.Set("rfc3339", offsetTimestamp.Format(time.RFC3339)); err != nil {
		return fmt.Errorf("error setting rfc3339: %s", err)
	}

	return resourceTimeOffsetRead(d, m)
}

func resourceTimeOffsetRead(d *schema.ResourceData, m interface{}) error {
	timestamp, err := time.Parse(time.RFC3339, d.Get("rfc3339").(string))

	if err != nil {
		return fmt.Errorf("error parsing offset timestamp (%s): %s", d.Get("rfc3339").(string), err)
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

func resourceTimeOffsetUpdate(d *schema.ResourceData, m interface{}) error {
	timestamp, err := time.Parse(time.RFC3339, d.Id())

	if err != nil {
		return fmt.Errorf("error parsing timestamp (%s): %s", d.Id(), err)
	}

	var offsetTimestamp time.Time

	if v, ok := d.GetOk("offset_days"); ok {
		offsetTimestamp = timestamp.AddDate(0, 0, v.(int))
	}

	if v, ok := d.GetOk("offset_hours"); ok {
		offsetTimestamp = timestamp.Add(time.Duration(v.(int)) * time.Hour)
	}

	if v, ok := d.GetOk("offset_minutes"); ok {
		offsetTimestamp = timestamp.Add(time.Duration(v.(int)) * time.Minute)
	}

	if v, ok := d.GetOk("offset_months"); ok {
		offsetTimestamp = timestamp.AddDate(0, v.(int), 0)
	}

	if v, ok := d.GetOk("offset_seconds"); ok {
		offsetTimestamp = timestamp.Add(time.Duration(v.(int)) * time.Second)
	}

	if v, ok := d.GetOk("offset_years"); ok {
		offsetTimestamp = timestamp.AddDate(v.(int), 0, 0)
	}

	if err := d.Set("rfc3339", offsetTimestamp.Format(time.RFC3339)); err != nil {
		return fmt.Errorf("error setting rfc3339: %s", err)
	}

	return resourceTimeOffsetRead(d, m)
}

func resourceTimeOffsetConditionExpirationChange(ctx context.Context, diff *schema.ResourceDiff, meta interface{}) bool {
	return diff.HasChange("offset_days") ||
		diff.HasChange("offset_hours") ||
		diff.HasChange("offset_minutes") ||
		diff.HasChange("offset_months") ||
		diff.HasChange("offset_seconds") ||
		diff.HasChange("offset_years")
}
