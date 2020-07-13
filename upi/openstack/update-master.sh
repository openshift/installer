for index in $(seq 0 3); do
    MASTER_HOSTNAME="$INFRA_ID-master-$index\n"
    sudo python3 -c "import base64, json, sys;
ignition = json.load(sys.stdin);
files = ignition['storage'].get('files', []);
files.append({'path': '/etc/hostname', 'mode': 420, 'contents': {'source': 'data:text/plain;charset=utf-8;base64,' + base64.standard_b64encode(b'$MASTER_HOSTNAME').decode().strip(), 'verification': {}}, 'filesystem': 'root'});

ifcfg_ens3 = 'TYPE=Ethernet\nDEVICE=ens3\nONBOOT=yes\nBOOTPROTO=dhcp\nDEFROUTE=yes\nPROXY_METHOD=none\nBROWSER_ONLY=no\nMTU=1700\nIPV4_FAILURE_FATAL=no\nIPV6INIT=no'.encode();

ifcfg_ens3_b64 = base64.standard_b64encode(ifcfg_ens3).decode().strip();

files.append({'path': '/etc/sysconfig/network-scripts/ifcfg-ens3','mode': 420,'contents': {'source': 'data:text/plain;charset=utf-8;base64,' + ifcfg_ens3_b64,'verification': {}},'filesystem': 'root',});


ifcfg_ens4 = 'TYPE=Ethernet\nDEVICE=ens4\nONBOOT=yes\nBOOTPROTO=dhcp\nDEFROUTE=yes\nPROXY_METHOD=none\nBROWSER_ONLY=no\nMTU=1700\nIPV4_FAILURE_FATAL=no\nIPV6INIT=no'.encode();

ifcfg_ens4_b64 = base64.standard_b64encode(ifcfg_ens4).decode().strip()

files.append({'path': '/etc/sysconfig/network-scripts/ifcfg-ens4','mode': 420,'contents': {'source': 'data:text/plain;charset=utf-8;base64,' + ifcfg_ens4_b64,'verification': {}},'filesystem': 'root',});


opflex_conn = 'VLAN=yes\nTYPE=Vlan\nPHYSDEV=ens4\nVLAN_ID=4093\nREORDER_HDR=yes\nGVRP=no\nMVRP=no\nPROXY_METHOD=none\nBROWSER_ONLY=no\nBOOTPROTO=dhcp\nDEFROUTE=no\nIPV4_FAILURE_FATAL=no\nIPV6INIT=no\nNAME=opflex-conn\nDEVICE=ens4.4093\nONBOOT=yes\nMTU=1600'.encode();

opflex_conn_b64 = base64.standard_b64encode(opflex_conn).decode().strip();

files.append({'path': '/etc/sysconfig/network-scripts/ifcfg-opflex-conn','mode': 420,'contents': {'source': 'data:text/plain;charset=utf-8;base64,' + opflex_conn_b64,'verification': {}},'filesystem': 'root',});

route_opflex_conn = 'ADDRESS0=224.0.0.0\nNETMASK0=240.0.0.0\nMETRIC0=1000'.encode();

route_opflex_conn_b64 = base64.standard_b64encode(route_opflex_conn).decode().strip();

files.append({'path': '/etc/sysconfig/network-scripts/route-opflex-conn','mode': 420,'contents': {'source': 'data:text/plain;charset=utf-8;base64,' + route_opflex_conn_b64,'verification': {}},'filesystem': 'root',})

ignition['storage']['files'] = files;
json.dump(ignition, sys.stdout)" <master.ign >"$INFRA_ID-master-$index-ignition.json"
done
