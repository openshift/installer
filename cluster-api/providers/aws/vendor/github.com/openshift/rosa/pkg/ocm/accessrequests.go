package ocm

import (
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
