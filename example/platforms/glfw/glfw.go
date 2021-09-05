package platforms

import (
	"fmt"
	"math"
	"runtime"

	"github.com/go-gl/glfw/v3.2/glfw"
	"github.com/splizard/imgui"
)

const (
	windowWidth  = 1280
	windowHeight = 720

	mouseButtonPrimary   = 0
	mouseButtonSecondary = 1
	mouseButtonTertiary  = 2
	mouseButtonCount     = 3
)

// StringError describes a basic error with static information.
type StringError string

// Error returns the string itself.
func (err StringError) Error() string {
	return string(err)
}

const (
	// ErrUnsupportedClientAPI is used in case the API is not available by the platform.
	ErrUnsupportedClientAPI = StringError("unsupported ClientAPI")
)

// GLFWClientAPI identifies the render system that shall be initialized.
type GLFWClientAPI string

// This is a list of GLFWClientAPI constants.
const (
	GLFWClientAPIOpenGL2 GLFWClientAPI = "OpenGL2"
	GLFWClientAPIOpenGL3 GLFWClientAPI = "OpenGL3"
)

// GLFW implements a platform based on github.com/go-gl/glfw (v3.2).
type GLFW struct {
	imguiIO *imgui.IO

	window *glfw.Window

	time             float64
	mouseJustPressed [3]bool
}

// NewGLFW attempts to initialize a GLFW context.
func NewGLFW(io *imgui.IO, clientAPI GLFWClientAPI) (*GLFW, error) {
	runtime.LockOSThread()

	err := glfw.Init()
	if err != nil {
		return nil, fmt.Errorf("failed to initialize glfw: %w", err)
	}

	switch clientAPI {
	case GLFWClientAPIOpenGL2:
		glfw.WindowHint(glfw.ContextVersionMajor, 2)
		glfw.WindowHint(glfw.ContextVersionMinor, 1)
	case GLFWClientAPIOpenGL3:
		glfw.WindowHint(glfw.ContextVersionMajor, 3)
		glfw.WindowHint(glfw.ContextVersionMinor, 2)
		glfw.WindowHint(glfw.OpenGLProfile, glfw.OpenGLCoreProfile)
		glfw.WindowHint(glfw.OpenGLForwardCompatible, 1)
	default:
		glfw.Terminate()
		return nil, ErrUnsupportedClientAPI
	}

	window, err := glfw.CreateWindow(windowWidth, windowHeight, "ImGui-Go GLFW+"+string(clientAPI)+" example", nil, nil)
	if err != nil {
		glfw.Terminate()
		return nil, fmt.Errorf("failed to create window: %w", err)
	}
	window.MakeContextCurrent()
	glfw.SwapInterval(1)

	platform := &GLFW{
		imguiIO: io,
		window:  window,
	}
	platform.setKeyMapping()
	platform.installCallbacks()

	return platform, nil
}

// Dispose cleans up the resources.
func (platform *GLFW) Dispose() {
	platform.window.Destroy()
	glfw.Terminate()
}

// ShouldStop returns true if the window is to be closed.
func (platform *GLFW) ShouldStop() bool {
	return platform.window.ShouldClose()
}

// ProcessEvents handles all pending window events.
func (platform *GLFW) ProcessEvents() {
	glfw.PollEvents()
}

// DisplaySize returns the dimension of the display.
func (platform *GLFW) DisplaySize() [2]float32 {
	w, h := platform.window.GetSize()
	return [2]float32{float32(w), float32(h)}
}

// FramebufferSize returns the dimension of the framebuffer.
func (platform *GLFW) FramebufferSize() [2]float32 {
	w, h := platform.window.GetFramebufferSize()
	return [2]float32{float32(w), float32(h)}
}

// NewFrame marks the begin of a render pass. It forwards all current state to imgui IO.
func (platform *GLFW) NewFrame() {
	// Setup display size (every frame to accommodate for window resizing)
	displaySize := platform.DisplaySize()
	platform.imguiIO.DisplaySize = imgui.Vec2{displaySize[0], displaySize[1]}

	// Setup time step
	currentTime := glfw.GetTime()
	if platform.time > 0 {
		platform.imguiIO.DeltaTime = float32(currentTime - platform.time)
	}
	platform.time = currentTime

	// Setup inputs
	if platform.window.GetAttrib(glfw.Focused) != 0 {
		x, y := platform.window.GetCursorPos()
		platform.imguiIO.MousePos = imgui.Vec2{float32(x), float32(y)}
	} else {
		platform.imguiIO.MousePos = imgui.Vec2{-math.MaxFloat32, -math.MaxFloat32}
	}

	for i := 0; i < len(platform.mouseJustPressed); i++ {
		down := platform.mouseJustPressed[i] || (platform.window.GetMouseButton(glfwButtonIDByIndex[i]) == glfw.Press)
		platform.imguiIO.MouseDown[i] = down
		platform.mouseJustPressed[i] = false
	}
}

// PostRender performs a buffer swap.
func (platform *GLFW) PostRender() {
	platform.window.SwapBuffers()
}

func (platform *GLFW) setKeyMapping() {
	// Keyboard mapping. ImGui will use those indices to peek into the io.KeysDown[] array.
	platform.imguiIO.KeyMap[imgui.KeyTab] = int(glfw.KeyTab)
	platform.imguiIO.KeyMap[imgui.KeyLeftArrow] = int(glfw.KeyLeft)
	platform.imguiIO.KeyMap[imgui.KeyRightArrow] = int(glfw.KeyRight)
	platform.imguiIO.KeyMap[imgui.KeyUpArrow] = int(glfw.KeyUp)
	platform.imguiIO.KeyMap[imgui.KeyDownArrow] = int(glfw.KeyDown)
	platform.imguiIO.KeyMap[imgui.KeyPageUp] = int(glfw.KeyPageUp)
	platform.imguiIO.KeyMap[imgui.KeyPageDown] = int(glfw.KeyPageDown)
	platform.imguiIO.KeyMap[imgui.KeyHome] = int(glfw.KeyHome)
	platform.imguiIO.KeyMap[imgui.KeyEnd] = int(glfw.KeyEnd)
	platform.imguiIO.KeyMap[imgui.KeyInsert] = int(glfw.KeyInsert)
	platform.imguiIO.KeyMap[imgui.KeyDelete] = int(glfw.KeyDelete)
	platform.imguiIO.KeyMap[imgui.KeyBackspace] = int(glfw.KeyBackspace)
	platform.imguiIO.KeyMap[imgui.KeySpace] = int(glfw.KeySpace)
	platform.imguiIO.KeyMap[imgui.KeyEnter] = int(glfw.KeyEnter)
	platform.imguiIO.KeyMap[imgui.KeyEscape] = int(glfw.KeyEscape)
	platform.imguiIO.KeyMap[imgui.KeyA] = int(glfw.KeyA)
	platform.imguiIO.KeyMap[imgui.KeyC] = int(glfw.KeyC)
	platform.imguiIO.KeyMap[imgui.KeyV] = int(glfw.KeyV)
	platform.imguiIO.KeyMap[imgui.KeyX] = int(glfw.KeyX)
	platform.imguiIO.KeyMap[imgui.KeyY] = int(glfw.KeyY)
	platform.imguiIO.KeyMap[imgui.KeyZ] = int(glfw.KeyZ)
}

func (platform *GLFW) installCallbacks() {
	platform.window.SetMouseButtonCallback(platform.mouseButtonChange)
	platform.window.SetScrollCallback(platform.mouseScrollChange)
	platform.window.SetKeyCallback(platform.keyChange)
	platform.window.SetCharCallback(platform.charChange)
}

var glfwButtonIndexByID = map[glfw.MouseButton]int{
	glfw.MouseButton1: mouseButtonPrimary,
	glfw.MouseButton2: mouseButtonSecondary,
	glfw.MouseButton3: mouseButtonTertiary,
}

var glfwButtonIDByIndex = map[int]glfw.MouseButton{
	mouseButtonPrimary:   glfw.MouseButton1,
	mouseButtonSecondary: glfw.MouseButton2,
	mouseButtonTertiary:  glfw.MouseButton3,
}

func (platform *GLFW) mouseButtonChange(window *glfw.Window, rawButton glfw.MouseButton, action glfw.Action, mods glfw.ModifierKey) {
	buttonIndex, known := glfwButtonIndexByID[rawButton]

	if known && (action == glfw.Press) {
		platform.mouseJustPressed[buttonIndex] = true
	}
}

func (platform *GLFW) mouseScrollChange(window *glfw.Window, x, y float64) {
	platform.imguiIO.MouseWheelH += float32(x)
	platform.imguiIO.MouseWheel += float32(y)
}

func (platform *GLFW) keyChange(window *glfw.Window, key glfw.Key, scancode int, action glfw.Action, mods glfw.ModifierKey) {
	if action == glfw.Press {
		platform.imguiIO.KeysDown[key] = true
	}
	if action == glfw.Release {
		platform.imguiIO.KeysDown[key] = false
	}

	platform.imguiIO.KeyCtrl = (mods & glfw.ModControl) != 0
	platform.imguiIO.KeyShift = (mods & glfw.ModShift) != 0
	platform.imguiIO.KeyAlt = (mods & glfw.ModAlt) != 0
	platform.imguiIO.KeySuper = (mods & glfw.ModSuper) != 0
}

func (platform *GLFW) charChange(window *glfw.Window, char rune) {
	platform.imguiIO.AddInputCharacters(string(char))
}

// ClipboardText returns the current clipboard text, if available.
func (platform *GLFW) ClipboardText() (string, error) {
	return platform.window.GetClipboardString()
}

// SetClipboardText sets the text as the current clipboard text.
func (platform *GLFW) SetClipboardText(text string) {
	platform.window.SetClipboardString(text)
}
