#!/bin/bash
# filter hides lines with words over MAX chars

MAX=${MAX:-24}
sed -e "s/[^\ ]\{${MAX},\}/***OMITTED***/g"

