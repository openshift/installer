package installerassets

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"path"

	"github.com/coreos/ignition/config/v2_2/types"
	"github.com/ghodss/yaml"
	"github.com/openshift/installer/pkg/assets"
	"github.com/vincent-petithory/dataurl"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// fileFromBytes creates an ignition-config file with the given
// contents.
func fileFromBytes(path string, mode int, contents []byte) types.File {
	return types.File{
		Node: types.Node{
			Filesystem: "root",
			Path:       path,
		},
		FileEmbedded1: types.FileEmbedded1{
			Mode: &mode,
			Contents: types.FileContents{
				Source: dataurl.EncodeBytes(contents),
			},
		},
	}
}

func bootstrapIgnRebuilder(ctx context.Context, getByName assets.GetByString) (*assets.Asset, error) {
	asset := &assets.Asset{
		Name:          "ignition/bootstrap.ign",
		RebuildHelper: bootstrapIgnRebuilder,
	}

	config := &types.Config{
		Ignition: types.Ignition{
			Version: types.MaxVersion.String(),
		},
	}

	parents, err := asset.GetParents(ctx, getByName, "ssh.pub")
	if err == nil {
		config.Passwd.Users = append(
			config.Passwd.Users,
			types.PasswdUser{
				Name: "core",
				SSHAuthorizedKeys: []types.SSHAuthorizedKey{
					types.SSHAuthorizedKey(parents["ssh.pub"].Data),
				},
			},
		)
	} else if !os.IsNotExist(err) {
		return nil, err
	}

	for _, entry := range []struct {
		name    string
		enabled bool
	}{
		{name: "systemd/units/bootkube.service"},
		{name: "systemd/units/tectonic.service"},
		{name: "systemd/units/progress.service", enabled: true},
		{name: "systemd/units/kubelet.service", enabled: true},
	} {
		parents, err := asset.GetParents(ctx, getByName, entry.name)
		if err != nil {
			return nil, err
		}

		unit := types.Unit{
			Name:     path.Base(parents[entry.name].Name),
			Contents: string(parents[entry.name].Data),
		}
		if entry.enabled {
			unit.Enabled = &entry.enabled
		}

		config.Systemd.Units = append(config.Systemd.Units, unit)
	}

	parents, err = asset.GetParents(ctx, getByName, "platform")
	if err != nil {
		return nil, err
	}
	platform := string(parents["platform"].Data)

	for _, entry := range []struct {
		path     string
		name     string
		mode     int
		append   bool
		user     string
		group    string
		platform string
	}{
		{path: "/etc/kubernetes/kubeconfig", name: "auth/kubeconfig-kubelet", mode: 0600},
		{path: "/etc/motd", name: "files/etc/motd", mode: 0644, append: true},
		{path: "/etc/ssl/etcd/ca.crt", name: "tls/etcd-client.crt", mode: 0600},
		{path: "/home/core/.bash_history", name: "files/home/core/.bash_history", mode: 0600, user: "core", group: "core"},
		{path: "/opt/tectonic/bootkube-config-overrides/kube-apiserver-config-overrides.yaml", name: "files/opt/tectonic/bootkube-config-overrides/kube-apiserver-config-overrides.yaml", mode: 0600},
		{path: "/opt/tectonic/bootkube-config-overrides/kube-controller-manager-config-overrides.yaml", name: "files/opt/tectonic/bootkube-config-overrides/kube-controller-manager-config-overrides.yaml", mode: 0600},
		{path: "/opt/tectonic/bootkube-config-overrides/kube-scheduler-config-overrides.yaml", name: "files/opt/tectonic/bootkube-config-overrides/kube-scheduler-config-overrides.yaml", mode: 0600},
		{path: "/opt/tectonic/auth/kubeconfig", name: "auth/kubeconfig-admin", mode: 0600},
		{path: "/opt/tectonic/manifests/03-openshift-web-console-namespace.yaml", name: "files/opt/tectonic/manifests/03-openshift-web-console-namespace.yaml", mode: 0600},
		{path: "/opt/tectonic/manifests/04-openshift-machine-config-operator.yaml", name: "files/opt/tectonic/manifests/04-openshift-machine-config-operator.yaml", mode: 0600},
		{path: "/opt/tectonic/manifests/05-openshift-cluster-api-namespace.yaml", name: "files/opt/tectonic/manifests/05-openshift-cluster-api-namespace.yaml", mode: 0600},
		{path: "/opt/tectonic/manifests/09-openshift-service-cert-signer-namespace.yaml", name: "files/opt/tectonic/manifests/09-openshift-service-cert-signer-namespace.yaml", mode: 0600},
		{path: "/opt/tectonic/manifests/cluster-config.yaml", name: "manifests/cluster-config.yaml", mode: 0600},
		{path: "/opt/tectonic/manifests/cluster-ingress-01-crd.yaml", name: "files/opt/tectonic/manifests/cluster-ingress-01-crd.yaml", mode: 0600},
		{path: "/opt/tectonic/manifests/cluster-ingress-02-config.yaml", name: "manifests/cluster-ingress-02-config.yaml", mode: 0600},
		{path: "/opt/tectonic/manifests/cluster-network-01-crd.yaml", name: "files/opt/tectonic/manifests/cluster-network-01-crd.yaml", mode: 0600},
		{path: "/opt/tectonic/manifests/cluster-network-02-config.yaml", name: "manifests/cluster-network-02-config.yaml", mode: 0600},
		{path: "/opt/tectonic/manifests/cvo-overrides.yaml", name: "files/opt/tectonic/manifests/cvo-overrides.yaml", mode: 0600},
		{path: "/opt/tectonic/manifests/etcd-service.yaml", name: "files/opt/tectonic/manifests/etcd-service.yaml", mode: 0600},
		{path: "/opt/tectonic/manifests/host-etcd-service.yaml", name: "files/opt/tectonic/manifests/host-etcd-service.yaml", mode: 0600},
		{path: "/opt/tectonic/manifests/host-etcd-service-endpoints.yaml", name: "files/opt/tectonic/manifests/host-etcd-service-endpoints.yaml", mode: 0600},
		{path: "/opt/tectonic/manifests/legacy-cvo-overrides.yaml", name: "files/opt/tectonic/manifests/legacy-cvo-overrides.yaml", mode: 0600},
		{path: "/opt/tectonic/manifests/kube-cloud-config.yaml", name: "files/opt/tectonic/manifests/kube-cloud-config.yaml", mode: 0600},
		{path: "/opt/tectonic/manifests/kube-system-secret-etcd-client.yaml", name: "files/opt/tectonic/manifests/kube-system-secret-etcd-client.yaml", mode: 0600},
		{path: "/opt/tectonic/manifests/kube-system-configmap-etcd-serving-ca.yaml", name: "files/opt/tectonic/manifests/kube-system-configmap-etcd-serving-ca.yaml", mode: 0600},
		{path: "/opt/tectonic/manifests/kube-system-configmap-root-ca.yaml", name: "files/opt/tectonic/manifests/kube-system-configmap-root-ca.yaml", mode: 0600},
		{path: "/opt/tectonic/manifests/machine-config-server-tls-secret.yaml", name: "files/opt/tectonic/manifests/machine-config-server-tls-secret.yaml", mode: 0600},
		{path: "/opt/tectonic/manifests/openshift-service-cert-signer-ca-secret.yaml", name: "files/opt/tectonic/manifests/openshift-service-cert-signer-ca-secret.yaml", mode: 0600},
		// FIXME: dup? {path: "/opt/tectonic/manifests/openshift-apiserver-secret.yaml", name: "manifests/openshift-apiserver-secret.yaml", mode: 0600},
		{path: "/opt/tectonic/manifests/pull.json", name: "manifests/pull.json", mode: 0600},
		{path: "/opt/tectonic/tectonic/99_binding-discovery.yaml", name: "files/opt/tectonic/tectonic/99_binding-discovery.yaml", mode: 0600},
		{path: "/opt/tectonic/tectonic/99_cloud-creds-secret.yaml", name: "files/opt/tectonic/tectonic/aws/99_cloud-creds-secret.yaml", mode: 0600, platform: "aws"},
		{path: "/opt/tectonic/tectonic/99_cloud-creds-secret.yaml", name: "files/opt/tectonic/tectonic/openstack/99_cloud-creds-secret.yaml", mode: 0600, platform: "openstack"},
		{path: "/opt/tectonic/tectonic/99_openshift-cluster-api_cluster.yaml", name: "manifests/99_openshift-cluster-api_cluster.yaml", mode: 0600},
		{path: "/opt/tectonic/tectonic/99_openshift-cluster-api_master-machines.yaml", name: "manifests/aws/99_openshift-cluster-api_master-machines.yaml", mode: 0600, platform: "aws"},
		{path: "/opt/tectonic/tectonic/99_openshift-cluster-api_master-machines.yaml", name: "manifests/libvirt/99_openshift-cluster-api_master-machines.yaml", mode: 0600, platform: "libvirt"},
		{path: "/opt/tectonic/tectonic/99_openshift-cluster-api_master-machines.yaml", name: "manifests/openstack/99_openshift-cluster-api_master-machines.yaml", mode: 0600, platform: "openstack"}, // FIXME
		{path: "/opt/tectonic/tectonic/99_openshift-cluster-api_master-user-data-secret.yaml", name: "manifests/99_openshift-cluster-api_master-user-data-secret.yaml"},
		{path: "/opt/tectonic/tectonic/99_openshift-cluster-api_worker-machinesets.yaml", name: "manifests/aws/99_openshift-cluster-api_worker-machinesets.yaml", mode: 0600, platform: "aws"},
		{path: "/opt/tectonic/tectonic/99_openshift-cluster-api_worker-machinesets.yaml", name: "manifests/libvirt/99_openshift-cluster-api_worker-machinesets.yaml", mode: 0600, platform: "libvirt"},
		{path: "/opt/tectonic/tectonic/99_openshift-cluster-api_worker-machinesets.yaml", name: "manifests/openstack/99_openshift-cluster-api_worker-machinesets.yaml", mode: 0600, platform: "openstack"}, // FIXME
		{path: "/opt/tectonic/tectonic/99_openshift-cluster-api_worker-user-data-secret.yaml", name: "manifests/99_openshift-cluster-api_worker-user-data-secret.yaml"},
		{path: "/opt/tectonic/tectonic/99_role-cloud-creds-secret-reader.yaml", name: "files/opt/tectonic/tectonic/aws/99_role-cloud-creds-secret-reader.yaml", mode: 0600, platform: "aws"},
		{path: "/opt/tectonic/tectonic/99_role-cloud-creds-secret-reader.yaml", name: "files/opt/tectonic/tectonic/openstack/99_role-cloud-creds-secret-reader.yaml", mode: 0600, platform: "openstack"},
		{path: "/opt/tectonic/tls/admin.crt", name: "tls/admin-client.crt", mode: 0600},
		{path: "/opt/tectonic/tls/admin.key", name: "tls/admin-client.key", mode: 0600},
		{path: "/opt/tectonic/tls/aggregator-ca.crt", name: "tls/aggregator-ca.crt", mode: 0600},
		{path: "/opt/tectonic/tls/aggregator-ca.key", name: "tls/aggregator-ca.key", mode: 0600},
		{path: "/opt/tectonic/tls/apiserver.crt", name: "tls/api-server-chain.crt", mode: 0600},
		{path: "/opt/tectonic/tls/apiserver.key", name: "tls/api-server.key", mode: 0600},
		{path: "/opt/tectonic/tls/apiserver-proxy.crt", name: "tls/api-server-proxy.crt", mode: 0600},
		{path: "/opt/tectonic/tls/apiserver-proxy.key", name: "tls/api-server-proxy.key", mode: 0600},
		{path: "/opt/tectonic/tls/etcd-client-ca.crt", name: "tls/etcd-ca.crt", mode: 0600},
		{path: "/opt/tectonic/tls/etcd-client-ca.key", name: "tls/etcd-ca.key", mode: 0600},
		{path: "/opt/tectonic/tls/etcd-client.crt", name: "tls/etcd-client.crt", mode: 0600},
		{path: "/opt/tectonic/tls/etcd-client.key", name: "tls/etcd-client.key", mode: 0600},
		{path: "/opt/tectonic/tls/kube-ca.crt", name: "tls/kube-ca.crt", mode: 0600},
		{path: "/opt/tectonic/tls/kube-ca.key", name: "tls/kube-ca.key", mode: 0600},
		{path: "/opt/tectonic/tls/kubelet.crt", name: "tls/kubelet-client.crt", mode: 0600},
		{path: "/opt/tectonic/tls/kubelet.key", name: "tls/kubelet-client.key", mode: 0600},
		{path: "/opt/tectonic/tls/machine-config-server.crt", name: "tls/machine-config-server.crt", mode: 0600},
		{path: "/opt/tectonic/tls/machine-config-server.key", name: "tls/machine-config-server.key", mode: 0600},
		{path: "/opt/tectonic/tls/root-ca.crt", name: "tls/root-ca.crt", mode: 0600},
		{path: "/opt/tectonic/tls/root-ca.key", name: "tls/root-ca.key", mode: 0600},
		{path: "/opt/tectonic/tls/service-account.key", name: "tls/service-account.key", mode: 0600},
		//{path: "/opt/tectonic/tls/service-account.pub", name: "", mode: 0600}, FIXME: do we need this?
		{path: "/opt/tectonic/tls/service-serving-ca.crt", name: "tls/service-serving-ca.crt", mode: 0600},
		{path: "/opt/tectonic/tls/service-serving-ca.key", name: "tls/service-serving-ca.key", mode: 0600},
		{path: "/usr/local/bin/bootkube.sh", name: "files/usr/local/bin/bootkube.sh", mode: 0555},
		{path: "/usr/local/bin/report-progress.sh", name: "files/usr/local/bin/report-progress.sh", mode: 0555},
		{path: "/usr/local/bin/tectonic.sh", name: "files/usr/local/bin/tectonic.sh", mode: 0555},
		{path: "/var/lib/kubelet/kubeconfig", name: "auth/kubeconfig-kubelet", mode: 0600},
	} {
		if entry.platform != "" && entry.platform != platform {
			continue
		}

		parents, err = asset.GetParents(ctx, getByName, entry.name)
		if err != nil {
			return nil, err
		}

		file := fileFromBytes(entry.path, entry.mode, parents[entry.name].Data)
		file.Append = entry.append
		if entry.user != "" {
			file.User = &types.NodeUser{Name: entry.user}
		}
		if entry.group != "" {
			file.Group = &types.NodeGroup{Name: entry.group}
		}
		config.Storage.Files = append(config.Storage.Files, file)
	}

	asset.Data, err = json.Marshal(config)
	if err != nil {
		return nil, err
	}

	return asset, nil
}

func pointerIgnitionRebuilder(role string) assets.Rebuild {
	return func(ctx context.Context, getByName assets.GetByString) (*assets.Asset, error) {
		asset := &assets.Asset{
			Name:          fmt.Sprintf("ignition/%s.ign", role),
			RebuildHelper: pointerIgnitionRebuilder(role),
		}

		parents, err := asset.GetParents(
			ctx,
			getByName,
			"base-domain",
			"cluster-name",
			"tls/root-ca.crt",
		)
		if err != nil {
			return nil, err
		}
		config := &types.Config{
			Ignition: types.Ignition{
				Version: types.MaxVersion.String(),
				Config: types.IgnitionConfig{
					Append: []types.ConfigReference{{
						Source: fmt.Sprintf("https://%s-api.%s:49500/config/%s", string(parents["cluster-name"].Data), string(parents["base-domain"].Data), role),
					}},
				},
				Security: types.Security{
					TLS: types.TLS{
						CertificateAuthorities: []types.CaReference{{
							Source: dataurl.EncodeBytes(parents["tls/root-ca.crt"].Data),
						}},
					},
				},
			},
		}

		// XXX: Remove this once MCO supports injecting SSH keys.
		parents, err = asset.GetParents(ctx, getByName, "ssh.pub")
		if err == nil {
			config.Passwd.Users = append(
				config.Passwd.Users,
				types.PasswdUser{
					Name: "core",
					SSHAuthorizedKeys: []types.SSHAuthorizedKey{
						types.SSHAuthorizedKey(parents["ssh.pub"].Data),
					},
				},
			)
		} else if !os.IsNotExist(err) {
			return nil, err
		}

		asset.Data, err = json.Marshal(config)
		if err != nil {
			return nil, err
		}

		return asset, nil
	}
}

func pointerIgnitionUserDataRebuilder(role string) assets.Rebuild {
	return func(ctx context.Context, getByName assets.GetByString) (*assets.Asset, error) {
		asset := &assets.Asset{
			Name:          fmt.Sprintf("manifests/99_openshift-cluster-api_%s-user-data-secret.yaml", role),
			RebuildHelper: pointerIgnitionUserDataRebuilder(role),
		}

		parentName := fmt.Sprintf("ignition/%s.ign", role)
		parents, err := asset.GetParents(
			ctx,
			getByName,
			parentName,
		)
		if err != nil {
			return nil, err
		}

		secret := corev1.Secret{
			TypeMeta: metav1.TypeMeta{
				Kind:       "Secret",
				APIVersion: "v1",
			},
			ObjectMeta: metav1.ObjectMeta{
				Name:      fmt.Sprintf("%s-user-data", role),
				Namespace: "openshift-cluster-api",
			},
			Data: map[string][]byte{
				"userData": parents[parentName].Data,
			},
			Type: corev1.SecretTypeOpaque,
		}

		asset.Data, err = yaml.Marshal(secret)
		if err != nil {
			return nil, err
		}

		return asset, nil
	}
}

func init() {
	Rebuilders["ignition/bootstrap.ign"] = bootstrapIgnRebuilder
	Rebuilders["ignition/master.ign"] = pointerIgnitionRebuilder("master")
	Rebuilders["ignition/worker.ign"] = pointerIgnitionRebuilder("worker")
	Rebuilders["manifests/99_openshift-cluster-api_master-user-data-secret.yaml"] = pointerIgnitionUserDataRebuilder("master")
	Rebuilders["manifests/99_openshift-cluster-api_worker-user-data-secret.yaml"] = pointerIgnitionUserDataRebuilder("worker")
}
