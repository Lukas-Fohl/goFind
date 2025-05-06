package finder

import (
	"fmt"
	"io"
	"os"
	"path"
	"path/filepath"
	"strconv"
	"strings"
	"sync"
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
	CheckFirst         bool
	ShowInfo           bool
	ShowColor          bool
	ShowPathOnly       bool
	PipeInput          bool
	ReadPipeFileList   bool
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
		CheckFirst:         false,
		ShowInfo:           true,
		ShowColor:          true,
		ShowPathOnly:       false,
		PipeInput:          false,
		ReadPipeFileList:   false,
		PathDepth:          0,
		Path:               "",
		SearchPattern:      "",
	}
}

func FlagHandle(args []string) Settings {
	flagSettings := DefaultSettings()

	if len(args) > 1 && args[1] == "--help" {
		PrintHelp()
		os.Exit(-1)
	} else if len(args) > 1 && args[1] == "--version" {
		PrintVersion()
		os.Exit(-1)
	}

	if len(args) < 2 {
		fmt.Println("Error: not enough arguments")
		os.Exit(-1)
	} else {
		//case no path is provided
		pathOut, err := os.Getwd()
		if err != nil {
			fmt.Println(err)
			os.Exit(-1)
		}

		flagSettings.Path = pathOut
		flagSettings.SearchPattern = args[1]
	}

	for i := 2; i < len(args) && len(args) > 2; i++ {
		switch args[i] {
		case "-i":
			flagSettings.CheckNormal = false
			flagSettings.CheckLetters = true
		case "-c":
			flagSettings.CheckNormal = false
			flagSettings.CheckFuzzy = true
		case "-f":
			flagSettings.CheckFileName = true
		case "-n":
			flagSettings.ShowInfo = false
			flagSettings.ShowColor = false
		case "-t":
			flagSettings.ShowColor = false
		case "-s":
			flagSettings.CheckCaseSensitive = false
		case "-fl":
			flagSettings.ReadPipeFileList = true
		case "-po":
			flagSettings.ShowPathOnly = true
			flagSettings.ShowColor = false
		case "-cf":
			flagSettings.CheckFirst = true
		case "--help":
			PrintHelp()
			os.Exit(-1)
		case "--version":
			PrintVersion()
			os.Exit(-1)
		case "-l":
			flagSettings.LevelRest = true
			if i < len(args)-1 {
				argToInt, err := strconv.Atoi(args[i+1])
				if err != nil {
					fmt.Println("Error: no size provided for depth")
					os.Exit(-1)
				}

				flagSettings.LevelRestLimit = argToInt
				i++
			} else {
				fmt.Println("Error: no size provided for depth")
				os.Exit(-1)
			}
		default:
			if i == 2 {
				//first two must be path and pattern
				flagSettings.Path = args[1]
				flagSettings.SearchPattern = args[2]

				fi, err := os.Stdin.Stat()
				if err != nil {
					fmt.Println(err)
					os.Exit(-1)
				}

				if len(args[i]) > 0 && args[i][0] == '-' {
					fmt.Println("Error: flag not found: " + args[i])
					os.Exit(-1)
				}

				if fi.Mode()&os.ModeNamedPipe != 0 {
					fmt.Println("Error: path in piped input")
					os.Exit(-1)
				}
			} else {
				fmt.Println("Error: flag not found: " + args[i])
				os.Exit(-1)
			}
		}
	}

	absPath, err := filepath.Abs(flagSettings.Path)
	if err != nil {
		fmt.Println(err)
		os.Exit(-1)
	}

	flagSettings.Path = absPath
	flagSettings.PathDepth = strings.Count(path.Join(flagSettings.Path), string(os.PathSeparator))

	return flagSettings
}

func Start() {
	flagSettings := FlagHandle(os.Args)

	var wg sync.WaitGroup

	var pipe string
	fi, err := os.Stdin.Stat()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	if fi.Mode()&os.ModeNamedPipe != 0 {
		if flagSettings.CheckFileName {
			fmt.Println("Error: -f in piped input")
			os.Exit(-1)
		}

		n, err := io.ReadAll(os.Stdin)
		if err != nil {
			fmt.Println(err)
			os.Exit(-1)
		}

		if !utf8.ValidString(string(n)) {
			fmt.Println("Error: pipe input not valid utf8")
			os.Exit(-1)
		}

		flagSettings.PipeInput = true
		pipe = string(n)
	}

	if flagSettings.PipeInput {
		wg.Add(1)
		go func() {
			defer wg.Done()
			if flagSettings.ReadPipeFileList {
				for _, lineIter := range strings.Split(pipe, "\n") {
					if len(lineIter) > 0 {
						for _, res := range FindTextInFile(lineIter, flagSettings) {
							PrintResult(res, flagSettings)
							if flagSettings.CheckFirst {
								break
							}
						}
					}
				}
			} else {
				for _, res := range FindTextInBuff(pipe, flagSettings) {
					PrintResult(res, flagSettings)
					if flagSettings.CheckFirst {
						break
					}
				}
			}
		}()
	} else {
		dat, err := os.Stat(flagSettings.Path)
		if err != nil {
			fmt.Println(err)
			os.Exit(-1)
		}

		switch pathType := dat.Mode(); {
		case pathType.IsDir():
			err := filepath.Walk(flagSettings.Path,
				func(pathIn string, info os.FileInfo, err error) error {
					if err != nil {
						return err
					}

					currentPathDepth := strings.Count(path.Join(pathIn), string(os.PathSeparator)) - flagSettings.PathDepth - 1
					stat, err := os.Stat(pathIn)
					if err == nil &&
						(((stat.Mode()&0111) == 0 || flagSettings.CheckFileName) && !stat.IsDir()) && //check if path is file and not executable
						((flagSettings.LevelRest && currentPathDepth <= flagSettings.LevelRestLimit) || !flagSettings.LevelRest) { //check path level
						wg.Add(1)
						go func() {
							defer wg.Done()
							for _, res := range FindTextInFile(pathIn, flagSettings) {
								PrintResult(res, flagSettings)
								if flagSettings.CheckFirst {
									break
								}
							}
						}()
					}

					return nil
				})

			if err != nil {
				fmt.Println(err)
				os.Exit(-1)
			}

		case pathType.IsRegular():
			wg.Add(1)
			go func() {
				defer wg.Done()
				for _, res := range FindTextInFile(flagSettings.Path, flagSettings) {
					PrintResult(res, flagSettings)
				}
			}()
		}
	}

	wg.Wait()
}
