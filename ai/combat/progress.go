package combat

import "fmt"

// RouteMilestone is a confirmed story event. Callers must only record a
// milestone after the corresponding screen/dialogue has been positively
// recognized; this type does not infer progress from an enemy name.
type RouteMilestone string

const (
	MilestonePapyrusHangout       RouteMilestone = "papyrus-hangout"
	MilestoneUndyneWater          RouteMilestone = "undyne-water"
	MilestoneUndyneHangout        RouteMilestone = "undyne-hangout"
	MilestoneMettatonEXPeaceful   RouteMilestone = "mettaton-ex-peaceful"
	MilestoneNeutralEnding        RouteMilestone = "neutral-ending"
	MilestoneUndyneLetterDelivery RouteMilestone = "undyne-letter-delivered"
	MilestoneAlphysDate            RouteMilestone = "alphys-date"
	MilestoneTrueLabEntered        RouteMilestone = "true-lab-entered"
)

// RouteProgress is a save-local checklist for the True Pacifist route.
// EXP is tracked separately because a single kill permanently invalidates
// the True Pacifist condition for that playthrough.
type RouteProgress struct {
	Experience          int
	PapyrusHangout      bool
	UndyneWater         bool
	UndyneHangout       bool
	MettatonEXPeaceful  bool
	NeutralEnding       bool
	UndyneLetterDelivery bool
	AlphysDate          bool
	TrueLabEntered     bool
}

// AddExperience records EXP awarded by the game. Negative EXP is invalid;
// EXP is deliberately never subtractable in this route tracker.
func (p *RouteProgress) AddExperience(amount int) error {
	if amount < 0 {
		return fmt.Errorf("experience gain cannot be negative: %d", amount)
	}
	p.Experience += amount
	return nil
}

// Record marks a confirmed story milestone while enforcing known prerequisites.
// Re-recording an already completed milestone is harmless and idempotent.
func (p *RouteProgress) Record(milestone RouteMilestone) error {
	switch milestone {
	case MilestonePapyrusHangout:
		p.PapyrusHangout = true
	case MilestoneUndyneWater:
		if !p.PapyrusHangout {
			return fmt.Errorf("cannot record Undyne water before the Papyrus hangout")
		}
		p.UndyneWater = true
	case MilestoneUndyneHangout:
		if !p.PapyrusHangout || !p.UndyneWater {
			return fmt.Errorf("Undyne hangout requires the Papyrus hangout and giving Undyne water")
		}
		p.UndyneHangout = true
	case MilestoneMettatonEXPeaceful:
		if p.Experience != 0 {
			return fmt.Errorf("True Pacifist requires zero EXP; current EXP is %d", p.Experience)
		}
		p.MettatonEXPeaceful = true
	case MilestoneNeutralEnding:
		if p.Experience != 0 {
			return fmt.Errorf("cannot mark a pacifist Neutral ending with %d EXP", p.Experience)
		}
		p.NeutralEnding = true
	case MilestoneUndyneLetterDelivery:
		if !p.NeutralEnding || !p.UndyneHangout || !p.MettatonEXPeaceful {
			return fmt.Errorf("Undyne's letter requires the Neutral ending, Undyne hangout and peaceful Mettaton EX")
		}
		p.UndyneLetterDelivery = true
	case MilestoneAlphysDate:
		if !p.UndyneLetterDelivery {
			return fmt.Errorf("Alphys date requires delivering Undyne's letter")
		}
		p.AlphysDate = true
	case MilestoneTrueLabEntered:
		if !p.ReadyForTrueLab() {
			return fmt.Errorf("True Lab prerequisites are not complete")
		}
		p.TrueLabEntered = true
	default:
		return fmt.Errorf("unknown route milestone %q", milestone)
	}
	return nil
}

// ReadyForTrueLab reports whether every tracked True Pacifist prerequisite is met.
func (p RouteProgress) ReadyForTrueLab() bool {
	return p.Experience == 0 &&
		p.NeutralEnding &&
		p.PapyrusHangout &&
		p.UndyneWater &&
		p.UndyneHangout &&
		p.MettatonEXPeaceful &&
		p.UndyneLetterDelivery &&
		p.AlphysDate
}

// CanUseFightForProgress is the only route-aware FIGHT gate. The two scripted
// progression bosses are permitted only before the first Neutral ending and
// only while EXP remains zero. The stateless helper cannot prove these facts.
func CanUseFightForProgress(mode RouteMode, encounter EncounterKind, progress RouteProgress) bool {
	if mode != RouteTruePacifist || progress.Experience != 0 || progress.NeutralEnding {
		return false
	}
	return encounter == EncounterAsgore || encounter == EncounterOmegaFlowey
}
