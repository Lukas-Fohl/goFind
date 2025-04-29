#!/bin/bash
echo "gfind ../tests "package" -cf -po | gfind "func*Test" -fl"
gfind ../tests "package" -cf -po | gfind "func*Test" -fl

: '
finds:
/tests/search_test.go:27,0:func TestFindExact(t *testing.T) {
/tests/search_test.go:72,0:func TestFindChars(t *testing.T) {
/tests/search_test.go:109,0:func TestFindFuzzy(t *testing.T) {
/tests/search_test.go:181,0:func TestFindResticted(t *testing.T) {
/tests/search_test.go:331,0:func TestFindExactProp(t *testing.T) {
/tests/search_test.go:375,0:func TestFindFuzzyProp(t *testing.T) {
/tests/search_test.go:416,0:func TestCaseSearch(t *testing.T) {
/tests/flags_test.go:8,0:func TestPath(t *testing.T) {
/tests/flags_test.go:46,0:func TestFlag(t *testing.T) {
'

echo "gfind ../search "*.go" -f -po | gfind "package f" -fl"
gfind ../search "*.go" -f -po | gfind "package f" -fl

: '
finds:
/search/util.go:0,0:package finder
/search/output.go:0,0:package finder
/search/search.go
'
