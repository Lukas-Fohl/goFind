package finder

import (
	"fmt"
	"path/filepath"
	"strconv"
	"strings"
)

func PrintResult(c chan Location, instSettings Settings) {
	for msg := range c {
		if !instSettings.ShowInfo {
			fmt.Println(msg.Line)
			continue
		}

		charIndex := -1
		if len(msg.CharNum) > 0 {
			charIndex = msg.CharNum[0]
		} else {
			break
		}

		absPath, err := filepath.Abs(msg.Path)
		if err != nil {
			panic(err)
		}

		if !instSettings.PipeInput {
			if instSettings.ShowColor {
				fmt.Print("\x1b[1;36m" + absPath + "\x1b[0m:")
			} else {
				fmt.Print(absPath)
			}

			if instSettings.CheckNormal {
				if !instSettings.CheckFileName {
					fmt.Print(strconv.FormatInt(int64(msg.LineNum), 10) + ",")
				}
				fmt.Print(strconv.FormatInt(int64(charIndex), 10))
			} else {
				fmt.Print(strconv.FormatInt(int64(msg.LineNum), 10))
			}
			fmt.Print(":")
		}

		coloredPrinted := 0
		splitLine := strings.Split(msg.Line, "")
		for i := 0; i < len(splitLine); i++ {
			if coloredPrinted < len(msg.CharNum) && i == msg.CharNum[coloredPrinted] && instSettings.ShowColor {
				fmt.Print("\x1b[1;31m" + string(splitLine[i]))
				coloredPrinted++
			} else if !instSettings.ShowColor {
				fmt.Print(string(splitLine[i]))
			} else {
				fmt.Print("\x1b[0m" + string(splitLine[i]))
			}
		}
		if !instSettings.ShowColor {
			fmt.Print("\n")
		} else {
			fmt.Print("\x1b[0m\n")
		}
	}
}
