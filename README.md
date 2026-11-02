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

1. Install 64-bit MinGW-w64 with POSIX threads and SEH exceptions. The pinned GoCV 0.14.0 `win_build_opencv.cmd` uses MinGW-w64 GCC 6.3.0 (`x86_64-6.3.0-release-posix-seh-rt_v5-rev1`). The CI may use the archived GCC 7.3.0 bundle from the Qt mirror only as a fallback when the original SourceForge archive is unavailable.
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
   go build -mod=mod -o Underbot.exe .
   ```

The GitHub Actions workflow `.github/workflows/windows-full-build.yml` provisions the native toolchain, builds the pinned OpenCV 3.4.2 version, runs the Go tests, checks that GoCV can load the OpenCV runtime, builds the original GUI and packages the executable with its runtime DLLs. Its downloadable artifact is named `Underbot-Windows-x64`. Version tags beginning with `v` create a **draft** GitHub release after the full build succeeds; a maintainer must review and publish that draft. The manual build above remains useful for local development. The current GoCV Windows installation page describes a newer GoCV/OpenCV combination; use it only as general toolchain background, not as a replacement for the versions pinned by this original source tree.

## Port scope and gameplay status

The Windows port now has a native Win32 window/input/capture backend, a
standalone capture diagnostic, and a CI workflow for building the original GUI.
The original game logic is preserved, with targeted improvements to the
prototype AI:

- Overworld navigation uses an image-bounds-aware A* grid, avoids tiles
  classified as black obstacles, executes the next movement step, and uses
  visible switches, room entrances, the training dummy and Toriel as local
  landmarks. The old right-edge target remains the fallback.
- Battle handling detects small, bright contours inside the fight box and
  applies a reactive dodge heuristic for red/blue souls. For green souls it
  turns the shield toward the nearest detected hazard.
- Pacifist mode is enabled by default. The four-button battle menu is ordered
  by screen position; ACT is preferred until MERCY is highlighted yellow.
  The automatic fight-timing confirmation is disabled in this mode, and
  unknown-state recovery only presses directional keys.

These behaviours are deliberately conservative prototypes, not a certified
playthrough system. Projectile detection can miss attacks or mistake visual
effects for hazards; green-shield direction is inferred from the nearest
visible object and not its trajectory. Enemy-specific ACT choices are not
recognized, so the first visible ACT submenu entry is selected. Many story
puzzles, room transitions, bosses, and scripted events still need explicit
recognition and handling. The bot does not yet have a complete room map or
dialogue/OCR system.

A successful Windows build or unit-test run therefore does **not** prove that
the bot can complete a True Pacifist route. The full real-game playthrough,
including all bosses and the ending sequence, remains unverified and is a
release blocker for claiming the bot is fully autonomous.
