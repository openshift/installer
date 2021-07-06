package cisv1

import (
	"fmt"
	"github.com/IBM-Cloud/bluemix-go/client"
)

type SettingsResult struct {
	Result   SettingsResObj `json:"result"`
	Success  bool           `json:"success"`
	Errors   []Error        `json:"errors"`
	Messages []string       `json:"messages"`
}

type SettingsResObj struct {
	Id                string `json:"id"`
	Value             string `json:"value"`
	Editable          bool   `json:"editable"`
	ModifiedDate      string `json:"modified_on"`
	CertificateStatus string `json:"certificate_status"`
}

type SettingsBody struct {
	Value string `json:"value"`
}

type Settings interface {
	GetSetting(cisId string, zoneId string, setting string) (*SettingsResObj, error)
	UpdateSetting(cisId string, zoneId string, setting string, settingsBody SettingsBody) (*SettingsResObj, error)
}

type settings struct {
	client *client.Client
}

func newSettingsAPI(c *client.Client) Settings {
	return &settings{
		client: c,
	}
}

func (r *settings) GetSetting(cisId string, zoneId string, setting string) (*SettingsResObj, error) {
	settingsResult := SettingsResult{}
	rawURL := fmt.Sprintf("/v1/%s/zones/%s/settings/%s", cisId, zoneId, setting)
	_, err := r.client.Get(rawURL, &settingsResult)
	if err != nil {
		return nil, err
	}
	return &settingsResult.Result, nil
}

func (r *settings) UpdateSetting(cisId string, zoneId string, setting string, settingsBody SettingsBody) (*SettingsResObj, error) {
	settingsResult := SettingsResult{}
	rawURL := fmt.Sprintf("/v1/%s/zones/%s/settings/%s", cisId, zoneId, setting)
	_, err := r.client.Patch(rawURL, &settingsBody, &settingsResult)
	if err != nil {
		return nil, err
	}
	return &settingsResult.Result, nil
}
