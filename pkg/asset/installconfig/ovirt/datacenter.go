package ovirt

import (
	ovirtsdk4 "github.com/ovirt/go-ovirt"

	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

type dataCenterConfig struct {
	conn         *ovirtsdk4.Connection // ovirtsdk4 connection data
	searchFilter string                // A search filter for datacenters
}

// datacentersAvailable look for all datacenters available in the system.
// Users can provide the filter for the search. Example: status=down
// If search filter not provided, the default filter will be "status=up"
// Returns type: *ovirtsdk.DataCentersServiceListResponse
func datacentersAvailable(d dataCenterConfig) (*ovirtsdk4.DataCentersServiceListResponse, error) {

	if d.searchFilter == "" {
		d.searchFilter = "status=up"
	}
	dcService := d.conn.SystemService().DataCentersService()

	logrus.Debugf("Searching for DataCenters with search filter: %s", d.searchFilter)
	dcResp, err := dcService.List().Search(d.searchFilter).Send()
	if err != nil {
		return nil, errors.Wrapf(err, "Failed to search available DataCenters!")
	}

	return dcResp, nil
}
