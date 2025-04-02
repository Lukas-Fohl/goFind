package main

import finder "finder/search"

func main() {
	finder.Start()
}

/*
fuzzy requirements:
  -new feats:
    - file name [x]
    - fix fuzzy char index
    - fix idxs name to indices
    - get input from buffer (some | goFind "test")
    - no color
    - get amount in file
    - no info output

  - test:
    - for each flag
    - check flag function behav with sameples
    - for flag behav

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
