import json
import base64
import sys
import os
from dotmap import DotMap
import yaml
from collections import OrderedDict 

resource_group = sys.argv[1]

with open('./gw/manifests/cloud-provider-config.yaml') as file:
      yamlx = yaml.load(file)
      jsondata = yamlx['data']['config']
      jsonx = json.loads(jsondata,object_pairs_hook=OrderedDict)
      config = DotMap(jsonx)
      config.resourceGroup = resource_group
      config.vnetName = "openshiftVnet"
      config.vnetResourceGroup = resource_group
      config.subnetName = "masterSubnet"
      config.securityGroupName = "master1nsg"
      config.routeTableName = ""
      jsondata = json.dumps(dict(**config.toDict()),indent='\t')
      jsonstr = str(jsondata)
      yamlx['data']['config'] =   jsonstr + '\n'
      yamlx['metadata']['creationTimestamp'] = None
      yamlstr = yaml.dump(yamlx, default_style='\"', width=4096)
      yamlstr = yamlstr.replace('!!null "null"','null')
      with open('./gw/manifests/cloud-provider-config.yaml', 'w') as outfile:
          outfile.write(yamlstr)

with open('./gw/manifests/cluster-infrastructure-02-config.yml') as file:
      yamlx = yaml.load(file)
      yamlx['status']['platformStatus']['azure']['resourceGroupName'] = resource_group   
      with open('./gw/manifests/cluster-infrastructure-02-config.yml','w') as outfile:
          yaml.dump(yamlx, outfile, default_flow_style=False)

dnsyml = "gw/manifests/cluster-dns-02-config.yml"
data = yaml.load(open(dnsyml));
del data["spec"]["publicZone"];
del data["spec"]["privateZone"];
open(dnsyml, "w").write(yaml.dump(data, default_flow_style=False))
