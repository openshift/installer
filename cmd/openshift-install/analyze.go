package main

import (
	"archive/tar"
	"bufio"
	"compress/gzip"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path"
	"regexp"
	"sort"
	"strings"

	"github.com/sirupsen/logrus"
	"k8s.io/utils/pointer"

	"github.com/openshift/installer/pkg/systemd"
)

var failedUnitRegexp = regexp.MustCompile(`^[*] (?P<unit>[^ ]*) .*`)
var failedSSHRegexp = regexp.MustCompile(`(?P<all>ssh: connect to host (?P<host>[\d.]+) port (?P<port>\d+): (?P<error>.*))$`)

type unit struct {
	// Name is the name of the unit.
	Name string

	// State is the unit state
	State string

	// Detail is a detailed explaination of the current unit state.
	Detail string

	// Log has systemd logs from the unit.
	Log *systemd.Log
}

func (unit *unit) Render(host string) {
	logrus.Warnf("%s had %s systemd unit %s", host, unit.State, unit.Name)
	if unit.Detail != "" {
		logrus.Warnf(unit.Detail)
	}

	if unit.Log != nil && unit.Log.Restarts(unit.Name) > 0 {
		for _, line := range unit.Log.Format(unit.Name, -2) {
			logrus.Warn(line)
		}
	}
}

type host struct {
	// Name is an identifier for the host.
	Name string

	// Role is an identifier for the host's role (e.g. 'control-plane').
	Role string

	// AccessError, if non-empty, is the error from a failed attempt to
	// SSH into the host.
	AccessError string

	// Units is a set of systemd units from the host.
	Units map[string]*unit
}

func (host *host) Render() {
	name := host.Name
	if host.Role != "" && host.Role != host.Name {
		name = fmt.Sprintf("%s %s", host.Role, name)
	}

	if host.AccessError != "" {
		logrus.Warnf("%s access error: %s", name, host.AccessError)
		return
	}

	if len(host.Units) == 0 {
		logrus.Infof("%s had no failing or restarting systemd units", name)
		return
	}

	units := make([]string, 0, len(host.Units))
	for unitKey := range host.Units {
		units = append(units, unitKey)
	}
	sort.Strings(units)

	logrus.Warnf("%s had failing or restarting systemd units: %s", name, strings.Join(units, ", "))
	for _, unitKey := range units {
		host.Units[unitKey].Render(name)
	}
}

type bootstrapLogBundle struct {
	// AccessError, if non-empty, is the error from a failed attempt to
	// SSH into the host.
	AccessError string

	// Hosts is a set of hosts referened from the bundle.
	Hosts map[string]*host

	// Log is the output of the SSH gather script
	Log string
}

func newBootstrapLogBundle(r io.Reader) (*bootstrapLogBundle, error) {
	tarReader := tar.NewReader(r)
	bundle := &bootstrapLogBundle{
		Hosts: map[string]*host{},
	}
	for {
		header, err := tarReader.Next()
		if err == io.EOF {
			break
		}
		if err != nil {
			return bundle, err
		}

		baseName := path.Base(header.Name)

		hostKey := "bootstrap"
		role := "bootstrap"
		var accessError *string
		if strings.Contains(header.Name, "/control-plane/") {
			segments := strings.Split(header.Name, "/")
			for i, segment := range segments {
				if segment == "control-plane" && len(segments) > i && segments[i+1] != "" {
					hostKey = segments[i+1]
					role = "control-plane"
					if len(segments) > i+1 && segments[i+2] != "" {
						accessError = pointer.StringPtr("")
					} else {
						accessError = pointer.StringPtr("failed to SSH into host")
					}
					break
				}
			}
		}

		if _, ok := bundle.Hosts[hostKey]; !ok {
			bundle.Hosts[hostKey] = &host{Name: hostKey, Role: role}
		}
		host := bundle.Hosts[hostKey]
		if accessError != nil {
			host.AccessError = *accessError
		}

		if strings.HasSuffix(header.Name, "/gather.log") {
			data, err := ioutil.ReadAll(tarReader)
			if err != nil {
				return bundle, fmt.Errorf("reading %q: %w", header.Name, err)
			}
			bundle.Log = string(data)
		}

		if baseName == "failed-units.txt" {
			if host.Units == nil {
				host.Units = map[string]*unit{}
			}
			scanner := bufio.NewScanner(tarReader)
			for scanner.Scan() {
				line := scanner.Text()
				matches := failedUnitRegexp.FindStringSubmatch(line)
				if matches == nil {
					continue
				}

				for i, name := range failedUnitRegexp.SubexpNames() {
					if name == "unit" {
						if u, ok := bundle.Hosts[hostKey].Units[matches[i]]; ok {
							u.State = "failed"
						} else {
							bundle.Hosts[hostKey].Units[matches[i]] = &unit{
								Name:  matches[i],
								State: "failed",
							}
						}
						break
					}
				}
			}
			if err := scanner.Err(); err != nil {
				return bundle, fmt.Errorf("reading line from %q: %w", header.Name, err)
			}
		}

		if path.Base(path.Dir(header.Name)) == "unit-status" {
			var extension string
			for _, ext := range []string{".txt", ".log.json"} {
				if strings.HasSuffix(baseName, ext) {
					extension = ext
					break
				}
			}
			if extension == "" {
				continue
			}

			if host.Units == nil {
				host.Units = map[string]*unit{}
			}

			unitKey := strings.TrimSuffix(baseName, extension) // crio.service.txt -> crio.service

			if _, ok := host.Units[unitKey]; !ok {
				host.Units[unitKey] = &unit{Name: unitKey}
			}
			u := host.Units[unitKey]
			if u.State == "" {
				u.State = "restarting" // for now, assume we're restarting unless we are failing via failed-units.txt
			}

			switch extension {
			case ".txt":
				data, err := ioutil.ReadAll(tarReader)
				if err != nil {
					return bundle, fmt.Errorf("reading %q: %w", header.Name, err)
				}
				u.Detail = string(data)
			case ".log.json":
				log, err := systemd.NewLog(tarReader)
				if err != nil {
					return bundle, fmt.Errorf("parsing systemd log JSON %q: %w", header.Name, err)
				}
				u.Log = log
			}
		}
	}

	scanner := bufio.NewScanner(strings.NewReader(bundle.Log))
	for scanner.Scan() {
		line := scanner.Text()
		matches := failedSSHRegexp.FindStringSubmatch(line)
		if matches == nil {
			continue
		}

		groups := map[string]string{}
		for i, name := range failedSSHRegexp.SubexpNames() {
			groups[name] = matches[i]
		}

		if host, ok := bundle.Hosts[groups["host"]]; ok {
			host.AccessError = strings.TrimSpace(groups["all"])
		} else if bundle.AccessError == "" {
			bundle.AccessError = strings.TrimSpace(groups["all"])
		}
	}
	if err := scanner.Err(); err != nil {
		return bundle, fmt.Errorf("reading line from gather log: %w", err)
	}

	return bundle, nil
}

func (bundle *bootstrapLogBundle) Render() {
	if bundle.AccessError != "" {
		logrus.Warnf("bootstrap access error: %s", bundle.AccessError)
		return
	}

	hosts := make([]string, 0, len(bundle.Hosts))
	for hostKey := range bundle.Hosts {
		hosts = append(hosts, hostKey)
	}
	sort.Strings(hosts)

	for _, hostKey := range hosts {
		bundle.Hosts[hostKey].Render()
	}

	if len(hosts) <= 1 {
		logrus.Warnf("no control-plane machines in the gathered tarball")
	}
}

func analyzeGatheredBootstrap(tarPath string) error {
	file, err := os.Open(tarPath)
	if err != nil {
		return err
	}
	defer file.Close()

	gzipReader, err := gzip.NewReader(file)
	if err != nil {
		return err
	}
	defer gzipReader.Close()

	bundle, err := newBootstrapLogBundle(gzipReader)
	if err != nil {
		return fmt.Errorf("parsing bootstrap log bundle %s: %w", tarPath, err)
	}

	bundle.Render()

	return nil
}
