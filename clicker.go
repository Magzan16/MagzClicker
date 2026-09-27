package main

import (
	"fmt"
	"strconv"
	"strings"
	"sync/atomic"
	"syscall"
	"time"
	"unsafe"
)

func showError(s string) {
	pMessageBox.Call(mainWnd, uintptr(unsafe.Pointer(w(s))), uintptr(unsafe.Pointer(w(appName))), MB_OK|MB_ICONERROR)
}

func parseField(id int) (int, bool) {
	v, err := strconv.Atoi(strings.TrimSpace(getText(dlg(id))))
	if err != nil || v < 0 {
		return 0, false
	}
	return v, true
}

type settings struct {
	interval  time.Duration
	button    int
	clicks    int
	maxClicks uint64 // 0 = unlimited
}

func intervalMilliseconds() (uint64, bool) {
	hours, okH := parseField(idHours)
	mins, okM := parseField(idMinutes)
	secs, okS := parseField(idSeconds)
	millis, okMS := parseField(idMillis)
	if !okH || !okM || !okS || !okMS {
		return 0, false
	}
	if hours > 999 || mins > 59 || secs > 59 || millis > 999 {
		return 0, false
	}
	totalMS := uint64(hours)*60*60*1000 + uint64(mins)*60*1000 + uint64(secs)*1000 + uint64(millis)
	return totalMS, totalMS >= 1
}

func parseClickLimit(showMessage bool) (uint64, bool) {
	t := strings.TrimSpace(getText(dlg(idClickLimit)))
	if t == "" {
		return 0, true
	}
	n, err := strconv.ParseUint(t, 10, 64)
	if err != nil || n == 0 {
		if showMessage {
			showError("Number of clicks must be a whole number greater than 0, or leave it blank to run until stopped.")
		}
		return 0, false
	}
	return n, true
}

func formatRunTime(ms uint64) string {
	if ms < 1000 {
		return fmt.Sprintf("%d ms", ms)
	}
	secs := ms / 1000
	days := secs / 86400
	secs %= 86400
	hours := secs / 3600
	secs %= 3600
	mins := secs / 60
	secs %= 60

	parts := make([]string, 0, 4)
	if days > 0 {
		parts = append(parts, fmt.Sprintf("%d d", days))
	}
	if hours > 0 {
		parts = append(parts, fmt.Sprintf("%d h", hours))
	}
	if mins > 0 {
		parts = append(parts, fmt.Sprintf("%d min", mins))
	}
	if secs > 0 || len(parts) == 0 {
		parts = append(parts, fmt.Sprintf("%d s", secs))
	}
	return strings.Join(parts, " ")
}

func updateEstimatedRunTime() {
	h := dlg(idEstimate)
	if h == 0 {
		return
	}
	limitText := strings.TrimSpace(getText(dlg(idClickLimit)))
	if limitText == "" {
		setText(h, "Estimated run time: Until stopped")
		return
	}
	limit, ok := parseClickLimit(false)
	if !ok {
		setText(h, "Estimated run time: —")
		return
	}
	intervalMS, ok := intervalMilliseconds()
	if !ok {
		setText(h, "Estimated run time: —")
		return
	}
	const maxUint64 = ^uint64(0)
	if limit > maxUint64/intervalMS {
		setText(h, "Estimated run time: Very long")
		return
	}
	setText(h, "Estimated run time: "+formatRunTime(limit*intervalMS))
}

func readSettings() (settings, bool) {
	hours, okH := parseField(idHours)
	mins, okM := parseField(idMinutes)
	secs, okS := parseField(idSeconds)
	millis, okMS := parseField(idMillis)
	if !okH || !okM || !okS || !okMS {
		showError("Click interval fields must contain whole, non-negative numbers.")
		return settings{}, false
	}
	if hours > 999 || mins > 59 || secs > 59 || millis > 999 {
		showError("Use 0–999 hours, 0–59 minutes, 0–59 seconds and 0–999 milliseconds.")
		return settings{}, false
	}

	totalMS := int64(hours)*60*60*1000 + int64(mins)*60*1000 + int64(secs)*1000 + int64(millis)
	if totalMS < 1 {
		showError("Click interval must be at least 1 millisecond.")
		return settings{}, false
	}

	maxClicks, ok := parseClickLimit(true)
	if !ok {
		return settings{}, false
	}

	return settings{
		interval:  time.Duration(totalMS) * time.Millisecond,
		button:    comboIndex(idButton),
		clicks:    comboIndex(idClickType) + 1,
		maxClicks: maxClicks,
	}, true
}

func updateUIState() {
	if atomic.LoadInt32(&running) == 1 {
		setText(dlg(idStatus), "●  Status: Running")
		textColors[dlg(idStatus)] = rgb(26, 127, 55)
	} else {
		setText(dlg(idStatus), "●  Status: Stopped")
		textColors[dlg(idStatus)] = rgb(105, 105, 105)
	}
	pInvalidateRect.Call(dlg(idStartStop), 0, 1)
	pInvalidateRect.Call(dlg(idStatus), 0, 1)
}

func updatePerformed() {
	setText(dlg(idPerformed), fmt.Sprintf("Clicks: %d", atomic.LoadUint64(&performed)))
}

func startClicking() {
	if atomic.LoadInt32(&running) == 1 {
		return
	}
	s, ok := readSettings()
	if !ok {
		return
	}

	atomic.StoreUint64(&performed, 0)
	updatePerformed()
	ch := make(chan struct{})
	stopCh = ch
	atomic.StoreInt32(&running, 1)
	updateUIState()
	go clickLoop(s, ch)
}

func stopClicking() {
	if atomic.CompareAndSwapInt32(&running, 1, 0) {
		ch := stopCh
		stopCh = nil
		if ch != nil {
			close(ch)
		}
		updateUIState()
		updatePerformed()
	}
}

func toggle() {
	if atomic.LoadInt32(&running) == 1 {
		stopClicking()
	} else {
		startClicking()
	}
}

func clickLoop(s settings, stop <-chan struct{}) {
	ticker := time.NewTicker(s.interval)
	defer ticker.Stop()

	lastUI := time.Now()
	lastFlash := time.Time{}
	for {
		select {
		case <-stop:
			pPostMessage.Call(mainWnd, WM_APP_COUNTER, 0, 0)
			return
		case <-ticker.C:
			if atomic.LoadInt32(&running) == 0 {
				return
			}
			for i := 0; i < s.clicks; i++ {
				mouseClick(s.button)
				count := atomic.AddUint64(&performed, 1)
				if lastFlash.IsZero() || time.Since(lastFlash) >= 45*time.Millisecond {
					pPostMessage.Call(mainWnd, WM_APP_CLICKFLASH, 0, 0)
					lastFlash = time.Now()
				}
				if s.maxClicks > 0 && count >= s.maxClicks {
					atomic.StoreInt32(&running, 0)
					pPostMessage.Call(mainWnd, WM_APP_FINISHED, 0, 0)
					return
				}
				if s.clicks > 1 && i == 0 {
					time.Sleep(70 * time.Millisecond)
				}
			}
			if time.Since(lastUI) >= 80*time.Millisecond {
				pPostMessage.Call(mainWnd, WM_APP_COUNTER, 0, 0)
				lastUI = time.Now()
			}
		}
	}
}

func mouseClick(btn int) {
	var down, up uint32
	switch btn {
	case 1:
		down, up = MOUSEEVENTF_RIGHTDOWN, MOUSEEVENTF_RIGHTUP
	case 2:
		down, up = MOUSEEVENTF_MIDDLEDOWN, MOUSEEVENTF_MIDDLEUP
	default:
		down, up = MOUSEEVENTF_LEFTDOWN, MOUSEEVENTF_LEFTUP
	}
	inputs := []INPUT{
		{Type: INPUT_MOUSE, Mi: MOUSEINPUT{DwFlags: down}},
		{Type: INPUT_MOUSE, Mi: MOUSEINPUT{DwFlags: up}},
	}
	pSendInput.Call(uintptr(len(inputs)), uintptr(unsafe.Pointer(&inputs[0])), unsafe.Sizeof(INPUT{}))
}

func keyDown(vk uintptr) bool {
	r, _, _ := pGetAsyncKeyState.Call(vk)
	return uint16(r)&0x8000 != 0
}

func hotkeyKeysDown() bool {
	ctrl := keyDown(VK_CONTROL) || keyDown(VK_LCONTROL) || keyDown(VK_RCONTROL)
	alt := keyDown(VK_MENU) || keyDown(VK_LMENU) || keyDown(VK_RMENU)
	return ctrl && alt && keyDown(VK_S)
}

func showHotkeyDetected() {
	if h := dlg(idHotkey); h != 0 {
		setText(h, "CTRL + ALT + S   ✓")
		pSetTimer.Call(mainWnd, HOTKEY_FEEDBACK_TIMER, 650, 0)
	}
}

func toggleFromHotkey() {
	now := time.Now().UnixNano()
	last := atomic.LoadInt64(&lastHotkeyToggleNS)
	if last != 0 && now-last < int64(350*time.Millisecond) {
		return
	}
	atomic.StoreInt64(&lastHotkeyToggleNS, now)
	showHotkeyDetected()
	toggle()
}

func pollGlobalHotkey() {
	down := hotkeyKeysDown()
	if down {
		if atomic.CompareAndSwapInt32(&pollHotkeyLatched, 0, 1) {
			toggleFromHotkey()
		}
	} else {
		atomic.StoreInt32(&pollHotkeyLatched, 0)
	}
}

func keyboardHookProc(nCode int, wparam, lparam uintptr) uintptr {
	if nCode == HC_ACTION && lparam != 0 {
		k := (*KBDLLHOOKSTRUCT)(unsafe.Pointer(lparam))
		msg := uint32(wparam)

		if k.VkCode == VK_S {
			if msg == WM_KEYUP || msg == WM_SYSKEYUP {
				atomic.StoreInt32(&hookHotkeyLatched, 0)
			} else if msg == WM_KEYDOWN || msg == WM_SYSKEYDOWN {
				if hotkeyKeysDown() && atomic.CompareAndSwapInt32(&hookHotkeyLatched, 0, 1) {
					pPostMessage.Call(mainWnd, WM_APP_TOGGLE, 0, 0)
				}
			}
		}
	}
	r, _, _ := pCallNextHookEx.Call(keyboardHook, uintptr(nCode), wparam, lparam)
	return r
}

func installGlobalHotkey(hinst uintptr) {
	r, _, _ := pRegisterHotKey.Call(0, HOTKEY_ID, MOD_CONTROL|MOD_ALT|MOD_NOREPEAT, VK_S)
	if r == 0 {
		r, _, _ = pRegisterHotKey.Call(0, HOTKEY_ID, MOD_CONTROL|MOD_ALT, VK_S)
	}
	if r != 0 {
		hotkeyRegistered = true
	}

	hookCallback = syscall.NewCallback(keyboardHookProc)
	h, _, _ := pSetWindowsHookEx.Call(WH_KEYBOARD_LL, hookCallback, hinst, 0)
	keyboardHook = h

	pSetTimer.Call(mainWnd, HOTKEY_POLL_TIMER, 15, 0)
}

func uninstallGlobalHotkey() {
	pKillTimer.Call(mainWnd, HOTKEY_POLL_TIMER)
	pKillTimer.Call(mainWnd, HOTKEY_FEEDBACK_TIMER)
	if hotkeyRegistered {
		pUnregisterHotKey.Call(0, HOTKEY_ID)
		hotkeyRegistered = false
	}
	if keyboardHook != 0 {
		pUnhookWindowsHookEx.Call(keyboardHook)
		keyboardHook = 0
	}
}
