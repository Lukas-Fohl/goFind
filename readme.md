### requirements
 - git
 - make
 - go `>=1.22.2`

### install
 - `git clone github.com/lukas-fohl/goFind`
 - `cd goFind && sudo make install`

### usage
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

### test:
 - `go test ./tests -v`
 - for each flag
 - check flag function behav with sameples

### flags:
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
  - `-s`:
    - not case sensitive

### features:
  - read from file
  - read from file name
  - read from piped input
</br>
</br>
  - search recursive in dir for file content (with limit)
</br>
</br>
  - check for exact input
  - check for fuzzy input
  - check for letters in line
  - check with/with out case sensitive
</br>
</br>
  - output with and without color/info

### TODO:
 - [x] impl fuzzy
 - [x] impl file name search
 - [x] split file into: main, search, output
 - [ ] write docs
 - [x] build test
 - [x] check for binary file
 - [ ] other stuff
 - [ ] change flags???

 - [x] file name 
 - [x] fix fuzzy char index
 - [x] fix indices name to indices
 - [x] get input from buffer (some | goFind "test")
 - [ ] get amount in file
 - [x] no info output
 - [ ] rewrite output
 - [ ] no colors as flag
 - [x] case sensitive
 - [ ] star pattern??????????????
