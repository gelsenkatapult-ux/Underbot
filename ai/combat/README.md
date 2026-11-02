# Undertale combat catalogue

This package is the first data-driven layer for the Underbot combat AI.

## Evidence and scope

The initial catalogue records documented attack *classes* per selected enemy.
It does not claim to contain every attack animation, phase, projectile shape,
frame timing, speed, or safe coordinate. Those require frame-level observations
from gameplay footage or live capture. Unknown measurements remain unset rather
than being guessed.

Primary references:
- https://undertale.wiki/w/Attack_type — contact semantics and enemies using
  white, gray, red, green, light-blue, orange and yellow attacks.
- https://undertale.wiki/w/Combat — encounter structure and SOUL modes.
- https://teemo.dev/game-design/undertale/systems/enemies/ — cross-check of
  attack classes and encounter mechanics.

## Next data collection stages

1. Expand profiles across all encounter entries and route variants.
2. Add one `AttackPattern` per visually distinct sequence, with source video
   timestamps and confidence.
3. Extract representative frames from public gameplay videos and label bullet
   masks, SOUL bounds, board bounds and telegraphs.
4. Calibrate velocity/period from multiple frames; never infer speed from one
   still image.
5. Connect observed projectiles to the predictor only after CV can detect and
   track them reliably.

The current legacy CV pipeline recognizes a small set of sprites and the SOUL;
it does not yet expose a projectile tracker. This package is therefore a tested
foundation, not a claim that autonomous dodging is already integrated.

## True Pacifist route policy

`route.go` adds a conservative, data-only route policy:

- Ordinary encounters and Toriel, Papyrus and Undyne prefer encounter-specific
  ACT/MERCY resolution. The policy does not authorize FIGHT for them.
- Mettaton EX is resolved through the peaceful ratings threshold. The yellow
  phone/shooter mechanic is only appropriate when its phase is recognized.
- Asgore and Omega Flowey are explicit progression exceptions for the required
  Neutral completion; those boss encounters do not award EXP.
- Asriel uses the SAVE/Lost Souls progression, not an ordinary attack strategy.
- Unknown encounter identities return no action and require human review.

True Pacifist is not available until a Neutral ending has been completed on the
save. The route also depends on keeping EXP at zero, completing Papyrus's
hangout, giving water to Undyne immediately after her collapse, completing
Undyne's hangout, reaching Mettaton EX's peaceful ratings ending, delivering
Undyne's letter to Alphys, completing the Alphys date, and progressing through
the True Lab. These are progression facts, not yet an automated quest/flag
tracker in Underbot.

The route policy is not connected to keyboard input yet. It deliberately does
not guess ACT sequences, MERCY timing, ratings, or encounter identity. The
caller must verify the encounter and game phase before acting; uncertain states
must remain idle. A FIGHT timing-bar recognizer is also not implemented, and
the route policy does not need it for ordinary encounters.
