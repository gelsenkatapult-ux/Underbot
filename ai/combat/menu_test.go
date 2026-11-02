package combat

import (
	"image"
	"testing"
)

func TestPacifistMenuTargetSelectsACTByPosition(t *testing.T) {
	options := []image.Rectangle{
		image.Rect(10, 200, 110, 240),   // FIGHT
		image.Rect(150, 200, 250, 240),  // ACT
		image.Rect(290, 200, 390, 240),  // ITEM
		image.Rect(430, 200, 530, 240),  // MERCY
	}
	got, ok := PacifistMenuTarget(options)
	if !ok {
		t.Fatal("expected standard battle menu to be recognized")
	}
	if got != (image.Point{X: 200, Y: 220}) {
		t.Fatalf("target=%v, want ACT center (200,220)", got)
	}
}

func TestPacifistMenuTargetRejectsAmbiguousMenus(t *testing.T) {
	cases := [][]image.Rectangle{
		nil,
		{image.Rect(0, 0, 10, 10)},
		{image.Rect(0, 0, 10, 10), image.Rect(20, 0, 30, 10), image.Rect(40, 0, 50, 10), image.Rect(60, 40, 70, 50)},
	}
	for _, options := range cases {
		if _, ok := PacifistMenuTarget(options); ok {
			t.Fatalf("ambiguous layout accepted: %v", options)
		}
	}
}

func TestPacifistMenuTargetDoesNotDependOnDetectorOrder(t *testing.T) {
	options := []image.Rectangle{
		image.Rect(430, 200, 530, 240),
		image.Rect(290, 200, 390, 240),
		image.Rect(10, 200, 110, 240),
		image.Rect(150, 200, 250, 240),
	}
	got, ok := PacifistMenuTarget(options)
	if !ok || got.X != 200 {
		t.Fatalf("target=%v, recognized=%v; want ACT x=200", got, ok)
	}
}
