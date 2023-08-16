package edgegrid

import (
	"bytes"
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"io/ioutil"
	"net/http"
	"sort"
	"strings"
	"time"

	"github.com/google/uuid"
	"go.uber.org/ratelimit"
)

type (
	// Signer is the request signer interface
	Signer interface {
		SignRequest(r *http.Request)
		CheckRequestLimit(requestLimit int)
	}

	authHeader struct {
		authType    string
		clientToken string
		accessToken string
		timestamp   string
		nonce       string
		signature   string
	}
)

const (
	authType = "EG1-HMAC-SHA256"
)

var (
	// rateLimit represents the maximum number of API requests per second the provider can make
	requestLimit ratelimit.Limiter
)

// SignRequest adds a signed authorization header to the http request
func (c Config) SignRequest(r *http.Request) {
	if r.URL.Host == "" {
		r.URL.Host = c.Host
	}
	if r.URL.Scheme == "" {
		r.URL.Scheme = "https"
	}
	r.URL.RawQuery = c.addAccountSwitchKey(r)
	r.Header.Set("Authorization", c.createAuthHeader(r).String())
}

// CheckRequestLimit waits if necessary to ensure that OpenAPI's request limit is not exceeded
func (c Config) CheckRequestLimit(limit int) {
	if limit > 0 {
		if requestLimit == nil {
			requestLimit = ratelimit.New(limit)
		}
		requestLimit.Take()
	}
}

func (c Config) createAuthHeader(r *http.Request) authHeader {
	timestamp := Timestamp(time.Now())

	auth := authHeader{
		authType:    authType,
		clientToken: c.ClientToken,
		accessToken: c.AccessToken,
		timestamp:   timestamp,
		nonce:       uuid.New().String(),
	}

	msgPath := r.URL.EscapedPath()
	if r.URL.RawQuery != "" {
		msgPath = fmt.Sprintf("%s?%s", msgPath, r.URL.RawQuery)
	}

	// create the message to be signed
	msgData := []string{
		r.Method,
		r.URL.Scheme,
		r.URL.Host,
		msgPath,
		canonicalizeHeaders(r.Header, c.HeaderToSign),
		createContentHash(r, c.MaxBody),
		auth.String(),
	}
	msg := strings.Join(msgData, "\t")

	key := createSignature(timestamp, c.ClientSecret)
	auth.signature = createSignature(msg, key)
	return auth
}

func canonicalizeHeaders(requestHeaders http.Header, headersToSign []string) string {
	var unsortedHeader []string
	var sortedHeader []string
	for k := range requestHeaders {
		unsortedHeader = append(unsortedHeader, k)
	}
	sort.Strings(unsortedHeader)
	for _, k := range unsortedHeader {
		for _, sign := range headersToSign {
			if sign == k {
				v := strings.TrimSpace(requestHeaders.Get(k))
				sortedHeader = append(sortedHeader, fmt.Sprintf("%s:%s", strings.ToLower(k), strings.ToLower(stringMinifier(v))))
			}
		}
	}
	return strings.Join(sortedHeader, "\t")
}

// The content hash is the base64-encoded SHA–256 hash of the POST body.
// For any other request methods, this field is empty. But the tab separator (\t) must be included.
// The size of the POST body must be less than or equal to the value specified by the service.
// Any request that does not meet this criteria SHOULD be rejected during the signing process,
// as the request will be rejected by EdgeGrid.
func createContentHash(r *http.Request, maxBody int) string {
	var (
		contentHash  string
		preparedBody string
		bodyBytes    []byte
	)

	if r.Body != nil {
		bodyBytes, _ = ioutil.ReadAll(r.Body)
		r.Body = ioutil.NopCloser(bytes.NewBuffer(bodyBytes))
		preparedBody = string(bodyBytes)
	}

	if r.Method == http.MethodPost && len(preparedBody) > 0 {
		if len(preparedBody) > maxBody {
			preparedBody = preparedBody[0:maxBody]
		}

		sum := sha256.Sum256([]byte(preparedBody))

		contentHash = base64.StdEncoding.EncodeToString(sum[:])
	}

	return contentHash
}

func (a authHeader) String() string {
	auth := fmt.Sprintf("%s client_token=%s;access_token=%s;timestamp=%s;nonce=%s;",
		a.authType,
		a.clientToken,
		a.accessToken,
		a.timestamp,
		a.nonce)
	if a.signature != "" {
		auth += fmt.Sprintf("signature=%s", a.signature)
	}
	return auth
}

func (c Config) addAccountSwitchKey(r *http.Request) string {
	if c.AccountKey != "" {
		values := r.URL.Query()
		values.Add("accountSwitchKey", c.AccountKey)
		r.URL.RawQuery = values.Encode()
	}
	return r.URL.RawQuery
}
