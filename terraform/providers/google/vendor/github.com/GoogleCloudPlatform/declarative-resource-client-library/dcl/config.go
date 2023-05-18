// Copyright 2023 Google LLC. All Rights Reserved.
// 
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
// 
//     http://www.apache.org/licenses/LICENSE-2.0
// 
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
package dcl

import (
	"context"
	"fmt"
	"math/rand"
	"net/http"
	"net/http/httputil"
	"regexp"
	"strings"
	"time"

	// glog aliased import is necessary since these packages will be open-sourced
	// and that is the public name of the google logging package.
	glog "github.com/golang/glog"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/option"
)

const ua = "DeclarativeClientLib/0.0.1"

const defaultTimeout = 15 * time.Minute

// ConfigOption is used to functionally configure Configs.
type ConfigOption func(*Config)

// Config is used to enclose the credentials and http client used to make
// requests to GCP APIs.
type Config struct {
	RetryProvider       RetryProvider
	codeRetryability    map[int]Retryability
	timeout             time.Duration
	header              http.Header
	clientOptions       []option.ClientOption
	userAgent           string
	contentType         string
	queryParams         map[string]string
	Logger              ContextLogger
	BasePath            string
	billingProject      string
	userOverrideProject bool
}

// Retryability holds the details for one error code to determine if it is retyable.
// The regex field is compiled for use in error handling.
// To be retryable, the boolean must be true and the regex must match.
type Retryability struct {
	Retryable bool
	Pattern   string
	regex     *regexp.Regexp
	Timeout   time.Duration
}

// UserAgent returns the user agent for the config, which will always include the
// declarative SDK name and version.
func (c *Config) UserAgent() string {
	if c.userAgent != "" {
		return fmt.Sprintf("%s %s", c.userAgent, ua)
	}
	return ua
}

// NewConfig creates a Config object.
func NewConfig(o ...ConfigOption) *Config {
	retryable := Retryability{
		Retryable: true,
		regex:     regexp.MustCompile(".*"),
		Timeout:   defaultTimeout,
	}
	nonretryable := Retryability{Retryable: false}
	c := &Config{
		codeRetryability: map[int]Retryability{
			400: Retryability{
				Retryable: true,
				regex:     regexp.MustCompile("The resource '[-/a-zA-Z0-9]*' is not ready"),
				Timeout:   defaultTimeout,
			},
			403: Retryability{
				Retryable: true,
				regex:     regexp.MustCompile(".*API request rate quota.*"),
				Timeout:   defaultTimeout,
			},
			404: nonretryable,
			409: nonretryable,
			429: retryable,
			500: retryable,
			502: retryable,
			503: retryable,
		},
		contentType: "application/json",
		queryParams: map[string]string{"alt": "json"},
		Logger: ContextLogger{
			logger: DefaultLogger(LoggerInfo),
		},
		RetryProvider: &BackoffRetryProvider{},
	}

	for _, opt := range o {
		opt(c)
	}

	return c
}

// Clone returns a copy of an existing Config with optional new values.
func (c *Config) Clone(o ...ConfigOption) *Config {
	result := &Config{
		RetryProvider:       c.RetryProvider,
		codeRetryability:    c.codeRetryability,
		timeout:             c.timeout,
		clientOptions:       c.clientOptions,
		userAgent:           c.userAgent,
		contentType:         c.contentType,
		queryParams:         c.queryParams,
		Logger:              c.Logger,
		BasePath:            c.BasePath,
		billingProject:      c.billingProject,
		userOverrideProject: c.userOverrideProject,
	}

	if c.header != nil {
		result.header = c.header.Clone()
	}

	for _, opt := range o {
		opt(result)
	}

	return result
}

// TimeoutOr returns a timeout for this config. If WithTimeout() was called, that timeout
// is used; if WithTimeout() was not called and a value is provided with `t`, that is used.
// Otherwise the default timeout is returned;
func (c *Config) TimeoutOr(t time.Duration) time.Duration {
	if c.timeout != 0 {
		return c.timeout
	} else if t != 0 {
		return t
	}
	return defaultTimeout
}

type loggingTransport struct {
	underlyingTransport http.RoundTripper
	logger              ContextLogger
}

func (t loggingTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	shouldLogRequest, err := ShouldLogRequest(req.Context())
	if err != nil {
		t.logger.Infof("Error fetching ShouldLogRequest value: %v", err)
	}
	reqDump, err := httputil.DumpRequestOut(req, true)
	randString := RandomString(5)
	if err == nil {
		if shouldLogRequest {
			t.logger.InfoWithContextf(req.Context(), "Google API Request: (id %s)\n-----------[REQUEST]----------\n%s\n-------[END REQUEST]--------", randString, strings.ReplaceAll(string(reqDump), "\r\n", "\n"))
		}
	} else {
		t.logger.WarningWithContextf(req.Context(), "Failed to make request (id %s): %s", randString, err)
	}
	resp, err := t.underlyingTransport.RoundTrip(req)
	if err == nil {
		respDump, err := httputil.DumpResponse(resp, true)
		if err == nil {
			respDumpStr := string(respDump)
			if shouldLogRequest {
				t.logger.InfoWithContextf(req.Context(), "Google API Response: (id %s) \n-----------[RESPONSE]----------\n%s\n-------[END RESPONSE]--------", randString, strings.ReplaceAll(respDumpStr, "\r\n", "\n"))
			} else if resp.StatusCode >= 400 || strings.Contains(respDumpStr, "error") {
				t.logger.InfoWithContextf(req.Context(), "Google API Request: (id %s)\n-----------[REQUEST]----------\n%s\n-------[END REQUEST]--------", randString, strings.ReplaceAll(string(reqDump), "\r\n", "\n"))
				t.logger.InfoWithContextf(req.Context(), "Google API Response: (id %s) \n-----------[RESPONSE]----------\n%s\n-------[END RESPONSE]--------", randString, strings.ReplaceAll(respDumpStr, "\r\n", "\n"))
			}
		} else {
			t.logger.WarningWithContextf(req.Context(), "Failed to parse response (id %s): %s", randString, err)
		}
	} else {
		t.logger.WarningWithContextf(req.Context(), "Failed to get response (id %s): %s", randString, err)
	}
	return resp, err
}

// ApplyOption is an option that is accepted by Apply() functions.
type ApplyOption interface {
	Apply(*ApplyOpts)
}

// ApplyOpts refers to options that are taken in the apply function.
type ApplyOpts struct {
	params    []LifecycleParam
	stateHint Resource
}

type lifecycleParamOption struct {
	param LifecycleParam
}

func (l lifecycleParamOption) Apply(o *ApplyOpts) {
	o.params = append(o.params, l.param)
}

// WithLifecycleParam allows a user to specify the proper lifecycle params.
func WithLifecycleParam(d LifecycleParam) ApplyOption {
	return lifecycleParamOption{param: d}
}

// FetchLifecycleParams returns the list of lifecycle params.
func FetchLifecycleParams(c []ApplyOption) []LifecycleParam {
	var o ApplyOpts
	for _, p := range c {
		p.Apply(&o)
	}
	return o.params
}

type stateHint struct {
	state Resource
}

func (s stateHint) Apply(o *ApplyOpts) {
	o.stateHint = s.state
}

// WithStateHint takes in a resource which will be used in place of the applied
// resource any time the current configuration of the resource is relevant.
// For instance, if an identity field will change, passing a state hint will ensure
// that the current resource is fetched (and possibly deleted).
func WithStateHint(r Resource) ApplyOption {
	return stateHint{state: r}
}

// FetchStateHint returns either nil or a dcl.Resource representing the pre-apply state.
func FetchStateHint(c []ApplyOption) Resource {
	var o ApplyOpts
	for _, p := range c {
		p.Apply(&o)
	}
	return o.stateHint
}

// WithRetryProvider allows a user to override default exponential backoff retry behavior.
func WithRetryProvider(r RetryProvider) ConfigOption {
	return func(c *Config) {
		c.RetryProvider = r
	}
}

// WithCodeRetryability allows a user to add additional retryable or non-retryable error codes.
// Each error code is mapped to a regexp which must match the error message to be retryable.
func WithCodeRetryability(cr map[int]Retryability) ConfigOption {
	return func(c *Config) {
		for code, retryability := range cr {
			// Non-retryable errors do not need a regex to check against.
			var re *regexp.Regexp
			if retryability.Retryable {
				re = regexp.MustCompile(retryability.Pattern)
			}
			// If timeout for this retryable error was not specified, assume default.
			to := defaultTimeout
			if retryability.Timeout > 0 {
				to = retryability.Timeout
			}
			c.codeRetryability[code] = Retryability{
				Retryable: retryability.Retryable,
				regex:     re,
				Timeout:   to,
			}
		}
	}
}

// WithTimeout allows a user to override default operation timeout.
func WithTimeout(to time.Duration) ConfigOption {
	return func(c *Config) {
		c.timeout = to
	}
}

// WithLogger allows a user to specify a custom logger.
func WithLogger(l Logger) ConfigOption {
	return func(c *Config) {
		c.Logger.logger = l
	}
}

// WithContextLogger allows a user to specify a custom context logger.
func WithContextLogger(l ContextLogger) ConfigOption {
	return func(c *Config) {
		c.Logger = l
	}
}

// WithBasePath allows a base path to be overridden.
func WithBasePath(b string) ConfigOption {
	return func(c *Config) {
		c.BasePath = b
	}
}

// WithHeader allows aribitrary HTTP headers to be addded to requests. Not all headers
// (e.g., "Content-Type") can be overridden. To set the User-Agent header, use WithUserAgent().
func WithHeader(header, value string) ConfigOption {
	return func(c *Config) {
		if c.header == nil {
			c.header = make(http.Header)
		}
		c.header.Add(header, value)
	}
}

// WithUserAgent allows a user to specify a custom user-agent.
func WithUserAgent(ua string) ConfigOption {
	return func(c *Config) {
		c.userAgent = ua
	}
}

// WithContentType allows a user to override the default Content-Type header.
func WithContentType(ct string) ConfigOption {
	return func(c *Config) {
		c.contentType = ct
	}
}

// WithQueryParams allows a user to override the default query parameters.
func WithQueryParams(ps map[string]string) ConfigOption {
	return func(c *Config) {
		c.queryParams = ps
	}
}

// WithAPIKey returns a ConfigOption that specifies an API key to be used as the basis for authentication.
func WithAPIKey(apiKey string) ConfigOption {
	return func(c *Config) {
		c.clientOptions = append(c.clientOptions, option.WithAPIKey(apiKey))
	}
}

// WithClientCertSource returns a ConfigOption that specifies a callback function for obtaining a TLS client certificate.
func WithClientCertSource(s option.ClientCertSource) ConfigOption {
	return func(c *Config) {
		c.clientOptions = append(c.clientOptions, option.WithClientCertSource(s))
	}
}

// WithCredentials returns a ConfigOption that authenticates API calls using a caller-supplier Credentials struct.
func WithCredentials(creds *google.Credentials) ConfigOption {
	return func(c *Config) {
		c.clientOptions = append(c.clientOptions, option.WithCredentials(creds))
	}
}

// WithCredentialsFile returns a ConfigOption that authenticates API calls with the given service account or refresh token JSON credentials file.
func WithCredentialsFile(filename string) ConfigOption {
	return func(c *Config) {
		c.clientOptions = append(c.clientOptions, option.WithCredentialsFile(filename))
	}
}

// WithCredentialsJSON returns a ConfigOption that authenticates API calls with the given service account or refresh token JSON credentials.
func WithCredentialsJSON(p []byte) ConfigOption {
	return func(c *Config) {
		c.clientOptions = append(c.clientOptions, option.WithCredentialsJSON(p))
	}
}

// WithHTTPClient returns a ConfigOption that specifies the HTTP client to use as the basis of communications.
// When used, the WithHTTPClient option takes precedent over all other supplied authentication options.
func WithHTTPClient(client *http.Client) ConfigOption {
	return func(c *Config) {
		c.clientOptions = append(c.clientOptions, option.WithHTTPClient(client))
	}
}

// WithBillingProject returns a ConfigOption that specifies the user override project.
// This will be used to set X-Goog-User-Project on API calls.
// This option will be ignored unless WithUserProjectOverride is also used.
func WithBillingProject(project string) ConfigOption {
	return func(c *Config) {
		c.billingProject = project
	}
}

// WithUserProjectOverride returns a ConfigOption that turns on WithUserProjectOverride.
// This will send the X-Goog-User-Project on API calls.
func WithUserProjectOverride() ConfigOption {
	return func(c *Config) {
		c.userOverrideProject = true
	}
}

// Logger is an interface for logging requests and responses.
type Logger interface {
	Fatal(args ...interface{})
	Fatalf(format string, args ...interface{})
	Info(args ...interface{})
	Infof(format string, args ...interface{})
	Warningf(format string, args ...interface{})
	Warning(args ...interface{})
}

// ContextLogger is the internal logger implementation.
type ContextLogger struct {
	logger Logger
}

// LoggerLevel is the most basic level that a logger should print.
// Anything at this level or more severe will be printed by this logger.
type LoggerLevel int32

const (
	// Fatal will print only Fatal logs.
	Fatal LoggerLevel = iota
	// Error will print Error and all Fatal logs.
	Error
	// Warning will print Warning and all Error logs.
	Warning
	// LoggerInfo will print Info and all Warning logs.
	LoggerInfo
)

// DefaultLogger returns the default logger for the Declarative Client Library.
func DefaultLogger(level LoggerLevel) Logger {
	return glogger{level: level}
}

type glogger struct {
	level LoggerLevel
}

// Fatal records Fatal errors.
func (l glogger) Fatal(args ...interface{}) {
	if l.level >= Fatal {
		glog.Fatal(args...)
	}
}

// Fatalf records Fatal errors with added arguments.
func (l glogger) Fatalf(format string, args ...interface{}) {
	if l.level >= Fatal {
		glog.Fatalf(format, HandleLogArgs(args...)...)
	}
}

// Info records Info errors.
func (l glogger) Info(args ...interface{}) {
	if l.level >= LoggerInfo {
		glog.Info(args...)
	}
}

// Infof records Info errors with added arguments.
func (l glogger) Infof(format string, args ...interface{}) {
	if l.level >= LoggerInfo {
		glog.Infof(format, HandleLogArgs(args...)...)
	}
}

// Warningf records Warning errors with added arguments.
func (l glogger) Warningf(format string, args ...interface{}) {
	if l.level >= Warning {
		glog.Warningf(format, HandleLogArgs(args...)...)
	}
}

// Warning records Warning errors.
func (l glogger) Warning(args ...interface{}) {
	if l.level >= Warning {
		glog.Warning(args...)
	}
}

// Fatal records Fatal errors.
func (l ContextLogger) Fatal(args ...interface{}) {
	l.logger.Fatal(args...)
}

// Fatalf records Fatal errors with added arguments.
func (l ContextLogger) Fatalf(format string, args ...interface{}) {
	l.logger.Fatalf(format, HandleLogArgs(args...)...)
}

// Info records Info errors.
func (l ContextLogger) Info(args ...interface{}) {
	l.logger.Info(args...)
}

// Infof records Info errors with added arguments.
func (l ContextLogger) Infof(format string, args ...interface{}) {
	l.logger.Infof(format, HandleLogArgs(args...)...)
}

// Warningf records Warning errors with added arguments.
func (l ContextLogger) Warningf(format string, args ...interface{}) {
	l.logger.Warningf(format, HandleLogArgs(args...)...)
}

// Warning records Warning errors.
func (l ContextLogger) Warning(args ...interface{}) {
	l.logger.Warning(args...)
}

// FatalWithContext records Fatal errors with context values.
func (l ContextLogger) FatalWithContext(ctx context.Context, args ...interface{}) {
	args = append([]interface{}{ConstructLogPrefixFromContext(ctx)}, args...)
	l.Fatal(args...)
}

// FatalWithContextf records Fatal errors with added arguments with context values.
func (l ContextLogger) FatalWithContextf(ctx context.Context, format string, args ...interface{}) {
	format = fmt.Sprintf("%s %s", ConstructLogPrefixFromContext(ctx), format)
	l.Fatalf(format, args...)
}

// InfoWithContext records Info errors with context values.
func (l ContextLogger) InfoWithContext(ctx context.Context, args ...interface{}) {
	args = append([]interface{}{ConstructLogPrefixFromContext(ctx)}, args...)
	l.Info(args...)
}

// InfoWithContextf records Info errors with added arguments with context values.
func (l ContextLogger) InfoWithContextf(ctx context.Context, format string, args ...interface{}) {
	format = fmt.Sprintf("%s %s", ConstructLogPrefixFromContext(ctx), format)
	l.Infof(format, args...)
}

// WarningWithContextf records Warning errors with added arguments with context values.
func (l ContextLogger) WarningWithContextf(ctx context.Context, format string, args ...interface{}) {
	format = fmt.Sprintf("%s %s", ConstructLogPrefixFromContext(ctx), format)
	l.Warningf(format, HandleLogArgs(args...)...)
}

// WarningWithContext records Warning errors with context values.
func (l ContextLogger) WarningWithContext(ctx context.Context, args ...interface{}) {
	args = append([]interface{}{ConstructLogPrefixFromContext(ctx)}, args...)
	l.Warning(args...)
}

// HandleLogArgs ensures that pointer arguments are dereferenced well.
func HandleLogArgs(args ...interface{}) []interface{} {
	a := make([]interface{}, len(args))
	for i, v := range args {
		if s, ok := v.(*string); ok && s != nil {
			a[i] = *s
		} else {
			a[i] = v
		}
	}
	return a
}

// ConstructLogPrefixFromContext constructs log prefix from info in context
func ConstructLogPrefixFromContext(ctx context.Context) string {
	return fmt.Sprintf("[RequestID:%s] ", APIRequestID(ctx))
}

// RandomString generates a random alpha-numeric string of input length.
func RandomString(length int) string {
	charset := "abcdefghijklmnoqrstuvwxyz0123456789"
	var seededRand *rand.Rand = rand.New(
		rand.NewSource(time.Now().UnixNano()))

	b := make([]byte, length)
	for i := range b {
		b[i] = charset[seededRand.Intn(len(charset))]
	}
	return string(b)
}

// CreateAPIRequestID creates a random APIRequestId.
func CreateAPIRequestID() string {
	return RandomString(8)
}
