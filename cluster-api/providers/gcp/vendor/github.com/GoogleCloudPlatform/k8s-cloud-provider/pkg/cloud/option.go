package cloud

import (
	"net/http"
)

// Option are optional parameters to the generated methods.
type Option func(*allOptions)

// allOptions that can be configured for the generated methods.
type allOptions struct {
	projectID            string
	addHeaders           http.Header
	returnPartialSuccess bool
}

func mergeOptions(options []Option) allOptions {
	var ret allOptions
	for _, opt := range options {
		opt(&ret)
	}
	return ret
}

// ReturnPartialSuccess only affects AggregatedList calls, and will have the call return success
// when fanout calls fail. For example, when partial success behavior is enabled, aggregatedList for
// a single zone scope either returns all resources in the zone or no resources, with an error code.
func ReturnPartialSuccess(returnPartialSuccess bool) Option {
	return func(opts *allOptions) {
		opts.returnPartialSuccess = returnPartialSuccess
	}
}

// ForceProjectID forces the projectID to be used in the call to be the one
// specified. This ignores the default routing done by the ProjectRouter.
func ForceProjectID(projectID string) Option {
	return func(opts *allOptions) {
		opts.projectID = projectID
	}
}

// WithHeaders sets the headers to be used in the call.
func AddHeaders(headers http.Header) Option {
	return func(opts *allOptions) {
		opts.addHeaders = headers
	}
}

func handleHeaderOptions(opts *allOptions, to http.Header) {
	for k, vals := range opts.addHeaders {
		for _, v := range vals {
			to.Add(k, v)
		}
	}
}
