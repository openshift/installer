/**
 * © Copyright IBM Corporation 2020, 2023. All Rights Reserved.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *      http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package auth

import (
	"net/http"
	"sync"
	"time"
)

// endOfTime is the time when session (non-persistent) cookies expire.
var endOfTime = time.Now().AddDate(0, 0, 14)

// session represent CouchDB AuthSession token and its expiration period.
type session struct {
	cookie       *http.Cookie
	expires      time.Time
	refreshTime  time.Time
	refreshMutex sync.Mutex
}

// newSession returns new session object constructed from AuthSession cookie.
func newSession(c *http.Cookie) (*session, error) {
	expires := c.Expires
	if expires.IsZero() {
		expires = endOfTime
	}

	// refreshTime is 20% of period between now and the expiration time
	return &session{
		cookie:      c,
		expires:     expires,
		refreshTime: expires.Add(-(time.Until(expires) * 20 / 100)),
	}, nil
}

func (s *session) getCookie() *http.Cookie {
	return s.cookie
}

// isValid checks if the auth cookie hasn't expired yet
func (s *session) isValid() bool {
	return time.Now().Before(s.expires)
}

// needsRefresh atomically identifies if the cookie is near of the expiration time
func (s *session) needsRefresh() bool {
	now := time.Now()
	if now.After(s.refreshTime) {
		s.refreshMutex.Lock()
		defer s.refreshMutex.Unlock()

		// advance refresh time by one minute to prevent a parallel process
		// that might be waiting on mutex right now
		// from starting the second refresh process
		if now.After(s.refreshTime) {
			s.refreshTime = time.Now().Add(time.Minute)
			return true
		}
		return false
	}
	return false
}
