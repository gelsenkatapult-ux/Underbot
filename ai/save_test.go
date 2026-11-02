package ai

import (
	"image"
	"reflect"
	"testing"

	"gitlab.com/256/Underbot/cv/object"
)

func saveScreenObject(id int, name string, bounds image.Rectangle) object.Object {
	parent := &object.Object{ID: id, Bounds: bounds}
	return object.Object{
		ID:         id,
		Bounds:     bounds,
		Recognized: true,
		RecogObj: object.RecognizedObject{
			Parent: parent,
			Type:   object.RecMap[name],
		},
	}
}

func TestSaveUpdateThrottlesRepeatedConfirmation(t *testing.T) {
	oldCooldown := saveMenuCooldownFrames
	defer func() { saveMenuCooldownFrames = oldCooldown }()

	tests := []struct {
		name       string
		heart      image.Rectangle
		wantInputs []string
	}{
		{
			name:       "save already selected",
			heart:      image.Rect(100, 120, 115, 135),
			wantInputs: []string{"z", "z"},
		},
		{
			name:       "move to save then confirm",
			heart:      image.Rect(500, 120, 515, 135),
			wantInputs: []string{"left", "z", "z"},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			saveMenuCooldownFrames = 0
			window := &pacifistTestWindow{}
			objects := []object.Object{
				saveScreenObject(1, "redHeart", test.heart),
				saveScreenObject(2, "saveBox", image.Rect(100, 100, 513, 263)),
			}

			if err := SaveUpdate(objects, window, image.NewRGBA(image.Rect(0, 0, 640, 480))); err != nil {
				t.Fatal(err)
			}
			if !reflect.DeepEqual(window.keys, test.wantInputs) {
				t.Fatalf("first save-screen inputs = %v, want %v", window.keys, test.wantInputs)
			}
			if saveMenuCooldownFrames <= 0 {
				t.Fatal("save confirmation did not set a cooldown")
			}

			// A further frame must not send another confirmation.
			if err := SaveUpdate(objects, window, image.NewRGBA(image.Rect(0, 0, 640, 480))); err != nil {
				t.Fatal(err)
			}
			if !reflect.DeepEqual(window.keys, test.wantInputs) {
				t.Fatalf("save-screen cooldown allowed duplicate input: %v", window.keys)
			}
		})
	}
}
