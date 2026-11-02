package pathfinding

import (
	"image"
	"image/color"
	"testing"

	"gitlab.com/256/Underbot/cv/params"
)

func resetGridState() {
	tiles = nil
	tileMap = nil
	gridWidth, gridHeight = 0, 0
	rows, columns = 0, 0
}

func TestMakeTilesCoversPartialEdges(t *testing.T) {
	oldSize := params.TileSize
	params.TileSize = 10
	defer func() { params.TileSize = oldSize; resetGridState() }()
	resetGridState()

	img := image.NewRGBA(image.Rect(0, 0, 25, 17))
	for y := 0; y < 17; y++ {
		for x := 0; x < 25; x++ {
			img.SetRGBA(x, y, color.RGBA{R: 255, G: 255, B: 255, A: 255})
		}
	}
	got, err := MakeTiles(*img)
	if err != nil { t.Fatalf("MakeTiles returned error: %v", err) }
	if len(got) != 9 { t.Fatalf("got %d tiles, want 9", len(got)) }
	if got[8].Rectangle != image.Rect(20, 10, 25, 17) {
		t.Fatalf("last tile bounds = %v, want (20,10)-(25,17)", got[8].Rectangle)
	}
}

func TestMakeTilesRejectsNonPositiveTileSize(t *testing.T) {
	oldSize := params.TileSize
	params.TileSize = 0
	defer func() { params.TileSize = oldSize; resetGridState() }()
	resetGridState()

	img := image.NewRGBA(image.Rect(0, 0, 20, 20))
	if _, err := MakeTiles(*img); err == nil {
		t.Fatal("MakeTiles accepted a zero tile size")
	}
}
