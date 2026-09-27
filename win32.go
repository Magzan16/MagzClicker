package main

import "syscall"

const (
	appName = "MagzClicker"
	version = "1.5.0"
	helpURL = "https://github.com/Magzan16/MagzClicker#readme"

	WS_OVERLAPPED  = 0x00000000
	WS_CAPTION     = 0x00C00000
	WS_SYSMENU     = 0x00080000
	WS_MINIMIZEBOX = 0x00020000
	WS_CHILD       = 0x40000000
	WS_VISIBLE     = 0x10000000
	WS_TABSTOP     = 0x00010000
	WS_BORDER      = 0x00800000
	WS_POPUP       = 0x80000000

	WS_EX_TRANSPARENT = 0x00000020
	WS_EX_TOOLWINDOW  = 0x00000080
	WS_EX_LAYERED     = 0x00080000
	WS_EX_NOACTIVATE  = 0x08000000

	BS_PUSHBUTTON = 0x00000000
	BS_GROUPBOX   = 0x00000007
	BS_OWNERDRAW  = 0x0000000B

	CBS_DROPDOWNLIST = 0x0003
	ES_NUMBER        = 0x2000
	ES_CENTER        = 0x0001

	SS_CENTER          = 0x00000001
	SS_ICON            = 0x00000003
	SS_CENTERIMAGE     = 0x00000200
	SS_REALSIZECONTROL = 0x00000040

	CW_USEDEFAULT = 0x80000000

	WM_DESTROY        = 0x0002
	WM_PAINT          = 0x000F
	WM_CLOSE          = 0x0010
	WM_TIMER          = 0x0113
	WM_COMMAND        = 0x0111
	WM_HOTKEY         = 0x0312
	WM_DRAWITEM       = 0x002B
	WM_SETFONT        = 0x0030
	WM_NCHITTEST      = 0x0084
	WM_CTLCOLORBTN    = 0x0135
	WM_CTLCOLORSTATIC = 0x0138
	WM_USER           = 0x0400
	WM_APP            = 0x8000

	WM_APP_COUNTER    = WM_APP + 1
	WM_APP_TOGGLE     = WM_APP + 2
	WM_APP_CLICKFLASH = WM_APP + 3
	WM_APP_FINISHED   = WM_APP + 4

	BN_CLICKED = 0
	EN_CHANGE  = 0x0300

	CB_ADDSTRING = 0x0143
	CB_SETCURSEL = 0x014E
	CB_GETCURSEL = 0x0147

	STM_SETIMAGE = 0x0172
	IMAGE_ICON   = 1

	VK_F1       = 0x70
	VK_CONTROL  = 0x11
	VK_MENU     = 0x12
	VK_S        = 0x53
	VK_LCONTROL = 0xA2
	VK_RCONTROL = 0xA3
	VK_LMENU    = 0xA4
	VK_RMENU    = 0xA5

	MOD_ALT               = 0x0001
	MOD_CONTROL           = 0x0002
	MOD_NOREPEAT          = 0x4000
	HOTKEY_ID             = 1
	HOTKEY_POLL_TIMER     = 9001
	HOTKEY_FEEDBACK_TIMER = 9002

	WH_KEYBOARD_LL = 13
	HC_ACTION      = 0
	WM_KEYDOWN     = 0x0100
	WM_KEYUP       = 0x0101
	WM_SYSKEYDOWN  = 0x0104
	WM_SYSKEYUP    = 0x0105

	MOUSEEVENTF_LEFTDOWN   = 0x0002
	MOUSEEVENTF_LEFTUP     = 0x0004
	MOUSEEVENTF_RIGHTDOWN  = 0x0008
	MOUSEEVENTF_RIGHTUP    = 0x0010
	MOUSEEVENTF_MIDDLEDOWN = 0x0020
	MOUSEEVENTF_MIDDLEUP   = 0x0040
	INPUT_MOUSE            = 0

	MB_OK        = 0x00000000
	MB_ICONERROR       = 0x00000010
	MB_ICONINFORMATION = 0x00000040

	COLOR_WINDOW = 5

	DEFAULT_GUI_FONT = 17
	NULL_BRUSH       = 5
	IDC_ARROW        = 32512

	LWA_COLORKEY      = 0x00000001
	SW_HIDE           = 0
	SW_SHOWNORMAL     = 1
	SW_SHOWNOACTIVATE = 4
	SWP_NOACTIVATE    = 0x0010

	ICC_UPDOWN_CLASS = 0x00000010
	UDS_SETBUDDYINT  = 0x0002
	UDS_ALIGNRIGHT   = 0x0004
	UDS_ARROWKEYS    = 0x0020
	UDS_NOTHOUSANDS  = 0x0080
	UDS_HOTTRACK     = 0x0100
	UDM_SETBUDDY     = WM_USER + 105
	UDM_SETRANGE32   = WM_USER + 111

	ODS_SELECTED  = 0x0001
	DT_CENTER     = 0x00000001
	DT_VCENTER    = 0x00000004
	DT_SINGLELINE = 0x00000020
	TRANSPARENT   = 1
)

const (
	idHours        = 101
	idMinutes      = 102
	idSeconds      = 103
	idMillis       = 104
	idButton       = 105
	idClickType    = 106
	idStartStop    = 107
	idExit         = 108
	idStatus       = 109
	idPerformed    = 110
	idLogo         = 111
	idHotkey       = 112
	idModeTitle    = 113
	idModeSub      = 114
	idTitleMagz    = 115
	idTitleClicker = 116
	idClickLimit   = 117
	idEstimate     = 118
	idAbout        = 119
)

type POINT struct{ X, Y int32 }
type RECT struct{ Left, Top, Right, Bottom int32 }
type MSG struct {
	Hwnd           uintptr
	Message        uint32
	WParam, LParam uintptr
	Time           uint32
	Pt             POINT
	LPrivate       uint32
}
type KBDLLHOOKSTRUCT struct {
	VkCode      uint32
	ScanCode    uint32
	Flags       uint32
	Time        uint32
	DwExtraInfo uintptr
}
type WNDCLASSEX struct {
	CbSize        uint32
	Style         uint32
	LpfnWndProc   uintptr
	CbClsExtra    int32
	CbWndExtra    int32
	HInstance     uintptr
	HIcon         uintptr
	HCursor       uintptr
	HbrBackground uintptr
	LpszMenuName  *uint16
	LpszClassName *uint16
	HIconSm       uintptr
}
type MOUSEINPUT struct {
	Dx, Dy                   int32
	MouseData, DwFlags, Time uint32
	DwExtraInfo              uintptr
}
type INPUT struct {
	Type uint32
	_    uint32 // alignment for the union on 64-bit Windows
	Mi   MOUSEINPUT
}
type INITCOMMONCONTROLSEX struct {
	DwSize uint32
	DwICC  uint32
}
type DRAWITEMSTRUCT struct {
	CtlType    uint32
	CtlID      uint32
	ItemID     uint32
	ItemAction uint32
	ItemState  uint32
	HwndItem   uintptr
	HDC        uintptr
	RcItem     RECT
	ItemData   uintptr
}
type PAINTSTRUCT struct {
	HDC         uintptr
	FErase      int32
	RcPaint     RECT
	FRestore    int32
	FIncUpdate  int32
	RgbReserved [32]byte
}

var (
	user32   = syscall.NewLazyDLL("user32.dll")
	kernel32 = syscall.NewLazyDLL("kernel32.dll")
	gdi32    = syscall.NewLazyDLL("gdi32.dll")
	comctl32 = syscall.NewLazyDLL("comctl32.dll")
	shell32  = syscall.NewLazyDLL("shell32.dll")

	pRegisterClassEx            = user32.NewProc("RegisterClassExW")
	pCreateWindowEx             = user32.NewProc("CreateWindowExW")
	pDefWindowProc              = user32.NewProc("DefWindowProcW")
	pShowWindow                 = user32.NewProc("ShowWindow")
	pUpdateWindow               = user32.NewProc("UpdateWindow")
	pGetMessage                 = user32.NewProc("GetMessageW")
	pTranslateMessage           = user32.NewProc("TranslateMessage")
	pDispatchMessage            = user32.NewProc("DispatchMessageW")
	pPostQuitMessage            = user32.NewProc("PostQuitMessage")
	pPostMessage                = user32.NewProc("PostMessageW")
	pSendMessage                = user32.NewProc("SendMessageW")
	pSetWindowText              = user32.NewProc("SetWindowTextW")
	pGetWindowText              = user32.NewProc("GetWindowTextW")
	pGetWindowTextLength        = user32.NewProc("GetWindowTextLengthW")
	pGetDlgItem                 = user32.NewProc("GetDlgItem")
	pGetAsyncKeyState           = user32.NewProc("GetAsyncKeyState")
	pRegisterHotKey             = user32.NewProc("RegisterHotKey")
	pUnregisterHotKey           = user32.NewProc("UnregisterHotKey")
	pSetWindowsHookEx           = user32.NewProc("SetWindowsHookExW")
	pCallNextHookEx             = user32.NewProc("CallNextHookEx")
	pUnhookWindowsHookEx        = user32.NewProc("UnhookWindowsHookEx")
	pGetCursorPos               = user32.NewProc("GetCursorPos")
	pLoadImage                  = user32.NewProc("LoadImageW")
	pLoadCursor                 = user32.NewProc("LoadCursorW")
	pSendInput                  = user32.NewProc("SendInput")
	pMessageBox                 = user32.NewProc("MessageBoxW")
	pDestroyWindow              = user32.NewProc("DestroyWindow")
	pSetProcessDPIAware         = user32.NewProc("SetProcessDPIAware")
	pInvalidateRect             = user32.NewProc("InvalidateRect")
	pFillRect                   = user32.NewProc("FillRect")
	pDrawText                   = user32.NewProc("DrawTextW")
	pSetLayeredWindowAttributes = user32.NewProc("SetLayeredWindowAttributes")
	pSetWindowPos               = user32.NewProc("SetWindowPos")
	pSetTimer                   = user32.NewProc("SetTimer")
	pKillTimer                  = user32.NewProc("KillTimer")
	pBeginPaint                 = user32.NewProc("BeginPaint")
	pEndPaint                   = user32.NewProc("EndPaint")

	pGetStockObject   = gdi32.NewProc("GetStockObject")
	pCreateFont       = gdi32.NewProc("CreateFontW")
	pDeleteObject     = gdi32.NewProc("DeleteObject")
	pCreateSolidBrush = gdi32.NewProc("CreateSolidBrush")
	pCreatePen        = gdi32.NewProc("CreatePen")
	pSelectObject     = gdi32.NewProc("SelectObject")
	pRoundRect        = gdi32.NewProc("RoundRect")
	pEllipse          = gdi32.NewProc("Ellipse")
	pSetBkMode        = gdi32.NewProc("SetBkMode")
	pSetTextColor     = gdi32.NewProc("SetTextColor")
	pGetModuleHandle  = kernel32.NewProc("GetModuleHandleW")

	pInitCommonControlsEx = comctl32.NewProc("InitCommonControlsEx")
	pShellExecute         = shell32.NewProc("ShellExecuteW")
)

var (
	mainWnd            uintptr
	flashWnd           uintptr
	running            int32
	stopCh             chan struct{}
	performed          uint64
	hotkeyRegistered   bool
	keyboardHook       uintptr
	hookCallback       uintptr
	hookHotkeyLatched  int32
	pollHotkeyLatched  int32
	lastHotkeyToggleNS int64

	regularFont uintptr
	inputFont   uintptr
	smallFont   uintptr
	sectionFont uintptr
	titleFont   uintptr
	hotkeyFont  uintptr
	buttonFont  uintptr

	whiteBrush     uintptr
	lightBlueBrush uintptr
	blackBrush     uintptr

	controls   []uintptr
	textColors = map[uintptr]uint32{}
	bgBrushes  = map[uintptr]uintptr{}
)
