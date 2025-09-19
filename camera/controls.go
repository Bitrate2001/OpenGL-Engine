package camera

import "github.com/go-gl/glfw/v3.3/glfw"

// ProcessKeyboard gestisce movimento WSAD + su/giù
func (c *Camera) ProcessKeyboard(window *glfw.Window) {
	if window.GetKey(glfw.KeyW) == glfw.Press {
		c.Position = c.Position.Add(c.Front.Mul(c.Speed))
	}
	if window.GetKey(glfw.KeyS) == glfw.Press {
		c.Position = c.Position.Sub(c.Front.Mul(c.Speed))
	}
	if window.GetKey(glfw.KeyA) == glfw.Press {
		c.Position = c.Position.Sub(c.Front.Cross(c.Up).Normalize().Mul(c.Speed))
	}
	if window.GetKey(glfw.KeyD) == glfw.Press {
		c.Position = c.Position.Add(c.Front.Cross(c.Up).Normalize().Mul(c.Speed))
	}
	if window.GetKey(glfw.KeyZ) == glfw.Press {
		c.Position = c.Position.Sub(c.Up.Mul(c.Speed))
	}
	if window.GetKey(glfw.KeyX) == glfw.Press {
		c.Position = c.Position.Add(c.Up.Mul(c.Speed))
	}
	if window.GetKey(glfw.KeyEscape) == glfw.Press {
		window.SetShouldClose(true)
	}
}
