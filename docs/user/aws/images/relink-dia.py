#!/usr/bin/env python3
#
# Replace file:///... references with relative links and add the targets to Git.

import base64
import os
import sys
import urllib.parse
import xml.etree.ElementTree


def slug(uri):
    return os.path.basename(uri)


def get_def(path):
    tree = xml.etree.ElementTree.parse(path)
    root = tree.getroot()
    root.set('id', slug(path))
    return root


def relink(path):
    images = set()
    tree = xml.etree.ElementTree.parse(path)
    root = tree.getroot()
    parents = {c: p for p in tree.getiterator() for c in p}
    defs = xml.etree.ElementTree.Element('{http://www.w3.org/2000/svg}defs')
    root.insert(0, defs)
    for image in list(root.findall('{http://www.w3.org/2000/svg}image')):
        uri = urllib.parse.unquote(image.get('{http://www.w3.org/1999/xlink}href')).split('/./', 1)[-1]
        if uri not in images:
            defs.append(get_def(path=uri))
            images.add(uri)
        parent = parents[image]
        parent_index = list(parent).index(image)
        parent.remove(image)
        attrib = dict(image.items())
        attrib['{http://www.w3.org/1999/xlink}href'] = '#' + slug(uri)
        parent.insert(
            parent_index,
            xml.etree.ElementTree.Element('{http://www.w3.org/2000/svg}use', attrib=attrib))
    tree.write(path)


for path in sys.argv[1:]:
    relink(path=path)
