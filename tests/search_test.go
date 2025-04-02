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
			found, idxs := finder.FindExact(&tc.text, tc.pattern)
			if found != tc.wantFind {
				t.Errorf("got found=%v, want %v", found, tc.wantFind)
			}
			if len(idxs) != tc.wantLen {
				t.Errorf("got %d indices, want %d", len(idxs), tc.wantLen)
			}
			if found {
				// Check if indices are sequential
				for i := 1; i < len(idxs); i++ {
					if idxs[i] != idxs[i-1]+1 {
						t.Errorf("indices not sequential: %v", idxs)
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
			found, idxs := finder.FindChars(&tc.text, tc.pattern)
			if found != tc.wantFind {
				t.Errorf("got found=%v, want %v", found, tc.wantFind)
			}
			if len(idxs) != tc.wantLen {
				t.Errorf("got %d indices, want %d", len(idxs), tc.wantLen)
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
			name:     "no match",
			text:     "ABCDEF",
			pattern:  "ABFFDEF",
			wantFind: false,
			wantLen:  0,
		},
		{
			name:     "no match",
			text:     "ABCFFDEF",
			pattern:  "ABCDEF",
			wantFind: false,
			wantLen:  0,
		},
		{
			name:     "no match",
			text:     "ABCEF",
			pattern:  "ABCDEF",
			wantFind: true,
			wantLen:  5,
		},
		{
			name:     "no match",
			text:     "ABCDEF",
			pattern:  "ABCEF",
			wantFind: true,
			wantLen:  4,
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			found, idxs := finder.FindFuzzy(&tc.text, tc.pattern)
			if found != tc.wantFind {
				t.Errorf("got found=%v, want %v", found, tc.wantFind)
			}
			if len(idxs) != tc.wantLen {
				t.Errorf("got %d indices, want %d", len(idxs), tc.wantLen)
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
			found, idxs := finder.FindExact(&tc.text, tc.pattern)
			if found != tc.wantFind {
				t.Errorf("got found=%v, want %v", found, tc.wantFind)
			}
			if len(idxs) != tc.wantLen {
				t.Errorf("got %d indices, want %d", len(idxs), tc.wantLen)
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
			found, idxs := finder.FindFuzzy(&tc.text, tc.pattern)
			if found != tc.wantFind {
				t.Errorf("got found=%v, want %v", found, tc.wantFind)
			}
			if len(idxs) != tc.wantLen {
				t.Errorf("got %d indices, want %d", len(idxs), tc.wantLen)
			}
		})
	}
}
