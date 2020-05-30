package ovirt

import (
	ovirtsdk4 "github.com/ovirt/go-ovirt"

	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

// datacentersAvailable look for all datacenters available in the system.
// Users can provide the criteria for the search. Example: status=up
// with the status criteria as UP. Returns type: *ovirtsdk.DataCentersServiceListResponse
func datacentersAvailable(c *ovirtsdk4.Connection, criteria string) (*ovirtsdk4.DataCentersServiceListResponse, error) {

	searchCriteria := criteria
	dcService := c.SystemService().DataCentersService()

	logrus.Debugf("Searching for DataCenters with search criteria: %s", criteria)
	dcResp, err := dcService.List().Search(searchCriteria).Send()
	if err != nil {
		return nil, errors.Wrapf(err, "Failed to search available DataCenters!")
	}

	return dcResp, nil
}
