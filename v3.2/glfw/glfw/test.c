#ifndef _cgo

#include <Windows.h>
#include <CommCtrl.h>

#pragma comment(lib, "comctl32.lib")

BOOL window1closed = FALSE;
BOOL window2closed = FALSE;
HWND hwndButton;
HWND hWndListView;

LRESULT CALLBACK WindowProcess1(HWND handle, UINT msg, WPARAM wParam, LPARAM lParam) {
    switch (msg) {
        case WM_DESTROY:
            window1closed = TRUE;
            return 0;
        default:
            return DefWindowProc(handle, msg, wParam, lParam);
    }
}

LRESULT CALLBACK WindowProcess2(HWND handle, UINT msg, WPARAM wParam, LPARAM lParam) {
    switch (msg) {
        case WM_DESTROY:
            window2closed = TRUE;
            return 0;
        default:
            return DefWindowProc(handle, msg, wParam, lParam);
    }
}

LRESULT CALLBACK WindowProcess3(HWND handle, UINT msg, WPARAM wParam, LPARAM lParam) {
    PAINTSTRUCT ps;
    HDC hdc;

    switch (msg) {
        case WM_PAINT:
            hdc = BeginPaint(handle, &ps);
            TextOut(hdc, 40, 40, "Hello, Windows!", 15);
            EndPaint(handle, &ps);
            return 0;
        case WM_COMMAND:
            if (hwndButton == (HWND) lParam)
                MessageBox(NULL, (LPCSTR) "You clicked me.", (LPCSTR) "Error", MB_ICONERROR);
            return 0;
        default:
            return DefWindowProc(handle, msg, wParam, lParam);
    }
}

int WINAPI WinMain(HINSTANCE hInst, HINSTANCE hPrevInst, LPSTR lpCmdLine, int nShowCmd) {
    //create window 1
    WNDCLASSEX windowclassforwindow1;
    ZeroMemory(&windowclassforwindow1, sizeof(WNDCLASSEX));
    windowclassforwindow1.cbSize = sizeof(WNDCLASSEX);
    windowclassforwindow1.hbrBackground = (HBRUSH) COLOR_WINDOW;
    windowclassforwindow1.hCursor = LoadCursor(NULL, IDC_ARROW);
    windowclassforwindow1.hInstance = hInst;
    windowclassforwindow1.lpfnWndProc = (WNDPROC) WindowProcess1;
    windowclassforwindow1.lpszClassName = (LPCSTR) "windowclass 1";
    windowclassforwindow1.style = CS_HREDRAW | CS_VREDRAW;

    if (!RegisterClassEx(&windowclassforwindow1))
        MessageBox(NULL, (LPCSTR) "Window class creation failed", (LPCSTR) "Window Class Failed", MB_ICONERROR);

    HWND handleforwindow1 = CreateWindowEx(0, windowclassforwindow1.lpszClassName,
                                           (LPCSTR) "Parent Window", WS_OVERLAPPEDWINDOW,
                                           200, 150, 640, 480, NULL, NULL, hInst, NULL
    );

    if (!handleforwindow1)
        MessageBox(NULL, (LPCSTR) "Window creation failed", (LPCSTR) "Window Creation Failed", MB_ICONERROR);

    // create window 2
    WNDCLASSEX windowclassforwindow2;
    ZeroMemory(&windowclassforwindow2, sizeof(WNDCLASSEX));
    windowclassforwindow2.cbSize = sizeof(WNDCLASSEX);
    windowclassforwindow2.hbrBackground = (HBRUSH) COLOR_WINDOW;
    windowclassforwindow2.hCursor = LoadCursor(NULL, IDC_ARROW);
    windowclassforwindow2.hInstance = hInst;
    windowclassforwindow2.lpfnWndProc = (WNDPROC) WindowProcess2;
    windowclassforwindow2.lpszClassName = (LPCSTR) "window class2";
    windowclassforwindow2.style = CS_HREDRAW | CS_VREDRAW;

    if (!RegisterClassEx(&windowclassforwindow2))
        MessageBox(NULL, (LPCSTR) "Window class creation failed for window 2", (LPCSTR) "Window Class Failed",
                   MB_ICONERROR);

    HWND handleforwindow2 = CreateWindowEx(0, windowclassforwindow2.lpszClassName,
                                           (LPCSTR) "Child Window", WS_OVERLAPPEDWINDOW,
                                           200, 150, 640, 480, NULL, NULL, hInst, NULL);

    if (!handleforwindow2)
        MessageBox(NULL, (LPCSTR) "Window creation failed", (LPCSTR) "Window Creation Failed", MB_ICONERROR);

    //create window 3
    WNDCLASSEX windowclassforwindow3;
    ZeroMemory(&windowclassforwindow3, sizeof(WNDCLASSEX));
    windowclassforwindow3.cbSize = sizeof(WNDCLASSEX);
    windowclassforwindow3.hbrBackground = CreateSolidBrush(RGB(255, 255, 255));
    windowclassforwindow3.hCursor = LoadCursor(NULL, IDC_ARROW);
    windowclassforwindow3.hInstance = hInst;
    windowclassforwindow3.lpfnWndProc = (WNDPROC) WindowProcess3;
    windowclassforwindow3.lpszClassName = (LPCSTR) "windowclass 3";
    windowclassforwindow3.style = CS_HREDRAW | CS_VREDRAW;

    if (!RegisterClassEx(&windowclassforwindow3))
        MessageBox(NULL, (LPCSTR) "Window class creation failed", (LPCSTR) "Window Class Failed", MB_ICONERROR);

    HWND handleforwindow3 = CreateWindowEx(0, windowclassforwindow3.lpszClassName,
                                           (LPCSTR) "", WS_CHILD | WS_VISIBLE,
                                           0, 0, 400, 300, handleforwindow1, NULL, hInst, NULL);

    if (!handleforwindow3)
        MessageBox(NULL, (LPCSTR) "Window creation failed", (LPCSTR) "Window Creation Failed", MB_ICONERROR);

    // (HINSTANCE) GetWindowLongPtr(handleforwindow3, GWLP_HINSTANCE)
    hwndButton = CreateWindowA(WC_BUTTONA, "Update",
                               WS_VISIBLE | WS_CHILD | WS_TABSTOP | BS_DEFPUSHBUTTON,
                               40, 80, 80, 40,
                               handleforwindow3, NULL, hInst, NULL);

    INITCOMMONCONTROLSEX icex;
    icex.dwICC = ICC_LISTVIEW_CLASSES;
    InitCommonControlsEx(&icex);

    //RECT rcClient;
    //GetClientRect(hwndParent, &rcClient);

    hWndListView = CreateWindowA(WC_LISTVIEWA, "",
                                 WS_CHILD | WS_VISIBLE | LVS_REPORT | LVS_EDITLABELS,
                                 0, 0, 640, 480,
                                 handleforwindow2, NULL, hInst, NULL);

    //LONG dwStyle = GetWindowLong(hWndListView, GWL_STYLE);
    //if ((dwStyle & LVS_TYPEMASK) != dwView) SetWindowLong(hWndListView, GWL_STYLE, (dwStyle & ~LVS_TYPEMASK) | dwView);

    //SetParent(handleforwindow3, handleforwindow1);
    ShowWindow(handleforwindow1, nShowCmd);
    ShowWindow(handleforwindow2, nShowCmd);

    BOOL end = FALSE;
    MSG msg;
    ZeroMemory(&msg, sizeof(MSG));

    while (end == FALSE) {
        if (GetMessage(&msg, NULL, 0, 0)) {
            TranslateMessage(&msg);
            DispatchMessage(&msg);
        }

        end = window1closed && window2closed;
    }

    return 0;
}

#endif
