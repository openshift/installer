import base64
import json
import os

with open('bootstrap.ign', 'r') as f:
    ignition = json.load(f)

files = ignition['storage'].get('files', [])

infra_id = os.environ.get('INFRA_ID', 'openshift').encode()
hostname_b64 = base64.standard_b64encode(infra_id + b'-bootstrap\n').decode().strip()
files.append(
{
    'path': '/etc/hostname',
    'mode': 420,
    'contents': {
        'source': 'data:text/plain;charset=utf-8;base64,' + hostname_b64,
        'verification': {}
    },
    'filesystem': 'root',
})

dhcp_client_conf_b64 = base64.standard_b64encode(b'[main]\ndhcp=dhclient\n').decode().strip()
files.append(
{
    'path': '/etc/NetworkManager/conf.d/dhcp-client.conf',
    'mode': 420,
    'contents': {
        'source': 'data:text/plain;charset=utf-8;base64,' + dhcp_client_conf_b64,
        'verification': {}
        },
    'filesystem': 'root',
})

dhclient_cont_b64 = base64.standard_b64encode(b'send dhcp-client-identifier = hardware;\nprepend domain-name-servers 127.0.0.1;\n').decode().strip()
files.append(
{
    'path': '/etc/dhcp/dhclient.conf',
    'mode': 420,
    'contents': {
        'source': 'data:text/plain;charset=utf-8;base64,' + dhclient_cont_b64,
        'verification': {}
        },
    'filesystem': 'root'
})

ifcfg_ens3 = """TYPE=Ethernet
DEVICE=ens3
ONBOOT=yes
BOOTPROTO=dhcp
DEFROUTE=yes
PROXY_METHOD=none
BROWSER_ONLY=no
MTU=1700
IPV4_FAILURE_FATAL=no
IPV6INIT=no""".encode()

ifcfg_ens3_b64 = base64.standard_b64encode(ifcfg_ens3).decode().strip()

files.append(
{
    'path': '/etc/sysconfig/network-scripts/ifcfg-ens3',
    'mode': 420,
    'contents': {
        'source': 'data:text/plain;charset=utf-8;base64,' + ifcfg_ens3_b64,
        'verification': {}
    },
    'filesystem': 'root',
})

ifcfg_ens4 = """TYPE=Ethernet
DEVICE=ens4
ONBOOT=yes
BOOTPROTO=dhcp
DEFROUTE=no
PROXY_METHOD=none
BROWSER_ONLY=no
MTU=1700
IPV4_FAILURE_FATAL=no
IPV6INIT=no""".encode()

ifcfg_ens4_b64 = base64.standard_b64encode(ifcfg_ens4).decode().strip()

files.append(
{
    'path': '/etc/sysconfig/network-scripts/ifcfg-ens4',
    'mode': 420,
    'contents': {
        'source': 'data:text/plain;charset=utf-8;base64,' + ifcfg_ens4_b64,
        'verification': {}
    },
    'filesystem': 'root',
})


opflex_conn = """VLAN=yes
TYPE=Vlan
PHYSDEV=ens4
VLAN_ID=4093
REORDER_HDR=yes
GVRP=no
MVRP=no
PROXY_METHOD=none
BROWSER_ONLY=no
BOOTPROTO=dhcp
DEFROUTE=no
IPV4_FAILURE_FATAL=no
IPV6INIT=no
NAME=opflex-conn
DEVICE=ens4.4093
ONBOOT=yes
MTU=1600""".encode()

opflex_conn_b64 = base64.standard_b64encode(opflex_conn).decode().strip()

files.append(
{
    'path': '/etc/sysconfig/network-scripts/ifcfg-opflex-conn',
    'mode': 420,
    'contents': {
        'source': 'data:text/plain;charset=utf-8;base64,' + opflex_conn_b64,
        'verification': {}
    },
    'filesystem': 'root',
})

route_opflex_conn = """ADDRESS0=224.0.0.0
NETMASK0=240.0.0.0
METRIC0=1000""".encode()

route_opflex_conn_b64 = base64.standard_b64encode(route_opflex_conn).decode().strip()

files.append(
{
    'path': '/etc/sysconfig/network-scripts/route-opflex-conn',
    'mode': 420,
    'contents': {
        'source': 'data:text/plain;charset=utf-8;base64,' + route_opflex_conn_b64,
        'verification': {}
    },
    'filesystem': 'root',
})


ignition['storage']['files'] = files;

with open('bootstrap.ign', 'w') as f:
    json.dump(ignition, f)

