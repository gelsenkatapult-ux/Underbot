# Projectile tracking and short-horizon dodge planning

This package provides a pure-Go first step toward real-time combat movement.

- `FindBattleGeometry` reads the legacy recognizer's `fightBox` and red/blue/green SOUL objects; it rejects missing or out-of-box geometry.
- `ObserveFrame` combines validated geometry, player exclusion, color-connected projectile candidates, and persistent tracking.
- `DetectProjectiles` remains an experimental candidate detector, not a trained sprite recognizer.
- `Tracker.Update` associates detections across frames.
- Tracks estimate horizontal and vertical velocity in pixels per second.
- `ProjectileTrack.Predict` extrapolates a short, constant-velocity trajectory.
- `ChooseDodge` evaluates nine short movement candidates against predicted
  **white damaging attacks** and a player safety margin. It deliberately does
  not treat every colored object as a generic collision hazard: gray is harmless,
  light-blue requires staying still, green is an interaction mechanic, orange
  requires movement, yellow is a shooter mechanic, and red is a warning class.

## Integration contract

The legacy CV pipeline recognizes the fight box and the red, blue, or green
SOUL. `ObserveFrame` now uses those recognized bounds to run the experimental
color-component detector and tracker in `InBattleUpdate`. This is connected
telemetry, not yet a validated attack recognizer: color effects, text, and
particles can still be mistaken for projectiles. Do not use these candidates
for autonomous input until tested against real battle captures.

The planner is intentionally independent of Win32 input and does not press
keys. The combat controller must additionally enforce SOUL mode constraints
(red free movement, blue gravity, green shield, yellow shooter, purple rails),
movement speed, attack-class rules, and game-specific board geometry. The current
planner only handles white attacks; it does not yet implement the required
stay-still response for light-blue, keep-moving response for orange, collection
or shield interaction for green, or shooting for yellow. Those classes are
intentionally not converted into generic dodge movements.

Constant-velocity extrapolation is only a short-horizon approximation. It is
not suitable by itself for homing, curved, accelerating, rotating, or beam
attacks. Such patterns need their own motion model and must be identified from
multiple observed frames before being used for autonomous dodging.
