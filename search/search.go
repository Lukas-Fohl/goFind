package finder

import (
	"fmt"
	"os"
	"path"
	"strings"
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

	indexList := []int{}
	charsFound := 0
	//iterate over line and check if each char is matched
	for i := 0; i < len(splitLine); i++ {
		if charsFound < len(splitPattern) && splitLine[i] == splitPattern[charsFound] {
			charsFound++
			indexList = append(indexList, i)
		}
	}

	if len(splitPattern) == charsFound {
		return true, indexList
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
	found, indices = FindChars(line, searchPattern)
	if found && (indices[len(indices)-1]-indices[0]) <= len(splitPattern) {
		return found, indices
	}

	//search pattern with each char missing -> one wrong char or one missing
	for i := 0; i < len(splitPattern); i++ {
		newSearch := strings.Join(splitPattern[:i], "") + strings.Join(splitPattern[i+1:], "")
		found, indices := FindChars(line, newSearch)
		if found && indices[len(indices)-1]-indices[0] < len(splitPattern) {
			return found, indices
		} else {
			found, indices = FindExact(line, newSearch)
			if found {
				return found, indices
			}
		}
	}

	return false, []int{}
}

func FindRestriced(line *string, searchPattern string) (bool, []int) {
	starSplit := strings.Split(searchPattern, "\\*")
	endSplit := strings.Split(searchPattern, "\\~")
	if len(endSplit) > 2 {
		fmt.Println("Error: to many end-symbols")
		os.Exit(-1)
	} else if len(endSplit) > 1 && endSplit[1] != "" {
		fmt.Println("Error: end-symbol at wrong position")
		os.Exit(-1)
	}

	//REM \\~ from star search if included
	listOfFound := [][]int{}
	if len(starSplit) > 1 {
		for _, elem := range starSplit {
			if len(elem) > 0 {
				found, indices := FindExact(line, strings.ReplaceAll(elem, "\\~", ""))
				if found {
					listOfFound = append(listOfFound, indices)
				} else {
					return false, []int{}
				}
			}
		}

		lastLast := -1
		lastFirst := -1
		for _, elem := range listOfFound {
			if elem[0] > lastFirst && elem[len(elem)-1] > lastLast && lastLast < elem[0] {
				lastFirst = elem[0]
				lastLast = elem[len(elem)-1]
			} else {
				return false, []int{}
			}
		}
	}

	if len(starSplit) > 1 {
		if len(endSplit) > 1 {
			lofLast := listOfFound[len(listOfFound)-1][len(listOfFound[len(listOfFound)-1])-1]
			if lofLast == len(*line)-1 {
				returnList := []int{}
				for _, i := range listOfFound {
					returnList = append(returnList, i...)
				}
				return true, returnList
			} else {
				return false, []int{}
			}
		} else {
			returnList := []int{}
			for _, i := range listOfFound {
				returnList = append(returnList, i...)
			}
			return true, returnList
		}
	} else {
		found, indices := FindExact(line, endSplit[0])
		if found && indices[len(indices)-1] == len(*line)-1 {
			return found, indices
		} else {
			return false, []int{}
		}
	}

	/*ops:
	- no split no end
	  - ¯\_(ツ)_/¯
	- no split yes end
	  - check for exact -> last elem = len-1
	- yes split no end
	  - check exact each elem -> no overlap + following -> return union
	- yes split yes end
	  - do check for split no end check if last
	*/
}

func FindTextInLine(line *string, settingsIn *Settings) (bool, []int) {
	tempLine := *line
	if !settingsIn.CheckCaseSensitive {
		tempLine = strings.ToLower(*line)
		settingsIn.SearchPattern = strings.ToLower(settingsIn.SearchPattern)
	}

	//call restriced search function
	starSplit := strings.Split(settingsIn.SearchPattern, "\\*")
	endSplit := strings.Split(settingsIn.SearchPattern, "\\~")
	if len(starSplit) > 1 || len(endSplit) > 1 {
		return FindRestriced(&tempLine, settingsIn.SearchPattern)
	}

	//check for right search
	if settingsIn.CheckNormal {
		return FindExact(&tempLine, settingsIn.SearchPattern)
	}

	if settingsIn.CheckLetters {
		return FindChars(&tempLine, settingsIn.SearchPattern)
	}

	if settingsIn.CheckFuzzy {
		return FindFuzzy(&tempLine, settingsIn.SearchPattern)
	}

	return false, []int{}
}

func FindTextInBuff(buffIn *string, settingsIn Settings) []Location {
	locationList := []Location{}
	if !utf8.ValidString(*buffIn) {
		return locationList
	}

	fileLines := strings.Split(string(*buffIn), "\n") //get lines
	for i, lineIter := range fileLines {
		var found bool
		var index []int
		found, index = FindTextInLine(&lineIter, &settingsIn)
		if found {
			locationList = append(locationList, Location{Path: "", Line: lineIter, LineNum: i, CharNum: index})
			if settingsIn.CheckFirst {
				return locationList
			}
		}
	}

	return locationList
}

func FindTextInFile(pathIn string, SettingsIn Settings) []Location {
	locationList := []Location{}
	if SettingsIn.CheckFileName {
		_, fileName := path.Split(pathIn)
		found, index := FindTextInLine(&fileName, &SettingsIn)
		if found {
			locationList = append(locationList, Location{Path: pathIn, Line: fileName, LineNum: 0, CharNum: index})
			if SettingsIn.CheckFirst {
				return locationList
			}
		}
		return locationList
	}

	dat, err := os.ReadFile(pathIn)
	if err != nil {
		fmt.Println(err)
	}

	if !utf8.ValidString(string(dat[len(dat)/5:])) {
		if pathIn == SettingsIn.Path {
			fmt.Printf("%s is a binary file\n", pathIn)
		}
		return locationList //check for binary-file
	}

	fileLines := strings.Split(string(dat), "\n") //get lines
	for i := 0; i < len(fileLines); i++ {
		found, index := FindTextInLine(&(fileLines[i]), &SettingsIn)
		if found {
			locationList = append(locationList, Location{Path: pathIn, Line: fileLines[i], LineNum: i, CharNum: index})
			if SettingsIn.CheckFirst {
				return locationList
			}
		}
	}

	return locationList
}
