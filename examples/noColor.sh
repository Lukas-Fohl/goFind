#!/bin/bash
gfind ./sample.txt "hunt" -t

: '
finds:
/examples/sample.txt:3748,0:hunting
/examples/sample.txt:4757,0:hunt
'
