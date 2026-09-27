# MagzClicker

**A simple, lightweight and open-source auto clicker for Windows.**

No installer. No telemetry. No ads. Just clicks.

MagzClicker clicks wherever your mouse cursor is currently located. There is no need to save a fixed screen position, and you can move the cursor freely while the clicker is running.

## Features

- Adjustable click interval using hours, minutes, seconds and milliseconds
- Left, right and middle mouse buttons
- Single or double click
- Global **Ctrl + Alt + S** start/stop shortcut
- Clicks at the mouse cursor's current position
- Visible orange click indicator around the cursor
- Optional click limit
- Leave the click-count field blank to run until manually stopped
- Automatic estimated run time based on click count and interval
- Built-in **About** dialog with the installed version
- Press **F1** to open the project help / README page
- Portable Windows executable
- No installer required

## How to use

1. Choose the click interval.
2. Select the mouse button and click type.
3. Optionally enter the number of clicks you want MagzClicker to perform.
4. Press **Start** or use **Ctrl + Alt + S**.
5. Move your mouse wherever you want the clicks to happen.
6. Press **Ctrl + Alt + S** again to stop at any time.

If a click limit is set, MagzClicker stops automatically when that number of clicks has been reached. Press **F1** at any time while MagzClicker is focused to open this help page in your default browser.

## Screenshot

A screenshot of the final public release will be added here before launch.

## Download

The current public release is **MagzClicker 1.5.0**.

- [Download MagzClicker.exe](https://github.com/Magzan16/MagzClicker/releases/latest/download/MagzClicker.exe)
- [Download SHA-256 checksum](https://github.com/Magzan16/MagzClicker/releases/latest/download/MagzClicker-SHA256.txt)
- [View all releases](https://github.com/Magzan16/MagzClicker/releases)

The executable keeps the same filename between releases: `MagzClicker.exe`.

## Privacy and security

MagzClicker is fully open source and intentionally small.

The application does **not** contain:

- Telemetry
- Analytics
- Advertising
- Network communication
- Automatic updates
- Autostart functionality
- Background services

MagzClicker does not require administrator privileges.

The source code used to build the application is available in this repository so that anyone can inspect how it works.

## Windows SmartScreen

New or unsigned Windows applications may trigger a Microsoft SmartScreen warning even when the software is safe.

MagzClicker is prepared for digitally signed releases through SignPath. Release files also include SHA-256 checksums so downloaded files can be verified.

## Code signing policy

Free code signing provided by SignPath.io, certificate by SignPath Foundation.

See [CODE_SIGNING.md](CODE_SIGNING.md) for the release signing process, project roles and privacy statement.

## Build from source

### Requirements

- Windows
- Go 1.23 or newer
- Windows PowerShell

From the repository folder, run:

```powershell
.\build.ps1
```

The build script creates:

```text
MagzClicker.exe
```

and prints the SHA-256 hash of the finished executable.

## Source layout

- `main.go` — application startup and Windows message loop
- `clicker.go` — click scheduling, click limits and global hotkey logic
- `ui.go` — application interface
- `windows.go` — window procedures and click indicator
- `win32.go` — Windows API declarations
- `helpers.go` — UI and Win32 helper functions
- `assets/MagzClicker.ico` — MagzClicker application icon
- `build.ps1` — reproducible Windows build command

## Issues and suggestions

Found a bug or have an idea for MagzClicker?

Use the repository's **Issues** section to report bugs or suggest improvements.

For security-related issues, see [SECURITY.md](SECURITY.md).

## License

MagzClicker is released under the [MIT License](LICENSE).
