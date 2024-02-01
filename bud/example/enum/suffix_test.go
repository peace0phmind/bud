package enum

//go:generate go run ../../../main.go --file-suffix .enum.gen

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// SuffixTest @ENUM{some_item}
type SuffixTest string

func TestSuffix(t *testing.T) {
	x := Suffix("")
	assert.Equal(t, "", x.String())

	assert.Equal(t, Suffix("gen"), SuffixGen)
}

func TestSuffixTest(t *testing.T) {
	x := SuffixTest("")
	assert.Equal(t, "", x.String())

	assert.Equal(t, SuffixTest("some_item"), SuffixTestSomeItem)
}
