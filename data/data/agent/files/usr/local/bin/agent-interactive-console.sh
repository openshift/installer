#!/bin/bash

# shellcheck disable=SC1091
source /etc/assisted/agent-installer.env

function has_registry_connectivity() {
    (>&2 echo "INFO: Checking OpenShift release image at $RELEASE_IMAGE is pullable.")
    (>&2 podman pull "$RELEASE_IMAGE")
    registry_check_status_code=$?
    if [[ "$registry_check_status_code" -ne 0 ]]; then
      (>&2 echo "WARNING: Unable to pull release image.")
      echo 0
    else
      (>&2 echo "INFO: Release image is pullable.")
      echo 1
    fi
}

function has_rendezvous_host_connectivity() {
    if (>&2 ping -c 4 "$NODE_ZERO_IP"); then
	  (>&2 echo "INFO: Successfully pinged rendezvous host at ${NODE_ZERO_IP}.")
	  echo 1
	else
	  (>&2 echo "WARNING: Failed to ping rendezvous host at ${NODE_ZERO_IP}.")
	  echo 0
	fi
}

function has_connectivity() {
    (>&2 echo "INFO: Checking connectivity.")
    problem_found=0
    if [[ $(has_registry_connectivity) -eq 0 ]]; then
      problem_found=1
    fi
    if [[ $(has_rendezvous_host_connectivity) -eq 0 ]]; then
      problem_found=1
    fi

    if [[ $problem_found -eq 1 ]]; then
      (>&2 echo "WARNING: Connectivity problem found.")
      echo 0
    else
      (>&2 echo "INFO: Host has connectivity to release image and rendezvous host.")
      echo 1
    fi
}

should_configure=1
while [ $should_configure -eq 1 ]; do
  while true; do
    if [[ $(has_connectivity) -eq 1 ]]; then
	    should_configure=0
    fi
    if [ $should_configure -eq 1 ]; then
      # Case 1: Connectivity checks failed, run the tui to allow user to update
      # network configuration.
      #
      # TODO: execute tui and remove sleep. If the tui does connectivity checks,
      # perhaps we can exit from the loop upon the tui exiting. For now, we always
      # exit to allow the automated flow to continue to preserve the current
      # behavior.
      #/usr/local/bin/agent-tui
      echo "TODO: execute agent-tui and remove 60s sleep"
      sleep 60
      # break here does not work, need exit
      exit 0
    else
      # Case 2: Connectivity checks passed. Give users 60s to use the tui
      # if they would like to make a change. After 60s, the prompt disappears
      # and automated flow will start.
      echo "INFO: Boot will continue in 60s if the next prompt goes unanswered."
      read -t 60 -e -p "Do you wish to enter NetworkManager terminal user interface to change the network configuration? (Y)es/(N)o: " ANSWER
      if [[ "$ANSWER" =~ ^[Yy](es)?$ ]]; then
	      # TODO: execute tui
        #/usr/local/bin/agent-tui
        # TODO: If the tui does connectivity checks, then we may not need to force
        # this script to do its own connectivty checks again. The next line may be
        # removed and can break from the loop.
        should_configure=1
      else
        # The prompt timed out. Proceed with automated flow.
        break
      fi
    fi
  done
done
