package cisv1

import (
	"fmt"

	"github.com/IBM-Cloud/bluemix-go/client"
)

type ResultsCount struct {
	Count int `json:"count"`
}

type Error struct {
	Code int    `json:"code"`
	Msg  string `json:"message"`
}

type NameServer struct {
	NameS int64 `json:"0"`
}

type Zone struct {
	Id                 string   `json:"id"`
	Name               string   `json:"name"`
	Status             string   `json:"status"`
	Paused             bool     `json:"paused"`
	NameServers        []string `json:"name_servers"`
	OriginalNameServer []string `json:"original_name_servers"`
}

type ZoneResults struct {
	ZoneList    []Zone       `json:"result"`
	ResultsInfo ResultsCount `json:"result_info"`
	Success     bool         `json:"success"`
	Errors      []Error      `json:"errors"`
}

type ZoneResult struct {
	Zone     Zone     `json:"result"`
	Success  bool     `json:"success"`
	Errors   []Error  `json:"errors"`
	Messages []string `json:"messages"`
}

type ZoneBody struct {
	Name string `json:"name"`
}

type ZoneDelete struct {
	Result struct {
		ZoneId string
	} `json:"result"`
	Success  bool     `json:"success"`
	Errors   []Error  `json:"errors"`
	Messages []string `json:"messages"`
}

type Zones interface {
	ListZones(cisId string) ([]Zone, error)
	GetZone(cisId string, zoneId string) (*Zone, error)
	CreateZone(cisId string, zoneBody ZoneBody) (*Zone, error)
	DeleteZone(cisId string, zoneId string) error
}

type zones struct {
	client *client.Client
}

func newZoneAPI(c *client.Client) Zones {
	return &zones{
		client: c,
	}
}

func (r *zones) ListZones(cisId string) ([]Zone, error) {
	zoneResults := ZoneResults{}
	rawURL := fmt.Sprintf("/v1/%s/zones?page=1", cisId)
	if _, err := r.client.GetPaginated(rawURL, NewDNSPaginatedResources(Zone{}), func(resource interface{}) bool {
		if zone, ok := resource.(Zone); ok {
			zoneResults.ZoneList = append(zoneResults.ZoneList, zone)
			return true
		}
		return false
	}); err != nil {
		return nil, fmt.Errorf("failed to list paginated dns records: %s", err)
	}
	return zoneResults.ZoneList, nil
}
func (r *zones) GetZone(cisId string, zoneId string) (*Zone, error) {
	zoneResult := ZoneResult{}
	rawURL := fmt.Sprintf("/v1/%s/zones/%s", cisId, zoneId)
	_, err := r.client.Get(rawURL, &zoneResult, nil)
	if err != nil {
		return nil, err
	}
	return &zoneResult.Zone, nil
}

func (r *zones) DeleteZone(cisId string, zoneId string) error {
	rawURL := fmt.Sprintf("/v1/%s/zones/%s", cisId, zoneId)
	_, err := r.client.Delete(rawURL)
	if err != nil {
		return err
	}
	return nil
}

func (r *zones) CreateZone(cisId string, zoneBody ZoneBody) (*Zone, error) {
	zoneResult := ZoneResult{}
	rawURL := fmt.Sprintf("/v1/%s/zones/", cisId)
	_, err := r.client.Post(rawURL, &zoneBody, &zoneResult)
	if err != nil {
		return nil, err
	}
	return &zoneResult.Zone, nil
}
