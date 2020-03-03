package main

import (
	libvirtplugin "github.com/openshift/installer/plugins/libvirt/plugin"
	"github.com/openshift/installer/plugins"
)

func NewPlugin() plugins.Plugin {
	return &libvirtplugin.LibvirtPlugin{}
}

