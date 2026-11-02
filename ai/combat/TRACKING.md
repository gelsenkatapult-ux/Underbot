# Projectile tracking and short-horizon dodge planning

This package provides a pure-Go first step toward real-time combat movement.

- `DetectProjectiles` extracts conservative color-connected candidates inside a supplied battle-box rectangle and supports player exclusion rectangles. This is an experimental candidate detector, not a trained sprite recognizer.
- `Tracker.Update` associates detections across frames.
- Tracks estimate horizontal and vertical velocity in pixels per second.
- `ProjectileTrack.Predict` extrapolates a short, constant-velocity trajectory.
- `ChooseDodge` evaluates nine short movement candidates against predicted
  projectile rectangles and a player safety margin.

## Integration contract

The current legacy CV pipeline does **not** yet identify projectiles. A future
detector must produce `ProjectileObservation` values from validated attack
sprites/masks, using the same client-area coordinates as the player and fight
box. Do not feed every contour into this tracker as if it were a projectile.

The planner is intentionally independent of Win32 input and does not press
keys. The combat controller must additionally enforce SOUL mode constraints
(red free movement, blue gravity, green shield, yellow shooter, purple rails),
movement speed, attack-class rules, and game-specific board geometry.

Constant-velocity extrapolation is only a short-horizon approximation. It is
not suitable by itself for homing, curved, accelerating, rotating, or beam
attacks. Such patterns need their own motion model and must be identified from
multiple observed frames before being used for autonomous dodging.
