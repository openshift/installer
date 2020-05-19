#!/bin/sh

printf 'Warning: This will destroy effectively all libvirt resources\nContinue [yN]? '
read -r CONTINUE
if test "${CONTINUE}" != y -a "${CONTINUE}" != Y
then
	  echo 'Aborted' >&2
	  exit 1
fi

CONNECT="${CONNECT:=qemu:///system}"

run()
{
	echo "$*"
	"$@"
}

for DOMAIN in $(virsh -c "${CONNECT}" list --all --name)
do
	run virsh -c "${CONNECT}" destroy "${DOMAIN}"
	run virsh -c "${CONNECT}" undefine "${DOMAIN}"
done

for POOL in $(virsh -c "${CONNECT}" pool-list --all --name)
do
	virsh -c "${CONNECT}" vol-list "${POOL}" | tail -n +3 | while read -r VOLUME _
	do
		if test -z "${VOLUME}"
		then
			continue
		fi
		run virsh -c "${CONNECT}" vol-delete --pool "${POOL}" "${VOLUME}"
	done
	run virsh -c "${CONNECT}" pool-destroy "${POOL}"
	run virsh -c "${CONNECT}" pool-undefine "${POOL}"
done

for NET in $(virsh -c "${CONNECT}" net-list --all --name)
do
	if test "${NET}" = default
	then
		continue
	fi
	run virsh -c "${CONNECT}" net-destroy "${NET}"
	run virsh -c "${CONNECT}" net-undefine "${NET}"
done
