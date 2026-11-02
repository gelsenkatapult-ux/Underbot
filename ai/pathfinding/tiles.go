package pathfinding

import (
	"fmt"
	"image"
	"image/color"
	"math"

	"github.com/beefsack/go-astar"
	"github.com/pkg/errors"
	"gitlab.com/256/Underbot/cv/params"
	"gitlab.com/256/Underbot/cv/rect"
)

// Tile is a section of the screen for pathfinding
type Tile struct {
	Pos       string      // String representation of the position. Ex: (1, 3)
	Coords    image.Point // The number representation of the above. This is not in pixels
	Rectangle image.Rectangle
	Color     color.Color // What color the tile should be rendered with
	Cost      int         // How long it takes to get past the tile
}

// How many rows and columns there are of tiles
var rows int
var columns int

var tiles []*Tile
var tileMap map[image.Point]*Tile
var gridWidth, gridHeight int

// MakeTiles creates a slice of tiles based on image dimensions
func MakeTiles(img image.RGBA) ([]*Tile, error) {
	if params.TileSize <= 0 {
		return nil, errors.New("tile size must be greater than zero")
	}
	width := img.Bounds().Dx()
	height := img.Bounds().Dy()
	if width <= 0 || height <= 0 {
		return nil, errors.New("cannot create a pathfinding grid from an empty image")
	}
	if len(tiles) == 0 || width != gridWidth || height != gridHeight {
		gridWidth, gridHeight = width, height
		rows = int(math.Ceil(float64(height) / float64(params.TileSize)))
		columns = int(math.Ceil(float64(width) / float64(params.TileSize)))
		tiles = nil
		tileMap = nil
		for r := 0; r < rows; r++ {
			for c := 0; c < columns; c++ {
				left, top := c*params.TileSize, r*params.TileSize
				right, bottom := left+params.TileSize, top+params.TileSize
				if right > width { right = width }
				if bottom > height { bottom = height }
				tileRect := image.Rect(left, top, right, bottom)
				newTile := Tile{Pos: fmt.Sprintf("(%v, %v)", r, c), Rectangle: tileRect, Coords: image.Point{X: c, Y: r}, Color: params.TileColor}
				cost, err := newTile.GetCost(img)
				if err != nil { return []*Tile{}, errors.Wrap(err, "failed to calculate the cost") }
				newTile.Cost = cost
				tiles = append(tiles, &newTile)
			}
		}
	} else {
		for _, tile := range tiles {
			cost, err := tile.GetCost(img)
			if err != nil {
				return []*Tile{}, errors.Wrap(err, "failed to calculate the cost")
			}
			tile.Cost = cost
		}
	}
	if len(tileMap) == 0 {
		tileMap = MapTiles(tiles)
	}
	return tiles, nil
}

// GetCost calculates the movement cost for a tile
func (t *Tile) GetCost(screen image.RGBA) (int, error) {
	centerColor, err := rect.CenterColor(t.GetImage(screen))
	if err != nil {
		return 0, errors.Wrap(err, "could not get the center color")
	}

	// If the tile is black, then set cost to 255, because it's impossible to go there
	if (centerColor == color.RGBA{0, 0, 0, 255}) {
		t.Color = color.RGBA{255, 0, 0, 255}
		return 255, nil
	}
	t.Color = params.TileColor
	return 0, nil
}

// PathNeighbors returns a slice of the tiles neighboring a tile
func (t *Tile) PathNeighbors() []astar.Pather {
	tileX := t.Coords.X
	tileY := t.Coords.Y

	neighbors := []astar.Pather{}
	for _, point := range []image.Point{
		{X: tileX, Y: tileY - 1},
		{X: tileX + 1, Y: tileY},
		{X: tileX, Y: tileY + 1},
		{X: tileX - 1, Y: tileY},
	} {
		if neighbor := tileMap[point]; neighbor != nil && neighbor.Cost < 255 {
			neighbors = append(neighbors, neighbor)
		}
	}
	return neighbors
}

// GetImage shows the image inside of a tile. The function should be passed the image of the whole screen
func (t *Tile) GetImage(screen image.RGBA) image.Image {
	return screen.SubImage(t.Rectangle)
}

func (t *Tile) PathNeighborCost(to astar.Pather) float64 {
	return float64(to.(*Tile).Cost + 1)
}

func (t *Tile) PathEstimatedCost(to astar.Pather) float64 {
	if t == nil {
		return 255
	}
	toT := to.(*Tile)
	absX := t.Coords.X - toT.Coords.X
	if absX < 0 {
		absX = -absX
	}
	absY := t.Coords.Y - toT.Coords.Y
	if absY < 0 {
		absY = -absY
	}
	return float64(absX + absY)
}

// MapTiles allows the access to a tile based on the point used
func MapTiles(tiles []*Tile) map[image.Point]*Tile {
	tileMap := make(map[image.Point]*Tile)
	for _, tile := range tiles {
		tileMap[tile.Coords] = tile
	}
	return tileMap
}

// GetPath uses A* pathfinding to get the best route from one tile to another
func GetPath(a *Tile, b *Tile) ([]*Tile, error) {
	if a == nil || b == nil {
		return nil, errors.New("pathfinding start or goal tile is missing")
	}
	if a.Cost >= 255 || b.Cost >= 255 {
		return nil, errors.New("pathfinding start or goal tile is blocked")
	}
	path, _, found := astar.Path(a, b)
	if !found {
		return nil, errors.New("failed to find route")
	}
	var tempTiles []*Tile
	for _, pather := range path {
		tile := pather.(*Tile)
		tile.Color = params.PathColor
		tempTiles = append(tempTiles, tile)
	}
	return tempTiles, nil
}

// GetGoal returns the tile that the pathfinding should aim for
func GetGoal() *Tile {
	if columns <= 0 || rows <= 0 || len(tileMap) == 0 {
		return nil
	}
	middlePoint := image.Point{X: columns - 1, Y: (rows - 1) / 2}
	goal := tileMap[middlePoint]
	if goal != nil && goal.Cost >= 255 {
		// Search outward from the intended exit tile for the nearest walkable tile.
		var nearest *Tile
		bestDistance := int(^uint(0) >> 1)
		for _, tile := range tiles {
			if tile.Cost >= 255 { continue }
			dx, dy := tile.Coords.X-middlePoint.X, tile.Coords.Y-middlePoint.Y
			distance := dx*dx + dy*dy
			if distance < bestDistance { bestDistance, nearest = distance, tile }
		}
		return nearest
	}
	return goal
}

// GetCurrentTile finds the tile that a point is in the most
func GetCurrentTile(fPoint image.Point) *Tile {
	if len(tiles) == 0 { return nil }
	var nearest *Tile
	bestDistance := int(^uint(0) >> 1)
	for _, tile := range tiles {
		if fPoint.In(tile.Rectangle) {
			tile.Color = color.RGBA{22, 193, 0, 255}
			return tile
		}
		center := rect.RectangleCenter(tile.Rectangle)
		dx, dy := center.X-fPoint.X, center.Y-fPoint.Y
		distance := dx*dx + dy*dy
		if distance < bestDistance {
			bestDistance = distance
			nearest = tile
		}
	}
	if nearest != nil { nearest.Color = color.RGBA{22, 193, 0, 255} }
	return nearest
}
