#!/usr/bin/env bash

start_time=$(date -u +%s)
sleep_time=2s
elapsed=0
MAX_WAIT=300
BOOTSTRAPIP=''
EXEC_CMD=''

await_trail=''

eval "$(ssh-agent -s)" > /dev/null || exit 1
ssh-add "${HOME}/.ssh/id_rsa" > /dev/null 2>&1 || exit 1

# First argument is assigned to MAX_WAIT
if [[ -n $1 ]]; then
    MAX_WAIT=$1
fi

# Second argument is assigned to EXEC_CMD
if [[ -n $2 ]]; then
    EXEC_CMD=$2
fi

while [[ elapsed -lt MAX_WAIT ]]
do
    if [[ -n $EXEC_CMD ]]; then
        await_trail="${await_trail}."
        if [[ ${#await_trail} -gt 3 ]]; then
            await_trail=''
        fi
        echo -ne "\\rAwaiting cluster availability${await_trail}    \\r"
    fi

    if [[ -z $BOOTSTRAPIP ]]; then

        TEMPBOOTSTRAPIP=$(virsh --connect qemu+tcp://192.168.122.1/system domifaddr bootstrap 2> /dev/null | awk '/192/{print $4}')
        if ! [[ -z $TEMPBOOTSTRAPIP ]]; then
            BOOTSTRAPIP=${TEMPBOOTSTRAPIP::${#TEMPBOOTSTRAPIP}-3}
        fi
    fi

    if [[ -n $BOOTSTRAPIP ]]; then
        msg=$(ssh -oStrictHostKeyChecking=no core@"${BOOTSTRAPIP}" journalctl -n 1 -u bootkube.service -u tectonic 2> /dev/null)
        if echo "$msg" | grep 'Tectonic installation is done'; then
            echo ''
            if ! [[ -z $EXEC_CMD ]]; then
                bash -c "${EXEC_CMD}"
                exit 0
            fi
            bash -c "ssh -oStrictHostKeyChecking=no core@${BOOTSTRAPIP} journalctl -f -u bootkube -u tectonic"
            exit 0
        fi
        if [[ -z $EXEC_CMD ]]; then
            echo "${msg}"
        fi
    fi
    elapsed=$(($(date -u +%s) - start_time))
    sleep $sleep_time
done
echo -e "\\nWatch stopped after elapsed time: ${elapsed}"
exit 1
