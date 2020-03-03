package plugins

import (
	"github.com/sirupsen/logrus"

	"github.com/openshift/installer/plugins/libvirt/loader"
)

func init() {
	exec := func() {
		lvp, err := loader.LoadPlugin()
		if err != nil {
			logrus.Fatalf(err.Error())
		}
		lvp.Init()
	}
	KnownPlugins["terraform-provider-libvirt"] = exec
}
