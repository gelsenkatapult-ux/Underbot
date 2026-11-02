package combat

import (
	"image"\n\t"time"

	"gitlab.com/256/Underbot/cv/object"
)

// BattleGeometry is the validated portion of a combat frame that can be
// derived from the legacy recognizer. Coordinates are in client-image space.
type BattleGeometry struct {
	Board    image.Rectangle
	Player   image.Rectangle
	SoulMode SoulMode
}

// FindBattleGeometry requires both the recognized fight box and a recognized
// SOUL. It never guesses geometry from arbitrary dark rectangles.
func FindBattleGeometry(objects []object.Object) (BattleGeometry, bool) {
	var result BattleGeometry
	boardFound, playerFound := false, false
	for i := range objects {
		detected := &objects[i]
		if !detected.Recognized || detected.RecogObj.Parent == nil {
			continue
		}
		name := detected.RecogObj.Type.Name
		bounds := detected.RecogObj.Parent.Bounds
		if bounds.Empty() {
			continue
		}
		if name == "fightBox" {
			if !boardFound || bounds.Dx()*bounds.Dy() > result.Board.Dx()*result.Board.Dy() {
				result.Board = bounds
				boardFound = true
			}
		}
		var mode SoulMode
		switch name {
		case "redHeart":
			mode = SoulRed
		case "blueHeart":
			mode = SoulBlue
		case "greenHeart":
			mode = SoulGreen
		default:
			continue
		}
		if !playerFound {
			result.Player = bounds
			result.SoulMode = mode
			playerFound = true
		} else if bounds.Dx()*bounds.Dy() > result.Player.Dx()*result.Player.Dy() {
			result.Player = bounds
			result.SoulMode = mode
		}
	}
	if !boardFound || !playerFound || result.Board.Empty() || result.Player.Empty() {
		return BattleGeometry{}, false
	}
	if !result.Player.In(result.Board) {
		return BattleGeometry{}, false
	}
	return result, true
}

// FrameObservation is a conservative combat snapshot. It deliberately carries
// detections and tracks without issuing movement or keypresses.
type FrameObservation struct {
	Geometry   BattleGeometry
	Projectiles []ProjectileObservation
	Tracks     []ProjectileTrack
}

// ObserveFrame validates the fight box and SOUL before running color detection.
// The tracker is supplied by the caller so separate game windows stay isolated.
func ObserveFrame(img *image.RGBA, objects []object.Object, tracker *Tracker) (FrameObservation, bool) {
	geometry, ok := FindBattleGeometry(objects)
	if !ok || img == nil {
		return FrameObservation{}, false
	}
	projectiles := DetectProjectiles(img, geometry.Board, []image.Rectangle{geometry.Player})
	if tracker == nil {
		tracker = NewTracker()
	}
	tracks := tracker.Update(projectiles, time.Now())
	return FrameObservation{
		Geometry: geometry,
		Projectiles: projectiles,
		Tracks: tracks,
	}, true
}
