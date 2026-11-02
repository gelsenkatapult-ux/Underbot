package combat

import (
	"image"
	"image/color"
	"testing"
)

func TestDetectProjectilesFindsColoredComponents(t *testing.T) {
	img := image.NewRGBA(image.Rect(0, 0, 80, 60))
	for y := 10; y < 15; y++ { for x := 12; x < 18; x++ { img.SetRGBA(x,y,color.RGBA{255,255,255,255}) } }
	for y := 30; y < 36; y++ { for x := 50; x < 56; x++ { img.SetRGBA(x,y,color.RGBA{255,140,0,255}) } }
	got := DetectProjectiles(img, image.Rect(0,0,80,60), nil)
	if len(got) != 2 { t.Fatalf("got %d components, want 2: %#v", len(got), got) }
	if got[0].Bounds != image.Rect(12,10,18,15) || got[0].Class != AttackWhite { t.Errorf("unexpected first component: %#v", got[0]) }
	if got[1].Bounds != image.Rect(50,30,56,36) || got[1].Class != AttackOrange { t.Errorf("unexpected second component: %#v", got[1]) }
}

func TestDetectProjectilesExcludesPlayerAndClipsBoard(t *testing.T) {
	img := image.NewRGBA(image.Rect(0,0,40,40))
	for y:=10; y<16; y++ { for x:=10; x<16; x++ { img.SetRGBA(x,y,color.RGBA{255,0,0,255}) } }
	got := DetectProjectiles(img, image.Rect(5,5,30,30), []image.Rectangle{image.Rect(9,9,17,17)})
	if len(got)!=0 { t.Fatalf("excluded SOUL still detected: %#v",got) }
}

func TestDetectProjectilesIgnoresTinyNoise(t *testing.T) {
	img := image.NewRGBA(image.Rect(0,0,20,20))
	img.SetRGBA(5,5,color.RGBA{255,255,255,255})
	got := DetectProjectiles(img,image.Rect(0,0,20,20),nil)
	if len(got)!=0 { t.Fatalf("tiny noise detected: %#v",got) }
}
