import os
import yaml

os.system('export INFRA_ID=$(jq -r .infraID metadata.json)')
# Read inventory.yaml for CiscoACI CNI variables
with open("inventory.yaml", 'r') as stream:
    try:
        inventory = yaml.safe_load(stream)['all']['hosts']['localhost']
    except yaml.YAMLError as exc:
        print(exc)

if 'neutron_network_mtu' not in inventory:
    neutron_network_mtu = "1500"
else:
    neutron_network_mtu = str(inventory['neutron_network_mtu'])

try:
    infra_vlan = str(inventory['infra_vlan'])
    node_interface = inventory['node_interface']
    opflex_interface = inventory['opflex_interface']
    os_cp_nodes_number = inventory['os_cp_nodes_number']
except:
    print("The inventory.yaml must have infra_vlan, node_interface and opflex_interface fields set")

os.system('./update-control.sh ' + str(os_cp_nodes_number) + ' ' + node_interface + ' ' + str(opflex_interface) + ' ' + infra_vlan + ' ' + str(neutron_network_mtu))
