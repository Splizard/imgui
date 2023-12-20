package stbrp

import (
	"reflect"
	"sort"
	"unsafe"

	"github.com/Splizard/imgui/golang"
)

//PORTING STATUS = DONE - Quentin Quaadgras

const _DEBUG = false

type double = float64
type int = int32
type uint = uint32
type float = float32
type size_t = uintptr
type char = byte

func isfalse(x int) bool {
	return x == 0
}

func istrue(x int) bool {
	return x != 0
}

func bool2int(b bool) int {
	if b {
		return 1
	}
	return 0
}

const STB_RECT_PACK_VERSION = 1

const STBRP_LARGE_RECTS = false

type Coord uint16

const STBRP__MAXVAL = 0xffff

// PackRects Assign packed locations to rectangles. The rectangles are of type
// 'stbrp_rect' defined below, stored in the array 'rects', and there
// are 'num_rects' many of them.
//
// Rectangles which are successfully packed have the 'was_packed' flag
// set to a non-zero value and 'x' and 'y' store the minimum location
// on each axis (i.e. bottom-left in cartesian coordinates, top-left
// if you imagine y increasing downwards). Rectangles which do not fit
// have the 'was_packed' flag set to 0.
//
// You should not try to access the 'rects' array from another thread
// while this function is running, as the function temporarily reorders
// the array while it executes.
//
// To pack into another rectangle, you need to call stbrp_init_target
// again. To continue packing into the same rectangle, you can call
// this function again. Calling this multiple times with multiple rect
// arrays will probably produce worse packing results than calling it
// a single time with the full rectangle array, but the option is
// available.
//
// The function returns 1 if all of the rectangles were successfully
// packed and 0 otherwise.
func PackRects(context *Context, rects []Rect, num_rects int) int {
	var i, all_rects_packed int = 0, 1

	// we use the 'was_packed' field internally to allow sorting/unsorting
	for i = 0; i < num_rects; i++ {
		rects[i].WasPacked = i
	}

	// sort according to heuristic
	STBRP_SORT(rects, num_rects, unsafe.Sizeof(rects[0]), rect_height_compare)

	for i = 0; i < num_rects; i++ {
		if rects[i].W == 0 || rects[i].H == 0 {
			rects[i].X, rects[i].Y = 0, 0 // empty rect needs no space
		} else {
			var fr = stbrp__skyline_pack_rectangle(context, int(rects[i].W), int(rects[i].H))
			if fr.prev_link != nil {
				rects[i].X = Coord(fr.x)
				rects[i].Y = Coord(fr.y)
			} else {
				rects[i].X, rects[i].Y = STBRP__MAXVAL, STBRP__MAXVAL
			}
		}
	}

	// unsort
	STBRP_SORT(rects, num_rects, unsafe.Sizeof(rects[0]), rect_original_order)

	// set was_packed flags and all_rects_packed status
	for i = 0; i < num_rects; i++ {
		rects[i].WasPacked = bool2int(!(rects[i].X == STBRP__MAXVAL && rects[i].Y == STBRP__MAXVAL))
		if isfalse(rects[i].WasPacked) {
			all_rects_packed = 0
		}
	}

	// return the all_rects_packed status
	return all_rects_packed
}

type Rect struct {
	// reserved for your use:
	ID int

	// input:
	W, H Coord

	// output:
	X, Y Coord

	WasPacked int
}

// InitTarget Initialize a rectangle packer to:
//
//	pack a rectangle that is 'width' by 'height' in dimensions
//	using temporary storage provided by the array 'nodes', which is 'num_nodes' long
//
// You must call this function every time you start packing into a new target.
//
// There is no "shutdown" function. The 'nodes' memory must stay valid for
// the following stbrp_pack_rects() call (or calls), but can be freed after
// the call (or calls) finish.
//
// Note: to guarantee best results, either:
//  1. make sure 'num_nodes' >= 'width'
//     or  2. call stbrp_allow_out_of_mem() defined below with 'allow_out_of_mem = 1'
//
// If you don't do either of the above things, widths will be quantized to multiples
// of small integers to guarantee the algorithm doesn't run out of temporary storage.
//
// If you do #2, then the non-quantized algorithm will be used, but the algorithm
// may run out of temporary storage and be unable to pack some rectangles.
func InitTarget(context *Context, width, height int, nodes []Node, num_nodes int) {
	var i int
	if !STBRP_LARGE_RECTS {
		STBRP_ASSERT(width <= 0xffff && height <= 0xffff)
	}

	for i = 0; i < num_nodes-1; i++ {
		nodes[i].next = &nodes[i+1]
	}
	nodes[i].next = nil
	context.init_mode = STBRP__INIT_skyline
	context.heuristic = STBRP_HEURISTIC_Skyline_default
	context.free_head = &nodes[0]
	context.active_head = &context.extra[0]
	context.width = width
	context.height = height
	context.num_nodes = num_nodes
	stbrp_setup_allow_out_of_mem(context, 0)

	// node 0 is the full width, node 1 is the sentinel (lets us not store width explicitly)
	context.extra[0].x = 0
	context.extra[0].y = 0
	context.extra[0].next = &context.extra[1]
	context.extra[1].x = Coord(width)
	if !STBRP_LARGE_RECTS {
		//context.extra[1].y = (1 << 30)
	} else {
		context.extra[1].y = 65535
	}
	context.extra[1].next = nil
}

// Optionally call this function after init but before doing any packing to
// change the handling of the out-of-temp-memory scenario, described above.
// If you call init again, this will be reset to the default (false).
func stbrp_setup_allow_out_of_mem(context *Context, allow_out_of_mem int) {
	if istrue(allow_out_of_mem) {
		// if it's ok to run out of memory, then don't bother aligning them;
		// this gives better packing, but may fail due to OOM (even though
		// the rectangles easily fit). @TODO a smarter approach would be to only
		// quantize once we've hit OOM, then we could get rid of this parameter.
		context.align = 1
	} else {
		// if it's not ok to run out of memory, then quantize the widths
		// so that num_nodes is always enough nodes.
		//
		// I.e. num_nodes * align >= width
		//                  align >= width / num_nodes
		//                  align = ceil(width/num_nodes)

		context.align = (context.width + context.num_nodes - 1) / context.num_nodes
	}
}

// Optionally select which packing heuristic the library should use. Different
// heuristics will produce better/worse results for different data sets.
// If you call init again, this will be reset to the default.
func stbrp_setup_heuristic(context *Context, heuristic int) {
	switch context.init_mode {
	case STBRP__INIT_skyline:
		STBRP_ASSERT(heuristic == STBRP_HEURISTIC_Skyline_BL_sortHeight || heuristic == STBRP_HEURISTIC_Skyline_BF_sortHeight)
		context.heuristic = heuristic
		break
	default:
		STBRP_ASSERT(false)
	}
}

const (
	STBRP_HEURISTIC_Skyline_default       = 0
	STBRP_HEURISTIC_Skyline_BL_sortHeight = STBRP_HEURISTIC_Skyline_default
	STBRP_HEURISTIC_Skyline_BF_sortHeight = 1
)

type Node struct {
	x, y Coord
	next *Node
}

type Context struct {
	width       int
	height      int
	align       int
	init_mode   int
	heuristic   int
	num_nodes   int
	active_head *Node
	free_head   *Node
	extra       [2]Node
}

func STBRP_ASSERT(cond bool) {
	if !cond {
		panic("assert failed")
	}
}

func STBRP_SORT(slice interface{}, _ int, _ uintptr, compare func(a, b interface{}) int) {
	sort.Slice(slice, func(i, j golang.Int) bool {
		return compare(reflect.ValueOf(slice).Index(i).Addr().Interface(), reflect.ValueOf(slice).Index(j).Addr().Interface()) < 0
	})
}

const (
	STBRP__INIT_skyline = 1
)

// find minimum y position if it starts at x1
func stbrp__skyline_find_min_y(c *Context, first *Node, x0, width int, pwaste *int) int {
	var node = first
	var x1 = x0 + width
	var min_y, visited_width, waste_area int

	STBRP_ASSERT(int(first.x) <= x0)
	STBRP_ASSERT(int(node.next.x) > x0)
	STBRP_ASSERT(int(node.x) <= x0)

	min_y = 0
	waste_area = 0
	visited_width = 0
	for int(node.x) < x1 {
		if int(node.y) > min_y {
			// raise min_y higher.
			// we've accounted for all waste up to min_y,
			// but we'll now add more waste for everything we've visted
			waste_area += visited_width * (int(node.y) - min_y)
			min_y = int(node.y)
			// the first time through, visited_width might be reduced
			if int(node.x) < x0 {
				visited_width += int(node.next.x) - x0
			} else {
				visited_width += int(node.next.x) - int(node.x)
			}
		} else {
			// add waste area
			var under_width = int(node.next.x) - int(node.x)
			if under_width+visited_width > width {
				under_width = width - visited_width
			}
			waste_area += under_width * (min_y - int(node.y))
			visited_width += under_width
		}
		node = node.next
	}

	*pwaste = waste_area
	return min_y
}

type stbrp__findresult struct {
	x, y      int
	prev_link **Node
}

func stbrp__skyline_find_best_pos(c *Context, width, height int) stbrp__findresult {
	var best_waste, best_x, best_y int = 1 << 30, 0, 1 << 30
	var fr stbrp__findresult
	var prev, best **Node
	var node, tail *Node

	// align to multiple of c.align
	width = width + c.align - 1
	width -= width % c.align
	STBRP_ASSERT(width%c.align == 0)

	// if it can't possibly fit, bail immediately
	if width > c.width || height > c.height {
		fr.prev_link = nil
		fr.x, fr.y = 0, 0
		return fr
	}

	node = c.active_head
	prev = &c.active_head
	for int(node.x)+width <= c.width {
		var y, waste int
		y = stbrp__skyline_find_min_y(c, node, int(node.x), width, &waste)
		if c.heuristic == STBRP_HEURISTIC_Skyline_BL_sortHeight { // actually just want to test BL
			// bottom left
			if y < best_y {
				best_y = y
				best = prev
			}
		} else {
			// best-fit
			if y+height <= c.height {
				// can only use it if it first vertically
				if y < best_y || (y == best_y && waste < best_waste) {
					best_y = y
					best_waste = waste
					best = prev
				}
			}
		}
		prev = &node.next
		node = node.next
	}

	if best == nil {
		best_x = 0
	} else {
		best_x = int((*best).x)
	}

	// if doing best-fit (BF), we also have to try aligning right edge to each node position
	//
	// e.g, if fitting
	//
	//     ____________________
	//    |____________________|
	//
	//            into
	//
	//   |                         |
	//   |             ____________|
	//   |____________|
	//
	// then right-aligned reduces waste, but bottom-left BL is always chooses left-aligned
	//
	// This makes BF take about 2x the time

	if c.heuristic == STBRP_HEURISTIC_Skyline_BF_sortHeight {
		tail = c.active_head
		node = c.active_head
		prev = &c.active_head
		// find first node that's admissible
		for int(tail.x) < width {
			tail = tail.next
		}
		for tail != nil {
			var xpos = int(tail.x) - width
			var y, waste int
			STBRP_ASSERT(xpos >= 0)
			// find the left position that matches this
			for int(node.next.x) <= xpos {
				prev = &node.next
				node = node.next
			}
			STBRP_ASSERT(int(node.next.x) > xpos && int(node.x) <= xpos)
			y = stbrp__skyline_find_min_y(c, node, xpos, width, &waste)
			if y+height <= c.height {
				if y <= best_y {
					if y < best_y || waste < best_waste || (waste == best_waste && xpos < best_x) {
						best_x = xpos
						STBRP_ASSERT(y <= best_y)
						best_y = y
						best_waste = waste
						best = prev
					}
				}
			}
			tail = tail.next
		}
	}

	fr.prev_link = best
	fr.x = best_x
	fr.y = best_y
	return fr
}

func stbrp__skyline_pack_rectangle(context *Context, width, height int) stbrp__findresult {
	// find best position according to heuristic
	var res = stbrp__skyline_find_best_pos(context, width, height)
	var node, cur *Node

	// bail if:
	//    1. it failed
	//    2. the best node doesn't fit (we don't always check this)
	//    3. we're out of memory
	if res.prev_link == nil || res.y+height > context.height || context.free_head == nil {
		res.prev_link = nil
		return res
	}

	// on success, create new node
	node = context.free_head
	node.x = Coord(res.x)
	node.y = Coord(res.y + height)

	context.free_head = node.next

	// insert the new node into the right starting point, and
	// let 'cur' point to the remaining nodes needing to be
	// stiched back in

	cur = *res.prev_link
	if int(cur.x) < res.x {
		// preserve the existing one, so start testing with the next one
		var next = cur.next
		cur.next = node
		cur = next
	} else {
		*res.prev_link = node
	}

	// from here, traverse cur and free the nodes, until we get to one
	// that shouldn't be freed
	for cur.next != nil && int(cur.next.x) <= res.x+width {
		var next = cur.next
		// move the current node to the free list
		cur.next = context.free_head
		context.free_head = cur
		cur = next
	}

	// stitch the list back in
	node.next = cur

	if int(cur.x) < res.x+width {
		cur.x = Coord(res.x + width)
	}

	if _DEBUG {
		cur = context.active_head
		for int(cur.x) < context.width {
			STBRP_ASSERT(cur.x < cur.next.x)
			cur = cur.next
		}
		STBRP_ASSERT(cur.next == nil)

		{
			var count int = 0
			cur = context.active_head
			for cur != nil {
				cur = cur.next
				count++
			}
			cur = context.free_head
			for cur != nil {
				cur = cur.next
				count++
			}
			STBRP_ASSERT(count == context.num_nodes+2)
		}
	}

	return res
}

func rect_height_compare(a, b interface{}) int {
	var p = a.(*Rect)
	var q = b.(*Rect)
	if p.H > q.H {
		return -1
	}
	if p.H < q.H {
		return 1
	}

	if p.W > q.W {
		return -1
	}

	if p.W < q.W {
		return 1
	}
	return -1
}

func rect_original_order(a, b interface{}) int {
	var p = a.(*Rect)
	var q = b.(*Rect)

	if p.WasPacked < q.WasPacked {
		return -1
	}

	if p.WasPacked > q.WasPacked {
		return 1
	}

	return -1
}

/*
------------------------------------------------------------------------------
This software is available under 2 licenses -- choose whichever you prefer.
------------------------------------------------------------------------------
ALTERNATIVE A - MIT License
Copyright (c) 2017 Sean Barrett
Copyright (c) 2021 Quentin Quaadgras
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
