#!/usr/bin/env python3
# This script updates the legacy metadata.  We hope to remove it soon.
# See docs/dev/pinned-coreos.md for more information.
import codecs,os,sys,json,argparse
import urllib.parse
import urllib.request

# An app running in the CI cluster exposes this public endpoint about ART RHCOS
# builds.  Do not try to e.g. point to RHT-internal endpoints.
RHCOS_RELEASES_APP = 'https://rhcos-redirector.apps.art.xq1c.p1.openshiftapps.com'

parser = argparse.ArgumentParser()
parser.add_argument("meta", action='store')
parser.add_argument("arch", action='store', choices=['amd64', 's390x', 'ppc64le', 'aarch64'])
args = parser.parse_args()

metadata_dir = os.path.join(os.path.dirname(sys.argv[0]), "../data/data")

if not args.meta.startswith(RHCOS_RELEASES_APP):
    raise SystemExit("URL must start with: " + RHCOS_RELEASES_APP)

with urllib.request.urlopen(args.meta) as f:
    string_f = codecs.getreader('utf-8')(f)  # support for Python < 3.6
    meta = json.load(string_f)
newmeta = {}

# x86_64 boot image bumps should contain more than just the ostree tarball and
# qemu qcow2 image. check the contents of the images dict for well known artifacts.
if args.arch == 'amd64':
    for k in ['aws', 'gcp', 'metal', 'metal4k']:
        if k not in meta['images'].keys():
            raise SystemExit("You appear to be missing most of the x86_64 boot "
                             + "image artifacts. Check the meta.json for this "
                             + "RHCOS version. Also check that boot images were "
                             + "enabled in the RHCOS build pipeline.")
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
