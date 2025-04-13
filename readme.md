> command-line tool for quickly finding text patterns in files, file-names and text
<br></br>
## Table of Contents
- [Table of Contents](#table-of-contents)
- [Requirements:](#requirements)
- [Installation:](#installation)
- [Usage:](#usage)
- [Features:](#features)
  - [Flags:](#flags)
- [Tests:](#tests)
- [TODO:](#todo)

## Requirements:
 - git
 - make
 - go `>=1.22.2`

## Installation:
 - linux:
```bash
git clone github.com/lukas-fohl/goFind
cd goFind && sudo make install
```

## Usage:
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

## Features:
 - read from file
 - read from file name
 - read from piped input
<br></br>
 - search recursive in dir for file content (with limit)
<br></br>
 - check for exact input
 - check for fuzzy input
 - check for letters in line
 - check with/with out case sensitive
<br></br>
- output with and without color/info


### Flags:
  - `-l`:
    - level depth of file tree search
  - `-f`:
    - check file name
  - `-i`:
    - check if letters in line
  - `-c`:
    - input can have 1 letter changed (missing, added, different)
  - `-n`:
    - no info in output, just the line
  - `-t`:
    - removes color from output
  - `-s`:
    - not case sensitive

## Tests:
 - `go test ./tests -v`
 - for each search
 - check flag function behaviour with sameples

## TODO:
 - [x] impl fuzzy
 - [x] impl file name search
 - [x] split file into: main, search, output
 - [ ] write docs
 - [x] build test
 - [x] check for binary file
 - [ ] other stuff
 - [ ] change flags???
 <br></br>
 - [x] file name 
 - [x] fix fuzzy char index
 - [x] fix indices name to indices
 - [x] get input from buffer (some | goFind "test")
 - [ ] get amount in file
 - [x] no info output
 - [ ] rewrite output
 - [x] no colors as flag!!!!!!!!!!!!!!
 - [x] case sensitive
 - [ ] star pattern!!!!!!!!!!!!!
   - pattern-transformer:
     - \ as escape for * (when using as normal)
     - * search for exact (between *)
 - [ ] check file permissions
 - [ ] replace panic with print when needed
 - [ ] when looking at single file and is binary the tell that is binary
 - [ ] oom error
 - [x] fix utf8 and unicode --> replace len(*line) with len(strings.Split(*line, ""))
