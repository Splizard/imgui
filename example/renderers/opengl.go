package renderers

import (
	_ "embed" // using embed for the shader sources
	"fmt"
	"os"
	"strings"
	"unsafe"

	"github.com/go-gl/gl/v3.2-core/gl"
	"github.com/jojomi/go-spew/spew"
	"github.com/zeozeozeo/imgui"
)

//go:embed shader/main.vert
var unversionedVertexShader string

//go:embed shader/main.frag
var unversionedFragmentShader string

// OpenGL3 implements a renderer based on github.com/go-gl/gl (v3.2-core).
type OpenGL3 struct {
	imguiIO *imgui.ImGuiIO

	glslVersion            string
	fontTexture            uint32
	shaderHandle           uint32
	vertHandle             uint32
	fragHandle             uint32
	attribLocationTex      int32
	attribLocationProjMtx  int32
	attribLocationPosition int32
	attribLocationUV       int32
	attribLocationColor    int32
	vboHandle              uint32
	elementsHandle         uint32
}

// NewOpenGL3 attempts to initialize a renderer.
// An OpenGL context has to be established before calling this function.
func NewOpenGL3(io *imgui.ImGuiIO) (*OpenGL3, error) {
	err := gl.Init()
	if err != nil {
		return nil, fmt.Errorf("failed to initialize OpenGL: %w", err)
	}

	renderer := &OpenGL3{
		imguiIO:     io,
		glslVersion: "#version 150",
	}
	renderer.createDeviceObjects()
	return renderer, nil
}

// Dispose cleans up the resources.
func (renderer *OpenGL3) Dispose() {
	renderer.invalidateDeviceObjects()
}

// PreRender clears the framebuffer.
func (renderer *OpenGL3) PreRender(clearColor [3]float32) {
	gl.ClearColor(clearColor[0], clearColor[1], clearColor[2], 1.0)
	gl.Clear(gl.COLOR_BUFFER_BIT)
}

var frames int

// Render translates the ImGui draw data to OpenGL3 commands.
func (renderer *OpenGL3) Render(displaySize [2]float32, framebufferSize [2]float32, drawData *imgui.ImDrawData) {
	// Avoid rendering when minimized, scale coordinates for retina displays (screen coordinates != framebuffer coordinates)
	displayWidth, displayHeight := displaySize[0], displaySize[1]
	fbWidth, fbHeight := framebufferSize[0], framebufferSize[1]
	if (fbWidth <= 0) || (fbHeight <= 0) {
		return
	}
	drawData.ScaleClipRects(imgui.NewImVec2(
		fbWidth/displayWidth,
		fbHeight/displayHeight,
	))

	// Backup GL state
	var lastActiveTexture int32
	gl.GetIntegerv(gl.ACTIVE_TEXTURE, &lastActiveTexture)
	gl.ActiveTexture(gl.TEXTURE0)
	var lastProgram int32
	gl.GetIntegerv(gl.CURRENT_PROGRAM, &lastProgram)
	var lastTexture int32
	gl.GetIntegerv(gl.TEXTURE_BINDING_2D, &lastTexture)
	var lastSampler int32
	gl.GetIntegerv(gl.SAMPLER_BINDING, &lastSampler)
	var lastArrayBuffer int32
	gl.GetIntegerv(gl.ARRAY_BUFFER_BINDING, &lastArrayBuffer)
	var lastElementArrayBuffer int32
	gl.GetIntegerv(gl.ELEMENT_ARRAY_BUFFER_BINDING, &lastElementArrayBuffer)
	var lastVertexArray int32
	gl.GetIntegerv(gl.VERTEX_ARRAY_BINDING, &lastVertexArray)
	var lastPolygonMode [2]int32
	gl.GetIntegerv(gl.POLYGON_MODE, &lastPolygonMode[0])
	var lastViewport [4]int32
	gl.GetIntegerv(gl.VIEWPORT, &lastViewport[0])
	var lastScissorBox [4]int32
	gl.GetIntegerv(gl.SCISSOR_BOX, &lastScissorBox[0])
	var lastBlendSrcRgb int32
	gl.GetIntegerv(gl.BLEND_SRC_RGB, &lastBlendSrcRgb)
	var lastBlendDstRgb int32
	gl.GetIntegerv(gl.BLEND_DST_RGB, &lastBlendDstRgb)
	var lastBlendSrcAlpha int32
	gl.GetIntegerv(gl.BLEND_SRC_ALPHA, &lastBlendSrcAlpha)
	var lastBlendDstAlpha int32
	gl.GetIntegerv(gl.BLEND_DST_ALPHA, &lastBlendDstAlpha)
	var lastBlendEquationRgb int32
	gl.GetIntegerv(gl.BLEND_EQUATION_RGB, &lastBlendEquationRgb)
	var lastBlendEquationAlpha int32
	gl.GetIntegerv(gl.BLEND_EQUATION_ALPHA, &lastBlendEquationAlpha)
	lastEnableBlend := gl.IsEnabled(gl.BLEND)
	lastEnableCullFace := gl.IsEnabled(gl.CULL_FACE)
	lastEnableDepthTest := gl.IsEnabled(gl.DEPTH_TEST)
	lastEnableScissorTest := gl.IsEnabled(gl.SCISSOR_TEST)

	// Setup render state: alpha-blending enabled, no face culling, no depth testing, scissor enabled, polygon fill
	gl.Enable(gl.BLEND)
	gl.BlendEquation(gl.FUNC_ADD)
	gl.BlendFunc(gl.SRC_ALPHA, gl.ONE_MINUS_SRC_ALPHA)
	gl.Disable(gl.CULL_FACE)
	gl.Disable(gl.DEPTH_TEST)
	gl.Enable(gl.SCISSOR_TEST)
	gl.PolygonMode(gl.FRONT_AND_BACK, gl.FILL)

	// Setup viewport, orthographic projection matrix
	// Our visible imgui space lies from draw_data->DisplayPos (top left) to draw_data->DisplayPos+data_data->DisplaySize (bottom right).
	// DisplayMin is typically (0,0) for single viewport apps.
	gl.Viewport(0, 0, int32(fbWidth), int32(fbHeight))
	orthoProjection := [4][4]float32{
		{2.0 / displayWidth, 0.0, 0.0, 0.0},
		{0.0, 2.0 / -displayHeight, 0.0, 0.0},
		{0.0, 0.0, -1.0, 0.0},
		{-1.0, 1.0, 0.0, 1.0},
	}
	gl.UseProgram(renderer.shaderHandle)
	gl.Uniform1i(renderer.attribLocationTex, 0)
	gl.UniformMatrix4fv(renderer.attribLocationProjMtx, 1, false, &orthoProjection[0][0])
	gl.BindSampler(0, 0) // Rely on combined texture/sampler state.

	// Recreate the VAO every time
	// (This is to easily allow multiple GL contexts. VAO are not shared among GL contexts, and
	// we don't track creation/deletion of windows so we don't have an obvious key to use to cache them.)
	var vaoHandle uint32
	gl.GenVertexArrays(1, &vaoHandle)
	gl.BindVertexArray(vaoHandle)
	gl.BindBuffer(gl.ARRAY_BUFFER, renderer.vboHandle)
	gl.EnableVertexAttribArray(uint32(renderer.attribLocationPosition))
	gl.EnableVertexAttribArray(uint32(renderer.attribLocationUV))
	gl.EnableVertexAttribArray(uint32(renderer.attribLocationColor))

	vertexSize, vertexOffsetPos, vertexOffsetUv, vertexOffsetCol := imgui.ImDrawVertSizeAndOffset()

	gl.VertexAttribPointerWithOffset(uint32(renderer.attribLocationPosition), 2, gl.FLOAT, false, int32(vertexSize), uintptr(vertexOffsetPos))
	gl.VertexAttribPointerWithOffset(uint32(renderer.attribLocationUV), 2, gl.FLOAT, false, int32(vertexSize), uintptr(vertexOffsetUv))
	gl.VertexAttribPointerWithOffset(uint32(renderer.attribLocationColor), 4, gl.UNSIGNED_BYTE, true, int32(vertexSize), uintptr(vertexOffsetCol))

	indexSize := unsafe.Sizeof(imgui.ImDrawIdx(0))
	drawType := gl.UNSIGNED_SHORT
	const bytesPerUint32 = 4
	if indexSize == bytesPerUint32 {
		drawType = gl.UNSIGNED_INT
	}

	// Draw

	for _, list := range drawData.CmdLists {
		var indexBufferOffset uintptr

		if os.Getenv("DUMP") != "" {
			if frames == 2 {
				fmt.Println(list.IdxBuffer)
				spew.Dump(list.VtxBuffer)
				os.Exit(0)
			} else {
				frames++
			}
		}

		gl.BindBuffer(gl.ARRAY_BUFFER, renderer.vboHandle)
		gl.BufferData(gl.ARRAY_BUFFER, len(list.VtxBuffer)*int(unsafe.Sizeof(list.VtxBuffer[0])), unsafe.Pointer(&list.VtxBuffer[0]), gl.STREAM_DRAW)

		gl.BindBuffer(gl.ELEMENT_ARRAY_BUFFER, renderer.elementsHandle)
		gl.BufferData(gl.ELEMENT_ARRAY_BUFFER, len(list.IdxBuffer)*int(unsafe.Sizeof(list.IdxBuffer[0])), unsafe.Pointer(&list.IdxBuffer[0]), gl.STREAM_DRAW)

		//fmt.Println(list.VtxBuffer)
		//fmt.Println(list.IdxBuffer)
		//os.Exit(0)

		//fmt.Println(len(list.CmdBuffer))

		for i, cmd := range list.CmdBuffer {
			if cmd.UserCallback != nil {
				cmd.UserCallback(list, &list.CmdBuffer[i])
			} else {
				gl.BindTexture(gl.TEXTURE_2D, uint32(cmd.TextureId))
				clipRect := cmd.ClipRect

				gl.Scissor(int32(clipRect.X()), int32(fbHeight)-int32(clipRect.W()), int32(clipRect.Z()-clipRect.X()), int32(clipRect.W()-clipRect.Y()))
				gl.DrawElementsWithOffset(gl.TRIANGLES, int32(cmd.ElemCount), uint32(drawType), indexBufferOffset)
			}
			indexBufferOffset += uintptr(cmd.ElemCount) * indexSize
		}
	}
	gl.DeleteVertexArrays(1, &vaoHandle)

	// Restore modified GL state
	gl.UseProgram(uint32(lastProgram))
	gl.BindTexture(gl.TEXTURE_2D, uint32(lastTexture))
	gl.BindSampler(0, uint32(lastSampler))
	gl.ActiveTexture(uint32(lastActiveTexture))
	gl.BindVertexArray(uint32(lastVertexArray))
	gl.BindBuffer(gl.ARRAY_BUFFER, uint32(lastArrayBuffer))
	gl.BindBuffer(gl.ELEMENT_ARRAY_BUFFER, uint32(lastElementArrayBuffer))
	gl.BlendEquationSeparate(uint32(lastBlendEquationRgb), uint32(lastBlendEquationAlpha))
	gl.BlendFuncSeparate(uint32(lastBlendSrcRgb), uint32(lastBlendDstRgb), uint32(lastBlendSrcAlpha), uint32(lastBlendDstAlpha))
	if lastEnableBlend {
		gl.Enable(gl.BLEND)
	} else {
		gl.Disable(gl.BLEND)
	}
	if lastEnableCullFace {
		gl.Enable(gl.CULL_FACE)
	} else {
		gl.Disable(gl.CULL_FACE)
	}
	if lastEnableDepthTest {
		gl.Enable(gl.DEPTH_TEST)
	} else {
		gl.Disable(gl.DEPTH_TEST)
	}
	if lastEnableScissorTest {
		gl.Enable(gl.SCISSOR_TEST)
	} else {
		gl.Disable(gl.SCISSOR_TEST)
	}
	gl.PolygonMode(gl.FRONT_AND_BACK, uint32(lastPolygonMode[0]))
	gl.Viewport(lastViewport[0], lastViewport[1], lastViewport[2], lastViewport[3])
	gl.Scissor(lastScissorBox[0], lastScissorBox[1], lastScissorBox[2], lastScissorBox[3])
}

func (renderer *OpenGL3) createDeviceObjects() {
	// Backup GL state
	var lastTexture int32
	var lastArrayBuffer int32
	var lastVertexArray int32
	gl.GetIntegerv(gl.TEXTURE_BINDING_2D, &lastTexture)
	gl.GetIntegerv(gl.ARRAY_BUFFER_BINDING, &lastArrayBuffer)
	gl.GetIntegerv(gl.VERTEX_ARRAY_BINDING, &lastVertexArray)

	vertexShader := renderer.glslVersion + "\n" + unversionedVertexShader
	fragmentShader := renderer.glslVersion + "\n" + unversionedFragmentShader

	renderer.shaderHandle = gl.CreateProgram()
	renderer.vertHandle = gl.CreateShader(gl.VERTEX_SHADER)
	renderer.fragHandle = gl.CreateShader(gl.FRAGMENT_SHADER)

	compileSource := func(kind uint32, source string) (uint32, error) {
		shader := gl.CreateShader(kind)

		csources, free := gl.Strs(string(source) + "\000")
		gl.ShaderSource(shader, 1, csources, nil)
		free()
		gl.CompileShader(shader)

		var status int32
		gl.GetShaderiv(shader, gl.COMPILE_STATUS, &status)
		if status == gl.FALSE {
			var logLength int32
			gl.GetShaderiv(shader, gl.INFO_LOG_LENGTH, &logLength)

			log := strings.Repeat("\x00", int(logLength+1))
			gl.GetShaderInfoLog(shader, logLength, nil, gl.Str(log))

			gl.DeleteShader(shader)

			return 0, fmt.Errorf("failed to compile shader: %v\n %v", log, source)
		}

		return shader, nil
	}

	var err error

	renderer.vertHandle, err = compileSource(gl.VERTEX_SHADER, vertexShader)
	if err != nil {
		panic(err)
	}
	renderer.fragHandle, err = compileSource(gl.FRAGMENT_SHADER, fragmentShader)
	if err != nil {
		panic(err)
	}

	fmt.Println("no errors")

	gl.CompileShader(renderer.vertHandle)
	gl.CompileShader(renderer.fragHandle)
	gl.AttachShader(renderer.shaderHandle, renderer.vertHandle)
	gl.AttachShader(renderer.shaderHandle, renderer.fragHandle)
	gl.LinkProgram(renderer.shaderHandle)

	renderer.attribLocationTex = gl.GetUniformLocation(renderer.shaderHandle, gl.Str("Texture"+"\x00"))
	renderer.attribLocationProjMtx = gl.GetUniformLocation(renderer.shaderHandle, gl.Str("ProjMtx"+"\x00"))
	renderer.attribLocationPosition = gl.GetAttribLocation(renderer.shaderHandle, gl.Str("Position"+"\x00"))
	renderer.attribLocationUV = gl.GetAttribLocation(renderer.shaderHandle, gl.Str("UV"+"\x00"))
	renderer.attribLocationColor = gl.GetAttribLocation(renderer.shaderHandle, gl.Str("Color"+"\x00"))

	gl.GenBuffers(1, &renderer.vboHandle)
	gl.GenBuffers(1, &renderer.elementsHandle)

	renderer.createFontsTexture()

	// Restore modified GL state
	gl.BindTexture(gl.TEXTURE_2D, uint32(lastTexture))
	gl.BindBuffer(gl.ARRAY_BUFFER, uint32(lastArrayBuffer))
	gl.BindVertexArray(uint32(lastVertexArray))
}

func (renderer *OpenGL3) createFontsTexture() {
	// Build texture atlas
	io := imgui.GetIO()
	var pixels []byte
	var width, height int32

	io.Fonts.GetTexDataAsAlpha8(&pixels, &width, &height, nil)

	// Upload texture to graphics system
	var lastTexture int32
	gl.GetIntegerv(gl.TEXTURE_BINDING_2D, &lastTexture)
	gl.GenTextures(1, &renderer.fontTexture)
	gl.BindTexture(gl.TEXTURE_2D, renderer.fontTexture)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MIN_FILTER, gl.LINEAR)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MAG_FILTER, gl.LINEAR)
	gl.PixelStorei(gl.UNPACK_ROW_LENGTH, 0)
	gl.TexImage2D(gl.TEXTURE_2D, 0, gl.RED, int32(width), int32(height),
		0, gl.RED, gl.UNSIGNED_BYTE, gl.Ptr(&pixels[0]))

	// Store our identifier
	io.Fonts.SetTexID(imgui.ImTextureID(renderer.fontTexture))

	// Restore state
	gl.BindTexture(gl.TEXTURE_2D, uint32(lastTexture))
}

func (renderer *OpenGL3) invalidateDeviceObjects() {
	if renderer.vboHandle != 0 {
		gl.DeleteBuffers(1, &renderer.vboHandle)
	}
	renderer.vboHandle = 0
	if renderer.elementsHandle != 0 {
		gl.DeleteBuffers(1, &renderer.elementsHandle)
	}
	renderer.elementsHandle = 0

	if (renderer.shaderHandle != 0) && (renderer.vertHandle != 0) {
		gl.DetachShader(renderer.shaderHandle, renderer.vertHandle)
	}
	if renderer.vertHandle != 0 {
		gl.DeleteShader(renderer.vertHandle)
	}
	renderer.vertHandle = 0

	if (renderer.shaderHandle != 0) && (renderer.fragHandle != 0) {
		gl.DetachShader(renderer.shaderHandle, renderer.fragHandle)
	}
	if renderer.fragHandle != 0 {
		gl.DeleteShader(renderer.fragHandle)
	}
	renderer.fragHandle = 0

	if renderer.shaderHandle != 0 {
		gl.DeleteProgram(renderer.shaderHandle)
	}
	renderer.shaderHandle = 0

	if renderer.fontTexture != 0 {
		gl.DeleteTextures(1, &renderer.fontTexture)
		imgui.GetIO().Fonts.SetTexID(0)
		renderer.fontTexture = 0
	}
}
