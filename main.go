package main

import "C"
import (
	"OpenGL/camera"
	"OpenGL/shaders"
	"fmt"
	"runtime"

	"github.com/go-gl/gl/v4.1-core/gl"
	"github.com/go-gl/glfw/v3.3/glfw"
	"github.com/go-gl/mathgl/mgl32"
)

const (
	width  = 1024
	height = 720
)

func init() {
	// OpenGL deve girare nel main thread
	runtime.LockOSThread()
}

func main() {
	// Inizializza GLFW
	if err := glfw.Init(); err != nil {
		panic(err)
	}
	defer glfw.Terminate()

	glfw.WindowHint(glfw.Resizable, glfw.False)
	glfw.WindowHint(glfw.ContextVersionMajor, 4)
	glfw.WindowHint(glfw.ContextVersionMinor, 1)
	glfw.WindowHint(glfw.OpenGLProfile, glfw.OpenGLCoreProfile)
	glfw.WindowHint(glfw.OpenGLForwardCompatible, glfw.True)

	window, err := glfw.CreateWindow(width, height, "Rotating Cube", nil, nil)
	if err != nil {
		panic(err)
	}
	window.MakeContextCurrent()

	// Inizializza OpenGL
	if err := gl.Init(); err != nil {
		panic(err)
	}
	fmt.Println("OpenGL", gl.GoStr(gl.GetString(gl.VERSION)))

	// Definiamo i vertici del cubo (posizione + colore)
	vertices := []float32{
		// posizioni        // colori (RGB)
		-0.5, -0.5, -0.5, 1, 0, 0,
		0.5, -0.5, -0.5, 0, 1, 0,
		0.5, 0.5, -0.5, 0, 0, 1,
		-0.5, 0.5, -0.5, 1, 1, 0,
		-0.5, -0.5, 0.5, 1, 0, 1,
		0.5, -0.5, 0.5, 0, 1, 1,
		0.5, 0.5, 0.5, 1, 1, 1,
		-0.5, 0.5, 0.5, 0, 0, 0,
	}
	indices := []uint32{
		0, 1, 2, 2, 3, 0,
		4, 5, 6, 6, 7, 4,
		0, 1, 5, 5, 4, 0,
		2, 3, 7, 7, 6, 2,
		0, 3, 7, 7, 4, 0,
		1, 2, 6, 6, 5, 1,
	}

	// VAO/VBO/EBO
	var vao, vbo, ebo uint32
	gl.GenVertexArrays(1, &vao)
	gl.GenBuffers(1, &vbo)
	gl.GenBuffers(1, &ebo)

	gl.BindVertexArray(vao)

	gl.BindBuffer(gl.ARRAY_BUFFER, vbo)
	gl.BufferData(gl.ARRAY_BUFFER, len(vertices)*4, gl.Ptr(vertices), gl.STATIC_DRAW)

	gl.BindBuffer(gl.ELEMENT_ARRAY_BUFFER, ebo)
	gl.BufferData(gl.ELEMENT_ARRAY_BUFFER, len(indices)*4, gl.Ptr(indices), gl.STATIC_DRAW)

	// Posizione (location=0), 3 float
	gl.EnableVertexAttribArray(0)
	gl.VertexAttribPointer(0, 3, gl.FLOAT, false, 6*4, gl.PtrOffset(0))
	// Colore (location=1), 3 float
	gl.EnableVertexAttribArray(1)
	gl.VertexAttribPointer(1, 3, gl.FLOAT, false, 6*4, gl.PtrOffset(3*4))

	// Shader program
	program := shaders.NewProgram(shaders.VertexShaderSource, shaders.FragmentShaderSource)
	gl.UseProgram(program)

	// Uniform locations
	modelLoc := gl.GetUniformLocation(program, gl.Str("model\x00"))
	viewLoc := gl.GetUniformLocation(program, gl.Str("view\x00"))
	projLoc := gl.GetUniformLocation(program, gl.Str("projection\x00"))

	fbWidth, fbHeight := window.GetFramebufferSize()
	gl.Viewport(0, 0, int32(fbWidth), int32(fbHeight))

	proj := mgl32.Perspective(mgl32.DegToRad(45), float32(fbWidth)/float32(fbHeight), 0.1, 100.0)

	cam := camera.New(
		mgl32.Vec3{0, 0, 5},  // inizialmente 5 unità davanti
		mgl32.Vec3{0, 0, -1}, // direzione iniziale
		mgl32.Vec3{0, 1, 0},  // asse Y come "up"
	)

	cam.AttachMouse(window)

	gl.Enable(gl.DEPTH_TEST)

	// Loop principale
	for !window.ShouldClose() {
		gl.ClearColor(0.2, 0.3, 0.3, 1.0)
		gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)

		cam.ProcessKeyboard(window)

		// Matrici di base
		view := cam.GetViewMatrix() // sposta la camera
		gl.UniformMatrix4fv(viewLoc, 1, false, &view[0])

		model := mgl32.HomogRotate3D(0, mgl32.Vec3{1.0, 1.0, 0.0}.Normalize())

		gl.UniformMatrix4fv(modelLoc, 1, false, &model[0])
		gl.UniformMatrix4fv(viewLoc, 1, false, &view[0])
		gl.UniformMatrix4fv(projLoc, 1, false, &proj[0])

		// Draw
		gl.BindVertexArray(vao)
		gl.DrawElements(gl.TRIANGLES, int32(len(indices)), gl.UNSIGNED_INT, gl.PtrOffset(0))

		window.SwapBuffers()
		glfw.PollEvents()
	}
}
