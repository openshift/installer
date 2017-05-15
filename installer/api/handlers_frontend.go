package api

import (
	"bytes"
	"encoding/json"
	"net/http"
	"path"
	"time"

	log "github.com/Sirupsen/logrus"

	"github.com/coreos/tectonic-installer/installer/assets"
)

func frontendHandler(assetDir string, platforms []string, devMode bool) http.Handler {
	mux := http.NewServeMux()
	assetHandler := http.HandlerFunc(servePublicAsset)
	if assetDir != "" {
		assetHandler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			serveAssetFromDir(assetDir, w, r)
		})
	}
	mux.Handle("/frontend/", http.StripPrefix("/frontend/", assetHandler))
	mux.Handle("/", serveIndex(platforms, devMode))
	return mux
}

func serveIndex(platforms []string, devMode bool) http.Handler {
	obj := struct {
		Platforms []string `json:"platforms"`
		DevMode   bool     `json:"devMode"`
	}{
		platforms, devMode,
	}
	jsonData, err := json.Marshal(obj)
	if err != nil {
		log.Errorf("Error marshalling config JSON: %v", err)
	}

	indexTmpl := assets.MustTemplateAsset("frontend/index.html.tmpl")
	b, err := assets.RenderTemplate(indexTmpl, struct {
		Config string
	}{
		string(jsonData),
	})
	if err != nil {
		log.Errorf("Error rendering index template: %v", err)
	}

	now := time.Now()
	fn := func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		w.Header().Set("Last-Modified", now.UTC().Format(http.TimeFormat))
		_, err := w.Write(b)
		if err != nil {
			log.Errorf("Error writing index: %v", err)
		}
	}
	return http.HandlerFunc(fn)
}

func servePublicAsset(w http.ResponseWriter, r *http.Request) {
	path := r.URL.Path
	assetName := "frontend/" + path
	asset, err := assets.Asset(assetName)
	log.Infof("Serving asset %s", assetName)
	if err != nil {
		http.Error(w, "no such asset", http.StatusNotFound)
		return
	}

	reader := bytes.NewReader(asset)
	http.ServeContent(w, r, path, time.Now(), reader)
}

func serveAssetFromDir(assetDir string, w http.ResponseWriter, r *http.Request) {
	servePath := path.Join(assetDir, r.URL.Path)
	log.Infof("Serving LOCAL FILE %s", servePath)
	http.ServeFile(w, r, servePath)
}
