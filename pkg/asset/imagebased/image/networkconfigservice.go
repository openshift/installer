package image

const (
	networkConfigService = `[Unit]
Description=Pre Network Manager Config
After=NetworkManager.service
Before=` + installationServiceName + `
[Service]
Type=oneshot
RemainAfterExit=yes
TimeoutSec=60
ExecStart=bash -c "until /usr/bin/nmstatectl apply /var/tmp/network-config.yaml; do sleep 1; done"
[Install]
WantedBy=multi-user.target
`
)
