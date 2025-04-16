package finder

import (
	"fmt"
	"io"
	"os"
	"path"
	"path/filepath"
	"strconv"
	"strings"
	"unicode/utf8"
)

type Settings struct {
	LevelRest          bool //level restriction
	LevelRestLimit     int  //value for ^
	CheckLetters       bool
	CheckFuzzy         bool
	CheckNormal        bool
	CheckFileName      bool
	CheckCaseSensitive bool
	ShowInfo           bool
	ShowColor          bool
	PipeInput          bool
	Path               string
	PathDepth          int
	SearchPattern      string
}

type Location struct {
	Line    string
	Path    string
	LineNum int
	CharNum []int
}

func DefaultSettings() Settings {
	return Settings{
		LevelRest:          false,
		LevelRestLimit:     -1,
		CheckLetters:       false,
		CheckFuzzy:         false,
		CheckNormal:        true,
		CheckFileName:      false,
		CheckCaseSensitive: true,
		ShowInfo:           true,
		ShowColor:          true,
		PipeInput:          false,
		PathDepth:          0,
		Path:               "",
		SearchPattern:      "",
	}
}

func FlagHandle(args []string) Settings {
	instSettings := DefaultSettings()

	if len(args) < 2 {
		fmt.Println("Error: not enougth arguments")
		os.Exit(-1)
	} else {
		//case no path is provided
		pathOut, err := os.Getwd()
		if err != nil {
			fmt.Println(err)
			os.Exit(-1)
		}

		instSettings.Path = pathOut
		instSettings.SearchPattern = args[1]
	}

	for i := 2; i < len(args) && len(args) > 2; i++ {
		switch args[i] {
		case "-i":
			instSettings.CheckNormal = false
			instSettings.CheckLetters = true
		case "-c":
			instSettings.CheckNormal = false
			instSettings.CheckFuzzy = true
		case "-f":
			instSettings.CheckFileName = true
		case "-n":
			instSettings.ShowInfo = false
		case "-t":
			instSettings.ShowColor = false
		case "-s":
			instSettings.CheckCaseSensitive = false
		case "-l":
			instSettings.LevelRest = true
			if i < len(args)-1 {
				argToInt, err := strconv.Atoi(args[i+1])
				if err != nil {
					fmt.Println("Error: no size provided for depth")
					os.Exit(-1)
				}
				instSettings.LevelRestLimit = argToInt
				i++
			} else {
				fmt.Println("Error: no size provided for depth")
				os.Exit(-1)
			}
		default:
			if i == 2 {
				//first two must be path and pattern
				instSettings.Path = args[1]
				instSettings.SearchPattern = args[2]

				fi, err := os.Stdin.Stat()
				if err != nil {
					panic(err)
				}

				if fi.Mode()&os.ModeNamedPipe != 0 {
					fmt.Println("Error: path in piped input")
					os.Exit(-1)
				}
			} else {
				fmt.Println("Error: flag not found")
				os.Exit(-1)
			}
		}
	}

	absPath, err := filepath.Abs(instSettings.Path)
	if err != nil {
		panic(err)
	}

	instSettings.Path = absPath
	instSettings.PathDepth = strings.Count(path.Join(instSettings.Path), string(os.PathSeparator))

	return instSettings
}

func Start() {
	instSettings := FlagHandle(os.Args)

	var pipe string
	fi, err := os.Stdin.Stat()
	if err != nil {
		panic(err)
	}
	if fi.Mode()&os.ModeNamedPipe != 0 {
		n, err := io.ReadAll(os.Stdin)
		if err != nil {
			panic(err)
		}
		if !utf8.ValidString(string(n)) {
			fmt.Println("Error: pipe input not valid utf8")
			os.Exit(-1)
		}
		instSettings.PipeInput = true
		pipe = string(n)
	}

	if instSettings.PipeInput {
		FindTextInBuff(&pipe, instSettings)
	} else {
		dat, err := os.Stat(instSettings.Path)
		if err != nil {
			panic(err)
		}

		switch pathType := dat.Mode(); {
		case pathType.IsDir():
			err := filepath.Walk(instSettings.Path,
				func(pathIn string, info os.FileInfo, err error) error {
					if err != nil {
						return err
					}

					stat, err := os.Stat(pathIn)
					if err == nil {
						if ((stat.Mode()&0111) == 0 || instSettings.CheckFileName) && !stat.IsDir() { //check if path is file and not executable
							currentPathDepth := strings.Count(path.Join(pathIn), string(os.PathSeparator)) - instSettings.PathDepth - 1
							if (instSettings.LevelRest && currentPathDepth <= instSettings.LevelRestLimit) || !instSettings.LevelRest {
								FindTextInFile(pathIn, instSettings)
							}
						}
					}
					return nil
				})

			if err != nil {
				panic(err)
			}

		case pathType.IsRegular():
			FindTextInFile(instSettings.Path, instSettings)
		}
	}

	//PrintResult(c, instSettings)
}
