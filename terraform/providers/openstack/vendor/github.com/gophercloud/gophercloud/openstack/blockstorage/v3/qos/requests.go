package qos

import (
	"github.com/gophercloud/gophercloud"
	"github.com/gophercloud/gophercloud/pagination"
)

type CreateOptsBuilder interface {
	ToQoSCreateMap() (map[string]interface{}, error)
}

// ListOptsBuilder allows extensions to add additional parameters to the
// List request.
type ListOptsBuilder interface {
	ToQoSListQuery() (string, error)
}

type QoSConsumer string

const (
	ConsumerFront QoSConsumer = "front-end"
	ConsumerBack  QoSConsumer = "back-end"
	ConsumerBoth  QoSConsumer = "both"
)

// CreateOpts contains options for creating a QoS specification.
// This object is passed to the qos.Create function.
type CreateOpts struct {
	// The name of the QoS spec
	Name string `json:"name"`
	// The consumer of the QoS spec. Possible values are
	// both, front-end, back-end.
	Consumer QoSConsumer `json:"consumer,omitempty"`
	// Specs is a collection of miscellaneous key/values used to set
	// specifications for the QoS
	Specs map[string]string `json:"-"`
}

// ToQoSCreateMap assembles a request body based on the contents of a
// CreateOpts.
func (opts CreateOpts) ToQoSCreateMap() (map[string]interface{}, error) {
	b, err := gophercloud.BuildRequestBody(opts, "qos_specs")
	if err != nil {
		return nil, err
	}

	if opts.Specs != nil {
		if v, ok := b["qos_specs"].(map[string]interface{}); ok {
			for key, value := range opts.Specs {
				v[key] = value
			}
		}
	}

	return b, nil
}

// Create will create a new QoS based on the values in CreateOpts. To extract
// the QoS object from the response, call the Extract method on the
// CreateResult.
func Create(client *gophercloud.ServiceClient, opts CreateOptsBuilder) (r CreateResult) {
	b, err := opts.ToQoSCreateMap()
	if err != nil {
		r.Err = err
		return
	}
	resp, err := client.Post(createURL(client), b, &r.Body, &gophercloud.RequestOpts{
		OkCodes: []int{200},
	})
	_, r.Header, r.Err = gophercloud.ParseResponse(resp, err)
	return
}

// DeleteOptsBuilder allows extensions to add additional parameters to the
// Delete request.
type DeleteOptsBuilder interface {
	ToQoSDeleteQuery() (string, error)
}

// DeleteOpts contains options for deleting a QoS. This object is passed to
// the qos.Delete function.
type DeleteOpts struct {
	// Delete a QoS specification even if it is in-use
	Force bool `q:"force"`
}

// ToQoSDeleteQuery formats a DeleteOpts into a query string.
func (opts DeleteOpts) ToQoSDeleteQuery() (string, error) {
	q, err := gophercloud.BuildQueryString(opts)
	return q.String(), err
}

// Delete will delete the existing QoS with the provided ID.
func Delete(client *gophercloud.ServiceClient, id string, opts DeleteOptsBuilder) (r DeleteResult) {
	url := deleteURL(client, id)
	if opts != nil {
		query, err := opts.ToQoSDeleteQuery()
		if err != nil {
			r.Err = err
			return
		}
		url += query
	}
	resp, err := client.Delete(url, nil)
	_, r.Header, r.Err = gophercloud.ParseResponse(resp, err)
	return
}

type ListOpts struct {
	// Sort is Comma-separated list of sort keys and optional sort
	// directions in the form of < key > [: < direction > ]. A valid
	//direction is asc (ascending) or desc (descending).
	Sort string `q:"sort"`

	// Marker and Limit control paging.
	// Marker instructs List where to start listing from.
	Marker string `q:"marker"`

	// Limit instructs List to refrain from sending excessively large lists of
	// QoS.
	Limit int `q:"limit"`
}

// ToQoSListQuery formats a ListOpts into a query string.
func (opts ListOpts) ToQoSListQuery() (string, error) {
	q, err := gophercloud.BuildQueryString(opts)
	return q.String(), err
}

// List instructs OpenStack to provide a list of QoS.
// You may provide criteria by which List curtails its results for easier
// processing.
func List(client *gophercloud.ServiceClient, opts ListOptsBuilder) pagination.Pager {
	url := listURL(client)
	if opts != nil {
		query, err := opts.ToQoSListQuery()
		if err != nil {
			return pagination.Pager{Err: err}
		}
		url += query
	}
	return pagination.NewPager(client, url, func(r pagination.PageResult) pagination.Page {
		return QoSPage{pagination.LinkedPageBase{PageResult: r}}
	})
}

// Get retrieves details of a single qos. Use Extract to convert its
// result into a QoS.
func Get(client *gophercloud.ServiceClient, id string) (r GetResult) {
	resp, err := client.Get(getURL(client, id), &r.Body, &gophercloud.RequestOpts{
		OkCodes: []int{200},
	})
	_, r.Header, r.Err = gophercloud.ParseResponse(resp, err)
	return
}

// CreateQosSpecsOptsBuilder allows extensions to add additional parameters to the
// CreateQosSpecs requests.
type CreateQosSpecsOptsBuilder interface {
	ToQosSpecsCreateMap() (map[string]interface{}, error)
}

// UpdateOpts contains options for creating a QoS specification.
// This object is passed to the qos.Update function.
type UpdateOpts struct {
	// The consumer of the QoS spec. Possible values are
	// both, front-end, back-end.
	Consumer QoSConsumer `json:"consumer,omitempty"`
	// Specs is a collection of miscellaneous key/values used to set
	// specifications for the QoS
	Specs map[string]string `json:"-"`
}

type UpdateOptsBuilder interface {
	ToQoSUpdateMap() (map[string]interface{}, error)
}

// ToQoSUpdateMap assembles a request body based on the contents of a
// UpdateOpts.
func (opts UpdateOpts) ToQoSUpdateMap() (map[string]interface{}, error) {
	b, err := gophercloud.BuildRequestBody(opts, "qos_specs")
	if err != nil {
		return nil, err
	}

	if opts.Specs != nil {
		if v, ok := b["qos_specs"].(map[string]interface{}); ok {
			for key, value := range opts.Specs {
				v[key] = value
			}
		}
	}

	return b, nil
}

// Update will update an existing QoS based on the values in UpdateOpts.
// To extract the QoS object from the response, call the Extract method
// on the UpdateResult.
func Update(client *gophercloud.ServiceClient, id string, opts UpdateOptsBuilder) (r updateResult) {
	b, err := opts.ToQoSUpdateMap()
	if err != nil {
		r.Err = err
		return
	}
	resp, err := client.Put(updateURL(client, id), b, &r.Body, &gophercloud.RequestOpts{
		OkCodes: []int{200},
	})
	_, r.Header, r.Err = gophercloud.ParseResponse(resp, err)
	return
}

// DeleteKeysOptsBuilder allows extensions to add additional parameters to the
// CreateExtraSpecs requests.
type DeleteKeysOptsBuilder interface {
	ToDeleteKeysCreateMap() (map[string]interface{}, error)
}

// DeleteKeysOpts is a string slice that contains keys to be deleted.
type DeleteKeysOpts []string

// ToDeleteKeysCreateMap assembles a body for a Create request based on
// the contents of ExtraSpecsOpts.
func (opts DeleteKeysOpts) ToDeleteKeysCreateMap() (map[string]interface{}, error) {
	return map[string]interface{}{"keys": opts}, nil
}

// DeleteKeys will delete the keys/specs from the specified QoS
func DeleteKeys(client *gophercloud.ServiceClient, qosID string, opts DeleteKeysOptsBuilder) (r DeleteResult) {
	b, err := opts.ToDeleteKeysCreateMap()
	if err != nil {
		r.Err = err
		return
	}
	resp, err := client.Put(deleteKeysURL(client, qosID), b, nil, &gophercloud.RequestOpts{
		OkCodes: []int{202},
	})
	_, r.Header, r.Err = gophercloud.ParseResponse(resp, err)
	return
}

// AssociateOpitsBuilder allows extensions to define volume type id
// to the associate query
type AssociateOptsBuilder interface {
	ToQosAssociateQuery() (string, error)
}

// AssociateOpts contains options for associating a QoS with a
// volume type
type AssociateOpts struct {
	VolumeTypeID string `q:"vol_type_id" required:"true"`
}

// ToQosAssociateQuery formats an AssociateOpts into a query string
func (opts AssociateOpts) ToQosAssociateQuery() (string, error) {
	q, err := gophercloud.BuildQueryString(opts)
	return q.String(), err
}

// Associate will associate a qos with a volute type
func Associate(client *gophercloud.ServiceClient, qosID string, opts AssociateOptsBuilder) (r AssociateResult) {
	url := associateURL(client, qosID)
	query, err := opts.ToQosAssociateQuery()
	if err != nil {
		r.Err = err
		return
	}
	url += query

	resp, err := client.Get(url, nil, &gophercloud.RequestOpts{
		OkCodes: []int{202},
	})
	_, r.Header, r.Err = gophercloud.ParseResponse(resp, err)
	return
}

// DisassociateOpitsBuilder allows extensions to define volume type id
// to the disassociate query
type DisassociateOptsBuilder interface {
	ToQosDisassociateQuery() (string, error)
}

// DisassociateOpts contains options for disassociating a QoS from a
// volume type
type DisassociateOpts struct {
	VolumeTypeID string `q:"vol_type_id" required:"true"`
}

// ToQosDisassociateQuery formats a DisassociateOpts into a query string
func (opts DisassociateOpts) ToQosDisassociateQuery() (string, error) {
	q, err := gophercloud.BuildQueryString(opts)
	return q.String(), err
}

// Disassociate will disassociate a qos from a volute type
func Disassociate(client *gophercloud.ServiceClient, qosID string, opts DisassociateOptsBuilder) (r DisassociateResult) {
	url := disassociateURL(client, qosID)
	query, err := opts.ToQosDisassociateQuery()
	if err != nil {
		r.Err = err
		return
	}
	url += query

	resp, err := client.Get(url, nil, &gophercloud.RequestOpts{
		OkCodes: []int{202},
	})
	_, r.Header, r.Err = gophercloud.ParseResponse(resp, err)
	return
}

// DisassociateAll will disassociate a qos from all volute types
func DisassociateAll(client *gophercloud.ServiceClient, qosID string) (r DisassociateAllResult) {
	resp, err := client.Get(disassociateAllURL(client, qosID), nil, &gophercloud.RequestOpts{
		OkCodes: []int{202},
	})
	_, r.Header, r.Err = gophercloud.ParseResponse(resp, err)
	return
}

// ListAssociations retrieves the associations of a QoS.
func ListAssociations(client *gophercloud.ServiceClient, qosID string) pagination.Pager {
	url := listAssociationsURL(client, qosID)

	return pagination.NewPager(client, url, func(r pagination.PageResult) pagination.Page {
		return AssociationPage{pagination.SinglePageBase(r)}
	})
}
