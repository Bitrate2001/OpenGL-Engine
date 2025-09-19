package camera

import (
	"math"

	"github.com/go-gl/glfw/v3.3/glfw"
	"github.com/go-gl/mathgl/mgl32"
)

type Camera struct {
	Position     mgl32.Vec3
	Front        mgl32.Vec3
	Up           mgl32.Vec3
	Yaw, Pitch   float32
	Sensitivity  float32
	Speed        float32
	firstMouse   bool
	lastX, lastY float64
}

func New(pos, front, up mgl32.Vec3) *Camera {
	return &Camera{
		Position:    pos,
		Front:       front,
		Up:          up,
		Yaw:         -90,
		Pitch:       0.0,
		Sensitivity: 0.1,
		Speed:       0.05,
		firstMouse:  true,
		lastX:       400,
		lastY:       300,
	}
}

func (c *Camera) GetViewMatrix() mgl32.Mat4 {
	return mgl32.LookAtV(c.Position, c.Position.Add(c.Front), c.Up)
}

// updateFront ricalcola il vettore Front da yaw/pitch
func (c *Camera) updateFront() {
	front := mgl32.Vec3{
		float32(math.Cos(float64(mgl32.DegToRad(c.Yaw))) * math.Cos(float64(mgl32.DegToRad(c.Pitch)))),
		float32(math.Sin(float64(mgl32.DegToRad(c.Pitch)))),
		float32(math.Sin(float64(mgl32.DegToRad(c.Yaw))) * math.Cos(float64(mgl32.DegToRad(c.Pitch)))),
	}
	c.Front = front.Normalize()
}

// AttachMouse aggiunge il callback per il movimento del mouse
func (c *Camera) AttachMouse(window *glfw.Window) {
	window.SetInputMode(glfw.CursorMode, glfw.CursorDisabled)

	window.SetCursorPosCallback(func(w *glfw.Window, xpos float64, ypos float64) {
		if c.firstMouse {
			c.lastX, c.lastY = xpos, ypos
			c.firstMouse = false
		}

		xoffset := float32(xpos-c.lastX) * c.Sensitivity
		yoffset := float32(c.lastY-ypos) * c.Sensitivity // inverti Y
		c.lastX, c.lastY = xpos, ypos

		c.Yaw += xoffset
		c.Pitch += yoffset

		if c.Pitch > 89.0 {
			c.Pitch = 89.0
		}
		if c.Pitch < -89.0 {
			c.Pitch = -89.0
		}

		c.updateFront()
	})
}
