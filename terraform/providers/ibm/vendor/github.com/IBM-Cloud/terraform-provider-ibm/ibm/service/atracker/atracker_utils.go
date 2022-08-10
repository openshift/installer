package atracker

import (
	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/conns"
	"github.com/IBM/platform-services-go-sdk/atrackerv1"
	"github.com/IBM/platform-services-go-sdk/atrackerv2"
)

const (
	REDACTED_TEXT       = "REDACTED"
	BLOCKED_V1_RESOURCE = "v2_resource_exists_v1_api_is_not_accessible"
)

func getAtrackerClients(meta interface{}) (
	atrackerClientv1 *atrackerv1.AtrackerV1, atrackerClientv2 *atrackerv2.AtrackerV2, err error) {
	atrackerClientv1, err = meta.(conns.ClientSession).AtrackerV1()
	if err != nil {
		return
	}
	atrackerClientv2, err = meta.(conns.ClientSession).AtrackerV2()
	if err != nil {
		return
	}

	_, err = meta.(conns.ClientSession).BluemixSession()
	if err != nil {
		return
	}

	return atrackerClientv1, atrackerClientv2, nil
}
