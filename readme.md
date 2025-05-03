> command-line tool for quickly finding text patterns in files, file-names and text
<br></br>
## Table of Contents
- [Table of Contents](#table-of-contents)
- [Requirements](#requirements)
- [Installation](#installation)
- [Usage](#usage)
- [Features](#features)
  - [Flags](#flags)
- [Examples](#examples)
- [Tests](#tests)
- [TODO](#todo)

## Requirements
 - git
 - make
 - go `>=1.22.2`

## Installation
 - **linux:**
```bash
git clone https://github.com/Lukas-Fohl/goFind
cd goFind && sudo make install
```
 - should work on MacOs if you have /usr/local/bin/

## Usage
 - `gfind [path] [pattern] [flags]`
 - `gfind [file] [pattern] [flags]`
 - `gfind [pattern] [flags]` (assumes current path)
 - `some_output | gfind [pattern] [flags]`
 - e.g.:
   - `gfind ./search "packge" -c`
   - `gfind ./main.go "start" -s -l 0`
   - `gfind "start" -s -l 0`
   - `cat main.go | gfind "package"`
   - `gfind ".go" -f`

## Features
 - read from path [[example]](./examples/basicFind.sh)
 - read from file name [[example]](./examples/basicFind.sh)
 - read from piped input [[example]](./examples/pipedFind.sh)
<br></br>
 - search recursive in dir for file content (with limit) [[example]](./examples/levelFind.sh)
 - search given file-list [[example]](./examples/fileList.sh)
<br></br>
 - check for exact input [[example]](./examples/basicFind.sh)
 - check for fuzzy input [[example]](./examples/fuzzyFind.sh)
 - check for letters in line [[example]](./examples/letterFind.sh)
 - check with/with out case sensitive [[example]](./examples/caseFind.sh)
 - check only for first occurrance [[example]](./examples/checkFirst.sh)
<br></br>
- output with and without color/info [[example]](./examples/noOutputFind.sh)
- output only path of pattern [[example]](./examples/pathOnly.sh)
<br></br>
- search with star-pattern [[example]](./examples/pattern.sh)
<br></br>
> check [Examples](#examples)

### Flags
  - `-l`:
    - level depth of file tree search
  - `-f`:
    - check file name
  - `-i`:
    - check if letters in line
  - `-c`:
    - input can have 1 letter changed (missing, added, different)
  - `-s`:
    - not case sensitive
  - `-fl`:
    - assumes input to be list of file-paths from stdin and searches those for a given patern (only works with piped input)
  - `-po`:
    - prints only path of the result
  - `cf`:
    - checks only for the first occurrence of an pattern in a file
  - `-n`:
    - no info in output, just the line
  - `-t`:
    - removes color from output
  - `--help`:
    - shows flags and usage
  - `--version`:
    - prints version


## Examples
 - initialize examples with `sh ./examples/init.sh`

## Tests
 - `go test ./tests -v`
 - for each search
 - check flag function behaviour with sameples

## TODO
 - [x] file name
 - [x] fix fuzzy char index
 - [x] fix indices name to indices
 - [x] get input from buffer (some | goFind "test")
 - [x] no info output
 - [x] no colors as flag!!!!!!!!!!!!!!
 - [x] case sensitive
 - [x] check file permissions
 - [x] replace panic with print when needed
 - [x] when looking at single file and is binary the tell that is binary
 - [x] oom error!!!!!!!!!!!!!
 - [x] readd gorutines!!!!!!!!!!!!!
 - [x] add flag test
 - [x] file list, check if len(string) > 0
 - [x] fix utf8 and unicode --> replace len(*line) with len(strings.Split(*line, ""))
 - [x] no tint with -po flag
 - [ ] get amount in file????????????? - no
 - [ ] rework flags????????????? - no
 - [ ] panic in debug????????????? - no
 - [x] star pattern?????????????
 - [x] eol pattern
 - [x] rewrite output!!!!!!!!!!!!!!
 - [x] error messages
 - [x] write --help, add to Error
 - [x] remove pointer to line
 - [x] double to lower?
 - [x] error when using -f in piped
 - [x] fuzzy bug
