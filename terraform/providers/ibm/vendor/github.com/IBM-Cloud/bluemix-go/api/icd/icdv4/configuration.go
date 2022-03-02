package icdv4

import (
	"fmt"

	"github.com/IBM-Cloud/bluemix-go/client"
	"github.com/IBM-Cloud/bluemix-go/utils"
)

type ConfigurationReq struct {
	Configuration interface{} `json:"configuration"`
}

type Configurations interface {
	UpdateConfiguration(icdId string, configurationReq ConfigurationReq) (Task, error)
	GetConfiguration(icdId string) (interface{}, error)
}

type configurations struct {
	client *client.Client
}

func newConfigurationsAPI(c *client.Client) Configurations {
	return &configurations{
		client: c,
	}
}

func (r *configurations) UpdateConfiguration(icdId string, configurationReq ConfigurationReq) (Task, error) {
	taskResult := TaskResult{}
	rawURL := fmt.Sprintf("/v4/ibm/deployments/%s/configuration", utils.EscapeUrlParm(icdId))
	_, err := r.client.Patch(rawURL, &configurationReq, &taskResult)
	if err != nil {
		return taskResult.Task, err
	}
	return taskResult.Task, nil
}

func (r *configurations) GetConfiguration(icdId string) (interface{}, error) {
	var taskResult interface{}
	rawURL := fmt.Sprintf("/v4/ibm/deployments/%s/configuration/schema", utils.EscapeUrlParm(icdId))
	_, err := r.client.Get(rawURL, &taskResult)
	if err != nil {
		return taskResult, err
	}
	return taskResult, nil
}
