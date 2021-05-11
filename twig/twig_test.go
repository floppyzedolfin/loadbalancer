package twig

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSwitch(t *testing.T) {
	// default state is OFF
	assert.False(t, debug)

	// switch ON
	Switch()
	assert.True(t, debug)

	// switch back OFF
	Switch()
	assert.False(t, debug)
}

func ExamplePrintf(t *testing.T) {
	// default state is OFF
	Printf("Hello, %s", "world")
	// Output:

	debug = true
	Printf("Hello, %s", "people")
	// Output: Hello, people
}
