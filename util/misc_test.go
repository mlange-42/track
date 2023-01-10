package util

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPair(t *testing.T) {
	pair := NewPair("a", 1)
	assert.Equal(t, Pair[string, int]{"a", 1}, pair, "Wrong pair")
}

func TestReverse(t *testing.T) {
	arr := []int{1, 2, 3}
	Reverse(arr)
	assert.Equal(t, []int{3, 2, 1}, arr, "Wrong reversed slice")
}

func TestUnique(t *testing.T) {
	arrInt := []int{1, 2, 3, 3, 1, 4}
	uniqueInt := Unique(arrInt)
	assert.Equal(t, []int{1, 2, 3, 4}, uniqueInt, "Wrong unique int values")

	arrStr := []string{"1", "2", "3", "3", "1", "4"}
	uniqueStr := Unique(arrStr)
	assert.Equal(t, []string{"1", "2", "3", "4"}, uniqueStr, "Wrong unique string values")
}
