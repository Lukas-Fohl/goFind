package tests

import (
	finder "finder/search"
	"path/filepath"
	"testing"
)

func TestPath(t *testing.T) {
	path, err := filepath.Abs(".")
	if err != nil {
		panic(err)
	}
	testCases := []struct {
		name       string
		flagsIn    []string
		pathOut    string
		patternOut string
	}{
		{
			name:       "implicit path",
			flagsIn:    []string{"main", "package", "-c"},
			pathOut:    path,
			patternOut: "package",
		},
		{
			name:       "explicit path",
			flagsIn:    []string{"main", ".", "package", "-i"},
			pathOut:    path,
			patternOut: "package",
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			res := finder.FlagHandle(tc.flagsIn)
			if res.Path != tc.pathOut {
				t.Errorf("wrong path, want: %s got: %s", tc.pathOut, res.Path)
			}

			if res.SearchPattern != tc.patternOut {
				t.Errorf("wrong pattern, want: %s got: %s", tc.patternOut, res.SearchPattern)
			}
		})
	}
}

func TestFlag(t *testing.T) {
	testCases := []struct {
		name    string
		flagsIn []string
		result  finder.Settings
	}{
		{
			name:    "fuzzy flag",
			flagsIn: []string{"main", "package", "-c"},
			result: finder.Settings{
				LevelRest:          false, //-l
				LevelRestLimit:     -1,
				CheckLetters:       false, //-i
				CheckFuzzy:         true,  //-c
				CheckNormal:        false,
				CheckFileName:      false, //-f
				CheckCaseSensitive: true,  //-s
				ShowInfo:           true,  //-n
				PathDepth:          0,
				Path:               "",
				SearchPattern:      "",
			},
		},
		{
			name:    "level flag",
			flagsIn: []string{"main", "package", "-l", "1"},
			result: finder.Settings{
				LevelRest:          true, //-l
				LevelRestLimit:     1,
				CheckLetters:       false, //-i
				CheckFuzzy:         false, //-c
				CheckNormal:        true,
				CheckFileName:      false, //-f
				CheckCaseSensitive: true,  //-s
				ShowInfo:           true,  //-n
				PathDepth:          0,
				Path:               "",
				SearchPattern:      "",
			},
		},
		{
			name:    "letter flag",
			flagsIn: []string{"main", "package", "-i"},
			result: finder.Settings{
				LevelRest:          false, //-l
				LevelRestLimit:     -1,
				CheckLetters:       true,  //-i
				CheckFuzzy:         false, //-c
				CheckNormal:        false,
				CheckFileName:      false, //-f
				CheckCaseSensitive: true,  //-s
				ShowInfo:           true,  //-n
				PathDepth:          0,
				Path:               "",
				SearchPattern:      "",
			},
		},
		{
			name:    "filename flag",
			flagsIn: []string{"main", "package", "-f"},
			result: finder.Settings{
				LevelRest:          false, //-l
				LevelRestLimit:     -1,
				CheckLetters:       false, //-i
				CheckFuzzy:         false, //-c
				CheckNormal:        true,
				CheckFileName:      true, //-f
				CheckCaseSensitive: true, //-s
				ShowInfo:           true, //-n
				PathDepth:          0,
				Path:               "",
				SearchPattern:      "",
			},
		},
		{
			name:    "check case flag",
			flagsIn: []string{"main", "package", "-s"},
			result: finder.Settings{
				LevelRest:          false, //-l
				LevelRestLimit:     -1,
				CheckLetters:       false, //-i
				CheckFuzzy:         false, //-c
				CheckNormal:        true,
				CheckFileName:      false, //-f
				CheckCaseSensitive: false, //-s
				ShowInfo:           true,  //-n
				PathDepth:          0,
				Path:               "",
				SearchPattern:      "",
			},
		},
		{
			name:    "show info flag",
			flagsIn: []string{"main", "package", "-n"},
			result: finder.Settings{
				LevelRest:          false, //-l
				LevelRestLimit:     -1,
				CheckLetters:       false, //-i
				CheckFuzzy:         false, //-c
				CheckNormal:        true,
				CheckFileName:      false, //-f
				CheckCaseSensitive: true,  //-s
				ShowInfo:           false, //-n
				PathDepth:          0,
				Path:               "",
				SearchPattern:      "",
			},
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			res := finder.FlagHandle(tc.flagsIn)
			if res.LevelRest != tc.result.LevelRest {
				t.Error("wrong -l handle")
			}

			if res.LevelRestLimit != tc.result.LevelRestLimit {
				t.Errorf("wrong level limit, want: %d, got: %d", tc.result.LevelRestLimit, res.LevelRestLimit)
			}

			if res.CheckLetters != tc.result.CheckLetters {
				t.Error("wrong -i handle")
			}

			if res.CheckFuzzy != tc.result.CheckFuzzy {
				t.Error("wrong -c handle")
			}

			if res.CheckFileName != tc.result.CheckFileName {
				t.Error("wrong -f handle")
			}

			if res.ShowInfo != tc.result.ShowInfo {
				t.Error("wrong -n handle")
			}

			if res.CheckCaseSensitive != tc.result.CheckCaseSensitive {
				t.Error("wrong -s handle")
			}

			if res.CheckNormal != tc.result.CheckNormal {
				t.Error("normal flag not set")
			}
		})
	}
}
