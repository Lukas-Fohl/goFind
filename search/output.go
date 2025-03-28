package finder

import (
	"fmt"
	"path/filepath"
	"strconv"
)

func PrintResult(c chan Loaction, instSettings Settings) {
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
