#!/bin/bash
gfind ./sample.txt "pack" -po

: '
finds:
/examples/sample.txt
/examples/sample.txt
/examples/sample.txt
'
