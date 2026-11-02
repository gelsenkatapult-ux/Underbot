package rect

import (
	"image"
	"image/color"
	"testing"
)

func TestGetRectangleUsesHalfOpenPixelBounds(t *testing.T) {
	got := GetRectangle([]image.Point{{X: 4, Y: 7}, {X: 8, Y: 10}, {X: 5, Y: 9}})
	want := image.Rect(4, 7, 9, 11)
	if got != want { t.Fatalf("GetRectangle() = %v, want %v", got, want) }
}
func TestGetRectangleSinglePointIsOnePixel(t *testing.T) {
	got := GetRectangle([]image.Point{{X: 3, Y: 6}})
	want := image.Rect(3, 6, 4, 7)
	if got != want { t.Fatalf("GetRectangle() = %v, want %v", got, want) }
}
func TestGetRectangleEmptyInput(t *testing.T) {
	if got := GetRectangle(nil); !got.Empty() { t.Fatalf("GetRectangle(nil) = %v, want empty", got) }
}
func TestDrawRectangleRespectsHalfOpenBounds(t *testing.T) {
	img := image.NewRGBA(image.Rect(0, 0, 4, 4))
	DrawRectangle(img, color.White, image.Rect(1, 1, 3, 3))
	if got := img.RGBAAt(3, 2); got != (color.RGBA{}) { t.Fatalf("pixel outside rectangle was drawn: %v", got) }
	if got := img.RGBAAt(2, 2); got != (color.RGBA{}) { t.Fatalf("rectangle interior should remain untouched: %v", got) }
	if got := img.RGBAAt(1, 1); got != color.White { t.Fatalf("rectangle border missing: %v", got) }
}
