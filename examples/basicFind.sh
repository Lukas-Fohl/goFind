#!/bin/bash
echo "gfind ./sample.txt "tool""
gfind ./sample.txt "tool"

: '
finds:
/examples/sample.txt:1713,0:tool
/examples/sample.txt:3239,1:stool
'

echo "gfind ../tests "package test""
gfind ../tests "package test"

: '
finds:
/tests/search_test.go:0,0:package tests
/tests/flags_test.go:0,0:package tests
'
