//go:build linux
// +build linux

package winmanage

import (
	"fmt"

	"github.com/go-vgo/robotgo"
)

func waitForSelect() error {
	ctrl := robotgo.AddEvent("shift")
	mleft := robotgo.AddEvent("mleft")
	if ctrl == 0 && mleft == 0 {
		return nil
	}
	return fmt.Errorf("could not register Shift + left mouse selection (shift=%d, left=%d)", ctrl, mleft)
}
