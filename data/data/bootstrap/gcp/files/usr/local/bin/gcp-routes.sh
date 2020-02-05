#!/bin/bash

declare -A routes
curler() {
  curl --silent -L -H "Metadata-Flavor: Google" "http://metadata.google.internal/computeMetadata/v1/instance/${1}"
}

get_ifname() {
  sysfs_path="/sys/class/net"
  for dev in $(find ${sysfs_path} -maxdepth 1  -mindepth 1);
  do
      local mac=$(<${dev}/address);
      local name="$(basename ${dev})"
      if [ "${mac}" == "${1}" ];
      then
          echo "${name}"
          return;
      fi
  done
}

set_routes() {
  local dev="${1}"
  read -a dev_routes <<< "${routes[$dev]}"
  for cur_route in $(ip route show dev ${dev} table local proto 66 | awk '{print$2}');
  do
      if [[ ! "${dev_routes[@]}" =~ "${cur_route}" ]];
      then
          echo "Removing stale forwarded IP ${cur_route}/32"
          ip route del ${cur_route}/32 dev ${dev} table local proto 66
      fi
  done
  for route in ${dev_routes[@]}
  do
      ip route replace to local ${route} dev $dev proto 66
  done
  unset dev_routes
}

del_routes() {
  local dev="${1}"
  read -a dev_routes <<< "${routes[$dev]}"
  for cur_route in $(ip route show dev ${dev} table local proto 66 | awk '{print$2}');
  do
      if [[ "${dev_routes[@]}" =~ "${cur_route}" ]];
      then
          echo "Removing forwarded IP ${cur_route}/32"
          ip route del ${cur_route}/32 dev ${dev} table local proto 66
      fi
  done
  unset dev_routes
}

run() {
  net_path="network-interfaces/"
  for vif in $(curler ${net_path}); do
      hw_addr=$(curler "${net_path}${vif}mac")
      fwip_path="${net_path}${vif}forwarded-ips/"
      dev_name="$(get_ifname ${hw_addr})"
      for level in $(curler ${fwip_path})
      do
          for fwip in $(curler ${fwip_path}${level})
          do
              echo "Processing route for NIC ${vif}${hw_addr} as ${dev_name} for ${fwip}"
              routes[$dev_name]+="${fwip} "
          done
      done
      $"${1}" ${dev_name}

      routes[$dev_name]=""
      unset hw_addr
      unset fwip_path
      unset dev_name
  done
}

case "$1" in
  start)
    while :; do
      run set_routes
      sleep 30
    done
    ;;
  cleanup)
    run del_routes
    ;;
  *)
    echo $"Usage: $0 {start|cleanup}"
    exit 1

esac
