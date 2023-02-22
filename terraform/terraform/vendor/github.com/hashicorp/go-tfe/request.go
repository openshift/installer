package tfe

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	retryablehttp "github.com/hashicorp/go-retryablehttp"
	"golang.org/x/time/rate"
)

// ClientRequest encapsulates a request sent by the Client
type ClientRequest struct {
	retryableRequest *retryablehttp.Request
	http             *retryablehttp.Client
	limiter          *rate.Limiter

	// Header are the headers that will be sent in this request
	Header http.Header
}

func (r ClientRequest) Do(ctx context.Context, model interface{}) error {
	// Wait will block until the limiter can obtain a new token
	// or returns an error if the given context is canceled.
	if err := r.limiter.Wait(ctx); err != nil {
		return err
	}

	// Add the context to the request.
	reqWithCxt := r.retryableRequest.WithContext(ctx)

	// Execute the request and check the response.
	resp, err := r.http.Do(reqWithCxt)
	if err != nil {
		// If we got an error, and the context has been canceled,
		// the context's error is probably more useful.
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
			return err
		}
	}
	defer resp.Body.Close()

	// Basic response checking.
	if err := checkResponseCode(resp); err != nil {
		return err
	}

	// Return here if decoding the response isn't needed.
	if model == nil {
		return nil
	}

	// If v implements io.Writer, write the raw response body.
	if w, ok := model.(io.Writer); ok {
		_, err := io.Copy(w, resp.Body)
		return err
	}

	return unmarshalResponse(resp.Body, model)
}

// doIpRanges is similar to Do except that The IP ranges API is not returning jsonapi
// like every other endpoint which means we need to handle it differently.
func (r *ClientRequest) doIpRanges(ctx context.Context, ir *IPRange) error {
	// Wait will block until the limiter can obtain a new token
	// or returns an error if the given context is canceled.
	if err := r.limiter.Wait(ctx); err != nil {
		return err
	}

	// Add the context to the request.
	contextReq := r.retryableRequest.WithContext(ctx)

	// Execute the request and check the response.
	resp, err := r.http.Do(contextReq)
	if err != nil {
		// If we got an error, and the context has been canceled,
		// the context's error is probably more useful.
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
			return err
		}
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 && resp.StatusCode >= 400 {
		return fmt.Errorf("error HTTP response while retrieving IP ranges: %d", resp.StatusCode)
	} else if resp.StatusCode == 304 {
		return nil
	}

	err = json.NewDecoder(resp.Body).Decode(ir)
	if err != nil {
		return err
	}
	return nil
}
