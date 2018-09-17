package main

import (
	"image/color"
	"testing"

	"github.com/hihebark/gore/core"
)

func TestRGBAtoXYZ(t *testing.T) {
	c := core.RGBAtoXYZ(color.RGBA{255, 255, 0, 0})
	if c.X == 0 && c.Y == 0 && c.Z == 0 {
		t.Errorf("TestRGBAtoXYZ error on converting to xyz nil detect XYZ = %v", c)
	}
}
