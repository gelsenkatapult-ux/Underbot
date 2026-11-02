package ai

import (
	"testing"

	"gitlab.com/256/Underbot/cv/object"
)

func TestGameOverScreenIsIdentified(t *testing.T) {
	screen := []object.Object{{
		Recognized: true,
		RecogObj: object.RecognizedObject{
			Type: object.RecMap["gameOverM"],
		},
	}}
	if got := identify(screen); got.Name != "gameOver" {
		t.Fatalf("identify(game-over screen) = %q, want gameOver", got.Name)
	}
}

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
