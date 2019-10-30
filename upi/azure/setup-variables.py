import json
import base64
import sys
import os
from dotmap import DotMap

url = sys.argv[1]
region = sys.argv[2]

ign = DotMap()
config = DotMap()
ign.ignition.version = "2.2.0"
config.replace.source = url
ign.ignition.config = config

ignstr = json.dumps(dict(**ign.toDict()))

sshpath = os.path.expanduser('~/.ssh/id_rsa.pub')
with open(sshpath,"r") as sshFile:
     sshkey = sshFile.read()
with open("gw/master.ign","r") as ignFile:
    master_ignition = json.load(ignFile)
with open("gw/worker.ign","r") as ignFile:
    worker_ignition = json.load(ignFile)
with open("azuredeploy.parameters.json", "r") as jsonFile:
    data = DotMap(json.load(jsonFile))


data.parameters.BootstrapIgnition.value =  base64.b64encode(ignstr.encode()).decode()
data.parameters.MasterIgnition.value =     base64.b64encode(json.dumps(master_ignition).encode()).decode()
data.parameters.WorkerIgnition.value =     base64.b64encode(json.dumps(worker_ignition).encode()).decode()
data.parameters.sshKeyData.value     =     sshkey.rstrip()
data.parameters.image.value          =     'https://sa' + region + '.blob.core.windows.net/vhd/rhcos.vhd'

jsondata = dict(**data.toDict())
with open("runit.parameters.json", "w") as jsonFile:
    json.dump(jsondata,jsonFile)

