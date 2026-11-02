// Package combat contains the data-driven combat model for Undertale encounters.
//
// The initial catalogue records documented attack semantics, not fabricated
// frame timings or trajectories. Runtime observations must supply those values.
package combat

// AttackClass describes the contact rule of an Undertale attack.
type AttackClass string

const (
	AttackWhite    AttackClass = "white"
	AttackGray     AttackClass = "gray"
	AttackRed      AttackClass = "red"
	AttackGreen    AttackClass = "green"
	AttackLightBlue AttackClass = "light-blue"
	AttackOrange   AttackClass = "orange"
	AttackYellow   AttackClass = "yellow"
)

// SoulMode describes the movement constraints imposed on the SOUL.
type SoulMode string

const (
	SoulRed    SoulMode = "red-free"
	SoulBlue   SoulMode = "blue-gravity"
	SoulGreen  SoulMode = "green-shield"
	SoulYellow SoulMode = "yellow-shooter"
	SoulPurple SoulMode = "purple-rail"
)

// Trajectory is deliberately descriptive. Exact speed and timing are learned
// from observed frames rather than guessed from a still image or wiki text.
type Trajectory string

const (
	TrajectoryUnknown  Trajectory = "unknown"
	TrajectoryStationary Trajectory = "stationary"
	TrajectoryLinear   Trajectory = "linear"
	TrajectoryHoming    Trajectory = "homing"
	TrajectoryWave      Trajectory = "wave"
	TrajectoryRotating  Trajectory = "rotating"
	TrajectoryBeam      Trajectory = "beam"
	TrajectoryFalling   Trajectory = "falling"
)

// AttackRule explains what the player should do when touching an attack.
type AttackRule string

const (
	RuleDodge       AttackRule = "dodge"
	RuleStayStill   AttackRule = "stay-still"
	RuleKeepMoving  AttackRule = "keep-moving"
	RuleCollect     AttackRule = "collect-or-interact"
	RuleShield      AttackRule = "block-with-shield"
	RuleShoot       AttackRule = "shoot-or-trigger"
	RuleWarning     AttackRule = "warning-only"
	RuleHarmless    AttackRule = "harmless"
)

// ClassRule returns the documented contact rule for an attack class.
func ClassRule(class AttackClass) AttackRule {
	switch class {
	case AttackWhite:
		return RuleDodge
	case AttackGray:
		return RuleHarmless
	case AttackRed:
		return RuleWarning
	case AttackGreen:
		return RuleCollect
	case AttackLightBlue:
		return RuleStayStill
	case AttackOrange:
		return RuleKeepMoving
	case AttackYellow:
		return RuleShoot
	default:
		return RuleDodge
	}
}

// AttackPattern is a reusable pattern description. Zero-valued measurements
// mean "not yet calibrated"; consumers must not interpret them as zero speed.
type AttackPattern struct {
	ID             string
	Name           string
	Classes        []AttackClass
	Trajectory     Trajectory
	Phase          string
	SoulMode       SoulMode
	SpeedPixelsSec float64
	PeriodMillis   int
	LeadMillis     int
	Notes          string
	SourceURL      string
}

// EnemyProfile contains only facts supported by the cited reference.
// Patterns can be added incrementally as frame-level observations are verified.
type EnemyProfile struct {
	Name         string
	Area         string
	AttackClasses []AttackClass
	Notes        string
	SourceURL    string
}

// Profiles is the initial, evidence-labelled enemy catalogue. It intentionally
// omits undocumented timing, geometry and per-phase attack sequences.
var Profiles = []EnemyProfile{
	{Name: "Doggo", Area: "Snowdin", AttackClasses: []AttackClass{AttackWhite, AttackLightBlue}, Notes: "Introduces light-blue attacks.", SourceURL: "https://undertale.wiki/w/Attack_type"},
	{Name: "Lesser Dog", Area: "Snowdin", AttackClasses: []AttackClass{AttackWhite, AttackLightBlue}, SourceURL: "https://undertale.wiki/w/Attack_type"},
	{Name: "Greater Dog", Area: "Snowdin", AttackClasses: []AttackClass{AttackWhite, AttackLightBlue}, SourceURL: "https://undertale.wiki/w/Attack_type"},
	{Name: "Dogamy", Area: "Snowdin", AttackClasses: []AttackClass{AttackWhite, AttackLightBlue}, SourceURL: "https://undertale.wiki/w/Attack_type"},
	{Name: "Dogaressa", Area: "Snowdin", AttackClasses: []AttackClass{AttackWhite, AttackLightBlue}, SourceURL: "https://undertale.wiki/w/Attack_type"},
	{Name: "Gyftrot", Area: "Snowdin", AttackClasses: []AttackClass{AttackWhite, AttackLightBlue}, SourceURL: "https://undertale.wiki/w/Attack_type"},
	{Name: "Papyrus", Area: "Snowdin", AttackClasses: []AttackClass{AttackWhite, AttackClass("light-blue")}, Notes: "Blue attacks interact with the blue SOUL gravity mode.", SourceURL: "https://undertale.wiki/w/Attack_type"},
	{Name: "Vegetoid", Area: "Ruins", AttackClasses: []AttackClass{AttackWhite, AttackGreen}, Notes: "Green attacks can heal and/or advance the encounter.", SourceURL: "https://undertale.wiki/w/Attack_type"},
	{Name: "Napstablook", Area: "Ruins", AttackClasses: []AttackClass{AttackWhite, AttackGray}, Notes: "Gray attack used for decorative text.", SourceURL: "https://undertale.wiki/w/Attack_type"},
	{Name: "Ice Cap", Area: "Snowdin", AttackClasses: []AttackClass{AttackWhite, AttackGray, AttackRed}, Notes: "Red warning attacks; one documented attack also uses gray.", SourceURL: "https://undertale.wiki/w/Attack_type"},
	{Name: "Woshua", Area: "Waterfall", AttackClasses: []AttackClass{AttackWhite, AttackGreen, AttackLightBlue}, SourceURL: "https://undertale.wiki/w/Attack_type"},
	{Name: "Whimsalot", Area: "Hotland", AttackClasses: []AttackClass{AttackWhite, AttackGreen}, SourceURL: "https://undertale.wiki/w/Attack_type"},
	{Name: "Parsnik", Area: "Hotland", AttackClasses: []AttackClass{AttackWhite, AttackGreen}, SourceURL: "https://undertale.wiki/w/Attack_type"},
	{Name: "Tsunderplane", Area: "Hotland", AttackClasses: []AttackClass{AttackWhite, AttackGreen}, Notes: "Green attacks are not necessarily healing attacks.", SourceURL: "https://undertale.wiki/w/Attack_type"},
	{Name: "Royal Guards", Area: "Hotland", AttackClasses: []AttackClass{AttackWhite, AttackGreen}, Notes: "Green attacks are encounter mechanics, not guaranteed healing.", SourceURL: "https://undertale.wiki/w/Attack_type"},
	{Name: "Pyrope", Area: "Hotland", AttackClasses: []AttackClass{AttackWhite, AttackOrange}, SourceURL: "https://undertale.wiki/w/Attack_type"},
	{Name: "Mettaton", Area: "Hotland", AttackClasses: []AttackClass{AttackWhite, AttackLightBlue}, SourceURL: "https://undertale.wiki/w/Attack_type"},
	{Name: "Mettaton EX", Area: "Hotland", AttackClasses: []AttackClass{AttackWhite, AttackGray, AttackRed, AttackLightBlue, AttackYellow}, Notes: "Yellow shooter mechanic and red telegraphs are encounter-specific.", SourceURL: "https://undertale.wiki/w/Attack_type"},
	{Name: "Asgore", Area: "New Home", AttackClasses: []AttackClass{AttackWhite, AttackRed, AttackGreen, AttackLightBlue, AttackOrange}, SourceURL: "https://undertale.wiki/w/Attack_type"},
	{Name: "Undyne", Area: "Waterfall", AttackClasses: []AttackClass{AttackWhite, AttackRed}, Notes: "Multicolored attacks behave as damaging white attacks.", SourceURL: "https://undertale.wiki/w/Attack_type"},
	{Name: "Sans", Area: "Judgment Hall", AttackClasses: []AttackClass{AttackWhite, AttackRed, AttackLightBlue}, Notes: "Nearly all attacks apply KARMA; exact attack sequence/timing is not yet encoded.", SourceURL: "https://undertale.wiki/w/Attack_type"},
	{Name: "So Sorry", Area: "Hotland", AttackClasses: []AttackClass{AttackWhite, AttackLightBlue, AttackOrange}, SourceURL: "https://undertale.wiki/w/Attack_type"},
	{Name: "Omega Flowey", Area: "True Lab", AttackClasses: []AttackClass{AttackWhite, AttackRed, AttackGreen}, Notes: "Multicolored attacks behave as damaging white attacks.", SourceURL: "https://undertale.wiki/w/Attack_type"},
	{Name: "Asriel", Area: "True Pacifist", AttackClasses: []AttackClass{AttackWhite, AttackRed}, Notes: "Multicolored attacks behave as damaging white attacks.", SourceURL: "https://undertale.wiki/w/Attack_type"},
}

// FindProfile returns a profile by its case-insensitive display name.
func FindProfile(name string) (EnemyProfile, bool) {
	for _, profile := range Profiles {
		if equalFold(profile.Name, name) {
			return profile, true
		}
	}
	return EnemyProfile{}, false
}

func equalFold(a, b string) bool {
	if len(a) != len(b) { return false }
	for i := range a {
		x, y := a[i], b[i]
		if x >= 'A' && x <= 'Z' { x += 'a' - 'A' }
		if y >= 'A' && y <= 'Z' { y += 'a' - 'A' }
		if x != y { return false }
	}
	return true
}
