package main

import (
	"github.com/stretchr/testify/assert"
)

func test() { /*
		actial :=1...;
		expected
		if actual !=expected{
			t.Errorf("wrong result expected %v, actual %v", expected, actual)
		}*/
	assert.Equal(t, New(1, 2, 3), New(1, 2, 3), "wrong sum")
}
