package finder

import (
	"bufio"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

func PrintResult(lin Location, instSettings Settings) {
	f := bufio.NewWriter(os.Stdout)
	defer f.Flush()

	charIndex := -1
	if len(lin.CharNum) > 0 {
		charIndex = lin.CharNum[0]
	} else {
		return
	}

	if !instSettings.ShowInfo {
		f.Write([]byte(lin.Line + "\n"))
		return
	}

	absPath, err := filepath.Abs(lin.Path)
	if err != nil {
		panic(err)
	}

	if !instSettings.PipeInput {
		if instSettings.ShowColor {
			f.Write([]byte("\x1b[1;36m" + absPath + "\x1b[0m:"))
		} else {
			f.Write([]byte(absPath))
		}

		if instSettings.CheckNormal {
			if !instSettings.CheckFileName {
				f.Write([]byte(strconv.FormatInt(int64(lin.LineNum), 10) + ","))
			}
			f.Write([]byte(strconv.FormatInt(int64(charIndex), 10)))
		} else {
			f.Write([]byte(strconv.FormatInt(int64(lin.LineNum), 10)))
		}
		f.Write([]byte(":"))
	}

	coloredPrinted := 0
	splitLine := strings.Split(lin.Line, "")
	for i := 0; i < len(splitLine); i++ {
		if coloredPrinted < len(lin.CharNum) && i == lin.CharNum[coloredPrinted] && instSettings.ShowColor {
			f.Write([]byte("\x1b[1;31m" + string(splitLine[i])))
			coloredPrinted++
		} else if !instSettings.ShowColor {
			f.Write([]byte(string(splitLine[i])))
		} else {
			f.Write([]byte("\x1b[0m" + string(splitLine[i])))
		}
	}

	if !instSettings.ShowColor {
		f.Write([]byte("\n"))
	} else {
		f.Write([]byte("\x1b[0m\n"))
	}
}
