#!/usr/bin/env python3
# Usage: ./hack/update-fcos-bootimage.py https://builds.coreos.fedoraproject.org/prod/streams/next/builds/32.20200517.1.0/x86_64/meta.json amd64
import codecs,os,sys,json,argparse
import urllib.parse
import urllib.request

# An app running in the CI cluster exposes this public endpoint about ART RHCOS
# builds.  Do not try to e.g. point to RHT-internal endpoints.
RHCOS_RELEASES_APP = 'https://builds.coreos.fedoraproject.org'

parser = argparse.ArgumentParser()
parser.add_argument("meta", action='store')
args = parser.parse_args()

metadata_dir = os.path.join(os.path.dirname(sys.argv[0]), "../data/data")

if not args.meta.startswith(RHCOS_RELEASES_APP):
    raise SystemExit("URL must start with: " + RHCOS_RELEASES_APP)

with urllib.request.urlopen(args.meta) as f:
    string_f = codecs.getreader('utf-8')(f)  # support for Python < 3.6
    meta = json.load(string_f)
newmeta = {}

# OKD doesn't yet support several platforms, so some image metadata should be skipped
include_images = ['aws', 'azure', 'gcp', 'metal', 'openstack', 'ostree', 'qemu', 'vmware']
rename_images = {
    "live-initramfs": "initramfs",
    "live-iso": "iso",
    "live-kernel": "kernel",
}

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

# Filter images
imgs = meta['images']
for k in list(imgs):
    if k in rename_images.keys():
        imgs[rename_images[k]] = imgs[k]
        del imgs[k]
        continue
    if k not in include_images:
        del imgs[k]

newmeta['images'] = imgs

# There are no official Azure images yet (https://github.com/coreos/fedora-coreos-tracker/issues/148)
# Installer has a workaround for that, so this section can be synthetized
newmeta['azure'] = {
    'image': newmeta['images']['azure']['path'].strip(".xz"),
    'url': '{}{}'.format(newmeta['baseURI'], newmeta['images']['azure']['path']),
}

with open(os.path.join(metadata_dir, 'fcos-amd64.json'), 'w') as f:
    json.dump(newmeta, f, sort_keys=True, indent=4)

# Continue to populate the legacy metadata file because there are still
# processes consuming this file directly. This normally could just be a symlink
# but some of these processes reference raw.githubusercontent.com which doesn't
# follow symlinks.
with open(os.path.join(metadata_dir, "fcos.json"), 'w') as f:
    json.dump(newmeta, f, sort_keys=True, indent=4)
