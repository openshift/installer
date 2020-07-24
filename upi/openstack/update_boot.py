import base64
import json
import os
import tarfile
import yaml

with open('bootstrap.ign', 'r') as f:
    ignition = json.load(f)

files = ignition['storage'].get('files', [])

# Read inventory.yaml for CiscoACI CNI variables
with open("inventory.yaml", 'r') as stream:
    try:
        inventory = yaml.safe_load(stream)['all']['hosts']['localhost']
    except yaml.YAMLError as exc:
        print(exc)

# Get accprovision tar path from inventory
try:
    acc_provision_tar = inventory['acc_provision_tar']
except:
    print("inventory.yaml should have acc_provision_tar field")

# Read acc-provision for vlan values
extract_to = './accProvisionTar'
tar = tarfile.open(acc_provision_tar, "r:gz")
tar.extractall(extract_to)
tar.close()

data = ''
for filename in os.listdir(extract_to):
    if 'ConfigMap-aci-containers-config' in filename:
        filepath = "%s/%s" % (extract_to, filename)
        with open(filepath, 'r') as stream:
            try:
                data = yaml.safe_load(stream)['data']['host-agent-config']
            except yaml.YAMLError as exc:
                print(exc)

# Extract host-agent-config and obtain vlan values
try:
    json_data = json.loads(data)
    aci_infra_vlan = json_data['aci-infra-vlan']
    service_vlan = json_data['service-vlan']
except:
    print("Couldn't extract host-agent-config from aci-containers ConfigMap")

# Set infra_vlan field in inventory.yaml using accprovision tar value
try:
    with open("inventory.yaml", 'r') as stream:
        cur_yaml = yaml.safe_load(stream)
        cur_yaml['all']['hosts']['localhost']['infra_vlan'] = aci_infra_vlan

    if cur_yaml:
        with open("inventory.yaml",'w') as yamlfile:
           yaml.safe_dump(cur_yaml, yamlfile)
except:
    print("Unable to edit inventory.yaml")
try:
    node_interface = inventory['node_interface']
    opflex_interface = inventory['opflex_interface']
except:
    print("The inventory.yaml must have node_interface and opflex_interface fields set")

if 'neutron_network_mtu' not in inventory:
    neutron_network_mtu = "1500"
else:
    neutron_network_mtu = str(inventory['neutron_network_mtu'])

infra_vlan = str(aci_infra_vlan)
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

ca_cert_path = os.environ.get('OS_CACERT', '')
if ca_cert_path:
    with open(ca_cert_path, 'r') as f:
        ca_cert = f.read().encode()
        ca_cert_b64 = base64.standard_b64encode(ca_cert).decode().strip()

    files.append(
    {
        'path': '/opt/openshift/tls/cloud-ca-cert.pem',
        'mode': 420,
        'contents': {
            'source': 'data:text/plain;charset=utf-8;base64,' + ca_cert_b64,
            'verification': {}
        },
        'filesystem': 'root',
    })

ifcfg_ens3 = ("""TYPE=Ethernet
DEVICE=""" + node_interface + """
ONBOOT=yes
BOOTPROTO=dhcp
DEFROUTE=yes
PROXY_METHOD=none
BROWSER_ONLY=no
MTU=""" + neutron_network_mtu + """
IPV4_FAILURE_FATAL=no
IPV6INIT=no""").encode()

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

ifcfg_ens4 = ("""TYPE=Ethernet
DEVICE=""" + opflex_interface + """
ONBOOT=yes
BOOTPROTO=dhcp
DEFROUTE=no
PROXY_METHOD=none
BROWSER_ONLY=no
MTU=""" + neutron_network_mtu + """
IPV4_FAILURE_FATAL=no
IPV6INIT=no""").encode()

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

opflex_conn = ("""VLAN=yes
TYPE=Vlan
PHYSDEV=""" + opflex_interface + """
VLAN_ID=""" + infra_vlan + """
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
DEVICE=""" + opflex_interface + """.""" + infra_vlan + """
ONBOOT=yes
MTU=""" + neutron_network_mtu).encode()

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

