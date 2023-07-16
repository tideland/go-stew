#!/bin/sh
# tidyall.sh - Run go mod tidy on all subdirectories
# Usage: tidyall.sh
for D in */; do
    echo "$D"
	cd $D
	go get -u
	go mod tidy
	cd ..
done
#
## EOF
#
