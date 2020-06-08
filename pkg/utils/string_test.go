package utils

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRemoveWhitepaces(t *testing.T) {
	initialString := " \t r\r\ng b    "
	expected := "rgb"

	assert.Equal(t, expected, RemoveWhitepaces(initialString))
}
