package ai

import (
	"image"

	"gitlab.com/256/Underbot/cv/object"
	"gitlab.com/256/Underbot/sys"
)

// PacifistMode prevents automatic Fight-menu selection. The AI alternates
// non-violent ACT turns with attempts to Spare from the Mercy menu.
var PacifistMode = true

// These counters are updated by the game-state loop; a short cooldown prevents
// the same visible menu entry from being confirmed repeatedly every frame.
var (
	pacifistMainMenuTurns  int
	battleMenuCooldownFrames int
)

func pacifistOptionIndex(optionCount int, chooseMercy bool) int {
	if optionCount <= 1 {
		return 0
	}
	if optionCount >= 4 {
		if chooseMercy {
			return optionCount - 1
		}
		return 1 // ACT is the second option in Undertale's main battle menu; the final option is MERCY.
	}
	// In an ACT submenu, prefer the first listed action rather than a random
	// menu entry. Game-specific ACT text still needs explicit recognition.
	return 0
}

func isMainBattleMenu(options []object.RecognizedObject) bool {
	if len(options) != 4 {
		return false
	}
	minY := options[0].Parent.Bounds.Min.Y + options[0].Parent.Bounds.Dy()/2
	maxY := minY
	for _, option := range options[1:] {
		centerY := option.Parent.Bounds.Min.Y + option.Parent.Bounds.Dy()/2
		if centerY < minY {
			minY = centerY
		}
		if centerY > maxY {
			maxY = centerY
		}
	}
	// The four main buttons share a horizontal baseline. ACT submenus are
	// generally arranged as a vertical list and must not be mistaken for them.
	return maxY-minY <= 12
}

// AttackGoalUpdate deliberately does not press the timing key in pacifist
// mode. If the fight-timing state appears, it must not be confirmed.
func AttackGoalUpdate(objects []object.Object, win sys.Window, img *image.RGBA) error {
	genUpdate()
	if PacifistMode {
		return nil
	}
	return BattleMenuUpdate(objects, win, img)
}
