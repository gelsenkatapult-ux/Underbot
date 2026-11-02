package combat

import (
	"image"
	"math"
	"sort"
	"time"
)

// ProjectileObservation is one projectile-sized object detected in a frame.
// Bounds must be in the same client-area coordinate system as the player.
type ProjectileObservation struct {
	Bounds image.Rectangle
	Class  AttackClass
}

// ProjectileTrack is a short-lived estimate of one projectile's motion.
type ProjectileTrack struct {
	ID         int
	Bounds     image.Rectangle
	Center     image.Point
	VelocityX  float64 // pixels per second
	VelocityY  float64 // pixels per second
	Class      AttackClass
	LastSeen   time.Time
	Confidence float64
}

// Tracker associates observations between frames and estimates velocity.
// It intentionally does not classify arbitrary CV contours as projectiles;
// callers must provide projectile observations from a validated detector.
type Tracker struct {
	Tracks          []ProjectileTrack
	NextID          int
	MaxAge          time.Duration
	MaxMatchDistance float64
}

func NewTracker() *Tracker {
	return &Tracker{NextID: 1, MaxAge: 250 * time.Millisecond, MaxMatchDistance: 90}
}

// Update matches each observation to at most one prior track. Unmatched
// observations create tracks; old unmatched tracks expire after MaxAge.
func (t *Tracker) Update(observations []ProjectileObservation, now time.Time) []ProjectileTrack {
	if t.MaxAge <= 0 { t.MaxAge = 250 * time.Millisecond }
	if t.MaxMatchDistance <= 0 { t.MaxMatchDistance = 90 }
	used := make([]bool, len(t.Tracks))
	next := make([]ProjectileTrack, 0, len(observations)+len(t.Tracks))

	for _, obs := range observations {
		if obs.Bounds.Empty() { continue }
		center := rectCenter(obs.Bounds)
		best, bestDistance := -1, math.MaxFloat64
		for i, old := range t.Tracks {
			if used[i] { continue }
			d := pointDistance(center, old.Center)
			if d < bestDistance && d <= t.MaxMatchDistance { best, bestDistance = i, d }
		}
		if best >= 0 {
			used[best] = true
			old := t.Tracks[best]
			dt := now.Sub(old.LastSeen).Seconds()
			vx, vy := old.VelocityX, old.VelocityY
			if dt > 0 {
				vx = float64(center.X-old.Center.X) / dt
				vy = float64(center.Y-old.Center.Y) / dt
			}
			next = append(next, ProjectileTrack{
				ID: old.ID, Bounds: obs.Bounds, Center: center,
				VelocityX: vx, VelocityY: vy, Class: obs.Class,
				LastSeen: now, Confidence: math.Min(1, old.Confidence+0.2),
			})
		} else {
			id := t.NextID
			if id <= 0 { id = 1 }
			t.NextID = id + 1
			next = append(next, ProjectileTrack{
				ID: id, Bounds: obs.Bounds, Center: center, Class: obs.Class,
				LastSeen: now, Confidence: 0.5,
			})
		}
	}
	for i, old := range t.Tracks {
		if used[i] { continue }
		if now.Sub(old.LastSeen) <= t.MaxAge {
			old.Confidence *= 0.8
			next = append(next, old)
		}
	}
	sort.Slice(next, func(i, j int) bool { return next[i].ID < next[j].ID })
	t.Tracks = next
	return append([]ProjectileTrack(nil), next...)
}

// Predict returns the projectile rectangle translated by its estimated
// constant velocity after elapsed. This is a local short-horizon model.
func (p ProjectileTrack) Predict(elapsed time.Duration) image.Rectangle {
	dt := elapsed.Seconds()
	dx := int(math.Round(p.VelocityX * dt))
	dy := int(math.Round(p.VelocityY * dt))
	return p.Bounds.Add(image.Pt(dx, dy))
}

type dodgeCandidate struct { dx, dy int }

// ChooseDodge evaluates short movement candidates against predicted projectile
// rectangles. It returns a relative displacement, not a keypress. The caller
// must still respect the battle-box bounds and movement mode.
func ChooseDodge(player, board image.Rectangle, tracks []ProjectileTrack, horizon time.Duration, margin int) (image.Point, bool) {
	if player.Empty() || board.Empty() || horizon <= 0 { return image.Point{}, false }
	candidates := []dodgeCandidate{
		{0, 0}, {0, -24}, {24, 0}, {0, 24}, {-24, 0},
		{17, -17}, {17, 17}, {-17, 17}, {-17, -17},
	}
	bestScore := math.Inf(-1)
	best := image.Point{}
	found := false
	steps := 8
	for _, candidate := range candidates {
		moved := player.Add(image.Pt(candidate.dx, candidate.dy))
		safeBoard := board.Inset(margin)
		if safeBoard.Empty() || !moved.Min.In(safeBoard) || !moved.Max.Sub(image.Pt(1, 1)).In(safeBoard) { continue }
		minDistance := math.Inf(1)
		collision := false
		for step := 0; step <= steps; step++ {
			fraction := float64(step) / float64(steps)
			playerAt := translateRect(player, int(math.Round(float64(candidate.dx)*fraction)), int(math.Round(float64(candidate.dy)*fraction)))
			elapsed := time.Duration(float64(horizon) * fraction)
			for _, projectile := range tracks {
				projectileAt := projectile.Predict(elapsed)
				distance := rectDistance(playerAt, projectileAt, margin)
				if distance == 0 { collision = true; break }
				if distance < minDistance { minDistance = distance }
			}
			if collision { break }
		}
		if collision { continue }
		if len(tracks) == 0 { minDistance = 1 }
		// Prefer the clearest path; break ties by avoiding unnecessary movement.
		score := minDistance - math.Hypot(float64(candidate.dx), float64(candidate.dy))*0.01
		if score > bestScore { bestScore, best, found = score, image.Pt(candidate.dx, candidate.dy), true }
	}
	return best, found
}

func rectCenter(r image.Rectangle) image.Point {
	return image.Pt((r.Min.X+r.Max.X)/2, (r.Min.Y+r.Max.Y)/2)
}

func pointDistance(a, b image.Point) float64 {
	return math.Hypot(float64(a.X-b.X), float64(a.Y-b.Y))
}

func translateRect(r image.Rectangle, dx, dy int) image.Rectangle {
	return r.Add(image.Pt(dx, dy))
}

// rectDistance returns zero for overlap (including the requested safety margin).
func rectDistance(a, b image.Rectangle, margin int) float64 {
	b = b.Inset(-margin)
	dx := 0
	if a.Max.X <= b.Min.X { dx = b.Min.X-a.Max.X } else if b.Max.X <= a.Min.X { dx = a.Min.X-b.Max.X }
	dy := 0
	if a.Max.Y <= b.Min.Y { dy = b.Min.Y-a.Max.Y } else if b.Max.Y <= a.Min.Y { dy = a.Min.Y-b.Max.Y }
	if dx == 0 && dy == 0 { return 0 }
	return math.Hypot(float64(dx), float64(dy))
}
