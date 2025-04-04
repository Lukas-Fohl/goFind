package finder

import (
	"os"
	"path"
	"strings"
	"sync"
	"unicode/utf8"
)

func FindExact(line *string, searchPattern string) (bool, []int) {
	if line == nil || len(*line) == 0 || len(searchPattern) == 0 {
		return false, []int{}
	}

	//iterate over line and check for match at each char
	returnList := []int{}
	for i := 0; i < len(*line)-len(searchPattern)+1; i++ {
		searchLength := 0
		for j := 0; j < len(searchPattern); j++ {
			if (*line)[i+j] != searchPattern[j] {
				break
			} else {
				searchLength++
			}
		}

		if searchLength == len(searchPattern) {
			for j := 0; j < len(searchPattern); j++ {
				returnList = append(returnList, i+j)
			}
			return true, returnList
		}
	}

	return false, []int{}
}

func FindChars(line *string, searchPattern string) (bool, []int) {
	if line == nil || len(*line) == 0 || len(searchPattern) == 0 {
		return false, []int{}
	}

	returnList := []int{}
	charsFound := 0
	//iterate over line and check if each char is matched
	for i := 0; i < len(*line); i++ {
		if charsFound < len(searchPattern) && (*line)[i] == searchPattern[charsFound] {
			charsFound++
			returnList = append(returnList, i)
		}
	}

	if len(searchPattern) == charsFound {
		return true, returnList
	}

	return false, []int{}
}

func FindFuzzy(line *string, searchPattern string) (bool, []int) {
	if line == nil || len(*line) == 0 || len(searchPattern) == 0 {
		return false, []int{}
	}

	found, indices := FindExact(line, searchPattern)
	if found {
		return found, indices
	}

	//search with one char added somewhere
	for i := 0; i < len(searchPattern); i++ {
		found, indices := FindChars(line, searchPattern)
		if found && indices[len(indices)-1]-indices[0] <= len(searchPattern) {
			return found, indices
		}
	}

	//search pattern with each char missing -> one wrong char or one missing
	for i := 0; i < len(searchPattern); i++ {
		newSearch := searchPattern[:i] + searchPattern[i+1:]
		found, indices := FindChars(line, newSearch)
		if found && indices[len(indices)-1]-indices[0] < len(searchPattern) {
			return found, indices
		}
	}

	return false, []int{}
}

func FindTextInLine(line *string, settingsIn *Settings) (bool, []int) {
	if !settingsIn.CheckCaseSensitive {
		*line = strings.ToLower(*line)
		settingsIn.SearchPattern = strings.ToLower(settingsIn.SearchPattern)
	}

	//check for right search
	if settingsIn.CheckNormal {
		found, index := FindExact(line, settingsIn.SearchPattern)
		return found, index
	}

	if settingsIn.CheckLetters {
		found, index := FindChars(line, settingsIn.SearchPattern)
		return found, index
	}

	if settingsIn.CheckFuzzy {
		found, index := FindFuzzy(line, settingsIn.SearchPattern)
		return found, index
	}

	return false, []int{}
}

func FindTextInBuff(buffIn *string, settingsIn Settings, c chan Location, wg *sync.WaitGroup) {
	defer wg.Done()

	if !utf8.ValidString(*buffIn) {
		return
	}

	if !settingsIn.CheckCaseSensitive {
		*buffIn = strings.ToLower(*buffIn)
		settingsIn.SearchPattern = strings.ToLower(settingsIn.SearchPattern)
	}

	fileLines := strings.Split(string(*buffIn), "\n") //get lines
	for i, lineIter := range fileLines {
		found, index := FindTextInLine(&lineIter, &settingsIn)
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
