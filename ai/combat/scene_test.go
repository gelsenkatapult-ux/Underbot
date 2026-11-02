package combat

import (
	"image"
	"testing"

	"gitlab.com/256/Underbot/cv/object"
)

func recognized(name string, bounds image.Rectangle) object.Object {
	parent := &object.Object{Bounds: bounds}
	return object.Object{
		Bounds: bounds,
		Recognized: true,
		RecogObj: object.RecognizedObject{
			Parent: parent,
			Type: object.RecognizableObject{Name: name},
		},
	}
}

func TestFindBattleGeometryRecognizesSoulModes(t *testing.T) {
	cases := []struct {
		name string
		want SoulMode
	}{
		{"redHeart", SoulRed},
		{"blueHeart", SoulBlue},
		{"greenHeart", SoulGreen},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			objects := []object.Object{
				recognized("fightBox", image.Rect(100, 100, 300, 250)),
				recognized(tc.name, image.Rect(140, 140, 155, 155)),
			}
			got, ok := FindBattleGeometry(objects)
			if !ok {
				t.Fatal("expected validated battle geometry")
			}
			if got.Board != image.Rect(100, 100, 300, 250) {
				t.Errorf("board=%v", got.Board)
			}
			if got.Player != image.Rect(140, 140, 155, 155) {
				t.Errorf("player=%v", got.Player)
			}
			if got.SoulMode != tc.want {
				t.Errorf("mode=%q, want %q", got.SoulMode, tc.want)
			}
		})
	}
}

func TestFindBattleGeometryRejectsMissingOrOutOfBoundsSoul(t *testing.T) {
	board := recognized("fightBox", image.Rect(100, 100, 300, 250))
	if _, ok := FindBattleGeometry([]object.Object{board}); ok {
		t.Fatal("geometry without a SOUL must be rejected")
	}
	outside := recognized("redHeart", image.Rect(10, 10, 25, 25))
	if _, ok := FindBattleGeometry([]object.Object{board, outside}); ok {
		t.Fatal("SOUL outside the fight box must be rejected")
	}
}

func TestFindBattleGeometryIgnoresUnrecognizedObjects(t *testing.T) {
	objects := []object.Object{
		recognized("fightBox", image.Rect(100, 100, 300, 250)),
		{Bounds: image.Rect(140, 140, 155, 155), Recognized: false},
	}
	if _, ok := FindBattleGeometry(objects); ok {
		t.Fatal("unrecognized object must not be treated as the SOUL")
	}
}
