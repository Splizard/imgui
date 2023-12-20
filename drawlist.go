package imgui

// ImDrawList Draw command list
// This is the low-level list of polygons that ImGui:: functions are filling. At the end of the frame,
// all command lists are passed to your ImGuiIO::RenderDrawListFn function for rendering.
// Each dear imgui window contains its own ImDrawList. You can use ImGui::GetWindowDrawList() to
// access the current window draw list and draw custom primitives.
// You can interleave normal ImGui:: calls and adding primitives to the current draw list.
// In single viewport mode, top-left is == GetMainViewport()->Pos (generally 0,0), bottom-right is == GetMainViewport()->Pos+Size (generally io.DisplaySize).
// You are totally free to apply whatever transformation matrix to want to the data (depending on the use of the transformation you may want to apply it to ClipRect as well!)
// Important: Primitives are always added to the list and not culled (culling is done at higher-level by ImGui:: functions), if you use this API a lot consider coarse culling your drawn objects.
type ImDrawList struct {
	// This is what you have to render
	CmdBuffer []ImDrawCmd     // Draw commands. Typically 1 command = 1 GPU draw call, unless the command is a callback.
	IdxBuffer []ImDrawIdx     // Index buffer. Each command consume ImDrawCmd::ElemCount of those
	VtxBuffer []ImDrawVert    // Vertex buffer.
	Flags     ImDrawListFlags // Flags, you may poke into these to adjust anti-aliasing settings per-primitive.

	// [Internal, used while building lists]
	_VtxCurrentIdx  uint                  // [Internal] generally == VtxBuffer.Size unless we are past 64K vertices, in which case this gets reset to 0.
	_Data           *ImDrawListSharedData // Pointer to shared draw data (you can use ImGui::GetDrawListSharedData() to get the one from current ImGui context)
	_OwnerName      string                // Pointer to owner window's name for debugging
	_VtxWritePtr    int                   // [Internal] point within VtxBuffer.Data after each add command (to avoid using the ImVector<> operators too much)
	_IdxWritePtr    int                   // [Internal] point within IdxBuffer.Data after each add command (to avoid using the ImVector<> operators too much)
	_ClipRectStack  []ImVec4              // [Internal]
	_TextureIdStack []ImTextureID         // [Internal]
	_Path           []ImVec2              // [Internal] current path building
	_CmdHeader      ImDrawCmdHeader       // [Internal] template of active commands. Fields should match those of CmdBuffer.back().
	_Splitter       ImDrawListSplitter    // [Internal] for channels api (note: prefer using your own persistent instance of ImDrawListSplitter!)
	_FringeScale    float                 // [Internal] anti-alias fringe is scaled by this value, this helps to keep things sharp while zooming at vertex buffer content
}

func NewImDrawList(shared_data *ImDrawListSharedData) ImDrawList {
	return ImDrawList{
		_Data: shared_data,
	}
}

func (this *ImDrawList) PushClipRect(cr_min, cr_max ImVec2, intersect_with_current_clip_rect bool) {
	var cr = ImVec4{cr_min.x, cr_min.y, cr_max.x, cr_max.y}
	if intersect_with_current_clip_rect {
		var current = this._CmdHeader.ClipRect
		if cr.x < current.x {
			cr.x = current.x
		}
		if cr.y < current.y {
			cr.y = current.y
		}
		if cr.z > current.z {
			cr.z = current.z
		}
		if cr.w > current.w {
			cr.w = current.w
		}
	}
	cr.z = ImMax(cr.x, cr.z)
	cr.w = ImMax(cr.y, cr.w)

	this._ClipRectStack = append(this._ClipRectStack, cr)
	this._CmdHeader.ClipRect = cr
	this._OnChangedClipRect()
}

// PushClipRectFullScreen Render-level scissoring. This is passed down to your render function but not used for CPU-side coarse clipping. Prefer using higher-level ImGui::PushClipRect() to affect logic (hit-testing and widget culling)
func (this *ImDrawList) PushClipRectFullScreen() {
	this.PushClipRect(ImVec2{this._Data.ClipRectFullscreen.x, this._Data.ClipRectFullscreen.y}, ImVec2{this._Data.ClipRectFullscreen.z, this._Data.ClipRectFullscreen.w}, false)
}

func (this *ImDrawList) PopClipRect() {
	this._ClipRectStack = this._ClipRectStack[len(this._ClipRectStack)-1:]
	if len(this._ClipRectStack) == 0 {
		this._CmdHeader.ClipRect = this._Data.ClipRectFullscreen
	} else {
		this._CmdHeader.ClipRect = this._ClipRectStack[len(this._ClipRectStack)-1]
	}
	this._OnChangedClipRect()
}

func (this *ImDrawList) GetClipRectMin() ImVec2 {
	var cr = &this._ClipRectStack[len(this._ClipRectStack)-1]
	return ImVec2{cr.x, cr.y}
}
func (this *ImDrawList) GetClipRectMax() ImVec2 {
	var cr = &this._ClipRectStack[len(this._ClipRectStack)-1]
	return ImVec2{cr.x, cr.y}
}

// AddLine Primitives
//   - For rectangular primitives, "p_min" and "p_max" represent the upper-left and lower-right corners.
//   - For circle primitives, use "num_segments == 0" to automatically calculate tessellation (preferred).
//     In older versions (until Dear ImGui 1.77) the AddCircle functions defaulted to num_segments == 12.
//     In future versions we will use textures to provide cheaper and higher-quality circles.
//     Use AddNgon() and AddNgonFilled() functions if you need to guaranteed a specific number of sides.
func (this *ImDrawList) AddLine(p1 *ImVec2, p2 *ImVec2, col ImU32, thickness float /*= 1.0f*/) {
	if (col & IM_COL32_A_MASK) == 0 {
		return
	}
	this.PathLineTo(p1.Add(ImVec2{0.5, 0.5}))
	this.PathLineTo(p2.Add(ImVec2{0.5, 0.5}))
	this.PathStroke(col, 0, thickness)
}

// AddRectFilledMultiColor p_min = upper-left, p_max = lower-right
func (this *ImDrawList) AddRectFilledMultiColor(p_min ImVec2, p_max ImVec2, col_upr_left, col_upr_right, col_bot_right, col_bot_left ImU32) {
	if ((col_upr_left | col_upr_right | col_bot_right | col_bot_left) & IM_COL32_A_MASK) == 0 {
		return
	}

	var uv = this._Data.TexUvWhitePixel
	this.PrimReserve(6, 4)
	this.PrimWriteIdx((ImDrawIdx)(this._VtxCurrentIdx))
	this.PrimWriteIdx((ImDrawIdx)(this._VtxCurrentIdx + 1))
	this.PrimWriteIdx((ImDrawIdx)(this._VtxCurrentIdx + 2))
	this.PrimWriteIdx((ImDrawIdx)(this._VtxCurrentIdx))
	this.PrimWriteIdx((ImDrawIdx)(this._VtxCurrentIdx + 2))
	this.PrimWriteIdx((ImDrawIdx)(this._VtxCurrentIdx + 3))
	this.PrimWriteVtx(p_min, &uv, col_upr_left)
	this.PrimWriteVtx(ImVec2{p_max.x, p_min.y}, &uv, col_upr_right)
	this.PrimWriteVtx(p_max, &uv, col_bot_right)
	this.PrimWriteVtx(ImVec2{p_min.x, p_max.y}, &uv, col_bot_left)
}

func (this *ImDrawList) AddQuad(p1 *ImVec2, p2 *ImVec2, p3 ImVec2, p4 ImVec2, col ImU32, thickness float /*= 1.0f*/) {
	if (col & IM_COL32_A_MASK) == 0 {
		return
	}

	this.PathLineTo(*p1)
	this.PathLineTo(*p2)
	this.PathLineTo(p3)
	this.PathLineTo(p4)
	this.PathStroke(col, ImDrawFlags_Closed, thickness)
}

func (this *ImDrawList) AddQuadFilled(p1 *ImVec2, p2 *ImVec2, p3 ImVec2, p4 ImVec2, col ImU32) {
	if (col & IM_COL32_A_MASK) == 0 {
		return
	}

	this.PathLineTo(*p1)
	this.PathLineTo(*p2)
	this.PathLineTo(p3)
	this.PathLineTo(p4)
	this.PathFillConvex(col)
}

func (this *ImDrawList) AddTriangle(p1 *ImVec2, p2 *ImVec2, p3 ImVec2, col ImU32, thickness float /*= 1.0f*/) {
	if (col & IM_COL32_A_MASK) == 0 {
		return
	}

	this.PathLineTo(*p1)
	this.PathLineTo(*p2)
	this.PathLineTo(p3)
	this.PathStroke(col, ImDrawFlags_Closed, thickness)
}

func (this *ImDrawList) AddTriangleFilled(p1 *ImVec2, p2 *ImVec2, p3 ImVec2, col ImU32) {
	if (col & IM_COL32_A_MASK) == 0 {
		return
	}

	this.PathLineTo(*p1)
	this.PathLineTo(*p2)
	this.PathLineTo(p3)
	this.PathFillConvex(col)
}

func (this *ImDrawList) AddCircle(center ImVec2, radius float, col ImU32, num_segments int, thickness float /*= 1.0f*/) {
	if col&IM_COL32_A_MASK == 0 || radius < 0.0 {
		return
	}

	if num_segments <= 0 {
		// Use arc with automatic segment count
		this.PathArcToFastEx(center, radius-0.5, 0, IM_DRAWLIST_ARCFAST_SAMPLE_MAX, 0)
		this._Path = this._Path[:len(this._Path)-1]
	} else {
		// Explicit segment count (still clamp to avoid drawing insanely tessellated shapes)
		num_segments = ImClampInt(num_segments, 3, IM_DRAWLIST_CIRCLE_AUTO_SEGMENT_MAX)

		// Because we are filling a closed shape we remove 1 from the count of segments/points
		var a_max = (IM_PI * 2.0) * (float(num_segments) - 1.0) / float(num_segments)
		this.PathArcTo(center, radius-0.5, 0.0, a_max, num_segments-1)
	}

	this.PathStroke(col, ImDrawFlags_Closed, thickness)
}

// AddNgon Guaranteed to honor 'num_segments'
func (this *ImDrawList) AddNgon(center ImVec2, radius float, col ImU32, num_segments int, thickness float /*= 1.0f*/) {
	if (col&IM_COL32_A_MASK) == 0 || num_segments <= 2 {
		return
	}

	// Because we are filling a closed shape we remove 1 from the count of segments/points
	var a_max = (IM_PI * 2.0) * ((float)(num_segments) - 1.0) / (float)(num_segments)
	this.PathArcTo(center, radius-0.5, 0.0, a_max, num_segments-1)
	this.PathStroke(col, ImDrawFlags_Closed, thickness)
}

// AddNgonFilled Guaranteed to honor 'num_segments'
func (this *ImDrawList) AddNgonFilled(center ImVec2, radius float, col ImU32, num_segments int) {
	if (col&IM_COL32_A_MASK) == 0 || num_segments <= 2 {
		return
	}

	// Because we are filling a closed shape we remove 1 from the count of segments/points
	var a_max = (IM_PI * 2.0) * ((float)(num_segments) - 1.0) / (float)(num_segments)
	this.PathArcTo(center, radius, 0.0, a_max, num_segments-1)
	this.PathFillConvex(col)
}

// AddBezierCubic Cubic Bezier takes 4 controls points
func (this *ImDrawList) AddBezierCubic(p1 *ImVec2, p2 *ImVec2, p3 ImVec2, p4 ImVec2, col ImU32, thickness float, num_segments int) {
	if (col & IM_COL32_A_MASK) == 0 {
		return
	}

	this.PathLineTo(*p1)
	this.PathBezierCubicCurveTo(p2, p3, p4, num_segments)
	this.PathStroke(col, 0, thickness)
}

// AddBezierQuadratic Quadratic Bezier takes 3 controls points
func (this *ImDrawList) AddBezierQuadratic(p1 *ImVec2, p2 *ImVec2, p3 ImVec2, col ImU32, thickness float, num_segments int) {
	if (col & IM_COL32_A_MASK) == 0 {
		return
	}

	this.PathLineTo(*p1)
	this.PathBezierQuadraticCurveTo(p2, p3, num_segments)
	this.PathStroke(col, 0, thickness)
}

// PathClear Stateful path API, add points then finish with PathFillConvex() or PathStroke()
func (this *ImDrawList) PathClear() {
	this._Path = this._Path[:0]
}
func (this *ImDrawList) PathLineTo(pos ImVec2) {
	this._Path = append(this._Path, pos)
}
func (this *ImDrawList) PathLineToMergeDuplicate(pos ImVec2) {
	if len(this._Path) == 0 || this._Path[len(this._Path)-1] == pos {
		this._Path = append(this._Path, pos)
	}
}

// PathFillConvex Note: Anti-aliased filling requires points to be in clockwise order.
func (this *ImDrawList) PathFillConvex(col ImU32) {
	this.AddConvexPolyFilled(this._Path, int(len(this._Path)), col)
	this._Path = this._Path[:0]
}

func (this *ImDrawList) PathStroke(col ImU32, flags ImDrawFlags, thickness float /*= 1.0f*/) {
	this.AddPolyline(this._Path, int(len(this._Path)), col, flags, thickness)
	this._Path = this._Path[:0]
}

// PathBezierCubicCurveTo Cubic Bezier (4 control points)
func (this *ImDrawList) PathBezierCubicCurveTo(p2 *ImVec2, p3 ImVec2, p4 ImVec2, num_segments int) {
	var p1 = this._Path[len(this._Path)-1]
	if num_segments == 0 {
		PathBezierCubicCurveToCasteljau(&this._Path, p1.x, p1.y, p2.x, p2.y, p3.x, p3.y, p4.x, p4.y, this._Data.CurveTessellationTol, 0) // Auto-tessellated
	} else {
		var t_step = 1.0 / (float)(num_segments)
		for i_step := int(1); i_step <= num_segments; i_step++ {
			this._Path = append(this._Path, ImBezierCubicCalc(&p1, p2, &p3, &p4, t_step*float(i_step)))
		}
	}
}

// PathBezierQuadraticCurveTo Quadratic Bezier (3 control points)
func (this *ImDrawList) PathBezierQuadraticCurveTo(p2 *ImVec2, p3 ImVec2, num_segments int) {
	var p1 = this._Path[len(this._Path)-1]
	if num_segments == 0 {
		PathBezierQuadraticCurveToCasteljau(&this._Path, p1.x, p1.y, p2.x, p2.y, p3.x, p3.y, this._Data.CurveTessellationTol, 0) // Auto-tessellated
	} else {
		var t_step = 1.0 / (float)(num_segments)
		for i_step := int(1); i_step <= num_segments; i_step++ {
			this._Path = append(this._Path, ImBezierQuadraticCalc(&p1, p2, &p3, t_step*float(i_step)))
		}
	}
}

func FixRectCornerFlags(flags ImDrawFlags) ImDrawFlags {
	// If this triggers, please update your code replacing hardcoded values with new ImDrawFlags_RoundCorners* values.
	// Note that ImDrawFlags_Closed (== 0x01) is an invalid flag for AddRect(), AddRectFilled(), PathRect() etc...
	IM_ASSERT_USER_ERROR((flags&0x0F) == 0, "Misuse of legacy hardcoded ImDrawCornerFlags values!")

	if (flags & ImDrawFlags_RoundCornersMask_) == 0 {
		flags |= ImDrawFlags_RoundCornersAll
	}

	return flags
}

func (this *ImDrawList) PathRect(a, b *ImVec2, rounding float, flags ImDrawFlags) {
	flags = FixRectCornerFlags(flags)

	var xamount, yamount float = 1, 1
	if ((flags & ImDrawFlags_RoundCornersTop) == ImDrawFlags_RoundCornersTop) || ((flags & ImDrawFlags_RoundCornersBottom) == ImDrawFlags_RoundCornersBottom) {
		xamount = 0.5
	}
	if ((flags & ImDrawFlags_RoundCornersLeft) == ImDrawFlags_RoundCornersLeft) || ((flags & ImDrawFlags_RoundCornersRight) == ImDrawFlags_RoundCornersRight) {
		yamount = 0.5
	}

	rounding = ImMin(rounding, ImFabs(b.x-a.x)*(xamount)-1.0)
	rounding = ImMin(rounding, ImFabs(b.y-a.y)*(yamount)-1.0)

	if rounding <= 0.0 || (flags&ImDrawFlags_RoundCornersMask_) == ImDrawFlags_RoundCornersNone {
		this.PathLineTo(*a)
		this.PathLineTo(ImVec2{b.x, a.y})
		this.PathLineTo(*b)
		this.PathLineTo(ImVec2{a.x, b.y})
	} else {
		var rounding_tl, rounding_tr, rounding_br, rounding_bl float
		if (flags & ImDrawFlags_RoundCornersTopLeft) != 0 {
			rounding_tl = rounding
		}
		if (flags & ImDrawFlags_RoundCornersTopRight) != 0 {
			rounding_tr = rounding
		}
		if (flags & ImDrawFlags_RoundCornersBottomRight) != 0 {
			rounding_br = rounding
		}
		if (flags & ImDrawFlags_RoundCornersBottomLeft) != 0 {
			rounding_bl = rounding
		}
		this.PathArcToFast(ImVec2{a.x + rounding_tl, a.y + rounding_tl}, rounding_tl, 6, 9)
		this.PathArcToFast(ImVec2{b.x - rounding_tr, a.y + rounding_tr}, rounding_tr, 9, 12)
		this.PathArcToFast(ImVec2{b.x - rounding_br, b.y - rounding_br}, rounding_br, 0, 3)
		this.PathArcToFast(ImVec2{a.x + rounding_bl, b.y - rounding_bl}, rounding_bl, 3, 6)
	}
}

// AddCallback Advanced
// Your rendering function must check for 'UserCallback' in ImDrawCmd and call the function instead of rendering triangles.
func (this *ImDrawList) AddCallback(callback ImDrawCallback, callback_data any) {
	var curr_cmd = &this.CmdBuffer[len(this.CmdBuffer)-1]
	IM_ASSERT(curr_cmd.UserCallback == nil)
	if curr_cmd.ElemCount != 0 {
		this.AddDrawCmd()
		curr_cmd = &this.CmdBuffer[len(this.CmdBuffer)-1]
	}
	curr_cmd.UserCallback = callback
	curr_cmd.UserCallbackData = callback_data

	this.AddDrawCmd() // Force a new command after us (see comment below)
}

// AddDrawCmd This is useful if you need to forcefully create a new draw call (to allow for dependent rendering / blending). Otherwise primitives are merged into the same draw-call as much as possible
func (this *ImDrawList) AddDrawCmd() {
	var draw_cmd ImDrawCmd
	draw_cmd.ClipRect = this._CmdHeader.ClipRect // Same as calling ImDrawCmd_HeaderCopy()
	draw_cmd.TextureId = this._CmdHeader.TextureId

	draw_cmd.VtxOffset = this._CmdHeader.VtxOffset
	draw_cmd.IdxOffset = uint(len(this.IdxBuffer))

	IM_ASSERT(draw_cmd.ClipRect.x <= draw_cmd.ClipRect.z && draw_cmd.ClipRect.y <= draw_cmd.ClipRect.w)
	this.CmdBuffer = append(this.CmdBuffer, draw_cmd)
}

// CloneOutput Create a clone of the CmdBuffer/IdxBuffer/VtxBuffer.
func (this *ImDrawList) CloneOutput() *ImDrawList {
	var dst = NewImDrawList(this._Data)
	dst.CmdBuffer = this.CmdBuffer
	dst.IdxBuffer = this.IdxBuffer
	dst.VtxBuffer = this.VtxBuffer
	dst.Flags = this.Flags
	return &dst
}

// ChannelsSplit Advanced: Channels
//   - Use to split render into layers. By switching channels to can render out-of-order (e.g. submit FG primitives before BG primitives)
//   - Use to minimize draw calls (e.g. if going back-and-forth between multiple clipping rectangles, prefer to append into separate channels then merge at the end)
//   - FIXME-OBSOLETE: This API shouldn't have been in ImDrawList in the first place!
//     Prefer using your own persistent instance of ImDrawListSplitter as you can stack them.
//     Using the ImDrawList::ChannelsXXXX you cannot stack a split over another.
func (this *ImDrawList) ChannelsSplit(count int)  { this._Splitter.Split(this, count) }
func (this *ImDrawList) ChannelsMerge()           { this._Splitter.Merge(this) }
func (this *ImDrawList) ChannelsSetCurrent(n int) { this._Splitter.SetCurrentChannel(this, n) }

// PrimUnreserve Release the a number of reserved vertices/indices from the end of the last reservation made with PrimReserve().
func (this *ImDrawList) PrimUnreserve(idx_count, vtx_count int) {
	IM_ASSERT(idx_count >= 0 && vtx_count >= 0)

	var draw_cmd = &this.CmdBuffer[len(this.CmdBuffer)-1]
	draw_cmd.ElemCount -= uint(idx_count)
	this.VtxBuffer = this.VtxBuffer[:int(len(this.VtxBuffer))-vtx_count]
	this.IdxBuffer = this.IdxBuffer[:int(len(this.IdxBuffer))-idx_count]
}

// Axis aligned rectangle (composed of two triangles)

func (this *ImDrawList) PrimQuadUV(a, b, c, d *ImVec2, uv_a, uv_b, uv_c, uv_d *ImVec2, col ImU32) {
	var idx = (ImDrawIdx)(this._VtxCurrentIdx)
	this.IdxBuffer[this._IdxWritePtr+0] = idx
	this.IdxBuffer[this._IdxWritePtr+1] = (ImDrawIdx)(idx + 1)
	this.IdxBuffer[this._IdxWritePtr+2] = (ImDrawIdx)(idx + 2)
	this.IdxBuffer[this._IdxWritePtr+3] = idx
	this.IdxBuffer[this._IdxWritePtr+4] = (ImDrawIdx)(idx + 2)
	this.IdxBuffer[this._IdxWritePtr+5] = (ImDrawIdx)(idx + 3)
	this.VtxBuffer[this._VtxWritePtr+0].Pos = *a
	this.VtxBuffer[this._VtxWritePtr+0].Uv = *uv_a
	this.VtxBuffer[this._VtxWritePtr+0].Col = col
	this.VtxBuffer[this._VtxWritePtr+1].Pos = *b
	this.VtxBuffer[this._VtxWritePtr+1].Uv = *uv_b
	this.VtxBuffer[this._VtxWritePtr+1].Col = col
	this.VtxBuffer[this._VtxWritePtr+2].Pos = *c
	this.VtxBuffer[this._VtxWritePtr+2].Uv = *uv_c
	this.VtxBuffer[this._VtxWritePtr+2].Col = col
	this.VtxBuffer[this._VtxWritePtr+3].Pos = *d
	this.VtxBuffer[this._VtxWritePtr+3].Uv = *uv_d
	this.VtxBuffer[this._VtxWritePtr+3].Col = col
	this._VtxWritePtr += 4
	this._VtxCurrentIdx += 4
	this._IdxWritePtr += 6
}

func (this *ImDrawList) PrimWriteVtx(pos ImVec2, uv *ImVec2, col ImU32) {
	this.VtxBuffer[this._VtxWritePtr].Pos = pos
	this.VtxBuffer[this._VtxWritePtr].Uv = *uv
	this.VtxBuffer[this._VtxWritePtr].Col = col
	this._VtxWritePtr++

	this._VtxCurrentIdx++
}
func (this *ImDrawList) PrimWriteIdx(idx ImDrawIdx) {
	this.IdxBuffer[this._IdxWritePtr] = idx
	this._IdxWritePtr++
}

func (this *ImDrawList) PrimVtx(pos ImVec2, uv *ImVec2, col ImU32) {
	this.PrimWriteIdx((ImDrawIdx)(this._VtxCurrentIdx))
	this.PrimWriteVtx(pos, uv, col)
} // Write vertex with unique index

func (this *ImDrawList) _ClearFreeMemory() {
	this.CmdBuffer = nil
	this.IdxBuffer = nil
	this.VtxBuffer = nil
	this.Flags = ImDrawListFlags_None
	this._VtxCurrentIdx = 0
	this._VtxWritePtr = 0
	this._IdxWritePtr = 0
	this._ClipRectStack = nil
	this._TextureIdStack = nil
	this._Path = nil
	this._Splitter.ClearFreeMemory()
}

// Pop trailing draw command (used before merging or presenting to user)
// Note that this leaves the ImDrawList in a state unfit for further commands, as most code assume that CmdBuffer.Size > 0 && CmdBuffer.back().UserCallback == nil
func (this *ImDrawList) _PopUnusedDrawCmd() {
	if len(this.CmdBuffer) == 0 {
		return
	}
	var curr_cmd = &this.CmdBuffer[len(this.CmdBuffer)-1]
	if curr_cmd.ElemCount == 0 && curr_cmd.UserCallback == nil {
		this.CmdBuffer = this.CmdBuffer[:len(this.CmdBuffer)-1]
	}
}

func (this *ImDrawList) _TryMergeDrawCmds() {
	var curr_cmd = &this.CmdBuffer[len(this.CmdBuffer)-1]
	var prev_cmd = &this.CmdBuffer[len(this.CmdBuffer)-2]
	if curr_cmd.HeaderEquals(prev_cmd) && curr_cmd.UserCallback == nil && prev_cmd.UserCallback == nil {
		prev_cmd.ElemCount += curr_cmd.ElemCount
		this.CmdBuffer = this.CmdBuffer[:len(this.CmdBuffer)-1]
	}
}

func (this *ImDrawList) _OnChangedVtxOffset() {
	// We don't need to compare curr_cmd.VtxOffset != _CmdHeader.VtxOffset because we know it'll be different at the time we call this.
	this._VtxCurrentIdx = 0
	var curr_cmd = &this.CmdBuffer[len(this.CmdBuffer)-1]
	//IM_ASSERT(curr_cmd.VtxOffset != _CmdHeader.VtxOffset); // See #3349
	if curr_cmd.ElemCount != 0 {
		this.AddDrawCmd()
		return
	}
	IM_ASSERT(curr_cmd.UserCallback == nil)
	curr_cmd.VtxOffset = this._CmdHeader.VtxOffset
}

func (this *ImDrawList) _CalcCircleAutoSegmentCount(radius float) int {
	// Automatic segment count
	var radius_idx = (int)(radius + 0.999999) // ceil to never reduce accuracy
	if radius_idx < int(len(this._Data.CircleSegmentCounts)) {
		return int(this._Data.CircleSegmentCounts[radius_idx]) // Use cached value
	} else {
		return int(IM_DRAWLIST_CIRCLE_AUTO_SEGMENT_CALC(radius, this._Data.CircleSegmentMaxError))
	}
}

func (this *ImDrawList) _PathArcToN(center ImVec2, radius, a_min, a_max float, num_segments int) {
	if radius <= 0.0 {
		this._Path = append(this._Path, center)
		return
	}

	// Note that we are adding a point at both a_min and a_max.
	// If you are trying to draw a full closed circle you don't want the overlapping points!
	this._Path = reserveVec2Slice(this._Path, int(len(this._Path))+(num_segments+1))
	for i := int(0); i <= num_segments; i++ {
		var a = a_min + ((float)(i)/(float)(num_segments))*(a_max-a_min)
		this._Path = append(this._Path, ImVec2{center.x + ImCos(a)*radius, center.y + ImSin(a)*radius})
	}
}

func (this *ImDrawList) AddImage(user_texture_id ImTextureID, p_min ImVec2, p_max ImVec2, uv_min *ImVec2, uv_max *ImVec2, col ImU32) {
	if (col & IM_COL32_A_MASK) == 0 {
		return
	}

	var push_texture_id = user_texture_id != this._CmdHeader.TextureId
	if push_texture_id {
		this.PushTextureID(user_texture_id)
	}

	this.PrimReserve(6, 4)
	this.PrimRectUV(&p_min, &p_max, uv_min, uv_max, col)

	if push_texture_id {
		this.PopTextureID()
	}
}

func (this *ImDrawList) PopTextureID() {
	this._TextureIdStack = this._TextureIdStack[:len(this._TextureIdStack)-1]
	if len(this._TextureIdStack) > 0 {
		this._CmdHeader.TextureId = this._TextureIdStack[len(this._TextureIdStack)-1]
	} else {
		this._CmdHeader.TextureId = 0
	}
	this._OnChangedTextureID()
}

func (this *ImDrawList) PrimRectUV(a, c, uv_a, uv_c *ImVec2, col ImU32) {
	var b, d, uv_b, uv_d = ImVec2{c.x, a.y}, ImVec2{a.x, c.y}, ImVec2{uv_c.x, uv_a.y}, ImVec2{uv_a.x, uv_c.y}
	var idx = (ImDrawIdx)(this._VtxCurrentIdx)
	this.IdxBuffer[this._IdxWritePtr+0] = idx
	this.IdxBuffer[this._IdxWritePtr+1] = (ImDrawIdx)(idx + 1)
	this.IdxBuffer[this._IdxWritePtr+2] = (ImDrawIdx)(idx + 2)
	this.IdxBuffer[this._IdxWritePtr+3] = idx
	this.IdxBuffer[this._IdxWritePtr+4] = (ImDrawIdx)(idx + 2)
	this.IdxBuffer[this._IdxWritePtr+5] = (ImDrawIdx)(idx + 3)
	this.VtxBuffer[this._VtxWritePtr+0].Pos = *a
	this.VtxBuffer[this._VtxWritePtr+0].Uv = *uv_a
	this.VtxBuffer[this._VtxWritePtr+0].Col = col
	this.VtxBuffer[this._VtxWritePtr+1].Pos = b
	this.VtxBuffer[this._VtxWritePtr+1].Uv = uv_b
	this.VtxBuffer[this._VtxWritePtr+1].Col = col
	this.VtxBuffer[this._VtxWritePtr+2].Pos = *c
	this.VtxBuffer[this._VtxWritePtr+2].Uv = *uv_c
	this.VtxBuffer[this._VtxWritePtr+2].Col = col
	this.VtxBuffer[this._VtxWritePtr+3].Pos = d
	this.VtxBuffer[this._VtxWritePtr+3].Uv = uv_d
	this.VtxBuffer[this._VtxWritePtr+3].Col = col
	this._VtxWritePtr += 4
	this._VtxCurrentIdx += 4
	this._IdxWritePtr += 6
}

func (this *ImDrawList) AddImageQuad(user_texture_id ImTextureID, p1 *ImVec2, p2 *ImVec2, p3 ImVec2, p4 ImVec2, uv1 *ImVec2, uv2 *ImVec2 /*= ImVec2(1, 0)*/, uv3 ImVec2 /*ImVec2(1, 1)*/, uv4 ImVec2 /*ImVec2(0, 1)*/, col ImU32) {
	if (col & IM_COL32_A_MASK) == 0 {
		return
	}

	var push_texture_id = user_texture_id != this._CmdHeader.TextureId
	if push_texture_id {
		this.PushTextureID(user_texture_id)
	}

	this.PrimReserve(6, 4)
	this.PrimQuadUV(p1, p2, &p3, &p4, uv1, uv2, &uv3, &uv4, col)

	if push_texture_id {
		this.PopTextureID()
	}
}

func (this *ImDrawList) AddImageRounded(user_texture_id ImTextureID, p_min ImVec2, p_max ImVec2, uv_min, uv_max *ImVec2, col ImU32, rounding float, flags ImDrawFlags) {
	if (col & IM_COL32_A_MASK) == 0 {
		return
	}

	flags = FixRectCornerFlags(flags)
	if rounding <= 0.0 || (flags&ImDrawFlags_RoundCornersMask_) == ImDrawFlags_RoundCornersNone {
		this.AddImage(user_texture_id, p_min, p_max, uv_min, uv_max, col)
		return
	}

	var push_texture_id = user_texture_id != this._CmdHeader.TextureId
	if push_texture_id {
		this.PushTextureID(user_texture_id)
	}

	var vert_start_idx = int(len(this.VtxBuffer))
	this.PathRect(&p_min, &p_max, rounding, flags)
	this.PathFillConvex(col)
	var vert_end_idx = int(len(this.VtxBuffer))
	ShadeVertsLinearUV(this, vert_start_idx, vert_end_idx, &p_min, &p_max, uv_min, uv_max, true)

	if push_texture_id {
		this.PopTextureID()
	}
}
