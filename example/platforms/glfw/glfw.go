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
	imguiIO *imgui.ImGuiIO

	window *glfw.Window

	time             float64
	mouseJustPressed [3]bool
}

// NewGLFW attempts to initialize a GLFW context.
func NewGLFW(io *imgui.ImGuiIO, clientAPI GLFWClientAPI) (*GLFW, error) {
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
	platform.imguiIO.DisplaySize = *imgui.NewImVec2(displaySize[0], displaySize[1])

	// Setup time step
	currentTime := glfw.GetTime()
	if platform.time > 0 {
		platform.imguiIO.DeltaTime = float32(currentTime - platform.time)
	}
	platform.time = currentTime

	// Setup inputs
	if platform.window.GetAttrib(glfw.Focused) != 0 {
		x, y := platform.window.GetCursorPos()
		platform.imguiIO.MousePos = *imgui.NewImVec2(float32(x), float32(y))
	} else {
		platform.imguiIO.MousePos = *imgui.NewImVec2(-math.MaxFloat32, -math.MaxFloat32)
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
	platform.imguiIO.KeyMap[imgui.ImGuiKey_Tab] = int32(glfw.KeyTab)
	platform.imguiIO.KeyMap[imgui.ImGuiKey_LeftArrow] = int32(glfw.KeyLeft)
	platform.imguiIO.KeyMap[imgui.ImGuiKey_RightArrow] = int32(glfw.KeyRight)
	platform.imguiIO.KeyMap[imgui.ImGuiKey_UpArrow] = int32(glfw.KeyUp)
	platform.imguiIO.KeyMap[imgui.ImGuiKey_DownArrow] = int32(glfw.KeyDown)
	platform.imguiIO.KeyMap[imgui.ImGuiKey_PageUp] = int32(glfw.KeyPageUp)
	platform.imguiIO.KeyMap[imgui.ImGuiKey_PageDown] = int32(glfw.KeyPageDown)
	platform.imguiIO.KeyMap[imgui.ImGuiKey_Home] = int32(glfw.KeyHome)
	platform.imguiIO.KeyMap[imgui.ImGuiKey_End] = int32(glfw.KeyEnd)
	platform.imguiIO.KeyMap[imgui.ImGuiKey_Insert] = int32(glfw.KeyInsert)
	platform.imguiIO.KeyMap[imgui.ImGuiKey_Delete] = int32(glfw.KeyDelete)
	platform.imguiIO.KeyMap[imgui.ImGuiKey_Backspace] = int32(glfw.KeyBackspace)
	platform.imguiIO.KeyMap[imgui.ImGuiKey_Space] = int32(glfw.KeySpace)
	platform.imguiIO.KeyMap[imgui.ImGuiKey_Enter] = int32(glfw.KeyEnter)
	platform.imguiIO.KeyMap[imgui.ImGuiKey_Escape] = int32(glfw.KeyEscape)
	platform.imguiIO.KeyMap[imgui.ImGuiKey_A] = int32(glfw.KeyA)
	platform.imguiIO.KeyMap[imgui.ImGuiKey_C] = int32(glfw.KeyC)
	platform.imguiIO.KeyMap[imgui.ImGuiKey_V] = int32(glfw.KeyV)
	platform.imguiIO.KeyMap[imgui.ImGuiKey_X] = int32(glfw.KeyX)
	platform.imguiIO.KeyMap[imgui.ImGuiKey_Y] = int32(glfw.KeyY)
	platform.imguiIO.KeyMap[imgui.ImGuiKey_Z] = int32(glfw.KeyZ)
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
	platform.imguiIO.AddInputCharacter(char)
}

// ClipboardText returns the current clipboard text, if available.
func (platform *GLFW) ClipboardText() (string, error) {
	return platform.window.GetClipboardString()
}

// SetClipboardText sets the text as the current clipboard text.
func (platform *GLFW) SetClipboardText(text string) {
	platform.window.SetClipboardString(text)
}
