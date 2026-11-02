package combat

import (
	"image"
	"image/color"
)

// DetectProjectiles extracts compact, color-coded components from the battle
// board. It is a conservative first-pass detector for Undertale's attack
// colors; callers should exclude the player SOUL and validate detections.
func DetectProjectiles(img *image.RGBA, board image.Rectangle, exclude []image.Rectangle) []ProjectileObservation {
	if img == nil || board.Empty() { return nil }
	board = board.Intersect(img.Bounds())
	if board.Empty() { return nil }
	width, height := board.Dx(), board.Dy()
	classes := make([]AttackClass, width*height)
	visited := make([]bool, width*height)
	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			point := image.Pt(board.Min.X+x, board.Min.Y+y)
			if inAnyRect(point, exclude) { continue }
			r, g, b, a := img.At(point.X, point.Y).RGBA()
			if a == 0 { continue }
			classes[y*width+x] = classifyAttackColor(color.RGBA{
				R: uint8(r >> 8), G: uint8(g >> 8), B: uint8(b >> 8), A: uint8(a >> 8),
			})
		}
	}

	var found []ProjectileObservation
	neighbors := [...]image.Point{{-1,-1},{0,-1},{1,-1},{-1,0},{1,0},{-1,1},{0,1},{1,1}}
	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			index := y*width+x
			class := classes[index]
			if class == "" || visited[index] { continue }
			queue := []image.Point{{X:x,Y:y}}
			visited[index] = true
			minX, maxX, minY, maxY, pixels := x, x, y, y, 0
			for len(queue) > 0 {
				current := queue[0]
				queue = queue[1:]
				pixels++
				if current.X < minX { minX = current.X }; if current.X > maxX { maxX = current.X }
				if current.Y < minY { minY = current.Y }; if current.Y > maxY { maxY = current.Y }
				for _, offset := range neighbors {
					nx, ny := current.X+offset.X, current.Y+offset.Y
					if nx < 0 || ny < 0 || nx >= width || ny >= height { continue }
					ni := ny*width+nx
					if !visited[ni] && classes[ni] == class {
						visited[ni] = true
						queue = append(queue, image.Pt(nx, ny))
					}
				}
			}
			// Reject isolated antialiasing noise and very large filled regions.
			if pixels < 3 || pixels > 5000 { continue }
			bounds := image.Rect(board.Min.X+minX, board.Min.Y+minY, board.Min.X+maxX+1, board.Min.Y+maxY+1)
			found = append(found, ProjectileObservation{Bounds: bounds, Class: class})
		}
	}
	return found
}

func inAnyRect(point image.Point, rectangles []image.Rectangle) bool {
	for _, r := range rectangles { if point.In(r) { return true } }
	return false
}

func classifyAttackColor(c color.RGBA) AttackClass {
	r, g, b := int(c.R), int(c.G), int(c.B)
	max, min := r, r
	if g > max { max = g }; if b > max { max = b }
	if g < min { min = g }; if b < min { min = b }
	if min >= 220 && max-min <= 35 { return AttackWhite }
	if min >= 70 && max <= 190 && max-min <= 24 { return AttackGray }
	if r >= 180 && g <= 100 && b <= 110 { return AttackWhite } // red attacks are damaging unless a verified warning glyph is known
	if g >= 155 && r <= 130 && b <= 155 { return AttackGreen }
	if b >= 155 && g >= 125 && r <= 125 { return AttackLightBlue }
	if r >= 180 && g >= 80 && g <= 190 && b <= 110 { return AttackOrange }
	if r >= 180 && g >= 160 && b <= 110 { return AttackYellow }
	return ""
}
