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
	pacifistMainMenuTurns    int
	battleMenuCooldownFrames int
	pacifistSubmenuPending   bool
)

func pacifistOptionIndex(optionCount int, chooseMercy bool) int {
	if optionCount >= 4 && chooseMercy {
		return optionCount - 1 // MERCY is the final button in the main menu.
	}
	if optionCount >= 2 {
		return 1 // ACT is second; safer than confirming the first (FIGHT) button.
	}
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
	// Main-menu buttons share a horizontal baseline. ACT submenus are generally
	// arranged as vertical lists and must not be mistaken for the main menu.
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
