package combat

import (
	"image"
	"testing"
	"time"
)

func TestTrackerEstimatesVelocity(t *testing.T) {
	tracker := NewTracker()
	start := time.Unix(100, 0)
	tracker.Update([]ProjectileObservation{{Bounds: image.Rect(10, 10, 20, 20), Class: AttackWhite}}, start)
	got := tracker.Update([]ProjectileObservation{{Bounds: image.Rect(20, 10, 30, 20), Class: AttackWhite}}, start.Add(100*time.Millisecond))
	if len(got) != 1 { t.Fatalf("got %d tracks, want 1", len(got)) }
	if got[0].VelocityX != 100 || got[0].VelocityY != 0 { t.Fatalf("velocity=(%v,%v), want (100,0)", got[0].VelocityX, got[0].VelocityY) }
	if got[0].ID != 1 { t.Fatalf("track ID=%d, want stable ID 1", got[0].ID) }
}

func TestTrackerExpiresLostProjectiles(t *testing.T) {
	tracker := NewTracker()
	start := time.Unix(200, 0)
	tracker.Update([]ProjectileObservation{{Bounds: image.Rect(1, 1, 5, 5)}}, start)
	got := tracker.Update(nil, start.Add(time.Second))
	if len(got) != 0 { t.Fatalf("stale track was retained: %#v", got) }
}

func TestTrackerDoesNotDuplicateOneObservation(t *testing.T) {
	tracker := NewTracker()
	now := time.Unix(300, 0)
	tracker.Update([]ProjectileObservation{
		{Bounds: image.Rect(10, 10, 20, 20)},
		{Bounds: image.Rect(100, 10, 110, 20)},
	}, now)
	got := tracker.Update([]ProjectileObservation{{Bounds: image.Rect(12, 10, 22, 20)}}, now.Add(50*time.Millisecond))
	if len(got) != 2 { t.Fatalf("expected one matched and one briefly retained track, got %d", len(got)) }
}

func TestPredictMovesProjectile(t *testing.T) {
	p := ProjectileTrack{Bounds: image.Rect(0, 0, 10, 10), VelocityX: 50, VelocityY: -20}
	got := p.Predict(200 * time.Millisecond)
	if got != image.Rect(10, -4, 20, 6) { t.Fatalf("prediction=%v", got) }
}

func TestChooseDodgeFindsSafeMovement(t *testing.T) {
	player := image.Rect(45, 45, 55, 55)
	board := image.Rect(0, 0, 100, 100)
	tracks := []ProjectileTrack{{
		Bounds: image.Rect(45, 10, 55, 20), VelocityY: 100,
	}}
	move, ok := ChooseDodge(player, board, tracks, 400*time.Millisecond, 2)
	if !ok { t.Fatal("expected at least one safe candidate") }
	if move == (image.Point{}) { t.Fatal("stationary candidate should be unsafe") }
}

func TestChooseDodgeRejectsInvalidGeometry(t *testing.T) {
	if _, ok := ChooseDodge(image.Rectangle{}, image.Rect(0, 0, 100, 100), nil, time.Second, 0); ok {
		t.Fatal("empty player rectangle must be rejected")
	}
}

func TestChooseDodgeDoesNotTreatOtherAttackClassesAsWhiteHazards(t *testing.T) {
	player := image.Rect(45, 45, 55, 55)
	board := image.Rect(0, 0, 100, 100)
	classes := []AttackClass{AttackGray, AttackRed, AttackGreen, AttackLightBlue, AttackOrange, AttackYellow}
	for _, class := range classes {
		tracks := []ProjectileTrack{{
			Bounds: image.Rect(45, 10, 55, 20), VelocityY: 100, Class: class,
		}}
		move, ok := ChooseDodge(player, board, tracks, 400*time.Millisecond, 2)
		if !ok { t.Fatalf("class %q: expected a valid conservative candidate", class) }
		if move != (image.Point{}) {
			t.Errorf("class %q caused generic dodge %v; class-specific behavior is required", class, move)
		}
	}
}

func TestChooseDodgeStillAvoidsWhiteAttacks(t *testing.T) {
	player := image.Rect(45, 45, 55, 55)
	board := image.Rect(0, 0, 100, 100)
	tracks := []ProjectileTrack{{
		Bounds: image.Rect(45, 10, 55, 20), VelocityY: 100, Class: AttackWhite,
	}}
	move, ok := ChooseDodge(player, board, tracks, 400*time.Millisecond, 2)
	if !ok { t.Fatal("expected a safe candidate") }
	if move == (image.Point{}) { t.Fatal("white damaging attack should trigger an evasive candidate") }
}
