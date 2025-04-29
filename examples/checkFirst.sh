#!/bin/bash
gfind ./sample.txt "pack" -cf

: '
finds:
/examples/sample.txt:688,0:pack
'
