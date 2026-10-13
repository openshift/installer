// go mod won't pull in code that isn't depended upon, but we have some code we don't depend on from code that must be included
// for our build to work.
package dependencymagnet

import (
	// this makes it easy to vendor the entire model.  If you add another dir, be sure to add it here
	_ "github.com/openshift-online/ocm-api-model/model/access_transparency/v1"
	_ "github.com/openshift-online/ocm-api-model/model/accounts_mgmt/v1"
	_ "github.com/openshift-online/ocm-api-model/model/addons_mgmt/v1"
	_ "github.com/openshift-online/ocm-api-model/model/aro_hcp/v1alpha1"
	_ "github.com/openshift-online/ocm-api-model/model/authorizations/v1"
	_ "github.com/openshift-online/ocm-api-model/model/clusters_mgmt/v1"
	_ "github.com/openshift-online/ocm-api-model/model/job_queue/v1"
	_ "github.com/openshift-online/ocm-api-model/model/osd_fleet_mgmt/v1"
	_ "github.com/openshift-online/ocm-api-model/model/service_logs/v1"
	_ "github.com/openshift-online/ocm-api-model/model/service_mgmt/v1"
	_ "github.com/openshift-online/ocm-api-model/model/status_board/v1"
	_ "github.com/openshift-online/ocm-api-model/model/web_rca/v1"
)
