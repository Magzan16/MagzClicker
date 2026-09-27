# MagzClicker

MagzClicker is a lightweight, portable auto clicker for Windows. It clicks at the mouse cursor's current position and does not save screen coordinates.

## Features

- Adjustable interval: hours, minutes, seconds and milliseconds
- Left, right and middle mouse buttons
- Single or double click
- Global **Ctrl + Alt + S** start/stop shortcut
- Visible click pulse around the cursor
- Optional click limit; leave it blank to run until stopped
- Estimated run time based on click count and interval
- No installer required

## Privacy

MagzClicker does not contain networking, telemetry, analytics, advertising, auto-update, autostart or background-service code. It does not require administrator rights.

## Build

Requirements:

- Go 1.23 or newer
- Python 3
- Windows PowerShell

Build from PowerShell:

```powershell
.\build.ps1
```

The build script creates `MagzClicker.exe` and prints its SHA-256 hash.

## Source layout

- `main.go` — application startup and message loop
- `clicker.go` — click scheduling, limits and hotkey logic
- `ui.go` / `windows.go` — Win32 interface and drawing
- `win32.go` / `helpers.go` — Windows API declarations and helpers
- `tools/make_icon.py` — generates the application icon
- `tools/embed_icon.py` — embeds the ICO resource in the Windows executable
- `build.ps1` — reproducible build command

## License

MIT. See [LICENSE](LICENSE).
