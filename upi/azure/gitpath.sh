#!/bin/bash

#get REPO top level dir
REPO=`git rev-parse --show-toplevel`

if [[ $? -ne 0 ]]
then
    echo "Not in git repo"
    exit 1
fi

git_branch=`git rev-parse --abbrev-ref HEAD 2>/dev/null`
#get relative path
relative_path=`git rev-parse --show-prefix|sed "s/\/$//"`
GIT_REMOTE_URL_UNFINISHED=`git config --get remote.origin.url|sed -s "s/^ssh/http/; s/git@//; s/.git$//;"`
GIT_REMOTE_URL="$(dirname $GIT_REMOTE_URL_UNFINISHED)/$(basename $GIT_REMOTE_URL_UNFINISHED)/$relative_path"

echo $GIT_REMOTE_URL

if [[ $git_branch != "master" ]]
then
    #if the feature branch name contains "/", transform it as per URL Encoding
    git_branch_in_url=`echo $git_branch| sed 's/\//%2F/g'`
    echo "Feature Branch URL"
    echo "${GIT_REMOTE_URL}/tree/$git_branch_in_url"
fi

