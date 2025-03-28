package finder

import (
	"os"
	"strings"
	"sync"
)

func FindExact(line *string, searchPattern string) (bool, []int) {
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
	returnList := []int{}
	charsFound := 0
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
	for i := 0; i < len(searchPattern); i++ {
		found, idxs := FindChars(line, searchPattern)
		if found && idxs[len(idxs)-1]-idxs[0] < len(searchPattern) {
			return found, idxs
		}
	}

	for i := 0; i < len(searchPattern); i++ {
		newSearch := searchPattern[:i] + searchPattern[i+1:]
		found, idxs := FindChars(line, newSearch)
		if found && idxs[len(idxs)-1]-idxs[0] < len(searchPattern) {
			return found, idxs
		}
	}
	return false, []int{}
}

func FindTextInLine(line *string, SettingsIn *Settings) (bool, []int) {
	if SettingsIn.checkNormal {
		found, index := FindExact(line, SettingsIn.searchPattern)
		return found, index
	}

	if SettingsIn.checkLetters {
		found, index := FindChars(line, SettingsIn.searchPattern)
		return found, index
	}

	if SettingsIn.checkFuzzy {
		found, index := FindFuzzy(line, SettingsIn.searchPattern)
		return found, index
	}

	return false, []int{}
}

func FindTextInFile(pathIn string, SettingsIn Settings, c chan Loaction, wg *sync.WaitGroup) {
	defer wg.Done()
	//fmt.Println(pathIn)
	dat, err := os.ReadFile(pathIn)
	if err != nil {
		panic(err)
	}
	fileLines := strings.Split(string(dat), "\n")
	for i := 0; i < len(fileLines); i++ {
		found, index := FindTextInLine(&(fileLines[i]), &SettingsIn)
		if found {
			c <- Loaction{path: pathIn, line: fileLines[i], lineNum: i, charNum: index}
		}
	}
}
