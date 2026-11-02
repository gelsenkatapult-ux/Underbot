package ai

import (
	"image"
	"image/color"
	"testing"

	"gitlab.com/256/Underbot/cv/object"
)

func battleObject(id int, bounds image.Rectangle, col color.Color, recognized bool) object.Object {
	return object.Object{
		ID:         id,
		Bounds:     bounds,
		Color:      col,
		Recognized: recognized,
	}
}

func TestFindBattleHazardsFiltersContours(t *testing.T) {
	arena := image.Rect(100, 100, 264, 239)
	objects := []object.Object{
		battleObject(1, image.Rect(170, 160, 180, 170), color.RGBA{R: 255, A: 255}, true), // player heart
		battleObject(2, image.Rect(200, 160, 210, 170), color.RGBA{R: 255, G: 255, B: 255, A: 255}, false), // white projectile
		battleObject(3, image.Rect(130, 180, 142, 192), color.RGBA{R: 255, G: 150, A: 255}, false), // orange projectile
		battleObject(4, image.Rect(150, 150, 220, 225), color.RGBA{R: 255, A: 255}, false), // too large
		battleObject(5, image.Rect(160, 160, 168, 168), color.RGBA{R: 20, G: 20, B: 20, A: 255}, false), // too dark
		battleObject(6, image.Rect(102, 170, 112, 180), color.RGBA{R: 255, A: 255}, false), // arena edge inset
		battleObject(7, image.Rect(270, 170, 280, 180), color.RGBA{R: 255, A: 255}, false), // outside
	}
	got := findBattleHazards(objects, arena, 1)
	if len(got) != 2 {
		t.Fatalf("findBattleHazards returned %d candidates, want 2", len(got))
	}
	if got[0].Center != image.Pt(205, 165) || got[1].Center != image.Pt(136, 186) {
		t.Fatalf("unexpected hazard centres: %#v", got)
	}
}

func TestChooseSafeBattlePointMovesAwayFromNearbyHazard(t *testing.T) {
	arena := image.Rect(0, 0, 164, 139)
	player := image.Pt(82, 70)
	hazards := []battleHazard{{Center: image.Pt(92, 70), Bounds: image.Rect(88, 66, 96, 74)}}
	got := chooseSafeBattlePoint(arena, player, hazards)
	if got == player {
		t.Fatal("expected a movement target when a projectile is close")
	}
	if !got.In(arena.Inset(12)) {
		t.Fatalf("chosen target %v is outside the safe arena interior", got)
	}
	if nearestHazardDistance(got, hazards) <= nearestHazardDistance(player, hazards) {
		t.Fatalf("chosen target %v is not safer than player position %v", got, player)
	}
}

func TestChooseSafeBattlePointDoesNothingWithoutNearbyHazards(t *testing.T) {
	arena := image.Rect(0, 0, 164, 139)
	player := image.Pt(82, 70)
	if got := chooseSafeBattlePoint(arena, player, nil); got != player {
		t.Fatalf("no-hazard target = %v, want unchanged point %v", got, player)
	}
	far := []battleHazard{{Center: image.Pt(150, 10), Bounds: image.Rect(147, 7, 153, 13)}}
	if got := chooseSafeBattlePoint(arena, player, far); got != player {
		t.Fatalf("distant-hazard target = %v, want unchanged point %v", got, player)
	}
}

func TestGreenHeartShieldFacesNearestHazard(t *testing.T) {
	player := image.Pt(80, 70)
	hazards := []battleHazard{
		{Center: image.Pt(150, 70)},
		{Center: image.Pt(80, 10)},
	}
	if got := nearestHazardDirection(player, hazards); got != "up" {
		t.Fatalf("nearestHazardDirection = %q, want up", got)
	}
}
