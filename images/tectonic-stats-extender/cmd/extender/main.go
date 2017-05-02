package main

import (
	"flag"
	"fmt"
	"time"

	log "github.com/Sirupsen/logrus"

	"github.com/coreos-inc/tectonic/stats-extender/pkg/extender"
	"github.com/coreos-inc/tectonic/stats-extender/pkg/version"
)

func main() {
	flags := struct {
		extensions extender.Extensions
		license    string
		logLevel   string
		output     string
		period     time.Duration
		publicKey  string
		version    bool
	}{extensions: extender.Extensions{}}

	flag.Var(&flags.extensions, "extension", "some extension to report in the form <key>:<value>")
	flag.StringVar(&flags.logLevel, "log-level", "info", "log level, e.g. \"debug\"")
	flag.StringVar(&flags.license, "license", "", "path to Tectonic license file")
	flag.StringVar(&flags.output, "output", "", "file to which to write extensions")
	flag.DurationVar(&flags.period, "period", time.Hour, "how often to send reports, eg 30m, 5h; set to 0 for one-shot mode")
	flag.StringVar(&flags.publicKey, "public-key", "", "path to the a public signing key for the Tectonic license; leave unset to use production key")
	flag.BoolVar(&flags.version, "version", false, "print version and exit")
	flag.Parse()

	lvl, err := log.ParseLevel(flags.logLevel)
	if err != nil {
		log.Fatalf("invalid log-level: %v", err)
	}

	if flags.version {
		fmt.Println(version.Version)
		return
	}

	if flags.output == "" {
		log.Fatal("--output is required")
	}

	if flags.license == "" {
		log.Fatal("--license is required")
	}

	fmt.Println(flags.extensions.String())
	l := log.New()
	l.Level = lvl
	e := extender.New(flags.extensions, flags.license, flags.output, flags.period, flags.publicKey)
	e.Run(l)
}
