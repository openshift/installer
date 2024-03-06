/*
Copyright 2023 The Kubernetes Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package utils

import (
	"fmt"
	"net/url"
)

// getStartToken parses the given url string and gets the 'start' query param.
func getStartToken(nextURLS string) (string, error) {
	nextURL, err := url.Parse(nextURLS)
	if err != nil || nextURL == nil {
		return "", fmt.Errorf("could not parse next url for getting next resources: %w", err)
	}

	start := nextURL.Query().Get("start")
	return start, nil
}

// PagingHelper while listing resources, can use this to get the start token for getting the next set of resources for processing
// start token will get fetched from nextURL returned by f and passed to the func f.
// f should take start as param and return three values isDone bool, nextURL string, e error.
// isDone  - represents no need to iterate for getting next set of resources.
// nextURL - if nextURL is present, will try to get the start token and pass it to f for next set of resource processing.
// e       - if e is not nil, will break and return the error.
func PagingHelper(f func(string) (bool, string, error)) error {
	start := ""
	var err error
	for {
		isDone, nextURL, e := f(start)

		if e != nil {
			err = e
			break
		}

		if isDone {
			break
		}

		// for paging over next set of resources getting the start token
		if nextURL != "" {
			start, err = getStartToken(nextURL)
			if err != nil {
				break
			}
		} else {
			break
		}
	}

	return err
}
