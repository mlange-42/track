package util

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFormat(t *testing.T) {
	assert.Equal(
		t,
		"foo baz bar",
		Format("foo {name} bar", map[string]string{"name": "baz"}),
		"Simple replacement not working",
	)
	assert.Equal(
		t,
		"foo baz bar baz",
		Format("foo {name} bar {name}", map[string]string{"name": "baz"}),
		"Repetitions not working",
	)
	assert.Equal(
		t,
		"foo baz bar foo",
		Format("foo {name} bar {name2}", map[string]string{"name": "baz", "name2": "foo"}),
		"Repetitions not working",
	)
}
