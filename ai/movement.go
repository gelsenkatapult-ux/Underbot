package ai

import (
	"image"
	"math"

	"gitlab.com/256/Underbot/cv/object"
)

// movementDirection chooses one axis at a time, matching the 4-way A* grid.
func movementDirection(from, to image.Point) string {
	dx := to.X - from.X
	dy := to.Y - from.Y
	if dx == 0 && dy == 0 {
		return ""
	}
	if abs(dx) >= abs(dy) && dx != 0 {
		if dx < 0 {
			return "left"
		}
		return "right"
	}
	if dy < 0 {
		return "up"
	}
	return "down"
}

func abs(value int) int {
	if value < 0 {
		return -value
	}
	return value
}

// chooseWorldGoal prefers visible, recognizable room objectives over the old
// unconditional "walk to the right-middle edge" target. It does not invent a
// room map: when no supported landmark is visible, the old exit heuristic is
// retained as a fallback.
func chooseWorldGoal(objects []object.Object, current, fallback image.Point) (image.Point, string) {
	priorities := []string{
		"yellowSwitch",
		"entrance",
		"dummy",
		"torielFront",
		"torielSide",
	}
	for _, wanted := range priorities {
		var nearest image.Point
		bestDistance := math.MaxInt
		found := false
		for _, candidate := range objects {
			if !candidate.Recognized || candidate.RecogObj.Type.Name != wanted {
				continue
			}
			point := image.Pt(
				(candidate.Bounds.Min.X+candidate.Bounds.Max.X)/2,
				(candidate.Bounds.Min.Y+candidate.Bounds.Max.Y)/2,
			)
			dx := point.X - current.X
			dy := point.Y - current.Y
			distance := dx*dx + dy*dy
			if distance < bestDistance {
				bestDistance = distance
				nearest = point
				found = true
			}
		}
		if found {
			return nearest, wanted
		}
	}
	return fallback, ""
}

func worldGoalNeedsInteraction(name string) bool {
	switch name {
	case "yellowSwitch", "dummy", "torielFront", "torielSide":
		return true
	default:
		return false
	}
}

func withinWorldInteractionRange(from, to image.Point) bool {
	dx := from.X - to.X
	dy := from.Y - to.Y
	return dx*dx+dy*dy <= 28*28
}
