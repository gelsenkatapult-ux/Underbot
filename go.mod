module gitlab.com/256/Underbot

go 1.22

// Versions below follow the original Gopkg.lock where possible. The Windows
// application still uses the original Ebiten 1.7 and GoCV 0.14 APIs.
require (
	github.com/BurntSushi/xgb v0.0.0-20160522181843-27f122750802
	github.com/BurntSushi/xgbutil v0.0.0-20160919175755-f7c97cef3b4e
	github.com/beefsack/go-astar v0.0.0-20171024231011-f324bbb0d6f7
	github.com/galaktor/gostwriter v0.0.0-20140828112032-b9e5df85bd4a
	github.com/go-vgo/robotgo v0.0.0-20180427105037-08341cc1aa28
	github.com/hajimehoshi/ebiten v1.7.0
	github.com/mitchellh/go-ps v0.0.0-20170309133038-4fdf99ab2936
	github.com/pkg/errors v0.8.0
	gocv.io/x/gocv v0.0.0-20180705114656-0d9b22a6450a
	golang.org/x/image v0.0.0-20180628062038-cc896f830ced
)
