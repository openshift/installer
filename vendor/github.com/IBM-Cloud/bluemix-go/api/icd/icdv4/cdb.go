package icdv4

import (
	"fmt"

	"github.com/IBM-Cloud/bluemix-go/client"
	"github.com/IBM-Cloud/bluemix-go/utils"
)

type CdbResult struct {
	Cdb Cdb `json:"deployment"`
}

type Cdb struct {
	Id              string          `json:"id"`
	Name            string          `json:"name"`
	Type            string          `json:"type"`
	PlatformOptions PlatformOptions `json:"platform_options"`
	Version         string          `json:"version"`
	AdminUser       string          `json:"admin_username"`
}

type PlatformOptions struct {
	KeyProtectKey          string `json:"key_protect_key_id"`
	DiskENcryptionKeyCrn   string `json:"disk_encryption_key_crn"`
	BackUpEncryptionKeyCrn string `json:"backup_encryption_key_crn"`
}

type Cdbs interface {
	GetCdb(icdId string) (Cdb, error)
}

type cdbs struct {
	client *client.Client
}

func newCdbAPI(c *client.Client) Cdbs {
	return &cdbs{
		client: c,
	}
}

func (r *cdbs) GetCdb(icdId string) (Cdb, error) {
	cdbResult := CdbResult{}
	rawURL := fmt.Sprintf("/v4/ibm/deployments/%s", utils.EscapeUrlParm(icdId))
	_, err := r.client.Get(rawURL, &cdbResult)
	if err != nil {
		return cdbResult.Cdb, err
	}
	return cdbResult.Cdb, nil
}
