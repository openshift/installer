package server

import (
	"bytes"
	"net/http"
	"time"

	"golang.org/x/net/context"

	"github.com/coreos/tectonic-installer/installer/server/asset"
	"github.com/coreos/tectonic-installer/installer/server/ctxh"
)

// zipAssetHandler returns a ContextHandler that writes assets to a zip file
// and serves a download response.
func zipAssetHandler(assets []asset.Asset) ctxh.ContextHandler {
	fn := func(ctx context.Context, w http.ResponseWriter, req *http.Request) *ctxh.AppError {
		buf, err := asset.ZipAssets(assets)
		if err != nil {
			return ctxh.NewAppError(err, "Cannot zip assets", http.StatusInternalServerError)
		}
		w.Header().Set("Content-Type", "application/zip")
		w.Header().Set("Content-Disposition", "attachment; filename=\"assets.zip\"")
		http.ServeContent(w, req, "assets.zip", time.Now(), bytes.NewReader(buf))
		return nil
	}
	return ctxh.ContextHandlerFuncWithError(fn)
}
