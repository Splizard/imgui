package imgui

import (
	"unsafe"
)

// Pass this to your backend rendering function! Valid after Render() and until the next call to NewFrame()
func GetDrawData() *ImDrawData {
	var g = GImGui
	var viewport = g.Viewports[0]

	if viewport.DrawDataP.Valid {
		return &viewport.DrawDataP
	} else {
		return nil
	}
} // valid after Render() and until the next call to NewFrame(). this is what you have to render.

// All draw data to render a Dear ImGui frame
// (NB: the style and the naming convention here is a little inconsistent, we currently preserve them for backward compatibility purpose,
// as this is one of the oldest structure exposed by the library! Basically, ImDrawList == CmdList)
type ImDrawData struct {
	Valid            bool          // Only valid after Render() is called and before the next NewFrame() is called.
	CmdListsCount    int           // Number of ImDrawList* to render
	TotalIdxCount    int           // For convenience, sum of all ImDrawList's IdxBuffer.Size
	TotalVtxCount    int           // For convenience, sum of all ImDrawList's VtxBuffer.Size
	CmdLists         []*ImDrawList // Array of ImDrawList* to render. The ImDrawList are owned by ImGuiContext and only pointed to from here.
	DisplayPos       ImVec2        // Top-left position of the viewport to render (== top-left of the orthogonal projection matrix to use) (== GetMainViewport()->Pos for the main viewport, == (0.0) in most single-viewport applications)
	DisplaySize      ImVec2        // Size of the viewport to render (== GetMainViewport()->Size for the main viewport, == io.DisplaySize in most single-viewport applications)
	FramebufferScale ImVec2        // Amount of pixels for each unit of DisplaySize. Based on io.DisplayFramebufferScale. Generally (1,1) on normal display, (2,2) on OSX with Retina display.
}

func (this *ImDrawData) Clear() {
	*this = ImDrawData{}
}

// Functions
func (ImDrawData) DeIndexAllBuffers() { panic("not implemented") } // Helper to convert all buffers from indexed to non-indexed, in case you cannot render indexed. Note: this is slow and most likely a waste of resources. Always prefer indexed rendering!

// Helper to scale the ClipRect field of each ImDrawCmd.
// Use if your final output buffer is at a different scale than draw_data.DisplaySize,
// or if there is a difference between your window resolution and framebuffer resolution.
func (this *ImDrawData) ScaleClipRects(fb_scale *ImVec2) {
	for i := int(0); i < this.CmdListsCount; i++ {
		var cmd_list *ImDrawList = this.CmdLists[i]
		for cmd_i := range cmd_list.CmdBuffer {
			var cmd *ImDrawCmd = &cmd_list.CmdBuffer[cmd_i]
			cmd.ClipRect = ImVec4{cmd.ClipRect.x * fb_scale.x, cmd.ClipRect.y * fb_scale.y, cmd.ClipRect.z * fb_scale.x, cmd.ClipRect.w * fb_scale.y}
		}
	}
} // Helper to scale the ClipRect field of each ImDrawCmd. Use if your final output buffer is at a different scale than Dear ImGui expects, or if there is a difference between your window resolution and framebuffer resolution.

// [Internal helpers]
func (this *ImDrawList) _ResetForNewFrame() {
	// Verify that the ImDrawCmd fields we want to memcmp() are contiguous in memory.
	// (those should be IM_STATIC_ASSERT() in theory but with our pre C++11 setup the whole check doesn't compile with GCC)
	IM_ASSERT(unsafe.Offsetof(ImDrawCmd{}.ClipRect) == 0)
	IM_ASSERT(unsafe.Offsetof(ImDrawCmd{}.TextureId) == unsafe.Sizeof(ImVec4{}))
	IM_ASSERT(unsafe.Offsetof(ImDrawCmd{}.VtxOffset) == unsafe.Sizeof(ImVec4{})+unsafe.Sizeof(ImTextureID(0)))

	this.CmdBuffer = this.CmdBuffer[:0]
	this.IdxBuffer = this.IdxBuffer[:0]
	this.VtxBuffer = this.VtxBuffer[:0]
	this.Flags = this._Data.InitialFlags
	this._CmdHeader = ImDrawCmdHeader{}

	this._VtxCurrentIdx = 0
	this._VtxWritePtr = 0
	this._IdxWritePtr = 0
	this._ClipRectStack = this._ClipRectStack[:0]
	this._TextureIdStack = this._TextureIdStack[:0]
	this._Path = this._Path[:0]
	this._Splitter.Clear()
	this.CmdBuffer = append(this.CmdBuffer, ImDrawCmd{})
	this._FringeScale = 1.0
}

func (this *ImDrawList) PushTextureID(texture_id ImTextureID) {
	this._TextureIdStack = append(this._TextureIdStack, texture_id)
	this._CmdHeader.TextureId = texture_id
	this._OnChangedTextureID()
}

func (this *ImDrawList) _OnChangedTextureID() {
	// If current command is used with different settings we need to add a new command
	var curr_cmd *ImDrawCmd = &this.CmdBuffer[len(this.CmdBuffer)-1]
	if curr_cmd.ElemCount != 0 && curr_cmd.TextureId != this._CmdHeader.TextureId {
		this.AddDrawCmd()
		return
	}
	IM_ASSERT(curr_cmd.UserCallback == nil)

	// Try to merge with previous command if it matches, else use current command
	var prev_cmd *ImDrawCmd = &this.CmdBuffer[len(this.CmdBuffer)-1]

	prevHeader := ImDrawCmdHeader{
		ClipRect:  prev_cmd.ClipRect,
		TextureId: prev_cmd.TextureId,
		VtxOffset: prev_cmd.ElemCount,
	}

	if curr_cmd.ElemCount == 0 && len(this.CmdBuffer) > 1 && this._CmdHeader == prevHeader && prev_cmd.UserCallback == nil {
		this.CmdBuffer = this.CmdBuffer[len(this.CmdBuffer)-1:]
		return
	}

	curr_cmd.TextureId = this._CmdHeader.TextureId

}

// Our scheme may appears a bit unusual, basically we want the most-common calls AddLine AddRect etc. to not have to perform any check so we always have a command ready in the stack.
// The cost of figuring out if a new command has to be added or if we can merge is paid in those Update** functions only.
func (this *ImDrawList) _OnChangedClipRect() {
	// If current command is used with different settings we need to add a new command
	var curr_cmd *ImDrawCmd = &this.CmdBuffer[len(this.CmdBuffer)-1]
	if curr_cmd.ElemCount != 0 && curr_cmd.ClipRect != this._CmdHeader.ClipRect {
		this.AddDrawCmd()
		return
	}
	IM_ASSERT(curr_cmd.UserCallback == nil)

	// Try to merge with previous command if it matches, else use current command
	var prev_cmd *ImDrawCmd = &this.CmdBuffer[len(this.CmdBuffer)-1]

	prevHeader := ImDrawCmdHeader{
		ClipRect:  prev_cmd.ClipRect,
		TextureId: prev_cmd.TextureId,
		VtxOffset: prev_cmd.ElemCount,
	}

	if curr_cmd.ElemCount == 0 && len(this.CmdBuffer) > 1 && this._CmdHeader == prevHeader && prev_cmd.UserCallback == nil {
		this.CmdBuffer = this.CmdBuffer[len(this.CmdBuffer)-1:]
		return
	}

	curr_cmd.ClipRect = this._CmdHeader.ClipRect
}

func (this *ImDrawList) AddRectFilled(p_min, p_max ImVec2, col ImU32, rounding float, flags ImDrawFlags) {

	if (col & IM_COL32_A_MASK) == 0 {
		return
	}

	if rounding <= 0.0 || (flags&ImDrawFlags_RoundCornersMask_) == ImDrawFlags_RoundCornersNone {
		this.PrimReserve(6, 4)
		this.PrimRect(&p_min, &p_max, col)
	} else {
		this.PathRect(&p_min, &p_max, rounding, flags)

		this.PathFillConvex(col)
	}

} // a: upper-left, b: lower-right (== upper-left + size)

func (this *ImDrawList) PathArcToFast(center ImVec2, radius float, a_min_sample, a_max_sample int) {
	if radius <= 0.0 {
		this._Path = append(this._Path, center)
		return
	}
	this.PathArcToFastEx(center, radius, a_min_sample*IM_DRAWLIST_ARCFAST_SAMPLE_MAX/12, a_max_sample*IM_DRAWLIST_ARCFAST_SAMPLE_MAX/12, 0)
}

// We intentionally avoid using ImVec2 and its math operators here to reduce cost to a minimum for debug/non-inlined builds.
func (this *ImDrawList) AddConvexPolyFilled(points []ImVec2, points_count int, col ImU32) {
	if points_count < 3 {
		return
	}

	var uv ImVec2 = this._Data.TexUvWhitePixel

	if this.Flags&ImDrawListFlags_AntiAliasedFill != 0 {

		// Anti-aliased Fill
		var AA_SIZE float = this._FringeScale
		var col_trans ImU32 = col &^ IM_COL32_A_MASK
		var idx_count int = (points_count-2)*3 + points_count*6
		var vtx_count int = (points_count * 2)
		this.PrimReserve(idx_count, vtx_count)

		// Add indexes for fill
		var vtx_inner_idx uint = this._VtxCurrentIdx
		var vtx_outer_idx uint = this._VtxCurrentIdx + 1
		for i := int(2); i < points_count; i++ {
			this.IdxBuffer[this._IdxWritePtr] = (ImDrawIdx)(vtx_inner_idx)
			this.IdxBuffer[this._IdxWritePtr+1] = (ImDrawIdx)(vtx_inner_idx + ((uint(i) - 1) << 1))
			this.IdxBuffer[this._IdxWritePtr+2] = (ImDrawIdx)(vtx_inner_idx + (uint(i) << 1))
			this._IdxWritePtr += 3
		}

		// Compute normals
		var temp_normals []ImVec2 = make([]ImVec2, points_count) //-V630
		for i0, i1 := points_count-1, int(0); i1 < points_count; i0, i1 = i1, i1+1 {

			var p0 *ImVec2 = &points[i0]
			var p1 *ImVec2 = &points[i1]
			var dx float = p1.x - p0.x
			var dy float = p1.y - p0.y
			IM_NORMALIZE2F_OVER_ZERO(&dx, &dy)
			temp_normals[i0].x = dy
			temp_normals[i0].y = -dx
		}

		for i0, i1 := points_count-1, int(0); i1 < points_count; i0, i1 = i1, i1+1 {
			// Average normals
			var n0 *ImVec2 = &temp_normals[i0]
			var n1 *ImVec2 = &temp_normals[i1]
			var dm_x float = (n0.x + n1.x) * 0.5
			var dm_y float = (n0.y + n1.y) * 0.5
			IM_FIXNORMAL2F(&dm_x, &dm_y)
			dm_x *= AA_SIZE * 0.5
			dm_y *= AA_SIZE * 0.5

			// Add vertices
			this.VtxBuffer[this._VtxWritePtr+0].pos.x = (points[i1].x - dm_x)
			this.VtxBuffer[this._VtxWritePtr+0].pos.y = (points[i1].y - dm_y)
			this.VtxBuffer[this._VtxWritePtr+0].uv = uv
			this.VtxBuffer[this._VtxWritePtr+0].col = col // Inner
			this.VtxBuffer[this._VtxWritePtr+1].pos.x = (points[i1].x + dm_x)
			this.VtxBuffer[this._VtxWritePtr+1].pos.y = (points[i1].y + dm_y)
			this.VtxBuffer[this._VtxWritePtr+1].uv = uv
			this.VtxBuffer[this._VtxWritePtr+1].col = col_trans // Outer

			this._VtxWritePtr += 2

			// Add indexes for fringes
			this.IdxBuffer[this._IdxWritePtr+0] = (ImDrawIdx)(vtx_inner_idx + uint(i1<<1))
			this.IdxBuffer[this._IdxWritePtr+1] = (ImDrawIdx)(vtx_inner_idx + uint(i0<<1))
			this.IdxBuffer[this._IdxWritePtr+2] = (ImDrawIdx)(vtx_outer_idx + uint(i0<<1))
			this.IdxBuffer[this._IdxWritePtr+3] = (ImDrawIdx)(vtx_outer_idx + uint(i0<<1))
			this.IdxBuffer[this._IdxWritePtr+4] = (ImDrawIdx)(vtx_outer_idx + uint(i1<<1))
			this.IdxBuffer[this._IdxWritePtr+5] = (ImDrawIdx)(vtx_inner_idx + uint(i1<<1))
			this._IdxWritePtr += 6
		}
		//printf("vtx_count %d\n", vtx_count)
		this._VtxCurrentIdx += uint(vtx_count)
	} else {

		// Non Anti-aliased Fill
		var idx_count int = (points_count - 2) * 3
		var vtx_count int = points_count
		this.PrimReserve(idx_count, vtx_count)
		for i := int(0); i < vtx_count; i++ {
			this.VtxBuffer[this._VtxWritePtr+0].pos = points[i]
			this.VtxBuffer[this._VtxWritePtr+0].uv = uv
			this.VtxBuffer[this._VtxWritePtr+0].col = col
			this._VtxWritePtr += 1
		}
		for i := uint(2); i < uint(points_count); i++ {
			this.IdxBuffer[this._IdxWritePtr+0] = (ImDrawIdx)(this._VtxCurrentIdx)
			this.IdxBuffer[this._IdxWritePtr+1] = (ImDrawIdx)(this._VtxCurrentIdx + i - 1)
			this.IdxBuffer[this._IdxWritePtr+2] = (ImDrawIdx)(this._VtxCurrentIdx + i)
			this._IdxWritePtr += 3
		}
		this._VtxCurrentIdx += uint(vtx_count)
	}
} // Note: Anti-aliased filling requires points to be in clockwise order.

func (this *ImDrawList) PathArcToFastEx(center ImVec2, radius float, a_min_sample, a_max_sample, a_step int) {
	if radius <= 0.0 {
		this._Path = append(this._Path, center)
		return
	}

	// Calculate arc auto segment step size
	if a_step <= 0 {
		a_step = IM_DRAWLIST_ARCFAST_SAMPLE_MAX / this._CalcCircleAutoSegmentCount(radius)
	}

	// Make sure we never do steps larger than one quarter of the circle
	a_step = int(ImClamp(float(a_step), 1, IM_DRAWLIST_ARCFAST_TABLE_SIZE/4))

	var sample_range int = ImAbsInt(a_max_sample - a_min_sample)
	var a_next_step = a_step

	var samples int = sample_range + 1
	var extra_max_sample bool = false
	if a_step > 1 {
		samples = sample_range/a_step + 1
		var overstep int = sample_range % a_step

		if overstep > 0 {
			extra_max_sample = true
			samples++

			// When we have overstep to avoid awkwardly looking one long line and one tiny one at the end,
			// distribute first step range evenly between them by reducing first step size.
			if sample_range > 0 {
				a_step -= (a_step - overstep) / 2
			}
		}
	}

	this._Path = append(this._Path, make([]ImVec2, samples)...)
	var out_ptr []ImVec2 = this._Path[(int(len(this._Path)) - samples):]

	var sample_index int = a_min_sample
	if sample_index < 0 || sample_index >= IM_DRAWLIST_ARCFAST_SAMPLE_MAX {
		sample_index = sample_index % IM_DRAWLIST_ARCFAST_SAMPLE_MAX
		if sample_index < 0 {
			sample_index += IM_DRAWLIST_ARCFAST_SAMPLE_MAX
		}
	}

	if a_max_sample >= a_min_sample {
		for a := a_min_sample; a <= a_max_sample; a, sample_index, a_step = a+a_step, sample_index+a_step, a_next_step {
			// a_step is clamped to IM_DRAWLIST_ARCFAST_SAMPLE_MAX, so we have guaranteed that it will not wrap over range twice or more
			if sample_index >= IM_DRAWLIST_ARCFAST_SAMPLE_MAX {
				sample_index -= IM_DRAWLIST_ARCFAST_SAMPLE_MAX
			}

			var s ImVec2 = this._Data.ArcFastVtx[sample_index]
			out_ptr[0].x = center.x + s.x*radius
			out_ptr[0].y = center.y + s.y*radius
			out_ptr = out_ptr[1:]
		}
	} else {
		for a := a_min_sample; a >= a_max_sample; a, sample_index, a_step = a-a_step, sample_index-a_step, a_next_step {
			// a_step is clamped to IM_DRAWLIST_ARCFAST_SAMPLE_MAX, so we have guaranteed that it will not wrap over range twice or more
			if sample_index < 0 {
				sample_index += IM_DRAWLIST_ARCFAST_SAMPLE_MAX
			}

			var s ImVec2 = this._Data.ArcFastVtx[sample_index]
			out_ptr[0].x = center.x + s.x*radius
			out_ptr[0].y = center.y + s.y*radius
			out_ptr = out_ptr[1:]
		}
	}

	if extra_max_sample {
		var normalized_max_sample int = a_max_sample % IM_DRAWLIST_ARCFAST_SAMPLE_MAX
		if normalized_max_sample < 0 {
			normalized_max_sample += IM_DRAWLIST_ARCFAST_SAMPLE_MAX
		}

		var s ImVec2 = this._Data.ArcFastVtx[normalized_max_sample]
		out_ptr[0].x = center.x + s.x*radius
		out_ptr[0].y = center.y + s.y*radius
		out_ptr = out_ptr[1:]
	}

	IM_ASSERT(len(out_ptr) == 0)
} // Use precomputed angles for a 12 steps circle

// Advanced: Primitives allocations
// - We render triangles (three vertices)
// - All primitives needs to be reserved via PrimReserve() beforehand.
// Reserve space for a number of vertices and indices.
// You must finish filling your reserved data before calling PrimReserve() again, as it may reallocate or
// submit the intermediate results. PrimUnreserve() can be used to release unused allocations.
func (this *ImDrawList) PrimReserve(idx_count, vtx_count int) {

	// Large mesh support (when enabled)
	IM_ASSERT(idx_count >= 0 && vtx_count >= 0)
	if unsafe.Sizeof(ImDrawIdx(0)) == 2 && (this._VtxCurrentIdx+uint(vtx_count) >= (1 << 16)) && (this.Flags&ImDrawListFlags_AllowVtxOffset != 0) {
		// FIXME: In theory we should be testing that vtx_count <64k here.
		// In practice, RenderText() relies on reserving ahead for a worst case scenario so it is currently useful for us
		// to not make that check until we rework the text functions to handle clipping and large horizontal lines better.
		this._CmdHeader.VtxOffset = uint(len(this.VtxBuffer))
		this._OnChangedVtxOffset()
	}

	var draw_cmd *ImDrawCmd = &this.CmdBuffer[len(this.CmdBuffer)-1]
	draw_cmd.ElemCount += uint(idx_count)

	var vtx_buffer_old_size int = int(len(this.VtxBuffer))
	this.VtxBuffer = append(this.VtxBuffer, make([]ImDrawVert, vtx_count)...)
	this._VtxWritePtr = vtx_buffer_old_size

	var idx_buffer_old_size int = int(len(this.IdxBuffer))
	this.IdxBuffer = append(this.IdxBuffer, make([]ImDrawIdx, idx_count)...)
	this._IdxWritePtr = idx_buffer_old_size
}

// p_min = upper-left, p_max = lower-right
// Note we don't render 1 pixels sized rectangles properly.
func (this *ImDrawList) AddRect(p_min ImVec2, p_max ImVec2, col ImU32, rounding float, flags ImDrawFlags, thickness float /*= 1.0f*/) {
	if (col & IM_COL32_A_MASK) == 0 {
		return
	}
	min := p_min.Add(ImVec2{0.50, 0.50})
	if this.Flags&ImDrawListFlags_AntiAliasedLines != 0 {
		max := p_max.Sub(ImVec2{0.50, 0.50})
		this.PathRect(&min, &max, rounding, flags)
	} else {
		max := p_max.Sub(ImVec2{0.49, 0.49})
		this.PathRect(&min, &max, rounding, flags)
	} // Better looking lower-right corner and rounded non-AA shapes.
	this.PathStroke(col, ImDrawFlags_Closed, thickness)
} // a: upper-left, b: lower-right (== upper-left + size)

func (this *ImDrawList) AddTriangleFilled(p1 *ImVec2, p2 *ImVec2, p3 ImVec2, col ImU32) {
	if (col & IM_COL32_A_MASK) == 0 {
		return
	}

	this.PathLineTo(*p1)
	this.PathLineTo(*p2)
	this.PathLineTo(p3)
	this.PathFillConvex(col)
}

func AddDrawListToDrawData(out_list *[]*ImDrawList, draw_list *ImDrawList) {
	// Remove trailing command if unused.
	// Technically we could return directly instead of popping, but this make things looks neat in Metrics/Debugger window as well.
	draw_list._PopUnusedDrawCmd()
	if len(draw_list.CmdBuffer) == 0 {
		return
	}

	// Draw list sanity check. Detect mismatch between PrimReserve() calls and incrementing _VtxCurrentIdx, _VtxWritePtr etc.
	// May trigger for you if you are using PrimXXX functions incorrectly.
	IM_ASSERT(len(draw_list.VtxBuffer) == 0 || int(draw_list._VtxWritePtr) == int(len(draw_list.VtxBuffer)))
	IM_ASSERT(len(draw_list.IdxBuffer) == 0 || int(draw_list._IdxWritePtr) == int(len(draw_list.IdxBuffer)))
	if 0 == (draw_list.Flags & ImDrawListFlags_AllowVtxOffset) {
		IM_ASSERT((int)(draw_list._VtxCurrentIdx) == int(len(draw_list.VtxBuffer)))
	}

	// Check that draw_list doesn't use more vertices than indexable (default ImDrawIdx = unsigned short = 2 bytes = 64K vertices per ImDrawList = per window)
	// If this assert triggers because you are drawing lots of stuff manually:
	// - First, make sure you are coarse clipping yourself and not trying to draw many things outside visible bounds.
	//   Be mindful that the ImDrawList API doesn't filter vertices. Use the Metrics/Debugger window to inspect draw list contents.
	// - If you want large meshes with more than 64K vertices, you can either:
	//   (A) Handle the ImDrawCmd::VtxOffset value in your renderer backend, and set 'io.BackendFlags |= ImGuiBackendFlags_RendererHasVtxOffset'.
	//       Most example backends already support this from 1.71. Pre-1.71 backends won't.
	//       Some graphics API such as GL ES 1/2 don't have a way to offset the starting vertex so it is not supported for them.
	//   (B) Or handle 32-bit indices in your renderer backend, and uncomment '#define ImDrawIdx unsigned int' line in imconfig.h.
	//       Most example backends already support this. For example, the OpenGL example code detect index size at compile-time:
	//         glDrawElements(GL_TRIANGLES, (GLsizei)pcmd.ElemCount, sizeof(ImDrawIdx) == 2 ? GL_UNSIGNED_SHORT : GL_UNSIGNED_INT, idx_buffer_offset);
	//       Your own engine or render API may use different parameters or function calls to specify index sizes.
	//       2 and 4 bytes indices are generally supported by most graphics API.
	// - If for some reason neither of those solutions works for you, a workaround is to call BeginChild()/EndChild() before reaching
	//   the 64K limit to split your draw commands in multiple draw lists.
	if unsafe.Sizeof(ImDrawIdx(0)) == 2 {
		IM_ASSERT_USER_ERROR(draw_list._VtxCurrentIdx < (1<<16), "Too many vertices in ImDrawList using 16-bit indices. Read comment above")
	}

	*out_list = append(*out_list, draw_list)
}

// Fully unrolled with inline call to keep our debug builds decently fast.
func (this *ImDrawList) PrimRect(a, c *ImVec2, col ImU32) {

	var b, d, uv ImVec2 = ImVec2{c.x, a.y}, ImVec2{a.x, c.y}, this._Data.TexUvWhitePixel
	var idx ImDrawIdx = (ImDrawIdx)(this._VtxCurrentIdx)
	this.IdxBuffer[this._IdxWritePtr+0] = idx
	this.IdxBuffer[this._IdxWritePtr+1] = (ImDrawIdx)(idx + 1)
	this.IdxBuffer[this._IdxWritePtr+2] = (ImDrawIdx)(idx + 2)
	this.IdxBuffer[this._IdxWritePtr+3] = idx
	this.IdxBuffer[this._IdxWritePtr+4] = (ImDrawIdx)(idx + 2)
	this.IdxBuffer[this._IdxWritePtr+5] = (ImDrawIdx)(idx + 3)
	this.VtxBuffer[this._VtxWritePtr+0].pos = *a
	this.VtxBuffer[this._VtxWritePtr+0].uv = uv
	this.VtxBuffer[this._VtxWritePtr+0].col = col
	this.VtxBuffer[this._VtxWritePtr+1].pos = b
	this.VtxBuffer[this._VtxWritePtr+1].uv = uv
	this.VtxBuffer[this._VtxWritePtr+1].col = col
	this.VtxBuffer[this._VtxWritePtr+2].pos = *c
	this.VtxBuffer[this._VtxWritePtr+2].uv = uv
	this.VtxBuffer[this._VtxWritePtr+2].col = col
	this.VtxBuffer[this._VtxWritePtr+3].pos = d
	this.VtxBuffer[this._VtxWritePtr+3].uv = uv
	this.VtxBuffer[this._VtxWritePtr+3].col = col
	this._VtxWritePtr += 4
	this._VtxCurrentIdx += 4
	this._IdxWritePtr += 6
}

// TODO: Thickness anti-aliased lines cap are missing their AA fringe.
// We avoid using the ImVec2 math operators here to reduce cost to a minimum for debug/non-inlined builds.
func (this *ImDrawList) AddPolyline(points []ImVec2, points_count int, col ImU32, flags ImDrawFlags, thickness float) {
	if points_count < 2 {
		return
	}

	var closed bool = (flags & ImDrawFlags_Closed) != 0
	var opaque_uv ImVec2 = this._Data.TexUvWhitePixel
	var count int // The number of line segments we need to draw
	if closed {
		count = points_count
	} else {
		count = points_count - 1
	}
	var thick_line bool = (thickness > this._FringeScale)

	if this.Flags&ImDrawListFlags_AntiAliasedLines != 0 {
		// Anti-aliased stroke
		var AA_SIZE float = this._FringeScale
		var col_trans ImU32 = col &^ IM_COL32_A_MASK

		// Thicknesses <1.0 should behave like thickness 1.0
		thickness = ImMax(thickness, 1.0)
		var integer_thickness int = (int)(thickness)
		var fractional_thickness float = thickness - float(integer_thickness)

		// Do we want to draw this line using a texture?
		// - For now, only draw integer-width lines using textures to avoid issues with the way scaling occurs, could be improved.
		// - If AA_SIZE is not 1.0f we cannot use the texture path.
		var use_texture bool = (this.Flags&ImDrawListFlags_AntiAliasedLinesUseTex != 0) && (integer_thickness < IM_DRAWLIST_TEX_LINES_WIDTH_MAX) && (fractional_thickness <= 0.00001) && (AA_SIZE == 1.0)

		// We should never hit this, because NewFrame() doesn't set ImDrawListFlags_AntiAliasedLinesUseTex unless ImFontAtlasFlags_NoBakedLines is off
		IM_ASSERT(!use_texture || 0 == (this._Data.Font.ContainerAtlas.Flags&ImFontAtlasFlags_NoBakedLines))

		var idx_count int
		var vtx_count int
		if use_texture {
			// Texture line
			idx_count = count * 6
			vtx_count = points_count * 2
		} else {
			if thick_line {
				// Thick anti-aliased lines
				idx_count = count * 18
				vtx_count = points_count * 4
			} else {
				// Anti-aliased lines
				idx_count = count * 12
				vtx_count = points_count * 3
			}
		}

		this.PrimReserve(idx_count, vtx_count)

		var num_normals int = 5
		if use_texture || !thick_line {
			num_normals = 3
		}

		// Temporary buffer
		// The first <points_count> items are normals at each line point, then after that there are either 2 or 4 temp points for each line point
		var temp_normals []ImVec2 = make([]ImVec2, points_count*num_normals) //-V630
		var temp_points []ImVec2 = temp_normals[points_count:]

		// Calculate normals (tangents) for each line segment
		for i1 := int(0); i1 < count; i1++ {

			var i2 int
			if (i1 + 1) != points_count {
				i2 = i1 + 1
			}
			var dx float = points[i2].x - points[i1].x
			var dy float = points[i2].y - points[i1].y
			IM_NORMALIZE2F_OVER_ZERO(&dx, &dy)
			temp_normals[i1].x = dy
			temp_normals[i1].y = -dx
		}
		if !closed {
			temp_normals[points_count-1] = temp_normals[points_count-2]
		}

		// If we are drawing a one-pixel-wide line without a texture, or a textured line of any width, we only need 2 or 3 vertices per point
		if use_texture || !thick_line {
			// [PATH 1] Texture-based lines (thick or non-thick)
			// [PATH 2] Non texture-based lines (non-thick)

			// The width of the geometry we need to draw - this is essentially <thickness> pixels for the line itself, plus "one pixel" for AA.
			// - In the texture-based path, we don't use AA_SIZE here because the +1 is tied to the generated texture
			//   (see ImFontAtlasBuildRenderLinesTexData() function), and so alternate values won't work without changes to that code.
			// - In the non texture-based paths, we would allow AA_SIZE to potentially be != 1.0f with a patch (e.g. fringe_scale patch to
			//   allow scaling geometry while preserving one-screen-pixel AA fringe).
			var half_draw_size float = AA_SIZE
			if use_texture {
				half_draw_size = thickness*0.5 + 1
			}

			// If line is not closed, the first and last points need to be generated differently as there are no normals to blend
			if !closed {
				temp_points[0] = points[0].Add(temp_normals[0].Scale(half_draw_size))
				temp_points[1] = points[0].Sub(temp_normals[0].Scale(half_draw_size))
				temp_points[(points_count-1)*2+0] = points[points_count-1].Add(temp_normals[points_count-1].Scale(half_draw_size))
				temp_points[(points_count-1)*2+1] = points[points_count-1].Sub(temp_normals[points_count-1].Scale(half_draw_size))
			}

			// Generate the indices to form a number of triangles for each line segment, and the vertices for the line edges
			// This takes points n and n+1 and writes into n+1, with the first point in a closed line being generated from the final one (as n+1 wraps)
			// FIXME-OPT: Merge the different loops, possibly remove the temporary buffer.
			var idx1 uint = this._VtxCurrentIdx  // Vertex index for start of line segment
			for i1 := int(0); i1 < count; i1++ { // i1 is the first point of the line segment

				var i2 int // i2 is the second point of the line segment
				if (i1 + 1) != points_count {
					i2 = i1 + 1
				}

				var idx2 uint // Vertex index for end of segment
				if (i1 + 1) == points_count {
					idx2 = this._VtxCurrentIdx
				} else if use_texture {
					idx2 = idx1 + 2
				} else {
					idx2 = idx1 + 3
				}

				//printf("i1, i2, idx2 %v %v %v, this._VtxCurrentIdx %v\n ", i1, i2, idx2, this._VtxCurrentIdx)

				// Average normals
				var dm_x float = (temp_normals[i1].x + temp_normals[i2].x) * 0.5
				var dm_y float = (temp_normals[i1].y + temp_normals[i2].y) * 0.5
				IM_FIXNORMAL2F(&dm_x, &dm_y)
				dm_x *= half_draw_size // dm_x, dm_y are offset to the outer edge of the AA area
				dm_y *= half_draw_size

				// Add temporary vertexes for the outer edges
				var out_vtx []ImVec2 = temp_points[i2*2:]
				out_vtx[0].x = points[i2].x + dm_x
				out_vtx[0].y = points[i2].y + dm_y
				out_vtx[1].x = points[i2].x - dm_x
				out_vtx[1].y = points[i2].y - dm_y

				if use_texture {
					// Add indices for two triangles
					this.IdxBuffer[this._IdxWritePtr+0] = (ImDrawIdx)(idx2 + 0)
					this.IdxBuffer[this._IdxWritePtr+1] = (ImDrawIdx)(idx1 + 0)
					this.IdxBuffer[this._IdxWritePtr+2] = (ImDrawIdx)(idx1 + 1) // Right tri
					this.IdxBuffer[this._IdxWritePtr+3] = (ImDrawIdx)(idx2 + 1)
					this.IdxBuffer[this._IdxWritePtr+4] = (ImDrawIdx)(idx1 + 1)
					this.IdxBuffer[this._IdxWritePtr+5] = (ImDrawIdx)(idx2 + 0) // Left tri
					this._IdxWritePtr += 6
				} else {
					// Add indexes for four triangles
					this.IdxBuffer[this._IdxWritePtr+0] = (ImDrawIdx)(idx2 + 0)
					this.IdxBuffer[this._IdxWritePtr+1] = (ImDrawIdx)(idx1 + 0)
					this.IdxBuffer[this._IdxWritePtr+2] = (ImDrawIdx)(idx1 + 2) // Right tri 1
					this.IdxBuffer[this._IdxWritePtr+3] = (ImDrawIdx)(idx1 + 2)
					this.IdxBuffer[this._IdxWritePtr+4] = (ImDrawIdx)(idx2 + 2)
					this.IdxBuffer[this._IdxWritePtr+5] = (ImDrawIdx)(idx2 + 0) // Right tri 2
					this.IdxBuffer[this._IdxWritePtr+6] = (ImDrawIdx)(idx2 + 1)
					this.IdxBuffer[this._IdxWritePtr+7] = (ImDrawIdx)(idx1 + 1)
					this.IdxBuffer[this._IdxWritePtr+8] = (ImDrawIdx)(idx1 + 0) // Left tri 1
					this.IdxBuffer[this._IdxWritePtr+9] = (ImDrawIdx)(idx1 + 0)
					this.IdxBuffer[this._IdxWritePtr+10] = (ImDrawIdx)(idx2 + 0)
					this.IdxBuffer[this._IdxWritePtr+11] = (ImDrawIdx)(idx2 + 1) // Left tri 2
					this._IdxWritePtr += 12
				}

				idx1 = idx2
			}

			// Add vertexes for each point on the line
			if use_texture {
				// If we're using textures we only need to emit the left/right edge vertices
				var tex_uvs ImVec4 = this._Data.TexUvLines[integer_thickness]
				/*if (fractional_thickness != 0.0f) // Currently always zero when use_texture==false!
				  {
				      const ImVec4 tex_uvs_1 = _Data.TexUvLines[integer_thickness + 1];
				      tex_uvs.x = tex_uvs.x + (tex_uvs_1.x - tex_uvs.x) * fractional_thickness; // inlined ImLerp()
				      tex_uvs.y = tex_uvs.y + (tex_uvs_1.y - tex_uvs.y) * fractional_thickness;
				      tex_uvs.z = tex_uvs.z + (tex_uvs_1.z - tex_uvs.z) * fractional_thickness;
				      tex_uvs.w = tex_uvs.w + (tex_uvs_1.w - tex_uvs.w) * fractional_thickness;
				  }*/
				var tex_uv0 = ImVec2{tex_uvs.x, tex_uvs.y}
				var tex_uv1 = ImVec2{tex_uvs.z, tex_uvs.w}
				for i := int(0); i < points_count; i++ {
					this.VtxBuffer[this._VtxWritePtr+0].pos = temp_points[i*2+0]
					this.VtxBuffer[this._VtxWritePtr+0].uv = tex_uv0
					this.VtxBuffer[this._VtxWritePtr+0].col = col // Left-side outer edge
					this.VtxBuffer[this._VtxWritePtr+1].pos = temp_points[i*2+1]
					this.VtxBuffer[this._VtxWritePtr+1].uv = tex_uv1
					this.VtxBuffer[this._VtxWritePtr+1].col = col // Right-side outer edge
					this._VtxWritePtr += 2
				}
			} else {
				// If we're not using a texture, we need the center vertex as well
				for i := int(0); i < points_count; i++ {
					this.VtxBuffer[this._VtxWritePtr+0].pos = points[i]
					this.VtxBuffer[this._VtxWritePtr+0].uv = opaque_uv
					this.VtxBuffer[this._VtxWritePtr+0].col = col // Center of line
					this.VtxBuffer[this._VtxWritePtr+1].pos = temp_points[i*2+0]
					this.VtxBuffer[this._VtxWritePtr+1].uv = opaque_uv
					this.VtxBuffer[this._VtxWritePtr+1].col = col_trans // Left-side outer edge
					this.VtxBuffer[this._VtxWritePtr+2].pos = temp_points[i*2+1]
					this.VtxBuffer[this._VtxWritePtr+2].uv = opaque_uv
					this.VtxBuffer[this._VtxWritePtr+2].col = col_trans // Right-side outer edge
					this._VtxWritePtr += 3
				}
			}
		} else {
			// [PATH 2] Non texture-based lines (thick): we need to draw the solid line core and thus require four vertices per point
			var half_inner_thickness float = (thickness - AA_SIZE) * 0.5

			// If line is not closed, the first and last points need to be generated differently as there are no normals to blend
			if !closed {
				var points_last int = points_count - 1
				temp_points[0] = points[0].Add(temp_normals[0].Scale(half_inner_thickness + AA_SIZE))
				temp_points[1] = points[0].Add(temp_normals[0].Scale(half_inner_thickness))
				temp_points[2] = points[0].Sub(temp_normals[0].Scale(half_inner_thickness))
				temp_points[3] = points[0].Sub(temp_normals[0].Scale(half_inner_thickness + AA_SIZE))
				temp_points[points_last*4+0] = points[points_last].Add(temp_normals[points_last].Scale(half_inner_thickness + AA_SIZE))
				temp_points[points_last*4+1] = points[points_last].Add(temp_normals[points_last].Scale(half_inner_thickness))
				temp_points[points_last*4+2] = points[points_last].Sub(temp_normals[points_last].Scale(half_inner_thickness))
				temp_points[points_last*4+3] = points[points_last].Sub(temp_normals[points_last].Scale(half_inner_thickness + AA_SIZE))
			}

			// Generate the indices to form a number of triangles for each line segment, and the vertices for the line edges
			// This takes points n and n+1 and writes into n+1, with the first point in a closed line being generated from the final one (as n+1 wraps)
			// FIXME-OPT: Merge the different loops, possibly remove the temporary buffer.
			var idx1 uint = this._VtxCurrentIdx  // Vertex index for start of line segment
			for i1 := int(0); i1 < count; i1++ { // i1 is the first point of the line segment

				var i2 int // i2 is the second point of the line segment
				if (i1 + 1) != points_count {
					i2 = (i1 + 1)
				}

				var idx2 uint // Vertex index for end of segment
				if (i1 + 1) == points_count {
					idx2 = this._VtxCurrentIdx
				} else {
					idx2 = (idx1 + 4)
				}

				// Average normals
				var dm_x float = (temp_normals[i1].x + temp_normals[i2].x) * 0.5
				var dm_y float = (temp_normals[i1].y + temp_normals[i2].y) * 0.5
				IM_FIXNORMAL2F(&dm_x, &dm_y)
				var dm_out_x float = dm_x * (half_inner_thickness + AA_SIZE)
				var dm_out_y float = dm_y * (half_inner_thickness + AA_SIZE)
				var dm_in_x float = dm_x * half_inner_thickness
				var dm_in_y float = dm_y * half_inner_thickness

				// Add temporary vertices
				var out_vtx []ImVec2 = temp_points[i2*4:]
				out_vtx[0].x = points[i2].x + dm_out_x
				out_vtx[0].y = points[i2].y + dm_out_y
				out_vtx[1].x = points[i2].x + dm_in_x
				out_vtx[1].y = points[i2].y + dm_in_y
				out_vtx[2].x = points[i2].x - dm_in_x
				out_vtx[2].y = points[i2].y - dm_in_y
				out_vtx[3].x = points[i2].x - dm_out_x
				out_vtx[3].y = points[i2].y - dm_out_y

				// Add indexes
				this.IdxBuffer[this._IdxWritePtr+0] = (ImDrawIdx)(idx2 + 1)
				this.IdxBuffer[this._IdxWritePtr+1] = (ImDrawIdx)(idx1 + 1)
				this.IdxBuffer[this._IdxWritePtr+2] = (ImDrawIdx)(idx1 + 2)
				this.IdxBuffer[this._IdxWritePtr+3] = (ImDrawIdx)(idx1 + 2)
				this.IdxBuffer[this._IdxWritePtr+4] = (ImDrawIdx)(idx2 + 2)
				this.IdxBuffer[this._IdxWritePtr+5] = (ImDrawIdx)(idx2 + 1)
				this.IdxBuffer[this._IdxWritePtr+6] = (ImDrawIdx)(idx2 + 1)
				this.IdxBuffer[this._IdxWritePtr+7] = (ImDrawIdx)(idx1 + 1)
				this.IdxBuffer[this._IdxWritePtr+8] = (ImDrawIdx)(idx1 + 0)
				this.IdxBuffer[this._IdxWritePtr+9] = (ImDrawIdx)(idx1 + 0)
				this.IdxBuffer[this._IdxWritePtr+10] = (ImDrawIdx)(idx2 + 0)
				this.IdxBuffer[this._IdxWritePtr+11] = (ImDrawIdx)(idx2 + 1)
				this.IdxBuffer[this._IdxWritePtr+12] = (ImDrawIdx)(idx2 + 2)
				this.IdxBuffer[this._IdxWritePtr+13] = (ImDrawIdx)(idx1 + 2)
				this.IdxBuffer[this._IdxWritePtr+14] = (ImDrawIdx)(idx1 + 3)
				this.IdxBuffer[this._IdxWritePtr+15] = (ImDrawIdx)(idx1 + 3)
				this.IdxBuffer[this._IdxWritePtr+16] = (ImDrawIdx)(idx2 + 3)
				this.IdxBuffer[this._IdxWritePtr+17] = (ImDrawIdx)(idx2 + 2)
				this._IdxWritePtr += 18

				idx1 = idx2
			}

			// Add vertices
			for i := int(0); i < points_count; i++ {
				this.VtxBuffer[this._VtxWritePtr+0].pos = temp_points[i*4+0]
				this.VtxBuffer[this._VtxWritePtr+0].uv = opaque_uv
				this.VtxBuffer[this._VtxWritePtr+0].col = col_trans
				this.VtxBuffer[this._VtxWritePtr+1].pos = temp_points[i*4+1]
				this.VtxBuffer[this._VtxWritePtr+1].uv = opaque_uv
				this.VtxBuffer[this._VtxWritePtr+1].col = col
				this.VtxBuffer[this._VtxWritePtr+2].pos = temp_points[i*4+2]
				this.VtxBuffer[this._VtxWritePtr+2].uv = opaque_uv
				this.VtxBuffer[this._VtxWritePtr+2].col = col
				this.VtxBuffer[this._VtxWritePtr+3].pos = temp_points[i*4+3]
				this.VtxBuffer[this._VtxWritePtr+3].uv = opaque_uv
				this.VtxBuffer[this._VtxWritePtr+3].col = col_trans
				this._VtxWritePtr += 4
			}
		}
		this._VtxCurrentIdx += uint((ImDrawIdx)(vtx_count))
	} else {

		// [PATH 4] Non texture-based, Non anti-aliased lines
		var idx_count int = count * 6
		var vtx_count int = count * 4 // FIXME-OPT: Not sharing edges
		this.PrimReserve(idx_count, vtx_count)

		for i1 := int(0); i1 < count; i1++ {
			var i2 int
			if (i1 + 1) != points_count {
				i2 = (i1 + 1)
			}

			var p1 *ImVec2 = &points[i1]
			var p2 *ImVec2 = &points[i2]

			var dx float = p2.x - p1.x
			var dy float = p2.y - p1.y
			IM_NORMALIZE2F_OVER_ZERO(&dx, &dy)
			dx *= (thickness * 0.5)
			dy *= (thickness * 0.5)

			this.VtxBuffer[this._VtxWritePtr+0].pos.x = p1.x + dy
			this.VtxBuffer[this._VtxWritePtr+0].pos.y = p1.y - dx
			this.VtxBuffer[this._VtxWritePtr+0].uv = opaque_uv
			this.VtxBuffer[this._VtxWritePtr+0].col = col
			this.VtxBuffer[this._VtxWritePtr+1].pos.x = p2.x + dy
			this.VtxBuffer[this._VtxWritePtr+1].pos.y = p2.y - dx
			this.VtxBuffer[this._VtxWritePtr+1].uv = opaque_uv
			this.VtxBuffer[this._VtxWritePtr+1].col = col
			this.VtxBuffer[this._VtxWritePtr+2].pos.x = p2.x - dy
			this.VtxBuffer[this._VtxWritePtr+2].pos.y = p2.y + dx
			this.VtxBuffer[this._VtxWritePtr+2].uv = opaque_uv
			this.VtxBuffer[this._VtxWritePtr+2].col = col
			this.VtxBuffer[this._VtxWritePtr+3].pos.x = p1.x - dy
			this.VtxBuffer[this._VtxWritePtr+3].pos.y = p1.y + dx
			this.VtxBuffer[this._VtxWritePtr+3].uv = opaque_uv
			this.VtxBuffer[this._VtxWritePtr+3].col = col
			this._VtxWritePtr += 4

			this.IdxBuffer[this._IdxWritePtr+0] = (ImDrawIdx)(this._VtxCurrentIdx)
			this.IdxBuffer[this._IdxWritePtr+1] = (ImDrawIdx)(this._VtxCurrentIdx + 1)
			this.IdxBuffer[this._IdxWritePtr+2] = (ImDrawIdx)(this._VtxCurrentIdx + 2)
			this.IdxBuffer[this._IdxWritePtr+3] = (ImDrawIdx)(this._VtxCurrentIdx)
			this.IdxBuffer[this._IdxWritePtr+4] = (ImDrawIdx)(this._VtxCurrentIdx + 2)
			this.IdxBuffer[this._IdxWritePtr+5] = (ImDrawIdx)(this._VtxCurrentIdx + 3)
			this._IdxWritePtr += 6
			this._VtxCurrentIdx += 4
		}
	}
}

func (this *ImDrawList) AddCircleFilled(center ImVec2, radius float, col ImU32, num_segments int) {
	if (col&IM_COL32_A_MASK) == 0 || radius <= 0.0 {
		return
	}

	if num_segments <= 0 {
		// Use arc with automatic segment count
		this.PathArcToFastEx(center, radius, 0, IM_DRAWLIST_ARCFAST_SAMPLE_MAX, 0)
		this._Path = this._Path[:len(this._Path)-1]
	} else {
		// Explicit segment count (still clamp to avoid drawing insanely tessellated shapes)
		num_segments = int(ImClamp(float(num_segments), 3, IM_DRAWLIST_CIRCLE_AUTO_SEGMENT_MAX))

		// Because we are filling a closed shape we remove 1 from the count of segments/points
		var a_max float = (IM_PI * 2.0) * ((float)(num_segments) - 1.0) / (float)(num_segments)
		this.PathArcTo(center, radius, 0.0, a_max, num_segments-1)
	}

	this.PathFillConvex(col)
}

func (this *ImDrawList) PathArcTo(center ImVec2, radius, a_min, a_max float, num_segments int) {
	if radius <= 0.0 {
		this._Path = append(this._Path, center)
		return
	}

	if num_segments > 0 {
		this._PathArcToN(center, radius, a_min, a_max, num_segments)
		return
	}

	// Automatic segment count
	if radius <= this._Data.ArcFastRadiusCutoff {
		var a_is_reverse bool = a_max < a_min

		// We are going to use precomputed values for mid samples.
		// Determine first and last sample in lookup table that belong to the arc.
		var a_min_sample_f float = IM_DRAWLIST_ARCFAST_SAMPLE_MAX * a_min / (IM_PI * 2.0)
		var a_max_sample_f float = IM_DRAWLIST_ARCFAST_SAMPLE_MAX * a_max / (IM_PI * 2.0)

		var a_min_sample int
		if a_is_reverse {
			a_min_sample = (int)(ImFloorSigned(a_min_sample_f))
		} else {
			a_min_sample = (int)(ImCeil(a_min_sample_f))
		}

		var a_max_sample int
		if a_is_reverse {
			a_max_sample = (int)(ImCeil(a_max_sample_f))
		} else {
			a_max_sample = (int)(ImFloorSigned(a_max_sample_f))
		}

		var a_mid_samples int
		if a_is_reverse {
			a_mid_samples = ImMaxInt(a_min_sample-a_max_sample, 0)
		} else {
			a_mid_samples = ImMaxInt(a_max_sample-a_min_sample, 0)
		}

		var a_min_segment_angle float = float(a_min_sample) * IM_PI * 2.0 / IM_DRAWLIST_ARCFAST_SAMPLE_MAX
		var a_max_segment_angle float = float(a_max_sample) * IM_PI * 2.0 / IM_DRAWLIST_ARCFAST_SAMPLE_MAX
		var a_emit_start bool = (a_min_segment_angle - a_min) != 0.0
		var a_emit_end bool = (a_max - a_max_segment_angle) != 0.0

		var emit int
		if a_emit_start {
			emit += 1
		}
		if a_emit_end {
			emit += 1
		}

		//grow slice if necessary (this._Path.reserve(_Path.Size + (a_mid_samples + 1 + emit)))
		this._Path = reserveVec2Slice(this._Path, int(len(this._Path))+(a_mid_samples+1+emit))

		if a_emit_start {
			this._Path = append(this._Path, ImVec2{center.x + ImCos(a_min)*radius, center.y + ImSin(a_min)*radius})
		}
		if a_mid_samples > 0 {
			this.PathArcToFastEx(center, radius, a_min_sample, a_max_sample, 0)
		}
		if a_emit_end {
			this._Path = append(this._Path, ImVec2{center.x + ImCos(a_max)*radius, center.y + ImSin(a_max)*radius})
		}
	} else {
		var arc_length float = ImAbs(a_max - a_min)
		var circle_segment_count int = this._CalcCircleAutoSegmentCount(radius)
		var arc_segment_count int = ImMaxInt((int)(ImCeil(float(circle_segment_count)*arc_length/(IM_PI*2.0))), (int)(2.0*IM_PI/arc_length))
		this._PathArcToN(center, radius, a_min, a_max, arc_segment_count)
	}
}
