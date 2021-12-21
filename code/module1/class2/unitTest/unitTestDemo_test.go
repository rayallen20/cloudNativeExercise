package main

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestAdd(t *testing.T) {
	t.Log("start testing")
	result := add(1, 2)
	assert.Equal(t, result, 3)
}
