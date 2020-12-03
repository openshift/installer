clean:
	./hack/virt-delete-sno.sh || true
	rm -rf mydir

generate:
	mkdir -p mydir
	cp ./install-config.yaml mydir/
	OPENSHIFT_INSTALL_RELEASE_IMAGE_OVERRIDE="quay.io/openshift-release-dev/ocp-release:4.6.1-x86_64" ./bin/openshift-install create ignition-configs --dir=mydir
embed:
	sudo podman run --pull=always --privileged --rm -v /dev:/dev -v /run/udev:/run/udev -v .:/data -w /data     quay.io/coreos/coreos-installer:release iso ignition embed /data/installer-image.iso -f --ignition-file /data/mydir/bootstrap.ign -o /data/installer-SNO-image.iso
	mv installer-SNO-image.iso /tmp/images/installer-SNO-image.iso

start-iso:
	./hack/virt-install-sno-iso-ign.sh
start:
	./hack/virt-install-sno-ign.sh ./mydir/bootstrap.ign

network:
	./hack/virt-create-net.sh

ssh:
	ssh -o UserKnownHostsFile=/dev/null -o StrictHostKeyChecking=no core@192.168.126.10

image:
	curl -O -L https://releases-art-rhcos.svc.ci.openshift.org/art/storage/releases/rhcos-4.6/46.82.202008181646-0/x86_64/rhcos-46.82.202008181646-0-qemu.x86_64.qcow2.gz
	mv rhcos-46.82.202008181646-0-qemu.x86_64.qcow2.gz /tmp
	sudo gunzip /tmp/rhcos-46.82.202008181646-0-qemu.x86_64.qcow2.gz
