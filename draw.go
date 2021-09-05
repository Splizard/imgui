package imgui

type TextureID uintptr
type DrawIdx int16

type drawCmdHeader struct {
	clipRect  Vec4
	textureID TextureID
	vtxOffset uint
}

type DrawCallback func(list *DrawList, cmd *DrawCmd)

//DrawCmd is a single draw command within a parent ImDrawList (generally maps to 1 GPU draw call, unless it is a callback)
type DrawCmd struct {
	ClipRect         Vec4         // 4*4  // Clipping rectangle (x1, y1, x2, y2). Subtract ImDrawData->DisplayPos to get clipping rectangle in "viewport" coordinates
	TextureId        TextureID    // 4-8  // User-provided texture ID. Set by user in ImfontAtlas::SetTexID() for fonts or passed to Image*() functions. Ignore if never using images or multiple fonts atlas.
	VtxOffset        uint         // 4    // Start offset in vertex buffer. ImGuiBackendFlags_RendererHasVtxOffset: always 0, otherwise may be >0 to support meshes larger than 64K vertices with 16-bit indices.
	IdxOffset        uint         // 4    // Start offset in index buffer. Always equal to sum of ElemCount drawn so far.
	ElemCount        uint         // 4    // Number of indices (multiple of 3) to be rendered as triangles. Vertices are stored in the callee ImDrawList's vtx_buffer[] array, indices in idx_buffer[].
	UserCallback     DrawCallback // 4-8  // If != NULL, call the function instead of rendering the vertices. clip_rect and texture_id will be set normally.
	UserCallbackData interface{}  // 4-8  // The draw callback code can access this.
}

type DrawList struct {
	// This is what you have to render
	CmdBuffer []DrawCmd     // Draw commands. Typically 1 command = 1 GPU draw call, unless the command is a callback.
	IdxBuffer []DrawIdx     // Index buffer. Each command consume ImDrawCmd::ElemCount of those
	VtxBuffer []DrawVert    // Vertex buffer.
	Flags     DrawListFlags // Flags, you may poke into these to adjust anti-aliasing settings per-primitive.

	// [Internal, used while building lists]
	vtxCurrentIdx  uint                // [Internal] generally == VtxBuffer.Size unless we are past 64K vertices, in which case this gets reset to 0.
	data           *DrawListSharedData // Pointer to shared draw data (you can use ImGui::GetDrawListSharedData() to get the one from current ImGui context)
	OwnerName      string              // Pointer to owner window's name for debugging
	vtxWritePtr    *DrawVert           // [Internal] point within VtxBuffer.Data after each add command (to avoid using the ImVector<> operators too much)
	idxWritePtr    *DrawIdx            // [Internal] point within IdxBuffer.Data after each add command (to avoid using the ImVector<> operators too much)
	clipRectStack  []Vec4              // [Internal]
	textureIdStack []TextureID         // [Internal]
	path           []Vec2              // [Internal] current path building
	cmdHeader      drawCmdHeader       // [Internal] template of active commands. Fields should match those of CmdBuffer.back().
	splitter       DrawListSplitter    // [Internal] for channels api (note: prefer using your own persistent instance of ImDrawListSplitter!)
	fringeScale    float32             // [Internal] anti-alias fringe is scaled by this value, this helps to keep things sharp while zooming at vertex buffer content
}

//DrawVert is a single vertex (pos + uv + col = 20 bytes by default).
type DrawVert struct {
	Pos Vec2
	UV  Vec2
	Col uint32
}

type DrawData struct {
	Valid            bool        // Only valid after Render() is called and before the next NewFrame() is called.
	CmdListsCount    int         // Number of ImDrawList* to render
	TotalIdxCount    int         // For convenience, sum of all ImDrawList's IdxBuffer.Size
	TotalVtxCount    int         // For convenience, sum of all ImDrawList's VtxBuffer.Size
	CmdLists         []*DrawList // Array of ImDrawList* to render. The ImDrawList are owned by ImGuiContext and only pointed to from here.
	DisplayPos       Vec2        // Top-left position of the viewport to render (== top-left of the orthogonal projection matrix to use) (== GetMainViewport()->Pos for the main viewport, == (0.0) in most single-viewport applications)
	DisplaySize      Vec2        // Size of the viewport to render (== GetMainViewport()->Size for the main viewport, == io.DisplaySize in most single-viewport applications)
	FramebufferScale Vec2        // Amount of pixels for each unit of DisplaySize. Based on io.DisplayFramebufferScale. Generally (1,1) on normal display, (2,2) on OSX with Retina display.
}

func (d *DrawData) DeIndexAllBuffers() { panic("not implemented") }

func (d *DrawData) ScaleClipRects(scale Vec2) { panic("not implemented") }
