package ai

import (
	"fmt"
	"image"
	"image/color"
	"time"

	"github.com/pkg/errors"
	"gitlab.com/256/Underbot/ai/pathfinding"\n\t"gitlab.com/256/Underbot/ai/combat"
	"gitlab.com/256/Underbot/cv/object"
	"gitlab.com/256/Underbot/cv/params"
	"gitlab.com/256/Underbot/cv/rect"
	"gitlab.com/256/Underbot/sys"
)

// GridShow determines if the pathfinding grid is enabled
var GridShow bool\n\n// combatTracker keeps projectile identities stable between consecutive battle frames.\nvar combatTracker = combat.NewTracker()

// Should be ran in every update function (other than unknown)
func genUpdate() {
	unknownFrames = 0
}

// BattleMenuUpdate is the function run every frame when the battle menu is detected
func BattleMenuUpdate(objects []object.Object, win sys.Window, img *image.RGBA) error {
	genUpdate()
	recObjects, err := GetWanted(objects, []object.RecognizableObject{
		object.RecMap["battleOption"],
		object.RecMap["redHeart"],
	})
	if err != nil {
		failedRetrieval++
		if failedRetrieval > params.FailedLimit {
			return errors.Wrap(err, "battle menu objects were not detected")
		}
		return nil
	}
	failedRetrieval = 0
	recMap := Map(recObjects)
	if len(recMap["battleOption"]) == 0 || len(recMap["redHeart"]) == 0 {
		return nil
	}

	var farthest int
	var fightOption object.RecognizedObject
	for i, option := range recMap["battleOption"] {
		if i == 0 {
			farthest = option.Parent.ID
			fightOption = option
			continue
		}
		if farthest < option.Parent.ID {
			farthest = option.Parent.ID
			fightOption = option
		}
	}

	redHeart := rect.RectangleCenter(recMap["redHeart"][0].Parent.Bounds)
	fightCenter := rect.RectangleCenter(fightOption.Parent.Bounds)

	if redHeart.X > fightCenter.X {
		err := win.Press("left")
		if err != nil {
			return errors.Wrap(err, "couldn't press left arrow key")
		}
	} else {
		err := win.Press("z")
		if err != nil {
			return errors.Wrap(err, "couldn't press z key")
		}
	}

	return nil
}

// EmptyUpdate does nothing
func EmptyUpdate(objects []object.Object, win sys.Window, img *image.RGBA) error {
	genUpdate()
	return nil
}

// DialogueUpdate is the function run every frame when dialogue is detected
func DialogueUpdate(objects []object.Object, win sys.Window, img *image.RGBA) error {
	genUpdate()
	err := win.Press("x")
	if err != nil {
		return errors.Wrap(err, "failed to press x key")
	}
	err = win.Press("z")
	if err != nil {
		return errors.Wrap(err, "failed to press z key")
	}
	return nil
}

// InBattleUpdate is the function run every frame when dialogue is detected
func InBattleUpdate(objects []object.Object, win sys.Window, img *image.RGBA) error {
	genUpdate()
	heartObjects, _ := GetWanted(objects, append([]object.RecognizableObject{}, object.Hearts...))
	if len(heartObjects) == 0 {
		// Stall until items can be found
		failedRetrieval++
		// If stalling takes too long, then try to get unstuck
		if failedRetrieval == params.FailedLimit+1 {
			fmt.Println("The AI cannot identify the player or battle state; waiting for a clear frame.")
		}
		return nil
	}
	failedRetrieval = 0

	// Build combat observations only when both the fight box and SOUL are
	// recognized. This stage is telemetry-only: no combat keypresses are sent.
	observation, observed := combat.ObserveFrame(img, objects, combatTracker)
	if observed && GridShow {
		for _, projectile := range observation.Projectiles {
			var marker color.RGBA
			switch projectile.Class {
			case combat.AttackWhite:
				marker = color.RGBA{255, 255, 255, 255}
			case combat.AttackGray:
				marker = color.RGBA{150, 150, 150, 255}
			case combat.AttackRed:
				marker = color.RGBA{255, 0, 0, 255}
			case combat.AttackGreen:
				marker = color.RGBA{0, 220, 0, 255}
			case combat.AttackLightBlue:
				marker = color.RGBA{80, 190, 255, 255}
			case combat.AttackOrange:
				marker = color.RGBA{255, 140, 0, 255}
			case combat.AttackYellow:
				marker = color.RGBA{255, 230, 0, 255}
			}
			rect.DrawRectangle(img, marker, projectile.Bounds)
		}
	}

	// Draw the tiles on the screen
	tiles, err := pathfinding.MakeTiles(*img)
	if err != nil {
		return errors.Wrap(err, "failed to create the screen tiles")
	}

	// Calculate which tile Frisk is in

	//currentTile := pathfinding.GetCurrentTile(image.Point{averageX, averageY})
	//_, err = pathfinding.GetPath(currentTile, pathfinding.GetGoal())
	if err != nil {
		return errors.Wrap(err, "failed to generate path")
	}
	/*
		for _, pathTile := range path {

		}
	*/
	if GridShow {
		for _, tile := range tiles {
			rect.DrawRectangle(img, tile.Color, tile.Rectangle)
		}
	}
	return nil
}

var failedRetrieval int

// SaveUpdate is the function that presses the save button
func SaveUpdate(objects []object.Object, win sys.Window, img *image.RGBA) error {
	genUpdate()
	recObjects, err := GetWanted(objects, []object.RecognizableObject{
		object.RecMap["redHeart"],
		object.RecMap["saveBox"],
	})
	if err != nil {
		failedRetrieval++
		if failedRetrieval == params.FailedLimit+1 {
			fmt.Println("The AI cannot identify the save screen; waiting for a clear frame.")
		}
		return nil
	}
	failedRetrieval = 0
	recMap := Map(recObjects)
	if len(recMap["redHeart"]) == 0 || len(recMap["saveBox"]) == 0 {
		return nil
	}

	// Determines if the heart is on the left side, indicating that it is selecting save
	if rect.RectangleCenter(recMap["redHeart"][0].Parent.Bounds).X < rect.RectangleCenter(recMap["saveBox"][0].Parent.Bounds).X {
		err := win.Press("z")
		if err != nil {
			return errors.Wrap(err, "failed to press z key")
		}
		err = win.Press("z")
		if err != nil {
			return errors.Wrap(err, "failed to press z key")
		}
	} else {
		err := win.Press("left")
		if err != nil {
			return errors.Wrap(err, "failed to press left key")
		}
		err = win.Press("z")
		if err != nil {
			return errors.Wrap(err, "failed to press z key")
		}
		err = win.Press("z")
		if err != nil {
			return errors.Wrap(err, "failed to press z key")
		}
	}
	return nil
}

// OutsideBattleUpdate is the function that pathfinds through the game
func OutsideBattleUpdate(objects []object.Object, win sys.Window, img *image.RGBA) error {
	genUpdate()
	friskObjects, _ := GetWanted(objects, object.Frisk)
	if len(friskObjects) == 0 {
		// Stall until items can be found
		failedRetrieval++
		// If stalling takes too long, then try to get unstuck
		if failedRetrieval == params.FailedLimit+1 {
			fmt.Println("The AI cannot identify the player or battle state; waiting for a clear frame.")
		}
		return nil
	}
	failedRetrieval = 0

	// The point in the middle of all the frisk objects
	var averageX int
	var averageY int
	for _, obj := range friskObjects {
		point := rect.RectangleCenter(obj.Parent.Bounds)
		averageX += point.X
		averageY += point.Y
	}
	averageX = averageX / len(friskObjects)
	averageY = averageY / len(friskObjects)

	img.Set(averageX, averageY, color.RGBA{255, 0, 0, 255})
	rect.VLine(img, color.RGBA{255, 0, 0, 255}, averageX, averageY-10, averageY+10)
	rect.HLine(img, color.RGBA{255, 0, 0, 255}, averageX-10, averageY, averageX+10)

	// Draw the tiles on the screen
	tiles, err := pathfinding.MakeTiles(*img)
	if err != nil {
		return errors.Wrap(err, "failed to create the screen tiles")
	}

	// Calculate which tile Frisk is in

	currentTile := pathfinding.GetCurrentTile(image.Point{averageX, averageY})
	goalTile := pathfinding.GetGoal()
	if currentTile == nil || goalTile == nil {
		return nil // Wait for a usable player position and walkable exit tile.
	}
	path, err := pathfinding.GetPath(currentTile, goalTile)
	if err != nil {
		return errors.Wrap(err, "failed to generate path")
	}
	// Follow one tile at a time. Short key holds allow the next frame to
	// re-evaluate the screen instead of committing to a long blind movement.
	if len(path) > 1 {
		target := rect.RectangleCenter(path[1].Rectangle)
		dx, dy := target.X-averageX, target.Y-averageY
		key := ""
		if abs(dx) >= abs(dy) {
			if dx > 0 { key = "right" } else if dx < 0 { key = "left" }
		} else {
			if dy > 0 { key = "down" } else if dy < 0 { key = "up" }
		}
		if key != "" {
			if err := pressFor(win, key, 90*time.Millisecond); err != nil {
				return errors.Wrapf(err, "failed to move %s toward the next path tile", key)
			}
		}
	}
	if GridShow {
		for _, tile := range tiles {
			rect.DrawRectangle(img, tile.Color, tile.Rectangle)
		}
	}
	return nil
}

// unstuck deliberately avoids arbitrary keypresses. Random inputs can confirm a
// destructive menu choice or walk Frisk into danger, so wait for recognition instead.
func unstuck(win sys.Window) error {
	return nil
}

// Counts how many frames the state has been unknown
var unknownFrames int

// UnknownUpdate is the update function that just calls unstuck
func UnknownUpdate(objects []object.Object, win sys.Window, img *image.RGBA) error {
	// Stall until items can be found
	unknownFrames++
	// If stalling takes too long, then try to get unstuck
	if unknownFrames == params.FailedLimit+1 {
		fmt.Println("The AI cannot identify the current screen; automatic inputs are suspended.")
	}
	return nil
}

// GetWanted returns a slice of objects based on what the recognized objects wanted are
func GetWanted(objects []object.Object, wanted []object.RecognizableObject) ([]object.RecognizedObject, error) {
	var found []object.RecognizedObject

	// Used to make sure that the types found are the same as the ones wanted
	var wantedCompare []object.RecognizableObject

	// Used to ensure that duplicate items of the same type aren't added to wantedCompare
	var typesUsed []object.RecognizableObject
	for _, wantObj := range wanted {
		for _, obj := range objects {
			if obj.RecogObj.Type == wantObj {
				found = append(found, obj.RecogObj)
				if !contains(typesUsed, obj.RecogObj.Type) {
					wantedCompare = append(wantedCompare, obj.RecogObj.Type)
					typesUsed = append(typesUsed, obj.RecogObj.Type)
				}
			}
		}
	}
	if !equal(wanted, wantedCompare) {
		return found, errors.New("could not get all the needed items")
	}
	return found, nil
}

// Map turns a slice of RecognizableObject into a map that can be searched with the string name
func Map(recObjects []object.RecognizedObject) (recMap map[string][]object.RecognizedObject) {
	recMap = make(map[string][]object.RecognizedObject)
	for _, obj := range recObjects {
		if recMap[obj.Type.Name] == nil {
			recMap[obj.Type.Name] = []object.RecognizedObject{obj}
		} else {
			recMap[obj.Type.Name] = append(recMap[obj.Type.Name], obj)
		}
	}
	return recMap
}

// equal tells whether a and b contain the same elements.
// A nil argument is equivalent to an empty slice.
func equal(a, b []object.RecognizableObject) bool {
	if len(a) != len(b) {
		return false
	}
	for i, v := range a {
		if v != b[i] {
			return false
		}
	}
	return true
}

func contains(s []object.RecognizableObject, e object.RecognizableObject) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}
	return false
}


// pressFor holds a direction briefly when the platform backend supports it.
// Other backends retain their existing single-press behavior.
func pressFor(win sys.Window, key string, duration time.Duration) error {
	if held, ok := win.(interface { PressFor(string, time.Duration) error }); ok {
		return held.PressFor(key, duration)
	}
	return win.Press(key)
}

func abs(value int) int {
	if value < 0 { return -value }
	return value
}
