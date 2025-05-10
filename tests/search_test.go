package tests

import (
	finder "finder/search"
	"math/rand/v2"
	"strings"
	"testing"
)

func reverse(strIn string) string {
	returnStr := ""
	for i := len(strIn) - 1; i >= 0; i-- {
		returnStr += string(strIn[i])
	}

	return returnStr
}

func randomString(l int) string {
	min := 65
	max := 90
	bytes := make([]byte, l)
	for i := 0; i < l; i++ {
		bytes[i] = byte(rand.IntN(max-min) + min)
	}

	return string(bytes)
}

func TestFindExact(t *testing.T) {
	testCases := []struct {
		name     string
		text     string
		pattern  string
		wantFind bool
		wantLen  int
	}{
		{
			name:     "simple match",
			text:     "ABCDEF",
			pattern:  "CD",
			wantFind: true,
			wantLen:  2,
		},
		{
			name:     "no match",
			text:     "ABCDEF",
			pattern:  "XY",
			wantFind: false,
			wantLen:  0,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			found, indices := finder.FindExact(tc.text, tc.pattern)
			if found != tc.wantFind {
				t.Errorf("got found=%v, want %v", found, tc.wantFind)
			}

			if len(indices) != tc.wantLen {
				t.Errorf("got %d indices, want %d", len(indices), tc.wantLen)
			}

			if found {
				// Check if indices are sequential
				for i := 1; i < len(indices); i++ {
					if indices[i] != indices[i-1]+1 {
						t.Errorf("indices not sequential: %v", indices)
					}
				}
			}
		})
	}
}

func TestFindChars(t *testing.T) {
	testCases := []struct {
		name     string
		text     string
		pattern  string
		wantFind bool
		wantLen  int
	}{
		{
			name:     "simple match",
			text:     "package",
			pattern:  "ack",
			wantFind: true,
			wantLen:  3,
		},
		{
			name:     "no match",
			text:     "ABCDEF",
			pattern:  "XY",
			wantFind: false,
			wantLen:  0,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			found, indices := finder.FindChars(tc.text, tc.pattern)
			if found != tc.wantFind {
				t.Errorf("got found=%v, want %v", found, tc.wantFind)
			}
			if len(indices) != tc.wantLen {
				t.Errorf("got %d indices, want %d", len(indices), tc.wantLen)
			}
		})
	}
}

func TestFindFuzzy(t *testing.T) {
	testCases := []struct {
		name     string
		text     string
		pattern  string
		wantFind bool
		wantLen  int
	}{
		{
			name:     "simple match",
			text:     "package",
			pattern:  "puckage",
			wantFind: true,
			wantLen:  6,
		},
		{
			name:     "no match bigger patten",
			text:     "ABCDEF",
			pattern:  "ABFFDEF",
			wantFind: false,
			wantLen:  0,
		},
		{
			name:     "no match bigger text",
			text:     "ABCFFDEF",
			pattern:  "ABCDEF",
			wantFind: false,
			wantLen:  0,
		},
		{
			name:     "full match",
			text:     "ABCEF",
			pattern:  "ABCDEF",
			wantFind: true,
			wantLen:  5,
		},
		{
			name:     "missing match",
			text:     "ABCDEF",
			pattern:  "ABCEF",
			wantFind: true,
			wantLen:  5,
		},
		{
			name:     "simple match other",
			text:     "sdfskdfjsakjfaks jtest adsfasfsafaf",
			pattern:  "teest",
			wantFind: true,
			wantLen:  4,
		},
		{
			name:     "simple match other",
			text:     "test",
			pattern:  "test",
			wantFind: true,
			wantLen:  4,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			found, indices := finder.FindFuzzy(tc.text, tc.pattern)
			if found != tc.wantFind {
				t.Errorf("got found=%v, want %v", found, tc.wantFind)
			}
			if len(indices) != tc.wantLen {
				t.Errorf("got %d indices, want %d", len(indices), tc.wantLen)
			}
		})
	}
}

func TestFindResticted(t *testing.T) {
	testCases := []struct {
		name     string
		text     string
		pattern  string
		wantFind bool
		wantLen  int
	}{
		{
			name:     "simple match",
			text:     "package",
			pattern:  "p*age",
			wantFind: true,
			wantLen:  4,
		},
		{
			name:     "simple match inside",
			text:     "ABCDEF",
			pattern:  "*CD*",
			wantFind: true,
			wantLen:  2,
		},
		{
			name:     "simple match around",
			text:     "ABCDEF",
			pattern:  "AB*DE",
			wantFind: true,
			wantLen:  4,
		},
		{
			name:     "no match end",
			text:     "ABCDEF",
			pattern:  "*FE*",
			wantFind: false,
			wantLen:  0,
		},
		{
			name:     "no match switched",
			text:     "ABCDEF",
			pattern:  "*EF*AB",
			wantFind: false,
			wantLen:  0,
		},
		{
			name:     "simple match end end",
			text:     "ABCDEF",
			pattern:  "EF~",
			wantFind: true,
			wantLen:  2,
		},
		{
			name:     "simple match end not end",
			text:     "ABCDEF",
			pattern:  "AB~",
			wantFind: false,
			wantLen:  0,
		},
		{
			name:     "simple mismatch end end",
			text:     "ABCDEF",
			pattern:  "FE~",
			wantFind: false,
			wantLen:  0,
		},
		{
			name:     "simple match all",
			text:     "ABCDEF",
			pattern:  "AB*EF~",
			wantFind: true,
			wantLen:  4,
		},
		{
			name:     "simple escape",
			text:     "ABCDEF",
			pattern:  "AB\\*EF\\~",
			wantFind: false,
			wantLen:  0,
		},
		{
			name:     "simple find match end",
			text:     "ABCDEF~",
			pattern:  "EF\\~",
			wantFind: true,
			wantLen:  3,
		},
		{
			name:     "simple find match star",
			text:     "ABCDEF*",
			pattern:  "EF\\*",
			wantFind: true,
			wantLen:  3,
		},
		{
			name:     "simple find match star + end",
			text:     "ABCD*EF~",
			pattern:  "CD\\*EF\\~",
			wantFind: true,
			wantLen:  6,
		},
		{
			name:     "simple in check star",
			text:     "test~test",
			pattern:  "st*te",
			wantFind: true,
			wantLen:  4,
		},
		{
			name:     "simple in check star2",
			text:     "testtest",
			pattern:  "st*te",
			wantFind: true,
			wantLen:  4,
		},
		{
			name:     "simple in check star end",
			text:     "test~test",
			pattern:  "st*test~",
			wantFind: true,
			wantLen:  6,
		},
		{
			name:     "simple in check star end symbol",
			text:     "test~test",
			pattern:  "st*\\~test~",
			wantFind: true,
			wantLen:  7,
		},
		{
			name:     "simple in check star end symbol double",
			text:     "test~test",
			pattern:  "st*\\~*test~",
			wantFind: true,
			wantLen:  7,
		},
		{
			name:     "end pattern + star in text not at end",
			text:     "package tests",
			pattern:  "pac*t*s~",
			wantFind: true,
			wantLen:  5,
		},
		{
			name:     "end pattern in text not at end",
			text:     "package tests",
			pattern:  "s~",
			wantFind: true,
			wantLen:  1,
		},
		{
			name:     "mismatch end test",
			text:     "package tests",
			pattern:  "st~",
			wantFind: false,
			wantLen:  0,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			found, indices := finder.FindRestriced(tc.text, tc.pattern)
			if found != tc.wantFind {
				t.Errorf("got found=%v, want %v", found, tc.wantFind)
				t.Error(tc.pattern)
			}
			if len(indices) != tc.wantLen {
				t.Errorf("got %d indices, want %d", len(indices), tc.wantLen)
			}
		})
	}
}

func TestFindExactProp(t *testing.T) {
	max := 20
	sta := 5
	end := 15
	str := randomString(max)
	tempLen := 0
	if reverse(str) == str {
		tempLen = max
	}
	testCases := []struct {
		name     string
		text     string
		pattern  string
		wantFind bool
		wantLen  int
	}{
		{
			name:     "simple match",
			text:     str,
			pattern:  str[sta:end],
			wantFind: true,
			wantLen:  (end - sta),
		},
		{
			name:     "reverse match",
			text:     str,
			pattern:  reverse(str),
			wantFind: reverse(str) == str,
			wantLen:  tempLen,
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			found, indices := finder.FindExact(tc.text, tc.pattern)
			if found != tc.wantFind {
				t.Errorf("got found=%v, want %v", found, tc.wantFind)
			}
			if len(indices) != tc.wantLen {
				t.Errorf("got %d indices, want %d", len(indices), tc.wantLen)
			}
		})
	}
}

func TestFindFuzzyProp(t *testing.T) {
	max := 20
	sta := 5
	str := randomString(max)
	testCases := []struct {
		name     string
		text     string
		pattern  string
		wantFind bool
		wantLen  int
	}{
		{
			name:     "missing match",
			text:     str,
			pattern:  str[:sta-1] + str[sta-1:sta+1],
			wantFind: true,
			wantLen:  sta + 1,
		},
		{
			name:     "wrong match",
			text:     str,
			pattern:  str[:sta-1] + strings.ToLower(string(str[sta])) + str[sta-1:sta+1],
			wantFind: true,
			wantLen:  sta + 1,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			found, indices := finder.FindFuzzy(tc.text, tc.pattern)
			if found != tc.wantFind {
				t.Errorf("got found=%v, want %v", found, tc.wantFind)
			}
			if len(indices) != tc.wantLen {
				t.Errorf("got %d indices, want %d", len(indices), tc.wantLen)
			}
		})
	}
}

func TestCaseSearch(t *testing.T) {
	testCases := []struct {
		name     string
		line     string
		settings finder.Settings
		result   bool
	}{
		{
			name: "not found lower",
			line: "TEST",
			settings: finder.Settings{
				LevelRest:          false, //-l
				LevelRestLimit:     -1,
				CheckLetters:       false, //-i
				CheckFuzzy:         true, //-c
				CheckNormal:        false,
				CheckFileName:      false, //-f
				CheckCaseSensitive: true,  //-s
				CheckFirst:         false, //-cf
				ShowInfo:           true,  //-n
				ShowColor:          true,  //-t
				ShowPathOnly:       false, //-po
				PipeInput:          true,
				ReadPipeFileList:   false, //-fl
				PathDepth:          0,
				Path:               "",
				SearchPattern:      "test",
			},
			result: false,
		},
		{
			name: "found lower",
			line: "TEST",
			settings: finder.Settings{
				LevelRest:          false, //-l
				LevelRestLimit:     -1,
				CheckLetters:       false, //-i
				CheckFuzzy:         true, //-c
				CheckNormal:        false,
				CheckFileName:      false, //-f
				CheckCaseSensitive: true,  //-s
				CheckFirst:         false, //-cf
				ShowInfo:           true,  //-n
				ShowColor:          true,  //-t
				ShowPathOnly:       false, //-po
				PipeInput:          true,
				ReadPipeFileList:   false, //-fl
				PathDepth:          0,
				Path:               "",
				SearchPattern:      "test",
			},
			result: true,
		},
		{
			name: "found lower fuzzy",
			line: "TEST",
			settings: finder.Settings{
				LevelRest:          false, //-l
				LevelRestLimit:     -1,
				CheckLetters:       false, //-i
				CheckFuzzy:         true, //-c
				CheckNormal:        false,
				CheckFileName:      false, //-f
				CheckCaseSensitive: true,  //-s
				CheckFirst:         false, //-cf
				ShowInfo:           true,  //-n
				ShowColor:          true,  //-t
				ShowPathOnly:       false, //-po
				PipeInput:          true,
				ReadPipeFileList:   false, //-fl
				PathDepth:          0,
				Path:               "",
				SearchPattern:      "tst",
			},
			result: true,
		},
		{
			name: "found lower in line",
			line: "ABCDEFGHIJKLMNOP",
			settings: finder.Settings{
				LevelRest:          false, //-l
				LevelRestLimit:     -1,
				CheckLetters:       true, //-i
				CheckFuzzy:         false, //-c
				CheckNormal:        false,
				CheckFileName:      false, //-f
				CheckCaseSensitive: true,  //-s
				CheckFirst:         false, //-cf
				ShowInfo:           true,  //-n
				ShowColor:          true,  //-t
				ShowPathOnly:       false, //-po
				PipeInput:          true,
				ReadPipeFileList:   false, //-fl
				PathDepth:          0,
				Path:               "",
				SearchPattern:      "afg",
			},
			result: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			locations := finder.FindTextInBuff(tc.line, tc.settings)
			for _, elem := range(locations) {
				if (len(elem.CharNum) != 0) != tc.result {
					t.Error("wrong reuslt in case sensitive test")
				}
			}
		})
	}
}
