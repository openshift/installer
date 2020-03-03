package plugins

import (
	"github.com/sirupsen/logrus"

	"github.com/openshift/installer/plugins/libvirt/loader"
)

func init() {
	KnownPlugins["terraform-provider-libvirt"] = func() {
		lvp, err := loader.LoadPlugin()
		if err != nil {
			logrus.Errorf(err.Error())
			return
		}
		lvp.Init()
	}
}
