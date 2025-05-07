package finder

import (
	"fmt"
	"os"
	"path"
	"strings"
	"unicode/utf8"
)

func FindExact(line string, searchPattern string) (bool, []int) {
	splitLine := strings.Split(line, "")
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

func FindChars(line string, searchPattern string) (bool, []int) {
	splitLine := strings.Split(line, "")
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

func FindFuzzy(line string, searchPattern string) (bool, []int) {
	splitLine := strings.Split(line, "")
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

func getPatternSplit(searchPattern string) ([]string, []string) {
	starSplit := []string{""}
	endSplit := []string{""}

	//split string wrt. special character
	splitSearch := strings.Split(searchPattern, "")
	for i := 0; i < len(splitSearch); i++ {
		switch splitSearch[i] {
		case "\\":
			if i < len(splitSearch)-1 {
				starSplit[len(starSplit)-1] += splitSearch[i+1]
				endSplit[len(endSplit)-1] += splitSearch[i+1]
				i++
			}
		case "*":
			starSplit = append(starSplit, "")
		case "~":
			endSplit = append(endSplit, "")
		default:
			starSplit[len(starSplit)-1] += splitSearch[i]
			endSplit[len(endSplit)-1] += splitSearch[i]
		}
	}
	return starSplit, endSplit
}

// check for findExact of parts split by '*'
func getStarPatternIndices(line string, starSplit []string) [][]int {
	listOfFound := [][]int{}

	remainingLine := line // Keep track of the remaining part of the line
	offset := 0           // Offset to adjust indices relative to the original string

	for _, elem := range starSplit {
		if len(elem) > 0 {
			found, indices := FindExact(remainingLine, elem)
			if found {
				// Adjust indices to be relative to the original string
				for i := range indices {
					indices[i] += offset
				}
				listOfFound = append(listOfFound, indices)

				// Update the remaining line and offset
				lastMatchEnd := indices[len(indices)-1] + 1
				remainingLine = remainingLine[lastMatchEnd-offset:]
				offset = lastMatchEnd
			} else {
				return [][]int{}
			}
		}
	}

	//check for order
	lastLast := -1
	lastFirst := -1
	for _, elem := range listOfFound {
		if elem[0] > lastFirst && elem[len(elem)-1] > lastLast && lastLast < elem[0] {
			lastFirst = elem[0]
			lastLast = elem[len(elem)-1]
		} else {
			return [][]int{}
		}
	}

	return listOfFound
}

// get true last position in end-check (search from right to left)
func getEndPatternIndices(line string, starSplit []string, endSplit []string) []int {
	elem := ""
	if len(endSplit) > 1 && len(starSplit) > 1 {
		elem = starSplit[len(starSplit)-1]
	} else {
		elem = endSplit[len(endSplit)-2]
	}

	searchLine := line
	searchOffset := 0
	lastOffset := 0
	found, indices := FindExact(searchLine, elem)
	for found {
		tempFound, tempIndices := FindExact(searchLine[searchOffset:], elem)
		if tempFound {
			indices = tempIndices
			searchOffset += tempIndices[len(tempIndices)-1] + 1
			lastOffset = tempIndices[len(tempIndices)-1] + 1
		} else {
			break
		}
	}

	searchOffset -= lastOffset
	if !found || len(indices) == 0 {
		return []int{}
	} else {
		for i := 0; i < len(indices); i++ {
			indices[i] += searchOffset
		}

		return indices
	}
}

func FindRestriced(line string, searchPattern string) (bool, []int) {
	starSplit, endSplit := getPatternSplit(searchPattern)

	if len(endSplit) > 2 {
		fmt.Println("Error: to many end-symbols")
		os.Exit(-1)
	} else if len(endSplit) > 1 && endSplit[1] != "" {
		fmt.Println("Error: end-symbol at wrong position")
		os.Exit(-1)
	}

	listOfFound := [][]int{}
	if len(starSplit) > 1 {
		listOfFound = getStarPatternIndices(line, starSplit)
		if len(listOfFound) == 0 {
			return false, []int{}
		}
	}

	endPatternIndices := []int{}
	if (len(endSplit) > 1 && len(starSplit) > 1) || len(endSplit) > 1 {
		endPatternIndices = getEndPatternIndices(line, starSplit, endSplit)
		if len(endPatternIndices) == 0 {
			return false, []int{}
		}
	}

	//return logic
	if len(starSplit) > 1 {
		if len(endSplit) > 1 {
			if endPatternIndices[len(endPatternIndices)-1] == len(line)-1 {
				//concat list
				returnList := []int{}
				for _, i := range listOfFound[:len(listOfFound)-1] { //all but last element
					returnList = append(returnList, i...)
				}
				returnList = append(returnList, endPatternIndices...)

				return true, returnList
			} else {
				return false, []int{}
			}
		} else {
			//concat list
			returnList := []int{}
			for _, i := range listOfFound {
				returnList = append(returnList, i...)
			}

			return true, returnList
		}
	} else {
		found, indices := FindExact(line, endSplit[0])
		if found && len(endPatternIndices) > 0 && endPatternIndices[len(endPatternIndices)-1] == len(line)-1 {
			return found, endPatternIndices
		} else if found && len(indices) > 0 && indices[len(indices)-1] == len(line)-1 {
			return found, indices
		} else {
			return false, []int{}
		}
	}
}

func FindTextInLine(line string, settingsIn Settings) (bool, []int) {
	tempLine := line
	if !settingsIn.CheckCaseSensitive {
		tempLine = strings.ToLower(line)
		settingsIn.SearchPattern = strings.ToLower(settingsIn.SearchPattern)
	}

	//call restriced search function
	lenStarSplit := len(strings.Split(settingsIn.SearchPattern, "*"))
	lenEndSplit := len(strings.Split(settingsIn.SearchPattern, "~"))
	lenEscStar := len(strings.Split(settingsIn.SearchPattern, "\\*"))
	lenEscEnd := len(strings.Split(settingsIn.SearchPattern, "\\~"))
	if lenStarSplit > lenEscStar || lenEndSplit > lenEscEnd {
		if settingsIn.CheckFuzzy || settingsIn.CheckLetters {
			fmt.Println("Error search restriction on pattern-search. Use --help")
			os.Exit(-1)
		}

		return FindRestriced(tempLine, settingsIn.SearchPattern)
	}

	tempSearch := strings.ReplaceAll(settingsIn.SearchPattern, "\\~", "~")
	tempSearch = strings.ReplaceAll(tempSearch, "\\*", "*")

	//check for right search
	if settingsIn.CheckNormal {
		return FindExact(tempLine, tempSearch)
	}

	if settingsIn.CheckLetters {
		return FindChars(tempLine, tempSearch)
	}

	if settingsIn.CheckFuzzy {
		return FindFuzzy(tempLine, tempSearch)
	}

	return false, []int{}
}

func FindTextInBuff(buffIn string, settingsIn Settings) []Location {
	locationList := []Location{}
	if !utf8.ValidString(buffIn) {
		return locationList
	}

	fileLines := strings.Split(string(buffIn), "\n") //get lines
	for i, lineIter := range fileLines {
		var found bool
		var index []int
		found, index = FindTextInLine(lineIter, settingsIn)
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
		if SettingsIn.CheckFirst {
			fmt.Println("Error: check first on file-name search. Use --help")
			os.Exit(-1)
		}

		_, fileName := path.Split(pathIn)
		found, index := FindTextInLine(fileName, SettingsIn)
		if found {
			locationList = append(locationList, Location{Path: pathIn, Line: fileName, LineNum: 0, CharNum: index})
		}

		return locationList
	}

	dat, err := os.ReadFile(pathIn)
	if err != nil {
		fmt.Println(err)
		os.Exit(-1)
	}

	if !utf8.ValidString(string(dat[len(dat)/5:])) { //check for binary-file
		if pathIn == SettingsIn.Path {
			fmt.Printf("%s is a binary file\n", pathIn)
		}

		return locationList
	}

	fileLines := strings.Split(string(dat), "\n") //get lines
	for i := 0; i < len(fileLines); i++ {
		found, index := FindTextInLine((fileLines[i]), SettingsIn)
		if found {
			locationList = append(locationList, Location{Path: pathIn, Line: fileLines[i], LineNum: i, CharNum: index})
			if SettingsIn.CheckFirst {
				return locationList
			}
		}
	}

	return locationList
}
