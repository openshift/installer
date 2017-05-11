#!/bin/bash
# filter hides lines with words over MAX chars

MAX=${MAX:-65}
sed -e "s/[a-zA-Z0-9\/+]\{${MAX},\}/***OMITTED***/g"

