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
