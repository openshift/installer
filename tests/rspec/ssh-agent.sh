#!/bin/bash
KEY_FILE=${HOME}/.ssh/id_rsa
ssh-keygen -f "${KEY_FILE}" -t rsa -N ''
ssh-add "${KEY_FILE}"
export TF_VAR_tectonic_azure_ssh_key="${KEY_FILE}.pub"
bash -i