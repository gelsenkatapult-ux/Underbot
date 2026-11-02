module gitlab.com/256/Underbot

go 1.22

// Exact revisions from the original Gopkg.lock are retained so the Windows
// port builds the original 2018 application rather than upgrading its stack.
require (
	github.com/BurntSushi/xgb v0.0.0-20160522181843-27f122750802
	github.com/BurntSushi/xgbutil v0.0.0-20160919175755-f7c97cef3b4e
	github.com/StackExchange/wmi v0.0.0-20180116203802-5d049714c4a6
	github.com/beefsack/go-astar v0.0.0-20171024231011-f324bbb0d6f7
	github.com/galaktor/gostwriter v0.0.0-20140828112032-b9e5df85bd4a
	github.com/go-gl/gl v0.0.0-20180407155706-68e253793080
	github.com/go-gl/glfw v0.0.0-20180426074136-46a8d530c326
	github.com/go-ole/go-ole v0.0.0-20171110160706-a41e3c4b706f
	github.com/go-vgo/robotgo v0.0.0-20180427105037-08341cc1aa28
	github.com/gopherjs/gopherjs v0.0.0-20180628210949-0892b62f0d9f
	github.com/gopherjs/webgl v0.0.0-20180508003723-39bd6d41eeb5
	github.com/hajimehoshi/ebiten v1.7.0
	github.com/mitchellh/go-ps v0.0.0-20170309133038-4fdf99ab2936
	github.com/pkg/errors v0.8.0
	github.com/shirou/gopsutil v0.0.0-20180625081143-4a180b209f5f
	github.com/shirou/w32 v0.0.0-20160930032740-bb4de0191aa4
	github.com/vcaesar/imgo v0.0.0-20180307100924-fece97d8745e
	gocv.io/x/gocv v0.0.0-20180705114656-0d9b22a6450a
	golang.org/x/exp v0.0.0-20180625033341-f9fa0fefb1e1
	golang.org/x/image v0.0.0-20180628062038-cc896f830ced
	golang.org/x/mobile v0.0.0-20180618222554-6621de06e1e9
	golang.org/x/sys v0.0.0-20180704094941-151529c776cd
	golang.org/x/tools v0.0.0-20180705090007-16f8f9bb7201
)
