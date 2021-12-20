package mccpv2

import (
	"net/http"
	"strings"

	"github.com/IBM-Cloud/bluemix-go/client"
	"github.com/IBM-Cloud/bluemix-go/models"
	"github.com/IBM-Cloud/bluemix-go/rest"
)

const mccpEndpointOfPublicBluemix = "https://mccp.ng.bluemix.net"

//go:generate counterfeiter . RegionRepository
type RegionRepository interface {
	PublicRegions() ([]models.Region, error)
	Regions() ([]models.Region, error)
	FindRegionByName(name string) (*models.Region, error)
	FindRegionById(id string) (*models.Region, error)
}

type region struct {
	client *client.Client
}

func newRegionRepositoryAPI(c *client.Client) RegionRepository {
	return &region{
		client: c,
	}
}

func (r *region) PublicRegions() ([]models.Region, error) {
	return r.regions(mccpEndpointOfPublicBluemix)
}

func (r *region) Regions() ([]models.Region, error) {
	return r.regions(*r.client.Config.Endpoint)
}

func (r *region) regions(endpoint string) ([]models.Region, error) {
	var result []models.Region
	resp, err := r.client.SendRequest(rest.GetRequest(endpoint+"/v2/regions"), &result)
	if resp.StatusCode == http.StatusNotFound {
		return []models.Region{}, nil
	}
	if err != nil {
		return []models.Region{}, err
	}
	return result, nil
}

func (r *region) FindRegionByName(name string) (*models.Region, error) {
	regions, err := r.Regions()
	if err != nil {
		return nil, err
	}
	for _, region := range regions {
		if strings.EqualFold(region.Name, name) {
			return &region, nil
		}
	}
	return nil, nil
}
func (r *region) FindRegionById(id string) (*models.Region, error) {
	regions, err := r.Regions()
	if err != nil {
		return nil, err
	}
	for _, region := range regions {
		if strings.EqualFold(region.ID, id) {
			return &region, nil
		}
	}
	return nil, nil
}
