package main

import (
	"math"
	"testing"
)

//Beispiel (Test funzt nicht, warum auch immer ?!)
func TestToRad(t *testing.T) {
	var res = ToRad(180)
	if res != math.Pi {
		t.Errorf("Sum was inccorect, got: %f, wanted: %f", res, math.Pi)
	}
}
func TestDistanceBetweenToPoints(t *testing.T) {
	var result = distanceBetweenTwoPoints(49.3547198000, 9.1508659200, 49.3546998700, 9.1509324100)
	if result != 5.301269557828212 {
		t.Errorf("Distance was wrong, got: %f, wanted %f", result, 5.301269557828212)
	}
}

//without timebetween
func TestSpeed(t *testing.T) {
	x, _ := speed(5.301269557828212, "2019-09-14T13:14:00.000Z", "2019-09-14T13:14:10.003Z")
	if x != 0.5299679653932032 {
		t.Errorf("Wrong Results, got: %f, wanted %f", x, 0.5299679653932032)
	}
}
