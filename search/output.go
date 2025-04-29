package finder

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

func PrintHelp() {
	fmt.Println(`
## Usage
 - gfind [path] [pattern] [flags]
 - gfind [file] [pattern] [flags]
 - gfind [pattern] [flags] (assumes current path)
 - some_output | gfind [pattern] [flags]
 - e.g.:
   - gfind ./search "packge" -c
   - gfind ./main.go "start" -s -l 0
   - gfind "start" -s -l 0
   - cat main.go | gfind "package"
   - gfind ".go" -f
	`)
	fmt.Println(`
### Flags
  - "-l <number>":
    - level depth of file tree search
  - "-f":
    - check file name
  - "-i":
    - check if letters in locationIne
  - "-c":
    - input can have 1 letter changed (missing, added, different)
  - "-s":
    - not case sensitive
  - "-fl":
    - assumes input to be list of file-paths from stdin and searches those for a given patern (only works with piped input)
  - "-po":
    - prints only path of the result
  - "cf":
    - checks only for the first occurrence of an pattern in a file
  - "-n":
    - no info in output, just the locationIne
  - "-t":
    - removes color from output
  - "--help":
    - shows flags and usage
	`)
}

func printPath(pathIn string, buffer *bufio.Writer, color bool) {
	if color {
		(*buffer).Write([]byte("\x1b[1;36m" + pathIn + "\x1b[0m"))
	} else {
		(*buffer).Write([]byte(pathIn))
	}
}

func printLocation(locationIn Location, buffer *bufio.Writer, settingsIn Settings) {
	if settingsIn.CheckNormal {
		if !settingsIn.CheckFileName {
			buffer.Write([]byte(strconv.FormatInt(int64(locationIn.LineNum), 10) + ","))
		}
		buffer.Write([]byte(strconv.FormatInt(int64(locationIn.CharNum[0]), 10)))
	} else {
		buffer.Write([]byte(strconv.FormatInt(int64(locationIn.LineNum), 10)))
	}

}

func printLine(locationIn Location, buffer *bufio.Writer, settingsIn Settings) {
	coloredPrinted := 0
	splitLine := strings.Split(locationIn.Line, "")
	for i := 0; i < len(splitLine); i++ {
		if coloredPrinted < len(locationIn.CharNum) && i == locationIn.CharNum[coloredPrinted] && settingsIn.ShowColor {
			buffer.Write([]byte("\x1b[1;31m" + string(splitLine[i])))
			coloredPrinted++
		} else if !settingsIn.ShowColor {
			buffer.Write([]byte(string(splitLine[i])))
		} else {
			buffer.Write([]byte("\x1b[0m" + string(splitLine[i])))
		}
	}
}

func PrintResult(locationIn Location, instSettings Settings) {
	f := bufio.NewWriter(os.Stdout)
	defer f.Flush()

	if len(locationIn.CharNum) < 1 {
		return
	}

	absPath, err := filepath.Abs(locationIn.Path)
	if err != nil {
		fmt.Println(err)
		os.Exit(-1)
	}

	if (!instSettings.PipeInput || instSettings.ReadPipeFileList) && instSettings.ShowInfo {
		printPath(absPath, f, instSettings.ShowColor)

		if instSettings.ShowPathOnly {
			f.Write([]byte(string("\n")))
			return
		}

		f.Write([]byte(string(":")))

		printLocation(locationIn, f, instSettings)

		f.Write([]byte(":"))
	}

	if instSettings.ShowPathOnly {
		if instSettings.PipeInput && !instSettings.ReadPipeFileList {
			f.Write([]byte(string("Error: piped input has no path\n")))
		}
		return
	}

	printLine(locationIn, f, instSettings)

	if !instSettings.ShowColor {
		f.Write([]byte("\n"))
	} else {
		f.Write([]byte("\x1b[0m\n"))
	}
}
