#!/bin/sh
set -eu
find . -name '*.go' ! -name '*_test.go' -print0 | xargs -0 awk 'BEGIN{n=0}{n++}END{print "non-test Go lines:",n; if(n<2000) exit 1}'
