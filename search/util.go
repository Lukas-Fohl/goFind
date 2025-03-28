package finder

import (
	"os"
	"path"
	"path/filepath"
	"strconv"
	"strings"
	"sync"
)

type Settings struct {
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

type Loaction struct {
	line    string
	path    string
	lineNum int
	charNum []int
}

func DefaultSettings() Settings {
	return Settings{
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

func FlagHandle(args []string) Settings {

	instSettings := DefaultSettings()

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

func MainCall() {
	instSettings := FlagHandle(os.Args)

	dat, err := os.Stat(instSettings.path)
	if err != nil {
		panic(err)
	}

	c := make(chan Loaction)
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

	PrintResult(c, instSettings)
}
