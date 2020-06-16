package packet

import (
	"net/http"
	"sort"
	"strings"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"

	"github.com/packethost/packngo"
)

// friendlyError improves error messages when the API error is blank or in an
// alternate format (as is the case with invalid token or loadbalancer errors)
func friendlyError(err error) error {
	if e, ok := err.(*packngo.ErrorResponse); ok {
		resp := e.Response
		errors := Errors(e.Errors)

		if 0 == len(errors) {
			errors = Errors{e.SingleError}
		}
		er := &ErrorResponse{
			StatusCode: resp.StatusCode,
			Errors:     errors,
		}
		respHead := resp.Header

		// this checks if the error comes from API (and not from cache/LB)
		if len(errors) > 0 {
			ct := respHead.Get("Content-Type")
			xrid := respHead.Get("X-Request-Id")
			if strings.Contains(ct, "application/json") && len(xrid) > 0 {
				er.IsAPIError = true
			}
		}
		return er
	}
	return err
}

func isForbidden(err error) bool {
	if r, ok := err.(*ErrorResponse); ok {
		return r.StatusCode == http.StatusForbidden
	}
	return false
}

func isNotFound(err error) bool {
	if r, ok := err.(*ErrorResponse); ok {
		return r.StatusCode == http.StatusNotFound && r.IsAPIError
	}
	return false
}

type Errors []string

func (e Errors) Error() string {
	return strings.Join(e, "; ")
}

type ErrorResponse struct {
	StatusCode int
	Errors
	IsAPIError bool
}

// setMap sets the map of values to ResourceData, checking and returning the
// errors. Typically d.Set is not error checked. This helper makes checking
// those errors less tedious. Because this works with a map, the order of the
// errors would not be predictable, to avoid this the errors will be sorted.
func setMap(d *schema.ResourceData, m map[string]interface{}) error {
	errs := &multierror.Error{}
	for key, v := range m {
		var err error
		if f, ok := v.(setFn); ok {
			err = f(d, key)
		} else {
			err = d.Set(key, v)
		}

		if err != nil {
			errs = multierror.Append(errs, err)
		}
	}
	sort.Sort(errs)

	return errs.ErrorOrNil()
}

type setFn = func(d *schema.ResourceData, key string) error
