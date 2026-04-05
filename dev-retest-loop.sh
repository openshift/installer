#!/usr/bin/env bash

# 1089

while true ; do
  for pr in 6576; do
    echo "Pull Request $pr"
    ./dev-retest.sh $pr 
  done
  sleep 3600
done

