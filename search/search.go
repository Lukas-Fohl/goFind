package finder

import (
	"os"
	"path"
	"strings"
	"sync"
	"unicode/utf8"
)

func FindExact(line *string, searchPattern string) (bool, []int) {
	if line == nil {
		return false, []int{}
	}
	splitLine := strings.Split(*line, "")
	splitPattern := strings.Split(searchPattern, "")
	if len(splitLine) == 0 || len(splitPattern) == 0 {
		return false, []int{}
	}

	//iterate over line and check for match at each char
	returnList := []int{}
	for i := 0; i < len(splitLine)-len(splitPattern)+1; i++ {
		searchLength := 0
		for j := 0; j < len(splitPattern); j++ {
			if splitLine[i+j] != splitPattern[j] {
				break
			} else {
				searchLength++
			}
		}

		if searchLength == len(splitPattern) {
			for j := 0; j < len(splitPattern); j++ {
				returnList = append(returnList, i+j)
			}
			return true, returnList
		}
	}

	return false, []int{}
}

func FindChars(line *string, searchPattern string) (bool, []int) {
	if line == nil {
		return false, []int{}
	}
	splitLine := strings.Split(*line, "")
	splitPattern := strings.Split(searchPattern, "")
	if len(splitLine) == 0 || len(splitPattern) == 0 {
		return false, []int{}
	}

	returnList := []int{}
	charsFound := 0
	//iterate over line and check if each char is matched
	for i := 0; i < len(splitLine); i++ {
		if charsFound < len(splitPattern) && splitLine[i] == splitPattern[charsFound] {
			charsFound++
			returnList = append(returnList, i)
		}
	}

	if len(splitPattern) == charsFound {
		return true, returnList
	}

	return false, []int{}
}

func FindFuzzy(line *string, searchPattern string) (bool, []int) {
	if line == nil {
		return false, []int{}
	}
	splitLine := strings.Split(*line, "")
	splitPattern := strings.Split(searchPattern, "")
	if len(splitLine) == 0 || len(splitPattern) == 0 {
		return false, []int{}
	}

	found, indices := FindExact(line, searchPattern)
	if found {
		return found, indices
	}

	//search with one char added somewhere
	for i := 0; i < len(splitPattern); i++ {
		found, indices := FindChars(line, searchPattern)
		if found && indices[len(indices)-1]-indices[0] <= len(splitPattern) {
			return found, indices
		}
	}

	//search pattern with each char missing -> one wrong char or one missing
	for i := 0; i < len(splitPattern); i++ {
		newSearch := strings.Join(splitPattern[:i], "") + strings.Join(splitPattern[i+1:], "")
		found, indices := FindChars(line, newSearch)
		if found && indices[len(indices)-1]-indices[0] < len(splitPattern) {
			return found, indices
		}
	}

	return false, []int{}
}

func FindTextInLine(line *string, settingsIn *Settings) (bool, []int) {
	tempLine := *line
	if !settingsIn.CheckCaseSensitive {
		tempLine = strings.ToLower(*line)
		settingsIn.SearchPattern = strings.ToLower(settingsIn.SearchPattern)
	}

	//check for right search
	if settingsIn.CheckNormal {
		found, index := FindExact(&tempLine, settingsIn.SearchPattern)
		return found, index
	}

	if settingsIn.CheckLetters {
		found, index := FindChars(&tempLine, settingsIn.SearchPattern)
		return found, index
	}

	if settingsIn.CheckFuzzy {
		found, index := FindFuzzy(&tempLine, settingsIn.SearchPattern)
		return found, index
	}

	return false, []int{}
}

func FindTextInBuff(buffIn *string, settingsIn Settings, c chan Location, wg *sync.WaitGroup) {
	defer wg.Done()

	if !utf8.ValidString(*buffIn) {
		return
	}

	fileLines := strings.Split(string(*buffIn), "\n") //get lines
	for i, lineIter := range fileLines {
		var found bool
		var index []int
		if !settingsIn.CheckCaseSensitive {
			settingsIn.SearchPattern = strings.ToLower(settingsIn.SearchPattern)
			temp := strings.ToLower(lineIter)
			found, index = FindTextInLine(&temp, &settingsIn)
		} else {
			found, index = FindTextInLine(&lineIter, &settingsIn)
		}
		if found {
			c <- Location{Path: "", Line: lineIter, LineNum: i, CharNum: index}
		}
	}

}

func FindTextInFile(pathIn string, SettingsIn Settings, c chan Location, wg *sync.WaitGroup) {
	defer wg.Done()

	if SettingsIn.CheckFileName {
		_, fileName := path.Split(pathIn)
		found, index := FindTextInLine(&fileName, &SettingsIn)
		if found {
			c <- Location{Path: pathIn, Line: fileName, LineNum: 0, CharNum: index}
		}
		return
	}

	dat, err := os.ReadFile(pathIn)
	if err != nil {
		panic(err)
	}

	if !utf8.ValidString(string(dat[len(dat)/5:])) {
		//fmt.Printf("%s is binary file\n", pathIn)
		return //check for binary-file
	}

	fileLines := strings.Split(string(dat), "\n") //get lines
	for i := 0; i < len(fileLines); i++ {
		found, index := FindTextInLine(&(fileLines[i]), &SettingsIn)
		if found {
			c <- Location{Path: pathIn, Line: fileLines[i], LineNum: i, CharNum: index} //add result to chanel
		}
	}
}
