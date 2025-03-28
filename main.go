package main

import finder "finder/search"

func main() {
	finder.MainCall()
}

/*
fuzzy requirements:
  - string as input
  - search files in child direcotries for lines with input in them
  - input not exact -> have same order of letters
  - sort out input for flags and text
  - how TO:
    - input:
      - "path" -> check if is path or file
    - get file tree
    - spawn thread for each file
    - search file:
      - look at each line -> look for first char -> for second in rest ...
      - if found send back in channel
    - return "nice" output for each found line (including file and number)

  - test (in python?):
    - for each flag
    - for an big file

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
	split file into: main, search, output
	write docs
	build test
	check for binary file
*/
