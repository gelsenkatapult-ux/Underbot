package combat

import "testing"

func TestTruePacifistDecisionKeepsOrdinaryEncountersNonLethal(t *testing.T) {
	for _, encounter := range []EncounterKind{
		EncounterOrdinary, EncounterToriel, EncounterPapyrus, EncounterUndyne,
		EncounterMettatonEX, EncounterAsriel,
	} {
		decision := TruePacifistDecision(encounter)
		if decision.AllowFight {
			t.Errorf("encounter %q unexpectedly allows FIGHT", encounter)
		}
		if decision.RequiresHumanReview {
			t.Errorf("known encounter %q unexpectedly requires review", encounter)
		}
	}
}

func TestTruePacifistBossProgressionExceptions(t *testing.T) {
	for _, encounter := range []EncounterKind{EncounterAsgore, EncounterOmegaFlowey} {
		decision := TruePacifistDecision(encounter)
		if !decision.AllowFight || decision.Strategy != StrategyBossProgress {
			t.Errorf("encounter %q should permit only the required boss progression, got %+v", encounter, decision)
		}
	}
	if CanUseFightForRoute(RouteTruePacifist, EncounterAsgore) != true {
		t.Fatal("Asgore progression must be permitted")
	}
	if CanUseFightForRoute(RouteTruePacifist, EncounterOrdinary) {
		t.Fatal("ordinary encounters must not permit FIGHT")
	}
	if CanUseFightForRoute(RouteMode("neutral"), EncounterAsgore) {
		t.Fatal("unknown route modes must fail closed")
	}
}

func TestTruePacifistUnknownEncounterFailsClosed(t *testing.T) {
	decision := TruePacifistDecision(EncounterKind("not-recognized"))
	if decision.Strategy != StrategyNoAction || !decision.RequiresHumanReview || decision.AllowFight {
		t.Fatalf("unknown encounter did not fail closed: %+v", decision)
	}
}

func TestTruePacifistSpecialStrategies(t *testing.T) {
	cases := map[EncounterKind]BattleStrategy{
		EncounterMettatonEX: StrategyRatings,
		EncounterAsriel: StrategySaveSouls,
		EncounterToriel: StrategyActThenSpare,
		EncounterPapyrus: StrategyActThenSpare,
		EncounterUndyne: StrategyActThenSpare,
	}
	for encounter, want := range cases {
		if got := TruePacifistDecision(encounter).Strategy; got != want {
			t.Errorf("TruePacifistDecision(%q) strategy=%q, want %q", encounter, got, want)
		}
	}
}
