#!/usr/bin/python3
# Usage: ./hack/update-rhcos-bootimage.py https://releases-rhcos.svc.ci.openshift.org/storage/releases/ootpa/410.8.20190401.0/meta.json
import codecs,os,sys,json,argparse
import urllib.parse
import urllib.request

dn = os.path.abspath(os.path.dirname(sys.argv[0]))

parser = argparse.ArgumentParser()
parser.add_argument("meta", action='store')
args = parser.parse_args()

with urllib.request.urlopen(args.meta) as f:
    string_f = codecs.getreader('utf-8')(f)  # support for Python < 3.6
    meta = json.load(string_f)
newmeta = {}
for k in ['images', 'buildid', 'oscontainer',
          'ostree-commit', 'ostree-version',
          'azure', 'gcp']:
    newmeta[k] = meta[k]
newmeta['amis'] = {
    entry['name']: {
        'hvm': entry['hvm'],
    }
    for entry in meta['amis']
}
newmeta['baseURI'] = urllib.parse.urljoin(args.meta, '.')
with open(os.path.join(dn, "../data/data/rhcos.json"), 'w') as f:
    json.dump(newmeta, f, sort_keys=True, indent=4)
