package azure

import (
	"github.com/coreos/ignition/config/util"
	ignition "github.com/coreos/ignition/config/v2_2/types"
)

//ModifyPointerIgnitionConfig modifies the ignitionconfig pointer as per the platform requirements
func ModifyPointerIgnitionConfig(ignitionConfig *ignition.Config) {
	//until afterburn sets the hostname for azure
	//and the fix is included in the azure image: https://github.com/coreos/afterburn/issues/197
	setHostname(ignitionConfig)
}

func setHostname(ignitionConfig *ignition.Config) {
	ignitionConfig.Systemd.Units = append(ignitionConfig.Systemd.Units, ignition.Unit{
		Name:     "setazurehostname.service",
		Enabled:  util.BoolToPtr(true),
		Enable:   true,
		Contents: "[Service]\nType=oneshot\nExecStart=curl -H Metadata:true \"http://169.254.169.254/metadata/instance/compute/name?api-version=2017-08-01&format=text\" -o /tmp/hostname\nExecStart=mv -f /tmp/hostname /etc/hostname\n\n[Install]\nWantedBy=multi-user.target",
	})
}
