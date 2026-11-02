package combat

import "testing"

func TestClassRules(t *testing.T) {
	cases := []struct{ class AttackClass; want AttackRule }{
		{AttackWhite, RuleDodge}, {AttackGray, RuleHarmless},
		{AttackRed, RuleWarning}, {AttackGreen, RuleCollect},
		{AttackLightBlue, RuleStayStill}, {AttackOrange, RuleKeepMoving},
		{AttackYellow, RuleShoot},
	}
	for _, tc := range cases {
		if got := ClassRule(tc.class); got != tc.want {
			t.Errorf("ClassRule(%q) = %q, want %q", tc.class, got, tc.want)
		}
	}
}

func TestEveryProfileHasEvidenceAndClasses(t *testing.T) {
	if len(Profiles) < 20 { t.Fatalf("catalogue unexpectedly small: %d profiles", len(Profiles)) }
	seen := map[string]bool{}
	for _, p := range Profiles {
		if p.Name == "" || p.Area == "" || p.SourceURL == "" || len(p.AttackClasses) == 0 {
			t.Errorf("incomplete evidence-labelled profile: %#v", p)
		}
		if seen[p.Name] { t.Errorf("duplicate profile %q", p.Name) }
		seen[p.Name] = true
	}
}

func TestFindProfile(t *testing.T) {
	p, ok := FindProfile("sAnS")
	if !ok || p.Name != "Sans" { t.Fatalf("FindProfile did not find Sans: %#v, %v", p, ok) }
	if _, ok := FindProfile("not an enemy"); ok { t.Fatal("unexpected profile match") }
}
