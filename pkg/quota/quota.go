package quota

import (
	"errors"
	"fmt"
	"math"
	"strings"
)

// Quota stores the used and limit for a resouce's quota.
type Quota struct {
	Service string
	Name    string
	Region  string

	InUse int64
	Limit int64

	Unlimited bool
}

// Constraint defines a check against availablity
// for a resource quota.
type Constraint struct {
	Name string
	// This should be global or specific region.
	Region string
	// Count is the number of the resource that is required
	// to be free for use.
	Count int64
}

// ConstraintReportResult provide one word result for the constraint
type ConstraintReportResult string

const (
	// Available corresponds to a constraint that is available.
	Available ConstraintReportResult = "Available"
	// AvailableButLow corresponds to a constraint that is available but there is very low overhead after satisfying the constraint.
	AvailableButLow ConstraintReportResult = "AvailableButLow"
	// NotAvailable corresponds to a constraint that is not available.
	NotAvailable ConstraintReportResult = "NotAvailable"
	// Unknown corresponds to a constraint that could not be calculated.
	Unknown ConstraintReportResult = "Unknown"
)

// ConstraintReport provides result for a given constraint.
type ConstraintReport struct {
	For     *Constraint
	Result  ConstraintReportResult
	Message string
}

// Check returns whether the checks constraints are possible gives the quotas.
// It returns an erros when any one of the constraint's result was `NotAvailable` or `Unknown`
func Check(quotas []Quota, checks []Constraint) ([]ConstraintReport, error) {
	reports := make([]ConstraintReport, 0, len(checks))

	match := func(check Constraint) (Quota, bool) {
		for _, q := range quotas {
			if !strings.EqualFold(check.Name, q.Name) {
				continue
			}
			if !strings.EqualFold(q.Region, check.Region) {
				continue
			}
			return q, true
		}
		return Quota{}, false
	}

	for idx, check := range checks {
		report := ConstraintReport{
			For: &checks[idx],
		}
		matched, ok := match(check)
		if !ok {
			report.Result = Unknown
			report.Message = "No matching quota found for the constraint"
			reports = append(reports, report)
			continue
		}

		if matched.Unlimited {
			report.Result = Available
			report.Message = "Unlimited quota found for the constraint"
			reports = append(reports, report)
			continue
		}

		if check.Count > matched.Limit {
			report.Result = NotAvailable
			report.Message = fmt.Sprintf("the required number of resources (%d) is more than the limit of %d", check.Count, matched.Limit)
			reports = append(reports, report)
			continue
		}

		avail := matched.Limit - matched.InUse
		availAfterUse := avail - check.Count
		headroom := int64(math.Ceil(0.2 * float64(matched.Limit)))
		if check.Count > avail {
			report.Result = NotAvailable
			report.Message = fmt.Sprintf("the required number of resources (%d) is more than remaining quota of %d", check.Count, avail)
			reports = append(reports, report)
			continue
		}
		if availAfterUse <= headroom {
			report.Result = AvailableButLow
			report.Message = fmt.Sprintf("the required number of resources is available but only %d will be leftover", availAfterUse)
			reports = append(reports, report)
			continue
		}
		report.Result = Available
		report.Message = "the required number of resources is available"
		reports = append(reports, report)
	}

	failed := 0
	for _, r := range reports {
		if r.Result == Unknown ||
			r.Result == NotAvailable {
			failed++
		}
	}
	if failed > 0 {
		return reports, errors.New("%d checks failed")
	}

	return reports, nil
}
