package main

import (
	"bytes"
	"flag"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"time"

	"github.com/coreos-inc/tectonic/ingress-error/binassets"
)

var (
	addr = flag.String("addr", "0.0.0.0:8080", "address to serve default backend.")

	errorPage = binassets.MustAsset("error.html")
	indexPage = binassets.MustAsset("index.html")
)

func handleErrorPage(w http.ResponseWriter, r *http.Request) {
	errorCode := r.Header.Get("X-Code")

	tmpl, _ := template.New("").Parse(string(errorPage))

	var errMsg string
	switch errorCode {
	case "400":
		errMsg = "Bad Request"
	case "401":
		errMsg = "Unauthorized Access"
	case "403":
		errMsg = "Forbidden"
	case "404":
		errMsg = "Not Found"
	case "500":
		errMsg = "Internal Server Error"
	case "503":
		errMsg = "Service Unavailable"
	case "504":
		errMsg = "Gateway Time-out"
	default:
		reader := bytes.NewReader(indexPage)
		http.ServeContent(w, r, r.URL.Path, time.Now(), reader)
		return
	}

	data := struct {
		ErrCode string
		ErrMsg  string
	}{errorCode, errMsg}

	if err := tmpl.Execute(w, data); err != nil {
		log.Printf("Unable to execute template.")
		return
	}

}

func main() {
	flag.Parse()
	http.HandleFunc("/", handleErrorPage)
	http.HandleFunc("/healthz", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		fmt.Fprint(w, "ok")
	})
	http.ListenAndServe(fmt.Sprintf("%s", *addr), nil)
}
