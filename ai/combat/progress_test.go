package combat

import "testing"

func TestRouteProgressRejectsOutOfOrderMilestones(t *testing.T) {
	var p RouteProgress
	if err := p.Record(MilestoneUndyneWater); err == nil {
		t.Fatal("Undyne water must require Papyrus hangout")
	}
	if err := p.Record(MilestoneUndyneHangout); err == nil {
		t.Fatal("Undyne hangout must require water and Papyrus hangout")
	}
	if err := p.Record(MilestoneUndyneLetterDelivery); err == nil {
		t.Fatal("letter delivery must require the Neutral ending and prior friendships")
	}
	if err := p.Record(MilestoneAlphysDate); err == nil {
		t.Fatal("Alphys date must require letter delivery")
	}
	if err := p.Record(MilestoneTrueLabEntered); err == nil {
		t.Fatal("True Lab must require all tracked prerequisites")
	}
}

func TestRouteProgressRecordsTruePacifistSequence(t *testing.T) {
	var p RouteProgress
	sequence := []RouteMilestone{
		MilestonePapyrusHangout,
		MilestoneUndyneWater,
		MilestoneUndyneHangout,
		MilestoneMettatonEXPeaceful,
		MilestoneNeutralEnding,
		MilestoneUndyneLetterDelivery,
		MilestoneAlphysDate,
		MilestoneTrueLabEntered,
	}
	for _, milestone := range sequence {
		if err := p.Record(milestone); err != nil {
			t.Fatalf("Record(%q): %v", milestone, err)
		}
	}
	if !p.ReadyForTrueLab() || !p.TrueLabEntered {
		t.Fatalf("expected completed route prerequisites, got %+v", p)
	}
}

func TestRouteProgressRejectsExperienceAndCannotEraseIt(t *testing.T) {
	var p RouteProgress
	if err := p.AddExperience(-1); err == nil {
		t.Fatal("negative EXP must be rejected")
	}
	if err := p.AddExperience(1); err != nil {
		t.Fatal(err)
	}
	if err := p.AddExperience(-1); err == nil {
		t.Fatal("EXP must not be erasable")
	}
	for _, milestone := range []RouteMilestone{MilestoneMettatonEXPeaceful, MilestoneNeutralEnding} {
		if err := p.Record(milestone); err == nil {
			t.Errorf("milestone %q must reject non-zero EXP", milestone)
		}
	}
	if p.ReadyForTrueLab() {
		t.Fatal("non-zero EXP must block True Pacifist readiness")
	}
}

func TestRouteProgressUnknownMilestoneFailsClosed(t *testing.T) {
	var p RouteProgress
	if err := p.Record(RouteMilestone("guessed-from-screen")); err == nil {
		t.Fatal("unknown milestone must be rejected")
	}
}

func TestRouteAwareFightGateRequiresProgressEvidence(t *testing.T) {
	var p RouteProgress
	for _, encounter := range []EncounterKind{EncounterAsgore, EncounterOmegaFlowey} {
		if !CanUseFightForProgress(RouteTruePacifist, encounter, p) {
			t.Errorf("pre-Neutral progression boss %q should be permitted at zero EXP", encounter)
		}
	}
	if CanUseFightForProgress(RouteTruePacifist, EncounterOrdinary, p) {
		t.Fatal("ordinary monster must never be authorized for FIGHT")
	}
	p.NeutralEnding = true
	if CanUseFightForProgress(RouteTruePacifist, EncounterAsgore, p) {
		t.Fatal("progression FIGHT must be blocked after Neutral ending")
	}
	p.NeutralEnding = false
	p.Experience = 1
	if CanUseFightForProgress(RouteTruePacifist, EncounterAsgore, p) {
		t.Fatal("progression FIGHT must be blocked when EXP is non-zero")
	}
}
