package cisv1

import (
	"fmt"

	"github.com/IBM-Cloud/bluemix-go/client"
)

type Pool struct {
	Id           string   `json:"id"`
	Description  string   `json:"description"`
	Name         string   `json:"name"`
	CheckRegions []string `json:"check_regions"`
	Enabled      bool     `json:"enabled"`
	MinOrigins   int      `json:"minimum_origins"`
	Monitor      string   `json:"monitor"`
	NotEmail     string   `json:"notification_email"`
	Origins      []Origin `json:"origins"`
	Health       string   `json:"health"`
	CreatedOn    string   `json:"created_on"`
	ModifiedOn   string   `json:"modified_on"`
}

type CheckRegion struct {
	Region string `json:"0"`
}

type Origin struct {
	Name    string `json:"name"`
	Address string `json:"address"`
	Enabled bool   `json:"enabled"`
	Weight  int    `json:"weight"`
	Healthy bool   `json:"healthy"`
}

type PoolResults struct {
	PoolList    []Pool       `json:"result"`
	ResultsInfo ResultsCount `json:"result_info"`
	Success     bool         `json:"success"`
	Errors      []Error      `json:"errors"`
}

type PoolResult struct {
	Pool     Pool     `json:"result"`
	Success  bool     `json:"success"`
	Errors   []Error  `json:"errors"`
	Messages []string `json:"messages"`
}

type PoolBody struct {
	Name         string   `json:"name"`
	Description  string   `json:"description,omitempty"`
	Origins      []Origin `json:"origins"`
	CheckRegions []string `json:"check_regions"`
	Enabled      bool     `json:"enabled"`
	MinOrigins   int      `json:"minimum_origins,omitempty"`
	Monitor      string   `json:"monitor,omitempty"`
	NotEmail     string   `json:"notification_email,omitempty"`
}

type PoolDelete struct {
	Result struct {
		PoolId string
	} `json:"result"`
	Success  bool     `json:"success"`
	Errors   []Error  `json:"errors"`
	Messages []string `json:"messages"`
}

type Pools interface {
	ListPools(cisId string) ([]Pool, error)
	GetPool(cisId string, poolId string) (*Pool, error)
	CreatePool(cisId string, poolBody PoolBody) (*Pool, error)
	DeletePool(cisId string, poolId string) error
	UpdatePool(cisId string, poolId string, poolBody PoolBody) (*Pool, error)
}

type pools struct {
	client *client.Client
}

func newPoolAPI(c *client.Client) Pools {
	return &pools{
		client: c,
	}
}

func (r *pools) ListPools(cisId string) ([]Pool, error) {
	poolResults := PoolResults{}
	rawURL := fmt.Sprintf("/v1/%s/load_balancers/pools/", cisId)
	_, err := r.client.Get(rawURL, &poolResults)
	if err != nil {
		return nil, err
	}
	return poolResults.PoolList, err
}

func (r *pools) GetPool(cisId string, poolId string) (*Pool, error) {
	poolResult := PoolResult{}
	rawURL := fmt.Sprintf("/v1/%s/load_balancers/pools/%s", cisId, poolId)
	_, err := r.client.Get(rawURL, &poolResult, nil)
	if err != nil {
		return nil, err
	}
	return &poolResult.Pool, nil
}

func (r *pools) DeletePool(cisId string, poolId string) error {
	rawURL := fmt.Sprintf("/v1/%s/load_balancers/pools/%s", cisId, poolId)
	_, err := r.client.Delete(rawURL)
	if err != nil {
		return err
	}
	return nil
}

func (r *pools) CreatePool(cisId string, poolBody PoolBody) (*Pool, error) {
	poolResult := PoolResult{}
	rawURL := fmt.Sprintf("/v1/%s/load_balancers/pools/", cisId)
	_, err := r.client.Post(rawURL, &poolBody, &poolResult)
	if err != nil {
		return nil, err
	}
	return &poolResult.Pool, nil
}

func (r *pools) UpdatePool(cisId string, poolId string, poolBody PoolBody) (*Pool, error) {
	poolResult := PoolResult{}
	rawURL := fmt.Sprintf("/v1/%s/load_balancers/pools/%s", cisId, poolId)
	_, err := r.client.Put(rawURL, &poolBody, &poolResult)
	if err != nil {
		return nil, err
	}
	return &poolResult.Pool, nil
}
