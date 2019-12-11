#!/usr/bin/env python3

import json
import base64
import sys
from dotmap import DotMap

bootstrapIgnitionURL = sys.argv[1]

ign = DotMap()
config = DotMap()
ign.ignition.version = "2.2.0"
config.replace.source = bootstrapIgnitionURL
ign.ignition.config = config

ignstr = json.dumps(dict(**ign.toDict()))

print(base64.b64encode(ignstr.encode()).decode())
