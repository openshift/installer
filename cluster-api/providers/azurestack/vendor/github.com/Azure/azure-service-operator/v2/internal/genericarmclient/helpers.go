/*
Copyright (c) Microsoft Corporation.
Licensed under the MIT license.
*/

package genericarmclient

import (
	"net/http"
	"strconv"
	"time"

	"github.com/pkg/errors"
)

var zeroDuration time.Duration = 0

func GetRetryAfter(resp *http.Response) time.Duration {
	if resp == nil {
		return zeroDuration
	}

	if retryAfterStr := resp.Header.Get("Retry-After"); retryAfterStr != "" {
		if retryAfterVal, parseErr := strconv.ParseInt(retryAfterStr, 10, 64); parseErr == nil {
			return time.Duration(retryAfterVal) * time.Second
		}

		if retryAfterTime, parseErr := parseHttpDate(retryAfterStr); parseErr == nil {
			result := time.Until(retryAfterTime)
			if result > 0 {
				return result
			}
		}
	}

	return 0
}

func parseHttpDate(s string) (time.Time, error) {
	if t, err := time.Parse("Mon, 02 Jan 2006 15:04:05 MST", s); err == nil {
		return t, nil
	} else if t, err = time.Parse("Monday, 02-Jan-06 15:04:05 MST", s); err == nil {
		return t, nil
	} else if t, err = time.Parse("Mon Jan  2 15:04:05 2006", s); err == nil {
		return t, nil
	}

	return time.Time{}, errors.New("unable to parse date")
}
