package icdv4

import (
	"fmt"
	"github.com/IBM-Cloud/bluemix-go/client"
	"github.com/IBM-Cloud/bluemix-go/utils"
)

type Whitelist struct {
	WhitelistEntrys []WhitelistEntry `json:"ip_addresses"`
}

type WhitelistEntry struct {
	Address     string `json:"address,omitempty"`
	Description string `json:"description,omitempty"`
}

type WhitelistReq struct {
	WhitelistEntry WhitelistEntry `json:"ip_address"`
}

type Whitelists interface {
	CreateWhitelist(icdId string, whitelistReq WhitelistReq) (Task, error)
	GetWhitelist(icdId string) (Whitelist, error)
	DeleteWhitelist(icdId string, ipAddress string) (Task, error)
}

type whitelists struct {
	client *client.Client
}

func newWhitelistAPI(c *client.Client) Whitelists {
	return &whitelists{
		client: c,
	}
}

func (r *whitelists) CreateWhitelist(icdId string, whitelistReq WhitelistReq) (Task, error) {
	taskResult := TaskResult{}
	rawURL := fmt.Sprintf("/v4/ibm/deployments/%s/whitelists/ip_addresses", utils.EscapeUrlParm(icdId))
	_, err := r.client.Post(rawURL, &whitelistReq, &taskResult)
	if err != nil {
		return taskResult.Task, err
	}
	return taskResult.Task, nil
}

func (r *whitelists) GetWhitelist(icdId string) (Whitelist, error) {
	whitelist := Whitelist{}
	rawURL := fmt.Sprintf("/v4/ibm/deployments/%s/whitelists/ip_addresses", utils.EscapeUrlParm(icdId))
	_, err := r.client.Get(rawURL, &whitelist)
	if err != nil {
		return whitelist, err
	}
	return whitelist, nil
}

func (r *whitelists) DeleteWhitelist(icdId string, ipAddress string) (Task, error) {
	taskResult := TaskResult{}
	rawURL := fmt.Sprintf("/v4/ibm/deployments/%s/whitelists/ip_addresses/%s", utils.EscapeUrlParm(icdId), utils.EscapeUrlParm(ipAddress))
	_, err := r.client.DeleteWithResp(rawURL, &taskResult)
	if err != nil {
		return taskResult.Task, err
	}
	return taskResult.Task, nil
}
