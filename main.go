package main

import (
	"runtime"
	"syscall"
	"unsafe"
)

func main() {
	runtime.LockOSThread()

	pSetProcessDPIAware.Call()
	icc := INITCOMMONCONTROLSEX{DwSize: uint32(unsafe.Sizeof(INITCOMMONCONTROLSEX{})), DwICC: ICC_UPDOWN_CLASS}
	pInitCommonControlsEx.Call(uintptr(unsafe.Pointer(&icc)))

	hinst, _, _ := pGetModuleHandle.Call(0)
	cls := w("MagzClickerMain")

	hIcon, _, _ := pLoadImage.Call(hinst, 1, IMAGE_ICON, 32, 32, 0)
	hIconSm, _, _ := pLoadImage.Call(hinst, 1, IMAGE_ICON, 16, 16, 0)
	hCursor, _, _ := pLoadCursor.Call(0, IDC_ARROW)

	wc := WNDCLASSEX{
		CbSize:        uint32(unsafe.Sizeof(WNDCLASSEX{})),
		LpfnWndProc:   syscall.NewCallback(wndProc),
		HInstance:     hinst,
		HIcon:         hIcon,
		HCursor:       hCursor,
		HbrBackground: COLOR_WINDOW + 1,
		LpszClassName: cls,
		HIconSm:       hIconSm,
	}
	if r, _, _ := pRegisterClassEx.Call(uintptr(unsafe.Pointer(&wc))); r == 0 {
		return
	}

	title := appName
	mainWnd, _, _ = pCreateWindowEx.Call(
		0,
		uintptr(unsafe.Pointer(cls)),
		uintptr(unsafe.Pointer(w(title))),
		WS_OVERLAPPED|WS_CAPTION|WS_SYSMENU|WS_MINIMIZEBOX,
		uintptr(CW_USEDEFAULT), uintptr(CW_USEDEFAULT),
		832, 780,
		0, 0, hinst, 0,
	)
	if mainWnd == 0 {
		return
	}

	buildUI(hinst)

	flashClass := w("MagzClickerFlash")
	flashWC := WNDCLASSEX{
		CbSize:        uint32(unsafe.Sizeof(WNDCLASSEX{})),
		LpfnWndProc:   syscall.NewCallback(flashWndProc),
		HInstance:     hinst,
		HCursor:       hCursor,
		HbrBackground: blackBrush,
		LpszClassName: flashClass,
	}
	pRegisterClassEx.Call(uintptr(unsafe.Pointer(&flashWC)))
	flashWnd, _, _ = pCreateWindowEx.Call(
		WS_EX_LAYERED|WS_EX_TRANSPARENT|WS_EX_TOOLWINDOW|WS_EX_NOACTIVATE,
		uintptr(unsafe.Pointer(flashClass)), uintptr(unsafe.Pointer(w(""))),
		WS_POPUP, 0, 0, 42, 42, 0, 0, hinst, 0,
	)
	if flashWnd != 0 {
		pSetLayeredWindowAttributes.Call(flashWnd, uintptr(rgb(0, 0, 0)), 0, LWA_COLORKEY)
	}

	installGlobalHotkey(hinst)

	pShowWindow.Call(mainWnd, 5)
	pUpdateWindow.Call(mainWnd)

	var msg MSG
	for {
		r, _, _ := pGetMessage.Call(uintptr(unsafe.Pointer(&msg)), 0, 0, 0)
		if int32(r) <= 0 {
			break
		}
		if msg.Message == WM_HOTKEY && int(msg.WParam) == HOTKEY_ID {
			toggleFromHotkey()
			continue
		}
		if (msg.Message == WM_KEYDOWN || msg.Message == WM_SYSKEYDOWN) && msg.WParam == VK_F1 {
			openHelp()
			continue
		}
		pTranslateMessage.Call(uintptr(unsafe.Pointer(&msg)))
		pDispatchMessage.Call(uintptr(unsafe.Pointer(&msg)))
	}
}
