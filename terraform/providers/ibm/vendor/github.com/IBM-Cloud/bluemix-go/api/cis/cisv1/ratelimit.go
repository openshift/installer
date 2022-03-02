package cisv1

import (
	"fmt"

	"github.com/IBM-Cloud/bluemix-go/client"
)

// RateLimitRecord is a policy than can be applied to limit traffic within a customer domain
type RateLimitRecord struct {
	ID          string              `json:"id,omitempty"`
	Disabled    bool                `json:"disabled,omitempty"`
	Description string              `json:"description,omitempty"`
	Bypass      []RateLimitByPass   `json:"bypass,omitempty"`
	Threshold   int                 `json:"threshold"`
	Period      int                 `json:"period"`
	Correlate   *RateLimitCorrelate `json:"correlate,omitempty"`
	Action      RateLimitAction     `json:"action"`
	Match       RateLimitMatch      `json:"match"`
}

//  RateLimitByPass ...
type RateLimitByPass struct {
	Name  string `json:"name"`
	Value string `json:"value"`
}

// RateLimitCorrelate ...
type RateLimitCorrelate struct {
	By string `json:"by"`
}

// RateLimitAction ...
type RateLimitAction struct {
	Mode     string          `json:"mode"`
	Timeout  int             `json:"timeout,omitempty"`
	Response *ActionResponse `json:"response,omitempty"`
}

// ActionResponse ...
type ActionResponse struct {
	ContentType string `json:"content_type,omitempty"`
	Body        string `json:"body,omitempty"`
}

// RateLimitMatch ...
type RateLimitMatch struct {
	Request  MatchRequest  `json:"request"`
	Response MatchResponse `json:"response"`
}

// MatchRequest ...
type MatchRequest struct {
	Methods []string `json:"methods,omitempty"`
	Schemes []string `json:"schemes,omitempty"`
	URL     string   `json:"url,omitempty"`
}

// MatchResponse ...
type MatchResponse struct {
	Statuses      []int                 `json:"status,omitempty"`
	OriginTraffic *bool                 `json:"origin_traffic,omitempty"` // api defaults to true so we need an explicit empty value
	Headers       []MatchResponseHeader `json:"headers,omitempty"`
}

// MatchResponseHeader ...
type MatchResponseHeader struct {
	Name  string `json:"name,omitempty"`
	Op    string `json:"op,omitempty"`
	Value string `json:"value,omitempty"`
}

//RateLimitResult ...
type RateLimitResult struct {
	RateLimit RateLimitRecord `json:"result"`
	Success   bool            `json:"success"`
	Errors    []Error         `json:"errors"`
	Messages  []string        `json:"messages"`
}

//RateLimitResults ...
type RateLimitResults struct {
	RateLimitList []RateLimitRecord `json:"result"`
	ResultsInfo   ResultsCount      `json:"result_info"`
	Success       bool              `json:"success"`
	Errors        []Error           `json:"errors"`
}

//RateLimit ...
type RateLimit interface {
	ListRateLimit(cisID string, zoneID string) ([]RateLimitRecord, error)
	GetRateLimit(cisID string, zoneID string, rateLimitID string) (*RateLimitRecord, error)
	CreateRateLimit(cisID string, zoneID string, rateLimitBody RateLimitRecord) (*RateLimitRecord, error)
	DeleteRateLimit(cisID string, zoneID string, rateLimitID string) error
	UpdateRateLimit(cisID string, zoneID string, rateLimitID string, rateLimitBody RateLimitRecord) (*RateLimitRecord, error)
}

//RateLimit ...
type ratelimit struct {
	client *client.Client
}

func newRateLimitAPI(c *client.Client) RateLimit {
	return &ratelimit{
		client: c,
	}
}

func (r *ratelimit) ListRateLimit(cisID string, zoneID string) ([]RateLimitRecord, error) {
	rateLimitResults := RateLimitResults{}

	var rawURL string
	rawURL = fmt.Sprintf("/v1/%s/zones/%s/rate_limits", cisID, zoneID)
	_, err := r.client.Get(rawURL, &rateLimitResults, nil)
	if err != nil {
		return nil, err
	}
	return rateLimitResults.RateLimitList, err
}

func (r *ratelimit) GetRateLimit(cisID string, zoneID string, rateLimitID string) (*RateLimitRecord, error) {
	rateLimitResult := RateLimitResult{}
	var rawURL string
	rawURL = fmt.Sprintf("/v1/%s/zones/%s/rate_limits/%s", cisID, zoneID, rateLimitID)
	_, err := r.client.Get(rawURL, &rateLimitResult, nil)
	if err != nil {
		return nil, err
	}
	return &rateLimitResult.RateLimit, nil
}

func (r *ratelimit) DeleteRateLimit(cisID string, zoneID string, rateLimitID string) error {
	var rawURL string
	rawURL = fmt.Sprintf("/v1/%s/zones/%s/rate_limits/%s", cisID, zoneID, rateLimitID)
	_, err := r.client.Delete(rawURL)
	if err != nil {
		return err
	}
	return nil
}

func (r *ratelimit) CreateRateLimit(cisID string, zoneID string, rateLimitBody RateLimitRecord) (*RateLimitRecord, error) {
	rateLimitResult := RateLimitResult{}
	var rawURL string
	rawURL = fmt.Sprintf("/v1/%s/zones/%s/rate_limits", cisID, zoneID)
	_, err := r.client.Post(rawURL, &rateLimitBody, &rateLimitResult)
	if err != nil {
		return nil, err
	}
	return &rateLimitResult.RateLimit, nil
}

func (r *ratelimit) UpdateRateLimit(cisID string, zoneID string, rateLimitID string, rateLimitBody RateLimitRecord) (*RateLimitRecord, error) {
	rateLimitResult := RateLimitResult{}
	var rawURL string
	rawURL = fmt.Sprintf("/v1/%s/zones/%s/rate_limits/%s", cisID, zoneID, rateLimitID)
	_, err := r.client.Put(rawURL, &rateLimitBody, &rateLimitResult)
	if err != nil {
		return nil, err
	}
	return &rateLimitResult.RateLimit, nil
}
