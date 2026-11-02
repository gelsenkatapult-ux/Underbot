package pathfinding

import (
	"image"
	"image/color"
	"testing"
)

func whiteImage(width, height int) *image.RGBA {
	img := image.NewRGBA(image.Rect(0, 0, width, height))
	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			img.SetRGBA(x, y, color.RGBA{R: 255, G: 255, B: 255, A: 255})
		}
	}
	return img
}

func TestNewGridCoversImageBoundsWithoutNegativeTiles(t *testing.T) {
	img := whiteImage(25, 21)
	grid, err := NewGrid(img, 10)
	if err != nil {
		t.Fatalf("NewGrid returned error: %v", err)
	}
	if got, want := len(grid.Tiles), 9; got != want {
		t.Fatalf("tile count = %d, want %d", got, want)
	}
	for _, tile := range grid.Tiles {
		if !tile.Rectangle.In(grid.Bounds) {
			t.Fatalf("tile %s rectangle %v exceeds image bounds %v", tile.Pos, tile.Rectangle, grid.Bounds)
		}
		if tile.Rectangle.Min.X < 0 || tile.Rectangle.Min.Y < 0 {
			t.Fatalf("tile %s has negative coordinates: %v", tile.Pos, tile.Rectangle)
		}
	}
	last := grid.byCoords[image.Pt(2, 2)]
	if last == nil || last.Rectangle != image.Rect(20, 20, 25, 21) {
		t.Fatalf("last clipped tile = %#v, want rectangle (20,20)-(25,21)", last)
	}
}

func TestTileAtUsesActualImageCoordinates(t *testing.T) {
	img := whiteImage(30, 20)
	grid, err := NewGrid(img, 10)
	if err != nil {
		t.Fatal(err)
	}
	if got := grid.TileAt(image.Pt(0, 0)); got == nil || got.Coords != image.Pt(0, 0) {
		t.Fatalf("TileAt(0,0) = %#v, want tile (0,0)", got)
	}
	if got := grid.TileAt(image.Pt(29, 19)); got == nil || got.Coords != image.Pt(2, 1) {
		t.Fatalf("TileAt(29,19) = %#v, want tile (2,1)", got)
	}
	if got := grid.TileAt(image.Pt(30, 19)); got != nil {
		t.Fatalf("TileAt(outside) = %#v, want nil", got)
	}
}

func TestPathAvoidsBlackObstacleTiles(t *testing.T) {
	img := whiteImage(50, 30)
	// Close the direct middle route. The bottom row remains a valid detour.
	for y := 0; y < 20; y++ {
		for x := 20; x < 30; x++ {
			img.SetRGBA(x, y, color.RGBA{A: 255})
		}
	}
	grid, err := NewGrid(img, 10)
	if err != nil {
		t.Fatal(err)
	}
	start := grid.TileAt(image.Pt(5, 15))
	goal := grid.TileAt(image.Pt(45, 15))
	if start == nil || goal == nil {
		t.Fatal("start or goal tile was not found")
	}
	blocked := grid.byCoords[image.Pt(2, 1)]
	if blocked == nil || blocked.Cost < blockedTileCost {
		t.Fatal("expected the middle obstacle tile to be blocked")
	}
	path, err := grid.Path(start, goal)
	if err != nil {
		t.Fatalf("Path failed to find a detour: %v", err)
	}
	if len(path) < 7 {
		t.Fatalf("path length = %d, expected a route around the obstacle", len(path))
	}
	for i, tile := range path {
		if tile.Cost >= blockedTileCost {
			t.Fatalf("path contains blocked tile %s at index %d", tile.Pos, i)
		}
		if i > 0 {
			dx := tile.Coords.X - path[i-1].Coords.X
			dy := tile.Coords.Y - path[i-1].Coords.Y
			if (dx < -1 || dx > 1) || (dy < -1 || dy > 1) || (dx == 0 && dy == 0) || (dx != 0 && dy != 0) {
				t.Fatalf("path has non-adjacent steps: %v then %v", path[i-1].Coords, tile.Coords)
			}
		}
	}
}

func TestRightEdgeGoalIsInsideAndWalkable(t *testing.T) {
	img := whiteImage(25, 21)
	grid, err := NewGrid(img, 10)
	if err != nil {
		t.Fatal(err)
	}
	goal := grid.RightEdgeGoal()
	if goal == nil {
		t.Fatal("RightEdgeGoal returned nil")
	}
	if !goal.Center().In(grid.Bounds) {
		t.Fatalf("goal centre %v is outside bounds %v", goal.Center(), grid.Bounds)
	}
	if goal.Cost >= blockedTileCost {
		t.Fatalf("goal tile %s is blocked", goal.Pos)
	}
}

func TestPathRejectsTilesFromDifferentGrids(t *testing.T) {
	first, err := NewGrid(whiteImage(20, 20), 10)
	if err != nil {
		t.Fatal(err)
	}
	second, err := NewGrid(whiteImage(20, 20), 10)
	if err != nil {
		t.Fatal(err)
	}
	if _, err := first.Path(first.Tiles[0], second.Tiles[1]); err == nil {
		t.Fatal("Path accepted a goal tile from another grid")
	}
}
