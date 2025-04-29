#!/bin/bash
echo "gfind ./sample.txt "to*ch~""
gfind ./sample.txt "to*ch~"

: '
finds:
/examples/sample.txt:139,0:touch
/examples/sample.txt:331,1:stomach
'

echo "gfind ./sample.txt "some~""
gfind ./sample.txt "some~"

: '
finds: (find something that ends with "some")
/examples/sample.txt:332,0:some
/examples/sample.txt:796,4:handsome
'

echo "gfind ./sample.txt "un*ion""
gfind ./sample.txt "un*ion"

: '
finds: (find something that "un" then something and then "ion")
/examples/sample.txt:132,1:function
/examples/sample.txt:464,0:union
/examples/sample.txt:1233,2:foundation
/examples/sample.txt:1675,4:communion
/examples/sample.txt:2493,1:junction
/examples/sample.txt:3698,4:communication
/examples/sample.txt:4578,4:conjunction
'

echo "gfind ./sample.txt "un\*ion\~""
gfind ./sample.txt "un\*ion\~"

: '
finds: (something with "un*ion~")
'
