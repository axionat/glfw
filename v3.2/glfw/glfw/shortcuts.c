#ifndef _cgo

#include <Windows.h>
#include <CommCtrl.h>

#pragma comment(lib, "comctl32.lib")

HWND CreateListView(HWND hwndParent, HINSTANCE hInstance) {
    INITCOMMONCONTROLSEX icex; // Structure for control initialization.
    icex.dwICC = ICC_LISTVIEW_CLASSES;
    InitCommonControlsEx(&icex);
    RECT rcClient; // The parent window's client area.
    GetClientRect(hwndParent, &rcClient);

    // Create the list-view window in report view with label editing enabled.
    HWND hWndListView = CreateWindow(WC_LISTVIEW, (LPCSTR) "", WS_CHILD | LVS_REPORT | LVS_EDITLABELS,
                                     0, 0, rcClient.right - rcClient.left, rcClient.bottom - rcClient.top,
                                     hwndParent, NULL, hInstance, NULL);

    return (hWndListView);
}

VOID SetView(HWND hWndListView, DWORD dwView) {
    // Retrieve the current window style.
    LONG dwStyle = GetWindowLong(hWndListView, GWL_STYLE);

    // Set the window style only if the view bits changed.
    if ((dwStyle & LVS_TYPEMASK) != dwView)
        SetWindowLong(hWndListView, GWL_STYLE, (dwStyle & ~LVS_TYPEMASK) | dwView);
}

BOOL window1closed = FALSE;
BOOL window2closed = FALSE;

LRESULT CALLBACK WindowProcess1(HWND handle, UINT msg, WPARAM wParam, LPARAM lParam) {
    switch (msg) {
        case WM_DESTROY:
            MessageBox(NULL, (LPCSTR) "Window 1 closed", (LPCSTR) "Message", MB_ICONINFORMATION);
            window1closed = TRUE;
            return 0;
        default:
            return DefWindowProc(handle, msg, wParam, lParam);
    }
}

LRESULT CALLBACK WindowProcess2(HWND handle, UINT msg, WPARAM wParam, LPARAM lParam) {
    switch (msg) {
        case WM_DESTROY:
            MessageBox(NULL, (LPCSTR) "Window 2 Closed", (LPCSTR) "Message", MB_ICONINFORMATION);
            window2closed = TRUE;
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

    ShowWindow(handleforwindow1, nShowCmd);
    ShowWindow(handleforwindow2, nShowCmd);
    //SetParent(handleforwindow2, handleforwindow1);

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

    MessageBox(NULL, (LPCSTR) "Both Windows are closed.  Program will now close.", (LPCSTR) "", MB_ICONINFORMATION);

    return 0;
}

#endif
