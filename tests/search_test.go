package tests

import (
	finder "finder/search"
	"math/rand/v2"
	"testing"
)

func randomString(l int) string {
	min := 65
	max := 90
	bytes := make([]byte, l)
	for i := 0; i < l; i++ {
		bytes[i] = byte(rand.IntN(max-min) + min)
	}
	return string(bytes)
}

func TestFindTextInLine(t *testing.T) {
	len := 20
	sta := 5
	end := 15
	prefix := "error in FindTextInLine - "
	testString := randomString(len)
	idxs := []int{}
	found := false
	found, idxs = finder.FindExact(&testString, testString[sta:end])
	if !found {
		t.Error(prefix + "not found")
	} else {
		if len(idxs) != end-sta {
			t.Error(prefix + "wrong length")
		}
	}
	if idxs[0]-idxs[len(idxs)-1] != end-sta {
		t.Error(prefix + "wrong chars found")
	}
}
