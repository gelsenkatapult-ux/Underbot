package ai

import (
	"image"
	"image/color"
	"os"
	"testing"

	"gitlab.com/256/Underbot/cv/object"
)

func TestPacifistOptionIndex(t *testing.T) {
	tests := []struct {
		name             string
		count            int
		mercyHighlighted bool
		want             int
	}{
		{"main menu prefers ACT", 4, false, 1},
		{"main menu chooses highlighted MERCY", 4, true, 3},
		{"short ACT submenu chooses first action", 2, false, 0},
		{"single option", 1, false, 0},
		{"empty options", 0, false, 0},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			if got := pacifistOptionIndex(test.count, test.mercyHighlighted); got != test.want {
				t.Fatalf("pacifistOptionIndex(%d, %t) = %d, want %d", test.count, test.mercyHighlighted, got, test.want)
			}
		})
	}
}

func TestHasYellowMenuText(t *testing.T) {
	img := image.NewRGBA(image.Rect(0, 0, 180, 60))
	button := image.Rect(20, 10, 130, 50)
	if hasYellowMenuText(img, button) {
		t.Fatal("black button unexpectedly detected as highlighted")
	}
	yellow := color.RGBA{R: 255, G: 255, B: 0, A: 255}
	for i := 0; i < 6; i++ {
		img.Set(40+i, 25, yellow)
	}
	if !hasYellowMenuText(img, button) {
		t.Fatal("yellow label pixels were not detected")
	}
}

type pacifistTestWindow struct {
	keys []string
}

func (w *pacifistTestWindow) GetImage() (image.RGBA, error) { return image.RGBA{}, nil }
func (w *pacifistTestWindow) Center() (image.Point, error) { return image.Point{}, nil }
func (w *pacifistTestWindow) Process() (*os.Process, error) { return nil, nil }
func (w *pacifistTestWindow) Name() (string, error) { return "test", nil }
func (w *pacifistTestWindow) Resize(int, int) error { return nil }
func (w *pacifistTestWindow) SetActive() error { return nil }
func (w *pacifistTestWindow) Pause() error { return nil }
func (w *pacifistTestWindow) Resume() error { return nil }
func (w *pacifistTestWindow) Press(key string) error {
	w.keys = append(w.keys, key)
	return nil
}
func (w *pacifistTestWindow) WxH() (int, int, error) { return 640, 480, nil }
func (w *pacifistTestWindow) ID() (int, error) { return 1, nil }

func menuObject(id int, name string, bounds image.Rectangle) object.Object {
	return object.Object{
		ID:         id,
		Bounds:     bounds,
		Recognized: true,
		RecogObj: object.RecognizedObject{
			Parent: &object.Object{ID: id, Bounds: bounds},
			Type:   object.RecognizableObject{Name: name},
		},
	}
}

func TestBattleMenuSelectsACTWhenMercyIsNotHighlighted(t *testing.T) {
	old := PacifistMode
	PacifistMode = true
	defer func() { PacifistMode = old }()

	objects := []object.Object{
		menuObject(1, "battleOption", image.Rect(10, 10, 117, 49)),
		menuObject(2, "battleOption", image.Rect(150, 10, 257, 49)),
		menuObject(3, "battleOption", image.Rect(290, 10, 397, 49)),
		menuObject(4, "battleOption", image.Rect(430, 10, 537, 49)),
		menuObject(5, "redHeart", image.Rect(198, 22, 213, 37)),
	}
	window := &pacifistTestWindow{}
	if err := BattleMenuUpdate(objects, window, image.NewRGBA(image.Rect(0, 0, 640, 480))); err != nil {
		t.Fatal(err)
	}
	if len(window.keys) != 1 || window.keys[0] != "z" {
		t.Fatalf("menu input = %v, want [z] to confirm ACT", window.keys)
	}
}

func TestBattleMenuChoosesMercyWhenYellowAndAligned(t *testing.T) {
	old := PacifistMode
	PacifistMode = true
	defer func() { PacifistMode = old }()

	objects := []object.Object{
		menuObject(1, "battleOption", image.Rect(10, 10, 117, 49)),
		menuObject(2, "battleOption", image.Rect(150, 10, 257, 49)),
		menuObject(3, "battleOption", image.Rect(290, 10, 397, 49)),
		menuObject(4, "battleOption", image.Rect(430, 10, 537, 49)),
		menuObject(5, "redHeart", image.Rect(475, 20, 490, 35)),
	}
	img := image.NewRGBA(image.Rect(0, 0, 640, 480))
	for i := 0; i < 6; i++ {
		img.Set(470+i, 28, color.RGBA{R: 255, G: 255, A: 255})
	}
	window := &pacifistTestWindow{}
	if err := BattleMenuUpdate(objects, window, img); err != nil {
		t.Fatal(err)
	}
	if len(window.keys) != 1 || window.keys[0] != "z" {
		t.Fatalf("menu input = %v, want [z] to confirm highlighted MERCY", window.keys)
	}
}

func TestAttackGoalNeverConfirmsInPacifistMode(t *testing.T) {
	old := PacifistMode
	PacifistMode = true
	defer func() { PacifistMode = old }()
	window := &pacifistTestWindow{}
	if err := AttackGoalUpdate(nil, window, nil); err != nil {
		t.Fatal(err)
	}
	if len(window.keys) != 0 {
		t.Fatalf("attack timing emitted inputs %v in pacifist mode", window.keys)
	}
}

func TestUnstuckOnlyUsesDirectionalInputs(t *testing.T) {
	window := &pacifistTestWindow{}
	for i := 0; i < 8; i++ {
		if err := unstuck(window); err != nil {
			t.Fatal(err)
		}
	}
	allowed := map[string]bool{"up": true, "right": true, "down": true, "left": true}
	for _, key := range window.keys {
		if !allowed[key] {
			t.Fatalf("unstuck issued unsafe confirmation key %q", key)
		}
	}
}
