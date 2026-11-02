package ai

import (
	"testing"

	"gitlab.com/256/Underbot/cv/object"
)

func recognized(name string) object.Object {
	return object.Object{RecogObj: object.RecognizedObject{Type: object.RecMap[name]}}
}

func TestIdentifyUnknownWithoutRecognizedObjects(t *testing.T) {
	if got := identify(nil); got.Name != "Unknown" { t.Fatalf("identify(nil) = %q, want Unknown", got.Name) }
}
func TestIdentifyOutsideBattleFromFrisk(t *testing.T) {
	got := identify([]object.Object{recognized("friskFrontFace")})
	if got.Name != "outsideBattle" { t.Fatalf("identify(Frisk) = %q, want outsideBattle", got.Name) }
}
func TestIdentifyBattleMenu(t *testing.T) {
	got := identify([]object.Object{recognized("narratorBox"), recognized("redHeart")})
	if got.Name != "battleMenu" { t.Fatalf("identify(battle menu) = %q, want battleMenu", got.Name) }
}
func TestIdentifyInBattle(t *testing.T) {
	got := identify([]object.Object{recognized("fightBox"), recognized("redHeart")})
	if got.Name != "inBattle" { t.Fatalf("identify(in battle) = %q, want inBattle", got.Name) }
}
func TestIdentifyDialogueAndAntiSigns(t *testing.T) {
	got := identify([]object.Object{recognized("dialogueBox")})
	if got.Name != "dialogue" { t.Fatalf("identify(dialogue) = %q, want dialogue", got.Name) }
	got = identify([]object.Object{recognized("dialogueBox"), recognized("attackGoal")})
	if got.Name == "dialogue" { t.Fatal("dialogue state matched despite an anti-sign") }
}
func TestGetWantedIsIndependentOfInputObjectOrder(t *testing.T) {
	wanted := []object.RecognizableObject{object.RecMap["redHeart"], object.RecMap["fightBox"]}
	found, err := GetWanted([]object.Object{recognized("fightBox"), recognized("redHeart")}, wanted)
	if err != nil { t.Fatalf("GetWanted returned an error: %v", err) }
	if len(found) != 2 || found[0].Type.Name != "redHeart" || found[1].Type.Name != "fightBox" {
		t.Fatalf("GetWanted returned unexpected order/content: %#v", found)
	}
}
func TestGetWantedReportsMissingTypes(t *testing.T) {
	wanted := []object.RecognizableObject{object.RecMap["redHeart"], object.RecMap["fightBox"]}
	found, err := GetWanted([]object.Object{recognized("redHeart")}, wanted)
	if err == nil { t.Fatal("GetWanted should report a missing fightBox") }
	if len(found) != 1 || found[0].Type.Name != "redHeart" { t.Fatalf("unexpected partial match: %#v", found) }
}
