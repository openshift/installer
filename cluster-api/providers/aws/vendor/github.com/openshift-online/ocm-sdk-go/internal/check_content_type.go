/*
Copyright (c) 2021 Red Hat, Inc.

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

// This file contains the functions used to check content types.

package internal

import (
	"fmt"
	"html"
	"io"
	"mime"
	"net/http"
	"regexp"
	"strings"

	"github.com/microcosm-cc/bluemonday"
)

var wsRegex = regexp.MustCompile(`\s+`)

// CheckContentType checks that the content type of the given response is JSON. Note that if the
// content type isn't JSON this method will consume the complete body in order to generate an error
// message containing a summary of the content.
func CheckContentType(response *http.Response) error {
	var err error
	var mediaType string
	contentType := response.Header.Get("Content-Type")
	if contentType != "" {
		mediaType, _, err = mime.ParseMediaType(contentType)
		if err != nil {
			return err
		}
	} else {
		mediaType = contentType
	}
	if !strings.EqualFold(mediaType, "application/json") {
		var summary string
		summary, err = contentSummary(mediaType, response)
		if err != nil {
			return fmt.Errorf(
				"expected response content type 'application/json' but received "+
					"'%s' and couldn't obtain content summary: %w",
				mediaType, err,
			)
		}
		return fmt.Errorf(
			"expected response content type 'application/json' but received '%s' and "+
				"content '%s'",
			mediaType, summary,
		)
	}
	return nil
}

// contentSummary reads the body of the given response and returns a summary it. The summary will
// be the complete body if it isn't too log. If it is too long then the summary will be the
// beginning of the content followed by ellipsis.
func contentSummary(mediaType string, response *http.Response) (summary string, err error) {
	var body []byte
	body, err = io.ReadAll(response.Body)
	if err != nil {
		return
	}
	limit := 250
	runes := []rune(string(body))
	if strings.EqualFold(mediaType, "text/html") && len(runes) > limit {
		content := html.UnescapeString(bluemonday.StrictPolicy().Sanitize(string(body)))
		content = wsRegex.ReplaceAllString(strings.TrimSpace(content), " ")
		runes = []rune(content)
	}
	if len(runes) > limit {
		summary = fmt.Sprintf("%s...", string(runes[:limit]))
	} else {
		summary = string(runes)
	}
	return
}
