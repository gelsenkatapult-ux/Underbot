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

The full application uses the same legacy dependency set as the original project. In particular, `Gopkg.lock` pins GoCV to the July 2018 v0.14.0 source revision, whose Windows cgo file expects **OpenCV 3.4.2** and MinGW libraries with the `342` suffix. Do not install a newer OpenCV release for this pinned GoCV revision.

To prepare a Windows machine for the original application:

1. Install 64-bit MinGW-w64 with POSIX threads and SEH exceptions. The original GoCV 0.14.0 Windows instructions were written for MinGW-w64 7.x (the 7.3.0 toolchain was used in the historical setup).
2. Install CMake.
3. Check out the exact GoCV source revision used by the original lock file and build its matching OpenCV version:

   ```powershell
   git clone https://github.com/hybridgroup/gocv.git
   cd gocv
   git checkout 0d9b22a6450a7d81761a4b763f58fb99e962ab47
   .\win_build_opencv.cmd
   ```

4. Add `C:\opencv\build\install\x64\mingw\bin` and the MinGW `bin` directory to the Windows `PATH`. The pinned GoCV cgo configuration expects OpenCV headers under `C:\opencv\build\install\include` and libraries under `C:\opencv\build\install\x64\mingw\lib`.
5. From the Underbot repository root, enable CGO, download the pinned Go modules and build the original GUI:

   ```powershell
   $env:CGO_ENABLED = '1'
   go mod download
   go build -o Underbot.exe .
   ```

This Windows GUI build is not yet verified in CI because it needs the matching native OpenCV 3.4.2 libraries and a CGO-capable MinGW compiler. The automated workflow currently verifies the Win32 backend and its diagnostic tool, not the complete GUI. The current GoCV Windows installation page describes a newer GoCV/OpenCV combination; use it only as general toolchain background, not as a replacement for the versions pinned by this original source tree.

## Port scope and gameplay status

The Windows branch is a platform port of the original project, not a replacement
or completion of its game-playing logic. The original AI remains a prototype:
outside battle it calculates a path but its movement-execution loop is still
commented out; during enemy attack phases it builds the tile view but does not
yet implement the movement/dodge controller. Battle interaction is limited to
the original menu/dialogue/save handling.

For that reason, a successful Windows build or backend smoke test must not be
interpreted as proof that the bot can complete an entire Undertale run. This
port preserves those original limitations deliberately rather than silently
changing the game logic. The Windows release workflow validates compilation,
repository tests, GoCV/OpenCV runtime loading and packaging; a full real-game
playthrough remains a separate, currently unverified test.
