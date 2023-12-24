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

func (l *ImDrawList) PushClipRect(cr_min, cr_max ImVec2, intersect_with_current_clip_rect bool) {
	var cr = ImVec4{cr_min.x, cr_min.y, cr_max.x, cr_max.y}
	if intersect_with_current_clip_rect {
		current := l._CmdHeader.ClipRect
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
	cr.z = max(cr.x, cr.z)
	cr.w = max(cr.y, cr.w)

	l._ClipRectStack = append(l._ClipRectStack, cr)
	l._CmdHeader.ClipRect = cr
	l._OnChangedClipRect()
}

// PushClipRectFullScreen Render-level scissoring. This is passed down to your render function but not used for CPU-side coarse clipping. Prefer using higher-level ImGui::PushClipRect() to affect logic (hit-testing and widget culling)
func (l *ImDrawList) PushClipRectFullScreen() {
	l.PushClipRect(ImVec2{l._Data.ClipRectFullscreen.x, l._Data.ClipRectFullscreen.y}, ImVec2{l._Data.ClipRectFullscreen.z, l._Data.ClipRectFullscreen.w}, false)
}

func (l *ImDrawList) PopClipRect() {
	l._ClipRectStack = l._ClipRectStack[len(l._ClipRectStack)-1:]
	if len(l._ClipRectStack) == 0 {
		l._CmdHeader.ClipRect = l._Data.ClipRectFullscreen
	} else {
		l._CmdHeader.ClipRect = l._ClipRectStack[len(l._ClipRectStack)-1]
	}
	l._OnChangedClipRect()
}

func (l *ImDrawList) GetClipRectMin() ImVec2 {
	var cr = &l._ClipRectStack[len(l._ClipRectStack)-1]
	return ImVec2{cr.x, cr.y}
}
func (l *ImDrawList) GetClipRectMax() ImVec2 {
	var cr = &l._ClipRectStack[len(l._ClipRectStack)-1]
	return ImVec2{cr.x, cr.y}
}

// AddLine Primitives
//   - For rectangular primitives, "p_min" and "p_max" represent the upper-left and lower-right corners.
//   - For circle primitives, use "num_segments == 0" to automatically calculate tessellation (preferred).
//     In older versions (until Dear ImGui 1.77) the AddCircle functions defaulted to num_segments == 12.
//     In future versions we will use textures to provide cheaper and higher-quality circles.
//     Use AddNgon() and AddNgonFilled() functions if you need to guaranteed a specific number of sides.
func (l *ImDrawList) AddLine(p1 *ImVec2, p2 *ImVec2, col ImU32, thickness float /*= 1.0f*/) {
	if (col & IM_COL32_A_MASK) == 0 {
		return
	}
	l.PathLineTo(p1.Add(ImVec2{0.5, 0.5}))
	l.PathLineTo(p2.Add(ImVec2{0.5, 0.5}))
	l.PathStroke(col, 0, thickness)
}

// AddRectFilledMultiColor p_min = upper-left, p_max = lower-right
func (l *ImDrawList) AddRectFilledMultiColor(p_min ImVec2, p_max ImVec2, col_upr_left, col_upr_right, col_bot_right, col_bot_left ImU32) {
	if ((col_upr_left | col_upr_right | col_bot_right | col_bot_left) & IM_COL32_A_MASK) == 0 {
		return
	}

	var uv = l._Data.TexUvWhitePixel
	l.PrimReserve(6, 4)
	l.PrimWriteIdx((ImDrawIdx)(l._VtxCurrentIdx))
	l.PrimWriteIdx((ImDrawIdx)(l._VtxCurrentIdx + 1))
	l.PrimWriteIdx((ImDrawIdx)(l._VtxCurrentIdx + 2))
	l.PrimWriteIdx((ImDrawIdx)(l._VtxCurrentIdx))
	l.PrimWriteIdx((ImDrawIdx)(l._VtxCurrentIdx + 2))
	l.PrimWriteIdx((ImDrawIdx)(l._VtxCurrentIdx + 3))
	l.PrimWriteVtx(p_min, &uv, col_upr_left)
	l.PrimWriteVtx(ImVec2{p_max.x, p_min.y}, &uv, col_upr_right)
	l.PrimWriteVtx(p_max, &uv, col_bot_right)
	l.PrimWriteVtx(ImVec2{p_min.x, p_max.y}, &uv, col_bot_left)
}

func (l *ImDrawList) AddQuad(p1 *ImVec2, p2 *ImVec2, p3 ImVec2, p4 ImVec2, col ImU32, thickness float /*= 1.0f*/) {
	if (col & IM_COL32_A_MASK) == 0 {
		return
	}

	l.PathLineTo(*p1)
	l.PathLineTo(*p2)
	l.PathLineTo(p3)
	l.PathLineTo(p4)
	l.PathStroke(col, ImDrawFlags_Closed, thickness)
}

func (l *ImDrawList) AddQuadFilled(p1 *ImVec2, p2 *ImVec2, p3 ImVec2, p4 ImVec2, col ImU32) {
	if (col & IM_COL32_A_MASK) == 0 {
		return
	}

	l.PathLineTo(*p1)
	l.PathLineTo(*p2)
	l.PathLineTo(p3)
	l.PathLineTo(p4)
	l.PathFillConvex(col)
}

func (l *ImDrawList) AddTriangle(p1 *ImVec2, p2 *ImVec2, p3 ImVec2, col ImU32, thickness float /*= 1.0f*/) {
	if (col & IM_COL32_A_MASK) == 0 {
		return
	}

	l.PathLineTo(*p1)
	l.PathLineTo(*p2)
	l.PathLineTo(p3)
	l.PathStroke(col, ImDrawFlags_Closed, thickness)
}

func (l *ImDrawList) AddTriangleFilled(p1 *ImVec2, p2 *ImVec2, p3 ImVec2, col ImU32) {
	if (col & IM_COL32_A_MASK) == 0 {
		return
	}

	l.PathLineTo(*p1)
	l.PathLineTo(*p2)
	l.PathLineTo(p3)
	l.PathFillConvex(col)
}

func (l *ImDrawList) AddCircle(center ImVec2, radius float, col ImU32, num_segments int, thickness float /*= 1.0f*/) {
	if col&IM_COL32_A_MASK == 0 || radius < 0.0 {
		return
	}

	if num_segments <= 0 {
		// Use arc with automatic segment count
		l.PathArcToFastEx(center, radius-0.5, 0, IM_DRAWLIST_ARCFAST_SAMPLE_MAX, 0)
		l._Path = l._Path[:len(l._Path)-1]
	} else {
		// Explicit segment count (still clamp to avoid drawing insanely tessellated shapes)
		num_segments = ImClampInt(num_segments, 3, IM_DRAWLIST_CIRCLE_AUTO_SEGMENT_MAX)

		// Because we are filling a closed shape we remove 1 from the count of segments/points
		var a_max = (IM_PI * 2.0) * (float(num_segments) - 1.0) / float(num_segments)
		l.PathArcTo(center, radius-0.5, 0.0, a_max, num_segments-1)
	}

	l.PathStroke(col, ImDrawFlags_Closed, thickness)
}

// AddNgon Guaranteed to honor 'num_segments'
func (l *ImDrawList) AddNgon(center ImVec2, radius float, col ImU32, num_segments int, thickness float /*= 1.0f*/) {
	if (col&IM_COL32_A_MASK) == 0 || num_segments <= 2 {
		return
	}

	// Because we are filling a closed shape we remove 1 from the count of segments/points
	a_max := (IM_PI * 2.0) * ((float)(num_segments) - 1.0) / (float)(num_segments)
	l.PathArcTo(center, radius-0.5, 0.0, a_max, num_segments-1)
	l.PathStroke(col, ImDrawFlags_Closed, thickness)
}

// AddNgonFilled Guaranteed to honor 'num_segments'
func (l *ImDrawList) AddNgonFilled(center ImVec2, radius float, col ImU32, num_segments int) {
	if (col&IM_COL32_A_MASK) == 0 || num_segments <= 2 {
		return
	}

	// Because we are filling a closed shape we remove 1 from the count of segments/points
	var a_max = (IM_PI * 2.0) * ((float)(num_segments) - 1.0) / (float)(num_segments)
	l.PathArcTo(center, radius, 0.0, a_max, num_segments-1)
	l.PathFillConvex(col)
}

// AddBezierCubic Cubic Bezier takes 4 controls points
func (l *ImDrawList) AddBezierCubic(p1 *ImVec2, p2 *ImVec2, p3 ImVec2, p4 ImVec2, col ImU32, thickness float, num_segments int) {
	if (col & IM_COL32_A_MASK) == 0 {
		return
	}

	l.PathLineTo(*p1)
	l.PathBezierCubicCurveTo(p2, p3, p4, num_segments)
	l.PathStroke(col, 0, thickness)
}

// AddBezierQuadratic Quadratic Bezier takes 3 controls points
func (l *ImDrawList) AddBezierQuadratic(p1 *ImVec2, p2 *ImVec2, p3 ImVec2, col ImU32, thickness float, num_segments int) {
	if (col & IM_COL32_A_MASK) == 0 {
		return
	}

	l.PathLineTo(*p1)
	l.PathBezierQuadraticCurveTo(p2, p3, num_segments)
	l.PathStroke(col, 0, thickness)
}

// PathClear Stateful path API, add points then finish with PathFillConvex() or PathStroke()
func (l *ImDrawList) PathClear() {
	l._Path = l._Path[:0]
}
func (l *ImDrawList) PathLineTo(pos ImVec2) {
	l._Path = append(l._Path, pos)
}
func (l *ImDrawList) PathLineToMergeDuplicate(pos ImVec2) {
	if len(l._Path) == 0 || l._Path[len(l._Path)-1] == pos {
		l._Path = append(l._Path, pos)
	}
}

// PathFillConvex Note: Anti-aliased filling requires points to be in clockwise order.
func (l *ImDrawList) PathFillConvex(col ImU32) {
	l.AddConvexPolyFilled(l._Path, int(len(l._Path)), col)
	l._Path = l._Path[:0]
}

func (l *ImDrawList) PathStroke(col ImU32, flags ImDrawFlags, thickness float /*= 1.0f*/) {
	l.AddPolyline(l._Path, int(len(l._Path)), col, flags, thickness)
	l._Path = l._Path[:0]
}

// PathBezierCubicCurveTo Cubic Bezier (4 control points)
func (l *ImDrawList) PathBezierCubicCurveTo(p2 *ImVec2, p3 ImVec2, p4 ImVec2, num_segments int) {
	p1 := l._Path[len(l._Path)-1]
	if num_segments == 0 {
		PathBezierCubicCurveToCasteljau(&l._Path, p1.x, p1.y, p2.x, p2.y, p3.x, p3.y, p4.x, p4.y, l._Data.CurveTessellationTol, 0) // Auto-tessellated
	} else {
		var t_step = 1.0 / (float)(num_segments)
		for i_step := int(1); i_step <= num_segments; i_step++ {
			l._Path = append(l._Path, ImBezierCubicCalc(&p1, p2, &p3, &p4, t_step*float(i_step)))
		}
	}
}

// PathBezierQuadraticCurveTo Quadratic Bezier (3 control points)
func (l *ImDrawList) PathBezierQuadraticCurveTo(p2 *ImVec2, p3 ImVec2, num_segments int) {
	p1 := l._Path[len(l._Path)-1]
	if num_segments == 0 {
		PathBezierQuadraticCurveToCasteljau(&l._Path, p1.x, p1.y, p2.x, p2.y, p3.x, p3.y, l._Data.CurveTessellationTol, 0) // Auto-tessellated
	} else {
		var t_step = 1.0 / (float)(num_segments)
		for i_step := int(1); i_step <= num_segments; i_step++ {
			l._Path = append(l._Path, ImBezierQuadraticCalc(&p1, p2, &p3, t_step*float(i_step)))
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

func (l *ImDrawList) PathRect(a, b *ImVec2, rounding float, flags ImDrawFlags) {
	flags = FixRectCornerFlags(flags)

	var xamount, yamount float = 1, 1
	if ((flags & ImDrawFlags_RoundCornersTop) == ImDrawFlags_RoundCornersTop) || ((flags & ImDrawFlags_RoundCornersBottom) == ImDrawFlags_RoundCornersBottom) {
		xamount = 0.5
	}
	if ((flags & ImDrawFlags_RoundCornersLeft) == ImDrawFlags_RoundCornersLeft) || ((flags & ImDrawFlags_RoundCornersRight) == ImDrawFlags_RoundCornersRight) {
		yamount = 0.5
	}

	rounding = min(rounding, ImFabs(b.x-a.x)*(xamount)-1.0)
	rounding = min(rounding, ImFabs(b.y-a.y)*(yamount)-1.0)

	if rounding <= 0.0 || (flags&ImDrawFlags_RoundCornersMask_) == ImDrawFlags_RoundCornersNone {
		l.PathLineTo(*a)
		l.PathLineTo(ImVec2{b.x, a.y})
		l.PathLineTo(*b)
		l.PathLineTo(ImVec2{a.x, b.y})
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
		l.PathArcToFast(ImVec2{a.x + rounding_tl, a.y + rounding_tl}, rounding_tl, 6, 9)
		l.PathArcToFast(ImVec2{b.x - rounding_tr, a.y + rounding_tr}, rounding_tr, 9, 12)
		l.PathArcToFast(ImVec2{b.x - rounding_br, b.y - rounding_br}, rounding_br, 0, 3)
		l.PathArcToFast(ImVec2{a.x + rounding_bl, b.y - rounding_bl}, rounding_bl, 3, 6)
	}
}

// AddCallback Advanced
// Your rendering function must check for 'UserCallback' in ImDrawCmd and call the function instead of rendering triangles.
func (l *ImDrawList) AddCallback(callback ImDrawCallback, callback_data any) {
	curr_cmd := &l.CmdBuffer[len(l.CmdBuffer)-1]
	IM_ASSERT(curr_cmd.UserCallback == nil)
	if curr_cmd.ElemCount != 0 {
		l.AddDrawCmd()
		curr_cmd = &l.CmdBuffer[len(l.CmdBuffer)-1]
	}
	curr_cmd.UserCallback = callback
	curr_cmd.UserCallbackData = callback_data

	l.AddDrawCmd() // Force a new command after us (see comment below)
}

// AddDrawCmd This is useful if you need to forcefully create a new draw call (to allow for dependent rendering / blending). Otherwise primitives are merged into the same draw-call as much as possible
func (l *ImDrawList) AddDrawCmd() {
	var draw_cmd ImDrawCmd
	draw_cmd.ClipRect = l._CmdHeader.ClipRect // Same as calling ImDrawCmd_HeaderCopy()
	draw_cmd.TextureId = l._CmdHeader.TextureId

	draw_cmd.VtxOffset = l._CmdHeader.VtxOffset
	draw_cmd.IdxOffset = uint(len(l.IdxBuffer))

	IM_ASSERT(draw_cmd.ClipRect.x <= draw_cmd.ClipRect.z && draw_cmd.ClipRect.y <= draw_cmd.ClipRect.w)
	l.CmdBuffer = append(l.CmdBuffer, draw_cmd)
}

// CloneOutput Create a clone of the CmdBuffer/IdxBuffer/VtxBuffer.
func (l *ImDrawList) CloneOutput() *ImDrawList {
	dst := NewImDrawList(l._Data)
	dst.CmdBuffer = l.CmdBuffer
	dst.IdxBuffer = l.IdxBuffer
	dst.VtxBuffer = l.VtxBuffer
	dst.Flags = l.Flags
	return &dst
}

// ChannelsSplit Advanced: Channels
//   - Use to split render into layers. By switching channels to can render out-of-order (e.guiContext. submit FG primitives before BG primitives)
//   - Use to minimize draw calls (e.guiContext. if going back-and-forth between multiple clipping rectangles, prefer to append into separate channels then merge at the end)
//   - FIXME-OBSOLETE: This API shouldn't have been in ImDrawList in the first place!
//     Prefer using your own persistent instance of ImDrawListSplitter as you can stack them.
//     Using the ImDrawList::ChannelsXXXX you cannot stack a split over another.
func (l *ImDrawList) ChannelsSplit(count int)  { l._Splitter.Split(l, count) }
func (l *ImDrawList) ChannelsMerge()           { l._Splitter.Merge(l) }
func (l *ImDrawList) ChannelsSetCurrent(n int) { l._Splitter.SetCurrentChannel(l, n) }

// PrimUnreserve Release the a number of reserved vertices/indices from the end of the last reservation made with PrimReserve().
func (l *ImDrawList) PrimUnreserve(idx_count, vtx_count int) {
	IM_ASSERT(idx_count >= 0 && vtx_count >= 0)

	var draw_cmd = &l.CmdBuffer[len(l.CmdBuffer)-1]
	draw_cmd.ElemCount -= uint(idx_count)
	l.VtxBuffer = l.VtxBuffer[:int(len(l.VtxBuffer))-vtx_count]
	l.IdxBuffer = l.IdxBuffer[:int(len(l.IdxBuffer))-idx_count]
}

// Axis aligned rectangle (composed of two triangles)

func (l *ImDrawList) PrimQuadUV(a, b, c, d *ImVec2, uv_a, uv_b, uv_c, uv_d *ImVec2, col ImU32) {
	var idx = (ImDrawIdx)(l._VtxCurrentIdx)
	l.IdxBuffer[l._IdxWritePtr+0] = idx
	l.IdxBuffer[l._IdxWritePtr+1] = (ImDrawIdx)(idx + 1)
	l.IdxBuffer[l._IdxWritePtr+2] = (ImDrawIdx)(idx + 2)
	l.IdxBuffer[l._IdxWritePtr+3] = idx
	l.IdxBuffer[l._IdxWritePtr+4] = (ImDrawIdx)(idx + 2)
	l.IdxBuffer[l._IdxWritePtr+5] = (ImDrawIdx)(idx + 3)
	l.VtxBuffer[l._VtxWritePtr+0].Pos = *a
	l.VtxBuffer[l._VtxWritePtr+0].Uv = *uv_a
	l.VtxBuffer[l._VtxWritePtr+0].Col = col
	l.VtxBuffer[l._VtxWritePtr+1].Pos = *b
	l.VtxBuffer[l._VtxWritePtr+1].Uv = *uv_b
	l.VtxBuffer[l._VtxWritePtr+1].Col = col
	l.VtxBuffer[l._VtxWritePtr+2].Pos = *c
	l.VtxBuffer[l._VtxWritePtr+2].Uv = *uv_c
	l.VtxBuffer[l._VtxWritePtr+2].Col = col
	l.VtxBuffer[l._VtxWritePtr+3].Pos = *d
	l.VtxBuffer[l._VtxWritePtr+3].Uv = *uv_d
	l.VtxBuffer[l._VtxWritePtr+3].Col = col
	l._VtxWritePtr += 4
	l._VtxCurrentIdx += 4
	l._IdxWritePtr += 6
}

func (l *ImDrawList) PrimWriteVtx(pos ImVec2, uv *ImVec2, col ImU32) {
	l.VtxBuffer[l._VtxWritePtr].Pos = pos
	l.VtxBuffer[l._VtxWritePtr].Uv = *uv
	l.VtxBuffer[l._VtxWritePtr].Col = col
	l._VtxWritePtr++

	l._VtxCurrentIdx++
}
func (l *ImDrawList) PrimWriteIdx(idx ImDrawIdx) {
	l.IdxBuffer[l._IdxWritePtr] = idx
	l._IdxWritePtr++
}

func (l *ImDrawList) PrimVtx(pos ImVec2, uv *ImVec2, col ImU32) {
	l.PrimWriteIdx((ImDrawIdx)(l._VtxCurrentIdx))
	l.PrimWriteVtx(pos, uv, col)
} // Write vertex with unique index

func (l *ImDrawList) _ClearFreeMemory() {
	l.CmdBuffer = nil
	l.IdxBuffer = nil
	l.VtxBuffer = nil
	l.Flags = ImDrawListFlags_None
	l._VtxCurrentIdx = 0
	l._VtxWritePtr = 0
	l._IdxWritePtr = 0
	l._ClipRectStack = nil
	l._TextureIdStack = nil
	l._Path = nil
	l._Splitter.ClearFreeMemory()
}

// Pop trailing draw command (used before merging or presenting to user)
// Note that this leaves the ImDrawList in a state unfit for further commands, as most code assume that CmdBuffer.Size > 0 && CmdBuffer.back().UserCallback == nil
func (l *ImDrawList) _PopUnusedDrawCmd() {
	if len(l.CmdBuffer) == 0 {
		return
	}
	var curr_cmd = &l.CmdBuffer[len(l.CmdBuffer)-1]
	if curr_cmd.ElemCount == 0 && curr_cmd.UserCallback == nil {
		l.CmdBuffer = l.CmdBuffer[:len(l.CmdBuffer)-1]
	}
}

func (l *ImDrawList) _TryMergeDrawCmds() {
	curr_cmd := &l.CmdBuffer[len(l.CmdBuffer)-1]
	prev_cmd := &l.CmdBuffer[len(l.CmdBuffer)-2]
	if curr_cmd.HeaderEquals(prev_cmd) && curr_cmd.UserCallback == nil && prev_cmd.UserCallback == nil {
		prev_cmd.ElemCount += curr_cmd.ElemCount
		l.CmdBuffer = l.CmdBuffer[:len(l.CmdBuffer)-1]
	}
}

func (l *ImDrawList) _OnChangedVtxOffset() {
	// We don't need to compare curr_cmd.VtxOffset != _CmdHeader.VtxOffset because we know it'll be different at the time we call l.
	l._VtxCurrentIdx = 0
	curr_cmd := &l.CmdBuffer[len(l.CmdBuffer)-1]
	//IM_ASSERT(curr_cmd.VtxOffset != _CmdHeader.VtxOffset); // See #3349
	if curr_cmd.ElemCount != 0 {
		l.AddDrawCmd()
		return
	}
	IM_ASSERT(curr_cmd.UserCallback == nil)
	curr_cmd.VtxOffset = l._CmdHeader.VtxOffset
}

func (l *ImDrawList) _CalcCircleAutoSegmentCount(radius float) int {
	// Automatic segment count
	var radius_idx = (int)(radius + 0.999999) // ceil to never reduce accuracy
	if radius_idx < int(len(l._Data.CircleSegmentCounts)) {
		return int(l._Data.CircleSegmentCounts[radius_idx]) // Use cached value
	} else {
		return int(IM_DRAWLIST_CIRCLE_AUTO_SEGMENT_CALC(radius, l._Data.CircleSegmentMaxError))
	}
}

func (l *ImDrawList) _PathArcToN(center ImVec2, radius, a_min, a_max float, num_segments int) {
	if radius <= 0.0 {
		l._Path = append(l._Path, center)
		return
	}

	// Note that we are adding a point at both a_min and a_max.
	// If you are trying to draw a full closed circle you don't want the overlapping points!
	l._Path = reserveVec2Slice(l._Path, int(len(l._Path))+(num_segments+1))
	for i := int(0); i <= num_segments; i++ {
		var a = a_min + ((float)(i)/(float)(num_segments))*(a_max-a_min)
		l._Path = append(l._Path, ImVec2{center.x + ImCos(a)*radius, center.y + ImSin(a)*radius})
	}
}

func (l *ImDrawList) AddImage(user_texture_id ImTextureID, p_min ImVec2, p_max ImVec2, uv_min *ImVec2, uv_max *ImVec2, col ImU32) {
	if (col & IM_COL32_A_MASK) == 0 {
		return
	}

	var push_texture_id = user_texture_id != l._CmdHeader.TextureId
	if push_texture_id {
		l.PushTextureID(user_texture_id)
	}

	l.PrimReserve(6, 4)
	l.PrimRectUV(&p_min, &p_max, uv_min, uv_max, col)

	if push_texture_id {
		l.PopTextureID()
	}
}

func (l *ImDrawList) PopTextureID() {
	l._TextureIdStack = l._TextureIdStack[:len(l._TextureIdStack)-1]
	if len(l._TextureIdStack) > 0 {
		l._CmdHeader.TextureId = l._TextureIdStack[len(l._TextureIdStack)-1]
	} else {
		l._CmdHeader.TextureId = 0
	}
	l._OnChangedTextureID()
}

func (l *ImDrawList) PrimRectUV(a, c, uv_a, uv_c *ImVec2, col ImU32) {
	b, d, uv_b, uv_d := ImVec2{c.x, a.y}, ImVec2{a.x, c.y}, ImVec2{uv_c.x, uv_a.y}, ImVec2{uv_a.x, uv_c.y}
	idx := (ImDrawIdx)(l._VtxCurrentIdx)
	l.IdxBuffer[l._IdxWritePtr+0] = idx
	l.IdxBuffer[l._IdxWritePtr+1] = (ImDrawIdx)(idx + 1)
	l.IdxBuffer[l._IdxWritePtr+2] = (ImDrawIdx)(idx + 2)
	l.IdxBuffer[l._IdxWritePtr+3] = idx
	l.IdxBuffer[l._IdxWritePtr+4] = (ImDrawIdx)(idx + 2)
	l.IdxBuffer[l._IdxWritePtr+5] = (ImDrawIdx)(idx + 3)
	l.VtxBuffer[l._VtxWritePtr+0].Pos = *a
	l.VtxBuffer[l._VtxWritePtr+0].Uv = *uv_a
	l.VtxBuffer[l._VtxWritePtr+0].Col = col
	l.VtxBuffer[l._VtxWritePtr+1].Pos = b
	l.VtxBuffer[l._VtxWritePtr+1].Uv = uv_b
	l.VtxBuffer[l._VtxWritePtr+1].Col = col
	l.VtxBuffer[l._VtxWritePtr+2].Pos = *c
	l.VtxBuffer[l._VtxWritePtr+2].Uv = *uv_c
	l.VtxBuffer[l._VtxWritePtr+2].Col = col
	l.VtxBuffer[l._VtxWritePtr+3].Pos = d
	l.VtxBuffer[l._VtxWritePtr+3].Uv = uv_d
	l.VtxBuffer[l._VtxWritePtr+3].Col = col
	l._VtxWritePtr += 4
	l._VtxCurrentIdx += 4
	l._IdxWritePtr += 6
}

func (l *ImDrawList) AddImageQuad(user_texture_id ImTextureID, p1 *ImVec2, p2 *ImVec2, p3 ImVec2, p4 ImVec2, uv1 *ImVec2, uv2 *ImVec2 /*= ImVec2(1, 0)*/, uv3 ImVec2 /*ImVec2(1, 1)*/, uv4 ImVec2 /*ImVec2(0, 1)*/, col ImU32) {
	if (col & IM_COL32_A_MASK) == 0 {
		return
	}

	push_texture_id := user_texture_id != l._CmdHeader.TextureId
	if push_texture_id {
		l.PushTextureID(user_texture_id)
	}

	l.PrimReserve(6, 4)
	l.PrimQuadUV(p1, p2, &p3, &p4, uv1, uv2, &uv3, &uv4, col)

	if push_texture_id {
		l.PopTextureID()
	}
}

func (l *ImDrawList) AddImageRounded(user_texture_id ImTextureID, p_min ImVec2, p_max ImVec2, uv_min, uv_max *ImVec2, col ImU32, rounding float, flags ImDrawFlags) {
	if (col & IM_COL32_A_MASK) == 0 {
		return
	}

	flags = FixRectCornerFlags(flags)
	if rounding <= 0.0 || (flags&ImDrawFlags_RoundCornersMask_) == ImDrawFlags_RoundCornersNone {
		l.AddImage(user_texture_id, p_min, p_max, uv_min, uv_max, col)
		return
	}

	push_texture_id := user_texture_id != l._CmdHeader.TextureId
	if push_texture_id {
		l.PushTextureID(user_texture_id)
	}

	vert_start_idx := int(len(l.VtxBuffer))
	l.PathRect(&p_min, &p_max, rounding, flags)
	l.PathFillConvex(col)
	var vert_end_idx = int(len(l.VtxBuffer))
	ShadeVertsLinearUV(l, vert_start_idx, vert_end_idx, &p_min, &p_max, uv_min, uv_max, true)

	if push_texture_id {
		l.PopTextureID()
	}
}
