package main

import (
	"sync/atomic"
	"unsafe"
)

func flashClickIndicator() {
	if flashWnd == 0 {
		return
	}
	var pt POINT
	if r, _, _ := pGetCursorPos.Call(uintptr(unsafe.Pointer(&pt))); r == 0 {
		return
	}
	const size = 42
	pSetWindowPos.Call(
		flashWnd,
		^uintptr(0), // HWND_TOPMOST
		uptrI32(pt.X-size/2), uptrI32(pt.Y-size/2),
		size, size,
		SWP_NOACTIVATE,
	)
	pInvalidateRect.Call(flashWnd, 0, 1)
	pShowWindow.Call(flashWnd, SW_SHOWNOACTIVATE)
	pSetTimer.Call(flashWnd, 1, 95, 0)
}

func flashWndProc(hwnd uintptr, msg uint32, wparam, lparam uintptr) uintptr {
	switch msg {
	case WM_PAINT:
		var ps PAINTSTRUCT
		hdc, _, _ := pBeginPaint.Call(hwnd, uintptr(unsafe.Pointer(&ps)))
		if hdc != 0 {
			r := RECT{0, 0, 42, 42}
			pFillRect.Call(hdc, uintptr(unsafe.Pointer(&r)), blackBrush)
			pen, _, _ := pCreatePen.Call(0, 4, uintptr(rgb(245, 86, 35)))
			oldPen, _, _ := pSelectObject.Call(hdc, pen)
			nullBrush, _, _ := pGetStockObject.Call(NULL_BRUSH)
			oldBrush, _, _ := pSelectObject.Call(hdc, nullBrush)
			pEllipse.Call(hdc, 3, 3, 39, 39)
			pSelectObject.Call(hdc, oldBrush)
			pSelectObject.Call(hdc, oldPen)
			pDeleteObject.Call(pen)
		}
		pEndPaint.Call(hwnd, uintptr(unsafe.Pointer(&ps)))
		return 0
	case WM_TIMER:
		pKillTimer.Call(hwnd, 1)
		pShowWindow.Call(hwnd, SW_HIDE)
		return 0
	case WM_NCHITTEST:
		return ^uintptr(0) // HTTRANSPARENT (-1): never intercept mouse input.
	}
	r, _, _ := pDefWindowProc.Call(hwnd, uintptr(msg), wparam, lparam)
	return r
}

func drawStartStopButton(dis *DRAWITEMSTRUCT) {
	if dis == nil {
		return
	}
	isRunning := atomic.LoadInt32(&running) == 1
	selected := dis.ItemState&ODS_SELECTED != 0

	var fill uint32
	var text string
	if isRunning {
		fill = rgb(202, 47, 47)
		if selected {
			fill = rgb(170, 35, 35)
		}
		text = "■  Stop"
	} else {
		fill = rgb(22, 163, 74)
		if selected {
			fill = rgb(18, 135, 60)
		}
		text = "▶  Start"
	}

	brush, _, _ := pCreateSolidBrush.Call(uintptr(fill))
	pen, _, _ := pCreatePen.Call(0, 1, uintptr(fill))
	oldBrush, _, _ := pSelectObject.Call(dis.HDC, brush)
	oldPen, _, _ := pSelectObject.Call(dis.HDC, pen)
	pRoundRect.Call(dis.HDC,
		uintptr(dis.RcItem.Left), uintptr(dis.RcItem.Top), uintptr(dis.RcItem.Right), uintptr(dis.RcItem.Bottom),
		14, 14)
	pSelectObject.Call(dis.HDC, oldBrush)
	pSelectObject.Call(dis.HDC, oldPen)
	pDeleteObject.Call(brush)
	pDeleteObject.Call(pen)

	if buttonFont != 0 {
		old, _, _ := pSelectObject.Call(dis.HDC, buttonFont)
		defer pSelectObject.Call(dis.HDC, old)
	}
	pSetBkMode.Call(dis.HDC, TRANSPARENT)
	pSetTextColor.Call(dis.HDC, uintptr(rgb(255, 255, 255)))
	r := dis.RcItem
	pDrawText.Call(dis.HDC, uintptr(unsafe.Pointer(w(text))), ^uintptr(0), uintptr(unsafe.Pointer(&r)), DT_CENTER|DT_VCENTER|DT_SINGLELINE)
}

func wndProc(hwnd uintptr, msg uint32, wparam, lparam uintptr) uintptr {
	switch msg {
	case WM_COMMAND:
		id := int(loword(wparam))
		code := hiword(wparam)
		if code == BN_CLICKED {
			switch id {
			case idStartStop:
				toggle()
			case idExit:
				pDestroyWindow.Call(mainWnd)
			}
		}
		if code == EN_CHANGE {
			switch id {
			case idHours, idMinutes, idSeconds, idMillis, idClickLimit:
				updateEstimatedRunTime()
			}
		}
		return 0

	case WM_HOTKEY:
		if int(wparam) == HOTKEY_ID {
			toggleFromHotkey()
		}
		return 0

	case WM_APP_TOGGLE:
		toggleFromHotkey()
		return 0

	case WM_TIMER:
		switch int(wparam) {
		case HOTKEY_POLL_TIMER:
			pollGlobalHotkey()
			return 0
		case HOTKEY_FEEDBACK_TIMER:
			pKillTimer.Call(mainWnd, HOTKEY_FEEDBACK_TIMER)
			if h := dlg(idHotkey); h != 0 {
				setText(h, "CTRL + ALT + S")
			}
			return 0
		}

	case WM_APP_CLICKFLASH:
		flashClickIndicator()
		return 0

	case WM_APP_COUNTER:
		updatePerformed()
		return 0

	case WM_APP_FINISHED:
		stopCh = nil
		updatePerformed()
		updateUIState()
		return 0

	case WM_DRAWITEM:
		dis := (*DRAWITEMSTRUCT)(unsafe.Pointer(lparam))
		if dis != nil && int(dis.CtlID) == idStartStop {
			drawStartStopButton(dis)
			return 1
		}

	case WM_CTLCOLORSTATIC:
		hdc := wparam
		ctl := lparam
		pSetBkMode.Call(hdc, TRANSPARENT)
		if c, ok := textColors[ctl]; ok {
			pSetTextColor.Call(hdc, uintptr(c))
		} else {
			pSetTextColor.Call(hdc, uintptr(rgb(45, 45, 45)))
		}
		if b, ok := bgBrushes[ctl]; ok && b != 0 {
			return b
		}
		if whiteBrush != 0 {
			return whiteBrush
		}

	case WM_CTLCOLORBTN:
		if whiteBrush != 0 {
			return whiteBrush
		}

	case WM_DESTROY:
		if atomic.LoadInt32(&running) == 1 {
			atomic.StoreInt32(&running, 0)
			if stopCh != nil {
				close(stopCh)
				stopCh = nil
			}
		}
		uninstallGlobalHotkey()
		for _, f := range []uintptr{regularFont, smallFont, sectionFont, titleFont, hotkeyFont, buttonFont} {
			if f != 0 {
				pDeleteObject.Call(f)
			}
		}
		for _, b := range []uintptr{whiteBrush, lightBlueBrush, blackBrush} {
			if b != 0 {
				pDeleteObject.Call(b)
			}
		}
		pPostQuitMessage.Call(0)
		return 0
	}

	r, _, _ := pDefWindowProc.Call(hwnd, uintptr(msg), wparam, lparam)
	return r
}
