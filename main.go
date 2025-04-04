package main

import finder "finder/search"

func main() {
	finder.Start()
}

/*
fuzzy requirements:
  -new feats:
    - file name [x]
    - fix fuzzy char index [x]
    - fix indices name to indices [x]
    - get input from buffer (some | goFind "test") [x]
    - get amount in file
    - no info output [x]
    - rewrite output
    - case sensitive [x]
    - fix lower case output when search
    - star pattern??????????????

  - test:
    - for each flag
    - check flag function behav with sameples

  - write docs
    - main readme with examples
      - installation
    - comments for functions
    - examples on how to use

  - flags:
    - "-l":
      - level depth of file tree search [x]
    - "-f"
      - check file name
    - "-i":
      - check if letters in line [x]
    - "-c":
      - input can have 1 letter changed (missing, added, different) [x]
    - "-n":
      - no info in output, just the line
    - "-s"
      - not case sensitive

  - feats:
    - read from file
    - read from file name
    - read from piped input

    - search recursive in dir for file content (with limit)

    - check for exact input
    - check for fuzzy input
    - check for letters in line
    - check with/with out case sensitive

    - output with and without color/info

TODO:
  [x] impl fuzzy
  [x]impl file name search
  [x] split file into: main, search, output
  write docs
  [x] build test
  [x] check for binary file
  other stuff
  change flags???
*/
