package main

import (
	"syscall"
	"unsafe"
)

func rgb(r, g, b byte) uint32 { return uint32(r) | uint32(g)<<8 | uint32(b)<<16 }

func w(s string) *uint16 {
	p, _ := syscall.UTF16PtrFromString(s)
	return p
}
func loword(v uintptr) uint16 { return uint16(v & 0xffff) }
func hiword(v uintptr) uint16 { return uint16((v >> 16) & 0xffff) }
func uptrI32(v int32) uintptr { return uintptr(uint32(v)) }

func getText(hwnd uintptr) string {
	n, _, _ := pGetWindowTextLength.Call(hwnd)
	b := make([]uint16, n+1)
	if len(b) > 0 {
		pGetWindowText.Call(hwnd, uintptr(unsafe.Pointer(&b[0])), n+1)
	}
	return syscall.UTF16ToString(b)
}
func setText(hwnd uintptr, s string) { pSetWindowText.Call(hwnd, uintptr(unsafe.Pointer(w(s)))) }
func dlg(id int) uintptr {
	r, _, _ := pGetDlgItem.Call(mainWnd, uintptr(id))
	return r
}
func comboIndex(id int) int {
	r, _, _ := pSendMessage.Call(dlg(id), CB_GETCURSEL, 0, 0)
	return int(r)
}

func createCtl(class, text string, style uint32, x, y, cx, cy, id int) uintptr {
	h, _, _ := pCreateWindowEx.Call(
		0,
		uintptr(unsafe.Pointer(w(class))),
		uintptr(unsafe.Pointer(w(text))),
		uintptr(style|WS_CHILD|WS_VISIBLE),
		uintptr(x), uintptr(y), uintptr(cx), uintptr(cy),
		mainWnd, uintptr(id), 0, 0,
	)
	if h != 0 {
		controls = append(controls, h)
	}
	return h
}
func label(text string, x, y, cx, cy int, id int) uintptr {
	return createCtl("STATIC", text, 0, x, y, cx, cy, id)
}
func centerLabel(text string, x, y, cx, cy int, id int) uintptr {
	return createCtl("STATIC", text, WS_BORDER|SS_CENTER|SS_CENTERIMAGE, x, y, cx, cy, id)
}
func edit(text string, x, y, cx, cy, id int) uintptr {
	return createCtl("EDIT", text, WS_BORDER|WS_TABSTOP|ES_NUMBER|ES_CENTER, x, y, cx, cy, id)
}
func button(text string, x, y, cx, cy, id int, ownerDraw bool) uintptr {
	style := uint32(WS_TABSTOP | BS_PUSHBUTTON)
	if ownerDraw {
		style = WS_TABSTOP | BS_OWNERDRAW
	}
	return createCtl("BUTTON", text, style, x, y, cx, cy, id)
}
func group(text string, x, y, cx, cy int) uintptr {
	h := createCtl("BUTTON", text, BS_GROUPBOX, x, y, cx, cy, 0)
	return h
}
func combo(items []string, x, y, cx, id, sel int) uintptr {
	h := createCtl("COMBOBOX", "", WS_TABSTOP|CBS_DROPDOWNLIST, x, y, cx, 200, id)
	for _, s := range items {
		pSendMessage.Call(h, CB_ADDSTRING, 0, uintptr(unsafe.Pointer(w(s))))
	}
	pSendMessage.Call(h, CB_SETCURSEL, uintptr(sel), 0)
	return h
}

func numericField(text string, x, y, cx, cy, id, min, max int) uintptr {
	e := edit(text, x, y, cx, cy, id)
	u := createCtl("msctls_updown32", "", UDS_SETBUDDYINT|UDS_ALIGNRIGHT|UDS_ARROWKEYS|UDS_NOTHOUSANDS|UDS_HOTTRACK,
		x+cx-22, y, 22, cy, 0)
	if u != 0 {
		pSendMessage.Call(u, UDM_SETBUDDY, e, 0)
		pSendMessage.Call(u, UDM_SETRANGE32, uintptr(min), uintptr(max))
	}
	return e
}

func applyFont(hwnd, font uintptr) {
	if hwnd != 0 && font != 0 {
		pSendMessage.Call(hwnd, WM_SETFONT, font, 1)
	}
}

func makeFont(height int32, weight int, face string) uintptr {
	f, _, _ := pCreateFont.Call(
		uptrI32(-height), 0, 0, 0, uintptr(weight),
		0, 0, 0, 1, 0, 0, 5, 0,
		uintptr(unsafe.Pointer(w(face))),
	)
	return f
}
