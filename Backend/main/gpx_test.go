package main

import (
	"github.com/stretchr/testify/assert"
	"math"
	"testing"
)

func TestToRad(t *testing.T) {
	var res = ToRad(180)
	assert.Equal(t, res, math.Pi, "Should be equal")
}

func TestDistanceBetweenToPoints(t *testing.T) {
	var result = distanceBetweenTwoPoints(49.3547198000, 9.1508659200, 49.3546998700, 9.1509324100)
	assert.Equal(t, result, 5.301269557828212, "Should be equal")
}

//without timebetween
func TestSpeed(t *testing.T) {
	x, _ := speed(5.301269557828212, "2019-09-14T13:14:00.000Z", "2019-09-14T13:14:10.003Z")
	assert.Equal(t, x, 0.5299679653932032, "Should be equal")
}
