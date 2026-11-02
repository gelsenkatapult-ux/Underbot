# Underbot
## Details
The bot that can play independently play Undertale (and some fangames) in it's entirety (that's the goal anyways)

## Build Instructions
### Linux (X11)
1. If you don't already have ```dep``` installed, use this command to do so: ```curl https://raw.githubusercontent.com/golang/dep/master/install.sh | sh```
2. Build OpenCV from the [GoCV project](https://github.com/hybridgroup/gocv) (OpenCV 3+ might work, but please try the GoCV version before reporting issues here)
3. Get this project with ```git clone https://gitlab.com/256/Underbot.git```
4.  Run ```dep ensure``` to get the needed packages in the root directory
5.  Finally, use ```go build``` and run the executable named ```Underbot```
### Linux (Wayland)
Until Wayland provides a method to interact with other windows, as it is designed to limit interaction between windows, this is unlikely to ever be in the future of this project.
### Other platforms
This is also not supported yet as the project relies on the X11 window system. However, porting is possible in the case of Windows and is being tracked in [Issue #7](https://gitlab.com/256/Underbot/issues/7). If you think you are up to the task, go ahead and give it a shot at porting it over.

## TODO
View tasks that need to be done at [the issue boards](https://gitlab.com/256/Underbot/boards?=)
### Windows (port in progress)

The `windows-port` branch contains a native Win32 backend and a small capture/input diagnostic. It is not yet a complete Windows build of the Underbot GUI: the main application still depends on GoCV/OpenCV and Ebiten, whose Windows toolchain and native libraries have not yet been integrated into the build pipeline.

To build the backend diagnostic on Windows with Go 1.22 or newer:

```powershell
go build -o Underbot-Windows-BackendCheck.exe ./cmd/winbackendcheck
.\Underbot-Windows-BackendCheck.exe
```

The diagnostic looks for the Undertale window (default title: `UNDERTALE`), captures its visible client area, and writes `underbot-window-test.png`. The game must be visible and in the foreground; this desktop capture does not read obscured or minimized windows. An optional exact window title can be passed as the first argument.

The full application build will require a matching GoCV/OpenCV installation and a CGO-capable Windows compiler. See the [GoCV Windows installation guide](https://gocv.io/getting-started/windows/) for the native OpenCV toolchain.
