package main

import (
	"fmt"
	"os"
	"path"
	"path/filepath"
	"strconv"
	"strings"
	"sync"
)

type settings struct {
	levelRest      bool //level restriction
	levelRestLimit int  //value for ^
	checkLetters   bool
	checkFuzzy     bool
	checkNormal    bool
	path           string
	pathDepth      int
	searchPattern  string
}

type location struct {
	line    string
	path    string
	lineNum int
	charNum int
}

func defaultSettings() settings {
	return settings{
		levelRest:      false,
		levelRestLimit: -1,
		checkLetters:   false,
		checkFuzzy:     false,
		checkNormal:    true,
		pathDepth:      0,
		path:           "",
		searchPattern:  "",
	}
}

func main() {
	args := os.Args

	instSettings := defaultSettings()

	if len(args) < 2 {
		panic("not enougth arguments")
	} else if len(args) == 2 {
		instSettings.searchPattern = args[1]
		pathOut, err := os.Getwd()
		if err != nil {
			panic(err)
		}
		instSettings.path = pathOut
	} else {
		instSettings.path = args[1]
		instSettings.searchPattern = args[2]
	}

	for i := 3; i < len(args) && len(args) > 3; i++ {
		switch args[i] {
		case "-i":
			instSettings.checkNormal = false
			instSettings.checkLetters = true
		case "-c":
			instSettings.checkNormal = false
			instSettings.checkFuzzy = true
		case "-l":
			instSettings.levelRest = true
			if i < len(args)-1 {
				argToInt, err := strconv.Atoi(args[i+1])
				if err != nil {
					panic("no size provided for depth")
				}
				instSettings.levelRestLimit = argToInt
				i++
			} else {
				panic("no size provided for depth")
			}
		default:
			panic("flag not found")
		}
	}

	instSettings.pathDepth = strings.Count(path.Join(instSettings.path), string(os.PathSeparator))
	//fmt.Println(instSettings)

	dat, err := os.Stat(instSettings.path)
	if err != nil {
		panic(err)
	}

	c := make(chan location)
	var wg sync.WaitGroup

	switch pathType := dat.Mode(); {
	case pathType.IsDir():
		//dir branch
		err := filepath.Walk(instSettings.path,
			func(pathIn string, info os.FileInfo, err error) error {
				if err != nil {
					return err
				}
				stat, err := os.Stat(pathIn)
				if err == nil {
					if (stat.Mode()&0111) == 0 && !stat.IsDir() {
						currentPathDepth := strings.Count(path.Join(pathIn), string(os.PathSeparator)) - instSettings.pathDepth - 1
						if (instSettings.levelRest && currentPathDepth <= instSettings.levelRestLimit) || !instSettings.levelRest {
							wg.Add(1)
							go findTextInFile(pathIn, instSettings, c, &wg)
						}
					}
				}
				return nil
			})
		if err != nil {
			panic(err)
		}

	case pathType.IsRegular():
		wg.Add(1)
		go findTextInFile(instSettings.path, instSettings, c, &wg)
	}
	go func() {
		wg.Wait()
		close(c)
	}()
	for msg := range c {
		if instSettings.checkNormal {
			fmt.Print("\x1b[1;36m" + msg.path + "\x1b[0m " + strconv.FormatInt(int64(msg.lineNum), 10) + "," + strconv.FormatInt(int64(msg.charNum), 10) + ":")
			for i := 0; i < len(msg.line); i++ {
				if i >= msg.charNum && i < msg.charNum+len(instSettings.searchPattern) {
					fmt.Print("\x1b[1;31m" + string(msg.line[i]))
				} else {
					fmt.Print("\x1b[0m" + string(msg.line[i]))
				}
			}
			fmt.Print("\x1b[0m\n")
		} else {
			fmt.Println(msg.path + " " + strconv.FormatInt(int64(msg.lineNum), 10) + ":" + msg.line)
		}
	}
}

func findExact(line *string, searchPattern string) (bool, int) {
	for i := 0; i < len(*line)-len(searchPattern)+1; i++ {
		searchLength := 0
		for j := 0; j < len(searchPattern); j++ {
			if (*line)[i+j] != searchPattern[j] {
				break
			} else {
				searchLength++
			}
		}
		if searchLength == len(searchPattern) {
			return true, i
		}
	}
	return false, -1
}

func findTextInLine(line *string, settingsIn *settings) (bool, int) {
	if settingsIn.checkNormal {
		found, index := findExact(line, settingsIn.searchPattern)
		return found, index
	}
	if settingsIn.checkLetters {
		charsFound := 0
		for i := 0; i < len(*line); i++ {
			if charsFound < len(settingsIn.searchPattern) && (*line)[i] == settingsIn.searchPattern[charsFound] {
				charsFound++
			}
		}
		if len(settingsIn.searchPattern) == charsFound {
			//fmt.Println(*line)
			return true, -1
		}
	}
	if settingsIn.checkFuzzy {
		return false, -1
	}

	return false, -1
}

func findTextInFile(pathIn string, settingsIn settings, c chan location, wg *sync.WaitGroup) {
	defer wg.Done()
	//fmt.Println(pathIn)
	dat, err := os.ReadFile(pathIn)
	if err != nil {
		panic(err)
	}
	fileLines := strings.Split(string(dat), "\n")
	for i := 0; i < len(fileLines); i++ {
		found, index := findTextInLine(&(fileLines[i]), &settingsIn)
		if found {
			c <- location{path: pathIn, line: fileLines[i], lineNum: i, charNum: index}
		}
	}
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
    - "-i":
      - check if letters in line [x]
    - "-c":
      - input can have 1 letter changed (missing, added, different)

TODO:
	impl fuzzy
	split file into: main, search, output
	write docs
	build test

*/
