package ai

import (
	"image"
	"math"

	"gitlab.com/256/Underbot/cv/object"
)

const (
	maxProjectileWidth  = 24
	maxProjectileHeight = 24
	maxProjectileArea   = 576
	battleDangerRadius  = 46.0
)

type battleHazard struct {
	Center image.Point
	Bounds image.Rectangle
}

// findBattleHazards extracts small, bright, unrecognized contours inside the
// soul arena. Undertale's projectile patterns are not uniform, so this is a
// conservative screen-space heuristic rather than a full attack classifier.
func findBattleHazards(objects []object.Object, arena image.Rectangle, heartID int) []battleHazard {
	hazards := make([]battleHazard, 0, 16)
	inner := arena.Inset(5)
	if inner.Empty() {
		return hazards
	}
	for _, candidate := range objects {
		if candidate.ID == heartID || candidate.Recognized || candidate.Color == nil {
			continue
		}
		bounds := candidate.Bounds
		if bounds.Empty() || bounds.Dx() > maxProjectileWidth || bounds.Dy() > maxProjectileHeight {
			continue
		}
		if bounds.Dx()*bounds.Dy() > maxProjectileArea || !bounds.In(inner) {
			continue
		}
		r, g, b, _ := candidate.Color.RGBA()
		maxChannel := maxUint32(r, maxUint32(g, b))
		minChannel := minUint32(r, minUint32(g, b))
		// Filter out dark floor details while retaining white, orange, cyan
		// and other high-contrast projectile colours.
		if maxChannel < 0xA000 || (minChannel < 0x6000 && maxChannel-minChannel < 0x5000) {
			continue
		}
		hazards = append(hazards, battleHazard{
			Center: image.Pt(bounds.Min.X+bounds.Dx()/2, bounds.Min.Y+bounds.Dy()/2),
			Bounds: bounds,
		})
	}
	return hazards
}

func maxUint32(a, b uint32) uint32 {
	if a > b {
		return a
	}
	return b
}

func minUint32(a, b uint32) uint32 {
	if a < b {
		return a
	}
	return b
}

func nearestHazardDistance(point image.Point, hazards []battleHazard) float64 {
	if len(hazards) == 0 {
		return math.Inf(1)
	}
	best := math.Inf(1)
	for _, hazard := range hazards {
		distance := math.Hypot(float64(point.X-hazard.Center.X), float64(point.Y-hazard.Center.Y))
		if distance < best {
			best = distance
		}
	}
	return best
}

// chooseSafeBattlePoint considers a small set of reachable positions inside
// the arena. It returns the current point unless a detected hazard is close
// enough to justify moving; distance from the player is penalized to avoid
// unnecessary long traversals across the box.
func chooseSafeBattlePoint(arena image.Rectangle, player image.Point, hazards []battleHazard) image.Point {
	if len(hazards) == 0 || arena.Empty() {
		return player
	}
	if nearestHazardDistance(player, hazards) > battleDangerRadius {
		return player
	}
	inner := arena.Inset(12)
	if inner.Empty() {
		return player
	}

	best := player
	bestScore := nearestHazardDistance(player, hazards)
	const columns, rows = 7, 5
	for row := 0; row < rows; row++ {
		for column := 0; column < columns; column++ {
			x := inner.Min.X + (inner.Dx()*(column+1))/(columns+1)
			y := inner.Min.Y + (inner.Dy()*(row+1))/(rows+1)
			candidate := image.Pt(x, y)
			clearance := nearestHazardDistance(candidate, hazards)
			travel := math.Hypot(float64(candidate.X-player.X), float64(candidate.Y-player.Y))
			score := clearance - travel*0.25
			if score > bestScore {
				bestScore = score
				best = candidate
			}
		}
	}
	return best
}

// nearestHazardDirection points a green-heart shield toward the closest
// detected projectile. The green soul is stationary in canonical Undertale
// fights; arrow input rotates the shield instead of moving the soul.
func nearestHazardDirection(player image.Point, hazards []battleHazard) string {
	if len(hazards) == 0 {
		return ""
	}
	nearest := hazards[0]
	best := math.Hypot(float64(player.X-nearest.Center.X), float64(player.Y-nearest.Center.Y))
	for _, hazard := range hazards[1:] {
		distance := math.Hypot(float64(player.X-hazard.Center.X), float64(player.Y-hazard.Center.Y))
		if distance < best {
			nearest = hazard
			best = distance
		}
	}
	return movementDirection(player, nearest.Center)
}
