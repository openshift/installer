#!/bin/bash

function connectivity_check() {
	declare target

	target="$1"

	(>&2 echo "Pinging ${target} to check the validity of the current network configuration")
	if ping -c 4 "$NODE_ZERO_IP"; then
		(>&2 echo "Successfully pinged ${target}.")
		return 0
	else
		(>&2 echo "Failed to ping ${target}.")
		return 1
	fi
}

function print_addressing_info() {
	declare with_prefix

	with_prefix="${1:-true}"
	echo "Currently configured addresses"
	if $with_prefix; then
		ip -j a show  | jq -r '. | map({name: .ifname, addresses: [.addr_info[] | .local + "/" + (.prefixlen | tostring)]}) | .[] | select(.addresses | length > 0) | select( .name != "lo") | .name + ": " + (.addresses | join(" "))'
	else
		ip -j a show  | jq -r '. | map({name: .ifname, addresses: [.addr_info[] | .local]}) | .[] | select(.addresses | length > 0) | select( .name != "lo") | .name + ": " + (.addresses | join(" "))'
	fi
}

NODE_ZERO=false
while true; do
	read -e -p "Do you wish for this node to be the one that runs the installation service (Only one node may perform this function)? (Y)es/(N)o: " NODE0_ANSWER
	if [[ "$NODE0_ANSWER" =~ ^[Yy](es)?$ ]]; then
		NODE_ZERO=true
		break
	elif [[ "$NODE0_ANSWER" =~ ^[Nn](o)?$ ]]; then
		(>&2 echo "This node will not run the installation service.")
		break
	else
		(>&2 echo "Invalid answer.")
	fi
done

if $NODE_ZERO; then
	print_addressing_info
	while true; do
		read -e -p "Do you wish to enter NetworkManager terminal user interface to change the network configuration? (Y)es/(N)o: " ANSWER
		if [[ "$ANSWER" =~ ^[Yy](es)?$ ]]; then
			nmtui
			print_addressing_info
		elif [[ "$ANSWER" =~ ^[Nn](o)?$ ]]; then
			echo "Other nodes should give one of the following addresses as the rendezvous IP"
			print_addressing_info false
			break
		else
			(>&2 echo "Answer not in (Y)es/(N)o")
		fi
	done
	exit 0
else
	CONFIRMED=false
	until $CONFIRMED; do
		read -e -p "Please, enter the rendezvous IP address: " NODE_ZERO_IP
		read -e -p "Please, enter the rendezvous IP address again to confirm: " CONFIRMATION
		if [[ "$NODE_ZERO_IP" == "$CONFIRMATION" ]]; then
			CONFIRMED=true
		else
			(>&2 echo "Confirmation failed")
		fi
	done
fi

SHOULD_CONFIGURE=true
while $SHOULD_CONFIGURE; do
	while true; do
		if connectivity_check "$NODE_ZERO_IP"; then
			SHOULD_CONFIGURE=false
		fi
		read -e -p "Do you wish to enter NetworkManager terminal user interface to change the network configuration? (Y)es/(N)o: " ANSWER
		if [[ "$ANSWER" =~ ^[Yy](es)?$ ]]; then
			SHOULD_CONFIGURE=true  # Setting to true in case the connectivity_check had passed but we chose to configure anyway
			nmtui
			if connectivity_check "$NODE_ZERO_IP"; then
				SHOULD_CONFIGURE=false
			fi
			break
		elif [[ "$ANSWER" =~ ^[Nn](o)?$ ]]; then
			if ! $SHOULD_CONFIGURE; then
				break
			fi
			(>&2 echo "No valid network configuration found")
		else
			(>&2 echo "Answer not in (Y)es/(N)o")
		fi
	done
done
