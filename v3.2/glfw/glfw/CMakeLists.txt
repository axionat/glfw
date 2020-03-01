cmake_minimum_required(VERSION 3.13)
project(glfw C)

set(CMAKE_C_STANDARD 99)

add_executable(glfw WIN32 main.c)
add_executable(test WIN32 test.c)