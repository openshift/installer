#!/usr/bin/env bash

printf "hostname: %s\n" "$(hostname)"

MASTER_FILE=/run/metadata/master
if [ -f $MASTER_FILE ]
then
  printf "master: %s\n" "$(cat $MASTER_FILE)"
else
  printf "%s does not exist!\n" $MASTER_FILE
fi

TECTONIC_DIR=/opt/tectonic
if [ -d $TECTONIC_DIR ]
then
  printf "ls %s: %s\n\n" $TECTONIC_DIR "$(ls -la $TECTONIC_DIR)"
else
  printf "%s does not exist!\n\n" $TECTONIC_DIR
fi

printf "Init assets:\n"
journalctl -u init-assets --no-tail | cat
printf "\n\n"

printf "Bootkube:\n"
journalctl -u bootkube --no-tail | cat
printf "\n\n"

printf "Tectonic:\n"
journalctl -u tectonic --no-tail | cat
printf "\n\n"
