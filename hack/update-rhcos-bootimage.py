#!/usr/bin/env python3
# Usage: ./hack/update-rhcos-bootimage.py https://releases-art-rhcos.svc.ci.openshift.org/storage/releases/ootpa/410.8.20190401.0/meta.json amd64
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
