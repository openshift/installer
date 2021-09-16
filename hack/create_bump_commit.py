#!/usr/bin/env python

'''
Create a commit for boot image bump for RHCOS
Usage: python create_bump_commit.py -v 4.6 -b 37124

Requires python3-GitPython

'''
import argparse
import git
import os
import subprocess
import sys

ARCH_LIST   = ['ppc64le', 's390x', 'amd64', 'aarch64','all']
URL         = 'https://rhcos-redirector.apps.art.xq1c.p1.openshiftapps.com/art/storage/releases/rhcos'
VERSIONS    = ('4.6', '4.7', '4.8', '4.9', '4.10')
COMMANDS    = []

OLD_ARCH_LIST   =  ['ppc64le', 's390x', 'amd64']
NEW_ARCH_LIST   =  OLD_ARCH_LIST.copy()
NEW_ARCH_LIST.append('aarch64')

ARCH_ALLOWED = {}
ARCH_ALLOWED['4.6']    = OLD_ARCH_LIST
ARCH_ALLOWED['4.7']    = OLD_ARCH_LIST
ARCH_ALLOWED['4.8']    = NEW_ARCH_LIST
ARCH_ALLOWED['4.9']    = NEW_ARCH_LIST
ARCH_ALLOWED['4.10']   = NEW_ARCH_LIST

def usage():
    '''
    # Handle parameters
    @param version string rhcos version
    '''

    parser = argparse.ArgumentParser(
        description='python boot-image-bumps.py -v 4.6',
    )

    parser.add_argument("-a", "--arch", choices=ARCH_LIST, default="all")
    parser.add_argument("-b", "--bugnumber", required=False)
    parser.add_argument("-v", "--version", choices=VERSIONS, required=True)

    args  = parser.parse_args()

    bugNumber = args.bugnumber
    version  = args.version

    if args.arch == "all":
        arches = ARCH_LIST
        arches.remove("all")
        if ("aarch64" in arches) and ("aarch64" not in ARCH_ALLOWED[version]):
            arches.remove("aarch64")
    else:
        arches = [args.arch]
    for arch in arches:
        if arch in ARCH_ALLOWED[version]:
            continue
        else:
            print("Arch %s not supported for version %s" % (arch, version))
            exit()

    return version, arches, bugNumber


def runCmd(cmd):
    '''
    Run the given command using subprocess.Popen and verify its return code.
    @param str command command to be executed
    '''

    try:
        process = subprocess.Popen(cmd,shell=True,stdout=subprocess.PIPE,stderr=subprocess.STDOUT)
        output = process.communicate()[0].strip().decode( "utf-8" ).replace('"', '')
    except subprocess.CalledProcessError as e:
        print('An exception h:as occurred: {0}'.format(e))
        sys.exit(1)
    return output

def getFullURL(arch, version):

     if (arch == "amd64" or  arch == "x86_64"):
         return '%s-%s' % (URL, version)
     return '%s-%s-%s' % (URL, version, arch)

def  getReleaseID(url):

     cmd = ('curl -Ls %s/builds.json | jq .builds[0].id') % (url)
     return runCmd(cmd)

def runScript(cmd):
    '''
    Run the given command using check_call and verify its return code.
    @param str command command to be executed
    '''
    try:
        subprocess.check_call(cmd.split())
    except subprocess.CalledProcessError as e:
        print('An exception h:as occurred: {0}'.format(e))
        sys.exit(1)
    COMMANDS.append(cmd)

def runUpdate(arch, version):

     url = getFullURL(arch, version)
     release = getReleaseID(url)
     print("Processing: %s - %s" % (arch, release))

     arch1 = arch
     if (arch == "amd64"):
         arch1 = "x86_64"
     cmd = 'python update-rhcos-bootimage.py %s/%s/%s/meta.json %s' % (url, release, arch1, arch)
     runScript(cmd)

def createCommit(version, bugNumber=None):
    ## Requires python3-GitPython
    repo = git.Repo("../")
    repo.git.add(update=True)
    if bugNumber:
        title = ('Bug %s: Bump boot images for RHCOS %s fixes' % (bugNumber,version))
    else:
        title = ('Bump boot images for RHCOS %s fixes' % (version))
    cmds = '\n'.join(map(str, COMMANDS))
    description = ('Changes generated with:')
    repo.git.commit('-m', title, '-m', description, '-m', cmds)
    print('Created: %s'% (repo.head.commit.message))

def main():

    version, arches, bugNumber = usage()

    for arch in arches:
        runUpdate(arch, version)
    #TODO - plume for 4.8+
    createCommit(version, bugNumber)

if __name__ == "__main__":
    main()

