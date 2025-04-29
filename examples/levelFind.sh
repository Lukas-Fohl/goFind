#!/bin/bash
echo "gfind .. "import" -l 0"
gfind .. "import" -l 0

: '
finds:
/main.go:2,0:import finder "finder/search"
'

echo "gfind .. "import" -l 1"
gfind .. "import" -l 1

: '
finds:
/main.go:2,0:import finder "finder/search"
/tests/search_test.go:2,0:import (
/search/util.go:2,0:import (
/search/search.go:2,0:import (
/search/output.go:2,0:import (
/tests/flags_test.go:2,0:import (
/examples/sample.txt:948,0:important
/examples/sample.txt:1941,0:importance
/examples/sample.txt:2583,0:import
'
