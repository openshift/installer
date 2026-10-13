package ocm

import (
	"fmt"
	"net/http"

	v1 "github.com/openshift-online/ocm-sdk-go/accesstransparency/v1"
)

func (c *Client) CreateDecision(accessRequest string, decision string, justification string) error {
	decisionSpec, err := v1.NewDecision().
		Decision(v1.DecisionDecision(decision)).
		Justification(justification).
		Build()
	if err != nil {
		return err
	}
	_, err = c.ocm.AccessTransparency().V1().AccessRequests().
		AccessRequest(accessRequest).Decisions().
		Add().
		Body(decisionSpec).
		Send()
	if err != nil {
		return err
	}
	return nil
}

func (c *Client) GetAccessRequest(id string) (*v1.AccessRequest, bool, error) {
	resp, err := c.ocm.AccessTransparency().V1().AccessRequests().
		AccessRequest(id).Get().Send()
	if resp.Status() == http.StatusNotFound {
		return nil, false, nil
	}
	if err != nil {
		return nil, false, err
	}
	return resp.Body(), true, nil
}

func (c *Client) ListAccessRequest(cluster string) ([]*v1.AccessRequest, error) {
	query := "status.state in ('Pending', 'Approved')"
	order := "status.state desc, updated_at desc"
	if cluster != "" {
		query = fmt.Sprintf("cluster_id='%s'", cluster)
		order = "updated_at desc"
	}
	page := 1
	size := 100
	accessRequests := []*v1.AccessRequest{}
	for {
		resp, err := c.ocm.AccessTransparency().V1().AccessRequests().
			List().
			Search(query).
			Order(order).
			Page(page).
			Size(size).
			Send()
		if err != nil {
			return accessRequests, err
		}
		resp.Items().Each(func(item *v1.AccessRequest) bool {
			accessRequests = append(accessRequests, item)
			return true
		})
		if resp.Size() < size {
			break
		}
		page++
	}
	return accessRequests, nil
}
