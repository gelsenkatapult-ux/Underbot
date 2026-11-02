package ai

import (
	"image"

	"gitlab.com/256/Underbot/cv/object"
	"gitlab.com/256/Underbot/sys"
)

// PacifistMode prevents automatic Fight-menu selection. The default strategy
// selects ACT when MERCY is not visibly available, and selects MERCY once its
// label is highlighted yellow by Undertale.
var PacifistMode = true

func pacifistOptionIndex(optionCount int, mercyHighlighted bool) int {
	if optionCount <= 1 {
		return 0
	}
	if optionCount >= 4 {
		if mercyHighlighted {
			return optionCount - 1
		}
		return 1 // ACT is the second option in Undertale's main battle menu.
	}
	// In an ACT submenu, prefer the first listed action rather than a random
	// menu entry. Game-specific ACT text still needs explicit recognition.
	return 0
}

// hasYellowMenuText checks the supplied menu button for the yellow text used
// by Undertale to indicate that the enemy can currently be spared.
func hasYellowMenuText(img *image.RGBA, bounds image.Rectangle) bool {
	if img == nil {
		return false
	}
	bounds = bounds.Intersect(img.Bounds()).Inset(5)
	if bounds.Empty() {
		return false
	}
	yellowPixels := 0
	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			r, g, b, _ := img.At(x, y).RGBA()
			if r >= 0xC000 && g >= 0x9000 && b <= 0x3000 {
				yellowPixels++
				if yellowPixels >= 4 {
					return true
				}
			}
		}
	}
	return false
}

// AttackGoalUpdate deliberately does not press the timing key in pacifist
// mode. If the fight-timing state appears, it was not selected by the
// non-violent menu policy and must not be confirmed automatically.
func AttackGoalUpdate(objects []object.Object, win sys.Window, img *image.RGBA) error {
	genUpdate()
	if PacifistMode {
		return nil
	}
	return BattleMenuUpdate(objects, win, img)
}
