package main

import (
	"fmt"
	"runtime"

	"github.com/axionat/eagle/piglet/glutils"
	"github.com/axionat/glfw/v3.2/glfw"
	"github.com/go-gl/gl/v4.1-core/gl"
)

func init() {
	runtime.LockOSThread()
}

func main() {
	fmt.Println("working...")
	err := glfw.Init()

	if err != nil {
		panic(err)
	}

	window, err := glfw.CreateWindow(640, 480, "Testing", nil, nil)

	if err != nil {
		panic(err)
	}

	window.MakeContextCurrent()

	err = gl.Init()
	if err != nil {
		panic("OpenGL initialization failed: " + err.Error())
	}
	gl.Enable(gl.BLEND)
	gl.BlendFunc(gl.SRC_ALPHA, gl.ONE_MINUS_SRC_ALPHA)
	glutils.CheckGlError()

	//	window.ShowToolsWindow()

	for !window.ShouldClose() {
		window.SwapBuffers()
		glfw.PollEvents()
	}

	glfw.Terminate()
	fmt.Println("Exiting normally...")
}
