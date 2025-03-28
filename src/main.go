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
	checkFileName  bool
	path           string
	pathDepth      int
	searchPattern  string
}

type location struct {
	line    string
	path    string
	lineNum int
	charNum []int
}

func defaultSettings() settings {
	return settings{
		levelRest:      false,
		levelRestLimit: -1,
		checkLetters:   false,
		checkFuzzy:     false,
		checkNormal:    true,
		checkFileName:  false,
		pathDepth:      0,
		path:           "",
		searchPattern:  "",
	}
}

func flagHandle(args []string) settings {

	instSettings := defaultSettings()

	if len(args) < 2 {
		panic("not enougth arguments")
	} else {
		pathOut, err := os.Getwd()
		if err != nil {
			panic(err)
		}
		instSettings.path = pathOut
		instSettings.searchPattern = args[1]
	}

	for i := 2; i < len(args) && len(args) > 2; i++ {
		switch args[i] {
		case "-i":
			instSettings.checkNormal = false
			instSettings.checkLetters = true
		case "-c":
			instSettings.checkNormal = false
			instSettings.checkFuzzy = true
		case "-f":
			instSettings.checkFileName = true
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
			if i == 2 {
				instSettings.path = args[1]
				instSettings.searchPattern = args[2]
			} else {
				panic("flag not found")
			}
		}
	}

	instSettings.pathDepth = strings.Count(path.Join(instSettings.path), string(os.PathSeparator))

	return instSettings
}

func printResult(c chan location, instSettings settings) {
	for msg := range c {
		charIndex := -1
		if len(msg.charNum) > 0 {
			charIndex = msg.charNum[0]
		} else {
			break
		}
		newPath, err := filepath.Abs(msg.path)
		if err != nil {
			panic(err)
		}
		fmt.Print("\x1b[1;36m" + newPath + "\x1b[0m:")
		if instSettings.checkNormal {
			fmt.Print(strconv.FormatInt(int64(msg.lineNum), 10) + "," + strconv.FormatInt(int64(charIndex), 10))
		} else {
			fmt.Print(strconv.FormatInt(int64(msg.lineNum), 10))
		}
		fmt.Print(":")
		coloredPrinted := 0
		for i := 0; i < len(msg.line); i++ {
			if coloredPrinted < len(msg.charNum) && i == msg.charNum[coloredPrinted] {
				fmt.Print("\x1b[1;31m" + string(msg.line[i]))
				coloredPrinted++
			} else {
				fmt.Print("\x1b[0m" + string(msg.line[i]))
			}
		}
		fmt.Print("\x1b[0m\n")
	}
}

func main() {
	instSettings := flagHandle(os.Args)

	dat, err := os.Stat(instSettings.path)
	if err != nil {
		panic(err)
	}

	c := make(chan location)
	var wg sync.WaitGroup

	switch pathType := dat.Mode(); {
	case pathType.IsDir():
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
							go FindTextInFile(pathIn, instSettings, c, &wg)
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
		go FindTextInFile(instSettings.path, instSettings, c, &wg)
	}

	go func() {
		wg.Wait()
		close(c)
	}()

	printResult(c, instSettings)
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
