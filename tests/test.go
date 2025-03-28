package tests

import (
	"math/rand/v2"
	"testing"

	main "github.com/lukas-fohl/goFind/src"
)

var letterRunes = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

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
	testString := randomString(len)
	found, idxs := main.FindExact(&testString, testString[sta:end])
	if !found {
		t.Error("error in FindTextInLine")
	}
	if len(idxs) != end-sta {
		t.Error("error in FindTextInLine arg len")
	}
}
