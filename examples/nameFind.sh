#!/bin/bash
gfind .. ".go" -f

: '
finds:
/main.go:4:main.go
/search/output.go:6:output.go
/search/search.go:6:search.go
/tests/flags_test.go:10:flags_test.go
/tests/search_test.go:11:search_test.go
/search/util.go:4:util.go
:
