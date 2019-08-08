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

void goMenuCallback(int code) {
    if (code == 33)
        MessageBox(NULL, TEXT("Hello, world!"), TEXT("Menu"), MB_OK);
}

static void error_callback(int error, const char* description) {
    fprintf(stderr, "Error: %s\n", description);
}

static void key_callback(GLFWwindow* window, int key, int scancode, int action, int mods) {
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

    HMENU mainMenu = createMenu();
    appendMenu(mainMenu, 33, TEXT("Push"));
    setMainMenu(glfwGetWin32Window(window), mainMenu);

    glfwSetKeyCallback(window, key_callback);
    glfwMakeContextCurrent(window);
    glfwSwapInterval(1);

    while (!glfwWindowShouldClose(window)) {
        //int width, height;
        //glfwGetFramebufferSize(window, &width, &height);
        //glViewport(0, 0, width, height);
        //glClearColor(1.f, 0.f, 0.f, 1.f);
        //glClear(GL_COLOR_BUFFER_BIT);
        glfwSwapBuffers(window);
        glfwPollEvents();
    }

    glfwDestroyWindow(window);
    glfwTerminate();
    exit(EXIT_SUCCESS);
}

#endif
