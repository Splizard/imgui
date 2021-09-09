package imgui

type short = int16

type (
	STB_TEXTEDIT_STRING   = ImGuiInputTextState
	STB_TEXTEDIT_CHARTYPE = ImWchar
)

const (
	STB_TEXTEDIT_NEWLINE          = '\n'
	STB_TEXTEDIT_GETWIDTH_NEWLINE = -1
	STB_TEXTEDIT_UNDOSTATECOUNT   = 99
	STB_TEXTEDIT_UNDOCHARCOUNT    = 999
)

func STB_TEXTEDIT_STRINGLEN(s *STB_TEXTEDIT_STRING) int {
	panic("not implemented")
}

func STB_TEXTEDIT_LAYOUTROW(r *StbTexteditRow, obj *STB_TEXTEDIT_STRING, n int) {
	panic("not implemented")
}

func STB_TEXTEDIT_GETWIDTH(str *STB_TEXTEDIT_STRING, i, k int) float {
	panic("not implemented")
}

func STB_TEXTEDIT_GETCHAR(str *STB_TEXTEDIT_STRING, i int) STB_TEXTEDIT_CHARTYPE {
	panic("not implemented")
}

func STB_TEXTEDIT_DELETECHARS(str *STB_TEXTEDIT_STRING, pos, n int) {
	panic("not implemented")
}

func STB_TEXTEDIT_INSERTCHARS(obj *STB_TEXTEDIT_STRING, i int, c []STB_TEXTEDIT_CHARTYPE, n int) int {
	panic("not implemented")
}

func STB_TEXTEDIT_KEYTOTEXT(key STB_TEXTEDIT_KEYTYPE) STB_TEXTEDIT_CHARTYPE {
	panic("not implemented")
}

func STB_TEXTEDIT_MOVEWORDLEFT(obj *STB_TEXTEDIT_STRING, i int) int {
	panic("not implemented")
}

func STB_TEXTEDIT_MOVEWORDRIGHT(obj *STB_TEXTEDIT_STRING, i int) int {
	panic("not implemented")
}

func STB_TEXTEDIT_IS_SPACE(char STB_TEXTEDIT_CHARTYPE) bool {
	return char == ' ' || char == '\t'
}

type (
	STB_TEXTEDIT_POSITIONTYPE = int
)

type StbUndoRecord struct {
	// private data
	where         STB_TEXTEDIT_POSITIONTYPE
	insert_length STB_TEXTEDIT_POSITIONTYPE
	delete_length STB_TEXTEDIT_POSITIONTYPE
	char_storage  int
}

type StbUndoState struct {
	// private data
	undo_rec                         [STB_TEXTEDIT_UNDOSTATECOUNT]StbUndoRecord
	undo_char                        [STB_TEXTEDIT_UNDOCHARCOUNT]STB_TEXTEDIT_CHARTYPE
	undo_point, redo_point           short
	undo_char_point, redo_char_point int
}

const STB_TEXTEDIT_K_SHIFT = 1 << 10

const (
	STB_TEXTEDIT_K_LEFT = iota
	STB_TEXTEDIT_K_RIGHT
	STB_TEXTEDIT_K_UP
	STB_TEXTEDIT_K_DOWN
	STB_TEXTEDIT_K_PGUP
	STB_TEXTEDIT_K_PGDOWN
	STB_TEXTEDIT_K_LINESTART
	STB_TEXTEDIT_K_LINEEND
	STB_TEXTEDIT_K_TEXTSTART
	STB_TEXTEDIT_K_TEXTEND
	STB_TEXTEDIT_K_DELETE
	STB_TEXTEDIT_K_BACKSPACE
	STB_TEXTEDIT_K_UNDO
	STB_TEXTEDIT_K_REDO

	STB_TEXTEDIT_K_INSERT
	STB_TEXTEDIT_K_WORDLEFT
	STB_TEXTEDIT_K_WORDRIGHT
	STB_TEXTEDIT_K_LINESTART2
	STB_TEXTEDIT_K_LINEEND2
	STB_TEXTEDIT_K_TEXTSTART2
	STB_TEXTEDIT_K_TEXTEND2
)

type STB_TexteditState struct {
	/////////////////////
	//
	// public data
	//

	cursor int
	// position of the text cursor within the string

	select_start int // selection start point
	select_end   int
	// selection start and end point in characters; if equal, no selection.
	// note that start may be less than or greater than end (e.g. when
	// dragging the mouse, start is where the initial click was, and you
	// can drag in either direction)

	insert_mode byte
	// each textfield keeps its own insert mode state. to keep an app-wide
	// insert mode, copy this value in/out of the app state

	row_count_per_page int
	// page size in number of row.
	// this value MUST be set to >0 for pageup or pagedown in multilines documents.

	/////////////////////
	//
	// private data
	//
	cursor_at_end_of_line        byte // not implemented yet
	initialized                  byte
	has_preferred_x              byte
	single_line                  byte
	padding1, padding2, padding3 byte
	preferred_x                  float // this determines where the cursor up/down tries to seek to along x
	undostate                    StbUndoState
}

////////////////////////////////////////////////////////////////////////
//
//     StbTexteditRow
//
// Result of layout query, used by stb_textedit to determine where
// the text in each row is.
// result of layout query
type StbTexteditRow struct {
	x0, x1           float // starting x location, end x location (allows for align=right, etc)
	baseline_y_delta float // position of baseline relative to previous row's baseline
	ymin, ymax       float // height of row above and below baseline
	num_chars        int
}

// traverse the layout to locate the nearest character to a display position
func stb_text_locate_coord(str *STB_TEXTEDIT_STRING, x, y float) int {
	var r StbTexteditRow
	var n int = STB_TEXTEDIT_STRINGLEN(str)
	var base_y, prev_x float
	var i, k int

	r.x0 = 0
	r.x1 = 0
	r.ymin = 0
	r.ymax = 0
	r.num_chars = 0

	// search rows to find one that straddles 'y'
	for i < n {
		STB_TEXTEDIT_LAYOUTROW(&r, str, i)
		if r.num_chars <= 0 {
			return n
		}

		if i == 0 && y < base_y+r.ymin {
			return 0
		}

		if y < base_y+r.ymax {
			break
		}

		i += r.num_chars
		base_y += r.baseline_y_delta
	}

	// below all text, return 'after' last character
	if i >= n {
		return n
	}

	// check if it's before the beginning of the line
	if x < r.x0 {
		return i
	}

	// check if it's before the end of the line
	if x < r.x1 {
		// search characters in row for one that straddles 'x'
		prev_x = r.x0
		for k = 0; k < r.num_chars; k++ {
			var w float = STB_TEXTEDIT_GETWIDTH(str, i, k)
			if x < prev_x+w {
				if x < prev_x+w/2 {
					return k + i
				} else {
					return k + i + 1
				}
			}
			prev_x += w
		}
		// shouldn't happen, but if it does, fall through to end-of-line case
	}

	// if the last character is a newline, return that. otherwise return 'after' the last character
	if STB_TEXTEDIT_GETCHAR(str, i+r.num_chars-1) == STB_TEXTEDIT_NEWLINE {
		return i + r.num_chars - 1
	} else {
		return i + r.num_chars
	}
}

// API click: on mouse down, move the cursor to the clicked location, and reset the selection
func stb_textedit_click(str *STB_TEXTEDIT_STRING, state *STB_TexteditState, x, y float) {
	// In single-line mode, just always make y = 0. This lets the drag keep working if the mouse
	// goes off the top or bottom of the text
	if state.single_line != 0 {
		var r StbTexteditRow
		STB_TEXTEDIT_LAYOUTROW(&r, str, 0)
		y = r.ymin
	}

	state.cursor = stb_text_locate_coord(str, x, y)
	state.select_start = state.cursor
	state.select_end = state.cursor
	state.has_preferred_x = 0
}

// API drag: on mouse drag, move the cursor and selection endpoint to the clicked location
func stb_textedit_drag(str *STB_TEXTEDIT_STRING, state *STB_TexteditState, x, y float) {
	var p int = 0

	// In single-line mode, just always make y = 0. This lets the drag keep working if the mouse
	// goes off the top or bottom of the text
	if state.single_line != 0 {
		var r StbTexteditRow
		STB_TEXTEDIT_LAYOUTROW(&r, str, 0)
		y = r.ymin
	}

	if state.select_start == state.select_end {
		state.select_start = state.cursor
	}

	p = stb_text_locate_coord(str, x, y)
	state.cursor = p
	state.select_end = p
}

func stb_text_undo(str *STB_TEXTEDIT_STRING, state *STB_TexteditState) {
	var s *StbUndoState = &state.undostate
	var u StbUndoRecord
	var r *StbUndoRecord
	if s.undo_point == 0 {
		return
	}

	// we need to do two things: apply the undo record, and create a redo record
	u = s.undo_rec[s.undo_point-1]
	r = &s.undo_rec[s.redo_point-1]
	r.char_storage = -1

	r.insert_length = u.delete_length
	r.delete_length = u.insert_length
	r.where = u.where

	if u.delete_length != 0 {
		// if the undo record says to delete characters, then the redo record will
		// need to re-insert the characters that get deleted, so we need to store
		// them.

		// there are three cases:
		//    there's enough room to store the characters
		//    characters stored for *redoing* don't leave room for redo
		//    characters stored for *undoing* don't leave room for redo
		// if the last is true, we have to bail

		if s.undo_char_point+u.delete_length >= STB_TEXTEDIT_UNDOCHARCOUNT {
			// the undo records take up too much character space; there's no space to store the redo characters
			r.insert_length = 0
		} else {
			var i int

			// there's definitely room to store the characters eventually
			for s.undo_char_point+u.delete_length > s.redo_char_point {
				// should never happen:
				if s.redo_point == STB_TEXTEDIT_UNDOSTATECOUNT {
					return
				}
				// there's currently not enough room, so discard a redo record
				stb_textedit_discard_redo(s)
			}
			r = &s.undo_rec[s.redo_point-1]

			r.char_storage = s.redo_char_point - u.delete_length
			s.redo_char_point = s.redo_char_point - u.delete_length

			// now save the characters
			for i = 0; i < u.delete_length; i++ {
				s.undo_char[r.char_storage+i] = STB_TEXTEDIT_GETCHAR(str, u.where+i)
			}
		}

		// now we can carry out the deletion
		STB_TEXTEDIT_DELETECHARS(str, u.where, u.delete_length)
	}

	// check type of recorded action:
	if u.insert_length != 0 {
		// easy case: was a deletion, so we need to insert n characters
		STB_TEXTEDIT_INSERTCHARS(str, u.where, s.undo_char[u.char_storage:], u.insert_length)
		s.undo_char_point -= u.insert_length
	}

	state.cursor = u.where + u.insert_length

	s.undo_point--
	s.redo_point--
}
func stb_text_redo(str *STB_TEXTEDIT_STRING, state *STB_TexteditState) {
	var s *StbUndoState = &state.undostate
	var u *StbUndoRecord
	var r StbUndoRecord
	if s.redo_point == STB_TEXTEDIT_UNDOSTATECOUNT {
		return
	}

	// we need to do two things: apply the redo record, and create an undo record
	u = &s.undo_rec[s.undo_point]
	r = s.undo_rec[s.redo_point]

	// we KNOW there must be room for the undo record, because the redo record
	// was derived from an undo record

	u.delete_length = r.insert_length
	u.insert_length = r.delete_length
	u.where = r.where
	u.char_storage = -1

	if r.delete_length != 0 {
		// the redo record requires us to delete characters, so the undo record
		// needs to store the characters

		if s.undo_char_point+u.insert_length > s.redo_char_point {
			u.insert_length = 0
			u.delete_length = 0
		} else {
			var i int
			u.char_storage = s.undo_char_point
			s.undo_char_point = s.undo_char_point + u.insert_length

			// now save the characters
			for i = 0; i < u.insert_length; i++ {
				s.undo_char[u.char_storage+i] = STB_TEXTEDIT_GETCHAR(str, u.where+i)
			}
		}

		STB_TEXTEDIT_DELETECHARS(str, r.where, r.delete_length)
	}

	if r.insert_length != 0 {
		// easy case: need to insert n characters
		STB_TEXTEDIT_INSERTCHARS(str, r.where, s.undo_char[r.char_storage:], r.insert_length)
		s.redo_char_point += r.insert_length
	}

	state.cursor = r.where + r.insert_length

	s.undo_point++
	s.redo_point++
}
func stb_text_makeundo_delete(str *STB_TEXTEDIT_STRING, state *STB_TexteditState, where, length int) {
	var i int
	var p []STB_TEXTEDIT_CHARTYPE = stb_text_createundo(&state.undostate, where, length, 0)
	if p != nil {
		for i = 0; i < length; i++ {
			p[i] = STB_TEXTEDIT_GETCHAR(str, where+i)
		}
	}
}
func stb_text_makeundo_insert(state *STB_TexteditState, where, length int) {
	stb_text_createundo(&state.undostate, where, 0, length)
}
func stb_text_makeundo_replace(str *STB_TEXTEDIT_STRING, state *STB_TexteditState, where, old_length, new_length int) {
	var i int
	var p []STB_TEXTEDIT_CHARTYPE = stb_text_createundo(&state.undostate, where, old_length, new_length)
	if p != nil {
		for i = 0; i < old_length; i++ {
			p[i] = STB_TEXTEDIT_GETCHAR(str, where+i)
		}
	}
}

type StbFindState struct {
	x, y               float // position of n'th character
	height             float // height of line
	first_char, length int   // first char of row, and length
	prev_first         int   // first char of previous row
}

// find the x/y location of a character, and remember info about the previous row in
// case we get a move-up event (for page up, we'll have to rescan)
func stb_textedit_find_charpos(find *StbFindState, str *STB_TEXTEDIT_STRING, n, single_line int) {
	var r StbTexteditRow
	var prev_start int = 0
	var z int = STB_TEXTEDIT_STRINGLEN(str)
	var i, first int

	if n == z {
		// if it's at the end, then find the last line -- simpler than trying to
		// explicitly handle this case in the regular code
		if single_line != 0 {
			STB_TEXTEDIT_LAYOUTROW(&r, str, 0)
			find.y = 0
			find.first_char = 0
			find.length = z
			find.height = r.ymax - r.ymin
			find.x = r.x1
		} else {
			find.y = 0
			find.x = 0
			find.height = 1
			for i < z {
				STB_TEXTEDIT_LAYOUTROW(&r, str, i)
				prev_start = i
				i += r.num_chars
			}
			find.first_char = i
			find.length = 0
			find.prev_first = prev_start
		}
		return
	}

	// search rows to find the one that straddles character n
	find.y = 0

	for {
		STB_TEXTEDIT_LAYOUTROW(&r, str, i)
		if n < i+r.num_chars {
			break
		}
		prev_start = i
		i += r.num_chars
		find.y += r.baseline_y_delta
	}

	find.first_char = i
	first = i
	find.length = r.num_chars
	find.height = r.ymax - r.ymin
	find.prev_first = prev_start

	// now scan to find xpos
	find.x = r.x0
	for i = 0; first+i < n; i++ {
		find.x += STB_TEXTEDIT_GETWIDTH(str, first, i)
	}
}

func STB_TEXT_HAS_SELECTION(s *STB_TexteditState) bool {
	return ((s).select_start != (s).select_end)
}

// make the selection/cursor state valid if client altered the string
func stb_textedit_clamp(str *STB_TEXTEDIT_STRING, state *STB_TexteditState) {
	var n int = STB_TEXTEDIT_STRINGLEN(str)
	if STB_TEXT_HAS_SELECTION(state) {
		if state.select_start > n {
			state.select_start = n
		}
		if state.select_end > n {
			state.select_end = n
		}
		// if clamping forced them to be equal, move the cursor to match
		if state.select_start == state.select_end {
			state.cursor = state.select_start
		}
	}
	if state.cursor > n {
		state.cursor = n
	}
}

// delete characters while updating undo
func stb_textedit_delete(str *STB_TEXTEDIT_STRING, state *STB_TexteditState, where, len int) {
	stb_text_makeundo_delete(str, state, where, len)
	STB_TEXTEDIT_DELETECHARS(str, where, len)
	state.has_preferred_x = 0
}

// delete the section
func stb_textedit_delete_selection(str *STB_TEXTEDIT_STRING, state *STB_TexteditState) {
	stb_textedit_clamp(str, state)
	if STB_TEXT_HAS_SELECTION(state) {
		if state.select_start < state.select_end {
			stb_textedit_delete(str, state, state.select_start, state.select_end-state.select_start)
			state.select_end = state.select_start
			state.cursor = state.select_start
		} else {
			stb_textedit_delete(str, state, state.select_end, state.select_start-state.select_end)
			state.select_start = state.select_end
			state.cursor = state.select_end
		}
		state.has_preferred_x = 0
	}
}

// canoncialize the selection so start <= end
func stb_textedit_sortselection(state *STB_TexteditState) {
	if state.select_end < state.select_start {
		var temp int = state.select_end
		state.select_end = state.select_start
		state.select_start = temp
	}
}

// move cursor to first character of selection
func stb_textedit_move_to_first(state *STB_TexteditState) {
	if STB_TEXT_HAS_SELECTION(state) {
		stb_textedit_sortselection(state)
		state.cursor = state.select_start
		state.select_end = state.select_start
		state.has_preferred_x = 0
	}
}

// move cursor to last character of selection
func stb_textedit_move_to_last(str *STB_TEXTEDIT_STRING, state *STB_TexteditState) {
	if STB_TEXT_HAS_SELECTION(state) {
		stb_textedit_sortselection(state)
		stb_textedit_clamp(str, state)
		state.cursor = state.select_end
		state.select_start = state.select_end
		state.has_preferred_x = 0
	}
}

func is_word_boundary(str *STB_TEXTEDIT_STRING, idx int) int {
	if idx > 0 {
		return bool2int(STB_TEXTEDIT_IS_SPACE(STB_TEXTEDIT_GETCHAR(str, idx-1)) && !STB_TEXTEDIT_IS_SPACE(STB_TEXTEDIT_GETCHAR(str, idx)))
	}
	return 1
}

func stb_textedit_move_to_word_previous(str *STB_TEXTEDIT_STRING, c int) int {
	c-- // always move at least one character
	for c >= 0 && is_word_boundary(str, c) == 0 {
		c--
	}

	if c < 0 {
		c = 0
	}

	return c
}

func stb_textedit_move_to_word_next(str *STB_TEXTEDIT_STRING, c int) int {
	var len int = STB_TEXTEDIT_STRINGLEN(str)
	c++ // always move at least one character
	for c < len && is_word_boundary(str, c) == 0 {
		c++
	}

	if c > len {
		c = len
	}

	return c
}

// update selection and cursor to match each other
func stb_textedit_prep_selection_at_cursor(state *STB_TexteditState) {
	if !STB_TEXT_HAS_SELECTION(state) {
		state.select_start = state.cursor
		state.select_end = state.cursor
	} else {
		state.cursor = state.select_end
	}
}

// API cut: delete selection
func stb_textedit_cut(str *STB_TEXTEDIT_STRING, state *STB_TexteditState) int {
	if STB_TEXT_HAS_SELECTION(state) {
		stb_textedit_delete_selection(str, state) // implicitly clamps
		state.has_preferred_x = 0
		return 1
	}
	return 0
}

// API paste: replace existing selection with passed-in text
func stb_textedit_paste_internal(str *STB_TEXTEDIT_STRING, state *STB_TexteditState, text []STB_TEXTEDIT_CHARTYPE, len int) int {
	// if there's a selection, the paste should delete it
	stb_textedit_clamp(str, state)
	stb_textedit_delete_selection(str, state)
	// try to insert the characters
	if STB_TEXTEDIT_INSERTCHARS(str, state.cursor, text, len) != 0 {
		stb_text_makeundo_insert(state, state.cursor, len)
		state.cursor += len
		state.has_preferred_x = 0
		return 1
	}
	// [DEAR IMGUI]
	//// remove the undo since we didn't actually insert the characters
	//if (state.undostate.undo_point)
	//   --state.undostate.undo_point;
	// note: paste failure will leave deleted selection, may be restored with an undo (see https://github.com/nothings/stb/issues/734 for details)
	return 0
}

type STB_TEXTEDIT_KEYTYPE int

// API key: process a keyboard input
func stb_textedit_key(str *STB_TEXTEDIT_STRING, state *STB_TexteditState, key STB_TEXTEDIT_KEYTYPE) {
retry:
	switch key {
	default:
		{
			var c int = int(STB_TEXTEDIT_KEYTOTEXT(key))
			if c > 0 {
				var ch [1]STB_TEXTEDIT_CHARTYPE = [1]STB_TEXTEDIT_CHARTYPE{(STB_TEXTEDIT_CHARTYPE)(c)}

				// can't add newline in single-line mode
				if c == '\n' && state.single_line != 0 {
					break
				}

				if state.insert_mode != 0 && !STB_TEXT_HAS_SELECTION(state) && state.cursor < STB_TEXTEDIT_STRINGLEN(str) {
					stb_text_makeundo_replace(str, state, state.cursor, 1, 1)
					STB_TEXTEDIT_DELETECHARS(str, state.cursor, 1)
					if STB_TEXTEDIT_INSERTCHARS(str, state.cursor, ch[:], 1) != 0 {
						state.cursor++
						state.has_preferred_x = 0
					}
				} else {
					stb_textedit_delete_selection(str, state) // implicitly clamps
					if STB_TEXTEDIT_INSERTCHARS(str, state.cursor, ch[:], 1) != 0 {
						stb_text_makeundo_insert(state, state.cursor, 1)
						state.cursor++
						state.has_preferred_x = 0
					}
				}
			}
			break
		}

	case STB_TEXTEDIT_K_INSERT:
		state.insert_mode = byte(bool2int(state.insert_mode == 0))
		break

	case STB_TEXTEDIT_K_UNDO:
		stb_text_undo(str, state)
		state.has_preferred_x = 0
		break

	case STB_TEXTEDIT_K_REDO:
		stb_text_redo(str, state)
		state.has_preferred_x = 0
		break

	case STB_TEXTEDIT_K_LEFT:
		// if currently there's a selection, move cursor to start of selection
		if STB_TEXT_HAS_SELECTION(state) {
			stb_textedit_move_to_first(state)
		} else {
			if state.cursor > 0 {
				state.cursor--
			}
		}
		state.has_preferred_x = 0
		break

	case STB_TEXTEDIT_K_RIGHT:
		// if currently there's a selection, move cursor to end of selection
		if STB_TEXT_HAS_SELECTION(state) {
			stb_textedit_move_to_last(str, state)
		} else {
			state.cursor++
		}
		stb_textedit_clamp(str, state)
		state.has_preferred_x = 0
		break

	case STB_TEXTEDIT_K_LEFT | STB_TEXTEDIT_K_SHIFT:
		stb_textedit_clamp(str, state)
		stb_textedit_prep_selection_at_cursor(state)
		// move selection left
		if state.select_end > 0 {
			state.select_end--
		}
		state.cursor = state.select_end
		state.has_preferred_x = 0
		break

	case STB_TEXTEDIT_K_WORDLEFT:
		if STB_TEXT_HAS_SELECTION(state) {
			stb_textedit_move_to_first(state)
		} else {
			state.cursor = STB_TEXTEDIT_MOVEWORDLEFT(str, state.cursor)
			stb_textedit_clamp(str, state)
		}
		break

	case STB_TEXTEDIT_K_WORDLEFT | STB_TEXTEDIT_K_SHIFT:
		if !STB_TEXT_HAS_SELECTION(state) {
			stb_textedit_prep_selection_at_cursor(state)
		}

		state.cursor = STB_TEXTEDIT_MOVEWORDLEFT(str, state.cursor)
		state.select_end = state.cursor

		stb_textedit_clamp(str, state)
		break

	case STB_TEXTEDIT_K_WORDRIGHT:
		if STB_TEXT_HAS_SELECTION(state) {
			stb_textedit_move_to_last(str, state)
		} else {
			state.cursor = STB_TEXTEDIT_MOVEWORDRIGHT(str, state.cursor)
			stb_textedit_clamp(str, state)
		}
		break

	case STB_TEXTEDIT_K_WORDRIGHT | STB_TEXTEDIT_K_SHIFT:
		if !STB_TEXT_HAS_SELECTION(state) {
			stb_textedit_prep_selection_at_cursor(state)
		}

		state.cursor = STB_TEXTEDIT_MOVEWORDRIGHT(str, state.cursor)
		state.select_end = state.cursor

		stb_textedit_clamp(str, state)
		break

	case STB_TEXTEDIT_K_RIGHT | STB_TEXTEDIT_K_SHIFT:
		stb_textedit_prep_selection_at_cursor(state)
		// move selection right
		state.select_end++
		stb_textedit_clamp(str, state)
		state.cursor = state.select_end
		state.has_preferred_x = 0
		break

	case STB_TEXTEDIT_K_DOWN:
		fallthrough
	case STB_TEXTEDIT_K_DOWN | STB_TEXTEDIT_K_SHIFT:
		fallthrough
	case STB_TEXTEDIT_K_PGDOWN:
		fallthrough
	case STB_TEXTEDIT_K_PGDOWN | STB_TEXTEDIT_K_SHIFT:
		{
			var find StbFindState
			var row StbTexteditRow
			var i, j, sel int = 0, 0, bool2int(key&STB_TEXTEDIT_K_SHIFT != 0)
			var is_page int = bool2int((key & ^STB_TEXTEDIT_K_SHIFT) == STB_TEXTEDIT_K_PGDOWN)
			var row_count int
			if is_page != 0 {
				row_count = state.row_count_per_page
			} else {
				row_count = 1
			}

			if is_page == 0 && state.single_line != 0 {
				// on windows, up&down in single-line behave like left&right
				key = STB_TEXTEDIT_K_RIGHT | (key & STB_TEXTEDIT_K_SHIFT)
				goto retry
			}

			if sel != 0 {
				stb_textedit_prep_selection_at_cursor(state)
			} else if STB_TEXT_HAS_SELECTION(state) {
				stb_textedit_move_to_last(str, state)
			}

			// compute current position of cursor point
			stb_textedit_clamp(str, state)
			stb_textedit_find_charpos(&find, str, state.cursor, int(state.single_line))

			for j = 0; j < row_count; j++ {
				var x, goal_x float
				if state.has_preferred_x != 0 {
					goal_x = state.preferred_x
				} else {
					goal_x = find.x
				}
				var start int = find.first_char + find.length

				if find.length == 0 {
					break
				}

				// [DEAR IMGUI]
				// going down while being on the last line shouldn't bring us to that line end
				if STB_TEXTEDIT_GETCHAR(str, find.first_char+find.length-1) != STB_TEXTEDIT_NEWLINE {
					break
				}

				// now find character position down a row
				state.cursor = start
				STB_TEXTEDIT_LAYOUTROW(&row, str, state.cursor)
				x = row.x0
				for i = 0; i < row.num_chars; i++ {
					var dx float = STB_TEXTEDIT_GETWIDTH(str, start, i)
					if dx == STB_TEXTEDIT_GETWIDTH_NEWLINE {
						break
					}
					x += dx
					if x > goal_x {
						break
					}
					state.cursor++
				}
				stb_textedit_clamp(str, state)

				state.has_preferred_x = 1
				state.preferred_x = goal_x

				if sel != 0x80 {
					state.select_end = state.cursor
				}

				// go to next line
				find.first_char = find.first_char + find.length
				find.length = row.num_chars
			}
			break
		}

	case STB_TEXTEDIT_K_UP:
		fallthrough
	case STB_TEXTEDIT_K_UP | STB_TEXTEDIT_K_SHIFT:
		fallthrough
	case STB_TEXTEDIT_K_PGUP:
		fallthrough
	case STB_TEXTEDIT_K_PGUP | STB_TEXTEDIT_K_SHIFT:
		{
			var find StbFindState
			var row StbTexteditRow
			var i, j, prev_scan, sel int = 0, 0, 0, bool2int(key&STB_TEXTEDIT_K_SHIFT != 0)
			var is_page int = bool2int((key & ^STB_TEXTEDIT_K_SHIFT) == STB_TEXTEDIT_K_PGUP)
			var row_count int
			if is_page != 0 {
				row_count = state.row_count_per_page
			} else {
				row_count = 1
			}

			if is_page == 0 && state.single_line != 0 {
				// on windows, up&down become left&right
				key = STB_TEXTEDIT_K_LEFT | (key & STB_TEXTEDIT_K_SHIFT)
				goto retry
			}

			if sel != 0 {
				stb_textedit_prep_selection_at_cursor(state)
			} else if STB_TEXT_HAS_SELECTION(state) {
				stb_textedit_move_to_first(state)
			}

			// compute current position of cursor point
			stb_textedit_clamp(str, state)
			stb_textedit_find_charpos(&find, str, state.cursor, int(state.single_line))

			for j = 0; j < row_count; j++ {
				var x, goal_x float
				if state.has_preferred_x != 0 {
					goal_x = state.preferred_x
				} else {
					goal_x = find.x
				}

				// can only go up if there's a previous row
				if find.prev_first == find.first_char {
					break
				}

				// now find character position up a row
				state.cursor = find.prev_first
				STB_TEXTEDIT_LAYOUTROW(&row, str, state.cursor)
				x = row.x0
				for i = 0; i < row.num_chars; i++ {
					var dx float = STB_TEXTEDIT_GETWIDTH(str, find.prev_first, i)
					if dx == STB_TEXTEDIT_GETWIDTH_NEWLINE {
						break
					}
					x += dx
					if x > goal_x {
						break
					}
					state.cursor++
				}
				stb_textedit_clamp(str, state)

				state.has_preferred_x = 1
				state.preferred_x = goal_x

				if sel != 0 {
					state.select_end = state.cursor
				}

				// go to previous line
				// (we need to scan previous line the hard way. maybe we could expose this as a new API function?)
				if find.prev_first > 0 {
					prev_scan = find.prev_first - 1
				} else {
					prev_scan = 0
				}
				for prev_scan > 0 && STB_TEXTEDIT_GETCHAR(str, prev_scan-1) != STB_TEXTEDIT_NEWLINE {
					prev_scan--
				}
				find.first_char = find.prev_first
				find.prev_first = prev_scan
			}
			break
		}

	case STB_TEXTEDIT_K_DELETE:
		fallthrough
	case STB_TEXTEDIT_K_DELETE | STB_TEXTEDIT_K_SHIFT:
		if STB_TEXT_HAS_SELECTION(state) {
			stb_textedit_delete_selection(str, state)
		} else {
			var n int = STB_TEXTEDIT_STRINGLEN(str)
			if state.cursor < n {
				stb_textedit_delete(str, state, state.cursor, 1)
			}
		}
		state.has_preferred_x = 0
		break

	case STB_TEXTEDIT_K_BACKSPACE:
		fallthrough
	case STB_TEXTEDIT_K_BACKSPACE | STB_TEXTEDIT_K_SHIFT:
		if STB_TEXT_HAS_SELECTION(state) {
			stb_textedit_delete_selection(str, state)
		} else {
			stb_textedit_clamp(str, state)
			if state.cursor > 0 {
				stb_textedit_delete(str, state, state.cursor-1, 1)
				state.cursor--
			}
		}
		state.has_preferred_x = 0
		break

	case STB_TEXTEDIT_K_TEXTSTART2:
		fallthrough
	case STB_TEXTEDIT_K_TEXTSTART:
		state.cursor = 0
		state.select_start = 0
		state.select_end = 0
		state.has_preferred_x = 0
		break

	case STB_TEXTEDIT_K_TEXTEND2:
		fallthrough
	case STB_TEXTEDIT_K_TEXTEND:
		state.cursor = STB_TEXTEDIT_STRINGLEN(str)
		state.select_start = 0
		state.select_end = 0
		state.has_preferred_x = 0
		break

	case STB_TEXTEDIT_K_TEXTSTART2 | STB_TEXTEDIT_K_SHIFT:
		fallthrough
	case STB_TEXTEDIT_K_TEXTSTART | STB_TEXTEDIT_K_SHIFT:
		stb_textedit_prep_selection_at_cursor(state)
		state.cursor = 0
		state.select_end = 0
		state.has_preferred_x = 0
		break

	case STB_TEXTEDIT_K_TEXTEND2 | STB_TEXTEDIT_K_SHIFT:
		fallthrough
	case STB_TEXTEDIT_K_TEXTEND | STB_TEXTEDIT_K_SHIFT:
		stb_textedit_prep_selection_at_cursor(state)
		state.cursor = STB_TEXTEDIT_STRINGLEN(str)
		state.select_end = STB_TEXTEDIT_STRINGLEN(str)
		state.has_preferred_x = 0
		break

	case STB_TEXTEDIT_K_LINESTART2:
		fallthrough
	case STB_TEXTEDIT_K_LINESTART:
		stb_textedit_clamp(str, state)
		stb_textedit_move_to_first(state)
		if state.single_line != 0 {
			state.cursor = 0
		} else {
			for state.cursor > 0 && STB_TEXTEDIT_GETCHAR(str, state.cursor-1) != STB_TEXTEDIT_NEWLINE {
				state.cursor--
			}
		}
		state.has_preferred_x = 0
		break

	case STB_TEXTEDIT_K_LINEEND2:
		fallthrough
	case STB_TEXTEDIT_K_LINEEND:
		{
			var n int = STB_TEXTEDIT_STRINGLEN(str)
			stb_textedit_clamp(str, state)
			stb_textedit_move_to_first(state)
			if state.single_line != 0 {
				state.cursor = n
			} else {
				for state.cursor < n && STB_TEXTEDIT_GETCHAR(str, state.cursor) != STB_TEXTEDIT_NEWLINE {
					state.cursor++
				}
			}
			state.has_preferred_x = 0
			break
		}

	case STB_TEXTEDIT_K_LINESTART2 | STB_TEXTEDIT_K_SHIFT:
		fallthrough
	case STB_TEXTEDIT_K_LINESTART | STB_TEXTEDIT_K_SHIFT:
		stb_textedit_clamp(str, state)
		stb_textedit_prep_selection_at_cursor(state)
		if state.single_line != 0 {
			state.cursor = 0
		} else {
			for state.cursor > 0 && STB_TEXTEDIT_GETCHAR(str, state.cursor-1) != STB_TEXTEDIT_NEWLINE {
				state.cursor--
			}
		}
		state.select_end = state.cursor
		state.has_preferred_x = 0
		break

	case STB_TEXTEDIT_K_LINEEND2 | STB_TEXTEDIT_K_SHIFT:
		fallthrough
	case STB_TEXTEDIT_K_LINEEND | STB_TEXTEDIT_K_SHIFT:
		{
			var n int = STB_TEXTEDIT_STRINGLEN(str)
			stb_textedit_clamp(str, state)
			stb_textedit_prep_selection_at_cursor(state)
			if state.single_line != 0 {
				state.cursor = n
			} else {
				for state.cursor < n && STB_TEXTEDIT_GETCHAR(str, state.cursor) != STB_TEXTEDIT_NEWLINE {
					state.cursor++
				}
			}
			state.select_end = state.cursor
			state.has_preferred_x = 0
			break
		}
	}
}

func stb_textedit_flush_redo(state *StbUndoState) {
	state.redo_point = STB_TEXTEDIT_UNDOSTATECOUNT
	state.redo_char_point = STB_TEXTEDIT_UNDOCHARCOUNT
}

// discard the oldest entry in the undo list
func stb_textedit_discard_undo(state *StbUndoState) {
	if state.undo_point > 0 {
		// if the 0th undo state has characters, clean those up
		if state.undo_rec[0].char_storage >= 0 {
			var n, i int = state.undo_rec[0].insert_length, 0
			// delete n characters from all other records
			state.undo_char_point -= n
			copy(state.undo_char[:], state.undo_char[n:n+state.undo_char_point]) //TODO/FIXME is this correct?
			for i = 0; i < int(state.undo_point); i++ {
				if state.undo_rec[i].char_storage >= 0 {
					state.undo_rec[i].char_storage -= n // @OPTIMIZE: get rid of char_storage and infer it
				}
			}
		}
		state.undo_point--
		copy(state.undo_rec[:], state.undo_rec[1:1+state.undo_point]) //TODO/FIXME is this correct?
	}
}

// discard the oldest entry in the redo list--it's bad if this
// ever happens, but because undo & redo have to store the actual
// characters in different cases, the redo character buffer can
// fill up even though the undo buffer didn't
func stb_textedit_discard_redo(state *StbUndoState) {
	var k int = STB_TEXTEDIT_UNDOSTATECOUNT - 1

	if int(state.redo_point) <= k {
		// if the k'th undo state has characters, clean those up
		if state.undo_rec[k].char_storage >= 0 {
			var n, i int = state.undo_rec[k].insert_length, 0
			// move the remaining redo character data to the end of the buffer
			state.redo_char_point += n
			copy(state.undo_char[state.redo_char_point:], state.undo_char[state.redo_char_point-n:state.redo_char_point-n+(STB_TEXTEDIT_UNDOCHARCOUNT-state.redo_char_point)])
			// adjust the position of all the other records to account for above memmove
			for i = int(state.redo_point); i < k; i++ {
				if state.undo_rec[i].char_storage >= 0 {
					state.undo_rec[i].char_storage += n
				}
			}
		}
		// now move all the redo records towards the end of the buffer; the first one is at 'redo_point'
		// [DEAR IMGUI]
		var move_size uintptr = (uintptr)((STB_TEXTEDIT_UNDOSTATECOUNT - state.redo_point - 1))
		//var buf_begin []byte = state.undo_rec[:]
		//var buf_end []byte = state.undo_rec[len(state.undo_rec):]
		//IM_ASSERT(((char*)(state.undo_rec + state.redo_point)) >= buf_begin); TODO/FIXME
		//IM_ASSERT(((char*)(state.undo_rec + state.redo_point + 1) + move_size) <= buf_end); TODO/FIXME
		copy(state.undo_rec[:state.redo_point+1], state.undo_rec[state.redo_point:uintptr(state.redo_point)+move_size])

		// now move redo_point to point to the new one
		state.redo_point++
	}
}

func stb_text_create_undo_record(state *StbUndoState, numchars int) *StbUndoRecord {
	// any time we create a new undo record, we discard redo
	stb_textedit_flush_redo(state)

	// if we have no free records, we have to make room, by sliding the
	// existing records down
	if state.undo_point == STB_TEXTEDIT_UNDOSTATECOUNT {
		stb_textedit_discard_undo(state)
	}

	// if the characters to store won't possibly fit in the buffer, we can't undo
	if numchars > STB_TEXTEDIT_UNDOCHARCOUNT {
		state.undo_point = 0
		state.undo_char_point = 0
		return nil
	}

	// if we don't have enough free characters in the buffer, we have to make room
	for state.undo_char_point+numchars > STB_TEXTEDIT_UNDOCHARCOUNT {
		stb_textedit_discard_undo(state)
	}

	state.undo_point++
	return &state.undo_rec[state.undo_point]
}

func stb_text_createundo(state *StbUndoState, pos, insert_len, delete_len int) []STB_TEXTEDIT_CHARTYPE {
	var r *StbUndoRecord = stb_text_create_undo_record(state, insert_len)
	if r == nil {
		return nil
	}

	r.where = pos
	r.insert_length = (STB_TEXTEDIT_POSITIONTYPE)(insert_len)
	r.delete_length = (STB_TEXTEDIT_POSITIONTYPE)(delete_len)

	if insert_len == 0 {
		r.char_storage = -1
		return nil
	} else {
		r.char_storage = state.undo_char_point
		state.undo_char_point += insert_len
		return state.undo_char[r.char_storage:]
	}
}

// reset the state to default
func stb_textedit_clear_state(state *STB_TexteditState, is_single_line int) {
	state.undostate.undo_point = 0
	state.undostate.undo_char_point = 0
	state.undostate.redo_point = STB_TEXTEDIT_UNDOSTATECOUNT
	state.undostate.redo_char_point = STB_TEXTEDIT_UNDOCHARCOUNT
	state.select_end = 0
	state.select_start = 0
	state.cursor = 0
	state.has_preferred_x = 0
	state.preferred_x = 0
	state.cursor_at_end_of_line = 0
	state.initialized = 1
	state.single_line = (byte)(is_single_line)
	state.insert_mode = 0
	state.row_count_per_page = 0
}

// API initialize
func stb_textedit_initialize_state(state *STB_TexteditState, is_single_line int) {
	stb_textedit_clear_state(state, is_single_line)
}

func stb_textedit_paste(str *STB_TEXTEDIT_STRING, state *STB_TexteditState, ctext []STB_TEXTEDIT_CHARTYPE, len int) int {
	return stb_textedit_paste_internal(str, state, ctext, len)
}

/*
------------------------------------------------------------------------------
This software is available under 2 licenses -- choose whichever you prefer.
------------------------------------------------------------------------------
ALTERNATIVE A - MIT License
Copyright (c) 2017 Sean Barrett
Permission is hereby granted, free of charge, to any person obtaining a copy of
this software and associated documentation files (the "Software"), to deal in
the Software without restriction, including without limitation the rights to
use, copy, modify, merge, publish, distribute, sublicense, and/or sell copies
of the Software, and to permit persons to whom the Software is furnished to do
so, subject to the following conditions:
The above copyright notice and this permission notice shall be included in all
copies or substantial portions of the Software.
THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
SOFTWARE.
------------------------------------------------------------------------------
ALTERNATIVE B - Public Domain (www.unlicense.org)
This is free and unencumbered software released into the public domain.
Anyone is free to copy, modify, publish, use, compile, sell, or distribute this
software, either in source code form or as a compiled binary, for any purpose,
commercial or non-commercial, and by any means.
In jurisdictions that recognize copyright laws, the author or authors of this
software dedicate any and all copyright interest in the software to the public
domain. We make this dedication for the benefit of the public at large and to
the detriment of our heirs and successors. We intend this dedication to be an
overt act of relinquishment in perpetuity of all present and future rights to
this software under copyright law.
THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER IN AN
ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN CONNECTION
WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.
------------------------------------------------------------------------------
*/
