package ai

import (
	"image"
	"testing"

	"gitlab.com/256/Underbot/cv/object"
)

func recognizedObject(name string, bounds image.Rectangle) object.Object {
	return object.Object{
		Bounds:     bounds,
		Recognized: true,
		RecogObj: object.RecognizedObject{
			Type: object.RecognizableObject{Name: name},
		},
	}
}

func TestMovementDirectionUsesSingleAxis(t *testing.T) {
	tests := []struct {
		name string
		from image.Point
		to   image.Point
		want string
	}{
		{"right", image.Pt(10, 10), image.Pt(20, 10), "right"},
		{"left", image.Pt(20, 10), image.Pt(10, 10), "left"},
		{"down", image.Pt(10, 10), image.Pt(10, 20), "down"},
		{"up", image.Pt(10, 20), image.Pt(10, 10), "up"},
		{"dominant x", image.Pt(10, 10), image.Pt(20, 15), "right"},
		{"dominant y", image.Pt(10, 10), image.Pt(15, 20), "down"},
		{"same point", image.Pt(10, 10), image.Pt(10, 10), ""},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			if got := movementDirection(test.from, test.to); got != test.want {
				t.Fatalf("movementDirection(%v, %v) = %q, want %q", test.from, test.to, got, test.want)
			}
		})
	}
}

func TestChooseWorldGoalPrefersKnownLandmarks(t *testing.T) {
	objects := []object.Object{
		recognizedObject("entrance", image.Rect(550, 150, 600, 250)),
		recognizedObject("yellowSwitch", image.Rect(100, 100, 120, 140)),
	}
	got, name := chooseWorldGoal(objects, image.Pt(200, 200), image.Pt(639, 240))
	if name != "yellowSwitch" {
		t.Fatalf("goal name = %q, want yellowSwitch", name)
	}
	if got != image.Pt(110, 120) {
		t.Fatalf("goal point = %v, want switch centre (110,120)", got)
	}
}

func TestChooseWorldGoalUsesNearestLandmarkOfSameType(t *testing.T) {
	objects := []object.Object{
		recognizedObject("entrance", image.Rect(500, 100, 520, 120)),
		recognizedObject("entrance", image.Rect(200, 200, 220, 220)),
	}
	got, name := chooseWorldGoal(objects, image.Pt(210, 210), image.Pt(639, 240))
	if name != "entrance" || got != image.Pt(210, 210) {
		t.Fatalf("goal = %v (%q), want nearest entrance centre (210,210)", got, name)
	}
}

func TestChooseWorldGoalIgnoresUnrecognizedObjects(t *testing.T) {
	candidate := recognizedObject("yellowSwitch", image.Rect(100, 100, 120, 120))
	candidate.Recognized = false
	fallback := image.Pt(639, 240)
	got, name := chooseWorldGoal([]object.Object{candidate}, image.Pt(20, 20), fallback)
	if got != fallback || name != "" {
		t.Fatalf("goal = %v (%q), want fallback %v", got, name, fallback)
	}
}

func TestInteractionRangeAndPolicy(t *testing.T) {
	if !worldGoalNeedsInteraction("yellowSwitch") || !worldGoalNeedsInteraction("dummy") {
		t.Fatal("switch and dummy landmarks should be interactable")
	}
	if worldGoalNeedsInteraction("entrance") {
		t.Fatal("room entrances should be entered by walking, not repeated Z presses")
	}
	if !withinWorldInteractionRange(image.Pt(100, 100), image.Pt(120, 100)) {
		t.Fatal("point 20 pixels away should be within interaction range")
	}
	if withinWorldInteractionRange(image.Pt(100, 100), image.Pt(140, 100)) {
		t.Fatal("point 40 pixels away should be outside interaction range")
	}
}
