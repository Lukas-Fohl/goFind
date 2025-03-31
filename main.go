package main

import finder "finder/search"

func main() {
	finder.MainCall()
}

/*
fuzzy requirements:
  - test:
    - for each flag
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
      - input can have 1 letter changed (missing, added, different)

TODO:
	[x] impl fuzzy
  impl file name search
	[x] split file into: main, search, output
	write docs
	[x] build test
	check for binary file
  other stuff
*/
