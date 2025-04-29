#!/bin/bash
gfind ./sample.txt "aeate" -i

: '
finds:
/fuzzy/examples/sample.txt:570:alternate
/fuzzy/examples/sample.txt:833:aggregate
/examples/sample.txt:2582:adequate
/examples/sample.txt:2686:alienate
'
