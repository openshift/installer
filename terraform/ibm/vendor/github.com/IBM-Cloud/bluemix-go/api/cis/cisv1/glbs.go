package cisv1

import (
	"fmt"
	"time"

	"github.com/IBM-Cloud/bluemix-go/client"
)

type Glb struct {
	Id              string              `json:"id"`
	Name            string              `json:"name"`
	Desc            string              `json:"description"`
	FallbackPool    string              `json:"fallback_pool"`
	DefaultPools    []string            `json:"default_pools"`
	Ttl             int                 `json:"ttl"`
	Proxied         bool                `json:"proxied"`
	CreatedOn       *time.Time          `json:"created_on,omitempty"`
	ModifiedOn      *time.Time          `json:"modified_on,omitempty"`
	SessionAffinity string              `json:"session_affinity"`
	Enabled         bool                `json:"enabled,omitempty"`
	RegionPools     map[string][]string `json:"region_pools,omitempty"`
	PopPools        map[string][]string `json:"pop_pools,omitempty"`
}

type GlbResults struct {
	GlbList     []Glb        `json:"result"`
	ResultsInfo ResultsCount `json:"result_info"`
	Success     bool         `json:"success"`
	Errors      []Error      `json:"errors"`
}

type GlbResult struct {
	Glb      Glb      `json:"result"`
	Success  bool     `json:"success"`
	Errors   []Error  `json:"errors"`
	Messages []string `json:"messages"`
}

type GlbBody struct {
	Desc            string              `json:"description,omitempty"`
	Proxied         bool                `json:"proxied,omitempty"`
	Name            string              `json:"name"`
	FallbackPool    string              `json:"fallback_pool"`
	DefaultPools    []string            `json:"default_pools"`
	SessionAffinity string              `json:"session_affinity,omitempty"`
	Ttl             int                 `json:"ttl,omitempty"`
	Enabled         bool                `json:"enabled,omitempty"`
	RegionPools     map[string][]string `json:"region_pools,omitempty"`
	PopPools        map[string][]string `json:"pop_pools,omitempty"`
}

type GlbDelete struct {
	Result struct {
		GlbId string
	} `json:"result"`
	Success  bool     `json:"success"`
	Errors   []Error  `json:"errors"`
	Messages []string `json:"messages"`
}

type Glbs interface {
	ListGlbs(cisId string, zoneId string) ([]Glb, error)
	GetGlb(cisId string, zoneId string, glbId string) (*Glb, error)
	CreateGlb(cisId string, zoneId string, glbBody GlbBody) (*Glb, error)
	DeleteGlb(cisId string, zoneId string, glbId string) error
	UpdateGlb(cisId string, zoneId string, glbId string, glbBody GlbBody) (*Glb, error)
}

type glbs struct {
	client *client.Client
}

func newGlbAPI(c *client.Client) Glbs {
	return &glbs{
		client: c,
	}
}

func (r *glbs) ListGlbs(cisId string, zoneId string) ([]Glb, error) {
	glbResults := GlbResults{}
	rawURL := fmt.Sprintf("/v1/%s/zones/%s/load_balancers", cisId, zoneId)
	_, err := r.client.Get(rawURL, &glbResults)
	if err != nil {
		return nil, err
	}
	return glbResults.GlbList, err
}

func (r *glbs) GetGlb(cisId string, zoneId string, glbId string) (*Glb, error) {
	glbResult := GlbResult{}
	rawURL := fmt.Sprintf("/v1/%s/zones/%s/load_balancers/%s", cisId, zoneId, glbId)
	_, err := r.client.Get(rawURL, &glbResult, nil)
	if err != nil {
		return nil, err
	}
	return &glbResult.Glb, nil
}

func (r *glbs) DeleteGlb(cisId string, zoneId string, glbId string) error {
	rawURL := fmt.Sprintf("/v1/%s/zones/%s/load_balancers/%s", cisId, zoneId, glbId)
	_, err := r.client.Delete(rawURL)
	if err != nil {
		return err
	}
	return nil
}

func (r *glbs) CreateGlb(cisId string, zoneId string, glbBody GlbBody) (*Glb, error) {
	glbResult := GlbResult{}
	rawURL := fmt.Sprintf("/v1/%s/zones/%s/load_balancers", cisId, zoneId)
	_, err := r.client.Post(rawURL, &glbBody, &glbResult)
	if err != nil {
		return nil, err
	}
	return &glbResult.Glb, nil
}

func (r *glbs) UpdateGlb(cisId string, zoneId string, glbId string, glbBody GlbBody) (*Glb, error) {
	glbResult := GlbResult{}
	rawURL := fmt.Sprintf("/v1/%s/zones/%s/load_balancers/%s", cisId, zoneId, glbId)
	_, err := r.client.Put(rawURL, &glbBody, &glbResult)
	if err != nil {
		return nil, err
	}
	return &glbResult.Glb, nil
}
