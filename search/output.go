package finder

import (
	"fmt"
	"path/filepath"
	"strconv"
)

func PrintResult(c chan Location, instSettings Settings) {
	for msg := range c {
		if !instSettings.ShowInfo {
			fmt.Println(msg.line)
			continue
		}
		charIndex := -1
		if len(msg.charNum) > 0 {
			charIndex = msg.charNum[0]
		} else {
			break
		}

		absPath, err := filepath.Abs(msg.path)
		if err != nil {
			panic(err)
		}

		fmt.Print("\x1b[1;36m" + absPath + "\x1b[0m:")
		if instSettings.CheckNormal {
			if !instSettings.CheckFileName {
				fmt.Print(strconv.FormatInt(int64(msg.lineNum), 10) + ",")
			}
			fmt.Print(strconv.FormatInt(int64(charIndex), 10))
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
