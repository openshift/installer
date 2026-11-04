package nethttplibrary

import (
	"context"
	"errors"
	"fmt"
	nethttp "net/http"
	"net/url"
	"strings"

	abs "github.com/microsoft/kiota-abstractions-go"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
)

// RedirectHandler handles redirect responses and follows them according to the options specified.
type RedirectHandler struct {
	// options to use when evaluating whether to redirect or not
	options RedirectHandlerOptions
}

// NewRedirectHandler creates a new redirect handler with the default options.
func NewRedirectHandler() *RedirectHandler {
	return NewRedirectHandlerWithOptions(RedirectHandlerOptions{
		MaxRedirects: defaultMaxRedirects,
		ShouldRedirect: func(req *nethttp.Request, res *nethttp.Response) bool {
			return true
		},
		ScrubSensitiveHeaders: DefaultScrubSensitiveHeaders,
	})
}

// NewRedirectHandlerWithOptions creates a new redirect handler with the specified options.
func NewRedirectHandlerWithOptions(options RedirectHandlerOptions) *RedirectHandler {
	if options.ScrubSensitiveHeaders == nil {
		options.ScrubSensitiveHeaders = DefaultScrubSensitiveHeaders
	}
	return &RedirectHandler{options: options}
}

// ScrubSensitiveHeaders is a callback function type for scrubbing sensitive headers during redirects.
// It receives the request to modify (which contains the new URL) and the original URL for comparison.
type ScrubSensitiveHeaders func(request *nethttp.Request, originalURL *url.URL)

// DefaultScrubSensitiveHeaders is the default implementation for scrubbing sensitive headers during redirects.
// This function removes Authorization and Cookie headers when the host, scheme, or port changes.
// Note: Proxy-Authorization is not handled here as proxy configuration in Go's net/http
// is managed at the transport level and not accessible to middleware.
var DefaultScrubSensitiveHeaders ScrubSensitiveHeaders = func(request *nethttp.Request, originalURL *url.URL) {
	if request == nil || originalURL == nil {
		return
	}

	newURL := request.URL
	if newURL == nil {
		return
	}

	// Remove Authorization and Cookie headers if the request's scheme, host, or port changes
	isDifferentOrigin := !strings.EqualFold(originalURL.Host, newURL.Host) ||
		!strings.EqualFold(originalURL.Scheme, newURL.Scheme) ||
		originalURL.Port() != newURL.Port()

	if isDifferentOrigin {
		request.Header.Del("Authorization")
		request.Header.Del("Cookie")
	}

	// Note: Proxy-Authorization is not handled here as proxy configuration in Go's net/http
	// is managed at the transport level (http.Transport.Proxy) and not accessible to middleware.
	// In environments where this matters, the proxy configuration should be managed at the HTTP client level.
}

// RedirectHandlerOptions to use when evaluating whether to redirect or not.
type RedirectHandlerOptions struct {
	// A callback that determines whether to redirect or not.
	ShouldRedirect func(req *nethttp.Request, res *nethttp.Response) bool
	// The maximum number of redirects to follow.
	MaxRedirects int
	// A callback for scrubbing sensitive headers during redirects.
	// Defaults to DefaultScrubSensitiveHeaders if not provided.
	ScrubSensitiveHeaders ScrubSensitiveHeaders
}

var redirectKeyValue = abs.RequestOptionKey{
	Key: "RedirectHandler",
}

type redirectHandlerOptionsInt interface {
	abs.RequestOption
	GetShouldRedirect() func(req *nethttp.Request, res *nethttp.Response) bool
	GetMaxRedirect() int
	GetScrubSensitiveHeaders() ScrubSensitiveHeaders
}

// GetKey returns the key value to be used when the option is added to the request context
func (options *RedirectHandlerOptions) GetKey() abs.RequestOptionKey {
	return redirectKeyValue
}

// GetShouldRedirect returns the redirection evaluation function.
func (options *RedirectHandlerOptions) GetShouldRedirect() func(req *nethttp.Request, res *nethttp.Response) bool {
	return options.ShouldRedirect
}

// GetMaxRedirect returns the maximum number of redirects to follow.
func (options *RedirectHandlerOptions) GetMaxRedirect() int {
	if options == nil || options.MaxRedirects < 1 {
		return defaultMaxRedirects
	} else if options.MaxRedirects > absoluteMaxRedirects {
		return absoluteMaxRedirects
	} else {
		return options.MaxRedirects
	}
}

// GetScrubSensitiveHeaders returns the header scrubbing function.
func (options *RedirectHandlerOptions) GetScrubSensitiveHeaders() ScrubSensitiveHeaders {
	if options == nil || options.ScrubSensitiveHeaders == nil {
		return DefaultScrubSensitiveHeaders
	}
	return options.ScrubSensitiveHeaders
}

const defaultMaxRedirects = 5
const absoluteMaxRedirects = 20
const movedPermanently = 301
const found = 302
const seeOther = 303
const temporaryRedirect = 307
const permanentRedirect = 308
const locationHeader = "Location"

// Intercept implements the interface and evaluates whether to follow a redirect response.
func (middleware RedirectHandler) Intercept(pipeline Pipeline, middlewareIndex int, req *nethttp.Request) (*nethttp.Response, error) {
	obsOptions := GetObservabilityOptionsFromRequest(req)
	ctx := req.Context()
	var span trace.Span
	var observabilityName string
	if obsOptions != nil {
		observabilityName = obsOptions.GetTracerInstrumentationName()
		ctx, span = otel.GetTracerProvider().Tracer(observabilityName).Start(ctx, "RedirectHandler_Intercept")
		span.SetAttributes(attribute.Bool("com.microsoft.kiota.handler.redirect.enable", true))
		defer span.End()
		req = req.WithContext(ctx)
	}
	response, err := pipeline.Next(req, middlewareIndex)
	if err != nil {
		return response, err
	}
	reqOption, ok := req.Context().Value(redirectKeyValue).(redirectHandlerOptionsInt)
	if !ok {
		reqOption = &middleware.options
	}
	return middleware.redirectRequest(ctx, pipeline, middlewareIndex, reqOption, req, response, 0, observabilityName)
}

func (middleware RedirectHandler) redirectRequest(ctx context.Context, pipeline Pipeline, middlewareIndex int, reqOption redirectHandlerOptionsInt, req *nethttp.Request, response *nethttp.Response, redirectCount int, observabilityName string) (*nethttp.Response, error) {
	shouldRedirect := reqOption.GetShouldRedirect() != nil && reqOption.GetShouldRedirect()(req, response) || reqOption.GetShouldRedirect() == nil
	if middleware.isRedirectResponse(response) &&
		redirectCount < reqOption.GetMaxRedirect() &&
		shouldRedirect {
		redirectCount++
		redirectRequest, err := middleware.getRedirectRequest(req, response)
		if err != nil {
			return response, err
		}
		if observabilityName != "" {
			ctx, span := otel.GetTracerProvider().Tracer(observabilityName).Start(ctx, "RedirectHandler_Intercept - redirect "+fmt.Sprint(redirectCount))
			span.SetAttributes(attribute.Int("com.microsoft.kiota.handler.redirect.count", redirectCount),
				httpResponseStatusCodeAttribute.Int(response.StatusCode),
			)
			defer span.End()
			redirectRequest = redirectRequest.WithContext(ctx)
		}

		result, err := pipeline.Next(redirectRequest, middlewareIndex)
		if err != nil {
			return result, err
		}
		return middleware.redirectRequest(ctx, pipeline, middlewareIndex, reqOption, redirectRequest, result, redirectCount, observabilityName)
	}
	return response, nil
}

func (middleware RedirectHandler) isRedirectResponse(response *nethttp.Response) bool {
	if response == nil {
		return false
	}
	locationHeader := response.Header.Get(locationHeader)
	if locationHeader == "" {
		return false
	}
	statusCode := response.StatusCode
	return statusCode == movedPermanently || statusCode == found || statusCode == seeOther || statusCode == temporaryRedirect || statusCode == permanentRedirect
}

func (middleware RedirectHandler) getRedirectRequest(request *nethttp.Request, response *nethttp.Response) (*nethttp.Request, error) {
	if request == nil || response == nil {
		return nil, errors.New("request or response is nil")
	}
	locationHeaderValue := response.Header.Get(locationHeader)
	if locationHeaderValue[0] == '/' {
		locationHeaderValue = request.URL.Scheme + "://" + request.URL.Host + locationHeaderValue
	}
	result := request.Clone(request.Context())
	targetUrl, err := url.Parse(locationHeaderValue)
	if err != nil {
		return nil, err
	}
	result.URL = targetUrl
	if result.Host != targetUrl.Host {
		result.Host = targetUrl.Host
	}

	// Scrub sensitive headers before following the redirect
	scrubber := middleware.options.GetScrubSensitiveHeaders()
	if scrubber != nil {
		scrubber(result, request.URL)
	}

	if response.StatusCode == seeOther {
		result.Method = nethttp.MethodGet
		result.Header.Del("Content-Type")
		result.Header.Del("Content-Length")
		result.Body = nil
	}
	return result, nil
}
