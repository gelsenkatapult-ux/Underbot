package ai

import (
	"fmt"
	"image"
	"image/color"
	"sort"

	"github.com/pkg/errors"
	"gitlab.com/256/Underbot/ai/pathfinding"
	"gitlab.com/256/Underbot/cv/object"
	"gitlab.com/256/Underbot/cv/params"
	"gitlab.com/256/Underbot/cv/rect"
	"gitlab.com/256/Underbot/sys"
)

// GridShow determines if the pathfinding grid is enabled
var GridShow bool

// Should be ran in every update function (other than unknown)
func genUpdate() {
	unknownFrames = 0
	if interactionCooldownFrames > 0 {
		interactionCooldownFrames--
	}
	if battleMenuCooldownFrames > 0 {
		battleMenuCooldownFrames--
	}
	if saveMenuCooldownFrames > 0 {
		saveMenuCooldownFrames--
	}
}

// BattleMenuUpdate only acts on a fully recognized four-button main menu or
// on a submenu that this pacifist controller opened itself. It alternates ACT
// turns with Spare attempts; unrecognized menus are left untouched.
func BattleMenuUpdate(objects []object.Object, win sys.Window, img *image.RGBA) error {
	genUpdate()
	if PacifistMode && battleMenuCooldownFrames > 0 {
		return nil
	}
	options := make([]object.RecognizedObject, 0, 4)
	var heart *object.Object
	for i := range objects {
		candidate := &objects[i]
		if !candidate.Recognized {
			continue
		}
		switch candidate.RecogObj.Type.Name {
		case "battleOption":
			if candidate.RecogObj.Parent != nil {
				options = append(options, candidate.RecogObj)
			}
		case "redHeart":
			if heart == nil {
				heart = candidate
			}
		}
	}
	if heart == nil {
		return nil
	}
	if len(options) == 0 {
		// Some ACT/MERCY submenu entries are text-only. Confirm their first
		// entry only when this controller just opened that submenu itself.
		if PacifistMode && pacifistSubmenuStepsRemaining > 0 {
			if err := win.Press("z"); err != nil {
				return errors.Wrap(err, "couldn't confirm the expected non-violent submenu entry")
			}
			pacifistSubmenuStepsRemaining--
			battleMenuCooldownFrames = 25
		}
		return nil
	}

	// Stable screen order is independent of OpenCV contour IDs, which are not
	// a meaningful menu ordering and caused the original random selection.
	sort.SliceStable(options, func(i, j int) bool {
		left := rect.RectangleCenter(options[i].Parent.Bounds)
		right := rect.RectangleCenter(options[j].Parent.Bounds)
		if left.Y == right.Y {
			return left.X < right.X
		}
		return left.Y < right.Y
	})

	mainMenu := isMainBattleMenu(options)
	if PacifistMode && !mainMenu && pacifistSubmenuStepsRemaining == 0 {
		// A partial/unrecognized menu must never be confirmed: it may be the
		// Fight button with three other contours missed by the CV detector.
		return nil
	}
	targetIndex := 0
	if PacifistMode && mainMenu {
		// Alternate between ACT and a Spare attempt. This makes progress on
		// ordinary encounters without ever selecting FIGHT automatically.
		chooseMercy := pacifistMainMenuTurns%2 == 1
		targetIndex = pacifistOptionIndex(len(options), chooseMercy)
	}
	player := rect.RectangleCenter(heart.Bounds)
	selectedIndex := closestMenuOptionIndex(options, player)
	if selectedIndex < 0 {
		return nil
	}
	if selectedIndex != targetIndex {
		current := rect.RectangleCenter(options[selectedIndex].Parent.Bounds)
		target := rect.RectangleCenter(options[targetIndex].Parent.Bounds)
		direction := movementDirection(current, target)
		if direction != "" {
			if err := win.Press(direction); err != nil {
				return errors.Wrap(err, fmt.Sprintf("couldn't move the menu cursor %s", direction))
			}
		}
		return nil
	}
	if err := win.Press("z"); err != nil {
		return errors.Wrap(err, "couldn't confirm the selected non-violent menu option")
	}
	if PacifistMode {
		battleMenuCooldownFrames = 25
		if mainMenu {
			pacifistMainMenuTurns++
			pacifistSubmenuStepsRemaining = 2
		} else {
			pacifistSubmenuStepsRemaining--
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

// InBattleUpdate performs conservative reactive dodging inside the arena.
// The recognized soul colour determines whether arrow input moves the soul or
// turns the green-heart shield. Unknown/missed hazards never trigger random
// movement.
func InBattleUpdate(objects []object.Object, win sys.Window, img *image.RGBA) error {
	genUpdate()

	var heart *object.Object
	heartName := ""
	var arena image.Rectangle
	for i := range objects {
		candidate := &objects[i]
		if !candidate.Recognized {
			continue
		}
		switch candidate.RecogObj.Type.Name {
		case "redHeart", "greenHeart", "blueHeart":
			if heart == nil {
				heart = candidate
				heartName = candidate.RecogObj.Type.Name
			}
		case "fightBox":
			if arena.Empty() {
				arena = candidate.Bounds
			}
		}
	}

	if heart == nil || arena.Empty() || arena.Dx() < 40 || arena.Dy() < 40 {
		failedRetrieval++
		if failedRetrieval > params.FailedLimit {
			fmt.Println("The AI cannot identify the battle arena or player soul. Holding position.")
			failedRetrieval = 0
		}
		return nil
	}
	failedRetrieval = 0

	player := rect.RectangleCenter(heart.Bounds)
	hazards := findBattleHazards(objects, arena, heart.ID)
	if len(hazards) == 0 {
		return nil
	}

	var direction string
	switch heartName {
	case "greenHeart":
		// The green soul cannot freely move; arrow input rotates its shield.
		if nearestHazardDistance(player, hazards) <= battleDangerRadius+20 {
			direction = nearestHazardDirection(player, hazards)
		}
	case "redHeart", "blueHeart":
		target := chooseSafeBattlePoint(arena, player, hazards)
		if abs(target.X-player.X) > 4 || abs(target.Y-player.Y) > 4 {
			direction = movementDirection(player, target)
		}
	}

	if direction != "" {
		if err := win.Press(direction); err != nil {
			return errors.Wrap(err, fmt.Sprintf("failed to dodge a detected battle hazard with %s", direction))
		}
	}

	if GridShow {
		for _, hazard := range hazards {
			rect.DrawRectangle(img, color.RGBA{R: 255, A: 255}, hazard.Bounds)
		}
	}
	return nil
}

var failedRetrieval int

// SaveUpdate is the function that presses the save button
func SaveUpdate(objects []object.Object, win sys.Window, img *image.RGBA) error {
	genUpdate()
	if saveMenuCooldownFrames > 0 {
		return nil
	}
	recObjects, err := GetWanted(objects, []object.RecognizableObject{
		object.RecMap["redHeart"],
		object.RecMap["saveBox"],
	})
	if err != nil {
		// Stall until items can be found
		failedRetrieval++
		// If stalling takes too long, then try to get unstuck
		if failedRetrieval > params.FailedLimit {
			fmt.Println("The AI is unsure about what is happening. Trying to get unstuck...")
			unstuck(win)
		}
		return nil
	}
	failedRetrieval = 0
	recMap := Map(recObjects)

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
	saveMenuCooldownFrames = 45
	return nil
}

// saveMenuCooldownFrames prevents rapid repeated confirmation at the save screen.
var saveMenuCooldownFrames int

// interactionCooldownFrames prevents repeated interaction inputs from being
// sent on every rendered frame while Frisk is standing next to an object.
var interactionCooldownFrames int

// OutsideBattleUpdate pathfinds through the current room and executes exactly
// one movement input per update. Recognized switches, room entrances and
// selected story characters are used as local landmarks; the right-hand edge
// remains the fallback when no landmark can be seen.
func OutsideBattleUpdate(objects []object.Object, win sys.Window, img *image.RGBA) error {
	genUpdate()
	pacifistMainMenuTurns = 0
	battleMenuCooldownFrames = 0
	pacifistSubmenuStepsRemaining = 0
	friskObjects, _ := GetWanted(objects, object.Frisk)
	if len(friskObjects) == 0 {
		failedRetrieval++
		if failedRetrieval > params.FailedLimit {
			fmt.Println("The AI cannot locate Frisk. Waiting for a recognizable frame...")
			failedRetrieval = 0
		}
		return nil
	}
	failedRetrieval = 0

	// Average all recognized body parts so side-facing and front-facing sprites
	// use the same stable player position.
	var playerX, playerY int
	for _, obj := range friskObjects {
		point := rect.RectangleCenter(obj.Parent.Bounds)
		playerX += point.X
		playerY += point.Y
	}
	playerPoint := image.Pt(playerX/len(friskObjects), playerY/len(friskObjects))

	grid, err := pathfinding.NewGrid(img, params.TileSize)
	if err != nil {
		return errors.Wrap(err, "failed to create the room navigation grid")
	}
	start := grid.ClosestWalkableTile(playerPoint)
	fallback := grid.RightEdgeGoal()
	if start == nil || fallback == nil {
		return nil
	}

	goalPoint, goalName := chooseWorldGoal(objects, playerPoint, fallback.Center())
	goal := grid.ClosestWalkableTile(goalPoint)
	if goal == nil {
		return nil
	}

	// Draw the player marker only after the walkability map was sampled; debug
	// pixels must not accidentally become walls in the navigation grid.
	img.Set(playerPoint.X, playerPoint.Y, color.RGBA{R: 255, A: 255})
	rect.VLine(img, color.RGBA{R: 255, A: 255}, playerPoint.X, playerPoint.Y-10, playerPoint.Y+10)
	rect.HLine(img, color.RGBA{R: 255, A: 255}, playerPoint.X-10, playerPoint.Y, playerPoint.X+10)

	path, err := grid.Path(start, goal)
	if err != nil {
		failedRetrieval++
		if failedRetrieval == 1 {
			fmt.Printf("No walkable route to %s; waiting instead of issuing random movement.\n", goalName)
		}
		if GridShow {
			for _, tile := range grid.Tiles {
				rect.DrawRectangle(img, tile.Color, tile.Rectangle)
			}
		}
		return nil
	}
	failedRetrieval = 0

	// Interact only with visible landmarks that need a confirmation key.
	// Room entrances are approached by movement so Z is not spammed at doors.
	if worldGoalNeedsInteraction(goalName) && withinWorldInteractionRange(playerPoint, goalPoint) {
		if interactionCooldownFrames == 0 {
			if err := win.Press("z"); err != nil {
				return errors.Wrap(err, "failed to interact with the current room landmark")
			}
			interactionCooldownFrames = 45
		}
	} else {
		// If start and goal fall inside the same tile, A* returns a one-tile
		// path. Continue steering toward the exact goal pixel instead of
		// stopping at the tile centre (especially important for doorways).
		moveTarget := goalPoint
		if len(path) > 1 {
			moveTarget = path[1].Center()
		}
		direction := movementDirection(playerPoint, moveTarget)
		if direction != "" && (abs(moveTarget.X-playerPoint.X) > 4 || abs(moveTarget.Y-playerPoint.Y) > 4) {
			if err := win.Press(direction); err != nil {
				return errors.Wrap(err, fmt.Sprintf("failed to move Frisk %s", direction))
			}
		}
	}

	if GridShow {
		for _, tile := range grid.Tiles {
			rect.DrawRectangle(img, tile.Color, tile.Rectangle)
		}
	}
	return nil
}

// unstuck performs a cautious directional nudge. It never presses Z or Enter,
// because an unknown state may be a battle menu and random confirmation could
// accidentally select FIGHT in a pacifist run.
var unstuckDirectionIndex int

func unstuck(win sys.Window) error {
	keys := []string{"up", "right", "down", "left"}
	key := keys[unstuckDirectionIndex%len(keys)]
	unstuckDirectionIndex++
	if err := win.Press(key); err != nil {
		return errors.Wrap(err, fmt.Sprintf("failed to press %s while recovering", key))
	}
	return nil
}

// Counts how many frames the state has been unknown
var unknownFrames int

// UnknownUpdate is the update function that just calls unstuck
func UnknownUpdate(objects []object.Object, win sys.Window, img *image.RGBA) error {
	unknownFrames++
	if unknownFrames == params.FailedLimit+1 {
		fmt.Println("The AI cannot identify the current game screen. It will only try directional recovery inputs.")
	}
	if unknownFrames > params.FailedLimit && (unknownFrames-params.FailedLimit)%15 == 0 {
		if err := unstuck(win); err != nil {
			return errors.Wrap(err, "failed to try directional recovery")
		}
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

 
// GameOverUpdate disables automated input as soon as the game-over screen is
// recognized. This prevents directional recovery from issuing meaningless
// inputs after a failed run.
func GameOverUpdate(objects []object.Object, win sys.Window, img *image.RGBA) error {
	if !Disabled {
		fmt.Println("Undertale Game Over detected. Underbot has disabled the AI.")
	}
	Disabled = true
	return nil
}
