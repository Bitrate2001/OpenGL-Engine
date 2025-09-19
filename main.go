package main

import (
	"OpenGL/shaders"
	"fmt"
	"math"
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
	program := shaders.NewProgram(vertexShaderSource, fragmentShaderSource)
	gl.UseProgram(program)

	// Uniform locations
	modelLoc := gl.GetUniformLocation(program, gl.Str("model\x00"))
	viewLoc := gl.GetUniformLocation(program, gl.Str("view\x00"))
	projLoc := gl.GetUniformLocation(program, gl.Str("projection\x00"))

	fbWidth, fbHeight := window.GetFramebufferSize()
	gl.Viewport(0, 0, int32(fbWidth), int32(fbHeight))

	proj := mgl32.Perspective(mgl32.DegToRad(45), float32(fbWidth)/float32(fbHeight), 0.1, 100.0)

	cameraPos := mgl32.Vec3{0, 0, 5}    // inizialmente 5 unità davanti
	cameraFront := mgl32.Vec3{0, 0, -1} // direzione iniziale
	cameraUp := mgl32.Vec3{0, 1, 0}     // asse Y come "up"

	var yaw float32 = -90.0 // guardo lungo -Z
	var pitch float32 = 0.0
	var lastX float64 = 400
	var lastY float64 = 300
	var firstMouse = true

	var sensitivity float32 = 0.1

	window.SetInputMode(glfw.CursorMode, glfw.CursorDisabled)

	window.SetCursorPosCallback(func(w *glfw.Window, xpos float64, ypos float64) {
		if firstMouse {
			lastX = xpos
			lastY = ypos
			firstMouse = false
		}

		xoffset := float32(xpos-lastX) * sensitivity
		yoffset := float32(lastY-ypos) * sensitivity // inverti Y
		lastX = xpos
		lastY = ypos

		yaw += xoffset
		pitch += yoffset

		if pitch > 89.0 {
			pitch = 89.0
		}
		if pitch < -89.0 {
			pitch = -89.0
		}

		front := mgl32.Vec3{
			float32(math.Cos(float64(mgl32.DegToRad(yaw))) * math.Cos(float64(mgl32.DegToRad(pitch)))),
			float32(math.Sin(float64(mgl32.DegToRad(pitch)))),
			float32(math.Sin(float64(mgl32.DegToRad(yaw))) * math.Cos(float64(mgl32.DegToRad(pitch)))),
		}
		cameraFront = front.Normalize()
	})

	gl.Enable(gl.DEPTH_TEST)

	//start := time.Now()
	//rotationPeriod := 3.0 // secondi per giro completo

	// Loop principale
	for !window.ShouldClose() {
		gl.ClearColor(0.2, 0.3, 0.3, 1.0)
		gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)

		if window.GetKey(glfw.KeyW) == glfw.Press {
			cameraPos = cameraPos.Add(cameraFront.Mul(0.05)) // muovi avanti
		}
		if window.GetKey(glfw.KeyS) == glfw.Press {
			cameraPos = cameraPos.Sub(cameraFront.Mul(0.05)) // muovi indietro
		}
		if window.GetKey(glfw.KeyA) == glfw.Press {
			cameraPos = cameraPos.Sub(cameraFront.Cross(cameraUp).Normalize().Mul(0.05)) // sinistra
		}
		if window.GetKey(glfw.KeyD) == glfw.Press {
			cameraPos = cameraPos.Add(cameraFront.Cross(cameraUp).Normalize().Mul(0.05)) // destra
		}
		if window.GetKey(glfw.KeyZ) == glfw.Press {
			cameraPos = cameraPos.Sub(cameraUp.Mul(0.05))
		}
		if window.GetKey(glfw.KeyX) == glfw.Press {
			cameraPos = cameraPos.Add(cameraUp.Mul(0.05))
		}
		if window.GetKey(glfw.KeyEscape) == glfw.Press {
			window.SetShouldClose(true)
		}

		// Matrici di base
		view := mgl32.LookAtV(cameraPos, cameraPos.Add(cameraFront), cameraUp) // sposta la camera
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

var vertexShaderSource = `
#version 410 core
layout (location = 0) in vec3 aPos;
layout (location = 1) in vec3 aColor;

out vec3 ourColor;

uniform mat4 model;
uniform mat4 view;
uniform mat4 projection;

void main() {
	gl_Position = projection * view * model * vec4(aPos, 1.0);
	ourColor = aColor;
}
` + "\x00"

var fragmentShaderSource = `
#version 410 core
in vec3 ourColor;
out vec4 FragColor;

void main() {
	FragColor = vec4(ourColor, 1.0);
}
` + "\x00"
