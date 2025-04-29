#!/bin/bash
gfind ./sample.txt "volunteer" -c
gfind ./sample.txt "volnteer" -c
gfind ./sample.txt "voluunteer" -c
gfind ./sample.txt "volonteer" -c

: '
all find:
/examples/sample.txt:4206:volunteer
'
