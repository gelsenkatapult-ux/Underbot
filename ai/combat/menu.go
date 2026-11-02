package combat

import (
	"image"
	"sort"
)

// PacifistMenuTarget returns the horizontal center of ACT in Undertale's
// standard four-option battle menu (FIGHT, ACT, ITEM, MERCY). It refuses
// ambiguous layouts instead of guessing from detector IDs.
func PacifistMenuTarget(options []image.Rectangle) (image.Point, bool) {
	if len(options) != 4 {
		return image.Point{}, false
	}
	ordered := append([]image.Rectangle(nil), options...)
	sort.Slice(ordered, func(i, j int) bool {
		return ordered[i].Min.X < ordered[j].Min.X
	})
	// The standard battle menu is a single horizontal row. Reject a layout
	// whose option centers are not at approximately the same vertical level.
	y := center(ordered[0]).Y
	for _, option := range ordered[1:] {
		if absInt(center(option).Y-y) > 12 {
			return image.Point{}, false
		}
	}
	return center(ordered[1]), true // ACT is the second option.
}

func center(r image.Rectangle) image.Point {
	return image.Point{X: (r.Min.X + r.Max.X) / 2, Y: (r.Min.Y + r.Max.Y) / 2}
}

func absInt(v int) int {
	if v < 0 { return -v }
	return v
}
