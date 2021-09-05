package stbrp

import "sort"

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

// Mostly for internal use, but this is the maximum supported coordinate value.
const maxVal = 0x7fffffff

const debug = false

type initMode int32

const (
	initSkyline initMode = 1
)

type Coord = int32

type Rect struct {
	ID int32 // reserved for your use

	// input:
	W, H Coord

	// output:
	X, Y Coord

	was_packed int32 // non-zero if valid packing

} // 16 bytes, nominally

type Node struct {
	x, y Coord
	next *Node
}

type Context struct {
	width, height          int32
	align                  int32
	init_mode              initMode
	heuristic              Heuristic
	num_nodes              int32
	active_head, free_head *Node
	extra                  [2]Node
}

// Assign packed locations to rectangles. The rectangles are of type
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
func PackRects(ctx *Context, rects []Rect) int32 {
	var all_rects_packed int32 = 1

	// we use the 'was_packed' field internally to allow sorting/unsorting
	for i := range rects {
		rects[i].was_packed = int32(i)
	}

	//heightCompare
	sort.Slice(rects, func(i, j int) bool {
		var p, q = rects[i], rects[j]
		if p.H > q.H {
			return true
		}
		if p.H < q.H {
			return false
		}
		if p.W > q.W {
			return true
		}
		return false
	})

	for i := range rects {
		if rects[i].W == 0 || rects[i].H == 0 {
			rects[i].X = 0
			rects[i].Y = 0 // empty rect needs no space
		} else {
			fr := skylinePackRectangle(ctx, rects[i].W, rects[i].H)
			if fr.prev_link != nil {
				rects[i].X = fr.x
				rects[i].Y = fr.y
			} else {
				rects[i].X = maxVal
				rects[i].Y = maxVal
			}
		}
	}

	// unsort
	sort.Slice(rects, func(i, j int) bool {
		return rects[i].was_packed < rects[j].was_packed
	})

	// set was_packed flags and all_rects_packed status
	for i := range rects {
		if !(rects[i].X == maxVal && rects[i].Y == maxVal) {
			rects[i].was_packed = 1
		}
		if rects[i].was_packed == 0 {
			all_rects_packed = 0
		}
	}

	// return the all_rects_packed status
	return all_rects_packed
}

// Initialize a rectangle packer to:
//    pack a rectangle that is 'width' by 'height' in dimensions
//    using temporary storage provided by the array 'nodes', which is 'num_nodes' long
//
// You must call this function every time you start packing into a new target.
//
// There is no "shutdown" function. The 'nodes' memory must stay valid for
// the following stbrp_pack_rects() call (or calls), but can be freed after
// the call (or calls) finish.
//
// Note: to guarantee best results, either:
//       1. make sure 'num_nodes' >= 'width'
//   or  2. call stbrp_allow_out_of_mem() defined below with 'allow_out_of_mem = 1'
//
// If you don't do either of the above things, widths will be quantized to multiples
// of small integers to guarantee the algorithm doesn't run out of temporary storage.
//
// If you do #2, then the non-quantized algorithm will be used, but the algorithm
// may run out of temporary storage and be unable to pack some rectangles.
func InitTarget(ctx *Context, width, height int32, nodes []Node) {
	var i int
	for i = range nodes {
		nodes[i].next = &nodes[i+1]
	}

	nodes[i].next = nil
	ctx.init_mode = initSkyline
	ctx.free_head = &nodes[0]
	ctx.active_head = &ctx.extra[0]
	ctx.width = width
	ctx.height = height
	ctx.num_nodes = int32(len(nodes))
	SetupAllowOutOfMem(ctx, false)

	// node 0 is the full width, node 1 is the sentinel (lets us not store width explicitly)
	ctx.extra[0].x = 0
	ctx.extra[0].y = 0
	ctx.extra[0].next = &ctx.extra[1]
	ctx.extra[1].x = Coord(width)
	ctx.extra[1].y = (1 << 30)
	ctx.extra[1].next = nil
}

// Optionally call this function after init but before doing any packing to
// change the handling of the out-of-temp-memory scenario, described above.
// If you call init again, this will be reset to the default (false).
func SetupAllowOutOfMem(ctx *Context, allow_out_of_mem bool) {
	if allow_out_of_mem {
		// if it's ok to run out of memory, then don't bother aligning them;
		// this gives better packing, but may fail due to OOM (even though
		// the rectangles easily fit). @TODO a smarter approach would be to only
		// quantize once we've hit OOM, then we could get rid of this parameter.
		ctx.align = 1
	} else {
		// if it's not ok to run out of memory, then quantize the widths
		// so that num_nodes is always enough nodes.
		//
		// I.e. num_nodes * align >= width
		//                  align >= width / num_nodes
		//                  align = ceil(width/num_nodes)
		ctx.align = (ctx.width + ctx.num_nodes - 1) / ctx.num_nodes
	}
}

type Heuristic int

const (
	SkylineBLSortHeight Heuristic = iota
	SkylineBFSortHeight
)

// Optionally select which packing heuristic the library should use. Different
// heuristics will produce better/worse results for different data sets.
// If you call init again, this will be reset to the default.
func SetupHeuristic(ctx *Context, heuristic Heuristic) {
	switch ctx.init_mode {
	case initSkyline:
		switch heuristic {
		case SkylineBLSortHeight, SkylineBFSortHeight:
			ctx.heuristic = heuristic
		default:
			panic("invalid heuristic")
		}
	default:
		panic("invalid init mode")
	}
}

type findResult struct {
	x, y      int32
	prev_link **Node
}

//find minimum y position if it starts at x1
func skylineFindMinY(ctx *Context, first *Node, x0, width int32) (min_y int32, pwaste int32) {
	var node = first
	var x1 = x0 + width
	var visited_width, waste_area int32

	if first.x > x0 {
		panic("first.x > x0")
	}

	if false {
		// skip in case we're past the node
		for node.next.x <= x0 {
			node = node.next
		}
	} else {
		// we ended up handling this in the caller for efficiency
		if node.next.x <= x0 {
			panic("node.next.x <= x0")
		}
	}

	if node.x > x0 {
		panic("node.x > x0")
	}

	min_y = 0
	waste_area = 0
	visited_width = 0
	for node.x < x1 {
		if node.y > min_y {
			// raise min_y higher.
			// we've accounted for all waste up to min_y,
			// but we'll now add more waste for everything we've visted
			waste_area += visited_width * (node.y - min_y)
			min_y = node.y
			// the first time through, visited_width might be reduced
			if node.x < x0 {
				visited_width += node.next.x - x0
			} else {
				visited_width += node.next.x - node.x
			}
		} else {
			// add waste area
			var under_width = node.next.x - node.x
			if under_width+visited_width > width {
				under_width = width - visited_width
			}
			waste_area += under_width * (min_y - node.y)
			visited_width += under_width
		}
		node = node.next
	}

	return min_y, waste_area
}

func skylineFindBestPos(ctx *Context, width, height int32) (fr findResult) {
	var (
		best_waste int32 = (1 << 30)
		best_x     int32
		best_y     int32 = (1 << 30)
	)
	var prev, best **Node
	var node, tail *Node

	// align to multiple of c.align
	width = (width + ctx.align - 1)
	width -= width % ctx.align

	if width%ctx.align != 0 {
		panic("width % ctx.align != 0")
	}

	// if it can't possibly fit, bail immediately
	if width > ctx.width || height > ctx.height {
		return fr
	}

	node = ctx.active_head
	prev = &ctx.active_head
	for node.x+width <= ctx.width {
		y, waste := skylineFindMinY(ctx, node, node.x, width)
		if ctx.heuristic == SkylineBLSortHeight { // actually just want to test BL
			// bottom left
			if y < best_y {
				best_y = y
				best = prev
			}
		} else {
			// best-fit
			if y+height <= ctx.height {
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

	if best != nil {
		best_x = (*best).x
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

	if ctx.heuristic == SkylineBFSortHeight {
		tail = ctx.active_head
		node = ctx.active_head
		prev = &ctx.active_head
		// find first node that's admissible
		for tail.x < width {
			tail = tail.next
		}
		for tail != nil {
			var xpos = tail.x - width

			if xpos < 0 {
				panic("xpos < 0")
			}

			// find the left position that matches this
			for node.next.x <= xpos {
				prev = &node.next
				node = node.next
			}

			if !(node.next.x > xpos && node.x <= xpos) {
				panic("!(node.next.x > xpos && node.x <= xpos)")
			}

			y, waste := skylineFindMinY(ctx, node, xpos, width)
			if y+height <= ctx.height {
				if y <= best_y {
					if y < best_y || waste < best_waste || (waste == best_waste && xpos < best_x) {
						best_x = xpos

						if !(y <= best_y) {
							panic("!(y <= best_y)")
						}

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

func skylinePackRectangle(ctx *Context, width, height int32) findResult {
	res := skylineFindBestPos(ctx, width, height)
	var node, cur *Node

	// bail if:
	//    1. it failed
	//    2. the best node doesn't fit (we don't always check this)
	//    3. we're out of memory
	if res.prev_link == nil || res.y+height > ctx.height || ctx.free_head == nil {
		res.prev_link = nil
		return res
	}

	// on success, create new node
	node = ctx.free_head
	node.x = res.x
	node.y = res.y + height

	ctx.free_head = node.next

	// insert the new node into the right starting point, and
	// let 'cur' point to the remaining nodes needing to be
	// stiched back in

	cur = *res.prev_link
	if cur.x < res.x {
		// preserve the existing one, so start testing with the next one
		var next = cur.next
		cur.next = node
		cur = next
	} else {
		*res.prev_link = node
	}

	// from here, traverse cur and free the nodes, until we get to one
	// that shouldn't be freed
	for cur.next != nil && cur.next.x <= res.x+width {
		var next = cur.next
		// move the current node to the free list
		cur.next = ctx.free_head
		ctx.free_head = cur
		cur = next
	}

	// stitch the list back in
	node.next = cur

	if cur.x < res.x+width {
		cur.x = res.x + width
	}

	if debug {
		cur = ctx.active_head
		for cur.x < ctx.width {
			if !(cur.x < cur.next.x) {
				panic("!(cur.x < cur.next.x)")
			}
			cur = cur.next
		}

		if !(cur.next == nil) {
			panic("!(cur.next == NULL)")
		}

		{
			var count int32 = 0
			cur = ctx.active_head
			for cur != nil {
				cur = cur.next
				count++
			}
			cur = ctx.free_head
			for cur != nil {
				cur = cur.next
				count++
			}

			if !(count == ctx.num_nodes+2) {
				panic("!(count == ctx.num_nodes+2)")
			}
		}
	}

	return res
}
