package server

import (
	"net/http"
	"net/http/httputil"
	"net/url"

	"golang.org/x/net/context"

	"github.com/coreos/tectonic-installer/installer/server/ctxh"
)

var proxyWhitelist = []string{
	"https://stable.release.core-os.net/amd64-usr/current/coreos_production_ami_all.json",
}

// proxy allows the client to access non-local domains without running afoul of cross-origin issues.
func proxy() ctxh.ContextHandler {
	fn := func(ctx context.Context, w http.ResponseWriter, req *http.Request) *ctxh.AppError {
		target := req.URL.Query().Get("target")
		if target == "" {
			return ctxh.NewAppError(nil, "proxy target missing", http.StatusBadRequest)
		}

		matched := ""
		for _, u := range proxyWhitelist {
			if u == target {
				matched = u
				break
			}
		}

		if matched == "" {
			return ctxh.NewAppError(nil, "proxy target not in whitelist", http.StatusBadRequest)
		}

		matchedURL, err := url.Parse(matched)
		if err != nil {
			return ctxh.NewAppError(err, "bug in url whitelist", http.StatusInternalServerError)
		}

		director := func(req *http.Request) {
			req.URL = matchedURL
			req.Host = matchedURL.Host // Cloudflare will reject requests without an accurate Host header
		}
		p := &httputil.ReverseProxy{Director: director}
		p.ServeHTTP(w, req)

		return nil
	}

	return ctxh.ContextHandlerFuncWithError(fn)
}
