package ai

import "testing"

func TestGameOverUpdateDisablesAI(t *testing.T) {
	old := Disabled
	Disabled = false
	defer func() { Disabled = old }()

	if err := GameOverUpdate(nil, nil, nil); err != nil {
		t.Fatalf("GameOverUpdate returned error: %v", err)
	}
	if !Disabled {
		t.Fatal("GameOverUpdate did not disable the AI")
	}
}
