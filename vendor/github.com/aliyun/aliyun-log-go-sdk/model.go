package sls

import (
	"encoding/json"
	"net/url"
	"strconv"
	"strings"
)

// GetLogRequest for GetLogsV2
type GetLogRequest struct {
	From     int64  // unix time, eg time.Now().Unix() - 900
	To       int64  // unix time, eg time.Now().Unix()
	Topic    string // @note topic is not used anymore, use __topic__ : xxx in query instead
	Lines    int64  // max 100; offset, lines and reverse is ignored when use SQL in query
	Offset   int64
	Reverse  bool
	Query    string
	PowerSQL bool
}

func (glr *GetLogRequest) ToURLParams() url.Values {
	urlVal := url.Values{}
	urlVal.Add("type", "log")
	urlVal.Add("from", strconv.Itoa(int(glr.From)))
	urlVal.Add("to", strconv.Itoa(int(glr.To)))
	urlVal.Add("topic", glr.Topic)
	urlVal.Add("line", strconv.Itoa(int(glr.Lines)))
	urlVal.Add("offset", strconv.Itoa(int(glr.Offset)))
	urlVal.Add("reverse", strconv.FormatBool(glr.Reverse))
	urlVal.Add("powerSql", strconv.FormatBool(glr.PowerSQL))
	urlVal.Add("query", glr.Query)
	return urlVal
}

// GetHistogramsResponse defines response from GetHistograms call
type SingleHistogram struct {
	Progress string `json:"progress"`
	Count    int64  `json:"count"`
	From     int64  `json:"from"`
	To       int64  `json:"to"`
}

type GetHistogramsResponse struct {
	Progress   string            `json:"progress"`
	Count      int64             `json:"count"`
	Histograms []SingleHistogram `json:"histograms"`
}

func (resp *GetHistogramsResponse) IsComplete() bool {
	return strings.ToLower(resp.Progress) == "complete"
}

// GetLogsResponse defines response from GetLogs call
type GetLogsResponse struct {
	Progress string              `json:"progress"`
	Count    int64               `json:"count"`
	Logs     []map[string]string `json:"logs"`
	Contents string              `json:"contents"`
	HasSQL   bool                `json:"hasSQL"`
}

// GetLogLinesResponse defines response from GetLogLines call
// note: GetLogLinesResponse.Logs is nil when use GetLogLinesResponse
type GetLogLinesResponse struct {
	GetLogsResponse
	Lines []json.RawMessage
}

func (resp *GetLogsResponse) IsComplete() bool {
	return strings.ToLower(resp.Progress) == "complete"
}

func (resp *GetLogsResponse) GetKeys() (error, []string) {
	type Content map[string][]interface{}
	var content Content
	err := json.Unmarshal([]byte(resp.Contents), &content)
	if err != nil {
		return err, nil
	}
	result := []string{}
	for _, v := range content["keys"] {
		result = append(result, v.(string))
	}
	return nil, result
}

type GetContextLogsResponse struct {
	Progress     string              `json:"progress"`
	TotalLines   int64               `json:"total_lines"`
	BackLines    int64               `json:"back_lines"`
	ForwardLines int64               `json:"forward_lines"`
	Logs         []map[string]string `json:"logs"`
}

func (resp *GetContextLogsResponse) IsComplete() bool {
	return strings.ToLower(resp.Progress) == "complete"
}

type JsonKey struct {
	Type     string `json:"type"`
	Alias    string `json:"alias,omitempty"`
	DocValue bool   `json:"doc_value,omitempty"`
}

// IndexKey ...
type IndexKey struct {
	Token         []string            `json:"token"` // tokens that split the log line.
	CaseSensitive bool                `json:"caseSensitive"`
	Type          string              `json:"type"` // text, long, double
	DocValue      bool                `json:"doc_value,omitempty"`
	Alias         string              `json:"alias,omitempty"`
	Chn           bool                `json:"chn"` // parse chinese or not
	JsonKeys      map[string]*JsonKey `json:"json_keys,omitempty"`
}

type IndexLine struct {
	Token         []string `json:"token"`
	CaseSensitive bool     `json:"caseSensitive"`
	IncludeKeys   []string `json:"include_keys,omitempty"`
	ExcludeKeys   []string `json:"exclude_keys,omitempty"`
	Chn           bool     `json:"chn"` // parse chinese or not
}

// Index is an index config for a log store.
type Index struct {
	Keys map[string]IndexKey `json:"keys,omitempty"`
	Line *IndexLine          `json:"line,omitempty"`
}

// CreateDefaultIndex return a full text index config
func CreateDefaultIndex() *Index {
	return &Index{
		Line: &IndexLine{
			Token:         []string{" ", "\n", "\t", "\r", ",", ";", "[", "]", "{", "}", "(", ")", "&", "^", "*", "#", "@", "~", "=", "<", ">", "/", "\\", "?", ":", "'", "\""},
			CaseSensitive: false,
		},
	}
}
