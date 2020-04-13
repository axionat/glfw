package main

import (
	"github.com/axionat/glfw/v3.2/glfw"
	"runtime"
)

func init() {
	runtime.LockOSThread()
}

func main() {
	err := glfw.Init()

	if err != nil {
		panic(err)
	}

	window, err := glfw.CreateWindow(640, 480, "Testing", nil, nil)

	if err != nil {
		panic(err)
	}

	window.MakeContextCurrent()
	window.ShowToolsWindow()

	for !window.ShouldClose() {
		window.SwapBuffers()
		glfw.PollEvents()
	}

	glfw.Terminate()
}
