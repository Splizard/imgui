package imgui

//-----------------------------------------------------------------------------
// [SECTION] ImDrawListSplitter
//-----------------------------------------------------------------------------
// FIXME: This may be a little confusing, trying to be a little too low-level/optimal instead of just doing vector swap..
//-----------------------------------------------------------------------------

// Split/Merge functions are used to split the draw list into different layers which can be drawn into out of order.
// This is used by the Columns/Tables API, so items of each column can be batched together in a same draw call.
type ImDrawListSplitter struct {
	_Current  int             // Current channel number (0)
	_Count    int             // Number of active channels (1+)
	_Channels []ImDrawChannel // Draw channels (not resized down so _Count might be < Channels.Size)
}

func (this *ImDrawListSplitter) Clear() {
	this._Current = 0
	this._Count = 1 // Do not clear Channels[] so our allocations are reused next frame
}

func (this *ImDrawListSplitter) ClearFreeMemory() {
	for i := int(0); i < int(len(this._Channels)); i++ {
		if i == this._Current {
			this._Channels[i] = ImDrawChannel{} // Current channel is a copy of CmdBuffer/IdxBuffer, don't destruct again
		}
		this._Channels[i]._CmdBuffer = nil
		this._Channels[i]._IdxBuffer = nil
	}
	this._Current = 0
	this._Count = 1
	this._Channels = nil
}

func (this *ImDrawListSplitter) Split(draw_list *ImDrawList, channels_count int) {
	IM_ASSERT_USER_ERROR(this._Current == 0 && this._Count <= 1, "Nested channel splitting is not supported. Please use separate instances of ImDrawListSplitter.")
	var old_channels_count = int(len(this._Channels))
	if old_channels_count < channels_count {
		this._Channels = append(this._Channels, make([]ImDrawChannel, channels_count-old_channels_count)...)
	}
	this._Count = channels_count

	// Channels[] (24/32 bytes each) hold storage that we'll swap with draw_list._CmdBuffer/_IdxBuffer
	// The content of Channels[0] at this point doesn't matter. We clear it to make state tidy in a debugger but we don't strictly need to.
	// When we switch to the next channel, we'll copy draw_list._CmdBuffer/_IdxBuffer into Channels[0] and then Channels[1] into draw_list.CmdBuffer/_IdxBuffer
	this._Channels[0] = ImDrawChannel{}
	for i := int(1); i < channels_count; i++ {
		if i >= old_channels_count {
			this._Channels[i] = ImDrawChannel{}
		} else {
			this._Channels[i]._CmdBuffer = this._Channels[i]._CmdBuffer[:0]
			this._Channels[i]._IdxBuffer = this._Channels[i]._IdxBuffer[:0]
		}
	}
}

func (this *ImDrawListSplitter) Merge(draw_list *ImDrawList) {
	// Note that we never use or rely on _Channels.Size because it is merely a buffer that we never shrink back to 0 to keep all sub-buffers ready for use.
	if this._Count <= 1 {
		return
	}

	this.SetCurrentChannel(draw_list, 0)
	draw_list._PopUnusedDrawCmd()

	// Calculate our final buffer sizes. Also fix the incorrect IdxOffset values in each command.
	var new_cmd_buffer_count int = 0
	var new_idx_buffer_count int = 0

	var last_cmd *ImDrawCmd
	if this._Count > 0 && len(draw_list.CmdBuffer) > 0 {
		last_cmd = &draw_list.CmdBuffer[len(draw_list.CmdBuffer)-1]
	}

	var idx_offset int
	if last_cmd != nil {
		idx_offset = int(last_cmd.IdxOffset + last_cmd.ElemCount)
	}

	for i := int(1); i < this._Count; i++ {
		var ch = &this._Channels[i]

		// Equivalent of PopUnusedDrawCmd() for this channel's cmdbuffer and except we don't need to test for UserCallback.
		if len(ch._CmdBuffer) > 0 && ch._CmdBuffer[len(ch._CmdBuffer)-1].ElemCount == 0 {
			ch._CmdBuffer = ch._CmdBuffer[:len(ch._CmdBuffer)-1]
		}

		if len(ch._CmdBuffer) > 0 && last_cmd != nil {
			var next_cmd = &ch._CmdBuffer[0]
			if last_cmd.HeaderEquals(next_cmd) && last_cmd.UserCallback == nil && next_cmd.UserCallback == nil {
				// Merge previous channel last draw command with current channel first draw command if matching.
				last_cmd.ElemCount += next_cmd.ElemCount
				idx_offset += int(next_cmd.ElemCount)
				copy(ch._CmdBuffer, ch._CmdBuffer[1:]) // FIXME-OPT: Improve for multiple merges.
			}
		}
		if len(ch._CmdBuffer) > 0 {
			last_cmd = &ch._CmdBuffer[len(ch._CmdBuffer)-1]
		}
		new_cmd_buffer_count += int(len(ch._CmdBuffer))
		new_idx_buffer_count += int(len(ch._IdxBuffer))
		for cmd_n := range ch._CmdBuffer {
			ch._CmdBuffer[cmd_n].IdxOffset = uint(idx_offset)
			idx_offset += int(ch._CmdBuffer[cmd_n].ElemCount)
		}
	}
	draw_list.CmdBuffer = append(draw_list.CmdBuffer, make([]ImDrawCmd, new_cmd_buffer_count)...)
	draw_list.IdxBuffer = append(draw_list.IdxBuffer, make([]ImDrawIdx, new_idx_buffer_count)...)

	// Write commands and indices in order (they are fairly small structures, we don't copy vertices only indices)
	var cmd_write = draw_list.CmdBuffer[int(len(draw_list.CmdBuffer))-new_cmd_buffer_count:]
	var idx_write = draw_list.IdxBuffer[int(len(draw_list.IdxBuffer))-new_idx_buffer_count:]
	for i := int(1); i < this._Count; i++ {
		var ch = this._Channels[i]
		if sz := len(ch._CmdBuffer); sz != 0 {
			copy(cmd_write, ch._CmdBuffer[:sz])
			cmd_write = cmd_write[sz:]
		}
		if sz := len(ch._IdxBuffer); sz != 0 {
			copy(idx_write, ch._IdxBuffer[:sz])
			idx_write = idx_write[sz:]
		}
	}
	draw_list._IdxWritePtr = int(len(draw_list.IdxBuffer) - len(idx_write))

	// Ensure there's always a non-callback draw command trailing the command-buffer
	if len(draw_list.CmdBuffer) == 0 || draw_list.CmdBuffer[len(draw_list.CmdBuffer)-1].UserCallback != nil {
		draw_list.AddDrawCmd()
	}

	// If current command is used with different settings we need to add a new command
	var curr_cmd = &draw_list.CmdBuffer[len(draw_list.CmdBuffer)-1]
	if curr_cmd.ElemCount == 0 {
		curr_cmd.HeaderCopyFromHeader(draw_list._CmdHeader)
	} else if curr_cmd.HeaderEqualsHeader(&draw_list._CmdHeader) {
		draw_list.AddDrawCmd()
	}

	this._Count = 1
}

func (this *ImDrawListSplitter) SetCurrentChannel(draw_list *ImDrawList, idx int) {
	IM_ASSERT(idx >= 0 && idx < this._Count)
	if this._Current == idx {
		return
	}

	// Swap the slice headers (pointer, length, capacity) between draw_list and the channel
	// This is equivalent to the C++ memcpy of the ImVector structure
	this._Channels[this._Current]._CmdBuffer, draw_list.CmdBuffer = draw_list.CmdBuffer, this._Channels[idx]._CmdBuffer
	this._Channels[this._Current]._IdxBuffer, draw_list.IdxBuffer = draw_list.IdxBuffer, this._Channels[idx]._IdxBuffer
	this._Current = idx
	draw_list._IdxWritePtr = int(len(draw_list.IdxBuffer))

	// If current command is used with different settings we need to add a new command
	var curr_cmd *ImDrawCmd
	if !(len(draw_list.CmdBuffer) == 0) {
		curr_cmd = &draw_list.CmdBuffer[len(draw_list.CmdBuffer)-1]
	}
	if curr_cmd == nil {
		draw_list.AddDrawCmd()
	} else if curr_cmd.ElemCount == 0 {
		curr_cmd.HeaderCopyFromHeader(draw_list._CmdHeader) // Copy ClipRect, TextureId, VtxOffset
	} else if curr_cmd.HeaderEqualsHeader(&draw_list._CmdHeader) {
		draw_list.AddDrawCmd()
	}
}
