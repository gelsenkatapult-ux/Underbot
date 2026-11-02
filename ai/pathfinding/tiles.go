package pathfinding

import (
	"fmt"
	"image"
	"image/color"
	"math"

	"github.com/beefsack/go-astar"
	"github.com/pkg/errors"
	"gitlab.com/256/Underbot/cv/params"
)

// Tile is one bounded section of the game image used for path planning.
type Tile struct {
	Pos       string
	Coords    image.Point
	Rectangle image.Rectangle
	Color     color.Color
	Cost      int

	grid *Grid
}

// Grid owns one image's pathfinding tiles. Unlike the original package-global
// grid, each screenshot gets fresh dimensions and bounds.
type Grid struct {
	Bounds   image.Rectangle
	TileSize int
	Tiles    []*Tile
	byCoords map[image.Point]*Tile
}

const blockedTileCost = 255

// NewGrid builds a clipped grid over the actual image bounds. The previous
// implementation created its first row/column at negative coordinates, used
// integer division before Ceil, and retained stale tiles if dimensions changed.
func NewGrid(img image.Image, tileSize int) (*Grid, error) {
	if img == nil {
		return nil, errors.New("cannot build a pathfinding grid from a nil image")
	}
	if tileSize <= 0 {
		return nil, errors.New("tile size must be greater than zero")
	}
	bounds := img.Bounds()
	if bounds.Empty() {
		return nil, errors.New("cannot build a pathfinding grid from an empty image")
	}

	rows := int(math.Ceil(float64(bounds.Dy()) / float64(tileSize)))
	columns := int(math.Ceil(float64(bounds.Dx()) / float64(tileSize)))
	grid := &Grid{
		Bounds:   bounds,
		TileSize: tileSize,
		byCoords: make(map[image.Point]*Tile, rows*columns),
		Tiles:    make([]*Tile, 0, rows*columns),
	}

	for row := 0; row < rows; row++ {
		for column := 0; column < columns; column++ {
			minX := bounds.Min.X + column*tileSize
			minY := bounds.Min.Y + row*tileSize
			rect := image.Rect(minX, minY, minX+tileSize, minY+tileSize).Intersect(bounds)
			if rect.Empty() {
				continue
			}
			tile := &Tile{
				Pos:       fmt.Sprintf("(%d, %d)", row, column),
				Coords:    image.Pt(column, row),
				Rectangle: rect,
				Color:     params.TileColor,
				grid:      grid,
			}
			tile.Cost = tileCost(img, rect)
			if tile.Cost >= blockedTileCost {
				tile.Color = color.RGBA{R: 255, A: 255}
			}
			grid.Tiles = append(grid.Tiles, tile)
			grid.byCoords[tile.Coords] = tile
		}
	}
	return grid, nil
}

// tileCost samples a 3x3 pattern in a tile. A majority of black samples marks
// a wall/obstacle; sampling only one centre pixel made narrow contours and
// anti-aliased edges unstable.
func tileCost(img image.Image, bounds image.Rectangle) int {
	black := 0
	for _, fy := range []int{1, 2, 3} {
		for _, fx := range []int{1, 2, 3} {
			x := bounds.Min.X + (bounds.Dx()*fx)/4
			y := bounds.Min.Y + (bounds.Dy()*fy)/4
			if x >= bounds.Max.X {
				x = bounds.Max.X - 1
			}
			if y >= bounds.Max.Y {
				y = bounds.Max.Y - 1
			}
			r, g, b, _ := img.At(x, y).RGBA()
			if r == 0 && g == 0 && b == 0 {
				black++
			}
		}
	}
	if black >= 5 {
		return blockedTileCost
	}
	return 0
}

// TileAt returns the tile containing a point, or nil when the point is outside
// the image bounds.
func (g *Grid) TileAt(point image.Point) *Tile {
	if g == nil || !point.In(g.Bounds) {
		return nil
	}
	coords := image.Pt(
		(point.X-g.Bounds.Min.X)/g.TileSize,
		(point.Y-g.Bounds.Min.Y)/g.TileSize,
	)
	return g.byCoords[coords]
}

// ClosestWalkableTile returns the walkable tile whose centre is nearest to the
// requested point. This also handles goals that fall on the black outline of
// an interactable sprite or doorway.
func (g *Grid) ClosestWalkableTile(point image.Point) *Tile {
	if g == nil {
		return nil
	}
	var best *Tile
	bestDistance := int(^uint(0) >> 1)
	for _, tile := range g.Tiles {
		if tile.Cost >= blockedTileCost {
			continue
		}
		dx := tile.Center().X - point.X
		dy := tile.Center().Y - point.Y
		distance := dx*dx + dy*dy
		if distance < bestDistance {
			best = tile
			bestDistance = distance
		}
	}
	return best
}

// RightEdgeGoal is the original prototype's default room-exit target, now
// constrained to a walkable tile that actually lies inside the screenshot.
func (g *Grid) RightEdgeGoal() *Tile {
	if g == nil {
		return nil
	}
	return g.ClosestWalkableTile(image.Pt(g.Bounds.Max.X-1, g.Bounds.Min.Y+g.Bounds.Dy()/2))
}

// Center returns the centre pixel of a tile.
func (t *Tile) Center() image.Point {
	if t == nil {
		return image.Point{}
	}
	return image.Pt(t.Rectangle.Min.X+t.Rectangle.Dx()/2, t.Rectangle.Min.Y+t.Rectangle.Dy()/2)
}

// GetImage returns the part of the screenshot represented by this tile.
func (t *Tile) GetImage(screen image.RGBA) image.Image {
	return screen.SubImage(t.Rectangle)
}

// PathNeighbors returns only existing, walkable 4-way neighbors. The original
// returned nil neighbors at edges and did not prevent A* from traversing walls.
func (t *Tile) PathNeighbors() []astar.Pather {
	if t == nil || t.grid == nil {
		return nil
	}
	offsets := [...]image.Point{
		image.Pt(0, -1),
		image.Pt(1, 0),
		image.Pt(0, 1),
		image.Pt(-1, 0),
	}
	neighbors := make([]astar.Pather, 0, 4)
	for _, offset := range offsets {
		neighbor := t.grid.byCoords[t.Coords.Add(offset)]
		if neighbor != nil && neighbor.Cost < blockedTileCost {
			neighbors = append(neighbors, neighbor)
		}
	}
	return neighbors
}

func (t *Tile) PathNeighborCost(to astar.Pather) float64 {
	next, ok := to.(*Tile)
	if !ok || next == nil {
		return math.Inf(1)
	}
	return 1 + float64(next.Cost)
}

func (t *Tile) PathEstimatedCost(to astar.Pather) float64 {
	next, ok := to.(*Tile)
	if t == nil || !ok || next == nil {
		return math.Inf(1)
	}
	dx := t.Coords.X - next.Coords.X
	dy := t.Coords.Y - next.Coords.Y
	if dx < 0 {
		dx = -dx
	}
	if dy < 0 {
		dy = -dy
	}
	return float64(dx + dy)
}

// Path finds an A* route between two tiles belonging to this grid.
func (g *Grid) Path(start, goal *Tile) ([]*Tile, error) {
	if g == nil || start == nil || goal == nil {
		return nil, errors.New("path start and goal must both exist")
	}
	if start.grid != g || goal.grid != g {
		return nil, errors.New("path start and goal must belong to the same grid")
	}
	if start.Cost >= blockedTileCost {
		return nil, errors.New("path start is blocked")
	}
	if goal.Cost >= blockedTileCost {
		return nil, errors.New("path goal is blocked")
	}
	if start == goal {
		return []*Tile{start}, nil
	}
	rawPath, _, found := astar.Path(start, goal)
	if !found {
		return nil, errors.New("failed to find a walkable route")
	}
	path := make([]*Tile, 0, len(rawPath))
	for _, node := range rawPath {
		tile, ok := node.(*Tile)
		if !ok || tile == nil {
			return nil, errors.New("A* returned an unexpected path node")
		}
		tile.Color = params.PathColor
		path = append(path, tile)
	}
	return path, nil
}

// MapTiles is retained for compatibility with callers of the old helper.
func MapTiles(tiles []*Tile) map[image.Point]*Tile {
	result := make(map[image.Point]*Tile, len(tiles))
	for _, tile := range tiles {
		if tile != nil {
			result[tile.Coords] = tile
		}
	}
	return result
}

// Legacy helpers operate on the most recently created grid. New code should
// keep and use the Grid value returned by NewGrid instead.
var defaultGrid *Grid

func MakeTiles(img image.RGBA) ([]*Tile, error) {
	grid, err := NewGrid(&img, params.TileSize)
	if err != nil {
		return nil, err
	}
	defaultGrid = grid
	return grid.Tiles, nil
}

func GetGoal() *Tile {
	if defaultGrid == nil {
		return nil
	}
	return defaultGrid.RightEdgeGoal()
}

func GetCurrentTile(point image.Point) *Tile {
	if defaultGrid == nil {
		return nil
	}
	return defaultGrid.ClosestWalkableTile(point)
}

func GetPath(start, goal *Tile) ([]*Tile, error) {
	if start == nil || start.grid == nil {
		return nil, errors.New("path start tile is not initialized")
	}
	return start.grid.Path(start, goal)
}
