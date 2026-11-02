package combat

import "testing"

func TestAttackIntentUsesAttackClassAndSoulMode(t *testing.T) {
	cases := []struct {
		mode SoulMode
		class AttackClass
		want CombatIntent
	}{
		{SoulRed, AttackWhite, IntentDodge},
		{SoulBlue, AttackWhite, IntentDodge},
		{SoulGreen, AttackWhite, IntentDodge},
		{SoulRed, AttackLightBlue, IntentStayStill},
		{SoulRed, AttackOrange, IntentKeepMoving},
		{SoulRed, AttackGreen, IntentCollect},
		{SoulGreen, AttackGreen, IntentShield},
		{SoulYellow, AttackYellow, IntentShoot},
		{SoulRed, AttackYellow, IntentWarning},
		{SoulRed, AttackGray, IntentNoAction},
		{SoulRed, AttackClass("unmapped"), IntentNoAction},
	}
	for _, tc := range cases {
		if got := AttackIntent(tc.mode, tc.class); got != tc.want {
			t.Errorf("AttackIntent(%q, %q)=%q, want %q", tc.mode, tc.class, got, tc.want)
		}
	}
}

func TestUnknownAttackClassIsNotACommand(t *testing.T) {
	if got := ClassRule(AttackClass("unknown-color")); got != RuleUnknown {
		t.Fatalf("unknown class rule=%q, want %q", got, RuleUnknown)
	}
	if got := AttackIntent(SoulRed, AttackClass("unknown-color")); got != IntentNoAction {
		t.Fatalf("unknown attack intent=%q, want no action", got)
	}
}

func TestModeCanMove(t *testing.T) {
	if !ModeCanMove(SoulRed) {
		t.Fatal("red SOUL should allow free movement")
	}
	for _, mode := range []SoulMode{SoulBlue, SoulGreen, SoulYellow, SoulPurple, SoulMode("unknown")} {
		if ModeCanMove(mode) {
			t.Errorf("mode %q should not be treated as free movement", mode)
		}
	}
}
