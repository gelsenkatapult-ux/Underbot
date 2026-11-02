package rect

import (
	"errors"
	"image"
	"image/color"

	"gitlab.com/256/Underbot/cv/object"
)

func CenterColor(tmpImg image.Image) (color.Color, error) {
	img, ok := tmpImg.(*image.RGBA)
	if !ok { return nil, errors.New("the underlying image is not of type image.RGBA") }
	if img.Bounds().Empty() { return nil, errors.New("cannot get the center color of an empty image") }
	centerPoint := RectangleCenter(img.Bounds())
	return img.At(centerPoint.X, centerPoint.Y), nil
}

func RectangleCenter(bounds image.Rectangle) image.Point {
	return image.Pt((bounds.Min.X+bounds.Max.X)/2, (bounds.Min.Y+bounds.Max.Y)/2)
}

func DrawObject(img *image.RGBA, col color.Color, obj object.Object) {
	DrawRectangle(img, col, obj.Bounds)
}

func DrawRectangle(img *image.RGBA, col color.Color, bounds image.Rectangle) {
	if img == nil || bounds.Empty() { return }
	Rect(img, col, bounds.Min.X, bounds.Min.Y, bounds.Max.X-1, bounds.Max.Y-1)
}

// GetRectangle creates the smallest half-open rectangle containing all points.
func GetRectangle(points []image.Point) image.Rectangle {
	if len(points) == 0 { return image.Rectangle{} }
	minX, minY := points[0].X, points[0].Y
	maxX, maxY := minX, minY
	for _, point := range points[1:] {
		if point.X < minX { minX = point.X }
		if point.Y < minY { minY = point.Y }
		if point.X > maxX { maxX = point.X }
		if point.Y > maxY { maxY = point.Y }
	}
	return image.Rect(minX, minY, maxX+1, maxY+1)
}

func HLine(img *image.RGBA, col color.Color, x1, y, x2 int) {
	if img == nil || y < img.Rect.Min.Y || y >= img.Rect.Max.Y { return }
	for x := x1; x <= x2; x++ { img.Set(x, y, col) }
}

func VLine(img *image.RGBA, col color.Color, x, y1, y2 int) {
	if img == nil || x < img.Rect.Min.X || x >= img.Rect.Max.X { return }
	for y := y1; y <= y2; y++ { img.Set(x, y, col) }
}

func Rect(img *image.RGBA, col color.Color, x1, y1, x2, y2 int) {
	if x2 < x1 || y2 < y1 { return }
	HLine(img, col, x1, y1, x2)
	HLine(img, col, x1, y2, x2)
	VLine(img, col, x1, y1, y2)
	VLine(img, col, x2, y1, y2)
}

func AverageSize(bounds image.Rectangle) int {
	return (bounds.Dx() + bounds.Dy()) / 2
}
