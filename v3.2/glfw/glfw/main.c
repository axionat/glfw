#ifndef _cgo

#include <stdio.h>
#include <stdlib.h>

#define _GLFW_WIN32

#include "src/context.c"
#include "src/egl_context.c"
#include "src/init.c"
#include "src/input.c"
#include "src/monitor.c"
#include "src/vulkan.c"
#include "src/wgl_context.c"
#include "src/win32_init.c"
#include "src/win32_joystick.c"
#include "src/win32_monitor.c"
#include "src/win32_time.c"
#include "src/win32_tls.c"
#include "src/win32_window.c"
#include "src/window.c"

#define TEST_CODE 33

void goMenuCallback(_GLFWwindow *window, int code) {
    char narrowString[] = "Cerrar sesión";
    int length = MultiByteToWideChar(CP_UTF8, 0, narrowString, -1, NULL, 0);
    wchar_t wideString[length];
    MultiByteToWideChar(CP_UTF8, 0, narrowString, -1, wideString, length);

    if (code == TEST_CODE) {
        MessageBox(NULL, wideString, L"Menu", MB_OK);
    }
}

BOOL goContextualMenuCallback(_GLFWwindow *window, long x, long y) {
    HMENU hMenu = CreatePopupMenu();
    AppendMenuW(hMenu, MF_STRING, 33, L"&New");
    AppendMenuW(hMenu, MF_STRING, 33, L"&Open");
    AppendMenuW(hMenu, MF_STRING, 33, L"&Quit");
    return showAndDestroyContextualMenu(hMenu, window->win32.handle, x, y);
}


static void error_callback(int error, const char *description) {
    fprintf(stderr, "Error: %s\n", description);
}

static void key_callback(GLFWwindow *window, int key, int scancode, int action, int mods) {
    if (key == GLFW_KEY_ESCAPE && action == GLFW_PRESS)
        glfwSetWindowShouldClose(window, GLFW_TRUE);
}

int main() {
    glfwSetErrorCallback(error_callback);
    if (!glfwInit()) exit(EXIT_FAILURE);
    glfwWindowHint(GLFW_CONTEXT_VERSION_MAJOR, 2);
    glfwWindowHint(GLFW_CONTEXT_VERSION_MINOR, 0);

    GLFWwindow *window = glfwCreateWindow(1200, 900, "Example", NULL, NULL);

    if (!window) {
        glfwTerminate();
        exit(EXIT_FAILURE);
    }

    HMENU mainMenu = CreateMenu();
    appendMenu(mainMenu, TEST_CODE, "Sesión");
    SetMenu(glfwGetWin32Window(window), mainMenu);

    glfwSetKeyCallback(window, key_callback);
    glfwMakeContextCurrent(window);
    glfwSwapInterval(1);

    while (!glfwWindowShouldClose(window)) {
        glfwSwapBuffers(window);
        glfwPollEvents();
    }

    glfwDestroyWindow(window);
    glfwTerminate();
    exit(EXIT_SUCCESS);
}

#endif
