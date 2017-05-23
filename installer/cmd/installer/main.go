// Package main is the Tectonic Installer binary app for end-users.
package main

import (
	"flag"
	"fmt"
	"net"
	"net/http"
	"os"

	log "github.com/Sirupsen/logrus"
	"github.com/coreos/pkg/flagutil"
	"github.com/toqueteos/webbrowser"

	"github.com/coreos/tectonic-installer/installer/api"
	"github.com/coreos/tectonic-installer/installer/pkg/terraform"
)

// Version is provided by compile time -ldflags.
var Version = "was not built properly"

func main() {
	// TerraForm entrypoint.
	if os.Getenv("TF_PLUGIN_MAGIC_COOKIE") != "" {
		terraform.ServePlugin(os.Args[1])
		return
	}

	flags := struct {
		address             string
		logLevel            string
		platforms           platformsValue
		cookieSigningSecret string
		disableSecureCookie bool
		openBrowser         bool
		help, version       bool
		// development features
		assetDir string
		devMode  bool
	}{}

	flag.StringVar(&flags.address, "address", "127.0.0.1:4444", "HTTP listen address (e.g. 127.0.0.1:4444)")
	flag.BoolVar(&flags.openBrowser, "open-browser", true, "open a browser window to interface with the installer")
	// Signing
	// downloadable binary should use a fixed secret and disable HTTPS-only cookies
	flag.StringVar(&flags.cookieSigningSecret, "cookie-signing-secret", "local", "Cookie signing/integrity secret")
	flag.BoolVar(&flags.disableSecureCookie, "disable-secure-cookie", true, "Allow cookies to be sent via HTTP")
	// Log levels https://github.com/Sirupsen/logrus/blob/master/logrus.go#L36
	flag.StringVar(&flags.logLevel, "log-level", "warn", "log level (e.g. \"debug\")")
	// Development
	flag.StringVar(&flags.assetDir, "asset-dir", "", "serve web assets from this directory rather than from internal storage")
	flag.BoolVar(&flags.devMode, "dev", false, "tell the front end that we're running in development mode")
	// Subcommands
	flag.BoolVar(&flags.version, "version", false, "print version and exit")
	flag.BoolVar(&flags.help, "help", false, "print usage and exit")

	flags.platforms = platformsValue{names: knownPlatforms}
	flag.Var(&flags.platforms, "platforms", "comma separated list of platforms to support")

	// parse command-line and environment variable arguments
	flag.Parse()
	if err := flagutil.SetFlagsFromEnv(flag.CommandLine, "INSTALLER"); err != nil {
		log.Fatal(err.Error())
	}

	if flags.help {
		fmt.Printf("%s: Tectonic Installer\n", os.Args[0])
		flag.PrintDefaults()
		return
	}

	if flags.version {
		fmt.Println(Version)
		return
	}

	// logging setup
	lvl, err := log.ParseLevel(flags.logLevel)
	if err != nil {
		log.Fatalf("invalid log-level: %v - use \"debug\" instead", err)
	}
	log.SetLevel(lvl)

	// HTTP server
	server, err := api.New(&api.Config{
		AssetDir:            flags.assetDir,
		DevMode:             flags.devMode,
		Platforms:           flags.platforms.names,
		CookieSigningSecret: flags.cookieSigningSecret,
		DisableSecureCookie: flags.disableSecureCookie,
	})
	if err != nil {
		log.Fatalf("invalid server config: %v", err)
	}

	fmt.Printf("Starting Tectonic Installer on %s\n", flags.address)
	ln, err := net.Listen("tcp", flags.address)
	if err != nil {
		log.Fatalf("failed to start listening: %v", err)
	}

	if flags.openBrowser {
		go func() {
			browseURL := "http://" + flags.address
			if err := webbrowser.Open(browseURL); err != nil {
				log.Fatalf("can't launch a web browser to view %s", browseURL)
			}
		}()
	}

	if err = http.Serve(ln, server); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
