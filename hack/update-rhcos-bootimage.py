#!/usr/bin/env python3
# As of 4.8 we are aiming to switch to stream metadata:
# https://github.com/openshift/enhancements/pull/679
# That transition hasn't yet fully completed; there are two copies of the
# RHCOS metadata:
# 
# - data/data/rhcos-4.8.json (stream format, 4.8+)
# - data/data/rhcos-$arch.json (openshift/installer specific, 4.7 and below)
# 
# See https://github.com/coreos/coreos-assembler/pull/2000 in particular.
# 
# The initial file data/data/rhcos-4.8 was generated this way:
# 
# $ plume cosa2stream --name rhcos-4.8 --distro rhcos  x86_64=48.83.202102230316-0 s390x=47.83.202102090311-0 ppc64le=47.83.202102091015-0 > data/data/rhcos-4.8.json
# 
# To update the bootimage for one or more architectures, use e.g.
# 
# $ plume cosa2stream --target data/data/rhcos-4.8.json --distro rhcos  x86_64=48.83.202102230316-0 s390x=47.83.202102090311-0 ppc64le=47.83.202102091015-0
#
# To update the legacy metadata, use:
# Usage: ./hack/update-rhcos-bootimage.py https://releases-art-rhcos.svc.ci.openshift.org/art/storage/releases/rhcos-4.6/46.82.202008260918-0/x86_64/meta.json amd64
import codecs,os,sys,json,argparse
import urllib.parse
import urllib.request

# An app running in the CI cluster exposes this public endpoint about ART RHCOS
# builds.  Do not try to e.g. point to RHT-internal endpoints.
RHCOS_RELEASES_APP = 'https://releases-art-rhcos.svc.ci.openshift.org'

parser = argparse.ArgumentParser()
parser.add_argument("meta", action='store')
parser.add_argument("arch", action='store', choices=['amd64', 's390x', 'ppc64le'])
args = parser.parse_args()

metadata_dir = os.path.join(os.path.dirname(sys.argv[0]), "../data/data")

if not args.meta.startswith(RHCOS_RELEASES_APP):
    raise SystemExit("URL must start with: " + RHCOS_RELEASES_APP)

with urllib.request.urlopen(args.meta) as f:
    string_f = codecs.getreader('utf-8')(f)  # support for Python < 3.6
    meta = json.load(string_f)
newmeta = {}
for k in ['images', 'buildid', 'oscontainer',
          'ostree-commit', 'ostree-version',
          'azure', 'gcp']:
    if meta.get(k):
        newmeta[k] = meta[k]
if meta.get(k):
    newmeta['amis'] = {
        entry['name']: {
            'hvm': entry['hvm'],
        }
        for entry in meta['amis']
    }
newmeta['baseURI'] = urllib.parse.urljoin(args.meta, '.')

with open(os.path.join(metadata_dir, 'rhcos-{}.json'.format(args.arch)), 'w') as f:
    json.dump(newmeta, f, sort_keys=True, indent=4)

# Continue to populate the legacy metadata file because there are still
# processes consuming this file directly. This normally could just be a symlink
# but some of these processes reference raw.githubusercontent.com which doesn't
# follow symlinks.
if args.arch == 'amd64':
    with open(os.path.join(metadata_dir, "rhcos.json"), 'w') as f:
        json.dump(newmeta, f, sort_keys=True, indent=4)
